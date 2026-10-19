package openai_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type errorResponseTransport struct{}

func (errorResponseTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusBadRequest,
		Status:     "400 synthetic-secret-status",
		Header:     http.Header{"Content-Type": {"application/json"}, "X-Request-Id": {"synthetic-secret-request-id"}},
		Body:       io.NopCloser(strings.NewReader(`{"error":{"message":"synthetic-secret-message","code":"synthetic-secret-code","type":"synthetic-secret-type","param":"synthetic-secret-param"}}`)),
		Request:    req,
	}, nil
}

func sensitiveAPIError(t *testing.T) error {
	t.Helper()
	clearOpenAIEnvironment(t)
	client := openai.NewClient(
		option.WithAPIKey("synthetic-secret-api-key"),
		option.WithBaseURL("https://synthetic-secret-user:synthetic-secret-password@synthetic-secret-host.invalid/synthetic-secret-path?api_key=synthetic-secret-query#synthetic-secret-fragment"),
		option.WithHTTPClient(&http.Client{Transport: errorResponseTransport{}}),
		option.WithMaxRetries(0),
		option.WithQuery("api_key", "synthetic-secret-query"),
	)
	_, err := client.Models.List(context.Background())
	var apiErr *openai.Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected API error, got %T", err)
	}
	return err
}

func TestAPIErrorImplicitDiagnostics(t *testing.T) {
	err := sensitiveAPIError(t)
	wrapped := fmt.Errorf("list models: %w", err)
	var apiErr *openai.Error
	if !errors.As(err, &apiErr) {
		t.Fatal("expected typed API error")
	}
	for name, value := range map[string]string{
		"Error": err.Error(), "string": fmt.Sprintf("%s", err), "value": fmt.Sprintf("%v", err),
		"detailed": fmt.Sprintf("%+v", err), "Go syntax": fmt.Sprintf("%#v", err), "quoted": fmt.Sprintf("%q", err), "wrapped": wrapped.Error(),
		"copied value": fmt.Sprintf("%v", *apiErr), "copied detailed": fmt.Sprintf("%+v", *apiErr),
		"copied Go syntax": fmt.Sprintf("%#v", *apiErr), "copied slice": fmt.Sprintf("%#v", []openai.Error{*apiErr}),
	} {
		t.Run(name, func(t *testing.T) {
			if strings.Contains(value, "synthetic-secret") {
				t.Fatalf("implicit diagnostics exposed sensitive data: %s", value)
			}
			if !strings.Contains(value, "400 Bad Request") {
				t.Fatalf("missing safe status: %s", value)
			}
		})
	}
	for name, value := range map[string]any{
		"pointer": apiErr, "wrapped": wrapped, "copied": *apiErr,
		"without transport metadata": openai.Error{StatusCode: 400, Message: "synthetic-secret-message"},
	} {
		t.Run("logging/"+name, func(t *testing.T) {
			for _, jsonHandler := range []bool{false, true} {
				var buf bytes.Buffer
				var handler slog.Handler = slog.NewTextHandler(&buf, nil)
				if jsonHandler {
					handler = slog.NewJSONHandler(&buf, nil)
				}
				slog.New(handler).Error("request failed", "err", value)
				if strings.Contains(buf.String(), "synthetic-secret") {
					t.Fatalf("structured logging exposed sensitive data: %s", buf.String())
				}
				if !strings.Contains(buf.String(), "400 Bad Request") {
					t.Fatalf("structured logging lost safe status: %s", buf.String())
				}
			}
		})
	}
	if !errors.As(wrapped, &apiErr) || apiErr.StatusCode != 400 {
		t.Fatal("wrapped API error lost its type or status")
	}
	if apiErr.Message != "synthetic-secret-message" || apiErr.Code != "synthetic-secret-code" || apiErr.Type != "synthetic-secret-type" || apiErr.Param != "synthetic-secret-param" {
		t.Fatal("structured diagnostics changed")
	}
	if !strings.Contains(apiErr.RawJSON(), "synthetic-secret-message") {
		t.Fatal("raw JSON unavailable")
	}
	if !strings.Contains(string(apiErr.DumpRequest(false)), "synthetic-secret-query") {
		t.Fatal("explicit request diagnostics unavailable")
	}
	if !strings.Contains(string(apiErr.DumpResponse(true)), "synthetic-secret-message") {
		t.Fatal("explicit response diagnostics unavailable")
	}
}

func TestAPIErrorFormatVerbs(t *testing.T) {
	var apiErr *openai.Error
	if !errors.As(sensitiveAPIError(t), &apiErr) {
		t.Fatal("expected typed API error")
	}
	summary := apiErr.Error()
	for name, value := range map[string]any{"pointer": apiErr, "copied": *apiErr} {
		for _, format := range []string{"%v", "%+v", "%#v", "%s", "%q", "%x", "%X", "%40s", "%-40.20s", "%d", "%f", "%t", "%c", "%U", "%b", "%o", "%e", "%g", "%a"} {
			t.Run(name+"/"+format, func(t *testing.T) {
				want := summary
				switch format[len(format)-1] {
				case 's', 'q', 'x', 'X':
					want = fmt.Sprintf(format, summary)
				}
				if got := fmt.Sprintf(format, value); got != want {
					t.Fatalf("formatted error got %q, want %q", got, want)
				}
				var output bytes.Buffer
				log.New(&output, "", 0).Printf(format, value)
				if got := output.String(); got != want+"\n" {
					t.Fatalf("logged error got %q, want %q", got, want+"\n")
				}
			})
		}
	}
}

func TestAPIErrorPanicDiagnostics(t *testing.T) {
	if os.Getenv("OPENAI_GO_TEST_ERROR_PANIC") == "1" {
		panic(sensitiveAPIError(t))
	}
	cmd := exec.Command(os.Args[0], "-test.run=^TestAPIErrorPanicDiagnostics$")
	cmd.Env = append(os.Environ(), "OPENAI_GO_TEST_ERROR_PANIC=1")
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected child panic")
	}
	if !bytes.Contains(output, []byte("panic:")) || !bytes.Contains(output, []byte("400 Bad Request")) {
		t.Fatalf("missing API panic: %s", output)
	}
	if bytes.Contains(output, []byte("synthetic-secret")) {
		t.Fatalf("panic exposed sensitive data: %s", output)
	}
}

func TestAPIErrorStatusSummary(t *testing.T) {
	for _, tc := range []struct {
		name string
		err  *openai.Error
		want string
	}{
		{"nil", nil, "OpenAI API error"},
		{"zero", &openai.Error{}, "OpenAI API error: 0"},
		{"status without request", &openai.Error{StatusCode: 401}, "OpenAI API error: 401 Unauthorized"},
		{"unknown status", &openai.Error{StatusCode: 499, Response: &http.Response{Status: "synthetic-secret"}}, "OpenAI API error: 499"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.err.Error(); got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
			formattedWant := tc.want
			if tc.err == nil {
				formattedWant = "<nil>"
			}
			for _, format := range []string{"%v", "%#v", "%d"} {
				if got := fmt.Sprintf(format, tc.err); got != formattedWant {
					t.Fatalf("format %s got %q, want %q", format, got, formattedWant)
				}
			}
		})
	}
}

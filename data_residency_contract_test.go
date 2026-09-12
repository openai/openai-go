// File generated from our OpenAPI spec by Castiron. See CONTRIBUTING.md for details.

package openai_test

import (
	"context"
	sdk "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"slices"
	"strings"
	"testing"
)

type castironResidencyRecorder struct {
	urls    []string
	bodies  []string
	headers []http.Header
}

func (r *castironResidencyRecorder) Do(req *http.Request) (*http.Response, error) {
	body := ""
	if req.Body != nil {
		content, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		body = string(content)
	}
	r.urls = append(r.urls, req.URL.String())
	r.bodies = append(r.bodies, body)
	r.headers = append(r.headers, req.Header.Clone())
	return &http.Response{StatusCode: http.StatusOK, Header: http.Header{"Content-Type": {"application/json"}}, Body: io.NopCloser(strings.NewReader("{}")), Request: req}, nil
}

func castironResidencyOptions(t *testing.T, recorder *castironResidencyRecorder) []option.RequestOption {
	t.Helper()
	for _, key := range []string{"OPENAI_API_KEY", "OPENAI_ADMIN_KEY", "OPENAI_BASE_URL", "OPENAI_CUSTOM_HEADERS", "OPENAI_ORG_ID", "OPENAI_PROJECT_ID", "OPENAI_WEBHOOK_SECRET"} {
		t.Setenv(key, "")
	}
	return []option.RequestOption{option.WithAPIKey("dummy-key"), option.WithHTTPClient(recorder), option.WithMaxRetries(0)}
}

func TestCastironDataResidencyMappings(t *testing.T) {
	for _, test := range []struct {
		region   option.DataResidency
		endpoint string
	}{
		{option.DataResidencyGlobal, "https://api.openai.com/v1/"},
		{option.DataResidencyUS, "https://us.api.openai.com/v1/"},
		{option.DataResidencyEU, "https://eu.api.openai.com/v1/"},
		{option.DataResidencyAE, "https://ae.api.openai.com/v1/"},
	} {
		t.Run(string(test.region), func(t *testing.T) {
			endpoint, err := url.Parse(test.endpoint)
			if err != nil {
				t.Fatalf("invalid configured residency URL %q: %v", test.endpoint, err)
			}
			if (endpoint.Scheme != "http" && endpoint.Scheme != "https") || endpoint.Host == "" {
				t.Fatalf("configured residency URL must be absolute HTTP(S): %q", test.endpoint)
			}
			recorder := &castironResidencyRecorder{}
			opts := castironResidencyOptions(t, recorder)
			client := sdk.NewClient(append(opts, option.WithDataResidency(test.region))...)
			if err := client.Post(context.Background(), "responses", map[string]string{"input": "hello"}, nil, requestconfig.WithDefaultBaseURL("https://resource.example/v1")); err != nil {
				t.Fatal(err)
			}
			if got, want := recorder.urls[0], test.endpoint+"responses"; got != want {
				t.Fatalf("URL = %q, want %q", got, want)
			}
			if strings.TrimSpace(recorder.bodies[0]) != `{"input":"hello"}` {
				t.Fatalf("unexpected wire body: %q", recorder.bodies[0])
			}
			for header := range recorder.headers[0] {
				if strings.Contains(strings.ToLower(header), "residency") {
					t.Fatalf("unexpected residency header %q", header)
				}
			}
		})
	}
}

func TestCastironDataResidencyInvalid(t *testing.T) {
	recorder := &castironResidencyRecorder{}
	client := sdk.NewClient(castironResidencyOptions(t, recorder)...)
	for _, value := range []option.DataResidency{"", "UNKNOWN", "not-a-region"} {
		err := client.Get(context.Background(), "models", nil, nil, option.WithDataResidency(value))
		if err == nil || !strings.Contains(err.Error(), "invalid data residency") {
			t.Fatalf("expected invalid residency, got %v", err)
		}
	}
	if len(recorder.urls) != 0 {
		t.Fatal("invalid options made an HTTP request")
	}
}

func TestCastironDataResidencyConflicts(t *testing.T) {
	for _, reverse := range []bool{false, true} {
		for _, layer := range []string{"client", "request"} {
			recorder := &castironResidencyRecorder{}
			opts := castironResidencyOptions(t, recorder)
			conflict := []option.RequestOption{option.WithBaseURL("https://custom.example/v1"), option.WithDataResidency(option.DataResidencyGlobal)}
			if reverse {
				slices.Reverse(conflict)
			}
			var err error
			if layer == "client" {
				client := sdk.NewClient(append(opts, conflict...)...)
				err = client.Get(context.Background(), "models", nil, nil)
			} else {
				client := sdk.NewClient(opts...)
				err = client.Get(context.Background(), "models", nil, nil, conflict...)
			}
			if err == nil || !strings.Contains(err.Error(), "mutually exclusive") {
				t.Fatalf("%s reverse=%t: expected conflict, got %v", layer, reverse, err)
			}
			if len(recorder.urls) != 0 {
				t.Fatal("conflicting options made an HTTP request")
			}
		}
	}
}

func TestCastironDataResidencyInheritedOverrides(t *testing.T) {
	recorder := &castironResidencyRecorder{}
	opts := castironResidencyOptions(t, recorder)
	t.Setenv("OPENAI_BASE_URL", "https://environment.example/v1")
	environment := sdk.NewClient(opts...)
	if err := environment.Get(context.Background(), "omitted", nil, nil); err != nil {
		t.Fatal(err)
	}
	client := sdk.NewClient(append(opts, option.WithDataResidency(option.DataResidencyGlobal))...)
	if err := client.Get(context.Background(), "original", nil, nil); err != nil {
		t.Fatal(err)
	}
	if err := client.Get(context.Background(), "request", nil, nil, option.WithBaseURL("https://request.example/v1")); err != nil {
		t.Fatal(err)
	}
	copy := sdk.NewClient(append(slices.Clone(client.Options), option.WithBaseURL("https://copy.example/v1"))...)
	if err := copy.Get(context.Background(), "copy", nil, nil); err != nil {
		t.Fatal(err)
	}
	selected := sdk.NewClient(append(slices.Clone(copy.Options), option.WithDataResidency(option.DataResidencyGlobal))...)
	if err := selected.Get(context.Background(), "selected", nil, nil); err != nil {
		t.Fatal(err)
	}
	if err := client.Get(context.Background(), "unchanged", nil, nil); err != nil {
		t.Fatal(err)
	}
	want := []string{"https://environment.example/v1/omitted", "https://api.openai.com/v1/" + "original", "https://request.example/v1/request", "https://copy.example/v1/copy", "https://api.openai.com/v1/" + "selected", "https://api.openai.com/v1/" + "unchanged"}
	if !reflect.DeepEqual(recorder.urls, want) {
		t.Fatalf("URLs = %#v, want %#v", recorder.urls, want)
	}
}

package openai_test

import (
	"context"
	"io"
	"net/http"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/azure"
	"github.com/openai/openai-go/v3/option"
)

type residencyRecorder struct {
	urls    []string
	bodies  []string
	headers []http.Header
}

func (r *residencyRecorder) Do(req *http.Request) (*http.Response, error) {
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
	return &http.Response{StatusCode: http.StatusOK, Header: http.Header{"Content-Type": {"application/json"}}, Body: io.NopCloser(strings.NewReader(`{"id":"test","object":"model"}`)), Request: req}, nil
}

func residencyTestOptions(t *testing.T, recorder *residencyRecorder) []option.RequestOption {
	t.Helper()
	for _, key := range []string{"OPENAI_API_KEY", "OPENAI_ADMIN_KEY", "OPENAI_BASE_URL", "OPENAI_CUSTOM_HEADERS", "OPENAI_ORG_ID", "OPENAI_PROJECT_ID"} {
		t.Setenv(key, "")
	}
	return []option.RequestOption{option.WithAPIKey("dummy-key"), option.WithHTTPClient(recorder), option.WithMaxRetries(0)}
}

func TestDataResidencyMappings(t *testing.T) {
	for _, test := range []struct {
		region option.DataResidency
		host   string
	}{
		{option.DataResidencyGlobal, "api.openai.com"},
		{option.DataResidencyUS, "us.api.openai.com"},
		{option.DataResidencyEU, "eu.api.openai.com"},
		{option.DataResidencyAE, "ae.api.openai.com"},
	} {
		t.Run(string(test.region), func(t *testing.T) {
			recorder := &residencyRecorder{}
			opts := residencyTestOptions(t, recorder)
			client := openai.NewClient(append(opts, option.WithDataResidency(test.region))...)
			if err := client.Post(context.Background(), "responses", map[string]string{"input": "hello"}, nil); err != nil {
				t.Fatal(err)
			}
			if want := "https://" + test.host + "/v1/responses"; recorder.urls[0] != want {
				t.Fatalf("URL = %q, want %q", recorder.urls[0], want)
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

func TestDataResidencyOptionLayers(t *testing.T) {
	recorder := &residencyRecorder{}
	opts := residencyTestOptions(t, recorder)
	t.Setenv("OPENAI_BASE_URL", "https://environment.example/v1")
	client := openai.NewClient(append(opts, option.WithDataResidency(option.DataResidencyEU))...)
	if _, err := client.Models.Get(context.Background(), "client"); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Models.Get(context.Background(), "request", option.WithBaseURL("https://request.example/v1")); err != nil {
		t.Fatal(err)
	}
	service := openai.NewModelService(append(slices.Clone(client.Options), option.WithBaseURL("https://service.example/v1"))...)
	if _, err := service.Get(context.Background(), "service"); err != nil {
		t.Fatal(err)
	}
	if _, err := service.Get(context.Background(), "global", option.WithDataResidency(option.DataResidencyGlobal)); err != nil {
		t.Fatal(err)
	}
	inherited := openai.NewClient(append(slices.Clone(client.Options), option.WithBaseURL("https://copy.example/v1"))...)
	if err := inherited.Get(context.Background(), "copy", nil, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Models.Get(context.Background(), "unchanged"); err != nil {
		t.Fatal(err)
	}
	want := []string{"https://eu.api.openai.com/v1/models/client", "https://request.example/v1/models/request", "https://service.example/v1/models/service", "https://api.openai.com/v1/models/global", "https://copy.example/v1/copy", "https://eu.api.openai.com/v1/models/unchanged"}
	if !reflect.DeepEqual(recorder.urls, want) {
		t.Fatalf("URLs = %#v, want %#v", recorder.urls, want)
	}
}

func TestDataResidencyConflicts(t *testing.T) {
	for _, reverse := range []bool{false, true} {
		for _, layer := range []string{"client", "service", "request", "inherited service"} {
			t.Run(layer+map[bool]string{false: "/base-first", true: "/residency-first"}[reverse], func(t *testing.T) {
				recorder := &residencyRecorder{}
				opts := residencyTestOptions(t, recorder)
				conflict := []option.RequestOption{option.WithBaseURL("https://custom.example/v1"), option.WithDataResidency(option.DataResidencyEU)}
				if reverse {
					slices.Reverse(conflict)
				}
				var err error
				switch layer {
				case "client":
					client := openai.NewClient(append(opts, conflict...)...)
					err = client.Get(context.Background(), "models", nil, nil)
				case "service":
					service := openai.NewModelService(append(opts, conflict...)...)
					_, err = service.Get(context.Background(), "test")
				case "request":
					client := openai.NewClient(opts...)
					_, err = client.Models.Get(context.Background(), "test", conflict...)
				case "inherited service":
					client := openai.NewClient(opts...)
					service := openai.NewModelService(append(slices.Clone(client.Options), conflict...)...)
					_, err = service.Get(context.Background(), "test", option.WithDataResidency(option.DataResidencyUS))
				}
				if err == nil || !strings.Contains(err.Error(), "mutually exclusive") {
					t.Fatalf("expected conflict, got %v", err)
				}
				if len(recorder.urls) != 0 {
					t.Fatal("conflicting options made an HTTP request")
				}
			})
		}
	}
}

func TestDataResidencyInvalid(t *testing.T) {
	recorder := &residencyRecorder{}
	client := openai.NewClient(residencyTestOptions(t, recorder)...)
	for _, value := range []option.DataResidency{"", "EU", "unknown"} {
		err := client.Get(context.Background(), "models", nil, nil, option.WithDataResidency(value))
		if err == nil || !strings.Contains(err.Error(), "invalid data residency") {
			t.Fatalf("expected invalid residency, got %v", err)
		}
	}
	if len(recorder.urls) != 0 {
		t.Fatal("invalid options made an HTTP request")
	}
}

func TestDataResidencyProviders(t *testing.T) {
	for _, provider := range []struct {
		name string
		opt  option.RequestOption
	}{
		{"Azure endpoint", azure.WithEndpoint("https://azure.example", "2024-06-01")},
		{"Azure API key", azure.WithAPIKey("dummy-azure")},
	} {
		for _, reverse := range []bool{false, true} {
			t.Run(provider.name+map[bool]string{false: "/provider-first", true: "/residency-first"}[reverse], func(t *testing.T) {
				recorder := &residencyRecorder{}
				opts := residencyTestOptions(t, recorder)
				conflict := []option.RequestOption{provider.opt, option.WithDataResidency(option.DataResidencyEU)}
				if reverse {
					slices.Reverse(conflict)
				}
				client := openai.NewClient(append(opts, conflict...)...)
				err := client.Get(context.Background(), "models", nil, nil)
				if err == nil || !strings.Contains(err.Error(), "provider") {
					t.Fatalf("expected provider error, got %v", err)
				}
				if len(recorder.urls) != 0 {
					t.Fatal("provider conflict made an HTTP request")
				}
			})
		}
	}
}

package azure

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func TestEndpointPathRouting(t *testing.T) {
	for _, tc := range []struct {
		name, endpoint, override, route, want string
		multipart, retry                      bool
	}{
		{name: "double leading slash", endpoint: "https://azure.example//apim", route: "embeddings", want: "//apim/openai/deployments/model/embeddings"},
		{name: "chat", endpoint: "https://azure.example/apim", route: "chat/completions", want: "/apim/openai/deployments/model/chat/completions"},
		{name: "embeddings trailing slash", endpoint: "https://azure.example/apim/", route: "embeddings", want: "/apim/openai/deployments/model/embeddings"},
		{name: "nested prefix", endpoint: "https://azure.example/tenant/service", route: "embeddings", want: "/tenant/service/openai/deployments/model/embeddings"},
		{name: "escaped prefix", endpoint: "https://azure.example/tenant%2Fservice", route: "embeddings", want: "/tenant%2Fservice/openai/deployments/model/embeddings"},
		{name: "escaped path parameter", endpoint: "https://azure.example/apim", route: "vector_stores/id%2Fpart%3Fq%23frag", want: "/apim/openai/vector_stores/id%2Fpart%3Fq%23frag"},
		{name: "non-deployment route", endpoint: "https://azure.example/apim", route: "models", want: "/apim/openai/models"},
		{name: "no prefix", endpoint: "https://azure.example", route: "embeddings", want: "/openai/deployments/model/embeddings"},
		{name: "root override", endpoint: "https://azure.example/apim", override: "https://proxy.example/", route: "chat/completions", want: "/openai/deployments/model/chat/completions"},
		{name: "prefixed override", endpoint: "https://azure.example/apim", override: "https://proxy.example/other/", route: "embeddings", want: "/other/openai/deployments/model/embeddings"},
		{name: "multipart", endpoint: "https://azure.example/apim", route: "audio/transcriptions", want: "/apim/openai/deployments/model/audio/transcriptions", multipart: true},
		{name: "retry", endpoint: "https://azure.example/apim", route: "embeddings", want: "/apim/openai/deployments/model/embeddings", retry: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			body := []byte(`{"model":"model","input":"hello"}`)
			contentType := "application/json"
			if tc.multipart {
				var b bytes.Buffer
				w := multipart.NewWriter(&b)
				if err := w.WriteField("model", "model"); err != nil {
					t.Fatal(err)
				}
				if err := w.Close(); err != nil {
					t.Fatal(err)
				}
				body, contentType = b.Bytes(), w.FormDataContentType()
			}
			attempts := 0
			client := openai.NewClient(WithEndpoint(tc.endpoint, "2024-10-21"), WithAPIKey("synthetic"), option.WithMaxRetries(1), option.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				attempts++
				if got := req.URL.EscapedPath(); got != tc.want {
					t.Errorf("path = %q, want %q", got, tc.want)
				}
				if got := req.URL.Query().Get("api-version"); got != "2024-10-21" {
					t.Errorf("api-version = %q", got)
				}
				if got := req.Header.Get("Api-Key"); got != "synthetic" {
					t.Error("Azure credential missing")
				}
				if req.Header.Get("Authorization") != "" {
					t.Error("unexpected authorization header")
				}
				got, err := io.ReadAll(req.Body)
				if err != nil {
					t.Fatal(err)
				}
				if !bytes.Equal(got, body) {
					t.Error("request body changed")
				}
				status := http.StatusOK
				if tc.retry && attempts == 1 {
					status = http.StatusTooManyRequests
				}
				return &http.Response{StatusCode: status, Header: http.Header{"Content-Type": []string{"application/json"}, "Retry-After-Ms": []string{"1"}}, Body: io.NopCloser(strings.NewReader(`{}`)), Request: req}, nil
			})}))
			opts := []option.RequestOption{option.WithRequestBody(contentType, body)}
			if tc.override != "" {
				opts = append(opts, option.WithBaseURL(tc.override))
			}
			if err := client.Post(context.Background(), tc.route, nil, nil, opts...); err != nil {
				t.Fatal(err)
			}
			wantAttempts := 1
			if tc.retry {
				wantAttempts = 2
			}
			if attempts != wantAttempts {
				t.Fatalf("attempts = %d, want %d", attempts, wantAttempts)
			}
		})
	}
}

func TestEndpointPathInvalidURL(t *testing.T) {
	for _, endpoint := range []string{"https://azure.example/%zz", "https://[invalid"} {
		t.Run(endpoint, func(t *testing.T) {
			client := openai.NewClient(WithEndpoint(endpoint, "2024-10-21"), WithAPIKey("synthetic"))
			err := client.Post(context.Background(), "embeddings", map[string]string{"model": "model"}, nil)
			if err == nil || !strings.Contains(err.Error(), "failed to parse url") {
				t.Fatalf("expected URL parse error, got %v", err)
			}
		})
	}
}

func TestEndpointPathMiddlewareOrder(t *testing.T) {
	var paths []string
	observe := option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
		paths = append(paths, req.URL.EscapedPath())
		return next(req)
	})
	client := openai.NewClient(observe, WithEndpoint("https://azure.example/apim", "2024-10-21"), WithAPIKey("synthetic"), observe, option.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": []string{"application/json"}}, Body: io.NopCloser(strings.NewReader(`{}`)), Request: req}, nil
	})}))
	if err := client.Post(context.Background(), "embeddings", map[string]string{"model": "model"}, nil); err != nil {
		t.Fatal(err)
	}
	if len(paths) != 2 || paths[0] != "/apim/embeddings" || paths[1] != "/apim/openai/deployments/model/embeddings" {
		t.Fatalf("unexpected middleware paths: %v", paths)
	}
}

func TestEndpointPathConcurrentOverrides(t *testing.T) {
	client := openai.NewClient(WithEndpoint("https://azure.example/original", "2024-10-21"), WithAPIKey("synthetic"), option.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		want := "/" + req.Header.Get("X-Test-Prefix") + "/openai/deployments/model/embeddings"
		if req.URL.EscapedPath() != want {
			t.Errorf("path = %q, want %q", req.URL.EscapedPath(), want)
		}
		return &http.Response{StatusCode: http.StatusOK, Header: http.Header{"Content-Type": []string{"application/json"}}, Body: io.NopCloser(strings.NewReader(`{}`)), Request: req}, nil
	})}))
	for _, prefix := range []string{"first", "second", "third"} {
		t.Run(prefix, func(t *testing.T) {
			t.Parallel()
			if err := client.Post(context.Background(), "embeddings", map[string]string{"model": "model"}, nil, option.WithBaseURL("https://azure.example/"+prefix+"/"), option.WithHeader("X-Test-Prefix", prefix)); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestEndpointPathMiddlewareReplacesPrefix(t *testing.T) {
	for _, path := range []string{"/other/embeddings", "/apim-other/embeddings"} {
		t.Run(path, func(t *testing.T) {
			client := openai.NewClient(option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
				req.URL.Path = path
				return next(req)
			}), WithEndpoint("https://azure.example/apim", "2024-10-21"), WithAPIKey("synthetic"), option.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				if got, want := req.URL.EscapedPath(), "/openai"+path; got != want {
					t.Errorf("path = %q, want %q", got, want)
				}
				return &http.Response{StatusCode: http.StatusOK, Header: http.Header{"Content-Type": []string{"application/json"}}, Body: io.NopCloser(strings.NewReader(`{}`)), Request: req}, nil
			})}))
			if err := client.Post(context.Background(), "embeddings", map[string]string{"model": "model"}, nil); err != nil {
				t.Fatal(err)
			}
		})
	}
}

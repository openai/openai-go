package bedrock

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/openai/openai-go/v3/option"
)

func TestSkipAuthLoopbackRedirectReturnsError(t *testing.T) {
	clearAWSEnvironment(t)
	for _, key := range []struct {
		name   string
		option option.RequestOption
	}{
		{"api", option.WithAPIKey("test-gateway")},
		{"admin", option.WithAdminAPIKey("test-gateway")},
	} {
		for _, level := range []string{"client", "request"} {
			t.Run(key.name+"/"+level, func(t *testing.T) {
				var targetCalls, redirectCalls, successCalls atomic.Int32
				target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					targetCalls.Add(1)
					w.WriteHeader(http.StatusNoContent)
				}))
				defer target.Close()
				source := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					if req.Header.Get("Authorization") != "Bearer test-gateway" {
						t.Error("missing gateway authentication")
					}
					if req.URL.Path == "/ok" {
						successCalls.Add(1)
						w.Header().Set("Content-Type", "application/json")
						_, _ = w.Write([]byte(`{}`))
						return
					}
					redirectCalls.Add(1)
					w.Header().Set("X-Should-Retry", "true")
					w.Header().Set("Retry-After-Ms", "1")
					http.Redirect(w, req, target.URL, http.StatusTemporaryRedirect)
				}))
				defer source.Close()

				opts := []option.RequestOption{option.WithMaxRetries(2)}
				var requestOpts []option.RequestOption
				if level == "client" {
					opts = append(opts, key.option)
				} else {
					requestOpts = append(requestOpts, key.option)
				}
				client, err := NewClient(t.Context(), Config{
					SkipAuth:        true,
					BaseURL:         source.URL,
					UnsafeAllowHTTP: true,
				}, opts...)
				if err != nil {
					t.Fatal(err)
				}
				body := map[string]string{"input": "synthetic prompt"}
				var result map[string]any
				if requestErr := client.Post(t.Context(), "/ok", body, &result, requestOpts...); requestErr != nil {
					t.Fatalf("direct loopback request failed: %v", requestErr)
				}
				if successCalls.Load() != 1 {
					t.Fatalf("direct loopback calls = %d, want 1", successCalls.Load())
				}

				var response *http.Response
				requestOpts = append(requestOpts, option.WithResponseInto(&response))
				err = client.Post(t.Context(), "/redirect", body, nil, requestOpts...)
				if response != nil && response.Body != nil {
					if closeErr := response.Body.Close(); closeErr != nil {
						t.Errorf("close response body: %v", closeErr)
					}
				}
				if err == nil || !strings.Contains(err.Error(), "authenticated loopback redirects are not allowed") {
					t.Errorf("redirect error = %v, want authenticated loopback redirect rejection", err)
				}
				if calls := redirectCalls.Load(); calls != 1 {
					t.Errorf("redirect source calls = %d, want 1 without retries", calls)
				}
				if calls := targetCalls.Load(); calls != 0 {
					t.Errorf("redirect target calls = %d, want 0", calls)
				}
			})
		}
	}
}

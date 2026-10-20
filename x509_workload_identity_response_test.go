package openai_test

import (
	"errors"
	"io"
	"net/http"
	"sync/atomic"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func TestX509WorkloadIdentityRedactsReturnedResponseRequests(t *testing.T) {
	for _, status := range []int{http.StatusOK, http.StatusUnauthorized} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			t.Setenv("OPENAI_BASE_URL", "https://mtls.api.openai.com/v1/")
			config, issuer, api := newX509WorkloadIdentityIntegration(t)
			originalHandler := api.server.Config.Handler
			var wireRequests, observedResponses atomic.Int32
			api.server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
				wireRequests.Add(1)
				if got := request.Header.Get("Authorization"); got != "Bearer "+x509ConformanceToken {
					t.Errorf("mutually authenticated API wire authorization = %q", got)
				}
				if status == http.StatusUnauthorized {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(status)
					_, _ = io.WriteString(w, `{"error":{"message":"synthetic unauthorized","type":"invalid_request_error"}}`)
					return
				}
				originalHandler.ServeHTTP(w, request)
			})

			var returned *http.Response
			client := openai.NewClient(
				option.WithX509WorkloadIdentity(config),
				option.WithMaxRetries(0),
				option.WithMiddleware(func(request *http.Request, next option.MiddlewareNext) (*http.Response, error) {
					response, err := next(request)
					if response == nil || response.Request == nil {
						t.Error("outer middleware did not receive its response request metadata")
					} else if got := response.Request.Header.Get("Authorization"); got != "" {
						t.Errorf("outer middleware observed a workload bearer in response metadata: %q", got)
					}
					observedResponses.Add(1)
					return response, err
				}),
			)
			_, err := client.Models.List(t.Context(), option.WithResponseInto(&returned))
			if status == http.StatusOK && err != nil {
				t.Fatalf("successful mutually authenticated API request: %v", err)
			}
			if status == http.StatusUnauthorized {
				var apiError *openai.Error
				if !errors.As(err, &apiError) || apiError.StatusCode != status {
					t.Fatalf("unauthorized mutually authenticated API request = %v", err)
				}
				if apiError.Response == nil || apiError.Response.Request == nil ||
					apiError.Response.Request.Header.Get("Authorization") != "" {
					t.Error("typed API error retained workload credentials in its response metadata")
				}
			}
			if returned == nil || returned.Request == nil || returned.Request.URL.Host != x509ConformanceAPIHost {
				t.Fatal("WithResponseInto did not preserve the response request and destination")
			}
			if got := returned.Request.Header.Get("Authorization"); got != "" {
				t.Errorf("WithResponseInto retained a workload bearer in response metadata: %q", got)
			}
			if wireRequests.Load() != 1 || observedResponses.Load() != 1 || len(issuer.requests()) != 1 {
				t.Errorf("API/middleware/issuer calls = %d/%d/%d, want 1/1/1",
					wireRequests.Load(), observedResponses.Load(), len(issuer.requests()))
			}
		})
	}
}

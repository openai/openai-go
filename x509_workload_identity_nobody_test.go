package openai_test

import (
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func TestX509WorkloadIdentityReplaysConfiguredHTTPNoBody(t *testing.T) {
	t.Setenv("OPENAI_BASE_URL", "https://mtls.api.openai.com/v1/")
	config, issuer, api := newX509WorkloadIdentityIntegration(t)
	var attempts, observations, exchanges atomic.Int32
	issuer.server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, x509IntegrationTokenResponse(
			fmt.Sprintf("synthetic-configured-empty-%d", exchanges.Add(1)),
		))
	})
	api.server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		if request.ContentLength != 0 {
			t.Errorf("configured empty request content length = %d, want zero", request.ContentLength)
		}
		w.Header().Set("Content-Type", "application/json")
		attempt := attempts.Add(1)
		if got := request.Header.Get("Authorization"); got != fmt.Sprintf("Bearer synthetic-configured-empty-%d", attempt) {
			t.Errorf("configured empty request attempt %d authorization = %q", attempt, got)
		}
		if attempt == 1 {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = io.WriteString(w, `{"error":{"message":"synthetic rejected bearer"}}`)
			return
		}
		_, _ = io.WriteString(w, `{"data":[]}`)
	})
	client := openai.NewClient(
		option.WithX509WorkloadIdentity(config),
		option.WithMaxRetries(1),
		option.WithMiddleware(func(request *http.Request, next option.MiddlewareNext) (*http.Response, error) {
			if request.Body != http.NoBody {
				t.Errorf("configured empty request lost the canonical http.NoBody sentinel: %T", request.Body)
			}
			observations.Add(1)
			return next(request)
		}),
	)
	var result map[string]any
	if err := client.Execute(t.Context(), http.MethodPost, "models", nil, &result,
		option.WithRequestBody("application/json", http.NoBody)); err != nil {
		t.Fatalf("configured http.NoBody request was not replayed after authentication failure: %v", err)
	}
	if attempts.Load() != 2 || observations.Load() != 2 || exchanges.Load() != 2 {
		t.Errorf("configured empty request API/middleware/issuer attempts = %d/%d/%d, want 2/2/2",
			attempts.Load(), observations.Load(), exchanges.Load())
	}
}

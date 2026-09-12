package openai_test

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/auth"
	"github.com/openai/openai-go/v3/option"
)

func TestX509WorkloadIdentityIgnoresInactiveProviderMiddlewareForBodyReplay(t *testing.T) {
	for _, test := range []struct {
		name             string
		callerMiddleware bool
		callerFirst      bool
		repeatedProvider bool
		wantAttempts     int32
	}{
		{name: "inactive inherited provider is safe to replay", wantAttempts: 2},
		{name: "multiple inactive inherited providers are safe to replay", repeatedProvider: true, wantAttempts: 2},
		{name: "caller middleware after inherited provider remains unsafe", callerMiddleware: true, wantAttempts: 1},
		{name: "caller middleware before inherited provider remains unsafe", callerMiddleware: true, callerFirst: true, wantAttempts: 1},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("OPENAI_BASE_URL", "https://mtls.api.openai.com/v1/")
			config, issuer, api := newX509WorkloadIdentityIntegration(t)
			var exchanges, requests atomic.Int32
			var bodiesMu sync.Mutex
			var bodies []string
			issuer.server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				_, _ = io.WriteString(w, x509IntegrationTokenResponse(
					fmt.Sprintf("synthetic-provider-replay-%d", exchanges.Add(1)),
				))
			})
			api.server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
				body, err := io.ReadAll(request.Body)
				if err != nil {
					t.Errorf("read mutually authenticated API request: %v", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				bodiesMu.Lock()
				bodies = append(bodies, string(body))
				bodiesMu.Unlock()
				w.Header().Set("Content-Type", "application/json")
				if requests.Add(1) == 1 {
					w.WriteHeader(http.StatusUnauthorized)
					_, _ = io.WriteString(w, `{"error":{"message":"synthetic rejected bearer"}}`)
					return
				}
				_, _ = io.WriteString(w, `{}`)
			})

			ordinary := option.WithWorkloadIdentity(auth.WorkloadIdentity{
				IdentityProviderID: "synthetic-inactive-provider",
				ServiceAccountID:   "synthetic-inactive-service-account",
				Provider:           originTestSubjectTokenProvider{},
			})
			options := []option.RequestOption{ordinary, option.WithMaxRetries(1)}
			if test.repeatedProvider {
				options = append(options, ordinary)
			}
			if test.callerMiddleware {
				caller := option.WithMiddleware(func(request *http.Request, next option.MiddlewareNext) (*http.Response, error) {
					if err := request.Body.Close(); err != nil {
						return nil, err
					}
					const transformed = "synthetic-middleware-transformed-body"
					request.Body = io.NopCloser(strings.NewReader(transformed))
					request.ContentLength = int64(len(transformed))
					return next(request)
				})
				if test.callerFirst {
					options = append([]option.RequestOption{caller}, options...)
				} else {
					options = append(options, caller)
				}
			}
			client := openai.NewClient(options...)
			err := client.Execute(
				t.Context(),
				http.MethodPost,
				"models",
				map[string]string{"input": "synthetic-original-body"},
				nil,
				option.WithX509WorkloadIdentity(config),
			)
			if gotSuccess, wantSuccess := err == nil, test.wantAttempts == 2; gotSuccess != wantSuccess {
				t.Errorf("inherited-provider unauthorized recovery error = %v, want success %v", err, wantSuccess)
			}
			if got := exchanges.Load(); got != test.wantAttempts {
				t.Errorf("issuer exchanges = %d, want %d", got, test.wantAttempts)
			}
			if got := requests.Load(); got != test.wantAttempts {
				t.Errorf("mutually authenticated API attempts = %d, want %d", got, test.wantAttempts)
			}
			bodiesMu.Lock()
			defer bodiesMu.Unlock()
			if len(bodies) != int(test.wantAttempts) {
				t.Fatalf("captured API bodies = %d, want %d", len(bodies), test.wantAttempts)
			}
			if test.callerMiddleware {
				if bodies[0] != "synthetic-middleware-transformed-body" {
					t.Errorf("API request body = %q, want caller-transformed body", bodies[0])
				}
			} else if bodies[0] != bodies[1] || !strings.Contains(bodies[0], "synthetic-original-body") {
				t.Errorf("replayed API bodies = %q, want identical serialized JSON", bodies)
			}
		})
	}
}

func TestWorkloadIdentityPreservesCallerMiddlewareOrdering(t *testing.T) {
	for _, test := range []struct {
		name              string
		callerFirst       bool
		wantAuthorization string
	}{
		{name: "middleware before authentication remains unsigned", callerFirst: true},
		{name: "middleware after authentication observes signed request", wantAuthorization: "Bearer synthetic-ordinary-token"},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("OPENAI_API_KEY", "")
			t.Setenv("OPENAI_BASE_URL", "https://api.openai.com/v1/")
			var observedAuthorization string
			caller := option.WithMiddleware(func(request *http.Request, next option.MiddlewareNext) (*http.Response, error) {
				observedAuthorization = request.Header.Get("Authorization")
				return next(request)
			})
			ordinary := option.WithWorkloadIdentity(auth.WorkloadIdentity{
				IdentityProviderID: "synthetic-ordinary-provider",
				ServiceAccountID:   "synthetic-ordinary-service-account",
				Provider:           originTestSubjectTokenProvider{},
			})
			options := []option.RequestOption{ordinary, caller}
			if test.callerFirst {
				options[0], options[1] = options[1], options[0]
			}
			options = append(options, option.WithHTTPClient(&http.Client{
				Transport: &closureTransport{fn: func(request *http.Request) (*http.Response, error) {
					body := `{}`
					if strings.Contains(request.URL.Path, "/oauth/token") {
						body = `{"access_token":"synthetic-ordinary-token","expires_in":3600}`
					} else if got := request.Header.Get("Authorization"); got != "Bearer synthetic-ordinary-token" {
						t.Errorf("ordinary authenticated API Authorization = %q, want signed bearer", got)
					}
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(body)),
						Header:     http.Header{"Content-Type": []string{"application/json"}},
					}, nil
				}},
			}))
			client := openai.NewClient(options...)
			if err := client.Execute(t.Context(), http.MethodGet, "models", nil, nil); err != nil {
				t.Fatalf("ordinary authenticated public API request: %v", err)
			}
			if observedAuthorization != test.wantAuthorization {
				t.Errorf("caller middleware observed Authorization = %q, want %q",
					observedAuthorization, test.wantAuthorization)
			}
		})
	}
}

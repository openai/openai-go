package openai_test

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/auth"
	"github.com/openai/openai-go/v3/option"
)

func TestX509WorkloadIdentityBoundsIssuerRetriesAcrossSDKRetries(t *testing.T) {
	for _, test := range []struct {
		name        string
		failures    int32
		wantIssuer  int32
		wantAPI     int32
		wantSuccess bool
	}{
		{name: "transient issuer recovers", failures: 2, wantIssuer: 3, wantAPI: 1, wantSuccess: true},
		{name: "transient issuer exhausts", failures: 10, wantIssuer: 3, wantAPI: 0},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("OPENAI_BASE_URL", "https://mtls.api.openai.com/v1/")
			config, issuer, api := newX509WorkloadIdentityIntegration(t)
			var exchanges atomic.Int32
			issuer.server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				if exchanges.Add(1) <= test.failures {
					w.WriteHeader(http.StatusServiceUnavailable)
					return
				}
				_, _ = io.WriteString(w, x509IntegrationTokenResponse("synthetic-recovered-token"))
			})
			client := openai.NewClient(option.WithX509WorkloadIdentity(config), option.WithMaxRetryDelay(time.Millisecond))
			_, err := client.Models.List(t.Context())
			if (err == nil) != test.wantSuccess {
				t.Errorf("issuer retry result error=%v, want success=%v", err, test.wantSuccess)
			}
			if got := exchanges.Load(); got != test.wantIssuer {
				t.Errorf("issuer attempts = %d, want %d without SDK retry multiplication", got, test.wantIssuer)
			}
			if got := int32(len(api.requests())); got != test.wantAPI {
				t.Errorf("API attempts = %d, want %d", got, test.wantAPI)
			}
		})
	}
}

func TestX509WorkloadIdentitySharesRetryBudgetAcrossMixedIssuerAPIAndUnauthorized(t *testing.T) {
	for _, test := range []struct {
		name        string
		issuerFails int32
		statuses    []int
		wantIssuer  int32
		wantAPI     int32
	}{
		{name: "API 500 then 401 replay then 500", statuses: []int{500, 401, 500}, wantIssuer: 2, wantAPI: 3},
		{name: "issuer exhausts retries before 401", issuerFails: 2, statuses: []int{401}, wantIssuer: 3, wantAPI: 1},
		{name: "401 replay then 500 then second 401", statuses: []int{401, 500, 401}, wantIssuer: 2, wantAPI: 3},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("OPENAI_BASE_URL", "https://mtls.api.openai.com/v1/")
			config, issuer, api := newX509WorkloadIdentityIntegration(t)
			var exchanges, requests atomic.Int32
			issuer.server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				number := exchanges.Add(1)
				if number <= test.issuerFails {
					w.WriteHeader(http.StatusServiceUnavailable)
					return
				}
				_, _ = io.WriteString(w, x509IntegrationTokenResponse(fmt.Sprintf("synthetic-scope-%d", number)))
			})
			api.server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				number := int(requests.Add(1)) - 1
				w.Header().Set("Content-Type", "application/json")
				if number >= len(test.statuses) {
					w.WriteHeader(http.StatusInternalServerError)
					_, _ = io.WriteString(w, `{"error":{"message":"unexpected extra attempt"}}`)
					return
				}
				w.WriteHeader(test.statuses[number])
				_, _ = io.WriteString(w, `{"error":{"message":"synthetic mixed-sequence failure"}}`)
			})
			client := openai.NewClient(option.WithX509WorkloadIdentity(config), option.WithMaxRetryDelay(time.Millisecond))
			if _, err := client.Models.List(t.Context()); err == nil {
				t.Fatal("mixed issuer/API failure unexpectedly succeeded")
			}
			if exchanges.Load() != test.wantIssuer || requests.Load() != test.wantAPI {
				t.Errorf("mixed issuer/API attempts = %d/%d, want %d/%d within one logical retry budget",
					exchanges.Load(), requests.Load(), test.wantIssuer, test.wantAPI)
			}
		})
	}
}

func TestX509WorkloadIdentityRefreshesUnauthorizedReplayableRequestsOnce(t *testing.T) {
	for _, test := range []struct {
		name         string
		body         any
		alwaysReject bool
		wantIssuer   int32
		wantAPI      int32
		wantSuccess  bool
	}{
		{name: "bodyless request", wantIssuer: 2, wantAPI: 2, wantSuccess: true},
		{name: "replayable JSON body", body: map[string]string{"input": "synthetic-body"},
			wantIssuer: 2, wantAPI: 2, wantSuccess: true},
		{name: "non-replayable body", body: strings.NewReader("synthetic-streaming-body"),
			wantIssuer: 1, wantAPI: 1},
		{name: "repeated unauthorized stops", alwaysReject: true, wantIssuer: 2, wantAPI: 2},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("OPENAI_BASE_URL", "https://mtls.api.openai.com/v1/")
			config, issuer, api := newX509WorkloadIdentityIntegration(t)
			var exchanges, requests atomic.Int32
			var mu sync.Mutex
			var bearers, bodies []string
			issuer.server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				number := exchanges.Add(1)
				_, _ = io.WriteString(w, x509IntegrationTokenResponse(fmt.Sprintf("synthetic-bearer-%d", number)))
			})
			api.server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
				body, _ := io.ReadAll(request.Body)
				mu.Lock()
				bearers = append(bearers, request.Header.Get("Authorization"))
				bodies = append(bodies, string(body))
				mu.Unlock()
				w.Header().Set("Content-Type", "application/json")
				if requests.Add(1) == 1 || test.alwaysReject {
					w.WriteHeader(http.StatusUnauthorized)
					_, _ = io.WriteString(w, `{"error":{"message":"synthetic rejected bearer"}}`)
					return
				}
				_, _ = io.WriteString(w, `{"data":[]}`)
			})
			client := openai.NewClient(option.WithX509WorkloadIdentity(config))
			var err error
			if test.body == nil {
				_, err = client.Models.List(t.Context())
			} else {
				err = client.Execute(t.Context(), http.MethodPost, "models", test.body, nil)
			}
			if (err == nil) != test.wantSuccess {
				t.Errorf("401 replay result error=%v, want success=%v", err, test.wantSuccess)
			}
			if got := exchanges.Load(); got != test.wantIssuer {
				t.Errorf("401 replay issuer exchanges = %d, want %d", got, test.wantIssuer)
			}
			if got := requests.Load(); got != test.wantAPI {
				t.Errorf("401 replay API attempts = %d, want %d", got, test.wantAPI)
			}
			mu.Lock()
			defer mu.Unlock()
			if len(bearers) == 2 && (bearers[0] != "Bearer synthetic-bearer-1" ||
				bearers[1] != "Bearer synthetic-bearer-2") {
				t.Errorf("401 replay bearer generations = %v", bearers)
			}
			if len(bodies) == 2 && bodies[0] != bodies[1] {
				t.Errorf("401 replay changed its request body: %q versus %q", bodies[0], bodies[1])
			}
		})
	}
}

func TestX509WorkloadIdentityInvalidatesNonreplayableAndRepeatedUnauthorizedBearers(t *testing.T) {
	for _, test := range []struct {
		name          string
		body          any
		firstStatuses int32
		wantIssuer    int32
		wantAPI       int32
	}{
		{name: "non-replayable first body", body: strings.NewReader("synthetic-streaming-body"),
			firstStatuses: 1, wantIssuer: 2, wantAPI: 2},
		{name: "replayed bearer also rejected", firstStatuses: 2, wantIssuer: 3, wantAPI: 3},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("OPENAI_BASE_URL", "https://mtls.api.openai.com/v1/")
			config, issuer, api := newX509WorkloadIdentityIntegration(t)
			var exchanges, requests atomic.Int32
			issuer.server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				_, _ = io.WriteString(w, x509IntegrationTokenResponse(
					fmt.Sprintf("synthetic-invalidated-%d", exchanges.Add(1)),
				))
			})
			api.server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				if requests.Add(1) <= test.firstStatuses {
					w.WriteHeader(http.StatusUnauthorized)
					_, _ = io.WriteString(w, `{"error":{"message":"synthetic rejected bearer"}}`)
					return
				}
				_, _ = io.WriteString(w, `{"data":[]}`)
			})
			client := openai.NewClient(option.WithX509WorkloadIdentity(config))
			var err error
			if test.body == nil {
				_, err = client.Models.List(t.Context())
			} else {
				err = client.Execute(t.Context(), http.MethodPost, "models", test.body, nil)
			}
			if err == nil {
				t.Fatal("initial unauthorized request unexpectedly succeeded")
			}
			if _, err := client.Models.List(t.Context()); err != nil {
				t.Fatalf("fresh request reused an invalidated bearer: %v", err)
			}
			if exchanges.Load() != test.wantIssuer || requests.Load() != test.wantAPI {
				t.Errorf("rejected bearer issuer/API requests = %d/%d, want %d/%d",
					exchanges.Load(), requests.Load(), test.wantIssuer, test.wantAPI)
			}
		})
	}
}

func TestX509WorkloadIdentityDoesNotReplayMiddlewareTransformedBodies(t *testing.T) {
	t.Setenv("OPENAI_BASE_URL", "https://mtls.api.openai.com/v1/")
	config, issuer, api := newX509WorkloadIdentityIntegration(t)
	var exchanges, requests atomic.Int32
	var observed string
	var mu sync.Mutex
	issuer.server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, x509IntegrationTokenResponse(fmt.Sprintf("synthetic-transform-%d", exchanges.Add(1))))
	})
	api.server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		body, _ := io.ReadAll(request.Body)
		mu.Lock()
		observed = string(body)
		mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		if requests.Add(1) == 1 {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = io.WriteString(w, `{"error":{"message":"synthetic transformed bearer rejection"}}`)
			return
		}
		_, _ = io.WriteString(w, `{"data":[]}`)
	})
	client := openai.NewClient(
		option.WithX509WorkloadIdentity(config),
		option.WithMiddleware(func(request *http.Request, next option.MiddlewareNext) (*http.Response, error) {
			if request.Body != nil {
				request.Body = io.NopCloser(strings.NewReader("middleware-transformed-body"))
				request.ContentLength = int64(len("middleware-transformed-body"))
			}
			return next(request)
		}),
	)
	if err := client.Execute(t.Context(), http.MethodPost, "models", map[string]string{"input": "original"}, nil); err == nil {
		t.Fatal("body transformed by caller middleware was replayed after unauthorized")
	}
	mu.Lock()
	if observed != "middleware-transformed-body" {
		t.Errorf("first wire body = %q, want caller-transformed bytes", observed)
	}
	mu.Unlock()
	if exchanges.Load() != 1 || requests.Load() != 1 {
		t.Errorf("transformed request issuer/API calls = %d/%d, want 1/1", exchanges.Load(), requests.Load())
	}
	if _, err := client.Models.List(t.Context()); err != nil {
		t.Fatalf("bodyless request did not recover after transformed-body invalidation: %v", err)
	}
	if exchanges.Load() != 2 || requests.Load() != 2 {
		t.Errorf("post-invalidation issuer/API calls = %d/%d, want 2/2", exchanges.Load(), requests.Load())
	}
}

func TestX509WorkloadIdentityDoesNotRetryRefreshFailureAfterUnauthorized(t *testing.T) {
	t.Setenv("OPENAI_BASE_URL", "https://mtls.api.openai.com/v1/")
	config, issuer, api := newX509WorkloadIdentityIntegration(t)
	var exchanges, requests atomic.Int32
	issuer.server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if exchanges.Add(1) == 1 {
			_, _ = io.WriteString(w, x509IntegrationTokenResponse("synthetic-initial-bearer"))
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, `{"error":"invalid_grant","error_description":"synthetic-private-detail"}`)
	})
	api.server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		requests.Add(1)
		w.WriteHeader(http.StatusUnauthorized)
	})
	client := openai.NewClient(option.WithX509WorkloadIdentity(config))
	_, err := client.Models.List(t.Context())
	var oauthError *auth.OAuthError
	if !errors.As(err, &oauthError) || oauthError.StatusCode != http.StatusBadRequest {
		t.Fatalf("401 refresh failure lost its typed OAuth cause: %v", err)
	}
	if exchanges.Load() != 2 || requests.Load() != 1 {
		t.Errorf("401 refresh issuer/API attempts = %d/%d, want 2/1", exchanges.Load(), requests.Load())
	}
}

func TestX509WorkloadIdentityPreservesOrdinaryAPIStatusRetries(t *testing.T) {
	t.Setenv("OPENAI_BASE_URL", "https://mtls.api.openai.com/v1/")
	config, issuer, api := newX509WorkloadIdentityIntegration(t)
	var requests atomic.Int32
	api.server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if requests.Add(1) == 1 {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = io.WriteString(w, `{"error":{"message":"synthetic transient API failure"}}`)
			return
		}
		_, _ = io.WriteString(w, `{"data":[]}`)
	})
	client := openai.NewClient(option.WithX509WorkloadIdentity(config), option.WithMaxRetryDelay(time.Millisecond))
	if _, err := client.Models.List(t.Context()); err != nil {
		t.Fatalf("ordinary API 500 did not preserve generic request retries: %v", err)
	}
	if len(issuer.requests()) != 1 || requests.Load() != 2 {
		t.Errorf("cached issuer/API retry attempts = %d/%d, want 1/2", len(issuer.requests()), requests.Load())
	}
}

func TestX509WorkloadIdentitySharesRefreshAfterConcurrentUnauthorized(t *testing.T) {
	t.Setenv("OPENAI_BASE_URL", "https://mtls.api.openai.com/v1/")
	config, issuer, api := newX509WorkloadIdentityIntegration(t)
	var exchanges, oldRequests, newRequests atomic.Int32
	bothRejected := make(chan struct{})
	issuer.server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		number := exchanges.Add(1)
		_, _ = io.WriteString(w, x509IntegrationTokenResponse(fmt.Sprintf("synthetic-concurrent-%d", number)))
	})
	api.server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		if request.Header.Get("Authorization") == "Bearer synthetic-concurrent-1" {
			if oldRequests.Add(1) == 2 {
				close(bothRejected)
			}
			<-bothRejected
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		newRequests.Add(1)
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"data":[]}`)
	})
	client := openai.NewClient(option.WithX509WorkloadIdentity(config))
	results := make(chan error, 2)
	for range 2 {
		go func() {
			_, err := client.Models.List(t.Context())
			results <- err
		}()
	}
	for range 2 {
		select {
		case err := <-results:
			if err != nil {
				t.Errorf("concurrent 401 replay: %v", err)
			}
		case <-time.After(5 * time.Second):
			t.Fatal("concurrent unauthorized requests did not complete")
		}
	}
	if exchanges.Load() != 2 || oldRequests.Load() != 2 || newRequests.Load() != 2 {
		t.Errorf("concurrent stale 401 issuer/old/new requests = %d/%d/%d, want 2/2/2",
			exchanges.Load(), oldRequests.Load(), newRequests.Load())
	}
}

func x509IntegrationTokenResponse(token string) string {
	return fmt.Sprintf(`{"access_token":%q,"token_type":"Bearer",`+
		`"issued_token_type":"urn:ietf:params:oauth:token-type:access_token","expires_in":60}`, token)
}

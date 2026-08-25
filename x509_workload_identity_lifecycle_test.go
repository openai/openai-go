package openai_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
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

func TestX509WorkloadIdentityCountsShortCircuitedCallerMiddlewareAttempts(t *testing.T) {
	t.Setenv("OPENAI_BASE_URL", "https://mtls.api.openai.com/v1/")
	config, issuer, api := newX509WorkloadIdentityIntegration(t)
	var middlewareAttempts, exchanges, dispatched atomic.Int32
	issuer.server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if exchanges.Add(1) == 1 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		_, _ = io.WriteString(w, x509IntegrationTokenResponse("synthetic-shared-attempt-token"))
	})
	api.server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		dispatched.Add(1)
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"data":[]}`)
	})
	client := openai.NewClient(
		option.WithX509WorkloadIdentity(config),
		option.WithMaxRetries(2),
		option.WithMaxRetryDelay(time.Millisecond),
		option.WithMiddleware(func(request *http.Request, next option.MiddlewareNext) (*http.Response, error) {
			if middlewareAttempts.Add(1) == 1 {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Header:     make(http.Header),
					Body:       io.NopCloser(strings.NewReader(`{"error":{"message":"synthetic caller failure"}}`)),
				}, nil
			}
			return next(request)
		}),
	)
	if _, err := client.Models.List(t.Context()); err != nil {
		t.Fatalf("shared caller-middleware and issuer retry budget did not recover: %v", err)
	}
	if middlewareAttempts.Load() != 2 || exchanges.Load() != 2 || dispatched.Load() != 1 {
		t.Errorf("caller/issuer/API attempts = %d/%d/%d, want 2/2/1 within one shared budget",
			middlewareAttempts.Load(), exchanges.Load(), dispatched.Load())
	}
}

func TestX509WorkloadIdentityRejectsCallerMiddlewareRemovingRetryScope(t *testing.T) {
	t.Setenv("OPENAI_BASE_URL", "https://mtls.api.openai.com/v1/")
	config, issuer, api := newX509WorkloadIdentityIntegration(t)
	client := openai.NewClient(
		option.WithX509WorkloadIdentity(config),
		option.WithMiddleware(func(request *http.Request, next option.MiddlewareNext) (*http.Response, error) {
			return next(request.WithContext(context.Background()))
		}),
	)
	if _, err := client.Models.List(t.Context()); err == nil {
		t.Fatal("caller middleware removed the request-owned retry scope")
	}
	assertX509WorkloadNoRequests(t, issuer, api)
}

func TestX509WorkloadIdentityUnauthorizedRecoveryReentersSDKAttempt(t *testing.T) {
	for _, test := range []struct {
		name        string
		maximum     int
		wantCalls   int32
		wantSuccess bool
	}{
		{name: "fresh middleware and attempt timeout", maximum: 1, wantCalls: 2, wantSuccess: true},
		{name: "zero retries prevents unauthorized recovery", maximum: 0, wantCalls: 1},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("OPENAI_BASE_URL", "https://mtls.api.openai.com/v1/")
			config, issuer, api := newX509WorkloadIdentityIntegration(t)
			var middlewareAttempts, exchanges, dispatched atomic.Int32
			var mu sync.Mutex
			var deadlines []time.Time
			issuer.server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				_, _ = io.WriteString(w, x509IntegrationTokenResponse(
					fmt.Sprintf("synthetic-outer-recovery-%d", exchanges.Add(1)),
				))
			})
			api.server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				if dispatched.Add(1) == 1 {
					time.Sleep(10 * time.Millisecond)
					w.WriteHeader(http.StatusUnauthorized)
					_, _ = io.WriteString(w, `{"error":{"message":"synthetic rejected bearer"}}`)
					return
				}
				_, _ = io.WriteString(w, `{"data":[]}`)
			})
			client := openai.NewClient(
				option.WithX509WorkloadIdentity(config),
				option.WithMaxRetries(test.maximum),
				option.WithMaxRetryDelay(time.Millisecond),
				option.WithRequestTimeout(time.Second),
				option.WithMiddleware(func(request *http.Request, next option.MiddlewareNext) (*http.Response, error) {
					middlewareAttempts.Add(1)
					if deadline, ok := request.Context().Deadline(); ok {
						mu.Lock()
						deadlines = append(deadlines, deadline)
						mu.Unlock()
					}
					return next(request)
				}),
			)
			_, err := client.Models.List(t.Context())
			if (err == nil) != test.wantSuccess {
				t.Errorf("outer unauthorized recovery error = %v, want success %v", err, test.wantSuccess)
			}
			if middlewareAttempts.Load() != test.wantCalls || exchanges.Load() != test.wantCalls ||
				dispatched.Load() != test.wantCalls {
				t.Errorf("caller/issuer/API attempts = %d/%d/%d, want %d/%d/%d",
					middlewareAttempts.Load(), exchanges.Load(), dispatched.Load(),
					test.wantCalls, test.wantCalls, test.wantCalls)
			}
			mu.Lock()
			defer mu.Unlock()
			if len(deadlines) != int(test.wantCalls) {
				t.Fatalf("caller middleware observed %d per-attempt deadlines, want %d", len(deadlines), test.wantCalls)
			}
			if len(deadlines) == 2 && !deadlines[1].After(deadlines[0].Add(5*time.Millisecond)) {
				t.Errorf("unauthorized recovery reused its first attempt deadline: %s then %s", deadlines[0], deadlines[1])
			}
		})
	}
}

func TestX509WorkloadIdentityRetriesOnlyTransientOAuthErrorCodes(t *testing.T) {
	for _, test := range []struct {
		name        string
		code        string
		status      int
		wantCalls   int32
		wantSuccess bool
	}{
		{name: "temporarily unavailable", code: "temporarily_unavailable", status: 400, wantCalls: 2, wantSuccess: true},
		{name: "server error", code: "server_error", status: 401, wantCalls: 2, wantSuccess: true},
		{name: "invalid grant remains permanent", code: "invalid_grant", status: 400, wantCalls: 1},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("OPENAI_BASE_URL", "https://mtls.api.openai.com/v1/")
			config, issuer, api := newX509WorkloadIdentityIntegration(t)
			var exchanges atomic.Int32
			issuer.server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				if exchanges.Add(1) == 1 {
					w.WriteHeader(test.status)
					_, _ = fmt.Fprintf(w, `{"error":%q,"error_description":"synthetic-private-issuer-detail"}`, test.code)
					return
				}
				_, _ = io.WriteString(w, x509IntegrationTokenResponse("synthetic-recovered-oauth-token"))
			})
			client := openai.NewClient(option.WithX509WorkloadIdentity(config),
				option.WithMaxRetries(1), option.WithMaxRetryDelay(time.Millisecond))
			_, err := client.Models.List(t.Context())
			if (err == nil) != test.wantSuccess {
				t.Errorf("OAuth error %q produced error=%v, want success %v", test.code, err, test.wantSuccess)
			}
			if err != nil && strings.Contains(err.Error(), "synthetic-private-issuer-detail") {
				t.Errorf("OAuth error disclosed sensitive issuer details: %v", err)
			}
			if got := exchanges.Load(); got != test.wantCalls {
				t.Errorf("OAuth error %q caused %d issuer attempts, want %d", test.code, got, test.wantCalls)
			}
			wantAPI := 0
			if test.wantSuccess {
				wantAPI = 1
			}
			if got := len(api.requests()); got != wantAPI {
				t.Errorf("OAuth error %q dispatched %d API requests, want %d", test.code, got, wantAPI)
			}
		})
	}
}

func TestX509WorkloadIdentityHonorsBoundedIssuerRetryAfter(t *testing.T) {
	for _, test := range []struct {
		name    string
		status  int
		header  string
		value   string
		maximum time.Duration
		minimum time.Duration
		ceiling time.Duration
	}{
		{name: "request timeout", status: http.StatusRequestTimeout,
			header: "Retry-After-Ms", value: "40", maximum: 200 * time.Millisecond, minimum: 30 * time.Millisecond},
		{name: "conflict", status: http.StatusConflict,
			header: "Retry-After-Ms", value: "40", maximum: 200 * time.Millisecond, minimum: 30 * time.Millisecond},
		{name: "rate limited", status: http.StatusTooManyRequests,
			header: "Retry-After", value: "0.04", maximum: 200 * time.Millisecond, minimum: 30 * time.Millisecond},
		{name: "service unavailable", status: http.StatusServiceUnavailable,
			header: "Retry-After-Ms", value: "40", maximum: 200 * time.Millisecond, minimum: 30 * time.Millisecond},
		{name: "caller delay clamps issuer hint", status: http.StatusTooManyRequests,
			header: "Retry-After", value: "8", maximum: 5 * time.Millisecond, ceiling: 250 * time.Millisecond},
		{name: "explicit zero is immediate", status: http.StatusServiceUnavailable,
			header: "Retry-After-Ms", value: "0", maximum: time.Second, ceiling: 250 * time.Millisecond},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("OPENAI_BASE_URL", "https://mtls.api.openai.com/v1/")
			config, issuer, api := newX509WorkloadIdentityIntegration(t)
			var exchanges atomic.Int32
			attempted := make(chan time.Time, 2)
			issuer.server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				attempted <- time.Now()
				if exchanges.Add(1) == 1 {
					w.Header().Set(test.header, test.value)
					w.WriteHeader(test.status)
					return
				}
				_, _ = io.WriteString(w, x509IntegrationTokenResponse("synthetic-retry-after-bearer"))
			})
			client := openai.NewClient(option.WithX509WorkloadIdentity(config), option.WithMaxRetryDelay(test.maximum))
			if _, err := client.Models.List(t.Context()); err != nil {
				t.Fatalf("issuer-directed retry failed: %v", err)
			}
			first, second := <-attempted, <-attempted
			delay := second.Sub(first)
			if delay < test.minimum {
				t.Errorf("issuer-directed retry waited %s, want at least %s", delay, test.minimum)
			}
			if test.ceiling > 0 && delay > test.ceiling {
				t.Errorf("bounded issuer-directed retry waited %s, want no more than %s", delay, test.ceiling)
			}
			if exchanges.Load() != 2 || len(api.requests()) != 1 {
				t.Errorf("issuer-directed retry issuer/API attempts = %d/%d", exchanges.Load(), len(api.requests()))
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

func TestX509WorkloadIdentityPaginationUsesIndependentLogicalRetryScopes(t *testing.T) {
	for _, test := range []struct {
		name      string
		maximum   int
		failPage  int
		wantCalls int32
	}{
		{name: "zero-retry budget still fetches four pages", maximum: 0, wantCalls: 4},
		{name: "default budget fetches four pages", maximum: 2, wantCalls: 4},
		{name: "later page has a fresh retry budget", maximum: 1, failPage: 3, wantCalls: 5},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("OPENAI_BASE_URL", "https://mtls.api.openai.com/v1/")
			config, issuer, api := newX509WorkloadIdentityIntegration(t)
			var requests atomic.Int32
			var failed atomic.Bool
			api.server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
				requests.Add(1)
				previous := request.URL.Query().Get("after")
				page := 1
				if previous != "" {
					parsed, err := strconv.Atoi(previous)
					if err != nil {
						w.WriteHeader(http.StatusBadRequest)
						return
					}
					page = parsed + 1
				}
				w.Header().Set("Content-Type", "application/json")
				if page == test.failPage && failed.CompareAndSwap(false, true) {
					w.WriteHeader(http.StatusInternalServerError)
					_, _ = io.WriteString(w, `{"error":{"message":"synthetic transient page failure"}}`)
					return
				}
				_, _ = fmt.Fprintf(w, `{"data":[{"id":%q}],"has_more":%t}`, strconv.Itoa(page), page < 4)
			})
			client := openai.NewClient(option.WithX509WorkloadIdentity(config),
				option.WithMaxRetries(test.maximum), option.WithMaxRetryDelay(time.Millisecond))
			pages := client.Batches.ListAutoPaging(t.Context(), openai.BatchListParams{})
			var received int
			for pages.Next() {
				received++
				if id := pages.Current().ID; id != strconv.Itoa(received) {
					t.Errorf("auto-paging item %d has ID %q", received, id)
				}
			}
			if err := pages.Err(); err != nil || received != 4 {
				t.Fatalf("independently scoped auto-paging returned %d pages, error = %v", received, err)
			}
			if requests.Load() != test.wantCalls || len(issuer.requests()) != 1 {
				t.Errorf("auto-paging issuer/API calls = %d/%d, want 1/%d",
					len(issuer.requests()), requests.Load(), test.wantCalls)
			}
		})
	}
}

func TestX509WorkloadIdentityConcurrentPageClonesHaveIndependentBudgets(t *testing.T) {
	t.Setenv("OPENAI_BASE_URL", "https://mtls.api.openai.com/v1/")
	config, issuer, api := newX509WorkloadIdentityIntegration(t)
	var subsequent atomic.Int32
	api.server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if request.URL.Query().Get("after") == "" {
			_, _ = io.WriteString(w, `{"data":[{"id":"1"}],"has_more":true}`)
			return
		}
		if subsequent.Add(1) <= 2 {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = io.WriteString(w, `{"error":{"message":"synthetic concurrent page failure"}}`)
			return
		}
		_, _ = io.WriteString(w, `{"data":[{"id":"2"}],"has_more":false}`)
	})
	client := openai.NewClient(option.WithX509WorkloadIdentity(config), option.WithMaxRetryDelay(time.Millisecond))
	page, err := client.Batches.List(t.Context(), openai.BatchListParams{})
	if err != nil {
		t.Fatalf("load first concurrently cloned page: %v", err)
	}
	results := make(chan error, 2)
	for range 2 {
		go func() {
			next, nextErr := page.GetNextPage()
			if nextErr == nil && (next == nil || len(next.Data) != 1 || next.Data[0].ID != "2") {
				nextErr = errors.New("concurrent page clone returned an unexpected item")
			}
			results <- nextErr
		}()
	}
	for range 2 {
		if nextErr := <-results; nextErr != nil {
			t.Errorf("concurrent independently scoped page clone: %v", nextErr)
		}
	}
	if len(issuer.requests()) != 1 {
		t.Errorf("concurrent page clones minted %d workload tokens", len(issuer.requests()))
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

func TestX509WorkloadIdentityDoesNotRestoreMiddlewareRemovedBody(t *testing.T) {
	t.Setenv("OPENAI_BASE_URL", "https://mtls.api.openai.com/v1/")
	config, issuer, api := newX509WorkloadIdentityIntegration(t)
	var requests atomic.Int32
	api.server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		body, err := io.ReadAll(request.Body)
		if err != nil || len(body) != 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if requests.Add(1) == 1 {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = io.WriteString(w, `{"error":{"message":"synthetic removed-body rejection"}}`)
			return
		}
		_, _ = io.WriteString(w, `{"data":[]}`)
	})
	client := openai.NewClient(
		option.WithX509WorkloadIdentity(config),
		option.WithMiddleware(func(request *http.Request, next option.MiddlewareNext) (*http.Response, error) {
			if request.Body != nil {
				request.Body = nil
				request.ContentLength = 0
			}
			return next(request)
		}),
	)
	if err := client.Execute(t.Context(), http.MethodPost, "models",
		map[string]string{"secret": "synthetic-removed-payload"}, nil); err != nil {
		t.Fatalf("bodyless unauthorized replay restored its removed payload: %v", err)
	}
	if requests.Load() != 2 || len(issuer.requests()) != 2 {
		t.Errorf("removed-body issuer/API attempts = %d/%d, want 2/2", len(issuer.requests()), requests.Load())
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

package auth_test

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	openai "github.com/openai/openai-go/v4"
	"github.com/openai/openai-go/v4/auth"
	"github.com/openai/openai-go/v4/internal/requestconfig"
	"github.com/openai/openai-go/v4/option"
)

func TestWorkloadIdentityPublicClientReplaySurvivesRetainedCacheEviction(t *testing.T) {
	identity := newOrdinaryWorkloadIdentity(t)
	started, failRefresh := make(chan struct{}), make(chan struct{})
	var failOnce sync.Once
	release := func() { failOnce.Do(func() { close(failRefresh) }) }
	t.Cleanup(release)
	var exchanges, apiCalls atomic.Int32
	transport := &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
		if req.URL.Host == "auth.openai.com" {
			call := exchanges.Add(1)
			if call == 2 {
				close(started)
				select {
				case <-failRefresh:
					return nil, errors.New("synthetic background refresh failure")
				case <-req.Context().Done():
					return nil, req.Context().Err()
				}
			}
			token := "synthetic-rejected"
			if call >= 4 {
				token = "synthetic-replacement"
			}
			return ordinaryWorkloadResponse(http.StatusOK, fmt.Sprintf(`{"access_token":%q,"expires_in":600}`, token)), nil
		}
		call := apiCalls.Add(1)
		if call == 1 {
			select {
			case <-started:
			case <-time.After(5 * time.Second):
				return nil, errors.New("background refresh never started")
			}
		}
		if req.Header.Get("Authorization") != "Bearer synthetic-replacement" {
			return ordinaryWorkloadResponse(http.StatusUnauthorized, `{"error":{"message":"synthetic unauthorized"}}`), nil
		}
		response := ordinaryWorkloadResponse(http.StatusOK, `{"data":[]}`)
		response.Header.Set("Content-Type", "application/json")
		return response, nil
	}}
	httpClient := &http.Client{Transport: transport}
	if _, err := identity.GetToken(t.Context(), httpClient); err != nil {
		t.Fatal(err)
	}
	evict := auth.WorkloadIdentityEvictionWindowForTest(t, identity, httpClient)
	// Install the public authentication middleware and the same logical-request
	// scope used by WithWorkloadIdentity, retaining direct access to the identity
	// solely to coordinate an exact eviction without sleeps or scheduler races.
	client := openai.NewClient(
		option.WithAPIKey(""), option.WithHTTPClient(httpClient), option.WithMaxRetries(2),
		option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
			response, err := next(req)
			if err == nil && response != nil && response.StatusCode == http.StatusUnauthorized && apiCalls.Load() == 1 {
				evict(release)
			}
			return response, err
		}),
		option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
			return auth.WorkloadIdentityMiddleware(identity, httpClient, req, next)
		}),
		requestconfig.WithRequestFinalizer(func(cfg *requestconfig.RequestConfig) error {
			cfg.InstallRequestRetryScope(true)
			cfg.InstallRequestAttemptMiddleware()
			return nil
		}),
	)
	if _, err := client.Models.List(t.Context()); err != nil {
		t.Fatal(err)
	}
	if exchanges.Load() != 4 || apiCalls.Load() != 2 {
		t.Fatalf("exchanges=%d API attempts=%d, want 4 and 2", exchanges.Load(), apiCalls.Load())
	}
}

func TestWorkloadIdentityNonComparableDoerRetainsReplayState(t *testing.T) {
	started, releaseRefresh := make(chan struct{}), make(chan struct{})
	var releaseOnce sync.Once
	release := func() { releaseOnce.Do(func() { close(releaseRefresh) }) }
	t.Cleanup(release)
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	var exchanges, apiCalls atomic.Int32
	doer := ordinaryOpaqueHTTPDoer(func(req *http.Request) (*http.Response, error) {
		if req.URL.Host == "auth.openai.com" {
			call := exchanges.Add(1)
			if call == 2 {
				close(started)
				select {
				case <-releaseRefresh:
				case <-req.Context().Done():
					return nil, req.Context().Err()
				}
			}
			token := "synthetic-rejected"
			if call >= 3 {
				token = "synthetic-replacement"
			}
			return ordinaryWorkloadResponse(http.StatusOK, fmt.Sprintf(`{"access_token":%q,"expires_in":600}`, token)), nil
		}
		call := apiCalls.Add(1)
		if call == 1 {
			return ordinaryWorkloadResponse(http.StatusInternalServerError, `{"error":{"message":"synthetic retry"}}`), nil
		}
		if call == 2 {
			select {
			case <-started:
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}
		if req.Header.Get("Authorization") != "Bearer synthetic-replacement" {
			return ordinaryWorkloadResponse(http.StatusUnauthorized, `{"error":{"message":"synthetic unauthorized"}}`), nil
		}
		response := ordinaryWorkloadResponse(http.StatusOK, `{"data":[]}`)
		response.Header.Set("Content-Type", "application/json")
		return response, nil
	})
	client := openai.NewClient(
		option.WithAPIKey(""), option.WithHTTPClient(doer), option.WithMaxRetries(3),
		option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
			response, err := next(req)
			if err == nil && response != nil && response.StatusCode == http.StatusUnauthorized {
				release()
			}
			return response, err
		}),
		option.WithWorkloadIdentity(auth.WorkloadIdentity{
			IdentityProviderID: "synthetic-provider", ServiceAccountID: "synthetic-service-account",
			Provider:             &mockProvider{token: "synthetic-subject", tokenType: auth.SubjectTokenTypeJWT},
			RefreshBufferSeconds: 3600,
		}),
	)
	if _, err := client.Models.List(ctx); err != nil {
		t.Fatalf("Models.List with non-comparable doer: got %v, want successful replay", err)
	}
	if exchanges.Load() != 3 || apiCalls.Load() != 3 {
		t.Errorf("replay exchanges=%d API calls=%d, want 3 and 3", exchanges.Load(), apiCalls.Load())
	}
	if _, err := client.Models.List(ctx); err != nil {
		t.Fatalf("second Models.List: got %v, want independent success", err)
	}
	if exchanges.Load() != 4 || apiCalls.Load() != 4 {
		t.Errorf("independent request exchanges=%d API calls=%d, want 4 and 4", exchanges.Load(), apiCalls.Load())
	}
}

func TestWorkloadIdentityNonComparableDoersStaySeparateWithinRequest(t *testing.T) {
	identity := newOrdinaryWorkloadIdentity(t)
	counts := map[string]int{}
	makeDoer := func(token string) auth.HTTPDoer {
		return ordinaryOpaqueHTTPDoer(func(*http.Request) (*http.Response, error) {
			counts[token]++
			return ordinaryWorkloadResponse(http.StatusOK, fmt.Sprintf(`{"access_token":%q,"expires_in":3600}`, token)), nil
		})
	}
	a, b := makeDoer("synthetic-a"), makeDoer("synthetic-b")
	scope := requestconfig.NewRequestRetryScope(0, 0, false, nil)
	ctx := requestconfig.WithRequestRetryScope(t.Context(), scope)
	for _, test := range []struct {
		doer  auth.HTTPDoer
		token string
	}{{a, "synthetic-a"}, {b, "synthetic-b"}, {a, "synthetic-a"}, {b, "synthetic-b"}} {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.openai.com/v1/models", nil)
		if err != nil {
			t.Fatal(err)
		}
		response, err := auth.WorkloadIdentityMiddleware(identity, test.doer, req, func(actual *http.Request) (*http.Response, error) {
			if got := actual.Header.Get("Authorization"); got != "Bearer "+test.token {
				t.Errorf("doer %s authorization=%q, want own bearer", test.token, got)
			}
			return ordinaryWorkloadResponse(http.StatusOK, ""), nil
		})
		if err != nil {
			t.Fatal(err)
		}
		_ = response.Body.Close()
	}
	if counts["synthetic-a"] != 1 || counts["synthetic-b"] != 1 {
		t.Errorf("per-request doer exchanges=%v, want one per identity", counts)
	}
}

// A slice-backed custom client is supported by HTTPDoer but cannot be a map key.
type retrySliceHTTPDoer []ordinaryOpaqueHTTPDoer

func (d retrySliceHTTPDoer) Do(req *http.Request) (*http.Response, error) {
	return d[0](req)
}

func TestWorkloadIdentityNonComparableDoerRetainsBearerAcrossOrdinaryRetries(t *testing.T) {
	for _, status := range []int{http.StatusTooManyRequests, http.StatusInternalServerError} {
		for _, kind := range []string{"function", "slice"} {
			t.Run(fmt.Sprintf("%s/%d", kind, status), func(t *testing.T) {
				var exchanges, apiCalls atomic.Int32
				fn := ordinaryOpaqueHTTPDoer(func(req *http.Request) (*http.Response, error) {
					if req.URL.Host == "auth.openai.com" {
						exchanges.Add(1)
						return ordinaryWorkloadResponse(http.StatusOK, `{"access_token":"synthetic-valid","expires_in":3600}`), nil
					}
					if req.Header.Get("Authorization") != "Bearer synthetic-valid" {
						t.Errorf("request did not retain its synthetic bearer")
					}
					if apiCalls.Add(1)%3 != 0 {
						response := ordinaryWorkloadResponse(status, `{"error":{"message":"synthetic retry"}}`)
						response.Header.Set("retry-after-ms", "1")
						return response, nil
					}
					response := ordinaryWorkloadResponse(http.StatusOK, `{"data":[]}`)
					response.Header.Set("Content-Type", "application/json")
					return response, nil
				})
				var doer auth.HTTPDoer = fn
				if kind == "slice" {
					doer = retrySliceHTTPDoer{fn}
				}
				client := openai.NewClient(
					option.WithAPIKey(""), option.WithHTTPClient(doer), option.WithMaxRetries(2),
					option.WithWorkloadIdentity(auth.WorkloadIdentity{
						IdentityProviderID: "synthetic-provider", ServiceAccountID: "synthetic-service-account",
						Provider: &mockProvider{token: "synthetic-subject", tokenType: auth.SubjectTokenTypeJWT},
					}),
				)
				for call := int32(1); call <= 2; call++ {
					if _, err := client.Models.List(t.Context()); err != nil {
						t.Fatal(err)
					}
					if exchanges.Load() != call || apiCalls.Load() != call*3 {
						t.Errorf("exchanges=%d API attempts=%d, want %d and %d", exchanges.Load(), apiCalls.Load(), call, call*3)
					}
				}
			})
		}
	}
}

// Comparable values can still be unsuitable equality keys: NaN != NaN.
type nonReflexiveHTTPIdentity struct {
	marker    float64
	transport *closureTransport
}

func (d nonReflexiveHTTPIdentity) Do(req *http.Request) (*http.Response, error) {
	return d.transport.RoundTrip(req)
}
func (d nonReflexiveHTTPIdentity) RoundTrip(req *http.Request) (*http.Response, error) {
	return d.transport.RoundTrip(req)
}

func TestWorkloadIdentityNonReflexiveIdentityRetainsCacheAndReplay(t *testing.T) {
	for _, native := range []bool{false, true} {
		t.Run(fmt.Sprintf("native=%t", native), func(t *testing.T) {
			var exchanges, calls atomic.Int32
			transport := &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
				if req.URL.Host == "auth.openai.com" {
					n := exchanges.Add(1)
					return ordinaryWorkloadResponse(http.StatusOK, fmt.Sprintf(`{"access_token":"synthetic-%d","expires_in":3600}`, n)), nil
				}
				n := calls.Add(1)
				if n == 2 {
					return ordinaryWorkloadResponse(http.StatusUnauthorized, `{"error":{"message":"synthetic unauthorized"}}`), nil
				}
				want := "Bearer synthetic-1"
				if n >= 3 {
					want = "Bearer synthetic-2"
				}
				if req.Header.Get("Authorization") != want {
					t.Errorf("attempt %d did not reuse the expected bearer", n)
				}
				r := ordinaryWorkloadResponse(http.StatusOK, `{"data":[]}`)
				r.Header.Set("Content-Type", "application/json")
				return r, nil
			}}
			boxed := nonReflexiveHTTPIdentity{marker: math.NaN(), transport: transport}
			var doer auth.HTTPDoer = boxed
			if native {
				doer = &http.Client{Transport: boxed}
			}
			client := openai.NewClient(option.WithAPIKey(""), option.WithHTTPClient(doer), option.WithMaxRetries(2),
				option.WithWorkloadIdentity(auth.WorkloadIdentity{
					IdentityProviderID: "synthetic-provider", ServiceAccountID: "synthetic-service-account",
					Provider: &mockProvider{token: "synthetic-subject", tokenType: auth.SubjectTokenTypeJWT},
				}))
			for i := 0; i < 3; i++ {
				if _, err := client.Models.List(t.Context()); err != nil {
					t.Fatal(err)
				}
			}
			if exchanges.Load() != 2 || calls.Load() != 4 {
				t.Errorf("exchanges=%d attempts=%d, want 2 and 4", exchanges.Load(), calls.Load())
			}
		})
	}
}

func TestWorkloadIdentityNonReflexiveIdentitiesStaySeparate(t *testing.T) {
	for _, native := range []bool{false, true} {
		t.Run(fmt.Sprintf("native=%t", native), func(t *testing.T) {
			identity := newOrdinaryWorkloadIdentity(t)
			var exchanges atomic.Int32
			makeIdentity := func(token string) nonReflexiveHTTPIdentity {
				return nonReflexiveHTTPIdentity{marker: math.NaN(), transport: &closureTransport{fn: func(*http.Request) (*http.Response, error) {
					exchanges.Add(1)
					return ordinaryWorkloadResponse(http.StatusOK, fmt.Sprintf(`{"access_token":%q,"expires_in":3600}`, token)), nil
				}}}
			}
			a, b := makeIdentity("synthetic-a"), makeIdentity("synthetic-b")
			var first, second auth.HTTPDoer = a, b
			client := &http.Client{}
			for _, item := range []struct {
				doer      auth.HTTPDoer
				transport nonReflexiveHTTPIdentity
				token     string
			}{
				{first, a, "synthetic-a"}, {second, b, "synthetic-b"}, {first, a, "synthetic-a"}, {second, b, "synthetic-b"},
			} {
				doer := item.doer
				if native {
					client.Transport = item.transport
					doer = client
				}
				got, err := identity.GetToken(t.Context(), doer)
				if err != nil || got != item.token {
					t.Fatalf("identity did not retain its own token: %v", err)
				}
			}
			// Re-boxing a value into a native transport may select a new identity.
			// Either behavior is safe, but tokens from the other transport must not leak.
			if !native && exchanges.Load() != 2 {
				t.Errorf("exchanges=%d, want 2", exchanges.Load())
			}
		})
	}
}

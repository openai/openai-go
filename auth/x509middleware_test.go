package auth

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/openai/openai-go/v3/internal/requestconfig"
)

func TestX509WorkloadIdentityMiddlewareRejectsMalformedRequestsBeforeExchange(t *testing.T) {
	var exchanges atomic.Int32
	fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		exchanges.Add(1)
		_, _ = io.WriteString(w, x509ValidExchangeResponse())
	}))
	identity := newX509LifecycleIdentity(t, fixture)
	for _, test := range []struct {
		name   string
		mutate func(*http.Request) *http.Request
	}{
		{name: "nil request", mutate: func(*http.Request) *http.Request { return nil }},
		{name: "nil headers", mutate: func(request *http.Request) *http.Request {
			request.Header = nil
			return request
		}},
		{name: "attacker origin", mutate: func(request *http.Request) *http.Request {
			request.URL.Host = "attacker.example.test"
			return request
		}},
		{name: "approved issuer is not API origin", mutate: func(request *http.Request) *http.Request {
			request.URL.Host = x509AuthenticationHost
			return request
		}},
		{name: "path traversal", mutate: func(request *http.Request) *http.Request {
			request.URL.Path = "/v1/../attacker"
			return request
		}},
		{name: "existing valid bearer", mutate: func(request *http.Request) *http.Request {
			request.Header.Set("Authorization", "Bearer synthetic-existing-token")
			return request
		}},
		{name: "credential alias", mutate: func(request *http.Request) *http.Request {
			request.Header.Set("X_API_KEY", "synthetic-existing-key")
			return request
		}},
		{name: "unsupported CONNECT method", mutate: func(request *http.Request) *http.Request {
			request.Method = http.MethodConnect
			return request
		}},
		{name: "custom HTTP framing", mutate: func(request *http.Request) *http.Request {
			request.Header.Set("Connection", "upgrade")
			return request
		}},
		{name: "custom transfer encoding", mutate: func(request *http.Request) *http.Request {
			request.TransferEncoding = []string{"chunked"}
			return request
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			request, err := http.NewRequestWithContext(t.Context(), http.MethodGet,
				"https://"+x509APIHost+"/v1/models", nil)
			if err != nil {
				t.Fatalf("construct malformed middleware request: %v", err)
			}
			response, middlewareErr := X509WorkloadIdentityMiddleware(identity, fixture.capability, test.mutate(request),
				func(*http.Request) (*http.Response, error) {
					t.Error("malformed request reached the authenticated middleware dispatch")
					return nil, nil
				})
			if response != nil && response.Body != nil {
				if closeErr := response.Body.Close(); closeErr != nil {
					t.Errorf("close unexpected malformed middleware response: %v", closeErr)
				}
			}
			if response != nil || middlewareErr == nil {
				t.Errorf("malformed direct middleware request = response:%v error:%v", response, middlewareErr)
			}
		})
	}
	if got := exchanges.Load(); got != 0 {
		t.Errorf("malformed direct middleware requests minted %d tokens", got)
	}
}

func TestX509WorkloadIdentityMiddlewareAcceptsExplicitDefaultHTTPSPort(t *testing.T) {
	fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, x509ValidExchangeResponse())
	}))
	identity := newX509LifecycleIdentity(t, fixture)
	request, err := http.NewRequestWithContext(t.Context(), http.MethodGet,
		"https://"+x509APIHost+":443/v1/models", nil)
	if err != nil {
		t.Fatalf("construct explicit-port middleware request: %v", err)
	}
	response, err := X509WorkloadIdentityMiddleware(identity, fixture.capability, request,
		func(authenticated *http.Request) (*http.Response, error) {
			if got := authenticated.Header.Get("Authorization"); got != "Bearer "+x509ExchangeSyntheticToken {
				t.Errorf("explicit-port middleware bearer = %q", got)
			}
			return &http.Response{StatusCode: http.StatusOK, Body: http.NoBody}, nil
		})
	if err != nil || response == nil || response.StatusCode != http.StatusOK {
		t.Fatalf("explicit-port middleware response=%v error=%v", response, err)
	}
	if closeErr := response.Body.Close(); closeErr != nil {
		t.Fatalf("close explicit-port middleware response: %v", closeErr)
	}
}

func TestX509WorkloadIdentityMiddlewareDoesNotReplayUnscopedTransformedBody(t *testing.T) {
	var exchanges, dispatched atomic.Int32
	fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		exchanges.Add(1)
		_, _ = io.WriteString(w, x509ValidExchangeResponse())
	}))
	identity := newX509LifecycleIdentity(t, fixture)
	request, err := http.NewRequestWithContext(t.Context(), http.MethodPost,
		"https://"+x509APIHost+"/v1/models", strings.NewReader("original-serialized-body"))
	if err != nil {
		t.Fatalf("construct replayable middleware request: %v", err)
	}
	if request.GetBody == nil {
		t.Fatal("negative-control middleware request lacks its original replay factory")
	}
	request.Body = io.NopCloser(strings.NewReader("caller-transformed-body"))
	response, err := X509WorkloadIdentityMiddleware(identity, fixture.capability, request,
		func(authenticated *http.Request) (*http.Response, error) {
			dispatched.Add(1)
			body, readErr := io.ReadAll(authenticated.Body)
			if readErr != nil || string(body) != "caller-transformed-body" {
				t.Errorf("unscoped middleware wire body = %q, error = %v", body, readErr)
			}
			return &http.Response{
				StatusCode: http.StatusUnauthorized,
				Body:       io.NopCloser(strings.NewReader("synthetic unauthorized response")),
			}, nil
		})
	if err != nil || response == nil || response.StatusCode != http.StatusUnauthorized {
		t.Fatalf("unscoped transformed request = response:%v error:%v", response, err)
	}
	if closeErr := response.Body.Close(); closeErr != nil {
		t.Fatalf("close preserved unauthorized response: %v", closeErr)
	}
	if exchanges.Load() != 1 || dispatched.Load() != 1 {
		t.Errorf("unscoped transformed request issuer/dispatch attempts = %d/%d", exchanges.Load(), dispatched.Load())
	}
	identity.mu.Lock()
	if identity.cached.value != "" {
		t.Error("unscoped unauthorized request did not invalidate its rejected bearer")
	}
	identity.mu.Unlock()
}

func TestX509WorkloadIdentityMiddlewareCanReplayUnscopedBodylessRequests(t *testing.T) {
	var exchanges, dispatched, restored atomic.Int32
	fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		token := fmt.Sprintf("synthetic-direct-token-%d", exchanges.Add(1))
		_, _ = io.WriteString(w, strings.Replace(x509ValidExchangeResponse(), x509ExchangeSyntheticToken, token, 1))
	}))
	identity := newX509LifecycleIdentity(t, fixture)
	request, err := http.NewRequestWithContext(t.Context(), http.MethodGet,
		"https://"+x509APIHost+"/v1/models", nil)
	if err != nil {
		t.Fatalf("construct bodyless direct middleware request: %v", err)
	}
	request.GetBody = func() (io.ReadCloser, error) {
		restored.Add(1)
		return io.NopCloser(strings.NewReader("synthetic-removed-private-body")), nil
	}
	response, err := X509WorkloadIdentityMiddleware(identity, fixture.capability, request,
		func(authenticated *http.Request) (*http.Response, error) {
			if authenticated.Body != nil {
				t.Error("bodyless direct middleware replay restored a removed request body")
			}
			status := http.StatusOK
			if dispatched.Add(1) == 1 {
				status = http.StatusUnauthorized
			}
			return &http.Response{StatusCode: status, Body: io.NopCloser(strings.NewReader("synthetic"))}, nil
		})
	if err != nil || response == nil || response.StatusCode != http.StatusOK {
		t.Fatalf("unscoped bodyless middleware replay = response:%v error:%v", response, err)
	}
	if closeErr := response.Body.Close(); closeErr != nil {
		t.Fatalf("close direct bodyless response: %v", closeErr)
	}
	if exchanges.Load() != 2 || dispatched.Load() != 2 {
		t.Errorf("unscoped bodyless middleware issuer/dispatch attempts = %d/%d", exchanges.Load(), dispatched.Load())
	}
	if got := restored.Load(); got != 0 {
		t.Errorf("bodyless middleware replay invoked its retained GetBody factory %d times", got)
	}
}

func TestX509WorkloadIdentityMiddlewareReplaysUnauthorizedResponsesWithoutBodies(t *testing.T) {
	for _, test := range []struct {
		name   string
		scoped bool
		noBody bool
	}{
		{name: "direct middleware"},
		{name: "request-scoped middleware", scoped: true},
		{name: "direct middleware with http.NoBody", noBody: true},
		{name: "request-scoped middleware with http.NoBody", scoped: true, noBody: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			var exchanges, dispatched, restored atomic.Int32
			fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				token := fmt.Sprintf("synthetic-bodyless-response-token-%d", exchanges.Add(1))
				_, _ = io.WriteString(w, strings.Replace(x509ValidExchangeResponse(), x509ExchangeSyntheticToken, token, 1))
			}))
			identity := newX509LifecycleIdentity(t, fixture)
			request, err := http.NewRequestWithContext(t.Context(), http.MethodGet,
				"https://"+x509APIHost+"/v1/models", nil)
			if err != nil {
				t.Fatalf("construct request for bodyless unauthorized response: %v", err)
			}
			if test.noBody {
				request.Body = http.NoBody
				request.GetBody = func() (io.ReadCloser, error) {
					restored.Add(1)
					return io.NopCloser(strings.NewReader("synthetic-removed-private-body")), nil
				}
			}
			if test.scoped {
				scope := requestconfig.NewRequestRetryScope(1, 0, false, nil)
				if !scope.BeginAttempt() {
					t.Fatal("begin initial scoped middleware attempt")
				}
				request = request.WithContext(requestconfig.WithRequestRetryScope(request.Context(), scope))
			}
			response, err := X509WorkloadIdentityMiddleware(identity, fixture.capability, request,
				func(authenticated *http.Request) (*http.Response, error) {
					if test.noBody && authenticated.Body != http.NoBody {
						t.Error("http.NoBody direct middleware request restored a removed payload")
					}
					attempt := dispatched.Add(1)
					want := fmt.Sprintf("Bearer synthetic-bodyless-response-token-%d", attempt)
					if got := authenticated.Header.Get("Authorization"); got != want {
						t.Errorf("authenticated bodyless response attempt %d bearer = %q, want %q", attempt, got, want)
					}
					status := http.StatusOK
					if attempt == 1 {
						status = http.StatusUnauthorized
					}
					return &http.Response{StatusCode: status}, nil
				})
			if response != nil && response.Body != nil {
				if closeErr := response.Body.Close(); closeErr != nil {
					t.Errorf("close unexpected synthetic response body: %v", closeErr)
				}
			}
			wantStatus := http.StatusOK
			wantAttempts := int32(2)
			if test.scoped {
				wantStatus = http.StatusUnauthorized
				wantAttempts = 1
				if response == nil || response.Header.Get("x-should-retry") != "true" {
					t.Error("scoped unauthorized response did not request a complete outer SDK retry")
				}
			}
			if err != nil || response == nil || response.StatusCode != wantStatus || response.Body != nil {
				t.Fatalf("bodyless unauthorized response replay = response:%v error:%v", response, err)
			}
			if exchanges.Load() != wantAttempts || dispatched.Load() != wantAttempts {
				t.Errorf("bodyless unauthorized response issuer/dispatch attempts = %d/%d", exchanges.Load(), dispatched.Load())
			}
			if got := request.Header.Get("Authorization"); got != "" {
				t.Errorf("direct middleware mutated caller-owned request Authorization: %q", got)
			}
			if got := restored.Load(); got != 0 {
				t.Errorf("bodyless middleware request invoked its retained GetBody factory %d times", got)
			}
		})
	}
}

func TestX509WorkloadIdentityMiddlewareNeverReturnsSignedResponseRequests(t *testing.T) {
	for _, test := range []struct {
		name           string
		status         int
		dispatchError  bool
		scoped         bool
		nilRequest     bool
		nilHeader      bool
		unauthorized   bool
		wantDispatches int32
	}{
		{name: "successful response", status: http.StatusOK, wantDispatches: 1},
		{name: "API error response", status: http.StatusForbidden, wantDispatches: 1},
		{name: "response and transport error", status: http.StatusOK, dispatchError: true, wantDispatches: 1},
		{name: "response without request", status: http.StatusOK, nilRequest: true, wantDispatches: 1},
		{name: "response request without headers", status: http.StatusOK, nilHeader: true, wantDispatches: 1},
		{name: "scoped unauthorized response", status: http.StatusUnauthorized, scoped: true, wantDispatches: 1},
		{name: "unscoped unauthorized replay", status: http.StatusOK, unauthorized: true, wantDispatches: 2},
	} {
		t.Run(test.name, func(t *testing.T) {
			var exchanges, dispatches atomic.Int32
			fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				token := fmt.Sprintf("synthetic-private-response-token-%d", exchanges.Add(1))
				_, _ = io.WriteString(w, strings.Replace(x509ValidExchangeResponse(), x509ExchangeSyntheticToken, token, 1))
			}))
			identity := newX509LifecycleIdentity(t, fixture)
			request, err := http.NewRequestWithContext(t.Context(), http.MethodGet,
				"https://"+x509APIHost+"/v1/models", nil)
			if err != nil {
				t.Fatalf("construct unsigned middleware request: %v", err)
			}
			request.Header.Set("X-Synthetic-Metadata", "preserved")
			if test.scoped {
				scope := requestconfig.NewRequestRetryScope(1, 0, true, nil)
				if !scope.BeginAttempt() {
					t.Fatal("begin request-scoped authentication attempt")
				}
				request = request.WithContext(requestconfig.WithRequestRetryScope(request.Context(), scope))
			}

			var wireResponses []*http.Response
			var returnedError = errors.New("synthetic dispatch failure")
			response, middlewareErr := X509WorkloadIdentityMiddleware(identity, fixture.capability, request,
				func(authenticated *http.Request) (*http.Response, error) {
					attempt := dispatches.Add(1)
					if got := authenticated.Header.Get("Authorization"); got != fmt.Sprintf("Bearer synthetic-private-response-token-%d", attempt) {
						t.Errorf("authenticated wire attempt %d bearer = %q", attempt, got)
					}
					status := test.status
					if test.unauthorized && attempt == 1 {
						status = http.StatusUnauthorized
					}
					wire := &http.Response{
						StatusCode: status,
						Header:     http.Header{"X-Synthetic-Response": {"preserved"}},
						Body:       io.NopCloser(strings.NewReader("synthetic-response-body")),
						Request:    authenticated,
					}
					if test.nilRequest {
						wire.Request = nil
					} else if test.nilHeader {
						wire.Request = authenticated.Clone(authenticated.Context())
						wire.Request.Header = nil
					} else {
						wire.Request.Header["aUtHoRiZaTiOn"] = []string{"Bearer synthetic-private-alias"}
					}
					wireResponses = append(wireResponses, wire)
					if test.dispatchError {
						return wire, returnedError
					}
					return wire, nil
				})
			if response == nil {
				t.Fatalf("signed middleware response unexpectedly nil: %v", middlewareErr)
			}
			if closeErr := response.Body.Close(); closeErr != nil {
				t.Errorf("close returned middleware response: %v", closeErr)
			}
			if test.dispatchError != errors.Is(middlewareErr, returnedError) {
				t.Errorf("middleware dispatch error = %v", middlewareErr)
			}
			if response.StatusCode != test.status || response.Header.Get("X-Synthetic-Response") != "preserved" {
				t.Errorf("returned response lost its status or headers: %v", response)
			}
			if test.nilRequest {
				if response.Request != nil {
					t.Error("response unexpectedly acquired a request")
				}
			} else if test.nilHeader {
				if response.Request == nil || response.Request.Header != nil {
					t.Error("response request unexpectedly acquired a header map")
				}
			} else {
				if response.Request == nil || response.Request.Header.Get("X-Synthetic-Metadata") != "preserved" {
					t.Fatal("response lost its unsigned request metadata")
				}
				for name := range response.Request.Header {
					if strings.EqualFold(strings.ReplaceAll(name, "_", "-"), "authorization") {
						t.Errorf("returned response exposed signed %q request credentials", name)
					}
				}
				for _, wire := range wireResponses {
					if wire.Request == nil || wire.Request.Header.Get("Authorization") == "" {
						t.Error("sanitizing the returned response modified its original signed wire request")
					}
				}
			}
			if test.scoped && response.Header.Get("x-should-retry") != "true" {
				t.Error("scoped unauthorized response lost its outer-retry signal")
			}
			if got := dispatches.Load(); got != test.wantDispatches {
				t.Errorf("authenticated dispatches = %d, want %d", got, test.wantDispatches)
			}
			if request.Header.Get("Authorization") != "" {
				t.Error("authentication modified the original unsigned caller request")
			}
		})
	}
}

func TestX509WorkloadIdentityMiddlewareRetainsUnscopedIssuerRefusal(t *testing.T) {
	var exchanges, dispatched atomic.Int32
	var allowIssuer atomic.Bool
	fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		exchanges.Add(1)
		if allowIssuer.Load() {
			_, _ = io.WriteString(w, x509ValidExchangeResponse())
			return
		}
		w.Header().Set("Retry-After", "90")
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	identity := newX509LifecycleIdentity(t, fixture)
	identity.cached = x509ExchangedToken{value: "synthetic-cached-bearer", expiresAt: time.Now().Add(time.Minute)}
	request, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "https://"+x509APIHost+"/v1/models", nil)
	if err != nil {
		t.Fatal(err)
	}
	response, err := X509WorkloadIdentityMiddleware(identity, fixture.capability, request, func(*http.Request) (*http.Response, error) {
		dispatched.Add(1)
		return &http.Response{StatusCode: http.StatusUnauthorized}, nil
	})

	if response != nil && response.Body != nil {
		if closeErr := response.Body.Close(); closeErr != nil {
			t.Errorf("close standalone middleware response: %v", closeErr)
		}
	}
	var status *x509ExchangeHTTPError
	if response != nil || !errors.As(err, &status) || status.statusCode != http.StatusServiceUnavailable {
		t.Errorf("standalone cached 401 replay = response:%v error:%v, want original issuer503", response, err)
	}
	if exchanges.Load() != 1 || dispatched.Load() != 1 {
		t.Errorf("standalone refusal issuer/API attempts = %d/%d, want 1/1", exchanges.Load(), dispatched.Load())
	}
	allowIssuer.Store(true)
	response, err = X509WorkloadIdentityMiddleware(identity, fixture.capability, request, func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK}, nil
	})

	if response != nil && response.Body != nil {
		if closeErr := response.Body.Close(); closeErr != nil {
			t.Errorf("close standalone middleware response: %v", closeErr)
		}
	}
	if err != nil || response == nil || response.StatusCode != http.StatusOK {
		t.Errorf("independent middleware invocation retained refusal: %v", err)
	}
	if got := exchanges.Load(); got != 2 {
		t.Errorf("independent middleware invocation issuer count = %d, want 2 total", got)
	}
}

func TestX509WorkloadIdentityMiddlewarePreservesUnscopedIssuerRetryBudget(t *testing.T) {
	var exchanges, dispatched atomic.Int32
	fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		attempt := exchanges.Add(1)
		if attempt%3 != 0 {
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		_, _ = io.WriteString(w, strings.Replace(x509ValidExchangeResponse(), x509ExchangeSyntheticToken, fmt.Sprintf("synthetic-bearer-%d", attempt), 1))
	}))
	identity := newX509LifecycleIdentity(t, fixture)
	request, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "https://"+x509APIHost+"/v1/models", nil)
	if err != nil {
		t.Fatal(err)
	}
	response, err := X509WorkloadIdentityMiddleware(identity, fixture.capability, request, func(*http.Request) (*http.Response, error) {
		status := http.StatusOK
		if dispatched.Add(1) == 1 {
			status = http.StatusUnauthorized
		}
		return &http.Response{StatusCode: status}, nil
	})

	if response != nil && response.Body != nil {
		if closeErr := response.Body.Close(); closeErr != nil {
			t.Errorf("close standalone middleware response: %v", closeErr)
		}
	}
	if err != nil || response == nil || response.StatusCode != http.StatusOK || exchanges.Load() != 6 || dispatched.Load() != 2 {
		t.Errorf("standalone issuer retry budgets: response:%v error:%v issuer/API=%d/%d, want success and6/2", response, err, exchanges.Load(), dispatched.Load())
	}
}

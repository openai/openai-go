package auth

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/openai/openai-go/v3/internal/requestconfig"
)

func TestX509WorkloadIdentityMiddlewareRejectsMalformedRequestsBeforeExchange(t *testing.T) {
	var exchanges atomic.Int32
	fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		exchanges.Add(1)
		_, _ = io.WriteString(w, x509ValidExchangeResponse())
	}))
	identity := newX509LifecycleIdentity(t, fixture)
	request, err := http.NewRequestWithContext(t.Context(), http.MethodGet,
		"https://"+x509APIHost+"/v1/models", nil)
	if err != nil {
		t.Fatalf("construct malformed middleware request: %v", err)
	}
	request.Header = nil
	for _, invalid := range []*http.Request{nil, request} {
		response, middlewareErr := X509WorkloadIdentityMiddleware(identity, fixture.capability, invalid,
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
	}
	if got := exchanges.Load(); got != 0 {
		t.Errorf("malformed direct middleware requests minted %d tokens", got)
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
	for _, scoped := range []bool{false, true} {
		name := "direct middleware"
		if scoped {
			name = "request-scoped middleware"
		}
		t.Run(name, func(t *testing.T) {
			var exchanges, dispatched atomic.Int32
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
			if scoped {
				scope := requestconfig.NewRequestRetryScope(1, 0, false, nil)
				if !scope.BeginAttempt() {
					t.Fatal("begin initial scoped middleware attempt")
				}
				request = request.WithContext(requestconfig.WithRequestRetryScope(request.Context(), scope))
			}
			response, err := X509WorkloadIdentityMiddleware(identity, fixture.capability, request,
				func(authenticated *http.Request) (*http.Response, error) {
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
			if err != nil || response == nil || response.StatusCode != http.StatusOK || response.Body != nil {
				t.Fatalf("bodyless unauthorized response replay = response:%v error:%v", response, err)
			}
			if exchanges.Load() != 2 || dispatched.Load() != 2 {
				t.Errorf("bodyless unauthorized response issuer/dispatch attempts = %d/%d", exchanges.Load(), dispatched.Load())
			}
			if got := request.Header.Get("Authorization"); got != "" {
				t.Errorf("direct middleware mutated caller-owned request Authorization: %q", got)
			}
		})
	}
}

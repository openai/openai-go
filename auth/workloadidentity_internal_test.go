package auth

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

type internalHTTPDoer struct {
	do func(*http.Request) (*http.Response, error)
}

func (d *internalHTTPDoer) Do(req *http.Request) (*http.Response, error) {
	return d.do(req)
}

func TestTokenExchangeRetryDelayHonorsRetryAfter(t *testing.T) {
	t.Run("seconds", func(t *testing.T) {
		resp := &http.Response{Header: http.Header{"Retry-After": []string{"2"}}}
		if got, want := tokenExchangeRetryDelay(resp, 0), 2*time.Second; got != want {
			t.Errorf("tokenExchangeRetryDelay() = %v, want %v", got, want)
		}
	})

	t.Run("capped", func(t *testing.T) {
		for _, retryAfter := range []string{"120", "1e300"} {
			resp := &http.Response{Header: http.Header{"Retry-After": []string{retryAfter}}}
			if got, want := tokenExchangeRetryDelay(resp, 0), tokenExchangeMaxRetryDelay; got != want {
				t.Errorf("tokenExchangeRetryDelay(%q) = %v, want %v", retryAfter, got, want)
			}
		}
	})

	t.Run("far future date", func(t *testing.T) {
		resp := &http.Response{Header: http.Header{"Retry-After": []string{"Fri, 31 Dec 9999 23:59:59 GMT"}}}
		if got, want := tokenExchangeRetryDelay(resp, 0), tokenExchangeMaxRetryDelay; got != want {
			t.Errorf("tokenExchangeRetryDelay() = %v, want %v", got, want)
		}
	})

	t.Run("malformed values remain bounded", func(t *testing.T) {
		for _, retryAfter := range []string{"1e999", "NaN", "+Inf", "-Inf", "+999999999999999999999", "-1"} {
			resp := &http.Response{Header: http.Header{"Retry-After": []string{retryAfter}}}
			if got := tokenExchangeRetryDelay(resp, 0); got < 0 || got > tokenExchangeMaxRetryDelay {
				t.Errorf("tokenExchangeRetryDelay(%q) = %v, want a bounded delay", retryAfter, got)
			}
		}
	})

	t.Run("http date", func(t *testing.T) {
		retryAt := time.Now().Add(5 * time.Second).UTC().Truncate(time.Second)
		resp := &http.Response{Header: http.Header{"Retry-After": []string{retryAt.Format(http.TimeFormat)}}}
		before := time.Until(retryAt)
		got := tokenExchangeRetryDelay(resp, 0)
		after := time.Until(retryAt)
		if got < after || got > before {
			t.Errorf("tokenExchangeRetryDelay() = %v, want between %v and %v", got, after, before)
		}
	})
}

func TestX509RefreshTimeoutCoversRetryDelayBudget(t *testing.T) {
	retryDelayBudget := tokenExchangeMaxRetryDelay * tokenExchangeMaxRetries
	if x509WorkloadIdentityRefreshTimeout <= retryDelayBudget {
		t.Fatalf(
			"X.509 refresh timeout = %v, must exceed retry delay budget %v",
			x509WorkloadIdentityRefreshTimeout,
			retryDelayBudget,
		)
	}
}

func TestCredentialSourcesHaveIndependentRefreshTimeouts(t *testing.T) {
	if got, want := (subjectTokenCredentialSource{}).refreshTimeout(), 30*time.Second; got != want {
		t.Errorf("subject-token refresh timeout = %v, want %v", got, want)
	}
	if got, want := (x509CredentialSource{}).refreshTimeout(), x509WorkloadIdentityRefreshTimeout; got != want {
		t.Errorf("X.509 refresh timeout = %v, want %v", got, want)
	}
}

func TestX509CanceledLeaderDoesNotCancelSharedRefresh(t *testing.T) {
	exchangeStarted := make(chan struct{})
	releaseExchange := make(chan struct{})
	var calls atomic.Int32
	httpClient := &internalHTTPDoer{do: func(req *http.Request) (*http.Response, error) {
		calls.Add(1)
		close(exchangeStarted)
		select {
		case <-releaseExchange:
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       io.NopCloser(strings.NewReader(`{"access_token":"shared-token","expires_in":60}`)),
			}, nil
		case <-req.Context().Done():
			return nil, req.Context().Err()
		}
	}}
	wia, err := NewX509WorkloadIdentityAuth(X509WorkloadIdentity{
		IdentityProviderID: "idp-test",
		ServiceAccountID:   "svc-test",
	})
	if err != nil {
		t.Fatalf("NewX509WorkloadIdentityAuth() error = %v", err)
	}
	leaderCtx, cancelLeader := context.WithCancel(t.Context())
	leaderDone := make(chan error, 1)
	go func() {
		_, getTokenErr := wia.GetToken(leaderCtx, httpClient)
		leaderDone <- getTokenErr
	}()
	<-exchangeStarted

	type tokenResult struct {
		token string
		err   error
	}
	waiterDone := make(chan tokenResult, 1)
	go func() {
		token, getTokenErr := wia.GetToken(t.Context(), httpClient)
		waiterDone <- tokenResult{token: token, err: getTokenErr}
	}()
	deadline := time.Now().Add(time.Second)
	for {
		wia.mu.Lock()
		waiters := wia.refreshInFlight.waiters
		wia.mu.Unlock()
		if waiters == 2 {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("timed out waiting for second refresh waiter")
		}
		time.Sleep(time.Millisecond)
	}

	cancelLeader()
	if err := <-leaderDone; !errors.Is(err, context.Canceled) {
		t.Errorf("leader GetToken() error = %v, want context.Canceled", err)
	}
	select {
	case result := <-waiterDone:
		t.Fatalf("waiter completed before shared exchange release: %+v", result)
	default:
	}

	close(releaseExchange)
	result := <-waiterDone
	if result.err != nil {
		t.Fatalf("waiter GetToken() error = %v", result.err)
	}
	if got, want := result.token, "shared-token"; got != want {
		t.Errorf("waiter GetToken() = %q, want %q", got, want)
	}
	if got, want := calls.Load(), int32(1); got != want {
		t.Errorf("token exchange calls = %d, want %d", got, want)
	}
}

func TestX509SharedRefreshPreservesInitiatorContextValues(t *testing.T) {
	type tenantContextKey struct{}
	const tenant = "tenant-a"
	httpClient := &internalHTTPDoer{do: func(req *http.Request) (*http.Response, error) {
		if got := req.Context().Value(tenantContextKey{}); got != tenant {
			return &http.Response{
				StatusCode: http.StatusBadRequest,
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader(`{"error":"missing_tenant"}`)),
			}, nil
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(`{"access_token":"tenant-token","expires_in":60}`)),
		}, nil
	}}
	wia, err := NewX509WorkloadIdentityAuth(X509WorkloadIdentity{
		IdentityProviderID: "idp-test",
		ServiceAccountID:   "svc-test",
	})
	if err != nil {
		t.Fatalf("NewX509WorkloadIdentityAuth() error = %v", err)
	}
	ctx := context.WithValue(t.Context(), tenantContextKey{}, tenant)

	token, err := wia.GetToken(ctx, httpClient)
	if err != nil {
		t.Fatalf("GetToken() error = %v", err)
	}
	if got, want := token, "tenant-token"; got != want {
		t.Fatalf("GetToken() = %q, want %q", got, want)
	}
}

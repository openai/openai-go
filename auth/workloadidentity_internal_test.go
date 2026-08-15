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

	t.Run("http date", func(t *testing.T) {
		retryAt := time.Now().Add(5 * time.Second).UTC().Truncate(time.Second)
		resp := &http.Response{Header: http.Header{"Retry-After": []string{retryAt.Format(http.TimeFormat)}}}
		got := tokenExchangeRetryDelay(resp, 0)
		if got < 4*time.Second || got > 5*time.Second {
			t.Errorf("tokenExchangeRetryDelay() = %v, want between 4s and 5s", got)
		}
	})
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

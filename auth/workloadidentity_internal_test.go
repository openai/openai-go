package auth

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"sync/atomic"
	"testing"
	"time"
)

type x509RefreshGenerationSource struct {
	calls         atomic.Int32
	firstStarted  chan struct{}
	firstCanceled chan struct{}
}

func (s *x509RefreshGenerationSource) exchange(ctx context.Context, _ HTTPDoer, _, _ string) (exchangedToken, error) {
	if s.calls.Add(1) == 1 {
		close(s.firstStarted)
		<-ctx.Done()
		close(s.firstCanceled)
		return exchangedToken{}, ctx.Err()
	}
	return exchangedToken{accessToken: "replacement-token", expiresIn: time.Hour}, nil
}

func (*x509RefreshGenerationSource) refreshBuffer(time.Duration) time.Duration {
	return time.Minute
}

func (*x509RefreshGenerationSource) refreshTimeout() time.Duration {
	return time.Minute
}

func (*x509RefreshGenerationSource) kind() workloadIdentityCredentialSourceKind {
	return workloadIdentityCredentialSourceX509
}

type internalHTTPDoer struct {
	do func(*http.Request) (*http.Response, error)
}

func (d *internalHTTPDoer) Do(req *http.Request) (*http.Response, error) {
	return d.do(req)
}

type internalX509CredentialSource struct {
	exchangeFunc func(context.Context) (exchangedToken, error)
}

func (s *internalX509CredentialSource) exchange(ctx context.Context, _ HTTPDoer, _, _ string) (exchangedToken, error) {
	return s.exchangeFunc(ctx)
}

func (*internalX509CredentialSource) refreshBuffer(time.Duration) time.Duration {
	return time.Minute
}

func (*internalX509CredentialSource) refreshTimeout() time.Duration {
	return time.Minute
}

func (*internalX509CredentialSource) kind() workloadIdentityCredentialSourceKind {
	return workloadIdentityCredentialSourceX509
}

func internalNativeX509HTTPClient() *http.Client {
	return &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{
		Certificates: []tls.Certificate{{Certificate: [][]byte{{1}}}},
	}}}
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
	source := &internalX509CredentialSource{exchangeFunc: func(ctx context.Context) (exchangedToken, error) {
		calls.Add(1)
		close(exchangeStarted)
		select {
		case <-releaseExchange:
			return exchangedToken{accessToken: "shared-token", expiresIn: time.Minute}, nil
		case <-ctx.Done():
			return exchangedToken{}, ctx.Err()
		}
	}}
	httpClient := internalNativeX509HTTPClient()
	wia := newWorkloadIdentityAuth("idp-test", "svc-test", source)
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

func TestX509CanceledOnlyWaiterCancelsSharedRefresh(t *testing.T) {
	exchangeStarted := make(chan struct{})
	exchangeCanceled := make(chan struct{})
	source := &internalX509CredentialSource{exchangeFunc: func(ctx context.Context) (exchangedToken, error) {
		close(exchangeStarted)
		<-ctx.Done()
		close(exchangeCanceled)
		return exchangedToken{}, ctx.Err()
	}}
	wia := newWorkloadIdentityAuth("idp-test", "svc-test", source)
	ctx, cancel := context.WithCancel(t.Context())
	done := make(chan error, 1)
	go func() {
		_, getTokenErr := wia.GetToken(ctx, internalNativeX509HTTPClient())
		done <- getTokenErr
	}()
	<-exchangeStarted
	cancel()
	if err := <-done; !errors.Is(err, context.Canceled) {
		t.Errorf("GetToken() error = %v, want context.Canceled", err)
	}
	select {
	case <-exchangeCanceled:
	case <-time.After(time.Second):
		t.Fatal("shared exchange context was not canceled")
	}
}

func TestX509PreCanceledContextDoesNotStartRefresh(t *testing.T) {
	exchangeStarted := make(chan struct{}, 1)
	source := &internalX509CredentialSource{exchangeFunc: func(ctx context.Context) (exchangedToken, error) {
		exchangeStarted <- struct{}{}
		<-ctx.Done()
		return exchangedToken{}, ctx.Err()
	}}
	wia := newWorkloadIdentityAuth("idp-test", "svc-test", source)
	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	_, err := wia.GetToken(ctx, internalNativeX509HTTPClient())
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("GetToken() error = %v, want context.Canceled", err)
	}
	select {
	case <-exchangeStarted:
		t.Fatal("token exchange started for an already-canceled context")
	case <-time.After(100 * time.Millisecond):
	}
}

func TestX509SharedRefreshPreservesInitiatorContextValues(t *testing.T) {
	type tenantContextKey struct{}
	const tenant = "tenant-a"
	source := &internalX509CredentialSource{exchangeFunc: func(ctx context.Context) (exchangedToken, error) {
		if got := ctx.Value(tenantContextKey{}); got != tenant {
			return exchangedToken{}, errors.New("missing tenant context")
		}
		return exchangedToken{accessToken: "tenant-token", expiresIn: time.Minute}, nil
	}}
	httpClient := internalNativeX509HTTPClient()
	wia := newWorkloadIdentityAuth("idp-test", "svc-test", source)
	ctx := context.WithValue(t.Context(), tenantContextKey{}, tenant)

	token, err := wia.GetToken(ctx, httpClient)
	if err != nil {
		t.Fatalf("GetToken() error = %v", err)
	}
	if got, want := token, "tenant-token"; got != want {
		t.Fatalf("GetToken() = %q, want %q", got, want)
	}
}

func TestX509InvalidationCancelsObsoleteProactiveRefresh(t *testing.T) {
	source := &x509RefreshGenerationSource{
		firstStarted:  make(chan struct{}),
		firstCanceled: make(chan struct{}),
	}
	wia := newWorkloadIdentityAuth("idp-test", "svc-test", source)
	wia.cachedToken = "stale-token"
	wia.tokenExpiry = time.Now().Add(time.Hour)
	wia.tokenRefreshAt = time.Now().Add(-time.Second)
	wia.tokenGeneration = 1
	httpClient := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{
		Certificates: []tls.Certificate{{Certificate: [][]byte{{1}}}},
	}}}

	token, err := wia.GetToken(t.Context(), httpClient)
	if err != nil {
		t.Fatalf("proactive GetToken() error = %v", err)
	}
	if token != "stale-token" {
		t.Fatalf("proactive GetToken() = %q, want stale-token", token)
	}
	<-source.firstStarted

	wia.invalidateToken("stale-token")
	ctx, cancel := context.WithTimeout(t.Context(), time.Second)
	defer cancel()
	token, err = wia.GetToken(ctx, httpClient)
	if err != nil {
		t.Fatalf("forced GetToken() error = %v", err)
	}
	if token != "replacement-token" {
		t.Fatalf("forced GetToken() = %q, want replacement-token", token)
	}
	select {
	case <-source.firstCanceled:
	case <-time.After(time.Second):
		t.Fatal("obsolete proactive refresh was not canceled")
	}
	if got, want := source.calls.Load(), int32(2); got != want {
		t.Fatalf("exchange calls = %d, want %d", got, want)
	}
}

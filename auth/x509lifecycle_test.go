package auth

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestX509WorkloadIdentitySharesConcurrentRefreshes(t *testing.T) {
	var exchanges atomic.Int32
	release := make(chan struct{})
	fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		exchanges.Add(1)
		<-release
		_, _ = io.WriteString(w, x509ValidExchangeResponse())
	}))
	identity := newX509LifecycleIdentity(t, fixture)
	const callers = 32
	var ready, finished sync.WaitGroup
	ready.Add(callers)
	finished.Add(callers)
	start := make(chan struct{})
	results := make(chan error, callers)
	for range callers {
		go func() {
			defer finished.Done()
			ready.Done()
			<-start
			token, err := identity.GetToken(t.Context(), fixture.capability)
			if err == nil && token != x509ExchangeSyntheticToken {
				err = errors.New("concurrent refresh returned an unexpected token")
			}
			results <- err
		}()
	}
	ready.Wait()
	close(start)
	close(release)
	finished.Wait()
	close(results)
	for err := range results {
		if err != nil {
			t.Errorf("concurrent token request: %v", err)
		}
	}
	if got := exchanges.Load(); got != 1 {
		t.Errorf("%d concurrent token callers performed %d exchanges", callers, got)
	}
}

func TestX509WorkloadIdentityRecoversAfterCanceledRefreshLeader(t *testing.T) {
	var exchanges atomic.Int32
	firstReached := make(chan struct{})
	releaseFirst := make(chan struct{})
	defer close(releaseFirst)
	fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if exchanges.Add(1) == 1 {
			close(firstReached)
			<-releaseFirst
			return
		}
		_, _ = io.WriteString(w, x509ValidExchangeResponse())
	}))
	identity := newX509LifecycleIdentity(t, fixture)
	leaderContext, cancelLeader := context.WithCancel(t.Context())
	defer cancelLeader()
	leaderResult := make(chan error, 1)
	go func() {
		_, err := identity.GetToken(leaderContext, fixture.capability)
		leaderResult <- err
	}()
	select {
	case <-firstReached:
	case <-time.After(5 * time.Second):
		t.Fatal("initial refresh leader never reached the issuer")
	}
	followerResult := make(chan error, 1)
	go func() {
		token, err := identity.GetToken(t.Context(), fixture.capability)
		if err == nil && token != x509ExchangeSyntheticToken {
			err = errors.New("healthy refresh follower received an unexpected token")
		}
		followerResult <- err
	}()
	cancelLeader()
	select {
	case err := <-leaderResult:
		if !errors.Is(err, context.Canceled) {
			t.Errorf("canceled leader error = %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("canceled refresh leader did not finish")
	}
	select {
	case err := <-followerResult:
		if err != nil {
			t.Errorf("healthy follower did not recover after leader cancellation: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("healthy refresh follower was poisoned by the canceled leader")
	}
	if got := exchanges.Load(); got != 2 {
		t.Errorf("canceled leader and healthy follower made %d exchanges, want exactly two", got)
	}
}

func TestX509WorkloadIdentityRefreshBufferHandlesShortLifetimes(t *testing.T) {
	var exchanges atomic.Int32
	fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		number := exchanges.Add(1)
		body := strings.Replace(x509ValidExchangeResponse(), x509ExchangeSyntheticToken,
			"synthetic-short-token-"+strconv.Itoa(int(number)), 1)
		body = strings.Replace(body, `"expires_in":60`, `"expires_in":1`, 1)
		_, _ = io.WriteString(w, body)
	}))
	identity := newX509LifecycleIdentity(t, fixture)
	first, err := identity.GetToken(t.Context(), fixture.capability)
	if err != nil {
		t.Fatalf("acquire short-lived token: %v", err)
	}
	second, err := identity.GetToken(t.Context(), fixture.capability)
	if err != nil || first != second || exchanges.Load() != 1 {
		t.Fatalf("short-lived token was not cached: first=%q second=%q exchanges=%d error=%v",
			first, second, exchanges.Load(), err)
	}
	identity.mu.Lock()
	identity.refreshAfter = time.Now().Add(-time.Millisecond)
	identity.mu.Unlock()
	third, err := identity.GetToken(t.Context(), fixture.capability)
	if err != nil || third == first || exchanges.Load() != 2 {
		t.Errorf("expired refresh window returned token=%q exchanges=%d error=%v", third, exchanges.Load(), err)
	}
}

func TestX509WorkloadIdentityRetriesOnlyBoundedTransientIssuerFailures(t *testing.T) {
	for _, test := range []struct {
		name         string
		status       int
		failures     int32
		wantAttempts int32
		wantSuccess  bool
		truncated    bool
	}{
		{name: "rate limited then succeeds", status: 429, failures: 2, wantAttempts: 3, wantSuccess: true},
		{name: "server error then succeeds", status: 500, failures: 2, wantAttempts: 3, wantSuccess: true},
		{name: "service unavailable exhausts budget", status: 503, failures: 3, wantAttempts: 3},
		{name: "body read failure then succeeds", failures: 2, wantAttempts: 3, wantSuccess: true, truncated: true},
		{name: "permanent OAuth failure", status: 400, failures: 3, wantAttempts: 1},
		{name: "permanent HTTP status", status: 404, failures: 3, wantAttempts: 1},
	} {
		t.Run(test.name, func(t *testing.T) {
			var attempts atomic.Int32
			fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				attempt := attempts.Add(1)
				if attempt <= test.failures {
					if test.truncated {
						w.Header().Set("Content-Length", "1024")
						_, _ = io.WriteString(w, "synthetic-truncated-issuer-response")
						return
					}
					w.WriteHeader(test.status)
					_, _ = io.WriteString(w, `{"error":"invalid_grant"}`)
					return
				}
				_, _ = io.WriteString(w, x509ValidExchangeResponse())
			}))
			identity := newX509LifecycleIdentity(t, fixture)
			token, err := identity.GetToken(t.Context(), fixture.capability)
			if (err == nil) != test.wantSuccess {
				t.Errorf("transient retry result token=%q error=%v", token, err)
			}
			if got := attempts.Load(); got != test.wantAttempts {
				t.Errorf("issuer attempts = %d, want %d", got, test.wantAttempts)
			}
		})
	}
}

func TestX509WorkloadIdentityCancelsRetryBackoff(t *testing.T) {
	var attempts atomic.Int32
	reached := make(chan struct{})
	fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if attempts.Add(1) == 1 {
			close(reached)
		}
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	identity := newX509LifecycleIdentity(t, fixture)
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	result := make(chan error, 1)
	go func() {
		_, err := identity.GetToken(ctx, fixture.capability)
		result <- err
	}()
	select {
	case <-reached:
		cancel()
	case <-time.After(5 * time.Second):
		t.Fatal("retryable issuer did not receive the first exchange")
	}
	select {
	case err := <-result:
		if !errors.Is(err, context.Canceled) {
			t.Errorf("canceled retry backoff error = %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("retry backoff did not respect request cancellation")
	}
	if got := attempts.Load(); got != 1 {
		t.Errorf("canceled retry backoff made %d issuer attempts", got)
	}
}

func TestX509WorkloadIdentityCacheIsGenerationScoped(t *testing.T) {
	var exchanges atomic.Int32
	fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		exchanges.Add(1)
		_, _ = io.WriteString(w, x509ValidExchangeResponse())
	}))
	for range 2 {
		identity := newX509LifecycleIdentity(t, fixture)
		if _, err := identity.GetToken(t.Context(), fixture.capability); err != nil {
			t.Fatalf("independent X.509 identity exchange: %v", err)
		}
	}
	if got := exchanges.Load(); got != 2 {
		t.Errorf("independent identity generations shared a cache: exchanges=%d", got)
	}
}

func TestX509WorkloadIdentityInvalidatesOnlyTheRejectedBearer(t *testing.T) {
	fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, x509ValidExchangeResponse())
	}))
	identity := newX509LifecycleIdentity(t, fixture)
	token, err := identity.GetToken(t.Context(), fixture.capability)
	if err != nil {
		t.Fatalf("prime X.509 bearer cache: %v", err)
	}
	identity.invalidateToken("stale-rejected-token")
	identity.mu.Lock()
	if identity.cached.value != token {
		t.Error("a stale concurrent 401 invalidated the newer cached bearer")
	}
	identity.mu.Unlock()
	identity.invalidateToken(token)
	identity.mu.Lock()
	if identity.cached.value != "" || !identity.refreshAfter.IsZero() {
		t.Error("the rejected current bearer was not removed from the cache")
	}
	identity.mu.Unlock()
}

func newX509LifecycleIdentity(t *testing.T, fixture *x509ExchangeFixture) *X509WorkloadIdentityAuth {
	t.Helper()
	identity, err := NewX509WorkloadIdentityAuth(X509WorkloadIdentity{
		IdentityProviderID: "synthetic-identity-provider",
		ServiceAccountID:   "synthetic-service-account",
		Transport:          fixture.capability,
	})
	if err != nil {
		t.Fatalf("construct ephemeral X.509 lifecycle identity: %v", err)
	}
	return identity
}

package auth

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/openai/openai-go/v3/internal/requestconfig"
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
	for _, hint := range []time.Duration{0, 100 * time.Millisecond} {
		for _, retries := range []int{-1, 0, 1} {
			t.Run(hint.String()+"/retries="+strconv.Itoa(retries), func(t *testing.T) {
				var exchanges, consumed atomic.Int32
				var firstAt atomic.Int64
				firstReached := make(chan struct{})
				releaseFirst := make(chan struct{})
				defer close(releaseFirst)
				fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					if exchanges.Add(1) == 1 {
						firstAt.Store(time.Now().UnixNano())
						if hint > 0 {
							w.Header().Set("Retry-After-Ms", strconv.FormatInt(hint.Milliseconds(), 10))
							w.WriteHeader(http.StatusServiceUnavailable)
							w.(http.Flusher).Flush()
						}
						close(firstReached)
						<-releaseFirst
						return
					}
					if elapsed := time.Since(time.Unix(0, firstAt.Load())); elapsed < hint {
						t.Errorf("healthy follower retried after %s, before %s", elapsed, hint)
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
				followerContext := t.Context()
				if retries >= 0 {
					followerContext = requestconfig.WithRequestRetryScope(followerContext, requestconfig.NewRequestRetryScope(retries, time.Second, true, func(n int) { consumed.Store(int32(n)) }))
				}
				waiting := make(chan struct{})
				followerContext = &x509ObservedDoneContext{Context: followerContext, observed: waiting}
				go func() {
					token, err := identity.GetToken(followerContext, fixture.capability)
					if err == nil && token != x509ExchangeSyntheticToken {
						err = errors.New("healthy refresh follower received an unexpected token")
					}
					followerResult <- err
				}()
				select {
				case <-waiting:
				case <-time.After(5 * time.Second):
					t.Fatal("healthy follower did not start waiting")
				}
				// Let the follower enter the shared-refresh wait before canceling its owner.
				time.Sleep(10 * time.Millisecond)
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
				if consumed.Load() != 0 {
					t.Errorf("leader cancellation consumed %d follower retries", consumed.Load())
				}

			})
		}
	}
}

func TestX509WorkloadIdentityFollowersUseTheirOwnRetryBudgets(t *testing.T) {
	for _, test := range []struct {
		name          string
		status        int
		followerLimit int
		wantExchanges int32
		wantSuccess   bool
	}{
		{name: "retry-enabled follower recovers", status: http.StatusServiceUnavailable,
			followerLimit: 2, wantExchanges: 2, wantSuccess: true},
		{name: "zero-retry follower fails", status: http.StatusServiceUnavailable, wantExchanges: 1},
		{name: "permanent failure is not retried", status: http.StatusBadRequest,
			followerLimit: 2, wantExchanges: 1},
	} {
		t.Run(test.name, func(t *testing.T) {
			var exchanges atomic.Int32
			firstReached := make(chan struct{})
			releaseFirst := make(chan struct{})
			fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				if exchanges.Add(1) == 1 {
					close(firstReached)
					<-releaseFirst
					w.WriteHeader(test.status)
					_, _ = io.WriteString(w, `{"error":{"code":"invalid_grant"}}`)
					return
				}
				_, _ = io.WriteString(w, x509ValidExchangeResponse())
			}))
			identity := newX509LifecycleIdentity(t, fixture)
			leaderScope := requestconfig.NewRequestRetryScope(0, time.Millisecond, true, nil)
			leaderContext := requestconfig.WithRequestRetryScope(t.Context(), leaderScope)
			if !leaderScope.BeginAttempt() {
				t.Fatal("zero-retry leader initial attempt was rejected")
			}
			leaderResult := make(chan error, 1)
			go func() {
				_, err := identity.GetToken(leaderContext, fixture.capability)
				leaderResult <- err
			}()
			select {
			case <-firstReached:
			case <-time.After(5 * time.Second):
				t.Fatal("zero-retry refresh leader never reached the issuer")
			}
			followerScope := requestconfig.NewRequestRetryScope(test.followerLimit, time.Millisecond, true, nil)
			followerParent := requestconfig.WithRequestRetryScope(t.Context(), followerScope)
			if !followerScope.BeginAttempt() {
				t.Fatal("follower initial attempt was rejected")
			}
			waiting := make(chan struct{})
			followerContext := &x509ObservedDoneContext{Context: followerParent, observed: waiting}
			followerResult := make(chan error, 1)
			go func() {
				token, err := identity.GetToken(followerContext, fixture.capability)
				if err == nil && token != x509ExchangeSyntheticToken {
					err = errors.New("retry-enabled follower received an unexpected bearer")
				}
				followerResult <- err
			}()
			select {
			case <-waiting:
				close(releaseFirst)
			case <-time.After(5 * time.Second):
				t.Fatal("refresh follower never waited for its shared leader")
			}
			select {
			case err := <-leaderResult:
				if err == nil {
					t.Error("zero-retry refresh leader unexpectedly recovered")
				}
			case <-time.After(5 * time.Second):
				t.Fatal("zero-retry refresh leader did not finish")
			}
			select {
			case err := <-followerResult:
				if (err == nil) != test.wantSuccess {
					t.Errorf("shared-failure follower error = %v, want success %t", err, test.wantSuccess)
				}
			case <-time.After(5 * time.Second):
				t.Fatal("refresh follower did not finish after its shared failure")
			}
			if got := exchanges.Load(); got != test.wantExchanges {
				t.Errorf("leader/follower issuer attempts = %d, want %d", got, test.wantExchanges)
			}
		})
	}
}

type x509ObservedDoneContext struct {
	context.Context
	observed chan struct{}
	once     sync.Once
}

func (ctx *x509ObservedDoneContext) Done() <-chan struct{} {
	ctx.once.Do(func() { close(ctx.observed) })
	return ctx.Context.Done()
}

func TestX509WorkloadIdentityNeverCachesExpiredExchange(t *testing.T) {
	reached := make(chan struct{})
	release := make(chan struct{})
	written := make(chan struct{})
	fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		close(reached)
		<-release
		body := strings.Replace(x509ValidExchangeResponse(), `"expires_in":60`, `"expires_in":1`, 1)
		_, _ = io.WriteString(w, body)
		if flush, ok := w.(http.Flusher); ok {
			flush.Flush()
		}
		close(written)
	}))
	identity := newX509LifecycleIdentity(t, fixture)
	started := time.Now()
	result := make(chan error, 1)
	go func() {
		token, err := identity.GetToken(t.Context(), fixture.capability)
		if token != "" {
			err = errors.New("expired exchange returned a bearer")
		}
		result <- err
	}()
	select {
	case <-reached:
	case <-time.After(5 * time.Second):
		t.Fatal("short-lived exchange never reached the issuer")
	}
	identity.mu.Lock()
	close(release)
	select {
	case <-written:
	case <-time.After(5 * time.Second):
		identity.mu.Unlock()
		t.Fatal("short-lived issuer did not finish its response")
	}
	timer := time.NewTimer(time.Until(started.Add(1200 * time.Millisecond)))
	<-timer.C
	identity.mu.Unlock()
	select {
	case err := <-result:
		if err == nil || !strings.Contains(err.Error(), "expired") {
			t.Errorf("exchange expiring before cache insertion returned error = %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("expired exchange did not finish")
	}
	identity.mu.Lock()
	if identity.cached.value != "" || !identity.refreshAfter.IsZero() {
		t.Error("already-expired exchange was inserted into the workload token cache")
	}
	identity.mu.Unlock()
}

func TestX509WorkloadIdentityRejectsExpiredCachedTokenDespiteFutureRefresh(t *testing.T) {
	var exchanges atomic.Int32
	fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		exchanges.Add(1)
		_, _ = io.WriteString(w, x509ValidExchangeResponse())
	}))
	identity := newX509LifecycleIdentity(t, fixture)
	identity.mu.Lock()
	identity.cached = x509ExchangedToken{value: "expired-synthetic-bearer", expiresAt: time.Now().Add(-time.Second)}
	identity.refreshAfter = time.Now().Add(time.Hour)
	identity.mu.Unlock()
	token, err := identity.GetToken(t.Context(), fixture.capability)
	if err != nil || token != x509ExchangeSyntheticToken || exchanges.Load() != 1 {
		t.Errorf("expired cached bearer result = %q, exchanges = %d, error = %v", token, exchanges.Load(), err)
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
		oauthCode    string
		failures     int32
		wantAttempts int32
		wantSuccess  bool
		truncated    bool
	}{
		{name: "request timeout then succeeds", status: 408, failures: 2, wantAttempts: 3, wantSuccess: true},
		{name: "conflict then succeeds", status: 409, failures: 2, wantAttempts: 3, wantSuccess: true},
		{name: "rate limited then succeeds", status: 429, failures: 2, wantAttempts: 3, wantSuccess: true},
		{name: "server error then succeeds", status: 500, failures: 2, wantAttempts: 3, wantSuccess: true},
		{name: "service unavailable exhausts budget", status: 503, failures: 3, wantAttempts: 3},
		{name: "body read failure then succeeds", failures: 2, wantAttempts: 3, wantSuccess: true, truncated: true},
		{name: "temporarily unavailable OAuth error then succeeds", status: 400,
			oauthCode: "temporarily_unavailable", failures: 2, wantAttempts: 3, wantSuccess: true},
		{name: "server error OAuth error then succeeds", status: 401,
			oauthCode: "server_error", failures: 2, wantAttempts: 3, wantSuccess: true},
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
					code := test.oauthCode
					if code == "" {
						code = "invalid_grant"
					}
					_, _ = io.WriteString(w, `{"error":"`+code+`"}`)
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

func TestX509WorkloadIdentitySelectsBoundedIssuerRetryDelays(t *testing.T) {
	for _, test := range []struct {
		name    string
		err     error
		attempt int
		scope   *requestconfig.RequestRetryScope
		want    time.Duration
	}{
		{name: "ordinary exponential delay", err: &x509ExchangeHTTPError{statusCode: 503},
			attempt: 1, want: 2 * x509InitialRetryDelay},
		{name: "issuer-directed delay", err: &x509ExchangeHTTPError{
			statusCode: 429, retryAfter: 70 * time.Millisecond, hasRetryAfter: true},
			want: 70 * time.Millisecond},
		{name: "explicit zero is immediate", err: &x509ExchangeHTTPError{
			statusCode: 429, retryAfter: 0, hasRetryAfter: true}},
		{name: "issuer hint is never shortened", err: &x509ExchangeHTTPError{
			statusCode: 503, retryAfter: time.Second, hasRetryAfter: true},
			scope: requestconfig.NewRequestRetryScope(2, 7*time.Millisecond, true, nil), want: time.Second},
		{name: "ordinary delay obeys caller maximum", err: &x509ExchangeReadError{},
			scope: requestconfig.NewRequestRetryScope(2, time.Millisecond, true, nil), want: time.Millisecond},
		{name: "standalone hint is never shortened", err: &x509ExchangeHTTPError{
			statusCode: 429, retryAfter: time.Hour, hasRetryAfter: true}, want: time.Hour},
	} {
		t.Run(test.name, func(t *testing.T) {
			if got := x509ExchangeRetryDelay(test.err, test.attempt, test.scope); got != test.want {
				t.Errorf("issuer retry delay = %s, want %s", got, test.want)
			}
		})
	}
}

func TestX509WorkloadIdentityCancelsRetryBackoff(t *testing.T) {
	var attempts atomic.Int32
	reached := make(chan struct{})
	fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		attempt := attempts.Add(1)
		w.Header().Set("Retry-After", "8")
		w.WriteHeader(http.StatusTooManyRequests)
		if flush, ok := w.(http.Flusher); ok {
			flush.Flush()
		}
		if attempt == 1 {
			close(reached)
		}
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

func TestX509WorkloadIdentityRetriesRedactedTransientNetworkFailure(t *testing.T) {
	fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, x509ValidExchangeResponse())
	}))
	originalDial := fixture.template.DialContext
	var attempts atomic.Int32
	fixture.template.DialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
		if attempts.Add(1) == 1 {
			return nil, errors.New("synthetic-private-first-dial-failure")
		}
		return originalDial(ctx, network, address)
	}
	transport, err := NewX509Transport(fixture.template)
	if err != nil {
		t.Fatalf("attest transient-dial workload transport: %v", err)
	}
	t.Cleanup(func() { _ = transport.Close() })
	identity, err := NewX509WorkloadIdentityAuth(X509WorkloadIdentity{
		IdentityProviderID: "synthetic-identity-provider",
		ServiceAccountID:   "synthetic-service-account",
		Transport:          transport,
	})
	if err != nil {
		t.Fatalf("construct transient-dial workload identity: %v", err)
	}
	token, err := identity.GetToken(t.Context(), transport)
	if err != nil || token != x509ExchangeSyntheticToken || attempts.Load() != 2 {
		t.Errorf("transient native dial token=%q attempts=%d error=%v", token, attempts.Load(), err)
	}
	for _, permanent := range []error{errX509Redirect, context.Canceled, context.DeadlineExceeded} {
		if retryableX509ExchangeError(&x509TransportError{cause: permanent}) {
			t.Errorf("non-transient native transport cause was retried: %v", permanent)
		}
	}
}

func TestX509WorkloadIdentityDoesNotRetryPermanentTLSVerificationFailures(t *testing.T) {
	const privateVerificationError = "Authorization: Bearer synthetic-private-verification-token"
	for _, test := range []struct {
		name              string
		configure         func(*tls.Config, *atomic.Int32)
		wantVerifications int32
		cached            bool
	}{
		{
			name: "caller connection verification",
			configure: func(config *tls.Config, verifications *atomic.Int32) {
				config.VerifyConnection = func(tls.ConnectionState) error {
					verifications.Add(1)
					return errors.New(privateVerificationError)
				}
			},
			wantVerifications: 1,
		},
		{
			name: "caller peer certificate verification",
			configure: func(config *tls.Config, verifications *atomic.Int32) {
				config.VerifyPeerCertificate = func([][]byte, [][]*x509.Certificate) error {
					verifications.Add(1)
					return errors.New(privateVerificationError)
				}
			},
			wantVerifications: 1,
		},
		{
			name: "native untrusted issuer certificate",
			configure: func(config *tls.Config, _ *atomic.Int32) {
				config.RootCAs = x509.NewCertPool()
			},
		},
		{
			name: "permanent verification failure does not reuse cached token",
			configure: func(config *tls.Config, verifications *atomic.Int32) {
				config.VerifyConnection = func(tls.ConnectionState) error {
					verifications.Add(1)
					return errors.New(privateVerificationError)
				}
			},
			wantVerifications: 1,
			cached:            true,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			var exchanges, attempts, verifications atomic.Int32
			fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				exchanges.Add(1)
				_, _ = io.WriteString(w, x509ValidExchangeResponse())
			}))
			test.configure(fixture.template.TLSClientConfig, &verifications)
			originalDial := fixture.template.DialContext
			fixture.template.DialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
				attempts.Add(1)
				return originalDial(ctx, network, address)
			}
			transport, err := NewX509Transport(fixture.template)
			if err != nil {
				t.Fatalf("attest TLS-verification workload transport: %v", err)
			}
			t.Cleanup(func() { _ = transport.Close() })
			identity, err := NewX509WorkloadIdentityAuth(X509WorkloadIdentity{
				IdentityProviderID: "synthetic-identity-provider",
				ServiceAccountID:   "synthetic-service-account",
				Transport:          transport,
			})
			if err != nil {
				t.Fatalf("construct TLS-verification workload identity: %v", err)
			}
			if test.cached {
				identity.cached = x509ExchangedToken{
					value:     x509ExchangeSyntheticToken,
					expiresAt: time.Now().Add(time.Minute),
				}
				identity.refreshAfter = time.Now().Add(-time.Second)
			}
			token, exchangeErr := identity.GetToken(t.Context(), transport)
			if exchangeErr == nil || token != "" {
				t.Errorf("permanent TLS-verification failure returned token=%q error=%v", token, exchangeErr)
			}
			for cause := exchangeErr; cause != nil; cause = errors.Unwrap(cause) {
				if strings.Contains(cause.Error(), privateVerificationError) ||
					strings.Contains(cause.Error(), "Authorization") {
					t.Errorf("TLS verification error disclosed caller-private credentials: %v", cause)
				}
			}
			if got := attempts.Load(); got != 1 {
				t.Errorf("permanent TLS-verification failure attempted %d issuer handshakes, want 1", got)
			}
			if got := verifications.Load(); got != test.wantVerifications {
				t.Errorf("TLS verification callback executed %d times, want %d", got, test.wantVerifications)
			}
			if got := exchanges.Load(); got != 0 {
				t.Errorf("failed TLS verification reached the issuer handler %d times", got)
			}
		})
	}
}

func TestX509WorkloadIdentityInvalidationPreservesIssuerMinimum(t *testing.T) {
	var exchanges atomic.Int32
	var firstFailure time.Time
	waiting := make(chan struct{})
	fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		switch exchanges.Add(1) {
		case 1:
			_, _ = io.WriteString(w, x509ValidExchangeResponse())
		case 2:
			firstFailure = time.Now()
			w.Header().Set("Retry-After-Ms", "100")
			w.WriteHeader(http.StatusTooManyRequests)
			if flush, ok := w.(http.Flusher); ok {
				flush.Flush()
			}
			close(waiting)
		default:
			if elapsed := time.Since(firstFailure); elapsed < 100*time.Millisecond {
				t.Errorf("invalidation shortened issuer wait to %s, want at least 100ms", elapsed)
			}
			_, _ = io.WriteString(w, strings.Replace(x509ValidExchangeResponse(),
				x509ExchangeSyntheticToken, "synthetic-refreshed-bearer", 1))
		}
	}))
	identity := newX509LifecycleIdentity(t, fixture)
	rejected, err := identity.GetToken(t.Context(), fixture.capability)
	if err != nil {
		t.Fatalf("prime proactive refresh bearer: %v", err)
	}
	identity.mu.Lock()
	identity.refreshAfter = time.Now().Add(-time.Second)
	identity.mu.Unlock()
	ctx, cancel := context.WithTimeout(t.Context(), 750*time.Millisecond)
	defer cancel()
	result := make(chan error, 1)
	go func() {
		token, refreshErr := identity.GetToken(ctx, fixture.capability)
		if refreshErr == nil && token != "synthetic-refreshed-bearer" {
			refreshErr = errors.New("woken proactive refresh returned its obsolete bearer")
		}
		result <- refreshErr
	}()
	select {
	case <-waiting:
	case <-ctx.Done():
		t.Fatal("obsolete proactive refresh never entered issuer-directed backoff")
	}
	identity.invalidateToken("synthetic-other-generation")
	identity.mu.Lock()
	current := identity.inFlight
	if current == nil || current.generation != rejected {
		identity.mu.Unlock()
		t.Fatal("proactive refresh lost its originating cached bearer generation")
	}
	select {
	case <-current.wake:
		identity.mu.Unlock()
		t.Fatal("stale-generation invalidation woke an unrelated proactive refresh")
	default:
	}
	identity.mu.Unlock()
	identity.invalidateToken(rejected)
	select {
	case refreshErr := <-result:
		if refreshErr != nil {
			t.Errorf("generation-specific invalidation failed after issuer-directed wait: %v", refreshErr)
		}
	case <-ctx.Done():
		t.Fatal("rejected cached bearer did not refresh after the issuer minimum")
	}
	if got := exchanges.Load(); got != 3 {
		t.Errorf("prime/obsolete/woken issuer attempts = %d, want 3", got)
	}
}

func TestX509WorkloadIdentityPreservesUnexpiredBearerDuringTransientRefreshFailure(t *testing.T) {
	var exchanges atomic.Int32
	var failureStatus atomic.Int32
	fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		exchanges.Add(1)
		if status := failureStatus.Load(); status != 0 {
			w.WriteHeader(int(status))
			_, _ = io.WriteString(w, `{"error":"invalid_grant"}`)
			return
		}
		_, _ = io.WriteString(w, x509ValidExchangeResponse())
	}))
	identity := newX509LifecycleIdentity(t, fixture)
	initial, err := identity.GetToken(t.Context(), fixture.capability)
	if err != nil {
		t.Fatalf("prime proactive-refresh token cache: %v", err)
	}
	failureStatus.Store(http.StatusServiceUnavailable)
	identity.mu.Lock()
	identity.refreshAfter = time.Now().Add(-time.Second)
	identity.mu.Unlock()
	fallback, err := identity.GetToken(t.Context(), fixture.capability)
	if err != nil || fallback != initial || exchanges.Load() != 4 {
		t.Fatalf("transient proactive refresh fallback=%q exchanges=%d error=%v", fallback, exchanges.Load(), err)
	}
	identity.mu.Lock()
	cooldown, expiry := identity.refreshAfter, identity.cached.expiresAt
	identity.mu.Unlock()
	if !cooldown.After(time.Now()) || cooldown.After(expiry) {
		t.Errorf("proactive retry cooldown = %v, token expiry = %v", cooldown, expiry)
	}
	if cached, cachedErr := identity.GetToken(t.Context(), fixture.capability); cachedErr != nil || cached != initial ||
		exchanges.Load() != 4 {
		t.Errorf("cooldown did not protect still-valid bearer: token=%q exchanges=%d error=%v",
			cached, exchanges.Load(), cachedErr)
	}
	identity.mu.Lock()
	identity.cached.expiresAt = time.Now().Add(-time.Second)
	identity.refreshAfter = time.Now().Add(-time.Second)
	identity.mu.Unlock()
	if token, exchangeErr := identity.GetToken(t.Context(), fixture.capability); token != "" || exchangeErr == nil {
		t.Errorf("expired bearer was used after refresh failure: token=%q error=%v", token, exchangeErr)
	}
}

func TestX509WorkloadIdentityPreservesBearerAfterLiveContextResponseHeaderTimeout(t *testing.T) {
	var exchanges atomic.Int32
	refreshStarted := make(chan struct{})
	releaseRefresh := make(chan struct{})
	defer close(releaseRefresh)
	var refreshOnce sync.Once
	fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if exchanges.Add(1) == 1 {
			_, _ = io.WriteString(w, x509ValidExchangeResponse())
			return
		}
		refreshOnce.Do(func() { close(refreshStarted) })
		<-releaseRefresh
	}))
	fixture.template.ResponseHeaderTimeout = 50 * time.Millisecond
	transport, err := NewX509Transport(fixture.template)
	if err != nil {
		t.Fatalf("attest response-timeout workload transport: %v", err)
	}
	t.Cleanup(func() { _ = transport.Close() })
	identity, err := NewX509WorkloadIdentityAuth(X509WorkloadIdentity{
		IdentityProviderID: "synthetic-identity-provider",
		ServiceAccountID:   "synthetic-service-account",
		Transport:          transport,
	})
	if err != nil {
		t.Fatalf("construct response-timeout workload identity: %v", err)
	}
	initial, err := identity.GetToken(t.Context(), transport)
	if err != nil {
		t.Fatalf("prime response-timeout token cache: %v", err)
	}
	identity.mu.Lock()
	identity.refreshAfter = time.Now().Add(-time.Second)
	identity.mu.Unlock()

	const callers = 8
	results := make(chan error, callers)
	refresh := func() {
		scope := requestconfig.NewRequestRetryScope(0, time.Millisecond, true, nil)
		if !scope.BeginAttempt() {
			results <- errors.New("zero-retry caller could not begin its initial attempt")
			return
		}
		ctx := requestconfig.WithRequestRetryScope(t.Context(), scope)
		token, tokenErr := identity.GetToken(ctx, transport)
		if tokenErr == nil && token != initial {
			tokenErr = errors.New("response-timeout refresh returned an unexpected bearer")
		}
		results <- tokenErr
	}
	go refresh()
	select {
	case <-refreshStarted:
	case <-time.After(5 * time.Second):
		t.Fatal("response-timeout refresh never reached the issuer")
	}
	for range callers - 1 {
		go refresh()
	}
	for range callers {
		select {
		case refreshErr := <-results:
			if refreshErr != nil {
				t.Errorf("response-timeout refresh did not reuse its valid cached bearer: %v", refreshErr)
			}
		case <-time.After(5 * time.Second):
			t.Fatal("response-timeout refresh caller did not finish")
		}
	}
	if got := exchanges.Load(); got != 2 {
		t.Errorf("prime and zero-retry response-timeout exchanges = %d, want 2", got)
	}
}

func TestX509WorkloadIdentityNeverFallsBackAfterPermanentOrInvalidatedFailure(t *testing.T) {
	for _, test := range []struct {
		name       string
		status     int
		invalidate bool
	}{
		{name: "permanent invalid grant", status: http.StatusBadRequest},
		{name: "invalidated bearer", status: http.StatusServiceUnavailable, invalidate: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			var failing atomic.Bool
			fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				if failing.Load() {
					w.WriteHeader(test.status)
					_, _ = io.WriteString(w, `{"error":"invalid_grant"}`)
					return
				}
				_, _ = io.WriteString(w, x509ValidExchangeResponse())
			}))
			identity := newX509LifecycleIdentity(t, fixture)
			token, err := identity.GetToken(t.Context(), fixture.capability)
			if err != nil {
				t.Fatalf("prime X.509 failure cache: %v", err)
			}
			identity.mu.Lock()
			identity.refreshAfter = time.Now().Add(-time.Second)
			identity.mu.Unlock()
			if test.invalidate {
				identity.invalidateToken(token)
			}
			failing.Store(true)
			if cached, exchangeErr := identity.GetToken(t.Context(), fixture.capability); cached != "" || exchangeErr == nil {
				t.Errorf("unsafe fallback returned bearer=%q error=%v", cached, exchangeErr)
			}
		})
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

func TestX509WorkloadIdentityFollowerKeepsOriginalIssuerMinimum(t *testing.T) {
	for _, test := range []struct {
		name        string
		hint        time.Duration
		bodyDelay   time.Duration
		maximum     time.Duration
		retries     int
		wantSuccess bool
	}{
		{name: "elapsed minimum", hint: 100 * time.Millisecond, bodyDelay: 200 * time.Millisecond, maximum: time.Second, retries: 1, wantSuccess: true},
		{name: "remaining minimum", hint: 600 * time.Millisecond, bodyDelay: 400 * time.Millisecond, maximum: time.Second, retries: 1, wantSuccess: true},
		{name: "no hint control", bodyDelay: 200 * time.Millisecond, maximum: time.Second, retries: 1, wantSuccess: true},
		{name: "zero budget control", hint: 100 * time.Millisecond, bodyDelay: 200 * time.Millisecond, maximum: time.Second},
		{name: "over cap control", hint: 100 * time.Millisecond, bodyDelay: 200 * time.Millisecond, maximum: 50 * time.Millisecond, retries: 1},
	} {
		t.Run(test.name, func(t *testing.T) {
			var exchanges, consumed atomic.Int32
			var failedAt atomic.Int64
			reached := make(chan struct{})
			release := make(chan struct{})
			var releaseOnce sync.Once
			defer releaseOnce.Do(func() { close(release) })
			fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				if exchanges.Add(1) == 1 {
					failedAt.Store(time.Now().UnixNano())
					if test.hint > 0 {
						w.Header().Set("Retry-After-Ms", strconv.FormatInt(test.hint.Milliseconds(), 10))
					}
					w.WriteHeader(http.StatusServiceUnavailable)
					w.(http.Flusher).Flush()
					close(reached)
					<-release
					return
				}
				if elapsed := time.Since(time.Unix(0, failedAt.Load())); elapsed < test.hint {
					t.Errorf("follower issuer retry after %s, want at least %s", elapsed, test.hint)
				}
				_, _ = io.WriteString(w, x509ValidExchangeResponse())
			}))
			identity := newX509LifecycleIdentity(t, fixture)
			leaderScope := requestconfig.NewRequestRetryScope(0, time.Second, true, nil)
			leaderContext := requestconfig.WithRequestRetryScope(t.Context(), leaderScope)
			leaderResult := make(chan error, 1)
			go func() {
				_, err := identity.GetToken(leaderContext, fixture.capability)
				leaderResult <- err
			}()
			select {
			case <-reached:
			case <-time.After(5 * time.Second):
				t.Fatal("issuer did not send leader response headers")
			}
			followerScope := requestconfig.NewRequestRetryScope(test.retries, test.maximum, true, func(n int) { consumed.Store(int32(n)) })
			parent := requestconfig.WithRequestRetryScope(t.Context(), followerScope)
			// A renewed 600ms minimum after the 400ms body drain exceeds this deadline.
			parent, cancel := context.WithTimeout(parent, 850*time.Millisecond)
			defer cancel()
			if !followerScope.BeginAttempt() {
				t.Fatal("follower initial attempt was rejected")
			}
			waiting := make(chan struct{})
			followerContext := &x509ObservedDoneContext{Context: parent, observed: waiting}
			followerResult := make(chan error, 1)
			go func() {
				token, err := identity.GetToken(followerContext, fixture.capability)
				if err == nil && token != x509ExchangeSyntheticToken {
					err = errors.New("follower received an unexpected bearer")
				}
				followerResult <- err
			}()
			select {
			case <-waiting:
			case <-time.After(5 * time.Second):
				t.Fatal("follower did not join the shared refresh")
			}
			time.Sleep(test.bodyDelay)
			releaseOnce.Do(func() { close(release) })
			select {
			case err := <-leaderResult:
				if err == nil {
					t.Error("zero-budget leader returned no error, want original issuer failure")
				}
			case <-time.After(5 * time.Second):
				t.Fatal("leader did not finish after the issuer response")
			}
			select {
			case err := <-followerResult:
				if (err == nil) != test.wantSuccess {
					t.Errorf("GetToken follower error = %v, want success %t", err, test.wantSuccess)
				}
				if !test.wantSuccess {
					var status *x509ExchangeHTTPError
					if !errors.As(err, &status) || status.statusCode != http.StatusServiceUnavailable {
						t.Errorf("GetToken refused follower error = %v, want original issuer 503", err)
					}
				}
			case <-time.After(5 * time.Second):
				t.Fatal("follower did not finish after the issuer response")
			}
			wantExchanges, wantConsumed := int32(1), int32(0)
			if test.wantSuccess {
				wantExchanges, wantConsumed = 2, 1
			}
			if exchanges.Load() != wantExchanges || consumed.Load() != wantConsumed {
				t.Errorf("follower issuer attempts/retries consumed = %d/%d, want %d/%d",
					exchanges.Load(), consumed.Load(), wantExchanges, wantConsumed)
			}
		})
	}
}

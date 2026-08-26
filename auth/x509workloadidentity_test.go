package auth

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/openai/openai-go/v3/internal/requestconfig"
)

func TestNewX509WorkloadIdentityAuthValidatesConfiguration(t *testing.T) {
	fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, x509ValidExchangeResponse())
	}))
	valid := X509WorkloadIdentity{
		IdentityProviderID: "synthetic-identity-provider",
		ServiceAccountID:   "synthetic-service-account",
		Transport:          fixture.capability,
	}
	for _, test := range []struct {
		name   string
		change func(*X509WorkloadIdentity)
	}{
		{name: "missing identity provider", change: func(config *X509WorkloadIdentity) { config.IdentityProviderID = "" }},
		{name: "blank identity provider", change: func(config *X509WorkloadIdentity) { config.IdentityProviderID = "  " }},
		{name: "missing service account", change: func(config *X509WorkloadIdentity) { config.ServiceAccountID = "" }},
		{name: "blank service account", change: func(config *X509WorkloadIdentity) { config.ServiceAccountID = "\t" }},
		{name: "missing transport", change: func(config *X509WorkloadIdentity) { config.Transport = nil }},
		{name: "negative refresh buffer", change: func(config *X509WorkloadIdentity) { config.RefreshBuffer = -time.Second }},
		{name: "excessive refresh buffer", change: func(config *X509WorkloadIdentity) { config.RefreshBuffer = time.Hour }},
		{name: "negative token exchange timeout", change: func(config *X509WorkloadIdentity) {
			config.TokenExchangeTimeout = -time.Second
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			config := valid
			test.change(&config)
			identity, err := NewX509WorkloadIdentityAuth(config)
			if identity != nil || err == nil {
				t.Fatalf("invalid X.509 workload configuration returned identity=%v error=%v", identity, err)
			}
		})
	}
	identity, err := NewX509WorkloadIdentityAuth(X509WorkloadIdentity{
		IdentityProviderID: "  synthetic-identity-provider\t",
		ServiceAccountID:   "\nsynthetic-service-account ",
		Transport:          fixture.capability,
	})
	if err != nil {
		t.Fatalf("valid X.509 workload configuration: %v", err)
	}
	if identity.config.IdentityProviderID != "synthetic-identity-provider" ||
		identity.config.ServiceAccountID != "synthetic-service-account" {
		t.Errorf("normalized identity identifiers = %q/%q", identity.config.IdentityProviderID,
			identity.config.ServiceAccountID)
	}
	if identity.config.TokenExchangeTimeout != x509DefaultTokenExchangeTimeout {
		t.Errorf("default token exchange timeout = %s, want %s", identity.config.TokenExchangeTimeout,
			x509DefaultTokenExchangeTimeout)
	}
}

func TestX509WorkloadIdentityAuthExchangesOnlyOnItsAttestedTransport(t *testing.T) {
	var exchanges atomic.Int32
	fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		exchanges.Add(1)
		_, _ = io.WriteString(w, x509ValidExchangeResponse())
	}))
	identity, err := NewX509WorkloadIdentityAuth(X509WorkloadIdentity{
		IdentityProviderID: "synthetic-identity-provider",
		ServiceAccountID:   "synthetic-service-account",
		Transport:          fixture.capability,
	})
	if err != nil {
		t.Fatalf("construct X.509 workload identity: %v", err)
	}
	for range 2 {
		token, tokenErr := identity.GetToken(t.Context(), fixture.capability)
		if tokenErr != nil || token != x509ExchangeSyntheticToken {
			t.Fatalf("exchange on attested transport token=%q error=%v", token, tokenErr)
		}
	}
	if got := exchanges.Load(); got != 1 {
		t.Errorf("cached workload identity made %d exchanges, want one per transport generation", got)
	}
	for _, doer := range []HTTPDoer{nil, http.DefaultClient} {
		if token, tokenErr := identity.GetToken(t.Context(), doer); token != "" || tokenErr == nil {
			t.Errorf("unattested HTTP client returned token=%q error=%v", token, tokenErr)
		}
	}
	another, err := NewX509Transport(fixture.template)
	if err != nil {
		t.Fatalf("attest a separate certificate generation: %v", err)
	}
	t.Cleanup(func() { _ = another.Close() })
	if token, tokenErr := identity.GetToken(t.Context(), another); token != "" || tokenErr == nil {
		t.Errorf("different attested generation returned token=%q error=%v", token, tokenErr)
	}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	if token, tokenErr := identity.GetToken(ctx, fixture.capability); token != "" || !errors.Is(tokenErr, context.Canceled) {
		t.Errorf("canceled X.509 exchange returned token=%q error=%v", token, tokenErr)
	}
	var invalid *X509WorkloadIdentityAuth
	if token, tokenErr := invalid.GetToken(t.Context(), fixture.capability); token != "" ||
		tokenErr == nil || !strings.Contains(tokenErr.Error(), "invalid") {
		t.Errorf("nil X.509 identity returned token=%q error=%v", token, tokenErr)
	}
	if got := exchanges.Load(); got != 1 {
		t.Errorf("rejected transports/cancellation caused unexpected exchanges: %d", got)
	}
}

func TestX509WorkloadIdentityNeverReusesRejectedBearerWithoutRefreshInFlight(t *testing.T) {
	var exchanges atomic.Int32
	fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		exchanges.Add(1)
		_, _ = io.WriteString(w, x509ValidExchangeResponse())
	}))
	identity := newX509LifecycleIdentity(t, fixture)
	rejected, err := identity.GetToken(t.Context(), fixture.capability)
	if err != nil || rejected != x509ExchangeSyntheticToken {
		t.Fatalf("prime rejected X.509 bearer = %q, error = %v", rejected, err)
	}

	identity.invalidateToken(rejected)
	identity.mu.Lock()
	if identity.cached.value != "" || identity.rejectedToken != rejected || identity.inFlight != nil {
		identity.mu.Unlock()
		t.Fatal("ordinary 401 did not retain the rejected bearer generation")
	}
	identity.mu.Unlock()

	token, err := identity.GetToken(t.Context(), fixture.capability)
	if token != "" || !errors.Is(err, errX509InvalidatedBearer) {
		t.Errorf("issuer repeatedly returned rejected bearer: token=%q error=%v", token, err)
	}
	if got, want := exchanges.Load(), int32(1+x509MaximumAttempts); got != want {
		t.Errorf("prime and bounded replacement exchanges = %d, want %d", got, want)
	}
}

func TestX509WorkloadIdentityBoundsCompleteTokenExchange(t *testing.T) {
	reached := make(chan struct{})
	release := make(chan struct{})
	fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(_ http.ResponseWriter, request *http.Request) {
		close(reached)
		select {
		case <-request.Context().Done():
		case <-release:
		}
	}))
	identity, err := NewX509WorkloadIdentityAuth(X509WorkloadIdentity{
		IdentityProviderID:   "synthetic-identity-provider",
		ServiceAccountID:     "synthetic-service-account",
		TokenExchangeTimeout: 50 * time.Millisecond,
		Transport:            fixture.capability,
	})
	if err != nil {
		t.Fatalf("construct bounded X.509 workload identity: %v", err)
	}
	started := time.Now()
	token, err := identity.GetToken(context.Background(), fixture.capability)
	close(release)
	if token != "" || !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("bounded X.509 exchange token=%q error=%v", token, err)
	}
	if elapsed := time.Since(started); elapsed > time.Second {
		t.Errorf("bounded X.509 exchange took %s, want less than one second", elapsed)
	}
	select {
	case <-reached:
	default:
		t.Error("bounded X.509 exchange did not reach its issuer")
	}
}

func TestX509WorkloadIdentityNeverRepublishesBearerInvalidatedDuringRefresh(t *testing.T) {
	const replacement = "synthetic-replacement-bearer"
	var exchanges atomic.Int32
	refreshStarted := make(chan struct{})
	releaseRefresh := make(chan struct{})
	fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		switch exchanges.Add(1) {
		case 1:
			_, _ = io.WriteString(w, x509ValidExchangeResponse())
		case 2:
			close(refreshStarted)
			select {
			case <-releaseRefresh:
				_, _ = io.WriteString(w, x509ValidExchangeResponse())
			case <-request.Context().Done():
			}
		default:
			_, _ = io.WriteString(w, strings.Replace(x509ValidExchangeResponse(),
				x509ExchangeSyntheticToken, replacement, 1))
		}
	}))
	identity := newX509LifecycleIdentity(t, fixture)
	rejected, err := identity.GetToken(t.Context(), fixture.capability)
	if err != nil || rejected != x509ExchangeSyntheticToken {
		t.Fatalf("prime proactively refreshed bearer = %q, error = %v", rejected, err)
	}
	identity.mu.Lock()
	identity.refreshAfter = time.Now().Add(-time.Second)
	identity.mu.Unlock()

	type result struct {
		token string
		err   error
	}
	results := make(chan result, 2)
	refresh := func(ctx context.Context) {
		token, tokenErr := identity.GetToken(ctx, fixture.capability)
		results <- result{token: token, err: tokenErr}
	}
	go refresh(t.Context())
	select {
	case <-refreshStarted:
	case <-t.Context().Done():
		t.Fatal("proactive refresh never reached its synchronized issuer response")
	}
	follower := &x509ObservedDoneContext{Context: t.Context(), observed: make(chan struct{})}
	go refresh(follower)
	select {
	case <-follower.observed:
	case <-t.Context().Done():
		t.Fatal("concurrent refresh follower never waited for the in-flight issuer response")
	}
	identity.invalidateToken(rejected)
	identity.mu.Lock()
	if identity.cached.value != "" || identity.inFlight == nil || identity.inFlight.generation != rejected {
		identity.mu.Unlock()
		t.Fatal("rejected bearer was not invalidated while its proactive refresh remained in flight")
	}
	identity.mu.Unlock()
	close(releaseRefresh)

	for range 2 {
		select {
		case got := <-results:
			if got.token != replacement || got.err != nil {
				t.Errorf("invalidated proactive refresh did not acquire a distinct bearer: token=%q error=%v",
					got.token, got.err)
			}
		case <-t.Context().Done():
			t.Fatal("invalidated proactive refresh did not release its concurrent callers")
		}
	}
	identity.mu.Lock()
	if identity.cached.value != replacement || identity.rejectedToken != "" || identity.inFlight != nil {
		t.Error("distinct replacement bearer was not safely published and cleared its rejected generation")
	}
	identity.mu.Unlock()
	if got := exchanges.Load(); got != 3 {
		t.Errorf("prime/invalidated/distinct replacement exchanges = %d, want three", got)
	}
	for range 2 {
		token, tokenErr := identity.GetToken(t.Context(), fixture.capability)
		if tokenErr != nil || token != replacement {
			t.Errorf("fresh post-invalidation bearer = %q, error = %v", token, tokenErr)
		}
	}
	if got := exchanges.Load(); got != 3 {
		t.Errorf("fresh post-invalidation bearer required %d exchanges, want three", got)
	}
}

func TestX509WorkloadIdentityBoundsRepeatedInvalidatedBearerResponses(t *testing.T) {
	for _, maximum := range []int{-1, 0, 1} {
		name := "unscoped"
		if maximum >= 0 {
			name = "scoped retries " + strconv.Itoa(maximum)
		}
		t.Run(name, func(t *testing.T) {
			var exchanges atomic.Int32
			refreshStarted := make(chan struct{})
			releaseRefresh := make(chan struct{})
			fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
				if exchanges.Add(1) == 2 {
					close(refreshStarted)
					select {
					case <-releaseRefresh:
					case <-request.Context().Done():
						return
					}
				}
				_, _ = io.WriteString(w, x509ValidExchangeResponse())
			}))
			identity := newX509LifecycleIdentity(t, fixture)
			rejected, err := identity.GetToken(t.Context(), fixture.capability)
			if err != nil {
				t.Fatalf("prime repeatedly invalidated workload bearer: %v", err)
			}
			identity.mu.Lock()
			identity.refreshAfter = time.Now().Add(-time.Second)
			identity.mu.Unlock()
			ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
			defer cancel()
			if maximum >= 0 {
				scope := requestconfig.NewRequestRetryScope(maximum, time.Millisecond, true, nil)
				if !scope.BeginAttempt() {
					t.Fatal("begin invalidated workload request attempt")
				}
				ctx = requestconfig.WithRequestRetryScope(ctx, scope)
			}
			result := make(chan error, 1)
			go func() {
				token, exchangeErr := identity.GetToken(ctx, fixture.capability)
				if token != "" {
					result <- errors.New("repeated issuer response republished its rejected bearer")
					return
				}
				result <- exchangeErr
			}()
			select {
			case <-refreshStarted:
			case <-ctx.Done():
				t.Fatal("repeatedly invalidated refresh never reached its issuer")
			}
			identity.invalidateToken(rejected)
			close(releaseRefresh)
			select {
			case exchangeErr := <-result:
				if !errors.Is(exchangeErr, errX509InvalidatedBearer) ||
					strings.Contains(exchangeErr.Error(), rejected) {
					t.Errorf("repeated rejected bearer error = %v, want a safe bounded failure", exchangeErr)
				}
			case <-ctx.Done():
				t.Fatal("repeated rejected bearer exhausted its context instead of bounded retry budget")
			}
			want := int32(1 + x509MaximumAttempts)
			if maximum >= 0 {
				want = int32(2 + maximum)
			}
			if got := exchanges.Load(); got != want {
				t.Errorf("bounded repeated-bearer issuer exchanges = %d, want %d", got, want)
			}
			identity.mu.Lock()
			if identity.cached.value != "" || identity.rejectedToken != rejected ||
				!identity.refreshAfter.IsZero() || identity.inFlight != nil {
				t.Error("repeatedly rejected bearer was republished or lost its invalidation record")
			}
			identity.mu.Unlock()
		})
	}
}

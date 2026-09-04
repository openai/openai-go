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

	"github.com/openai/openai-go/v3/internal/requestconfig"
)

type x509ObservedErrorContext struct {
	context.Context
	errObserved chan struct{}
	once        sync.Once
}

func (ctx *x509ObservedErrorContext) Err() error {
	ctx.once.Do(func() { close(ctx.errObserved) })
	return ctx.Context.Err()
}

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
	if identity.config.IdentityProviderID != "  synthetic-identity-provider\t" ||
		identity.config.ServiceAccountID != "\nsynthetic-service-account " {
		t.Errorf("configured identity identifiers changed to %q/%q", identity.config.IdentityProviderID,
			identity.config.ServiceAccountID)
	}
	if identity.tokenExchangeTimeout != x509DefaultTokenExchangeTimeout {
		t.Errorf("default token exchange timeout = %s, want %s", identity.tokenExchangeTimeout,
			x509DefaultTokenExchangeTimeout)
	}
	if x509DefaultTokenExchangeTimeout != 30*time.Second {
		t.Errorf("X.509 token exchange default = %s, want 30s", x509DefaultTokenExchangeTimeout)
	}

	defaultCtx, cancelDefault := identity.exchangeContext(context.Background())
	defer cancelDefault()
	defaultDeadline, hasDefaultDeadline := defaultCtx.Deadline()
	if remaining := time.Until(defaultDeadline); !hasDefaultDeadline ||
		remaining <= x509DefaultTokenExchangeTimeout-time.Second || remaining > x509DefaultTokenExchangeTimeout {
		t.Errorf("default exchange deadline remaining = %s, present = %t", remaining, hasDefaultDeadline)
	}

	callerCtx, cancelCaller := context.WithTimeout(context.Background(), 2*x509DefaultTokenExchangeTimeout)
	defer cancelCaller()
	callerDeadline, _ := callerCtx.Deadline()
	exchangeCtx, cancelExchange := identity.exchangeContext(callerCtx)
	defer cancelExchange()
	exchangeDeadline, hasExchangeDeadline := exchangeCtx.Deadline()
	if !hasExchangeDeadline || !exchangeDeadline.Equal(callerDeadline) {
		t.Errorf("explicit caller deadline changed from %s to %s", callerDeadline, exchangeDeadline)
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
	if identity.cached.value != "" || !identity.bearers.isRejected(rejected, time.Now()) || identity.inFlight != nil {
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

func TestX509WorkloadIdentityRetainsLateRejectionAcrossRotations(t *testing.T) {
	const (
		first  = "synthetic-first-bearer"
		second = "synthetic-second-bearer"
		third  = "synthetic-third-bearer"
	)
	tokens := []string{first, second, first, third}
	var exchanges atomic.Int32
	fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		index := int(exchanges.Add(1)) - 1
		if index >= len(tokens) {
			http.Error(w, "unexpected exchange", http.StatusInternalServerError)
			return
		}
		_, _ = io.WriteString(w, strings.Replace(x509ValidExchangeResponse(),
			x509ExchangeSyntheticToken, tokens[index], 1))
	}))
	identity := newX509LifecycleIdentity(t, fixture)

	prime, err := identity.GetToken(t.Context(), fixture.capability)
	if err != nil || prime != first {
		t.Fatalf("prime first X.509 bearer = %q, error = %v", prime, err)
	}
	forceRefresh := func() string {
		identity.mu.Lock()
		identity.refreshAfter = time.Now().Add(-time.Second)
		identity.mu.Unlock()
		token, refreshErr := identity.GetToken(t.Context(), fixture.capability)
		if refreshErr != nil {
			t.Fatalf("force X.509 bearer refresh: %v", refreshErr)
		}
		return token
	}
	if token := forceRefresh(); token != second {
		t.Fatalf("rotated X.509 bearer = %q, want %q", token, second)
	}

	identity.invalidateToken(first)
	identity.mu.Lock()
	secondRemainsCached := identity.cached.value == second
	firstRemainsRejected := identity.bearers.isRejected(first, time.Now())
	identity.mu.Unlock()
	if !secondRemainsCached || !firstRemainsRejected {
		t.Fatal("late first-generation rejection evicted its replacement or lost its tombstone")
	}
	if token := forceRefresh(); token != third {
		t.Errorf("reissued rejected X.509 bearer was accepted: token=%q, want %q", token, third)
	}
	if got := exchanges.Load(); got != int32(len(tokens)) {
		t.Errorf("X.509 rotation exchanges = %d, want %d", got, len(tokens))
	}
}

func TestX509WorkloadIdentityBoundsCompleteTokenExchange(t *testing.T) {
	reached := make(chan struct{})
	release := make(chan struct{})
	var reachedOnce sync.Once
	var releaseOnce sync.Once
	releaseHandler := func() { releaseOnce.Do(func() { close(release) }) }
	defer releaseHandler()
	fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(_ http.ResponseWriter, request *http.Request) {
		reachedOnce.Do(func() { close(reached) })
		select {
		case <-request.Context().Done():
		case <-release:
		}
	}))
	identity, err := NewX509WorkloadIdentityAuth(X509WorkloadIdentity{
		IdentityProviderID: "synthetic-identity-provider",
		ServiceAccountID:   "synthetic-service-account",
		Transport:          fixture.capability,
	})
	if err != nil {
		t.Fatalf("construct bounded X.509 workload identity: %v", err)
	}
	identity.tokenExchangeTimeout = 250 * time.Millisecond
	result := make(chan struct {
		token string
		err   error
	}, 1)
	started := time.Now()
	go func() {
		token, tokenErr := identity.GetToken(context.Background(), fixture.capability)
		result <- struct {
			token string
			err   error
		}{token, tokenErr}
	}()
	var tokenResult struct {
		token string
		err   error
	}
	select {
	case tokenResult = <-result:
	case <-time.After(4 * time.Second):
		releaseHandler()
		t.Fatal("bounded X.509 exchange ignored both its internal timeout and caller watchdog")
	}
	releaseHandler()
	token, err := tokenResult.token, tokenResult.err
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

func TestX509WorkloadIdentityFallsBackOnlyAfterInternalExchangeTimeout(t *testing.T) {
	var exchanges atomic.Int32
	release := make(chan struct{})
	var releaseOnce sync.Once
	releaseHandler := func() { releaseOnce.Do(func() { close(release) }) }
	defer releaseHandler()
	fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if exchanges.Add(1) == 1 {
			_, _ = io.WriteString(w, x509ValidExchangeResponse())
			return
		}
		<-release
	}))
	identity, err := NewX509WorkloadIdentityAuth(X509WorkloadIdentity{
		IdentityProviderID: "synthetic-identity-provider",
		ServiceAccountID:   "synthetic-service-account",
		Transport:          fixture.capability,
	})
	if err != nil {
		t.Fatalf("construct bounded X.509 workload identity: %v", err)
	}
	identity.tokenExchangeTimeout = 250 * time.Millisecond
	initial, err := identity.GetToken(t.Context(), fixture.capability)
	if err != nil {
		t.Fatalf("prime proactive-refresh token cache: %v", err)
	}
	identity.mu.Lock()
	identity.refreshAfter = time.Now().Add(-time.Second)
	identity.mu.Unlock()

	type tokenResult struct {
		token string
		err   error
	}
	getToken := func(ctx context.Context) <-chan tokenResult {
		result := make(chan tokenResult, 1)
		go func() {
			token, tokenErr := identity.GetToken(ctx, fixture.capability)
			result <- tokenResult{token: token, err: tokenErr}
		}()
		return result
	}
	var refreshed tokenResult
	select {
	case refreshed = <-getToken(context.Background()):
	case <-time.After(4 * time.Second):
		releaseHandler()
		t.Fatal("proactive refresh ignored both its internal timeout and caller watchdog")
	}
	if refreshed.err != nil || refreshed.token != initial {
		t.Fatalf("internal exchange timeout fallback=%q, want %q; error=%v",
			refreshed.token, initial, refreshed.err)
	}

	identity.tokenExchangeTimeout = 2 * time.Second
	identity.mu.Lock()
	identity.refreshAfter = time.Now().Add(-time.Second)
	identity.mu.Unlock()
	callerCtx, cancel := context.WithTimeout(t.Context(), 250*time.Millisecond)
	defer cancel()
	var caller tokenResult
	select {
	case caller = <-getToken(callerCtx):
	case <-time.After(4 * time.Second):
		releaseHandler()
		t.Fatal("caller timeout did not release its token exchange")
	}
	if caller.token != "" || !errors.Is(caller.err, context.DeadlineExceeded) {
		t.Errorf("caller timeout returned token=%q error=%v", caller.token, caller.err)
	}
	if got := exchanges.Load(); got != 3 {
		t.Errorf("prime and timed exchange attempts = %d, want 3", got)
	}
}

func TestX509WorkloadIdentityCachedFallbackHonorsCallerCancellation(t *testing.T) {
	retryable := &x509ExchangeHTTPError{statusCode: http.StatusServiceUnavailable}
	active := context.Background()
	if !x509CanFallBackToCachedToken(retryable, active, active) {
		t.Error("active caller was denied fallback after a retryable exchange failure")
	}

	canceled, cancel := context.WithCancel(active)
	cancel()
	if x509CanFallBackToCachedToken(retryable, canceled, canceled) {
		t.Error("canceled caller was allowed fallback after a retryable exchange failure")
	}
	deadline, cancelDeadline := context.WithTimeout(active, 0)
	defer cancelDeadline()
	<-deadline.Done()
	if x509CanFallBackToCachedToken(retryable, deadline, deadline) {
		t.Error("expired caller deadline was allowed fallback after a retryable exchange failure")
	}

	internal, cancelInternal := context.WithTimeoutCause(active, 0, errX509TokenExchangeTimeout)
	defer cancelInternal()
	<-internal.Done()
	if !x509CanFallBackToCachedToken(context.DeadlineExceeded, internal, active) {
		t.Error("active caller was denied fallback after the internal exchange timeout")
	}
}

func TestX509WorkloadIdentityCachedFallbackRechecksCallerAfterLock(t *testing.T) {
	baseCaller, cancelCaller := context.WithCancel(context.Background())
	caller := &x509ObservedErrorContext{
		Context:     baseCaller,
		errObserved: make(chan struct{}),
	}
	exchange, cancelExchange := context.WithTimeoutCause(
		context.Background(), 0, errX509TokenExchangeTimeout,
	)
	defer cancelExchange()
	<-exchange.Done()

	identity := &X509WorkloadIdentityAuth{
		cached: x509ExchangedToken{
			value:     "synthetic-cached-bearer",
			expiresAt: time.Now().Add(time.Hour),
		},
	}
	identity.mu.Lock()
	result := make(chan struct {
		token string
		err   error
	}, 1)
	go func() {
		token, err := identity.tokenAfterExchangeContextDone(exchange, caller)
		result <- struct {
			token string
			err   error
		}{token: token, err: err}
	}()
	<-caller.errObserved
	cancelCaller()
	identity.mu.Unlock()

	got := <-result
	if got.token != "" || !errors.Is(got.err, context.Canceled) {
		t.Errorf("fallback after caller cancellation = token:%q error:%v", got.token, got.err)
	}
}

func TestX509WorkloadIdentityFollowerUsesCacheAfterInternalTimeout(t *testing.T) {
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
		t.Fatalf("construct concurrent X.509 workload identity: %v", err)
	}
	identity.tokenExchangeTimeout = 50 * time.Millisecond
	const cached = "synthetic-cached-bearer"
	refresh := &x509TokenRefresh{
		done:       make(chan struct{}),
		wake:       make(chan struct{}),
		generation: cached,
	}
	identity.mu.Lock()
	identity.cached = x509ExchangedToken{value: cached, expiresAt: time.Now().Add(time.Minute)}
	identity.refreshAfter = time.Now().Add(-time.Second)
	identity.inFlight = refresh
	identity.mu.Unlock()
	finishLeader := make(chan struct{})
	leaderFinished := make(chan struct{})
	var finishOnce sync.Once
	finish := func() { finishOnce.Do(func() { close(finishLeader) }) }
	go func() {
		<-finishLeader
		identity.mu.Lock()
		if identity.inFlight == refresh {
			identity.inFlight = nil
			refresh.err = context.DeadlineExceeded
			close(refresh.done)
		}
		identity.mu.Unlock()
		close(leaderFinished)
	}()
	defer func() {
		finish()
		<-leaderFinished
	}()

	followerCtx, followerCancel := context.WithCancel(context.Background())
	defer followerCancel()
	type followerTokenResult struct {
		token string
		err   error
	}
	followerResult := make(chan followerTokenResult, 1)
	go func() {
		token, followerErr := identity.GetToken(followerCtx, fixture.capability)
		followerResult <- followerTokenResult{token: token, err: followerErr}
	}()
	select {
	case result := <-followerResult:
		if result.err != nil || result.token != cached {
			t.Errorf("internally timed-out refresh follower token=%q, want %q; error=%v",
				result.token, cached, result.err)
		}
	case <-time.After(3 * time.Second):
		followerCancel()
		t.Fatal("internally timed-out refresh follower did not return")
	}
	if got := exchanges.Load(); got != 0 {
		t.Errorf("timed-out refresh follower made %d redundant exchanges, want 0", got)
	}
	finish()
	select {
	case <-leaderFinished:
	case <-time.After(3 * time.Second):
		t.Fatal("synthetic refresh leader did not publish completion")
	}
	identity.mu.Lock()
	inFlight := identity.inFlight
	identity.mu.Unlock()
	if inFlight != nil {
		t.Error("synthetic refresh leader did not clear the in-flight refresh")
	}
	select {
	case <-refresh.done:
	default:
		t.Error("synthetic refresh leader did not release refresh followers")
	}
}

func TestX509WorkloadIdentityInvalidationKeepsWakeChannelStable(t *testing.T) {
	const bearer = "synthetic-bearer"
	wake := make(chan struct{})
	refresh := &x509TokenRefresh{wake: wake, generation: bearer}
	identity := &X509WorkloadIdentityAuth{
		cached:   x509ExchangedToken{value: bearer, expiresAt: time.Now().Add(time.Minute)},
		inFlight: refresh,
	}

	identity.invalidateToken(bearer)
	if refresh.wake != wake {
		t.Fatal("bearer invalidation mutated the refresh wake channel")
	}
	select {
	case <-wake:
	default:
		t.Fatal("bearer invalidation did not wake its refresh generation")
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
	identity.mu.Lock()
	follower := &x509ObservedDoneContext{Context: t.Context(), observed: make(chan struct{})}
	go refresh(follower)
	select {
	case <-follower.observed:
	case <-time.After(3 * time.Second):
		identity.mu.Unlock()
		t.Fatal("concurrent refresh follower never reached the synchronized identity state")
	}
	identity.invalidateTokenLocked(rejected)
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
	if identity.cached.value != replacement || !identity.bearers.isRejected(rejected, time.Now()) ||
		identity.inFlight != nil {
		t.Error("distinct replacement bearer was not safely published while retaining its rejected predecessor")
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
			if identity.cached.value != "" || !identity.bearers.isRejected(rejected, time.Now()) ||
				!identity.refreshAfter.IsZero() || identity.inFlight != nil {
				t.Error("repeatedly rejected bearer was republished or lost its invalidation record")
			}
			identity.mu.Unlock()
		})
	}
}

func TestX509WorkloadIdentityStopsExcessiveIssuerHints(t *testing.T) {
	for _, statusCode := range []int{429, 503} {
		for _, test := range []struct {
			name, header, value string
			maximum             time.Duration
		}{
			{name: "default cap", header: "Retry-After", value: "9"},
			{name: "caller cap", header: "Retry-After-Ms", value: "10", maximum: time.Millisecond},
			{name: "overflow", header: "Retry-After", value: "1e999", maximum: time.Millisecond},
			{name: "preferred overflow", header: "Retry-After-Ms", value: "1e999", maximum: time.Millisecond},
			{name: "date", header: "Retry-After", value: time.Now().Add(time.Hour).UTC().Format(http.TimeFormat)},
		} {
			t.Run(strconv.Itoa(statusCode)+"/"+test.name, func(t *testing.T) {
				var attempts atomic.Int32
				fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					attempts.Add(1)
					w.Header().Set("Retry-After", "0")
					w.Header().Set(test.header, test.value)
					w.WriteHeader(statusCode)
					_, _ = io.WriteString(w, "private-issuer-body")
				}))
				identity := newX509LifecycleIdentity(t, fixture)
				ctx, cancel := context.WithTimeout(t.Context(), 250*time.Millisecond)
				defer cancel()
				if test.maximum != 0 {
					ctx = requestconfig.WithRequestRetryScope(ctx, requestconfig.NewRequestRetryScope(2, test.maximum, true, nil))
				}
				token, err := identity.GetToken(ctx, fixture.capability)
				var status *x509ExchangeHTTPError
				if token != "" || !errors.As(err, &status) || status.statusCode != statusCode {
					t.Errorf("GetToken(%s) = (%q, %v), want original HTTP %d error", test.name, token, err, statusCode)
				}
				if got := attempts.Load(); got != 1 {
					t.Errorf("GetToken(%s) issuer attempts = %d, want 1", test.name, got)
				}
				if err != nil && strings.Contains(err.Error(), "private-") {
					t.Errorf("GetToken(%s) exposed issuer body", test.name)
				}
			})
		}
	}
}

func TestX509WorkloadIdentityPreservesIssuerWaitAfterTruncatedBodyAndWake(t *testing.T) {
	var attempts atomic.Int32
	var first time.Time
	fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if attempts.Add(1) == 1 {
			first = time.Now()
			w.Header().Set("Retry-After-Ms", "100")
			w.Header().Set("Content-Length", "1000")
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = io.WriteString(w, "private-truncated-issuer-body")
			return
		}
		if elapsed := time.Since(first); elapsed < 100*time.Millisecond {
			t.Errorf("issuer replay after truncated body = %s, want at least 100ms", elapsed)
		}
		_, _ = io.WriteString(w, x509ValidExchangeResponse())
	}))
	identity := newX509LifecycleIdentity(t, fixture)
	wake := make(chan struct{})
	close(wake)
	token, err := identity.exchangeWithRetry(t.Context(), fixture.capability, &x509TokenRefresh{wake: wake})
	if err != nil || token.value != x509ExchangeSyntheticToken || attempts.Load() != 2 {
		t.Errorf("exchangeWithRetry(truncated body, wake) attempts=%d error=%v, want successful second attempt", attempts.Load(), err)
	}
}

func TestX509WorkloadIdentityFollowerDoesNotReplayIssuerHint(t *testing.T) {
	var attempts atomic.Int32
	fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		attempts.Add(1)
		_, _ = io.WriteString(w, x509ValidExchangeResponse())
	}))
	identity := newX509LifecycleIdentity(t, fixture)
	issuerError := &x509ExchangeHTTPError{statusCode: 503, hasRetryAfter: true, retryAfter: time.Second}
	done := make(chan struct{})
	close(done)
	for _, ownerError := range []error{nil, context.Canceled} {
		identity.inFlight = &x509TokenRefresh{done: done, err: errors.Join(issuerError, ownerError), ownerContextErr: ownerError}
		var retries int
		ctx := requestconfig.WithRequestRetryScope(t.Context(), requestconfig.NewRequestRetryScope(1, time.Second, true, func(count int) { retries = count }))
		_, err := identity.GetToken(ctx, fixture.capability)
		if !errors.Is(err, issuerError) || errors.Is(err, context.Canceled) || attempts.Load() != 0 || retries != 0 {
			t.Errorf("GetToken(follower with issuer hint, owner=%v) = %v, attempts=%d retries=%d, want original error and zero attempts/retries", ownerError, err, attempts.Load(), retries)
		}
	}
}

func TestX509WorkloadIdentityRetainsIssuerHintWhenWaitIsCanceled(t *testing.T) {
	for _, statusCode := range []int{429, 503} {
		t.Run(strconv.Itoa(statusCode), func(t *testing.T) {
			var attempts atomic.Int32
			fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				attempts.Add(1)
				w.Header().Set("Retry-After", "1")
				w.WriteHeader(statusCode)
			}))
			identity := newX509LifecycleIdentity(t, fixture)
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			scope := requestconfig.NewRequestRetryScope(1, time.Second, true, func(int) { cancel() })
			ctx = requestconfig.WithRequestRetryScope(ctx, scope)
			_, err := identity.exchangeWithRetry(ctx, fixture.capability, &x509TokenRefresh{})
			var status *x509ExchangeHTTPError
			if !errors.Is(err, context.Canceled) || !errors.As(err, &status) || status.statusCode != statusCode || status.retryAfter != time.Second || attempts.Load() != 1 {
				t.Fatalf("canceled wait error=%v attempts=%d, want cancellation and retained issuer minimum after one request", err, attempts.Load())
			}
		})
	}
}

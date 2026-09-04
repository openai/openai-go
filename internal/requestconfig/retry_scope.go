package requestconfig

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"
)

type requestRetryScopeContextKey struct{}

type requestRetryScopeFactory struct {
	maximum         int
	maximumDelay    time.Duration
	allowBodyReplay bool
}

// RequestRetryScope shares one retry budget across provider authentication,
// ordinary SDK request attempts, and a single unauthorized-response replay.
// This is internal API and may change without notice.
type RequestRetryScope struct {
	mu                  sync.Mutex
	maximum             int
	outerRetries        int
	internalRetries     int
	attempts            int
	replayUsed          bool
	allowBodyReplay     bool
	onInternalRetry     func(int)
	authenticationRetry *AuthenticationRetryState
}

// AuthenticationRetryState preserves only an issuer minimum for one logical request.
// Background refreshes may retain it without retaining request payloads or retry callbacks.
// This is internal API and may change without notice.
type AuthenticationRetryState struct {
	mu           sync.Mutex
	maximumDelay time.Duration
	refusal      error
	until        time.Time
}

// NewRequestRetryScope creates the retry scope for one logical SDK request.
func NewRequestRetryScope(maximum int, maximumDelay time.Duration, allowBodyReplay bool, onInternalRetry func(int)) *RequestRetryScope {
	if maximumDelay <= 0 {
		maximumDelay = DefaultMaxServerDelay
	}
	return &RequestRetryScope{
		maximum:             max(0, maximum),
		allowBodyReplay:     allowBodyReplay,
		onInternalRetry:     onInternalRetry,
		authenticationRetry: &AuthenticationRetryState{maximumDelay: maximumDelay},
	}
}

// InstallRequestRetryScope attaches a fresh logical-request budget to cfg and
// ensures configuration clones receive independent budgets and callbacks.
func (cfg *RequestConfig) InstallRequestRetryScope(allowBodyReplay bool) *RequestRetryScope {
	factory := &requestRetryScopeFactory{
		maximum:         cfg.MaxRetries,
		maximumDelay:    cfg.MaxRetryDelay,
		allowBodyReplay: allowBodyReplay,
	}
	cfg.authentication.retryScopeFactory = factory
	return factory.install(cfg)
}

// InstallRequestAttemptMiddleware counts complete SDK attempts before caller
// middleware can short-circuit authentication. Configuration clones inherit
// this middleware and resolve their own independently installed scope.
func (cfg *RequestConfig) InstallRequestAttemptMiddleware() {
	cfg.Middlewares = append([]middleware{func(request *http.Request, next middlewareNext) (*http.Response, error) {
		if request == nil {
			return nil, WithNoRetryError(errors.New("X.509 workload identity requires a non-nil request"))
		}
		scope := RequestRetryScopeFromContext(request.Context())
		if scope == nil || !scope.BeginAttempt() {
			return nil, WithNoRetryError(errors.New("X.509 workload identity retry budget exhausted"))
		}
		return next(request)
	}}, cfg.Middlewares...)
}

func (factory *requestRetryScopeFactory) install(cfg *RequestConfig) *RequestRetryScope {
	cfg.MaxRetries = factory.maximum
	scope := NewRequestRetryScope(factory.maximum, factory.maximumDelay, factory.allowBodyReplay, func(consumed int) {
		cfg.MaxRetries = factory.maximum - consumed
	})
	cfg.Request = cfg.Request.WithContext(WithRequestRetryScope(cfg.Request.Context(), scope))
	return scope
}

// WithRequestRetryScope associates a scope with an existing request context.
func WithRequestRetryScope(ctx context.Context, scope *RequestRetryScope) context.Context {
	return context.WithValue(ctx, requestRetryScopeContextKey{}, scope)
}

// RequestRetryScopeFromContext returns the logical request's retry scope.
func RequestRetryScopeFromContext(ctx context.Context) *RequestRetryScope {
	scope, _ := ctx.Value(requestRetryScopeContextKey{}).(*RequestRetryScope)
	return scope
}

// AuthenticationRetryState returns the request's lightweight issuer-refusal storage.
func (scope *RequestRetryScope) AuthenticationRetryState() *AuthenticationRetryState {
	return scope.authenticationRetry
}

// RefuseAuthenticationRetry preserves an issuer refusal for this logical request.
// A zero notBefore preserves the refusal until this logical request ends.
func (scope *RequestRetryScope) RefuseAuthenticationRetry(err error, notBefore time.Time) {
	scope.authenticationRetry.RefuseAuthenticationRetry(err, notBefore)
}

// AuthenticationRetryRefusal returns the original issuer refusal until its minimum elapses.
func (scope *RequestRetryScope) AuthenticationRetryRefusal() error {
	return scope.authenticationRetry.AuthenticationRetryRefusal()
}

// RefuseAuthenticationRetry records the first issuer refusal and its minimum.
func (state *AuthenticationRetryState) RefuseAuthenticationRetry(err error, notBefore time.Time) {
	state.mu.Lock()
	defer state.mu.Unlock()
	if state.refusal == nil {
		state.refusal = err
		state.until = notBefore
	}
}

// AuthenticationRetryRefusal returns the original refusal until its minimum elapses.
func (state *AuthenticationRetryState) AuthenticationRetryRefusal() error {
	state.mu.Lock()
	defer state.mu.Unlock()
	if !state.until.IsZero() && !time.Now().Before(state.until) {
		state.refusal = nil
		state.until = time.Time{}
	}
	return state.refusal
}

// MaxRetryDelay reports the logical request's configured server-delay cap.
func (state *AuthenticationRetryState) MaxRetryDelay() time.Duration {
	return state.maximumDelay
}

// BeginAttempt records the first dispatch or one ordinary SDK retry attempt.
func (scope *RequestRetryScope) BeginAttempt() bool {
	scope.mu.Lock()
	defer scope.mu.Unlock()
	if scope.attempts != 0 {
		if scope.outerRetries+scope.internalRetries >= scope.maximum {
			return false
		}
		scope.outerRetries++
	}
	scope.attempts++
	return true
}

// TryRetry consumes one retry for an internal provider-authentication attempt.
func (scope *RequestRetryScope) TryRetry() bool {
	scope.mu.Lock()
	defer scope.mu.Unlock()
	return scope.tryInternalRetry()
}

// TryOuterReplay reserves the sole unauthorized recovery without consuming a
// retry. The SDK's next complete attempt accounts for the retry instead.
func (scope *RequestRetryScope) TryOuterReplay() bool {
	scope.mu.Lock()
	defer scope.mu.Unlock()
	if scope.replayUsed || scope.outerRetries+scope.internalRetries >= scope.maximum {
		return false
	}
	scope.replayUsed = true
	return true
}

// AllowBodyReplay reports whether caller middleware could have transformed the
// body after it became replayable.
func (scope *RequestRetryScope) AllowBodyReplay() bool {
	return scope.allowBodyReplay
}

// MaxRetryDelay bounds issuer-directed and SDK-selected retry waits for this
// logical request.
func (scope *RequestRetryScope) MaxRetryDelay() time.Duration {
	return scope.authenticationRetry.MaxRetryDelay()
}

func (scope *RequestRetryScope) tryInternalRetry() bool {
	if scope.outerRetries+scope.internalRetries >= scope.maximum {
		return false
	}
	scope.internalRetries++
	if scope.onInternalRetry != nil {
		scope.onInternalRetry(scope.internalRetries)
	}
	return true
}

// AuthenticationRetryDelay shares ordinary response hint parsing with auth replay.
// A missing or invalid hint permits immediate auth recovery; an excessive valid
// hint refuses the replay instead of shortening the server's minimum.
func AuthenticationRetryDelay(response *http.Response, maximumDelay time.Duration) (delay time.Duration, hasHint, allowed bool) {
	delay, hasHint, exceeds := parseRetryAfterHeader(response, maximumDelay)
	return delay, hasHint, !exceeds
}

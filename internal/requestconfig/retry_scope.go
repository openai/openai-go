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
	mu              sync.Mutex
	maximum         int
	maximumDelay    time.Duration
	outerRetries    int
	internalRetries int
	attempts        int
	replayUsed      bool
	allowBodyReplay bool
	onInternalRetry func(int)
}

// NewRequestRetryScope creates the retry scope for one logical SDK request.
func NewRequestRetryScope(maximum int, maximumDelay time.Duration, allowBodyReplay bool, onInternalRetry func(int)) *RequestRetryScope {
	if maximumDelay <= 0 {
		maximumDelay = DefaultMaxServerDelay
	}
	return &RequestRetryScope{
		maximum:         max(0, maximum),
		maximumDelay:    maximumDelay,
		allowBodyReplay: allowBodyReplay,
		onInternalRetry: onInternalRetry,
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
	return scope.maximumDelay
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

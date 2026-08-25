package requestconfig

import (
	"context"
	"sync"
)

type requestRetryScopeContextKey struct{}

// RequestRetryScope shares one retry budget across provider authentication,
// ordinary SDK request attempts, and a single unauthorized-response replay.
// This is internal API and may change without notice.
type RequestRetryScope struct {
	mu              sync.Mutex
	maximum         int
	outerRetries    int
	internalRetries int
	attempts        int
	replayUsed      bool
	allowBodyReplay bool
	onInternalRetry func(int)
}

// NewRequestRetryScope creates the retry scope for one logical SDK request.
func NewRequestRetryScope(maximum int, allowBodyReplay bool, onInternalRetry func(int)) *RequestRetryScope {
	return &RequestRetryScope{
		maximum:         max(0, maximum),
		allowBodyReplay: allowBodyReplay,
		onInternalRetry: onInternalRetry,
	}
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

// TryReplay consumes one retry for the only allowed unauthorized-response
// replay associated with this logical request.
func (scope *RequestRetryScope) TryReplay() bool {
	scope.mu.Lock()
	defer scope.mu.Unlock()
	if scope.replayUsed || !scope.tryInternalRetry() {
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

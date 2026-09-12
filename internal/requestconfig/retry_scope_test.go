package requestconfig

import (
	"net/http"
	"sync"
	"testing"
	"time"
)

func TestRequestRetryScopeSharesLogicalAttemptBudget(t *testing.T) {
	var consumed []int
	scope := NewRequestRetryScope(2, time.Second, true, func(value int) { consumed = append(consumed, value) })
	for range 2 {
		if !scope.BeginAttempt() {
			t.Fatal("initial request and first outer retry were not admitted")
		}
	}
	if !scope.TryOuterReplay() {
		t.Fatal("remaining retry budget did not admit the first 401 replay")
	}
	if scope.TryOuterReplay() || !scope.BeginAttempt() || scope.TryRetry() || scope.BeginAttempt() {
		t.Error("combined retry scope admitted an attempt beyond its logical budget")
	}
	if len(consumed) != 0 {
		t.Errorf("outer retry incorrectly consumed internal retry callbacks: %v", consumed)
	}
}

func TestRequestRetryScopePreservesContextAndReplayPolicy(t *testing.T) {
	scope := NewRequestRetryScope(2, 0, false, nil)
	ctx := WithRequestRetryScope(t.Context(), scope)
	if RequestRetryScopeFromContext(ctx) != scope || scope.AllowBodyReplay() {
		t.Error("request context lost its retry scope or transformed-body replay policy")
	}
	if RequestRetryScopeFromContext(t.Context()) != nil {
		t.Error("ordinary contexts unexpectedly acquired an X.509 retry scope")
	}
	if scope.MaxRetryDelay() != DefaultMaxServerDelay {
		t.Errorf("default maximum retry delay = %s, want %s", scope.MaxRetryDelay(), DefaultMaxServerDelay)
	}
	if !scope.BeginAttempt() {
		t.Fatal("initial issuer attempt was rejected")
	}
	for range 2 {
		if !scope.TryRetry() {
			t.Fatal("issuer retry budget was rejected prematurely")
		}
	}
	if scope.TryRetry() {
		t.Error("issuer retries did not stop at the logical request budget")
	}
}

func TestRequestRetryScopeReservesOnlyOneCompleteUnauthorizedAttempt(t *testing.T) {
	var consumed []int
	scope := NewRequestRetryScope(2, time.Second, true, func(value int) { consumed = append(consumed, value) })
	if !scope.BeginAttempt() || !scope.TryOuterReplay() {
		t.Fatal("initial attempt did not reserve its unauthorized outer retry")
	}
	if scope.TryOuterReplay() {
		t.Error("one logical request reserved multiple unauthorized recoveries")
	}
	if len(consumed) != 0 {
		t.Errorf("outer recovery incorrectly consumed an internal retry: %v", consumed)
	}
	if !scope.BeginAttempt() || !scope.TryRetry() {
		t.Error("reserved outer attempt did not preserve the remaining issuer retry")
	}
	if scope.BeginAttempt() || scope.TryRetry() {
		t.Error("outer recovery exceeded its shared logical retry budget")
	}
	if len(consumed) != 1 || consumed[0] != 1 {
		t.Errorf("mixed outer/internal retry callbacks = %v, want [1]", consumed)
	}
	zero := NewRequestRetryScope(0, time.Second, true, nil)
	if !zero.BeginAttempt() || zero.TryOuterReplay() {
		t.Error("zero-retry scope incorrectly reserved an unauthorized recovery")
	}
}

func TestRequestAttemptMiddlewarePreservesIndependentCloneScopes(t *testing.T) {
	request, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "https://example.com", nil)
	if err != nil {
		t.Fatalf("construct counted request: %v", err)
	}
	parent := RequestConfig{Request: request, MaxRetries: 1}
	parentScope := parent.InstallRequestRetryScope(true)
	parent.InstallRequestAttemptMiddleware()
	if len(parent.Middlewares) != 1 {
		t.Fatalf("parent installed %d attempt counters, want exactly one", len(parent.Middlewares))
	}
	var calls int
	next := func(*http.Request) (*http.Response, error) {
		calls++
		return &http.Response{StatusCode: http.StatusOK}, nil
	}
	attempt := func(cfg *RequestConfig) error {
		response, attemptErr := cfg.Middlewares[0](cfg.Request, next)
		if response != nil && response.Body != nil {
			if closeErr := response.Body.Close(); closeErr != nil {
				t.Errorf("close synthetic counted-attempt response: %v", closeErr)
			}
		}
		return attemptErr
	}
	for range 2 {
		if attemptErr := attempt(&parent); attemptErr != nil {
			t.Fatalf("counted parent attempt failed: %v", attemptErr)
		}
	}
	if attemptErr := attempt(&parent); attemptErr == nil {
		t.Error("attempt counter accepted a caller-middleware short circuit beyond its budget")
	}
	clone := parent.Clone(t.Context())
	if clone == nil || len(clone.Middlewares) != 1 {
		t.Fatalf("cloned request duplicated its outer attempt counter: %v", clone)
	}
	cloneScope := RequestRetryScopeFromContext(clone.Request.Context())
	if cloneScope == nil || cloneScope == parentScope {
		t.Fatal("cloned outer attempt counter reused its parent's scope")
	}
	if attemptErr := attempt(clone); attemptErr != nil {
		t.Fatalf("cloned attempt counter did not resolve its fresh context scope: %v", attemptErr)
	}
	if calls != 3 {
		t.Errorf("parent and clone dispatched %d admitted attempts, want 3", calls)
	}
}

func TestRequestRetryScopePreservesConfiguredMaximumDelay(t *testing.T) {
	scope := NewRequestRetryScope(2, 7*time.Millisecond, true, nil)
	if got := scope.MaxRetryDelay(); got != 7*time.Millisecond {
		t.Errorf("request maximum retry delay = %s, want 7ms", got)
	}
}

func TestRequestRetryScopeClonesHaveIndependentBudgetsAndCallbacks(t *testing.T) {
	request, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "https://example.com", nil)
	if err != nil {
		t.Fatalf("construct retry-scope clone request: %v", err)
	}
	parent := RequestConfig{Request: request, MaxRetries: 2, MaxRetryDelay: 7 * time.Millisecond}
	parentScope := parent.InstallRequestRetryScope(true)
	if !parentScope.BeginAttempt() || !parentScope.TryRetry() || parent.MaxRetries != 1 {
		t.Fatal("parent retry scope did not consume its own logical budget")
	}
	clone := parent.Clone(t.Context())
	if clone == nil || clone.MaxRetries != 2 {
		t.Fatalf("cloned logical request maximum = %v, want reset original maximum 2", clone)
	}
	cloneScope := RequestRetryScopeFromContext(clone.Request.Context())
	if cloneScope == nil || cloneScope == parentScope || cloneScope.MaxRetryDelay() != 7*time.Millisecond {
		t.Fatal("request clone reused its parent's scope or lost its configured maximum delay")
	}
	if !cloneScope.BeginAttempt() || !cloneScope.TryRetry() || clone.MaxRetries != 1 {
		t.Fatal("cloned logical request did not consume its own independent retry")
	}
	if parent.MaxRetries != 1 {
		t.Errorf("cloned retry mutated ancestor MaxRetries = %d", parent.MaxRetries)
	}
	var finished sync.WaitGroup
	const concurrentClones = 12
	finished.Add(concurrentClones)
	for range concurrentClones {
		go func() {
			defer finished.Done()
			current := parent.Clone(t.Context())
			if current == nil {
				t.Error("concurrent request clone returned nil")
				return
			}
			scope := RequestRetryScopeFromContext(current.Request.Context())
			if scope == nil || !scope.BeginAttempt() || !scope.TryRetry() || current.MaxRetries != 1 {
				t.Error("concurrent request clone lost its isolated retry budget")
			}
		}()
	}
	finished.Wait()
	if parent.MaxRetries != 1 {
		t.Errorf("concurrent request clones mutated ancestor MaxRetries = %d", parent.MaxRetries)
	}
}

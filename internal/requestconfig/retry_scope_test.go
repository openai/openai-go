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
	if !scope.TryReplay() {
		t.Fatal("remaining retry budget did not admit the first 401 replay")
	}
	if scope.TryReplay() || scope.TryRetry() || scope.BeginAttempt() {
		t.Error("combined retry scope admitted an attempt beyond its logical budget")
	}
	if len(consumed) != 1 || consumed[0] != 1 {
		t.Errorf("internal retry callbacks = %v, want [1]", consumed)
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

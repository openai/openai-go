package requestconfig

import "testing"

func TestRequestRetryScopeSharesLogicalAttemptBudget(t *testing.T) {
	var consumed []int
	scope := NewRequestRetryScope(2, true, func(value int) { consumed = append(consumed, value) })
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
	scope := NewRequestRetryScope(2, false, nil)
	ctx := WithRequestRetryScope(t.Context(), scope)
	if RequestRetryScopeFromContext(ctx) != scope || scope.AllowBodyReplay() {
		t.Error("request context lost its retry scope or transformed-body replay policy")
	}
	if RequestRetryScopeFromContext(t.Context()) != nil {
		t.Error("ordinary contexts unexpectedly acquired an X.509 retry scope")
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

package auth

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/openai/openai-go/v4/internal/requestconfig"
)

type cacheTestProvider struct {
	started chan struct{}
	release chan struct{}
}

func (*cacheTestProvider) TokenType() SubjectTokenType { return SubjectTokenTypeJWT }

func (p *cacheTestProvider) GetToken(ctx context.Context, _ HTTPDoer) (string, error) {
	if p.started != nil {
		close(p.started)
		select {
		case <-p.release:
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}
	return "synthetic-subject", nil
}

type cacheTestTransport string

func (transport cacheTestTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Body != nil {
		defer func() { _ = req.Body.Close() }()
	}
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body: io.NopCloser(strings.NewReader(fmt.Sprintf(
			`{"access_token":"synthetic-%s","expires_in":3600}`, transport))),
	}, nil
}

func TestWorkloadIdentityBackgroundRefreshKeepsSelectedTransport(t *testing.T) {
	provider := &cacheTestProvider{started: make(chan struct{}), release: make(chan struct{})}
	identity, err := NewWorkloadIdentityAuth(WorkloadIdentity{
		IdentityProviderID: "synthetic-provider", ServiceAccountID: "synthetic-account",
		Provider: provider, RefreshBufferSeconds: 60,
	})
	if err != nil {
		t.Fatal(err)
	}
	client := &http.Client{Transport: cacheTestTransport("a")}
	cache, _, release := identity.acquireCache(client)
	cache.cachedToken = "synthetic-prime"
	cache.tokenExpiry = time.Now().Add(30 * time.Second)
	release()
	if token, tokenErr := identity.GetToken(t.Context(), client); tokenErr != nil || token != "synthetic-prime" {
		t.Fatalf("proactive refresh did not preserve the live cached token: %v", tokenErr)
	}
	cache.mu.Lock()
	refresh := cache.refreshInFlight
	cache.mu.Unlock()
	if refresh == nil {
		t.Fatal("proactive refresh was not started")
	}
	defer func() {
		close(provider.release)
		select {
		case <-refresh.done:
			if refresh.result.err != nil || refresh.result.token != "synthetic-a" {
				t.Errorf("background refresh changed identity: %v", refresh.result.err)
			}
		case <-time.After(5 * time.Second):
			t.Error("proactive refresh did not finish")
		}
	}()
	select {
	case <-provider.started:
	case <-time.After(5 * time.Second):
		t.Fatal("proactive provider was not called")
	}
	// The initiating SDK call has returned, but its background refresh still
	// owns transport A. A sequential client replacement must not retarget it.
	client.Transport = cacheTestTransport("b")
}

// Retention and eviction are internal resource contracts; public tests cover
// the credentials observed by each effective transport.
func TestWorkloadIdentityCacheRetentionAndActiveOwnership(t *testing.T) {
	identity := &WorkloadIdentityAuth{}
	activeClient := &http.Client{}
	active, _, releaseActive := identity.acquireCache(activeClient)
	defer releaseActive()
	for range workloadMaximumCachedIdentities * 3 {
		_, _, release := identity.acquireCache(&http.Client{})
		release()
	}
	if len(identity.partitions.entries) != workloadMaximumCachedIdentities {
		t.Fatal("identity cache did not retain its bounded capacity")
	}
	again, _, releaseAgain := identity.acquireCache(activeClient)
	defer releaseAgain()
	if again != active {
		t.Fatal("active request lost its authentication state during cache churn")
	}
}

func TestWorkloadIdentityCacheDoesNotEvictBackgroundRefresh(t *testing.T) {
	identity := &WorkloadIdentityAuth{}
	client := &http.Client{}
	cache, _, release := identity.acquireCache(client)
	cache.refreshInFlight = &tokenRefreshState{done: make(chan struct{})}
	release()
	for range workloadMaximumCachedIdentities * 2 {
		_, _, releaseOther := identity.acquireCache(&http.Client{})
		releaseOther()
	}
	again, _, releaseAgain := identity.acquireCache(client)
	defer releaseAgain()
	if again != cache {
		t.Fatal("background refresh was evicted while still in flight")
	}
}

func TestWorkloadIdentityCacheOverflowUsesIndependentState(t *testing.T) {
	identity := &WorkloadIdentityAuth{}
	for range workloadMaximumCachedIdentities {
		_, _, release := identity.acquireCache(&http.Client{})
		defer release()
	}
	client := &http.Client{}
	first, _, releaseFirst := identity.acquireCache(client)
	defer releaseFirst()
	second, _, releaseSecond := identity.acquireCache(client)
	defer releaseSecond()
	if first == second || len(identity.partitions.entries) != workloadMaximumCachedIdentities {
		t.Fatal("overflow shared credentials or exceeded retained identity capacity")
	}
}

func TestWorkloadIdentityCacheInvalidationDoesNotChangeOtherPartitions(t *testing.T) {
	identity := &WorkloadIdentityAuth{}
	a, _, releaseA := identity.acquireCache(&http.Client{})
	defer releaseA()
	b, _, releaseB := identity.acquireCache(&http.Client{})
	defer releaseB()
	for _, cache := range []*workloadTokenCache{a, b} {
		cache.cachedToken = "synthetic-equal-bearer"
		cache.tokenExpiry = time.Now().Add(time.Hour)
		cache.beginRefreshLocked()
	}
	a.invalidateToken("synthetic-equal-bearer")
	if a.cachedToken != "" || a.rejectedToken == "" || b.cachedToken == "" || b.rejectedToken != "" {
		t.Fatal("in-flight token rejection crossed identity partitions")
	}
}

type cacheFunctionTransport func(*http.Request) (*http.Response, error)

func (transport cacheFunctionTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return transport(req)
}

func TestWorkloadIdentityNativeFunctionTransportReusesToken(t *testing.T) {
	identity, err := NewWorkloadIdentityAuth(WorkloadIdentity{
		IdentityProviderID: "synthetic-provider", ServiceAccountID: "synthetic-account",
		Provider: &cacheTestProvider{},
	})
	if err != nil {
		t.Fatal(err)
	}
	exchanges := 0
	client := &http.Client{Transport: cacheFunctionTransport(func(req *http.Request) (*http.Response, error) {
		exchanges++
		return cacheTestTransport("native-function").RoundTrip(req)
	})}
	for range 2 {
		if _, err := identity.GetToken(t.Context(), client); err != nil {
			t.Fatal(err)
		}
	}
	if exchanges != 1 {
		t.Fatalf("unchanged native client exchanged %d tokens, want 1", exchanges)
	}
}

func TestWorkloadIdentityOverflowPreservesCacheAcrossOuterReplay(t *testing.T) {
	identity, err := NewWorkloadIdentityAuth(WorkloadIdentity{
		IdentityProviderID: "synthetic-provider", ServiceAccountID: "synthetic-account",
		Provider: &cacheTestProvider{},
	})
	if err != nil {
		t.Fatal(err)
	}
	for range workloadMaximumCachedIdentities {
		_, _, release := identity.acquireCache(&http.Client{})
		defer release()
	}
	exchanges := 0
	client := &http.Client{Transport: cacheFunctionTransport(func(req *http.Request) (*http.Response, error) {
		exchanges++
		token := "rejected"
		if exchanges >= 2 {
			token = "fresh"
		}
		return cacheTestTransport(token).RoundTrip(req)
	})}
	// Use a comparable transport wrapper so this exercises saturated retention,
	// independently of the native function-transport case above.
	client.Transport = &cacheRoundTripperWrapper{client.Transport}
	scope := requestconfig.NewRequestRetryScope(3, time.Second, true, nil)
	ctx := requestconfig.WithRequestRetryScope(t.Context(), scope)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://example.invalid/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	next := func(req *http.Request) (*http.Response, error) {
		status := http.StatusUnauthorized
		if req.Header.Get("Authorization") == "Bearer synthetic-fresh" {
			status = http.StatusOK
		}
		return &http.Response{StatusCode: status, Header: make(http.Header), Body: http.NoBody}, nil
	}
	first, err := WorkloadIdentityMiddleware(identity, client, req, next)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = first.Body.Close() }()
	if first.StatusCode != http.StatusUnauthorized {
		t.Fatalf("first status = %d", first.StatusCode)
	}
	second, err := WorkloadIdentityMiddleware(identity, client, req, next)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = second.Body.Close() }()
	if second.StatusCode != http.StatusOK || exchanges != 2 {
		t.Fatalf("replay status=%d exchanges=%d, want 200 and 2", second.StatusCode, exchanges)
	}
	// A subsequent ordinary attempt in the same logical request must retain
	// the replacement cache instead of creating another overflow partition.
	third, err := WorkloadIdentityMiddleware(identity, client, req, next)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = third.Body.Close() }()
	if third.StatusCode != http.StatusOK || exchanges != 2 {
		t.Fatalf("retained replay cache: response=%v exchanges=%d error=%v", third, exchanges, err)
	}
}

type cacheRoundTripperWrapper struct{ http.RoundTripper }

func TestWorkloadIdentityNativeFunctionTransportReplacementSeparatesClosures(t *testing.T) {
	identity, err := NewWorkloadIdentityAuth(WorkloadIdentity{
		IdentityProviderID: "synthetic-provider", ServiceAccountID: "synthetic-account", Provider: &cacheTestProvider{},
	})
	if err != nil {
		t.Fatal(err)
	}
	calls := 0
	makeTransport := func(name string) http.RoundTripper {
		return cacheFunctionTransport(func(req *http.Request) (*http.Response, error) {
			calls++
			return cacheTestTransport(name).RoundTrip(req)
		})
	}
	a, b := makeTransport("a"), makeTransport("b")
	client := &http.Client{}
	for _, step := range []struct {
		transport http.RoundTripper
		token     string
	}{{a, "synthetic-a"}, {a, "synthetic-a"}, {b, "synthetic-b"}, {b, "synthetic-b"}, {a, "synthetic-a"}} {
		client.Transport = step.transport
		token, err := identity.GetToken(t.Context(), client)
		if err != nil || token != step.token {
			t.Fatalf("wrong transport token: %q %v", token, err)
		}
	}
	if calls != 2 {
		t.Fatalf("transport exchanges=%d, want 2", calls)
	}
}

type cacheSliceTransport struct {
	labels []string
	calls  *int
}

func (transport cacheSliceTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	*transport.calls++
	return cacheTestTransport(transport.labels[0]).RoundTrip(req)
}

func TestWorkloadIdentityNativeSliceTransportIsolation(t *testing.T) {
	identity, err := NewWorkloadIdentityAuth(WorkloadIdentity{
		IdentityProviderID: "synthetic-provider", ServiceAccountID: "synthetic-account", Provider: &cacheTestProvider{},
	})
	if err != nil {
		t.Fatal(err)
	}
	calls := 0
	a := http.RoundTripper(cacheSliceTransport{labels: []string{"a"}, calls: &calls})
	b := http.RoundTripper(cacheSliceTransport{labels: []string{"b"}, calls: &calls})
	first, second := &http.Client{Transport: a}, &http.Client{Transport: a}
	for _, client := range []*http.Client{first, first, second, second} {
		if token, err := identity.GetToken(t.Context(), client); err != nil || token != "synthetic-a" {
			t.Fatalf("native slice transport: token=%q error=%v", token, err)
		}
	}
	if calls != 2 {
		t.Fatalf("distinct clients shared cache: calls=%d", calls)
	}
	first.Transport = b
	for range 2 {
		if token, err := identity.GetToken(t.Context(), first); err != nil || token != "synthetic-b" {
			t.Fatalf("replacement slice transport: token=%q error=%v", token, err)
		}
	}
	if calls != 3 {
		t.Fatalf("replacement transport calls=%d, want 3", calls)
	}
}

// WorkloadIdentityEvictionWindowForTest coordinates the internal eviction window
// while an external-package test drives authentication through the public client.
func WorkloadIdentityEvictionWindowForTest(t *testing.T, identity *WorkloadIdentityAuth, client HTTPDoer) func(func()) {
	t.Helper()
	cache, _, release := identity.acquireCache(client)
	release()
	for range workloadMaximumCachedIdentities - 1 {
		_, _, release := identity.acquireCache(&http.Client{})
		t.Cleanup(release)
	}
	return func(failRefresh func()) {
		t.Helper()
		cache.mu.Lock()
		refresh := cache.refreshInFlight
		cache.mu.Unlock()
		if refresh == nil {
			t.Fatal("expected proactive refresh before unauthorized recovery")
			return
		}
		failRefresh()
		select {
		case <-refresh.done:
			if refresh.result.err == nil {
				t.Fatal("expected failed background refresh")
			}
		case <-time.After(5 * time.Second):
			t.Fatal("background refresh did not complete")
		}
		_, _, release := identity.acquireCache(&http.Client{})
		release()
		identity.mu.Lock()
		defer identity.mu.Unlock()
		for _, partition := range identity.partitions.entries {
			if partition.cache == cache {
				t.Fatal("retained cache was not evicted before outer replay")
			}
		}
	}
}

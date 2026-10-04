package auth

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
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

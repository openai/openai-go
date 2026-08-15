package option

import (
	"net/http"
	"testing"

	"github.com/openai/openai-go/v3/auth"
)

func TestX509WorkloadIdentityAuthCacheIsBoundedAndLeastRecentlyUsed(t *testing.T) {
	config := auth.X509WorkloadIdentity{
		IdentityProviderID: "idp-test",
		ServiceAccountID:   "svc-test",
	}
	cache := x509WorkloadIdentityAuthCache{}
	httpClients := make([]*http.Client, x509WorkloadIdentityAuthCacheCapacity+1)
	authStates := make([]*auth.WorkloadIdentityAuth, x509WorkloadIdentityAuthCacheCapacity)
	for i := range authStates {
		httpClients[i] = &http.Client{}
		var err error
		authStates[i], err = cache.get(httpClients[i], config)
		if err != nil {
			t.Fatalf("cache.get(%d) error = %v", i, err)
		}
	}

	// Refresh the first entry so the second becomes least recently used.
	first, err := cache.get(httpClients[0], config)
	if err != nil {
		t.Fatalf("cache.get(first) error = %v", err)
	}
	if first != authStates[0] {
		t.Fatal("cache.get(first) returned a new auth state")
	}

	httpClients[len(httpClients)-1] = &http.Client{}
	_, err = cache.get(httpClients[len(httpClients)-1], config)
	if err != nil {
		t.Fatalf("cache.get(new) error = %v", err)
	}
	if got, want := len(cache.entries), x509WorkloadIdentityAuthCacheCapacity; got != want {
		t.Fatalf("cache entry count = %d, want %d", got, want)
	}

	first, err = cache.get(httpClients[0], config)
	if err != nil {
		t.Fatalf("cache.get(first after eviction) error = %v", err)
	}
	if first != authStates[0] {
		t.Fatal("recently used auth state was evicted")
	}
	second, err := cache.get(httpClients[1], config)
	if err != nil {
		t.Fatalf("cache.get(second after eviction) error = %v", err)
	}
	if second == authStates[1] {
		t.Fatal("least recently used auth state was retained")
	}
	if got, want := len(cache.entries), x509WorkloadIdentityAuthCacheCapacity; got != want {
		t.Fatalf("cache entry count after replacement = %d, want %d", got, want)
	}
}

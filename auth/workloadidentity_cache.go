package auth

import (
	"net/http"
	"reflect"
)

// This limits retained transport identities, not the number of clients callers
// may use. When every slot is busy, an independent cache serves the caller.
const workloadMaximumCachedIdentities = 32

type workloadIdentityKey struct {
	doer      HTTPDoer
	transport http.RoundTripper
}

// Keep the exported authenticator comparable, as in the existing API.
type workloadIdentityPartitions struct {
	entries []*workloadIdentityPartition
}

type workloadIdentityPartition struct {
	key   workloadIdentityKey
	cache *workloadTokenCache
	users int
}

func workloadHTTPIdentity(doer HTTPDoer) (workloadIdentityKey, HTTPDoer, bool) {
	if doer == nil {
		doer = http.DefaultClient
	}
	key := workloadIdentityKey{doer: doer}
	if client, ok := doer.(*http.Client); ok && client != nil {
		key.transport = client.Transport
		if key.transport == nil {
			key.transport = http.DefaultTransport
		}
		// Freeze the selected transport for provider calls and background refresh.
		// A later sequential client.Transport change must not change an in-flight
		// exchange's identity after its cache partition has been selected.
		snapshot := *client
		snapshot.Transport = key.transport
		doer = &snapshot
	}
	// Value.Comparable also checks dynamic values nested inside interface fields.
	// Type.Comparable alone can still allow a panic during map lookup/equality.
	return key, doer, reflect.ValueOf(key).Comparable()
}

func (w *WorkloadIdentityAuth) acquireCache(doer HTTPDoer) (*workloadTokenCache, HTTPDoer, func()) {
	key, selected, comparable := workloadHTTPIdentity(doer)
	if !comparable {
		return &workloadTokenCache{config: w.config}, selected, func() {}
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.partitions == nil {
		w.partitions = &workloadIdentityPartitions{}
	}
	for index, partition := range w.partitions.entries {
		if partition.key == key {
			// Keep most recently used identities at the end of the bounded slice.
			copy(w.partitions.entries[index:], w.partitions.entries[index+1:])
			w.partitions.entries[len(w.partitions.entries)-1] = partition
			partition.users++
			return partition.cache, selected, func() { w.releaseCache(partition) }
		}
	}
	cache := &workloadTokenCache{config: w.config}
	if len(w.partitions.entries) == workloadMaximumCachedIdentities {
		evict := -1
		for index, partition := range w.partitions.entries {
			partition.cache.mu.Lock()
			idle := partition.users == 0 && partition.cache.refreshInFlight == nil
			partition.cache.mu.Unlock()
			if idle {
				evict = index
				break
			}
		}
		if evict < 0 {
			return cache, selected, func() {}
		}
		copy(w.partitions.entries[evict:], w.partitions.entries[evict+1:])
		w.partitions.entries = w.partitions.entries[:len(w.partitions.entries)-1]
	}
	partition := &workloadIdentityPartition{key: key, cache: cache, users: 1}
	w.partitions.entries = append(w.partitions.entries, partition)
	return cache, selected, func() { w.releaseCache(partition) }
}

func (w *WorkloadIdentityAuth) releaseCache(partition *workloadIdentityPartition) {
	w.mu.Lock()
	defer w.mu.Unlock()
	partition.users--
}

package auth_test

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestWorkloadIdentityRefreshIsIndependentAcrossClients(t *testing.T) {
	identity := newOrdinaryWorkloadIdentity(t)
	started, release := make(chan struct{}), make(chan struct{})
	var aExchanges atomic.Int32
	a := &http.Client{Transport: &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
		if aExchanges.Add(1) == 1 {
			close(started)
		}
		select {
		case <-release:
			return ordinaryWorkloadResponse(http.StatusOK, `{"access_token":"synthetic-a","expires_in":3600}`), nil
		case <-req.Context().Done():
			return nil, req.Context().Err()
		}
	}}}
	b := ordinaryWorkloadIssuer(t, func() string { return "synthetic-b" })
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	var wg sync.WaitGroup
	defer wg.Wait()
	defer close(release)
	wg.Go(func() {
		token, err := identity.GetToken(ctx, a)
		if err != nil || token != "synthetic-a" {
			t.Errorf("identity A acquisition failed: %v", err)
		}
	})
	select {
	case <-started:
	case <-ctx.Done():
		t.Fatal("identity A did not start its exchange")
	}
	// B must finish while A's exchange is deliberately held open.
	if token, err := identity.GetToken(ctx, b); err != nil || token != "synthetic-b" {
		t.Fatalf("identity B did not refresh independently: %v", err)
	}
	followerCtx, cancelFollower := context.WithCancel(ctx)
	cancelFollower()
	if _, err := identity.GetToken(followerCtx, a); !errors.Is(err, context.Canceled) {
		t.Fatalf("canceled A follower error=%v", err)
	}
	for range 16 {
		wg.Go(func() {
			token, err := identity.GetToken(ctx, a)
			if err != nil || token != "synthetic-a" {
				t.Errorf("identity A follower acquisition failed: %v", err)
			}
		})
	}
	// Deferred release and wait ensure all followers finish before this check.
	t.Cleanup(func() {
		if got := aExchanges.Load(); got != 1 {
			t.Errorf("same-identity exchanges=%d, want 1", got)
		}
	})
}

func TestWorkloadIdentityDirectCallsKeepNativeClientsSeparate(t *testing.T) {
	identity := newOrdinaryWorkloadIdentity(t)
	var calls atomic.Int32
	transport := &closureTransport{fn: func(*http.Request) (*http.Response, error) {
		calls.Add(1)
		return ordinaryWorkloadResponse(http.StatusOK, `{"access_token":"synthetic-shared-transport","expires_in":3600}`), nil
	}}
	a, b := &http.Client{Transport: transport}, &http.Client{Transport: transport}
	for _, client := range []*http.Client{a, b, a, b} {
		if _, err := identity.GetToken(t.Context(), client); err != nil {
			t.Fatal(err)
		}
	}
	if calls.Load() != 2 {
		t.Errorf("two clients sharing a transport made %d exchanges, want 2", calls.Load())
	}
}

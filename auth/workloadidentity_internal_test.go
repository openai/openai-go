package auth

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestWorkloadIdentityNeverPublishesBearerInvalidatedDuringRefreshCompletion(t *testing.T) {
	const (
		previous    = "synthetic-previous-bearer"
		replacement = "synthetic-replacement-bearer"
	)
	identity := &WorkloadIdentityAuth{
		cachedToken: replacement,
		tokenExpiry: time.Now().Add(time.Hour),
	}
	refresh := &tokenRefreshState{
		done:       make(chan struct{}),
		generation: previous,
	}
	identity.refreshInFlight = refresh

	identity.invalidateToken(replacement)
	identity.finishRefresh(refresh, replacement, nil)

	if refresh.result.token != "" || !errors.Is(refresh.result.err, errInvalidatedWorkloadBearer) {
		t.Errorf("invalidated refresh result = token:%q error:%v", refresh.result.token, refresh.result.err)
	}
	if identity.cachedToken != "" || !identity.tokenExpiry.IsZero() ||
		!identity.bearers.isRejected(replacement, time.Now()) || identity.refreshInFlight != nil {
		t.Error("invalidated bearer was republished or its rejected generation was lost")
	}
}

func TestWorkloadIdentityRecordsLateRejectionWithoutEvictingReplacement(t *testing.T) {
	const (
		previous    = "synthetic-previous-bearer"
		replacement = "synthetic-replacement-bearer"
	)
	now := time.Now()
	expiresAt := now.Add(time.Hour)
	identity := &WorkloadIdentityAuth{
		cachedToken: replacement,
		tokenExpiry: expiresAt,
	}
	history := ensureWorkloadBearerHistory(&identity.bearers)
	if err := history.recordIssued(previous, expiresAt, now); err != nil {
		t.Fatalf("record previous bearer: %v", err)
	}
	if err := history.recordIssued(replacement, expiresAt, now); err != nil {
		t.Fatalf("record replacement bearer: %v", err)
	}

	identity.invalidateToken(previous)
	identity.mu.Lock()
	replacementRemainsCached := identity.cachedToken == replacement
	previousIsRejected := identity.bearers.isRejected(previous, time.Now())
	refresh := identity.beginRefreshLocked()
	identity.mu.Unlock()
	if !replacementRemainsCached || !previousIsRejected {
		t.Fatal("late rejection evicted a newer bearer or failed to retain the rejected generation")
	}

	identity.finishRefresh(refresh, previous, nil)
	if refresh.result.token != "" || !errors.Is(refresh.result.err, errInvalidatedWorkloadBearer) {
		t.Errorf("late rejected refresh result = token:%q error:%v", refresh.result.token, refresh.result.err)
	}
	if identity.cachedToken != replacement || identity.refreshInFlight != nil {
		t.Error("late rejected refresh disturbed the cached replacement generation")
	}
}

func TestWorkloadBearerHistoryExpiresTombstoneWithOriginalGeneration(t *testing.T) {
	const bearer = "synthetic-reusable-bearer"
	now := time.Unix(1_700_000_000, 0)
	originalExpiry := now.Add(time.Minute)
	history := &workloadBearerHistory{}
	if err := history.recordIssued(bearer, originalExpiry, now); err != nil {
		t.Fatalf("record original bearer: %v", err)
	}
	if !history.reject(bearer, time.Time{}, now) || !history.isRejected(bearer, now) {
		t.Fatal("issued bearer was not rejected through its original expiry")
	}
	afterExpiry := originalExpiry.Add(time.Nanosecond)
	if history.isRejected(bearer, afterExpiry) {
		t.Fatal("expired bearer tombstone prevented a later independent generation")
	}
	if err := history.recordIssued(bearer, afterExpiry.Add(time.Hour), afterExpiry); err != nil {
		t.Fatalf("record independent bearer: %v", err)
	}
	if history.isRejected(bearer, afterExpiry) {
		t.Fatal("new bearer generation inherited an expired rejection tombstone")
	}
}

func TestWorkloadBearerHistoryFailsClosedAtCapacityAndReleasesExpiredStorage(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	expiresAt := now.Add(time.Minute)
	history := &workloadBearerHistory{}
	for index := range workloadMaximumTrackedBearerGenerations {
		if err := history.recordIssued(fmt.Sprintf("synthetic-bearer-%d", index), expiresAt, now); err != nil {
			t.Fatalf("record bearer %d within capacity: %v", index, err)
		}
	}
	if err := history.recordIssued("synthetic-overflow-bearer", expiresAt, now); !errors.Is(err, errWorkloadBearerHistoryCapacity) {
		t.Fatalf("overflow history error = %v, want capacity error", err)
	}
	if got := len(history.generations); got != workloadMaximumTrackedBearerGenerations {
		t.Fatalf("tracked bearer generations = %d, want %d", got, workloadMaximumTrackedBearerGenerations)
	}

	afterExpiry := expiresAt.Add(time.Nanosecond)
	if history.isRejected("synthetic-untracked-bearer", afterExpiry) {
		t.Fatal("untracked bearer was rejected after capacity expiry")
	}
	if history.generations != nil || !history.nextExpiry.IsZero() {
		t.Fatal("expired bearer history retained its peak map allocation")
	}
	if err := history.recordIssued("synthetic-post-expiry-bearer", afterExpiry.Add(time.Hour), afterExpiry); err != nil {
		t.Fatalf("record bearer after capacity expired: %v", err)
	}
}

package auth

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

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

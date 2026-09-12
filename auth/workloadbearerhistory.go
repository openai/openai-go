package auth

import (
	"crypto/sha256"
	"errors"
	"time"
)

const workloadMaximumTrackedBearerGenerations = 1024

var errWorkloadBearerHistoryCapacity = errors.New("workload identity bearer history capacity exceeded")

// workloadBearerHistory retains issued bearer expirations so a late API
// rejection can tombstone the exact generation without evicting a newer one.
// Only fixed-size digests are retained; expired generations are pruned when
// the earliest known expiry is reached.
type workloadBearerHistory struct {
	generations map[[sha256.Size]byte]workloadBearerGeneration
	nextExpiry  time.Time
}

type workloadBearerGeneration struct {
	expiresAt time.Time
	rejected  bool
}

func ensureWorkloadBearerHistory(history **workloadBearerHistory) *workloadBearerHistory {
	if *history == nil {
		*history = &workloadBearerHistory{}
	}
	return *history
}

func (history *workloadBearerHistory) recordIssued(value string, expiresAt, now time.Time) error {
	if history == nil {
		return errWorkloadBearerHistoryCapacity
	}
	history.pruneExpired(now)
	if value == "" || !now.Before(expiresAt) {
		return nil
	}
	digest := sha256.Sum256([]byte(value))
	generation, tracked := history.generations[digest]
	if generation.rejected || !generation.expiresAt.Before(expiresAt) {
		return nil
	}
	if !tracked && len(history.generations) >= workloadMaximumTrackedBearerGenerations {
		return errWorkloadBearerHistoryCapacity
	}
	if history.generations == nil {
		history.generations = make(map[[sha256.Size]byte]workloadBearerGeneration)
	}
	generation.expiresAt = expiresAt
	history.generations[digest] = generation
	history.noteExpiry(expiresAt)
	return nil
}

func (history *workloadBearerHistory) reject(value string, knownExpiry, now time.Time) bool {
	if history == nil {
		return false
	}
	history.pruneExpired(now)
	if value == "" {
		return false
	}
	digest := sha256.Sum256([]byte(value))
	generation := history.generations[digest]
	if knownExpiry.Before(generation.expiresAt) {
		knownExpiry = generation.expiresAt
	}
	if !now.Before(knownExpiry) {
		return false
	}
	if history.generations == nil {
		history.generations = make(map[[sha256.Size]byte]workloadBearerGeneration)
	}
	history.generations[digest] = workloadBearerGeneration{expiresAt: knownExpiry, rejected: true}
	history.noteExpiry(knownExpiry)
	return true
}

func (history *workloadBearerHistory) isRejected(value string, now time.Time) bool {
	if history == nil || value == "" {
		return false
	}
	history.pruneExpired(now)
	generation, ok := history.generations[sha256.Sum256([]byte(value))]
	return ok && generation.rejected && now.Before(generation.expiresAt)
}

func (history *workloadBearerHistory) noteExpiry(expiresAt time.Time) {
	if history.nextExpiry.IsZero() || expiresAt.Before(history.nextExpiry) {
		history.nextExpiry = expiresAt
	}
}

func (history *workloadBearerHistory) pruneExpired(now time.Time) {
	if history.nextExpiry.IsZero() || now.Before(history.nextExpiry) {
		return
	}
	history.nextExpiry = time.Time{}
	previousSize := len(history.generations)
	for digest, generation := range history.generations {
		if !now.Before(generation.expiresAt) {
			delete(history.generations, digest)
			continue
		}
		history.noteExpiry(generation.expiresAt)
	}
	remaining := len(history.generations)
	if remaining == 0 {
		history.generations = nil
	} else if remaining*2 < previousSize {
		compacted := make(map[[sha256.Size]byte]workloadBearerGeneration, remaining)
		for digest, generation := range history.generations {
			compacted[digest] = generation
		}
		history.generations = compacted
	}
}

package auth

import (
	"errors"
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
		identity.rejectedToken != replacement || identity.refreshInFlight != nil {
		t.Error("invalidated bearer was republished or its rejected generation was lost")
	}
}

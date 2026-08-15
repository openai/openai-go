package auth

import (
	"net/http"
	"testing"
	"time"
)

func TestTokenExchangeRetryDelayHonorsRetryAfter(t *testing.T) {
	t.Run("seconds", func(t *testing.T) {
		resp := &http.Response{Header: http.Header{"Retry-After": []string{"2"}}}
		if got, want := tokenExchangeRetryDelay(resp, 0), 2*time.Second; got != want {
			t.Errorf("tokenExchangeRetryDelay() = %v, want %v", got, want)
		}
	})

	t.Run("capped", func(t *testing.T) {
		for _, retryAfter := range []string{"120", "1e300"} {
			resp := &http.Response{Header: http.Header{"Retry-After": []string{retryAfter}}}
			if got, want := tokenExchangeRetryDelay(resp, 0), tokenExchangeMaxRetryDelay; got != want {
				t.Errorf("tokenExchangeRetryDelay(%q) = %v, want %v", retryAfter, got, want)
			}
		}
	})

	t.Run("http date", func(t *testing.T) {
		retryAt := time.Now().Add(5 * time.Second).UTC().Truncate(time.Second)
		resp := &http.Response{Header: http.Header{"Retry-After": []string{retryAt.Format(http.TimeFormat)}}}
		got := tokenExchangeRetryDelay(resp, 0)
		if got < 4*time.Second || got > 5*time.Second {
			t.Errorf("tokenExchangeRetryDelay() = %v, want between 4s and 5s", got)
		}
	})
}

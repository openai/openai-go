package auth

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/openai/openai-go/v3/internal/requestconfig"
)

func TestX509OAuthIssuerRetryAfter(t *testing.T) {
	for _, status := range []int{400, 401, 403} {
		for _, test := range []struct {
			name    string
			code    string
			header  string
			value   string
			maximum time.Duration
			minimum time.Duration
			refuse  bool
		}{
			{name: "allowed seconds", code: "temporarily_unavailable", header: "Retry-After", value: "0.0755", minimum: 75500 * time.Microsecond},
			{name: "allowed milliseconds", code: "server_error", header: "Retry-After-Ms", value: "75.5", minimum: 75500 * time.Microsecond},
			{name: "default ceiling", code: "temporarily_unavailable", header: "Retry-After", value: "90", refuse: true},
			{name: "caller ceiling", code: "server_error", header: "Retry-After-Ms", value: "75.5", maximum: 50 * time.Millisecond, refuse: true},
			{name: "overflow", code: "server_error", header: "Retry-After", value: "1e999", refuse: true},
			{name: "malformed fallback", code: "temporarily_unavailable", header: "Retry-After", value: "synthetic-invalid", minimum: x509InitialRetryDelay},
			{name: "permanent error", code: "invalid_grant", header: "Retry-After", value: "90", refuse: true},
		} {
			t.Run(strconv.Itoa(status)+"/"+test.name, func(t *testing.T) {
				var attempts atomic.Int32
				var first atomic.Int64
				delays := make(chan time.Duration, 1)
				fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					if attempts.Add(1) == 1 {
						first.Store(time.Now().UnixNano())
						w.Header().Set(test.header, test.value)
						w.WriteHeader(status)
						_, _ = io.WriteString(w, `{"error":"`+test.code+`","error_description":"private-synthetic-body"}`)
						return
					}
					delays <- time.Since(time.Unix(0, first.Load()))
					_, _ = io.WriteString(w, x509ValidExchangeResponse())
				}))
				identity := newX509LifecycleIdentity(t, fixture)
				ctx := t.Context()
				if test.maximum > 0 {
					ctx = requestconfig.WithRequestRetryScope(ctx, requestconfig.NewRequestRetryScope(2, test.maximum, true, nil))
				}
				token, err := identity.GetToken(ctx, fixture.capability)
				if test.refuse {
					var oauth *OAuthError
					if token != "" || !errors.As(err, &oauth) || oauth.StatusCode != status || string(oauth.ErrorCode) != test.code {
						t.Fatalf("refusal lost original OAuth error: token=%q err=%v", token, err)
					}
					if strings.Contains(err.Error(), "private-synthetic-body") {
						t.Fatal("issuer error exposed response body")
					}
					if attempts.Load() != 1 {
						t.Fatalf("refused minimum made %d issuer attempts", attempts.Load())
					}
					return
				}
				if err != nil || token != x509ExchangeSyntheticToken || attempts.Load() != 2 {
					t.Fatalf("allowed retry: token=%q err=%v attempts=%d", token, err, attempts.Load())
				}
				if delay := <-delays; delay < test.minimum {
					t.Errorf("issuer retry after %s, before minimum %s", delay, test.minimum)
				}
			})
		}
	}
}

func TestX509IssuerHonorsRaisedRetryCap(t *testing.T) {
	var attempts atomic.Int32
	var first time.Time
	fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if attempts.Add(1) == 1 {
			first = time.Now()
			w.Header().Set("Retry-After", "10")
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		if elapsed := time.Since(first); elapsed < 10*time.Second {
			t.Errorf("issuer retry delay = %s, want at least 10s", elapsed)
		}
		_, _ = io.WriteString(w, x509ValidExchangeResponse())
	}))
	identity := newX509LifecycleIdentity(t, fixture)
	ctx, cancel := context.WithTimeout(t.Context(), 20*time.Second)
	defer cancel()
	ctx = requestconfig.WithRequestRetryScope(ctx, requestconfig.NewRequestRetryScope(1, 30*time.Second, true, nil))
	token, err := identity.GetToken(ctx, fixture.capability)
	if err != nil || token != x509ExchangeSyntheticToken || attempts.Load() != 2 {
		t.Errorf("GetToken(10s issuer minimum, 30s cap) = %q, %v, attempts=%d; want success after two attempts", token, err, attempts.Load())
	}
}

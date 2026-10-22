package websockettest

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	wire "github.com/coder/websocket"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/auth"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

func TestAuthenticatedWebSocketOpeningAttemptTimeout(t *testing.T) {
	const budget = 300 * time.Millisecond
	for _, tc := range []struct {
		name                                 string
		perAttempt, caller, native, httpOnly bool
		wantSuccess                          bool
	}{
		{"ws overall control", false, false, false, false, true},
		{"ws per attempt", true, false, false, false, true},
		{"http per attempt control", true, false, false, true, true},
		{"ws caller total bound", false, true, false, false, false},
		{"ws native total bound", false, false, true, false, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var attempts, exchanges atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				attempt := attempts.Add(1)
				// Each ordinary request completes inside its attempt budget,
				// while two requests take longer than a total opening budget.
				timer := time.NewTimer(200 * time.Millisecond)
				defer timer.Stop()
				select {
				case <-r.Context().Done():
					return
				case <-timer.C:
				}
				if r.Header.Get("Authorization") != "Bearer synthetic-token-2" {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				if attempt != 2 || r.Header.Get("X-Stainless-Retry-Count") != "1" {
					t.Errorf("second request count=%d, retry header=%q", attempt, r.Header.Get("X-Stainless-Retry-Count"))
				}
				if tc.httpOnly {
					w.Header().Set("Content-Type", "application/json")
					_, _ = io.WriteString(w, "{}")
					return
				}
				socket, err := wire.Accept(w, r, nil)
				if err != nil {
					t.Errorf("accept: %v", err)
					return
				}
				defer func() { _ = socket.CloseNow() }()
				_, _, _ = socket.Read(r.Context())
			}))
			defer server.Close()
			transport := server.Client().Transport
			httpClient := &http.Client{Transport: websocketAuthTransport(func(r *http.Request) (*http.Response, error) {
				if r.URL.String() == auth.TokenExchangeURL {
					token := exchanges.Add(1)
					return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(fmt.Sprintf("{\"access_token\":\"synthetic-token-%d\",\"expires_in\":3600}", token)))}, nil
				}
				return transport.RoundTrip(r)
			})}
			if tc.native {
				httpClient.Timeout = budget
			}
			opts := []option.RequestOption{option.WithUnsafeAllowHTTP(), option.WithBaseURL(server.URL),
				option.WithHTTPClient(httpClient), option.WithMaxRetries(1),
				option.WithWorkloadIdentity(auth.WorkloadIdentity{IdentityProviderID: "synthetic-idp", ServiceAccountID: "synthetic-account", Provider: &websocketSubjectTokenProvider{}})}
			if tc.perAttempt {
				opts = append(opts, option.WithRequestTimeout(budget))
			}
			client := openai.NewClient(opts...)
			ctx, cancel := context.WithTimeout(t.Context(), 3*time.Second)
			defer cancel()
			if tc.caller {
				var stop context.CancelFunc
				ctx, stop = context.WithTimeout(ctx, budget)
				defer stop()
			}
			started := time.Now()
			var err error
			if tc.httpOnly {
				var result map[string]any
				err = client.Get(ctx, "responses", nil, &result)
			} else {
				var conn *responses.ResponseConnection
				conn, err = client.Responses.Connect(ctx, responses.ResponseConnectionOptions{})
				if conn != nil {
					defer conn.Abort()
				}
			}
			t.Logf("attempts=%d exchanges=%d elapsed=%v error=%v", attempts.Load(), exchanges.Load(), time.Since(started), err)
			if tc.wantSuccess {
				if err != nil || attempts.Load() != 2 || exchanges.Load() != 2 {
					t.Errorf("expected refreshed success, attempts=2 exchanges=2; got %v, %d, %d", err, attempts.Load(), exchanges.Load())
				}
			} else if !errors.Is(err, context.DeadlineExceeded) || attempts.Load() != 2 || exchanges.Load() != 2 {
				t.Errorf("expected total deadline during second attempt, got %v, %d, %d", err, attempts.Load(), exchanges.Load())
			}
		})
	}
}

package websockettest

import (
	"context"
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

func TestResponsesWebSocketResponseInto(t *testing.T) {
	for _, status := range []int{http.StatusSwitchingProtocols, http.StatusUnauthorized, http.StatusTooManyRequests} {
		t.Run(fmt.Sprint(status), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Request-Id", "synthetic-handshake")
				if status != http.StatusSwitchingProtocols {
					w.WriteHeader(status)
					return
				}
				socket, err := wire.Accept(w, r, nil)
				if err != nil {
					t.Errorf("accept handshake: %v", err)
					return
				}
				defer func() { _ = socket.CloseNow() }()
				_, _, _ = socket.Read(r.Context())
			}))
			defer server.Close()
			client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("synthetic-key"))
			ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
			defer cancel()
			var response *http.Response
			connection, err := client.Responses.Connect(ctx, responses.ResponseConnectionOptions{}, option.WithResponseInto(&response))
			if connection != nil {
				defer connection.Abort()
			}
			if (err == nil) != (status == http.StatusSwitchingProtocols) {
				t.Errorf("Connect(status=%d) error = %v", status, err)
			}
			if response == nil {
				t.Fatal("WithResponseInto = nil, want opening response")
			}
			if response.StatusCode != status || response.Header.Get("X-Request-Id") != "synthetic-handshake" {
				t.Errorf("opening response status/header = %d/%q, want %d/synthetic-handshake", response.StatusCode, response.Header.Get("X-Request-Id"), status)
			}
		})
	}
}

func TestResponsesWebSocketRefreshesRejectedWorkloadToken(t *testing.T) {
	for _, retries := range []int{0, 1} {
		t.Run(fmt.Sprint(retries), func(t *testing.T) {
			var upgrades, exchanges atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				upgrades.Add(1)
				if r.Header.Get("Authorization") != "Bearer synthetic-token-2" {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				socket, err := wire.Accept(w, r, nil)
				if err != nil {
					t.Errorf("accept refreshed handshake: %v", err)
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
					return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(fmt.Sprintf(`{"access_token":"synthetic-token-%d","expires_in":3600}`, token)))}, nil
				}
				return transport.RoundTrip(r)
			})}
			client := openai.NewClient(option.WithBaseURL(server.URL), option.WithHTTPClient(httpClient), option.WithMaxRetries(retries), option.WithWorkloadIdentity(auth.WorkloadIdentity{IdentityProviderID: "synthetic-idp", ServiceAccountID: "synthetic-account", Provider: &websocketSubjectTokenProvider{}}))
			ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
			defer cancel()
			connection, err := client.Responses.Connect(ctx, responses.ResponseConnectionOptions{})
			if connection != nil {
				defer connection.Abort()
			}
			if (err == nil) != (retries == 1) {
				t.Errorf("Connect(retries=%d) error = %v", retries, err)
			}
			if upgrades.Load() != int32(retries+1) || exchanges.Load() != int32(retries+1) {
				t.Errorf("Connect(retries=%d) upgrades/exchanges = %d/%d, want %d/%d", retries, upgrades.Load(), exchanges.Load(), retries+1, retries+1)
			}
		})
	}
}

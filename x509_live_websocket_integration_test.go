package openai_test

import (
	"context"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	wire "github.com/coder/websocket"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/live"
	"github.com/openai/openai-go/v3/option"
)

func TestX509LiveWebSocketEscapedSessionID(t *testing.T) {
	for _, api := range []struct {
		name    string
		path    string
		connect func(context.Context, openai.Client) error
	}{
		{"fork", "/v1/live/sessions/sess%2Fa/fork", func(ctx context.Context, client openai.Client) error {
			connection, err := client.Live.Forks.Connect(ctx, "sess/a", live.ForkConnectionOptions{})
			if err != nil {
				return err
			}
			return connection.Close()
		}},
		{"sideband", "/v1/live/sessions/sess%2Fa/attach", func(ctx context.Context, client openai.Client) error {
			connection, err := client.Live.Sideband.Connect(ctx, "sess/a", live.SidebandConnectionOptions{})
			if err != nil {
				return err
			}
			return connection.Close()
		}},
	} {
		for _, mutation := range []string{"none", "decoded path", "escaped path"} {
			t.Run(api.name+"/"+mutation, func(t *testing.T) {
				var exchanges, upgrades atomic.Int32
				config := newX509WebSocketConfig(t, x509WebSocketIssuer(t, &exchanges), func(w http.ResponseWriter, r *http.Request) {
					upgrades.Add(1)
					if r.URL.EscapedPath() != api.path || r.Header.Get("Authorization") != "Bearer synthetic-websocket-token-1" {
						t.Error("Live X.509 handshake changed the escaped path or bearer")
						w.WriteHeader(http.StatusBadRequest)
						return
					}
					socket, err := wire.Accept(w, r, nil)
					if err != nil {
						t.Errorf("accept Live WebSocket: %v", err)
						return
					}
					defer func() { _ = socket.CloseNow() }()
					_, _, _ = socket.Read(r.Context())
				})
				client := openai.NewClient(option.WithX509WorkloadIdentity(config), option.WithMaxRetries(0),
					option.WithMiddleware(func(r *http.Request, next option.MiddlewareNext) (*http.Response, error) {
						switch mutation {
						case "decoded path":
							r.URL.Path = "/v1/models"
						case "escaped path":
							r.URL.RawPath = strings.Replace(api.path, "sess%2F", "%73ess%2F", 1)
						}
						return next(r)
					}))
				ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
				defer cancel()
				err := api.connect(ctx, client)
				if mutation == "none" {
					if err != nil {
						t.Fatalf("Connect escaped session ID: %v", err)
					}
					if exchanges.Load() != 1 || upgrades.Load() != 1 {
						t.Error("expected one token exchange and upgrade")
					}
				} else if err == nil || exchanges.Load() != 0 || upgrades.Load() != 0 {
					t.Error("mutated Live path reached credential exchange or upgrade")
				}
			})
		}
	}
}

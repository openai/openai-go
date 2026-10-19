package live_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	wire "github.com/coder/websocket"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/live"
	"github.com/openai/openai-go/v3/option"
)

func TestLivePrimaryAndForkStartupAndDrain(t *testing.T) {
	for _, role := range []string{"primary", "fork"} {
		t.Run(role, func(t *testing.T) {
			closed := make(chan struct{})
			server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				defer close(closed)
				path := "/live/sessions"
				if role == "fork" {
					path += "/live_stored/fork"
				}
				if r.URL.Path != path || r.Header.Get("Authorization") != "Bearer synthetic" {
					t.Errorf("unexpected request policy")
				}
				socket, err := wire.Accept(w, r, nil)
				if err != nil {
					t.Error(err)
					return
				}
				defer func() { _ = socket.CloseNow() }()
				_, command, err := socket.Read(r.Context())
				if err != nil {
					t.Error(err)
					return
				}
				var start struct {
					Type    string
					Session map[string]any
				}
				if err := json.Unmarshal(command, &start); err != nil {
					t.Error(err)
					return
				}
				if start.Type != "session.start" {
					t.Errorf("wrong startup command: %s", command)
				}
				if role == "primary" && start.Session["model"] != "synthetic" {
					t.Errorf("primary model lost")
				}
				if role == "fork" && len(start.Session) != 0 {
					t.Errorf("fork overrides changed: %s", command)
				}
				for _, frame := range []string{
					`{"type":"session.started","event_id":"start","session":{"id":"live_new","model":"synthetic","expires_at":2000000000,"status":"active"}}`,
					`{"type":"session.output_transcript.delta","event_id":"text","delta":"hello","start_ms":0,"end_ms":20}`,
					`{"type":"transport.ringing","event_id":"ring","session_id":"live_new"}`,
					`{"type":"session.closed","event_id":"end"}`,
					`{"type":"future.live.event","marker":42}`,
				} {
					if err := socket.Write(r.Context(), wire.MessageText, []byte(frame)); err != nil {
						t.Error(err)
						return
					}
				}
				_, _, _ = socket.Read(r.Context())
			}))
			defer server.Close()
			client := openai.NewClient(option.WithBaseURL(server.URL), option.WithHTTPClient(server.Client()), option.WithAPIKey("synthetic"))
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			var recv func() (string, string, error)
			var closeConnection func() error
			if role == "primary" {
				connection, err := client.Live.Connect(ctx, live.LiveConnectionOptions{})
				if err != nil {
					t.Fatal(err)
				}
				defer connection.Abort()
				var start live.ClientEventUnionParam
				if err := json.Unmarshal([]byte(`{"type":"session.start","session":{"model":"synthetic"}}`), &start); err != nil {
					t.Fatal(err)
				}
				if err := connection.Send(ctx, start); err != nil {
					t.Fatal(err)
				}
				recv = func() (string, string, error) {
					e, err := connection.Recv(ctx)
					if e.Type == "transport.ringing" && e.AsServerEventTransportRinging().SessionID != "live_new" {
						t.Error("primary ringing session lost")
					}
					return e.Type, e.AsOutputTranscriptDeltaEvent().Delta, err
				}
				closeConnection = connection.Close
			} else {
				connection, err := client.Live.Forks.Connect(ctx, "live_stored", live.ForkConnectionOptions{})
				if err != nil {
					t.Fatal(err)
				}
				defer connection.Abort()
				var start live.ForkClientEventUnionParam
				if err := json.Unmarshal([]byte(`{"type":"session.start","session":{}}`), &start); err != nil {
					t.Fatal(err)
				}
				if err := connection.Send(ctx, start); err != nil {
					t.Fatal(err)
				}
				recv = func() (string, string, error) {
					e, err := connection.Recv(ctx)
					if e.Type == "transport.ringing" && e.AsTransportRinging().SessionID != "live_new" {
						t.Error("fork ringing session lost")
					}
					return e.Type, e.AsSessionOutputTranscriptDelta().Delta, err
				}
				closeConnection = connection.Close
			}
			for _, want := range []string{"session.started", "session.output_transcript.delta", "transport.ringing", "session.closed", "future.live.event"} {
				kind, text, err := recv()
				if err != nil || kind != want {
					t.Fatalf("event %s: got %s, %v", want, kind, err)
				}
				if kind == "session.output_transcript.delta" && text != "hello" {
					t.Fatalf("transcript lost: %q", text)
				}
			}
			if err := closeConnection(); err != nil {
				t.Fatal(err)
			}
			select {
			case <-closed:
			case <-ctx.Done():
				t.Fatal("socket not cleaned up")
			}
		})
	}
}

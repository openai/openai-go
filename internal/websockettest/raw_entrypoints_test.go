package websockettest

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	wire "github.com/coder/websocket"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/live"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

// Only adapt the differing typed event unions; all operations use public SDK entrypoints.
type rawEntrypoints struct {
	send  func(context.Context, json.RawMessage) error
	recv  func(context.Context) (json.RawMessage, error)
	typed func(context.Context) (string, error)
	abort func()
}

func rawEntrypointsFor[E interface{ RawJSON() string }](c interface {
	SendRaw(context.Context, json.RawMessage) error
	RecvRaw(context.Context) (json.RawMessage, error)
	Recv(context.Context) (E, error)
	Abort()
}) rawEntrypoints {
	return rawEntrypoints{send: c.SendRaw, recv: c.RecvRaw, abort: c.Abort, typed: func(ctx context.Context) (string, error) {
		event, err := c.Recv(ctx)
		return event.RawJSON(), err
	}}
}

func TestRawWebsocketEntrypoints(t *testing.T) {
	for _, test := range []struct {
		name    string
		connect func(context.Context, *openai.Client) (rawEntrypoints, error)
	}{
		{name: "responses", connect: func(ctx context.Context, client *openai.Client) (rawEntrypoints, error) {
			c, err := client.Responses.Connect(ctx, responses.ResponseConnectionOptions{})
			if err != nil {
				return rawEntrypoints{}, err
			}
			return rawEntrypointsFor(c), nil
		}},
		{name: "live", connect: func(ctx context.Context, client *openai.Client) (rawEntrypoints, error) {
			c, err := client.Live.Connect(ctx, live.LiveConnectionOptions{})
			if err != nil {
				return rawEntrypoints{}, err
			}
			return rawEntrypointsFor(c), nil
		}},
		{name: "sideband", connect: func(ctx context.Context, client *openai.Client) (rawEntrypoints, error) {
			c, err := client.Live.Sideband.Connect(ctx, "session_synthetic", live.SidebandConnectionOptions{})
			if err != nil {
				return rawEntrypoints{}, err
			}
			return rawEntrypointsFor(c), nil
		}},
		{name: "fork", connect: func(ctx context.Context, client *openai.Client) (rawEntrypoints, error) {
			c, err := client.Live.Forks.Connect(ctx, "session_synthetic", live.ForkConnectionOptions{})
			if err != nil {
				return rawEntrypoints{}, err
			}
			return rawEntrypointsFor(c), nil
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			t.Cleanup(cancel)
			captured := make(chan []byte, 1)
			release := make(chan struct{})
			frames := []string{
				` { "type" : "response.future", "value": 1.00, "escaped": "\u0061" } `,
				`{"type":"response.future","sequence_number":2}`,
				`{"type":"response.future","sequence_number":3}`,
			}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				socket, err := wire.Accept(w, r, nil)
				if err != nil {
					t.Errorf("Accept() = %v, want success", err)
					return
				}
				defer func() { _ = socket.CloseNow() }()
				_, data, err := socket.Read(ctx)
				if err != nil {
					t.Errorf("Read(raw command) = %v, want success", err)
					return
				}
				captured <- data
				select {
				case <-release:
				case <-ctx.Done():
					return
				}
				for _, frame := range frames {
					if err := socket.Write(ctx, wire.MessageText, []byte(frame)); err != nil {
						t.Errorf("Write(%q) = %v, want success", frame, err)
						return
					}
				}
				_, _, _ = socket.Read(ctx)
			}))
			t.Cleanup(server.Close)
			client := openai.NewClient(option.WithUnsafeAllowHTTP(), option.WithBaseURL(server.URL), option.WithAPIKey("synthetic-key"))
			c, err := test.connect(ctx, &client)
			if err != nil {
				t.Fatalf("Connect(local server) = %v, want success", err)
			}
			t.Cleanup(c.abort)
			// Cancel before frames arrive: only this receive should be canceled.
			canceled, stop := context.WithCancel(ctx)
			stop()
			if _, err = c.recv(canceled); !errors.Is(err, context.Canceled) {
				t.Fatalf("RecvRaw(canceled) = %v, want context.Canceled", err)
			}
			payload := json.RawMessage(` { "type" : "synthetic.command", "value": 1.00, "escaped": "\u0061" } `)
			if err = c.send(ctx, payload); err != nil {
				t.Fatalf("SendRaw(%q) = %v, want success", payload, err)
			}
			select {
			case got := <-captured:
				if string(got) != string(payload) {
					t.Errorf("SendRaw wire bytes = %q, want %q", got, payload)
				}
			case <-ctx.Done():
				t.Fatal("SendRaw capture timed out")
			}
			close(release)
			first, err := c.recv(ctx)
			if err != nil || string(first) != frames[0] {
				t.Fatalf("RecvRaw(first) = %q, %v, want %q, nil", first, err, frames[0])
			}
			if _, err = c.recv(canceled); !errors.Is(err, context.Canceled) {
				t.Fatalf("RecvRaw(canceled between events) = %v, want context.Canceled", err)
			}
			second, err := c.typed(ctx)
			if err != nil || second != frames[1] {
				t.Fatalf("Recv(after raw) = %q, %v, want %q, nil", second, err, frames[1])
			}
			third, err := c.recv(ctx)
			if err != nil || string(third) != frames[2] {
				t.Fatalf("RecvRaw(after typed) = %q, %v, want %q, nil", third, err, frames[2])
			}
		})
	}
}

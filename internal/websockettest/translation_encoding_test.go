package websockettest

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	wire "github.com/coder/websocket"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/realtime"
)

var translationEncodingFailure = errors.New("synthetic translation encoder failure")

type translationFailingEncoder struct{ panicInstead bool }

func (e translationFailingEncoder) MarshalJSON() ([]byte, error) {
	if e.panicInstead {
		panic(translationEncodingFailure)
	}
	return nil, translationEncodingFailure
}

func TestTranslationEncodingFailureReleasesSendCapacity(t *testing.T) {
	for _, tc := range []struct {
		name         string
		panicInstead bool
	}{
		{"returned_error_control", false},
		{"recovered_panic", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			conn, ctx, result := translationPeer(t, realtime.TranslationConnectionOptions{MaxPendingSends: 1},
				func(ctx context.Context, socket *wire.Conn) error {
					for _, want := range []string{"client.future", "session.close"} {
						if err := translationCommand(ctx, socket, want); err != nil {
							return err
						}
					}
					if err := socket.Write(ctx, wire.MessageText, []byte(`{"type":"session.closed","event_id":"done"}`)); err != nil {
						return err
					}
					return translationClosed(ctx, socket)
				},
			)
			typ := param.Override[realtime.RealtimeTranslationClientEventUnionParam](translationFailingEncoder{tc.panicInstead})
			if tc.panicInstead {
				func() {
					defer func() {
						if got := recover(); got != translationEncodingFailure { //nolint:errorlint // Assert exact panic identity, not wrapped-error equivalence.
							t.Errorf("panic = %v; want original encoder panic", got)
						}
					}()
					if err := conn.Send(ctx, typ); err != nil {
						t.Errorf("panic converted to unexpected error: %v", err)
					}
				}()
			} else if err := conn.Send(ctx, typ); !errors.Is(err, translationEncodingFailure) {
				t.Fatalf("Send(encoding fails) = %v; want original error", err)
			}
			if err := conn.SendRaw(ctx, json.RawMessage(`{"type":"client.future"}`)); err != nil {
				t.Fatalf("SendRaw after failed, never-attempted encoding = %v; want success", err)
			}
			if err := conn.Finish(ctx, func(realtime.RealtimeTranslationServerEventUnion) error { return nil }); err != nil {
				t.Fatalf("Finish after failed, never-attempted encoding = %v; want success", err)
			}
			translationResult(t, ctx, result)
		})
	}
}

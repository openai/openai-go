package websockettest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"runtime"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	wire "github.com/coder/websocket"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/packages/websocket"
	"github.com/openai/openai-go/v3/realtime"
)

func translationPeer(t *testing.T, options realtime.TranslationConnectionOptions, script func(context.Context, *wire.Conn) error) (*realtime.TranslationSession, context.Context, <-chan error) {
	t.Helper()
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	t.Cleanup(cancel)
	result := make(chan error, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		socket, err := wire.Accept(w, r, nil)
		if err != nil {
			result <- err
			return
		}
		defer func() { _ = socket.CloseNow() }()
		result <- script(ctx, socket)
	}))
	t.Cleanup(server.Close)
	client := openai.NewClient(option.WithBaseURL(server.URL+"/v1"), option.WithAPIKey("synthetic-key"))
	options.CloseTimeout = 100 * time.Millisecond
	connection, err := client.Realtime.ConnectTranslation(ctx, "gpt-realtime-translate", options)
	if err != nil {
		t.Fatalf("ConnectTranslation(loopback) = %v, want success", err)
	}
	t.Cleanup(connection.Abort)
	return connection, ctx, result
}

func translationCommand(ctx context.Context, socket *wire.Conn, want string) error {
	_, data, err := socket.Read(ctx)
	if err != nil {
		return err
	}
	var command struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &command); err != nil {
		return err
	}
	if command.Type != want {
		return fmt.Errorf("command type = %q, want %q", command.Type, want)
	}
	return nil
}

func translationClosed(ctx context.Context, socket *wire.Conn) error {
	_, data, err := socket.Read(ctx)
	if wire.CloseStatus(err) != wire.StatusNormalClosure {
		return fmt.Errorf("after drain got data %q, error %v; want normal close", data, err)
	}
	return nil
}

func translationResult(t *testing.T, ctx context.Context, result <-chan error) {
	t.Helper()
	select {
	case err := <-result:
		if err != nil {
			t.Errorf("Translation peer = %v, want success", err)
		}
	case <-ctx.Done():
		t.Errorf("Translation peer = %v, want completion", ctx.Err())
	}
}

func TestTranslationFinishForwardsBufferedEventsAndOwnsDelivery(t *testing.T) {
	allWritten := make(chan struct{})
	want := []string{"session.input_transcript.delta", "session.output_transcript.delta", "session.output_audio.delta", "future.event", "error", "session.closed"}
	terminalDelivered := atomic.Bool{}
	connection, ctx, result := translationPeer(t, realtime.TranslationConnectionOptions{}, func(ctx context.Context, socket *wire.Conn) error {
		if err := socket.Write(ctx, wire.MessageText, []byte(`{"type":"session.created","event_id":"created","session":{"id":"session","type":"translation","model":"gpt-realtime-translate","audio":{},"expires_at":123}}`)); err != nil {
			return err
		}
		for _, expected := range []string{
			`{"type":"session.update","session":{"audio":{"output":{"language":"es"}}}}`,
			`{"type":"session.input_audio_buffer.append","audio":"AAAA"}`,
			`{"type":"client.future","extra":true}`,
			`{"type":"session.close"}`,
		} {
			_, data, err := socket.Read(ctx)
			if err != nil {
				return err
			}
			var got, want map[string]any
			if err := json.Unmarshal(data, &got); err != nil {
				return err
			}
			if err := json.Unmarshal([]byte(expected), &want); err != nil {
				return err
			}
			if !reflect.DeepEqual(got, want) {
				return fmt.Errorf("wire command = %s, want %s", data, expected)
			}
		}
		for _, kind := range want {
			data, err := json.Marshal(map[string]any{"type": kind, "event_id": "synthetic", "delta": "AAAA", "extra": "preserved", "error": map[string]string{"type": "invalid_request_error", "message": "synthetic API error"}})
			if err != nil {
				return err
			}
			if err := socket.Write(ctx, wire.MessageText, data); err != nil {
				return err
			}
		}
		close(allWritten)
		if err := translationClosed(ctx, socket); err != nil {
			return err
		}
		if !terminalDelivered.Load() {
			return errors.New("transport closed before terminal callback")
		}
		return nil
	})
	event, receiveErr := connection.Recv(ctx)
	if receiveErr != nil || event.Type != "session.created" {
		t.Fatalf("Recv(initial) = %s, %v, want session.created", event.RawJSON(), receiveErr)
	}
	created, ok := event.AsAny().(realtime.RealtimeTranslationSessionCreatedEvent)
	if !ok || created.Session.ID != "session" {
		t.Error("Recv(initial).AsAny() did not preserve the typed session.created variant")
	}
	if !event.JSON.OfRealtimeTranslationSessionCreatedEvent.Valid() || event.JSON.OfRealtimeTranslationSessionUpdatedEvent.Valid() || event.OfRealtimeTranslationSessionUpdatedEvent.Type != "" {
		t.Error("Recv(initial) populated an unrelated event variant")
	}
	update := realtime.RealtimeTranslationClientEventParamOfSessionUpdate(realtime.RealtimeTranslationSessionUpdateRequestParam{
		Audio: realtime.RealtimeTranslationSessionUpdateRequestAudioParam{
			Output: realtime.RealtimeTranslationSessionUpdateRequestAudioOutputParam{Language: openai.String("es")},
		},
	})
	if err := connection.Send(ctx, update); err != nil {
		t.Fatalf("Send(update) = %v, want success", err)
	}
	if err := connection.Send(ctx, realtime.RealtimeTranslationClientEventParamOfSessionInputAudioBufferAppend("AAAA")); err != nil {
		t.Fatalf("Send(audio) = %v, want success", err)
	}
	if err := connection.SendRaw(ctx, json.RawMessage(`{"type":"client.future","extra":true}`)); err != nil {
		t.Fatalf("SendRaw(future) = %v, want success", err)
	}
	var got []string
	err := connection.Finish(ctx, func(event realtime.RealtimeTranslationServerEventUnion) error {
		select {
		case <-allWritten:
		case <-ctx.Done():
			return ctx.Err()
		}
		if _, err := connection.Recv(ctx); !errors.Is(err, realtime.ErrTranslationReceiving) {
			return fmt.Errorf("competing Recv = %v, want receiver rejection", err)
		}
		if err := connection.Finish(ctx, func(realtime.RealtimeTranslationServerEventUnion) error { return nil }); !errors.Is(err, realtime.ErrTranslationReceiving) {
			return fmt.Errorf("competing Finish = %v, want receiver rejection", err)
		}
		if err := connection.SendRaw(ctx, json.RawMessage(`{"type":"session.input_audio_buffer.append","audio":"late"}`)); !errors.Is(err, realtime.ErrTranslationInputClosed) {
			return fmt.Errorf("late Send = %v, want closed input", err)
		}
		kind := event.Type
		got = append(got, kind)
		switch kind {
		case "future.event":
			if !strings.Contains(event.RawJSON(), `"extra":"preserved"`) || event.AsAny() != nil {
				return errors.New("future event did not preserve unknown JSON")
			}
		case "error":
			if event.AsRealtimeErrorEvent().Error.Message != "synthetic API error" {
				return errors.New("API error was not inspectable")
			}
		case "session.output_audio.delta":
			if event.AsRealtimeTranslationOutputAudioDeltaEvent().Delta != "AAAA" {
				return errors.New("typed audio delta was lost")
			}
			if !event.JSON.OfRealtimeTranslationOutputAudioDeltaEvent.Valid() || event.JSON.OfRealtimeTranslationInputTranscriptDeltaEvent.Valid() || event.OfRealtimeTranslationInputTranscriptDeltaEvent.Delta != "" {
				return errors.New("Finish populated an unrelated event variant")
			}
		case "session.closed":
			terminalDelivered.Store(true)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Finish(buffered output) = %v, want success", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Finish events = %v, want %v", got, want)
	}
	if err := connection.Finish(ctx, func(realtime.RealtimeTranslationServerEventUnion) error { return errors.New("duplicate callback") }); err != nil {
		t.Errorf("Finish(repeated) = %v, want success", err)
	}
	if err := connection.Finish(ctx, nil); err != nil {
		t.Errorf("Finish(repeated, nil) = %v, want success", err)
	}
	translationResult(t, ctx, result)
}

func TestTranslationRawSendPreservesWireBytesAndBudget(t *testing.T) {
	const raw = " \n{\"type\":\"client.future\", \"payload\":\"<>&<>&\"} \n"
	const rawClose = " \n{\"type\":\"session.close\", \"payload\":\"<>&\"} \n"
	connection, ctx, result := translationPeer(t, realtime.TranslationConnectionOptions{MaxMessageBytes: int64(len(raw))}, func(ctx context.Context, socket *wire.Conn) error {
		for _, want := range []string{raw, rawClose} {
			_, data, err := socket.Read(ctx)
			if err != nil {
				return err
			}
			if string(data) != want {
				return fmt.Errorf("raw wire command = %q, want %q", data, want)
			}
		}
		if err := socket.Write(ctx, wire.MessageText, []byte(`{"type":"session.closed","event_id":"closed"}`)); err != nil {
			return err
		}
		return translationClosed(ctx, socket)
	})
	if err := connection.SendRaw(ctx, json.RawMessage(raw)); err != nil {
		t.Fatalf("SendRaw(at byte limit) = %v, want success", err)
	}
	if err := connection.SendRaw(ctx, json.RawMessage(rawClose)); err != nil {
		t.Fatalf("SendRaw(session.close) = %v, want success", err)
	}
	if err := connection.Finish(ctx, func(realtime.RealtimeTranslationServerEventUnion) error { return nil }); err != nil {
		t.Fatalf("Finish(after raw send) = %v, want success", err)
	}
	translationResult(t, ctx, result)
}

func TestTranslationUsesExactTypeKeyForControlEvents(t *testing.T) {
	const rawClose = `{"type":"session.close","TYPE":"future.command"}`
	connection, ctx, result := translationPeer(t, realtime.TranslationConnectionOptions{}, func(ctx context.Context, socket *wire.Conn) error {
		_, data, err := socket.Read(ctx)
		if err != nil {
			return err
		}
		if string(data) != rawClose {
			return fmt.Errorf("close command = %q, want %q", data, rawClose)
		}
		for _, frame := range []string{`{"type":"future.event","TYPE":"session.closed"}`, `{"type":"session.closed","event_id":"closed"}`} {
			if err := socket.Write(ctx, wire.MessageText, []byte(frame)); err != nil {
				return err
			}
		}
		return translationClosed(ctx, socket)
	})
	if err := connection.SendRaw(ctx, json.RawMessage(rawClose)); err != nil {
		t.Fatalf("SendRaw(close with future field) = %v, want success", err)
	}
	var kinds []string
	if err := connection.Finish(ctx, func(event realtime.RealtimeTranslationServerEventUnion) error {
		kinds = append(kinds, event.Type)
		return nil
	}); err != nil {
		t.Fatalf("Finish(conflicting future fields) = %v, want success", err)
	}
	if want := []string{"future.event", "session.closed"}; !reflect.DeepEqual(kinds, want) {
		t.Errorf("Finish delivered event types %v, want %v", kinds, want)
	}
	translationResult(t, ctx, result)
}

func TestTranslationRawReceiveCancellationAndRemoteTerminal(t *testing.T) {
	release := make(chan struct{})
	connection, ctx, result := translationPeer(t, realtime.TranslationConnectionOptions{}, func(ctx context.Context, socket *wire.Conn) error {
		select {
		case <-release:
		case <-ctx.Done():
			return ctx.Err()
		}
		for _, data := range []string{`{"type":"future.event","unknown":[1,true]}`, `{"type":"session.closed","event_id":"closed"}`} {
			if err := socket.Write(ctx, wire.MessageText, []byte(data)); err != nil {
				return err
			}
		}
		return translationClosed(ctx, socket)
	})
	canceled, cancel := context.WithCancel(ctx)
	cancel()
	if _, err := connection.Recv(canceled); !errors.Is(err, context.Canceled) {
		t.Errorf("Recv(canceled) = %v, want cancellation", err)
	}
	close(release)
	data, err := connection.RecvRaw(ctx)
	if err != nil || string(data) != `{"type":"future.event","unknown":[1,true]}` {
		t.Errorf("RecvRaw(future) = %s, %v, want unchanged JSON", data, err)
	}
	if event, err := connection.Recv(ctx); err != nil || event.Type != "session.closed" {
		t.Fatalf("Recv(terminal) = %s, %v, want session.closed", event.RawJSON(), err)
	}
	if _, err := connection.RecvRaw(ctx); !errors.Is(err, io.EOF) {
		t.Errorf("RecvRaw(after terminal) = %v, want EOF", err)
	}
	if err := connection.Finish(ctx, func(realtime.RealtimeTranslationServerEventUnion) error {
		return errors.New("terminal delivered twice")
	}); err != nil {
		t.Errorf("Finish(remote terminal) = %v, want success without another close command", err)
	}
	translationResult(t, ctx, result)
}

func TestTranslationFinishAfterConsumedTerminal(t *testing.T) {
	for _, receive := range []string{"typed", "raw"} {
		for _, expiration := range []string{"canceled", "deadline"} {
			t.Run(receive+"/"+expiration, func(t *testing.T) {
				connection, ctx, result := translationPeer(t, realtime.TranslationConnectionOptions{}, func(ctx context.Context, socket *wire.Conn) error {
					if err := socket.Write(ctx, wire.MessageText, []byte(`{"type":"session.closed"}`)); err != nil {
						return err
					}
					return translationClosed(ctx, socket)
				})
				if receive == "raw" {
					if _, err := connection.RecvRaw(ctx); err != nil {
						t.Fatal(err)
					}
				} else if _, err := connection.Recv(ctx); err != nil {
					t.Fatal(err)
				}
				finishCtx, cancel := context.WithCancel(ctx)
				if expiration == "deadline" {
					cancel()
					finishCtx, cancel = context.WithDeadline(ctx, time.Now().Add(-time.Second))
				}
				cancel()
				if err := connection.Finish(finishCtx, func(realtime.RealtimeTranslationServerEventUnion) error {
					return errors.New("terminal event delivered twice")
				}); err != nil {
					t.Errorf("Finish after terminal = %v, want successful cleanup", err)
				}
				if err := connection.Finish(finishCtx, nil); err != nil {
					t.Errorf("repeated Finish = %v, want recorded success", err)
				}
				translationResult(t, ctx, result)
			})
		}
	}
}

func TestTranslationFinishTerminalCallbackCancelsContext(t *testing.T) {
	for _, callbackErr := range []error{nil, errors.New("callback failed")} {
		t.Run(fmt.Sprint(callbackErr), func(t *testing.T) {
			connection, ctx, result := translationPeer(t, realtime.TranslationConnectionOptions{}, func(ctx context.Context, socket *wire.Conn) error {
				if err := translationCommand(ctx, socket, "session.close"); err != nil {
					return err
				}
				if err := socket.Write(ctx, wire.MessageText, []byte(`{"type":"session.closed"}`)); err != nil {
					return err
				}
				if callbackErr == nil {
					return translationClosed(ctx, socket)
				}
				if _, _, err := socket.Read(ctx); err == nil {
					return errors.New("received an extra command after callback failure")
				}
				return nil
			})
			finishCtx, cancel := context.WithCancel(ctx)
			defer cancel()
			calls := 0
			err := connection.Finish(finishCtx, func(event realtime.RealtimeTranslationServerEventUnion) error {
				calls++
				if event.Type != "session.closed" {
					t.Errorf("callback event = %q, want session.closed", event.Type)
				}
				cancel()
				return callbackErr
			})
			if !errors.Is(err, callbackErr) {
				t.Errorf("Finish after terminal callback = %v, want %v", err, callbackErr)
			}
			if again := connection.Finish(ctx, nil); !errors.Is(again, callbackErr) {
				t.Errorf("repeated Finish = %v, want recorded %v", again, callbackErr)
			}
			if calls != 1 {
				t.Errorf("callback calls = %d, want 1", calls)
			}
			translationResult(t, ctx, result)
		})
	}
}

func TestTranslationFinishFailureNeverReplays(t *testing.T) {
	for _, scenario := range []struct {
		name  string
		frame string
	}{
		{name: "early EOF"},
		{name: "malformed JSON", frame: `{"type":`},
		{name: "missing type", frame: `{}`},
		{name: "uppercase type", frame: `{"TYPE":"session.closed"}`},
		{name: "wrong type", frame: `{"type":123}`},
		{name: "callback", frame: `{"type":"session.output_audio.delta","event_id":"audio","delta":"AAAA"}`},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			connection, ctx, result := translationPeer(t, realtime.TranslationConnectionOptions{}, func(ctx context.Context, socket *wire.Conn) error {
				if err := translationCommand(ctx, socket, "session.close"); err != nil {
					return err
				}
				if scenario.frame == "" {
					return socket.Close(wire.StatusNormalClosure, "")
				}
				if err := socket.Write(ctx, wire.MessageText, []byte(scenario.frame)); err != nil {
					return err
				}
				if _, _, err := socket.Read(ctx); err == nil {
					return errors.New("received extra command after failure")
				}
				return nil
			})
			callbackErr := errors.New("caller processing failed")
			err := connection.Finish(ctx, func(realtime.RealtimeTranslationServerEventUnion) error { return callbackErr })
			if err == nil {
				t.Fatal("Finish(failure) = nil, want failure")
			}
			if scenario.name == "callback" && !errors.Is(err, callbackErr) {
				t.Errorf("Finish(callback failure) = %v, want original error", err)
			}
			if scenario.name == "malformed JSON" {
				var syntax *json.SyntaxError
				if !errors.As(err, &syntax) {
					t.Errorf("Finish(malformed JSON) = %v, want syntax error", err)
				}
			}
			if again := connection.Finish(ctx, func(realtime.RealtimeTranslationServerEventUnion) error { return nil }); !errors.Is(again, err) {
				t.Errorf("Finish(repeated failure) = %v, want original %v", again, err)
			}
			if again := connection.Finish(ctx, nil); !errors.Is(again, err) {
				t.Errorf("Finish(repeated failure, nil) = %v, want original %v", again, err)
			}
			translationResult(t, ctx, result)
		})
	}
}

func TestTranslationFinishCancellationAndDeadline(t *testing.T) {
	for _, deadline := range []bool{false, true} {
		t.Run(fmt.Sprintf("deadline=%t", deadline), func(t *testing.T) {
			closeSeen := make(chan struct{})
			connection, ctx, result := translationPeer(t, realtime.TranslationConnectionOptions{}, func(ctx context.Context, socket *wire.Conn) error {
				if err := translationCommand(ctx, socket, "session.close"); err != nil {
					return err
				}
				close(closeSeen)
				if _, _, err := socket.Read(ctx); err == nil {
					return errors.New("received extra command after cancellation")
				}
				return nil
			})
			finishCtx, cancel := context.WithCancel(ctx)
			if deadline {
				cancel()
				finishCtx, cancel = context.WithTimeout(ctx, 100*time.Millisecond)
			}
			defer cancel()
			finished := make(chan error, 1)
			go func() {
				finished <- connection.Finish(finishCtx, func(realtime.RealtimeTranslationServerEventUnion) error { return nil })
			}()
			select {
			case <-closeSeen:
			case <-ctx.Done():
				t.Fatal("Finish never sent session.close")
			}
			if !deadline {
				cancel()
			}
			want := context.Canceled
			if deadline {
				want = context.DeadlineExceeded
			}
			select {
			case err := <-finished:
				if !errors.Is(err, want) {
					t.Errorf("Finish(interrupted) = %v, want %v", err, want)
				}
			case <-ctx.Done():
				t.Fatal("Finish did not stop after context expiration")
			}
			translationResult(t, ctx, result)
		})
	}
}

type translationBlockingCommand struct {
	started chan<- struct{}
	release <-chan struct{}
	done    <-chan struct{}
}

func (c translationBlockingCommand) MarshalJSON() ([]byte, error) {
	close(c.started)
	select {
	case <-c.release:
		return []byte(`{"type":"session.input_audio_buffer.append","audio":"AAAA"}`), nil
	case <-c.done:
		return nil, context.Canceled
	}
}

func TestTranslationFinishWaitsForAdmittedSend(t *testing.T) {
	connection, ctx, result := translationPeer(t, realtime.TranslationConnectionOptions{MaxPendingSends: 1}, func(ctx context.Context, socket *wire.Conn) error {
		for _, kind := range []string{"session.input_audio_buffer.append", "session.close"} {
			if err := translationCommand(ctx, socket, kind); err != nil {
				return err
			}
		}
		if err := socket.Write(ctx, wire.MessageText, []byte(`{"type":"session.closed","event_id":"closed"}`)); err != nil {
			return err
		}
		return translationClosed(ctx, socket)
	})
	started, release := make(chan struct{}), make(chan struct{})
	sent, finished := make(chan error, 1), make(chan error, 1)
	command := param.Override[realtime.RealtimeTranslationClientEventUnionParam](translationBlockingCommand{started: started, release: release, done: ctx.Done()})
	go func() { sent <- connection.Send(ctx, command) }()
	select {
	case <-started:
	case <-ctx.Done():
		t.Fatal("Send never began serialization")
	}
	if err := connection.SendRaw(ctx, json.RawMessage(`{"type":"client.future"}`)); !errors.Is(err, websocket.ErrPendingLimit) {
		t.Errorf("Send(concurrent serialization) = %v, want pending limit", err)
	}
	go func() {
		finished <- connection.Finish(ctx, func(realtime.RealtimeTranslationServerEventUnion) error { return nil })
	}()
	for {
		err := connection.SendRaw(ctx, json.RawMessage(`{"type":"client.future"}`))
		if errors.Is(err, realtime.ErrTranslationInputClosed) {
			break
		}
		if !errors.Is(err, websocket.ErrPendingLimit) {
			t.Fatalf("Send(while finish starts) = %v, want pending or input-closed error", err)
		}
		if ctx.Err() != nil {
			t.Fatal("Finish never stopped input admission")
		}
		runtime.Gosched()
	}
	close(release)
	if err := <-sent; err != nil {
		t.Errorf("Send(admitted before Finish) = %v, want success", err)
	}
	if err := <-finished; err != nil {
		t.Errorf("Finish(after admitted send) = %v, want success", err)
	}
	translationResult(t, ctx, result)
}

func TestTranslationExplicitSessionCloseAndTransportClose(t *testing.T) {
	for _, protocol := range []bool{true, false} {
		t.Run(fmt.Sprintf("protocol=%t", protocol), func(t *testing.T) {
			connection, ctx, result := translationPeer(t, realtime.TranslationConnectionOptions{}, func(ctx context.Context, socket *wire.Conn) error {
				if protocol {
					if err := translationCommand(ctx, socket, "session.close"); err != nil {
						return err
					}
					if err := socket.Write(ctx, wire.MessageText, []byte(`{"type":"session.closed","event_id":"closed"}`)); err != nil {
						return err
					}
				}
				return translationClosed(ctx, socket)
			})
			if protocol {
				if err := connection.Send(ctx, realtime.RealtimeTranslationClientEventUnionParam{OfSessionClose: &realtime.RealtimeTranslationSessionCloseEventParam{}}); err != nil {
					t.Fatalf("Send(session.close) = %v, want success", err)
				}
				if err := connection.Finish(ctx, func(realtime.RealtimeTranslationServerEventUnion) error { return nil }); err != nil {
					t.Errorf("Finish(after explicit session.close) = %v, want success", err)
				}
			} else if err := connection.Close(); err != nil {
				t.Errorf("Close(transport only) = %v, want success", err)
			}
			translationResult(t, ctx, result)
		})
	}
}

func TestTranslationCleanupDeadline(t *testing.T) {
	release := make(chan struct{})
	connection, ctx, result := translationPeer(t, realtime.TranslationConnectionOptions{}, func(ctx context.Context, socket *wire.Conn) error {
		if err := translationCommand(ctx, socket, "session.close"); err != nil {
			return err
		}
		if err := socket.Write(ctx, wire.MessageText, []byte(`{"type":"session.closed","event_id":"closed"}`)); err != nil {
			return err
		}
		// Deliberately do not read or acknowledge the transport closing handshake.
		select {
		case <-release:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	})
	start := time.Now()
	if err := connection.Finish(ctx, func(realtime.RealtimeTranslationServerEventUnion) error { return nil }); err != nil {
		t.Errorf("Finish(unresponsive close peer) = %v, want success", err)
	}
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Errorf("Finish cleanup = %v, want bounded by close timeout", elapsed)
	}
	close(release)
	translationResult(t, ctx, result)
}

type translationHTTPDoer struct {
	client *http.Client
	calls  atomic.Int32
}

func (c *translationHTTPDoer) Do(*http.Request) (*http.Response, error) {
	c.calls.Add(1)
	return nil, errors.New("unexpected custom HTTP Do call")
}

func (c *translationHTTPDoer) WebSocketHTTPClient() *http.Client { return c.client }

func TestTranslationOpeningPreservesOptionsAndCustomTLSClient(t *testing.T) {
	var upgrades, transports atomic.Int32
	var originalRequests atomic.Int32
	original := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		originalRequests.Add(1)
		w.WriteHeader(http.StatusForbidden)
	}))
	defer original.Close()
	result := make(chan error, 1)
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrades.Add(1)
		if r.URL.Path != "/v1/realtime/translations" || !reflect.DeepEqual(r.URL.Query()["model"], []string{"gpt-realtime-translate"}) || r.URL.Query().Get("base") != "retained" || !reflect.DeepEqual(r.URL.Query()["constructor"], []string{"once"}) {
			result <- fmt.Errorf("opening path/query = %s, want Translation path, single explicit model, base query", r.URL)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if r.Header.Get("Authorization") != "Bearer synthetic-service-key" || r.Header.Get("OpenAI-Organization") != "synthetic-org" || r.Header.Get("OpenAI-Project") != "synthetic-project" || r.Header.Get("X-Service") != "service" || r.Header.Get("X-Custom") != "connection" || r.Header.Get("X-Transport") != "preserved" {
			result <- errors.New("opening did not preserve auth, inherited headers, or per-connection precedence")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		socket, err := wire.Accept(w, r, nil)
		if err != nil {
			result <- err
			return
		}
		defer func() { _ = socket.CloseNow() }()
		ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
		defer cancel()
		if err := socket.Write(ctx, wire.MessageText, []byte(`{"type":"session.created","event_id":"created","session":{}}`)); err != nil {
			result <- err
			return
		}
		result <- translationClosed(ctx, socket)
	}))
	defer server.Close()
	type contextKey struct{}
	httpClient := server.Client()
	transport := httpClient.Transport
	httpClient.Transport = websocketAuthTransport(func(r *http.Request) (*http.Response, error) {
		transports.Add(1)
		if r.Context().Value(contextKey{}) != "preserved" {
			return nil, errors.New("prepared request context was lost")
		}
		r.Header.Set("X-Transport", "preserved")
		return transport.RoundTrip(r)
	})
	doer := &translationHTTPDoer{client: httpClient}
	client := openai.NewClient(option.WithBaseURL(original.URL+"/original"), option.WithAPIKey("synthetic-key"), option.WithOrganization("synthetic-org"), option.WithProject("synthetic-project"), option.WithHeader("X-Custom", "client"), option.WithQueryAdd("constructor", "once"))
	client.Realtime.Options = append(client.Realtime.Options,
		option.WithBaseURL(server.URL+"/v1?base=retained&model=old"), option.WithHTTPClient(doer),
		option.WithAPIKey("synthetic-service-key"), option.WithHeader("X-Service", "service"), option.WithHeader("X-Custom", "service"))
	opening, cancel := context.WithCancel(t.Context())
	defer cancel()
	connection, err := client.Realtime.ConnectTranslation(opening, "gpt-realtime-translate", realtime.TranslationConnectionOptions{}, option.WithHeader("X-Custom", "connection"), option.WithQueryAdd("model", "ignored"), requestconfig.RequestOptionFunc(func(cfg *requestconfig.RequestConfig) error {
		cfg.Request = cfg.Request.WithContext(context.WithValue(cfg.Request.Context(), contextKey{}, "preserved"))
		return nil
	}))
	if err != nil {
		t.Fatalf("ConnectTranslation(custom TLS) = %v, want success", err)
	}
	defer connection.Abort()
	cancel() // Opening context cancellation must not cancel the established socket.
	ctx, stop := context.WithTimeout(t.Context(), 5*time.Second)
	defer stop()
	if event, err := connection.Recv(ctx); err != nil || event.Type != "session.created" {
		t.Errorf("Recv(after opening cancellation) = %s, %v, want session.created", event.RawJSON(), err)
	}
	if err := connection.Close(); err != nil {
		t.Errorf("Close(custom TLS) = %v, want success", err)
	}
	translationResult(t, ctx, result)
	if originalRequests.Load() != 0 {
		t.Errorf("requests to original endpoint = %d, want 0", originalRequests.Load())
	}
	if upgrades.Load() != 1 || transports.Load() != 1 || doer.calls.Load() != 0 {
		t.Errorf("opening upgrade/transport/custom Do calls = %d/%d/%d, want 1/1/0", upgrades.Load(), transports.Load(), doer.calls.Load())
	}
}

type translationUnavailableDoer struct{ calls atomic.Int32 }

func (c *translationUnavailableDoer) Do(*http.Request) (*http.Response, error) {
	c.calls.Add(1)
	return nil, errors.New("unexpected custom HTTP Do call")
}

func TestTranslationRejectsUnsupportedClientAndProvider(t *testing.T) {
	doer := &translationUnavailableDoer{}
	client := openai.NewClient(option.WithBaseURL("https://translation.invalid/v1"), option.WithAPIKey("synthetic-key"), option.WithHTTPClient(doer))
	if _, err := client.Realtime.ConnectTranslation(t.Context(), "gpt-realtime-translate", realtime.TranslationConnectionOptions{}); err == nil || !strings.Contains(err.Error(), "WebSocketHTTPClient") {
		t.Errorf("ConnectTranslation(unsupported HTTP client) = %v, want capability rejection", err)
	}
	if doer.calls.Load() != 0 {
		t.Errorf("unsupported HTTP Do calls = %d, want 0", doer.calls.Load())
	}
	var credentialWork atomic.Int32
	client = openai.NewClient(option.WithBaseURL("https://translation.invalid/v1"), option.WithAPIKey("synthetic-key"), requestconfig.WithEndpointProvider("azure"), requestconfig.WithRequestFinalizer(func(*requestconfig.RequestConfig) error {
		credentialWork.Add(1)
		return errors.New("unexpected credential finalizer")
	}))
	if _, err := client.Realtime.ConnectTranslation(t.Context(), "gpt-realtime-translate", realtime.TranslationConnectionOptions{}); err == nil || !strings.Contains(err.Error(), "provider does not support Translation") {
		t.Errorf("ConnectTranslation(unsupported provider) = %v, want provider rejection", err)
	}
	if credentialWork.Load() != 0 {
		t.Errorf("unsupported-provider credential finalizers = %d, want 0", credentialWork.Load())
	}
}

func TestTranslationRejectsCrossOriginRedirect(t *testing.T) {
	var leaked atomic.Int32
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { leaked.Add(1) }))
	defer target.Close()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, target.URL+"/v1/realtime/translations", http.StatusTemporaryRedirect)
	}))
	defer server.Close()
	client := openai.NewClient(option.WithBaseURL(server.URL+"/v1"), option.WithAPIKey("synthetic-key"))
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	if connection, err := client.Realtime.ConnectTranslation(ctx, "gpt-realtime-translate", realtime.TranslationConnectionOptions{}); err == nil {
		connection.Abort()
		t.Error("ConnectTranslation(cross-origin redirect) succeeded, want rejection")
	}
	if leaked.Load() != 0 {
		t.Errorf("cross-origin requests = %d, want 0", leaked.Load())
	}
}

func TestTranslationFinishOverflowRetainsAcceptedEvents(t *testing.T) {
	callbackStarted, peerClosed := make(chan struct{}), make(chan struct{})
	connection, ctx, result := translationPeer(t, realtime.TranslationConnectionOptions{MaxBufferedEvents: 1}, func(ctx context.Context, socket *wire.Conn) error {
		defer close(peerClosed)
		if err := translationCommand(ctx, socket, "session.close"); err != nil {
			return err
		}
		if err := socket.Write(ctx, wire.MessageText, []byte(`{"type":"future.first"}`)); err != nil {
			return err
		}
		select {
		case <-callbackStarted:
		case <-ctx.Done():
			return ctx.Err()
		}
		for _, frame := range []string{`{"type":"future.second"}`, `{"type":"future.overflow"}`} {
			if err := socket.Write(ctx, wire.MessageText, []byte(frame)); err != nil {
				return err
			}
		}
		if _, _, err := socket.Read(ctx); err == nil {
			return errors.New("expected overflow to abort transport")
		}
		return nil
	})
	var got []string
	err := connection.Finish(ctx, func(event realtime.RealtimeTranslationServerEventUnion) error {
		got = append(got, event.Type)
		if len(got) == 1 {
			close(callbackStarted)
			select {
			case <-peerClosed:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		return nil
	})
	if !errors.Is(err, websocket.ErrBufferLimit) {
		t.Errorf("Finish(overflow) = %v, want buffer-limit error", err)
	}
	if want := []string{"future.first", "future.second"}; !reflect.DeepEqual(got, want) {
		t.Errorf("Finish(overflow) events = %v, want %v", got, want)
	}
	translationResult(t, ctx, result)
}

func TestTranslationFinishCanceledWhileWaitingForSend(t *testing.T) {
	connection, ctx, result := translationPeer(t, realtime.TranslationConnectionOptions{}, func(ctx context.Context, socket *wire.Conn) error {
		if _, _, err := socket.Read(ctx); err == nil {
			return errors.New("received a command after canceled finish")
		}
		return nil
	})
	started, release := make(chan struct{}), make(chan struct{})
	sent := make(chan error, 1)
	command := param.Override[realtime.RealtimeTranslationClientEventUnionParam](translationBlockingCommand{started: started, release: release, done: ctx.Done()})
	go func() { sent <- connection.Send(ctx, command) }()
	select {
	case <-started:
	case <-ctx.Done():
		t.Fatal("Send never began serialization")
	}
	finishCtx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()
	err := connection.Finish(finishCtx, func(realtime.RealtimeTranslationServerEventUnion) error { return nil })
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("Finish(waiting for send) = %v, want deadline", err)
	}
	close(release)
	if err := <-sent; !errors.Is(err, websocket.ErrClosed) {
		t.Errorf("Send(after canceled finish) = %v, want closed transport", err)
	}
	translationResult(t, ctx, result)
}

func TestTranslationUnsentCloseFailureIsNotRetried(t *testing.T) {
	connection, ctx, result := translationPeer(t, realtime.TranslationConnectionOptions{MaxMessageBytes: 1}, func(ctx context.Context, socket *wire.Conn) error {
		if _, _, err := socket.Read(ctx); err == nil {
			return errors.New("oversized close command was transmitted")
		}
		return nil
	})
	err := connection.Finish(ctx, func(realtime.RealtimeTranslationServerEventUnion) error { return nil })
	if !errors.Is(err, websocket.ErrMessageLimit) {
		t.Errorf("Finish(oversized close) = %v, want message-limit error", err)
	}
	var delivery *websocket.DeliveryError
	if !errors.As(err, &delivery) || delivery.MayHaveBeenSent {
		t.Errorf("Finish(oversized close) = %v, want definitely-unsent delivery error", err)
	}
	if again := connection.Finish(ctx, func(realtime.RealtimeTranslationServerEventUnion) error { return nil }); !errors.Is(again, err) {
		t.Errorf("Finish(repeated unsent close) = %v, want original %v", again, err)
	}
	translationResult(t, ctx, result)
}

func TestTranslationFinishAfterUnconsumedRemoteTerminal(t *testing.T) {
	connection, ctx, result := translationPeer(t, realtime.TranslationConnectionOptions{}, func(ctx context.Context, socket *wire.Conn) error {
		for _, frame := range []string{
			`{"type":"session.output_transcript.delta","delta":"Hola"}`,
			`{"type":"session.closed","event_id":"remote-close"}`,
		} {
			if err := socket.Write(ctx, wire.MessageText, []byte(frame)); err != nil {
				return err
			}
		}
		// The client's owned reader must consume both preceding frames to answer
		// this close handshake. No application Recv has claimed them.
		return socket.Close(wire.StatusNormalClosure, "")
	})
	translationResult(t, ctx, result)
	var got []string
	err := connection.Finish(ctx, func(event realtime.RealtimeTranslationServerEventUnion) error {
		got = append(got, event.Type)
		return nil
	})
	if err != nil {
		t.Errorf("Finish(after remote terminal and transport close) = %v, want success", err)
	}
	if want := []string{"session.output_transcript.delta", "session.closed"}; !reflect.DeepEqual(got, want) {
		t.Errorf("Finish(remote output) = %v, want %v", got, want)
	}
}

type translationHeldClose struct {
	started chan<- struct{}
	release <-chan struct{}
	done    <-chan struct{}
}

func (c translationHeldClose) MarshalJSON() ([]byte, error) {
	close(c.started)
	select {
	case <-c.release:
		return []byte(`{"type":"session.close"}`), nil
	case <-c.done:
		return nil, context.Canceled
	}
}

func TestTranslationCloseWaitPreservesRemoteCompletion(t *testing.T) {
	connection, ctx, result := translationPeer(t, realtime.TranslationConnectionOptions{}, func(ctx context.Context, socket *wire.Conn) error {
		if err := socket.Write(ctx, wire.MessageText, []byte(`{"type":"session.closed"}`)); err != nil {
			return err
		}
		if err := translationCommand(ctx, socket, "session.input_audio_buffer.append"); err != nil {
			return err
		}
		return translationClosed(ctx, socket)
	})
	started, release := make(chan struct{}), make(chan struct{})
	sent, closed := make(chan error, 1), make(chan error, 1)
	command := param.Override[realtime.RealtimeTranslationClientEventUnionParam](translationBlockingCommand{started: started, release: release, done: ctx.Done()})
	go func() { sent <- connection.Send(ctx, command) }()
	<-started
	go func() { closed <- connection.SendRaw(ctx, json.RawMessage(`{"type":"session.close"}`)) }()
	for {
		err := connection.SendRaw(ctx, json.RawMessage(`{"type":""}`))
		if errors.Is(err, realtime.ErrTranslationInputClosed) {
			break
		}
		if ctx.Err() != nil {
			t.Fatal("close never stopped input")
		}
		runtime.Gosched()
	}
	if event, err := connection.Recv(ctx); err != nil || event.Type != "session.closed" {
		t.Fatalf("Recv = %s, %v, want session.closed", event.Type, err)
	}
	close(release)
	if err := <-sent; err != nil {
		t.Errorf("admitted send = %v, want success", err)
	}
	if err := <-closed; err != nil {
		t.Errorf("close after remote completion = %v, want success", err)
	}
	if err := connection.Finish(ctx, func(realtime.RealtimeTranslationServerEventUnion) error {
		return errors.New("terminal event delivered twice")
	}); err != nil {
		t.Errorf("Finish = %v, want successful cleanup", err)
	}
	translationResult(t, ctx, result)
}

func TestTranslationSharedCloseDeliveryCancellation(t *testing.T) {
	connection, ctx, result := translationPeer(t, realtime.TranslationConnectionOptions{}, func(ctx context.Context, socket *wire.Conn) error {
		for _, kind := range []string{"session.input_audio_buffer.append", "session.close"} {
			if err := translationCommand(ctx, socket, kind); err != nil {
				return err
			}
		}
		if err := socket.Write(ctx, wire.MessageText, []byte(`{"type":"session.closed"}`)); err != nil {
			return err
		}
		return translationClosed(ctx, socket)
	})
	appendStarted, appendRelease := make(chan struct{}), make(chan struct{})
	appendResult := make(chan error, 1)
	appendCommand := param.Override[realtime.RealtimeTranslationClientEventUnionParam](translationBlockingCommand{started: appendStarted, release: appendRelease, done: ctx.Done()})
	go func() { appendResult <- connection.Send(ctx, appendCommand) }()
	<-appendStarted

	// Admit the follower before the first close changes input state. Hold its
	// serialization so the first close owns the sole protocol write.
	closeStarted, closeRelease := make(chan struct{}), make(chan struct{})
	follower, first := make(chan error, 1), make(chan error, 1)
	followerCtx, cancelFollower := context.WithCancel(ctx)
	defer cancelFollower()
	held := param.Override[realtime.RealtimeTranslationClientEventUnionParam](translationHeldClose{started: closeStarted, release: closeRelease, done: ctx.Done()})
	go func() { follower <- connection.Send(followerCtx, held) }()
	<-closeStarted
	go func() { first <- connection.SendRaw(ctx, json.RawMessage(`{"type":"session.close"}`)) }()
	for {
		err := connection.SendRaw(ctx, json.RawMessage(`{"type":""}`))
		if errors.Is(err, realtime.ErrTranslationInputClosed) {
			break
		}
		if ctx.Err() != nil {
			t.Fatal("first close never stopped input")
		}
		runtime.Gosched()
	}
	close(closeRelease)
	cancelFollower()
	err := <-follower
	var delivery *websocket.DeliveryError
	if !errors.Is(err, context.Canceled) || !errors.As(err, &delivery) || !delivery.MayHaveBeenSent {
		t.Errorf("Send(canceled follower) = %v, want delivery-uncertain cancellation", err)
	}
	close(appendRelease)
	if err := <-appendResult; err != nil {
		t.Errorf("admitted append = %v", err)
	}
	if err := <-first; err != nil {
		t.Errorf("first close = %v", err)
	}
	if err := connection.Finish(ctx, func(realtime.RealtimeTranslationServerEventUnion) error { return nil }); err != nil {
		t.Errorf("Finish(shared close) = %v", err)
	}
	translationResult(t, ctx, result)
}

func TestTranslationCloseCanceledBeforeWriteIsUnsent(t *testing.T) {
	connection, ctx, result := translationPeer(t, realtime.TranslationConnectionOptions{}, func(ctx context.Context, socket *wire.Conn) error {
		_, _, err := socket.Read(ctx)
		if err == nil {
			return errors.New("unexpected write after interrupted close")
		}
		return nil
	})
	started, release := make(chan struct{}), make(chan struct{})
	sent := make(chan error, 1)
	command := param.Override[realtime.RealtimeTranslationClientEventUnionParam](translationBlockingCommand{started: started, release: release, done: ctx.Done()})
	go func() { sent <- connection.Send(ctx, command) }()
	<-started
	closeCtx, cancel := context.WithTimeout(ctx, 20*time.Millisecond)
	defer cancel()
	err := connection.SendRaw(closeCtx, json.RawMessage(`{"type":"session.close"}`))
	var delivery *websocket.DeliveryError
	if !errors.Is(err, context.DeadlineExceeded) || !errors.As(err, &delivery) || delivery.MayHaveBeenSent {
		t.Errorf("Send(close canceled before write) = %v, want definitely-unsent deadline", err)
	}
	// Finish must clean up without retrying the rejected close.
	if again := connection.Finish(ctx, func(realtime.RealtimeTranslationServerEventUnion) error { return nil }); !errors.Is(again, err) {
		t.Errorf("Finish(failed close) = %v, want original %v", again, err)
	}
	close(release)
	if err := <-sent; !errors.Is(err, websocket.ErrClosed) {
		t.Errorf("Send(after failed close) = %v, want closed", err)
	}
	translationResult(t, ctx, result)
}

func TestTranslationFinishDrainsAfterUncertainCloseCancellation(t *testing.T) {
	queued := make(chan struct{})
	release := make(chan struct{})
	connection, ctx, result := translationPeer(t, realtime.TranslationConnectionOptions{}, func(ctx context.Context, socket *wire.Conn) error {
		for _, frame := range []string{
			`{"type":"session.output_transcript.delta","delta":"Hola"}`,
			`{"type":"session.closed","event_id":"server-close"}`,
		} {
			if err := socket.Write(ctx, wire.MessageText, []byte(frame)); err != nil {
				return err
			}
		}
		close(queued)
		// Do not read the caller's closing frame. A sufficiently large future
		// field forces real socket backpressure until the caller cancels.
		select {
		case <-release:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	})
	<-queued
	command := json.RawMessage(`{"type":"session.close","future":"` + strings.Repeat("x", 16<<20) + `"}`)
	// Allow the large JSON command to be validated even under the race
	// detector; the peer never reads, so the socket write cannot complete.
	closeCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	err := connection.SendRaw(closeCtx, command)
	var delivery *websocket.DeliveryError
	if !errors.Is(err, context.DeadlineExceeded) || !errors.As(err, &delivery) || !delivery.MayHaveBeenSent {
		t.Fatalf("Send(close with socket backpressure) = %v, want delivery-uncertain deadline", err)
	}
	var got []string
	err = connection.Finish(ctx, func(event realtime.RealtimeTranslationServerEventUnion) error {
		got = append(got, event.Type)
		return nil
	})
	if err != nil {
		t.Errorf("Finish(with a fresh context) = %v, want accepted final output", err)
	}
	if want := []string{"session.output_transcript.delta", "session.closed"}; !reflect.DeepEqual(got, want) {
		t.Errorf("Finish(after canceled close) events = %v, want %v", got, want)
	}
	close(release)
	translationResult(t, ctx, result)
}

package websockettest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync/atomic"
	"testing"
	"time"

	wire "github.com/coder/websocket"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/websocket"
	"github.com/openai/openai-go/v3/responses"
)

func TestResponsesWS21ServiceErrorCodes(t *testing.T) {
	for _, tc := range []struct {
		code     string
		status   int64
		streamID string
	}{
		{code: "previous_response_not_found", status: 400, streamID: "ws21-1"},
		{code: "invalid_stream_id", status: 400, streamID: "ws21-2"},
		{code: "websocket_stream_limit_reached", status: 429, streamID: "ws21-3"},
		{code: "websocket_connection_limit_reached", status: 429, streamID: "ws21-4"},
	} {
		t.Run(tc.code, func(t *testing.T) {
			turn := fixtureScenario(t, "nested_error_then_completed").Turns[0]
			var request, envelope map[string]any
			if err := json.Unmarshal(turn.Request, &request); err != nil {
				t.Fatal(err)
			}
			if err := json.Unmarshal(turn.Frames[0], &envelope); err != nil {
				t.Fatal(err)
			}
			const errorType = "invalid_request_error"
			const message = "Synthetic invalid input"
			const param = "input"
			details, ok := envelope["error"].(map[string]any)
			if !ok || details["type"] != errorType || details["message"] != message || details["param"] != param {
				t.Fatal("fixture lost its synthetic nested error details")
			}
			// Status and stream ID are synthetic inputs, not service policy assertions.
			request["stream_id"] = tc.streamID
			envelope["stream_id"] = tc.streamID
			envelope["status"] = tc.status
			details["code"] = tc.code
			command, err := json.Marshal(request)
			if err != nil {
				t.Fatal(err)
			}
			frame, err := json.Marshal(envelope)
			if err != nil {
				t.Fatal(err)
			}

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			t.Cleanup(cancel)
			type peerResult struct {
				commands, deliveries int
				requestMatched       bool
				closeStatus          wire.StatusCode
				closeNowErr, err     error
			}
			// Exactly one handler owns this buffered completion slot.
			done := make(chan peerResult, 1)
			var attempts, upgrades atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if attempts.Add(1) != 1 {
					http.Error(w, "unexpected second upgrade", http.StatusBadRequest)
					return
				}
				var result peerResult
				defer func() { done <- result }()
				peerCtx, peerCancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer peerCancel()
				socket, acceptErr := wire.Accept(w, r.WithContext(peerCtx), nil)
				if acceptErr != nil {
					result.err = fmt.Errorf("upgrade: %w", acceptErr)
					return
				}
				upgrades.Add(1)
				defer func() { result.closeNowErr = socket.CloseNow() }()
				kind, data, readErr := socket.Read(peerCtx)
				if readErr != nil {
					result.err = fmt.Errorf("expected Create: %w", readErr)
					return
				}
				result.commands++
				var got map[string]any
				if decodeErr := json.Unmarshal(data, &got); decodeErr != nil {
					result.err = fmt.Errorf("decode Create: %w", decodeErr)
					return
				}
				result.requestMatched = kind == wire.MessageText && got["type"] == "response.create" && reflect.DeepEqual(got, request)
				if !result.requestMatched {
					result.err = errors.New("Create did not match the synthetic request")
					return
				}
				// Script two identical errors so Recv and FinalResponse each consume
				// one frame after the same Create. This does not model service duplication.
				for range 2 {
					if writeErr := socket.Write(peerCtx, wire.MessageText, frame); writeErr != nil {
						result.err = fmt.Errorf("scripted error delivery: %w", writeErr)
						return
					}
					result.deliveries++
				}
				_, _, readErr = socket.Read(peerCtx)
				result.closeStatus = wire.CloseStatus(readErr)
				if readErr == nil {
					result.commands++
					result.err = errors.New("unexpected additional command")
				} else if result.closeStatus != wire.StatusNormalClosure {
					result.err = fmt.Errorf("expected client close: %w", readErr)
				}
			}))
			var connection *responses.ResponseConnection
			// Registered before Connect: failures still close, abort, await the
			// peer's CloseNow result, and stop the server.
			t.Cleanup(func() {
				cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cleanupCancel()
				var closeErr error
				if connection != nil {
					closeErr = connection.Close()
					connection.Abort()
				}
				if closeErr != nil {
					t.Errorf("connection Close: %v", closeErr)
				}
				select {
				case result := <-done:
					t.Logf("WS21 peer upgrades=%d attempts=%d commands=%d request_matched=%t deliveries=%d close_status=%d close_now_error=%v peer_error=%v",
						upgrades.Load(), attempts.Load(), result.commands, result.requestMatched, result.deliveries, result.closeStatus, result.closeNowErr, result.err)
					if upgrades.Load() != 1 || attempts.Load() != 1 || result.commands != 1 || !result.requestMatched || result.deliveries != 2 || result.err != nil {
						t.Error("peer did not complete the one-Create/two-error script cleanly")
					}
					if result.closeNowErr != nil && !errors.Is(result.closeNowErr, net.ErrClosed) {
						t.Errorf("peer CloseNow: %v", result.closeNowErr)
					}
				case <-cleanupCtx.Done():
					t.Errorf("peer completion: %v", cleanupCtx.Err())
				}
				server.Close()
				t.Logf("WS21 cleanup connection_close_error=%v abort_called=%t server_closed=true assertion_failed=%t", closeErr, connection != nil, t.Failed())
			})
			client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("fixture-key"))
			// Keep the local closing handshake well within the five-second waits.
			connection, err = client.Responses.Connect(ctx, responses.ResponseConnectionOptions{CloseTimeout: 100 * time.Millisecond})
			if err != nil {
				t.Fatalf("Connect: %v", err)
			}
			// No lane is attached: existing routing delivers this supplied ID to Recv.
			sendFixture(t, ctx, connection.Create, command)

			checkEvent := func(path string, event responses.ResponsesServerEventUnion) {
				t.Helper()
				got, ok := event.AsAny().(responses.ResponsesServerEventResponseWsError)
				if !ok {
					t.Errorf("%s AsAny = %T, want ResponsesServerEventResponseWsError", path, event.AsAny())
					return
				}
				valid := event.JSON.Type.Valid() && event.JSON.OfResponsesServerEventResponseWsError.Valid() &&
					got.JSON.Type.Valid() && got.JSON.Error.Valid() && got.JSON.Status.Valid() && got.JSON.StreamID.Valid() &&
					got.Error.JSON.Code.Valid() && got.Error.JSON.Type.Valid() && got.Error.JSON.Message.Valid() && got.Error.JSON.Param.Valid()
				t.Logf("WS21 path=%s code=%s status=%d stream_id=%s nested_type=%s message=%q param=%s supplied_fields_valid=%t",
					path, got.Error.Code, got.Status, got.StreamID, got.Error.Type, got.Error.Message, got.Error.Param, valid)
				if event.Type != "error" || got.Type != "error" || got.Error.Code != tc.code || got.Status != tc.status || got.StreamID != tc.streamID ||
					got.Error.Type != errorType || got.Error.Message != message || got.Error.Param != param || !valid {
					t.Errorf("%s typed event did not preserve the supplied %s fields and presence metadata", path, tc.code)
				}
			}
			event, receiveErr := connection.Recv(ctx)
			t.Logf("WS21 Recv receive_error=%v", receiveErr)
			if receiveErr != nil {
				t.Errorf("Recv error = %v, want nil", receiveErr)
			} else {
				checkEvent("Recv", event)
			}

			response, finalErr := connection.FinalResponse(ctx)
			var protocolErr *responses.ResponseProtocolError
			isAPIError := errors.As(finalErr, &protocolErr)
			t.Logf("WS21 FinalResponse response_nil=%t protocol_error=%t", response == nil, isAPIError)
			if response != nil || !isAPIError {
				t.Errorf("FinalResponse = response-nil:%t, error:%T; want nil response and *ResponseProtocolError", response == nil, finalErr)
			}
			if isAPIError {
				checkEvent("FinalResponse", protocolErr.Event)
			}
			// Classify the actual returned API error, without testing additional faults.
			wrongSentinel := false
			for _, other := range []error{context.Canceled, context.DeadlineExceeded, io.EOF, io.ErrUnexpectedEOF, net.ErrClosed, websocket.ErrClosed} {
				if errors.Is(finalErr, other) {
					wrongSentinel = true
					t.Errorf("FinalResponse API error unexpectedly matches %v", other)
				}
			}
			var syntaxErr *json.SyntaxError
			var typeErr *json.UnmarshalTypeError
			var closeError wire.CloseError
			var transportErr net.Error
			var deliveryErr *websocket.DeliveryError
			wrongCategory := errors.As(finalErr, &syntaxErr) || errors.As(finalErr, &typeErr) ||
				errors.As(finalErr, &closeError) || errors.As(finalErr, &transportErr) || errors.As(finalErr, &deliveryErr)
			if wrongCategory {
				t.Errorf("FinalResponse API error matches a decode, close, transport, or delivery error: %T", finalErr)
			}
			t.Logf("WS21 negative_classification wrong_sentinel=%t wrong_typed_category=%t", wrongSentinel, wrongCategory)
		})
	}
}

package websockettest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync/atomic"
	"testing"
	"time"

	wire "github.com/coder/websocket"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

func fixtureScenario(t *testing.T, id string) scenario {
	t.Helper()
	for _, item := range loadScenarios(t) {
		if item.ID == id {
			return item
		}
	}
	t.Fatalf("missing scenario %s", id)
	return scenario{}
}

func TestResponsesFinalSnapshots(t *testing.T) {
	for _, item := range []struct{ id, status string }{{"completed_then_completed", "completed"}, {"failed_then_completed", "failed"}, {"incomplete_then_completed", "incomplete"}} {
		t.Run(item.status, func(t *testing.T) {
			turn := fixtureScenario(t, item.id).Turns[0]
			done := make(chan struct{})
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				defer close(done)
				socket, err := wire.Accept(w, r, nil)
				if err != nil {
					t.Errorf("upgrade: %v", err)
					return
				}
				defer func() { _ = socket.CloseNow() }()
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				if _, _, err = socket.Read(ctx); err != nil {
					t.Errorf("command: %v", err)
					return
				}
				for _, frame := range turn.Frames {
					if err = socket.Write(ctx, wire.MessageText, frame); err != nil {
						t.Errorf("event: %v", err)
						return
					}
				}
				_, _, _ = socket.Read(ctx)
			}))
			defer server.Close()
			client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("fixture-key"))
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			connection, err := client.Responses.Connect(ctx, responses.ResponseConnectionOptions{})
			if err != nil {
				t.Fatal(err)
			}
			defer connection.Abort()
			sendFixture(t, ctx, connection.Create, turn.Request)
			result, err := connection.FinalResponse(ctx)
			if err != nil {
				t.Fatal(err)
			}
			if result.ID != "resp_first" || string(result.Status) != item.status {
				t.Fatalf("final result = %s/%s", result.ID, result.Status)
			}
			if item.status == "completed" && result.OutputText() != "Synthetic text" {
				t.Fatalf("final output text = %q", result.OutputText())
			}
			_ = connection.Close()
			<-done
		})
	}
}

func TestResponsesFinalSnapshotRetainsToolsAndMultipartOutput(t *testing.T) {
	turn := fixtureScenario(t, "completed_then_completed").Turns[0]
	var terminal map[string]any
	if err := json.Unmarshal(turn.Frames[len(turn.Frames)-1], &terminal); err != nil {
		t.Fatal(err)
	}
	response := terminal["response"].(map[string]any)
	response["output"] = []any{
		map[string]any{"type": "function_call", "id": "fc_synthetic", "call_id": "call_synthetic", "name": "lookup", "arguments": "{}", "status": "completed"},
		map[string]any{"type": "message", "id": "msg_synthetic", "role": "assistant", "status": "completed", "content": []any{
			map[string]any{"type": "output_text", "text": "first", "annotations": []any{}},
			map[string]any{"type": "output_text", "text": "second", "annotations": []any{}},
		}},
	}
	payload, err := json.Marshal(terminal)
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		socket, acceptErr := wire.Accept(w, r, nil)
		if acceptErr != nil {
			t.Errorf("upgrade: %v", acceptErr)
			return
		}
		defer func() { _ = socket.CloseNow() }()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if _, _, readErr := socket.Read(ctx); readErr != nil {
			t.Errorf("command: %v", readErr)
			return
		}
		for _, frame := range [][]byte{[]byte(`{"type":"response.future_event","extra":"preserved"}`), payload} {
			if writeErr := socket.Write(ctx, wire.MessageText, frame); writeErr != nil {
				t.Errorf("event: %v", writeErr)
				return
			}
		}
		_, _, _ = socket.Read(ctx)
	}))
	defer server.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("fixture-key"))
	connection, err := client.Responses.Connect(ctx, responses.ResponseConnectionOptions{})
	if err != nil {
		t.Fatal(err)
	}
	defer connection.Abort()
	sendFixture(t, ctx, connection.Create, turn.Request)
	result, err := connection.FinalResponse(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Output) != 2 || result.OutputText() != "firstsecond" || result.Output[0].AsFunctionCall().Name != "lookup" {
		t.Fatal("final snapshot lost tool or multi-part output")
	}
	_ = connection.Close()
}

func withStreamID(t *testing.T, data json.RawMessage, streamID string) json.RawMessage {
	t.Helper()
	var object map[string]any
	if err := json.Unmarshal(data, &object); err != nil {
		t.Fatal(err)
	}
	object["stream_id"] = streamID
	result, err := json.Marshal(object)
	if err != nil {
		t.Fatal(err)
	}
	return result
}

func TestResponsesInterleavedLanes(t *testing.T) {
	for _, failA := range []bool{false, true} {
		t.Run(fmt.Sprintf("lane_a_fails_%t", failA), func(t *testing.T) {
			turn := fixtureScenario(t, "completed_then_completed").Turns[0]
			framesA, framesB := make([]json.RawMessage, len(turn.Frames)), make([]json.RawMessage, len(turn.Frames))
			for i, frame := range turn.Frames {
				framesA[i] = withStreamID(t, frame, "a")
				framesB[i] = withStreamID(t, frame, "b")
			}
			if failA {
				framesA = []json.RawMessage{json.RawMessage(`{"type":"error","stream_id":"a","error":{"type":"invalid_request_error","code":"synthetic","message":"synthetic lane failure"}}`)}
			}
			done := make(chan struct{})
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				defer close(done)
				socket, err := wire.Accept(w, r, nil)
				if err != nil {
					t.Errorf("upgrade: %v", err)
					return
				}
				defer func() { _ = socket.CloseNow() }()
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				for range 2 {
					if _, _, err = socket.Read(ctx); err != nil {
						t.Errorf("command: %v", err)
						return
					}
				}
				for i := range framesB {
					frames := []json.RawMessage{framesB[i]}
					if i < len(framesA) {
						frames = append([]json.RawMessage{framesA[i]}, frames...)
					}
					for _, frame := range frames {
						if err = socket.Write(ctx, wire.MessageText, frame); err != nil {
							t.Errorf("event: %v", err)
							return
						}
					}
				}
				_, _, _ = socket.Read(ctx)
			}))
			defer server.Close()
			client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("fixture-key"))
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			connection, err := client.Responses.Connect(ctx, responses.ResponseConnectionOptions{})
			if err != nil {
				t.Fatal(err)
			}
			defer connection.Abort()
			laneA, err := connection.Lane("a")
			if err != nil {
				t.Fatal(err)
			}
			laneB, err := connection.Lane("b")
			if err != nil {
				t.Fatal(err)
			}
			defer laneB.Close()
			sendFixture(t, ctx, connection.Create, withStreamID(t, turn.Request, "a"))
			sendFixture(t, ctx, connection.Create, withStreamID(t, turn.Request, "b"))
			resultA, resultErr := laneA.FinalResponse(ctx)
			if failA {
				var protocolError *responses.ResponseProtocolError
				if !errors.As(resultErr, &protocolError) {
					t.Fatalf("lane a error = %v, want typed protocol error", resultErr)
				}
			} else if resultErr != nil || resultA.OutputText() != "Synthetic text" {
				t.Fatalf("lane a final = %v,%v", resultA, resultErr)
			}
			laneA.Close()
			if result, laneBErr := laneB.FinalResponse(ctx); laneBErr != nil || result.OutputText() != "Synthetic text" {
				t.Fatalf("lane b final = %v,%v", result, laneBErr)
			}
			_ = connection.Close()
			<-done
		})
	}
}

func TestResponsesRecoveryRefreshesCredentialsWithoutReplay(t *testing.T) {
	turn := fixtureScenario(t, "completed_then_completed").Turns[0]
	var attempts, credentials, commands atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempt := attempts.Add(1)
		if r.Header.Get("Authorization") != fmt.Sprintf("Bearer fixture-%d", attempt) {
			t.Errorf("credential was not refreshed on attempt %d", attempt)
		}
		if attempt == 2 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		socket, err := wire.Accept(w, r, nil)
		if err != nil {
			t.Errorf("upgrade: %v", err)
			return
		}
		defer func() { _ = socket.CloseNow() }()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if _, _, err = socket.Read(ctx); err != nil {
			t.Errorf("command: %v", err)
			return
		}
		commands.Add(1)
		if attempt == 1 {
			_ = socket.Close(wire.StatusInternalError, "")
			return
		}
		for _, frame := range turn.Frames {
			if err = socket.Write(ctx, wire.MessageText, frame); err != nil {
				t.Errorf("event: %v", err)
				return
			}
		}
		if _, _, err = socket.Read(ctx); err == nil {
			t.Error("unexpected replayed command")
		}
	}))
	defer server.Close()
	client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("fixture-key"), option.WithMaxRetries(9), option.WithMiddleware(func(r *http.Request, next option.MiddlewareNext) (*http.Response, error) {
		r.Header.Set("Authorization", fmt.Sprintf("Bearer fixture-%d", credentials.Add(1)))
		return next(r)
	}))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	connection, err := client.Responses.Connect(ctx, responses.ResponseConnectionOptions{})
	if err != nil {
		t.Fatal(err)
	}
	defer connection.Abort()
	sendFixture(t, ctx, connection.Create, turn.Request)
	if _, err = connection.Recv(ctx); err == nil {
		t.Fatal("disconnect was not surfaced")
	}
	var restored bool
	replacement, err := connection.Recover(ctx, responses.ResponseRecoveryOptions{MaxAttempts: 2, Restore: func(ctx context.Context, next *responses.ResponseConnection) error {
		restored = true
		sendFixture(t, ctx, next.Create, turn.Request)
		response, restoreErr := next.FinalResponse(ctx)
		if restoreErr != nil {
			return restoreErr
		}
		if response.ID != "resp_first" {
			return errors.New("unexpected restored response")
		}
		return nil
	}})
	if err != nil {
		t.Fatal(err)
	}
	defer replacement.Abort()
	if !restored || attempts.Load() != 3 || credentials.Load() != 3 || commands.Load() != 2 {
		t.Fatalf("restored=%t attempts=%d credentials=%d commands=%d", restored, attempts.Load(), credentials.Load(), commands.Load())
	}
	_ = replacement.Close()
}

func TestResponsesCustomHeaderPrecedenceAndReconnect(t *testing.T) {
	headers := make(chan http.Header, 2)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headers <- r.Header.Clone()
		socket, err := wire.Accept(w, r, nil)
		if err != nil {
			t.Errorf("upgrade: %v", err)
			return
		}
		defer func() { _ = socket.CloseNow() }()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, _, _ = socket.Read(ctx)
	}))
	defer server.Close()
	var refresh atomic.Int32
	client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("fixture-key"), option.WithHeader("X-Application", "client-value"), option.WithHeader("X-Remove", "remove-me"), option.WithMiddleware(func(r *http.Request, next option.MiddlewareNext) (*http.Response, error) {
		r.Header.Set("X-Refresh", fmt.Sprint(refresh.Add(1)))
		return next(r)
	}))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	connection, err := client.Responses.Connect(ctx, responses.ResponseConnectionOptions{}, option.WithHeader("X-Application", "connection-value"), option.WithHeaderAdd("X-Multiple", "first"), option.WithHeaderAdd("X-Multiple", "second"), option.WithHeaderDel("X-Remove"))
	if err != nil {
		t.Fatal(err)
	}
	defer connection.Abort()
	check := func(expectedRefresh string) {
		t.Helper()
		select {
		case got := <-headers:
			if got.Get("X-Application") != "connection-value" || got.Get("Authorization") != "Bearer fixture-key" || got.Get("X-Remove") != "" || got.Get("X-Refresh") != expectedRefresh {
				t.Errorf("custom headers/auth/refresh were not preserved: application=%q removed=%q refresh=%q auth-present=%t", got.Get("X-Application"), got.Get("X-Remove"), got.Get("X-Refresh"), got.Get("Authorization") != "")
			}
			if values := got.Values("X-Multiple"); len(values) != 2 || values[0] != "first" || values[1] != "second" {
				t.Errorf("multi-value custom header = %q", values)
			}
		case <-ctx.Done():
			t.Fatal("no upgrade headers captured")
		}
	}
	check("1")
	replacement, err := connection.Reconnect(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer replacement.Abort()
	check("2")
	_ = replacement.Close()
}

func TestResponsesTypedCommandsPreserveWireShape(t *testing.T) {
	requests := []json.RawMessage{
		json.RawMessage(`{"type":"response.create","model":"gpt-4o-mini","input":"warmup","generate":false,"stream_id":"warmup"}`),
		json.RawMessage(`{"type":"response.steer","previous_response_id":"resp_parent","input":[{"type":"message","role":"user","content":[{"type":"input_text","text":"Prioritize testing."}]}]}`),
	}
	captured := make(chan []byte, len(requests))
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		socket, err := wire.Accept(w, r, nil)
		if err != nil {
			t.Errorf("upgrade: %v", err)
			return
		}
		defer func() { _ = socket.CloseNow() }()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		for range requests {
			_, data, readErr := socket.Read(ctx)
			if readErr != nil {
				t.Errorf("command: %v", readErr)
				return
			}
			captured <- data
		}
		_, _, _ = socket.Read(ctx)
	}))
	defer server.Close()
	client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("fixture-key"))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	connection, err := client.Responses.Connect(ctx, responses.ResponseConnectionOptions{})
	if err != nil {
		t.Fatal(err)
	}
	defer connection.Abort()
	for _, request := range requests {
		sendFixture(t, ctx, func(ctx context.Context, event responses.ResponsesClientEventUnionParam) error {
			if event.OfResponseCreate != nil {
				// Warmup is supplied through the existing extra-field API.
				event.OfResponseCreate.SetExtraFields(map[string]any{"generate": false})
			}
			return connection.Send(ctx, event)
		}, request)
		select {
		case data := <-captured:
			var want, got any
			if decodeErr := json.Unmarshal(request, &want); decodeErr != nil {
				t.Fatal(decodeErr)
			}
			if decodeErr := json.Unmarshal(data, &got); decodeErr != nil {
				t.Fatal(decodeErr)
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("typed command wire shape = %s, want %s", data, request)
			}
		case <-ctx.Done():
			t.Fatal("typed command was not received")
		}
	}
	_ = connection.Close()
}

package websockettest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	wire "github.com/coder/websocket"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

func TestResponsesDefaultErrorWithTwoActiveLanes(t *testing.T) {
	turn := fixtureScenario(t, "completed_then_completed").Turns[0]
	lanes := []struct {
		streamID, responseID, text string
		request                    json.RawMessage
		frames                     [3]json.RawMessage
	}{
		{streamID: "default-error-a", responseID: "resp_default_error_a", text: "lane A completed"},
		{streamID: "default-error-b", responseID: "resp_default_error_b", text: "lane B completed"},
	}
	for i := range lanes {
		item := &lanes[i]
		item.request = withStreamID(t, turn.Request, item.streamID)
		for j, eventType := range []string{"response.created", "response.in_progress", "response.completed"} {
			source := turn.Frames[0]
			if j == 2 {
				source = turn.Frames[len(turn.Frames)-1]
			}
			var frame map[string]any
			if err := json.Unmarshal(source, &frame); err != nil {
				t.Fatalf("decode %s fixture: %v", eventType, err)
			}
			frame["type"] = eventType
			frame["sequence_number"] = j
			response := frame["response"].(map[string]any)
			response["id"] = item.responseID
			if j == 2 {
				message := response["output"].([]any)[0].(map[string]any)
				message["id"] = "msg_" + item.streamID
				message["content"].([]any)[0].(map[string]any)["text"] = item.text
			}
			payload, err := json.Marshal(frame)
			if err != nil {
				t.Fatalf("encode %s fixture: %v", eventType, err)
			}
			item.frames[j] = withStreamID(t, payload, item.streamID)
		}
	}

	const errorFrame = `{"type":"error","sequence_number":42,"status":400,"error":{"type":"invalid_request_error","code":"invalid_value","param":"input","message":"Synthetic default connection error"}}`
	bothActive := make(chan struct{})
	defaultReceived := make(chan struct{})
	peerDone := make(chan struct{})
	peerResult := make(chan error, 1)
	peerCtx, cancelPeer := context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(cancelPeer)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer close(peerDone)
		peerResult <- func() (peerErr error) {
			socket, err := wire.Accept(w, r, nil)
			if err != nil {
				return fmt.Errorf("upgrade: %w", err)
			}
			// CloseNow also joins the peer's internal goroutines on every return path.
			defer func() {
				if err := socket.CloseNow(); err != nil && !errors.Is(err, net.ErrClosed) {
					peerErr = errors.Join(peerErr, fmt.Errorf("peer CloseNow: %w", err))
				}
			}()
			for _, item := range lanes {
				kind, request, err := socket.Read(peerCtx)
				if err != nil {
					return fmt.Errorf("read Create for %s: %w", item.streamID, err)
				}
				var got, want any
				if err := json.Unmarshal(request, &got); err != nil {
					return fmt.Errorf("decode Create for %s: %w", item.streamID, err)
				}
				if err := json.Unmarshal(item.request, &want); err != nil {
					return fmt.Errorf("decode expected Create: %w", err)
				}
				if kind != wire.MessageText || !reflect.DeepEqual(got, want) {
					return fmt.Errorf("Create for %s = %s, want %s", item.streamID, request, item.request)
				}
				t.Logf("peer: received Create %s", item.streamID)
				if err := socket.Write(peerCtx, wire.MessageText, item.frames[0]); err != nil {
					return fmt.Errorf("write created for %s: %w", item.streamID, err)
				}
				t.Logf("peer: sent created %s", item.responseID)
			}
			select {
			case <-bothActive:
				t.Log("peer: both distinct response starts consumed; both lanes remain active")
			case <-peerCtx.Done():
				return fmt.Errorf("wait for both active lanes: %w", peerCtx.Err())
			}
			if err := socket.Write(peerCtx, wire.MessageText, []byte(errorFrame)); err != nil {
				return fmt.Errorf("write default error: %w", err)
			}
			t.Logf("peer: sent exactly one unkeyed error %s", errorFrame)
			select {
			case <-defaultReceived:
				t.Log("peer: default error receipt confirmed")
			case <-peerCtx.Done():
				return fmt.Errorf("wait for default receipt: %w", peerCtx.Err())
			}
			// Interleave progress A, progress B, completed A, completed B.
			for j := 1; j < 3; j++ {
				for _, item := range lanes {
					if err := socket.Write(peerCtx, wire.MessageText, item.frames[j]); err != nil {
						return fmt.Errorf("write frame %d for %s: %w", j, item.streamID, err)
					}
					t.Logf("peer: sent frame %d for %s", j, item.responseID)
				}
			}
			// A third data message would violate the exactly-two-Create contract.
			if _, _, err := socket.Read(peerCtx); wire.CloseStatus(err) != wire.StatusNormalClosure {
				return fmt.Errorf("expected normal close after exactly two Creates, got %v", err)
			}
			t.Log("peer: observed normal client close after exactly two Creates")
			return nil
		}()
	}))
	t.Cleanup(func() {
		cancelPeer()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		select {
		case <-peerDone:
			if err := <-peerResult; err != nil {
				t.Errorf("peer failed: %v", err)
			}
			t.Log("cleanup: peer completed after CloseNow")
		case <-ctx.Done():
			t.Errorf("peer completion: %v", ctx.Err())
		}
		server.Close()
		t.Log("cleanup: httptest server closed")
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(cancel)
	client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("fixture-key"))
	// Keep assertion-failure cleanup short while allowing a loopback closing handshake.
	connection, err := client.Responses.Connect(ctx, responses.ResponseConnectionOptions{CloseTimeout: 100 * time.Millisecond})
	if err != nil {
		t.Fatalf("Connect: %v", err)
	}
	t.Cleanup(func() {
		connection.Abort()
		t.Log("cleanup: connection Abort fallback completed")
	})
	receivers := make([]*responses.ResponseLane, len(lanes))
	for i := range lanes {
		receivers[i], err = connection.Lane(lanes[i].streamID)
		if err != nil {
			t.Fatalf("attach lane %s: %v", lanes[i].streamID, err)
		}
		t.Cleanup(receivers[i].Close)
	}
	for i, item := range lanes {
		sendFixture(t, ctx, connection.Create, item.request)
		event, recvErr := receivers[i].Recv(ctx)
		if recvErr != nil {
			t.Fatalf("Recv created for %s: %v", item.streamID, recvErr)
		}
		created, ok := event.AsAny().(responses.ResponsesServerEventResponseWsCreated)
		if !ok {
			t.Fatalf("Recv created for %s = %T, want ResponsesServerEventResponseWsCreated", item.streamID, event.AsAny())
		}
		if created.Type != "response.created" || created.StreamID != item.streamID || created.Response.ID != item.responseID || created.Response.Status != "in_progress" {
			t.Fatalf("created for %s = %s/%s/%s/%s", item.streamID, created.Type, created.StreamID, created.Response.ID, created.Response.Status)
		}
		t.Logf("client: consumed created %s on %s", created.Response.ID, item.streamID)
	}
	close(bothActive)
	event, err := connection.Recv(ctx)
	if err != nil {
		t.Fatalf("default Recv: %v, want nil Go error", err)
	}
	apiError, ok := event.AsAny().(responses.ResponsesServerEventResponseWsError)
	if !ok {
		t.Fatalf("default Recv = %T, want ResponsesServerEventResponseWsError", event.AsAny())
	}
	if apiError.Type != "error" || apiError.Status != 400 || apiError.SequenceNumber != 42 ||
		apiError.Error.Type != "invalid_request_error" || apiError.Error.Code != "invalid_value" ||
		apiError.Error.Param != "input" || apiError.Error.Message != "Synthetic default connection error" {
		t.Fatalf("default error fields = %+v", apiError)
	}
	if !apiError.JSON.Type.Valid() || !apiError.JSON.Status.Valid() || !apiError.JSON.SequenceNumber.Valid() ||
		!apiError.JSON.Error.Valid() || !apiError.Error.JSON.Type.Valid() || !apiError.Error.JSON.Code.Valid() ||
		!apiError.Error.JSON.Param.Valid() || !apiError.Error.JSON.Message.Valid() {
		t.Fatal("default error lost valid presence metadata for supplied fields")
	}
	if apiError.StreamID != "" || apiError.JSON.StreamID.Valid() || apiError.JSON.StreamID.Raw() != "" {
		t.Fatalf("default error StreamID = %q, presence valid = %t, field raw = %q; want empty and absent", apiError.StreamID, apiError.JSON.StreamID.Valid(), apiError.JSON.StreamID.Raw())
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal([]byte(apiError.RawJSON()), &raw); err != nil {
		t.Fatalf("decode default error RawJSON: %v", err)
	}
	if _, present := raw["stream_id"]; present {
		t.Fatal("default error RawJSON contains stream_id, want omitted key")
	}
	t.Log("client: default Recv preserved typed error, supplied-field presence, and absent stream_id with nil Go error")
	close(defaultReceived)

	for i, item := range lanes {
		event, err := receivers[i].Recv(ctx)
		if err != nil {
			t.Fatalf("Recv progress for %s: %v", item.streamID, err)
		}
		progress, ok := event.AsAny().(responses.ResponsesServerEventResponseInWsProgress)
		if !ok {
			t.Fatalf("next event for %s = %T, want ResponsesServerEventResponseInWsProgress", item.streamID, event.AsAny())
		}
		if progress.Type != "response.in_progress" || progress.SequenceNumber != 1 || progress.StreamID != item.streamID ||
			progress.Response.ID != item.responseID || progress.Response.Status != "in_progress" {
			t.Fatalf("progress for %s = %s/%d/%s/%s/%s", item.streamID, progress.Type, progress.SequenceNumber, progress.StreamID, progress.Response.ID, progress.Response.Status)
		}
		t.Logf("client: consumed distinct progress %s on %s", progress.Response.ID, item.streamID)
	}
	for i, item := range lanes {
		result, err := receivers[i].FinalResponse(ctx)
		if err != nil || result == nil {
			t.Fatalf("FinalResponse for %s = %v, %v; want nonnil completed response", item.streamID, result, err)
		}
		if result.ID != item.responseID || result.Status != "completed" || result.OutputText() != item.text {
			t.Fatalf("FinalResponse for %s = %s/%s/%q, want %s/completed/%q", item.streamID, result.ID, result.Status, result.OutputText(), item.responseID, item.text)
		}
		t.Logf("client: completed %s with output %q", result.ID, result.OutputText())
	}
	if err := connection.Close(); err != nil {
		t.Errorf("connection Close: %v", err)
	}
	t.Log("client: bounded connection Close returned")
}

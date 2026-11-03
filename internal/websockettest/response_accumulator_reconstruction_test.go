package websockettest

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	wire "github.com/coder/websocket"
	"github.com/openai/openai-go/v3/responses"
)

// These fixtures match Ruby's final-item collector. Missing/null output cases
// exercise its compatibility behavior, not standard Response schema validity.
func TestResponsesAccumulatorReconstructedResponse(t *testing.T) {
	data, err := os.ReadFile("testdata/responses_reconstruction.json")
	if err != nil {
		t.Fatal(err)
	}
	var corpus struct {
		Cases []struct {
			Name        string            `json:"name"`
			ResetBefore []int             `json:"reset_before_event_indices"`
			Events      []json.RawMessage `json:"events"`
			Expected    struct {
				Result   string          `json:"result"`
				Response json.RawMessage `json:"response"`
			} `json:"expected"`
		} `json:"cases"`
	}
	if err := json.Unmarshal(data, &corpus); err != nil {
		t.Fatal(err)
	}
	for _, tc := range corpus.Cases {
		t.Run(tc.Name, func(t *testing.T) {
			frames := make([]string, len(tc.Events))
			for i, raw := range tc.Events {
				frames[i] = string(raw)
			}
			_, lane, _ := reconstructionConnection(t, frames)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			var acc responses.ResponseAccumulator
			for index, frame := range frames {
				for _, resetIndex := range tc.ResetBefore {
					if index == resetIndex {
						acc.Reset()
					}
				}
				event, err := lane.Recv(ctx)
				if err != nil {
					t.Fatal(err)
				}
				before := event.RawJSON()
				err = acc.AddEvent(event)
				if event.Type == "error" {
					var protocol *responses.ResponseProtocolError
					if !errors.As(err, &protocol) {
						t.Fatalf("AddEvent(error) = %v", err)
					}
				} else if err != nil {
					t.Fatal(err)
				}
				if event.RawJSON() != before {
					t.Fatal("AddEvent changed original event")
				}
				assertReconstructionJSON(t, before, frame)
				if acc.Snapshot().TerminalEvent == "" {
					if response, err := acc.ReconstructedResponse(); err == nil || response != nil {
						t.Fatal("unfinished turn became a response")
					}
				}
			}
			before := acc.Snapshot()
			details := acc.DetailedSnapshot()
			response, err := acc.ReconstructedResponse()
			if tc.Expected.Result != "response" {
				if response != nil || err == nil {
					t.Fatal("missing terminal or protocol error became a response")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			assertReconstructionJSON(t, response.RawJSON(), string(tc.Expected.Response))
			if !reflect.DeepEqual(acc.Snapshot(), before) || !reflect.DeepEqual(acc.DetailedSnapshot(), details) {
				t.Fatal("accessor changed existing projections")
			}
			// Each call returns fresh typed slices and unknown JSON fields.
			response.ID = "caller mutation"
			if len(response.Output) > 0 {
				response.Output[0].ID = "caller mutation"
			}
			if response.JSON.ExtraFields != nil {
				delete(response.JSON.ExtraFields, "future_metadata")
			}
			again, err := acc.ReconstructedResponse()
			if err != nil {
				t.Fatal(err)
			}
			assertReconstructionJSON(t, again.RawJSON(), string(tc.Expected.Response))
			if again.ID == "caller mutation" || (len(again.Output) > 0 && again.Output[0].ID == "caller mutation") {
				t.Fatal("returned response aliases retained state")
			}
			acc.Reset()
			if response, err := acc.ReconstructedResponse(); response != nil || err == nil {
				t.Fatal("Reset retained a terminal response")
			}
			assertReconstructionJSON(t, again.RawJSON(), string(tc.Expected.Response))
		})
	}
}

func reconstructionConnection(t *testing.T, frames []string) (*responses.ResponseConnection, *responses.ResponseLane, *responses.ResponseLane) {
	t.Helper()
	start := make(chan struct{})
	conn := laneTestConnection(t, responses.ResponseConnectionOptions{}, func(ctx context.Context, socket *wire.Conn) {
		select {
		case <-start:
		case <-ctx.Done():
			return
		}
		for _, frame := range frames {
			if err := socket.Write(ctx, wire.MessageText, []byte(frame)); err != nil {
				t.Errorf("write fixture: %v", err)
				return
			}
		}
		_ = socket.Close(wire.StatusNormalClosure, "")
	})
	lane, err := conn.Lane("reconstruction")
	if err != nil {
		t.Fatal(err)
	}
	other, err := conn.Lane("other")
	if err != nil {
		t.Fatal(err)
	}
	close(start)
	return conn, lane, other
}

func assertReconstructionJSON(t *testing.T, got, want string) {
	t.Helper()
	var a, b any
	if json.Unmarshal([]byte(got), &a) != nil || json.Unmarshal([]byte(want), &b) != nil || !reflect.DeepEqual(a, b) {
		t.Fatal("reconstructed response differs from expected wire fields")
	}
}

func TestResponsesAccumulatorReconstructionLaneErrorAndEOF(t *testing.T) {
	const done = `{"type":"response.output_item.done","stream_id":"reconstruction","output_index":0,"item":{"type":"message","id":"m","content":[{"type":"output_text","text":"final","annotations":[]}]}}`
	const terminal = `{"type":"response.completed","stream_id":"reconstruction","response":{"id":"r","status":"completed","output":null}}`
	for _, scenario := range []string{"other lane", "error then terminal", "EOF", "canceled receive"} {
		t.Run(scenario, func(t *testing.T) {
			frames := []string{done}
			switch scenario {
			case "other lane":
				frames = append(frames, strings.ReplaceAll(done, `"reconstruction"`, `"other"`), terminal)
			case "error then terminal":
				frames = append(frames, `{"type":"error","stream_id":"reconstruction","error":{"type":"server_error","message":"synthetic"}}`, terminal)
			}
			_, lane, other := reconstructionConnection(t, frames)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			var acc responses.ResponseAccumulator
			event, err := lane.Recv(ctx)
			if err != nil {
				t.Fatal(err)
			}
			if err = acc.AddEvent(event); err != nil {
				t.Fatal(err)
			}
			if scenario == "other lane" {
				event, err = other.Recv(ctx)
				if err != nil {
					t.Fatal(err)
				}
				if acc.AddEvent(event) == nil {
					t.Fatal("different lane accepted")
				}
			}
			if scenario == "error then terminal" {
				event, err = lane.Recv(ctx)
				if err != nil {
					t.Fatal(err)
				}
				if acc.AddEvent(event) == nil {
					t.Fatal("protocol error accepted")
				}
			}
			if scenario == "EOF" || scenario == "canceled receive" {
				receiveCtx := ctx
				if scenario == "canceled receive" {
					c, stop := context.WithCancel(ctx)
					stop()
					receiveCtx = c
				}
				if _, err = lane.Recv(receiveCtx); err == nil {
					t.Fatal("receive unexpectedly succeeded")
				}
				if response, reconstructionErr := acc.ReconstructedResponse(); response != nil || reconstructionErr == nil {
					t.Fatal("receive failure fabricated completion")
				}
				return
			}
			event, err = lane.Recv(ctx)
			if err != nil {
				t.Fatal(err)
			}
			if err = acc.AddEvent(event); err != nil {
				t.Fatal(err)
			}
			response, err := acc.ReconstructedResponse()
			if scenario == "error then terminal" {
				var protocol *responses.ResponseProtocolError
				if response != nil || !errors.As(err, &protocol) {
					t.Fatal("protocol error did not remain latched")
				}
			} else if err != nil || response.OutputText() != "final" {
				t.Fatal("other lane changed finalized records")
			}
		})
	}
}

func TestResponsesAccumulatorReconstructionDoesNotChangeNativeFinalResponse(t *testing.T) {
	frames := []string{`{"type":"response.output_item.done","stream_id":"reconstruction","output_index":0,"item":{"type":"message","id":"m","content":[{"type":"output_text","text":"final"}]}}`, `{"type":"response.completed","stream_id":"reconstruction","response":{"id":"r","status":"completed","output":null}}`}
	_, lane, _ := reconstructionConnection(t, frames)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	response, err := lane.FinalResponse(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(response.Output) != 0 || response.JSON.Output.Raw() != "null" {
		t.Fatal("native FinalResponse reconstructed output")
	}
}

func TestResponsesAccumulatorReconstructionRetainsOnlyFinalRecords(t *testing.T) {
	frames := []string{
		`{"type":"response.output_item.done","stream_id":"reconstruction","output_index":0,"item":{"type":"message","id":"m","content":[{"type":"output_text","text":"final","annotations":[]}]}}`,
		`{"type":"response.output_item.added","stream_id":"reconstruction","output_index":0,"item":{"type":"message","id":"replacement","content":[{"type":"output_text","text":"draft","annotations":[]}]}}`,
		`{"type":"response.output_item.added","stream_id":"reconstruction","output_index":1,"item":{"type":"message","id":"unfinished","content":[{"type":"output_text","text":"partial","annotations":[]}]}}`,
		`{"type":"response.output_item.done","stream_id":"reconstruction","output_index":9223372036854775700,"item":{"type":"future_tool","id":"future","future_metadata":{"number":9007199254740993}}}`,
		`{"type":"response.completed","stream_id":"reconstruction","response":{"id":"r","status":"completed","output":null}}`,
	}
	_, lane, _ := reconstructionConnection(t, frames)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var acc responses.ResponseAccumulator
	for range frames {
		event, err := lane.Recv(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if err := acc.AddEvent(event); err != nil {
			t.Fatal(err)
		}
		if event.Type == "response.output_item.done" {
			item := &event.OfResponsesServerEventResponseOutputItemWsDone.Item
			if len(item.Content) > 0 {
				item.Content[0].Text = "caller mutation"
			}
			item.ID = "caller mutation"
		}
	}
	response, err := acc.ReconstructedResponse()
	if err != nil {
		t.Fatal(err)
	}
	if response.OutputText() != "final" || len(response.Output) != 2 ||
		response.Output[1].Type != "future_tool" || !strings.Contains(response.Output[1].RawJSON(), "9007199254740993") ||
		acc.Snapshot().OutputText() != "draftpartial" {
		t.Fatal("unfinished projection replaced a finalized record, or unknown exact fields were lost")
	}
	response.Output[0].Content[0].Text = "caller mutation"
	again, err := acc.ReconstructedResponse()
	if err != nil || again.OutputText() != "final" {
		t.Fatal("returned nested content aliases retained state")
	}
}

func TestResponsesAccumulatorReconstructedLargeResponse(t *testing.T) {
	text := strings.Repeat("x", 65<<20)
	item := map[string]any{"type": "message", "id": "m", "content": []any{map[string]any{"type": "output_text", "text": text, "annotations": []any{}}}}
	frame, err := json.Marshal(map[string]any{"type": "response.output_item.done", "stream_id": "reconstruction", "output_index": 0, "item": item})
	if err != nil {
		t.Fatal(err)
	}
	// Keep reconstruction coverage above 64 MiB without a socket deadline.
	// TestResponsesWebSocketLargeFinalResponse covers the large-frame transport;
	// the fixture tests above cover reconstruction through a real connection.
	frames := [][]byte{frame, []byte(`{"type":"response.completed","stream_id":"reconstruction","response":{"id":"r","status":"completed"}}`)}
	var acc responses.ResponseAccumulator
	for _, raw := range frames {
		var event responses.ResponsesServerEventUnion
		if decodeErr := json.Unmarshal(raw, &event); decodeErr != nil {
			t.Fatal(decodeErr)
		}
		if addErr := acc.AddEvent(event); addErr != nil {
			t.Fatal(addErr)
		}
	}
	response, err := acc.ReconstructedResponse()
	if err != nil || response.OutputText() != text {
		t.Fatalf("large reconstruction failed: %v", err)
	}
}

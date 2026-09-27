package websockettest

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	wire "github.com/coder/websocket"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

func TestResponsesAccumulatorInterleavedTurns(t *testing.T) {
	// Deltas precede any terminal output. The peer corrects all three selected
	// fields in done events and omits output from both final responses.
	frames := []string{
		`{"type":"response.created","response":{"id":"resp_first","status":"in_progress"}}`,
		`{"type":"response.output_text.delta","item_id":"msg","output_index":4,"content_index":1,"delta":"part"}`,
		`{"type":"response.output_text.delta","item_id":"msg","output_index":4,"content_index":1,"delta":"ial"}`,
		`{"type":"response.output_item.added","output_index":0,"item":{"type":"function_call","id":"fc","call_id":"call_fc","name":"lookup","arguments":""}}`,
		`{"type":"response.function_call_arguments.delta","output_index":0,"item_id":"fc","delta":"{incomplete"}`,
		`{"type":"response.function_call_arguments.done","output_index":0,"item_id":"fc","arguments":"{\"value\":1}"}`,
		`{"type":"response.output_item.added","output_index":2,"item":{"type":"custom_tool_call","id":"ct","call_id":"call_ct","name":"freeform","input":""}}`,
		`{"type":"response.custom_tool_call_input.delta","output_index":2,"item_id":"ct","delta":"draft"}`,
		`{"type":"response.custom_tool_call_input.done","output_index":2,"item_id":"ct","input":"corrected input"}`,
		`{"type":"response.output_text.done","item_id":"msg","output_index":4,"content_index":1,"text":"corrected text"}`,
		`{"type":"response.future_event","private_field":"leave on raw event"}`,
		`{"type":"response.completed","response":{"id":"resp_first","status":"completed"}}`,
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		socket, err := wire.Accept(w, r, nil)
		if err != nil {
			t.Errorf("upgrade: %v", err)
			return
		}
		defer func() { _ = socket.CloseNow() }()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		for range 2 {
			if _, _, err := socket.Read(ctx); err != nil {
				t.Errorf("create: %v", err)
				return
			}
		}
		for _, frame := range frames {
			for _, id := range []string{"a", "b"} {
				body := frame
				if id == "b" {
					body = strings.ReplaceAll(body, "resp_first", "resp_other")
					body = strings.ReplaceAll(body, "corrected text", "other lane")
					body = strings.ReplaceAll(body, "response.completed", "response.incomplete")
					body = strings.ReplaceAll(body, `"status":"completed"`, `"status":"incomplete"`)
				}
				if err := socket.Write(ctx, wire.MessageText, withStreamID(t, []byte(body), id)); err != nil {
					t.Errorf("send %s: %v", id, err)
					return
				}
			}
		}
		// Reset must not close the socket: a second response uses the same lane.
		if _, _, err := socket.Read(ctx); err != nil {
			t.Errorf("second create: %v", err)
			return
		}
		for _, frame := range []string{
			`{"type":"response.created","response":{"id":"resp_next","status":"in_progress"}}`,
			`{"type":"response.output_text.delta","item_id":"temp","output_index":900000000,"content_index":70000000,"delta":"superseded"}`,
			`{"type":"response.failed","response":{"id":"resp_next","status":"failed","output":[{"type":"message","id":"final","role":"assistant","status":"incomplete","content":[{"type":"output_text","text":"authoritative","annotations":[]}]}]}}`,
		} {
			if err := socket.Write(ctx, wire.MessageText, withStreamID(t, []byte(frame), "a")); err != nil {
				t.Errorf("next turn: %v", err)
				return
			}
		}
		_, _, _ = socket.Read(ctx)
	}))
	defer server.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("fixture-key"))
	conn, err := client.Responses.Connect(ctx, responses.ResponseConnectionOptions{})
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Abort()
	laneA, err := conn.Lane("a")
	if err != nil {
		t.Fatal(err)
	}
	laneB, err := conn.Lane("b")
	if err != nil {
		t.Fatal(err)
	}
	turn := fixtureScenario(t, "completed_then_completed").Turns[0]
	for _, id := range []string{"a", "b"} {
		sendFixture(t, ctx, conn.Create, withStreamID(t, turn.Request, id))
	}
	var a, b responses.ResponseAccumulator
	var saved responses.ResponseAccumulatorSnapshot
	for i := range frames {
		ea, err := laneA.Recv(ctx)
		if err != nil {
			t.Fatal(err)
		}
		eb, err := laneB.Recv(ctx)
		if err != nil {
			t.Fatal(err)
		}
		raw := ea.RawJSON()
		if err := a.AddEvent(ea); err != nil {
			t.Fatal(err)
		}
		if i == 0 {
			if err := a.AddEvent(eb); err == nil {
				t.Fatal("accepted an event from a different lane")
			}
		}
		if err := b.AddEvent(eb); err != nil {
			t.Fatal(err)
		}
		if ea.RawJSON() != raw {
			t.Fatal("accumulation altered the received event")
		}
		if i == 1 {
			saved = a.Snapshot()
			if saved.OutputText() != "part" || saved.TerminalEvent != "" || saved.ResponseID != "resp_first" {
				t.Fatalf("before terminal: %+v", saved)
			}
		}
		if i == 2 && a.Snapshot().OutputText() != "partial" {
			t.Fatalf("incremental text: %q", a.Snapshot().OutputText())
		}
		if ea.Type == "response.future_event" && !strings.Contains(raw, "leave on raw event") {
			t.Fatal("unknown event data unavailable")
		}
	}
	got, other := a.Snapshot(), b.Snapshot()
	if got.TerminalEvent != "response.completed" || got.StreamID != "a" || got.OutputText() != "corrected text" ||
		other.TerminalEvent != "response.incomplete" || other.StreamID != "b" || other.ResponseID != "resp_other" ||
		other.OutputText() != "other lane" {
		t.Fatalf("lane snapshots: %+v, %+v", got, other)
	}
	if len(got.Output) != 3 || got.Output[0].OutputIndex != 0 || got.Output[0].Arguments != `{"value":1}` ||
		got.Output[0].CallID != "call_fc" || got.Output[0].Name != "lookup" ||
		got.Output[1].OutputIndex != 2 || got.Output[1].Input != "corrected input" || got.Output[1].CallID != "call_ct" ||
		got.Output[2].OutputIndex != 4 {
		t.Fatalf("tools and content lost or reordered: %+v", got.Output)
	}
	if saved.OutputText() != "part" {
		t.Fatalf("previously obtained snapshot mutated: %q", saved.OutputText())
	}
	got.Output[2].Text[1] = "caller mutation"
	if a.Snapshot().OutputText() != "corrected text" {
		t.Fatal("caller-owned snapshot changed accumulator state")
	}
	a.Reset()
	sendFixture(t, ctx, conn.Create, withStreamID(t, turn.Request, "a"))
	for range 3 {
		event, err := laneA.Recv(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if err := a.AddEvent(event); err != nil {
			t.Fatal(err)
		}
	}
	next := a.Snapshot()
	if next.ResponseID != "resp_next" || next.TerminalEvent != "response.failed" ||
		next.OutputText() != "authoritative" || len(next.Output) != 1 || saved.OutputText() != "part" {
		t.Fatalf("reset or supplied terminal output leaked prior data: %+v", next)
	}
}

func TestResponsesAccumulatorNeverInventsCompletion(t *testing.T) {
	for _, tc := range []struct {
		name, last  string
		protocolErr bool
	}{
		{name: "EOF"},
		{name: "missing response", last: `{"type":"response.completed"}`},
		{name: "null response", last: `{"type":"response.failed","response":null}`},
		{name: "invalid response", last: `{"type":"response.incomplete","response":42}`},
		{name: "API error", last: `{"type":"error","error":{"type":"invalid_request_error","code":"test","message":"test"}}`, protocolErr: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				socket, err := wire.Accept(w, r, nil)
				if err != nil {
					t.Errorf("upgrade: %v", err)
					return
				}
				defer func() { _ = socket.CloseNow() }()
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				for _, frame := range []string{
					`{"type":"response.created","response":{"id":"resp_partial","status":"in_progress"}}`,
					`{"type":"response.output_text.delta","item_id":"msg","output_index":0,"content_index":0,"delta":"partial"}`,
					tc.last,
				} {
					if frame != "" {
						if err := socket.Write(ctx, wire.MessageText, []byte(frame)); err != nil {
							t.Errorf("write: %v", err)
							return
						}
					}
				}
			}))
			defer server.Close()
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("fixture-key"))
			conn, err := client.Responses.Connect(ctx, responses.ResponseConnectionOptions{})
			if err != nil {
				t.Fatal(err)
			}
			defer conn.Abort()
			var acc responses.ResponseAccumulator
			for range 2 {
				event, err := conn.Recv(ctx)
				if err != nil {
					t.Fatal(err)
				}
				if err := acc.AddEvent(event); err != nil {
					t.Fatal(err)
				}
			}
			event, recvErr := conn.Recv(ctx)
			if tc.last == "" {
				if recvErr == nil {
					t.Fatal("EOF returned as event")
				}
			} else {
				if recvErr != nil {
					t.Fatal(recvErr)
				}
				addErr := acc.AddEvent(event)
				if addErr == nil {
					t.Fatal("accepted missing/invalid terminal or API error")
				}
				var protocol *responses.ResponseProtocolError
				if errors.As(addErr, &protocol) != tc.protocolErr {
					t.Fatalf("error identity: %v", addErr)
				}
			}
			result := acc.Snapshot()
			if result.ResponseID != "resp_partial" || result.OutputText() != "partial" || result.TerminalEvent != "" {
				t.Fatalf("interruption reported as completion: %+v", result)
			}
		})
	}
}

func TestResponsesAccumulatorReplacements(t *testing.T) {
	// All corrections arrive after a caller has already observed provisional
	// data, so silently leaving stale text or tool input would be visible.
	for _, tc := range []struct {
		name                          string
		initial, correction           string
		wantID, wantText, wantArgs    string
		wantInput, wantType, terminal string
		wantItems                     int
	}{
		{
			name:       "explicit empty terminal clears provisional output",
			initial:    `{"type":"response.output_text.delta","item_id":"msg","output_index":0,"content_index":0,"delta":"provisional"}`,
			correction: `{"type":"response.completed","response":{"id":"resp","status":"completed","output":[]}}`,
			terminal:   "response.completed",
		},
		{
			name:       "omitted terminal output retains projection",
			initial:    `{"type":"response.output_text.delta","item_id":"msg","output_index":0,"content_index":0,"delta":"provisional"}`,
			correction: `{"type":"response.completed","response":{"id":"resp","status":"completed"}}`,
			wantItems:  1, wantID: "msg", wantText: "provisional", terminal: "response.completed",
		},
		{
			name:       "null terminal output remains nullish",
			initial:    `{"type":"response.output_text.delta","item_id":"msg","output_index":0,"content_index":0,"delta":"provisional"}`,
			correction: `{"type":"response.incomplete","response":{"id":"resp","status":"incomplete","output":null}}`,
			wantItems:  1, wantID: "msg", wantText: "provisional", terminal: "response.incomplete",
		},
		{
			name:       "done message removes excess content",
			initial:    `{"type":"response.output_item.added","output_index":0,"item":{"type":"message","id":"msg","content":[{"type":"output_text","text":"first"},{"type":"output_text","text":"stale"}]}}`,
			correction: `{"type":"response.output_item.done","output_index":0,"item":{"type":"message","id":"msg","content":[{"type":"output_text","text":"corrected"}]}}`,
			wantItems:  1, wantID: "msg", wantText: "corrected", wantType: "message",
		},
		{
			name:       "done message with empty content",
			initial:    `{"type":"response.output_item.added","output_index":0,"item":{"type":"message","id":"msg","content":[{"type":"output_text","text":"remove"}]}}`,
			correction: `{"type":"response.output_item.done","output_index":0,"item":{"type":"message","id":"msg","content":[]}}`,
			wantItems:  1, wantID: "msg", wantType: "message",
		},
		{
			name:       "done item changes from message to function",
			initial:    `{"type":"response.output_item.added","output_index":0,"item":{"type":"message","id":"item","content":[{"type":"output_text","text":"stale"}]}}`,
			correction: `{"type":"response.output_item.done","output_index":0,"item":{"type":"function_call","id":"item","call_id":"call","name":"lookup","arguments":"{}"}}`,
			wantItems:  1, wantID: "item", wantArgs: "{}", wantType: "function_call",
		},
		{
			name:       "done item changes from function to custom input",
			initial:    `{"type":"response.output_item.added","output_index":0,"item":{"type":"function_call","id":"item","call_id":"call","name":"lookup","arguments":"{\"stale\":true}"}}`,
			correction: `{"type":"response.output_item.done","output_index":0,"item":{"type":"custom_tool_call","id":"item","call_id":"call","name":"corrected","input":"free form"}}`,
			wantItems:  1, wantID: "item", wantInput: "free form", wantType: "custom_tool_call",
		},
		{
			name:       "delta for different identity starts a fresh index",
			initial:    `{"type":"response.output_item.added","output_index":0,"item":{"type":"message","id":"msg_old","content":[{"type":"output_text","text":"stale"}]}}`,
			correction: `{"type":"response.output_text.delta","item_id":"msg_new","output_index":0,"content_index":0,"delta":"fresh"}`,
			wantItems:  1, wantID: "msg_new", wantText: "fresh",
		},
		{
			name:       "tool delta with different identity does not concatenate",
			initial:    `{"type":"response.output_item.added","output_index":0,"item":{"type":"function_call","id":"fc_old","call_id":"old","name":"lookup","arguments":"{\"old\":1}"}}`,
			correction: `{"type":"response.function_call_arguments.delta","item_id":"fc_new","output_index":0,"delta":"{\"new\":"}`,
			wantItems:  1, wantID: "fc_new", wantArgs: `{"new":`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			conn := laneTestConnection(t, responses.ResponseConnectionOptions{}, func(ctx context.Context, socket *wire.Conn) {
				for _, frame := range []string{
					`{"type":"response.created","response":{"id":"resp","status":"in_progress"}}`,
					tc.initial, tc.correction,
				} {
					if err := socket.Write(ctx, wire.MessageText, []byte(frame)); err != nil {
						t.Errorf("write event: %v", err)
						return
					}
				}
				_, _, _ = socket.Read(ctx)
			})
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			var acc responses.ResponseAccumulator
			for range 2 {
				event, err := conn.Recv(ctx)
				if err != nil {
					t.Fatal(err)
				}
				if err := acc.AddEvent(event); err != nil {
					t.Fatal(err)
				}
			}
			before := acc.Snapshot()
			beforeText := before.OutputText()
			if len(before.Output) != 1 {
				t.Fatalf("initial item not available: %+v", before)
			}
			event, err := conn.Recv(ctx)
			if err != nil {
				t.Fatal(err)
			}
			if err := acc.AddEvent(event); err != nil {
				t.Fatal(err)
			}
			after := acc.Snapshot()
			if len(after.Output) != tc.wantItems || after.OutputText() != tc.wantText || after.TerminalEvent != tc.terminal {
				t.Fatalf("after correction: %+v, text %q; want %d items, text %q, terminal %q",
					after, after.OutputText(), tc.wantItems, tc.wantText, tc.terminal)
			}
			if tc.wantItems > 0 {
				item := after.Output[0]
				if item.ItemID != tc.wantID || item.Type != tc.wantType || item.Arguments != tc.wantArgs || item.Input != tc.wantInput {
					t.Fatalf("corrected item: %+v; want ID %s, type %s, args %q, input %q",
						item, tc.wantID, tc.wantType, tc.wantArgs, tc.wantInput)
				}
			}
			if before.OutputText() != beforeText {
				t.Fatal("item replacement changed previously published snapshot")
			}
		})
	}
}

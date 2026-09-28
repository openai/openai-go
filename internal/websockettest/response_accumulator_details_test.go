package websockettest

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"
	"time"

	wire "github.com/coder/websocket"
	"github.com/openai/openai-go/v3/responses"
)

func TestResponsesAccumulatorDetailsCarryWireProgressWithoutBecomingAResponse(t *testing.T) {
	const large int64 = 9223372036854775700
	frames := []string{
		`{"type":"response.output_text.annotation.added","stream_id":"other","item_id":"other","output_index":1,"content_index":0,"annotation_index":0,"annotation":{"type":"future_fixture"}}`,
		`{"type":"response.created","response":{"id":"resp","status":"in_progress","output":null,"model":"fixture_model","metadata":{"large_fixture":9007199254740993},"future_fixture":null}}`,
		`{"type":"response.output_text.delta","item_id":"msg","output_index":9223372036854775700,"content_index":9223372036854775701,"delta":"hel","logprobs":[{"token":"hel","logprob":-0.3,"extra_fixture":[true]}]}`,
		`{"type":"response.output_text.delta","item_id":"msg","output_index":9223372036854775700,"content_index":9223372036854775701,"delta":"lo","logprobs":[{"token":"lo","logprob":-0.1,"top_logprobs":null}]}`,
		`{"type":"response.output_text.annotation.added","item_id":"msg","output_index":9223372036854775700,"content_index":9223372036854775701,"annotation_index":9223372036854775702,"annotation":{"type":"url_citation","title":"fixture","url":"https://example.com/fixture","start_index":0,"end_index":5,"future_fixture":[1]}}`,
		`{"type":"response.output_text.annotation.added","item_id":"wrong","output_index":9223372036854775700,"content_index":9223372036854775701,"annotation_index":0}`,
		`{"type":"response.output_text.done","item_id":"msg","output_index":9223372036854775700,"content_index":9223372036854775701,"text":"corrected","logprobs":[{"token":"corrected","logprob":-0.2}]}`,
		`{"type":"response.incomplete","response":{"id":"resp","status":"incomplete","output":null,"metadata":{"after_fixture":true}}}`,
	}
	conn := laneTestConnection(t, responses.ResponseConnectionOptions{}, func(ctx context.Context, socket *wire.Conn) {
		for _, frame := range frames {
			if err := socket.Write(ctx, wire.MessageText, []byte(frame)); err != nil {
				t.Errorf("write synthetic progress: %v", err)
				return
			}
		}
		_, _, _ = socket.Read(ctx)
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var acc responses.ResponseAccumulator
	for range 6 {
		event, err := conn.Recv(ctx)
		if err != nil {
			t.Fatal(err)
		}
		before := event.RawJSON()
		if err := acc.AddEvent(event); err != nil {
			t.Fatal(err)
		}
		if before != event.RawJSON() {
			t.Fatal("mutated original event")
		}
	}
	legacy := acc.Snapshot()
	saved := acc.DetailedSnapshot()
	if legacy.OutputText() != "hello" || saved.TerminalEvent != "" || saved.ResponseID != "resp" {
		t.Fatal("annotation changed legacy turn, or invented a terminal")
	}
	if string(saved.Response["metadata"]) != `{"large_fixture":9007199254740993}` ||
		string(saved.Response["future_fixture"]) != "null" || saved.Response["output"] != nil {
		t.Fatal("lost response metadata, exact number, or conflated output with a Response")
	}
	if len(saved.Output) != 1 || saved.Output[0].OutputIndex != large {
		t.Fatal("lost sparse output position")
	}
	item := &saved.Output[0]
	if string(item.Item["id"]) != `"msg"` || item.Item["type"] != nil {
		t.Fatal("missing scaffolding must not invent an item type")
	}
	part := item.Content[large+1]
	if string(part["text"]) != `"hello"` ||
		string(part["logprobs"]) != `[{"token":"hel","logprob":-0.3,"extra_fixture":[true]},{"token":"lo","logprob":-0.1,"top_logprobs":null}]` {
		t.Fatal("lost provisional logprobs or coerced into unsupported final Logprob")
	}
	if len(item.Annotations[large+1]) != 1 || len(item.Annotations[large+1][large+2]) == 0 {
		t.Fatal("lost sparse annotation")
	}
	// RawMessage is mutable. A copy of every nested caller-owned value is necessary.
	saved.Response["metadata"][0] = 'x'
	part["logprobs"][1] = 'x'
	item.Annotations[large+1][large+2][0] = 'x'
	clean := acc.DetailedSnapshot()
	if !json.Valid(clean.Response["metadata"]) || !json.Valid(clean.Output[0].Content[large+1]["logprobs"]) ||
		!json.Valid(clean.Output[0].Annotations[large+1][large+2]) {
		t.Fatal("caller mutation reached accumulator state")
	}
	done, err := conn.Recv(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if err = acc.AddEvent(done); err != nil {
		t.Fatal(err)
	}
	projected := acc.DetailedSnapshot()
	if string(projected.Output[0].Content[large+1]["logprobs"]) != `[{"token":"corrected","logprob":-0.2}]` ||
		len(projected.Output[0].Annotations[large+1]) != 1 ||
		acc.Snapshot().OutputText() != "corrected" || legacy.OutputText() != "hello" {
		t.Fatal("text done did not replace only supplied text/logprobs")
	}
	terminal, err := conn.Recv(ctx)
	if err != nil {
		t.Fatal(err)
	}
	rawTerminal := terminal.RawJSON()
	if err = acc.AddEvent(terminal); err != nil {
		t.Fatal(err)
	}
	if acc.DetailedSnapshot().TerminalEvent != "response.incomplete" ||
		string(acc.DetailedSnapshot().Response["metadata"]) != `{"after_fixture":true}` || terminal.RawJSON() != rawTerminal {
		t.Fatal("terminal authority lost")
	}
	acc.Reset()
	empty := acc.DetailedSnapshot()
	if empty.Response != nil || len(empty.Output) != 0 || !json.Valid(clean.Response["metadata"]) {
		t.Fatal("reset retained current data or altered an earlier snapshot")
	}
}

func TestResponsesAccumulatorDetailsReplacePartsAndKeepOriginalUnknownFields(t *testing.T) {
	frames := []string{
		`{"type":"response.output_item.added","output_index":5,"item":{"type":"message","id":"m","role":"assistant","status":"in_progress","fixture_score":9007199254740993,"content":[{"type":"output_text","text":"early","annotations":[{"type":"file_citation","file_id":"file_fixture","filename":"synthetic.txt","index":0}],"logprobs":null},{"type":"refusal","refusal":"synthetic refusal"},{"type":"future_part","fixture_data":{"maybe":null}}]}}`,
		`{"type":"response.output_text.annotation.added","item_id":"m","output_index":5,"content_index":0,"annotation_index":9,"annotation":{"type":"future_citation","fixture_data":[1]}}`,
		`{"type":"response.content_part.done","item_id":"m","output_index":5,"content_index":0,"part":{"type":"output_text","text":"final","annotations":[],"logprobs":[]}}`,
		`{"type":"response.output_item.done","output_index":5,"item":{"type":"mcp_call","id":"mcp","name":"fixture_search","server_label":"fixture_server","arguments":"{}","output":null,"error":null,"fixture_number":9007199254740993}}`,
		`{"type":"response.output_item.done","output_index":3,"item":{"type":"future_fixture","id":"new","content":{"value":null}}}`,
		`{"type":"response.failed","response":{"id":"r","status":"failed","output":[]}}`,
	}
	conn := laneTestConnection(t, responses.ResponseConnectionOptions{}, func(ctx context.Context, socket *wire.Conn) {
		for _, frame := range frames {
			if err := socket.Write(ctx, wire.MessageText, []byte(frame)); err != nil {
				t.Errorf("write synthetic replacement: %v", err)
				return
			}
		}
		_, _, _ = socket.Read(ctx)
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var acc responses.ResponseAccumulator
	for i := range frames {
		event, err := conn.Recv(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if err = acc.AddEvent(event); err != nil {
			t.Fatalf("frame %d: %v", i, err)
		}
		d := acc.DetailedSnapshot()
		switch i {
		case 0:
			p := d.Output[0]
			if p.Item["content"] != nil || string(p.Item["fixture_score"]) != "9007199254740993" ||
				string(p.Content[1]["refusal"]) != `"synthetic refusal"` ||
				string(p.Content[2]["fixture_data"]) != `{"maybe":null}` ||
				string(p.Content[0]["logprobs"]) != "null" || len(p.Annotations[0]) != 1 {
				t.Fatal("incomplete full item")
			}
		case 1:
			if len(d.Output[0].Annotations[0]) != 2 {
				t.Fatal("streamed annotation did not enrich original citations")
			}
		case 2:
			if len(d.Output[0].Annotations[0]) != 0 || string(d.Output[0].Content[0]["logprobs"]) != "[]" ||
				string(d.Output[0].Content[1]["refusal"]) != `"synthetic refusal"` {
				t.Fatal("part replacement retained old annotations or erased its neighbor")
			}
		case 3:
			p := d.Output[0]
			if len(p.Content) != 0 || len(p.Annotations) != 0 || string(p.Item["server_label"]) != `"fixture_server"` ||
				string(p.Item["arguments"]) != `"{}"` || string(p.Item["error"]) != "null" {
				t.Fatal("MCP metadata lost or obsolete text/citations retained")
			}
		case 4:
			if !reflect.DeepEqual(d.Output[0].Item["content"], json.RawMessage(`{"value":null}`)) {
				t.Fatal("unknown item's content was treated as a message")
			}
		case 5:
			if len(d.Output) != 0 || d.TerminalEvent != "response.failed" {
				t.Fatal("explicit empty terminal did not supersede projection")
			}
		}
	}
}

func TestResponsesAccumulatorDetailsDoNotIndexNullOrMalformedOriginalAnnotations(t *testing.T) {
	for _, original := range []string{"null", `{"unexpected":"fixture"}`} {
		t.Run(original, func(t *testing.T) {
			conn := laneTestConnection(t, responses.ResponseConnectionOptions{}, func(ctx context.Context, socket *wire.Conn) {
				for _, frame := range []string{
					`{"type":"response.output_item.added","output_index":0,"item":{"type":"message","id":"m","content":[{"type":"output_text","text":"hi","annotations":` + original + `}]}}`,
					`{"type":"response.output_text.annotation.added","output_index":0,"content_index":0,"item_id":"m","annotation_index":9,"annotation":{"type":"future_fixture","data":null}}`,
					`{"type":"response.incomplete","response":{"id":"resp","status":"incomplete","output":null}}`,
				} {
					if err := socket.Write(ctx, wire.MessageText, []byte(frame)); err != nil {
						t.Errorf("write synthetic null annotations: %v", err)
						return
					}
				}
				_, _, _ = socket.Read(ctx)
			})
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			var acc responses.ResponseAccumulator
			for range 3 {
				event, err := conn.Recv(ctx)
				if err != nil {
					t.Fatal(err)
				}
				if err = acc.AddEvent(event); err != nil {
					t.Fatal(err)
				}
			}
			snapshot := acc.DetailedSnapshot()
			item := snapshot.Output[0]
			if string(item.Content[0]["annotations"]) != original ||
				len(item.Annotations[0]) != 1 ||
				string(item.Annotations[0][9]) != `{"type":"future_fixture","data":null}` {
				t.Fatal("scalar or object annotations invented a sequential annotation")
			}
		})
	}
}

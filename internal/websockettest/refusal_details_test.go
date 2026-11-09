package websockettest

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/openai/openai-go/v3/responses"
)

func TestResponsesAccumulatorRefusalLifecycle(t *testing.T) {
	const large int64 = 9223372036854775700
	var acc responses.ResponseAccumulator
	add := func(raw string) {
		t.Helper()
		var event responses.ResponsesServerEventUnion
		if err := json.Unmarshal([]byte(raw), &event); err != nil {
			t.Fatal(err)
		}
		if err := acc.AddEvent(event); err != nil {
			t.Fatal(err)
		}
		if event.RawJSON() != raw {
			t.Fatal("caller event changed")
		}
	}
	part := func() map[string]json.RawMessage {
		t.Helper()
		s := acc.DetailedSnapshot()
		if len(s.Output) != 1 || s.Output[0].OutputIndex != large || len(s.Output[0].Content) != 1 {
			t.Fatal("sparse observed output/content did not stay bounded")
		}
		return s.Output[0].Content[large+1]
	}
	add(`{"type":"response.created","response":{"id":"resp_refusal","status":"in_progress","object":"response"}}`)
	add(`{"type":"response.output_item.added","output_index":9223372036854775700,"item":{"id":"msg","type":"message","role":"assistant","status":"in_progress","content":[]}}`)
	add(`{"type":"response.content_part.added","output_index":9223372036854775700,"content_index":9223372036854775701,"item_id":"msg","part":{"type":"refusal","refusal":"already ","fixture":9007199254740993}}`)
	add(`{"type":"response.refusal.delta","output_index":9223372036854775700,"content_index":9223372036854775701,"item_id":"msg","delta":"拒"}`)
	saved := acc.DetailedSnapshot()
	if string(part()["refusal"]) != `"already 拒"` {
		t.Fatalf("initial part seed + delta: %s", part()["refusal"])
	}
	add(`{"type":"response.refusal.delta","output_index":9223372036854775700,"content_index":9223372036854775701,"item_id":"msg","delta":"否🙂"}`)
	p := part()
	if string(p["refusal"]) != `"already 拒否🙂"` || string(p["fixture"]) != "9007199254740993" {
		t.Fatal("refusal data or original part fields lost")
	}
	if string(saved.Output[0].Content[large+1]["refusal"]) != `"already 拒"` ||
		saved.TerminalEvent != "" || acc.Snapshot().OutputText() != "" {
		t.Fatal("mutated older snapshot, invented final, or projected refusal as output text")
	}
	p["refusal"][1] = 'X'
	if string(part()["refusal"]) != `"already 拒否🙂"` {
		t.Fatal("caller mutation reached live accumulator")
	}
	add(`{"type":"response.refusal.done","output_index":9223372036854775700,"content_index":9223372036854775701,"item_id":"msg","refusal":""}`)
	if string(part()["refusal"]) != `""` {
		t.Fatal("empty done did not authoritatively replace progress")
	}
	add(`{"type":"response.content_part.done","output_index":9223372036854775700,"content_index":9223372036854775701,"item_id":"msg","part":{"type":"output_text","text":"later text","annotations":[]}}`)
	if part()["refusal"] != nil || string(part()["text"]) != `"later text"` {
		t.Fatal("replaced part retained refusal")
	}
	add(`{"type":"response.completed","response":{"id":"resp_refusal","status":"completed","output":[]}}`)
	if d := acc.DetailedSnapshot(); d.TerminalEvent != "response.completed" || len(d.Output) != 0 {
		t.Fatal("valid terminal empty output did not replace progress")
	}
	acc.Reset()
	if d := acc.DetailedSnapshot(); d.TerminalEvent != "" || len(d.Output) != 0 {
		t.Fatal("reset retained refusal data or binding")
	}
}

func TestResponsesAccumulatorRefusalScopeAndValidation(t *testing.T) {
	for _, tc := range []struct {
		name string
		part string
		bad  string
	}{
		{"other lane", `{"type":"refusal","refusal":"seed"}`, `{"type":"response.refusal.delta","stream_id":"other","item_id":"msg","output_index":0,"content_index":0,"delta":"wrong"}`},
		{"fractional index", `{"type":"refusal","refusal":"seed"}`, `{"type":"response.refusal.delta","item_id":"msg","output_index":0.5,"content_index":0,"delta":"wrong"}`},
		{"numeric delta", `{"type":"refusal","refusal":"seed"}`, `{"type":"response.refusal.delta","item_id":"msg","output_index":0,"content_index":0,"delta":7}`},
		{"missing done refusal", `{"type":"refusal","refusal":"seed"}`, `{"type":"response.refusal.done","item_id":"msg","output_index":0,"content_index":0}`},
		{"malformed prior refusal", `{"type":"refusal","refusal":7}`, `{"type":"response.refusal.delta","item_id":"msg","output_index":0,"content_index":0,"delta":"wrong"}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var acc responses.ResponseAccumulator
			read := func(raw string) responses.ResponsesServerEventUnion {
				t.Helper()
				var event responses.ResponsesServerEventUnion
				if err := json.Unmarshal([]byte(raw), &event); err != nil {
					t.Fatal(err)
				}
				return event
			}
			initial := read(`{"type":"response.output_item.added","output_index":0,"item":{"id":"msg","type":"message","content":[` + tc.part + `]}}`)
			if err := acc.AddEvent(initial); err != nil {
				t.Fatal("previously raw nonselected refusal must remain accepted", err)
			}
			before := acc.DetailedSnapshot()
			event := read(tc.bad)
			if err := acc.AddEvent(event); err == nil {
				t.Fatal("accepted invalid selected update")
			}
			if !reflect.DeepEqual(before, acc.DetailedSnapshot()) || event.RawJSON() != tc.bad {
				t.Fatal("rejected update mutated owned state or input")
			}
			done := read(`{"type":"response.refusal.done","item_id":"msg","output_index":0,"content_index":0,"refusal":"corrected"}`)
			if err := acc.AddEvent(done); err != nil {
				t.Fatal("valid replacement of malformed/old data must recover", err)
			}
			if got := string(acc.DetailedSnapshot().Output[0].Content[0]["refusal"]); got != `"corrected"` {
				t.Fatalf("after correction: %s", got)
			}
		})
	}
}

func TestResponsesAccumulatorRefusalWithoutPartAndUnknownReplacement(t *testing.T) {
	var acc responses.ResponseAccumulator
	for _, raw := range []string{
		`{"type":"response.refusal.delta","item_id":"first","output_index":1,"content_index":2,"delta":"first"}`,
		`{"type":"response.refusal.done","item_id":"new","output_index":1,"content_index":2,"refusal":"second"}`,
	} {
		var event responses.ResponsesServerEventUnion
		if err := json.Unmarshal([]byte(raw), &event); err != nil {
			t.Fatal(err)
		}
		if err := acc.AddEvent(event); err != nil {
			t.Fatal(err)
		}
	}
	d := acc.DetailedSnapshot()
	if len(d.Output) != 1 || string(d.Output[0].Item["id"]) != `"new"` ||
		string(d.Output[0].Content[2]["refusal"]) != `"second"` || d.Output[0].Content[2]["type"] != nil {
		t.Fatal("missing part or reused index did not preserve actual data without fabricating kind")
	}
	for _, raw := range []string{
		`{"type":"response.content_part.done","item_id":"new","output_index":1,"content_index":2,"part":{"type":"future_part","refusal":{"fixture":"untyped"},"more":null}}`,
		`{"type":"response.refusal.done","item_id":"new","output_index":1,"content_index":2,"refusal":"must not coerce future part"}`,
	} {
		var event responses.ResponsesServerEventUnion
		if err := json.Unmarshal([]byte(raw), &event); err != nil {
			t.Fatal(err)
		}
		if err := acc.AddEvent(event); err != nil {
			t.Fatal(err)
		}
	}
	p := acc.DetailedSnapshot().Output[0].Content[2]
	if string(p["refusal"]) != `{"fixture":"untyped"}` || string(p["more"]) != "null" ||
		string(d.Output[0].Content[2]["refusal"]) != `"second"` {
		t.Fatal("future part fields or prior snapshot overwritten by mismatch")
	}
}

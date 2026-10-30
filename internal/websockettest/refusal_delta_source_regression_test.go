package websockettest

import (
	"encoding/json"
	"github.com/openai/openai-go/v3/responses"
	"testing"
)

func TestDetailedSnapshotAccumulatesKnownRefusalDeltas(t *testing.T) {
	var a responses.ResponseAccumulator
	for _, raw := range []string{
		`{"type":"response.created","response":{"id":"resp_refusal","status":"in_progress","object":"response"}}`,
		`{"type":"response.output_item.added","output_index":0,"item":{"id":"msg","type":"message","role":"assistant","status":"in_progress","content":[]}}`,
		`{"type":"response.content_part.added","output_index":0,"content_index":0,"item_id":"msg","part":{"type":"refusal","refusal":"","part_extension":"kept"}}`,
		`{"type":"response.refusal.delta","output_index":0,"content_index":0,"item_id":"msg","delta":"syn"}`,
		`{"type":"response.refusal.delta","output_index":0,"content_index":0,"item_id":"msg","delta":"thetic"}`,
	} {
		func() {
			var event responses.ResponsesServerEventUnion
			if err := json.Unmarshal([]byte(raw), &event); err != nil {
				t.Fatal(err)
			}
			if err := a.AddEvent(event); err != nil {
				t.Fatal(err)
			}
			if event.RawJSON() != raw {
				t.Fatal("caller event mutated")
			}
		}()
	}
	d := a.DetailedSnapshot()
	if len(d.Output) != 1 || len(d.Output[0].Content) != 1 {
		t.Fatal("missing actual refusal part")
	}
	if string(d.Output[0].Content[0]["part_extension"]) != `"kept"` {
		t.Fatal("part field not retained")
	}
	if string(d.Output[0].Content[0]["refusal"]) != `"synthetic"` {
		t.Fatalf("two known refusal deltas produced %s, expected the synthetic joined value", d.Output[0].Content[0]["refusal"])
	}
}

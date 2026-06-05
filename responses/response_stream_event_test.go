package responses_test

import (
	"encoding/json"
	"testing"

	"github.com/openai/openai-go/v3/responses"
)

func TestResponseStreamEventShellCallCommandDelta(t *testing.T) {
	raw := []byte(`{
		"type": "response.shell_call_command.delta",
		"delta": "echo hello",
		"item_id": "rs_123",
		"output_index": 0,
		"sequence_number": 1
	}`)

	var event responses.ResponseStreamEventUnion
	if err := json.Unmarshal(raw, &event); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	variant, ok := event.AsAny().(responses.ResponseShellCallCommandDeltaEvent)
	if !ok {
		t.Fatalf("AsAny() = %T, want ResponseShellCallCommandDeltaEvent", event.AsAny())
	}
	if variant.Delta != "echo hello" {
		t.Fatalf("Delta = %q, want %q", variant.Delta, "echo hello")
	}
	if variant.ItemID != "rs_123" {
		t.Fatalf("ItemID = %q, want %q", variant.ItemID, "rs_123")
	}
	if variant.OutputIndex != 0 {
		t.Fatalf("OutputIndex = %d, want 0", variant.OutputIndex)
	}
	if variant.SequenceNumber != 1 {
		t.Fatalf("SequenceNumber = %d, want 1", variant.SequenceNumber)
	}
}

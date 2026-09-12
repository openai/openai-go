package openai_test

import (
	"encoding/json"
	"testing"

	openai "github.com/openai/openai-go/v3"
)

func TestBetaResponseStreamEventShellCallOutputContentDelta(t *testing.T) {
	raw := []byte(`{
		"type": "response.shell_call_output_content.delta",
		"sequence_number": 1,
		"item_id": "sh_123",
		"output_index": 0,
		"command_index": 0,
		"delta": {"stdout": "hello\n", "stderr": "warning\n"}
	}`)

	var event openai.BetaResponseStreamEventUnion
	if err := json.Unmarshal(raw, &event); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	variant, ok := event.AsAny().(openai.BetaResponseShellCallOutputContentDeltaEvent)
	if !ok {
		t.Fatalf("AsAny() = %T, want BetaResponseShellCallOutputContentDeltaEvent", event.AsAny())
	}
	if variant.Delta.Stdout != "hello\n" {
		t.Fatalf("Delta.Stdout = %q, want %q", variant.Delta.Stdout, "hello\n")
	}
	if variant.Delta.Stderr != "warning\n" {
		t.Fatalf("Delta.Stderr = %q, want %q", variant.Delta.Stderr, "warning\n")
	}
}

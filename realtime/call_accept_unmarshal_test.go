package realtime_test

import (
	"encoding/json"
	"testing"

	"github.com/openai/openai-go/v3/realtime"
)

func TestCallAcceptParamsJSONRoundTrip(t *testing.T) {
	input := []byte(`{"type":"realtime","model":"gpt-realtime","instructions":"hello"}`)

	var params realtime.CallAcceptParams
	if err := json.Unmarshal(input, &params); err != nil {
		t.Fatalf("unmarshal CallAcceptParams: %v", err)
	}
	if got := params.RealtimeSessionCreateRequest.Model; got != realtime.RealtimeSessionCreateRequestModelGPTRealtime {
		t.Errorf("Model = %q, want %q", got, realtime.RealtimeSessionCreateRequestModelGPTRealtime)
	}
	if got := params.RealtimeSessionCreateRequest.Instructions.Value; got != "hello" {
		t.Errorf("Instructions = %q, want %q", got, "hello")
	}

	output, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("marshal CallAcceptParams: %v", err)
	}
	var roundTrip map[string]any
	if err := json.Unmarshal(output, &roundTrip); err != nil {
		t.Fatalf("unmarshal round-trip JSON: %v", err)
	}
	if got := roundTrip["model"]; got != "gpt-realtime" {
		t.Errorf("round-trip model = %q, want %q", got, "gpt-realtime")
	}
	if got := roundTrip["instructions"]; got != "hello" {
		t.Errorf("round-trip instructions = %q, want %q", got, "hello")
	}
}

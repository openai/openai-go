package realtime_test

import (
	"testing"

	"github.com/openai/openai-go/v3/realtime"
)

// Existing callers can choose a model dynamically and consume it as a string.
func TestRealtimeSessionModelStringCompatibility(t *testing.T) {
	model := "gpt-realtime-caller-selected"
	request := realtime.RealtimeSessionCreateRequestParam{Model: model}
	var got string = request.Model
	if got != model {
		t.Errorf("model = %q, want %q", got, model)
	}

	var advertised string = realtime.RealtimeSessionCreateRequestModelGPTRealtime
	if advertised != "gpt-realtime" {
		t.Errorf("advertised model = %q, want %q", advertised, "gpt-realtime")
	}
}

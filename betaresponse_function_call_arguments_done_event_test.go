package openai

import "testing"

func TestBetaResponseFunctionCallArgumentsDoneEventPreservesPositionalLayout(t *testing.T) {
	// Existing callers may use the field order from before Name became optional.
	event := BetaResponseFunctionCallArgumentsDoneEvent{
		"{}", "fc_test", "get_weather", 2, 14, "response.function_call_arguments.done",
		BetaResponseFunctionCallArgumentsDoneEventAgent{},
		BetaResponseFunctionCallArgumentsDoneEvent{}.JSON,
	}
	if event.Name != "get_weather" || event.ItemID != "fc_test" || event.OutputIndex != 2 || event.SequenceNumber != 14 {
		t.Fatalf("positional event fields changed: %+v", event)
	}
}

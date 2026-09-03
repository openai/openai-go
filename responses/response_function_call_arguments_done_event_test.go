package responses

import "testing"

func TestResponseFunctionCallArgumentsDoneEventPreservesPositionalLayout(t *testing.T) {
	// Existing callers may use the field order from before Name became optional.
	event := ResponseFunctionCallArgumentsDoneEvent{
		"{}", "fc_test", "get_weather", 2, 14, "response.function_call_arguments.done",
		ResponseFunctionCallArgumentsDoneEvent{}.JSON,
	}
	if event.Name != "get_weather" || event.ItemID != "fc_test" || event.OutputIndex != 2 || event.SequenceNumber != 14 {
		t.Fatalf("positional event fields changed: %+v", event)
	}
}

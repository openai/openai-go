package openai_test

import (
	"encoding/json"
	"reflect"
	"testing"

	openai "github.com/openai/openai-go/v3"
)

func TestChatCompletionParamsPreserveUnknownAssistantContentParts(t *testing.T) {
	body := []byte(`{
		"model":"gpt-4o-mini",
		"messages":[{
			"role":"assistant",
			"content":[
				{"type":"extension","text":"unknown","metadata":{"source":"fixture"}},
				{"type":"text","text":"kept"}
			]
		}]
	}`)

	var params openai.ChatCompletionNewParams
	if err := json.Unmarshal(body, &params); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(params.Messages) != 1 || params.Messages[0].OfAssistant == nil {
		t.Fatalf("assistant message was not decoded: %#v", params.Messages)
	}
	parts := params.Messages[0].OfAssistant.Content.OfArrayOfContentParts
	if len(parts) != 2 {
		t.Fatalf("decoded content parts = %d, want 2", len(parts))
	}
	unknown, ok := parts[0].Overrides()
	if !ok {
		t.Fatal("unknown content part was not preserved as raw JSON")
	}
	var gotUnknown any
	if err := json.Unmarshal(unknown.(json.RawMessage), &gotUnknown); err != nil {
		t.Fatalf("decode preserved content part: %v", err)
	}
	var wantUnknown any
	if err := json.Unmarshal([]byte(`{"type":"extension","text":"unknown","metadata":{"source":"fixture"}}`), &wantUnknown); err != nil {
		t.Fatalf("decode expected content part: %v", err)
	}
	if !reflect.DeepEqual(gotUnknown, wantUnknown) {
		t.Fatalf("preserved content part = %#v, want %#v", gotUnknown, wantUnknown)
	}
	if got := parts[1].GetText(); got == nil || *got != "kept" {
		t.Fatalf("supported text content was not preserved: %#v", parts[1])
	}

	encoded, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var roundTrip map[string]any
	if err := json.Unmarshal(encoded, &roundTrip); err != nil {
		t.Fatalf("decode marshaled params: %v", err)
	}
	messages := roundTrip["messages"].([]any)
	content := messages[0].(map[string]any)["content"].([]any)
	if !reflect.DeepEqual(content[0], wantUnknown) {
		t.Fatalf("round-tripped content part = %#v, want %#v", content[0], wantUnknown)
	}
}

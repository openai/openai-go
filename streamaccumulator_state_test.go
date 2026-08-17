package openai_test

import (
	"testing"

	openai "github.com/openai/openai-go/v3"
)

const (
	pr688EmptyToolCallsChunk = `{"id":"test","choices":[{"index":0,"delta":{"tool_calls":[]}}]}`
	pr688NullToolCallsChunk  = `{"id":"test","choices":[{"index":0,"delta":{"tool_calls":null}}]}`
	pr688ToolCallChunk       = `{"id":"test","choices":[{"index":0,"delta":{"tool_calls":[{"index":0,"id":"call_123","type":"function","function":{"name":"test_func","arguments":"{}"}}]}}]}`
	pr688FinishedChunk       = `{"id":"test","choices":[{"index":0,"delta":{},"finish_reason":"tool_calls"}]}`
)

func TestAccumulatorEmptyAndNullToolCallStateTransitions(t *testing.T) {
	for _, tc := range []struct {
		name  string
		chunk string
	}{
		{name: "empty array", chunk: pr688EmptyToolCallsChunk},
		{name: "null", chunk: pr688NullToolCallsChunk},
	} {
		tc := tc

		t.Run(tc.name+"/without active tool call", func(t *testing.T) {
			var acc openai.ChatCompletionAccumulator

			pr688AddChunk(t, &acc, tc.chunk)
			pr688AssertNoFinishedToolCall(t, &acc)

			pr688AddChunk(t, &acc, pr688FinishedChunk)
			pr688AssertNoFinishedToolCall(t, &acc)
		})

		t.Run(tc.name+"/before active tool call", func(t *testing.T) {
			var acc openai.ChatCompletionAccumulator

			pr688AddChunk(t, &acc, tc.chunk)
			pr688AssertNoFinishedToolCall(t, &acc)

			pr688AddChunk(t, &acc, pr688ToolCallChunk)
			pr688AssertNoFinishedToolCall(t, &acc)

			pr688AddChunk(t, &acc, pr688FinishedChunk)
			pr688AssertFinishedToolCall(t, &acc)

			pr688AddChunk(t, &acc, pr688FinishedChunk)
			pr688AssertNoFinishedToolCall(t, &acc)
		})

		t.Run(tc.name+"/after active tool call", func(t *testing.T) {
			var acc openai.ChatCompletionAccumulator

			pr688AddChunk(t, &acc, pr688ToolCallChunk)
			pr688AssertNoFinishedToolCall(t, &acc)

			pr688AddChunk(t, &acc, tc.chunk)
			pr688AssertFinishedToolCall(t, &acc)

			pr688AddChunk(t, &acc, pr688FinishedChunk)
			pr688AssertNoFinishedToolCall(t, &acc)
		})
	}
}

func pr688AddChunk(t *testing.T, acc *openai.ChatCompletionAccumulator, raw string) {
	t.Helper()

	var chunk openai.ChatCompletionChunk
	if err := chunk.UnmarshalJSON([]byte(raw)); err != nil {
		t.Fatalf("failed to unmarshal chunk: %v", err)
	}
	if !acc.AddChunk(chunk) {
		t.Fatal("AddChunk returned false")
	}
}

func pr688AssertNoFinishedToolCall(t *testing.T, acc *openai.ChatCompletionAccumulator) {
	t.Helper()

	if toolCall, ok := acc.JustFinishedToolCall(); ok {
		t.Fatalf("JustFinishedToolCall returned unexpected tool call: %+v", toolCall)
	}
}

func pr688AssertFinishedToolCall(t *testing.T, acc *openai.ChatCompletionAccumulator) {
	t.Helper()

	toolCall, ok := acc.JustFinishedToolCall()
	if !ok {
		t.Fatal("JustFinishedToolCall did not return the completed tool call")
	}
	if toolCall.Index != 0 {
		t.Fatalf("tool call index: expected 0, got %d", toolCall.Index)
	}
	if toolCall.ID != "call_123" {
		t.Fatalf("tool call ID: expected %q, got %q", "call_123", toolCall.ID)
	}
	if toolCall.Name != "test_func" {
		t.Fatalf("tool call name: expected %q, got %q", "test_func", toolCall.Name)
	}
	if toolCall.Arguments != "{}" {
		t.Fatalf("tool call arguments: expected %q, got %q", "{}", toolCall.Arguments)
	}
}

package openai_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	openai "github.com/openai/openai-go/v3"
)

func TestAccumulatorRejectsInvalidIndicesWithoutMutation(t *testing.T) {
	for _, test := range []struct {
		name string
		raw  string
	}{
		{
			name: "negative choice",
			raw:  `{"id":"test","choices":[{"index":0,"delta":{"content":" mutated"}},{"index":-1,"delta":{}}]}`,
		},
		{
			name: "choice above limit",
			raw:  `{"id":"test","choices":[{"index":0,"delta":{"content":" mutated"}},{"index":128,"delta":{}}]}`,
		},
		{
			name: "maximum int64 choice",
			raw:  `{"id":"test","choices":[{"index":0,"delta":{"content":" mutated"}},{"index":9223372036854775807,"delta":{}}]}`,
		},
		{
			name: "negative tool call",
			raw:  `{"id":"test","choices":[{"index":0,"delta":{"content":" mutated","tool_calls":[{"index":-2}]}}]}`,
		},
		{
			name: "sparse tool call",
			raw:  `{"id":"test","choices":[{"index":0,"delta":{"content":" mutated","tool_calls":[{"index":128}]}}]}`,
		},
		{
			name: "maximum int64 tool call",
			raw:  `{"id":"test","choices":[{"index":0,"delta":{"content":" mutated","tool_calls":[{"index":9223372036854775807}]}}]}`,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			var acc openai.ChatCompletionAccumulator
			addAccumulatorChunk(t, &acc, `{"id":"test","choices":[{"index":0,"delta":{"content":"safe"}}]}`)
			addAccumulatorChunk(t, &acc, `{"id":"test","choices":[{"index":0,"delta":{},"finish_reason":"stop"}]}`)

			before, err := json.Marshal(acc.ChatCompletion)
			if err != nil {
				t.Fatalf("failed to snapshot accumulator: %v", err)
			}

			var chunk openai.ChatCompletionChunk
			if unmarshalErr := chunk.UnmarshalJSON([]byte(test.raw)); unmarshalErr != nil {
				t.Fatalf("failed to unmarshal invalid chunk: %v", unmarshalErr)
			}
			if acc.AddChunk(chunk) {
				t.Fatal("AddChunk accepted an invalid index")
			}

			after, err := json.Marshal(acc.ChatCompletion)
			if err != nil {
				t.Fatalf("failed to snapshot accumulator after rejection: %v", err)
			}
			if !bytes.Equal(after, before) {
				t.Fatalf("rejected chunk mutated the accumulator:\nbefore: %s\n after: %s", before, after)
			}
			if content, ok := acc.JustFinishedContent(); ok {
				t.Fatalf("rejected chunk retained the finished event: got %q", content)
			}
		})
	}
}

func TestAccumulatorRejectedChunksDoNotReplayFinishedToolCall(t *testing.T) {
	var acc openai.ChatCompletionAccumulator
	addAccumulatorChunk(t, &acc, `{"id":"test","choices":[{"index":0,"delta":{"tool_calls":[{"index":0,"id":"call_0","type":"function","function":{"name":"lookup","arguments":"{}"}}]}}]}`)
	addAccumulatorChunk(t, &acc, `{"id":"test","choices":[{"index":0,"delta":{},"finish_reason":"tool_calls"}]}`)

	if toolCall, ok := acc.JustFinishedToolCall(); !ok || toolCall.ID != "call_0" {
		t.Fatalf("expected initial finished tool call, got %+v, %v", toolCall, ok)
	}
	if toolCall, ok := acc.JustFinishedToolCallForChoice(0); !ok || toolCall.ID != "call_0" {
		t.Fatalf("expected initial choice-specific finished tool call, got %+v, %v", toolCall, ok)
	}

	var rejected openai.ChatCompletionChunk
	if err := rejected.UnmarshalJSON([]byte(`{"id":"test","choices":[{"index":-1,"delta":{}}]}`)); err != nil {
		t.Fatalf("failed to unmarshal rejected chunk: %v", err)
	}
	for attempt := 1; attempt <= 2; attempt++ {
		if acc.AddChunk(rejected) {
			t.Fatalf("attempt %d: AddChunk accepted an invalid index", attempt)
		}
		if toolCall, ok := acc.JustFinishedToolCall(); ok {
			t.Fatalf("attempt %d: rejected chunk replayed finished tool call %+v", attempt, toolCall)
		}
		if toolCall, ok := acc.JustFinishedToolCallForChoice(0); ok {
			t.Fatalf("attempt %d: rejected chunk replayed choice-specific finished tool call %+v", attempt, toolCall)
		}
	}
}

func TestAccumulatorAcceptsMaximumChoiceIndex(t *testing.T) {
	var acc openai.ChatCompletionAccumulator
	addAccumulatorChunk(t, &acc, `{"id":"test","choices":[{"index":127,"delta":{"tool_calls":[{"index":0,"id":"call_0","type":"function","function":{"name":"lookup","arguments":"{}"}}]}}]}`)

	if len(acc.Choices) != 128 {
		t.Fatalf("choices: expected 128 slots, got %d", len(acc.Choices))
	}
	toolCalls := acc.Choices[127].Message.ToolCalls
	if len(toolCalls) != 1 {
		t.Fatalf("tool calls: expected 1 slot, got %d", len(toolCalls))
	}
	if toolCalls[0].ID != "call_0" || toolCalls[0].Function.Name != "lookup" {
		t.Fatalf("unexpected tool call at accepted choice boundary: %+v", toolCalls[0])
	}
}

func TestAccumulatorAcceptsDenseToolCallsBeyondChoiceLimit(t *testing.T) {
	toolCalls := make([]map[string]any, 129)
	for i := range toolCalls {
		toolCalls[i] = map[string]any{
			"index": i,
			"id":    fmt.Sprintf("call_%d", i),
			"type":  "function",
			"function": map[string]string{
				"name":      "lookup",
				"arguments": "{}",
			},
		}
	}
	raw, err := json.Marshal(map[string]any{
		"id": "test",
		"choices": []any{map[string]any{
			"index": 0,
			"delta": map[string]any{"tool_calls": toolCalls},
		}},
	})
	if err != nil {
		t.Fatalf("failed to marshal dense tool calls: %v", err)
	}

	var acc openai.ChatCompletionAccumulator
	addAccumulatorChunk(t, &acc, string(raw))

	accumulated := acc.Choices[0].Message.ToolCalls
	if len(accumulated) != 129 {
		t.Fatalf("tool calls: expected 129 slots, got %d", len(accumulated))
	}
	if accumulated[128].ID != "call_128" || accumulated[128].Function.Name != "lookup" {
		t.Fatalf("unexpected final dense tool call: %+v", accumulated[128])
	}
}

func addAccumulatorChunk(t *testing.T, acc *openai.ChatCompletionAccumulator, raw string) {
	t.Helper()

	var chunk openai.ChatCompletionChunk
	if err := chunk.UnmarshalJSON([]byte(raw)); err != nil {
		t.Fatalf("failed to unmarshal chunk: %v", err)
	}
	if !acc.AddChunk(chunk) {
		t.Fatal("AddChunk returned false for a valid chunk")
	}
}

package openai_test

import (
	"bytes"
	"encoding/json"
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
			name: "tool call above limit",
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
			if content, ok := acc.JustFinishedContent(); !ok || content != "safe" {
				t.Fatalf("rejected chunk changed the finished event: got %q, %v", content, ok)
			}
		})
	}
}

func TestAccumulatorAcceptsMaximumIndex(t *testing.T) {
	var acc openai.ChatCompletionAccumulator
	addAccumulatorChunk(t, &acc, `{"id":"test","choices":[{"index":127,"delta":{"tool_calls":[{"index":127,"id":"call_127","type":"function","function":{"name":"lookup","arguments":"{}"}}]}}]}`)

	if len(acc.Choices) != 128 {
		t.Fatalf("choices: expected 128 slots, got %d", len(acc.Choices))
	}
	toolCalls := acc.Choices[127].Message.ToolCalls
	if len(toolCalls) != 128 {
		t.Fatalf("tool calls: expected 128 slots, got %d", len(toolCalls))
	}
	if toolCalls[127].ID != "call_127" || toolCalls[127].Function.Name != "lookup" {
		t.Fatalf("unexpected tool call at accepted boundary: %+v", toolCalls[127])
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

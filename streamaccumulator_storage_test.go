package openai

import (
	"strings"
	"testing"
	"unsafe"
)

func TestAccumulatorReleasesClearedStringStorageBeforeNextChunk(t *testing.T) {
	text := strings.Repeat("x", 32<<10)
	tests := []struct {
		name    string
		initial ChatCompletionChunk
		value   func(*ChatCompletionAccumulator) *string
		state   func(*ChatCompletionAccumulator) *chatCompletionString
	}{
		{
			name:    "content",
			initial: storageTestChunk(ChatCompletionChunkChoiceDelta{Content: text}),
			value: func(acc *ChatCompletionAccumulator) *string {
				return &acc.Choices[0].Message.Content
			},
			state: func(acc *ChatCompletionAccumulator) *chatCompletionString {
				return &acc.stringState.choices[0].content
			},
		},
		{
			name:    "refusal",
			initial: storageTestChunk(ChatCompletionChunkChoiceDelta{Refusal: text}),
			value: func(acc *ChatCompletionAccumulator) *string {
				return &acc.Choices[0].Message.Refusal
			},
			state: func(acc *ChatCompletionAccumulator) *chatCompletionString {
				return &acc.stringState.choices[0].refusal
			},
		},
		{
			name: "tool_name",
			initial: storageTestChunk(ChatCompletionChunkChoiceDelta{
				ToolCalls: []ChatCompletionChunkChoiceDeltaToolCall{{
					Function: ChatCompletionChunkChoiceDeltaToolCallFunction{Name: text},
				}},
			}),
			value: func(acc *ChatCompletionAccumulator) *string {
				return &acc.Choices[0].Message.ToolCalls[0].Function.Name
			},
			state: func(acc *ChatCompletionAccumulator) *chatCompletionString {
				return &acc.stringState.choices[0].toolCalls[0].name
			},
		},
		{
			name: "tool_arguments",
			initial: storageTestChunk(ChatCompletionChunkChoiceDelta{
				ToolCalls: []ChatCompletionChunkChoiceDeltaToolCall{{
					Function: ChatCompletionChunkChoiceDeltaToolCallFunction{Arguments: text},
				}},
			}),
			value: func(acc *ChatCompletionAccumulator) *string {
				return &acc.Choices[0].Message.ToolCalls[0].Function.Arguments
			},
			state: func(acc *ChatCompletionAccumulator) *chatCompletionString {
				return &acc.stringState.choices[0].toolCalls[0].arguments
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var acc ChatCompletionAccumulator
			if !acc.AddChunk(test.initial) {
				t.Fatal("AddChunk rejected the initial chunk")
			}
			if capacity := test.state(&acc).builder.Cap(); capacity < len(text) {
				t.Fatalf("initial builder capacity: got %d, want at least %d", capacity, len(text))
			}

			*test.value(&acc) = ""
			if !acc.AddChunk(ChatCompletionChunk{ID: test.initial.ID}) {
				t.Fatal("AddChunk rejected the chunk after the public string was cleared")
			}
			state := test.state(&acc)
			if state.builder.Cap() != 0 || state.published != "" {
				t.Fatalf("cleared string storage was retained: capacity %d, published length %d", state.builder.Cap(), len(state.published))
			}
		})
	}
}

func TestAccumulatorCanonicalizesEqualPublicStringBacking(t *testing.T) {
	var acc ChatCompletionAccumulator
	if !acc.AddChunk(storageTestChunk(ChatCompletionChunkChoiceDelta{Content: strings.Repeat("x", 32<<10)})) {
		t.Fatal("AddChunk rejected the initial chunk")
	}

	published := acc.stringState.choices[0].content.published
	clone := strings.Clone(acc.Choices[0].Message.Content)
	if unsafe.StringData(clone) == unsafe.StringData(published) {
		t.Fatal("strings.Clone did not produce distinct test backing storage")
	}
	acc.Choices[0].Message.Content = clone

	if !acc.AddChunk(ChatCompletionChunk{ID: "chatcmpl-storage-reconciliation"}) {
		t.Fatal("AddChunk rejected the chunk after an equal public string replacement")
	}
	if unsafe.StringData(acc.Choices[0].Message.Content) != unsafe.StringData(published) {
		t.Fatal("AddChunk retained distinct backing storage for an equal public string")
	}
}

func storageTestChunk(delta ChatCompletionChunkChoiceDelta) ChatCompletionChunk {
	return ChatCompletionChunk{
		ID: "chatcmpl-storage-reconciliation",
		Choices: []ChatCompletionChunkChoice{{
			Delta: delta,
		}},
	}
}

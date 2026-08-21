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

func TestAccumulatorTracksOnlyPopulatedSparseState(t *testing.T) {
	chunk := storageTestChunk(ChatCompletionChunkChoiceDelta{Content: "x"})
	chunk.Choices[0].Index = maxChatCompletionAccumulatorStructuralSlots - 1

	var acc ChatCompletionAccumulator
	if !acc.AddChunk(chunk) {
		t.Fatal("AddChunk rejected the sparse choice")
	}
	if len(acc.stringState.activeChoices) != 1 ||
		acc.stringState.activeChoices[0] != maxChatCompletionAccumulatorStructuralSlots-1 {
		t.Fatalf("active choices = %v, want only the populated sparse choice", acc.stringState.activeChoices)
	}

	empty := ChatCompletionChunk{ID: chunk.ID}
	for range maxChatCompletionAccumulatorChunks - 1 {
		if !acc.AddChunk(empty) {
			t.Fatal("AddChunk rejected an empty chunk within the documented budget")
		}
	}
	if len(acc.stringState.activeChoices) != 1 {
		t.Fatalf("empty chunks expanded active state to %d choices, want 1", len(acc.stringState.activeChoices))
	}
}

func TestAccumulatorBoundsLogprobReconciliationWork(t *testing.T) {
	chunk := storageTestChunk(ChatCompletionChunkChoiceDelta{})
	chunk.Choices[0].Logprobs.Content = []ChatCompletionTokenLogprob{{Token: "initial"}}

	var acc ChatCompletionAccumulator
	if !acc.AddChunk(chunk) {
		t.Fatal("AddChunk rejected the initial logprob")
	}
	acc.logprobState.reconciliationWork = maxChatCompletionAccumulatorLogprobWork - 1
	acc.Choices[0].Logprobs.Content = append([]ChatCompletionTokenLogprob(nil), acc.Choices[0].Logprobs.Content...)
	if acc.AddChunk(ChatCompletionChunk{ID: chunk.ID}) {
		t.Fatal("AddChunk accepted work beyond the logprob reconciliation budget")
	}
	if got := len(acc.Choices[0].Logprobs.Content); got != 1 {
		t.Fatalf("rejected chunk changed logprobs to length %d, want 1", got)
	}
}

func TestAccumulatorLogprobStreamingUsesIncrementalAccounting(t *testing.T) {
	chunk := storageTestChunk(ChatCompletionChunkChoiceDelta{})
	chunk.Choices[0].Logprobs.Content = []ChatCompletionTokenLogprob{{Token: "x"}}

	var acc ChatCompletionAccumulator
	const chunkCount = 12_000
	for i := range chunkCount {
		if !acc.AddChunk(chunk) {
			t.Fatalf("AddChunk rejected logprob chunk %d within the documented budgets", i+1)
		}
	}
	if got := len(acc.Choices[0].Logprobs.Content); got != chunkCount {
		t.Fatalf("accumulated logprobs = %d, want %d", got, chunkCount)
	}
	if acc.logprobState.reconciliationWork != 0 {
		t.Fatalf("normal streaming used %d public-replacement reconciliation steps, want 0", acc.logprobState.reconciliationWork)
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

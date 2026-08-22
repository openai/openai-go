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
			if capacity := cap(test.state(&acc).buffer); capacity < len(text) {
				t.Fatalf("initial buffer capacity: got %d, want at least %d", capacity, len(text))
			}

			*test.value(&acc) = ""
			if !acc.AddChunk(ChatCompletionChunk{ID: test.initial.ID}) {
				t.Fatal("AddChunk rejected the chunk after the public string was cleared")
			}
			state := test.state(&acc)
			if cap(state.buffer) != 0 || state.published != "" {
				t.Fatalf("cleared string storage was retained: capacity %d, published length %d", cap(state.buffer), len(state.published))
			}
		})
	}
}

func TestAccumulatorBoundsRetainedTextBufferCapacity(t *testing.T) {
	const fragmentBytes = 1 << 10
	fragment := strings.Repeat("x", fragmentBytes)
	chunk := storageTestChunk(ChatCompletionChunkChoiceDelta{Content: fragment})

	var acc ChatCompletionAccumulator
	for range maxChatCompletionAccumulatorTextBytes / fragmentBytes {
		if !acc.AddChunk(chunk) {
			t.Fatal("AddChunk rejected text within the documented logical and retained-capacity budgets")
		}
	}
	if got := len(acc.Choices[0].Message.Content); got != maxChatCompletionAccumulatorTextBytes {
		t.Fatalf("content length = %d, want %d", got, maxChatCompletionAccumulatorTextBytes)
	}
	if got := cap(acc.stringState.choices[0].content.buffer); got > maxChatCompletionAccumulatorTextCapacity {
		t.Fatalf("retained buffer capacity = %d, limit %d", got, maxChatCompletionAccumulatorTextCapacity)
	}

	capacity := cap(acc.stringState.choices[0].content.buffer)
	if acc.AddChunk(chunk) {
		t.Fatal("AddChunk accepted text beyond the documented logical budget")
	}
	if got := cap(acc.stringState.choices[0].content.buffer); got != capacity {
		t.Fatalf("rejected chunk changed retained capacity from %d to %d", capacity, got)
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

func TestAccumulatorBoundsEqualPublicStringComparisonWork(t *testing.T) {
	const contentBytes = 1 << 20
	chunk := storageTestChunk(ChatCompletionChunkChoiceDelta{Content: strings.Repeat("x", contentBytes)})

	var acc ChatCompletionAccumulator
	if !acc.AddChunk(chunk) {
		t.Fatal("AddChunk rejected the initial chunk")
	}
	clone := strings.Clone(acc.Choices[0].Message.Content)
	if unsafe.StringData(clone) == unsafe.StringData(acc.stringState.choices[0].content.published) {
		t.Fatal("strings.Clone did not produce distinct test backing storage")
	}
	acc.Choices[0].Message.Content = clone

	empty := ChatCompletionChunk{ID: chunk.ID}
	const activeChoicePasses = 6
	requiredWork := activeChoicePasses + len(empty.ID) + 2*contentBytes
	acc.reconciliationWork = maxChatCompletionAccumulatorReconcileWork - requiredWork + 1
	beforeWork := acc.reconciliationWork
	if acc.AddChunk(empty) {
		t.Fatal("AddChunk accepted equal-string comparison work beyond the documented budget")
	}
	if acc.reconciliationWork != beforeWork {
		t.Fatal("rejected chunk changed the reconciliation budget")
	}
	if unsafe.StringData(acc.Choices[0].Message.Content) != unsafe.StringData(clone) {
		t.Fatal("rejected chunk changed the public string backing")
	}
}

func TestAccumulatorCanonicalizesEqualPublicMetadataBacking(t *testing.T) {
	chunk := storageTestChunk(ChatCompletionChunkChoiceDelta{
		Role: "assistant",
		ToolCalls: []ChatCompletionChunkChoiceDeltaToolCall{{
			ID:   "tool-id",
			Type: "function",
		}},
	})
	chunk.Model = "model"
	chunk.SystemFingerprint = "fingerprint"
	chunk.ServiceTier = "default"
	chunk.Object = chunk.Object.Default()
	chunk.Choices[0].FinishReason = "stop"

	var acc ChatCompletionAccumulator
	if !acc.AddChunk(chunk) {
		t.Fatal("AddChunk rejected the initial metadata")
	}
	replacements := []struct {
		name   string
		source *byte
		get    func() string
	}{
		{name: "id", source: replaceWithEqualSubstring(&acc.ID), get: func() string { return acc.ID }},
		{name: "model", source: replaceWithEqualSubstring(&acc.Model), get: func() string { return acc.Model }},
		{name: "fingerprint", source: replaceWithEqualSubstring(&acc.SystemFingerprint), get: func() string { return acc.SystemFingerprint }},
		{name: "service_tier", source: replaceWithEqualSubstring(&acc.ServiceTier), get: func() string { return string(acc.ServiceTier) }},
		{name: "object", source: replaceWithEqualSubstring(&acc.Object), get: func() string { return string(acc.Object) }},
		{name: "finish_reason", source: replaceWithEqualSubstring(&acc.Choices[0].FinishReason), get: func() string { return acc.Choices[0].FinishReason }},
		{name: "role", source: replaceWithEqualSubstring(&acc.Choices[0].Message.Role), get: func() string { return string(acc.Choices[0].Message.Role) }},
		{name: "tool_id", source: replaceWithEqualSubstring(&acc.Choices[0].Message.ToolCalls[0].ID), get: func() string { return acc.Choices[0].Message.ToolCalls[0].ID }},
		{name: "tool_type", source: replaceWithEqualSubstring(&acc.Choices[0].Message.ToolCalls[0].Type), get: func() string { return acc.Choices[0].Message.ToolCalls[0].Type }},
	}

	if !acc.AddChunk(chunk) {
		t.Fatal("AddChunk rejected equal public metadata replacements")
	}
	for _, replacement := range replacements {
		if unsafe.StringData(replacement.get()) == replacement.source {
			t.Errorf("%s still retains replacement backing", replacement.name)
		}
	}

	emptyBacking := strings.Repeat("x", 1<<20)[:0]
	acc.SystemFingerprint = emptyBacking
	if !acc.AddChunk(chunk) {
		t.Fatal("AddChunk rejected an empty public metadata replacement")
	}
	if unsafe.StringData(acc.SystemFingerprint) == unsafe.StringData(emptyBacking) {
		t.Error("empty metadata still retains replacement backing")
	}
}

func replaceWithEqualSubstring[T ~string](dst *T) *byte {
	value := string(*dst)
	backing := value + strings.Repeat("x", 1<<20)
	replacement := backing[:len(value)]
	*dst = T(replacement)
	return unsafe.StringData(replacement)
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

func TestAccumulatorTracksOnlyPopulatedToolState(t *testing.T) {
	chunk := storageTestChunk(ChatCompletionChunkChoiceDelta{
		ToolCalls: []ChatCompletionChunkChoiceDeltaToolCall{{
			Index: maxChatCompletionAccumulatorStructuralSlots - 2,
			Function: ChatCompletionChunkChoiceDeltaToolCallFunction{
				Name: "tool",
			},
		}},
	})

	var acc ChatCompletionAccumulator
	if !acc.AddChunk(chunk) {
		t.Fatal("AddChunk rejected the sparse tool call")
	}
	if acc.stringState.activeTools != 1 {
		t.Fatalf("active tools = %d, want 1", acc.stringState.activeTools)
	}
	acc.Choices[0].Message.ToolCalls = nil
	if !acc.AddChunk(ChatCompletionChunk{ID: chunk.ID}) {
		t.Fatal("AddChunk rejected the truncated tool state")
	}
	if acc.stringState.activeTools != 0 {
		t.Fatalf("active tools after truncation = %d, want 0", acc.stringState.activeTools)
	}
}

func TestAccumulatorBoundsLogprobReconciliationWork(t *testing.T) {
	chunk := storageTestChunk(ChatCompletionChunkChoiceDelta{})
	chunk.Choices[0].Logprobs.Content = []ChatCompletionTokenLogprob{{Token: "initial"}}

	var acc ChatCompletionAccumulator
	if !acc.AddChunk(chunk) {
		t.Fatal("AddChunk rejected the initial logprob")
	}
	acc.reconciliationWork = maxChatCompletionAccumulatorReconcileWork - 12
	acc.Choices[0].Logprobs.Content = append([]ChatCompletionTokenLogprob(nil), acc.Choices[0].Logprobs.Content...)
	if acc.AddChunk(ChatCompletionChunk{ID: chunk.ID}) {
		t.Fatal("AddChunk accepted copy work beyond the logprob reconciliation budget")
	}
	if got := len(acc.Choices[0].Logprobs.Content); got != 1 {
		t.Fatalf("rejected chunk changed logprobs to length %d, want 1", got)
	}
}

func TestAccumulatorChargesAllLogprobReconciliationPasses(t *testing.T) {
	chunk := storageTestChunk(ChatCompletionChunkChoiceDelta{})
	chunk.Choices[0].Logprobs.Content = []ChatCompletionTokenLogprob{{Token: "initial"}}

	var acc ChatCompletionAccumulator
	if !acc.AddChunk(chunk) {
		t.Fatal("AddChunk rejected the initial logprob")
	}
	before := acc.reconciliationWork
	acc.Choices[0].Logprobs.Content = []ChatCompletionTokenLogprob{}
	if !acc.AddChunk(ChatCompletionChunk{ID: chunk.ID}) {
		t.Fatal("AddChunk rejected an empty whole-slice replacement")
	}
	wantWork := 10 + len(chunk.ID) // six normal passes, two staged-slice passes, measurement, sparse commit, and incoming ID
	if got := acc.reconciliationWork - before; got != wantWork {
		t.Fatalf("reconciliation work = %d, want %d", got, wantWork)
	}
}

func TestAccumulatorChargesPublicReplacementCopyWork(t *testing.T) {
	chunk := storageTestChunk(ChatCompletionChunkChoiceDelta{Content: "initial"})
	chunk.Model = "initial"

	var acc ChatCompletionAccumulator
	if !acc.AddChunk(chunk) {
		t.Fatal("AddChunk rejected the initial strings")
	}
	const replacementBytes = 1_024
	acc.Choices[0].Message.Content = strings.Repeat("c", replacementBytes)
	acc.Model = strings.Repeat("m", replacementBytes)
	chunk.Choices = nil
	chunk.Model = acc.Model
	before := acc.reconciliationWork
	if !acc.AddChunk(chunk) {
		t.Fatal("AddChunk rejected supported public replacements")
	}
	const normalPasses = 6
	wantWork := normalPasses + len(chunk.ID) + 5*replacementBytes
	if got := acc.reconciliationWork - before; got != wantWork {
		t.Fatalf("reconciliation work = %d, want %d", got, wantWork)
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
	wantWork := chunkCount*(chatCompletionAccumulatorChoiceWork+len(chunk.ID)) + (chunkCount-1)*7
	if acc.reconciliationWork != wantWork {
		t.Fatalf("normal streaming used %d reconciliation steps, want %d", acc.reconciliationWork, wantWork)
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

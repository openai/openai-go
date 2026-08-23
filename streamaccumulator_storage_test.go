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

func TestAccumulatorTextReplacementRetainsGeometricHeadroom(t *testing.T) {
	chunk := storageTestChunk(ChatCompletionChunkChoiceDelta{Content: "x"})
	var acc ChatCompletionAccumulator
	if !acc.AddChunk(chunk) {
		t.Fatal("AddChunk rejected the initial chunk")
	}

	replacementBytes := maxChatCompletionAccumulatorTextBytes/2 + 1
	acc.Choices[0].Message.Content = strings.Repeat("x", replacementBytes)
	chunk.Choices[0].Delta.Content = "y"
	if !acc.AddChunk(chunk) {
		t.Fatal("AddChunk rejected text after a large supported public replacement")
	}
	if got := cap(acc.stringState.choices[0].content.buffer); got != maxChatCompletionAccumulatorTextBytes {
		t.Fatalf("retained text capacity = %d, want geometric headroom %d", got, maxChatCompletionAccumulatorTextBytes)
	}
}

func TestAccumulatorLogprobsRetainGeometricHeadroomAtFinalGrowth(t *testing.T) {
	logprobSize := int(unsafe.Sizeof(ChatCompletionTokenLogprob{}))
	maxLogprobs := maxChatCompletionAccumulatorLogprobBytes / logprobSize
	chunk := storageTestChunk(ChatCompletionChunkChoiceDelta{})
	chunk.Choices[0].Logprobs.Content = make([]ChatCompletionTokenLogprob, 1<<16)

	var acc ChatCompletionAccumulator
	if !acc.AddChunk(chunk) {
		t.Fatal("AddChunk rejected logprobs at the last power-of-two capacity")
	}
	chunk.Choices[0].Logprobs.Content = []ChatCompletionTokenLogprob{{}}
	if !acc.AddChunk(chunk) {
		t.Fatal("AddChunk rejected the final bounded geometric logprob growth")
	}
	if got := cap(acc.Choices[0].Logprobs.Content); got != maxLogprobs {
		t.Fatalf("retained logprob capacity = %d, want bounded headroom %d", got, maxLogprobs)
	}
}

func TestAccumulatorBoundsDuplicateMetadataAssignmentWork(t *testing.T) {
	initial := storageTestChunk(ChatCompletionChunkChoiceDelta{})
	var acc ChatCompletionAccumulator
	if !acc.AddChunk(initial) {
		t.Fatal("AddChunk rejected the initial chunk")
	}

	const metadataBytes = 1 << 20
	next := storageTestChunk(ChatCompletionChunkChoiceDelta{})
	next.Choices = append(next.Choices, next.Choices[0])
	next.Choices[0].FinishReason = strings.Repeat("a", metadataBytes-1) + "x"
	next.Choices[1].FinishReason = strings.Repeat("a", metadataBytes-1) + "y"
	acc.reconciliationWork = maxChatCompletionAccumulatorReconcileWork - 3*metadataBytes
	before := acc.reconciliationWork
	if acc.AddChunk(next) {
		t.Fatal("AddChunk accepted duplicate metadata work beyond the cumulative budget")
	}
	if acc.reconciliationWork != before || acc.Choices[0].FinishReason != "" {
		t.Fatal("rejected duplicate metadata changed the accumulator")
	}
}

func TestAccumulatorAccountsMetadataAssignmentAfterPublicReplacement(t *testing.T) {
	const metadataBytes = 256 << 10
	original := strings.Repeat("a", metadataBytes-1) + "x"
	replacement := strings.Repeat("a", metadataBytes-1) + "y"
	tests := []struct {
		name     string
		initial  func() ChatCompletionChunk
		get      func(*ChatCompletionAccumulator) string
		set      func(*ChatCompletionAccumulator, string)
		setChunk func(*ChatCompletionChunk, string)
	}{
		{
			name: "model",
			initial: func() ChatCompletionChunk {
				chunk := storageTestChunk(ChatCompletionChunkChoiceDelta{})
				chunk.Model = original
				return chunk
			},
			get:      func(acc *ChatCompletionAccumulator) string { return acc.Model },
			set:      func(acc *ChatCompletionAccumulator, value string) { acc.Model = value },
			setChunk: func(chunk *ChatCompletionChunk, value string) { chunk.Model = value },
		},
		{
			name: "tool_id",
			initial: func() ChatCompletionChunk {
				return storageTestChunk(ChatCompletionChunkChoiceDelta{ToolCalls: []ChatCompletionChunkChoiceDeltaToolCall{{ID: original}}})
			},
			get: func(acc *ChatCompletionAccumulator) string {
				return acc.Choices[0].Message.ToolCalls[0].ID
			},
			set: func(acc *ChatCompletionAccumulator, value string) {
				acc.Choices[0].Message.ToolCalls[0].ID = value
			},
			setChunk: func(chunk *ChatCompletionChunk, value string) {
				chunk.Choices[0].Delta.ToolCalls[0].ID = value
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			initial := test.initial()
			var acc ChatCompletionAccumulator
			if !acc.AddChunk(initial) {
				t.Fatal("AddChunk rejected the initial metadata")
			}
			published := test.get(&acc)
			test.set(&acc, replacement)
			next := test.initial()
			test.setChunk(&next, published)

			acc.reconciliationWork = maxChatCompletionAccumulatorReconcileWork - 4*metadataBytes
			beforeWork := acc.reconciliationWork
			if acc.AddChunk(next) {
				t.Fatal("AddChunk accepted metadata work beyond the post-reconciliation budget")
			}
			if acc.reconciliationWork != beforeWork || test.get(&acc) != replacement {
				t.Fatal("rejected metadata assignment changed the accumulator")
			}
		})
	}
}

func TestAccumulatorAccountsGrowthAfterPublicTextReplacement(t *testing.T) {
	const replacementBytes = 256 << 10
	replacement := strings.Repeat("r", replacementBytes)
	tests := []struct {
		name    string
		initial ChatCompletionChunkChoiceDelta
		next    ChatCompletionChunkChoiceDelta
		value   func(*ChatCompletionAccumulator) *string
	}{
		{
			name:    "content",
			initial: ChatCompletionChunkChoiceDelta{Content: "initial"},
			next:    ChatCompletionChunkChoiceDelta{Content: "x"},
			value: func(acc *ChatCompletionAccumulator) *string {
				return &acc.Choices[0].Message.Content
			},
		},
		{
			name:    "refusal",
			initial: ChatCompletionChunkChoiceDelta{Refusal: "initial"},
			next:    ChatCompletionChunkChoiceDelta{Refusal: "x"},
			value: func(acc *ChatCompletionAccumulator) *string {
				return &acc.Choices[0].Message.Refusal
			},
		},
		{
			name: "tool_name",
			initial: ChatCompletionChunkChoiceDelta{ToolCalls: []ChatCompletionChunkChoiceDeltaToolCall{{
				Function: ChatCompletionChunkChoiceDeltaToolCallFunction{Name: "initial"},
			}}},
			next: ChatCompletionChunkChoiceDelta{ToolCalls: []ChatCompletionChunkChoiceDeltaToolCall{{
				Function: ChatCompletionChunkChoiceDeltaToolCallFunction{Name: "x"},
			}}},
			value: func(acc *ChatCompletionAccumulator) *string {
				return &acc.Choices[0].Message.ToolCalls[0].Function.Name
			},
		},
		{
			name: "tool_arguments",
			initial: ChatCompletionChunkChoiceDelta{ToolCalls: []ChatCompletionChunkChoiceDeltaToolCall{{
				Function: ChatCompletionChunkChoiceDeltaToolCallFunction{Arguments: "initial"},
			}}},
			next: ChatCompletionChunkChoiceDelta{ToolCalls: []ChatCompletionChunkChoiceDeltaToolCall{{
				Function: ChatCompletionChunkChoiceDeltaToolCallFunction{Arguments: "x"},
			}}},
			value: func(acc *ChatCompletionAccumulator) *string {
				return &acc.Choices[0].Message.ToolCalls[0].Function.Arguments
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			initial := storageTestChunk(test.initial)
			var acc ChatCompletionAccumulator
			if !acc.AddChunk(initial) {
				t.Fatal("AddChunk rejected the initial text")
			}
			*test.value(&acc) = replacement
			next := storageTestChunk(test.next)
			acc.reconciliationWork = maxChatCompletionAccumulatorReconcileWork - (2*replacementBytes + 4_096)
			beforeWork := acc.reconciliationWork
			if acc.AddChunk(next) {
				t.Fatal("AddChunk accepted an unbudgeted post-reconciliation growth copy")
			}
			if acc.reconciliationWork != beforeWork || *test.value(&acc) != replacement {
				t.Fatal("rejected text append changed the accumulator")
			}
		})
	}
}

func TestAccumulatorAccountsGrowthAfterPublicLogprobReplacement(t *testing.T) {
	initial := storageTestChunk(ChatCompletionChunkChoiceDelta{})
	initial.Choices[0].Logprobs.Content = []ChatCompletionTokenLogprob{{Token: "initial"}}
	var acc ChatCompletionAccumulator
	if !acc.AddChunk(initial) {
		t.Fatal("AddChunk rejected the initial logprob")
	}

	const replacementBytes = 1 << 20
	logprobSize := int(unsafe.Sizeof(ChatCompletionTokenLogprob{}))
	replacementCount := replacementBytes / logprobSize
	replacement := make([]ChatCompletionTokenLogprob, replacementCount)
	acc.Choices[0].Logprobs.Content = replacement
	next := storageTestChunk(ChatCompletionChunkChoiceDelta{})
	next.Choices[0].Logprobs.Content = []ChatCompletionTokenLogprob{{Token: "x"}}
	oldProjectionWork := replacementCount*logprobSize + 2*replacementCount
	acc.reconciliationWork = maxChatCompletionAccumulatorReconcileWork - oldProjectionWork - 4_096
	beforeWork := acc.reconciliationWork
	if acc.AddChunk(next) {
		t.Fatal("AddChunk accepted an unbudgeted post-reconciliation logprob growth copy")
	}
	if acc.reconciliationWork != beforeWork || len(acc.Choices[0].Logprobs.Content) != replacementCount {
		t.Fatal("rejected logprob append changed the accumulator")
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
	chunk.Choices[0].Index = maxStreamAccumulatorChoiceIndex

	var acc ChatCompletionAccumulator
	if !acc.AddChunk(chunk) {
		t.Fatal("AddChunk rejected the sparse choice")
	}
	if len(acc.stringState.activeChoices) != 1 ||
		acc.stringState.activeChoices[0] != maxStreamAccumulatorChoiceIndex {
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
			Index: maxStreamAccumulatorToolCallGrowth - 1,
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

func TestAccumulatorValueCopyLogprobStateIsolated(t *testing.T) {
	initial := storageTestChunk(ChatCompletionChunkChoiceDelta{})
	initial.Choices[0].Index = 1
	initial.Choices[0].Logprobs.Content = []ChatCompletionTokenLogprob{{Token: "existing"}}

	var original ChatCompletionAccumulator
	if !original.AddChunk(initial) {
		t.Fatal("AddChunk rejected the initial sparse choice")
	}

	branch := original
	branch.Choices = cloneAccumulatorSlice(branch.Choices)
	branchChunk := storageTestChunk(ChatCompletionChunkChoiceDelta{})
	branchChunk.Choices[0].Logprobs.Content = []ChatCompletionTokenLogprob{{Token: "branch"}}
	if !branch.AddChunk(branchChunk) {
		t.Fatal("AddChunk rejected the branch logprob")
	}
	if length := original.logprobState.choices[0].content.length; length != 0 {
		t.Fatalf("branch changed original logprob length to %d", length)
	}

	originalChunk := storageTestChunk(ChatCompletionChunkChoiceDelta{})
	originalChunk.Choices[0].Logprobs.Content = []ChatCompletionTokenLogprob{{Token: "original"}}
	if !original.AddChunk(originalChunk) {
		t.Fatal("AddChunk rejected the original logprob")
	}
	if length := original.logprobState.choices[0].content.length; length != 1 {
		t.Fatalf("original logprob length = %d, want 1", length)
	}
}

func TestAccumulatorValueCopyToolActivationStateIsolated(t *testing.T) {
	initial := storageTestChunk(ChatCompletionChunkChoiceDelta{
		ToolCalls: []ChatCompletionChunkChoiceDeltaToolCall{{Index: 1}},
	})

	var original ChatCompletionAccumulator
	if !original.AddChunk(initial) {
		t.Fatal("AddChunk rejected the initial sparse tool call")
	}

	branch := original
	branch.Choices = cloneAccumulatorSlice(branch.Choices)
	branch.Choices[0].Message.ToolCalls = cloneAccumulatorSlice(branch.Choices[0].Message.ToolCalls)
	branchChunk := storageTestChunk(ChatCompletionChunkChoiceDelta{
		ToolCalls: []ChatCompletionChunkChoiceDeltaToolCall{{Index: 0}},
	})
	if !branch.AddChunk(branchChunk) {
		t.Fatal("AddChunk rejected the branch tool call")
	}

	choiceState := original.stringState.choices[0]
	if original.stringState.activeTools != 1 || len(choiceState.activeToolCalls) != 1 || choiceState.activeToolCalls[0] != 1 {
		t.Fatalf(
			"branch changed original tool state: active tools %d, indices %v",
			original.stringState.activeTools,
			choiceState.activeToolCalls,
		)
	}
}

func TestAccumulatorToolActivationStateGrowsAmortized(t *testing.T) {
	const toolCount = 128

	var acc ChatCompletionAccumulator
	var previousToolCalls **chatCompletionToolCallStringState
	var previousActiveToolCalls *int
	toolCallBackingChanges := 0
	activeToolCallBackingChanges := 0
	for i := range toolCount {
		toolCall := ChatCompletionChunkChoiceDeltaToolCall{Index: int64(i)}
		chunk := storageTestChunk(ChatCompletionChunkChoiceDelta{
			ToolCalls: []ChatCompletionChunkChoiceDeltaToolCall{toolCall, toolCall},
		})
		if !acc.AddChunk(chunk) {
			t.Fatalf("AddChunk rejected dense tool call %d", i)
		}

		choiceState := acc.stringState.choices[0]
		toolCalls := unsafe.SliceData(choiceState.toolCalls)
		if toolCalls != previousToolCalls {
			toolCallBackingChanges++
			previousToolCalls = toolCalls
		}
		activeToolCalls := unsafe.SliceData(choiceState.activeToolCalls)
		if activeToolCalls != previousActiveToolCalls {
			activeToolCallBackingChanges++
			previousActiveToolCalls = activeToolCalls
		}
	}

	if acc.stringState.activeTools != toolCount {
		t.Fatalf("active tools = %d, want %d", acc.stringState.activeTools, toolCount)
	}
	if toolCallBackingChanges > 16 || activeToolCallBackingChanges > 16 {
		t.Fatalf(
			"tool activation repeatedly reallocated state: tool calls %d, active indices %d backing changes",
			toolCallBackingChanges,
			activeToolCallBackingChanges,
		)
	}
}

func TestAccumulatorSteadyLogprobStateDoesNotAllocate(t *testing.T) {
	chunk := storageTestChunk(ChatCompletionChunkChoiceDelta{})
	chunk.Choices[0].Logprobs.Content = []ChatCompletionTokenLogprob{{Token: "token"}}
	completion := ChatCompletion{Choices: []ChatCompletionChoice{{
		Logprobs: ChatCompletionChoiceLogprobs{Content: []ChatCompletionTokenLogprob{{Token: "token"}}},
	}}}
	state := chatCompletionAccumulatorLogprobState{
		choices: []chatCompletionChoiceLogprobState{{
			content: chatCompletionLogprobSliceState{length: 1, capacity: 1},
		}},
	}

	allocations := testing.AllocsPerRun(100, func() {
		state.acceptChunk(&completion, &chunk)
	})
	if allocations != 0 {
		t.Fatalf("steady logprob state allocated %.0f times per chunk, want 0", allocations)
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
	wantWork := 10 + 2*len(chunk.ID) // six normal passes, two staged-slice passes, measurement, sparse commit, and established ID
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
	wantWork := normalPasses + 2*len(chunk.ID) + 5*replacementBytes
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
	wantWork := chunkCount*chatCompletionAccumulatorChoiceWork + len(chunk.ID) +
		(chunkCount-1)*(7+2*len(chunk.ID))
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

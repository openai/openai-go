package openai

import (
	"strings"
	"testing"
	"unsafe"
)

func TestAccumulatorDenseToolReconciliationIsDeltaLocal(t *testing.T) {
	const toolCount = 4_096
	chunk := storageTestChunk(ChatCompletionChunkChoiceDelta{
		ToolCalls: make([]ChatCompletionChunkChoiceDeltaToolCall, toolCount),
	})
	for i := range chunk.Choices[0].Delta.ToolCalls {
		chunk.Choices[0].Delta.ToolCalls[i].Index = int64(i)
	}

	var acc ChatCompletionAccumulator
	if !acc.AddChunk(chunk) {
		t.Fatal("AddChunk rejected dense tool-call state")
	}
	before := acc.reconciliationWork
	if !acc.AddChunk(ChatCompletionChunk{ID: chunk.ID}) {
		t.Fatal("AddChunk rejected the next empty chunk")
	}
	if work := acc.reconciliationWork - before; work >= toolCount {
		t.Fatalf("empty chunk scanned %d populated tool calls; want delta-local reconciliation", work)
	}
}

func TestAccumulatorRoundRobinReleasesUntouchedToolBacking(t *testing.T) {
	chunk := storageTestChunk(ChatCompletionChunkChoiceDelta{
		ToolCalls: make([]ChatCompletionChunkChoiceDeltaToolCall, 4),
	})
	for i := range chunk.Choices[0].Delta.ToolCalls {
		chunk.Choices[0].Delta.ToolCalls[i].Index = int64(i)
		chunk.Choices[0].Delta.ToolCalls[i].Function.Arguments = strings.Repeat("x", 1<<10)
	}

	var acc ChatCompletionAccumulator
	if !acc.AddChunk(chunk) {
		t.Fatal("AddChunk rejected the initial dense tool calls")
	}
	acc.Choices[0].Message.ToolCalls[3].Function.Arguments = ""
	for range len(chunk.Choices[0].Delta.ToolCalls) {
		if !acc.AddChunk(ChatCompletionChunk{ID: chunk.ID}) {
			t.Fatal("AddChunk rejected a bounded tool-reconciliation sweep")
		}
	}
	if state := &acc.stringState.choices[0].toolCallState(3).arguments; cap(state.buffer) != 0 || state.published != "" {
		t.Fatalf("untouched cleared tool backing remained: capacity %d, published length %d", cap(state.buffer), len(state.published))
	}
}

func TestAccumulatorReleasesClearedToolBackingAtIncomingGrowthRate(t *testing.T) {
	const toolCount = 16
	chunk := storageTestChunk(ChatCompletionChunkChoiceDelta{
		ToolCalls: make([]ChatCompletionChunkChoiceDeltaToolCall, toolCount),
	})
	for i := range chunk.Choices[0].Delta.ToolCalls {
		chunk.Choices[0].Delta.ToolCalls[i].Index = int64(i)
		chunk.Choices[0].Delta.ToolCalls[i].Function.Name = strings.Repeat("x", 1<<10)
		chunk.Choices[0].Delta.ToolCalls[i].Function.Arguments = strings.Repeat("x", 1<<10)
	}

	var acc ChatCompletionAccumulator
	if !acc.AddChunk(chunk) {
		t.Fatal("AddChunk rejected the initial dense tool batch")
	}
	for i := range acc.Choices[0].Message.ToolCalls {
		acc.Choices[0].Message.ToolCalls[i].Function.Name = ""
		acc.Choices[0].Message.ToolCalls[i].Function.Arguments = ""
	}
	for i := range chunk.Choices[0].Delta.ToolCalls {
		chunk.Choices[0].Delta.ToolCalls[i].Index = int64(toolCount + i)
	}
	if !acc.AddChunk(chunk) {
		t.Fatal("AddChunk rejected the next dense tool batch")
	}
	for i := range toolCount {
		state := acc.stringState.choices[0].toolCallState(i)
		if cap(state.name.buffer) != 0 || cap(state.arguments.buffer) != 0 {
			t.Fatalf("cleared tool %d retains name capacity %d and arguments capacity %d", i,
				cap(state.name.buffer), cap(state.arguments.buffer))
		}
	}
}

func TestAccumulatorToolProjectionUsesDistinctChoiceAndDenseIndices(t *testing.T) {
	var projections chatCompletionToolProjectionTable
	first := lookupChatCompletionToolMetadataProjection(&projections, 0, 1_024)
	first.fields = projectedToolID
	second := lookupChatCompletionToolMetadataProjection(&projections, 1, 0)
	second.fields = projectedToolType

	if first == second || first.fields != projectedToolID || second.fields != projectedToolType {
		t.Fatal("projection aliased dense tool indices belonging to different choices")
	}
}

func TestAccumulatorProjectsOnlyToolsWithMetadata(t *testing.T) {
	const toolCount = 2*chatCompletionAccumulatorInlineProjectionSlots + 1
	chunk := storageTestChunk(ChatCompletionChunkChoiceDelta{
		ToolCalls: make([]ChatCompletionChunkChoiceDeltaToolCall, toolCount),
	})
	for i := range chunk.Choices[0].Delta.ToolCalls {
		chunk.Choices[0].Delta.ToolCalls[i].Index = int64(i)
	}
	chunk.Choices[0].Delta.ToolCalls[toolCount-1].ID = "tool-id"
	var acc ChatCompletionAccumulator
	allocations := testing.AllocsPerRun(5, func() {
		work := 0
		if !acc.addChatCompletionToolMetadataWork(&work, &chunk) {
			panic("metadata projection rejected a valid dense tool collection")
		}
	})
	if allocations >= 10 {
		t.Fatalf("sparse metadata projection allocated %.0f times for one metadata-bearing tool", allocations)
	}
}

func TestAccumulatorLogprobAccountingRejectsIntegerOverflow(t *testing.T) {
	total := maxChatCompletionAccumulatorInt - 1
	if !addAccumulatorLogprobBytes(&total, 1) || total != maxChatCompletionAccumulatorInt {
		t.Fatal("logprob accounting rejected the largest representable total")
	}
	if addAccumulatorLogprobBytes(&total, 1) || total != maxChatCompletionAccumulatorInt {
		t.Fatal("logprob accounting accepted or mutated an overflowing total")
	}
	total = maxChatCompletionAccumulatorInt - 1
	if addAccumulatorLogprobStorage(&total, 2, 1) || total != maxChatCompletionAccumulatorInt-1 {
		t.Fatal("logprob storage accounting accepted or mutated an overflowing total")
	}
}

func TestAccumulatorRejectedLogprobCopyDoesNotDetachSharedBacking(t *testing.T) {
	initial := storageTestChunk(ChatCompletionChunkChoiceDelta{})
	initial.Choices[0].Logprobs.Content = []ChatCompletionTokenLogprob{
		{Token: "first"}, {Token: "second"}, {Token: "third"},
	}
	var original ChatCompletionAccumulator
	if !original.AddChunk(initial) {
		t.Fatal("AddChunk rejected the initial logprob slice")
	}

	branch := original
	branch.Choices = cloneAccumulatorSlice(branch.Choices)
	if !branch.AddChunk(ChatCompletionChunk{ID: initial.ID}) {
		t.Fatal("AddChunk rejected the copied accumulator's empty chunk")
	}
	before := unsafe.SliceData(branch.Choices[0].Logprobs.Content)
	branch.logprobBytes = maxChatCompletionAccumulatorInt
	next := storageTestChunk(ChatCompletionChunkChoiceDelta{})
	next.Choices[0].Logprobs.Content = []ChatCompletionTokenLogprob{{Token: "branch"}}
	if branch.AddChunk(next) {
		t.Fatal("AddChunk accepted overflowing copied-logprob accounting")
	}
	if got := unsafe.SliceData(branch.Choices[0].Logprobs.Content); got != before {
		t.Fatal("rejected copied-logprob append detached public backing")
	}
	if state := branch.logprobState.choices[0].content; !state.shared || state.data.Value() != before {
		t.Fatal("rejected copied-logprob append changed private ownership state")
	}
}

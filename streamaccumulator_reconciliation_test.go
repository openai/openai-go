package openai

import (
	"strings"
	"testing"
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
	if state := &acc.stringState.choices[0].toolCalls[3].arguments; cap(state.buffer) != 0 || state.published != "" {
		t.Fatalf("untouched cleared tool backing remained: capacity %d, published length %d", cap(state.buffer), len(state.published))
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

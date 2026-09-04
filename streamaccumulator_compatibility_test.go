package openai_test

import (
	"strings"
	"testing"

	openai "github.com/openai/openai-go/v3"
)

func TestAccumulatorAcceptsStreamsBeyondFormerChunkBudget(t *testing.T) {
	var acc openai.ChatCompletionAccumulator
	chunk := accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{Content: "x"})

	for i := range testAccumulatorMaxChunks + 1 {
		if !acc.AddChunk(chunk) {
			t.Fatalf("AddChunk rejected historically supported chunk %d", i+1)
		}
	}
	if got := len(acc.Choices[0].Message.Content); got != testAccumulatorMaxChunks+1 {
		t.Fatalf("accumulated content length = %d, want %d", got, testAccumulatorMaxChunks+1)
	}
}

func TestAccumulatorAcceptsDenseToolCallsBeyondFormerStructuralBudget(t *testing.T) {
	const toolCount = testAccumulatorMaxStructuralSlots + 1
	var acc openai.ChatCompletionAccumulator
	chunk := accumulatorDenseToolStringChunk(toolCount)

	if !acc.AddChunk(chunk) {
		t.Fatalf("AddChunk rejected %d historically supported dense tool calls", toolCount)
	}
	if got := len(acc.Choices[0].Message.ToolCalls); got != toolCount {
		t.Fatalf("accumulated tool calls = %d, want %d", got, toolCount)
	}
}

func TestAccumulatorAcceptsDenseToolMetadataBeyondFormerProjectionTable(t *testing.T) {
	const toolCount = 2*testAccumulatorMaxStructuralSlots + 1
	chunk := accumulatorDenseToolStringChunk(toolCount)
	for i := range chunk.Choices[0].Delta.ToolCalls {
		chunk.Choices[0].Delta.ToolCalls[i].ID = "tool-id"
	}

	var acc openai.ChatCompletionAccumulator
	if !acc.AddChunk(chunk) {
		t.Fatalf("AddChunk rejected %d historically supported dense tool metadata entries", toolCount)
	}
	if got := len(acc.Choices[0].Message.ToolCalls); got != toolCount {
		t.Fatalf("accumulated tool calls = %d, want %d", got, toolCount)
	}
}

func TestAccumulatorPreservesLargeDenseToolCompletionIndex(t *testing.T) {
	const toolIndex = 1 << 13
	var acc openai.ChatCompletionAccumulator
	initial := accumulatorDenseToolStringChunk(toolIndex + 1)
	if !acc.AddChunk(initial) {
		t.Fatal("AddChunk rejected the large dense tool-call collection")
	}

	active := accumulatorToolStringChunk("last-tool", "")
	active.Choices[0].Delta.ToolCalls[0].Index = toolIndex
	if !acc.AddChunk(active) {
		t.Fatal("AddChunk rejected the large dense tool-call index")
	}
	if !acc.AddChunk(accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{})) {
		t.Fatal("AddChunk rejected the tool completion chunk")
	}
	finished, ok := acc.JustFinishedToolCall()
	if !ok || finished.Index != toolIndex {
		t.Fatalf("finished tool index = %d, ok = %t; want %d, true", finished.Index, ok, toolIndex)
	}
}

func TestAccumulatorAcceptsLogprobsBeyondFormerBudget(t *testing.T) {
	var acc openai.ChatCompletionAccumulator
	chunk := accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{})
	chunk.Choices[0].Logprobs.Content = []openai.ChatCompletionTokenLogprob{{
		Token: strings.Repeat("x", testAccumulatorMaxLogprobBytes/2+1),
	}}

	for i := range 2 {
		if !acc.AddChunk(chunk) {
			t.Fatalf("AddChunk rejected historically supported logprob chunk %d", i+1)
		}
	}
	if got := len(acc.Choices[0].Logprobs.Content); got != 2 {
		t.Fatalf("accumulated logprobs = %d, want 2", got)
	}
}

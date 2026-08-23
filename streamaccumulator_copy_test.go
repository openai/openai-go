package openai_test

import (
	"slices"
	"strings"
	"testing"

	openai "github.com/openai/openai-go/v3"
)

func TestAccumulatorValueCopyTruncationPreservesOriginalState(t *testing.T) {
	t.Run("choices", func(t *testing.T) {
		var original openai.ChatCompletionAccumulator
		initial := accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{Content: "x"})
		if !original.AddChunk(initial) {
			t.Fatal("AddChunk rejected the initial chunk")
		}

		truncated := original
		truncated.Choices = truncated.Choices[:0]
		if !truncated.AddChunk(openai.ChatCompletionChunk{ID: initial.ID}) {
			t.Fatal("AddChunk rejected the chunk after the copied accumulator was truncated")
		}
		if !original.AddChunk(openai.ChatCompletionChunk{ID: initial.ID}) {
			t.Fatal("AddChunk rejected the original accumulator after its copy was truncated")
		}
		if content := original.Choices[0].Message.Content; content != "x" {
			t.Fatalf("original content = %q, want %q", content, "x")
		}
	})

	t.Run("tool_calls", func(t *testing.T) {
		var original openai.ChatCompletionAccumulator
		initial := accumulatorToolStringChunk("tool", "arguments")
		if !original.AddChunk(initial) {
			t.Fatal("AddChunk rejected the initial chunk")
		}

		truncated := original
		truncated.Choices = slices.Clone(truncated.Choices)
		truncated.Choices[0].Message.ToolCalls = truncated.Choices[0].Message.ToolCalls[:0]
		if !truncated.AddChunk(openai.ChatCompletionChunk{ID: initial.ID}) {
			t.Fatal("AddChunk rejected the chunk after the copied tool calls were truncated")
		}
		if !original.AddChunk(openai.ChatCompletionChunk{ID: initial.ID}) {
			t.Fatal("AddChunk rejected the original accumulator after its copy's tool calls were truncated")
		}
		if name := original.Choices[0].Message.ToolCalls[0].Function.Name; name != "tool" {
			t.Fatalf("original tool name = %q, want %q", name, "tool")
		}
	})
}

func TestAccumulatorValueCopyChoiceActivationPreservesAccounting(t *testing.T) {
	var original openai.ChatCompletionAccumulator
	initial := accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{Content: "existing"})
	initial.Choices[0].Index = 1
	if !original.AddChunk(initial) {
		t.Fatal("AddChunk rejected the initial sparse choice")
	}

	branch := original
	branch.Choices = slices.Clone(branch.Choices)
	branchChunk := accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{Content: "branch"})
	if !branch.AddChunk(branchChunk) {
		t.Fatal("AddChunk rejected the branch choice")
	}
	originalChunk := accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{Content: "original"})
	if !original.AddChunk(originalChunk) {
		t.Fatal("AddChunk rejected the original choice")
	}

	original.Choices[0].Message.Content = strings.Repeat("x", testAccumulatorMaxTextBytes)
	if original.AddChunk(accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{Content: "x"})) {
		t.Fatal("AddChunk omitted a choice activated after a value copy from the live-text budget")
	}
}

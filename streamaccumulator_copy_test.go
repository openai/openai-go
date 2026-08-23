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

func TestAccumulatorValueCopyStringBuffersIsolated(t *testing.T) {
	tests := []struct {
		name     string
		initial  openai.ChatCompletionChunk
		branch   openai.ChatCompletionChunk
		original openai.ChatCompletionChunk
		value    func(*openai.ChatCompletionAccumulator) string
	}{
		{
			name:     "content",
			initial:  accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{Content: "initial"}),
			branch:   accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{Content: "-branch"}),
			original: accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{Content: "-original"}),
			value: func(acc *openai.ChatCompletionAccumulator) string {
				return acc.Choices[0].Message.Content
			},
		},
		{
			name:     "refusal",
			initial:  accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{Refusal: "initial"}),
			branch:   accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{Refusal: "-branch"}),
			original: accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{Refusal: "-original"}),
			value: func(acc *openai.ChatCompletionAccumulator) string {
				return acc.Choices[0].Message.Refusal
			},
		},
		{
			name:     "tool_name",
			initial:  accumulatorToolStringChunk("initial", ""),
			branch:   accumulatorToolStringChunk("-branch", ""),
			original: accumulatorToolStringChunk("-original", ""),
			value: func(acc *openai.ChatCompletionAccumulator) string {
				return acc.Choices[0].Message.ToolCalls[0].Function.Name
			},
		},
		{
			name:     "tool_arguments",
			initial:  accumulatorToolStringChunk("", "initial"),
			branch:   accumulatorToolStringChunk("", "-branch"),
			original: accumulatorToolStringChunk("", "-original"),
			value: func(acc *openai.ChatCompletionAccumulator) string {
				return acc.Choices[0].Message.ToolCalls[0].Function.Arguments
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var original openai.ChatCompletionAccumulator
			if !original.AddChunk(test.initial) {
				t.Fatal("AddChunk rejected the initial chunk")
			}

			branch := original
			branch.Choices = slices.Clone(branch.Choices)
			branch.Choices[0].Message.ToolCalls = slices.Clone(branch.Choices[0].Message.ToolCalls)
			if !branch.AddChunk(test.branch) {
				t.Fatal("AddChunk rejected the branch suffix")
			}
			if !original.AddChunk(test.original) {
				t.Fatal("AddChunk rejected the original suffix")
			}

			if got := test.value(&branch); got != "initial-branch" {
				t.Fatalf("branch value = %q, want %q", got, "initial-branch")
			}
			if got := test.value(&original); got != "initial-original" {
				t.Fatalf("original value = %q, want %q", got, "initial-original")
			}
		})
	}
}

func TestAccumulatorValueCopyResponseStateIsolated(t *testing.T) {
	var original openai.ChatCompletionAccumulator
	initial := accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{Content: "existing"})
	initial.Choices[0].Index = 1
	if !original.AddChunk(initial) {
		t.Fatal("AddChunk rejected the initial sparse choice")
	}

	branch := original
	branch.Choices = slices.Clone(branch.Choices)
	if !branch.AddChunk(accumulatorToolStringChunk("tool", "arguments")) {
		t.Fatal("AddChunk rejected the branch tool call")
	}
	if !original.AddChunk(accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{})) {
		t.Fatal("AddChunk rejected the original empty delta")
	}

	if _, ok := original.JustFinishedToolCall(); ok {
		t.Fatal("original reported a tool call activated only in its copied branch")
	}
	if _, ok := original.JustFinishedToolCallForChoice(0); ok {
		t.Fatal("original reported a per-choice tool call activated only in its copied branch")
	}
}

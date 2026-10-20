package openai_test

import (
	"runtime"
	"slices"
	"strings"
	"testing"
	"time"

	openai "github.com/openai/openai-go/v3"
)

func TestAccumulatorValueCopyDoesNotRetainSource(t *testing.T) {
	finalized := make(chan struct{}, 1)
	var dormant openai.ChatCompletionAccumulator
	func() {
		source := &openai.ChatCompletionAccumulator{}
		initial := accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{Content: "existing"})
		initial.Choices[0].Index = 1
		if !source.AddChunk(initial) {
			t.Fatal("AddChunk rejected the initial sparse choice")
		}

		dormant = *source
		source.Choices = slices.Clone(source.Choices)
		later := accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{
			Content: strings.Repeat("x", 1<<20),
		})
		if !source.AddChunk(later) {
			t.Fatal("AddChunk rejected storage accumulated after the value copy")
		}
		runtime.SetFinalizer(source, func(*openai.ChatCompletionAccumulator) {
			finalized <- struct{}{}
		})
	}()

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		runtime.GC()
		select {
		case <-finalized:
			runtime.KeepAlive(&dormant)
			return
		default:
			runtime.Gosched()
		}
	}
	runtime.KeepAlive(&dormant)
	t.Fatal("a dormant accumulator copy retained its later-mutated source")
}

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

	t.Run("original_choices", func(t *testing.T) {
		var original openai.ChatCompletionAccumulator
		initial := accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{Content: "x"})
		if !original.AddChunk(initial) {
			t.Fatal("AddChunk rejected the initial chunk")
		}

		copy := original
		copy.Choices = slices.Clone(copy.Choices)
		original.Choices = original.Choices[:0]
		if !original.AddChunk(openai.ChatCompletionChunk{ID: initial.ID}) {
			t.Fatal("AddChunk rejected the chunk after the original accumulator was truncated")
		}
		if !copy.AddChunk(openai.ChatCompletionChunk{ID: initial.ID}) {
			t.Fatal("AddChunk rejected the copy after the original accumulator was truncated")
		}
		if content := copy.Choices[0].Message.Content; content != "x" {
			t.Fatalf("copied content = %q, want %q", content, "x")
		}
	})

	t.Run("original_tool_calls", func(t *testing.T) {
		var original openai.ChatCompletionAccumulator
		initial := accumulatorToolStringChunk("tool", "arguments")
		if !original.AddChunk(initial) {
			t.Fatal("AddChunk rejected the initial chunk")
		}

		copy := original
		copy.Choices = slices.Clone(copy.Choices)
		copy.Choices[0].Message.ToolCalls = slices.Clone(copy.Choices[0].Message.ToolCalls)
		original.Choices[0].Message.ToolCalls = original.Choices[0].Message.ToolCalls[:0]
		if !original.AddChunk(openai.ChatCompletionChunk{ID: initial.ID}) {
			t.Fatal("AddChunk rejected the chunk after the original tool calls were truncated")
		}
		if !copy.AddChunk(openai.ChatCompletionChunk{ID: initial.ID}) {
			t.Fatal("AddChunk rejected the copy after the original tool calls were truncated")
		}
		if name := copy.Choices[0].Message.ToolCalls[0].Function.Name; name != "tool" {
			t.Fatalf("copied tool name = %q, want %q", name, "tool")
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

	replacement := strings.Repeat("x", 2*testAccumulatorLargeTextBytes+1)
	original.Choices[0].Message.Content = replacement
	if !original.AddChunk(accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{Content: "x"})) {
		t.Fatal("AddChunk rejected large public text in a choice activated after a value copy")
	}
	if got := original.Choices[0].Message.Content; got != replacement+"x" {
		t.Fatal("AddChunk did not preserve the activated choice's large replacement")
	}
	if got := branch.Choices[0].Message.Content; got != "branch" {
		t.Fatalf("AddChunk changed the copied accumulator's content to %q", got)
	}
}

func TestAccumulatorValueCopyAllowsLargeTextPrefixes(t *testing.T) {
	prefix := strings.Repeat("x", testAccumulatorMaxReconcileWork+1)
	initial := accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{Content: prefix})
	var original openai.ChatCompletionAccumulator
	if !original.AddChunk(initial) {
		t.Fatal("AddChunk rejected the original large text prefix")
	}

	copy := original
	copy.Choices = slices.Clone(copy.Choices)
	if !copy.AddChunk(accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{Content: "y"})) {
		t.Fatal("AddChunk rejected a copied text prefix beyond the structural-work budget")
	}
	if got := copy.Choices[0].Message.Content; got != prefix+"y" {
		t.Fatal("AddChunk did not preserve the copied large text prefix")
	}
	if got := original.Choices[0].Message.Content; got != prefix {
		t.Fatal("AddChunk changed the original accumulator's large text prefix")
	}
}

func TestAccumulatorValueCopyLogprobsDetachSharedSpareCapacity(t *testing.T) {
	tests := []struct {
		name string
		set  func(*openai.ChatCompletionChunkChoice, []openai.ChatCompletionTokenLogprob)
		get  func(*openai.ChatCompletionAccumulator) []openai.ChatCompletionTokenLogprob
	}{
		{
			name: "content",
			set: func(choice *openai.ChatCompletionChunkChoice, values []openai.ChatCompletionTokenLogprob) {
				choice.Logprobs.Content = values
			},
			get: func(acc *openai.ChatCompletionAccumulator) []openai.ChatCompletionTokenLogprob {
				return acc.Choices[0].Logprobs.Content
			},
		},
		{
			name: "refusal",
			set: func(choice *openai.ChatCompletionChunkChoice, values []openai.ChatCompletionTokenLogprob) {
				choice.Logprobs.Refusal = values
			},
			get: func(acc *openai.ChatCompletionAccumulator) []openai.ChatCompletionTokenLogprob {
				return acc.Choices[0].Logprobs.Refusal
			},
		},
	}

	for _, test := range tests {
		for _, delayed := range []bool{false, true} {
			name := "immediate"
			if delayed {
				name = "delayed"
			}
			t.Run(test.name+"/"+name, func(t *testing.T) {
				initial := accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{})
				test.set(&initial.Choices[0], []openai.ChatCompletionTokenLogprob{
					{Token: "first"}, {Token: "second"}, {Token: "third"},
				})
				var original openai.ChatCompletionAccumulator
				if !original.AddChunk(initial) {
					t.Fatal("AddChunk rejected the initial logprobs")
				}
				if logprobs := test.get(&original); len(logprobs) != 3 || cap(logprobs) != 4 {
					t.Fatalf("initial logprob backing: len %d cap %d, want 3 and 4", len(logprobs), cap(logprobs))
				}

				branch := original
				branch.Choices = slices.Clone(branch.Choices)
				if delayed && !branch.AddChunk(openai.ChatCompletionChunk{ID: initial.ID}) {
					t.Fatal("AddChunk rejected the copied accumulator's initial empty chunk")
				}
				branchChunk := accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{})
				test.set(&branchChunk.Choices[0], []openai.ChatCompletionTokenLogprob{{Token: "branch"}})
				if !branch.AddChunk(branchChunk) {
					t.Fatal("AddChunk rejected the copied accumulator's in-capacity append")
				}

				originalChunk := accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{})
				test.set(&originalChunk.Choices[0], []openai.ChatCompletionTokenLogprob{{Token: "original"}})
				if !original.AddChunk(originalChunk) {
					t.Fatal("AddChunk rejected the original accumulator's in-capacity append")
				}
				if got := test.get(&branch)[3].Token; got != "branch" {
					t.Fatalf("original append changed the copied logprob to %q, want branch", got)
				}
				if got := test.get(&original)[3].Token; got != "original" {
					t.Fatalf("original appended logprob = %q, want original", got)
				}
			})
		}
	}
}

func TestAccumulatorValueCopyPreservesUntouchedLogprobOwnership(t *testing.T) {
	for _, firstContent := range []bool{true, false} {
		name := "content_then_refusal"
		if !firstContent {
			name = "refusal_then_content"
		}
		t.Run(name, func(t *testing.T) {
			initial := accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{})
			initial.Choices[0].Logprobs.Content = []openai.ChatCompletionTokenLogprob{
				{Token: "first"}, {Token: "second"}, {Token: "third"},
			}
			initial.Choices[0].Logprobs.Refusal = []openai.ChatCompletionTokenLogprob{
				{Token: "first"}, {Token: "second"}, {Token: "third"},
			}
			var original openai.ChatCompletionAccumulator
			if !original.AddChunk(initial) {
				t.Fatal("AddChunk rejected the initial Content and Refusal logprobs")
			}
			branch := original
			branch.Choices = slices.Clone(branch.Choices)

			first := accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{})
			second := accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{})
			other := accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{})
			var branchValues func() []openai.ChatCompletionTokenLogprob
			if firstContent {
				first.Choices[0].Logprobs.Content = []openai.ChatCompletionTokenLogprob{{Token: "first-branch"}}
				second.Choices[0].Logprobs.Refusal = []openai.ChatCompletionTokenLogprob{{Token: "branch"}}
				other.Choices[0].Logprobs.Refusal = []openai.ChatCompletionTokenLogprob{{Token: "original"}}
				branchValues = func() []openai.ChatCompletionTokenLogprob { return branch.Choices[0].Logprobs.Refusal }
			} else {
				first.Choices[0].Logprobs.Refusal = []openai.ChatCompletionTokenLogprob{{Token: "first-branch"}}
				second.Choices[0].Logprobs.Content = []openai.ChatCompletionTokenLogprob{{Token: "branch"}}
				other.Choices[0].Logprobs.Content = []openai.ChatCompletionTokenLogprob{{Token: "original"}}
				branchValues = func() []openai.ChatCompletionTokenLogprob { return branch.Choices[0].Logprobs.Content }
			}
			if !branch.AddChunk(first) || !branch.AddChunk(second) {
				t.Fatal("AddChunk rejected copied logprobs in different fields")
			}
			if !original.AddChunk(other) {
				t.Fatal("AddChunk rejected the original accumulator's untouched-field append")
			}
			if got := branchValues()[3].Token; got != "branch" {
				t.Fatalf("original append changed the copied sibling logprob to %q, want branch", got)
			}
		})
	}

	t.Run("different_choices", func(t *testing.T) {
		initial := accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{})
		initial.Choices[0].Logprobs.Content = []openai.ChatCompletionTokenLogprob{
			{Token: "first"}, {Token: "second"}, {Token: "third"},
		}
		secondChoice := initial.Choices[0]
		secondChoice.Index = 1
		secondChoice.Logprobs.Content = slices.Clone(initial.Choices[0].Logprobs.Content)
		initial.Choices = append(initial.Choices, secondChoice)
		var original openai.ChatCompletionAccumulator
		if !original.AddChunk(initial) {
			t.Fatal("AddChunk rejected initial logprobs for both choices")
		}
		branch := original
		branch.Choices = slices.Clone(branch.Choices)
		first := accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{})
		first.Choices[0].Logprobs.Content = []openai.ChatCompletionTokenLogprob{{Token: "first-branch"}}
		if !branch.AddChunk(first) {
			t.Fatal("AddChunk rejected the copied accumulator's first choice")
		}
		next := accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{})
		next.Choices[0].Index = 1
		next.Choices[0].Logprobs.Content = []openai.ChatCompletionTokenLogprob{{Token: "branch"}}
		if !branch.AddChunk(next) {
			t.Fatal("AddChunk rejected the copied accumulator's second choice")
		}
		next.Choices[0].Logprobs.Content[0].Token = "original"
		if !original.AddChunk(next) {
			t.Fatal("AddChunk rejected the original accumulator's second choice")
		}
		if got := branch.Choices[1].Logprobs.Content[3].Token; got != "branch" {
			t.Fatalf("original append changed copied second-choice logprob to %q, want branch", got)
		}
	})
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
			for _, branchFirst := range []bool{true, false} {
				name := "original_first"
				if branchFirst {
					name = "branch_first"
				}
				t.Run(name, func(t *testing.T) {
					var original openai.ChatCompletionAccumulator
					if !original.AddChunk(test.initial) {
						t.Fatal("AddChunk rejected the initial chunk")
					}

					branch := original
					branch.Choices = slices.Clone(branch.Choices)
					branch.Choices[0].Message.ToolCalls = slices.Clone(branch.Choices[0].Message.ToolCalls)
					if branchFirst {
						if !branch.AddChunk(test.branch) {
							t.Fatal("AddChunk rejected the branch suffix")
						}
						if !original.AddChunk(test.original) {
							t.Fatal("AddChunk rejected the original suffix")
						}
					} else {
						if !original.AddChunk(test.original) {
							t.Fatal("AddChunk rejected the original suffix")
						}
						if !branch.AddChunk(test.branch) {
							t.Fatal("AddChunk rejected the branch suffix")
						}
					}

					if got := test.value(&branch); got != "initial-branch" {
						t.Fatalf("branch value = %q, want %q", got, "initial-branch")
					}
					if got := test.value(&original); got != "initial-original" {
						t.Fatalf("original value = %q, want %q", got, "initial-original")
					}
				})
			}
		})
	}
}

func TestAccumulatorValueCopyResponseStateIsolated(t *testing.T) {
	for _, branchFirst := range []bool{true, false} {
		name := "original_first"
		if branchFirst {
			name = "branch_first"
		}
		t.Run(name, func(t *testing.T) {
			var original openai.ChatCompletionAccumulator
			initial := accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{Content: "existing"})
			initial.Choices[0].Index = 1
			if !original.AddChunk(initial) {
				t.Fatal("AddChunk rejected the initial sparse choice")
			}

			branch := original
			branch.Choices = slices.Clone(branch.Choices)
			if branchFirst {
				if !branch.AddChunk(accumulatorToolStringChunk("tool", "arguments")) {
					t.Fatal("AddChunk rejected the branch tool call")
				}
				if !original.AddChunk(accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{})) {
					t.Fatal("AddChunk rejected the original empty delta")
				}
			} else {
				if !original.AddChunk(accumulatorToolStringChunk("tool", "arguments")) {
					t.Fatal("AddChunk rejected the original tool call")
				}
				if !branch.AddChunk(accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{})) {
					t.Fatal("AddChunk rejected the branch empty delta")
				}
			}

			unmodified := &original
			if !branchFirst {
				unmodified = &branch
			}
			if _, ok := unmodified.JustFinishedToolCall(); ok {
				t.Fatal("unmodified accumulator reported a tool call activated only in its copy")
			}
			if _, ok := unmodified.JustFinishedToolCallForChoice(0); ok {
				t.Fatal("unmodified accumulator reported a per-choice tool call activated only in its copy")
			}
		})
	}
}

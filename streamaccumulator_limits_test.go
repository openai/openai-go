package openai_test

import (
	"fmt"
	"strings"
	"testing"
	"unsafe"

	openai "github.com/openai/openai-go/v3"
)

const (
	testAccumulatorMaxChunks          = 100_000
	testAccumulatorMaxStructuralSlots = 1_024
	testAccumulatorMaxTextBytes       = 16 << 20
	testAccumulatorMaxLogprobBytes    = 16 << 20
)

func TestAccumulatorStringGrowthIsAmortized(t *testing.T) {
	const fragmentCount = 2_048
	chunk := openai.ChatCompletionChunk{
		ID: "chatcmpl-linear-growth",
		Choices: []openai.ChatCompletionChunkChoice{{
			Delta: openai.ChatCompletionChunkChoiceDelta{
				Content: "c",
				Refusal: "r",
				ToolCalls: []openai.ChatCompletionChunkChoiceDeltaToolCall{{
					Index: 0,
					Function: openai.ChatCompletionChunkChoiceDeltaToolCallFunction{
						Name:      "n",
						Arguments: "a",
					},
				}},
			},
		}},
	}

	allocs := testing.AllocsPerRun(3, func() {
		var acc openai.ChatCompletionAccumulator
		for range fragmentCount {
			if !acc.AddChunk(chunk) {
				panic("AddChunk rejected a chunk within the documented budgets")
			}
		}

		message := acc.Choices[0].Message
		if len(message.Content) != fragmentCount || len(message.Refusal) != fragmentCount ||
			len(message.ToolCalls[0].Function.Name) != fragmentCount ||
			len(message.ToolCalls[0].Function.Arguments) != fragmentCount {
			panic("AddChunk did not publish the complete accumulated strings")
		}
	})

	if allocs >= fragmentCount {
		t.Fatalf("AddChunk allocated %.0f times for %d fragments; want amortized growth", allocs, fragmentCount)
	}
}

func TestAccumulatorPublishesStringsAfterEveryChunk(t *testing.T) {
	var acc openai.ChatCompletionAccumulator
	chunk := openai.ChatCompletionChunk{
		ID: "chatcmpl-current-strings",
		Choices: []openai.ChatCompletionChunkChoice{{
			Delta: openai.ChatCompletionChunkChoiceDelta{
				Content: "content-1",
				Refusal: "refusal-1",
				ToolCalls: []openai.ChatCompletionChunkChoiceDeltaToolCall{{
					Index: 0,
					Function: openai.ChatCompletionChunkChoiceDeltaToolCallFunction{
						Name:      "name-1",
						Arguments: "arguments-1",
					},
				}},
			},
		}},
	}

	if !acc.AddChunk(chunk) {
		t.Fatal("AddChunk rejected the first chunk")
	}
	assertAccumulatorStrings(t, &acc, "content-1", "refusal-1", "name-1", "arguments-1")

	message := &acc.Choices[0].Message
	message.Content = "replacement-content"
	message.Refusal = "replacement-refusal"
	message.ToolCalls[0].Function.Name = "replacement-name"
	message.ToolCalls[0].Function.Arguments = "replacement-arguments"
	chunk.Choices[0].Delta.Content = "-2"
	chunk.Choices[0].Delta.Refusal = "-2"
	chunk.Choices[0].Delta.ToolCalls[0].Function.Name = "-2"
	chunk.Choices[0].Delta.ToolCalls[0].Function.Arguments = "-2"

	if !acc.AddChunk(chunk) {
		t.Fatal("AddChunk rejected the second chunk")
	}
	assertAccumulatorStrings(
		t,
		&acc,
		"replacement-content-2",
		"replacement-refusal-2",
		"replacement-name-2",
		"replacement-arguments-2",
	)
}

func TestAccumulatorBudgetsPublicStringReplacements(t *testing.T) {
	atLimit := strings.Repeat("x", testAccumulatorMaxTextBytes)
	tests := []struct {
		name  string
		chunk func(string) openai.ChatCompletionChunk
		value func(*openai.ChatCompletionAccumulator) *string
	}{
		{
			name: "content",
			chunk: func(text string) openai.ChatCompletionChunk {
				return accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{Content: text})
			},
			value: func(acc *openai.ChatCompletionAccumulator) *string {
				return &acc.Choices[0].Message.Content
			},
		},
		{
			name: "refusal",
			chunk: func(text string) openai.ChatCompletionChunk {
				return accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{Refusal: text})
			},
			value: func(acc *openai.ChatCompletionAccumulator) *string {
				return &acc.Choices[0].Message.Refusal
			},
		},
		{
			name: "tool_name",
			chunk: func(text string) openai.ChatCompletionChunk {
				return accumulatorToolStringChunk(text, "")
			},
			value: func(acc *openai.ChatCompletionAccumulator) *string {
				return &acc.Choices[0].Message.ToolCalls[0].Function.Name
			},
		},
		{
			name: "tool_arguments",
			chunk: func(text string) openai.ChatCompletionChunk {
				return accumulatorToolStringChunk("", text)
			},
			value: func(acc *openai.ChatCompletionAccumulator) *string {
				return &acc.Choices[0].Message.ToolCalls[0].Function.Arguments
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name+"/replacement_is_rejected", func(t *testing.T) {
			var acc openai.ChatCompletionAccumulator
			chunk := test.chunk("x")
			chunk.Model = "accepted-model"
			if !acc.AddChunk(chunk) {
				t.Fatal("AddChunk rejected the initial chunk")
			}

			*test.value(&acc) = atLimit
			chunk = test.chunk("x")
			chunk.Model = "rejected-model"
			if acc.AddChunk(chunk) {
				t.Fatal("AddChunk accepted text beyond the live public-string budget")
			}
			if acc.Model != "accepted-model" || *test.value(&acc) != atLimit {
				t.Fatal("AddChunk mutated the accumulator after rejecting the chunk")
			}
		})

		t.Run(test.name+"/clearing_recovers_budget", func(t *testing.T) {
			var acc openai.ChatCompletionAccumulator
			chunk := test.chunk(atLimit)
			chunk.Model = "initial-model"
			if !acc.AddChunk(chunk) {
				t.Fatal("AddChunk rejected text at the documented aggregate budget")
			}

			*test.value(&acc) = ""
			chunk = test.chunk("x")
			chunk.Model = "recovered-model"
			if !acc.AddChunk(chunk) {
				t.Fatal("AddChunk did not recover budget after the public string was cleared")
			}
			if acc.Model != "recovered-model" || *test.value(&acc) != "x" {
				t.Fatal("AddChunk did not accumulate the accepted chunk after budget recovery")
			}
		})
	}
}

func TestAccumulatorRejectsTextBeyondBudgetWithoutMutation(t *testing.T) {
	var acc openai.ChatCompletionAccumulator
	chunk := openai.ChatCompletionChunk{
		ID:    "chatcmpl-oversized",
		Model: "accepted-model",
		Choices: []openai.ChatCompletionChunkChoice{{
			Delta: openai.ChatCompletionChunkChoiceDelta{
				Content: strings.Repeat("x", testAccumulatorMaxTextBytes),
			},
		}},
	}

	if !acc.AddChunk(chunk) {
		t.Fatal("AddChunk rejected text at the documented aggregate budget")
	}
	chunk.Model = "rejected-model"
	chunk.Choices[0].Delta.Content = "x"
	if acc.AddChunk(chunk) {
		t.Fatal("AddChunk accepted text beyond the documented aggregate budget")
	}
	if acc.Model != "accepted-model" || len(acc.Choices[0].Message.Content) != testAccumulatorMaxTextBytes {
		t.Fatalf("AddChunk mutated the accumulator after rejecting the chunk: model %q, content length %d", acc.Model, len(acc.Choices[0].Message.Content))
	}
}

func TestAccumulatorRejectsLogprobsBeyondBudgetWithoutMutation(t *testing.T) {
	var acc openai.ChatCompletionAccumulator
	toolChunk := accumulatorToolStringChunk("tool", "arguments")
	if !acc.AddChunk(toolChunk) {
		t.Fatal("AddChunk rejected the initial tool-call chunk")
	}

	logprobOverhead := 2 * int(unsafe.Sizeof(openai.ChatCompletionTokenLogprob{}))
	byteCount := (testAccumulatorMaxLogprobBytes - logprobOverhead) / int(unsafe.Sizeof(int64(0)))
	remaining := testAccumulatorMaxLogprobBytes - logprobOverhead - byteCount*int(unsafe.Sizeof(int64(0)))
	atLimit := accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{})
	atLimit.Model = "accepted-model"
	atLimit.Choices[0].Logprobs.Content = []openai.ChatCompletionTokenLogprob{{
		Token: strings.Repeat("x", remaining),
		Bytes: make([]int64, byteCount),
	}}
	if !acc.AddChunk(atLimit) {
		t.Fatal("AddChunk rejected logprobs at the documented aggregate budget")
	}
	if toolCall, ok := acc.JustFinishedToolCall(); !ok || toolCall.Name != "tool" {
		t.Fatal("AddChunk did not publish the expected pre-rejection tool-call event")
	}

	beyondLimit := accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{})
	beyondLimit.Model = "rejected-model"
	beyondLimit.Choices[0].Logprobs.Content = []openai.ChatCompletionTokenLogprob{{Token: "x"}}
	if acc.AddChunk(beyondLimit) {
		t.Fatal("AddChunk accepted logprobs beyond the documented aggregate budget")
	}
	if acc.Model != "accepted-model" || len(acc.Choices[0].Logprobs.Content) != 1 {
		t.Fatal("AddChunk mutated the accumulator after rejecting excessive logprobs")
	}
	if toolCall, ok := acc.JustFinishedToolCall(); !ok || toolCall.Name != "tool" {
		t.Fatal("AddChunk changed the current completion event after rejecting excessive logprobs")
	}

	afterRejection := accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{Content: "-after"})
	if !acc.AddChunk(afterRejection) {
		t.Fatal("AddChunk consumed unrelated text or chunk budget after rejecting excessive logprobs")
	}
	if got := acc.Choices[0].Message.Content; got != "-after" {
		t.Fatalf("content after rejected logprobs = %q, want %q", got, "-after")
	}
}

func TestAccumulatorRejectsChunksBeyondBudgetWithoutMutation(t *testing.T) {
	var acc openai.ChatCompletionAccumulator
	chunk := openai.ChatCompletionChunk{ID: "chatcmpl-too-many-chunks", Model: "accepted-model"}

	for i := range testAccumulatorMaxChunks {
		if !acc.AddChunk(chunk) {
			t.Fatalf("AddChunk rejected chunk %d within the documented budget", i+1)
		}
	}

	chunk.Model = "rejected-model"
	if acc.AddChunk(chunk) {
		t.Fatal("AddChunk accepted a chunk beyond the documented budget")
	}
	if acc.Model != "accepted-model" {
		t.Fatalf("AddChunk mutated model after rejecting the chunk: got %q", acc.Model)
	}
}

func TestAccumulatorRejectsSparseToolCallIndexWithoutMutation(t *testing.T) {
	var acc openai.ChatCompletionAccumulator
	chunk := accumulatorToolStringChunk("name", "arguments")
	chunk.Model = "rejected-model"
	chunk.Choices[0].Delta.ToolCalls[0].Index = 100_000

	if acc.AddChunk(chunk) {
		t.Fatal("AddChunk accepted a sparse tool-call index beyond the structural budget")
	}
	if acc.ID != "" || acc.Model != "" || len(acc.Choices) != 0 {
		t.Fatal("AddChunk mutated the accumulator after rejecting a sparse tool-call index")
	}
}

func TestAccumulatorStructuralBudgetBoundary(t *testing.T) {
	t.Run("accepts_limit", func(t *testing.T) {
		var acc openai.ChatCompletionAccumulator
		chunk := accumulatorToolStringChunk("name", "arguments")
		chunk.Choices[0].Delta.ToolCalls[0].Index = testAccumulatorMaxStructuralSlots - 2

		if !acc.AddChunk(chunk) {
			t.Fatal("AddChunk rejected the documented structural budget")
		}
		if len(acc.Choices) != 1 || len(acc.Choices[0].Message.ToolCalls) != testAccumulatorMaxStructuralSlots-1 {
			t.Fatal("AddChunk did not accumulate the structural budget boundary")
		}
	})

	t.Run("rejects_beyond_limit", func(t *testing.T) {
		var acc openai.ChatCompletionAccumulator
		chunk := accumulatorToolStringChunk("name", "arguments")
		chunk.Choices[0].Delta.ToolCalls[0].Index = testAccumulatorMaxStructuralSlots - 1

		if acc.AddChunk(chunk) {
			t.Fatal("AddChunk accepted structures beyond the documented budget")
		}
		if acc.ID != "" || len(acc.Choices) != 0 {
			t.Fatal("AddChunk mutated the accumulator after rejecting structures beyond the budget")
		}
	})
}

func TestAccumulatorDoesNotResurrectTruncatedPublicStrings(t *testing.T) {
	atLimit := strings.Repeat("x", testAccumulatorMaxTextBytes)

	t.Run("choices", func(t *testing.T) {
		var acc openai.ChatCompletionAccumulator
		if !acc.AddChunk(accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{Content: atLimit})) {
			t.Fatal("AddChunk rejected the initial chunk")
		}

		acc.Choices = acc.Choices[:0]
		if !acc.AddChunk(accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{Content: "x"})) {
			t.Fatal("AddChunk rejected text after the public choices were truncated")
		}
		if content := acc.Choices[0].Message.Content; content != "x" {
			t.Fatalf("AddChunk resurrected truncated choice content with length %d", len(content))
		}
	})

	t.Run("tool_calls", func(t *testing.T) {
		var acc openai.ChatCompletionAccumulator
		if !acc.AddChunk(accumulatorToolStringChunk("", atLimit)) {
			t.Fatal("AddChunk rejected the initial chunk")
		}

		message := &acc.Choices[0].Message
		message.ToolCalls = message.ToolCalls[:0]
		if !acc.AddChunk(accumulatorToolStringChunk("", "x")) {
			t.Fatal("AddChunk rejected text after the public tool calls were truncated")
		}
		if arguments := acc.Choices[0].Message.ToolCalls[0].Function.Arguments; arguments != "x" {
			t.Fatalf("AddChunk resurrected truncated tool arguments with length %d", len(arguments))
		}
	})
}

func TestAccumulatorReleasesTruncatedPublicBacking(t *testing.T) {
	halfLimit := strings.Repeat("x", testAccumulatorMaxTextBytes/2)

	t.Run("choices", func(t *testing.T) {
		var acc openai.ChatCompletionAccumulator
		chunk := accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{Content: halfLimit})
		secondChoice := chunk.Choices[0]
		secondChoice.Index = 1
		chunk.Choices = append(chunk.Choices, secondChoice)
		if !acc.AddChunk(chunk) {
			t.Fatal("AddChunk rejected the initial choices")
		}

		firstChoice := &acc.Choices[0]
		acc.Choices = acc.Choices[:1:1]
		if !acc.AddChunk(openai.ChatCompletionChunk{ID: chunk.ID}) {
			t.Fatal("AddChunk rejected the chunk after the public choices were truncated")
		}
		if &acc.Choices[0] == firstChoice {
			t.Fatal("capacity-clipped choices still use the truncated backing array")
		}
		if capacity := cap(acc.Choices); capacity != len(acc.Choices) {
			t.Fatalf("truncated choice backing was retained: length %d, capacity %d", len(acc.Choices), capacity)
		}
	})

	t.Run("tool_calls", func(t *testing.T) {
		var acc openai.ChatCompletionAccumulator
		chunk := accumulatorToolStringChunk("", halfLimit)
		secondToolCall := chunk.Choices[0].Delta.ToolCalls[0]
		secondToolCall.Index = 1
		chunk.Choices[0].Delta.ToolCalls = append(chunk.Choices[0].Delta.ToolCalls, secondToolCall)
		if !acc.AddChunk(chunk) {
			t.Fatal("AddChunk rejected the initial tool calls")
		}

		message := &acc.Choices[0].Message
		firstToolCall := &message.ToolCalls[0]
		message.ToolCalls = message.ToolCalls[:1:1]
		if !acc.AddChunk(openai.ChatCompletionChunk{ID: chunk.ID}) {
			t.Fatal("AddChunk rejected the chunk after the public tool calls were truncated")
		}
		if &message.ToolCalls[0] == firstToolCall {
			t.Fatal("capacity-clipped tool calls still use the truncated backing array")
		}
		if capacity := cap(message.ToolCalls); capacity != len(message.ToolCalls) {
			t.Fatalf("truncated tool-call backing was retained: length %d, capacity %d", len(message.ToolCalls), capacity)
		}
	})
}

func TestAccumulatorTruncationInvalidatesFinishedToolCallState(t *testing.T) {
	t.Run("legacy_tool_calls", func(t *testing.T) {
		var acc openai.ChatCompletionAccumulator
		initial := accumulatorToolStringChunk("tool", "arguments")
		if !acc.AddChunk(initial) {
			t.Fatal("AddChunk rejected the initial tool call")
		}

		acc.Choices[0].Message.ToolCalls = acc.Choices[0].Message.ToolCalls[:0]
		if !acc.AddChunk(accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{})) {
			t.Fatal("AddChunk rejected the chunk after the public tool calls were truncated")
		}
		if _, ok := acc.JustFinishedToolCall(); ok {
			t.Fatal("JustFinishedToolCall reported a removed tool call")
		}
	})

	t.Run("choice_aware_choices", func(t *testing.T) {
		var acc openai.ChatCompletionAccumulator
		initial := accumulatorToolStringChunk("tool", "arguments")
		initial.Choices[0].Index = 1
		if !acc.AddChunk(initial) {
			t.Fatal("AddChunk rejected the initial sparse choice")
		}

		acc.Choices = acc.Choices[:1]
		next := accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{})
		next.Choices[0].Index = 1
		if !acc.AddChunk(next) {
			t.Fatal("AddChunk rejected the chunk after the public choices were truncated")
		}
		if _, ok := acc.JustFinishedToolCallForChoice(1); ok {
			t.Fatal("JustFinishedToolCallForChoice reported a removed tool call")
		}
	})
}

func BenchmarkAccumulatorOneByteChunks(b *testing.B) {
	for _, fragmentCount := range []int{128, 1_024, 8_192} {
		b.Run(fmt.Sprintf("fragments=%d", fragmentCount), func(b *testing.B) {
			chunk := openai.ChatCompletionChunk{
				ID: "chatcmpl-benchmark",
				Choices: []openai.ChatCompletionChunkChoice{{
					Delta: openai.ChatCompletionChunkChoiceDelta{Content: "x"},
				}},
			}

			b.ReportAllocs()
			for b.Loop() {
				var acc openai.ChatCompletionAccumulator
				for range fragmentCount {
					if !acc.AddChunk(chunk) {
						b.Fatal("AddChunk rejected a benchmark chunk")
					}
				}
			}
		})
	}
}

func accumulatorStringChunk(delta openai.ChatCompletionChunkChoiceDelta) openai.ChatCompletionChunk {
	return openai.ChatCompletionChunk{
		ID: "chatcmpl-public-string-budget",
		Choices: []openai.ChatCompletionChunkChoice{{
			Delta: delta,
		}},
	}
}

func accumulatorToolStringChunk(name, arguments string) openai.ChatCompletionChunk {
	return accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{
		ToolCalls: []openai.ChatCompletionChunkChoiceDeltaToolCall{{
			Index: 0,
			Function: openai.ChatCompletionChunkChoiceDeltaToolCallFunction{
				Name:      name,
				Arguments: arguments,
			},
		}},
	})
}

func assertAccumulatorStrings(t *testing.T, acc *openai.ChatCompletionAccumulator, content, refusal, name, arguments string) {
	t.Helper()

	message := acc.Choices[0].Message
	if message.Content != content || message.Refusal != refusal ||
		message.ToolCalls[0].Function.Name != name ||
		message.ToolCalls[0].Function.Arguments != arguments {
		t.Fatalf(
			"accumulated strings: got content %q, refusal %q, name %q, arguments %q",
			message.Content,
			message.Refusal,
			message.ToolCalls[0].Function.Name,
			message.ToolCalls[0].Function.Arguments,
		)
	}
}

package openai_test

import (
	"fmt"
	"strings"
	"testing"

	openai "github.com/openai/openai-go/v3"
)

const (
	testAccumulatorMaxChunks    = 100_000
	testAccumulatorMaxTextBytes = 16 << 20
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

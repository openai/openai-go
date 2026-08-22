package openai_test

import (
	"strings"
	"testing"
	"unsafe"

	openai "github.com/openai/openai-go/v3"
)

func TestAccumulatorRejectsLogprobsBeyondBudgetWithoutMutation(t *testing.T) {
	var acc openai.ChatCompletionAccumulator
	toolChunk := accumulatorToolStringChunk("tool", "arguments")
	if !acc.AddChunk(toolChunk) {
		t.Fatal("AddChunk rejected the initial tool-call chunk")
	}

	logprobOverhead := int(unsafe.Sizeof(openai.ChatCompletionTokenLogprob{}))
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
	if _, ok := acc.JustFinishedToolCall(); ok {
		t.Fatal("AddChunk retained a completion event after rejecting excessive logprobs")
	}

	afterRejection := accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{Content: "-after"})
	if !acc.AddChunk(afterRejection) {
		t.Fatal("AddChunk consumed unrelated text or chunk budget after rejecting excessive logprobs")
	}
	if got := acc.Choices[0].Message.Content; got != "-after" {
		t.Fatalf("content after rejected logprobs = %q, want %q", got, "-after")
	}

	afterClearing := accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{})
	afterClearing.Choices[0].Logprobs.Content = []openai.ChatCompletionTokenLogprob{{Token: "x"}}
	acc.Choices[0].Logprobs.Content = acc.Choices[0].Logprobs.Content[:0]
	if !acc.AddChunk(accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{})) {
		t.Fatal("AddChunk did not release retained logprob backing after the public slice was truncated")
	}
	if got := len(acc.Choices[0].Logprobs.Content); got != 0 {
		t.Fatalf("logprobs after truncation = %d, want 0", got)
	}
	if acc.AddChunk(afterClearing) {
		t.Fatal("AddChunk replenished the cumulative remote-logprob budget after public truncation")
	}
	if got := len(acc.Choices[0].Logprobs.Content); got != 0 {
		t.Fatalf("rejected logprobs changed accumulated length to %d", got)
	}
}

func TestAccumulatorBudgetsPublicLogprobReplacement(t *testing.T) {
	var acc openai.ChatCompletionAccumulator
	if !acc.AddChunk(accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{})) {
		t.Fatal("AddChunk rejected the initial chunk")
	}

	logprobOverhead := int(unsafe.Sizeof(openai.ChatCompletionTokenLogprob{}))
	byteCount := (testAccumulatorMaxLogprobBytes - logprobOverhead) / int(unsafe.Sizeof(int64(0)))
	remaining := testAccumulatorMaxLogprobBytes - logprobOverhead - byteCount*int(unsafe.Sizeof(int64(0)))
	acc.Choices[0].Logprobs.Content = []openai.ChatCompletionTokenLogprob{{
		Token: strings.Repeat("x", remaining),
		Bytes: make([]int64, byteCount),
	}}

	beyondLimit := accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{})
	beyondLimit.Choices[0].Logprobs.Content = []openai.ChatCompletionTokenLogprob{{Token: "x"}}
	if acc.AddChunk(beyondLimit) {
		t.Fatal("AddChunk accepted logprobs beyond the live public logprob budget")
	}
	if got := len(acc.Choices[0].Logprobs.Content); got != 1 {
		t.Fatalf("public replacement changed after rejection: got %d logprobs, want 1", got)
	}
}

func TestAccumulatorBudgetsPublicLogprobInsertionBeforeTextChunk(t *testing.T) {
	var acc openai.ChatCompletionAccumulator
	if !acc.AddChunk(accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{})) {
		t.Fatal("AddChunk rejected the initial chunk")
	}

	logprobOverhead := int(unsafe.Sizeof(openai.ChatCompletionTokenLogprob{}))
	byteCount := (testAccumulatorMaxLogprobBytes-logprobOverhead)/int(unsafe.Sizeof(int64(0))) + 1
	acc.Choices[0].Logprobs.Content = []openai.ChatCompletionTokenLogprob{{Bytes: make([]int64, byteCount)}}

	textChunk := accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{Content: "text"})
	if acc.AddChunk(textChunk) {
		t.Fatal("AddChunk accepted text while a supported public logprob insertion exceeded the budget")
	}
	if got := acc.Choices[0].Message.Content; got != "" {
		t.Fatalf("rejected chunk changed content to %q", got)
	}
}

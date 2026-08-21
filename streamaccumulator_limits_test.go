package openai_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"unsafe"

	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/packages/ssestream"
)

const (
	testAccumulatorMaxChunks          = 100_000
	testAccumulatorMaxStructuralSlots = 1_024
	testAccumulatorMaxTextBytes       = 16 << 20
	testAccumulatorMaxMetadataBytes   = 16 << 20
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

func TestAccumulatorRejectsRetainedToolMetadataBeyondBudgetWithoutMutation(t *testing.T) {
	tests := []struct {
		name string
		set  func(*openai.ChatCompletionChunkChoiceDeltaToolCall, string)
	}{
		{
			name: "id",
			set: func(toolCall *openai.ChatCompletionChunkChoiceDeltaToolCall, text string) {
				toolCall.ID = text
			},
		},
		{
			name: "type",
			set: func(toolCall *openai.ChatCompletionChunkChoiceDeltaToolCall, text string) {
				toolCall.Type = text
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var acc openai.ChatCompletionAccumulator
			atLimit := accumulatorToolStringChunk("", "")
			atLimit.Created = 1
			test.set(&atLimit.Choices[0].Delta.ToolCalls[0], strings.Repeat("x", testAccumulatorMaxMetadataBytes-len(atLimit.ID)))
			if !acc.AddChunk(atLimit) {
				t.Fatal("AddChunk rejected retained tool metadata at the documented aggregate budget")
			}

			beyondLimit := accumulatorToolStringChunk("", "")
			beyondLimit.Created = 2
			beyondLimit.Choices[0].Delta.ToolCalls[0].Index = 1
			test.set(&beyondLimit.Choices[0].Delta.ToolCalls[0], "x")
			if acc.AddChunk(beyondLimit) {
				t.Fatal("AddChunk accepted retained tool metadata beyond the documented aggregate budget")
			}
			if acc.Created != 1 || len(acc.Choices[0].Message.ToolCalls) != 1 {
				t.Fatal("AddChunk mutated the accumulator after rejecting excessive retained tool metadata")
			}
		})
	}
}

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

	afterClearing := accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{})
	afterClearing.Choices[0].Logprobs.Content = []openai.ChatCompletionTokenLogprob{{Token: "x"}}
	acc.Choices[0].Logprobs.Content = acc.Choices[0].Logprobs.Content[:0]
	if acc.AddChunk(afterClearing) {
		t.Fatal("AddChunk ignored retained logprob backing after the public slice was truncated")
	}

	acc.Choices[0].Logprobs.Content = nil
	if !acc.AddChunk(afterClearing) {
		t.Fatal("AddChunk did not recover logprob budget after the public logprobs were cleared")
	}
	if got := len(acc.Choices[0].Logprobs.Content); got != 1 {
		t.Fatalf("logprobs after clearing = %d, want 1", got)
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

func TestAccumulatorDetachesCapacityClippedLogprobBacking(t *testing.T) {
	tests := []struct {
		start   int
		visible int
	}{
		{start: 0, visible: 0},
		{start: 0, visible: 1},
		{start: 1, visible: 0},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("start=%d/visible=%d", test.start, test.visible), func(t *testing.T) {
			var acc openai.ChatCompletionAccumulator
			initial := accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{})
			initial.Choices[0].Logprobs.Content = []openai.ChatCompletionTokenLogprob{{Token: "visible"}, {Token: "hidden"}}
			if !acc.AddChunk(initial) {
				t.Fatal("AddChunk rejected the initial logprobs")
			}

			before := unsafe.SliceData(acc.Choices[0].Logprobs.Content)
			end := test.start + test.visible
			acc.Choices[0].Logprobs.Content = acc.Choices[0].Logprobs.Content[test.start:end:end]
			if !acc.AddChunk(accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{})) {
				t.Fatal("AddChunk rejected the chunk after the public logprobs were capacity-clipped")
			}
			if after := unsafe.SliceData(acc.Choices[0].Logprobs.Content); after == before {
				t.Fatal("capacity-clipped public logprobs still retain the previous backing array")
			}
			if got := len(acc.Choices[0].Logprobs.Content); got != test.visible {
				t.Fatalf("visible logprobs after detachment = %d, want %d", got, test.visible)
			}
		})
	}
}

func TestAccumulatorDetachesDecodedChunkBacking(t *testing.T) {
	payload := fmt.Sprintf(`{
		"id":"chatcmpl-detach-logprob-json",
		"choices":[{
			"delta":{"role":"assistant","content":"content","tool_calls":[{"index":0,"id":"tool-id","type":"function","function":{"name":"tool-name","arguments":"{}"}}]},
			"finish_reason":"stop",
			"index":0,
			"logprobs":{"content":[{"token":"token","bytes":[],"logprob":0,"top_logprobs":[]}],"refusal":[]}
		}],
		"created":0,
		"model":"model",
		"object":"chat.completion.chunk",
		"service_tier":"default",
		"system_fingerprint":"fingerprint",
		"ignored":"%s"
	}`, strings.Repeat("x", 1<<20))

	decoder := ssestream.NewDecoder(&http.Response{
		Body: io.NopCloser(strings.NewReader("data: " + strings.ReplaceAll(payload, "\n", "") + "\n\n")),
	})
	defer func() {
		if err := decoder.Close(); err != nil {
			t.Errorf("close decoder: %v", err)
		}
	}()
	if !decoder.Next() {
		t.Fatalf("decode event: %v", decoder.Err())
	}
	var chunk openai.ChatCompletionChunk
	if err := json.Unmarshal(decoder.Event().Data, &chunk); err != nil {
		t.Fatalf("unmarshal chunk: %v", err)
	}
	source := chunk.Choices[0].Logprobs.Content[0]
	var acc openai.ChatCompletionAccumulator
	if !acc.AddChunk(chunk) {
		t.Fatal("AddChunk rejected the decoded logprob chunk")
	}

	retained := acc.Choices[0].Logprobs.Content[0]
	if unsafe.StringData(retained.Token) == unsafe.StringData(source.Token) {
		t.Fatal("accumulated token still aliases the decoded chunk backing")
	}
	if unsafe.StringData(retained.RawJSON()) == unsafe.StringData(source.RawJSON()) {
		t.Fatal("accumulated logprob metadata still aliases the decoded event backing")
	}
	metadata := []struct {
		name   string
		got    string
		source string
	}{
		{name: "id", got: acc.ID, source: chunk.ID},
		{name: "model", got: acc.Model, source: chunk.Model},
		{name: "system_fingerprint", got: acc.SystemFingerprint, source: chunk.SystemFingerprint},
		{name: "service_tier", got: string(acc.ServiceTier), source: string(chunk.ServiceTier)},
		{name: "finish_reason", got: acc.Choices[0].FinishReason, source: chunk.Choices[0].FinishReason},
		{name: "role", got: string(acc.Choices[0].Message.Role), source: chunk.Choices[0].Delta.Role},
		{name: "tool_id", got: acc.Choices[0].Message.ToolCalls[0].ID, source: chunk.Choices[0].Delta.ToolCalls[0].ID},
		{name: "tool_type", got: acc.Choices[0].Message.ToolCalls[0].Type, source: chunk.Choices[0].Delta.ToolCalls[0].Type},
	}
	for _, field := range metadata {
		if unsafe.StringData(field.got) == unsafe.StringData(field.source) {
			t.Errorf("accumulated %s still aliases the decoded event backing", field.name)
		}
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

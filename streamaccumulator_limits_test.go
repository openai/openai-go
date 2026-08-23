package openai_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
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

		t.Run(test.name+"/clearing_does_not_replenish_output_budget", func(t *testing.T) {
			var acc openai.ChatCompletionAccumulator
			chunk := test.chunk(atLimit)
			chunk.Model = "initial-model"
			if !acc.AddChunk(chunk) {
				t.Fatal("AddChunk rejected text at the documented aggregate budget")
			}

			*test.value(&acc) = ""
			chunk = test.chunk("x")
			chunk.Model = "rejected-model"
			if acc.AddChunk(chunk) {
				t.Fatal("AddChunk replenished the cumulative remote-output budget after the public string was cleared")
			}
			if acc.Model != "initial-model" || *test.value(&acc) != "" {
				t.Fatal("AddChunk mutated the accumulator after rejecting remote text beyond the cumulative budget")
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

func TestAccumulatorPublicIDMutationContract(t *testing.T) {
	chunk := accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{})
	var acc openai.ChatCompletionAccumulator
	if !acc.AddChunk(chunk) {
		t.Fatal("AddChunk rejected the initial chunk")
	}

	acc.ID = ""
	differentStream := chunk
	differentStream.ID = "chatcmpl-different-stream"
	differentStream.Model = "rejected-model"
	if acc.AddChunk(differentStream) {
		t.Fatal("AddChunk accepted a different stream after the public ID was cleared")
	}
	if acc.ID != "" || acc.Model == "rejected-model" {
		t.Fatal("AddChunk mutated the accumulator after rejecting the different stream")
	}
	if !acc.AddChunk(chunk) || acc.ID != chunk.ID {
		t.Fatal("AddChunk did not restore a cleared public ID")
	}

	acc.ID = "different-id"
	if acc.AddChunk(chunk) {
		t.Fatal("AddChunk accepted a different nonempty public ID")
	}
	if acc.ID != "different-id" {
		t.Fatal("AddChunk changed the rejected public ID replacement")
	}
}

func TestAccumulatorChargesOnlyClonedNestedLogprobStorage(t *testing.T) {
	tests := []struct {
		name string
		set  func(*openai.ChatCompletionTokenLogprob)
	}{
		{
			name: "bytes",
			set: func(logprob *openai.ChatCompletionTokenLogprob) {
				logprob.Bytes = make([]int64, 1, testAccumulatorMaxLogprobBytes/int(unsafe.Sizeof(int64(0))))
			},
		},
		{
			name: "top_logprobs",
			set: func(logprob *openai.ChatCompletionTokenLogprob) {
				logprob.TopLogprobs = make([]openai.ChatCompletionTokenLogprobTopLogprob, 1,
					testAccumulatorMaxLogprobBytes/int(unsafe.Sizeof(openai.ChatCompletionTokenLogprobTopLogprob{})))
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			chunk := accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{})
			chunk.Choices[0].Logprobs.Content = []openai.ChatCompletionTokenLogprob{{Token: "token"}}
			test.set(&chunk.Choices[0].Logprobs.Content[0])

			var acc openai.ChatCompletionAccumulator
			if !acc.AddChunk(chunk) {
				t.Fatal("AddChunk charged nested spare capacity that the accumulator does not retain")
			}
			retained := acc.Choices[0].Logprobs.Content[0]
			if cap(retained.Bytes) != len(retained.Bytes) || cap(retained.TopLogprobs) != len(retained.TopLogprobs) {
				t.Fatal("accumulated logprob retained nested spare capacity")
			}
		})
	}
}

func TestAccumulatorLogprobReconciliationIsAtomic(t *testing.T) {
	var acc openai.ChatCompletionAccumulator
	initial := accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{})
	initial.Choices[0].Logprobs.Content = []openai.ChatCompletionTokenLogprob{{Token: "visible"}, {Token: "hidden"}}
	if !acc.AddChunk(initial) {
		t.Fatal("AddChunk rejected the initial logprobs")
	}

	acc.Choices[0].Logprobs.Content = acc.Choices[0].Logprobs.Content[:0:0]
	contentData := unsafe.SliceData(acc.Choices[0].Logprobs.Content)
	byteCount := testAccumulatorMaxLogprobBytes / int(unsafe.Sizeof(int64(0)))
	acc.Choices[0].Logprobs.Refusal = []openai.ChatCompletionTokenLogprob{{Bytes: make([]int64, byteCount)}}
	if acc.AddChunk(accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{})) {
		t.Fatal("AddChunk accepted public logprobs beyond the aggregate budget")
	}
	content := acc.Choices[0].Logprobs.Content
	if len(content) != 0 || cap(content) != 0 || unsafe.SliceData(content) != contentData {
		t.Fatal("AddChunk detached public logprob backing before rejecting the chunk")
	}

	acc.Choices[0].Logprobs.Refusal = nil
	if !acc.AddChunk(accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{})) {
		t.Fatal("AddChunk rejected the valid chunk after excessive public logprobs were cleared")
	}
	if unsafe.SliceData(acc.Choices[0].Logprobs.Content) == contentData {
		t.Fatal("AddChunk did not apply the staged detachment after accepting the chunk")
	}
}

func TestAccumulatorMetadataUsesProjectedRetainedState(t *testing.T) {
	tests := []struct {
		name             string
		chunk            func() openai.ChatCompletionChunk
		setChunk         func(*openai.ChatCompletionChunk, string)
		get              func(*openai.ChatCompletionAccumulator) string
		clearAccumulator func(*openai.ChatCompletionAccumulator)
	}{
		{
			name: "model",
			chunk: func() openai.ChatCompletionChunk {
				return accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{})
			},
			setChunk: func(chunk *openai.ChatCompletionChunk, value string) { chunk.Model = value },
			get:      func(acc *openai.ChatCompletionAccumulator) string { return acc.Model },
		},
		{
			name: "system_fingerprint",
			chunk: func() openai.ChatCompletionChunk {
				return accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{})
			},
			setChunk: func(chunk *openai.ChatCompletionChunk, value string) { chunk.SystemFingerprint = value },
			get:      func(acc *openai.ChatCompletionAccumulator) string { return acc.SystemFingerprint },
		},
		{
			name: "finish_reason",
			chunk: func() openai.ChatCompletionChunk {
				return accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{})
			},
			setChunk: func(chunk *openai.ChatCompletionChunk, value string) { chunk.Choices[0].FinishReason = value },
			get:      func(acc *openai.ChatCompletionAccumulator) string { return acc.Choices[0].FinishReason },
		},
		{
			name: "role",
			chunk: func() openai.ChatCompletionChunk {
				return accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{})
			},
			setChunk:         func(chunk *openai.ChatCompletionChunk, value string) { chunk.Choices[0].Delta.Role = value },
			get:              func(acc *openai.ChatCompletionAccumulator) string { return string(acc.Choices[0].Message.Role) },
			clearAccumulator: func(acc *openai.ChatCompletionAccumulator) { acc.Choices[0].Message.Role = "" },
		},
		{
			name:             "tool_id",
			chunk:            func() openai.ChatCompletionChunk { return accumulatorToolStringChunk("", "") },
			setChunk:         func(chunk *openai.ChatCompletionChunk, value string) { chunk.Choices[0].Delta.ToolCalls[0].ID = value },
			get:              func(acc *openai.ChatCompletionAccumulator) string { return acc.Choices[0].Message.ToolCalls[0].ID },
			clearAccumulator: func(acc *openai.ChatCompletionAccumulator) { acc.Choices[0].Message.ToolCalls[0].ID = "" },
		},
		{
			name:  "tool_type",
			chunk: func() openai.ChatCompletionChunk { return accumulatorToolStringChunk("", "") },
			setChunk: func(chunk *openai.ChatCompletionChunk, value string) {
				chunk.Choices[0].Delta.ToolCalls[0].Type = value
			},
			get:              func(acc *openai.ChatCompletionAccumulator) string { return acc.Choices[0].Message.ToolCalls[0].Type },
			clearAccumulator: func(acc *openai.ChatCompletionAccumulator) { acc.Choices[0].Message.ToolCalls[0].Type = "" },
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			chunk := test.chunk()
			atLimit := strings.Repeat("x", testAccumulatorMaxMetadataBytes-len(chunk.ID))
			test.setChunk(&chunk, atLimit)
			var acc openai.ChatCompletionAccumulator
			if !acc.AddChunk(chunk) {
				t.Fatal("AddChunk rejected metadata at the aggregate budget")
			}
			if !acc.AddChunk(chunk) {
				t.Fatal("AddChunk double-counted repeated replacement metadata")
			}

			test.setChunk(&chunk, "replacement")
			if !acc.AddChunk(chunk) || test.get(&acc) != "replacement" {
				t.Fatal("AddChunk did not accept shorter replacement metadata")
			}

			if test.clearAccumulator != nil {
				test.clearAccumulator(&acc)
			} else {
				test.setChunk(&chunk, "")
			}
			if test.clearAccumulator != nil {
				test.setChunk(&chunk, "")
			}
			if !acc.AddChunk(chunk) || test.get(&acc) != "" {
				t.Fatal("AddChunk did not preserve cleared replacement metadata")
			}
		})
	}
}

func TestAccumulatorMetadataProjectionUsesLastWriteSemantics(t *testing.T) {
	tests := []struct {
		name string
		set  func(*openai.ChatCompletionChunkChoice, string)
		get  func(*openai.ChatCompletionAccumulator) string
	}{
		{
			name: "finish_reason",
			set:  func(choice *openai.ChatCompletionChunkChoice, value string) { choice.FinishReason = value },
			get:  func(acc *openai.ChatCompletionAccumulator) string { return acc.Choices[0].FinishReason },
		},
		{
			name: "tool_id",
			set:  func(choice *openai.ChatCompletionChunkChoice, value string) { choice.Delta.ToolCalls[0].ID = value },
			get:  func(acc *openai.ChatCompletionAccumulator) string { return acc.Choices[0].Message.ToolCalls[0].ID },
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			chunk := accumulatorToolStringChunk("", "")
			test.set(&chunk.Choices[0], strings.Repeat("x", testAccumulatorMaxMetadataBytes-len(chunk.ID)))
			last := chunk.Choices[0]
			test.set(&last, "final")
			chunk.Choices = append(chunk.Choices, last)

			var acc openai.ChatCompletionAccumulator
			if !acc.AddChunk(chunk) {
				t.Fatal("AddChunk charged metadata overwritten by duplicate deltas")
			}
			if got := test.get(&acc); got != "final" {
				t.Fatalf("last metadata value = %q, want final", got)
			}
		})
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

func TestAccumulatorDoesNotReplenishLogprobBudgetAfterDetachment(t *testing.T) {
	logprobOverhead := 2 * int(unsafe.Sizeof(openai.ChatCompletionTokenLogprob{}))
	initial := accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{})
	initial.Choices[0].Logprobs.Content = []openai.ChatCompletionTokenLogprob{
		{Token: strings.Repeat("x", testAccumulatorMaxLogprobBytes-logprobOverhead)},
		{},
	}

	var acc openai.ChatCompletionAccumulator
	if !acc.AddChunk(initial) {
		t.Fatal("AddChunk rejected initial logprobs at the aggregate budget")
	}
	acc.Choices[0].Logprobs.Content = acc.Choices[0].Logprobs.Content[:0:1]
	if !acc.AddChunk(accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{})) {
		t.Fatal("AddChunk rejected the empty chunk that detaches hidden logprobs")
	}
	if got := len(acc.Choices[0].Logprobs.Content); got != 0 {
		t.Fatalf("visible logprobs after detachment = %d, want 0", got)
	}

	next := accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{})
	next.Choices[0].Logprobs.Content = []openai.ChatCompletionTokenLogprob{{Token: "replacement"}}
	if acc.AddChunk(next) {
		t.Fatal("AddChunk replenished the cumulative remote-logprob budget after detachment")
	}
	if got := len(acc.Choices[0].Logprobs.Content); got != 0 {
		t.Fatalf("rejected logprobs changed accumulated length to %d", got)
	}
}

func TestAccumulatorCanonicalizesWholeLogprobSliceReplacement(t *testing.T) {
	var acc openai.ChatCompletionAccumulator
	initial := accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{})
	initial.Choices[0].Logprobs.Content = []openai.ChatCompletionTokenLogprob{{Token: "initial"}}
	if !acc.AddChunk(initial) {
		t.Fatal("AddChunk rejected the initial logprob")
	}

	logprobOverhead := 2 * int(unsafe.Sizeof(openai.ChatCompletionTokenLogprob{}))
	replacement := make([]openai.ChatCompletionTokenLogprob, 1, 2)
	replacement[0].Token = "visible"
	hidden := replacement[:2]
	hidden[1].Token = strings.Repeat("x", testAccumulatorMaxLogprobBytes-logprobOverhead-len(replacement[0].Token))
	acc.Choices[0].Logprobs.Content = replacement
	if !acc.AddChunk(accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{})) {
		t.Fatal("AddChunk rejected a supported whole-slice replacement at the aggregate budget")
	}
	if capacity := cap(acc.Choices[0].Logprobs.Content); capacity != len(replacement) {
		t.Fatalf("canonical replacement capacity = %d, want %d", capacity, len(replacement))
	}

	next := accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{})
	next.Choices[0].Logprobs.Content = []openai.ChatCompletionTokenLogprob{{Token: "next"}}
	if !acc.AddChunk(next) {
		t.Fatal("AddChunk charged hidden replacement data that canonicalization releases")
	}
	if got := len(acc.Choices[0].Logprobs.Content); got != 2 {
		t.Fatalf("accumulated logprobs = %d, want 2", got)
	}

	largeCapacity := testAccumulatorMaxLogprobBytes / int(unsafe.Sizeof(openai.ChatCompletionTokenLogprob{}))
	acc.Choices[0].Logprobs.Content = make([]openai.ChatCompletionTokenLogprob, 0, largeCapacity)
	if !acc.AddChunk(accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{})) {
		t.Fatal("AddChunk rejected an empty whole-slice replacement with spare capacity")
	}
	if capacity := cap(acc.Choices[0].Logprobs.Content); capacity != 0 {
		t.Fatalf("empty canonical replacement capacity = %d, want 0", capacity)
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

func TestAccumulatorPreservesNilNestedLogprobSlices(t *testing.T) {
	chunk := accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{})
	chunk.Choices[0].Logprobs.Content = []openai.ChatCompletionTokenLogprob{{Token: "token"}}

	var acc openai.ChatCompletionAccumulator
	if !acc.AddChunk(chunk) {
		t.Fatal("AddChunk rejected the logprob chunk")
	}
	logprob := acc.Choices[0].Logprobs.Content[0]
	if logprob.Bytes != nil || logprob.TopLogprobs != nil {
		t.Fatal("AddChunk changed nil nested logprob slices into non-nil empty slices")
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
		chunk := accumulatorDenseToolStringChunk(testAccumulatorMaxStructuralSlots - 1)

		if !acc.AddChunk(chunk) {
			t.Fatal("AddChunk rejected the documented structural budget")
		}
		if len(acc.Choices) != 1 || len(acc.Choices[0].Message.ToolCalls) != testAccumulatorMaxStructuralSlots-1 {
			t.Fatal("AddChunk did not accumulate the structural budget boundary")
		}
	})

	t.Run("rejects_beyond_limit", func(t *testing.T) {
		var acc openai.ChatCompletionAccumulator
		chunk := accumulatorDenseToolStringChunk(testAccumulatorMaxStructuralSlots)

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
		empty := accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{})
		if !acc.AddChunk(empty) {
			t.Fatal("AddChunk rejected the empty chunk after the public choices were truncated")
		}
		if content := acc.Choices[0].Message.Content; content != "" {
			t.Fatalf("AddChunk resurrected truncated choice content with length %d", len(content))
		}
		if acc.AddChunk(accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{Content: "x"})) {
			t.Fatal("AddChunk replenished the cumulative remote-output budget after choices were truncated")
		}
	})

	t.Run("tool_calls", func(t *testing.T) {
		var acc openai.ChatCompletionAccumulator
		if !acc.AddChunk(accumulatorToolStringChunk("", atLimit)) {
			t.Fatal("AddChunk rejected the initial chunk")
		}

		message := &acc.Choices[0].Message
		message.ToolCalls = message.ToolCalls[:0]
		empty := accumulatorStringChunk(openai.ChatCompletionChunkChoiceDelta{})
		if !acc.AddChunk(empty) {
			t.Fatal("AddChunk rejected the empty chunk after the public tool calls were truncated")
		}
		if len(acc.Choices[0].Message.ToolCalls) != 0 {
			t.Fatal("AddChunk resurrected a truncated tool call")
		}
		if acc.AddChunk(accumulatorToolStringChunk("", "x")) {
			t.Fatal("AddChunk replenished the cumulative remote-output budget after tool calls were truncated")
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

func accumulatorDenseToolStringChunk(count int) openai.ChatCompletionChunk {
	chunk := accumulatorToolStringChunk("name", "arguments")
	toolCall := chunk.Choices[0].Delta.ToolCalls[0]
	chunk.Choices[0].Delta.ToolCalls = make([]openai.ChatCompletionChunkChoiceDeltaToolCall, count)
	for i := range chunk.Choices[0].Delta.ToolCalls {
		toolCall.Index = int64(i)
		chunk.Choices[0].Delta.ToolCalls[i] = toolCall
	}
	return chunk
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

package openai

import (
	"runtime"
	"strings"
	"testing"
	"unsafe"
	"weak"
)

func TestAccumulatorValueCopyChoiceActivationDoesNotPopulateDormantSpareBacking(t *testing.T) {
	var original ChatCompletionAccumulator
	for index := range 3 {
		chunk := storageTestChunk(ChatCompletionChunkChoiceDelta{})
		chunk.Choices[0].Index = int64(index)
		if !original.AddChunk(chunk) {
			t.Fatalf("AddChunk rejected initial choice %d", index)
		}
	}
	if length, capacity := len(original.stringState.choices), cap(original.stringState.choices); length != 3 || capacity <= length {
		t.Fatalf("initial private choice backing: len %d cap %d, want len 3 with spare capacity", length, capacity)
	}
	dormant := original
	dormant.Choices = cloneAccumulatorSlice(dormant.Choices)
	activation := storageTestChunk(ChatCompletionChunkChoiceDelta{Content: strings.Repeat("x", 32<<10)})
	activation.Choices[0].Index = 3
	if !original.AddChunk(activation) {
		t.Fatal("AddChunk rejected the original accumulator's dense choice activation")
	}
	if dormantBacking := dormant.stringState.choices[:cap(dormant.stringState.choices)]; dormantBacking[3] != nil {
		t.Fatal("dense choice activation populated private backing retained by a dormant accumulator copy")
	}
}

func TestAccumulatorValueCopyToolActivationDoesNotPopulateDormantSpareBacking(t *testing.T) {
	var original ChatCompletionAccumulator
	for index := range 3 {
		chunk := storageTestChunk(ChatCompletionChunkChoiceDelta{
			ToolCalls: []ChatCompletionChunkChoiceDeltaToolCall{{Index: int64(index)}},
		})
		if !original.AddChunk(chunk) {
			t.Fatalf("AddChunk rejected initial tool call %d", index)
		}
	}
	toolCalls := original.stringState.choices[0].toolCalls
	if length, capacity := len(toolCalls), cap(toolCalls); length != 3 || capacity <= length {
		t.Fatalf("initial private tool backing: len %d cap %d, want len 3 with spare capacity", length, capacity)
	}
	dormant := original
	dormant.Choices = cloneAccumulatorSlice(dormant.Choices)
	dormant.Choices[0].Message.ToolCalls = cloneAccumulatorSlice(dormant.Choices[0].Message.ToolCalls)
	activation := storageTestChunk(ChatCompletionChunkChoiceDelta{
		ToolCalls: []ChatCompletionChunkChoiceDeltaToolCall{{
			Index:    3,
			Function: ChatCompletionChunkChoiceDeltaToolCallFunction{Arguments: strings.Repeat("x", 32<<10)},
		}},
	})
	if !original.AddChunk(activation) {
		t.Fatal("AddChunk rejected the original accumulator's dense tool activation")
	}
	lifetime := weak.Make(original.stringState.choices[0].toolCallState(3))
	original = ChatCompletionAccumulator{}
	runtime.GC()
	if lifetime.Value() != nil {
		t.Fatal("dense tool activation retained private state through a dormant accumulator copy's spare backing")
	}
	if !dormant.AddChunk(activation) {
		t.Fatal("AddChunk rejected the dormant accumulator after the original was collected")
	}
}

func TestAccumulatorValueCopyDoesNotRetainOriginalLogprobGrowth(t *testing.T) {
	for _, refusal := range []bool{false, true} {
		name := "content"
		if refusal {
			name = "refusal"
		}
		t.Run(name, func(t *testing.T) {
			initial := storageTestChunk(ChatCompletionChunkChoiceDelta{})
			values := []ChatCompletionTokenLogprob{{Token: "existing"}}
			if refusal {
				initial.Choices[0].Logprobs.Refusal = values
			} else {
				initial.Choices[0].Logprobs.Content = values
			}
			var original ChatCompletionAccumulator
			if !original.AddChunk(initial) {
				t.Fatal("AddChunk rejected the initial logprob")
			}
			dormant := original
			dormant.Choices = cloneAccumulatorSlice(dormant.Choices)
			if refusal {
				dormant.Choices[0].Logprobs.Refusal = cloneAccumulatorSlice(dormant.Choices[0].Logprobs.Refusal)
			} else {
				dormant.Choices[0].Logprobs.Content = cloneAccumulatorSlice(dormant.Choices[0].Logprobs.Content)
			}
			growth := storageTestChunk(ChatCompletionChunkChoiceDelta{})
			values = []ChatCompletionTokenLogprob{{Token: strings.Repeat("x", 32<<10)}}
			if refusal {
				growth.Choices[0].Logprobs.Refusal = values
			} else {
				growth.Choices[0].Logprobs.Content = values
			}
			if !original.AddChunk(growth) {
				t.Fatal("AddChunk rejected growth of the original accumulator's existing logprob")
			}
			logprobs := original.Choices[0].Logprobs.Content
			if refusal {
				logprobs = original.Choices[0].Logprobs.Refusal
			}
			lifetime := weak.Make(&logprobs[1])
			original = ChatCompletionAccumulator{}
			logprobs = nil
			runtime.GC()
			if lifetime.Value() != nil {
				t.Fatal("dormant accumulator copy retained the original accumulator's later logprob backing")
			}
			runtime.KeepAlive(&dormant)
		})
	}
}

func TestAccumulatorValueCopyDoesNotRetainOriginalTextGrowth(t *testing.T) {
	tests := []struct {
		name  string
		chunk func(string) ChatCompletionChunk
		value func(*ChatCompletionAccumulator) string
	}{
		{
			name: "content",
			chunk: func(value string) ChatCompletionChunk {
				return storageTestChunk(ChatCompletionChunkChoiceDelta{Content: value})
			},
			value: func(acc *ChatCompletionAccumulator) string { return acc.Choices[0].Message.Content },
		},
		{
			name: "refusal",
			chunk: func(value string) ChatCompletionChunk {
				return storageTestChunk(ChatCompletionChunkChoiceDelta{Refusal: value})
			},
			value: func(acc *ChatCompletionAccumulator) string { return acc.Choices[0].Message.Refusal },
		},
		{
			name: "tool_name",
			chunk: func(value string) ChatCompletionChunk {
				return storageTestChunk(ChatCompletionChunkChoiceDelta{ToolCalls: []ChatCompletionChunkChoiceDeltaToolCall{{
					Function: ChatCompletionChunkChoiceDeltaToolCallFunction{Name: value},
				}}})
			},
			value: func(acc *ChatCompletionAccumulator) string { return acc.Choices[0].Message.ToolCalls[0].Function.Name },
		},
		{
			name: "tool_arguments",
			chunk: func(value string) ChatCompletionChunk {
				return storageTestChunk(ChatCompletionChunkChoiceDelta{ToolCalls: []ChatCompletionChunkChoiceDeltaToolCall{{
					Function: ChatCompletionChunkChoiceDeltaToolCallFunction{Arguments: value},
				}}})
			},
			value: func(acc *ChatCompletionAccumulator) string {
				return acc.Choices[0].Message.ToolCalls[0].Function.Arguments
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var original ChatCompletionAccumulator
			if !original.AddChunk(test.chunk("existing")) {
				t.Fatal("AddChunk rejected the initial text")
			}
			dormant := original
			dormant.Choices = cloneAccumulatorSlice(dormant.Choices)
			dormant.Choices[0].Message.ToolCalls = cloneAccumulatorSlice(dormant.Choices[0].Message.ToolCalls)
			if !original.AddChunk(test.chunk(strings.Repeat("x", 32<<10))) {
				t.Fatal("AddChunk rejected growth of the original accumulator's existing text")
			}
			lifetime := weak.Make(unsafe.StringData(test.value(&original)))
			original = ChatCompletionAccumulator{}
			runtime.GC()
			if lifetime.Value() != nil {
				t.Fatal("dormant accumulator copy retained the original accumulator's later text backing")
			}
			runtime.KeepAlive(&dormant)
		})
	}
}

func TestAccumulatorValueCopyDoesNotRetainOriginalMetadataReplacement(t *testing.T) {
	tests := []struct {
		name  string
		chunk func(string) ChatCompletionChunk
		value func(*ChatCompletionAccumulator) string
	}{
		{
			name: "finish_reason",
			chunk: func(value string) ChatCompletionChunk {
				chunk := storageTestChunk(ChatCompletionChunkChoiceDelta{})
				chunk.Choices[0].FinishReason = value
				return chunk
			},
			value: func(acc *ChatCompletionAccumulator) string { return acc.Choices[0].FinishReason },
		},
		{
			name: "role",
			chunk: func(value string) ChatCompletionChunk {
				return storageTestChunk(ChatCompletionChunkChoiceDelta{Role: value})
			},
			value: func(acc *ChatCompletionAccumulator) string { return string(acc.Choices[0].Message.Role) },
		},
		{
			name: "tool_id",
			chunk: func(value string) ChatCompletionChunk {
				return storageTestChunk(ChatCompletionChunkChoiceDelta{
					ToolCalls: []ChatCompletionChunkChoiceDeltaToolCall{{ID: value}},
				})
			},
			value: func(acc *ChatCompletionAccumulator) string { return acc.Choices[0].Message.ToolCalls[0].ID },
		},
		{
			name: "tool_type",
			chunk: func(value string) ChatCompletionChunk {
				return storageTestChunk(ChatCompletionChunkChoiceDelta{
					ToolCalls: []ChatCompletionChunkChoiceDeltaToolCall{{Type: value}},
				})
			},
			value: func(acc *ChatCompletionAccumulator) string { return acc.Choices[0].Message.ToolCalls[0].Type },
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var original ChatCompletionAccumulator
			if !original.AddChunk(test.chunk("existing")) {
				t.Fatal("AddChunk rejected the initial metadata")
			}
			dormant := original
			dormant.Choices = cloneAccumulatorSlice(dormant.Choices)
			dormant.Choices[0].Message.ToolCalls = cloneAccumulatorSlice(dormant.Choices[0].Message.ToolCalls)
			if !original.AddChunk(test.chunk(strings.Repeat("x", 32<<10))) {
				t.Fatal("AddChunk rejected replacement of the original accumulator's existing metadata")
			}
			lifetime := weak.Make(unsafe.StringData(test.value(&original)))
			original = ChatCompletionAccumulator{}
			runtime.GC()
			if lifetime.Value() != nil {
				t.Fatal("dormant accumulator copy retained the original accumulator's later metadata backing")
			}
			runtime.KeepAlive(&dormant)
		})
	}
}

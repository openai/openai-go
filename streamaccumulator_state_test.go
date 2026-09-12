package openai_test

import (
	"testing"

	openai "github.com/openai/openai-go/v3"
)

const (
	pr688EmptyToolCallsChunk = `{"id":"test","choices":[{"index":0,"delta":{"tool_calls":[]}}]}`
	pr688NullToolCallsChunk  = `{"id":"test","choices":[{"index":0,"delta":{"tool_calls":null}}]}`
	pr688ToolCallChunk       = `{"id":"test","choices":[{"index":0,"delta":{"tool_calls":[{"index":0,"id":"call_123","type":"function","function":{"name":"test_func","arguments":"{}"}}]}}]}`
	pr688FinishedChunk       = `{"id":"test","choices":[{"index":0,"delta":{},"finish_reason":"tool_calls"}]}`
)

func TestAccumulatorEmptyAndNullToolCallStateTransitions(t *testing.T) {
	for _, tc := range []struct {
		name  string
		chunk string
	}{
		{name: "empty array", chunk: pr688EmptyToolCallsChunk},
		{name: "null", chunk: pr688NullToolCallsChunk},
	} {
		tc := tc

		t.Run(tc.name+"/without active tool call", func(t *testing.T) {
			var acc openai.ChatCompletionAccumulator

			pr688AddChunk(t, &acc, tc.chunk)
			pr688AssertNoFinishedToolCall(t, &acc)

			pr688AddChunk(t, &acc, pr688FinishedChunk)
			pr688AssertNoFinishedToolCall(t, &acc)
		})

		t.Run(tc.name+"/before active tool call", func(t *testing.T) {
			var acc openai.ChatCompletionAccumulator

			pr688AddChunk(t, &acc, tc.chunk)
			pr688AssertNoFinishedToolCall(t, &acc)

			pr688AddChunk(t, &acc, pr688ToolCallChunk)
			pr688AssertNoFinishedToolCall(t, &acc)

			pr688AddChunk(t, &acc, pr688FinishedChunk)
			pr688AssertFinishedToolCall(t, &acc)

			pr688AddChunk(t, &acc, pr688FinishedChunk)
			pr688AssertNoFinishedToolCall(t, &acc)
		})

		t.Run(tc.name+"/after active tool call", func(t *testing.T) {
			var acc openai.ChatCompletionAccumulator

			pr688AddChunk(t, &acc, pr688ToolCallChunk)
			pr688AssertNoFinishedToolCall(t, &acc)

			pr688AddChunk(t, &acc, tc.chunk)
			pr688AssertFinishedToolCall(t, &acc)

			pr688AddChunk(t, &acc, pr688FinishedChunk)
			pr688AssertNoFinishedToolCall(t, &acc)
		})
	}
}

func TestAccumulatorTracksFinishedEventsForEveryChoice(t *testing.T) {
	var acc openai.ChatCompletionAccumulator

	pr688AddChunk(t, &acc, `{"id":"test","choices":[{"index":0,"delta":{"content":"choice zero"}},{"index":1,"delta":{"refusal":"choice one refusal"}},{"index":2,"delta":{"tool_calls":[{"index":1,"id":"call_choice_two","type":"function","function":{"name":"choice_two_tool","arguments":"{\"value\":2}"}}]}}]}`)
	pr688AddChunk(t, &acc, `{"id":"test","choices":[{"index":0,"delta":{"content":" continues"}},{"index":1,"delta":{},"finish_reason":"stop"},{"index":2,"delta":{},"finish_reason":"tool_calls"}]}`)

	if content, ok := acc.JustFinishedContent(); ok {
		t.Fatalf("JustFinishedContent returned unexpected content: %q", content)
	}

	if refusal, ok := acc.JustFinishedRefusal(); ok {
		t.Fatalf("JustFinishedRefusal returned alternate choice refusal: %q", refusal)
	}
	refusal, ok := acc.JustFinishedRefusalForChoice(1)
	if !ok {
		t.Fatal("JustFinishedRefusal did not return the completed refusal")
	}
	if refusal != "choice one refusal" {
		t.Fatalf("refusal: expected %q, got %q", "choice one refusal", refusal)
	}

	if toolCall, toolCallOK := acc.JustFinishedToolCall(); toolCallOK {
		t.Fatalf("JustFinishedToolCall returned alternate choice tool call: %+v", toolCall)
	}
	toolCall, ok := acc.JustFinishedToolCallForChoice(2)
	if !ok {
		t.Fatal("JustFinishedToolCall did not return the completed tool call")
	}
	if toolCall.Index != 1 {
		t.Fatalf("tool call index: expected 1, got %d", toolCall.Index)
	}
	if toolCall.ID != "call_choice_two" {
		t.Fatalf("tool call ID: expected %q, got %q", "call_choice_two", toolCall.ID)
	}
	if toolCall.Name != "choice_two_tool" {
		t.Fatalf("tool call name: expected %q, got %q", "choice_two_tool", toolCall.Name)
	}
	if toolCall.Arguments != `{"value":2}` {
		t.Fatalf("tool call arguments: expected %q, got %q", `{"value":2}`, toolCall.Arguments)
	}

	pr688AddChunk(t, &acc, `{"id":"test","choices":[{"index":0,"delta":{},"finish_reason":"stop"},{"index":1,"delta":{},"finish_reason":"stop"},{"index":2,"delta":{},"finish_reason":"tool_calls"}]}`)

	content, ok := acc.JustFinishedContent()
	if !ok {
		t.Fatal("JustFinishedContent did not return the completed content")
	}
	if content != "choice zero continues" {
		t.Fatalf("content: expected %q, got %q", "choice zero continues", content)
	}
	if refusal, ok := acc.JustFinishedRefusalForChoice(1); ok {
		t.Fatalf("JustFinishedRefusal returned repeated refusal: %q", refusal)
	}
	if toolCall, ok := acc.JustFinishedToolCallForChoice(2); ok {
		t.Fatalf("JustFinishedToolCall returned repeated tool call: %+v", toolCall)
	}
}

func TestAccumulatorSingularFinishedEventsUseFirstWireChoice(t *testing.T) {
	var acc openai.ChatCompletionAccumulator

	pr688AddChunk(t, &acc, `{"id":"test","choices":[{"index":0,"delta":{"content":"selected"}},{"index":1,"delta":{"tool_calls":[{"index":0,"id":"alternate_tool","type":"function","function":{"name":"mutate","arguments":"{}"}}]}},{"index":2,"delta":{"refusal":"alternate refusal"}}]}`)
	pr688AddChunk(t, &acc, `{"id":"test","choices":[{"index":0,"delta":{"content":" response"}},{"index":1,"delta":{},"finish_reason":"tool_calls"},{"index":2,"delta":{},"finish_reason":"stop"}]}`)

	if toolCall, ok := acc.JustFinishedToolCall(); ok {
		t.Errorf("JustFinishedToolCall returned alternate choice tool call: %+v", toolCall)
	}
	if refusal, ok := acc.JustFinishedRefusal(); ok {
		t.Errorf("JustFinishedRefusal returned alternate choice refusal: %q", refusal)
	}
	if content, ok := acc.JustFinishedContent(); ok {
		t.Errorf("JustFinishedContent returned unfinished first choice content: %q", content)
	}

	finish := `{"id":"test","choices":[{"index":0,"delta":{},"finish_reason":"stop"},{"index":1,"delta":{},"finish_reason":"tool_calls"},{"index":2,"delta":{},"finish_reason":"stop"}]}`
	pr688AddChunk(t, &acc, finish)

	content, ok := acc.JustFinishedContent()
	if !ok {
		t.Fatal("JustFinishedContent did not return completed first choice content")
	}
	if content != "selected response" {
		t.Fatalf("content: expected %q, got %q", "selected response", content)
	}
	if toolCall, ok := acc.JustFinishedToolCall(); ok {
		t.Errorf("JustFinishedToolCall repeated alternate choice tool call: %+v", toolCall)
	}
	if refusal, ok := acc.JustFinishedRefusal(); ok {
		t.Errorf("JustFinishedRefusal repeated alternate choice refusal: %q", refusal)
	}

	pr688AddChunk(t, &acc, finish)
	if content, ok := acc.JustFinishedContent(); ok {
		t.Errorf("JustFinishedContent repeated completed first choice content: %q", content)
	}
}

func TestAccumulatorSingularToolCallIgnoresAlternateChoiceWhenFirstChoiceIsOmitted(t *testing.T) {
	var acc openai.ChatCompletionAccumulator

	pr688AddChunk(t, &acc, `{"id":"test","choices":[{"index":0,"delta":{"content":"selected"}},{"index":1,"delta":{"tool_calls":[{"index":0,"id":"alternate_call","type":"function","function":{"name":"lookup","arguments":"{\"id\":\"example\"}"}}]}}]}`)
	pr688AssertNoFinishedToolCall(t, &acc)

	finishAlternate := `{"id":"test","choices":[{"index":1,"delta":{},"finish_reason":"tool_calls"}]}`
	pr688AddChunk(t, &acc, finishAlternate)
	pr688AssertNoFinishedToolCall(t, &acc)

	toolCall, ok := acc.JustFinishedToolCallForChoice(1)
	if !ok {
		t.Fatal("JustFinishedToolCallForChoice did not return the alternate choice tool call")
	}
	if toolCall.ID != "alternate_call" || toolCall.Name != "lookup" {
		t.Fatalf("alternate choice tool call: got ID %q and name %q", toolCall.ID, toolCall.Name)
	}

	pr688AddChunk(t, &acc, finishAlternate)
	pr688AssertNoFinishedToolCall(t, &acc)
	if toolCall, ok := acc.JustFinishedToolCallForChoice(1); ok {
		t.Fatalf("JustFinishedToolCallForChoice returned repeated tool call: %+v", toolCall)
	}
}

func TestAccumulatorSingularToolCallKeepsIndependentLegacyChoiceState(t *testing.T) {
	var acc openai.ChatCompletionAccumulator

	pr688AddChunk(t, &acc, `{"id":"test","choices":[{"index":0,"delta":{"tool_calls":[{"index":0,"id":"selected_call","type":"function","function":{"name":"selected_tool","arguments":"{}"}}]}},{"index":1,"delta":{"content":"other"}}]}`)
	pr688AssertNoFinishedToolCall(t, &acc)

	pr688AddChunk(t, &acc, `{"id":"test","choices":[{"index":1,"delta":{"content":" continues"}},{"index":0,"delta":{"tool_calls":null}}]}`)
	pr688AssertNoFinishedToolCall(t, &acc)
	choiceToolCall, ok := acc.JustFinishedToolCallForChoice(0)
	if !ok {
		t.Fatal("JustFinishedToolCallForChoice did not return the selected choice tool call")
	}
	if choiceToolCall.ID != "selected_call" || choiceToolCall.Name != "selected_tool" {
		t.Fatalf("choice-aware tool call: got ID %q and name %q", choiceToolCall.ID, choiceToolCall.Name)
	}

	finishSelected := `{"id":"test","choices":[{"index":0,"delta":{},"finish_reason":"tool_calls"}]}`
	pr688AddChunk(t, &acc, finishSelected)
	legacyToolCall, ok := acc.JustFinishedToolCall()
	if !ok {
		t.Fatal("JustFinishedToolCall did not return the selected choice tool call")
	}
	if legacyToolCall.ID != "selected_call" || legacyToolCall.Name != "selected_tool" {
		t.Fatalf("legacy tool call: got ID %q and name %q", legacyToolCall.ID, legacyToolCall.Name)
	}
	if toolCall, ok := acc.JustFinishedToolCallForChoice(0); ok {
		t.Fatalf("JustFinishedToolCallForChoice returned repeated tool call: %+v", toolCall)
	}

	pr688AddChunk(t, &acc, finishSelected)
	pr688AssertNoFinishedToolCall(t, &acc)
}

func TestAccumulatorSingularContentUsesFirstWireChoice(t *testing.T) {
	var acc openai.ChatCompletionAccumulator

	pr688AddChunk(t, &acc, `{"id":"test","choices":[{"index":2,"delta":{"content":"choice two"}},{"index":1,"delta":{"content":"choice one"}}]}`)
	pr688AddChunk(t, &acc, `{"id":"test","choices":[{"index":2,"delta":{},"finish_reason":"stop"},{"index":1,"delta":{},"finish_reason":"stop"}]}`)

	content, ok := acc.JustFinishedContent()
	if !ok {
		t.Fatal("JustFinishedContent did not return completed content")
	}
	if content != "choice two" {
		t.Fatalf("content: expected %q, got %q", "choice two", content)
	}

	content, ok = acc.JustFinishedContentForChoice(1)
	if !ok {
		t.Fatal("JustFinishedContentForChoice did not return second choice content")
	}
	if content != "choice one" {
		t.Fatalf("choice 1 content: expected %q, got %q", "choice one", content)
	}
}

func pr688AddChunk(t *testing.T, acc *openai.ChatCompletionAccumulator, raw string) {
	t.Helper()

	var chunk openai.ChatCompletionChunk
	if err := chunk.UnmarshalJSON([]byte(raw)); err != nil {
		t.Fatalf("failed to unmarshal chunk: %v", err)
	}
	if !acc.AddChunk(chunk) {
		t.Fatal("AddChunk returned false")
	}
}

func pr688AssertNoFinishedToolCall(t *testing.T, acc *openai.ChatCompletionAccumulator) {
	t.Helper()

	if toolCall, ok := acc.JustFinishedToolCall(); ok {
		t.Fatalf("JustFinishedToolCall returned unexpected tool call: %+v", toolCall)
	}
}

func pr688AssertFinishedToolCall(t *testing.T, acc *openai.ChatCompletionAccumulator) {
	t.Helper()

	toolCall, ok := acc.JustFinishedToolCall()
	if !ok {
		t.Fatal("JustFinishedToolCall did not return the completed tool call")
	}
	if toolCall.Index != 0 {
		t.Fatalf("tool call index: expected 0, got %d", toolCall.Index)
	}
	if toolCall.ID != "call_123" {
		t.Fatalf("tool call ID: expected %q, got %q", "call_123", toolCall.ID)
	}
	if toolCall.Name != "test_func" {
		t.Fatalf("tool call name: expected %q, got %q", "test_func", toolCall.Name)
	}
	if toolCall.Arguments != "{}" {
		t.Fatalf("tool call arguments: expected %q, got %q", "{}", toolCall.Arguments)
	}
}

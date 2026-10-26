package openai_test

import (
	"encoding/json"
	"reflect"
	"testing"

	openai "github.com/openai/openai-go/v3"
)

func checkChatCompletionMessageConversion(t *testing.T, message openai.ChatCompletionMessage, want openai.ChatCompletionAssistantMessageParam, wantJSON string) {
	t.Helper()
	assistant := message.ToAssistantMessageParam()
	if !reflect.DeepEqual(assistant, want) {
		t.Errorf("assistant = %#v, want %#v", assistant, want)
	}
	union := message.ToParam()
	if !reflect.DeepEqual(union, openai.ChatCompletionMessageParamUnion{OfAssistant: &want}) {
		t.Errorf("message union = %#v, want only the assistant variant", union)
	}
	for _, value := range []any{assistant, union} {
		data, err := json.Marshal(value)
		if err != nil {
			t.Fatal(err)
		}
		var actual, expected any
		if err := json.Unmarshal(data, &actual); err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal([]byte(wantJSON), &expected); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("JSON = %s, want %s", data, wantJSON)
		}
	}
}

func TestChatCompletionMessageConversionEmptyValues(t *testing.T) {
	for _, test := range []struct {
		name, body string
	}{
		{"missing", `{}`},
		{"null", `{"content":null,"refusal":null,"audio":null,"function_call":null,"tool_calls":null}`},
		{"empty", `{"content":"","refusal":"","audio":{"id":""},"function_call":{"name":"","arguments":""},"tool_calls":[]}`},
	} {
		t.Run(test.name, func(t *testing.T) {
			var message openai.ChatCompletionMessage
			if err := json.Unmarshal([]byte(test.body), &message); err != nil {
				t.Fatal(err)
			}
			checkChatCompletionMessageConversion(t, message, openai.ChatCompletionAssistantMessageParam{}, `{"role":"assistant"}`)
		})
	}
	for _, toolCalls := range [][]openai.ChatCompletionMessageToolCallUnion{nil, {}} {
		message := openai.ChatCompletionMessage{ToolCalls: toolCalls}
		checkChatCompletionMessageConversion(t, message, openai.ChatCompletionAssistantMessageParam{}, `{"role":"assistant"}`)
	}
}

func TestChatCompletionMessageConversionUsesCurrentFields(t *testing.T) {
	const body = `{
		"role":"assistant","content":"original","refusal":"no",
		"audio":{"id":"audio_1","data":"synthetic","expires_at":123,"transcript":"not copied"},
		"function_call":{"name":"legacy","arguments":"{}"},
		"tool_calls":[
			{"id":"call_1","type":"function","function":{"name":"lookup","arguments":"{}"}},
			{"id":"call_2","type":"custom","custom":{"name":"format","input":"original"}}
		],
		"annotations":[{"type":"url_citation","url_citation":{"start_index":0,"end_index":1,"title":"fixture","url":"https://example.com"}}],
		"extension":"not copied"
	}`
	var decoded openai.ChatCompletionMessage
	if err := json.Unmarshal([]byte(body), &decoded); err != nil {
		t.Fatal(err)
	}
	decoded.Content = "edited"
	decoded.ToolCalls[0].Function.Arguments = `{"x":1}`
	decoded.ToolCalls[1].Custom.Input = "edited input"
	manual := openai.ChatCompletionMessage{
		Role:    "assistant",
		Content: "edited",
		Refusal: "no",
		Audio: openai.ChatCompletionAudio{
			ID: "audio_1", Data: "synthetic", ExpiresAt: 123, Transcript: "not copied",
		},
		FunctionCall: openai.ChatCompletionMessageFunctionCall{Name: "legacy", Arguments: "{}"},
		ToolCalls: []openai.ChatCompletionMessageToolCallUnion{
			{ID: "call_1", Type: "function", Function: openai.ChatCompletionMessageFunctionToolCallFunction{Name: "lookup", Arguments: `{"x":1}`}},
			{ID: "call_2", Type: "custom", Custom: openai.ChatCompletionMessageCustomToolCallCustom{Name: "format", Input: "edited input"}},
		},
	}
	want := openai.ChatCompletionAssistantMessageParam{
		Role:    "assistant",
		Content: openai.ChatCompletionAssistantMessageParamContentUnion{OfString: openai.String("edited")},
		Refusal: openai.String("no"),
		Audio:   openai.ChatCompletionAssistantMessageParamAudio{ID: "audio_1"},
		FunctionCall: openai.ChatCompletionAssistantMessageParamFunctionCall{
			Name: "legacy", Arguments: "{}",
		},
		ToolCalls: []openai.ChatCompletionMessageToolCallUnionParam{
			{OfFunction: &openai.ChatCompletionMessageFunctionToolCallParam{
				ID: "call_1", Function: openai.ChatCompletionMessageFunctionToolCallFunctionParam{Name: "lookup", Arguments: `{"x":1}`},
			}},
			{OfCustom: &openai.ChatCompletionMessageCustomToolCallParam{
				ID: "call_2", Custom: openai.ChatCompletionMessageCustomToolCallCustomParam{Name: "format", Input: "edited input"},
			}},
		},
	}
	const wantJSON = `{
		"role":"assistant","content":"edited","refusal":"no","audio":{"id":"audio_1"},
		"function_call":{"name":"legacy","arguments":"{}"},
		"tool_calls":[
			{"id":"call_1","type":"function","function":{"name":"lookup","arguments":"{\"x\":1}"}},
			{"id":"call_2","type":"custom","custom":{"name":"format","input":"edited input"}}
		]
	}`
	t.Run("decoded", func(t *testing.T) {
		if !decoded.JSON.Content.Valid() || decoded.JSON.ExtraFields["extension"].Raw() != `"not copied"` {
			t.Fatal("fixture did not retain response metadata")
		}
		checkChatCompletionMessageConversion(t, decoded, want, wantJSON)
	})
	t.Run("manual", func(t *testing.T) {
		if manual.JSON.Content.Valid() || manual.ToolCalls[0].RawJSON() != "" {
			t.Fatal("manual fixture unexpectedly has response metadata")
		}
		checkChatCompletionMessageConversion(t, manual, want, wantJSON)
	})
}

func TestChatCompletionMessageConversionFromAccumulator(t *testing.T) {
	var accumulator openai.ChatCompletionAccumulator
	for _, body := range []string{
		`{"id":"chatcmpl_fixture","choices":[{"index":0,"delta":{"role":"assistant","content":"hel","refusal":"n","tool_calls":[{"index":0,"id":"call_1","type":"function","function":{"name":"lookup","arguments":"{\"x\":"}}]}}]}`,
		`{"id":"chatcmpl_fixture","choices":[{"index":0,"delta":{"content":"lo","refusal":"o","tool_calls":[{"index":0,"function":{"arguments":"1}"}}]}}]}`,
	} {
		var chunk openai.ChatCompletionChunk
		if err := json.Unmarshal([]byte(body), &chunk); err != nil {
			t.Fatal(err)
		}
		if !accumulator.AddChunk(chunk) {
			t.Fatal("chunk was not accumulated")
		}
	}
	if len(accumulator.Choices) != 1 {
		t.Fatalf("choices = %d, want 1", len(accumulator.Choices))
	}
	message := accumulator.Choices[0].Message
	if message.JSON.Content.Valid() || message.RawJSON() != "" {
		t.Fatal("accumulator unexpectedly populated message JSON metadata")
	}
	want := openai.ChatCompletionAssistantMessageParam{
		Role:    "assistant",
		Content: openai.ChatCompletionAssistantMessageParamContentUnion{OfString: openai.String("hello")},
		Refusal: openai.String("no"),
		ToolCalls: []openai.ChatCompletionMessageToolCallUnionParam{{
			OfFunction: &openai.ChatCompletionMessageFunctionToolCallParam{
				ID: "call_1", Function: openai.ChatCompletionMessageFunctionToolCallFunctionParam{Name: "lookup", Arguments: `{"x":1}`},
			},
		}},
	}
	checkChatCompletionMessageConversion(t, message, want,
		`{"role":"assistant","content":"hello","refusal":"no","tool_calls":[{"id":"call_1","type":"function","function":{"name":"lookup","arguments":"{\"x\":1}"}}]}`)
}

func TestChatCompletionMessageConversionUnknownToolVariants(t *testing.T) {
	message := openai.ChatCompletionMessage{
		ToolCalls: []openai.ChatCompletionMessageToolCallUnion{
			{ID: "unknown", Type: "future"},
			{ID: "known", Type: "function", Function: openai.ChatCompletionMessageFunctionToolCallFunction{Name: "lookup", Arguments: "{}"}},
			{ID: "missing type"},
		},
	}
	want := openai.ChatCompletionAssistantMessageParam{
		ToolCalls: []openai.ChatCompletionMessageToolCallUnionParam{
			{},
			{OfFunction: &openai.ChatCompletionMessageFunctionToolCallParam{
				ID: "known", Function: openai.ChatCompletionMessageFunctionToolCallFunctionParam{Name: "lookup", Arguments: "{}"},
			}},
			{},
		},
	}
	checkChatCompletionMessageConversion(t, message, want,
		`{"role":"assistant","tool_calls":[null,{"id":"known","type":"function","function":{"name":"lookup","arguments":"{}"}},null]}`)
}

func TestChatCompletionMessageConversionCopiesToolCalls(t *testing.T) {
	message := openai.ChatCompletionMessage{
		ToolCalls: []openai.ChatCompletionMessageToolCallUnion{
			{ID: "call_1", Type: "function", Function: openai.ChatCompletionMessageFunctionToolCallFunction{Name: "lookup", Arguments: "{}"}},
			{ID: "call_2", Type: "custom", Custom: openai.ChatCompletionMessageCustomToolCallCustom{Name: "format", Input: "original"}},
		},
	}
	first := message.ToAssistantMessageParam()
	second := message.ToAssistantMessageParam()
	if len(first.ToolCalls) != 2 || len(second.ToolCalls) != 2 ||
		first.ToolCalls[0].OfFunction == nil || first.ToolCalls[1].OfCustom == nil ||
		second.ToolCalls[0].OfFunction == nil || second.ToolCalls[1].OfCustom == nil {
		t.Fatal("conversion did not populate both tool variants")
	}
	first.ToolCalls[0].OfFunction.Function.Name = "changed output"
	first.ToolCalls[1].OfCustom.Custom.Input = "changed output"
	if message.ToolCalls[0].Function.Name != "lookup" || message.ToolCalls[1].Custom.Input != "original" ||
		second.ToolCalls[0].OfFunction.Function.Name != "lookup" || second.ToolCalls[1].OfCustom.Custom.Input != "original" {
		t.Fatal("conversions share mutable tool-call state")
	}
	message.ToolCalls[0].ID = "changed source"
	message.ToolCalls[1].Custom.Name = "changed source"
	if first.ToolCalls[0].OfFunction.ID != "call_1" || first.ToolCalls[1].OfCustom.Custom.Name != "format" {
		t.Fatal("conversion retained mutable source tool-call state")
	}
}

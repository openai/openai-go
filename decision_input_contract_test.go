package openai_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/live"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/shared"
)

func TestDecisionExpandedInputWireContract(t *testing.T) {
	requests := make(chan map[string]any, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Error(err)
		}
		requests <- body
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	}))
	t.Cleanup(server.Close)

	params := openai.DecisionNewParams{
		Model: "test-model",
		Input: openai.DecisionNewParamsInputUnion{OfDecisionNewsInputArray: []openai.DecisionNewParamsInputArrayItemUnion{
			{OfDecisionInputUserMessage: &openai.DecisionNewParamsInputArrayItemDecisionInputUserMessage{
				Content: openai.DecisionNewParamsInputArrayItemDecisionInputUserMessageContentUnion{OfString: openai.String("review")},
			}},
			{OfDecisionInputDeveloperMessage: &openai.DecisionNewParamsInputArrayItemDecisionInputDeveloperMessage{
				Content: openai.DecisionNewParamsInputArrayItemDecisionInputDeveloperMessageContentUnion{
					OfParts: []openai.DecisionNewParamsInputArrayItemDecisionInputDeveloperMessageContentParts{{Text: "developer guidance"}},
				},
			}},
			{OfDecisionInputSystemMessage: &openai.DecisionNewParamsInputArrayItemDecisionInputSystemMessage{
				Content: openai.DecisionNewParamsInputArrayItemDecisionInputSystemMessageContentUnion{
					OfParts: []openai.DecisionNewParamsInputArrayItemDecisionInputSystemMessageContentParts{{Text: "system guidance"}},
				},
			}},
			{OfFunctionCall: &openai.DecisionNewParamsInputArrayItemFunctionCall{CallID: "call-1", Name: "lookup", Arguments: "{}"}},
			{OfFunctionCallOutput: &openai.DecisionNewParamsInputArrayItemFunctionCallOutput{
				CallID: "call-1", Output: openai.DecisionNewParamsInputArrayItemFunctionCallOutputOutputUnion{OfString: openai.String("found")},
			}},
			{OfCustomToolCall: &openai.DecisionNewParamsInputArrayItemCustomToolCall{CallID: "call-2", Name: "custom", Input: "payload"}},
			{OfCustomToolCallOutput: &openai.DecisionNewParamsInputArrayItemCustomToolCallOutput{
				CallID: "call-2", Output: openai.DecisionNewParamsInputArrayItemCustomToolCallOutputOutputUnion{OfString: openai.String("done")},
			}},
		}},
		Questions: []openai.DecisionNewParamsQuestionUnion{{OfPredicate: &openai.DecisionNewParamsQuestionPredicate{
			Instructions: openai.DecisionNewParamsQuestionPredicateInstructionsUnion{
				OfDecisionQuestionInputPartArray: []openai.DecisionQuestionInputPartUnionParam{
					shared.DecisionQuestionInputPartParamOfInputJSON(map[string]any{"evidence": map[string]any{"ok": true}}),
				},
			},
		}}},
	}
	client := openai.NewClient(option.WithAPIKey("fake-test-key"), option.WithBaseURL(server.URL), option.WithUnsafeAllowHTTP(), option.WithMaxRetries(0))
	if _, err := client.Decisions.New(context.Background(), params); err != nil {
		t.Fatal(err)
	}
	got := <-requests
	want := map[string]any{
		"model": "test-model",
		"input": []any{
			map[string]any{"role": "user", "content": "review"},
			map[string]any{"role": "developer", "content": []any{map[string]any{"type": "input_text", "text": "developer guidance"}}},
			map[string]any{"role": "system", "content": []any{map[string]any{"type": "input_text", "text": "system guidance"}}},
			map[string]any{"type": "function_call", "call_id": "call-1", "name": "lookup", "arguments": "{}"},
			map[string]any{"type": "function_call_output", "call_id": "call-1", "output": "found"},
			map[string]any{"type": "custom_tool_call", "call_id": "call-2", "name": "custom", "input": "payload"},
			map[string]any{"type": "custom_tool_call_output", "call_id": "call-2", "output": "done"},
		},
		"questions": []any{map[string]any{
			"type": "predicate", "instructions": []any{map[string]any{"type": "input_json", "data": map[string]any{"evidence": map[string]any{"ok": true}}}},
		}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("request body = %#v, want %#v", got, want)
	}
}

func TestLiveDecisionQuestionInputPartResponseVariant(t *testing.T) {
	var part live.DecisionQuestionChoiceInstructionsPartsItemUnion
	if err := json.Unmarshal([]byte(`{"type":"input_text","text":"review"}`), &part); err != nil {
		t.Fatal(err)
	}
	text := part.AsDecisionInputText()
	if !part.JSON.OfDecisionInputText.Valid() || text.Text != "review" {
		t.Fatalf("question input = %#v", part)
	}
}

func TestDecisionInputMessagesWithoutType(t *testing.T) {
	var items []openai.DecisionNewParamsInputArrayItemUnion
	data := []byte(`[{"role":"user","content":"hello"},{"role":"developer","content":"hello"},{"role":"system","content":"hello"},{"role":"assistant","content":"hello"}]`)
	if err := json.Unmarshal(data, &items); err != nil {
		t.Fatal(err)
	}
	if len(items) != 4 || items[0].OfDecisionInputUserMessage == nil || items[1].OfDecisionInputDeveloperMessage == nil || items[2].OfDecisionInputSystemMessage == nil || items[3].OfDecisionInputAssistantMessage == nil {
		t.Fatalf("message variants = %#v", items)
	}
}

func TestDecisionMessageHelpersAsRequestInput(t *testing.T) {
	input := openai.DecisionNewParamsInputUnion{OfDecisionNewsInputArray: []openai.DecisionNewParamsInputArrayItemUnion{
		{OfDecisionInputMessage: openai.Ptr(openai.DecisionInputMessageParamOfUser("user input"))},
		{OfDecisionInputMessage: openai.Ptr(openai.DecisionInputMessageParamOfDeveloper("developer input"))},
		{OfDecisionInputMessage: openai.Ptr(openai.DecisionInputMessageParamOfSystem("system input"))},
		{OfDecisionInputMessage: openai.Ptr(openai.DecisionInputMessageParamOfAssistant("assistant input"))},
	}}
	data, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}
	var got []map[string]any
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	want := []map[string]any{
		{"role": "user", "content": "user input"},
		{"role": "developer", "content": "developer input"},
		{"role": "system", "content": "system input"},
		{"role": "assistant", "content": "assistant input"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("input = %#v, want %#v", got, want)
	}
}

func TestLiveDecisionInputJSONResponseShapes(t *testing.T) {
	var object live.DecisionQuestionChoiceInstructionsPartInputJSON
	if err := json.Unmarshal([]byte(`{"type":"input_json","data":{"key":"value"}}`), &object); err != nil {
		t.Fatal(err)
	}
	if object.Data.OfLiveDecisionJSONObject["key"] != "value" || object.Data.OfArrayOfDecisionJSONObject != nil {
		t.Fatalf("object data = %#v", object.Data)
	}

	var array live.DecisionQuestionChoiceInstructionsPartInputJSON
	if err := json.Unmarshal([]byte(`{"type":"input_json","data":[{"key":"value"}]}`), &array); err != nil {
		t.Fatal(err)
	}
	if len(array.Data.OfArrayOfDecisionJSONObject) != 1 || array.Data.OfArrayOfDecisionJSONObject[0]["key"] != "value" || array.Data.OfLiveDecisionJSONObject != nil {
		t.Fatalf("array data = %#v", array.Data)
	}
}

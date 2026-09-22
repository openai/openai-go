package responses_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
)

func TestResponseOutputTextAggregation(t *testing.T) {
	for _, test := range []struct {
		name   string
		output []responses.ResponseOutputItemUnion
		want   string
	}{
		{"nil output", nil, ""},
		{"empty output", []responses.ResponseOutputItemUnion{}, ""},
		{"nil and empty content", []responses.ResponseOutputItemUnion{{}, {Content: []responses.ResponseOutputMessageContentUnion{}}}, ""},
		{
			"only output text in item and content order",
			[]responses.ResponseOutputItemUnion{
				{Type: "message", Content: []responses.ResponseOutputMessageContentUnion{
					{Type: "output_text", Text: "α\n"},
					{Type: "refusal", Refusal: "no", Text: "not copied"},
					{Type: "output_text", Text: ""},
					{Type: "future", Text: "not copied"},
					{Text: "not copied"},
					{Type: "output_text", Text: "β"},
				}},
				{},
				{Type: "message", Content: []responses.ResponseOutputMessageContentUnion{{Type: "output_text", Text: "!"}}},
			},
			"α\nβ!",
		},
		{
			"content on manually constructed item without message discriminator",
			[]responses.ResponseOutputItemUnion{{Content: []responses.ResponseOutputMessageContentUnion{{Type: "output_text", Text: "kept"}}}},
			"kept",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			response := responses.Response{Output: test.output}
			if got := response.OutputText(); got != test.want {
				t.Fatalf("OutputText() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestResponseOutputTextUsesCurrentFields(t *testing.T) {
	const body = `{"output":[{"type":"message","content":[{"type":"output_text","text":"original"},{"type":"output_text","text":"removed"}]}]}`
	var response responses.Response
	if err := json.Unmarshal([]byte(body), &response); err != nil {
		t.Fatal(err)
	}
	if len(response.Output) != 1 || len(response.Output[0].Content) != 2 || response.RawJSON() != body {
		t.Fatal("fixture did not decode with response metadata")
	}
	response.Output[0].Content[0].Text = "edited"
	response.Output[0].Content[1].Type = "refusal"
	wantContent := append([]responses.ResponseOutputMessageContentUnion(nil), response.Output[0].Content...)
	for range 2 {
		if got := response.OutputText(); got != "edited" {
			t.Fatalf("OutputText() = %q, want current text field", got)
		}
	}
	if !reflect.DeepEqual(response.Output[0].Content, wantContent) || response.RawJSON() != body {
		t.Fatal("OutputText mutated the input or response metadata")
	}
}

func TestResponseToInputPreservesCompleteOutputForFollowupRequest(t *testing.T) {
	const body = `{"output":[{"id":"rs_123","type":"reasoning","summary":[{"type":"summary_text","text":"reasoned"}],"encrypted_content":"encrypted-reasoning","status":"completed","future_reasoning_field":{"version":1}},{"id":"msg_123","type":"message","role":"assistant","phase":"commentary","content":[{"type":"output_text","text":"I'll check that."}],"status":"completed","future_message_field":true},{"id":"fc_123","type":"function_call","call_id":"call_123","name":"lookup","arguments":"{\\"city\\":\\"Paris\\"}","status":"completed","future_tool_field":{"version":2}}]}`

	var response responses.Response
	if err := json.Unmarshal([]byte(body), &response); err != nil {
		t.Fatal(err)
	}

	input := response.ToInput()
	if len(input) != 3 {
		t.Fatalf("ToInput() returned %d items, want 3", len(input))
	}
	for i, item := range input {
		if _, ok := item.Overrides(); !ok {
			t.Fatalf("ToInput()[%d] was not preserved as a raw API item", i)
		}
	}

	toolOutput := responses.ResponseInputItemParamOfFunctionCallOutput(`{"temperature_c":20}`)
	toolOutput.OfFunctionCallOutput.CallID = openai.String("call_123")
	input = append(input, toolOutput)
	input = append(input, responses.ResponseInputItemParamOfMessage(
		"Use the tool result.",
		responses.EasyInputMessageRoleUser,
	))

	data, err := json.Marshal(responses.ResponseNewParams{
		Model: openai.ChatModelGPT5_2,
		Store: openai.Bool(false),
		Input: responses.ResponseNewParamsInputUnion{OfInputItemList: input},
	})
	if err != nil {
		t.Fatal(err)
	}
	var request struct {
		Input []map[string]json.RawMessage `json:"input"`
	}
	if err := json.Unmarshal(data, &request); err != nil {
		t.Fatal(err)
	}

	items := request.Input
	if len(items) != 5 ||
		string(items[0]["type"]) != `"reasoning"` ||
		string(items[1]["type"]) != `"message"` ||
		string(items[1]["phase"]) != `"commentary"` ||
		string(items[2]["type"]) != `"function_call"` ||
		string(items[2]["call_id"]) != `"call_123"` ||
		string(items[3]["type"]) != `"function_call_output"` ||
		string(items[3]["call_id"]) != `"call_123"` ||
		string(items[4]["role"]) != `"user"` ||
		string(items[4]["content"]) != `"Use the tool result."` {
		t.Fatalf("ToInput() changed complete replay ordering: %s", data)
	}
	if string(items[0]["encrypted_content"]) != `"encrypted-reasoning"` {
		t.Fatalf("ToInput() dropped encrypted reasoning content: %s", data)
	}
	if string(items[0]["future_reasoning_field"]) != `{"version":1}` ||
		string(items[1]["future_message_field"]) != "true" ||
		string(items[2]["future_tool_field"]) != `{"version":2}` {
		t.Fatalf("ToInput() dropped output fields: %s", data)
	}
}

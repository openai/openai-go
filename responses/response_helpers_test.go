package responses_test

import (
	"encoding/json"
	"reflect"
	"testing"

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

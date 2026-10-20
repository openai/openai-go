package openai_test

import (
	"encoding/json"
	"reflect"
	"testing"

	openai "github.com/openai/openai-go/v3"
)

func checkAssistantMessageConstructor[T string | []openai.ChatCompletionAssistantMessageParamContentArrayOfContentPartUnion](
	t *testing.T,
	content T,
	wantContent openai.ChatCompletionAssistantMessageParamContentUnion,
	wantJSON string,
) {
	t.Helper()
	got := openai.AssistantMessage(content)
	generated := openai.ChatCompletionMessageParamOfAssistant(content)
	if !reflect.DeepEqual(got, generated) {
		t.Errorf("compatibility constructor = %#v, generated constructor = %#v", got, generated)
	}
	if got.OfAssistant == nil {
		t.Fatal("assistant variant is nil")
	}
	if !reflect.DeepEqual(got.OfAssistant.Content, wantContent) {
		t.Errorf("content = %#v, want %#v", got.OfAssistant.Content, wantContent)
	}
	if got.OfDeveloper != nil || got.OfSystem != nil || got.OfUser != nil || got.OfTool != nil || got.OfFunction != nil {
		t.Error("constructor populated another message variant")
	}
	data, err := json.Marshal(got)
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

func TestAssistantMessageConstructorStrings(t *testing.T) {
	for _, test := range []struct {
		name, content, wantJSON string
	}{
		{"text", "hello", `{"content":"hello","role":"assistant"}`},
		{"empty text", "", `{"content":"","role":"assistant"}`},
	} {
		t.Run(test.name, func(t *testing.T) {
			checkAssistantMessageConstructor(t, test.content,
				openai.ChatCompletionAssistantMessageParamContentUnion{
					OfString: openai.String(test.content),
				}, test.wantJSON)
		})
	}
}

func TestAssistantMessageConstructorContentParts(t *testing.T) {
	for _, test := range []struct {
		name     string
		content  []openai.ChatCompletionAssistantMessageParamContentArrayOfContentPartUnion
		wantJSON string
	}{
		{"nil", nil, `{"role":"assistant"}`},
		{"empty", []openai.ChatCompletionAssistantMessageParamContentArrayOfContentPartUnion{}, `{"content":[],"role":"assistant"}`},
		{
			"text and refusal",
			[]openai.ChatCompletionAssistantMessageParamContentArrayOfContentPartUnion{
				{OfText: &openai.ChatCompletionContentPartTextParam{Text: "hello"}},
				{OfRefusal: &openai.ChatCompletionContentPartRefusalParam{Refusal: "no"}},
			},
			`{"content":[{"text":"hello","type":"text"},{"refusal":"no","type":"refusal"}],"role":"assistant"}`,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			checkAssistantMessageConstructor(t, test.content,
				openai.ChatCompletionAssistantMessageParamContentUnion{
					OfArrayOfContentParts: test.content,
				}, test.wantJSON)
		})
	}
}

func TestAssistantMessageConstructorRetainsContentSlice(t *testing.T) {
	content := []openai.ChatCompletionAssistantMessageParamContentArrayOfContentPartUnion{
		{OfText: &openai.ChatCompletionContentPartTextParam{Text: "before"}},
	}
	got := openai.AssistantMessage(content)
	refusal := &openai.ChatCompletionContentPartRefusalParam{Refusal: "after"}
	content[0] = openai.ChatCompletionAssistantMessageParamContentArrayOfContentPartUnion{OfRefusal: refusal}
	if got.OfAssistant == nil || len(got.OfAssistant.Content.OfArrayOfContentParts) != 1 ||
		got.OfAssistant.Content.OfArrayOfContentParts[0].OfRefusal != refusal {
		t.Fatalf("constructor did not retain the caller's content slice: %#v", got)
	}
}

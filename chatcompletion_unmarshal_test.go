package openai_test

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/openai/openai-go/v3"
)

func TestChatCompletionChunkUnmarshalValidatesTopLevelShape(t *testing.T) {
	t.Run("object", func(t *testing.T) {
		var chunk openai.ChatCompletionChunk
		if err := json.Unmarshal([]byte(`{"id":"chatcmpl_test","choices":[]}`), &chunk); err != nil {
			t.Fatalf("unmarshal object: %v", err)
		}
		if chunk.ID != "chatcmpl_test" {
			t.Fatalf("chunk ID = %q, want %q", chunk.ID, "chatcmpl_test")
		}
	})

	t.Run("null", func(t *testing.T) {
		chunk := openai.ChatCompletionChunk{ID: "unchanged"}
		if err := json.Unmarshal([]byte(`null`), &chunk); err != nil {
			t.Fatalf("unmarshal null: %v", err)
		}
		if chunk.ID != "unchanged" {
			t.Fatalf("chunk ID = %q after null, want unchanged", chunk.ID)
		}
	})

	for name, test := range map[string]struct {
		input string
		value string
	}{
		"array":   {input: `[]`, value: "array"},
		"boolean": {input: `true`, value: "bool"},
		"number":  {input: `42`, value: "number"},
		"string":  {input: `"chunk"`, value: "string"},
	} {
		t.Run(name, func(t *testing.T) {
			chunk := openai.ChatCompletionChunk{ID: "unchanged"}
			err := json.Unmarshal([]byte(test.input), &chunk)
			if err == nil {
				t.Fatalf("unmarshal %s succeeded, want a type error", name)
			}
			var typeErr *json.UnmarshalTypeError
			if !errors.As(err, &typeErr) {
				t.Fatalf("error type = %T, want *json.UnmarshalTypeError", err)
			}
			if typeErr.Value != test.value {
				t.Fatalf("error value = %q, want %q", typeErr.Value, test.value)
			}
			if typeErr.Type != reflect.TypeOf(chunk) {
				t.Fatalf("error target type = %v, want %v", typeErr.Type, reflect.TypeOf(chunk))
			}
			if chunk.ID != "unchanged" {
				t.Fatalf("chunk ID = %q after failed unmarshal, want unchanged", chunk.ID)
			}
			if got := chunk.RawJSON(); got != "" {
				t.Fatalf("raw JSON = %q after failed unmarshal, want empty", got)
			}
		})
	}
}

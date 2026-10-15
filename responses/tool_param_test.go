package responses_test

import (
	"encoding/json"
	"testing"

	"github.com/openai/openai-go/v3/responses"
)

func TestToolParamOfFunctionWithDescription(t *testing.T) {
	tool := responses.ToolParamOfFunctionWithDescription(
		"search_docs",
		"Search the documentation",
		map[string]any{"type": "object"},
		true,
	)

	data, err := json.Marshal(tool)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	var body map[string]any
	if err := json.Unmarshal(data, &body); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if body["name"] != "search_docs" {
		t.Fatalf("name = %v, want %q", body["name"], "search_docs")
	}
	if body["description"] != "Search the documentation" {
		t.Fatalf("description = %v, want %q", body["description"], "Search the documentation")
	}
	if body["strict"] != true {
		t.Fatalf("strict = %v, want true", body["strict"])
	}
}

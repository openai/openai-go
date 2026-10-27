package openai_test

import (
	"encoding/json"
	"testing"

	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
)

func TestReleasedToolSearchOutputSlices(t *testing.T) {
	stable := []responses.ResponseToolSearchOutputItemParamToolUnion{
		{OfMcp: &responses.ResponseToolSearchOutputItemParamToolMcpParam{ServerLabel: "docs"}},
	}
	beta := []openai.BetaResponseToolSearchOutputItemParamToolUnion{
		{OfMcp: &openai.BetaResponseToolSearchOutputItemParamToolMcpParam{ServerLabel: "docs"}},
	}

	for _, tc := range []struct {
		name  string
		value any
	}{
		{"stable direct", responses.ResponseToolSearchOutputItemParam{Tools: stable}},
		{"stable constructor", responses.ResponseInputItemParamOfToolSearchOutput(stable)},
		{"beta direct", openai.BetaResponseToolSearchOutputItemParam{Tools: beta}},
		{"beta constructor", openai.BetaResponseInputItemParamOfToolSearchOutput(beta)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			data, err := json.Marshal(tc.value)
			if err != nil {
				t.Fatal(err)
			}
			var body struct {
				Type  string `json:"type"`
				Tools []struct {
					Type        string `json:"type"`
					ServerLabel string `json:"server_label"`
				} `json:"tools"`
			}
			if err := json.Unmarshal(data, &body); err != nil {
				t.Fatal(err)
			}
			if body.Type != "tool_search_output" || len(body.Tools) != 1 ||
				body.Tools[0].Type != "mcp" || body.Tools[0].ServerLabel != "docs" {
				t.Fatalf("released tool-search output did not reach the wire: %s", data)
			}
		})
	}
}

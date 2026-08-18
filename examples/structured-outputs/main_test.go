package main

import (
	"encoding/json"
	"testing"
)

func TestHistoricalComputerParamsSerializesStructuredOutputSchema(t *testing.T) {
	schema, err := GenerateSchema[HistoricalComputer]()
	if err != nil {
		t.Fatalf("GenerateSchema() error = %v", err)
	}

	data, err := json.Marshal(historicalComputerParams("question", schema))
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	var request struct {
		Text struct {
			Format struct {
				Name        string         `json:"name"`
				Description string         `json:"description"`
				Schema      map[string]any `json:"schema"`
				Strict      bool           `json:"strict"`
				Type        string         `json:"type"`
			} `json:"format"`
		} `json:"text"`
	}
	if err := json.Unmarshal(data, &request); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	format := request.Text.Format
	if format.Type != "json_schema" {
		t.Errorf("format.type = %q, want %q", format.Type, "json_schema")
	}
	if format.Name != "historical_computer" {
		t.Errorf("format.name = %q, want %q", format.Name, "historical_computer")
	}
	if format.Description == "" {
		t.Error("format.description is empty")
	}
	if !format.Strict {
		t.Error("format.strict = false, want true")
	}
	if format.Schema["type"] != "object" {
		t.Errorf("format.schema.type = %v, want %q", format.Schema["type"], "object")
	}
	if format.Schema["additionalProperties"] != false {
		t.Errorf("format.schema.additionalProperties = %v, want false", format.Schema["additionalProperties"])
	}
}

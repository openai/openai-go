package main

import (
	"encoding/json"
	"strings"
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

func TestGenerateSchemaPreserves64BitIntegerEnums(t *testing.T) {
	type exactIntegerEnums struct {
		Signed   int64  `json:"signed" jsonschema:"enum=-9007199254740993,enum=9007199254740993"`
		Unsigned uint64 `json:"unsigned" jsonschema:"enum=9007199254740993,enum=18446744073709551615"`
	}

	schema, err := GenerateSchema[exactIntegerEnums]()
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
				Schema struct {
					Properties map[string]struct {
						Enum []json.RawMessage `json:"enum"`
					} `json:"properties"`
				} `json:"schema"`
			} `json:"format"`
		} `json:"text"`
	}
	if err := json.Unmarshal(data, &request); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	for _, test := range []struct {
		property string
		want     []string
	}{
		{property: "signed", want: []string{"-9007199254740993", "9007199254740993"}},
		{property: "unsigned", want: []string{"9007199254740993", "18446744073709551615"}},
	} {
		enum := request.Text.Format.Schema.Properties[test.property].Enum
		if len(enum) != len(test.want) {
			t.Fatalf("%s enum has %d values, want %d: %s", test.property, len(enum), len(test.want), data)
		}
		for i, value := range enum {
			if string(value) != test.want[i] {
				t.Errorf("%s enum[%d] = %s, want exact numeric token %s: %s", test.property, i, value, test.want[i], data)
			}
		}
	}
}

func TestGenerateSchemaPreservesRecursiveReferences(t *testing.T) {
	type recursiveSchemaNode struct {
		Children []recursiveSchemaNode `json:"children"`
	}

	schema, err := GenerateSchema[recursiveSchemaNode]()
	if err != nil {
		t.Fatalf("GenerateSchema() error = %v", err)
	}

	data, err := json.Marshal(schema)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	var decoded struct {
		Type string `json:"type"`
		Defs map[string]struct {
			Properties map[string]struct {
				Items struct {
					Ref string `json:"$ref"`
				} `json:"items"`
			} `json:"properties"`
		} `json:"$defs"`
	}
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if decoded.Type != "object" {
		t.Fatalf("root type = %q, want %q", decoded.Type, "object")
	}
	if len(decoded.Defs) == 0 {
		t.Fatal("recursive schema has no $defs")
	}

	const refPrefix = "#/$defs/"
	for definitionName, definition := range decoded.Defs {
		children, ok := definition.Properties["children"]
		if !ok {
			continue
		}
		if !strings.HasPrefix(children.Items.Ref, refPrefix) {
			t.Fatalf("%s children items $ref = %q, want prefix %q", definitionName, children.Items.Ref, refPrefix)
		}
		if strings.TrimPrefix(children.Items.Ref, refPrefix) != definitionName {
			t.Fatalf("%s children items $ref = %q, want self-reference", definitionName, children.Items.Ref)
		}
		return
	}

	t.Fatalf("recursive $defs has no definition with children property: %s", data)
}

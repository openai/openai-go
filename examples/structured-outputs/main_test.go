package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
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

func TestHistoricalComputerFromResponse(t *testing.T) {
	const valid = `{"origin":{"year_of_construction":1951,"organization":"Harvard"},"full_name":"SNARC","legacy":"positive","notable_facts":["Neural network computer"]}`
	textPart := func(text string) string {
		return fmt.Sprintf(`{"type":"output_text","text":%q,"annotations":[]}`, text)
	}
	message := func(phase, content string) string {
		return `{"type":"message","id":"msg_fixture","role":"assistant","status":"completed"` + phase + `,"content":[` + content + `]}`
	}
	final := message(`,"phase":"final_answer"`, textPart(valid))
	refusal := `{"type":"refusal","refusal":"Synthetic refusal"}`
	for _, test := range []struct {
		name    string
		status  string
		output  string
		wantErr string
	}{
		{name: "final answer", output: final},
		{name: "missing phase", output: message("", textPart(valid))},
		{name: "null phase", output: message(`,"phase":null`, textPart(valid))},
		{name: "commentary before final", output: message(`,"phase":"commentary"`, textPart("Let me check.")) + "," + final},
		{name: "reasoning before final", output: `{"type":"reasoning","id":"rs_fixture","summary":[]},` + final},
		{name: "incomplete response with valid JSON", status: "incomplete", output: final, wantErr: "response did not complete: incomplete"},
		{name: "failed response", status: "failed", output: final, wantErr: "response did not complete: failed"},
		{name: "cancelled response", status: "cancelled", output: final, wantErr: "response did not complete: cancelled"},
		{name: "refusal", output: message("", refusal), wantErr: "model refused"},
		{name: "refusal after text", output: message("", textPart(valid)+","+refusal), wantErr: "model refused"},
		{name: "empty output", wantErr: "no final JSON text"},
		{name: "commentary only", output: message(`,"phase":"commentary"`, textPart(valid)), wantErr: "no final JSON text"},
		{name: "unknown phase", output: message(`,"phase":"future_phase"`, textPart(valid)), wantErr: "no final JSON text"},
		{name: "malformed JSON", output: message("", textPart("{")), wantErr: "decode historical computer"},
		{name: "wrong field type", output: message("", textPart(`{"full_name":42}`)), wantErr: "decode historical computer"},
		{name: "null JSON", output: message("", textPart("null")), wantErr: "received null"},
		{name: "multiple final messages", output: final + "," + final, wantErr: "expected one final JSON text part"},
		{name: "multiple text parts", output: message("", textPart(valid)+","+textPart(valid)), wantErr: "expected one final JSON text part"},
		{name: "incomplete message", output: strings.Replace(final, `"status":"completed"`, `"status":"incomplete"`, 1), wantErr: "final output message did not complete"},
	} {
		t.Run(test.name, func(t *testing.T) {
			status := test.status
			if status == "" {
				status = "completed"
			}
			server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost || r.URL.Path != "/responses" {
					t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
				}
				w.Header().Set("Content-Type", "application/json")
				if _, err := fmt.Fprintf(w, `{"id":"resp_fixture","object":"response","status":%q,"output":[%s]}`, status, test.output); err != nil {
					t.Errorf("write fixture: %v", err)
				}
			}))
			defer server.Close()
			client := openai.NewClient(
				option.WithAPIKey("test-key"),
				option.WithBaseURL(server.URL),
				option.WithHTTPClient(server.Client()),
				option.WithMaxRetries(0),
			)
			schema, err := GenerateSchema[HistoricalComputer]()
			if err != nil {
				t.Fatal(err)
			}
			response, err := client.Responses.New(context.Background(), historicalComputerParams("test question", schema))
			if err != nil {
				t.Fatal(err)
			}
			computer, err := historicalComputerFromResponse(response)
			if test.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), test.wantErr) || computer != nil {
					t.Fatalf("got computer=%v, err=%v; want nil computer and error containing %q", computer, err, test.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if computer.Name != "SNARC" || computer.Origin.YearBuilt != 1951 || len(computer.NotableFacts) != 1 {
				t.Fatalf("unexpected decoded computer: %+v", computer)
			}
		})
	}
}

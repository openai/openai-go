package main

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/openai/openai-go/v3"
)

func TestBetaAgentOutputUsesExistingGeneratedSchema(t *testing.T) {
	schema, err := GenerateSchema[HistoricalComputer]()
	if err != nil {
		t.Fatal(err)
	}
	output, err := openai.NewBetaAgentOutput(schema, func(data []byte) (HistoricalComputer, error) {
		var result HistoricalComputer
		err := json.Unmarshal(data, &result)
		return result, err
	})
	if err != nil {
		t.Fatal(err)
	}
	params := openai.BetaAgentSessionNewParams{Agent: openai.BetaAgentSessionNewParamsAgent{Text: openai.AgentTextParam{Format: output.Format()}}}
	wire, err := json.Marshal(params)
	if err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{`"$id":`, `"$schema":`, `"origin":`, `"notable_facts":`, `"enum":["positive","neutral","negative"]`} {
		if !strings.Contains(string(wire), expected) {
			t.Errorf("missing %s in %s", expected, wire)
		}
	}
	format, err := json.Marshal(output.Format())
	if err != nil {
		t.Fatal(err)
	}
	var envelope map[string]json.RawMessage
	if err := json.Unmarshal(format, &envelope); err != nil {
		t.Fatal(err)
	}
	if len(envelope) != 2 || string(envelope["type"]) != `"json_schema"` {
		t.Fatalf("wrong Agents format: %s", format)
	}
}

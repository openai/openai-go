package openai_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type betaReport struct {
	Summary string `json:"summary"`
}

func betaReportParser(data []byte) (betaReport, error) {
	var value betaReport
	err := json.Unmarshal(data, &value)
	if err == nil && value.Summary == "" {
		err = errors.New("summary is required")
	}
	return value, err
}
func betaOutputAdapter(t *testing.T) *openai.BetaAgentOutput[betaReport] {
	t.Helper()
	output, err := openai.NewBetaAgentOutput(map[string]any{"type": "object", "properties": map[string]any{"summary": map[string]any{"type": "string"}}}, betaReportParser)
	if err != nil {
		t.Fatal(err)
	}
	return output
}
func betaOutputEvents(text string) []string {
	return []string{betaResultTurn("created", "root", "null"), betaResultMessage("done", "answer", "root", `"final_answer"`, text, 0), betaResultTurn("completed", "root", "null"), betaResultIdle()}
}
func TestBetaAgentTypedCreationWireAndResult(t *testing.T) {
	output := betaOutputAdapter(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Error(err)
		}
		format := body["agent"].(map[string]any)["text"].(map[string]any)["format"].(map[string]any)
		if len(format) != 2 || format["type"] != "json_schema" {
			t.Errorf("wrong Agents format: %#v", format)
		}
		schema := format["schema"].(map[string]any)
		if schema["additionalProperties"] != false || len(schema["required"].([]any)) != 1 {
			t.Errorf("not normalized: %#v", schema)
		}
		if _, ok := body["output_type"]; ok {
			t.Error("local parser sent on wire")
		}
		w.Header().Set("Content-Type", "text/event-stream")
		for _, event := range betaOutputEvents(`{"summary":"ready"}`) {
			_, _ = fmt.Fprintf(w, "data: %s\n\n", event)
		}
	}))
	defer server.Close()
	client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"))
	stream := client.Beta.Agents.Sessions.NewStreaming(context.Background(), openai.BetaAgentSessionNewParams{Agent: openai.BetaAgentSessionNewParamsAgent{Text: openai.AgentTextParam{Format: output.Format()}}})
	defer func() { _ = stream.Close() }()
	result, err := output.FinalResult(stream)
	if err != nil || result.OutputParsed.Summary != "ready" || result.OutputText() != `{"summary":"ready"}` || result.TurnID() != "root" {
		t.Fatalf("result=%#v err=%v", result, err)
	}
	raw, err := openai.BetaAgentSessionFinalResult(stream)
	if err != nil || raw != result.BetaAgentTurnResult {
		t.Fatal("raw result was not preserved", err)
	}
}
func TestBetaAgentTypedFollowupAndErrors(t *testing.T) {
	output := betaOutputAdapter(t)
	_, client := newAgentHelperServer(t, betaOutputEvents(`{"summary":"followup"}`)...)
	stream := client.Beta.Agents.Sessions.Stream(context.Background(), "session", openai.AgentSessionStreamParams{Input: "continue"}).WithResultCollection()
	if !stream.Next() {
		t.Fatal(stream.Err())
	}
	result, err := output.FinalResult(stream)
	if err != nil || result.OutputParsed.Summary != "followup" {
		t.Fatalf("result=%v err=%v", result, err)
	}
	for _, text := range []string{`not json`, `{}`, `{"summary":42}`} {
		t.Run(text, func(t *testing.T) {
			rawStream := betaResultCreateStream(t, betaOutputEvents(text), false)
			_, parseErr := output.FinalResult(rawStream)
			var parseError *openai.BetaAgentOutputParseError
			if !errors.As(parseErr, &parseError) || parseError.Result.OutputText() != text || parseError.Result.Turn.Status != "completed" {
				t.Fatalf("lost completed result: %v", parseErr)
			}
		})
	}
	failed := betaResultCreateStream(t, []string{betaResultTurn("created", "root", "null"), betaResultTurn("failed", "root", "null")}, false)
	_, err = output.FinalResult(failed)
	var failure *openai.BetaAgentTurnResultError
	if !errors.As(err, &failure) || failure.Reason != "turn_failed" {
		t.Fatalf("hosted error became parse error: %v", err)
	}
}
func TestBetaAgentOutputSchemaValidation(t *testing.T) {
	invalid := []string{
		`{"type":"array","items":{"type":"string"}}`,
		`{"type":"object","properties":{},"anyOf":[]}`,
		`{"type":"object","properties":{"x":{"type":"string","default":"x"}}}`,
		`{"type":"object","properties":{},"additionalProperties":true}`,
		`{"type":"object","properties":{"x":{"$ref":"https://example.com/schema"}}}`,
		`{"type":"object","properties":{"x":{"$ref":"#/$defs/missing"}}}`,
		`{"type":"object","properties":{"x":{"type":"array","items":[]}}}`,
		`{"type":"object","properties":{"x":{"type":"string","pattern":"x"}}}`,
		`{"type":"object","properties":{"x":{"type":["string","integer"]}}}`,
	}
	for _, schema := range invalid {
		t.Run(schema, func(t *testing.T) {
			_, err := openai.NewBetaAgentOutput(json.RawMessage(schema), betaReportParser)
			if err == nil {
				t.Fatal("accepted unsupported schema")
			}
		})
	}
	supported := json.RawMessage(`{"type":"object","$defs":{"name":{"type":["string","null"]}},"properties":{"nested":{"type":"object","properties":{"names":{"type":"array","items":{"$ref":"#/$defs/name"}}}},"choice":{"anyOf":[{"type":"integer","enum":[1,2]},{"type":"null"}]}}}`)
	output, err := openai.NewBetaAgentOutput(supported, betaReportParser)
	if err != nil {
		t.Fatal(err)
	}
	first := output.Format()
	first.OfParamJSONSchema.Schema["type"] = "array"
	raw, err := json.Marshal(output.Format())
	if err != nil || !strings.Contains(string(raw), `"type":"object"`) {
		t.Fatal("caller mutation escaped schema copy", err)
	}
	if _, err := openai.NewBetaAgentOutput[betaReport](supported, nil); err == nil {
		t.Fatal("nil parser accepted")
	}
}

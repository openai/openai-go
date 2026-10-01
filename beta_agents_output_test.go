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
		`{"type":"object","properties":{"x":{"type":"string","format":"uri"}}}`,
		`{"type":"object","properties":{"x":{"type":"object","patternProperties":{".*":{"type":"string"}}}}}`,
		`{"type":"object","title":12}`,
		`{"type":"object","description":null}`,
		`{"type":"object","examples":{}}`,
		`{"type":"object","properties":{"a\"b":{"type":"string"}}}`,
		`{"type":"object","properties":{"x":{"type":"string","enum":["a\nb"]}}}`,
		`{"type":"object","$defs":{"a/b":{"type":"string"}},"properties":{"x":{"$ref":"#/$defs/a~1b"}}}`,
		`{"type":"object","$defs":{"a/b":{"type":"string"}},"properties":{"x":{"$ref":"#/$defs/a/b"}}}`,
		`{"type":"object","properties":{"x":{"oneOf":[{"type":"string"}]}}}`,
		`{"type":"object","$defs":{"name":{"type":"string"}},"properties":{"x":{"$ref":"#/$defs/name","description":"name"}}}`,
		`{"type":"object","properties":{"x":{"type":"string","enum":[1]}}}`,
		`{"type":"object","properties":{"x":{"type":"integer","const":1.0000000000000001}}}`,
		`{"type":"object","properties":{"x":{"type":"object","properties":{},"const":{}}}}`,
		`{"type":"array","items":{"type":"string"}}`,
		`{"type":"object","properties":{},"anyOf":[]}`,
		`{"type":"object","properties":{"x":{"type":"string","allOf":[]}}}`,
		`{"type":"object","properties":{},"additionalProperties":true}`,
		`{"type":"object","properties":{"x":{"$ref":"https://example.com/schema"}}}`,
		`{"type":"object","properties":{"x":{"$ref":"#/$defs/missing"}}}`,
		`{"type":"object","properties":{"x":{"type":"array","items":[]}}}`,
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
	supported := json.RawMessage(`{"type":"object","default":{},"$defs":{"name":{"type":["string","null"]}},"properties":{"nested":{"type":"object","properties":{"names":{"type":"array","items":{"$ref":"#/$defs/name"}}}},"choice":{"anyOf":[{"type":"integer","enum":[1,18446744073709551615]},{"type":"null"}]}}}`)
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

func TestBetaAgentOutputPreservesSupportedSchema(t *testing.T) {
	schema := json.RawMessage(`{"type":"object","$id":"https://example.com/report","$defs":{"a~1b":{"type":"string","pattern":"^[a-z]+$","minLength":1,"maxLength":20}},"properties":{"name":{"$ref":"#/$defs/a~1b"},"empty":{"type":"object"},"choice":{"anyOf":[{"type":"string"},{"type":"null"}],"default":null,"examples":["yes"]},"count":{"type":"integer","minimum":0,"maximum":18446744073709551615,"multipleOf":1},"tuple":{"type":"array","items":{"type":"string"},"prefixItems":[{"type":"object"}],"minItems":1,"maxItems":3},"recursive":{"$ref":"#"}}}`)
	output, err := openai.NewBetaAgentOutput(schema, betaReportParser)
	if err != nil {
		t.Fatal(err)
	}
	wire, err := json.Marshal(output.Format())
	if err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{`"$id":"https://example.com/report"`, `"default":null`, `"examples":["yes"]`, `"maximum":18446744073709551615`, `"pattern":"^[a-z]+$"`, `"properties":{}`, `"required":[]`} {
		if !strings.Contains(string(wire), expected) {
			t.Errorf("missing %s in %s", expected, wire)
		}
	}
}

func TestBetaAgentOutputParseErrorDoesNotExposeOutput(t *testing.T) {
	const secret = "private-output-canary"
	output, err := openai.NewBetaAgentOutput(map[string]any{"type": "object"}, func([]byte) (betaReport, error) { return betaReport{}, errors.New(secret) })
	if err != nil {
		t.Fatal(err)
	}
	_, err = output.FinalResult(betaResultCreateStream(t, betaOutputEvents(secret), false))
	if err == nil || strings.Contains(fmt.Sprint(err), secret) {
		t.Fatalf("unsafe parse error: %v", err)
	}
	if errors.Unwrap(err).Error() != secret {
		t.Fatal("parser cause unavailable")
	}
}

func TestBetaAgentTypedOutputWithTypedToolHandler(t *testing.T) {
	type LookupArgs struct {
		ID string `json:"id"`
	}
	type Receipt struct {
		Summary string `json:"summary"`
	}
	calls := 0
	handler := func(_ context.Context, arguments map[string]any) (any, error) {
		data, err := json.Marshal(arguments)
		if err != nil {
			return nil, err
		}
		var args LookupArgs
		if parseErr := json.Unmarshal(data, &args); parseErr != nil {
			return nil, parseErr
		}
		if args.ID != "A123" {
			return nil, errors.New("invalid lookup ID")
		}
		calls++
		data, err = json.Marshal(Receipt{Summary: "found"})
		return string(data), err
	}
	events := []string{betaResultTurn("created", "root", "null"), agentCall("call-event", "root", "call", "lookup", `{"id":"A123"}`), betaResultMessage("done", "answer", "root", `"final_answer"`, `{"summary":"found"}`, 0), betaResultTurn("completed", "root", "null"), betaResultIdle()}
	mock, client := newAgentHelperServer(t, events...)
	stream := client.Beta.Agents.Sessions.Stream(context.Background(), "session", openai.AgentSessionStreamParams{Input: "lookup", ToolHandlers: map[string]openai.AgentToolHandler{"lookup": handler}})
	result, err := betaOutputAdapter(t).FinalResult(stream)
	if err != nil || result.OutputParsed.Summary != "found" || calls != 1 {
		t.Fatalf("result=%v calls=%d err=%v", result, calls, err)
	}
	posts := mock.submissions()
	if len(posts) != 2 {
		t.Fatalf("expected input and result posts, got %d", len(posts))
	}
	wire, err := json.Marshal(posts[1].body)
	if err != nil || len(posts) != 2 || !strings.Contains(string(wire), `tool_result`) || !strings.Contains(string(wire), `found`) {
		t.Fatalf("typed receipt not submitted: %s err=%v", wire, err)
	}
}

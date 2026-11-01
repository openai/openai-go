package websockettest

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/openai/openai-go/v3/responses"
)

func addMCPProgressEvent(t *testing.T, a *responses.ResponseAccumulator, raw string) responses.ResponsesServerEventUnion {
	t.Helper()
	var event responses.ResponsesServerEventUnion
	if err := json.Unmarshal([]byte(raw), &event); err != nil {
		t.Fatal(err)
	}
	if err := a.AddEvent(event); err != nil {
		t.Fatal(err)
	}
	if event.RawJSON() != raw {
		t.Fatal("accumulator mutated the received event")
	}
	return event
}

func TestResponsesAccumulatorMCPArgumentProgress(t *testing.T) {
	var a responses.ResponseAccumulator
	addMCPProgressEvent(t, &a, `{"type":"response.created","sequence_number":0,"response":{"id":"resp_fixture","object":"response","access_programs":{"cyber":"standard"},"created_at":0,"status":"in_progress","error":null,"incomplete_details":null,"instructions":null,"metadata":{},"model":"fixture_model","output":[],"parallel_tool_calls":true,"temperature":1,"tool_choice":"auto","tools":[],"top_p":1}}`)
	addMCPProgressEvent(t, &a, `{"type":"response.output_item.added","sequence_number":1,"output_index":0,"item":{"id":"mcp_fixture","type":"mcp_call","name":"lookup","server_label":"fixture","arguments":"{","status":"in_progress","approval_request_id":null,"error":null,"output":null}}`)
	saved := a.Snapshot()
	detailed := a.DetailedSnapshot()
	for _, step := range []struct{ name, raw, want string }{
		{"first delta", `{"type":"response.mcp_call_arguments.delta","sequence_number":2,"output_index":0,"item_id":"mcp_fixture","delta":"\"value\":"}`, `{"value":`},
		{"second delta", `{"type":"response.mcp_call_arguments.delta","sequence_number":3,"output_index":0,"item_id":"mcp_fixture","delta":"1}"}`, `{"value":1}`},
		{"done replaces delta", `{"type":"response.mcp_call_arguments.done","sequence_number":4,"output_index":0,"item_id":"mcp_fixture","arguments":"{\"value\":2}"}`, `{"value":2}`},
	} {
		addMCPProgressEvent(t, &a, step.raw)
		got := a.Snapshot()
		if len(got.Output) != 1 {
			t.Fatal("missing MCP item")
		}
		if got.Output[0].Arguments != step.want {
			t.Errorf("%s: compact arguments = %q, want %q", step.name, got.Output[0].Arguments, step.want)
		}
		var arguments string
		if err := json.Unmarshal(a.DetailedSnapshot().Output[0].Item["arguments"], &arguments); err != nil {
			t.Fatal(err)
		}
		if arguments != step.want {
			t.Errorf("%s: detailed arguments = %q, want %q", step.name, arguments, step.want)
		}
	}
	if saved.Output[0].Arguments != "{" || string(detailed.Output[0].Item["arguments"]) != `"{"` {
		t.Fatal("earlier snapshots changed after argument progress")
	}
	addMCPProgressEvent(t, &a, `{"type":"response.output_item.done","sequence_number":5,"output_index":0,"item":{"id":"mcp_fixture","type":"mcp_call","name":"lookup","server_label":"fixture","arguments":"{\"value\":3}","status":"completed","approval_request_id":null,"error":null,"output":"synthetic result"}}`)
	full := a.Snapshot()
	fullDetails := a.DetailedSnapshot()
	if full.Output[0].Arguments != `{"value":3}` || string(fullDetails.Output[0].Item["arguments"]) != `"{\"value\":3}"` || string(fullDetails.Output[0].Item["server_label"]) != `"fixture"` || string(fullDetails.Output[0].Item["output"]) != `"synthetic result"` {
		t.Fatal("full item did not replace argument progress and preserve metadata")
	}
	terminal := addMCPProgressEvent(t, &a, `{"type":"response.completed","sequence_number":6,"response":{"id":"resp_fixture","object":"response","access_programs":{"cyber":"standard"},"created_at":0,"status":"completed","error":null,"incomplete_details":null,"instructions":null,"metadata":{},"model":"fixture_model","output":[{"id":"mcp_fixture","type":"mcp_call","name":"lookup","server_label":"fixture","arguments":"{\"value\":4}","status":"completed","approval_request_id":null,"error":null,"output":"terminal result"}],"parallel_tool_calls":true,"temperature":1,"tool_choice":"auto","tools":[],"top_p":1}}`)
	if a.Snapshot().Output[0].Arguments != `{"value":4}` || string(a.DetailedSnapshot().Output[0].Item["arguments"]) != `"{\"value\":4}"` || a.Snapshot().TerminalEvent != "response.completed" {
		t.Fatal("terminal arguments did not replace the full item")
	}
	if terminal.OfResponsesServerEventResponseWsCompleted.Response.Output[0].Arguments.OfString != `{"value":4}` || full.Output[0].Arguments != `{"value":3}` || string(fullDetails.Output[0].Item["arguments"]) != `"{\"value\":3}"` {
		t.Fatal("terminal or retained full-item snapshot changed")
	}
}

func TestResponsesAccumulatorMCPCompleteItemsAndTerminalAuthority(t *testing.T) {
	var a responses.ResponseAccumulator
	addMCPProgressEvent(t, &a, `{"type":"response.output_item.done","sequence_number":3,"output_index":0,"item":{"id":"mcp_fixture","type":"mcp_call","name":"lookup","server_label":"fixture","arguments":"{\"value\":3}","status":"completed","approval_request_id":null,"error":null,"output":"synthetic result"}}`)
	saved := a.DetailedSnapshot()
	if string(saved.Output[0].Item["arguments"]) != `"{\"value\":3}"` || string(saved.Output[0].Item["server_label"]) != `"fixture"` || string(saved.Output[0].Item["output"]) != `"synthetic result"` {
		t.Fatal("complete MCP item lost its arguments or metadata")
	}
	terminal := addMCPProgressEvent(t, &a, `{"type":"response.completed","sequence_number":4,"response":{"id":"resp_fixture","object":"response","access_programs":{"cyber":"standard"},"created_at":0,"status":"completed","error":null,"incomplete_details":null,"instructions":null,"metadata":{},"model":"fixture_model","output":[],"parallel_tool_calls":true,"temperature":1,"tool_choice":"auto","tools":[],"top_p":1}}`)
	if a.Snapshot().TerminalEvent != "response.completed" || len(a.Snapshot().Output) != 0 || len(a.DetailedSnapshot().Output) != 0 {
		t.Fatal("authoritative empty terminal output did not replace prior MCP output")
	}
	if len(terminal.OfResponsesServerEventResponseWsCompleted.Response.Output) != 0 || terminal.OfResponsesServerEventResponseWsCompleted.Response.JSON.Output.Raw() != "[]" || string(saved.Output[0].Item["arguments"]) != `"{\"value\":3}"` {
		t.Fatal("terminal or earlier detailed snapshot was rewritten")
	}
}

func TestResponsesAccumulatorMCPProgressRequiresCurrentScaffold(t *testing.T) {
	for _, tc := range []struct{ name, scaffold string }{
		{"no scaffold", ""},
		{"wrong kind", `{"type":"response.output_item.added","sequence_number":0,"output_index":0,"item":{"id":"mcp_fixture","type":"message","role":"assistant","status":"in_progress","content":[]}}`},
		{"newer MCP item", `{"type":"response.output_item.added","sequence_number":0,"output_index":0,"item":{"id":"mcp_new","type":"mcp_call","name":"lookup","server_label":"new_fixture","arguments":"new","status":"in_progress"}}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var a responses.ResponseAccumulator
			if tc.scaffold != "" {
				addMCPProgressEvent(t, &a, tc.scaffold)
			}
			before, details := a.Snapshot(), a.DetailedSnapshot()
			for _, raw := range []string{
				`{"type":"response.mcp_call_arguments.delta","sequence_number":1,"output_index":0,"item_id":"mcp_fixture","delta":"stale"}`,
				`{"type":"response.mcp_call_arguments.done","sequence_number":2,"output_index":0,"item_id":"mcp_fixture","arguments":"stale done"}`,
			} {
				addMCPProgressEvent(t, &a, raw)
			}
			if !reflect.DeepEqual(before, a.Snapshot()) || !reflect.DeepEqual(details, a.DetailedSnapshot()) {
				t.Fatal("MCP progress created or replaced another item")
			}
		})
	}
}

func TestResponsesAccumulatorMCPMalformedProgressDoesNotMutate(t *testing.T) {
	for _, raw := range []string{
		`{"type":"response.mcp_call_arguments.delta","sequence_number":1,"output_index":-1,"item_id":"mcp_fixture","delta":"bad"}`,
		`{"type":"response.mcp_call_arguments.delta","sequence_number":1,"output_index":"0","item_id":"mcp_fixture","delta":"bad"}`,
		`{"type":"response.mcp_call_arguments.delta","sequence_number":1,"output_index":0,"item_id":7,"delta":"bad"}`,
		`{"type":"response.mcp_call_arguments.delta","sequence_number":1,"output_index":0,"item_id":"mcp_fixture","delta":7}`,
		`{"type":"response.mcp_call_arguments.done","sequence_number":1,"output_index":0,"item_id":"mcp_fixture"}`,
		`{"type":"response.mcp_call_arguments.done","sequence_number":1,"output_index":0,"item_id":"mcp_fixture","arguments":null}`,
	} {
		var a responses.ResponseAccumulator
		addMCPProgressEvent(t, &a, `{"type":"response.output_item.added","sequence_number":0,"output_index":0,"item":{"id":"mcp_fixture","type":"mcp_call","name":"lookup","server_label":"fixture","arguments":"seed","status":"in_progress"}}`)
		before, details := a.Snapshot(), a.DetailedSnapshot()
		var event responses.ResponsesServerEventUnion
		if err := json.Unmarshal([]byte(raw), &event); err != nil {
			t.Fatal(err)
		}
		if err := a.AddEvent(event); err == nil {
			t.Error("accepted malformed selected MCP progress")
		}
		if event.RawJSON() != raw || !reflect.DeepEqual(before, a.Snapshot()) || !reflect.DeepEqual(details, a.DetailedSnapshot()) {
			t.Fatal("malformed MCP progress changed observed data")
		}
	}
}

func TestResponsesAccumulatorMCPRawAdmissionRemainsPermissive(t *testing.T) {
	for _, item := range []string{
		`{"id":"mcp_fixture","type":"mcp_call","arguments":{},"future_fixture":{"kept":true}}`,
		`{"id":"approval_fixture","type":"mcp_approval_request","arguments":"approval","name":"lookup","server_label":"fixture","future_fixture":{"kept":true}}`,
	} {
		var a responses.ResponseAccumulator
		addMCPProgressEvent(t, &a, `{"type":"response.output_item.added","sequence_number":0,"output_index":0,"item":`+item+`}`)
		var fields map[string]json.RawMessage
		if err := json.Unmarshal([]byte(item), &fields); err != nil {
			t.Fatal(err)
		}
		got := a.DetailedSnapshot().Output[0].Item
		if !reflect.DeepEqual(got, fields) {
			t.Fatal("raw MCP admission or unselected variant changed")
		}
	}
}

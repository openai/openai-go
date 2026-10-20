package websockettest

import (
	"encoding/json"
	"testing"

	"github.com/openai/openai-go/v3/responses"
)

func TestResponsesNestedSteeringDispatch(t *testing.T) {
	for _, test := range []struct{ kind, want string }{
		{"message", "message"},
		{"function_call_output", "function_call_output"},
		{"future_input", ""},
	} {
		t.Run("input/"+test.kind, func(t *testing.T) {
			raw := `{"type":"` + test.kind + `","role":"user","content":"hello","call_id":"call_1","output":"result","future":7}`
			var input responses.ResponseSteerInputResponseSteerInputItemListUnion
			if err := json.Unmarshal([]byte(raw), &input); err != nil {
				t.Fatal(err)
			}
			var got string
			switch variant := input.AsAny().(type) {
			case responses.ResponseSteerInputResponseSteerInputItemListMessage:
				got = "message"
				if variant.Content.AsString() != "hello" {
					t.Fatalf("lost message content: %#v", variant.Content)
				}
			case responses.ResponseSteerInputResponseSteerInputItemListFunctionCallOutput:
				got = "function_call_output"
				if variant.CallID != "call_1" {
					t.Fatalf("lost call ID: %q", variant.CallID)
				}
			case nil:
			default:
				t.Fatalf("unexpected variant: %T", variant)
			}
			if got != test.want {
				t.Fatalf("variant %q, want %q", got, test.want)
			}
			if input.RawJSON() != raw {
				t.Fatalf("raw JSON changed: %s", input.RawJSON())
			}
		})
	}
	for _, kind := range []string{
		"function_call_output", "custom_tool_call_output", "computer_call_output",
		"shell_call_output", "apply_patch_call_output", "tool_search_output", "mcp_approval_response", "future_required_input",
	} {
		t.Run("required/"+kind, func(t *testing.T) {
			raw := `{"type":"` + kind + `","call_id":"call_2","name":"tool","approval_request_id":"approval_1","future":true}`
			var input responses.ResponseSteerRequiredInputUnion
			if err := json.Unmarshal([]byte(raw), &input); err != nil {
				t.Fatal(err)
			}
			var got string
			switch variant := input.AsAny().(type) {
			case responses.ResponseSteerRequiredInputFunctionCallOutput:
				got = "function_call_output"
			case responses.ResponseSteerRequiredInputCustomToolCallOutput:
				got = "custom_tool_call_output"
			case responses.ResponseSteerRequiredInputComputerCallOutput:
				got = "computer_call_output"
			case responses.ResponseSteerRequiredInputShellCallOutput:
				got = "shell_call_output"
			case responses.ResponseSteerRequiredInputApplyPatchCallOutput:
				got = "apply_patch_call_output"
			case responses.ResponseSteerRequiredInputToolSearchOutput:
				got = "tool_search_output"
			case responses.ResponseSteerRequiredInputMcpApprovalResponse:
				got = "mcp_approval_response"
			case nil:
			default:
				t.Fatalf("unexpected variant: %T", variant)
			}
			want := kind
			if kind == "future_required_input" {
				want = ""
			}
			if got != want {
				t.Fatalf("variant %q, want %q", got, want)
			}
			if input.RawJSON() != raw {
				t.Fatalf("raw JSON changed: %s", input.RawJSON())
			}
		})
	}
}

package live_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/openai/openai-go/v3/live"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/packages/respjson"
)

func TestLiveToolChoiceParamsCompatibility(t *testing.T) {
	function := live.ResponsesDelegationConfigToolChoiceFunctionParam{Name: "function"}
	function.SetExtraFields(map[string]any{"future": "function-value"})
	mcp := live.ResponsesDelegationConfigToolChoiceMcpParam{Name: param.NewOpt("tool"), ServerLabel: "box"}
	mcp.SetExtraFields(map[string]any{"future": "mcp-value"})

	for _, tt := range []struct {
		name  string
		param json.Marshaler
		want  map[string]any
	}{
		{
			name:  "function extra fields",
			param: function,
			want:  map[string]any{"name": "function", "type": "function", "future": "function-value"},
		},
		{
			name:  "mcp extra fields",
			param: mcp,
			want:  map[string]any{"name": "tool", "server_label": "box", "type": "mcp", "future": "mcp-value"},
		},
		{
			name: "function override",
			param: param.Override[live.ResponsesDelegationConfigToolChoiceFunctionParam](
				json.RawMessage(`{"name":"function","type":"function","future":"raw-function"}`)),
			want: map[string]any{"name": "function", "type": "function", "future": "raw-function"},
		},
		{
			name: "mcp override",
			param: param.Override[live.ResponsesDelegationConfigToolChoiceMcpParam](
				json.RawMessage(`{"name":"tool","type":"mcp","server_label":"box","future":"raw-mcp"}`)),
			want: map[string]any{"name": "tool", "server_label": "box", "type": "mcp", "future": "raw-mcp"},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			data, err := tt.param.MarshalJSON()
			if err != nil {
				t.Fatal(err)
			}
			var got map[string]any
			if err := json.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tool choice = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestLiveResponsesDelegationMCPToolChoiceName(t *testing.T) {
	for _, tc := range []struct {
		name  string
		value param.Opt[string]
		want  map[string]any
		raw   string
	}{
		{"omitted", param.Opt[string]{}, map[string]any{"type": "mcp", "server_label": "box"}, respjson.Omitted},
		{"null", param.Null[string](), map[string]any{"type": "mcp", "server_label": "box", "name": nil}, respjson.Null},
		{"value", param.NewOpt("lookup"), map[string]any{"type": "mcp", "server_label": "box", "name": "lookup"}, `"lookup"`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			create := live.ResponsesDelegationConfigParam{
				Model: "synthetic-model",
				ToolChoice: live.ResponsesDelegationConfigToolChoiceUnionParam{
					OfResponsesDelegationConfigToolChoiceMcp: &live.ResponsesDelegationConfigToolChoiceMcpParam{ServerLabel: "box", Name: tc.value},
				},
			}
			update := live.ResponsesDelegationUpdateConfigParam{
				ToolChoice: live.ResponsesDelegationUpdateConfigToolChoiceUnionParam{
					OfResponsesDelegationUpdateConfigToolChoiceMcp: &live.ResponsesDelegationUpdateConfigToolChoiceMcpParam{ServerLabel: "box", Name: tc.value},
				},
			}
			wantValue, wantValid := tc.want["name"].(string)
			for _, request := range []struct {
				name     string
				value    any
				toolName *string
			}{{"create", create, create.ToolChoice.GetName()}, {"update", update, update.ToolChoice.GetName()}} {
				t.Run(request.name, func(t *testing.T) {
					data, err := json.Marshal(request.value)
					if err != nil {
						t.Fatal(err)
					}
					var body map[string]any
					if err := json.Unmarshal(data, &body); err != nil {
						t.Fatal(err)
					}
					if !reflect.DeepEqual(body["tool_choice"], tc.want) {
						t.Fatalf("tool_choice = %#v, want %#v", body["tool_choice"], tc.want)
					}
					if wantValid {
						if request.toolName == nil || *request.toolName != wantValue {
							t.Fatalf("GetName() = %v, want %q", request.toolName, wantValue)
						}
					} else if request.toolName != nil {
						t.Fatalf("GetName() = %q, want nil", *request.toolName)
					}
				})
			}

			choiceJSON, err := json.Marshal(tc.want)
			if err != nil {
				t.Fatal(err)
			}
			frame := fmt.Sprintf(`{"type":"session.updated","event_id":"evt_test","session":{"id":"sess_test","expires_at":123,"model":"gpt-live-1","status":"active","delegation":{"type":"responses","responses":{"model":"synthetic-model","tool_choice":%s}}}}`, choiceJSON)
			var event live.ServerEventUnion
			if err := json.Unmarshal([]byte(frame), &event); err != nil {
				t.Fatal(err)
			}
			updated, ok := event.AsAny().(live.SessionUpdatedEvent)
			if !ok {
				t.Fatalf("event = %T, want SessionUpdatedEvent", event.AsAny())
			}
			choice := updated.Session.Delegation.AsSessionResourceDelegationResponses().Responses.ToolChoice.AsResponsesDelegationConfigToolChoiceMcp()
			if choice.ServerLabel != "box" || choice.Name != wantValue || choice.JSON.Name.Valid() != wantValid || choice.JSON.Name.Raw() != tc.raw {
				t.Fatalf("MCP choice = %#v; name valid = %t, raw = %q", choice, choice.JSON.Name.Valid(), choice.JSON.Name.Raw())
			}
			if choice.RawJSON() != string(choiceJSON) {
				t.Fatalf("RawJSON() = %s, want %s", choice.RawJSON(), choiceJSON)
			}
			var roundTrip live.ResponsesDelegationConfigToolChoiceMcp
			if err := json.Unmarshal([]byte(choice.RawJSON()), &roundTrip); err != nil {
				t.Fatal(err)
			}
			if roundTrip.Name != wantValue || roundTrip.JSON.Name.Raw() != tc.raw {
				t.Fatalf("round-trip name = %q, raw = %q", roundTrip.Name, roundTrip.JSON.Name.Raw())
			}
		})
	}
}

package live_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/openai/openai-go/v3/live"
	"github.com/openai/openai-go/v3/packages/param"
)

func TestLiveToolChoiceParamsCompatibility(t *testing.T) {
	function := live.ResponsesDelegationConfigToolChoiceLiveFunctionToolChoiceParam{Name: "function"}
	function.SetExtraFields(map[string]any{"future": "function-value"})
	mcp := live.ResponsesDelegationConfigToolChoiceLiveMcpToolChoiceParam{Name: "tool", ServerLabel: "box"}
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
			param: param.Override[live.ResponsesDelegationConfigToolChoiceLiveFunctionToolChoiceParam](
				json.RawMessage(`{"name":"function","type":"function","future":"raw-function"}`)),
			want: map[string]any{"name": "function", "type": "function", "future": "raw-function"},
		},
		{
			name: "mcp override",
			param: param.Override[live.ResponsesDelegationConfigToolChoiceLiveMcpToolChoiceParam](
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

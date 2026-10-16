package live_test

import (
	"encoding/json"
	"testing"

	"github.com/openai/openai-go/v3/live"
	"github.com/openai/openai-go/v3/packages/param"
)

func TestLiveAppendEventsNullableDelegation(t *testing.T) {
	for _, delegation := range []struct {
		name string
		id   param.Opt[string]
		want any
	}{
		{name: "general session", id: param.Null[string](), want: nil},
		{name: "client delegation", id: param.NewOpt("synthetic-delegation"), want: "synthetic-delegation"},
	} {
		t.Run(delegation.name, func(t *testing.T) {
			instructions := &live.InstructionsAppendEventParam{Content: "instructions", DelegationID: delegation.id}
			thinking := &live.ThinkingAppendEventParam{Content: "thinking", DelegationID: delegation.id}
			commentary := &live.CommentaryAppendEventParam{Content: "commentary", DelegationID: delegation.id}
			for _, event := range []struct {
				name    string
				value   any
				content string
			}{
				{"live instructions", live.ClientEventUnionParam{OfSessionInstructionsAppend: instructions}, "instructions"},
				{"live thinking", live.ClientEventUnionParam{OfSessionThinkingAppend: thinking}, "thinking"},
				{"live commentary", live.ClientEventUnionParam{OfSessionCommentaryAppend: commentary}, "commentary"},
				{"fork instructions", live.ForkClientEventUnionParam{OfSessionInstructionsAppend: instructions}, "instructions"},
				{"fork thinking", live.ForkClientEventUnionParam{OfSessionThinkingAppend: thinking}, "thinking"},
				{"fork commentary", live.ForkClientEventUnionParam{OfSessionCommentaryAppend: commentary}, "commentary"},
				{"sideband instructions", live.ConnectClientEventUnionParam{OfSessionInstructionsAppend: instructions}, "instructions"},
				{"sideband thinking", live.ConnectClientEventUnionParam{OfSessionThinkingAppend: thinking}, "thinking"},
				{"sideband commentary", live.ConnectClientEventUnionParam{OfSessionCommentaryAppend: commentary}, "commentary"},
			} {
				t.Run(event.name, func(t *testing.T) {
					data, err := json.Marshal(event.value)
					if err != nil {
						t.Fatal(err)
					}
					var got map[string]any
					if err := json.Unmarshal(data, &got); err != nil {
						t.Fatal(err)
					}
					if id, present := got["delegation_id"]; !present || id != delegation.want {
						t.Errorf("delegation_id = %v (present %v), want %v", id, present, delegation.want)
					}
					if got["content"] != event.content || got["type"] != "session."+event.content+".append" {
						t.Errorf("unexpected event: %s", data)
					}
				})
			}
		})
	}
}

package live_test

import (
	"encoding/json"
	"testing"

	"github.com/openai/openai-go/v3/live"
)

func TestMCPFilterPresence(t *testing.T) {
	for _, tc := range []struct {
		name, allowed, approval, namesRaw, ruleRaw string
	}{
		{"omitted", "{}", "{}", "", ""},
		{"null", `{"tool_names":null}`, `{"always":null,"never":null}`, "null", "null"},
		{"empty", `{"tool_names":[]}`, `{"always":{},"never":{}}`, "[]", "{}"},
		{"populated", `{"tool_names":["lookup"]}`, `{"always":{"tool_names":["lookup"]},"never":{"tool_names":["lookup"]}}`, `["lookup"]`, `{"tool_names":["lookup"]}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var cfg live.ResponsesDelegationConfig
			payload := `{"tools":[{"type":"mcp","allowed_tools":` + tc.allowed + `,"require_approval":` + tc.approval + `}]}`
			if err := json.Unmarshal([]byte(payload), &cfg); err != nil {
				t.Fatal(err)
			}
			mcp := cfg.Tools[0].AsResponsesDelegationConfigToolMcp()
			filter := mcp.AllowedTools.AsMcpToolFilter()
			rules := mcp.RequireApproval.AsMcpRequireApprovalFilter()
			if got := filter.JSON.ToolNames.Raw(); got != tc.namesRaw {
				t.Errorf("tool_names = %q, want %q", got, tc.namesRaw)
			}
			if got := rules.JSON.Always.Raw(); got != tc.ruleRaw {
				t.Errorf("always = %q, want %q", got, tc.ruleRaw)
			}
			if got := rules.JSON.Never.Raw(); got != tc.ruleRaw {
				t.Errorf("never = %q, want %q", got, tc.ruleRaw)
			}
			if tc.name == "populated" && (rules.Always.ToolNames[0] != "lookup" || rules.Never.ToolNames[0] != "lookup") {
				t.Error("nested tool names were not decoded")
			}
		})
	}
}

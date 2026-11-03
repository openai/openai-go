package live_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	wire "github.com/coder/websocket"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/live"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/respjson"
)

func TestLiveResponsesDelegationReceivesNullableToolFields(t *testing.T) {
	configs := []string{
		`{"model":"gpt-4.1","tool_choice":{"type":"allowed_tools","mode":"required","tools":[{"type":"file_search"}]},"tools":[
			{"type":"shell","environment":{"type":"container_auto","file_ids":["file_test"],"memory_limit":"4g","network_policy":{"type":"allowlist","allowed_domains":["example.com"]},"skills":[{"type":"skill_reference","skill_id":"skill_value","version":"latest"},{"type":"skill_reference","skill_id":"skill_null","version":null},{"type":"skill_reference","skill_id":"skill_omitted"}]}},
			{"type":"shell","environment":{"type":"container_auto","file_ids":null,"memory_limit":null,"network_policy":null,"skills":null}},
			{"type":"shell","environment":{"type":"container_auto"}},
			{"type":"shell","environment":{"type":"local","skills":[{"name":"lookup","description":"Lookup information","path":"/skills/lookup"}]}},
			{"type":"shell","environment":{"type":"local","skills":null}},
			{"type":"shell","environment":{"type":"local"}},
			{"type":"shell","environment":null},
			{"type":"shell"}
		]}`,
		`{"model":"gpt-4.1","tool_choice":{"type":"allowed_tools","mode":null,"tools":[{"type":"file_search"}]}}`,
		`{"model":"gpt-4.1","tool_choice":{"type":"allowed_tools","tools":[{"type":"file_search"}]}}`,
	}
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	exercised := make(chan struct{}, 1)
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		socket, err := wire.Accept(w, r, nil)
		if err != nil {
			t.Errorf("accept: %v", err)
			return
		}
		defer func() { _ = socket.CloseNow() }()
		exercised <- struct{}{}
		for _, config := range configs {
			frame := fmt.Sprintf(`{"type":"session.started","event_id":"evt_test","session":{"id":"sess_test","expires_at":123,"model":"gpt-live-1","status":"active","delegation":{"type":"responses","responses":%s}}}`, config)
			if err := socket.Write(ctx, wire.MessageText, []byte(frame)); err != nil {
				t.Errorf("write: %v", err)
				return
			}
		}
		_, _, _ = socket.Read(ctx)
	}))
	defer server.Close()
	client := openai.NewClient(option.WithAPIKey("synthetic-test-key"), option.WithBaseURL(server.URL), option.WithHTTPClient(server.Client()))
	connection, err := client.Live.Connect(ctx, live.LiveConnectionOptions{})
	if err != nil {
		t.Fatal(err)
	}
	defer connection.Abort()

	receive := func(t *testing.T) live.ResponsesDelegationConfig {
		t.Helper()
		event, err := connection.Recv(ctx)
		if err != nil {
			t.Fatal(err)
		}
		return event.AsSessionStartedEvent().Session.Delegation.AsSessionResourceDelegationResponses().Responses
	}
	response := receive(t)
	if got := response.ToolChoice.AsResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParam().Mode; got != "required" {
		t.Fatalf("tool choice mode = %q, want required", got)
	}
	if len(response.Tools) != 8 {
		t.Fatalf("got %d tools, want 8", len(response.Tools))
	}
	container := response.Tools[0].AsResponsesDelegationConfigToolShell().Environment.AsResponsesDelegationConfigToolShellEnvironmentContainerAuto()
	if len(container.FileIDs) != 1 || container.FileIDs[0] != "file_test" || container.MemoryLimit != "4g" {
		t.Fatalf("container file IDs = %v, memory limit = %q", container.FileIDs, container.MemoryLimit)
	}
	policy := container.NetworkPolicy.AsResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyAllowlist()
	if len(policy.AllowedDomains) != 1 || policy.AllowedDomains[0] != "example.com" {
		t.Fatalf("allowed domains = %v", policy.AllowedDomains)
	}
	if len(container.Skills) != 3 {
		t.Fatalf("got %d container skills, want 3", len(container.Skills))
	}
	if got := container.Skills[0].AsResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillSkillReference().Version; got != "latest" {
		t.Fatalf("skill version = %q, want latest", got)
	}
	local := response.Tools[3].AsResponsesDelegationConfigToolShell().Environment.AsResponsesDelegationConfigToolShellEnvironmentLocal()
	if len(local.Skills) != 1 || local.Skills[0].Name != "lookup" {
		t.Fatalf("local skills = %+v", local.Skills)
	}
	for _, tc := range []struct {
		name                           string
		container, local, shell, skill int
		want                           string
	}{
		{"null", 1, 4, 6, 1, respjson.Null},
		{"omitted", 2, 5, 7, 2, respjson.Omitted},
	} {
		t.Run(tc.name, func(t *testing.T) {
			auto := response.Tools[tc.container].AsResponsesDelegationConfigToolShell().Environment.AsResponsesDelegationConfigToolShellEnvironmentContainerAuto()
			local := response.Tools[tc.local].AsResponsesDelegationConfigToolShell().Environment.AsResponsesDelegationConfigToolShellEnvironmentLocal()
			shell := response.Tools[tc.shell].AsResponsesDelegationConfigToolShell()
			skill := container.Skills[tc.skill].AsResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillSkillReference()
			for name, field := range map[string]respjson.Field{
				"file_ids": auto.JSON.FileIDs, "memory_limit": auto.JSON.MemoryLimit,
				"network_policy": auto.JSON.NetworkPolicy, "container skills": auto.JSON.Skills,
				"local skills": local.JSON.Skills, "environment": shell.JSON.Environment, "version": skill.JSON.Version,
			} {
				if field.Valid() || field.Raw() != tc.want {
					t.Errorf("%s valid = %t, raw = %q; want false, %q", name, field.Valid(), field.Raw(), tc.want)
				}
			}
			if got := receive(t).ToolChoice.AsResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParam().Mode; got != "" {
				t.Errorf("tool choice mode = %q, want empty", got)
			}
		})
	}
	select {
	case <-exercised:
	default:
		t.Fatal("the local WebSocket server was not exercised")
	}
}

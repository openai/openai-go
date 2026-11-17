package live_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	wire "github.com/coder/websocket"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/live"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/respjson"
)

func TestMCPFieldsSurviveSessionEventTransport(t *testing.T) {
	const tools = `[{"type":"mcp","server_label":"fixture","server_url":"https://example.test/mcp","server_description":"MCP response fixture","connector_id":"connector_dropbox","tunnel_id":"tunnel_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","allowed_callers":["direct"],"defer_loading":false,"allowed_tools":{"tool_names":["lookup"],"read_only":true},"require_approval":{"always":{"tool_names":["write"],"read_only":false},"never":{"tool_names":["lookup"],"read_only":true}}},{"type":"mcp","server_url":null,"server_description":null,"connector_id":null,"tunnel_id":null,"allowed_callers":null,"defer_loading":null,"allowed_tools":null,"require_approval":null},{"type":"mcp"},{"type":"mcp","allowed_tools":{"tool_names":null,"read_only":null},"require_approval":{"always":null,"never":null}},{"type":"mcp","allowed_tools":{},"require_approval":{}}]`
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
		for _, kind := range []string{"session.started", "session.updated"} {
			frame := fmt.Sprintf(`{"type":%q,"event_id":"event_fixture","session":{"id":"live_fixture","model":"gpt-live-1","status":"active","expires_at":123,"delegation":{"type":"responses","responses":{"model":"model_fixture","tools":%s}}}}`, kind, tools)
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
	for i := range 2 {
		event, err := connection.Recv(ctx)
		if err != nil {
			t.Fatal(err)
		}
		session := event.AsSessionStartedEvent().Session
		if i == 1 {
			session = event.AsSessionUpdatedEvent().Session
		}
		observed := session.Delegation.AsSessionResourceDelegationResponses().Responses.Tools
		if len(observed) != 5 {
			t.Fatalf("got %d MCP tools, want 5", len(observed))
		}
		tool := observed[0].AsResponsesDelegationConfigToolMcp()
		if tool.ServerDescription != "MCP response fixture" {
			t.Fatalf("server_description = %q", tool.ServerDescription)
		}
		if tool.ServerLabel != "fixture" || tool.ServerURL != "https://example.test/mcp" || tool.ConnectorID != "connector_dropbox" || tool.TunnelID != "tunnel_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" {
			t.Fatalf("server settings did not survive transport")
		}
		if !reflect.DeepEqual(tool.AllowedCallers, []string{"direct"}) || tool.DeferLoading || !tool.JSON.DeferLoading.Valid() {
			t.Fatal("caller settings or explicit false lost")
		}
		filter := tool.AllowedTools.AsMcpToolFilter()
		if !reflect.DeepEqual(filter.ToolNames, []string{"lookup"}) || !filter.ReadOnly {
			t.Fatal("allowed-tool filter lost")
		}
		rules := tool.RequireApproval.AsMcpRequireApprovalFilter()
		if !reflect.DeepEqual(rules.Always.ToolNames, []string{"write"}) || rules.Always.ReadOnly || !rules.Always.JSON.ReadOnly.Valid() || !reflect.DeepEqual(rules.Never.ToolNames, []string{"lookup"}) || !rules.Never.ReadOnly {
			t.Fatal("approval filters lost")
		}
		for index, want := range map[int]string{1: respjson.Null, 2: respjson.Omitted} {
			absent := observed[index].AsResponsesDelegationConfigToolMcp()
			for name, field := range map[string]respjson.Field{"server_url": absent.JSON.ServerURL, "server_description": absent.JSON.ServerDescription, "connector_id": absent.JSON.ConnectorID, "tunnel_id": absent.JSON.TunnelID, "allowed_callers": absent.JSON.AllowedCallers, "defer_loading": absent.JSON.DeferLoading, "allowed_tools": absent.JSON.AllowedTools, "require_approval": absent.JSON.RequireApproval} {
				if field.Valid() || field.Raw() != want {
					t.Errorf("%s presence = %q, want %q", name, field.Raw(), want)
				}
			}
		}
		for index, want := range map[int]string{3: respjson.Null, 4: respjson.Omitted} {
			nested := observed[index].AsResponsesDelegationConfigToolMcp()
			filter := nested.AllowedTools.AsMcpToolFilter()
			rules := nested.RequireApproval.AsMcpRequireApprovalFilter()
			for name, field := range map[string]respjson.Field{"tool_names": filter.JSON.ToolNames, "read_only": filter.JSON.ReadOnly, "always": rules.JSON.Always, "never": rules.JSON.Never} {
				if field.Valid() || field.Raw() != want {
					t.Errorf("%s presence = %q, want %q", name, field.Raw(), want)
				}
			}
		}
	}
	select {
	case <-exercised:
	default:
		t.Fatal("local WebSocket server was not exercised")
	}
}

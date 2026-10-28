package live_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	wire "github.com/coder/websocket"
	"github.com/openai/openai-go/v4"
	"github.com/openai/openai-go/v4/live"
	"github.com/openai/openai-go/v4/option"
)

func TestDelegationResponseNullableFields(t *testing.T) {
	for _, tc := range []struct {
		name, nameField, environmentField, wantName, wantContainer string
		null                                                       bool
	}{
		{name: "values", nameField: `,"name":"lookup_fixture"`, environmentField: `,"environment":{"type":"container_reference","container_id":"cntr_fixture"}`, wantName: "lookup_fixture", wantContainer: "cntr_fixture"},
		{name: "nulls", nameField: `,"name":null`, environmentField: `,"environment":null`, null: true},
		{name: "omitted"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			frame := []byte(fmt.Sprintf(`{"type":"session.updated","event_id":"event_fixture","session":{"id":"session_fixture","model":"synthetic-live","status":"active","expires_at":123,"delegation":{"type":"responses","responses":{"model":"synthetic-responses","tool_choice":{"type":"mcp","server_label":"fixture"%s},"tools":[{"type":"shell"%s}]}}}}`, tc.nameField, tc.environmentField))
			server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				socket, err := wire.Accept(w, r, nil)
				if err != nil {
					t.Error(err)
					return
				}
				defer func() { _ = socket.CloseNow() }()
				if err := socket.Write(r.Context(), wire.MessageText, frame); err != nil {
					t.Error(err)
					return
				}
				_, _, _ = socket.Read(r.Context())
			}))
			defer server.Close()
			client := openai.NewClient(option.WithAPIKey("synthetic"), option.WithBaseURL(server.URL), option.WithHTTPClient(server.Client()))
			ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
			defer cancel()
			connection, err := client.Live.Sideband.Connect(ctx, "session_fixture", live.SidebandConnectionOptions{})
			if err != nil {
				t.Fatal(err)
			}
			defer connection.Abort()
			event, err := connection.Recv(ctx)
			if err != nil {
				t.Fatal(err)
			}
			result := event.AsSessionUpdatedEvent().Session.Delegation.Responses
			mcp := result.ToolChoice.AsResponsesDelegationConfigToolChoiceMcp()
			if mcp.Name != tc.wantName {
				t.Errorf("Recv(%s).Mcp.Name = %q, want %q", tc.name, mcp.Name, tc.wantName)
			}
			if len(result.Tools) != 1 {
				t.Fatalf("Recv(%s).Tools length = %d, want 1", tc.name, len(result.Tools))
			}
			shell := result.Tools[0].AsResponsesDelegationConfigToolShell()
			container := shell.Environment.AsResponsesDelegationConfigToolShellEnvironmentContainerReference()
			if container.ContainerID != tc.wantContainer {
				t.Errorf("Recv(%s).ContainerID = %q, want %q", tc.name, container.ContainerID, tc.wantContainer)
			}
			if got := mcp.JSON.Name.Raw() == "null"; got != tc.null {
				t.Errorf("Recv(%s).Name null = %t, want %t", tc.name, got, tc.null)
			}
			if got := shell.JSON.Environment.Raw() == "null"; got != tc.null {
				t.Errorf("Recv(%s).Environment null = %t, want %t", tc.name, got, tc.null)
			}
			if tc.name == "omitted" && (mcp.JSON.Name.Valid() || shell.JSON.Environment.Valid()) {
				t.Error("Recv(omitted) has present nullable fields, want absent")
			}
		})
	}
}

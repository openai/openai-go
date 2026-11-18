package openai_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync/atomic"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func TestAgentSessionMCPApprovalResponse(t *testing.T) {
	cases := []struct {
		name      string
		wireName  any
		arguments any
		wantName  string
		isNull    bool
	}{
		{"non_null", "lookup_probe_492", map[string]any{"city": "probe_city_492", "limit": float64(7)}, "lookup_probe_492", false},
		{"explicit_null", nil, nil, "", true},
		{"empty_values", "", []any{}, "", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var requests atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requests.Add(1)
				if r.Method != http.MethodGet || r.URL.Path != "/agents/sessions/session_probe" {
					t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
					http.Error(w, "unexpected request", http.StatusBadRequest)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				body := map[string]any{
					"id": "session_probe", "object": "agent.session", "status": "requires_action",
					"required_actions": []any{map[string]any{
						"type": "mcp_approval_request", "request_id": "approval_probe", "turn_id": "turn_probe",
						"request": map[string]any{
							"type": "mcp_tool_call", "server_label": "probe_server", "message": "Approve this synthetic request",
							"name": tc.wireName, "arguments": tc.arguments,
						},
					}},
				}
				if err := json.NewEncoder(w).Encode(body); err != nil {
					t.Errorf("write response: %v", err)
				}
			}))
			t.Cleanup(server.Close)
			client := openai.NewClient(
				option.WithAPIKey("test-key"),
				option.WithUnsafeAllowHTTP(),
				option.WithBaseURL(server.URL),
				option.WithHTTPClient(server.Client()),
				option.WithMaxRetries(0),
			)
			session, err := client.Beta.Agents.Sessions.Get(t.Context(), "session_probe")
			if err != nil {
				t.Fatal(err)
			}
			if requests.Load() != 1 || len(session.RequiredActions) != 1 {
				t.Fatalf("requests = %d, actions = %d; want one of each", requests.Load(), len(session.RequiredActions))
			}
			approval := session.RequiredActions[0].AsMcpApprovalRequest()
			if approval.Type != "mcp_approval_request" || approval.RequestID != "approval_probe" {
				t.Fatalf("approval identity = %q, %q", approval.Type, approval.RequestID)
			}
			request := approval.Request
			if request.Name != tc.wantName {
				t.Errorf("name = %q; want %q", request.Name, tc.wantName)
			}
			if !reflect.DeepEqual(request.Arguments, tc.arguments) {
				t.Errorf("arguments = %#v; want %#v", request.Arguments, tc.arguments)
			}
			if request.JSON.Name.Valid() == tc.isNull || request.JSON.Arguments.Valid() == tc.isNull {
				t.Errorf("field validity: name = %v, arguments = %v; isNull = %v", request.JSON.Name.Valid(), request.JSON.Arguments.Valid(), tc.isNull)
			}
			if tc.isNull && (request.JSON.Name.Raw() != "null" || request.JSON.Arguments.Raw() != "null") {
				t.Errorf("null presence metadata: name = %q, arguments = %q", request.JSON.Name.Raw(), request.JSON.Arguments.Raw())
			}
		})
	}
}

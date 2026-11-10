package lib_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/param"
)

const automationFixture = `{"id":"auto_test","object":"agent.automation","type":"periodic","name":"offline schedule","expression":"0 9 * * *","timezone":"UTC","created_at":1700000000,"created_by":{"type":"developer"},"status":"active","metadata":{"marker":"preserved"},"target":{"type":"existing_session","session_id":"sess_target","instructions":"offline"}}`

func TestAutomationResponseMetadata(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/agents/automations/auto_test" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(automationFixture))
	}))
	defer server.Close()
	client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("test-key"))
	got, err := client.Beta.Agents.Automations.Get(context.Background(), "auto_test")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got.Metadata, map[string]string{"marker": "preserved"}) {
		t.Fatalf("Metadata = %#v, want marker: preserved map", got.Metadata)
	}
}

func TestAutomationCreateUnionBody(t *testing.T) {
	cases := []struct {
		name   string
		body   openai.AutomationScheduleParamUnion
		target string
	}{
		{name: "periodic", target: "existing_session", body: openai.AutomationScheduleParamUnion{OfParamPeriodic: &openai.AutomationScheduleParamPeriodic{
			Name: "periodic test", Expression: "0 9 * * *", Timezone: "UTC", Metadata: map[string]string{"marker": "request"},
			Target: openai.AutomationTargetParamUnion{OfParamExistingSession: &openai.AutomationTargetParamExistingSession{SessionID: "sess_target", Instructions: "offline"}},
		}}},
		{name: "once", target: "new_session", body: openai.AutomationScheduleParamUnion{OfParamOnce: &openai.AutomationScheduleParamOnce{
			Name: "once test", RunAt: 1800000000, Metadata: map[string]string{"marker": "request"},
			Target: openai.AutomationTargetParamUnion{OfParamNewSession: &openai.AutomationTargetParamNewSession{Environment: openai.EnvironmentParamUnion{OfParamNone: &openai.EnvironmentParamNone{}}, Input: []openai.AgentSessionInputMessageParam{}, AgentID: openai.String("agent_test")}},
		}}},
	}
	for _, tc := range cases {
		for _, mode := range []string{"typed", "raw_override", "extra_fields"} {
			t.Run(tc.name+"/"+mode, func(t *testing.T) {
				params := openai.BetaAgentAutomationNewParams{CreateAutomationParams: tc.body}
				switch mode {
				case "raw_override":
					raw, err := json.Marshal(tc.body)
					if err != nil {
						t.Fatal(err)
					}
					params = openai.BetaAgentAutomationNewParams{}
					param.SetJSON(raw, &params)
				case "extra_fields":
					params.SetExtraFields(map[string]any{"wrapper_marker": "retained", "name": "overridden"})
				}
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.Method != http.MethodPost || r.URL.Path != "/agents/automations" {
						t.Errorf("unexpected request: %s %s", r.Method, r.URL)
					}
					if r.Header.Get("OpenAI-Beta") != "agents=v1" {
						t.Errorf("missing Agents beta header")
					}
					var body map[string]any
					if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
						t.Error(err)
						w.WriteHeader(400)
						return
					}
					if body["type"] != tc.name {
						t.Errorf("request type = %v, want %s", body["type"], tc.name)
					}
					target, ok := body["target"].(map[string]any)
					if !ok || target["type"] != tc.target {
						t.Errorf("target = %#v, want %s", body["target"], tc.target)
					}
					if !reflect.DeepEqual(body["metadata"], map[string]any{"marker": "request"}) {
						t.Errorf("metadata = %#v", body["metadata"])
					}
					if mode == "extra_fields" && (body["wrapper_marker"] != "retained" || body["name"] != "overridden") {
						t.Errorf("wrapper extra fields were lost: %#v", body)
					}
					w.Header().Set("Content-Type", "application/json")
					_, _ = w.Write([]byte(automationFixture))
				}))
				defer server.Close()
				client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("test-key"))
				if _, err := client.Beta.Agents.Automations.New(context.Background(), params); err != nil {
					t.Fatal(err)
				}
			})
		}
	}
}

func TestAutomationInvocationNullableFields(t *testing.T) {
	for _, nulls := range []bool{false, true} {
		name := "values"
		if nulls {
			name = "nulls"
		}
		t.Run(name, func(t *testing.T) {
			fixture := map[string]any{"id": "inv_test", "automation_id": "auto_test", "object": "agent.automation.invocation", "created_at": 1700000001, "status": "failed", "source": map[string]any{"type": "manual"}, "scheduled_at": int64(1700000123), "session_id": "sess_preserved", "skip_reason": "pending_delivery", "error": map[string]string{"code": "synthetic_failure", "message": "distinctive explanation"}}
			if nulls {
				for _, key := range []string{"scheduled_at", "session_id", "skip_reason", "error"} {
					fixture[key] = nil
				}
			}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet || r.URL.Path != "/agents/automations/auto_test/invocations/inv_test" {
					t.Errorf("unexpected request: %s %s", r.Method, r.URL)
				}
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(fixture); err != nil {
					t.Error(err)
				}
			}))
			defer server.Close()
			client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("test-key"))
			got, err := client.Beta.Agents.Automations.Invocations.Get(context.Background(), "auto_test", "inv_test")
			if err != nil {
				t.Fatal(err)
			}
			if nulls {
				if got.JSON.ScheduledAt.Raw() != "null" || got.JSON.SessionID.Raw() != "null" || got.JSON.SkipReason.Raw() != "null" || got.JSON.Error.Raw() != "null" {
					t.Fatalf("nullable field lost its JSON null: %#v", got.JSON)
				}
			} else {
				if got.ScheduledAt != 1700000123 {
					t.Errorf("ScheduledAt = %d", got.ScheduledAt)
				}
				if got.SessionID != "sess_preserved" {
					t.Errorf("SessionID = %q", got.SessionID)
				}
				if string(got.SkipReason) != "pending_delivery" {
					t.Errorf("SkipReason = %q", got.SkipReason)
				}
				if got.Error.Code != "synthetic_failure" || got.Error.Message != "distinctive explanation" {
					t.Errorf("Error = %#v", got.Error)
				}
			}
		})
	}
}

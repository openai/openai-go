package openai_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/ssestream"
)

func TestAgentBudgetEventStreamingReplay(t *testing.T) {
	for _, mode := range []string{"create", "subscribe"} {
		t.Run(mode, func(t *testing.T) {
			limitJSON, wantLimit := "500", int64(500)
			if mode == "subscribe" {
				limitJSON, wantLimit = "null", 0
			}
			payloads := []string{
				`{"type":"error","event_id":"before","session_id":"session","error":{"code":"session_budget_exceeded","message":"Synthetic budget reached"}}`,
				fmt.Sprintf(`{"type":"agent.session.budget_reached","event_id":"budget","session":{"id":"session","status":"budget_reached","spend_control":{"limit":%s,"consumed":500}},"obfuscation":"synthetic-padding"}`, limitJSON),
				`{"type":"agent.session.idle","event_id":"after","session":{"id":"session","status":"idle"}}`,
			}
			server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				method, path := http.MethodGet, "/agents/sessions/session/events"
				if mode == "create" {
					method, path = http.MethodPost, "/agents/sessions"
					var body map[string]any
					if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body["stream"] != true {
						t.Errorf("creation body=%v, error=%v", body, err)
					}
				}
				if r.Method != method || r.URL.Path != path {
					t.Errorf("request=%s %s, want %s %s", r.Method, r.URL.Path, method, path)
				}
				w.Header().Set("Content-Type", "text/event-stream")
				for _, payload := range payloads {
					var event struct {
						Type string `json:"type"`
					}
					if err := json.Unmarshal([]byte(payload), &event); err != nil {
						t.Error(err)
						return
					}
					if _, err := fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event.Type, payload); err != nil {
						t.Error(err)
						return
					}
				}
			}))
			defer server.Close()
			client := openai.NewClient(option.WithBaseURL(server.URL), option.WithHTTPClient(server.Client()), option.WithAPIKey("synthetic"), option.WithMaxRetries(0))
			var stream *ssestream.Stream[openai.AgentSessionEventUnion]
			if mode == "create" {
				stream = client.Beta.Agents.Sessions.NewStreaming(t.Context(), openai.BetaAgentSessionNewParams{
					Environment: openai.EnvironmentParamUnion{OfParamNone: &openai.EnvironmentParamNone{}},
				})
			} else {
				stream = client.Beta.Agents.Sessions.Events.StreamStreaming(t.Context(), "session")
			}
			defer func() { _ = stream.Close() }()
			var ids []string
			for stream.Next() {
				event := stream.Current()
				ids = append(ids, event.EventID)
				if event.EventID == "before" {
					failure, ok := event.AsAny().(openai.AgentSessionErrorEvent)
					if !ok || failure.Type != "error" || failure.SessionID != "session" || failure.Error.Code != "session_budget_exceeded" || failure.Error.Message != "Synthetic budget reached" {
						t.Fatalf("error event type=%T payload=%#v", event.AsAny(), failure)
					}
				}
				if event.EventID != "budget" {
					continue
				}
				budget, ok := event.AsAny().(openai.AgentSessionBudgetReachedEvent)
				if !ok {
					t.Fatalf("budget event type=%T", event.AsAny())
				}
				if budget.Type != "agent.session.budget_reached" || budget.EventID != "budget" || budget.Session.ID != "session" || budget.Session.Status != "budget_reached" || budget.Session.SpendControl.Limit != wantLimit || budget.Session.SpendControl.Consumed != 500 || budget.Obfuscation != "synthetic-padding" {
					t.Fatalf("budget payload=%#v", budget)
				}
				limit := budget.Session.SpendControl.JSON.Limit
				if limit.Raw() != limitJSON || limit.Valid() != (limitJSON != "null") {
					t.Fatalf("limit raw=%q valid=%t, want %q valid=%t", limit.Raw(), limit.Valid(), limitJSON, limitJSON != "null")
				}
				if budget.RawJSON() != payloads[1] {
					t.Fatalf("raw payload=%q, want %q", budget.RawJSON(), payloads[1])
				}
			}
			if err := stream.Err(); err != nil {
				t.Fatal(err)
			}
			if want := []string{"before", "budget", "after"}; !reflect.DeepEqual(ids, want) {
				t.Fatalf("event order=%v, want %v", ids, want)
			}
		})
	}
}

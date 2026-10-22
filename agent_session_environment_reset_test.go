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
	"github.com/openai/openai-go/v3/packages/respjson"
)

func TestAgentSessionEnvironmentResetEvent(t *testing.T) {
	turnID := "turn_synthetic"
	tests := []struct {
		name       string
		turnID     *string
		resetCount int64
	}{
		{name: "null_turn_zero_count", turnID: nil, resetCount: 0},
		{name: "associated_turn", turnID: &turnID, resetCount: 7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload, err := json.Marshal(map[string]any{
				"type":           "agent.session.environment.reset",
				"event_id":       "event_synthetic",
				"session_id":     "session_synthetic",
				"environment_id": "environment_synthetic",
				"turn_id":        tt.turnID,
				"reset_count":    tt.resetCount,
			})
			if err != nil {
				t.Fatalf("json.Marshal(reset event %s) error = %v, want nil", tt.name, err)
			}
			wantTurnJSON, err := json.Marshal(tt.turnID)
			if err != nil {
				t.Fatalf("json.Marshal(turn ID %s) error = %v, want nil", tt.name, err)
			}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet || r.URL.Path != "/agents/sessions/session_synthetic/events" {
					t.Errorf("StreamStreaming(%s) request = %s %s, want GET /agents/sessions/session_synthetic/events", tt.name, r.Method, r.URL.Path)
					w.WriteHeader(http.StatusMethodNotAllowed)
					return
				}
				w.Header().Set("Content-Type", "text/event-stream")
				if _, err := fmt.Fprintf(w, "data: %s\n\n", payload); err != nil {
					t.Errorf("write reset event(%s) error = %v, want nil", tt.name, err)
				}
			}))
			t.Cleanup(server.Close)
			client := openai.NewClient(
				option.WithUnsafeAllowHTTP(), option.WithBaseURL(server.URL),
				option.WithAPIKey("synthetic"),
				option.WithMaxRetries(0),
			)
			stream := client.Beta.Agents.Sessions.Events.StreamStreaming(t.Context(), "session_synthetic")
			t.Cleanup(func() {
				if err := stream.Close(); err != nil {
					t.Errorf("StreamStreaming(%s).Close() error = %v, want nil", tt.name, err)
				}
			})
			if !stream.Next() {
				t.Fatalf("StreamStreaming(%s).Next() = false, error = %v, want one event", tt.name, stream.Err())
			}
			event := stream.Current()
			if _, ok := event.AsAny().(openai.AgentSessionEnvironmentResetEvent); !ok {
				t.Fatalf("StreamStreaming(%s).AsAny() type = %T, want AgentSessionEnvironmentResetEvent", tt.name, event.AsAny())
			}
			reset := event.AsAgentSessionEnvironmentReset()
			if reset.EnvironmentID != "environment_synthetic" || reset.EventID != "event_synthetic" || reset.SessionID != "session_synthetic" {
				t.Errorf("StreamStreaming(%s) identifiers = (%q, %q, %q), want synthetic environment, event, and session IDs", tt.name, reset.EnvironmentID, reset.EventID, reset.SessionID)
			}
			if reset.ResetCount != tt.resetCount || !reset.JSON.ResetCount.Valid() {
				t.Errorf("StreamStreaming(%s) reset count = %d, valid = %t, want %d, true", tt.name, reset.ResetCount, reset.JSON.ResetCount.Valid(), tt.resetCount)
			}
			if got := reset.JSON.TurnID.Raw(); got != string(wantTurnJSON) {
				t.Errorf("StreamStreaming(%s) turn ID metadata = %q, want %q", tt.name, got, wantTurnJSON)
			}
			if stream.Next() || stream.Err() != nil {
				t.Errorf("StreamStreaming(%s) after reset error = %v, want clean end of stream", tt.name, stream.Err())
			}
		})
	}
}

func TestAgentSessionEnvironmentLifecycleEvents(t *testing.T) {
	turnID := "turn_synthetic"
	for _, status := range []string{"suspended", "expired"} {
		for _, associatedTurn := range []*string{nil, &turnID} {
			name := status + "/null_turn"
			if associatedTurn != nil {
				name = status + "/associated_turn"
			}
			t.Run(name, func(t *testing.T) {
				eventType := "agent.session.environment." + status
				payload, err := json.Marshal(map[string]any{
					"type":       eventType,
					"event_id":   "event_synthetic",
					"session_id": "session_synthetic",
					"turn_id":    associatedTurn,
					"environment": map[string]any{
						"id":     "environment_synthetic",
						"type":   "openai_hosted",
						"status": status,
						"error":  nil,
					},
				})
				if err != nil {
					t.Fatal(err)
				}
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.Method != http.MethodGet || r.URL.Path != "/agents/sessions/session_synthetic/events" {
						t.Errorf("request = %s %s, want GET session events", r.Method, r.URL.Path)
						w.WriteHeader(http.StatusMethodNotAllowed)
						return
					}
					w.Header().Set("Content-Type", "text/event-stream")
					if _, err := fmt.Fprintf(w, "data: %s\n\n", payload); err != nil {
						t.Errorf("write event: %v", err)
					}
				}))
				t.Cleanup(server.Close)
				client := openai.NewClient(option.WithUnsafeAllowHTTP(), option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"), option.WithMaxRetries(0))
				stream := client.Beta.Agents.Sessions.Events.StreamStreaming(t.Context(), "session_synthetic")
				t.Cleanup(func() {
					if err := stream.Close(); err != nil {
						t.Errorf("Close(): %v", err)
					}
				})
				if !stream.Next() {
					t.Fatalf("Next() = false, error = %v, want one event", stream.Err())
				}
				event := stream.Current()
				var typed any
				var environment openai.AgentSessionEnvironmentState
				var eventID, sessionID, decodedTurnID string
				var turnField, environmentField respjson.Field
				switch status {
				case "suspended":
					v := event.AsAgentSessionEnvironmentSuspended()
					typed, environment = v, v.Environment
					eventID, sessionID, decodedTurnID = v.EventID, v.SessionID, v.TurnID
					turnField, environmentField = v.JSON.TurnID, v.JSON.Environment
				case "expired":
					v := event.AsAgentSessionEnvironmentExpired()
					typed, environment = v, v.Environment
					eventID, sessionID, decodedTurnID = v.EventID, v.SessionID, v.TurnID
					turnField, environmentField = v.JSON.TurnID, v.JSON.Environment
				}
				if event.Type != eventType || eventID != "event_synthetic" || sessionID != "session_synthetic" {
					t.Fatalf("wrong event envelope: %#v", typed)
				}
				if got := event.AsAny(); !reflect.DeepEqual(got, typed) {
					t.Fatalf("AsAny() = %#v (%T), want %#v (%T)", got, got, typed, typed)
				}
				if !environmentField.Valid() || environment.ID != "environment_synthetic" || environment.Type != "openai_hosted" || string(environment.Status) != status {
					t.Fatalf("wrong environment: %#v", environment)
				}
				if environment.JSON.Error.Valid() || environment.JSON.Error.Raw() != "null" {
					t.Fatalf("environment error metadata = %#v, want null", environment.JSON.Error)
				}
				if associatedTurn == nil {
					if decodedTurnID != "" || turnField.Valid() || turnField.Raw() != "null" {
						t.Fatalf("turn = %q, metadata = %#v, want null", decodedTurnID, turnField)
					}
				} else if decodedTurnID != *associatedTurn || !turnField.Valid() || turnField.Raw() != `"turn_synthetic"` {
					t.Fatalf("turn = %q, metadata = %#v, want associated turn", decodedTurnID, turnField)
				}
				if stream.Next() || stream.Err() != nil {
					t.Fatalf("after event error = %v, want clean end of stream", stream.Err())
				}
			})
		}
	}
}

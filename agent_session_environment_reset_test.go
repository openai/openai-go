package openai_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
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
				option.WithBaseURL(server.URL),
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

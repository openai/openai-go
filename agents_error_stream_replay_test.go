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

func TestAgentErrorEventStreamingReplay(t *testing.T) {
	for _, mode := range []string{"create", "subscribe"} {
		t.Run(mode, func(t *testing.T) {
			payloads := []string{
				`{"type":"error","event_id":"before","session_id":"session","error":{"code":"internal_error","message":"Synthetic domain failure"}}`,
				`{"type":"agent.session.future_state","event_id":"future","session":{"id":"session"}}`,
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
					if !ok || failure.Type != "error" || failure.SessionID != "session" || failure.Error.Code != "internal_error" || failure.Error.Message != "Synthetic domain failure" {
						t.Fatalf("error event type=%T payload=%#v", event.AsAny(), failure)
					}
					if failure.RawJSON() != payloads[0] {
						t.Fatalf("raw error payload=%q, want %q", failure.RawJSON(), payloads[0])
					}
				}
			}
			if err := stream.Err(); err != nil {
				t.Fatal(err)
			}
			if want := []string{"before", "future", "after"}; !reflect.DeepEqual(ids, want) {
				t.Fatalf("event order=%v, want %v", ids, want)
			}
		})
	}
}

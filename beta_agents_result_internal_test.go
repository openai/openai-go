package openai

import (
	"context"
	"net/http"
	"net/http/httptest"

	"github.com/openai/openai-go/v3/option"

	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestBetaAgentCollectorRetainsOnlyFinalPayload(t *testing.T) {
	c := &betaAgentTurnCollector{}
	c.enable()
	observe := func(raw string) {
		t.Helper()
		var event AgentSessionEventUnion
		if err := json.Unmarshal([]byte(raw), &event); err != nil {
			t.Fatal(err)
		}
		c.Accumulate(event)
	}
	observe(`{"type":"agent.session.turn.created","turn_id":"root","turn":{"id":"root","session_id":"session","subagent_id":null,"status":"in_progress"}}`)
	payload := strings.Repeat("progress", 100_000)
	for _, phase := range []string{"final_answer", "commentary"} {
		observe(fmt.Sprintf(`{"type":"agent.session.turn.item.added","turn_id":"root","item":{"id":%q,"type":"message","role":"assistant","turn_id":"root","status":"in_progress","phase":%q,"content":[{"type":"output_text","text":%q}]}}`, phase, phase, payload))
	}
	observe(fmt.Sprintf(`{"type":"agent.session.turn.item.done","turn_id":"root","output_index":1,"item":{"id":"commentary","type":"message","role":"assistant","turn_id":"root","status":"completed","phase":"commentary","content":[{"type":"output_text","text":%q}]}}`, payload))
	if len(c.messages) != 0 {
		t.Fatal("retained pending or commentary payload")
	}
	observe(`{"type":"agent.session.turn.item.done","turn_id":"root","output_index":0,"item":{"id":"final_answer","type":"message","role":"assistant","turn_id":"root","status":"completed","phase":"final_answer","content":[{"type":"output_text","text":"answer"}]}}`)
	observe(`{"type":"agent.session.turn.completed","turn_id":"root","turn":{"id":"root","session_id":"session","subagent_id":null,"status":"completed"}}`)
	observe(`{"type":"agent.session.idle"}`)
	result, err := c.finalResult(nil)
	if err != nil || result.OutputText() != "answer" {
		t.Fatalf("result=%v err=%v", result, err)
	}
	if c.messages != nil || c.turn != nil || c.required != nil {
		t.Fatal("collector retained transferred result state")
	}
}

func TestBetaAgentCollectionOptInPublicStream(t *testing.T) {
	const count = 32
	payload := strings.Repeat("answer", 4096)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/agents/sessions/session" {
			w.Header().Set("Content-Type", "application/json")
			_, _ = fmt.Fprint(w, `{"id":"session","status":"idle"}`)
			return
		}
		if r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/events") {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		emit := func(event string) { _, _ = fmt.Fprintf(w, "data: %s\n\n", event) }
		emit(`{"type":"agent.session.turn.created","event_id":"created","turn_id":"root","turn":{"id":"root","session_id":"session","subagent_id":null,"status":"in_progress"}}`)
		for i := 0; i < count; i++ {
			emit(fmt.Sprintf(`{"type":"agent.session.turn.item.done","event_id":"item%d","turn_id":"root","output_index":%d,"item":{"id":"m%d","turn_id":"root","type":"message","role":"assistant","status":"completed","phase":"final_answer","content":[{"type":"output_text","text":%q}]}}`, i, i, i, payload))
		}
		emit(`{"type":"agent.session.turn.completed","event_id":"completed","turn_id":"root","turn":{"id":"root","session_id":"session","status":"completed"}}`)
		emit(`{"type":"agent.session.idle","event_id":"idle"}`)
	}))
	defer server.Close()
	client := NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"), option.WithMaxRetries(0))
	for _, entrypoint := range []string{"creation", "followup"} {
		for _, mode := range []string{"raw", "explicit", "fresh_getter"} {
			t.Run(entrypoint+"/"+mode, func(t *testing.T) {
				var events interface {
					Next() bool
					Err() error
					Close() error
				}
				var collector *betaAgentTurnCollector
				var collect func() (*BetaAgentTurnResult, error)
				if entrypoint == "creation" {
					stream := client.Beta.Agents.Sessions.NewStreaming(context.Background(), BetaAgentSessionNewParams{})
					if mode == "explicit" && BetaAgentSessionWithResultCollection(stream) != stream {
						t.Fatal("different stream")
					}
					events, collector = stream, stream.BetaAccumulator().(*betaAgentTurnCollector)
					collect = func() (*BetaAgentTurnResult, error) { return BetaAgentSessionFinalResult(stream) }
				} else {
					stream := client.Beta.Agents.Sessions.Stream(context.Background(), "session", AgentSessionStreamParams{Input: "question"})
					if mode == "explicit" && stream.WithResultCollection() != stream {
						t.Fatal("different stream")
					}
					events, collector, collect = stream, &stream.collector, stream.FinalResult
				}
				defer func() { _ = events.Close() }()
				if mode != "fresh_getter" {
					for events.Next() {
					}
					if events.Err() != nil {
						t.Fatal(events.Err())
					}
				}
				if mode == "raw" {
					if !collector.started || collector.enabled || collector.turn != nil || collector.messages != nil || collector.required != nil || collector.sessionID != "" {
						t.Fatal("raw iteration retained result state")
					}
					return
				}
				result, err := collect()
				if err != nil || len(result.Messages) != count || len(result.OutputText()) != count*len(payload) {
					t.Fatalf("result collection failed: %v", err)
				}
			})
		}
	}
}

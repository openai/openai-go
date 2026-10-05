package openai_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/ssestream"
)

func TestBetaAgentCreationTools(t *testing.T) {
	for _, mode := range []string{"result", "iteration", "handler-error", "close", "submit-error"} {
		t.Run(mode, func(t *testing.T) {
			var calls, posts atomic.Int32
			posted := make(chan struct{})
			server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/agents/sessions":
					var body map[string]any
					if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
						t.Error(err)
					}
					if r.Method != "POST" || body["stream"] != true || body["input"] != "Look up A123" || strings.Contains(fmt.Sprint(body), "ToolHandlers") {
						t.Errorf("unexpected creation: %s %#v", r.Method, body)
					}
					if r.Header.Get("Idempotency-Key") != "creation-key" {
						t.Error("lost creation key")
					}
					w.Header().Set("Content-Type", "text/event-stream")
					events := []string{
						agentEvent("created", "created", `,"session":{"id":"session","status":"in_progress"}`),
						betaResultTurn("created", "root", "null"),
						agentCall("call-event", "root", "call", "lookup", `{"id":"A123"}`),
						agentCall("duplicate", "root", "call", "lookup", `{"id":"A123"}`),
						agentEvent("requires_action", "pending", `,"session":{"id":"session","required_actions":[{"type":"function_call","turn_id":"root","call_id":"call","name":"lookup","arguments":{}}]}`),
					}
					for _, event := range events {
						_, _ = fmt.Fprintf(w, "data: %s\n\n", event)
						w.(http.Flusher).Flush()
					}
					select {
					case <-r.Context().Done():
						return
					case <-posted:
					}
					for _, event := range []string{betaResultMessage("done", "m", "root", `"final_answer"`, "ready", 0), betaResultTurn("completed", "root", "null"), betaResultIdle()} {
						_, _ = fmt.Fprintf(w, "data: %s\n\n", event)
						w.(http.Flusher).Flush()
					}
					<-r.Context().Done()
				case "/agents/sessions/session/events":
					posts.Add(1)
					if r.Header.Get("Idempotency-Key") == "" || r.Header.Get("Idempotency-Key") == "creation-key" {
						t.Error("tool result requires distinct key")
					}
					var body map[string]any
					if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
						t.Error(err)
					}
					if body["stream"] != nil || body["input"] != nil {
						t.Error("creation params leaked into result request")
					}
					event := body["events"].([]any)[0].(map[string]any)
					if event["call_id"] != "call" || event["turn_id"] != "root" || event["success"] != (mode != "handler-error") {
						t.Errorf("unexpected tool result: %#v", event)
					}
					if strings.Contains(fmt.Sprint(body), "private-error") {
						t.Error("handler error leaked")
					}
					if mode == "submit-error" {
						w.WriteHeader(503)
						return
					}
					w.WriteHeader(204)
					close(posted)
				default:
					t.Errorf("unexpected request: %s", r.URL.Path)
					w.WriteHeader(404)
				}
			}))
			defer server.Close()
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			client := openai.NewClient(option.WithBaseURL(server.URL), option.WithHTTPClient(server.Client()), option.WithAPIKey("synthetic"), option.WithMaxRetries(0))
			handlers := map[string]openai.AgentToolHandler{"lookup": func(_ context.Context, args map[string]any) (any, error) {
				calls.Add(1)
				if args["id"] != "A123" {
					t.Errorf("arguments mutated: %#v", args)
				}
				if mode == "handler-error" {
					return nil, errors.New("private-error")
				}
				return "ready", nil
			}}
			// The concrete stream and existing result helpers remain source compatible.
			var stream *ssestream.Stream[openai.AgentSessionEventUnion] = client.Beta.Agents.Sessions.NewStreaming(ctx, openai.BetaAgentSessionNewParams{
				Environment: openai.EnvironmentParamUnion{OfParamNone: &openai.EnvironmentParamNone{}},
				Input:       openai.BetaAgentSessionNewParamsInputUnion{OfString: openai.String("Look up A123")}, ToolHandlers: handlers,
			}, option.WithHeader("Idempotency-Key", "creation-key"))
			defer func() { _ = stream.Close() }()
			delete(handlers, "lookup")
			if mode == "iteration" || mode == "close" {
				openai.BetaAgentSessionWithResultCollection(stream)
				for stream.Next() {
					event := stream.Current()
					if event.Type == "agent.session.turn.item.added" {
						event.Item.Arguments.(map[string]any)["id"] = "mutated"
						if mode == "close" {
							_ = stream.Close()
							break
						}
					}
				}
			}
			if mode == "close" {
				if stream.Next() || calls.Load() != 0 || posts.Load() != 0 {
					t.Fatal("Close must discard pending callbacks")
				}
				return
			}
			result, err := openai.BetaAgentSessionFinalResult(stream)
			if mode == "submit-error" {
				var apiErr *openai.Error
				if !errors.As(err, &apiErr) || apiErr.StatusCode != 503 || stream.Err() == nil {
					t.Fatalf("lost submission error: %v", err)
				}
				return
			}
			if err != nil || result.OutputText() != "ready" {
				t.Fatalf("result=%v error=%v", result, err)
			}
			cached, err := openai.BetaAgentSessionFinalResult(stream)
			if err != nil || cached != result || calls.Load() != 1 || posts.Load() != 1 {
				t.Fatal("expected one callback/post and cached result")
			}
		})
	}
}

func TestBetaAgentCreationToolHandlersAreLocal(t *testing.T) {
	params := openai.BetaAgentSessionNewParams{ToolHandlers: map[string]openai.AgentToolHandler{
		"lookup": func(context.Context, map[string]any) (any, error) { return nil, nil },
	}}
	body, err := json.Marshal(params)
	if err != nil || strings.Contains(string(body), "lookup") || strings.Contains(string(body), "ToolHandlers") {
		t.Fatalf("handlers serialized: %s %v", body, err)
	}
	client := openai.NewClient(option.WithAPIKey("synthetic"))
	if _, err := client.Beta.Agents.Sessions.New(context.Background(), params); err == nil {
		t.Fatal("non-streaming creation must not silently ignore handlers")
	}
}

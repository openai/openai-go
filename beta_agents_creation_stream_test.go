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
	for _, mode := range []string{"result", "iteration", "handler-error", "close", "submit-error", "json-options", "body-options", "inherited-options", "events-options", "events-constructor-options", "events-wrapped-options", "session-options", "session-constructor-options"} {
		t.Run(mode, func(t *testing.T) {
			var calls, posts atomic.Int32
			posted := make(chan struct{})
			server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("X-Application") != "test" || r.Header.Get("Authorization") != "Bearer synthetic" {
					t.Error("lost request header/auth options")
				}
				if strings.HasPrefix(mode, "session-") && r.Header.Get("X-Session") != "kept" {
					t.Error("lost session-scoped header")
				}
				switch r.URL.Path {
				case "/agents/sessions":
					var body map[string]any
					if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
						t.Error(err)
					}
					if r.Method != "POST" || body["event_marker"] != nil || body["stream"] != true || body["input"] != "Look up A123" || strings.Contains(fmt.Sprint(body), "ToolHandlers") {
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
					if r.Header.Get("Content-Type") != "application/json" {
						t.Error("creation content type leaked into tool result")
					}
					if r.Header.Get("Idempotency-Key") == "" || r.Header.Get("Idempotency-Key") == "creation-key" || r.Header.Get("Idempotency-Key") == "event-key" {
						t.Error("tool result requires distinct key")
					}
					var body map[string]any
					if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
						t.Error(err)
					}
					if (body["event_marker"] == "yes") != (strings.HasPrefix(mode, "events-")) {
						t.Error("lost Events-scoped body option")
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
			var creationResponse, eventResponse *http.Response
			input := "Look up A123"
			opts := []option.RequestOption{option.WithHeader("Idempotency-Key", "creation-key"), option.WithHeader("X-Application", "test")}
			if strings.HasSuffix(mode, "-options") {
				input = "This should be replaced by the creation body option"
				creationOpts := []option.RequestOption{option.WithResponseInto(&creationResponse), option.WithJSONSet("input", "Look up A123"), option.WithJSONDel("events")}
				if mode == "body-options" {
					creationOpts = []option.RequestOption{option.WithResponseInto(&creationResponse), option.WithRequestBody("application/vnd.test+json", []byte(`{"environment":{"type":"none"},"input":"Look up A123"}`))}
				}
				if mode == "inherited-options" || mode == "events-wrapped-options" {
					client = openai.NewClient(append([]option.RequestOption{option.WithBaseURL(server.URL), option.WithHTTPClient(server.Client()), option.WithAPIKey("synthetic"), option.WithMaxRetries(0)}, creationOpts...)...)
				} else if strings.HasPrefix(mode, "session-") {
					client = openai.NewClient(option.WithBaseURL(server.URL+"/unused"), option.WithAPIKey("original-key"))
					sessionOpts := append([]option.RequestOption{option.WithBaseURL(server.URL), option.WithHTTPClient(server.Client()), option.WithAPIKey("synthetic"), option.WithHeader("X-Session", "kept")}, creationOpts...)
					if mode == "session-constructor-options" {
						sessionOpts = openai.NewBetaAgentSessionService(sessionOpts...).Options
					}
					client.Beta.Agents.Sessions.Options = append(client.Beta.Agents.Sessions.Options, sessionOpts...)
				} else {
					opts = append(opts, creationOpts...)
				}
			}
			if strings.HasPrefix(mode, "events-") {
				eventOpts := []option.RequestOption{option.WithResponseInto(&eventResponse), option.WithJSONSet("event_marker", "yes"), option.WithHeader("Idempotency-Key", "event-key")}
				switch mode {
				case "events-constructor-options":
					defaults := []option.RequestOption{option.WithBaseURL(server.URL), option.WithHTTPClient(server.Client()), option.WithAPIKey("synthetic"), option.WithMaxRetries(0)}
					client.Beta.Agents.Sessions.Events = openai.NewBetaAgentSessionEventService(append(defaults, eventOpts...)...)
				case "events-wrapped-options":
					client.Beta.Agents.Sessions.Events = openai.NewBetaAgentSessionEventService(append(client.Beta.Agents.Sessions.Events.Options, eventOpts...)...)
				default:
					client.Beta.Agents.Sessions.Events.Options = append(client.Beta.Agents.Sessions.Events.Options, eventOpts...)
				}
			}
			// The concrete stream and existing result helpers remain source compatible.
			var stream *ssestream.Stream[openai.AgentSessionEventUnion] = client.Beta.Agents.Sessions.NewStreaming(ctx, openai.BetaAgentSessionNewParams{
				Environment: openai.EnvironmentParamUnion{OfParamNone: &openai.EnvironmentParamNone{}},
				Input:       openai.BetaAgentSessionNewParamsInputUnion{OfString: openai.String(input)}, ToolHandlers: handlers,
			}, opts...)
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
			if strings.HasSuffix(mode, "-options") && (creationResponse == nil || creationResponse.Request.URL.Path != "/agents/sessions") {
				t.Fatal("tool result overwrote creation response capture")
			}
			if strings.HasPrefix(mode, "events-") && (eventResponse == nil || eventResponse.Request.URL.Path != "/agents/sessions/session/events") {
				t.Fatal("lost Events-scoped response capture")
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

func TestBetaAgentCreationToolsRejectResponseBodyOverrides(t *testing.T) {
	for _, inherited := range []bool{false, true} {
		for _, destination := range []any{new([]byte), new(*http.Response)} {
			t.Run(fmt.Sprintf("inherited=%v/%T", inherited, destination), func(t *testing.T) {
				server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					t.Error("response override must be rejected before transport")
					w.WriteHeader(500)
				}))
				defer server.Close()
				client := openai.NewClient(option.WithBaseURL(server.URL), option.WithHTTPClient(server.Client()), option.WithAPIKey("synthetic"))
				opts := []option.RequestOption{option.WithResponseBodyInto(destination)}
				if inherited {
					client.Beta.Agents.Sessions.Options = append(client.Beta.Agents.Sessions.Options, opts...)
					opts = nil
				}
				stream := client.Beta.Agents.Sessions.NewStreaming(context.Background(), openai.BetaAgentSessionNewParams{
					ToolHandlers: map[string]openai.AgentToolHandler{"lookup": func(context.Context, map[string]any) (any, error) {
						t.Error("handler must not execute")
						return nil, nil
					}},
				}, opts...)
				defer func() { _ = stream.Close() }()
				if stream.Err() == nil || !strings.Contains(stream.Err().Error(), "WithResponseBodyInto") || stream.Next() {
					t.Fatalf("expected response override error, got %v", stream.Err())
				}
			})
		}
	}
}

func TestBetaAgentCreationToolsRejectReorderedOptions(t *testing.T) {
	for _, service := range []string{"sessions", "events", "wrapped-events"} {
		t.Run(service, func(t *testing.T) {
			server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				t.Error("invalid option order must be rejected before transport")
				w.WriteHeader(500)
			}))
			defer server.Close()
			client := openai.NewClient(option.WithBaseURL(server.URL), option.WithHTTPClient(server.Client()), option.WithAPIKey("synthetic"))
			prefix := []option.RequestOption{option.WithBaseURL(server.URL + "/other")}
			switch service {
			case "sessions":
				client.Beta.Agents.Sessions.Options = append(prefix, client.Beta.Agents.Sessions.Options...)
			case "events":
				client.Beta.Agents.Sessions.Events.Options = append(prefix, client.Beta.Agents.Sessions.Events.Options...)
			case "wrapped-events":
				client.Beta.Agents.Sessions.Events = openai.NewBetaAgentSessionEventService(append(prefix, client.Beta.Agents.Sessions.Events.Options...)...)
			}
			stream := client.Beta.Agents.Sessions.NewStreaming(context.Background(), openai.BetaAgentSessionNewParams{
				ToolHandlers: map[string]openai.AgentToolHandler{"lookup": func(context.Context, map[string]any) (any, error) {
					t.Error("handler must not execute")
					return nil, nil
				}},
			})
			defer func() { _ = stream.Close() }()
			if stream.Err() == nil || !strings.Contains(stream.Err().Error(), "append options") || stream.Next() {
				t.Fatalf("expected option-order error, got %v", stream.Err())
			}
		})
	}
}

func TestBetaAgentCreationToolsWithoutInput(t *testing.T) {
	for _, done := range []bool{false, true} {
		t.Run(fmt.Sprintf("done-frame=%v", done), func(t *testing.T) {
			server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/event-stream")
				for _, event := range []string{agentEvent("created", "created", `,"session":{"id":"session","status":"idle"}`), betaResultIdle()} {
					_, _ = fmt.Fprintf(w, "data: %s\n\n", event)
				}
				if done {
					_, _ = fmt.Fprint(w, "data: [DONE]\n\n")
				}
			}))
			defer server.Close()
			client := openai.NewClient(option.WithBaseURL(server.URL), option.WithHTTPClient(server.Client()), option.WithAPIKey("synthetic"))
			stream := client.Beta.Agents.Sessions.NewStreaming(context.Background(), openai.BetaAgentSessionNewParams{
				Environment: openai.EnvironmentParamOfParamSelfHosted("/workspace"),
				ToolHandlers: map[string]openai.AgentToolHandler{"lookup": func(context.Context, map[string]any) (any, error) {
					t.Error("no-input creation must not call tools")
					return nil, nil
				}},
			})
			defer func() { _ = stream.Close() }()
			openai.BetaAgentSessionWithResultCollection(stream)
			events := 0
			for stream.Next() {
				events++
			}
			if stream.Err() != nil || events != 2 {
				t.Fatalf("raw completion: events=%d err=%v", events, stream.Err())
			}
			result, err := openai.BetaAgentSessionFinalResult(stream)
			var incomplete *openai.BetaAgentTurnResultError
			if result != nil || !errors.As(err, &incomplete) || incomplete.Reason != "observation_incomplete" {
				t.Fatalf("no turn must not produce a final result: result=%v err=%v", result, err)
			}
		})
	}
}

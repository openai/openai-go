package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type recordingCatalog struct {
	args  []lookupArguments
	cause error
}

func TestRunReportsStreamFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()
	client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"), option.WithMaxRetries(0))
	if err := run(context.Background(), client, "session"); err == nil {
		t.Fatal("failed stream reported a successful run")
	}
}

func TestRunReportsFailedEvents(t *testing.T) {
	for _, tc := range []struct {
		name, terminal string
		wantError      bool
	}{
		{"session failed", `{"type":"agent.session.failed","event_id":"failed","session":{"status":"failed"}}`, true},
		{"turn failed", `{"type":"agent.session.turn.failed","event_id":"failed","turn_id":"turn","turn":{"id":"turn","subagent_id":null}}`, true},
		{"turn cancelled", `{"type":"agent.session.turn.cancelled","event_id":"cancelled","turn_id":"turn","turn":{"id":"turn","subagent_id":null}}`, true},
		{"child failed", `{"type":"agent.session.turn.failed","event_id":"failed","turn_id":"child","turn":{"id":"child","subagent_id":"subagent"}}`, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.Method + " " + r.URL.Path {
				case "GET /agents/sessions/session":
					w.Header().Set("Content-Type", "application/json")
					_, _ = fmt.Fprint(w, `{"id":"session","status":"idle"}`)
				case "GET /agents/sessions/session/events":
					w.Header().Set("Content-Type", "text/event-stream")
					for _, event := range []string{
						`{"type":"agent.session.turn.created","event_id":"created","turn_id":"turn","turn":{"id":"turn","subagent_id":null}}`,
						tc.terminal,
						`{"type":"agent.session.turn.completed","event_id":"done","turn_id":"turn"}`,
						`{"type":"agent.session.idle","event_id":"idle","session":{"status":"idle"}}`,
					} {
						_, _ = fmt.Fprintf(w, "data: %s\n\n", event)
					}
				case "POST /agents/sessions/session/events":
					w.WriteHeader(http.StatusNoContent)
				default:
					t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
					w.WriteHeader(http.StatusNotFound)
				}
			}))
			defer server.Close()
			client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"), option.WithMaxRetries(0))
			if err := run(context.Background(), client, "session"); (err != nil) != tc.wantError {
				t.Fatalf("run error = %v, want error: %t", err, tc.wantError)
			}
		})
	}
}

func (c *recordingCatalog) Lookup(_ context.Context, args lookupArguments) (itemRecord, error) {
	c.args = append(c.args, args)
	return itemRecord{Catalog: "bound-catalog", ItemID: args.ItemID, Availability: "available"}, c.cause
}

func TestLookupActionRejectsInvalidArguments(t *testing.T) {
	for name, args := range map[string]map[string]any{
		"nil": nil, "missing": {}, "null": {"item_id": nil},
		"wrong type": {"item_id": 123}, "unsupported item": {"item_id": "unknown-item"},
		"unexpected field": {"item_id": "item-1", "catalog": "other-catalog"},
		"case mismatch":    {"ItemID": "item-1"},
	} {
		t.Run(name, func(t *testing.T) {
			catalog := &recordingCatalog{}
			_, err := (lookupAction{catalog: catalog}).handle(context.Background(), args)
			if err == nil || len(catalog.args) != 0 {
				t.Fatalf("invalid arguments reached the catalog: calls=%d, error=%v", len(catalog.args), err)
			}
		})
	}
}

func TestLookupActionStream(t *testing.T) {
	for _, fail := range []bool{false, true} {
		t.Run(fmt.Sprintf("provider_error=%t", fail), func(t *testing.T) {
			catalog := &recordingCatalog{}
			if fail {
				catalog.cause = errors.New("private provider detail")
			}
			action := lookupAction{catalog: catalog}
			tool, err := action.tool()
			if err != nil {
				t.Fatal(err)
			}
			definition := tool.OfParamFunction
			schema := definition.Parameters
			properties := schema["properties"].(map[string]any)
			if schema["type"] != "object" || schema["additionalProperties"] != false || len(properties) != 1 ||
				!reflect.DeepEqual(schema["required"], []any{"item_id"}) ||
				!reflect.DeepEqual(properties["item_id"].(map[string]any)["enum"], []any{"item-1", "item-2"}) {
				t.Fatalf("schema does not describe the typed arguments: %v", schema)
			}

			posts := make(chan map[string]any, 2)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.Method + " " + r.URL.Path {
				case "GET /agents/sessions/session":
					w.Header().Set("Content-Type", "application/json")
					_, _ = fmt.Fprint(w, `{"id":"session","status":"idle"}`)
				case "GET /agents/sessions/session/events":
					w.Header().Set("Content-Type", "text/event-stream")
					for _, event := range []string{
						`{"type":"agent.session.turn.created","event_id":"created","turn_id":"turn","turn":{"id":"turn","subagent_id":null}}`,
						`{"type":"agent.session.turn.item.added","event_id":"call","turn_id":"turn","item":{"id":"item","type":"function_call","call_id":"call","name":"catalog_lookup","arguments":{"item_id":"item-1"},"status":"in_progress"}}`,
						`{"type":"agent.session.turn.completed","event_id":"done","turn_id":"turn"}`,
						`{"type":"agent.session.idle","event_id":"idle","session":{"status":"idle"}}`,
					} {
						_, _ = fmt.Fprintf(w, "data: %s\n\n", event)
					}
				case "POST /agents/sessions/session/events":
					var body map[string]any
					if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
						t.Errorf("decode request: %v", err)
					}
					posts <- body
					w.WriteHeader(http.StatusNoContent)
				default:
					t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
					w.WriteHeader(http.StatusNotFound)
				}
			}))
			defer server.Close()
			client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"), option.WithMaxRetries(0))
			stream := client.Beta.Agents.Sessions.Stream(context.Background(), "session", openai.AgentSessionStreamParams{
				Input:        "Look up item-1 in the catalog.",
				ToolHandlers: map[string]openai.AgentToolHandler{definition.Name: action.handle},
			})
			defer func() { _ = stream.Close() }()
			for stream.Next() {
			}
			if err := stream.Err(); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(catalog.args, []lookupArguments{{ItemID: "item-1"}}) || len(posts) != 2 {
				t.Fatalf("unexpected dispatch: calls=%v, posts=%d", catalog.args, len(posts))
			}
			<-posts // The user input.
			result := (<-posts)["events"].([]any)[0].(map[string]any)
			if fail {
				if result["error"] != "Tool handler failed." || result["success"] == true {
					t.Fatalf("provider failure not safely reported: %v", result)
				}
			} else {
				var record itemRecord
				if err := json.Unmarshal([]byte(result["output"].(string)), &record); err != nil {
					t.Fatal(err)
				}
				if result["success"] != true || record != (itemRecord{Catalog: "bound-catalog", ItemID: "item-1", Availability: "available"}) {
					t.Fatalf("unexpected typed record: %v", result)
				}
			}
		})
	}
}

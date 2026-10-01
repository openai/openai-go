package openai_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func TestBetaAgentAttachPendingAndDurableOutput(t *testing.T) {
	var posts, pages atomic.Int32
	events := []string{agentCall("recovered", "root", "pending", "lookup", `{"id":"A123"}`), agentCall("overlap", "root", "pending", "lookup", `{"id":"A123"}`), betaResultTurn("completed", "root", "null"), betaResultIdle()}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Attachment") != "kept" {
			t.Error("request options lost")
		}
		switch {
		case strings.HasSuffix(r.URL.Path, "/events") && r.Method == http.MethodPost:
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Error(err)
			}
			input := body["events"].([]any)[0].(map[string]any)
			if input["type"] != "agent.session.input.tool_result" {
				t.Errorf("attachment submitted input: %#v", body)
			}
			posts.Add(1)
			w.WriteHeader(http.StatusNoContent)
		case strings.HasSuffix(r.URL.Path, "/events"):
			w.Header().Set("Content-Type", "text/event-stream")
			for _, event := range events {
				_, _ = fmt.Fprintf(w, "data: %s\n\n", event)
			}
		case strings.HasSuffix(r.URL.Path, "/turns"):
			w.Header().Set("Content-Type", "application/json")
			_, _ = fmt.Fprint(w, `{"data":[{"id":"root","session_id":"session","status":"waiting","subagent_id":null}],"has_more":false}`)
		case strings.HasSuffix(r.URL.Path, "/turns/root"):
			w.Header().Set("Content-Type", "application/json")
			status := "waiting"
			if posts.Load() > 0 {
				status = "completed"
			}
			_, _ = fmt.Fprintf(w, `{"id":"root","session_id":"session","status":%q,"subagent_id":null}`, status)
		case strings.HasSuffix(r.URL.Path, "/items"):
			w.Header().Set("Content-Type", "application/json")
			pages.Add(1)
			if r.URL.Query().Get("order") != "asc" {
				t.Error("wrong durable message ordering")
			}
			if r.URL.Query().Get("after") == "" {
				_, _ = fmt.Fprint(w, `{"data":[{"id":"old","type":"message","turn_id":"older","role":"assistant","status":"completed","content":[{"type":"output_text","text":"wrong"}]},{"id":"m1","type":"message","turn_id":"root","role":"assistant","status":"completed","content":[{"type":"output_text","text":"earlier "}]}],"has_more":true,"last_id":"m1"}`)
			} else {
				_, _ = fmt.Fprint(w, `{"data":[{"id":"m2","type":"message","turn_id":"root","role":"assistant","status":"completed","phase":"final_answer","content":[{"type":"output_text","text":"answer"}]}],"has_more":false}`)
			}
		default:
			w.Header().Set("Content-Type", "application/json")
			// The stale snapshot action is not an execution queue.
			_, _ = fmt.Fprint(w, `{"id":"session","status":"requires_action","required_actions":[{"type":"function_call","turn_id":"root","call_id":"answered","name":"lookup","arguments":{}}]}`)
		}
	}))
	defer server.Close()
	client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"), option.WithMaxRetries(0))
	calls := 0
	stream := client.Beta.Agents.Sessions.Stream(context.Background(), "session", openai.AgentSessionStreamParams{ToolHandlers: map[string]openai.AgentToolHandler{"lookup": func(context.Context, map[string]any) (any, error) { calls++; return "found", nil }}}, option.WithHeader("X-Attachment", "kept"))
	result, err := stream.FinalResult()
	if err != nil || result.OutputText() != "earlier answer" || result.TurnID() != "root" {
		t.Fatalf("result=%v err=%v", result, err)
	}
	if calls != 1 || posts.Load() != 1 || pages.Load() != 2 {
		t.Fatalf("calls=%d posts=%d pages=%d", calls, posts.Load(), pages.Load())
	}
	again, err := stream.FinalResult()
	if err != nil || again != result || pages.Load() != 2 {
		t.Fatal("getter was not cached", err)
	}
}

func TestBetaAgentAttachAlreadyIdle(t *testing.T) {
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		if r.Method != http.MethodGet {
			t.Errorf("unexpected idle request %s %s", r.Method, r.URL.Path)
		}
		if strings.HasSuffix(r.URL.Path, "/events") {
			w.Header().Set("Content-Type", "text/event-stream")
			w.WriteHeader(200)
			w.(http.Flusher).Flush()
			<-r.Context().Done()
			return
		}
		if strings.HasSuffix(r.URL.Path, "/turns") {
			w.Header().Set("Content-Type", "application/json")
			_, _ = fmt.Fprint(w, `{"data":[{"id":"previous","session_id":"session","status":"completed"}],"has_more":false}`)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"id":"session","status":"idle"}`)
	}))
	defer server.Close()
	client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"))
	stream := client.Beta.Agents.Sessions.Stream(context.Background(), "session", openai.AgentSessionStreamParams{})
	if stream.Next() || stream.Err() != nil {
		t.Fatal("idle attach should settle", stream.Err())
	}
	_, err := stream.FinalResult()
	var failure *openai.BetaAgentTurnResultError
	if !errors.As(err, &failure) || failure.Reason != "no_turn_selected" {
		t.Fatal("idle attach selected historical answer", err)
	}
	if requests.Load() != 5 {
		t.Fatal("idle attach did unnecessary history requests")
	}
}

func TestBetaAgentAttachCompletionGap(t *testing.T) {
	var sessionReads atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case strings.HasSuffix(r.URL.Path, "/events"):
			w.Header().Set("Content-Type", "text/event-stream")
			w.WriteHeader(http.StatusOK)
			w.(http.Flusher).Flush()
			<-r.Context().Done()
		case strings.HasSuffix(r.URL.Path, "/turns"):
			_, _ = fmt.Fprint(w, `{"data":[{"id":"root","session_id":"session","status":"in_progress"}],"has_more":false}`)
		case strings.HasSuffix(r.URL.Path, "/turns/root"):
			_, _ = fmt.Fprint(w, `{"id":"root","session_id":"session","status":"completed"}`)
		case strings.HasSuffix(r.URL.Path, "/items"):
			_, _ = fmt.Fprint(w, `{"data":[{"id":"m","turn_id":"root","type":"message","role":"assistant","status":"completed","content":[{"type":"output_text","text":"finished"}]}],"has_more":false}`)
		default:
			status := "in_progress"
			if sessionReads.Add(1) > 1 {
				status = "idle"
			}
			_, _ = fmt.Fprintf(w, `{"id":"session","status":%q}`, status)
		}
	}))
	defer server.Close()
	client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	stream := client.Beta.Agents.Sessions.Stream(ctx, "session", openai.AgentSessionStreamParams{})
	result, err := stream.FinalResult()
	if err != nil || result.OutputText() != "finished" {
		t.Fatalf("completion gap result=%v err=%v", result, err)
	}
}

func TestBetaAgentAttachDoesNotReplayAcceptedResultAfterDisconnect(t *testing.T) {
	var connections, posts atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case strings.HasSuffix(r.URL.Path, "/events") && r.Method == http.MethodPost:
			posts.Add(1)
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = fmt.Fprint(w, `{"error":{"message":"synthetic response loss","type":"server_error"}}`)
		case strings.HasSuffix(r.URL.Path, "/events"):
			w.Header().Set("Content-Type", "text/event-stream")
			if connections.Add(1) == 1 {
				_, _ = fmt.Fprintf(w, "data: %s\n\n", agentCall("pending", "root", "call", "lookup", `{}`))
			} else {
				_, _ = fmt.Fprintf(w, "data: %s\n\ndata: %s\n\n", betaResultTurn("completed", "root", "null"), betaResultIdle())
			}
		case strings.HasSuffix(r.URL.Path, "/turns"):
			_, _ = fmt.Fprint(w, `{"data":[{"id":"root","session_id":"session","status":"waiting"}],"has_more":false}`)
		case strings.HasSuffix(r.URL.Path, "/turns/root"):
			status := "in_progress"
			if connections.Load() > 1 {
				status = "completed"
			}
			_, _ = fmt.Fprintf(w, `{"id":"root","session_id":"session","status":%q}`, status)
		case strings.HasSuffix(r.URL.Path, "/items"):
			_, _ = fmt.Fprint(w, `{"data":[],"has_more":false}`)
		default:
			_, _ = fmt.Fprint(w, `{"id":"session","status":"requires_action","required_actions":[{"type":"function_call","turn_id":"root","call_id":"call","name":"lookup","arguments":{}}]}`)
		}
	}))
	defer server.Close()
	client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"), option.WithMaxRetries(0))
	calls := 0
	params := openai.AgentSessionStreamParams{ToolHandlers: map[string]openai.AgentToolHandler{"lookup": func(context.Context, map[string]any) (any, error) { calls++; return "saved receipt", nil }}}
	first := client.Beta.Agents.Sessions.Stream(context.Background(), "session", params)
	_, err := first.FinalResult()
	var failure *openai.BetaAgentTurnResultError
	if !errors.As(err, &failure) || failure.Reason != "observation_failed" {
		t.Fatal("lost tool-result transport failure", err)
	}
	second := client.Beta.Agents.Sessions.Stream(context.Background(), "session", params)
	result, err := second.FinalResult()
	if err != nil || result.TurnID() != "root" || calls != 1 || posts.Load() != 1 {
		t.Fatalf("retried side effect: result=%v err=%v calls=%d posts=%d", result, err, calls, posts.Load())
	}
}

func TestBetaAgentAttachUnhandledFirstFrame(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case strings.HasSuffix(r.URL.Path, "/events"):
			w.Header().Set("Content-Type", "text/event-stream")
			_, _ = fmt.Fprintf(w, "data: %s\n\n", agentCall("pending", "root", "call", "unknown", `{"key":"value"}`))
			w.(http.Flusher).Flush()
			<-r.Context().Done()
		case strings.HasSuffix(r.URL.Path, "/turns"):
			_, _ = fmt.Fprint(w, `{"data":[{"id":"root","session_id":"session","status":"waiting"}],"has_more":false}`)
		case strings.HasSuffix(r.URL.Path, "/turns/root"):
			_, _ = fmt.Fprint(w, `{"id":"root","session_id":"session","status":"waiting"}`)
		default:
			_, _ = fmt.Fprint(w, `{"id":"session","status":"requires_action"}`)
		}
	}))
	defer server.Close()
	client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := client.Beta.Agents.Sessions.Stream(ctx, "session", openai.AgentSessionStreamParams{}).FinalResult()
	var failure *openai.BetaAgentTurnResultError
	if !errors.As(err, &failure) || failure.Reason != "requires_action" || len(failure.RequiredActions) != 1 || failure.RequiredActions[0].CallID != "call" {
		t.Fatal("missing unhandled recovered action", err)
	}
}

func TestBetaAgentAttachSelectsNewRootDuringHandshake(t *testing.T) {
	var lists atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case strings.HasSuffix(r.URL.Path, "/events"):
			w.Header().Set("Content-Type", "text/event-stream")
			w.WriteHeader(200)
			w.(http.Flusher).Flush()
			<-r.Context().Done()
		case strings.HasSuffix(r.URL.Path, "/turns"):
			id := "previous"
			if lists.Add(1) > 1 {
				id = "new"
			}
			_, _ = fmt.Fprintf(w, `{"data":[{"id":%q,"session_id":"session","status":"completed"}],"has_more":false}`, id)
		case strings.HasSuffix(r.URL.Path, "/turns/new"):
			_, _ = fmt.Fprint(w, `{"id":"new","session_id":"session","status":"completed"}`)
		case strings.HasSuffix(r.URL.Path, "/items"):
			_, _ = fmt.Fprint(w, `{"data":[],"has_more":false}`)
		default:
			// A successor can keep the session active after the selected turn finishes.
			_, _ = fmt.Fprint(w, `{"id":"session","status":"in_progress"}`)
		}
	}))
	defer server.Close()
	client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	result, err := client.Beta.Agents.Sessions.Stream(ctx, "session", openai.AgentSessionStreamParams{}).FinalResult()
	if err != nil || result.TurnID() != "new" {
		t.Fatalf("missed handshake root result=%v err=%v", result, err)
	}
}

func TestBetaAgentAttachReadFailureReconciliation(t *testing.T) {
	for _, outcome := range []string{"completed", "active", "read_failed"} {
		t.Run(outcome, func(t *testing.T) {
			var reads atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				switch {
				case strings.HasSuffix(r.URL.Path, "/events"):
					w.Header().Set("Content-Type", "text/event-stream")
					w.WriteHeader(200)
				case strings.HasSuffix(r.URL.Path, "/turns"):
					_, _ = fmt.Fprint(w, `{"data":[{"id":"root","session_id":"session","status":"in_progress"}],"has_more":false}`)
				case strings.HasSuffix(r.URL.Path, "/turns/root"):
					status := "in_progress"
					if reads.Add(1) > 1 {
						if outcome == "completed" {
							status = "completed"
						}
						if outcome == "read_failed" {
							w.WriteHeader(500)
							_, _ = fmt.Fprint(w, `{"error":{"message":"synthetic"}}`)
							return
						}
					}
					_, _ = fmt.Fprintf(w, `{"id":"root","session_id":"session","status":%q}`, status)
				case strings.HasSuffix(r.URL.Path, "/items"):
					_, _ = fmt.Fprint(w, `{"data":[],"has_more":false}`)
				default:
					_, _ = fmt.Fprint(w, `{"id":"session","status":"in_progress"}`)
				}
			}))
			defer server.Close()
			client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"), option.WithMaxRetries(0))
			result, err := client.Beta.Agents.Sessions.Stream(context.Background(), "session", openai.AgentSessionStreamParams{}).FinalResult()
			if outcome == "completed" {
				if err != nil || result.TurnID() != "root" {
					t.Fatalf("missed durable completion: result=%v err=%v", result, err)
				}
			} else if !errors.Is(err, io.ErrUnexpectedEOF) {
				t.Fatal("lost original observation cause", err)
			}
		})
	}
}

func TestBetaAgentAttachRawIterationDoesNotReadOutputHistory(t *testing.T) {
	var items atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case strings.HasSuffix(r.URL.Path, "/events"):
			w.Header().Set("Content-Type", "text/event-stream")
			_, _ = fmt.Fprintf(w, "data: %s\n\ndata: %s\n\n", betaResultMessage("done", "answer", "root", `"final_answer"`, "private-answer", 0), betaResultTurn("completed", "root", "null"))
		case strings.HasSuffix(r.URL.Path, "/turns"):
			_, _ = fmt.Fprint(w, `{"data":[{"id":"root","session_id":"session","status":"in_progress"}],"has_more":false}`)
		case strings.HasSuffix(r.URL.Path, "/turns/root"):
			_, _ = fmt.Fprint(w, `{"id":"root","session_id":"session","status":"in_progress"}`)
		case strings.HasSuffix(r.URL.Path, "/items"):
			items.Add(1)
			_, _ = fmt.Fprint(w, `{"data":[],"has_more":false}`)
		default:
			_, _ = fmt.Fprint(w, `{"id":"session","status":"in_progress"}`)
		}
	}))
	defer server.Close()
	client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"))
	stream := client.Beta.Agents.Sessions.Stream(context.Background(), "session", openai.AgentSessionStreamParams{})
	for stream.Next() {
	}
	if stream.Err() != nil || items.Load() != 0 {
		t.Fatal("raw iteration fetched result payloads", stream.Err())
	}
	_, err := stream.FinalResult()
	var failure *openai.BetaAgentTurnResultError
	if !errors.As(err, &failure) || failure.Reason != "collection_not_enabled" || items.Load() != 0 {
		t.Fatal("late collection retained output", err)
	}
}

func TestBetaAgentAttachManualSnapshotIsSelectedRootOnly(t *testing.T) {
	for _, tc := range []struct {
		name, action      string
		manual, successor bool
	}{
		{name: "current approval", action: `{"type":"computer_use_approval_request","turn_id":"root","request_id":"approval"}`, manual: true},
		{name: "old approval", action: `{"type":"computer_use_approval_request","turn_id":"older","request_id":"approval"}`},
		{name: "current environment", action: `{"type":"environment_connection","environment_id":"env"}`, manual: true},
		{name: "successor environment", action: `{"type":"environment_connection","environment_id":"env"}`, successor: true},
		{name: "stale function", action: `{"type":"function_call","turn_id":"root","name":"unknown","call_id":"answered","arguments":{}}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var reads, lists atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				switch {
				case strings.HasSuffix(r.URL.Path, "/events"):
					w.Header().Set("Content-Type", "text/event-stream")
					if !tc.manual {
						_, _ = fmt.Fprintf(w, "data: %s\n\n", betaResultTurn("completed", "root", "null"))
					}
					w.WriteHeader(200)
					w.(http.Flusher).Flush()
					<-r.Context().Done()
				case strings.HasSuffix(r.URL.Path, "/turns"):
					id := "root"
					if lists.Add(1) > 1 && tc.successor {
						id = "successor"
					}
					_, _ = fmt.Fprintf(w, `{"data":[{"id":%q,"session_id":"session","status":"waiting"}],"has_more":false}`, id)
				case strings.HasSuffix(r.URL.Path, "/turns/root"):
					status := "waiting"
					if reads.Add(1) > 1 && !tc.manual {
						status = "completed"
					}
					_, _ = fmt.Fprintf(w, `{"id":"root","session_id":"session","status":%q}`, status)
				case strings.HasSuffix(r.URL.Path, "/items"):
					_, _ = fmt.Fprint(w, `{"data":[],"has_more":false}`)
				default:
					_, _ = fmt.Fprintf(w, `{"id":"session","status":"requires_action","required_actions":[%s]}`, tc.action)
				}
			}))
			defer server.Close()
			client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"))
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			result, err := client.Beta.Agents.Sessions.Stream(ctx, "session", openai.AgentSessionStreamParams{}).FinalResult()
			if tc.manual {
				var failure *openai.BetaAgentTurnResultError
				if !errors.As(err, &failure) || failure.Reason != "requires_action" || len(failure.RequiredActions) != 1 {
					t.Fatal("missing current manual action", err)
				}
			} else if err != nil || result.TurnID() != "root" {
				t.Fatalf("unrelated snapshot blocked selected root: result=%v err=%v", result, err)
			}
		})
	}
}

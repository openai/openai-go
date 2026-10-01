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
	var lists, snapshots atomic.Int32
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
			status := "in_progress"
			if snapshots.Add(1) > 1 {
				status = "idle"
			}
			_, _ = fmt.Fprintf(w, `{"id":"session","status":%q}`, status)
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
		name, action              string
		manual, successor, closed bool
	}{
		{name: "closed environment", action: `{"type":"environment_connection","environment_id":"env"}`, manual: true, closed: true},
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
					w.WriteHeader(200)
					if !tc.manual {
						_, _ = fmt.Fprintf(w, "data: %s\n\n", betaResultTurn("completed", "root", "null"))
					}
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
			stream := client.Beta.Agents.Sessions.Stream(ctx, "session", openai.AgentSessionStreamParams{})
			if tc.closed {
				_ = stream.Close()
			}
			result, err := stream.FinalResult()
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

func TestBetaAgentAttachIgnoresStaleIdle(t *testing.T) {
	var reads atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case strings.HasSuffix(r.URL.Path, "/events"):
			w.Header().Set("Content-Type", "text/event-stream")
			_, _ = fmt.Fprintf(w, "data: %s\n\ndata: %s\n\n", betaResultIdle(), betaResultTurn("completed", "root", "null"))
		case strings.HasSuffix(r.URL.Path, "/turns"):
			_, _ = fmt.Fprint(w, `{"data":[{"id":"root","session_id":"session","status":"in_progress"}],"has_more":false}`)
		case strings.HasSuffix(r.URL.Path, "/turns/root"):
			status := "in_progress"
			if reads.Add(1) > 2 {
				status = "completed"
			}
			_, _ = fmt.Fprintf(w, `{"id":"root","session_id":"session","status":%q}`, status)
		case strings.HasSuffix(r.URL.Path, "/items"):
			_, _ = fmt.Fprint(w, `{"data":[],"has_more":false}`)
		default:
			_, _ = fmt.Fprint(w, `{"id":"session","status":"idle"}`)
		}
	}))
	defer server.Close()
	client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"))
	result, err := client.Beta.Agents.Sessions.Stream(context.Background(), "session", openai.AgentSessionStreamParams{}).FinalResult()
	if err != nil || result.TurnID() != "root" || reads.Load() != 3 {
		t.Fatalf("stale idle stopped active root: result=%v reads=%d err=%v", result, reads.Load(), err)
	}
}

func TestBetaAgentAttachTerminalErrorFidelity(t *testing.T) {
	for _, status := range []string{"failed", "cancelled"} {
		for _, getFails := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/getFails=%t", status, getFails), func(t *testing.T) {
				var reads atomic.Int32
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					switch {
					case strings.HasSuffix(r.URL.Path, "/events"):
						w.Header().Set("Content-Type", "text/event-stream")
						_, _ = fmt.Fprintf(w, "data: %s\n\ndata: %s\n\n", betaResultMessage("done", "answer", "root", `"final_answer"`, "observed", 0), betaResultTurn(status, "root", "null"))
					case strings.HasSuffix(r.URL.Path, "/turns"):
						_, _ = fmt.Fprint(w, `{"data":[{"id":"root","session_id":"session","status":"in_progress"}],"has_more":false}`)
					case strings.HasSuffix(r.URL.Path, "/turns/root"):
						current := "in_progress"
						if reads.Add(1) > 1 {
							current = status
							if getFails {
								w.WriteHeader(500)
								_, _ = fmt.Fprint(w, `{"error":{"message":"synthetic"}}`)
								return
							}
						}
						_, _ = fmt.Fprintf(w, `{"id":"root","session_id":"session","status":%q}`, current)
					case strings.HasSuffix(r.URL.Path, "/items"):
						_, _ = fmt.Fprint(w, `{"data":[{"id":"earlier","turn_id":"root","type":"message","role":"assistant","status":"completed","content":[{"type":"output_text","text":"durable earlier"}]}],"has_more":false}`)
					default:
						_, _ = fmt.Fprint(w, `{"id":"session","status":"in_progress"}`)
					}
				}))
				defer server.Close()
				client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"), option.WithMaxRetries(0))
				_, err := client.Beta.Agents.Sessions.Stream(context.Background(), "session", openai.AgentSessionStreamParams{}).FinalResult()
				var failure *openai.BetaAgentTurnResultError
				wantMessages := 2
				if getFails {
					wantMessages = 1
				}
				if !errors.As(err, &failure) || failure.Reason != "turn_"+status || failure.Turn.Status != openai.TurnStatus(status) || len(failure.Messages) != wantMessages {
					t.Fatal("terminal status/messages lost", err)
				}
				text := failure.Messages[0].Content[0].Text
				expected := "durable earlier"
				if getFails {
					expected = "observed"
				}
				if text != expected {
					t.Fatalf("wrong partial output: %q", text)
				}
			})
		}
	}
}

func TestBetaAgentAttachHistoricalRootDoesNotSelectWork(t *testing.T) {
	for _, historical := range []string{"baseline", "older"} {
		t.Run(historical, func(t *testing.T) {
			var posts atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				switch {
				case strings.HasSuffix(r.URL.Path, "/events") && r.Method == http.MethodPost:
					posts.Add(1)
					w.WriteHeader(204)
				case strings.HasSuffix(r.URL.Path, "/events"):
					w.Header().Set("Content-Type", "text/event-stream")
					old := agentEvent("turn.item.added", "historical", fmt.Sprintf(`,"turn_id":%q,"item":{"type":"tool_call","turn_id":%q}`, historical, historical))
					for _, event := range []string{old, agentCall("pending", "new", "call", "lookup", `{}`), betaResultTurn("completed", "new", "null")} {
						_, _ = fmt.Fprintf(w, "data: %s\n\n", event)
					}
				case strings.HasSuffix(r.URL.Path, "/turns"):
					_, _ = fmt.Fprint(w, `{"data":[{"id":"baseline","session_id":"session","status":"completed"}],"has_more":false}`)
				case strings.HasSuffix(r.URL.Path, "/turns/new"):
					status := "in_progress"
					if posts.Load() > 0 {
						status = "completed"
					}
					_, _ = fmt.Fprintf(w, `{"id":"new","session_id":"session","status":%q}`, status)
				case strings.HasSuffix(r.URL.Path, "/turns/"+historical):
					_, _ = fmt.Fprintf(w, `{"id":%q,"session_id":"session","status":"completed"}`, historical)
				case strings.HasSuffix(r.URL.Path, "/items"):
					_, _ = fmt.Fprint(w, `{"data":[],"has_more":false}`)
				default:
					_, _ = fmt.Fprint(w, `{"id":"session","status":"in_progress"}`)
				}
			}))
			defer server.Close()
			client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"))
			calls := 0
			result, err := client.Beta.Agents.Sessions.Stream(context.Background(), "session", openai.AgentSessionStreamParams{ToolHandlers: map[string]openai.AgentToolHandler{"lookup": func(context.Context, map[string]any) (any, error) { calls++; return "receipt", nil }}}).FinalResult()
			if err != nil || result.TurnID() != "new" || calls != 1 {
				t.Fatalf("historical root selected: result=%v calls=%d err=%v", result, calls, err)
			}
		})
	}
}

func TestBetaAgentAttachEarlyFailureKeepsSessionID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"id":"session","status":"failed"}`)
	}))
	defer server.Close()
	client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"))
	_, err := client.Beta.Agents.Sessions.Stream(context.Background(), "session", openai.AgentSessionStreamParams{}).FinalResult()
	var failure *openai.BetaAgentTurnResultError
	if !errors.As(err, &failure) || failure.Reason != "session_failed" || failure.SessionID != "session" {
		t.Fatal("lost known session ID", err)
	}
}

func TestBetaAgentAttachSuccessorEnvironmentEvent(t *testing.T) {
	var lists, reads atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case strings.HasSuffix(r.URL.Path, "/events"):
			w.Header().Set("Content-Type", "text/event-stream")
			action := agentEvent("requires_action", "successor-action", `,"session":{"id":"session","required_actions":[{"type":"environment_connection","environment_id":"env"}]}`)
			_, _ = fmt.Fprintf(w, "data: %s\n\ndata: %s\n\n", action, betaResultTurn("completed", "root", "null"))
		case strings.HasSuffix(r.URL.Path, "/turns"):
			id, status := "root", "in_progress"
			if lists.Add(1) > 1 {
				id, status = "successor", "waiting"
			}
			_, _ = fmt.Fprintf(w, `{"data":[{"id":%q,"session_id":"session","status":%q}],"has_more":false}`, id, status)
		case strings.HasSuffix(r.URL.Path, "/turns/root"):
			status := "in_progress"
			if reads.Add(1) > 1 {
				status = "completed"
			}
			_, _ = fmt.Fprintf(w, `{"id":"root","session_id":"session","status":%q}`, status)
		case strings.HasSuffix(r.URL.Path, "/items"):
			_, _ = fmt.Fprint(w, `{"data":[],"has_more":false}`)
		default:
			_, _ = fmt.Fprint(w, `{"id":"session","status":"in_progress"}`)
		}
	}))
	defer server.Close()
	client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"))
	result, err := client.Beta.Agents.Sessions.Stream(context.Background(), "session", openai.AgentSessionStreamParams{}).FinalResult()
	if err != nil || result.TurnID() != "root" {
		t.Fatalf("successor action stopped selected root: result=%v err=%v", result, err)
	}
}

func TestBetaAgentAttachRejectsConflictingItemTurn(t *testing.T) {
	var posts atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodPost:
			posts.Add(1)
			w.WriteHeader(204)
		case strings.HasSuffix(r.URL.Path, "/events"):
			w.Header().Set("Content-Type", "text/event-stream")
			event := agentEvent("turn.item.added", "bad", `,"turn_id":"root","item":{"type":"function_call","turn_id":"other","call_id":"call","name":"lookup","arguments":{}}`)
			_, _ = fmt.Fprintf(w, "data: %s\n\n", event)
		case strings.HasSuffix(r.URL.Path, "/turns"):
			_, _ = fmt.Fprint(w, `{"data":[{"id":"root","session_id":"session","status":"in_progress"}],"has_more":false}`)
		case strings.HasSuffix(r.URL.Path, "/turns/root"):
			_, _ = fmt.Fprint(w, `{"id":"root","session_id":"session","status":"in_progress"}`)
		default:
			_, _ = fmt.Fprint(w, `{"id":"session","status":"in_progress"}`)
		}
	}))
	defer server.Close()
	client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"))
	calls := 0
	_, err := client.Beta.Agents.Sessions.Stream(context.Background(), "session", openai.AgentSessionStreamParams{ToolHandlers: map[string]openai.AgentToolHandler{"lookup": func(context.Context, map[string]any) (any, error) { calls++; return "wrong", nil }}}).FinalResult()
	if err == nil || calls != 0 || posts.Load() != 0 {
		t.Fatalf("conflicting identity dispatched: calls=%d posts=%d err=%v", calls, posts.Load(), err)
	}
}

func TestBetaAgentAttachStaleIdleWithoutVisibleRoot(t *testing.T) {
	var posts atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodPost:
			posts.Add(1)
			w.WriteHeader(204)
		case strings.HasSuffix(r.URL.Path, "/events"):
			w.Header().Set("Content-Type", "text/event-stream")
			for _, event := range []string{betaResultIdle(), agentCall("pending", "new", "call", "lookup", `{}`), betaResultTurn("completed", "new", "null")} {
				_, _ = fmt.Fprintf(w, "data: %s\n\n", event)
			}
		case strings.HasSuffix(r.URL.Path, "/turns"):
			_, _ = fmt.Fprint(w, `{"data":[],"has_more":false}`)
		case strings.HasSuffix(r.URL.Path, "/turns/new"):
			status := "in_progress"
			if posts.Load() > 0 {
				status = "completed"
			}
			_, _ = fmt.Fprintf(w, `{"id":"new","session_id":"session","status":%q}`, status)
		case strings.HasSuffix(r.URL.Path, "/items"):
			_, _ = fmt.Fprint(w, `{"data":[],"has_more":false}`)
		default:
			_, _ = fmt.Fprint(w, `{"id":"session","status":"in_progress"}`)
		}
	}))
	defer server.Close()
	client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"))
	result, err := client.Beta.Agents.Sessions.Stream(context.Background(), "session", openai.AgentSessionStreamParams{ToolHandlers: map[string]openai.AgentToolHandler{"lookup": func(context.Context, map[string]any) (any, error) { return "receipt", nil }}}).FinalResult()
	if err != nil || result.TurnID() != "new" || posts.Load() != 1 {
		t.Fatalf("stale idle stopped new work: %v %v", result, err)
	}
}

func TestBetaAgentAttachSettlesBeforeSuccessorCall(t *testing.T) {
	for _, lagged := range []bool{false, true} {
		t.Run(fmt.Sprint(lagged), func(t *testing.T) {
			var reads atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				switch {
				case strings.HasSuffix(r.URL.Path, "/events"):
					w.Header().Set("Content-Type", "text/event-stream")
					if lagged {
						_, _ = fmt.Fprintf(w, "data: %s\n\n", agentEvent("turn.created", "successor-created", `,"turn_id":"other","turn":{"id":"other","session_id":"session","status":"in_progress"}`))
					}
					_, _ = fmt.Fprintf(w, "data: %s\n\n", agentCall("successor", "other", "call", "lookup", `{}`))
					w.(http.Flusher).Flush()
					<-r.Context().Done()
				case strings.HasSuffix(r.URL.Path, "/turns"):
					_, _ = fmt.Fprint(w, `{"data":[{"id":"root","session_id":"session","status":"in_progress"}],"has_more":false}`)
				case strings.HasSuffix(r.URL.Path, "/turns/root"):
					status := "in_progress"
					threshold := int32(1)
					if lagged {
						threshold = 2
					}
					if reads.Add(1) > threshold {
						status = "completed"
					}
					_, _ = fmt.Fprintf(w, `{"id":"root","session_id":"session","status":%q}`, status)
				case strings.HasSuffix(r.URL.Path, "/items"):
					_, _ = fmt.Fprint(w, `{"data":[],"has_more":false}`)
				default:
					_, _ = fmt.Fprint(w, `{"id":"session","status":"in_progress"}`)
				}
			}))
			defer server.Close()
			client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"))
			calls := 0
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			result, err := client.Beta.Agents.Sessions.Stream(ctx, "session", openai.AgentSessionStreamParams{ToolHandlers: map[string]openai.AgentToolHandler{"lookup": func(context.Context, map[string]any) (any, error) { calls++; return "wrong", nil }}}).FinalResult()
			if err != nil || result.TurnID() != "root" || calls != 0 {
				t.Fatalf("successor dispatched or hung: %v %v", result, err)
			}

		})
	}
}

func TestBetaAgentAttachEnvironmentWithoutActiveTurn(t *testing.T) {
	for _, history := range []string{`{"data":[],"has_more":false}`, `{"data":[{"id":"old","session_id":"session","status":"completed"}],"has_more":false}`} {
		t.Run(history, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				switch {
				case strings.HasSuffix(r.URL.Path, "/events"):
					w.Header().Set("Content-Type", "text/event-stream")
					w.WriteHeader(200)
					w.(http.Flusher).Flush()
					<-r.Context().Done()
				case strings.HasSuffix(r.URL.Path, "/turns"):
					_, _ = fmt.Fprint(w, history)
				default:
					_, _ = fmt.Fprint(w, `{"id":"session","status":"requires_action","required_actions":[{"type":"environment_connection","environment_id":"env"}]}`)
				}
			}))
			defer server.Close()
			client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"))
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			_, err := client.Beta.Agents.Sessions.Stream(ctx, "session", openai.AgentSessionStreamParams{}).FinalResult()
			var failure *openai.BetaAgentTurnResultError
			if !errors.As(err, &failure) || failure.Reason != "requires_action" || failure.SessionID != "session" || failure.Turn != nil {
				t.Fatal("lost pre-input environment diagnostic", err)
			}

		})
	}
}

func TestBetaAgentAttachCloseDuringReconciliation(t *testing.T) {
	for _, parentCancel := range []bool{false, true} {
		t.Run(fmt.Sprint(parentCancel), func(t *testing.T) {
			pageStarted := make(chan struct{})
			var pages atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				switch {
				case strings.HasSuffix(r.URL.Path, "/events"):
					w.Header().Set("Content-Type", "text/event-stream")
					_, _ = fmt.Fprintf(w, "data: %s\n\n", betaResultTurn("completed", "root", "null"))
				case strings.HasSuffix(r.URL.Path, "/turns"):
					_, _ = fmt.Fprint(w, `{"data":[{"id":"root","session_id":"session","status":"in_progress"}],"has_more":false}`)
				case strings.HasSuffix(r.URL.Path, "/turns/root"):
					_, _ = fmt.Fprint(w, `{"id":"root","session_id":"session","status":"completed"}`)
				case strings.HasSuffix(r.URL.Path, "/items"):
					pages.Add(1)
					close(pageStarted)
					<-r.Context().Done()
				default:
					_, _ = fmt.Fprint(w, `{"id":"session","status":"in_progress"}`)
				}
			}))
			defer server.Close()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"))
			stream := client.Beta.Agents.Sessions.Stream(ctx, "session", openai.AgentSessionStreamParams{})
			_ = stream.Close() // An explicit new getter still starts bounded recovery.
			done := make(chan error, 1)
			go func() { _, err := stream.FinalResult(); done <- err }()
			select {
			case <-pageStarted:
			case <-time.After(time.Second):
				t.Fatal("recovery did not start")
			}
			if parentCancel {
				cancel()
			} else {
				_ = stream.Close()
			}
			select {
			case err := <-done:
				if !errors.Is(err, context.Canceled) {
					t.Fatal("lost cancellation", err)
				}
			case <-time.After(time.Second):
				t.Fatal("close did not interrupt recovery")
			}
			if pages.Load() != 1 {
				t.Fatal("requested extra pages")
			}
		})
	}
}

func TestBetaAgentAttachProgressActionSatisfiedBeforeResult(t *testing.T) {
	var resumed atomic.Bool
	var reads atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case strings.HasSuffix(r.URL.Path, "/events"):
			w.Header().Set("Content-Type", "text/event-stream")
			action := agentEvent("requires_action", "manual", `,"session":{"id":"session","required_actions":[{"type":"computer_use_approval_request","turn_id":"root","request_id":"approval"}]}`)
			_, _ = fmt.Fprintf(w, "data: %s\n\ndata: %s\n\n", action, betaResultTurn("completed", "root", "null"))
		case strings.HasSuffix(r.URL.Path, "/turns"):
			_, _ = fmt.Fprint(w, `{"data":[{"id":"root","session_id":"session","status":"waiting"}],"has_more":false}`)
		case strings.HasSuffix(r.URL.Path, "/turns/root"):
			status := "waiting"
			if reads.Add(1) > 1 {
				status = "completed"
			}
			_, _ = fmt.Fprintf(w, `{"id":"root","session_id":"session","status":%q}`, status)
		case strings.HasSuffix(r.URL.Path, "/items"):
			_, _ = fmt.Fprint(w, `{"data":[],"has_more":false}`)
		default:
			status := "requires_action"
			if resumed.Load() {
				status = "in_progress"
			}
			_, _ = fmt.Fprintf(w, `{"id":"session","status":%q}`, status)
		}
	}))
	defer server.Close()
	client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"))
	stream := client.Beta.Agents.Sessions.Stream(context.Background(), "session", openai.AgentSessionStreamParams{}).WithResultCollection()
	if !stream.Next() {
		t.Fatal(stream.Err())
	}
	resumed.Store(true)
	result, err := stream.FinalResult()
	if err != nil || result.TurnID() != "root" {
		t.Fatalf("resolved action stopped result: %v %v", result, err)
	}
}

func TestBetaAgentAttachFirstItemDiscoversManualApproval(t *testing.T) {
	var lists atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case strings.HasSuffix(r.URL.Path, "/events"):
			w.Header().Set("Content-Type", "text/event-stream")
			event := agentEvent("turn.item.added", "approval", `,"item":{"type":"tool_call","turn_id":"root"}`)
			_, _ = fmt.Fprintf(w, "data: %s\n\n", event)
			w.(http.Flusher).Flush()
			<-r.Context().Done()
		case strings.HasSuffix(r.URL.Path, "/turns"):
			lists.Add(1)
			_, _ = fmt.Fprint(w, `{"data":[],"has_more":false}`)
		case strings.HasSuffix(r.URL.Path, "/turns/root"):
			_, _ = fmt.Fprint(w, `{"id":"root","session_id":"session","status":"waiting"}`)
		default:
			if lists.Load() < 2 {
				_, _ = fmt.Fprint(w, `{"id":"session","status":"in_progress"}`)
			} else {
				_, _ = fmt.Fprint(w, `{"id":"session","status":"requires_action","required_actions":[{"type":"computer_use_approval_request","turn_id":"root","request_id":"approval"}]}`)
			}
		}
	}))
	defer server.Close()
	client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := client.Beta.Agents.Sessions.Stream(ctx, "session", openai.AgentSessionStreamParams{}).FinalResult()
	var failure *openai.BetaAgentTurnResultError
	if !errors.As(err, &failure) || failure.Reason != "requires_action" || len(failure.RequiredActions) != 1 || (failure.Turn == nil || failure.Turn.ID != "root") {
		t.Fatal("first-item approval not diagnosed", err)
	}
}

func TestBetaAgentAttachActionFrameAndProgressSnapshot(t *testing.T) {
	for _, firstFrame := range []bool{false, true} {
		t.Run(fmt.Sprint(firstFrame), func(t *testing.T) {
			var lists atomic.Int32
			action := `{"type":"computer_use_approval_request","turn_id":"root","request_id":"approval"}`
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				switch {
				case strings.HasSuffix(r.URL.Path, "/events"):
					w.Header().Set("Content-Type", "text/event-stream")
					event := agentEvent("requires_action", "approval", `,"session":{"id":"session","required_actions":[`+action+`]}`)
					_, _ = fmt.Fprintf(w, "data: %s\n\n", event)
					w.(http.Flusher).Flush()
					<-r.Context().Done()
				case strings.HasSuffix(r.URL.Path, "/turns"):
					lists.Add(1)
					if firstFrame {
						_, _ = fmt.Fprint(w, `{"data":[],"has_more":false}`)
					} else {
						_, _ = fmt.Fprint(w, `{"data":[{"id":"root","session_id":"session","status":"waiting"}],"has_more":false}`)
					}
				case strings.HasSuffix(r.URL.Path, "/turns/root"):
					_, _ = fmt.Fprint(w, `{"id":"root","session_id":"session","status":"waiting"}`)
				default:
					if firstFrame && lists.Load() < 3 {
						_, _ = fmt.Fprint(w, `{"id":"session","status":"in_progress"}`)
					} else {
						_, _ = fmt.Fprintf(w, `{"id":"session","status":"requires_action","required_actions":[%s]}`, action)
					}
				}
			}))
			defer server.Close()
			client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"))
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			stream := client.Beta.Agents.Sessions.Stream(ctx, "session", openai.AgentSessionStreamParams{}).WithResultCollection()
			if !firstFrame && !stream.Next() {
				t.Fatal(stream.Err())
			}
			_, err := stream.FinalResult()
			var failure *openai.BetaAgentTurnResultError
			if !errors.As(err, &failure) || failure.Reason != "requires_action" || len(failure.RequiredActions) != 1 || failure.Turn == nil || failure.Turn.ID != "root" {
				t.Fatal("missing or duplicated approval", err)
			}
		})
	}
}

func TestBetaAgentAttachRefreshesWaitingActionsAfterClose(t *testing.T) {
	var waiting atomic.Bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case strings.HasSuffix(r.URL.Path, "/events"):
			w.Header().Set("Content-Type", "text/event-stream")
			w.WriteHeader(200)
			w.(http.Flusher).Flush()
			<-r.Context().Done()
		case strings.HasSuffix(r.URL.Path, "/turns"):
			_, _ = fmt.Fprint(w, `{"data":[{"id":"root","session_id":"session","status":"in_progress"}],"has_more":false}`)
		case strings.HasSuffix(r.URL.Path, "/turns/root"):
			status := "in_progress"
			if waiting.Load() {
				status = "waiting"
			}
			_, _ = fmt.Fprintf(w, `{"id":"root","session_id":"session","status":%q}`, status)
		default:
			_, _ = fmt.Fprint(w, `{"id":"session","status":"requires_action","required_actions":[{"type":"computer_use_approval_request","turn_id":"root","request_id":"approval"}]}`)
		}
	}))
	defer server.Close()
	client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"))
	stream := client.Beta.Agents.Sessions.Stream(context.Background(), "session", openai.AgentSessionStreamParams{})
	_ = stream.Close()
	waiting.Store(true)
	_, err := stream.FinalResult()
	var failure *openai.BetaAgentTurnResultError
	if !errors.As(err, &failure) || failure.Reason != "requires_action" || len(failure.RequiredActions) != 1 {
		t.Fatal("closed recovery lost new approval", err)
	}
}

func TestBetaAgentAttachPrunesAcknowledgedFunctionDiagnostic(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case strings.HasSuffix(r.URL.Path, "/events"):
			w.Header().Set("Content-Type", "text/event-stream")
			_, _ = fmt.Fprintf(w, "data: %s\n\n", agentCall("call-event", "root", "answered", "unknown", `{}`))
			w.(http.Flusher).Flush()
			<-r.Context().Done()
		case strings.HasSuffix(r.URL.Path, "/turns"):
			_, _ = fmt.Fprint(w, `{"data":[{"id":"root","session_id":"session","status":"waiting"}],"has_more":false}`)
		case strings.HasSuffix(r.URL.Path, "/turns/root"):
			_, _ = fmt.Fprint(w, `{"id":"root","session_id":"session","status":"waiting"}`)
		default:
			_, _ = fmt.Fprint(w, `{"id":"session","status":"requires_action","required_actions":[{"type":"computer_use_approval_request","turn_id":"root","request_id":"approval"}]}`)
		}
	}))
	defer server.Close()
	client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"))
	stream := client.Beta.Agents.Sessions.Stream(context.Background(), "session", openai.AgentSessionStreamParams{}).WithResultCollection()
	if !stream.Next() {
		t.Fatal(stream.Err())
	}
	_, err := stream.FinalResult()
	var failure *openai.BetaAgentTurnResultError
	if !errors.As(err, &failure) || failure.Reason != "requires_action" || len(failure.RequiredActions) != 1 || failure.RequiredActions[0].Type != "computer_use_approval_request" {
		t.Fatal("resolved function diagnostic retained", err)
	}
}

func TestBetaAgentAttachRejectsConflictingEmbeddedTurn(t *testing.T) {
	var posts atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodPost:
			posts.Add(1)
			w.WriteHeader(204)
		case strings.HasSuffix(r.URL.Path, "/events"):
			w.Header().Set("Content-Type", "text/event-stream")
			event := agentEvent("turn.completed", "bad", `,"turn_id":"root","turn":{"id":"other","session_id":"session","status":"completed"}`)
			_, _ = fmt.Fprintf(w, "data: %s\n\n", event)
		case strings.HasSuffix(r.URL.Path, "/turns"):
			_, _ = fmt.Fprint(w, `{"data":[{"id":"root","session_id":"session","status":"in_progress"}],"has_more":false}`)
		case strings.HasSuffix(r.URL.Path, "/turns/root"):
			_, _ = fmt.Fprint(w, `{"id":"root","session_id":"session","status":"in_progress"}`)
		default:
			_, _ = fmt.Fprint(w, `{"id":"session","status":"in_progress"}`)
		}
	}))
	defer server.Close()
	client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"))
	calls := 0
	_, err := client.Beta.Agents.Sessions.Stream(context.Background(), "session", openai.AgentSessionStreamParams{ToolHandlers: map[string]openai.AgentToolHandler{"lookup": func(context.Context, map[string]any) (any, error) { calls++; return "wrong", nil }}}).FinalResult()
	if err == nil || calls != 0 || posts.Load() != 0 {
		t.Fatalf("conflicting identity dispatched: calls=%d posts=%d err=%v", calls, posts.Load(), err)
	}
}

func TestBetaAgentAttachDeduplicatesPendingFunctionDiagnostics(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case strings.HasSuffix(r.URL.Path, "/events"):
			w.Header().Set("Content-Type", "text/event-stream")
			events := []string{
				agentCall("first", "root", "answered", "unknown", `{}`),
				agentCall("replayed", "root", "answered", "unknown", `{}`),
				agentEvent("in_progress", "reset", `,"session":{"id":"session","status":"in_progress"}`),
				agentCall("after-reset", "root", "answered", "unknown", `{}`),
				agentEvent("requires_action", "replace", `,"session":{"id":"session","status":"requires_action","required_actions":[{"type":"function_call","turn_id":"root","call_id":"answered","name":"unknown","arguments":{}}]}`),
				agentCall("after-replacement", "root", "answered", "unknown", `{}`),
			}
			for _, event := range events {
				_, _ = fmt.Fprintf(w, "data: %s\n\n", event)
			}
			w.(http.Flusher).Flush()
			<-r.Context().Done()
		case strings.HasSuffix(r.URL.Path, "/turns"):
			_, _ = fmt.Fprint(w, `{"data":[{"id":"root","session_id":"session","status":"waiting"}],"has_more":false}`)
		case strings.HasSuffix(r.URL.Path, "/turns/root"):
			_, _ = fmt.Fprint(w, `{"id":"root","session_id":"session","status":"waiting"}`)
		default:
			_, _ = fmt.Fprint(w, `{"id":"session","status":"requires_action","required_actions":[{"type":"function_call","turn_id":"root","call_id":"answered","name":"unknown","arguments":{}}]}`)
		}
	}))
	defer server.Close()
	client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"))
	stream := client.Beta.Agents.Sessions.Stream(context.Background(), "session", openai.AgentSessionStreamParams{}).WithResultCollection()
	for range 6 {
		if !stream.Next() {
			t.Fatal(stream.Err())
		}
	}
	_, err := stream.FinalResult()
	var failure *openai.BetaAgentTurnResultError
	if !errors.As(err, &failure) || failure.Reason != "requires_action" || len(failure.RequiredActions) != 1 || failure.RequiredActions[0].Type != "function_call" {
		t.Fatal("replayed function diagnostic duplicated", err)
	}
}

func TestBetaAgentAttachKnownSessionFailureDoesNotRefresh(t *testing.T) {
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if requests.Add(1) > 1 {
			t.Error("requested snapshot after known session failure")
			http.Error(w, "unavailable", 503)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"id":"session","status":"failed"}`)
	}))
	defer server.Close()
	client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"), option.WithMaxRetries(0))
	_, err := client.Beta.Agents.Sessions.Stream(context.Background(), "session", openai.AgentSessionStreamParams{}).FinalResult()
	var failure *openai.BetaAgentTurnResultError
	if !errors.As(err, &failure) || failure.Reason != "session_failed" || requests.Load() != 1 {
		t.Fatalf("requests=%d err=%v", requests.Load(), err)
	}
}

func TestBetaAgentAttachRecoversAfterSubscriptionFailure(t *testing.T) {
	var subscriptions atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case strings.HasSuffix(r.URL.Path, "/events"):
			subscriptions.Add(1)
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = fmt.Fprint(w, `{"error":{"message":"synthetic temporary failure"}}`)
		case strings.HasSuffix(r.URL.Path, "/turns"):
			_, _ = fmt.Fprint(w, `{"data":[{"id":"root","session_id":"session","status":"in_progress"}],"has_more":false}`)
		case strings.HasSuffix(r.URL.Path, "/turns/root"):
			_, _ = fmt.Fprint(w, `{"id":"root","session_id":"session","status":"completed"}`)
		case strings.HasSuffix(r.URL.Path, "/items"):
			_, _ = fmt.Fprint(w, `{"data":[{"id":"m","turn_id":"root","type":"message","role":"assistant","status":"completed","content":[{"type":"output_text","text":"finished"}]}],"has_more":false}`)
		default:
			_, _ = fmt.Fprint(w, `{"id":"session","status":"in_progress"}`)
		}
	}))
	defer server.Close()
	client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"), option.WithMaxRetries(0))
	stream := client.Beta.Agents.Sessions.Stream(context.Background(), "session", openai.AgentSessionStreamParams{})
	if stream.Err() == nil {
		t.Fatal("raw stream lost subscription error")
	}
	result, err := stream.FinalResult()
	if err != nil || result.OutputText() != "finished" || subscriptions.Load() != 1 {
		t.Fatalf("result=%v subscriptions=%d err=%v", result, subscriptions.Load(), err)
	}
}

func TestBetaAgentAttachRecoversAfterSnapshotFailure(t *testing.T) {
	for _, failedRead := range []string{"turn", "session", "idle", "manual"} {
		t.Run(failedRead, func(t *testing.T) {
			var turnReads, sessionReads atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				fail := func() {
					w.WriteHeader(http.StatusServiceUnavailable)
					_, _ = fmt.Fprint(w, `{"error":{"message":"synthetic temporary failure"}}`)
				}
				switch {
				case strings.HasSuffix(r.URL.Path, "/events"):
					w.Header().Set("Content-Type", "text/event-stream")
					_, _ = fmt.Fprintf(w, "data: %s\n\n", betaResultIdle())
					w.(http.Flusher).Flush()
					<-r.Context().Done()
				case strings.HasSuffix(r.URL.Path, "/turns"):
					_, _ = fmt.Fprint(w, `{"data":[{"id":"root","session_id":"session","status":"in_progress"}],"has_more":false}`)
				case strings.HasSuffix(r.URL.Path, "/turns/root"):
					count := turnReads.Add(1)
					if (failedRead == "turn" && count == 1) || (failedRead == "idle" && count == 2) {
						fail()
						return
					}
					status := "in_progress"
					if failedRead == "manual" {
						status = "waiting"
					}
					if (failedRead != "idle" && count >= 2) || (failedRead == "idle" && count >= 3) {
						status = "completed"
					}
					_, _ = fmt.Fprintf(w, `{"id":"root","session_id":"session","status":%q}`, status)
				case strings.HasSuffix(r.URL.Path, "/items"):
					_, _ = fmt.Fprint(w, `{"data":[{"id":"m","turn_id":"root","type":"message","role":"assistant","status":"completed","content":[{"type":"output_text","text":"finished"}]}],"has_more":false}`)
				default:
					count := sessionReads.Add(1)
					if (count == 2 && failedRead == "session") || (count == 3 && failedRead == "manual") {
						fail()
						return
					}
					_, _ = fmt.Fprint(w, `{"id":"session","status":"in_progress"}`)
				}
			}))
			defer server.Close()
			client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"), option.WithMaxRetries(0))
			result, err := client.Beta.Agents.Sessions.Stream(context.Background(), "session", openai.AgentSessionStreamParams{}).FinalResult()
			if err != nil || result.OutputText() != "finished" {
				t.Fatalf("result=%v err=%v", result, err)
			}
		})
	}
}

func TestBetaAgentAttachStreamedSessionFailureDoesNotRefresh(t *testing.T) {
	var turnReads atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case strings.HasSuffix(r.URL.Path, "/events"):
			w.Header().Set("Content-Type", "text/event-stream")
			_, _ = fmt.Fprintf(w, "data: %s\n\n", agentEvent("failed", "failure", `,"session":{"id":"session","status":"failed"}`))
		case strings.HasSuffix(r.URL.Path, "/turns"):
			_, _ = fmt.Fprint(w, `{"data":[{"id":"root","session_id":"session","status":"in_progress"}],"has_more":false}`)
		case strings.HasSuffix(r.URL.Path, "/turns/root"):
			if turnReads.Add(1) > 1 {
				t.Error("requested snapshot after streamed session failure")
				http.Error(w, "unavailable", 503)
				return
			}
			_, _ = fmt.Fprint(w, `{"id":"root","session_id":"session","status":"in_progress"}`)
		default:
			_, _ = fmt.Fprint(w, `{"id":"session","status":"in_progress"}`)
		}
	}))
	defer server.Close()
	client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"), option.WithMaxRetries(0))
	_, err := client.Beta.Agents.Sessions.Stream(context.Background(), "session", openai.AgentSessionStreamParams{}).FinalResult()
	var failure *openai.BetaAgentTurnResultError
	if !errors.As(err, &failure) || failure.Reason != "session_failed" || turnReads.Load() != 1 {
		t.Fatalf("turn_reads=%d err=%v", turnReads.Load(), err)
	}
}

func TestBetaAgentAttachLargeFunctionDiagnostics(t *testing.T) {
	const count = 512
	var actions []map[string]any
	for i := range count {
		turn := "root"
		if i%2 != 0 {
			turn = "other"
		}
		actions = append(actions, map[string]any{"type": "function_call", "turn_id": turn, "call_id": fmt.Sprint(i), "name": "snapshot", "arguments": map[string]any{}})
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case strings.HasSuffix(r.URL.Path, "/events"):
			w.Header().Set("Content-Type", "text/event-stream")
			for i := range count {
				for repeat := range 2 {
					_, _ = fmt.Fprintf(w, "data: %s\n\n", agentCall(fmt.Sprintf("event-%d-%d", i, repeat), "root", fmt.Sprint(i), "observed", `{}`))
				}
			}
			w.(http.Flusher).Flush()
			<-r.Context().Done()
		case strings.HasSuffix(r.URL.Path, "/turns"):
			_, _ = fmt.Fprint(w, `{"data":[{"id":"root","session_id":"session","status":"waiting"}],"has_more":false}`)
		case strings.HasSuffix(r.URL.Path, "/turns/root"):
			_, _ = fmt.Fprint(w, `{"id":"root","session_id":"session","status":"waiting"}`)
		default:
			_ = json.NewEncoder(w).Encode(map[string]any{"id": "session", "status": "requires_action", "required_actions": actions})
		}
	}))
	defer server.Close()
	client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"))
	stream := client.Beta.Agents.Sessions.Stream(context.Background(), "session", openai.AgentSessionStreamParams{}).WithResultCollection()
	for range count * 2 {
		if !stream.Next() {
			t.Fatal(stream.Err())
		}
	}
	_, err := stream.FinalResult()
	var failure *openai.BetaAgentTurnResultError
	if !errors.As(err, &failure) || failure.Reason != "requires_action" || len(failure.RequiredActions) != count/2 {
		t.Fatal("incorrect pending diagnostics", err)
	}
	for i, action := range failure.RequiredActions {
		if action.CallID != fmt.Sprint(i*2) || action.Name != "observed" {
			t.Fatalf("diagnostic order or SSE payload changed: index=%d call=%s name=%s", i, action.CallID, action.Name)
		}
	}
}

func TestBetaAgentAttachPreservesObservedCompletedMessages(t *testing.T) {
	for _, tc := range []struct {
		name     string
		history  []string
		observed []string
		want     string
	}{
		{"empty history", nil, []string{"b", "c"}, "bc"},
		{"missing tail", []string{"a", "b"}, []string{"b", "c"}, "abc"},
		{"shared anchors", []string{"a", "c", "e"}, []string{"b", "c", "d", "e", "f"}, "abcdef"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var reads atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				switch {
				case strings.HasSuffix(r.URL.Path, "/events"):
					w.Header().Set("Content-Type", "text/event-stream")
					for i, id := range tc.observed {
						_, _ = fmt.Fprintf(w, "data: %s\n\n", betaResultMessage("done", id, "root", `"final_answer"`, id, i*3))
					}
					_, _ = fmt.Fprintf(w, "data: %s\n\n", betaResultTurn("completed", "root", "null"))
				case strings.HasSuffix(r.URL.Path, "/turns"):
					_, _ = fmt.Fprint(w, `{"data":[{"id":"root","session_id":"session","status":"in_progress"}],"has_more":false}`)
				case strings.HasSuffix(r.URL.Path, "/turns/root"):
					status := "in_progress"
					if reads.Add(1) > 1 {
						status = "completed"
					}
					_, _ = fmt.Fprintf(w, `{"id":"root","session_id":"session","status":%q}`, status)
				case strings.HasSuffix(r.URL.Path, "/items"):
					items := make([]map[string]any, 0, len(tc.history))
					for _, id := range tc.history {
						text := id
						for _, live := range tc.observed {
							if live == id {
								text = "stale"
							}
						}
						items = append(items, map[string]any{"id": id, "type": "message", "turn_id": "root", "role": "assistant", "status": "completed", "phase": "final_answer", "content": []map[string]string{{"type": "output_text", "text": text}}})
					}
					_ = json.NewEncoder(w).Encode(map[string]any{"data": items, "has_more": false})
				default:
					_, _ = fmt.Fprint(w, `{"id":"session","status":"in_progress"}`)
				}
			}))
			defer server.Close()
			client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"))
			stream := client.Beta.Agents.Sessions.Stream(context.Background(), "session", openai.AgentSessionStreamParams{})
			result, err := stream.FinalResult()
			if err != nil || result.OutputText() != tc.want {
				t.Fatalf("result=%v err=%v want=%s", result, err, tc.want)
			}
		})
	}
}

func TestBetaAgentAttachSnapshotFailureAfterCloseOrStaleIdle(t *testing.T) {
	for _, closeFirst := range []bool{false, true} {
		t.Run(fmt.Sprintf("close=%t", closeFirst), func(t *testing.T) {
			var failed atomic.Bool
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				switch {
				case strings.HasSuffix(r.URL.Path, "/events"):
					w.Header().Set("Content-Type", "text/event-stream")
					_, _ = fmt.Fprintf(w, "data: %s\n\n", betaResultIdle())
					w.(http.Flusher).Flush()
					<-r.Context().Done()
				case strings.HasSuffix(r.URL.Path, "/turns"):
					data := `[]`
					if closeFirst {
						data = `[{"id":"root","session_id":"session","status":"waiting"}]`
					}
					_, _ = fmt.Fprintf(w, `{"data":%s,"has_more":false}`, data)
				case strings.HasSuffix(r.URL.Path, "/turns/root"):
					if failed.Load() {
						t.Error("read a turn after definitive session failure")
					}
					_, _ = fmt.Fprint(w, `{"id":"root","session_id":"session","status":"waiting"}`)
				default:
					status := "in_progress"
					if failed.Load() {
						status = "failed"
					}
					_, _ = fmt.Fprintf(w, `{"id":"session","status":%q}`, status)
				}
			}))
			defer server.Close()
			client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"))
			stream := client.Beta.Agents.Sessions.Stream(context.Background(), "session", openai.AgentSessionStreamParams{}).WithResultCollection()
			failed.Store(true)
			if closeFirst {
				_ = stream.Close()
			} else if !stream.Next() {
				t.Fatal(stream.Err())
			}
			_, err := stream.FinalResult()
			var failure *openai.BetaAgentTurnResultError
			if !errors.As(err, &failure) || failure.Reason != "session_failed" || failure.SessionID != "session" {
				t.Fatal("snapshot session failure lost", err)
			}
		})
	}
}

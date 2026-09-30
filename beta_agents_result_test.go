package openai_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/ssestream"
)

func betaResultTurn(kind, id, subagent string) string {
	return agentEvent("turn."+kind, kind+id, fmt.Sprintf(`,"session_id":"session","turn_id":%q,"turn":{"id":%q,"session_id":"session","subagent_id":%s,"status":%q,"usage":{"input_tokens":7,"output_tokens":3,"total_tokens":10}}`, id, id, subagent, map[string]string{"created": "in_progress", "completed": "completed", "failed": "failed", "cancelled": "cancelled"}[kind]))
}
func betaResultMessage(kind, id, turn, phase, text string, index int) string {
	return agentEvent("turn.item."+kind, kind+id, fmt.Sprintf(`,"session_id":"session","turn_id":%q,"output_index":%d,"item":{"id":%q,"type":"message","role":"assistant","turn_id":%q,"status":%q,"phase":%s,"content":[{"type":"output_text","text":%q,"annotations":[{"type":"url_citation","url":"https://example.com","title":"policy","start_index":0,"end_index":1}]}]}`, turn, index, id, turn, map[string]string{"added": "in_progress", "done": "completed"}[kind], phase, text))
}
func betaResultIdle() string {
	return agentEvent("idle", "idle", `,"session":{"id":"session","status":"idle"}`)
}
func betaResultCreateStream(t *testing.T, events []string, stayOpen bool) *ssestream.Stream[openai.AgentSessionEventUnion] {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		for _, event := range events {
			_, _ = fmt.Fprintf(w, "data: %s\n\n", event)
			w.(http.Flusher).Flush()
		}
		if stayOpen {
			<-r.Context().Done()
		}
	}))
	t.Cleanup(server.Close)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(cancel)
	client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"), option.WithMaxRetries(0))
	// Keep a concrete assignment: changing NewStreaming's return type is a source break.
	var stream *ssestream.Stream[openai.AgentSessionEventUnion] = client.Beta.Agents.Sessions.NewStreaming(ctx, openai.BetaAgentSessionNewParams{})
	t.Cleanup(func() { _ = stream.Close() })
	return stream
}

func TestBetaAgentFinalResultCreation(t *testing.T) {
	events := []string{
		betaResultTurn("created", "child", `"worker"`), betaResultTurn("completed", "child", `"worker"`), betaResultIdle(),
		betaResultTurn("created", "root", "null"),
		betaResultMessage("added", "m1", "root", `"final_answer"`, "", 0),
		betaResultMessage("done", "comment", "root", `"commentary"`, "thinking", 1),
		betaResultMessage("done", "child-answer", "child", `"final_answer"`, "wrong", 0),
		betaResultMessage("done", "m2", "root", `"final_answer"`, " world", 3),
		betaResultMessage("done", "m1", "root", `"final_answer"`, "hello", 0),
		betaResultTurn("completed", "root", "null"), betaResultIdle(),
	}
	for _, prefix := range []int{0, 5, len(events)} {
		t.Run(fmt.Sprint(prefix), func(t *testing.T) {
			stream := betaResultCreateStream(t, events, true)
			for i := 0; i < prefix; i++ {
				if !stream.Next() {
					t.Fatal(stream.Err())
				}
				event := stream.Current()
				if event.Item.Type == "message" {
					event.Item.Phase = "commentary"
				}
				event.Turn.Status = "failed"
			}
			result, err := openai.BetaAgentSessionFinalResult(stream)
			if err != nil {
				t.Fatal(err)
			}
			if result.OutputText() != "hello world" || result.SessionID() != "session" || result.TurnID() != "root" || len(result.Messages) != 2 {
				t.Fatalf("wrong result: %#v", result)
			}
			if !strings.Contains(result.Messages[0].RawJSON(), "url_citation") || result.Turn.Usage.TotalTokens != 10 {
				t.Fatal("metadata lost")
			}
			second, err := openai.BetaAgentSessionFinalResult(stream)
			if err != nil || second != result {
				t.Fatal("result not cached")
			}
			if stream.Next() {
				t.Fatal("getter must close observation at idle boundary")
			}
		})
	}
}

func TestBetaAgentFinalResultOutcomes(t *testing.T) {
	base := []string{betaResultTurn("created", "root", "null")}
	cases := []struct {
		name   string
		tail   []string
		reason string
	}{
		{"empty", []string{betaResultTurn("completed", "root", "null"), betaResultIdle()}, ""},
		{"failed", []string{betaResultTurn("failed", "root", "null"), betaResultIdle()}, "turn_failed"},
		{"cancelled", []string{betaResultTurn("cancelled", "root", "null"), betaResultIdle()}, "turn_cancelled"},
		{"missing-idle", []string{betaResultTurn("completed", "root", "null")}, "observation_incomplete"},
		{"null-phase", []string{betaResultMessage("done", "m", "root", "null", "answer", 0), betaResultTurn("completed", "root", "null"), betaResultIdle()}, ""},
		{"session-failed", []string{agentEvent("failed", "failure", `,"session":{"id":"session","status":"failed"}`)}, "session_failed"},
		{"protocol-error", []string{`{"type":"error","error":{"message":"synthetic delivery error"}}`}, "observation_failed"},
		{"required-action", []string{agentEvent("requires_action", "pending", `,"session":{"id":"session","required_actions":[{"type":"function_call","turn_id":"root","call_id":"call","name":"lookup","arguments":{}}]}`)}, "requires_action"},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			stream := betaResultCreateStream(t, append(append([]string{}, base...), test.tail...), test.name == "required-action")
			result, err := openai.BetaAgentSessionFinalResult(stream)
			if test.reason == "" {
				expectedText := ""
				if test.name == "null-phase" {
					expectedText = "answer"
				}
				if err != nil || result.OutputText() != expectedText {
					t.Fatalf("result=%v err=%v", result, err)
				}
				return
			}
			var failure *openai.BetaAgentTurnResultError
			if !errors.As(err, &failure) || failure.Reason != test.reason || failure.Turn == nil || failure.Turn.ID != "root" {
				t.Fatalf("error=%#v", err)
			}
			if test.name == "missing-idle" && !errors.Is(err, io.ErrUnexpectedEOF) {
				t.Fatal("lost EOF cause")
			}
			if test.name == "protocol-error" {
				var cause *ssestream.StreamError
				if !errors.As(err, &cause) {
					t.Fatal("lost transport cause")
				}
			}
			if test.name == "required-action" && len(failure.RequiredActions) != 1 {
				t.Fatal("lost pending action")
			}
			_, again := openai.BetaAgentSessionFinalResult(stream)
			if !errors.Is(again, err) {
				t.Fatal("error not cached")
			}
		})
	}
}

func TestBetaAgentFinalResultFollowupHandlers(t *testing.T) {
	events := []string{agentEvent("idle", "initial-idle", ""), betaResultTurn("created", "root", "null"), agentCall("call-event", "root", "call", "lookup", `{}`),
		agentEvent("requires_action", "pending", `,"session":{"id":"session","required_actions":[{"type":"function_call","turn_id":"root","call_id":"call","name":"lookup","arguments":{}}]}`),
		betaResultMessage("done", "answer", "root", `"final_answer"`, "found", 0), betaResultTurn("completed", "root", "null"), betaResultIdle()}
	mock, client := newAgentHelperServer(t, events...)
	calls := 0
	stream := client.Beta.Agents.Sessions.Stream(context.Background(), "session", openai.AgentSessionStreamParams{Input: "lookup", ToolHandlers: map[string]openai.AgentToolHandler{"lookup": func(context.Context, map[string]any) (any, error) { calls++; return "found", nil }}})
	result, err := stream.FinalResult()
	if err != nil || result.OutputText() != "found" {
		t.Fatalf("result=%v err=%v", result, err)
	}
	_, err = stream.FinalResult()
	if err != nil || calls != 1 || len(mock.submissions()) != 2 {
		t.Fatalf("calls=%d submissions=%d err=%v", calls, len(mock.submissions()), err)
	}
}

func TestBetaAgentFinalResultEarlyClose(t *testing.T) {
	stream := betaResultCreateStream(t, []string{betaResultTurn("created", "root", "null")}, true)
	if !stream.Next() {
		t.Fatal(stream.Err())
	}
	_ = stream.Close()
	_, err := openai.BetaAgentSessionFinalResult(stream)
	var failure *openai.BetaAgentTurnResultError
	if !errors.As(err, &failure) || failure.Reason != "observation_incomplete" {
		t.Fatal(err)
	}
}

func TestBetaAgentFinalResultAfterIteration(t *testing.T) {
	events := []string{betaResultTurn("created", "root", "null"), betaResultMessage("done", "m", "root", `"final_answer"`, "answer", 0), betaResultTurn("completed", "root", "null"), betaResultIdle()}
	stream := betaResultCreateStream(t, events, false)
	for stream.Next() {
		item := stream.Current().Item
		if len(item.Content.OfAgentSessionMessageContentArray) > 0 {
			item.Content.OfAgentSessionMessageContentArray[0].Text = "mutated"
		}
		if len(item.Content.OfOutputTextArray) > 0 {
			item.Content.OfOutputTextArray[0].Text = "mutated"
		}
	}
	result, err := openai.BetaAgentSessionFinalResult(stream)
	if err != nil || result.OutputText() != "answer" {
		t.Fatalf("result=%v err=%v", result, err)
	}
	_, client := newAgentHelperServer(t, events...)
	followup := client.Beta.Agents.Sessions.Stream(context.Background(), "session", openai.AgentSessionStreamParams{Input: "question"})
	for followup.Next() {
	}
	result, err = followup.FinalResult()
	if err != nil || result.OutputText() != "answer" {
		t.Fatalf("result=%v err=%v", result, err)
	}
}

func TestBetaAgentFinalResultUnhandledFollowup(t *testing.T) {
	_, client := newAgentHelperServer(t, betaResultTurn("created", "root", "null"), agentEvent("requires_action", "pending", `,"session":{"id":"session","required_actions":[{"type":"function_call","turn_id":"root","call_id":"call","name":"lookup","arguments":{}}]}`))
	stream := client.Beta.Agents.Sessions.Stream(context.Background(), "session", openai.AgentSessionStreamParams{Input: "question"})
	_, err := stream.FinalResult()
	var failure *openai.BetaAgentTurnResultError
	if !errors.As(err, &failure) || failure.Reason != "requires_action" {
		t.Fatal(err)
	}
}

func TestBetaAgentFinalResultIgnoresLaterTurns(t *testing.T) {
	events := []string{betaResultTurn("created", "root", "null"), betaResultMessage("done", "m", "root", `"final_answer"`, "answer", 0), betaResultTurn("completed", "root", "null"), betaResultIdle(), betaResultTurn("created", "later", "null"), agentEvent("failed", "later-failure", `,"session":{"id":"session","status":"failed"}`), `{"type":"error","error":{"message":"later delivery error"}}`}
	stream := betaResultCreateStream(t, events, false)
	for stream.Next() {
	}
	if stream.Err() == nil {
		t.Fatal("raw stream must retain its later delivery error")
	}
	result, err := openai.BetaAgentSessionFinalResult(stream)
	if err != nil || result.OutputText() != "answer" || result.TurnID() != "root" {
		t.Fatalf("result=%v err=%v", result, err)
	}
}

func TestBetaAgentFinalResultUnknownAddedThenCommentary(t *testing.T) {
	events := []string{betaResultTurn("created", "root", "null"), betaResultMessage("added", "m", "root", "null", "", 0), betaResultMessage("done", "m", "root", `"commentary"`, "thinking", 0), betaResultTurn("completed", "root", "null"), betaResultIdle()}
	stream := betaResultCreateStream(t, events, true)
	result, err := openai.BetaAgentSessionFinalResult(stream)
	if err != nil || result.OutputText() != "" || len(result.Messages) != 0 {
		t.Fatalf("result=%v err=%v", result, err)
	}
}

func TestBetaAgentFinalResultKnownFailureDoesNotWaitForIdle(t *testing.T) {
	for _, status := range []string{"failed", "cancelled"} {
		t.Run(status, func(t *testing.T) {
			events := []string{betaResultTurn("created", "root", "null"), betaResultTurn(status, "root", "null")}
			stream := betaResultCreateStream(t, events, true)
			_, err := openai.BetaAgentSessionFinalResult(stream)
			var failure *openai.BetaAgentTurnResultError
			if !errors.As(err, &failure) || failure.Reason != "turn_"+status || failure.Cause != nil {
				t.Fatalf("creation waited for a transport failure: %v", err)
			}
			mock, client := newAgentHelperServer(t, events...)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			followup := client.Beta.Agents.Sessions.Stream(ctx, "session", openai.AgentSessionStreamParams{Input: "question"})
			_, err = followup.FinalResult()
			if !errors.As(err, &failure) || failure.Reason != "turn_"+status || failure.Cause != nil || mock.closed.Load() != 1 {
				t.Fatalf("follow-up waited for a transport failure: %v", err)
			}
		})
	}
}

func TestBetaAgentFinalResultUsesCompletedMessageTurn(t *testing.T) {
	done := strings.Replace(betaResultMessage("done", "m", "root", `"final_answer"`, "answer", 0), `"turn_id":"root"`, `"turn_id":null`, 1)
	stream := betaResultCreateStream(t, []string{betaResultTurn("created", "root", "null"), done, betaResultTurn("completed", "root", "null"), betaResultIdle()}, true)
	result, err := openai.BetaAgentSessionFinalResult(stream)
	if err != nil || result.OutputText() != "answer" {
		t.Fatalf("result=%v err=%v", result, err)
	}
}

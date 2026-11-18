package openai_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func budgetEvent(id, sessionID string) string {
	return agentEvent("budget_reached", id, `,"session":{"id":"`+sessionID+`","status":"budget_reached"}`)
}

func TestAgentsHelperBudgetReachedStandalone(t *testing.T) {
	// The stream remains open after the matching budget event. Closing it must
	// neither require a turn nor mistake another session's budget for this one.
	events := []string{budgetEvent("foreign", "other-session"), budgetEvent("matching", "session")}
	mock, client := newAgentHelperServer(t, events...)
	stream := client.Beta.Agents.Sessions.Stream(context.Background(), "session", openai.AgentSessionStreamParams{Input: "Question"})
	defer func() { _ = stream.Close() }()
	for index, want := range []string{"agent.session.budget_reached", "agent.session.budget_reached"} {
		if !stream.Next() || stream.Current().Type != want {
			t.Fatalf("event %d: got %q, error=%v, want %q", index, stream.Current().Type, stream.Err(), want)
		}
		if stream.Current().Session.ID == "" {
			t.Fatal("budget snapshot lost its session identity")
		}
	}
	if stream.Next() || stream.Err() != nil || mock.closed.Load() != 1 {
		t.Fatalf("stream error=%v closed=%d", stream.Err(), mock.closed.Load())
	}
}

func TestBetaAgentBudgetBoundaryResults(t *testing.T) {
	for _, mode := range []string{"creation", "follow-up"} {
		for _, status := range []string{"completed", "failed", "cancelled", "incomplete"} {
			t.Run(mode+"/"+status, func(t *testing.T) {
				var events []string
				if mode == "creation" {
					events = append(events, agentEvent("created", "session-created", `,"session":{"id":"session"}`))
				}
				if status != "incomplete" {
					events = append(events, betaResultTurn("created", "root", "null"), budgetEvent("early", "session"))
					events = append(events, betaResultMessage("done", "answer", "root", `"final_answer"`, "answer", 0), betaResultTurn(status, "root", "null"))
				}
				events = append(events, budgetEvent("foreign-terminal", "other-session"), budgetEvent("matching", "session"))
				var result *openai.BetaAgentTurnResult
				var err error
				var finalEventID string
				if mode == "creation" {
					// A budget event ends observation before the turn starts or
					// after it ends, even if the SSE subscription stays open.
					stream := betaResultCreateStream(t, events, true)
					result, err = openai.BetaAgentSessionFinalResult(stream)
					finalEventID = stream.Current().EventID
				} else {
					_, client := newAgentHelperServer(t, events...)
					stream := client.Beta.Agents.Sessions.Stream(context.Background(), "session", openai.AgentSessionStreamParams{Input: "Question"})
					defer func() { _ = stream.Close() }()
					result, err = stream.FinalResult()
					finalEventID = stream.Current().EventID
				}
				if (status == "completed" || status == "incomplete") && finalEventID != "matching" {
					t.Fatalf("stopped at event %q, want matching-session budget", finalEventID)
				}
				if status == "completed" {
					if err != nil || result == nil || result.OutputText() != "answer" || result.Turn.Status != "completed" {
						t.Fatalf("result=%#v error=%v", result, err)
					}
					return
				}
				want := map[string]string{"failed": "turn_failed", "cancelled": "turn_cancelled", "incomplete": "observation_incomplete"}[status]
				var failure *openai.BetaAgentTurnResultError
				if result != nil || !errors.As(err, &failure) || failure.Reason != want {
					t.Fatalf("result=%#v error=%#v, want %s", result, err, want)
				}
				if status == "incomplete" && (failure.SessionID != "session" || failure.Cause != nil) {
					t.Fatalf("budget pause result reported missing session or transport failure: %#v", failure)
				}
			})
		}
	}
}

func TestBetaAgentBudgetBoundaryWithoutRoot(t *testing.T) {
	for _, mode := range []string{"budget-first", "error-before-budget"} {
		t.Run(mode, func(t *testing.T) {
			events := []string{}
			if mode == "error-before-budget" {
				events = append(events, `{"type":"error","event_id":"error","session_id":"session","error":{"code":"session_budget_exceeded","message":"Synthetic budget reached"}}`)
			}
			events = append(events, budgetEvent("matching", "session"))
			stream := betaResultCreateStream(t, events, true)
			result, err := openai.BetaAgentSessionFinalResult(stream)
			var failure *openai.BetaAgentTurnResultError
			if result != nil || !errors.As(err, &failure) || failure.Reason != "observation_incomplete" || failure.SessionID != "session" || failure.Cause != nil || failure.Turn != nil {
				t.Fatalf("result=%#v error=%#v", result, err)
			}
		})
	}
}

func TestAgentsCreationToolsBudgetWithoutCreated(t *testing.T) {
	for _, missingID := range []bool{false, true} {
		t.Run(fmt.Sprintf("missing-id=%v", missingID), func(t *testing.T) {

			sessionID := "session"
			if missingID {
				sessionID = ""
			}
			payloads := []string{
				`{"type":"error","event_id":"error","session_id":"session","error":{"code":"session_budget_exceeded","message":"Synthetic budget reached"}}`,
				budgetEvent("matching", sessionID),
			}
			wantTypes := []string{"error", "agent.session.budget_reached"}
			if missingID {
				payloads = append(payloads, agentEvent("failed", "failed", `,"session":{"id":"session","status":"failed"}`))
				wantTypes = append(wantTypes, "agent.session.failed")
			}
			closed := make(chan struct{})
			server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/event-stream")
				for _, payload := range payloads {
					_, _ = fmt.Fprintf(w, "data: %s\n\n", payload)
					w.(http.Flusher).Flush()
				}
				<-r.Context().Done()
				close(closed)
			}))
			defer server.Close()
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			client := openai.NewClient(option.WithBaseURL(server.URL), option.WithHTTPClient(server.Client()), option.WithAPIKey("synthetic"), option.WithMaxRetries(0))
			stream := client.Beta.Agents.Sessions.NewStreaming(ctx, openai.BetaAgentSessionNewParams{
				Environment: openai.EnvironmentParamUnion{OfParamNone: &openai.EnvironmentParamNone{}},
				Input:       openai.BetaAgentSessionNewParamsInputUnion{OfString: openai.String("Question")},
				ToolHandlers: map[string]openai.AgentToolHandler{"unused": func(context.Context, map[string]any) (any, error) {
					t.Error("unexpected tool call")
					return nil, nil
				}},
			})
			defer func() { _ = stream.Close() }()
			for _, want := range wantTypes {
				if !stream.Next() || stream.Current().Type != want {
					t.Fatalf("type=%q error=%v, want %q", stream.Current().Type, stream.Err(), want)
				}
				if want == "error" {
					if _, ok := stream.Current().AsAny().(openai.AgentSessionErrorEvent); !ok {
						t.Fatalf("error type=%T", stream.Current().AsAny())
					}
				}
			}
			if stream.Current().Session.ID != "session" {
				t.Fatalf("budget session=%q", stream.Current().Session.ID)
			}
			select {
			case <-closed:
			case <-time.After(time.Second):
				t.Fatal("matching budget did not close the live response body")
			}
			if stream.Next() || stream.Err() != nil {
				t.Fatalf("matching budget did not stop cleanly: err=%v", stream.Err())
			}

		})
	}
}

func TestBetaAgentBudgetWithoutSessionIDDoesNotEndCreation(t *testing.T) {
	stream := betaResultCreateStream(t, []string{budgetEvent("missing", ""), agentEvent("failed", "failed", `,"session":{"id":"session","status":"failed"}`)}, true)
	_, err := openai.BetaAgentSessionFinalResult(stream)
	var failure *openai.BetaAgentTurnResultError
	if !errors.As(err, &failure) || failure.Reason != "session_failed" || stream.Current().Type != "agent.session.failed" {
		t.Fatalf("error=%#v last event=%q", err, stream.Current().Type)
	}
}

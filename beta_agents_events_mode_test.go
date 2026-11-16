package openai_test

import (
	"testing"

	openai "github.com/openai/openai-go/v3"
)

func TestBetaAgentCreationResultThroughDomainError(t *testing.T) {
	events := []string{
		betaResultTurn("created", "root", "null"),
		betaResultMessage("done", "answer", "root", `"final_answer"`, "answer", 0),
		betaResultTurn("completed", "root", "null"),
		`{"type":"error","event_id":"domain_error","session_id":"session","error":{"type":"server_error","code":"internal_error","message":"Synthetic domain failure."}}`,
		betaResultIdle(),
	}
	stream := betaResultCreateStream(t, events, true)
	openai.BetaAgentSessionWithResultCollection(stream)
	for i, want := range []string{"agent.session.turn.created", "agent.session.turn.item.done", "agent.session.turn.completed", "error"} {
		if !stream.Next() || stream.Current().Type != want {
			t.Fatalf("event %d: got %q, error=%v, want %q", i, stream.Current().Type, stream.Err(), want)
		}
	}
	result, err := openai.BetaAgentSessionFinalResult(stream)
	if err != nil {
		t.Fatal(err)
	}
	if result.Turn.Status != "completed" || result.OutputText() != "answer" {
		t.Fatalf("result=%#v", result)
	}
	if stream.Next() || stream.Err() != nil {
		t.Fatalf("finished stream error=%v", stream.Err())
	}
}

package openai

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestBetaAgentCollectorRetainsOnlyFinalPayload(t *testing.T) {
	c := &betaAgentTurnCollector{}
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

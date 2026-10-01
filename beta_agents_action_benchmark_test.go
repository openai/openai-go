package openai

import (
	"encoding/json"
	"fmt"
	"testing"
)

func BenchmarkBetaAgentAttachmentActionIndex(b *testing.B) {
	for _, count := range []int{1000, 8000} {
		b.Run(fmt.Sprint(count), func(b *testing.B) {
			events := make([]AgentSessionEventUnion, count)
			for i := range events {
				if err := json.Unmarshal([]byte(fmt.Sprintf(`{"type":"agent.session.turn.item.added","item":{"type":"function_call","turn_id":"root","call_id":"call-%d","name":"unknown","arguments":{}}}`, i)), &events[i]); err != nil {
					b.Fatal(err)
				}
			}
			b.ResetTimer()
			for range b.N {
				stream := &AgentSessionStream{attachment: &betaAgentAttachment{}}
				stream.collector.enable()
				for _, event := range events {
					stream.attachmentActions(event)
				}
			}
		})
	}
}

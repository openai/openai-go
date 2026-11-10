package openai_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
)

func TestCompactionMessageRoles(t *testing.T) {
	for _, role := range []string{"user", "system", "developer", "assistant"} {
		content := "input_text"
		if role == "assistant" {
			content = "output_text"
		}
		raw := []byte(fmt.Sprintf(`{"output":[{"type":"message","id":"msg_1","status":"completed","role":%q,"content":[{"type":%q,"text":"hello"}]},{"type":"compaction","id":"cmp_1","encrypted_content":"record"}]}`, role, content))
		t.Run("GA/"+role, func(t *testing.T) {
			var result responses.CompactedResponse
			if err := json.Unmarshal(raw, &result); err != nil {
				t.Fatal(err)
			}
			var items []responses.ResponseItemUnion = result.Output
			if len(items) != 2 || items[1].AsCompaction().EncryptedContent != "record" {
				t.Fatalf("unexpected compaction output: %+v", items)
			}
			if role == "assistant" {
				if _, ok := items[0].AsAny().(responses.ResponseOutputMessage); !ok {
					t.Errorf("AsAny returned %T", items[0].AsAny())
				}
				return
			}
			msg := items[0].AsMessage()
			if string(msg.Role) != role || len(msg.Content) != 1 || msg.Content[0].Text != "hello" {
				t.Fatalf("unexpected input message: %+v", msg)
			}
			if _, ok := items[0].AsAny().(responses.ResponseInputMessageItem); !ok {
				t.Errorf("AsAny returned %T for role %s", items[0].AsAny(), role)
			}
		})
		t.Run("beta/"+role, func(t *testing.T) {
			var result openai.BetaCompactedResponse
			if err := json.Unmarshal(raw, &result); err != nil {
				t.Fatal(err)
			}
			var items []openai.BetaResponseItemUnion = result.Output
			if len(items) != 2 || items[1].AsCompaction().EncryptedContent != "record" {
				t.Fatalf("unexpected compaction output: %+v", items)
			}
			if role == "assistant" {
				if _, ok := items[0].AsAny().(openai.BetaResponseOutputMessage); !ok {
					t.Errorf("AsAny returned %T", items[0].AsAny())
				}
				return
			}
			msg := items[0].AsMessage()
			if string(msg.Role) != role || len(msg.Content) != 1 || msg.Content[0].Text != "hello" {
				t.Fatalf("unexpected input message: %+v", msg)
			}
			if _, ok := items[0].AsAny().(openai.BetaResponseInputMessageItem); !ok {
				t.Errorf("AsAny returned %T for role %s", items[0].AsAny(), role)
			}
		})
	}
}

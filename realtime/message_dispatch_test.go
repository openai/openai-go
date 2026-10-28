package realtime_test

import (
	"encoding/json"
	"testing"

	"github.com/openai/openai-go/v4/realtime"
)

func TestConversationItemMessageDispatch(t *testing.T) {
	for _, role := range []string{"assistant", "system", "user"} {
		t.Run(role, func(t *testing.T) {
			content := `{"type":"input_text","text":"hello"}`
			if role == "user" {
				content = `{"type":"input_image","image_url":"data:image/png;base64,c3ludGhldGlj"}`
			}
			data := []byte(`{"type":"message","role":"` + role + `","content":[` + content + `]}`)
			var item realtime.ConversationItemUnion
			if err := json.Unmarshal(data, &item); err != nil {
				t.Fatal(err)
			}
			switch message := item.AsAny().(type) {
			case realtime.RealtimeConversationItemAssistantMessage:
				if role != "assistant" {
					t.Fatalf("%s selected assistant", role)
				}
			case realtime.RealtimeConversationItemSystemMessage:
				if role != "system" {
					t.Fatalf("%s selected system", role)
				}
			case realtime.RealtimeConversationItemUserMessage:
				if role != "user" {
					t.Fatalf("%s selected user", role)
				}
				if len(message.Content) != 1 || message.Content[0].ImageURL != "data:image/png;base64,c3ludGhldGlj" {
					t.Fatalf("image content lost: %+v", message.Content)
				}
			default:
				t.Fatalf("%s selected %T", role, message)
			}
		})
	}
}

func TestConversationItemMessageDispatchRejectsUnknownRole(t *testing.T) {
	for _, data := range []string{
		`{"type":"message"}`,
		`{"type":"message","role":null}`,
		`{"type":"message","role":""}`,
		`{"type":"message","role":"future_role"}`,
		`{"type":"message","role":42}`,
	} {
		t.Run(data, func(t *testing.T) {
			var item realtime.ConversationItemUnion
			if err := json.Unmarshal([]byte(data), &item); err != nil {
				return
			}
			if got := item.AsAny(); got != nil {
				t.Fatalf("unexpected variant %T", got)
			}
		})
	}
}

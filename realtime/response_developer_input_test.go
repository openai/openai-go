package realtime_test

import (
	"encoding/json"
	"reflect"
	"testing"

	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/realtime"
)

func TestResponseDeveloperInputKeepsTypedUserAndWireShape(t *testing.T) {
	peer := newRealtimePeer(t, false)
	client := openai.NewClient(option.WithAPIKey("local-test-key"), option.WithBaseURL(peer.server.URL+"/v1"))
	conn, err := client.Realtime.Connect(peer.ctx, realtime.RealtimeConnectionOptions{}, option.WithQuery("model", "fixture"))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Abort()

	user := realtime.ConversationItemParamOfRealtimeUserMessageItem(
		[]realtime.RealtimeConversationItemUserMessageContentParam{
			{Type: "input_text", Text: param.NewOpt("Say hello.")},
		})
	before, err := json.Marshal(user)
	if err != nil {
		t.Fatal(err)
	}
	role, typ, content := user.GetRole(), user.GetType(), user.GetContent()
	// A caller's existing typed slice remains the response input API.
	original := []realtime.ConversationItemUnionParam{user}
	response := realtime.RealtimeResponseCreateParams{
		Input: original, Conversation: realtime.RealtimeResponseCreateParamsConversationNone,
		OutputModalities: []string{"text"},
	}
	response.AddDeveloperInput("Use a \"quoted\" word, then a newline:\nUnicode: café")
	response.AddDeveloperInput("End with a period.")
	response.AddDeveloperInput("")
	if len(original) != 1 || len(response.Input) != 4 {
		t.Fatal("append changed the old slice or lost an input")
	}
	stillUser := response.Input[0]
	after, err := json.Marshal(stillUser)
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != string(before) || !reflect.DeepEqual(role, stillUser.GetRole()) ||
		!reflect.DeepEqual(typ, stillUser.GetType()) || !reflect.DeepEqual(content, stillUser.GetContent()) {
		t.Fatal("existing user caller or getters changed")
	}
	if stillUser.OfRealtimeUserMessageItem == nil || stillUser.OfRealtimeSystemMessageItem != nil {
		t.Fatal("ordinary user input selected a different variant")
	}
	err = conn.Send(peer.ctx, realtime.RealtimeClientEventUnionParam{
		OfResponseCreate: &realtime.ResponseCreateEventParam{Response: response},
	})
	if err != nil {
		t.Fatal(err)
	}
	data := peerValue(t, peer.ctx, peer.messages)
	var event struct {
		Type     string `json:"type"`
		Response struct {
			Conversation     string   `json:"conversation"`
			OutputModalities []string `json:"output_modalities"`
			Input            []struct {
				Type    string          `json:"type"`
				Role    string          `json:"role"`
				ID      json.RawMessage `json:"id"`
				Status  json.RawMessage `json:"status"`
				Content []struct {
					Type string  `json:"type"`
					Text *string `json:"text"`
				} `json:"content"`
			} `json:"input"`
		} `json:"response"`
	}
	if err := json.Unmarshal(data, &event); err != nil {
		t.Fatal(err)
	}
	if event.Type != "response.create" || event.Response.Conversation != "none" ||
		!reflect.DeepEqual(event.Response.OutputModalities, []string{"text"}) || len(event.Response.Input) != 4 {
		t.Fatal("incorrect response frame shape")
	}
	want := []struct{ role, text string }{
		{"user", "Say hello."},
		{"developer", "Use a \"quoted\" word, then a newline:\nUnicode: café"},
		{"developer", "End with a period."},
		{"developer", ""},
	}
	for i, expected := range want {
		item := event.Response.Input[i]
		if item.Type != "message" || item.Role != expected.role || len(item.Content) != 1 ||
			item.Content[0].Type != "input_text" || item.Content[0].Text == nil || *item.Content[0].Text != expected.text ||
			(i > 0 && (len(item.ID) != 0 || len(item.Status) != 0)) {
			t.Errorf("response input %d has incorrect role or content", i)
		}
	}
}

// Existing creation still uses its real typed user variant, on the same transport.
func TestResponseDeveloperInputDoesNotAlterConversationCreation(t *testing.T) {
	peer := newRealtimePeer(t, false)
	client := openai.NewClient(option.WithAPIKey("local-test-key"), option.WithBaseURL(peer.server.URL+"/v1"))
	conn, err := client.Realtime.Connect(peer.ctx, realtime.RealtimeConnectionOptions{}, option.WithQuery("model", "fixture"))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Abort()
	item := realtime.RealtimeConversationItemUserMessageParam{
		Content: []realtime.RealtimeConversationItemUserMessageContentParam{{Type: "input_text", Text: param.NewOpt("Good day.")}},
	}
	if err := conn.Send(peer.ctx, realtime.RealtimeClientEventParamOfConversationItemCreate(item)); err != nil {
		t.Fatal(err)
	}
	var event struct {
		Type string `json:"type"`
		Item struct{ Type, Role string }
	}
	if err := json.Unmarshal(peerValue(t, peer.ctx, peer.messages), &event); err != nil {
		t.Fatal(err)
	}
	if event.Type != "conversation.item.create" || event.Item.Type != "message" || event.Item.Role != "user" {
		t.Fatal("typed conversation creation changed")
	}
}

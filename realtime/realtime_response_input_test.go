package realtime_test

import (
	"encoding/json"
	"reflect"
	"testing"

	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/realtime"
)

func TestRealtimeWebsocketResponseInputReferences(t *testing.T) {
	p := newRealtimePeer(t, false)
	client := openai.NewClient(option.WithBaseURL(p.server.URL+"/v1/"), option.WithAPIKey("fake-key"), option.WithUnsafeAllowHTTP())
	conn, err := client.Realtime.Connect(p.ctx, realtime.RealtimeConnectionOptions{})
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()

	user := realtime.ConversationItemParamOfRealtimeUserMessageItem([]realtime.RealtimeConversationItemUserMessageContentParam{
		{Type: "input_text", Text: openai.String("hello")},
	})
	reference := realtime.ConversationItemParamOfItemReference("item_123")
	for _, tc := range []struct {
		name  string
		input []realtime.ConversationItemUnionParam
		wire  string
	}{
		{
			"reference",
			[]realtime.ConversationItemUnionParam{reference},
			`[{"id":"item_123","type":"item_reference"}]`,
		},
		{
			"user",
			[]realtime.ConversationItemUnionParam{user},
			`[{"type":"message","role":"user","content":[{"type":"input_text","text":"hello"}]}]`,
		},
		{
			"mixed",
			[]realtime.ConversationItemUnionParam{user, reference},
			`[{"type":"message","role":"user","content":[{"type":"input_text","text":"hello"}]},{"id":"item_123","type":"item_reference"}]`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			event := realtime.RealtimeClientEventUnionParam{
				OfResponseCreate: &realtime.ResponseCreateEventParam{
					Response: realtime.RealtimeResponseCreateParams{Conversation: "none", Input: tc.input},
				},
			}
			if err := conn.Send(p.ctx, event); err != nil {
				t.Fatal(err)
			}
			wire := peerValue(t, p.ctx, p.messages)
			expected := `{"type":"response.create","response":{"conversation":"none","input":` + tc.wire + `}}`
			var got, want any
			if err := json.Unmarshal(wire, &got); err != nil {
				t.Fatal(err)
			}
			if err := json.Unmarshal([]byte(expected), &want); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("response.create = %s, want %s", wire, expected)
			}
		})
	}
}

package realtime

import "github.com/openai/openai-go/v3/packages/param"

// AddDeveloperInput appends one developer message with one input_text part to
// this response's Input. It omits the message's optional id and status. An empty
// string is sent as text. Call again to append another message.
//
// This sets explicit response input; it does not add the default conversation.
// Developer messages are accepted only in response.create input, not in
// conversation.item.create. The new item uses param.Override, so its union
// variant fields and getters are unset; existing input variants stay intact.
func (r *RealtimeResponseCreateParams) AddDeveloperInput(text string) {
	r.Input = append(r.Input, param.Override[ConversationItemUnionParam](map[string]any{
		"type": "message",
		"role": "developer",
		"content": []map[string]string{
			{"type": "input_text", "text": text},
		},
	}))
}

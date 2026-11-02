package realtime

import (
	"context"
	"encoding/json"
	"errors"
)

// ResponseInput is a conversation item or existing-item reference for
// [RealtimeConnection.SendResponseCreate].
type ResponseInput struct {
	item        *ConversationItemUnionParam
	referenceID *string
}

// ResponseInputItem includes an ordinary item in a response's input.
func ResponseInputItem(item ConversationItemUnionParam) ResponseInput {
	return ResponseInput{item: &item}
}

// ResponseInputReference refers to an item already in the same Realtime session.
func ResponseInputReference(itemID string) ResponseInput {
	return ResponseInput{referenceID: &itemID}
}

// MarshalJSON encodes the selected response input.
func (input ResponseInput) MarshalJSON() ([]byte, error) {
	if input.item != nil {
		return json.Marshal(input.item)
	}
	if input.referenceID != nil {
		return json.Marshal(struct {
			Type string `json:"type"`
			ID   string `json:"id"`
		}{Type: "item_reference", ID: *input.referenceID})
	}
	return nil, errors.New("realtime: uninitialized response input")
}

// SendResponseCreate sends response.create with typed ordinary items and existing
// item references. inputs replaces only response.input; nil or empty inputs clear
// the response context. Other effective event fields and caller values are kept.
// Admission, cancellation, recovery queues and delivery errors follow Send.
func (c *RealtimeConnection) SendResponseCreate(ctx context.Context, event ResponseCreateEventParam, inputs []ResponseInput) error {
	return c.connection.SendJSON(ctx, responseCreateInput{event: event, inputs: inputs})
}

type responseCreateInput struct {
	event  ResponseCreateEventParam
	inputs []ResponseInput
}

func (request responseCreateInput) MarshalJSON() ([]byte, error) {
	encoded, err := json.Marshal(request.event)
	if err != nil {
		return nil, err
	}
	var event map[string]json.RawMessage
	if err = json.Unmarshal(encoded, &event); err != nil || event == nil {
		return nil, errors.New("realtime: response event must encode an object")
	}
	response := map[string]json.RawMessage{}
	if raw, present := event["response"]; present {
		if err = json.Unmarshal(raw, &response); err != nil || response == nil {
			return nil, errors.New("realtime: response parameters must encode an object")
		}
	}
	inputs := request.inputs
	if inputs == nil {
		inputs = []ResponseInput{}
	}
	response["input"], err = json.Marshal(inputs)
	if err != nil {
		return nil, err
	}
	event["response"], err = json.Marshal(response)
	if err != nil {
		return nil, err
	}
	return json.Marshal(event)
}

package realtime

import (
	"context"
	"encoding/json"
	"errors"
	"maps"

	"github.com/openai/openai-go/v3/packages/param"
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
// Known input fields and plain map[string]any metadata are replaced before
// encoding. Other encoders retain their normal error behavior.
// Admission, cancellation, recovery queues and delivery errors follow Send.
func (c *RealtimeConnection) SendResponseCreate(ctx context.Context, event ResponseCreateEventParam, inputs []ResponseInput) error {
	return c.connection.SendJSON(ctx, responseCreateInput{event: event, inputs: inputs})
}

type responseCreateInput struct {
	event  ResponseCreateEventParam
	inputs []ResponseInput
}

func (request responseCreateInput) MarshalJSON() ([]byte, error) {
	eventTemplate := request.event
	eventTemplate.Response, _ = responseCreateInputParams(eventTemplate.Response)
	if override, ok := eventTemplate.Overrides(); ok {
		switch object := override.(type) {
		case map[string]any:
			if projected, changed := responseCreateInputEventObject(object); changed {
				eventTemplate = param.Override[ResponseCreateEventParam](projected)
			}
		case *map[string]any:
			if object != nil {
				if projected, changed := responseCreateInputEventObject(*object); changed {
					eventTemplate = param.Override[ResponseCreateEventParam](&projected)
				}
			}
		}
	} else if extras := eventTemplate.ExtraFields(); extras != nil {
		if projected, changed := responseCreateInputEventObject(extras); changed {
			eventTemplate.SetExtraFields(projected)
		}
	}
	encoded, err := json.Marshal(eventTemplate)
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

func responseCreateInputEventObject(object map[string]any) (map[string]any, bool) {
	projected, changed := responseCreateInputValue(object["response"])
	if !changed {
		return object, false
	}
	object = maps.Clone(object)
	object["response"] = projected
	return object, true
}

func responseCreateInputValue(value any) (any, bool) {
	switch response := value.(type) {
	case RealtimeResponseCreateParams:
		return responseCreateInputParams(response)
	case *RealtimeResponseCreateParams:
		if response != nil {
			if projected, changed := responseCreateInputParams(*response); changed {
				return &projected, true
			}
		}
	case map[string]any:
		return responseCreateInputObject(response)
	case *map[string]any:
		if response != nil {
			if projected, changed := responseCreateInputObject(*response); changed {
				return &projected, true
			}
		}
	}
	return value, false
}

func responseCreateInputParams(response RealtimeResponseCreateParams) (RealtimeResponseCreateParams, bool) {
	changed := response.Input != nil
	response.Input = nil
	if override, ok := response.Overrides(); ok {
		switch object := override.(type) {
		case map[string]any:
			if projected, edited := responseCreateInputObject(object); edited {
				response = param.Override[RealtimeResponseCreateParams](projected)
				changed = true
			}
		case *map[string]any:
			if object != nil {
				if projected, edited := responseCreateInputObject(*object); edited {
					response = param.Override[RealtimeResponseCreateParams](&projected)
					changed = true
				}
			}
		}
	} else if extras := response.ExtraFields(); extras != nil {
		if projected, edited := responseCreateInputObject(extras); edited {
			response.SetExtraFields(projected)
			changed = true
		}
	}
	return response, changed
}

func responseCreateInputObject(object map[string]any) (map[string]any, bool) {
	if _, exists := object["input"]; !exists {
		return object, false
	}
	object = maps.Clone(object)
	delete(object, "input")
	return object, true
}

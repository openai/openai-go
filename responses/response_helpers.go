package responses

import (
	"encoding/json"
	"strings"

	"github.com/openai/openai-go/v3/packages/param"
)

func responseOutputText(r Response) string {
	var outputText strings.Builder
	for _, item := range r.Output {
		for _, content := range item.Content {
			if content.Type == "output_text" {
				outputText.WriteString(content.Text)
			}
		}
	}
	return outputText.String()
}

// ToInput converts this response's output items to request input while preserving
// every item in its original order.
//
// Use ToInput when manually managing stateless conversation history. It preserves
// reasoning items together with their following messages, including encrypted
// content and unknown fields. Append new input to the returned value and start a
// new request without PreviousResponseID.
//
// ToInput is intended for a completed response received from the API. For a
// streamed response, use the completed response or output-item-done events rather
// than in-progress output-item-added events.
func (r Response) ToInput() ResponseInputParam {
	input := make(ResponseInputParam, len(r.Output))
	for i, item := range r.Output {
		input[i] = param.Override[ResponseInputItemUnionParam](json.RawMessage(item.RawJSON()))
	}
	return input
}

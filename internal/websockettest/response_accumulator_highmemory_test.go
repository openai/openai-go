//go:build highmemory

package websockettest

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/openai/openai-go/v3/responses"
)

func TestResponsesAccumulatorReconstructedLargeResponse(t *testing.T) {
	text := strings.Repeat("x", 65<<20)
	item := map[string]any{"type": "message", "id": "m", "content": []any{map[string]any{"type": "output_text", "text": text, "annotations": []any{}}}}
	frame, err := json.Marshal(map[string]any{"type": "response.output_item.done", "stream_id": "reconstruction", "output_index": 0, "item": item})
	if err != nil {
		t.Fatal(err)
	}
	// Keep reconstruction coverage above 64 MiB without a socket deadline.
	// TestResponsesWebSocketLargeFinalResponse covers the large-frame transport;
	// the other reconstruction fixtures cover a real connection.
	frames := [][]byte{frame, []byte(`{"type":"response.completed","stream_id":"reconstruction","response":{"id":"r","status":"completed"}}`)}
	var acc responses.ResponseAccumulator
	for _, raw := range frames {
		var event responses.ResponsesServerEventUnion
		if decodeErr := json.Unmarshal(raw, &event); decodeErr != nil {
			t.Fatal(decodeErr)
		}
		if addErr := acc.AddEvent(event); addErr != nil {
			t.Fatal(addErr)
		}
	}
	response, err := acc.ReconstructedResponse()
	if err != nil || response.OutputText() != text {
		t.Fatalf("large reconstruction failed: %v", err)
	}
}

package websockettest

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

// TestResponsesLive is an explicit opt-in smoke test using synthetic input.
// Credentials remain in the environment and are never logged or recorded.
func TestResponsesLive(t *testing.T) {
	if os.Getenv("OPENAI_RESPONSES_WEBSOCKET_LIVE") != "1" {
		t.Skip("live WebSocket smoke test is opt-in")
	}
	if os.Getenv("OPENAI_API_KEY") == "" {
		t.Fatal("live WebSocket smoke test requires OPENAI_API_KEY")
	}
	model := os.Getenv("OPENAI_WEBSOCKET_MODEL")
	if model == "" {
		t.Fatal("live WebSocket smoke test requires OPENAI_WEBSOCKET_MODEL")
	}
	client := openai.NewClient(option.WithHeader("X-SDK-Integration-Test", "responses-websocket-go"))
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	connection, err := client.Responses.Connect(ctx, responses.ResponseConnectionOptions{}, option.WithHeader("X-SDK-Connection-Test", "custom-header"))
	if err != nil {
		t.Fatalf("live WebSocket upgrade failed: %T", err)
	}
	defer connection.Abort()
	previous := ""
	for range 2 {
		request := map[string]any{"type": "response.create", "model": model, "input": "Reply with the text websocket-ok.", "stream_id": "smoke", "store": false}
		if previous != "" {
			request["previous_response_id"] = previous
		}
		payload, marshalErr := json.Marshal(request)
		if marshalErr != nil {
			t.Fatal(marshalErr)
		}
		sendFixture(t, ctx, connection.Create, payload)
		response, receiveErr := connection.FinalResponse(ctx)
		if receiveErr != nil {
			t.Fatalf("live WebSocket response failed: %T", receiveErr)
		}
		if response.ID == "" || response.ID == previous || response.Status != "completed" || response.OutputText() == "" {
			t.Fatal("live WebSocket response was not a new completed text response")
		}
		previous = response.ID
	}
	if err = connection.Close(); err != nil {
		t.Fatalf("live WebSocket close failed: %T", err)
	}
}

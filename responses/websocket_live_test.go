//go:build live

package responses_test

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

func requireLiveResponsesWebSocket(t *testing.T) {
	t.Helper()
	if os.Getenv("RUN_LIVE_API_TESTS") != "1" {
		t.Skip("set RUN_LIVE_API_TESTS=1 to run live API tests")
	}
	if os.Getenv("OPENAI_API_KEY") == "" {
		t.Skip("set OPENAI_API_KEY to run live API tests")
	}
}

func liveResponsesModel() openai.ResponsesModel {
	if model := os.Getenv("OPENAI_MODEL"); model != "" {
		return openai.ResponsesModel(model)
	}
	return openai.ResponsesModelGPT5Codex
}

func TestLiveResponseWebSocketSimpleText(t *testing.T) {
	requireLiveResponsesWebSocket(t)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client := openai.NewClient()
	conn, err := client.Responses.ConnectWebSocket(ctx)
	if err != nil {
		t.Fatalf("ConnectWebSocket() error = %v", err)
	}
	defer conn.Close()

	for _, key := range []string{"x-reasoning-included", "openai-model", "x-models-etag"} {
		if value := conn.HandshakeHeader().Get(key); value != "" {
			t.Logf("handshake %s: %s", key, value)
		}
	}

	stream, err := conn.New(ctx, responses.ResponseNewParams{
		Input: responses.ResponseNewParamsInputUnion{OfString: openai.String("Reply with exactly: websocket ok")},
		Model: liveResponsesModel(),
	})
	if err != nil {
		t.Fatalf("conn.New() error = %v", err)
	}
	defer stream.Close()

	for stream.Next() {
		event := stream.Current()
		var raw map[string]any
		sequenceNumber := any(nil)
		if err := json.Unmarshal(stream.CurrentRaw(), &raw); err == nil {
			sequenceNumber = raw["sequence_number"]
		}
		t.Logf("event type=%s sequence_number=%v", event.Type, sequenceNumber)
	}
	if err := stream.Err(); err != nil {
		t.Fatalf("stream.Err() = %v", err)
	}
}

func TestLiveResponseWebSocketBadAPIKeyDialError(t *testing.T) {
	if os.Getenv("RUN_LIVE_API_TESTS") != "1" {
		t.Skip("set RUN_LIVE_API_TESTS=1 to run live API tests")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client := openai.NewClient(option.WithAPIKey("sk-this-key-is-deliberately-invalid"))
	_, err := client.Responses.ConnectWebSocket(ctx)
	if err == nil {
		t.Fatal("ConnectWebSocket() error = nil, want bad-key error")
	}
	var dialErr *responses.WebSocketDialError
	if !errors.As(err, &dialErr) {
		t.Fatalf("ConnectWebSocket() error = %T %[1]v, want *responses.WebSocketDialError", err)
	}
	t.Logf("status=%d request-id=%s body=%s", dialErr.StatusCode, dialErr.Header.Get("x-request-id"), dialErr.Body)
}

func TestLiveResponseWebSocketInvalidModelClassifiesAsAPIError(t *testing.T) {
	requireLiveResponsesWebSocket(t)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client := openai.NewClient()
	conn, err := client.Responses.ConnectWebSocket(ctx)
	if err != nil {
		t.Fatalf("ConnectWebSocket() error = %v", err)
	}
	defer conn.Close()

	stream, err := conn.New(ctx, responses.ResponseNewParams{
		Input: responses.ResponseNewParamsInputUnion{OfString: openai.String("hello")},
		Model: openai.ResponsesModel("definitely-not-a-real-model"),
	})
	if err != nil {
		t.Fatalf("conn.New() error = %v", err)
	}
	defer stream.Close()

	var sawFailed bool
	for stream.Next() {
		event := stream.Current()
		t.Logf("event type=%s raw=%s", event.Type, string(stream.CurrentRaw()))
		if event.Type == "response.failed" {
			sawFailed = true
		}
	}
	if err := stream.Err(); err != nil {
		var wsErr *responses.WebSocketError
		if !errors.As(err, &wsErr) {
			var closeErr *responses.WebSocketCloseError
			if errors.As(err, &closeErr) {
				t.Fatalf("stream.Err() = %T %[1]v, want API/model failure, not close error", err)
			}
			var transportErr *responses.WebSocketTransportError
			if errors.As(err, &transportErr) {
				t.Fatalf("stream.Err() = %T %[1]v, want API/model failure, not transport failure", err)
			}
			t.Fatalf("stream.Err() = %T %[1]v, want *responses.WebSocketError or response.failed event", err)
		}
		t.Logf("websocket api error status=%d code=%s message=%s", wsErr.Status, wsErr.Code, wsErr.Message)
		return
	}
	if !sawFailed {
		t.Fatal("invalid model produced no WebSocketError and no response.failed event")
	}
}

func TestLiveResponseWebSocketHandshakeMetadata(t *testing.T) {
	requireLiveResponsesWebSocket(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client := openai.NewClient()
	conn, err := client.Responses.ConnectWebSocket(ctx)
	if err != nil {
		t.Fatalf("ConnectWebSocket() error = %v", err)
	}
	defer conn.Close()

	for _, key := range []string{"x-reasoning-included", "openai-model", "x-models-etag"} {
		t.Logf("handshake %s: %q", key, conn.HandshakeHeader().Get(key))
	}
}

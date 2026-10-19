package openai_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	wire "github.com/coder/websocket"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func TestBetaResponsesWebSocketRoundTrip(t *testing.T) {
	closed := make(chan struct{})
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer close(closed)
		if r.URL.Path != "/responses" || r.URL.Query().Get("test") != "beta" || r.Header.Get("Authorization") != "Bearer synthetic-override" || r.Header.Get("X-Request-Hook") != "called" {
			t.Errorf("request policy was not applied")
		}
		socket, err := wire.Accept(w, r, nil)
		if err != nil {
			t.Error(err)
			return
		}
		defer func() { _ = socket.CloseNow() }()
		_, body, err := socket.Read(r.Context())
		if err != nil {
			t.Error(err)
			return
		}
		var raw map[string]any
		if err := json.Unmarshal(body, &raw); err != nil {
			t.Error(err)
			return
		}
		if raw["type"] != "response.inject" || raw["response_id"] != "resp_beta" {
			t.Errorf("wrong beta command: %s", body)
		}
		for _, frame := range []string{
			`{"type":"response.inject.created","response_id":"resp_beta","sequence_number":7,"stream_id":"stream_beta"}`,
			`{"type":"response.completed","response":{"id":"resp_beta","status":"completed"},"sequence_number":8}`,
			`{"type":"future.beta.event","value":42}`,
		} {
			if err := socket.Write(r.Context(), wire.MessageText, []byte(frame)); err != nil {
				t.Error(err)
				return
			}
		}
		_, _, _ = socket.Read(r.Context())
	}))
	defer server.Close()
	client := openai.NewClient(option.WithAPIKey("synthetic-default"), option.WithBaseURL(server.URL), option.WithHTTPClient(server.Client()), option.WithMiddleware(func(r *http.Request, next option.MiddlewareNext) (*http.Response, error) {
		r.Header.Set("X-Request-Hook", "called")
		return next(r)
	}))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	connection, err := client.Beta.Responses.Connect(ctx, openai.BetaResponseConnectionOptions{}, option.WithHeader("Authorization", "Bearer synthetic-override"), option.WithQuery("test", "beta"))
	if err != nil {
		t.Fatal(err)
	}
	defer connection.Abort()
	if err = connection.Send(ctx, openai.BetaResponsesClientEventUnionParam{OfResponseInject: &openai.BetaResponseInjectEventParam{ResponseID: "resp_beta", Input: []openai.BetaResponseInputItemUnionParam{}}}); err != nil {
		t.Fatal(err)
	}
	event, err := connection.Recv(ctx)
	if err != nil {
		t.Fatal(err)
	}
	injected := event.AsBetaResponseInjectCreatedEvent()
	if injected.ResponseID != "resp_beta" || injected.StreamID != "stream_beta" || injected.SequenceNumber != 7 {
		t.Fatalf("lost beta event: %s", event.RawJSON())
	}
	event, err = connection.Recv(ctx)
	if err != nil || event.Type != "response.completed" {
		t.Fatalf("terminal event: %s, %v", event.Type, err)
	}
	raw, err := connection.RecvRaw(ctx)
	if err != nil || string(raw) != `{"type":"future.beta.event","value":42}` {
		t.Fatalf("socket closed at response completion: %s, %v", raw, err)
	}
	if err := connection.Close(); err != nil {
		t.Fatal(err)
	}
	select {
	case <-closed:
	case <-ctx.Done():
		t.Fatal("socket was not cleaned up")
	}
}

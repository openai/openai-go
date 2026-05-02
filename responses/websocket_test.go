package responses_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

type responsesWebSocketTestServer struct {
	t      *testing.T
	server *httptest.Server

	mu       sync.Mutex
	path     string
	header   http.Header
	messages [][]byte

	accepted chan struct{}
	received chan []byte
	resume   chan struct{}
}

func newResponsesWebSocketTestServer(t *testing.T, handler func(context.Context, *websocket.Conn)) *responsesWebSocketTestServer {
	t.Helper()

	ts := &responsesWebSocketTestServer{
		t:        t,
		accepted: make(chan struct{}),
		received: make(chan []byte, 1),
		resume:   make(chan struct{}),
	}
	ts.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ts.mu.Lock()
		ts.path = r.URL.Path
		ts.header = r.Header.Clone()
		ts.mu.Unlock()

		conn, err := websocket.Accept(w, r, nil)
		if err != nil {
			t.Logf("websocket accept failed: %v", err)
			return
		}
		defer conn.Close(websocket.StatusNormalClosure, "")
		close(ts.accepted)

		if handler != nil {
			handler(r.Context(), conn)
		}
	}))
	t.Cleanup(ts.server.Close)
	return ts
}

func (ts *responsesWebSocketTestServer) recordOneAndWait(ctx context.Context, conn *websocket.Conn) {
	_, msg, err := conn.Read(ctx)
	if err != nil {
		ts.t.Logf("websocket read failed: %v", err)
		return
	}
	ts.mu.Lock()
	ts.messages = append(ts.messages, append([]byte(nil), msg...))
	ts.mu.Unlock()
	ts.received <- msg

	select {
	case <-ts.resume:
	case <-ctx.Done():
	}
}

func (ts *responsesWebSocketTestServer) headerValue(key string) string {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	return ts.header.Get(key)
}

func (ts *responsesWebSocketTestServer) requestPath() string {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	return ts.path
}

func (ts *responsesWebSocketTestServer) firstMessage() []byte {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	if len(ts.messages) == 0 {
		return nil
	}
	return append([]byte(nil), ts.messages[0]...)
}

func (ts *responsesWebSocketTestServer) waitAccepted(t *testing.T) {
	t.Helper()
	select {
	case <-ts.accepted:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for websocket accept")
	}
}

func (ts *responsesWebSocketTestServer) waitReceived(t *testing.T) []byte {
	t.Helper()
	select {
	case msg := <-ts.received:
		return msg
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for websocket message")
		return nil
	}
}

func connectTestWebSocket(t *testing.T, ts *responsesWebSocketTestServer, opts ...option.RequestOption) *responses.WebSocketConn {
	t.Helper()
	client := openai.NewClient(append([]option.RequestOption{
		option.WithBaseURL(ts.server.URL),
		option.WithAPIKey("test-key"),
	}, opts...)...)
	conn, err := client.Responses.ConnectWebSocket(context.Background())
	if err != nil {
		t.Fatalf("ConnectWebSocket() error = %v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })
	return conn
}

func newTestResponseStream(t *testing.T, conn *responses.WebSocketConn) *responses.WebSocketStream {
	t.Helper()
	stream, err := conn.New(context.Background(), responses.ResponseNewParams{
		Input: responses.ResponseNewParamsInputUnion{OfString: openai.String("hello")},
		Model: openai.ChatModelGPT4oMini,
	})
	if err != nil {
		t.Fatalf("conn.New() error = %v", err)
	}
	return stream
}

func TestResponseWebSocketURLConversion(t *testing.T) {
	ts := newResponsesWebSocketTestServer(t, func(ctx context.Context, conn *websocket.Conn) {
		<-ctx.Done()
	})

	conn := connectTestWebSocket(t, ts)
	defer conn.Close()

	ts.waitAccepted(t)
	if got, want := ts.requestPath(), "/responses"; got != want {
		t.Fatalf("request path = %q, want %q", got, want)
	}
}

func TestResponseWebSocketHeadersAndAuth(t *testing.T) {
	ts := newResponsesWebSocketTestServer(t, func(ctx context.Context, conn *websocket.Conn) {
		<-ctx.Done()
	})

	client := openai.NewClient(
		option.WithBaseURL(ts.server.URL),
		option.WithAPIKey("test-key"),
		option.WithOrganization("org_123"),
		option.WithProject("proj_123"),
		option.WithHeader("X-Custom-Header", "custom-value"),
	)
	conn, err := client.Responses.ConnectWebSocket(context.Background())
	if err != nil {
		t.Fatalf("ConnectWebSocket() error = %v", err)
	}
	defer conn.Close()

	ts.waitAccepted(t)
	if got, want := ts.headerValue("Authorization"), "Bearer test-key"; got != want {
		t.Fatalf("Authorization = %q, want %q", got, want)
	}
	if got, want := ts.headerValue("OpenAI-Organization"), "org_123"; got != want {
		t.Fatalf("OpenAI-Organization = %q, want %q", got, want)
	}
	if got, want := ts.headerValue("OpenAI-Project"), "proj_123"; got != want {
		t.Fatalf("OpenAI-Project = %q, want %q", got, want)
	}
	if got, want := ts.headerValue("X-Custom-Header"), "custom-value"; got != want {
		t.Fatalf("X-Custom-Header = %q, want %q", got, want)
	}
}

func TestResponseWebSocketNewSendsResponseCreate(t *testing.T) {
	var ts *responsesWebSocketTestServer
	ts = newResponsesWebSocketTestServer(t, func(ctx context.Context, conn *websocket.Conn) {
		ts.recordOneAndWait(ctx, conn)
	})
	conn := connectTestWebSocket(t, ts)

	stream := newTestResponseStream(t, conn)
	defer stream.Close()
	msg := ts.waitReceived(t)

	var payload map[string]any
	if err := json.Unmarshal(msg, &payload); err != nil {
		t.Fatalf("message was not JSON: %v", err)
	}
	if got, want := payload["type"], "response.create"; got != want {
		t.Fatalf("type = %v, want %q", got, want)
	}
	if got, want := payload["model"], string(openai.ChatModelGPT4oMini); got != want {
		t.Fatalf("model = %v, want %q", got, want)
	}
	if got, want := payload["input"], "hello"; got != want {
		t.Fatalf("input = %v, want %q", got, want)
	}
	close(ts.resume)
}

func TestResponseWebSocketNewDoesNotForceStreamTrue(t *testing.T) {
	var ts *responsesWebSocketTestServer
	ts = newResponsesWebSocketTestServer(t, func(ctx context.Context, conn *websocket.Conn) {
		ts.recordOneAndWait(ctx, conn)
	})
	conn := connectTestWebSocket(t, ts)

	stream := newTestResponseStream(t, conn)
	defer stream.Close()
	msg := ts.waitReceived(t)

	var payload map[string]any
	if err := json.Unmarshal(msg, &payload); err != nil {
		t.Fatalf("message was not JSON: %v", err)
	}
	if v, ok := payload["stream"]; ok {
		t.Fatalf("stream was set to %v; want omitted", v)
	}
	close(ts.resume)
}

func TestResponseWebSocketEventsDecode(t *testing.T) {
	ts := newResponsesWebSocketTestServer(t, func(ctx context.Context, conn *websocket.Conn) {
		_, _, err := conn.Read(ctx)
		if err != nil {
			t.Logf("websocket read failed: %v", err)
			return
		}
		event := []byte(`{"type":"response.output_text.delta","sequence_number":1,"item_id":"item_123","output_index":0,"content_index":0,"delta":"hi","logprobs":[]}`)
		if err := conn.Write(ctx, websocket.MessageText, event); err != nil {
			t.Logf("websocket write failed: %v", err)
		}
	})
	conn := connectTestWebSocket(t, ts)

	stream := newTestResponseStream(t, conn)
	if !stream.Next() {
		t.Fatalf("stream.Next() = false, err = %v", stream.Err())
	}
	event := stream.Current()
	if got, want := event.Type, "response.output_text.delta"; got != want {
		t.Fatalf("event.Type = %q, want %q", got, want)
	}
	if got, want := event.Delta, "hi"; got != want {
		t.Fatalf("event.Delta = %q, want %q", got, want)
	}
	if _, ok := event.AsAny().(responses.ResponseTextDeltaEvent); !ok {
		t.Fatalf("event.AsAny() = %T, want responses.ResponseTextDeltaEvent", event.AsAny())
	}
}

func TestResponseWebSocketOneInFlightEnforced(t *testing.T) {
	var ts *responsesWebSocketTestServer
	ts = newResponsesWebSocketTestServer(t, func(ctx context.Context, conn *websocket.Conn) {
		ts.recordOneAndWait(ctx, conn)
	})
	conn := connectTestWebSocket(t, ts)

	stream := newTestResponseStream(t, conn)
	defer stream.Close()
	ts.waitReceived(t)

	_, err := conn.New(context.Background(), responses.ResponseNewParams{
		Input: responses.ResponseNewParamsInputUnion{OfString: openai.String("second")},
		Model: openai.ChatModelGPT4oMini,
	})
	if err == nil {
		t.Fatal("conn.New() while stream active returned nil error")
	}
	if !strings.Contains(err.Error(), "already active") {
		t.Fatalf("conn.New() error = %q, want already active", err.Error())
	}
	close(ts.resume)
}

func TestResponseWebSocketDocumentedError(t *testing.T) {
	ts := newResponsesWebSocketTestServer(t, func(ctx context.Context, conn *websocket.Conn) {
		_, _, err := conn.Read(ctx)
		if err != nil {
			t.Logf("websocket read failed: %v", err)
			return
		}
		errEvent := []byte(`{"type":"error","status":400,"error":{"code":"previous_response_not_found","message":"missing response","param":"previous_response_id"}}`)
		if err := conn.Write(ctx, websocket.MessageText, errEvent); err != nil {
			t.Logf("websocket write failed: %v", err)
		}
	})
	conn := connectTestWebSocket(t, ts)

	stream := newTestResponseStream(t, conn)
	if stream.Next() {
		t.Fatal("stream.Next() = true, want false")
	}
	var wsErr *responses.WebSocketError
	if !errors.As(stream.Err(), &wsErr) {
		t.Fatalf("stream.Err() = %T %[1]v, want *responses.WebSocketError", stream.Err())
	}
	if wsErr.Status != 400 || wsErr.Code != "previous_response_not_found" || wsErr.Param != "previous_response_id" {
		t.Fatalf("unexpected WebSocketError: %+v", wsErr)
	}
	if !strings.Contains(wsErr.Error(), "missing response") {
		t.Fatalf("WebSocketError.Error() = %q, want message", wsErr.Error())
	}
}

func TestResponseWebSocketCloseIdempotent(t *testing.T) {
	ts := newResponsesWebSocketTestServer(t, func(ctx context.Context, conn *websocket.Conn) {
		_, _, _ = conn.Read(ctx)
	})
	conn := connectTestWebSocket(t, ts)

	if err := conn.Close(); err != nil {
		t.Fatalf("first Close() error = %v", err)
	}
	if err := conn.Close(); err != nil {
		t.Fatalf("second Close() error = %v", err)
	}
}

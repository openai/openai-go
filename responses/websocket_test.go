package responses_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
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

func newResponsesWebSocketHeaderTestServer(t *testing.T, header http.Header, handler func(context.Context, *websocket.Conn)) *responsesWebSocketTestServer {
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

		for key, values := range header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
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
	if got := string(stream.CurrentRaw()); !strings.Contains(got, `"response.output_text.delta"`) {
		t.Fatalf("stream.CurrentRaw() = %q, want raw event", got)
	}
	raw := stream.CurrentRaw()
	raw[0] = '!'
	if got := string(stream.CurrentRaw()); !strings.HasPrefix(got, `{"type"`) {
		t.Fatalf("mutating CurrentRaw result changed internal raw payload: %q", got)
	}
}

func TestResponseWebSocketOneInFlightReturnsSentinel(t *testing.T) {
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
	if !errors.Is(err, responses.ErrWebSocketStreamActive) {
		t.Fatalf("conn.New() error = %T %[1]v, want ErrWebSocketStreamActive", err)
	}
	close(ts.resume)
}

func TestResponseWebSocketNewOnClosedConnReturnsSentinel(t *testing.T) {
	ts := newResponsesWebSocketTestServer(t, func(ctx context.Context, conn *websocket.Conn) {
		<-ctx.Done()
	})
	conn := connectTestWebSocket(t, ts)
	_ = conn.Close()

	_, err := conn.New(context.Background(), responses.ResponseNewParams{
		Input: responses.ResponseNewParamsInputUnion{OfString: openai.String("second")},
		Model: openai.ChatModelGPT4oMini,
	})
	if !errors.Is(err, responses.ErrWebSocketConnectionClosed) {
		t.Fatalf("conn.New() error = %T %[1]v, want ErrWebSocketConnectionClosed", err)
	}
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
	if got := string(wsErr.Raw); !strings.Contains(got, "previous_response_not_found") {
		t.Fatalf("WebSocketError.Raw = %q, want raw payload", got)
	}
}

func TestResponseWebSocketDocumentedErrorPreservesDetails(t *testing.T) {
	ts := newResponsesWebSocketTestServer(t, func(ctx context.Context, conn *websocket.Conn) {
		_, _, err := conn.Read(ctx)
		if err != nil {
			t.Logf("websocket read failed: %v", err)
			return
		}
		errEvent := []byte(`{"type":"error","status":429,"headers":{"x-request-id":"req_123","x-ratelimit-reset":["1","2"]},"error":{"code":"rate_limit_exceeded","message":"slow down","param":"input"}}`)
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
	if got, want := wsErr.Status, http.StatusTooManyRequests; got != want {
		t.Fatalf("Status = %d, want %d", got, want)
	}
	if got, want := wsErr.Code, "rate_limit_exceeded"; got != want {
		t.Fatalf("Code = %q, want %q", got, want)
	}
	if got, want := wsErr.Message, "slow down"; got != want {
		t.Fatalf("Message = %q, want %q", got, want)
	}
	if got, want := wsErr.Param, "input"; got != want {
		t.Fatalf("Param = %q, want %q", got, want)
	}
	if got, want := wsErr.Header.Get("x-request-id"), "req_123"; got != want {
		t.Fatalf("Header.Get(x-request-id) = %q, want %q", got, want)
	}
	if got, want := wsErr.Header.Values("x-ratelimit-reset"), []string{"1", "2"}; strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("Header.Values(x-ratelimit-reset) = %v, want %v", got, want)
	}
	if got := string(wsErr.Raw); !strings.Contains(got, `"status":429`) {
		t.Fatalf("Raw = %q, want status payload", got)
	}
}

func TestResponseWebSocketErrorParsesHeaders(t *testing.T) {
	ts := newResponsesWebSocketTestServer(t, func(ctx context.Context, conn *websocket.Conn) {
		_, _, err := conn.Read(ctx)
		if err != nil {
			t.Logf("websocket read failed: %v", err)
			return
		}
		errEvent := []byte(`{"type":"error","headers":{"x-request-id":"req_header"},"error":{"code":"bad","message":"bad"}}`)
		_ = conn.Write(ctx, websocket.MessageText, errEvent)
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
	if got, want := wsErr.Header.Get("x-request-id"), "req_header"; got != want {
		t.Fatalf("Header.Get(x-request-id) = %q, want %q", got, want)
	}
}

func TestResponseWebSocketErrorParsesStatusCodeAlias(t *testing.T) {
	ts := newResponsesWebSocketTestServer(t, func(ctx context.Context, conn *websocket.Conn) {
		_, _, err := conn.Read(ctx)
		if err != nil {
			t.Logf("websocket read failed: %v", err)
			return
		}
		errEvent := []byte(`{"type":"error","status_code":500,"error":{"code":"server_error","message":"try later"}}`)
		_ = conn.Write(ctx, websocket.MessageText, errEvent)
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
	if got, want := wsErr.Status, http.StatusInternalServerError; got != want {
		t.Fatalf("Status = %d, want %d", got, want)
	}
}

func TestResponseWebSocketDialErrorIncludesHandshakeResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Request-ID", "req_123")
		http.Error(w, `{"error":"bad auth"}`, http.StatusUnauthorized)
	}))
	t.Cleanup(server.Close)

	client := openai.NewClient(
		option.WithBaseURL(server.URL),
		option.WithAPIKey("test-key"),
	)
	_, err := client.Responses.ConnectWebSocket(context.Background())
	if err == nil {
		t.Fatal("ConnectWebSocket() error = nil, want error")
	}
	var dialErr *responses.WebSocketDialError
	if !errors.As(err, &dialErr) {
		t.Fatalf("ConnectWebSocket() error = %T %[1]v, want *responses.WebSocketDialError", err)
	}
	if got, want := dialErr.StatusCode, http.StatusUnauthorized; got != want {
		t.Fatalf("StatusCode = %d, want %d", got, want)
	}
	if got, want := dialErr.Header.Get("X-Request-ID"), "req_123"; got != want {
		t.Fatalf("Header.Get(X-Request-ID) = %q, want %q", got, want)
	}
	if !strings.Contains(dialErr.Body, "bad auth") {
		t.Fatalf("Body = %q, want bad auth", dialErr.Body)
	}
	if errors.Unwrap(dialErr) == nil {
		t.Fatal("errors.Unwrap(dialErr) = nil, want original dial error")
	}
}

func TestResponseWebSocketUnexpectedCloseError(t *testing.T) {
	ts := newResponsesWebSocketTestServer(t, func(ctx context.Context, conn *websocket.Conn) {
		_, _, err := conn.Read(ctx)
		if err != nil {
			t.Logf("websocket read failed: %v", err)
			return
		}
		_ = conn.Close(websocket.StatusPolicyViolation, "keepalive timeout")
	})
	conn := connectTestWebSocket(t, ts)

	stream := newTestResponseStream(t, conn)
	if stream.Next() {
		t.Fatal("stream.Next() = true, want false")
	}
	var closeErr *responses.WebSocketCloseError
	if !errors.As(stream.Err(), &closeErr) {
		t.Fatalf("stream.Err() = %T %[1]v, want *responses.WebSocketCloseError", stream.Err())
	}
	if got, want := closeErr.Status, websocket.StatusPolicyViolation; got != want {
		t.Fatalf("Status = %d, want %d", got, want)
	}
	if got, want := closeErr.Reason, "keepalive timeout"; got != want {
		t.Fatalf("Reason = %q, want %q", got, want)
	}
}

func TestResponseWebSocketNormalCloseBeforeTerminalEventIsError(t *testing.T) {
	ts := newResponsesWebSocketTestServer(t, func(ctx context.Context, conn *websocket.Conn) {
		_, _, err := conn.Read(ctx)
		if err != nil {
			t.Logf("websocket read failed: %v", err)
			return
		}
		_ = conn.Close(websocket.StatusNormalClosure, "done early")
	})
	conn := connectTestWebSocket(t, ts)

	stream := newTestResponseStream(t, conn)
	if stream.Next() {
		t.Fatal("stream.Next() = true, want false")
	}
	var closeErr *responses.WebSocketCloseError
	if !errors.As(stream.Err(), &closeErr) {
		t.Fatalf("stream.Err() = %T %[1]v, want *responses.WebSocketCloseError", stream.Err())
	}
	if got, want := closeErr.Status, websocket.StatusNormalClosure; got != want {
		t.Fatalf("Status = %d, want %d", got, want)
	}
	if got, want := closeErr.Reason, "done early"; got != want {
		t.Fatalf("Reason = %q, want %q", got, want)
	}
}

func TestResponseWebSocketCompletedThenCloseIsClean(t *testing.T) {
	ts := newResponsesWebSocketTestServer(t, func(ctx context.Context, conn *websocket.Conn) {
		_, _, err := conn.Read(ctx)
		if err != nil {
			t.Logf("websocket read failed: %v", err)
			return
		}
		event := []byte(`{"type":"response.completed","sequence_number":1,"response":{"id":"resp_123","created_at":123,"model":"gpt-4o-mini","object":"response","output":[],"parallel_tool_calls":true,"tool_choice":"auto","tools":[]}}`)
		if err := conn.Write(ctx, websocket.MessageText, event); err != nil {
			t.Logf("websocket write failed: %v", err)
			return
		}
		_ = conn.Close(websocket.StatusNormalClosure, "")
	})
	conn := connectTestWebSocket(t, ts)

	stream := newTestResponseStream(t, conn)
	if !stream.Next() {
		t.Fatalf("stream.Next() = false, err = %v", stream.Err())
	}
	if got, want := stream.Current().Type, "response.completed"; got != want {
		t.Fatalf("event.Type = %q, want %q", got, want)
	}
	if stream.Next() {
		t.Fatal("second stream.Next() = true, want false")
	}
	if err := stream.Err(); err != nil {
		t.Fatalf("stream.Err() = %v, want nil", err)
	}
}

func TestResponseWebSocketFailedEventIsTerminalAndReleasesStream(t *testing.T) {
	ts := newResponsesWebSocketTestServer(t, func(ctx context.Context, conn *websocket.Conn) {
		if _, _, err := conn.Read(ctx); err != nil {
			t.Logf("websocket read failed: %v", err)
			return
		}
		failed := []byte(`{"type":"response.failed","sequence_number":1,"response":{"id":"resp_failed","created_at":123,"model":"gpt-4o-mini","object":"response","output":[],"parallel_tool_calls":true,"tool_choice":"auto","tools":[],"status":"failed"}}`)
		if err := conn.Write(ctx, websocket.MessageText, failed); err != nil {
			t.Logf("websocket write failed: %v", err)
			return
		}

		if _, _, err := conn.Read(ctx); err != nil {
			t.Logf("second websocket read failed: %v", err)
			return
		}
		completed := []byte(`{"type":"response.completed","sequence_number":2,"response":{"id":"resp_completed","created_at":123,"model":"gpt-4o-mini","object":"response","output":[],"parallel_tool_calls":true,"tool_choice":"auto","tools":[]}}`)
		_ = conn.Write(ctx, websocket.MessageText, completed)
	})
	conn := connectTestWebSocket(t, ts)

	stream := newTestResponseStream(t, conn)
	if !stream.Next() {
		t.Fatalf("stream.Next() = false, err = %v", stream.Err())
	}
	if got, want := stream.Current().Type, "response.failed"; got != want {
		t.Fatalf("event.Type = %q, want %q", got, want)
	}
	if stream.Next() {
		t.Fatal("second stream.Next() = true, want false after terminal response.failed")
	}
	if err := stream.Err(); err != nil {
		t.Fatalf("stream.Err() = %v, want nil for response.failed event", err)
	}

	next := newTestResponseStream(t, conn)
	if !next.Next() {
		t.Fatalf("next stream.Next() = false, err = %v", next.Err())
	}
	if got, want := next.Current().Type, "response.completed"; got != want {
		t.Fatalf("next event.Type = %q, want %q", got, want)
	}
}

func TestResponseWebSocketDropBeforeTerminalEventIsError(t *testing.T) {
	ts := newResponsesWebSocketTestServer(t, func(ctx context.Context, conn *websocket.Conn) {
		_, _, err := conn.Read(ctx)
		if err != nil {
			t.Logf("websocket read failed: %v", err)
			return
		}
		_ = conn.CloseNow()
	})
	conn := connectTestWebSocket(t, ts)

	stream := newTestResponseStream(t, conn)
	if stream.Next() {
		t.Fatal("stream.Next() = true, want false")
	}
	if stream.Err() == nil {
		t.Fatal("stream.Err() = nil, want transport error")
	}
	if !errors.Is(stream.Err(), io.EOF) {
		var transportErr *responses.WebSocketTransportError
		if !errors.As(stream.Err(), &transportErr) {
			t.Fatalf("stream.Err() = %T %[1]v, want io.EOF or *responses.WebSocketTransportError", stream.Err())
		}
	}
}

func TestResponseWebSocketDecodeErrorIncludesRawPayload(t *testing.T) {
	ts := newResponsesWebSocketTestServer(t, func(ctx context.Context, conn *websocket.Conn) {
		_, _, err := conn.Read(ctx)
		if err != nil {
			t.Logf("websocket read failed: %v", err)
			return
		}
		_ = conn.Write(ctx, websocket.MessageText, []byte(`not json`))
	})
	conn := connectTestWebSocket(t, ts)

	stream := newTestResponseStream(t, conn)
	if stream.Next() {
		t.Fatal("stream.Next() = true, want false")
	}
	var decodeErr *responses.WebSocketDecodeError
	if !errors.As(stream.Err(), &decodeErr) {
		t.Fatalf("stream.Err() = %T %[1]v, want *responses.WebSocketDecodeError", stream.Err())
	}
	if got, want := string(decodeErr.Data), "not json"; got != want {
		t.Fatalf("Data = %q, want %q", got, want)
	}
}

func TestResponseWebSocketUnknownEventIsYieldedWithRawPayload(t *testing.T) {
	ts := newResponsesWebSocketTestServer(t, func(ctx context.Context, conn *websocket.Conn) {
		_, _, err := conn.Read(ctx)
		if err != nil {
			t.Logf("websocket read failed: %v", err)
			return
		}
		unknown := []byte(`{"type":"response.future_event","sequence_number":99,"new_field":"preserved"}`)
		if err := conn.Write(ctx, websocket.MessageText, unknown); err != nil {
			t.Logf("websocket write failed: %v", err)
			return
		}
		completed := []byte(`{"type":"response.completed","sequence_number":100,"response":{"id":"resp_123","created_at":123,"model":"gpt-4o-mini","object":"response","output":[],"parallel_tool_calls":true,"tool_choice":"auto","tools":[]}}`)
		_ = conn.Write(ctx, websocket.MessageText, completed)
	})
	conn := connectTestWebSocket(t, ts)

	stream := newTestResponseStream(t, conn)
	if !stream.Next() {
		t.Fatalf("stream.Next() = false, err = %v", stream.Err())
	}
	if got, want := stream.Current().Type, "response.future_event"; got != want {
		t.Fatalf("event.Type = %q, want %q", got, want)
	}
	if got := string(stream.CurrentRaw()); !strings.Contains(got, `"new_field":"preserved"`) {
		t.Fatalf("CurrentRaw() = %q, want unknown event raw payload", got)
	}
	if !stream.Next() {
		t.Fatalf("second stream.Next() = false, err = %v", stream.Err())
	}
	if got, want := stream.Current().Type, "response.completed"; got != want {
		t.Fatalf("event.Type = %q, want %q", got, want)
	}
	if err := stream.Err(); err != nil {
		t.Fatalf("stream.Err() = %v, want nil", err)
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

func TestResponseWebSocketHandshakeHeader(t *testing.T) {
	header := http.Header{}
	header.Set("X-Reasoning-Included", "true")
	header.Set("OpenAI-Model", "gpt-test")
	ts := newResponsesWebSocketHeaderTestServer(t, header, func(ctx context.Context, conn *websocket.Conn) {
		<-ctx.Done()
	})
	conn := connectTestWebSocket(t, ts)

	if got, want := conn.HandshakeHeader().Get("X-Reasoning-Included"), "true"; got != want {
		t.Fatalf("HandshakeHeader().Get(X-Reasoning-Included) = %q, want %q", got, want)
	}
	if got, want := conn.HandshakeHeader().Get("OpenAI-Model"), "gpt-test"; got != want {
		t.Fatalf("HandshakeHeader().Get(OpenAI-Model) = %q, want %q", got, want)
	}
	got := conn.HandshakeHeader()
	got.Set("OpenAI-Model", "mutated")
	if got, want := conn.HandshakeHeader().Get("OpenAI-Model"), "gpt-test"; got != want {
		t.Fatalf("mutating returned header changed internal header: got %q, want %q", got, want)
	}
}

func TestIsWebSocketRetryableError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{name: "nil", err: nil, want: false},
		{name: "connection closed", err: responses.ErrWebSocketConnectionClosed, want: true},
		{name: "stream active", err: responses.ErrWebSocketStreamActive, want: false},
		{name: "context canceled", err: context.Canceled, want: false},
		{name: "context deadline", err: context.DeadlineExceeded, want: false},
		{name: "dial transport failure", err: &responses.WebSocketDialError{Err: io.EOF}, want: true},
		{name: "dial 401", err: &responses.WebSocketDialError{StatusCode: http.StatusUnauthorized}, want: false},
		{name: "dial 429", err: &responses.WebSocketDialError{StatusCode: http.StatusTooManyRequests}, want: true},
		{name: "dial 500", err: &responses.WebSocketDialError{StatusCode: http.StatusInternalServerError}, want: true},
		{name: "close abnormal", err: &responses.WebSocketCloseError{Status: websocket.StatusAbnormalClosure}, want: true},
		{name: "close normal", err: &responses.WebSocketCloseError{Status: websocket.StatusNormalClosure}, want: false},
		{name: "transport EOF", err: &responses.WebSocketTransportError{Op: "read", Err: io.EOF}, want: true},
		{name: "wrapped EOF", err: io.EOF, want: true},
		{name: "api 400", err: &responses.WebSocketError{Status: http.StatusBadRequest}, want: false},
		{name: "api 408", err: &responses.WebSocketError{Status: http.StatusRequestTimeout}, want: true},
		{name: "api 409", err: &responses.WebSocketError{Status: http.StatusConflict}, want: true},
		{name: "api 429", err: &responses.WebSocketError{Status: http.StatusTooManyRequests}, want: true},
		{name: "api 500", err: &responses.WebSocketError{Status: http.StatusInternalServerError}, want: true},
		{name: "api retryable code without status", err: &responses.WebSocketError{Code: "rate_limit_exceeded"}, want: true},
		{name: "api validation code without status", err: &responses.WebSocketError{Code: "invalid_request_error"}, want: false},
		{name: "arbitrary", err: errors.New("boom"), want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := responses.IsWebSocketRetryableError(tt.err); got != tt.want {
				t.Fatalf("IsWebSocketRetryableError(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}

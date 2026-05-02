package responses

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"slices"
	"strings"
	"sync"

	"github.com/coder/websocket"
	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// ConnectWebSocket connects to the Responses API WebSocket endpoint.
func (r *ResponseService) ConnectWebSocket(ctx context.Context, opts ...option.RequestOption) (*WebSocketConn, error) {
	opts = slices.Concat(r.Options, opts)

	cfg, err := requestconfig.NewRequestConfig(ctx, http.MethodGet, "responses", nil, nil, opts...)
	if err != nil {
		return nil, err
	}

	wsURL, err := responsesWebSocketURL(cfg)
	if err != nil {
		return nil, err
	}

	dialCtx := cfg.Request.Context()
	if cfg.RequestTimeout > 0 {
		var cancel context.CancelFunc
		dialCtx, cancel = context.WithTimeout(dialCtx, cfg.RequestTimeout)
		defer cancel()
	}

	conn, _, err := websocket.Dial(dialCtx, wsURL, &websocket.DialOptions{
		HTTPClient: websocketHTTPClient(cfg),
		HTTPHeader: cfg.Request.Header.Clone(),
	})
	if err != nil {
		return nil, err
	}

	return newWebSocketConn(conn), nil
}

func responsesWebSocketURL(cfg *requestconfig.RequestConfig) (string, error) {
	baseURL := cfg.BaseURL
	if baseURL == nil {
		baseURL = cfg.DefaultBaseURL
	}
	if baseURL == nil {
		return "", errors.New("requestconfig: base url is not set")
	}

	u, err := baseURL.Parse(strings.TrimLeft(cfg.Request.URL.String(), "/"))
	if err != nil {
		return "", err
	}

	switch u.Scheme {
	case "https":
		u.Scheme = "wss"
	case "http":
		u.Scheme = "ws"
	case "wss", "ws":
	default:
		return "", fmt.Errorf("responses websocket: unsupported base URL scheme %q", u.Scheme)
	}

	return u.String(), nil
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func websocketHTTPClient(cfg *requestconfig.RequestConfig) *http.Client {
	if cfg.CustomHTTPDoer == nil && len(cfg.Middlewares) == 0 {
		return cfg.HTTPClient
	}

	handler := cfg.HTTPClient.Do
	if cfg.CustomHTTPDoer != nil {
		handler = cfg.CustomHTTPDoer.Do
	}
	for i := len(cfg.Middlewares) - 1; i >= 0; i-- {
		middleware := cfg.Middlewares[i]
		next := handler
		handler = func(req *http.Request) (*http.Response, error) {
			return middleware(req, next)
		}
	}

	return &http.Client{Transport: roundTripperFunc(handler)}
}

type wsMessage struct {
	messageType websocket.MessageType
	data        []byte
	err         error
}

// WebSocketConn is a WebSocket connection to the Responses API.
type WebSocketConn struct {
	conn *websocket.Conn

	msgCh   chan wsMessage
	writeMu sync.Mutex

	mu       sync.Mutex
	inFlight bool
	closed   bool

	closeOnce sync.Once
	closeErr  error
}

func newWebSocketConn(conn *websocket.Conn) *WebSocketConn {
	c := &WebSocketConn{
		conn:  conn,
		msgCh: make(chan wsMessage, 16),
	}
	go c.readPump()
	return c
}

func (c *WebSocketConn) readPump() {
	defer close(c.msgCh)
	for {
		messageType, data, err := c.conn.Read(context.Background())
		if err != nil {
			c.markClosed()
			c.msgCh <- wsMessage{err: err}
			return
		}
		c.msgCh <- wsMessage{messageType: messageType, data: data}
	}
}

// New creates a response over the WebSocket connection.
func (c *WebSocketConn) New(ctx context.Context, body ResponseNewParams) (*WebSocketStream, error) {
	if err := c.acquireStream(); err != nil {
		return nil, err
	}

	payload, err := body.MarshalJSON()
	if err != nil {
		c.releaseStream()
		return nil, err
	}
	payload, err = sjson.SetBytes(payload, "type", "response.create")
	if err != nil {
		c.releaseStream()
		return nil, err
	}

	c.writeMu.Lock()
	err = c.conn.Write(ctx, websocket.MessageText, payload)
	c.writeMu.Unlock()
	if err != nil {
		c.releaseStream()
		c.markClosed()
		return nil, err
	}

	return &WebSocketStream{conn: c, ctx: ctx}, nil
}

func (c *WebSocketConn) acquireStream() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return errors.New("responses websocket: connection is closed")
	}
	if c.inFlight {
		return errors.New("responses websocket: another response stream is already active on this connection")
	}
	c.inFlight = true
	return nil
}

func (c *WebSocketConn) releaseStream() {
	c.mu.Lock()
	c.inFlight = false
	c.mu.Unlock()
}

func (c *WebSocketConn) markClosed() {
	c.mu.Lock()
	c.closed = true
	c.inFlight = false
	c.mu.Unlock()
}

func (c *WebSocketConn) isClosed() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.closed
}

// Close closes the WebSocket connection. It is safe to call multiple times.
func (c *WebSocketConn) Close() error {
	c.closeOnce.Do(func() {
		c.markClosed()
		err := c.conn.Close(websocket.StatusNormalClosure, "")
		if errors.Is(err, net.ErrClosed) {
			err = nil
		}
		c.closeErr = err
	})
	return c.closeErr
}

// WebSocketStream is a stream of Responses API events received over a WebSocket.
type WebSocketStream struct {
	conn *WebSocketConn
	ctx  context.Context

	cur ResponseStreamEventUnion
	err error

	done     bool
	released bool

	closeOnce sync.Once
	closeErr  error
}

// Next advances the stream to the next event.
func (s *WebSocketStream) Next() bool {
	if s.err != nil || s.done {
		s.release()
		return false
	}

	for {
		msg, ok := s.nextMessage()
		if !ok {
			return false
		}
		if msg.err != nil {
			s.handleReadError(msg.err)
			return false
		}
		messageType, data := msg.messageType, msg.data
		if messageType != websocket.MessageText && messageType != websocket.MessageBinary {
			continue
		}
		if strings.HasPrefix(string(data), "[DONE]") {
			s.done = true
			s.release()
			return false
		}

		if wsErr := parseWebSocketError(data); wsErr != nil {
			s.err = wsErr
			s.done = true
			s.release()
			return false
		}

		var event ResponseStreamEventUnion
		if err := json.Unmarshal(data, &event); err != nil {
			s.err = fmt.Errorf("responses websocket: error decoding event: %w", err)
			s.done = true
			s.release()
			return false
		}

		s.cur = event
		if isTerminalResponseStreamEvent(event.Type) {
			s.done = true
			s.release()
		}
		return true
	}
}

func (s *WebSocketStream) nextMessage() (wsMessage, bool) {
	select {
	case msg, ok := <-s.conn.msgCh:
		if !ok {
			s.handleReadError(net.ErrClosed)
			return wsMessage{}, false
		}
		return msg, true
	case <-s.ctx.Done():
		s.err = s.ctx.Err()
		s.done = true
		s.release()
		return wsMessage{}, false
	}
}

func (s *WebSocketStream) handleReadError(err error) {
	status := websocket.CloseStatus(err)
	if status == websocket.StatusNormalClosure || status == websocket.StatusGoingAway || status == websocket.StatusNoStatusRcvd || errors.Is(err, net.ErrClosed) || s.conn.isClosed() {
		s.err = nil
	} else {
		s.err = err
	}
	if status != -1 || errors.Is(err, net.ErrClosed) {
		s.conn.markClosed()
	}
	s.done = true
	s.release()
}

func isTerminalResponseStreamEvent(eventType string) bool {
	switch eventType {
	case "response.completed", "response.failed", "response.incomplete", "response.cancelled", "error":
		return true
	default:
		return false
	}
}

// Current returns the current event.
func (s *WebSocketStream) Current() ResponseStreamEventUnion {
	return s.cur
}

// Err returns the stream error, if any.
func (s *WebSocketStream) Err() error {
	return s.err
}

// Close closes the stream. If the stream is still active, it closes the underlying
// WebSocket connection because unread events cannot be safely skipped.
func (s *WebSocketStream) Close() error {
	s.closeOnce.Do(func() {
		if s.done {
			s.release()
			return
		}
		s.done = true
		s.release()
		s.closeErr = s.conn.Close()
	})
	return s.closeErr
}

func (s *WebSocketStream) release() {
	if s.released {
		return
	}
	s.released = true
	s.conn.releaseStream()
}

// WebSocketError is returned by WebSocketStream.Err for documented Responses API
// WebSocket error messages.
type WebSocketError struct {
	Status  int
	Code    string
	Message string
	Param   string
}

func (e *WebSocketError) Error() string {
	msg := "responses websocket error"
	if e.Status != 0 {
		msg += fmt.Sprintf(": status %d", e.Status)
	}
	if e.Code != "" {
		msg += fmt.Sprintf(": %s", e.Code)
	}
	if e.Message != "" {
		msg += fmt.Sprintf(": %s", e.Message)
	}
	if e.Param != "" {
		msg += fmt.Sprintf(" (param: %s)", e.Param)
	}
	return msg
}

func parseWebSocketError(data []byte) *WebSocketError {
	if gjson.GetBytes(data, "type").String() != "error" {
		return nil
	}
	errorObject := gjson.GetBytes(data, "error")
	if !errorObject.Exists() || !errorObject.IsObject() {
		return nil
	}

	return &WebSocketError{
		Status:  int(gjson.GetBytes(data, "status").Int()),
		Code:    errorObject.Get("code").String(),
		Message: errorObject.Get("message").String(),
		Param:   errorObject.Get("param").String(),
	}
}

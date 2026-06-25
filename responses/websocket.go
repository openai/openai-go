package responses

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

	conn, resp, err := websocket.Dial(dialCtx, wsURL, &websocket.DialOptions{
		HTTPClient: websocketHTTPClient(cfg),
		HTTPHeader: cfg.Request.Header.Clone(),
	})
	if err != nil {
		return nil, newWebSocketDialError(wsURL, resp, err)
	}

	var header http.Header
	if resp != nil {
		header = resp.Header
	}
	return newWebSocketConn(conn, header), nil
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
	conn   *websocket.Conn
	header http.Header

	msgCh   chan wsMessage
	writeMu sync.Mutex

	mu       sync.Mutex
	inFlight bool
	closed   bool

	closeOnce sync.Once
	closeErr  error
}

var (
	// ErrWebSocketConnectionClosed is returned when a new stream is attempted on
	// a Responses WebSocket connection the SDK already knows is closed.
	ErrWebSocketConnectionClosed = errors.New("responses websocket: connection is closed")

	// ErrWebSocketStreamActive is returned when a new stream is attempted while
	// another response stream is still active on the same WebSocket connection.
	ErrWebSocketStreamActive = errors.New("responses websocket: another response stream is already active on this connection")
)

func newWebSocketConn(conn *websocket.Conn, header http.Header) *WebSocketConn {
	c := &WebSocketConn{
		conn:   conn,
		header: header.Clone(),
		msgCh:  make(chan wsMessage, 16),
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
		return nil, &WebSocketTransportError{Op: "write", Err: err}
	}

	return &WebSocketStream{conn: c, ctx: ctx}, nil
}

func (c *WebSocketConn) acquireStream() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return ErrWebSocketConnectionClosed
	}
	if c.inFlight {
		return ErrWebSocketStreamActive
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

// HandshakeHeader returns the HTTP response headers from the successful
// WebSocket upgrade. The returned header is a copy and may be mutated by the
// caller.
func (c *WebSocketConn) HandshakeHeader() http.Header {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.header.Clone()
}

// WebSocketStream is a stream of Responses API events received over a WebSocket.
type WebSocketStream struct {
	conn *WebSocketConn
	ctx  context.Context

	cur ResponseStreamEventUnion
	raw []byte
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
			s.err = &WebSocketDecodeError{Data: append([]byte(nil), data...), Err: err}
			s.done = true
			s.release()
			return false
		}

		s.cur = event
		s.raw = append(s.raw[:0], data...)
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
	if errors.Is(err, net.ErrClosed) {
		s.err = nil
	} else if status != -1 {
		var closeErr websocket.CloseError
		_ = errors.As(err, &closeErr)
		s.err = &WebSocketCloseError{
			Status: status,
			Reason: closeErr.Reason,
			Err:    err,
		}
	} else {
		s.err = &WebSocketTransportError{Op: "read", Err: err}
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

// CurrentRaw returns the raw JSON payload for the current event. It can be used
// to inspect fields that are not yet modeled by this SDK or to debug ordering
// issues between event types.
func (s *WebSocketStream) CurrentRaw() []byte {
	return append([]byte(nil), s.raw...)
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
	// Type is the top-level event type. It is usually "error".
	Type    string
	Status  int
	Code    string
	Message string
	Param   string
	Header  http.Header
	// Raw is the complete JSON payload received from the server.
	Raw []byte
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
	eventType := gjson.GetBytes(data, "type").String()
	if eventType != "error" {
		return nil
	}
	errorObject := gjson.GetBytes(data, "error")
	if !errorObject.Exists() || !errorObject.IsObject() {
		return nil
	}

	header := make(http.Header)
	headersObject := gjson.GetBytes(data, "headers")
	if headersObject.Exists() && headersObject.IsObject() {
		headersObject.ForEach(func(key, value gjson.Result) bool {
			if value.IsArray() {
				value.ForEach(func(_, item gjson.Result) bool {
					header.Add(key.String(), item.String())
					return true
				})
				return true
			}
			header.Add(key.String(), value.String())
			return true
		})
	}

	status := int(gjson.GetBytes(data, "status").Int())
	if status == 0 {
		status = int(gjson.GetBytes(data, "status_code").Int())
	}

	return &WebSocketError{
		Type:    eventType,
		Status:  status,
		Code:    errorObject.Get("code").String(),
		Message: errorObject.Get("message").String(),
		Param:   errorObject.Get("param").String(),
		Header:  header.Clone(),
		Raw:     append([]byte(nil), data...),
	}
}

// WebSocketDialError is returned by ConnectWebSocket when the WebSocket
// handshake fails. When the server returned an HTTP response, StatusCode,
// Header, and Body contain the handshake response details.
type WebSocketDialError struct {
	URL        string
	StatusCode int
	Header     http.Header
	Body       string
	Err        error
}

func newWebSocketDialError(url string, resp *http.Response, err error) error {
	dialErr := &WebSocketDialError{URL: url, Err: err}
	if resp != nil {
		dialErr.StatusCode = resp.StatusCode
		dialErr.Header = resp.Header.Clone()
		if resp.Body != nil {
			body, readErr := io.ReadAll(resp.Body)
			_ = resp.Body.Close()
			if readErr == nil {
				dialErr.Body = string(body)
			}
		}
	}
	return dialErr
}

func (e *WebSocketDialError) Error() string {
	msg := "responses websocket: dial failed"
	if e.StatusCode != 0 {
		msg += fmt.Sprintf(": status %d", e.StatusCode)
	}
	if e.URL != "" {
		msg += fmt.Sprintf(": %s", e.URL)
	}
	if e.Body != "" {
		msg += fmt.Sprintf(": %s", e.Body)
	}
	if e.Err != nil {
		msg += fmt.Sprintf(": %v", e.Err)
	}
	return msg
}

func (e *WebSocketDialError) Unwrap() error {
	return e.Err
}

// WebSocketCloseError is returned by WebSocketStream.Err when the server closes
// the connection unexpectedly. It exposes the WebSocket close status and reason
// separately from the underlying error string.
type WebSocketCloseError struct {
	Status websocket.StatusCode
	Reason string
	Err    error
}

func (e *WebSocketCloseError) Error() string {
	msg := fmt.Sprintf("responses websocket: closed unexpectedly: status %d", e.Status)
	if e.Reason != "" {
		msg += fmt.Sprintf(": %s", e.Reason)
	}
	if e.Err != nil {
		msg += fmt.Sprintf(": %v", e.Err)
	}
	return msg
}

func (e *WebSocketCloseError) Unwrap() error {
	return e.Err
}

// WebSocketTransportError is returned by WebSocketStream.Err when the
// underlying WebSocket transport fails outside the close-frame path.
type WebSocketTransportError struct {
	Op  string
	Err error
}

func (e *WebSocketTransportError) Error() string {
	if e.Err == nil {
		return "responses websocket: transport error"
	}
	if e.Op == "" {
		return fmt.Sprintf("responses websocket: transport error: %v", e.Err)
	}
	return fmt.Sprintf("responses websocket: transport %s error: %v", e.Op, e.Err)
}

func (e *WebSocketTransportError) Unwrap() error {
	return e.Err
}

// WebSocketDecodeError is returned by WebSocketStream.Err when an event payload
// cannot be decoded into a Responses API stream event. Data contains the raw
// payload for debugging malformed or not-yet-modeled events.
type WebSocketDecodeError struct {
	Data []byte
	Err  error
}

func (e *WebSocketDecodeError) Error() string {
	if e.Err == nil {
		return "responses websocket: error decoding event"
	}
	return fmt.Sprintf("responses websocket: error decoding event: %v: %s", e.Err, string(e.Data))
}

func (e *WebSocketDecodeError) Unwrap() error {
	return e.Err
}

// IsWebSocketRetryableError reports whether err is generally retryable at the
// WebSocket transport/API layer. Callers must still decide whether replaying a
// streamed request is safe for their application, for example based on whether
// any output has already been emitted.
func IsWebSocketRetryableError(err error) bool {
	if err == nil || errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}
	if errors.Is(err, ErrWebSocketStreamActive) {
		return false
	}
	if errors.Is(err, ErrWebSocketConnectionClosed) || errors.Is(err, io.EOF) {
		return true
	}
	var dialErr *WebSocketDialError
	if errors.As(err, &dialErr) {
		if dialErr.StatusCode == 0 {
			return true
		}
		return dialErr.StatusCode == http.StatusRequestTimeout ||
			dialErr.StatusCode == http.StatusConflict ||
			dialErr.StatusCode == http.StatusTooManyRequests ||
			dialErr.StatusCode >= http.StatusInternalServerError
	}
	var transportErr *WebSocketTransportError
	if errors.As(err, &transportErr) {
		return true
	}
	var closeErr *WebSocketCloseError
	if errors.As(err, &closeErr) {
		return closeErr.Status != websocket.StatusNormalClosure
	}
	var wsErr *WebSocketError
	if errors.As(err, &wsErr) {
		switch wsErr.Code {
		case "rate_limit_exceeded", "server_error", "internal_server_error",
			"service_unavailable", "timeout", "request_timeout",
			"connection_limit_exceeded", "too_many_connections":
			return true
		}
		return wsErr.Status == http.StatusRequestTimeout ||
			wsErr.Status == http.StatusConflict ||
			wsErr.Status == http.StatusTooManyRequests ||
			wsErr.Status >= http.StatusInternalServerError
	}
	return false
}

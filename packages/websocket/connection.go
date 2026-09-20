// Package websocket provides ordered WebSocket delivery with configurable limits.
package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	wire "github.com/coder/websocket"
	"github.com/openai/openai-go/v3/internal/websockettransport"
)

var (
	// ErrClosed means the caller closed or aborted the connection.
	ErrClosed = errors.New("websocket: connection closed")
	// ErrBufferLimit means ordered delivery exceeded its configured queue bound.
	ErrBufferLimit = errors.New("websocket: receive buffer limit exceeded")
	// ErrPendingLimit means a send was rejected before serialization or transmission.
	ErrPendingLimit = errors.New("websocket: pending send limit exceeded")
	// ErrMessageLimit means an outgoing message exceeded its configured bound.
	ErrMessageLimit = errors.New("websocket: message limit exceeded")
	// ErrLaneClosed means a lane consumer has been detached.
	ErrLaneClosed = errors.New("websocket: lane detached")
)

// Options applies to new connections only. Zero values select the documented
// defaults. Message and receive buffer limits are opt-in so large Responses
// payloads remain supported.
type Options struct {
	// MaxMessageBytes bounds a complete incoming or outgoing message, including
	// decompressed incoming data. Zero means no byte limit.
	MaxMessageBytes int64
	// MaxBufferedBytes bounds queued incoming data across all lanes. Zero means no byte limit.
	MaxBufferedBytes int64
	// MaxBufferedEvents bounds queued events across all lanes. Zero means no event limit.
	MaxBufferedEvents int
	// MaxPendingSends bounds admitted concurrent writes. Default: 16.
	MaxPendingSends int
	// MaxLanes bounds distinct lane keys over the connection's lifetime, including
	// closed lanes. The default receiver does not count toward this limit. Default: 64.
	MaxLanes int
	// CloseTimeout bounds the closing handshake. Default: five seconds.
	CloseTimeout time.Duration
}

func (o Options) normalized() (Options, error) {
	if o.MaxPendingSends == 0 {
		o.MaxPendingSends = 16
	}
	if o.MaxLanes == 0 {
		o.MaxLanes = 64
	}
	if o.CloseTimeout == 0 {
		o.CloseTimeout = 5 * time.Second
	}
	if o.MaxMessageBytes < 0 || o.MaxBufferedBytes < 0 || o.MaxBufferedEvents < 0 || o.MaxPendingSends < 1 || o.MaxLanes < 1 || o.CloseTimeout < 0 {
		return o, errors.New("websocket: limits and timeout must not be negative")
	}
	return o, nil
}

// DeliveryError reports whether a failed write might have reached the server.
// Such writes must not be replayed automatically.
type DeliveryError struct {
	Cause           error
	MayHaveBeenSent bool
}

func (e *DeliveryError) Error() string {
	return fmt.Sprintf("websocket: send failed (delivery uncertain: %t): %v", e.MayHaveBeenSent, e.Cause)
}
func (e *DeliveryError) Unwrap() error { return e.Cause }

type entry[T any] struct {
	event T
	bytes int64
}
type queue[T any] struct {
	events []entry[T]
	closed bool
}

// upgradeTransport retains the upgraded stream so a close deadline can interrupt
// the WebSocket library's closing handshake, including its pending writes.
type upgradeTransport struct {
	http.RoundTripper
	body *upgradeBody
}

func (t *upgradeTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	response, err := websockettransport.RoundTrip(t.RoundTripper, request)
	if err == nil && response != nil && response.StatusCode == http.StatusSwitchingProtocols {
		if body, ok := response.Body.(io.ReadWriteCloser); ok {
			t.body = &upgradeBody{ReadWriteCloser: body}
			response.Body = t.body
		}
	}
	return response, err
}

type upgradeBody struct {
	io.ReadWriteCloser
	once sync.Once
	err  error
}

func (b *upgradeBody) Close() error {
	b.once.Do(func() { b.err = b.ReadWriteCloser.Close() })
	return b.err
}

// dialError omits URL-bearing wrappers from its display while retaining the
// original error chain for errors.Is/As and structured diagnostics.
type dialError struct{ cause error }

func (e *dialError) Error() string {
	cause := e.cause
	for {
		var urlError *url.Error
		if !errors.As(cause, &urlError) {
			return "websocket: handshake failed: " + cause.Error()
		}
		cause = urlError.Err
	}
}

func (e *dialError) Unwrap() error { return e.cause }

// Connection owns exactly one reader. Recv cancellation affects only the wait.
// Close or Abort ends the connection and unblocks every operation.
type Connection[T any] struct {
	socket          *wire.Conn
	transport       *upgradeBody
	options         Options
	cancel          context.CancelFunc
	readerDone      chan struct{}
	pending         chan struct{}
	writer          chan struct{}
	mu              sync.Mutex
	notify          chan struct{}
	err             error
	queues          map[string]*queue[T]
	bufferedBytes   int64
	bufferedEvents  int
	closeOnce       sync.Once
	closeDone       chan struct{}
	userClosed      bool
	reconnecting    bool
	reconnectCancel context.CancelFunc
}

// Dial upgrades request using client, preserving its transport. decode and route
// run on the owned reader and must not block. ctx bounds only opening.
func Dial[T any](ctx context.Context, request *http.Request, client *http.Client, options Options, decode func([]byte) (T, error), route func(T) string) (*Connection[T], error) {
	if request == nil || request.URL == nil || client == nil || decode == nil || route == nil {
		return nil, errors.New("websocket: request, HTTP client, decoder and router are required")
	}
	options, err := options.normalized()
	if err != nil {
		return nil, err
	}
	if client.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, client.Timeout)
		defer cancel()
	}
	handshakeClient := *client
	handshakeClient.Timeout = 0
	transport := &upgradeTransport{RoundTripper: client.Transport}
	if transport.RoundTripper == nil {
		transport.RoundTripper = http.DefaultTransport
	}
	handshakeClient.Transport = transport
	socket, response, err := dialSocket(ctx, request, &handshakeClient)
	if err != nil {
		if response != nil && response.Body != nil {
			_ = response.Body.Close()
		}
		return nil, &dialError{cause: err}
	}
	if options.MaxMessageBytes == 0 {
		socket.SetReadLimit(-1)
	} else {
		socket.SetReadLimit(options.MaxMessageBytes)
	}
	lifetime, cancel := context.WithCancel(context.Background())
	c := &Connection[T]{socket: socket, transport: transport.body, options: options, cancel: cancel, readerDone: make(chan struct{}), pending: make(chan struct{}, options.MaxPendingSends), writer: make(chan struct{}, 1), notify: make(chan struct{}), queues: map[string]*queue[T]{"": {}}, closeDone: make(chan struct{})}
	go c.read(lifetime, decode, route)
	return c, nil
}

func (c *Connection[T]) signal() { close(c.notify); c.notify = make(chan struct{}) }
func (c *Connection[T]) fail(err error) {
	c.mu.Lock()
	if c.err == nil {
		c.err = err
		c.signal()
	}
	c.mu.Unlock()
	c.cancel()
	_ = c.transport.Close() // Interrupt an in-progress closing handshake as well.
	_ = c.socket.CloseNow() // Best effort: the initiating error remains authoritative.
}

func (c *Connection[T]) read(ctx context.Context, decode func([]byte) (T, error), route func(T) string) {
	defer close(c.readerDone)
	for {
		kind, data, err := c.socket.Read(ctx)
		if err != nil {
			c.fail(err)
			return
		}
		if kind != wire.MessageText {
			c.fail(errors.New("websocket: expected a text message"))
			return
		}
		event, err := decode(data)
		if err != nil {
			c.fail(fmt.Errorf("websocket: decode event: %w", err))
			return
		}
		c.mu.Lock()
		if c.err != nil {
			c.mu.Unlock()
			return
		}
		if (c.options.MaxBufferedEvents > 0 && c.bufferedEvents >= c.options.MaxBufferedEvents) || (c.options.MaxBufferedBytes > 0 && int64(len(data)) > c.options.MaxBufferedBytes-c.bufferedBytes) {
			c.mu.Unlock()
			c.fail(ErrBufferLimit)
			return
		}
		q := c.queues[route(event)]
		if q == nil || q.closed {
			q = c.queues[""]
		}
		q.events = append(q.events, entry[T]{event: event, bytes: int64(len(data))})
		c.bufferedEvents++
		c.bufferedBytes += int64(len(data))
		c.signal()
		c.mu.Unlock()
	}
}

// Send writes one text message. It serializes writers and never retries.
// Cancellation after a write starts aborts the connection and is delivery-uncertain.
func (c *Connection[T]) Send(ctx context.Context, data []byte) error {
	return c.send(ctx, func() ([]byte, error) { return data, nil })
}

// SendJSON admits a write before serializing it, bounding concurrent pending
// serialization as well as network writes.
func (c *Connection[T]) SendJSON(ctx context.Context, value any) error {
	return c.send(ctx, func() ([]byte, error) { return json.Marshal(value) })
}

func (c *Connection[T]) send(ctx context.Context, encode func() ([]byte, error)) error {
	if err := ctx.Err(); err != nil {
		return &DeliveryError{Cause: err}
	}
	if err := c.failure(); err != nil {
		return &DeliveryError{Cause: err}
	}
	select {
	case c.pending <- struct{}{}:
	default:
		return &DeliveryError{Cause: ErrPendingLimit}
	}
	defer func() { <-c.pending }()
	data, err := encode()
	if err != nil {
		return &DeliveryError{Cause: err}
	}
	if c.options.MaxMessageBytes > 0 && int64(len(data)) > c.options.MaxMessageBytes {
		return &DeliveryError{Cause: ErrMessageLimit}
	}

	select {
	case c.writer <- struct{}{}:
	case <-ctx.Done():
		return &DeliveryError{Cause: ctx.Err()}
	case <-c.readerDone:
		return &DeliveryError{Cause: c.failure()}
	}
	defer func() { <-c.writer }()
	if err := c.failure(); err != nil {
		return &DeliveryError{Cause: err}
	}
	if err := ctx.Err(); err != nil {
		return &DeliveryError{Cause: err}
	}
	if err := c.socket.Write(ctx, wire.MessageText, data); err != nil {
		c.fail(err)
		return &DeliveryError{Cause: err, MayHaveBeenSent: true}
	}
	return nil
}

func (c *Connection[T]) failure() error { c.mu.Lock(); defer c.mu.Unlock(); return c.err }

// Recv returns the next event not claimed by an attached lane. Without lanes it
// returns every event. Cancellation does not consume an event or stop the reader.
func (c *Connection[T]) Recv(ctx context.Context) (T, error) {
	c.mu.Lock()
	q := c.queues[""]
	c.mu.Unlock()
	return c.recv(ctx, q)
}
func (c *Connection[T]) recv(ctx context.Context, q *queue[T]) (T, error) {
	var zero T
	for {
		c.mu.Lock()
		if err := ctx.Err(); err != nil {
			c.mu.Unlock()
			return zero, err
		}
		if q.closed {
			c.mu.Unlock()
			return zero, ErrLaneClosed
		}
		if len(q.events) != 0 {
			item := q.events[0]
			q.events[0] = entry[T]{}
			q.events = q.events[1:]
			c.bufferedEvents--
			c.bufferedBytes -= item.bytes
			c.mu.Unlock()
			return item.event, nil
		}
		if c.err != nil {
			err := c.err
			c.mu.Unlock()
			return zero, err
		}
		notify := c.notify
		c.mu.Unlock()
		select {
		case <-notify:
		case <-ctx.Done():
			return zero, ctx.Err()
		}
	}
}

// Lane consumes only events carrying its routing key. Attach it before sending
// the corresponding request. A lane has no reader goroutine of its own.
type Lane[T any] struct {
	connection *Connection[T]
	queue      *queue[T]
}

// Lane attaches a consumer for a nonempty routing key. Keys remain reserved for
// this connection's lifetime, including after Close. Reuse an open lane to receive
// sequential responses on the same key; duplicate attachments fail.
func (c *Connection[T]) Lane(key string) (*Lane[T], error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.err != nil {
		return nil, c.err
	}
	if key == "" {
		return nil, errors.New("websocket: lane key must not be empty")
	}
	if c.queues[key] != nil {
		return nil, errors.New("websocket: lane key already used on this connection")
	}
	if len(c.queues)-1 >= c.options.MaxLanes {
		return nil, errors.New("websocket: lane limit exceeded")
	}
	q := &queue[T]{}
	c.queues[key] = q
	return &Lane[T]{connection: c, queue: q}, nil
}

// Recv returns this lane's next event; canceling ctx leaves the lane attached.
func (l *Lane[T]) Recv(ctx context.Context) (T, error) { return l.connection.recv(ctx, l.queue) }

// Close detaches this consumer, discarding its buffered events. It neither
// cancels server work nor closes the connection. Later events use default Recv.
// The key remains reserved because remote work can continue after detachment.
func (l *Lane[T]) Close() {
	c := l.connection
	c.mu.Lock()
	defer c.mu.Unlock()
	if l.queue.closed {
		return
	}
	l.queue.closed = true
	for _, item := range l.queue.events {
		c.bufferedEvents--
		c.bufferedBytes -= item.bytes
	}
	l.queue.events = nil
	c.signal()
}

// Close performs a bounded closing handshake and releases all connection tasks.
// Repeated and concurrent calls are safe.
func (c *Connection[T]) Close() error { return c.close(true) }

func (c *Connection[T]) close(userInitiated bool) error {
	if userInitiated {
		c.mu.Lock()
		c.userClosed = true
		if c.reconnectCancel != nil {
			c.reconnectCancel()
		}
		c.mu.Unlock()
	}
	c.closeOnce.Do(func() {
		c.mu.Lock()
		if c.err == nil {
			c.err = ErrClosed
			c.signal()
		}
		c.mu.Unlock()
		timer := time.AfterFunc(c.options.CloseTimeout, func() { _ = c.transport.Close() })
		_ = c.socket.Close(wire.StatusNormalClosure, "") // Cleanup errors do not replace the initiating failure.
		timer.Stop()
		c.cancel()
		_ = c.socket.CloseNow()
		<-c.readerDone
		c.mu.Lock()
		for _, q := range c.queues {
			q.events = nil
		}
		c.bufferedBytes = 0
		c.bufferedEvents = 0
		c.mu.Unlock()
		close(c.closeDone)
	})
	<-c.closeDone
	return nil
}

// Abort immediately tears down the socket without waiting for a close handshake.
func (c *Connection[T]) Abort() { c.fail(ErrClosed); _ = c.Close() }

// Reconnect explicitly retires this socket and opens a replacement. It never
// replays a message. Closing or aborting this connection while open is running
// cancels the attempt and closes any late replacement. On success the caller
// takes ownership of the returned connection; the old one remains retired.
func (c *Connection[T]) Reconnect(ctx context.Context, open func(context.Context) (*Connection[T], error)) (*Connection[T], error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if open == nil {
		return nil, errors.New("websocket: reconnect opener is required")
	}
	c.mu.Lock()
	if c.userClosed {
		c.mu.Unlock()
		return nil, ErrClosed
	}
	if c.reconnecting {
		c.mu.Unlock()
		return nil, errors.New("websocket: reconnect already in progress")
	}
	ctx, cancel := context.WithCancel(ctx)
	c.reconnecting = true
	c.reconnectCancel = cancel
	c.mu.Unlock()
	defer cancel()
	stopRetirement := context.AfterFunc(ctx, func() { _ = c.transport.Close() })
	_ = c.close(false)
	stopRetirement()
	var next *Connection[T]
	err := ctx.Err()
	if err == nil {
		next, err = open(ctx)
	}
	if next == nil && err == nil {
		err = errors.New("websocket: reconnect opener returned no connection")
	}
	if next != nil && err != nil {
		next.Abort()
		next = nil
	}
	c.mu.Lock()
	closed := c.userClosed
	canceled := ctx.Err()
	c.reconnecting = false
	c.reconnectCancel = nil
	if err == nil && !closed && canceled == nil {
		c.userClosed = true
	}
	c.mu.Unlock()
	if closed || canceled != nil {
		if next != nil {
			next.Abort()
		}
		if closed {
			return nil, ErrClosed
		}
		return nil, canceled
	}
	return next, err
}

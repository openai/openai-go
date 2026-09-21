package responses

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/websocket"
	"net/http"
	"slices"
	"time"
)

// ResponseConnectionOptions configures a Responses WebSocket connection.
// Zero values use the transport's documented defaults.
type ResponseConnectionOptions = websocket.Options

// ResponseConnection owns a Responses API socket. Terminal response events do
// not close it. Use Close or Abort to release its reader and network resources.
type ResponseConnection struct {
	connection *websocket.Connection[responseWebsocketEvent]
	open       func(context.Context) (*ResponseConnection, error)
}

// Connect opens a Responses API WebSocket. ctx and HTTP client timeouts bound
// opening only. The caller owns the connection lifetime through Close or Abort.
func (r *ResponseService) Connect(ctx context.Context, options ResponseConnectionOptions, opts ...option.RequestOption) (*ResponseConnection, error) {
	requestOptions := slices.Concat([]option.RequestOption{requestconfig.WithBearerAuthSecurity()}, r.Options, opts)
	var open func(context.Context) (*ResponseConnection, error)
	open = func(ctx context.Context) (*ResponseConnection, error) {
		cfg, err := requestconfig.NewRequestConfig(ctx, http.MethodGet, "responses", nil, nil, requestOptions...)
		if err != nil {
			return nil, err
		}
		request, client, err := cfg.PrepareWebSocket()
		if err != nil {
			return nil, err
		}
		ctx = request.Context()
		if timeout := cfg.WebSocketOpenTimeout(); timeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, timeout)
			defer cancel()
		}
		connection, err := websocket.Dial(ctx, request, client, options, decodeResponseWebsocketEvent, func(event responseWebsocketEvent) string { return event.streamID })
		if err != nil {
			return nil, err
		}
		return &ResponseConnection{connection: connection, open: open}, nil
	}
	return open(ctx)
}

// Queues retain wire bytes rather than all inline variants of the public union.
// Validate JSON and its event type before queueing so malformed events still
// fail the connection immediately, even when their lane is not being consumed.
type responseWebsocketEvent struct {
	data     []byte
	streamID string
}

func decodeResponseWebsocketEvent(data []byte) (responseWebsocketEvent, error) {
	var event responseWebsocketEvent
	var envelope struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &envelope); err != nil {
		return event, err
	}
	if envelope.Type == "" {
		return event, errors.New("responses websocket event is missing type")
	}
	var route struct {
		StreamID string `json:"stream_id"`
	}
	// Preserve the existing permissive routing of missing or invalid stream IDs.
	_ = json.Unmarshal(data, &route)
	return responseWebsocketEvent{data: data, streamID: route.StreamID}, nil
}

func (e responseWebsocketEvent) decode() (ResponsesServerEventUnion, error) {
	// The union preserves schema mismatches as field metadata; JSON and event
	// type errors have already been rejected by the reader.
	var event ResponsesServerEventUnion
	err := event.unmarshalWebsocketEvent(e.data)
	return event, err
}

// Send writes a typed client event once. Success is a local write, not a server
// acknowledgment. A websocket.DeliveryError identifies uncertain delivery.
func (c *ResponseConnection) Send(ctx context.Context, event ResponsesClientEventUnionParam) error {
	return c.connection.SendJSON(ctx, event)
}

// Create sends a typed response.create command, including routing and
// continuation fields. It does not wait for response completion.
func (c *ResponseConnection) Create(ctx context.Context, event ResponsesClientEventResponseCreateParam) error {
	return c.connection.SendJSON(ctx, event)
}

// Recv receives the next event not claimed by an attached lane. Canceling ctx
// cancels only this wait and leaves the next event available.
func (c *ResponseConnection) Recv(ctx context.Context) (ResponsesServerEventUnion, error) {
	event, err := c.connection.Recv(ctx)
	if err != nil {
		return ResponsesServerEventUnion{}, err
	}
	return event.decode()
}

// Close performs a bounded closing handshake and releases connection resources.
func (c *ResponseConnection) Close() error { return c.connection.Close() }

// Abort immediately closes the connection and cancels pending operations.
func (c *ResponseConnection) Abort() { c.connection.Abort() }

// Reconnect explicitly replaces this connection. It refreshes handshake
// credentials but never replays commands or restores connection-local service
// state. Applications must restore any required state on the returned connection.
func (c *ResponseConnection) Reconnect(ctx context.Context) (*ResponseConnection, error) {
	return c.Recover(ctx, ResponseRecoveryOptions{})
}

// ResponseRecoveryOptions configures explicit, opt-in transport recovery.
// Recovery never restores server caches or replays previous client commands.
type ResponseRecoveryOptions struct {
	// MaxAttempts defaults to one. Every attempt resolves handshake credentials.
	MaxAttempts int
	// RetryDelay is a fixed, cancellable delay between attempts.
	RetryDelay time.Duration
	// Restore runs on a fresh connection before it is returned. It may run once
	// per attempt and must honor its context. The application owns any commands
	// it deliberately sends here, including idempotency and replay decisions.
	Restore func(context.Context, *ResponseConnection) error
}

// Recover replaces the transport with a bounded number of attempts. A failed
// restoration closes that attempted connection. Closing the original connection
// while recovery runs cancels opening, retry delays and restoration context.
func (c *ResponseConnection) Recover(ctx context.Context, options ResponseRecoveryOptions) (*ResponseConnection, error) {
	if options.MaxAttempts == 0 {
		options.MaxAttempts = 1
	}
	if options.MaxAttempts < 1 || options.RetryDelay < 0 {
		return nil, errors.New("websocket: invalid recovery options")
	}
	connection, err := c.connection.Reconnect(ctx, func(ctx context.Context) (*websocket.Connection[responseWebsocketEvent], error) {
		var lastErr error
		for attempt := 0; attempt < options.MaxAttempts; attempt++ {
			if err := ctx.Err(); err != nil {
				return nil, err
			}
			if attempt > 0 && options.RetryDelay > 0 {
				timer := time.NewTimer(options.RetryDelay)
				select {
				case <-timer.C:
				case <-ctx.Done():
					timer.Stop()
					return nil, ctx.Err()
				}
			}
			next, err := c.open(ctx)
			if err != nil {
				lastErr = err
				continue
			}
			if options.Restore != nil {
				if restoreErr := options.Restore(ctx, next); restoreErr != nil {
					next.Abort()
					lastErr = restoreErr
					continue
				}
			}
			return next.connection, nil
		}
		return nil, lastErr
	})
	if err != nil {
		return nil, err
	}
	return &ResponseConnection{connection: connection, open: c.open}, nil
}

// ResponseLane consumes events for one stream_id without a second socket reader.
// Attach before sending the corresponding command.
type ResponseLane struct {
	lane *websocket.Lane[responseWebsocketEvent]
}

// Lane attaches a consumer for a stream_id of 1-256 ASCII letters, digits,
// underscores, hyphens, or periods. The ID remains reserved until this connection
// ends, including after the lane closes. Reuse an open lane for sequential responses.
func (c *ResponseConnection) Lane(streamID string) (*ResponseLane, error) {
	if !validResponseStreamID(streamID) {
		return nil, errors.New("websocket: stream_id must contain 1-256 ASCII letters, digits, underscores, hyphens, or periods")
	}
	lane, err := c.connection.Lane(streamID)
	if err != nil {
		return nil, err
	}
	return &ResponseLane{lane: lane}, nil
}

func validResponseStreamID(streamID string) bool {
	if len(streamID) == 0 || len(streamID) > 256 {
		return false
	}
	for i := range len(streamID) {
		ch := streamID[i]
		if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '_' || ch == '-' || ch == '.') {
			return false
		}
	}
	return true
}

// Recv returns this lane's next event without reading the physical socket itself.
func (l *ResponseLane) Recv(ctx context.Context) (ResponsesServerEventUnion, error) {
	event, err := l.lane.Recv(ctx)
	if err != nil {
		return ResponsesServerEventUnion{}, err
	}
	return event.decode()
}

// Close detaches only this lane; it does not cancel remote work or the socket.
// Its stream_id stays reserved on this connection; later events use connection Recv.
func (l *ResponseLane) Close() { l.lane.Close() }

// ResponseProtocolError retains the complete typed API error event.
type ResponseProtocolError struct{ Event ResponsesServerEventUnion }

func (e *ResponseProtocolError) Error() string { return "responses websocket API error" }

// FinalResponse consumes events until a completed, failed, or incomplete
// response supplies its authoritative final snapshot, including tool and
// multi-part output. EOF is never completion. Use Recv to observe intermediate
// events directly. This consumes unclaimed events, including named streams;
// attached lanes collect their own events. Cancellation leaves the connection usable.
func (c *ResponseConnection) FinalResponse(ctx context.Context) (*Response, error) {
	return finalWebsocketResponse(ctx, c.Recv)
}

// FinalResponse receives this lane's authoritative final response snapshot.
func (l *ResponseLane) FinalResponse(ctx context.Context) (*Response, error) {
	return finalWebsocketResponse(ctx, l.Recv)
}
func finalWebsocketResponse(ctx context.Context, receive func(context.Context) (ResponsesServerEventUnion, error)) (*Response, error) {
	for {
		event, err := receive(ctx)
		if err != nil {
			return nil, err
		}
		// Use the already decoded snapshot without parsing large output again.
		// Keep the raw path below for missing or invalid fields so its errors
		// and permissive handling of other envelope fields remain unchanged.
		switch event.Type {
		case "response.completed":
			if event.OfResponsesServerEventResponseWsCompleted.JSON.Response.Valid() {
				response := event.OfResponsesServerEventResponseWsCompleted.Response
				return &response, nil
			}
		case "response.failed":
			if event.OfResponsesServerEventResponseWsFailed.JSON.Response.Valid() {
				response := event.OfResponsesServerEventResponseWsFailed.Response
				return &response, nil
			}
		case "response.incomplete":
			if event.OfResponsesServerEventResponseWsIncomplete.JSON.Response.Valid() {
				response := event.OfResponsesServerEventResponseWsIncomplete.Response
				return &response, nil
			}
		}
		var envelope struct {
			Type     string    `json:"type"`
			Response *Response `json:"response"`
		}
		if decodeErr := json.Unmarshal([]byte(event.RawJSON()), &envelope); decodeErr != nil {
			return nil, decodeErr
		}
		switch envelope.Type {
		case "response.completed", "response.failed", "response.incomplete":
			if envelope.Response == nil {
				return nil, errors.New("responses terminal event is missing response")
			}
			return envelope.Response, nil
		case "error":
			return nil, &ResponseProtocolError{Event: event}
		}
	}
}

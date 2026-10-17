// File generated from our OpenAPI spec by Castiron. See CONTRIBUTING.md for details.

package live

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/websocket"
	"github.com/tidwall/gjson"
)

// LiveConnectionOptions configures a WebSocket connection.
// Zero values use the transport's documented defaults.
type LiveConnectionOptions = websocket.Options

// LiveConnection owns a WebSocket. Receiving a terminal API event
// does not close it. Use Close or Abort to release its reader and socket.
type LiveConnection struct {
	connection *websocket.Connection[liveWebsocketEvent]
	open       func(context.Context) (*LiveConnection, error)
}

// Connect opens a WebSocket. ctx and HTTP client timeouts bound opening only.
// The caller owns its lifetime through Close or Abort. Use option.WithQuery
// and option.WithQueryAdd for API query parameters.
func (r *LiveService) Connect(ctx context.Context, options LiveConnectionOptions, opts ...option.RequestOption) (*LiveConnection, error) {
	newRequest := r.newWebSocketRequest(opts...)
	var open func(context.Context) (*LiveConnection, error)
	open = func(ctx context.Context) (*LiveConnection, error) {
		cfg, err := newRequest(ctx)
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
		connection, err := websocket.Dial(ctx, request, client, options, decodeLiveWebsocketEvent, func(liveWebsocketEvent) string { return "" })
		if err != nil {
			return nil, err
		}
		return &LiveConnection{connection: connection, open: open}, nil
	}
	return open(ctx)
}

// Queues retain wire bytes instead of every inline variant of the event union.
type liveWebsocketEvent struct{ data []byte }

func decodeLiveWebsocketEvent(data []byte) (liveWebsocketEvent, error) {
	var envelope struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &envelope); err != nil {
		return liveWebsocketEvent{}, err
	}
	// Typed decoding uses the first exact key; keep legacy envelope validation too.
	kind := gjson.GetBytes(data, "type")
	if envelope.Type == "" || (kind.Exists() && (kind.Type != gjson.String || kind.Str == "")) {
		return liveWebsocketEvent{}, errors.New("websocket event is missing type")
	}
	return liveWebsocketEvent{data: data}, nil
}

func (e liveWebsocketEvent) decode() (ServerEventUnion, error) {
	var event ServerEventUnion
	err := event.unmarshalWebsocketEvent(e.data)
	return event, err
}

// Send writes an event once. Success means a local write, not acknowledgment.
// A websocket.DeliveryError identifies uncertain delivery.
func (c *LiveConnection) Send(ctx context.Context, event ClientEventUnionParam) error {
	return c.connection.SendJSON(ctx, event)
}

// Recv receives the next event. Canceling ctx cancels only this wait.
func (c *LiveConnection) Recv(ctx context.Context) (ServerEventUnion, error) {
	event, err := c.connection.Recv(ctx)
	if err != nil {
		var zero ServerEventUnion
		return zero, err
	}
	return event.decode()
}

// Close performs a bounded closing handshake and releases connection resources.
func (c *LiveConnection) Close() error { return c.connection.Close() }

// Abort immediately closes the connection and cancels pending operations.
func (c *LiveConnection) Abort() { c.connection.Abort() }

// Reconnect explicitly replaces this connection and refreshes its credentials.
// It never replays commands or restores API session state.
func (c *LiveConnection) Reconnect(ctx context.Context) (*LiveConnection, error) {
	connection, err := c.connection.Reconnect(ctx, func(ctx context.Context) (*websocket.Connection[liveWebsocketEvent], error) {
		next, err := c.open(ctx)
		if err != nil {
			return nil, err
		}
		return next.connection, nil
	})
	if err != nil {
		return nil, err
	}
	return &LiveConnection{connection: connection, open: c.open}, nil
}

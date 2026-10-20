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

// ForkConnectionOptions configures a WebSocket connection.
// Zero values use the transport's documented defaults.
type ForkConnectionOptions = websocket.Options

// ForkConnection owns a WebSocket. Receiving a terminal API event
// does not close it. Use Close or Abort to release its reader and socket.
type ForkConnection struct {
	connection *websocket.Connection[forkWebsocketEvent]
	open       func(context.Context) (*ForkConnection, error)
}

// Connect opens a WebSocket. ctx and HTTP client timeouts bound opening only.
// The caller owns its lifetime through Close or Abort. Use option.WithQuery
// and option.WithQueryAdd for API query parameters.
func (r *ForkService) Connect(ctx context.Context, sessionID string, options ForkConnectionOptions, opts ...option.RequestOption) (*ForkConnection, error) {
	newRequest := r.newWebSocketRequest(sessionID, opts...)
	var open func(context.Context) (*ForkConnection, error)
	open = func(ctx context.Context) (*ForkConnection, error) {
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
		connection, err := websocket.Dial(ctx, request, client, options, decodeForkWebsocketEvent, func(forkWebsocketEvent) string { return "" })
		if err != nil {
			return nil, err
		}
		return &ForkConnection{connection: connection, open: open}, nil
	}
	return open(ctx)
}

// Queues retain wire bytes instead of every inline variant of the event union.
type forkWebsocketEvent struct{ data []byte }

func decodeForkWebsocketEvent(data []byte) (forkWebsocketEvent, error) {
	var envelope struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &envelope); err != nil {
		return forkWebsocketEvent{}, err
	}
	// Typed decoding uses the first exact key; keep legacy envelope validation too.
	kind := gjson.GetBytes(data, "type")
	if envelope.Type == "" || (kind.Exists() && (kind.Type != gjson.String || kind.Str == "")) {
		return forkWebsocketEvent{}, errors.New("websocket event is missing type")
	}
	return forkWebsocketEvent{data: data}, nil
}

func (e forkWebsocketEvent) decode() (ForkServerEventUnion, error) {
	var event ForkServerEventUnion
	err := event.unmarshalWebsocketEvent(e.data)
	return event, err
}

// Send writes an event once. Success means a local write, not acknowledgment.
// A websocket.DeliveryError identifies uncertain delivery.
func (c *ForkConnection) Send(ctx context.Context, event ForkClientEventUnionParam) error {
	return c.connection.SendJSON(ctx, event)
}

// Recv receives the next event. Canceling ctx cancels only this wait.
func (c *ForkConnection) Recv(ctx context.Context) (ForkServerEventUnion, error) {
	event, err := c.connection.Recv(ctx)
	if err != nil {
		var zero ForkServerEventUnion
		return zero, err
	}
	return event.decode()
}

// SendRaw writes a JSON event once without changing its bytes. Success means a
// local write, not acknowledgment. A websocket.DeliveryError reports uncertainty.
func (c *ForkConnection) SendRaw(ctx context.Context, data json.RawMessage) error {
	return c.connection.Send(ctx, data)
}

// RecvRaw receives the next event not claimed by an attached lane, preserving
// its original JSON. Canceling ctx cancels only this wait. Typed and raw reads
// consume the same ordered event queue.
func (c *ForkConnection) RecvRaw(ctx context.Context) (json.RawMessage, error) {
	event, err := c.connection.Recv(ctx)
	if err != nil {
		return nil, err
	}
	return event.data, nil
}

// Close performs a bounded closing handshake and releases connection resources.
func (c *ForkConnection) Close() error { return c.connection.Close() }

// Abort immediately closes the connection and cancels pending operations.
func (c *ForkConnection) Abort() { c.connection.Abort() }

// Reconnect explicitly replaces this connection and refreshes its credentials.
// It never replays commands or restores API session state.
func (c *ForkConnection) Reconnect(ctx context.Context) (*ForkConnection, error) {
	connection, err := c.connection.Reconnect(ctx, func(ctx context.Context) (*websocket.Connection[forkWebsocketEvent], error) {
		next, err := c.open(ctx)
		if err != nil {
			return nil, err
		}
		return next.connection, nil
	})
	if err != nil {
		return nil, err
	}
	return &ForkConnection{connection: connection, open: c.open}, nil
}

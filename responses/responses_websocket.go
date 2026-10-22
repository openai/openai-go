// File generated from our OpenAPI spec by Castiron. See CONTRIBUTING.md for details.

package responses

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/websocket"
)

// ResponseConnectionOptions configures a WebSocket connection.
// Zero values use the transport's documented defaults.
type ResponseConnectionOptions = websocket.Options

// ResponseConnection owns a WebSocket. Receiving a terminal API event
// does not close it. Use Close or Abort to release its reader and socket.
type ResponseConnection struct {
	connection *websocket.Connection[responseWebsocketEvent]
	open       func(context.Context) (*ResponseConnection, error)
}

// Connect opens a WebSocket. ctx and HTTP client timeouts bound opening only.
// The caller owns its lifetime through Close or Abort. Use option.WithQuery
// and option.WithQueryAdd for API query parameters.
func (r *ResponseService) Connect(ctx context.Context, options ResponseConnectionOptions, opts ...option.RequestOption) (*ResponseConnection, error) {
	newRequest := r.newWebSocketRequest(opts...)
	var open func(context.Context) (*ResponseConnection, error)
	open = func(ctx context.Context) (*ResponseConnection, error) {
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
		connection, err := websocket.Dial(ctx, request, client, options, decodeResponseWebsocketEvent, func(responseWebsocketEvent) string { return "" })
		if err != nil {
			return nil, err
		}
		return &ResponseConnection{connection: connection, open: open}, nil
	}
	return open(ctx)
}

// Queues retain wire bytes instead of every inline variant of the event union.
type responseWebsocketEvent struct{ data []byte }

func decodeResponseWebsocketEvent(data []byte) (responseWebsocketEvent, error) {
	var envelope struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &envelope); err != nil {
		return responseWebsocketEvent{}, err
	}
	if envelope.Type == "" {
		return responseWebsocketEvent{}, errors.New("websocket event is missing type")
	}
	return responseWebsocketEvent{data: data}, nil
}

func (e responseWebsocketEvent) decode() (ResponsesServerEventUnion, error) {
	var event ResponsesServerEventUnion
	err := event.unmarshalWebsocketEvent(e.data)
	return event, err
}

// Send writes an event once. Success means a local write, not acknowledgment.
// A websocket.DeliveryError identifies uncertain delivery.
func (c *ResponseConnection) Send(ctx context.Context, event ResponsesClientEventUnionParam) error {
	return c.connection.SendJSON(ctx, event)
}

// Recv receives the next event. Canceling ctx cancels only this wait.
func (c *ResponseConnection) Recv(ctx context.Context) (ResponsesServerEventUnion, error) {
	event, err := c.connection.Recv(ctx)
	if err != nil {
		var zero ResponsesServerEventUnion
		return zero, err
	}
	return event.decode()
}

// SendRaw writes a JSON event once without changing its bytes. Success means a
// local write, not acknowledgment. A websocket.DeliveryError reports uncertainty.
func (c *ResponseConnection) SendRaw(ctx context.Context, data json.RawMessage) error {
	return c.connection.Send(ctx, data)
}

// RecvRaw receives the next event's original JSON. Canceling ctx cancels only
// this wait. Typed and raw reads consume the same ordered event queue.
func (c *ResponseConnection) RecvRaw(ctx context.Context) (json.RawMessage, error) {
	event, err := c.connection.Recv(ctx)
	if err != nil {
		return nil, err
	}
	return event.data, nil
}

// Close performs a bounded closing handshake and releases connection resources.
func (c *ResponseConnection) Close() error { return c.connection.Close() }

// Abort immediately closes the connection and cancels pending operations.
func (c *ResponseConnection) Abort() { c.connection.Abort() }

// Reconnect explicitly replaces this connection and refreshes its credentials.
// It never replays commands or restores API session state.
func (c *ResponseConnection) Reconnect(ctx context.Context) (*ResponseConnection, error) {
	connection, err := c.connection.Reconnect(ctx, func(ctx context.Context) (*websocket.Connection[responseWebsocketEvent], error) {
		next, err := c.open(ctx)
		if err != nil {
			return nil, err
		}
		return next.connection, nil
	})
	if err != nil {
		return nil, err
	}
	return &ResponseConnection{connection: connection, open: c.open}, nil
}

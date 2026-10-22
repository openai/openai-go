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

// SidebandConnectionOptions configures a WebSocket connection.
// Zero values use the transport's documented defaults.
type SidebandConnectionOptions = websocket.Options

// SidebandConnection owns a WebSocket. Receiving a terminal API event
// does not close it. Use Close or Abort to release its reader and socket.
type SidebandConnection struct {
	connection *websocket.Connection[sidebandWebsocketEvent]
	open       func(context.Context) (*SidebandConnection, error)
}

// Connect opens a WebSocket. ctx and HTTP client timeouts bound opening only.
// The caller owns its lifetime through Close or Abort. Use option.WithQuery
// and option.WithQueryAdd for API query parameters.
func (r *SidebandService) Connect(ctx context.Context, sessionID string, options SidebandConnectionOptions, opts ...option.RequestOption) (*SidebandConnection, error) {
	newRequest := r.newWebSocketRequest(sessionID, opts...)
	var open func(context.Context) (*SidebandConnection, error)
	open = func(ctx context.Context) (*SidebandConnection, error) {
		cfg, err := newRequest(ctx)
		if err != nil {
			return nil, err
		}
		request, client, err := cfg.PrepareWebSocket()
		if err != nil {
			return nil, err
		}
		ctx = request.Context()
		connection, err := websocket.Dial(ctx, request, client, options, decodeSidebandWebsocketEvent, func(sidebandWebsocketEvent) string { return "" })
		if err != nil {
			return nil, err
		}
		return &SidebandConnection{connection: connection, open: open}, nil
	}
	return open(ctx)
}

// Queues retain wire bytes instead of every inline variant of the event union.
type sidebandWebsocketEvent struct{ data []byte }

func decodeSidebandWebsocketEvent(data []byte) (sidebandWebsocketEvent, error) {
	var envelope struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &envelope); err != nil {
		return sidebandWebsocketEvent{}, err
	}
	// Typed decoding uses the first exact key; keep legacy envelope validation too.
	kind := gjson.GetBytes(data, "type")
	if envelope.Type == "" || (kind.Exists() && (kind.Type != gjson.String || kind.Str == "")) {
		return sidebandWebsocketEvent{}, errors.New("websocket event is missing type")
	}
	return sidebandWebsocketEvent{data: data}, nil
}

func (e sidebandWebsocketEvent) decode() (ConnectServerEventUnion, error) {
	var event ConnectServerEventUnion
	err := event.unmarshalWebsocketEvent(e.data)
	return event, err
}

// Send writes an event once. Success means a local write, not acknowledgment.
// A websocket.DeliveryError identifies uncertain delivery.
func (c *SidebandConnection) Send(ctx context.Context, event ConnectClientEventUnionParam) error {
	return c.connection.SendJSON(ctx, event)
}

// Recv receives the next event. Canceling ctx cancels only this wait.
func (c *SidebandConnection) Recv(ctx context.Context) (ConnectServerEventUnion, error) {
	event, err := c.connection.Recv(ctx)
	if err != nil {
		var zero ConnectServerEventUnion
		return zero, err
	}
	return event.decode()
}

// SendRaw writes a JSON event once without changing its bytes. Success means a
// local write, not acknowledgment. A websocket.DeliveryError reports uncertainty.
func (c *SidebandConnection) SendRaw(ctx context.Context, data json.RawMessage) error {
	return c.connection.Send(ctx, data)
}

// RecvRaw receives the next event not claimed by an attached lane, preserving
// its original JSON. Canceling ctx cancels only this wait. Typed and raw reads
// consume the same ordered event queue.
func (c *SidebandConnection) RecvRaw(ctx context.Context) (json.RawMessage, error) {
	event, err := c.connection.Recv(ctx)
	if err != nil {
		return nil, err
	}
	return event.data, nil
}

// Close performs a bounded closing handshake and releases connection resources.
func (c *SidebandConnection) Close() error { return c.connection.Close() }

// Abort immediately closes the connection and cancels pending operations.
func (c *SidebandConnection) Abort() { c.connection.Abort() }

// Reconnect explicitly replaces this connection and refreshes its credentials.
// It never replays commands or restores API session state.
func (c *SidebandConnection) Reconnect(ctx context.Context) (*SidebandConnection, error) {
	connection, err := c.connection.Reconnect(ctx, func(ctx context.Context) (*websocket.Connection[sidebandWebsocketEvent], error) {
		next, err := c.open(ctx)
		if err != nil {
			return nil, err
		}
		return next.connection, nil
	})
	if err != nil {
		return nil, err
	}
	return &SidebandConnection{connection: connection, open: c.open}, nil
}

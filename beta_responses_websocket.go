// File generated from our OpenAPI spec by Castiron. See CONTRIBUTING.md for details.

package openai

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/websocket"
	"github.com/tidwall/gjson"
)

// BetaResponseConnectionOptions configures a WebSocket connection.
// Zero values use the transport's documented defaults.
type BetaResponseConnectionOptions = websocket.Options

// BetaResponseConnection owns a WebSocket. Receiving a terminal API event
// does not close it. Use Close or Abort to release its reader and socket.
type BetaResponseConnection struct {
	connection *websocket.Connection[betaResponseWebsocketEvent]
	open       func(context.Context) (*BetaResponseConnection, error)
}

// Connect opens a WebSocket. ctx and HTTP client timeouts bound opening only.
// The caller owns its lifetime through Close or Abort. Use option.WithQuery
// and option.WithQueryAdd for API query parameters.
func (r *BetaResponseService) Connect(ctx context.Context, options BetaResponseConnectionOptions, opts ...option.RequestOption) (*BetaResponseConnection, error) {
	newRequest := r.newWebSocketRequest(opts...)
	var open func(context.Context) (*BetaResponseConnection, error)
	open = func(ctx context.Context) (*BetaResponseConnection, error) {
		cfg, err := newRequest(ctx)
		if err != nil {
			return nil, err
		}
		request, client, err := cfg.PrepareWebSocket()
		if err != nil {
			return nil, err
		}
		ctx = request.Context()
		connection, err := websocket.Dial(ctx, request, client, options, decodeBetaResponseWebsocketEvent, func(betaResponseWebsocketEvent) string { return "" })
		if err != nil {
			return nil, err
		}
		return &BetaResponseConnection{connection: connection, open: open}, nil
	}
	return open(ctx)
}

// Queues retain wire bytes instead of every inline variant of the event union.
type betaResponseWebsocketEvent struct{ data []byte }

func decodeBetaResponseWebsocketEvent(data []byte) (betaResponseWebsocketEvent, error) {
	var envelope struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &envelope); err != nil {
		return betaResponseWebsocketEvent{}, err
	}
	// Typed decoding uses the first exact key; keep legacy envelope validation too.
	kind := gjson.GetBytes(data, "type")
	if envelope.Type == "" || (kind.Exists() && (kind.Type != gjson.String || kind.Str == "")) {
		return betaResponseWebsocketEvent{}, errors.New("websocket event is missing type")
	}
	return betaResponseWebsocketEvent{data: data}, nil
}

func (e betaResponseWebsocketEvent) decode() (BetaResponsesServerEventUnion, error) {
	var event BetaResponsesServerEventUnion
	err := event.unmarshalWebsocketEvent(e.data)
	return event, err
}

// Send writes an event once. Success means a local write, not acknowledgment.
// A websocket.DeliveryError identifies uncertain delivery.
func (c *BetaResponseConnection) Send(ctx context.Context, event BetaResponsesClientEventUnionParam) error {
	return c.connection.SendJSON(ctx, event)
}

// Recv receives the next event. Canceling ctx cancels only this wait.
func (c *BetaResponseConnection) Recv(ctx context.Context) (BetaResponsesServerEventUnion, error) {
	event, err := c.connection.Recv(ctx)
	if err != nil {
		var zero BetaResponsesServerEventUnion
		return zero, err
	}
	return event.decode()
}

// SendRaw writes a JSON event once without changing its bytes. Success means a
// local write, not acknowledgment. A websocket.DeliveryError reports uncertainty.
func (c *BetaResponseConnection) SendRaw(ctx context.Context, data json.RawMessage) error {
	return c.connection.Send(ctx, data)
}

// RecvRaw receives the next event not claimed by an attached lane, preserving
// its original JSON. Canceling ctx cancels only this wait. Typed and raw reads
// consume the same ordered event queue.
func (c *BetaResponseConnection) RecvRaw(ctx context.Context) (json.RawMessage, error) {
	event, err := c.connection.Recv(ctx)
	if err != nil {
		return nil, err
	}
	return event.data, nil
}

// Close performs a bounded closing handshake and releases connection resources.
func (c *BetaResponseConnection) Close() error { return c.connection.Close() }

// Abort immediately closes the connection and cancels pending operations.
func (c *BetaResponseConnection) Abort() { c.connection.Abort() }

// Reconnect explicitly replaces this connection and refreshes its credentials.
// It never replays commands or restores API session state.
func (c *BetaResponseConnection) Reconnect(ctx context.Context) (*BetaResponseConnection, error) {
	connection, err := c.connection.Reconnect(ctx, func(ctx context.Context) (*websocket.Connection[betaResponseWebsocketEvent], error) {
		next, err := c.open(ctx)
		if err != nil {
			return nil, err
		}
		return next.connection, nil
	})
	if err != nil {
		return nil, err
	}
	return &BetaResponseConnection{connection: connection, open: c.open}, nil
}

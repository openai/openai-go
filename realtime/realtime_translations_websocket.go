// File generated from our OpenAPI spec by Castiron. See CONTRIBUTING.md for details.

package realtime

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/websocket"
	"github.com/tidwall/gjson"
)

// TranslationConnectionOptions configures a WebSocket connection.
// Zero values use the transport's documented defaults.
type TranslationConnectionOptions = websocket.Options

// TranslationConnection owns a WebSocket. Receiving a terminal API event
// does not close it. Use Close or Abort to release its reader and socket.
type TranslationConnection struct {
	connection *websocket.Connection[translationWebsocketEvent]
	open       func(context.Context) (*TranslationConnection, error)
}

// Connect opens a WebSocket. ctx and HTTP client timeouts bound opening only.
// The caller owns its lifetime through Close or Abort. Use option.WithQuery
// and option.WithQueryAdd for API query parameters.
func (r *TranslationService) Connect(ctx context.Context, options TranslationConnectionOptions, opts ...option.RequestOption) (*TranslationConnection, error) {
	newRequest := r.newWebSocketRequest(opts...)
	var open func(context.Context) (*TranslationConnection, error)
	open = func(ctx context.Context) (*TranslationConnection, error) {
		cfg, err := newRequest(ctx)
		if err != nil {
			return nil, err
		}
		request, client, err := cfg.PrepareWebSocket()
		if err != nil {
			return nil, err
		}
		ctx = request.Context()
		connection, err := websocket.Dial(ctx, request, client, options, decodeTranslationWebsocketEvent, func(translationWebsocketEvent) string { return "" })
		if err != nil {
			return nil, err
		}
		return &TranslationConnection{connection: connection, open: open}, nil
	}
	return open(ctx)
}

// Queues retain wire bytes instead of every inline variant of the event union.
type translationWebsocketEvent struct{ data []byte }

func decodeTranslationWebsocketEvent(data []byte) (translationWebsocketEvent, error) {
	var envelope struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &envelope); err != nil {
		return translationWebsocketEvent{}, err
	}
	// Typed decoding uses the first exact key; keep legacy envelope validation too.
	kind := gjson.GetBytes(data, "type")
	if envelope.Type == "" || (kind.Exists() && (kind.Type != gjson.String || kind.Str == "")) {
		return translationWebsocketEvent{}, errors.New("websocket event is missing type")
	}
	return translationWebsocketEvent{data: data}, nil
}

func (e translationWebsocketEvent) decode() (RealtimeTranslationServerEventUnion, error) {
	var event RealtimeTranslationServerEventUnion
	err := event.unmarshalWebsocketEvent(e.data)
	return event, err
}

// Send writes an event once. Success means a local write, not acknowledgment.
// A websocket.DeliveryError identifies uncertain delivery.
func (c *TranslationConnection) Send(ctx context.Context, event RealtimeTranslationClientEventUnionParam) error {
	return c.connection.SendJSON(ctx, event)
}

// Recv receives the next event. Canceling ctx cancels only this wait.
func (c *TranslationConnection) Recv(ctx context.Context) (RealtimeTranslationServerEventUnion, error) {
	event, err := c.connection.Recv(ctx)
	if err != nil {
		var zero RealtimeTranslationServerEventUnion
		return zero, err
	}
	return event.decode()
}

// SendRaw writes a JSON event once without changing its bytes. Success means a
// local write, not acknowledgment. A websocket.DeliveryError reports uncertainty.
func (c *TranslationConnection) SendRaw(ctx context.Context, data json.RawMessage) error {
	return c.connection.Send(ctx, data)
}

// RecvRaw receives the next event not claimed by an attached lane, preserving
// its original JSON. Canceling ctx cancels only this wait. Typed and raw reads
// consume the same ordered event queue.
func (c *TranslationConnection) RecvRaw(ctx context.Context) (json.RawMessage, error) {
	event, err := c.connection.Recv(ctx)
	if err != nil {
		return nil, err
	}
	return event.data, nil
}

// Close performs a bounded closing handshake and releases connection resources.
func (c *TranslationConnection) Close() error { return c.connection.Close() }

// Abort immediately closes the connection and cancels pending operations.
func (c *TranslationConnection) Abort() { c.connection.Abort() }

// Reconnect explicitly replaces this connection and refreshes its credentials.
// It never replays commands or restores API session state.
func (c *TranslationConnection) Reconnect(ctx context.Context) (*TranslationConnection, error) {
	connection, err := c.connection.Reconnect(ctx, func(ctx context.Context) (*websocket.Connection[translationWebsocketEvent], error) {
		next, err := c.open(ctx)
		if err != nil {
			return nil, err
		}
		return next.connection, nil
	})
	if err != nil {
		return nil, err
	}
	return &TranslationConnection{connection: connection, open: c.open}, nil
}

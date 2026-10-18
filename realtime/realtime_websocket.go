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

// RealtimeConnectionOptions configures a WebSocket connection.
// Zero values use the transport's documented defaults.
type RealtimeConnectionOptions = websocket.Options

// RealtimeConnection owns a WebSocket. Receiving a terminal API event
// does not close it. Use Close or Abort to release its reader and socket.
type RealtimeConnection struct {
	connection *websocket.Session[realtimeWebsocketEvent]
}

// Connect opens a WebSocket. ctx and HTTP client timeouts bound opening only.
// The caller owns its lifetime through Close or Abort. Use option.WithQuery
// and option.WithQueryAdd for API query parameters.
func (r *RealtimeService) Connect(ctx context.Context, options RealtimeConnectionOptions, opts ...option.RequestOption) (*RealtimeConnection, error) {
	newRequest := r.newWebSocketRequest(opts...)
	first, firstErr := newRequest(ctx)
	if firstErr != nil {
		return nil, firstErr
	}
	recovery := first.WebSocketRecovery
	open := func(ctx context.Context) (*websocket.Connection[realtimeWebsocketEvent], error) {
		cfg := first
		if cfg != nil {
			first = nil
		} else {
			var err error
			cfg, err = newRequest(ctx)
			if err != nil {
				return nil, err
			}
		}
		request, client, err := cfg.PrepareWebSocketWithRecovery(recovery.Enabled, "intent", "call_id")
		if err != nil {
			return nil, err
		}
		if recovery.Enabled {
			if request.URL.Query().Has("intent") {
				return nil, errors.New("websocket: automatic recovery does not support this session selector")
			}
			if request.URL.Query().Has("call_id") {
				return nil, errors.New("websocket: automatic recovery does not support this session selector")
			}
		}
		ctx = request.Context()
		if timeout := cfg.WebSocketOpenTimeout(); timeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, timeout)
			defer cancel()
		}
		return websocket.Dial(ctx, request, client, options, decodeRealtimeWebsocketEvent, func(realtimeWebsocketEvent) string { return "" })
	}
	connection, err := websocket.OpenSession(ctx, open, recovery)
	if err != nil {
		return nil, err
	}
	return &RealtimeConnection{connection: connection}, nil

}

// Queues retain wire bytes instead of every inline variant of the event union.
type realtimeWebsocketEvent struct{ data []byte }

func decodeRealtimeWebsocketEvent(data []byte) (realtimeWebsocketEvent, error) {
	var envelope struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &envelope); err != nil {
		return realtimeWebsocketEvent{}, err
	}
	// Typed decoding uses the first exact key; keep legacy envelope validation too.
	kind := gjson.GetBytes(data, "type")
	if envelope.Type == "" || (kind.Exists() && (kind.Type != gjson.String || kind.Str == "")) {
		return realtimeWebsocketEvent{}, errors.New("websocket event is missing type")
	}
	return realtimeWebsocketEvent{data: data}, nil
}

func (e realtimeWebsocketEvent) decode() (RealtimeServerEventUnion, error) {
	var event RealtimeServerEventUnion
	err := event.unmarshalWebsocketEvent(e.data)
	return event, err
}

// Send writes an event once. Success means a local write, not acknowledgment.
// A websocket.DeliveryError identifies uncertain delivery.
func (c *RealtimeConnection) Send(ctx context.Context, event RealtimeClientEventUnionParam) error {
	return c.connection.SendJSON(ctx, event)
}

// Recv receives the next event. Canceling ctx cancels only this wait.
func (c *RealtimeConnection) Recv(ctx context.Context) (RealtimeServerEventUnion, error) {
	event, err := c.connection.Recv(ctx)
	if err != nil {
		var zero RealtimeServerEventUnion
		return zero, err
	}
	return event.decode()
}

// SendRaw writes a JSON event once without changing its bytes. Success means a
// local write, not acknowledgment. A websocket.DeliveryError reports uncertainty.
func (c *RealtimeConnection) SendRaw(ctx context.Context, data json.RawMessage) error {
	return c.connection.Send(ctx, data)
}

// RecvRaw receives the next event's original JSON. Canceling ctx cancels only
// this wait. Typed and raw reads consume the same ordered event queue.
func (c *RealtimeConnection) RecvRaw(ctx context.Context) (json.RawMessage, error) {
	event, err := c.connection.Recv(ctx)
	if err != nil {
		return nil, err
	}
	return event.data, nil
}

// Close performs a bounded closing handshake and releases connection resources.
func (c *RealtimeConnection) Close() error { return c.connection.Close() }

// Abort immediately closes the connection and cancels pending operations.
func (c *RealtimeConnection) Abort() { c.connection.Abort() }

// Reconnect explicitly replaces this connection and refreshes its credentials.
// It never replays commands or restores API session state.
func (c *RealtimeConnection) Reconnect(ctx context.Context) (*RealtimeConnection, error) {
	connection, err := c.connection.Reconnect(ctx)
	if err != nil {
		return nil, err
	}
	return &RealtimeConnection{connection: connection}, nil

}

// FlushPending sends never-attempted commands after the caller restores session state.
func (c *RealtimeConnection) FlushPending(ctx context.Context) error {
	return c.connection.FlushPending(ctx)
}

// TakePending removes retained commands without sending them and releases their callers.
func (c *RealtimeConnection) TakePending() []json.RawMessage { return c.connection.TakePending() }

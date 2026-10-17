package realtime

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"slices"
	"sync"

	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/websocket"
	"github.com/tidwall/gjson"
)

var (
	// ErrTranslationReceiving means another Recv, RecvRaw, or Finish owns delivery.
	ErrTranslationReceiving = errors.New("translation: another receiver is active")
	// ErrTranslationInputClosed means protocol or transport closure has begun.
	ErrTranslationInputClosed = errors.New("translation: input is closed")
)

// TranslationSession owns one ordered event queue. Sends may run concurrently,
// but only one Recv, RecvRaw, or Finish may run at a time. The caller must release
// the connection with Finish, Close, or Abort.
type TranslationSession struct {
	connection *TranslationConnection
	mu         sync.Mutex
	receiving  bool
	terminal   bool
	closed     bool
	stopping   bool
	sending    int
	maxSending int
	idle       chan struct{}
	done       chan struct{}
	closeDone  chan struct{}
	closeErr   error
	finishDone bool
	finishErr  error
}

// ConnectTranslation opens /realtime/translations. model replaces any model query
// option. ctx and HTTP timeouts bound opening only. The first session.created
// event remains available to Recv; opening does not wait for session readiness.
func (r *RealtimeService) ConnectTranslation(ctx context.Context, model string, options TranslationConnectionOptions, opts ...option.RequestOption) (*TranslationSession, error) {
	if model == "" {
		return nil, errors.New("translation: model is required")
	}
	requestOptions := append(slices.Clone(opts),
		option.WithQuery("model", model),
		requestconfig.RequestOptionFunc(func(cfg *requestconfig.RequestConfig) error {
			if cfg.EndpointProvider() != "" {
				return errors.New("translation: this provider does not support Translation upgrades")
			}
			return nil
		}),
	)
	translations := r.Translations
	translations.Options = r.Options
	connection, err := translations.Connect(ctx, options, requestOptions...)
	if err != nil {
		return nil, err
	}
	maxSending := options.MaxPendingSends
	if maxSending == 0 {
		maxSending = 16 // Match websocket.Options; bound serialization before inspecting command type.
	}
	return &TranslationSession{connection: connection, maxSending: maxSending, done: make(chan struct{})}, nil
}

type translationSessionEvent struct {
	data []byte
	kind string
}

func decodeTranslationSessionEvent(data []byte) (translationSessionEvent, error) {
	var object struct{}
	if err := json.Unmarshal(data, &object); err != nil {
		return translationSessionEvent{}, err
	}
	kind := gjson.GetBytes(data, "type")
	if kind.Type != gjson.String || kind.Str == "" {
		return translationSessionEvent{}, errors.New("translation: event is missing a nonempty string type")
	}
	return translationSessionEvent{data: data, kind: kind.Str}, nil
}

// Send writes a typed command once. session.close stops further input and waits
// for admitted writes, but does not drain output. Use Finish to drain and clean up.
// Success acknowledges a local write only; [websocket.DeliveryError] reports uncertain writes.
func (c *TranslationSession) Send(ctx context.Context, event RealtimeTranslationClientEventUnionParam) error {
	return c.send(ctx, event)
}

// SendRaw writes a JSON command unchanged, including future command types. The command must
// have a nonempty string type. session.close has the same semantics as typed Send.
func (c *TranslationSession) SendRaw(ctx context.Context, data json.RawMessage) error {
	return c.send(ctx, data)
}

func (c *TranslationSession) send(ctx context.Context, value any) error {
	if err := ctx.Err(); err != nil {
		return &websocket.DeliveryError{Cause: err}
	}
	c.mu.Lock()
	if c.stopping || c.terminal || c.closed {
		c.mu.Unlock()
		return &websocket.DeliveryError{Cause: ErrTranslationInputClosed}
	}
	if c.sending >= c.maxSending {
		c.mu.Unlock()
		return &websocket.DeliveryError{Cause: websocket.ErrPendingLimit}
	}
	if c.sending == 0 {
		c.idle = make(chan struct{})
	}
	c.sending++
	c.mu.Unlock()
	var data []byte
	var err error
	if raw, ok := value.(json.RawMessage); ok {
		data = slices.Clone(raw)
	} else {
		data, err = json.Marshal(value)
	}
	var event translationSessionEvent
	if err == nil {
		event, err = decodeTranslationSessionEvent(data)
	}
	if err != nil || event.kind == "session.close" {
		c.sent()
		if err != nil {
			return &websocket.DeliveryError{Cause: err}
		}
		return c.requestClose(ctx, data)
	}
	defer c.sent()
	return c.connection.SendRaw(ctx, data)
}

func (c *TranslationSession) sent() {
	c.mu.Lock()
	c.sending--
	if c.sending == 0 {
		close(c.idle)
	}
	c.mu.Unlock()
}

// requestClose shares the first close attempt with Finish and never retries it.
func (c *TranslationSession) requestClose(ctx context.Context, data []byte) error {
	c.mu.Lock()
	if c.closeDone != nil {
		done := c.closeDone
		c.mu.Unlock()
		select {
		case <-done:
			return c.closeErr // Published by closing done.
		case <-ctx.Done():
			// The first caller still owns the shared write. Its outcome is
			// unknown to a follower which stops waiting.
			return &websocket.DeliveryError{Cause: ctx.Err(), MayHaveBeenSent: true}
		}
	}
	if c.closed {
		c.mu.Unlock()
		return &websocket.DeliveryError{Cause: websocket.ErrClosed}
	}
	if c.terminal {
		c.mu.Unlock()
		return nil
	}
	c.stopping = true
	c.closeDone = make(chan struct{})
	idle, sending := c.idle, c.sending
	c.mu.Unlock()
	var err error
	if sending != 0 {
		select {
		case <-idle:
		case <-ctx.Done():
			err = &websocket.DeliveryError{Cause: ctx.Err()}
		case <-c.done:
			err = &websocket.DeliveryError{Cause: websocket.ErrClosed}
		}
	}
	if err == nil {
		c.mu.Lock()
		terminal := c.terminal
		c.mu.Unlock()
		if !terminal {
			err = c.connection.SendRaw(ctx, data)
		}
	}
	c.mu.Lock()
	c.closeErr = err
	close(c.closeDone)
	c.mu.Unlock()
	return err
}

func (c *TranslationSession) claim() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.receiving {
		return ErrTranslationReceiving
	}
	c.receiving = true
	return nil
}

func (c *TranslationSession) release() {
	c.mu.Lock()
	c.receiving = false
	c.mu.Unlock()
}

func (c *TranslationSession) receive(ctx context.Context) (translationSessionEvent, error) {
	c.mu.Lock()
	terminal := c.terminal
	c.mu.Unlock()
	if terminal {
		return translationSessionEvent{}, io.EOF
	}
	frame, err := c.connection.RecvRaw(ctx)
	var event translationSessionEvent
	if err == nil {
		event, err = decodeTranslationSessionEvent(frame)
	}
	if err == nil && event.kind == "session.closed" {
		c.mu.Lock()
		c.terminal = true
		c.mu.Unlock()
	}
	return event, err
}

// Recv receives the next typed event. Unknown events retain their RawJSON.
// Cancellation affects only this wait and leaves subsequent events available.
func (c *TranslationSession) Recv(ctx context.Context) (RealtimeTranslationServerEventUnion, error) {
	var event RealtimeTranslationServerEventUnion
	if err := c.claim(); err != nil {
		return event, err
	}
	defer c.release()
	value, err := c.receive(ctx)
	if err == nil {
		err = event.unmarshalWebsocketEvent(value.data)
	}
	return event, err
}

// RecvRaw receives the next event as its original JSON, without schema conversion.
func (c *TranslationSession) RecvRaw(ctx context.Context) (json.RawMessage, error) {
	if err := c.claim(); err != nil {
		return nil, err
	}
	defer c.release()
	event, err := c.receive(ctx)
	return event.data, err
}

// Finish owns event delivery, sends session.close once, and forwards every event
// through session.closed before closing the transport. onEvent runs synchronously
// and must return promptly. API errors are delivered as events. Cancellation,
// early EOF, decoding errors, and callback errors abort the transport and remain
// the result of subsequent Finish calls. Successful Finish is also idempotent.
// The context bounds draining; CloseTimeout separately bounds transport cleanup.
func (c *TranslationSession) Finish(ctx context.Context, onEvent func(RealtimeTranslationServerEventUnion) error) error {
	if err := c.claim(); err != nil {
		return err
	}
	defer c.release()
	c.mu.Lock()
	finished, previousErr, terminal := c.finishDone, c.finishErr, c.terminal
	c.mu.Unlock()
	if finished {
		return previousErr
	}
	if onEvent == nil {
		return errors.New("translation: an event consumer is required")
	}
	err := c.finish(ctx, terminal, onEvent)
	if err != nil {
		c.Abort()
	} else {
		err = c.Close()
	}
	c.mu.Lock()
	c.finishDone, c.finishErr = true, err
	c.mu.Unlock()
	return err
}

func (c *TranslationSession) finish(ctx context.Context, terminal bool, onEvent func(RealtimeTranslationServerEventUnion) error) error {
	if terminal {
		return nil
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	data, err := json.Marshal(RealtimeTranslationSessionCloseEventParam{})
	if err != nil {
		return err
	}
	closeErr := c.requestClose(ctx, data)
	// Local send failures cannot finish the protocol. A failed transport,
	// however, may have already queued the server's terminal and final output.
	if errors.Is(closeErr, context.Canceled) || errors.Is(closeErr, context.DeadlineExceeded) ||
		errors.Is(closeErr, websocket.ErrMessageLimit) || errors.Is(closeErr, websocket.ErrPendingLimit) {
		var delivery *websocket.DeliveryError
		if !errors.As(closeErr, &delivery) || !delivery.MayHaveBeenSent {
			return closeErr
		}
	}
	for {
		value, err := c.receive(ctx)
		if err != nil {
			if closeErr != nil {
				return closeErr
			}
			return err
		}
		var event RealtimeTranslationServerEventUnion
		if err := event.unmarshalWebsocketEvent(value.data); err != nil {
			return err
		}
		if err := onEvent(event); err != nil {
			return err
		}
		if value.kind == "session.closed" {
			return nil
		}
		if err := ctx.Err(); err != nil {
			return err
		}
	}
}

func (c *TranslationSession) stop() {
	c.mu.Lock()
	if !c.closed {
		c.closed = true
		close(c.done)
	}
	c.mu.Unlock()
}

// Close performs bounded transport cleanup. It does not send session.close or
// drain translated output; use Finish for protocol completion.
func (c *TranslationSession) Close() error {
	c.stop()
	return c.connection.Close()
}

// Abort immediately releases transport resources without protocol completion.
func (c *TranslationSession) Abort() {
	c.stop()
	c.connection.Abort()
}

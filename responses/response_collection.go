package responses

import (
	"context"
	"encoding/json"
	"errors"
	"sync/atomic"

	"github.com/tidwall/gjson"
)

// ErrResponseCollectionLimit means a response exceeded its caller-selected
// cumulative collection budget. Raw receiving and other lanes remain usable.
var ErrResponseCollectionLimit = errors.New("responses collection: response byte limit exceeded")

// ResponseCollection optionally collects final responses while receiving events.
// Use only its Recv or RecvRaw to consume its receiver; direct receives bypass
// collection. Do not copy it. Receive methods and Reset reject overlap; Err must
// not run concurrently with another method. It owns no reader goroutine or socket;
// close the original connection or lane when finished.
type ResponseCollection struct {
	receive     func(context.Context) (responseWebsocketEvent, error)
	localClose  func() error
	defaultLane bool
	maxBytes    int64
	bytes       int64
	accumulator ResponseAccumulator
	err         error
	receiving   atomic.Bool
}

// CollectResponse creates an opt-in final collector for the default stream.
// Unclaimed named-stream events remain receivable but do not change its state.
// maxBytes must be positive. It counts cumulative incoming UTF-8 message bytes
// for response.created, response.output_item.done and terminal responses in one
// turn, including metadata. Other progress and API error events are not charged.
// No existing connection limit or default receive behavior changes.
func (c *ResponseConnection) CollectResponse(maxBytes int64) (*ResponseCollection, error) {
	return newResponseCollection(c.connection.Recv, c.connection.LocalCloseError, true, maxBytes)
}

// CollectResponse creates an opt-in final collector for this lane. It requires
// the same explicit cumulative byte budget and exclusive receive use as the
// connection method. A detached lane or retired connection invalidates its cache.
func (l *ResponseLane) CollectResponse(maxBytes int64) (*ResponseCollection, error) {
	return newResponseCollection(l.lane.Recv, l.lane.LocalCloseError, false, maxBytes)
}

func newResponseCollection(receive func(context.Context) (responseWebsocketEvent, error), localClose func() error, defaultLane bool, maxBytes int64) (*ResponseCollection, error) {
	if maxBytes <= 0 {
		return nil, errors.New("responses collection: response byte limit must be positive")
	}
	return &ResponseCollection{receive: receive, localClose: localClose, defaultLane: defaultLane, maxBytes: maxBytes}, nil
}

// Reset clears collection state and errors, without consuming events or changing
// the socket. A consumed response.created also resets state, even when an earlier
// cached result has not been taken. Merely queued successor events do not reset it.
// Reset returns an error without mutation while a receive method is active.
func (c *ResponseCollection) Reset() error {
	if err := c.beginReceive(context.Background()); err != nil {
		return err
	}
	defer c.receiving.Store(false)
	c.resetState()
	return nil
}

func (c *ResponseCollection) resetState() {
	c.accumulator.Reset()
	c.bytes, c.err = 0, nil
}

// Err reports a collection error or explicit receiver closure. Collection errors
// remain until Reset or a consumed response.created. Receive cancellation does
// not become a collection error.
func (c *ResponseCollection) Err() error {
	if err := c.localClose(); err != nil {
		c.fail(err)
		return err
	}
	return c.err
}

func (c *ResponseCollection) observe(event responseWebsocketEvent) {
	if c.defaultLane && event.streamID != "" {
		return
	}
	kind := gjson.GetBytes(event.data, "type").Str
	if kind == "response.created" {
		c.resetState()
	}
	if c.err != nil {
		return
	}
	switch kind {
	case "error":
	case "response.created", "response.output_item.done", "response.completed", "response.failed", "response.incomplete":
		if int64(len(event.data)) > c.maxBytes-c.bytes {
			c.fail(ErrResponseCollectionLimit)
			return
		}
		c.bytes += int64(len(event.data))
	default:
		return
	}
	typed, err := event.decode()
	if err == nil {
		err = c.accumulator.AddEvent(typed)
	}
	if err != nil {
		c.fail(err)
	}
}

func (c *ResponseCollection) fail(err error) {
	c.accumulator.Reset()
	c.bytes, c.err = 0, err
}

func (c *ResponseCollection) recv(ctx context.Context) (responseWebsocketEvent, error) {
	if err := c.localClose(); err != nil {
		c.fail(err)
		return responseWebsocketEvent{}, err
	}
	event, err := c.receive(ctx)
	if err != nil {
		// A canceled wait can resume. An actual terminal already received remains
		// available even when the peer closes immediately afterward.
		if ctx.Err() == nil && c.err == nil && c.accumulator.terminalEvent == "" {
			c.fail(err)
		}
		return event, err
	}
	c.observe(event)
	return event, nil
}

func (c *ResponseCollection) beginReceive(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if !c.receiving.CompareAndSwap(false, true) {
		return errors.New("responses collection: receive already in progress")
	}
	return nil
}

// RecvRaw receives original JSON and automatically collects final-response data.
// A collection failure does not suppress the raw event or stop the receiver.
func (c *ResponseCollection) RecvRaw(ctx context.Context) (json.RawMessage, error) {
	if err := c.beginReceive(ctx); err != nil {
		return nil, err
	}
	defer c.receiving.Store(false)
	event, err := c.recv(ctx)
	if err != nil {
		return nil, err
	}
	return event.data, nil
}

// Recv receives a typed event and automatically collects final-response data.
// A collection failure remains available through Err and FinalResponse; typed
// decoding errors are returned normally. Original JSON remains on the event.
func (c *ResponseCollection) Recv(ctx context.Context) (ResponsesServerEventUnion, error) {
	if err := c.beginReceive(ctx); err != nil {
		return ResponsesServerEventUnion{}, err
	}
	defer c.receiving.Store(false)
	event, err := c.recv(ctx)
	if err != nil {
		return ResponsesServerEventUnion{}, err
	}
	return event.decode()
}

// FinalResponse takes an already collected terminal result once, or receives
// until the next result. Subsequent calls wait for the next turn. Missing/null
// output is reconstructed only from finalized items; supplied output is
// authoritative. Failed and incomplete statuses remain unchanged. Cancellation
// preserves collected progress for a later call. Explicit local close or
// recovery invalidates the cache; peer EOF alone does not erase a cached terminal.
func (c *ResponseCollection) FinalResponse(ctx context.Context) (*Response, error) {
	if err := c.beginReceive(ctx); err != nil {
		return nil, err
	}
	defer c.receiving.Store(false)
	for {
		if err := c.Err(); err != nil {
			return nil, err
		}
		if c.accumulator.terminalEvent != "" {
			response, err := c.accumulator.ReconstructedResponse()
			if err != nil {
				c.fail(err)
				return nil, err
			}
			c.resetState()
			return response, nil
		}
		if _, err := c.recv(ctx); err != nil {
			return nil, err
		}
	}
}

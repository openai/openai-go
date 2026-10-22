package websocket

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"errors"
	"io"
	"net"
	"slices"
	"sync"
	"time"

	wire "github.com/coder/websocket"
)

// RecoveryOptions controls optional recovery of a logical WebSocket session.
// Enablement does not restore server state or authorize replay of attempted writes.
type RecoveryOptions struct {
	Enabled bool
	// MaxAttempts counts additional attempts over the entire session. Zero selects five.
	MaxAttempts int
	// InitialDelay and MaxDelay default to 500 milliseconds and eight seconds.
	InitialDelay time.Duration
	MaxDelay     time.Duration
	// MaxPendingBytes bounds never-attempted writes waiting for an explicit flush.
	// Zero rejects sends during recovery. Send still waits for an actual write.
	MaxPendingBytes int64
	// OnReconnected may restore application state and explicitly flush pending writes.
	// It runs outside session locks, after the new transport is available.
	OnReconnected func()
}

type readError struct{ cause error }

func (e *readError) Error() string { return e.cause.Error() }
func (e *readError) Unwrap() error { return e.cause }

func retryable(err error) bool {
	var certificate x509.CertificateInvalidError
	var unknown x509.UnknownAuthorityError
	var hostname x509.HostnameError
	if errors.As(err, &certificate) || errors.As(err, &unknown) || errors.As(err, &hostname) {
		return false
	}
	var handshake *dialError
	if errors.As(err, &handshake) && handshake.statusCode != 0 {
		switch handshake.statusCode {
		case 408, 429, 500, 502, 503, 504:
			return true
		default:
			return false
		}
	}
	var read *readError
	if errors.As(err, &read) {
		switch wire.CloseStatus(read.cause) {
		case wire.StatusAbnormalClosure, wire.StatusInternalError, wire.StatusServiceRestart, wire.StatusTryAgainLater:
			return true
		case -1: // An EOF or network failure can occur without a close frame.
		default:
			return false
		}
	} else if handshake == nil {
		return false
	}
	var network net.Error
	return errors.As(err, &network) || errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF)
}

type pendingWrite struct {
	data   json.RawMessage
	ctx    context.Context
	result chan error
}

// Session retains one logical receiver while replacing failed physical connections.
// Recovery is driven by Recv; canceling a receive cancels only that wait. API events
// never trigger recovery. Keep receiving to observe transport failure and recovery.
type Session[T any] struct {
	mu           sync.Mutex
	active       *Connection[T]
	open         func(context.Context) (*Connection[T], error)
	options      RecoveryOptions
	ctx          context.Context
	cancel       context.CancelFunc
	recovering   bool
	ready        chan struct{}
	attempts     int
	err          error
	pending      []*pendingWrite
	pendingBytes int64
	pendingEpoch uint64
	flushing     *Connection[T]
	preparing    int
	writing      int
}

// OpenSession opens a physical connection once, then optionally recovers it on
// retryable transport failures. open must construct a fresh authenticated request.
func OpenSession[T any](ctx context.Context, open func(context.Context) (*Connection[T], error), options RecoveryOptions) (*Session[T], error) {
	if open == nil {
		return nil, errors.New("websocket: session opener is required")
	}
	if options.MaxAttempts < 0 || options.InitialDelay < 0 || options.MaxDelay < 0 || options.MaxPendingBytes < 0 {
		return nil, errors.New("websocket: recovery limits must not be negative")
	}
	if options.MaxAttempts == 0 {
		options.MaxAttempts = 5
	}
	if options.InitialDelay == 0 {
		options.InitialDelay = 500 * time.Millisecond
	}
	if options.MaxDelay == 0 {
		options.MaxDelay = 8 * time.Second
	}
	first, err := open(ctx)
	if err != nil {
		if first != nil {
			first.Abort()
		}
		return nil, err
	}
	if first == nil {
		return nil, errors.New("websocket: session opener returned no connection")
	}
	lifetime, cancel := context.WithCancel(context.WithoutCancel(ctx))
	ready := make(chan struct{})
	close(ready)
	return &Session[T]{active: first, open: open, options: options, ctx: lifetime, cancel: cancel, ready: ready}, nil
}

// Recv drains received events before recovering a failed transport. It preserves
// event order and does not interpret response or session termination events.
func (s *Session[T]) Recv(ctx context.Context) (T, error) {
	if !s.options.Enabled {
		return s.active.Recv(ctx)
	}
	var zero T
	for {
		s.mu.Lock()
		active, ready, recovering, terminal := s.active, s.ready, s.recovering, s.err
		s.mu.Unlock()
		if terminal != nil {
			return zero, terminal
		}
		if recovering {
			select {
			case <-ctx.Done():
				return zero, ctx.Err()
			case <-ready:
				continue
			}
		}
		event, err := active.Recv(ctx)
		if err == nil {
			return event, nil
		}
		if ctx.Err() != nil {
			return zero, ctx.Err()
		}
		s.mu.Lock()
		if s.active != active || s.recovering {
			s.mu.Unlock()
			continue
		}
		if s.err != nil {
			terminal = s.err
			s.mu.Unlock()
			return zero, terminal
		}
		if !s.options.Enabled || !retryable(err) {
			s.mu.Unlock()
			// The failed connection already owns its cleanup. Retire the session
			// without making this receive wait for the transport to close.
			s.stop(err)
			return zero, err
		}
		s.recovering = true
		s.ready = make(chan struct{})
		s.mu.Unlock()
		go s.recover(active, err)
	}
}

func (s *Session[T]) recover(old *Connection[T], original error) {
	old.Abort()
	err := original
	for {
		s.mu.Lock()
		if s.err != nil {
			s.mu.Unlock()
			return
		}
		if s.attempts >= s.options.MaxAttempts {
			s.mu.Unlock()
			s.fail(err)
			return
		}
		s.attempts++
		delay := s.options.InitialDelay
		for i := 1; i < s.attempts && delay < s.options.MaxDelay; i++ {
			if delay > s.options.MaxDelay/2 {
				delay = s.options.MaxDelay
			} else {
				delay *= 2
			}
		}
		if delay > s.options.MaxDelay {
			delay = s.options.MaxDelay
		}
		s.mu.Unlock()
		timer := time.NewTimer(delay)
		select {
		case <-s.ctx.Done():
			timer.Stop()
			return
		case <-timer.C:
		}
		next, failure := s.open(s.ctx)
		if failure == nil && next == nil {
			failure = errors.New("websocket: session opener returned no connection")
		}
		if failure != nil {
			if next != nil {
				next.Abort()
			}
			err = failure
			if !retryable(err) {
				s.fail(err)
				return
			}
			continue
		}
		s.mu.Lock()
		if s.err != nil {
			s.mu.Unlock()
			next.Abort()
			return
		}
		s.active = next
		s.writing = 0
		s.flushing = nil
		s.recovering = false
		close(s.ready)
		s.mu.Unlock()
		if s.options.OnReconnected != nil {
			s.options.OnReconnected()
		}
		return
	}
}

func (s *Session[T]) stop(err error) *Connection[T] {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.stopLocked(err)
}

func (s *Session[T]) stopLocked(err error) *Connection[T] {
	if s.err != nil {
		return s.active
	}
	s.err = err
	s.cancel()
	if s.recovering {
		s.recovering = false
		close(s.ready)
	}
	for _, p := range s.pending {
		select {
		case p.result <- &DeliveryError{Cause: err}:
		default:
		}
	}
	return s.active
}

func (s *Session[T]) fail(err error) { s.stop(err).fail(err) }

// SendJSON retains the existing local-write guarantee of Connection.SendJSON.
func (s *Session[T]) SendJSON(ctx context.Context, event any) error {
	if err := ctx.Err(); err != nil {
		return &DeliveryError{Cause: err}
	}
	s.mu.Lock()
	if s.err != nil {
		err := s.err
		s.mu.Unlock()
		return &DeliveryError{Cause: err}
	}
	if s.flushing != nil {
		s.mu.Unlock()
		return &DeliveryError{Cause: ErrPendingLimit}
	}
	if !s.recovering {
		active := s.active
		s.writing++
		s.mu.Unlock()
		return s.write(active, func() error { return active.SendJSON(ctx, event) })
	}
	// Reserve admission before calling user serialization. Recovery may finish while
	// it runs, but this command must still wait for explicit restoration and flush.
	if s.options.MaxPendingBytes == 0 || len(s.pending)+s.preparing >= cap(s.active.pending) {
		s.mu.Unlock()
		return &DeliveryError{Cause: ErrPendingLimit}
	}
	s.preparing++
	epoch := s.pendingEpoch
	s.mu.Unlock()
	reserved := true
	defer func() {
		if reserved {
			s.mu.Lock()
			s.preparing--
			s.mu.Unlock()
		}
	}()
	data, err := json.Marshal(event)
	s.mu.Lock()
	s.preparing--
	reserved = false
	if err != nil {
		s.mu.Unlock()
		return &DeliveryError{Cause: err}
	}
	p, err := s.enqueue(ctx, data, epoch)
	s.mu.Unlock()
	if err != nil {
		return err
	}
	return s.wait(ctx, p)
}

func (s *Session[T]) write(active *Connection[T], send func() error) (err error) {
	defer func() {
		s.failWrite(active, err)
		s.mu.Lock()
		if s.active == active {
			s.writing--
		}
		s.mu.Unlock()
	}()
	return send()
}

func (s *Session[T]) failWrite(active *Connection[T], err error) {
	var delivery *DeliveryError
	if !errors.As(err, &delivery) || !delivery.MayHaveBeenSent {
		return
	}
	s.mu.Lock()
	if s.active == active {
		s.stopLocked(err)
	}
	s.mu.Unlock()
}

// Send writes once, or waits for an explicitly flushed recovery queue entry.
// Cancellation removes an unattempted entry; an uncertain write is never retried.
func (s *Session[T]) Send(ctx context.Context, data json.RawMessage) error {
	if err := ctx.Err(); err != nil {
		return &DeliveryError{Cause: err}
	}
	s.mu.Lock()
	if s.err != nil {
		err := s.err
		s.mu.Unlock()
		return &DeliveryError{Cause: err}
	}
	if s.flushing != nil {
		s.mu.Unlock()
		return &DeliveryError{Cause: ErrPendingLimit}
	}
	if !s.recovering {
		active := s.active
		s.writing++
		s.mu.Unlock()
		return s.write(active, func() error { return active.Send(ctx, data) })
	}
	p, err := s.enqueue(ctx, data, s.pendingEpoch)
	s.mu.Unlock()
	if err != nil {
		return err
	}
	return s.wait(ctx, p)
}

// enqueue is called with mu held, including after serialization reserved admission.
func (s *Session[T]) enqueue(ctx context.Context, data json.RawMessage, epoch uint64) (*pendingWrite, error) {
	if s.err != nil {
		return nil, &DeliveryError{Cause: s.err}
	}
	if err := ctx.Err(); err != nil {
		return nil, &DeliveryError{Cause: err}
	}
	if epoch != s.pendingEpoch {
		return nil, &DeliveryError{Cause: ErrPendingTaken}
	}
	if int64(len(data)) > s.options.MaxPendingBytes-s.pendingBytes || s.options.MaxPendingBytes == 0 || len(s.pending)+s.preparing >= cap(s.active.pending) {
		return nil, &DeliveryError{Cause: ErrPendingLimit}
	}
	if limit := s.active.options.MaxMessageBytes; limit > 0 && int64(len(data)) > limit {
		return nil, &DeliveryError{Cause: ErrMessageLimit}
	}
	p := &pendingWrite{data: append(json.RawMessage(nil), data...), ctx: ctx, result: make(chan error, 1)}
	s.pending = append(s.pending, p)
	s.pendingBytes += int64(len(data))
	return p, nil
}

func (s *Session[T]) wait(ctx context.Context, p *pendingWrite) error {
	select {
	case err := <-p.result:
		return err
	case <-ctx.Done():
		s.mu.Lock()
		for i, item := range s.pending {
			if item == p {
				s.pending = slices.Delete(s.pending, i, i+1)
				s.pendingBytes -= int64(len(p.data))
				s.mu.Unlock()
				return &DeliveryError{Cause: ctx.Err()}
			}
		}
		s.mu.Unlock()
		return <-p.result
	}
}

// TakePending removes never-attempted commands. Their waiting Send calls return
// ErrPendingTaken with MayHaveBeenSent false. The returned bytes belong to the caller.
// An in-progress serialization cannot be returned and will not be queued.
func (s *Session[T]) TakePending() []json.RawMessage {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make([]json.RawMessage, 0, len(s.pending))
	for _, p := range s.pending {
		result = append(result, p.data)
		select {
		case p.result <- &DeliveryError{Cause: ErrPendingTaken}:
		default:
		}
	}
	s.pending = nil
	s.pendingBytes = 0
	s.pendingEpoch++
	return result
}

// ErrPendingTaken means TakePending removed or invalidated a never-attempted command.
var ErrPendingTaken = errors.New("websocket: pending command taken by application")

// FlushPending writes retained commands after the application restores server state.
func (s *Session[T]) FlushPending(ctx context.Context) error {
	s.mu.Lock()
	if s.err != nil {
		err := s.err
		s.mu.Unlock()
		return err
	}
	if s.recovering || s.flushing != nil || s.writing > 0 || s.preparing > 0 {
		s.mu.Unlock()
		return ErrPendingLimit
	}
	owner := s.active
	s.flushing = owner
	s.mu.Unlock()
	defer func() {
		s.mu.Lock()
		if s.flushing == owner {
			s.flushing = nil
		}
		s.mu.Unlock()
	}()
	for {
		s.mu.Lock()
		if s.active != owner {
			s.mu.Unlock()
			return ErrClosed
		}
		if s.err != nil {
			err := s.err
			s.mu.Unlock()
			return err
		}
		if s.recovering {
			s.mu.Unlock()
			return ErrPendingLimit
		}
		if len(s.pending) == 0 {
			s.mu.Unlock()
			return nil
		}
		p := s.pending[0]
		s.pending[0] = nil
		s.pending = s.pending[1:]
		s.pendingBytes -= int64(len(p.data))
		active := s.active
		epoch := s.pendingEpoch
		s.mu.Unlock()
		writeCtx, cancel := context.WithCancel(ctx)
		stop := context.AfterFunc(p.ctx, cancel)
		if p.ctx.Err() != nil {
			cancel()
		}
		err := active.Send(writeCtx, p.data)
		stop()
		cancel()
		var delivery *DeliveryError
		if errors.As(err, &delivery) && !delivery.MayHaveBeenSent && p.ctx.Err() == nil {
			s.mu.Lock()
			if p.ctx.Err() != nil {
				s.mu.Unlock()
				p.result <- &DeliveryError{Cause: p.ctx.Err()}
				return err
			}
			if epoch != s.pendingEpoch {
				cause := s.err
				if cause == nil {
					cause = ErrPendingTaken
				}
				s.mu.Unlock()
				p.result <- &DeliveryError{Cause: cause}
				return err
			}
			// Another recovery can admit commands while this old write finishes.
			// Returning an unsent command must not exceed the retained-work limit.
			if int64(len(p.data)) > s.options.MaxPendingBytes-s.pendingBytes || len(s.pending)+s.preparing >= cap(s.active.pending) {
				s.mu.Unlock()
				p.result <- &DeliveryError{Cause: ErrPendingLimit}
				return err
			}
			s.pending = append([]*pendingWrite{p}, s.pending...)
			s.pendingBytes += int64(len(p.data))
			terminal := s.err
			s.mu.Unlock()
			if terminal != nil {
				select {
				case p.result <- &DeliveryError{Cause: terminal}:
				default:
				}
			}
			return err
		}
		p.result <- err
		if err != nil {
			s.failWrite(active, err)
			return err
		}
	}
}

// Close cancels recovery and releases the current socket. Pending bytes remain
// available through TakePending, but all waiting Send calls are released.
func (s *Session[T]) Close() error { return s.stop(ErrClosed).Close() }
func (s *Session[T]) Abort()       { s.stop(ErrClosed).Abort() }

// Reconnect explicitly retires this logical connection and returns a new one.
func (s *Session[T]) Reconnect(ctx context.Context) (*Session[T], error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	s.mu.Lock()
	// Without recovery, preserve Connection.Reconnect's handling of active writes.
	if s.recovering || s.flushing != nil || s.preparing > 0 || (s.options.Enabled && s.writing > 0) {
		s.mu.Unlock()
		return nil, ErrPendingLimit
	}
	active := s.stopLocked(ErrClosed)
	s.mu.Unlock()
	next, err := active.Reconnect(ctx, s.open)
	if err != nil {
		active.fail(ErrClosed)
		return nil, err
	}
	replacement, err := OpenSession(ctx, func(context.Context) (*Connection[T], error) { return next, nil }, s.options)
	if replacement != nil {
		replacement.open = s.open
	}
	return replacement, err
}

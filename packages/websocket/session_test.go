package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	wire "github.com/coder/websocket"
)

type sessionPeer struct {
	socket        *wire.Conn
	authorization string
}

func sessionFixture(t *testing.T, options RecoveryOptions, beforeOpen func(context.Context, int32) error) (*Session[testEvent], <-chan sessionPeer, <-chan string, *atomic.Int32) {
	t.Helper()
	peers := make(chan sessionPeer, 8)
	writes := make(chan string, 16)
	attempts := new(atomic.Int32)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := wire.Accept(w, r, nil)
		if err != nil {
			return
		}
		defer func() { _ = c.CloseNow() }()
		peers <- sessionPeer{c, r.Header.Get("Authorization")}
		for {
			_, data, err := c.Read(r.Context())
			if err != nil {
				return
			}
			writes <- string(data)
		}
	}))
	t.Cleanup(server.Close)
	open := func(ctx context.Context) (*Connection[testEvent], error) {
		attempt := attempts.Add(1)
		if beforeOpen != nil {
			if err := beforeOpen(ctx, attempt); err != nil {
				return nil, err
			}
		}
		req, _ := http.NewRequestWithContext(ctx, "GET", server.URL, nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer synthetic-%d", attempt))
		return Dial(ctx, req, server.Client(), Options{}, func(data []byte) (testEvent, error) {
			var event testEvent
			err := json.Unmarshal(data, &event)
			return event, err
		}, func(testEvent) string { return "" })
	}
	s, err := OpenSession(context.Background(), open, options)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(s.Abort)
	return s, peers, writes, attempts
}

func receiveWithin[T any](t *testing.T, ch <-chan T) T {
	t.Helper()
	select {
	case value := <-ch:
		return value
	case <-time.After(5 * time.Second):
		t.Fatal("operation did not finish")
		var zero T
		return zero
	}
}

func waitPending(t *testing.T, s *Session[testEvent], count int) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		s.mu.Lock()
		n := len(s.pending)
		s.mu.Unlock()
		if n == count {
			return
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatal("send was not admitted to the pending queue")
}

func TestSessionRecoveryRequiresExplicitFlushAndRefreshesCredentials(t *testing.T) {
	opening, allow := make(chan struct{}), make(chan struct{})
	recovered := make(chan struct{}, 1)
	s, peers, writes, attempts := sessionFixture(t, RecoveryOptions{Enabled: true, MaxAttempts: 1, InitialDelay: time.Millisecond, MaxPendingBytes: 256, OnReconnected: func() { recovered <- struct{}{} }}, func(ctx context.Context, attempt int32) error {
		if attempt == 2 {
			close(opening)
			select {
			case <-allow:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		return nil
	})
	first := receiveWithin(t, peers)
	if first.authorization != "Bearer synthetic-1" {
		t.Fatal(first.authorization)
	}
	if err := s.Send(context.Background(), json.RawMessage(`{"type":"old"}`)); err != nil {
		t.Fatal(err)
	}
	if got := receiveWithin(t, writes); got != `{"type":"old"}` {
		t.Fatal(got)
	}
	read := make(chan error, 1)
	go func() { _, err := s.Recv(context.Background()); read <- err }()
	_ = first.socket.Close(wire.StatusServiceRestart, "restart")
	receiveWithin(t, opening)
	sent := make(chan error, 1)
	go func() { sent <- s.Send(context.Background(), json.RawMessage(`{"type":"held"}`)) }()
	waitPending(t, s, 1)
	close(allow)
	second := receiveWithin(t, peers)
	receiveWithin(t, recovered)
	if second.authorization != "Bearer synthetic-2" {
		t.Fatal(second.authorization)
	}
	select {
	case err := <-sent:
		t.Fatalf("send completed before flush: %v", err)
	default:
	}
	select {
	case got := <-writes:
		t.Fatalf("automatic replay: %s", got)
	default:
	}
	if err := s.Send(context.Background(), json.RawMessage(`{"type":"restore"}`)); err != nil {
		t.Fatal(err)
	}
	if err := s.FlushPending(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := receiveWithin(t, sent); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{`{"type":"restore"}`, `{"type":"held"}`} {
		if got := receiveWithin(t, writes); got != want {
			t.Fatalf("got %s want %s", got, want)
		}
	}
	_ = second.socket.Close(wire.StatusServiceRestart, "again")
	if err := receiveWithin(t, read); err == nil {
		t.Fatal("retry exhaustion returned no error")
	}
	if attempts.Load() != 2 {
		t.Fatal("retry budget reset after recovery")
	}
}

type blockingSessionEvent struct {
	started chan struct{}
	release chan struct{}
}

func (e blockingSessionEvent) MarshalJSON() ([]byte, error) {
	close(e.started)
	<-e.release
	return []byte(`{"type":"held"}`), nil
}

func TestSessionSerializationCannotBypassRestoration(t *testing.T) {
	opening, allow := make(chan struct{}), make(chan struct{})
	recovered := make(chan struct{}, 1)
	s, peers, writes, _ := sessionFixture(t, RecoveryOptions{Enabled: true, InitialDelay: time.Millisecond, MaxPendingBytes: 256, OnReconnected: func() { recovered <- struct{}{} }}, func(ctx context.Context, n int32) error {
		if n == 2 {
			close(opening)
			select {
			case <-allow:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		return nil
	})
	first := receiveWithin(t, peers)
	read := make(chan error, 1)
	go func() { _, err := s.Recv(context.Background()); read <- err }()
	_ = first.socket.Close(wire.StatusServiceRestart, "restart")
	receiveWithin(t, opening)
	event := blockingSessionEvent{make(chan struct{}), make(chan struct{})}
	sent := make(chan error, 1)
	go func() { sent <- s.SendJSON(context.Background(), event) }()
	receiveWithin(t, event.started)
	close(allow)
	receiveWithin(t, peers)
	receiveWithin(t, recovered)
	close(event.release)
	waitPending(t, s, 1)
	select {
	case got := <-writes:
		t.Fatalf("write escaped restoration: %s", got)
	default:
	}
	pending := s.TakePending()
	if len(pending) != 1 || string(pending[0]) != `{"type":"held"}` {
		t.Fatalf("pending = %s", pending)
	}
	if err := receiveWithin(t, sent); !errors.Is(err, ErrPendingTaken) {
		t.Fatal(err)
	}
	s.Abort()
	if err := receiveWithin(t, read); !errors.Is(err, ErrClosed) {
		t.Fatal(err)
	}
}

func TestSessionOldSerializerCannotBlockReplacementFlush(t *testing.T) {
	opening, allow := make(chan struct{}), make(chan struct{})
	recovered := make(chan struct{}, 1)
	s, peers, writes, _ := sessionFixture(t, RecoveryOptions{Enabled: true, InitialDelay: time.Millisecond, MaxPendingBytes: 256, OnReconnected: func() { recovered <- struct{}{} }}, func(ctx context.Context, n int32) error {
		if n == 2 {
			close(opening)
			select {
			case <-allow:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		return nil
	})
	first := receiveWithin(t, peers)
	event := blockingSessionEvent{make(chan struct{}), make(chan struct{})}
	oldSend := make(chan error, 1)
	go func() { oldSend <- s.SendJSON(context.Background(), event) }()
	receiveWithin(t, event.started)
	read := make(chan error, 1)
	go func() { _, err := s.Recv(context.Background()); read <- err }()
	_ = first.socket.Close(wire.StatusServiceRestart, "restart")
	receiveWithin(t, opening)
	queued := make(chan error, 1)
	go func() { queued <- s.Send(context.Background(), json.RawMessage(`{"type":"new"}`)) }()
	waitPending(t, s, 1)
	close(allow)
	receiveWithin(t, peers)
	receiveWithin(t, recovered)
	if err := s.FlushPending(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := receiveWithin(t, queued); err != nil {
		t.Fatal(err)
	}
	if got := receiveWithin(t, writes); got != `{"type":"new"}` {
		t.Fatal(got)
	}
	close(event.release)
	var delivery *DeliveryError
	if err := receiveWithin(t, oldSend); !errors.As(err, &delivery) || delivery.MayHaveBeenSent {
		t.Fatalf("old serializer result = %v", err)
	}
	if err := s.Send(context.Background(), json.RawMessage(`{"type":"still-open"}`)); err != nil {
		t.Fatal(err)
	}
	if got := receiveWithin(t, writes); got != `{"type":"still-open"}` {
		t.Fatal(got)
	}
	s.Abort()
	if err := receiveWithin(t, read); !errors.Is(err, ErrClosed) {
		t.Fatal(err)
	}
}

func TestSessionCloseCancelsProviderAndRejectsDefaultQueue(t *testing.T) {
	opening, canceled := make(chan struct{}), make(chan struct{})
	s, peers, _, attempts := sessionFixture(t, RecoveryOptions{Enabled: true, InitialDelay: time.Millisecond}, func(ctx context.Context, n int32) error {
		if n == 2 {
			close(opening)
			<-ctx.Done()
			close(canceled)
			return ctx.Err()
		}
		return nil
	})
	first := receiveWithin(t, peers)
	read := make(chan error, 1)
	go func() { _, err := s.Recv(context.Background()); read <- err }()
	_ = first.socket.Close(wire.StatusServiceRestart, "restart")
	receiveWithin(t, opening)
	err := s.Send(context.Background(), json.RawMessage(`{}`))
	var delivery *DeliveryError
	if !errors.As(err, &delivery) || delivery.MayHaveBeenSent || !errors.Is(err, ErrPendingLimit) {
		t.Fatal(err)
	}
	s.Abort()
	receiveWithin(t, canceled)
	if err := receiveWithin(t, read); !errors.Is(err, ErrClosed) {
		t.Fatal(err)
	}
	if attempts.Load() != 2 {
		t.Fatal(attempts.Load())
	}
}

func TestSessionDoesNotRecoverNormalCloseOrDecodeErrors(t *testing.T) {
	for _, scenario := range []string{"normal", "policy", "decode", "disabled"} {
		t.Run(scenario, func(t *testing.T) {
			s, peers, _, attempts := sessionFixture(t, RecoveryOptions{Enabled: scenario != "disabled", InitialDelay: time.Millisecond}, nil)
			peer := receiveWithin(t, peers)
			switch scenario {
			case "decode":
				_ = peer.socket.Write(context.Background(), wire.MessageText, []byte("invalid"))
			case "normal":
				go func() { _ = peer.socket.Close(wire.StatusNormalClosure, "done") }()
			case "policy":
				go func() { _ = peer.socket.Close(wire.StatusPolicyViolation, "policy") }()
			case "disabled":
				_ = peer.socket.CloseNow()
			}
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			if _, err := s.Recv(ctx); err == nil || errors.Is(err, context.DeadlineExceeded) {
				t.Fatal(err)
			}
			if attempts.Load() != 1 {
				t.Fatal("unexpected recovery")
			}
		})
	}
}

func TestSessionBoundsSerializationAndCancelsRecoveryWithoutWaitingForIt(t *testing.T) {
	opening := make(chan struct{})
	canceled := make(chan struct{})
	s, peers, _, _ := sessionFixture(t, RecoveryOptions{Enabled: true, InitialDelay: time.Millisecond, MaxPendingBytes: 1024}, func(ctx context.Context, n int32) error {
		if n == 2 {
			close(opening)
			<-ctx.Done()
			close(canceled)
			return ctx.Err()
		}
		return nil
	})
	peer := receiveWithin(t, peers)
	reading := make(chan error, 1)
	go func() { _, err := s.Recv(context.Background()); reading <- err }()
	_ = peer.socket.CloseNow()
	receiveWithin(t, opening)
	release := make(chan struct{})
	results := make(chan error, 16)
	for i := 0; i < 16; i++ {
		event := blockingSessionEvent{make(chan struct{}), release}
		go func() { results <- s.SendJSON(context.Background(), event) }()
		receiveWithin(t, event.started)
	}
	extra := blockingSessionEvent{make(chan struct{}), release}
	if err := s.SendJSON(context.Background(), extra); !errors.Is(err, ErrPendingLimit) {
		t.Fatal(err)
	}
	s.Abort()
	receiveWithin(t, canceled)
	if err := receiveWithin(t, reading); !errors.Is(err, ErrClosed) {
		t.Fatal(err)
	}
	close(release)
	for i := 0; i < 16; i++ {
		if err := receiveWithin(t, results); !errors.Is(err, ErrClosed) {
			t.Fatal(err)
		}
	}
}

func TestSessionOverflowAndCancellationPreserveOtherPendingCommands(t *testing.T) {
	opening, allow := make(chan struct{}), make(chan struct{})
	s, peers, _, _ := sessionFixture(t, RecoveryOptions{Enabled: true, InitialDelay: time.Millisecond, MaxPendingBytes: 4}, func(ctx context.Context, n int32) error {
		if n == 2 {
			close(opening)
			select {
			case <-allow:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		return nil
	})
	peer := receiveWithin(t, peers)
	reading := make(chan error, 1)
	go func() { _, err := s.Recv(context.Background()); reading <- err }()
	_ = peer.socket.CloseNow()
	receiveWithin(t, opening)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	canceled := make(chan error, 1)
	retained := make(chan error, 1)
	go func() { canceled <- s.Send(ctx, json.RawMessage(`{}`)) }()
	waitPending(t, s, 1)
	go func() { retained <- s.Send(context.Background(), json.RawMessage(`[]`)) }()
	waitPending(t, s, 2)
	if err := s.Send(context.Background(), json.RawMessage(`{}`)); !errors.Is(err, ErrPendingLimit) {
		t.Fatal(err)
	}
	cancel()
	var delivery *DeliveryError
	if err := receiveWithin(t, canceled); !errors.As(err, &delivery) || delivery.MayHaveBeenSent || !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
	pending := s.TakePending()
	if len(pending) != 1 || string(pending[0]) != `[]` {
		t.Fatalf("pending=%s", pending)
	}
	if err := receiveWithin(t, retained); !errors.Is(err, ErrPendingTaken) {
		t.Fatal(err)
	}
	s.Abort()
	if err := receiveWithin(t, reading); !errors.Is(err, ErrClosed) {
		t.Fatal(err)
	}
}

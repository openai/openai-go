package websocket

import (
	"context"
	"crypto/tls"
	"crypto/x509"
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
	return sessionFixtureContext(t, context.Background(), options, beforeOpen)
}

func blockSecondOpen(opening chan struct{}, allow <-chan struct{}) func(context.Context, int32) error {
	return func(ctx context.Context, attempt int32) error {
		if attempt == 2 {
			close(opening)
			select {
			case <-allow:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		return nil
	}
}

func sessionFixtureContext(t *testing.T, ctx context.Context, options RecoveryOptions, beforeOpen func(context.Context, int32) error) (*Session[testEvent], <-chan sessionPeer, <-chan string, *atomic.Int32) {
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
	s, err := OpenSession(ctx, open, options)
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
	type contextKey struct{}
	openingCtx, cancelOpening := context.WithCancel(context.WithValue(context.Background(), contextKey{}, "original"))
	defer cancelOpening()
	opening, allow := make(chan struct{}), make(chan struct{})
	recovered := make(chan struct{}, 1)
	beforeOpen := blockSecondOpen(opening, allow)
	s, peers, writes, attempts := sessionFixtureContext(t, openingCtx, RecoveryOptions{Enabled: true, MaxAttempts: 1, InitialDelay: time.Millisecond, MaxPendingBytes: 256, OnReconnected: func() { recovered <- struct{}{} }}, func(ctx context.Context, attempt int32) error {
		if attempt == 2 {
			if got := ctx.Value(contextKey{}); got != "original" {
				t.Errorf("recovery context value = %v", got)
			}
		}
		return beforeOpen(ctx, attempt)
	})
	first := receiveWithin(t, peers)
	cancelOpening()
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

type panickingSessionEvent struct{}

func (panickingSessionEvent) MarshalJSON() ([]byte, error) { panic("serialization") }

func sendPanickingEvent(t *testing.T, s *Session[testEvent]) {
	t.Helper()
	defer func() {
		if got := recover(); got != "serialization" {
			t.Fatalf("serialization panic = %v", got)
		}
	}()
	_ = s.SendJSON(context.Background(), panickingSessionEvent{})
}

func TestSessionPanicDoesNotBlockFlush(t *testing.T) {
	s, _, _, _ := sessionFixture(t, RecoveryOptions{Enabled: true}, nil)
	sendPanickingEvent(t, s)
	if err := s.FlushPending(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestSessionSerializationCannotBypassRestoration(t *testing.T) {
	opening, allow := make(chan struct{}), make(chan struct{})
	recovered := make(chan struct{}, 1)
	s, peers, writes, _ := sessionFixture(t, RecoveryOptions{Enabled: true, InitialDelay: time.Millisecond, MaxPendingBytes: 256, OnReconnected: func() { recovered <- struct{}{} }}, blockSecondOpen(opening, allow))
	first := receiveWithin(t, peers)
	read := make(chan error, 1)
	go func() { _, err := s.Recv(context.Background()); read <- err }()
	_ = first.socket.Close(wire.StatusServiceRestart, "restart")
	receiveWithin(t, opening)
	sendPanickingEvent(t, s)
	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	never := blockingSessionEvent{make(chan struct{}), make(chan struct{})}
	close(never.release)
	if err := s.SendJSON(canceled, never); !errors.Is(err, context.Canceled) {
		t.Errorf("canceled send = %v", err)
	}
	select {
	case <-never.started:
		t.Error("canceled send invoked serialization")
	default:
	}
	event := blockingSessionEvent{make(chan struct{}), make(chan struct{})}
	sent := make(chan error, 1)
	go func() { sent <- s.SendJSON(context.Background(), event) }()
	receiveWithin(t, event.started)
	close(allow)
	receiveWithin(t, peers)
	receiveWithin(t, recovered)
	if err := s.FlushPending(context.Background()); !errors.Is(err, ErrPendingLimit) {
		t.Errorf("flush while recovery-admitted serialization is running = %v", err)
	}
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
	if err := s.FlushPending(context.Background()); err != nil {
		t.Fatal(err)
	}
	s.Abort()
	if err := receiveWithin(t, read); !errors.Is(err, ErrClosed) {
		t.Fatal(err)
	}
}

func TestSessionTakePendingRejectsInProgressSerialization(t *testing.T) {
	opening, allow, recovered := make(chan struct{}), make(chan struct{}), make(chan struct{})
	s, peers, writes, _ := sessionFixture(t, RecoveryOptions{Enabled: true, InitialDelay: time.Millisecond, MaxPendingBytes: 256, OnReconnected: func() { close(recovered) }}, blockSecondOpen(opening, allow))
	first := receiveWithin(t, peers)
	go func() { _, _ = s.Recv(context.Background()) }()
	_ = first.socket.CloseNow()
	receiveWithin(t, opening)
	event := blockingSessionEvent{make(chan struct{}), make(chan struct{}, 1)}
	defer close(event.release)
	sent := make(chan error, 1)
	go func() { sent <- s.SendJSON(context.Background(), event) }()
	receiveWithin(t, event.started)
	if pending := s.TakePending(); len(pending) != 0 {
		t.Fatalf("pending before serialization finished = %s", pending)
	}
	close(allow)
	receiveWithin(t, recovered)
	if err := s.FlushPending(context.Background()); !errors.Is(err, ErrPendingLimit) {
		t.Fatalf("flush during serialization = %v", err)
	}
	event.release <- struct{}{}
	var delivery *DeliveryError
	if err := receiveWithin(t, sent); !errors.Is(err, ErrPendingTaken) || !errors.As(err, &delivery) || delivery.MayHaveBeenSent {
		t.Fatalf("send = %v", err)
	}
	if pending := s.TakePending(); len(pending) != 0 {
		t.Fatalf("pending after serialization finished = %s", pending)
	}
	if err := s.FlushPending(context.Background()); err != nil {
		t.Fatalf("flush after serialization = %v", err)
	}
	select {
	case got := <-writes:
		t.Fatalf("removed command was sent: %s", got)
	default:
	}
}

func TestSessionOldSerializerCannotBlockReplacementFlush(t *testing.T) {
	opening, allow := make(chan struct{}), make(chan struct{})
	recovered := make(chan struct{}, 1)
	s, peers, writes, _ := sessionFixture(t, RecoveryOptions{Enabled: true, InitialDelay: time.Millisecond, MaxPendingBytes: 256, OnReconnected: func() { recovered <- struct{}{} }}, blockSecondOpen(opening, allow))
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

func TestSessionRecoveryClassifiesHandshakeFailures(t *testing.T) {
	for _, tc := range []struct {
		status       int
		wantAttempts int32
	}{{http.StatusUnauthorized, 2}, {http.StatusServiceUnavailable, 3}} {
		t.Run(http.StatusText(tc.status), func(t *testing.T) {
			var attempts atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				attempt := attempts.Add(1)
				if attempt == 2 {
					w.WriteHeader(tc.status)
					return
				}
				peer, err := wire.Accept(w, r, nil)
				if err != nil {
					t.Errorf("accept: %v", err)
					return
				}
				defer func() { _ = peer.CloseNow() }()
				if attempt == 3 {
					if err := peer.Write(r.Context(), wire.MessageText, []byte("recovered")); err != nil {
						t.Errorf("write recovered event: %v", err)
					}
				}
			}))
			defer server.Close()
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			s, err := OpenSession(ctx, func(ctx context.Context) (*Connection[string], error) {
				request, err := http.NewRequestWithContext(ctx, http.MethodGet, server.URL, nil)
				if err != nil {
					return nil, err
				}
				return Dial(ctx, request, server.Client(), Options{}, func(data []byte) (string, error) { return string(data), nil }, func(string) string { return "" })
			}, RecoveryOptions{Enabled: true, MaxAttempts: 2, InitialDelay: time.Millisecond})
			if err != nil {
				t.Fatal(err)
			}
			defer s.Abort()
			event, err := s.Recv(ctx)
			if tc.wantAttempts == 3 {
				if err != nil || event != "recovered" {
					t.Fatalf("receive = %q, %v", event, err)
				}
			} else if err == nil || errors.Is(err, context.DeadlineExceeded) {
				t.Fatalf("receive = %q, %v", event, err)
			}
			if got := attempts.Load(); got != tc.wantAttempts {
				t.Fatalf("handshake attempts = %d, want %d", got, tc.wantAttempts)
			}
		})
	}
	t.Run("untrusted certificate", func(t *testing.T) {
		server := httptest.NewTLSServer(http.NotFoundHandler())
		defer server.Close()
		client := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{RootCAs: x509.NewCertPool()}}}
		s, peers, _, attempts := sessionFixture(t, RecoveryOptions{Enabled: true, MaxAttempts: 2, InitialDelay: time.Millisecond}, func(ctx context.Context, attempt int32) error {
			if attempt == 1 {
				return nil
			}
			request, err := http.NewRequestWithContext(ctx, http.MethodGet, server.URL, nil)
			if err != nil {
				return err
			}
			_, err = Dial(ctx, request, client, Options{}, func([]byte) (testEvent, error) { return testEvent{}, nil }, func(testEvent) string { return "" })
			return err
		})
		_ = receiveWithin(t, peers).socket.CloseNow()
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_, err := s.Recv(ctx)
		var certificate x509.UnknownAuthorityError
		if !errors.As(err, &certificate) {
			t.Fatalf("receive = %v, want untrusted certificate", err)
		}
		if got := attempts.Load(); got != 2 {
			t.Fatalf("connection attempts = %d, want 2", got)
		}
	})
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
	s, peers, _, _ := sessionFixture(t, RecoveryOptions{Enabled: true, InitialDelay: time.Millisecond, MaxPendingBytes: 4}, blockSecondOpen(opening, allow))
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

func TestSessionCanceledFlushBeforeWriteIsDefinitelyUnsent(t *testing.T) {
	for _, scenario := range []string{"send canceled", "taken then flush canceled"} {
		t.Run(scenario, func(t *testing.T) {
			opening, allow, recovered := make(chan struct{}), make(chan struct{}), make(chan struct{})
			s, peers, _, _ := sessionFixture(t, RecoveryOptions{Enabled: true, InitialDelay: time.Millisecond, MaxPendingBytes: 256, OnReconnected: func() { close(recovered) }}, blockSecondOpen(opening, allow))
			first := receiveWithin(t, peers)
			go func() { _, _ = s.Recv(context.Background()) }()
			_ = first.socket.CloseNow()
			receiveWithin(t, opening)
			sendCtx, cancelSend := context.WithCancel(context.Background())
			defer cancelSend()
			sent := make(chan error, 1)
			go func() { sent <- s.Send(sendCtx, json.RawMessage(`{}`)) }()
			waitPending(t, s, 1)
			close(allow)
			receiveWithin(t, recovered)
			s.active.writer <- struct{}{}
			defer func() { <-s.active.writer }()
			flushCtx, cancelFlush := context.WithCancel(context.Background())
			defer cancelFlush()
			flushed := make(chan error, 1)
			go func() { flushed <- s.FlushPending(flushCtx) }()
			waitPending(t, s, 0)
			want := context.Canceled
			if scenario == "send canceled" {
				cancelSend()
			} else {
				if pending := s.TakePending(); len(pending) != 0 {
					t.Fatalf("in-flight command was taken: %s", pending)
				}
				cancelFlush()
				want = ErrPendingTaken
			}
			if err := receiveWithin(t, flushed); !errors.Is(err, context.Canceled) {
				t.Fatalf("flush = %v", err)
			}
			if pending := s.TakePending(); len(pending) != 0 {
				t.Fatalf("command reappeared after take: %s", pending)
			}
			var delivery *DeliveryError
			if err := receiveWithin(t, sent); !errors.Is(err, want) || !errors.As(err, &delivery) || delivery.MayHaveBeenSent {
				t.Fatalf("send = %v, want definitely unsent %v", err, want)
			}
		})
	}
}

type cancelOnSecondCheck struct {
	context.Context
	cancel context.CancelFunc
	checks atomic.Int32
}

func (c *cancelOnSecondCheck) Err() error {
	if c.checks.Add(1) == 2 {
		c.cancel()
	}
	return c.Context.Err()
}

func TestSessionReconnectDuringSendWithoutRecovery(t *testing.T) {
	s, _, writes, _ := sessionFixture(t, RecoveryOptions{}, nil)
	event := blockingSessionEvent{make(chan struct{}), make(chan struct{}, 1)}
	defer close(event.release)
	sent := make(chan error, 1)
	go func() { sent <- s.SendJSON(context.Background(), event) }()
	receiveWithin(t, event.started)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	next, err := s.Reconnect(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(next.Abort)
	event.release <- struct{}{}
	var delivery *DeliveryError
	if err := receiveWithin(t, sent); !errors.Is(err, ErrClosed) || !errors.As(err, &delivery) || delivery.MayHaveBeenSent {
		t.Fatalf("old send = %v, want definitely unsent ErrClosed", err)
	}
	if err := next.Send(ctx, json.RawMessage(`{"type":"new"}`)); err != nil {
		t.Fatal(err)
	}
	if got := receiveWithin(t, writes); got != `{"type":"new"}` {
		t.Fatal(got)
	}
}

func TestSessionCanceledReconnectRetiresSocketAndCanRetry(t *testing.T) {
	for _, enabled := range []bool{false, true} {
		t.Run(fmt.Sprintf("recovery=%t", enabled), func(t *testing.T) {
			s, peers, writes, _ := sessionFixture(t, RecoveryOptions{Enabled: enabled}, nil)
			receiveWithin(t, peers)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			if _, err := s.Reconnect(&cancelOnSecondCheck{Context: ctx, cancel: cancel}); !errors.Is(err, context.Canceled) {
				t.Fatalf("reconnect = %v", err)
			}
			receiveWithin(t, s.active.readerDone)
			next, err := s.Reconnect(context.Background())
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(next.Abort)
			if err := next.Send(context.Background(), json.RawMessage(`{"type":"new"}`)); err != nil {
				t.Fatal(err)
			}
			if got := receiveWithin(t, writes); got != `{"type":"new"}` {
				t.Fatal(got)
			}
		})
	}
}

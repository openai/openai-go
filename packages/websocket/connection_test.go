package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	wire "github.com/coder/websocket"
)

type testEvent struct {
	Type     string `json:"type"`
	StreamID string `json:"stream_id"`
	Text     string `json:"text"`
}

func testConnection(t *testing.T, options Options, handler func(context.Context, *wire.Conn)) *Connection[testEvent] {
	t.Helper()
	serverDone := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer close(serverDone)
		socket, err := wire.Accept(w, r, nil)
		if err != nil {
			t.Errorf("accept upgrade: %v", err)
			return
		}
		defer func() { _ = socket.CloseNow() }()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		handler(ctx, socket)
	}))
	t.Cleanup(server.Close)
	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	connection, err := Dial(ctx, req, server.Client(), options, func(data []byte) (testEvent, error) {
		var event testEvent
		decodeErr := json.Unmarshal(data, &event)
		return event, decodeErr
	}, func(event testEvent) string { return event.StreamID })
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		connection.Abort()
		select {
		case <-serverDone:
		case <-time.After(10 * time.Second):
			t.Error("server did not stop")
		}
	})
	return connection
}

func writeEvent(ctx context.Context, socket *wire.Conn, event testEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return socket.Write(ctx, wire.MessageText, data)
}

func TestFailedReconnectClosesPartialReplacement(t *testing.T) {
	wait := func(ctx context.Context, socket *wire.Conn) { _, _, _ = socket.Read(ctx) }
	connection := testConnection(t, Options{}, wait)
	replacement := testConnection(t, Options{}, wait)
	want := errors.New("restoration failed")
	next, err := connection.Reconnect(context.Background(), func(context.Context) (*Connection[testEvent], error) {
		return replacement, want
	})
	if next != nil || !errors.Is(err, want) {
		t.Fatalf("reconnect = %v, %v", next, err)
	}
	select {
	case <-replacement.readerDone:
	default:
		t.Fatal("failed reconnect leaked a partially opened replacement")
	}
}

func TestCanceledReconnectCanBeRetriedExplicitly(t *testing.T) {
	wait := func(ctx context.Context, socket *wire.Conn) { _, _, _ = socket.Read(ctx) }
	connection := testConnection(t, Options{}, wait)
	replacement := testConnection(t, Options{}, wait)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err := connection.Reconnect(ctx, func(context.Context) (*Connection[testEvent], error) {
		cancel()
		return replacement, nil
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("canceled reconnect = %v", err)
	}
	select {
	case <-replacement.readerDone:
	default:
		t.Fatal("canceled reconnect leaked its replacement")
	}
	retry := testConnection(t, Options{}, wait)
	got, err := connection.Reconnect(context.Background(), func(context.Context) (*Connection[testEvent], error) {
		return retry, nil
	})
	if err != nil || got != retry {
		t.Fatalf("explicit retry = %v, %v", got, err)
	}
}

func TestReceiveCancellationDoesNotConsumeNextEvent(t *testing.T) {
	release := make(chan struct{})
	connection := testConnection(t, Options{}, func(ctx context.Context, socket *wire.Conn) {
		<-release
		if err := writeEvent(ctx, socket, testEvent{Type: "response.completed", Text: "still connected"}); err != nil {
			t.Errorf("write: %v", err)
		}
		_, _, _ = socket.Read(ctx)
	})
	cancelled, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := connection.Recv(cancelled); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancelled receive = %v, want context.Canceled", err)
	}
	close(release)
	ctx, done := context.WithTimeout(context.Background(), time.Second)
	defer done()
	event, err := connection.Recv(ctx)
	if err != nil || event.Text != "still connected" {
		t.Fatalf("next receive = %+v, %v", event, err)
	}
}

func TestLanesDetachIndependently(t *testing.T) {
	release := make(chan struct{})
	connection := testConnection(t, Options{}, func(ctx context.Context, socket *wire.Conn) {
		<-release
		for _, event := range []testEvent{{StreamID: "a", Text: "first"}, {StreamID: "b", Text: "second"}, {Text: "default"}} {
			if err := writeEvent(ctx, socket, event); err != nil {
				t.Errorf("write: %v", err)
				return
			}
		}
		_, _, _ = socket.Read(ctx)
	})
	laneA, err := connection.Lane("a")
	if err != nil {
		t.Fatal(err)
	}
	laneB, err := connection.Lane("b")
	if err != nil {
		t.Fatal(err)
	}
	close(release)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if event, err := laneA.Recv(ctx); err != nil || event.Text != "first" {
		t.Fatalf("lane a = %+v,%v", event, err)
	}
	laneA.Close()
	laneA.Close()
	if _, err := laneA.Recv(ctx); !errors.Is(err, ErrLaneClosed) {
		t.Fatalf("detached lane = %v", err)
	}
	if event, err := laneB.Recv(ctx); err != nil || event.Text != "second" {
		t.Fatalf("lane b = %+v,%v", event, err)
	}
	if event, err := connection.Recv(ctx); err != nil || event.Text != "default" {
		t.Fatalf("default lane = %+v,%v", event, err)
	}
}

func TestLaneCloseReservesKeyAndReleasesBufferedEvents(t *testing.T) {
	const retiredKey = "retired/key" // Generic keys are not Responses stream IDs.
	payload := strings.Repeat("x", 256)
	connection := testConnection(t, Options{MaxLanes: 2, MaxBufferedEvents: 3, MaxBufferedBytes: 800}, func(ctx context.Context, socket *wire.Conn) {
		for _, events := range [][]testEvent{
			{{StreamID: retiredKey, Text: payload}, {StreamID: retiredKey, Text: payload}, {Text: "buffered"}},
			{{StreamID: retiredKey, Text: payload}, {StreamID: "active", Text: payload}, {StreamID: "unclaimed", Text: "fallback"}},
		} {
			if _, _, err := socket.Read(ctx); err != nil {
				t.Errorf("read release command: %v", err)
				return
			}
			for _, event := range events {
				if err := writeEvent(ctx, socket, event); err != nil {
					t.Errorf("write event for %q: %v", event.StreamID, err)
					return
				}
			}
		}
		_, _, _ = socket.Read(ctx)
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	lane, laneErr := connection.Lane(retiredKey)
	if laneErr != nil {
		t.Fatal(laneErr)
	}
	canceled, stop := context.WithCancel(ctx)
	stop()
	if _, err := lane.Recv(canceled); !errors.Is(err, context.Canceled) {
		t.Errorf("Recv(canceled) = %v, want context.Canceled", err)
	}
	if err := connection.Send(ctx, []byte("buffer")); err != nil {
		t.Fatal(err)
	}
	// Reading the last frame proves both earlier lane events are buffered.
	if event, err := connection.Recv(ctx); err != nil || event.Text != "buffered" {
		t.Fatalf("Recv(buffer barrier) = %+v, %v, want buffered", event, err)
	}
	lane.Close()
	lane.Close()
	if _, err := lane.Recv(ctx); !errors.Is(err, ErrLaneClosed) {
		t.Errorf("Recv(closed lane) = %v, want ErrLaneClosed", err)
	}
	if replacement, err := connection.Lane(retiredKey); err == nil {
		replacement.Close()
		t.Error("Lane(retired) succeeded after Close, want reserved-key error")
	}
	active, err := connection.Lane("active")
	if err != nil {
		t.Fatal(err)
	}
	if err := connection.Send(ctx, []byte("after close")); err != nil {
		t.Fatal(err)
	}
	if event, err := connection.Recv(ctx); err != nil || event.StreamID != retiredKey || event.Text != payload {
		t.Fatalf("Recv(retired fallback) = %+v, %v, want late retired event", event, err)
	}
	if event, err := connection.Recv(ctx); err != nil || event.StreamID != "unclaimed" || event.Text != "fallback" {
		t.Fatalf("Recv(unclaimed fallback) = %+v, %v, want unclaimed event", event, err)
	}
	if event, err := active.Recv(ctx); err != nil || event.Text != payload {
		t.Fatalf("Recv(active lane) = %+v, %v, want active event", event, err)
	}
	active.Close()
	if extra, err := connection.Lane("extra"); err == nil {
		extra.Close()
		t.Error("Lane(extra) succeeded after both slots closed, want lane limit error")
	}
}

func TestReceiveLimits(t *testing.T) {
	for _, test := range []struct {
		name    string
		options Options
		payload string
		want    error
	}{
		{name: "message", options: Options{MaxMessageBytes: 128}, payload: strings.Repeat("x", 1024), want: wire.ErrMessageTooBig},
		{name: "buffer", options: Options{MaxBufferedBytes: 128}, payload: strings.Repeat("x", 1024), want: ErrBufferLimit},
	} {
		t.Run(test.name, func(t *testing.T) {
			connection := testConnection(t, test.options, func(ctx context.Context, socket *wire.Conn) {
				_ = writeEvent(ctx, socket, testEvent{Text: test.payload})
				_, _, _ = socket.Read(ctx)
			})
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			_, err := connection.Recv(ctx)
			if !errors.Is(err, test.want) {
				t.Fatalf("receive limit = %v, want %v", err, test.want)
			}
		})
	}
}

func TestLargeResponseMessage(t *testing.T) {
	text := strings.Repeat("large response ", 1<<20)
	connection := testConnection(t, Options{}, func(ctx context.Context, socket *wire.Conn) {
		if err := writeEvent(ctx, socket, testEvent{Text: text}); err != nil {
			t.Errorf("large write: %v", err)
		}
		_, _, _ = socket.Read(ctx)
	})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	event, err := connection.Recv(ctx)
	if err != nil || event.Text != text {
		t.Fatalf("large receive length = %d, error = %v, want %d", len(event.Text), err, len(text))
	}
}

func TestCloseTimeoutInterruptsUnresponsivePeer(t *testing.T) {
	release := make(chan struct{})
	connection := testConnection(t, Options{CloseTimeout: 50 * time.Millisecond}, func(_ context.Context, _ *wire.Conn) {
		<-release // Deliberately never read or acknowledge the close frame.
	})
	defer close(release)
	start := time.Now()
	if err := connection.Close(); err != nil {
		t.Fatalf("Close() = %v, want nil", err)
	}
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Errorf("Close() with a 50ms timeout took %v, want under one second", elapsed)
	}
	select {
	case <-connection.readerDone:
	default:
		t.Error("Close() returned before the reader exited")
	}
}

func TestReconnectDeadlineInterruptsRetirement(t *testing.T) {
	release := make(chan struct{})
	connection := testConnection(t, Options{CloseTimeout: 3 * time.Second}, func(_ context.Context, _ *wire.Conn) {
		<-release // Deliberately never acknowledge the close frame.
	})
	defer close(release)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	opened := false
	start := time.Now()
	next, err := connection.Reconnect(ctx, func(context.Context) (*Connection[testEvent], error) {
		opened = true
		return nil, errors.New("unexpected reconnect opener")
	})
	if next != nil || !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Reconnect() = %v, %v, want nil, context.DeadlineExceeded", next, err)
	}
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Errorf("Reconnect() with a 100ms deadline took %v, want under one second", elapsed)
	}
	if opened {
		t.Error("Reconnect() called the opener after the deadline expired")
	}
	select {
	case <-connection.readerDone:
	default:
		t.Error("Reconnect() returned before the old reader exited")
	}
}

func TestConcurrentSendAndClose(t *testing.T) {
	connection := testConnection(t, Options{}, func(ctx context.Context, socket *wire.Conn) {
		for {
			if _, _, err := socket.Read(ctx); err != nil {
				return
			}
		}
	})
	var workers sync.WaitGroup
	for range 64 {
		workers.Go(func() {
			err := connection.Send(context.Background(), []byte(`{"type":"response.create"}`))
			if err != nil {
				var delivery *DeliveryError
				if !errors.As(err, &delivery) {
					t.Errorf("send error type = %T", err)
				}
			}
		})
	}
	workers.Go(func() { _ = connection.Close() })
	workers.Go(func() { connection.Abort() })
	workers.Wait()
	if err := connection.Send(context.Background(), []byte(`{}`)); !errors.Is(err, ErrClosed) {
		t.Errorf("send after close = %v, want closed", err)
	}
}

func TestFragmentedMessageLimitBeforeFinalFrame(t *testing.T) {
	release := make(chan struct{})
	connection := testConnection(t, Options{MaxMessageBytes: 128}, func(ctx context.Context, socket *wire.Conn) {
		writer, err := socket.Writer(ctx, wire.MessageText)
		if err != nil {
			t.Errorf("writer: %v", err)
			return
		}
		// Writer emits non-final fragments once its buffer fills. The reader must
		// enforce the bound even when the final fragment is never sent.
		_, _ = writer.Write([]byte(strings.Repeat("x", 8192)))
		<-release
		_ = writer.Close()
	})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := connection.Recv(ctx)
	close(release)
	if !errors.Is(err, wire.ErrMessageTooBig) {
		t.Fatalf("fragmented limit = %v, want message-too-big", err)
	}
}

func TestOverflowRetainsAcceptedEventsInOrder(t *testing.T) {
	connection := testConnection(t, Options{MaxBufferedEvents: 2}, func(ctx context.Context, socket *wire.Conn) {
		for _, value := range []string{"one", "two", "overflow"} {
			if err := writeEvent(ctx, socket, testEvent{Text: value}); err != nil {
				return
			}
		}
		_, _, _ = socket.Read(ctx)
	})
	select {
	case <-connection.readerDone:
	case <-time.After(time.Second):
		t.Fatal("overflow did not stop reader")
	}
	for _, want := range []string{"one", "two"} {
		event, err := connection.Recv(context.Background())
		if err != nil || event.Text != want {
			t.Fatalf("accepted event = %+v,%v, want %q", event, err, want)
		}
	}
	if _, err := connection.Recv(context.Background()); !errors.Is(err, ErrBufferLimit) {
		t.Fatalf("overflow = %v, want buffer limit", err)
	}
}

func TestDefaultReceiveBufferSupportsSlowConsumer(t *testing.T) {
	const count = 2048
	ready := make(chan error, 1)
	connection := testConnection(t, Options{}, func(ctx context.Context, socket *wire.Conn) {
		for i := 0; i < count; i++ {
			if err := writeEvent(ctx, socket, testEvent{Text: "queued"}); err != nil {
				ready <- err
				return
			}
		}
		// A pong proves the reader has processed every preceding event while
		// the application deliberately has not called Recv yet.
		closed := socket.CloseRead(ctx)
		ready <- socket.Ping(ctx)
		<-closed.Done()
	})
	if err := <-ready; err != nil {
		t.Fatalf("default receive buffer interrupted stream: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	for i := 0; i < count; i++ {
		event, err := connection.Recv(ctx)
		if err != nil || event.Text != "queued" {
			t.Fatalf("event %d = %+v, %v", i, event, err)
		}
	}
}

func TestCanceledSendDoesNotReachServer(t *testing.T) {
	connection := testConnection(t, Options{}, func(ctx context.Context, socket *wire.Conn) {
		_, data, err := socket.Read(ctx)
		if err == nil {
			t.Errorf("cancelled send reached server: %q", data)
		}
	})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := connection.Send(ctx, []byte(`{"type":"response.create"}`))
	var delivery *DeliveryError
	if !errors.As(err, &delivery) || delivery.MayHaveBeenSent || !errors.Is(err, context.Canceled) {
		t.Fatalf("cancelled send = %v, want known-not-sent cancellation", err)
	}
	_ = connection.Close()
}

func TestCloseCancelsPendingReconnectUpgrade(t *testing.T) {
	connection := testConnection(t, Options{}, func(ctx context.Context, socket *wire.Conn) { _, _, _ = socket.Read(ctx) })
	requested := make(chan struct{})
	finished := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(requested)
		<-r.Context().Done()
		close(finished)
	}))
	defer server.Close()
	result := make(chan error, 1)
	go func() {
		_, err := connection.Reconnect(context.Background(), func(ctx context.Context) (*Connection[testEvent], error) {
			req, err := http.NewRequest(http.MethodGet, server.URL, nil)
			if err != nil {
				return nil, err
			}
			return Dial(ctx, req, server.Client(), Options{}, func([]byte) (testEvent, error) { return testEvent{}, nil }, func(testEvent) string { return "" })
		})
		result <- err
	}()
	select {
	case <-requested:
	case <-time.After(time.Second):
		t.Fatal("replacement upgrade did not start")
	}
	_ = connection.Close()
	select {
	case err := <-result:
		if !errors.Is(err, ErrClosed) {
			t.Errorf("reconnect error = %v, want closed", err)
		}
	case <-time.After(time.Second):
		t.Fatal("close did not cancel replacement upgrade")
	}
	select {
	case <-finished:
	case <-time.After(time.Second):
		t.Fatal("replacement server did not observe cancellation")
	}
}

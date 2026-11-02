package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	wire "github.com/coder/websocket"
)

func TestIsolatedOverflowReclaimsBacklogAndRetainsFailedID(t *testing.T) {
	for _, dimension := range []string{"events", "bytes"} {
		t.Run(dimension, func(t *testing.T) {
			slowEvent := testEvent{Type: "delta", StreamID: "slow", Text: strings.Repeat("界", 80)}
			fastEvent := testEvent{Type: "done", StreamID: "fast", Text: "preserved metadata"}
			options := Options{IsolateLaneOverflow: true, MaxBufferedEvents: 2}
			if dimension == "bytes" {
				slowJSON, err := json.Marshal(slowEvent)
				if err != nil {
					t.Fatal(err)
				}
				fastJSON, err := json.Marshal(fastEvent)
				if err != nil {
					t.Fatal(err)
				}
				options.MaxBufferedEvents = 0
				options.MaxBufferedBytes = int64(2*len(slowJSON) + len(fastJSON) - 1)
			}
			connection := testConnection(t, options, func(ctx context.Context, socket *wire.Conn) {
				for _, frames := range [][]testEvent{
					{slowEvent, slowEvent, fastEvent},
					{{StreamID: "slow", Text: "late failed-lane frame"}, {Type: "default", Text: "still open"}},
				} {
					if _, _, err := socket.Read(ctx); err != nil {
						t.Error(err)
						return
					}
					for _, event := range frames {
						if err := writeEvent(ctx, socket, event); err != nil {
							t.Error(err)
							return
						}
					}
				}
				_, _, _ = socket.Read(ctx)
			})
			slow, err := connection.Lane("slow")
			if err != nil {
				t.Fatal(err)
			}
			fast, err := connection.Lane("fast")
			if err != nil {
				t.Fatal(err)
			}
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := connection.Send(ctx, []byte("first")); err != nil {
				t.Fatal(err)
			}
			// Receiving the final frame proves the older backlog was already considered.
			if got, err := fast.Recv(ctx); err != nil || got != fastEvent {
				t.Fatalf("fast = %+v, %v", got, err)
			}
			_, first := slow.Recv(ctx)
			if !errors.Is(first, ErrBufferLimit) {
				t.Fatalf("slow = %v", first)
			}
			slow.Close()
			if _, err := slow.Recv(ctx); !errors.Is(err, first) {
				t.Fatalf("first cause changed: %v", err)
			}
			if _, err := connection.Lane("slow"); err == nil {
				t.Fatal("failed ID was reused")
			}
			if err := connection.Send(ctx, []byte("second")); err != nil {
				t.Fatal(err)
			}
			if got, err := connection.Recv(ctx); err != nil || got.Type != "default" || got.Text != "still open" {
				t.Fatalf("default = %+v, %v", got, err)
			}
		})
	}
}

func TestIsolatedOverflowUsesRegistrationOrderForTies(t *testing.T) {
	connection := testConnection(t, Options{IsolateLaneOverflow: true, MaxBufferedEvents: 2}, func(ctx context.Context, socket *wire.Conn) {
		if _, _, err := socket.Read(ctx); err != nil {
			t.Error(err)
			return
		}
		for _, key := range []string{"z-first", "a-second", "incoming"} {
			if err := writeEvent(ctx, socket, testEvent{StreamID: key, Text: key}); err != nil {
				t.Error(err)
				return
			}
		}
		_, _, _ = socket.Read(ctx)
	})
	lanes := make([]*Lane[testEvent], 0, 3)
	for _, key := range []string{"z-first", "a-second", "incoming"} {
		lane, err := connection.Lane(key)
		if err != nil {
			t.Fatal(err)
		}
		lanes = append(lanes, lane)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := connection.Send(ctx, []byte("start")); err != nil {
		t.Fatal(err)
	}
	if got, err := lanes[2].Recv(ctx); err != nil || got.Text != "incoming" {
		t.Fatalf("incoming = %+v, %v", got, err)
	}
	if _, err := lanes[0].Recv(ctx); !errors.Is(err, ErrBufferLimit) {
		t.Fatalf("first = %v", err)
	}
	if got, err := lanes[1].Recv(ctx); err != nil || got.Text != "a-second" {
		t.Fatalf("second = %+v, %v", got, err)
	}
}

func TestIsolatedOverflowCanFailOnlyTheDefaultLane(t *testing.T) {
	connection := testConnection(t, Options{IsolateLaneOverflow: true, MaxBufferedEvents: 1}, func(ctx context.Context, socket *wire.Conn) {
		if _, _, err := socket.Read(ctx); err != nil {
			t.Error(err)
			return
		}
		for _, event := range []testEvent{{Text: "older default"}, {StreamID: "named", Text: "survives"}} {
			if err := writeEvent(ctx, socket, event); err != nil {
				t.Error(err)
				return
			}
		}
		_, _, _ = socket.Read(ctx)
	})
	lane, err := connection.Lane("named")
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := connection.Send(ctx, []byte("start")); err != nil {
		t.Fatal(err)
	}
	if got, err := lane.Recv(ctx); err != nil || got.Text != "survives" {
		t.Fatalf("named = %+v, %v", got, err)
	}
	if _, err := connection.Recv(ctx); !errors.Is(err, ErrBufferLimit) {
		t.Fatalf("default = %v", err)
	}
	if err := connection.Send(ctx, []byte("socket remains usable")); err != nil {
		t.Fatal(err)
	}
}

func TestIsolatedOversizedFrameWakesItsReceiver(t *testing.T) {
	connection := testConnection(t, Options{IsolateLaneOverflow: true, MaxBufferedBytes: 100}, func(ctx context.Context, socket *wire.Conn) {
		if _, _, err := socket.Read(ctx); err != nil {
			t.Error(err)
			return
		}
		if err := writeEvent(ctx, socket, testEvent{StreamID: "oversized", Text: strings.Repeat("x", 200)}); err != nil {
			t.Error(err)
			return
		}
		if err := writeEvent(ctx, socket, testEvent{Text: "healthy"}); err != nil {
			t.Error(err)
			return
		}
		_, _, _ = socket.Read(ctx)
	})
	lane, err := connection.Lane("oversized")
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result := make(chan error, 1)
	go func() { _, err := lane.Recv(ctx); result <- err }()
	if err := connection.Send(ctx, []byte("start")); err != nil {
		t.Fatal(err)
	}
	if err := <-result; !errors.Is(err, ErrBufferLimit) {
		t.Fatalf("receiver = %v", err)
	}
	if got, err := connection.Recv(ctx); err != nil || got.Text != "healthy" {
		t.Fatalf("healthy = %+v, %v", got, err)
	}
}

func TestIsolationWithZeroLimitsRemainsUnlimited(t *testing.T) {
	connection := testConnection(t, Options{IsolateLaneOverflow: true}, func(ctx context.Context, socket *wire.Conn) {
		if _, _, err := socket.Read(ctx); err != nil {
			t.Error(err)
			return
		}
		for range 20 {
			if err := writeEvent(ctx, socket, testEvent{Text: strings.Repeat("x", 1024)}); err != nil {
				t.Error(err)
				return
			}
		}
		_, _, _ = socket.Read(ctx)
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := connection.Send(ctx, []byte("start")); err != nil {
		t.Fatal(err)
	}
	for range 20 {
		if got, err := connection.Recv(ctx); err != nil || len(got.Text) != 1024 {
			t.Fatalf("unlimited receive = %+v, %v", got, err)
		}
	}
}

func TestOverflowStillFailsTheConnectionByDefault(t *testing.T) {
	connection := testConnection(t, Options{MaxBufferedEvents: 1}, func(ctx context.Context, socket *wire.Conn) {
		if _, _, err := socket.Read(ctx); err != nil {
			t.Error(err)
			return
		}
		for _, event := range []testEvent{{StreamID: "first", Text: "accepted"}, {StreamID: "second", Text: "overflow"}} {
			if err := writeEvent(ctx, socket, event); err != nil {
				t.Error(err)
				return
			}
		}
		_, _, _ = socket.Read(ctx)
	})
	first, err := connection.Lane("first")
	if err != nil {
		t.Fatal(err)
	}
	second, err := connection.Lane("second")
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := connection.Send(ctx, []byte("start")); err != nil {
		t.Fatal(err)
	}
	if _, err := second.Recv(ctx); !errors.Is(err, ErrBufferLimit) {
		t.Fatalf("second = %v", err)
	}
	if got, err := first.Recv(ctx); err != nil || got.Text != "accepted" {
		t.Fatalf("accepted = %+v, %v", got, err)
	}
	if _, err := first.Recv(ctx); !errors.Is(err, ErrBufferLimit) {
		t.Fatalf("first = %v", err)
	}
	if _, err := connection.Recv(ctx); !errors.Is(err, ErrBufferLimit) {
		t.Fatalf("default = %v", err)
	}
}

func TestIsolationKeepsDecodeAndPhysicalFailuresGlobal(t *testing.T) {
	for _, mode := range []string{"decode", "physical"} {
		t.Run(mode, func(t *testing.T) {
			connection := testConnection(t, Options{IsolateLaneOverflow: true}, func(ctx context.Context, socket *wire.Conn) {
				if _, _, err := socket.Read(ctx); err != nil {
					t.Error(err)
					return
				}
				if mode == "decode" {
					if err := socket.Write(ctx, wire.MessageText, []byte("{")); err != nil {
						t.Error(err)
						return
					}
					_, _, _ = socket.Read(ctx)
				} else {
					// A synthetic peer failure intentionally ends transport delivery.
					_ = socket.CloseNow()
				}
			})
			lane, err := connection.Lane("named")
			if err != nil {
				t.Fatal(err)
			}
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := connection.Send(ctx, []byte("start")); err != nil {
				t.Fatal(err)
			}
			_, first := lane.Recv(ctx)
			if first == nil || errors.Is(first, ErrBufferLimit) {
				t.Fatalf("global failure = %v", first)
			}
			if _, err := connection.Recv(ctx); !errors.Is(err, first) {
				t.Fatalf("default lost global cause: %v", err)
			}
			if err := connection.Send(ctx, []byte("must not write")); !errors.Is(err, first) {
				t.Fatalf("send lost global cause: %v", err)
			}
		})
	}
}

package ssestream

import (
	"bufio"
	"io"
	"strings"
	"testing"
)

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

func newTestDecoder(raw string) *eventStreamDecoder {
	rc := nopCloser{strings.NewReader(raw)}
	scn := bufio.NewScanner(rc)
	return &eventStreamDecoder{rc: rc, scn: scn}
}

func TestCommentOnlyBlockDoesNotDispatchEvent(t *testing.T) {
	// A comment line followed by an empty line should NOT produce an event.
	raw := ": ping\n\n"
	dec := newTestDecoder(raw)

	if dec.Next() {
		t.Fatalf("expected no event for comment-only block, got event: Type=%q Data=%q", dec.Event().Type, dec.Event().Data)
	}
	if dec.Err() != nil {
		t.Fatalf("unexpected error: %v", dec.Err())
	}
}

func TestCommentBeforeDataEventIsIgnored(t *testing.T) {
	// A comment followed by a real data event should only produce the data event.
	raw := ": keep-alive\n\ndata: {\"id\":\"1\"}\n\n"
	dec := newTestDecoder(raw)

	if !dec.Next() {
		t.Fatalf("expected an event but got none; err=%v", dec.Err())
	}
	evt := dec.Event()
	if evt.Type != "" {
		t.Errorf("expected empty event type, got %q", evt.Type)
	}
	expected := "{\"id\":\"1\"}\n"
	if string(evt.Data) != expected {
		t.Errorf("expected data %q, got %q", expected, string(evt.Data))
	}
}

func TestMultipleCommentsDoNotDispatchEvents(t *testing.T) {
	raw := ": ping\n\n: ping\n\ndata: hello\n\n"
	dec := newTestDecoder(raw)

	if !dec.Next() {
		t.Fatalf("expected an event but got none; err=%v", dec.Err())
	}
	evt := dec.Event()
	expected := "hello\n"
	if string(evt.Data) != expected {
		t.Errorf("expected data %q, got %q", expected, string(evt.Data))
	}

	if dec.Next() {
		t.Fatalf("expected no more events, got Type=%q Data=%q", dec.Event().Type, dec.Event().Data)
	}
}

func TestEventTypeOnlyDispatchesEvent(t *testing.T) {
	// An event with only a type (no data) should still dispatch per SSE spec,
	// since the event type was explicitly set.
	raw := "event: ping\n\n"
	dec := newTestDecoder(raw)

	if !dec.Next() {
		t.Fatalf("expected an event but got none; err=%v", dec.Err())
	}
	evt := dec.Event()
	if evt.Type != "ping" {
		t.Errorf("expected event type %q, got %q", "ping", evt.Type)
	}
	if len(evt.Data) != 0 {
		t.Errorf("expected empty data, got %q", string(evt.Data))
	}
}

func TestNormalDataEvent(t *testing.T) {
	raw := "data: {\"msg\":\"hi\"}\n\n"
	dec := newTestDecoder(raw)

	if !dec.Next() {
		t.Fatalf("expected an event but got none; err=%v", dec.Err())
	}
	evt := dec.Event()
	expected := "{\"msg\":\"hi\"}\n"
	if string(evt.Data) != expected {
		t.Errorf("expected data %q, got %q", expected, string(evt.Data))
	}
}

func TestEventWithTypeAndData(t *testing.T) {
	raw := "event: message\ndata: {\"text\":\"hello\"}\n\n"
	dec := newTestDecoder(raw)

	if !dec.Next() {
		t.Fatalf("expected an event but got none; err=%v", dec.Err())
	}
	evt := dec.Event()
	if evt.Type != "message" {
		t.Errorf("expected event type %q, got %q", "message", evt.Type)
	}
	expected := "{\"text\":\"hello\"}\n"
	if string(evt.Data) != expected {
		t.Errorf("expected data %q, got %q", expected, string(evt.Data))
	}
}

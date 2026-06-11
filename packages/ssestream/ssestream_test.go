package ssestream

import (
	"bytes"
	"io"
	"net/http"
	"testing"
)

// mockBody creates an io.ReadCloser from a string
func mockBody(s string) io.ReadCloser {
	return io.NopCloser(bytes.NewReader([]byte(s)))
}

func TestEventStreamDecoder_SkipsEmptyDataEvents(t *testing.T) {
	// SSE stream with a comment-only event (keep-alive) followed by a real data event.
	// The ": PROCESSING\n\n" block should be silently discarded per HTML5 SSE spec.
	input := ": PROCESSING\n\ndata: {\"id\":\"chatcmpl-1\"}\n\n"

	resp := &http.Response{
		Header: http.Header{"Content-Type": []string{"text/event-stream"}},
		Body:   mockBody(input),
	}
	dec := NewDecoder(resp)
	if dec == nil {
		t.Fatal("NewDecoder returned nil")
	}
	defer dec.Close()

	// First event should be the data event, not the empty comment event
	if !dec.Next() {
		t.Fatalf("expected Next()=true, got false, err=%v", dec.Err())
	}
	evt := dec.Event()
	if string(evt.Data) != "{\"id\":\"chatcmpl-1\"}\n" {
		t.Errorf("expected data event, got: %q", string(evt.Data))
	}

	// No more events
	if dec.Next() {
		t.Error("expected no more events")
	}
	if dec.Err() != nil {
		t.Errorf("unexpected error: %v", dec.Err())
	}
}

func TestEventStreamDecoder_MultipleCommentsThenData(t *testing.T) {
	input := ": PROCESSING\n\n: STILL PROCESSING\n\ndata: hello\n\n"

	resp := &http.Response{
		Header: http.Header{"Content-Type": []string{"text/event-stream"}},
		Body:   mockBody(input),
	}
	dec := NewDecoder(resp)
	defer dec.Close()

	if !dec.Next() {
		t.Fatalf("expected Next()=true, err=%v", dec.Err())
	}
	if string(dec.Event().Data) != "hello\n" {
		t.Errorf("expected 'hello\\n', got: %q", string(dec.Event().Data))
	}
}

func TestEventStreamDecoder_MultipleDataEventsPreserved(t *testing.T) {
	input := "data: first\n\ndata: second\n\n"

	resp := &http.Response{
		Header: http.Header{"Content-Type": []string{"text/event-stream"}},
		Body:   mockBody(input),
	}
	dec := NewDecoder(resp)
	defer dec.Close()

	if !dec.Next() {
		t.Fatal("first event expected")
	}
	if string(dec.Event().Data) != "first\n" {
		t.Errorf("first event: %q", string(dec.Event().Data))
	}

	if !dec.Next() {
		t.Fatal("second event expected")
	}
	if string(dec.Event().Data) != "second\n" {
		t.Errorf("second event: %q", string(dec.Event().Data))
	}
}

func TestStreamSkipsEmptyDataEvents(t *testing.T) {
	type testChunk struct {
		ID string `json:"id"`
	}

	input := ": PROCESSING\n\ndata: {\"id\":\"1\"}\n\n"

	resp := &http.Response{
		Header: http.Header{"Content-Type": []string{"text/event-stream"}},
		Body:   mockBody(input),
	}
	dec := NewDecoder(resp)
	stream := NewStream[testChunk](dec, nil)
	defer stream.Close()

	if !stream.Next() {
		t.Fatalf("expected stream event, err=%v", stream.Err())
	}
	if stream.Current().ID != "1" {
		t.Errorf("expected ID=1, got %q", stream.Current().ID)
	}
	if stream.Next() {
		t.Error("expected no more events")
	}
	if stream.Err() != nil {
		t.Errorf("unexpected error: %v", stream.Err())
	}
}

func TestStreamEmptyDataDoesNotCrash(t *testing.T) {
	// Previously this caused "unexpected end of JSON input" panic
	input := ": PROCESSING\n\n"

	resp := &http.Response{
		Header: http.Header{"Content-Type": []string{"text/event-stream"}},
		Body:   mockBody(input),
	}
	dec := NewDecoder(resp)
	stream := NewStream[map[string]any](dec, nil)
	defer stream.Close()

	if stream.Next() {
		t.Error("expected no events from comment-only stream")
	}
	if stream.Err() != nil {
		t.Errorf("unexpected error: %v", stream.Err())
	}
}

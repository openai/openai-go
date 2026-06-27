package ssestream

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestEventStreamDecoderEmitsFinalEventAtEOF(t *testing.T) {
	res := &http.Response{
		Header: http.Header{
			"Content-Type": []string{"text/event-stream"},
		},
		Body: io.NopCloser(strings.NewReader("event: update\ndata: hello")),
	}

	decoder := NewDecoder(res)
	if decoder == nil {
		t.Fatal("expected decoder")
	}
	if !decoder.Next() {
		t.Fatalf("expected final event at EOF, err=%v", decoder.Err())
	}

	evt := decoder.Event()
	if evt.Type != "update" {
		t.Fatalf("event type = %q, want %q", evt.Type, "update")
	}
	if got := string(evt.Data); got != "hello\n" {
		t.Fatalf("event data = %q, want %q", got, "hello\n")
	}

	if decoder.Next() {
		t.Fatal("unexpected extra event")
	}
	if err := decoder.Err(); err != nil {
		t.Fatalf("unexpected decoder error: %v", err)
	}
}

func TestStreamEmitsFinalEventAtEOF(t *testing.T) {
	res := &http.Response{
		Header: http.Header{
			"Content-Type": []string{"text/event-stream"},
		},
		Body: io.NopCloser(strings.NewReader("data: {\"value\":\"last\"}")),
	}

	type payload struct {
		Value string `json:"value"`
	}

	stream := NewStream[payload](NewDecoder(res), nil)
	if !stream.Next() {
		t.Fatalf("expected final event at EOF, err=%v", stream.Err())
	}
	if got := stream.Current().Value; got != "last" {
		t.Fatalf("stream value = %q, want %q", got, "last")
	}

	if stream.Next() {
		t.Fatal("unexpected extra stream event")
	}
	if err := stream.Err(); err != nil {
		t.Fatalf("unexpected stream error: %v", err)
	}
}

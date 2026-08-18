package ssestream

import (
	"errors"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestStreamSkipsEventsWithoutData(t *testing.T) {
	tests := map[string]struct {
		body string
		want []string
	}{
		"comment": {
			body: ": OPENROUTER PROCESSING\n\ndata: {\"value\":\"first\"}\n\ndata: [DONE]\n\n",
			want: []string{"first"},
		},
		"retry directive": {
			body: "retry: 3000\n\ndata: {\"value\":\"first\"}\n\ndata: [DONE]\n\n",
			want: []string{"first"},
		},
		"CRLF comment": {
			body: ": OPENROUTER PROCESSING\r\n\r\ndata: {\"value\":\"first\"}\r\n\r\ndata: [DONE]\r\n\r\n",
			want: []string{"first"},
		},
		"multiple empty events": {
			body: "data: {\"value\":\"first\"}\n\n: keep-alive\n\nretry: 3000\n\ndata: {\"value\":\"second\"}\n\ndata: [DONE]\n\n",
			want: []string{"first", "second"},
		},
		"only empty events": {
			body: ": keep-alive\n\nretry: 3000\n\n",
			want: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res := &http.Response{
				Header: http.Header{"Content-Type": []string{"text/event-stream"}},
				Body:   io.NopCloser(strings.NewReader(test.body)),
			}
			stream := NewStream[struct {
				Value string `json:"value"`
			}](NewDecoder(res), nil)

			var got []string
			for stream.Next() {
				got = append(got, stream.Current().Value)
			}

			if err := stream.Err(); err != nil {
				t.Fatalf("stream ended with error: %v", err)
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("received values %v, want %v", got, test.want)
			}
		})
	}
}

func TestDecoderSkipsBlocksWithoutData(t *testing.T) {
	decoder := NewDecoder(&http.Response{
		Header: http.Header{"Content-Type": []string{"text/event-stream"}},
		Body: io.NopCloser(strings.NewReader(
			": keep-alive\n\nretry: 3000\n\ndata: {\"value\":\"first\"}\n\n",
		)),
	})

	if !decoder.Next() {
		t.Fatalf("decoder stopped before data event: %v", decoder.Err())
	}
	if got, want := string(decoder.Event().Data), "{\"value\":\"first\"}\n"; got != want {
		t.Fatalf("event data = %q, want %q", got, want)
	}
	if decoder.Next() {
		t.Fatalf("unexpected additional event: %+v", decoder.Event())
	}
	if err := decoder.Err(); err != nil {
		t.Fatalf("decoder ended with error: %v", err)
	}
}

type testDecoder struct {
	events  []Event
	current Event
	next    int
}

func (d *testDecoder) Next() bool {
	if d.next == len(d.events) {
		return false
	}
	d.current = d.events[d.next]
	d.next++
	return true
}

func (d *testDecoder) Event() Event { return d.current }
func (d *testDecoder) Close() error { return nil }
func (d *testDecoder) Err() error   { return nil }

func TestSynthesizedStreamPreservesCustomEventsWithoutData(t *testing.T) {
	decoder := &testDecoder{events: []Event{{Type: "custom.heartbeat"}}}
	stream := NewStreamWithSynthesizeEventData[struct {
		Event string `json:"event"`
		Data  any    `json:"data"`
	}](decoder, nil)

	if !stream.Next() {
		t.Fatalf("stream stopped before custom event: %v", stream.Err())
	}
	if got := stream.Current(); got.Event != "custom.heartbeat" || got.Data != nil {
		t.Fatalf("custom event = %#v, want event custom.heartbeat with nil data", got)
	}
}

func TestStreamPreservesErrorAfterBlockWithoutData(t *testing.T) {
	stream := NewStream[map[string]any](NewDecoder(&http.Response{
		Header: http.Header{"Content-Type": []string{"text/event-stream"}},
		Body: io.NopCloser(strings.NewReader(
			": keep-alive\n\ndata: {\"error\":{\"message\":\"bad\"}}\n\n",
		)),
	}), nil)

	if stream.Next() {
		t.Fatal("error event unexpectedly produced a stream value")
	}
	var streamErr *StreamError
	if !errors.As(stream.Err(), &streamErr) {
		t.Fatalf("stream error = %v, want *StreamError", stream.Err())
	}
}

var errTestReader = errors.New("test reader error")

type testErrorReadCloser struct {
	*strings.Reader
}

func (r *testErrorReadCloser) Read(p []byte) (int, error) {
	if r.Len() == 0 {
		return 0, errTestReader
	}
	return r.Reader.Read(p)
}

func (r *testErrorReadCloser) Close() error { return nil }

func TestStreamPreservesReaderErrorAfterBlockWithoutData(t *testing.T) {
	stream := NewStream[map[string]any](NewDecoder(&http.Response{
		Header: http.Header{"Content-Type": []string{"text/event-stream"}},
		Body: &testErrorReadCloser{
			Reader: strings.NewReader(": keep-alive\n\n"),
		},
	}), nil)

	if stream.Next() {
		t.Fatal("reader error unexpectedly produced a stream value")
	}
	if !errors.Is(stream.Err(), errTestReader) {
		t.Fatalf("stream error = %v, want %v", stream.Err(), errTestReader)
	}
}

func TestStreamRejectsEmptyDataField(t *testing.T) {
	stream := NewStream[map[string]any](NewDecoder(&http.Response{
		Header: http.Header{"Content-Type": []string{"text/event-stream"}},
		Body:   io.NopCloser(strings.NewReader("data:\n\n")),
	}), nil)

	if stream.Next() {
		t.Fatal("empty data field unexpectedly produced a stream value")
	}
	if err := stream.Err(); err == nil || !strings.Contains(err.Error(), "unexpected end of JSON input") {
		t.Fatalf("stream error = %v, want unexpected end of JSON input", err)
	}
}

func TestEventStreamDecoderDiscardsIncompleteEventAtEOF(t *testing.T) {
	for name, body := range map[string]string{
		"event with data":    "event: update\ndata: hello",
		"event without data": "event: update",
	} {
		t.Run(name, func(t *testing.T) {
			res := &http.Response{
				Header: http.Header{
					"Content-Type": []string{"text/event-stream"},
				},
				Body: io.NopCloser(strings.NewReader(body)),
			}

			decoder := NewDecoder(res)
			if decoder == nil {
				t.Fatal("expected decoder")
			}
			if decoder.Next() {
				t.Fatal("unexpected incomplete event")
			}
			if err := decoder.Err(); err != nil {
				t.Fatalf("unexpected decoder error: %v", err)
			}
		})
	}
}

func TestStreamDiscardsIncompleteEventAtEOF(t *testing.T) {
	res := &http.Response{
		Header: http.Header{
			"Content-Type": []string{"text/event-stream"},
		},
		Body: io.NopCloser(strings.NewReader("event: update")),
	}

	type payload struct {
		Value string `json:"value"`
	}

	stream := NewStream[payload](NewDecoder(res), nil)
	if stream.Next() {
		t.Fatal("unexpected incomplete stream event")
	}
	if err := stream.Err(); err != nil {
		t.Fatalf("unexpected stream error: %v", err)
	}
}

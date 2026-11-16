package ssestream

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

type domainEvent struct {
	Type  string `json:"type"`
	Error *struct {
		Code string `json:"code"`
	} `json:"error,omitempty"`
}

type domainAccumulator struct{ events []domainEvent }

func (c *domainAccumulator) Accumulate(event domainEvent) { c.events = append(c.events, event) }

func TestEventsModeDeliversDomainErrorsAndSubsequentEvents(t *testing.T) {
	for _, accumulated := range []bool{false, true} {
		for _, named := range []bool{false, true} {
			name := "plain"
			if accumulated {
				name = "accumulated"
			}
			if named {
				name += "/named"
			} else {
				name += "/unnamed"
			}
			t.Run(name, func(t *testing.T) {
				wire := ""
				for _, event := range []string{`{"type":"error","error":{"code":"internal_error"}}`, `{"type":"agent.session.idle"}`} {
					if named {
						var value domainEvent
						_ = json.Unmarshal([]byte(event), &value)
						wire += "event: " + value.Type + "\n"
					}
					wire += "data: " + event + "\n\n"
				}
				body := &closeTrackingReadCloser{Reader: strings.NewReader(wire)}
				decoder := NewDecoder(&http.Response{Header: http.Header{"Content-Type": {"text/event-stream"}}, Body: body})
				collector := &domainAccumulator{}
				stream := NewStream[domainEvent](decoder, nil)
				if accumulated {
					stream = NewStreamWithBetaAccumulator[domainEvent](decoder, nil, collector)
				}
				stream.WithEventsMode()
				var got []domainEvent
				for stream.Next() {
					got = append(got, stream.Current())
				}
				if stream.Err() != nil {
					t.Fatal(stream.Err())
				}
				if len(got) != 2 || got[0].Type != "error" || got[0].Error == nil || got[0].Error.Code != "internal_error" || got[1].Type != "agent.session.idle" {
					t.Fatalf("events = %#v", got)
				}
				if accumulated && !reflect.DeepEqual(collector.events, got) {
					t.Fatalf("collector = %#v, want %#v", collector.events, got)
				}
				if body.closeCalls != 1 {
					t.Fatalf("body closed %d times", body.closeCalls)
				}
			})
		}
	}
}

func TestEventsModePreservesDecodeAndTransportFailures(t *testing.T) {
	for _, test := range []struct {
		name    string
		data    string
		failure error
		syntax  bool
	}{
		{name: "malformed domain event", data: `{"type":"error","error":`, syntax: true},
		{name: "transport failure after domain event", data: `{"type":"error","error":{"code":"internal_error"}}`, failure: io.ErrUnexpectedEOF},
	} {
		t.Run(test.name, func(t *testing.T) {
			decoder := &testDecoder{events: []Event{{Data: []byte(test.data)}}, err: test.failure}
			stream := NewStream[domainEvent](decoder, nil).WithEventsMode()
			if test.syntax {
				if stream.Next() {
					t.Fatal("malformed JSON produced a value")
				}
				var syntax *json.SyntaxError
				if !errors.As(stream.Err(), &syntax) {
					t.Fatalf("error = %v, want JSON syntax error", stream.Err())
				}
			} else {
				if !stream.Next() || stream.Current().Type != "error" {
					t.Fatalf("domain event was not delivered: %v", stream.Err())
				}
				if stream.Next() || !errors.Is(stream.Err(), test.failure) {
					t.Fatalf("error = %v, want %v", stream.Err(), test.failure)
				}
			}
			if decoder.closeCalls != 1 {
				t.Fatalf("decoder closed %d times", decoder.closeCalls)
			}
		})
	}
}

func TestEventsModePreservesRequestCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = io.WriteString(w, "data: {\"type\":\"error\",\"error\":{\"code\":\"internal_error\"}}\n\n")
		w.(http.Flusher).Flush()
		<-r.Context().Done()
	}))
	defer server.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	request, _ := http.NewRequestWithContext(ctx, http.MethodGet, server.URL, nil)
	response, err := server.Client().Do(request)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if closeErr := response.Body.Close(); closeErr != nil {
			t.Errorf("close synthetic stream response: %v", closeErr)
		}
	}()
	stream := NewStream[domainEvent](NewDecoder(response), nil).WithEventsMode()
	if !stream.Next() || stream.Current().Type != "error" {
		t.Fatalf("domain event was not delivered: %v", stream.Err())
	}
	next := make(chan bool, 1)
	go func() { next <- stream.Next() }()
	cancel()
	select {
	case value := <-next:
		if value || !errors.Is(stream.Err(), context.Canceled) {
			t.Fatalf("Next=%v, error=%v", value, stream.Err())
		}
	case <-time.After(time.Second):
		t.Fatal("request cancellation did not interrupt the stream")
	}
}

func TestEventsModeUsesRegisteredDecoder(t *testing.T) {
	const contentType = "application/x-events-mode-test"
	decoder := &testDecoder{events: []Event{
		{Data: []byte(`{"type":"error","error":{"code":"internal_error"}}`)},
		{Data: []byte(`{"type":"agent.session.idle"}`)},
	}}
	RegisterDecoder(contentType, func(io.ReadCloser) Decoder { return decoder })
	t.Cleanup(func() { delete(decoderTypes, contentType) })
	stream := NewStream[domainEvent](NewDecoder(&http.Response{
		Header: http.Header{"Content-Type": {contentType}},
		Body:   io.NopCloser(strings.NewReader("custom wire format")),
	}), nil).WithEventsMode()
	var got []string
	for stream.Next() {
		got = append(got, stream.Current().Type)
	}
	if !reflect.DeepEqual(got, []string{"error", "agent.session.idle"}) || stream.Err() != nil {
		t.Fatalf("events=%v error=%v", got, stream.Err())
	}
	if decoder.closeCalls != 1 {
		t.Fatalf("registered decoder closed %d times", decoder.closeCalls)
	}
}

type domainIterator struct {
	values  []domainEvent
	next    int
	failure error
	closed  int
}

func (i *domainIterator) Next() bool {
	if i.next == len(i.values) {
		return false
	}
	i.next++
	return true
}
func (i *domainIterator) Current() domainEvent { return i.values[i.next-1] }
func (i *domainIterator) Err() error {
	if i.next == len(i.values) {
		return i.failure
	}
	return nil
}
func (i *domainIterator) Close() error { i.closed++; return nil }

func TestEventsModePreservesTypedIterator(t *testing.T) {
	iterator := &domainIterator{values: []domainEvent{{Type: "error"}, {Type: "agent.session.idle"}}, failure: io.ErrUnexpectedEOF}
	stream := NewStreamWithBetaIterator[domainEvent](iterator, nil).WithEventsMode()
	var got []domainEvent
	for stream.Next() {
		got = append(got, stream.Current())
	}
	if !reflect.DeepEqual(got, iterator.values) || !errors.Is(stream.Err(), io.ErrUnexpectedEOF) {
		t.Fatalf("events=%v error=%v", got, stream.Err())
	}
	if iterator.closed != 1 {
		t.Fatalf("iterator closed %d times", iterator.closed)
	}
}

func TestEventsModeKeepsNamedThreadEventsAsRawData(t *testing.T) {
	for _, synthesize := range []bool{false, true} {
		name := "raw"
		if synthesize {
			name = "explicit synthesis"
		}
		t.Run(name, func(t *testing.T) {
			decoder := &testDecoder{events: []Event{{Type: "thread.error", Data: []byte(`{"type":"error","error":{"code":"internal_error"}}`)}}}
			stream := NewStream[map[string]any](decoder, nil)
			if synthesize {
				stream = NewStreamWithSynthesizeEventData[map[string]any](decoder, nil)
			}
			stream.WithEventsMode()
			if !stream.Next() {
				t.Fatal(stream.Err())
			}
			got := stream.Current()
			if synthesize {
				if got["event"] != "thread.error" || got["data"].(map[string]any)["type"] != "error" {
					t.Fatalf("synthesized event=%v", got)
				}
			} else if got["type"] != "error" || got["error"].(map[string]any)["code"] != "internal_error" {
				t.Fatalf("raw event=%v", got)
			}
			if stream.Next() || stream.Err() != nil {
				t.Fatalf("stream error=%v", stream.Err())
			}
		})
	}
}

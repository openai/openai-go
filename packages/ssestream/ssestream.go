package ssestream

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"

	shimjson "github.com/openai/openai-go/v3/internal/encoding/json"
	"github.com/tidwall/gjson"
)

type Decoder interface {
	Event() Event
	Next() bool
	Close() error
	Err() error
}

func NewDecoder(res *http.Response) Decoder {
	if res == nil || res.Body == nil {
		return nil
	}

	contentType, mediaType := decoderContentTypes(res.Header.Get("content-type"))
	if t, ok := decoderTypes[contentType]; ok {
		return t(res.Body)
	}

	// Preserve parameter-specific registrations while allowing a bare media
	// type registration to match standard Content-Type parameters.
	if mediaType != "" {
		if t, ok := decoderTypes[mediaType]; ok {
			return t(res.Body)
		}
	}

	scn := bufio.NewScanner(res.Body)
	scn.Buffer(nil, bufio.MaxScanTokenSize<<9)
	return &eventStreamDecoder{rc: res.Body, scn: scn}
}

var decoderTypes = map[string](func(io.ReadCloser) Decoder){}

func RegisterDecoder(contentType string, decoder func(io.ReadCloser) Decoder) {
	decoderTypes[strings.ToLower(contentType)] = decoder
}

func decoderContentTypes(contentType string) (string, string) {
	base, _, _ := strings.Cut(contentType, ";")
	exactType := strings.ToLower(base) + contentType[len(base):]

	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return exactType, ""
	}
	return exactType, mediaType
}

type Event struct {
	Type string
	Data []byte
}

// StreamError represents an error event that occurred during streaming,
// preserving the original event data for structured access.
type StreamError struct {
	Message string
	Event   Event
}

func (e *StreamError) Error() string {
	return e.Message
}

// A base implementation of a Decoder for text/event-stream.
type eventStreamDecoder struct {
	evt Event
	rc  io.ReadCloser
	scn *bufio.Scanner
	err error
}

func (s *eventStreamDecoder) Next() bool {
	if s.err != nil {
		return false
	}

	event := ""
	var data []byte

	for s.scn.Scan() {
		txt := s.scn.Bytes()

		// Dispatch event on an empty line
		if len(txt) == 0 {
			if len(data) == 0 {
				event = ""
				continue
			}
			s.evt = Event{
				Type: event,
				Data: data,
			}
			return true
		}

		// Split a string like "event: bar" into name="event" and value=" bar".
		name, value, _ := bytes.Cut(txt, []byte(":"))

		// Consume an optional space after the colon if it exists.
		if len(value) > 0 && value[0] == ' ' {
			value = value[1:]
		}

		switch string(name) {
		case "":
			// An empty line in the for ": something" is a comment and should be ignored.
			continue
		case "event":
			event = string(value)
		case "data":
			data = append(data, value...)
			data = append(data, '\n')
		}
	}

	if s.scn.Err() != nil {
		s.err = s.scn.Err()
	}

	return false
}

func (s *eventStreamDecoder) Event() Event {
	return s.evt
}

func (s *eventStreamDecoder) Close() error {
	return s.rc.Close()
}

func (s *eventStreamDecoder) Err() error {
	return s.err
}

type Stream[T any] struct {
	decoder             Decoder
	cur                 T
	err                 error
	closeErr            error
	closeOnce           sync.Once
	done                atomic.Bool
	synthesizeEventData bool
}

func NewStream[T any](decoder Decoder, err error) *Stream[T] {
	return &Stream[T]{
		decoder: decoder,
		err:     err,
	}
}

func NewStreamWithSynthesizeEventData[T any](decoder Decoder, err error) *Stream[T] {
	return &Stream[T]{
		decoder:             decoder,
		err:                 err,
		synthesizeEventData: true,
	}
}

// Next returns false if the stream has ended or an error occurred.
// Call Stream.Current() to get the current value.
// Call Stream.Err() to get the error.
// The stream closes automatically when it reaches a terminal event or error.
// Call Stream.Close() if iteration stops before Next returns false.
//
//		for stream.Next() {
//			data := stream.Current()
//		}
//
//	 	if stream.Err() != nil {
//			...
//	 	}
func (s *Stream[T]) Next() bool {
	if s.err != nil {
		return s.finish(s.err)
	}
	decoder := s.decoder
	if s.done.Load() || decoder == nil {
		return false
	}

	if !decoder.Next() {
		// decoder.Next() may be false because of an error
		return s.finish(decoder.Err())
	}

	event := decoder.Event()
	if bytes.HasPrefix(event.Data, []byte("[DONE]")) {
		return s.finish(nil)
	}

	ep := gjson.GetBytes(event.Data, "error")
	if ep.Exists() {
		return s.finish(&StreamError{
			Message: fmt.Sprintf("received error while streaming: %s", ep.String()),
			Event:   event,
		})
	}
	var nxt T
	data := event.Data
	if s.synthesizeEventData || strings.HasPrefix(event.Type, "thread.") {
		synthesized := map[string]any{
			"event": event.Type,
			"data":  json.RawMessage(data),
		}
		var err error
		data, err = shimjson.Marshal(synthesized)
		if err != nil {
			return s.finish(err)
		}
	}
	if err := json.Unmarshal(data, &nxt); err != nil {
		return s.finish(err)
	}
	s.cur = nxt
	return true
}

func (s *Stream[T]) finish(err error) bool {
	s.err = err
	_ = s.Close()
	return false
}

func (s *Stream[T]) Current() T {
	return s.cur
}

func (s *Stream[T]) Err() error {
	return s.err
}

// Close releases the stream's decoder. Repeated calls return the first close
// result without closing the decoder again.
func (s *Stream[T]) Close() error {
	s.closeOnce.Do(func() {
		s.done.Store(true)
		if s.decoder != nil {
			s.closeErr = s.decoder.Close()
		}
	})
	return s.closeErr
}

// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package ssestream

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

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

	var decoder Decoder
	contentType := res.Header.Get("content-type")
	if t, ok := decoderTypes[contentType]; ok {
		decoder = t(res.Body)
	} else {
		decoder = &eventStreamDecoder{rc: res.Body, rdr: bufio.NewReader(res.Body)}
	}
	return decoder
}

var decoderTypes = map[string](func(io.ReadCloser) Decoder){}

func RegisterDecoder(contentType string, decoder func(io.ReadCloser) Decoder) {
	decoderTypes[strings.ToLower(contentType)] = decoder
}

type Event struct {
	Type string
	Data []byte
}

// A base implementation of a Decoder for text/event-stream.
type eventStreamDecoder struct {
	evt Event
	rc  io.ReadCloser
	rdr *bufio.Reader
	err error
}

func line(r *bufio.Reader) ([]byte, error) {
	var overflow bytes.Buffer

	// To prevent infinite loops, the failsafe stops when a line is
	// 100 times longer than the [io.Reader] default buffer size,
	// or after 10 failed attempts to find an end of line.
	for f := 0; f < 100; f++ {
		part, isPrefix, err := r.ReadLine()
		if err != nil {
			return nil, err
		}

		// Happy case, the line fits in the default buffer.
		if !isPrefix && overflow.Len() == 0 {
			return part, nil
		}

		// Overflow case, append to the buffer.
		if isPrefix || overflow.Len() > 0 {
			n, err := overflow.Write(part)
			if err != nil {
				return nil, err
			}

			// Didn't find an end of line, heavily increment the failsafe.
			if n != r.Size() {
				f += 10
			}
		}

		if !isPrefix {
			return overflow.Bytes(), nil
		}
	}

	return nil, fmt.Errorf("ssestream: too many attempts to read a line")
}

func (s *eventStreamDecoder) Next() bool {
	if s.err != nil {
		return false
	}

	event := ""
	data := bytes.NewBuffer(nil)

	for {
		txt, err := line(s.rdr)
		if err == io.EOF {
			return false
		}

		if err != nil {
			s.err = err
			break
		}

		// Dispatch event on an empty line
		if len(txt) == 0 {
			s.evt = Event{
				Type: event,
				Data: data.Bytes(),
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
			_, s.err = data.Write(value)
			if s.err != nil {
				break
			}
			_, s.err = data.WriteRune('\n')
			if s.err != nil {
				break
			}
		}
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
	decoder Decoder
	cur     T
	err     error
	done    bool
}

func NewStream[T any](decoder Decoder, err error) *Stream[T] {
	return &Stream[T]{
		decoder: decoder,
		err:     err,
	}
}

// Next returns false if the stream has ended or an error occurred.
// Call Stream.Current() to get the current value.
// Call Stream.Err() to get the error.
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
		return false
	}

	for s.decoder.Next() {
		if s.done {
			continue
		}

		if bytes.HasPrefix(s.decoder.Event().Data, []byte("[DONE]")) {
			// In this case we don't break because we still want to iterate through the full stream.
			s.done = true
			continue
		}

		ep := gjson.GetBytes(s.decoder.Event().Data, "error")
		if ep.Exists() {
			s.err = fmt.Errorf("received error while streaming: %s", ep.String())
			return false
		}
		var nxt T
		s.err = json.Unmarshal(s.decoder.Event().Data, &nxt)
		if s.err != nil {
			return false
		}
		s.cur = nxt
		return true
	}

	// decoder.Next() may be false because of an error
	s.err = s.decoder.Err()

	return false
}

func (s *Stream[T]) Current() T {
	return s.cur
}

func (s *Stream[T]) Err() error {
	return s.err
}

func (s *Stream[T]) Close() error {
	if s.decoder == nil {
		// already closed
		return nil
	}
	return s.decoder.Close()
}

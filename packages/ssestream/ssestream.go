package ssestream

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"mime"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"

	shimjson "github.com/openai/openai-go/v3/internal/encoding/json"
	"github.com/tidwall/gjson"
)

const (
	// DefaultMaxEventBytes is the maximum aggregate size of a built-in SSE
	// event, including its parsed fields and logical line endings.
	DefaultMaxEventBytes = bufio.MaxScanTokenSize << 9
	// DefaultMaxEventLines is the maximum number of non-empty lines in a
	// built-in SSE event.
	DefaultMaxEventLines = 4096
)

// ErrEventTooLarge is returned when a built-in SSE event exceeds its
// configured aggregate byte or line limit.
var ErrEventTooLarge = errors.New("ssestream: event exceeds configured limit")

// DecoderOptions configures the built-in SSE decoder. Non-positive limits use
// their corresponding defaults. Registered custom decoders manage their own
// limits and ignore these options.
type DecoderOptions struct {
	// MaxEventBytes is the maximum aggregate size of a parsed SSE event.
	MaxEventBytes int
	// MaxEventLines is the maximum number of non-empty lines in a parsed SSE
	// event.
	MaxEventLines int
}

type Decoder interface {
	Event() Event
	Next() bool
	Close() error
	Err() error
}

func NewDecoder(res *http.Response) Decoder {
	return NewDecoderWithOptions(res, DecoderOptions{})
}

// NewDecoderWithOptions returns a decoder with limits applied to the built-in
// SSE parser. Registered custom decoders retain their existing behavior.
func NewDecoderWithOptions(res *http.Response, options DecoderOptions) Decoder {
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

	if options.MaxEventBytes <= 0 {
		options.MaxEventBytes = DefaultMaxEventBytes
	}
	if options.MaxEventLines <= 0 {
		options.MaxEventLines = DefaultMaxEventLines
	}
	maxScanTokenBytes := options.MaxEventBytes
	if maxScanTokenBytes < math.MaxInt {
		maxScanTokenBytes++
	}
	scn := bufio.NewScanner(res.Body)
	scn.Buffer(nil, maxScanTokenBytes)
	decoder := &eventStreamDecoder{
		rc:            res.Body,
		maxEventBytes: options.MaxEventBytes,
		maxEventLines: options.MaxEventLines,
	}
	scn.Split(decoder.scanLines)
	decoder.scn = scn
	return decoder
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
	evt            Event
	rc             io.ReadCloser
	scn            *bufio.Scanner
	err            error
	maxEventBytes  int
	maxEventLines  int
	remainingBytes int
}

func (s *eventStreamDecoder) scanLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if newline := bytes.IndexByte(data, '\n'); newline >= 0 {
		line := data[:newline]
		if !isEventDelimiter(line) && !logicalLineFits(line, s.remainingBytes) {
			return 0, nil, s.eventTooLargeError()
		}
		return bufio.ScanLines(data, atEOF)
	}

	if atEOF {
		if isEventDelimiter(data) {
			return bufio.ScanLines(data, atEOF)
		}
		if len(data) > s.remainingBytes {
			return 0, nil, s.eventTooLargeError()
		}
		return len(data), nil, nil
	}

	if len(data) > s.remainingBytes && !(s.remainingBytes == 0 && isEventDelimiter(data)) {
		return 0, nil, s.eventTooLargeError()
	}
	return 0, nil, nil
}

func (s *eventStreamDecoder) eventTooLargeError() error {
	return fmt.Errorf("%w: maximum event size is %d bytes", ErrEventTooLarge, s.maxEventBytes)
}

func isEventDelimiter(line []byte) bool {
	return len(line) == 0 || len(line) == 1 && line[0] == '\r'
}

// logicalLineFits counts either an omitted LF or a trailing CR as one logical
// line ending, matching bufio.ScanLines normalization and event accounting.
func logicalLineFits(line []byte, maxBytes int) bool {
	if len(line) > 0 && line[len(line)-1] == '\r' {
		return len(line) <= maxBytes
	}
	return len(line) < maxBytes
}

func (s *eventStreamDecoder) Next() bool {
	if s.err != nil {
		return false
	}

	s.remainingBytes = s.maxEventBytes
	event := ""
	var data []byte
	eventLines := 0

	for s.scn.Scan() {
		txt := s.scn.Bytes()

		// Dispatch event on an empty line
		if len(txt) == 0 {
			if len(data) == 0 {
				event = ""
				s.remainingBytes = s.maxEventBytes
				eventLines = 0
				continue
			}
			s.evt = Event{
				Type: event,
				Data: data,
			}
			return true
		}

		eventLines++
		if eventLines > s.maxEventLines {
			s.err = fmt.Errorf("%w: maximum event line count is %d", ErrEventTooLarge, s.maxEventLines)
			return false
		}

		lineBytes := len(txt) + 1
		if lineBytes > s.remainingBytes {
			s.err = s.eventTooLargeError()
			return false
		}
		s.remainingBytes -= lineBytes

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

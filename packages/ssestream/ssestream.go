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
	decoderTypes[decoderContentTypeKey(contentType)] = decoder
}

func decoderContentTypes(contentType string) (string, string) {
	exactType := decoderContentTypeKey(contentType)

	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return exactType, ""
	}
	return exactType, mediaType
}

// decoderContentTypeKey normalizes only MIME components whose case is
// semantically insignificant. Parameter values remain case-sensitive unless
// their parameter defines otherwise. Extended parameter percent-encoding is
// normalized without changing unescaped value bytes.
func decoderContentTypeKey(contentType string) string {
	base, params, found := strings.Cut(contentType, ";")
	if !found {
		return strings.ToLower(contentType)
	}
	normalizedBase := strings.ToLower(base)
	return normalizedBase + ";" + normalizeMediaParameterTail(normalizedBase, params)
}

func normalizeMediaParameterTail(mediaType string, params string) string {
	var normalized strings.Builder
	segmentStart := 0
	inQuotes := false
	escaped := false

	for i := 0; i <= len(params); i++ {
		if i == len(params) || (!inQuotes && params[i] == ';') {
			normalized.WriteString(normalizeMediaParameter(mediaType, params[segmentStart:i]))
			if i < len(params) {
				normalized.WriteByte(';')
			}
			segmentStart = i + 1
			continue
		}

		switch params[i] {
		case '\\':
			if inQuotes && !escaped {
				escaped = true
				continue
			}
		case '"':
			if !escaped {
				inQuotes = !inQuotes
			}
		}
		escaped = false
	}

	return normalized.String()
}

func normalizeMediaParameter(mediaType string, param string) string {
	equals := strings.IndexByte(param, '=')
	if equals < 0 {
		return param
	}

	namePart := param[:equals]
	nameStart, nameEnd := trimOWSBounds(namePart)
	if nameStart == nameEnd {
		return param
	}
	name := namePart[nameStart:nameEnd]
	logicalName := mediaParameterLogicalName(name)

	var normalized strings.Builder
	normalized.WriteString(namePart[:nameStart])
	normalized.WriteString(strings.ToLower(name))
	normalized.WriteString(namePart[nameEnd:])
	normalized.WriteByte('=')

	value := param[equals+1:]
	switch {
	case isCaseInsensitiveMediaParameterValue(mediaType, logicalName):
		if strings.HasSuffix(name, "*") {
			normalized.WriteString(normalizeCaseInsensitiveExtendedParameterValue(value))
		} else {
			valueStart, valueEnd := trimOWSBounds(value)
			normalized.WriteString(value[:valueStart])
			normalized.WriteString(strings.ToLower(value[valueStart:valueEnd]))
			normalized.WriteString(value[valueEnd:])
		}
	case strings.HasSuffix(name, "*"):
		normalized.WriteString(normalizeExtendedParameterValue(value))
	default:
		normalized.WriteString(value)
	}

	return normalized.String()
}

func mediaParameterLogicalName(name string) string {
	logicalName := strings.TrimSuffix(name, "*")
	section := strings.LastIndexByte(logicalName, '*')
	if section < 0 || !isRFC2231Section(logicalName[section+1:]) {
		return logicalName
	}
	return logicalName[:section]
}

func isRFC2231Section(section string) bool {
	if section == "0" {
		return true
	}
	if len(section) == 0 || section[0] < '1' || section[0] > '9' {
		return false
	}
	for i := 1; i < len(section); i++ {
		if section[i] < '0' || section[i] > '9' {
			return false
		}
	}
	return true
}

func isCaseInsensitiveMediaParameterValue(mediaType string, name string) bool {
	if strings.EqualFold(name, "charset") {
		return true
	}

	switch strings.ToLower(strings.TrimSpace(mediaType)) {
	case "message/external-body":
		switch strings.ToLower(name) {
		case "access-type", "permission", "mode":
			return true
		}
	case "text/plain":
		switch strings.ToLower(name) {
		case "format", "delsp":
			return true
		}
	}
	return false
}

func normalizeCaseInsensitiveExtendedParameterValue(value string) string {
	valueStart, valueEnd := trimOWSBounds(value)
	core := value[valueStart:valueEnd]
	if strings.HasPrefix(core, "\"") {
		return value
	}

	firstQuote := strings.IndexByte(core, '\'')
	secondQuote := -1
	if firstQuote >= 0 {
		if offset := strings.IndexByte(core[firstQuote+1:], '\''); offset >= 0 {
			secondQuote = firstQuote + 1 + offset
		}
	}

	var normalized string
	if firstQuote >= 0 && secondQuote >= 0 {
		normalized = strings.ToLower(core[:firstQuote]) + "'" +
			strings.ToLower(core[firstQuote+1:secondQuote]) + "'" +
			normalizeCaseInsensitiveExtendedData(core[secondQuote+1:])
	} else {
		normalized = normalizeCaseInsensitiveExtendedData(core)
	}

	return value[:valueStart] + normalized + value[valueEnd:]
}

func normalizeCaseInsensitiveExtendedData(value string) string {
	bytes := []byte(value)
	for i := 0; i < len(bytes); i++ {
		if bytes[i] == '%' && i+2 < len(bytes) && isHexDigit(bytes[i+1]) && isHexDigit(bytes[i+2]) {
			decoded := hexValue(bytes[i+1])<<4 | hexValue(bytes[i+2])
			if decoded >= 'A' && decoded <= 'Z' {
				decoded += 'a' - 'A'
			}
			bytes[i+1] = hexDigit(decoded >> 4)
			bytes[i+2] = hexDigit(decoded & 0x0f)
			i += 2
			continue
		}
		if bytes[i] >= 'A' && bytes[i] <= 'Z' {
			bytes[i] += 'a' - 'A'
		}
	}
	return string(bytes)
}

func hexValue(value byte) byte {
	switch {
	case value >= '0' && value <= '9':
		return value - '0'
	case value >= 'a' && value <= 'f':
		return value - 'a' + 10
	default:
		return value - 'A' + 10
	}
}

func hexDigit(value byte) byte {
	if value < 10 {
		return '0' + value
	}
	return 'a' + value - 10
}

func normalizeExtendedParameterValue(value string) string {
	valueStart, valueEnd := trimOWSBounds(value)
	core := value[valueStart:valueEnd]
	if strings.HasPrefix(core, "\"") {
		return value
	}

	firstQuote := strings.IndexByte(core, '\'')
	secondQuote := -1
	if firstQuote >= 0 {
		if offset := strings.IndexByte(core[firstQuote+1:], '\''); offset >= 0 {
			secondQuote = firstQuote + 1 + offset
		}
	}

	var normalized string
	if firstQuote >= 0 && secondQuote >= 0 {
		normalized = strings.ToLower(core[:firstQuote]) + "'" +
			strings.ToLower(core[firstQuote+1:secondQuote]) + "'" +
			normalizePercentEncoding(core[secondQuote+1:])
	} else {
		normalized = normalizePercentEncoding(core)
	}

	return value[:valueStart] + normalized + value[valueEnd:]
}

func normalizePercentEncoding(value string) string {
	bytes := []byte(value)
	for i := 0; i+2 < len(bytes); i++ {
		if bytes[i] != '%' || !isHexDigit(bytes[i+1]) || !isHexDigit(bytes[i+2]) {
			continue
		}
		bytes[i+1] = lowerHexDigit(bytes[i+1])
		bytes[i+2] = lowerHexDigit(bytes[i+2])
		i += 2
	}
	return string(bytes)
}

func isHexDigit(value byte) bool {
	return value >= '0' && value <= '9' ||
		value >= 'a' && value <= 'f' ||
		value >= 'A' && value <= 'F'
}

func lowerHexDigit(value byte) byte {
	if value >= 'A' && value <= 'F' {
		return value + ('a' - 'A')
	}
	return value
}

func trimOWSBounds(value string) (int, int) {
	start := 0
	end := len(value)
	for start < end && (value[start] == ' ' || value[start] == '\t') {
		start++
	}
	for end > start && (value[end-1] == ' ' || value[end-1] == '\t') {
		end--
	}
	return start, end
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

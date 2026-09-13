package ssestream

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	shimjson "github.com/openai/openai-go/v3/internal/encoding/json"
	"github.com/tidwall/gjson"
	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/encoding/unicode/utf32"
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
		return asciiLower(contentType)
	}
	normalizedBase := asciiLower(base)
	externalBodyAccessType := ""
	externalBodyAccessLanguage := ""
	hasExternalBodyAccessType := false
	if strings.EqualFold(strings.TrimSpace(normalizedBase), "message/external-body") {
		externalBodyAccessType, externalBodyAccessLanguage, hasExternalBodyAccessType = parseExternalBodyAccessType(contentType, params)
	}
	return normalizedBase + ";" + normalizeMediaParameterTail(normalizedBase, params, externalBodyAccessType, externalBodyAccessLanguage, hasExternalBodyAccessType)
}

func parseExternalBodyAccessType(contentType string, params string) (string, string, bool) {
	if accessType, language, found, ok := decodeExtendedMediaParameter(params, "access-type"); found {
		if !ok || !isMIMEToken(accessType) {
			return "", "", false
		}
		return asciiLower(accessType), asciiLower(language), true
	}

	_, parsedParams, err := mime.ParseMediaType(contentType)
	if err != nil {
		return "", "", false
	}
	accessType, ok := parsedParams["access-type"]
	if !ok || !isMIMEToken(accessType) {
		return "", "", false
	}
	return asciiLower(accessType), "", true
}

type extendedMediaParameterSegment struct {
	encoded bool
	value   string
}

func decodeExtendedMediaParameter(params string, logicalName string) (string, string, bool, bool) {
	var single extendedMediaParameterSegment
	hasSingle := false
	sections := map[int]extendedMediaParameterSegment{}
	plainFound := false
	found := false
	valid := true

	forEachMediaParameter(params, func(param string) {
		equals := strings.IndexByte(param, '=')
		if !valid || equals < 0 {
			return
		}
		nameStart, nameEnd := trimOWSBounds(param[:equals])
		name := param[nameStart:nameEnd]
		if !strings.EqualFold(mediaParameterLogicalName(name), logicalName) {
			return
		}
		if strings.EqualFold(name, logicalName) {
			if found {
				valid = false
			}
			plainFound = true
			return
		}
		if plainFound {
			found = true
			valid = false
			return
		}
		found = true

		encoded := strings.HasSuffix(name, "*")
		sectionName := strings.TrimSuffix(name, "*")
		segment := extendedMediaParameterSegment{encoded: encoded, value: param[equals+1:]}
		if strings.EqualFold(sectionName, logicalName) {
			if hasSingle || len(sections) != 0 || !encoded {
				valid = false
				return
			}
			single = segment
			hasSingle = true
			return
		}

		star := strings.LastIndexByte(sectionName, '*')
		if star < 0 {
			valid = false
			return
		}
		section, err := strconv.Atoi(sectionName[star+1:])
		if err != nil || hasSingle {
			valid = false
			return
		}
		if _, duplicate := sections[section]; duplicate {
			valid = false
			return
		}
		sections[section] = segment
	})

	if !found {
		return "", "", false, false
	}
	if !valid {
		return "", "", true, false
	}
	if hasSingle {
		core, quoted, ok := mediaParameterValueCore(single.value)
		if !ok {
			return "", "", true, false
		}
		charset, language, data, ok := splitExtendedInitialValue(core)
		if !ok || (!quoted && !validRFC2231ExtendedData(data)) {
			return "", "", true, false
		}
		raw, ok := decodeExtendedOctets(data)
		if !ok {
			return "", "", true, false
		}
		decoded, ok := decodeMIMEParameterValue(charset, raw)
		return decoded, language, true, ok
	}
	if len(sections) == 0 {
		return "", "", true, false
	}

	var decoded strings.Builder
	var encodedRun []byte
	charset := ""
	language := ""
	flushEncoded := func() bool {
		if len(encodedRun) == 0 {
			return true
		}
		text, ok := decodeMIMEParameterValue(charset, encodedRun)
		if !ok {
			return false
		}
		decoded.WriteString(text)
		encodedRun = encodedRun[:0]
		return true
	}

	for section := 0; section < len(sections); section++ {
		segment, ok := sections[section]
		if !ok {
			return "", "", true, false
		}
		core, quoted, ok := mediaParameterValueCore(segment.value)
		if !ok {
			return "", "", true, false
		}
		data := core
		if section == 0 && segment.encoded {
			charset, language, data, ok = splitExtendedInitialValue(core)
			if !ok {
				return "", "", true, false
			}
		}

		if segment.encoded {
			if !quoted && !validRFC2231ExtendedData(data) {
				return "", "", true, false
			}
			octets, ok := decodeExtendedOctets(data)
			if !ok {
				return "", "", true, false
			}
			encodedRun = append(encodedRun, octets...)
			continue
		}

		if !quoted && !isMIMEToken(data) {
			return "", "", true, false
		}
		if !flushEncoded() {
			return "", "", true, false
		}
		decoded.WriteString(data)
	}
	if !flushEncoded() {
		return "", "", true, false
	}
	return decoded.String(), language, true, true
}

func mediaParameterValueCore(value string) (string, bool, bool) {
	valueStart, valueEnd := trimOWSBounds(value)
	core := value[valueStart:valueEnd]
	if !strings.HasPrefix(core, "\"") {
		return core, false, true
	}
	contents, ok := quotedMediaParameterContents(core)
	if !ok {
		return "", false, false
	}

	var decoded strings.Builder
	for i := 0; i < len(contents); i++ {
		if contents[i] == '\\' {
			if i+1 >= len(contents) {
				return "", false, false
			}
			i++
		}
		decoded.WriteByte(contents[i])
	}
	return decoded.String(), true, true
}

func splitExtendedInitialValue(value string) (string, string, string, bool) {
	firstQuote := strings.IndexByte(value, '\'')
	if firstQuote < 0 {
		return "", "", "", false
	}
	secondOffset := strings.IndexByte(value[firstQuote+1:], '\'')
	if secondOffset < 0 {
		return "", "", "", false
	}
	secondQuote := firstQuote + secondOffset + 1
	charset := value[:firstQuote]
	language := value[firstQuote+1 : secondQuote]
	if charset != "" {
		if charset != strings.TrimSpace(charset) {
			return "", "", "", false
		}
		if _, err := ianaindex.IANA.Encoding(charset); err != nil {
			return "", "", "", false
		}
	}
	if language != "" && !isRFC1766LanguageTag(language) {
		return "", "", "", false
	}
	return charset, language, value[secondQuote+1:], true
}

func validRFC2231ExtendedData(value string) bool {
	for i := 0; i < len(value); i++ {
		if value[i] == '%' {
			if i+2 >= len(value) || !isHexDigit(value[i+1]) || !isHexDigit(value[i+2]) {
				return false
			}
			i += 2
			continue
		}
		if !isRFC2231AttributeChar(value[i]) {
			return false
		}
	}
	return true
}

func isRFC2231AttributeChar(value byte) bool {
	return isMIMETokenChar(value) && value != '*' && value != '\'' && value != '%'
}

func isMIMEToken(value string) bool {
	if value == "" {
		return false
	}
	for i := 0; i < len(value); i++ {
		if !isMIMETokenChar(value[i]) {
			return false
		}
	}
	return true
}

func isMIMETokenChar(value byte) bool {
	if value <= ' ' || value >= 0x7f {
		return false
	}
	switch value {
	case '(', ')', '<', '>', '@', ',', ';', ':', '\\', '"', '/', '[', ']', '?', '=':
		return false
	default:
		return true
	}
}

func validMIMEMediaType(value string) bool {
	typePart, subtype, found := strings.Cut(value, "/")
	return found && !strings.Contains(subtype, "/") && isMIMEToken(typePart) && isMIMEToken(subtype)
}

func asciiLower(value string) string {
	bytes := []byte(value)
	for i := range bytes {
		if bytes[i] >= 'A' && bytes[i] <= 'Z' {
			bytes[i] += 'a' - 'A'
		}
	}
	return string(bytes)
}

func decodeExtendedOctets(value string) ([]byte, bool) {
	decoded := make([]byte, 0, len(value))
	for i := 0; i < len(value); i++ {
		if value[i] != '%' {
			decoded = append(decoded, value[i])
			continue
		}
		if i+2 >= len(value) || !isHexDigit(value[i+1]) || !isHexDigit(value[i+2]) {
			return nil, false
		}
		decoded = append(decoded, hexValue(value[i+1])<<4|hexValue(value[i+2]))
		i += 2
	}
	return decoded, true
}

func decodeMIMEParameterValue(charset string, value []byte) (string, bool) {
	if charset == "" {
		return string(value), true
	}

	decoderEncoding, err := ianaindex.IANA.Encoding(charset)
	if err != nil {
		return "", false
	}
	if decoderEncoding == nil {
		switch asciiLower(charset) {
		case "utf-32be", "csutf32be":
			decoderEncoding = utf32.UTF32(utf32.BigEndian, utf32.IgnoreBOM)
		case "utf-32le", "csutf32le":
			decoderEncoding = utf32.UTF32(utf32.LittleEndian, utf32.IgnoreBOM)
		case "utf-32", "csutf32":
			decoderEncoding = utf32.UTF32(utf32.BigEndian, utf32.ExpectBOM)
		default:
			return "", false
		}
	}
	decoded, err := decoderEncoding.NewDecoder().Bytes(value)
	if err != nil {
		return "", false
	}
	return string(decoded), true
}

func isRFC1766LanguageTag(language string) bool {
	// RFC 1766: Primary-tag = 1*8ALPHA, Subtag = 1*8ALPHA.
	// Digits are not valid in either the primary tag or its subtags.
	partLength := 0
	for i := 0; i <= len(language); i++ {
		if i == len(language) || language[i] == '-' {
			if partLength == 0 || partLength > 8 {
				return false
			}
			partLength = 0
			continue
		}
		c := language[i]
		if c < 'A' || c > 'Z' {
			if c < 'a' || c > 'z' {
				return false
			}
		}
		partLength++
	}
	return true
}

type decodedExtendedParameterState struct {
	value    string
	language string
	found    bool
	decoded  bool
}

func normalizeMediaParameterTail(mediaType string, params string, externalBodyAccessType string, externalBodyAccessLanguage string, hasExternalBodyAccessType bool) string {
	decodedParameters := map[string]decodedExtendedParameterState{}
	var normalized strings.Builder
	first := true
	forEachMediaParameter(params, func(param string) {
		if !first {
			normalized.WriteByte(';')
		}
		first = false
		normalized.WriteString(normalizeMediaParameter(mediaType, params, param, externalBodyAccessType, externalBodyAccessLanguage, hasExternalBodyAccessType, decodedParameters))
	})
	return normalized.String()
}

func forEachMediaParameter(params string, visit func(string)) {
	segmentStart := 0
	inQuotes := false
	escaped := false

	for i := 0; i <= len(params); i++ {
		if i == len(params) || (!inQuotes && params[i] == ';') {
			visit(params[segmentStart:i])
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
}

func normalizeMediaParameter(mediaType string, params string, param string, externalBodyAccessType string, externalBodyAccessLanguage string, hasExternalBodyAccessType bool, decodedParameters map[string]decodedExtendedParameterState) string {
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
	normalized.WriteString(asciiLower(name))
	normalized.WriteString(namePart[nameEnd:])
	normalized.WriteByte('=')

	value := param[equals+1:]
	isExternalBodyAccessType := strings.TrimSpace(mediaType) == "message/external-body" && strings.EqualFold(logicalName, "access-type")
	switch {
	case isExternalBodyAccessType:
		normalized.WriteString(normalizeCanonicalMediaParameterValue(name, value, "access-type", externalBodyAccessType, externalBodyAccessLanguage, hasExternalBodyAccessType))
	case isCaseInsensitiveMediaParameterValue(mediaType, logicalName, externalBodyAccessType):
		state := decodedCaseInsensitiveParameter(mediaType, params, logicalName, externalBodyAccessType, decodedParameters)
		if state.found {
			normalized.WriteString(normalizeCanonicalMediaParameterValue(name, value, logicalName, state.value, state.language, state.decoded))
		} else {
			core, _, ok := mediaParameterValueCore(value)
			if ok && validCaseInsensitiveMediaParameterValue(mediaType, logicalName, core, externalBodyAccessType) {
				valueStart, valueEnd := trimOWSBounds(value)
				normalized.WriteString(value[:valueStart])
				normalized.WriteString(asciiLower(value[valueStart:valueEnd]))
				normalized.WriteString(value[valueEnd:])
			} else {
				normalized.WriteString(value)
			}
		}
	case strings.HasSuffix(name, "*"):
		normalized.WriteString(normalizeExtendedParameterValue(value, extendedMediaParameterHasMetadata(name)))
	default:
		normalized.WriteString(value)
	}

	return normalized.String()
}

func decodedCaseInsensitiveParameter(mediaType string, params string, logicalName string, externalBodyAccessType string, cache map[string]decodedExtendedParameterState) decodedExtendedParameterState {
	key := asciiLower(logicalName)
	if state, ok := cache[key]; ok {
		return state
	}
	value, language, found, decoded := decodeExtendedMediaParameter(params, logicalName)
	if decoded {
		if validCaseInsensitiveMediaParameterValue(mediaType, logicalName, value, externalBodyAccessType) {
			value = asciiLower(value)
			language = asciiLower(language)
		} else {
			decoded = false
		}
	}
	state := decodedExtendedParameterState{value: value, language: language, found: found, decoded: decoded}
	cache[key] = state
	return state
}

func normalizeCanonicalMediaParameterValue(name string, value string, logicalName string, decodedValue string, decodedLanguage string, decoded bool) string {
	if decoded {
		canonical := "d"
		if isInitialMediaParameterSegment(name, logicalName) {
			canonical = encodeDecodedDecoderKeyValue(decodedLanguage, decodedValue)
		}
		valueStart, valueEnd := trimOWSBounds(value)
		return value[:valueStart] + canonical + value[valueEnd:]
	}

	normalizedValue := value
	if strings.HasSuffix(name, "*") {
		normalizedValue = normalizeExtendedParameterValue(value, extendedMediaParameterHasMetadata(name))
	}
	valueStart, valueEnd := trimOWSBounds(normalizedValue)
	return normalizedValue[:valueStart] + encodeDecoderKeyValue('r', normalizedValue[valueStart:valueEnd]) + normalizedValue[valueEnd:]
}

func encodeDecodedDecoderKeyValue(language string, value string) string {
	languageHex := make([]byte, hex.EncodedLen(len(language)))
	hex.Encode(languageHex, []byte(language))
	valueHex := make([]byte, hex.EncodedLen(len(value)))
	hex.Encode(valueHex, []byte(value))
	return "d" + string(languageHex) + "g" + string(valueHex)
}

func encodeDecoderKeyValue(prefix byte, value string) string {
	encoded := make([]byte, 1+hex.EncodedLen(len(value)))
	encoded[0] = prefix
	hex.Encode(encoded[1:], []byte(value))
	return string(encoded)
}

func isInitialMediaParameterSegment(name string, logicalName string) bool {
	if strings.EqualFold(name, logicalName) {
		return true
	}
	sectionName := strings.TrimSuffix(name, "*")
	if strings.EqualFold(sectionName, logicalName) {
		return true
	}
	star := strings.LastIndexByte(sectionName, '*')
	return star >= 0 && strings.EqualFold(sectionName[:star], logicalName) && sectionName[star+1:] == "0"
}

func mediaParameterLogicalName(name string) string {
	logicalName := strings.TrimSuffix(name, "*")
	section := strings.LastIndexByte(logicalName, '*')
	if section < 0 || !isRFC2231Section(logicalName[section+1:]) {
		return logicalName
	}
	return logicalName[:section]
}

func extendedMediaParameterHasMetadata(name string) bool {
	if !strings.HasSuffix(name, "*") {
		return false
	}

	encodedName := strings.TrimSuffix(name, "*")
	section := strings.LastIndexByte(encodedName, '*')
	if section < 0 {
		return true
	}
	return encodedName[section+1:] == "0"
}

func isRFC2231Section(section string) bool {
	// RFC 2231 section 3 explicitly forbids leading zeroes and gaps.
	// The initial section is exactly "0"; subsequent sections start 1-9.
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

func isCaseInsensitiveMediaParameterValue(mediaType string, name string, externalBodyAccessType string) bool {
	if strings.EqualFold(name, "charset") {
		return true
	}

	switch strings.TrimSpace(mediaType) {
	case "message/external-body":
		switch asciiLower(name) {
		case "access-type", "permission":
			return true
		case "mode":
			switch externalBodyAccessType {
			case "ftp", "anon-ftp", "tftp":
				return true
			}
		}
	case "multipart/encrypted":
		return strings.EqualFold(name, "protocol")
	case "multipart/signed":
		// RFC 1847 makes micalg value syntax and semantics protocol-defined.
		// RFC 2045 therefore leaves micalg case-sensitive unless that selected
		// protocol explicitly defines otherwise; do not fold it generically.
		return strings.EqualFold(name, "protocol")
	case "multipart/report":
		return strings.EqualFold(name, "report-type")
	case "multipart/related":
		return strings.EqualFold(name, "type")
	case "text/plain":
		switch asciiLower(name) {
		case "format", "delsp":
			return true
		}
	case "text/calendar":
		switch asciiLower(name) {
		case "method", "component":
			return true
		}
	}
	return false
}

func validCaseInsensitiveMediaParameterValue(mediaType string, name string, value string, externalBodyAccessType string) bool {
	if strings.EqualFold(name, "charset") {
		// RFC 2046 defines charset values as case-insensitive and permits
		// private X-* names; registry lookup is only required when decoding
		// RFC 2231 octets, not when comparing the charset parameter value.
		return isMIMEToken(value)
	}

	switch strings.TrimSpace(mediaType) {
	case "message/external-body":
		switch asciiLower(name) {
		case "access-type", "permission", "mode":
			return isMIMEToken(value)
		}
	case "multipart/encrypted", "multipart/signed":
		if strings.EqualFold(name, "protocol") {
			return validMIMEMediaType(value)
		}
	case "multipart/report":
		if strings.EqualFold(name, "report-type") {
			return isMIMEToken(value)
		}
	case "multipart/related":
		if strings.EqualFold(name, "type") {
			return validMIMEMediaType(value)
		}
	case "text/plain":
		switch asciiLower(name) {
		case "format", "delsp":
			return isMIMEToken(value)
		}
	case "text/calendar":
		switch asciiLower(name) {
		case "method", "component":
			return isMIMEToken(value)
		}
	}
	return false
}

func normalizeExtendedParameterValue(value string, hasMetadata bool) string {
	valueStart, valueEnd := trimOWSBounds(value)
	core := value[valueStart:valueEnd]
	quoted := false
	if strings.HasPrefix(core, "\"") {
		var ok bool
		core, ok = quotedMediaParameterContents(core)
		if !ok {
			return value
		}
		quoted = true
	}

	firstQuote := -1
	secondQuote := -1
	if hasMetadata {
		firstQuote = strings.IndexByte(core, '\'')
		if firstQuote >= 0 {
			if offset := strings.IndexByte(core[firstQuote+1:], '\''); offset >= 0 {
				secondQuote = firstQuote + 1 + offset
			}
		}
	}

	var normalized string
	if firstQuote >= 0 && secondQuote >= 0 {
		normalized = asciiLower(core[:firstQuote]) + "'" +
			asciiLower(core[firstQuote+1:secondQuote]) + "'" +
			normalizePercentEncoding(core[secondQuote+1:])
	} else {
		normalized = normalizePercentEncoding(core)
	}
	if quoted {
		normalized = "\"" + normalized + "\""
	}

	return value[:valueStart] + normalized + value[valueEnd:]
}

func quotedMediaParameterContents(value string) (string, bool) {
	if len(value) < 2 || value[0] != '"' || value[len(value)-1] != '"' {
		return "", false
	}

	escaped := false
	for i := 1; i < len(value)-1; i++ {
		switch value[i] {
		case '\\':
			escaped = !escaped
		case '"':
			if !escaped {
				return "", false
			}
			escaped = false
		default:
			escaped = false
		}
	}
	if escaped {
		return "", false
	}
	return value[1 : len(value)-1], true
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

package ssestream

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/mail"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	shimjson "github.com/openai/openai-go/v3/internal/encoding/json"
	"github.com/tidwall/gjson"
	textencoding "golang.org/x/text/encoding"
	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/encoding/unicode/utf32"
	"golang.org/x/text/transform"
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
	semanticMediaType := strings.TrimSpace(normalizedBase)
	externalBodyAccessType := ""
	externalBodyAccessLanguage := ""
	hasExternalBodyAccessType := false
	if semanticMediaType == "message/external-body" {
		externalBodyAccessType, externalBodyAccessLanguage, hasExternalBodyAccessType = parseExternalBodyAccessType(params)
	}
	return normalizedBase + ";" + normalizeMediaParameterTail(semanticMediaType, params, externalBodyAccessType, externalBodyAccessLanguage, hasExternalBodyAccessType)
}

func parseExternalBodyAccessType(params string) (string, string, bool) {
	if accessType, language, found, ok := decodeExtendedMediaParameter(params, "access-type"); found {
		if !ok || !isMIMEToken(accessType) {
			return "", "", false
		}
		return asciiLower(accessType), asciiLower(language), true
	}

	accessType, count, equal, valid := scanPlainMediaParameter(params, "access-type")
	if count == 0 || !equal || !valid || !isMIMEToken(accessType) {
		return "", "", false
	}
	return asciiLower(accessType), "", true
}

type extendedMediaParameterSegment struct {
	encoded bool
	value   string
}

type extendedMediaParameterSection struct {
	value   string
	section uint32
	encoded bool
}

func decodeExtendedMediaParameter(params string, logicalName string) (string, string, bool, bool) {
	var single extendedMediaParameterSegment
	singleCore := ""
	hasSingle := false
	sectionCount := 0
	sectionsOrdered := true
	plainValue := ""
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
			if plainFound {
				duplicate, quoted, ok := mediaParameterValueCore(param[equals+1:])
				if !ok || (!quoted && !isMIMEToken(duplicate)) || (quoted && strings.ContainsAny(duplicate, "\r\n")) || duplicate != plainValue {
					valid = false
				}
				return
			}
			parsed, ok := parsePlainMediaParameterValue(param[equals+1:])
			if !ok {
				valid = false
				return
			}
			plainValue = parsed
			plainFound = true
			return
		}
		found = true

		encoded := strings.HasSuffix(name, "*")
		sectionName := strings.TrimSuffix(name, "*")
		segment := extendedMediaParameterSegment{encoded: encoded, value: param[equals+1:]}
		if strings.EqualFold(sectionName, logicalName) {
			if sectionCount != 0 || !encoded {
				valid = false
				return
			}
			core, _, ok := mediaParameterValueCore(segment.value)
			if !ok {
				valid = false
				return
			}
			if hasSingle {
				if core != singleCore {
					valid = false
				}
				return
			}
			single = segment
			singleCore = core
			hasSingle = true
			return
		}

		star := strings.LastIndexByte(sectionName, '*')
		if star < 0 {
			valid = false
			return
		}
		sectionValue, err := strconv.ParseUint(sectionName[star+1:], 10, 32)
		if err != nil || hasSingle {
			valid = false
			return
		}
		section := uint32(sectionValue)
		if section != uint32(sectionCount) {
			sectionsOrdered = false
		}
		sectionCount++
	})

	if !found {
		return "", "", false, false
	}
	if !valid {
		return "", "", true, false
	}
	if hasSingle {
		fallbackPlain := func() (string, string, bool, bool) {
			if plainFound {
				return plainValue, "", true, true
			}
			return "", "", true, false
		}
		core, quoted, ok := mediaParameterValueCore(single.value)
		if !ok {
			return "", "", true, false
		}
		charset, language, data, ok := splitExtendedInitialValue(core)
		if !ok {
			return fallbackPlain()
		}
		if !quoted && !validRFC2231ExtendedData(data) {
			if malformedRFC2231PercentEncodingOnly(data) {
				return fallbackPlain()
			}
			return "", "", true, false
		}
		raw, ok := decodeExtendedOctets(data)
		if !ok {
			return fallbackPlain()
		}
		decoded, ok := decodeMIMEParameterValue(charset, raw)
		if !ok {
			return fallbackPlain()
		}
		return decoded, language, true, true
	}
	if sectionCount == 0 {
		return "", "", true, false
	}

	var decoder extendedMediaParameterDecoder
	if sectionsOrdered {
		nextSection := 0
		orderedValid := true
		forEachMediaParameter(params, func(param string) {
			if !orderedValid {
				return
			}
			segment, ok := extendedContinuationSection(param, logicalName)
			if !ok {
				return
			}
			if segment.section != uint32(nextSection) || !decoder.consume(nextSection, segment) {
				orderedValid = false
				return
			}
			nextSection++
		})
		if !orderedValid || nextSection != sectionCount {
			return "", "", true, false
		}
		decoded, language, ok := decoder.finish()
		return decoded, language, true, ok
	}

	sections := make([]extendedMediaParameterSection, 0, sectionCount)
	forEachMediaParameter(params, func(param string) {
		if segment, ok := extendedContinuationSection(param, logicalName); ok {
			sections = append(sections, segment)
		}
	})
	if len(sections) != sectionCount {
		return "", "", true, false
	}
	sort.SliceStable(sections, func(i, j int) bool {
		return sections[i].section < sections[j].section
	})
	nextSection := 0
	for i := 0; i < len(sections); {
		segment := sections[i]
		if segment.section != uint32(nextSection) {
			return "", "", true, false
		}
		core, _, ok := mediaParameterValueCore(segment.value)
		if !ok {
			return "", "", true, false
		}
		j := i + 1
		for j < len(sections) && sections[j].section == segment.section {
			duplicateCore, _, duplicateOK := mediaParameterValueCore(sections[j].value)
			if !duplicateOK || sections[j].encoded != segment.encoded || duplicateCore != core {
				return "", "", true, false
			}
			j++
		}
		if !decoder.consume(nextSection, segment) {
			return "", "", true, false
		}
		nextSection++
		i = j
	}
	decoded, language, ok := decoder.finish()
	return decoded, language, true, ok
}

type extendedMediaParameterDecoder struct {
	decoded        strings.Builder
	encodedRun     []byte
	decodedRun     []byte
	charsetDecoder *textencoding.Decoder
	charset        string
	language       string
}

func (decoder *extendedMediaParameterDecoder) consume(section int, segment extendedMediaParameterSection) bool {
	core, quoted, ok := mediaParameterValueCore(segment.value)
	if !ok {
		return false
	}
	data := core
	if section == 0 && segment.encoded {
		decoder.charset, decoder.language, data, ok = splitExtendedInitialValue(core)
		if !ok {
			return false
		}
	}

	if segment.encoded {
		if !quoted && !validRFC2231ExtendedData(data) {
			return false
		}
		decoder.encodedRun, ok = appendDecodedExtendedOctets(decoder.encodedRun, data)
		return ok
	}

	if !quoted && !isMIMEToken(data) {
		return false
	}
	if !decoder.flushEncoded() {
		return false
	}
	decoder.decoded.WriteString(data)
	return true
}

func (decoder *extendedMediaParameterDecoder) flushEncoded() bool {
	if len(decoder.encodedRun) == 0 {
		return true
	}
	if decoder.charset == "" {
		decoder.decoded.Write(decoder.encodedRun)
		decoder.encodedRun = decoder.encodedRun[:0]
		return true
	}
	if decoder.charsetDecoder == nil {
		decoderEncoding, ok := mimeParameterEncoding(decoder.charset)
		if !ok || decoderEncoding == nil {
			return false
		}
		decoder.charsetDecoder = decoderEncoding.NewDecoder()
	}
	decoder.decodedRun = decoder.decodedRun[:0]
	var err error
	decoder.decodedRun, _, err = transform.Append(decoder.charsetDecoder, decoder.decodedRun, decoder.encodedRun)
	if err != nil {
		return false
	}
	decoder.decoded.Write(decoder.decodedRun)
	decoder.encodedRun = decoder.encodedRun[:0]
	return true
}

func (decoder *extendedMediaParameterDecoder) finish() (string, string, bool) {
	if !decoder.flushEncoded() {
		return "", "", false
	}
	return decoder.decoded.String(), decoder.language, true
}

func extendedContinuationSection(param string, logicalName string) (extendedMediaParameterSection, bool) {
	equals := strings.IndexByte(param, '=')
	if equals < 0 {
		return extendedMediaParameterSection{}, false
	}
	nameStart, nameEnd := trimOWSBounds(param[:equals])
	name := param[nameStart:nameEnd]
	if !strings.EqualFold(mediaParameterLogicalName(name), logicalName) {
		return extendedMediaParameterSection{}, false
	}
	encoded := strings.HasSuffix(name, "*")
	sectionName := strings.TrimSuffix(name, "*")
	if strings.EqualFold(sectionName, logicalName) {
		return extendedMediaParameterSection{}, false
	}
	star := strings.LastIndexByte(sectionName, '*')
	if star < 0 {
		return extendedMediaParameterSection{}, false
	}
	sectionValue, err := strconv.ParseUint(sectionName[star+1:], 10, 32)
	if err != nil {
		return extendedMediaParameterSection{}, false
	}
	return extendedMediaParameterSection{
		value:   param[equals+1:],
		section: uint32(sectionValue),
		encoded: encoded,
	}, true
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
	if !strings.Contains(contents, "\\") {
		return contents, true, true
	}

	var decoded strings.Builder
	for i := 0; i < len(contents); i++ {
		if contents[i] == '\\' && i+1 < len(contents) && isMIMETSpecial(contents[i+1]) {
			i++
		}
		decoded.WriteByte(contents[i])
	}
	return decoded.String(), true, true
}

func isMIMETSpecial(value byte) bool {
	switch value {
	case '(', ')', '<', '>', '@', ',', ';', ':', '\\', '"', '/', '[', ']', '?', '=':
		return true
	}
	return false
}

func parsePlainMediaParameterValue(value string) (string, bool) {
	_, params, err := mime.ParseMediaType("application/octet-stream; x=" + value)
	if err != nil {
		return "", false
	}
	parsed, ok := params["x"]
	return parsed, ok
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

func malformedRFC2231PercentEncodingOnly(value string) bool {
	malformedPercent := false
	for i := 0; i < len(value); i++ {
		if value[i] == '%' {
			if i+2 >= len(value) || !isHexDigit(value[i+1]) || !isHexDigit(value[i+2]) {
				malformedPercent = true
				continue
			}
			i += 2
			continue
		}
		if !isRFC2231AttributeChar(value[i]) {
			return false
		}
	}
	return malformedPercent
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
	firstUpper := -1
	for i := 0; i < len(value); i++ {
		if value[i] >= 'A' && value[i] <= 'Z' {
			firstUpper = i
			break
		}
	}
	if firstUpper < 0 {
		return value
	}
	bytes := []byte(value)
	for i := firstUpper; i < len(bytes); i++ {
		if bytes[i] >= 'A' && bytes[i] <= 'Z' {
			bytes[i] += 'a' - 'A'
		}
	}
	return string(bytes)
}

func asciiUpper(value string) string {
	bytes := []byte(value)
	for i := range bytes {
		if bytes[i] >= 0x80 {
			return ""
		}
		if bytes[i] >= 'a' && bytes[i] <= 'z' {
			bytes[i] -= 'a' - 'A'
		}
	}
	return string(bytes)
}

func validRFC822DateTime(value string) bool {
	// RFC 822 date-time values necessarily contain a time separator and
	// whitespace separating the date/time components. Reject obvious junk
	// before the comparatively allocation-heavy mail.ParseDate fallback.
	if strings.IndexByte(value, ':') < 0 || strings.IndexAny(value, " \t") < 0 {
		return false
	}
	normalized := asciiUpper(value)
	if normalized == "" {
		return false
	}
	parsed, err := mail.ParseDate(normalized)
	if err != nil {
		zoneStart := strings.LastIndexAny(normalized, " \t")
		if zoneStart < 0 || len(normalized)-zoneStart-1 != 1 {
			return false
		}
		zone := normalized[len(normalized)-1]
		if zone < 'A' || zone > 'Z' || zone == 'J' {
			return false
		}
		parsed, err = mail.ParseDate(normalized[:zoneStart+1] + "+0000")
		if err != nil {
			return false
		}
	}
	if comma := strings.IndexByte(normalized, ','); comma >= 0 {
		day := strings.TrimSpace(normalized[:comma])
		wantDay := asciiUpper(parsed.Weekday().String()[:3])
		if day != wantDay {
			return false
		}
	}
	return true
}

func validFixedLengthHex(value string, length int) bool {
	if len(value) != length {
		return false
	}
	for i := 0; i < len(value); i++ {
		if !isHexDigit(value[i]) {
			return false
		}
	}
	return true
}

func validH264ProfileLevelID(value string) bool {
	return validFixedLengthHex(value, 6)
}

func validH264MaxReceiveLevel(value string) bool {
	// RFC 6184 and RFC 6190 define max-recv-level as the base16
	// representation of profile-iop plus level_idc: two bytes, four digits.
	return validFixedLengthHex(value, 4)
}

func decodeExtendedOctets(value string) ([]byte, bool) {
	return appendDecodedExtendedOctets(make([]byte, 0, len(value)), value)
}

func appendDecodedExtendedOctets(decoded []byte, value string) ([]byte, bool) {
	start := len(decoded)
	for i := 0; i < len(value); i++ {
		if value[i] != '%' {
			decoded = append(decoded, value[i])
			continue
		}
		if i+2 >= len(value) || !isHexDigit(value[i+1]) || !isHexDigit(value[i+2]) {
			return decoded[:start], false
		}
		decoded = append(decoded, hexValue(value[i+1])<<4|hexValue(value[i+2]))
		i += 2
	}
	return decoded, true
}

func mimeParameterEncoding(charset string) (textencoding.Encoding, bool) {
	if charset == "" {
		return nil, true
	}

	decoderEncoding, err := ianaindex.IANA.Encoding(charset)
	if err != nil {
		return nil, false
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
			return nil, false
		}
	}
	return decoderEncoding, true
}

func decodeMIMEParameterValue(charset string, value []byte) (string, bool) {
	decoderEncoding, ok := mimeParameterEncoding(charset)
	if !ok {
		return "", false
	}
	if decoderEncoding == nil {
		return string(value), true
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
	value          string
	language       string
	found          bool
	decoded        bool
	duplicatePlain bool
}

func normalizeMediaParameterTail(mediaType string, params string, externalBodyAccessType string, externalBodyAccessLanguage string, hasExternalBodyAccessType bool) string {
	decodedParameters := map[string]decodedExtendedParameterState{}
	var normalized strings.Builder
	normalized.Grow(len(params))
	first := true
	forEachMediaParameter(params, func(param string) {
		if !first {
			normalized.WriteByte(';')
		}
		first = false
		writeNormalizedMediaParameter(&normalized, mediaType, params, param, externalBodyAccessType, externalBodyAccessLanguage, hasExternalBodyAccessType, decodedParameters)
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

func writeNormalizedMediaParameter(normalized *strings.Builder, mediaType string, params string, param string, externalBodyAccessType string, externalBodyAccessLanguage string, hasExternalBodyAccessType bool, decodedParameters map[string]decodedExtendedParameterState) {
	equals := strings.IndexByte(param, '=')
	if equals < 0 {
		normalized.WriteString(param)
		return
	}

	namePart := param[:equals]
	nameStart, nameEnd := trimOWSBounds(namePart)
	if nameStart == nameEnd {
		normalized.WriteString(param)
		return
	}
	name := namePart[nameStart:nameEnd]
	logicalName := mediaParameterLogicalName(name)

	normalized.WriteString(namePart[:nameStart])
	writeASCIILower(normalized, name)
	normalized.WriteString(namePart[nameEnd:])
	normalized.WriteByte('=')

	value := param[equals+1:]
	isExternalBodyAccessType := mediaType == "message/external-body" && strings.EqualFold(logicalName, "access-type")
	switch {
	case isExternalBodyAccessType:
		writeCanonicalMediaParameterValue(normalized, name, value, "access-type", externalBodyAccessType, externalBodyAccessLanguage, hasExternalBodyAccessType)
	case isCaseInsensitiveMediaParameterValue(mediaType, logicalName, externalBodyAccessType):
		state := decodedCaseInsensitiveParameter(mediaType, params, logicalName, externalBodyAccessType, decodedParameters)
		if state.duplicatePlain {
			normalized.WriteString(value)
		} else if state.found {
			writeCanonicalMediaParameterValue(normalized, name, value, logicalName, state.value, state.language, state.decoded)
		} else {
			core, _, parsed := mediaParameterValueCore(value)
			canonical, ok := normalizeCaseInsensitiveMediaParameterValue(mediaType, logicalName, core, externalBodyAccessType)
			if parsed && ok {
				valueStart, valueEnd := trimOWSBounds(value)
				normalized.WriteString(value[:valueStart])
				if isCanonicalMediaTypeParameterValue(mediaType, logicalName) {
					writeEncodedDecoderKeyValue(normalized, 'm', canonical)
				} else {
					writeASCIILower(normalized, value[valueStart:valueEnd])
				}
				normalized.WriteString(value[valueEnd:])
			} else {
				normalized.WriteString(value)
			}
		}
	case strings.HasSuffix(name, "*"):
		writeNormalizedExtendedParameterValue(normalized, value, extendedMediaParameterHasMetadata(name))
	default:
		normalized.WriteString(value)
	}

}

func writeASCIILower(dst *strings.Builder, value string) {
	start := 0
	for i := 0; i < len(value); i++ {
		if value[i] < 'A' || value[i] > 'Z' {
			continue
		}
		if start < i {
			dst.WriteString(value[start:i])
		}
		dst.WriteByte(value[i] + ('a' - 'A'))
		start = i + 1
	}
	if start < len(value) {
		dst.WriteString(value[start:])
	}
}

func decodedCaseInsensitiveParameter(mediaType string, params string, logicalName string, externalBodyAccessType string, cache map[string]decodedExtendedParameterState) decodedExtendedParameterState {
	key := asciiLower(logicalName)
	if state, ok := cache[key]; ok {
		return state
	}
	value, language, found, decoded := decodeExtendedMediaParameter(params, logicalName)
	if !found {
		plainValue, count, equal, valid := scanPlainMediaParameter(params, logicalName)
		if count > 1 {
			if equal && valid {
				if canonical, ok := normalizeCaseInsensitiveMediaParameterValue(mediaType, logicalName, plainValue, externalBodyAccessType); ok {
					state := decodedExtendedParameterState{value: canonical, found: true, decoded: true}
					cache[key] = state
					return state
				}
			}
			state := decodedExtendedParameterState{found: true, duplicatePlain: true}
			cache[key] = state
			return state
		}
	}
	if decoded {
		if canonical, ok := normalizeCaseInsensitiveMediaParameterValue(mediaType, logicalName, value, externalBodyAccessType); ok {
			value = canonical
			language = asciiLower(language)
		} else {
			decoded = false
		}
	}
	state := decodedExtendedParameterState{value: value, language: language, found: found, decoded: decoded}
	cache[key] = state
	return state
}

func scanPlainMediaParameter(params string, logicalName string) (string, int, bool, bool) {
	value := ""
	count := 0
	equal := true
	valid := true
	forEachMediaParameter(params, func(param string) {
		equals := strings.IndexByte(param, '=')
		if equals < 0 {
			return
		}
		nameStart, nameEnd := trimOWSBounds(param[:equals])
		if !strings.EqualFold(param[nameStart:nameEnd], logicalName) {
			return
		}
		core, _, ok := mediaParameterValueCore(param[equals+1:])
		if !ok {
			valid = false
			count++
			return
		}
		if count == 0 {
			value = core
		} else if core != value {
			equal = false
		}
		count++
	})
	return value, count, equal, valid
}

func writeCanonicalMediaParameterValue(dst *strings.Builder, name string, value string, logicalName string, decodedValue string, decodedLanguage string, decoded bool) {
	if decoded {
		valueStart, valueEnd := trimOWSBounds(value)
		dst.WriteString(value[:valueStart])
		if isInitialMediaParameterSegment(name, logicalName) {
			writeDecodedDecoderKeyValue(dst, decodedLanguage, decodedValue)
		} else {
			dst.WriteByte('d')
		}
		dst.WriteString(value[valueEnd:])
		return
	}

	normalizedValue := value
	if strings.HasSuffix(name, "*") {
		var normalized strings.Builder
		normalized.Grow(len(value))
		writeNormalizedExtendedParameterValue(&normalized, value, extendedMediaParameterHasMetadata(name))
		normalizedValue = normalized.String()
	}
	valueStart, valueEnd := trimOWSBounds(normalizedValue)
	dst.WriteString(normalizedValue[:valueStart])
	writeEncodedDecoderKeyValue(dst, 'r', normalizedValue[valueStart:valueEnd])
	dst.WriteString(normalizedValue[valueEnd:])
}

func writeDecodedDecoderKeyValue(dst *strings.Builder, language string, value string) {
	dst.WriteByte('d')
	writeHex(dst, language)
	dst.WriteByte('g')
	writeHex(dst, value)
}

func writeEncodedDecoderKeyValue(dst *strings.Builder, prefix byte, value string) {
	dst.WriteByte(prefix)
	writeHex(dst, value)
}

func writeHex(dst *strings.Builder, value string) {
	const digits = "0123456789abcdef"
	for i := 0; i < len(value); i++ {
		dst.WriteByte(digits[value[i]>>4])
		dst.WriteByte(digits[value[i]&0x0f])
	}
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

	switch mediaType {
	case "message/external-body":
		switch asciiLower(name) {
		case "access-type", "permission", "expiration":
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
	case "application/xop+xml":
		return strings.EqualFold(name, "type")
	case "text/plain":
		switch asciiLower(name) {
		case "format", "delsp":
			return true
		}
	case "text/directory":
		return strings.EqualFold(name, "profile")
	case "text/calendar":
		switch asciiLower(name) {
		case "method", "component":
			return true
		}
	case "video/h264":
		switch asciiLower(name) {
		case "profile-level-id", "max-recv-level":
			return true
		}
	case "video/h264-svc":
		switch asciiLower(name) {
		case "profile-level-id", "max-recv-level", "max-recv-base-level":
			return true
		}
	}
	return false
}

func normalizeCaseInsensitiveMediaParameterValue(mediaType string, name string, value string, externalBodyAccessType string) (string, bool) {
	if isCanonicalMediaTypeParameterValue(mediaType, name) {
		return canonicalMIMEMediaTypeValue(value)
	}
	if !validCaseInsensitiveMediaParameterValue(mediaType, name, value, externalBodyAccessType) {
		return "", false
	}
	return asciiLower(value), true
}

func isCanonicalMediaTypeParameterValue(mediaType string, name string) bool {
	return mediaType == "application/xop+xml" && strings.EqualFold(name, "type")
}

func canonicalMIMEMediaTypeValue(value string) (string, bool) {
	mediaType, params, err := mime.ParseMediaType(value)
	if err != nil {
		return "", false
	}
	_, rawParams, hasParams := strings.Cut(value, ";")

	externalBodyAccessType := ""
	if mediaType == "message/external-body" {
		if hasParams {
			if accessType, _, ok := parseExternalBodyAccessType(rawParams); ok {
				externalBodyAccessType = accessType
			}
		} else if accessType, ok := params["access-type"]; ok && isMIMEToken(accessType) {
			externalBodyAccessType = asciiLower(accessType)
		}
	}

	// mime.ParseMediaType understands only a narrow RFC 2231 charset set.
	// Re-decode extended values for the small set of parameters whose value
	// semantics are case-insensitive so supported SDK charsets (for example
	// ISO-8859-1) do not degrade to a lossy tail-only value. Each logical
	// name is decoded at most once, keeping this pass linear for duplicates.
	if hasParams && strings.Contains(rawParams, "*") {
		seen := map[string]struct{}{}
		extendedValid := true
		forEachMediaParameter(rawParams, func(param string) {
			if !extendedValid {
				return
			}
			equals := strings.IndexByte(param, '=')
			if equals < 0 {
				return
			}
			nameStart, nameEnd := trimOWSBounds(param[:equals])
			name := param[nameStart:nameEnd]
			if !strings.Contains(name, "*") {
				return
			}
			logicalName := mediaParameterLogicalName(name)
			if mediaType == "application/xop+xml" && strings.EqualFold(logicalName, "type") {
				return
			}
			key := asciiLower(logicalName)
			if _, ok := seen[key]; ok {
				return
			}
			if !isCaseInsensitiveMediaParameterValue(mediaType, logicalName, externalBodyAccessType) {
				return
			}
			seen[key] = struct{}{}
			decoded, _, found, ok := decodeExtendedMediaParameter(rawParams, logicalName)
			if !found {
				return
			}
			if !ok {
				extendedValid = false
				return
			}
			canonical, ok := normalizeCaseInsensitiveMediaParameterValue(mediaType, logicalName, decoded, externalBodyAccessType)
			if !ok {
				extendedValid = false
				return
			}
			params[key] = canonical
		})
		if !extendedValid {
			return "", false
		}
	}

	for name, paramValue := range params {
		// Avoid recursively canonicalizing nested XOP type parameters here.
		// The outer XOP type is normalized once, while its case-sensitive
		// nested parameter values remain intact.
		if mediaType == "application/xop+xml" && strings.EqualFold(name, "type") {
			continue
		}
		if !isCaseInsensitiveMediaParameterValue(mediaType, name, externalBodyAccessType) {
			continue
		}
		canonical, ok := normalizeCaseInsensitiveMediaParameterValue(mediaType, name, paramValue, externalBodyAccessType)
		if !ok {
			return "", false
		}
		params[name] = canonical
	}

	formatted := mime.FormatMediaType(mediaType, params)
	if formatted == "" {
		return "", false
	}
	if !hasParams || !strings.Contains(rawParams, "*") {
		return formatted, true
	}
	extendedIdentity, ok := canonicalRFC2231MediaParameterIdentity(rawParams, len(params), mediaType, externalBodyAccessType)
	if !ok {
		return "", false
	}
	if extendedIdentity == "" {
		return formatted, true
	}
	return formatted + "\x00r" + extendedIdentity, true
}

type rfc2231MediaParameterIdentity struct {
	name          string
	logicalName   string
	language      string
	fallbackIndex uint32
	section       uint32
	continuation  bool
}

type rfc2231MediaParameterFallback struct {
	charset string
	data    string
}

func canonicalRFC2231MediaParameterIdentity(params string, estimatedIdentities int, mediaType string, externalBodyAccessType string) (string, bool) {
	identities := make([]rfc2231MediaParameterIdentity, 0, estimatedIdentities)
	var fallbacks []rfc2231MediaParameterFallback
	valid := true
	encodedLength := 0
	forEachMediaParameter(params, func(param string) {
		if !valid {
			return
		}
		equals := strings.IndexByte(param, '=')
		if equals < 0 {
			return
		}
		nameStart, nameEnd := trimOWSBounds(param[:equals])
		if nameStart == nameEnd {
			return
		}
		name := param[nameStart:nameEnd]
		logicalName := mediaParameterLogicalName(name)
		if !strings.Contains(name, "*") {
			return
		}

		encoded := strings.HasSuffix(name, "*")
		hasMetadata := extendedMediaParameterHasMetadata(name)
		var core string
		var quoted bool
		if encoded {
			var ok bool
			core, quoted, ok = mediaParameterValueCore(param[equals+1:])
			if !ok {
				valid = false
				return
			}
			// Go's mime.ParseMediaType can keep a successfully decoded prefix
			// while dropping a malformed later RFC 2231 continuation. Do not let
			// that lossy parse collapse malformed response metadata onto a valid
			// registered decoder key.
			if !hasMetadata && (!validPercentHexEscapes(core) || (!quoted && !validRFC2231ExtendedData(core))) {
				valid = false
				return
			}
		}

		identity := rfc2231MediaParameterIdentity{name: name, logicalName: logicalName}
		if segment, ok := extendedContinuationSection(param, logicalName); ok {
			identity.section = segment.section
			identity.continuation = true
		}
		if !hasMetadata && !identity.continuation {
			// ParseMediaType silently drops malformed RFC 2231 star forms.
			// Keep them out of canonical keys rather than losing their payload.
			valid = false
			return
		}
		if hasMetadata {
			charset, language, data, ok := splitRFC2231IdentityMetadata(core)
			if !ok {
				valid = false
				return
			}
			identity.language = language
			lowerCharset := asciiLower(charset)
			caseInsensitive := isCaseInsensitiveMediaParameterValue(mediaType, logicalName, externalBodyAccessType)
			needsFallback := !caseInsensitive && ((lowerCharset != "utf-8" && lowerCharset != "us-ascii") || !validPercentHexEscapes(data))
			if needsFallback {
				fallbacks = append(fallbacks, rfc2231MediaParameterFallback{charset: charset, data: data})
				identity.fallbackIndex = uint32(len(fallbacks))
				encodedLength += len(charset)*2 + len(data)*2 + 2
			}
		}

		identities = append(identities, identity)
		encodedLength += len(name)*2 + len(identity.language)*2 + 2
	})
	if !valid {
		return "", false
	}
	if len(identities) == 0 {
		return "", true
	}
	sort.Slice(identities, func(i, j int) bool {
		if comparison := compareASCIIFold(identities[i].logicalName, identities[j].logicalName); comparison != 0 {
			return comparison < 0
		}
		if identities[i].continuation != identities[j].continuation {
			return !identities[i].continuation
		}
		if identities[i].continuation && identities[i].section != identities[j].section {
			return identities[i].section < identities[j].section
		}
		if comparison := compareASCIIFold(identities[i].name, identities[j].name); comparison != 0 {
			return comparison < 0
		}
		return compareASCIIFold(identities[i].language, identities[j].language) < 0
	})

	var currentLogicalName string
	var previousContinuationName string
	var previousSection uint32
	haveLogicalName := false
	haveSingleExtended := false
	haveContinuation := false
	for _, identity := range identities {
		if !haveLogicalName || !strings.EqualFold(identity.logicalName, currentLogicalName) {
			currentLogicalName = identity.logicalName
			previousContinuationName = ""
			haveLogicalName = true
			haveSingleExtended = false
			haveContinuation = false
		}

		if !identity.continuation {
			haveSingleExtended = true
			continue
		}
		// ParseMediaType prefers a single extended value over continuations.
		// Reject that ambiguous mix so ignored payload cannot collapse keys.
		if haveSingleExtended {
			return "", false
		}
		if !haveContinuation {
			if identity.section != 0 {
				return "", false
			}
			previousContinuationName = identity.name
			previousSection = 0
			haveContinuation = true
			continue
		}
		if identity.section == previousSection {
			// ParseMediaType prefers the unencoded spelling when both forms
			// exist for one section; reject that lossy representation too.
			if !strings.EqualFold(identity.name, previousContinuationName) {
				return "", false
			}
			continue
		}
		if previousSection == ^uint32(0) || identity.section != previousSection+1 {
			return "", false
		}
		previousContinuationName = identity.name
		previousSection = identity.section
	}

	var encoded strings.Builder
	encoded.Grow(encodedLength)
	for i, identity := range identities {
		if i != 0 {
			encoded.WriteByte('h')
		}
		writeASCIILowerHex(&encoded, identity.name)
		encoded.WriteByte('g')
		writeASCIILowerHex(&encoded, identity.language)
		if identity.fallbackIndex != 0 {
			fallback := fallbacks[identity.fallbackIndex-1]
			encoded.WriteByte('v')
			writeASCIILowerHex(&encoded, fallback.charset)
			encoded.WriteByte('q')
			writeNormalizedPercentEncodingHex(&encoded, fallback.data)
		}
	}
	return encoded.String(), true
}

func splitRFC2231IdentityMetadata(value string) (string, string, string, bool) {
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
	if charset != strings.TrimSpace(charset) {
		return "", "", "", false
	}
	if language != "" && !isRFC1766LanguageTag(language) {
		return "", "", "", false
	}
	return charset, language, value[secondQuote+1:], true
}

func validPercentHexEscapes(value string) bool {
	for i := 0; i < len(value); i++ {
		if value[i] != '%' {
			continue
		}
		if i+2 >= len(value) || !isHexDigit(value[i+1]) || !isHexDigit(value[i+2]) {
			return false
		}
		i += 2
	}
	return true
}

func writeNormalizedPercentEncodingHex(dst *strings.Builder, value string) {
	const digits = "0123456789abcdef"
	writeByteHex := func(current byte) {
		dst.WriteByte(digits[current>>4])
		dst.WriteByte(digits[current&0x0f])
	}
	for i := 0; i < len(value); i++ {
		writeByteHex(value[i])
		if value[i] == '%' && i+2 < len(value) && isHexDigit(value[i+1]) && isHexDigit(value[i+2]) {
			writeByteHex(lowerHexDigit(value[i+1]))
			writeByteHex(lowerHexDigit(value[i+2]))
			i += 2
		}
	}
}

func compareASCIIFold(first string, second string) int {
	limit := len(first)
	if len(second) < limit {
		limit = len(second)
	}
	for i := 0; i < limit; i++ {
		firstByte := asciiLowerByte(first[i])
		secondByte := asciiLowerByte(second[i])
		if firstByte < secondByte {
			return -1
		}
		if firstByte > secondByte {
			return 1
		}
	}
	if len(first) < len(second) {
		return -1
	}
	if len(first) > len(second) {
		return 1
	}
	return 0
}

func asciiLowerByte(value byte) byte {
	if value >= 'A' && value <= 'Z' {
		return value + ('a' - 'A')
	}
	return value
}

func writeASCIILowerHex(dst *strings.Builder, value string) {
	const digits = "0123456789abcdef"
	for i := 0; i < len(value); i++ {
		current := asciiLowerByte(value[i])
		dst.WriteByte(digits[current>>4])
		dst.WriteByte(digits[current&0x0f])
	}
}

func validCaseInsensitiveMediaParameterValue(mediaType string, name string, value string, externalBodyAccessType string) bool {
	if strings.EqualFold(name, "charset") {
		// RFC 2046 defines charset values as case-insensitive and permits
		// private X-* names; registry lookup is only required when decoding
		// RFC 2231 octets, not when comparing the charset parameter value.
		return isMIMEToken(value)
	}

	switch mediaType {
	case "message/external-body":
		switch asciiLower(name) {
		case "access-type":
			return isMIMEToken(value)
		case "permission":
			switch asciiLower(value) {
			case "read", "read-write":
				return true
			}
		case "mode":
			return validExternalBodyMode(externalBodyAccessType, value)
		case "expiration":
			return validRFC822DateTime(value)
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
		case "format":
			switch asciiLower(value) {
			case "fixed", "flowed":
				return true
			}
		case "delsp":
			switch asciiLower(value) {
			case "yes", "no":
				return true
			}
		}
	case "text/directory":
		if strings.EqualFold(name, "profile") {
			return isMIMEToken(value)
		}
	case "text/calendar":
		switch asciiLower(name) {
		case "method", "component":
			return isMIMEToken(value)
		}
	case "video/h264":
		switch asciiLower(name) {
		case "profile-level-id":
			return validH264ProfileLevelID(value)
		case "max-recv-level":
			return validH264MaxReceiveLevel(value)
		}
	case "video/h264-svc":
		switch asciiLower(name) {
		case "profile-level-id":
			return validH264ProfileLevelID(value)
		case "max-recv-level", "max-recv-base-level":
			return validH264MaxReceiveLevel(value)
		}
	}
	return false
}

func validExternalBodyMode(accessType string, value string) bool {
	mode := asciiLower(value)
	switch accessType {
	case "ftp", "anon-ftp":
		switch mode {
		case "ascii", "ebcdic", "image":
			return true
		}
		if !strings.HasPrefix(mode, "local") || len(mode) == len("local") {
			return false
		}
		for i := len("local"); i < len(mode); i++ {
			if mode[i] < '0' || mode[i] > '9' {
				return false
			}
		}
		return true
	case "tftp":
		switch mode {
		case "netascii", "octet", "mail":
			return true
		}
	}
	return false
}

func writeNormalizedExtendedParameterValue(dst *strings.Builder, value string, hasMetadata bool) {
	valueStart, valueEnd := trimOWSBounds(value)
	core := value[valueStart:valueEnd]
	quoted := false
	if strings.HasPrefix(core, "\"") {
		var ok bool
		core, ok = quotedMediaParameterContents(core)
		if !ok {
			dst.WriteString(value)
			return
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

	dst.WriteString(value[:valueStart])
	if quoted {
		dst.WriteByte('"')
	}
	if firstQuote >= 0 && secondQuote >= 0 {
		writeASCIILower(dst, core[:firstQuote])
		dst.WriteByte('\'')
		writeASCIILower(dst, core[firstQuote+1:secondQuote])
		dst.WriteByte('\'')
		writeNormalizedPercentEncoding(dst, core[secondQuote+1:])
	} else {
		writeNormalizedPercentEncoding(dst, core)
	}
	if quoted {
		dst.WriteByte('"')
	}
	dst.WriteString(value[valueEnd:])
}

func writeNormalizedPercentEncoding(dst *strings.Builder, value string) {
	start := 0
	for i := 0; i+2 < len(value); i++ {
		if value[i] != '%' || !isHexDigit(value[i+1]) || !isHexDigit(value[i+2]) {
			continue
		}
		dst.WriteString(value[start : i+1])
		dst.WriteByte(lowerHexDigit(value[i+1]))
		dst.WriteByte(lowerHexDigit(value[i+2]))
		i += 2
		start = i + 1
	}
	dst.WriteString(value[start:])
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

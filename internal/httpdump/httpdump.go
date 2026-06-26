package httpdump

import (
	"net/http"
	"strings"
)

var sensitiveHeaders = map[string]struct{}{
	"authorization":       {},
	"proxy-authorization": {},
	"api-key":             {},
	"x-api-key":           {},
	"cookie":              {},
	"set-cookie":          {},
	"openai-organization": {},
	"openai-project":      {},
	"webhook-id":          {},
	"webhook-signature":   {},
	"webhook-timestamp":   {},
}

// RedactSensitiveHeaders returns a header map with known sensitive values
// replaced. If no sensitive headers are present, the original map is returned.
func RedactSensitiveHeaders(headers http.Header) http.Header {
	var redacted http.Header
	for name, values := range headers {
		if len(values) == 0 {
			continue
		}
		if _, ok := sensitiveHeaders[strings.ToLower(name)]; !ok {
			continue
		}
		if redacted == nil {
			redacted = headers.Clone()
		}
		redacted[name] = make([]string, len(values))
		for i := range values {
			redacted[name][i] = "***"
		}
	}
	if redacted == nil {
		return headers
	}
	return redacted
}

package option

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type debugRequestMetadata struct {
	Method        string      `json:"method"`
	URL           string      `json:"url"`
	ContentLength int64       `json:"content_length"`
	Headers       http.Header `json:"headers"`
}

type debugResponseMetadata struct {
	Status        string      `json:"status"`
	StatusCode    int         `json:"status_code"`
	ContentLength int64       `json:"content_length"`
	Headers       http.Header `json:"headers"`
}

// WithDebugLog logs allowlisted HTTP request and response metadata.
// If the logger parameter is nil, it uses the default logger.
//
// Request and response bodies, URL user information and query strings, and
// headers outside the allowlist are omitted. Credential-bearing headers are
// represented by a placeholder for each original value.
//
// WithDebugLog is for debugging and development purposes only.
// It should not be used in production code. The behavior and interface
// of WithDebugLog is not guaranteed to be stable.
func WithDebugLog(logger *log.Logger) RequestOption {
	if logger == nil {
		logger = log.Default()
	}

	return WithMiddleware(func(req *http.Request, nxt MiddlewareNext) (*http.Response, error) {
		requestMetadata, marshalErr := json.Marshal(debugRequestMetadata{
			Method:        req.Method,
			URL:           debugLogURL(req.URL),
			ContentLength: req.ContentLength,
			Headers:       debugLogHeaders(req.Header),
		})
		if marshalErr == nil {
			logger.Printf("Request Metadata: %s\n", requestMetadata)
		}

		resp, err := nxt(req)
		if err != nil {
			return resp, err
		}

		if resp != nil {
			responseMetadata, responseMarshalErr := json.Marshal(debugResponseMetadata{
				Status:        resp.Status,
				StatusCode:    resp.StatusCode,
				ContentLength: resp.ContentLength,
				Headers:       debugLogHeaders(resp.Header),
			})
			if responseMarshalErr == nil {
				logger.Printf("Response Metadata: %s\n", responseMetadata)
			}
		}

		return resp, err
	})
}

func debugLogURL(original *url.URL) string {
	if original == nil {
		return ""
	}
	return (&url.URL{
		Scheme:  original.Scheme,
		Host:    original.Host,
		Path:    original.Path,
		RawPath: original.RawPath,
	}).String()
}

func debugLogHeaders(headers http.Header) http.Header {
	result := make(http.Header)
	for name, values := range headers {
		canonicalName := http.CanonicalHeaderKey(name)
		switch strings.ToLower(name) {
		case "authorization", "proxy-authorization", "api-key", "x-api-key", "x-amz-security-token", "cookie", "set-cookie":
			for range values {
				result.Add(canonicalName, "***")
			}
		case "accept", "content-length", "content-type", "openai-processing-ms", "openai-version", "request-id", "retry-after", "user-agent", "x-request-id", "x-stainless-retry-count", "x-stainless-timeout":
			for _, value := range values {
				result.Add(canonicalName, value)
			}
		}
	}
	return result
}

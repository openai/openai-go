package option

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

const debugLogRedacted = "***"

type debugRequestMetadata struct {
	Method        string      `json:"method"`
	ContentLength int64       `json:"content_length"`
	Headers       http.Header `json:"headers"`
}

type debugResponseMetadata struct {
	StatusCode    int         `json:"status_code"`
	ContentLength int64       `json:"content_length"`
	Headers       http.Header `json:"headers"`
}

// WithDebugLog logs allowlisted HTTP request and response metadata.
// If the logger parameter is nil, it uses the default logger.
//
// Request and response bodies, URLs, free-form response status text, and all
// original header values are omitted. Unrecognized request methods and values
// for recognized credential-bearing header names are represented by a
// placeholder.
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
			Method:        debugLogMethod(req.Method),
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

func debugLogMethod(method string) string {
	switch method {
	case http.MethodConnect,
		http.MethodDelete,
		http.MethodGet,
		http.MethodHead,
		http.MethodOptions,
		http.MethodPatch,
		http.MethodPost,
		http.MethodPut,
		http.MethodTrace:
		return method
	default:
		return debugLogRedacted
	}
}

func debugLogHeaders(headers http.Header) http.Header {
	result := make(http.Header)
	for name, values := range headers {
		var redactedName string
		switch strings.ToLower(name) {
		case "authorization":
			redactedName = "Authorization"
		case "proxy-authorization":
			redactedName = "Proxy-Authorization"
		case "api-key":
			redactedName = "Api-Key"
		case "x-api-key":
			redactedName = "X-Api-Key"
		case "x-amz-security-token":
			redactedName = "X-Amz-Security-Token"
		case "cookie":
			redactedName = "Cookie"
		case "set-cookie":
			redactedName = "Set-Cookie"
		default:
			continue
		}
		for range values {
			result.Add(redactedName, debugLogRedacted)
		}
	}
	return result
}

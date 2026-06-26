// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package option

import (
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/openai/openai-go/v3/internal/httpdump"
)

// WithDebugLog logs the HTTP request and response content.
// If the logger parameter is nil, it uses the default logger.
//
// WithDebugLog is for debugging and development purposes only.
// It should not be used in production code. The behavior and interface
// of WithDebugLog is not guaranteed to be stable.
func WithDebugLog(logger *log.Logger) RequestOption {
	return WithMiddleware(func(req *http.Request, nxt MiddlewareNext) (*http.Response, error) {
		if logger == nil {
			logger = log.Default()
		}

		if reqBytes, err := dumpRedactedRequest(req); err == nil {
			logger.Printf("Request Content:\n%s\n", reqBytes)
		}

		resp, err := nxt(req)
		if err != nil {
			return resp, err
		}

		if respBytes, err := dumpRedactedResponse(resp); err == nil {
			logger.Printf("Response Content:\n%s\n", respBytes)
		}

		return resp, err
	})
}

// dumpRedactedRequest dumps req with sensitive headers replaced. The
// original headers are restored via defer so a panic in DumpRequest cannot
// leak the placeholder map into the live request sent downstream.
func dumpRedactedRequest(req *http.Request) ([]byte, error) {
	origHeaders := req.Header
	req.Header = httpdump.RedactSensitiveHeaders(origHeaders)
	defer func() { req.Header = origHeaders }()
	return httputil.DumpRequest(req, true)
}

func dumpRedactedResponse(resp *http.Response) ([]byte, error) {
	origHeaders := resp.Header
	resp.Header = httpdump.RedactSensitiveHeaders(origHeaders)
	defer func() { resp.Header = origHeaders }()
	return httputil.DumpResponse(resp, true)
}

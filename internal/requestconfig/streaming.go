package requestconfig

import (
	"context"
	"net/http"
)

type sseMaxEventBytesKey struct{}

// WithSSEMaxEventBytes attaches an SSE event budget to req. This function is
// internal API and may change without notice.
func WithSSEMaxEventBytes(req *http.Request, maxBytes int) *http.Request {
	return req.WithContext(context.WithValue(req.Context(), sseMaxEventBytesKey{}, maxBytes))
}

// SSEMaxEventBytes returns the SSE event budget attached to req. This function
// is internal API and may change without notice.
func SSEMaxEventBytes(req *http.Request) int {
	if req == nil {
		return 0
	}
	maxBytes, _ := req.Context().Value(sseMaxEventBytesKey{}).(int)
	return maxBytes
}

func attachSSEMaxEventBytes(res *http.Response, fallbackRequest *http.Request, maxBytes int) {
	if res == nil || maxBytes <= 0 {
		return
	}
	request := res.Request
	if request == nil {
		request = fallbackRequest
	}
	res.Request = WithSSEMaxEventBytes(request, maxBytes)
}

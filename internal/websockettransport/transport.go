// Package websockettransport preserves HTTP cancellation during upgrades.
package websockettransport

import (
	"context"
	"io"
	"net/http"
	"sync"
)

// RoundTrip honors legacy transport cancellation without an HTTP client timer
// that would keep running after the response body becomes a WebSocket.
func RoundTrip(transport http.RoundTripper, request *http.Request) (*http.Response, error) {
	if err := request.Context().Err(); err != nil {
		return nil, err
	}
	canceler, ok := transport.(interface{ CancelRequest(*http.Request) })
	if !ok {
		return transport.RoundTrip(request)
	}
	done := make(chan struct{})
	stop := context.AfterFunc(request.Context(), func() {
		canceler.CancelRequest(request)
		close(done)
	})
	var once sync.Once
	stopCancel := func() {
		once.Do(func() {
			if !stop() {
				<-done
			}
		})
	}
	response, err := transport.RoundTrip(request)
	if err != nil || response == nil || response.Body == nil || response.StatusCode == http.StatusSwitchingProtocols {
		stopCancel()
	} else {
		// Redirect/rejection bodies are still part of opening the connection.
		// Keep legacy cancellation active until the HTTP client closes them.
		response.Body = &cancelBody{ReadCloser: response.Body, stop: stopCancel}
	}
	return response, err
}

type cancelBody struct {
	io.ReadCloser
	stop func()
}

func (b *cancelBody) Close() error {
	defer b.stop()
	return b.ReadCloser.Close()
}

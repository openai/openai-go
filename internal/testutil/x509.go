package testutil

import (
	"context"
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

// NewNativeX509HTTPClient routes HTTPS requests through a hermetic TLS server
// while retaining a real native transport with one static client identity.
func NewNativeX509HTTPClient(
	t testing.TB,
	roundTrip func(*http.Request) (*http.Response, error),
) *http.Client {
	t.Helper()
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(writer http.ResponseWriter, incoming *http.Request) {
		request := incoming.Clone(incoming.Context())
		request.URL.Scheme = "https"
		request.URL.Host = incoming.Host
		response, err := roundTrip(request)
		if err != nil {
			connection, _, hijackErr := writer.(http.Hijacker).Hijack()
			if hijackErr == nil {
				_ = connection.Close()
			}
			return
		}
		if response == nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		if response.Body != nil {
			defer func() { _ = response.Body.Close() }()
		}
		for name, values := range response.Header {
			for _, value := range values {
				writer.Header().Add(name, value)
			}
		}
		statusCode := response.StatusCode
		if statusCode == 0 {
			statusCode = http.StatusOK
		}
		writer.WriteHeader(statusCode)
		if response.Body != nil {
			_, _ = io.Copy(writer, response.Body)
		}
	}))
	server.StartTLS()
	t.Cleanup(server.Close)

	serverTransport := server.Client().Transport.(*http.Transport)
	transport := serverTransport.Clone()
	transport.TLSClientConfig = serverTransport.TLSClientConfig.Clone()
	transport.TLSClientConfig.ServerName = "example.com"
	transport.TLSClientConfig.Certificates = []tls.Certificate{{Certificate: [][]byte{{1}}}}
	transport.TLSClientConfig.ClientSessionCache = nil
	dialer := &net.Dialer{}
	transport.DialContext = func(ctx context.Context, _, _ string) (net.Conn, error) {
		return dialer.DialContext(ctx, "tcp", server.Listener.Addr().String())
	}
	return &http.Client{Transport: transport}
}

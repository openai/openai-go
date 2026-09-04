package openai

import (
	"net/http"
	"time"

	"github.com/openai/openai-go/v3/internal/requestconfig"
)

// defaultResponseHeaderTimeout bounds the time between a fully written request
// and the server's response headers. It does not apply to the response body,
// so long-running streams are unaffected. Without this, a server that accepts
// the connection but never responds would hang the request indefinitely.
const defaultResponseHeaderTimeout = 10 * time.Minute

// defaultHTTPClient returns an internally marked [*http.Client] used when the
// caller does not supply one via [option.WithHTTPClient]. When
// [http.DefaultTransport] is the stdlib [*http.Transport], it is cloned and a
// [http.Transport.ResponseHeaderTimeout] is set so stuck connections fail fast
// instead of compounding across retries.
// If [http.DefaultTransport] has been wrapped (for example by otelhttp for
// distributed tracing), the wrapping is preserved and the header timeout is
// skipped.
func defaultHTTPClient() *requestconfig.DefaultHTTPClient {
	if t, ok := http.DefaultTransport.(*http.Transport); ok {
		t = t.Clone()
		t.ResponseHeaderTimeout = defaultResponseHeaderTimeout
		return &requestconfig.DefaultHTTPClient{Client: &http.Client{Transport: t}}
	}
	return &requestconfig.DefaultHTTPClient{Client: &http.Client{Transport: http.DefaultTransport}}
}

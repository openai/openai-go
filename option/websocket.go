package option

import (
	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/packages/websocket"
	"net/http"
)

// WebSocketHTTPClient is an optional capability for a custom HTTPClient used
// with Responses WebSockets. Its returned client must preserve the custom
// client's authentication, proxy, TLS and transport policies and support
// upgraded connections with writable response bodies. The SDK does not close
// the supplied client or its transport.
//
// Pass implementations to WithHTTPClient as usual. A standard *http.Client
// already supports the capability and needs no wrapper. Bespoke HTTPClient
// implementations without this capability are rejected before an upgrade is
// attempted; they continue to work for ordinary HTTP and SSE requests.
type WebSocketHTTPClient interface {
	HTTPClient
	WebSocketHTTPClient() *http.Client
}

// WithWebSocketReconnect enables optional transport recovery for bindings that
// support it. The default is false. Recovery never restores server state or
// replays a command whose delivery is uncertain. Other request types ignore it.
func WithWebSocketReconnect(enabled bool) RequestOption {
	return requestconfig.RequestOptionFunc(func(cfg *requestconfig.RequestConfig) error {
		cfg.WebSocketRecovery.Enabled = enabled
		return nil
	})
}

// WithWebSocketRecovery configures supported bindings' recovery and pending-write
// limits. Set Enabled to opt in. Later request options override earlier ones.
func WithWebSocketRecovery(options websocket.RecoveryOptions) RequestOption {
	return requestconfig.RequestOptionFunc(func(cfg *requestconfig.RequestConfig) error {
		cfg.WebSocketRecovery = options
		return nil
	})
}

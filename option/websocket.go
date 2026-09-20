package option

import "net/http"

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

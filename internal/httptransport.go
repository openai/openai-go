package internal

import "net/http"

// NativeHTTPTransport returns the effective RoundTripper for a native
// *http.Client. Opaque HTTP doers do not expose a transport identity.
func NativeHTTPTransport(httpDoer any) http.RoundTripper {
	client, ok := httpDoer.(*http.Client)
	if !ok {
		return nil
	}
	if client.Transport == nil {
		return http.DefaultTransport
	}
	return client.Transport
}

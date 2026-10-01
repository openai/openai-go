package requestconfig

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type unsafeLoopbackContextKey struct{}

// WithUnsafeAllowHTTP permits authenticated local development over a dedicated
// direct loopback transport. Its connection pool is shared by this option's users.
func WithUnsafeAllowHTTP() RequestOption {
	transport := &http.Transport{
		DialContext:           dialCredentialLoopback,
		IdleConnTimeout:       90 * time.Second,
		ResponseHeaderTimeout: 10 * time.Minute,
		ExpectContinueTimeout: time.Second,
	}
	return RequestOptionFunc(func(cfg *RequestConfig) error {
		cfg.unsafeLoopbackTransport = transport
		return nil
	})
}

func credentialTransportError() error {
	return WithNoRetryError(errors.New("openai: authenticated requests require HTTPS; option.WithUnsafeAllowHTTP permits HTTP only for direct loopback development endpoints"))
}

func credentialLoopbackURL(u *url.URL) bool {
	if u == nil || u.Scheme != "http" || u.Opaque != "" {
		return false
	}
	if strings.EqualFold(u.Hostname(), "localhost") {
		return true
	}
	ip := net.ParseIP(u.Hostname())
	return ip != nil && ip.IsLoopback()
}

func validateCredentialURL(u *url.URL, allowLoopback bool) error {
	if u != nil && u.Scheme == "https" && u.Hostname() != "" && u.Opaque == "" {
		return nil
	}
	if allowLoopback && credentialLoopbackURL(u) {
		return nil
	}
	return credentialTransportError()
}

// ValidateOpenAICredentialRequest checks a workload destination before acquiring
// its token. The SDK installs the loopback permission only with a direct transport.
// Rejected requests release their body using the SDK's close-once ownership.
func ValidateOpenAICredentialRequest(req *http.Request) error {
	if !requestHasCanonicalTarget(req) {
		closeRequestBody(req)
		return credentialTransportError()
	}
	allowLoopback, _ := req.Context().Value(unsafeLoopbackContextKey{}).(bool)
	if err := validateCredentialURL(req.URL, allowLoopback); err != nil {
		closeRequestBody(req)
		return err
	}
	return nil
}

func (cfg *RequestConfig) configureCredentialTransport() error {
	// Azure and Bedrock own their credential and transport policies.
	if cfg.endpointProvider != "" {
		return nil
	}
	// Workload acquisition must be denied before middleware can request a token.
	// Static credentials are checked at dispatch, after middleware can remove them.
	_, cfg.requireSecureTransport = cfg.ProviderAuth("OpenAI")
	if !cfg.requireSecureTransport {
		return nil
	}
	base := cfg.BaseURL
	if base == nil {
		base = cfg.DefaultBaseURL
	}
	if base == nil {
		// Keep the existing missing-endpoint error from Execute/PrepareWebSocket.
		return nil
	}
	if base.Scheme == "ws" || base.Scheme == "wss" {
		copy := *base
		if copy.Scheme == "ws" {
			copy.Scheme = "http"
		} else {
			copy.Scheme = "https"
		}
		base = &copy
	}
	return validateCredentialURL(base, cfg.unsafeLoopbackTransport != nil)
}

func dialCredentialLoopback(ctx context.Context, network, address string) (net.Conn, error) {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return nil, credentialTransportError()
	}
	dialer := net.Dialer{Timeout: 30 * time.Second, KeepAlive: 30 * time.Second}
	if !strings.EqualFold(host, "localhost") {
		ip := net.ParseIP(host)
		if ip == nil || !ip.IsLoopback() {
			return nil, credentialTransportError()
		}
		return dialer.DialContext(ctx, network, net.JoinHostPort(ip.String(), port))
	}
	// Never resolve localhost through DNS or a caller-supplied resolver.
	conn, err := dialer.DialContext(ctx, network, net.JoinHostPort("127.0.0.1", port))
	if err == nil {
		return conn, nil
	}
	return dialer.DialContext(ctx, network, net.JoinHostPort("::1", port))
}

func hasAuthorization(headers http.Header) bool {
	for name, values := range headers {
		if strings.EqualFold(name, "Authorization") {
			for _, value := range values {
				if value != "" {
					return true
				}
			}
		}
	}
	return false
}

// credentialTransport validates the final middleware output and selects the
// dedicated direct transport when a credential-bearing request uses local HTTP.
// A nil transport preserves the caller's configured client.
func (cfg *RequestConfig) credentialTransport(req *http.Request) (http.RoundTripper, error) {
	if cfg.endpointProvider != "" || (!cfg.requireSecureTransport && !hasAuthorization(req.Header) && (req.URL == nil || req.URL.User == nil)) {
		return nil, nil
	}
	if !requestHasCanonicalTarget(req) {
		closeRequestBody(req)
		return nil, credentialTransportError()
	}
	if err := validateCredentialURL(req.URL, cfg.unsafeLoopbackTransport != nil); err != nil {
		closeRequestBody(req)
		return nil, err
	}
	if req.URL.Scheme == "http" {
		return cfg.unsafeLoopbackTransport, nil
	}
	return nil, nil
}

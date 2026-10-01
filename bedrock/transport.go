package bedrock

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"
)

func validateEndpointTransport(endpoint *url.URL, unsafeAllowHTTP bool) error {
	if endpoint.Scheme == "https" {
		return nil
	}
	if endpoint.Scheme == "http" && unsafeAllowHTTP && isLoopbackHost(endpoint.Hostname()) {
		return nil
	}
	return errors.New("bedrock: authenticated endpoints require HTTPS; UnsafeAllowHTTP permits HTTP only for local development on loopback endpoints")
}

func isLoopbackHost(host string) bool {
	if host == "localhost" {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

// Finalizers run per request. Reuse direct connection pools, but bound the cache
// so temporary method-level transports do not live as long as the client.
type loopbackTransportCache struct {
	mu         sync.Mutex
	transports [8]struct {
		source    *http.Transport // nil means the default
		transport *http.Transport
	}
	next int
}

func (c *loopbackTransportCache) get(base http.RoundTripper) (*http.Transport, error) {
	var source *http.Transport
	if base != nil {
		var ok bool
		source, ok = base.(*http.Transport)
		if !ok || source == nil {
			return nil, errors.New("bedrock: authenticated loopback HTTP requires a nil Transport or an *http.Transport")
		}
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, cached := range c.transports {
		if cached.transport != nil && cached.source == source {
			return cached.transport, nil
		}
	}
	transport := &http.Transport{
		ResponseHeaderTimeout: 10 * time.Minute,
		ExpectContinueTimeout: time.Second,
	}
	if source != nil {
		transport = source.Clone()
	}
	transport.Proxy = nil
	transport.DialContext = dialLoopback
	// Restrict this development transport to HTTP/1 so cloned HTTP/2 protocol
	// callbacks cannot hand authenticated requests to an opaque RoundTripper.
	transport.TLSNextProto = nil
	transport.Protocols = &http.Protocols{}
	transport.Protocols.SetHTTP1(true)
	if transport.IdleConnTimeout <= 0 {
		transport.IdleConnTimeout = 90 * time.Second
	}
	entry := &c.transports[c.next]
	if entry.transport != nil {
		// Evict only our clone; the caller owns the source and its connections.
		entry.transport.CloseIdleConnections()
	}
	entry.source, entry.transport = source, transport
	c.next = (c.next + 1) % len(c.transports)
	return transport, nil
}

func dialLoopback(ctx context.Context, network, address string) (net.Conn, error) {
	host, port, err := net.SplitHostPort(address)
	if err != nil || !isLoopbackHost(host) {
		return nil, errors.New("bedrock: refusing a non-loopback plaintext connection")
	}
	hosts := []string{host}
	if host == "localhost" {
		// Never resolve localhost through DNS or the caller's custom dialer.
		hosts = []string{"127.0.0.1", "::1"}
	}
	dialer := &net.Dialer{Timeout: 30 * time.Second, KeepAlive: 30 * time.Second}
	for _, ip := range hosts {
		var conn net.Conn
		conn, err = dialer.DialContext(ctx, network, net.JoinHostPort(ip, port))
		if err == nil {
			return conn, nil
		}
	}
	return nil, err
}

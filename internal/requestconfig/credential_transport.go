package requestconfig

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Only opted-in requests can use this pool. Creating a request option must
// not create a separate idle connection pool for every call.
var directLoopbackTransport = &http.Transport{
	DialContext:           dialCredentialLoopback,
	MaxIdleConns:          100,
	IdleConnTimeout:       90 * time.Second,
	ResponseHeaderTimeout: 10 * time.Minute,
	ExpectContinueTimeout: time.Second,
}

// WithUnsafeAllowHTTP permits authenticated local development over direct loopback.
func WithUnsafeAllowHTTP() RequestOption {
	return RequestOptionFunc(func(cfg *RequestConfig) error {
		cfg.allowLoopbackHTTP = true
		return nil
	})
}

var errCredentialTransport = WithNoRetryError(errors.New("openai: authenticated requests require HTTPS; option.WithUnsafeAllowHTTP permits HTTP only for direct loopback development endpoints"))

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
	return errCredentialTransport
}

// A middleware replacement can be shared by two request clones without
// sharing the SDK's close-once wrapper. Let Execute recognize the body that
// the rejection already closed, even if the clone has a different context.
type rejectedCredentialBody struct {
	error
	body atomic.Pointer[io.ReadCloser]
}

func (e *rejectedCredentialBody) Unwrap() error { return e.error }

func rejectCredentialRequest(req *http.Request, err error) error {
	if req == nil || req.Body == nil || req.Body == http.NoBody {
		return err
	}
	body := req.Body
	closeRequestBody(req)
	if _, wrapped := body.(*closeOnceReadCloser); wrapped {
		return err
	}
	rejected := &rejectedCredentialBody{error: err}
	rejected.body.Store(&body)
	return rejected
}

// Middleware may translate a rejection into an unrelated error. Remember the
// closed reader as the rejection crosses back out of each SDK middleware call.
type rejectedBodies struct {
	closed sync.Map // reflect.Value of a boxed io.ReadCloser, including non-comparable readers
}

func (b *rejectedBodies) record(err error) {
	var rejected *rejectedCredentialBody
	if err != nil && errors.As(err, &rejected) {
		if body := rejected.body.Swap(nil); body != nil {
			b.closed.Store(reflect.ValueOf(*body), true)
		}
	}
}

func (b *rejectedBodies) observe(next middlewareNext) middlewareNext {
	return func(req *http.Request) (*http.Response, error) {
		res, err := next(req)
		b.record(err)
		return res, err
	}
}

func (b *rejectedBodies) finishAttempt(err error, bodies ...io.ReadCloser) {
	if err != nil {
		b.record(err)
		for _, body := range bodies {
			if body != nil {
				if _, closed := b.closed.LoadOrStore(reflect.ValueOf(body), true); !closed {
					_ = body.Close()
				}
			}
		}
	}
	b.closed.Clear()
}

// ValidateOpenAICredentialRequest checks the request passed to SDK-managed
// workload authentication before acquiring a token.
func (cfg *RequestConfig) ValidateOpenAICredentialRequest(req *http.Request) error {
	if !requestHasCanonicalTarget(req) {
		return rejectCredentialRequest(req, errCredentialTransport)
	}
	base := cfg.credentialBaseURL()
	allowLoopback := cfg.allowLoopbackHTTP && SameOrigin(req.URL, base)
	if err := validateCredentialURL(req.URL, allowLoopback); err != nil {
		return rejectCredentialRequest(req, err)
	}
	if !SameOrigin(req.URL, base) {
		return rejectCredentialRequest(req, requestOriginError())
	}
	return nil
}

// UseSelectedOpenAIWorkloadTransport applies OpenAI transport rules when the
// effective authentication option acquires an OpenAI workload token.
func (cfg *RequestConfig) UseSelectedOpenAIWorkloadTransport() {
	_, cfg.requireSecureTransport = cfg.ProviderAuth("OpenAI")
}

func (cfg *RequestConfig) managesCredentialTransport() bool {
	return cfg.endpointProvider == "" || cfg.requireSecureTransport
}

func (cfg *RequestConfig) configureCredentialTransport() error {
	if cfg.endpointProvider == "" {
		cfg.UseSelectedOpenAIWorkloadTransport()
	}
	if !cfg.managesCredentialTransport() {
		return nil
	}
	// Check the destination before ApplySecurity attaches a key or workload
	// middleware obtains a token. Explicit header options have already run.
	if !cfg.requireSecureTransport && cfg.staticBearerKey() == "" {
		return nil
	}
	base := cfg.credentialBaseURL()
	if base == nil {
		return nil // Execute and PrepareWebSocket report unconfigured endpoints.
	}
	return validateCredentialURL(base, cfg.allowLoopbackHTTP)
}

func (cfg *RequestConfig) credentialBaseURL() *url.URL {
	base := cfg.BaseURL
	if base == nil {
		base = cfg.DefaultBaseURL
	}
	if base == nil {
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
	return base
}

func dialCredentialLoopback(ctx context.Context, network, address string) (net.Conn, error) {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return nil, errCredentialTransport
	}
	dialer := net.Dialer{Timeout: 30 * time.Second, KeepAlive: 30 * time.Second}
	if !strings.EqualFold(host, "localhost") {
		ip := net.ParseIP(host)
		if ip == nil || !ip.IsLoopback() {
			return nil, errCredentialTransport
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

func hasAuthorization(req *http.Request) bool {
	// Trailer values may be supplied while streaming the body, after dispatch.
	for name := range req.Trailer {
		if strings.EqualFold(name, "Authorization") {
			return true
		}
	}
	for name, values := range req.Header {
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
	if !cfg.managesCredentialTransport() || (!cfg.requireSecureTransport && !hasAuthorization(req) && (req.URL == nil || req.URL.User == nil)) {
		return nil, nil
	}
	if !requestHasCanonicalTarget(req) {
		return nil, errCredentialTransport
	}
	if err := validateCredentialURL(req.URL, cfg.allowLoopbackHTTP); err != nil {
		return nil, err
	}
	if req.URL.Scheme == "http" {
		return directLoopbackTransport, nil
	}
	return nil, nil
}

func (cfg *RequestConfig) credentialTransportForDispatch(req *http.Request) (http.RoundTripper, error) {
	direct, err := cfg.credentialTransport(req)
	if err != nil {
		return nil, rejectCredentialRequest(req, err)
	}
	return direct, nil
}

// newHTTPHandler composes the per-attempt middleware and transport guards.
func (cfg *RequestConfig) newHTTPHandler(rejected *rejectedBodies) middlewareNext {
	client := *cfg.HTTPClient
	transport := client.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}
	if cfg.managesCredentialTransport() {
		guard := &credentialHTTPTransport{cfg: cfg, next: transport}
		if client.Timeout > 0 {
			guard.legacy, _ = transport.(interface{ CancelRequest(*http.Request) })
		}
		transport = guard
	}
	if hasCredentialRedirectGuard(transport) {
		client.Transport = transport
	} else {
		client.Transport = originTransport{origin: cfg.BaseURL, next: transport}
	}
	// Provider transports may implement a narrower, explicitly configured
	// redirect policy. Every initial request must still retain the configured
	// origin after all SDK and caller middleware has run.
	handler := enforceRequestOrigin(cfg.BaseURL, func(req *http.Request) (*http.Response, error) {
		direct, transportErr := cfg.credentialTransportForDispatch(req)
		if transportErr != nil {
			return nil, transportErr
		}
		selected := client
		if direct != nil {
			// The native client owns the full redirect exchange. Even a hop
			// without credentials stays on the direct same-origin connection path.
			selected.Transport = originTransport{origin: cfg.BaseURL, next: direct}
		} else if cfg.CustomHTTPDoer != nil {
			return cfg.CustomHTTPDoer.Do(req)
		}
		res, err := selected.Do(req)
		if errors.Is(err, errCredentialTransport) {
			// net/http includes the rejected redirect URL in its error.
			// Return the SDK policy failure without the URL's userinfo or query.
			var urlErr *url.Error
			if errors.As(err, &urlErr) {
				return res, urlErr.Err
			}
		}
		return res, err
	})
	for i := len(cfg.Middlewares) - 1; i >= 0; i -= 1 {
		handler = applyMiddleware(cfg.Middlewares[i], rejected.observe(handler))
	}
	if cfg.requireSecureTransport {
		next := handler
		handler = func(req *http.Request) (*http.Response, error) {
			if transportErr := cfg.ValidateOpenAICredentialRequest(req); transportErr != nil {
				return nil, transportErr
			}
			return next(req)
		}
	}
	return handler
}

// Native clients can gain credentials after SDK middleware has run, for
// example from redirect userinfo or CheckRedirect. Enforce the same policy at
// each transport call, after net/http has applied those credentials.
type credentialHTTPTransport struct {
	cfg  *RequestConfig
	next http.RoundTripper

	// net/http uses CancelRequest for legacy transports when Client.Timeout
	// is set. The guard belongs to one Execute call, including its redirects.
	legacy         interface{ CancelRequest(*http.Request) }
	legacyRequests sync.Map // *http.Request -> true if dispatched, false if canceled
}

func (t *credentialHTTPTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	direct, err := t.cfg.credentialTransportForDispatch(req)
	if err != nil {
		return nil, err
	}
	if direct != nil {
		return direct.RoundTrip(req)
	}
	if t.legacy != nil {
		// net/http can cancel while the credential check is still running.
		if state, loaded := t.legacyRequests.LoadOrStore(req, true); loaded && !state.(bool) {
			err := req.Context().Err()
			if err == nil {
				err = context.Canceled
			}
			return nil, rejectCredentialRequest(req, err)
		}
	}
	return t.next.RoundTrip(req)
}

func (t *credentialHTTPTransport) CancelRequest(req *http.Request) {
	// The direct loopback transport observes req.Cancel and its context. Only
	// forward the legacy hook for a request dispatched to that exact transport.
	if t.legacy == nil {
		return
	}
	if state, loaded := t.legacyRequests.Swap(req, false); loaded && state.(bool) {
		t.legacy.CancelRequest(req)
	}
}

package requestconfig

import (
	"context"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/openai/openai-go/v3/internal/websockettransport"
)

// WebSocketHTTPClient is an optional capability for bespoke HTTPDoers. The
// returned client must preserve the doer's transport and authentication policy
// and support upgrades with writable response bodies. The SDK never closes it.
type WebSocketHTTPClient interface {
	WebSocketHTTPClient() *http.Client
}

type websocketRoundTripper func(*http.Request) (*http.Response, error)

func (f websocketRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	response, err := f(r)
	// Client.Do may return both a closed redirect response and an error, but
	// RoundTrip must return only the error. Also release any middleware body.
	if err != nil && response != nil {
		if response.Body != nil {
			_ = response.Body.Close()
		}
		return nil, err
	}
	return response, err
}

// RecordQueryChange remembers whether explicit options append to or replace
// a base-URL query value. Ordinary HTTP request construction remains unchanged.
func (cfg *RequestConfig) RecordQueryChange(key string, replace bool) {
	if cfg.queryChanges == nil {
		cfg.queryChanges = make(map[string]bool)
	}
	cfg.queryChanges[key] = cfg.queryChanges[key] || replace
}

// PrepareWebSocket prepares an upgrade without response decoding. Only a rejected
// workload-identity handshake can retry, within its existing request retry scope.
// Provider-specific clients require an independently tested adapter.
func (cfg *RequestConfig) PrepareWebSocket() (*http.Request, *http.Client, error) {
	return cfg.PrepareWebSocketWithRecovery(false)
}

// PrepareWebSocketWithRecovery rejects excluded session selectors on every
// recovery-enabled upgrade before handing it to the configured transport.
func (cfg *RequestConfig) PrepareWebSocketWithRecovery(recoveryEnabled bool, excludedQueryParams ...string) (*http.Request, *http.Client, error) {
	if cfg.cloneError != nil {
		return nil, nil, cfg.cloneError
	}
	if cfg.EndpointProvider() != "" {
		return nil, nil, errors.New("websocket: this provider does not support Responses upgrades")
	}
	base := cfg.BaseURL
	if base == nil {
		base = cfg.DefaultBaseURL
	}
	if base == nil {
		return nil, nil, errors.New("websocket: base URL is not configured")
	}
	origin := *base
	switch origin.Scheme {
	case "ws":
		origin.Scheme = "http"
	case "wss":
		origin.Scheme = "https"
	case "http", "https":
	default:
		return nil, nil, errors.New("websocket: unsupported URL scheme")
	}
	if origin.User != nil || origin.Fragment != "" {
		return nil, nil, errors.New("websocket: URL must not contain user information or a fragment")
	}
	target, err := origin.Parse(strings.TrimLeft(cfg.Request.URL.String(), "/"))
	if err != nil {
		return nil, nil, err
	}
	query := origin.Query()
	targetQuery := target.Query()
	for key, replace := range cfg.queryChanges {
		if replace {
			delete(query, key)
		} else {
			query[key] = append(query[key], targetQuery[key]...)
			delete(targetQuery, key)
		}
	}
	for key, values := range targetQuery {
		query[key] = values
	}
	target.RawQuery = query.Encode()
	request := cfg.Request.Clone(cfg.Request.Context())
	request.URL = target
	selected := cfg.HTTPClient
	if cfg.CustomHTTPDoer != nil {
		capable, ok := cfg.CustomHTTPDoer.(WebSocketHTTPClient)
		if !ok {
			return nil, nil, errors.New("websocket: custom HTTP client must implement WebSocketHTTPClient")
		}
		selected = capable.WebSocketHTTPClient()
	}
	if selected == nil {
		return nil, nil, errors.New("websocket: nil HTTP client")
	}
	client := *selected
	transport := client.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}
	client.Transport = websocketRoundTripper(enforceRequestOrigin(&origin, func(req *http.Request) (*http.Response, error) {
		if recoveryEnabled && len(excludedQueryParams) > 0 {
			query := req.URL.Query()
			for _, key := range excludedQueryParams {
				if query.Has(key) {
					return nil, errors.New("websocket: automatic recovery does not support this session selector")
				}
			}
		}
		return websockettransport.RoundTrip(transport, req)
	}))
	previousRedirect := client.CheckRedirect
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		switch req.URL.Scheme {
		case "ws":
			req.URL.Scheme = "http"
		case "wss":
			req.URL.Scheme = "https"
		}
		if !RequestHasOrigin(req, &origin) {
			return requestOriginError()
		}
		if previousRedirect != nil {
			return previousRedirect(req, via)
		}
		if len(via) >= 10 {
			return errors.New("websocket: stopped after 10 redirects")
		}
		return nil
	}
	// Run the SDK middleware once around the complete HTTP exchange, just like
	// Execute. Redirect hops retain their transport-level origin validation.
	// Only the outer client's timeout bounds opening; neither client may keep
	// a body timer running after the response becomes an established socket.
	client.Timeout = 0
	handler := enforceRequestOrigin(&origin, client.Do)
	for i := len(cfg.Middlewares) - 1; i >= 0; i-- {
		handler = applyMiddleware(cfg.Middlewares[i], handler)
	}
	// Authenticate only the endpoint prepared by this binding, even if middleware
	// changes the request URL before the credential exchange.
	expectedPath := websocketUpgradePath{path: target.Path, escapedPath: target.EscapedPath()}
	handshakeClient := &http.Client{
		Transport: websocketRoundTripper(func(req *http.Request) (*http.Response, error) {
			req = req.WithContext(context.WithValue(req.Context(), websocketUpgradeContextKey{}, expectedPath))
			var cancel context.CancelFunc
			defer func() {
				if cancel != nil {
					cancel()
				}
			}()
			for retryCount := 0; ; retryCount++ {
				attemptContext := req.Context()
				if cfg.RequestTimeout > 0 {
					attemptContext, cancel = context.WithTimeout(attemptContext, cfg.RequestTimeout)
				}
				attempt := req.Clone(attemptContext)
				if req.Header.Get("X-Stainless-Retry-Count") == "0" {
					attempt.Header.Set("X-Stainless-Retry-Count", strconv.Itoa(retryCount))
				}
				response, err := handler(attempt)
				if cfg.ResponseInto != nil {
					*cfg.ResponseInto = response
				}
				if attemptErr := attemptContext.Err(); attemptErr != nil {
					err = attemptErr
				}
				if err == nil && response != nil && response.StatusCode == http.StatusUnauthorized &&
					response.Header.Get("X-Should-Retry") == "true" &&
					RequestRetryScopeFromContext(req.Context()) != nil && retryCount < cfg.MaxRetries {
					if delay, retry := retryDelay(response, retryCount, cfg.MaxRetryDelay); retry {
						if response.Body != nil {
							_ = response.Body.Close()
						}
						// As in Execute, rejection and its backoff belong to this attempt.
						waitErr := WaitForDelay(attemptContext, delay)
						if cancel != nil {
							cancel()
							cancel = nil
						}
						if waitErr != nil {
							return nil, waitErr
						}
						continue
					}
				}
				if cancel != nil && response != nil && response.Body != nil {
					if body, ok := response.Body.(io.ReadWriteCloser); ok && err == nil && response.StatusCode == http.StatusSwitchingProtocols {
						response.Body = &websocketOpeningBody{ReadWriteCloser: body, cancel: cancel}
					} else {
						// Release rejected bodies before canceling their opening request.
						response.Body = &bodyWithTimeout{rc: response.Body, stop: cancel}
					}
					cancel = nil // Ownership moves to the response body or upgraded socket.
				}
				return response, err
			}
		}),
		Timeout: selected.Timeout,
		// The inner client already applied the caller's redirect policy.
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
	}
	return request, handshakeClient, nil
}

// websocketOpeningBody keeps the attempt alive through handshake validation.
// It must not add write support to a non-writable response body.
type websocketOpeningBody struct {
	io.ReadWriteCloser
	cancel context.CancelFunc
}

func (b *websocketOpeningBody) CompleteWebSocketOpening() { b.cancel() }

func (b *websocketOpeningBody) Close() error {
	defer b.CompleteWebSocketOpening()
	return b.ReadWriteCloser.Close()
}

// websocketUpgradeContextKey marks an opening prepared by the SDK, not an
// ordinary HTTP request with user-supplied protocol upgrade headers.
type websocketUpgradeContextKey struct{}

type websocketUpgradePath struct {
	path        string
	escapedPath string
}

// WebSocketValidationRequest returns an ordinary HTTP view of an SDK-owned
// upgrade for authentication policy checks. It validates the complete
// upgrade framing before omitting only the two protocol-switch headers. The
// actual request and its headers remain unchanged.
func WebSocketValidationRequest(request *http.Request) (*http.Request, error) {
	if request == nil {
		return request, nil
	}
	expectedPath, prepared := request.Context().Value(websocketUpgradeContextKey{}).(websocketUpgradePath)
	if !prepared {
		return request, nil
	}
	if request.URL == nil || request.Method != http.MethodGet || request.URL.Path != expectedPath.path ||
		request.URL.EscapedPath() != expectedPath.escapedPath || request.ContentLength != 0 ||
		(request.Body != nil && request.Body != http.NoBody) {
		return nil, errors.New("websocket: X.509 upgrades require a bodyless GET to " + expectedPath.path)
	}
	counts := make(map[string]int)
	for name, values := range request.Header {
		normalized := strings.ToLower(strings.ReplaceAll(name, "_", "-"))
		switch normalized {
		case "connection", "upgrade", "sec-websocket-version", "sec-websocket-key":
			counts[normalized] += len(values)
			if name != http.CanonicalHeaderKey(normalized) || len(values) != 1 {
				return nil, errors.New("websocket: invalid X.509 upgrade headers")
			}
		}
	}
	key, err := base64.StdEncoding.DecodeString(request.Header.Get("Sec-WebSocket-Key"))
	if err != nil || len(key) != 16 || counts["connection"] != 1 || counts["upgrade"] != 1 ||
		counts["sec-websocket-key"] != 1 || counts["sec-websocket-version"] != 1 ||
		!strings.EqualFold(request.Header.Get("Connection"), "Upgrade") ||
		!strings.EqualFold(request.Header.Get("Upgrade"), "websocket") ||
		request.Header.Get("Sec-WebSocket-Version") != "13" {
		return nil, errors.New("websocket: invalid X.509 upgrade headers")
	}
	validation := request.Clone(request.Context())
	validation.Header.Del("Connection")
	validation.Header.Del("Upgrade")
	return validation, nil
}

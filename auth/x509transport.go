package auth

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"crypto/tls"
	"encoding/binary"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptrace"
	"path"
	"slices"
	"strings"
	"sync/atomic"
)

const (
	x509AuthenticationHost = "mtls.auth.openai.com"
	x509APIHost            = "mtls.api.openai.com"
)

var errX509Redirect = errors.New("X.509 workload identity does not follow redirects")

type x509TransportError struct {
	cause error
}

func (err *x509TransportError) Error() string { return "X.509 transport request failed" }
func (err *x509TransportError) Unwrap() error { return err.cause }

// x509TraceFreeContext preserves cancellation, deadlines, and ordinary context
// values while preventing a mutable parent from revealing HTTP trace hooks
// after the request's final security checks.
type x509TraceFreeContext struct {
	context.Context
}

func (ctx x509TraceFreeContext) Value(key any) any {
	value := ctx.Context.Value(key)
	if _, ok := value.(*httptrace.ClientTrace); ok {
		return nil
	}
	return value
}

// X509Transport represents one explicitly attested mTLS configuration
// generation. Its identity is the capability pointer, not a certificate-bound
// OAuth token. Rotating credentials requires a new transport and capability.
//
// The application owns its template transport, certificate, private key, and
// trust roots. The template is never used, modified, or closed by this
// capability. X509Transport creates an isolated native connection pool that
// cannot inherit privately registered HTTPS handlers; call Close to release
// that pool. The template and its TLS configuration must remain unchanged after
// attestation. HTTPS proxies are unsupported because net/http uses the same
// client TLS configuration for a proxy and its origin.
type X509Transport struct {
	template          *http.Transport
	templateTLS       *tls.Config
	transport         *http.Transport
	tlsConfig         *tls.Config
	certificateDigest [sha256.Size]byte
	client            *http.Client
	closed            atomic.Bool
}

// NewX509Transport attests a direct, native HTTP transport template configured
// with one static client certificate. The template and TLS credentials remain
// caller-owned; the returned capability owns a separate clean connection pool.
func NewX509Transport(template *http.Transport) (*X509Transport, error) {
	if err := validateX509NativeTransport(template); err != nil {
		return nil, err
	}
	config := template.TLSClientConfig.Clone()
	config.Certificates = slices.Clone(config.Certificates)
	for index := range config.Certificates {
		certificate := &config.Certificates[index]
		certificate.Certificate = slices.Clone(certificate.Certificate)
		for chainIndex := range certificate.Certificate {
			certificate.Certificate[chainIndex] = slices.Clone(certificate.Certificate[chainIndex])
		}
		certificate.OCSPStaple = slices.Clone(certificate.OCSPStaple)
		certificate.SignedCertificateTimestamps = slices.Clone(certificate.SignedCertificateTimestamps)
		for timestampIndex := range certificate.SignedCertificateTimestamps {
			certificate.SignedCertificateTimestamps[timestampIndex] = slices.Clone(
				certificate.SignedCertificateTimestamps[timestampIndex],
			)
		}
	}
	if config.RootCAs != nil {
		config.RootCAs = config.RootCAs.Clone()
	}
	config.NextProtos = slices.Clone(config.NextProtos)
	transport := &http.Transport{
		DialContext:            template.DialContext,
		TLSClientConfig:        config,
		TLSHandshakeTimeout:    template.TLSHandshakeTimeout,
		DisableKeepAlives:      template.DisableKeepAlives,
		DisableCompression:     template.DisableCompression,
		MaxIdleConns:           template.MaxIdleConns,
		MaxIdleConnsPerHost:    template.MaxIdleConnsPerHost,
		MaxConnsPerHost:        template.MaxConnsPerHost,
		IdleConnTimeout:        template.IdleConnTimeout,
		ResponseHeaderTimeout:  template.ResponseHeaderTimeout,
		ExpectContinueTimeout:  template.ExpectContinueTimeout,
		MaxResponseHeaderBytes: template.MaxResponseHeaderBytes,
		WriteBufferSize:        template.WriteBufferSize,
		ReadBufferSize:         template.ReadBufferSize,
		ForceAttemptHTTP2:      template.ForceAttemptHTTP2,
	}
	//nolint:staticcheck // Preserve a legacy TCP dialer; unlike DialTLS, it cannot bypass native TLS.
	transport.Dial = template.Dial
	if template.Protocols != nil {
		protocols := *template.Protocols
		transport.Protocols = &protocols
	}
	if template.HTTP2 != nil {
		configuration := *template.HTTP2
		transport.HTTP2 = &configuration
	}

	return &X509Transport{
		template:          template,
		templateTLS:       template.TLSClientConfig,
		transport:         transport,
		tlsConfig:         config,
		certificateDigest: x509CertificateDigest(config.Certificates[0].Certificate),
		client: &http.Client{
			Transport: transport,
			CheckRedirect: func(*http.Request, []*http.Request) error {
				return errX509Redirect
			},
		},
	}, nil
}

func x509CertificateDigest(chain [][]byte) [sha256.Size]byte {
	digest := sha256.New()
	var length [8]byte
	for _, certificate := range chain {
		binary.BigEndian.PutUint64(length[:], uint64(len(certificate)))
		_, _ = digest.Write(length[:])
		_, _ = digest.Write(certificate)
	}
	var result [sha256.Size]byte
	copy(result[:], digest.Sum(nil))
	return result
}

func validateX509NativeTransport(transport *http.Transport) error {
	if transport == nil {
		return errors.New("X.509 transport requires a non-nil native HTTP transport")
	}
	if transport.TLSClientConfig == nil {
		return errors.New("X.509 transport requires an explicit TLS configuration")
	}
	config := transport.TLSClientConfig
	if len(config.Certificates) != 1 || len(config.Certificates[0].Certificate) == 0 ||
		len(config.Certificates[0].Certificate[0]) == 0 || config.Certificates[0].PrivateKey == nil {
		return errors.New("X.509 transport requires exactly one static certificate and private key")
	}
	if config.GetClientCertificate != nil {
		return errors.New("X.509 transport does not support dynamic client-certificate callbacks")
	}
	if transport.DialTLSContext != nil {
		return errors.New("X.509 transport does not support custom TLS dialers")
	}
	//nolint:staticcheck // Deprecated DialTLS still bypasses TLSClientConfig and must be rejected.
	if transport.DialTLS != nil {
		return errors.New("X.509 transport does not support deprecated TLS dialers")
	}
	if config.ClientSessionCache != nil {
		return errors.New("X.509 transport does not support shared TLS session caches")
	}
	if transport.Proxy != nil {
		return errors.New("X.509 transport currently supports direct connections only")
	}
	for protocol, handler := range transport.TLSNextProto {
		if (protocol != "h2" && protocol != "unencrypted_http2") || handler == nil {
			return errors.New("X.509 transport does not support custom TLS protocol handlers")
		}
	}
	if len(transport.TLSNextProto) != 0 && transport.TLSNextProto["h2"] == nil {
		return errors.New("X.509 transport does not support custom TLS protocol handlers")
	}
	if transport.HTTP2 != nil && transport.HTTP2.CountError != nil {
		return errors.New("X.509 transport does not support HTTP/2 error callbacks")
	}
	if config.InsecureSkipVerify {
		return errors.New("X.509 transport requires TLS certificate and hostname verification")
	}
	if config.ServerName != "" {
		return errors.New("X.509 transport does not support overriding the TLS server name")
	}
	return nil
}

func (transport *X509Transport) validateAttestation() error {
	if transport == nil || transport.client == nil || transport.transport == nil || transport.template == nil {
		return errors.New("X.509 transport capability is invalid")
	}
	if transport.closed.Load() {
		return errors.New("X.509 transport capability is closed")
	}
	if err := validateX509NativeTransport(transport.template); err != nil {
		return fmt.Errorf("X.509 transport attestation changed: %w", err)
	}
	if transport.template.TLSClientConfig != transport.templateTLS ||
		transport.transport.TLSClientConfig != transport.tlsConfig {
		return errors.New("X.509 transport TLS configuration changed after attestation")
	}
	for _, config := range []*tls.Config{transport.templateTLS, transport.tlsConfig} {
		if len(config.Certificates) != 1 {
			return errors.New("X.509 transport certificate changed after attestation")
		}
		digest := x509CertificateDigest(config.Certificates[0].Certificate)
		if subtle.ConstantTimeCompare(digest[:], transport.certificateDigest[:]) != 1 {
			return errors.New("X.509 transport certificate changed after attestation")
		}
	}
	return nil
}

// Close releases this capability's isolated idle connections. It is
// idempotent, never closes the caller's template, and prevents future requests.
func (transport *X509Transport) Close() error {
	if transport == nil || transport.transport == nil {
		return errors.New("X.509 transport capability is invalid")
	}
	if transport.closed.CompareAndSwap(false, true) {
		transport.transport.CloseIdleConnections()
	}
	return nil
}

// Do sends a request to an approved global OpenAI mTLS endpoint without
// following redirects. The request is snapshotted before validation so caller
// hooks cannot alter its URL or credentials after the final safety checks.
// HTTP tracing callbacks and request trailers are unsupported because mutable
// trace hooks can expose a live connection after request validation.
// Request context and caller-owned TLS credentials are preserved; the
// capability owns its isolated pool. OAuth authentication is configured
// separately.
func (transport *X509Transport) Do(request *http.Request) (*http.Response, error) {
	if request == nil {
		return nil, errors.New("X.509 transport requires a non-nil HTTP request")
	}
	body := request.Body
	bodyTransferred := false
	defer func() {
		if !bodyTransferred && body != nil {
			_ = body.Close()
		}
	}()
	if err := request.Context().Err(); err != nil {
		return nil, err
	}
	if httptrace.ContextClientTrace(request.Context()) != nil {
		return nil, errors.New("X.509 transport does not support HTTP trace callbacks")
	}
	request = request.Clone(x509TraceFreeContext{request.Context()})
	if err := transport.validateAttestation(); err != nil {
		return nil, err
	}
	if err := validateX509Request(request); err != nil {
		return nil, err
	}
	bodyTransferred = true
	response, err := transport.client.Do(request)
	if err != nil {
		redacted := &x509TransportError{}
		switch {
		case errors.Is(err, errX509Redirect):
			redacted.cause = errX509Redirect
		case errors.Is(err, context.Canceled):
			redacted.cause = context.Canceled
		case errors.Is(err, context.DeadlineExceeded):
			redacted.cause = context.DeadlineExceeded
		}
		return nil, redacted
	}
	if response.StatusCode == http.StatusSwitchingProtocols {
		_ = response.Body.Close()
		return nil, errors.New("X.509 transport does not support protocol upgrades")
	}
	if response.StatusCode >= http.StatusMultipleChoices &&
		response.StatusCode < http.StatusBadRequest && response.StatusCode != http.StatusNotModified {
		_ = response.Body.Close()
		return nil, &x509TransportError{cause: errX509Redirect}
	}
	return response, nil
}

func validateX509Request(request *http.Request) error {
	if len(request.Trailer) != 0 {
		return errors.New("X.509 requests do not support HTTP trailers")
	}
	if request.URL == nil || request.URL.Scheme != "https" || request.URL.User != nil ||
		request.URL.Opaque != "" || request.URL.Fragment != "" || request.URL.RawFragment != "" {
		return errors.New("X.509 requests require an absolute HTTPS URL without credentials or fragments")
	}
	host := request.URL.Hostname()
	if host != x509AuthenticationHost && host != x509APIHost {
		return errors.New("X.509 requests are limited to approved global OpenAI mTLS origins")
	}
	if request.URL.Port() != "" && request.URL.Port() != "443" {
		return errors.New("X.509 requests require the default HTTPS port")
	}
	if request.URL.Host != host && request.URL.Host != host+":443" {
		return errors.New("X.509 requests require an exact mTLS URL authority")
	}
	if request.Host != "" && request.Host != request.URL.Host {
		return errors.New("X.509 request Host must match its URL authority")
	}

	authorizationCount := 0
	for name, values := range request.Header {
		normalized := strings.ToLower(strings.ReplaceAll(name, "_", "-"))
		switch normalized {
		case "api-key", "x-api-key", "proxy-authorization", "cookie", "set-cookie", "host", "x-amz-security-token":
			return errors.New("X.509 requests cannot contain alternate credentials, cookies, or authority headers")
		case "authorization":
			authorizationCount += len(values)
			if host == x509AuthenticationHost || len(values) != 1 ||
				authorizationCount != 1 || !validX509BearerHeader(values[0]) {
				return errors.New("X.509 requests contain invalid Authorization credentials")
			}
		case "openai-organization", "openai-project":
			if host == x509AuthenticationHost {
				return errors.New("X.509 token exchange cannot contain organization or project headers")
			}
		}
		if strings.HasPrefix(normalized, ":") {
			return errors.New("X.509 requests cannot override HTTP authority pseudo-headers")
		}
	}

	if host == x509AuthenticationHost {
		if request.Method != http.MethodPost || request.URL.Path != "/oauth/token" ||
			request.URL.EscapedPath() != "/oauth/token" || request.URL.RawQuery != "" || request.URL.ForceQuery {
			return errors.New("X.509 token exchange requires POST to the exact pinned OAuth endpoint")
		}
		return nil
	}
	if !strings.HasPrefix(request.URL.Path, "/v1/") ||
		!strings.HasPrefix(request.URL.EscapedPath(), "/v1/") || path.Clean(request.URL.Path) != request.URL.Path {
		return errors.New("X.509 API requests must remain inside the /v1/ path")
	}
	return nil
}

func validX509BearerHeader(value string) bool {
	token, ok := strings.CutPrefix(value, "Bearer ")
	if !ok || token == "" {
		return false
	}
	padding := false
	for _, value := range []byte(token) {
		switch {
		case value == '=':
			padding = true
		case value >= 'A' && value <= 'Z', value >= 'a' && value <= 'z', value >= '0' && value <= '9',
			value == '-', value == '.', value == '_', value == '~', value == '+', value == '/':
			if padding {
				return false
			}
		default:
			return false
		}
	}
	return token[0] != '='
}

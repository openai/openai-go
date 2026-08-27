package auth

import (
	"context"
	"crypto"
	"crypto/sha256"
	"crypto/subtle"
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptrace"
	"path"
	"reflect"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/openai/openai-go/v3/internal/requestconfig"
)

const (
	x509AuthenticationHost           = "mtls.auth.openai.com"
	x509APIHost                      = "mtls.api.openai.com"
	x509DefaultDialTimeout           = 30 * time.Second
	x509DefaultTLSHandshakeTimeout   = 10 * time.Second
	x509DefaultResponseHeaderTimeout = 10 * time.Minute
	x509MaximumConcurrentDials       = 32
)

var (
	errX509Redirect        = errors.New("X.509 workload identity does not follow redirects")
	errX509TLSVerification = errors.New("X.509 transport TLS verification failed")
)

type x509DialAdmissionContextKey struct{}

type x509TransportError struct {
	cause     error
	timeout   bool
	temporary bool
	retryable bool
}

func (err *x509TransportError) Error() string   { return "X.509 transport request failed" }
func (err *x509TransportError) Unwrap() error   { return err.cause }
func (err *x509TransportError) Timeout() bool   { return err.timeout }
func (err *x509TransportError) Temporary() bool { return err.temporary }

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
// trust roots. The template is never used for network I/O, modified, or closed
// by this capability. X509Transport creates an isolated native connection pool
// that cannot inherit privately registered HTTPS handlers; call Close to
// release that pool. The template and its TLS configuration are revalidated and
// must remain unchanged after attestation. HTTPS proxies are unsupported because
// net/http uses the same client TLS configuration for a proxy and its origin.
type X509Transport struct {
	template          *http.Transport
	templateTLS       *tls.Config
	transport         *http.Transport
	tlsConfig         *tls.Config
	certificateDigest [sha256.Size]byte
	certificateSelect uintptr
	client            *http.Client
	lifecycle         *x509TransportLifecycle
}

type x509TransportLifecycle struct {
	closed         atomic.Bool
	mu             sync.Mutex
	activeRequests int
	connections    map[*x509TrackedConnection]struct{}
	dialSemaphore  chan struct{}
	dialContext    func(context.Context, string, string) (net.Conn, error)
	transport      *http.Transport
}

// NewX509Transport attests a direct, native HTTP transport template configured
// with one static client certificate. The template and TLS credentials remain
// caller-owned; the returned capability owns a separate clean connection pool.
func NewX509Transport(template *http.Transport) (*X509Transport, error) {
	if err := validateX509NativeTransport(template); err != nil {
		return nil, err
	}
	if err := validateX509ApplicationProtocols(template); err != nil {
		return nil, err
	}
	if err := validateX509TLSProtocolHandlers(template.TLSNextProto); err != nil {
		return nil, err
	}
	config := template.TLSClientConfig.Clone()
	config.Certificates = slices.Clone(config.Certificates)
	for index := range config.Certificates {
		certificate := &config.Certificates[index]
		certificate.Leaf = nil
		certificate.Certificate = slices.Clone(certificate.Certificate)
		for chainIndex := range certificate.Certificate {
			certificate.Certificate[chainIndex] = slices.Clone(certificate.Certificate[chainIndex])
		}
		certificate.OCSPStaple = slices.Clone(certificate.OCSPStaple)
		certificate.SupportedSignatureAlgorithms = slices.Clone(certificate.SupportedSignatureAlgorithms)
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
	config.CipherSuites = slices.Clone(config.CipherSuites)
	config.CurvePreferences = slices.Clone(config.CurvePreferences)
	config.EncryptedClientHelloConfigList = slices.Clone(config.EncryptedClientHelloConfigList)
	// The sole attested certificate remains static even when the server advertises
	// an unrelated acceptable-CA hint. Caller-controlled selectors remain forbidden.
	config.GetClientCertificate = func(*tls.CertificateRequestInfo) (*tls.Certificate, error) {
		return &config.Certificates[0], nil
	}
	dialContext := template.DialContext
	if dialContext == nil {
		//nolint:staticcheck // Preserve the legacy TCP dialer accepted by the original X.509 contract.
		if legacyDial := template.Dial; legacyDial != nil {
			dialContext = func(_ context.Context, network, address string) (net.Conn, error) {
				return legacyDial(network, address)
			}
		} else {
			dialContext = (&net.Dialer{
				Timeout:   x509DefaultDialTimeout,
				KeepAlive: x509DefaultDialTimeout,
			}).DialContext
		}
	}
	transport := &http.Transport{
		DialContext:            dialContext,
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
	if transport.TLSHandshakeTimeout == 0 {
		transport.TLSHandshakeTimeout = x509DefaultTLSHandshakeTimeout
	}
	if transport.ResponseHeaderTimeout == 0 {
		transport.ResponseHeaderTimeout = x509DefaultResponseHeaderTimeout
	}
	if template.Protocols != nil {
		protocols := *template.Protocols
		transport.Protocols = &protocols
	} else if template.TLSNextProto != nil {
		protocols := new(http.Protocols)
		protocols.SetHTTP1(true)
		protocols.SetHTTP2(template.TLSNextProto["h2"] != nil)
		transport.Protocols = protocols
	}
	if template.HTTP2 != nil {
		configuration := *template.HTTP2
		transport.HTTP2 = &configuration
	}
	http1, http2 := x509EffectiveHTTPProtocols(template)
	verifyPeerCertificate := config.VerifyPeerCertificate
	if verifyPeerCertificate != nil {
		config.VerifyPeerCertificate = func(certificates [][]byte, chains [][]*x509.Certificate) error {
			if err := verifyPeerCertificate(certificates, chains); err != nil {
				return errX509TLSVerification
			}
			return nil
		}
	}
	verifyConnection := config.VerifyConnection
	config.VerifyConnection = func(state tls.ConnectionState) error {
		switch state.NegotiatedProtocol {
		case "h2":
			if !http2 || transport.TLSNextProto["h2"] == nil {
				return errX509TLSVerification
			}
		case "", "http/1.1":
			if !http1 {
				return errX509TLSVerification
			}
		default:
			return errX509TLSVerification
		}
		if verifyConnection != nil {
			if err := verifyConnection(state); err != nil {
				return errX509TLSVerification
			}
		}
		return nil
	}

	lifecycle := &x509TransportLifecycle{
		connections:   make(map[*x509TrackedConnection]struct{}),
		dialSemaphore: make(chan struct{}, x509MaximumConcurrentDials),
		dialContext:   transport.DialContext,
		transport:     transport,
	}
	capability := &X509Transport{
		template:          template,
		templateTLS:       template.TLSClientConfig,
		transport:         transport,
		tlsConfig:         config,
		certificateDigest: x509CertificateDigest(config.Certificates[0].Certificate),
		certificateSelect: reflect.ValueOf(config.GetClientCertificate).Pointer(),
		client: &http.Client{
			Transport: transport,
			CheckRedirect: func(*http.Request, []*http.Request) error {
				return errX509Redirect
			},
		},
		lifecycle: lifecycle,
	}
	transport.DialContext = lifecycle.dialTrackedConnection
	return capability, nil
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
	if err := validateX509TLSConfig(transport.TLSClientConfig); err != nil {
		return err
	}
	if transport.DialTLSContext != nil {
		return errors.New("X.509 transport does not support custom TLS dialers")
	}
	//nolint:staticcheck // Deprecated DialTLS still bypasses TLSClientConfig and must be rejected.
	if transport.DialTLS != nil {
		return errors.New("X.509 transport does not support deprecated TLS dialers")
	}
	if transport.Proxy != nil {
		return errors.New("X.509 transport currently supports direct connections only")
	}
	if transport.HTTP2 != nil && transport.HTTP2.CountError != nil {
		return errors.New("X.509 transport does not support HTTP/2 error callbacks")
	}
	return nil
}

func validateX509TLSConfig(config *tls.Config) error {
	if config.GetClientCertificate != nil {
		return errors.New("X.509 transport does not support dynamic client-certificate callbacks")
	}
	return validateX509StaticTLSConfig(config)
}

func validateX509StaticTLSConfig(config *tls.Config) error {
	var signer crypto.Signer
	if len(config.Certificates) == 1 {
		signer, _ = config.Certificates[0].PrivateKey.(crypto.Signer)
	}
	if len(config.Certificates) != 1 || len(config.Certificates[0].Certificate) == 0 ||
		len(config.Certificates[0].Certificate[0]) == 0 || x509PrivateKeyIsNil(signer) {
		return errors.New("X.509 transport requires exactly one static certificate and private key")
	}
	if len(config.EncryptedClientHelloConfigList) != 0 || config.EncryptedClientHelloRejectionVerify != nil {
		return errors.New("X.509 transport does not support encrypted client hello")
	}
	if config.KeyLogWriter != nil {
		return errors.New("X.509 transport does not support TLS session key logging")
	}
	if config.Rand != nil {
		return errors.New("X.509 transport does not support custom TLS randomness")
	}
	if config.Time != nil {
		return errors.New("X.509 transport does not support a custom TLS verification clock")
	}
	if (config.MinVersion != 0 && config.MinVersion < tls.VersionTLS12) ||
		(config.MaxVersion != 0 && config.MaxVersion < tls.VersionTLS12) {
		return errors.New("X.509 transport requires TLS version 1.2 or newer")
	}
	if config.ClientSessionCache != nil {
		return errors.New("X.509 transport does not support shared TLS session caches")
	}
	if config.InsecureSkipVerify {
		return errors.New("X.509 transport requires TLS certificate and hostname verification")
	}
	if config.ServerName != "" {
		return errors.New("X.509 transport does not support overriding the TLS server name")
	}
	return nil
}

func validateX509ApplicationProtocols(template *http.Transport) error {
	http1, http2 := x509EffectiveHTTPProtocols(template)
	if !http1 && !http2 {
		return errors.New("X.509 transport requires an enabled HTTPS-compatible HTTP protocol")
	}
	for _, protocol := range template.TLSClientConfig.NextProtos {
		if protocol != "h2" && protocol != "http/1.1" {
			return errors.New("X.509 transport does not support non-HTTP TLS application protocols")
		}
		if (protocol == "h2" && !http2) || (protocol == "http/1.1" && !http1) {
			return errors.New("X.509 transport TLS application protocol is disabled by its HTTP configuration")
		}
	}
	return nil
}

func x509EffectiveHTTPProtocols(template *http.Transport) (bool, bool) {
	http1, http2 := true, template.ForceAttemptHTTP2
	if template.Protocols != nil {
		http1, http2 = template.Protocols.HTTP1(), template.Protocols.HTTP2()
	} else if template.TLSNextProto != nil {
		http2 = template.TLSNextProto["h2"] != nil
	}
	return http1, http2
}

func x509PrivateKeyIsNil(key any) bool {
	if key == nil {
		return true
	}
	value := reflect.ValueOf(key)
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice, reflect.UnsafePointer:
		return value.IsNil()
	default:
		return false
	}
}

func validateX509TLSProtocolHandlers(handlers map[string]func(string, *tls.Conn) http.RoundTripper) error {
	for protocol, handler := range handlers {
		if (protocol != "h2" && protocol != "unencrypted_http2") || handler == nil {
			return errors.New("X.509 transport does not support custom TLS protocol handlers")
		}
	}
	if len(handlers) != 0 && handlers["h2"] == nil {
		return errors.New("X.509 transport does not support custom TLS protocol handlers")
	}
	return nil
}

func (transport *X509Transport) validateAttestation() error {
	if transport == nil || transport.client == nil || transport.transport == nil || transport.template == nil ||
		transport.lifecycle == nil {
		return errors.New("X.509 transport capability is invalid")
	}
	if transport.lifecycle.closed.Load() {
		return errors.New("X.509 transport capability is closed")
	}
	if err := validateX509NativeTransport(transport.template); err != nil {
		return fmt.Errorf("X.509 transport attestation changed: %w", err)
	}
	if transport.template.TLSClientConfig != transport.templateTLS ||
		transport.transport.TLSClientConfig != transport.tlsConfig {
		return errors.New("X.509 transport TLS configuration changed after attestation")
	}
	if err := validateX509StaticTLSConfig(transport.tlsConfig); err != nil {
		return fmt.Errorf("X.509 transport TLS configuration changed after attestation: %w", err)
	}
	selector := transport.tlsConfig.GetClientCertificate
	if selector == nil || reflect.ValueOf(selector).Pointer() != transport.certificateSelect {
		return errors.New("X.509 transport client-certificate selection changed after attestation")
	}
	selected, err := selector(new(tls.CertificateRequestInfo))
	if err != nil || selected != &transport.tlsConfig.Certificates[0] {
		return errors.New("X.509 transport client-certificate selection changed after attestation")
	}
	for _, config := range []*tls.Config{transport.templateTLS, transport.tlsConfig} {
		digest := x509CertificateDigest(config.Certificates[0].Certificate)
		if subtle.ConstantTimeCompare(digest[:], transport.certificateDigest[:]) != 1 {
			return errors.New("X.509 transport certificate changed after attestation")
		}
	}
	return nil
}

// Close releases this capability's isolated idle connections. It is
// idempotent, never closes the caller's template, and causes Do calls made
// after Close returns to fail. Close does not cancel or wait for calls already
// in progress; an in-progress request may be sent or complete after Close
// returns, at which point its connection is removed from the isolated pool.
func (transport *X509Transport) Close() error {
	if transport == nil || transport.transport == nil || transport.lifecycle == nil {
		return errors.New("X.509 transport capability is invalid")
	}
	transport.lifecycle.close()
	return nil
}

type x509TrackedConnection struct {
	net.Conn
	lifecycle *x509TransportLifecycle
	once      sync.Once
	err       error
}

func (connection *x509TrackedConnection) Close() error {
	connection.once.Do(func() {
		connection.err = connection.Conn.Close()
		connection.lifecycle.mu.Lock()
		delete(connection.lifecycle.connections, connection)
		connection.lifecycle.mu.Unlock()
	})
	return connection.err
}

func (lifecycle *x509TransportLifecycle) dialTrackedConnection(
	ctx context.Context, network, address string,
) (net.Conn, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	requestDone, _ := ctx.Value(x509DialAdmissionContextKey{}).(<-chan struct{})
	select {
	case <-requestDone:
		return nil, context.Canceled
	default:
	}
	select {
	case lifecycle.dialSemaphore <- struct{}{}:
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-requestDone:
		return nil, context.Canceled
	}
	if err := ctx.Err(); err != nil {
		<-lifecycle.dialSemaphore
		return nil, err
	}
	select {
	case <-requestDone:
		<-lifecycle.dialSemaphore
		return nil, context.Canceled
	default:
	}
	defer func() { <-lifecycle.dialSemaphore }()
	connection, err := lifecycle.dialContext(ctx, network, address)
	if err != nil {
		if !x509ConnectionIsNil(connection) {
			_ = connection.Close()
		}
		return nil, err
	}
	if x509ConnectionIsNil(connection) {
		return nil, errors.New("X.509 transport dialer returned a nil connection")
	}
	tracked := &x509TrackedConnection{Conn: connection, lifecycle: lifecycle}
	lifecycle.mu.Lock()
	if lifecycle.closed.Load() && lifecycle.activeRequests == 0 {
		lifecycle.mu.Unlock()
		_ = connection.Close()
		return nil, errors.New("X.509 transport capability closed before its dial completed")
	}
	lifecycle.connections[tracked] = struct{}{}
	lifecycle.mu.Unlock()
	return tracked, nil
}

func x509ConnectionIsNil(connection net.Conn) bool {
	if connection == nil {
		return true
	}
	value := reflect.ValueOf(connection)
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice,
		reflect.UnsafePointer:
		return value.IsNil()
	default:
		return false
	}
}

func (lifecycle *x509TransportLifecycle) beginRequest() bool {
	lifecycle.mu.Lock()
	defer lifecycle.mu.Unlock()
	if lifecycle.closed.Load() {
		return false
	}
	lifecycle.activeRequests++
	return true
}

func (lifecycle *x509TransportLifecycle) endRequest() {
	lifecycle.mu.Lock()
	lifecycle.activeRequests--
	idle := lifecycle.activeRequests == 0
	lifecycle.mu.Unlock()
	closed := lifecycle.closed.Load()
	if closed {
		lifecycle.transport.CloseIdleConnections()
	}
	if idle && closed {
		lifecycle.closeTrackedConnectionsIfIdle()
	}
}

func (lifecycle *x509TransportLifecycle) close() {
	lifecycle.mu.Lock()
	lifecycle.closed.Store(true)
	lifecycle.mu.Unlock()
	lifecycle.transport.CloseIdleConnections()
	lifecycle.closeTrackedConnectionsIfIdle()
}

func (lifecycle *x509TransportLifecycle) closeTrackedConnectionsIfIdle() {
	lifecycle.mu.Lock()
	if lifecycle.activeRequests != 0 {
		lifecycle.mu.Unlock()
		return
	}
	connections := make([]*x509TrackedConnection, 0, len(lifecycle.connections))
	for connection := range lifecycle.connections {
		connections = append(connections, connection)
		delete(lifecycle.connections, connection)
	}
	lifecycle.mu.Unlock()
	for _, connection := range connections {
		_ = connection.Close()
	}
}

type x509ResponseBody struct {
	io.ReadCloser
	lifecycle *x509TransportLifecycle
	finish    sync.Once
}

func (body *x509ResponseBody) Read(data []byte) (int, error) {
	read, err := body.ReadCloser.Read(data)
	if err != nil {
		body.finishRequest()
	}
	return read, err
}

func (body *x509ResponseBody) Close() error {
	err := body.ReadCloser.Close()
	body.finishRequest()
	return err
}

func (body *x509ResponseBody) finishRequest() {
	body.finish.Do(body.lifecycle.endRequest)
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
	dialAdmissionDone := make(chan struct{})
	defer close(dialAdmissionDone)
	requestContext := context.WithValue(
		request.Context(), x509DialAdmissionContextKey{}, (<-chan struct{})(dialAdmissionDone),
	)
	request = request.Clone(x509TraceFreeContext{requestContext})
	if err := transport.validateAttestation(); err != nil {
		return nil, err
	}
	if err := validateX509Request(request); err != nil {
		return nil, err
	}
	if !transport.lifecycle.beginRequest() {
		return nil, errors.New("X.509 transport capability is closed")
	}
	requestActive := true
	defer func() {
		if requestActive {
			transport.lifecycle.endRequest()
		}
	}()
	bodyTransferred = true
	response, err := transport.client.Do(request)
	if err != nil {
		redacted := &x509TransportError{}
		requestContextErr := request.Context().Err()
		var verificationError *tls.CertificateVerificationError
		var networkError net.Error
		if errors.As(err, &networkError) {
			redacted.timeout = networkError.Timeout()
			if temporary, ok := err.(interface{ Temporary() bool }); ok {
				redacted.temporary = temporary.Temporary()
			}
		}
		switch {
		case errors.Is(err, errX509Redirect):
			redacted.cause = errX509Redirect
		case errors.Is(err, context.Canceled):
			if requestContextErr != nil {
				redacted.cause = requestContextErr
			} else {
				redacted.cause = context.Canceled
			}
		case errors.Is(err, context.DeadlineExceeded):
			redacted.cause = context.DeadlineExceeded
			redacted.retryable = requestContextErr == nil
		case errors.Is(err, errX509TLSVerification) || errors.As(err, &verificationError) ||
			x509PermanentClientCertificateAlert(err):
			redacted.cause = errX509TLSVerification
		}
		if redacted.cause == nil {
			redacted.retryable = true
		}
		return nil, redacted
	}
	if response.Body == nil {
		transport.lifecycle.endRequest()
	} else {
		response.Body = &x509ResponseBody{ReadCloser: response.Body, lifecycle: transport.lifecycle}
	}
	requestActive = false
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

func x509PermanentClientCertificateAlert(err error) bool {
	var operation *net.OpError
	if !errors.As(err, &operation) || operation.Op != "remote error" || operation.Err == nil {
		return false
	}
	switch operation.Err.Error() {
	case "tls: bad certificate", "tls: unsupported certificate", "tls: revoked certificate",
		"tls: expired certificate", "tls: unknown certificate", "tls: unknown certificate authority",
		"tls: bad certificate status response", "tls: bad certificate hash value",
		"tls: certificate unobtainable", "tls: certificate required":
		return true
	default:
		return false
	}
}

func validateX509Request(request *http.Request) error {
	switch request.Method {
	case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch,
		http.MethodDelete, http.MethodHead, http.MethodOptions:
	default:
		return errors.New("X.509 requests require a supported non-tunneling HTTP method")
	}
	if len(request.Trailer) != 0 {
		return errors.New("X.509 requests do not support HTTP trailers")
	}
	hasBody := request.Body != nil && request.Body != http.NoBody
	if request.ContentLength < -1 || (!hasBody && request.ContentLength != 0) {
		return errors.New("X.509 requests require consistent HTTP body framing")
	}
	if len(request.TransferEncoding) != 0 {
		return errors.New("X.509 requests do not support custom HTTP transfer encodings")
	}
	if request.RequestURI != "" {
		return errors.New("X.509 client requests cannot specify an explicit request URI")
	}
	if request.Cancel != nil { //nolint:staticcheck // Rejecting caller-controlled legacy cancellation requires its sole, deprecated accessor.
		return errors.New("X.509 requests require context-based cancellation")
	}
	if request.URL == nil || request.URL.Scheme != "https" || request.URL.User != nil ||
		request.URL.Opaque != "" || request.URL.Fragment != "" || request.URL.RawFragment != "" {
		return errors.New("X.509 requests require an absolute HTTPS URL without credentials or fragments")
	}
	for _, character := range []byte(request.URL.RawQuery) {
		if character < ' ' || character == 0x7f {
			return errors.New("X.509 requests cannot contain HTTP control characters in their query")
		}
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
		if !requestconfig.ValidHTTPHeaderName(name) {
			return errors.New("X.509 requests require valid HTTP header names")
		}
		for _, value := range values {
			if !requestconfig.ValidHTTPHeaderValue(value) {
				return errors.New("X.509 requests require valid HTTP header values")
			}
		}
		normalized := strings.ToLower(strings.ReplaceAll(name, "_", "-"))
		switch normalized {
		case "api-key", "x-api-key", "proxy-authorization", "cookie", "set-cookie", "host", "x-amz-security-token":
			return errors.New("X.509 requests cannot contain alternate credentials, cookies, or authority headers")
		case "transfer-encoding", "content-length", "connection", "upgrade", "trailer", "te",
			"proxy-connection", "keep-alive", "http2-settings":
			return errors.New("X.509 requests do not support custom HTTP framing or protocol upgrades")
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

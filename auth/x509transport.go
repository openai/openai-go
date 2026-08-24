package auth

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"path"
	"strings"
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

// X509Transport represents one explicitly attested, caller-owned mTLS transport
// generation. Its identity is the capability pointer, not a certificate-bound
// OAuth token. Rotating credentials requires a new transport and capability.
//
// The application owns the transport, certificate, private key, trust roots,
// connection pool, and lifecycle. It must not mutate the transport or its TLS
// configuration after attestation. HTTPS proxies are unsupported because
// net/http uses the same client TLS configuration for a proxy and its origin.
type X509Transport struct {
	transport *http.Transport
	tlsConfig *tls.Config
	client    *http.Client
}

// NewX509Transport attests a direct, native HTTP transport configured with one
// static client certificate. The transport remains owned by its caller.
func NewX509Transport(transport *http.Transport) (*X509Transport, error) {
	if err := validateX509NativeTransport(transport); err != nil {
		return nil, err
	}

	return &X509Transport{
		transport: transport,
		tlsConfig: transport.TLSClientConfig,
		client: &http.Client{
			Transport: transport,
			CheckRedirect: func(*http.Request, []*http.Request) error {
				return errX509Redirect
			},
		},
	}, nil
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
	if config.InsecureSkipVerify {
		return errors.New("X.509 transport requires TLS certificate and hostname verification")
	}
	if config.ServerName != "" {
		return errors.New("X.509 transport does not support overriding the TLS server name")
	}
	return nil
}

func (transport *X509Transport) validateAttestation() error {
	if transport == nil || transport.client == nil || transport.transport == nil {
		return errors.New("X.509 transport capability is invalid")
	}
	if err := validateX509NativeTransport(transport.transport); err != nil {
		return fmt.Errorf("X.509 transport attestation changed: %w", err)
	}
	if transport.transport.TLSClientConfig != transport.tlsConfig {
		return errors.New("X.509 transport TLS configuration changed after attestation")
	}
	return nil
}

// Do sends a request to an approved global OpenAI mTLS endpoint without
// following redirects. Request context, transport ownership, and pooling are
// preserved. OAuth exchange and client authentication are configured separately.
func (transport *X509Transport) Do(request *http.Request) (*http.Response, error) {
	if request == nil {
		return nil, errors.New("X.509 transport requires a non-nil HTTP request")
	}
	if err := request.Context().Err(); err != nil {
		return nil, err
	}
	if err := transport.validateAttestation(); err != nil {
		return nil, err
	}
	if err := validateX509Request(request); err != nil {
		return nil, err
	}
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
	return response, nil
}

func validateX509Request(request *http.Request) error {
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
			if host == x509AuthenticationHost || authorizationCount != 1 {
				return errors.New("X.509 requests contain invalid Authorization credentials")
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

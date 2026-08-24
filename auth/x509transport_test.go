package auth_test

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"net/http/httptrace"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/openai/openai-go/v3/auth"
)

const (
	x509TransportIssuer = "mtls.auth.openai.com"
	x509TransportAPI    = "mtls.api.openai.com"
)

func TestNewX509TransportRequiresStaticNativeConfiguration(t *testing.T) {
	fixture := newX509TransportFixture(t)
	tests := []struct {
		name   string
		change func(*http.Transport)
		want   string
	}{
		{
			name:   "missing TLS configuration",
			change: func(transport *http.Transport) { transport.TLSClientConfig = nil },
			want:   "TLS configuration",
		},
		{
			name:   "missing client certificate",
			change: func(transport *http.Transport) { transport.TLSClientConfig.Certificates = nil },
			want:   "exactly one static certificate",
		},
		{
			name: "multiple client certificates",
			change: func(transport *http.Transport) {
				transport.TLSClientConfig.Certificates = append(transport.TLSClientConfig.Certificates, fixture.client)
			},
			want: "exactly one static certificate",
		},
		{
			name:   "empty certificate chain",
			change: func(transport *http.Transport) { transport.TLSClientConfig.Certificates[0].Certificate = nil },
			want:   "exactly one static certificate",
		},
		{
			name: "empty leaf certificate",
			change: func(transport *http.Transport) {
				transport.TLSClientConfig.Certificates[0].Certificate = [][]byte{{}}
			},
			want: "exactly one static certificate",
		},
		{
			name:   "missing private key",
			change: func(transport *http.Transport) { transport.TLSClientConfig.Certificates[0].PrivateKey = nil },
			want:   "exactly one static certificate",
		},
		{
			name: "dynamic certificate callback",
			change: func(transport *http.Transport) {
				transport.TLSClientConfig.GetClientCertificate = func(*tls.CertificateRequestInfo) (*tls.Certificate, error) {
					return &fixture.client, nil
				}
			},
			want: "dynamic client-certificate",
		},
		{
			name: "custom TLS context dialer",
			change: func(transport *http.Transport) {
				transport.DialTLSContext = func(context.Context, string, string) (net.Conn, error) {
					return nil, errors.New("synthetic TLS dialer")
				}
			},
			want: "custom TLS dialers",
		},
		{
			name: "deprecated TLS dialer",
			change: func(transport *http.Transport) {
				//nolint:staticcheck // The regression must exercise the deprecated dialer that bypasses TLSClientConfig.
				transport.DialTLS = func(string, string) (net.Conn, error) {
					return nil, errors.New("synthetic deprecated TLS dialer")
				}
			},
			want: "deprecated TLS dialers",
		},
		{
			name: "shared TLS session cache",
			change: func(transport *http.Transport) {
				transport.TLSClientConfig.ClientSessionCache = tls.NewLRUClientSessionCache(1)
			},
			want: "session caches",
		},
		{
			name: "configured proxy",
			change: func(transport *http.Transport) {
				transport.Proxy = http.ProxyURL(&url.URL{Scheme: "http", Host: "proxy.example.test"})
			},
			want: "direct connections",
		},
		{
			name: "custom TLS protocol handler",
			change: func(transport *http.Transport) {
				transport.TLSNextProto = map[string]func(string, *tls.Conn) http.RoundTripper{
					"attacker": func(string, *tls.Conn) http.RoundTripper { return nil },
				}
			},
			want: "TLS protocol handlers",
		},
		{
			name:   "disabled certificate verification",
			change: func(transport *http.Transport) { transport.TLSClientConfig.InsecureSkipVerify = true },
			want:   "hostname verification",
		},
		{
			name:   "fixed TLS server name",
			change: func(transport *http.Transport) { transport.TLSClientConfig.ServerName = x509TransportIssuer },
			want:   "TLS server name",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			transport := fixture.transport(t, nil)
			test.change(transport)
			capability, err := auth.NewX509Transport(transport)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("NewX509Transport error = %v, want %q", err, test.want)
			}
			if capability != nil {
				t.Error("invalid transport returned a non-nil capability")
			}
		})
	}

	t.Run("nil native transport", func(t *testing.T) {
		if capability, err := auth.NewX509Transport(nil); err == nil || capability != nil {
			t.Fatalf("NewX509Transport(nil) = %v, %v; want nil capability and error", capability, err)
		}
	})
}

func TestX509TransportAttestationDetectsUnsafeMutations(t *testing.T) {
	fixture := newX509TransportFixture(t)
	tests := []struct {
		name   string
		change func(*http.Transport)
	}{
		{
			name:   "replace TLS configuration",
			change: func(transport *http.Transport) { transport.TLSClientConfig = transport.TLSClientConfig.Clone() },
		},
		{
			name: "replace static certificate generation",
			change: func(transport *http.Transport) {
				transport.TLSClientConfig.Certificates = []tls.Certificate{fixture.certificate(t, "rotated-workload", nil, true)}
			},
		},
		{
			name: "replace certificate chain",
			change: func(transport *http.Transport) {
				transport.TLSClientConfig.Certificates[0].Certificate = append(
					transport.TLSClientConfig.Certificates[0].Certificate,
					[]byte("synthetic alternate intermediate"),
				)
			},
		},
		{
			name: "install dynamic certificate callback",
			change: func(transport *http.Transport) {
				transport.TLSClientConfig.GetClientCertificate = func(*tls.CertificateRequestInfo) (*tls.Certificate, error) {
					return &fixture.client, nil
				}
			},
		},
		{
			name: "install TLS session cache",
			change: func(transport *http.Transport) {
				transport.TLSClientConfig.ClientSessionCache = tls.NewLRUClientSessionCache(1)
			},
		},
		{
			name: "install proxy",
			change: func(transport *http.Transport) {
				transport.Proxy = http.ProxyURL(&url.URL{Scheme: "https", Host: "proxy.example.test"})
			},
		},
		{
			name: "install custom TLS dialer",
			change: func(transport *http.Transport) {
				transport.DialTLSContext = func(context.Context, string, string) (net.Conn, error) {
					return nil, errors.New("synthetic TLS dialer")
				}
			},
		},
		{
			name: "install custom TLS protocol handler",
			change: func(transport *http.Transport) {
				transport.TLSNextProto = map[string]func(string, *tls.Conn) http.RoundTripper{
					"attacker": func(string, *tls.Conn) http.RoundTripper { return nil },
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			transport := fixture.transport(t, nil)
			capability := newX509Capability(t, transport)
			test.change(transport)
			request := x509TransportRequest(t, http.MethodGet, "https://"+x509TransportAPI+"/v1/models")
			_ = x509Rejected(t, capability, request)
		})
	}
}

func TestX509TransportRejectsCertificateRotationBeforeMutualTLS(t *testing.T) {
	fixture := newX509TransportFixture(t)
	var requests atomic.Int32
	server := fixture.server(t, http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		requests.Add(1)
		w.WriteHeader(http.StatusOK)
	}))
	template := fixture.transport(t, server)
	capability := newX509Capability(t, template)
	request := x509TransportRequest(t, http.MethodGet, "https://"+x509TransportAPI+"/v1/models")
	request.Header.Set("Authorization", "Bearer initial-workload-token")
	first, err := capability.Do(request)
	if err != nil {
		t.Fatalf("initial hostname-verified mTLS request: %v", err)
	}
	if closeErr := first.Body.Close(); closeErr != nil {
		t.Fatalf("close initial mTLS response: %v", closeErr)
	}

	template.TLSClientConfig.Certificates[0] = fixture.certificate(t, "rotated-workload", nil, true)
	if rotationErr := x509Rejected(t, capability, request); !strings.Contains(rotationErr.Error(), "certificate changed") {
		t.Fatalf("certificate rotation error = %v, want changed-generation refusal", rotationErr)
	}
	if got := requests.Load(); got != 1 {
		t.Errorf("mTLS server received %d requests after unauthorized rotation, want one original request", got)
	}
}

func TestX509TransportRejectsUnsafeOriginsBeforeDial(t *testing.T) {
	fixture := newX509TransportFixture(t)
	tests := []struct {
		name   string
		method string
		target string
		change func(*http.Request)
	}{
		{name: "plaintext API", method: http.MethodGet, target: "http://" + x509TransportAPI + "/v1/models"},
		{name: "unapproved API host", method: http.MethodGet, target: "https://api.openai.com/v1/models"},
		{name: "lookalike host", method: http.MethodGet, target: "https://mtls.api.openai.com.attacker.test/v1/models"},
		{name: "Azure host", method: http.MethodGet, target: "https://resource.openai.azure.com/v1/models"},
		{name: "URL user credentials", method: http.MethodGet, target: "https://user:password@" + x509TransportAPI + "/v1/models"},
		{name: "nondefault API port", method: http.MethodGet, target: "https://" + x509TransportAPI + ":444/v1/models"},
		{name: "URL fragment", method: http.MethodGet, target: "https://" + x509TransportAPI + "/v1/models#private"},
		{name: "outside API version", method: http.MethodGet, target: "https://" + x509TransportAPI + "/v2/models"},
		{name: "API path traversal", method: http.MethodGet, target: "https://" + x509TransportAPI + "/v1/../private"},
		{name: "encoded API traversal", method: http.MethodGet, target: "https://" + x509TransportAPI + "/v1/%2e%2e/private"},
		{name: "wrong exchange method", method: http.MethodGet, target: "https://" + x509TransportIssuer + "/oauth/token"},
		{name: "wrong exchange path", method: http.MethodPost, target: "https://" + x509TransportIssuer + "/oauth/other"},
		{name: "encoded exchange path", method: http.MethodPost, target: "https://" + x509TransportIssuer + "/oauth/%74oken"},
		{name: "exchange query", method: http.MethodPost, target: "https://" + x509TransportIssuer + "/oauth/token?audience=unsafe"},
		{
			name:   "mismatched Host authority",
			method: http.MethodGet,
			target: "https://" + x509TransportAPI + "/v1/models",
			change: func(request *http.Request) { request.Host = "attacker.example.test" },
		},
		{
			name:   "URL opaque authority",
			method: http.MethodGet,
			target: "https://" + x509TransportAPI + "/v1/models",
			change: func(request *http.Request) { request.URL.Opaque = "//attacker.example.test/v1/models" },
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var dials atomic.Int32
			transport := fixture.transport(t, nil)
			transport.DialContext = func(context.Context, string, string) (net.Conn, error) {
				dials.Add(1)
				return nil, errors.New("unexpected network dial")
			}
			capability := newX509Capability(t, transport)
			request := x509TransportRequest(t, test.method, test.target)
			if test.change != nil {
				test.change(request)
			}
			_ = x509Rejected(t, capability, request)
			if count := dials.Load(); count != 0 {
				t.Errorf("unsafe request opened %d network connections", count)
			}
		})
	}
}

func TestX509TransportRejectsUnsafeCredentialAliases(t *testing.T) {
	fixture := newX509TransportFixture(t)
	tests := []struct {
		name   string
		host   string
		header string
		values []string
	}{
		{name: "issuer bearer", host: x509TransportIssuer, header: "Authorization", values: []string{"Bearer synthetic"}},
		{name: "issuer bearer alias", host: x509TransportIssuer, header: "aUtHoRiZaTiOn", values: []string{"Bearer synthetic"}},
		{name: "issuer cookie", host: x509TransportIssuer, header: "Cookie", values: []string{"session=synthetic"}},
		{name: "issuer Set-Cookie", host: x509TransportIssuer, header: "Set-Cookie", values: []string{"session=synthetic"}},
		{name: "issuer API key", host: x509TransportIssuer, header: "Api-Key", values: []string{"synthetic"}},
		{name: "issuer API key underscore", host: x509TransportIssuer, header: "API_KEY", values: []string{"synthetic"}},
		{name: "issuer X API key", host: x509TransportIssuer, header: "X-Api-Key", values: []string{"synthetic"}},
		{name: "issuer X API key underscore", host: x509TransportIssuer, header: "x_api_key", values: []string{"synthetic"}},
		{name: "issuer proxy bearer", host: x509TransportIssuer, header: "Proxy-Authorization", values: []string{"Basic synthetic"}},
		{name: "issuer proxy bearer underscore", host: x509TransportIssuer, header: "proxy_authorization", values: []string{"Basic synthetic"}},
		{name: "issuer AWS credential", host: x509TransportIssuer, header: "X-Amz-Security-Token", values: []string{"synthetic"}},
		{name: "issuer forged Host", host: x509TransportIssuer, header: "Host", values: []string{"attacker.test"}},
		{name: "issuer pseudo authority", host: x509TransportIssuer, header: ":authority", values: []string{"attacker.test"}},
		{name: "issuer organization", host: x509TransportIssuer, header: "OpenAI-Organization", values: []string{"synthetic-org"}},
		{name: "issuer organization alias", host: x509TransportIssuer, header: "openai_organization", values: []string{"synthetic-org"}},
		{name: "issuer project", host: x509TransportIssuer, header: "OpenAI-Project", values: []string{"synthetic-project"}},
		{name: "issuer project alias", host: x509TransportIssuer, header: "openai_project", values: []string{"synthetic-project"}},
		{name: "API API key", host: x509TransportAPI, header: "api_key", values: []string{"synthetic"}},
		{name: "API proxy authorization", host: x509TransportAPI, header: "Proxy_Authorization", values: []string{"synthetic"}},
		{name: "API cookie", host: x509TransportAPI, header: "Cookie", values: []string{"session=synthetic"}},
		{name: "multiple API bearers", host: x509TransportAPI, header: "Authorization", values: []string{"Bearer one", "Bearer two"}},
		{name: "API Basic credential", host: x509TransportAPI, header: "Authorization", values: []string{"Basic synthetic"}},
		{name: "API comma-separated bearer", host: x509TransportAPI, header: "Authorization", values: []string{"Bearer first,Bearer second"}},
		{name: "API whitespace bearer", host: x509TransportAPI, header: "Authorization", values: []string{"Bearer token with spaces"}},
		{name: "API empty bearer", host: x509TransportAPI, header: "Authorization", values: []string{"Bearer "}},
		{name: "API malformed bearer padding", host: x509TransportAPI, header: "Authorization", values: []string{"Bearer token=unsafe"}},
		{name: "API padding-only bearer", host: x509TransportAPI, header: "Authorization", values: []string{"Bearer ==="}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			capability := newX509Capability(t, fixture.transport(t, nil))
			method, requestPath := http.MethodPost, "/oauth/token"
			if test.host == x509TransportAPI {
				method, requestPath = http.MethodGet, "/v1/models"
			}
			request := x509TransportRequest(t, method, "https://"+test.host+requestPath)
			request.Header[test.header] = test.values
			_ = x509Rejected(t, capability, request)
		})
	}

	t.Run("API bearer duplicated across differently cased keys", func(t *testing.T) {
		capability := newX509Capability(t, fixture.transport(t, nil))
		request := x509TransportRequest(t, http.MethodGet, "https://"+x509TransportAPI+"/v1/models")
		request.Header["Authorization"] = []string{"Bearer first-synthetic-token"}
		request.Header["authorization"] = []string{"Bearer second-synthetic-token"}
		_ = x509Rejected(t, capability, request)
	})
}

func TestX509TransportUsesIsolatedNativePoolAndCallerOwnedCredentials(t *testing.T) {
	fixture := newX509TransportFixture(t)
	type requestRecord struct {
		host          string
		path          string
		serverName    string
		peer          string
		authorization string
	}
	var mu sync.Mutex
	var received []requestRecord
	server := fixture.server(t, http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		mu.Lock()
		received = append(received, requestRecord{
			host:          request.Host,
			path:          request.URL.Path,
			serverName:    request.TLS.ServerName,
			peer:          request.TLS.PeerCertificates[0].Subject.CommonName,
			authorization: request.Header.Get("Authorization"),
		})
		mu.Unlock()
		w.WriteHeader(http.StatusOK)
	}))

	transport := fixture.transport(t, server)
	config := transport.TLSClientConfig
	type contextKey struct{}
	var contextValues sync.Map
	dialer := &net.Dialer{Timeout: 5 * time.Second}
	transport.DialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
		contextValues.Store(address, ctx.Value(contextKey{}))
		return dialer.DialContext(ctx, network, server.Listener.Addr().String())
	}
	capability := newX509Capability(t, transport)
	if transport.TLSClientConfig != config || transport.Proxy != nil {
		t.Fatal("creating an X.509 capability mutated its caller-owned transport")
	}
	if next := newX509Capability(t, transport); capability == next {
		t.Error("separate attestations unexpectedly share a capability generation")
	}

	for _, test := range []struct {
		host          string
		method        string
		path          string
		authorization string
	}{
		{host: x509TransportIssuer, method: http.MethodPost, path: "/oauth/token"},
		{host: x509TransportAPI, method: http.MethodGet, path: "/v1/models", authorization: "Bearer synthetic-access-token"},
	} {
		ctx := context.WithValue(t.Context(), contextKey{}, "caller-context-value")
		request := x509TransportRequest(t, test.method, "https://"+test.host+test.path).WithContext(ctx)
		if test.authorization != "" {
			request.Header.Set("Authorization", test.authorization)
		}
		response, err := capability.Do(request)
		if err != nil {
			t.Fatalf("%s request: %v", test.host, err)
		}
		if err := response.Body.Close(); err != nil {
			t.Fatalf("close %s response: %v", test.host, err)
		}
		if got, ok := contextValues.Load(test.host + ":443"); !ok || got != "caller-context-value" {
			t.Errorf("%s dial context value = %v, want preserved caller context", test.host, got)
		}
	}

	mu.Lock()
	defer mu.Unlock()
	if len(received) != 2 {
		t.Fatalf("mTLS server received %d requests, want two", len(received))
	}
	for _, record := range received {
		if record.host != record.serverName {
			t.Errorf("request Host/SNI = %q/%q", record.host, record.serverName)
		}
		if record.peer != "static-workload" {
			t.Errorf("request client certificate = %q, want static-workload", record.peer)
		}
	}
	if received[0].authorization != "" || received[1].authorization != "Bearer synthetic-access-token" {
		t.Errorf("issuer/API authorization = %q/%q", received[0].authorization, received[1].authorization)
	}
}

func TestX509TransportIgnoresHiddenRegisteredHTTPSProtocol(t *testing.T) {
	fixture := newX509TransportFixture(t)
	var delivered atomic.Int32
	server := fixture.server(t, http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		if request.Header.Get("Authorization") == "Bearer safe-workload-token" {
			delivered.Add(1)
		}
		w.WriteHeader(http.StatusOK)
	}))
	template := fixture.transport(t, server)
	var captured atomic.Int32
	template.RegisterProtocol("https", &closureTransport{fn: func(request *http.Request) (*http.Response, error) {
		captured.Add(1)
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader("{}")),
			Request:    request,
		}, nil
	}})
	capability := newX509Capability(t, template)
	request := x509TransportRequest(t, http.MethodGet, "https://"+x509TransportAPI+"/v1/models")
	request.Header.Set("Authorization", "Bearer safe-workload-token")

	response, err := capability.Do(request)
	if err != nil {
		t.Fatalf("dispatch through clean native adapter: %v", err)
	}
	if closeErr := response.Body.Close(); closeErr != nil {
		t.Fatalf("close clean-adapter response: %v", closeErr)
	}
	if got := captured.Load(); got != 0 {
		t.Fatalf("caller template's hidden HTTPS handler captured %d protected requests", got)
	}
	if got := delivered.Load(); got != 1 {
		t.Fatalf("real mutually authenticated server received %d protected requests", got)
	}

	originalRequest := x509TransportRequest(t, http.MethodGet, "https://"+x509TransportAPI+"/v1/original")
	originalResponse, err := (&http.Client{Transport: template}).Do(originalRequest)
	if err != nil {
		t.Fatalf("caller-owned original transport: %v", err)
	}
	if err := originalResponse.Body.Close(); err != nil {
		t.Fatalf("close original transport response: %v", err)
	}
	if got := captured.Load(); got != 1 {
		t.Errorf("original caller-owned HTTPS handler was changed or removed; calls = %d", got)
	}
}

func TestX509TransportSnapshotsRequestBeforeTraceMutation(t *testing.T) {
	fixture := newX509TransportFixture(t)
	type observedRequest struct {
		host          string
		authorization string
		serverName    string
	}
	observed := make(chan observedRequest, 1)
	server := fixture.server(t, http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		observed <- observedRequest{
			host:          request.Host,
			authorization: request.Header.Get("Authorization"),
			serverName:    request.TLS.ServerName,
		}
		w.WriteHeader(http.StatusOK)
	}))
	capability := newX509Capability(t, fixture.transport(t, server))
	request := x509TransportRequest(t, http.MethodGet, "https://"+x509TransportAPI+"/v1/models")
	request.Header.Set("Authorization", "Bearer original-workload-token")
	trace := &httptrace.ClientTrace{
		GetConn: func(string) {
			request.URL.Host = "attacker.example.test"
			request.Host = "attacker.example.test"
			request.Header.Set("Authorization", "Bearer attacker-replaced-token")
		},
	}
	request = request.WithContext(httptrace.WithClientTrace(request.Context(), trace))

	response, err := capability.Do(request)
	if err != nil {
		t.Fatalf("dispatch with late request-mutating trace: %v", err)
	}
	if err := response.Body.Close(); err != nil {
		t.Fatalf("close trace-mutating response: %v", err)
	}
	select {
	case record := <-observed:
		if record.host != x509TransportAPI || record.serverName != x509TransportAPI {
			t.Errorf("protected Host/SNI = %q/%q", record.host, record.serverName)
		}
		if record.authorization != "Bearer original-workload-token" {
			t.Errorf("protected bearer was changed by late trace: %q", record.authorization)
		}
	default:
		t.Fatal("real mTLS server did not receive the protected request")
	}
	if request.Host != "attacker.example.test" ||
		request.Header.Get("Authorization") != "Bearer attacker-replaced-token" {
		t.Fatal("negative-control trace did not mutate the original caller request")
	}
}

func TestX509TransportCloseOwnsOnlyItsIsolatedPool(t *testing.T) {
	fixture := newX509TransportFixture(t)
	var mu sync.Mutex
	var connections []string
	server := fixture.server(t, http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		mu.Lock()
		connections = append(connections, request.RemoteAddr)
		mu.Unlock()
		w.WriteHeader(http.StatusOK)
	}))
	template := fixture.transport(t, server)
	caller := &http.Client{Transport: template}
	request := x509TransportRequest(t, http.MethodGet, "https://"+x509TransportAPI+"/v1/models")
	original, err := caller.Do(request)
	if err != nil {
		t.Fatalf("initial caller-owned request: %v", err)
	}
	if closeErr := original.Body.Close(); closeErr != nil {
		t.Fatalf("close caller-owned response: %v", closeErr)
	}

	capability := newX509Capability(t, template)
	protected, err := capability.Do(request)
	if err != nil {
		t.Fatalf("capability-owned request: %v", err)
	}
	if closeErr := protected.Body.Close(); closeErr != nil {
		t.Fatalf("close capability-owned response: %v", closeErr)
	}
	if closeErr := capability.Close(); closeErr != nil {
		t.Fatalf("close capability-owned pool: %v", closeErr)
	}
	if closeErr := capability.Close(); closeErr != nil {
		t.Fatalf("idempotent capability close: %v", closeErr)
	}
	if rejectedErr := x509Rejected(t, capability, request); !strings.Contains(rejectedErr.Error(), "closed") {
		t.Fatalf("closed capability error = %v", rejectedErr)
	}

	remaining, err := caller.Do(request)
	if err != nil {
		t.Fatalf("caller-owned request after capability close: %v", err)
	}
	if err := remaining.Body.Close(); err != nil {
		t.Fatalf("close remaining caller-owned response: %v", err)
	}
	mu.Lock()
	defer mu.Unlock()
	if len(connections) != 3 || connections[0] != connections[2] || connections[0] == connections[1] {
		t.Errorf("original/capability/original connection pools = %v; want original reuse around an isolated pool", connections)
	}
}

func TestX509TransportRejectsRedirectsWithConcreteError(t *testing.T) {
	for _, source := range []string{x509TransportIssuer, x509TransportAPI} {
		t.Run(source, func(t *testing.T) {
			fixture := newX509TransportFixture(t)
			var requests atomic.Int32
			server := fixture.server(t, http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
				requests.Add(1)
				target := "https://" + x509TransportAPI + "/v1/capture"
				if source == x509TransportAPI {
					target = "https://" + x509TransportIssuer + "/oauth/token"
				}
				http.Redirect(w, request, target, http.StatusTemporaryRedirect)
			}))
			capability := newX509Capability(t, fixture.transport(t, server))
			method, requestPath := http.MethodPost, "/oauth/token"
			if source == x509TransportAPI {
				method, requestPath = http.MethodGet, "/v1/models"
			}
			request := x509TransportRequest(t, method, "https://"+source+requestPath)
			if source == x509TransportAPI {
				request.Header.Set("Authorization", "Bearer synthetic-access-token")
			}

			err := x509Rejected(t, capability, request)
			refused := errors.Unwrap(err)
			if refused == nil || !strings.Contains(refused.Error(), "does not follow redirects") ||
				!errors.Is(err, refused) {
				t.Fatalf("redirect error = %v, want preserved concrete refusal error", err)
			}
			var unsafeURL *url.Error
			if errors.As(err, &unsafeURL) {
				t.Error("redirect error retained an unsafe native URL error")
			}
			if errors.Is(err, http.ErrUseLastResponse) {
				t.Error("redirect refusal incorrectly used http.ErrUseLastResponse")
			}
			if count := requests.Load(); count != 1 {
				t.Errorf("redirect chain received %d requests, want only the original", count)
			}
		})
	}
}

func TestX509TransportRedactsNativeTransportFailures(t *testing.T) {
	fixture := newX509TransportFixture(t)
	transport := fixture.transport(t, nil)
	sensitiveCause := errors.New("Authorization: Bearer synthetic-private-transport-token")
	transport.DialContext = func(context.Context, string, string) (net.Conn, error) {
		return nil, sensitiveCause
	}
	capability := newX509Capability(t, transport)
	request := x509TransportRequest(
		t,
		http.MethodGet,
		"https://"+x509TransportAPI+"/v1/models?signature=synthetic-private-url-token",
	)

	err := x509Rejected(t, capability, request)
	if errors.Is(err, sensitiveCause) {
		t.Fatal("native transport error retained its sensitive underlying cause")
	}
	for cause := err; cause != nil; cause = errors.Unwrap(cause) {
		if strings.Contains(cause.Error(), "synthetic-private-transport-token") ||
			strings.Contains(cause.Error(), "synthetic-private-url-token") || strings.Contains(cause.Error(), "Authorization") {
			t.Errorf("transport error chain disclosed sensitive native details: %q", cause.Error())
		}
	}
	var unsafeURL *url.Error
	if errors.As(err, &unsafeURL) {
		t.Error("redacted transport error retained an unsafe native URL error")
	}
}

func TestX509TransportRejectsInvalidCapabilitiesAndCanceledRequests(t *testing.T) {
	fixture := newX509TransportFixture(t)
	request := x509TransportRequest(t, http.MethodGet, "https://"+x509TransportAPI+"/v1/models")

	var nilCapability *auth.X509Transport
	_ = x509Rejected(t, nilCapability, request)
	var forged auth.X509Transport
	_ = x509Rejected(t, &forged, request)
	capability := newX509Capability(t, fixture.transport(t, nil))
	_ = x509Rejected(t, capability, nil)
	withoutURL := request.Clone(request.Context())
	withoutURL.URL = nil
	_ = x509Rejected(t, capability, withoutURL)
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	if err := x509Rejected(t, capability, request.WithContext(ctx)); !errors.Is(err, context.Canceled) {
		t.Fatalf("canceled request error = %v, want context.Canceled", err)
	}
}

func x509Rejected(t *testing.T, capability *auth.X509Transport, request *http.Request) error {
	t.Helper()
	response, err := capability.Do(request)
	if response != nil {
		if closeErr := response.Body.Close(); closeErr != nil {
			t.Fatalf("close unexpected X.509 response: %v", closeErr)
		}
		t.Fatalf("rejected X.509 request returned an unexpected response: %v", response.Status)
	}
	if err == nil {
		t.Fatal("unsafe X.509 request unexpectedly succeeded")
	}
	return err
}

func newX509Capability(t *testing.T, transport *http.Transport) *auth.X509Transport {
	t.Helper()
	capability, err := auth.NewX509Transport(transport)
	if err != nil {
		t.Fatalf("NewX509Transport: %v", err)
	}
	t.Cleanup(func() {
		if err := capability.Close(); err != nil {
			t.Errorf("close isolated X.509 transport: %v", err)
		}
	})
	return capability
}

func x509TransportRequest(t *testing.T, method, target string) *http.Request {
	t.Helper()
	request, err := http.NewRequestWithContext(t.Context(), method, target, nil)
	if err != nil {
		t.Fatalf("create %s %s: %v", method, target, err)
	}
	return request
}

type x509TransportFixture struct {
	root    *x509.Certificate
	rootKey *ecdsa.PrivateKey
	trust   *x509.CertPool
	client  tls.Certificate
}

func newX509TransportFixture(t *testing.T) *x509TransportFixture {
	t.Helper()
	rootKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate ephemeral root key: %v", err)
	}
	template := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "synthetic X.509 transport root"},
		NotBefore:             time.Now().Add(-time.Minute),
		NotAfter:              time.Now().Add(time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	encoded, err := x509.CreateCertificate(rand.Reader, template, template, &rootKey.PublicKey, rootKey)
	if err != nil {
		t.Fatalf("issue ephemeral root certificate: %v", err)
	}
	root, err := x509.ParseCertificate(encoded)
	if err != nil {
		t.Fatalf("parse ephemeral root certificate: %v", err)
	}
	fixture := &x509TransportFixture{root: root, rootKey: rootKey, trust: x509.NewCertPool()}
	fixture.trust.AddCert(root)
	fixture.client = fixture.certificate(t, "static-workload", nil, true)
	return fixture
}

func (fixture *x509TransportFixture) certificate(t *testing.T, name string, hostnames []string, client bool) tls.Certificate {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate ephemeral %q key: %v", name, err)
	}
	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		t.Fatalf("generate ephemeral %q serial: %v", name, err)
	}
	publicKey, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		t.Fatalf("encode ephemeral %q public key: %v", name, err)
	}
	keyID := sha256.Sum256(publicKey)
	template := &x509.Certificate{
		SerialNumber:   serial,
		Subject:        pkix.Name{CommonName: name},
		NotBefore:      time.Now().Add(-time.Minute),
		NotAfter:       time.Now().Add(time.Hour),
		KeyUsage:       x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		DNSNames:       hostnames,
		ExtKeyUsage:    []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		SubjectKeyId:   keyID[:],
		AuthorityKeyId: fixture.root.SubjectKeyId,
	}
	if client {
		template.DNSNames = []string{name + ".workload.example.test"}
		template.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}
	}
	encoded, err := x509.CreateCertificate(rand.Reader, template, fixture.root, &key.PublicKey, fixture.rootKey)
	if err != nil {
		t.Fatalf("issue ephemeral %q certificate: %v", name, err)
	}
	certificate, err := x509.ParseCertificate(encoded)
	if err != nil {
		t.Fatalf("parse ephemeral %q certificate: %v", name, err)
	}
	if client && (len(certificate.DNSNames) == 0 || len(certificate.SubjectKeyId) == 0 ||
		!bytes.Equal(certificate.AuthorityKeyId, fixture.root.SubjectKeyId) ||
		certificate.KeyUsage&(x509.KeyUsageDigitalSignature|x509.KeyUsageKeyEncipherment) !=
			(x509.KeyUsageDigitalSignature|x509.KeyUsageKeyEncipherment) ||
		len(certificate.ExtKeyUsage) != 1 || certificate.ExtKeyUsage[0] != x509.ExtKeyUsageClientAuth) {
		t.Fatalf("ephemeral workload certificate %q does not satisfy X.509 enrollment requirements", name)
	}
	return tls.Certificate{Certificate: [][]byte{encoded}, PrivateKey: key, Leaf: certificate}
}

func (fixture *x509TransportFixture) transport(t *testing.T, server *httptest.Server) *http.Transport {
	t.Helper()
	certificate := fixture.client
	certificate.Certificate = append([][]byte(nil), certificate.Certificate...)
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion:   tls.VersionTLS12,
			RootCAs:      fixture.trust,
			Certificates: []tls.Certificate{certificate},
		},
	}
	if server != nil {
		dialer := &net.Dialer{Timeout: 5 * time.Second}
		transport.DialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
			if address != x509TransportIssuer+":443" && address != x509TransportAPI+":443" {
				return nil, fmt.Errorf("unexpected hermetic X.509 address %q", address)
			}
			return dialer.DialContext(ctx, network, server.Listener.Addr().String())
		}
	}
	t.Cleanup(transport.CloseIdleConnections)
	return transport
}

func (fixture *x509TransportFixture) server(t *testing.T, handler http.Handler) *httptest.Server {
	t.Helper()
	server := httptest.NewUnstartedServer(handler)
	server.TLS = &tls.Config{
		MinVersion:   tls.VersionTLS12,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    fixture.trust,
		Certificates: []tls.Certificate{fixture.certificate(t, "synthetic server", []string{x509TransportIssuer, x509TransportAPI}, false)},
	}
	server.Config.ErrorLog = log.New(io.Discard, "", 0)
	server.StartTLS()
	t.Cleanup(server.Close)
	return server
}

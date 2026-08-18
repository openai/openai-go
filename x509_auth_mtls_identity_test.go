package openai_test

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"io"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type x509ClientIdentityContextKey struct{}

type x509RoundTripperDecorator struct {
	inner http.RoundTripper
}

func (d *x509RoundTripperDecorator) RoundTrip(req *http.Request) (*http.Response, error) {
	return d.inner.RoundTrip(req)
}

func TestClientX509RejectsContextSelectedCertificatesBeforeExchange(t *testing.T) {
	pki := newX509IdentityTestPKI(t)
	server, requests := newX509IdentityTestServer(t, pki)

	transport := &http.Transport{
		Proxy: nil,
		DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
			var dialer net.Dialer
			return dialer.DialContext(ctx, "tcp", server.Listener.Addr().String())
		},
		TLSClientConfig: &tls.Config{
			RootCAs:    pki.roots,
			MinVersion: tls.VersionTLS12,
			GetClientCertificate: func(info *tls.CertificateRequestInfo) (*tls.Certificate, error) {
				identity, _ := info.Context().Value(x509ClientIdentityContextKey{}).(string)
				certificate := pki.clients[identity]
				return &certificate, nil
			},
		},
		DisableKeepAlives: true,
	}
	t.Cleanup(transport.CloseIdleConnections)
	httpClient := &http.Client{
		Transport: transport,
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	tenantAContext := context.WithValue(t.Context(), x509ClientIdentityContextKey{}, "tenant-a")
	tokenRequest, err := http.NewRequestWithContext(tenantAContext, http.MethodPost, "https://mtls.auth.openai.com/oauth/token", nil)
	if err != nil {
		t.Fatalf("http.NewRequestWithContext(token exchange) error = %v", err)
	}
	tokenResponse, err := httpClient.Do(tokenRequest)
	if err != nil {
		t.Fatalf("raw token exchange error = %v", err)
	}
	var token struct {
		AccessToken string `json:"access_token"`
	}
	if decodeErr := json.NewDecoder(tokenResponse.Body).Decode(&token); decodeErr != nil {
		_ = tokenResponse.Body.Close()
		t.Fatalf("decode raw token exchange response error = %v", decodeErr)
	}
	if closeErr := tokenResponse.Body.Close(); closeErr != nil {
		t.Fatalf("close raw token exchange response error = %v", closeErr)
	}

	tenantBContext := context.WithValue(t.Context(), x509ClientIdentityContextKey{}, "tenant-b")
	apiRequest, err := http.NewRequestWithContext(tenantBContext, http.MethodGet, "https://mtls.api.openai.com/v1/models", nil)
	if err != nil {
		t.Fatalf("http.NewRequestWithContext(API) error = %v", err)
	}
	apiRequest.Header.Set("Authorization", "Bearer "+token.AccessToken)
	apiResponse, err := httpClient.Do(apiRequest)
	if err != nil {
		t.Fatalf("raw API request error = %v", err)
	}
	_, _ = io.Copy(io.Discard, apiResponse.Body)
	if err := apiResponse.Body.Close(); err != nil {
		t.Fatalf("close raw API response error = %v", err)
	}
	if got, want := apiResponse.StatusCode, http.StatusUnauthorized; got != want {
		t.Fatalf("raw API status = %d, want %d", got, want)
	}
	if got, want := requests.mismatches.Load(), int32(1); got != want {
		t.Fatalf("raw certificate/bearer mismatches = %d, want %d", got, want)
	}
	requests.reset()

	client := openai.NewClient(
		option.WithHTTPClient(httpClient),
		option.WithX509WorkloadIdentity(clientX509WorkloadIdentity()),
	)

	if _, err := client.Models.List(tenantAContext); err == nil {
		t.Error("Models.List() error = nil")
	}
	if got := requests.total.Load(); got != 0 {
		t.Errorf("mTLS requests = %d, want 0", got)
	}
}

func TestClientX509UsesOneStaticCertificateForExchangeAndAPI(t *testing.T) {
	pki := newX509IdentityTestPKI(t)
	server, requests := newX509IdentityTestServer(t, pki)
	transport := &http.Transport{
		Proxy: nil,
		DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
			var dialer net.Dialer
			return dialer.DialContext(ctx, "tcp", server.Listener.Addr().String())
		},
		TLSClientConfig: &tls.Config{
			Certificates: []tls.Certificate{pki.clients["tenant-a"]},
			RootCAs:      pki.roots,
			MinVersion:   tls.VersionTLS12,
		},
		DisableKeepAlives: true,
	}
	t.Cleanup(transport.CloseIdleConnections)
	client := openai.NewClient(
		option.WithHTTPClient(&http.Client{
			Transport: transport,
			CheckRedirect: func(*http.Request, []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}),
		option.WithX509WorkloadIdentity(clientX509WorkloadIdentity()),
	)

	for range 2 {
		if _, err := client.Models.List(t.Context()); err != nil {
			t.Fatalf("Models.List() error = %v", err)
		}
	}
	if got, want := requests.exchanges.Load(), int32(1); got != want {
		t.Errorf("token exchanges = %d, want %d", got, want)
	}
	if got, want := requests.api.Load(), int32(2); got != want {
		t.Errorf("API requests = %d, want %d", got, want)
	}
	if got := requests.mismatches.Load(); got != 0 {
		t.Errorf("certificate/bearer mismatches = %d, want 0", got)
	}
}

func TestClientX509RejectsStaticCertificateRotationBeforeRealMTLS(t *testing.T) {
	pki := newX509IdentityTestPKI(t)
	server, requests := newX509IdentityTestServer(t, pki)
	transport := &http.Transport{
		Proxy: nil,
		DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
			var dialer net.Dialer
			return dialer.DialContext(ctx, "tcp", server.Listener.Addr().String())
		},
		TLSClientConfig: &tls.Config{
			Certificates: []tls.Certificate{pki.clients["tenant-a"]},
			RootCAs:      pki.roots,
			MinVersion:   tls.VersionTLS12,
		},
		DisableKeepAlives: true,
	}
	t.Cleanup(transport.CloseIdleConnections)
	client := openai.NewClient(
		option.WithHTTPClient(&http.Client{Transport: transport}),
		option.WithX509WorkloadIdentity(clientX509WorkloadIdentity()),
	)

	if _, err := client.Models.List(t.Context()); err != nil {
		t.Fatalf("first Models.List() error = %v", err)
	}
	transport.TLSClientConfig.Certificates[0] = pki.clients["tenant-b"]
	if _, err := client.Models.List(t.Context()); err == nil {
		t.Fatal("Models.List() after certificate rotation error = nil")
	}
	if got, want := requests.total.Load(), int32(2); got != want {
		t.Fatalf("mTLS requests = %d, want %d", got, want)
	}
	if got := requests.mismatches.Load(); got != 0 {
		t.Fatalf("certificate/bearer mismatches = %d, want 0", got)
	}
}

func TestClientX509RejectsOpaqueTransportAndSessionCacheBeforeRealMTLS(t *testing.T) {
	pki := newX509IdentityTestPKI(t)
	server, requests := newX509IdentityTestServer(t, pki)
	newTransport := func() *http.Transport {
		return &http.Transport{
			Proxy: nil,
			DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
				var dialer net.Dialer
				return dialer.DialContext(ctx, "tcp", server.Listener.Addr().String())
			},
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{pki.clients["tenant-a"]},
				RootCAs:      pki.roots,
				MinVersion:   tls.VersionTLS12,
			},
		}
	}

	t.Run("opaque decorator", func(t *testing.T) {
		transport := newTransport()
		t.Cleanup(transport.CloseIdleConnections)
		client := openai.NewClient(
			option.WithHTTPClient(&http.Client{Transport: &x509RoundTripperDecorator{inner: transport}}),
			option.WithX509WorkloadIdentity(clientX509WorkloadIdentity()),
		)
		if _, err := client.Models.List(t.Context()); err == nil {
			t.Fatal("Models.List() error = nil")
		}
		if got := requests.total.Load(); got != 0 {
			t.Fatalf("mTLS requests = %d, want 0", got)
		}
	})

	t.Run("client session cache", func(t *testing.T) {
		transport := newTransport()
		transport.TLSClientConfig.ClientSessionCache = tls.NewLRUClientSessionCache(1)
		t.Cleanup(transport.CloseIdleConnections)
		client := openai.NewClient(
			option.WithHTTPClient(&http.Client{Transport: transport}),
			option.WithX509WorkloadIdentity(clientX509WorkloadIdentity()),
		)
		if _, err := client.Models.List(t.Context()); err == nil {
			t.Fatal("Models.List() error = nil")
		}
		if got := requests.total.Load(); got != 0 {
			t.Fatalf("mTLS requests = %d, want 0", got)
		}
	})
}

type x509IdentityTestRequestCounts struct {
	total      atomic.Int32
	exchanges  atomic.Int32
	api        atomic.Int32
	mismatches atomic.Int32
}

func (r *x509IdentityTestRequestCounts) reset() {
	r.total.Store(0)
	r.exchanges.Store(0)
	r.api.Store(0)
	r.mismatches.Store(0)
}

func newX509IdentityTestServer(t *testing.T, pki x509IdentityTestPKI) (*httptest.Server, *x509IdentityTestRequestCounts) {
	t.Helper()
	requests := &x509IdentityTestRequestCounts{}
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.total.Add(1)
		identity := r.TLS.PeerCertificates[0].Subject.CommonName
		if r.URL.Path == "/oauth/token" {
			requests.exchanges.Add(1)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"access_token": "token-" + identity,
				"expires_in":   3600,
			})
			return
		}
		requests.api.Add(1)
		if got, want := r.Header.Get("Authorization"), "Bearer token-"+identity; got != want {
			requests.mismatches.Add(1)
			http.Error(w, "certificate/bearer mismatch", http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"object":"list","data":[]}`))
	}))
	server.TLS = &tls.Config{
		Certificates: []tls.Certificate{pki.server},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    pki.roots,
		MinVersion:   tls.VersionTLS12,
	}
	server.StartTLS()
	t.Cleanup(server.Close)
	return server, requests
}

type x509IdentityTestPKI struct {
	roots   *x509.CertPool
	server  tls.Certificate
	clients map[string]tls.Certificate
}

func newX509IdentityTestPKI(t *testing.T) x509IdentityTestPKI {
	t.Helper()
	now := time.Now()
	caKey := newX509IdentityTestKey(t)
	caTemplate := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "test-root"},
		NotBefore:             now.Add(-time.Minute),
		NotAfter:              now.Add(time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	caDER, err := x509.CreateCertificate(rand.Reader, caTemplate, caTemplate, &caKey.PublicKey, caKey)
	if err != nil {
		t.Fatalf("x509.CreateCertificate(CA) error = %v", err)
	}
	ca, err := x509.ParseCertificate(caDER)
	if err != nil {
		t.Fatalf("x509.ParseCertificate(CA) error = %v", err)
	}
	roots := x509.NewCertPool()
	roots.AddCert(ca)

	return x509IdentityTestPKI{
		roots: roots,
		server: newX509IdentityTestCertificate(t, ca, caKey, 2, "test-server", []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}, []string{
			"mtls.auth.openai.com",
			"mtls.api.openai.com",
		}),
		clients: map[string]tls.Certificate{
			"tenant-a": newX509IdentityTestCertificate(t, ca, caKey, 3, "tenant-a", []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}, nil),
			"tenant-b": newX509IdentityTestCertificate(t, ca, caKey, 4, "tenant-b", []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}, nil),
		},
	}
}

func newX509IdentityTestKey(t *testing.T) *ecdsa.PrivateKey {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("ecdsa.GenerateKey() error = %v", err)
	}
	return key
}

func newX509IdentityTestCertificate(
	t *testing.T,
	ca *x509.Certificate,
	caKey *ecdsa.PrivateKey,
	serial int64,
	commonName string,
	extKeyUsage []x509.ExtKeyUsage,
	dnsNames []string,
) tls.Certificate {
	t.Helper()
	key := newX509IdentityTestKey(t)
	template := &x509.Certificate{
		SerialNumber: big.NewInt(serial),
		Subject:      pkix.Name{CommonName: commonName},
		DNSNames:     dnsNames,
		NotBefore:    time.Now().Add(-time.Minute),
		NotAfter:     time.Now().Add(time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  extKeyUsage,
	}
	certificateDER, err := x509.CreateCertificate(rand.Reader, template, ca, &key.PublicKey, caKey)
	if err != nil {
		t.Fatalf("x509.CreateCertificate(%q) error = %v", commonName, err)
	}
	return tls.Certificate{
		Certificate: [][]byte{certificateDER, ca.Raw},
		PrivateKey:  key,
	}
}

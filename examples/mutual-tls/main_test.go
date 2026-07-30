package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func TestNativeMutualTLSHTTPClient(t *testing.T) {
	fixture := newTLSFixture(t)
	tempDir := t.TempDir()
	certificatePath := filepath.Join(tempDir, "client-chain.pem")
	privateKeyPath := filepath.Join(tempDir, "client.key")
	if err := os.WriteFile(certificatePath, fixture.certificatePEM, 0o600); err != nil {
		t.Fatalf("os.WriteFile(%q) error = %v", certificatePath, err)
	}
	if err := os.WriteFile(privateKeyPath, fixture.privateKeyPEM, 0o600); err != nil {
		t.Fatalf("os.WriteFile(%q) error = %v", privateKeyPath, err)
	}

	certificate, err := tls.LoadX509KeyPair(certificatePath, privateKeyPath)
	if err != nil {
		t.Fatalf("tls.LoadX509KeyPair(%q, %q) error = %v", certificatePath, privateKeyPath, err)
	}
	if got, want := len(certificate.Certificate), 2; got != want {
		t.Fatalf("tls.LoadX509KeyPair() chain length = %d, want %d", got, want)
	}

	originalDefaultTransport := http.DefaultTransport
	baseTransport, ok := originalDefaultTransport.(*http.Transport)
	if !ok {
		t.Skip("http.DefaultTransport is not an *http.Transport")
	}
	defaultTransport := baseTransport.Clone()
	defaultTLSConfig := &tls.Config{}
	if defaultTransport.TLSClientConfig != nil {
		defaultTLSConfig = defaultTransport.TLSClientConfig.Clone()
	}
	defaultTLSConfig.GetClientCertificate = func(*tls.CertificateRequestInfo) (*tls.Certificate, error) {
		return nil, nil
	}
	defaultTransport.TLSClientConfig = defaultTLSConfig
	defaultTransport.Proxy = http.ProxyFromEnvironment
	http.DefaultTransport = defaultTransport
	t.Cleanup(func() {
		http.DefaultTransport = originalDefaultTransport
	})

	originalCertificateCount := len(defaultTransport.TLSClientConfig.Certificates)
	originalResponseHeaderTimeout := defaultTransport.ResponseHeaderTimeout
	httpClient, err := newMutualTLSHTTPClient(certificate)
	if err != nil {
		t.Fatalf("newMutualTLSHTTPClient() error = %v", err)
	}
	t.Cleanup(httpClient.CloseIdleConnections)

	transport, ok := httpClient.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("newMutualTLSHTTPClient().Transport type = %T, want *http.Transport", httpClient.Transport)
	}
	if transport == defaultTransport {
		t.Fatal("newMutualTLSHTTPClient() reused http.DefaultTransport")
	}
	if transport.Proxy != nil {
		t.Error("newMutualTLSHTTPClient().Transport.Proxy is non-nil")
	}
	if got, want := transport.ResponseHeaderTimeout, responseHeaderTimeout; got != want {
		t.Errorf("newMutualTLSHTTPClient().Transport.ResponseHeaderTimeout = %v, want %v", got, want)
	}
	if transport.TLSClientConfig.GetClientCertificate != nil {
		t.Error("newMutualTLSHTTPClient().Transport.TLSClientConfig.GetClientCertificate is non-nil")
	}
	defaultCertificateCount := len(defaultTransport.TLSClientConfig.Certificates)
	if defaultCertificateCount != originalCertificateCount {
		t.Errorf(
			"newMutualTLSHTTPClient() default certificate count = %d, want %d",
			defaultCertificateCount,
			originalCertificateCount,
		)
	}
	if defaultTransport.Proxy == nil {
		t.Error("newMutualTLSHTTPClient() cleared http.DefaultTransport.Proxy")
	}
	if defaultTransport.ResponseHeaderTimeout != originalResponseHeaderTimeout {
		t.Errorf(
			"newMutualTLSHTTPClient() default response header timeout = %v, want %v",
			defaultTransport.ResponseHeaderTimeout,
			originalResponseHeaderTimeout,
		)
	}
	if defaultTransport.TLSClientConfig.GetClientCertificate == nil {
		t.Error("newMutualTLSHTTPClient() cleared http.DefaultTransport.TLSClientConfig.GetClientCertificate")
	}

	rootPool := x509.NewCertPool()
	rootPool.AddCert(fixture.root)
	transport.TLSClientConfig.RootCAs = rootPool

	peerCertificateCount := 0
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		peerCertificateCount = len(request.TLS.PeerCertificates)
		if got, want := request.Header.Get("Authorization"), "Bearer test-api-key"; got != want {
			t.Errorf("Authorization header = %q, want %q", got, want)
		}
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(`{"ok":true}`)); err != nil {
			t.Errorf("ResponseWriter.Write() error = %v", err)
		}
	}))
	server.TLS = &tls.Config{
		Certificates: []tls.Certificate{fixture.server},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    rootPool,
		MinVersion:   tls.VersionTLS12,
	}
	server.StartTLS()
	t.Cleanup(server.Close)

	client := openai.NewClient(
		option.WithAPIKey("test-api-key"),
		option.WithBaseURL(server.URL+"/v1/"),
		option.WithHTTPClient(httpClient),
		option.WithMaxRetries(0),
	)
	var response struct {
		OK bool `json:"ok"`
	}
	if err := client.Get(t.Context(), "test", nil, &response); err != nil {
		t.Fatalf("Client.Get() error = %v", err)
	}
	if !response.OK {
		t.Error("Client.Get().OK = false, want true")
	}
	if got, want := peerCertificateCount, 2; got != want {
		t.Errorf("Client.Get() peer certificate count = %d, want %d", got, want)
	}

	redirectRequest := &http.Request{}
	if err := httpClient.CheckRedirect(redirectRequest, nil); err != http.ErrUseLastResponse {
		t.Errorf("CheckRedirect() error = %v, want %v", err, http.ErrUseLastResponse)
	}
}

type tlsFixture struct {
	root           *x509.Certificate
	server         tls.Certificate
	certificatePEM []byte
	privateKeyPEM  []byte
}

func newTLSFixture(t *testing.T) tlsFixture {
	t.Helper()
	now := time.Now()

	rootKey := newECDSAKey(t)
	rootTemplate := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "mTLS test root"},
		NotBefore:             now.Add(-time.Hour),
		NotAfter:              now.Add(24 * time.Hour),
		IsCA:                  true,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageCertSign,
	}
	rootDER := createCertificate(t, rootTemplate, rootTemplate, &rootKey.PublicKey, rootKey)
	root := parseCertificate(t, rootDER)

	intermediateKey := newECDSAKey(t)
	intermediateTemplate := &x509.Certificate{
		SerialNumber:          big.NewInt(2),
		Subject:               pkix.Name{CommonName: "mTLS test intermediate"},
		NotBefore:             now.Add(-time.Hour),
		NotAfter:              now.Add(12 * time.Hour),
		IsCA:                  true,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageCertSign,
	}
	intermediateDER := createCertificate(t, intermediateTemplate, root, &intermediateKey.PublicKey, rootKey)
	intermediate := parseCertificate(t, intermediateDER)

	clientKey := newECDSAKey(t)
	clientTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(3),
		Subject:      pkix.Name{CommonName: "mTLS test client"},
		NotBefore:    now.Add(-time.Hour),
		NotAfter:     now.Add(12 * time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}
	clientDER := createCertificate(t, clientTemplate, intermediate, &clientKey.PublicKey, intermediateKey)
	certificatePEM := append(
		pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: clientDER}),
		pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: intermediateDER})...,
	)
	privateKeyDER, err := x509.MarshalECPrivateKey(clientKey)
	if err != nil {
		t.Fatalf("x509.MarshalECPrivateKey() error = %v", err)
	}

	serverKey := newECDSAKey(t)
	serverTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(4),
		Subject:      pkix.Name{CommonName: "mTLS test server"},
		NotBefore:    now.Add(-time.Hour),
		NotAfter:     now.Add(12 * time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:     []string{"localhost"},
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
	}
	serverDER := createCertificate(t, serverTemplate, root, &serverKey.PublicKey, rootKey)

	return tlsFixture{
		root: root,
		server: tls.Certificate{
			Certificate: [][]byte{serverDER},
			PrivateKey:  serverKey,
		},
		certificatePEM: certificatePEM,
		privateKeyPEM: pem.EncodeToMemory(&pem.Block{
			Type:  "EC PRIVATE KEY",
			Bytes: privateKeyDER,
		}),
	}
}

func newECDSAKey(t *testing.T) *ecdsa.PrivateKey {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("ecdsa.GenerateKey() error = %v", err)
	}
	return key
}

func createCertificate(t *testing.T, template, parent *x509.Certificate, publicKey, signer any) []byte {
	t.Helper()
	der, err := x509.CreateCertificate(rand.Reader, template, parent, publicKey, signer)
	if err != nil {
		t.Fatalf("x509.CreateCertificate() error = %v", err)
	}
	return der
}

func parseCertificate(t *testing.T, der []byte) *x509.Certificate {
	t.Helper()
	certificate, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatalf("x509.ParseCertificate() error = %v", err)
	}
	return certificate
}

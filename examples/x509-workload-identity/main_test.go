package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewX509HTTPClientOwnsTransportConfiguration(t *testing.T) {
	certificatePath, privateKeyPath := writeTestKeyPair(t)

	originalDefaultTransport := http.DefaultTransport
	baseTransport, ok := originalDefaultTransport.(*http.Transport)
	if !ok {
		t.Skip("http.DefaultTransport is not an *http.Transport")
	}
	defaultTransport := baseTransport.Clone()
	defaultTransport.Proxy = http.ProxyFromEnvironment
	defaultTransport.ResponseHeaderTimeout = time.Second
	defaultTransport.TLSClientConfig = &tls.Config{ServerName: "inherited.invalid"}
	http.DefaultTransport = defaultTransport
	t.Cleanup(func() { http.DefaultTransport = originalDefaultTransport })

	httpClient, err := newX509HTTPClient(certificatePath, privateKeyPath)
	if err != nil {
		t.Fatalf("newX509HTTPClient() error = %v", err)
	}
	t.Cleanup(httpClient.CloseIdleConnections)

	transport, ok := httpClient.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("newX509HTTPClient().Transport type = %T, want *http.Transport", httpClient.Transport)
	}
	if transport == defaultTransport {
		t.Fatal("newX509HTTPClient() reused http.DefaultTransport")
	}
	if transport.Proxy != nil {
		t.Error("newX509HTTPClient().Transport.Proxy is non-nil")
	}
	if got, want := transport.ResponseHeaderTimeout, responseHeaderTimeout; got != want {
		t.Errorf("ResponseHeaderTimeout = %v, want %v", got, want)
	}
	if got, want := transport.TLSClientConfig.MinVersion, uint16(tls.VersionTLS12); got != want {
		t.Errorf("TLS MinVersion = %d, want %d", got, want)
	}
	if got := transport.TLSClientConfig.ServerName; got != "" {
		t.Errorf("TLS ServerName = %q, want empty", got)
	}
	if got, want := len(transport.TLSClientConfig.Certificates), 1; got != want {
		t.Fatalf("TLS certificate count = %d, want %d", got, want)
	}
	if transport.TLSClientConfig.Certificates[0].PrivateKey == nil {
		t.Error("transport-owned certificate has no private key")
	}
	if transport.TLSClientConfig.GetClientCertificate == nil {
		t.Fatal("GetClientCertificate is nil")
	}
	certificate, err := transport.TLSClientConfig.GetClientCertificate(&tls.CertificateRequestInfo{})
	if err != nil {
		t.Fatalf("GetClientCertificate() error = %v", err)
	}
	if certificate.PrivateKey != transport.TLSClientConfig.Certificates[0].PrivateKey {
		t.Error("GetClientCertificate() returned a different private key")
	}
	if err := httpClient.CheckRedirect(nil, nil); !errors.Is(err, http.ErrUseLastResponse) {
		t.Errorf("CheckRedirect() error = %v, want http.ErrUseLastResponse", err)
	}

	if defaultTransport.Proxy == nil {
		t.Error("newX509HTTPClient() mutated http.DefaultTransport.Proxy")
	}
	if got, want := defaultTransport.ResponseHeaderTimeout, time.Second; got != want {
		t.Errorf("default ResponseHeaderTimeout = %v, want %v", got, want)
	}
	if got, want := defaultTransport.TLSClientConfig.ServerName, "inherited.invalid"; got != want {
		t.Errorf("default TLS ServerName = %q, want %q", got, want)
	}
	if got := len(defaultTransport.TLSClientConfig.Certificates); got != 0 {
		t.Errorf("default TLS certificate count = %d, want 0", got)
	}
}

func writeTestKeyPair(t *testing.T) (string, string) {
	t.Helper()
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("ecdsa.GenerateKey() error = %v", err)
	}
	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		NotBefore:    time.Now().Add(-time.Minute),
		NotAfter:     time.Now().Add(time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}
	certificateDER, err := x509.CreateCertificate(rand.Reader, template, template, &privateKey.PublicKey, privateKey)
	if err != nil {
		t.Fatalf("x509.CreateCertificate() error = %v", err)
	}
	privateKeyDER, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		t.Fatalf("x509.MarshalECPrivateKey() error = %v", err)
	}

	tempDir := t.TempDir()
	certificatePath := filepath.Join(tempDir, "client.pem")
	privateKeyPath := filepath.Join(tempDir, "client.key")
	certificatePEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certificateDER})
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: privateKeyDER})
	if err := os.WriteFile(certificatePath, certificatePEM, 0o600); err != nil {
		t.Fatalf("os.WriteFile(%q) error = %v", certificatePath, err)
	}
	if err := os.WriteFile(privateKeyPath, privateKeyPEM, 0o600); err != nil {
		t.Fatalf("os.WriteFile(%q) error = %v", privateKeyPath, err)
	}
	return certificatePath, privateKeyPath
}

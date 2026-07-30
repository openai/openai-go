// The mutual-tls command demonstrates native Go mutual TLS configuration with
// option.WithHTTPClient.
package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

const responseHeaderTimeout = 10 * time.Minute

func main() {
	if err := run(context.Background()); err != nil {
		slog.Error("mutual TLS example failed", "err", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	// When certificate-chain support is enabled for the organization,
	// client-chain.pem must contain the leaf certificate first, followed by
	// every required intermediate. Otherwise, the leaf must be signed directly
	// by an active uploaded certificate.
	certificate, err := tls.LoadX509KeyPair(
		"/secrets/openai/client-chain.pem",
		"/secrets/openai/client.key",
	)
	if err != nil {
		return fmt.Errorf("load client certificate: %w", err)
	}

	httpClient, err := newMutualTLSHTTPClient(certificate)
	if err != nil {
		return fmt.Errorf("configure mutual TLS HTTP client: %w", err)
	}

	client := openai.NewClient(
		option.WithBaseURL("https://mtls.api.openai.com/v1"),
		option.WithHTTPClient(httpClient),
	)

	models, err := client.Models.List(ctx)
	if err != nil {
		return fmt.Errorf("list models: %w", err)
	}
	slog.Info("mutual TLS request succeeded", "models", len(models.Data))
	return nil
}

func newMutualTLSHTTPClient(certificate tls.Certificate) (*http.Client, error) {
	defaultTransport, ok := http.DefaultTransport.(*http.Transport)
	if !ok {
		return nil, errors.New("http.DefaultTransport is not an *http.Transport")
	}
	transport := defaultTransport.Clone()
	transport.Proxy = nil
	transport.DialTLS = nil
	transport.DialTLSContext = nil
	transport.ResponseHeaderTimeout = responseHeaderTimeout

	transport.TLSClientConfig = &tls.Config{
		Certificates: []tls.Certificate{certificate},
		// Always return this certificate; automatic selection can reject it when
		// the server's acceptable-CA hint does not match the local chain.
		GetClientCertificate: func(*tls.CertificateRequestInfo) (*tls.Certificate, error) {
			return &certificate, nil
		},
	}

	return &http.Client{
		Transport: transport,
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}, nil
}

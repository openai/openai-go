// The x509-workload-identity command demonstrates an application-owned rollout
// toggle between API-key authentication and X.509 workload identity federation.
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
	"github.com/openai/openai-go/v3/auth"
	"github.com/openai/openai-go/v3/option"
)

const responseHeaderTimeout = 10 * time.Minute

func main() {
	if err := run(context.Background()); err != nil {
		slog.Error("OpenAI request failed", "err", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	var client openai.Client
	switch mode := os.Getenv("OPENAI_AUTH_MODE"); mode {
	case "", "api_key":
		client = openai.NewClient()
	case "x509":
		certificateChainPath, err := requiredEnvironmentVariable("OPENAI_MTLS_CERTIFICATE_CHAIN")
		if err != nil {
			return err
		}
		privateKeyPath, err := requiredEnvironmentVariable("OPENAI_MTLS_PRIVATE_KEY")
		if err != nil {
			return err
		}
		identityProviderID, err := requiredEnvironmentVariable("OPENAI_IDENTITY_PROVIDER_ID")
		if err != nil {
			return err
		}
		serviceAccountID, err := requiredEnvironmentVariable("OPENAI_SERVICE_ACCOUNT_ID")
		if err != nil {
			return err
		}

		httpClient, err := newX509HTTPClient(certificateChainPath, privateKeyPath)
		if err != nil {
			return err
		}
		defer httpClient.CloseIdleConnections()

		client = openai.NewClient(
			option.WithHTTPClient(httpClient),
			option.WithX509WorkloadIdentity(auth.X509WorkloadIdentity{
				IdentityProviderID: identityProviderID,
				ServiceAccountID:   serviceAccountID,
			}),
		)
	default:
		return fmt.Errorf("unsupported OPENAI_AUTH_MODE %q (want api_key or x509)", mode)
	}

	models, err := client.Models.List(ctx)
	if err != nil {
		return fmt.Errorf("list models: %w", err)
	}
	slog.Info("OpenAI request succeeded", "models", len(models.Data))
	return nil
}

func requiredEnvironmentVariable(name string) (string, error) {
	value := os.Getenv(name)
	if value == "" {
		return "", fmt.Errorf("%s is required", name)
	}
	return value, nil
}

func newX509HTTPClient(certificateChainPath string, privateKeyPath string) (*http.Client, error) {
	certificate, err := tls.LoadX509KeyPair(certificateChainPath, privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("load X.509 client certificate: %w", err)
	}

	defaultTransport, ok := http.DefaultTransport.(*http.Transport)
	if !ok {
		return nil, errors.New("http.DefaultTransport is not an *http.Transport")
	}
	transport := defaultTransport.Clone()
	transport.Proxy = nil
	transport.DialTLS = nil //nolint:staticcheck // Clear an inherited legacy hook so TLSClientConfig remains authoritative.
	transport.DialTLSContext = nil
	transport.ResponseHeaderTimeout = responseHeaderTimeout
	transport.TLSClientConfig = &tls.Config{
		Certificates: []tls.Certificate{certificate},
		MinVersion:   tls.VersionTLS12,
	}

	return &http.Client{
		Transport: transport,
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}, nil
}

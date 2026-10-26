package azure

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"math/big"
	"net"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/auth"
	"github.com/openai/openai-go/v3/option"
)

func TestX509WorkloadIdentityRejectsAzureBeforeExchange(t *testing.T) {
	for _, test := range []struct {
		name string
		opts []option.RequestOption
	}{
		{
			name: "endpoint without credentials",
			opts: []option.RequestOption{
				WithEndpoint("https://resource.openai.azure.com", "2024-06-01"),
			},
		},
		{
			name: "API key",
			opts: []option.RequestOption{
				WithEndpoint("https://resource.openai.azure.com", "2024-06-01"),
				WithAPIKey("synthetic-azure-key"),
			},
		},
	} {
		for _, order := range []struct {
			name        string
			x509First   bool
			wantMessage string
		}{
			{name: "Azure first"},
			{name: "X.509 first", x509First: true, wantMessage: "cannot be combined"},
		} {
			t.Run(test.name+"/"+order.name, func(t *testing.T) {
				t.Setenv("OPENAI_BASE_URL", "https://mtls.api.openai.com/v1/")
				config, connections := syntheticX509WorkloadIdentity(t)
				x509Opts := []option.RequestOption{
					option.WithX509WorkloadIdentity(config),
					option.WithMaxRetries(0),
				}
				opts := append(test.opts, x509Opts...)
				if order.x509First {
					opts = append(x509Opts, test.opts...)
				}
				client := openai.NewClient(opts...)
				_, err := client.Models.List(t.Context())
				if err == nil || order.wantMessage != "" && !strings.Contains(err.Error(), order.wantMessage) {
					t.Errorf("Models.List() error = %v, want Azure/X.509 conflict", err)
				}
				if got := connections.Load(); got != 0 {
					t.Errorf("Models.List() attempted %d issuer or API connections, want 0", got)
				}
			})
		}
	}
}

func syntheticX509WorkloadIdentity(t *testing.T) (auth.X509WorkloadIdentity, *atomic.Int32) {
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
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}
	certificate, err := x509.CreateCertificate(rand.Reader, template, template, &privateKey.PublicKey, privateKey)
	if err != nil {
		t.Fatalf("x509.CreateCertificate() error = %v", err)
	}
	connections := &atomic.Int32{}
	transport, err := auth.NewX509Transport(&http.Transport{
		DialContext: func(context.Context, string, string) (net.Conn, error) {
			connections.Add(1)
			return nil, errors.New("unexpected synthetic workload connection")
		},
		TLSClientConfig: &tls.Config{
			Certificates: []tls.Certificate{{Certificate: [][]byte{certificate}, PrivateKey: privateKey}},
			MinVersion:   tls.VersionTLS12,
		},
	})
	if err != nil {
		t.Fatalf("auth.NewX509Transport() error = %v", err)
	}
	t.Cleanup(func() {
		if err := transport.Close(); err != nil {
			t.Errorf("X509Transport.Close() error = %v", err)
		}
	})
	return auth.X509WorkloadIdentity{
		IdentityProviderID: "synthetic-identity-provider",
		ServiceAccountID:   "synthetic-service-account",
		Transport:          transport,
	}, connections
}

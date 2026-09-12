package bedrock

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
	"sync/atomic"
	"testing"
	"time"

	"github.com/openai/openai-go/v3/auth"
	"github.com/openai/openai-go/v3/option"
)

func TestSkipAuthRejectsX509WorkloadIdentityBeforeNetwork(t *testing.T) {
	t.Setenv("OPENAI_BASE_URL", "https://mtls.api.openai.com/v1/")
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate synthetic workload private key: %v", err)
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
		t.Fatalf("issue synthetic workload certificate: %v", err)
	}
	var connections atomic.Int32
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
		t.Fatalf("attest synthetic workload transport: %v", err)
	}
	t.Cleanup(func() {
		if closeErr := transport.Close(); closeErr != nil {
			t.Errorf("close synthetic workload transport: %v", closeErr)
		}
	})

	client, err := NewClient(t.Context(), Config{AWSRegion: "us-east-1", SkipAuth: true},
		option.WithX509WorkloadIdentity(auth.X509WorkloadIdentity{
			IdentityProviderID: "synthetic-identity-provider",
			ServiceAccountID:   "synthetic-service-account",
			Transport:          transport,
		}),
		option.WithMaxRetries(0),
	)
	if err == nil {
		_, err = client.Models.List(t.Context())
	}
	if err == nil {
		t.Fatal("Bedrock SkipAuth accepted an X.509 workload identity")
	}
	if got := connections.Load(); got != 0 {
		t.Errorf("rejected workload identity attempted %d issuer or API connections", got)
	}
}

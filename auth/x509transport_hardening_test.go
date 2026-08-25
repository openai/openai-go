package auth_test

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/openai/openai-go/v3/auth"
)

func TestX509TransportPresentsStaticCertificateDespiteMismatchedCAHints(t *testing.T) {
	for _, host := range []string{x509TransportIssuer, x509TransportAPI} {
		t.Run(host, func(t *testing.T) {
			fixture := newX509TransportFixture(t)
			advertised := x509.NewCertPool()
			advertised.AddCert(fixture.certificate(t, "synthetic unrelated CA hint", nil, false).Leaf)
			var presented atomic.Int32
			server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
				presented.Store(int32(len(request.TLS.PeerCertificates)))
				w.WriteHeader(http.StatusOK)
			}))
			server.Config.ErrorLog = log.New(io.Discard, "", 0)
			server.TLS = &tls.Config{
				MinVersion:   tls.VersionTLS12,
				Certificates: []tls.Certificate{fixture.certificate(t, "synthetic hint server", []string{host}, false)},
				ClientAuth:   tls.RequireAnyClientCert,
				ClientCAs:    advertised,
				VerifyConnection: func(state tls.ConnectionState) error {
					if len(state.PeerCertificates) == 0 {
						return errors.New("the attested client certificate was not presented")
					}
					_, err := state.PeerCertificates[0].Verify(x509.VerifyOptions{
						Roots:     fixture.trust,
						KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
					})
					return err
				},
			}
			server.StartTLS()
			t.Cleanup(server.Close)
			capability := newX509Capability(t, fixture.transport(t, server))
			method, target := http.MethodGet, "https://"+host+"/v1/models"
			if host == x509TransportIssuer {
				method, target = http.MethodPost, "https://"+host+"/oauth/token"
			}
			request := x509TransportRequest(t, method, target)
			if host == x509TransportAPI {
				request.Header.Set("Authorization", "Bearer synthetic-attested-token")
			}
			response, err := capability.Do(request)
			if err != nil {
				t.Fatalf("mutually authenticated handshake with an unrelated acceptable-CA hint: %v", err)
			}
			if err := response.Body.Close(); err != nil {
				t.Fatalf("close mutually authenticated response: %v", err)
			}
			if got := presented.Load(); got != 1 {
				t.Errorf("presented workload certificates = %d, want one", got)
			}
		})
	}
}

func TestX509WorkloadIdentityRejectsPermanentRemoteCertificateAlertsWithoutRetry(t *testing.T) {
	fixture := newX509TransportFixture(t)
	unrelated := x509.NewCertPool()
	unrelated.AddCert(fixture.certificate(t, "synthetic unrelated client authority", nil, false).Leaf)
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		t.Error("untrusted workload certificate reached the issuer handler")
		w.WriteHeader(http.StatusOK)
	}))
	server.Config.ErrorLog = log.New(io.Discard, "", 0)
	server.TLS = &tls.Config{
		MinVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{fixture.certificate(t, "synthetic refusing issuer", []string{x509TransportIssuer}, false)},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    unrelated,
	}
	server.StartTLS()
	t.Cleanup(server.Close)
	template := fixture.transport(t, server)
	dial := template.DialContext
	var handshakes atomic.Int32
	template.DialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
		handshakes.Add(1)
		return dial(ctx, network, address)
	}
	capability := newX509Capability(t, template)
	identity, err := auth.NewX509WorkloadIdentityAuth(auth.X509WorkloadIdentity{
		IdentityProviderID: "synthetic-identity-provider",
		ServiceAccountID:   "synthetic-service-account",
		Transport:          capability,
	})
	if err != nil {
		t.Fatalf("create statically attested workload identity: %v", err)
	}
	token, err := identity.GetToken(t.Context(), capability)
	if err == nil || token != "" {
		t.Fatalf("permanently rejected workload certificate returned token=%q error=%v", token, err)
	}
	if got := handshakes.Load(); got != 1 {
		t.Errorf("permanent client-certificate rejection attempted %d TLS handshakes, want one", got)
	}
	if cause := errors.Unwrap(err); cause == nil || strings.Contains(cause.Error(), "synthetic") {
		t.Errorf("permanent TLS classification = %v, want a safe non-sensitive sentinel", cause)
	}
}

func TestX509TransportResponseHeaderTimeoutDoesNotLimitStreamingBodies(t *testing.T) {
	fixture := newX509TransportFixture(t)
	server := fixture.server(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		flusher, ok := w.(http.Flusher)
		if !ok {
			t.Error("synthetic streaming response cannot flush headers")
			return
		}
		flusher.Flush()
		timer := time.NewTimer(75 * time.Millisecond)
		defer timer.Stop()
		<-timer.C
		_, _ = io.WriteString(w, "synthetic delayed streaming payload")
	}))
	template := fixture.transport(t, server)
	template.ResponseHeaderTimeout = 20 * time.Millisecond
	capability := newX509Capability(t, template)
	request := x509TransportRequest(t, http.MethodGet, "https://"+x509TransportAPI+"/v1/models")
	response, err := capability.Do(request)
	if err != nil {
		t.Fatalf("start bounded-header streaming response: %v", err)
	}
	defer func() { _ = response.Body.Close() }()
	body, err := io.ReadAll(response.Body)
	if err != nil || string(body) != "synthetic delayed streaming payload" {
		t.Errorf("streaming body beyond header deadline = %q, error=%v", body, err)
	}
}

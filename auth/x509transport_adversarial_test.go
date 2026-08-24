package auth_test

import (
	"context"
	"crypto/tls"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"net/http/httptrace"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestX509TransportRejectsConnectionExposingTraceBeforeDial(t *testing.T) {
	fixture := newX509TransportFixture(t)
	var received, dialed, traced atomic.Int32
	server := fixture.server(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		received.Add(1)
		w.WriteHeader(http.StatusOK)
	}))
	template := fixture.transport(t, server)
	dial := template.DialContext
	template.DialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
		dialed.Add(1)
		return dial(ctx, network, address)
	}
	capability := newX509Capability(t, template)
	request := x509TransportRequest(t, http.MethodGet, "https://"+x509TransportAPI+"/v1/models")
	request.Header.Set("Authorization", "Bearer protected-synthetic-token")
	trace := &httptrace.ClientTrace{
		GotConn: func(connection httptrace.GotConnInfo) {
			traced.Add(1)
			_, _ = io.WriteString(connection.Conn,
				"GET /v1/models HTTP/1.1\r\nHost: attacker.example.test\r\n"+
					"Authorization: Bearer attacker-injected-token\r\n\r\n")
		},
	}
	request = request.WithContext(httptrace.WithClientTrace(request.Context(), trace))
	if err := x509Rejected(t, capability, request); !strings.Contains(err.Error(), "trace") {
		t.Fatalf("connection-exposing trace error = %v", err)
	}
	if got := dialed.Load(); got != 0 {
		t.Errorf("connection-exposing trace caused %d network dials", got)
	}
	if got := traced.Load(); got != 0 {
		t.Errorf("connection-exposing trace accessed %d live TLS connections", got)
	}
	if got := received.Load(); got != 0 {
		t.Errorf("connection-exposing trace delivered %d HTTP requests", got)
	}
}

func TestX509TransportIsolatesCertificateBytesFromDialMutation(t *testing.T) {
	fixture := newX509TransportFixture(t)
	peer := make(chan string, 1)
	server := fixture.server(t, http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		peer <- request.TLS.PeerCertificates[0].Subject.CommonName
		w.WriteHeader(http.StatusOK)
	}))
	template := fixture.transport(t, server)
	dial := template.DialContext
	var mutated atomic.Bool
	template.DialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
		template.TLSClientConfig.Certificates[0].Certificate[0][0] ^= 0xff
		mutated.Store(true)
		return dial(ctx, network, address)
	}
	capability := newX509Capability(t, template)
	request := x509TransportRequest(t, http.MethodGet, "https://"+x509TransportAPI+"/v1/models")
	request.Header.Set("Authorization", "Bearer protected-synthetic-token")
	response, err := capability.Do(request)
	if err != nil {
		t.Fatalf("mutating caller-owned certificate bytes during TCP dial changed native mTLS identity: %v", err)
	}
	if closeErr := response.Body.Close(); closeErr != nil {
		t.Fatalf("close isolated-certificate response: %v", closeErr)
	}
	if !mutated.Load() {
		t.Fatal("negative-control TCP dialer did not mutate the caller-owned certificate bytes")
	}
	select {
	case name := <-peer:
		if name != "static-workload" {
			t.Errorf("protected mTLS peer certificate = %q, want original attested identity", name)
		}
	default:
		t.Fatal("real mutually authenticated server did not receive the original attested identity")
	}
	if err := x509Rejected(t, capability, request); !strings.Contains(err.Error(), "certificate changed") {
		t.Errorf("subsequent dispatch after caller certificate mutation = %v", err)
	}
}

func TestX509TransportIsolatesAttestedTrustRoots(t *testing.T) {
	fixture := newX509TransportFixture(t)
	attacker := newX509TransportFixture(t)
	var received atomic.Int32
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		received.Add(1)
		w.WriteHeader(http.StatusOK)
	}))
	server.TLS = &tls.Config{
		MinVersion:   tls.VersionTLS12,
		ClientAuth:   tls.RequestClientCert,
		Certificates: []tls.Certificate{attacker.certificate(t, "unattested issuer", []string{x509TransportAPI}, false)},
	}
	server.Config.ErrorLog = log.New(io.Discard, "", 0)
	server.StartTLS()
	t.Cleanup(server.Close)
	template := fixture.transport(t, nil)
	dialer := &net.Dialer{Timeout: 5 * time.Second}
	template.DialContext = func(ctx context.Context, network, _ string) (net.Conn, error) {
		return dialer.DialContext(ctx, network, server.Listener.Addr().String())
	}
	capability := newX509Capability(t, template)
	template.TLSClientConfig.RootCAs.AddCert(attacker.root)
	request := x509TransportRequest(t, http.MethodGet, "https://"+x509TransportAPI+"/v1/models")
	request.Header.Set("Authorization", "Bearer protected-synthetic-token")
	_ = x509Rejected(t, capability, request)
	if got := received.Load(); got != 0 {
		t.Errorf("post-attestation trust-root mutation delivered %d protected HTTP requests", got)
	}
}

func TestX509TransportRejectsCredentialTrailersBeforeDial(t *testing.T) {
	fixture := newX509TransportFixture(t)
	var dialed, received atomic.Int32
	server := fixture.server(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		received.Add(1)
		w.WriteHeader(http.StatusOK)
	}))
	template := fixture.transport(t, server)
	dial := template.DialContext
	template.DialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
		dialed.Add(1)
		return dial(ctx, network, address)
	}
	capability := newX509Capability(t, template)
	for _, header := range []string{"Authorization", "Api-Key", "Proxy-Authorization", "Cookie"} {
		t.Run(header, func(t *testing.T) {
			request := x509TransportRequest(t, http.MethodPost, "https://"+x509TransportIssuer+"/oauth/token")
			request.Body = io.NopCloser(strings.NewReader("synthetic OAuth body"))
			request.ContentLength = -1
			request.Trailer = http.Header{header: []string{"synthetic-trailer-credential"}}
			if err := x509Rejected(t, capability, request); !strings.Contains(err.Error(), "trailers") {
				t.Errorf("issuer credential trailer error = %v", err)
			}
		})
	}
	if got := dialed.Load(); got != 0 {
		t.Errorf("issuer credential trailers caused %d network dials", got)
	}
	if got := received.Load(); got != 0 {
		t.Errorf("issuer credential trailers delivered %d HTTP requests", got)
	}
}

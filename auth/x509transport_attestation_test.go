package auth

import (
	"context"
	"crypto/tls"
	"errors"
	"io"
	"net"
	"net/http"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestX509TransportUsesBoundedNativeDefaultsWithoutMutatingItsTemplate(t *testing.T) {
	fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, x509ValidExchangeResponse())
	}))
	template := fixture.template
	template.DialContext = nil
	capability, err := NewX509Transport(template)
	if err != nil {
		t.Fatalf("attest timeout-default workload transport: %v", err)
	}
	t.Cleanup(func() { _ = capability.Close() })
	if capability.transport.DialContext == nil {
		t.Fatal("attested native transport did not install a bounded default TCP dialer")
	}
	if got := capability.transport.TLSHandshakeTimeout; got != x509DefaultTLSHandshakeTimeout {
		t.Errorf("default TLS handshake timeout = %s, want %s", got, x509DefaultTLSHandshakeTimeout)
	}
	if got := capability.transport.ResponseHeaderTimeout; got != x509DefaultResponseHeaderTimeout {
		t.Errorf("default response-header timeout = %s, want %s", got, x509DefaultResponseHeaderTimeout)
	}
	if template.DialContext != nil || template.TLSHandshakeTimeout != 0 || template.ResponseHeaderTimeout != 0 {
		t.Error("installing private bounded defaults mutated the caller-owned transport template")
	}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	if _, err := capability.transport.DialContext(ctx, "tcp", "127.0.0.1:1"); !errors.Is(err, context.Canceled) {
		t.Errorf("bounded TCP dialer ignored caller cancellation: %v", err)
	}
}

func TestX509TransportPreservesExplicitNativeTimeoutAndDialConfiguration(t *testing.T) {
	fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, x509ValidExchangeResponse())
	}))
	template := fixture.template
	var dialed atomic.Int32
	original := template.DialContext
	template.DialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
		dialed.Add(1)
		return original(ctx, network, address)
	}
	template.TLSHandshakeTimeout = 3 * time.Second
	template.ResponseHeaderTimeout = 4 * time.Second
	capability, err := NewX509Transport(template)
	if err != nil {
		t.Fatalf("attest explicitly configured workload transport: %v", err)
	}
	t.Cleanup(func() { _ = capability.Close() })
	if got := capability.transport.TLSHandshakeTimeout; got != 3*time.Second {
		t.Errorf("explicit TLS handshake timeout = %s, want three seconds", got)
	}
	if got := capability.transport.ResponseHeaderTimeout; got != 4*time.Second {
		t.Errorf("explicit response-header timeout = %s, want four seconds", got)
	}
	if _, err := x509Exchange(t.Context(), capability, "synthetic-provider", "synthetic-account"); err != nil {
		t.Fatalf("exchange using the preserved caller TCP dialer: %v", err)
	}
	if got := dialed.Load(); got != 1 {
		t.Errorf("preserved custom TCP dialer ran %d times, want one", got)
	}
}

func TestX509TransportRejectsReplacedCapabilityOwnedCertificateSelector(t *testing.T) {
	for _, test := range []struct {
		name    string
		replace func(*X509Transport, *X509Transport, *atomic.Int32)
	}{
		{name: "removed internal selector", replace: func(first, _ *X509Transport, _ *atomic.Int32) {
			first.tlsConfig.GetClientCertificate = nil
		}},
		{name: "caller-controlled callback", replace: func(first, _ *X509Transport, invoked *atomic.Int32) {
			first.tlsConfig.GetClientCertificate = func(*tls.CertificateRequestInfo) (*tls.Certificate, error) {
				invoked.Add(1)
				return &first.tlsConfig.Certificates[0], nil
			}
		}},
		{name: "another capability with the same closure code", replace: func(first, second *X509Transport, _ *atomic.Int32) {
			if reflect.ValueOf(first.tlsConfig.GetClientCertificate).Pointer() !=
				reflect.ValueOf(second.tlsConfig.GetClientCertificate).Pointer() {
				t.Fatal("separate capabilities unexpectedly used different static-selector function code")
			}
			first.tlsConfig.GetClientCertificate = second.tlsConfig.GetClientCertificate
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			first := newX509ExchangeFixture(t, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
			second := newX509ExchangeFixture(t, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
			var invoked atomic.Int32
			test.replace(first.capability, second.capability, &invoked)
			if err := first.capability.validateAttestation(); err == nil || !strings.Contains(err.Error(), "selection changed") {
				t.Fatalf("replaced private selector attestation error = %v", err)
			}
			if got := invoked.Load(); got != 0 {
				t.Errorf("untrusted replacement selector executed %d times", got)
			}
		})
	}
}

func TestX509TransportRejectsMissingAttestedCertificatesWithoutPanic(t *testing.T) {
	for _, test := range []struct {
		name   string
		remove func(*X509Transport)
	}{
		{name: "caller template", remove: func(transport *X509Transport) {
			transport.templateTLS.Certificates = nil
		}},
		{name: "capability-owned configuration", remove: func(transport *X509Transport) {
			transport.tlsConfig.Certificates = nil
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
			test.remove(fixture.capability)
			if err := fixture.capability.validateAttestation(); err == nil ||
				!strings.Contains(err.Error(), "requires exactly one static certificate") {
				t.Fatalf("missing attested certificate error = %v, want a safe validation failure", err)
			}
		})
	}
}

func TestX509TransportClassifiesOnlyPermanentRemoteCertificateAlerts(t *testing.T) {
	for _, test := range []struct {
		name      string
		operation string
		message   string
		permanent bool
	}{
		{name: "bad certificate", operation: "remote error", message: "tls: bad certificate", permanent: true},
		{name: "unsupported certificate", operation: "remote error", message: "tls: unsupported certificate", permanent: true},
		{name: "revoked certificate", operation: "remote error", message: "tls: revoked certificate", permanent: true},
		{name: "expired certificate", operation: "remote error", message: "tls: expired certificate", permanent: true},
		{name: "unknown certificate", operation: "remote error", message: "tls: unknown certificate", permanent: true},
		{name: "unknown authority", operation: "remote error", message: "tls: unknown certificate authority", permanent: true},
		{name: "required certificate", operation: "remote error", message: "tls: certificate required", permanent: true},
		{name: "transient internal TLS error", operation: "remote error", message: "tls: internal error"},
		{name: "local certificate-shaped error", operation: "dial", message: "tls: bad certificate"},
		{name: "sensitive unrelated network error", operation: "remote error", message: "synthetic-private-token"},
	} {
		t.Run(test.name, func(t *testing.T) {
			err := &net.OpError{Op: test.operation, Err: errors.New(test.message)}
			if got := x509PermanentClientCertificateAlert(err); got != test.permanent {
				t.Errorf("remote TLS alert permanent = %v, want %v", got, test.permanent)
			}
		})
	}
	if x509PermanentClientCertificateAlert(nil) {
		t.Error("nil transport error was classified as a permanent TLS rejection")
	}
}

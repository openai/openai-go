package auth_test

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/openai/openai-go/v3/auth"
)

func TestX509WorkloadIdentityAuthRejectsHTTPTransportChange(t *testing.T) {
	t.Run("explicit transport", func(t *testing.T) {
		var callsA atomic.Int32
		var callsB atomic.Int32
		clientA := nativeX509HTTPClient(t, func(*http.Request) (*http.Response, error) {
			callsA.Add(1)
			return tokenResponse("token-a", 60), nil
		})
		clientB := nativeX509HTTPClient(t, func(*http.Request) (*http.Response, error) {
			callsB.Add(1)
			return tokenResponse("token-b", 60), nil
		})
		httpClient := &http.Client{Transport: clientA.Transport}
		wia, err := auth.NewX509WorkloadIdentityAuth(testX509WorkloadIdentity())
		if err != nil {
			t.Fatalf("NewX509WorkloadIdentityAuth() error = %v", err)
		}

		token, err := wia.GetToken(t.Context(), httpClient)
		if err != nil {
			t.Fatalf("first GetToken() error = %v", err)
		}
		if token != "token-a" {
			t.Fatalf("first GetToken() = %q, want token-a", token)
		}
		token, err = wia.GetToken(t.Context(), httpClient)
		if err != nil {
			t.Fatalf("cached GetToken() error = %v", err)
		}
		if token != "token-a" {
			t.Fatalf("cached GetToken() = %q, want token-a", token)
		}

		httpClient.Transport = clientB.Transport
		ctx := t.Context()
		const callers = 16
		errs := make(chan error, callers)
		for range callers {
			go func() {
				_, getErr := wia.GetToken(ctx, httpClient)
				errs <- getErr
			}()
		}
		for range callers {
			if getErr := <-errs; getErr == nil {
				t.Error("GetToken() after HTTP transport change error = nil")
			}
		}
		if got, want := callsA.Load(), int32(1); got != want {
			t.Errorf("HTTP transport A exchange calls = %d, want %d", got, want)
		}
		if got := callsB.Load(); got != 0 {
			t.Errorf("HTTP transport B exchange calls = %d, want 0", got)
		}
	})

	t.Run("default transport", func(t *testing.T) {
		originalDefaultTransport := http.DefaultTransport
		t.Cleanup(func() {
			http.DefaultTransport = originalDefaultTransport
		})

		var callsA atomic.Int32
		var callsB atomic.Int32
		clientA := nativeX509HTTPClient(t, func(*http.Request) (*http.Response, error) {
			callsA.Add(1)
			return tokenResponse("token-a", 60), nil
		})
		http.DefaultTransport = clientA.Transport
		httpClient := &http.Client{}
		wia, err := auth.NewX509WorkloadIdentityAuth(testX509WorkloadIdentity())
		if err != nil {
			t.Fatalf("NewX509WorkloadIdentityAuth() error = %v", err)
		}
		if _, err := wia.GetToken(t.Context(), httpClient); err != nil {
			t.Fatalf("first GetToken() error = %v", err)
		}

		clientB := nativeX509HTTPClient(t, func(*http.Request) (*http.Response, error) {
			callsB.Add(1)
			return tokenResponse("token-b", 60), nil
		})
		http.DefaultTransport = clientB.Transport
		if _, err := wia.GetToken(t.Context(), httpClient); err == nil {
			t.Fatal("GetToken() after default HTTP transport change error = nil")
		}
		if got, want := callsA.Load(), int32(1); got != want {
			t.Errorf("default HTTP transport A exchange calls = %d, want %d", got, want)
		}
		if got := callsB.Load(); got != 0 {
			t.Errorf("default HTTP transport B exchange calls = %d, want 0", got)
		}
	})
}

func TestX509WorkloadIdentityAuthRequiresImmutableNativeTransportIdentity(t *testing.T) {
	staticCertificate := []tls.Certificate{{Certificate: [][]byte{{1}}, PrivateKey: struct{}{}}}
	testCases := []struct {
		name      string
		transport *http.Transport
	}{
		{name: "missing static certificate", transport: &http.Transport{}},
		{
			name: "certificate selection hook",
			transport: &http.Transport{TLSClientConfig: &tls.Config{
				Certificates: staticCertificate,
				GetClientCertificate: func(*tls.CertificateRequestInfo) (*tls.Certificate, error) {
					return &staticCertificate[0], nil
				},
			}},
		},
		{
			name: "multiple static certificates",
			transport: &http.Transport{TLSClientConfig: &tls.Config{
				Certificates: append(staticCertificate, tls.Certificate{}),
			}},
		},
		{
			name: "legacy TLS dial hook",
			transport: &http.Transport{
				TLSClientConfig: &tls.Config{Certificates: staticCertificate},
				DialTLS: func(string, string) (net.Conn, error) {
					return nil, errors.New("must not be called")
				},
			},
		},
		{
			name: "context TLS dial hook",
			transport: &http.Transport{
				TLSClientConfig: &tls.Config{Certificates: staticCertificate},
				DialTLSContext: func(context.Context, string, string) (net.Conn, error) {
					return nil, errors.New("must not be called")
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			wia, err := auth.NewX509WorkloadIdentityAuth(testX509WorkloadIdentity())
			if err != nil {
				t.Fatalf("NewX509WorkloadIdentityAuth() error = %v", err)
			}
			if _, err := wia.GetToken(t.Context(), &http.Client{Transport: testCase.transport}); err == nil {
				t.Fatal("GetToken() error = nil")
			}
		})
	}
}

func TestX509WorkloadIdentityAuthRejectsClientSessionCacheBeforeExchange(t *testing.T) {
	var dialCalls atomic.Int32
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates:       []tls.Certificate{{Certificate: [][]byte{{1}}}},
			ClientSessionCache: tls.NewLRUClientSessionCache(1),
		},
		DialContext: func(context.Context, string, string) (net.Conn, error) {
			dialCalls.Add(1)
			return nil, errors.New("must not dial")
		},
	}
	wia, err := auth.NewX509WorkloadIdentityAuth(testX509WorkloadIdentity())
	if err != nil {
		t.Fatalf("NewX509WorkloadIdentityAuth() error = %v", err)
	}

	if _, err := wia.GetToken(t.Context(), &http.Client{Transport: transport}); err == nil || !strings.Contains(err.Error(), "session caching") {
		t.Fatalf("GetToken() error = %v, want client session cache rejection", err)
	}
	if got := dialCalls.Load(); got != 0 {
		t.Fatalf("dial calls = %d, want 0", got)
	}
}

func TestX509WorkloadIdentityAuthRejectsStaticCertificateMutation(t *testing.T) {
	var calls atomic.Int32
	httpClient := nativeX509HTTPClient(t, func(*http.Request) (*http.Response, error) {
		calls.Add(1)
		return tokenResponse("token-a", 60), nil
	})
	wia, err := auth.NewX509WorkloadIdentityAuth(testX509WorkloadIdentity())
	if err != nil {
		t.Fatalf("NewX509WorkloadIdentityAuth() error = %v", err)
	}
	if _, err := wia.GetToken(t.Context(), httpClient); err != nil {
		t.Fatalf("first GetToken() error = %v", err)
	}

	transport := httpClient.Transport.(*http.Transport)
	transport.TLSClientConfig.Certificates[0].Certificate[0] = []byte{2}
	if _, err := wia.GetToken(t.Context(), httpClient); err == nil || !strings.Contains(err.Error(), "cannot change client certificates") {
		t.Fatalf("GetToken() after certificate mutation error = %v, want identity rejection", err)
	}
	if got, want := calls.Load(), int32(1); got != want {
		t.Fatalf("token exchange calls = %d, want %d", got, want)
	}
}

func TestX509WorkloadIdentityAuthRejectsOpaqueRoundTripperBeforeExchange(t *testing.T) {
	var calls atomic.Int32
	httpClient := &http.Client{Transport: &closureTransport{fn: func(*http.Request) (*http.Response, error) {
		calls.Add(1)
		return tokenResponse("must-not-be-used", 60), nil
	}}}
	wia, err := auth.NewX509WorkloadIdentityAuth(testX509WorkloadIdentity())
	if err != nil {
		t.Fatalf("NewX509WorkloadIdentityAuth() error = %v", err)
	}

	if _, err := wia.GetToken(t.Context(), httpClient); err == nil || !strings.Contains(err.Error(), "native *http.Transport") {
		t.Fatalf("GetToken() error = %v, want opaque transport rejection", err)
	}
	if got := calls.Load(); got != 0 {
		t.Fatalf("token exchange calls = %d, want 0", got)
	}
}

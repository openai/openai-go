package auth_test

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/openai/openai-go/v3/auth"
)

func TestX509TransportRejectsTunnelsAndAmbiguousFramingBeforeDial(t *testing.T) {
	fixture := newX509TransportFixture(t)
	var received, dialed atomic.Int32
	server := fixture.server(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		received.Add(1)
		w.WriteHeader(http.StatusOK)
	}))
	template := fixture.transport(t, server)
	originalDial := template.DialContext
	template.DialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
		dialed.Add(1)
		return originalDial(ctx, network, address)
	}
	capability := newX509Capability(t, template)
	for _, test := range []struct {
		name   string
		change func(*http.Request)
	}{
		{name: "CONNECT tunnel", change: func(request *http.Request) { request.Method = http.MethodConnect }},
		{name: "TRACE tunnel", change: func(request *http.Request) { request.Method = http.MethodTrace }},
		{name: "HTTP/2 preface", change: func(request *http.Request) { request.Method = "PRI" }},
		{name: "custom tunnel method", change: func(request *http.Request) { request.Method = "SYNTHETIC-TUNNEL" }},
		{name: "identity transfer encoding", change: func(request *http.Request) {
			request.TransferEncoding = []string{"identity"}
		}},
		{name: "multiple transfer encodings", change: func(request *http.Request) {
			request.TransferEncoding = []string{"chunked", "identity"}
		}},
		{name: "unsupported transfer encoding", change: func(request *http.Request) {
			request.TransferEncoding = []string{"gzip"}
		}},
		{name: "explicit chunked transfer encoding", change: func(request *http.Request) {
			request.TransferEncoding = []string{"chunked"}
		}},
		{name: "chunked known-length conflict", change: func(request *http.Request) {
			request.TransferEncoding = []string{"chunked"}
			request.ContentLength = 5
		}},
		{name: "invalid negative content length", change: func(request *http.Request) {
			request.ContentLength = -2
		}},
		{name: "missing body with content length", change: func(request *http.Request) {
			_ = request.Body.Close()
			request.Body = nil
			request.ContentLength = 5
		}},
		{name: "chunked without body", change: func(request *http.Request) {
			_ = request.Body.Close()
			request.Body = nil
			request.ContentLength = 0
			request.TransferEncoding = []string{"chunked"}
		}},
		{name: "canonical transfer encoding header", change: x509UnsafeFramingHeader("Transfer-Encoding", "identity")},
		{name: "lowercase transfer encoding alias", change: x509UnsafeFramingHeader("transfer-encoding", "identity")},
		{name: "underscore transfer encoding alias", change: x509UnsafeFramingHeader("transfer_encoding", "identity")},
		{name: "canonical content length header", change: x509UnsafeFramingHeader("Content-Length", "0")},
		{name: "lowercase content length alias", change: x509UnsafeFramingHeader("content-length", "0")},
		{name: "underscore content length alias", change: x509UnsafeFramingHeader("content_length", "0")},
		{name: "lowercase trailer alias", change: x509UnsafeFramingHeader("trailer", "Authorization")},
		{name: "underscore trailer alias", change: x509UnsafeFramingHeader("trailer", "Host")},
		{name: "connection upgrade", change: x509UnsafeFramingHeader("connection", "Upgrade")},
		{name: "protocol upgrade", change: x509UnsafeFramingHeader("upgrade", "attacker-protocol")},
		{name: "HTTP/2 upgrade settings", change: x509UnsafeFramingHeader("http2_settings", "attacker-settings")},
		{name: "proxy connection", change: x509UnsafeFramingHeader("proxy_connection", "keep-alive")},
		{name: "hop-by-hop TE", change: x509UnsafeFramingHeader("te", "trailers")},
		{name: "lowercase Host alias", change: x509UnsafeFramingHeader("host", "attacker.example.test")},
	} {
		t.Run(test.name, func(t *testing.T) {
			request := x509TransportRequest(t, http.MethodPost, "https://"+x509TransportAPI+"/v1/models")
			request.Header.Set("Authorization", "Bearer protected-synthetic-token")
			body := &x509TrackedRequestBody{Reader: strings.NewReader(
				"synthetic-body\r\nGET /v1/stolen HTTP/1.1\r\n" +
					"Host: attacker.example.test\r\nAuthorization: Bearer injected-secret\r\n\r\n",
			)}
			request.Body = body
			request.ContentLength = -1
			test.change(request)
			_ = x509Rejected(t, capability, request)
			if got := body.closed.Load(); got != 1 {
				t.Errorf("rejected framing closed its request body %d times, want once", got)
			}
			if got := dialed.Load(); got != 0 {
				t.Errorf("ambiguous HTTP framing caused %d network dials", got)
			}
			if got := received.Load(); got != 0 {
				t.Errorf("ambiguous HTTP framing delivered %d tunneled or smuggled requests", got)
			}
		})
	}
}

func x509UnsafeFramingHeader(name, value string) func(*http.Request) {
	return func(request *http.Request) {
		request.Header[name] = []string{value}
	}
}

func TestX509TransportPreservesSafelyChunkedStreamingBodies(t *testing.T) {
	fixture := newX509TransportFixture(t)
	var received atomic.Int32
	payload := "synthetic stream\r\nGET /v1/stolen HTTP/1.1\r\nHost: attacker.example.test\r\n\r\n"
	server := fixture.server(t, http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		received.Add(1)
		body, err := io.ReadAll(request.Body)
		if err != nil || string(body) != payload {
			t.Errorf("safely chunked workload body = %q, error = %v", body, err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	capability := newX509Capability(t, fixture.transport(t, server))
	for _, knownLength := range []bool{false, true} {
		request := x509TransportRequest(t, http.MethodPost, "https://"+x509TransportAPI+"/v1/models")
		request.Header.Set("Authorization", "Bearer protected-synthetic-token")
		request.Body = io.NopCloser(strings.NewReader(payload))
		request.ContentLength = -1
		if knownLength {
			request.ContentLength = int64(len(payload))
		}
		response, err := capability.Do(request)
		if err != nil {
			t.Fatalf("safely framed streaming request (known length=%t): %v", knownLength, err)
		}
		if closeErr := response.Body.Close(); closeErr != nil {
			t.Fatalf("close safely chunked streaming response: %v", closeErr)
		}
	}
	if got := received.Load(); got != 2 {
		t.Errorf("safe chunked streams delivered %d HTTP requests, want exactly two", got)
	}
}

func TestX509TransportRejectsSessionKeyLoggingBeforeTLS(t *testing.T) {
	fixture := newX509TransportFixture(t)
	var received atomic.Int32
	server := fixture.server(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		received.Add(1)
		w.WriteHeader(http.StatusOK)
	}))
	var leaked bytes.Buffer
	template := fixture.transport(t, server)
	template.TLSClientConfig.KeyLogWriter = &leaked
	if capability, err := auth.NewX509Transport(template); capability != nil || err == nil ||
		!strings.Contains(err.Error(), "session key logging") {
		t.Fatalf("attest caller-controlled TLS key logger = capability:%v error:%v", capability, err)
	}
	template.TLSClientConfig.KeyLogWriter = nil
	capability := newX509Capability(t, template)
	template.TLSClientConfig.KeyLogWriter = &leaked
	request := x509TransportRequest(t, http.MethodGet, "https://"+x509TransportAPI+"/v1/models")
	request.Header.Set("Authorization", "Bearer protected-synthetic-token")
	if err := x509Rejected(t, capability, request); !strings.Contains(err.Error(), "session key logging") {
		t.Fatalf("post-attestation TLS key logger error = %v", err)
	}
	if leaked.Len() != 0 || received.Load() != 0 {
		t.Errorf("rejected TLS key logger captured %d bytes and delivered %d requests", leaked.Len(), received.Load())
	}
}

func TestX509TransportFreezesMutableTLSPolicySlices(t *testing.T) {
	for _, test := range []struct {
		name    string
		prepare func(*http.Transport)
		mutate  func(*http.Transport)
	}{
		{
			name: "cipher suite preferences",
			prepare: func(template *http.Transport) {
				template.TLSClientConfig.MinVersion = tls.VersionTLS12
				template.TLSClientConfig.MaxVersion = tls.VersionTLS12
				template.TLSClientConfig.CipherSuites = []uint16{tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256}
			},
			mutate: func(template *http.Transport) {
				template.TLSClientConfig.CipherSuites[0] = tls.TLS_RSA_WITH_AES_128_CBC_SHA
			},
		},
		{
			name: "curve preferences",
			prepare: func(template *http.Transport) {
				template.TLSClientConfig.CurvePreferences = []tls.CurveID{tls.CurveP256}
			},
			mutate: func(template *http.Transport) {
				template.TLSClientConfig.CurvePreferences[0] = tls.CurveID(0xffff)
			},
		},
		{
			name: "TLS application protocol preferences",
			prepare: func(template *http.Transport) {
				template.TLSClientConfig.NextProtos = []string{"http/1.1"}
			},
			mutate: func(template *http.Transport) {
				template.TLSClientConfig.NextProtos[0] = "synthetic-attacker-protocol"
			},
		},
		{
			name: "client certificate signature algorithms",
			prepare: func(template *http.Transport) {
				template.TLSClientConfig.Certificates[0].SupportedSignatureAlgorithms = []tls.SignatureScheme{
					tls.ECDSAWithP256AndSHA256,
				}
			},
			mutate: func(template *http.Transport) {
				template.TLSClientConfig.Certificates[0].SupportedSignatureAlgorithms[0] = tls.PKCS1WithSHA256
			},
		},
		{
			name: "parsed client certificate leaf",
			mutate: func(template *http.Transport) {
				template.TLSClientConfig.Certificates[0].Leaf.RawIssuer = []byte("synthetic-untrusted-issuer")
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			fixture := newX509TransportFixture(t)
			var received atomic.Int32
			server := fixture.server(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				received.Add(1)
				w.WriteHeader(http.StatusOK)
			}))
			template := fixture.transport(t, server)
			if test.prepare != nil {
				test.prepare(template)
			}
			capability := newX509Capability(t, template)
			test.mutate(template)
			request := x509TransportRequest(t, http.MethodGet, "https://"+x509TransportAPI+"/v1/models")
			request.Header.Set("Authorization", "Bearer protected-synthetic-token")
			response, err := capability.Do(request)
			if err != nil {
				t.Fatalf("caller TLS-policy slice mutation changed the attested handshake: %v", err)
			}
			if closeErr := response.Body.Close(); closeErr != nil {
				t.Fatalf("close isolated TLS-policy response: %v", closeErr)
			}
			if got := received.Load(); got != 1 {
				t.Errorf("isolated TLS-policy mutation delivered %d requests", got)
			}
		})
	}
}

func TestX509TransportPreservesSanitizedNetworkTimeoutClassification(t *testing.T) {
	for _, test := range []struct {
		name             string
		handshakeTimeout bool
		contextDeadline  bool
	}{
		{name: "TLS handshake timeout", handshakeTimeout: true},
		{name: "response header timeout"},
		{name: "request context deadline", contextDeadline: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			fixture := newX509TransportFixture(t)
			var template *http.Transport
			if test.handshakeTimeout {
				listener, err := net.Listen("tcp", "127.0.0.1:0")
				if err != nil {
					t.Fatalf("start stalled TLS handshake listener: %v", err)
				}
				release := make(chan struct{})
				finished := make(chan struct{})
				go func() {
					defer close(finished)
					connection, acceptErr := listener.Accept()
					if acceptErr != nil {
						return
					}
					defer func() { _ = connection.Close() }()
					<-release
				}()
				t.Cleanup(func() {
					close(release)
					_ = listener.Close()
					<-finished
				})
				template = fixture.transport(t, nil)
				template.DialContext = func(ctx context.Context, network, _ string) (net.Conn, error) {
					return (&net.Dialer{}).DialContext(ctx, network, listener.Addr().String())
				}
				template.TLSHandshakeTimeout = 75 * time.Millisecond
			} else {
				server := fixture.server(t, http.HandlerFunc(func(_ http.ResponseWriter, request *http.Request) {
					<-request.Context().Done()
				}))
				template = fixture.transport(t, server)
				if !test.contextDeadline {
					template.ResponseHeaderTimeout = 75 * time.Millisecond
				}
			}
			capability := newX509Capability(t, template)
			request := x509TransportRequest(t, http.MethodGet,
				"https://"+x509TransportAPI+"/v1/models?private-query-token=synthetic-secret")
			request.Header.Set("Authorization", "Bearer private-header-token")
			if test.contextDeadline {
				ctx, cancel := context.WithTimeout(request.Context(), 75*time.Millisecond)
				defer cancel()
				request = request.WithContext(ctx)
			}
			err := x509Rejected(t, capability, request)
			var networkError net.Error
			if !errors.As(err, &networkError) || !networkError.Timeout() || !os.IsTimeout(err) {
				t.Errorf("sanitized timeout classification = error:%v net.Error:%t os.IsTimeout:%t",
					err, networkError != nil && networkError.Timeout(), os.IsTimeout(err))
			}
			if test.contextDeadline && !errors.Is(err, context.DeadlineExceeded) {
				t.Errorf("request deadline lost context identity: %v", err)
			}
			if cause := errors.Unwrap(err); cause != nil && !errors.Is(cause, context.DeadlineExceeded) {
				t.Errorf("native timeout retained its potentially sensitive underlying cause: %v", err)
			}
			for cause := err; cause != nil; cause = errors.Unwrap(cause) {
				if strings.Contains(cause.Error(), "private-") || strings.Contains(cause.Error(), "synthetic-secret") {
					t.Errorf("network timeout exposed request credentials: %q", cause.Error())
				}
			}
		})
	}
}

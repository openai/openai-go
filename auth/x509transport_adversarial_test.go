package auth_test

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
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

	"github.com/openai/openai-go/v3/auth"
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

func TestX509TransportRejectsLateInstalledConnectionTraceBeforeDial(t *testing.T) {
	fixture := newX509TransportFixture(t)
	var received, dialed, installed, exposed atomic.Int32
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
	trace := new(httptrace.ClientTrace)
	trace.GetConn = func(string) {
		installed.Add(1)
		trace.GotConn = func(connection httptrace.GotConnInfo) {
			exposed.Add(1)
			_, _ = io.WriteString(connection.Conn,
				"GET /v1/models HTTP/1.1\r\nHost: attacker.example.test\r\n"+
					"Authorization: Bearer attacker-injected-token\r\n\r\n")
		}
	}
	request = request.WithContext(httptrace.WithClientTrace(request.Context(), trace))
	if trace.GotConn != nil {
		t.Fatal("negative-control connection hook must initially be absent")
	}
	if err := x509Rejected(t, capability, request); !strings.Contains(err.Error(), "trace") {
		t.Fatalf("late-installing trace error = %v", err)
	}
	if got := installed.Load(); got != 0 {
		t.Errorf("late-installing trace executed %d initial callbacks", got)
	}
	if got := exposed.Load(); got != 0 {
		t.Errorf("late-installing trace accessed %d live TLS connections", got)
	}
	if got := dialed.Load(); got != 0 {
		t.Errorf("late-installing trace caused %d network dials", got)
	}
	if got := received.Load(); got != 0 {
		t.Errorf("late-installing trace delivered %d HTTP requests", got)
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

func TestX509TransportRejectsEmptyAuthorizationAliasesWithoutPanicking(t *testing.T) {
	fixture := newX509TransportFixture(t)
	capability := newX509Capability(t, fixture.transport(t, nil))
	for _, values := range [][]string{nil, {}} {
		request := x509TransportRequest(t, http.MethodGet, "https://"+x509TransportAPI+"/v1/models")
		request.Header["Authorization"] = []string{"Bearer synthetic-valid-token"}
		request.Header["authorization"] = values
		if err := x509Rejected(t, capability, request); !strings.Contains(err.Error(), "Authorization") {
			t.Errorf("empty differently cased Authorization alias error = %v", err)
		}
	}
}

func TestX509TransportFiltersLateAppearingTraceAndPreservesContext(t *testing.T) {
	fixture := newX509TransportFixture(t)
	var received, exposed, contextPreserved atomic.Int32
	server := fixture.server(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		received.Add(1)
		w.WriteHeader(http.StatusOK)
	}))
	template := fixture.transport(t, server)
	type contextKey struct{}
	parent, cancel := context.WithTimeout(context.WithValue(t.Context(), contextKey{}, "synthetic-context-value"), 5*time.Second)
	defer cancel()
	dial := template.DialContext
	template.DialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
		if ctx.Value(contextKey{}) == "synthetic-context-value" {
			contextPreserved.Add(1)
		}
		return dial(ctx, network, address)
	}
	capability := newX509Capability(t, template)
	trace := &httptrace.ClientTrace{
		GotConn: func(connection httptrace.GotConnInfo) {
			exposed.Add(1)
			_, _ = io.WriteString(connection.Conn,
				"GET /v1/models HTTP/1.1\r\nHost: attacker.example.test\r\n"+
					"Authorization: Bearer attacker-injected-token\r\n\r\n")
		},
	}
	mutable := &x509EmergingTraceContext{
		Context: parent,
		traced:  httptrace.WithClientTrace(context.Background(), trace),
	}
	request := x509TransportRequest(t, http.MethodGet, "https://"+x509TransportAPI+"/v1/models")
	request.Header.Set("Authorization", "Bearer protected-synthetic-token")
	response, err := capability.Do(request.WithContext(mutable))
	if err != nil {
		t.Fatalf("request with a late-appearing filtered HTTP trace: %v", err)
	}
	if err := response.Body.Close(); err != nil {
		t.Fatalf("close filtered-trace response: %v", err)
	}
	if got := mutable.lookups.Load(); got < 2 {
		t.Errorf("mutable trace was queried %d times; expected the native dispatch to query again", got)
	}
	if got := exposed.Load(); got != 0 {
		t.Errorf("late-appearing HTTP trace accessed %d authenticated sockets", got)
	}
	if got := contextPreserved.Load(); got != 1 {
		t.Errorf("native TCP dial observed the preserved context value %d times", got)
	}
	if got := received.Load(); got != 1 {
		t.Errorf("filtered-trace request delivered %d HTTP requests", got)
	}
}

type x509EmergingTraceContext struct {
	context.Context
	traced  context.Context
	lookups atomic.Int32
}

func (ctx *x509EmergingTraceContext) Value(key any) any {
	value := ctx.traced.Value(key)
	if _, ok := value.(*httptrace.ClientTrace); ok {
		if ctx.lookups.Add(1) == 1 {
			return nil
		}
		return value
	}
	return ctx.Context.Value(key)
}

func TestX509TransportRejectsProtocolUpgradeWithoutExposingSocket(t *testing.T) {
	fixture := newX509TransportFixture(t)
	closed := make(chan error, 1)
	server := fixture.server(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		hijacker, ok := w.(http.Hijacker)
		if !ok {
			closed <- errors.New("synthetic server does not support hijacking")
			return
		}
		connection, buffered, err := hijacker.Hijack()
		if err != nil {
			closed <- fmt.Errorf("hijack synthetic upgrade: %w", err)
			return
		}
		defer func() { _ = connection.Close() }()
		if _, writeErr := fmt.Fprint(buffered,
			"HTTP/1.1 101 Switching Protocols\r\nConnection: Upgrade\r\nUpgrade: synthetic\r\n\r\n",
		); writeErr != nil {
			closed <- fmt.Errorf("write synthetic upgrade: %w", writeErr)
			return
		}
		if flushErr := buffered.Flush(); flushErr != nil {
			closed <- fmt.Errorf("flush synthetic upgrade: %w", flushErr)
			return
		}
		if deadlineErr := connection.SetReadDeadline(time.Now().Add(5 * time.Second)); deadlineErr != nil {
			closed <- fmt.Errorf("bound synthetic upgrade socket wait: %w", deadlineErr)
			return
		}
		var leaked [1]byte
		_, err = connection.Read(leaked[:])
		if err == nil {
			closed <- errors.New("authenticated upgraded socket remained writable")
			return
		}
		var timeout net.Error
		if errors.As(err, &timeout) && timeout.Timeout() {
			closed <- errors.New("authenticated upgraded socket was not closed")
			return
		}
		closed <- nil
	}))
	capability := newX509Capability(t, fixture.transport(t, server))
	request := x509TransportRequest(t, http.MethodGet, "https://"+x509TransportAPI+"/v1/models")
	request.Header.Set("Authorization", "Bearer protected-synthetic-token")
	if err := x509Rejected(t, capability, request); !strings.Contains(err.Error(), "protocol upgrades") &&
		!strings.Contains(err.Error(), "transport request failed") {
		t.Fatalf("authenticated protocol-upgrade error = %v", err)
	}
	select {
	case err := <-closed:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("authenticated upgraded socket was not released")
	}
}

func TestX509TransportRejectsNonReplayableRedirectResponses(t *testing.T) {
	for _, status := range []int{
		http.StatusMultipleChoices, http.StatusMovedPermanently, http.StatusFound,
		http.StatusSeeOther, http.StatusUseProxy, http.StatusTemporaryRedirect,
		http.StatusPermanentRedirect,
	} {
		t.Run(fmt.Sprint(status), func(t *testing.T) {
			fixture := newX509TransportFixture(t)
			var received atomic.Int32
			server := fixture.server(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				received.Add(1)
				w.Header().Set("Location", "https://attacker.example.test/private-synthetic-redirect-token")
				w.WriteHeader(status)
				_, _ = io.WriteString(w, "synthetic redirect response")
			}))
			capability := newX509Capability(t, fixture.transport(t, server))
			request := x509TransportRequest(t, http.MethodPost, "https://"+x509TransportIssuer+"/oauth/token")
			request.Body = io.NopCloser(strings.NewReader("synthetic non-replayable body"))
			request.ContentLength = -1
			if request.GetBody != nil {
				t.Fatal("synthetic redirect request unexpectedly became replayable")
			}
			err := x509Rejected(t, capability, request)
			cause := errors.Unwrap(err)
			if cause == nil || !strings.Contains(cause.Error(), "does not follow redirects") {
				t.Fatalf("non-replayable %d redirect error = %v, cause = %v", status, err, cause)
			}
			if strings.Contains(err.Error(), "private-synthetic-redirect-token") {
				t.Error("redirect error disclosed its sensitive destination")
			}
			if got := received.Load(); got != 1 {
				t.Errorf("non-replayable %d redirect delivered %d HTTP requests", status, got)
			}
		})
	}
}

func TestX509TransportClosesBodiesOnEveryPreDispatchFailure(t *testing.T) {
	for _, test := range []struct {
		name   string
		change func(*testing.T, *auth.X509Transport, *http.Request) *http.Request
	}{
		{
			name: "canceled context",
			change: func(t *testing.T, _ *auth.X509Transport, request *http.Request) *http.Request {
				t.Helper()
				ctx, cancel := context.WithCancel(request.Context())
				cancel()
				return request.WithContext(ctx)
			},
		},
		{
			name: "HTTP trace",
			change: func(_ *testing.T, _ *auth.X509Transport, request *http.Request) *http.Request {
				return request.WithContext(httptrace.WithClientTrace(request.Context(), new(httptrace.ClientTrace)))
			},
		},
		{
			name: "closed capability",
			change: func(t *testing.T, capability *auth.X509Transport, request *http.Request) *http.Request {
				t.Helper()
				if err := capability.Close(); err != nil {
					t.Fatalf("close capability before synthetic request: %v", err)
				}
				return request
			},
		},
		{
			name: "unsafe authority",
			change: func(_ *testing.T, _ *auth.X509Transport, request *http.Request) *http.Request {
				request.URL.Host = "attacker.example.test"
				return request
			},
		},
		{
			name: "credential trailer",
			change: func(_ *testing.T, _ *auth.X509Transport, request *http.Request) *http.Request {
				request.Trailer = http.Header{"Authorization": []string{"Bearer unsafe"}}
				return request
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			fixture := newX509TransportFixture(t)
			capability := newX509Capability(t, fixture.transport(t, nil))
			request := x509TransportRequest(t, http.MethodPost, "https://"+x509TransportIssuer+"/oauth/token")
			body := &x509TrackedRequestBody{Reader: strings.NewReader("synthetic request body")}
			request.Body = body
			request = test.change(t, capability, request)
			_ = x509Rejected(t, capability, request)
			if got := body.closed.Load(); got != 1 {
				t.Errorf("rejected request body closed %d times, want exactly once", got)
			}
		})
	}
}

type x509TrackedRequestBody struct {
	io.Reader
	closed atomic.Int32
}

func (body *x509TrackedRequestBody) Close() error {
	body.closed.Add(1)
	return nil
}

func TestX509TransportPreservesExplicitHTTP2ProtocolPolicy(t *testing.T) {
	fixture := newX509TransportFixture(t)
	var observedProtocol atomic.Int32
	server := newX509HTTP2TransportServer(t, fixture, http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		observedProtocol.Store(int32(request.ProtoMajor))
		w.WriteHeader(http.StatusOK)
	}))
	template := fixture.transport(t, server)
	template.Protocols = new(http.Protocols)
	template.Protocols.SetHTTP2(true)
	template.HTTP2 = &http.HTTP2Config{
		MaxReadFrameSize: 32 << 10,
		PingTimeout:      3 * time.Second,
	}
	capability := newX509Capability(t, template)
	template.Protocols.SetHTTP1(true)
	template.Protocols.SetHTTP2(false)
	template.HTTP2.MaxReadFrameSize = 64 << 10
	request := x509TransportRequest(t, http.MethodGet, "https://"+x509TransportAPI+"/v1/models")
	response, err := capability.Do(request)
	if err != nil {
		t.Fatalf("dispatch with an isolated explicit HTTP/2-only protocol policy: %v", err)
	}
	if err := response.Body.Close(); err != nil {
		t.Fatalf("close HTTP/2-only response: %v", err)
	}
	if got := observedProtocol.Load(); got != 2 {
		t.Errorf("isolated explicit HTTP/2-only request negotiated HTTP/%d", got)
	}
}

func TestX509TransportPreservesImplicitTLSNextProtoHTTP2Policy(t *testing.T) {
	for _, test := range []struct {
		name       string
		forceHTTP2 bool
		handler    bool
		want       int32
	}{
		{name: "empty protocol map disables forced HTTP/2", forceHTTP2: true, want: 1},
		{name: "HTTP/2 handler enables HTTP/2 without force", handler: true, want: 2},
	} {
		t.Run(test.name, func(t *testing.T) {
			fixture := newX509TransportFixture(t)
			var observed, intercepted atomic.Int32
			server := newX509HTTP2TransportServer(t, fixture, http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
				observed.Store(int32(request.ProtoMajor))
				w.WriteHeader(http.StatusOK)
			}))
			template := fixture.transport(t, server)
			template.ForceAttemptHTTP2 = test.forceHTTP2
			template.TLSNextProto = make(map[string]func(string, *tls.Conn) http.RoundTripper)
			if test.handler {
				template.TLSNextProto["h2"] = func(string, *tls.Conn) http.RoundTripper {
					intercepted.Add(1)
					return nil
				}
			}
			capability := newX509Capability(t, template)
			request := x509TransportRequest(t, http.MethodGet, "https://"+x509TransportAPI+"/v1/models")
			response, err := capability.Do(request)
			if err != nil {
				t.Fatalf("dispatch with implicit TLS protocol policy: %v", err)
			}
			if closeErr := response.Body.Close(); closeErr != nil {
				t.Fatalf("close implicit TLS protocol response: %v", closeErr)
			}
			if got := observed.Load(); got != test.want {
				t.Errorf("implicit TLS protocol policy negotiated HTTP/%d, want HTTP/%d", got, test.want)
			}
			if got := intercepted.Load(); got != 0 {
				t.Errorf("caller-owned TLS protocol handlers intercepted %d requests", got)
			}
		})
	}
}

func TestX509TransportAcceptsNativeInitializedHTTP2Templates(t *testing.T) {
	for _, initializeBefore := range []bool{false, true} {
		name := "initialized after attestation"
		if initializeBefore {
			name = "initialized before attestation"
		}
		t.Run(name, func(t *testing.T) {
			fixture := newX509TransportFixture(t)
			var received atomic.Int32
			server := newX509HTTP2TransportServer(t, fixture, http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
				if request.ProtoMajor == 2 {
					received.Add(1)
				}
				w.WriteHeader(http.StatusOK)
			}))
			template := fixture.transport(t, server)
			template.ForceAttemptHTTP2 = true
			var capability *auth.X509Transport
			if !initializeBefore {
				capability = newX509Capability(t, template)
			}
			request := x509TransportRequest(t, http.MethodGet, "https://"+x509TransportAPI+"/v1/models")
			response, err := (&http.Client{Transport: template}).Do(request)
			if err != nil {
				t.Fatalf("initialize the caller-owned native HTTP/2 template: %v", err)
			}
			if closeErr := response.Body.Close(); closeErr != nil {
				t.Fatalf("close caller-owned native HTTP/2 response: %v", closeErr)
			}
			if template.TLSNextProto["h2"] == nil {
				t.Fatal("negative-control native transport did not install its bundled HTTP/2 handler")
			}
			if initializeBefore {
				capability = newX509Capability(t, template)
			}
			request = x509TransportRequest(t, http.MethodGet, "https://"+x509TransportAPI+"/v1/models")
			response, err = capability.Do(request)
			if err != nil {
				t.Fatalf("dispatch after native HTTP/2 template initialization: %v", err)
			}
			if err := response.Body.Close(); err != nil {
				t.Fatalf("close isolated native HTTP/2 response: %v", err)
			}
			if got := received.Load(); got != 2 {
				t.Errorf("native and isolated adapters delivered %d mutually authenticated HTTP/2 requests", got)
			}
		})
	}
}

func TestX509TransportDoesNotRaceConcurrentNativeHTTP2Initialization(t *testing.T) {
	fixture := newX509TransportFixture(t)
	var received atomic.Int32
	server := newX509HTTP2TransportServer(t, fixture, http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		if request.ProtoMajor == 2 {
			received.Add(1)
		}
		w.WriteHeader(http.StatusOK)
	}))
	template := fixture.transport(t, server)
	template.ForceAttemptHTTP2 = true
	capability := newX509Capability(t, template)
	request := x509TransportRequest(t, http.MethodGet, "https://"+x509TransportAPI+"/v1/models")

	const isolatedRequests = 32
	ready := make(chan struct{}, isolatedRequests+1)
	start := make(chan struct{})
	results := make(chan error, isolatedRequests+1)
	go func() {
		ready <- struct{}{}
		<-start
		response, err := (&http.Client{Transport: template}).Do(request.Clone(request.Context()))
		if err == nil {
			err = response.Body.Close()
		}
		results <- err
	}()
	for range isolatedRequests {
		go func() {
			ready <- struct{}{}
			<-start
			response, err := capability.Do(request.Clone(request.Context()))
			if err == nil {
				err = response.Body.Close()
			}
			results <- err
		}()
	}
	for range isolatedRequests + 1 {
		<-ready
	}
	close(start)
	for range isolatedRequests + 1 {
		if err := <-results; err != nil {
			t.Errorf("concurrent native or isolated HTTP/2 request: %v", err)
		}
	}
	if template.TLSNextProto["h2"] == nil {
		t.Fatal("caller-owned transport did not initialize its native HTTP/2 protocol map")
	}
	if got := received.Load(); got != isolatedRequests+1 {
		t.Errorf("native and isolated mutually authenticated HTTP/2 requests = %d, want %d", got, isolatedRequests+1)
	}
}

func TestX509TransportDoesNotInheritNativeShapedHTTP2Handlers(t *testing.T) {
	fixture := newX509TransportFixture(t)
	server := newX509HTTP2TransportServer(t, fixture, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	template := fixture.transport(t, server)
	template.ForceAttemptHTTP2 = true
	var intercepted atomic.Int32
	template.TLSNextProto = map[string]func(string, *tls.Conn) http.RoundTripper{
		"h2": func(string, *tls.Conn) http.RoundTripper {
			intercepted.Add(1)
			return nil
		},
		"unencrypted_http2": func(string, *tls.Conn) http.RoundTripper {
			intercepted.Add(1)
			return nil
		},
	}
	capability := newX509Capability(t, template)
	request := x509TransportRequest(t, http.MethodGet, "https://"+x509TransportAPI+"/v1/models")
	response, err := capability.Do(request)
	if err != nil {
		t.Fatalf("dispatch without inheriting native-shaped caller protocol callbacks: %v", err)
	}
	if err := response.Body.Close(); err != nil {
		t.Fatalf("close isolated native-shaped HTTP/2 response: %v", err)
	}
	if got := intercepted.Load(); got != 0 {
		t.Errorf("caller-owned native-shaped HTTP/2 handlers intercepted %d requests", got)
	}
}

func TestX509TransportDoesNotInheritLaterCallerTLSProtocolHandlers(t *testing.T) {
	fixture := newX509TransportFixture(t)
	var received, intercepted atomic.Int32
	server := newX509HTTP2TransportServer(t, fixture, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		received.Add(1)
		w.WriteHeader(http.StatusOK)
	}))
	template := fixture.transport(t, server)
	template.ForceAttemptHTTP2 = true
	capability := newX509Capability(t, template)
	template.TLSNextProto = map[string]func(string, *tls.Conn) http.RoundTripper{
		"attacker": func(string, *tls.Conn) http.RoundTripper {
			intercepted.Add(1)
			return nil
		},
	}
	request := x509TransportRequest(t, http.MethodGet, "https://"+x509TransportAPI+"/v1/models")
	response, err := capability.Do(request)
	if err != nil {
		t.Fatalf("dispatch after an isolated caller TLS protocol handler was installed: %v", err)
	}
	if err := response.Body.Close(); err != nil {
		t.Fatalf("close isolated caller protocol response: %v", err)
	}
	if got := intercepted.Load(); got != 0 {
		t.Errorf("caller-owned post-attestation TLS protocol handler intercepted %d requests", got)
	}
	if got := received.Load(); got != 1 {
		t.Errorf("isolated post-attestation TLS protocol request reached the API %d times", got)
	}
}

func TestX509TransportRejectsHTTP2ErrorCallbacks(t *testing.T) {
	fixture := newX509TransportFixture(t)
	template := fixture.transport(t, nil)
	var invoked atomic.Int32
	template.HTTP2 = &http.HTTP2Config{CountError: func(string) { invoked.Add(1) }}
	if capability, err := auth.NewX509Transport(template); capability != nil ||
		err == nil || !strings.Contains(err.Error(), "HTTP/2 error callbacks") {
		t.Fatalf("caller-controlled HTTP/2 error hook: capability = %v, error = %v", capability, err)
	}
	if got := invoked.Load(); got != 0 {
		t.Errorf("rejected caller-controlled HTTP/2 error callback ran %d times", got)
	}
}

func newX509HTTP2TransportServer(t *testing.T, fixture *x509TransportFixture, handler http.Handler) *httptest.Server {
	t.Helper()
	server := httptest.NewUnstartedServer(handler)
	server.EnableHTTP2 = true
	server.TLS = &tls.Config{
		MinVersion:   tls.VersionTLS12,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    fixture.trust,
		Certificates: []tls.Certificate{fixture.certificate(t, "synthetic HTTP/2 server", []string{x509TransportAPI}, false)},
	}
	server.Config.ErrorLog = log.New(io.Discard, "", 0)
	server.StartTLS()
	t.Cleanup(server.Close)
	return server
}

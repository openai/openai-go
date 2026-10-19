package openai_test

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httptrace"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	wire "github.com/coder/websocket"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/auth"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

func newX509WebSocketConfig(t *testing.T, issuerHandler, apiHandler http.HandlerFunc) auth.X509WorkloadIdentity {
	t.Helper()
	t.Setenv("OPENAI_BASE_URL", "https://mtls.api.openai.com/v1/")
	lab := newX509ConformanceLab(t)
	routes := make(map[string]string)
	for host, handler := range map[string]http.HandlerFunc{x509ConformanceIssuerHost: issuerHandler, x509ConformanceAPIHost: apiHandler} {
		server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.TLS == nil || len(r.TLS.VerifiedChains) == 0 || r.TLS.PeerCertificates[0].Subject.CommonName != "websocket-workload" || r.TLS.ServerName != host {
				t.Error("X.509 WebSocket request did not preserve verified workload identity and TLS origin")
				w.WriteHeader(http.StatusForbidden)
				return
			}
			handler(w, r)
		}))
		leaf, key := lab.issue(t, host, lab.intermediate, lab.issuerKey, false, []string{host})
		server.TLS = &tls.Config{MinVersion: tls.VersionTLS12, ClientAuth: tls.RequireAndVerifyClientCert, ClientCAs: lab.trust,
			Certificates: []tls.Certificate{{Certificate: [][]byte{leaf.Raw, lab.intermediate.Raw}, PrivateKey: key}}}
		server.EnableHTTP2 = true
		server.Config.ErrorLog = log.New(io.Discard, "", 0)
		server.StartTLS()
		t.Cleanup(server.Close)
		routes[host+":443"] = server.Listener.Addr().String()
	}
	template := lab.transport(t, routes, lab.identity(t, "websocket-workload", true))
	template.ForceAttemptHTTP2 = true
	transport, err := auth.NewX509Transport(template)
	if err != nil {
		t.Fatalf("attest WebSocket workload transport: %v", err)
	}
	t.Cleanup(func() { _ = transport.Close() })
	return auth.X509WorkloadIdentity{IdentityProviderID: "synthetic-provider", ServiceAccountID: "synthetic-account", Transport: transport}
}

func x509WebSocketIssuer(t *testing.T, calls *atomic.Int32) http.HandlerFunc {
	t.Helper()
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/oauth/token" || r.Method != http.MethodPost || r.Header.Get("Authorization") != "" || r.Header.Get("X-Custom-Handshake") != "" || r.Header.Get("Upgrade") != "" {
			t.Error("WebSocket token exchange leaked API headers or changed the pinned endpoint")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var fields map[string]any
		if err := json.NewDecoder(r.Body).Decode(&fields); err != nil || len(fields) != 4 {
			t.Error("WebSocket token exchange did not contain the expected identity fields")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, _ = fmt.Fprintf(w, `{"access_token":"synthetic-websocket-token-%d","token_type":"Bearer","issued_token_type":"urn:ietf:params:oauth:token-type:access_token","expires_in":3600}`, calls.Add(1))
	}
}

func TestX509WebSocketPublicClientHeadersMessagesAndLifetime(t *testing.T) {
	var exchanges, upgrades atomic.Int32
	closed := make(chan struct{})
	config := newX509WebSocketConfig(t, x509WebSocketIssuer(t, &exchanges), func(w http.ResponseWriter, r *http.Request) {
		upgrades.Add(1)
		if r.ProtoMajor != 1 || r.URL.Path != "/v1/responses" || r.Header.Get("Authorization") != "Bearer synthetic-websocket-token-1" || r.Header.Get("X-Custom-Handshake") != "synthetic-value" {
			t.Error("X.509 WebSocket handshake lost its origin, bearer or custom header")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		socket, err := wire.Accept(w, r, nil)
		if err != nil {
			t.Errorf("accept X.509 WebSocket: %v", err)
			return
		}
		defer close(closed)
		defer func() { _ = socket.CloseNow() }()
		for range 2 {
			_, payload, err := socket.Read(r.Context())
			if err != nil {
				t.Errorf("read authenticated WebSocket command: %v", err)
				return
			}
			var command map[string]any
			if err := json.Unmarshal(payload, &command); err != nil || command["type"] != "response.create" {
				t.Error("authenticated WebSocket received an invalid response.create command")
				return
			}
			if err := socket.Write(r.Context(), wire.MessageText, []byte(`{"type":"response.completed","response":{"id":"synthetic-response","status":"completed","output":[]}}`)); err != nil {
				t.Errorf("write authenticated WebSocket response: %v", err)
				return
			}
		}
		_, _, _ = socket.Read(r.Context())
	})
	client := openai.NewClient(option.WithX509WorkloadIdentity(config), option.WithMaxRetries(0))
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	opening, stopOpening := context.WithCancel(ctx)
	var response *http.Response
	connection, err := client.Responses.Connect(opening, responses.ResponseConnectionOptions{}, option.WithHeader("X-Custom-Handshake", "synthetic-value"), option.WithResponseInto(&response))
	if err != nil {
		t.Fatalf("Connect with attested X.509 transport: %v", err)
	}
	defer connection.Abort()
	stopOpening()
	if response == nil || response.StatusCode != http.StatusSwitchingProtocols || response.Request.Header.Get("Authorization") != "" {
		t.Error("X.509 opening response did not preserve status and redact the bearer")
	}
	// Closing the capability rejects new requests but preserves active streams,
	// matching its ordinary HTTP request lifecycle.
	if err := config.Transport.Close(); err != nil {
		t.Fatalf("close X.509 capability: %v", err)
	}
	for range 2 {
		if err := connection.Create(ctx, responses.ResponsesClientEventResponseCreateParam{}); err != nil {
			t.Fatalf("write after opening context/capability close: %v", err)
		}
		result, err := connection.FinalResponse(ctx)
		if err != nil || result.ID != "synthetic-response" {
			t.Fatalf("read authenticated final response: result=%v error=%v", result, err)
		}
	}
	if next, err := client.Responses.Connect(ctx, responses.ResponseConnectionOptions{}); err == nil {
		next.Abort()
		t.Error("closed X.509 capability accepted another connection")
	}
	if err := connection.Close(); err != nil {
		t.Fatalf("close authenticated WebSocket: %v", err)
	}
	select {
	case <-closed:
	case <-ctx.Done():
		t.Fatal("WebSocket Close did not release the authenticated stream")
	}
	if exchanges.Load() != 1 || upgrades.Load() != 1 {
		t.Errorf("two commands used exchanges/upgrades=%d/%d, want 1/1", exchanges.Load(), upgrades.Load())
	}
}

func TestX509WebSocketRejectedBearerUsesOpeningBudget(t *testing.T) {
	for _, retries := range []int{0, 1} {
		t.Run(fmt.Sprint(retries), func(t *testing.T) {
			var exchanges, upgrades atomic.Int32
			config := newX509WebSocketConfig(t, x509WebSocketIssuer(t, &exchanges), func(w http.ResponseWriter, r *http.Request) {
				upgrades.Add(1)
				if r.Header.Get("Authorization") != "Bearer synthetic-websocket-token-2" {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				socket, err := wire.Accept(w, r, nil)
				if err != nil {
					t.Errorf("accept refreshed X.509 WebSocket: %v", err)
					return
				}
				defer func() { _ = socket.CloseNow() }()
				_, _, _ = socket.Read(r.Context())
			})
			client := openai.NewClient(option.WithX509WorkloadIdentity(config), option.WithMaxRetries(retries))
			ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
			defer cancel()
			connection, err := client.Responses.Connect(ctx, responses.ResponseConnectionOptions{})
			if connection != nil {
				defer connection.Abort()
			}
			if (err == nil) != (retries == 1) {
				t.Errorf("Connect(retries=%d) error=%v", retries, err)
			}
			if exchanges.Load() != int32(retries+1) || upgrades.Load() != int32(retries+1) {
				t.Errorf("Connect(retries=%d) exchanges/upgrades=%d/%d, want %d/%d", retries, exchanges.Load(), upgrades.Load(), retries+1, retries+1)
			}
		})
	}
}

func TestX509WebSocketRejectsUnsafeRequestsBeforeExchange(t *testing.T) {
	for _, test := range []struct {
		name   string
		mutate func(*http.Request)
	}{
		{"different endpoint", func(r *http.Request) { r.URL.Path = "/v1/models" }},
		{"issuer endpoint", func(r *http.Request) { r.URL.Host = x509ConformanceIssuerHost; r.URL.Path = "/oauth/token" }},
		{"different origin", func(r *http.Request) { r.URL.Host = "attacker.example.test" }},
		{"different protocol", func(r *http.Request) { r.Header.Set("Upgrade", "attacker") }},
		{"ambiguous connection", func(r *http.Request) { r.Header["connection"] = []string{"attacker"} }},
		{"invalid key", func(r *http.Request) { r.Header.Set("Sec-WebSocket-Key", "invalid") }},
		{"framing header", func(r *http.Request) { r.Header.Set("Content-Length", "5") }},
		{"alternate credentials", func(r *http.Request) { r.Header.Set("Cookie", "synthetic-cookie") }},
		{"request body", func(r *http.Request) { r.Body = io.NopCloser(strings.NewReader("synthetic")) }},
		{"trace callback", func(r *http.Request) {
			*r = *r.WithContext(httptrace.WithClientTrace(r.Context(), &httptrace.ClientTrace{}))
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			var exchanges, upgrades atomic.Int32
			config := newX509WebSocketConfig(t, x509WebSocketIssuer(t, &exchanges), func(w http.ResponseWriter, _ *http.Request) { upgrades.Add(1); w.WriteHeader(http.StatusBadRequest) })
			client := openai.NewClient(option.WithX509WorkloadIdentity(config), option.WithMaxRetries(0), option.WithMiddleware(func(r *http.Request, next option.MiddlewareNext) (*http.Response, error) {
				test.mutate(r)
				return next(r)
			}))
			connection, err := client.Responses.Connect(t.Context(), responses.ResponseConnectionOptions{})
			if connection != nil {
				connection.Abort()
			}
			if err == nil {
				t.Error("unsafe X.509 upgrade was accepted")
			}
			if exchanges.Load() != 0 || upgrades.Load() != 0 {
				t.Errorf("unsafe upgrade made exchanges/upgrades=%d/%d, want 0/0", exchanges.Load(), upgrades.Load())
			}
		})
	}
}

func TestX509WebSocketCancellationStopsOpening(t *testing.T) {
	started, stopped := make(chan struct{}), make(chan struct{})
	var exchanges atomic.Int32
	config := newX509WebSocketConfig(t, x509WebSocketIssuer(t, &exchanges), func(_ http.ResponseWriter, r *http.Request) {
		close(started)
		<-r.Context().Done()
		close(stopped)
	})
	client := openai.NewClient(option.WithX509WorkloadIdentity(config), option.WithMaxRetries(0))
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	result := make(chan error, 1)
	go func() {
		connection, err := client.Responses.Connect(ctx, responses.ResponseConnectionOptions{})
		if connection != nil {
			connection.Abort()
		}
		result <- err
	}()
	select {
	case <-started:
	case <-time.After(5 * time.Second):
		t.Fatal("X.509 WebSocket opening did not reach the server")
	}
	cancel()
	select {
	case err := <-result:
		if !errors.Is(err, context.Canceled) {
			t.Errorf("canceled X.509 opening error=%v, want context.Canceled", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("X.509 WebSocket opening ignored cancellation")
	}
	select {
	case <-stopped:
	case <-time.After(5 * time.Second):
		t.Fatal("X.509 WebSocket cancellation did not release the opening request")
	}
}

func TestX509WebSocketRejectsRedirectsAndRedactsErrors(t *testing.T) {
	for _, destination := range []string{"https://mtls.api.openai.com/v1/responses?secret=synthetic-sensitive", "https://attacker.example.test/?secret=synthetic-sensitive"} {
		t.Run(destination, func(t *testing.T) {
			var exchanges, upgrades atomic.Int32
			config := newX509WebSocketConfig(t, x509WebSocketIssuer(t, &exchanges), func(w http.ResponseWriter, r *http.Request) {
				upgrades.Add(1)
				http.Redirect(w, r, destination, http.StatusTemporaryRedirect)
			})
			client := openai.NewClient(option.WithX509WorkloadIdentity(config), option.WithMaxRetries(2))
			connection, err := client.Responses.Connect(t.Context(), responses.ResponseConnectionOptions{})
			if connection != nil {
				connection.Abort()
			}
			if err == nil {
				t.Fatal("X.509 WebSocket followed a redirect")
			}
			if strings.Contains(err.Error(), "synthetic-sensitive") || strings.Contains(err.Error(), "synthetic-websocket-token") {
				t.Error("X.509 WebSocket error revealed redirect metadata or credentials")
			}
			if exchanges.Load() != 1 || upgrades.Load() != 1 {
				t.Errorf("redirect made exchanges/upgrades=%d/%d, want 1/1", exchanges.Load(), upgrades.Load())
			}
		})
	}
}

func TestX509WebSocketAdapterDoesNotRelaxOrdinaryHTTP(t *testing.T) {
	var exchanges, upgrades atomic.Int32
	config := newX509WebSocketConfig(t, x509WebSocketIssuer(t, &exchanges), func(w http.ResponseWriter, _ *http.Request) {
		upgrades.Add(1)
		w.WriteHeader(http.StatusBadRequest)
	})
	client := openai.NewClient(option.WithX509WorkloadIdentity(config), option.WithMaxRetries(0))
	if err := client.Get(t.Context(), "responses", nil, option.WithHeader("Connection", "Upgrade"), option.WithHeader("Upgrade", "websocket")); err == nil {
		t.Error("ordinary HTTP request accepted protocol upgrade headers")
	}
	adapter, ok := any(config.Transport).(interface{ WebSocketHTTPClient() *http.Client })
	if !ok {
		t.Fatal("attested transport does not expose a WebSocket adapter")
	}
	request, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "https://mtls.api.openai.com/v1/responses", nil)
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Connection", "Upgrade")
	request.Header.Set("Upgrade", "websocket")
	request.Header.Set("Sec-WebSocket-Version", "13")
	request.Header.Set("Sec-WebSocket-Key", "c3ludGhldGljLW5vbmNlIQ==")
	response, err := adapter.WebSocketHTTPClient().Do(request)
	if response != nil {
		_ = response.Body.Close()
	}
	if err == nil {
		t.Error("direct adapter call accepted an upgrade outside the SDK opening flow")
	}
	if exchanges.Load() != 0 || upgrades.Load() != 0 {
		t.Errorf("non-SDK upgrade made exchanges/upgrades=%d/%d, want 0/0", exchanges.Load(), upgrades.Load())
	}
}

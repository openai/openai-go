package requestconfig_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	wire "github.com/coder/websocket"
	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
	transport "github.com/openai/openai-go/v3/packages/websocket"
)

func dialPrepared(t *testing.T, opts ...option.RequestOption) (*transport.Connection[string], error) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cfg, err := requestconfig.NewRequestConfig(ctx, http.MethodGet, "responses", nil, nil, append([]option.RequestOption{requestconfig.WithBearerAuthSecurity()}, opts...)...)
	if err != nil {
		return nil, err
	}
	req, client, err := cfg.PrepareWebSocket()
	if err != nil {
		return nil, err
	}
	return transport.Dial(req.Context(), req, client, transport.Options{}, func(data []byte) (string, error) { return string(data), nil }, func(string) string { return "" })
}

func TestWebSocketUpgradePreservesRequestPolicyAndTLS(t *testing.T) {
	done := make(chan struct{})
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer close(done)
		if r.URL.Path != "/custom/v1/responses" || r.URL.Query().Get("base") != "yes" || r.URL.Query().Get("extra") != "present" {
			t.Errorf("upgrade URL = %s", r.URL)
		}
		if r.Header.Get("Authorization") != "Bearer fake-key" || r.Header.Get("OpenAI-Organization") != "org-test" || r.Header.Get("OpenAI-Project") != "proj-test" || r.Header.Get("X-Middleware") != "applied" {
			t.Error("upgrade did not preserve configured headers")
		}
		socket, err := wire.Accept(w, r, nil)
		if err != nil {
			t.Errorf("accept: %v", err)
			return
		}
		defer func() { _ = socket.CloseNow() }()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		_, _, err = socket.Read(ctx)
		if err != nil {
			t.Errorf("read: %v", err)
			return
		}
		_ = socket.Write(ctx, wire.MessageText, []byte("response.completed"))
		_, _, _ = socket.Read(ctx)
	}))
	defer server.Close()
	client := server.Client()
	client.Timeout = 100 * time.Millisecond
	connection, err := dialPrepared(t, option.WithBaseURL(server.URL+"/custom/v1?base=yes"), option.WithAPIKey("fake-key"), option.WithOrganization("org-test"), option.WithProject("proj-test"), option.WithQuery("extra", "present"), option.WithHTTPClient(client), option.WithMiddleware(func(r *http.Request, next option.MiddlewareNext) (*http.Response, error) {
		r.Header.Set("X-Middleware", "applied")
		return next(r)
	}))
	if err != nil {
		t.Fatal(err)
	}
	defer connection.Abort()
	// The configured HTTP timeout must stop applying after the upgrade.
	time.Sleep(150 * time.Millisecond)
	if err := connection.Send(context.Background(), []byte(`{}`)); err != nil {
		t.Fatal(err)
	}
	if event, err := connection.Recv(context.Background()); err != nil || event != "response.completed" {
		t.Fatalf("receive = %q,%v", event, err)
	}
	_ = connection.Close()
	<-done
	if client.Timeout != 100*time.Millisecond {
		t.Error("caller-owned HTTP client was mutated")
	}
}

func TestWebSocketRejectsCrossOriginRedirect(t *testing.T) {
	var leaked atomic.Int32
	other := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { leaked.Add(1); w.WriteHeader(http.StatusBadRequest) }))
	defer other.Close()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, other.URL, http.StatusTemporaryRedirect)
	}))
	defer server.Close()
	_, err := dialPrepared(t, option.WithBaseURL(server.URL), option.WithAPIKey("fake-key"))
	if err == nil {
		t.Fatal("cross-origin redirect unexpectedly succeeded")
	}
	if leaked.Load() != 0 {
		t.Fatalf("cross-origin requests = %d, want 0", leaked.Load())
	}
}

type unavailableHTTPDoer struct{ calls atomic.Int32 }

func (c *unavailableHTTPDoer) Do(*http.Request) (*http.Response, error) {
	c.calls.Add(1)
	return nil, fmt.Errorf("unexpected request")
}

func TestWebSocketRejectsUnsupportedHTTPClientBeforeNetwork(t *testing.T) {
	client := &unavailableHTTPDoer{}
	_, err := dialPrepared(t, option.WithBaseURL("https://example.invalid/v1"), option.WithHTTPClient(client))
	if err == nil || !strings.Contains(err.Error(), "WebSocketHTTPClient") {
		t.Fatalf("custom HTTP client error = %v", err)
	}
	if client.calls.Load() != 0 {
		t.Error("unsupported HTTP client was called")
	}
}

func TestWebSocketRejectsOriginChangeFromMiddleware(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1); w.WriteHeader(http.StatusBadRequest) }))
	defer server.Close()
	_, err := dialPrepared(t, option.WithBaseURL(server.URL), option.WithMiddleware(func(r *http.Request, next option.MiddlewareNext) (*http.Response, error) {
		r.URL.Host = "example.invalid"
		return next(r)
	}))
	if err == nil {
		t.Fatal("middleware origin change unexpectedly succeeded")
	}
	if calls.Load() != 0 {
		t.Error("origin-changing request reached the transport")
	}
}

type capableHTTPDoer struct {
	client *http.Client
	calls  atomic.Int32
}

func (c *capableHTTPDoer) Do(*http.Request) (*http.Response, error) {
	c.calls.Add(1)
	return nil, fmt.Errorf("unexpected HTTP request")
}
func (c *capableHTTPDoer) WebSocketHTTPClient() *http.Client { return c.client }

func TestWebSocketCustomCapabilityUsesConfiguredProxy(t *testing.T) {
	var calls atomic.Int32
	proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		if r.URL.Host != "responses.example" || r.URL.Path != "/v1/responses" {
			t.Errorf("proxy target = %s", r.URL)
		}
		socket, err := wire.Accept(w, r, nil)
		if err != nil {
			t.Errorf("accept: %v", err)
			return
		}
		defer func() { _ = socket.CloseNow() }()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		_ = socket.Write(ctx, wire.MessageText, []byte("proxy-used"))
		_, _, _ = socket.Read(ctx)
	}))
	defer proxy.Close()
	proxyURL, _ := url.Parse(proxy.URL)
	httpTransport := http.DefaultTransport.(*http.Transport).Clone()
	httpTransport.Proxy = http.ProxyURL(proxyURL)
	defer httpTransport.CloseIdleConnections()
	capable := &capableHTTPDoer{client: &http.Client{Transport: httpTransport}}
	connection, err := dialPrepared(t, option.WithBaseURL("http://responses.example/v1"), option.WithHTTPClient(capable))
	if err != nil {
		t.Fatal(err)
	}
	defer connection.Abort()
	if event, err := connection.Recv(context.Background()); err != nil || event != "proxy-used" {
		t.Fatalf("proxy receive = %q,%v", event, err)
	}
	if calls.Load() != 1 || capable.calls.Load() != 0 {
		t.Errorf("proxy calls = %d, HTTP Do calls = %d", calls.Load(), capable.calls.Load())
	}
}

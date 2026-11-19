package openai_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"

	wire "github.com/coder/websocket"
	"github.com/openai/openai-go/v3/responses"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/auth"
	"github.com/openai/openai-go/v3/bedrock"
	"github.com/openai/openai-go/v3/option"
)

func TestCredentialsRejectRemotePlaintextEndpoint(t *testing.T) {
	for _, source := range []string{"client", "environment", "request"} {
		for _, credential := range []string{"api", "admin", "workload"} {
			t.Run(source+"/"+credential, func(t *testing.T) {
				t.Setenv("OPENAI_API_KEY", "")
				t.Setenv("OPENAI_ADMIN_KEY", "")
				t.Setenv("OPENAI_BASE_URL", "https://api.openai.com/v1")
				provider := &mockSubjectTokenProvider{token: "synthetic-subject", tokenType: auth.SubjectTokenTypeJWT}
				apiCalls, exchangeCalls, credentialCalls := 0, 0, 0
				transport := originTestRoundTripper(func(req *http.Request) (*http.Response, error) {
					body := `{"data":[],"object":"list"}`
					if req.URL.String() == auth.TokenExchangeURL {
						exchangeCalls++
						body = `{"access_token":"synthetic-access","expires_in":3600}`
					} else {
						apiCalls++
						if req.Header.Get("Authorization") != "" {
							credentialCalls++
						}
					}
					return &http.Response{StatusCode: http.StatusOK, Header: http.Header{"Content-Type": {"application/json"}}, Body: io.NopCloser(strings.NewReader(body)), Request: req}, nil
				})
				opts := []option.RequestOption{
					option.WithHTTPClient(&http.Client{Transport: transport}),
					option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
						t.Error("blocked credential reached middleware")
						return next(req)
					}),
				}
				switch credential {
				case "api":
					opts = append(opts, option.WithAPIKey("synthetic-api"))
				case "admin":
					opts = append(opts, option.WithAdminAPIKey("synthetic-admin"))
				case "workload":
					opts = append(opts, option.WithWorkloadIdentity(testWorkloadIdentity(provider)))
				}
				const endpoint = "http://remote.invalid/v1"
				var requestOpts []option.RequestOption
				switch source {
				case "client":
					opts = append(opts, option.WithBaseURL(endpoint))
				case "environment":
					t.Setenv("OPENAI_BASE_URL", endpoint)
				case "request":
					requestOpts = append(requestOpts, option.WithBaseURL(endpoint))
				}
				client := openai.NewClient(opts...)
				var response map[string]any
				err := client.Get(t.Context(), "models", nil, &response, requestOpts...)
				if err == nil || !strings.Contains(err.Error(), "HTTPS") {
					t.Errorf("expected HTTPS policy error, got %v", err)
				}
				if apiCalls != 0 || credentialCalls != 0 || exchangeCalls != 0 || provider.GetCallCount() != 0 {
					t.Errorf("blocked request made API=%d credential-bearing=%d exchange=%d provider=%d calls", apiCalls, credentialCalls, exchangeCalls, provider.GetCallCount())
				}
			})
		}
	}
}

func TestPlaintextCredentialRejectionClosesUpload(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("OPENAI_ADMIN_KEY", "")
	for _, credential := range []string{"api", "admin", "workload"} {
		t.Run(credential, func(t *testing.T) {
			body := &plaintextEndpointBody{Reader: strings.NewReader("synthetic-body")}
			provider := &mockSubjectTokenProvider{token: "synthetic-subject", tokenType: auth.SubjectTokenTypeJWT}
			opts := []option.RequestOption{option.WithBaseURL("http://remote.invalid/v1"), option.WithMaxRetries(0)}
			switch credential {
			case "api":
				opts = append(opts, option.WithAPIKey("synthetic-api"))
			case "admin":
				opts = append(opts, option.WithAdminAPIKey("synthetic-admin"))
			case "workload":
				opts = append(opts, option.WithWorkloadIdentity(testWorkloadIdentity(provider)))
			}
			calls := 0
			opts = append(opts, option.WithHTTPClient(plaintextEndpointDoer(func(*http.Request) (*http.Response, error) {
				calls++
				return nil, io.EOF
			})))
			client := openai.NewClient(opts...)
			err := client.Post(t.Context(), "models", body, nil)
			if err == nil || !strings.Contains(err.Error(), "HTTPS") || calls != 0 || provider.GetCallCount() != 0 || body.closes != 1 {
				t.Fatalf("error=%v dispatches=%d provider=%d closes=%d", err, calls, provider.GetCallCount(), body.closes)
			}
		})
	}
}

func TestRequestOptionCanRemoveStaticCredentials(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("OPENAI_ADMIN_KEY", "")
	for _, credential := range []option.RequestOption{option.WithAPIKey("synthetic-api"), option.WithAdminAPIKey("synthetic-admin")} {
		calls := 0
		client := openai.NewClient(credential, option.WithBaseURL("http://remote.invalid/v1"), option.WithMaxRetries(0),
			option.WithHTTPClient(plaintextEndpointDoer(func(req *http.Request) (*http.Response, error) {
				calls++
				if req.Header.Get("Authorization") != "" {
					t.Error("removed credential reached transport")
				}
				return &http.Response{StatusCode: http.StatusOK, Header: http.Header{}, Body: http.NoBody, Request: req}, nil
			})))
		if err := client.Get(t.Context(), "models", nil, nil, option.WithHeaderDel("Authorization")); err != nil || calls != 1 {
			t.Fatalf("error=%v dispatches=%d", err, calls)
		}
	}
}

func TestLaterStaticKeyRequiresSecureEndpoint(t *testing.T) {
	for _, key := range []option.RequestOption{option.WithAPIKey("synthetic"), option.WithAdminAPIKey("synthetic")} {
		client := openai.NewClient(option.WithBaseURL("http://remote.invalid"), option.WithHeaderDel("Authorization"),
			option.WithMiddleware(func(*http.Request, option.MiddlewareNext) (*http.Response, error) {
				t.Error("blocked request reached middleware")
				return nil, io.EOF
			}))
		if err := client.Get(t.Context(), "models", nil, nil, key); err == nil || !strings.Contains(err.Error(), "HTTPS") {
			t.Fatalf("expected HTTPS policy error, got %v", err)
		}
	}
}

func TestUnsupportedSecureWebSocketClientRejectsBeforeWorkloadAuth(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("OPENAI_ADMIN_KEY", "")
	provider := &mockSubjectTokenProvider{token: "synthetic", tokenType: auth.SubjectTokenTypeJWT}
	client := openai.NewClient(option.WithBaseURL("https://service.invalid"), option.WithWorkloadIdentity(testWorkloadIdentity(provider)),
		option.WithHTTPClient(plaintextEndpointDoer(func(*http.Request) (*http.Response, error) {
			t.Error("unsupported WebSocket client received a request")
			return nil, io.EOF
		})))
	_, err := client.Responses.Connect(t.Context(), responses.ResponseConnectionOptions{})
	if err == nil || !strings.Contains(err.Error(), "WebSocketHTTPClient") || provider.GetCallCount() != 0 {
		t.Fatalf("adapter error=%v provider calls=%d", err, provider.GetCallCount())
	}
}

func TestURLCredentialsRequireSecureTransport(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("OPENAI_ADMIN_KEY", "")
	for _, scheme := range []string{"http", "https"} {
		t.Run(scheme, func(t *testing.T) {
			calls := 0
			client := openai.NewClient(option.WithBaseURL(scheme+"://synthetic:password@remote.invalid/v1"), option.WithMaxRetries(0),
				option.WithHTTPClient(&http.Client{Transport: originTestRoundTripper(func(req *http.Request) (*http.Response, error) {
					calls++
					if _, _, ok := req.BasicAuth(); !ok {
						t.Error("HTTPS URL credential was not preserved")
					}
					return &http.Response{StatusCode: http.StatusOK, Header: http.Header{}, Body: http.NoBody, Request: req}, nil
				})}))
			err := client.Get(t.Context(), "models", nil, nil)
			if scheme == "http" {
				if err == nil || !strings.Contains(err.Error(), "HTTPS") || calls != 0 {
					t.Fatalf("expected rejection before dispatch: error=%v calls=%d", err, calls)
				}
			} else if err != nil || calls != 1 {
				t.Fatalf("HTTPS error=%v calls=%d", err, calls)
			}
		})
	}
}

type plaintextWebSocketClient struct{ *http.Client }

func (c plaintextWebSocketClient) WebSocketHTTPClient() *http.Client { return c.Client }

func TestWebSocketMiddlewareCanRemoveStaticCredentials(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("OPENAI_ADMIN_KEY", "")
	var received atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Header.Get("Authorization") != "" {
			t.Error("removed credential reached transport")
		}
		received.Add(1)
		socket, err := wire.Accept(w, req, nil)
		if err != nil {
			t.Error(err)
			return
		}
		defer func() { _ = socket.CloseNow() }()
		_, _, _ = socket.Read(req.Context())
	}))
	defer server.Close()
	adapterCalls := 0
	adapter := plaintextWebSocketClient{&http.Client{Transport: originTestRoundTripper(func(req *http.Request) (*http.Response, error) {
		adapterCalls++
		return http.DefaultTransport.RoundTrip(req)
	})}}
	client := openai.NewClient(option.WithAPIKey("synthetic-api"), option.WithBaseURL(server.URL), option.WithUnsafeAllowHTTP(), option.WithHTTPClient(adapter),
		option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
			req.Header.Del("Authorization")
			return next(req)
		}))
	socket, err := client.Responses.Connect(t.Context(), responses.ResponseConnectionOptions{})
	if socket != nil {
		socket.Abort()
	}
	if err != nil || received.Load() != 1 || adapterCalls != 1 {
		t.Fatalf("error=%v dispatches=%d adapter calls=%d", err, received.Load(), adapterCalls)
	}
	body := &plaintextEndpointBody{Reader: strings.NewReader("fixture")}
	socket, err = client.Responses.Connect(t.Context(), responses.ResponseConnectionOptions{},
		option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
			req.Body = body
			return next(req)
		}),
		option.WithHTTPClient(plaintextEndpointDoer(func(*http.Request) (*http.Response, error) { return nil, io.EOF })))
	if socket != nil {
		socket.Abort()
	}
	if err == nil || !strings.Contains(err.Error(), "WebSocketHTTPClient") || received.Load() != 1 {
		t.Fatalf("unsupported client error=%v dispatches=%d", err, received.Load())
	}
	if body.closes != 1 {
		t.Errorf("unsupported client closed the middleware body %d times, want 1", body.closes)
	}
}

func TestUnsafeHTTPRejectsNonLoopbackEndpoints(t *testing.T) {
	for _, endpoint := range []string{
		"http://remote.invalid", "http://192.0.2.1", "http://[2001:db8::1]",
		"http://localhost.remote.invalid", "http://localhost.", "http://2130706433",
		"http://127.0.0.1@remote.invalid", "http://[::ffff:192.0.2.1]", "ftp://127.0.0.1",
	} {
		t.Run(endpoint, func(t *testing.T) {
			calls := 0
			client := openai.NewClient(option.WithAPIKey("synthetic-key"), option.WithBaseURL(endpoint), option.WithUnsafeAllowHTTP(),
				option.WithHTTPClient(&http.Client{Transport: originTestRoundTripper(func(req *http.Request) (*http.Response, error) {
					calls++
					return nil, io.EOF
				})}))
			_, err := client.Models.List(t.Context())
			if err == nil || !strings.Contains(err.Error(), "HTTPS") || calls != 0 {
				t.Fatalf("error=%v transport calls=%d; want HTTPS rejection without transport", err, calls)
			}
		})
	}
}

type plaintextEndpointDoer func(*http.Request) (*http.Response, error)

func (f plaintextEndpointDoer) Do(req *http.Request) (*http.Response, error) { return f(req) }

func TestLocalHTTPRequiresOptInAndDirectConnection(t *testing.T) {
	for _, host := range []string{"127.0.0.1", "localhost", "LOCALHOST", "LocalHost", "[::ffff:127.0.0.1]"} {
		for _, credential := range []string{"api", "admin", "workload", "middleware"} {
			t.Run(host+"/"+credential, func(t *testing.T) {
				t.Setenv("OPENAI_API_KEY", "")
				t.Setenv("OPENAI_ADMIN_KEY", "")
				var received, connections atomic.Int32
				server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.Header.Get("Authorization") != "Bearer synthetic-"+credential {
						t.Error("missing expected synthetic credential")
					}
					received.Add(1)
					w.Header().Set("Content-Type", "application/json")
					_, _ = io.WriteString(w, `{"data":[],"object":"list"}`)
				}))
				server.Config.ConnState = func(_ net.Conn, state http.ConnState) {
					if state == http.StateNew {
						connections.Add(1)
					}
				}
				server.Start()
				defer server.Close()
				u, err := url.Parse(server.URL)
				if err != nil {
					t.Fatal(err)
				}
				endpoint := "http://" + host + ":" + u.Port()
				provider := &mockSubjectTokenProvider{token: "synthetic-subject", tokenType: auth.SubjectTokenTypeJWT}
				opts := []option.RequestOption{option.WithBaseURL(endpoint)}
				switch credential {
				case "middleware":
					opts = append(opts, option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
						req.Header.Set("Authorization", "Bearer synthetic-middleware")
						return next(req)
					}))
				case "api":
					opts = append(opts, option.WithAPIKey("synthetic-api"))
				case "admin":
					opts = append(opts, option.WithAdminAPIKey("synthetic-admin"))
				case "workload":
					opts = append(opts, option.WithWorkloadIdentity(testWorkloadIdentity(provider)))
				}
				var response map[string]any
				client := openai.NewClient(opts...)
				err = client.Get(t.Context(), "models", nil, &response)
				if err == nil || received.Load() != 0 || provider.GetCallCount() != 0 {
					t.Fatalf("without opt-in: error=%v received=%d provider=%d", err, received.Load(), provider.GetCallCount())
				}
				var bypassCalls, exchanges atomic.Int32
				do := func(req *http.Request) (*http.Response, error) {
					if req.URL.String() == auth.TokenExchangeURL {
						exchanges.Add(1)
						return &http.Response{StatusCode: http.StatusOK, Header: http.Header{"Content-Type": {"application/json"}}, Body: io.NopCloser(strings.NewReader(`{"access_token":"synthetic-workload","expires_in":3600}`)), Request: req}, nil
					}
					bypassCalls.Add(1)
					return nil, io.EOF
				}
				native := &http.Client{Transport: &http.Transport{
					Proxy:       func(*http.Request) (*url.URL, error) { bypassCalls.Add(1); return u, nil },
					DialContext: func(context.Context, string, string) (net.Conn, error) { bypassCalls.Add(1); return nil, io.EOF },
				}}
				// A custom doer can handle workload exchanges but must never receive
				// the plaintext API request. Native proxy/dial hooks also stay unused.
				clients := []option.HTTPClient{native, plaintextEndpointDoer(do)}
				if credential == "workload" {
					clients = []option.HTTPClient{&http.Client{Transport: originTestRoundTripper(do)}, plaintextEndpointDoer(do)}
				}
				for _, httpClient := range clients {
					err = client.Get(t.Context(), "models", nil, &response, option.WithUnsafeAllowHTTP(), option.WithHTTPClient(httpClient))
					if err != nil {
						t.Fatal(err)
					}
				}
				if received.Load() != int32(len(clients)) || bypassCalls.Load() != 0 {
					t.Fatalf("received=%d routing hook calls=%d", received.Load(), bypassCalls.Load())
				}
				if connections.Load() != 1 {
					t.Errorf("successive opted-in calls opened %d connections, want one", connections.Load())
				}
				if credential == "workload" && exchanges.Load() == 0 {
					t.Error("workload credential was not exchanged")
				}
			})
		}
	}
}

func TestBedrockSkipAuthWithOpenAIWorkloadIdentity(t *testing.T) {
	var received atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer synthetic-workload" {
			t.Error("missing workload credential")
		}
		received.Add(1)
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"data":[],"object":"list"}`)
	}))
	defer server.Close()

	provider := &mockSubjectTokenProvider{token: "synthetic-subject", tokenType: auth.SubjectTokenTypeJWT}
	exchange := plaintextEndpointDoer(func(req *http.Request) (*http.Response, error) {
		if req.URL.String() != auth.TokenExchangeURL {
			return nil, errors.New("expected the direct loopback transport for API requests")
		}
		return &http.Response{StatusCode: http.StatusOK, Header: http.Header{"Content-Type": {"application/json"}}, Body: io.NopCloser(strings.NewReader(`{"access_token":"synthetic-workload","expires_in":3600}`)), Request: req}, nil
	})
	client, err := bedrock.NewClient(t.Context(), bedrock.Config{SkipAuth: true, BaseURL: server.URL},
		option.WithWorkloadIdentity(testWorkloadIdentity(provider)), option.WithHTTPClient(exchange), option.WithMaxRetries(0))
	if err != nil {
		t.Fatal(err)
	}

	if _, err := client.Models.List(t.Context()); err == nil || !strings.Contains(err.Error(), "HTTPS") {
		t.Fatalf("expected HTTPS policy error, got %v", err)
	}
	if provider.GetCallCount() != 0 || received.Load() != 0 {
		t.Fatal("denied request acquired a credential or reached the server")
	}
	if _, err := client.Models.List(t.Context(), option.WithUnsafeAllowHTTP()); err != nil {
		t.Fatal(err)
	}
	if received.Load() != 1 {
		t.Errorf("received %d authenticated requests, want 1", received.Load())
	}
}

func TestHTTPSPreservesCredentialTransport(t *testing.T) {
	for _, credential := range []string{"api", "admin", "workload", "middleware"} {
		t.Run(credential, func(t *testing.T) {
			t.Setenv("OPENAI_API_KEY", "")
			t.Setenv("OPENAI_ADMIN_KEY", "")
			calls := 0
			opts := []option.RequestOption{option.WithBaseURL("https://remote.invalid/v1"), option.WithUnsafeAllowHTTP()}
			switch credential {
			case "middleware":
				opts = append(opts, option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
					req.Header.Set("Authorization", "Bearer synthetic-middleware")
					return next(req)
				}))
			case "api":
				opts = append(opts, option.WithAPIKey("synthetic-api"))
			case "admin":
				opts = append(opts, option.WithAdminAPIKey("synthetic-admin"))
			case "workload":
				opts = append(opts, option.WithWorkloadIdentity(testWorkloadIdentity(&mockSubjectTokenProvider{token: "synthetic-subject", tokenType: auth.SubjectTokenTypeJWT})))
			}
			opts = append(opts, option.WithHTTPClient(plaintextEndpointDoer(func(req *http.Request) (*http.Response, error) {
				body := `{"data":[],"object":"list"}`
				if req.URL.String() == auth.TokenExchangeURL {
					body = `{"access_token":"synthetic-workload","expires_in":3600}`
				} else {
					calls++
					if req.URL.Scheme != "https" || req.Header.Get("Authorization") != "Bearer synthetic-"+credential {
						t.Error("HTTPS credential transport changed")
					}
				}
				return &http.Response{StatusCode: http.StatusOK, Header: http.Header{"Content-Type": {"application/json"}}, Body: io.NopCloser(strings.NewReader(body)), Request: req}, nil
			})))
			client := openai.NewClient(opts...)
			var response map[string]any
			if err := client.Get(t.Context(), "models", nil, &response); err != nil || calls != 1 {
				t.Fatalf("HTTPS error=%v calls=%d", err, calls)
			}
		})
	}
}

func TestPlaintextWithoutAuthorizationPreservesCustomClient(t *testing.T) {
	for _, endpoint := range []string{"http://remote.invalid/v1", "http://127.0.0.1:1/v1"} {
		for _, deletion := range []string{"client", "request", "empty-client", "empty-request"} {
			t.Run(endpoint+"/"+deletion, func(t *testing.T) {
				t.Setenv("OPENAI_API_KEY", "synthetic-api")
				t.Setenv("OPENAI_ADMIN_KEY", "synthetic-admin")
				calls := 0
				opts := []option.RequestOption{option.WithBaseURL(endpoint), option.WithUnsafeAllowHTTP(),
					option.WithHTTPClient(plaintextEndpointDoer(func(req *http.Request) (*http.Response, error) {
						calls++
						if req.Header.Get("Authorization") != "" {
							t.Error("deleted authorization was restored")
						}
						return &http.Response{StatusCode: http.StatusOK, Header: http.Header{"Content-Type": {"application/json"}}, Body: io.NopCloser(strings.NewReader(`{"data":[]}`)), Request: req}, nil
					}))}
				var requestOpts []option.RequestOption
				override := option.WithHeaderDel("Authorization")
				if strings.HasPrefix(deletion, "empty-") {
					override = option.WithHeader("Authorization", "")
				}
				if strings.HasSuffix(deletion, "client") {
					opts = append(opts, override)
				} else {
					requestOpts = append(requestOpts, override)
				}
				client := openai.NewClient(opts...)
				var response map[string]any
				if err := client.Get(t.Context(), "models", nil, &response, requestOpts...); err != nil || calls != 1 {
					t.Fatalf("unauthenticated error=%v calls=%d", err, calls)
				}
			})
		}
	}
}

type plaintextEndpointBody struct {
	io.Reader
	closes int
}

func (b *plaintextEndpointBody) Close() error { b.closes++; return nil }

func TestHTTPSMiddlewareErrorClosesReplacedBody(t *testing.T) {
	original := &plaintextEndpointBody{Reader: strings.NewReader("original")}
	replacement := &plaintextEndpointBody{Reader: strings.NewReader("replacement")}
	blocked := errors.New("middleware blocked request")
	client := openai.NewClient(option.WithAPIKey("synthetic-api"), option.WithBaseURL("https://remote.invalid/v1"), option.WithMaxRetries(0),
		option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
			req.Body = replacement
			return nil, blocked
		}))
	err := client.Post(t.Context(), "models", original, nil)
	if !errors.Is(err, blocked) || original.closes != 1 || replacement.closes != 1 {
		t.Fatalf("error=%v original closes=%d replacement closes=%d", err, original.closes, replacement.closes)
	}
}

func TestMiddlewareCredentialsRequireSecureTransport(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("OPENAI_ADMIN_KEY", "")
	for _, scheme := range []string{"http", "https"} {
		t.Run(scheme, func(t *testing.T) {
			calls := 0
			body := &plaintextEndpointBody{Reader: strings.NewReader("synthetic-body")}
			client := openai.NewClient(option.WithBaseURL(scheme+"://remote.invalid/v1"), option.WithMaxRetries(0),
				option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
					req.Header["authorization"] = []string{"", "Bearer synthetic-middleware"}
					return next(req)
				}),
				option.WithHTTPClient(plaintextEndpointDoer(func(req *http.Request) (*http.Response, error) {
					calls++
					if req.Body != nil {
						_ = req.Body.Close()
					}
					return &http.Response{StatusCode: http.StatusOK, Header: http.Header{}, Body: http.NoBody, Request: req}, nil
				})))
			err := client.Post(t.Context(), "models", body, nil)
			if scheme == "http" {
				if err == nil || !strings.Contains(err.Error(), "HTTPS") || calls != 0 {
					t.Fatalf("error=%v dispatches=%d; want HTTPS rejection before dispatch", err, calls)
				}
			} else if err != nil || calls != 1 {
				t.Fatalf("HTTPS error=%v dispatches=%d", err, calls)
			}
			if body.closes != 1 {
				t.Fatalf("body closed %d times", body.closes)
			}
		})
	}
}

func TestDeclaredAuthorizationTrailerRequiresOptIn(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("OPENAI_ADMIN_KEY", "")
	client := openai.NewClient(option.WithBaseURL("http://127.0.0.1"), option.WithMaxRetries(0),
		option.WithHTTPClient(plaintextEndpointDoer(func(*http.Request) (*http.Response, error) {
			t.Fatal("request with declared Authorization trailer reached the HTTP doer")
			return nil, io.EOF
		})),
		option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
			req.Trailer = http.Header{"Authorization": nil}
			return next(req)
		}))
	if err := client.Post(t.Context(), "models", strings.NewReader("fixture"), nil); err == nil || !strings.Contains(err.Error(), "HTTPS") {
		t.Fatalf("expected the HTTPS policy error, got %v", err)
	}
}

func TestWebSocketMiddlewareCredentialsTransport(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("OPENAI_ADMIN_KEY", "")
	var received atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received.Add(1)
		if r.Header.Get("Authorization") != "Bearer synthetic-middleware" {
			t.Error("missing synthetic credential")
		}
		socket, err := wire.Accept(w, r, nil)
		if err != nil {
			t.Error(err)
			return
		}
		defer func() { _ = socket.CloseNow() }()
		_, _, _ = socket.Read(r.Context())
	}))
	defer server.Close()
	for _, endpoint := range []string{"http://remote.invalid/v1", server.URL + "/v1"} {
		for _, optIn := range []bool{false, true} {
			t.Run(endpoint+"/"+fmt.Sprint(optIn), func(t *testing.T) {
				calls := 0
				opts := []option.RequestOption{option.WithBaseURL(endpoint), option.WithMaxRetries(0),
					option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
						req.Header.Set("Authorization", "Bearer synthetic-middleware")
						return next(req)
					}), option.WithHTTPClient(plaintextEndpointDoer(func(*http.Request) (*http.Response, error) { calls++; return nil, io.EOF }))}
				if optIn {
					opts = append(opts, option.WithUnsafeAllowHTTP())
				}
				before := received.Load()
				client := openai.NewClient(opts...)
				socket, err := client.Responses.Connect(t.Context(), responses.ResponseConnectionOptions{})
				if socket != nil {
					socket.Abort()
				}
				if endpoint == server.URL+"/v1" && optIn {
					if err != nil || received.Load() != before+1 {
						t.Fatalf("local upgrade error=%v requests=%d", err, received.Load()-before)
					}
				} else if err == nil || !strings.Contains(err.Error(), "HTTPS") || received.Load() != before {
					t.Fatalf("error=%v requests=%d; want HTTPS rejection", err, received.Load()-before)
				}
				if calls != 0 {
					t.Fatalf("custom doer received %d requests", calls)
				}
			})
		}
	}
}

func TestPlaintextCredentialRejectionClosesMiddlewareBody(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("OPENAI_ADMIN_KEY", "")
	for _, clone := range []string{"in-place", "clone-before-replacement", "clone-after-replacement"} {
		t.Run(clone, func(t *testing.T) {
			original := &plaintextEndpointBody{Reader: strings.NewReader("original")}
			replacement := &plaintextEndpointBody{Reader: strings.NewReader("replacement")}
			middlewareErr := errors.New("upload middleware")
			client := openai.NewClient(option.WithBaseURL("http://remote.invalid/v1"), option.WithMaxRetries(0),
				option.WithHTTPClient(plaintextEndpointDoer(func(*http.Request) (*http.Response, error) {
					t.Fatal("rejected request reached transport")
					return nil, io.EOF
				})),
				option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
					if clone == "clone-before-replacement" {
						req = req.Clone(req.Context())
					}
					// A body implemented by a value may contain a slice.
					req.Body = struct {
						io.ReadCloser
						_ []byte
					}{ReadCloser: replacement}
					if clone == "clone-after-replacement" {
						req = req.Clone(t.Context())
					}
					req.Header.Set("Authorization", "Bearer synthetic-middleware")
					res, err := next(req)
					if err != nil {
						if clone == "clone-after-replacement" {
							return res, middlewareErr
						}
						return res, fmt.Errorf("%w: %w", middlewareErr, err)
					}
					return res, nil
				}))
			err := client.Post(t.Context(), "models", original, nil)
			if !errors.Is(err, middlewareErr) || original.closes != 1 || replacement.closes != 1 {
				t.Fatalf("error=%v original closes=%d replacement closes=%d", err, original.closes, replacement.closes)
			}
		})
	}
}

func TestComposedWorkloadMiddlewareHTTPPermission(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("OPENAI_ADMIN_KEY", "")
	for _, websocket := range []bool{false, true} {
		for _, local := range []bool{false, true} {
			for _, optIn := range []bool{false, true} {
				t.Run(fmt.Sprintf("websocket=%t/local=%t/optIn=%t", websocket, local, optIn), func(t *testing.T) {
					var received atomic.Int32
					server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						received.Add(1)
						if r.Header.Get("Authorization") != "Bearer synthetic-workload" {
							t.Error("missing synthetic workload credential")
						}
						if websocket {
							socket, err := wire.Accept(w, r, nil)
							if err != nil {
								t.Error(err)
								return
							}
							defer func() { _ = socket.CloseNow() }()
							_, _, _ = socket.Read(r.Context())
							return
						}
						w.Header().Set("Content-Type", "application/json")
						_, _ = io.WriteString(w, `{"data":[]}`)
					}))
					defer server.Close()
					endpoint := "http://remote.invalid/v1"
					if local {
						endpoint = server.URL
					}
					provider := &mockSubjectTokenProvider{token: "synthetic-subject", tokenType: auth.SubjectTokenTypeJWT}
					identity, err := auth.NewWorkloadIdentityAuth(testWorkloadIdentity(provider))
					if err != nil {
						t.Fatal(err)
					}
					exchanges, dispatches := 0, 0
					issuer := plaintextEndpointDoer(func(req *http.Request) (*http.Response, error) {
						exchanges++
						return &http.Response{StatusCode: http.StatusOK, Header: http.Header{"Content-Type": {"application/json"}}, Body: io.NopCloser(strings.NewReader(`{"access_token":"synthetic-workload","expires_in":3600}`)), Request: req}, nil
					})
					opts := []option.RequestOption{option.WithBaseURL(endpoint), option.WithMaxRetries(0),
						option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
							return auth.WorkloadIdentityMiddleware(identity, issuer, req, next)
						}),
						option.WithHTTPClient(plaintextEndpointDoer(func(*http.Request) (*http.Response, error) { dispatches++; return nil, io.EOF }))}
					if optIn {
						opts = append(opts, option.WithUnsafeAllowHTTP())
					}
					client := openai.NewClient(opts...)
					if websocket {
						socket, connectErr := client.Responses.Connect(t.Context(), responses.ResponseConnectionOptions{})
						err = connectErr
						if socket != nil {
							socket.Abort()
						}
					} else {
						_, err = client.Models.List(t.Context())
					}
					if local && optIn {
						if err != nil || received.Load() != 1 || exchanges != 1 || provider.GetCallCount() != 1 {
							t.Fatalf("opted-in request: error=%v received=%d exchanges=%d provider=%d", err, received.Load(), exchanges, provider.GetCallCount())
						}
					} else if err == nil || !strings.Contains(err.Error(), "HTTPS") || received.Load() != 0 {
						t.Fatalf("blocked request: error=%v received=%d exchanges=%d provider=%d", err, received.Load(), exchanges, provider.GetCallCount())
					}
					// The direct helper authenticates before yielding to the SDK transport.
					if exchanges != 1 || provider.GetCallCount() != 1 {
						t.Fatalf("direct helper exchanges=%d provider=%d", exchanges, provider.GetCallCount())
					}
					if dispatches != 0 {
						t.Fatalf("custom client received %d requests", dispatches)
					}
				})
			}
		}
	}
}

func TestWorkloadCannotInheritAnotherClientsHTTPPermission(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("OPENAI_ADMIN_KEY", "")
	for _, websocket := range []bool{false, true} {
		t.Run(fmt.Sprintf("websocket=%t", websocket), func(t *testing.T) {
			provider := &mockSubjectTokenProvider{token: "synthetic-subject", tokenType: auth.SubjectTokenTypeJWT}
			dispatches := 0
			inner := openai.NewClient(option.WithBaseURL("http://127.0.0.1:9999/v1"), option.WithMaxRetries(0),
				option.WithWorkloadIdentity(testWorkloadIdentity(provider)),
				option.WithHTTPClient(plaintextEndpointDoer(func(*http.Request) (*http.Response, error) {
					dispatches++
					return nil, io.EOF
				})))
			// Share the request context, as instrumentation and nested clients
			// commonly do. Only the outer client has opted in to local HTTP.
			outer := openai.NewClient(option.WithBaseURL("http://127.0.0.1:9998/v1"), option.WithUnsafeAllowHTTP(),
				option.WithMiddleware(func(req *http.Request, _ option.MiddlewareNext) (*http.Response, error) {
					var innerErr error
					if websocket {
						socket, connectErr := inner.Responses.Connect(req.Context(), responses.ResponseConnectionOptions{})
						innerErr = connectErr
						if socket != nil {
							socket.Abort()
						}
					} else {
						_, innerErr = inner.Models.List(req.Context())
					}
					if innerErr == nil || !strings.Contains(innerErr.Error(), "HTTPS") {
						t.Errorf("unopted inner request should be rejected before authentication: %v", innerErr)
					}
					return &http.Response{StatusCode: http.StatusOK, Header: http.Header{}, Body: http.NoBody, Request: req}, nil
				}))
			if err := outer.Get(t.Context(), "models", nil, nil); err != nil {
				t.Fatal(err)
			}
			if provider.GetCallCount() != 0 || dispatches != 0 {
				t.Fatalf("unopted inner request made provider=%d HTTP=%d calls", provider.GetCallCount(), dispatches)
			}
		})
	}
}

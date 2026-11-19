package openai_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"

	wire "github.com/coder/websocket"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

// The first request is intentionally unauthenticated. Redirect handling runs
// inside http.Client.Do, after the SDK middleware has finished.
func TestRedirectAddedCredentialsRequireSecureTransport(t *testing.T) {
	for _, credential := range []string{"check-redirect", "url-user"} {
		for _, destination := range []string{"remote-http", "loopback-http", "https"} {
			t.Run(credential+"/"+destination, func(t *testing.T) {
				t.Setenv("OPENAI_API_KEY", "")
				t.Setenv("OPENAI_ADMIN_KEY", "")

				var initial, configuredRedirect, directRedirect atomic.Int32
				base := "http://redirect.invalid/v1/"
				if destination == "https" {
					base = "https://redirect.invalid/v1/"
				}
				if destination == "loopback-http" {
					server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
						directRedirect.Add(1)
						if req.URL.Path != "/v1/after" || req.Header.Get("Authorization") == "" {
							t.Error("direct loopback request did not contain redirected credentials")
						}
						w.Header().Set("Content-Type", "application/json")
						_, _ = io.WriteString(w, `{"data":[],"object":"list"}`)
					}))
					t.Cleanup(server.Close)
					base = server.URL + "/v1/"
				}

				transport := originTestRoundTripper(func(req *http.Request) (*http.Response, error) {
					if req.URL.Path == "/v1/before" {
						initial.Add(1)
						if req.Header.Get("Authorization") != "" {
							t.Error("initial request unexpectedly authenticated")
						}
						location, err := url.Parse(base + "after?fixture=query-marker")
						if err != nil {
							t.Fatal("could not build synthetic redirect")
						}
						if credential == "url-user" {
							location.User = url.UserPassword("synthetic-user", "synthetic-password")
						}
						return &http.Response{
							StatusCode: http.StatusFound,
							Header:     http.Header{"Location": {location.String()}},
							Body:       io.NopCloser(strings.NewReader("")),
							Request:    req,
						}, nil
					}
					configuredRedirect.Add(1)
					if req.URL.Path != "/v1/after" || req.Header.Get("Authorization") == "" {
						t.Error("configured transport request did not contain redirected credentials")
					}
					return &http.Response{
						StatusCode: http.StatusOK,
						Header:     http.Header{"Content-Type": {"application/json"}},
						Body:       io.NopCloser(strings.NewReader(`{"data":[],"object":"list"}`)),
						Request:    req,
					}, nil
				})
				native := &http.Client{Transport: transport}
				if credential == "check-redirect" {
					native.CheckRedirect = func(req *http.Request, _ []*http.Request) error {
						req.Header.Set("Authorization", "Bearer synthetic-redirect")
						return nil
					}
				}
				opts := []option.RequestOption{
					option.WithBaseURL(base),
					option.WithHTTPClient(native),
					option.WithMaxRetries(0),
				}
				if destination == "loopback-http" {
					opts = append(opts, option.WithUnsafeAllowHTTP())
				}
				client := openai.NewClient(opts...)
				var response map[string]any
				err := client.Get(t.Context(), "before", nil, &response)
				if initial.Load() != 1 {
					t.Fatalf("initial configured calls = %d; want 1", initial.Load())
				}
				switch destination {
				case "remote-http":
					if err == nil || !strings.Contains(err.Error(), "HTTPS") || configuredRedirect.Load() != 0 {
						t.Fatalf("remote HTTP: rejected=%t configured redirect calls=%d; want HTTPS rejection before dispatch",
							err != nil, configuredRedirect.Load())
					}
					if strings.Contains(err.Error(), "query-marker") || strings.Contains(err.Error(), "synthetic-user") {
						t.Error("policy error included the redirect URL")
					}
				case "loopback-http":
					if err != nil || configuredRedirect.Load() != 0 || directRedirect.Load() != 1 {
						t.Fatalf("loopback HTTP: rejected=%t configured redirect calls=%d direct calls=%d; want only direct dispatch",
							err != nil, configuredRedirect.Load(), directRedirect.Load())
					}
				case "https":
					if err != nil || configuredRedirect.Load() != 1 || directRedirect.Load() != 0 {
						t.Fatalf("HTTPS: rejected=%t configured redirect calls=%d direct calls=%d; want configured transport",
							err != nil, configuredRedirect.Load(), directRedirect.Load())
					}
				}
			})
		}
	}
}

func TestDirectLoopbackExchangeKeepsTransportAcrossRedirect(t *testing.T) {
	for _, protocol := range []string{"HTTP", "WebSocket"} {
		t.Run(protocol, func(t *testing.T) {
			var first, last atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				if req.URL.Path != "/finish" {
					first.Add(1)
					if req.Header.Get("Authorization") != "Bearer fixture" {
						t.Error("initial request missing fixture key")
					}
					http.Redirect(w, req, "/finish", http.StatusFound)
					return
				}
				last.Add(1)
				if req.Header.Get("Authorization") != "" {
					t.Error("redirect retained a removed credential")
				}
				if protocol == "WebSocket" {
					socket, err := wire.Accept(w, req, nil)
					if err != nil {
						t.Error(err)
						return
					}
					defer func() { _ = socket.CloseNow() }()
					_, _, _ = socket.Read(req.Context())
				}
			}))
			defer server.Close()
			redirects := 0
			native := &http.Client{
				CheckRedirect: func(req *http.Request, _ []*http.Request) error {
					redirects++
					req.Header.Del("Authorization")
					return nil
				},
				Transport: originTestRoundTripper(func(*http.Request) (*http.Response, error) {
					return nil, errors.New("direct exchange reached the configured transport")
				}),
			}
			client := openai.NewClient(option.WithBaseURL(server.URL), option.WithUnsafeAllowHTTP(), option.WithAPIKey("fixture"),
				option.WithMaxRetries(0), option.WithHTTPClient(native), option.WithHTTPClient(plaintextEndpointDoer(func(*http.Request) (*http.Response, error) {
					return nil, errors.New("direct exchange reached the custom doer")
				})))
			if protocol == "HTTP" {
				if err := client.Get(t.Context(), "start", nil, nil); err != nil {
					t.Fatal(err)
				}
			} else {
				socket, err := client.Responses.Connect(t.Context(), responses.ResponseConnectionOptions{})
				if err != nil {
					t.Fatal(err)
				}
				socket.Abort()
			}
			if first.Load() != 1 || redirects != 1 || last.Load() != 1 {
				t.Fatalf("initial=%d redirects=%d final=%d, want 1 each", first.Load(), redirects, last.Load())
			}
		})
	}
}

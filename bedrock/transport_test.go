package bedrock

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/openai/openai-go/v3/option"
)

func TestAuthenticatedEndpointRequiresHTTPS(t *testing.T) {
	clearAWSEnvironment(t)
	for _, source := range []string{"config", "environment"} {
		for _, auth := range []string{"bearer", "sigv4"} {
			for _, endpoint := range []Endpoint{EndpointMantle, EndpointRuntime} {
				for _, scheme := range []string{"http", "https"} {
					t.Run(source+"/"+auth+"/"+string(endpoint)+"/"+scheme, func(t *testing.T) {
						providerCalls, transportCalls := 0, 0
						cfg := Config{AWSRegion: "us-east-1", Endpoint: endpoint}
						if auth == "bearer" {
							cfg.BedrockTokenProvider = func(context.Context) (string, error) {
								providerCalls++
								return "test-bearer", nil
							}
						} else {
							cfg.AWSCredentialsProvider = aws.CredentialsProviderFunc(func(context.Context) (aws.Credentials, error) {
								providerCalls++
								return aws.Credentials{AccessKeyID: "test-access", SecretAccessKey: "test-secret", SessionToken: "test-session"}, nil
							})
						}
						baseURL := scheme + "://gateway.example/openai/v1"
						if source == "config" {
							cfg.BaseURL = baseURL
						} else {
							t.Setenv("AWS_BEDROCK_BASE_URL", baseURL)
						}
						client, err := NewClient(t.Context(), cfg, option.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
							transportCalls++
							if req.URL.Scheme != scheme || req.URL.Host != "gateway.example" || req.Header.Get("Authorization") == "" {
								t.Error("expected authenticated request to the configured gateway")
							}
							if auth == "sigv4" && req.Header.Get("X-Amz-Security-Token") != "test-session" {
								t.Error("missing temporary session token")
							}
							return successfulResponse(req), nil
						})}))
						if scheme == "http" {
							if err == nil {
								var response map[string]any
								if callErr := client.Post(t.Context(), "/responses", map[string]string{"input": "synthetic prompt"}, &response); callErr != nil {
									t.Fatalf("unexpected request error: %v", callErr)
								}
							}
							if err == nil || !strings.Contains(err.Error(), "require HTTPS") || providerCalls != 0 || transportCalls != 0 {
								t.Fatalf("HTTP must fail before credentials or transport: error=%v provider calls=%d transport calls=%d", err, providerCalls, transportCalls)
							}
							return
						}
						if err != nil {
							t.Fatal(err)
						}
						var response map[string]any
						if err := client.Get(t.Context(), "/models", nil, &response); err != nil {
							t.Fatal(err)
						}
						if providerCalls == 0 || transportCalls != 1 {
							t.Fatalf("HTTPS request not authenticated and sent: provider calls=%d transport calls=%d", providerCalls, transportCalls)
						}
					})
				}
			}
		}
	}
}

func TestUnsafeHTTPRequiresLoopback(t *testing.T) {
	clearAWSEnvironment(t)
	for _, host := range []string{"gateway.example", "192.0.2.1", "10.0.0.1", "localhost.example", "localhost.", "127.0.0.1.example", "[::]", "127.1"} {
		t.Run(host, func(t *testing.T) {
			_, err := NewClient(t.Context(), Config{APIKey: "test-bearer", BaseURL: "http://" + host, UnsafeAllowHTTP: true})
			if err == nil || !strings.Contains(err.Error(), "require HTTPS") {
				t.Fatalf("remote or ambiguous HTTP endpoint accepted: %v", err)
			}
		})
	}
	for _, host := range []string{"localhost", "127.0.0.1", "[::1]"} {
		t.Run(host+"/default", func(t *testing.T) {
			_, err := NewClient(t.Context(), Config{APIKey: "test-bearer", BaseURL: "http://" + host})
			if err == nil || !strings.Contains(err.Error(), "require HTTPS") {
				t.Fatalf("loopback HTTP accepted without opt-in: %v", err)
			}
		})
	}
}

func TestUnsafeHTTPConnectsDirectlyToLoopback(t *testing.T) {
	clearAWSEnvironment(t)
	for _, host := range []string{"127.0.0.1", "localhost", "::1"} {
		for _, auth := range []string{"bearer", "sigv4", "gateway"} {
			t.Run(host+"/"+auth, func(t *testing.T) {
				var calls, bypassCalls atomic.Int32
				server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					calls.Add(1)
					if req.Header.Get("Authorization") == "" {
						t.Error("missing authentication")
					}
					if auth == "sigv4" && req.Header.Get("X-Amz-Security-Token") != "test-session" {
						t.Error("missing session token")
					}
					w.Header().Set("Content-Type", "application/json")
					_, _ = w.Write([]byte(`{}`))
				}))
				if host == "::1" {
					if err := server.Listener.Close(); err != nil {
						t.Fatal(err)
					}
					listener, err := net.Listen("tcp6", "[::1]:0")
					if err != nil {
						t.Skipf("IPv6 loopback unavailable: %v", err)
					}
					server.Listener = listener
				}
				server.Start()
				defer server.Close()
				baseURL := server.URL
				if host == "localhost" {
					baseURL = strings.Replace(baseURL, "127.0.0.1", "localhost", 1)
				}
				cfg := Config{BaseURL: baseURL, UnsafeAllowHTTP: true}
				opts := []option.RequestOption{}
				switch auth {
				case "bearer":
					cfg.APIKey = "test-bearer"
				case "sigv4":
					cfg.AWSRegion, cfg.AWSAccessKeyID, cfg.AWSSecretAccessKey, cfg.AWSSessionToken = "us-east-1", "test-access", "test-secret", "test-session"
				case "gateway":
					cfg.SkipAuth = true
					opts = append(opts, option.WithAPIKey("test-gateway"))
				}
				transport := &http.Transport{
					Proxy: func(*http.Request) (*url.URL, error) {
						bypassCalls.Add(1)
						return nil, errors.New("proxy must not be used")
					},
					DialContext: func(context.Context, string, string) (net.Conn, error) {
						bypassCalls.Add(1)
						return nil, errors.New("custom dial must not be used")
					},
				}
				transport.RegisterProtocol("http", roundTripFunc(func(req *http.Request) (*http.Response, error) {
					bypassCalls.Add(1)
					return successfulResponse(req), nil
				}))
				transport.Protocols = &http.Protocols{}
				transport.Protocols.SetUnencryptedHTTP2(true)
				protocolHandler := func(_ string, conn *tls.Conn) http.RoundTripper {
					bypassCalls.Add(1)
					if raw, ok := conn.NetConn().(interface{ UnencryptedNetConn() net.Conn }); ok {
						_ = raw.UnencryptedNetConn().Close()
					}
					return roundTripFunc(func(req *http.Request) (*http.Response, error) {
						if req.Header.Get("Authorization") != "" {
							bypassCalls.Add(1)
						}
						return successfulResponse(req), nil
					})
				}
				transport.TLSNextProto = map[string]func(string, *tls.Conn) http.RoundTripper{
					"h2":                protocolHandler,
					"unencrypted_http2": protocolHandler,
				}
				opts = append(opts, option.WithHTTPClient(&http.Client{Transport: transport}))
				client, err := NewClient(t.Context(), cfg, opts...)
				if err != nil {
					t.Fatal(err)
				}
				for range 2 {
					var response map[string]any
					if err := client.Get(t.Context(), "/models", nil, &response); err != nil {
						t.Fatal(err)
					}
				}
				if calls.Load() != 2 || bypassCalls.Load() != 0 {
					t.Fatalf("direct calls=%d bypass calls=%d", calls.Load(), bypassCalls.Load())
				}
			})
		}
	}
}

func TestUnsafeHTTPRejectsOpaqueTransport(t *testing.T) {
	transportCalls := 0
	client, err := NewClient(t.Context(), Config{APIKey: "test-bearer", BaseURL: "http://127.0.0.1", UnsafeAllowHTTP: true},
		option.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			transportCalls++
			return successfulResponse(req), nil
		})}))
	if err != nil {
		t.Fatal(err)
	}
	var response map[string]any
	err = client.Get(t.Context(), "/models", nil, &response)
	if err == nil || !strings.Contains(err.Error(), "Transport") || transportCalls != 0 {
		t.Fatalf("opaque transport must be rejected: error=%v calls=%d", err, transportCalls)
	}
}

func TestSkipAuthHTTPGatewayCredentials(t *testing.T) {
	clearAWSEnvironment(t)
	for _, key := range []string{"none", "api", "admin"} {
		for _, requestLevel := range []bool{false, true} {
			t.Run(key+"/"+map[bool]string{false: "client", true: "request"}[requestLevel], func(t *testing.T) {
				transportCalls := 0
				var authOpts []option.RequestOption
				switch key {
				case "api":
					authOpts = append(authOpts, option.WithAPIKey("test-gateway"))
				case "admin":
					authOpts = append(authOpts, option.WithAdminAPIKey("test-gateway-admin"))
				}
				opts := []option.RequestOption{option.WithHTTPClient(httpDoerFunc(func(req *http.Request) (*http.Response, error) {
					transportCalls++
					return successfulResponse(req), nil
				}))}
				if !requestLevel {
					opts = append(opts, authOpts...)
					authOpts = nil
				}
				client, err := NewClient(t.Context(), Config{SkipAuth: true, BaseURL: "http://gateway.example"}, opts...)
				if err != nil {
					t.Fatal(err)
				}
				var response map[string]any
				err = client.Get(t.Context(), "/models", nil, &response, authOpts...)
				if key == "none" {
					if err != nil || transportCalls != 1 {
						t.Fatalf("unauthenticated HTTP failed: error=%v calls=%d", err, transportCalls)
					}
				} else if err == nil || !strings.Contains(err.Error(), "require HTTPS") || transportCalls != 0 {
					t.Fatalf("HTTP gateway credentials must be rejected: error=%v calls=%d", err, transportCalls)
				}
			})
		}
	}
}

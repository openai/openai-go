package bedrock

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type finalizerTestBody struct {
	io.Reader
	closes   int
	closeErr error
}

func (b *finalizerTestBody) Close() error {
	b.closes++
	return b.closeErr
}

func TestTransportFinalizerErrorClosesBody(t *testing.T) {
	clearAWSEnvironment(t)
	for _, key := range []string{"api", "admin"} {
		credential := option.WithAPIKey("synthetic-gateway")
		if key == "admin" {
			credential = option.WithAdminAPIKey("synthetic-gateway-admin")
		}
		for _, level := range []string{"client", "request"} {
			for _, failure := range []string{"remote-http", "loopback-custom-doer", "loopback-opaque-transport"} {
				for _, closeFails := range []bool{false, true} {
					name := key + "/" + level + "/" + failure
					if closeFails {
						name += "/close-error"
					}
					t.Run(name, func(t *testing.T) {
						body := &finalizerTestBody{Reader: strings.NewReader("synthetic request")}
						if closeFails {
							body.closeErr = errors.New("synthetic close error")
						}
						calls := 0
						dispatch := func(req *http.Request) (*http.Response, error) {
							calls++
							return successfulResponse(req), nil
						}
						cfg := Config{SkipAuth: true, BaseURL: "http://gateway.example"}
						clientOption := option.WithHTTPClient(httpDoerFunc(dispatch))
						wantError := "require HTTPS"
						if failure != "remote-http" {
							cfg.BaseURL, cfg.UnsafeAllowHTTP = "http://127.0.0.1", true
							wantError = "requires an *http.Client"
						}
						if failure == "loopback-opaque-transport" {
							clientOption = option.WithHTTPClient(&http.Client{Transport: roundTripFunc(dispatch)})
							wantError = "requires a nil Transport or an *http.Transport"
						}
						opts := []option.RequestOption{clientOption}
						requestOpts := []option.RequestOption{credential}
						if level == "client" {
							opts = append(opts, credential)
							requestOpts = nil
						}
						client, err := NewClient(t.Context(), cfg, opts...)
						if err != nil {
							t.Fatal(err)
						}
						err = client.Post(t.Context(), "/responses", body, nil, requestOpts...)
						if err == nil || !strings.Contains(err.Error(), wantError) || calls != 0 || body.closes != 1 {
							t.Fatalf("error=%v dispatches=%d closes=%d", err, calls, body.closes)
						}
					})
				}
			}
		}
	}
}

func TestSkipAuthHTTPGatewayAuthorizationOverrides(t *testing.T) {
	clearAWSEnvironment(t)
	for _, key := range []string{"api", "admin"} {
		credential := option.WithAPIKey("synthetic-gateway")
		if key == "admin" {
			credential = option.WithAdminAPIKey("synthetic-gateway-admin")
		}
		for _, overrideName := range []string{"delete", "empty", "caller-managed"} {
			override := option.WithHeaderDel("Authorization")
			wantAuthorization := ""
			if overrideName == "empty" {
				override = option.WithHeader("Authorization", "")
			} else if overrideName == "caller-managed" {
				wantAuthorization = "Basic synthetic-caller-managed"
				override = option.WithHeader("Authorization", wantAuthorization)
			}
			for _, test := range []struct {
				name           string
				clientOptions  []option.RequestOption
				requestOptions []option.RequestOption
				authenticated  bool
			}{
				{"client-override-only", []option.RequestOption{override}, nil, false},
				{"request-override-only", nil, []option.RequestOption{override}, false},
				{"client-removes", []option.RequestOption{credential, override}, nil, false},
				{"request-removes", []option.RequestOption{credential}, []option.RequestOption{override}, false},
				{"request-key-then-removes", nil, []option.RequestOption{credential, override}, false},
				{"client-restores", []option.RequestOption{override, credential}, nil, true},
				{"request-restores", []option.RequestOption{credential, override}, []option.RequestOption{credential}, true},
				{"request-removes-then-restores", []option.RequestOption{credential}, []option.RequestOption{override, credential}, true},
			} {
				for _, endpoint := range []string{"http://gateway.example", "http://127.0.0.1"} {
					t.Run(key+"/"+overrideName+"/"+test.name+"/"+endpoint, func(t *testing.T) {
						calls := 0
						opts := append([]option.RequestOption{option.WithHTTPClient(httpDoerFunc(func(req *http.Request) (*http.Response, error) {
							calls++
							if got := req.Header.Get("Authorization"); got != wantAuthorization {
								t.Errorf("Authorization = %q, want %q", got, wantAuthorization)
							}
							return successfulResponse(req), nil
						}))}, test.clientOptions...)
						client, err := NewClient(t.Context(), Config{SkipAuth: true, BaseURL: endpoint, UnsafeAllowHTTP: true}, opts...)
						if err != nil {
							t.Fatal(err)
						}
						err = client.Get(t.Context(), "/models", nil, nil, test.requestOptions...)
						if test.authenticated {
							wantError := "require HTTPS"
							if endpoint == "http://127.0.0.1" {
								wantError = "requires an *http.Client"
							}
							if err == nil || !strings.Contains(err.Error(), wantError) || calls != 0 {
								t.Fatalf("authenticated error=%v dispatches=%d", err, calls)
							}
						} else if err != nil || calls != 1 {
							t.Fatalf("unauthenticated error=%v dispatches=%d", err, calls)
						}
					})
				}
			}
		}
	}
}

func TestSkipAuthHTTPGatewayEndpointSecurity(t *testing.T) {
	clearAWSEnvironment(t)
	for _, key := range []string{"api", "admin"} {
		for _, endpoint := range []string{"model", "certificate"} {
			t.Run(key+"/"+endpoint, func(t *testing.T) {
				calls := 0
				credential := option.WithAPIKey("synthetic-gateway")
				if key == "admin" {
					credential = option.WithAdminAPIKey("synthetic-gateway-admin")
				}
				client, err := NewClient(t.Context(), Config{SkipAuth: true, BaseURL: "http://gateway.example"}, credential,
					option.WithHTTPClient(httpDoerFunc(func(req *http.Request) (*http.Response, error) {
						calls++
						if req.Header.Get("Authorization") != "" {
							t.Error("endpoint emitted an incompatible credential")
						}
						return successfulResponse(req), nil
					})))
				if err != nil {
					t.Fatal(err)
				}
				if endpoint == "model" {
					_, err = client.Models.Get(t.Context(), "synthetic-model")
				} else {
					_, err = client.Admin.Organization.Certificates.Get(t.Context(), "synthetic-certificate", openai.AdminOrganizationCertificateGetParams{})
				}
				authenticated := key == "api" && endpoint == "model" || key == "admin" && endpoint == "certificate"
				if authenticated {
					if err == nil || !strings.Contains(err.Error(), "require HTTPS") || calls != 0 {
						t.Fatalf("authenticated error=%v dispatches=%d", err, calls)
					}
				} else if err != nil || calls != 1 {
					t.Fatalf("unauthenticated error=%v dispatches=%d", err, calls)
				}
			})
		}
	}
}

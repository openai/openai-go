package azure

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func TestAzureRejectsOriginChangesBeforeTokenAcquisition(t *testing.T) {
	for _, tc := range []struct {
		name     string
		endpoint string
		mutate   func(*http.Request)
	}{
		{name: "hostname", mutate: func(req *http.Request) { req.URL.Host = "other.example" }},
		{name: "port", mutate: func(req *http.Request) { req.URL.Host = "azure.example:444" }},
		{name: "scheme", mutate: func(req *http.Request) { req.URL.Scheme = "http" }},
		{name: "authority", mutate: func(req *http.Request) { req.Host = "other.example" }},
		{name: "opaque target", mutate: func(req *http.Request) { req.URL.Opaque = "//other.example/models" }},
		{name: "request URI", mutate: func(req *http.Request) { req.RequestURI = "/models" }},
		{name: "loopback port", endpoint: "http://127.0.0.1:8080", mutate: func(req *http.Request) { req.URL.Host = "127.0.0.1:8081" }},
		{name: "loopback scheme", endpoint: "http://127.0.0.1:8080", mutate: func(req *http.Request) { req.URL.Scheme = "https" }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			credential := &countingTokenCredential{}
			transportCalls, middlewareCalls := 0, 0
			endpoint := tc.endpoint
			if endpoint == "" {
				endpoint = "https://azure.example"
			}
			client := openai.NewClient(
				WithEndpoint(endpoint, "2024-10-21"),
				WithTokenCredential(credential),
				WithUnsafeAllowHTTP(),
				option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
					middlewareCalls++
					tc.mutate(req)
					return next(req)
				}),
				option.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
					transportCalls++
					return successfulAzureResponse(req), nil
				})}),
			)
			methods := []struct {
				name string
				call func(context.Context, string, any, any, ...option.RequestOption) error
			}{
				{"Get", client.Get}, {"Post", client.Post}, {"Put", client.Put},
				{"Patch", client.Patch}, {"Delete", client.Delete},
				{"Execute", func(ctx context.Context, path string, params, res any, opts ...option.RequestOption) error {
					return client.Execute(ctx, http.MethodHead, path, params, res, opts...)
				}},
			}
			for _, method := range methods {
				if err := method.call(context.Background(), "models", nil, nil); err == nil || !strings.Contains(err.Error(), "request URL origin must match the configured base URL") {
					t.Errorf("%s: expected origin rejection, got %v", method.name, err)
				}
			}
			if calls := credential.calls.Load(); calls != 0 {
				t.Errorf("token acquisition calls = %d, want 0", calls)
			}
			if transportCalls != 0 || middlewareCalls != len(methods) {
				t.Errorf("transport calls = %d, middleware calls = %d; want 0 and %d", transportCalls, middlewareCalls, len(methods))
			}
		})
	}
}

func TestAzureTrustedEndpointOverrides(t *testing.T) {
	for _, mode := range []string{"API key", "token"} {
		t.Run(mode, func(t *testing.T) {
			credential := &countingTokenCredential{}
			auth := WithAPIKey("synthetic-azure-key")
			header, value := "Api-Key", "synthetic-azure-key"
			if mode == "token" {
				auth = WithTokenCredential(credential)
				header, value = "Authorization", "Bearer counting_token"
			}
			baseOptions := []option.RequestOption{
				WithEndpoint("https://azure.example/original", "2024-10-21"),
				auth,
			}
			for _, layer := range []string{"client", "service", "request"} {
				t.Run(layer, func(t *testing.T) {
					t.Parallel()
					transportCalls := 0
					endpoint := "https://" + layer + ".example/gateway"
					opts := append([]option.RequestOption(nil), baseOptions...)
					opts = append(opts, option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
						req.URL.Host = strings.ToUpper(req.URL.Hostname()) + ":443"
						return next(req)
					}))
					opts = append(opts, option.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
						transportCalls++
						if !strings.EqualFold(req.URL.Hostname(), layer+".example") || req.URL.Path != "/gateway/openai/models/test" {
							t.Errorf("unexpected destination: %s", req.URL.Redacted())
						}
						if req.Header.Get(header) != value {
							t.Error("synthetic Azure credential missing")
						}
						return successfulAzureResponse(req), nil
					})}))
					var requestOptions []option.RequestOption
					if layer == "client" {
						opts = append(opts, option.WithBaseURL(endpoint))
					} else if layer == "request" {
						requestOptions = append(requestOptions, option.WithBaseURL(endpoint))
					}
					client := openai.NewClient(opts...)
					service := client.Models
					if layer == "service" {
						service = openai.NewModelService(append(client.Options, option.WithBaseURL(endpoint))...)
					}
					if _, err := service.Get(context.Background(), "test", requestOptions...); err != nil {
						t.Fatal(err)
					}
					if transportCalls != 1 {
						t.Errorf("transport calls = %d, want 1", transportCalls)
					}
				})
			}
		})
	}
}

func TestAzureOriginRejectionClosesReplacementBody(t *testing.T) {
	for _, mode := range []string{"API key", "token"} {
		for _, clone := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/clone=%t", mode, clone), func(t *testing.T) {
				closed := make(chan struct{}, 2)
				auth := WithAPIKey("synthetic-azure-key")
				if mode == "token" {
					auth = WithTokenCredential(&countingTokenCredential{})
				}
				client := openai.NewClient(
					WithEndpoint("https://azure.example", "2024-10-21"), auth,
					option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
						if clone {
							req = req.Clone(req.Context())
						}
						req.Body = &closeTrackingBody{Reader: strings.NewReader("synthetic body"), closed: closed}
						req.URL.Host = "other.example"
						return next(req)
					}),
					option.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
						t.Error("rejected request reached transport")
						return successfulAzureResponse(req), nil
					})}),
				)
				if err := client.Post(context.Background(), "models", nil, nil); err == nil || !strings.Contains(err.Error(), "request URL origin must match the configured base URL") {
					t.Fatalf("expected origin rejection, got %v", err)
				}
				if calls := len(closed); calls != 1 {
					t.Errorf("replacement body close calls = %d, want 1", calls)
				}
			})
		}
	}
}

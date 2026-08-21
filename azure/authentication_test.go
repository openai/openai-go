package azure

import (
	"context"
	"io"
	"net/http"
	"slices"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func TestEndpointRequiresExplicitAzureAuthentication(t *testing.T) {
	tests := []struct {
		name        string
		apiKey      string
		adminAPIKey string
	}{
		{name: "OpenAI API key", apiKey: "ambient-openai-key"},
		{name: "OpenAI admin API key", adminAPIKey: "ambient-openai-admin-key"},
		{name: "both OpenAI keys", apiKey: "ambient-openai-key", adminAPIKey: "ambient-openai-admin-key"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("OPENAI_API_KEY", test.apiKey)
			t.Setenv("OPENAI_ADMIN_KEY", test.adminAPIKey)

			middlewareCalls := 0
			transportCalls := 0
			client := openai.NewClient(
				WithEndpoint("https://my-resource.openai.azure.com", "2024-10-21"),
				option.WithMaxRetries(0),
				option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
					middlewareCalls++
					return next(req)
				}),
				option.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
					transportCalls++
					return successfulAzureResponse(req), nil
				})}),
			)

			var res map[string]any
			err := client.Execute(context.Background(), http.MethodGet, "models", nil, &res)
			if err == nil || !strings.Contains(err.Error(), "authentication is required") {
				t.Fatalf("error = %v", err)
			}
			if middlewareCalls != 0 || transportCalls != 0 {
				t.Fatalf("middleware calls = %d, transport calls = %d", middlewareCalls, transportCalls)
			}
		})
	}
}

func TestAzureAuthenticationRequiresEndpoint(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("OPENAI_ADMIN_KEY", "")
	tests := []struct {
		name string
		auth option.RequestOption
	}{
		{name: "Azure API key", auth: WithAPIKey("azure-api-key")},
		{name: "Azure token credential", auth: WithTokenCredential(&fake.TokenCredential{})},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			middlewareCalls := 0
			transportCalls := 0
			client := openai.NewClient(
				test.auth,
				option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
					middlewareCalls++
					return next(req)
				}),
				option.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
					transportCalls++
					return successfulAzureResponse(req), nil
				})}),
			)

			var res map[string]any
			err := client.Execute(context.Background(), http.MethodGet, "models", nil, &res)
			if err == nil || !strings.Contains(err.Error(), "requires azure.WithEndpoint") {
				t.Fatalf("error = %v", err)
			}
			if middlewareCalls != 0 || transportCalls != 0 {
				t.Fatalf("middleware calls = %d, transport calls = %d", middlewareCalls, transportCalls)
			}
		})
	}
}

func TestAzureRejectsNilTokenCredentialBeforeMiddleware(t *testing.T) {
	var typedNil *fake.TokenCredential
	for _, test := range []struct {
		name       string
		credential azcore.TokenCredential
	}{
		{name: "nil interface"},
		{name: "typed nil", credential: typedNil},
	} {
		t.Run(test.name, func(t *testing.T) {
			middlewareCalls := 0
			transportCalls := 0
			client := openai.NewClient(
				WithEndpoint("https://my-resource.openai.azure.com", "2024-10-21"),
				WithTokenCredential(test.credential),
				option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
					middlewareCalls++
					return next(req)
				}),
				option.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
					transportCalls++
					return successfulAzureResponse(req), nil
				})}),
			)

			var res map[string]any
			err := client.Execute(context.Background(), http.MethodGet, "models", nil, &res)
			if err == nil || !strings.Contains(err.Error(), "must not be nil") {
				t.Fatalf("error = %v", err)
			}
			if middlewareCalls != 0 || transportCalls != 0 {
				t.Fatalf("middleware calls = %d, transport calls = %d", middlewareCalls, transportCalls)
			}
		})
	}
}

func TestReusedTokenCredentialOptionAuthenticatesOnce(t *testing.T) {
	credential := &countingTokenCredential{}
	auth := WithTokenCredential(credential)
	var authorization string
	client := openai.NewClient(
		WithEndpoint("https://my-resource.openai.azure.com", "2024-10-21"),
		auth,
		option.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			authorization = req.Header.Get("Authorization")
			return successfulAzureResponse(req), nil
		})}),
	)

	var res map[string]any
	if err := client.Execute(context.Background(), http.MethodGet, "models", nil, &res, auth); err != nil {
		t.Fatal(err)
	}
	if calls := credential.calls.Load(); calls != 1 {
		t.Fatalf("token acquisition calls = %d, want 1", calls)
	}
	if authorization != "Bearer counting_token" {
		t.Fatalf("Authorization = %q", authorization)
	}
}

func TestAzureAuthenticationIsolatesOpenAIEnvironment(t *testing.T) {
	environments := []struct {
		name        string
		apiKey      string
		adminAPIKey string
	}{
		{name: "OpenAI API key", apiKey: "ambient-openai-key"},
		{name: "OpenAI admin API key", adminAPIKey: "ambient-openai-admin-key"},
		{name: "both OpenAI keys", apiKey: "ambient-openai-key", adminAPIKey: "ambient-openai-admin-key"},
	}
	authModes := []struct {
		name              string
		option            func() option.RequestOption
		wantAuthorization string
		wantAPIKey        string
	}{
		{
			name:       "Azure API key",
			option:     func() option.RequestOption { return WithAPIKey("azure-api-key") },
			wantAPIKey: "azure-api-key",
		},
		{
			name:              "Azure token credential",
			option:            func() option.RequestOption { return WithTokenCredential(&fake.TokenCredential{}) },
			wantAuthorization: "Bearer fake_token",
		},
	}

	for _, environment := range environments {
		for _, authMode := range authModes {
			t.Run(environment.name+"/"+authMode.name, func(t *testing.T) {
				t.Setenv("OPENAI_API_KEY", environment.apiKey)
				t.Setenv("OPENAI_ADMIN_KEY", environment.adminAPIKey)

				var middlewareAuthorization string
				var middlewareAPIKey string
				var transportAuthorization string
				var transportAPIKey string
				client := openai.NewClient(
					WithEndpoint("https://my-resource.openai.azure.com", "2024-10-21"),
					option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
						middlewareAuthorization = req.Header.Get("Authorization")
						middlewareAPIKey = req.Header.Get("Api-Key")
						return next(req)
					}),
					authMode.option(),
					option.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
						transportAuthorization = req.Header.Get("Authorization")
						transportAPIKey = req.Header.Get("Api-Key")
						return successfulAzureResponse(req), nil
					})}),
				)

				var res map[string]any
				if err := client.Execute(context.Background(), http.MethodGet, "models", nil, &res); err != nil {
					t.Fatal(err)
				}
				if middlewareAuthorization != "" {
					t.Fatalf("middleware Authorization = %q, want empty before Azure authentication", middlewareAuthorization)
				}
				if authMode.wantAPIKey == "" && middlewareAPIKey != "" {
					t.Fatalf("middleware Api-Key = %q, want empty", middlewareAPIKey)
				}
				if transportAuthorization != authMode.wantAuthorization {
					t.Fatalf("transport Authorization = %q, want %q", transportAuthorization, authMode.wantAuthorization)
				}
				if transportAPIKey != authMode.wantAPIKey {
					t.Fatalf("transport Api-Key = %q, want %q", transportAPIKey, authMode.wantAPIKey)
				}
			})
		}
	}
}

func TestAzureAuthenticationReplacesInheritedOpenAICredentials(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("OPENAI_ADMIN_KEY", "")
	credentials := []struct {
		name   string
		option func() option.RequestOption
	}{
		{name: "OpenAI API key", option: func() option.RequestOption { return option.WithAPIKey("openai-api-key") }},
		{name: "OpenAI admin API key", option: func() option.RequestOption { return option.WithAdminAPIKey("openai-admin-key") }},
		{name: "custom Authorization header", option: func() option.RequestOption {
			return option.WithHeader("Authorization", "Bearer custom-token")
		}},
	}
	authModes := []struct {
		name              string
		option            func() option.RequestOption
		wantAuthorization string
		wantAPIKey        string
	}{
		{
			name:       "Azure API key",
			option:     func() option.RequestOption { return WithAPIKey("azure-api-key") },
			wantAPIKey: "azure-api-key",
		},
		{
			name:              "Azure token credential",
			option:            func() option.RequestOption { return WithTokenCredential(&fake.TokenCredential{}) },
			wantAuthorization: "Bearer fake_token",
		},
	}
	entrypoints := []string{"request", "service"}

	for _, credential := range credentials {
		for _, authMode := range authModes {
			for _, entrypoint := range entrypoints {
				t.Run(credential.name+"/"+authMode.name+"/"+entrypoint, func(t *testing.T) {
					var middlewareAuthorization string
					var transportAuthorization string
					var transportAPIKey string
					base := openai.NewClient(
						credential.option(),
						option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
							middlewareAuthorization = req.Header.Get("Authorization")
							return next(req)
						}),
						option.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
							transportAuthorization = req.Header.Get("Authorization")
							transportAPIKey = req.Header.Get("Api-Key")
							return successfulAzureResponse(req), nil
						})}),
					)
					azureOptions := []option.RequestOption{
						WithEndpoint("https://my-resource.openai.azure.com", "2024-10-21"),
						authMode.option(),
					}

					var err error
					if entrypoint == "request" {
						var res map[string]any
						err = base.Execute(context.Background(), http.MethodGet, "models", nil, &res, azureOptions...)
					} else {
						serviceOptions := append(slices.Clone(base.Options), azureOptions...)
						service := openai.NewModelService(serviceOptions...)
						_, err = service.List(context.Background())
					}
					if err != nil {
						t.Fatal(err)
					}
					if middlewareAuthorization != "" {
						t.Fatalf("middleware Authorization = %q, want empty before Azure authentication", middlewareAuthorization)
					}
					if transportAuthorization != authMode.wantAuthorization {
						t.Fatalf("transport Authorization = %q, want %q", transportAuthorization, authMode.wantAuthorization)
					}
					if transportAPIKey != authMode.wantAPIKey {
						t.Fatalf("transport Api-Key = %q, want %q", transportAPIKey, authMode.wantAPIKey)
					}
				})
			}
		}
	}
}

func TestAzureRejectsAmbiguousCredentialsBeforeMiddleware(t *testing.T) {
	tests := []struct {
		name    string
		options []option.RequestOption
		message string
	}{
		{
			name:    "Azure API key and token credential",
			options: []option.RequestOption{WithAPIKey("azure-api-key"), WithTokenCredential(&fake.TokenCredential{})},
			message: "authentication is ambiguous",
		},
		{
			name:    "Azure API key and OpenAI API key",
			options: []option.RequestOption{WithAPIKey("azure-api-key"), option.WithAPIKey("openai-api-key")},
			message: "cannot be combined",
		},
		{
			name:    "OpenAI API key and Azure API key",
			options: []option.RequestOption{option.WithAPIKey("openai-api-key"), WithAPIKey("azure-api-key")},
			message: "cannot be combined",
		},
		{
			name:    "Azure token credential and OpenAI admin key",
			options: []option.RequestOption{WithTokenCredential(&fake.TokenCredential{}), option.WithAdminAPIKey("openai-admin-key")},
			message: "cannot be combined",
		},
		{
			name:    "OpenAI admin key and Azure token credential",
			options: []option.RequestOption{option.WithAdminAPIKey("openai-admin-key"), WithTokenCredential(&fake.TokenCredential{})},
			message: "cannot be combined",
		},
		{
			name:    "Azure token credential and Api-Key header",
			options: []option.RequestOption{WithTokenCredential(&fake.TokenCredential{}), option.WithHeader("Api-Key", "other-azure-key")},
			message: "cannot be combined",
		},
		{
			name:    "Azure token credential and Authorization header",
			options: []option.RequestOption{WithTokenCredential(&fake.TokenCredential{}), option.WithHeader("Authorization", "Bearer other-token")},
			message: "cannot be combined",
		},
		{
			name:    "Authorization header and Azure token credential",
			options: []option.RequestOption{option.WithHeader("Authorization", "Bearer other-token"), WithTokenCredential(&fake.TokenCredential{})},
			message: "cannot be combined",
		},
		{
			name:    "Azure API key and Authorization header",
			options: []option.RequestOption{WithAPIKey("azure-api-key"), option.WithHeader("Authorization", "Bearer other-token")},
			message: "cannot be combined",
		},
		{
			name:    "Authorization header and Azure API key",
			options: []option.RequestOption{option.WithHeader("Authorization", "Bearer other-token"), WithAPIKey("azure-api-key")},
			message: "cannot be combined",
		},
		{
			name: "blank-leading Authorization values",
			options: []option.RequestOption{
				WithTokenCredential(&fake.TokenCredential{}),
				option.WithHeader("Authorization", ""),
				option.WithHeaderAdd("Authorization", "Bearer other-token"),
			},
			message: "cannot be combined",
		},
		{
			name: "blank-leading Api-Key values",
			options: []option.RequestOption{
				WithTokenCredential(&fake.TokenCredential{}),
				option.WithHeader("Api-Key", ""),
				option.WithHeaderAdd("Api-Key", "other-azure-key"),
			},
			message: "cannot be combined",
		},
		{
			name:    "multiple API key values",
			options: []option.RequestOption{WithAPIKey("azure-api-key"), option.WithHeaderAdd("Api-Key", "other-azure-key")},
			message: "exactly one non-empty API key",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("OPENAI_API_KEY", "ambient-openai-key")
			t.Setenv("OPENAI_ADMIN_KEY", "ambient-openai-admin-key")
			middlewareCalls := 0
			transportCalls := 0
			opts := []option.RequestOption{WithEndpoint("https://my-resource.openai.azure.com", "2024-10-21")}
			opts = append(opts, test.options...)
			opts = append(opts,
				option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
					middlewareCalls++
					return next(req)
				}),
				option.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
					transportCalls++
					return successfulAzureResponse(req), nil
				})}),
			)
			client := openai.NewClient(opts...)

			var res map[string]any
			err := client.Execute(context.Background(), http.MethodGet, "models", nil, &res)
			if err == nil || !strings.Contains(err.Error(), test.message) {
				t.Fatalf("error = %v, want substring %q", err, test.message)
			}
			if middlewareCalls != 0 || transportCalls != 0 {
				t.Fatalf("middleware calls = %d, transport calls = %d", middlewareCalls, transportCalls)
			}
		})
	}
}

func TestAzureAuthenticationRequestOptionOverridesInheritedMode(t *testing.T) {
	tests := []struct {
		name              string
		clientAuth        option.RequestOption
		requestAuth       option.RequestOption
		wantAuthorization string
		wantAPIKey        string
	}{
		{
			name:              "token credential replaces API key",
			clientAuth:        WithAPIKey("client-api-key"),
			requestAuth:       WithTokenCredential(&fake.TokenCredential{}),
			wantAuthorization: "Bearer fake_token",
		},
		{
			name:        "API key replaces token credential",
			clientAuth:  WithTokenCredential(&fake.TokenCredential{}),
			requestAuth: WithAPIKey("request-api-key"),
			wantAPIKey:  "request-api-key",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var authorization string
			var apiKey string
			client := openai.NewClient(
				WithEndpoint("https://my-resource.openai.azure.com", "2024-10-21"),
				test.clientAuth,
				option.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
					authorization = req.Header.Get("Authorization")
					apiKey = req.Header.Get("Api-Key")
					return successfulAzureResponse(req), nil
				})}),
			)

			var res map[string]any
			if err := client.Execute(context.Background(), http.MethodGet, "models", nil, &res, test.requestAuth); err != nil {
				t.Fatal(err)
			}
			if authorization != test.wantAuthorization {
				t.Fatalf("Authorization = %q, want %q", authorization, test.wantAuthorization)
			}
			if apiKey != test.wantAPIKey {
				t.Fatalf("Api-Key = %q, want %q", apiKey, test.wantAPIKey)
			}
		})
	}
}

func successfulAzureResponse(req *http.Request) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": {"application/json"}},
		Body:       io.NopCloser(strings.NewReader(`{"ok":true}`)),
		Request:    req,
	}
}

type countingTokenCredential struct {
	calls atomic.Int32
}

func (c *countingTokenCredential) GetToken(context.Context, policy.TokenRequestOptions) (azcore.AccessToken, error) {
	c.calls.Add(1)
	return azcore.AccessToken{Token: "counting_token", ExpiresOn: time.Unix(1<<31, 0)}, nil
}

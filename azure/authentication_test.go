package azure

import (
	"compress/gzip"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/openai/openai-go/v3"
	openaiauth "github.com/openai/openai-go/v3/auth"
	"github.com/openai/openai-go/v3/option"
)

type compressionDisabledRoundTripper struct {
	http.RoundTripper
}

func (compressionDisabledRoundTripper) CompressionDisabled() bool { return true }

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

func TestAzureTokenCredentialBodyTimeoutBoundsRetryableResponse(t *testing.T) {
	tests := []struct {
		name          string
		overrideRetry bool
	}{
		{name: "default Azure retry policy"},
		{name: "caller Azure retry override", overrideRetry: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var attempts atomic.Int32
			server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				attempts.Add(1)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				flusher, ok := w.(http.Flusher)
				if !ok {
					t.Error("response writer does not support flushing")
					return
				}
				flusher.Flush()
				<-req.Context().Done()
			}))
			t.Cleanup(server.Close)

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			if test.overrideRetry {
				ctx = policy.WithRetryOptions(ctx, policy.RetryOptions{MaxRetries: 1})
			}
			client := openai.NewClient(
				WithEndpoint(server.URL, "2024-10-21"),
				WithTokenCredential(&fake.TokenCredential{}),
				option.WithMaxRetries(0),
				option.WithResponseBodyTimeout(20*time.Millisecond),
				option.WithHTTPClient(server.Client()),
			)

			var response map[string]any
			err := client.Execute(ctx, http.MethodGet, "models", nil, &response)
			if !errors.Is(err, context.DeadlineExceeded) {
				t.Fatalf("Execute() error = %v, want response body deadline exceeded", err)
			}
			if ctx.Err() != nil {
				t.Fatal("Azure retry consumed the response body until the caller safety deadline")
			}
			if !strings.Contains(err.Error(), "response body read timed out") {
				t.Fatalf("Execute() error = %v, want SDK response body timeout", err)
			}
			if got := attempts.Load(); got != 1 {
				t.Fatalf("request attempts = %d, want 1", got)
			}
		})
	}
}

func TestAzureTokenCredentialUsesSDKRetryPolicy(t *testing.T) {
	var attempts atomic.Int32
	retryCounts := make(chan string, 2)
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		retryCounts <- req.Header.Get("X-Stainless-Retry-Count")
		w.Header().Set("Content-Type", "application/json")
		if attempts.Add(1) == 1 {
			w.Header().Set("Retry-After-Ms", "1")
			w.WriteHeader(http.StatusTooManyRequests)
			if _, err := io.WriteString(w, `{"error":{"message":"retry"}}`); err != nil {
				t.Errorf("write retry response: %v", err)
			}
			return
		}
		if _, err := io.WriteString(w, `{"ok":true}`); err != nil {
			t.Errorf("write success response: %v", err)
		}
	}))
	t.Cleanup(server.Close)

	client := openai.NewClient(
		WithEndpoint(server.URL, "2024-10-21"),
		WithTokenCredential(&fake.TokenCredential{}),
		option.WithMaxRetries(1),
		option.WithMaxRetryDelay(time.Millisecond),
		option.WithHTTPClient(server.Client()),
	)
	var response map[string]any
	if err := client.Execute(context.Background(), http.MethodGet, "models", nil, &response); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if got := attempts.Load(); got != 2 {
		t.Fatalf("request attempts = %d, want 2", got)
	}
	for _, want := range []string{"0", "1"} {
		if got := <-retryCounts; got != want {
			t.Fatalf("X-Stainless-Retry-Count = %q, want %q", got, want)
		}
	}
}

func TestAzureAuthenticationPreservesDisabledNativeCompression(t *testing.T) {
	authModes := []struct {
		name    string
		auth    option.RequestOption
		wrapped bool
	}{
		{name: "API key/native transport", auth: WithAPIKey("azure-api-key")},
		{name: "API key/wrapped transport", auth: WithAPIKey("azure-api-key"), wrapped: true},
		{name: "token credential/native transport", auth: WithTokenCredential(&fake.TokenCredential{})},
		{name: "token credential/wrapped transport", auth: WithTokenCredential(&fake.TokenCredential{}), wrapped: true},
	}

	for _, authMode := range authModes {
		t.Run(authMode.name, func(t *testing.T) {
			var acceptEncoding atomicString
			server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				acceptEncoding.Store(req.Header.Get("Accept-Encoding"))
				w.Header().Set("Content-Type", "application/json")
				if strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") {
					w.Header().Set("Content-Encoding", "gzip")
					compressed := gzip.NewWriter(w)
					if _, err := compressed.Write([]byte("{}")); err != nil {
						t.Errorf("write gzip response: %v", err)
						return
					}
					if err := compressed.Close(); err != nil {
						t.Errorf("close gzip response: %v", err)
					}
					return
				}
				if _, err := io.WriteString(w, "{}"); err != nil {
					t.Errorf("write response: %v", err)
				}
			}))
			t.Cleanup(server.Close)

			serverTransport, ok := server.Client().Transport.(*http.Transport)
			if !ok {
				t.Fatalf("server transport = %T, want *http.Transport", server.Client().Transport)
			}
			transport := serverTransport.Clone()
			transport.DisableCompression = true
			var selectedTransport http.RoundTripper = transport
			if authMode.wrapped {
				selectedTransport = compressionDisabledRoundTripper{RoundTripper: transport}
			}
			client := openai.NewClient(
				WithEndpoint(server.URL, "2024-10-21"),
				authMode.auth,
				option.WithMaxRetries(0),
				option.WithMaxResponseBodyBytes(2),
				option.WithHTTPClient(&http.Client{Transport: selectedTransport}),
			)

			var response map[string]any
			if err := client.Execute(context.Background(), http.MethodGet, "models", nil, &response); err != nil {
				t.Fatalf("Execute() error = %v", err)
			}
			if got := acceptEncoding.Load(); got != "" {
				t.Fatalf("Accept-Encoding = %q, want no compression negotiation", got)
			}
		})
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

func TestAzureAuthenticationReplacesInheritedOpenAIAuthentication(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("OPENAI_ADMIN_KEY", "")
	credentials := []struct {
		name   string
		option func() (option.RequestOption, *countingSubjectTokenProvider)
	}{
		{name: "OpenAI API key", option: func() (option.RequestOption, *countingSubjectTokenProvider) {
			return option.WithAPIKey("openai-api-key"), nil
		}},
		{name: "OpenAI admin API key", option: func() (option.RequestOption, *countingSubjectTokenProvider) {
			return option.WithAdminAPIKey("openai-admin-key"), nil
		}},
		{name: "custom Authorization header", option: func() (option.RequestOption, *countingSubjectTokenProvider) {
			return option.WithHeader("Authorization", "Bearer custom-token"), nil
		}},
		{name: "OpenAI workload identity", option: func() (option.RequestOption, *countingSubjectTokenProvider) {
			provider := &countingSubjectTokenProvider{}
			return withWorkloadIdentity(provider), provider
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
					credentialOption, subjectTokenProvider := credential.option()
					var middlewareAuthorization string
					var transportAuthorization string
					var transportAPIKey string
					base := openai.NewClient(
						credentialOption,
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
					if subjectTokenProvider != nil && subjectTokenProvider.calls.Load() != 0 {
						t.Fatalf("subject token calls = %d, want 0", subjectTokenProvider.calls.Load())
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
			name:    "Api-Key header and Azure token credential",
			options: []option.RequestOption{option.WithHeader("Api-Key", "other-azure-key"), WithTokenCredential(&fake.TokenCredential{})},
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

func TestAzureRejectsWorkloadIdentityBeforeMiddleware(t *testing.T) {
	authModes := []struct {
		name   string
		option func() option.RequestOption
	}{
		{name: "Azure API key", option: func() option.RequestOption { return WithAPIKey("azure-api-key") }},
		{name: "Azure token credential", option: func() option.RequestOption {
			return WithTokenCredential(&fake.TokenCredential{})
		}},
	}

	for _, authMode := range authModes {
		for _, workloadIdentityFirst := range []bool{false, true} {
			order := "Azure auth first"
			if workloadIdentityFirst {
				order = "workload identity first"
			}
			t.Run(authMode.name+"/"+order, func(t *testing.T) {
				provider := &countingSubjectTokenProvider{}
				workloadIdentity := withWorkloadIdentity(provider)
				authOptions := []option.RequestOption{authMode.option(), workloadIdentity}
				if workloadIdentityFirst {
					authOptions[0], authOptions[1] = authOptions[1], authOptions[0]
				}

				middlewareCalls := 0
				transportCalls := 0
				opts := []option.RequestOption{WithEndpoint("https://my-resource.openai.azure.com", "2024-10-21")}
				opts = append(opts, authOptions...)
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
				if err == nil || !strings.Contains(err.Error(), "authentication is ambiguous") {
					t.Fatalf("error = %v", err)
				}
				if provider.calls.Load() != 0 || middlewareCalls != 0 || transportCalls != 0 {
					t.Fatalf(
						"subject token calls = %d, middleware calls = %d, transport calls = %d",
						provider.calls.Load(), middlewareCalls, transportCalls,
					)
				}
			})
		}
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

type countingSubjectTokenProvider struct {
	calls atomic.Int32
}

func (c *countingSubjectTokenProvider) TokenType() openaiauth.SubjectTokenType {
	return openaiauth.SubjectTokenTypeJWT
}

func (c *countingSubjectTokenProvider) GetToken(context.Context, openaiauth.HTTPDoer) (string, error) {
	c.calls.Add(1)
	return "subject-token", nil
}

func withWorkloadIdentity(provider *countingSubjectTokenProvider) option.RequestOption {
	return option.WithWorkloadIdentity(openaiauth.WorkloadIdentity{
		IdentityProviderID: "test-idp-id",
		ServiceAccountID:   "test-service-account-id",
		Provider:           provider,
	})
}

func (c *countingTokenCredential) GetToken(context.Context, policy.TokenRequestOptions) (azcore.AccessToken, error) {
	c.calls.Add(1)
	return azcore.AccessToken{Token: "counting_token", ExpiresOn: time.Unix(1<<31, 0)}, nil
}

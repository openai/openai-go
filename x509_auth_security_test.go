package openai_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/azure"
	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
)

type countingAzureTokenCredential struct {
	calls atomic.Int32
}

func (c *countingAzureTokenCredential) GetToken(context.Context, policy.TokenRequestOptions) (azcore.AccessToken, error) {
	c.calls.Add(1)
	return azcore.AccessToken{Token: "azure-provider-token", ExpiresOn: time.Now().Add(time.Hour)}, nil
}

func TestClientX509WorkloadIdentityRejectsUnsafeBaseURLBeforeExchange(t *testing.T) {
	testCases := []struct {
		name    string
		baseURL string
		fromEnv bool
	}{
		{name: "explicit plaintext", baseURL: "http://plaintext.example/v1"},
		{name: "environment plaintext", baseURL: "http://plaintext.example/v1", fromEnv: true},
		{name: "scheme relative", baseURL: "//scheme-relative.example/v1"},
		{name: "opaque", baseURL: "https:opaque-api"},
		{name: "userinfo", baseURL: "https://mtls.api.openai.com.@attacker.invalid/v1"},
		{name: "backslash userinfo", baseURL: `https://mtls.api.openai.com\@attacker.invalid/v1`},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var calls atomic.Int32
			httpClient := &http.Client{Transport: &closureTransport{fn: func(*http.Request) (*http.Response, error) {
				calls.Add(1)
				return nil, errors.New("HTTP request must not be attempted")
			}}}
			opts := []option.RequestOption{
				option.WithX509WorkloadIdentity(clientX509WorkloadIdentity()),
				option.WithHTTPClient(httpClient),
			}
			if testCase.fromEnv {
				t.Setenv("OPENAI_BASE_URL", testCase.baseURL)
			} else {
				opts = append(opts, option.WithBaseURL(testCase.baseURL))
			}
			client := openai.NewClient(opts...)

			if _, err := client.Models.List(t.Context()); err == nil {
				t.Fatal("Models.List() error = nil")
			}
			if got := calls.Load(); got != 0 {
				t.Fatalf("HTTP calls = %d, want 0", got)
			}
		})
	}

	t.Run("method option plaintext", func(t *testing.T) {
		var calls atomic.Int32
		httpClient := &http.Client{Transport: &closureTransport{fn: func(*http.Request) (*http.Response, error) {
			calls.Add(1)
			return nil, errors.New("HTTP request must not be attempted")
		}}}
		client := openai.NewClient(
			option.WithX509WorkloadIdentity(clientX509WorkloadIdentity()),
			option.WithHTTPClient(httpClient),
		)
		if _, err := client.Models.List(t.Context(), option.WithBaseURL("http://plaintext.example/v1")); err == nil {
			t.Fatal("Models.List() error = nil")
		}
		if got := calls.Load(); got != 0 {
			t.Fatalf("HTTP calls = %d, want 0", got)
		}
	})
}

func TestClientX509WorkloadIdentityRejectsProviderOwnedBaseURLBeforeExchange(t *testing.T) {
	testCases := []struct {
		name    string
		baseURL string
		fromEnv bool
	}{
		{name: "Azure OpenAI", baseURL: "https://resource.openai.azure.com/openai/v1"},
		{name: "Azure absolute DNS name", baseURL: "https://resource.openai.azure.com./openai/v1"},
		{name: "Azure US sovereign", baseURL: "https://resource.openai.azure.us/openai/v1"},
		{name: "Azure China sovereign", baseURL: "https://resource.openai.azure.cn/openai/v1"},
		{name: "Azure Cognitive Services", baseURL: "https://resource.cognitiveservices.azure.com/openai/v1"},
		{name: "Azure Foundry", baseURL: "https://resource.services.ai.azure.com/openai/v1"},
		{name: "Azure API Management", baseURL: "https://resource.azure-api.net/openai/v1"},
		{name: "Bedrock runtime", baseURL: "https://bedrock-runtime.us-east-1.amazonaws.com/openai/v1"},
		{name: "Bedrock control plane", baseURL: "https://bedrock.us-east-1.amazonaws.com/openai/v1"},
		{name: "Bedrock control plane FIPS", baseURL: "https://bedrock-fips.us-east-1.amazonaws.com/openai/v1"},
		{name: "Bedrock agent", baseURL: "https://bedrock-agent.us-east-1.amazonaws.com/openai/v1"},
		{name: "Bedrock agent FIPS", baseURL: "https://bedrock-agent-fips.us-east-1.amazonaws.com/openai/v1"},
		{name: "Bedrock agent runtime", baseURL: "https://bedrock-agent-runtime.us-east-1.amazonaws.com/openai/v1"},
		{name: "Bedrock agent runtime FIPS", baseURL: "https://bedrock-agent-runtime-fips.us-east-1.amazonaws.com/openai/v1"},
		{name: "Bedrock data automation", baseURL: "https://bedrock-data-automation.us-east-1.amazonaws.com/openai/v1"},
		{name: "Bedrock data automation runtime FIPS", baseURL: "https://bedrock-data-automation-runtime-fips.us-east-1.api.aws/openai/v1"},
		{name: "Bedrock AgentCore control", baseURL: "https://bedrock-agentcore-control.us-east-1.amazonaws.com/openai/v1"},
		{name: "Bedrock AgentCore gateway", baseURL: "https://gateway-id.gateway.bedrock-agentcore.us-east-1.amazonaws.com/openai/v1"},
		{name: "Bedrock China", baseURL: "https://bedrock-runtime.cn-north-1.amazonaws.com.cn/openai/v1"},
		{name: "Bedrock EU sovereign", baseURL: "https://bedrock-runtime.eusc-de-east-1.amazonaws.eu/openai/v1"},
		{name: "Bedrock mantle", baseURL: "https://bedrock-mantle.us-east-1.api.aws/openai/v1"},
		{name: "Bedrock PrivateLink runtime", baseURL: "https://vpce-0123456789abcdef0.bedrock-runtime.us-east-1.vpce.amazonaws.com/openai/v1"},
		{name: "Bedrock PrivateLink mantle", baseURL: "https://vpce-0123456789abcdef0.bedrock-mantle.us-east-1.vpce.amazonaws.com/openai/v1"},
		{name: "environment provider URL", baseURL: "https://resource.openai.azure.com/openai/v1", fromEnv: true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var calls atomic.Int32
			httpClient := &http.Client{Transport: &closureTransport{fn: func(*http.Request) (*http.Response, error) {
				calls.Add(1)
				return nil, errors.New("HTTP request must not be attempted")
			}}}
			opts := []option.RequestOption{
				option.WithX509WorkloadIdentity(clientX509WorkloadIdentity()),
				option.WithHTTPClient(httpClient),
			}
			if testCase.fromEnv {
				t.Setenv("OPENAI_BASE_URL", testCase.baseURL)
			} else {
				opts = append(opts, option.WithBaseURL(testCase.baseURL))
			}

			client := openai.NewClient(opts...)
			if _, err := client.Models.List(t.Context()); err == nil || !strings.Contains(err.Error(), "provider-owned") {
				t.Fatalf("Models.List() error = %v, want provider-owned URL rejection", err)
			}
			if got := calls.Load(); got != 0 {
				t.Fatalf("HTTP calls = %d, want 0", got)
			}
		})
	}
}

func TestClientX509WorkloadIdentityAllowsProviderSuffixLookalikeGateway(t *testing.T) {
	testCases := []struct {
		name    string
		baseURL string
	}{
		{name: "Azure suffix", baseURL: "https://resource.openai.azure.com.example/v1"},
		{name: "Bedrock PrivateLink suffix", baseURL: "https://vpce-0123456789abcdef0.bedrock-runtime.us-east-1.vpce.amazonaws.com.example/v1"},
		{name: "Bedrock PrivateLink service", baseURL: "https://vpce-0123456789abcdef0.bedrock-runtime-gateway.us-east-1.vpce.amazonaws.com/v1"},
		{name: "Bedrock PrivateLink endpoint ID", baseURL: "https://gateway.bedrock-runtime.us-east-1.vpce.amazonaws.com/v1"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var exchangeCalls atomic.Int32
			var apiCalls atomic.Int32
			httpClient := nativeX509HTTPClient(t, func(req *http.Request) (*http.Response, error) {
				if req.URL.Hostname() == "mtls.auth.openai.com" {
					exchangeCalls.Add(1)
					return &http.Response{
						StatusCode: http.StatusOK,
						Header:     http.Header{"Content-Type": []string{"application/json"}},
						Body:       io.NopCloser(strings.NewReader(`{"access_token":"x509-token","expires_in":3600}`)),
					}, nil
				}
				apiCalls.Add(1)
				return modelsListResponse(), nil
			})
			client := openai.NewClient(
				option.WithBaseURL(testCase.baseURL),
				option.WithX509WorkloadIdentity(clientX509WorkloadIdentity()),
				option.WithHTTPClient(httpClient),
			)

			if _, err := client.Models.List(t.Context()); err != nil {
				t.Fatalf("Models.List() error = %v", err)
			}
			if got, want := exchangeCalls.Load(), int32(1); got != want {
				t.Fatalf("exchange calls = %d, want %d", got, want)
			}
			if got, want := apiCalls.Load(), int32(1); got != want {
				t.Fatalf("API calls = %d, want %d", got, want)
			}
		})
	}
}

func TestClientX509WorkloadIdentityRejectsBedrockPrivateLinkBeforeExchange(t *testing.T) {
	var calls atomic.Int32
	httpClient := nativeX509HTTPClient(t, func(req *http.Request) (*http.Response, error) {
		calls.Add(1)
		if req.URL.Hostname() == "mtls.auth.openai.com" {
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       io.NopCloser(strings.NewReader(`{"access_token":"x509-token","expires_in":3600}`)),
			}, nil
		}
		return modelsListResponse(), nil
	})
	client := openai.NewClient(
		option.WithBaseURL("https://vpce-0123456789abcdef0.bedrock-runtime.us-east-1.vpce.amazonaws.com/openai/v1"),
		option.WithX509WorkloadIdentity(clientX509WorkloadIdentity()),
		option.WithHTTPClient(httpClient),
	)

	if _, err := client.Models.List(t.Context()); err == nil || !strings.Contains(err.Error(), "provider-owned") {
		t.Errorf("Models.List() error = %v, want provider-owned URL rejection", err)
	}
	if got := calls.Load(); got != 0 {
		t.Fatalf("HTTP calls = %d, want 0", got)
	}
}

func TestClientX509WorkloadIdentityBindsAbsoluteURLsToConfiguredOrigin(t *testing.T) {
	testCases := []struct {
		name       string
		requestURL string
		useGet     bool
		wantOK     bool
	}{
		{name: "plaintext", requestURL: "http://trusted.example/v1/models"},
		{name: "plaintext public Get", requestURL: "http://trusted.example/v1/models", useGet: true},
		{name: "different TLS origin", requestURL: "https://attacker.invalid/v1/models"},
		{name: "userinfo lookalike", requestURL: "https://trusted.example@attacker.invalid/v1/models"},
		{name: "normalized configured origin", requestURL: "https://TRUSTED.EXAMPLE/v1/models", wantOK: true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var exchangeCalls atomic.Int32
			var apiCalls atomic.Int32
			httpClient := nativeX509HTTPClient(t, func(req *http.Request) (*http.Response, error) {
				if req.URL.Hostname() == "mtls.auth.openai.com" {
					exchangeCalls.Add(1)
					return &http.Response{
						StatusCode: http.StatusOK,
						Header:     http.Header{"Content-Type": []string{"application/json"}},
						Body:       io.NopCloser(strings.NewReader(`{"access_token":"x509-token","expires_in":3600}`)),
					}, nil
				}
				apiCalls.Add(1)
				return modelsListResponse(), nil
			})
			client := openai.NewClient(
				option.WithBaseURL("https://trusted.example:443/v1"),
				option.WithX509WorkloadIdentity(clientX509WorkloadIdentity()),
				option.WithHTTPClient(httpClient),
			)

			var err error
			if testCase.useGet {
				err = client.Get(t.Context(), testCase.requestURL, nil, nil)
			} else {
				err = client.Execute(t.Context(), http.MethodGet, testCase.requestURL, nil, nil)
			}
			if testCase.wantOK {
				if err != nil {
					t.Fatalf("Execute() error = %v", err)
				}
				if got := exchangeCalls.Load(); got != 1 {
					t.Fatalf("exchange calls = %d, want 1", got)
				}
				if got := apiCalls.Load(); got != 1 {
					t.Fatalf("API calls = %d, want 1", got)
				}
				return
			}
			if err == nil {
				t.Fatal("Execute() error = nil")
			}
			if got := exchangeCalls.Load(); got != 0 {
				t.Fatalf("exchange calls = %d, want 0", got)
			}
			if got := apiCalls.Load(); got != 0 {
				t.Fatalf("API calls = %d, want 0", got)
			}
		})
	}
}

func TestClientX509WorkloadIdentityGuardsFinalDispatch(t *testing.T) {
	t.Run("origin mutation", func(t *testing.T) {
		var apiCalls atomic.Int32
		httpClient := nativeX509HTTPClient(t, func(req *http.Request) (*http.Response, error) {
			if req.URL.Hostname() == "mtls.auth.openai.com" {
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     http.Header{"Content-Type": []string{"application/json"}},
					Body:       io.NopCloser(strings.NewReader(`{"access_token":"x509-token","expires_in":3600}`)),
				}, nil
			}
			apiCalls.Add(1)
			return modelsListResponse(), nil
		})
		client := openai.NewClient(
			option.WithBaseURL("https://trusted.example/v1"),
			option.WithX509WorkloadIdentity(clientX509WorkloadIdentity()),
			option.WithHTTPClient(httpClient),
			option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
				rewritten, err := url.Parse("https://attacker.invalid/v1/models")
				if err != nil {
					return nil, err
				}
				req.URL = rewritten
				return next(req)
			}),
		)

		if _, err := client.Models.List(t.Context()); err == nil {
			t.Fatal("Models.List() error = nil")
		}
		if got := apiCalls.Load(); got != 0 {
			t.Fatalf("API calls = %d, want 0", got)
		}
	})

	t.Run("authority mutation", func(t *testing.T) {
		var apiCalls atomic.Int32
		httpClient := nativeX509HTTPClient(t, func(req *http.Request) (*http.Response, error) {
			if req.URL.Hostname() == "mtls.auth.openai.com" {
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     http.Header{"Content-Type": []string{"application/json"}},
					Body:       io.NopCloser(strings.NewReader(`{"access_token":"x509-token","expires_in":3600}`)),
				}, nil
			}
			apiCalls.Add(1)
			return modelsListResponse(), nil
		})
		client := openai.NewClient(
			option.WithBaseURL("https://trusted.example/v1"),
			option.WithX509WorkloadIdentity(clientX509WorkloadIdentity()),
			option.WithHTTPClient(httpClient),
			option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
				req.Host = "attacker.invalid"
				return next(req)
			}),
		)

		if _, err := client.Models.List(t.Context()); err == nil {
			t.Fatal("Models.List() error = nil")
		}
		if got := apiCalls.Load(); got != 0 {
			t.Fatalf("API calls = %d, want 0", got)
		}
	})

	t.Run("normalized authority", func(t *testing.T) {
		var apiCalls atomic.Int32
		httpClient := nativeX509HTTPClient(t, func(req *http.Request) (*http.Response, error) {
			if req.URL.Hostname() == "mtls.auth.openai.com" {
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     http.Header{"Content-Type": []string{"application/json"}},
					Body:       io.NopCloser(strings.NewReader(`{"access_token":"x509-token","expires_in":3600}`)),
				}, nil
			}
			apiCalls.Add(1)
			return modelsListResponse(), nil
		})
		client := openai.NewClient(
			option.WithBaseURL("https://trusted.example/v1"),
			option.WithX509WorkloadIdentity(clientX509WorkloadIdentity()),
			option.WithHTTPClient(httpClient),
			option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
				req.Host = "TRUSTED.EXAMPLE:443"
				return next(req)
			}),
		)

		if _, err := client.Models.List(t.Context()); err != nil {
			t.Fatalf("Models.List() error = %v", err)
		}
		if got := apiCalls.Load(); got != 1 {
			t.Fatalf("API calls = %d, want 1", got)
		}
	})

	t.Run("authorization provenance", func(t *testing.T) {
		var authorization string
		var exchangeAuthorization string
		var exchangeCookie string
		httpClient := nativeX509HTTPClient(t, func(req *http.Request) (*http.Response, error) {
			if req.URL.Hostname() == "mtls.auth.openai.com" {
				exchangeAuthorization = req.Header.Get("Authorization")
				exchangeCookie = req.Header.Get("Cookie")
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     http.Header{"Content-Type": []string{"application/json"}},
					Body:       io.NopCloser(strings.NewReader(`{"access_token":"x509-token","expires_in":3600}`)),
				}, nil
			}
			authorization = req.Header.Get("Authorization")
			return modelsListResponse(), nil
		})
		client := openai.NewClient(
			option.WithBaseURL("https://trusted.example/v1"),
			option.WithX509WorkloadIdentity(clientX509WorkloadIdentity()),
			option.WithHTTPClient(httpClient),
			option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
				req.Header.Set("Authorization", "Bearer attacker-override")
				req.Header.Set("Cookie", "api-session=secret")
				return next(req)
			}),
		)

		if _, err := client.Models.List(t.Context()); err == nil {
			t.Fatal("Models.List() error = nil")
		}
		if authorization != "" {
			t.Fatalf("API Authorization = %q, want no dispatch", authorization)
		}
		if exchangeAuthorization != "" || exchangeCookie != "" {
			t.Fatalf("token exchange inherited API middleware headers: Authorization=%q Cookie=%q", exchangeAuthorization, exchangeCookie)
		}
	})

	for _, headerName := range []string{
		"Api-Key", "X-Api-Key", "api-key", "x-api-key",
		"api_key", "API_KEY", "x_api_key", "X_API_KEY", "x_api-key", "x-api_key",
		"Proxy-Authorization", "proxy-authorization", "PROXY_AUTHORIZATION", "proxy_authorization", "Proxy_Authorization",
	} {
		t.Run("secondary API credential "+headerName, func(t *testing.T) {
			var exchangeCalls atomic.Int32
			var apiCalls atomic.Int32
			httpClient := nativeX509HTTPClient(t, func(req *http.Request) (*http.Response, error) {
				if req.URL.Hostname() == "mtls.auth.openai.com" {
					exchangeCalls.Add(1)
					return &http.Response{
						StatusCode: http.StatusOK,
						Header:     http.Header{"Content-Type": []string{"application/json"}},
						Body:       io.NopCloser(strings.NewReader(`{"access_token":"x509-token","expires_in":3600}`)),
					}, nil
				}
				apiCalls.Add(1)
				return modelsListResponse(), nil
			})
			client := openai.NewClient(
				option.WithBaseURL("https://trusted.example/v1"),
				option.WithX509WorkloadIdentity(clientX509WorkloadIdentity()),
				option.WithHTTPClient(httpClient),
				option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
					// Assign the raw map key to cover non-canonical header names
					// that Header.Get does not find but net/http can transmit.
					req.Header[headerName] = []string{"other-credential"}
					return next(req)
				}),
			)

			if _, err := client.Models.List(t.Context()); err == nil {
				t.Fatal("Models.List() error = nil")
			}
			if got := exchangeCalls.Load(); got != 1 {
				t.Fatalf("exchange calls = %d, want 1", got)
			}
			if got := apiCalls.Load(); got != 0 {
				t.Fatalf("API calls = %d, want 0", got)
			}
		})
	}
}

func TestDefaultHTTPTransportUsesRequestHostOverride(t *testing.T) {
	hosts := make(chan string, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		hosts <- req.Host
		w.WriteHeader(http.StatusNoContent)
	}))
	t.Cleanup(server.Close)
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, server.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Host = "attacker.invalid"
	resp, err := server.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close()
	if got, want := <-hosts, "attacker.invalid"; got != want {
		t.Fatalf("server-observed Host = %q, want %q", got, want)
	}
}

func TestClientX509WorkloadIdentityPreservesEndpointAuthorizationProvenance(t *testing.T) {
	newClient := func(middleware option.Middleware) (openai.Client, *atomic.Int32, *atomic.Int32, *string) {
		var exchangeCalls atomic.Int32
		var apiCalls atomic.Int32
		var authorization string
		httpClient := &http.Client{Transport: &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
			if req.URL.Hostname() == "mtls.auth.openai.com" {
				exchangeCalls.Add(1)
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     http.Header{"Content-Type": []string{"application/json"}},
					Body:       io.NopCloser(strings.NewReader(`{"access_token":"x509-token","expires_in":3600}`)),
				}, nil
			}
			apiCalls.Add(1)
			authorization = req.Header.Get("Authorization")
			return modelsListResponse(), nil
		}}}
		opts := []option.RequestOption{
			option.WithBaseURL("https://trusted.example/v1"),
			option.WithAdminAPIKey("admin-secret"),
			option.WithX509WorkloadIdentity(clientX509WorkloadIdentity()),
			option.WithHTTPClient(httpClient),
		}
		if middleware != nil {
			opts = append(opts, option.WithMiddleware(middleware))
		}
		client := openai.NewClient(opts...)
		return client, &exchangeCalls, &apiCalls, &authorization
	}

	t.Run("admin authorization", func(t *testing.T) {
		client, exchangeCalls, apiCalls, authorization := newClient(nil)
		err := client.Get(
			t.Context(),
			"organization/audit_logs",
			nil,
			nil,
			requestconfig.WithAdminAPIKeyAuthSecurity(),
		)
		if err != nil {
			t.Fatalf("Get() error = %v", err)
		}
		if got := exchangeCalls.Load(); got != 0 {
			t.Fatalf("token exchange calls = %d, want 0", got)
		}
		if got := apiCalls.Load(); got != 1 {
			t.Fatalf("API calls = %d, want 1", got)
		}
		if got, want := *authorization, "Bearer admin-secret"; got != want {
			t.Fatalf("Authorization = %q, want %q", got, want)
		}
	})

	t.Run("admin mutation", func(t *testing.T) {
		client, exchangeCalls, apiCalls, _ := newClient(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
			req.Header.Set("Authorization", "Basic customer-secret")
			return next(req)
		})
		err := client.Get(
			t.Context(),
			"organization/audit_logs",
			nil,
			nil,
			requestconfig.WithAdminAPIKeyAuthSecurity(),
		)
		if err == nil {
			t.Fatal("Get() error = nil")
		}
		if got := exchangeCalls.Load(); got != 0 {
			t.Fatalf("token exchange calls = %d, want 0", got)
		}
		if got := apiCalls.Load(); got != 0 {
			t.Fatalf("API calls = %d, want 0", got)
		}
	})

	t.Run("headerless mutation", func(t *testing.T) {
		client, exchangeCalls, apiCalls, _ := newClient(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
			req.Header["authorization"] = []string{"Basic customer-secret"}
			return next(req)
		})
		err := client.Get(
			t.Context(),
			"health",
			nil,
			nil,
			requestconfig.WithSecurity(requestconfig.Security{}),
		)
		if err == nil {
			t.Fatal("Get() error = nil")
		}
		if got := exchangeCalls.Load(); got != 0 {
			t.Fatalf("token exchange calls = %d, want 0", got)
		}
		if got := apiCalls.Load(); got != 0 {
			t.Fatalf("API calls = %d, want 0", got)
		}
	})

	t.Run("headerless request", func(t *testing.T) {
		client, exchangeCalls, apiCalls, authorization := newClient(nil)
		err := client.Get(
			t.Context(),
			"health",
			nil,
			nil,
			requestconfig.WithSecurity(requestconfig.Security{}),
		)
		if err != nil {
			t.Fatalf("Get() error = %v", err)
		}
		if got := exchangeCalls.Load(); got != 0 {
			t.Fatalf("token exchange calls = %d, want 0", got)
		}
		if got := apiCalls.Load(); got != 1 {
			t.Fatalf("API calls = %d, want 1", got)
		}
		if *authorization != "" {
			t.Fatalf("Authorization = %q, want empty", *authorization)
		}
	})
}

func TestClientX509WorkloadIdentityHonorsMethodAdminCredentialPreference(t *testing.T) {
	var exchangeCalls atomic.Int32
	var apiCalls atomic.Int32
	var authorization string
	httpClient := nativeX509HTTPClient(t, func(req *http.Request) (*http.Response, error) {
		if req.URL.Hostname() == "mtls.auth.openai.com" {
			exchangeCalls.Add(1)
			return nil, errors.New("token exchange must not be attempted")
		}
		apiCalls.Add(1)
		authorization = req.Header.Get("Authorization")
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(`{"ok":true}`)),
		}, nil
	})
	client := openai.NewClient(
		option.WithBaseURL("https://trusted.example/v1"),
		option.WithX509WorkloadIdentity(clientX509WorkloadIdentity()),
		option.WithHTTPClient(httpClient),
	)

	var result map[string]any
	err := client.Execute(
		t.Context(),
		http.MethodGet,
		"organization/audit_logs",
		nil,
		&result,
		option.WithAdminAPIKey("method-admin-secret"),
	)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if got := exchangeCalls.Load(); got != 0 {
		t.Fatalf("token exchange calls = %d, want 0", got)
	}
	if got := apiCalls.Load(); got != 1 {
		t.Fatalf("API calls = %d, want 1", got)
	}
	if got, want := authorization, "Bearer method-admin-secret"; got != want {
		t.Fatalf("Authorization = %q, want %q", got, want)
	}
}

func TestClientX509WorkloadIdentityHonorsExplicitAuthorizationDeletion(t *testing.T) {
	var exchangeCalls atomic.Int32
	var apiCalls atomic.Int32
	var authorization string
	httpClient := nativeX509HTTPClient(t, func(req *http.Request) (*http.Response, error) {
		if req.URL.Hostname() == "mtls.auth.openai.com" {
			exchangeCalls.Add(1)
			return nil, errors.New("token exchange must not be attempted")
		}
		apiCalls.Add(1)
		authorization = req.Header.Get("Authorization")
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(`{"ok":true}`)),
		}, nil
	})
	client := openai.NewClient(
		option.WithBaseURL("https://trusted.example/v1"),
		option.WithX509WorkloadIdentity(clientX509WorkloadIdentity()),
		option.WithHTTPClient(httpClient),
	)

	var result map[string]any
	err := client.Execute(
		t.Context(),
		http.MethodGet,
		"models",
		nil,
		&result,
		option.WithHeaderDel("Authorization"),
	)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if got := exchangeCalls.Load(); got != 0 {
		t.Fatalf("token exchange calls = %d, want 0", got)
	}
	if got := apiCalls.Load(); got != 1 {
		t.Fatalf("API calls = %d, want 1", got)
	}
	if authorization != "" {
		t.Fatalf("Authorization = %q, want empty", authorization)
	}
}

func TestClientRejectsAzureAuthenticationWithX509BeforeExchange(t *testing.T) {
	credential := &countingAzureTokenCredential{}
	testCases := []struct {
		name string
		auth option.RequestOption
	}{
		{name: "API key", auth: azure.WithAPIKey("azure-secret")},
		{name: "token credential", auth: azure.WithTokenCredential(credential)},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var calls atomic.Int32
			httpClient := &http.Client{Transport: &closureTransport{fn: func(*http.Request) (*http.Response, error) {
				calls.Add(1)
				return nil, errors.New("HTTP request must not be attempted")
			}}}
			client := openai.NewClient(
				azure.WithEndpoint("https://resource.openai.azure.com", "2024-10-21"),
				testCase.auth,
				option.WithX509WorkloadIdentity(clientX509WorkloadIdentity()),
				option.WithHTTPClient(httpClient),
			)

			if _, err := client.Models.List(t.Context()); err == nil {
				t.Fatal("Models.List() error = nil")
			}
			if got := calls.Load(); got != 0 {
				t.Fatalf("HTTP calls = %d, want 0", got)
			}
		})
	}
	if got := credential.calls.Load(); got != 0 {
		t.Fatalf("Azure token credential calls = %d, want 0", got)
	}
}

func TestClientRejectsAzureEndpointWithX509BeforeExchange(t *testing.T) {
	endpoint := azure.WithEndpoint("https://resource.openai.azure.com", "2024-10-21")
	x509 := option.WithX509WorkloadIdentity(clientX509WorkloadIdentity())
	testCases := []struct {
		name string
		opts []option.RequestOption
	}{
		{name: "endpoint before X.509", opts: []option.RequestOption{endpoint, x509}},
		{name: "endpoint after X.509", opts: []option.RequestOption{x509, endpoint}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var calls atomic.Int32
			httpClient := &http.Client{Transport: &closureTransport{fn: func(*http.Request) (*http.Response, error) {
				calls.Add(1)
				return nil, errors.New("HTTP request must not be attempted")
			}}}
			opts := append([]option.RequestOption{}, testCase.opts...)
			opts = append(opts, option.WithHTTPClient(httpClient))
			client := openai.NewClient(opts...)

			if _, err := client.Models.List(t.Context()); err == nil {
				t.Fatal("Models.List() error = nil")
			}
			if got := calls.Load(); got != 0 {
				t.Fatalf("HTTP calls = %d, want 0", got)
			}
		})
	}
}

func TestClientX509RejectsLaterAPIKeyCredentialsBeforeExchange(t *testing.T) {
	testCases := []struct {
		name        string
		opt         option.RequestOption
		clientLevel bool
	}{
		{name: "API key", opt: option.WithAPIKey("other-api-key")},
		{name: "method API key header", opt: option.WithHeader("api-key", "other-api-key")},
		{name: "method X-API-Key header", opt: option.WithHeader("X-API-KEY", "other-api-key")},
		{name: "method API key underscore alias", opt: option.WithHeader("API_KEY", "other-api-key")},
		{name: "method X-API-Key underscore alias", opt: option.WithHeader("x_api_key", "other-api-key")},
		{name: "method proxy authorization", opt: option.WithHeader("Proxy-Authorization", "Basic proxy-secret")},
		{name: "method proxy authorization alias", opt: option.WithHeader("proxy_authorization", "Basic proxy-secret")},
		{name: "method authorization override", opt: option.WithHeader("Authorization", "Basic customer-secret")},
		{name: "client API key header", opt: option.WithHeader("API-KEY", "other-api-key"), clientLevel: true},
		{name: "client X-API-Key header", opt: option.WithHeader("x-api-key", "other-api-key"), clientLevel: true},
		{name: "client API key underscore alias", opt: option.WithHeader("api_key", "other-api-key"), clientLevel: true},
		{name: "client X-API-Key hybrid alias", opt: option.WithHeader("X_API-Key", "other-api-key"), clientLevel: true},
		{name: "client proxy authorization", opt: option.WithHeader("PROXY-AUTHORIZATION", "Basic proxy-secret"), clientLevel: true},
		{name: "client proxy authorization alias", opt: option.WithHeader("Proxy_Authorization", "Basic proxy-secret"), clientLevel: true},
		{name: "client authorization override", opt: option.WithHeader("authorization", "Basic customer-secret"), clientLevel: true},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var calls atomic.Int32
			httpClient := &http.Client{Transport: &closureTransport{fn: func(*http.Request) (*http.Response, error) {
				calls.Add(1)
				return nil, errors.New("HTTP request must not be attempted")
			}}}
			opts := []option.RequestOption{
				option.WithX509WorkloadIdentity(clientX509WorkloadIdentity()),
				option.WithHTTPClient(httpClient),
			}
			if testCase.clientLevel {
				opts = append(opts, testCase.opt)
			}
			client := openai.NewClient(opts...)

			var requestOpts []option.RequestOption
			if !testCase.clientLevel {
				requestOpts = append(requestOpts, testCase.opt)
			}
			if _, err := client.Models.List(t.Context(), requestOpts...); err == nil {
				t.Fatal("Models.List() error = nil")
			}
			if got := calls.Load(); got != 0 {
				t.Fatalf("HTTP calls = %d, want 0", got)
			}
		})
	}
}

package openai_test

import (
	"errors"
	"net/http"
	"sync/atomic"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/auth"
	"github.com/openai/openai-go/v3/azure"
	"github.com/openai/openai-go/v3/bedrock"
	"github.com/openai/openai-go/v3/option"
)

func TestX509WorkloadIdentityAuthenticatesPublicClientOverMutualTLS(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "synthetic-ambient-api-key")
	t.Setenv("OPENAI_ADMIN_KEY", "synthetic-ambient-admin-key")
	t.Setenv("OPENAI_CUSTOM_HEADERS", "Authorization: Bearer synthetic-ambient-header")
	t.Setenv("OPENAI_ORG_ID", "synthetic-ambient-organization")
	t.Setenv("OPENAI_PROJECT_ID", "synthetic-ambient-project")
	t.Setenv("OPENAI_BASE_URL", "https://mtls.api.openai.com/v1/")
	config, issuer, api := newX509WorkloadIdentityIntegration(t)
	client := openai.NewClient(option.WithX509WorkloadIdentity(config), option.WithMaxRetries(0))
	for range 2 {
		if _, err := client.Models.List(t.Context()); err != nil {
			t.Fatalf("mutually authenticated public client request: %v", err)
		}
	}
	issuerRequests := issuer.requests()
	apiRequests := api.requests()
	if len(issuerRequests) != 2 || len(apiRequests) != 2 {
		t.Fatalf("pre-cache client made issuer/API requests = %d/%d, want 2/2", len(issuerRequests), len(apiRequests))
	}
	for _, request := range issuerRequests {
		if request.authorization != "" || request.host != x509ConformanceIssuerHost ||
			request.peerCommonName != "integrated-workload" || len(request.exchangeFields) != 4 {
			t.Errorf("issuer received unsafe workload exchange: %+v", request)
		}
	}
	for _, request := range apiRequests {
		if request.authorization != "Bearer "+x509ConformanceToken || request.host != x509ConformanceAPIHost ||
			request.peerCommonName != "integrated-workload" {
			t.Errorf("API received unsafe workload request: %+v", request)
		}
	}
}

func TestX509WorkloadIdentityRejectsUnsafeConfiguredOriginsBeforeExchange(t *testing.T) {
	for _, test := range []struct {
		name string
		env  string
		opts []option.RequestOption
	}{
		{name: "environment attacker", env: "https://attacker.example.test/v1/"},
		{name: "environment plaintext", env: "http://mtls.api.openai.com/v1/"},
		{name: "environment Azure", env: "https://resource.openai.azure.com/v1/"},
		{name: "environment EU", env: "https://eu.api.openai.com/v1/"},
		{name: "explicit attacker", opts: []option.RequestOption{option.WithBaseURL("https://attacker.example.test/v1/")}},
		{name: "explicit plaintext", opts: []option.RequestOption{option.WithBaseURL("http://mtls.api.openai.com/v1/")}},
		{name: "data residency EU", opts: []option.RequestOption{option.WithDataResidency(option.DataResidencyEU)}},
		{name: "data residency global", opts: []option.RequestOption{option.WithDataResidency(option.DataResidencyGlobal)}},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("OPENAI_BASE_URL", test.env)
			config, issuer, api := newX509WorkloadIdentityIntegration(t)
			opts := []option.RequestOption{option.WithX509WorkloadIdentity(config), option.WithMaxRetries(0)}
			opts = append(opts, test.opts...)
			client := openai.NewClient(opts...)
			if _, err := client.Models.List(t.Context()); err == nil {
				t.Fatal("unsafe workload-identity endpoint was accepted")
			}
			assertX509WorkloadNoRequests(t, issuer, api)
		})
	}
}

func TestX509WorkloadIdentityRejectsConflictingCredentialsAndClients(t *testing.T) {
	for _, test := range []struct {
		name   string
		before func(auth.X509WorkloadIdentity) []option.RequestOption
		after  func(auth.X509WorkloadIdentity) []option.RequestOption
	}{
		{name: "later API key", after: func(auth.X509WorkloadIdentity) []option.RequestOption {
			return []option.RequestOption{option.WithAPIKey("synthetic-explicit-api-key")}
		}},
		{name: "earlier API key", before: func(auth.X509WorkloadIdentity) []option.RequestOption {
			return []option.RequestOption{option.WithAPIKey("synthetic-explicit-api-key")}
		}},
		{name: "later admin key", after: func(auth.X509WorkloadIdentity) []option.RequestOption {
			return []option.RequestOption{option.WithAdminAPIKey("synthetic-explicit-admin-key")}
		}},
		{name: "custom Authorization", after: func(auth.X509WorkloadIdentity) []option.RequestOption {
			return []option.RequestOption{option.WithHeader("Authorization", "Bearer synthetic-header-token")}
		}},
		{name: "custom API key alias", after: func(auth.X509WorkloadIdentity) []option.RequestOption {
			return []option.RequestOption{option.WithHeader("x_api_key", "synthetic-header-key")}
		}},
		{name: "custom cookie", after: func(auth.X509WorkloadIdentity) []option.RequestOption {
			return []option.RequestOption{option.WithHeader("Cookie", "synthetic=session")}
		}},
		{name: "native HTTP client override", after: func(auth.X509WorkloadIdentity) []option.RequestOption {
			return []option.RequestOption{option.WithHTTPClient(&http.Client{})}
		}},
		{name: "opaque HTTP client override", after: func(auth.X509WorkloadIdentity) []option.RequestOption {
			return []option.RequestOption{option.WithHTTPClient(x509WorkloadRejectedDoer{})}
		}},
		{name: "other workload identity", after: func(config auth.X509WorkloadIdentity) []option.RequestOption {
			return []option.RequestOption{option.WithWorkloadIdentity(auth.WorkloadIdentity{
				IdentityProviderID: config.IdentityProviderID,
				ServiceAccountID:   config.ServiceAccountID,
			})}
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("OPENAI_BASE_URL", "https://mtls.api.openai.com/v1/")
			config, issuer, api := newX509WorkloadIdentityIntegration(t)
			var opts []option.RequestOption
			if test.before != nil {
				opts = append(opts, test.before(config)...)
			}
			opts = append(opts, option.WithX509WorkloadIdentity(config), option.WithMaxRetries(0))
			if test.after != nil {
				opts = append(opts, test.after(config)...)
			}
			client := openai.NewClient(opts...)
			if _, err := client.Models.List(t.Context()); err == nil {
				t.Fatal("conflicting X.509 credentials or client were accepted")
			}
			assertX509WorkloadNoRequests(t, issuer, api)
		})
	}
}

func TestX509WorkloadIdentityRejectsLateMiddlewareMutationBeforeExchange(t *testing.T) {
	for _, test := range []struct {
		name   string
		mutate func(*http.Request)
	}{
		{name: "attacker host", mutate: func(request *http.Request) { request.URL.Host = "attacker.example.test" }},
		{name: "attacker Host header", mutate: func(request *http.Request) { request.Host = "attacker.example.test" }},
		{name: "attacker Authorization", mutate: func(request *http.Request) {
			request.Header.Set("Authorization", "Bearer attacker-token")
		}},
		{name: "attacker API key", mutate: func(request *http.Request) {
			request.Header.Set("X-Api-Key", "synthetic-attacker-key")
		}},
		{name: "path traversal", mutate: func(request *http.Request) {
			request.URL.Path = "/v1/../attacker"
		}},
		{name: "credential trailer", mutate: func(request *http.Request) {
			request.Trailer = http.Header{"Authorization": []string{"Bearer attacker-token"}}
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("OPENAI_BASE_URL", "https://mtls.api.openai.com/v1/")
			config, issuer, api := newX509WorkloadIdentityIntegration(t)
			var invoked atomic.Int32
			client := openai.NewClient(
				option.WithX509WorkloadIdentity(config),
				option.WithMaxRetries(0),
				option.WithMiddleware(func(request *http.Request, next option.MiddlewareNext) (*http.Response, error) {
					invoked.Add(1)
					test.mutate(request)
					return next(request)
				}),
			)
			if _, err := client.Models.List(t.Context()); err == nil {
				t.Fatal("unsafe late middleware mutation was accepted")
			}
			if got := invoked.Load(); got != 1 {
				t.Errorf("adversarial middleware invoked %d times", got)
			}
			assertX509WorkloadNoRequests(t, issuer, api)
		})
	}
}

func TestX509WorkloadIdentityRejectsMethodLevelOverridesBeforeExchange(t *testing.T) {
	for _, test := range []struct {
		name string
		opt  option.RequestOption
	}{
		{name: "attacker endpoint", opt: option.WithBaseURL("https://attacker.example.test/v1/")},
		{name: "regional endpoint", opt: option.WithDataResidency(option.DataResidencyEU)},
		{name: "API key", opt: option.WithAPIKey("synthetic-method-api-key")},
		{name: "admin API key", opt: option.WithAdminAPIKey("synthetic-method-admin-key")},
		{name: "custom bearer", opt: option.WithHeader("Authorization", "Bearer synthetic-method-token")},
		{name: "HTTP client", opt: option.WithHTTPClient(&http.Client{})},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("OPENAI_BASE_URL", "https://mtls.api.openai.com/v1/")
			config, issuer, api := newX509WorkloadIdentityIntegration(t)
			client := openai.NewClient(option.WithX509WorkloadIdentity(config), option.WithMaxRetries(0))
			if _, err := client.Models.List(t.Context(), test.opt); err == nil {
				t.Fatal("unsafe method-level request option was accepted")
			}
			assertX509WorkloadNoRequests(t, issuer, api)
		})
	}
}

func TestX509WorkloadIdentityRejectsAzureAndBedrockBeforeExchange(t *testing.T) {
	t.Run("Azure endpoint without credentials", func(t *testing.T) {
		t.Setenv("OPENAI_BASE_URL", "https://mtls.api.openai.com/v1/")
		config, issuer, api := newX509WorkloadIdentityIntegration(t)
		client := openai.NewClient(
			azure.WithEndpoint("https://resource.openai.azure.com", "2024-06-01"),
			option.WithX509WorkloadIdentity(config),
			option.WithMaxRetries(0),
		)
		if _, err := client.Models.List(t.Context()); err == nil {
			t.Fatal("Azure endpoint accepted an X.509 workload identity")
		}
		assertX509WorkloadNoRequests(t, issuer, api)
	})
	t.Run("Azure API key", func(t *testing.T) {
		t.Setenv("OPENAI_BASE_URL", "https://mtls.api.openai.com/v1/")
		config, issuer, api := newX509WorkloadIdentityIntegration(t)
		client := openai.NewClient(
			azure.WithEndpoint("https://resource.openai.azure.com", "2024-06-01"),
			azure.WithAPIKey("synthetic-azure-key"),
			option.WithX509WorkloadIdentity(config),
			option.WithMaxRetries(0),
		)
		if _, err := client.Models.List(t.Context()); err == nil {
			t.Fatal("Azure API-key provider accepted an X.509 workload identity")
		}
		assertX509WorkloadNoRequests(t, issuer, api)
	})
	t.Run("Bedrock SkipAuth", func(t *testing.T) {
		t.Setenv("OPENAI_BASE_URL", "https://mtls.api.openai.com/v1/")
		config, issuer, api := newX509WorkloadIdentityIntegration(t)
		client, err := bedrock.NewClient(t.Context(), bedrock.Config{
			AWSRegion: "us-east-1",
			SkipAuth:  true,
		}, option.WithX509WorkloadIdentity(config), option.WithMaxRetries(0))
		if err == nil {
			_, err = client.Models.List(t.Context())
		}
		if err == nil {
			t.Fatal("Bedrock SkipAuth accepted an X.509 workload identity")
		}
		assertX509WorkloadNoRequests(t, issuer, api)
	})
}

type x509WorkloadRejectedDoer struct{}

func (x509WorkloadRejectedDoer) Do(*http.Request) (*http.Response, error) {
	return nil, errors.New("opaque synthetic transport unexpectedly invoked")
}

func newX509WorkloadIdentityIntegration(t *testing.T) (
	auth.X509WorkloadIdentity,
	*x509ConformanceServer,
	*x509ConformanceServer,
) {
	t.Helper()
	lab := newX509ConformanceLab(t)
	issuer := lab.server(t, x509ConformanceIssuerHost, true)
	api := lab.server(t, x509ConformanceAPIHost, false)
	template := lab.transport(t, x509ConformanceRoutes(issuer, api), lab.identity(t, "integrated-workload", true))
	transport, err := auth.NewX509Transport(template)
	if err != nil {
		t.Fatalf("attest integration workload transport: %v", err)
	}
	t.Cleanup(func() {
		if err := transport.Close(); err != nil {
			t.Errorf("close integration workload transport: %v", err)
		}
	})
	return auth.X509WorkloadIdentity{
		IdentityProviderID: "synthetic-identity-provider",
		ServiceAccountID:   "synthetic-service-account",
		Transport:          transport,
	}, issuer, api
}

func assertX509WorkloadNoRequests(t *testing.T, issuer, api *x509ConformanceServer) {
	t.Helper()
	if requests := issuer.requests(); len(requests) != 0 {
		t.Errorf("rejected configuration sent %d token-exchange requests", len(requests))
	}
	if requests := api.requests(); len(requests) != 0 {
		t.Errorf("rejected configuration sent %d protected API requests", len(requests))
	}
}

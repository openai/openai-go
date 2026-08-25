package openai_test

import (
	"errors"
	"io"
	"net/http"
	"sync/atomic"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/auth"
	"github.com/openai/openai-go/v3/azure"
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
		{name: "earlier native HTTP client override", before: func(auth.X509WorkloadIdentity) []option.RequestOption {
			return []option.RequestOption{option.WithHTTPClient(&http.Client{})}
		}},
		{name: "opaque HTTP client override", after: func(auth.X509WorkloadIdentity) []option.RequestOption {
			return []option.RequestOption{option.WithHTTPClient(x509WorkloadRejectedDoer{})}
		}},
		{name: "earlier opaque HTTP client override", before: func(auth.X509WorkloadIdentity) []option.RequestOption {
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
		{name: "nil header map", mutate: func(request *http.Request) { request.Header = nil }},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("OPENAI_BASE_URL", "https://mtls.api.openai.com/v1/")
			config, issuer, api := newX509WorkloadIdentityIntegration(t)
			var invoked atomic.Int32
			client := openai.NewClient(
				option.WithX509WorkloadIdentity(config),
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

func TestX509WorkloadIdentityPreservesStandaloneHTTPClientProvenance(t *testing.T) {
	for _, test := range []struct {
		name      string
		service   func(auth.X509WorkloadIdentity) openai.ModelService
		method    func(auth.X509WorkloadIdentity) []option.RequestOption
		wantCalls bool
	}{
		{
			name: "standalone implicit native default",
			service: func(config auth.X509WorkloadIdentity) openai.ModelService {
				return openai.NewModelService(option.WithX509WorkloadIdentity(config))
			},
			wantCalls: true,
		},
		{
			name: "standalone same-layer custom client",
			service: func(config auth.X509WorkloadIdentity) openai.ModelService {
				return openai.NewModelService(option.WithHTTPClient(&http.Client{}), option.WithX509WorkloadIdentity(config))
			},
		},
		{
			name: "standalone inherited custom client",
			service: func(auth.X509WorkloadIdentity) openai.ModelService {
				return openai.NewModelService(option.WithHTTPClient(&http.Client{}))
			},
			method: func(config auth.X509WorkloadIdentity) []option.RequestOption {
				return []option.RequestOption{option.WithX509WorkloadIdentity(config)}
			},
		},
		{
			name: "standalone later custom client",
			service: func(config auth.X509WorkloadIdentity) openai.ModelService {
				return openai.NewModelService(option.WithX509WorkloadIdentity(config), option.WithHTTPClient(&http.Client{}))
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			config, issuer, api := newX509WorkloadIdentityIntegration(t)
			service := test.service(config)
			var opts []option.RequestOption
			if test.method != nil {
				opts = test.method(config)
			}
			_, err := service.List(t.Context(), opts...)
			if test.wantCalls {
				if err != nil {
					t.Fatalf("standalone model service with its implicit native default: %v", err)
				}
				if len(issuer.requests()) != 1 || len(api.requests()) != 1 {
					t.Errorf("standalone model service issuer/API calls = %d/%d", len(issuer.requests()), len(api.requests()))
				}
				return
			}
			if err == nil {
				t.Fatal("standalone X.509 service accepted an explicit custom HTTP client")
			}
			assertX509WorkloadNoRequests(t, issuer, api)
		})
	}
}

func TestX509WorkloadIdentityUsesOnlyTheLastAuthenticationOption(t *testing.T) {
	t.Run("invalid earlier identity is never initialized", func(t *testing.T) {
		t.Setenv("OPENAI_BASE_URL", "https://mtls.api.openai.com/v1/")
		config, issuer, api := newX509WorkloadIdentityIntegration(t)
		invalid := config
		invalid.IdentityProviderID = ""
		client := openai.NewClient(option.WithX509WorkloadIdentity(invalid), option.WithX509WorkloadIdentity(config))
		if _, err := client.Models.List(t.Context()); err != nil {
			t.Fatalf("later valid X.509 option did not replace an invalid earlier identity: %v", err)
		}
		if len(issuer.requests()) != 1 || len(api.requests()) != 1 {
			t.Errorf("winning X.509 option issuer/API calls = %d/%d", len(issuer.requests()), len(api.requests()))
		}
	})
	t.Run("later client identity owns exchange and dispatch", func(t *testing.T) {
		t.Setenv("OPENAI_BASE_URL", "https://mtls.api.openai.com/v1/")
		first, firstIssuer, firstAPI := newX509WorkloadIdentityIntegration(t)
		second, secondIssuer, secondAPI := newX509WorkloadIdentityIntegration(t)
		client := openai.NewClient(option.WithX509WorkloadIdentity(first), option.WithX509WorkloadIdentity(second))
		if _, err := client.Models.List(t.Context()); err != nil {
			t.Fatalf("later X.509 transport did not replace an earlier generation: %v", err)
		}
		assertX509WorkloadNoRequests(t, firstIssuer, firstAPI)
		if len(secondIssuer.requests()) != 1 || len(secondAPI.requests()) != 1 {
			t.Errorf("winning X.509 transport issuer/API calls = %d/%d", len(secondIssuer.requests()), len(secondAPI.requests()))
		}
	})
	t.Run("method identity replaces inherited identity", func(t *testing.T) {
		t.Setenv("OPENAI_BASE_URL", "https://mtls.api.openai.com/v1/")
		first, firstIssuer, firstAPI := newX509WorkloadIdentityIntegration(t)
		second, secondIssuer, secondAPI := newX509WorkloadIdentityIntegration(t)
		client := openai.NewClient(option.WithX509WorkloadIdentity(first))
		if _, err := client.Models.List(t.Context(), option.WithX509WorkloadIdentity(second)); err != nil {
			t.Fatalf("method-level X.509 transport did not replace its inherited generation: %v", err)
		}
		assertX509WorkloadNoRequests(t, firstIssuer, firstAPI)
		if len(secondIssuer.requests()) != 1 || len(secondAPI.requests()) != 1 {
			t.Errorf("winning method X.509 transport issuer/API calls = %d/%d", len(secondIssuer.requests()), len(secondAPI.requests()))
		}
	})
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
		{name: "explicit headerless request", opt: option.WithHeaderDel("Authorization")},
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

func TestX509WorkloadIdentityPreservesHeaderlessRequestAfterEmptyCredentials(t *testing.T) {
	for _, test := range []struct {
		name string
		opt  option.RequestOption
	}{
		{name: "empty API key", opt: option.WithAPIKey("")},
		{name: "empty admin API key", opt: option.WithAdminAPIKey("")},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("OPENAI_BASE_URL", "https://mtls.api.openai.com/v1/")
			config, issuer, api := newX509WorkloadIdentityIntegration(t)
			client := openai.NewClient(option.WithX509WorkloadIdentity(config))
			if _, err := client.Models.List(t.Context(), option.WithHeaderDel("Authorization"), test.opt); err == nil {
				t.Fatal("empty credentials erased the caller's explicit headerless-request policy")
			}
			assertX509WorkloadNoRequests(t, issuer, api)
		})
	}
}

func TestX509WorkloadIdentityRejectsAdminOnlyOperationsBeforeExchange(t *testing.T) {
	t.Setenv("OPENAI_BASE_URL", "https://mtls.api.openai.com/v1/")
	config, issuer, api := newX509WorkloadIdentityIntegration(t)
	client := openai.NewClient(option.WithX509WorkloadIdentity(config))
	if _, err := client.Admin.Organization.DataRetention.Get(t.Context()); err == nil {
		t.Fatal("admin-only data-retention endpoint accepted X.509 bearer authentication")
	}
	assertX509WorkloadNoRequests(t, issuer, api)
}

func TestX509WorkloadIdentityDoesNotRetryPermanentIssuerFailure(t *testing.T) {
	t.Setenv("OPENAI_BASE_URL", "https://mtls.api.openai.com/v1/")
	config, issuer, api := newX509WorkloadIdentityIntegration(t)
	var requests atomic.Int32
	issuer.server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		requests.Add(1)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, `{"error":"invalid_grant","error_description":"synthetic-private-issuer-detail"}`)
	})
	client := openai.NewClient(option.WithX509WorkloadIdentity(config))
	_, err := client.Models.List(t.Context())
	var oauthError *auth.OAuthError
	if !errors.As(err, &oauthError) || oauthError.StatusCode != http.StatusBadRequest {
		t.Fatalf("permanent issuer error lost its typed OAuth identity: %v", err)
	}
	if got := requests.Load(); got != 1 {
		t.Errorf("default SDK retries repeated a permanent issuer failure %d times", got)
	}
	if requests := api.requests(); len(requests) != 0 {
		t.Errorf("permanent issuer failure sent %d protected API requests", len(requests))
	}
}

func TestX509WorkloadIdentityDoesNotMutateCallerOwnedRequest(t *testing.T) {
	t.Setenv("OPENAI_BASE_URL", "https://mtls.api.openai.com/v1/")
	config, issuer, api := newX509WorkloadIdentityIntegration(t)
	var observed atomic.Int32
	client := openai.NewClient(
		option.WithX509WorkloadIdentity(config),
		option.WithMiddleware(func(request *http.Request, next option.MiddlewareNext) (*http.Response, error) {
			response, err := next(request)
			if request.Header.Get("Authorization") != "" {
				t.Error("X.509 authentication wrote a bearer credential into the caller-owned request")
			}
			observed.Add(1)
			return response, err
		}),
	)
	if _, err := client.Models.List(t.Context()); err != nil {
		t.Fatalf("request with observer middleware: %v", err)
	}
	if observed.Load() != 1 || len(issuer.requests()) != 1 || len(api.requests()) != 1 {
		t.Errorf("request observer/issuer/API calls = %d/%d/%d", observed.Load(), len(issuer.requests()), len(api.requests()))
	}
	if request := api.requests()[0]; request.authorization != "Bearer "+x509ConformanceToken {
		t.Error("cloned API request did not receive its expected bearer credential")
	}
}

func TestX509WorkloadIdentityRejectsAzureBeforeExchange(t *testing.T) {
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

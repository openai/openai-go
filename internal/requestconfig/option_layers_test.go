package requestconfig

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestInheritedOptionsPreservePreRequestOptions(t *testing.T) {
	marker := errors.New("pre-request failure")
	seen := []string{}
	pre := func(value string) RequestOption {
		return PreRequestOptionFunc(func(cfg *RequestConfig) error { seen = append(seen, value); cfg.WebhookSecret = value; return nil })
	}
	ordinary := RequestOptionFunc(func(*RequestConfig) error { t.Fatal("ordinary option executed during pre-request scan"); return nil })
	inner := InheritedOptions(pre("inner"), ordinary, WithEnvironmentDefaultsDisabled())
	opts := InheritedOptions(append(inner, pre("outer"))...)
	cfg, err := PreRequestOptions(opts...)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.WebhookSecret != "outer" || !reflect.DeepEqual(seen, []string{"inner", "outer"}) {
		t.Fatalf("pre-request order lost: %#v / %q", seen, cfg.WebhookSecret)
	}
	if !EnvironmentDefaultsDisabled(opts...) {
		t.Fatal("nested provider marker lost")
	}
	_, err = PreRequestOptions(InheritedOptions(PreRequestOptionFunc(func(*RequestConfig) error { return marker }))...)
	if !errors.Is(err, marker) {
		t.Fatalf("pre-request error lost: %v", err)
	}
}

func TestInheritedOptionsCaptureSlice(t *testing.T) {
	first := RequestOptionFunc(func(cfg *RequestConfig) error { cfg.WebhookSecret = "first"; return nil })
	second := RequestOptionFunc(func(cfg *RequestConfig) error { cfg.WebhookSecret = "second"; return nil })
	original := []RequestOption{first}
	captured := InheritedOptions(original...)
	original[0] = second
	cfg := RequestConfig{}
	if err := cfg.Apply(captured...); err != nil {
		t.Fatal(err)
	}
	if cfg.WebhookSecret != "first" {
		t.Fatal("inherited slice was mutated")
	}
}

func TestInheritedOptionsRestoreOuterLayer(t *testing.T) {
	endpoint := &url.URL{Scheme: "https", Host: "example.com"}
	residency := RequestOptionFunc(func(cfg *RequestConfig) error { return cfg.SetEndpoint("data_residency", endpoint) })
	base := RequestOptionFunc(func(cfg *RequestConfig) error { return cfg.SetEndpoint("base_url", endpoint) })
	cfg := RequestConfig{}
	opts := append([]RequestOption{base}, InheritedOptions(residency)...)
	opts = append(opts, residency)
	if err := cfg.Apply(opts...); err == nil {
		t.Fatal("nested layer concealed an outer conflict")
	}
}

func TestEnvironmentDefaultsDisabledWrapperPreservesOptionBehavior(t *testing.T) {
	applied := false
	preApplied := false
	wrapped := WithEnvironmentDefaultsDisabled(
		RequestOptionFunc(func(*RequestConfig) error {
			applied = true
			return nil
		}),
		PreRequestOptionFunc(func(*RequestConfig) error {
			preApplied = true
			return nil
		}),
	)
	if !EnvironmentDefaultsDisabled(InheritedOptions(wrapped)...) {
		t.Fatal("wrapped environment marker was not detected")
	}
	cfg := RequestConfig{}
	if err := cfg.Apply(wrapped); err != nil {
		t.Fatal(err)
	}
	if !applied {
		t.Fatal("wrapped request option was not applied")
	}
	if _, err := PreRequestOptions(wrapped); err != nil {
		t.Fatal(err)
	}
	if !preApplied {
		t.Fatal("wrapped pre-request option was not applied")
	}
}

func TestEndpointProviderReportsProviderRouting(t *testing.T) {
	for _, provider := range []string{"", "Azure", "Bedrock"} {
		t.Run(provider, func(t *testing.T) {
			cfg := RequestConfig{}
			if provider != "" {
				if err := WithEndpointProvider(provider).Apply(&cfg); err != nil {
					t.Fatalf("configure provider routing: %v", err)
				}
			}
			if got := cfg.EndpointProvider(); got != provider {
				t.Errorf("endpoint provider = %q, want %q", got, provider)
			}
		})
	}
}

func TestAuthorizationHeaderOverriddenPreservesExplicitDeletion(t *testing.T) {
	request, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "https://example.com", nil)
	if err != nil {
		t.Fatalf("construct synthetic request: %v", err)
	}
	cfg := RequestConfig{Request: request}
	if cfg.AuthorizationHeaderOverridden() {
		t.Error("unmodified request unexpectedly has an explicit Authorization override")
	}
	cfg.DelHeader("Authorization")
	if !cfg.AuthorizationHeaderOverridden() || len(cfg.Request.Header.Values("Authorization")) != 0 {
		t.Error("explicit header deletion was not preserved as an Authorization override")
	}
	for _, test := range []struct {
		name string
		set  func(*RequestConfig)
	}{
		{name: "empty API key", set: func(config *RequestConfig) { config.SetAPIKey("") }},
		{name: "empty admin API key", set: func(config *RequestConfig) { config.SetAdminAPIKey("") }},
	} {
		t.Run(test.name, func(t *testing.T) {
			current := RequestConfig{Request: request.Clone(request.Context())}
			current.DelHeader("Authorization")
			test.set(&current)
			if current.authentication.headerOverride || current.authentication.authorizationHeaderLayer != nil {
				t.Error("negative control did not erase the ordinary override and nil-layer marker")
			}
			if !current.AuthorizationHeaderOverridden() {
				t.Error("empty credential selection erased explicit headerless-request provenance")
			}
			cloned := current.Clone(request.Context())
			if cloned == nil || !cloned.AuthorizationHeaderOverridden() {
				t.Fatal("request clone lost explicit headerless-request provenance")
			}
			cloned.ClearInheritedAuthentication()
			if cloned.AuthorizationHeaderOverridden() {
				t.Error("request clone did not clear inherited headerless-request provenance")
			}
		})
	}
	inherited := RequestConfig{Request: request.Clone(request.Context())}
	if err := inherited.Apply(InheritedOptions(RequestOptionFunc(func(config *RequestConfig) error {
		config.DelHeader("Authorization")
		return nil
	}))...); err != nil {
		t.Fatalf("configure inherited Authorization override: %v", err)
	}
	inherited.ClearInheritedAuthentication()
	if inherited.AuthorizationHeaderOverridden() {
		t.Error("inherited Authorization override was not cleared")
	}
}

func TestHTTPClientSelectionProvenanceDistinguishesDefaultsAndOverrides(t *testing.T) {
	trusted := RequestOptionFunc(func(cfg *RequestConfig) error {
		cfg.RecordHTTPClientSelection(true)
		return nil
	})
	explicit := RequestOptionFunc(func(cfg *RequestConfig) error {
		cfg.RecordHTTPClientSelection(false)
		return nil
	})
	t.Run("implicit native default", func(t *testing.T) {
		cfg := RequestConfig{}
		if cfg.HTTPClientExplicitlySelected() {
			t.Error("unselected implicit native HTTP client was classified as explicit")
		}
	})
	t.Run("marked generated default", func(t *testing.T) {
		cfg := RequestConfig{}
		if err := cfg.Apply(InheritedOptions(InheritedOptions(trusted)...)...); err != nil {
			t.Fatalf("apply marked generated HTTP client default: %v", err)
		}
		if cfg.HTTPClientExplicitlySelected() {
			t.Error("marked generated HTTP client default was classified as explicit")
		}
	})
	t.Run("standalone inherited custom client", func(t *testing.T) {
		cfg := RequestConfig{}
		if err := cfg.Apply(InheritedOptions(explicit)...); err != nil {
			t.Fatalf("apply standalone inherited custom HTTP client: %v", err)
		}
		if !cfg.HTTPClientExplicitlySelected() {
			t.Error("standalone inherited custom HTTP client lost its explicit provenance")
		}
	})
	t.Run("arbitrarily nested custom client", func(t *testing.T) {
		cfg := RequestConfig{}
		if err := cfg.Apply(InheritedOptions(InheritedOptions(InheritedOptions(explicit)...)...)...); err != nil {
			t.Fatalf("apply arbitrarily nested custom HTTP client: %v", err)
		}
		if !cfg.HTTPClientExplicitlySelected() {
			t.Error("nested custom HTTP client was misclassified as a trusted SDK default")
		}
	})
	t.Run("explicit override after marked default", func(t *testing.T) {
		cfg := RequestConfig{}
		if err := cfg.Apply(InheritedOptions(trusted)...); err != nil {
			t.Fatalf("apply marked HTTP client default: %v", err)
		}
		cfg.RecordHTTPClientSelection(false)
		if !cfg.HTTPClientExplicitlySelected() {
			t.Error("HTTP client override was not classified as explicit")
		}
	})
	t.Run("later default cannot erase explicit client", func(t *testing.T) {
		cfg := RequestConfig{}
		cfg.RecordHTTPClientSelection(false)
		cfg.RecordHTTPClientSelection(true)
		if !cfg.HTTPClientExplicitlySelected() {
			t.Error("later trusted default erased explicit HTTP client provenance")
		}
	})
	t.Run("clone preserves custom-client provenance", func(t *testing.T) {
		request, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "https://example.com", nil)
		if err != nil {
			t.Fatalf("construct clone provenance request: %v", err)
		}
		cfg := RequestConfig{Request: request}
		cfg.RecordHTTPClientSelection(false)
		cloned := cfg.Clone(t.Context())
		if cloned == nil || !cloned.HTTPClientExplicitlySelected() {
			t.Fatal("request clone lost explicit HTTP client selection provenance")
		}
	})
}

func TestProviderAuthOptionsPreserveConfigurationLayers(t *testing.T) {
	apiKey := NewProviderAuthOption("Azure", "azure.WithAPIKey")
	secondAPIKey := NewProviderAuthOption("Azure", "azure.WithAPIKey")
	token := NewProviderAuthOption("Azure", "azure.WithTokenCredential")

	t.Run("same mode uses the last option", func(t *testing.T) {
		cfg := RequestConfig{}
		if err := cfg.Apply(apiKey, secondAPIKey); err != nil {
			t.Fatal(err)
		}
		if !secondAPIKey.Selected(&cfg) {
			t.Fatal("last same-mode option was not selected")
		}
	})

	t.Run("different modes in one layer are rejected", func(t *testing.T) {
		cfg := RequestConfig{}
		err := cfg.Apply(apiKey, token)
		if err == nil {
			t.Fatal("expected ambiguous authentication error")
		}
	})

	t.Run("later layer replaces inherited mode", func(t *testing.T) {
		cfg := RequestConfig{}
		opts := append(InheritedOptions(apiKey), token)
		if err := cfg.Apply(opts...); err != nil {
			t.Fatal(err)
		}
		if !token.Selected(&cfg) {
			t.Fatal("request-layer authentication did not replace inherited mode")
		}
		if got, ok := cfg.ProviderAuth("Azure"); !ok || got != "azure.WithTokenCredential" {
			t.Fatalf("selected authentication = %q, %t", got, ok)
		}
	})

	t.Run("nested layer does not obscure the outer selection", func(t *testing.T) {
		cfg := RequestConfig{}
		opts := append([]RequestOption{apiKey}, InheritedOptions(token)...)
		opts = append(opts, secondAPIKey)
		if err := cfg.Apply(opts...); err != nil {
			t.Fatal(err)
		}
		if !secondAPIKey.Selected(&cfg) {
			t.Fatal("last outer-layer option was not selected")
		}
	})
}

func TestClearInheritedAuthenticationPreservesLayerConflicts(t *testing.T) {
	t.Run("inherited credentials are cleared", func(t *testing.T) {
		cfg := RequestConfig{}
		credentials := RequestOptionFunc(func(cfg *RequestConfig) error {
			cfg.SetAPIKey("openai-api-key")
			cfg.SetAdminAPIKey("openai-admin-key")
			return nil
		})
		if err := cfg.Apply(InheritedOptions(credentials)...); err != nil {
			t.Fatal(err)
		}
		cfg.ClearInheritedAuthentication()
		if cfg.APIKey != "" || cfg.AdminAPIKey != "" {
			t.Fatal("credentials were not cleared")
		}
	})

	for _, test := range []struct {
		name string
		set  func(*RequestConfig)
	}{
		{name: "API key", set: func(cfg *RequestConfig) { cfg.SetAPIKey("openai-api-key") }},
		{name: "admin API key", set: func(cfg *RequestConfig) { cfg.SetAdminAPIKey("openai-admin-key") }},
	} {
		t.Run(test.name+" in the current layer is preserved", func(t *testing.T) {
			cfg := RequestConfig{}
			test.set(&cfg)
			cfg.ClearInheritedAuthentication()
			if cfg.APIKey == "" && cfg.AdminAPIKey == "" {
				t.Fatal("same-layer credential was cleared")
			}
		})
	}
}

func TestClearInheritedAuthenticationHandlesAuthorizationHeaderLayers(t *testing.T) {
	newConfig := func(t *testing.T) RequestConfig {
		t.Helper()
		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://example.com", nil)
		if err != nil {
			t.Fatal(err)
		}
		return RequestConfig{Request: req}
	}

	t.Run("inherited header is cleared", func(t *testing.T) {
		cfg := newConfig(t)
		if err := cfg.Apply(InheritedOptions(RequestOptionFunc(func(cfg *RequestConfig) error {
			cfg.SetHeader("Authorization", "Bearer inherited-token")
			return nil
		}))...); err != nil {
			t.Fatal(err)
		}
		cfg.ClearInheritedAuthentication()
		if got := cfg.Request.Header.Values("Authorization"); len(got) != 0 {
			t.Fatalf("Authorization values = %q, want none", got)
		}
		if cfg.authentication.headerOverride {
			t.Fatal("inherited authorization override was not cleared")
		}
	})

	t.Run("current layer header is preserved", func(t *testing.T) {
		cfg := newConfig(t)
		cfg.SetHeader("Authorization", "Bearer current-token")
		cfg.ClearInheritedAuthentication()
		if got := cfg.Request.Header.Get("Authorization"); got != "Bearer current-token" {
			t.Fatalf("Authorization = %q", got)
		}
	})
}

func TestClearInheritedAuthenticationHandlesAPIKeyHeaderLayers(t *testing.T) {
	newConfig := func(t *testing.T) RequestConfig {
		t.Helper()
		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://example.com", nil)
		if err != nil {
			t.Fatal(err)
		}
		return RequestConfig{Request: req}
	}

	t.Run("inherited header is cleared", func(t *testing.T) {
		cfg := newConfig(t)
		if err := cfg.Apply(InheritedOptions(RequestOptionFunc(func(cfg *RequestConfig) error {
			cfg.SetHeader("Api-Key", "inherited-key")
			return nil
		}))...); err != nil {
			t.Fatal(err)
		}
		cfg.ClearInheritedAuthentication()
		if got := cfg.Request.Header.Values("Api-Key"); len(got) != 0 {
			t.Fatalf("Api-Key values = %q, want none", got)
		}
	})

	t.Run("current layer header is preserved", func(t *testing.T) {
		cfg := newConfig(t)
		cfg.SetHeader("Api-Key", "current-key")
		cfg.ClearInheritedAuthentication()
		if got := cfg.Request.Header.Get("Api-Key"); got != "current-key" {
			t.Fatalf("Api-Key = %q", got)
		}
	})
}

func TestProviderSelectionIsPreservedByClone(t *testing.T) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://example.com", nil)
	if err != nil {
		t.Fatal(err)
	}
	auth := NewProviderAuthOption("Azure", "azure.WithAPIKey")
	cfg := RequestConfig{Request: req}
	if err := cfg.Apply(WithProviderEndpointConfigured("Azure"), auth); err != nil {
		t.Fatal(err)
	}

	clone := cfg.Clone(context.Background())
	if clone == nil {
		t.Fatal("request config clone is nil")
	}
	if !clone.ProviderEndpointConfigured("Azure") {
		t.Fatal("provider endpoint was not preserved")
	}
	if got, ok := clone.ProviderAuth("Azure"); !ok || got != "azure.WithAPIKey" {
		t.Fatalf("cloned authentication = %q, %t", got, ok)
	}
}

func TestCloneTreatsOpenAICredentialsAsInherited(t *testing.T) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://example.com", nil)
	if err != nil {
		t.Fatal(err)
	}
	cfg := RequestConfig{Request: req}
	cfg.SetAPIKey("openai-api-key")
	cfg.SetAdminAPIKey("openai-admin-key")

	clone := cfg.Clone(context.Background())
	if clone == nil {
		t.Fatal("request config clone is nil")
	}
	clone.ClearInheritedAuthentication()
	if clone.APIKey != "" || clone.AdminAPIKey != "" {
		t.Fatal("cloned credentials were not cleared")
	}
}

func TestCloneTreatsAuthorizationHeaderAsInherited(t *testing.T) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://example.com", nil)
	if err != nil {
		t.Fatal(err)
	}
	cfg := RequestConfig{Request: req}
	cfg.SetHeader("Authorization", "Bearer custom-token")

	clone := cfg.Clone(context.Background())
	if clone == nil {
		t.Fatal("request config clone is nil")
	}
	clone.ClearInheritedAuthentication()
	if got := clone.Request.Header.Values("Authorization"); len(got) != 0 {
		t.Fatalf("Authorization values = %q, want none", got)
	}
}

func TestCloneTreatsAPIKeyHeaderAsInherited(t *testing.T) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://example.com", nil)
	if err != nil {
		t.Fatal(err)
	}
	cfg := RequestConfig{Request: req}
	cfg.SetHeader("Api-Key", "custom-key")

	clone := cfg.Clone(context.Background())
	if clone == nil {
		t.Fatal("request config clone is nil")
	}
	clone.ClearInheritedAuthentication()
	if got := clone.Request.Header.Values("Api-Key"); len(got) != 0 {
		t.Fatalf("Api-Key values = %q, want none", got)
	}
}

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

func TestClearInheritedOpenAICredentialsPreservesLayerConflicts(t *testing.T) {
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
		cfg.ClearInheritedOpenAICredentials()
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
			cfg.ClearInheritedOpenAICredentials()
			if cfg.APIKey == "" && cfg.AdminAPIKey == "" {
				t.Fatal("same-layer credential was cleared")
			}
		})
	}
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
	clone.ClearInheritedOpenAICredentials()
	if clone.APIKey != "" || clone.AdminAPIKey != "" {
		t.Fatal("cloned credentials were not cleared")
	}
}

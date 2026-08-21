package requestconfig

import (
	"fmt"
	"net/url"
	"slices"
)

// optionLayer preserves a configuration-call boundary after generated services
// concatenate inherited and request options. Its contents remain immutable.
type optionLayer []RequestOption

type optionLayerIdentity byte

func (opts optionLayer) Apply(cfg *RequestConfig) error {
	previousEndpoint := cfg.endpointSelector
	previousAuth := cfg.providerAuthLayer
	previousLayer := cfg.optionLayer
	cfg.endpointSelector = ""
	cfg.providerAuthLayer = nil
	cfg.optionLayer = new(optionLayerIdentity)
	defer func() {
		cfg.endpointSelector = previousEndpoint
		cfg.providerAuthLayer = previousAuth
		cfg.optionLayer = previousLayer
	}()
	return cfg.Apply(opts...)
}

// InheritedOptions captures one inherited configuration layer. This is internal
// API; generated client and service constructors use it before storing options.
func InheritedOptions(opts ...RequestOption) []RequestOption {
	if len(opts) == 0 {
		return nil
	}
	return []RequestOption{optionLayer(slices.Clone(opts))}
}

// endpointOption carries inspectable endpoint configuration without evaluating
// arbitrary request-option callbacks during provider construction.
type endpointOption struct {
	selector string
	endpoint *url.URL
	err      error
}

func (opt endpointOption) Apply(cfg *RequestConfig) error {
	if opt.err != nil {
		return opt.err
	}
	return cfg.SetEndpoint(opt.selector, opt.endpoint)
}

// WithEndpointSelection constructs an inspectable endpoint selector. This is
// internal API; public endpoint options own parsing and value validation.
func WithEndpointSelection(selector string, endpoint *url.URL, err error) RequestOption {
	return endpointOption{selector: selector, endpoint: endpoint, err: err}
}

// ValidateEndpointOptions checks built-in selectors before a provider performs
// credential discovery. It preserves inherited option layers and deliberately
// does not invoke user-defined RequestOption callbacks. Normal request setup
// still validates every option, including selectors applied by custom callbacks.
func ValidateEndpointOptions(provider string, opts ...RequestOption) error {
	cfg := RequestConfig{endpointProvider: provider}
	return cfg.Apply(endpointOptionsOnly(opts)...)
}

func endpointOptionsOnly(opts []RequestOption) []RequestOption {
	filtered := make([]RequestOption, 0, len(opts))
	for _, opt := range opts {
		switch opt := opt.(type) {
		case endpointOption:
			filtered = append(filtered, opt)
		case optionLayer:
			filtered = append(filtered, optionLayer(endpointOptionsOnly(opt)))
		}
	}
	return filtered
}

// SetEndpoint selects an explicit endpoint within the current option layer.
// Different selectors in one layer are an error regardless of their order.
func (cfg *RequestConfig) SetEndpoint(selector string, endpoint *url.URL) error {
	if selector == "data_residency" && cfg.endpointProvider != "" {
		return fmt.Errorf("requestoption: WithDataResidency cannot be used with the %s provider", cfg.endpointProvider)
	}
	if cfg.endpointSelector != "" && cfg.endpointSelector != selector {
		return fmt.Errorf("requestoption: WithDataResidency and WithBaseURL are mutually exclusive in the same configuration call")
	}
	cfg.endpointSelector = selector
	cfg.dataResidencyEndpoint = selector == "data_residency"
	cfg.BaseURL = endpoint
	return nil
}

// WithEndpointProvider marks a third-party provider whose routing and credentials
// cannot be combined with OpenAI data residency. This is internal API.
func WithEndpointProvider(provider string) RequestOption {
	return RequestOptionFunc(func(cfg *RequestConfig) error {
		if cfg.dataResidencyEndpoint {
			return fmt.Errorf("requestoption: WithDataResidency cannot be used with the %s provider", provider)
		}
		cfg.endpointProvider = provider
		return nil
	})
}

// WithProviderEndpointConfigured records that a provider-specific endpoint
// option was applied. Authentication options can still use WithEndpointProvider
// to retain provider conflict checks without claiming that routing was
// configured.
func WithProviderEndpointConfigured(provider string) RequestOption {
	return RequestOptionFunc(func(cfg *RequestConfig) error {
		if err := WithEndpointProvider(provider).Apply(cfg); err != nil {
			return err
		}
		cfg.configuredProviderEndpoint = provider
		return nil
	})
}

// ProviderEndpointConfigured reports whether the provider-specific endpoint
// option was applied. It does not validate the effective request origin after
// other routing options or redirects.
func (cfg *RequestConfig) ProviderEndpointConfigured(provider string) bool {
	return cfg.configuredProviderEndpoint == provider
}

// ProviderAuthOption identifies one provider authentication mode. Instances are
// immutable and can be shared across concurrent requests.
type ProviderAuthOption struct {
	provider string
	selector string
}

// NewProviderAuthOption constructs an inspectable provider authentication
// option. Different authentication selectors in one option layer are rejected;
// a later inherited or request layer may explicitly replace an earlier mode.
func NewProviderAuthOption(provider string, selector string) *ProviderAuthOption {
	return &ProviderAuthOption{provider: provider, selector: selector}
}

func (opt *ProviderAuthOption) Apply(cfg *RequestConfig) error {
	if previous := cfg.providerAuthLayer; previous != nil &&
		(previous.provider != opt.provider || previous.selector != opt.selector) {
		return fmt.Errorf(
			"requestconfig: %s authentication is ambiguous; %s and %s cannot be combined in the same configuration call",
			opt.provider,
			previous.selector,
			opt.selector,
		)
	}
	cfg.providerAuthLayer = opt
	cfg.providerAuth = opt
	return nil
}

// Selected reports whether this authentication option won after every option
// layer was applied to cfg.
func (opt *ProviderAuthOption) Selected(cfg *RequestConfig) bool {
	return cfg.providerAuth == opt
}

// ProviderAuth reports the selected authentication mode for provider.
func (cfg *RequestConfig) ProviderAuth(provider string) (string, bool) {
	if cfg.providerAuth == nil || cfg.providerAuth.provider != provider {
		return "", false
	}
	return cfg.providerAuth.selector, true
}

// ClearInheritedAuthentication removes OpenAI credentials and a custom
// Authorization header selected by an earlier option layer. Authentication
// selected in the current layer remains so a provider finalizer can reject the
// ambiguity after all options are applied.
func (cfg *RequestConfig) ClearInheritedAuthentication() {
	if cfg.APIKey != "" && cfg.apiKeyLayer != cfg.optionLayer {
		cfg.APIKey = ""
		cfg.apiKeyLayer = nil
	}
	if cfg.AdminAPIKey != "" && cfg.adminAPIKeyLayer != cfg.optionLayer {
		cfg.AdminAPIKey = ""
		cfg.adminAPIKeyLayer = nil
	}
	if cfg.APIKey == "" && cfg.AdminAPIKey == "" {
		cfg.authPreference = authCredentialPreferenceNone
	}
	if cfg.authorizationHeaderLayer != cfg.optionLayer {
		cfg.Request.Header.Del("Authorization")
		cfg.authorizationHeaderLayer = nil
		cfg.authHeaderOverride = false
	}
}

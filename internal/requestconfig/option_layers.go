package requestconfig

import (
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"
)

// optionLayer preserves a configuration-call boundary after generated services
// concatenate inherited and request options. Its contents remain immutable.
type optionLayer []RequestOption

type optionLayerIdentity byte

type authenticationState struct {
	currentLayer             *optionLayerIdentity
	providerInCurrentLayer   *ProviderAuthOption
	selectedProvider         *ProviderAuthOption
	retryScopeFactory        *requestRetryScopeFactory
	httpClientExplicit       bool
	apiKeyLayer              *optionLayerIdentity
	adminAPIKeyLayer         *optionLayerIdentity
	authorizationHeaderLayer *optionLayerIdentity
	apiKeyHeaderLayer        *optionLayerIdentity
	headerOverride           bool
	authorizationExplicit    bool
	preference               authCredentialPreference
}

func (state *authenticationState) enterLayer() func() {
	previousLayer := state.currentLayer
	previousProvider := state.providerInCurrentLayer
	state.currentLayer = new(optionLayerIdentity)
	state.providerInCurrentLayer = nil
	return func() {
		state.currentLayer = previousLayer
		state.providerInCurrentLayer = previousProvider
	}
}

func (opts optionLayer) Apply(cfg *RequestConfig) error {
	previousEndpoint := cfg.endpointSelector
	cfg.endpointSelector = ""
	restoreAuthenticationLayer := cfg.authentication.enterLayer()
	defer func() {
		cfg.endpointSelector = previousEndpoint
		restoreAuthenticationLayer()
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

// EndpointProvider reports the third-party provider associated with this
// request's routing, including providers that intentionally skip authentication.
func (cfg *RequestConfig) EndpointProvider() string {
	return cfg.endpointProvider
}

// AuthorizationHeaderOverridden reports whether an explicit request option
// selected or deleted the Authorization header.
func (cfg *RequestConfig) AuthorizationHeaderOverridden() bool {
	return cfg.authentication.headerOverride || cfg.authentication.authorizationExplicit
}

// DefaultHTTPClient marks a native client created by the SDK's own defaults.
// Applications outside the SDK cannot construct this internal marker.
type DefaultHTTPClient struct {
	*http.Client
}

// RecordHTTPClientSelection records whether an HTTP client was explicitly
// supplied. A trusted SDK default never clears an earlier explicit selection.
func (cfg *RequestConfig) RecordHTTPClientSelection(trustedDefault bool) {
	if !trustedDefault {
		cfg.authentication.httpClientExplicit = true
	}
}

// HTTPClientExplicitlySelected reports whether any caller supplied a native or
// custom client, regardless of option-layer nesting or later SDK defaults.
func (cfg *RequestConfig) HTTPClientExplicitlySelected() bool {
	return cfg.authentication.httpClientExplicit
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
	if previous := cfg.authentication.providerInCurrentLayer; previous != nil &&
		(previous.provider != opt.provider || previous.selector != opt.selector) {
		return fmt.Errorf(
			"requestconfig: %s authentication is ambiguous; %s and %s cannot be combined in the same configuration call",
			opt.provider,
			previous.selector,
			opt.selector,
		)
	}
	cfg.authentication.providerInCurrentLayer = opt
	cfg.authentication.selectedProvider = opt
	return nil
}

// Selected reports whether this authentication option won after every option
// layer was applied to cfg.
func (opt *ProviderAuthOption) Selected(cfg *RequestConfig) bool {
	return cfg.authentication.selectedProvider == opt
}

// ProviderAuth reports the selected authentication mode for provider.
func (cfg *RequestConfig) ProviderAuth(provider string) (string, bool) {
	selected := cfg.authentication.selectedProvider
	if selected == nil || selected.provider != provider {
		return "", false
	}
	return selected.selector, true
}

func (state *authenticationState) recordHeader(name string) {
	switch {
	case strings.EqualFold(name, "Authorization"):
		state.headerOverride = true
		state.authorizationExplicit = true
		state.authorizationHeaderLayer = state.currentLayer
	case strings.EqualFold(name, "Api-Key"):
		state.apiKeyHeaderLayer = state.currentLayer
	}
}

func (state *authenticationState) recordAPIKey() {
	state.apiKeyLayer = state.currentLayer
	state.preference = authCredentialPreferenceBearer
	state.headerOverride = false
}

func (state *authenticationState) recordAdminAPIKey() {
	state.adminAPIKeyLayer = state.currentLayer
	state.preference = authCredentialPreferenceAdmin
	state.headerOverride = false
}

func (state authenticationState) cloneAsInherited(cfg *RequestConfig) authenticationState {
	state.currentLayer = nil
	state.providerInCurrentLayer = nil
	state.apiKeyLayer = nil
	state.adminAPIKeyLayer = nil
	state.authorizationHeaderLayer = nil
	state.apiKeyHeaderLayer = nil
	if state.retryScopeFactory != nil {
		state.retryScopeFactory.install(cfg)
	}

	hasAuthorizationOverride := state.headerOverride || state.authorizationExplicit ||
		len(cfg.Request.Header.Values("Authorization")) != 0
	hasAPIKeyHeader := len(cfg.Request.Header.Values("Api-Key")) != 0
	if cfg.APIKey == "" && cfg.AdminAPIKey == "" && !hasAuthorizationOverride && !hasAPIKeyHeader {
		return state
	}

	inheritedLayer := new(optionLayerIdentity)
	if cfg.APIKey != "" {
		state.apiKeyLayer = inheritedLayer
	}
	if cfg.AdminAPIKey != "" {
		state.adminAPIKeyLayer = inheritedLayer
	}
	if hasAuthorizationOverride {
		state.authorizationHeaderLayer = inheritedLayer
	}
	if hasAPIKeyHeader {
		state.apiKeyHeaderLayer = inheritedLayer
	}
	return state
}

// ClearInheritedAuthentication removes OpenAI credentials and a custom
// Authorization header selected by an earlier option layer. Authentication
// selected in the current layer remains so a provider finalizer can reject the
// ambiguity after all options are applied.
func (cfg *RequestConfig) ClearInheritedAuthentication() {
	state := &cfg.authentication
	if cfg.APIKey != "" && state.apiKeyLayer != state.currentLayer {
		cfg.APIKey = ""
		state.apiKeyLayer = nil
	}
	if cfg.AdminAPIKey != "" && state.adminAPIKeyLayer != state.currentLayer {
		cfg.AdminAPIKey = ""
		state.adminAPIKeyLayer = nil
	}
	if cfg.APIKey == "" && cfg.AdminAPIKey == "" {
		state.preference = authCredentialPreferenceNone
	}
	if state.authorizationHeaderLayer != state.currentLayer {
		cfg.Request.Header.Del("Authorization")
		state.authorizationHeaderLayer = nil
		state.headerOverride = false
		state.authorizationExplicit = false
	}
	if state.apiKeyHeaderLayer != state.currentLayer {
		cfg.Request.Header.Del("Api-Key")
		state.apiKeyHeaderLayer = nil
	}
}

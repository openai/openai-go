package requestconfig

import (
	"fmt"
	"net/url"
	"slices"
)

// optionLayer preserves a configuration-call boundary after generated services
// concatenate inherited and request options. Its contents remain immutable.
type optionLayer []RequestOption

func (opts optionLayer) Apply(cfg *RequestConfig) error {
	previous := cfg.endpointSelector
	cfg.endpointSelector = ""
	defer func() { cfg.endpointSelector = previous }()
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

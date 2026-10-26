// File generated from our OpenAPI spec by Castiron. See CONTRIBUTING.md for details.

package option

import (
	"fmt"
	"github.com/openai/openai-go/v3/internal/requestconfig"
	"net/url"
)

// DataResidency selects an OpenAI API endpoint. Endpoint selection does not
// grant access to a region or change project and model eligibility.
type DataResidency string

const (
	DataResidencyGlobal DataResidency = "global"
	DataResidencyUS     DataResidency = "us"
	DataResidencyEU     DataResidency = "eu"
	DataResidencyAE     DataResidency = "ae"
)

// WithDataResidency selects the OpenAI endpoint for a client, service, or request.
// It cannot be combined with WithBaseURL in the same configuration call, but may
// replace an inherited endpoint. It sends no additional request fields or headers
// and never falls back to a different region. An empty or unknown value is invalid.
func WithDataResidency(residency DataResidency) RequestOption {
	var base string
	switch residency {
	case DataResidencyGlobal:
		base = "https://api.openai.com/v1/"
	case DataResidencyUS:
		base = "https://us.api.openai.com/v1/"
	case DataResidencyEU:
		base = "https://eu.api.openai.com/v1/"
	case DataResidencyAE:
		base = "https://ae.api.openai.com/v1/"
	default:
		return requestconfig.WithEndpointSelection("data_residency", nil,
			fmt.Errorf("requestoption: invalid data residency %q (expected global, us, eu, or ae)", residency))
	}
	endpoint, _ := url.Parse(base)
	return requestconfig.WithEndpointSelection("data_residency", endpoint, nil)
}

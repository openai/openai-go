// File generated from our OpenAPI spec by Castiron. See CONTRIBUTING.md for details.

package live

import (
	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
)

// SidebandService contains methods and other services that help with interacting
// with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewSidebandService] method instead.
type SidebandService struct {
	Options []option.RequestOption
}

// NewSidebandService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewSidebandService(opts ...option.RequestOption) (r SidebandService) {
	r = SidebandService{}
	r.Options = requestconfig.InheritedOptions(opts...)
	return
}

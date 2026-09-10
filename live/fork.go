// File generated from our OpenAPI spec by Castiron. See CONTRIBUTING.md for details.

package live

import (
	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
)

// ForkService contains methods and other services that help with interacting with
// the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewForkService] method instead.
type ForkService struct {
	Options []option.RequestOption
}

// NewForkService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewForkService(opts ...option.RequestOption) (r ForkService) {
	r = ForkService{}
	r.Options = requestconfig.InheritedOptions(opts...)
	return
}

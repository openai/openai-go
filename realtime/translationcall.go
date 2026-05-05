// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package realtime

import (
	"github.com/openai/openai-go/v3/option"
)

// TranslationCallService contains methods and other services that help with
// interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewTranslationCallService] method instead.
type TranslationCallService struct {
	Options []option.RequestOption
}

// NewTranslationCallService generates a new service that applies the given options
// to each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewTranslationCallService(opts ...option.RequestOption) (r TranslationCallService) {
	r = TranslationCallService{}
	r.Options = opts
	return
}

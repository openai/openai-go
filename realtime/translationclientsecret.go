// File generated from our OpenAPI spec by Castiron. See CONTRIBUTING.md for details.

package realtime

import (
	"context"
	"net/http"
	"slices"

	"github.com/openai/openai-go/v3/internal/apijson"
	shimjson "github.com/openai/openai-go/v3/internal/encoding/json"
	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
)

// TranslationClientSecretService contains methods and other services that help
// with interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewTranslationClientSecretService] method instead.
type TranslationClientSecretService struct {
	Options []option.RequestOption
}

// NewTranslationClientSecretService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewTranslationClientSecretService(opts ...option.RequestOption) (r TranslationClientSecretService) {
	r = TranslationClientSecretService{}
	r.Options = requestconfig.InheritedOptions(opts...)
	return
}

// Create a Realtime translation client secret with an associated translation
// session configuration.
//
// Client secrets are short-lived tokens that can be passed to a client app, such
// as a web frontend or mobile client, which grants access to the Realtime
// Translation API without leaking your main API key. You can configure a custom
// TTL for each client secret.
//
// Returns the created client secret and the effective translation session object.
// The client secret is a string that looks like `ek_1234`.
func (r *TranslationClientSecretService) New(ctx context.Context, body TranslationClientSecretNewParams, opts ...option.RequestOption) (res *RealtimeTranslationClientSecretCreateResponse, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	path := "realtime/translations/client_secrets"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

type TranslationClientSecretNewParams struct {
	// Create a translation session and client secret for the Realtime API.
	RealtimeTranslationClientSecretCreateRequest RealtimeTranslationClientSecretCreateRequestParam
	paramObj
}

func (r TranslationClientSecretNewParams) MarshalJSON() (data []byte, err error) {
	return shimjson.Marshal(r.RealtimeTranslationClientSecretCreateRequest)
}
func (r *TranslationClientSecretNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// File generated from our OpenAPI spec by Castiron. See CONTRIBUTING.md for details.

package openai

import (
	"context"
	"errors"
	"net/http"
	"slices"

	"github.com/openai/openai-go/v3/internal/apijson"
	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/respjson"
	"github.com/openai/openai-go/v3/shared/constant"
)

// SafetyCaseService contains methods and other services that help with interacting
// with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewSafetyCaseService] method instead.
type SafetyCaseService struct {
	Options []option.RequestOption
}

// NewSafetyCaseService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewSafetyCaseService(opts ...option.RequestOption) (r SafetyCaseService) {
	r = SafetyCaseService{}
	r.Options = requestconfig.InheritedOptions(opts...)
	return
}

// Get a safety case by ID.
func (r *SafetyCaseService) Get(ctx context.Context, id string, opts ...option.RequestOption) (res *SafetyCase, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("safety/cases/%s", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

type SafetyCase struct {
	ID               string              `json:"id" api:"required"`
	CreatedAt        int64               `json:"created_at" api:"required" format:"unixtime"`
	EntityIdentifier string              `json:"entity_identifier" api:"required"`
	Notice           SafetyCaseNotice    `json:"notice" api:"required"`
	Object           constant.SafetyCase `json:"object" default:"safety.case"`
	Reason           string              `json:"reason" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID               respjson.Field
		CreatedAt        respjson.Field
		EntityIdentifier respjson.Field
		Notice           respjson.Field
		Object           respjson.Field
		Reason           respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r SafetyCase) RawJSON() string { return r.JSON.raw }
func (r *SafetyCase) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type SafetyCaseNotice struct {
	// Any of "warning", "deactivation".
	Type string `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r SafetyCaseNotice) RawJSON() string { return r.JSON.raw }
func (r *SafetyCaseNotice) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai

import (
	"context"
	"errors"
	"net/http"
	"slices"

	"github.com/openai/openai-go/v3/internal/apijson"
	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/packages/respjson"
	"github.com/openai/openai-go/v3/shared/constant"
)

// AdminOrganizationProjectServiceAccountAPIKeyService contains methods and other
// services that help with interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAdminOrganizationProjectServiceAccountAPIKeyService] method instead.
type AdminOrganizationProjectServiceAccountAPIKeyService struct {
	Options []option.RequestOption
}

// NewAdminOrganizationProjectServiceAccountAPIKeyService generates a new service
// that applies the given options to each request. These options are applied after
// the parent client's options (if there is one), and before any request-specific
// options.
func NewAdminOrganizationProjectServiceAccountAPIKeyService(opts ...option.RequestOption) (r AdminOrganizationProjectServiceAccountAPIKeyService) {
	r = AdminOrganizationProjectServiceAccountAPIKeyService{}
	r.Options = opts
	return
}

// Creates an API key for a service account in the project.
func (r *AdminOrganizationProjectServiceAccountAPIKeyService) New(ctx context.Context, projectID string, serviceAccountID string, body AdminOrganizationProjectServiceAccountAPIKeyNewParams, opts ...option.RequestOption) (res *AdminOrganizationProjectServiceAccountAPIKeyNewResponse, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if projectID == "" {
		err = errors.New("missing required project_id parameter")
		return nil, err
	}
	if serviceAccountID == "" {
		err = errors.New("missing required service_account_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("organization/projects/%s/service_accounts/%s/api_keys", projectID, serviceAccountID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

type AdminOrganizationProjectServiceAccountAPIKeyNewResponse struct {
	ID        string `json:"id" api:"required"`
	CreatedAt int64  `json:"created_at" api:"required" format:"unixtime"`
	Name      string `json:"name" api:"required"`
	// The object type, which is always `organization.project.service_account.api_key`
	Object constant.OrganizationProjectServiceAccountAPIKey `json:"object" default:"organization.project.service_account.api_key"`
	Value  string                                           `json:"value" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		Name        respjson.Field
		Object      respjson.Field
		Value       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationProjectServiceAccountAPIKeyNewResponse) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationProjectServiceAccountAPIKeyNewResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationProjectServiceAccountAPIKeyNewParams struct {
	// API key name.
	Name param.Opt[string] `json:"name,omitzero"`
	// API key scopes.
	Scopes []string `json:"scopes,omitzero"`
	paramObj
}

func (r AdminOrganizationProjectServiceAccountAPIKeyNewParams) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationProjectServiceAccountAPIKeyNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationProjectServiceAccountAPIKeyNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

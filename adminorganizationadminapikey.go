// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"

	"github.com/openai/openai-go/v3/internal/apijson"
	"github.com/openai/openai-go/v3/internal/apiquery"
	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/pagination"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/packages/respjson"
)

// AdminOrganizationAdminAPIKeyService contains methods and other services that
// help with interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAdminOrganizationAdminAPIKeyService] method instead.
type AdminOrganizationAdminAPIKeyService struct {
	Options []option.RequestOption
}

// NewAdminOrganizationAdminAPIKeyService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewAdminOrganizationAdminAPIKeyService(opts ...option.RequestOption) (r AdminOrganizationAdminAPIKeyService) {
	r = AdminOrganizationAdminAPIKeyService{}
	r.Options = opts
	return
}

// Create an organization admin API key
func (r *AdminOrganizationAdminAPIKeyService) New(ctx context.Context, body AdminOrganizationAdminAPIKeyNewParams, opts ...option.RequestOption) (res *AdminAPIKey, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	path := "organization/admin_api_keys"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Retrieve a single organization API key
func (r *AdminOrganizationAdminAPIKeyService) Get(ctx context.Context, keyID string, opts ...option.RequestOption) (res *AdminAPIKey, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if keyID == "" {
		err = errors.New("missing required key_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("organization/admin_api_keys/%s", keyID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// List organization API keys
func (r *AdminOrganizationAdminAPIKeyService) List(ctx context.Context, query AdminOrganizationAdminAPIKeyListParams, opts ...option.RequestOption) (res *pagination.CursorPage[AdminAPIKey], err error) {
	var raw *http.Response
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "organization/admin_api_keys"
	cfg, err := requestconfig.NewRequestConfig(ctx, http.MethodGet, path, query, &res, opts...)
	if err != nil {
		return nil, err
	}
	err = cfg.Execute()
	if err != nil {
		return nil, err
	}
	res.SetPageConfig(cfg, raw)
	return res, nil
}

// List organization API keys
func (r *AdminOrganizationAdminAPIKeyService) ListAutoPaging(ctx context.Context, query AdminOrganizationAdminAPIKeyListParams, opts ...option.RequestOption) *pagination.CursorPageAutoPager[AdminAPIKey] {
	return pagination.NewCursorPageAutoPager(r.List(ctx, query, opts...))
}

// Delete an organization admin API key
func (r *AdminOrganizationAdminAPIKeyService) Delete(ctx context.Context, keyID string, opts ...option.RequestOption) (res *AdminOrganizationAdminAPIKeyDeleteResponse, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if keyID == "" {
		err = errors.New("missing required key_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("organization/admin_api_keys/%s", keyID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Represents an individual Admin API key in an org.
type AdminAPIKey struct {
	// The identifier, which can be referenced in API endpoints
	ID string `json:"id" api:"required"`
	// The Unix timestamp (in seconds) of when the API key was created
	CreatedAt int64 `json:"created_at" api:"required" format:"unixtime"`
	// The Unix timestamp (in seconds) of when the API key was last used
	LastUsedAt int64 `json:"last_used_at" api:"required" format:"unixtime"`
	// The name of the API key
	Name string `json:"name" api:"required"`
	// The object type, which is always `organization.admin_api_key`
	Object string           `json:"object" api:"required"`
	Owner  AdminAPIKeyOwner `json:"owner" api:"required"`
	// The redacted value of the API key
	RedactedValue string `json:"redacted_value" api:"required"`
	// The value of the API key. Only shown on create.
	Value string `json:"value"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID            respjson.Field
		CreatedAt     respjson.Field
		LastUsedAt    respjson.Field
		Name          respjson.Field
		Object        respjson.Field
		Owner         respjson.Field
		RedactedValue respjson.Field
		Value         respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminAPIKey) RawJSON() string { return r.JSON.raw }
func (r *AdminAPIKey) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminAPIKeyOwner struct {
	// The identifier, which can be referenced in API endpoints
	ID string `json:"id"`
	// The Unix timestamp (in seconds) of when the user was created
	CreatedAt int64 `json:"created_at" format:"unixtime"`
	// The name of the user
	Name string `json:"name"`
	// The object type, which is always organization.user
	Object string `json:"object"`
	// Always `owner`
	Role string `json:"role"`
	// Always `user`
	Type string `json:"type"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		Name        respjson.Field
		Object      respjson.Field
		Role        respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminAPIKeyOwner) RawJSON() string { return r.JSON.raw }
func (r *AdminAPIKeyOwner) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationAdminAPIKeyDeleteResponse struct {
	ID      string `json:"id"`
	Deleted bool   `json:"deleted"`
	Object  string `json:"object"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Deleted     respjson.Field
		Object      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAdminAPIKeyDeleteResponse) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationAdminAPIKeyDeleteResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationAdminAPIKeyNewParams struct {
	Name string `json:"name" api:"required"`
	paramObj
}

func (r AdminOrganizationAdminAPIKeyNewParams) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationAdminAPIKeyNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationAdminAPIKeyNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationAdminAPIKeyListParams struct {
	// Return keys with IDs that come after this ID in the pagination order.
	After param.Opt[string] `query:"after,omitzero" json:"-"`
	// Maximum number of keys to return.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Order results by creation time, ascending or descending.
	//
	// Any of "asc", "desc".
	Order AdminOrganizationAdminAPIKeyListParamsOrder `query:"order,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [AdminOrganizationAdminAPIKeyListParams]'s query parameters
// as `url.Values`.
func (r AdminOrganizationAdminAPIKeyListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Order results by creation time, ascending or descending.
type AdminOrganizationAdminAPIKeyListParamsOrder string

const (
	AdminOrganizationAdminAPIKeyListParamsOrderAsc  AdminOrganizationAdminAPIKeyListParamsOrder = "asc"
	AdminOrganizationAdminAPIKeyListParamsOrderDesc AdminOrganizationAdminAPIKeyListParamsOrder = "desc"
)

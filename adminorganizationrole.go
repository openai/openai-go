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
	"github.com/openai/openai-go/v3/shared/constant"
)

// AdminOrganizationRoleService contains methods and other services that help with
// interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAdminOrganizationRoleService] method instead.
type AdminOrganizationRoleService struct {
	Options []option.RequestOption
}

// NewAdminOrganizationRoleService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewAdminOrganizationRoleService(opts ...option.RequestOption) (r AdminOrganizationRoleService) {
	r = AdminOrganizationRoleService{}
	r.Options = opts
	return
}

// Creates a custom role for the organization.
func (r *AdminOrganizationRoleService) New(ctx context.Context, body AdminOrganizationRoleNewParams, opts ...option.RequestOption) (res *Role, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	path := "organization/roles"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Updates an existing organization role.
func (r *AdminOrganizationRoleService) Update(ctx context.Context, roleID string, body AdminOrganizationRoleUpdateParams, opts ...option.RequestOption) (res *Role, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if roleID == "" {
		err = errors.New("missing required role_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("organization/roles/%s", roleID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Lists the roles configured for the organization.
func (r *AdminOrganizationRoleService) List(ctx context.Context, query AdminOrganizationRoleListParams, opts ...option.RequestOption) (res *pagination.CursorPage[Role], err error) {
	var raw *http.Response
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "organization/roles"
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

// Lists the roles configured for the organization.
func (r *AdminOrganizationRoleService) ListAutoPaging(ctx context.Context, query AdminOrganizationRoleListParams, opts ...option.RequestOption) *pagination.CursorPageAutoPager[Role] {
	return pagination.NewCursorPageAutoPager(r.List(ctx, query, opts...))
}

// Deletes a custom role from the organization.
func (r *AdminOrganizationRoleService) Delete(ctx context.Context, roleID string, opts ...option.RequestOption) (res *AdminOrganizationRoleDeleteResponse, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if roleID == "" {
		err = errors.New("missing required role_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("organization/roles/%s", roleID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Details about a role that can be assigned through the public Roles API.
type Role struct {
	// Identifier for the role.
	ID string `json:"id" api:"required"`
	// Optional description of the role.
	Description string `json:"description" api:"required"`
	// Unique name for the role.
	Name string `json:"name" api:"required"`
	// Always `role`.
	Object constant.Role `json:"object" default:"role"`
	// Permissions granted by the role.
	Permissions []string `json:"permissions" api:"required"`
	// Whether the role is predefined and managed by OpenAI.
	PredefinedRole bool `json:"predefined_role" api:"required"`
	// Resource type the role is bound to (for example `api.organization` or
	// `api.project`).
	ResourceType string `json:"resource_type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID             respjson.Field
		Description    respjson.Field
		Name           respjson.Field
		Object         respjson.Field
		Permissions    respjson.Field
		PredefinedRole respjson.Field
		ResourceType   respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r Role) RawJSON() string { return r.JSON.raw }
func (r *Role) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Confirmation payload returned after deleting a role.
type AdminOrganizationRoleDeleteResponse struct {
	// Identifier of the deleted role.
	ID string `json:"id" api:"required"`
	// Whether the role was deleted.
	Deleted bool `json:"deleted" api:"required"`
	// Always `role.deleted`.
	Object constant.RoleDeleted `json:"object" default:"role.deleted"`
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
func (r AdminOrganizationRoleDeleteResponse) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationRoleDeleteResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationRoleNewParams struct {
	// Permissions to grant to the role.
	Permissions []string `json:"permissions,omitzero" api:"required"`
	// Unique name for the role.
	RoleName string `json:"role_name" api:"required"`
	// Optional description of the role.
	Description param.Opt[string] `json:"description,omitzero"`
	paramObj
}

func (r AdminOrganizationRoleNewParams) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationRoleNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationRoleNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationRoleUpdateParams struct {
	// New description for the role.
	Description param.Opt[string] `json:"description,omitzero"`
	// New name for the role.
	RoleName param.Opt[string] `json:"role_name,omitzero"`
	// Updated set of permissions for the role.
	Permissions []string `json:"permissions,omitzero"`
	paramObj
}

func (r AdminOrganizationRoleUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationRoleUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationRoleUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationRoleListParams struct {
	// Cursor for pagination. Provide the value from the previous response's `next`
	// field to continue listing roles.
	After param.Opt[string] `query:"after,omitzero" json:"-"`
	// A limit on the number of roles to return. Defaults to 1000.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Sort order for the returned roles.
	//
	// Any of "asc", "desc".
	Order AdminOrganizationRoleListParamsOrder `query:"order,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [AdminOrganizationRoleListParams]'s query parameters as
// `url.Values`.
func (r AdminOrganizationRoleListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Sort order for the returned roles.
type AdminOrganizationRoleListParamsOrder string

const (
	AdminOrganizationRoleListParamsOrderAsc  AdminOrganizationRoleListParamsOrder = "asc"
	AdminOrganizationRoleListParamsOrderDesc AdminOrganizationRoleListParamsOrder = "desc"
)

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

// AdminOrganizationProjectRoleService contains methods and other services that
// help with interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAdminOrganizationProjectRoleService] method instead.
type AdminOrganizationProjectRoleService struct {
	Options []option.RequestOption
}

// NewAdminOrganizationProjectRoleService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewAdminOrganizationProjectRoleService(opts ...option.RequestOption) (r AdminOrganizationProjectRoleService) {
	r = AdminOrganizationProjectRoleService{}
	r.Options = opts
	return
}

// Creates a custom role for a project.
func (r *AdminOrganizationProjectRoleService) New(ctx context.Context, projectID string, body AdminOrganizationProjectRoleNewParams, opts ...option.RequestOption) (res *Role, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if projectID == "" {
		err = errors.New("missing required project_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("projects/%s/roles", projectID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Updates an existing project role.
func (r *AdminOrganizationProjectRoleService) Update(ctx context.Context, projectID string, roleID string, body AdminOrganizationProjectRoleUpdateParams, opts ...option.RequestOption) (res *Role, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if projectID == "" {
		err = errors.New("missing required project_id parameter")
		return nil, err
	}
	if roleID == "" {
		err = errors.New("missing required role_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("projects/%s/roles/%s", projectID, roleID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Lists the roles configured for a project.
func (r *AdminOrganizationProjectRoleService) List(ctx context.Context, projectID string, query AdminOrganizationProjectRoleListParams, opts ...option.RequestOption) (res *pagination.CursorPage[Role], err error) {
	var raw *http.Response
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	if projectID == "" {
		err = errors.New("missing required project_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("projects/%s/roles", projectID)
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

// Lists the roles configured for a project.
func (r *AdminOrganizationProjectRoleService) ListAutoPaging(ctx context.Context, projectID string, query AdminOrganizationProjectRoleListParams, opts ...option.RequestOption) *pagination.CursorPageAutoPager[Role] {
	return pagination.NewCursorPageAutoPager(r.List(ctx, projectID, query, opts...))
}

// Deletes a custom role from a project.
func (r *AdminOrganizationProjectRoleService) Delete(ctx context.Context, projectID string, roleID string, opts ...option.RequestOption) (res *AdminOrganizationProjectRoleDeleteResponse, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if projectID == "" {
		err = errors.New("missing required project_id parameter")
		return nil, err
	}
	if roleID == "" {
		err = errors.New("missing required role_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("projects/%s/roles/%s", projectID, roleID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Confirmation payload returned after deleting a role.
type AdminOrganizationProjectRoleDeleteResponse struct {
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
func (r AdminOrganizationProjectRoleDeleteResponse) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationProjectRoleDeleteResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationProjectRoleNewParams struct {
	// Permissions to grant to the role.
	Permissions []string `json:"permissions,omitzero" api:"required"`
	// Unique name for the role.
	RoleName string `json:"role_name" api:"required"`
	// Optional description of the role.
	Description param.Opt[string] `json:"description,omitzero"`
	paramObj
}

func (r AdminOrganizationProjectRoleNewParams) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationProjectRoleNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationProjectRoleNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationProjectRoleUpdateParams struct {
	// New description for the role.
	Description param.Opt[string] `json:"description,omitzero"`
	// New name for the role.
	RoleName param.Opt[string] `json:"role_name,omitzero"`
	// Updated set of permissions for the role.
	Permissions []string `json:"permissions,omitzero"`
	paramObj
}

func (r AdminOrganizationProjectRoleUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationProjectRoleUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationProjectRoleUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationProjectRoleListParams struct {
	// Cursor for pagination. Provide the value from the previous response's `next`
	// field to continue listing roles.
	After param.Opt[string] `query:"after,omitzero" json:"-"`
	// A limit on the number of roles to return. Defaults to 1000.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Sort order for the returned roles.
	//
	// Any of "asc", "desc".
	Order AdminOrganizationProjectRoleListParamsOrder `query:"order,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [AdminOrganizationProjectRoleListParams]'s query parameters
// as `url.Values`.
func (r AdminOrganizationProjectRoleListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Sort order for the returned roles.
type AdminOrganizationProjectRoleListParamsOrder string

const (
	AdminOrganizationProjectRoleListParamsOrderAsc  AdminOrganizationProjectRoleListParamsOrder = "asc"
	AdminOrganizationProjectRoleListParamsOrderDesc AdminOrganizationProjectRoleListParamsOrder = "desc"
)

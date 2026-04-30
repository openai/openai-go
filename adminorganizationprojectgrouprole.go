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

// AdminOrganizationProjectGroupRoleService contains methods and other services
// that help with interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAdminOrganizationProjectGroupRoleService] method instead.
type AdminOrganizationProjectGroupRoleService struct {
	Options []option.RequestOption
}

// NewAdminOrganizationProjectGroupRoleService generates a new service that applies
// the given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewAdminOrganizationProjectGroupRoleService(opts ...option.RequestOption) (r AdminOrganizationProjectGroupRoleService) {
	r = AdminOrganizationProjectGroupRoleService{}
	r.Options = opts
	return
}

// Assigns a project role to a group within a project.
func (r *AdminOrganizationProjectGroupRoleService) New(ctx context.Context, projectID string, groupID string, body AdminOrganizationProjectGroupRoleNewParams, opts ...option.RequestOption) (res *AdminOrganizationProjectGroupRoleNewResponse, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if projectID == "" {
		err = errors.New("missing required project_id parameter")
		return nil, err
	}
	if groupID == "" {
		err = errors.New("missing required group_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("projects/%s/groups/%s/roles", projectID, groupID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Lists the project roles assigned to a group within a project.
func (r *AdminOrganizationProjectGroupRoleService) List(ctx context.Context, projectID string, groupID string, query AdminOrganizationProjectGroupRoleListParams, opts ...option.RequestOption) (res *pagination.CursorPage[AdminOrganizationProjectGroupRoleListResponse], err error) {
	var raw *http.Response
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	if projectID == "" {
		err = errors.New("missing required project_id parameter")
		return nil, err
	}
	if groupID == "" {
		err = errors.New("missing required group_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("projects/%s/groups/%s/roles", projectID, groupID)
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

// Lists the project roles assigned to a group within a project.
func (r *AdminOrganizationProjectGroupRoleService) ListAutoPaging(ctx context.Context, projectID string, groupID string, query AdminOrganizationProjectGroupRoleListParams, opts ...option.RequestOption) *pagination.CursorPageAutoPager[AdminOrganizationProjectGroupRoleListResponse] {
	return pagination.NewCursorPageAutoPager(r.List(ctx, projectID, groupID, query, opts...))
}

// Unassigns a project role from a group within a project.
func (r *AdminOrganizationProjectGroupRoleService) Delete(ctx context.Context, projectID string, groupID string, roleID string, opts ...option.RequestOption) (res *AdminOrganizationProjectGroupRoleDeleteResponse, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if projectID == "" {
		err = errors.New("missing required project_id parameter")
		return nil, err
	}
	if groupID == "" {
		err = errors.New("missing required group_id parameter")
		return nil, err
	}
	if roleID == "" {
		err = errors.New("missing required role_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("projects/%s/groups/%s/roles/%s", projectID, groupID, roleID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Role assignment linking a group to a role.
type AdminOrganizationProjectGroupRoleNewResponse struct {
	// Summary information about a group returned in role assignment responses.
	Group AdminOrganizationProjectGroupRoleNewResponseGroup `json:"group" api:"required"`
	// Always `group.role`.
	Object constant.GroupRole `json:"object" default:"group.role"`
	// Details about a role that can be assigned through the public Roles API.
	Role Role `json:"role" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Group       respjson.Field
		Object      respjson.Field
		Role        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationProjectGroupRoleNewResponse) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationProjectGroupRoleNewResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Summary information about a group returned in role assignment responses.
type AdminOrganizationProjectGroupRoleNewResponseGroup struct {
	// Identifier for the group.
	ID string `json:"id" api:"required"`
	// Unix timestamp (in seconds) when the group was created.
	CreatedAt int64 `json:"created_at" api:"required" format:"unixtime"`
	// Display name of the group.
	Name string `json:"name" api:"required"`
	// Always `group`.
	Object constant.Group `json:"object" default:"group"`
	// Whether the group is managed through SCIM.
	ScimManaged bool `json:"scim_managed" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		Name        respjson.Field
		Object      respjson.Field
		ScimManaged respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationProjectGroupRoleNewResponseGroup) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationProjectGroupRoleNewResponseGroup) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Detailed information about a role assignment entry returned when listing
// assignments.
type AdminOrganizationProjectGroupRoleListResponse struct {
	// Identifier for the role.
	ID string `json:"id" api:"required"`
	// When the role was created.
	CreatedAt int64 `json:"created_at" api:"required" format:"unixtime"`
	// Identifier of the actor who created the role.
	CreatedBy string `json:"created_by" api:"required"`
	// User details for the actor that created the role, when available.
	CreatedByUserObj map[string]any `json:"created_by_user_obj" api:"required"`
	// Description of the role.
	Description string `json:"description" api:"required"`
	// Arbitrary metadata stored on the role.
	Metadata map[string]any `json:"metadata" api:"required"`
	// Name of the role.
	Name string `json:"name" api:"required"`
	// Permissions associated with the role.
	Permissions []string `json:"permissions" api:"required"`
	// Whether the role is predefined by OpenAI.
	PredefinedRole bool `json:"predefined_role" api:"required"`
	// Resource type the role applies to.
	ResourceType string `json:"resource_type" api:"required"`
	// When the role was last updated.
	UpdatedAt int64 `json:"updated_at" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID               respjson.Field
		CreatedAt        respjson.Field
		CreatedBy        respjson.Field
		CreatedByUserObj respjson.Field
		Description      respjson.Field
		Metadata         respjson.Field
		Name             respjson.Field
		Permissions      respjson.Field
		PredefinedRole   respjson.Field
		ResourceType     respjson.Field
		UpdatedAt        respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationProjectGroupRoleListResponse) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationProjectGroupRoleListResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Confirmation payload returned after unassigning a role.
type AdminOrganizationProjectGroupRoleDeleteResponse struct {
	// Whether the assignment was removed.
	Deleted bool `json:"deleted" api:"required"`
	// Identifier for the deleted assignment, such as `group.role.deleted` or
	// `user.role.deleted`.
	Object string `json:"object" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Deleted     respjson.Field
		Object      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationProjectGroupRoleDeleteResponse) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationProjectGroupRoleDeleteResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationProjectGroupRoleNewParams struct {
	// Identifier of the role to assign.
	RoleID string `json:"role_id" api:"required"`
	paramObj
}

func (r AdminOrganizationProjectGroupRoleNewParams) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationProjectGroupRoleNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationProjectGroupRoleNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationProjectGroupRoleListParams struct {
	// Cursor for pagination. Provide the value from the previous response's `next`
	// field to continue listing project roles.
	After param.Opt[string] `query:"after,omitzero" json:"-"`
	// A limit on the number of project role assignments to return.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Sort order for the returned project roles.
	//
	// Any of "asc", "desc".
	Order AdminOrganizationProjectGroupRoleListParamsOrder `query:"order,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [AdminOrganizationProjectGroupRoleListParams]'s query
// parameters as `url.Values`.
func (r AdminOrganizationProjectGroupRoleListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Sort order for the returned project roles.
type AdminOrganizationProjectGroupRoleListParamsOrder string

const (
	AdminOrganizationProjectGroupRoleListParamsOrderAsc  AdminOrganizationProjectGroupRoleListParamsOrder = "asc"
	AdminOrganizationProjectGroupRoleListParamsOrderDesc AdminOrganizationProjectGroupRoleListParamsOrder = "desc"
)

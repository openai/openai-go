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

// AdminOrganizationUserRoleService contains methods and other services that help
// with interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAdminOrganizationUserRoleService] method instead.
type AdminOrganizationUserRoleService struct {
	Options []option.RequestOption
}

// NewAdminOrganizationUserRoleService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewAdminOrganizationUserRoleService(opts ...option.RequestOption) (r AdminOrganizationUserRoleService) {
	r = AdminOrganizationUserRoleService{}
	r.Options = opts
	return
}

// Assigns an organization role to a user within the organization.
func (r *AdminOrganizationUserRoleService) New(ctx context.Context, userID string, body AdminOrganizationUserRoleNewParams, opts ...option.RequestOption) (res *AdminOrganizationUserRoleNewResponse, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if userID == "" {
		err = errors.New("missing required user_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("organization/users/%s/roles", userID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Lists the organization roles assigned to a user within the organization.
func (r *AdminOrganizationUserRoleService) List(ctx context.Context, userID string, query AdminOrganizationUserRoleListParams, opts ...option.RequestOption) (res *pagination.NextCursorPage[AdminOrganizationUserRoleListResponse], err error) {
	var raw *http.Response
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	if userID == "" {
		err = errors.New("missing required user_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("organization/users/%s/roles", userID)
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

// Lists the organization roles assigned to a user within the organization.
func (r *AdminOrganizationUserRoleService) ListAutoPaging(ctx context.Context, userID string, query AdminOrganizationUserRoleListParams, opts ...option.RequestOption) *pagination.NextCursorPageAutoPager[AdminOrganizationUserRoleListResponse] {
	return pagination.NewNextCursorPageAutoPager(r.List(ctx, userID, query, opts...))
}

// Unassigns an organization role from a user within the organization.
func (r *AdminOrganizationUserRoleService) Delete(ctx context.Context, userID string, roleID string, opts ...option.RequestOption) (res *AdminOrganizationUserRoleDeleteResponse, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if userID == "" {
		err = errors.New("missing required user_id parameter")
		return nil, err
	}
	if roleID == "" {
		err = errors.New("missing required role_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("organization/users/%s/roles/%s", userID, roleID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Role assignment linking a user to a role.
type AdminOrganizationUserRoleNewResponse struct {
	// Always `user.role`.
	Object constant.UserRole `json:"object" default:"user.role"`
	// Details about a role that can be assigned through the public Roles API.
	Role Role `json:"role" api:"required"`
	// Represents an individual `user` within an organization.
	User OrganizationUser `json:"user" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Object      respjson.Field
		Role        respjson.Field
		User        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUserRoleNewResponse) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationUserRoleNewResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Detailed information about a role assignment entry returned when listing
// assignments.
type AdminOrganizationUserRoleListResponse struct {
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
func (r AdminOrganizationUserRoleListResponse) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationUserRoleListResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Confirmation payload returned after unassigning a role.
type AdminOrganizationUserRoleDeleteResponse struct {
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
func (r AdminOrganizationUserRoleDeleteResponse) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationUserRoleDeleteResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationUserRoleNewParams struct {
	// Identifier of the role to assign.
	RoleID string `json:"role_id" api:"required"`
	paramObj
}

func (r AdminOrganizationUserRoleNewParams) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationUserRoleNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationUserRoleNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationUserRoleListParams struct {
	// Cursor for pagination. Provide the value from the previous response's `next`
	// field to continue listing organization roles.
	After param.Opt[string] `query:"after,omitzero" json:"-"`
	// A limit on the number of organization role assignments to return.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Sort order for the returned organization roles.
	//
	// Any of "asc", "desc".
	Order AdminOrganizationUserRoleListParamsOrder `query:"order,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [AdminOrganizationUserRoleListParams]'s query parameters as
// `url.Values`.
func (r AdminOrganizationUserRoleListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Sort order for the returned organization roles.
type AdminOrganizationUserRoleListParamsOrder string

const (
	AdminOrganizationUserRoleListParamsOrderAsc  AdminOrganizationUserRoleListParamsOrder = "asc"
	AdminOrganizationUserRoleListParamsOrderDesc AdminOrganizationUserRoleListParamsOrder = "desc"
)

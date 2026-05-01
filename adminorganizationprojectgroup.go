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

// AdminOrganizationProjectGroupService contains methods and other services that
// help with interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAdminOrganizationProjectGroupService] method instead.
type AdminOrganizationProjectGroupService struct {
	Options []option.RequestOption
	Roles   AdminOrganizationProjectGroupRoleService
}

// NewAdminOrganizationProjectGroupService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewAdminOrganizationProjectGroupService(opts ...option.RequestOption) (r AdminOrganizationProjectGroupService) {
	r = AdminOrganizationProjectGroupService{}
	r.Options = opts
	r.Roles = NewAdminOrganizationProjectGroupRoleService(opts...)
	return
}

// Grants a group access to a project.
func (r *AdminOrganizationProjectGroupService) New(ctx context.Context, projectID string, body AdminOrganizationProjectGroupNewParams, opts ...option.RequestOption) (res *ProjectGroup, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if projectID == "" {
		err = errors.New("missing required project_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("organization/projects/%s/groups", projectID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Lists the groups that have access to a project.
func (r *AdminOrganizationProjectGroupService) List(ctx context.Context, projectID string, query AdminOrganizationProjectGroupListParams, opts ...option.RequestOption) (res *pagination.NextCursorPage[ProjectGroup], err error) {
	var raw *http.Response
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	if projectID == "" {
		err = errors.New("missing required project_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("organization/projects/%s/groups", projectID)
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

// Lists the groups that have access to a project.
func (r *AdminOrganizationProjectGroupService) ListAutoPaging(ctx context.Context, projectID string, query AdminOrganizationProjectGroupListParams, opts ...option.RequestOption) *pagination.NextCursorPageAutoPager[ProjectGroup] {
	return pagination.NewNextCursorPageAutoPager(r.List(ctx, projectID, query, opts...))
}

// Revokes a group's access to a project.
func (r *AdminOrganizationProjectGroupService) Delete(ctx context.Context, projectID string, groupID string, opts ...option.RequestOption) (res *AdminOrganizationProjectGroupDeleteResponse, err error) {
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
	path := fmt.Sprintf("organization/projects/%s/groups/%s", projectID, groupID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Details about a group's membership in a project.
type ProjectGroup struct {
	// Unix timestamp (in seconds) when the group was granted project access.
	CreatedAt int64 `json:"created_at" api:"required" format:"unixtime"`
	// Identifier of the group that has access to the project.
	GroupID string `json:"group_id" api:"required"`
	// Display name of the group.
	GroupName string `json:"group_name" api:"required"`
	// Always `project.group`.
	Object constant.ProjectGroup `json:"object" default:"project.group"`
	// Identifier of the project.
	ProjectID string `json:"project_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		CreatedAt   respjson.Field
		GroupID     respjson.Field
		GroupName   respjson.Field
		Object      respjson.Field
		ProjectID   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProjectGroup) RawJSON() string { return r.JSON.raw }
func (r *ProjectGroup) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Confirmation payload returned after removing a group from a project.
type AdminOrganizationProjectGroupDeleteResponse struct {
	// Whether the group membership in the project was removed.
	Deleted bool `json:"deleted" api:"required"`
	// Always `project.group.deleted`.
	Object constant.ProjectGroupDeleted `json:"object" default:"project.group.deleted"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Deleted     respjson.Field
		Object      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationProjectGroupDeleteResponse) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationProjectGroupDeleteResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationProjectGroupNewParams struct {
	// Identifier of the group to add to the project.
	GroupID string `json:"group_id" api:"required"`
	// Identifier of the project role to grant to the group.
	Role string `json:"role" api:"required"`
	paramObj
}

func (r AdminOrganizationProjectGroupNewParams) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationProjectGroupNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationProjectGroupNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationProjectGroupListParams struct {
	// Cursor for pagination. Provide the ID of the last group from the previous
	// response to fetch the next page.
	After param.Opt[string] `query:"after,omitzero" json:"-"`
	// A limit on the number of project groups to return. Defaults to 20.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Sort order for the returned groups.
	//
	// Any of "asc", "desc".
	Order AdminOrganizationProjectGroupListParamsOrder `query:"order,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [AdminOrganizationProjectGroupListParams]'s query parameters
// as `url.Values`.
func (r AdminOrganizationProjectGroupListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Sort order for the returned groups.
type AdminOrganizationProjectGroupListParamsOrder string

const (
	AdminOrganizationProjectGroupListParamsOrderAsc  AdminOrganizationProjectGroupListParamsOrder = "asc"
	AdminOrganizationProjectGroupListParamsOrderDesc AdminOrganizationProjectGroupListParamsOrder = "desc"
)

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

// AdminOrganizationGroupService contains methods and other services that help with
// interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAdminOrganizationGroupService] method instead.
type AdminOrganizationGroupService struct {
	Options []option.RequestOption
	Users   AdminOrganizationGroupUserService
	Roles   AdminOrganizationGroupRoleService
}

// NewAdminOrganizationGroupService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewAdminOrganizationGroupService(opts ...option.RequestOption) (r AdminOrganizationGroupService) {
	r = AdminOrganizationGroupService{}
	r.Options = opts
	r.Users = NewAdminOrganizationGroupUserService(opts...)
	r.Roles = NewAdminOrganizationGroupRoleService(opts...)
	return
}

// Creates a new group in the organization.
func (r *AdminOrganizationGroupService) New(ctx context.Context, body AdminOrganizationGroupNewParams, opts ...option.RequestOption) (res *Group, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	path := "organization/groups"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Updates a group's information.
func (r *AdminOrganizationGroupService) Update(ctx context.Context, groupID string, body AdminOrganizationGroupUpdateParams, opts ...option.RequestOption) (res *AdminOrganizationGroupUpdateResponse, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if groupID == "" {
		err = errors.New("missing required group_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("organization/groups/%s", groupID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Lists all groups in the organization.
func (r *AdminOrganizationGroupService) List(ctx context.Context, query AdminOrganizationGroupListParams, opts ...option.RequestOption) (res *pagination.NextCursorPage[Group], err error) {
	var raw *http.Response
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "organization/groups"
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

// Lists all groups in the organization.
func (r *AdminOrganizationGroupService) ListAutoPaging(ctx context.Context, query AdminOrganizationGroupListParams, opts ...option.RequestOption) *pagination.NextCursorPageAutoPager[Group] {
	return pagination.NewNextCursorPageAutoPager(r.List(ctx, query, opts...))
}

// Deletes a group from the organization.
func (r *AdminOrganizationGroupService) Delete(ctx context.Context, groupID string, opts ...option.RequestOption) (res *AdminOrganizationGroupDeleteResponse, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if groupID == "" {
		err = errors.New("missing required group_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("organization/groups/%s", groupID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Details about an organization group.
type Group struct {
	// Identifier for the group.
	ID string `json:"id" api:"required"`
	// Unix timestamp (in seconds) when the group was created.
	CreatedAt int64 `json:"created_at" api:"required" format:"unixtime"`
	// Whether the group is managed through SCIM and controlled by your identity
	// provider.
	IsScimManaged bool `json:"is_scim_managed" api:"required"`
	// Display name of the group.
	Name string `json:"name" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID            respjson.Field
		CreatedAt     respjson.Field
		IsScimManaged respjson.Field
		Name          respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r Group) RawJSON() string { return r.JSON.raw }
func (r *Group) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Response returned after updating a group.
type AdminOrganizationGroupUpdateResponse struct {
	// Identifier for the group.
	ID string `json:"id" api:"required"`
	// Unix timestamp (in seconds) when the group was created.
	CreatedAt int64 `json:"created_at" api:"required" format:"unixtime"`
	// Whether the group is managed through SCIM and controlled by your identity
	// provider.
	IsScimManaged bool `json:"is_scim_managed" api:"required"`
	// Updated display name for the group.
	Name string `json:"name" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID            respjson.Field
		CreatedAt     respjson.Field
		IsScimManaged respjson.Field
		Name          respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationGroupUpdateResponse) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationGroupUpdateResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Confirmation payload returned after deleting a group.
type AdminOrganizationGroupDeleteResponse struct {
	// Identifier of the deleted group.
	ID string `json:"id" api:"required"`
	// Whether the group was deleted.
	Deleted bool `json:"deleted" api:"required"`
	// Always `group.deleted`.
	Object constant.GroupDeleted `json:"object" default:"group.deleted"`
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
func (r AdminOrganizationGroupDeleteResponse) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationGroupDeleteResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationGroupNewParams struct {
	// Human readable name for the group.
	Name string `json:"name" api:"required"`
	paramObj
}

func (r AdminOrganizationGroupNewParams) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationGroupNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationGroupNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationGroupUpdateParams struct {
	// New display name for the group.
	Name string `json:"name" api:"required"`
	paramObj
}

func (r AdminOrganizationGroupUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationGroupUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationGroupUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationGroupListParams struct {
	// A cursor for use in pagination. `after` is a group ID that defines your place in
	// the list. For instance, if you make a list request and receive 100 objects,
	// ending with group_abc, your subsequent call can include `after=group_abc` in
	// order to fetch the next page of the list.
	After param.Opt[string] `query:"after,omitzero" json:"-"`
	// A limit on the number of groups to be returned. Limit can range between 0 and
	// 1000, and the default is 100.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Specifies the sort order of the returned groups.
	//
	// Any of "asc", "desc".
	Order AdminOrganizationGroupListParamsOrder `query:"order,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [AdminOrganizationGroupListParams]'s query parameters as
// `url.Values`.
func (r AdminOrganizationGroupListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Specifies the sort order of the returned groups.
type AdminOrganizationGroupListParamsOrder string

const (
	AdminOrganizationGroupListParamsOrderAsc  AdminOrganizationGroupListParamsOrder = "asc"
	AdminOrganizationGroupListParamsOrderDesc AdminOrganizationGroupListParamsOrder = "desc"
)

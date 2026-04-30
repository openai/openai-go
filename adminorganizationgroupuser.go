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

// AdminOrganizationGroupUserService contains methods and other services that help
// with interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAdminOrganizationGroupUserService] method instead.
type AdminOrganizationGroupUserService struct {
	Options []option.RequestOption
}

// NewAdminOrganizationGroupUserService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewAdminOrganizationGroupUserService(opts ...option.RequestOption) (r AdminOrganizationGroupUserService) {
	r = AdminOrganizationGroupUserService{}
	r.Options = opts
	return
}

// Adds a user to a group.
func (r *AdminOrganizationGroupUserService) New(ctx context.Context, groupID string, body AdminOrganizationGroupUserNewParams, opts ...option.RequestOption) (res *AdminOrganizationGroupUserNewResponse, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if groupID == "" {
		err = errors.New("missing required group_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("organization/groups/%s/users", groupID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Lists the users assigned to a group.
func (r *AdminOrganizationGroupUserService) List(ctx context.Context, groupID string, query AdminOrganizationGroupUserListParams, opts ...option.RequestOption) (res *pagination.CursorPage[OrganizationUser], err error) {
	var raw *http.Response
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	if groupID == "" {
		err = errors.New("missing required group_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("organization/groups/%s/users", groupID)
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

// Lists the users assigned to a group.
func (r *AdminOrganizationGroupUserService) ListAutoPaging(ctx context.Context, groupID string, query AdminOrganizationGroupUserListParams, opts ...option.RequestOption) *pagination.CursorPageAutoPager[OrganizationUser] {
	return pagination.NewCursorPageAutoPager(r.List(ctx, groupID, query, opts...))
}

// Removes a user from a group.
func (r *AdminOrganizationGroupUserService) Delete(ctx context.Context, groupID string, userID string, opts ...option.RequestOption) (res *AdminOrganizationGroupUserDeleteResponse, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if groupID == "" {
		err = errors.New("missing required group_id parameter")
		return nil, err
	}
	if userID == "" {
		err = errors.New("missing required user_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("organization/groups/%s/users/%s", groupID, userID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Confirmation payload returned after adding a user to a group.
type AdminOrganizationGroupUserNewResponse struct {
	// Identifier of the group the user was added to.
	GroupID string `json:"group_id" api:"required"`
	// Always `group.user`.
	Object constant.GroupUser `json:"object" default:"group.user"`
	// Identifier of the user that was added.
	UserID string `json:"user_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		GroupID     respjson.Field
		Object      respjson.Field
		UserID      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationGroupUserNewResponse) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationGroupUserNewResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Confirmation payload returned after removing a user from a group.
type AdminOrganizationGroupUserDeleteResponse struct {
	// Whether the group membership was removed.
	Deleted bool `json:"deleted" api:"required"`
	// Always `group.user.deleted`.
	Object constant.GroupUserDeleted `json:"object" default:"group.user.deleted"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Deleted     respjson.Field
		Object      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationGroupUserDeleteResponse) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationGroupUserDeleteResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationGroupUserNewParams struct {
	// Identifier of the user to add to the group.
	UserID string `json:"user_id" api:"required"`
	paramObj
}

func (r AdminOrganizationGroupUserNewParams) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationGroupUserNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationGroupUserNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationGroupUserListParams struct {
	// A cursor for use in pagination. Provide the ID of the last user from the
	// previous list response to retrieve the next page.
	After param.Opt[string] `query:"after,omitzero" json:"-"`
	// A limit on the number of users to be returned. Limit can range between 0 and
	// 1000, and the default is 100.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Specifies the sort order of users in the list.
	//
	// Any of "asc", "desc".
	Order AdminOrganizationGroupUserListParamsOrder `query:"order,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [AdminOrganizationGroupUserListParams]'s query parameters as
// `url.Values`.
func (r AdminOrganizationGroupUserListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Specifies the sort order of users in the list.
type AdminOrganizationGroupUserListParamsOrder string

const (
	AdminOrganizationGroupUserListParamsOrderAsc  AdminOrganizationGroupUserListParamsOrder = "asc"
	AdminOrganizationGroupUserListParamsOrderDesc AdminOrganizationGroupUserListParamsOrder = "desc"
)

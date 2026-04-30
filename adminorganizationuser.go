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

// AdminOrganizationUserService contains methods and other services that help with
// interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAdminOrganizationUserService] method instead.
type AdminOrganizationUserService struct {
	Options []option.RequestOption
	Roles   AdminOrganizationUserRoleService
}

// NewAdminOrganizationUserService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewAdminOrganizationUserService(opts ...option.RequestOption) (r AdminOrganizationUserService) {
	r = AdminOrganizationUserService{}
	r.Options = opts
	r.Roles = NewAdminOrganizationUserRoleService(opts...)
	return
}

// Retrieves a user by their identifier.
func (r *AdminOrganizationUserService) Get(ctx context.Context, userID string, opts ...option.RequestOption) (res *OrganizationUser, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if userID == "" {
		err = errors.New("missing required user_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("organization/users/%s", userID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Modifies a user's role in the organization.
func (r *AdminOrganizationUserService) Update(ctx context.Context, userID string, body AdminOrganizationUserUpdateParams, opts ...option.RequestOption) (res *OrganizationUser, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if userID == "" {
		err = errors.New("missing required user_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("organization/users/%s", userID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Lists all of the users in the organization.
func (r *AdminOrganizationUserService) List(ctx context.Context, query AdminOrganizationUserListParams, opts ...option.RequestOption) (res *pagination.ConversationCursorPage[OrganizationUser], err error) {
	var raw *http.Response
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "organization/users"
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

// Lists all of the users in the organization.
func (r *AdminOrganizationUserService) ListAutoPaging(ctx context.Context, query AdminOrganizationUserListParams, opts ...option.RequestOption) *pagination.ConversationCursorPageAutoPager[OrganizationUser] {
	return pagination.NewConversationCursorPageAutoPager(r.List(ctx, query, opts...))
}

// Deletes a user from the organization.
func (r *AdminOrganizationUserService) Delete(ctx context.Context, userID string, opts ...option.RequestOption) (res *AdminOrganizationUserDeleteResponse, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if userID == "" {
		err = errors.New("missing required user_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("organization/users/%s", userID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Represents an individual `user` within an organization.
type OrganizationUser struct {
	// The identifier, which can be referenced in API endpoints
	ID string `json:"id" api:"required"`
	// The Unix timestamp (in seconds) of when the user was added.
	AddedAt int64 `json:"added_at" api:"required" format:"unixtime"`
	// The email address of the user
	Email string `json:"email" api:"required"`
	// The name of the user
	Name string `json:"name" api:"required"`
	// The object type, which is always `organization.user`
	Object constant.OrganizationUser `json:"object" default:"organization.user"`
	// `owner` or `reader`
	//
	// Any of "owner", "reader".
	Role OrganizationUserRole `json:"role" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		AddedAt     respjson.Field
		Email       respjson.Field
		Name        respjson.Field
		Object      respjson.Field
		Role        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r OrganizationUser) RawJSON() string { return r.JSON.raw }
func (r *OrganizationUser) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// `owner` or `reader`
type OrganizationUserRole string

const (
	OrganizationUserRoleOwner  OrganizationUserRole = "owner"
	OrganizationUserRoleReader OrganizationUserRole = "reader"
)

type AdminOrganizationUserDeleteResponse struct {
	ID      string                           `json:"id" api:"required"`
	Deleted bool                             `json:"deleted" api:"required"`
	Object  constant.OrganizationUserDeleted `json:"object" default:"organization.user.deleted"`
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
func (r AdminOrganizationUserDeleteResponse) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationUserDeleteResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationUserUpdateParams struct {
	// `owner` or `reader`
	//
	// Any of "owner", "reader".
	Role AdminOrganizationUserUpdateParamsRole `json:"role,omitzero" api:"required"`
	paramObj
}

func (r AdminOrganizationUserUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationUserUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationUserUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// `owner` or `reader`
type AdminOrganizationUserUpdateParamsRole string

const (
	AdminOrganizationUserUpdateParamsRoleOwner  AdminOrganizationUserUpdateParamsRole = "owner"
	AdminOrganizationUserUpdateParamsRoleReader AdminOrganizationUserUpdateParamsRole = "reader"
)

type AdminOrganizationUserListParams struct {
	// A cursor for use in pagination. `after` is an object ID that defines your place
	// in the list. For instance, if you make a list request and receive 100 objects,
	// ending with obj_foo, your subsequent call can include after=obj_foo in order to
	// fetch the next page of the list.
	After param.Opt[string] `query:"after,omitzero" json:"-"`
	// A limit on the number of objects to be returned. Limit can range between 1 and
	// 100, and the default is 20.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Filter by the email address of users.
	Emails []string `query:"emails,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [AdminOrganizationUserListParams]'s query parameters as
// `url.Values`.
func (r AdminOrganizationUserListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

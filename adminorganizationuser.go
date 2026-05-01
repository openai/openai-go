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
	// The object type, which is always `organization.user`
	Object constant.OrganizationUser `json:"object" default:"organization.user"`
	// The Unix timestamp (in seconds) of the user's last API key usage.
	APIKeyLastUsedAt int64 `json:"api_key_last_used_at" api:"nullable" format:"unixtime"`
	// The Unix timestamp (in seconds) of when the user was created.
	Created int64 `json:"created" format:"unixtime"`
	// The developer persona metadata for the user.
	DeveloperPersona string `json:"developer_persona" api:"nullable"`
	// The email address of the user
	Email string `json:"email" api:"nullable"`
	// Whether this is the organization's default user.
	IsDefault bool `json:"is_default"`
	// Whether the user is an authorized purchaser for Scale Tier.
	IsScaleTierAuthorizedPurchaser bool `json:"is_scale_tier_authorized_purchaser" api:"nullable"`
	// Whether the user is managed through SCIM.
	IsScimManaged bool `json:"is_scim_managed"`
	// Whether the user is a service account.
	IsServiceAccount bool `json:"is_service_account"`
	// The name of the user
	Name string `json:"name" api:"nullable"`
	// Projects associated with the user, if included.
	Projects OrganizationUserProjects `json:"projects" api:"nullable"`
	// `owner` or `reader`
	Role string `json:"role" api:"nullable"`
	// The technical level metadata for the user.
	TechnicalLevel string `json:"technical_level" api:"nullable"`
	// Nested user details.
	User OrganizationUserUser `json:"user"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                             respjson.Field
		AddedAt                        respjson.Field
		Object                         respjson.Field
		APIKeyLastUsedAt               respjson.Field
		Created                        respjson.Field
		DeveloperPersona               respjson.Field
		Email                          respjson.Field
		IsDefault                      respjson.Field
		IsScaleTierAuthorizedPurchaser respjson.Field
		IsScimManaged                  respjson.Field
		IsServiceAccount               respjson.Field
		Name                           respjson.Field
		Projects                       respjson.Field
		Role                           respjson.Field
		TechnicalLevel                 respjson.Field
		User                           respjson.Field
		ExtraFields                    map[string]respjson.Field
		raw                            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r OrganizationUser) RawJSON() string { return r.JSON.raw }
func (r *OrganizationUser) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Projects associated with the user, if included.
type OrganizationUserProjects struct {
	Data   []OrganizationUserProjectsData `json:"data" api:"required"`
	Object constant.List                  `json:"object" default:"list"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Data        respjson.Field
		Object      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r OrganizationUserProjects) RawJSON() string { return r.JSON.raw }
func (r *OrganizationUserProjects) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrganizationUserProjectsData struct {
	ID   string `json:"id" api:"nullable"`
	Name string `json:"name" api:"nullable"`
	Role string `json:"role" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Name        respjson.Field
		Role        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r OrganizationUserProjectsData) RawJSON() string { return r.JSON.raw }
func (r *OrganizationUserProjectsData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Nested user details.
type OrganizationUserUser struct {
	ID       string        `json:"id" api:"required"`
	Object   constant.User `json:"object" default:"user"`
	Banned   bool          `json:"banned" api:"nullable"`
	BannedAt int64         `json:"banned_at" api:"nullable" format:"unixtime"`
	Email    string        `json:"email" api:"nullable"`
	Enabled  bool          `json:"enabled" api:"nullable"`
	Name     string        `json:"name" api:"nullable"`
	Picture  string        `json:"picture" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Object      respjson.Field
		Banned      respjson.Field
		BannedAt    respjson.Field
		Email       respjson.Field
		Enabled     respjson.Field
		Name        respjson.Field
		Picture     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r OrganizationUserUser) RawJSON() string { return r.JSON.raw }
func (r *OrganizationUserUser) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

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
	// Developer persona metadata.
	DeveloperPersona param.Opt[string] `json:"developer_persona,omitzero"`
	// `owner` or `reader`
	Role param.Opt[string] `json:"role,omitzero"`
	// Role ID to assign to the user.
	RoleID param.Opt[string] `json:"role_id,omitzero"`
	// Technical level metadata.
	TechnicalLevel param.Opt[string] `json:"technical_level,omitzero"`
	paramObj
}

func (r AdminOrganizationUserUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationUserUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationUserUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

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

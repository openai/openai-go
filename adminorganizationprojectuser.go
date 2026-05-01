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

// AdminOrganizationProjectUserService contains methods and other services that
// help with interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAdminOrganizationProjectUserService] method instead.
type AdminOrganizationProjectUserService struct {
	Options []option.RequestOption
	Roles   AdminOrganizationProjectUserRoleService
}

// NewAdminOrganizationProjectUserService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewAdminOrganizationProjectUserService(opts ...option.RequestOption) (r AdminOrganizationProjectUserService) {
	r = AdminOrganizationProjectUserService{}
	r.Options = opts
	r.Roles = NewAdminOrganizationProjectUserRoleService(opts...)
	return
}

// Adds a user to the project. Users must already be members of the organization to
// be added to a project.
func (r *AdminOrganizationProjectUserService) New(ctx context.Context, projectID string, body AdminOrganizationProjectUserNewParams, opts ...option.RequestOption) (res *ProjectUser, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if projectID == "" {
		err = errors.New("missing required project_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("organization/projects/%s/users", projectID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Retrieves a user in the project.
func (r *AdminOrganizationProjectUserService) Get(ctx context.Context, projectID string, userID string, opts ...option.RequestOption) (res *ProjectUser, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if projectID == "" {
		err = errors.New("missing required project_id parameter")
		return nil, err
	}
	if userID == "" {
		err = errors.New("missing required user_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("organization/projects/%s/users/%s", projectID, userID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Modifies a user's role in the project.
func (r *AdminOrganizationProjectUserService) Update(ctx context.Context, projectID string, userID string, body AdminOrganizationProjectUserUpdateParams, opts ...option.RequestOption) (res *ProjectUser, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if projectID == "" {
		err = errors.New("missing required project_id parameter")
		return nil, err
	}
	if userID == "" {
		err = errors.New("missing required user_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("organization/projects/%s/users/%s", projectID, userID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Returns a list of users in the project.
func (r *AdminOrganizationProjectUserService) List(ctx context.Context, projectID string, query AdminOrganizationProjectUserListParams, opts ...option.RequestOption) (res *pagination.ConversationCursorPage[ProjectUser], err error) {
	var raw *http.Response
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	if projectID == "" {
		err = errors.New("missing required project_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("organization/projects/%s/users", projectID)
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

// Returns a list of users in the project.
func (r *AdminOrganizationProjectUserService) ListAutoPaging(ctx context.Context, projectID string, query AdminOrganizationProjectUserListParams, opts ...option.RequestOption) *pagination.ConversationCursorPageAutoPager[ProjectUser] {
	return pagination.NewConversationCursorPageAutoPager(r.List(ctx, projectID, query, opts...))
}

// Deletes a user from the project.
//
// Returns confirmation of project user deletion, or an error if the project is
// archived (archived projects have no users).
func (r *AdminOrganizationProjectUserService) Delete(ctx context.Context, projectID string, userID string, opts ...option.RequestOption) (res *AdminOrganizationProjectUserDeleteResponse, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if projectID == "" {
		err = errors.New("missing required project_id parameter")
		return nil, err
	}
	if userID == "" {
		err = errors.New("missing required user_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("organization/projects/%s/users/%s", projectID, userID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Represents an individual user in a project.
type ProjectUser struct {
	// The identifier, which can be referenced in API endpoints
	ID string `json:"id" api:"required"`
	// The Unix timestamp (in seconds) of when the project was added.
	AddedAt int64 `json:"added_at" api:"required" format:"unixtime"`
	// The object type, which is always `organization.project.user`
	Object constant.OrganizationProjectUser `json:"object" default:"organization.project.user"`
	// `owner` or `member`
	Role string `json:"role" api:"required"`
	// The email address of the user
	Email string `json:"email" api:"nullable"`
	// The name of the user
	Name string `json:"name" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		AddedAt     respjson.Field
		Object      respjson.Field
		Role        respjson.Field
		Email       respjson.Field
		Name        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProjectUser) RawJSON() string { return r.JSON.raw }
func (r *ProjectUser) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationProjectUserDeleteResponse struct {
	ID      string                                  `json:"id" api:"required"`
	Deleted bool                                    `json:"deleted" api:"required"`
	Object  constant.OrganizationProjectUserDeleted `json:"object" default:"organization.project.user.deleted"`
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
func (r AdminOrganizationProjectUserDeleteResponse) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationProjectUserDeleteResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationProjectUserNewParams struct {
	// `owner` or `member`
	Role string `json:"role" api:"required"`
	// Email of the user to add.
	Email param.Opt[string] `json:"email,omitzero"`
	// The ID of the user.
	UserID param.Opt[string] `json:"user_id,omitzero"`
	paramObj
}

func (r AdminOrganizationProjectUserNewParams) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationProjectUserNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationProjectUserNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationProjectUserUpdateParams struct {
	// `owner` or `member`
	Role param.Opt[string] `json:"role,omitzero"`
	paramObj
}

func (r AdminOrganizationProjectUserUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationProjectUserUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationProjectUserUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationProjectUserListParams struct {
	// A cursor for use in pagination. `after` is an object ID that defines your place
	// in the list. For instance, if you make a list request and receive 100 objects,
	// ending with obj_foo, your subsequent call can include after=obj_foo in order to
	// fetch the next page of the list.
	After param.Opt[string] `query:"after,omitzero" json:"-"`
	// A limit on the number of objects to be returned. Limit can range between 1 and
	// 100, and the default is 20.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [AdminOrganizationProjectUserListParams]'s query parameters
// as `url.Values`.
func (r AdminOrganizationProjectUserListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

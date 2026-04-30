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

// AdminOrganizationProjectServiceAccountService contains methods and other
// services that help with interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAdminOrganizationProjectServiceAccountService] method instead.
type AdminOrganizationProjectServiceAccountService struct {
	Options []option.RequestOption
}

// NewAdminOrganizationProjectServiceAccountService generates a new service that
// applies the given options to each request. These options are applied after the
// parent client's options (if there is one), and before any request-specific
// options.
func NewAdminOrganizationProjectServiceAccountService(opts ...option.RequestOption) (r AdminOrganizationProjectServiceAccountService) {
	r = AdminOrganizationProjectServiceAccountService{}
	r.Options = opts
	return
}

// Creates a new service account in the project. This also returns an unredacted
// API key for the service account.
func (r *AdminOrganizationProjectServiceAccountService) New(ctx context.Context, projectID string, body AdminOrganizationProjectServiceAccountNewParams, opts ...option.RequestOption) (res *AdminOrganizationProjectServiceAccountNewResponse, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if projectID == "" {
		err = errors.New("missing required project_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("organization/projects/%s/service_accounts", projectID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Retrieves a service account in the project.
func (r *AdminOrganizationProjectServiceAccountService) Get(ctx context.Context, projectID string, serviceAccountID string, opts ...option.RequestOption) (res *ProjectServiceAccount, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if projectID == "" {
		err = errors.New("missing required project_id parameter")
		return nil, err
	}
	if serviceAccountID == "" {
		err = errors.New("missing required service_account_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("organization/projects/%s/service_accounts/%s", projectID, serviceAccountID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Returns a list of service accounts in the project.
func (r *AdminOrganizationProjectServiceAccountService) List(ctx context.Context, projectID string, query AdminOrganizationProjectServiceAccountListParams, opts ...option.RequestOption) (res *pagination.ConversationCursorPage[ProjectServiceAccount], err error) {
	var raw *http.Response
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	if projectID == "" {
		err = errors.New("missing required project_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("organization/projects/%s/service_accounts", projectID)
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

// Returns a list of service accounts in the project.
func (r *AdminOrganizationProjectServiceAccountService) ListAutoPaging(ctx context.Context, projectID string, query AdminOrganizationProjectServiceAccountListParams, opts ...option.RequestOption) *pagination.ConversationCursorPageAutoPager[ProjectServiceAccount] {
	return pagination.NewConversationCursorPageAutoPager(r.List(ctx, projectID, query, opts...))
}

// Deletes a service account from the project.
//
// Returns confirmation of service account deletion, or an error if the project is
// archived (archived projects have no service accounts).
func (r *AdminOrganizationProjectServiceAccountService) Delete(ctx context.Context, projectID string, serviceAccountID string, opts ...option.RequestOption) (res *AdminOrganizationProjectServiceAccountDeleteResponse, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if projectID == "" {
		err = errors.New("missing required project_id parameter")
		return nil, err
	}
	if serviceAccountID == "" {
		err = errors.New("missing required service_account_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("organization/projects/%s/service_accounts/%s", projectID, serviceAccountID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Represents an individual service account in a project.
type ProjectServiceAccount struct {
	// The identifier, which can be referenced in API endpoints
	ID string `json:"id" api:"required"`
	// The Unix timestamp (in seconds) of when the service account was created
	CreatedAt int64 `json:"created_at" api:"required" format:"unixtime"`
	// The name of the service account
	Name string `json:"name" api:"required"`
	// The object type, which is always `organization.project.service_account`
	Object constant.OrganizationProjectServiceAccount `json:"object" default:"organization.project.service_account"`
	// `owner` or `member`
	//
	// Any of "owner", "member".
	Role ProjectServiceAccountRole `json:"role" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		Name        respjson.Field
		Object      respjson.Field
		Role        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProjectServiceAccount) RawJSON() string { return r.JSON.raw }
func (r *ProjectServiceAccount) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// `owner` or `member`
type ProjectServiceAccountRole string

const (
	ProjectServiceAccountRoleOwner  ProjectServiceAccountRole = "owner"
	ProjectServiceAccountRoleMember ProjectServiceAccountRole = "member"
)

type AdminOrganizationProjectServiceAccountNewResponse struct {
	ID        string                                                  `json:"id" api:"required"`
	APIKey    AdminOrganizationProjectServiceAccountNewResponseAPIKey `json:"api_key" api:"required"`
	CreatedAt int64                                                   `json:"created_at" api:"required" format:"unixtime"`
	Name      string                                                  `json:"name" api:"required"`
	Object    constant.OrganizationProjectServiceAccount              `json:"object" default:"organization.project.service_account"`
	// Service accounts can only have one role of type `member`
	Role constant.Member `json:"role" default:"member"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		APIKey      respjson.Field
		CreatedAt   respjson.Field
		Name        respjson.Field
		Object      respjson.Field
		Role        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationProjectServiceAccountNewResponse) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationProjectServiceAccountNewResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationProjectServiceAccountNewResponseAPIKey struct {
	ID        string `json:"id" api:"required"`
	CreatedAt int64  `json:"created_at" api:"required" format:"unixtime"`
	Name      string `json:"name" api:"required"`
	// The object type, which is always `organization.project.service_account.api_key`
	Object constant.OrganizationProjectServiceAccountAPIKey `json:"object" default:"organization.project.service_account.api_key"`
	Value  string                                           `json:"value" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		Name        respjson.Field
		Object      respjson.Field
		Value       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationProjectServiceAccountNewResponseAPIKey) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationProjectServiceAccountNewResponseAPIKey) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationProjectServiceAccountDeleteResponse struct {
	ID      string                                            `json:"id" api:"required"`
	Deleted bool                                              `json:"deleted" api:"required"`
	Object  constant.OrganizationProjectServiceAccountDeleted `json:"object" default:"organization.project.service_account.deleted"`
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
func (r AdminOrganizationProjectServiceAccountDeleteResponse) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationProjectServiceAccountDeleteResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationProjectServiceAccountNewParams struct {
	// The name of the service account being created.
	Name string `json:"name" api:"required"`
	paramObj
}

func (r AdminOrganizationProjectServiceAccountNewParams) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationProjectServiceAccountNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationProjectServiceAccountNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationProjectServiceAccountListParams struct {
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

// URLQuery serializes [AdminOrganizationProjectServiceAccountListParams]'s query
// parameters as `url.Values`.
func (r AdminOrganizationProjectServiceAccountListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

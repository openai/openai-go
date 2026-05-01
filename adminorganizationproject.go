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

// AdminOrganizationProjectService contains methods and other services that help
// with interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAdminOrganizationProjectService] method instead.
type AdminOrganizationProjectService struct {
	Options         []option.RequestOption
	Users           AdminOrganizationProjectUserService
	ServiceAccounts AdminOrganizationProjectServiceAccountService
	APIKeys         AdminOrganizationProjectAPIKeyService
	RateLimits      AdminOrganizationProjectRateLimitService
	Groups          AdminOrganizationProjectGroupService
	Roles           AdminOrganizationProjectRoleService
	Certificates    AdminOrganizationProjectCertificateService
}

// NewAdminOrganizationProjectService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewAdminOrganizationProjectService(opts ...option.RequestOption) (r AdminOrganizationProjectService) {
	r = AdminOrganizationProjectService{}
	r.Options = opts
	r.Users = NewAdminOrganizationProjectUserService(opts...)
	r.ServiceAccounts = NewAdminOrganizationProjectServiceAccountService(opts...)
	r.APIKeys = NewAdminOrganizationProjectAPIKeyService(opts...)
	r.RateLimits = NewAdminOrganizationProjectRateLimitService(opts...)
	r.Groups = NewAdminOrganizationProjectGroupService(opts...)
	r.Roles = NewAdminOrganizationProjectRoleService(opts...)
	r.Certificates = NewAdminOrganizationProjectCertificateService(opts...)
	return
}

// Create a new project in the organization. Projects can be created and archived,
// but cannot be deleted.
func (r *AdminOrganizationProjectService) New(ctx context.Context, body AdminOrganizationProjectNewParams, opts ...option.RequestOption) (res *Project, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	path := "organization/projects"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Retrieves a project.
func (r *AdminOrganizationProjectService) Get(ctx context.Context, projectID string, opts ...option.RequestOption) (res *Project, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if projectID == "" {
		err = errors.New("missing required project_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("organization/projects/%s", projectID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Modifies a project in the organization.
func (r *AdminOrganizationProjectService) Update(ctx context.Context, projectID string, body AdminOrganizationProjectUpdateParams, opts ...option.RequestOption) (res *Project, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if projectID == "" {
		err = errors.New("missing required project_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("organization/projects/%s", projectID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Returns a list of projects.
func (r *AdminOrganizationProjectService) List(ctx context.Context, query AdminOrganizationProjectListParams, opts ...option.RequestOption) (res *pagination.ConversationCursorPage[Project], err error) {
	var raw *http.Response
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "organization/projects"
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

// Returns a list of projects.
func (r *AdminOrganizationProjectService) ListAutoPaging(ctx context.Context, query AdminOrganizationProjectListParams, opts ...option.RequestOption) *pagination.ConversationCursorPageAutoPager[Project] {
	return pagination.NewConversationCursorPageAutoPager(r.List(ctx, query, opts...))
}

// Archives a project in the organization. Archived projects cannot be used or
// updated.
func (r *AdminOrganizationProjectService) Archive(ctx context.Context, projectID string, opts ...option.RequestOption) (res *Project, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if projectID == "" {
		err = errors.New("missing required project_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("organization/projects/%s/archive", projectID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// Represents an individual project.
type Project struct {
	// The identifier, which can be referenced in API endpoints
	ID string `json:"id" api:"required"`
	// The Unix timestamp (in seconds) of when the project was created.
	CreatedAt int64 `json:"created_at" api:"required" format:"unixtime"`
	// The object type, which is always `organization.project`
	Object constant.OrganizationProject `json:"object" default:"organization.project"`
	// The Unix timestamp (in seconds) of when the project was archived or `null`.
	ArchivedAt int64 `json:"archived_at" api:"nullable" format:"unixtime"`
	// The external key associated with the project.
	ExternalKeyID string `json:"external_key_id" api:"nullable"`
	// The name of the project. This appears in reporting.
	Name string `json:"name" api:"nullable"`
	// `active` or `archived`
	Status string `json:"status" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID            respjson.Field
		CreatedAt     respjson.Field
		Object        respjson.Field
		ArchivedAt    respjson.Field
		ExternalKeyID respjson.Field
		Name          respjson.Field
		Status        respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r Project) RawJSON() string { return r.JSON.raw }
func (r *Project) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationProjectNewParams struct {
	// The friendly name of the project, this name appears in reports.
	Name string `json:"name" api:"required"`
	// External key ID to associate with the project.
	ExternalKeyID param.Opt[string] `json:"external_key_id,omitzero"`
	// Create the project with the specified data residency region. Your organization
	// must have access to Data residency functionality in order to use. See
	// [data residency controls](https://platform.openai.com/docs/guides/your-data#data-residency-controls)
	// to review the functionality and limitations of setting this field.
	Geography param.Opt[string] `json:"geography,omitzero"`
	paramObj
}

func (r AdminOrganizationProjectNewParams) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationProjectNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationProjectNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationProjectUpdateParams struct {
	// External key ID to associate with the project.
	ExternalKeyID param.Opt[string] `json:"external_key_id,omitzero"`
	// Geography for the project.
	Geography param.Opt[string] `json:"geography,omitzero"`
	// The updated name of the project, this name appears in reports.
	Name param.Opt[string] `json:"name,omitzero"`
	paramObj
}

func (r AdminOrganizationProjectUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationProjectUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationProjectUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationProjectListParams struct {
	// A cursor for use in pagination. `after` is an object ID that defines your place
	// in the list. For instance, if you make a list request and receive 100 objects,
	// ending with obj_foo, your subsequent call can include after=obj_foo in order to
	// fetch the next page of the list.
	After param.Opt[string] `query:"after,omitzero" json:"-"`
	// If `true` returns all projects including those that have been `archived`.
	// Archived projects are not included by default.
	IncludeArchived param.Opt[bool] `query:"include_archived,omitzero" json:"-"`
	// A limit on the number of objects to be returned. Limit can range between 1 and
	// 100, and the default is 20.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [AdminOrganizationProjectListParams]'s query parameters as
// `url.Values`.
func (r AdminOrganizationProjectListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

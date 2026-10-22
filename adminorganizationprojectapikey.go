// File generated from our OpenAPI spec by Castiron. See CONTRIBUTING.md for details.

package openai

import (
	"context"
	"errors"
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

// AdminOrganizationProjectAPIKeyService contains methods and other services that
// help with interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAdminOrganizationProjectAPIKeyService] method instead.
type AdminOrganizationProjectAPIKeyService struct {
	Options []option.RequestOption
}

// NewAdminOrganizationProjectAPIKeyService generates a new service that applies
// the given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewAdminOrganizationProjectAPIKeyService(opts ...option.RequestOption) (r AdminOrganizationProjectAPIKeyService) {
	r = AdminOrganizationProjectAPIKeyService{}
	r.Options = requestconfig.InheritedOptions(opts...)
	return
}

// Retrieves an API key in the project.
func (r *AdminOrganizationProjectAPIKeyService) Get(ctx context.Context, projectID string, apiKeyID string, opts ...option.RequestOption) (res *ProjectAPIKey, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if projectID == "" {
		err = errors.New("missing required project_id parameter")
		return nil, err
	}
	if apiKeyID == "" {
		err = errors.New("missing required api_key_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("organization/projects/%s/api_keys/%s", projectID, apiKeyID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Updates the scopes and/or expiration of an existing API key in the project,
// including a service-account-owned key. Requires permission to write project API
// keys.
//
// Provide at least one of `scopes` or `expires_in_seconds`. Omitted fields stay
// unchanged; neither field accepts `null`. Both changes are saved together.
//
// `scopes` replaces the key's scope restrictions using the same scope identifiers
// as the dashboard. An empty array means All permissions, subject to the owner's
// effective access; it does not mean no access. Use `["api.all.read"]` for
// read-only access or specific scopes such as `["api.files.read"]` for restricted
// access. This does not change user or service-account roles or grant admin
// permissions. Scope changes follow the same authentication-cache propagation as
// dashboard edits.
//
// Provide `expires_in_seconds` to set the expiration relative to this request. The
// value must satisfy the organization's and project's expiry policies. Expired
// keys cannot be updated. Retrying resets the expiration relative to the new
// request. The response contains only the redacted key, not its secret.
func (r *AdminOrganizationProjectAPIKeyService) Update(ctx context.Context, projectID string, apiKeyID string, body AdminOrganizationProjectAPIKeyUpdateParams, opts ...option.RequestOption) (res *ProjectAPIKey, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if projectID == "" {
		err = errors.New("missing required project_id parameter")
		return nil, err
	}
	if apiKeyID == "" {
		err = errors.New("missing required api_key_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("organization/projects/%s/api_keys/%s", projectID, apiKeyID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Returns a list of API keys in the project.
func (r *AdminOrganizationProjectAPIKeyService) List(ctx context.Context, projectID string, query AdminOrganizationProjectAPIKeyListParams, opts ...option.RequestOption) (res *pagination.ConversationCursorPage[ProjectAPIKey], err error) {
	var raw *http.Response
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	if projectID == "" {
		err = errors.New("missing required project_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("organization/projects/%s/api_keys", projectID)
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

// Returns a list of API keys in the project.
func (r *AdminOrganizationProjectAPIKeyService) ListAutoPaging(ctx context.Context, projectID string, query AdminOrganizationProjectAPIKeyListParams, opts ...option.RequestOption) *pagination.ConversationCursorPageAutoPager[ProjectAPIKey] {
	return pagination.NewConversationCursorPageAutoPager(r.List(ctx, projectID, query, opts...))
}

// Deletes an API key from the project, including keys owned by service accounts.
//
// Requires permission to write project API keys. Deleting a service-account key
// does not delete the service account or its other keys.
//
// Returns confirmation of the key deletion.
func (r *AdminOrganizationProjectAPIKeyService) Delete(ctx context.Context, projectID string, apiKeyID string, opts ...option.RequestOption) (res *AdminOrganizationProjectAPIKeyDeleteResponse, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if projectID == "" {
		err = errors.New("missing required project_id parameter")
		return nil, err
	}
	if apiKeyID == "" {
		err = errors.New("missing required api_key_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("organization/projects/%s/api_keys/%s", projectID, apiKeyID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Represents an individual API key in a project.
type ProjectAPIKey struct {
	// The identifier, which can be referenced in API endpoints
	ID string `json:"id" api:"required"`
	// The Unix timestamp (in seconds) of when the API key was created
	CreatedAt int64 `json:"created_at" api:"required" format:"unixtime"`
	// The Unix timestamp (in seconds) of when the API key was last used.
	LastUsedAt int64 `json:"last_used_at" api:"required" format:"unixtime"`
	// The name of the API key
	Name string `json:"name" api:"required"`
	// The object type, which is always `organization.project.api_key`
	Object constant.OrganizationProjectAPIKey `json:"object" default:"organization.project.api_key"`
	Owner  ProjectAPIKeyOwner                 `json:"owner" api:"required"`
	// Whether the API key's owner currently has effective access to the project.
	//
	// Any of "active", "inactive".
	OwnerProjectAccess ProjectAPIKeyOwnerProjectAccess `json:"owner_project_access" api:"required"`
	// The redacted value of the API key
	RedactedValue string `json:"redacted_value" api:"required"`
	// The Unix timestamp (in seconds) when the API key expires, or null if it does not
	// expire.
	ExpiresAt int64 `json:"expires_at" api:"nullable" format:"unixtime"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                 respjson.Field
		CreatedAt          respjson.Field
		LastUsedAt         respjson.Field
		Name               respjson.Field
		Object             respjson.Field
		Owner              respjson.Field
		OwnerProjectAccess respjson.Field
		RedactedValue      respjson.Field
		ExpiresAt          respjson.Field
		ExtraFields        map[string]respjson.Field
		raw                string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProjectAPIKey) RawJSON() string { return r.JSON.raw }
func (r *ProjectAPIKey) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ProjectAPIKeyOwner struct {
	// The service account that owns a project API key.
	ServiceAccount ProjectAPIKeyOwnerServiceAccount `json:"service_account"`
	// `user` or `service_account`
	//
	// Any of "user", "service_account".
	Type string `json:"type"`
	// The user that owns a project API key.
	User ProjectAPIKeyOwnerUser `json:"user"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ServiceAccount respjson.Field
		Type           respjson.Field
		User           respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProjectAPIKeyOwner) RawJSON() string { return r.JSON.raw }
func (r *ProjectAPIKeyOwner) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The service account that owns a project API key.
type ProjectAPIKeyOwnerServiceAccount struct {
	// The identifier, which can be referenced in API endpoints
	ID string `json:"id" api:"required"`
	// The Unix timestamp (in seconds) of when the service account was created.
	CreatedAt int64 `json:"created_at" api:"required" format:"unixtime"`
	// The name of the service account.
	Name string `json:"name" api:"required"`
	// The service account's project role.
	Role string `json:"role" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		Name        respjson.Field
		Role        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProjectAPIKeyOwnerServiceAccount) RawJSON() string { return r.JSON.raw }
func (r *ProjectAPIKeyOwnerServiceAccount) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The user that owns a project API key.
type ProjectAPIKeyOwnerUser struct {
	// The identifier, which can be referenced in API endpoints
	ID string `json:"id" api:"required"`
	// The Unix timestamp (in seconds) of when the user was created.
	CreatedAt int64 `json:"created_at" api:"required" format:"unixtime"`
	// The email address of the user.
	Email string `json:"email" api:"required"`
	// The name of the user.
	Name string `json:"name" api:"required"`
	// The user's project role.
	Role string `json:"role" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		Email       respjson.Field
		Name        respjson.Field
		Role        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProjectAPIKeyOwnerUser) RawJSON() string { return r.JSON.raw }
func (r *ProjectAPIKeyOwnerUser) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Whether the API key's owner currently has effective access to the project.
type ProjectAPIKeyOwnerProjectAccess string

const (
	ProjectAPIKeyOwnerProjectAccessActive   ProjectAPIKeyOwnerProjectAccess = "active"
	ProjectAPIKeyOwnerProjectAccessInactive ProjectAPIKeyOwnerProjectAccess = "inactive"
)

type AdminOrganizationProjectAPIKeyDeleteResponse struct {
	ID      string                                    `json:"id" api:"required"`
	Deleted bool                                      `json:"deleted" api:"required"`
	Object  constant.OrganizationProjectAPIKeyDeleted `json:"object" default:"organization.project.api_key.deleted"`
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
func (r AdminOrganizationProjectAPIKeyDeleteResponse) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationProjectAPIKeyDeleteResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationProjectAPIKeyUpdateParams struct {
	// Number of seconds from this request until the API key expires. Must comply with
	// the organization and project expiry policies. Expired keys cannot be updated.
	// Retrying resets the expiry relative to the new request. Omit to keep the current
	// expiration.
	ExpiresInSeconds int64 `json:"expires_in_seconds,omitzero"`
	// Replace the API key scope restrictions. Omit to leave scopes unchanged. An empty
	// array removes key restrictions (All permissions), subject to the owner's
	// effective permissions. Use ["api.all.read"] for read-only access. Scopes do not
	// change the owner's roles or grant administrative permissions.
	Scopes []string `json:"scopes,omitzero"`
	paramObj
}

func (r AdminOrganizationProjectAPIKeyUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationProjectAPIKeyUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationProjectAPIKeyUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationProjectAPIKeyListParams struct {
	// A cursor for use in pagination. `after` is an object ID that defines your place
	// in the list. For instance, if you make a list request and receive 100 objects,
	// ending with obj_foo, your subsequent call can include after=obj_foo in order to
	// fetch the next page of the list.
	After param.Opt[string] `query:"after,omitzero" json:"-"`
	// A limit on the number of objects to be returned. Limit can range between 1 and
	// 100, and the default is 20.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Filter API keys by whether the owner currently has effective access to the
	// project. Use `active` for owners with access, `inactive` for owners without
	// access, or `any` for all enabled project API keys. If omitted, the endpoint
	// applies its existing membership-based visibility rules, which may exclude some
	// enabled keys.
	//
	// Any of "active", "inactive", "any".
	OwnerProjectAccess AdminOrganizationProjectAPIKeyListParamsOwnerProjectAccess `query:"owner_project_access,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [AdminOrganizationProjectAPIKeyListParams]'s query
// parameters as `url.Values`.
func (r AdminOrganizationProjectAPIKeyListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Filter API keys by whether the owner currently has effective access to the
// project. Use `active` for owners with access, `inactive` for owners without
// access, or `any` for all enabled project API keys. If omitted, the endpoint
// applies its existing membership-based visibility rules, which may exclude some
// enabled keys.
type AdminOrganizationProjectAPIKeyListParamsOwnerProjectAccess string

const (
	AdminOrganizationProjectAPIKeyListParamsOwnerProjectAccessActive   AdminOrganizationProjectAPIKeyListParamsOwnerProjectAccess = "active"
	AdminOrganizationProjectAPIKeyListParamsOwnerProjectAccessInactive AdminOrganizationProjectAPIKeyListParamsOwnerProjectAccess = "inactive"
	AdminOrganizationProjectAPIKeyListParamsOwnerProjectAccessAny      AdminOrganizationProjectAPIKeyListParamsOwnerProjectAccess = "any"
)

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
	r.Options = opts
	return
}

// Retrieves an API key in the project.
func (r *AdminOrganizationProjectAPIKeyService) Get(ctx context.Context, projectID string, keyID string, opts ...option.RequestOption) (res *ProjectAPIKey, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if projectID == "" {
		err = errors.New("missing required project_id parameter")
		return nil, err
	}
	if keyID == "" {
		err = errors.New("missing required key_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("organization/projects/%s/api_keys/%s", projectID, keyID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
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
	path := fmt.Sprintf("organization/projects/%s/api_keys", projectID)
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

// Deletes an API key from the project.
//
// Returns confirmation of the key deletion, or an error if the key belonged to a
// service account.
func (r *AdminOrganizationProjectAPIKeyService) Delete(ctx context.Context, projectID string, keyID string, opts ...option.RequestOption) (res *AdminOrganizationProjectAPIKeyDeleteResponse, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if projectID == "" {
		err = errors.New("missing required project_id parameter")
		return nil, err
	}
	if keyID == "" {
		err = errors.New("missing required key_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("organization/projects/%s/api_keys/%s", projectID, keyID)
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
	// The redacted value of the API key
	RedactedValue string `json:"redacted_value" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID            respjson.Field
		CreatedAt     respjson.Field
		LastUsedAt    respjson.Field
		Name          respjson.Field
		Object        respjson.Field
		Owner         respjson.Field
		RedactedValue respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProjectAPIKey) RawJSON() string { return r.JSON.raw }
func (r *ProjectAPIKey) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ProjectAPIKeyOwner struct {
	// Represents an individual service account in a project.
	ServiceAccount ProjectServiceAccount `json:"service_account"`
	// `user` or `service_account`
	//
	// Any of "user", "service_account".
	Type string `json:"type"`
	// Represents an individual user in a project.
	User ProjectUser `json:"user"`
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

type AdminOrganizationProjectAPIKeyListParams struct {
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

// URLQuery serializes [AdminOrganizationProjectAPIKeyListParams]'s query
// parameters as `url.Values`.
func (r AdminOrganizationProjectAPIKeyListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

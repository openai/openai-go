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

// AdminOrganizationProjectRateLimitService contains methods and other services
// that help with interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAdminOrganizationProjectRateLimitService] method instead.
type AdminOrganizationProjectRateLimitService struct {
	Options []option.RequestOption
}

// NewAdminOrganizationProjectRateLimitService generates a new service that applies
// the given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewAdminOrganizationProjectRateLimitService(opts ...option.RequestOption) (r AdminOrganizationProjectRateLimitService) {
	r = AdminOrganizationProjectRateLimitService{}
	r.Options = opts
	return
}

// Returns the rate limits per model for a project.
func (r *AdminOrganizationProjectRateLimitService) ListRateLimits(ctx context.Context, projectID string, query AdminOrganizationProjectRateLimitListRateLimitsParams, opts ...option.RequestOption) (res *pagination.ConversationCursorPage[ProjectRateLimit], err error) {
	var raw *http.Response
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	if projectID == "" {
		err = errors.New("missing required project_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("organization/projects/%s/rate_limits", projectID)
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

// Returns the rate limits per model for a project.
func (r *AdminOrganizationProjectRateLimitService) ListRateLimitsAutoPaging(ctx context.Context, projectID string, query AdminOrganizationProjectRateLimitListRateLimitsParams, opts ...option.RequestOption) *pagination.ConversationCursorPageAutoPager[ProjectRateLimit] {
	return pagination.NewConversationCursorPageAutoPager(r.ListRateLimits(ctx, projectID, query, opts...))
}

// Updates a project rate limit.
func (r *AdminOrganizationProjectRateLimitService) UpdateRateLimit(ctx context.Context, projectID string, rateLimitID string, body AdminOrganizationProjectRateLimitUpdateRateLimitParams, opts ...option.RequestOption) (res *ProjectRateLimit, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if projectID == "" {
		err = errors.New("missing required project_id parameter")
		return nil, err
	}
	if rateLimitID == "" {
		err = errors.New("missing required rate_limit_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("organization/projects/%s/rate_limits/%s", projectID, rateLimitID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Represents a project rate limit config.
type ProjectRateLimit struct {
	// The identifier, which can be referenced in API endpoints.
	ID string `json:"id" api:"required"`
	// The maximum requests per minute.
	MaxRequestsPer1Minute int64 `json:"max_requests_per_1_minute" api:"required"`
	// The maximum tokens per minute.
	MaxTokensPer1Minute int64 `json:"max_tokens_per_1_minute" api:"required"`
	// The model this rate limit applies to.
	Model string `json:"model" api:"required"`
	// The object type, which is always `project.rate_limit`
	Object constant.ProjectRateLimit `json:"object" default:"project.rate_limit"`
	// The maximum batch input tokens per day. Only present for relevant models.
	Batch1DayMaxInputTokens int64 `json:"batch_1_day_max_input_tokens"`
	// The maximum audio megabytes per minute. Only present for relevant models.
	MaxAudioMegabytesPer1Minute int64 `json:"max_audio_megabytes_per_1_minute"`
	// The maximum images per minute. Only present for relevant models.
	MaxImagesPer1Minute int64 `json:"max_images_per_1_minute"`
	// The maximum requests per day. Only present for relevant models.
	MaxRequestsPer1Day int64 `json:"max_requests_per_1_day"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                          respjson.Field
		MaxRequestsPer1Minute       respjson.Field
		MaxTokensPer1Minute         respjson.Field
		Model                       respjson.Field
		Object                      respjson.Field
		Batch1DayMaxInputTokens     respjson.Field
		MaxAudioMegabytesPer1Minute respjson.Field
		MaxImagesPer1Minute         respjson.Field
		MaxRequestsPer1Day          respjson.Field
		ExtraFields                 map[string]respjson.Field
		raw                         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProjectRateLimit) RawJSON() string { return r.JSON.raw }
func (r *ProjectRateLimit) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationProjectRateLimitListRateLimitsParams struct {
	// A cursor for use in pagination. `after` is an object ID that defines your place
	// in the list. For instance, if you make a list request and receive 100 objects,
	// ending with obj_foo, your subsequent call can include after=obj_foo in order to
	// fetch the next page of the list.
	After param.Opt[string] `query:"after,omitzero" json:"-"`
	// A cursor for use in pagination. `before` is an object ID that defines your place
	// in the list. For instance, if you make a list request and receive 100 objects,
	// beginning with obj_foo, your subsequent call can include before=obj_foo in order
	// to fetch the previous page of the list.
	Before param.Opt[string] `query:"before,omitzero" json:"-"`
	// A limit on the number of objects to be returned. The default is 100.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [AdminOrganizationProjectRateLimitListRateLimitsParams]'s
// query parameters as `url.Values`.
func (r AdminOrganizationProjectRateLimitListRateLimitsParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type AdminOrganizationProjectRateLimitUpdateRateLimitParams struct {
	// The maximum batch input tokens per day. Only relevant for certain models.
	Batch1DayMaxInputTokens param.Opt[int64] `json:"batch_1_day_max_input_tokens,omitzero"`
	// The maximum audio megabytes per minute. Only relevant for certain models.
	MaxAudioMegabytesPer1Minute param.Opt[int64] `json:"max_audio_megabytes_per_1_minute,omitzero"`
	// The maximum images per minute. Only relevant for certain models.
	MaxImagesPer1Minute param.Opt[int64] `json:"max_images_per_1_minute,omitzero"`
	// The maximum requests per day. Only relevant for certain models.
	MaxRequestsPer1Day param.Opt[int64] `json:"max_requests_per_1_day,omitzero"`
	// The maximum requests per minute.
	MaxRequestsPer1Minute param.Opt[int64] `json:"max_requests_per_1_minute,omitzero"`
	// The maximum tokens per minute.
	MaxTokensPer1Minute param.Opt[int64] `json:"max_tokens_per_1_minute,omitzero"`
	paramObj
}

func (r AdminOrganizationProjectRateLimitUpdateRateLimitParams) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationProjectRateLimitUpdateRateLimitParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationProjectRateLimitUpdateRateLimitParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

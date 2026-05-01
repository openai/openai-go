// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"slices"

	"github.com/openai/openai-go/v3/internal/apijson"
	"github.com/openai/openai-go/v3/internal/apiquery"
	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/packages/respjson"
	"github.com/openai/openai-go/v3/shared/constant"
)

// AdminOrganizationUsageService contains methods and other services that help with
// interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAdminOrganizationUsageService] method instead.
type AdminOrganizationUsageService struct {
	Options []option.RequestOption
}

// NewAdminOrganizationUsageService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewAdminOrganizationUsageService(opts ...option.RequestOption) (r AdminOrganizationUsageService) {
	r = AdminOrganizationUsageService{}
	r.Options = opts
	return
}

// Get audio speeches usage details for the organization.
func (r *AdminOrganizationUsageService) AudioSpeeches(ctx context.Context, query AdminOrganizationUsageAudioSpeechesParams, opts ...option.RequestOption) (res *AdminOrganizationUsageAudioSpeechesResponse, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	path := "organization/usage/audio_speeches"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

// Get audio transcriptions usage details for the organization.
func (r *AdminOrganizationUsageService) AudioTranscriptions(ctx context.Context, query AdminOrganizationUsageAudioTranscriptionsParams, opts ...option.RequestOption) (res *AdminOrganizationUsageAudioTranscriptionsResponse, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	path := "organization/usage/audio_transcriptions"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

// Get code interpreter sessions usage details for the organization.
func (r *AdminOrganizationUsageService) CodeInterpreterSessions(ctx context.Context, query AdminOrganizationUsageCodeInterpreterSessionsParams, opts ...option.RequestOption) (res *AdminOrganizationUsageCodeInterpreterSessionsResponse, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	path := "organization/usage/code_interpreter_sessions"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

// Get completions usage details for the organization.
func (r *AdminOrganizationUsageService) Completions(ctx context.Context, query AdminOrganizationUsageCompletionsParams, opts ...option.RequestOption) (res *AdminOrganizationUsageCompletionsResponse, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	path := "organization/usage/completions"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

// Get costs details for the organization.
func (r *AdminOrganizationUsageService) Costs(ctx context.Context, query AdminOrganizationUsageCostsParams, opts ...option.RequestOption) (res *AdminOrganizationUsageCostsResponse, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	path := "organization/costs"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

// Get embeddings usage details for the organization.
func (r *AdminOrganizationUsageService) Embeddings(ctx context.Context, query AdminOrganizationUsageEmbeddingsParams, opts ...option.RequestOption) (res *AdminOrganizationUsageEmbeddingsResponse, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	path := "organization/usage/embeddings"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

// Get images usage details for the organization.
func (r *AdminOrganizationUsageService) Images(ctx context.Context, query AdminOrganizationUsageImagesParams, opts ...option.RequestOption) (res *AdminOrganizationUsageImagesResponse, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	path := "organization/usage/images"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

// Get moderations usage details for the organization.
func (r *AdminOrganizationUsageService) Moderations(ctx context.Context, query AdminOrganizationUsageModerationsParams, opts ...option.RequestOption) (res *AdminOrganizationUsageModerationsResponse, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	path := "organization/usage/moderations"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

// Get vector stores usage details for the organization.
func (r *AdminOrganizationUsageService) VectorStores(ctx context.Context, query AdminOrganizationUsageVectorStoresParams, opts ...option.RequestOption) (res *AdminOrganizationUsageVectorStoresResponse, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	path := "organization/usage/vector_stores"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

type AdminOrganizationUsageAudioSpeechesResponse struct {
	Data     []AdminOrganizationUsageAudioSpeechesResponseData `json:"data" api:"required"`
	HasMore  bool                                              `json:"has_more" api:"required"`
	NextPage string                                            `json:"next_page" api:"required"`
	Object   constant.Page                                     `json:"object" default:"page"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Data        respjson.Field
		HasMore     respjson.Field
		NextPage    respjson.Field
		Object      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageAudioSpeechesResponse) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationUsageAudioSpeechesResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationUsageAudioSpeechesResponseData struct {
	EndTime   int64                                                        `json:"end_time" api:"required"`
	Object    constant.Bucket                                              `json:"object" default:"bucket"`
	Results   []AdminOrganizationUsageAudioSpeechesResponseDataResultUnion `json:"results" api:"required"`
	StartTime int64                                                        `json:"start_time" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EndTime     respjson.Field
		Object      respjson.Field
		Results     respjson.Field
		StartTime   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageAudioSpeechesResponseData) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationUsageAudioSpeechesResponseData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// AdminOrganizationUsageAudioSpeechesResponseDataResultUnion contains all possible
// properties and values from
// [AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageCompletionsResult],
// [AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageEmbeddingsResult],
// [AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageModerationsResult],
// [AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageImagesResult],
// [AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageAudioSpeechesResult],
// [AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageAudioTranscriptionsResult],
// [AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageVectorStoresResult],
// [AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageCodeInterpreterSessionsResult],
// [AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationCostsResult].
//
// Use the [AdminOrganizationUsageAudioSpeechesResponseDataResultUnion.AsAny]
// method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type AdminOrganizationUsageAudioSpeechesResponseDataResultUnion struct {
	InputTokens      int64 `json:"input_tokens"`
	NumModelRequests int64 `json:"num_model_requests"`
	// Any of "organization.usage.completions.result",
	// "organization.usage.embeddings.result", "organization.usage.moderations.result",
	// "organization.usage.images.result", "organization.usage.audio_speeches.result",
	// "organization.usage.audio_transcriptions.result",
	// "organization.usage.vector_stores.result",
	// "organization.usage.code_interpreter_sessions.result",
	// "organization.costs.result".
	Object string `json:"object"`
	// This field is from variant
	// [AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageCompletionsResult].
	OutputTokens int64  `json:"output_tokens"`
	APIKeyID     string `json:"api_key_id"`
	// This field is from variant
	// [AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageCompletionsResult].
	Batch bool `json:"batch"`
	// This field is from variant
	// [AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageCompletionsResult].
	InputAudioTokens int64 `json:"input_audio_tokens"`
	// This field is from variant
	// [AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageCompletionsResult].
	InputCachedTokens int64  `json:"input_cached_tokens"`
	Model             string `json:"model"`
	// This field is from variant
	// [AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageCompletionsResult].
	OutputAudioTokens int64  `json:"output_audio_tokens"`
	ProjectID         string `json:"project_id"`
	// This field is from variant
	// [AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageCompletionsResult].
	ServiceTier string `json:"service_tier"`
	UserID      string `json:"user_id"`
	// This field is from variant
	// [AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageImagesResult].
	Images int64 `json:"images"`
	// This field is from variant
	// [AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageImagesResult].
	Size string `json:"size"`
	// This field is from variant
	// [AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageImagesResult].
	Source string `json:"source"`
	// This field is from variant
	// [AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageAudioSpeechesResult].
	Characters int64 `json:"characters"`
	// This field is from variant
	// [AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageAudioTranscriptionsResult].
	Seconds int64 `json:"seconds"`
	// This field is from variant
	// [AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageVectorStoresResult].
	UsageBytes int64 `json:"usage_bytes"`
	// This field is from variant
	// [AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageCodeInterpreterSessionsResult].
	NumSessions int64 `json:"num_sessions"`
	// This field is from variant
	// [AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationCostsResult].
	Amount AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationCostsResultAmount `json:"amount"`
	// This field is from variant
	// [AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationCostsResult].
	LineItem string `json:"line_item"`
	JSON     struct {
		InputTokens       respjson.Field
		NumModelRequests  respjson.Field
		Object            respjson.Field
		OutputTokens      respjson.Field
		APIKeyID          respjson.Field
		Batch             respjson.Field
		InputAudioTokens  respjson.Field
		InputCachedTokens respjson.Field
		Model             respjson.Field
		OutputAudioTokens respjson.Field
		ProjectID         respjson.Field
		ServiceTier       respjson.Field
		UserID            respjson.Field
		Images            respjson.Field
		Size              respjson.Field
		Source            respjson.Field
		Characters        respjson.Field
		Seconds           respjson.Field
		UsageBytes        respjson.Field
		NumSessions       respjson.Field
		Amount            respjson.Field
		LineItem          respjson.Field
		raw               string
	} `json:"-"`
}

// anyAdminOrganizationUsageAudioSpeechesResponseDataResult is implemented by each
// variant of [AdminOrganizationUsageAudioSpeechesResponseDataResultUnion] to add
// type safety for the return type of
// [AdminOrganizationUsageAudioSpeechesResponseDataResultUnion.AsAny]
type anyAdminOrganizationUsageAudioSpeechesResponseDataResult interface {
	implAdminOrganizationUsageAudioSpeechesResponseDataResultUnion()
}

func (AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageCompletionsResult) implAdminOrganizationUsageAudioSpeechesResponseDataResultUnion() {
}
func (AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageEmbeddingsResult) implAdminOrganizationUsageAudioSpeechesResponseDataResultUnion() {
}
func (AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageModerationsResult) implAdminOrganizationUsageAudioSpeechesResponseDataResultUnion() {
}
func (AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageImagesResult) implAdminOrganizationUsageAudioSpeechesResponseDataResultUnion() {
}
func (AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageAudioSpeechesResult) implAdminOrganizationUsageAudioSpeechesResponseDataResultUnion() {
}
func (AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageAudioTranscriptionsResult) implAdminOrganizationUsageAudioSpeechesResponseDataResultUnion() {
}
func (AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageVectorStoresResult) implAdminOrganizationUsageAudioSpeechesResponseDataResultUnion() {
}
func (AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageCodeInterpreterSessionsResult) implAdminOrganizationUsageAudioSpeechesResponseDataResultUnion() {
}
func (AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationCostsResult) implAdminOrganizationUsageAudioSpeechesResponseDataResultUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := AdminOrganizationUsageAudioSpeechesResponseDataResultUnion.AsAny().(type) {
//	case openai.AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageCompletionsResult:
//	case openai.AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageEmbeddingsResult:
//	case openai.AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageModerationsResult:
//	case openai.AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageImagesResult:
//	case openai.AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageAudioSpeechesResult:
//	case openai.AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageAudioTranscriptionsResult:
//	case openai.AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageVectorStoresResult:
//	case openai.AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageCodeInterpreterSessionsResult:
//	case openai.AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationCostsResult:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u AdminOrganizationUsageAudioSpeechesResponseDataResultUnion) AsAny() anyAdminOrganizationUsageAudioSpeechesResponseDataResult {
	switch u.Object {
	case "organization.usage.completions.result":
		return u.AsOrganizationUsageCompletionsResult()
	case "organization.usage.embeddings.result":
		return u.AsOrganizationUsageEmbeddingsResult()
	case "organization.usage.moderations.result":
		return u.AsOrganizationUsageModerationsResult()
	case "organization.usage.images.result":
		return u.AsOrganizationUsageImagesResult()
	case "organization.usage.audio_speeches.result":
		return u.AsOrganizationUsageAudioSpeechesResult()
	case "organization.usage.audio_transcriptions.result":
		return u.AsOrganizationUsageAudioTranscriptionsResult()
	case "organization.usage.vector_stores.result":
		return u.AsOrganizationUsageVectorStoresResult()
	case "organization.usage.code_interpreter_sessions.result":
		return u.AsOrganizationUsageCodeInterpreterSessionsResult()
	case "organization.costs.result":
		return u.AsOrganizationCostsResult()
	}
	return nil
}

func (u AdminOrganizationUsageAudioSpeechesResponseDataResultUnion) AsOrganizationUsageCompletionsResult() (v AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageCompletionsResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageAudioSpeechesResponseDataResultUnion) AsOrganizationUsageEmbeddingsResult() (v AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageEmbeddingsResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageAudioSpeechesResponseDataResultUnion) AsOrganizationUsageModerationsResult() (v AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageModerationsResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageAudioSpeechesResponseDataResultUnion) AsOrganizationUsageImagesResult() (v AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageImagesResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageAudioSpeechesResponseDataResultUnion) AsOrganizationUsageAudioSpeechesResult() (v AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageAudioSpeechesResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageAudioSpeechesResponseDataResultUnion) AsOrganizationUsageAudioTranscriptionsResult() (v AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageAudioTranscriptionsResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageAudioSpeechesResponseDataResultUnion) AsOrganizationUsageVectorStoresResult() (v AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageVectorStoresResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageAudioSpeechesResponseDataResultUnion) AsOrganizationUsageCodeInterpreterSessionsResult() (v AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageCodeInterpreterSessionsResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageAudioSpeechesResponseDataResultUnion) AsOrganizationCostsResult() (v AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationCostsResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u AdminOrganizationUsageAudioSpeechesResponseDataResultUnion) RawJSON() string {
	return u.JSON.raw
}

func (r *AdminOrganizationUsageAudioSpeechesResponseDataResultUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated completions usage details of the specific time bucket.
type AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageCompletionsResult struct {
	// The aggregated number of text input tokens used, including cached tokens. For
	// customers subscribe to scale tier, this includes scale tier tokens.
	InputTokens int64 `json:"input_tokens" api:"required"`
	// The count of requests made to the model.
	NumModelRequests int64                                       `json:"num_model_requests" api:"required"`
	Object           constant.OrganizationUsageCompletionsResult `json:"object" default:"organization.usage.completions.result"`
	// The aggregated number of text output tokens used. For customers subscribe to
	// scale tier, this includes scale tier tokens.
	OutputTokens int64 `json:"output_tokens" api:"required"`
	// When `group_by=api_key_id`, this field provides the API key ID of the grouped
	// usage result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=batch`, this field tells whether the grouped usage result is
	// batch or not.
	Batch bool `json:"batch" api:"nullable"`
	// The aggregated number of audio input tokens used, including cached tokens.
	InputAudioTokens int64 `json:"input_audio_tokens"`
	// The aggregated number of text input tokens that has been cached from previous
	// requests. For customers subscribe to scale tier, this includes scale tier
	// tokens.
	InputCachedTokens int64 `json:"input_cached_tokens"`
	// When `group_by=model`, this field provides the model name of the grouped usage
	// result.
	Model string `json:"model" api:"nullable"`
	// The aggregated number of audio output tokens used.
	OutputAudioTokens int64 `json:"output_audio_tokens"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// When `group_by=service_tier`, this field provides the service tier of the
	// grouped usage result.
	ServiceTier string `json:"service_tier" api:"nullable"`
	// When `group_by=user_id`, this field provides the user ID of the grouped usage
	// result.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		InputTokens       respjson.Field
		NumModelRequests  respjson.Field
		Object            respjson.Field
		OutputTokens      respjson.Field
		APIKeyID          respjson.Field
		Batch             respjson.Field
		InputAudioTokens  respjson.Field
		InputCachedTokens respjson.Field
		Model             respjson.Field
		OutputAudioTokens respjson.Field
		ProjectID         respjson.Field
		ServiceTier       respjson.Field
		UserID            respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageCompletionsResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageCompletionsResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated embeddings usage details of the specific time bucket.
type AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageEmbeddingsResult struct {
	// The aggregated number of input tokens used.
	InputTokens int64 `json:"input_tokens" api:"required"`
	// The count of requests made to the model.
	NumModelRequests int64                                      `json:"num_model_requests" api:"required"`
	Object           constant.OrganizationUsageEmbeddingsResult `json:"object" default:"organization.usage.embeddings.result"`
	// When `group_by=api_key_id`, this field provides the API key ID of the grouped
	// usage result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=model`, this field provides the model name of the grouped usage
	// result.
	Model string `json:"model" api:"nullable"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// When `group_by=user_id`, this field provides the user ID of the grouped usage
	// result.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		InputTokens      respjson.Field
		NumModelRequests respjson.Field
		Object           respjson.Field
		APIKeyID         respjson.Field
		Model            respjson.Field
		ProjectID        respjson.Field
		UserID           respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageEmbeddingsResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageEmbeddingsResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated moderations usage details of the specific time bucket.
type AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageModerationsResult struct {
	// The aggregated number of input tokens used.
	InputTokens int64 `json:"input_tokens" api:"required"`
	// The count of requests made to the model.
	NumModelRequests int64                                       `json:"num_model_requests" api:"required"`
	Object           constant.OrganizationUsageModerationsResult `json:"object" default:"organization.usage.moderations.result"`
	// When `group_by=api_key_id`, this field provides the API key ID of the grouped
	// usage result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=model`, this field provides the model name of the grouped usage
	// result.
	Model string `json:"model" api:"nullable"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// When `group_by=user_id`, this field provides the user ID of the grouped usage
	// result.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		InputTokens      respjson.Field
		NumModelRequests respjson.Field
		Object           respjson.Field
		APIKeyID         respjson.Field
		Model            respjson.Field
		ProjectID        respjson.Field
		UserID           respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageModerationsResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageModerationsResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated images usage details of the specific time bucket.
type AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageImagesResult struct {
	// The number of images processed.
	Images int64 `json:"images" api:"required"`
	// The count of requests made to the model.
	NumModelRequests int64                                  `json:"num_model_requests" api:"required"`
	Object           constant.OrganizationUsageImagesResult `json:"object" default:"organization.usage.images.result"`
	// When `group_by=api_key_id`, this field provides the API key ID of the grouped
	// usage result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=model`, this field provides the model name of the grouped usage
	// result.
	Model string `json:"model" api:"nullable"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// When `group_by=size`, this field provides the image size of the grouped usage
	// result.
	Size string `json:"size" api:"nullable"`
	// When `group_by=source`, this field provides the source of the grouped usage
	// result, possible values are `image.generation`, `image.edit`, `image.variation`.
	Source string `json:"source" api:"nullable"`
	// When `group_by=user_id`, this field provides the user ID of the grouped usage
	// result.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Images           respjson.Field
		NumModelRequests respjson.Field
		Object           respjson.Field
		APIKeyID         respjson.Field
		Model            respjson.Field
		ProjectID        respjson.Field
		Size             respjson.Field
		Source           respjson.Field
		UserID           respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageImagesResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageImagesResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated audio speeches usage details of the specific time bucket.
type AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageAudioSpeechesResult struct {
	// The number of characters processed.
	Characters int64 `json:"characters" api:"required"`
	// The count of requests made to the model.
	NumModelRequests int64                                         `json:"num_model_requests" api:"required"`
	Object           constant.OrganizationUsageAudioSpeechesResult `json:"object" default:"organization.usage.audio_speeches.result"`
	// When `group_by=api_key_id`, this field provides the API key ID of the grouped
	// usage result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=model`, this field provides the model name of the grouped usage
	// result.
	Model string `json:"model" api:"nullable"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// When `group_by=user_id`, this field provides the user ID of the grouped usage
	// result.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Characters       respjson.Field
		NumModelRequests respjson.Field
		Object           respjson.Field
		APIKeyID         respjson.Field
		Model            respjson.Field
		ProjectID        respjson.Field
		UserID           respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageAudioSpeechesResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageAudioSpeechesResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated audio transcriptions usage details of the specific time bucket.
type AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageAudioTranscriptionsResult struct {
	// The count of requests made to the model.
	NumModelRequests int64                                               `json:"num_model_requests" api:"required"`
	Object           constant.OrganizationUsageAudioTranscriptionsResult `json:"object" default:"organization.usage.audio_transcriptions.result"`
	// The number of seconds processed.
	Seconds int64 `json:"seconds" api:"required"`
	// When `group_by=api_key_id`, this field provides the API key ID of the grouped
	// usage result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=model`, this field provides the model name of the grouped usage
	// result.
	Model string `json:"model" api:"nullable"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// When `group_by=user_id`, this field provides the user ID of the grouped usage
	// result.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		NumModelRequests respjson.Field
		Object           respjson.Field
		Seconds          respjson.Field
		APIKeyID         respjson.Field
		Model            respjson.Field
		ProjectID        respjson.Field
		UserID           respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageAudioTranscriptionsResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageAudioTranscriptionsResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated vector stores usage details of the specific time bucket.
type AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageVectorStoresResult struct {
	Object constant.OrganizationUsageVectorStoresResult `json:"object" default:"organization.usage.vector_stores.result"`
	// The vector stores usage in bytes.
	UsageBytes int64 `json:"usage_bytes" api:"required"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Object      respjson.Field
		UsageBytes  respjson.Field
		ProjectID   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageVectorStoresResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageVectorStoresResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated code interpreter sessions usage details of the specific time
// bucket.
type AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageCodeInterpreterSessionsResult struct {
	// The number of code interpreter sessions.
	NumSessions int64                                                   `json:"num_sessions" api:"required"`
	Object      constant.OrganizationUsageCodeInterpreterSessionsResult `json:"object" default:"organization.usage.code_interpreter_sessions.result"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		NumSessions respjson.Field
		Object      respjson.Field
		ProjectID   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageCodeInterpreterSessionsResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationUsageCodeInterpreterSessionsResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated costs details of the specific time bucket.
type AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationCostsResult struct {
	Object constant.OrganizationCostsResult `json:"object" default:"organization.costs.result"`
	// The monetary value in its associated currency.
	Amount AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationCostsResultAmount `json:"amount"`
	// When `group_by=api_key_id`, this field provides the API Key ID of the grouped
	// costs result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=line_item`, this field provides the line item of the grouped
	// costs result.
	LineItem string `json:"line_item" api:"nullable"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// costs result.
	ProjectID string `json:"project_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Object      respjson.Field
		Amount      respjson.Field
		APIKeyID    respjson.Field
		LineItem    respjson.Field
		ProjectID   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationCostsResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationCostsResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The monetary value in its associated currency.
type AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationCostsResultAmount struct {
	// Lowercase ISO-4217 currency e.g. "usd"
	Currency string `json:"currency"`
	// The numeric value of the cost.
	Value float64 `json:"value"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Currency    respjson.Field
		Value       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationCostsResultAmount) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageAudioSpeechesResponseDataResultOrganizationCostsResultAmount) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationUsageAudioTranscriptionsResponse struct {
	Data     []AdminOrganizationUsageAudioTranscriptionsResponseData `json:"data" api:"required"`
	HasMore  bool                                                    `json:"has_more" api:"required"`
	NextPage string                                                  `json:"next_page" api:"required"`
	Object   constant.Page                                           `json:"object" default:"page"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Data        respjson.Field
		HasMore     respjson.Field
		NextPage    respjson.Field
		Object      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageAudioTranscriptionsResponse) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationUsageAudioTranscriptionsResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationUsageAudioTranscriptionsResponseData struct {
	EndTime   int64                                                              `json:"end_time" api:"required"`
	Object    constant.Bucket                                                    `json:"object" default:"bucket"`
	Results   []AdminOrganizationUsageAudioTranscriptionsResponseDataResultUnion `json:"results" api:"required"`
	StartTime int64                                                              `json:"start_time" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EndTime     respjson.Field
		Object      respjson.Field
		Results     respjson.Field
		StartTime   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageAudioTranscriptionsResponseData) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationUsageAudioTranscriptionsResponseData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// AdminOrganizationUsageAudioTranscriptionsResponseDataResultUnion contains all
// possible properties and values from
// [AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageCompletionsResult],
// [AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageEmbeddingsResult],
// [AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageModerationsResult],
// [AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageImagesResult],
// [AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageAudioSpeechesResult],
// [AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageAudioTranscriptionsResult],
// [AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageVectorStoresResult],
// [AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageCodeInterpreterSessionsResult],
// [AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationCostsResult].
//
// Use the [AdminOrganizationUsageAudioTranscriptionsResponseDataResultUnion.AsAny]
// method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type AdminOrganizationUsageAudioTranscriptionsResponseDataResultUnion struct {
	InputTokens      int64 `json:"input_tokens"`
	NumModelRequests int64 `json:"num_model_requests"`
	// Any of "organization.usage.completions.result",
	// "organization.usage.embeddings.result", "organization.usage.moderations.result",
	// "organization.usage.images.result", "organization.usage.audio_speeches.result",
	// "organization.usage.audio_transcriptions.result",
	// "organization.usage.vector_stores.result",
	// "organization.usage.code_interpreter_sessions.result",
	// "organization.costs.result".
	Object string `json:"object"`
	// This field is from variant
	// [AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageCompletionsResult].
	OutputTokens int64  `json:"output_tokens"`
	APIKeyID     string `json:"api_key_id"`
	// This field is from variant
	// [AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageCompletionsResult].
	Batch bool `json:"batch"`
	// This field is from variant
	// [AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageCompletionsResult].
	InputAudioTokens int64 `json:"input_audio_tokens"`
	// This field is from variant
	// [AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageCompletionsResult].
	InputCachedTokens int64  `json:"input_cached_tokens"`
	Model             string `json:"model"`
	// This field is from variant
	// [AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageCompletionsResult].
	OutputAudioTokens int64  `json:"output_audio_tokens"`
	ProjectID         string `json:"project_id"`
	// This field is from variant
	// [AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageCompletionsResult].
	ServiceTier string `json:"service_tier"`
	UserID      string `json:"user_id"`
	// This field is from variant
	// [AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageImagesResult].
	Images int64 `json:"images"`
	// This field is from variant
	// [AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageImagesResult].
	Size string `json:"size"`
	// This field is from variant
	// [AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageImagesResult].
	Source string `json:"source"`
	// This field is from variant
	// [AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageAudioSpeechesResult].
	Characters int64 `json:"characters"`
	// This field is from variant
	// [AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageAudioTranscriptionsResult].
	Seconds int64 `json:"seconds"`
	// This field is from variant
	// [AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageVectorStoresResult].
	UsageBytes int64 `json:"usage_bytes"`
	// This field is from variant
	// [AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageCodeInterpreterSessionsResult].
	NumSessions int64 `json:"num_sessions"`
	// This field is from variant
	// [AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationCostsResult].
	Amount AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationCostsResultAmount `json:"amount"`
	// This field is from variant
	// [AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationCostsResult].
	LineItem string `json:"line_item"`
	JSON     struct {
		InputTokens       respjson.Field
		NumModelRequests  respjson.Field
		Object            respjson.Field
		OutputTokens      respjson.Field
		APIKeyID          respjson.Field
		Batch             respjson.Field
		InputAudioTokens  respjson.Field
		InputCachedTokens respjson.Field
		Model             respjson.Field
		OutputAudioTokens respjson.Field
		ProjectID         respjson.Field
		ServiceTier       respjson.Field
		UserID            respjson.Field
		Images            respjson.Field
		Size              respjson.Field
		Source            respjson.Field
		Characters        respjson.Field
		Seconds           respjson.Field
		UsageBytes        respjson.Field
		NumSessions       respjson.Field
		Amount            respjson.Field
		LineItem          respjson.Field
		raw               string
	} `json:"-"`
}

// anyAdminOrganizationUsageAudioTranscriptionsResponseDataResult is implemented by
// each variant of
// [AdminOrganizationUsageAudioTranscriptionsResponseDataResultUnion] to add type
// safety for the return type of
// [AdminOrganizationUsageAudioTranscriptionsResponseDataResultUnion.AsAny]
type anyAdminOrganizationUsageAudioTranscriptionsResponseDataResult interface {
	implAdminOrganizationUsageAudioTranscriptionsResponseDataResultUnion()
}

func (AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageCompletionsResult) implAdminOrganizationUsageAudioTranscriptionsResponseDataResultUnion() {
}
func (AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageEmbeddingsResult) implAdminOrganizationUsageAudioTranscriptionsResponseDataResultUnion() {
}
func (AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageModerationsResult) implAdminOrganizationUsageAudioTranscriptionsResponseDataResultUnion() {
}
func (AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageImagesResult) implAdminOrganizationUsageAudioTranscriptionsResponseDataResultUnion() {
}
func (AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageAudioSpeechesResult) implAdminOrganizationUsageAudioTranscriptionsResponseDataResultUnion() {
}
func (AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageAudioTranscriptionsResult) implAdminOrganizationUsageAudioTranscriptionsResponseDataResultUnion() {
}
func (AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageVectorStoresResult) implAdminOrganizationUsageAudioTranscriptionsResponseDataResultUnion() {
}
func (AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageCodeInterpreterSessionsResult) implAdminOrganizationUsageAudioTranscriptionsResponseDataResultUnion() {
}
func (AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationCostsResult) implAdminOrganizationUsageAudioTranscriptionsResponseDataResultUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := AdminOrganizationUsageAudioTranscriptionsResponseDataResultUnion.AsAny().(type) {
//	case openai.AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageCompletionsResult:
//	case openai.AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageEmbeddingsResult:
//	case openai.AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageModerationsResult:
//	case openai.AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageImagesResult:
//	case openai.AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageAudioSpeechesResult:
//	case openai.AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageAudioTranscriptionsResult:
//	case openai.AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageVectorStoresResult:
//	case openai.AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageCodeInterpreterSessionsResult:
//	case openai.AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationCostsResult:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u AdminOrganizationUsageAudioTranscriptionsResponseDataResultUnion) AsAny() anyAdminOrganizationUsageAudioTranscriptionsResponseDataResult {
	switch u.Object {
	case "organization.usage.completions.result":
		return u.AsOrganizationUsageCompletionsResult()
	case "organization.usage.embeddings.result":
		return u.AsOrganizationUsageEmbeddingsResult()
	case "organization.usage.moderations.result":
		return u.AsOrganizationUsageModerationsResult()
	case "organization.usage.images.result":
		return u.AsOrganizationUsageImagesResult()
	case "organization.usage.audio_speeches.result":
		return u.AsOrganizationUsageAudioSpeechesResult()
	case "organization.usage.audio_transcriptions.result":
		return u.AsOrganizationUsageAudioTranscriptionsResult()
	case "organization.usage.vector_stores.result":
		return u.AsOrganizationUsageVectorStoresResult()
	case "organization.usage.code_interpreter_sessions.result":
		return u.AsOrganizationUsageCodeInterpreterSessionsResult()
	case "organization.costs.result":
		return u.AsOrganizationCostsResult()
	}
	return nil
}

func (u AdminOrganizationUsageAudioTranscriptionsResponseDataResultUnion) AsOrganizationUsageCompletionsResult() (v AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageCompletionsResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageAudioTranscriptionsResponseDataResultUnion) AsOrganizationUsageEmbeddingsResult() (v AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageEmbeddingsResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageAudioTranscriptionsResponseDataResultUnion) AsOrganizationUsageModerationsResult() (v AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageModerationsResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageAudioTranscriptionsResponseDataResultUnion) AsOrganizationUsageImagesResult() (v AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageImagesResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageAudioTranscriptionsResponseDataResultUnion) AsOrganizationUsageAudioSpeechesResult() (v AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageAudioSpeechesResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageAudioTranscriptionsResponseDataResultUnion) AsOrganizationUsageAudioTranscriptionsResult() (v AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageAudioTranscriptionsResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageAudioTranscriptionsResponseDataResultUnion) AsOrganizationUsageVectorStoresResult() (v AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageVectorStoresResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageAudioTranscriptionsResponseDataResultUnion) AsOrganizationUsageCodeInterpreterSessionsResult() (v AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageCodeInterpreterSessionsResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageAudioTranscriptionsResponseDataResultUnion) AsOrganizationCostsResult() (v AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationCostsResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u AdminOrganizationUsageAudioTranscriptionsResponseDataResultUnion) RawJSON() string {
	return u.JSON.raw
}

func (r *AdminOrganizationUsageAudioTranscriptionsResponseDataResultUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated completions usage details of the specific time bucket.
type AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageCompletionsResult struct {
	// The aggregated number of text input tokens used, including cached tokens. For
	// customers subscribe to scale tier, this includes scale tier tokens.
	InputTokens int64 `json:"input_tokens" api:"required"`
	// The count of requests made to the model.
	NumModelRequests int64                                       `json:"num_model_requests" api:"required"`
	Object           constant.OrganizationUsageCompletionsResult `json:"object" default:"organization.usage.completions.result"`
	// The aggregated number of text output tokens used. For customers subscribe to
	// scale tier, this includes scale tier tokens.
	OutputTokens int64 `json:"output_tokens" api:"required"`
	// When `group_by=api_key_id`, this field provides the API key ID of the grouped
	// usage result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=batch`, this field tells whether the grouped usage result is
	// batch or not.
	Batch bool `json:"batch" api:"nullable"`
	// The aggregated number of audio input tokens used, including cached tokens.
	InputAudioTokens int64 `json:"input_audio_tokens"`
	// The aggregated number of text input tokens that has been cached from previous
	// requests. For customers subscribe to scale tier, this includes scale tier
	// tokens.
	InputCachedTokens int64 `json:"input_cached_tokens"`
	// When `group_by=model`, this field provides the model name of the grouped usage
	// result.
	Model string `json:"model" api:"nullable"`
	// The aggregated number of audio output tokens used.
	OutputAudioTokens int64 `json:"output_audio_tokens"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// When `group_by=service_tier`, this field provides the service tier of the
	// grouped usage result.
	ServiceTier string `json:"service_tier" api:"nullable"`
	// When `group_by=user_id`, this field provides the user ID of the grouped usage
	// result.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		InputTokens       respjson.Field
		NumModelRequests  respjson.Field
		Object            respjson.Field
		OutputTokens      respjson.Field
		APIKeyID          respjson.Field
		Batch             respjson.Field
		InputAudioTokens  respjson.Field
		InputCachedTokens respjson.Field
		Model             respjson.Field
		OutputAudioTokens respjson.Field
		ProjectID         respjson.Field
		ServiceTier       respjson.Field
		UserID            respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageCompletionsResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageCompletionsResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated embeddings usage details of the specific time bucket.
type AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageEmbeddingsResult struct {
	// The aggregated number of input tokens used.
	InputTokens int64 `json:"input_tokens" api:"required"`
	// The count of requests made to the model.
	NumModelRequests int64                                      `json:"num_model_requests" api:"required"`
	Object           constant.OrganizationUsageEmbeddingsResult `json:"object" default:"organization.usage.embeddings.result"`
	// When `group_by=api_key_id`, this field provides the API key ID of the grouped
	// usage result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=model`, this field provides the model name of the grouped usage
	// result.
	Model string `json:"model" api:"nullable"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// When `group_by=user_id`, this field provides the user ID of the grouped usage
	// result.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		InputTokens      respjson.Field
		NumModelRequests respjson.Field
		Object           respjson.Field
		APIKeyID         respjson.Field
		Model            respjson.Field
		ProjectID        respjson.Field
		UserID           respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageEmbeddingsResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageEmbeddingsResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated moderations usage details of the specific time bucket.
type AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageModerationsResult struct {
	// The aggregated number of input tokens used.
	InputTokens int64 `json:"input_tokens" api:"required"`
	// The count of requests made to the model.
	NumModelRequests int64                                       `json:"num_model_requests" api:"required"`
	Object           constant.OrganizationUsageModerationsResult `json:"object" default:"organization.usage.moderations.result"`
	// When `group_by=api_key_id`, this field provides the API key ID of the grouped
	// usage result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=model`, this field provides the model name of the grouped usage
	// result.
	Model string `json:"model" api:"nullable"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// When `group_by=user_id`, this field provides the user ID of the grouped usage
	// result.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		InputTokens      respjson.Field
		NumModelRequests respjson.Field
		Object           respjson.Field
		APIKeyID         respjson.Field
		Model            respjson.Field
		ProjectID        respjson.Field
		UserID           respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageModerationsResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageModerationsResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated images usage details of the specific time bucket.
type AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageImagesResult struct {
	// The number of images processed.
	Images int64 `json:"images" api:"required"`
	// The count of requests made to the model.
	NumModelRequests int64                                  `json:"num_model_requests" api:"required"`
	Object           constant.OrganizationUsageImagesResult `json:"object" default:"organization.usage.images.result"`
	// When `group_by=api_key_id`, this field provides the API key ID of the grouped
	// usage result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=model`, this field provides the model name of the grouped usage
	// result.
	Model string `json:"model" api:"nullable"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// When `group_by=size`, this field provides the image size of the grouped usage
	// result.
	Size string `json:"size" api:"nullable"`
	// When `group_by=source`, this field provides the source of the grouped usage
	// result, possible values are `image.generation`, `image.edit`, `image.variation`.
	Source string `json:"source" api:"nullable"`
	// When `group_by=user_id`, this field provides the user ID of the grouped usage
	// result.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Images           respjson.Field
		NumModelRequests respjson.Field
		Object           respjson.Field
		APIKeyID         respjson.Field
		Model            respjson.Field
		ProjectID        respjson.Field
		Size             respjson.Field
		Source           respjson.Field
		UserID           respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageImagesResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageImagesResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated audio speeches usage details of the specific time bucket.
type AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageAudioSpeechesResult struct {
	// The number of characters processed.
	Characters int64 `json:"characters" api:"required"`
	// The count of requests made to the model.
	NumModelRequests int64                                         `json:"num_model_requests" api:"required"`
	Object           constant.OrganizationUsageAudioSpeechesResult `json:"object" default:"organization.usage.audio_speeches.result"`
	// When `group_by=api_key_id`, this field provides the API key ID of the grouped
	// usage result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=model`, this field provides the model name of the grouped usage
	// result.
	Model string `json:"model" api:"nullable"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// When `group_by=user_id`, this field provides the user ID of the grouped usage
	// result.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Characters       respjson.Field
		NumModelRequests respjson.Field
		Object           respjson.Field
		APIKeyID         respjson.Field
		Model            respjson.Field
		ProjectID        respjson.Field
		UserID           respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageAudioSpeechesResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageAudioSpeechesResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated audio transcriptions usage details of the specific time bucket.
type AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageAudioTranscriptionsResult struct {
	// The count of requests made to the model.
	NumModelRequests int64                                               `json:"num_model_requests" api:"required"`
	Object           constant.OrganizationUsageAudioTranscriptionsResult `json:"object" default:"organization.usage.audio_transcriptions.result"`
	// The number of seconds processed.
	Seconds int64 `json:"seconds" api:"required"`
	// When `group_by=api_key_id`, this field provides the API key ID of the grouped
	// usage result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=model`, this field provides the model name of the grouped usage
	// result.
	Model string `json:"model" api:"nullable"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// When `group_by=user_id`, this field provides the user ID of the grouped usage
	// result.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		NumModelRequests respjson.Field
		Object           respjson.Field
		Seconds          respjson.Field
		APIKeyID         respjson.Field
		Model            respjson.Field
		ProjectID        respjson.Field
		UserID           respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageAudioTranscriptionsResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageAudioTranscriptionsResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated vector stores usage details of the specific time bucket.
type AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageVectorStoresResult struct {
	Object constant.OrganizationUsageVectorStoresResult `json:"object" default:"organization.usage.vector_stores.result"`
	// The vector stores usage in bytes.
	UsageBytes int64 `json:"usage_bytes" api:"required"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Object      respjson.Field
		UsageBytes  respjson.Field
		ProjectID   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageVectorStoresResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageVectorStoresResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated code interpreter sessions usage details of the specific time
// bucket.
type AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageCodeInterpreterSessionsResult struct {
	// The number of code interpreter sessions.
	NumSessions int64                                                   `json:"num_sessions" api:"required"`
	Object      constant.OrganizationUsageCodeInterpreterSessionsResult `json:"object" default:"organization.usage.code_interpreter_sessions.result"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		NumSessions respjson.Field
		Object      respjson.Field
		ProjectID   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageCodeInterpreterSessionsResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationUsageCodeInterpreterSessionsResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated costs details of the specific time bucket.
type AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationCostsResult struct {
	Object constant.OrganizationCostsResult `json:"object" default:"organization.costs.result"`
	// The monetary value in its associated currency.
	Amount AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationCostsResultAmount `json:"amount"`
	// When `group_by=api_key_id`, this field provides the API Key ID of the grouped
	// costs result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=line_item`, this field provides the line item of the grouped
	// costs result.
	LineItem string `json:"line_item" api:"nullable"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// costs result.
	ProjectID string `json:"project_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Object      respjson.Field
		Amount      respjson.Field
		APIKeyID    respjson.Field
		LineItem    respjson.Field
		ProjectID   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationCostsResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationCostsResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The monetary value in its associated currency.
type AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationCostsResultAmount struct {
	// Lowercase ISO-4217 currency e.g. "usd"
	Currency string `json:"currency"`
	// The numeric value of the cost.
	Value float64 `json:"value"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Currency    respjson.Field
		Value       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationCostsResultAmount) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageAudioTranscriptionsResponseDataResultOrganizationCostsResultAmount) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationUsageCodeInterpreterSessionsResponse struct {
	Data     []AdminOrganizationUsageCodeInterpreterSessionsResponseData `json:"data" api:"required"`
	HasMore  bool                                                        `json:"has_more" api:"required"`
	NextPage string                                                      `json:"next_page" api:"required"`
	Object   constant.Page                                               `json:"object" default:"page"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Data        respjson.Field
		HasMore     respjson.Field
		NextPage    respjson.Field
		Object      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageCodeInterpreterSessionsResponse) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationUsageCodeInterpreterSessionsResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationUsageCodeInterpreterSessionsResponseData struct {
	EndTime   int64                                                                  `json:"end_time" api:"required"`
	Object    constant.Bucket                                                        `json:"object" default:"bucket"`
	Results   []AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultUnion `json:"results" api:"required"`
	StartTime int64                                                                  `json:"start_time" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EndTime     respjson.Field
		Object      respjson.Field
		Results     respjson.Field
		StartTime   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageCodeInterpreterSessionsResponseData) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageCodeInterpreterSessionsResponseData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultUnion contains
// all possible properties and values from
// [AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageCompletionsResult],
// [AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageEmbeddingsResult],
// [AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageModerationsResult],
// [AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageImagesResult],
// [AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageAudioSpeechesResult],
// [AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageAudioTranscriptionsResult],
// [AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageVectorStoresResult],
// [AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageCodeInterpreterSessionsResult],
// [AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationCostsResult].
//
// Use the
// [AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultUnion.AsAny]
// method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultUnion struct {
	InputTokens      int64 `json:"input_tokens"`
	NumModelRequests int64 `json:"num_model_requests"`
	// Any of "organization.usage.completions.result",
	// "organization.usage.embeddings.result", "organization.usage.moderations.result",
	// "organization.usage.images.result", "organization.usage.audio_speeches.result",
	// "organization.usage.audio_transcriptions.result",
	// "organization.usage.vector_stores.result",
	// "organization.usage.code_interpreter_sessions.result",
	// "organization.costs.result".
	Object string `json:"object"`
	// This field is from variant
	// [AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageCompletionsResult].
	OutputTokens int64  `json:"output_tokens"`
	APIKeyID     string `json:"api_key_id"`
	// This field is from variant
	// [AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageCompletionsResult].
	Batch bool `json:"batch"`
	// This field is from variant
	// [AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageCompletionsResult].
	InputAudioTokens int64 `json:"input_audio_tokens"`
	// This field is from variant
	// [AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageCompletionsResult].
	InputCachedTokens int64  `json:"input_cached_tokens"`
	Model             string `json:"model"`
	// This field is from variant
	// [AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageCompletionsResult].
	OutputAudioTokens int64  `json:"output_audio_tokens"`
	ProjectID         string `json:"project_id"`
	// This field is from variant
	// [AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageCompletionsResult].
	ServiceTier string `json:"service_tier"`
	UserID      string `json:"user_id"`
	// This field is from variant
	// [AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageImagesResult].
	Images int64 `json:"images"`
	// This field is from variant
	// [AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageImagesResult].
	Size string `json:"size"`
	// This field is from variant
	// [AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageImagesResult].
	Source string `json:"source"`
	// This field is from variant
	// [AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageAudioSpeechesResult].
	Characters int64 `json:"characters"`
	// This field is from variant
	// [AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageAudioTranscriptionsResult].
	Seconds int64 `json:"seconds"`
	// This field is from variant
	// [AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageVectorStoresResult].
	UsageBytes int64 `json:"usage_bytes"`
	// This field is from variant
	// [AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageCodeInterpreterSessionsResult].
	NumSessions int64 `json:"num_sessions"`
	// This field is from variant
	// [AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationCostsResult].
	Amount AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationCostsResultAmount `json:"amount"`
	// This field is from variant
	// [AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationCostsResult].
	LineItem string `json:"line_item"`
	JSON     struct {
		InputTokens       respjson.Field
		NumModelRequests  respjson.Field
		Object            respjson.Field
		OutputTokens      respjson.Field
		APIKeyID          respjson.Field
		Batch             respjson.Field
		InputAudioTokens  respjson.Field
		InputCachedTokens respjson.Field
		Model             respjson.Field
		OutputAudioTokens respjson.Field
		ProjectID         respjson.Field
		ServiceTier       respjson.Field
		UserID            respjson.Field
		Images            respjson.Field
		Size              respjson.Field
		Source            respjson.Field
		Characters        respjson.Field
		Seconds           respjson.Field
		UsageBytes        respjson.Field
		NumSessions       respjson.Field
		Amount            respjson.Field
		LineItem          respjson.Field
		raw               string
	} `json:"-"`
}

// anyAdminOrganizationUsageCodeInterpreterSessionsResponseDataResult is
// implemented by each variant of
// [AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultUnion] to add
// type safety for the return type of
// [AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultUnion.AsAny]
type anyAdminOrganizationUsageCodeInterpreterSessionsResponseDataResult interface {
	implAdminOrganizationUsageCodeInterpreterSessionsResponseDataResultUnion()
}

func (AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageCompletionsResult) implAdminOrganizationUsageCodeInterpreterSessionsResponseDataResultUnion() {
}
func (AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageEmbeddingsResult) implAdminOrganizationUsageCodeInterpreterSessionsResponseDataResultUnion() {
}
func (AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageModerationsResult) implAdminOrganizationUsageCodeInterpreterSessionsResponseDataResultUnion() {
}
func (AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageImagesResult) implAdminOrganizationUsageCodeInterpreterSessionsResponseDataResultUnion() {
}
func (AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageAudioSpeechesResult) implAdminOrganizationUsageCodeInterpreterSessionsResponseDataResultUnion() {
}
func (AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageAudioTranscriptionsResult) implAdminOrganizationUsageCodeInterpreterSessionsResponseDataResultUnion() {
}
func (AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageVectorStoresResult) implAdminOrganizationUsageCodeInterpreterSessionsResponseDataResultUnion() {
}
func (AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageCodeInterpreterSessionsResult) implAdminOrganizationUsageCodeInterpreterSessionsResponseDataResultUnion() {
}
func (AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationCostsResult) implAdminOrganizationUsageCodeInterpreterSessionsResponseDataResultUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultUnion.AsAny().(type) {
//	case openai.AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageCompletionsResult:
//	case openai.AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageEmbeddingsResult:
//	case openai.AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageModerationsResult:
//	case openai.AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageImagesResult:
//	case openai.AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageAudioSpeechesResult:
//	case openai.AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageAudioTranscriptionsResult:
//	case openai.AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageVectorStoresResult:
//	case openai.AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageCodeInterpreterSessionsResult:
//	case openai.AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationCostsResult:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultUnion) AsAny() anyAdminOrganizationUsageCodeInterpreterSessionsResponseDataResult {
	switch u.Object {
	case "organization.usage.completions.result":
		return u.AsOrganizationUsageCompletionsResult()
	case "organization.usage.embeddings.result":
		return u.AsOrganizationUsageEmbeddingsResult()
	case "organization.usage.moderations.result":
		return u.AsOrganizationUsageModerationsResult()
	case "organization.usage.images.result":
		return u.AsOrganizationUsageImagesResult()
	case "organization.usage.audio_speeches.result":
		return u.AsOrganizationUsageAudioSpeechesResult()
	case "organization.usage.audio_transcriptions.result":
		return u.AsOrganizationUsageAudioTranscriptionsResult()
	case "organization.usage.vector_stores.result":
		return u.AsOrganizationUsageVectorStoresResult()
	case "organization.usage.code_interpreter_sessions.result":
		return u.AsOrganizationUsageCodeInterpreterSessionsResult()
	case "organization.costs.result":
		return u.AsOrganizationCostsResult()
	}
	return nil
}

func (u AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultUnion) AsOrganizationUsageCompletionsResult() (v AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageCompletionsResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultUnion) AsOrganizationUsageEmbeddingsResult() (v AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageEmbeddingsResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultUnion) AsOrganizationUsageModerationsResult() (v AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageModerationsResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultUnion) AsOrganizationUsageImagesResult() (v AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageImagesResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultUnion) AsOrganizationUsageAudioSpeechesResult() (v AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageAudioSpeechesResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultUnion) AsOrganizationUsageAudioTranscriptionsResult() (v AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageAudioTranscriptionsResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultUnion) AsOrganizationUsageVectorStoresResult() (v AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageVectorStoresResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultUnion) AsOrganizationUsageCodeInterpreterSessionsResult() (v AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageCodeInterpreterSessionsResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultUnion) AsOrganizationCostsResult() (v AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationCostsResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultUnion) RawJSON() string {
	return u.JSON.raw
}

func (r *AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated completions usage details of the specific time bucket.
type AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageCompletionsResult struct {
	// The aggregated number of text input tokens used, including cached tokens. For
	// customers subscribe to scale tier, this includes scale tier tokens.
	InputTokens int64 `json:"input_tokens" api:"required"`
	// The count of requests made to the model.
	NumModelRequests int64                                       `json:"num_model_requests" api:"required"`
	Object           constant.OrganizationUsageCompletionsResult `json:"object" default:"organization.usage.completions.result"`
	// The aggregated number of text output tokens used. For customers subscribe to
	// scale tier, this includes scale tier tokens.
	OutputTokens int64 `json:"output_tokens" api:"required"`
	// When `group_by=api_key_id`, this field provides the API key ID of the grouped
	// usage result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=batch`, this field tells whether the grouped usage result is
	// batch or not.
	Batch bool `json:"batch" api:"nullable"`
	// The aggregated number of audio input tokens used, including cached tokens.
	InputAudioTokens int64 `json:"input_audio_tokens"`
	// The aggregated number of text input tokens that has been cached from previous
	// requests. For customers subscribe to scale tier, this includes scale tier
	// tokens.
	InputCachedTokens int64 `json:"input_cached_tokens"`
	// When `group_by=model`, this field provides the model name of the grouped usage
	// result.
	Model string `json:"model" api:"nullable"`
	// The aggregated number of audio output tokens used.
	OutputAudioTokens int64 `json:"output_audio_tokens"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// When `group_by=service_tier`, this field provides the service tier of the
	// grouped usage result.
	ServiceTier string `json:"service_tier" api:"nullable"`
	// When `group_by=user_id`, this field provides the user ID of the grouped usage
	// result.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		InputTokens       respjson.Field
		NumModelRequests  respjson.Field
		Object            respjson.Field
		OutputTokens      respjson.Field
		APIKeyID          respjson.Field
		Batch             respjson.Field
		InputAudioTokens  respjson.Field
		InputCachedTokens respjson.Field
		Model             respjson.Field
		OutputAudioTokens respjson.Field
		ProjectID         respjson.Field
		ServiceTier       respjson.Field
		UserID            respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageCompletionsResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageCompletionsResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated embeddings usage details of the specific time bucket.
type AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageEmbeddingsResult struct {
	// The aggregated number of input tokens used.
	InputTokens int64 `json:"input_tokens" api:"required"`
	// The count of requests made to the model.
	NumModelRequests int64                                      `json:"num_model_requests" api:"required"`
	Object           constant.OrganizationUsageEmbeddingsResult `json:"object" default:"organization.usage.embeddings.result"`
	// When `group_by=api_key_id`, this field provides the API key ID of the grouped
	// usage result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=model`, this field provides the model name of the grouped usage
	// result.
	Model string `json:"model" api:"nullable"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// When `group_by=user_id`, this field provides the user ID of the grouped usage
	// result.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		InputTokens      respjson.Field
		NumModelRequests respjson.Field
		Object           respjson.Field
		APIKeyID         respjson.Field
		Model            respjson.Field
		ProjectID        respjson.Field
		UserID           respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageEmbeddingsResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageEmbeddingsResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated moderations usage details of the specific time bucket.
type AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageModerationsResult struct {
	// The aggregated number of input tokens used.
	InputTokens int64 `json:"input_tokens" api:"required"`
	// The count of requests made to the model.
	NumModelRequests int64                                       `json:"num_model_requests" api:"required"`
	Object           constant.OrganizationUsageModerationsResult `json:"object" default:"organization.usage.moderations.result"`
	// When `group_by=api_key_id`, this field provides the API key ID of the grouped
	// usage result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=model`, this field provides the model name of the grouped usage
	// result.
	Model string `json:"model" api:"nullable"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// When `group_by=user_id`, this field provides the user ID of the grouped usage
	// result.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		InputTokens      respjson.Field
		NumModelRequests respjson.Field
		Object           respjson.Field
		APIKeyID         respjson.Field
		Model            respjson.Field
		ProjectID        respjson.Field
		UserID           respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageModerationsResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageModerationsResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated images usage details of the specific time bucket.
type AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageImagesResult struct {
	// The number of images processed.
	Images int64 `json:"images" api:"required"`
	// The count of requests made to the model.
	NumModelRequests int64                                  `json:"num_model_requests" api:"required"`
	Object           constant.OrganizationUsageImagesResult `json:"object" default:"organization.usage.images.result"`
	// When `group_by=api_key_id`, this field provides the API key ID of the grouped
	// usage result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=model`, this field provides the model name of the grouped usage
	// result.
	Model string `json:"model" api:"nullable"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// When `group_by=size`, this field provides the image size of the grouped usage
	// result.
	Size string `json:"size" api:"nullable"`
	// When `group_by=source`, this field provides the source of the grouped usage
	// result, possible values are `image.generation`, `image.edit`, `image.variation`.
	Source string `json:"source" api:"nullable"`
	// When `group_by=user_id`, this field provides the user ID of the grouped usage
	// result.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Images           respjson.Field
		NumModelRequests respjson.Field
		Object           respjson.Field
		APIKeyID         respjson.Field
		Model            respjson.Field
		ProjectID        respjson.Field
		Size             respjson.Field
		Source           respjson.Field
		UserID           respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageImagesResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageImagesResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated audio speeches usage details of the specific time bucket.
type AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageAudioSpeechesResult struct {
	// The number of characters processed.
	Characters int64 `json:"characters" api:"required"`
	// The count of requests made to the model.
	NumModelRequests int64                                         `json:"num_model_requests" api:"required"`
	Object           constant.OrganizationUsageAudioSpeechesResult `json:"object" default:"organization.usage.audio_speeches.result"`
	// When `group_by=api_key_id`, this field provides the API key ID of the grouped
	// usage result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=model`, this field provides the model name of the grouped usage
	// result.
	Model string `json:"model" api:"nullable"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// When `group_by=user_id`, this field provides the user ID of the grouped usage
	// result.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Characters       respjson.Field
		NumModelRequests respjson.Field
		Object           respjson.Field
		APIKeyID         respjson.Field
		Model            respjson.Field
		ProjectID        respjson.Field
		UserID           respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageAudioSpeechesResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageAudioSpeechesResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated audio transcriptions usage details of the specific time bucket.
type AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageAudioTranscriptionsResult struct {
	// The count of requests made to the model.
	NumModelRequests int64                                               `json:"num_model_requests" api:"required"`
	Object           constant.OrganizationUsageAudioTranscriptionsResult `json:"object" default:"organization.usage.audio_transcriptions.result"`
	// The number of seconds processed.
	Seconds int64 `json:"seconds" api:"required"`
	// When `group_by=api_key_id`, this field provides the API key ID of the grouped
	// usage result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=model`, this field provides the model name of the grouped usage
	// result.
	Model string `json:"model" api:"nullable"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// When `group_by=user_id`, this field provides the user ID of the grouped usage
	// result.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		NumModelRequests respjson.Field
		Object           respjson.Field
		Seconds          respjson.Field
		APIKeyID         respjson.Field
		Model            respjson.Field
		ProjectID        respjson.Field
		UserID           respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageAudioTranscriptionsResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageAudioTranscriptionsResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated vector stores usage details of the specific time bucket.
type AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageVectorStoresResult struct {
	Object constant.OrganizationUsageVectorStoresResult `json:"object" default:"organization.usage.vector_stores.result"`
	// The vector stores usage in bytes.
	UsageBytes int64 `json:"usage_bytes" api:"required"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Object      respjson.Field
		UsageBytes  respjson.Field
		ProjectID   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageVectorStoresResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageVectorStoresResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated code interpreter sessions usage details of the specific time
// bucket.
type AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageCodeInterpreterSessionsResult struct {
	// The number of code interpreter sessions.
	NumSessions int64                                                   `json:"num_sessions" api:"required"`
	Object      constant.OrganizationUsageCodeInterpreterSessionsResult `json:"object" default:"organization.usage.code_interpreter_sessions.result"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		NumSessions respjson.Field
		Object      respjson.Field
		ProjectID   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageCodeInterpreterSessionsResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationUsageCodeInterpreterSessionsResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated costs details of the specific time bucket.
type AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationCostsResult struct {
	Object constant.OrganizationCostsResult `json:"object" default:"organization.costs.result"`
	// The monetary value in its associated currency.
	Amount AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationCostsResultAmount `json:"amount"`
	// When `group_by=api_key_id`, this field provides the API Key ID of the grouped
	// costs result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=line_item`, this field provides the line item of the grouped
	// costs result.
	LineItem string `json:"line_item" api:"nullable"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// costs result.
	ProjectID string `json:"project_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Object      respjson.Field
		Amount      respjson.Field
		APIKeyID    respjson.Field
		LineItem    respjson.Field
		ProjectID   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationCostsResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationCostsResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The monetary value in its associated currency.
type AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationCostsResultAmount struct {
	// Lowercase ISO-4217 currency e.g. "usd"
	Currency string `json:"currency"`
	// The numeric value of the cost.
	Value float64 `json:"value"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Currency    respjson.Field
		Value       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationCostsResultAmount) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageCodeInterpreterSessionsResponseDataResultOrganizationCostsResultAmount) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationUsageCompletionsResponse struct {
	Data     []AdminOrganizationUsageCompletionsResponseData `json:"data" api:"required"`
	HasMore  bool                                            `json:"has_more" api:"required"`
	NextPage string                                          `json:"next_page" api:"required"`
	Object   constant.Page                                   `json:"object" default:"page"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Data        respjson.Field
		HasMore     respjson.Field
		NextPage    respjson.Field
		Object      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageCompletionsResponse) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationUsageCompletionsResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationUsageCompletionsResponseData struct {
	EndTime   int64                                                      `json:"end_time" api:"required"`
	Object    constant.Bucket                                            `json:"object" default:"bucket"`
	Results   []AdminOrganizationUsageCompletionsResponseDataResultUnion `json:"results" api:"required"`
	StartTime int64                                                      `json:"start_time" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EndTime     respjson.Field
		Object      respjson.Field
		Results     respjson.Field
		StartTime   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageCompletionsResponseData) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationUsageCompletionsResponseData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// AdminOrganizationUsageCompletionsResponseDataResultUnion contains all possible
// properties and values from
// [AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageCompletionsResult],
// [AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageEmbeddingsResult],
// [AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageModerationsResult],
// [AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageImagesResult],
// [AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageAudioSpeechesResult],
// [AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageAudioTranscriptionsResult],
// [AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageVectorStoresResult],
// [AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageCodeInterpreterSessionsResult],
// [AdminOrganizationUsageCompletionsResponseDataResultOrganizationCostsResult].
//
// Use the [AdminOrganizationUsageCompletionsResponseDataResultUnion.AsAny] method
// to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type AdminOrganizationUsageCompletionsResponseDataResultUnion struct {
	InputTokens      int64 `json:"input_tokens"`
	NumModelRequests int64 `json:"num_model_requests"`
	// Any of "organization.usage.completions.result",
	// "organization.usage.embeddings.result", "organization.usage.moderations.result",
	// "organization.usage.images.result", "organization.usage.audio_speeches.result",
	// "organization.usage.audio_transcriptions.result",
	// "organization.usage.vector_stores.result",
	// "organization.usage.code_interpreter_sessions.result",
	// "organization.costs.result".
	Object string `json:"object"`
	// This field is from variant
	// [AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageCompletionsResult].
	OutputTokens int64  `json:"output_tokens"`
	APIKeyID     string `json:"api_key_id"`
	// This field is from variant
	// [AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageCompletionsResult].
	Batch bool `json:"batch"`
	// This field is from variant
	// [AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageCompletionsResult].
	InputAudioTokens int64 `json:"input_audio_tokens"`
	// This field is from variant
	// [AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageCompletionsResult].
	InputCachedTokens int64  `json:"input_cached_tokens"`
	Model             string `json:"model"`
	// This field is from variant
	// [AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageCompletionsResult].
	OutputAudioTokens int64  `json:"output_audio_tokens"`
	ProjectID         string `json:"project_id"`
	// This field is from variant
	// [AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageCompletionsResult].
	ServiceTier string `json:"service_tier"`
	UserID      string `json:"user_id"`
	// This field is from variant
	// [AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageImagesResult].
	Images int64 `json:"images"`
	// This field is from variant
	// [AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageImagesResult].
	Size string `json:"size"`
	// This field is from variant
	// [AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageImagesResult].
	Source string `json:"source"`
	// This field is from variant
	// [AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageAudioSpeechesResult].
	Characters int64 `json:"characters"`
	// This field is from variant
	// [AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageAudioTranscriptionsResult].
	Seconds int64 `json:"seconds"`
	// This field is from variant
	// [AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageVectorStoresResult].
	UsageBytes int64 `json:"usage_bytes"`
	// This field is from variant
	// [AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageCodeInterpreterSessionsResult].
	NumSessions int64 `json:"num_sessions"`
	// This field is from variant
	// [AdminOrganizationUsageCompletionsResponseDataResultOrganizationCostsResult].
	Amount AdminOrganizationUsageCompletionsResponseDataResultOrganizationCostsResultAmount `json:"amount"`
	// This field is from variant
	// [AdminOrganizationUsageCompletionsResponseDataResultOrganizationCostsResult].
	LineItem string `json:"line_item"`
	JSON     struct {
		InputTokens       respjson.Field
		NumModelRequests  respjson.Field
		Object            respjson.Field
		OutputTokens      respjson.Field
		APIKeyID          respjson.Field
		Batch             respjson.Field
		InputAudioTokens  respjson.Field
		InputCachedTokens respjson.Field
		Model             respjson.Field
		OutputAudioTokens respjson.Field
		ProjectID         respjson.Field
		ServiceTier       respjson.Field
		UserID            respjson.Field
		Images            respjson.Field
		Size              respjson.Field
		Source            respjson.Field
		Characters        respjson.Field
		Seconds           respjson.Field
		UsageBytes        respjson.Field
		NumSessions       respjson.Field
		Amount            respjson.Field
		LineItem          respjson.Field
		raw               string
	} `json:"-"`
}

// anyAdminOrganizationUsageCompletionsResponseDataResult is implemented by each
// variant of [AdminOrganizationUsageCompletionsResponseDataResultUnion] to add
// type safety for the return type of
// [AdminOrganizationUsageCompletionsResponseDataResultUnion.AsAny]
type anyAdminOrganizationUsageCompletionsResponseDataResult interface {
	implAdminOrganizationUsageCompletionsResponseDataResultUnion()
}

func (AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageCompletionsResult) implAdminOrganizationUsageCompletionsResponseDataResultUnion() {
}
func (AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageEmbeddingsResult) implAdminOrganizationUsageCompletionsResponseDataResultUnion() {
}
func (AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageModerationsResult) implAdminOrganizationUsageCompletionsResponseDataResultUnion() {
}
func (AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageImagesResult) implAdminOrganizationUsageCompletionsResponseDataResultUnion() {
}
func (AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageAudioSpeechesResult) implAdminOrganizationUsageCompletionsResponseDataResultUnion() {
}
func (AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageAudioTranscriptionsResult) implAdminOrganizationUsageCompletionsResponseDataResultUnion() {
}
func (AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageVectorStoresResult) implAdminOrganizationUsageCompletionsResponseDataResultUnion() {
}
func (AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageCodeInterpreterSessionsResult) implAdminOrganizationUsageCompletionsResponseDataResultUnion() {
}
func (AdminOrganizationUsageCompletionsResponseDataResultOrganizationCostsResult) implAdminOrganizationUsageCompletionsResponseDataResultUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := AdminOrganizationUsageCompletionsResponseDataResultUnion.AsAny().(type) {
//	case openai.AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageCompletionsResult:
//	case openai.AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageEmbeddingsResult:
//	case openai.AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageModerationsResult:
//	case openai.AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageImagesResult:
//	case openai.AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageAudioSpeechesResult:
//	case openai.AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageAudioTranscriptionsResult:
//	case openai.AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageVectorStoresResult:
//	case openai.AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageCodeInterpreterSessionsResult:
//	case openai.AdminOrganizationUsageCompletionsResponseDataResultOrganizationCostsResult:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u AdminOrganizationUsageCompletionsResponseDataResultUnion) AsAny() anyAdminOrganizationUsageCompletionsResponseDataResult {
	switch u.Object {
	case "organization.usage.completions.result":
		return u.AsOrganizationUsageCompletionsResult()
	case "organization.usage.embeddings.result":
		return u.AsOrganizationUsageEmbeddingsResult()
	case "organization.usage.moderations.result":
		return u.AsOrganizationUsageModerationsResult()
	case "organization.usage.images.result":
		return u.AsOrganizationUsageImagesResult()
	case "organization.usage.audio_speeches.result":
		return u.AsOrganizationUsageAudioSpeechesResult()
	case "organization.usage.audio_transcriptions.result":
		return u.AsOrganizationUsageAudioTranscriptionsResult()
	case "organization.usage.vector_stores.result":
		return u.AsOrganizationUsageVectorStoresResult()
	case "organization.usage.code_interpreter_sessions.result":
		return u.AsOrganizationUsageCodeInterpreterSessionsResult()
	case "organization.costs.result":
		return u.AsOrganizationCostsResult()
	}
	return nil
}

func (u AdminOrganizationUsageCompletionsResponseDataResultUnion) AsOrganizationUsageCompletionsResult() (v AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageCompletionsResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageCompletionsResponseDataResultUnion) AsOrganizationUsageEmbeddingsResult() (v AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageEmbeddingsResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageCompletionsResponseDataResultUnion) AsOrganizationUsageModerationsResult() (v AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageModerationsResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageCompletionsResponseDataResultUnion) AsOrganizationUsageImagesResult() (v AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageImagesResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageCompletionsResponseDataResultUnion) AsOrganizationUsageAudioSpeechesResult() (v AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageAudioSpeechesResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageCompletionsResponseDataResultUnion) AsOrganizationUsageAudioTranscriptionsResult() (v AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageAudioTranscriptionsResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageCompletionsResponseDataResultUnion) AsOrganizationUsageVectorStoresResult() (v AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageVectorStoresResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageCompletionsResponseDataResultUnion) AsOrganizationUsageCodeInterpreterSessionsResult() (v AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageCodeInterpreterSessionsResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageCompletionsResponseDataResultUnion) AsOrganizationCostsResult() (v AdminOrganizationUsageCompletionsResponseDataResultOrganizationCostsResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u AdminOrganizationUsageCompletionsResponseDataResultUnion) RawJSON() string { return u.JSON.raw }

func (r *AdminOrganizationUsageCompletionsResponseDataResultUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated completions usage details of the specific time bucket.
type AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageCompletionsResult struct {
	// The aggregated number of text input tokens used, including cached tokens. For
	// customers subscribe to scale tier, this includes scale tier tokens.
	InputTokens int64 `json:"input_tokens" api:"required"`
	// The count of requests made to the model.
	NumModelRequests int64                                       `json:"num_model_requests" api:"required"`
	Object           constant.OrganizationUsageCompletionsResult `json:"object" default:"organization.usage.completions.result"`
	// The aggregated number of text output tokens used. For customers subscribe to
	// scale tier, this includes scale tier tokens.
	OutputTokens int64 `json:"output_tokens" api:"required"`
	// When `group_by=api_key_id`, this field provides the API key ID of the grouped
	// usage result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=batch`, this field tells whether the grouped usage result is
	// batch or not.
	Batch bool `json:"batch" api:"nullable"`
	// The aggregated number of audio input tokens used, including cached tokens.
	InputAudioTokens int64 `json:"input_audio_tokens"`
	// The aggregated number of text input tokens that has been cached from previous
	// requests. For customers subscribe to scale tier, this includes scale tier
	// tokens.
	InputCachedTokens int64 `json:"input_cached_tokens"`
	// When `group_by=model`, this field provides the model name of the grouped usage
	// result.
	Model string `json:"model" api:"nullable"`
	// The aggregated number of audio output tokens used.
	OutputAudioTokens int64 `json:"output_audio_tokens"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// When `group_by=service_tier`, this field provides the service tier of the
	// grouped usage result.
	ServiceTier string `json:"service_tier" api:"nullable"`
	// When `group_by=user_id`, this field provides the user ID of the grouped usage
	// result.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		InputTokens       respjson.Field
		NumModelRequests  respjson.Field
		Object            respjson.Field
		OutputTokens      respjson.Field
		APIKeyID          respjson.Field
		Batch             respjson.Field
		InputAudioTokens  respjson.Field
		InputCachedTokens respjson.Field
		Model             respjson.Field
		OutputAudioTokens respjson.Field
		ProjectID         respjson.Field
		ServiceTier       respjson.Field
		UserID            respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageCompletionsResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageCompletionsResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated embeddings usage details of the specific time bucket.
type AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageEmbeddingsResult struct {
	// The aggregated number of input tokens used.
	InputTokens int64 `json:"input_tokens" api:"required"`
	// The count of requests made to the model.
	NumModelRequests int64                                      `json:"num_model_requests" api:"required"`
	Object           constant.OrganizationUsageEmbeddingsResult `json:"object" default:"organization.usage.embeddings.result"`
	// When `group_by=api_key_id`, this field provides the API key ID of the grouped
	// usage result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=model`, this field provides the model name of the grouped usage
	// result.
	Model string `json:"model" api:"nullable"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// When `group_by=user_id`, this field provides the user ID of the grouped usage
	// result.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		InputTokens      respjson.Field
		NumModelRequests respjson.Field
		Object           respjson.Field
		APIKeyID         respjson.Field
		Model            respjson.Field
		ProjectID        respjson.Field
		UserID           respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageEmbeddingsResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageEmbeddingsResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated moderations usage details of the specific time bucket.
type AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageModerationsResult struct {
	// The aggregated number of input tokens used.
	InputTokens int64 `json:"input_tokens" api:"required"`
	// The count of requests made to the model.
	NumModelRequests int64                                       `json:"num_model_requests" api:"required"`
	Object           constant.OrganizationUsageModerationsResult `json:"object" default:"organization.usage.moderations.result"`
	// When `group_by=api_key_id`, this field provides the API key ID of the grouped
	// usage result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=model`, this field provides the model name of the grouped usage
	// result.
	Model string `json:"model" api:"nullable"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// When `group_by=user_id`, this field provides the user ID of the grouped usage
	// result.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		InputTokens      respjson.Field
		NumModelRequests respjson.Field
		Object           respjson.Field
		APIKeyID         respjson.Field
		Model            respjson.Field
		ProjectID        respjson.Field
		UserID           respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageModerationsResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageModerationsResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated images usage details of the specific time bucket.
type AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageImagesResult struct {
	// The number of images processed.
	Images int64 `json:"images" api:"required"`
	// The count of requests made to the model.
	NumModelRequests int64                                  `json:"num_model_requests" api:"required"`
	Object           constant.OrganizationUsageImagesResult `json:"object" default:"organization.usage.images.result"`
	// When `group_by=api_key_id`, this field provides the API key ID of the grouped
	// usage result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=model`, this field provides the model name of the grouped usage
	// result.
	Model string `json:"model" api:"nullable"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// When `group_by=size`, this field provides the image size of the grouped usage
	// result.
	Size string `json:"size" api:"nullable"`
	// When `group_by=source`, this field provides the source of the grouped usage
	// result, possible values are `image.generation`, `image.edit`, `image.variation`.
	Source string `json:"source" api:"nullable"`
	// When `group_by=user_id`, this field provides the user ID of the grouped usage
	// result.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Images           respjson.Field
		NumModelRequests respjson.Field
		Object           respjson.Field
		APIKeyID         respjson.Field
		Model            respjson.Field
		ProjectID        respjson.Field
		Size             respjson.Field
		Source           respjson.Field
		UserID           respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageImagesResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageImagesResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated audio speeches usage details of the specific time bucket.
type AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageAudioSpeechesResult struct {
	// The number of characters processed.
	Characters int64 `json:"characters" api:"required"`
	// The count of requests made to the model.
	NumModelRequests int64                                         `json:"num_model_requests" api:"required"`
	Object           constant.OrganizationUsageAudioSpeechesResult `json:"object" default:"organization.usage.audio_speeches.result"`
	// When `group_by=api_key_id`, this field provides the API key ID of the grouped
	// usage result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=model`, this field provides the model name of the grouped usage
	// result.
	Model string `json:"model" api:"nullable"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// When `group_by=user_id`, this field provides the user ID of the grouped usage
	// result.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Characters       respjson.Field
		NumModelRequests respjson.Field
		Object           respjson.Field
		APIKeyID         respjson.Field
		Model            respjson.Field
		ProjectID        respjson.Field
		UserID           respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageAudioSpeechesResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageAudioSpeechesResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated audio transcriptions usage details of the specific time bucket.
type AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageAudioTranscriptionsResult struct {
	// The count of requests made to the model.
	NumModelRequests int64                                               `json:"num_model_requests" api:"required"`
	Object           constant.OrganizationUsageAudioTranscriptionsResult `json:"object" default:"organization.usage.audio_transcriptions.result"`
	// The number of seconds processed.
	Seconds int64 `json:"seconds" api:"required"`
	// When `group_by=api_key_id`, this field provides the API key ID of the grouped
	// usage result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=model`, this field provides the model name of the grouped usage
	// result.
	Model string `json:"model" api:"nullable"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// When `group_by=user_id`, this field provides the user ID of the grouped usage
	// result.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		NumModelRequests respjson.Field
		Object           respjson.Field
		Seconds          respjson.Field
		APIKeyID         respjson.Field
		Model            respjson.Field
		ProjectID        respjson.Field
		UserID           respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageAudioTranscriptionsResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageAudioTranscriptionsResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated vector stores usage details of the specific time bucket.
type AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageVectorStoresResult struct {
	Object constant.OrganizationUsageVectorStoresResult `json:"object" default:"organization.usage.vector_stores.result"`
	// The vector stores usage in bytes.
	UsageBytes int64 `json:"usage_bytes" api:"required"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Object      respjson.Field
		UsageBytes  respjson.Field
		ProjectID   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageVectorStoresResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageVectorStoresResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated code interpreter sessions usage details of the specific time
// bucket.
type AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageCodeInterpreterSessionsResult struct {
	// The number of code interpreter sessions.
	NumSessions int64                                                   `json:"num_sessions" api:"required"`
	Object      constant.OrganizationUsageCodeInterpreterSessionsResult `json:"object" default:"organization.usage.code_interpreter_sessions.result"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		NumSessions respjson.Field
		Object      respjson.Field
		ProjectID   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageCodeInterpreterSessionsResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageCompletionsResponseDataResultOrganizationUsageCodeInterpreterSessionsResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated costs details of the specific time bucket.
type AdminOrganizationUsageCompletionsResponseDataResultOrganizationCostsResult struct {
	Object constant.OrganizationCostsResult `json:"object" default:"organization.costs.result"`
	// The monetary value in its associated currency.
	Amount AdminOrganizationUsageCompletionsResponseDataResultOrganizationCostsResultAmount `json:"amount"`
	// When `group_by=api_key_id`, this field provides the API Key ID of the grouped
	// costs result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=line_item`, this field provides the line item of the grouped
	// costs result.
	LineItem string `json:"line_item" api:"nullable"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// costs result.
	ProjectID string `json:"project_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Object      respjson.Field
		Amount      respjson.Field
		APIKeyID    respjson.Field
		LineItem    respjson.Field
		ProjectID   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageCompletionsResponseDataResultOrganizationCostsResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageCompletionsResponseDataResultOrganizationCostsResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The monetary value in its associated currency.
type AdminOrganizationUsageCompletionsResponseDataResultOrganizationCostsResultAmount struct {
	// Lowercase ISO-4217 currency e.g. "usd"
	Currency string `json:"currency"`
	// The numeric value of the cost.
	Value float64 `json:"value"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Currency    respjson.Field
		Value       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageCompletionsResponseDataResultOrganizationCostsResultAmount) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageCompletionsResponseDataResultOrganizationCostsResultAmount) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationUsageCostsResponse struct {
	Data     []AdminOrganizationUsageCostsResponseData `json:"data" api:"required"`
	HasMore  bool                                      `json:"has_more" api:"required"`
	NextPage string                                    `json:"next_page" api:"required"`
	Object   constant.Page                             `json:"object" default:"page"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Data        respjson.Field
		HasMore     respjson.Field
		NextPage    respjson.Field
		Object      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageCostsResponse) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationUsageCostsResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationUsageCostsResponseData struct {
	EndTime   int64                                                `json:"end_time" api:"required"`
	Object    constant.Bucket                                      `json:"object" default:"bucket"`
	Results   []AdminOrganizationUsageCostsResponseDataResultUnion `json:"results" api:"required"`
	StartTime int64                                                `json:"start_time" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EndTime     respjson.Field
		Object      respjson.Field
		Results     respjson.Field
		StartTime   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageCostsResponseData) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationUsageCostsResponseData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// AdminOrganizationUsageCostsResponseDataResultUnion contains all possible
// properties and values from
// [AdminOrganizationUsageCostsResponseDataResultOrganizationUsageCompletionsResult],
// [AdminOrganizationUsageCostsResponseDataResultOrganizationUsageEmbeddingsResult],
// [AdminOrganizationUsageCostsResponseDataResultOrganizationUsageModerationsResult],
// [AdminOrganizationUsageCostsResponseDataResultOrganizationUsageImagesResult],
// [AdminOrganizationUsageCostsResponseDataResultOrganizationUsageAudioSpeechesResult],
// [AdminOrganizationUsageCostsResponseDataResultOrganizationUsageAudioTranscriptionsResult],
// [AdminOrganizationUsageCostsResponseDataResultOrganizationUsageVectorStoresResult],
// [AdminOrganizationUsageCostsResponseDataResultOrganizationUsageCodeInterpreterSessionsResult],
// [AdminOrganizationUsageCostsResponseDataResultOrganizationCostsResult].
//
// Use the [AdminOrganizationUsageCostsResponseDataResultUnion.AsAny] method to
// switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type AdminOrganizationUsageCostsResponseDataResultUnion struct {
	InputTokens      int64 `json:"input_tokens"`
	NumModelRequests int64 `json:"num_model_requests"`
	// Any of "organization.usage.completions.result",
	// "organization.usage.embeddings.result", "organization.usage.moderations.result",
	// "organization.usage.images.result", "organization.usage.audio_speeches.result",
	// "organization.usage.audio_transcriptions.result",
	// "organization.usage.vector_stores.result",
	// "organization.usage.code_interpreter_sessions.result",
	// "organization.costs.result".
	Object string `json:"object"`
	// This field is from variant
	// [AdminOrganizationUsageCostsResponseDataResultOrganizationUsageCompletionsResult].
	OutputTokens int64  `json:"output_tokens"`
	APIKeyID     string `json:"api_key_id"`
	// This field is from variant
	// [AdminOrganizationUsageCostsResponseDataResultOrganizationUsageCompletionsResult].
	Batch bool `json:"batch"`
	// This field is from variant
	// [AdminOrganizationUsageCostsResponseDataResultOrganizationUsageCompletionsResult].
	InputAudioTokens int64 `json:"input_audio_tokens"`
	// This field is from variant
	// [AdminOrganizationUsageCostsResponseDataResultOrganizationUsageCompletionsResult].
	InputCachedTokens int64  `json:"input_cached_tokens"`
	Model             string `json:"model"`
	// This field is from variant
	// [AdminOrganizationUsageCostsResponseDataResultOrganizationUsageCompletionsResult].
	OutputAudioTokens int64  `json:"output_audio_tokens"`
	ProjectID         string `json:"project_id"`
	// This field is from variant
	// [AdminOrganizationUsageCostsResponseDataResultOrganizationUsageCompletionsResult].
	ServiceTier string `json:"service_tier"`
	UserID      string `json:"user_id"`
	// This field is from variant
	// [AdminOrganizationUsageCostsResponseDataResultOrganizationUsageImagesResult].
	Images int64 `json:"images"`
	// This field is from variant
	// [AdminOrganizationUsageCostsResponseDataResultOrganizationUsageImagesResult].
	Size string `json:"size"`
	// This field is from variant
	// [AdminOrganizationUsageCostsResponseDataResultOrganizationUsageImagesResult].
	Source string `json:"source"`
	// This field is from variant
	// [AdminOrganizationUsageCostsResponseDataResultOrganizationUsageAudioSpeechesResult].
	Characters int64 `json:"characters"`
	// This field is from variant
	// [AdminOrganizationUsageCostsResponseDataResultOrganizationUsageAudioTranscriptionsResult].
	Seconds int64 `json:"seconds"`
	// This field is from variant
	// [AdminOrganizationUsageCostsResponseDataResultOrganizationUsageVectorStoresResult].
	UsageBytes int64 `json:"usage_bytes"`
	// This field is from variant
	// [AdminOrganizationUsageCostsResponseDataResultOrganizationUsageCodeInterpreterSessionsResult].
	NumSessions int64 `json:"num_sessions"`
	// This field is from variant
	// [AdminOrganizationUsageCostsResponseDataResultOrganizationCostsResult].
	Amount AdminOrganizationUsageCostsResponseDataResultOrganizationCostsResultAmount `json:"amount"`
	// This field is from variant
	// [AdminOrganizationUsageCostsResponseDataResultOrganizationCostsResult].
	LineItem string `json:"line_item"`
	JSON     struct {
		InputTokens       respjson.Field
		NumModelRequests  respjson.Field
		Object            respjson.Field
		OutputTokens      respjson.Field
		APIKeyID          respjson.Field
		Batch             respjson.Field
		InputAudioTokens  respjson.Field
		InputCachedTokens respjson.Field
		Model             respjson.Field
		OutputAudioTokens respjson.Field
		ProjectID         respjson.Field
		ServiceTier       respjson.Field
		UserID            respjson.Field
		Images            respjson.Field
		Size              respjson.Field
		Source            respjson.Field
		Characters        respjson.Field
		Seconds           respjson.Field
		UsageBytes        respjson.Field
		NumSessions       respjson.Field
		Amount            respjson.Field
		LineItem          respjson.Field
		raw               string
	} `json:"-"`
}

// anyAdminOrganizationUsageCostsResponseDataResult is implemented by each variant
// of [AdminOrganizationUsageCostsResponseDataResultUnion] to add type safety for
// the return type of [AdminOrganizationUsageCostsResponseDataResultUnion.AsAny]
type anyAdminOrganizationUsageCostsResponseDataResult interface {
	implAdminOrganizationUsageCostsResponseDataResultUnion()
}

func (AdminOrganizationUsageCostsResponseDataResultOrganizationUsageCompletionsResult) implAdminOrganizationUsageCostsResponseDataResultUnion() {
}
func (AdminOrganizationUsageCostsResponseDataResultOrganizationUsageEmbeddingsResult) implAdminOrganizationUsageCostsResponseDataResultUnion() {
}
func (AdminOrganizationUsageCostsResponseDataResultOrganizationUsageModerationsResult) implAdminOrganizationUsageCostsResponseDataResultUnion() {
}
func (AdminOrganizationUsageCostsResponseDataResultOrganizationUsageImagesResult) implAdminOrganizationUsageCostsResponseDataResultUnion() {
}
func (AdminOrganizationUsageCostsResponseDataResultOrganizationUsageAudioSpeechesResult) implAdminOrganizationUsageCostsResponseDataResultUnion() {
}
func (AdminOrganizationUsageCostsResponseDataResultOrganizationUsageAudioTranscriptionsResult) implAdminOrganizationUsageCostsResponseDataResultUnion() {
}
func (AdminOrganizationUsageCostsResponseDataResultOrganizationUsageVectorStoresResult) implAdminOrganizationUsageCostsResponseDataResultUnion() {
}
func (AdminOrganizationUsageCostsResponseDataResultOrganizationUsageCodeInterpreterSessionsResult) implAdminOrganizationUsageCostsResponseDataResultUnion() {
}
func (AdminOrganizationUsageCostsResponseDataResultOrganizationCostsResult) implAdminOrganizationUsageCostsResponseDataResultUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := AdminOrganizationUsageCostsResponseDataResultUnion.AsAny().(type) {
//	case openai.AdminOrganizationUsageCostsResponseDataResultOrganizationUsageCompletionsResult:
//	case openai.AdminOrganizationUsageCostsResponseDataResultOrganizationUsageEmbeddingsResult:
//	case openai.AdminOrganizationUsageCostsResponseDataResultOrganizationUsageModerationsResult:
//	case openai.AdminOrganizationUsageCostsResponseDataResultOrganizationUsageImagesResult:
//	case openai.AdminOrganizationUsageCostsResponseDataResultOrganizationUsageAudioSpeechesResult:
//	case openai.AdminOrganizationUsageCostsResponseDataResultOrganizationUsageAudioTranscriptionsResult:
//	case openai.AdminOrganizationUsageCostsResponseDataResultOrganizationUsageVectorStoresResult:
//	case openai.AdminOrganizationUsageCostsResponseDataResultOrganizationUsageCodeInterpreterSessionsResult:
//	case openai.AdminOrganizationUsageCostsResponseDataResultOrganizationCostsResult:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u AdminOrganizationUsageCostsResponseDataResultUnion) AsAny() anyAdminOrganizationUsageCostsResponseDataResult {
	switch u.Object {
	case "organization.usage.completions.result":
		return u.AsOrganizationUsageCompletionsResult()
	case "organization.usage.embeddings.result":
		return u.AsOrganizationUsageEmbeddingsResult()
	case "organization.usage.moderations.result":
		return u.AsOrganizationUsageModerationsResult()
	case "organization.usage.images.result":
		return u.AsOrganizationUsageImagesResult()
	case "organization.usage.audio_speeches.result":
		return u.AsOrganizationUsageAudioSpeechesResult()
	case "organization.usage.audio_transcriptions.result":
		return u.AsOrganizationUsageAudioTranscriptionsResult()
	case "organization.usage.vector_stores.result":
		return u.AsOrganizationUsageVectorStoresResult()
	case "organization.usage.code_interpreter_sessions.result":
		return u.AsOrganizationUsageCodeInterpreterSessionsResult()
	case "organization.costs.result":
		return u.AsOrganizationCostsResult()
	}
	return nil
}

func (u AdminOrganizationUsageCostsResponseDataResultUnion) AsOrganizationUsageCompletionsResult() (v AdminOrganizationUsageCostsResponseDataResultOrganizationUsageCompletionsResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageCostsResponseDataResultUnion) AsOrganizationUsageEmbeddingsResult() (v AdminOrganizationUsageCostsResponseDataResultOrganizationUsageEmbeddingsResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageCostsResponseDataResultUnion) AsOrganizationUsageModerationsResult() (v AdminOrganizationUsageCostsResponseDataResultOrganizationUsageModerationsResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageCostsResponseDataResultUnion) AsOrganizationUsageImagesResult() (v AdminOrganizationUsageCostsResponseDataResultOrganizationUsageImagesResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageCostsResponseDataResultUnion) AsOrganizationUsageAudioSpeechesResult() (v AdminOrganizationUsageCostsResponseDataResultOrganizationUsageAudioSpeechesResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageCostsResponseDataResultUnion) AsOrganizationUsageAudioTranscriptionsResult() (v AdminOrganizationUsageCostsResponseDataResultOrganizationUsageAudioTranscriptionsResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageCostsResponseDataResultUnion) AsOrganizationUsageVectorStoresResult() (v AdminOrganizationUsageCostsResponseDataResultOrganizationUsageVectorStoresResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageCostsResponseDataResultUnion) AsOrganizationUsageCodeInterpreterSessionsResult() (v AdminOrganizationUsageCostsResponseDataResultOrganizationUsageCodeInterpreterSessionsResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageCostsResponseDataResultUnion) AsOrganizationCostsResult() (v AdminOrganizationUsageCostsResponseDataResultOrganizationCostsResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u AdminOrganizationUsageCostsResponseDataResultUnion) RawJSON() string { return u.JSON.raw }

func (r *AdminOrganizationUsageCostsResponseDataResultUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated completions usage details of the specific time bucket.
type AdminOrganizationUsageCostsResponseDataResultOrganizationUsageCompletionsResult struct {
	// The aggregated number of text input tokens used, including cached tokens. For
	// customers subscribe to scale tier, this includes scale tier tokens.
	InputTokens int64 `json:"input_tokens" api:"required"`
	// The count of requests made to the model.
	NumModelRequests int64                                       `json:"num_model_requests" api:"required"`
	Object           constant.OrganizationUsageCompletionsResult `json:"object" default:"organization.usage.completions.result"`
	// The aggregated number of text output tokens used. For customers subscribe to
	// scale tier, this includes scale tier tokens.
	OutputTokens int64 `json:"output_tokens" api:"required"`
	// When `group_by=api_key_id`, this field provides the API key ID of the grouped
	// usage result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=batch`, this field tells whether the grouped usage result is
	// batch or not.
	Batch bool `json:"batch" api:"nullable"`
	// The aggregated number of audio input tokens used, including cached tokens.
	InputAudioTokens int64 `json:"input_audio_tokens"`
	// The aggregated number of text input tokens that has been cached from previous
	// requests. For customers subscribe to scale tier, this includes scale tier
	// tokens.
	InputCachedTokens int64 `json:"input_cached_tokens"`
	// When `group_by=model`, this field provides the model name of the grouped usage
	// result.
	Model string `json:"model" api:"nullable"`
	// The aggregated number of audio output tokens used.
	OutputAudioTokens int64 `json:"output_audio_tokens"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// When `group_by=service_tier`, this field provides the service tier of the
	// grouped usage result.
	ServiceTier string `json:"service_tier" api:"nullable"`
	// When `group_by=user_id`, this field provides the user ID of the grouped usage
	// result.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		InputTokens       respjson.Field
		NumModelRequests  respjson.Field
		Object            respjson.Field
		OutputTokens      respjson.Field
		APIKeyID          respjson.Field
		Batch             respjson.Field
		InputAudioTokens  respjson.Field
		InputCachedTokens respjson.Field
		Model             respjson.Field
		OutputAudioTokens respjson.Field
		ProjectID         respjson.Field
		ServiceTier       respjson.Field
		UserID            respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageCostsResponseDataResultOrganizationUsageCompletionsResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageCostsResponseDataResultOrganizationUsageCompletionsResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated embeddings usage details of the specific time bucket.
type AdminOrganizationUsageCostsResponseDataResultOrganizationUsageEmbeddingsResult struct {
	// The aggregated number of input tokens used.
	InputTokens int64 `json:"input_tokens" api:"required"`
	// The count of requests made to the model.
	NumModelRequests int64                                      `json:"num_model_requests" api:"required"`
	Object           constant.OrganizationUsageEmbeddingsResult `json:"object" default:"organization.usage.embeddings.result"`
	// When `group_by=api_key_id`, this field provides the API key ID of the grouped
	// usage result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=model`, this field provides the model name of the grouped usage
	// result.
	Model string `json:"model" api:"nullable"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// When `group_by=user_id`, this field provides the user ID of the grouped usage
	// result.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		InputTokens      respjson.Field
		NumModelRequests respjson.Field
		Object           respjson.Field
		APIKeyID         respjson.Field
		Model            respjson.Field
		ProjectID        respjson.Field
		UserID           respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageCostsResponseDataResultOrganizationUsageEmbeddingsResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageCostsResponseDataResultOrganizationUsageEmbeddingsResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated moderations usage details of the specific time bucket.
type AdminOrganizationUsageCostsResponseDataResultOrganizationUsageModerationsResult struct {
	// The aggregated number of input tokens used.
	InputTokens int64 `json:"input_tokens" api:"required"`
	// The count of requests made to the model.
	NumModelRequests int64                                       `json:"num_model_requests" api:"required"`
	Object           constant.OrganizationUsageModerationsResult `json:"object" default:"organization.usage.moderations.result"`
	// When `group_by=api_key_id`, this field provides the API key ID of the grouped
	// usage result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=model`, this field provides the model name of the grouped usage
	// result.
	Model string `json:"model" api:"nullable"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// When `group_by=user_id`, this field provides the user ID of the grouped usage
	// result.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		InputTokens      respjson.Field
		NumModelRequests respjson.Field
		Object           respjson.Field
		APIKeyID         respjson.Field
		Model            respjson.Field
		ProjectID        respjson.Field
		UserID           respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageCostsResponseDataResultOrganizationUsageModerationsResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageCostsResponseDataResultOrganizationUsageModerationsResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated images usage details of the specific time bucket.
type AdminOrganizationUsageCostsResponseDataResultOrganizationUsageImagesResult struct {
	// The number of images processed.
	Images int64 `json:"images" api:"required"`
	// The count of requests made to the model.
	NumModelRequests int64                                  `json:"num_model_requests" api:"required"`
	Object           constant.OrganizationUsageImagesResult `json:"object" default:"organization.usage.images.result"`
	// When `group_by=api_key_id`, this field provides the API key ID of the grouped
	// usage result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=model`, this field provides the model name of the grouped usage
	// result.
	Model string `json:"model" api:"nullable"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// When `group_by=size`, this field provides the image size of the grouped usage
	// result.
	Size string `json:"size" api:"nullable"`
	// When `group_by=source`, this field provides the source of the grouped usage
	// result, possible values are `image.generation`, `image.edit`, `image.variation`.
	Source string `json:"source" api:"nullable"`
	// When `group_by=user_id`, this field provides the user ID of the grouped usage
	// result.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Images           respjson.Field
		NumModelRequests respjson.Field
		Object           respjson.Field
		APIKeyID         respjson.Field
		Model            respjson.Field
		ProjectID        respjson.Field
		Size             respjson.Field
		Source           respjson.Field
		UserID           respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageCostsResponseDataResultOrganizationUsageImagesResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageCostsResponseDataResultOrganizationUsageImagesResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated audio speeches usage details of the specific time bucket.
type AdminOrganizationUsageCostsResponseDataResultOrganizationUsageAudioSpeechesResult struct {
	// The number of characters processed.
	Characters int64 `json:"characters" api:"required"`
	// The count of requests made to the model.
	NumModelRequests int64                                         `json:"num_model_requests" api:"required"`
	Object           constant.OrganizationUsageAudioSpeechesResult `json:"object" default:"organization.usage.audio_speeches.result"`
	// When `group_by=api_key_id`, this field provides the API key ID of the grouped
	// usage result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=model`, this field provides the model name of the grouped usage
	// result.
	Model string `json:"model" api:"nullable"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// When `group_by=user_id`, this field provides the user ID of the grouped usage
	// result.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Characters       respjson.Field
		NumModelRequests respjson.Field
		Object           respjson.Field
		APIKeyID         respjson.Field
		Model            respjson.Field
		ProjectID        respjson.Field
		UserID           respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageCostsResponseDataResultOrganizationUsageAudioSpeechesResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageCostsResponseDataResultOrganizationUsageAudioSpeechesResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated audio transcriptions usage details of the specific time bucket.
type AdminOrganizationUsageCostsResponseDataResultOrganizationUsageAudioTranscriptionsResult struct {
	// The count of requests made to the model.
	NumModelRequests int64                                               `json:"num_model_requests" api:"required"`
	Object           constant.OrganizationUsageAudioTranscriptionsResult `json:"object" default:"organization.usage.audio_transcriptions.result"`
	// The number of seconds processed.
	Seconds int64 `json:"seconds" api:"required"`
	// When `group_by=api_key_id`, this field provides the API key ID of the grouped
	// usage result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=model`, this field provides the model name of the grouped usage
	// result.
	Model string `json:"model" api:"nullable"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// When `group_by=user_id`, this field provides the user ID of the grouped usage
	// result.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		NumModelRequests respjson.Field
		Object           respjson.Field
		Seconds          respjson.Field
		APIKeyID         respjson.Field
		Model            respjson.Field
		ProjectID        respjson.Field
		UserID           respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageCostsResponseDataResultOrganizationUsageAudioTranscriptionsResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageCostsResponseDataResultOrganizationUsageAudioTranscriptionsResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated vector stores usage details of the specific time bucket.
type AdminOrganizationUsageCostsResponseDataResultOrganizationUsageVectorStoresResult struct {
	Object constant.OrganizationUsageVectorStoresResult `json:"object" default:"organization.usage.vector_stores.result"`
	// The vector stores usage in bytes.
	UsageBytes int64 `json:"usage_bytes" api:"required"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Object      respjson.Field
		UsageBytes  respjson.Field
		ProjectID   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageCostsResponseDataResultOrganizationUsageVectorStoresResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageCostsResponseDataResultOrganizationUsageVectorStoresResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated code interpreter sessions usage details of the specific time
// bucket.
type AdminOrganizationUsageCostsResponseDataResultOrganizationUsageCodeInterpreterSessionsResult struct {
	// The number of code interpreter sessions.
	NumSessions int64                                                   `json:"num_sessions" api:"required"`
	Object      constant.OrganizationUsageCodeInterpreterSessionsResult `json:"object" default:"organization.usage.code_interpreter_sessions.result"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		NumSessions respjson.Field
		Object      respjson.Field
		ProjectID   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageCostsResponseDataResultOrganizationUsageCodeInterpreterSessionsResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageCostsResponseDataResultOrganizationUsageCodeInterpreterSessionsResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated costs details of the specific time bucket.
type AdminOrganizationUsageCostsResponseDataResultOrganizationCostsResult struct {
	Object constant.OrganizationCostsResult `json:"object" default:"organization.costs.result"`
	// The monetary value in its associated currency.
	Amount AdminOrganizationUsageCostsResponseDataResultOrganizationCostsResultAmount `json:"amount"`
	// When `group_by=api_key_id`, this field provides the API Key ID of the grouped
	// costs result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=line_item`, this field provides the line item of the grouped
	// costs result.
	LineItem string `json:"line_item" api:"nullable"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// costs result.
	ProjectID string `json:"project_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Object      respjson.Field
		Amount      respjson.Field
		APIKeyID    respjson.Field
		LineItem    respjson.Field
		ProjectID   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageCostsResponseDataResultOrganizationCostsResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageCostsResponseDataResultOrganizationCostsResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The monetary value in its associated currency.
type AdminOrganizationUsageCostsResponseDataResultOrganizationCostsResultAmount struct {
	// Lowercase ISO-4217 currency e.g. "usd"
	Currency string `json:"currency"`
	// The numeric value of the cost.
	Value float64 `json:"value"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Currency    respjson.Field
		Value       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageCostsResponseDataResultOrganizationCostsResultAmount) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageCostsResponseDataResultOrganizationCostsResultAmount) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationUsageEmbeddingsResponse struct {
	Data     []AdminOrganizationUsageEmbeddingsResponseData `json:"data" api:"required"`
	HasMore  bool                                           `json:"has_more" api:"required"`
	NextPage string                                         `json:"next_page" api:"required"`
	Object   constant.Page                                  `json:"object" default:"page"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Data        respjson.Field
		HasMore     respjson.Field
		NextPage    respjson.Field
		Object      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageEmbeddingsResponse) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationUsageEmbeddingsResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationUsageEmbeddingsResponseData struct {
	EndTime   int64                                                     `json:"end_time" api:"required"`
	Object    constant.Bucket                                           `json:"object" default:"bucket"`
	Results   []AdminOrganizationUsageEmbeddingsResponseDataResultUnion `json:"results" api:"required"`
	StartTime int64                                                     `json:"start_time" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EndTime     respjson.Field
		Object      respjson.Field
		Results     respjson.Field
		StartTime   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageEmbeddingsResponseData) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationUsageEmbeddingsResponseData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// AdminOrganizationUsageEmbeddingsResponseDataResultUnion contains all possible
// properties and values from
// [AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageCompletionsResult],
// [AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageEmbeddingsResult],
// [AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageModerationsResult],
// [AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageImagesResult],
// [AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageAudioSpeechesResult],
// [AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageAudioTranscriptionsResult],
// [AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageVectorStoresResult],
// [AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageCodeInterpreterSessionsResult],
// [AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationCostsResult].
//
// Use the [AdminOrganizationUsageEmbeddingsResponseDataResultUnion.AsAny] method
// to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type AdminOrganizationUsageEmbeddingsResponseDataResultUnion struct {
	InputTokens      int64 `json:"input_tokens"`
	NumModelRequests int64 `json:"num_model_requests"`
	// Any of "organization.usage.completions.result",
	// "organization.usage.embeddings.result", "organization.usage.moderations.result",
	// "organization.usage.images.result", "organization.usage.audio_speeches.result",
	// "organization.usage.audio_transcriptions.result",
	// "organization.usage.vector_stores.result",
	// "organization.usage.code_interpreter_sessions.result",
	// "organization.costs.result".
	Object string `json:"object"`
	// This field is from variant
	// [AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageCompletionsResult].
	OutputTokens int64  `json:"output_tokens"`
	APIKeyID     string `json:"api_key_id"`
	// This field is from variant
	// [AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageCompletionsResult].
	Batch bool `json:"batch"`
	// This field is from variant
	// [AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageCompletionsResult].
	InputAudioTokens int64 `json:"input_audio_tokens"`
	// This field is from variant
	// [AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageCompletionsResult].
	InputCachedTokens int64  `json:"input_cached_tokens"`
	Model             string `json:"model"`
	// This field is from variant
	// [AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageCompletionsResult].
	OutputAudioTokens int64  `json:"output_audio_tokens"`
	ProjectID         string `json:"project_id"`
	// This field is from variant
	// [AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageCompletionsResult].
	ServiceTier string `json:"service_tier"`
	UserID      string `json:"user_id"`
	// This field is from variant
	// [AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageImagesResult].
	Images int64 `json:"images"`
	// This field is from variant
	// [AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageImagesResult].
	Size string `json:"size"`
	// This field is from variant
	// [AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageImagesResult].
	Source string `json:"source"`
	// This field is from variant
	// [AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageAudioSpeechesResult].
	Characters int64 `json:"characters"`
	// This field is from variant
	// [AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageAudioTranscriptionsResult].
	Seconds int64 `json:"seconds"`
	// This field is from variant
	// [AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageVectorStoresResult].
	UsageBytes int64 `json:"usage_bytes"`
	// This field is from variant
	// [AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageCodeInterpreterSessionsResult].
	NumSessions int64 `json:"num_sessions"`
	// This field is from variant
	// [AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationCostsResult].
	Amount AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationCostsResultAmount `json:"amount"`
	// This field is from variant
	// [AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationCostsResult].
	LineItem string `json:"line_item"`
	JSON     struct {
		InputTokens       respjson.Field
		NumModelRequests  respjson.Field
		Object            respjson.Field
		OutputTokens      respjson.Field
		APIKeyID          respjson.Field
		Batch             respjson.Field
		InputAudioTokens  respjson.Field
		InputCachedTokens respjson.Field
		Model             respjson.Field
		OutputAudioTokens respjson.Field
		ProjectID         respjson.Field
		ServiceTier       respjson.Field
		UserID            respjson.Field
		Images            respjson.Field
		Size              respjson.Field
		Source            respjson.Field
		Characters        respjson.Field
		Seconds           respjson.Field
		UsageBytes        respjson.Field
		NumSessions       respjson.Field
		Amount            respjson.Field
		LineItem          respjson.Field
		raw               string
	} `json:"-"`
}

// anyAdminOrganizationUsageEmbeddingsResponseDataResult is implemented by each
// variant of [AdminOrganizationUsageEmbeddingsResponseDataResultUnion] to add type
// safety for the return type of
// [AdminOrganizationUsageEmbeddingsResponseDataResultUnion.AsAny]
type anyAdminOrganizationUsageEmbeddingsResponseDataResult interface {
	implAdminOrganizationUsageEmbeddingsResponseDataResultUnion()
}

func (AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageCompletionsResult) implAdminOrganizationUsageEmbeddingsResponseDataResultUnion() {
}
func (AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageEmbeddingsResult) implAdminOrganizationUsageEmbeddingsResponseDataResultUnion() {
}
func (AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageModerationsResult) implAdminOrganizationUsageEmbeddingsResponseDataResultUnion() {
}
func (AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageImagesResult) implAdminOrganizationUsageEmbeddingsResponseDataResultUnion() {
}
func (AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageAudioSpeechesResult) implAdminOrganizationUsageEmbeddingsResponseDataResultUnion() {
}
func (AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageAudioTranscriptionsResult) implAdminOrganizationUsageEmbeddingsResponseDataResultUnion() {
}
func (AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageVectorStoresResult) implAdminOrganizationUsageEmbeddingsResponseDataResultUnion() {
}
func (AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageCodeInterpreterSessionsResult) implAdminOrganizationUsageEmbeddingsResponseDataResultUnion() {
}
func (AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationCostsResult) implAdminOrganizationUsageEmbeddingsResponseDataResultUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := AdminOrganizationUsageEmbeddingsResponseDataResultUnion.AsAny().(type) {
//	case openai.AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageCompletionsResult:
//	case openai.AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageEmbeddingsResult:
//	case openai.AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageModerationsResult:
//	case openai.AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageImagesResult:
//	case openai.AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageAudioSpeechesResult:
//	case openai.AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageAudioTranscriptionsResult:
//	case openai.AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageVectorStoresResult:
//	case openai.AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageCodeInterpreterSessionsResult:
//	case openai.AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationCostsResult:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u AdminOrganizationUsageEmbeddingsResponseDataResultUnion) AsAny() anyAdminOrganizationUsageEmbeddingsResponseDataResult {
	switch u.Object {
	case "organization.usage.completions.result":
		return u.AsOrganizationUsageCompletionsResult()
	case "organization.usage.embeddings.result":
		return u.AsOrganizationUsageEmbeddingsResult()
	case "organization.usage.moderations.result":
		return u.AsOrganizationUsageModerationsResult()
	case "organization.usage.images.result":
		return u.AsOrganizationUsageImagesResult()
	case "organization.usage.audio_speeches.result":
		return u.AsOrganizationUsageAudioSpeechesResult()
	case "organization.usage.audio_transcriptions.result":
		return u.AsOrganizationUsageAudioTranscriptionsResult()
	case "organization.usage.vector_stores.result":
		return u.AsOrganizationUsageVectorStoresResult()
	case "organization.usage.code_interpreter_sessions.result":
		return u.AsOrganizationUsageCodeInterpreterSessionsResult()
	case "organization.costs.result":
		return u.AsOrganizationCostsResult()
	}
	return nil
}

func (u AdminOrganizationUsageEmbeddingsResponseDataResultUnion) AsOrganizationUsageCompletionsResult() (v AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageCompletionsResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageEmbeddingsResponseDataResultUnion) AsOrganizationUsageEmbeddingsResult() (v AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageEmbeddingsResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageEmbeddingsResponseDataResultUnion) AsOrganizationUsageModerationsResult() (v AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageModerationsResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageEmbeddingsResponseDataResultUnion) AsOrganizationUsageImagesResult() (v AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageImagesResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageEmbeddingsResponseDataResultUnion) AsOrganizationUsageAudioSpeechesResult() (v AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageAudioSpeechesResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageEmbeddingsResponseDataResultUnion) AsOrganizationUsageAudioTranscriptionsResult() (v AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageAudioTranscriptionsResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageEmbeddingsResponseDataResultUnion) AsOrganizationUsageVectorStoresResult() (v AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageVectorStoresResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageEmbeddingsResponseDataResultUnion) AsOrganizationUsageCodeInterpreterSessionsResult() (v AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageCodeInterpreterSessionsResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageEmbeddingsResponseDataResultUnion) AsOrganizationCostsResult() (v AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationCostsResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u AdminOrganizationUsageEmbeddingsResponseDataResultUnion) RawJSON() string { return u.JSON.raw }

func (r *AdminOrganizationUsageEmbeddingsResponseDataResultUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated completions usage details of the specific time bucket.
type AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageCompletionsResult struct {
	// The aggregated number of text input tokens used, including cached tokens. For
	// customers subscribe to scale tier, this includes scale tier tokens.
	InputTokens int64 `json:"input_tokens" api:"required"`
	// The count of requests made to the model.
	NumModelRequests int64                                       `json:"num_model_requests" api:"required"`
	Object           constant.OrganizationUsageCompletionsResult `json:"object" default:"organization.usage.completions.result"`
	// The aggregated number of text output tokens used. For customers subscribe to
	// scale tier, this includes scale tier tokens.
	OutputTokens int64 `json:"output_tokens" api:"required"`
	// When `group_by=api_key_id`, this field provides the API key ID of the grouped
	// usage result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=batch`, this field tells whether the grouped usage result is
	// batch or not.
	Batch bool `json:"batch" api:"nullable"`
	// The aggregated number of audio input tokens used, including cached tokens.
	InputAudioTokens int64 `json:"input_audio_tokens"`
	// The aggregated number of text input tokens that has been cached from previous
	// requests. For customers subscribe to scale tier, this includes scale tier
	// tokens.
	InputCachedTokens int64 `json:"input_cached_tokens"`
	// When `group_by=model`, this field provides the model name of the grouped usage
	// result.
	Model string `json:"model" api:"nullable"`
	// The aggregated number of audio output tokens used.
	OutputAudioTokens int64 `json:"output_audio_tokens"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// When `group_by=service_tier`, this field provides the service tier of the
	// grouped usage result.
	ServiceTier string `json:"service_tier" api:"nullable"`
	// When `group_by=user_id`, this field provides the user ID of the grouped usage
	// result.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		InputTokens       respjson.Field
		NumModelRequests  respjson.Field
		Object            respjson.Field
		OutputTokens      respjson.Field
		APIKeyID          respjson.Field
		Batch             respjson.Field
		InputAudioTokens  respjson.Field
		InputCachedTokens respjson.Field
		Model             respjson.Field
		OutputAudioTokens respjson.Field
		ProjectID         respjson.Field
		ServiceTier       respjson.Field
		UserID            respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageCompletionsResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageCompletionsResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated embeddings usage details of the specific time bucket.
type AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageEmbeddingsResult struct {
	// The aggregated number of input tokens used.
	InputTokens int64 `json:"input_tokens" api:"required"`
	// The count of requests made to the model.
	NumModelRequests int64                                      `json:"num_model_requests" api:"required"`
	Object           constant.OrganizationUsageEmbeddingsResult `json:"object" default:"organization.usage.embeddings.result"`
	// When `group_by=api_key_id`, this field provides the API key ID of the grouped
	// usage result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=model`, this field provides the model name of the grouped usage
	// result.
	Model string `json:"model" api:"nullable"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// When `group_by=user_id`, this field provides the user ID of the grouped usage
	// result.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		InputTokens      respjson.Field
		NumModelRequests respjson.Field
		Object           respjson.Field
		APIKeyID         respjson.Field
		Model            respjson.Field
		ProjectID        respjson.Field
		UserID           respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageEmbeddingsResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageEmbeddingsResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated moderations usage details of the specific time bucket.
type AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageModerationsResult struct {
	// The aggregated number of input tokens used.
	InputTokens int64 `json:"input_tokens" api:"required"`
	// The count of requests made to the model.
	NumModelRequests int64                                       `json:"num_model_requests" api:"required"`
	Object           constant.OrganizationUsageModerationsResult `json:"object" default:"organization.usage.moderations.result"`
	// When `group_by=api_key_id`, this field provides the API key ID of the grouped
	// usage result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=model`, this field provides the model name of the grouped usage
	// result.
	Model string `json:"model" api:"nullable"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// When `group_by=user_id`, this field provides the user ID of the grouped usage
	// result.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		InputTokens      respjson.Field
		NumModelRequests respjson.Field
		Object           respjson.Field
		APIKeyID         respjson.Field
		Model            respjson.Field
		ProjectID        respjson.Field
		UserID           respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageModerationsResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageModerationsResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated images usage details of the specific time bucket.
type AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageImagesResult struct {
	// The number of images processed.
	Images int64 `json:"images" api:"required"`
	// The count of requests made to the model.
	NumModelRequests int64                                  `json:"num_model_requests" api:"required"`
	Object           constant.OrganizationUsageImagesResult `json:"object" default:"organization.usage.images.result"`
	// When `group_by=api_key_id`, this field provides the API key ID of the grouped
	// usage result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=model`, this field provides the model name of the grouped usage
	// result.
	Model string `json:"model" api:"nullable"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// When `group_by=size`, this field provides the image size of the grouped usage
	// result.
	Size string `json:"size" api:"nullable"`
	// When `group_by=source`, this field provides the source of the grouped usage
	// result, possible values are `image.generation`, `image.edit`, `image.variation`.
	Source string `json:"source" api:"nullable"`
	// When `group_by=user_id`, this field provides the user ID of the grouped usage
	// result.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Images           respjson.Field
		NumModelRequests respjson.Field
		Object           respjson.Field
		APIKeyID         respjson.Field
		Model            respjson.Field
		ProjectID        respjson.Field
		Size             respjson.Field
		Source           respjson.Field
		UserID           respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageImagesResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageImagesResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated audio speeches usage details of the specific time bucket.
type AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageAudioSpeechesResult struct {
	// The number of characters processed.
	Characters int64 `json:"characters" api:"required"`
	// The count of requests made to the model.
	NumModelRequests int64                                         `json:"num_model_requests" api:"required"`
	Object           constant.OrganizationUsageAudioSpeechesResult `json:"object" default:"organization.usage.audio_speeches.result"`
	// When `group_by=api_key_id`, this field provides the API key ID of the grouped
	// usage result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=model`, this field provides the model name of the grouped usage
	// result.
	Model string `json:"model" api:"nullable"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// When `group_by=user_id`, this field provides the user ID of the grouped usage
	// result.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Characters       respjson.Field
		NumModelRequests respjson.Field
		Object           respjson.Field
		APIKeyID         respjson.Field
		Model            respjson.Field
		ProjectID        respjson.Field
		UserID           respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageAudioSpeechesResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageAudioSpeechesResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated audio transcriptions usage details of the specific time bucket.
type AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageAudioTranscriptionsResult struct {
	// The count of requests made to the model.
	NumModelRequests int64                                               `json:"num_model_requests" api:"required"`
	Object           constant.OrganizationUsageAudioTranscriptionsResult `json:"object" default:"organization.usage.audio_transcriptions.result"`
	// The number of seconds processed.
	Seconds int64 `json:"seconds" api:"required"`
	// When `group_by=api_key_id`, this field provides the API key ID of the grouped
	// usage result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=model`, this field provides the model name of the grouped usage
	// result.
	Model string `json:"model" api:"nullable"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// When `group_by=user_id`, this field provides the user ID of the grouped usage
	// result.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		NumModelRequests respjson.Field
		Object           respjson.Field
		Seconds          respjson.Field
		APIKeyID         respjson.Field
		Model            respjson.Field
		ProjectID        respjson.Field
		UserID           respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageAudioTranscriptionsResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageAudioTranscriptionsResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated vector stores usage details of the specific time bucket.
type AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageVectorStoresResult struct {
	Object constant.OrganizationUsageVectorStoresResult `json:"object" default:"organization.usage.vector_stores.result"`
	// The vector stores usage in bytes.
	UsageBytes int64 `json:"usage_bytes" api:"required"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Object      respjson.Field
		UsageBytes  respjson.Field
		ProjectID   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageVectorStoresResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageVectorStoresResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated code interpreter sessions usage details of the specific time
// bucket.
type AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageCodeInterpreterSessionsResult struct {
	// The number of code interpreter sessions.
	NumSessions int64                                                   `json:"num_sessions" api:"required"`
	Object      constant.OrganizationUsageCodeInterpreterSessionsResult `json:"object" default:"organization.usage.code_interpreter_sessions.result"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		NumSessions respjson.Field
		Object      respjson.Field
		ProjectID   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageCodeInterpreterSessionsResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationUsageCodeInterpreterSessionsResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated costs details of the specific time bucket.
type AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationCostsResult struct {
	Object constant.OrganizationCostsResult `json:"object" default:"organization.costs.result"`
	// The monetary value in its associated currency.
	Amount AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationCostsResultAmount `json:"amount"`
	// When `group_by=api_key_id`, this field provides the API Key ID of the grouped
	// costs result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=line_item`, this field provides the line item of the grouped
	// costs result.
	LineItem string `json:"line_item" api:"nullable"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// costs result.
	ProjectID string `json:"project_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Object      respjson.Field
		Amount      respjson.Field
		APIKeyID    respjson.Field
		LineItem    respjson.Field
		ProjectID   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationCostsResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationCostsResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The monetary value in its associated currency.
type AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationCostsResultAmount struct {
	// Lowercase ISO-4217 currency e.g. "usd"
	Currency string `json:"currency"`
	// The numeric value of the cost.
	Value float64 `json:"value"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Currency    respjson.Field
		Value       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationCostsResultAmount) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageEmbeddingsResponseDataResultOrganizationCostsResultAmount) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationUsageImagesResponse struct {
	Data     []AdminOrganizationUsageImagesResponseData `json:"data" api:"required"`
	HasMore  bool                                       `json:"has_more" api:"required"`
	NextPage string                                     `json:"next_page" api:"required"`
	Object   constant.Page                              `json:"object" default:"page"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Data        respjson.Field
		HasMore     respjson.Field
		NextPage    respjson.Field
		Object      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageImagesResponse) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationUsageImagesResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationUsageImagesResponseData struct {
	EndTime   int64                                                 `json:"end_time" api:"required"`
	Object    constant.Bucket                                       `json:"object" default:"bucket"`
	Results   []AdminOrganizationUsageImagesResponseDataResultUnion `json:"results" api:"required"`
	StartTime int64                                                 `json:"start_time" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EndTime     respjson.Field
		Object      respjson.Field
		Results     respjson.Field
		StartTime   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageImagesResponseData) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationUsageImagesResponseData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// AdminOrganizationUsageImagesResponseDataResultUnion contains all possible
// properties and values from
// [AdminOrganizationUsageImagesResponseDataResultOrganizationUsageCompletionsResult],
// [AdminOrganizationUsageImagesResponseDataResultOrganizationUsageEmbeddingsResult],
// [AdminOrganizationUsageImagesResponseDataResultOrganizationUsageModerationsResult],
// [AdminOrganizationUsageImagesResponseDataResultOrganizationUsageImagesResult],
// [AdminOrganizationUsageImagesResponseDataResultOrganizationUsageAudioSpeechesResult],
// [AdminOrganizationUsageImagesResponseDataResultOrganizationUsageAudioTranscriptionsResult],
// [AdminOrganizationUsageImagesResponseDataResultOrganizationUsageVectorStoresResult],
// [AdminOrganizationUsageImagesResponseDataResultOrganizationUsageCodeInterpreterSessionsResult],
// [AdminOrganizationUsageImagesResponseDataResultOrganizationCostsResult].
//
// Use the [AdminOrganizationUsageImagesResponseDataResultUnion.AsAny] method to
// switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type AdminOrganizationUsageImagesResponseDataResultUnion struct {
	InputTokens      int64 `json:"input_tokens"`
	NumModelRequests int64 `json:"num_model_requests"`
	// Any of "organization.usage.completions.result",
	// "organization.usage.embeddings.result", "organization.usage.moderations.result",
	// "organization.usage.images.result", "organization.usage.audio_speeches.result",
	// "organization.usage.audio_transcriptions.result",
	// "organization.usage.vector_stores.result",
	// "organization.usage.code_interpreter_sessions.result",
	// "organization.costs.result".
	Object string `json:"object"`
	// This field is from variant
	// [AdminOrganizationUsageImagesResponseDataResultOrganizationUsageCompletionsResult].
	OutputTokens int64  `json:"output_tokens"`
	APIKeyID     string `json:"api_key_id"`
	// This field is from variant
	// [AdminOrganizationUsageImagesResponseDataResultOrganizationUsageCompletionsResult].
	Batch bool `json:"batch"`
	// This field is from variant
	// [AdminOrganizationUsageImagesResponseDataResultOrganizationUsageCompletionsResult].
	InputAudioTokens int64 `json:"input_audio_tokens"`
	// This field is from variant
	// [AdminOrganizationUsageImagesResponseDataResultOrganizationUsageCompletionsResult].
	InputCachedTokens int64  `json:"input_cached_tokens"`
	Model             string `json:"model"`
	// This field is from variant
	// [AdminOrganizationUsageImagesResponseDataResultOrganizationUsageCompletionsResult].
	OutputAudioTokens int64  `json:"output_audio_tokens"`
	ProjectID         string `json:"project_id"`
	// This field is from variant
	// [AdminOrganizationUsageImagesResponseDataResultOrganizationUsageCompletionsResult].
	ServiceTier string `json:"service_tier"`
	UserID      string `json:"user_id"`
	// This field is from variant
	// [AdminOrganizationUsageImagesResponseDataResultOrganizationUsageImagesResult].
	Images int64 `json:"images"`
	// This field is from variant
	// [AdminOrganizationUsageImagesResponseDataResultOrganizationUsageImagesResult].
	Size string `json:"size"`
	// This field is from variant
	// [AdminOrganizationUsageImagesResponseDataResultOrganizationUsageImagesResult].
	Source string `json:"source"`
	// This field is from variant
	// [AdminOrganizationUsageImagesResponseDataResultOrganizationUsageAudioSpeechesResult].
	Characters int64 `json:"characters"`
	// This field is from variant
	// [AdminOrganizationUsageImagesResponseDataResultOrganizationUsageAudioTranscriptionsResult].
	Seconds int64 `json:"seconds"`
	// This field is from variant
	// [AdminOrganizationUsageImagesResponseDataResultOrganizationUsageVectorStoresResult].
	UsageBytes int64 `json:"usage_bytes"`
	// This field is from variant
	// [AdminOrganizationUsageImagesResponseDataResultOrganizationUsageCodeInterpreterSessionsResult].
	NumSessions int64 `json:"num_sessions"`
	// This field is from variant
	// [AdminOrganizationUsageImagesResponseDataResultOrganizationCostsResult].
	Amount AdminOrganizationUsageImagesResponseDataResultOrganizationCostsResultAmount `json:"amount"`
	// This field is from variant
	// [AdminOrganizationUsageImagesResponseDataResultOrganizationCostsResult].
	LineItem string `json:"line_item"`
	JSON     struct {
		InputTokens       respjson.Field
		NumModelRequests  respjson.Field
		Object            respjson.Field
		OutputTokens      respjson.Field
		APIKeyID          respjson.Field
		Batch             respjson.Field
		InputAudioTokens  respjson.Field
		InputCachedTokens respjson.Field
		Model             respjson.Field
		OutputAudioTokens respjson.Field
		ProjectID         respjson.Field
		ServiceTier       respjson.Field
		UserID            respjson.Field
		Images            respjson.Field
		Size              respjson.Field
		Source            respjson.Field
		Characters        respjson.Field
		Seconds           respjson.Field
		UsageBytes        respjson.Field
		NumSessions       respjson.Field
		Amount            respjson.Field
		LineItem          respjson.Field
		raw               string
	} `json:"-"`
}

// anyAdminOrganizationUsageImagesResponseDataResult is implemented by each variant
// of [AdminOrganizationUsageImagesResponseDataResultUnion] to add type safety for
// the return type of [AdminOrganizationUsageImagesResponseDataResultUnion.AsAny]
type anyAdminOrganizationUsageImagesResponseDataResult interface {
	implAdminOrganizationUsageImagesResponseDataResultUnion()
}

func (AdminOrganizationUsageImagesResponseDataResultOrganizationUsageCompletionsResult) implAdminOrganizationUsageImagesResponseDataResultUnion() {
}
func (AdminOrganizationUsageImagesResponseDataResultOrganizationUsageEmbeddingsResult) implAdminOrganizationUsageImagesResponseDataResultUnion() {
}
func (AdminOrganizationUsageImagesResponseDataResultOrganizationUsageModerationsResult) implAdminOrganizationUsageImagesResponseDataResultUnion() {
}
func (AdminOrganizationUsageImagesResponseDataResultOrganizationUsageImagesResult) implAdminOrganizationUsageImagesResponseDataResultUnion() {
}
func (AdminOrganizationUsageImagesResponseDataResultOrganizationUsageAudioSpeechesResult) implAdminOrganizationUsageImagesResponseDataResultUnion() {
}
func (AdminOrganizationUsageImagesResponseDataResultOrganizationUsageAudioTranscriptionsResult) implAdminOrganizationUsageImagesResponseDataResultUnion() {
}
func (AdminOrganizationUsageImagesResponseDataResultOrganizationUsageVectorStoresResult) implAdminOrganizationUsageImagesResponseDataResultUnion() {
}
func (AdminOrganizationUsageImagesResponseDataResultOrganizationUsageCodeInterpreterSessionsResult) implAdminOrganizationUsageImagesResponseDataResultUnion() {
}
func (AdminOrganizationUsageImagesResponseDataResultOrganizationCostsResult) implAdminOrganizationUsageImagesResponseDataResultUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := AdminOrganizationUsageImagesResponseDataResultUnion.AsAny().(type) {
//	case openai.AdminOrganizationUsageImagesResponseDataResultOrganizationUsageCompletionsResult:
//	case openai.AdminOrganizationUsageImagesResponseDataResultOrganizationUsageEmbeddingsResult:
//	case openai.AdminOrganizationUsageImagesResponseDataResultOrganizationUsageModerationsResult:
//	case openai.AdminOrganizationUsageImagesResponseDataResultOrganizationUsageImagesResult:
//	case openai.AdminOrganizationUsageImagesResponseDataResultOrganizationUsageAudioSpeechesResult:
//	case openai.AdminOrganizationUsageImagesResponseDataResultOrganizationUsageAudioTranscriptionsResult:
//	case openai.AdminOrganizationUsageImagesResponseDataResultOrganizationUsageVectorStoresResult:
//	case openai.AdminOrganizationUsageImagesResponseDataResultOrganizationUsageCodeInterpreterSessionsResult:
//	case openai.AdminOrganizationUsageImagesResponseDataResultOrganizationCostsResult:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u AdminOrganizationUsageImagesResponseDataResultUnion) AsAny() anyAdminOrganizationUsageImagesResponseDataResult {
	switch u.Object {
	case "organization.usage.completions.result":
		return u.AsOrganizationUsageCompletionsResult()
	case "organization.usage.embeddings.result":
		return u.AsOrganizationUsageEmbeddingsResult()
	case "organization.usage.moderations.result":
		return u.AsOrganizationUsageModerationsResult()
	case "organization.usage.images.result":
		return u.AsOrganizationUsageImagesResult()
	case "organization.usage.audio_speeches.result":
		return u.AsOrganizationUsageAudioSpeechesResult()
	case "organization.usage.audio_transcriptions.result":
		return u.AsOrganizationUsageAudioTranscriptionsResult()
	case "organization.usage.vector_stores.result":
		return u.AsOrganizationUsageVectorStoresResult()
	case "organization.usage.code_interpreter_sessions.result":
		return u.AsOrganizationUsageCodeInterpreterSessionsResult()
	case "organization.costs.result":
		return u.AsOrganizationCostsResult()
	}
	return nil
}

func (u AdminOrganizationUsageImagesResponseDataResultUnion) AsOrganizationUsageCompletionsResult() (v AdminOrganizationUsageImagesResponseDataResultOrganizationUsageCompletionsResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageImagesResponseDataResultUnion) AsOrganizationUsageEmbeddingsResult() (v AdminOrganizationUsageImagesResponseDataResultOrganizationUsageEmbeddingsResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageImagesResponseDataResultUnion) AsOrganizationUsageModerationsResult() (v AdminOrganizationUsageImagesResponseDataResultOrganizationUsageModerationsResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageImagesResponseDataResultUnion) AsOrganizationUsageImagesResult() (v AdminOrganizationUsageImagesResponseDataResultOrganizationUsageImagesResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageImagesResponseDataResultUnion) AsOrganizationUsageAudioSpeechesResult() (v AdminOrganizationUsageImagesResponseDataResultOrganizationUsageAudioSpeechesResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageImagesResponseDataResultUnion) AsOrganizationUsageAudioTranscriptionsResult() (v AdminOrganizationUsageImagesResponseDataResultOrganizationUsageAudioTranscriptionsResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageImagesResponseDataResultUnion) AsOrganizationUsageVectorStoresResult() (v AdminOrganizationUsageImagesResponseDataResultOrganizationUsageVectorStoresResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageImagesResponseDataResultUnion) AsOrganizationUsageCodeInterpreterSessionsResult() (v AdminOrganizationUsageImagesResponseDataResultOrganizationUsageCodeInterpreterSessionsResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageImagesResponseDataResultUnion) AsOrganizationCostsResult() (v AdminOrganizationUsageImagesResponseDataResultOrganizationCostsResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u AdminOrganizationUsageImagesResponseDataResultUnion) RawJSON() string { return u.JSON.raw }

func (r *AdminOrganizationUsageImagesResponseDataResultUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated completions usage details of the specific time bucket.
type AdminOrganizationUsageImagesResponseDataResultOrganizationUsageCompletionsResult struct {
	// The aggregated number of text input tokens used, including cached tokens. For
	// customers subscribe to scale tier, this includes scale tier tokens.
	InputTokens int64 `json:"input_tokens" api:"required"`
	// The count of requests made to the model.
	NumModelRequests int64                                       `json:"num_model_requests" api:"required"`
	Object           constant.OrganizationUsageCompletionsResult `json:"object" default:"organization.usage.completions.result"`
	// The aggregated number of text output tokens used. For customers subscribe to
	// scale tier, this includes scale tier tokens.
	OutputTokens int64 `json:"output_tokens" api:"required"`
	// When `group_by=api_key_id`, this field provides the API key ID of the grouped
	// usage result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=batch`, this field tells whether the grouped usage result is
	// batch or not.
	Batch bool `json:"batch" api:"nullable"`
	// The aggregated number of audio input tokens used, including cached tokens.
	InputAudioTokens int64 `json:"input_audio_tokens"`
	// The aggregated number of text input tokens that has been cached from previous
	// requests. For customers subscribe to scale tier, this includes scale tier
	// tokens.
	InputCachedTokens int64 `json:"input_cached_tokens"`
	// When `group_by=model`, this field provides the model name of the grouped usage
	// result.
	Model string `json:"model" api:"nullable"`
	// The aggregated number of audio output tokens used.
	OutputAudioTokens int64 `json:"output_audio_tokens"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// When `group_by=service_tier`, this field provides the service tier of the
	// grouped usage result.
	ServiceTier string `json:"service_tier" api:"nullable"`
	// When `group_by=user_id`, this field provides the user ID of the grouped usage
	// result.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		InputTokens       respjson.Field
		NumModelRequests  respjson.Field
		Object            respjson.Field
		OutputTokens      respjson.Field
		APIKeyID          respjson.Field
		Batch             respjson.Field
		InputAudioTokens  respjson.Field
		InputCachedTokens respjson.Field
		Model             respjson.Field
		OutputAudioTokens respjson.Field
		ProjectID         respjson.Field
		ServiceTier       respjson.Field
		UserID            respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageImagesResponseDataResultOrganizationUsageCompletionsResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageImagesResponseDataResultOrganizationUsageCompletionsResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated embeddings usage details of the specific time bucket.
type AdminOrganizationUsageImagesResponseDataResultOrganizationUsageEmbeddingsResult struct {
	// The aggregated number of input tokens used.
	InputTokens int64 `json:"input_tokens" api:"required"`
	// The count of requests made to the model.
	NumModelRequests int64                                      `json:"num_model_requests" api:"required"`
	Object           constant.OrganizationUsageEmbeddingsResult `json:"object" default:"organization.usage.embeddings.result"`
	// When `group_by=api_key_id`, this field provides the API key ID of the grouped
	// usage result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=model`, this field provides the model name of the grouped usage
	// result.
	Model string `json:"model" api:"nullable"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// When `group_by=user_id`, this field provides the user ID of the grouped usage
	// result.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		InputTokens      respjson.Field
		NumModelRequests respjson.Field
		Object           respjson.Field
		APIKeyID         respjson.Field
		Model            respjson.Field
		ProjectID        respjson.Field
		UserID           respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageImagesResponseDataResultOrganizationUsageEmbeddingsResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageImagesResponseDataResultOrganizationUsageEmbeddingsResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated moderations usage details of the specific time bucket.
type AdminOrganizationUsageImagesResponseDataResultOrganizationUsageModerationsResult struct {
	// The aggregated number of input tokens used.
	InputTokens int64 `json:"input_tokens" api:"required"`
	// The count of requests made to the model.
	NumModelRequests int64                                       `json:"num_model_requests" api:"required"`
	Object           constant.OrganizationUsageModerationsResult `json:"object" default:"organization.usage.moderations.result"`
	// When `group_by=api_key_id`, this field provides the API key ID of the grouped
	// usage result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=model`, this field provides the model name of the grouped usage
	// result.
	Model string `json:"model" api:"nullable"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// When `group_by=user_id`, this field provides the user ID of the grouped usage
	// result.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		InputTokens      respjson.Field
		NumModelRequests respjson.Field
		Object           respjson.Field
		APIKeyID         respjson.Field
		Model            respjson.Field
		ProjectID        respjson.Field
		UserID           respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageImagesResponseDataResultOrganizationUsageModerationsResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageImagesResponseDataResultOrganizationUsageModerationsResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated images usage details of the specific time bucket.
type AdminOrganizationUsageImagesResponseDataResultOrganizationUsageImagesResult struct {
	// The number of images processed.
	Images int64 `json:"images" api:"required"`
	// The count of requests made to the model.
	NumModelRequests int64                                  `json:"num_model_requests" api:"required"`
	Object           constant.OrganizationUsageImagesResult `json:"object" default:"organization.usage.images.result"`
	// When `group_by=api_key_id`, this field provides the API key ID of the grouped
	// usage result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=model`, this field provides the model name of the grouped usage
	// result.
	Model string `json:"model" api:"nullable"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// When `group_by=size`, this field provides the image size of the grouped usage
	// result.
	Size string `json:"size" api:"nullable"`
	// When `group_by=source`, this field provides the source of the grouped usage
	// result, possible values are `image.generation`, `image.edit`, `image.variation`.
	Source string `json:"source" api:"nullable"`
	// When `group_by=user_id`, this field provides the user ID of the grouped usage
	// result.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Images           respjson.Field
		NumModelRequests respjson.Field
		Object           respjson.Field
		APIKeyID         respjson.Field
		Model            respjson.Field
		ProjectID        respjson.Field
		Size             respjson.Field
		Source           respjson.Field
		UserID           respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageImagesResponseDataResultOrganizationUsageImagesResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageImagesResponseDataResultOrganizationUsageImagesResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated audio speeches usage details of the specific time bucket.
type AdminOrganizationUsageImagesResponseDataResultOrganizationUsageAudioSpeechesResult struct {
	// The number of characters processed.
	Characters int64 `json:"characters" api:"required"`
	// The count of requests made to the model.
	NumModelRequests int64                                         `json:"num_model_requests" api:"required"`
	Object           constant.OrganizationUsageAudioSpeechesResult `json:"object" default:"organization.usage.audio_speeches.result"`
	// When `group_by=api_key_id`, this field provides the API key ID of the grouped
	// usage result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=model`, this field provides the model name of the grouped usage
	// result.
	Model string `json:"model" api:"nullable"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// When `group_by=user_id`, this field provides the user ID of the grouped usage
	// result.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Characters       respjson.Field
		NumModelRequests respjson.Field
		Object           respjson.Field
		APIKeyID         respjson.Field
		Model            respjson.Field
		ProjectID        respjson.Field
		UserID           respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageImagesResponseDataResultOrganizationUsageAudioSpeechesResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageImagesResponseDataResultOrganizationUsageAudioSpeechesResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated audio transcriptions usage details of the specific time bucket.
type AdminOrganizationUsageImagesResponseDataResultOrganizationUsageAudioTranscriptionsResult struct {
	// The count of requests made to the model.
	NumModelRequests int64                                               `json:"num_model_requests" api:"required"`
	Object           constant.OrganizationUsageAudioTranscriptionsResult `json:"object" default:"organization.usage.audio_transcriptions.result"`
	// The number of seconds processed.
	Seconds int64 `json:"seconds" api:"required"`
	// When `group_by=api_key_id`, this field provides the API key ID of the grouped
	// usage result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=model`, this field provides the model name of the grouped usage
	// result.
	Model string `json:"model" api:"nullable"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// When `group_by=user_id`, this field provides the user ID of the grouped usage
	// result.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		NumModelRequests respjson.Field
		Object           respjson.Field
		Seconds          respjson.Field
		APIKeyID         respjson.Field
		Model            respjson.Field
		ProjectID        respjson.Field
		UserID           respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageImagesResponseDataResultOrganizationUsageAudioTranscriptionsResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageImagesResponseDataResultOrganizationUsageAudioTranscriptionsResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated vector stores usage details of the specific time bucket.
type AdminOrganizationUsageImagesResponseDataResultOrganizationUsageVectorStoresResult struct {
	Object constant.OrganizationUsageVectorStoresResult `json:"object" default:"organization.usage.vector_stores.result"`
	// The vector stores usage in bytes.
	UsageBytes int64 `json:"usage_bytes" api:"required"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Object      respjson.Field
		UsageBytes  respjson.Field
		ProjectID   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageImagesResponseDataResultOrganizationUsageVectorStoresResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageImagesResponseDataResultOrganizationUsageVectorStoresResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated code interpreter sessions usage details of the specific time
// bucket.
type AdminOrganizationUsageImagesResponseDataResultOrganizationUsageCodeInterpreterSessionsResult struct {
	// The number of code interpreter sessions.
	NumSessions int64                                                   `json:"num_sessions" api:"required"`
	Object      constant.OrganizationUsageCodeInterpreterSessionsResult `json:"object" default:"organization.usage.code_interpreter_sessions.result"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		NumSessions respjson.Field
		Object      respjson.Field
		ProjectID   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageImagesResponseDataResultOrganizationUsageCodeInterpreterSessionsResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageImagesResponseDataResultOrganizationUsageCodeInterpreterSessionsResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated costs details of the specific time bucket.
type AdminOrganizationUsageImagesResponseDataResultOrganizationCostsResult struct {
	Object constant.OrganizationCostsResult `json:"object" default:"organization.costs.result"`
	// The monetary value in its associated currency.
	Amount AdminOrganizationUsageImagesResponseDataResultOrganizationCostsResultAmount `json:"amount"`
	// When `group_by=api_key_id`, this field provides the API Key ID of the grouped
	// costs result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=line_item`, this field provides the line item of the grouped
	// costs result.
	LineItem string `json:"line_item" api:"nullable"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// costs result.
	ProjectID string `json:"project_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Object      respjson.Field
		Amount      respjson.Field
		APIKeyID    respjson.Field
		LineItem    respjson.Field
		ProjectID   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageImagesResponseDataResultOrganizationCostsResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageImagesResponseDataResultOrganizationCostsResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The monetary value in its associated currency.
type AdminOrganizationUsageImagesResponseDataResultOrganizationCostsResultAmount struct {
	// Lowercase ISO-4217 currency e.g. "usd"
	Currency string `json:"currency"`
	// The numeric value of the cost.
	Value float64 `json:"value"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Currency    respjson.Field
		Value       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageImagesResponseDataResultOrganizationCostsResultAmount) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageImagesResponseDataResultOrganizationCostsResultAmount) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationUsageModerationsResponse struct {
	Data     []AdminOrganizationUsageModerationsResponseData `json:"data" api:"required"`
	HasMore  bool                                            `json:"has_more" api:"required"`
	NextPage string                                          `json:"next_page" api:"required"`
	Object   constant.Page                                   `json:"object" default:"page"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Data        respjson.Field
		HasMore     respjson.Field
		NextPage    respjson.Field
		Object      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageModerationsResponse) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationUsageModerationsResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationUsageModerationsResponseData struct {
	EndTime   int64                                                      `json:"end_time" api:"required"`
	Object    constant.Bucket                                            `json:"object" default:"bucket"`
	Results   []AdminOrganizationUsageModerationsResponseDataResultUnion `json:"results" api:"required"`
	StartTime int64                                                      `json:"start_time" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EndTime     respjson.Field
		Object      respjson.Field
		Results     respjson.Field
		StartTime   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageModerationsResponseData) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationUsageModerationsResponseData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// AdminOrganizationUsageModerationsResponseDataResultUnion contains all possible
// properties and values from
// [AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageCompletionsResult],
// [AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageEmbeddingsResult],
// [AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageModerationsResult],
// [AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageImagesResult],
// [AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageAudioSpeechesResult],
// [AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageAudioTranscriptionsResult],
// [AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageVectorStoresResult],
// [AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageCodeInterpreterSessionsResult],
// [AdminOrganizationUsageModerationsResponseDataResultOrganizationCostsResult].
//
// Use the [AdminOrganizationUsageModerationsResponseDataResultUnion.AsAny] method
// to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type AdminOrganizationUsageModerationsResponseDataResultUnion struct {
	InputTokens      int64 `json:"input_tokens"`
	NumModelRequests int64 `json:"num_model_requests"`
	// Any of "organization.usage.completions.result",
	// "organization.usage.embeddings.result", "organization.usage.moderations.result",
	// "organization.usage.images.result", "organization.usage.audio_speeches.result",
	// "organization.usage.audio_transcriptions.result",
	// "organization.usage.vector_stores.result",
	// "organization.usage.code_interpreter_sessions.result",
	// "organization.costs.result".
	Object string `json:"object"`
	// This field is from variant
	// [AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageCompletionsResult].
	OutputTokens int64  `json:"output_tokens"`
	APIKeyID     string `json:"api_key_id"`
	// This field is from variant
	// [AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageCompletionsResult].
	Batch bool `json:"batch"`
	// This field is from variant
	// [AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageCompletionsResult].
	InputAudioTokens int64 `json:"input_audio_tokens"`
	// This field is from variant
	// [AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageCompletionsResult].
	InputCachedTokens int64  `json:"input_cached_tokens"`
	Model             string `json:"model"`
	// This field is from variant
	// [AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageCompletionsResult].
	OutputAudioTokens int64  `json:"output_audio_tokens"`
	ProjectID         string `json:"project_id"`
	// This field is from variant
	// [AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageCompletionsResult].
	ServiceTier string `json:"service_tier"`
	UserID      string `json:"user_id"`
	// This field is from variant
	// [AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageImagesResult].
	Images int64 `json:"images"`
	// This field is from variant
	// [AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageImagesResult].
	Size string `json:"size"`
	// This field is from variant
	// [AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageImagesResult].
	Source string `json:"source"`
	// This field is from variant
	// [AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageAudioSpeechesResult].
	Characters int64 `json:"characters"`
	// This field is from variant
	// [AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageAudioTranscriptionsResult].
	Seconds int64 `json:"seconds"`
	// This field is from variant
	// [AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageVectorStoresResult].
	UsageBytes int64 `json:"usage_bytes"`
	// This field is from variant
	// [AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageCodeInterpreterSessionsResult].
	NumSessions int64 `json:"num_sessions"`
	// This field is from variant
	// [AdminOrganizationUsageModerationsResponseDataResultOrganizationCostsResult].
	Amount AdminOrganizationUsageModerationsResponseDataResultOrganizationCostsResultAmount `json:"amount"`
	// This field is from variant
	// [AdminOrganizationUsageModerationsResponseDataResultOrganizationCostsResult].
	LineItem string `json:"line_item"`
	JSON     struct {
		InputTokens       respjson.Field
		NumModelRequests  respjson.Field
		Object            respjson.Field
		OutputTokens      respjson.Field
		APIKeyID          respjson.Field
		Batch             respjson.Field
		InputAudioTokens  respjson.Field
		InputCachedTokens respjson.Field
		Model             respjson.Field
		OutputAudioTokens respjson.Field
		ProjectID         respjson.Field
		ServiceTier       respjson.Field
		UserID            respjson.Field
		Images            respjson.Field
		Size              respjson.Field
		Source            respjson.Field
		Characters        respjson.Field
		Seconds           respjson.Field
		UsageBytes        respjson.Field
		NumSessions       respjson.Field
		Amount            respjson.Field
		LineItem          respjson.Field
		raw               string
	} `json:"-"`
}

// anyAdminOrganizationUsageModerationsResponseDataResult is implemented by each
// variant of [AdminOrganizationUsageModerationsResponseDataResultUnion] to add
// type safety for the return type of
// [AdminOrganizationUsageModerationsResponseDataResultUnion.AsAny]
type anyAdminOrganizationUsageModerationsResponseDataResult interface {
	implAdminOrganizationUsageModerationsResponseDataResultUnion()
}

func (AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageCompletionsResult) implAdminOrganizationUsageModerationsResponseDataResultUnion() {
}
func (AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageEmbeddingsResult) implAdminOrganizationUsageModerationsResponseDataResultUnion() {
}
func (AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageModerationsResult) implAdminOrganizationUsageModerationsResponseDataResultUnion() {
}
func (AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageImagesResult) implAdminOrganizationUsageModerationsResponseDataResultUnion() {
}
func (AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageAudioSpeechesResult) implAdminOrganizationUsageModerationsResponseDataResultUnion() {
}
func (AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageAudioTranscriptionsResult) implAdminOrganizationUsageModerationsResponseDataResultUnion() {
}
func (AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageVectorStoresResult) implAdminOrganizationUsageModerationsResponseDataResultUnion() {
}
func (AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageCodeInterpreterSessionsResult) implAdminOrganizationUsageModerationsResponseDataResultUnion() {
}
func (AdminOrganizationUsageModerationsResponseDataResultOrganizationCostsResult) implAdminOrganizationUsageModerationsResponseDataResultUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := AdminOrganizationUsageModerationsResponseDataResultUnion.AsAny().(type) {
//	case openai.AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageCompletionsResult:
//	case openai.AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageEmbeddingsResult:
//	case openai.AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageModerationsResult:
//	case openai.AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageImagesResult:
//	case openai.AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageAudioSpeechesResult:
//	case openai.AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageAudioTranscriptionsResult:
//	case openai.AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageVectorStoresResult:
//	case openai.AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageCodeInterpreterSessionsResult:
//	case openai.AdminOrganizationUsageModerationsResponseDataResultOrganizationCostsResult:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u AdminOrganizationUsageModerationsResponseDataResultUnion) AsAny() anyAdminOrganizationUsageModerationsResponseDataResult {
	switch u.Object {
	case "organization.usage.completions.result":
		return u.AsOrganizationUsageCompletionsResult()
	case "organization.usage.embeddings.result":
		return u.AsOrganizationUsageEmbeddingsResult()
	case "organization.usage.moderations.result":
		return u.AsOrganizationUsageModerationsResult()
	case "organization.usage.images.result":
		return u.AsOrganizationUsageImagesResult()
	case "organization.usage.audio_speeches.result":
		return u.AsOrganizationUsageAudioSpeechesResult()
	case "organization.usage.audio_transcriptions.result":
		return u.AsOrganizationUsageAudioTranscriptionsResult()
	case "organization.usage.vector_stores.result":
		return u.AsOrganizationUsageVectorStoresResult()
	case "organization.usage.code_interpreter_sessions.result":
		return u.AsOrganizationUsageCodeInterpreterSessionsResult()
	case "organization.costs.result":
		return u.AsOrganizationCostsResult()
	}
	return nil
}

func (u AdminOrganizationUsageModerationsResponseDataResultUnion) AsOrganizationUsageCompletionsResult() (v AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageCompletionsResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageModerationsResponseDataResultUnion) AsOrganizationUsageEmbeddingsResult() (v AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageEmbeddingsResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageModerationsResponseDataResultUnion) AsOrganizationUsageModerationsResult() (v AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageModerationsResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageModerationsResponseDataResultUnion) AsOrganizationUsageImagesResult() (v AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageImagesResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageModerationsResponseDataResultUnion) AsOrganizationUsageAudioSpeechesResult() (v AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageAudioSpeechesResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageModerationsResponseDataResultUnion) AsOrganizationUsageAudioTranscriptionsResult() (v AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageAudioTranscriptionsResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageModerationsResponseDataResultUnion) AsOrganizationUsageVectorStoresResult() (v AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageVectorStoresResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageModerationsResponseDataResultUnion) AsOrganizationUsageCodeInterpreterSessionsResult() (v AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageCodeInterpreterSessionsResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageModerationsResponseDataResultUnion) AsOrganizationCostsResult() (v AdminOrganizationUsageModerationsResponseDataResultOrganizationCostsResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u AdminOrganizationUsageModerationsResponseDataResultUnion) RawJSON() string { return u.JSON.raw }

func (r *AdminOrganizationUsageModerationsResponseDataResultUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated completions usage details of the specific time bucket.
type AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageCompletionsResult struct {
	// The aggregated number of text input tokens used, including cached tokens. For
	// customers subscribe to scale tier, this includes scale tier tokens.
	InputTokens int64 `json:"input_tokens" api:"required"`
	// The count of requests made to the model.
	NumModelRequests int64                                       `json:"num_model_requests" api:"required"`
	Object           constant.OrganizationUsageCompletionsResult `json:"object" default:"organization.usage.completions.result"`
	// The aggregated number of text output tokens used. For customers subscribe to
	// scale tier, this includes scale tier tokens.
	OutputTokens int64 `json:"output_tokens" api:"required"`
	// When `group_by=api_key_id`, this field provides the API key ID of the grouped
	// usage result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=batch`, this field tells whether the grouped usage result is
	// batch or not.
	Batch bool `json:"batch" api:"nullable"`
	// The aggregated number of audio input tokens used, including cached tokens.
	InputAudioTokens int64 `json:"input_audio_tokens"`
	// The aggregated number of text input tokens that has been cached from previous
	// requests. For customers subscribe to scale tier, this includes scale tier
	// tokens.
	InputCachedTokens int64 `json:"input_cached_tokens"`
	// When `group_by=model`, this field provides the model name of the grouped usage
	// result.
	Model string `json:"model" api:"nullable"`
	// The aggregated number of audio output tokens used.
	OutputAudioTokens int64 `json:"output_audio_tokens"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// When `group_by=service_tier`, this field provides the service tier of the
	// grouped usage result.
	ServiceTier string `json:"service_tier" api:"nullable"`
	// When `group_by=user_id`, this field provides the user ID of the grouped usage
	// result.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		InputTokens       respjson.Field
		NumModelRequests  respjson.Field
		Object            respjson.Field
		OutputTokens      respjson.Field
		APIKeyID          respjson.Field
		Batch             respjson.Field
		InputAudioTokens  respjson.Field
		InputCachedTokens respjson.Field
		Model             respjson.Field
		OutputAudioTokens respjson.Field
		ProjectID         respjson.Field
		ServiceTier       respjson.Field
		UserID            respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageCompletionsResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageCompletionsResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated embeddings usage details of the specific time bucket.
type AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageEmbeddingsResult struct {
	// The aggregated number of input tokens used.
	InputTokens int64 `json:"input_tokens" api:"required"`
	// The count of requests made to the model.
	NumModelRequests int64                                      `json:"num_model_requests" api:"required"`
	Object           constant.OrganizationUsageEmbeddingsResult `json:"object" default:"organization.usage.embeddings.result"`
	// When `group_by=api_key_id`, this field provides the API key ID of the grouped
	// usage result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=model`, this field provides the model name of the grouped usage
	// result.
	Model string `json:"model" api:"nullable"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// When `group_by=user_id`, this field provides the user ID of the grouped usage
	// result.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		InputTokens      respjson.Field
		NumModelRequests respjson.Field
		Object           respjson.Field
		APIKeyID         respjson.Field
		Model            respjson.Field
		ProjectID        respjson.Field
		UserID           respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageEmbeddingsResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageEmbeddingsResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated moderations usage details of the specific time bucket.
type AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageModerationsResult struct {
	// The aggregated number of input tokens used.
	InputTokens int64 `json:"input_tokens" api:"required"`
	// The count of requests made to the model.
	NumModelRequests int64                                       `json:"num_model_requests" api:"required"`
	Object           constant.OrganizationUsageModerationsResult `json:"object" default:"organization.usage.moderations.result"`
	// When `group_by=api_key_id`, this field provides the API key ID of the grouped
	// usage result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=model`, this field provides the model name of the grouped usage
	// result.
	Model string `json:"model" api:"nullable"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// When `group_by=user_id`, this field provides the user ID of the grouped usage
	// result.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		InputTokens      respjson.Field
		NumModelRequests respjson.Field
		Object           respjson.Field
		APIKeyID         respjson.Field
		Model            respjson.Field
		ProjectID        respjson.Field
		UserID           respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageModerationsResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageModerationsResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated images usage details of the specific time bucket.
type AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageImagesResult struct {
	// The number of images processed.
	Images int64 `json:"images" api:"required"`
	// The count of requests made to the model.
	NumModelRequests int64                                  `json:"num_model_requests" api:"required"`
	Object           constant.OrganizationUsageImagesResult `json:"object" default:"organization.usage.images.result"`
	// When `group_by=api_key_id`, this field provides the API key ID of the grouped
	// usage result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=model`, this field provides the model name of the grouped usage
	// result.
	Model string `json:"model" api:"nullable"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// When `group_by=size`, this field provides the image size of the grouped usage
	// result.
	Size string `json:"size" api:"nullable"`
	// When `group_by=source`, this field provides the source of the grouped usage
	// result, possible values are `image.generation`, `image.edit`, `image.variation`.
	Source string `json:"source" api:"nullable"`
	// When `group_by=user_id`, this field provides the user ID of the grouped usage
	// result.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Images           respjson.Field
		NumModelRequests respjson.Field
		Object           respjson.Field
		APIKeyID         respjson.Field
		Model            respjson.Field
		ProjectID        respjson.Field
		Size             respjson.Field
		Source           respjson.Field
		UserID           respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageImagesResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageImagesResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated audio speeches usage details of the specific time bucket.
type AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageAudioSpeechesResult struct {
	// The number of characters processed.
	Characters int64 `json:"characters" api:"required"`
	// The count of requests made to the model.
	NumModelRequests int64                                         `json:"num_model_requests" api:"required"`
	Object           constant.OrganizationUsageAudioSpeechesResult `json:"object" default:"organization.usage.audio_speeches.result"`
	// When `group_by=api_key_id`, this field provides the API key ID of the grouped
	// usage result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=model`, this field provides the model name of the grouped usage
	// result.
	Model string `json:"model" api:"nullable"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// When `group_by=user_id`, this field provides the user ID of the grouped usage
	// result.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Characters       respjson.Field
		NumModelRequests respjson.Field
		Object           respjson.Field
		APIKeyID         respjson.Field
		Model            respjson.Field
		ProjectID        respjson.Field
		UserID           respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageAudioSpeechesResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageAudioSpeechesResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated audio transcriptions usage details of the specific time bucket.
type AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageAudioTranscriptionsResult struct {
	// The count of requests made to the model.
	NumModelRequests int64                                               `json:"num_model_requests" api:"required"`
	Object           constant.OrganizationUsageAudioTranscriptionsResult `json:"object" default:"organization.usage.audio_transcriptions.result"`
	// The number of seconds processed.
	Seconds int64 `json:"seconds" api:"required"`
	// When `group_by=api_key_id`, this field provides the API key ID of the grouped
	// usage result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=model`, this field provides the model name of the grouped usage
	// result.
	Model string `json:"model" api:"nullable"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// When `group_by=user_id`, this field provides the user ID of the grouped usage
	// result.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		NumModelRequests respjson.Field
		Object           respjson.Field
		Seconds          respjson.Field
		APIKeyID         respjson.Field
		Model            respjson.Field
		ProjectID        respjson.Field
		UserID           respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageAudioTranscriptionsResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageAudioTranscriptionsResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated vector stores usage details of the specific time bucket.
type AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageVectorStoresResult struct {
	Object constant.OrganizationUsageVectorStoresResult `json:"object" default:"organization.usage.vector_stores.result"`
	// The vector stores usage in bytes.
	UsageBytes int64 `json:"usage_bytes" api:"required"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Object      respjson.Field
		UsageBytes  respjson.Field
		ProjectID   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageVectorStoresResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageVectorStoresResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated code interpreter sessions usage details of the specific time
// bucket.
type AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageCodeInterpreterSessionsResult struct {
	// The number of code interpreter sessions.
	NumSessions int64                                                   `json:"num_sessions" api:"required"`
	Object      constant.OrganizationUsageCodeInterpreterSessionsResult `json:"object" default:"organization.usage.code_interpreter_sessions.result"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		NumSessions respjson.Field
		Object      respjson.Field
		ProjectID   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageCodeInterpreterSessionsResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageModerationsResponseDataResultOrganizationUsageCodeInterpreterSessionsResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated costs details of the specific time bucket.
type AdminOrganizationUsageModerationsResponseDataResultOrganizationCostsResult struct {
	Object constant.OrganizationCostsResult `json:"object" default:"organization.costs.result"`
	// The monetary value in its associated currency.
	Amount AdminOrganizationUsageModerationsResponseDataResultOrganizationCostsResultAmount `json:"amount"`
	// When `group_by=api_key_id`, this field provides the API Key ID of the grouped
	// costs result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=line_item`, this field provides the line item of the grouped
	// costs result.
	LineItem string `json:"line_item" api:"nullable"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// costs result.
	ProjectID string `json:"project_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Object      respjson.Field
		Amount      respjson.Field
		APIKeyID    respjson.Field
		LineItem    respjson.Field
		ProjectID   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageModerationsResponseDataResultOrganizationCostsResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageModerationsResponseDataResultOrganizationCostsResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The monetary value in its associated currency.
type AdminOrganizationUsageModerationsResponseDataResultOrganizationCostsResultAmount struct {
	// Lowercase ISO-4217 currency e.g. "usd"
	Currency string `json:"currency"`
	// The numeric value of the cost.
	Value float64 `json:"value"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Currency    respjson.Field
		Value       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageModerationsResponseDataResultOrganizationCostsResultAmount) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageModerationsResponseDataResultOrganizationCostsResultAmount) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationUsageVectorStoresResponse struct {
	Data     []AdminOrganizationUsageVectorStoresResponseData `json:"data" api:"required"`
	HasMore  bool                                             `json:"has_more" api:"required"`
	NextPage string                                           `json:"next_page" api:"required"`
	Object   constant.Page                                    `json:"object" default:"page"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Data        respjson.Field
		HasMore     respjson.Field
		NextPage    respjson.Field
		Object      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageVectorStoresResponse) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationUsageVectorStoresResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationUsageVectorStoresResponseData struct {
	EndTime   int64                                                       `json:"end_time" api:"required"`
	Object    constant.Bucket                                             `json:"object" default:"bucket"`
	Results   []AdminOrganizationUsageVectorStoresResponseDataResultUnion `json:"results" api:"required"`
	StartTime int64                                                       `json:"start_time" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EndTime     respjson.Field
		Object      respjson.Field
		Results     respjson.Field
		StartTime   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageVectorStoresResponseData) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationUsageVectorStoresResponseData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// AdminOrganizationUsageVectorStoresResponseDataResultUnion contains all possible
// properties and values from
// [AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageCompletionsResult],
// [AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageEmbeddingsResult],
// [AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageModerationsResult],
// [AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageImagesResult],
// [AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageAudioSpeechesResult],
// [AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageAudioTranscriptionsResult],
// [AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageVectorStoresResult],
// [AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageCodeInterpreterSessionsResult],
// [AdminOrganizationUsageVectorStoresResponseDataResultOrganizationCostsResult].
//
// Use the [AdminOrganizationUsageVectorStoresResponseDataResultUnion.AsAny] method
// to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type AdminOrganizationUsageVectorStoresResponseDataResultUnion struct {
	InputTokens      int64 `json:"input_tokens"`
	NumModelRequests int64 `json:"num_model_requests"`
	// Any of "organization.usage.completions.result",
	// "organization.usage.embeddings.result", "organization.usage.moderations.result",
	// "organization.usage.images.result", "organization.usage.audio_speeches.result",
	// "organization.usage.audio_transcriptions.result",
	// "organization.usage.vector_stores.result",
	// "organization.usage.code_interpreter_sessions.result",
	// "organization.costs.result".
	Object string `json:"object"`
	// This field is from variant
	// [AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageCompletionsResult].
	OutputTokens int64  `json:"output_tokens"`
	APIKeyID     string `json:"api_key_id"`
	// This field is from variant
	// [AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageCompletionsResult].
	Batch bool `json:"batch"`
	// This field is from variant
	// [AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageCompletionsResult].
	InputAudioTokens int64 `json:"input_audio_tokens"`
	// This field is from variant
	// [AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageCompletionsResult].
	InputCachedTokens int64  `json:"input_cached_tokens"`
	Model             string `json:"model"`
	// This field is from variant
	// [AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageCompletionsResult].
	OutputAudioTokens int64  `json:"output_audio_tokens"`
	ProjectID         string `json:"project_id"`
	// This field is from variant
	// [AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageCompletionsResult].
	ServiceTier string `json:"service_tier"`
	UserID      string `json:"user_id"`
	// This field is from variant
	// [AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageImagesResult].
	Images int64 `json:"images"`
	// This field is from variant
	// [AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageImagesResult].
	Size string `json:"size"`
	// This field is from variant
	// [AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageImagesResult].
	Source string `json:"source"`
	// This field is from variant
	// [AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageAudioSpeechesResult].
	Characters int64 `json:"characters"`
	// This field is from variant
	// [AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageAudioTranscriptionsResult].
	Seconds int64 `json:"seconds"`
	// This field is from variant
	// [AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageVectorStoresResult].
	UsageBytes int64 `json:"usage_bytes"`
	// This field is from variant
	// [AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageCodeInterpreterSessionsResult].
	NumSessions int64 `json:"num_sessions"`
	// This field is from variant
	// [AdminOrganizationUsageVectorStoresResponseDataResultOrganizationCostsResult].
	Amount AdminOrganizationUsageVectorStoresResponseDataResultOrganizationCostsResultAmount `json:"amount"`
	// This field is from variant
	// [AdminOrganizationUsageVectorStoresResponseDataResultOrganizationCostsResult].
	LineItem string `json:"line_item"`
	JSON     struct {
		InputTokens       respjson.Field
		NumModelRequests  respjson.Field
		Object            respjson.Field
		OutputTokens      respjson.Field
		APIKeyID          respjson.Field
		Batch             respjson.Field
		InputAudioTokens  respjson.Field
		InputCachedTokens respjson.Field
		Model             respjson.Field
		OutputAudioTokens respjson.Field
		ProjectID         respjson.Field
		ServiceTier       respjson.Field
		UserID            respjson.Field
		Images            respjson.Field
		Size              respjson.Field
		Source            respjson.Field
		Characters        respjson.Field
		Seconds           respjson.Field
		UsageBytes        respjson.Field
		NumSessions       respjson.Field
		Amount            respjson.Field
		LineItem          respjson.Field
		raw               string
	} `json:"-"`
}

// anyAdminOrganizationUsageVectorStoresResponseDataResult is implemented by each
// variant of [AdminOrganizationUsageVectorStoresResponseDataResultUnion] to add
// type safety for the return type of
// [AdminOrganizationUsageVectorStoresResponseDataResultUnion.AsAny]
type anyAdminOrganizationUsageVectorStoresResponseDataResult interface {
	implAdminOrganizationUsageVectorStoresResponseDataResultUnion()
}

func (AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageCompletionsResult) implAdminOrganizationUsageVectorStoresResponseDataResultUnion() {
}
func (AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageEmbeddingsResult) implAdminOrganizationUsageVectorStoresResponseDataResultUnion() {
}
func (AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageModerationsResult) implAdminOrganizationUsageVectorStoresResponseDataResultUnion() {
}
func (AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageImagesResult) implAdminOrganizationUsageVectorStoresResponseDataResultUnion() {
}
func (AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageAudioSpeechesResult) implAdminOrganizationUsageVectorStoresResponseDataResultUnion() {
}
func (AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageAudioTranscriptionsResult) implAdminOrganizationUsageVectorStoresResponseDataResultUnion() {
}
func (AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageVectorStoresResult) implAdminOrganizationUsageVectorStoresResponseDataResultUnion() {
}
func (AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageCodeInterpreterSessionsResult) implAdminOrganizationUsageVectorStoresResponseDataResultUnion() {
}
func (AdminOrganizationUsageVectorStoresResponseDataResultOrganizationCostsResult) implAdminOrganizationUsageVectorStoresResponseDataResultUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := AdminOrganizationUsageVectorStoresResponseDataResultUnion.AsAny().(type) {
//	case openai.AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageCompletionsResult:
//	case openai.AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageEmbeddingsResult:
//	case openai.AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageModerationsResult:
//	case openai.AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageImagesResult:
//	case openai.AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageAudioSpeechesResult:
//	case openai.AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageAudioTranscriptionsResult:
//	case openai.AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageVectorStoresResult:
//	case openai.AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageCodeInterpreterSessionsResult:
//	case openai.AdminOrganizationUsageVectorStoresResponseDataResultOrganizationCostsResult:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u AdminOrganizationUsageVectorStoresResponseDataResultUnion) AsAny() anyAdminOrganizationUsageVectorStoresResponseDataResult {
	switch u.Object {
	case "organization.usage.completions.result":
		return u.AsOrganizationUsageCompletionsResult()
	case "organization.usage.embeddings.result":
		return u.AsOrganizationUsageEmbeddingsResult()
	case "organization.usage.moderations.result":
		return u.AsOrganizationUsageModerationsResult()
	case "organization.usage.images.result":
		return u.AsOrganizationUsageImagesResult()
	case "organization.usage.audio_speeches.result":
		return u.AsOrganizationUsageAudioSpeechesResult()
	case "organization.usage.audio_transcriptions.result":
		return u.AsOrganizationUsageAudioTranscriptionsResult()
	case "organization.usage.vector_stores.result":
		return u.AsOrganizationUsageVectorStoresResult()
	case "organization.usage.code_interpreter_sessions.result":
		return u.AsOrganizationUsageCodeInterpreterSessionsResult()
	case "organization.costs.result":
		return u.AsOrganizationCostsResult()
	}
	return nil
}

func (u AdminOrganizationUsageVectorStoresResponseDataResultUnion) AsOrganizationUsageCompletionsResult() (v AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageCompletionsResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageVectorStoresResponseDataResultUnion) AsOrganizationUsageEmbeddingsResult() (v AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageEmbeddingsResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageVectorStoresResponseDataResultUnion) AsOrganizationUsageModerationsResult() (v AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageModerationsResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageVectorStoresResponseDataResultUnion) AsOrganizationUsageImagesResult() (v AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageImagesResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageVectorStoresResponseDataResultUnion) AsOrganizationUsageAudioSpeechesResult() (v AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageAudioSpeechesResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageVectorStoresResponseDataResultUnion) AsOrganizationUsageAudioTranscriptionsResult() (v AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageAudioTranscriptionsResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageVectorStoresResponseDataResultUnion) AsOrganizationUsageVectorStoresResult() (v AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageVectorStoresResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageVectorStoresResponseDataResultUnion) AsOrganizationUsageCodeInterpreterSessionsResult() (v AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageCodeInterpreterSessionsResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AdminOrganizationUsageVectorStoresResponseDataResultUnion) AsOrganizationCostsResult() (v AdminOrganizationUsageVectorStoresResponseDataResultOrganizationCostsResult) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u AdminOrganizationUsageVectorStoresResponseDataResultUnion) RawJSON() string {
	return u.JSON.raw
}

func (r *AdminOrganizationUsageVectorStoresResponseDataResultUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated completions usage details of the specific time bucket.
type AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageCompletionsResult struct {
	// The aggregated number of text input tokens used, including cached tokens. For
	// customers subscribe to scale tier, this includes scale tier tokens.
	InputTokens int64 `json:"input_tokens" api:"required"`
	// The count of requests made to the model.
	NumModelRequests int64                                       `json:"num_model_requests" api:"required"`
	Object           constant.OrganizationUsageCompletionsResult `json:"object" default:"organization.usage.completions.result"`
	// The aggregated number of text output tokens used. For customers subscribe to
	// scale tier, this includes scale tier tokens.
	OutputTokens int64 `json:"output_tokens" api:"required"`
	// When `group_by=api_key_id`, this field provides the API key ID of the grouped
	// usage result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=batch`, this field tells whether the grouped usage result is
	// batch or not.
	Batch bool `json:"batch" api:"nullable"`
	// The aggregated number of audio input tokens used, including cached tokens.
	InputAudioTokens int64 `json:"input_audio_tokens"`
	// The aggregated number of text input tokens that has been cached from previous
	// requests. For customers subscribe to scale tier, this includes scale tier
	// tokens.
	InputCachedTokens int64 `json:"input_cached_tokens"`
	// When `group_by=model`, this field provides the model name of the grouped usage
	// result.
	Model string `json:"model" api:"nullable"`
	// The aggregated number of audio output tokens used.
	OutputAudioTokens int64 `json:"output_audio_tokens"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// When `group_by=service_tier`, this field provides the service tier of the
	// grouped usage result.
	ServiceTier string `json:"service_tier" api:"nullable"`
	// When `group_by=user_id`, this field provides the user ID of the grouped usage
	// result.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		InputTokens       respjson.Field
		NumModelRequests  respjson.Field
		Object            respjson.Field
		OutputTokens      respjson.Field
		APIKeyID          respjson.Field
		Batch             respjson.Field
		InputAudioTokens  respjson.Field
		InputCachedTokens respjson.Field
		Model             respjson.Field
		OutputAudioTokens respjson.Field
		ProjectID         respjson.Field
		ServiceTier       respjson.Field
		UserID            respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageCompletionsResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageCompletionsResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated embeddings usage details of the specific time bucket.
type AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageEmbeddingsResult struct {
	// The aggregated number of input tokens used.
	InputTokens int64 `json:"input_tokens" api:"required"`
	// The count of requests made to the model.
	NumModelRequests int64                                      `json:"num_model_requests" api:"required"`
	Object           constant.OrganizationUsageEmbeddingsResult `json:"object" default:"organization.usage.embeddings.result"`
	// When `group_by=api_key_id`, this field provides the API key ID of the grouped
	// usage result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=model`, this field provides the model name of the grouped usage
	// result.
	Model string `json:"model" api:"nullable"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// When `group_by=user_id`, this field provides the user ID of the grouped usage
	// result.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		InputTokens      respjson.Field
		NumModelRequests respjson.Field
		Object           respjson.Field
		APIKeyID         respjson.Field
		Model            respjson.Field
		ProjectID        respjson.Field
		UserID           respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageEmbeddingsResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageEmbeddingsResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated moderations usage details of the specific time bucket.
type AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageModerationsResult struct {
	// The aggregated number of input tokens used.
	InputTokens int64 `json:"input_tokens" api:"required"`
	// The count of requests made to the model.
	NumModelRequests int64                                       `json:"num_model_requests" api:"required"`
	Object           constant.OrganizationUsageModerationsResult `json:"object" default:"organization.usage.moderations.result"`
	// When `group_by=api_key_id`, this field provides the API key ID of the grouped
	// usage result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=model`, this field provides the model name of the grouped usage
	// result.
	Model string `json:"model" api:"nullable"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// When `group_by=user_id`, this field provides the user ID of the grouped usage
	// result.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		InputTokens      respjson.Field
		NumModelRequests respjson.Field
		Object           respjson.Field
		APIKeyID         respjson.Field
		Model            respjson.Field
		ProjectID        respjson.Field
		UserID           respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageModerationsResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageModerationsResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated images usage details of the specific time bucket.
type AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageImagesResult struct {
	// The number of images processed.
	Images int64 `json:"images" api:"required"`
	// The count of requests made to the model.
	NumModelRequests int64                                  `json:"num_model_requests" api:"required"`
	Object           constant.OrganizationUsageImagesResult `json:"object" default:"organization.usage.images.result"`
	// When `group_by=api_key_id`, this field provides the API key ID of the grouped
	// usage result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=model`, this field provides the model name of the grouped usage
	// result.
	Model string `json:"model" api:"nullable"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// When `group_by=size`, this field provides the image size of the grouped usage
	// result.
	Size string `json:"size" api:"nullable"`
	// When `group_by=source`, this field provides the source of the grouped usage
	// result, possible values are `image.generation`, `image.edit`, `image.variation`.
	Source string `json:"source" api:"nullable"`
	// When `group_by=user_id`, this field provides the user ID of the grouped usage
	// result.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Images           respjson.Field
		NumModelRequests respjson.Field
		Object           respjson.Field
		APIKeyID         respjson.Field
		Model            respjson.Field
		ProjectID        respjson.Field
		Size             respjson.Field
		Source           respjson.Field
		UserID           respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageImagesResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageImagesResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated audio speeches usage details of the specific time bucket.
type AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageAudioSpeechesResult struct {
	// The number of characters processed.
	Characters int64 `json:"characters" api:"required"`
	// The count of requests made to the model.
	NumModelRequests int64                                         `json:"num_model_requests" api:"required"`
	Object           constant.OrganizationUsageAudioSpeechesResult `json:"object" default:"organization.usage.audio_speeches.result"`
	// When `group_by=api_key_id`, this field provides the API key ID of the grouped
	// usage result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=model`, this field provides the model name of the grouped usage
	// result.
	Model string `json:"model" api:"nullable"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// When `group_by=user_id`, this field provides the user ID of the grouped usage
	// result.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Characters       respjson.Field
		NumModelRequests respjson.Field
		Object           respjson.Field
		APIKeyID         respjson.Field
		Model            respjson.Field
		ProjectID        respjson.Field
		UserID           respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageAudioSpeechesResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageAudioSpeechesResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated audio transcriptions usage details of the specific time bucket.
type AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageAudioTranscriptionsResult struct {
	// The count of requests made to the model.
	NumModelRequests int64                                               `json:"num_model_requests" api:"required"`
	Object           constant.OrganizationUsageAudioTranscriptionsResult `json:"object" default:"organization.usage.audio_transcriptions.result"`
	// The number of seconds processed.
	Seconds int64 `json:"seconds" api:"required"`
	// When `group_by=api_key_id`, this field provides the API key ID of the grouped
	// usage result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=model`, this field provides the model name of the grouped usage
	// result.
	Model string `json:"model" api:"nullable"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// When `group_by=user_id`, this field provides the user ID of the grouped usage
	// result.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		NumModelRequests respjson.Field
		Object           respjson.Field
		Seconds          respjson.Field
		APIKeyID         respjson.Field
		Model            respjson.Field
		ProjectID        respjson.Field
		UserID           respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageAudioTranscriptionsResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageAudioTranscriptionsResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated vector stores usage details of the specific time bucket.
type AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageVectorStoresResult struct {
	Object constant.OrganizationUsageVectorStoresResult `json:"object" default:"organization.usage.vector_stores.result"`
	// The vector stores usage in bytes.
	UsageBytes int64 `json:"usage_bytes" api:"required"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Object      respjson.Field
		UsageBytes  respjson.Field
		ProjectID   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageVectorStoresResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageVectorStoresResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated code interpreter sessions usage details of the specific time
// bucket.
type AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageCodeInterpreterSessionsResult struct {
	// The number of code interpreter sessions.
	NumSessions int64                                                   `json:"num_sessions" api:"required"`
	Object      constant.OrganizationUsageCodeInterpreterSessionsResult `json:"object" default:"organization.usage.code_interpreter_sessions.result"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// usage result.
	ProjectID string `json:"project_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		NumSessions respjson.Field
		Object      respjson.Field
		ProjectID   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageCodeInterpreterSessionsResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageVectorStoresResponseDataResultOrganizationUsageCodeInterpreterSessionsResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The aggregated costs details of the specific time bucket.
type AdminOrganizationUsageVectorStoresResponseDataResultOrganizationCostsResult struct {
	Object constant.OrganizationCostsResult `json:"object" default:"organization.costs.result"`
	// The monetary value in its associated currency.
	Amount AdminOrganizationUsageVectorStoresResponseDataResultOrganizationCostsResultAmount `json:"amount"`
	// When `group_by=api_key_id`, this field provides the API Key ID of the grouped
	// costs result.
	APIKeyID string `json:"api_key_id" api:"nullable"`
	// When `group_by=line_item`, this field provides the line item of the grouped
	// costs result.
	LineItem string `json:"line_item" api:"nullable"`
	// When `group_by=project_id`, this field provides the project ID of the grouped
	// costs result.
	ProjectID string `json:"project_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Object      respjson.Field
		Amount      respjson.Field
		APIKeyID    respjson.Field
		LineItem    respjson.Field
		ProjectID   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageVectorStoresResponseDataResultOrganizationCostsResult) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageVectorStoresResponseDataResultOrganizationCostsResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The monetary value in its associated currency.
type AdminOrganizationUsageVectorStoresResponseDataResultOrganizationCostsResultAmount struct {
	// Lowercase ISO-4217 currency e.g. "usd"
	Currency string `json:"currency"`
	// The numeric value of the cost.
	Value float64 `json:"value"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Currency    respjson.Field
		Value       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationUsageVectorStoresResponseDataResultOrganizationCostsResultAmount) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationUsageVectorStoresResponseDataResultOrganizationCostsResultAmount) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationUsageAudioSpeechesParams struct {
	// Start time (Unix seconds) of the query time range, inclusive.
	StartTime int64 `query:"start_time" api:"required" json:"-"`
	// End time (Unix seconds) of the query time range, exclusive.
	EndTime param.Opt[int64] `query:"end_time,omitzero" json:"-"`
	// Specifies the number of buckets to return.
	//
	// - `bucket_width=1d`: default: 7, max: 31
	// - `bucket_width=1h`: default: 24, max: 168
	// - `bucket_width=1m`: default: 60, max: 1440
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// A cursor for use in pagination. Corresponding to the `next_page` field from the
	// previous response.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Return only usage for these API keys.
	APIKeyIDs []string `query:"api_key_ids,omitzero" json:"-"`
	// Width of each time bucket in response. Currently `1m`, `1h` and `1d` are
	// supported, default to `1d`.
	//
	// Any of "1m", "1h", "1d".
	BucketWidth AdminOrganizationUsageAudioSpeechesParamsBucketWidth `query:"bucket_width,omitzero" json:"-"`
	// Group the usage data by the specified fields. Support fields include
	// `project_id`, `user_id`, `api_key_id`, `model` or any combination of them.
	//
	// Any of "project_id", "user_id", "api_key_id", "model".
	GroupBy []string `query:"group_by,omitzero" json:"-"`
	// Return only usage for these models.
	Models []string `query:"models,omitzero" json:"-"`
	// Return only usage for these projects.
	ProjectIDs []string `query:"project_ids,omitzero" json:"-"`
	// Return only usage for these users.
	UserIDs []string `query:"user_ids,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [AdminOrganizationUsageAudioSpeechesParams]'s query
// parameters as `url.Values`.
func (r AdminOrganizationUsageAudioSpeechesParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Width of each time bucket in response. Currently `1m`, `1h` and `1d` are
// supported, default to `1d`.
type AdminOrganizationUsageAudioSpeechesParamsBucketWidth string

const (
	AdminOrganizationUsageAudioSpeechesParamsBucketWidth1m AdminOrganizationUsageAudioSpeechesParamsBucketWidth = "1m"
	AdminOrganizationUsageAudioSpeechesParamsBucketWidth1h AdminOrganizationUsageAudioSpeechesParamsBucketWidth = "1h"
	AdminOrganizationUsageAudioSpeechesParamsBucketWidth1d AdminOrganizationUsageAudioSpeechesParamsBucketWidth = "1d"
)

type AdminOrganizationUsageAudioTranscriptionsParams struct {
	// Start time (Unix seconds) of the query time range, inclusive.
	StartTime int64 `query:"start_time" api:"required" json:"-"`
	// End time (Unix seconds) of the query time range, exclusive.
	EndTime param.Opt[int64] `query:"end_time,omitzero" json:"-"`
	// Specifies the number of buckets to return.
	//
	// - `bucket_width=1d`: default: 7, max: 31
	// - `bucket_width=1h`: default: 24, max: 168
	// - `bucket_width=1m`: default: 60, max: 1440
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// A cursor for use in pagination. Corresponding to the `next_page` field from the
	// previous response.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Return only usage for these API keys.
	APIKeyIDs []string `query:"api_key_ids,omitzero" json:"-"`
	// Width of each time bucket in response. Currently `1m`, `1h` and `1d` are
	// supported, default to `1d`.
	//
	// Any of "1m", "1h", "1d".
	BucketWidth AdminOrganizationUsageAudioTranscriptionsParamsBucketWidth `query:"bucket_width,omitzero" json:"-"`
	// Group the usage data by the specified fields. Support fields include
	// `project_id`, `user_id`, `api_key_id`, `model` or any combination of them.
	//
	// Any of "project_id", "user_id", "api_key_id", "model".
	GroupBy []string `query:"group_by,omitzero" json:"-"`
	// Return only usage for these models.
	Models []string `query:"models,omitzero" json:"-"`
	// Return only usage for these projects.
	ProjectIDs []string `query:"project_ids,omitzero" json:"-"`
	// Return only usage for these users.
	UserIDs []string `query:"user_ids,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [AdminOrganizationUsageAudioTranscriptionsParams]'s query
// parameters as `url.Values`.
func (r AdminOrganizationUsageAudioTranscriptionsParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Width of each time bucket in response. Currently `1m`, `1h` and `1d` are
// supported, default to `1d`.
type AdminOrganizationUsageAudioTranscriptionsParamsBucketWidth string

const (
	AdminOrganizationUsageAudioTranscriptionsParamsBucketWidth1m AdminOrganizationUsageAudioTranscriptionsParamsBucketWidth = "1m"
	AdminOrganizationUsageAudioTranscriptionsParamsBucketWidth1h AdminOrganizationUsageAudioTranscriptionsParamsBucketWidth = "1h"
	AdminOrganizationUsageAudioTranscriptionsParamsBucketWidth1d AdminOrganizationUsageAudioTranscriptionsParamsBucketWidth = "1d"
)

type AdminOrganizationUsageCodeInterpreterSessionsParams struct {
	// Start time (Unix seconds) of the query time range, inclusive.
	StartTime int64 `query:"start_time" api:"required" json:"-"`
	// End time (Unix seconds) of the query time range, exclusive.
	EndTime param.Opt[int64] `query:"end_time,omitzero" json:"-"`
	// Specifies the number of buckets to return.
	//
	// - `bucket_width=1d`: default: 7, max: 31
	// - `bucket_width=1h`: default: 24, max: 168
	// - `bucket_width=1m`: default: 60, max: 1440
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// A cursor for use in pagination. Corresponding to the `next_page` field from the
	// previous response.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Width of each time bucket in response. Currently `1m`, `1h` and `1d` are
	// supported, default to `1d`.
	//
	// Any of "1m", "1h", "1d".
	BucketWidth AdminOrganizationUsageCodeInterpreterSessionsParamsBucketWidth `query:"bucket_width,omitzero" json:"-"`
	// Group the usage data by the specified fields. Support fields include
	// `project_id`.
	//
	// Any of "project_id".
	GroupBy []string `query:"group_by,omitzero" json:"-"`
	// Return only usage for these projects.
	ProjectIDs []string `query:"project_ids,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [AdminOrganizationUsageCodeInterpreterSessionsParams]'s
// query parameters as `url.Values`.
func (r AdminOrganizationUsageCodeInterpreterSessionsParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Width of each time bucket in response. Currently `1m`, `1h` and `1d` are
// supported, default to `1d`.
type AdminOrganizationUsageCodeInterpreterSessionsParamsBucketWidth string

const (
	AdminOrganizationUsageCodeInterpreterSessionsParamsBucketWidth1m AdminOrganizationUsageCodeInterpreterSessionsParamsBucketWidth = "1m"
	AdminOrganizationUsageCodeInterpreterSessionsParamsBucketWidth1h AdminOrganizationUsageCodeInterpreterSessionsParamsBucketWidth = "1h"
	AdminOrganizationUsageCodeInterpreterSessionsParamsBucketWidth1d AdminOrganizationUsageCodeInterpreterSessionsParamsBucketWidth = "1d"
)

type AdminOrganizationUsageCompletionsParams struct {
	// Start time (Unix seconds) of the query time range, inclusive.
	StartTime int64 `query:"start_time" api:"required" json:"-"`
	// If `true`, return batch jobs only. If `false`, return non-batch jobs only. By
	// default, return both.
	Batch param.Opt[bool] `query:"batch,omitzero" json:"-"`
	// End time (Unix seconds) of the query time range, exclusive.
	EndTime param.Opt[int64] `query:"end_time,omitzero" json:"-"`
	// Specifies the number of buckets to return.
	//
	// - `bucket_width=1d`: default: 7, max: 31
	// - `bucket_width=1h`: default: 24, max: 168
	// - `bucket_width=1m`: default: 60, max: 1440
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// A cursor for use in pagination. Corresponding to the `next_page` field from the
	// previous response.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Return only usage for these API keys.
	APIKeyIDs []string `query:"api_key_ids,omitzero" json:"-"`
	// Width of each time bucket in response. Currently `1m`, `1h` and `1d` are
	// supported, default to `1d`.
	//
	// Any of "1m", "1h", "1d".
	BucketWidth AdminOrganizationUsageCompletionsParamsBucketWidth `query:"bucket_width,omitzero" json:"-"`
	// Group the usage data by the specified fields. Support fields include
	// `project_id`, `user_id`, `api_key_id`, `model`, `batch`, `service_tier` or any
	// combination of them.
	//
	// Any of "project_id", "user_id", "api_key_id", "model", "batch", "service_tier".
	GroupBy []string `query:"group_by,omitzero" json:"-"`
	// Return only usage for these models.
	Models []string `query:"models,omitzero" json:"-"`
	// Return only usage for these projects.
	ProjectIDs []string `query:"project_ids,omitzero" json:"-"`
	// Return only usage for these users.
	UserIDs []string `query:"user_ids,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [AdminOrganizationUsageCompletionsParams]'s query parameters
// as `url.Values`.
func (r AdminOrganizationUsageCompletionsParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Width of each time bucket in response. Currently `1m`, `1h` and `1d` are
// supported, default to `1d`.
type AdminOrganizationUsageCompletionsParamsBucketWidth string

const (
	AdminOrganizationUsageCompletionsParamsBucketWidth1m AdminOrganizationUsageCompletionsParamsBucketWidth = "1m"
	AdminOrganizationUsageCompletionsParamsBucketWidth1h AdminOrganizationUsageCompletionsParamsBucketWidth = "1h"
	AdminOrganizationUsageCompletionsParamsBucketWidth1d AdminOrganizationUsageCompletionsParamsBucketWidth = "1d"
)

type AdminOrganizationUsageCostsParams struct {
	// Start time (Unix seconds) of the query time range, inclusive.
	StartTime int64 `query:"start_time" api:"required" json:"-"`
	// End time (Unix seconds) of the query time range, exclusive.
	EndTime param.Opt[int64] `query:"end_time,omitzero" json:"-"`
	// A limit on the number of buckets to be returned. Limit can range between 1 and
	// 180, and the default is 7.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// A cursor for use in pagination. Corresponding to the `next_page` field from the
	// previous response.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Return only costs for these API keys.
	APIKeyIDs []string `query:"api_key_ids,omitzero" json:"-"`
	// Width of each time bucket in response. Currently only `1d` is supported, default
	// to `1d`.
	//
	// Any of "1d".
	BucketWidth AdminOrganizationUsageCostsParamsBucketWidth `query:"bucket_width,omitzero" json:"-"`
	// Group the costs by the specified fields. Support fields include `project_id`,
	// `line_item`, `api_key_id` and any combination of them.
	//
	// Any of "project_id", "line_item", "api_key_id".
	GroupBy []string `query:"group_by,omitzero" json:"-"`
	// Return only costs for these projects.
	ProjectIDs []string `query:"project_ids,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [AdminOrganizationUsageCostsParams]'s query parameters as
// `url.Values`.
func (r AdminOrganizationUsageCostsParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Width of each time bucket in response. Currently only `1d` is supported, default
// to `1d`.
type AdminOrganizationUsageCostsParamsBucketWidth string

const (
	AdminOrganizationUsageCostsParamsBucketWidth1d AdminOrganizationUsageCostsParamsBucketWidth = "1d"
)

type AdminOrganizationUsageEmbeddingsParams struct {
	// Start time (Unix seconds) of the query time range, inclusive.
	StartTime int64 `query:"start_time" api:"required" json:"-"`
	// End time (Unix seconds) of the query time range, exclusive.
	EndTime param.Opt[int64] `query:"end_time,omitzero" json:"-"`
	// Specifies the number of buckets to return.
	//
	// - `bucket_width=1d`: default: 7, max: 31
	// - `bucket_width=1h`: default: 24, max: 168
	// - `bucket_width=1m`: default: 60, max: 1440
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// A cursor for use in pagination. Corresponding to the `next_page` field from the
	// previous response.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Return only usage for these API keys.
	APIKeyIDs []string `query:"api_key_ids,omitzero" json:"-"`
	// Width of each time bucket in response. Currently `1m`, `1h` and `1d` are
	// supported, default to `1d`.
	//
	// Any of "1m", "1h", "1d".
	BucketWidth AdminOrganizationUsageEmbeddingsParamsBucketWidth `query:"bucket_width,omitzero" json:"-"`
	// Group the usage data by the specified fields. Support fields include
	// `project_id`, `user_id`, `api_key_id`, `model` or any combination of them.
	//
	// Any of "project_id", "user_id", "api_key_id", "model".
	GroupBy []string `query:"group_by,omitzero" json:"-"`
	// Return only usage for these models.
	Models []string `query:"models,omitzero" json:"-"`
	// Return only usage for these projects.
	ProjectIDs []string `query:"project_ids,omitzero" json:"-"`
	// Return only usage for these users.
	UserIDs []string `query:"user_ids,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [AdminOrganizationUsageEmbeddingsParams]'s query parameters
// as `url.Values`.
func (r AdminOrganizationUsageEmbeddingsParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Width of each time bucket in response. Currently `1m`, `1h` and `1d` are
// supported, default to `1d`.
type AdminOrganizationUsageEmbeddingsParamsBucketWidth string

const (
	AdminOrganizationUsageEmbeddingsParamsBucketWidth1m AdminOrganizationUsageEmbeddingsParamsBucketWidth = "1m"
	AdminOrganizationUsageEmbeddingsParamsBucketWidth1h AdminOrganizationUsageEmbeddingsParamsBucketWidth = "1h"
	AdminOrganizationUsageEmbeddingsParamsBucketWidth1d AdminOrganizationUsageEmbeddingsParamsBucketWidth = "1d"
)

type AdminOrganizationUsageImagesParams struct {
	// Start time (Unix seconds) of the query time range, inclusive.
	StartTime int64 `query:"start_time" api:"required" json:"-"`
	// End time (Unix seconds) of the query time range, exclusive.
	EndTime param.Opt[int64] `query:"end_time,omitzero" json:"-"`
	// Specifies the number of buckets to return.
	//
	// - `bucket_width=1d`: default: 7, max: 31
	// - `bucket_width=1h`: default: 24, max: 168
	// - `bucket_width=1m`: default: 60, max: 1440
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// A cursor for use in pagination. Corresponding to the `next_page` field from the
	// previous response.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Return only usage for these API keys.
	APIKeyIDs []string `query:"api_key_ids,omitzero" json:"-"`
	// Width of each time bucket in response. Currently `1m`, `1h` and `1d` are
	// supported, default to `1d`.
	//
	// Any of "1m", "1h", "1d".
	BucketWidth AdminOrganizationUsageImagesParamsBucketWidth `query:"bucket_width,omitzero" json:"-"`
	// Group the usage data by the specified fields. Support fields include
	// `project_id`, `user_id`, `api_key_id`, `model`, `size`, `source` or any
	// combination of them.
	//
	// Any of "project_id", "user_id", "api_key_id", "model", "size", "source".
	GroupBy []string `query:"group_by,omitzero" json:"-"`
	// Return only usage for these models.
	Models []string `query:"models,omitzero" json:"-"`
	// Return only usage for these projects.
	ProjectIDs []string `query:"project_ids,omitzero" json:"-"`
	// Return only usages for these image sizes. Possible values are `256x256`,
	// `512x512`, `1024x1024`, `1792x1792`, `1024x1792` or any combination of them.
	//
	// Any of "256x256", "512x512", "1024x1024", "1792x1792", "1024x1792".
	Sizes []string `query:"sizes,omitzero" json:"-"`
	// Return only usages for these sources. Possible values are `image.generation`,
	// `image.edit`, `image.variation` or any combination of them.
	//
	// Any of "image.generation", "image.edit", "image.variation".
	Sources []string `query:"sources,omitzero" json:"-"`
	// Return only usage for these users.
	UserIDs []string `query:"user_ids,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [AdminOrganizationUsageImagesParams]'s query parameters as
// `url.Values`.
func (r AdminOrganizationUsageImagesParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Width of each time bucket in response. Currently `1m`, `1h` and `1d` are
// supported, default to `1d`.
type AdminOrganizationUsageImagesParamsBucketWidth string

const (
	AdminOrganizationUsageImagesParamsBucketWidth1m AdminOrganizationUsageImagesParamsBucketWidth = "1m"
	AdminOrganizationUsageImagesParamsBucketWidth1h AdminOrganizationUsageImagesParamsBucketWidth = "1h"
	AdminOrganizationUsageImagesParamsBucketWidth1d AdminOrganizationUsageImagesParamsBucketWidth = "1d"
)

type AdminOrganizationUsageModerationsParams struct {
	// Start time (Unix seconds) of the query time range, inclusive.
	StartTime int64 `query:"start_time" api:"required" json:"-"`
	// End time (Unix seconds) of the query time range, exclusive.
	EndTime param.Opt[int64] `query:"end_time,omitzero" json:"-"`
	// Specifies the number of buckets to return.
	//
	// - `bucket_width=1d`: default: 7, max: 31
	// - `bucket_width=1h`: default: 24, max: 168
	// - `bucket_width=1m`: default: 60, max: 1440
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// A cursor for use in pagination. Corresponding to the `next_page` field from the
	// previous response.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Return only usage for these API keys.
	APIKeyIDs []string `query:"api_key_ids,omitzero" json:"-"`
	// Width of each time bucket in response. Currently `1m`, `1h` and `1d` are
	// supported, default to `1d`.
	//
	// Any of "1m", "1h", "1d".
	BucketWidth AdminOrganizationUsageModerationsParamsBucketWidth `query:"bucket_width,omitzero" json:"-"`
	// Group the usage data by the specified fields. Support fields include
	// `project_id`, `user_id`, `api_key_id`, `model` or any combination of them.
	//
	// Any of "project_id", "user_id", "api_key_id", "model".
	GroupBy []string `query:"group_by,omitzero" json:"-"`
	// Return only usage for these models.
	Models []string `query:"models,omitzero" json:"-"`
	// Return only usage for these projects.
	ProjectIDs []string `query:"project_ids,omitzero" json:"-"`
	// Return only usage for these users.
	UserIDs []string `query:"user_ids,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [AdminOrganizationUsageModerationsParams]'s query parameters
// as `url.Values`.
func (r AdminOrganizationUsageModerationsParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Width of each time bucket in response. Currently `1m`, `1h` and `1d` are
// supported, default to `1d`.
type AdminOrganizationUsageModerationsParamsBucketWidth string

const (
	AdminOrganizationUsageModerationsParamsBucketWidth1m AdminOrganizationUsageModerationsParamsBucketWidth = "1m"
	AdminOrganizationUsageModerationsParamsBucketWidth1h AdminOrganizationUsageModerationsParamsBucketWidth = "1h"
	AdminOrganizationUsageModerationsParamsBucketWidth1d AdminOrganizationUsageModerationsParamsBucketWidth = "1d"
)

type AdminOrganizationUsageVectorStoresParams struct {
	// Start time (Unix seconds) of the query time range, inclusive.
	StartTime int64 `query:"start_time" api:"required" json:"-"`
	// End time (Unix seconds) of the query time range, exclusive.
	EndTime param.Opt[int64] `query:"end_time,omitzero" json:"-"`
	// Specifies the number of buckets to return.
	//
	// - `bucket_width=1d`: default: 7, max: 31
	// - `bucket_width=1h`: default: 24, max: 168
	// - `bucket_width=1m`: default: 60, max: 1440
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// A cursor for use in pagination. Corresponding to the `next_page` field from the
	// previous response.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Width of each time bucket in response. Currently `1m`, `1h` and `1d` are
	// supported, default to `1d`.
	//
	// Any of "1m", "1h", "1d".
	BucketWidth AdminOrganizationUsageVectorStoresParamsBucketWidth `query:"bucket_width,omitzero" json:"-"`
	// Group the usage data by the specified fields. Support fields include
	// `project_id`.
	//
	// Any of "project_id".
	GroupBy []string `query:"group_by,omitzero" json:"-"`
	// Return only usage for these projects.
	ProjectIDs []string `query:"project_ids,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [AdminOrganizationUsageVectorStoresParams]'s query
// parameters as `url.Values`.
func (r AdminOrganizationUsageVectorStoresParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Width of each time bucket in response. Currently `1m`, `1h` and `1d` are
// supported, default to `1d`.
type AdminOrganizationUsageVectorStoresParamsBucketWidth string

const (
	AdminOrganizationUsageVectorStoresParamsBucketWidth1m AdminOrganizationUsageVectorStoresParamsBucketWidth = "1m"
	AdminOrganizationUsageVectorStoresParamsBucketWidth1h AdminOrganizationUsageVectorStoresParamsBucketWidth = "1h"
	AdminOrganizationUsageVectorStoresParamsBucketWidth1d AdminOrganizationUsageVectorStoresParamsBucketWidth = "1d"
)

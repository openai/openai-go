// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/openai/openai-go/internal/apijson"
	"github.com/openai/openai-go/internal/apiquery"
	"github.com/openai/openai-go/internal/requestconfig"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/pagination"
	"github.com/openai/openai-go/packages/param"
	"github.com/openai/openai-go/packages/resp"
	"github.com/openai/openai-go/shared/constant"
)

// EvalRunOutputItemService contains methods and other services that help with
// interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewEvalRunOutputItemService] method instead.
type EvalRunOutputItemService struct {
	Options []option.RequestOption
}

// NewEvalRunOutputItemService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewEvalRunOutputItemService(opts ...option.RequestOption) (r EvalRunOutputItemService) {
	r = EvalRunOutputItemService{}
	r.Options = opts
	return
}

// Get an evaluation run output item by ID.
func (r *EvalRunOutputItemService) Get(ctx context.Context, evalID string, runID string, outputItemID string, opts ...option.RequestOption) (res *EvalRunOutputItemGetResponse, err error) {
	opts = append(r.Options[:], opts...)
	if evalID == "" {
		err = errors.New("missing required eval_id parameter")
		return
	}
	if runID == "" {
		err = errors.New("missing required run_id parameter")
		return
	}
	if outputItemID == "" {
		err = errors.New("missing required output_item_id parameter")
		return
	}
	path := fmt.Sprintf("evals/%s/runs/%s/output_items/%s", evalID, runID, outputItemID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// Get a list of output items for an evaluation run.
func (r *EvalRunOutputItemService) List(ctx context.Context, evalID string, runID string, query EvalRunOutputItemListParams, opts ...option.RequestOption) (res *pagination.CursorPage[EvalRunOutputItemListResponse], err error) {
	var raw *http.Response
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	if evalID == "" {
		err = errors.New("missing required eval_id parameter")
		return
	}
	if runID == "" {
		err = errors.New("missing required run_id parameter")
		return
	}
	path := fmt.Sprintf("evals/%s/runs/%s/output_items", evalID, runID)
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

// Get a list of output items for an evaluation run.
func (r *EvalRunOutputItemService) ListAutoPaging(ctx context.Context, evalID string, runID string, query EvalRunOutputItemListParams, opts ...option.RequestOption) *pagination.CursorPageAutoPager[EvalRunOutputItemListResponse] {
	return pagination.NewCursorPageAutoPager(r.List(ctx, evalID, runID, query, opts...))
}

// A schema representing an evaluation run output item.
type EvalRunOutputItemGetResponse struct {
	// Unique identifier for the evaluation run output item.
	ID string `json:"id,required"`
	// Unix timestamp (in seconds) when the evaluation run was created.
	CreatedAt int64 `json:"created_at,required"`
	// Details of the input data source item.
	DatasourceItem map[string]interface{} `json:"datasource_item,required"`
	// The identifier for the data source item.
	DatasourceItemID int64 `json:"datasource_item_id,required"`
	// The identifier of the evaluation group.
	EvalID string `json:"eval_id,required"`
	// The type of the object. Always "eval.run.output_item".
	Object constant.EvalRunOutputItem `json:"object,required"`
	// A list of results from the evaluation run.
	Results []map[string]interface{} `json:"results,required"`
	// The identifier of the evaluation run associated with this output item.
	RunID string `json:"run_id,required"`
	// A sample containing the input and output of the evaluation run.
	Sample EvalRunOutputItemGetResponseSample `json:"sample,required"`
	// The status of the evaluation run.
	Status string `json:"status,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		ID               resp.Field
		CreatedAt        resp.Field
		DatasourceItem   resp.Field
		DatasourceItemID resp.Field
		EvalID           resp.Field
		Object           resp.Field
		Results          resp.Field
		RunID            resp.Field
		Sample           resp.Field
		Status           resp.Field
		ExtraFields      map[string]resp.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EvalRunOutputItemGetResponse) RawJSON() string { return r.JSON.raw }
func (r *EvalRunOutputItemGetResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A sample containing the input and output of the evaluation run.
type EvalRunOutputItemGetResponseSample struct {
	// An object representing an error response from the Eval API.
	Error EvalAPIError `json:"error,required"`
	// The reason why the sample generation was finished.
	FinishReason string `json:"finish_reason,required"`
	// An array of input messages.
	Input []EvalRunOutputItemGetResponseSampleInput `json:"input,required"`
	// The maximum number of tokens allowed for completion.
	MaxCompletionTokens int64 `json:"max_completion_tokens,required"`
	// The model used for generating the sample.
	Model string `json:"model,required"`
	// An array of output messages.
	Output []EvalRunOutputItemGetResponseSampleOutput `json:"output,required"`
	// The seed used for generating the sample.
	Seed int64 `json:"seed,required"`
	// The sampling temperature used.
	Temperature float64 `json:"temperature,required"`
	// The top_p value used for sampling.
	TopP float64 `json:"top_p,required"`
	// Token usage details for the sample.
	Usage EvalRunOutputItemGetResponseSampleUsage `json:"usage,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		Error               resp.Field
		FinishReason        resp.Field
		Input               resp.Field
		MaxCompletionTokens resp.Field
		Model               resp.Field
		Output              resp.Field
		Seed                resp.Field
		Temperature         resp.Field
		TopP                resp.Field
		Usage               resp.Field
		ExtraFields         map[string]resp.Field
		raw                 string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EvalRunOutputItemGetResponseSample) RawJSON() string { return r.JSON.raw }
func (r *EvalRunOutputItemGetResponseSample) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// An input message.
type EvalRunOutputItemGetResponseSampleInput struct {
	// The content of the message.
	Content string `json:"content,required"`
	// The role of the message sender (e.g., system, user, developer).
	Role string `json:"role,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		Content     resp.Field
		Role        resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EvalRunOutputItemGetResponseSampleInput) RawJSON() string { return r.JSON.raw }
func (r *EvalRunOutputItemGetResponseSampleInput) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type EvalRunOutputItemGetResponseSampleOutput struct {
	// The content of the message.
	Content string `json:"content"`
	// The role of the message (e.g. "system", "assistant", "user").
	Role string `json:"role"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		Content     resp.Field
		Role        resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EvalRunOutputItemGetResponseSampleOutput) RawJSON() string { return r.JSON.raw }
func (r *EvalRunOutputItemGetResponseSampleOutput) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Token usage details for the sample.
type EvalRunOutputItemGetResponseSampleUsage struct {
	// The number of tokens retrieved from cache.
	CachedTokens int64 `json:"cached_tokens,required"`
	// The number of completion tokens generated.
	CompletionTokens int64 `json:"completion_tokens,required"`
	// The number of prompt tokens used.
	PromptTokens int64 `json:"prompt_tokens,required"`
	// The total number of tokens used.
	TotalTokens int64 `json:"total_tokens,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		CachedTokens     resp.Field
		CompletionTokens resp.Field
		PromptTokens     resp.Field
		TotalTokens      resp.Field
		ExtraFields      map[string]resp.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EvalRunOutputItemGetResponseSampleUsage) RawJSON() string { return r.JSON.raw }
func (r *EvalRunOutputItemGetResponseSampleUsage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A schema representing an evaluation run output item.
type EvalRunOutputItemListResponse struct {
	// Unique identifier for the evaluation run output item.
	ID string `json:"id,required"`
	// Unix timestamp (in seconds) when the evaluation run was created.
	CreatedAt int64 `json:"created_at,required"`
	// Details of the input data source item.
	DatasourceItem map[string]interface{} `json:"datasource_item,required"`
	// The identifier for the data source item.
	DatasourceItemID int64 `json:"datasource_item_id,required"`
	// The identifier of the evaluation group.
	EvalID string `json:"eval_id,required"`
	// The type of the object. Always "eval.run.output_item".
	Object constant.EvalRunOutputItem `json:"object,required"`
	// A list of results from the evaluation run.
	Results []map[string]interface{} `json:"results,required"`
	// The identifier of the evaluation run associated with this output item.
	RunID string `json:"run_id,required"`
	// A sample containing the input and output of the evaluation run.
	Sample EvalRunOutputItemListResponseSample `json:"sample,required"`
	// The status of the evaluation run.
	Status string `json:"status,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		ID               resp.Field
		CreatedAt        resp.Field
		DatasourceItem   resp.Field
		DatasourceItemID resp.Field
		EvalID           resp.Field
		Object           resp.Field
		Results          resp.Field
		RunID            resp.Field
		Sample           resp.Field
		Status           resp.Field
		ExtraFields      map[string]resp.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EvalRunOutputItemListResponse) RawJSON() string { return r.JSON.raw }
func (r *EvalRunOutputItemListResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A sample containing the input and output of the evaluation run.
type EvalRunOutputItemListResponseSample struct {
	// An object representing an error response from the Eval API.
	Error EvalAPIError `json:"error,required"`
	// The reason why the sample generation was finished.
	FinishReason string `json:"finish_reason,required"`
	// An array of input messages.
	Input []EvalRunOutputItemListResponseSampleInput `json:"input,required"`
	// The maximum number of tokens allowed for completion.
	MaxCompletionTokens int64 `json:"max_completion_tokens,required"`
	// The model used for generating the sample.
	Model string `json:"model,required"`
	// An array of output messages.
	Output []EvalRunOutputItemListResponseSampleOutput `json:"output,required"`
	// The seed used for generating the sample.
	Seed int64 `json:"seed,required"`
	// The sampling temperature used.
	Temperature float64 `json:"temperature,required"`
	// The top_p value used for sampling.
	TopP float64 `json:"top_p,required"`
	// Token usage details for the sample.
	Usage EvalRunOutputItemListResponseSampleUsage `json:"usage,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		Error               resp.Field
		FinishReason        resp.Field
		Input               resp.Field
		MaxCompletionTokens resp.Field
		Model               resp.Field
		Output              resp.Field
		Seed                resp.Field
		Temperature         resp.Field
		TopP                resp.Field
		Usage               resp.Field
		ExtraFields         map[string]resp.Field
		raw                 string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EvalRunOutputItemListResponseSample) RawJSON() string { return r.JSON.raw }
func (r *EvalRunOutputItemListResponseSample) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// An input message.
type EvalRunOutputItemListResponseSampleInput struct {
	// The content of the message.
	Content string `json:"content,required"`
	// The role of the message sender (e.g., system, user, developer).
	Role string `json:"role,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		Content     resp.Field
		Role        resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EvalRunOutputItemListResponseSampleInput) RawJSON() string { return r.JSON.raw }
func (r *EvalRunOutputItemListResponseSampleInput) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type EvalRunOutputItemListResponseSampleOutput struct {
	// The content of the message.
	Content string `json:"content"`
	// The role of the message (e.g. "system", "assistant", "user").
	Role string `json:"role"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		Content     resp.Field
		Role        resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EvalRunOutputItemListResponseSampleOutput) RawJSON() string { return r.JSON.raw }
func (r *EvalRunOutputItemListResponseSampleOutput) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Token usage details for the sample.
type EvalRunOutputItemListResponseSampleUsage struct {
	// The number of tokens retrieved from cache.
	CachedTokens int64 `json:"cached_tokens,required"`
	// The number of completion tokens generated.
	CompletionTokens int64 `json:"completion_tokens,required"`
	// The number of prompt tokens used.
	PromptTokens int64 `json:"prompt_tokens,required"`
	// The total number of tokens used.
	TotalTokens int64 `json:"total_tokens,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		CachedTokens     resp.Field
		CompletionTokens resp.Field
		PromptTokens     resp.Field
		TotalTokens      resp.Field
		ExtraFields      map[string]resp.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EvalRunOutputItemListResponseSampleUsage) RawJSON() string { return r.JSON.raw }
func (r *EvalRunOutputItemListResponseSampleUsage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type EvalRunOutputItemListParams struct {
	// Identifier for the last output item from the previous pagination request.
	After param.Opt[string] `query:"after,omitzero" json:"-"`
	// Number of output items to retrieve.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Sort order for output items by timestamp. Use `asc` for ascending order or
	// `desc` for descending order. Defaults to `asc`.
	//
	// Any of "asc", "desc".
	Order EvalRunOutputItemListParamsOrder `query:"order,omitzero" json:"-"`
	// Filter output items by status. Use `failed` to filter by failed output items or
	// `pass` to filter by passed output items.
	//
	// Any of "fail", "pass".
	Status EvalRunOutputItemListParamsStatus `query:"status,omitzero" json:"-"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f EvalRunOutputItemListParams) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }

// URLQuery serializes [EvalRunOutputItemListParams]'s query parameters as
// `url.Values`.
func (r EvalRunOutputItemListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Sort order for output items by timestamp. Use `asc` for ascending order or
// `desc` for descending order. Defaults to `asc`.
type EvalRunOutputItemListParamsOrder string

const (
	EvalRunOutputItemListParamsOrderAsc  EvalRunOutputItemListParamsOrder = "asc"
	EvalRunOutputItemListParamsOrderDesc EvalRunOutputItemListParamsOrder = "desc"
)

// Filter output items by status. Use `failed` to filter by failed output items or
// `pass` to filter by passed output items.
type EvalRunOutputItemListParamsStatus string

const (
	EvalRunOutputItemListParamsStatusFail EvalRunOutputItemListParamsStatus = "fail"
	EvalRunOutputItemListParamsStatusPass EvalRunOutputItemListParamsStatus = "pass"
)

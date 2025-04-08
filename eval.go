// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"

	"github.com/openai/openai-go/internal/apijson"
	"github.com/openai/openai-go/internal/apiquery"
	"github.com/openai/openai-go/internal/requestconfig"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/pagination"
	"github.com/openai/openai-go/packages/param"
	"github.com/openai/openai-go/packages/resp"
	"github.com/openai/openai-go/shared"
	"github.com/openai/openai-go/shared/constant"
	"github.com/tidwall/gjson"
)

// EvalService contains methods and other services that help with interacting with
// the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewEvalService] method instead.
type EvalService struct {
	Options []option.RequestOption
	Runs    EvalRunService
}

// NewEvalService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewEvalService(opts ...option.RequestOption) (r EvalService) {
	r = EvalService{}
	r.Options = opts
	r.Runs = NewEvalRunService(opts...)
	return
}

// Create the structure of an evaluation that can be used to test a model's
// performance. An evaluation is a set of testing criteria and a datasource. After
// creating an evaluation, you can run it on different models and model parameters.
// We support several types of graders and datasources. For more information, see
// the [Evals guide](https://platform.openai.com/docs/guides/evals).
func (r *EvalService) New(ctx context.Context, body EvalNewParams, opts ...option.RequestOption) (res *EvalNewResponse, err error) {
	opts = append(r.Options[:], opts...)
	path := "evals"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Get an evaluation by ID.
func (r *EvalService) Get(ctx context.Context, evalID string, opts ...option.RequestOption) (res *EvalGetResponse, err error) {
	opts = append(r.Options[:], opts...)
	if evalID == "" {
		err = errors.New("missing required eval_id parameter")
		return
	}
	path := fmt.Sprintf("evals/%s", evalID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// Update certain properties of an evaluation.
func (r *EvalService) Update(ctx context.Context, evalID string, body EvalUpdateParams, opts ...option.RequestOption) (res *EvalUpdateResponse, err error) {
	opts = append(r.Options[:], opts...)
	if evalID == "" {
		err = errors.New("missing required eval_id parameter")
		return
	}
	path := fmt.Sprintf("evals/%s", evalID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// List evaluations for a project.
func (r *EvalService) List(ctx context.Context, query EvalListParams, opts ...option.RequestOption) (res *pagination.CursorPage[EvalListResponse], err error) {
	var raw *http.Response
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "evals"
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

// List evaluations for a project.
func (r *EvalService) ListAutoPaging(ctx context.Context, query EvalListParams, opts ...option.RequestOption) *pagination.CursorPageAutoPager[EvalListResponse] {
	return pagination.NewCursorPageAutoPager(r.List(ctx, query, opts...))
}

// Delete an evaluation.
func (r *EvalService) Delete(ctx context.Context, evalID string, opts ...option.RequestOption) (res *EvalDeleteResponse, err error) {
	opts = append(r.Options[:], opts...)
	if evalID == "" {
		err = errors.New("missing required eval_id parameter")
		return
	}
	path := fmt.Sprintf("evals/%s", evalID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return
}

// A CustomDataSourceConfig which specifies the schema of your `item` and
// optionally `sample` namespaces. The response schema defines the shape of the
// data that will be:
//
// - Used to define your testing criteria and
// - What data is required when creating a run
type EvalCustomDataSourceConfig struct {
	// The json schema for the run data source items. Learn how to build JSON schemas
	// [here](https://json-schema.org/).
	Schema map[string]interface{} `json:"schema,required"`
	// The type of data source. Always `custom`.
	Type constant.Custom `json:"type,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		Schema      resp.Field
		Type        resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EvalCustomDataSourceConfig) RawJSON() string { return r.JSON.raw }
func (r *EvalCustomDataSourceConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A LabelModelGrader object which uses a model to assign labels to each item in
// the evaluation.
type EvalLabelModelGrader struct {
	Input []EvalLabelModelGraderInputUnion `json:"input,required"`
	// The labels to assign to each item in the evaluation.
	Labels []string `json:"labels,required"`
	// The model to use for the evaluation. Must support structured outputs.
	Model string `json:"model,required"`
	// The name of the grader.
	Name string `json:"name,required"`
	// The labels that indicate a passing result. Must be a subset of labels.
	PassingLabels []string `json:"passing_labels,required"`
	// The object type, which is always `label_model`.
	Type constant.LabelModel `json:"type,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		Input         resp.Field
		Labels        resp.Field
		Model         resp.Field
		Name          resp.Field
		PassingLabels resp.Field
		Type          resp.Field
		ExtraFields   map[string]resp.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EvalLabelModelGrader) RawJSON() string { return r.JSON.raw }
func (r *EvalLabelModelGrader) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// EvalLabelModelGraderInputUnion contains all possible properties and values from
// [EvalLabelModelGraderInputInputMessage], [EvalLabelModelGraderInputAssistant].
//
// Use the [EvalLabelModelGraderInputUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type EvalLabelModelGraderInputUnion struct {
	// This field is a union of [EvalLabelModelGraderInputInputMessageContent],
	// [EvalLabelModelGraderInputAssistantContent]
	Content EvalLabelModelGraderInputUnionContent `json:"content"`
	// Any of nil, "assistant".
	Role string `json:"role"`
	Type string `json:"type"`
	JSON struct {
		Content resp.Field
		Role    resp.Field
		Type    resp.Field
		raw     string
	} `json:"-"`
}

func (u EvalLabelModelGraderInputUnion) AsInputMessage() (v EvalLabelModelGraderInputInputMessage) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u EvalLabelModelGraderInputUnion) AsAssistant() (v EvalLabelModelGraderInputAssistant) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u EvalLabelModelGraderInputUnion) RawJSON() string { return u.JSON.raw }

func (r *EvalLabelModelGraderInputUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// EvalLabelModelGraderInputUnionContent is an implicit subunion of
// [EvalLabelModelGraderInputUnion]. EvalLabelModelGraderInputUnionContent provides
// convenient access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [EvalLabelModelGraderInputUnion].
type EvalLabelModelGraderInputUnionContent struct {
	Text string `json:"text"`
	Type string `json:"type"`
	JSON struct {
		Text resp.Field
		Type resp.Field
		raw  string
	} `json:"-"`
}

func (r *EvalLabelModelGraderInputUnionContent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type EvalLabelModelGraderInputInputMessage struct {
	Content EvalLabelModelGraderInputInputMessageContent `json:"content,required"`
	// The role of the message. One of `user`, `system`, or `developer`.
	//
	// Any of "user", "system", "developer".
	Role string `json:"role,required"`
	// The type of item, which is always `message`.
	//
	// Any of "message".
	Type string `json:"type,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		Content     resp.Field
		Role        resp.Field
		Type        resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EvalLabelModelGraderInputInputMessage) RawJSON() string { return r.JSON.raw }
func (r *EvalLabelModelGraderInputInputMessage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type EvalLabelModelGraderInputInputMessageContent struct {
	// The text content.
	Text string `json:"text,required"`
	// The type of content, which is always `input_text`.
	//
	// Any of "input_text".
	Type string `json:"type,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		Text        resp.Field
		Type        resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EvalLabelModelGraderInputInputMessageContent) RawJSON() string { return r.JSON.raw }
func (r *EvalLabelModelGraderInputInputMessageContent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type EvalLabelModelGraderInputAssistant struct {
	Content EvalLabelModelGraderInputAssistantContent `json:"content,required"`
	// The role of the message. Must be `assistant` for output.
	Role constant.Assistant `json:"role,required"`
	// The type of item, which is always `message`.
	//
	// Any of "message".
	Type string `json:"type,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		Content     resp.Field
		Role        resp.Field
		Type        resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EvalLabelModelGraderInputAssistant) RawJSON() string { return r.JSON.raw }
func (r *EvalLabelModelGraderInputAssistant) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type EvalLabelModelGraderInputAssistantContent struct {
	// The text content.
	Text string `json:"text,required"`
	// The type of content, which is always `output_text`.
	//
	// Any of "output_text".
	Type string `json:"type,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		Text        resp.Field
		Type        resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EvalLabelModelGraderInputAssistantContent) RawJSON() string { return r.JSON.raw }
func (r *EvalLabelModelGraderInputAssistantContent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A StoredCompletionsDataSourceConfig which specifies the metadata property of
// your stored completions query. This is usually metadata like `usecase=chatbot`
// or `prompt-version=v2`, etc. The schema returned by this data source config is
// used to defined what variables are available in your evals. `item` and `sample`
// are both defined when using this data source config.
type EvalStoredCompletionsDataSourceConfig struct {
	// The json schema for the run data source items. Learn how to build JSON schemas
	// [here](https://json-schema.org/).
	Schema map[string]interface{} `json:"schema,required"`
	// The type of data source. Always `stored_completions`.
	Type constant.StoredCompletions `json:"type,required"`
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard.
	//
	// Keys are strings with a maximum length of 64 characters. Values are strings with
	// a maximum length of 512 characters.
	Metadata shared.Metadata `json:"metadata,nullable"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		Schema      resp.Field
		Type        resp.Field
		Metadata    resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EvalStoredCompletionsDataSourceConfig) RawJSON() string { return r.JSON.raw }
func (r *EvalStoredCompletionsDataSourceConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A StringCheckGrader object that performs a string comparison between input and
// reference using a specified operation.
type EvalStringCheckGrader struct {
	// The input text. This may include template strings.
	Input string `json:"input,required"`
	// The name of the grader.
	Name string `json:"name,required"`
	// The string check operation to perform. One of `eq`, `ne`, `like`, or `ilike`.
	//
	// Any of "eq", "ne", "like", "ilike".
	Operation EvalStringCheckGraderOperation `json:"operation,required"`
	// The reference text. This may include template strings.
	Reference string `json:"reference,required"`
	// The object type, which is always `string_check`.
	Type constant.StringCheck `json:"type,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		Input       resp.Field
		Name        resp.Field
		Operation   resp.Field
		Reference   resp.Field
		Type        resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EvalStringCheckGrader) RawJSON() string { return r.JSON.raw }
func (r *EvalStringCheckGrader) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this EvalStringCheckGrader to a EvalStringCheckGraderParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// EvalStringCheckGraderParam.IsOverridden()
func (r EvalStringCheckGrader) ToParam() EvalStringCheckGraderParam {
	return param.OverrideObj[EvalStringCheckGraderParam](r.RawJSON())
}

// The string check operation to perform. One of `eq`, `ne`, `like`, or `ilike`.
type EvalStringCheckGraderOperation string

const (
	EvalStringCheckGraderOperationEq    EvalStringCheckGraderOperation = "eq"
	EvalStringCheckGraderOperationNe    EvalStringCheckGraderOperation = "ne"
	EvalStringCheckGraderOperationLike  EvalStringCheckGraderOperation = "like"
	EvalStringCheckGraderOperationIlike EvalStringCheckGraderOperation = "ilike"
)

// A StringCheckGrader object that performs a string comparison between input and
// reference using a specified operation.
//
// The properties Input, Name, Operation, Reference, Type are required.
type EvalStringCheckGraderParam struct {
	// The input text. This may include template strings.
	Input string `json:"input,required"`
	// The name of the grader.
	Name string `json:"name,required"`
	// The string check operation to perform. One of `eq`, `ne`, `like`, or `ilike`.
	//
	// Any of "eq", "ne", "like", "ilike".
	Operation EvalStringCheckGraderOperation `json:"operation,omitzero,required"`
	// The reference text. This may include template strings.
	Reference string `json:"reference,required"`
	// The object type, which is always `string_check`.
	//
	// This field can be elided, and will marshal its zero value as "string_check".
	Type constant.StringCheck `json:"type,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f EvalStringCheckGraderParam) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r EvalStringCheckGraderParam) MarshalJSON() (data []byte, err error) {
	type shadow EvalStringCheckGraderParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// A TextSimilarityGrader object which grades text based on similarity metrics.
type EvalTextSimilarityGrader struct {
	// The evaluation metric to use. One of `cosine`, `fuzzy_match`, `bleu`, `gleu`,
	// `meteor`, `rouge_1`, `rouge_2`, `rouge_3`, `rouge_4`, `rouge_5`, or `rouge_l`.
	//
	// Any of "fuzzy_match", "bleu", "gleu", "meteor", "rouge_1", "rouge_2", "rouge_3",
	// "rouge_4", "rouge_5", "rouge_l", "cosine".
	EvaluationMetric EvalTextSimilarityGraderEvaluationMetric `json:"evaluation_metric,required"`
	// The text being graded.
	Input string `json:"input,required"`
	// A float score where a value greater than or equal indicates a passing grade.
	PassThreshold float64 `json:"pass_threshold,required"`
	// The text being graded against.
	Reference string `json:"reference,required"`
	// The type of grader.
	Type constant.TextSimilarity `json:"type,required"`
	// The name of the grader.
	Name string `json:"name"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		EvaluationMetric resp.Field
		Input            resp.Field
		PassThreshold    resp.Field
		Reference        resp.Field
		Type             resp.Field
		Name             resp.Field
		ExtraFields      map[string]resp.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EvalTextSimilarityGrader) RawJSON() string { return r.JSON.raw }
func (r *EvalTextSimilarityGrader) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this EvalTextSimilarityGrader to a
// EvalTextSimilarityGraderParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// EvalTextSimilarityGraderParam.IsOverridden()
func (r EvalTextSimilarityGrader) ToParam() EvalTextSimilarityGraderParam {
	return param.OverrideObj[EvalTextSimilarityGraderParam](r.RawJSON())
}

// The evaluation metric to use. One of `cosine`, `fuzzy_match`, `bleu`, `gleu`,
// `meteor`, `rouge_1`, `rouge_2`, `rouge_3`, `rouge_4`, `rouge_5`, or `rouge_l`.
type EvalTextSimilarityGraderEvaluationMetric string

const (
	EvalTextSimilarityGraderEvaluationMetricFuzzyMatch EvalTextSimilarityGraderEvaluationMetric = "fuzzy_match"
	EvalTextSimilarityGraderEvaluationMetricBleu       EvalTextSimilarityGraderEvaluationMetric = "bleu"
	EvalTextSimilarityGraderEvaluationMetricGleu       EvalTextSimilarityGraderEvaluationMetric = "gleu"
	EvalTextSimilarityGraderEvaluationMetricMeteor     EvalTextSimilarityGraderEvaluationMetric = "meteor"
	EvalTextSimilarityGraderEvaluationMetricRouge1     EvalTextSimilarityGraderEvaluationMetric = "rouge_1"
	EvalTextSimilarityGraderEvaluationMetricRouge2     EvalTextSimilarityGraderEvaluationMetric = "rouge_2"
	EvalTextSimilarityGraderEvaluationMetricRouge3     EvalTextSimilarityGraderEvaluationMetric = "rouge_3"
	EvalTextSimilarityGraderEvaluationMetricRouge4     EvalTextSimilarityGraderEvaluationMetric = "rouge_4"
	EvalTextSimilarityGraderEvaluationMetricRouge5     EvalTextSimilarityGraderEvaluationMetric = "rouge_5"
	EvalTextSimilarityGraderEvaluationMetricRougeL     EvalTextSimilarityGraderEvaluationMetric = "rouge_l"
	EvalTextSimilarityGraderEvaluationMetricCosine     EvalTextSimilarityGraderEvaluationMetric = "cosine"
)

// A TextSimilarityGrader object which grades text based on similarity metrics.
//
// The properties EvaluationMetric, Input, PassThreshold, Reference, Type are
// required.
type EvalTextSimilarityGraderParam struct {
	// The evaluation metric to use. One of `cosine`, `fuzzy_match`, `bleu`, `gleu`,
	// `meteor`, `rouge_1`, `rouge_2`, `rouge_3`, `rouge_4`, `rouge_5`, or `rouge_l`.
	//
	// Any of "fuzzy_match", "bleu", "gleu", "meteor", "rouge_1", "rouge_2", "rouge_3",
	// "rouge_4", "rouge_5", "rouge_l", "cosine".
	EvaluationMetric EvalTextSimilarityGraderEvaluationMetric `json:"evaluation_metric,omitzero,required"`
	// The text being graded.
	Input string `json:"input,required"`
	// A float score where a value greater than or equal indicates a passing grade.
	PassThreshold float64 `json:"pass_threshold,required"`
	// The text being graded against.
	Reference string `json:"reference,required"`
	// The name of the grader.
	Name param.Opt[string] `json:"name,omitzero"`
	// The type of grader.
	//
	// This field can be elided, and will marshal its zero value as "text_similarity".
	Type constant.TextSimilarity `json:"type,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f EvalTextSimilarityGraderParam) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r EvalTextSimilarityGraderParam) MarshalJSON() (data []byte, err error) {
	type shadow EvalTextSimilarityGraderParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// An Eval object with a data source config and testing criteria. An Eval
// represents a task to be done for your LLM integration. Like:
//
// - Improve the quality of my chatbot
// - See how well my chatbot handles customer support
// - Check if o3-mini is better at my usecase than gpt-4o
type EvalNewResponse struct {
	// Unique identifier for the evaluation.
	ID string `json:"id,required"`
	// The Unix timestamp (in seconds) for when the eval was created.
	CreatedAt int64 `json:"created_at,required"`
	// Configuration of data sources used in runs of the evaluation.
	DataSourceConfig EvalNewResponseDataSourceConfigUnion `json:"data_source_config,required"`
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard.
	//
	// Keys are strings with a maximum length of 64 characters. Values are strings with
	// a maximum length of 512 characters.
	Metadata shared.Metadata `json:"metadata,required"`
	// The name of the evaluation.
	Name string `json:"name,required"`
	// The object type.
	Object constant.Eval `json:"object,required"`
	// Indicates whether the evaluation is shared with OpenAI.
	ShareWithOpenAI bool `json:"share_with_openai,required"`
	// A list of testing criteria.
	TestingCriteria []EvalNewResponseTestingCriterionUnion `json:"testing_criteria,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		ID               resp.Field
		CreatedAt        resp.Field
		DataSourceConfig resp.Field
		Metadata         resp.Field
		Name             resp.Field
		Object           resp.Field
		ShareWithOpenAI  resp.Field
		TestingCriteria  resp.Field
		ExtraFields      map[string]resp.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EvalNewResponse) RawJSON() string { return r.JSON.raw }
func (r *EvalNewResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// EvalNewResponseDataSourceConfigUnion contains all possible properties and values
// from [EvalCustomDataSourceConfig], [EvalStoredCompletionsDataSourceConfig].
//
// Use the [EvalNewResponseDataSourceConfigUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type EvalNewResponseDataSourceConfigUnion struct {
	// This field is a union of [map[string]interface{}], [map[string]interface{}]
	Schema EvalNewResponseDataSourceConfigUnionSchema `json:"schema"`
	// Any of "custom", "stored_completions".
	Type string `json:"type"`
	// This field is from variant [EvalStoredCompletionsDataSourceConfig].
	Metadata shared.Metadata `json:"metadata"`
	JSON     struct {
		Schema   resp.Field
		Type     resp.Field
		Metadata resp.Field
		raw      string
	} `json:"-"`
}

// anyEvalNewResponseDataSourceConfig is implemented by each variant of
// [EvalNewResponseDataSourceConfigUnion] to add type safety for the return type of
// [EvalNewResponseDataSourceConfigUnion.AsAny]
type anyEvalNewResponseDataSourceConfig interface {
	implEvalNewResponseDataSourceConfigUnion()
}

func (EvalCustomDataSourceConfig) implEvalNewResponseDataSourceConfigUnion()            {}
func (EvalStoredCompletionsDataSourceConfig) implEvalNewResponseDataSourceConfigUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := EvalNewResponseDataSourceConfigUnion.AsAny().(type) {
//	case EvalCustomDataSourceConfig:
//	case EvalStoredCompletionsDataSourceConfig:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u EvalNewResponseDataSourceConfigUnion) AsAny() anyEvalNewResponseDataSourceConfig {
	switch u.Type {
	case "custom":
		return u.AsCustom()
	case "stored_completions":
		return u.AsStoredCompletions()
	}
	return nil
}

func (u EvalNewResponseDataSourceConfigUnion) AsCustom() (v EvalCustomDataSourceConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u EvalNewResponseDataSourceConfigUnion) AsStoredCompletions() (v EvalStoredCompletionsDataSourceConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u EvalNewResponseDataSourceConfigUnion) RawJSON() string { return u.JSON.raw }

func (r *EvalNewResponseDataSourceConfigUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// EvalNewResponseDataSourceConfigUnionSchema is an implicit subunion of
// [EvalNewResponseDataSourceConfigUnion].
// EvalNewResponseDataSourceConfigUnionSchema provides convenient access to the
// sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [EvalNewResponseDataSourceConfigUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfEvalCustomDataSourceConfigSchema
// OfEvalStoredCompletionsDataSourceConfigSchema]
type EvalNewResponseDataSourceConfigUnionSchema struct {
	// This field will be present if the value is a [interface{}] instead of an object.
	OfEvalCustomDataSourceConfigSchema interface{} `json:",inline"`
	// This field will be present if the value is a [interface{}] instead of an object.
	OfEvalStoredCompletionsDataSourceConfigSchema interface{} `json:",inline"`
	JSON                                          struct {
		OfEvalCustomDataSourceConfigSchema            resp.Field
		OfEvalStoredCompletionsDataSourceConfigSchema resp.Field
		raw                                           string
	} `json:"-"`
}

func (r *EvalNewResponseDataSourceConfigUnionSchema) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// EvalNewResponseTestingCriterionUnion contains all possible properties and values
// from [EvalLabelModelGrader], [EvalStringCheckGrader],
// [EvalTextSimilarityGrader].
//
// Use the [EvalNewResponseTestingCriterionUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type EvalNewResponseTestingCriterionUnion struct {
	// This field is a union of [[]EvalLabelModelGraderInputUnion], [string], [string]
	Input EvalNewResponseTestingCriterionUnionInput `json:"input"`
	// This field is from variant [EvalLabelModelGrader].
	Labels []string `json:"labels"`
	// This field is from variant [EvalLabelModelGrader].
	Model string `json:"model"`
	Name  string `json:"name"`
	// This field is from variant [EvalLabelModelGrader].
	PassingLabels []string `json:"passing_labels"`
	// Any of "label_model", "string_check", "text_similarity".
	Type string `json:"type"`
	// This field is from variant [EvalStringCheckGrader].
	Operation EvalStringCheckGraderOperation `json:"operation"`
	Reference string                         `json:"reference"`
	// This field is from variant [EvalTextSimilarityGrader].
	EvaluationMetric EvalTextSimilarityGraderEvaluationMetric `json:"evaluation_metric"`
	// This field is from variant [EvalTextSimilarityGrader].
	PassThreshold float64 `json:"pass_threshold"`
	JSON          struct {
		Input            resp.Field
		Labels           resp.Field
		Model            resp.Field
		Name             resp.Field
		PassingLabels    resp.Field
		Type             resp.Field
		Operation        resp.Field
		Reference        resp.Field
		EvaluationMetric resp.Field
		PassThreshold    resp.Field
		raw              string
	} `json:"-"`
}

// anyEvalNewResponseTestingCriterion is implemented by each variant of
// [EvalNewResponseTestingCriterionUnion] to add type safety for the return type of
// [EvalNewResponseTestingCriterionUnion.AsAny]
type anyEvalNewResponseTestingCriterion interface {
	implEvalNewResponseTestingCriterionUnion()
}

func (EvalLabelModelGrader) implEvalNewResponseTestingCriterionUnion()     {}
func (EvalStringCheckGrader) implEvalNewResponseTestingCriterionUnion()    {}
func (EvalTextSimilarityGrader) implEvalNewResponseTestingCriterionUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := EvalNewResponseTestingCriterionUnion.AsAny().(type) {
//	case EvalLabelModelGrader:
//	case EvalStringCheckGrader:
//	case EvalTextSimilarityGrader:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u EvalNewResponseTestingCriterionUnion) AsAny() anyEvalNewResponseTestingCriterion {
	switch u.Type {
	case "label_model":
		return u.AsLabelModel()
	case "string_check":
		return u.AsStringCheck()
	case "text_similarity":
		return u.AsTextSimilarity()
	}
	return nil
}

func (u EvalNewResponseTestingCriterionUnion) AsLabelModel() (v EvalLabelModelGrader) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u EvalNewResponseTestingCriterionUnion) AsStringCheck() (v EvalStringCheckGrader) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u EvalNewResponseTestingCriterionUnion) AsTextSimilarity() (v EvalTextSimilarityGrader) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u EvalNewResponseTestingCriterionUnion) RawJSON() string { return u.JSON.raw }

func (r *EvalNewResponseTestingCriterionUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// EvalNewResponseTestingCriterionUnionInput is an implicit subunion of
// [EvalNewResponseTestingCriterionUnion].
// EvalNewResponseTestingCriterionUnionInput provides convenient access to the
// sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [EvalNewResponseTestingCriterionUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfEvalLabelModelGraderInput OfString]
type EvalNewResponseTestingCriterionUnionInput struct {
	// This field will be present if the value is a [[]EvalLabelModelGraderInputUnion]
	// instead of an object.
	OfEvalLabelModelGraderInput []EvalLabelModelGraderInputUnion `json:",inline"`
	// This field will be present if the value is a [string] instead of an object.
	OfString string `json:",inline"`
	JSON     struct {
		OfEvalLabelModelGraderInput resp.Field
		OfString                    resp.Field
		raw                         string
	} `json:"-"`
}

func (r *EvalNewResponseTestingCriterionUnionInput) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// An Eval object with a data source config and testing criteria. An Eval
// represents a task to be done for your LLM integration. Like:
//
// - Improve the quality of my chatbot
// - See how well my chatbot handles customer support
// - Check if o3-mini is better at my usecase than gpt-4o
type EvalGetResponse struct {
	// Unique identifier for the evaluation.
	ID string `json:"id,required"`
	// The Unix timestamp (in seconds) for when the eval was created.
	CreatedAt int64 `json:"created_at,required"`
	// Configuration of data sources used in runs of the evaluation.
	DataSourceConfig EvalGetResponseDataSourceConfigUnion `json:"data_source_config,required"`
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard.
	//
	// Keys are strings with a maximum length of 64 characters. Values are strings with
	// a maximum length of 512 characters.
	Metadata shared.Metadata `json:"metadata,required"`
	// The name of the evaluation.
	Name string `json:"name,required"`
	// The object type.
	Object constant.Eval `json:"object,required"`
	// Indicates whether the evaluation is shared with OpenAI.
	ShareWithOpenAI bool `json:"share_with_openai,required"`
	// A list of testing criteria.
	TestingCriteria []EvalGetResponseTestingCriterionUnion `json:"testing_criteria,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		ID               resp.Field
		CreatedAt        resp.Field
		DataSourceConfig resp.Field
		Metadata         resp.Field
		Name             resp.Field
		Object           resp.Field
		ShareWithOpenAI  resp.Field
		TestingCriteria  resp.Field
		ExtraFields      map[string]resp.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EvalGetResponse) RawJSON() string { return r.JSON.raw }
func (r *EvalGetResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// EvalGetResponseDataSourceConfigUnion contains all possible properties and values
// from [EvalCustomDataSourceConfig], [EvalStoredCompletionsDataSourceConfig].
//
// Use the [EvalGetResponseDataSourceConfigUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type EvalGetResponseDataSourceConfigUnion struct {
	// This field is a union of [map[string]interface{}], [map[string]interface{}]
	Schema EvalGetResponseDataSourceConfigUnionSchema `json:"schema"`
	// Any of "custom", "stored_completions".
	Type string `json:"type"`
	// This field is from variant [EvalStoredCompletionsDataSourceConfig].
	Metadata shared.Metadata `json:"metadata"`
	JSON     struct {
		Schema   resp.Field
		Type     resp.Field
		Metadata resp.Field
		raw      string
	} `json:"-"`
}

// anyEvalGetResponseDataSourceConfig is implemented by each variant of
// [EvalGetResponseDataSourceConfigUnion] to add type safety for the return type of
// [EvalGetResponseDataSourceConfigUnion.AsAny]
type anyEvalGetResponseDataSourceConfig interface {
	implEvalGetResponseDataSourceConfigUnion()
}

func (EvalCustomDataSourceConfig) implEvalGetResponseDataSourceConfigUnion()            {}
func (EvalStoredCompletionsDataSourceConfig) implEvalGetResponseDataSourceConfigUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := EvalGetResponseDataSourceConfigUnion.AsAny().(type) {
//	case EvalCustomDataSourceConfig:
//	case EvalStoredCompletionsDataSourceConfig:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u EvalGetResponseDataSourceConfigUnion) AsAny() anyEvalGetResponseDataSourceConfig {
	switch u.Type {
	case "custom":
		return u.AsCustom()
	case "stored_completions":
		return u.AsStoredCompletions()
	}
	return nil
}

func (u EvalGetResponseDataSourceConfigUnion) AsCustom() (v EvalCustomDataSourceConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u EvalGetResponseDataSourceConfigUnion) AsStoredCompletions() (v EvalStoredCompletionsDataSourceConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u EvalGetResponseDataSourceConfigUnion) RawJSON() string { return u.JSON.raw }

func (r *EvalGetResponseDataSourceConfigUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// EvalGetResponseDataSourceConfigUnionSchema is an implicit subunion of
// [EvalGetResponseDataSourceConfigUnion].
// EvalGetResponseDataSourceConfigUnionSchema provides convenient access to the
// sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [EvalGetResponseDataSourceConfigUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfEvalCustomDataSourceConfigSchema
// OfEvalStoredCompletionsDataSourceConfigSchema]
type EvalGetResponseDataSourceConfigUnionSchema struct {
	// This field will be present if the value is a [interface{}] instead of an object.
	OfEvalCustomDataSourceConfigSchema interface{} `json:",inline"`
	// This field will be present if the value is a [interface{}] instead of an object.
	OfEvalStoredCompletionsDataSourceConfigSchema interface{} `json:",inline"`
	JSON                                          struct {
		OfEvalCustomDataSourceConfigSchema            resp.Field
		OfEvalStoredCompletionsDataSourceConfigSchema resp.Field
		raw                                           string
	} `json:"-"`
}

func (r *EvalGetResponseDataSourceConfigUnionSchema) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// EvalGetResponseTestingCriterionUnion contains all possible properties and values
// from [EvalLabelModelGrader], [EvalStringCheckGrader],
// [EvalTextSimilarityGrader].
//
// Use the [EvalGetResponseTestingCriterionUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type EvalGetResponseTestingCriterionUnion struct {
	// This field is a union of [[]EvalLabelModelGraderInputUnion], [string], [string]
	Input EvalGetResponseTestingCriterionUnionInput `json:"input"`
	// This field is from variant [EvalLabelModelGrader].
	Labels []string `json:"labels"`
	// This field is from variant [EvalLabelModelGrader].
	Model string `json:"model"`
	Name  string `json:"name"`
	// This field is from variant [EvalLabelModelGrader].
	PassingLabels []string `json:"passing_labels"`
	// Any of "label_model", "string_check", "text_similarity".
	Type string `json:"type"`
	// This field is from variant [EvalStringCheckGrader].
	Operation EvalStringCheckGraderOperation `json:"operation"`
	Reference string                         `json:"reference"`
	// This field is from variant [EvalTextSimilarityGrader].
	EvaluationMetric EvalTextSimilarityGraderEvaluationMetric `json:"evaluation_metric"`
	// This field is from variant [EvalTextSimilarityGrader].
	PassThreshold float64 `json:"pass_threshold"`
	JSON          struct {
		Input            resp.Field
		Labels           resp.Field
		Model            resp.Field
		Name             resp.Field
		PassingLabels    resp.Field
		Type             resp.Field
		Operation        resp.Field
		Reference        resp.Field
		EvaluationMetric resp.Field
		PassThreshold    resp.Field
		raw              string
	} `json:"-"`
}

// anyEvalGetResponseTestingCriterion is implemented by each variant of
// [EvalGetResponseTestingCriterionUnion] to add type safety for the return type of
// [EvalGetResponseTestingCriterionUnion.AsAny]
type anyEvalGetResponseTestingCriterion interface {
	implEvalGetResponseTestingCriterionUnion()
}

func (EvalLabelModelGrader) implEvalGetResponseTestingCriterionUnion()     {}
func (EvalStringCheckGrader) implEvalGetResponseTestingCriterionUnion()    {}
func (EvalTextSimilarityGrader) implEvalGetResponseTestingCriterionUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := EvalGetResponseTestingCriterionUnion.AsAny().(type) {
//	case EvalLabelModelGrader:
//	case EvalStringCheckGrader:
//	case EvalTextSimilarityGrader:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u EvalGetResponseTestingCriterionUnion) AsAny() anyEvalGetResponseTestingCriterion {
	switch u.Type {
	case "label_model":
		return u.AsLabelModel()
	case "string_check":
		return u.AsStringCheck()
	case "text_similarity":
		return u.AsTextSimilarity()
	}
	return nil
}

func (u EvalGetResponseTestingCriterionUnion) AsLabelModel() (v EvalLabelModelGrader) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u EvalGetResponseTestingCriterionUnion) AsStringCheck() (v EvalStringCheckGrader) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u EvalGetResponseTestingCriterionUnion) AsTextSimilarity() (v EvalTextSimilarityGrader) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u EvalGetResponseTestingCriterionUnion) RawJSON() string { return u.JSON.raw }

func (r *EvalGetResponseTestingCriterionUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// EvalGetResponseTestingCriterionUnionInput is an implicit subunion of
// [EvalGetResponseTestingCriterionUnion].
// EvalGetResponseTestingCriterionUnionInput provides convenient access to the
// sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [EvalGetResponseTestingCriterionUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfEvalLabelModelGraderInput OfString]
type EvalGetResponseTestingCriterionUnionInput struct {
	// This field will be present if the value is a [[]EvalLabelModelGraderInputUnion]
	// instead of an object.
	OfEvalLabelModelGraderInput []EvalLabelModelGraderInputUnion `json:",inline"`
	// This field will be present if the value is a [string] instead of an object.
	OfString string `json:",inline"`
	JSON     struct {
		OfEvalLabelModelGraderInput resp.Field
		OfString                    resp.Field
		raw                         string
	} `json:"-"`
}

func (r *EvalGetResponseTestingCriterionUnionInput) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// An Eval object with a data source config and testing criteria. An Eval
// represents a task to be done for your LLM integration. Like:
//
// - Improve the quality of my chatbot
// - See how well my chatbot handles customer support
// - Check if o3-mini is better at my usecase than gpt-4o
type EvalUpdateResponse struct {
	// Unique identifier for the evaluation.
	ID string `json:"id,required"`
	// The Unix timestamp (in seconds) for when the eval was created.
	CreatedAt int64 `json:"created_at,required"`
	// Configuration of data sources used in runs of the evaluation.
	DataSourceConfig EvalUpdateResponseDataSourceConfigUnion `json:"data_source_config,required"`
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard.
	//
	// Keys are strings with a maximum length of 64 characters. Values are strings with
	// a maximum length of 512 characters.
	Metadata shared.Metadata `json:"metadata,required"`
	// The name of the evaluation.
	Name string `json:"name,required"`
	// The object type.
	Object constant.Eval `json:"object,required"`
	// Indicates whether the evaluation is shared with OpenAI.
	ShareWithOpenAI bool `json:"share_with_openai,required"`
	// A list of testing criteria.
	TestingCriteria []EvalUpdateResponseTestingCriterionUnion `json:"testing_criteria,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		ID               resp.Field
		CreatedAt        resp.Field
		DataSourceConfig resp.Field
		Metadata         resp.Field
		Name             resp.Field
		Object           resp.Field
		ShareWithOpenAI  resp.Field
		TestingCriteria  resp.Field
		ExtraFields      map[string]resp.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EvalUpdateResponse) RawJSON() string { return r.JSON.raw }
func (r *EvalUpdateResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// EvalUpdateResponseDataSourceConfigUnion contains all possible properties and
// values from [EvalCustomDataSourceConfig],
// [EvalStoredCompletionsDataSourceConfig].
//
// Use the [EvalUpdateResponseDataSourceConfigUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type EvalUpdateResponseDataSourceConfigUnion struct {
	// This field is a union of [map[string]interface{}], [map[string]interface{}]
	Schema EvalUpdateResponseDataSourceConfigUnionSchema `json:"schema"`
	// Any of "custom", "stored_completions".
	Type string `json:"type"`
	// This field is from variant [EvalStoredCompletionsDataSourceConfig].
	Metadata shared.Metadata `json:"metadata"`
	JSON     struct {
		Schema   resp.Field
		Type     resp.Field
		Metadata resp.Field
		raw      string
	} `json:"-"`
}

// anyEvalUpdateResponseDataSourceConfig is implemented by each variant of
// [EvalUpdateResponseDataSourceConfigUnion] to add type safety for the return type
// of [EvalUpdateResponseDataSourceConfigUnion.AsAny]
type anyEvalUpdateResponseDataSourceConfig interface {
	implEvalUpdateResponseDataSourceConfigUnion()
}

func (EvalCustomDataSourceConfig) implEvalUpdateResponseDataSourceConfigUnion()            {}
func (EvalStoredCompletionsDataSourceConfig) implEvalUpdateResponseDataSourceConfigUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := EvalUpdateResponseDataSourceConfigUnion.AsAny().(type) {
//	case EvalCustomDataSourceConfig:
//	case EvalStoredCompletionsDataSourceConfig:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u EvalUpdateResponseDataSourceConfigUnion) AsAny() anyEvalUpdateResponseDataSourceConfig {
	switch u.Type {
	case "custom":
		return u.AsCustom()
	case "stored_completions":
		return u.AsStoredCompletions()
	}
	return nil
}

func (u EvalUpdateResponseDataSourceConfigUnion) AsCustom() (v EvalCustomDataSourceConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u EvalUpdateResponseDataSourceConfigUnion) AsStoredCompletions() (v EvalStoredCompletionsDataSourceConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u EvalUpdateResponseDataSourceConfigUnion) RawJSON() string { return u.JSON.raw }

func (r *EvalUpdateResponseDataSourceConfigUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// EvalUpdateResponseDataSourceConfigUnionSchema is an implicit subunion of
// [EvalUpdateResponseDataSourceConfigUnion].
// EvalUpdateResponseDataSourceConfigUnionSchema provides convenient access to the
// sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [EvalUpdateResponseDataSourceConfigUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfEvalCustomDataSourceConfigSchema
// OfEvalStoredCompletionsDataSourceConfigSchema]
type EvalUpdateResponseDataSourceConfigUnionSchema struct {
	// This field will be present if the value is a [interface{}] instead of an object.
	OfEvalCustomDataSourceConfigSchema interface{} `json:",inline"`
	// This field will be present if the value is a [interface{}] instead of an object.
	OfEvalStoredCompletionsDataSourceConfigSchema interface{} `json:",inline"`
	JSON                                          struct {
		OfEvalCustomDataSourceConfigSchema            resp.Field
		OfEvalStoredCompletionsDataSourceConfigSchema resp.Field
		raw                                           string
	} `json:"-"`
}

func (r *EvalUpdateResponseDataSourceConfigUnionSchema) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// EvalUpdateResponseTestingCriterionUnion contains all possible properties and
// values from [EvalLabelModelGrader], [EvalStringCheckGrader],
// [EvalTextSimilarityGrader].
//
// Use the [EvalUpdateResponseTestingCriterionUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type EvalUpdateResponseTestingCriterionUnion struct {
	// This field is a union of [[]EvalLabelModelGraderInputUnion], [string], [string]
	Input EvalUpdateResponseTestingCriterionUnionInput `json:"input"`
	// This field is from variant [EvalLabelModelGrader].
	Labels []string `json:"labels"`
	// This field is from variant [EvalLabelModelGrader].
	Model string `json:"model"`
	Name  string `json:"name"`
	// This field is from variant [EvalLabelModelGrader].
	PassingLabels []string `json:"passing_labels"`
	// Any of "label_model", "string_check", "text_similarity".
	Type string `json:"type"`
	// This field is from variant [EvalStringCheckGrader].
	Operation EvalStringCheckGraderOperation `json:"operation"`
	Reference string                         `json:"reference"`
	// This field is from variant [EvalTextSimilarityGrader].
	EvaluationMetric EvalTextSimilarityGraderEvaluationMetric `json:"evaluation_metric"`
	// This field is from variant [EvalTextSimilarityGrader].
	PassThreshold float64 `json:"pass_threshold"`
	JSON          struct {
		Input            resp.Field
		Labels           resp.Field
		Model            resp.Field
		Name             resp.Field
		PassingLabels    resp.Field
		Type             resp.Field
		Operation        resp.Field
		Reference        resp.Field
		EvaluationMetric resp.Field
		PassThreshold    resp.Field
		raw              string
	} `json:"-"`
}

// anyEvalUpdateResponseTestingCriterion is implemented by each variant of
// [EvalUpdateResponseTestingCriterionUnion] to add type safety for the return type
// of [EvalUpdateResponseTestingCriterionUnion.AsAny]
type anyEvalUpdateResponseTestingCriterion interface {
	implEvalUpdateResponseTestingCriterionUnion()
}

func (EvalLabelModelGrader) implEvalUpdateResponseTestingCriterionUnion()     {}
func (EvalStringCheckGrader) implEvalUpdateResponseTestingCriterionUnion()    {}
func (EvalTextSimilarityGrader) implEvalUpdateResponseTestingCriterionUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := EvalUpdateResponseTestingCriterionUnion.AsAny().(type) {
//	case EvalLabelModelGrader:
//	case EvalStringCheckGrader:
//	case EvalTextSimilarityGrader:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u EvalUpdateResponseTestingCriterionUnion) AsAny() anyEvalUpdateResponseTestingCriterion {
	switch u.Type {
	case "label_model":
		return u.AsLabelModel()
	case "string_check":
		return u.AsStringCheck()
	case "text_similarity":
		return u.AsTextSimilarity()
	}
	return nil
}

func (u EvalUpdateResponseTestingCriterionUnion) AsLabelModel() (v EvalLabelModelGrader) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u EvalUpdateResponseTestingCriterionUnion) AsStringCheck() (v EvalStringCheckGrader) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u EvalUpdateResponseTestingCriterionUnion) AsTextSimilarity() (v EvalTextSimilarityGrader) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u EvalUpdateResponseTestingCriterionUnion) RawJSON() string { return u.JSON.raw }

func (r *EvalUpdateResponseTestingCriterionUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// EvalUpdateResponseTestingCriterionUnionInput is an implicit subunion of
// [EvalUpdateResponseTestingCriterionUnion].
// EvalUpdateResponseTestingCriterionUnionInput provides convenient access to the
// sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [EvalUpdateResponseTestingCriterionUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfEvalLabelModelGraderInput OfString]
type EvalUpdateResponseTestingCriterionUnionInput struct {
	// This field will be present if the value is a [[]EvalLabelModelGraderInputUnion]
	// instead of an object.
	OfEvalLabelModelGraderInput []EvalLabelModelGraderInputUnion `json:",inline"`
	// This field will be present if the value is a [string] instead of an object.
	OfString string `json:",inline"`
	JSON     struct {
		OfEvalLabelModelGraderInput resp.Field
		OfString                    resp.Field
		raw                         string
	} `json:"-"`
}

func (r *EvalUpdateResponseTestingCriterionUnionInput) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// An Eval object with a data source config and testing criteria. An Eval
// represents a task to be done for your LLM integration. Like:
//
// - Improve the quality of my chatbot
// - See how well my chatbot handles customer support
// - Check if o3-mini is better at my usecase than gpt-4o
type EvalListResponse struct {
	// Unique identifier for the evaluation.
	ID string `json:"id,required"`
	// The Unix timestamp (in seconds) for when the eval was created.
	CreatedAt int64 `json:"created_at,required"`
	// Configuration of data sources used in runs of the evaluation.
	DataSourceConfig EvalListResponseDataSourceConfigUnion `json:"data_source_config,required"`
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard.
	//
	// Keys are strings with a maximum length of 64 characters. Values are strings with
	// a maximum length of 512 characters.
	Metadata shared.Metadata `json:"metadata,required"`
	// The name of the evaluation.
	Name string `json:"name,required"`
	// The object type.
	Object constant.Eval `json:"object,required"`
	// Indicates whether the evaluation is shared with OpenAI.
	ShareWithOpenAI bool `json:"share_with_openai,required"`
	// A list of testing criteria.
	TestingCriteria []EvalListResponseTestingCriterionUnion `json:"testing_criteria,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		ID               resp.Field
		CreatedAt        resp.Field
		DataSourceConfig resp.Field
		Metadata         resp.Field
		Name             resp.Field
		Object           resp.Field
		ShareWithOpenAI  resp.Field
		TestingCriteria  resp.Field
		ExtraFields      map[string]resp.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EvalListResponse) RawJSON() string { return r.JSON.raw }
func (r *EvalListResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// EvalListResponseDataSourceConfigUnion contains all possible properties and
// values from [EvalCustomDataSourceConfig],
// [EvalStoredCompletionsDataSourceConfig].
//
// Use the [EvalListResponseDataSourceConfigUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type EvalListResponseDataSourceConfigUnion struct {
	// This field is a union of [map[string]interface{}], [map[string]interface{}]
	Schema EvalListResponseDataSourceConfigUnionSchema `json:"schema"`
	// Any of "custom", "stored_completions".
	Type string `json:"type"`
	// This field is from variant [EvalStoredCompletionsDataSourceConfig].
	Metadata shared.Metadata `json:"metadata"`
	JSON     struct {
		Schema   resp.Field
		Type     resp.Field
		Metadata resp.Field
		raw      string
	} `json:"-"`
}

// anyEvalListResponseDataSourceConfig is implemented by each variant of
// [EvalListResponseDataSourceConfigUnion] to add type safety for the return type
// of [EvalListResponseDataSourceConfigUnion.AsAny]
type anyEvalListResponseDataSourceConfig interface {
	implEvalListResponseDataSourceConfigUnion()
}

func (EvalCustomDataSourceConfig) implEvalListResponseDataSourceConfigUnion()            {}
func (EvalStoredCompletionsDataSourceConfig) implEvalListResponseDataSourceConfigUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := EvalListResponseDataSourceConfigUnion.AsAny().(type) {
//	case EvalCustomDataSourceConfig:
//	case EvalStoredCompletionsDataSourceConfig:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u EvalListResponseDataSourceConfigUnion) AsAny() anyEvalListResponseDataSourceConfig {
	switch u.Type {
	case "custom":
		return u.AsCustom()
	case "stored_completions":
		return u.AsStoredCompletions()
	}
	return nil
}

func (u EvalListResponseDataSourceConfigUnion) AsCustom() (v EvalCustomDataSourceConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u EvalListResponseDataSourceConfigUnion) AsStoredCompletions() (v EvalStoredCompletionsDataSourceConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u EvalListResponseDataSourceConfigUnion) RawJSON() string { return u.JSON.raw }

func (r *EvalListResponseDataSourceConfigUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// EvalListResponseDataSourceConfigUnionSchema is an implicit subunion of
// [EvalListResponseDataSourceConfigUnion].
// EvalListResponseDataSourceConfigUnionSchema provides convenient access to the
// sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [EvalListResponseDataSourceConfigUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfEvalCustomDataSourceConfigSchema
// OfEvalStoredCompletionsDataSourceConfigSchema]
type EvalListResponseDataSourceConfigUnionSchema struct {
	// This field will be present if the value is a [interface{}] instead of an object.
	OfEvalCustomDataSourceConfigSchema interface{} `json:",inline"`
	// This field will be present if the value is a [interface{}] instead of an object.
	OfEvalStoredCompletionsDataSourceConfigSchema interface{} `json:",inline"`
	JSON                                          struct {
		OfEvalCustomDataSourceConfigSchema            resp.Field
		OfEvalStoredCompletionsDataSourceConfigSchema resp.Field
		raw                                           string
	} `json:"-"`
}

func (r *EvalListResponseDataSourceConfigUnionSchema) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// EvalListResponseTestingCriterionUnion contains all possible properties and
// values from [EvalLabelModelGrader], [EvalStringCheckGrader],
// [EvalTextSimilarityGrader].
//
// Use the [EvalListResponseTestingCriterionUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type EvalListResponseTestingCriterionUnion struct {
	// This field is a union of [[]EvalLabelModelGraderInputUnion], [string], [string]
	Input EvalListResponseTestingCriterionUnionInput `json:"input"`
	// This field is from variant [EvalLabelModelGrader].
	Labels []string `json:"labels"`
	// This field is from variant [EvalLabelModelGrader].
	Model string `json:"model"`
	Name  string `json:"name"`
	// This field is from variant [EvalLabelModelGrader].
	PassingLabels []string `json:"passing_labels"`
	// Any of "label_model", "string_check", "text_similarity".
	Type string `json:"type"`
	// This field is from variant [EvalStringCheckGrader].
	Operation EvalStringCheckGraderOperation `json:"operation"`
	Reference string                         `json:"reference"`
	// This field is from variant [EvalTextSimilarityGrader].
	EvaluationMetric EvalTextSimilarityGraderEvaluationMetric `json:"evaluation_metric"`
	// This field is from variant [EvalTextSimilarityGrader].
	PassThreshold float64 `json:"pass_threshold"`
	JSON          struct {
		Input            resp.Field
		Labels           resp.Field
		Model            resp.Field
		Name             resp.Field
		PassingLabels    resp.Field
		Type             resp.Field
		Operation        resp.Field
		Reference        resp.Field
		EvaluationMetric resp.Field
		PassThreshold    resp.Field
		raw              string
	} `json:"-"`
}

// anyEvalListResponseTestingCriterion is implemented by each variant of
// [EvalListResponseTestingCriterionUnion] to add type safety for the return type
// of [EvalListResponseTestingCriterionUnion.AsAny]
type anyEvalListResponseTestingCriterion interface {
	implEvalListResponseTestingCriterionUnion()
}

func (EvalLabelModelGrader) implEvalListResponseTestingCriterionUnion()     {}
func (EvalStringCheckGrader) implEvalListResponseTestingCriterionUnion()    {}
func (EvalTextSimilarityGrader) implEvalListResponseTestingCriterionUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := EvalListResponseTestingCriterionUnion.AsAny().(type) {
//	case EvalLabelModelGrader:
//	case EvalStringCheckGrader:
//	case EvalTextSimilarityGrader:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u EvalListResponseTestingCriterionUnion) AsAny() anyEvalListResponseTestingCriterion {
	switch u.Type {
	case "label_model":
		return u.AsLabelModel()
	case "string_check":
		return u.AsStringCheck()
	case "text_similarity":
		return u.AsTextSimilarity()
	}
	return nil
}

func (u EvalListResponseTestingCriterionUnion) AsLabelModel() (v EvalLabelModelGrader) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u EvalListResponseTestingCriterionUnion) AsStringCheck() (v EvalStringCheckGrader) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u EvalListResponseTestingCriterionUnion) AsTextSimilarity() (v EvalTextSimilarityGrader) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u EvalListResponseTestingCriterionUnion) RawJSON() string { return u.JSON.raw }

func (r *EvalListResponseTestingCriterionUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// EvalListResponseTestingCriterionUnionInput is an implicit subunion of
// [EvalListResponseTestingCriterionUnion].
// EvalListResponseTestingCriterionUnionInput provides convenient access to the
// sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [EvalListResponseTestingCriterionUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfEvalLabelModelGraderInput OfString]
type EvalListResponseTestingCriterionUnionInput struct {
	// This field will be present if the value is a [[]EvalLabelModelGraderInputUnion]
	// instead of an object.
	OfEvalLabelModelGraderInput []EvalLabelModelGraderInputUnion `json:",inline"`
	// This field will be present if the value is a [string] instead of an object.
	OfString string `json:",inline"`
	JSON     struct {
		OfEvalLabelModelGraderInput resp.Field
		OfString                    resp.Field
		raw                         string
	} `json:"-"`
}

func (r *EvalListResponseTestingCriterionUnionInput) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type EvalDeleteResponse struct {
	Deleted bool   `json:"deleted,required"`
	EvalID  string `json:"eval_id,required"`
	Object  string `json:"object,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		Deleted     resp.Field
		EvalID      resp.Field
		Object      resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EvalDeleteResponse) RawJSON() string { return r.JSON.raw }
func (r *EvalDeleteResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type EvalNewParams struct {
	// The configuration for the data source used for the evaluation runs.
	DataSourceConfig EvalNewParamsDataSourceConfigUnion `json:"data_source_config,omitzero,required"`
	// A list of graders for all eval runs in this group.
	TestingCriteria []EvalNewParamsTestingCriterionUnion `json:"testing_criteria,omitzero,required"`
	// The name of the evaluation.
	Name param.Opt[string] `json:"name,omitzero"`
	// Indicates whether the evaluation is shared with OpenAI.
	ShareWithOpenAI param.Opt[bool] `json:"share_with_openai,omitzero"`
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard.
	//
	// Keys are strings with a maximum length of 64 characters. Values are strings with
	// a maximum length of 512 characters.
	Metadata shared.MetadataParam `json:"metadata,omitzero"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f EvalNewParams) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }

func (r EvalNewParams) MarshalJSON() (data []byte, err error) {
	type shadow EvalNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type EvalNewParamsDataSourceConfigUnion struct {
	OfCustom            *EvalNewParamsDataSourceConfigCustom            `json:",omitzero,inline"`
	OfStoredCompletions *EvalNewParamsDataSourceConfigStoredCompletions `json:",omitzero,inline"`
	paramUnion
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (u EvalNewParamsDataSourceConfigUnion) IsPresent() bool {
	return !param.IsOmitted(u) && !u.IsNull()
}
func (u EvalNewParamsDataSourceConfigUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion[EvalNewParamsDataSourceConfigUnion](u.OfCustom, u.OfStoredCompletions)
}

func (u *EvalNewParamsDataSourceConfigUnion) asAny() any {
	if !param.IsOmitted(u.OfCustom) {
		return u.OfCustom
	} else if !param.IsOmitted(u.OfStoredCompletions) {
		return u.OfStoredCompletions
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u EvalNewParamsDataSourceConfigUnion) GetItemSchema() map[string]interface{} {
	if vt := u.OfCustom; vt != nil {
		return vt.ItemSchema
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u EvalNewParamsDataSourceConfigUnion) GetIncludeSampleSchema() *bool {
	if vt := u.OfCustom; vt != nil && vt.IncludeSampleSchema.IsPresent() {
		return &vt.IncludeSampleSchema.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u EvalNewParamsDataSourceConfigUnion) GetMetadata() shared.MetadataParam {
	if vt := u.OfStoredCompletions; vt != nil {
		return vt.Metadata
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u EvalNewParamsDataSourceConfigUnion) GetType() *string {
	if vt := u.OfCustom; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfStoredCompletions; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[EvalNewParamsDataSourceConfigUnion](
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(EvalNewParamsDataSourceConfigCustom{}),
			DiscriminatorValue: "custom",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(EvalNewParamsDataSourceConfigStoredCompletions{}),
			DiscriminatorValue: "stored_completions",
		},
	)
}

// A CustomDataSourceConfig object that defines the schema for the data source used
// for the evaluation runs. This schema is used to define the shape of the data
// that will be:
//
// - Used to define your testing criteria and
// - What data is required when creating a run
//
// The properties ItemSchema, Type are required.
type EvalNewParamsDataSourceConfigCustom struct {
	// The json schema for the run data source items.
	ItemSchema map[string]interface{} `json:"item_schema,omitzero,required"`
	// Whether to include the sample schema in the data source.
	IncludeSampleSchema param.Opt[bool] `json:"include_sample_schema,omitzero"`
	// The type of data source. Always `custom`.
	//
	// This field can be elided, and will marshal its zero value as "custom".
	Type constant.Custom `json:"type,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f EvalNewParamsDataSourceConfigCustom) IsPresent() bool {
	return !param.IsOmitted(f) && !f.IsNull()
}
func (r EvalNewParamsDataSourceConfigCustom) MarshalJSON() (data []byte, err error) {
	type shadow EvalNewParamsDataSourceConfigCustom
	return param.MarshalObject(r, (*shadow)(&r))
}

// A data source config which specifies the metadata property of your stored
// completions query. This is usually metadata like `usecase=chatbot` or
// `prompt-version=v2`, etc.
//
// The property Type is required.
type EvalNewParamsDataSourceConfigStoredCompletions struct {
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard.
	//
	// Keys are strings with a maximum length of 64 characters. Values are strings with
	// a maximum length of 512 characters.
	Metadata shared.MetadataParam `json:"metadata,omitzero"`
	// The type of data source. Always `stored_completions`.
	//
	// This field can be elided, and will marshal its zero value as
	// "stored_completions".
	Type constant.StoredCompletions `json:"type,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f EvalNewParamsDataSourceConfigStoredCompletions) IsPresent() bool {
	return !param.IsOmitted(f) && !f.IsNull()
}
func (r EvalNewParamsDataSourceConfigStoredCompletions) MarshalJSON() (data []byte, err error) {
	type shadow EvalNewParamsDataSourceConfigStoredCompletions
	return param.MarshalObject(r, (*shadow)(&r))
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type EvalNewParamsTestingCriterionUnion struct {
	OfLabelModel     *EvalNewParamsTestingCriterionLabelModel `json:",omitzero,inline"`
	OfStringCheck    *EvalStringCheckGraderParam              `json:",omitzero,inline"`
	OfTextSimilarity *EvalTextSimilarityGraderParam           `json:",omitzero,inline"`
	paramUnion
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (u EvalNewParamsTestingCriterionUnion) IsPresent() bool {
	return !param.IsOmitted(u) && !u.IsNull()
}
func (u EvalNewParamsTestingCriterionUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion[EvalNewParamsTestingCriterionUnion](u.OfLabelModel, u.OfStringCheck, u.OfTextSimilarity)
}

func (u *EvalNewParamsTestingCriterionUnion) asAny() any {
	if !param.IsOmitted(u.OfLabelModel) {
		return u.OfLabelModel
	} else if !param.IsOmitted(u.OfStringCheck) {
		return u.OfStringCheck
	} else if !param.IsOmitted(u.OfTextSimilarity) {
		return u.OfTextSimilarity
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u EvalNewParamsTestingCriterionUnion) GetLabels() []string {
	if vt := u.OfLabelModel; vt != nil {
		return vt.Labels
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u EvalNewParamsTestingCriterionUnion) GetModel() *string {
	if vt := u.OfLabelModel; vt != nil {
		return &vt.Model
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u EvalNewParamsTestingCriterionUnion) GetPassingLabels() []string {
	if vt := u.OfLabelModel; vt != nil {
		return vt.PassingLabels
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u EvalNewParamsTestingCriterionUnion) GetOperation() *string {
	if vt := u.OfStringCheck; vt != nil {
		return (*string)(&vt.Operation)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u EvalNewParamsTestingCriterionUnion) GetEvaluationMetric() *string {
	if vt := u.OfTextSimilarity; vt != nil {
		return (*string)(&vt.EvaluationMetric)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u EvalNewParamsTestingCriterionUnion) GetPassThreshold() *float64 {
	if vt := u.OfTextSimilarity; vt != nil {
		return &vt.PassThreshold
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u EvalNewParamsTestingCriterionUnion) GetName() *string {
	if vt := u.OfLabelModel; vt != nil {
		return (*string)(&vt.Name)
	} else if vt := u.OfStringCheck; vt != nil {
		return (*string)(&vt.Name)
	} else if vt := u.OfTextSimilarity; vt != nil && vt.Name.IsPresent() {
		return &vt.Name.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u EvalNewParamsTestingCriterionUnion) GetType() *string {
	if vt := u.OfLabelModel; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfStringCheck; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfTextSimilarity; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u EvalNewParamsTestingCriterionUnion) GetReference() *string {
	if vt := u.OfStringCheck; vt != nil {
		return (*string)(&vt.Reference)
	} else if vt := u.OfTextSimilarity; vt != nil {
		return (*string)(&vt.Reference)
	}
	return nil
}

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u EvalNewParamsTestingCriterionUnion) GetInput() (res evalNewParamsTestingCriterionUnionInput) {
	if vt := u.OfLabelModel; vt != nil {
		res.ofEvalNewsTestingCriterionLabelModelInput = &vt.Input
	} else if vt := u.OfStringCheck; vt != nil {
		res.ofString = &vt.Input
	} else if vt := u.OfTextSimilarity; vt != nil {
		res.ofString = &vt.Input
	}
	return
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type evalNewParamsTestingCriterionUnionInput struct {
	ofEvalNewsTestingCriterionLabelModelInput *[]EvalNewParamsTestingCriterionLabelModelInputUnion
	ofString                                  *string
}

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *[]openai.EvalNewParamsTestingCriterionLabelModelInputUnion:
//	case *string:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u evalNewParamsTestingCriterionUnionInput) AsAny() any {
	if !param.IsOmitted(u.ofEvalNewsTestingCriterionLabelModelInput) {
		return u.ofEvalNewsTestingCriterionLabelModelInput
	} else if !param.IsOmitted(u.ofString) {
		return u.ofString
	} else if !param.IsOmitted(u.ofString) {
		return u.ofString
	}
	return nil
}

func init() {
	apijson.RegisterUnion[EvalNewParamsTestingCriterionUnion](
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(EvalNewParamsTestingCriterionLabelModel{}),
			DiscriminatorValue: "label_model",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(EvalStringCheckGraderParam{}),
			DiscriminatorValue: "string_check",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(EvalTextSimilarityGraderParam{}),
			DiscriminatorValue: "text_similarity",
		},
	)
}

// A LabelModelGrader object which uses a model to assign labels to each item in
// the evaluation.
//
// The properties Input, Labels, Model, Name, PassingLabels, Type are required.
type EvalNewParamsTestingCriterionLabelModel struct {
	Input []EvalNewParamsTestingCriterionLabelModelInputUnion `json:"input,omitzero,required"`
	// The labels to classify to each item in the evaluation.
	Labels []string `json:"labels,omitzero,required"`
	// The model to use for the evaluation. Must support structured outputs.
	Model string `json:"model,required"`
	// The name of the grader.
	Name string `json:"name,required"`
	// The labels that indicate a passing result. Must be a subset of labels.
	PassingLabels []string `json:"passing_labels,omitzero,required"`
	// The object type, which is always `label_model`.
	//
	// This field can be elided, and will marshal its zero value as "label_model".
	Type constant.LabelModel `json:"type,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f EvalNewParamsTestingCriterionLabelModel) IsPresent() bool {
	return !param.IsOmitted(f) && !f.IsNull()
}
func (r EvalNewParamsTestingCriterionLabelModel) MarshalJSON() (data []byte, err error) {
	type shadow EvalNewParamsTestingCriterionLabelModel
	return param.MarshalObject(r, (*shadow)(&r))
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type EvalNewParamsTestingCriterionLabelModelInputUnion struct {
	OfSimpleInputMessage *EvalNewParamsTestingCriterionLabelModelInputSimpleInputMessage `json:",omitzero,inline"`
	OfInputMessage       *EvalNewParamsTestingCriterionLabelModelInputInputMessage       `json:",omitzero,inline"`
	OfOutputMessage      *EvalNewParamsTestingCriterionLabelModelInputOutputMessage      `json:",omitzero,inline"`
	paramUnion
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (u EvalNewParamsTestingCriterionLabelModelInputUnion) IsPresent() bool {
	return !param.IsOmitted(u) && !u.IsNull()
}
func (u EvalNewParamsTestingCriterionLabelModelInputUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion[EvalNewParamsTestingCriterionLabelModelInputUnion](u.OfSimpleInputMessage, u.OfInputMessage, u.OfOutputMessage)
}

func (u *EvalNewParamsTestingCriterionLabelModelInputUnion) asAny() any {
	if !param.IsOmitted(u.OfSimpleInputMessage) {
		return u.OfSimpleInputMessage
	} else if !param.IsOmitted(u.OfInputMessage) {
		return u.OfInputMessage
	} else if !param.IsOmitted(u.OfOutputMessage) {
		return u.OfOutputMessage
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u EvalNewParamsTestingCriterionLabelModelInputUnion) GetRole() *string {
	if vt := u.OfSimpleInputMessage; vt != nil {
		return (*string)(&vt.Role)
	} else if vt := u.OfInputMessage; vt != nil {
		return (*string)(&vt.Role)
	} else if vt := u.OfOutputMessage; vt != nil {
		return (*string)(&vt.Role)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u EvalNewParamsTestingCriterionLabelModelInputUnion) GetType() *string {
	if vt := u.OfInputMessage; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfOutputMessage; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u EvalNewParamsTestingCriterionLabelModelInputUnion) GetContent() (res evalNewParamsTestingCriterionLabelModelInputUnionContent) {
	if vt := u.OfSimpleInputMessage; vt != nil {
		res.ofString = &vt.Content
	} else if vt := u.OfInputMessage; vt != nil {
		res.ofEvalNewsTestingCriterionLabelModelInputInputMessageContent = &vt.Content
	} else if vt := u.OfOutputMessage; vt != nil {
		res.ofEvalNewsTestingCriterionLabelModelInputOutputMessageContent = &vt.Content
	}
	return
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type evalNewParamsTestingCriterionLabelModelInputUnionContent struct {
	ofString                                                      *string
	ofEvalNewsTestingCriterionLabelModelInputInputMessageContent  *EvalNewParamsTestingCriterionLabelModelInputInputMessageContent
	ofEvalNewsTestingCriterionLabelModelInputOutputMessageContent *EvalNewParamsTestingCriterionLabelModelInputOutputMessageContent
}

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *string:
//	case *openai.EvalNewParamsTestingCriterionLabelModelInputInputMessageContent:
//	case *openai.EvalNewParamsTestingCriterionLabelModelInputOutputMessageContent:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u evalNewParamsTestingCriterionLabelModelInputUnionContent) AsAny() any {
	if !param.IsOmitted(u.ofString) {
		return u.ofString
	} else if !param.IsOmitted(u.ofEvalNewsTestingCriterionLabelModelInputInputMessageContent) {
		return u.ofEvalNewsTestingCriterionLabelModelInputInputMessageContent
	} else if !param.IsOmitted(u.ofEvalNewsTestingCriterionLabelModelInputOutputMessageContent) {
		return u.ofEvalNewsTestingCriterionLabelModelInputOutputMessageContent
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u evalNewParamsTestingCriterionLabelModelInputUnionContent) GetText() *string {
	if vt := u.ofEvalNewsTestingCriterionLabelModelInputInputMessageContent; vt != nil {
		return (*string)(&vt.Text)
	} else if vt := u.ofEvalNewsTestingCriterionLabelModelInputOutputMessageContent; vt != nil {
		return (*string)(&vt.Text)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u evalNewParamsTestingCriterionLabelModelInputUnionContent) GetType() *string {
	if vt := u.ofEvalNewsTestingCriterionLabelModelInputInputMessageContent; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.ofEvalNewsTestingCriterionLabelModelInputOutputMessageContent; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// The properties Content, Role are required.
type EvalNewParamsTestingCriterionLabelModelInputSimpleInputMessage struct {
	// The content of the message.
	Content string `json:"content,required"`
	// The role of the message (e.g. "system", "assistant", "user").
	Role string `json:"role,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f EvalNewParamsTestingCriterionLabelModelInputSimpleInputMessage) IsPresent() bool {
	return !param.IsOmitted(f) && !f.IsNull()
}
func (r EvalNewParamsTestingCriterionLabelModelInputSimpleInputMessage) MarshalJSON() (data []byte, err error) {
	type shadow EvalNewParamsTestingCriterionLabelModelInputSimpleInputMessage
	return param.MarshalObject(r, (*shadow)(&r))
}

// The properties Content, Role, Type are required.
type EvalNewParamsTestingCriterionLabelModelInputInputMessage struct {
	Content EvalNewParamsTestingCriterionLabelModelInputInputMessageContent `json:"content,omitzero,required"`
	// The role of the message. One of `user`, `system`, or `developer`.
	//
	// Any of "user", "system", "developer".
	Role string `json:"role,omitzero,required"`
	// The type of item, which is always `message`.
	//
	// Any of "message".
	Type string `json:"type,omitzero,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f EvalNewParamsTestingCriterionLabelModelInputInputMessage) IsPresent() bool {
	return !param.IsOmitted(f) && !f.IsNull()
}
func (r EvalNewParamsTestingCriterionLabelModelInputInputMessage) MarshalJSON() (data []byte, err error) {
	type shadow EvalNewParamsTestingCriterionLabelModelInputInputMessage
	return param.MarshalObject(r, (*shadow)(&r))
}

func init() {
	apijson.RegisterFieldValidator[EvalNewParamsTestingCriterionLabelModelInputInputMessage](
		"Role", false, "user", "system", "developer",
	)
	apijson.RegisterFieldValidator[EvalNewParamsTestingCriterionLabelModelInputInputMessage](
		"Type", false, "message",
	)
}

// The properties Text, Type are required.
type EvalNewParamsTestingCriterionLabelModelInputInputMessageContent struct {
	// The text content.
	Text string `json:"text,required"`
	// The type of content, which is always `input_text`.
	//
	// Any of "input_text".
	Type string `json:"type,omitzero,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f EvalNewParamsTestingCriterionLabelModelInputInputMessageContent) IsPresent() bool {
	return !param.IsOmitted(f) && !f.IsNull()
}
func (r EvalNewParamsTestingCriterionLabelModelInputInputMessageContent) MarshalJSON() (data []byte, err error) {
	type shadow EvalNewParamsTestingCriterionLabelModelInputInputMessageContent
	return param.MarshalObject(r, (*shadow)(&r))
}

func init() {
	apijson.RegisterFieldValidator[EvalNewParamsTestingCriterionLabelModelInputInputMessageContent](
		"Type", false, "input_text",
	)
}

// The properties Content, Role, Type are required.
type EvalNewParamsTestingCriterionLabelModelInputOutputMessage struct {
	Content EvalNewParamsTestingCriterionLabelModelInputOutputMessageContent `json:"content,omitzero,required"`
	// The role of the message. Must be `assistant` for output.
	//
	// Any of "assistant".
	Role string `json:"role,omitzero,required"`
	// The type of item, which is always `message`.
	//
	// Any of "message".
	Type string `json:"type,omitzero,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f EvalNewParamsTestingCriterionLabelModelInputOutputMessage) IsPresent() bool {
	return !param.IsOmitted(f) && !f.IsNull()
}
func (r EvalNewParamsTestingCriterionLabelModelInputOutputMessage) MarshalJSON() (data []byte, err error) {
	type shadow EvalNewParamsTestingCriterionLabelModelInputOutputMessage
	return param.MarshalObject(r, (*shadow)(&r))
}

func init() {
	apijson.RegisterFieldValidator[EvalNewParamsTestingCriterionLabelModelInputOutputMessage](
		"Role", false, "assistant",
	)
	apijson.RegisterFieldValidator[EvalNewParamsTestingCriterionLabelModelInputOutputMessage](
		"Type", false, "message",
	)
}

// The properties Text, Type are required.
type EvalNewParamsTestingCriterionLabelModelInputOutputMessageContent struct {
	// The text content.
	Text string `json:"text,required"`
	// The type of content, which is always `output_text`.
	//
	// Any of "output_text".
	Type string `json:"type,omitzero,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f EvalNewParamsTestingCriterionLabelModelInputOutputMessageContent) IsPresent() bool {
	return !param.IsOmitted(f) && !f.IsNull()
}
func (r EvalNewParamsTestingCriterionLabelModelInputOutputMessageContent) MarshalJSON() (data []byte, err error) {
	type shadow EvalNewParamsTestingCriterionLabelModelInputOutputMessageContent
	return param.MarshalObject(r, (*shadow)(&r))
}

func init() {
	apijson.RegisterFieldValidator[EvalNewParamsTestingCriterionLabelModelInputOutputMessageContent](
		"Type", false, "output_text",
	)
}

type EvalUpdateParams struct {
	// Rename the evaluation.
	Name param.Opt[string] `json:"name,omitzero"`
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard.
	//
	// Keys are strings with a maximum length of 64 characters. Values are strings with
	// a maximum length of 512 characters.
	Metadata shared.MetadataParam `json:"metadata,omitzero"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f EvalUpdateParams) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }

func (r EvalUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow EvalUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}

type EvalListParams struct {
	// Identifier for the last eval from the previous pagination request.
	After param.Opt[string] `query:"after,omitzero" json:"-"`
	// Number of evals to retrieve.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Sort order for evals by timestamp. Use `asc` for ascending order or `desc` for
	// descending order.
	//
	// Any of "asc", "desc".
	Order EvalListParamsOrder `query:"order,omitzero" json:"-"`
	// Evals can be ordered by creation time or last updated time. Use `created_at` for
	// creation time or `updated_at` for last updated time.
	//
	// Any of "created_at", "updated_at".
	OrderBy EvalListParamsOrderBy `query:"order_by,omitzero" json:"-"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f EvalListParams) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }

// URLQuery serializes [EvalListParams]'s query parameters as `url.Values`.
func (r EvalListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Sort order for evals by timestamp. Use `asc` for ascending order or `desc` for
// descending order.
type EvalListParamsOrder string

const (
	EvalListParamsOrderAsc  EvalListParamsOrder = "asc"
	EvalListParamsOrderDesc EvalListParamsOrder = "desc"
)

// Evals can be ordered by creation time or last updated time. Use `created_at` for
// creation time or `updated_at` for last updated time.
type EvalListParamsOrderBy string

const (
	EvalListParamsOrderByCreatedAt EvalListParamsOrderBy = "created_at"
	EvalListParamsOrderByUpdatedAt EvalListParamsOrderBy = "updated_at"
)

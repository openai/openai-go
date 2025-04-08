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

// EvalRunService contains methods and other services that help with interacting
// with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewEvalRunService] method instead.
type EvalRunService struct {
	Options     []option.RequestOption
	OutputItems EvalRunOutputItemService
}

// NewEvalRunService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewEvalRunService(opts ...option.RequestOption) (r EvalRunService) {
	r = EvalRunService{}
	r.Options = opts
	r.OutputItems = NewEvalRunOutputItemService(opts...)
	return
}

// Create a new evaluation run. This is the endpoint that will kick off grading.
func (r *EvalRunService) New(ctx context.Context, evalID string, body EvalRunNewParams, opts ...option.RequestOption) (res *EvalRunNewResponse, err error) {
	opts = append(r.Options[:], opts...)
	if evalID == "" {
		err = errors.New("missing required eval_id parameter")
		return
	}
	path := fmt.Sprintf("evals/%s/runs", evalID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Get an evaluation run by ID.
func (r *EvalRunService) Get(ctx context.Context, evalID string, runID string, opts ...option.RequestOption) (res *EvalRunGetResponse, err error) {
	opts = append(r.Options[:], opts...)
	if evalID == "" {
		err = errors.New("missing required eval_id parameter")
		return
	}
	if runID == "" {
		err = errors.New("missing required run_id parameter")
		return
	}
	path := fmt.Sprintf("evals/%s/runs/%s", evalID, runID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// Get a list of runs for an evaluation.
func (r *EvalRunService) List(ctx context.Context, evalID string, query EvalRunListParams, opts ...option.RequestOption) (res *pagination.CursorPage[EvalRunListResponse], err error) {
	var raw *http.Response
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	if evalID == "" {
		err = errors.New("missing required eval_id parameter")
		return
	}
	path := fmt.Sprintf("evals/%s/runs", evalID)
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

// Get a list of runs for an evaluation.
func (r *EvalRunService) ListAutoPaging(ctx context.Context, evalID string, query EvalRunListParams, opts ...option.RequestOption) *pagination.CursorPageAutoPager[EvalRunListResponse] {
	return pagination.NewCursorPageAutoPager(r.List(ctx, evalID, query, opts...))
}

// Delete an eval run.
func (r *EvalRunService) Delete(ctx context.Context, evalID string, runID string, opts ...option.RequestOption) (res *EvalRunDeleteResponse, err error) {
	opts = append(r.Options[:], opts...)
	if evalID == "" {
		err = errors.New("missing required eval_id parameter")
		return
	}
	if runID == "" {
		err = errors.New("missing required run_id parameter")
		return
	}
	path := fmt.Sprintf("evals/%s/runs/%s", evalID, runID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return
}

// Cancel an ongoing evaluation run.
func (r *EvalRunService) Cancel(ctx context.Context, evalID string, runID string, opts ...option.RequestOption) (res *EvalRunCancelResponse, err error) {
	opts = append(r.Options[:], opts...)
	if evalID == "" {
		err = errors.New("missing required eval_id parameter")
		return
	}
	if runID == "" {
		err = errors.New("missing required run_id parameter")
		return
	}
	path := fmt.Sprintf("evals/%s/runs/%s", evalID, runID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return
}

// A CompletionsRunDataSource object describing a model sampling configuration.
type CreateEvalCompletionsRunDataSource struct {
	InputMessages CreateEvalCompletionsRunDataSourceInputMessagesUnion `json:"input_messages,required"`
	// The name of the model to use for generating completions (e.g. "o3-mini").
	Model string `json:"model,required"`
	// A StoredCompletionsRunDataSource configuration describing a set of filters
	Source CreateEvalCompletionsRunDataSourceSourceUnion `json:"source,required"`
	// The type of run data source. Always `completions`.
	//
	// Any of "completions".
	Type           CreateEvalCompletionsRunDataSourceType           `json:"type,required"`
	SamplingParams CreateEvalCompletionsRunDataSourceSamplingParams `json:"sampling_params"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		InputMessages  resp.Field
		Model          resp.Field
		Source         resp.Field
		Type           resp.Field
		SamplingParams resp.Field
		ExtraFields    map[string]resp.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CreateEvalCompletionsRunDataSource) RawJSON() string { return r.JSON.raw }
func (r *CreateEvalCompletionsRunDataSource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this CreateEvalCompletionsRunDataSource to a
// CreateEvalCompletionsRunDataSourceParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// CreateEvalCompletionsRunDataSourceParam.IsOverridden()
func (r CreateEvalCompletionsRunDataSource) ToParam() CreateEvalCompletionsRunDataSourceParam {
	return param.OverrideObj[CreateEvalCompletionsRunDataSourceParam](r.RawJSON())
}

// CreateEvalCompletionsRunDataSourceInputMessagesUnion contains all possible
// properties and values from
// [CreateEvalCompletionsRunDataSourceInputMessagesTemplate],
// [CreateEvalCompletionsRunDataSourceInputMessagesItemReference].
//
// Use the [CreateEvalCompletionsRunDataSourceInputMessagesUnion.AsAny] method to
// switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type CreateEvalCompletionsRunDataSourceInputMessagesUnion struct {
	// This field is from variant
	// [CreateEvalCompletionsRunDataSourceInputMessagesTemplate].
	Template []CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateUnion `json:"template"`
	// Any of "template", "item_reference".
	Type string `json:"type"`
	// This field is from variant
	// [CreateEvalCompletionsRunDataSourceInputMessagesItemReference].
	ItemReference string `json:"item_reference"`
	JSON          struct {
		Template      resp.Field
		Type          resp.Field
		ItemReference resp.Field
		raw           string
	} `json:"-"`
}

// anyCreateEvalCompletionsRunDataSourceInputMessages is implemented by each
// variant of [CreateEvalCompletionsRunDataSourceInputMessagesUnion] to add type
// safety for the return type of
// [CreateEvalCompletionsRunDataSourceInputMessagesUnion.AsAny]
type anyCreateEvalCompletionsRunDataSourceInputMessages interface {
	implCreateEvalCompletionsRunDataSourceInputMessagesUnion()
}

func (CreateEvalCompletionsRunDataSourceInputMessagesTemplate) implCreateEvalCompletionsRunDataSourceInputMessagesUnion() {
}
func (CreateEvalCompletionsRunDataSourceInputMessagesItemReference) implCreateEvalCompletionsRunDataSourceInputMessagesUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := CreateEvalCompletionsRunDataSourceInputMessagesUnion.AsAny().(type) {
//	case CreateEvalCompletionsRunDataSourceInputMessagesTemplate:
//	case CreateEvalCompletionsRunDataSourceInputMessagesItemReference:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u CreateEvalCompletionsRunDataSourceInputMessagesUnion) AsAny() anyCreateEvalCompletionsRunDataSourceInputMessages {
	switch u.Type {
	case "template":
		return u.AsTemplate()
	case "item_reference":
		return u.AsItemReference()
	}
	return nil
}

func (u CreateEvalCompletionsRunDataSourceInputMessagesUnion) AsTemplate() (v CreateEvalCompletionsRunDataSourceInputMessagesTemplate) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u CreateEvalCompletionsRunDataSourceInputMessagesUnion) AsItemReference() (v CreateEvalCompletionsRunDataSourceInputMessagesItemReference) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u CreateEvalCompletionsRunDataSourceInputMessagesUnion) RawJSON() string { return u.JSON.raw }

func (r *CreateEvalCompletionsRunDataSourceInputMessagesUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type CreateEvalCompletionsRunDataSourceInputMessagesTemplate struct {
	// A list of chat messages forming the prompt or context. May include variable
	// references to the "item" namespace, ie {{item.name}}.
	Template []CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateUnion `json:"template,required"`
	// The type of input messages. Always `template`.
	Type constant.Template `json:"type,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		Template    resp.Field
		Type        resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CreateEvalCompletionsRunDataSourceInputMessagesTemplate) RawJSON() string { return r.JSON.raw }
func (r *CreateEvalCompletionsRunDataSourceInputMessagesTemplate) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateUnion contains
// all possible properties and values from
// [CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateChatMessage],
// [CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateInputMessage],
// [CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateOutputMessage].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateUnion struct {
	// This field is a union of [string],
	// [CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateInputMessageContent],
	// [CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateOutputMessageContent]
	Content CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateUnionContent `json:"content"`
	Role    string                                                                      `json:"role"`
	Type    string                                                                      `json:"type"`
	JSON    struct {
		Content resp.Field
		Role    resp.Field
		Type    resp.Field
		raw     string
	} `json:"-"`
}

func (u CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateUnion) AsChatMessage() (v CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateChatMessage) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateUnion) AsInputMessage() (v CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateInputMessage) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateUnion) AsOutputMessage() (v CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateOutputMessage) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateUnion) RawJSON() string {
	return u.JSON.raw
}

func (r *CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateUnionContent is
// an implicit subunion of
// [CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateUnion].
// CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateUnionContent
// provides convenient access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfString]
type CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateUnionContent struct {
	// This field will be present if the value is a [string] instead of an object.
	OfString string `json:",inline"`
	Text     string `json:"text"`
	Type     string `json:"type"`
	JSON     struct {
		OfString resp.Field
		Text     resp.Field
		Type     resp.Field
		raw      string
	} `json:"-"`
}

func (r *CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateUnionContent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateChatMessage struct {
	// The content of the message.
	Content string `json:"content,required"`
	// The role of the message (e.g. "system", "assistant", "user").
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
func (r CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateChatMessage) RawJSON() string {
	return r.JSON.raw
}
func (r *CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateChatMessage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateInputMessage struct {
	Content CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateInputMessageContent `json:"content,required"`
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
func (r CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateInputMessage) RawJSON() string {
	return r.JSON.raw
}
func (r *CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateInputMessage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateInputMessageContent struct {
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
func (r CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateInputMessageContent) RawJSON() string {
	return r.JSON.raw
}
func (r *CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateInputMessageContent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateOutputMessage struct {
	Content CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateOutputMessageContent `json:"content,required"`
	// The role of the message. Must be `assistant` for output.
	//
	// Any of "assistant".
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
func (r CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateOutputMessage) RawJSON() string {
	return r.JSON.raw
}
func (r *CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateOutputMessage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateOutputMessageContent struct {
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
func (r CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateOutputMessageContent) RawJSON() string {
	return r.JSON.raw
}
func (r *CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateOutputMessageContent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type CreateEvalCompletionsRunDataSourceInputMessagesItemReference struct {
	// A reference to a variable in the "item" namespace. Ie, "item.name"
	ItemReference string `json:"item_reference,required"`
	// The type of input messages. Always `item_reference`.
	Type constant.ItemReference `json:"type,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		ItemReference resp.Field
		Type          resp.Field
		ExtraFields   map[string]resp.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CreateEvalCompletionsRunDataSourceInputMessagesItemReference) RawJSON() string {
	return r.JSON.raw
}
func (r *CreateEvalCompletionsRunDataSourceInputMessagesItemReference) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// CreateEvalCompletionsRunDataSourceSourceUnion contains all possible properties
// and values from [CreateEvalCompletionsRunDataSourceSourceFileContent],
// [CreateEvalCompletionsRunDataSourceSourceFileID],
// [CreateEvalCompletionsRunDataSourceSourceStoredCompletions].
//
// Use the [CreateEvalCompletionsRunDataSourceSourceUnion.AsAny] method to switch
// on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type CreateEvalCompletionsRunDataSourceSourceUnion struct {
	// This field is from variant
	// [CreateEvalCompletionsRunDataSourceSourceFileContent].
	Content []CreateEvalCompletionsRunDataSourceSourceFileContentContent `json:"content"`
	// Any of "file_content", "file_id", "stored_completions".
	Type string `json:"type"`
	// This field is from variant [CreateEvalCompletionsRunDataSourceSourceFileID].
	ID string `json:"id"`
	// This field is from variant
	// [CreateEvalCompletionsRunDataSourceSourceStoredCompletions].
	CreatedAfter int64 `json:"created_after"`
	// This field is from variant
	// [CreateEvalCompletionsRunDataSourceSourceStoredCompletions].
	CreatedBefore int64 `json:"created_before"`
	// This field is from variant
	// [CreateEvalCompletionsRunDataSourceSourceStoredCompletions].
	Limit int64 `json:"limit"`
	// This field is from variant
	// [CreateEvalCompletionsRunDataSourceSourceStoredCompletions].
	Metadata shared.Metadata `json:"metadata"`
	// This field is from variant
	// [CreateEvalCompletionsRunDataSourceSourceStoredCompletions].
	Model string `json:"model"`
	JSON  struct {
		Content       resp.Field
		Type          resp.Field
		ID            resp.Field
		CreatedAfter  resp.Field
		CreatedBefore resp.Field
		Limit         resp.Field
		Metadata      resp.Field
		Model         resp.Field
		raw           string
	} `json:"-"`
}

// anyCreateEvalCompletionsRunDataSourceSource is implemented by each variant of
// [CreateEvalCompletionsRunDataSourceSourceUnion] to add type safety for the
// return type of [CreateEvalCompletionsRunDataSourceSourceUnion.AsAny]
type anyCreateEvalCompletionsRunDataSourceSource interface {
	implCreateEvalCompletionsRunDataSourceSourceUnion()
}

func (CreateEvalCompletionsRunDataSourceSourceFileContent) implCreateEvalCompletionsRunDataSourceSourceUnion() {
}
func (CreateEvalCompletionsRunDataSourceSourceFileID) implCreateEvalCompletionsRunDataSourceSourceUnion() {
}
func (CreateEvalCompletionsRunDataSourceSourceStoredCompletions) implCreateEvalCompletionsRunDataSourceSourceUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := CreateEvalCompletionsRunDataSourceSourceUnion.AsAny().(type) {
//	case CreateEvalCompletionsRunDataSourceSourceFileContent:
//	case CreateEvalCompletionsRunDataSourceSourceFileID:
//	case CreateEvalCompletionsRunDataSourceSourceStoredCompletions:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u CreateEvalCompletionsRunDataSourceSourceUnion) AsAny() anyCreateEvalCompletionsRunDataSourceSource {
	switch u.Type {
	case "file_content":
		return u.AsFileContent()
	case "file_id":
		return u.AsFileID()
	case "stored_completions":
		return u.AsStoredCompletions()
	}
	return nil
}

func (u CreateEvalCompletionsRunDataSourceSourceUnion) AsFileContent() (v CreateEvalCompletionsRunDataSourceSourceFileContent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u CreateEvalCompletionsRunDataSourceSourceUnion) AsFileID() (v CreateEvalCompletionsRunDataSourceSourceFileID) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u CreateEvalCompletionsRunDataSourceSourceUnion) AsStoredCompletions() (v CreateEvalCompletionsRunDataSourceSourceStoredCompletions) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u CreateEvalCompletionsRunDataSourceSourceUnion) RawJSON() string { return u.JSON.raw }

func (r *CreateEvalCompletionsRunDataSourceSourceUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type CreateEvalCompletionsRunDataSourceSourceFileContent struct {
	// The content of the jsonl file.
	Content []CreateEvalCompletionsRunDataSourceSourceFileContentContent `json:"content,required"`
	// The type of jsonl source. Always `file_content`.
	Type constant.FileContent `json:"type,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		Content     resp.Field
		Type        resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CreateEvalCompletionsRunDataSourceSourceFileContent) RawJSON() string { return r.JSON.raw }
func (r *CreateEvalCompletionsRunDataSourceSourceFileContent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type CreateEvalCompletionsRunDataSourceSourceFileContentContent struct {
	Item   map[string]interface{} `json:"item,required"`
	Sample map[string]interface{} `json:"sample"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		Item        resp.Field
		Sample      resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CreateEvalCompletionsRunDataSourceSourceFileContentContent) RawJSON() string {
	return r.JSON.raw
}
func (r *CreateEvalCompletionsRunDataSourceSourceFileContentContent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type CreateEvalCompletionsRunDataSourceSourceFileID struct {
	// The identifier of the file.
	ID string `json:"id,required"`
	// The type of jsonl source. Always `file_id`.
	Type constant.FileID `json:"type,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		ID          resp.Field
		Type        resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CreateEvalCompletionsRunDataSourceSourceFileID) RawJSON() string { return r.JSON.raw }
func (r *CreateEvalCompletionsRunDataSourceSourceFileID) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A StoredCompletionsRunDataSource configuration describing a set of filters
type CreateEvalCompletionsRunDataSourceSourceStoredCompletions struct {
	// An optional Unix timestamp to filter items created after this time.
	CreatedAfter int64 `json:"created_after,required"`
	// An optional Unix timestamp to filter items created before this time.
	CreatedBefore int64 `json:"created_before,required"`
	// An optional maximum number of items to return.
	Limit int64 `json:"limit,required"`
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard.
	//
	// Keys are strings with a maximum length of 64 characters. Values are strings with
	// a maximum length of 512 characters.
	Metadata shared.Metadata `json:"metadata,required"`
	// An optional model to filter by (e.g., 'gpt-4o').
	Model string `json:"model,required"`
	// The type of source. Always `stored_completions`.
	Type constant.StoredCompletions `json:"type,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		CreatedAfter  resp.Field
		CreatedBefore resp.Field
		Limit         resp.Field
		Metadata      resp.Field
		Model         resp.Field
		Type          resp.Field
		ExtraFields   map[string]resp.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CreateEvalCompletionsRunDataSourceSourceStoredCompletions) RawJSON() string {
	return r.JSON.raw
}
func (r *CreateEvalCompletionsRunDataSourceSourceStoredCompletions) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The type of run data source. Always `completions`.
type CreateEvalCompletionsRunDataSourceType string

const (
	CreateEvalCompletionsRunDataSourceTypeCompletions CreateEvalCompletionsRunDataSourceType = "completions"
)

type CreateEvalCompletionsRunDataSourceSamplingParams struct {
	// The maximum number of tokens in the generated output.
	MaxCompletionTokens int64 `json:"max_completion_tokens"`
	// A seed value to initialize the randomness, during sampling.
	Seed int64 `json:"seed"`
	// A higher temperature increases randomness in the outputs.
	Temperature float64 `json:"temperature"`
	// An alternative to temperature for nucleus sampling; 1.0 includes all tokens.
	TopP float64 `json:"top_p"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		MaxCompletionTokens resp.Field
		Seed                resp.Field
		Temperature         resp.Field
		TopP                resp.Field
		ExtraFields         map[string]resp.Field
		raw                 string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CreateEvalCompletionsRunDataSourceSamplingParams) RawJSON() string { return r.JSON.raw }
func (r *CreateEvalCompletionsRunDataSourceSamplingParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A CompletionsRunDataSource object describing a model sampling configuration.
//
// The properties InputMessages, Model, Source, Type are required.
type CreateEvalCompletionsRunDataSourceParam struct {
	InputMessages CreateEvalCompletionsRunDataSourceInputMessagesUnionParam `json:"input_messages,omitzero,required"`
	// The name of the model to use for generating completions (e.g. "o3-mini").
	Model string `json:"model,required"`
	// A StoredCompletionsRunDataSource configuration describing a set of filters
	Source CreateEvalCompletionsRunDataSourceSourceUnionParam `json:"source,omitzero,required"`
	// The type of run data source. Always `completions`.
	//
	// Any of "completions".
	Type           CreateEvalCompletionsRunDataSourceType                `json:"type,omitzero,required"`
	SamplingParams CreateEvalCompletionsRunDataSourceSamplingParamsParam `json:"sampling_params,omitzero"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f CreateEvalCompletionsRunDataSourceParam) IsPresent() bool {
	return !param.IsOmitted(f) && !f.IsNull()
}
func (r CreateEvalCompletionsRunDataSourceParam) MarshalJSON() (data []byte, err error) {
	type shadow CreateEvalCompletionsRunDataSourceParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type CreateEvalCompletionsRunDataSourceInputMessagesUnionParam struct {
	OfTemplate      *CreateEvalCompletionsRunDataSourceInputMessagesTemplateParam      `json:",omitzero,inline"`
	OfItemReference *CreateEvalCompletionsRunDataSourceInputMessagesItemReferenceParam `json:",omitzero,inline"`
	paramUnion
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (u CreateEvalCompletionsRunDataSourceInputMessagesUnionParam) IsPresent() bool {
	return !param.IsOmitted(u) && !u.IsNull()
}
func (u CreateEvalCompletionsRunDataSourceInputMessagesUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion[CreateEvalCompletionsRunDataSourceInputMessagesUnionParam](u.OfTemplate, u.OfItemReference)
}

func (u *CreateEvalCompletionsRunDataSourceInputMessagesUnionParam) asAny() any {
	if !param.IsOmitted(u.OfTemplate) {
		return u.OfTemplate
	} else if !param.IsOmitted(u.OfItemReference) {
		return u.OfItemReference
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u CreateEvalCompletionsRunDataSourceInputMessagesUnionParam) GetTemplate() []CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateUnionParam {
	if vt := u.OfTemplate; vt != nil {
		return vt.Template
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u CreateEvalCompletionsRunDataSourceInputMessagesUnionParam) GetItemReference() *string {
	if vt := u.OfItemReference; vt != nil {
		return &vt.ItemReference
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u CreateEvalCompletionsRunDataSourceInputMessagesUnionParam) GetType() *string {
	if vt := u.OfTemplate; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfItemReference; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[CreateEvalCompletionsRunDataSourceInputMessagesUnionParam](
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(CreateEvalCompletionsRunDataSourceInputMessagesTemplateParam{}),
			DiscriminatorValue: "template",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(CreateEvalCompletionsRunDataSourceInputMessagesItemReferenceParam{}),
			DiscriminatorValue: "item_reference",
		},
	)
}

// The properties Template, Type are required.
type CreateEvalCompletionsRunDataSourceInputMessagesTemplateParam struct {
	// A list of chat messages forming the prompt or context. May include variable
	// references to the "item" namespace, ie {{item.name}}.
	Template []CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateUnionParam `json:"template,omitzero,required"`
	// The type of input messages. Always `template`.
	//
	// This field can be elided, and will marshal its zero value as "template".
	Type constant.Template `json:"type,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f CreateEvalCompletionsRunDataSourceInputMessagesTemplateParam) IsPresent() bool {
	return !param.IsOmitted(f) && !f.IsNull()
}
func (r CreateEvalCompletionsRunDataSourceInputMessagesTemplateParam) MarshalJSON() (data []byte, err error) {
	type shadow CreateEvalCompletionsRunDataSourceInputMessagesTemplateParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateUnionParam struct {
	OfChatMessage   *CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateChatMessageParam   `json:",omitzero,inline"`
	OfInputMessage  *CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateInputMessageParam  `json:",omitzero,inline"`
	OfOutputMessage *CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateOutputMessageParam `json:",omitzero,inline"`
	paramUnion
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (u CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateUnionParam) IsPresent() bool {
	return !param.IsOmitted(u) && !u.IsNull()
}
func (u CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion[CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateUnionParam](u.OfChatMessage, u.OfInputMessage, u.OfOutputMessage)
}

func (u *CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateUnionParam) asAny() any {
	if !param.IsOmitted(u.OfChatMessage) {
		return u.OfChatMessage
	} else if !param.IsOmitted(u.OfInputMessage) {
		return u.OfInputMessage
	} else if !param.IsOmitted(u.OfOutputMessage) {
		return u.OfOutputMessage
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateUnionParam) GetRole() *string {
	if vt := u.OfChatMessage; vt != nil {
		return (*string)(&vt.Role)
	} else if vt := u.OfInputMessage; vt != nil {
		return (*string)(&vt.Role)
	} else if vt := u.OfOutputMessage; vt != nil {
		return (*string)(&vt.Role)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateUnionParam) GetType() *string {
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
func (u CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateUnionParam) GetContent() (res createEvalCompletionsRunDataSourceInputMessagesTemplateTemplateUnionParamContent) {
	if vt := u.OfChatMessage; vt != nil {
		res.ofString = &vt.Content
	} else if vt := u.OfInputMessage; vt != nil {
		res.ofCreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateInputMessageContent = &vt.Content
	} else if vt := u.OfOutputMessage; vt != nil {
		res.ofCreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateOutputMessageContent = &vt.Content
	}
	return
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type createEvalCompletionsRunDataSourceInputMessagesTemplateTemplateUnionParamContent struct {
	ofString                                                                              *string
	ofCreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateInputMessageContent  *CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateInputMessageContentParam
	ofCreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateOutputMessageContent *CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateOutputMessageContentParam
}

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *string:
//	case *openai.CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateInputMessageContentParam:
//	case *openai.CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateOutputMessageContentParam:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u createEvalCompletionsRunDataSourceInputMessagesTemplateTemplateUnionParamContent) AsAny() any {
	if !param.IsOmitted(u.ofString) {
		return u.ofString
	} else if !param.IsOmitted(u.ofCreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateInputMessageContent) {
		return u.ofCreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateInputMessageContent
	} else if !param.IsOmitted(u.ofCreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateOutputMessageContent) {
		return u.ofCreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateOutputMessageContent
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u createEvalCompletionsRunDataSourceInputMessagesTemplateTemplateUnionParamContent) GetText() *string {
	if vt := u.ofCreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateInputMessageContent; vt != nil {
		return (*string)(&vt.Text)
	} else if vt := u.ofCreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateOutputMessageContent; vt != nil {
		return (*string)(&vt.Text)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u createEvalCompletionsRunDataSourceInputMessagesTemplateTemplateUnionParamContent) GetType() *string {
	if vt := u.ofCreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateInputMessageContent; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.ofCreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateOutputMessageContent; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// The properties Content, Role are required.
type CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateChatMessageParam struct {
	// The content of the message.
	Content string `json:"content,required"`
	// The role of the message (e.g. "system", "assistant", "user").
	Role string `json:"role,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateChatMessageParam) IsPresent() bool {
	return !param.IsOmitted(f) && !f.IsNull()
}
func (r CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateChatMessageParam) MarshalJSON() (data []byte, err error) {
	type shadow CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateChatMessageParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// The properties Content, Role, Type are required.
type CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateInputMessageParam struct {
	Content CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateInputMessageContentParam `json:"content,omitzero,required"`
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
func (f CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateInputMessageParam) IsPresent() bool {
	return !param.IsOmitted(f) && !f.IsNull()
}
func (r CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateInputMessageParam) MarshalJSON() (data []byte, err error) {
	type shadow CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateInputMessageParam
	return param.MarshalObject(r, (*shadow)(&r))
}

func init() {
	apijson.RegisterFieldValidator[CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateInputMessageParam](
		"Role", false, "user", "system", "developer",
	)
	apijson.RegisterFieldValidator[CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateInputMessageParam](
		"Type", false, "message",
	)
}

// The properties Text, Type are required.
type CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateInputMessageContentParam struct {
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
func (f CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateInputMessageContentParam) IsPresent() bool {
	return !param.IsOmitted(f) && !f.IsNull()
}
func (r CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateInputMessageContentParam) MarshalJSON() (data []byte, err error) {
	type shadow CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateInputMessageContentParam
	return param.MarshalObject(r, (*shadow)(&r))
}

func init() {
	apijson.RegisterFieldValidator[CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateInputMessageContentParam](
		"Type", false, "input_text",
	)
}

// The properties Content, Role, Type are required.
type CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateOutputMessageParam struct {
	Content CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateOutputMessageContentParam `json:"content,omitzero,required"`
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
func (f CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateOutputMessageParam) IsPresent() bool {
	return !param.IsOmitted(f) && !f.IsNull()
}
func (r CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateOutputMessageParam) MarshalJSON() (data []byte, err error) {
	type shadow CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateOutputMessageParam
	return param.MarshalObject(r, (*shadow)(&r))
}

func init() {
	apijson.RegisterFieldValidator[CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateOutputMessageParam](
		"Role", false, "assistant",
	)
	apijson.RegisterFieldValidator[CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateOutputMessageParam](
		"Type", false, "message",
	)
}

// The properties Text, Type are required.
type CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateOutputMessageContentParam struct {
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
func (f CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateOutputMessageContentParam) IsPresent() bool {
	return !param.IsOmitted(f) && !f.IsNull()
}
func (r CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateOutputMessageContentParam) MarshalJSON() (data []byte, err error) {
	type shadow CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateOutputMessageContentParam
	return param.MarshalObject(r, (*shadow)(&r))
}

func init() {
	apijson.RegisterFieldValidator[CreateEvalCompletionsRunDataSourceInputMessagesTemplateTemplateOutputMessageContentParam](
		"Type", false, "output_text",
	)
}

// The properties ItemReference, Type are required.
type CreateEvalCompletionsRunDataSourceInputMessagesItemReferenceParam struct {
	// A reference to a variable in the "item" namespace. Ie, "item.name"
	ItemReference string `json:"item_reference,required"`
	// The type of input messages. Always `item_reference`.
	//
	// This field can be elided, and will marshal its zero value as "item_reference".
	Type constant.ItemReference `json:"type,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f CreateEvalCompletionsRunDataSourceInputMessagesItemReferenceParam) IsPresent() bool {
	return !param.IsOmitted(f) && !f.IsNull()
}
func (r CreateEvalCompletionsRunDataSourceInputMessagesItemReferenceParam) MarshalJSON() (data []byte, err error) {
	type shadow CreateEvalCompletionsRunDataSourceInputMessagesItemReferenceParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type CreateEvalCompletionsRunDataSourceSourceUnionParam struct {
	OfFileContent       *CreateEvalCompletionsRunDataSourceSourceFileContentParam       `json:",omitzero,inline"`
	OfFileID            *CreateEvalCompletionsRunDataSourceSourceFileIDParam            `json:",omitzero,inline"`
	OfStoredCompletions *CreateEvalCompletionsRunDataSourceSourceStoredCompletionsParam `json:",omitzero,inline"`
	paramUnion
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (u CreateEvalCompletionsRunDataSourceSourceUnionParam) IsPresent() bool {
	return !param.IsOmitted(u) && !u.IsNull()
}
func (u CreateEvalCompletionsRunDataSourceSourceUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion[CreateEvalCompletionsRunDataSourceSourceUnionParam](u.OfFileContent, u.OfFileID, u.OfStoredCompletions)
}

func (u *CreateEvalCompletionsRunDataSourceSourceUnionParam) asAny() any {
	if !param.IsOmitted(u.OfFileContent) {
		return u.OfFileContent
	} else if !param.IsOmitted(u.OfFileID) {
		return u.OfFileID
	} else if !param.IsOmitted(u.OfStoredCompletions) {
		return u.OfStoredCompletions
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u CreateEvalCompletionsRunDataSourceSourceUnionParam) GetContent() []CreateEvalCompletionsRunDataSourceSourceFileContentContentParam {
	if vt := u.OfFileContent; vt != nil {
		return vt.Content
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u CreateEvalCompletionsRunDataSourceSourceUnionParam) GetID() *string {
	if vt := u.OfFileID; vt != nil {
		return &vt.ID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u CreateEvalCompletionsRunDataSourceSourceUnionParam) GetCreatedAfter() *int64 {
	if vt := u.OfStoredCompletions; vt != nil && vt.CreatedAfter.IsPresent() {
		return &vt.CreatedAfter.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u CreateEvalCompletionsRunDataSourceSourceUnionParam) GetCreatedBefore() *int64 {
	if vt := u.OfStoredCompletions; vt != nil && vt.CreatedBefore.IsPresent() {
		return &vt.CreatedBefore.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u CreateEvalCompletionsRunDataSourceSourceUnionParam) GetLimit() *int64 {
	if vt := u.OfStoredCompletions; vt != nil && vt.Limit.IsPresent() {
		return &vt.Limit.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u CreateEvalCompletionsRunDataSourceSourceUnionParam) GetMetadata() shared.MetadataParam {
	if vt := u.OfStoredCompletions; vt != nil {
		return vt.Metadata
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u CreateEvalCompletionsRunDataSourceSourceUnionParam) GetModel() *string {
	if vt := u.OfStoredCompletions; vt != nil && vt.Model.IsPresent() {
		return &vt.Model.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u CreateEvalCompletionsRunDataSourceSourceUnionParam) GetType() *string {
	if vt := u.OfFileContent; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfFileID; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfStoredCompletions; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[CreateEvalCompletionsRunDataSourceSourceUnionParam](
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(CreateEvalCompletionsRunDataSourceSourceFileContentParam{}),
			DiscriminatorValue: "file_content",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(CreateEvalCompletionsRunDataSourceSourceFileIDParam{}),
			DiscriminatorValue: "file_id",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(CreateEvalCompletionsRunDataSourceSourceStoredCompletionsParam{}),
			DiscriminatorValue: "stored_completions",
		},
	)
}

// The properties Content, Type are required.
type CreateEvalCompletionsRunDataSourceSourceFileContentParam struct {
	// The content of the jsonl file.
	Content []CreateEvalCompletionsRunDataSourceSourceFileContentContentParam `json:"content,omitzero,required"`
	// The type of jsonl source. Always `file_content`.
	//
	// This field can be elided, and will marshal its zero value as "file_content".
	Type constant.FileContent `json:"type,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f CreateEvalCompletionsRunDataSourceSourceFileContentParam) IsPresent() bool {
	return !param.IsOmitted(f) && !f.IsNull()
}
func (r CreateEvalCompletionsRunDataSourceSourceFileContentParam) MarshalJSON() (data []byte, err error) {
	type shadow CreateEvalCompletionsRunDataSourceSourceFileContentParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// The property Item is required.
type CreateEvalCompletionsRunDataSourceSourceFileContentContentParam struct {
	Item   map[string]interface{} `json:"item,omitzero,required"`
	Sample map[string]interface{} `json:"sample,omitzero"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f CreateEvalCompletionsRunDataSourceSourceFileContentContentParam) IsPresent() bool {
	return !param.IsOmitted(f) && !f.IsNull()
}
func (r CreateEvalCompletionsRunDataSourceSourceFileContentContentParam) MarshalJSON() (data []byte, err error) {
	type shadow CreateEvalCompletionsRunDataSourceSourceFileContentContentParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// The properties ID, Type are required.
type CreateEvalCompletionsRunDataSourceSourceFileIDParam struct {
	// The identifier of the file.
	ID string `json:"id,required"`
	// The type of jsonl source. Always `file_id`.
	//
	// This field can be elided, and will marshal its zero value as "file_id".
	Type constant.FileID `json:"type,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f CreateEvalCompletionsRunDataSourceSourceFileIDParam) IsPresent() bool {
	return !param.IsOmitted(f) && !f.IsNull()
}
func (r CreateEvalCompletionsRunDataSourceSourceFileIDParam) MarshalJSON() (data []byte, err error) {
	type shadow CreateEvalCompletionsRunDataSourceSourceFileIDParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// A StoredCompletionsRunDataSource configuration describing a set of filters
//
// The properties CreatedAfter, CreatedBefore, Limit, Metadata, Model, Type are
// required.
type CreateEvalCompletionsRunDataSourceSourceStoredCompletionsParam struct {
	// An optional Unix timestamp to filter items created after this time.
	CreatedAfter param.Opt[int64] `json:"created_after,omitzero,required"`
	// An optional Unix timestamp to filter items created before this time.
	CreatedBefore param.Opt[int64] `json:"created_before,omitzero,required"`
	// An optional maximum number of items to return.
	Limit param.Opt[int64] `json:"limit,omitzero,required"`
	// An optional model to filter by (e.g., 'gpt-4o').
	Model param.Opt[string] `json:"model,omitzero,required"`
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard.
	//
	// Keys are strings with a maximum length of 64 characters. Values are strings with
	// a maximum length of 512 characters.
	Metadata shared.MetadataParam `json:"metadata,omitzero,required"`
	// The type of source. Always `stored_completions`.
	//
	// This field can be elided, and will marshal its zero value as
	// "stored_completions".
	Type constant.StoredCompletions `json:"type,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f CreateEvalCompletionsRunDataSourceSourceStoredCompletionsParam) IsPresent() bool {
	return !param.IsOmitted(f) && !f.IsNull()
}
func (r CreateEvalCompletionsRunDataSourceSourceStoredCompletionsParam) MarshalJSON() (data []byte, err error) {
	type shadow CreateEvalCompletionsRunDataSourceSourceStoredCompletionsParam
	return param.MarshalObject(r, (*shadow)(&r))
}

type CreateEvalCompletionsRunDataSourceSamplingParamsParam struct {
	// The maximum number of tokens in the generated output.
	MaxCompletionTokens param.Opt[int64] `json:"max_completion_tokens,omitzero"`
	// A seed value to initialize the randomness, during sampling.
	Seed param.Opt[int64] `json:"seed,omitzero"`
	// A higher temperature increases randomness in the outputs.
	Temperature param.Opt[float64] `json:"temperature,omitzero"`
	// An alternative to temperature for nucleus sampling; 1.0 includes all tokens.
	TopP param.Opt[float64] `json:"top_p,omitzero"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f CreateEvalCompletionsRunDataSourceSamplingParamsParam) IsPresent() bool {
	return !param.IsOmitted(f) && !f.IsNull()
}
func (r CreateEvalCompletionsRunDataSourceSamplingParamsParam) MarshalJSON() (data []byte, err error) {
	type shadow CreateEvalCompletionsRunDataSourceSamplingParamsParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// A JsonlRunDataSource object with that specifies a JSONL file that matches the
// eval
type CreateEvalJSONLRunDataSource struct {
	Source CreateEvalJSONLRunDataSourceSourceUnion `json:"source,required"`
	// The type of data source. Always `jsonl`.
	Type constant.JSONL `json:"type,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		Source      resp.Field
		Type        resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CreateEvalJSONLRunDataSource) RawJSON() string { return r.JSON.raw }
func (r *CreateEvalJSONLRunDataSource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this CreateEvalJSONLRunDataSource to a
// CreateEvalJSONLRunDataSourceParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// CreateEvalJSONLRunDataSourceParam.IsOverridden()
func (r CreateEvalJSONLRunDataSource) ToParam() CreateEvalJSONLRunDataSourceParam {
	return param.OverrideObj[CreateEvalJSONLRunDataSourceParam](r.RawJSON())
}

// CreateEvalJSONLRunDataSourceSourceUnion contains all possible properties and
// values from [CreateEvalJSONLRunDataSourceSourceFileContent],
// [CreateEvalJSONLRunDataSourceSourceFileID].
//
// Use the [CreateEvalJSONLRunDataSourceSourceUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type CreateEvalJSONLRunDataSourceSourceUnion struct {
	// This field is from variant [CreateEvalJSONLRunDataSourceSourceFileContent].
	Content []CreateEvalJSONLRunDataSourceSourceFileContentContent `json:"content"`
	// Any of "file_content", "file_id".
	Type string `json:"type"`
	// This field is from variant [CreateEvalJSONLRunDataSourceSourceFileID].
	ID   string `json:"id"`
	JSON struct {
		Content resp.Field
		Type    resp.Field
		ID      resp.Field
		raw     string
	} `json:"-"`
}

// anyCreateEvalJSONLRunDataSourceSource is implemented by each variant of
// [CreateEvalJSONLRunDataSourceSourceUnion] to add type safety for the return type
// of [CreateEvalJSONLRunDataSourceSourceUnion.AsAny]
type anyCreateEvalJSONLRunDataSourceSource interface {
	implCreateEvalJSONLRunDataSourceSourceUnion()
}

func (CreateEvalJSONLRunDataSourceSourceFileContent) implCreateEvalJSONLRunDataSourceSourceUnion() {}
func (CreateEvalJSONLRunDataSourceSourceFileID) implCreateEvalJSONLRunDataSourceSourceUnion()      {}

// Use the following switch statement to find the correct variant
//
//	switch variant := CreateEvalJSONLRunDataSourceSourceUnion.AsAny().(type) {
//	case CreateEvalJSONLRunDataSourceSourceFileContent:
//	case CreateEvalJSONLRunDataSourceSourceFileID:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u CreateEvalJSONLRunDataSourceSourceUnion) AsAny() anyCreateEvalJSONLRunDataSourceSource {
	switch u.Type {
	case "file_content":
		return u.AsFileContent()
	case "file_id":
		return u.AsFileID()
	}
	return nil
}

func (u CreateEvalJSONLRunDataSourceSourceUnion) AsFileContent() (v CreateEvalJSONLRunDataSourceSourceFileContent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u CreateEvalJSONLRunDataSourceSourceUnion) AsFileID() (v CreateEvalJSONLRunDataSourceSourceFileID) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u CreateEvalJSONLRunDataSourceSourceUnion) RawJSON() string { return u.JSON.raw }

func (r *CreateEvalJSONLRunDataSourceSourceUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type CreateEvalJSONLRunDataSourceSourceFileContent struct {
	// The content of the jsonl file.
	Content []CreateEvalJSONLRunDataSourceSourceFileContentContent `json:"content,required"`
	// The type of jsonl source. Always `file_content`.
	Type constant.FileContent `json:"type,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		Content     resp.Field
		Type        resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CreateEvalJSONLRunDataSourceSourceFileContent) RawJSON() string { return r.JSON.raw }
func (r *CreateEvalJSONLRunDataSourceSourceFileContent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type CreateEvalJSONLRunDataSourceSourceFileContentContent struct {
	Item   map[string]interface{} `json:"item,required"`
	Sample map[string]interface{} `json:"sample"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		Item        resp.Field
		Sample      resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CreateEvalJSONLRunDataSourceSourceFileContentContent) RawJSON() string { return r.JSON.raw }
func (r *CreateEvalJSONLRunDataSourceSourceFileContentContent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type CreateEvalJSONLRunDataSourceSourceFileID struct {
	// The identifier of the file.
	ID string `json:"id,required"`
	// The type of jsonl source. Always `file_id`.
	Type constant.FileID `json:"type,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		ID          resp.Field
		Type        resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CreateEvalJSONLRunDataSourceSourceFileID) RawJSON() string { return r.JSON.raw }
func (r *CreateEvalJSONLRunDataSourceSourceFileID) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A JsonlRunDataSource object with that specifies a JSONL file that matches the
// eval
//
// The properties Source, Type are required.
type CreateEvalJSONLRunDataSourceParam struct {
	Source CreateEvalJSONLRunDataSourceSourceUnionParam `json:"source,omitzero,required"`
	// The type of data source. Always `jsonl`.
	//
	// This field can be elided, and will marshal its zero value as "jsonl".
	Type constant.JSONL `json:"type,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f CreateEvalJSONLRunDataSourceParam) IsPresent() bool {
	return !param.IsOmitted(f) && !f.IsNull()
}
func (r CreateEvalJSONLRunDataSourceParam) MarshalJSON() (data []byte, err error) {
	type shadow CreateEvalJSONLRunDataSourceParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type CreateEvalJSONLRunDataSourceSourceUnionParam struct {
	OfFileContent *CreateEvalJSONLRunDataSourceSourceFileContentParam `json:",omitzero,inline"`
	OfFileID      *CreateEvalJSONLRunDataSourceSourceFileIDParam      `json:",omitzero,inline"`
	paramUnion
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (u CreateEvalJSONLRunDataSourceSourceUnionParam) IsPresent() bool {
	return !param.IsOmitted(u) && !u.IsNull()
}
func (u CreateEvalJSONLRunDataSourceSourceUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion[CreateEvalJSONLRunDataSourceSourceUnionParam](u.OfFileContent, u.OfFileID)
}

func (u *CreateEvalJSONLRunDataSourceSourceUnionParam) asAny() any {
	if !param.IsOmitted(u.OfFileContent) {
		return u.OfFileContent
	} else if !param.IsOmitted(u.OfFileID) {
		return u.OfFileID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u CreateEvalJSONLRunDataSourceSourceUnionParam) GetContent() []CreateEvalJSONLRunDataSourceSourceFileContentContentParam {
	if vt := u.OfFileContent; vt != nil {
		return vt.Content
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u CreateEvalJSONLRunDataSourceSourceUnionParam) GetID() *string {
	if vt := u.OfFileID; vt != nil {
		return &vt.ID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u CreateEvalJSONLRunDataSourceSourceUnionParam) GetType() *string {
	if vt := u.OfFileContent; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfFileID; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[CreateEvalJSONLRunDataSourceSourceUnionParam](
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(CreateEvalJSONLRunDataSourceSourceFileContentParam{}),
			DiscriminatorValue: "file_content",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(CreateEvalJSONLRunDataSourceSourceFileIDParam{}),
			DiscriminatorValue: "file_id",
		},
	)
}

// The properties Content, Type are required.
type CreateEvalJSONLRunDataSourceSourceFileContentParam struct {
	// The content of the jsonl file.
	Content []CreateEvalJSONLRunDataSourceSourceFileContentContentParam `json:"content,omitzero,required"`
	// The type of jsonl source. Always `file_content`.
	//
	// This field can be elided, and will marshal its zero value as "file_content".
	Type constant.FileContent `json:"type,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f CreateEvalJSONLRunDataSourceSourceFileContentParam) IsPresent() bool {
	return !param.IsOmitted(f) && !f.IsNull()
}
func (r CreateEvalJSONLRunDataSourceSourceFileContentParam) MarshalJSON() (data []byte, err error) {
	type shadow CreateEvalJSONLRunDataSourceSourceFileContentParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// The property Item is required.
type CreateEvalJSONLRunDataSourceSourceFileContentContentParam struct {
	Item   map[string]interface{} `json:"item,omitzero,required"`
	Sample map[string]interface{} `json:"sample,omitzero"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f CreateEvalJSONLRunDataSourceSourceFileContentContentParam) IsPresent() bool {
	return !param.IsOmitted(f) && !f.IsNull()
}
func (r CreateEvalJSONLRunDataSourceSourceFileContentContentParam) MarshalJSON() (data []byte, err error) {
	type shadow CreateEvalJSONLRunDataSourceSourceFileContentContentParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// The properties ID, Type are required.
type CreateEvalJSONLRunDataSourceSourceFileIDParam struct {
	// The identifier of the file.
	ID string `json:"id,required"`
	// The type of jsonl source. Always `file_id`.
	//
	// This field can be elided, and will marshal its zero value as "file_id".
	Type constant.FileID `json:"type,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f CreateEvalJSONLRunDataSourceSourceFileIDParam) IsPresent() bool {
	return !param.IsOmitted(f) && !f.IsNull()
}
func (r CreateEvalJSONLRunDataSourceSourceFileIDParam) MarshalJSON() (data []byte, err error) {
	type shadow CreateEvalJSONLRunDataSourceSourceFileIDParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// An object representing an error response from the Eval API.
type EvalAPIError struct {
	// The error code.
	Code string `json:"code,required"`
	// The error message.
	Message string `json:"message,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		Code        resp.Field
		Message     resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EvalAPIError) RawJSON() string { return r.JSON.raw }
func (r *EvalAPIError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A schema representing an evaluation run.
type EvalRunNewResponse struct {
	// Unique identifier for the evaluation run.
	ID string `json:"id,required"`
	// Unix timestamp (in seconds) when the evaluation run was created.
	CreatedAt int64 `json:"created_at,required"`
	// Information about the run's data source.
	DataSource EvalRunNewResponseDataSourceUnion `json:"data_source,required"`
	// An object representing an error response from the Eval API.
	Error EvalAPIError `json:"error,required"`
	// The identifier of the associated evaluation.
	EvalID string `json:"eval_id,required"`
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard.
	//
	// Keys are strings with a maximum length of 64 characters. Values are strings with
	// a maximum length of 512 characters.
	Metadata shared.Metadata `json:"metadata,required"`
	// The model that is evaluated, if applicable.
	Model string `json:"model,required"`
	// The name of the evaluation run.
	Name string `json:"name,required"`
	// The type of the object. Always "eval.run".
	Object constant.EvalRun `json:"object,required"`
	// Usage statistics for each model during the evaluation run.
	PerModelUsage []EvalRunNewResponsePerModelUsage `json:"per_model_usage,required"`
	// Results per testing criteria applied during the evaluation run.
	PerTestingCriteriaResults []EvalRunNewResponsePerTestingCriteriaResult `json:"per_testing_criteria_results,required"`
	// The URL to the rendered evaluation run report on the UI dashboard.
	ReportURL string `json:"report_url,required"`
	// Counters summarizing the outcomes of the evaluation run.
	ResultCounts EvalRunNewResponseResultCounts `json:"result_counts,required"`
	// The status of the evaluation run.
	Status string `json:"status,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		ID                        resp.Field
		CreatedAt                 resp.Field
		DataSource                resp.Field
		Error                     resp.Field
		EvalID                    resp.Field
		Metadata                  resp.Field
		Model                     resp.Field
		Name                      resp.Field
		Object                    resp.Field
		PerModelUsage             resp.Field
		PerTestingCriteriaResults resp.Field
		ReportURL                 resp.Field
		ResultCounts              resp.Field
		Status                    resp.Field
		ExtraFields               map[string]resp.Field
		raw                       string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EvalRunNewResponse) RawJSON() string { return r.JSON.raw }
func (r *EvalRunNewResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// EvalRunNewResponseDataSourceUnion contains all possible properties and values
// from [CreateEvalJSONLRunDataSource], [CreateEvalCompletionsRunDataSource].
//
// Use the [EvalRunNewResponseDataSourceUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type EvalRunNewResponseDataSourceUnion struct {
	// This field is a union of [CreateEvalJSONLRunDataSourceSourceUnion],
	// [CreateEvalCompletionsRunDataSourceSourceUnion]
	Source EvalRunNewResponseDataSourceUnionSource `json:"source"`
	// Any of "jsonl", "completions".
	Type string `json:"type"`
	// This field is from variant [CreateEvalCompletionsRunDataSource].
	InputMessages CreateEvalCompletionsRunDataSourceInputMessagesUnion `json:"input_messages"`
	// This field is from variant [CreateEvalCompletionsRunDataSource].
	Model string `json:"model"`
	// This field is from variant [CreateEvalCompletionsRunDataSource].
	SamplingParams CreateEvalCompletionsRunDataSourceSamplingParams `json:"sampling_params"`
	JSON           struct {
		Source         resp.Field
		Type           resp.Field
		InputMessages  resp.Field
		Model          resp.Field
		SamplingParams resp.Field
		raw            string
	} `json:"-"`
}

// anyEvalRunNewResponseDataSource is implemented by each variant of
// [EvalRunNewResponseDataSourceUnion] to add type safety for the return type of
// [EvalRunNewResponseDataSourceUnion.AsAny]
type anyEvalRunNewResponseDataSource interface {
	implEvalRunNewResponseDataSourceUnion()
}

func (CreateEvalJSONLRunDataSource) implEvalRunNewResponseDataSourceUnion()       {}
func (CreateEvalCompletionsRunDataSource) implEvalRunNewResponseDataSourceUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := EvalRunNewResponseDataSourceUnion.AsAny().(type) {
//	case CreateEvalJSONLRunDataSource:
//	case CreateEvalCompletionsRunDataSource:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u EvalRunNewResponseDataSourceUnion) AsAny() anyEvalRunNewResponseDataSource {
	switch u.Type {
	case "jsonl":
		return u.AsJSONL()
	case "completions":
		return u.AsCompletions()
	}
	return nil
}

func (u EvalRunNewResponseDataSourceUnion) AsJSONL() (v CreateEvalJSONLRunDataSource) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u EvalRunNewResponseDataSourceUnion) AsCompletions() (v CreateEvalCompletionsRunDataSource) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u EvalRunNewResponseDataSourceUnion) RawJSON() string { return u.JSON.raw }

func (r *EvalRunNewResponseDataSourceUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// EvalRunNewResponseDataSourceUnionSource is an implicit subunion of
// [EvalRunNewResponseDataSourceUnion]. EvalRunNewResponseDataSourceUnionSource
// provides convenient access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [EvalRunNewResponseDataSourceUnion].
type EvalRunNewResponseDataSourceUnionSource struct {
	// This field is a union of
	// [[]CreateEvalJSONLRunDataSourceSourceFileContentContent],
	// [[]CreateEvalCompletionsRunDataSourceSourceFileContentContent]
	Content EvalRunNewResponseDataSourceUnionSourceContent `json:"content"`
	Type    string                                         `json:"type"`
	ID      string                                         `json:"id"`
	// This field is from variant [CreateEvalCompletionsRunDataSourceSourceUnion].
	CreatedAfter int64 `json:"created_after"`
	// This field is from variant [CreateEvalCompletionsRunDataSourceSourceUnion].
	CreatedBefore int64 `json:"created_before"`
	// This field is from variant [CreateEvalCompletionsRunDataSourceSourceUnion].
	Limit int64 `json:"limit"`
	// This field is from variant [CreateEvalCompletionsRunDataSourceSourceUnion].
	Metadata shared.Metadata `json:"metadata"`
	// This field is from variant [CreateEvalCompletionsRunDataSourceSourceUnion].
	Model string `json:"model"`
	JSON  struct {
		Content       resp.Field
		Type          resp.Field
		ID            resp.Field
		CreatedAfter  resp.Field
		CreatedBefore resp.Field
		Limit         resp.Field
		Metadata      resp.Field
		Model         resp.Field
		raw           string
	} `json:"-"`
}

func (r *EvalRunNewResponseDataSourceUnionSource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// EvalRunNewResponseDataSourceUnionSourceContent is an implicit subunion of
// [EvalRunNewResponseDataSourceUnion].
// EvalRunNewResponseDataSourceUnionSourceContent provides convenient access to the
// sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [EvalRunNewResponseDataSourceUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfCreateEvalJSONLRunDataSourceSourceFileContentContent
// OfCreateEvalCompletionsRunDataSourceSourceFileContentContent]
type EvalRunNewResponseDataSourceUnionSourceContent struct {
	// This field will be present if the value is a
	// [[]CreateEvalJSONLRunDataSourceSourceFileContentContent] instead of an object.
	OfCreateEvalJSONLRunDataSourceSourceFileContentContent []CreateEvalJSONLRunDataSourceSourceFileContentContent `json:",inline"`
	// This field will be present if the value is a
	// [[]CreateEvalCompletionsRunDataSourceSourceFileContentContent] instead of an
	// object.
	OfCreateEvalCompletionsRunDataSourceSourceFileContentContent []CreateEvalCompletionsRunDataSourceSourceFileContentContent `json:",inline"`
	JSON                                                         struct {
		OfCreateEvalJSONLRunDataSourceSourceFileContentContent       resp.Field
		OfCreateEvalCompletionsRunDataSourceSourceFileContentContent resp.Field
		raw                                                          string
	} `json:"-"`
}

func (r *EvalRunNewResponseDataSourceUnionSourceContent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type EvalRunNewResponsePerModelUsage struct {
	// The number of tokens retrieved from cache.
	CachedTokens int64 `json:"cached_tokens,required"`
	// The number of completion tokens generated.
	CompletionTokens int64 `json:"completion_tokens,required"`
	// The number of invocations.
	InvocationCount int64 `json:"invocation_count,required"`
	// The name of the model.
	ModelName string `json:"model_name,required"`
	// The number of prompt tokens used.
	PromptTokens int64 `json:"prompt_tokens,required"`
	// The total number of tokens used.
	TotalTokens int64 `json:"total_tokens,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		CachedTokens     resp.Field
		CompletionTokens resp.Field
		InvocationCount  resp.Field
		ModelName        resp.Field
		PromptTokens     resp.Field
		TotalTokens      resp.Field
		ExtraFields      map[string]resp.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EvalRunNewResponsePerModelUsage) RawJSON() string { return r.JSON.raw }
func (r *EvalRunNewResponsePerModelUsage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type EvalRunNewResponsePerTestingCriteriaResult struct {
	// Number of tests failed for this criteria.
	Failed int64 `json:"failed,required"`
	// Number of tests passed for this criteria.
	Passed int64 `json:"passed,required"`
	// A description of the testing criteria.
	TestingCriteria string `json:"testing_criteria,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		Failed          resp.Field
		Passed          resp.Field
		TestingCriteria resp.Field
		ExtraFields     map[string]resp.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EvalRunNewResponsePerTestingCriteriaResult) RawJSON() string { return r.JSON.raw }
func (r *EvalRunNewResponsePerTestingCriteriaResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Counters summarizing the outcomes of the evaluation run.
type EvalRunNewResponseResultCounts struct {
	// Number of output items that resulted in an error.
	Errored int64 `json:"errored,required"`
	// Number of output items that failed to pass the evaluation.
	Failed int64 `json:"failed,required"`
	// Number of output items that passed the evaluation.
	Passed int64 `json:"passed,required"`
	// Total number of executed output items.
	Total int64 `json:"total,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		Errored     resp.Field
		Failed      resp.Field
		Passed      resp.Field
		Total       resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EvalRunNewResponseResultCounts) RawJSON() string { return r.JSON.raw }
func (r *EvalRunNewResponseResultCounts) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A schema representing an evaluation run.
type EvalRunGetResponse struct {
	// Unique identifier for the evaluation run.
	ID string `json:"id,required"`
	// Unix timestamp (in seconds) when the evaluation run was created.
	CreatedAt int64 `json:"created_at,required"`
	// Information about the run's data source.
	DataSource EvalRunGetResponseDataSourceUnion `json:"data_source,required"`
	// An object representing an error response from the Eval API.
	Error EvalAPIError `json:"error,required"`
	// The identifier of the associated evaluation.
	EvalID string `json:"eval_id,required"`
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard.
	//
	// Keys are strings with a maximum length of 64 characters. Values are strings with
	// a maximum length of 512 characters.
	Metadata shared.Metadata `json:"metadata,required"`
	// The model that is evaluated, if applicable.
	Model string `json:"model,required"`
	// The name of the evaluation run.
	Name string `json:"name,required"`
	// The type of the object. Always "eval.run".
	Object constant.EvalRun `json:"object,required"`
	// Usage statistics for each model during the evaluation run.
	PerModelUsage []EvalRunGetResponsePerModelUsage `json:"per_model_usage,required"`
	// Results per testing criteria applied during the evaluation run.
	PerTestingCriteriaResults []EvalRunGetResponsePerTestingCriteriaResult `json:"per_testing_criteria_results,required"`
	// The URL to the rendered evaluation run report on the UI dashboard.
	ReportURL string `json:"report_url,required"`
	// Counters summarizing the outcomes of the evaluation run.
	ResultCounts EvalRunGetResponseResultCounts `json:"result_counts,required"`
	// The status of the evaluation run.
	Status string `json:"status,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		ID                        resp.Field
		CreatedAt                 resp.Field
		DataSource                resp.Field
		Error                     resp.Field
		EvalID                    resp.Field
		Metadata                  resp.Field
		Model                     resp.Field
		Name                      resp.Field
		Object                    resp.Field
		PerModelUsage             resp.Field
		PerTestingCriteriaResults resp.Field
		ReportURL                 resp.Field
		ResultCounts              resp.Field
		Status                    resp.Field
		ExtraFields               map[string]resp.Field
		raw                       string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EvalRunGetResponse) RawJSON() string { return r.JSON.raw }
func (r *EvalRunGetResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// EvalRunGetResponseDataSourceUnion contains all possible properties and values
// from [CreateEvalJSONLRunDataSource], [CreateEvalCompletionsRunDataSource].
//
// Use the [EvalRunGetResponseDataSourceUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type EvalRunGetResponseDataSourceUnion struct {
	// This field is a union of [CreateEvalJSONLRunDataSourceSourceUnion],
	// [CreateEvalCompletionsRunDataSourceSourceUnion]
	Source EvalRunGetResponseDataSourceUnionSource `json:"source"`
	// Any of "jsonl", "completions".
	Type string `json:"type"`
	// This field is from variant [CreateEvalCompletionsRunDataSource].
	InputMessages CreateEvalCompletionsRunDataSourceInputMessagesUnion `json:"input_messages"`
	// This field is from variant [CreateEvalCompletionsRunDataSource].
	Model string `json:"model"`
	// This field is from variant [CreateEvalCompletionsRunDataSource].
	SamplingParams CreateEvalCompletionsRunDataSourceSamplingParams `json:"sampling_params"`
	JSON           struct {
		Source         resp.Field
		Type           resp.Field
		InputMessages  resp.Field
		Model          resp.Field
		SamplingParams resp.Field
		raw            string
	} `json:"-"`
}

// anyEvalRunGetResponseDataSource is implemented by each variant of
// [EvalRunGetResponseDataSourceUnion] to add type safety for the return type of
// [EvalRunGetResponseDataSourceUnion.AsAny]
type anyEvalRunGetResponseDataSource interface {
	implEvalRunGetResponseDataSourceUnion()
}

func (CreateEvalJSONLRunDataSource) implEvalRunGetResponseDataSourceUnion()       {}
func (CreateEvalCompletionsRunDataSource) implEvalRunGetResponseDataSourceUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := EvalRunGetResponseDataSourceUnion.AsAny().(type) {
//	case CreateEvalJSONLRunDataSource:
//	case CreateEvalCompletionsRunDataSource:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u EvalRunGetResponseDataSourceUnion) AsAny() anyEvalRunGetResponseDataSource {
	switch u.Type {
	case "jsonl":
		return u.AsJSONL()
	case "completions":
		return u.AsCompletions()
	}
	return nil
}

func (u EvalRunGetResponseDataSourceUnion) AsJSONL() (v CreateEvalJSONLRunDataSource) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u EvalRunGetResponseDataSourceUnion) AsCompletions() (v CreateEvalCompletionsRunDataSource) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u EvalRunGetResponseDataSourceUnion) RawJSON() string { return u.JSON.raw }

func (r *EvalRunGetResponseDataSourceUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// EvalRunGetResponseDataSourceUnionSource is an implicit subunion of
// [EvalRunGetResponseDataSourceUnion]. EvalRunGetResponseDataSourceUnionSource
// provides convenient access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [EvalRunGetResponseDataSourceUnion].
type EvalRunGetResponseDataSourceUnionSource struct {
	// This field is a union of
	// [[]CreateEvalJSONLRunDataSourceSourceFileContentContent],
	// [[]CreateEvalCompletionsRunDataSourceSourceFileContentContent]
	Content EvalRunGetResponseDataSourceUnionSourceContent `json:"content"`
	Type    string                                         `json:"type"`
	ID      string                                         `json:"id"`
	// This field is from variant [CreateEvalCompletionsRunDataSourceSourceUnion].
	CreatedAfter int64 `json:"created_after"`
	// This field is from variant [CreateEvalCompletionsRunDataSourceSourceUnion].
	CreatedBefore int64 `json:"created_before"`
	// This field is from variant [CreateEvalCompletionsRunDataSourceSourceUnion].
	Limit int64 `json:"limit"`
	// This field is from variant [CreateEvalCompletionsRunDataSourceSourceUnion].
	Metadata shared.Metadata `json:"metadata"`
	// This field is from variant [CreateEvalCompletionsRunDataSourceSourceUnion].
	Model string `json:"model"`
	JSON  struct {
		Content       resp.Field
		Type          resp.Field
		ID            resp.Field
		CreatedAfter  resp.Field
		CreatedBefore resp.Field
		Limit         resp.Field
		Metadata      resp.Field
		Model         resp.Field
		raw           string
	} `json:"-"`
}

func (r *EvalRunGetResponseDataSourceUnionSource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// EvalRunGetResponseDataSourceUnionSourceContent is an implicit subunion of
// [EvalRunGetResponseDataSourceUnion].
// EvalRunGetResponseDataSourceUnionSourceContent provides convenient access to the
// sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [EvalRunGetResponseDataSourceUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfCreateEvalJSONLRunDataSourceSourceFileContentContent
// OfCreateEvalCompletionsRunDataSourceSourceFileContentContent]
type EvalRunGetResponseDataSourceUnionSourceContent struct {
	// This field will be present if the value is a
	// [[]CreateEvalJSONLRunDataSourceSourceFileContentContent] instead of an object.
	OfCreateEvalJSONLRunDataSourceSourceFileContentContent []CreateEvalJSONLRunDataSourceSourceFileContentContent `json:",inline"`
	// This field will be present if the value is a
	// [[]CreateEvalCompletionsRunDataSourceSourceFileContentContent] instead of an
	// object.
	OfCreateEvalCompletionsRunDataSourceSourceFileContentContent []CreateEvalCompletionsRunDataSourceSourceFileContentContent `json:",inline"`
	JSON                                                         struct {
		OfCreateEvalJSONLRunDataSourceSourceFileContentContent       resp.Field
		OfCreateEvalCompletionsRunDataSourceSourceFileContentContent resp.Field
		raw                                                          string
	} `json:"-"`
}

func (r *EvalRunGetResponseDataSourceUnionSourceContent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type EvalRunGetResponsePerModelUsage struct {
	// The number of tokens retrieved from cache.
	CachedTokens int64 `json:"cached_tokens,required"`
	// The number of completion tokens generated.
	CompletionTokens int64 `json:"completion_tokens,required"`
	// The number of invocations.
	InvocationCount int64 `json:"invocation_count,required"`
	// The name of the model.
	ModelName string `json:"model_name,required"`
	// The number of prompt tokens used.
	PromptTokens int64 `json:"prompt_tokens,required"`
	// The total number of tokens used.
	TotalTokens int64 `json:"total_tokens,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		CachedTokens     resp.Field
		CompletionTokens resp.Field
		InvocationCount  resp.Field
		ModelName        resp.Field
		PromptTokens     resp.Field
		TotalTokens      resp.Field
		ExtraFields      map[string]resp.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EvalRunGetResponsePerModelUsage) RawJSON() string { return r.JSON.raw }
func (r *EvalRunGetResponsePerModelUsage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type EvalRunGetResponsePerTestingCriteriaResult struct {
	// Number of tests failed for this criteria.
	Failed int64 `json:"failed,required"`
	// Number of tests passed for this criteria.
	Passed int64 `json:"passed,required"`
	// A description of the testing criteria.
	TestingCriteria string `json:"testing_criteria,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		Failed          resp.Field
		Passed          resp.Field
		TestingCriteria resp.Field
		ExtraFields     map[string]resp.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EvalRunGetResponsePerTestingCriteriaResult) RawJSON() string { return r.JSON.raw }
func (r *EvalRunGetResponsePerTestingCriteriaResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Counters summarizing the outcomes of the evaluation run.
type EvalRunGetResponseResultCounts struct {
	// Number of output items that resulted in an error.
	Errored int64 `json:"errored,required"`
	// Number of output items that failed to pass the evaluation.
	Failed int64 `json:"failed,required"`
	// Number of output items that passed the evaluation.
	Passed int64 `json:"passed,required"`
	// Total number of executed output items.
	Total int64 `json:"total,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		Errored     resp.Field
		Failed      resp.Field
		Passed      resp.Field
		Total       resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EvalRunGetResponseResultCounts) RawJSON() string { return r.JSON.raw }
func (r *EvalRunGetResponseResultCounts) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A schema representing an evaluation run.
type EvalRunListResponse struct {
	// Unique identifier for the evaluation run.
	ID string `json:"id,required"`
	// Unix timestamp (in seconds) when the evaluation run was created.
	CreatedAt int64 `json:"created_at,required"`
	// Information about the run's data source.
	DataSource EvalRunListResponseDataSourceUnion `json:"data_source,required"`
	// An object representing an error response from the Eval API.
	Error EvalAPIError `json:"error,required"`
	// The identifier of the associated evaluation.
	EvalID string `json:"eval_id,required"`
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard.
	//
	// Keys are strings with a maximum length of 64 characters. Values are strings with
	// a maximum length of 512 characters.
	Metadata shared.Metadata `json:"metadata,required"`
	// The model that is evaluated, if applicable.
	Model string `json:"model,required"`
	// The name of the evaluation run.
	Name string `json:"name,required"`
	// The type of the object. Always "eval.run".
	Object constant.EvalRun `json:"object,required"`
	// Usage statistics for each model during the evaluation run.
	PerModelUsage []EvalRunListResponsePerModelUsage `json:"per_model_usage,required"`
	// Results per testing criteria applied during the evaluation run.
	PerTestingCriteriaResults []EvalRunListResponsePerTestingCriteriaResult `json:"per_testing_criteria_results,required"`
	// The URL to the rendered evaluation run report on the UI dashboard.
	ReportURL string `json:"report_url,required"`
	// Counters summarizing the outcomes of the evaluation run.
	ResultCounts EvalRunListResponseResultCounts `json:"result_counts,required"`
	// The status of the evaluation run.
	Status string `json:"status,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		ID                        resp.Field
		CreatedAt                 resp.Field
		DataSource                resp.Field
		Error                     resp.Field
		EvalID                    resp.Field
		Metadata                  resp.Field
		Model                     resp.Field
		Name                      resp.Field
		Object                    resp.Field
		PerModelUsage             resp.Field
		PerTestingCriteriaResults resp.Field
		ReportURL                 resp.Field
		ResultCounts              resp.Field
		Status                    resp.Field
		ExtraFields               map[string]resp.Field
		raw                       string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EvalRunListResponse) RawJSON() string { return r.JSON.raw }
func (r *EvalRunListResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// EvalRunListResponseDataSourceUnion contains all possible properties and values
// from [CreateEvalJSONLRunDataSource], [CreateEvalCompletionsRunDataSource].
//
// Use the [EvalRunListResponseDataSourceUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type EvalRunListResponseDataSourceUnion struct {
	// This field is a union of [CreateEvalJSONLRunDataSourceSourceUnion],
	// [CreateEvalCompletionsRunDataSourceSourceUnion]
	Source EvalRunListResponseDataSourceUnionSource `json:"source"`
	// Any of "jsonl", "completions".
	Type string `json:"type"`
	// This field is from variant [CreateEvalCompletionsRunDataSource].
	InputMessages CreateEvalCompletionsRunDataSourceInputMessagesUnion `json:"input_messages"`
	// This field is from variant [CreateEvalCompletionsRunDataSource].
	Model string `json:"model"`
	// This field is from variant [CreateEvalCompletionsRunDataSource].
	SamplingParams CreateEvalCompletionsRunDataSourceSamplingParams `json:"sampling_params"`
	JSON           struct {
		Source         resp.Field
		Type           resp.Field
		InputMessages  resp.Field
		Model          resp.Field
		SamplingParams resp.Field
		raw            string
	} `json:"-"`
}

// anyEvalRunListResponseDataSource is implemented by each variant of
// [EvalRunListResponseDataSourceUnion] to add type safety for the return type of
// [EvalRunListResponseDataSourceUnion.AsAny]
type anyEvalRunListResponseDataSource interface {
	implEvalRunListResponseDataSourceUnion()
}

func (CreateEvalJSONLRunDataSource) implEvalRunListResponseDataSourceUnion()       {}
func (CreateEvalCompletionsRunDataSource) implEvalRunListResponseDataSourceUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := EvalRunListResponseDataSourceUnion.AsAny().(type) {
//	case CreateEvalJSONLRunDataSource:
//	case CreateEvalCompletionsRunDataSource:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u EvalRunListResponseDataSourceUnion) AsAny() anyEvalRunListResponseDataSource {
	switch u.Type {
	case "jsonl":
		return u.AsJSONL()
	case "completions":
		return u.AsCompletions()
	}
	return nil
}

func (u EvalRunListResponseDataSourceUnion) AsJSONL() (v CreateEvalJSONLRunDataSource) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u EvalRunListResponseDataSourceUnion) AsCompletions() (v CreateEvalCompletionsRunDataSource) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u EvalRunListResponseDataSourceUnion) RawJSON() string { return u.JSON.raw }

func (r *EvalRunListResponseDataSourceUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// EvalRunListResponseDataSourceUnionSource is an implicit subunion of
// [EvalRunListResponseDataSourceUnion]. EvalRunListResponseDataSourceUnionSource
// provides convenient access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [EvalRunListResponseDataSourceUnion].
type EvalRunListResponseDataSourceUnionSource struct {
	// This field is a union of
	// [[]CreateEvalJSONLRunDataSourceSourceFileContentContent],
	// [[]CreateEvalCompletionsRunDataSourceSourceFileContentContent]
	Content EvalRunListResponseDataSourceUnionSourceContent `json:"content"`
	Type    string                                          `json:"type"`
	ID      string                                          `json:"id"`
	// This field is from variant [CreateEvalCompletionsRunDataSourceSourceUnion].
	CreatedAfter int64 `json:"created_after"`
	// This field is from variant [CreateEvalCompletionsRunDataSourceSourceUnion].
	CreatedBefore int64 `json:"created_before"`
	// This field is from variant [CreateEvalCompletionsRunDataSourceSourceUnion].
	Limit int64 `json:"limit"`
	// This field is from variant [CreateEvalCompletionsRunDataSourceSourceUnion].
	Metadata shared.Metadata `json:"metadata"`
	// This field is from variant [CreateEvalCompletionsRunDataSourceSourceUnion].
	Model string `json:"model"`
	JSON  struct {
		Content       resp.Field
		Type          resp.Field
		ID            resp.Field
		CreatedAfter  resp.Field
		CreatedBefore resp.Field
		Limit         resp.Field
		Metadata      resp.Field
		Model         resp.Field
		raw           string
	} `json:"-"`
}

func (r *EvalRunListResponseDataSourceUnionSource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// EvalRunListResponseDataSourceUnionSourceContent is an implicit subunion of
// [EvalRunListResponseDataSourceUnion].
// EvalRunListResponseDataSourceUnionSourceContent provides convenient access to
// the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [EvalRunListResponseDataSourceUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfCreateEvalJSONLRunDataSourceSourceFileContentContent
// OfCreateEvalCompletionsRunDataSourceSourceFileContentContent]
type EvalRunListResponseDataSourceUnionSourceContent struct {
	// This field will be present if the value is a
	// [[]CreateEvalJSONLRunDataSourceSourceFileContentContent] instead of an object.
	OfCreateEvalJSONLRunDataSourceSourceFileContentContent []CreateEvalJSONLRunDataSourceSourceFileContentContent `json:",inline"`
	// This field will be present if the value is a
	// [[]CreateEvalCompletionsRunDataSourceSourceFileContentContent] instead of an
	// object.
	OfCreateEvalCompletionsRunDataSourceSourceFileContentContent []CreateEvalCompletionsRunDataSourceSourceFileContentContent `json:",inline"`
	JSON                                                         struct {
		OfCreateEvalJSONLRunDataSourceSourceFileContentContent       resp.Field
		OfCreateEvalCompletionsRunDataSourceSourceFileContentContent resp.Field
		raw                                                          string
	} `json:"-"`
}

func (r *EvalRunListResponseDataSourceUnionSourceContent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type EvalRunListResponsePerModelUsage struct {
	// The number of tokens retrieved from cache.
	CachedTokens int64 `json:"cached_tokens,required"`
	// The number of completion tokens generated.
	CompletionTokens int64 `json:"completion_tokens,required"`
	// The number of invocations.
	InvocationCount int64 `json:"invocation_count,required"`
	// The name of the model.
	ModelName string `json:"model_name,required"`
	// The number of prompt tokens used.
	PromptTokens int64 `json:"prompt_tokens,required"`
	// The total number of tokens used.
	TotalTokens int64 `json:"total_tokens,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		CachedTokens     resp.Field
		CompletionTokens resp.Field
		InvocationCount  resp.Field
		ModelName        resp.Field
		PromptTokens     resp.Field
		TotalTokens      resp.Field
		ExtraFields      map[string]resp.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EvalRunListResponsePerModelUsage) RawJSON() string { return r.JSON.raw }
func (r *EvalRunListResponsePerModelUsage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type EvalRunListResponsePerTestingCriteriaResult struct {
	// Number of tests failed for this criteria.
	Failed int64 `json:"failed,required"`
	// Number of tests passed for this criteria.
	Passed int64 `json:"passed,required"`
	// A description of the testing criteria.
	TestingCriteria string `json:"testing_criteria,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		Failed          resp.Field
		Passed          resp.Field
		TestingCriteria resp.Field
		ExtraFields     map[string]resp.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EvalRunListResponsePerTestingCriteriaResult) RawJSON() string { return r.JSON.raw }
func (r *EvalRunListResponsePerTestingCriteriaResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Counters summarizing the outcomes of the evaluation run.
type EvalRunListResponseResultCounts struct {
	// Number of output items that resulted in an error.
	Errored int64 `json:"errored,required"`
	// Number of output items that failed to pass the evaluation.
	Failed int64 `json:"failed,required"`
	// Number of output items that passed the evaluation.
	Passed int64 `json:"passed,required"`
	// Total number of executed output items.
	Total int64 `json:"total,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		Errored     resp.Field
		Failed      resp.Field
		Passed      resp.Field
		Total       resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EvalRunListResponseResultCounts) RawJSON() string { return r.JSON.raw }
func (r *EvalRunListResponseResultCounts) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type EvalRunDeleteResponse struct {
	Deleted bool   `json:"deleted"`
	Object  string `json:"object"`
	RunID   string `json:"run_id"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		Deleted     resp.Field
		Object      resp.Field
		RunID       resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EvalRunDeleteResponse) RawJSON() string { return r.JSON.raw }
func (r *EvalRunDeleteResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A schema representing an evaluation run.
type EvalRunCancelResponse struct {
	// Unique identifier for the evaluation run.
	ID string `json:"id,required"`
	// Unix timestamp (in seconds) when the evaluation run was created.
	CreatedAt int64 `json:"created_at,required"`
	// Information about the run's data source.
	DataSource EvalRunCancelResponseDataSourceUnion `json:"data_source,required"`
	// An object representing an error response from the Eval API.
	Error EvalAPIError `json:"error,required"`
	// The identifier of the associated evaluation.
	EvalID string `json:"eval_id,required"`
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard.
	//
	// Keys are strings with a maximum length of 64 characters. Values are strings with
	// a maximum length of 512 characters.
	Metadata shared.Metadata `json:"metadata,required"`
	// The model that is evaluated, if applicable.
	Model string `json:"model,required"`
	// The name of the evaluation run.
	Name string `json:"name,required"`
	// The type of the object. Always "eval.run".
	Object constant.EvalRun `json:"object,required"`
	// Usage statistics for each model during the evaluation run.
	PerModelUsage []EvalRunCancelResponsePerModelUsage `json:"per_model_usage,required"`
	// Results per testing criteria applied during the evaluation run.
	PerTestingCriteriaResults []EvalRunCancelResponsePerTestingCriteriaResult `json:"per_testing_criteria_results,required"`
	// The URL to the rendered evaluation run report on the UI dashboard.
	ReportURL string `json:"report_url,required"`
	// Counters summarizing the outcomes of the evaluation run.
	ResultCounts EvalRunCancelResponseResultCounts `json:"result_counts,required"`
	// The status of the evaluation run.
	Status string `json:"status,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		ID                        resp.Field
		CreatedAt                 resp.Field
		DataSource                resp.Field
		Error                     resp.Field
		EvalID                    resp.Field
		Metadata                  resp.Field
		Model                     resp.Field
		Name                      resp.Field
		Object                    resp.Field
		PerModelUsage             resp.Field
		PerTestingCriteriaResults resp.Field
		ReportURL                 resp.Field
		ResultCounts              resp.Field
		Status                    resp.Field
		ExtraFields               map[string]resp.Field
		raw                       string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EvalRunCancelResponse) RawJSON() string { return r.JSON.raw }
func (r *EvalRunCancelResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// EvalRunCancelResponseDataSourceUnion contains all possible properties and values
// from [CreateEvalJSONLRunDataSource], [CreateEvalCompletionsRunDataSource].
//
// Use the [EvalRunCancelResponseDataSourceUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type EvalRunCancelResponseDataSourceUnion struct {
	// This field is a union of [CreateEvalJSONLRunDataSourceSourceUnion],
	// [CreateEvalCompletionsRunDataSourceSourceUnion]
	Source EvalRunCancelResponseDataSourceUnionSource `json:"source"`
	// Any of "jsonl", "completions".
	Type string `json:"type"`
	// This field is from variant [CreateEvalCompletionsRunDataSource].
	InputMessages CreateEvalCompletionsRunDataSourceInputMessagesUnion `json:"input_messages"`
	// This field is from variant [CreateEvalCompletionsRunDataSource].
	Model string `json:"model"`
	// This field is from variant [CreateEvalCompletionsRunDataSource].
	SamplingParams CreateEvalCompletionsRunDataSourceSamplingParams `json:"sampling_params"`
	JSON           struct {
		Source         resp.Field
		Type           resp.Field
		InputMessages  resp.Field
		Model          resp.Field
		SamplingParams resp.Field
		raw            string
	} `json:"-"`
}

// anyEvalRunCancelResponseDataSource is implemented by each variant of
// [EvalRunCancelResponseDataSourceUnion] to add type safety for the return type of
// [EvalRunCancelResponseDataSourceUnion.AsAny]
type anyEvalRunCancelResponseDataSource interface {
	implEvalRunCancelResponseDataSourceUnion()
}

func (CreateEvalJSONLRunDataSource) implEvalRunCancelResponseDataSourceUnion()       {}
func (CreateEvalCompletionsRunDataSource) implEvalRunCancelResponseDataSourceUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := EvalRunCancelResponseDataSourceUnion.AsAny().(type) {
//	case CreateEvalJSONLRunDataSource:
//	case CreateEvalCompletionsRunDataSource:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u EvalRunCancelResponseDataSourceUnion) AsAny() anyEvalRunCancelResponseDataSource {
	switch u.Type {
	case "jsonl":
		return u.AsJSONL()
	case "completions":
		return u.AsCompletions()
	}
	return nil
}

func (u EvalRunCancelResponseDataSourceUnion) AsJSONL() (v CreateEvalJSONLRunDataSource) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u EvalRunCancelResponseDataSourceUnion) AsCompletions() (v CreateEvalCompletionsRunDataSource) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u EvalRunCancelResponseDataSourceUnion) RawJSON() string { return u.JSON.raw }

func (r *EvalRunCancelResponseDataSourceUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// EvalRunCancelResponseDataSourceUnionSource is an implicit subunion of
// [EvalRunCancelResponseDataSourceUnion].
// EvalRunCancelResponseDataSourceUnionSource provides convenient access to the
// sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [EvalRunCancelResponseDataSourceUnion].
type EvalRunCancelResponseDataSourceUnionSource struct {
	// This field is a union of
	// [[]CreateEvalJSONLRunDataSourceSourceFileContentContent],
	// [[]CreateEvalCompletionsRunDataSourceSourceFileContentContent]
	Content EvalRunCancelResponseDataSourceUnionSourceContent `json:"content"`
	Type    string                                            `json:"type"`
	ID      string                                            `json:"id"`
	// This field is from variant [CreateEvalCompletionsRunDataSourceSourceUnion].
	CreatedAfter int64 `json:"created_after"`
	// This field is from variant [CreateEvalCompletionsRunDataSourceSourceUnion].
	CreatedBefore int64 `json:"created_before"`
	// This field is from variant [CreateEvalCompletionsRunDataSourceSourceUnion].
	Limit int64 `json:"limit"`
	// This field is from variant [CreateEvalCompletionsRunDataSourceSourceUnion].
	Metadata shared.Metadata `json:"metadata"`
	// This field is from variant [CreateEvalCompletionsRunDataSourceSourceUnion].
	Model string `json:"model"`
	JSON  struct {
		Content       resp.Field
		Type          resp.Field
		ID            resp.Field
		CreatedAfter  resp.Field
		CreatedBefore resp.Field
		Limit         resp.Field
		Metadata      resp.Field
		Model         resp.Field
		raw           string
	} `json:"-"`
}

func (r *EvalRunCancelResponseDataSourceUnionSource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// EvalRunCancelResponseDataSourceUnionSourceContent is an implicit subunion of
// [EvalRunCancelResponseDataSourceUnion].
// EvalRunCancelResponseDataSourceUnionSourceContent provides convenient access to
// the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [EvalRunCancelResponseDataSourceUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfCreateEvalJSONLRunDataSourceSourceFileContentContent
// OfCreateEvalCompletionsRunDataSourceSourceFileContentContent]
type EvalRunCancelResponseDataSourceUnionSourceContent struct {
	// This field will be present if the value is a
	// [[]CreateEvalJSONLRunDataSourceSourceFileContentContent] instead of an object.
	OfCreateEvalJSONLRunDataSourceSourceFileContentContent []CreateEvalJSONLRunDataSourceSourceFileContentContent `json:",inline"`
	// This field will be present if the value is a
	// [[]CreateEvalCompletionsRunDataSourceSourceFileContentContent] instead of an
	// object.
	OfCreateEvalCompletionsRunDataSourceSourceFileContentContent []CreateEvalCompletionsRunDataSourceSourceFileContentContent `json:",inline"`
	JSON                                                         struct {
		OfCreateEvalJSONLRunDataSourceSourceFileContentContent       resp.Field
		OfCreateEvalCompletionsRunDataSourceSourceFileContentContent resp.Field
		raw                                                          string
	} `json:"-"`
}

func (r *EvalRunCancelResponseDataSourceUnionSourceContent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type EvalRunCancelResponsePerModelUsage struct {
	// The number of tokens retrieved from cache.
	CachedTokens int64 `json:"cached_tokens,required"`
	// The number of completion tokens generated.
	CompletionTokens int64 `json:"completion_tokens,required"`
	// The number of invocations.
	InvocationCount int64 `json:"invocation_count,required"`
	// The name of the model.
	ModelName string `json:"model_name,required"`
	// The number of prompt tokens used.
	PromptTokens int64 `json:"prompt_tokens,required"`
	// The total number of tokens used.
	TotalTokens int64 `json:"total_tokens,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		CachedTokens     resp.Field
		CompletionTokens resp.Field
		InvocationCount  resp.Field
		ModelName        resp.Field
		PromptTokens     resp.Field
		TotalTokens      resp.Field
		ExtraFields      map[string]resp.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EvalRunCancelResponsePerModelUsage) RawJSON() string { return r.JSON.raw }
func (r *EvalRunCancelResponsePerModelUsage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type EvalRunCancelResponsePerTestingCriteriaResult struct {
	// Number of tests failed for this criteria.
	Failed int64 `json:"failed,required"`
	// Number of tests passed for this criteria.
	Passed int64 `json:"passed,required"`
	// A description of the testing criteria.
	TestingCriteria string `json:"testing_criteria,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		Failed          resp.Field
		Passed          resp.Field
		TestingCriteria resp.Field
		ExtraFields     map[string]resp.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EvalRunCancelResponsePerTestingCriteriaResult) RawJSON() string { return r.JSON.raw }
func (r *EvalRunCancelResponsePerTestingCriteriaResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Counters summarizing the outcomes of the evaluation run.
type EvalRunCancelResponseResultCounts struct {
	// Number of output items that resulted in an error.
	Errored int64 `json:"errored,required"`
	// Number of output items that failed to pass the evaluation.
	Failed int64 `json:"failed,required"`
	// Number of output items that passed the evaluation.
	Passed int64 `json:"passed,required"`
	// Total number of executed output items.
	Total int64 `json:"total,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		Errored     resp.Field
		Failed      resp.Field
		Passed      resp.Field
		Total       resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EvalRunCancelResponseResultCounts) RawJSON() string { return r.JSON.raw }
func (r *EvalRunCancelResponseResultCounts) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type EvalRunNewParams struct {
	// Details about the run's data source.
	DataSource EvalRunNewParamsDataSourceUnion `json:"data_source,omitzero,required"`
	// The name of the run.
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
func (f EvalRunNewParams) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }

func (r EvalRunNewParams) MarshalJSON() (data []byte, err error) {
	type shadow EvalRunNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type EvalRunNewParamsDataSourceUnion struct {
	OfJSONLRunDataSource       *CreateEvalJSONLRunDataSourceParam       `json:",omitzero,inline"`
	OfCompletionsRunDataSource *CreateEvalCompletionsRunDataSourceParam `json:",omitzero,inline"`
	paramUnion
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (u EvalRunNewParamsDataSourceUnion) IsPresent() bool { return !param.IsOmitted(u) && !u.IsNull() }
func (u EvalRunNewParamsDataSourceUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion[EvalRunNewParamsDataSourceUnion](u.OfJSONLRunDataSource, u.OfCompletionsRunDataSource)
}

func (u *EvalRunNewParamsDataSourceUnion) asAny() any {
	if !param.IsOmitted(u.OfJSONLRunDataSource) {
		return u.OfJSONLRunDataSource
	} else if !param.IsOmitted(u.OfCompletionsRunDataSource) {
		return u.OfCompletionsRunDataSource
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u EvalRunNewParamsDataSourceUnion) GetInputMessages() *CreateEvalCompletionsRunDataSourceInputMessagesUnionParam {
	if vt := u.OfCompletionsRunDataSource; vt != nil {
		return &vt.InputMessages
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u EvalRunNewParamsDataSourceUnion) GetModel() *string {
	if vt := u.OfCompletionsRunDataSource; vt != nil {
		return &vt.Model
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u EvalRunNewParamsDataSourceUnion) GetSamplingParams() *CreateEvalCompletionsRunDataSourceSamplingParamsParam {
	if vt := u.OfCompletionsRunDataSource; vt != nil {
		return &vt.SamplingParams
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u EvalRunNewParamsDataSourceUnion) GetType() *string {
	if vt := u.OfJSONLRunDataSource; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfCompletionsRunDataSource; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u EvalRunNewParamsDataSourceUnion) GetSource() (res evalRunNewParamsDataSourceUnionSource) {
	if vt := u.OfJSONLRunDataSource; vt != nil {
		res.ofCreateEvalJSONLRunDataSourceSourceUnion = &vt.Source
	} else if vt := u.OfCompletionsRunDataSource; vt != nil {
		res.ofCreateEvalCompletionsRunDataSourceSourceUnion = &vt.Source
	}
	return
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type evalRunNewParamsDataSourceUnionSource struct {
	ofCreateEvalJSONLRunDataSourceSourceUnion       *CreateEvalJSONLRunDataSourceSourceUnionParam
	ofCreateEvalCompletionsRunDataSourceSourceUnion *CreateEvalCompletionsRunDataSourceSourceUnionParam
}

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *openai.CreateEvalJSONLRunDataSourceSourceFileContentParam:
//	case *openai.CreateEvalJSONLRunDataSourceSourceFileIDParam:
//	case *openai.CreateEvalCompletionsRunDataSourceSourceFileContentParam:
//	case *openai.CreateEvalCompletionsRunDataSourceSourceFileIDParam:
//	case *openai.CreateEvalCompletionsRunDataSourceSourceStoredCompletionsParam:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u evalRunNewParamsDataSourceUnionSource) AsAny() any {
	if !param.IsOmitted(u.ofCreateEvalJSONLRunDataSourceSourceUnion) {
		return u.ofCreateEvalJSONLRunDataSourceSourceUnion.asAny()
	} else if !param.IsOmitted(u.ofCreateEvalCompletionsRunDataSourceSourceUnion) {
		return u.ofCreateEvalCompletionsRunDataSourceSourceUnion.asAny()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u evalRunNewParamsDataSourceUnionSource) GetCreatedAfter() *int64 {
	if u.ofCreateEvalCompletionsRunDataSourceSourceUnion != nil {
		return u.ofCreateEvalCompletionsRunDataSourceSourceUnion.GetCreatedAfter()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u evalRunNewParamsDataSourceUnionSource) GetCreatedBefore() *int64 {
	if u.ofCreateEvalCompletionsRunDataSourceSourceUnion != nil {
		return u.ofCreateEvalCompletionsRunDataSourceSourceUnion.GetCreatedBefore()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u evalRunNewParamsDataSourceUnionSource) GetLimit() *int64 {
	if u.ofCreateEvalCompletionsRunDataSourceSourceUnion != nil {
		return u.ofCreateEvalCompletionsRunDataSourceSourceUnion.GetLimit()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u evalRunNewParamsDataSourceUnionSource) GetMetadata() shared.MetadataParam {
	if u.ofCreateEvalCompletionsRunDataSourceSourceUnion != nil {
		return u.ofCreateEvalCompletionsRunDataSourceSourceUnion.GetMetadata()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u evalRunNewParamsDataSourceUnionSource) GetModel() *string {
	if u.ofCreateEvalCompletionsRunDataSourceSourceUnion != nil {
		return u.ofCreateEvalCompletionsRunDataSourceSourceUnion.GetModel()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u evalRunNewParamsDataSourceUnionSource) GetType() *string {
	if u.ofCreateEvalJSONLRunDataSourceSourceUnion != nil {
		return u.ofCreateEvalJSONLRunDataSourceSourceUnion.GetType()
	} else if u.ofCreateEvalJSONLRunDataSourceSourceUnion != nil {
		return u.ofCreateEvalJSONLRunDataSourceSourceUnion.GetType()
	} else if u.ofCreateEvalCompletionsRunDataSourceSourceUnion != nil {
		return u.ofCreateEvalCompletionsRunDataSourceSourceUnion.GetType()
	} else if u.ofCreateEvalCompletionsRunDataSourceSourceUnion != nil {
		return u.ofCreateEvalCompletionsRunDataSourceSourceUnion.GetType()
	} else if u.ofCreateEvalCompletionsRunDataSourceSourceUnion != nil {
		return u.ofCreateEvalCompletionsRunDataSourceSourceUnion.GetType()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u evalRunNewParamsDataSourceUnionSource) GetID() *string {
	if u.ofCreateEvalJSONLRunDataSourceSourceUnion != nil {
		return u.ofCreateEvalJSONLRunDataSourceSourceUnion.GetID()
	} else if u.ofCreateEvalCompletionsRunDataSourceSourceUnion != nil {
		return u.ofCreateEvalCompletionsRunDataSourceSourceUnion.GetID()
	}
	return nil
}

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u evalRunNewParamsDataSourceUnionSource) GetContent() (res evalRunNewParamsDataSourceUnionSourceContent) {
	if u.ofCreateEvalJSONLRunDataSourceSourceUnion != nil {
		res.ofCreateEvalJSONLRunDataSourceSourceFileContentContent = u.ofCreateEvalJSONLRunDataSourceSourceUnion.GetContent()
	} else if u.ofCreateEvalCompletionsRunDataSourceSourceUnion != nil {
		res.ofCreateEvalCompletionsRunDataSourceSourceFileContentContent = u.ofCreateEvalCompletionsRunDataSourceSourceUnion.GetContent()
	}
	return
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type evalRunNewParamsDataSourceUnionSourceContent struct {
	ofCreateEvalJSONLRunDataSourceSourceFileContentContent       *[]CreateEvalJSONLRunDataSourceSourceFileContentContentParam
	ofCreateEvalCompletionsRunDataSourceSourceFileContentContent *[]CreateEvalCompletionsRunDataSourceSourceFileContentContentParam
}

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *[]openai.CreateEvalJSONLRunDataSourceSourceFileContentContentParam:
//	case *[]openai.CreateEvalCompletionsRunDataSourceSourceFileContentContentParam:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u evalRunNewParamsDataSourceUnionSourceContent) AsAny() any {
	if !param.IsOmitted(u.ofCreateEvalJSONLRunDataSourceSourceFileContentContent) {
		return u.ofCreateEvalJSONLRunDataSourceSourceFileContentContent
	} else if !param.IsOmitted(u.ofCreateEvalCompletionsRunDataSourceSourceFileContentContent) {
		return u.ofCreateEvalCompletionsRunDataSourceSourceFileContentContent
	}
	return nil
}

type EvalRunListParams struct {
	// Identifier for the last run from the previous pagination request.
	After param.Opt[string] `query:"after,omitzero" json:"-"`
	// Number of runs to retrieve.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Sort order for runs by timestamp. Use `asc` for ascending order or `desc` for
	// descending order. Defaults to `asc`.
	//
	// Any of "asc", "desc".
	Order EvalRunListParamsOrder `query:"order,omitzero" json:"-"`
	// Filter runs by status. Use "queued" | "in_progress" | "failed" | "completed" |
	// "canceled".
	//
	// Any of "queued", "in_progress", "completed", "canceled", "failed".
	Status EvalRunListParamsStatus `query:"status,omitzero" json:"-"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f EvalRunListParams) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }

// URLQuery serializes [EvalRunListParams]'s query parameters as `url.Values`.
func (r EvalRunListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Sort order for runs by timestamp. Use `asc` for ascending order or `desc` for
// descending order. Defaults to `asc`.
type EvalRunListParamsOrder string

const (
	EvalRunListParamsOrderAsc  EvalRunListParamsOrder = "asc"
	EvalRunListParamsOrderDesc EvalRunListParamsOrder = "desc"
)

// Filter runs by status. Use "queued" | "in_progress" | "failed" | "completed" |
// "canceled".
type EvalRunListParamsStatus string

const (
	EvalRunListParamsStatusQueued     EvalRunListParamsStatus = "queued"
	EvalRunListParamsStatusInProgress EvalRunListParamsStatus = "in_progress"
	EvalRunListParamsStatusCompleted  EvalRunListParamsStatus = "completed"
	EvalRunListParamsStatusCanceled   EvalRunListParamsStatus = "canceled"
	EvalRunListParamsStatusFailed     EvalRunListParamsStatus = "failed"
)

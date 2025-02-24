// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai

import (
	"context"
	"encoding/json"
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
	"github.com/openai/openai-go/shared"
	"github.com/openai/openai-go/shared/constant"
)

// BetaAssistantService contains methods and other services that help with
// interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaAssistantService] method instead.
type BetaAssistantService struct {
	Options []option.RequestOption
}

// NewBetaAssistantService generates a new service that applies the given options
// to each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewBetaAssistantService(opts ...option.RequestOption) (r BetaAssistantService) {
	r = BetaAssistantService{}
	r.Options = opts
	return
}

// Create an assistant with a model and instructions.
func (r *BetaAssistantService) New(ctx context.Context, body BetaAssistantNewParams, opts ...option.RequestOption) (res *Assistant, err error) {
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "assistants=v2")}, opts...)
	path := "assistants"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Retrieves an assistant.
func (r *BetaAssistantService) Get(ctx context.Context, assistantID string, opts ...option.RequestOption) (res *Assistant, err error) {
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "assistants=v2")}, opts...)
	if assistantID == "" {
		err = errors.New("missing required assistant_id parameter")
		return
	}
	path := fmt.Sprintf("assistants/%s", assistantID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// Modifies an assistant.
func (r *BetaAssistantService) Update(ctx context.Context, assistantID string, body BetaAssistantUpdateParams, opts ...option.RequestOption) (res *Assistant, err error) {
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "assistants=v2")}, opts...)
	if assistantID == "" {
		err = errors.New("missing required assistant_id parameter")
		return
	}
	path := fmt.Sprintf("assistants/%s", assistantID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Returns a list of assistants.
func (r *BetaAssistantService) List(ctx context.Context, query BetaAssistantListParams, opts ...option.RequestOption) (res *pagination.CursorPage[Assistant], err error) {
	var raw *http.Response
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "assistants=v2"), option.WithResponseInto(&raw)}, opts...)
	path := "assistants"
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

// Returns a list of assistants.
func (r *BetaAssistantService) ListAutoPaging(ctx context.Context, query BetaAssistantListParams, opts ...option.RequestOption) *pagination.CursorPageAutoPager[Assistant] {
	return pagination.NewCursorPageAutoPager(r.List(ctx, query, opts...))
}

// Delete an assistant.
func (r *BetaAssistantService) Delete(ctx context.Context, assistantID string, opts ...option.RequestOption) (res *AssistantDeleted, err error) {
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "assistants=v2")}, opts...)
	if assistantID == "" {
		err = errors.New("missing required assistant_id parameter")
		return
	}
	path := fmt.Sprintf("assistants/%s", assistantID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return
}

// Represents an `assistant` that can call the model and use tools.
type Assistant struct {
	// The identifier, which can be referenced in API endpoints.
	ID string `json:"id,omitzero,required"`
	// The Unix timestamp (in seconds) for when the assistant was created.
	CreatedAt int64 `json:"created_at,omitzero,required"`
	// The description of the assistant. The maximum length is 512 characters.
	Description string `json:"description,omitzero,required,nullable"`
	// The system instructions that the assistant uses. The maximum length is 256,000
	// characters.
	Instructions string `json:"instructions,omitzero,required,nullable"`
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard.
	//
	// Keys are strings with a maximum length of 64 characters. Values are strings with
	// a maximum length of 512 characters.
	Metadata shared.Metadata `json:"metadata,omitzero,required,nullable"`
	// ID of the model to use. You can use the
	// [List models](https://platform.openai.com/docs/api-reference/models/list) API to
	// see all of your available models, or see our
	// [Model overview](https://platform.openai.com/docs/models) for descriptions of
	// them.
	Model string `json:"model,omitzero,required"`
	// The name of the assistant. The maximum length is 256 characters.
	Name string `json:"name,omitzero,required,nullable"`
	// The object type, which is always `assistant`.
	//
	// This field can be elided, and will be automatically set as "assistant".
	Object constant.Assistant `json:"object,required"`
	// A list of tool enabled on the assistant. There can be a maximum of 128 tools per
	// assistant. Tools can be of types `code_interpreter`, `file_search`, or
	// `function`.
	Tools []AssistantToolUnion `json:"tools,omitzero,required"`
	// What sampling temperature to use, between 0 and 2. Higher values like 0.8 will
	// make the output more random, while lower values like 0.2 will make it more
	// focused and deterministic.
	Temperature float64 `json:"temperature,omitzero,nullable"`
	// A set of resources that are used by the assistant's tools. The resources are
	// specific to the type of tool. For example, the `code_interpreter` tool requires
	// a list of file IDs, while the `file_search` tool requires a list of vector store
	// IDs.
	ToolResources AssistantToolResources `json:"tool_resources,omitzero,nullable"`
	// An alternative to sampling with temperature, called nucleus sampling, where the
	// model considers the results of the tokens with top_p probability mass. So 0.1
	// means only the tokens comprising the top 10% probability mass are considered.
	//
	// We generally recommend altering this or temperature but not both.
	TopP float64 `json:"top_p,omitzero,nullable"`
	JSON struct {
		ID            resp.Field
		CreatedAt     resp.Field
		Description   resp.Field
		Instructions  resp.Field
		Metadata      resp.Field
		Model         resp.Field
		Name          resp.Field
		Object        resp.Field
		Tools         resp.Field
		Temperature   resp.Field
		ToolResources resp.Field
		TopP          resp.Field
		raw           string
	} `json:"-"`
}

func (r Assistant) RawJSON() string { return r.JSON.raw }
func (r *Assistant) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A set of resources that are used by the assistant's tools. The resources are
// specific to the type of tool. For example, the `code_interpreter` tool requires
// a list of file IDs, while the `file_search` tool requires a list of vector store
// IDs.
type AssistantToolResources struct {
	CodeInterpreter AssistantToolResourcesCodeInterpreter `json:"code_interpreter,omitzero"`
	FileSearch      AssistantToolResourcesFileSearch      `json:"file_search,omitzero"`
	JSON            struct {
		CodeInterpreter resp.Field
		FileSearch      resp.Field
		raw             string
	} `json:"-"`
}

func (r AssistantToolResources) RawJSON() string { return r.JSON.raw }
func (r *AssistantToolResources) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AssistantToolResourcesCodeInterpreter struct {
	// A list of [file](https://platform.openai.com/docs/api-reference/files) IDs made
	// available to the `code_interpreter“ tool. There can be a maximum of 20 files
	// associated with the tool.
	FileIDs []string `json:"file_ids,omitzero"`
	JSON    struct {
		FileIDs resp.Field
		raw     string
	} `json:"-"`
}

func (r AssistantToolResourcesCodeInterpreter) RawJSON() string { return r.JSON.raw }
func (r *AssistantToolResourcesCodeInterpreter) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AssistantToolResourcesFileSearch struct {
	// The ID of the
	// [vector store](https://platform.openai.com/docs/api-reference/vector-stores/object)
	// attached to this assistant. There can be a maximum of 1 vector store attached to
	// the assistant.
	VectorStoreIDs []string `json:"vector_store_ids,omitzero"`
	JSON           struct {
		VectorStoreIDs resp.Field
		raw            string
	} `json:"-"`
}

func (r AssistantToolResourcesFileSearch) RawJSON() string { return r.JSON.raw }
func (r *AssistantToolResourcesFileSearch) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AssistantDeleted struct {
	ID      string `json:"id,omitzero,required"`
	Deleted bool   `json:"deleted,omitzero,required"`
	// This field can be elided, and will be automatically set as "assistant.deleted".
	Object constant.AssistantDeleted `json:"object,required"`
	JSON   struct {
		ID      resp.Field
		Deleted resp.Field
		Object  resp.Field
		raw     string
	} `json:"-"`
}

func (r AssistantDeleted) RawJSON() string { return r.JSON.raw }
func (r *AssistantDeleted) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AssistantStreamEventUnion struct {
	// This field is a union of
	// [Thread,Run,RunStep,RunStepDeltaEvent,Message,MessageDeltaEvent,shared.ErrorObject]
	Data    AssistantStreamEventUnionData `json:"data"`
	Event   string                        `json:"event"`
	Enabled bool                          `json:"enabled"`
	JSON    struct {
		Data    resp.Field
		Event   resp.Field
		Enabled resp.Field
		raw     string
	} `json:"-"`
}

// note: this function is generated only for discriminated unions
func (u AssistantStreamEventUnion) Variant() (res struct {
	OfThreadCreated           *AssistantStreamEventThreadCreated
	OfThreadRunCreated        *AssistantStreamEventThreadRunCreated
	OfThreadRunQueued         *AssistantStreamEventThreadRunQueued
	OfThreadRunInProgress     *AssistantStreamEventThreadRunInProgress
	OfThreadRunRequiresAction *AssistantStreamEventThreadRunRequiresAction
	OfThreadRunCompleted      *AssistantStreamEventThreadRunCompleted
	OfThreadRunIncomplete     *AssistantStreamEventThreadRunIncomplete
	OfThreadRunFailed         *AssistantStreamEventThreadRunFailed
	OfThreadRunCancelling     *AssistantStreamEventThreadRunCancelling
	OfThreadRunCancelled      *AssistantStreamEventThreadRunCancelled
	OfThreadRunExpired        *AssistantStreamEventThreadRunExpired
	OfThreadRunStepCreated    *AssistantStreamEventThreadRunStepCreated
	OfThreadRunStepInProgress *AssistantStreamEventThreadRunStepInProgress
	OfThreadRunStepDelta      *AssistantStreamEventThreadRunStepDelta
	OfThreadRunStepCompleted  *AssistantStreamEventThreadRunStepCompleted
	OfThreadRunStepFailed     *AssistantStreamEventThreadRunStepFailed
	OfThreadRunStepCancelled  *AssistantStreamEventThreadRunStepCancelled
	OfThreadRunStepExpired    *AssistantStreamEventThreadRunStepExpired
	OfThreadMessageCreated    *AssistantStreamEventThreadMessageCreated
	OfThreadMessageInProgress *AssistantStreamEventThreadMessageInProgress
	OfThreadMessageDelta      *AssistantStreamEventThreadMessageDelta
	OfThreadMessageCompleted  *AssistantStreamEventThreadMessageCompleted
	OfThreadMessageIncomplete *AssistantStreamEventThreadMessageIncomplete
	OfErrorEvent              *AssistantStreamEventErrorEvent
}) {
	switch u.Event {
	case "thread.created":
		v := u.AsThreadCreated()
		res.OfThreadCreated = &v
	case "thread.run.created":
		v := u.AsThreadRunCreated()
		res.OfThreadRunCreated = &v
	case "thread.run.queued":
		v := u.AsThreadRunQueued()
		res.OfThreadRunQueued = &v
	case "thread.run.in_progress":
		v := u.AsThreadRunInProgress()
		res.OfThreadRunInProgress = &v
	case "thread.run.requires_action":
		v := u.AsThreadRunRequiresAction()
		res.OfThreadRunRequiresAction = &v
	case "thread.run.completed":
		v := u.AsThreadRunCompleted()
		res.OfThreadRunCompleted = &v
	case "thread.run.incomplete":
		v := u.AsThreadRunIncomplete()
		res.OfThreadRunIncomplete = &v
	case "thread.run.failed":
		v := u.AsThreadRunFailed()
		res.OfThreadRunFailed = &v
	case "thread.run.cancelling":
		v := u.AsThreadRunCancelling()
		res.OfThreadRunCancelling = &v
	case "thread.run.cancelled":
		v := u.AsThreadRunCancelled()
		res.OfThreadRunCancelled = &v
	case "thread.run.expired":
		v := u.AsThreadRunExpired()
		res.OfThreadRunExpired = &v
	case "thread.run.step.created":
		v := u.AsThreadRunStepCreated()
		res.OfThreadRunStepCreated = &v
	case "thread.run.step.in_progress":
		v := u.AsThreadRunStepInProgress()
		res.OfThreadRunStepInProgress = &v
	case "thread.run.step.delta":
		v := u.AsThreadRunStepDelta()
		res.OfThreadRunStepDelta = &v
	case "thread.run.step.completed":
		v := u.AsThreadRunStepCompleted()
		res.OfThreadRunStepCompleted = &v
	case "thread.run.step.failed":
		v := u.AsThreadRunStepFailed()
		res.OfThreadRunStepFailed = &v
	case "thread.run.step.cancelled":
		v := u.AsThreadRunStepCancelled()
		res.OfThreadRunStepCancelled = &v
	case "thread.run.step.expired":
		v := u.AsThreadRunStepExpired()
		res.OfThreadRunStepExpired = &v
	case "thread.message.created":
		v := u.AsThreadMessageCreated()
		res.OfThreadMessageCreated = &v
	case "thread.message.in_progress":
		v := u.AsThreadMessageInProgress()
		res.OfThreadMessageInProgress = &v
	case "thread.message.delta":
		v := u.AsThreadMessageDelta()
		res.OfThreadMessageDelta = &v
	case "thread.message.completed":
		v := u.AsThreadMessageCompleted()
		res.OfThreadMessageCompleted = &v
	case "thread.message.incomplete":
		v := u.AsThreadMessageIncomplete()
		res.OfThreadMessageIncomplete = &v
	case "error":
		v := u.AsErrorEvent()
		res.OfErrorEvent = &v
	}
	return
}

func (u AssistantStreamEventUnion) WhichKind() string {
	return u.Event
}

func (u AssistantStreamEventUnion) AsThreadCreated() (v AssistantStreamEventThreadCreated) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AssistantStreamEventUnion) AsThreadRunCreated() (v AssistantStreamEventThreadRunCreated) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AssistantStreamEventUnion) AsThreadRunQueued() (v AssistantStreamEventThreadRunQueued) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AssistantStreamEventUnion) AsThreadRunInProgress() (v AssistantStreamEventThreadRunInProgress) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AssistantStreamEventUnion) AsThreadRunRequiresAction() (v AssistantStreamEventThreadRunRequiresAction) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AssistantStreamEventUnion) AsThreadRunCompleted() (v AssistantStreamEventThreadRunCompleted) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AssistantStreamEventUnion) AsThreadRunIncomplete() (v AssistantStreamEventThreadRunIncomplete) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AssistantStreamEventUnion) AsThreadRunFailed() (v AssistantStreamEventThreadRunFailed) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AssistantStreamEventUnion) AsThreadRunCancelling() (v AssistantStreamEventThreadRunCancelling) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AssistantStreamEventUnion) AsThreadRunCancelled() (v AssistantStreamEventThreadRunCancelled) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AssistantStreamEventUnion) AsThreadRunExpired() (v AssistantStreamEventThreadRunExpired) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AssistantStreamEventUnion) AsThreadRunStepCreated() (v AssistantStreamEventThreadRunStepCreated) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AssistantStreamEventUnion) AsThreadRunStepInProgress() (v AssistantStreamEventThreadRunStepInProgress) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AssistantStreamEventUnion) AsThreadRunStepDelta() (v AssistantStreamEventThreadRunStepDelta) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AssistantStreamEventUnion) AsThreadRunStepCompleted() (v AssistantStreamEventThreadRunStepCompleted) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AssistantStreamEventUnion) AsThreadRunStepFailed() (v AssistantStreamEventThreadRunStepFailed) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AssistantStreamEventUnion) AsThreadRunStepCancelled() (v AssistantStreamEventThreadRunStepCancelled) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AssistantStreamEventUnion) AsThreadRunStepExpired() (v AssistantStreamEventThreadRunStepExpired) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AssistantStreamEventUnion) AsThreadMessageCreated() (v AssistantStreamEventThreadMessageCreated) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AssistantStreamEventUnion) AsThreadMessageInProgress() (v AssistantStreamEventThreadMessageInProgress) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AssistantStreamEventUnion) AsThreadMessageDelta() (v AssistantStreamEventThreadMessageDelta) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AssistantStreamEventUnion) AsThreadMessageCompleted() (v AssistantStreamEventThreadMessageCompleted) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AssistantStreamEventUnion) AsThreadMessageIncomplete() (v AssistantStreamEventThreadMessageIncomplete) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AssistantStreamEventUnion) AsErrorEvent() (v AssistantStreamEventErrorEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AssistantStreamEventUnion) RawJSON() string { return u.JSON.raw }

func (r *AssistantStreamEventUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AssistantStreamEventUnionData struct {
	ID            string              `json:"id"`
	CreatedAt     int64               `json:"created_at"`
	Metadata      shared.Metadata     `json:"metadata"`
	Object        string              `json:"object"`
	ToolResources ThreadToolResources `json:"tool_resources"`
	AssistantID   string              `json:"assistant_id"`
	CancelledAt   int64               `json:"cancelled_at"`
	CompletedAt   int64               `json:"completed_at"`
	ExpiresAt     int64               `json:"expires_at"`
	FailedAt      int64               `json:"failed_at"`
	// This field is a union of [RunIncompleteDetails,MessageIncompleteDetails]
	IncompleteDetails AssistantStreamEventUnionDataIncompleteDetails `json:"incomplete_details"`
	Instructions      string                                         `json:"instructions"`
	// This field is a union of [RunLastError,RunStepLastError]
	LastError           AssistantStreamEventUnionDataLastError `json:"last_error"`
	MaxCompletionTokens int64                                  `json:"max_completion_tokens"`
	MaxPromptTokens     int64                                  `json:"max_prompt_tokens"`
	Model               string                                 `json:"model"`
	ParallelToolCalls   bool                                   `json:"parallel_tool_calls"`
	RequiredAction      RunRequiredAction                      `json:"required_action"`
	StartedAt           int64                                  `json:"started_at"`
	Status              string                                 `json:"status"`
	ThreadID            string                                 `json:"thread_id"`
	ToolChoice          AssistantToolChoiceOptionUnion         `json:"tool_choice"`
	Tools               []AssistantToolUnion                   `json:"tools"`
	TruncationStrategy  RunTruncationStrategy                  `json:"truncation_strategy"`
	Usage               RunUsage                               `json:"usage"`
	Temperature         float64                                `json:"temperature"`
	TopP                float64                                `json:"top_p"`
	ExpiredAt           int64                                  `json:"expired_at"`
	RunID               string                                 `json:"run_id"`
	StepDetails         RunStepStepDetailsUnion                `json:"step_details"`
	Type                string                                 `json:"type"`
	// This field is a union of [RunStepDelta,MessageDelta]
	Delta        AssistantStreamEventUnionDataDelta `json:"delta"`
	Attachments  []MessageAttachment                `json:"attachments"`
	Content      []MessageContentUnion              `json:"content"`
	IncompleteAt int64                              `json:"incomplete_at"`
	Role         string                             `json:"role"`
	Code         string                             `json:"code"`
	Message      string                             `json:"message"`
	Param        string                             `json:"param"`
	JSON         struct {
		ID                  resp.Field
		CreatedAt           resp.Field
		Metadata            resp.Field
		Object              resp.Field
		ToolResources       resp.Field
		AssistantID         resp.Field
		CancelledAt         resp.Field
		CompletedAt         resp.Field
		ExpiresAt           resp.Field
		FailedAt            resp.Field
		IncompleteDetails   resp.Field
		Instructions        resp.Field
		LastError           resp.Field
		MaxCompletionTokens resp.Field
		MaxPromptTokens     resp.Field
		Model               resp.Field
		ParallelToolCalls   resp.Field
		RequiredAction      resp.Field
		StartedAt           resp.Field
		Status              resp.Field
		ThreadID            resp.Field
		ToolChoice          resp.Field
		Tools               resp.Field
		TruncationStrategy  resp.Field
		Usage               resp.Field
		Temperature         resp.Field
		TopP                resp.Field
		ExpiredAt           resp.Field
		RunID               resp.Field
		StepDetails         resp.Field
		Type                resp.Field
		Delta               resp.Field
		Attachments         resp.Field
		Content             resp.Field
		IncompleteAt        resp.Field
		Role                resp.Field
		Code                resp.Field
		Message             resp.Field
		Param               resp.Field
		raw                 string
	} `json:"-"`
}

func (r *AssistantStreamEventUnionData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AssistantStreamEventUnionDataIncompleteDetails struct {
	Reason string `json:"reason"`
	JSON   struct {
		Reason resp.Field
		raw    string
	} `json:"-"`
}

func (r *AssistantStreamEventUnionDataIncompleteDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AssistantStreamEventUnionDataLastError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	JSON    struct {
		Code    resp.Field
		Message resp.Field
		raw     string
	} `json:"-"`
}

func (r *AssistantStreamEventUnionDataLastError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AssistantStreamEventUnionDataDelta struct {
	StepDetails RunStepDeltaStepDetailsUnion `json:"step_details"`
	Content     []MessageContentDeltaUnion   `json:"content"`
	Role        string                       `json:"role"`
	JSON        struct {
		StepDetails resp.Field
		Content     resp.Field
		Role        resp.Field
		raw         string
	} `json:"-"`
}

func (r *AssistantStreamEventUnionDataDelta) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Occurs when a new
// [thread](https://platform.openai.com/docs/api-reference/threads/object) is
// created.
type AssistantStreamEventThreadCreated struct {
	// Represents a thread that contains
	// [messages](https://platform.openai.com/docs/api-reference/messages).
	Data Thread `json:"data,omitzero,required"`
	// This field can be elided, and will be automatically set as "thread.created".
	Event constant.ThreadCreated `json:"event,required"`
	// Whether to enable input audio transcription.
	Enabled bool `json:"enabled,omitzero"`
	JSON    struct {
		Data    resp.Field
		Event   resp.Field
		Enabled resp.Field
		raw     string
	} `json:"-"`
}

func (r AssistantStreamEventThreadCreated) RawJSON() string { return r.JSON.raw }
func (r *AssistantStreamEventThreadCreated) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Occurs when a new
// [run](https://platform.openai.com/docs/api-reference/runs/object) is created.
type AssistantStreamEventThreadRunCreated struct {
	// Represents an execution run on a
	// [thread](https://platform.openai.com/docs/api-reference/threads).
	Data Run `json:"data,omitzero,required"`
	// This field can be elided, and will be automatically set as "thread.run.created".
	Event constant.ThreadRunCreated `json:"event,required"`
	JSON  struct {
		Data  resp.Field
		Event resp.Field
		raw   string
	} `json:"-"`
}

func (r AssistantStreamEventThreadRunCreated) RawJSON() string { return r.JSON.raw }
func (r *AssistantStreamEventThreadRunCreated) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Occurs when a [run](https://platform.openai.com/docs/api-reference/runs/object)
// moves to a `queued` status.
type AssistantStreamEventThreadRunQueued struct {
	// Represents an execution run on a
	// [thread](https://platform.openai.com/docs/api-reference/threads).
	Data Run `json:"data,omitzero,required"`
	// This field can be elided, and will be automatically set as "thread.run.queued".
	Event constant.ThreadRunQueued `json:"event,required"`
	JSON  struct {
		Data  resp.Field
		Event resp.Field
		raw   string
	} `json:"-"`
}

func (r AssistantStreamEventThreadRunQueued) RawJSON() string { return r.JSON.raw }
func (r *AssistantStreamEventThreadRunQueued) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Occurs when a [run](https://platform.openai.com/docs/api-reference/runs/object)
// moves to an `in_progress` status.
type AssistantStreamEventThreadRunInProgress struct {
	// Represents an execution run on a
	// [thread](https://platform.openai.com/docs/api-reference/threads).
	Data Run `json:"data,omitzero,required"`
	// This field can be elided, and will be automatically set as
	// "thread.run.in_progress".
	Event constant.ThreadRunInProgress `json:"event,required"`
	JSON  struct {
		Data  resp.Field
		Event resp.Field
		raw   string
	} `json:"-"`
}

func (r AssistantStreamEventThreadRunInProgress) RawJSON() string { return r.JSON.raw }
func (r *AssistantStreamEventThreadRunInProgress) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Occurs when a [run](https://platform.openai.com/docs/api-reference/runs/object)
// moves to a `requires_action` status.
type AssistantStreamEventThreadRunRequiresAction struct {
	// Represents an execution run on a
	// [thread](https://platform.openai.com/docs/api-reference/threads).
	Data Run `json:"data,omitzero,required"`
	// This field can be elided, and will be automatically set as
	// "thread.run.requires_action".
	Event constant.ThreadRunRequiresAction `json:"event,required"`
	JSON  struct {
		Data  resp.Field
		Event resp.Field
		raw   string
	} `json:"-"`
}

func (r AssistantStreamEventThreadRunRequiresAction) RawJSON() string { return r.JSON.raw }
func (r *AssistantStreamEventThreadRunRequiresAction) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Occurs when a [run](https://platform.openai.com/docs/api-reference/runs/object)
// is completed.
type AssistantStreamEventThreadRunCompleted struct {
	// Represents an execution run on a
	// [thread](https://platform.openai.com/docs/api-reference/threads).
	Data Run `json:"data,omitzero,required"`
	// This field can be elided, and will be automatically set as
	// "thread.run.completed".
	Event constant.ThreadRunCompleted `json:"event,required"`
	JSON  struct {
		Data  resp.Field
		Event resp.Field
		raw   string
	} `json:"-"`
}

func (r AssistantStreamEventThreadRunCompleted) RawJSON() string { return r.JSON.raw }
func (r *AssistantStreamEventThreadRunCompleted) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Occurs when a [run](https://platform.openai.com/docs/api-reference/runs/object)
// ends with status `incomplete`.
type AssistantStreamEventThreadRunIncomplete struct {
	// Represents an execution run on a
	// [thread](https://platform.openai.com/docs/api-reference/threads).
	Data Run `json:"data,omitzero,required"`
	// This field can be elided, and will be automatically set as
	// "thread.run.incomplete".
	Event constant.ThreadRunIncomplete `json:"event,required"`
	JSON  struct {
		Data  resp.Field
		Event resp.Field
		raw   string
	} `json:"-"`
}

func (r AssistantStreamEventThreadRunIncomplete) RawJSON() string { return r.JSON.raw }
func (r *AssistantStreamEventThreadRunIncomplete) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Occurs when a [run](https://platform.openai.com/docs/api-reference/runs/object)
// fails.
type AssistantStreamEventThreadRunFailed struct {
	// Represents an execution run on a
	// [thread](https://platform.openai.com/docs/api-reference/threads).
	Data Run `json:"data,omitzero,required"`
	// This field can be elided, and will be automatically set as "thread.run.failed".
	Event constant.ThreadRunFailed `json:"event,required"`
	JSON  struct {
		Data  resp.Field
		Event resp.Field
		raw   string
	} `json:"-"`
}

func (r AssistantStreamEventThreadRunFailed) RawJSON() string { return r.JSON.raw }
func (r *AssistantStreamEventThreadRunFailed) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Occurs when a [run](https://platform.openai.com/docs/api-reference/runs/object)
// moves to a `cancelling` status.
type AssistantStreamEventThreadRunCancelling struct {
	// Represents an execution run on a
	// [thread](https://platform.openai.com/docs/api-reference/threads).
	Data Run `json:"data,omitzero,required"`
	// This field can be elided, and will be automatically set as
	// "thread.run.cancelling".
	Event constant.ThreadRunCancelling `json:"event,required"`
	JSON  struct {
		Data  resp.Field
		Event resp.Field
		raw   string
	} `json:"-"`
}

func (r AssistantStreamEventThreadRunCancelling) RawJSON() string { return r.JSON.raw }
func (r *AssistantStreamEventThreadRunCancelling) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Occurs when a [run](https://platform.openai.com/docs/api-reference/runs/object)
// is cancelled.
type AssistantStreamEventThreadRunCancelled struct {
	// Represents an execution run on a
	// [thread](https://platform.openai.com/docs/api-reference/threads).
	Data Run `json:"data,omitzero,required"`
	// This field can be elided, and will be automatically set as
	// "thread.run.cancelled".
	Event constant.ThreadRunCancelled `json:"event,required"`
	JSON  struct {
		Data  resp.Field
		Event resp.Field
		raw   string
	} `json:"-"`
}

func (r AssistantStreamEventThreadRunCancelled) RawJSON() string { return r.JSON.raw }
func (r *AssistantStreamEventThreadRunCancelled) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Occurs when a [run](https://platform.openai.com/docs/api-reference/runs/object)
// expires.
type AssistantStreamEventThreadRunExpired struct {
	// Represents an execution run on a
	// [thread](https://platform.openai.com/docs/api-reference/threads).
	Data Run `json:"data,omitzero,required"`
	// This field can be elided, and will be automatically set as "thread.run.expired".
	Event constant.ThreadRunExpired `json:"event,required"`
	JSON  struct {
		Data  resp.Field
		Event resp.Field
		raw   string
	} `json:"-"`
}

func (r AssistantStreamEventThreadRunExpired) RawJSON() string { return r.JSON.raw }
func (r *AssistantStreamEventThreadRunExpired) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Occurs when a
// [run step](https://platform.openai.com/docs/api-reference/run-steps/step-object)
// is created.
type AssistantStreamEventThreadRunStepCreated struct {
	// Represents a step in execution of a run.
	Data RunStep `json:"data,omitzero,required"`
	// This field can be elided, and will be automatically set as
	// "thread.run.step.created".
	Event constant.ThreadRunStepCreated `json:"event,required"`
	JSON  struct {
		Data  resp.Field
		Event resp.Field
		raw   string
	} `json:"-"`
}

func (r AssistantStreamEventThreadRunStepCreated) RawJSON() string { return r.JSON.raw }
func (r *AssistantStreamEventThreadRunStepCreated) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Occurs when a
// [run step](https://platform.openai.com/docs/api-reference/run-steps/step-object)
// moves to an `in_progress` state.
type AssistantStreamEventThreadRunStepInProgress struct {
	// Represents a step in execution of a run.
	Data RunStep `json:"data,omitzero,required"`
	// This field can be elided, and will be automatically set as
	// "thread.run.step.in_progress".
	Event constant.ThreadRunStepInProgress `json:"event,required"`
	JSON  struct {
		Data  resp.Field
		Event resp.Field
		raw   string
	} `json:"-"`
}

func (r AssistantStreamEventThreadRunStepInProgress) RawJSON() string { return r.JSON.raw }
func (r *AssistantStreamEventThreadRunStepInProgress) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Occurs when parts of a
// [run step](https://platform.openai.com/docs/api-reference/run-steps/step-object)
// are being streamed.
type AssistantStreamEventThreadRunStepDelta struct {
	// Represents a run step delta i.e. any changed fields on a run step during
	// streaming.
	Data RunStepDeltaEvent `json:"data,omitzero,required"`
	// This field can be elided, and will be automatically set as
	// "thread.run.step.delta".
	Event constant.ThreadRunStepDelta `json:"event,required"`
	JSON  struct {
		Data  resp.Field
		Event resp.Field
		raw   string
	} `json:"-"`
}

func (r AssistantStreamEventThreadRunStepDelta) RawJSON() string { return r.JSON.raw }
func (r *AssistantStreamEventThreadRunStepDelta) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Occurs when a
// [run step](https://platform.openai.com/docs/api-reference/run-steps/step-object)
// is completed.
type AssistantStreamEventThreadRunStepCompleted struct {
	// Represents a step in execution of a run.
	Data RunStep `json:"data,omitzero,required"`
	// This field can be elided, and will be automatically set as
	// "thread.run.step.completed".
	Event constant.ThreadRunStepCompleted `json:"event,required"`
	JSON  struct {
		Data  resp.Field
		Event resp.Field
		raw   string
	} `json:"-"`
}

func (r AssistantStreamEventThreadRunStepCompleted) RawJSON() string { return r.JSON.raw }
func (r *AssistantStreamEventThreadRunStepCompleted) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Occurs when a
// [run step](https://platform.openai.com/docs/api-reference/run-steps/step-object)
// fails.
type AssistantStreamEventThreadRunStepFailed struct {
	// Represents a step in execution of a run.
	Data RunStep `json:"data,omitzero,required"`
	// This field can be elided, and will be automatically set as
	// "thread.run.step.failed".
	Event constant.ThreadRunStepFailed `json:"event,required"`
	JSON  struct {
		Data  resp.Field
		Event resp.Field
		raw   string
	} `json:"-"`
}

func (r AssistantStreamEventThreadRunStepFailed) RawJSON() string { return r.JSON.raw }
func (r *AssistantStreamEventThreadRunStepFailed) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Occurs when a
// [run step](https://platform.openai.com/docs/api-reference/run-steps/step-object)
// is cancelled.
type AssistantStreamEventThreadRunStepCancelled struct {
	// Represents a step in execution of a run.
	Data RunStep `json:"data,omitzero,required"`
	// This field can be elided, and will be automatically set as
	// "thread.run.step.cancelled".
	Event constant.ThreadRunStepCancelled `json:"event,required"`
	JSON  struct {
		Data  resp.Field
		Event resp.Field
		raw   string
	} `json:"-"`
}

func (r AssistantStreamEventThreadRunStepCancelled) RawJSON() string { return r.JSON.raw }
func (r *AssistantStreamEventThreadRunStepCancelled) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Occurs when a
// [run step](https://platform.openai.com/docs/api-reference/run-steps/step-object)
// expires.
type AssistantStreamEventThreadRunStepExpired struct {
	// Represents a step in execution of a run.
	Data RunStep `json:"data,omitzero,required"`
	// This field can be elided, and will be automatically set as
	// "thread.run.step.expired".
	Event constant.ThreadRunStepExpired `json:"event,required"`
	JSON  struct {
		Data  resp.Field
		Event resp.Field
		raw   string
	} `json:"-"`
}

func (r AssistantStreamEventThreadRunStepExpired) RawJSON() string { return r.JSON.raw }
func (r *AssistantStreamEventThreadRunStepExpired) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Occurs when a
// [message](https://platform.openai.com/docs/api-reference/messages/object) is
// created.
type AssistantStreamEventThreadMessageCreated struct {
	// Represents a message within a
	// [thread](https://platform.openai.com/docs/api-reference/threads).
	Data Message `json:"data,omitzero,required"`
	// This field can be elided, and will be automatically set as
	// "thread.message.created".
	Event constant.ThreadMessageCreated `json:"event,required"`
	JSON  struct {
		Data  resp.Field
		Event resp.Field
		raw   string
	} `json:"-"`
}

func (r AssistantStreamEventThreadMessageCreated) RawJSON() string { return r.JSON.raw }
func (r *AssistantStreamEventThreadMessageCreated) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Occurs when a
// [message](https://platform.openai.com/docs/api-reference/messages/object) moves
// to an `in_progress` state.
type AssistantStreamEventThreadMessageInProgress struct {
	// Represents a message within a
	// [thread](https://platform.openai.com/docs/api-reference/threads).
	Data Message `json:"data,omitzero,required"`
	// This field can be elided, and will be automatically set as
	// "thread.message.in_progress".
	Event constant.ThreadMessageInProgress `json:"event,required"`
	JSON  struct {
		Data  resp.Field
		Event resp.Field
		raw   string
	} `json:"-"`
}

func (r AssistantStreamEventThreadMessageInProgress) RawJSON() string { return r.JSON.raw }
func (r *AssistantStreamEventThreadMessageInProgress) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Occurs when parts of a
// [Message](https://platform.openai.com/docs/api-reference/messages/object) are
// being streamed.
type AssistantStreamEventThreadMessageDelta struct {
	// Represents a message delta i.e. any changed fields on a message during
	// streaming.
	Data MessageDeltaEvent `json:"data,omitzero,required"`
	// This field can be elided, and will be automatically set as
	// "thread.message.delta".
	Event constant.ThreadMessageDelta `json:"event,required"`
	JSON  struct {
		Data  resp.Field
		Event resp.Field
		raw   string
	} `json:"-"`
}

func (r AssistantStreamEventThreadMessageDelta) RawJSON() string { return r.JSON.raw }
func (r *AssistantStreamEventThreadMessageDelta) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Occurs when a
// [message](https://platform.openai.com/docs/api-reference/messages/object) is
// completed.
type AssistantStreamEventThreadMessageCompleted struct {
	// Represents a message within a
	// [thread](https://platform.openai.com/docs/api-reference/threads).
	Data Message `json:"data,omitzero,required"`
	// This field can be elided, and will be automatically set as
	// "thread.message.completed".
	Event constant.ThreadMessageCompleted `json:"event,required"`
	JSON  struct {
		Data  resp.Field
		Event resp.Field
		raw   string
	} `json:"-"`
}

func (r AssistantStreamEventThreadMessageCompleted) RawJSON() string { return r.JSON.raw }
func (r *AssistantStreamEventThreadMessageCompleted) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Occurs when a
// [message](https://platform.openai.com/docs/api-reference/messages/object) ends
// before it is completed.
type AssistantStreamEventThreadMessageIncomplete struct {
	// Represents a message within a
	// [thread](https://platform.openai.com/docs/api-reference/threads).
	Data Message `json:"data,omitzero,required"`
	// This field can be elided, and will be automatically set as
	// "thread.message.incomplete".
	Event constant.ThreadMessageIncomplete `json:"event,required"`
	JSON  struct {
		Data  resp.Field
		Event resp.Field
		raw   string
	} `json:"-"`
}

func (r AssistantStreamEventThreadMessageIncomplete) RawJSON() string { return r.JSON.raw }
func (r *AssistantStreamEventThreadMessageIncomplete) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Occurs when an
// [error](https://platform.openai.com/docs/guides/error-codes#api-errors) occurs.
// This can happen due to an internal server error or a timeout.
type AssistantStreamEventErrorEvent struct {
	Data shared.ErrorObject `json:"data,omitzero,required"`
	// This field can be elided, and will be automatically set as "error".
	Event constant.Error `json:"event,required"`
	JSON  struct {
		Data  resp.Field
		Event resp.Field
		raw   string
	} `json:"-"`
}

func (r AssistantStreamEventErrorEvent) RawJSON() string { return r.JSON.raw }
func (r *AssistantStreamEventErrorEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AssistantToolUnion struct {
	Type       string                    `json:"type"`
	FileSearch FileSearchToolFileSearch  `json:"file_search"`
	Function   shared.FunctionDefinition `json:"function"`
	JSON       struct {
		Type       resp.Field
		FileSearch resp.Field
		Function   resp.Field
		raw        string
	} `json:"-"`
}

// note: this function is generated only for discriminated unions
func (u AssistantToolUnion) Variant() (res struct {
	OfCodeInterpreter *CodeInterpreterTool
	OfFileSearch      *FileSearchTool
	OfFunction        *FunctionTool
}) {
	switch u.Type {
	case "code_interpreter":
		v := u.AsCodeInterpreter()
		res.OfCodeInterpreter = &v
	case "file_search":
		v := u.AsFileSearch()
		res.OfFileSearch = &v
	case "function":
		v := u.AsFunction()
		res.OfFunction = &v
	}
	return
}

func (u AssistantToolUnion) WhichKind() string {
	return u.Type
}

func (u AssistantToolUnion) AsCodeInterpreter() (v CodeInterpreterTool) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AssistantToolUnion) AsFileSearch() (v FileSearchTool) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AssistantToolUnion) AsFunction() (v FunctionTool) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AssistantToolUnion) RawJSON() string { return u.JSON.raw }

func (r *AssistantToolUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this AssistantToolUnion to a AssistantToolUnionParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// AssistantToolUnionParam.IsOverridden()
func (r AssistantToolUnion) ToParam() AssistantToolUnionParam {
	return param.Override[AssistantToolUnionParam](r.RawJSON())
}

func NewAssistantToolOfFunction(function shared.FunctionDefinitionParam) AssistantToolUnionParam {
	var variant FunctionToolParam
	variant.Function = function
	return AssistantToolUnionParam{OfFunction: &variant}
}

// Only one field can be non-zero
type AssistantToolUnionParam struct {
	OfCodeInterpreter *CodeInterpreterToolParam
	OfFileSearch      *FileSearchToolParam
	OfFunction        *FunctionToolParam
	apiunion
}

func (u AssistantToolUnionParam) IsMissing() bool { return param.IsOmitted(u) || u.IsNull() }

func (u AssistantToolUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion[AssistantToolUnionParam](u.OfCodeInterpreter, u.OfFileSearch, u.OfFunction)
}

func (u AssistantToolUnionParam) GetFileSearch() *FileSearchToolFileSearchParam {
	if vt := u.OfFileSearch; vt != nil {
		return &vt.FileSearch
	}
	return nil
}

func (u AssistantToolUnionParam) GetFunction() *shared.FunctionDefinitionParam {
	if vt := u.OfFunction; vt != nil {
		return &vt.Function
	}
	return nil
}

func (u AssistantToolUnionParam) GetType() *string {
	if vt := u.OfCodeInterpreter; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfFileSearch; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfFunction; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

type CodeInterpreterTool struct {
	// The type of tool being defined: `code_interpreter`
	//
	// This field can be elided, and will be automatically set as "code_interpreter".
	Type constant.CodeInterpreter `json:"type,required"`
	JSON struct {
		Type resp.Field
		raw  string
	} `json:"-"`
}

func (r CodeInterpreterTool) RawJSON() string { return r.JSON.raw }
func (r *CodeInterpreterTool) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this CodeInterpreterTool to a CodeInterpreterToolParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// CodeInterpreterToolParam.IsOverridden()
func (r CodeInterpreterTool) ToParam() CodeInterpreterToolParam {
	return param.Override[CodeInterpreterToolParam](r.RawJSON())
}

type CodeInterpreterToolParam struct {
	// The type of tool being defined: `code_interpreter`
	//
	// This field can be elided, and will be automatically set as "code_interpreter".
	Type constant.CodeInterpreter `json:"type,required"`
	apiobject
}

func (f CodeInterpreterToolParam) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

func (r CodeInterpreterToolParam) MarshalJSON() (data []byte, err error) {
	type shadow CodeInterpreterToolParam
	return param.MarshalObject(r, (*shadow)(&r))
}

type FileSearchTool struct {
	// The type of tool being defined: `file_search`
	//
	// This field can be elided, and will be automatically set as "file_search".
	Type constant.FileSearch `json:"type,required"`
	// Overrides for the file search tool.
	FileSearch FileSearchToolFileSearch `json:"file_search,omitzero"`
	JSON       struct {
		Type       resp.Field
		FileSearch resp.Field
		raw        string
	} `json:"-"`
}

func (r FileSearchTool) RawJSON() string { return r.JSON.raw }
func (r *FileSearchTool) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this FileSearchTool to a FileSearchToolParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// FileSearchToolParam.IsOverridden()
func (r FileSearchTool) ToParam() FileSearchToolParam {
	return param.Override[FileSearchToolParam](r.RawJSON())
}

// Overrides for the file search tool.
type FileSearchToolFileSearch struct {
	// The maximum number of results the file search tool should output. The default is
	// 20 for `gpt-4*` models and 5 for `gpt-3.5-turbo`. This number should be between
	// 1 and 50 inclusive.
	//
	// Note that the file search tool may output fewer than `max_num_results` results.
	// See the
	// [file search tool documentation](https://platform.openai.com/docs/assistants/tools/file-search#customizing-file-search-settings)
	// for more information.
	MaxNumResults int64 `json:"max_num_results,omitzero"`
	// The ranking options for the file search. If not specified, the file search tool
	// will use the `auto` ranker and a score_threshold of 0.
	//
	// See the
	// [file search tool documentation](https://platform.openai.com/docs/assistants/tools/file-search#customizing-file-search-settings)
	// for more information.
	RankingOptions FileSearchToolFileSearchRankingOptions `json:"ranking_options,omitzero"`
	JSON           struct {
		MaxNumResults  resp.Field
		RankingOptions resp.Field
		raw            string
	} `json:"-"`
}

func (r FileSearchToolFileSearch) RawJSON() string { return r.JSON.raw }
func (r *FileSearchToolFileSearch) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The ranking options for the file search. If not specified, the file search tool
// will use the `auto` ranker and a score_threshold of 0.
//
// See the
// [file search tool documentation](https://platform.openai.com/docs/assistants/tools/file-search#customizing-file-search-settings)
// for more information.
type FileSearchToolFileSearchRankingOptions struct {
	// The score threshold for the file search. All values must be a floating point
	// number between 0 and 1.
	ScoreThreshold float64 `json:"score_threshold,omitzero,required"`
	// The ranker to use for the file search. If not specified will use the `auto`
	// ranker.
	//
	// Any of "auto", "default_2024_08_21"
	Ranker string `json:"ranker,omitzero"`
	JSON   struct {
		ScoreThreshold resp.Field
		Ranker         resp.Field
		raw            string
	} `json:"-"`
}

func (r FileSearchToolFileSearchRankingOptions) RawJSON() string { return r.JSON.raw }
func (r *FileSearchToolFileSearchRankingOptions) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The ranker to use for the file search. If not specified will use the `auto`
// ranker.
type FileSearchToolFileSearchRankingOptionsRanker = string

const (
	FileSearchToolFileSearchRankingOptionsRankerAuto              FileSearchToolFileSearchRankingOptionsRanker = "auto"
	FileSearchToolFileSearchRankingOptionsRankerDefault2024_08_21 FileSearchToolFileSearchRankingOptionsRanker = "default_2024_08_21"
)

type FileSearchToolParam struct {
	// The type of tool being defined: `file_search`
	//
	// This field can be elided, and will be automatically set as "file_search".
	Type constant.FileSearch `json:"type,required"`
	// Overrides for the file search tool.
	FileSearch FileSearchToolFileSearchParam `json:"file_search,omitzero"`
	apiobject
}

func (f FileSearchToolParam) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

func (r FileSearchToolParam) MarshalJSON() (data []byte, err error) {
	type shadow FileSearchToolParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// Overrides for the file search tool.
type FileSearchToolFileSearchParam struct {
	// The maximum number of results the file search tool should output. The default is
	// 20 for `gpt-4*` models and 5 for `gpt-3.5-turbo`. This number should be between
	// 1 and 50 inclusive.
	//
	// Note that the file search tool may output fewer than `max_num_results` results.
	// See the
	// [file search tool documentation](https://platform.openai.com/docs/assistants/tools/file-search#customizing-file-search-settings)
	// for more information.
	MaxNumResults param.Int `json:"max_num_results,omitzero"`
	// The ranking options for the file search. If not specified, the file search tool
	// will use the `auto` ranker and a score_threshold of 0.
	//
	// See the
	// [file search tool documentation](https://platform.openai.com/docs/assistants/tools/file-search#customizing-file-search-settings)
	// for more information.
	RankingOptions FileSearchToolFileSearchRankingOptionsParam `json:"ranking_options,omitzero"`
	apiobject
}

func (f FileSearchToolFileSearchParam) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

func (r FileSearchToolFileSearchParam) MarshalJSON() (data []byte, err error) {
	type shadow FileSearchToolFileSearchParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// The ranking options for the file search. If not specified, the file search tool
// will use the `auto` ranker and a score_threshold of 0.
//
// See the
// [file search tool documentation](https://platform.openai.com/docs/assistants/tools/file-search#customizing-file-search-settings)
// for more information.
type FileSearchToolFileSearchRankingOptionsParam struct {
	// The score threshold for the file search. All values must be a floating point
	// number between 0 and 1.
	ScoreThreshold param.Float `json:"score_threshold,omitzero,required"`
	// The ranker to use for the file search. If not specified will use the `auto`
	// ranker.
	//
	// Any of "auto", "default_2024_08_21"
	Ranker string `json:"ranker,omitzero"`
	apiobject
}

func (f FileSearchToolFileSearchRankingOptionsParam) IsMissing() bool {
	return param.IsOmitted(f) || f.IsNull()
}

func (r FileSearchToolFileSearchRankingOptionsParam) MarshalJSON() (data []byte, err error) {
	type shadow FileSearchToolFileSearchRankingOptionsParam
	return param.MarshalObject(r, (*shadow)(&r))
}

type FunctionTool struct {
	Function shared.FunctionDefinition `json:"function,omitzero,required"`
	// The type of tool being defined: `function`
	//
	// This field can be elided, and will be automatically set as "function".
	Type constant.Function `json:"type,required"`
	JSON struct {
		Function resp.Field
		Type     resp.Field
		raw      string
	} `json:"-"`
}

func (r FunctionTool) RawJSON() string { return r.JSON.raw }
func (r *FunctionTool) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this FunctionTool to a FunctionToolParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// FunctionToolParam.IsOverridden()
func (r FunctionTool) ToParam() FunctionToolParam {
	return param.Override[FunctionToolParam](r.RawJSON())
}

type FunctionToolParam struct {
	Function shared.FunctionDefinitionParam `json:"function,omitzero,required"`
	// The type of tool being defined: `function`
	//
	// This field can be elided, and will be automatically set as "function".
	Type constant.Function `json:"type,required"`
	apiobject
}

func (f FunctionToolParam) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

func (r FunctionToolParam) MarshalJSON() (data []byte, err error) {
	type shadow FunctionToolParam
	return param.MarshalObject(r, (*shadow)(&r))
}

type BetaAssistantNewParams struct {
	// ID of the model to use. You can use the
	// [List models](https://platform.openai.com/docs/api-reference/models/list) API to
	// see all of your available models, or see our
	// [Model overview](https://platform.openai.com/docs/models) for descriptions of
	// them.
	Model ChatModel `json:"model,omitzero,required"`
	// The description of the assistant. The maximum length is 512 characters.
	Description param.String `json:"description,omitzero"`
	// The system instructions that the assistant uses. The maximum length is 256,000
	// characters.
	Instructions param.String `json:"instructions,omitzero"`
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard.
	//
	// Keys are strings with a maximum length of 64 characters. Values are strings with
	// a maximum length of 512 characters.
	Metadata shared.MetadataParam `json:"metadata,omitzero"`
	// The name of the assistant. The maximum length is 256 characters.
	Name param.String `json:"name,omitzero"`
	// **o1 and o3-mini models only**
	//
	// Constrains effort on reasoning for
	// [reasoning models](https://platform.openai.com/docs/guides/reasoning). Currently
	// supported values are `low`, `medium`, and `high`. Reducing reasoning effort can
	// result in faster responses and fewer tokens used on reasoning in a response.
	//
	// Any of "low", "medium", "high"
	ReasoningEffort BetaAssistantNewParamsReasoningEffort `json:"reasoning_effort,omitzero"`
	// What sampling temperature to use, between 0 and 2. Higher values like 0.8 will
	// make the output more random, while lower values like 0.2 will make it more
	// focused and deterministic.
	Temperature param.Float `json:"temperature,omitzero"`
	// A set of resources that are used by the assistant's tools. The resources are
	// specific to the type of tool. For example, the `code_interpreter` tool requires
	// a list of file IDs, while the `file_search` tool requires a list of vector store
	// IDs.
	ToolResources BetaAssistantNewParamsToolResources `json:"tool_resources,omitzero"`
	// A list of tool enabled on the assistant. There can be a maximum of 128 tools per
	// assistant. Tools can be of types `code_interpreter`, `file_search`, or
	// `function`.
	Tools []AssistantToolUnionParam `json:"tools,omitzero"`
	// An alternative to sampling with temperature, called nucleus sampling, where the
	// model considers the results of the tokens with top_p probability mass. So 0.1
	// means only the tokens comprising the top 10% probability mass are considered.
	//
	// We generally recommend altering this or temperature but not both.
	TopP param.Float `json:"top_p,omitzero"`
	apiobject
}

func (f BetaAssistantNewParams) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

func (r BetaAssistantNewParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaAssistantNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}

// **o1 and o3-mini models only**
//
// Constrains effort on reasoning for
// [reasoning models](https://platform.openai.com/docs/guides/reasoning). Currently
// supported values are `low`, `medium`, and `high`. Reducing reasoning effort can
// result in faster responses and fewer tokens used on reasoning in a response.
type BetaAssistantNewParamsReasoningEffort string

const (
	BetaAssistantNewParamsReasoningEffortLow    BetaAssistantNewParamsReasoningEffort = "low"
	BetaAssistantNewParamsReasoningEffortMedium BetaAssistantNewParamsReasoningEffort = "medium"
	BetaAssistantNewParamsReasoningEffortHigh   BetaAssistantNewParamsReasoningEffort = "high"
)

// A set of resources that are used by the assistant's tools. The resources are
// specific to the type of tool. For example, the `code_interpreter` tool requires
// a list of file IDs, while the `file_search` tool requires a list of vector store
// IDs.
type BetaAssistantNewParamsToolResources struct {
	CodeInterpreter BetaAssistantNewParamsToolResourcesCodeInterpreter `json:"code_interpreter,omitzero"`
	FileSearch      BetaAssistantNewParamsToolResourcesFileSearch      `json:"file_search,omitzero"`
	apiobject
}

func (f BetaAssistantNewParamsToolResources) IsMissing() bool {
	return param.IsOmitted(f) || f.IsNull()
}

func (r BetaAssistantNewParamsToolResources) MarshalJSON() (data []byte, err error) {
	type shadow BetaAssistantNewParamsToolResources
	return param.MarshalObject(r, (*shadow)(&r))
}

type BetaAssistantNewParamsToolResourcesCodeInterpreter struct {
	// A list of [file](https://platform.openai.com/docs/api-reference/files) IDs made
	// available to the `code_interpreter` tool. There can be a maximum of 20 files
	// associated with the tool.
	FileIDs []string `json:"file_ids,omitzero"`
	apiobject
}

func (f BetaAssistantNewParamsToolResourcesCodeInterpreter) IsMissing() bool {
	return param.IsOmitted(f) || f.IsNull()
}

func (r BetaAssistantNewParamsToolResourcesCodeInterpreter) MarshalJSON() (data []byte, err error) {
	type shadow BetaAssistantNewParamsToolResourcesCodeInterpreter
	return param.MarshalObject(r, (*shadow)(&r))
}

type BetaAssistantNewParamsToolResourcesFileSearch struct {
	// The
	// [vector store](https://platform.openai.com/docs/api-reference/vector-stores/object)
	// attached to this assistant. There can be a maximum of 1 vector store attached to
	// the assistant.
	VectorStoreIDs []string `json:"vector_store_ids,omitzero"`
	// A helper to create a
	// [vector store](https://platform.openai.com/docs/api-reference/vector-stores/object)
	// with file_ids and attach it to this assistant. There can be a maximum of 1
	// vector store attached to the assistant.
	VectorStores []BetaAssistantNewParamsToolResourcesFileSearchVectorStore `json:"vector_stores,omitzero"`
	apiobject
}

func (f BetaAssistantNewParamsToolResourcesFileSearch) IsMissing() bool {
	return param.IsOmitted(f) || f.IsNull()
}

func (r BetaAssistantNewParamsToolResourcesFileSearch) MarshalJSON() (data []byte, err error) {
	type shadow BetaAssistantNewParamsToolResourcesFileSearch
	return param.MarshalObject(r, (*shadow)(&r))
}

type BetaAssistantNewParamsToolResourcesFileSearchVectorStore struct {
	// The chunking strategy used to chunk the file(s). If not set, will use the `auto`
	// strategy. Only applicable if `file_ids` is non-empty.
	ChunkingStrategy FileChunkingStrategyParamUnion `json:"chunking_strategy,omitzero"`
	// A list of [file](https://platform.openai.com/docs/api-reference/files) IDs to
	// add to the vector store. There can be a maximum of 10000 files in a vector
	// store.
	FileIDs []string `json:"file_ids,omitzero"`
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard.
	//
	// Keys are strings with a maximum length of 64 characters. Values are strings with
	// a maximum length of 512 characters.
	Metadata shared.MetadataParam `json:"metadata,omitzero"`
	apiobject
}

func (f BetaAssistantNewParamsToolResourcesFileSearchVectorStore) IsMissing() bool {
	return param.IsOmitted(f) || f.IsNull()
}

func (r BetaAssistantNewParamsToolResourcesFileSearchVectorStore) MarshalJSON() (data []byte, err error) {
	type shadow BetaAssistantNewParamsToolResourcesFileSearchVectorStore
	return param.MarshalObject(r, (*shadow)(&r))
}

type BetaAssistantUpdateParams struct {
	// The description of the assistant. The maximum length is 512 characters.
	Description param.String `json:"description,omitzero"`
	// The system instructions that the assistant uses. The maximum length is 256,000
	// characters.
	Instructions param.String `json:"instructions,omitzero"`
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard.
	//
	// Keys are strings with a maximum length of 64 characters. Values are strings with
	// a maximum length of 512 characters.
	Metadata shared.MetadataParam `json:"metadata,omitzero"`
	// ID of the model to use. You can use the
	// [List models](https://platform.openai.com/docs/api-reference/models/list) API to
	// see all of your available models, or see our
	// [Model overview](https://platform.openai.com/docs/models) for descriptions of
	// them.
	Model string `json:"model,omitzero"`
	// The name of the assistant. The maximum length is 256 characters.
	Name param.String `json:"name,omitzero"`
	// **o1 and o3-mini models only**
	//
	// Constrains effort on reasoning for
	// [reasoning models](https://platform.openai.com/docs/guides/reasoning). Currently
	// supported values are `low`, `medium`, and `high`. Reducing reasoning effort can
	// result in faster responses and fewer tokens used on reasoning in a response.
	//
	// Any of "low", "medium", "high"
	ReasoningEffort BetaAssistantUpdateParamsReasoningEffort `json:"reasoning_effort,omitzero"`
	// What sampling temperature to use, between 0 and 2. Higher values like 0.8 will
	// make the output more random, while lower values like 0.2 will make it more
	// focused and deterministic.
	Temperature param.Float `json:"temperature,omitzero"`
	// A set of resources that are used by the assistant's tools. The resources are
	// specific to the type of tool. For example, the `code_interpreter` tool requires
	// a list of file IDs, while the `file_search` tool requires a list of vector store
	// IDs.
	ToolResources BetaAssistantUpdateParamsToolResources `json:"tool_resources,omitzero"`
	// A list of tool enabled on the assistant. There can be a maximum of 128 tools per
	// assistant. Tools can be of types `code_interpreter`, `file_search`, or
	// `function`.
	Tools []AssistantToolUnionParam `json:"tools,omitzero"`
	// An alternative to sampling with temperature, called nucleus sampling, where the
	// model considers the results of the tokens with top_p probability mass. So 0.1
	// means only the tokens comprising the top 10% probability mass are considered.
	//
	// We generally recommend altering this or temperature but not both.
	TopP param.Float `json:"top_p,omitzero"`
	apiobject
}

func (f BetaAssistantUpdateParams) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

func (r BetaAssistantUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaAssistantUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}

// ID of the model to use. You can use the
// [List models](https://platform.openai.com/docs/api-reference/models/list) API to
// see all of your available models, or see our
// [Model overview](https://platform.openai.com/docs/models) for descriptions of
// them.
type BetaAssistantUpdateParamsModel = string

const (
	BetaAssistantUpdateParamsModelO3Mini              BetaAssistantUpdateParamsModel = "o3-mini"
	BetaAssistantUpdateParamsModelO3Mini2025_01_31    BetaAssistantUpdateParamsModel = "o3-mini-2025-01-31"
	BetaAssistantUpdateParamsModelO1                  BetaAssistantUpdateParamsModel = "o1"
	BetaAssistantUpdateParamsModelO1_2024_12_17       BetaAssistantUpdateParamsModel = "o1-2024-12-17"
	BetaAssistantUpdateParamsModelGPT4o               BetaAssistantUpdateParamsModel = "gpt-4o"
	BetaAssistantUpdateParamsModelGPT4o2024_11_20     BetaAssistantUpdateParamsModel = "gpt-4o-2024-11-20"
	BetaAssistantUpdateParamsModelGPT4o2024_08_06     BetaAssistantUpdateParamsModel = "gpt-4o-2024-08-06"
	BetaAssistantUpdateParamsModelGPT4o2024_05_13     BetaAssistantUpdateParamsModel = "gpt-4o-2024-05-13"
	BetaAssistantUpdateParamsModelGPT4oMini           BetaAssistantUpdateParamsModel = "gpt-4o-mini"
	BetaAssistantUpdateParamsModelGPT4oMini2024_07_18 BetaAssistantUpdateParamsModel = "gpt-4o-mini-2024-07-18"
	BetaAssistantUpdateParamsModelGPT4Turbo           BetaAssistantUpdateParamsModel = "gpt-4-turbo"
	BetaAssistantUpdateParamsModelGPT4Turbo2024_04_09 BetaAssistantUpdateParamsModel = "gpt-4-turbo-2024-04-09"
	BetaAssistantUpdateParamsModelGPT4_0125Preview    BetaAssistantUpdateParamsModel = "gpt-4-0125-preview"
	BetaAssistantUpdateParamsModelGPT4TurboPreview    BetaAssistantUpdateParamsModel = "gpt-4-turbo-preview"
	BetaAssistantUpdateParamsModelGPT4_1106Preview    BetaAssistantUpdateParamsModel = "gpt-4-1106-preview"
	BetaAssistantUpdateParamsModelGPT4VisionPreview   BetaAssistantUpdateParamsModel = "gpt-4-vision-preview"
	BetaAssistantUpdateParamsModelGPT4                BetaAssistantUpdateParamsModel = "gpt-4"
	BetaAssistantUpdateParamsModelGPT4_0314           BetaAssistantUpdateParamsModel = "gpt-4-0314"
	BetaAssistantUpdateParamsModelGPT4_0613           BetaAssistantUpdateParamsModel = "gpt-4-0613"
	BetaAssistantUpdateParamsModelGPT4_32k            BetaAssistantUpdateParamsModel = "gpt-4-32k"
	BetaAssistantUpdateParamsModelGPT4_32k0314        BetaAssistantUpdateParamsModel = "gpt-4-32k-0314"
	BetaAssistantUpdateParamsModelGPT4_32k0613        BetaAssistantUpdateParamsModel = "gpt-4-32k-0613"
	BetaAssistantUpdateParamsModelGPT3_5Turbo         BetaAssistantUpdateParamsModel = "gpt-3.5-turbo"
	BetaAssistantUpdateParamsModelGPT3_5Turbo16k      BetaAssistantUpdateParamsModel = "gpt-3.5-turbo-16k"
	BetaAssistantUpdateParamsModelGPT3_5Turbo0613     BetaAssistantUpdateParamsModel = "gpt-3.5-turbo-0613"
	BetaAssistantUpdateParamsModelGPT3_5Turbo1106     BetaAssistantUpdateParamsModel = "gpt-3.5-turbo-1106"
	BetaAssistantUpdateParamsModelGPT3_5Turbo0125     BetaAssistantUpdateParamsModel = "gpt-3.5-turbo-0125"
	BetaAssistantUpdateParamsModelGPT3_5Turbo16k0613  BetaAssistantUpdateParamsModel = "gpt-3.5-turbo-16k-0613"
)

// **o1 and o3-mini models only**
//
// Constrains effort on reasoning for
// [reasoning models](https://platform.openai.com/docs/guides/reasoning). Currently
// supported values are `low`, `medium`, and `high`. Reducing reasoning effort can
// result in faster responses and fewer tokens used on reasoning in a response.
type BetaAssistantUpdateParamsReasoningEffort string

const (
	BetaAssistantUpdateParamsReasoningEffortLow    BetaAssistantUpdateParamsReasoningEffort = "low"
	BetaAssistantUpdateParamsReasoningEffortMedium BetaAssistantUpdateParamsReasoningEffort = "medium"
	BetaAssistantUpdateParamsReasoningEffortHigh   BetaAssistantUpdateParamsReasoningEffort = "high"
)

// A set of resources that are used by the assistant's tools. The resources are
// specific to the type of tool. For example, the `code_interpreter` tool requires
// a list of file IDs, while the `file_search` tool requires a list of vector store
// IDs.
type BetaAssistantUpdateParamsToolResources struct {
	CodeInterpreter BetaAssistantUpdateParamsToolResourcesCodeInterpreter `json:"code_interpreter,omitzero"`
	FileSearch      BetaAssistantUpdateParamsToolResourcesFileSearch      `json:"file_search,omitzero"`
	apiobject
}

func (f BetaAssistantUpdateParamsToolResources) IsMissing() bool {
	return param.IsOmitted(f) || f.IsNull()
}

func (r BetaAssistantUpdateParamsToolResources) MarshalJSON() (data []byte, err error) {
	type shadow BetaAssistantUpdateParamsToolResources
	return param.MarshalObject(r, (*shadow)(&r))
}

type BetaAssistantUpdateParamsToolResourcesCodeInterpreter struct {
	// Overrides the list of
	// [file](https://platform.openai.com/docs/api-reference/files) IDs made available
	// to the `code_interpreter` tool. There can be a maximum of 20 files associated
	// with the tool.
	FileIDs []string `json:"file_ids,omitzero"`
	apiobject
}

func (f BetaAssistantUpdateParamsToolResourcesCodeInterpreter) IsMissing() bool {
	return param.IsOmitted(f) || f.IsNull()
}

func (r BetaAssistantUpdateParamsToolResourcesCodeInterpreter) MarshalJSON() (data []byte, err error) {
	type shadow BetaAssistantUpdateParamsToolResourcesCodeInterpreter
	return param.MarshalObject(r, (*shadow)(&r))
}

type BetaAssistantUpdateParamsToolResourcesFileSearch struct {
	// Overrides the
	// [vector store](https://platform.openai.com/docs/api-reference/vector-stores/object)
	// attached to this assistant. There can be a maximum of 1 vector store attached to
	// the assistant.
	VectorStoreIDs []string `json:"vector_store_ids,omitzero"`
	apiobject
}

func (f BetaAssistantUpdateParamsToolResourcesFileSearch) IsMissing() bool {
	return param.IsOmitted(f) || f.IsNull()
}

func (r BetaAssistantUpdateParamsToolResourcesFileSearch) MarshalJSON() (data []byte, err error) {
	type shadow BetaAssistantUpdateParamsToolResourcesFileSearch
	return param.MarshalObject(r, (*shadow)(&r))
}

type BetaAssistantListParams struct {
	// A cursor for use in pagination. `after` is an object ID that defines your place
	// in the list. For instance, if you make a list request and receive 100 objects,
	// ending with obj_foo, your subsequent call can include after=obj_foo in order to
	// fetch the next page of the list.
	After param.String `query:"after,omitzero"`
	// A cursor for use in pagination. `before` is an object ID that defines your place
	// in the list. For instance, if you make a list request and receive 100 objects,
	// starting with obj_foo, your subsequent call can include before=obj_foo in order
	// to fetch the previous page of the list.
	Before param.String `query:"before,omitzero"`
	// A limit on the number of objects to be returned. Limit can range between 1 and
	// 100, and the default is 20.
	Limit param.Int `query:"limit,omitzero"`
	// Sort order by the `created_at` timestamp of the objects. `asc` for ascending
	// order and `desc` for descending order.
	//
	// Any of "asc", "desc"
	Order BetaAssistantListParamsOrder `query:"order,omitzero"`
	apiobject
}

func (f BetaAssistantListParams) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

// URLQuery serializes [BetaAssistantListParams]'s query parameters as
// `url.Values`.
func (r BetaAssistantListParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Sort order by the `created_at` timestamp of the objects. `asc` for ascending
// order and `desc` for descending order.
type BetaAssistantListParamsOrder string

const (
	BetaAssistantListParamsOrderAsc  BetaAssistantListParamsOrder = "asc"
	BetaAssistantListParamsOrderDesc BetaAssistantListParamsOrder = "desc"
)

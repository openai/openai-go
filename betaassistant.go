// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"

	"github.com/openai/openai-go/internal/apijson"
	"github.com/openai/openai-go/internal/apiquery"
	"github.com/openai/openai-go/internal/param"
	"github.com/openai/openai-go/internal/requestconfig"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/pagination"
	"github.com/openai/openai-go/shared"
	"github.com/tidwall/gjson"
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
func NewBetaAssistantService(opts ...option.RequestOption) (r *BetaAssistantService) {
	r = &BetaAssistantService{}
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
	ID string `json:"id,required"`
	// The Unix timestamp (in seconds) for when the assistant was created.
	CreatedAt int64 `json:"created_at,required"`
	// The description of the assistant. The maximum length is 512 characters.
	Description string `json:"description,required,nullable"`
	// The system instructions that the assistant uses. The maximum length is 256,000
	// characters.
	Instructions string `json:"instructions,required,nullable"`
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard.
	//
	// Keys are strings with a maximum length of 64 characters. Values are strings with
	// a maximum length of 512 characters.
	Metadata shared.Metadata `json:"metadata,required,nullable"`
	// ID of the model to use. You can use the
	// [List models](https://platform.openai.com/docs/api-reference/models/list) API to
	// see all of your available models, or see our
	// [Model overview](https://platform.openai.com/docs/models) for descriptions of
	// them.
	Model string `json:"model,required"`
	// The name of the assistant. The maximum length is 256 characters.
	Name string `json:"name,required,nullable"`
	// The object type, which is always `assistant`.
	Object AssistantObject `json:"object,required"`
	// A list of tool enabled on the assistant. There can be a maximum of 128 tools per
	// assistant. Tools can be of types `code_interpreter`, `file_search`, or
	// `function`.
	Tools []AssistantTool `json:"tools,required"`
	// What sampling temperature to use, between 0 and 2. Higher values like 0.8 will
	// make the output more random, while lower values like 0.2 will make it more
	// focused and deterministic.
	Temperature float64 `json:"temperature,nullable"`
	// A set of resources that are used by the assistant's tools. The resources are
	// specific to the type of tool. For example, the `code_interpreter` tool requires
	// a list of file IDs, while the `file_search` tool requires a list of vector store
	// IDs.
	ToolResources AssistantToolResources `json:"tool_resources,nullable"`
	// An alternative to sampling with temperature, called nucleus sampling, where the
	// model considers the results of the tokens with top_p probability mass. So 0.1
	// means only the tokens comprising the top 10% probability mass are considered.
	//
	// We generally recommend altering this or temperature but not both.
	TopP float64       `json:"top_p,nullable"`
	JSON assistantJSON `json:"-"`
}

// assistantJSON contains the JSON metadata for the struct [Assistant]
type assistantJSON struct {
	ID            apijson.Field
	CreatedAt     apijson.Field
	Description   apijson.Field
	Instructions  apijson.Field
	Metadata      apijson.Field
	Model         apijson.Field
	Name          apijson.Field
	Object        apijson.Field
	Tools         apijson.Field
	Temperature   apijson.Field
	ToolResources apijson.Field
	TopP          apijson.Field
	raw           string
	ExtraFields   map[string]apijson.Field
}

func (r *Assistant) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r assistantJSON) RawJSON() string {
	return r.raw
}

// The object type, which is always `assistant`.
type AssistantObject string

const (
	AssistantObjectAssistant AssistantObject = "assistant"
)

func (r AssistantObject) IsKnown() bool {
	switch r {
	case AssistantObjectAssistant:
		return true
	}
	return false
}

// A set of resources that are used by the assistant's tools. The resources are
// specific to the type of tool. For example, the `code_interpreter` tool requires
// a list of file IDs, while the `file_search` tool requires a list of vector store
// IDs.
type AssistantToolResources struct {
	CodeInterpreter AssistantToolResourcesCodeInterpreter `json:"code_interpreter"`
	FileSearch      AssistantToolResourcesFileSearch      `json:"file_search"`
	JSON            assistantToolResourcesJSON            `json:"-"`
}

// assistantToolResourcesJSON contains the JSON metadata for the struct
// [AssistantToolResources]
type assistantToolResourcesJSON struct {
	CodeInterpreter apijson.Field
	FileSearch      apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r *AssistantToolResources) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r assistantToolResourcesJSON) RawJSON() string {
	return r.raw
}

type AssistantToolResourcesCodeInterpreter struct {
	// A list of [file](https://platform.openai.com/docs/api-reference/files) IDs made
	// available to the `code_interpreter“ tool. There can be a maximum of 20 files
	// associated with the tool.
	FileIDs []string                                  `json:"file_ids"`
	JSON    assistantToolResourcesCodeInterpreterJSON `json:"-"`
}

// assistantToolResourcesCodeInterpreterJSON contains the JSON metadata for the
// struct [AssistantToolResourcesCodeInterpreter]
type assistantToolResourcesCodeInterpreterJSON struct {
	FileIDs     apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *AssistantToolResourcesCodeInterpreter) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r assistantToolResourcesCodeInterpreterJSON) RawJSON() string {
	return r.raw
}

type AssistantToolResourcesFileSearch struct {
	// The ID of the
	// [vector store](https://platform.openai.com/docs/api-reference/vector-stores/object)
	// attached to this assistant. There can be a maximum of 1 vector store attached to
	// the assistant.
	VectorStoreIDs []string                             `json:"vector_store_ids"`
	JSON           assistantToolResourcesFileSearchJSON `json:"-"`
}

// assistantToolResourcesFileSearchJSON contains the JSON metadata for the struct
// [AssistantToolResourcesFileSearch]
type assistantToolResourcesFileSearchJSON struct {
	VectorStoreIDs apijson.Field
	raw            string
	ExtraFields    map[string]apijson.Field
}

func (r *AssistantToolResourcesFileSearch) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r assistantToolResourcesFileSearchJSON) RawJSON() string {
	return r.raw
}

type AssistantDeleted struct {
	ID      string                 `json:"id,required"`
	Deleted bool                   `json:"deleted,required"`
	Object  AssistantDeletedObject `json:"object,required"`
	JSON    assistantDeletedJSON   `json:"-"`
}

// assistantDeletedJSON contains the JSON metadata for the struct
// [AssistantDeleted]
type assistantDeletedJSON struct {
	ID          apijson.Field
	Deleted     apijson.Field
	Object      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *AssistantDeleted) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r assistantDeletedJSON) RawJSON() string {
	return r.raw
}

type AssistantDeletedObject string

const (
	AssistantDeletedObjectAssistantDeleted AssistantDeletedObject = "assistant.deleted"
)

func (r AssistantDeletedObject) IsKnown() bool {
	switch r {
	case AssistantDeletedObjectAssistantDeleted:
		return true
	}
	return false
}

// Represents an event emitted when streaming a Run.
//
// Each event in a server-sent events stream has an `event` and `data` property:
//
// ```
// event: thread.created
// data: {"id": "thread_123", "object": "thread", ...}
// ```
//
// We emit events whenever a new object is created, transitions to a new state, or
// is being streamed in parts (deltas). For example, we emit `thread.run.created`
// when a new run is created, `thread.run.completed` when a run completes, and so
// on. When an Assistant chooses to create a message during a run, we emit a
// `thread.message.created event`, a `thread.message.in_progress` event, many
// `thread.message.delta` events, and finally a `thread.message.completed` event.
//
// We may add additional events over time, so we recommend handling unknown events
// gracefully in your code. See the
// [Assistants API quickstart](https://platform.openai.com/docs/assistants/overview)
// to learn how to integrate the Assistants API with streaming.
type AssistantStreamEvent struct {
	// This field can have the runtime type of [Thread], [Run], [RunStep],
	// [RunStepDeltaEvent], [Message], [MessageDeltaEvent], [shared.ErrorObject].
	Data  interface{}               `json:"data,required"`
	Event AssistantStreamEventEvent `json:"event,required"`
	// Whether to enable input audio transcription.
	Enabled bool                     `json:"enabled"`
	JSON    assistantStreamEventJSON `json:"-"`
	union   AssistantStreamEventUnion
}

// assistantStreamEventJSON contains the JSON metadata for the struct
// [AssistantStreamEvent]
type assistantStreamEventJSON struct {
	Data        apijson.Field
	Event       apijson.Field
	Enabled     apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r assistantStreamEventJSON) RawJSON() string {
	return r.raw
}

func (r *AssistantStreamEvent) UnmarshalJSON(data []byte) (err error) {
	*r = AssistantStreamEvent{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [AssistantStreamEventUnion] interface which you can cast to
// the specific types for more type safety.
//
// Possible runtime types of the union are [AssistantStreamEventThreadCreated],
// [AssistantStreamEventThreadRunCreated], [AssistantStreamEventThreadRunQueued],
// [AssistantStreamEventThreadRunInProgress],
// [AssistantStreamEventThreadRunRequiresAction],
// [AssistantStreamEventThreadRunCompleted],
// [AssistantStreamEventThreadRunIncomplete],
// [AssistantStreamEventThreadRunFailed],
// [AssistantStreamEventThreadRunCancelling],
// [AssistantStreamEventThreadRunCancelled],
// [AssistantStreamEventThreadRunExpired],
// [AssistantStreamEventThreadRunStepCreated],
// [AssistantStreamEventThreadRunStepInProgress],
// [AssistantStreamEventThreadRunStepDelta],
// [AssistantStreamEventThreadRunStepCompleted],
// [AssistantStreamEventThreadRunStepFailed],
// [AssistantStreamEventThreadRunStepCancelled],
// [AssistantStreamEventThreadRunStepExpired],
// [AssistantStreamEventThreadMessageCreated],
// [AssistantStreamEventThreadMessageInProgress],
// [AssistantStreamEventThreadMessageDelta],
// [AssistantStreamEventThreadMessageCompleted],
// [AssistantStreamEventThreadMessageIncomplete], [AssistantStreamEventErrorEvent].
func (r AssistantStreamEvent) AsUnion() AssistantStreamEventUnion {
	return r.union
}

// Represents an event emitted when streaming a Run.
//
// Each event in a server-sent events stream has an `event` and `data` property:
//
// ```
// event: thread.created
// data: {"id": "thread_123", "object": "thread", ...}
// ```
//
// We emit events whenever a new object is created, transitions to a new state, or
// is being streamed in parts (deltas). For example, we emit `thread.run.created`
// when a new run is created, `thread.run.completed` when a run completes, and so
// on. When an Assistant chooses to create a message during a run, we emit a
// `thread.message.created event`, a `thread.message.in_progress` event, many
// `thread.message.delta` events, and finally a `thread.message.completed` event.
//
// We may add additional events over time, so we recommend handling unknown events
// gracefully in your code. See the
// [Assistants API quickstart](https://platform.openai.com/docs/assistants/overview)
// to learn how to integrate the Assistants API with streaming.
//
// Union satisfied by [AssistantStreamEventThreadCreated],
// [AssistantStreamEventThreadRunCreated], [AssistantStreamEventThreadRunQueued],
// [AssistantStreamEventThreadRunInProgress],
// [AssistantStreamEventThreadRunRequiresAction],
// [AssistantStreamEventThreadRunCompleted],
// [AssistantStreamEventThreadRunIncomplete],
// [AssistantStreamEventThreadRunFailed],
// [AssistantStreamEventThreadRunCancelling],
// [AssistantStreamEventThreadRunCancelled],
// [AssistantStreamEventThreadRunExpired],
// [AssistantStreamEventThreadRunStepCreated],
// [AssistantStreamEventThreadRunStepInProgress],
// [AssistantStreamEventThreadRunStepDelta],
// [AssistantStreamEventThreadRunStepCompleted],
// [AssistantStreamEventThreadRunStepFailed],
// [AssistantStreamEventThreadRunStepCancelled],
// [AssistantStreamEventThreadRunStepExpired],
// [AssistantStreamEventThreadMessageCreated],
// [AssistantStreamEventThreadMessageInProgress],
// [AssistantStreamEventThreadMessageDelta],
// [AssistantStreamEventThreadMessageCompleted],
// [AssistantStreamEventThreadMessageIncomplete] or
// [AssistantStreamEventErrorEvent].
type AssistantStreamEventUnion interface {
	implementsAssistantStreamEvent()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*AssistantStreamEventUnion)(nil)).Elem(),
		"event",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(AssistantStreamEventThreadCreated{}),
			DiscriminatorValue: "thread.created",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(AssistantStreamEventThreadRunCreated{}),
			DiscriminatorValue: "thread.run.created",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(AssistantStreamEventThreadRunQueued{}),
			DiscriminatorValue: "thread.run.queued",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(AssistantStreamEventThreadRunInProgress{}),
			DiscriminatorValue: "thread.run.in_progress",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(AssistantStreamEventThreadRunRequiresAction{}),
			DiscriminatorValue: "thread.run.requires_action",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(AssistantStreamEventThreadRunCompleted{}),
			DiscriminatorValue: "thread.run.completed",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(AssistantStreamEventThreadRunIncomplete{}),
			DiscriminatorValue: "thread.run.incomplete",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(AssistantStreamEventThreadRunFailed{}),
			DiscriminatorValue: "thread.run.failed",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(AssistantStreamEventThreadRunCancelling{}),
			DiscriminatorValue: "thread.run.cancelling",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(AssistantStreamEventThreadRunCancelled{}),
			DiscriminatorValue: "thread.run.cancelled",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(AssistantStreamEventThreadRunExpired{}),
			DiscriminatorValue: "thread.run.expired",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(AssistantStreamEventThreadRunStepCreated{}),
			DiscriminatorValue: "thread.run.step.created",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(AssistantStreamEventThreadRunStepInProgress{}),
			DiscriminatorValue: "thread.run.step.in_progress",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(AssistantStreamEventThreadRunStepDelta{}),
			DiscriminatorValue: "thread.run.step.delta",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(AssistantStreamEventThreadRunStepCompleted{}),
			DiscriminatorValue: "thread.run.step.completed",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(AssistantStreamEventThreadRunStepFailed{}),
			DiscriminatorValue: "thread.run.step.failed",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(AssistantStreamEventThreadRunStepCancelled{}),
			DiscriminatorValue: "thread.run.step.cancelled",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(AssistantStreamEventThreadRunStepExpired{}),
			DiscriminatorValue: "thread.run.step.expired",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(AssistantStreamEventThreadMessageCreated{}),
			DiscriminatorValue: "thread.message.created",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(AssistantStreamEventThreadMessageInProgress{}),
			DiscriminatorValue: "thread.message.in_progress",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(AssistantStreamEventThreadMessageDelta{}),
			DiscriminatorValue: "thread.message.delta",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(AssistantStreamEventThreadMessageCompleted{}),
			DiscriminatorValue: "thread.message.completed",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(AssistantStreamEventThreadMessageIncomplete{}),
			DiscriminatorValue: "thread.message.incomplete",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(AssistantStreamEventErrorEvent{}),
			DiscriminatorValue: "error",
		},
	)
}

// Occurs when a new
// [thread](https://platform.openai.com/docs/api-reference/threads/object) is
// created.
type AssistantStreamEventThreadCreated struct {
	// Represents a thread that contains
	// [messages](https://platform.openai.com/docs/api-reference/messages).
	Data  Thread                                 `json:"data,required"`
	Event AssistantStreamEventThreadCreatedEvent `json:"event,required"`
	// Whether to enable input audio transcription.
	Enabled bool                                  `json:"enabled"`
	JSON    assistantStreamEventThreadCreatedJSON `json:"-"`
}

// assistantStreamEventThreadCreatedJSON contains the JSON metadata for the struct
// [AssistantStreamEventThreadCreated]
type assistantStreamEventThreadCreatedJSON struct {
	Data        apijson.Field
	Event       apijson.Field
	Enabled     apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *AssistantStreamEventThreadCreated) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r assistantStreamEventThreadCreatedJSON) RawJSON() string {
	return r.raw
}

func (r AssistantStreamEventThreadCreated) implementsAssistantStreamEvent() {}

type AssistantStreamEventThreadCreatedEvent string

const (
	AssistantStreamEventThreadCreatedEventThreadCreated AssistantStreamEventThreadCreatedEvent = "thread.created"
)

func (r AssistantStreamEventThreadCreatedEvent) IsKnown() bool {
	switch r {
	case AssistantStreamEventThreadCreatedEventThreadCreated:
		return true
	}
	return false
}

// Occurs when a new
// [run](https://platform.openai.com/docs/api-reference/runs/object) is created.
type AssistantStreamEventThreadRunCreated struct {
	// Represents an execution run on a
	// [thread](https://platform.openai.com/docs/api-reference/threads).
	Data  Run                                       `json:"data,required"`
	Event AssistantStreamEventThreadRunCreatedEvent `json:"event,required"`
	JSON  assistantStreamEventThreadRunCreatedJSON  `json:"-"`
}

// assistantStreamEventThreadRunCreatedJSON contains the JSON metadata for the
// struct [AssistantStreamEventThreadRunCreated]
type assistantStreamEventThreadRunCreatedJSON struct {
	Data        apijson.Field
	Event       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *AssistantStreamEventThreadRunCreated) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r assistantStreamEventThreadRunCreatedJSON) RawJSON() string {
	return r.raw
}

func (r AssistantStreamEventThreadRunCreated) implementsAssistantStreamEvent() {}

type AssistantStreamEventThreadRunCreatedEvent string

const (
	AssistantStreamEventThreadRunCreatedEventThreadRunCreated AssistantStreamEventThreadRunCreatedEvent = "thread.run.created"
)

func (r AssistantStreamEventThreadRunCreatedEvent) IsKnown() bool {
	switch r {
	case AssistantStreamEventThreadRunCreatedEventThreadRunCreated:
		return true
	}
	return false
}

// Occurs when a [run](https://platform.openai.com/docs/api-reference/runs/object)
// moves to a `queued` status.
type AssistantStreamEventThreadRunQueued struct {
	// Represents an execution run on a
	// [thread](https://platform.openai.com/docs/api-reference/threads).
	Data  Run                                      `json:"data,required"`
	Event AssistantStreamEventThreadRunQueuedEvent `json:"event,required"`
	JSON  assistantStreamEventThreadRunQueuedJSON  `json:"-"`
}

// assistantStreamEventThreadRunQueuedJSON contains the JSON metadata for the
// struct [AssistantStreamEventThreadRunQueued]
type assistantStreamEventThreadRunQueuedJSON struct {
	Data        apijson.Field
	Event       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *AssistantStreamEventThreadRunQueued) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r assistantStreamEventThreadRunQueuedJSON) RawJSON() string {
	return r.raw
}

func (r AssistantStreamEventThreadRunQueued) implementsAssistantStreamEvent() {}

type AssistantStreamEventThreadRunQueuedEvent string

const (
	AssistantStreamEventThreadRunQueuedEventThreadRunQueued AssistantStreamEventThreadRunQueuedEvent = "thread.run.queued"
)

func (r AssistantStreamEventThreadRunQueuedEvent) IsKnown() bool {
	switch r {
	case AssistantStreamEventThreadRunQueuedEventThreadRunQueued:
		return true
	}
	return false
}

// Occurs when a [run](https://platform.openai.com/docs/api-reference/runs/object)
// moves to an `in_progress` status.
type AssistantStreamEventThreadRunInProgress struct {
	// Represents an execution run on a
	// [thread](https://platform.openai.com/docs/api-reference/threads).
	Data  Run                                          `json:"data,required"`
	Event AssistantStreamEventThreadRunInProgressEvent `json:"event,required"`
	JSON  assistantStreamEventThreadRunInProgressJSON  `json:"-"`
}

// assistantStreamEventThreadRunInProgressJSON contains the JSON metadata for the
// struct [AssistantStreamEventThreadRunInProgress]
type assistantStreamEventThreadRunInProgressJSON struct {
	Data        apijson.Field
	Event       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *AssistantStreamEventThreadRunInProgress) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r assistantStreamEventThreadRunInProgressJSON) RawJSON() string {
	return r.raw
}

func (r AssistantStreamEventThreadRunInProgress) implementsAssistantStreamEvent() {}

type AssistantStreamEventThreadRunInProgressEvent string

const (
	AssistantStreamEventThreadRunInProgressEventThreadRunInProgress AssistantStreamEventThreadRunInProgressEvent = "thread.run.in_progress"
)

func (r AssistantStreamEventThreadRunInProgressEvent) IsKnown() bool {
	switch r {
	case AssistantStreamEventThreadRunInProgressEventThreadRunInProgress:
		return true
	}
	return false
}

// Occurs when a [run](https://platform.openai.com/docs/api-reference/runs/object)
// moves to a `requires_action` status.
type AssistantStreamEventThreadRunRequiresAction struct {
	// Represents an execution run on a
	// [thread](https://platform.openai.com/docs/api-reference/threads).
	Data  Run                                              `json:"data,required"`
	Event AssistantStreamEventThreadRunRequiresActionEvent `json:"event,required"`
	JSON  assistantStreamEventThreadRunRequiresActionJSON  `json:"-"`
}

// assistantStreamEventThreadRunRequiresActionJSON contains the JSON metadata for
// the struct [AssistantStreamEventThreadRunRequiresAction]
type assistantStreamEventThreadRunRequiresActionJSON struct {
	Data        apijson.Field
	Event       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *AssistantStreamEventThreadRunRequiresAction) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r assistantStreamEventThreadRunRequiresActionJSON) RawJSON() string {
	return r.raw
}

func (r AssistantStreamEventThreadRunRequiresAction) implementsAssistantStreamEvent() {}

type AssistantStreamEventThreadRunRequiresActionEvent string

const (
	AssistantStreamEventThreadRunRequiresActionEventThreadRunRequiresAction AssistantStreamEventThreadRunRequiresActionEvent = "thread.run.requires_action"
)

func (r AssistantStreamEventThreadRunRequiresActionEvent) IsKnown() bool {
	switch r {
	case AssistantStreamEventThreadRunRequiresActionEventThreadRunRequiresAction:
		return true
	}
	return false
}

// Occurs when a [run](https://platform.openai.com/docs/api-reference/runs/object)
// is completed.
type AssistantStreamEventThreadRunCompleted struct {
	// Represents an execution run on a
	// [thread](https://platform.openai.com/docs/api-reference/threads).
	Data  Run                                         `json:"data,required"`
	Event AssistantStreamEventThreadRunCompletedEvent `json:"event,required"`
	JSON  assistantStreamEventThreadRunCompletedJSON  `json:"-"`
}

// assistantStreamEventThreadRunCompletedJSON contains the JSON metadata for the
// struct [AssistantStreamEventThreadRunCompleted]
type assistantStreamEventThreadRunCompletedJSON struct {
	Data        apijson.Field
	Event       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *AssistantStreamEventThreadRunCompleted) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r assistantStreamEventThreadRunCompletedJSON) RawJSON() string {
	return r.raw
}

func (r AssistantStreamEventThreadRunCompleted) implementsAssistantStreamEvent() {}

type AssistantStreamEventThreadRunCompletedEvent string

const (
	AssistantStreamEventThreadRunCompletedEventThreadRunCompleted AssistantStreamEventThreadRunCompletedEvent = "thread.run.completed"
)

func (r AssistantStreamEventThreadRunCompletedEvent) IsKnown() bool {
	switch r {
	case AssistantStreamEventThreadRunCompletedEventThreadRunCompleted:
		return true
	}
	return false
}

// Occurs when a [run](https://platform.openai.com/docs/api-reference/runs/object)
// ends with status `incomplete`.
type AssistantStreamEventThreadRunIncomplete struct {
	// Represents an execution run on a
	// [thread](https://platform.openai.com/docs/api-reference/threads).
	Data  Run                                          `json:"data,required"`
	Event AssistantStreamEventThreadRunIncompleteEvent `json:"event,required"`
	JSON  assistantStreamEventThreadRunIncompleteJSON  `json:"-"`
}

// assistantStreamEventThreadRunIncompleteJSON contains the JSON metadata for the
// struct [AssistantStreamEventThreadRunIncomplete]
type assistantStreamEventThreadRunIncompleteJSON struct {
	Data        apijson.Field
	Event       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *AssistantStreamEventThreadRunIncomplete) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r assistantStreamEventThreadRunIncompleteJSON) RawJSON() string {
	return r.raw
}

func (r AssistantStreamEventThreadRunIncomplete) implementsAssistantStreamEvent() {}

type AssistantStreamEventThreadRunIncompleteEvent string

const (
	AssistantStreamEventThreadRunIncompleteEventThreadRunIncomplete AssistantStreamEventThreadRunIncompleteEvent = "thread.run.incomplete"
)

func (r AssistantStreamEventThreadRunIncompleteEvent) IsKnown() bool {
	switch r {
	case AssistantStreamEventThreadRunIncompleteEventThreadRunIncomplete:
		return true
	}
	return false
}

// Occurs when a [run](https://platform.openai.com/docs/api-reference/runs/object)
// fails.
type AssistantStreamEventThreadRunFailed struct {
	// Represents an execution run on a
	// [thread](https://platform.openai.com/docs/api-reference/threads).
	Data  Run                                      `json:"data,required"`
	Event AssistantStreamEventThreadRunFailedEvent `json:"event,required"`
	JSON  assistantStreamEventThreadRunFailedJSON  `json:"-"`
}

// assistantStreamEventThreadRunFailedJSON contains the JSON metadata for the
// struct [AssistantStreamEventThreadRunFailed]
type assistantStreamEventThreadRunFailedJSON struct {
	Data        apijson.Field
	Event       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *AssistantStreamEventThreadRunFailed) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r assistantStreamEventThreadRunFailedJSON) RawJSON() string {
	return r.raw
}

func (r AssistantStreamEventThreadRunFailed) implementsAssistantStreamEvent() {}

type AssistantStreamEventThreadRunFailedEvent string

const (
	AssistantStreamEventThreadRunFailedEventThreadRunFailed AssistantStreamEventThreadRunFailedEvent = "thread.run.failed"
)

func (r AssistantStreamEventThreadRunFailedEvent) IsKnown() bool {
	switch r {
	case AssistantStreamEventThreadRunFailedEventThreadRunFailed:
		return true
	}
	return false
}

// Occurs when a [run](https://platform.openai.com/docs/api-reference/runs/object)
// moves to a `cancelling` status.
type AssistantStreamEventThreadRunCancelling struct {
	// Represents an execution run on a
	// [thread](https://platform.openai.com/docs/api-reference/threads).
	Data  Run                                          `json:"data,required"`
	Event AssistantStreamEventThreadRunCancellingEvent `json:"event,required"`
	JSON  assistantStreamEventThreadRunCancellingJSON  `json:"-"`
}

// assistantStreamEventThreadRunCancellingJSON contains the JSON metadata for the
// struct [AssistantStreamEventThreadRunCancelling]
type assistantStreamEventThreadRunCancellingJSON struct {
	Data        apijson.Field
	Event       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *AssistantStreamEventThreadRunCancelling) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r assistantStreamEventThreadRunCancellingJSON) RawJSON() string {
	return r.raw
}

func (r AssistantStreamEventThreadRunCancelling) implementsAssistantStreamEvent() {}

type AssistantStreamEventThreadRunCancellingEvent string

const (
	AssistantStreamEventThreadRunCancellingEventThreadRunCancelling AssistantStreamEventThreadRunCancellingEvent = "thread.run.cancelling"
)

func (r AssistantStreamEventThreadRunCancellingEvent) IsKnown() bool {
	switch r {
	case AssistantStreamEventThreadRunCancellingEventThreadRunCancelling:
		return true
	}
	return false
}

// Occurs when a [run](https://platform.openai.com/docs/api-reference/runs/object)
// is cancelled.
type AssistantStreamEventThreadRunCancelled struct {
	// Represents an execution run on a
	// [thread](https://platform.openai.com/docs/api-reference/threads).
	Data  Run                                         `json:"data,required"`
	Event AssistantStreamEventThreadRunCancelledEvent `json:"event,required"`
	JSON  assistantStreamEventThreadRunCancelledJSON  `json:"-"`
}

// assistantStreamEventThreadRunCancelledJSON contains the JSON metadata for the
// struct [AssistantStreamEventThreadRunCancelled]
type assistantStreamEventThreadRunCancelledJSON struct {
	Data        apijson.Field
	Event       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *AssistantStreamEventThreadRunCancelled) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r assistantStreamEventThreadRunCancelledJSON) RawJSON() string {
	return r.raw
}

func (r AssistantStreamEventThreadRunCancelled) implementsAssistantStreamEvent() {}

type AssistantStreamEventThreadRunCancelledEvent string

const (
	AssistantStreamEventThreadRunCancelledEventThreadRunCancelled AssistantStreamEventThreadRunCancelledEvent = "thread.run.cancelled"
)

func (r AssistantStreamEventThreadRunCancelledEvent) IsKnown() bool {
	switch r {
	case AssistantStreamEventThreadRunCancelledEventThreadRunCancelled:
		return true
	}
	return false
}

// Occurs when a [run](https://platform.openai.com/docs/api-reference/runs/object)
// expires.
type AssistantStreamEventThreadRunExpired struct {
	// Represents an execution run on a
	// [thread](https://platform.openai.com/docs/api-reference/threads).
	Data  Run                                       `json:"data,required"`
	Event AssistantStreamEventThreadRunExpiredEvent `json:"event,required"`
	JSON  assistantStreamEventThreadRunExpiredJSON  `json:"-"`
}

// assistantStreamEventThreadRunExpiredJSON contains the JSON metadata for the
// struct [AssistantStreamEventThreadRunExpired]
type assistantStreamEventThreadRunExpiredJSON struct {
	Data        apijson.Field
	Event       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *AssistantStreamEventThreadRunExpired) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r assistantStreamEventThreadRunExpiredJSON) RawJSON() string {
	return r.raw
}

func (r AssistantStreamEventThreadRunExpired) implementsAssistantStreamEvent() {}

type AssistantStreamEventThreadRunExpiredEvent string

const (
	AssistantStreamEventThreadRunExpiredEventThreadRunExpired AssistantStreamEventThreadRunExpiredEvent = "thread.run.expired"
)

func (r AssistantStreamEventThreadRunExpiredEvent) IsKnown() bool {
	switch r {
	case AssistantStreamEventThreadRunExpiredEventThreadRunExpired:
		return true
	}
	return false
}

// Occurs when a
// [run step](https://platform.openai.com/docs/api-reference/run-steps/step-object)
// is created.
type AssistantStreamEventThreadRunStepCreated struct {
	// Represents a step in execution of a run.
	Data  RunStep                                       `json:"data,required"`
	Event AssistantStreamEventThreadRunStepCreatedEvent `json:"event,required"`
	JSON  assistantStreamEventThreadRunStepCreatedJSON  `json:"-"`
}

// assistantStreamEventThreadRunStepCreatedJSON contains the JSON metadata for the
// struct [AssistantStreamEventThreadRunStepCreated]
type assistantStreamEventThreadRunStepCreatedJSON struct {
	Data        apijson.Field
	Event       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *AssistantStreamEventThreadRunStepCreated) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r assistantStreamEventThreadRunStepCreatedJSON) RawJSON() string {
	return r.raw
}

func (r AssistantStreamEventThreadRunStepCreated) implementsAssistantStreamEvent() {}

type AssistantStreamEventThreadRunStepCreatedEvent string

const (
	AssistantStreamEventThreadRunStepCreatedEventThreadRunStepCreated AssistantStreamEventThreadRunStepCreatedEvent = "thread.run.step.created"
)

func (r AssistantStreamEventThreadRunStepCreatedEvent) IsKnown() bool {
	switch r {
	case AssistantStreamEventThreadRunStepCreatedEventThreadRunStepCreated:
		return true
	}
	return false
}

// Occurs when a
// [run step](https://platform.openai.com/docs/api-reference/run-steps/step-object)
// moves to an `in_progress` state.
type AssistantStreamEventThreadRunStepInProgress struct {
	// Represents a step in execution of a run.
	Data  RunStep                                          `json:"data,required"`
	Event AssistantStreamEventThreadRunStepInProgressEvent `json:"event,required"`
	JSON  assistantStreamEventThreadRunStepInProgressJSON  `json:"-"`
}

// assistantStreamEventThreadRunStepInProgressJSON contains the JSON metadata for
// the struct [AssistantStreamEventThreadRunStepInProgress]
type assistantStreamEventThreadRunStepInProgressJSON struct {
	Data        apijson.Field
	Event       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *AssistantStreamEventThreadRunStepInProgress) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r assistantStreamEventThreadRunStepInProgressJSON) RawJSON() string {
	return r.raw
}

func (r AssistantStreamEventThreadRunStepInProgress) implementsAssistantStreamEvent() {}

type AssistantStreamEventThreadRunStepInProgressEvent string

const (
	AssistantStreamEventThreadRunStepInProgressEventThreadRunStepInProgress AssistantStreamEventThreadRunStepInProgressEvent = "thread.run.step.in_progress"
)

func (r AssistantStreamEventThreadRunStepInProgressEvent) IsKnown() bool {
	switch r {
	case AssistantStreamEventThreadRunStepInProgressEventThreadRunStepInProgress:
		return true
	}
	return false
}

// Occurs when parts of a
// [run step](https://platform.openai.com/docs/api-reference/run-steps/step-object)
// are being streamed.
type AssistantStreamEventThreadRunStepDelta struct {
	// Represents a run step delta i.e. any changed fields on a run step during
	// streaming.
	Data  RunStepDeltaEvent                           `json:"data,required"`
	Event AssistantStreamEventThreadRunStepDeltaEvent `json:"event,required"`
	JSON  assistantStreamEventThreadRunStepDeltaJSON  `json:"-"`
}

// assistantStreamEventThreadRunStepDeltaJSON contains the JSON metadata for the
// struct [AssistantStreamEventThreadRunStepDelta]
type assistantStreamEventThreadRunStepDeltaJSON struct {
	Data        apijson.Field
	Event       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *AssistantStreamEventThreadRunStepDelta) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r assistantStreamEventThreadRunStepDeltaJSON) RawJSON() string {
	return r.raw
}

func (r AssistantStreamEventThreadRunStepDelta) implementsAssistantStreamEvent() {}

type AssistantStreamEventThreadRunStepDeltaEvent string

const (
	AssistantStreamEventThreadRunStepDeltaEventThreadRunStepDelta AssistantStreamEventThreadRunStepDeltaEvent = "thread.run.step.delta"
)

func (r AssistantStreamEventThreadRunStepDeltaEvent) IsKnown() bool {
	switch r {
	case AssistantStreamEventThreadRunStepDeltaEventThreadRunStepDelta:
		return true
	}
	return false
}

// Occurs when a
// [run step](https://platform.openai.com/docs/api-reference/run-steps/step-object)
// is completed.
type AssistantStreamEventThreadRunStepCompleted struct {
	// Represents a step in execution of a run.
	Data  RunStep                                         `json:"data,required"`
	Event AssistantStreamEventThreadRunStepCompletedEvent `json:"event,required"`
	JSON  assistantStreamEventThreadRunStepCompletedJSON  `json:"-"`
}

// assistantStreamEventThreadRunStepCompletedJSON contains the JSON metadata for
// the struct [AssistantStreamEventThreadRunStepCompleted]
type assistantStreamEventThreadRunStepCompletedJSON struct {
	Data        apijson.Field
	Event       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *AssistantStreamEventThreadRunStepCompleted) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r assistantStreamEventThreadRunStepCompletedJSON) RawJSON() string {
	return r.raw
}

func (r AssistantStreamEventThreadRunStepCompleted) implementsAssistantStreamEvent() {}

type AssistantStreamEventThreadRunStepCompletedEvent string

const (
	AssistantStreamEventThreadRunStepCompletedEventThreadRunStepCompleted AssistantStreamEventThreadRunStepCompletedEvent = "thread.run.step.completed"
)

func (r AssistantStreamEventThreadRunStepCompletedEvent) IsKnown() bool {
	switch r {
	case AssistantStreamEventThreadRunStepCompletedEventThreadRunStepCompleted:
		return true
	}
	return false
}

// Occurs when a
// [run step](https://platform.openai.com/docs/api-reference/run-steps/step-object)
// fails.
type AssistantStreamEventThreadRunStepFailed struct {
	// Represents a step in execution of a run.
	Data  RunStep                                      `json:"data,required"`
	Event AssistantStreamEventThreadRunStepFailedEvent `json:"event,required"`
	JSON  assistantStreamEventThreadRunStepFailedJSON  `json:"-"`
}

// assistantStreamEventThreadRunStepFailedJSON contains the JSON metadata for the
// struct [AssistantStreamEventThreadRunStepFailed]
type assistantStreamEventThreadRunStepFailedJSON struct {
	Data        apijson.Field
	Event       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *AssistantStreamEventThreadRunStepFailed) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r assistantStreamEventThreadRunStepFailedJSON) RawJSON() string {
	return r.raw
}

func (r AssistantStreamEventThreadRunStepFailed) implementsAssistantStreamEvent() {}

type AssistantStreamEventThreadRunStepFailedEvent string

const (
	AssistantStreamEventThreadRunStepFailedEventThreadRunStepFailed AssistantStreamEventThreadRunStepFailedEvent = "thread.run.step.failed"
)

func (r AssistantStreamEventThreadRunStepFailedEvent) IsKnown() bool {
	switch r {
	case AssistantStreamEventThreadRunStepFailedEventThreadRunStepFailed:
		return true
	}
	return false
}

// Occurs when a
// [run step](https://platform.openai.com/docs/api-reference/run-steps/step-object)
// is cancelled.
type AssistantStreamEventThreadRunStepCancelled struct {
	// Represents a step in execution of a run.
	Data  RunStep                                         `json:"data,required"`
	Event AssistantStreamEventThreadRunStepCancelledEvent `json:"event,required"`
	JSON  assistantStreamEventThreadRunStepCancelledJSON  `json:"-"`
}

// assistantStreamEventThreadRunStepCancelledJSON contains the JSON metadata for
// the struct [AssistantStreamEventThreadRunStepCancelled]
type assistantStreamEventThreadRunStepCancelledJSON struct {
	Data        apijson.Field
	Event       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *AssistantStreamEventThreadRunStepCancelled) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r assistantStreamEventThreadRunStepCancelledJSON) RawJSON() string {
	return r.raw
}

func (r AssistantStreamEventThreadRunStepCancelled) implementsAssistantStreamEvent() {}

type AssistantStreamEventThreadRunStepCancelledEvent string

const (
	AssistantStreamEventThreadRunStepCancelledEventThreadRunStepCancelled AssistantStreamEventThreadRunStepCancelledEvent = "thread.run.step.cancelled"
)

func (r AssistantStreamEventThreadRunStepCancelledEvent) IsKnown() bool {
	switch r {
	case AssistantStreamEventThreadRunStepCancelledEventThreadRunStepCancelled:
		return true
	}
	return false
}

// Occurs when a
// [run step](https://platform.openai.com/docs/api-reference/run-steps/step-object)
// expires.
type AssistantStreamEventThreadRunStepExpired struct {
	// Represents a step in execution of a run.
	Data  RunStep                                       `json:"data,required"`
	Event AssistantStreamEventThreadRunStepExpiredEvent `json:"event,required"`
	JSON  assistantStreamEventThreadRunStepExpiredJSON  `json:"-"`
}

// assistantStreamEventThreadRunStepExpiredJSON contains the JSON metadata for the
// struct [AssistantStreamEventThreadRunStepExpired]
type assistantStreamEventThreadRunStepExpiredJSON struct {
	Data        apijson.Field
	Event       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *AssistantStreamEventThreadRunStepExpired) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r assistantStreamEventThreadRunStepExpiredJSON) RawJSON() string {
	return r.raw
}

func (r AssistantStreamEventThreadRunStepExpired) implementsAssistantStreamEvent() {}

type AssistantStreamEventThreadRunStepExpiredEvent string

const (
	AssistantStreamEventThreadRunStepExpiredEventThreadRunStepExpired AssistantStreamEventThreadRunStepExpiredEvent = "thread.run.step.expired"
)

func (r AssistantStreamEventThreadRunStepExpiredEvent) IsKnown() bool {
	switch r {
	case AssistantStreamEventThreadRunStepExpiredEventThreadRunStepExpired:
		return true
	}
	return false
}

// Occurs when a
// [message](https://platform.openai.com/docs/api-reference/messages/object) is
// created.
type AssistantStreamEventThreadMessageCreated struct {
	// Represents a message within a
	// [thread](https://platform.openai.com/docs/api-reference/threads).
	Data  Message                                       `json:"data,required"`
	Event AssistantStreamEventThreadMessageCreatedEvent `json:"event,required"`
	JSON  assistantStreamEventThreadMessageCreatedJSON  `json:"-"`
}

// assistantStreamEventThreadMessageCreatedJSON contains the JSON metadata for the
// struct [AssistantStreamEventThreadMessageCreated]
type assistantStreamEventThreadMessageCreatedJSON struct {
	Data        apijson.Field
	Event       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *AssistantStreamEventThreadMessageCreated) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r assistantStreamEventThreadMessageCreatedJSON) RawJSON() string {
	return r.raw
}

func (r AssistantStreamEventThreadMessageCreated) implementsAssistantStreamEvent() {}

type AssistantStreamEventThreadMessageCreatedEvent string

const (
	AssistantStreamEventThreadMessageCreatedEventThreadMessageCreated AssistantStreamEventThreadMessageCreatedEvent = "thread.message.created"
)

func (r AssistantStreamEventThreadMessageCreatedEvent) IsKnown() bool {
	switch r {
	case AssistantStreamEventThreadMessageCreatedEventThreadMessageCreated:
		return true
	}
	return false
}

// Occurs when a
// [message](https://platform.openai.com/docs/api-reference/messages/object) moves
// to an `in_progress` state.
type AssistantStreamEventThreadMessageInProgress struct {
	// Represents a message within a
	// [thread](https://platform.openai.com/docs/api-reference/threads).
	Data  Message                                          `json:"data,required"`
	Event AssistantStreamEventThreadMessageInProgressEvent `json:"event,required"`
	JSON  assistantStreamEventThreadMessageInProgressJSON  `json:"-"`
}

// assistantStreamEventThreadMessageInProgressJSON contains the JSON metadata for
// the struct [AssistantStreamEventThreadMessageInProgress]
type assistantStreamEventThreadMessageInProgressJSON struct {
	Data        apijson.Field
	Event       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *AssistantStreamEventThreadMessageInProgress) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r assistantStreamEventThreadMessageInProgressJSON) RawJSON() string {
	return r.raw
}

func (r AssistantStreamEventThreadMessageInProgress) implementsAssistantStreamEvent() {}

type AssistantStreamEventThreadMessageInProgressEvent string

const (
	AssistantStreamEventThreadMessageInProgressEventThreadMessageInProgress AssistantStreamEventThreadMessageInProgressEvent = "thread.message.in_progress"
)

func (r AssistantStreamEventThreadMessageInProgressEvent) IsKnown() bool {
	switch r {
	case AssistantStreamEventThreadMessageInProgressEventThreadMessageInProgress:
		return true
	}
	return false
}

// Occurs when parts of a
// [Message](https://platform.openai.com/docs/api-reference/messages/object) are
// being streamed.
type AssistantStreamEventThreadMessageDelta struct {
	// Represents a message delta i.e. any changed fields on a message during
	// streaming.
	Data  MessageDeltaEvent                           `json:"data,required"`
	Event AssistantStreamEventThreadMessageDeltaEvent `json:"event,required"`
	JSON  assistantStreamEventThreadMessageDeltaJSON  `json:"-"`
}

// assistantStreamEventThreadMessageDeltaJSON contains the JSON metadata for the
// struct [AssistantStreamEventThreadMessageDelta]
type assistantStreamEventThreadMessageDeltaJSON struct {
	Data        apijson.Field
	Event       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *AssistantStreamEventThreadMessageDelta) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r assistantStreamEventThreadMessageDeltaJSON) RawJSON() string {
	return r.raw
}

func (r AssistantStreamEventThreadMessageDelta) implementsAssistantStreamEvent() {}

type AssistantStreamEventThreadMessageDeltaEvent string

const (
	AssistantStreamEventThreadMessageDeltaEventThreadMessageDelta AssistantStreamEventThreadMessageDeltaEvent = "thread.message.delta"
)

func (r AssistantStreamEventThreadMessageDeltaEvent) IsKnown() bool {
	switch r {
	case AssistantStreamEventThreadMessageDeltaEventThreadMessageDelta:
		return true
	}
	return false
}

// Occurs when a
// [message](https://platform.openai.com/docs/api-reference/messages/object) is
// completed.
type AssistantStreamEventThreadMessageCompleted struct {
	// Represents a message within a
	// [thread](https://platform.openai.com/docs/api-reference/threads).
	Data  Message                                         `json:"data,required"`
	Event AssistantStreamEventThreadMessageCompletedEvent `json:"event,required"`
	JSON  assistantStreamEventThreadMessageCompletedJSON  `json:"-"`
}

// assistantStreamEventThreadMessageCompletedJSON contains the JSON metadata for
// the struct [AssistantStreamEventThreadMessageCompleted]
type assistantStreamEventThreadMessageCompletedJSON struct {
	Data        apijson.Field
	Event       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *AssistantStreamEventThreadMessageCompleted) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r assistantStreamEventThreadMessageCompletedJSON) RawJSON() string {
	return r.raw
}

func (r AssistantStreamEventThreadMessageCompleted) implementsAssistantStreamEvent() {}

type AssistantStreamEventThreadMessageCompletedEvent string

const (
	AssistantStreamEventThreadMessageCompletedEventThreadMessageCompleted AssistantStreamEventThreadMessageCompletedEvent = "thread.message.completed"
)

func (r AssistantStreamEventThreadMessageCompletedEvent) IsKnown() bool {
	switch r {
	case AssistantStreamEventThreadMessageCompletedEventThreadMessageCompleted:
		return true
	}
	return false
}

// Occurs when a
// [message](https://platform.openai.com/docs/api-reference/messages/object) ends
// before it is completed.
type AssistantStreamEventThreadMessageIncomplete struct {
	// Represents a message within a
	// [thread](https://platform.openai.com/docs/api-reference/threads).
	Data  Message                                          `json:"data,required"`
	Event AssistantStreamEventThreadMessageIncompleteEvent `json:"event,required"`
	JSON  assistantStreamEventThreadMessageIncompleteJSON  `json:"-"`
}

// assistantStreamEventThreadMessageIncompleteJSON contains the JSON metadata for
// the struct [AssistantStreamEventThreadMessageIncomplete]
type assistantStreamEventThreadMessageIncompleteJSON struct {
	Data        apijson.Field
	Event       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *AssistantStreamEventThreadMessageIncomplete) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r assistantStreamEventThreadMessageIncompleteJSON) RawJSON() string {
	return r.raw
}

func (r AssistantStreamEventThreadMessageIncomplete) implementsAssistantStreamEvent() {}

type AssistantStreamEventThreadMessageIncompleteEvent string

const (
	AssistantStreamEventThreadMessageIncompleteEventThreadMessageIncomplete AssistantStreamEventThreadMessageIncompleteEvent = "thread.message.incomplete"
)

func (r AssistantStreamEventThreadMessageIncompleteEvent) IsKnown() bool {
	switch r {
	case AssistantStreamEventThreadMessageIncompleteEventThreadMessageIncomplete:
		return true
	}
	return false
}

// Occurs when an
// [error](https://platform.openai.com/docs/guides/error-codes#api-errors) occurs.
// This can happen due to an internal server error or a timeout.
type AssistantStreamEventErrorEvent struct {
	Data  shared.ErrorObject                  `json:"data,required"`
	Event AssistantStreamEventErrorEventEvent `json:"event,required"`
	JSON  assistantStreamEventErrorEventJSON  `json:"-"`
}

// assistantStreamEventErrorEventJSON contains the JSON metadata for the struct
// [AssistantStreamEventErrorEvent]
type assistantStreamEventErrorEventJSON struct {
	Data        apijson.Field
	Event       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *AssistantStreamEventErrorEvent) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r assistantStreamEventErrorEventJSON) RawJSON() string {
	return r.raw
}

func (r AssistantStreamEventErrorEvent) implementsAssistantStreamEvent() {}

type AssistantStreamEventErrorEventEvent string

const (
	AssistantStreamEventErrorEventEventError AssistantStreamEventErrorEventEvent = "error"
)

func (r AssistantStreamEventErrorEventEvent) IsKnown() bool {
	switch r {
	case AssistantStreamEventErrorEventEventError:
		return true
	}
	return false
}

type AssistantStreamEventEvent string

const (
	AssistantStreamEventEventThreadCreated           AssistantStreamEventEvent = "thread.created"
	AssistantStreamEventEventThreadRunCreated        AssistantStreamEventEvent = "thread.run.created"
	AssistantStreamEventEventThreadRunQueued         AssistantStreamEventEvent = "thread.run.queued"
	AssistantStreamEventEventThreadRunInProgress     AssistantStreamEventEvent = "thread.run.in_progress"
	AssistantStreamEventEventThreadRunRequiresAction AssistantStreamEventEvent = "thread.run.requires_action"
	AssistantStreamEventEventThreadRunCompleted      AssistantStreamEventEvent = "thread.run.completed"
	AssistantStreamEventEventThreadRunIncomplete     AssistantStreamEventEvent = "thread.run.incomplete"
	AssistantStreamEventEventThreadRunFailed         AssistantStreamEventEvent = "thread.run.failed"
	AssistantStreamEventEventThreadRunCancelling     AssistantStreamEventEvent = "thread.run.cancelling"
	AssistantStreamEventEventThreadRunCancelled      AssistantStreamEventEvent = "thread.run.cancelled"
	AssistantStreamEventEventThreadRunExpired        AssistantStreamEventEvent = "thread.run.expired"
	AssistantStreamEventEventThreadRunStepCreated    AssistantStreamEventEvent = "thread.run.step.created"
	AssistantStreamEventEventThreadRunStepInProgress AssistantStreamEventEvent = "thread.run.step.in_progress"
	AssistantStreamEventEventThreadRunStepDelta      AssistantStreamEventEvent = "thread.run.step.delta"
	AssistantStreamEventEventThreadRunStepCompleted  AssistantStreamEventEvent = "thread.run.step.completed"
	AssistantStreamEventEventThreadRunStepFailed     AssistantStreamEventEvent = "thread.run.step.failed"
	AssistantStreamEventEventThreadRunStepCancelled  AssistantStreamEventEvent = "thread.run.step.cancelled"
	AssistantStreamEventEventThreadRunStepExpired    AssistantStreamEventEvent = "thread.run.step.expired"
	AssistantStreamEventEventThreadMessageCreated    AssistantStreamEventEvent = "thread.message.created"
	AssistantStreamEventEventThreadMessageInProgress AssistantStreamEventEvent = "thread.message.in_progress"
	AssistantStreamEventEventThreadMessageDelta      AssistantStreamEventEvent = "thread.message.delta"
	AssistantStreamEventEventThreadMessageCompleted  AssistantStreamEventEvent = "thread.message.completed"
	AssistantStreamEventEventThreadMessageIncomplete AssistantStreamEventEvent = "thread.message.incomplete"
	AssistantStreamEventEventError                   AssistantStreamEventEvent = "error"
)

func (r AssistantStreamEventEvent) IsKnown() bool {
	switch r {
	case AssistantStreamEventEventThreadCreated, AssistantStreamEventEventThreadRunCreated, AssistantStreamEventEventThreadRunQueued, AssistantStreamEventEventThreadRunInProgress, AssistantStreamEventEventThreadRunRequiresAction, AssistantStreamEventEventThreadRunCompleted, AssistantStreamEventEventThreadRunIncomplete, AssistantStreamEventEventThreadRunFailed, AssistantStreamEventEventThreadRunCancelling, AssistantStreamEventEventThreadRunCancelled, AssistantStreamEventEventThreadRunExpired, AssistantStreamEventEventThreadRunStepCreated, AssistantStreamEventEventThreadRunStepInProgress, AssistantStreamEventEventThreadRunStepDelta, AssistantStreamEventEventThreadRunStepCompleted, AssistantStreamEventEventThreadRunStepFailed, AssistantStreamEventEventThreadRunStepCancelled, AssistantStreamEventEventThreadRunStepExpired, AssistantStreamEventEventThreadMessageCreated, AssistantStreamEventEventThreadMessageInProgress, AssistantStreamEventEventThreadMessageDelta, AssistantStreamEventEventThreadMessageCompleted, AssistantStreamEventEventThreadMessageIncomplete, AssistantStreamEventEventError:
		return true
	}
	return false
}

type AssistantTool struct {
	// The type of tool being defined: `code_interpreter`
	Type AssistantToolType `json:"type,required"`
	// This field can have the runtime type of [FileSearchToolFileSearch].
	FileSearch interface{}               `json:"file_search"`
	Function   shared.FunctionDefinition `json:"function"`
	JSON       assistantToolJSON         `json:"-"`
	union      AssistantToolUnion
}

// assistantToolJSON contains the JSON metadata for the struct [AssistantTool]
type assistantToolJSON struct {
	Type        apijson.Field
	FileSearch  apijson.Field
	Function    apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r assistantToolJSON) RawJSON() string {
	return r.raw
}

func (r *AssistantTool) UnmarshalJSON(data []byte) (err error) {
	*r = AssistantTool{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [AssistantToolUnion] interface which you can cast to the
// specific types for more type safety.
//
// Possible runtime types of the union are [CodeInterpreterTool], [FileSearchTool],
// [FunctionTool].
func (r AssistantTool) AsUnion() AssistantToolUnion {
	return r.union
}

// Union satisfied by [CodeInterpreterTool], [FileSearchTool] or [FunctionTool].
type AssistantToolUnion interface {
	implementsAssistantTool()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*AssistantToolUnion)(nil)).Elem(),
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(CodeInterpreterTool{}),
			DiscriminatorValue: "code_interpreter",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(FileSearchTool{}),
			DiscriminatorValue: "file_search",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(FunctionTool{}),
			DiscriminatorValue: "function",
		},
	)
}

// The type of tool being defined: `code_interpreter`
type AssistantToolType string

const (
	AssistantToolTypeCodeInterpreter AssistantToolType = "code_interpreter"
	AssistantToolTypeFileSearch      AssistantToolType = "file_search"
	AssistantToolTypeFunction        AssistantToolType = "function"
)

func (r AssistantToolType) IsKnown() bool {
	switch r {
	case AssistantToolTypeCodeInterpreter, AssistantToolTypeFileSearch, AssistantToolTypeFunction:
		return true
	}
	return false
}

type AssistantToolParam struct {
	// The type of tool being defined: `code_interpreter`
	Type       param.Field[AssistantToolType]              `json:"type,required"`
	FileSearch param.Field[interface{}]                    `json:"file_search"`
	Function   param.Field[shared.FunctionDefinitionParam] `json:"function"`
}

func (r AssistantToolParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r AssistantToolParam) implementsAssistantToolUnionParam() {}

// Satisfied by [CodeInterpreterToolParam], [FileSearchToolParam],
// [FunctionToolParam], [AssistantToolParam].
type AssistantToolUnionParam interface {
	implementsAssistantToolUnionParam()
}

type CodeInterpreterTool struct {
	// The type of tool being defined: `code_interpreter`
	Type CodeInterpreterToolType `json:"type,required"`
	JSON codeInterpreterToolJSON `json:"-"`
}

// codeInterpreterToolJSON contains the JSON metadata for the struct
// [CodeInterpreterTool]
type codeInterpreterToolJSON struct {
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CodeInterpreterTool) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r codeInterpreterToolJSON) RawJSON() string {
	return r.raw
}

func (r CodeInterpreterTool) implementsAssistantTool() {}

func (r CodeInterpreterTool) implementsMessageAttachmentsTool() {}

// The type of tool being defined: `code_interpreter`
type CodeInterpreterToolType string

const (
	CodeInterpreterToolTypeCodeInterpreter CodeInterpreterToolType = "code_interpreter"
)

func (r CodeInterpreterToolType) IsKnown() bool {
	switch r {
	case CodeInterpreterToolTypeCodeInterpreter:
		return true
	}
	return false
}

type CodeInterpreterToolParam struct {
	// The type of tool being defined: `code_interpreter`
	Type param.Field[CodeInterpreterToolType] `json:"type,required"`
}

func (r CodeInterpreterToolParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r CodeInterpreterToolParam) implementsAssistantToolUnionParam() {}

func (r CodeInterpreterToolParam) implementsBetaThreadNewParamsMessagesAttachmentsToolUnion() {}

func (r CodeInterpreterToolParam) implementsBetaThreadNewAndRunParamsThreadMessagesAttachmentsToolUnion() {
}

func (r CodeInterpreterToolParam) implementsBetaThreadNewAndRunParamsToolUnion() {}

func (r CodeInterpreterToolParam) implementsBetaThreadRunNewParamsAdditionalMessagesAttachmentsToolUnion() {
}

func (r CodeInterpreterToolParam) implementsBetaThreadMessageNewParamsAttachmentsToolUnion() {}

type FileSearchTool struct {
	// The type of tool being defined: `file_search`
	Type FileSearchToolType `json:"type,required"`
	// Overrides for the file search tool.
	FileSearch FileSearchToolFileSearch `json:"file_search"`
	JSON       fileSearchToolJSON       `json:"-"`
}

// fileSearchToolJSON contains the JSON metadata for the struct [FileSearchTool]
type fileSearchToolJSON struct {
	Type        apijson.Field
	FileSearch  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *FileSearchTool) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r fileSearchToolJSON) RawJSON() string {
	return r.raw
}

func (r FileSearchTool) implementsAssistantTool() {}

// The type of tool being defined: `file_search`
type FileSearchToolType string

const (
	FileSearchToolTypeFileSearch FileSearchToolType = "file_search"
)

func (r FileSearchToolType) IsKnown() bool {
	switch r {
	case FileSearchToolTypeFileSearch:
		return true
	}
	return false
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
	MaxNumResults int64 `json:"max_num_results"`
	// The ranking options for the file search. If not specified, the file search tool
	// will use the `auto` ranker and a score_threshold of 0.
	//
	// See the
	// [file search tool documentation](https://platform.openai.com/docs/assistants/tools/file-search#customizing-file-search-settings)
	// for more information.
	RankingOptions FileSearchToolFileSearchRankingOptions `json:"ranking_options"`
	JSON           fileSearchToolFileSearchJSON           `json:"-"`
}

// fileSearchToolFileSearchJSON contains the JSON metadata for the struct
// [FileSearchToolFileSearch]
type fileSearchToolFileSearchJSON struct {
	MaxNumResults  apijson.Field
	RankingOptions apijson.Field
	raw            string
	ExtraFields    map[string]apijson.Field
}

func (r *FileSearchToolFileSearch) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r fileSearchToolFileSearchJSON) RawJSON() string {
	return r.raw
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
	ScoreThreshold float64 `json:"score_threshold,required"`
	// The ranker to use for the file search. If not specified will use the `auto`
	// ranker.
	Ranker FileSearchToolFileSearchRankingOptionsRanker `json:"ranker"`
	JSON   fileSearchToolFileSearchRankingOptionsJSON   `json:"-"`
}

// fileSearchToolFileSearchRankingOptionsJSON contains the JSON metadata for the
// struct [FileSearchToolFileSearchRankingOptions]
type fileSearchToolFileSearchRankingOptionsJSON struct {
	ScoreThreshold apijson.Field
	Ranker         apijson.Field
	raw            string
	ExtraFields    map[string]apijson.Field
}

func (r *FileSearchToolFileSearchRankingOptions) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r fileSearchToolFileSearchRankingOptionsJSON) RawJSON() string {
	return r.raw
}

// The ranker to use for the file search. If not specified will use the `auto`
// ranker.
type FileSearchToolFileSearchRankingOptionsRanker string

const (
	FileSearchToolFileSearchRankingOptionsRankerAuto              FileSearchToolFileSearchRankingOptionsRanker = "auto"
	FileSearchToolFileSearchRankingOptionsRankerDefault2024_08_21 FileSearchToolFileSearchRankingOptionsRanker = "default_2024_08_21"
)

func (r FileSearchToolFileSearchRankingOptionsRanker) IsKnown() bool {
	switch r {
	case FileSearchToolFileSearchRankingOptionsRankerAuto, FileSearchToolFileSearchRankingOptionsRankerDefault2024_08_21:
		return true
	}
	return false
}

type FileSearchToolParam struct {
	// The type of tool being defined: `file_search`
	Type param.Field[FileSearchToolType] `json:"type,required"`
	// Overrides for the file search tool.
	FileSearch param.Field[FileSearchToolFileSearchParam] `json:"file_search"`
}

func (r FileSearchToolParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r FileSearchToolParam) implementsAssistantToolUnionParam() {}

func (r FileSearchToolParam) implementsBetaThreadNewAndRunParamsToolUnion() {}

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
	MaxNumResults param.Field[int64] `json:"max_num_results"`
	// The ranking options for the file search. If not specified, the file search tool
	// will use the `auto` ranker and a score_threshold of 0.
	//
	// See the
	// [file search tool documentation](https://platform.openai.com/docs/assistants/tools/file-search#customizing-file-search-settings)
	// for more information.
	RankingOptions param.Field[FileSearchToolFileSearchRankingOptionsParam] `json:"ranking_options"`
}

func (r FileSearchToolFileSearchParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
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
	ScoreThreshold param.Field[float64] `json:"score_threshold,required"`
	// The ranker to use for the file search. If not specified will use the `auto`
	// ranker.
	Ranker param.Field[FileSearchToolFileSearchRankingOptionsRanker] `json:"ranker"`
}

func (r FileSearchToolFileSearchRankingOptionsParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type FunctionTool struct {
	Function shared.FunctionDefinition `json:"function,required"`
	// The type of tool being defined: `function`
	Type FunctionToolType `json:"type,required"`
	JSON functionToolJSON `json:"-"`
}

// functionToolJSON contains the JSON metadata for the struct [FunctionTool]
type functionToolJSON struct {
	Function    apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *FunctionTool) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r functionToolJSON) RawJSON() string {
	return r.raw
}

func (r FunctionTool) implementsAssistantTool() {}

// The type of tool being defined: `function`
type FunctionToolType string

const (
	FunctionToolTypeFunction FunctionToolType = "function"
)

func (r FunctionToolType) IsKnown() bool {
	switch r {
	case FunctionToolTypeFunction:
		return true
	}
	return false
}

type FunctionToolParam struct {
	Function param.Field[shared.FunctionDefinitionParam] `json:"function,required"`
	// The type of tool being defined: `function`
	Type param.Field[FunctionToolType] `json:"type,required"`
}

func (r FunctionToolParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r FunctionToolParam) implementsAssistantToolUnionParam() {}

func (r FunctionToolParam) implementsBetaThreadNewAndRunParamsToolUnion() {}

type BetaAssistantNewParams struct {
	// ID of the model to use. You can use the
	// [List models](https://platform.openai.com/docs/api-reference/models/list) API to
	// see all of your available models, or see our
	// [Model overview](https://platform.openai.com/docs/models) for descriptions of
	// them.
	Model param.Field[shared.ChatModel] `json:"model,required"`
	// The description of the assistant. The maximum length is 512 characters.
	Description param.Field[string] `json:"description"`
	// The system instructions that the assistant uses. The maximum length is 256,000
	// characters.
	Instructions param.Field[string] `json:"instructions"`
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard.
	//
	// Keys are strings with a maximum length of 64 characters. Values are strings with
	// a maximum length of 512 characters.
	Metadata param.Field[shared.MetadataParam] `json:"metadata"`
	// The name of the assistant. The maximum length is 256 characters.
	Name param.Field[string] `json:"name"`
	// **o1 and o3-mini models only**
	//
	// Constrains effort on reasoning for
	// [reasoning models](https://platform.openai.com/docs/guides/reasoning). Currently
	// supported values are `low`, `medium`, and `high`. Reducing reasoning effort can
	// result in faster responses and fewer tokens used on reasoning in a response.
	ReasoningEffort param.Field[BetaAssistantNewParamsReasoningEffort] `json:"reasoning_effort"`
	// What sampling temperature to use, between 0 and 2. Higher values like 0.8 will
	// make the output more random, while lower values like 0.2 will make it more
	// focused and deterministic.
	Temperature param.Field[float64] `json:"temperature"`
	// A set of resources that are used by the assistant's tools. The resources are
	// specific to the type of tool. For example, the `code_interpreter` tool requires
	// a list of file IDs, while the `file_search` tool requires a list of vector store
	// IDs.
	ToolResources param.Field[BetaAssistantNewParamsToolResources] `json:"tool_resources"`
	// A list of tool enabled on the assistant. There can be a maximum of 128 tools per
	// assistant. Tools can be of types `code_interpreter`, `file_search`, or
	// `function`.
	Tools param.Field[[]AssistantToolUnionParam] `json:"tools"`
	// An alternative to sampling with temperature, called nucleus sampling, where the
	// model considers the results of the tokens with top_p probability mass. So 0.1
	// means only the tokens comprising the top 10% probability mass are considered.
	//
	// We generally recommend altering this or temperature but not both.
	TopP param.Field[float64] `json:"top_p"`
}

func (r BetaAssistantNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
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

func (r BetaAssistantNewParamsReasoningEffort) IsKnown() bool {
	switch r {
	case BetaAssistantNewParamsReasoningEffortLow, BetaAssistantNewParamsReasoningEffortMedium, BetaAssistantNewParamsReasoningEffortHigh:
		return true
	}
	return false
}

// A set of resources that are used by the assistant's tools. The resources are
// specific to the type of tool. For example, the `code_interpreter` tool requires
// a list of file IDs, while the `file_search` tool requires a list of vector store
// IDs.
type BetaAssistantNewParamsToolResources struct {
	CodeInterpreter param.Field[BetaAssistantNewParamsToolResourcesCodeInterpreter] `json:"code_interpreter"`
	FileSearch      param.Field[BetaAssistantNewParamsToolResourcesFileSearch]      `json:"file_search"`
}

func (r BetaAssistantNewParamsToolResources) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type BetaAssistantNewParamsToolResourcesCodeInterpreter struct {
	// A list of [file](https://platform.openai.com/docs/api-reference/files) IDs made
	// available to the `code_interpreter` tool. There can be a maximum of 20 files
	// associated with the tool.
	FileIDs param.Field[[]string] `json:"file_ids"`
}

func (r BetaAssistantNewParamsToolResourcesCodeInterpreter) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type BetaAssistantNewParamsToolResourcesFileSearch struct {
	// The
	// [vector store](https://platform.openai.com/docs/api-reference/vector-stores/object)
	// attached to this assistant. There can be a maximum of 1 vector store attached to
	// the assistant.
	VectorStoreIDs param.Field[[]string] `json:"vector_store_ids"`
	// A helper to create a
	// [vector store](https://platform.openai.com/docs/api-reference/vector-stores/object)
	// with file_ids and attach it to this assistant. There can be a maximum of 1
	// vector store attached to the assistant.
	VectorStores param.Field[[]BetaAssistantNewParamsToolResourcesFileSearchVectorStore] `json:"vector_stores"`
}

func (r BetaAssistantNewParamsToolResourcesFileSearch) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type BetaAssistantNewParamsToolResourcesFileSearchVectorStore struct {
	// The chunking strategy used to chunk the file(s). If not set, will use the `auto`
	// strategy. Only applicable if `file_ids` is non-empty.
	ChunkingStrategy param.Field[FileChunkingStrategyParamUnion] `json:"chunking_strategy"`
	// A list of [file](https://platform.openai.com/docs/api-reference/files) IDs to
	// add to the vector store. There can be a maximum of 10000 files in a vector
	// store.
	FileIDs param.Field[[]string] `json:"file_ids"`
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard.
	//
	// Keys are strings with a maximum length of 64 characters. Values are strings with
	// a maximum length of 512 characters.
	Metadata param.Field[shared.MetadataParam] `json:"metadata"`
}

func (r BetaAssistantNewParamsToolResourcesFileSearchVectorStore) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type BetaAssistantUpdateParams struct {
	// The description of the assistant. The maximum length is 512 characters.
	Description param.Field[string] `json:"description"`
	// The system instructions that the assistant uses. The maximum length is 256,000
	// characters.
	Instructions param.Field[string] `json:"instructions"`
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard.
	//
	// Keys are strings with a maximum length of 64 characters. Values are strings with
	// a maximum length of 512 characters.
	Metadata param.Field[shared.MetadataParam] `json:"metadata"`
	// ID of the model to use. You can use the
	// [List models](https://platform.openai.com/docs/api-reference/models/list) API to
	// see all of your available models, or see our
	// [Model overview](https://platform.openai.com/docs/models) for descriptions of
	// them.
	Model param.Field[BetaAssistantUpdateParamsModel] `json:"model"`
	// The name of the assistant. The maximum length is 256 characters.
	Name param.Field[string] `json:"name"`
	// **o1 and o3-mini models only**
	//
	// Constrains effort on reasoning for
	// [reasoning models](https://platform.openai.com/docs/guides/reasoning). Currently
	// supported values are `low`, `medium`, and `high`. Reducing reasoning effort can
	// result in faster responses and fewer tokens used on reasoning in a response.
	ReasoningEffort param.Field[BetaAssistantUpdateParamsReasoningEffort] `json:"reasoning_effort"`
	// What sampling temperature to use, between 0 and 2. Higher values like 0.8 will
	// make the output more random, while lower values like 0.2 will make it more
	// focused and deterministic.
	Temperature param.Field[float64] `json:"temperature"`
	// A set of resources that are used by the assistant's tools. The resources are
	// specific to the type of tool. For example, the `code_interpreter` tool requires
	// a list of file IDs, while the `file_search` tool requires a list of vector store
	// IDs.
	ToolResources param.Field[BetaAssistantUpdateParamsToolResources] `json:"tool_resources"`
	// A list of tool enabled on the assistant. There can be a maximum of 128 tools per
	// assistant. Tools can be of types `code_interpreter`, `file_search`, or
	// `function`.
	Tools param.Field[[]AssistantToolUnionParam] `json:"tools"`
	// An alternative to sampling with temperature, called nucleus sampling, where the
	// model considers the results of the tokens with top_p probability mass. So 0.1
	// means only the tokens comprising the top 10% probability mass are considered.
	//
	// We generally recommend altering this or temperature but not both.
	TopP param.Field[float64] `json:"top_p"`
}

func (r BetaAssistantUpdateParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// ID of the model to use. You can use the
// [List models](https://platform.openai.com/docs/api-reference/models/list) API to
// see all of your available models, or see our
// [Model overview](https://platform.openai.com/docs/models) for descriptions of
// them.
type BetaAssistantUpdateParamsModel string

const (
	BetaAssistantUpdateParamsModelO3Mini                  BetaAssistantUpdateParamsModel = "o3-mini"
	BetaAssistantUpdateParamsModelO3Mini2025_01_31        BetaAssistantUpdateParamsModel = "o3-mini-2025-01-31"
	BetaAssistantUpdateParamsModelO1                      BetaAssistantUpdateParamsModel = "o1"
	BetaAssistantUpdateParamsModelO1_2024_12_17           BetaAssistantUpdateParamsModel = "o1-2024-12-17"
	BetaAssistantUpdateParamsModelGPT4o                   BetaAssistantUpdateParamsModel = "gpt-4o"
	BetaAssistantUpdateParamsModelGPT4o2024_11_20         BetaAssistantUpdateParamsModel = "gpt-4o-2024-11-20"
	BetaAssistantUpdateParamsModelGPT4o2024_08_06         BetaAssistantUpdateParamsModel = "gpt-4o-2024-08-06"
	BetaAssistantUpdateParamsModelGPT4o2024_05_13         BetaAssistantUpdateParamsModel = "gpt-4o-2024-05-13"
	BetaAssistantUpdateParamsModelGPT4oMini               BetaAssistantUpdateParamsModel = "gpt-4o-mini"
	BetaAssistantUpdateParamsModelGPT4oMini2024_07_18     BetaAssistantUpdateParamsModel = "gpt-4o-mini-2024-07-18"
	BetaAssistantUpdateParamsModelGPT4_5Preview           BetaAssistantUpdateParamsModel = "gpt-4.5-preview"
	BetaAssistantUpdateParamsModelGPT4_5Preview2025_02_27 BetaAssistantUpdateParamsModel = "gpt-4.5-preview-2025-02-27"
	BetaAssistantUpdateParamsModelGPT4Turbo               BetaAssistantUpdateParamsModel = "gpt-4-turbo"
	BetaAssistantUpdateParamsModelGPT4Turbo2024_04_09     BetaAssistantUpdateParamsModel = "gpt-4-turbo-2024-04-09"
	BetaAssistantUpdateParamsModelGPT4_0125Preview        BetaAssistantUpdateParamsModel = "gpt-4-0125-preview"
	BetaAssistantUpdateParamsModelGPT4TurboPreview        BetaAssistantUpdateParamsModel = "gpt-4-turbo-preview"
	BetaAssistantUpdateParamsModelGPT4_1106Preview        BetaAssistantUpdateParamsModel = "gpt-4-1106-preview"
	BetaAssistantUpdateParamsModelGPT4VisionPreview       BetaAssistantUpdateParamsModel = "gpt-4-vision-preview"
	BetaAssistantUpdateParamsModelGPT4                    BetaAssistantUpdateParamsModel = "gpt-4"
	BetaAssistantUpdateParamsModelGPT4_0314               BetaAssistantUpdateParamsModel = "gpt-4-0314"
	BetaAssistantUpdateParamsModelGPT4_0613               BetaAssistantUpdateParamsModel = "gpt-4-0613"
	BetaAssistantUpdateParamsModelGPT4_32k                BetaAssistantUpdateParamsModel = "gpt-4-32k"
	BetaAssistantUpdateParamsModelGPT4_32k0314            BetaAssistantUpdateParamsModel = "gpt-4-32k-0314"
	BetaAssistantUpdateParamsModelGPT4_32k0613            BetaAssistantUpdateParamsModel = "gpt-4-32k-0613"
	BetaAssistantUpdateParamsModelGPT3_5Turbo             BetaAssistantUpdateParamsModel = "gpt-3.5-turbo"
	BetaAssistantUpdateParamsModelGPT3_5Turbo16k          BetaAssistantUpdateParamsModel = "gpt-3.5-turbo-16k"
	BetaAssistantUpdateParamsModelGPT3_5Turbo0613         BetaAssistantUpdateParamsModel = "gpt-3.5-turbo-0613"
	BetaAssistantUpdateParamsModelGPT3_5Turbo1106         BetaAssistantUpdateParamsModel = "gpt-3.5-turbo-1106"
	BetaAssistantUpdateParamsModelGPT3_5Turbo0125         BetaAssistantUpdateParamsModel = "gpt-3.5-turbo-0125"
	BetaAssistantUpdateParamsModelGPT3_5Turbo16k0613      BetaAssistantUpdateParamsModel = "gpt-3.5-turbo-16k-0613"
)

func (r BetaAssistantUpdateParamsModel) IsKnown() bool {
	switch r {
	case BetaAssistantUpdateParamsModelO3Mini, BetaAssistantUpdateParamsModelO3Mini2025_01_31, BetaAssistantUpdateParamsModelO1, BetaAssistantUpdateParamsModelO1_2024_12_17, BetaAssistantUpdateParamsModelGPT4o, BetaAssistantUpdateParamsModelGPT4o2024_11_20, BetaAssistantUpdateParamsModelGPT4o2024_08_06, BetaAssistantUpdateParamsModelGPT4o2024_05_13, BetaAssistantUpdateParamsModelGPT4oMini, BetaAssistantUpdateParamsModelGPT4oMini2024_07_18, BetaAssistantUpdateParamsModelGPT4_5Preview, BetaAssistantUpdateParamsModelGPT4_5Preview2025_02_27, BetaAssistantUpdateParamsModelGPT4Turbo, BetaAssistantUpdateParamsModelGPT4Turbo2024_04_09, BetaAssistantUpdateParamsModelGPT4_0125Preview, BetaAssistantUpdateParamsModelGPT4TurboPreview, BetaAssistantUpdateParamsModelGPT4_1106Preview, BetaAssistantUpdateParamsModelGPT4VisionPreview, BetaAssistantUpdateParamsModelGPT4, BetaAssistantUpdateParamsModelGPT4_0314, BetaAssistantUpdateParamsModelGPT4_0613, BetaAssistantUpdateParamsModelGPT4_32k, BetaAssistantUpdateParamsModelGPT4_32k0314, BetaAssistantUpdateParamsModelGPT4_32k0613, BetaAssistantUpdateParamsModelGPT3_5Turbo, BetaAssistantUpdateParamsModelGPT3_5Turbo16k, BetaAssistantUpdateParamsModelGPT3_5Turbo0613, BetaAssistantUpdateParamsModelGPT3_5Turbo1106, BetaAssistantUpdateParamsModelGPT3_5Turbo0125, BetaAssistantUpdateParamsModelGPT3_5Turbo16k0613:
		return true
	}
	return false
}

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

func (r BetaAssistantUpdateParamsReasoningEffort) IsKnown() bool {
	switch r {
	case BetaAssistantUpdateParamsReasoningEffortLow, BetaAssistantUpdateParamsReasoningEffortMedium, BetaAssistantUpdateParamsReasoningEffortHigh:
		return true
	}
	return false
}

// A set of resources that are used by the assistant's tools. The resources are
// specific to the type of tool. For example, the `code_interpreter` tool requires
// a list of file IDs, while the `file_search` tool requires a list of vector store
// IDs.
type BetaAssistantUpdateParamsToolResources struct {
	CodeInterpreter param.Field[BetaAssistantUpdateParamsToolResourcesCodeInterpreter] `json:"code_interpreter"`
	FileSearch      param.Field[BetaAssistantUpdateParamsToolResourcesFileSearch]      `json:"file_search"`
}

func (r BetaAssistantUpdateParamsToolResources) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type BetaAssistantUpdateParamsToolResourcesCodeInterpreter struct {
	// Overrides the list of
	// [file](https://platform.openai.com/docs/api-reference/files) IDs made available
	// to the `code_interpreter` tool. There can be a maximum of 20 files associated
	// with the tool.
	FileIDs param.Field[[]string] `json:"file_ids"`
}

func (r BetaAssistantUpdateParamsToolResourcesCodeInterpreter) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type BetaAssistantUpdateParamsToolResourcesFileSearch struct {
	// Overrides the
	// [vector store](https://platform.openai.com/docs/api-reference/vector-stores/object)
	// attached to this assistant. There can be a maximum of 1 vector store attached to
	// the assistant.
	VectorStoreIDs param.Field[[]string] `json:"vector_store_ids"`
}

func (r BetaAssistantUpdateParamsToolResourcesFileSearch) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type BetaAssistantListParams struct {
	// A cursor for use in pagination. `after` is an object ID that defines your place
	// in the list. For instance, if you make a list request and receive 100 objects,
	// ending with obj_foo, your subsequent call can include after=obj_foo in order to
	// fetch the next page of the list.
	After param.Field[string] `query:"after"`
	// A cursor for use in pagination. `before` is an object ID that defines your place
	// in the list. For instance, if you make a list request and receive 100 objects,
	// starting with obj_foo, your subsequent call can include before=obj_foo in order
	// to fetch the previous page of the list.
	Before param.Field[string] `query:"before"`
	// A limit on the number of objects to be returned. Limit can range between 1 and
	// 100, and the default is 20.
	Limit param.Field[int64] `query:"limit"`
	// Sort order by the `created_at` timestamp of the objects. `asc` for ascending
	// order and `desc` for descending order.
	Order param.Field[BetaAssistantListParamsOrder] `query:"order"`
}

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

func (r BetaAssistantListParamsOrder) IsKnown() bool {
	switch r {
	case BetaAssistantListParamsOrderAsc, BetaAssistantListParamsOrderDesc:
		return true
	}
	return false
}

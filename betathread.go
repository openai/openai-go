// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/openai/openai-go/internal/apijson"
	"github.com/openai/openai-go/internal/requestconfig"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/param"
	"github.com/openai/openai-go/packages/resp"
	"github.com/openai/openai-go/packages/ssestream"
	"github.com/openai/openai-go/shared"
	"github.com/openai/openai-go/shared/constant"
)

// BetaThreadService contains methods and other services that help with interacting
// with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaThreadService] method instead.
type BetaThreadService struct {
	Options  []option.RequestOption
	Runs     BetaThreadRunService
	Messages BetaThreadMessageService
}

// NewBetaThreadService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewBetaThreadService(opts ...option.RequestOption) (r BetaThreadService) {
	r = BetaThreadService{}
	r.Options = opts
	r.Runs = NewBetaThreadRunService(opts...)
	r.Messages = NewBetaThreadMessageService(opts...)
	return
}

// Create a thread.
func (r *BetaThreadService) New(ctx context.Context, body BetaThreadNewParams, opts ...option.RequestOption) (res *Thread, err error) {
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "assistants=v2")}, opts...)
	path := "threads"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Retrieves a thread.
func (r *BetaThreadService) Get(ctx context.Context, threadID string, opts ...option.RequestOption) (res *Thread, err error) {
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "assistants=v2")}, opts...)
	if threadID == "" {
		err = errors.New("missing required thread_id parameter")
		return
	}
	path := fmt.Sprintf("threads/%s", threadID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// Modifies a thread.
func (r *BetaThreadService) Update(ctx context.Context, threadID string, body BetaThreadUpdateParams, opts ...option.RequestOption) (res *Thread, err error) {
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "assistants=v2")}, opts...)
	if threadID == "" {
		err = errors.New("missing required thread_id parameter")
		return
	}
	path := fmt.Sprintf("threads/%s", threadID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Delete a thread.
func (r *BetaThreadService) Delete(ctx context.Context, threadID string, opts ...option.RequestOption) (res *ThreadDeleted, err error) {
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "assistants=v2")}, opts...)
	if threadID == "" {
		err = errors.New("missing required thread_id parameter")
		return
	}
	path := fmt.Sprintf("threads/%s", threadID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return
}

// Create a thread and run it in one request.
func (r *BetaThreadService) NewAndRun(ctx context.Context, body BetaThreadNewAndRunParams, opts ...option.RequestOption) (res *Run, err error) {
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "assistants=v2")}, opts...)
	path := "threads/runs"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Create a thread and run it in one request.
func (r *BetaThreadService) NewAndRunStreaming(ctx context.Context, body BetaThreadNewAndRunParams, opts ...option.RequestOption) (stream *ssestream.Stream[AssistantStreamEventUnion]) {
	var (
		raw *http.Response
		err error
	)
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "assistants=v2"), option.WithJSONSet("stream", true)}, opts...)
	path := "threads/runs"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &raw, opts...)
	return ssestream.NewStream[AssistantStreamEventUnion](ssestream.NewDecoder(raw), err)
}

// Specifies a tool the model should use. Use to force the model to call a specific
// tool.
type AssistantToolChoice struct {
	// The type of the tool. If type is `function`, the function name must be set
	//
	// Any of "function", "code_interpreter", "file_search"
	Type     string                      `json:"type,omitzero,required"`
	Function AssistantToolChoiceFunction `json:"function,omitzero"`
	JSON     struct {
		Type     resp.Field
		Function resp.Field
		raw      string
	} `json:"-"`
}

func (r AssistantToolChoice) RawJSON() string { return r.JSON.raw }
func (r *AssistantToolChoice) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this AssistantToolChoice to a AssistantToolChoiceParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// AssistantToolChoiceParam.IsOverridden()
func (r AssistantToolChoice) ToParam() AssistantToolChoiceParam {
	return param.Override[AssistantToolChoiceParam](r.RawJSON())
}

// The type of the tool. If type is `function`, the function name must be set
type AssistantToolChoiceType = string

const (
	AssistantToolChoiceTypeFunction        AssistantToolChoiceType = "function"
	AssistantToolChoiceTypeCodeInterpreter AssistantToolChoiceType = "code_interpreter"
	AssistantToolChoiceTypeFileSearch      AssistantToolChoiceType = "file_search"
)

// Specifies a tool the model should use. Use to force the model to call a specific
// tool.
type AssistantToolChoiceParam struct {
	// The type of the tool. If type is `function`, the function name must be set
	//
	// Any of "function", "code_interpreter", "file_search"
	Type     string                           `json:"type,omitzero,required"`
	Function AssistantToolChoiceFunctionParam `json:"function,omitzero"`
	apiobject
}

func (f AssistantToolChoiceParam) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

func (r AssistantToolChoiceParam) MarshalJSON() (data []byte, err error) {
	type shadow AssistantToolChoiceParam
	return param.MarshalObject(r, (*shadow)(&r))
}

type AssistantToolChoiceFunction struct {
	// The name of the function to call.
	Name string `json:"name,omitzero,required"`
	JSON struct {
		Name resp.Field
		raw  string
	} `json:"-"`
}

func (r AssistantToolChoiceFunction) RawJSON() string { return r.JSON.raw }
func (r *AssistantToolChoiceFunction) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this AssistantToolChoiceFunction to a
// AssistantToolChoiceFunctionParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// AssistantToolChoiceFunctionParam.IsOverridden()
func (r AssistantToolChoiceFunction) ToParam() AssistantToolChoiceFunctionParam {
	return param.Override[AssistantToolChoiceFunctionParam](r.RawJSON())
}

type AssistantToolChoiceFunctionParam struct {
	// The name of the function to call.
	Name param.String `json:"name,omitzero,required"`
	apiobject
}

func (f AssistantToolChoiceFunctionParam) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

func (r AssistantToolChoiceFunctionParam) MarshalJSON() (data []byte, err error) {
	type shadow AssistantToolChoiceFunctionParam
	return param.MarshalObject(r, (*shadow)(&r))
}

type AssistantToolChoiceOptionUnion struct {
	OfString string                      `json:",inline"`
	Type     string                      `json:"type"`
	Function AssistantToolChoiceFunction `json:"function"`
	JSON     struct {
		OfString resp.Field
		Type     resp.Field
		Function resp.Field
		raw      string
	} `json:"-"`
}

func (u AssistantToolChoiceOptionUnion) AsAuto() (v string) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AssistantToolChoiceOptionUnion) AsAssistantToolChoice() (v AssistantToolChoice) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AssistantToolChoiceOptionUnion) RawJSON() string { return u.JSON.raw }

func (r *AssistantToolChoiceOptionUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this AssistantToolChoiceOptionUnion to a
// AssistantToolChoiceOptionUnionParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// AssistantToolChoiceOptionUnionParam.IsOverridden()
func (r AssistantToolChoiceOptionUnion) ToParam() AssistantToolChoiceOptionUnionParam {
	return param.Override[AssistantToolChoiceOptionUnionParam](r.RawJSON())
}

// `none` means the model will not call any tools and instead generates a message.
// `auto` means the model can pick between generating a message or calling one or
// more tools. `required` means the model must call one or more tools before
// responding to the user.
type AssistantToolChoiceOptionAuto = string

const (
	AssistantToolChoiceOptionAutoNone     AssistantToolChoiceOptionAuto = "none"
	AssistantToolChoiceOptionAutoAuto     AssistantToolChoiceOptionAuto = "auto"
	AssistantToolChoiceOptionAutoRequired AssistantToolChoiceOptionAuto = "required"
)

func NewAssistantToolChoiceOptionOfAssistantToolChoice(type_ string) AssistantToolChoiceOptionUnionParam {
	var variant AssistantToolChoiceParam
	variant.Type = type_
	return AssistantToolChoiceOptionUnionParam{OfAssistantToolChoice: &variant}
}

// Only one field can be non-zero
type AssistantToolChoiceOptionUnionParam struct {
	// Check if union is this variant with !param.IsOmitted(union.OfAuto)
	OfAuto                string
	OfAssistantToolChoice *AssistantToolChoiceParam
	apiunion
}

func (u AssistantToolChoiceOptionUnionParam) IsMissing() bool {
	return param.IsOmitted(u) || u.IsNull()
}

func (u AssistantToolChoiceOptionUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion[AssistantToolChoiceOptionUnionParam](u.OfAuto, u.OfAssistantToolChoice)
}

func (u AssistantToolChoiceOptionUnionParam) GetType() *string {
	if vt := u.OfAssistantToolChoice; vt != nil {
		return &vt.Type
	}
	return nil
}

func (u AssistantToolChoiceOptionUnionParam) GetFunction() *AssistantToolChoiceFunctionParam {
	if vt := u.OfAssistantToolChoice; vt != nil {
		return &vt.Function
	}
	return nil
}

// Represents a thread that contains
// [messages](https://platform.openai.com/docs/api-reference/messages).
type Thread struct {
	// The identifier, which can be referenced in API endpoints.
	ID string `json:"id,omitzero,required"`
	// The Unix timestamp (in seconds) for when the thread was created.
	CreatedAt int64 `json:"created_at,omitzero,required"`
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard.
	//
	// Keys are strings with a maximum length of 64 characters. Values are strings with
	// a maximum length of 512 characters.
	Metadata shared.Metadata `json:"metadata,omitzero,required,nullable"`
	// The object type, which is always `thread`.
	//
	// This field can be elided, and will be automatically set as "thread".
	Object constant.Thread `json:"object,required"`
	// A set of resources that are made available to the assistant's tools in this
	// thread. The resources are specific to the type of tool. For example, the
	// `code_interpreter` tool requires a list of file IDs, while the `file_search`
	// tool requires a list of vector store IDs.
	ToolResources ThreadToolResources `json:"tool_resources,omitzero,required,nullable"`
	JSON          struct {
		ID            resp.Field
		CreatedAt     resp.Field
		Metadata      resp.Field
		Object        resp.Field
		ToolResources resp.Field
		raw           string
	} `json:"-"`
}

func (r Thread) RawJSON() string { return r.JSON.raw }
func (r *Thread) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A set of resources that are made available to the assistant's tools in this
// thread. The resources are specific to the type of tool. For example, the
// `code_interpreter` tool requires a list of file IDs, while the `file_search`
// tool requires a list of vector store IDs.
type ThreadToolResources struct {
	CodeInterpreter ThreadToolResourcesCodeInterpreter `json:"code_interpreter,omitzero"`
	FileSearch      ThreadToolResourcesFileSearch      `json:"file_search,omitzero"`
	JSON            struct {
		CodeInterpreter resp.Field
		FileSearch      resp.Field
		raw             string
	} `json:"-"`
}

func (r ThreadToolResources) RawJSON() string { return r.JSON.raw }
func (r *ThreadToolResources) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ThreadToolResourcesCodeInterpreter struct {
	// A list of [file](https://platform.openai.com/docs/api-reference/files) IDs made
	// available to the `code_interpreter` tool. There can be a maximum of 20 files
	// associated with the tool.
	FileIDs []string `json:"file_ids,omitzero"`
	JSON    struct {
		FileIDs resp.Field
		raw     string
	} `json:"-"`
}

func (r ThreadToolResourcesCodeInterpreter) RawJSON() string { return r.JSON.raw }
func (r *ThreadToolResourcesCodeInterpreter) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ThreadToolResourcesFileSearch struct {
	// The
	// [vector store](https://platform.openai.com/docs/api-reference/vector-stores/object)
	// attached to this thread. There can be a maximum of 1 vector store attached to
	// the thread.
	VectorStoreIDs []string `json:"vector_store_ids,omitzero"`
	JSON           struct {
		VectorStoreIDs resp.Field
		raw            string
	} `json:"-"`
}

func (r ThreadToolResourcesFileSearch) RawJSON() string { return r.JSON.raw }
func (r *ThreadToolResourcesFileSearch) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ThreadDeleted struct {
	ID      string `json:"id,omitzero,required"`
	Deleted bool   `json:"deleted,omitzero,required"`
	// This field can be elided, and will be automatically set as "thread.deleted".
	Object constant.ThreadDeleted `json:"object,required"`
	JSON   struct {
		ID      resp.Field
		Deleted resp.Field
		Object  resp.Field
		raw     string
	} `json:"-"`
}

func (r ThreadDeleted) RawJSON() string { return r.JSON.raw }
func (r *ThreadDeleted) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaThreadNewParams struct {
	// A list of [messages](https://platform.openai.com/docs/api-reference/messages) to
	// start the thread with.
	Messages []BetaThreadNewParamsMessage `json:"messages,omitzero"`
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard.
	//
	// Keys are strings with a maximum length of 64 characters. Values are strings with
	// a maximum length of 512 characters.
	Metadata shared.MetadataParam `json:"metadata,omitzero"`
	// A set of resources that are made available to the assistant's tools in this
	// thread. The resources are specific to the type of tool. For example, the
	// `code_interpreter` tool requires a list of file IDs, while the `file_search`
	// tool requires a list of vector store IDs.
	ToolResources BetaThreadNewParamsToolResources `json:"tool_resources,omitzero"`
	apiobject
}

func (f BetaThreadNewParams) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

func (r BetaThreadNewParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaThreadNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}

type BetaThreadNewParamsMessage struct {
	// An array of content parts with a defined type, each can be of type `text` or
	// images can be passed with `image_url` or `image_file`. Image types are only
	// supported on
	// [Vision-compatible models](https://platform.openai.com/docs/models).
	Content []MessageContentPartParamUnion `json:"content,omitzero,required"`
	// The role of the entity that is creating the message. Allowed values include:
	//
	//   - `user`: Indicates the message is sent by an actual user and should be used in
	//     most cases to represent user-generated messages.
	//   - `assistant`: Indicates the message is generated by the assistant. Use this
	//     value to insert messages from the assistant into the conversation.
	//
	// Any of "user", "assistant"
	Role string `json:"role,omitzero,required"`
	// A list of files attached to the message, and the tools they should be added to.
	Attachments []BetaThreadNewParamsMessagesAttachment `json:"attachments,omitzero"`
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard.
	//
	// Keys are strings with a maximum length of 64 characters. Values are strings with
	// a maximum length of 512 characters.
	Metadata shared.MetadataParam `json:"metadata,omitzero"`
	apiobject
}

func (f BetaThreadNewParamsMessage) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

func (r BetaThreadNewParamsMessage) MarshalJSON() (data []byte, err error) {
	type shadow BetaThreadNewParamsMessage
	return param.MarshalObject(r, (*shadow)(&r))
}

// The role of the entity that is creating the message. Allowed values include:
//
//   - `user`: Indicates the message is sent by an actual user and should be used in
//     most cases to represent user-generated messages.
//   - `assistant`: Indicates the message is generated by the assistant. Use this
//     value to insert messages from the assistant into the conversation.
type BetaThreadNewParamsMessagesRole = string

const (
	BetaThreadNewParamsMessagesRoleUser      BetaThreadNewParamsMessagesRole = "user"
	BetaThreadNewParamsMessagesRoleAssistant BetaThreadNewParamsMessagesRole = "assistant"
)

type BetaThreadNewParamsMessagesAttachment struct {
	// The ID of the file to attach to the message.
	FileID param.String `json:"file_id,omitzero"`
	// The tools to add this file to.
	Tools []BetaThreadNewParamsMessagesAttachmentsToolUnion `json:"tools,omitzero"`
	apiobject
}

func (f BetaThreadNewParamsMessagesAttachment) IsMissing() bool {
	return param.IsOmitted(f) || f.IsNull()
}

func (r BetaThreadNewParamsMessagesAttachment) MarshalJSON() (data []byte, err error) {
	type shadow BetaThreadNewParamsMessagesAttachment
	return param.MarshalObject(r, (*shadow)(&r))
}

// Only one field can be non-zero
type BetaThreadNewParamsMessagesAttachmentsToolUnion struct {
	OfCodeInterpreter *CodeInterpreterToolParam
	OfFileSearch      *BetaThreadNewParamsMessagesAttachmentsToolsFileSearch
	apiunion
}

func (u BetaThreadNewParamsMessagesAttachmentsToolUnion) IsMissing() bool {
	return param.IsOmitted(u) || u.IsNull()
}

func (u BetaThreadNewParamsMessagesAttachmentsToolUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion[BetaThreadNewParamsMessagesAttachmentsToolUnion](u.OfCodeInterpreter, u.OfFileSearch)
}

func (u BetaThreadNewParamsMessagesAttachmentsToolUnion) GetType() *string {
	if vt := u.OfCodeInterpreter; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfFileSearch; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

type BetaThreadNewParamsMessagesAttachmentsToolsFileSearch struct {
	// The type of tool being defined: `file_search`
	//
	// This field can be elided, and will be automatically set as "file_search".
	Type constant.FileSearch `json:"type,required"`
	apiobject
}

func (f BetaThreadNewParamsMessagesAttachmentsToolsFileSearch) IsMissing() bool {
	return param.IsOmitted(f) || f.IsNull()
}

func (r BetaThreadNewParamsMessagesAttachmentsToolsFileSearch) MarshalJSON() (data []byte, err error) {
	type shadow BetaThreadNewParamsMessagesAttachmentsToolsFileSearch
	return param.MarshalObject(r, (*shadow)(&r))
}

// A set of resources that are made available to the assistant's tools in this
// thread. The resources are specific to the type of tool. For example, the
// `code_interpreter` tool requires a list of file IDs, while the `file_search`
// tool requires a list of vector store IDs.
type BetaThreadNewParamsToolResources struct {
	CodeInterpreter BetaThreadNewParamsToolResourcesCodeInterpreter `json:"code_interpreter,omitzero"`
	FileSearch      BetaThreadNewParamsToolResourcesFileSearch      `json:"file_search,omitzero"`
	apiobject
}

func (f BetaThreadNewParamsToolResources) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

func (r BetaThreadNewParamsToolResources) MarshalJSON() (data []byte, err error) {
	type shadow BetaThreadNewParamsToolResources
	return param.MarshalObject(r, (*shadow)(&r))
}

type BetaThreadNewParamsToolResourcesCodeInterpreter struct {
	// A list of [file](https://platform.openai.com/docs/api-reference/files) IDs made
	// available to the `code_interpreter` tool. There can be a maximum of 20 files
	// associated with the tool.
	FileIDs []string `json:"file_ids,omitzero"`
	apiobject
}

func (f BetaThreadNewParamsToolResourcesCodeInterpreter) IsMissing() bool {
	return param.IsOmitted(f) || f.IsNull()
}

func (r BetaThreadNewParamsToolResourcesCodeInterpreter) MarshalJSON() (data []byte, err error) {
	type shadow BetaThreadNewParamsToolResourcesCodeInterpreter
	return param.MarshalObject(r, (*shadow)(&r))
}

type BetaThreadNewParamsToolResourcesFileSearch struct {
	// The
	// [vector store](https://platform.openai.com/docs/api-reference/vector-stores/object)
	// attached to this thread. There can be a maximum of 1 vector store attached to
	// the thread.
	VectorStoreIDs []string `json:"vector_store_ids,omitzero"`
	// A helper to create a
	// [vector store](https://platform.openai.com/docs/api-reference/vector-stores/object)
	// with file_ids and attach it to this thread. There can be a maximum of 1 vector
	// store attached to the thread.
	VectorStores []BetaThreadNewParamsToolResourcesFileSearchVectorStore `json:"vector_stores,omitzero"`
	apiobject
}

func (f BetaThreadNewParamsToolResourcesFileSearch) IsMissing() bool {
	return param.IsOmitted(f) || f.IsNull()
}

func (r BetaThreadNewParamsToolResourcesFileSearch) MarshalJSON() (data []byte, err error) {
	type shadow BetaThreadNewParamsToolResourcesFileSearch
	return param.MarshalObject(r, (*shadow)(&r))
}

type BetaThreadNewParamsToolResourcesFileSearchVectorStore struct {
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

func (f BetaThreadNewParamsToolResourcesFileSearchVectorStore) IsMissing() bool {
	return param.IsOmitted(f) || f.IsNull()
}

func (r BetaThreadNewParamsToolResourcesFileSearchVectorStore) MarshalJSON() (data []byte, err error) {
	type shadow BetaThreadNewParamsToolResourcesFileSearchVectorStore
	return param.MarshalObject(r, (*shadow)(&r))
}

type BetaThreadUpdateParams struct {
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard.
	//
	// Keys are strings with a maximum length of 64 characters. Values are strings with
	// a maximum length of 512 characters.
	Metadata shared.MetadataParam `json:"metadata,omitzero"`
	// A set of resources that are made available to the assistant's tools in this
	// thread. The resources are specific to the type of tool. For example, the
	// `code_interpreter` tool requires a list of file IDs, while the `file_search`
	// tool requires a list of vector store IDs.
	ToolResources BetaThreadUpdateParamsToolResources `json:"tool_resources,omitzero"`
	apiobject
}

func (f BetaThreadUpdateParams) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

func (r BetaThreadUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaThreadUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}

// A set of resources that are made available to the assistant's tools in this
// thread. The resources are specific to the type of tool. For example, the
// `code_interpreter` tool requires a list of file IDs, while the `file_search`
// tool requires a list of vector store IDs.
type BetaThreadUpdateParamsToolResources struct {
	CodeInterpreter BetaThreadUpdateParamsToolResourcesCodeInterpreter `json:"code_interpreter,omitzero"`
	FileSearch      BetaThreadUpdateParamsToolResourcesFileSearch      `json:"file_search,omitzero"`
	apiobject
}

func (f BetaThreadUpdateParamsToolResources) IsMissing() bool {
	return param.IsOmitted(f) || f.IsNull()
}

func (r BetaThreadUpdateParamsToolResources) MarshalJSON() (data []byte, err error) {
	type shadow BetaThreadUpdateParamsToolResources
	return param.MarshalObject(r, (*shadow)(&r))
}

type BetaThreadUpdateParamsToolResourcesCodeInterpreter struct {
	// A list of [file](https://platform.openai.com/docs/api-reference/files) IDs made
	// available to the `code_interpreter` tool. There can be a maximum of 20 files
	// associated with the tool.
	FileIDs []string `json:"file_ids,omitzero"`
	apiobject
}

func (f BetaThreadUpdateParamsToolResourcesCodeInterpreter) IsMissing() bool {
	return param.IsOmitted(f) || f.IsNull()
}

func (r BetaThreadUpdateParamsToolResourcesCodeInterpreter) MarshalJSON() (data []byte, err error) {
	type shadow BetaThreadUpdateParamsToolResourcesCodeInterpreter
	return param.MarshalObject(r, (*shadow)(&r))
}

type BetaThreadUpdateParamsToolResourcesFileSearch struct {
	// The
	// [vector store](https://platform.openai.com/docs/api-reference/vector-stores/object)
	// attached to this thread. There can be a maximum of 1 vector store attached to
	// the thread.
	VectorStoreIDs []string `json:"vector_store_ids,omitzero"`
	apiobject
}

func (f BetaThreadUpdateParamsToolResourcesFileSearch) IsMissing() bool {
	return param.IsOmitted(f) || f.IsNull()
}

func (r BetaThreadUpdateParamsToolResourcesFileSearch) MarshalJSON() (data []byte, err error) {
	type shadow BetaThreadUpdateParamsToolResourcesFileSearch
	return param.MarshalObject(r, (*shadow)(&r))
}

type BetaThreadNewAndRunParams struct {
	// The ID of the
	// [assistant](https://platform.openai.com/docs/api-reference/assistants) to use to
	// execute this run.
	AssistantID param.String `json:"assistant_id,omitzero,required"`
	// Override the default system message of the assistant. This is useful for
	// modifying the behavior on a per-run basis.
	Instructions param.String `json:"instructions,omitzero"`
	// The maximum number of completion tokens that may be used over the course of the
	// run. The run will make a best effort to use only the number of completion tokens
	// specified, across multiple turns of the run. If the run exceeds the number of
	// completion tokens specified, the run will end with status `incomplete`. See
	// `incomplete_details` for more info.
	MaxCompletionTokens param.Int `json:"max_completion_tokens,omitzero"`
	// The maximum number of prompt tokens that may be used over the course of the run.
	// The run will make a best effort to use only the number of prompt tokens
	// specified, across multiple turns of the run. If the run exceeds the number of
	// prompt tokens specified, the run will end with status `incomplete`. See
	// `incomplete_details` for more info.
	MaxPromptTokens param.Int `json:"max_prompt_tokens,omitzero"`
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard.
	//
	// Keys are strings with a maximum length of 64 characters. Values are strings with
	// a maximum length of 512 characters.
	Metadata shared.MetadataParam `json:"metadata,omitzero"`
	// The ID of the [Model](https://platform.openai.com/docs/api-reference/models) to
	// be used to execute this run. If a value is provided here, it will override the
	// model associated with the assistant. If not, the model associated with the
	// assistant will be used.
	Model ChatModel `json:"model,omitzero"`
	// Whether to enable
	// [parallel function calling](https://platform.openai.com/docs/guides/function-calling#configuring-parallel-function-calling)
	// during tool use.
	ParallelToolCalls param.Bool `json:"parallel_tool_calls,omitzero"`
	// What sampling temperature to use, between 0 and 2. Higher values like 0.8 will
	// make the output more random, while lower values like 0.2 will make it more
	// focused and deterministic.
	Temperature param.Float `json:"temperature,omitzero"`
	// Options to create a new thread. If no thread is provided when running a request,
	// an empty thread will be created.
	Thread BetaThreadNewAndRunParamsThread `json:"thread,omitzero"`
	// Controls which (if any) tool is called by the model. `none` means the model will
	// not call any tools and instead generates a message. `auto` is the default value
	// and means the model can pick between generating a message or calling one or more
	// tools. `required` means the model must call one or more tools before responding
	// to the user. Specifying a particular tool like `{"type": "file_search"}` or
	// `{"type": "function", "function": {"name": "my_function"}}` forces the model to
	// call that tool.
	ToolChoice AssistantToolChoiceOptionUnionParam `json:"tool_choice,omitzero"`
	// A set of resources that are used by the assistant's tools. The resources are
	// specific to the type of tool. For example, the `code_interpreter` tool requires
	// a list of file IDs, while the `file_search` tool requires a list of vector store
	// IDs.
	ToolResources BetaThreadNewAndRunParamsToolResources `json:"tool_resources,omitzero"`
	// Override the tools the assistant can use for this run. This is useful for
	// modifying the behavior on a per-run basis.
	Tools []BetaThreadNewAndRunParamsToolUnion `json:"tools,omitzero"`
	// An alternative to sampling with temperature, called nucleus sampling, where the
	// model considers the results of the tokens with top_p probability mass. So 0.1
	// means only the tokens comprising the top 10% probability mass are considered.
	//
	// We generally recommend altering this or temperature but not both.
	TopP param.Float `json:"top_p,omitzero"`
	// Controls for how a thread will be truncated prior to the run. Use this to
	// control the intial context window of the run.
	TruncationStrategy BetaThreadNewAndRunParamsTruncationStrategy `json:"truncation_strategy,omitzero"`
	apiobject
}

func (f BetaThreadNewAndRunParams) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

func (r BetaThreadNewAndRunParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaThreadNewAndRunParams
	return param.MarshalObject(r, (*shadow)(&r))
}

// Options to create a new thread. If no thread is provided when running a request,
// an empty thread will be created.
type BetaThreadNewAndRunParamsThread struct {
	// A list of [messages](https://platform.openai.com/docs/api-reference/messages) to
	// start the thread with.
	Messages []BetaThreadNewAndRunParamsThreadMessage `json:"messages,omitzero"`
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard.
	//
	// Keys are strings with a maximum length of 64 characters. Values are strings with
	// a maximum length of 512 characters.
	Metadata shared.MetadataParam `json:"metadata,omitzero"`
	// A set of resources that are made available to the assistant's tools in this
	// thread. The resources are specific to the type of tool. For example, the
	// `code_interpreter` tool requires a list of file IDs, while the `file_search`
	// tool requires a list of vector store IDs.
	ToolResources BetaThreadNewAndRunParamsThreadToolResources `json:"tool_resources,omitzero"`
	apiobject
}

func (f BetaThreadNewAndRunParamsThread) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

func (r BetaThreadNewAndRunParamsThread) MarshalJSON() (data []byte, err error) {
	type shadow BetaThreadNewAndRunParamsThread
	return param.MarshalObject(r, (*shadow)(&r))
}

type BetaThreadNewAndRunParamsThreadMessage struct {
	// An array of content parts with a defined type, each can be of type `text` or
	// images can be passed with `image_url` or `image_file`. Image types are only
	// supported on
	// [Vision-compatible models](https://platform.openai.com/docs/models).
	Content []MessageContentPartParamUnion `json:"content,omitzero,required"`
	// The role of the entity that is creating the message. Allowed values include:
	//
	//   - `user`: Indicates the message is sent by an actual user and should be used in
	//     most cases to represent user-generated messages.
	//   - `assistant`: Indicates the message is generated by the assistant. Use this
	//     value to insert messages from the assistant into the conversation.
	//
	// Any of "user", "assistant"
	Role string `json:"role,omitzero,required"`
	// A list of files attached to the message, and the tools they should be added to.
	Attachments []BetaThreadNewAndRunParamsThreadMessagesAttachment `json:"attachments,omitzero"`
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard.
	//
	// Keys are strings with a maximum length of 64 characters. Values are strings with
	// a maximum length of 512 characters.
	Metadata shared.MetadataParam `json:"metadata,omitzero"`
	apiobject
}

func (f BetaThreadNewAndRunParamsThreadMessage) IsMissing() bool {
	return param.IsOmitted(f) || f.IsNull()
}

func (r BetaThreadNewAndRunParamsThreadMessage) MarshalJSON() (data []byte, err error) {
	type shadow BetaThreadNewAndRunParamsThreadMessage
	return param.MarshalObject(r, (*shadow)(&r))
}

// The role of the entity that is creating the message. Allowed values include:
//
//   - `user`: Indicates the message is sent by an actual user and should be used in
//     most cases to represent user-generated messages.
//   - `assistant`: Indicates the message is generated by the assistant. Use this
//     value to insert messages from the assistant into the conversation.
type BetaThreadNewAndRunParamsThreadMessagesRole = string

const (
	BetaThreadNewAndRunParamsThreadMessagesRoleUser      BetaThreadNewAndRunParamsThreadMessagesRole = "user"
	BetaThreadNewAndRunParamsThreadMessagesRoleAssistant BetaThreadNewAndRunParamsThreadMessagesRole = "assistant"
)

type BetaThreadNewAndRunParamsThreadMessagesAttachment struct {
	// The ID of the file to attach to the message.
	FileID param.String `json:"file_id,omitzero"`
	// The tools to add this file to.
	Tools []BetaThreadNewAndRunParamsThreadMessagesAttachmentsToolUnion `json:"tools,omitzero"`
	apiobject
}

func (f BetaThreadNewAndRunParamsThreadMessagesAttachment) IsMissing() bool {
	return param.IsOmitted(f) || f.IsNull()
}

func (r BetaThreadNewAndRunParamsThreadMessagesAttachment) MarshalJSON() (data []byte, err error) {
	type shadow BetaThreadNewAndRunParamsThreadMessagesAttachment
	return param.MarshalObject(r, (*shadow)(&r))
}

// Only one field can be non-zero
type BetaThreadNewAndRunParamsThreadMessagesAttachmentsToolUnion struct {
	OfCodeInterpreter *CodeInterpreterToolParam
	OfFileSearch      *BetaThreadNewAndRunParamsThreadMessagesAttachmentsToolsFileSearch
	apiunion
}

func (u BetaThreadNewAndRunParamsThreadMessagesAttachmentsToolUnion) IsMissing() bool {
	return param.IsOmitted(u) || u.IsNull()
}

func (u BetaThreadNewAndRunParamsThreadMessagesAttachmentsToolUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion[BetaThreadNewAndRunParamsThreadMessagesAttachmentsToolUnion](u.OfCodeInterpreter, u.OfFileSearch)
}

func (u BetaThreadNewAndRunParamsThreadMessagesAttachmentsToolUnion) GetType() *string {
	if vt := u.OfCodeInterpreter; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfFileSearch; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

type BetaThreadNewAndRunParamsThreadMessagesAttachmentsToolsFileSearch struct {
	// The type of tool being defined: `file_search`
	//
	// This field can be elided, and will be automatically set as "file_search".
	Type constant.FileSearch `json:"type,required"`
	apiobject
}

func (f BetaThreadNewAndRunParamsThreadMessagesAttachmentsToolsFileSearch) IsMissing() bool {
	return param.IsOmitted(f) || f.IsNull()
}

func (r BetaThreadNewAndRunParamsThreadMessagesAttachmentsToolsFileSearch) MarshalJSON() (data []byte, err error) {
	type shadow BetaThreadNewAndRunParamsThreadMessagesAttachmentsToolsFileSearch
	return param.MarshalObject(r, (*shadow)(&r))
}

// A set of resources that are made available to the assistant's tools in this
// thread. The resources are specific to the type of tool. For example, the
// `code_interpreter` tool requires a list of file IDs, while the `file_search`
// tool requires a list of vector store IDs.
type BetaThreadNewAndRunParamsThreadToolResources struct {
	CodeInterpreter BetaThreadNewAndRunParamsThreadToolResourcesCodeInterpreter `json:"code_interpreter,omitzero"`
	FileSearch      BetaThreadNewAndRunParamsThreadToolResourcesFileSearch      `json:"file_search,omitzero"`
	apiobject
}

func (f BetaThreadNewAndRunParamsThreadToolResources) IsMissing() bool {
	return param.IsOmitted(f) || f.IsNull()
}

func (r BetaThreadNewAndRunParamsThreadToolResources) MarshalJSON() (data []byte, err error) {
	type shadow BetaThreadNewAndRunParamsThreadToolResources
	return param.MarshalObject(r, (*shadow)(&r))
}

type BetaThreadNewAndRunParamsThreadToolResourcesCodeInterpreter struct {
	// A list of [file](https://platform.openai.com/docs/api-reference/files) IDs made
	// available to the `code_interpreter` tool. There can be a maximum of 20 files
	// associated with the tool.
	FileIDs []string `json:"file_ids,omitzero"`
	apiobject
}

func (f BetaThreadNewAndRunParamsThreadToolResourcesCodeInterpreter) IsMissing() bool {
	return param.IsOmitted(f) || f.IsNull()
}

func (r BetaThreadNewAndRunParamsThreadToolResourcesCodeInterpreter) MarshalJSON() (data []byte, err error) {
	type shadow BetaThreadNewAndRunParamsThreadToolResourcesCodeInterpreter
	return param.MarshalObject(r, (*shadow)(&r))
}

type BetaThreadNewAndRunParamsThreadToolResourcesFileSearch struct {
	// The
	// [vector store](https://platform.openai.com/docs/api-reference/vector-stores/object)
	// attached to this thread. There can be a maximum of 1 vector store attached to
	// the thread.
	VectorStoreIDs []string `json:"vector_store_ids,omitzero"`
	// A helper to create a
	// [vector store](https://platform.openai.com/docs/api-reference/vector-stores/object)
	// with file_ids and attach it to this thread. There can be a maximum of 1 vector
	// store attached to the thread.
	VectorStores []BetaThreadNewAndRunParamsThreadToolResourcesFileSearchVectorStore `json:"vector_stores,omitzero"`
	apiobject
}

func (f BetaThreadNewAndRunParamsThreadToolResourcesFileSearch) IsMissing() bool {
	return param.IsOmitted(f) || f.IsNull()
}

func (r BetaThreadNewAndRunParamsThreadToolResourcesFileSearch) MarshalJSON() (data []byte, err error) {
	type shadow BetaThreadNewAndRunParamsThreadToolResourcesFileSearch
	return param.MarshalObject(r, (*shadow)(&r))
}

type BetaThreadNewAndRunParamsThreadToolResourcesFileSearchVectorStore struct {
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

func (f BetaThreadNewAndRunParamsThreadToolResourcesFileSearchVectorStore) IsMissing() bool {
	return param.IsOmitted(f) || f.IsNull()
}

func (r BetaThreadNewAndRunParamsThreadToolResourcesFileSearchVectorStore) MarshalJSON() (data []byte, err error) {
	type shadow BetaThreadNewAndRunParamsThreadToolResourcesFileSearchVectorStore
	return param.MarshalObject(r, (*shadow)(&r))
}

// A set of resources that are used by the assistant's tools. The resources are
// specific to the type of tool. For example, the `code_interpreter` tool requires
// a list of file IDs, while the `file_search` tool requires a list of vector store
// IDs.
type BetaThreadNewAndRunParamsToolResources struct {
	CodeInterpreter BetaThreadNewAndRunParamsToolResourcesCodeInterpreter `json:"code_interpreter,omitzero"`
	FileSearch      BetaThreadNewAndRunParamsToolResourcesFileSearch      `json:"file_search,omitzero"`
	apiobject
}

func (f BetaThreadNewAndRunParamsToolResources) IsMissing() bool {
	return param.IsOmitted(f) || f.IsNull()
}

func (r BetaThreadNewAndRunParamsToolResources) MarshalJSON() (data []byte, err error) {
	type shadow BetaThreadNewAndRunParamsToolResources
	return param.MarshalObject(r, (*shadow)(&r))
}

type BetaThreadNewAndRunParamsToolResourcesCodeInterpreter struct {
	// A list of [file](https://platform.openai.com/docs/api-reference/files) IDs made
	// available to the `code_interpreter` tool. There can be a maximum of 20 files
	// associated with the tool.
	FileIDs []string `json:"file_ids,omitzero"`
	apiobject
}

func (f BetaThreadNewAndRunParamsToolResourcesCodeInterpreter) IsMissing() bool {
	return param.IsOmitted(f) || f.IsNull()
}

func (r BetaThreadNewAndRunParamsToolResourcesCodeInterpreter) MarshalJSON() (data []byte, err error) {
	type shadow BetaThreadNewAndRunParamsToolResourcesCodeInterpreter
	return param.MarshalObject(r, (*shadow)(&r))
}

type BetaThreadNewAndRunParamsToolResourcesFileSearch struct {
	// The ID of the
	// [vector store](https://platform.openai.com/docs/api-reference/vector-stores/object)
	// attached to this assistant. There can be a maximum of 1 vector store attached to
	// the assistant.
	VectorStoreIDs []string `json:"vector_store_ids,omitzero"`
	apiobject
}

func (f BetaThreadNewAndRunParamsToolResourcesFileSearch) IsMissing() bool {
	return param.IsOmitted(f) || f.IsNull()
}

func (r BetaThreadNewAndRunParamsToolResourcesFileSearch) MarshalJSON() (data []byte, err error) {
	type shadow BetaThreadNewAndRunParamsToolResourcesFileSearch
	return param.MarshalObject(r, (*shadow)(&r))
}

// Only one field can be non-zero
type BetaThreadNewAndRunParamsToolUnion struct {
	OfCodeInterpreterTool *CodeInterpreterToolParam
	OfFileSearchTool      *FileSearchToolParam
	OfFunctionTool        *FunctionToolParam
	apiunion
}

func (u BetaThreadNewAndRunParamsToolUnion) IsMissing() bool { return param.IsOmitted(u) || u.IsNull() }

func (u BetaThreadNewAndRunParamsToolUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion[BetaThreadNewAndRunParamsToolUnion](u.OfCodeInterpreterTool, u.OfFileSearchTool, u.OfFunctionTool)
}

func (u BetaThreadNewAndRunParamsToolUnion) GetFileSearch() *FileSearchToolFileSearchParam {
	if vt := u.OfFileSearchTool; vt != nil {
		return &vt.FileSearch
	}
	return nil
}

func (u BetaThreadNewAndRunParamsToolUnion) GetFunction() *shared.FunctionDefinitionParam {
	if vt := u.OfFunctionTool; vt != nil {
		return &vt.Function
	}
	return nil
}

func (u BetaThreadNewAndRunParamsToolUnion) GetType() *string {
	if vt := u.OfCodeInterpreterTool; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfFileSearchTool; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfFunctionTool; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Controls for how a thread will be truncated prior to the run. Use this to
// control the intial context window of the run.
type BetaThreadNewAndRunParamsTruncationStrategy struct {
	// The truncation strategy to use for the thread. The default is `auto`. If set to
	// `last_messages`, the thread will be truncated to the n most recent messages in
	// the thread. When set to `auto`, messages in the middle of the thread will be
	// dropped to fit the context length of the model, `max_prompt_tokens`.
	//
	// Any of "auto", "last_messages"
	Type string `json:"type,omitzero,required"`
	// The number of most recent messages from the thread when constructing the context
	// for the run.
	LastMessages param.Int `json:"last_messages,omitzero"`
	apiobject
}

func (f BetaThreadNewAndRunParamsTruncationStrategy) IsMissing() bool {
	return param.IsOmitted(f) || f.IsNull()
}

func (r BetaThreadNewAndRunParamsTruncationStrategy) MarshalJSON() (data []byte, err error) {
	type shadow BetaThreadNewAndRunParamsTruncationStrategy
	return param.MarshalObject(r, (*shadow)(&r))
}

// The truncation strategy to use for the thread. The default is `auto`. If set to
// `last_messages`, the thread will be truncated to the n most recent messages in
// the thread. When set to `auto`, messages in the middle of the thread will be
// dropped to fit the context length of the model, `max_prompt_tokens`.
type BetaThreadNewAndRunParamsTruncationStrategyType = string

const (
	BetaThreadNewAndRunParamsTruncationStrategyTypeAuto         BetaThreadNewAndRunParamsTruncationStrategyType = "auto"
	BetaThreadNewAndRunParamsTruncationStrategyTypeLastMessages BetaThreadNewAndRunParamsTruncationStrategyType = "last_messages"
)

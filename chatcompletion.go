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
	"github.com/openai/openai-go/packages/ssestream"
	"github.com/openai/openai-go/shared"
	"github.com/openai/openai-go/shared/constant"
)

// ChatCompletionService contains methods and other services that help with
// interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewChatCompletionService] method instead.
type ChatCompletionService struct {
	Options  []option.RequestOption
	Messages ChatCompletionMessageService
}

// NewChatCompletionService generates a new service that applies the given options
// to each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewChatCompletionService(opts ...option.RequestOption) (r ChatCompletionService) {
	r = ChatCompletionService{}
	r.Options = opts
	r.Messages = NewChatCompletionMessageService(opts...)
	return
}

// Creates a model response for the given chat conversation. Learn more in the
// [text generation](https://platform.openai.com/docs/guides/text-generation),
// [vision](https://platform.openai.com/docs/guides/vision), and
// [audio](https://platform.openai.com/docs/guides/audio) guides.
//
// Parameter support can differ depending on the model used to generate the
// response, particularly for newer reasoning models. Parameters that are only
// supported for reasoning models are noted below. For the current state of
// unsupported parameters in reasoning models,
// [refer to the reasoning guide](https://platform.openai.com/docs/guides/reasoning).
func (r *ChatCompletionService) New(ctx context.Context, body ChatCompletionNewParams, opts ...option.RequestOption) (res *ChatCompletion, err error) {
	opts = append(r.Options[:], opts...)
	path := "chat/completions"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Creates a model response for the given chat conversation. Learn more in the
// [text generation](https://platform.openai.com/docs/guides/text-generation),
// [vision](https://platform.openai.com/docs/guides/vision), and
// [audio](https://platform.openai.com/docs/guides/audio) guides.
//
// Parameter support can differ depending on the model used to generate the
// response, particularly for newer reasoning models. Parameters that are only
// supported for reasoning models are noted below. For the current state of
// unsupported parameters in reasoning models,
// [refer to the reasoning guide](https://platform.openai.com/docs/guides/reasoning).
func (r *ChatCompletionService) NewStreaming(ctx context.Context, body ChatCompletionNewParams, opts ...option.RequestOption) (stream *ssestream.Stream[ChatCompletionChunk]) {
	var (
		raw *http.Response
		err error
	)
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithJSONSet("stream", true)}, opts...)
	path := "chat/completions"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &raw, opts...)
	return ssestream.NewStream[ChatCompletionChunk](ssestream.NewDecoder(raw), err)
}

// Get a stored chat completion. Only chat completions that have been created with
// the `store` parameter set to `true` will be returned.
func (r *ChatCompletionService) Get(ctx context.Context, completionID string, opts ...option.RequestOption) (res *ChatCompletion, err error) {
	opts = append(r.Options[:], opts...)
	if completionID == "" {
		err = errors.New("missing required completion_id parameter")
		return
	}
	path := fmt.Sprintf("chat/completions/%s", completionID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// Modify a stored chat completion. Only chat completions that have been created
// with the `store` parameter set to `true` can be modified. Currently, the only
// supported modification is to update the `metadata` field.
func (r *ChatCompletionService) Update(ctx context.Context, completionID string, body ChatCompletionUpdateParams, opts ...option.RequestOption) (res *ChatCompletion, err error) {
	opts = append(r.Options[:], opts...)
	if completionID == "" {
		err = errors.New("missing required completion_id parameter")
		return
	}
	path := fmt.Sprintf("chat/completions/%s", completionID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// List stored chat completions. Only chat completions that have been stored with
// the `store` parameter set to `true` will be returned.
func (r *ChatCompletionService) List(ctx context.Context, query ChatCompletionListParams, opts ...option.RequestOption) (res *pagination.CursorPage[ChatCompletion], err error) {
	var raw *http.Response
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "chat/completions"
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

// List stored chat completions. Only chat completions that have been stored with
// the `store` parameter set to `true` will be returned.
func (r *ChatCompletionService) ListAutoPaging(ctx context.Context, query ChatCompletionListParams, opts ...option.RequestOption) *pagination.CursorPageAutoPager[ChatCompletion] {
	return pagination.NewCursorPageAutoPager(r.List(ctx, query, opts...))
}

// Delete a stored chat completion. Only chat completions that have been created
// with the `store` parameter set to `true` can be deleted.
func (r *ChatCompletionService) Delete(ctx context.Context, completionID string, opts ...option.RequestOption) (res *ChatCompletionDeleted, err error) {
	opts = append(r.Options[:], opts...)
	if completionID == "" {
		err = errors.New("missing required completion_id parameter")
		return
	}
	path := fmt.Sprintf("chat/completions/%s", completionID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return
}

// Represents a chat completion response returned by model, based on the provided
// input.
type ChatCompletion struct {
	// A unique identifier for the chat completion.
	ID string `json:"id,omitzero,required"`
	// A list of chat completion choices. Can be more than one if `n` is greater
	// than 1.
	Choices []ChatCompletionChoice `json:"choices,omitzero,required"`
	// The Unix timestamp (in seconds) of when the chat completion was created.
	Created int64 `json:"created,omitzero,required"`
	// The model used for the chat completion.
	Model string `json:"model,omitzero,required"`
	// The object type, which is always `chat.completion`.
	//
	// This field can be elided, and will be automatically set as "chat.completion".
	Object constant.ChatCompletion `json:"object,required"`
	// The service tier used for processing the request.
	//
	// Any of "scale", "default"
	ServiceTier string `json:"service_tier,omitzero,nullable"`
	// This fingerprint represents the backend configuration that the model runs with.
	//
	// Can be used in conjunction with the `seed` request parameter to understand when
	// backend changes have been made that might impact determinism.
	SystemFingerprint string `json:"system_fingerprint,omitzero"`
	// Usage statistics for the completion request.
	Usage CompletionUsage `json:"usage,omitzero"`
	JSON  struct {
		ID                resp.Field
		Choices           resp.Field
		Created           resp.Field
		Model             resp.Field
		Object            resp.Field
		ServiceTier       resp.Field
		SystemFingerprint resp.Field
		Usage             resp.Field
		raw               string
	} `json:"-"`
}

func (r ChatCompletion) RawJSON() string { return r.JSON.raw }
func (r *ChatCompletion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ChatCompletionChoice struct {
	// The reason the model stopped generating tokens. This will be `stop` if the model
	// hit a natural stop point or a provided stop sequence, `length` if the maximum
	// number of tokens specified in the request was reached, `content_filter` if
	// content was omitted due to a flag from our content filters, `tool_calls` if the
	// model called a tool, or `function_call` (deprecated) if the model called a
	// function.
	//
	// Any of "stop", "length", "tool_calls", "content_filter", "function_call"
	FinishReason string `json:"finish_reason,omitzero,required"`
	// The index of the choice in the list of choices.
	Index int64 `json:"index,omitzero,required"`
	// Log probability information for the choice.
	Logprobs ChatCompletionChoicesLogprobs `json:"logprobs,omitzero,required,nullable"`
	// A chat completion message generated by the model.
	Message ChatCompletionMessage `json:"message,omitzero,required"`
	JSON    struct {
		FinishReason resp.Field
		Index        resp.Field
		Logprobs     resp.Field
		Message      resp.Field
		raw          string
	} `json:"-"`
}

func (r ChatCompletionChoice) RawJSON() string { return r.JSON.raw }
func (r *ChatCompletionChoice) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The reason the model stopped generating tokens. This will be `stop` if the model
// hit a natural stop point or a provided stop sequence, `length` if the maximum
// number of tokens specified in the request was reached, `content_filter` if
// content was omitted due to a flag from our content filters, `tool_calls` if the
// model called a tool, or `function_call` (deprecated) if the model called a
// function.
type ChatCompletionChoicesFinishReason = string

const (
	ChatCompletionChoicesFinishReasonStop          ChatCompletionChoicesFinishReason = "stop"
	ChatCompletionChoicesFinishReasonLength        ChatCompletionChoicesFinishReason = "length"
	ChatCompletionChoicesFinishReasonToolCalls     ChatCompletionChoicesFinishReason = "tool_calls"
	ChatCompletionChoicesFinishReasonContentFilter ChatCompletionChoicesFinishReason = "content_filter"
	ChatCompletionChoicesFinishReasonFunctionCall  ChatCompletionChoicesFinishReason = "function_call"
)

// Log probability information for the choice.
type ChatCompletionChoicesLogprobs struct {
	// A list of message content tokens with log probability information.
	Content []ChatCompletionTokenLogprob `json:"content,omitzero,required,nullable"`
	// A list of message refusal tokens with log probability information.
	Refusal []ChatCompletionTokenLogprob `json:"refusal,omitzero,required,nullable"`
	JSON    struct {
		Content resp.Field
		Refusal resp.Field
		raw     string
	} `json:"-"`
}

func (r ChatCompletionChoicesLogprobs) RawJSON() string { return r.JSON.raw }
func (r *ChatCompletionChoicesLogprobs) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The service tier used for processing the request.
type ChatCompletionServiceTier = string

const (
	ChatCompletionServiceTierScale   ChatCompletionServiceTier = "scale"
	ChatCompletionServiceTierDefault ChatCompletionServiceTier = "default"
)

// Messages sent by the model in response to user messages.
type ChatCompletionAssistantMessageParam struct {
	// The role of the messages author, in this case `assistant`.
	//
	// This field can be elided, and will be automatically set as "assistant".
	Role constant.Assistant `json:"role,required"`
	// Data about a previous audio response from the model.
	// [Learn more](https://platform.openai.com/docs/guides/audio).
	Audio ChatCompletionAssistantMessageParamAudio `json:"audio,omitzero"`
	// The contents of the assistant message. Required unless `tool_calls` or
	// `function_call` is specified.
	Content []ChatCompletionAssistantMessageParamContentUnion `json:"content,omitzero"`
	// Deprecated and replaced by `tool_calls`. The name and arguments of a function
	// that should be called, as generated by the model.
	//
	// Deprecated: deprecated
	FunctionCall ChatCompletionAssistantMessageParamFunctionCall `json:"function_call,omitzero"`
	// An optional name for the participant. Provides the model information to
	// differentiate between participants of the same role.
	Name param.String `json:"name,omitzero"`
	// The refusal message by the assistant.
	Refusal param.String `json:"refusal,omitzero"`
	// The tool calls generated by the model, such as function calls.
	ToolCalls []ChatCompletionMessageToolCallParam `json:"tool_calls,omitzero"`
	apiobject
}

func (f ChatCompletionAssistantMessageParam) IsMissing() bool {
	return param.IsOmitted(f) || f.IsNull()
}

func (r ChatCompletionAssistantMessageParam) MarshalJSON() (data []byte, err error) {
	type shadow ChatCompletionAssistantMessageParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// Data about a previous audio response from the model.
// [Learn more](https://platform.openai.com/docs/guides/audio).
type ChatCompletionAssistantMessageParamAudio struct {
	// Unique identifier for a previous audio response from the model.
	ID param.String `json:"id,omitzero,required"`
	apiobject
}

func (f ChatCompletionAssistantMessageParamAudio) IsMissing() bool {
	return param.IsOmitted(f) || f.IsNull()
}

func (r ChatCompletionAssistantMessageParamAudio) MarshalJSON() (data []byte, err error) {
	type shadow ChatCompletionAssistantMessageParamAudio
	return param.MarshalObject(r, (*shadow)(&r))
}

// Only one field can be non-zero
type ChatCompletionAssistantMessageParamContentUnion struct {
	OfText    *ChatCompletionContentPartTextParam
	OfRefusal *ChatCompletionContentPartRefusalParam
	apiunion
}

func (u ChatCompletionAssistantMessageParamContentUnion) IsMissing() bool {
	return param.IsOmitted(u) || u.IsNull()
}

func (u ChatCompletionAssistantMessageParamContentUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion[ChatCompletionAssistantMessageParamContentUnion](u.OfText, u.OfRefusal)
}

func (u ChatCompletionAssistantMessageParamContentUnion) GetText() *string {
	if vt := u.OfText; vt != nil && !vt.Text.IsOmitted() {
		return &vt.Text.V
	}
	return nil
}

func (u ChatCompletionAssistantMessageParamContentUnion) GetRefusal() *string {
	if vt := u.OfRefusal; vt != nil && !vt.Refusal.IsOmitted() {
		return &vt.Refusal.V
	}
	return nil
}

func (u ChatCompletionAssistantMessageParamContentUnion) GetType() *string {
	if vt := u.OfText; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfRefusal; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Deprecated and replaced by `tool_calls`. The name and arguments of a function
// that should be called, as generated by the model.
//
// Deprecated: deprecated
type ChatCompletionAssistantMessageParamFunctionCall struct {
	// The arguments to call the function with, as generated by the model in JSON
	// format. Note that the model does not always generate valid JSON, and may
	// hallucinate parameters not defined by your function schema. Validate the
	// arguments in your code before calling your function.
	Arguments param.String `json:"arguments,omitzero,required"`
	// The name of the function to call.
	Name param.String `json:"name,omitzero,required"`
	apiobject
}

func (f ChatCompletionAssistantMessageParamFunctionCall) IsMissing() bool {
	return param.IsOmitted(f) || f.IsNull()
}

func (r ChatCompletionAssistantMessageParamFunctionCall) MarshalJSON() (data []byte, err error) {
	type shadow ChatCompletionAssistantMessageParamFunctionCall
	return param.MarshalObject(r, (*shadow)(&r))
}

// If the audio output modality is requested, this object contains data about the
// audio response from the model.
// [Learn more](https://platform.openai.com/docs/guides/audio).
type ChatCompletionAudio struct {
	// Unique identifier for this audio response.
	ID string `json:"id,omitzero,required"`
	// Base64 encoded audio bytes generated by the model, in the format specified in
	// the request.
	Data string `json:"data,omitzero,required"`
	// The Unix timestamp (in seconds) for when this audio response will no longer be
	// accessible on the server for use in multi-turn conversations.
	ExpiresAt int64 `json:"expires_at,omitzero,required"`
	// Transcript of the audio generated by the model.
	Transcript string `json:"transcript,omitzero,required"`
	JSON       struct {
		ID         resp.Field
		Data       resp.Field
		ExpiresAt  resp.Field
		Transcript resp.Field
		raw        string
	} `json:"-"`
}

func (r ChatCompletionAudio) RawJSON() string { return r.JSON.raw }
func (r *ChatCompletionAudio) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Parameters for audio output. Required when audio output is requested with
// `modalities: ["audio"]`.
// [Learn more](https://platform.openai.com/docs/guides/audio).
type ChatCompletionAudioParam struct {
	// Specifies the output audio format. Must be one of `wav`, `mp3`, `flac`, `opus`,
	// or `pcm16`.
	//
	// Any of "wav", "mp3", "flac", "opus", "pcm16"
	Format string `json:"format,omitzero,required"`
	// The voice the model uses to respond. Supported voices are `ash`, `ballad`,
	// `coral`, `sage`, and `verse` (also supported but not recommended are `alloy`,
	// `echo`, and `shimmer`; these voices are less expressive).
	//
	// Any of "alloy", "ash", "ballad", "coral", "echo", "sage", "shimmer", "verse"
	Voice string `json:"voice,omitzero,required"`
	apiobject
}

func (f ChatCompletionAudioParam) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

func (r ChatCompletionAudioParam) MarshalJSON() (data []byte, err error) {
	type shadow ChatCompletionAudioParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// Specifies the output audio format. Must be one of `wav`, `mp3`, `flac`, `opus`,
// or `pcm16`.
type ChatCompletionAudioParamFormat = string

const (
	ChatCompletionAudioParamFormatWAV   ChatCompletionAudioParamFormat = "wav"
	ChatCompletionAudioParamFormatMP3   ChatCompletionAudioParamFormat = "mp3"
	ChatCompletionAudioParamFormatFLAC  ChatCompletionAudioParamFormat = "flac"
	ChatCompletionAudioParamFormatOpus  ChatCompletionAudioParamFormat = "opus"
	ChatCompletionAudioParamFormatPcm16 ChatCompletionAudioParamFormat = "pcm16"
)

// The voice the model uses to respond. Supported voices are `ash`, `ballad`,
// `coral`, `sage`, and `verse` (also supported but not recommended are `alloy`,
// `echo`, and `shimmer`; these voices are less expressive).
type ChatCompletionAudioParamVoice = string

const (
	ChatCompletionAudioParamVoiceAlloy   ChatCompletionAudioParamVoice = "alloy"
	ChatCompletionAudioParamVoiceAsh     ChatCompletionAudioParamVoice = "ash"
	ChatCompletionAudioParamVoiceBallad  ChatCompletionAudioParamVoice = "ballad"
	ChatCompletionAudioParamVoiceCoral   ChatCompletionAudioParamVoice = "coral"
	ChatCompletionAudioParamVoiceEcho    ChatCompletionAudioParamVoice = "echo"
	ChatCompletionAudioParamVoiceSage    ChatCompletionAudioParamVoice = "sage"
	ChatCompletionAudioParamVoiceShimmer ChatCompletionAudioParamVoice = "shimmer"
	ChatCompletionAudioParamVoiceVerse   ChatCompletionAudioParamVoice = "verse"
)

// Represents a streamed chunk of a chat completion response returned by model,
// based on the provided input.
type ChatCompletionChunk struct {
	// A unique identifier for the chat completion. Each chunk has the same ID.
	ID string `json:"id,omitzero,required"`
	// A list of chat completion choices. Can contain more than one elements if `n` is
	// greater than 1. Can also be empty for the last chunk if you set
	// `stream_options: {"include_usage": true}`.
	Choices []ChatCompletionChunkChoice `json:"choices,omitzero,required"`
	// The Unix timestamp (in seconds) of when the chat completion was created. Each
	// chunk has the same timestamp.
	Created int64 `json:"created,omitzero,required"`
	// The model to generate the completion.
	Model string `json:"model,omitzero,required"`
	// The object type, which is always `chat.completion.chunk`.
	//
	// This field can be elided, and will be automatically set as
	// "chat.completion.chunk".
	Object constant.ChatCompletionChunk `json:"object,required"`
	// The service tier used for processing the request.
	//
	// Any of "scale", "default"
	ServiceTier string `json:"service_tier,omitzero,nullable"`
	// This fingerprint represents the backend configuration that the model runs with.
	// Can be used in conjunction with the `seed` request parameter to understand when
	// backend changes have been made that might impact determinism.
	SystemFingerprint string `json:"system_fingerprint,omitzero"`
	// An optional field that will only be present when you set
	// `stream_options: {"include_usage": true}` in your request. When present, it
	// contains a null value except for the last chunk which contains the token usage
	// statistics for the entire request.
	Usage CompletionUsage `json:"usage,omitzero,nullable"`
	JSON  struct {
		ID                resp.Field
		Choices           resp.Field
		Created           resp.Field
		Model             resp.Field
		Object            resp.Field
		ServiceTier       resp.Field
		SystemFingerprint resp.Field
		Usage             resp.Field
		raw               string
	} `json:"-"`
}

func (r ChatCompletionChunk) RawJSON() string { return r.JSON.raw }
func (r *ChatCompletionChunk) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ChatCompletionChunkChoice struct {
	// A chat completion delta generated by streamed model responses.
	Delta ChatCompletionChunkChoicesDelta `json:"delta,omitzero,required"`
	// The reason the model stopped generating tokens. This will be `stop` if the model
	// hit a natural stop point or a provided stop sequence, `length` if the maximum
	// number of tokens specified in the request was reached, `content_filter` if
	// content was omitted due to a flag from our content filters, `tool_calls` if the
	// model called a tool, or `function_call` (deprecated) if the model called a
	// function.
	//
	// Any of "stop", "length", "tool_calls", "content_filter", "function_call"
	FinishReason string `json:"finish_reason,omitzero,required,nullable"`
	// The index of the choice in the list of choices.
	Index int64 `json:"index,omitzero,required"`
	// Log probability information for the choice.
	Logprobs ChatCompletionChunkChoicesLogprobs `json:"logprobs,omitzero,nullable"`
	JSON     struct {
		Delta        resp.Field
		FinishReason resp.Field
		Index        resp.Field
		Logprobs     resp.Field
		raw          string
	} `json:"-"`
}

func (r ChatCompletionChunkChoice) RawJSON() string { return r.JSON.raw }
func (r *ChatCompletionChunkChoice) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A chat completion delta generated by streamed model responses.
type ChatCompletionChunkChoicesDelta struct {
	// The contents of the chunk message.
	Content string `json:"content,omitzero,nullable"`
	// Deprecated and replaced by `tool_calls`. The name and arguments of a function
	// that should be called, as generated by the model.
	//
	// Deprecated: deprecated
	FunctionCall ChatCompletionChunkChoicesDeltaFunctionCall `json:"function_call,omitzero"`
	// The refusal message generated by the model.
	Refusal string `json:"refusal,omitzero,nullable"`
	// The role of the author of this message.
	//
	// Any of "developer", "system", "user", "assistant", "tool"
	Role      string                                    `json:"role,omitzero"`
	ToolCalls []ChatCompletionChunkChoicesDeltaToolCall `json:"tool_calls,omitzero"`
	JSON      struct {
		Content      resp.Field
		FunctionCall resp.Field
		Refusal      resp.Field
		Role         resp.Field
		ToolCalls    resp.Field
		raw          string
	} `json:"-"`
}

func (r ChatCompletionChunkChoicesDelta) RawJSON() string { return r.JSON.raw }
func (r *ChatCompletionChunkChoicesDelta) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Deprecated and replaced by `tool_calls`. The name and arguments of a function
// that should be called, as generated by the model.
//
// Deprecated: deprecated
type ChatCompletionChunkChoicesDeltaFunctionCall struct {
	// The arguments to call the function with, as generated by the model in JSON
	// format. Note that the model does not always generate valid JSON, and may
	// hallucinate parameters not defined by your function schema. Validate the
	// arguments in your code before calling your function.
	Arguments string `json:"arguments,omitzero"`
	// The name of the function to call.
	Name string `json:"name,omitzero"`
	JSON struct {
		Arguments resp.Field
		Name      resp.Field
		raw       string
	} `json:"-"`
}

func (r ChatCompletionChunkChoicesDeltaFunctionCall) RawJSON() string { return r.JSON.raw }
func (r *ChatCompletionChunkChoicesDeltaFunctionCall) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The role of the author of this message.
type ChatCompletionChunkChoicesDeltaRole = string

const (
	ChatCompletionChunkChoicesDeltaRoleDeveloper ChatCompletionChunkChoicesDeltaRole = "developer"
	ChatCompletionChunkChoicesDeltaRoleSystem    ChatCompletionChunkChoicesDeltaRole = "system"
	ChatCompletionChunkChoicesDeltaRoleUser      ChatCompletionChunkChoicesDeltaRole = "user"
	ChatCompletionChunkChoicesDeltaRoleAssistant ChatCompletionChunkChoicesDeltaRole = "assistant"
	ChatCompletionChunkChoicesDeltaRoleTool      ChatCompletionChunkChoicesDeltaRole = "tool"
)

type ChatCompletionChunkChoicesDeltaToolCall struct {
	Index int64 `json:"index,omitzero,required"`
	// The ID of the tool call.
	ID       string                                           `json:"id,omitzero"`
	Function ChatCompletionChunkChoicesDeltaToolCallsFunction `json:"function,omitzero"`
	// The type of the tool. Currently, only `function` is supported.
	//
	// Any of "function"
	Type string `json:"type"`
	JSON struct {
		Index    resp.Field
		ID       resp.Field
		Function resp.Field
		Type     resp.Field
		raw      string
	} `json:"-"`
}

func (r ChatCompletionChunkChoicesDeltaToolCall) RawJSON() string { return r.JSON.raw }
func (r *ChatCompletionChunkChoicesDeltaToolCall) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ChatCompletionChunkChoicesDeltaToolCallsFunction struct {
	// The arguments to call the function with, as generated by the model in JSON
	// format. Note that the model does not always generate valid JSON, and may
	// hallucinate parameters not defined by your function schema. Validate the
	// arguments in your code before calling your function.
	Arguments string `json:"arguments,omitzero"`
	// The name of the function to call.
	Name string `json:"name,omitzero"`
	JSON struct {
		Arguments resp.Field
		Name      resp.Field
		raw       string
	} `json:"-"`
}

func (r ChatCompletionChunkChoicesDeltaToolCallsFunction) RawJSON() string { return r.JSON.raw }
func (r *ChatCompletionChunkChoicesDeltaToolCallsFunction) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The type of the tool. Currently, only `function` is supported.
type ChatCompletionChunkChoicesDeltaToolCallsType = string

const (
	ChatCompletionChunkChoicesDeltaToolCallsTypeFunction ChatCompletionChunkChoicesDeltaToolCallsType = "function"
)

// The reason the model stopped generating tokens. This will be `stop` if the model
// hit a natural stop point or a provided stop sequence, `length` if the maximum
// number of tokens specified in the request was reached, `content_filter` if
// content was omitted due to a flag from our content filters, `tool_calls` if the
// model called a tool, or `function_call` (deprecated) if the model called a
// function.
type ChatCompletionChunkChoicesFinishReason = string

const (
	ChatCompletionChunkChoicesFinishReasonStop          ChatCompletionChunkChoicesFinishReason = "stop"
	ChatCompletionChunkChoicesFinishReasonLength        ChatCompletionChunkChoicesFinishReason = "length"
	ChatCompletionChunkChoicesFinishReasonToolCalls     ChatCompletionChunkChoicesFinishReason = "tool_calls"
	ChatCompletionChunkChoicesFinishReasonContentFilter ChatCompletionChunkChoicesFinishReason = "content_filter"
	ChatCompletionChunkChoicesFinishReasonFunctionCall  ChatCompletionChunkChoicesFinishReason = "function_call"
)

// Log probability information for the choice.
type ChatCompletionChunkChoicesLogprobs struct {
	// A list of message content tokens with log probability information.
	Content []ChatCompletionTokenLogprob `json:"content,omitzero,required,nullable"`
	// A list of message refusal tokens with log probability information.
	Refusal []ChatCompletionTokenLogprob `json:"refusal,omitzero,required,nullable"`
	JSON    struct {
		Content resp.Field
		Refusal resp.Field
		raw     string
	} `json:"-"`
}

func (r ChatCompletionChunkChoicesLogprobs) RawJSON() string { return r.JSON.raw }
func (r *ChatCompletionChunkChoicesLogprobs) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The service tier used for processing the request.
type ChatCompletionChunkServiceTier = string

const (
	ChatCompletionChunkServiceTierScale   ChatCompletionChunkServiceTier = "scale"
	ChatCompletionChunkServiceTierDefault ChatCompletionChunkServiceTier = "default"
)

func NewChatCompletionContentPartOfText(text string) ChatCompletionContentPartUnionParam {
	var variant ChatCompletionContentPartTextParam
	variant.Text = newString(text)
	return ChatCompletionContentPartUnionParam{OfText: &variant}
}

func NewChatCompletionContentPartOfImageURL(imageURL ChatCompletionContentPartImageImageURLParam) ChatCompletionContentPartUnionParam {
	var image_url ChatCompletionContentPartImageParam
	image_url.ImageURL = imageURL
	return ChatCompletionContentPartUnionParam{OfImageURL: &image_url}
}

func NewChatCompletionContentPartOfInputAudio(inputAudio ChatCompletionContentPartInputAudioInputAudioParam) ChatCompletionContentPartUnionParam {
	var input_audio ChatCompletionContentPartInputAudioParam
	input_audio.InputAudio = inputAudio
	return ChatCompletionContentPartUnionParam{OfInputAudio: &input_audio}
}

// Only one field can be non-zero
type ChatCompletionContentPartUnionParam struct {
	OfText       *ChatCompletionContentPartTextParam
	OfImageURL   *ChatCompletionContentPartImageParam
	OfInputAudio *ChatCompletionContentPartInputAudioParam
	apiunion
}

func (u ChatCompletionContentPartUnionParam) IsMissing() bool {
	return param.IsOmitted(u) || u.IsNull()
}

func (u ChatCompletionContentPartUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion[ChatCompletionContentPartUnionParam](u.OfText, u.OfImageURL, u.OfInputAudio)
}

func (u ChatCompletionContentPartUnionParam) GetText() *string {
	if vt := u.OfText; vt != nil && !vt.Text.IsOmitted() {
		return &vt.Text.V
	}
	return nil
}

func (u ChatCompletionContentPartUnionParam) GetImageURL() *ChatCompletionContentPartImageImageURLParam {
	if vt := u.OfImageURL; vt != nil {
		return &vt.ImageURL
	}
	return nil
}

func (u ChatCompletionContentPartUnionParam) GetInputAudio() *ChatCompletionContentPartInputAudioInputAudioParam {
	if vt := u.OfInputAudio; vt != nil {
		return &vt.InputAudio
	}
	return nil
}

func (u ChatCompletionContentPartUnionParam) GetType() *string {
	if vt := u.OfText; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfImageURL; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfInputAudio; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Learn about [image inputs](https://platform.openai.com/docs/guides/vision).
type ChatCompletionContentPartImageParam struct {
	ImageURL ChatCompletionContentPartImageImageURLParam `json:"image_url,omitzero,required"`
	// The type of the content part.
	//
	// This field can be elided, and will be automatically set as "image_url".
	Type constant.ImageURL `json:"type,required"`
	apiobject
}

func (f ChatCompletionContentPartImageParam) IsMissing() bool {
	return param.IsOmitted(f) || f.IsNull()
}

func (r ChatCompletionContentPartImageParam) MarshalJSON() (data []byte, err error) {
	type shadow ChatCompletionContentPartImageParam
	return param.MarshalObject(r, (*shadow)(&r))
}

type ChatCompletionContentPartImageImageURLParam struct {
	// Either a URL of the image or the base64 encoded image data.
	URL param.String `json:"url,omitzero,required" format:"uri"`
	// Specifies the detail level of the image. Learn more in the
	// [Vision guide](https://platform.openai.com/docs/guides/vision#low-or-high-fidelity-image-understanding).
	//
	// Any of "auto", "low", "high"
	Detail string `json:"detail,omitzero"`
	apiobject
}

func (f ChatCompletionContentPartImageImageURLParam) IsMissing() bool {
	return param.IsOmitted(f) || f.IsNull()
}

func (r ChatCompletionContentPartImageImageURLParam) MarshalJSON() (data []byte, err error) {
	type shadow ChatCompletionContentPartImageImageURLParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// Specifies the detail level of the image. Learn more in the
// [Vision guide](https://platform.openai.com/docs/guides/vision#low-or-high-fidelity-image-understanding).
type ChatCompletionContentPartImageImageURLDetail = string

const (
	ChatCompletionContentPartImageImageURLDetailAuto ChatCompletionContentPartImageImageURLDetail = "auto"
	ChatCompletionContentPartImageImageURLDetailLow  ChatCompletionContentPartImageImageURLDetail = "low"
	ChatCompletionContentPartImageImageURLDetailHigh ChatCompletionContentPartImageImageURLDetail = "high"
)

// Learn about [audio inputs](https://platform.openai.com/docs/guides/audio).
type ChatCompletionContentPartInputAudioParam struct {
	InputAudio ChatCompletionContentPartInputAudioInputAudioParam `json:"input_audio,omitzero,required"`
	// The type of the content part. Always `input_audio`.
	//
	// This field can be elided, and will be automatically set as "input_audio".
	Type constant.InputAudio `json:"type,required"`
	apiobject
}

func (f ChatCompletionContentPartInputAudioParam) IsMissing() bool {
	return param.IsOmitted(f) || f.IsNull()
}

func (r ChatCompletionContentPartInputAudioParam) MarshalJSON() (data []byte, err error) {
	type shadow ChatCompletionContentPartInputAudioParam
	return param.MarshalObject(r, (*shadow)(&r))
}

type ChatCompletionContentPartInputAudioInputAudioParam struct {
	// Base64 encoded audio data.
	Data param.String `json:"data,omitzero,required"`
	// The format of the encoded audio data. Currently supports "wav" and "mp3".
	//
	// Any of "wav", "mp3"
	Format string `json:"format,omitzero,required"`
	apiobject
}

func (f ChatCompletionContentPartInputAudioInputAudioParam) IsMissing() bool {
	return param.IsOmitted(f) || f.IsNull()
}

func (r ChatCompletionContentPartInputAudioInputAudioParam) MarshalJSON() (data []byte, err error) {
	type shadow ChatCompletionContentPartInputAudioInputAudioParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// The format of the encoded audio data. Currently supports "wav" and "mp3".
type ChatCompletionContentPartInputAudioInputAudioFormat = string

const (
	ChatCompletionContentPartInputAudioInputAudioFormatWAV ChatCompletionContentPartInputAudioInputAudioFormat = "wav"
	ChatCompletionContentPartInputAudioInputAudioFormatMP3 ChatCompletionContentPartInputAudioInputAudioFormat = "mp3"
)

type ChatCompletionContentPartRefusalParam struct {
	// The refusal message generated by the model.
	Refusal param.String `json:"refusal,omitzero,required"`
	// The type of the content part.
	//
	// This field can be elided, and will be automatically set as "refusal".
	Type constant.Refusal `json:"type,required"`
	apiobject
}

func (f ChatCompletionContentPartRefusalParam) IsMissing() bool {
	return param.IsOmitted(f) || f.IsNull()
}

func (r ChatCompletionContentPartRefusalParam) MarshalJSON() (data []byte, err error) {
	type shadow ChatCompletionContentPartRefusalParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// Learn about
// [text inputs](https://platform.openai.com/docs/guides/text-generation).
type ChatCompletionContentPartTextParam struct {
	// The text content.
	Text param.String `json:"text,omitzero,required"`
	// The type of the content part.
	//
	// This field can be elided, and will be automatically set as "text".
	Type constant.Text `json:"type,required"`
	apiobject
}

func (f ChatCompletionContentPartTextParam) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

func (r ChatCompletionContentPartTextParam) MarshalJSON() (data []byte, err error) {
	type shadow ChatCompletionContentPartTextParam
	return param.MarshalObject(r, (*shadow)(&r))
}

type ChatCompletionDeleted struct {
	// The ID of the chat completion that was deleted.
	ID string `json:"id,omitzero,required"`
	// Whether the chat completion was deleted.
	Deleted bool `json:"deleted,omitzero,required"`
	// The type of object being deleted.
	//
	// This field can be elided, and will be automatically set as
	// "chat.completion.deleted".
	Object constant.ChatCompletionDeleted `json:"object,required"`
	JSON   struct {
		ID      resp.Field
		Deleted resp.Field
		Object  resp.Field
		raw     string
	} `json:"-"`
}

func (r ChatCompletionDeleted) RawJSON() string { return r.JSON.raw }
func (r *ChatCompletionDeleted) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Developer-provided instructions that the model should follow, regardless of
// messages sent by the user. With o1 models and newer, `developer` messages
// replace the previous `system` messages.
type ChatCompletionDeveloperMessageParam struct {
	// The contents of the developer message.
	Content []ChatCompletionContentPartTextParam `json:"content,omitzero,required"`
	// The role of the messages author, in this case `developer`.
	//
	// This field can be elided, and will be automatically set as "developer".
	Role constant.Developer `json:"role,required"`
	// An optional name for the participant. Provides the model information to
	// differentiate between participants of the same role.
	Name param.String `json:"name,omitzero"`
	apiobject
}

func (f ChatCompletionDeveloperMessageParam) IsMissing() bool {
	return param.IsOmitted(f) || f.IsNull()
}

func (r ChatCompletionDeveloperMessageParam) MarshalJSON() (data []byte, err error) {
	type shadow ChatCompletionDeveloperMessageParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// Specifying a particular function via `{"name": "my_function"}` forces the model
// to call that function.
type ChatCompletionFunctionCallOptionParam struct {
	// The name of the function to call.
	Name param.String `json:"name,omitzero,required"`
	apiobject
}

func (f ChatCompletionFunctionCallOptionParam) IsMissing() bool {
	return param.IsOmitted(f) || f.IsNull()
}

func (r ChatCompletionFunctionCallOptionParam) MarshalJSON() (data []byte, err error) {
	type shadow ChatCompletionFunctionCallOptionParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// Deprecated: deprecated
type ChatCompletionFunctionMessageParam struct {
	// The contents of the function message.
	Content param.String `json:"content,omitzero,required"`
	// The name of the function to call.
	Name param.String `json:"name,omitzero,required"`
	// The role of the messages author, in this case `function`.
	//
	// This field can be elided, and will be automatically set as "function".
	Role constant.Function `json:"role,required"`
	apiobject
}

func (f ChatCompletionFunctionMessageParam) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

func (r ChatCompletionFunctionMessageParam) MarshalJSON() (data []byte, err error) {
	type shadow ChatCompletionFunctionMessageParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// A chat completion message generated by the model.
type ChatCompletionMessage struct {
	// The contents of the message.
	Content string `json:"content,omitzero,required,nullable"`
	// The refusal message generated by the model.
	Refusal string `json:"refusal,omitzero,required,nullable"`
	// The role of the author of this message.
	//
	// This field can be elided, and will be automatically set as "assistant".
	Role constant.Assistant `json:"role,required"`
	// If the audio output modality is requested, this object contains data about the
	// audio response from the model.
	// [Learn more](https://platform.openai.com/docs/guides/audio).
	Audio ChatCompletionAudio `json:"audio,omitzero,nullable"`
	// Deprecated and replaced by `tool_calls`. The name and arguments of a function
	// that should be called, as generated by the model.
	//
	// Deprecated: deprecated
	FunctionCall ChatCompletionMessageFunctionCall `json:"function_call,omitzero"`
	// The tool calls generated by the model, such as function calls.
	ToolCalls []ChatCompletionMessageToolCall `json:"tool_calls,omitzero"`
	JSON      struct {
		Content      resp.Field
		Refusal      resp.Field
		Role         resp.Field
		Audio        resp.Field
		FunctionCall resp.Field
		ToolCalls    resp.Field
		raw          string
	} `json:"-"`
}

func (r ChatCompletionMessage) RawJSON() string { return r.JSON.raw }
func (r *ChatCompletionMessage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (r ChatCompletionMessage) ToParam() ChatCompletionAssistantMessageParam {
	var p ChatCompletionAssistantMessageParam
	p.Audio.ID = toParamString(r.Audio.ID, r.Audio.JSON.ID)
	p.FunctionCall.Arguments = toParamString(r.FunctionCall.Arguments, r.FunctionCall.JSON.Arguments)
	p.FunctionCall.Name = toParamString(r.FunctionCall.Name, r.FunctionCall.JSON.Name)
	p.Refusal = toParamString(r.Refusal, r.JSON.Refusal)
	_ = p.ToolCalls
	_ = r.ToolCalls
	return p
}

// Deprecated and replaced by `tool_calls`. The name and arguments of a function
// that should be called, as generated by the model.
//
// Deprecated: deprecated
type ChatCompletionMessageFunctionCall struct {
	// The arguments to call the function with, as generated by the model in JSON
	// format. Note that the model does not always generate valid JSON, and may
	// hallucinate parameters not defined by your function schema. Validate the
	// arguments in your code before calling your function.
	Arguments string `json:"arguments,omitzero,required"`
	// The name of the function to call.
	Name string `json:"name,omitzero,required"`
	JSON struct {
		Arguments resp.Field
		Name      resp.Field
		raw       string
	} `json:"-"`
}

func (r ChatCompletionMessageFunctionCall) RawJSON() string { return r.JSON.raw }
func (r *ChatCompletionMessageFunctionCall) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewChatCompletionMessageParamOfDeveloper(content []ChatCompletionContentPartTextParam) ChatCompletionMessageParamUnion {
	var developer ChatCompletionDeveloperMessageParam
	developer.Content = content
	return ChatCompletionMessageParamUnion{OfDeveloper: &developer}
}

func NewChatCompletionMessageParamOfSystem(content []ChatCompletionContentPartTextParam) ChatCompletionMessageParamUnion {
	var system ChatCompletionSystemMessageParam
	system.Content = content
	return ChatCompletionMessageParamUnion{OfSystem: &system}
}

func NewChatCompletionMessageParamOfUser(content []ChatCompletionContentPartUnionParam) ChatCompletionMessageParamUnion {
	var user ChatCompletionUserMessageParam
	user.Content = content
	return ChatCompletionMessageParamUnion{OfUser: &user}
}

func NewChatCompletionMessageParamOfTool(content []ChatCompletionContentPartTextParam, toolCallID string) ChatCompletionMessageParamUnion {
	var tool ChatCompletionToolMessageParam
	tool.Content = content
	tool.ToolCallID = newString(toolCallID)
	return ChatCompletionMessageParamUnion{OfTool: &tool}
}

func NewChatCompletionMessageParamOfFunction(content string, name string) ChatCompletionMessageParamUnion {
	var function ChatCompletionFunctionMessageParam
	function.Content = newString(content)
	function.Name = newString(name)
	return ChatCompletionMessageParamUnion{OfFunction: &function}
}

// Only one field can be non-zero
type ChatCompletionMessageParamUnion struct {
	OfDeveloper *ChatCompletionDeveloperMessageParam
	OfSystem    *ChatCompletionSystemMessageParam
	OfUser      *ChatCompletionUserMessageParam
	OfAssistant *ChatCompletionAssistantMessageParam
	OfTool      *ChatCompletionToolMessageParam
	OfFunction  *ChatCompletionFunctionMessageParam
	apiunion
}

func (u ChatCompletionMessageParamUnion) IsMissing() bool { return param.IsOmitted(u) || u.IsNull() }

func (u ChatCompletionMessageParamUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion[ChatCompletionMessageParamUnion](u.OfDeveloper, u.OfSystem, u.OfUser, u.OfAssistant, u.OfTool, u.OfFunction)
}

func (u ChatCompletionMessageParamUnion) GetName() *string {
	if vt := u.OfDeveloper; vt != nil && !vt.Name.IsOmitted() {
		return &vt.Name.V
	} else if vt := u.OfSystem; vt != nil && !vt.Name.IsOmitted() {
		return &vt.Name.V
	} else if vt := u.OfUser; vt != nil && !vt.Name.IsOmitted() {
		return &vt.Name.V
	} else if vt := u.OfAssistant; vt != nil && !vt.Name.IsOmitted() {
		return &vt.Name.V
	} else if vt := u.OfFunction; vt != nil && !vt.Name.IsOmitted() {
		return &vt.Name.V
	}
	return nil
}

func (u ChatCompletionMessageParamUnion) GetAudio() *ChatCompletionAssistantMessageParamAudio {
	if vt := u.OfAssistant; vt != nil {
		return &vt.Audio
	}
	return nil
}

func (u ChatCompletionMessageParamUnion) GetFunctionCall() *ChatCompletionAssistantMessageParamFunctionCall {
	if vt := u.OfAssistant; vt != nil {
		return &vt.FunctionCall
	}
	return nil
}

func (u ChatCompletionMessageParamUnion) GetRefusal() *string {
	if vt := u.OfAssistant; vt != nil && !vt.Refusal.IsOmitted() {
		return &vt.Refusal.V
	}
	return nil
}

func (u ChatCompletionMessageParamUnion) GetToolCalls() []ChatCompletionMessageToolCallParam {
	if vt := u.OfAssistant; vt != nil {
		return vt.ToolCalls
	}
	return nil
}

func (u ChatCompletionMessageParamUnion) GetToolCallID() *string {
	if vt := u.OfTool; vt != nil && !vt.ToolCallID.IsOmitted() {
		return &vt.ToolCallID.V
	}
	return nil
}

func (u ChatCompletionMessageParamUnion) GetRole() *string {
	if vt := u.OfDeveloper; vt != nil {
		return (*string)(&vt.Role)
	} else if vt := u.OfSystem; vt != nil {
		return (*string)(&vt.Role)
	} else if vt := u.OfUser; vt != nil {
		return (*string)(&vt.Role)
	} else if vt := u.OfAssistant; vt != nil {
		return (*string)(&vt.Role)
	} else if vt := u.OfTool; vt != nil {
		return (*string)(&vt.Role)
	} else if vt := u.OfFunction; vt != nil {
		return (*string)(&vt.Role)
	}
	return nil
}

func (u ChatCompletionMessageParamUnion) GetContent() (res chatCompletionMessageParamUnionContent) {
	if vt := u.OfDeveloper; vt != nil {
		res.OfChatCompletionDeveloperMessageContent = &vt.Content
	} else if vt := u.OfSystem; vt != nil {
		res.OfChatCompletionDeveloperMessageContent = &vt.Content
	} else if vt := u.OfTool; vt != nil {
		res.OfChatCompletionDeveloperMessageContent = &vt.Content
	} else if vt := u.OfUser; vt != nil {
		res.OfChatCompletionUserMessageContent = &vt.Content
	} else if vt := u.OfAssistant; vt != nil {
		res.OfChatCompletionAssistantMessageContent = &vt.Content
	} else if vt := u.OfFunction; vt != nil {
		res.OfString = &vt.Content
	}
	return
}

// Only one field can be non-zero
type chatCompletionMessageParamUnionContent struct {
	OfChatCompletionDeveloperMessageContent *[]ChatCompletionContentPartTextParam
	OfChatCompletionUserMessageContent      *[]ChatCompletionContentPartUnionParam
	OfChatCompletionAssistantMessageContent *[]ChatCompletionAssistantMessageParamContentUnion
	OfString                                *param.String
}

type ChatCompletionMessageToolCall struct {
	// The ID of the tool call.
	ID string `json:"id,omitzero,required"`
	// The function that the model called.
	Function ChatCompletionMessageToolCallFunction `json:"function,omitzero,required"`
	// The type of the tool. Currently, only `function` is supported.
	//
	// This field can be elided, and will be automatically set as "function".
	Type constant.Function `json:"type,required"`
	JSON struct {
		ID       resp.Field
		Function resp.Field
		Type     resp.Field
		raw      string
	} `json:"-"`
}

func (r ChatCompletionMessageToolCall) RawJSON() string { return r.JSON.raw }
func (r *ChatCompletionMessageToolCall) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this ChatCompletionMessageToolCall to a
// ChatCompletionMessageToolCallParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// ChatCompletionMessageToolCallParam.IsOverridden()
func (r ChatCompletionMessageToolCall) ToParam() ChatCompletionMessageToolCallParam {
	return param.Override[ChatCompletionMessageToolCallParam](r.RawJSON())
}

// The function that the model called.
type ChatCompletionMessageToolCallFunction struct {
	// The arguments to call the function with, as generated by the model in JSON
	// format. Note that the model does not always generate valid JSON, and may
	// hallucinate parameters not defined by your function schema. Validate the
	// arguments in your code before calling your function.
	Arguments string `json:"arguments,omitzero,required"`
	// The name of the function to call.
	Name string `json:"name,omitzero,required"`
	JSON struct {
		Arguments resp.Field
		Name      resp.Field
		raw       string
	} `json:"-"`
}

func (r ChatCompletionMessageToolCallFunction) RawJSON() string { return r.JSON.raw }
func (r *ChatCompletionMessageToolCallFunction) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ChatCompletionMessageToolCallParam struct {
	// The ID of the tool call.
	ID param.String `json:"id,omitzero,required"`
	// The function that the model called.
	Function ChatCompletionMessageToolCallFunctionParam `json:"function,omitzero,required"`
	// The type of the tool. Currently, only `function` is supported.
	//
	// This field can be elided, and will be automatically set as "function".
	Type constant.Function `json:"type,required"`
	apiobject
}

func (f ChatCompletionMessageToolCallParam) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

func (r ChatCompletionMessageToolCallParam) MarshalJSON() (data []byte, err error) {
	type shadow ChatCompletionMessageToolCallParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// The function that the model called.
type ChatCompletionMessageToolCallFunctionParam struct {
	// The arguments to call the function with, as generated by the model in JSON
	// format. Note that the model does not always generate valid JSON, and may
	// hallucinate parameters not defined by your function schema. Validate the
	// arguments in your code before calling your function.
	Arguments param.String `json:"arguments,omitzero,required"`
	// The name of the function to call.
	Name param.String `json:"name,omitzero,required"`
	apiobject
}

func (f ChatCompletionMessageToolCallFunctionParam) IsMissing() bool {
	return param.IsOmitted(f) || f.IsNull()
}

func (r ChatCompletionMessageToolCallFunctionParam) MarshalJSON() (data []byte, err error) {
	type shadow ChatCompletionMessageToolCallFunctionParam
	return param.MarshalObject(r, (*shadow)(&r))
}

type ChatCompletionModality string

const (
	ChatCompletionModalityText  ChatCompletionModality = "text"
	ChatCompletionModalityAudio ChatCompletionModality = "audio"
)

// Specifies a tool the model should use. Use to force the model to call a specific
// function.
type ChatCompletionNamedToolChoiceParam struct {
	Function ChatCompletionNamedToolChoiceFunctionParam `json:"function,omitzero,required"`
	// The type of the tool. Currently, only `function` is supported.
	//
	// This field can be elided, and will be automatically set as "function".
	Type constant.Function `json:"type,required"`
	apiobject
}

func (f ChatCompletionNamedToolChoiceParam) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

func (r ChatCompletionNamedToolChoiceParam) MarshalJSON() (data []byte, err error) {
	type shadow ChatCompletionNamedToolChoiceParam
	return param.MarshalObject(r, (*shadow)(&r))
}

type ChatCompletionNamedToolChoiceFunctionParam struct {
	// The name of the function to call.
	Name param.String `json:"name,omitzero,required"`
	apiobject
}

func (f ChatCompletionNamedToolChoiceFunctionParam) IsMissing() bool {
	return param.IsOmitted(f) || f.IsNull()
}

func (r ChatCompletionNamedToolChoiceFunctionParam) MarshalJSON() (data []byte, err error) {
	type shadow ChatCompletionNamedToolChoiceFunctionParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// Static predicted output content, such as the content of a text file that is
// being regenerated.
type ChatCompletionPredictionContentParam struct {
	// The content that should be matched when generating a model response. If
	// generated tokens would match this content, the entire model response can be
	// returned much more quickly.
	Content []ChatCompletionContentPartTextParam `json:"content,omitzero,required"`
	// The type of the predicted content you want to provide. This type is currently
	// always `content`.
	//
	// This field can be elided, and will be automatically set as "content".
	Type constant.Content `json:"type,required"`
	apiobject
}

func (f ChatCompletionPredictionContentParam) IsMissing() bool {
	return param.IsOmitted(f) || f.IsNull()
}

func (r ChatCompletionPredictionContentParam) MarshalJSON() (data []byte, err error) {
	type shadow ChatCompletionPredictionContentParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// **o1 and o3-mini models only**
//
// Constrains effort on reasoning for
// [reasoning models](https://platform.openai.com/docs/guides/reasoning). Currently
// supported values are `low`, `medium`, and `high`. Reducing reasoning effort can
// result in faster responses and fewer tokens used on reasoning in a response.
type ChatCompletionReasoningEffort string

const (
	ChatCompletionReasoningEffortLow    ChatCompletionReasoningEffort = "low"
	ChatCompletionReasoningEffortMedium ChatCompletionReasoningEffort = "medium"
	ChatCompletionReasoningEffortHigh   ChatCompletionReasoningEffort = "high"
)

// A chat completion message generated by the model.
type ChatCompletionStoreMessage struct {
	// The identifier of the chat message.
	ID   string `json:"id,omitzero,required"`
	JSON struct {
		ID  resp.Field
		raw string
	} `json:"-"`
	ChatCompletionMessage
}

func (r ChatCompletionStoreMessage) RawJSON() string { return r.JSON.raw }
func (r *ChatCompletionStoreMessage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Options for streaming response. Only set this when you set `stream: true`.
type ChatCompletionStreamOptionsParam struct {
	// If set, an additional chunk will be streamed before the `data: [DONE]` message.
	// The `usage` field on this chunk shows the token usage statistics for the entire
	// request, and the `choices` field will always be an empty array. All other chunks
	// will also include a `usage` field, but with a null value.
	IncludeUsage param.Bool `json:"include_usage,omitzero"`
	apiobject
}

func (f ChatCompletionStreamOptionsParam) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

func (r ChatCompletionStreamOptionsParam) MarshalJSON() (data []byte, err error) {
	type shadow ChatCompletionStreamOptionsParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// Developer-provided instructions that the model should follow, regardless of
// messages sent by the user. With o1 models and newer, use `developer` messages
// for this purpose instead.
type ChatCompletionSystemMessageParam struct {
	// The contents of the system message.
	Content []ChatCompletionContentPartTextParam `json:"content,omitzero,required"`
	// The role of the messages author, in this case `system`.
	//
	// This field can be elided, and will be automatically set as "system".
	Role constant.System `json:"role,required"`
	// An optional name for the participant. Provides the model information to
	// differentiate between participants of the same role.
	Name param.String `json:"name,omitzero"`
	apiobject
}

func (f ChatCompletionSystemMessageParam) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

func (r ChatCompletionSystemMessageParam) MarshalJSON() (data []byte, err error) {
	type shadow ChatCompletionSystemMessageParam
	return param.MarshalObject(r, (*shadow)(&r))
}

type ChatCompletionTokenLogprob struct {
	// The token.
	Token string `json:"token,omitzero,required"`
	// A list of integers representing the UTF-8 bytes representation of the token.
	// Useful in instances where characters are represented by multiple tokens and
	// their byte representations must be combined to generate the correct text
	// representation. Can be `null` if there is no bytes representation for the token.
	Bytes []int64 `json:"bytes,omitzero,required,nullable"`
	// The log probability of this token, if it is within the top 20 most likely
	// tokens. Otherwise, the value `-9999.0` is used to signify that the token is very
	// unlikely.
	Logprob float64 `json:"logprob,omitzero,required"`
	// List of the most likely tokens and their log probability, at this token
	// position. In rare cases, there may be fewer than the number of requested
	// `top_logprobs` returned.
	TopLogprobs []ChatCompletionTokenLogprobTopLogprob `json:"top_logprobs,omitzero,required"`
	JSON        struct {
		Token       resp.Field
		Bytes       resp.Field
		Logprob     resp.Field
		TopLogprobs resp.Field
		raw         string
	} `json:"-"`
}

func (r ChatCompletionTokenLogprob) RawJSON() string { return r.JSON.raw }
func (r *ChatCompletionTokenLogprob) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ChatCompletionTokenLogprobTopLogprob struct {
	// The token.
	Token string `json:"token,omitzero,required"`
	// A list of integers representing the UTF-8 bytes representation of the token.
	// Useful in instances where characters are represented by multiple tokens and
	// their byte representations must be combined to generate the correct text
	// representation. Can be `null` if there is no bytes representation for the token.
	Bytes []int64 `json:"bytes,omitzero,required,nullable"`
	// The log probability of this token, if it is within the top 20 most likely
	// tokens. Otherwise, the value `-9999.0` is used to signify that the token is very
	// unlikely.
	Logprob float64 `json:"logprob,omitzero,required"`
	JSON    struct {
		Token   resp.Field
		Bytes   resp.Field
		Logprob resp.Field
		raw     string
	} `json:"-"`
}

func (r ChatCompletionTokenLogprobTopLogprob) RawJSON() string { return r.JSON.raw }
func (r *ChatCompletionTokenLogprobTopLogprob) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ChatCompletionToolParam struct {
	Function shared.FunctionDefinitionParam `json:"function,omitzero,required"`
	// The type of the tool. Currently, only `function` is supported.
	//
	// This field can be elided, and will be automatically set as "function".
	Type constant.Function `json:"type,required"`
	apiobject
}

func (f ChatCompletionToolParam) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

func (r ChatCompletionToolParam) MarshalJSON() (data []byte, err error) {
	type shadow ChatCompletionToolParam
	return param.MarshalObject(r, (*shadow)(&r))
}

func NewChatCompletionToolChoiceOptionOfChatCompletionNamedToolChoice(function ChatCompletionNamedToolChoiceFunctionParam) ChatCompletionToolChoiceOptionUnionParam {
	var variant ChatCompletionNamedToolChoiceParam
	variant.Function = function
	return ChatCompletionToolChoiceOptionUnionParam{OfChatCompletionNamedToolChoice: &variant}
}

// Only one field can be non-zero
type ChatCompletionToolChoiceOptionUnionParam struct {
	// Check if union is this variant with !param.IsOmitted(union.OfAuto)
	OfAuto                          string
	OfChatCompletionNamedToolChoice *ChatCompletionNamedToolChoiceParam
	apiunion
}

func (u ChatCompletionToolChoiceOptionUnionParam) IsMissing() bool {
	return param.IsOmitted(u) || u.IsNull()
}

func (u ChatCompletionToolChoiceOptionUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion[ChatCompletionToolChoiceOptionUnionParam](u.OfAuto, u.OfChatCompletionNamedToolChoice)
}

func (u ChatCompletionToolChoiceOptionUnionParam) GetFunction() *ChatCompletionNamedToolChoiceFunctionParam {
	if vt := u.OfChatCompletionNamedToolChoice; vt != nil {
		return &vt.Function
	}
	return nil
}

func (u ChatCompletionToolChoiceOptionUnionParam) GetType() *constant.Function {
	if vt := u.OfChatCompletionNamedToolChoice; vt != nil {
		return &vt.Type
	}
	return nil
}

// `none` means the model will not call any tool and instead generates a message.
// `auto` means the model can pick between generating a message or calling one or
// more tools. `required` means the model must call one or more tools.
type ChatCompletionToolChoiceOptionAuto = string

const (
	ChatCompletionToolChoiceOptionAutoNone     ChatCompletionToolChoiceOptionAuto = "none"
	ChatCompletionToolChoiceOptionAutoAuto     ChatCompletionToolChoiceOptionAuto = "auto"
	ChatCompletionToolChoiceOptionAutoRequired ChatCompletionToolChoiceOptionAuto = "required"
)

type ChatCompletionToolMessageParam struct {
	// The contents of the tool message.
	Content []ChatCompletionContentPartTextParam `json:"content,omitzero,required"`
	// The role of the messages author, in this case `tool`.
	//
	// This field can be elided, and will be automatically set as "tool".
	Role constant.Tool `json:"role,required"`
	// Tool call that this message is responding to.
	ToolCallID param.String `json:"tool_call_id,omitzero,required"`
	apiobject
}

func (f ChatCompletionToolMessageParam) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

func (r ChatCompletionToolMessageParam) MarshalJSON() (data []byte, err error) {
	type shadow ChatCompletionToolMessageParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// Messages sent by an end user, containing prompts or additional context
// information.
type ChatCompletionUserMessageParam struct {
	// The contents of the user message.
	Content []ChatCompletionContentPartUnionParam `json:"content,omitzero,required"`
	// The role of the messages author, in this case `user`.
	//
	// This field can be elided, and will be automatically set as "user".
	Role constant.User `json:"role,required"`
	// An optional name for the participant. Provides the model information to
	// differentiate between participants of the same role.
	Name param.String `json:"name,omitzero"`
	apiobject
}

func (f ChatCompletionUserMessageParam) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

func (r ChatCompletionUserMessageParam) MarshalJSON() (data []byte, err error) {
	type shadow ChatCompletionUserMessageParam
	return param.MarshalObject(r, (*shadow)(&r))
}

type ChatCompletionNewParams struct {
	// A list of messages comprising the conversation so far. Depending on the
	// [model](https://platform.openai.com/docs/models) you use, different message
	// types (modalities) are supported, like
	// [text](https://platform.openai.com/docs/guides/text-generation),
	// [images](https://platform.openai.com/docs/guides/vision), and
	// [audio](https://platform.openai.com/docs/guides/audio).
	Messages []ChatCompletionMessageParamUnion `json:"messages,omitzero,required"`
	// ID of the model to use. See the
	// [model endpoint compatibility](https://platform.openai.com/docs/models#model-endpoint-compatibility)
	// table for details on which models work with the Chat API.
	Model ChatModel `json:"model,omitzero,required"`
	// Parameters for audio output. Required when audio output is requested with
	// `modalities: ["audio"]`.
	// [Learn more](https://platform.openai.com/docs/guides/audio).
	Audio ChatCompletionAudioParam `json:"audio,omitzero"`
	// Number between -2.0 and 2.0. Positive values penalize new tokens based on their
	// existing frequency in the text so far, decreasing the model's likelihood to
	// repeat the same line verbatim.
	FrequencyPenalty param.Float `json:"frequency_penalty,omitzero"`
	// Deprecated in favor of `tool_choice`.
	//
	// Controls which (if any) function is called by the model.
	//
	// `none` means the model will not call a function and instead generates a message.
	//
	// `auto` means the model can pick between generating a message or calling a
	// function.
	//
	// Specifying a particular function via `{"name": "my_function"}` forces the model
	// to call that function.
	//
	// `none` is the default when no functions are present. `auto` is the default if
	// functions are present.
	FunctionCall ChatCompletionNewParamsFunctionCallUnion `json:"function_call,omitzero"`
	// Deprecated in favor of `tools`.
	//
	// A list of functions the model may generate JSON inputs for.
	Functions []ChatCompletionNewParamsFunction `json:"functions,omitzero"`
	// Modify the likelihood of specified tokens appearing in the completion.
	//
	// Accepts a JSON object that maps tokens (specified by their token ID in the
	// tokenizer) to an associated bias value from -100 to 100. Mathematically, the
	// bias is added to the logits generated by the model prior to sampling. The exact
	// effect will vary per model, but values between -1 and 1 should decrease or
	// increase likelihood of selection; values like -100 or 100 should result in a ban
	// or exclusive selection of the relevant token.
	LogitBias map[string]int64 `json:"logit_bias,omitzero"`
	// Whether to return log probabilities of the output tokens or not. If true,
	// returns the log probabilities of each output token returned in the `content` of
	// `message`.
	Logprobs param.Bool `json:"logprobs,omitzero"`
	// An upper bound for the number of tokens that can be generated for a completion,
	// including visible output tokens and
	// [reasoning tokens](https://platform.openai.com/docs/guides/reasoning).
	MaxCompletionTokens param.Int `json:"max_completion_tokens,omitzero"`
	// The maximum number of [tokens](/tokenizer) that can be generated in the chat
	// completion. This value can be used to control
	// [costs](https://openai.com/api/pricing/) for text generated via API.
	//
	// This value is now deprecated in favor of `max_completion_tokens`, and is not
	// compatible with
	// [o1 series models](https://platform.openai.com/docs/guides/reasoning).
	MaxTokens param.Int `json:"max_tokens,omitzero"`
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard.
	//
	// Keys are strings with a maximum length of 64 characters. Values are strings with
	// a maximum length of 512 characters.
	Metadata shared.MetadataParam `json:"metadata,omitzero"`
	// Output types that you would like the model to generate for this request. Most
	// models are capable of generating text, which is the default:
	//
	// `["text"]`
	//
	// The `gpt-4o-audio-preview` model can also be used to
	// [generate audio](https://platform.openai.com/docs/guides/audio). To request that
	// this model generate both text and audio responses, you can use:
	//
	// `["text", "audio"]`
	Modalities []ChatCompletionModality `json:"modalities,omitzero"`
	// How many chat completion choices to generate for each input message. Note that
	// you will be charged based on the number of generated tokens across all of the
	// choices. Keep `n` as `1` to minimize costs.
	N param.Int `json:"n,omitzero"`
	// Whether to enable
	// [parallel function calling](https://platform.openai.com/docs/guides/function-calling#configuring-parallel-function-calling)
	// during tool use.
	ParallelToolCalls param.Bool `json:"parallel_tool_calls,omitzero"`
	// Static predicted output content, such as the content of a text file that is
	// being regenerated.
	Prediction ChatCompletionPredictionContentParam `json:"prediction,omitzero"`
	// Number between -2.0 and 2.0. Positive values penalize new tokens based on
	// whether they appear in the text so far, increasing the model's likelihood to
	// talk about new topics.
	PresencePenalty param.Float `json:"presence_penalty,omitzero"`
	// **o1 and o3-mini models only**
	//
	// Constrains effort on reasoning for
	// [reasoning models](https://platform.openai.com/docs/guides/reasoning). Currently
	// supported values are `low`, `medium`, and `high`. Reducing reasoning effort can
	// result in faster responses and fewer tokens used on reasoning in a response.
	//
	// Any of "low", "medium", "high"
	ReasoningEffort ChatCompletionReasoningEffort `json:"reasoning_effort,omitzero"`
	// An object specifying the format that the model must output.
	//
	// Setting to `{ "type": "json_schema", "json_schema": {...} }` enables Structured
	// Outputs which ensures the model will match your supplied JSON schema. Learn more
	// in the
	// [Structured Outputs guide](https://platform.openai.com/docs/guides/structured-outputs).
	//
	// Setting to `{ "type": "json_object" }` enables JSON mode, which ensures the
	// message the model generates is valid JSON.
	//
	// **Important:** when using JSON mode, you **must** also instruct the model to
	// produce JSON yourself via a system or user message. Without this, the model may
	// generate an unending stream of whitespace until the generation reaches the token
	// limit, resulting in a long-running and seemingly "stuck" request. Also note that
	// the message content may be partially cut off if `finish_reason="length"`, which
	// indicates the generation exceeded `max_tokens` or the conversation exceeded the
	// max context length.
	ResponseFormat ChatCompletionNewParamsResponseFormatUnion `json:"response_format,omitzero"`
	// This feature is in Beta. If specified, our system will make a best effort to
	// sample deterministically, such that repeated requests with the same `seed` and
	// parameters should return the same result. Determinism is not guaranteed, and you
	// should refer to the `system_fingerprint` response parameter to monitor changes
	// in the backend.
	Seed param.Int `json:"seed,omitzero"`
	// Specifies the latency tier to use for processing the request. This parameter is
	// relevant for customers subscribed to the scale tier service:
	//
	//   - If set to 'auto', and the Project is Scale tier enabled, the system will
	//     utilize scale tier credits until they are exhausted.
	//   - If set to 'auto', and the Project is not Scale tier enabled, the request will
	//     be processed using the default service tier with a lower uptime SLA and no
	//     latency guarantee.
	//   - If set to 'default', the request will be processed using the default service
	//     tier with a lower uptime SLA and no latency guarantee.
	//   - When not set, the default behavior is 'auto'.
	//
	// Any of "auto", "default"
	ServiceTier ChatCompletionNewParamsServiceTier `json:"service_tier,omitzero"`
	// Up to 4 sequences where the API will stop generating further tokens.
	Stop ChatCompletionNewParamsStopUnion `json:"stop,omitzero"`
	// Whether or not to store the output of this chat completion request for use in
	// our [model distillation](https://platform.openai.com/docs/guides/distillation)
	// or [evals](https://platform.openai.com/docs/guides/evals) products.
	Store param.Bool `json:"store,omitzero"`
	// Options for streaming response. Only set this when you set `stream: true`.
	StreamOptions ChatCompletionStreamOptionsParam `json:"stream_options,omitzero"`
	// What sampling temperature to use, between 0 and 2. Higher values like 0.8 will
	// make the output more random, while lower values like 0.2 will make it more
	// focused and deterministic. We generally recommend altering this or `top_p` but
	// not both.
	Temperature param.Float `json:"temperature,omitzero"`
	// Controls which (if any) tool is called by the model. `none` means the model will
	// not call any tool and instead generates a message. `auto` means the model can
	// pick between generating a message or calling one or more tools. `required` means
	// the model must call one or more tools. Specifying a particular tool via
	// `{"type": "function", "function": {"name": "my_function"}}` forces the model to
	// call that tool.
	//
	// `none` is the default when no tools are present. `auto` is the default if tools
	// are present.
	ToolChoice ChatCompletionToolChoiceOptionUnionParam `json:"tool_choice,omitzero"`
	// A list of tools the model may call. Currently, only functions are supported as a
	// tool. Use this to provide a list of functions the model may generate JSON inputs
	// for. A max of 128 functions are supported.
	Tools []ChatCompletionToolParam `json:"tools,omitzero"`
	// An integer between 0 and 20 specifying the number of most likely tokens to
	// return at each token position, each with an associated log probability.
	// `logprobs` must be set to `true` if this parameter is used.
	TopLogprobs param.Int `json:"top_logprobs,omitzero"`
	// An alternative to sampling with temperature, called nucleus sampling, where the
	// model considers the results of the tokens with top_p probability mass. So 0.1
	// means only the tokens comprising the top 10% probability mass are considered.
	//
	// We generally recommend altering this or `temperature` but not both.
	TopP param.Float `json:"top_p,omitzero"`
	// A unique identifier representing your end-user, which can help OpenAI to monitor
	// and detect abuse.
	// [Learn more](https://platform.openai.com/docs/guides/safety-best-practices#end-user-ids).
	User param.String `json:"user,omitzero"`
	apiobject
}

func (f ChatCompletionNewParams) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

func (r ChatCompletionNewParams) MarshalJSON() (data []byte, err error) {
	type shadow ChatCompletionNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}

// Only one field can be non-zero
type ChatCompletionNewParamsFunctionCallUnion struct {
	// Check if union is this variant with !param.IsOmitted(union.OfAuto)
	OfAuto               string
	OfFunctionCallOption *ChatCompletionFunctionCallOptionParam
	apiunion
}

func (u ChatCompletionNewParamsFunctionCallUnion) IsMissing() bool {
	return param.IsOmitted(u) || u.IsNull()
}

func (u ChatCompletionNewParamsFunctionCallUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion[ChatCompletionNewParamsFunctionCallUnion](u.OfAuto, u.OfFunctionCallOption)
}

func (u ChatCompletionNewParamsFunctionCallUnion) GetName() *string {
	if vt := u.OfFunctionCallOption; vt != nil && !vt.Name.IsOmitted() {
		return &vt.Name.V
	}
	return nil
}

// `none` means the model will not call a function and instead generates a message.
// `auto` means the model can pick between generating a message or calling a
// function.
type ChatCompletionNewParamsFunctionCallAuto = string

const (
	ChatCompletionNewParamsFunctionCallAutoNone ChatCompletionNewParamsFunctionCallAuto = "none"
	ChatCompletionNewParamsFunctionCallAutoAuto ChatCompletionNewParamsFunctionCallAuto = "auto"
)

// Deprecated: deprecated
type ChatCompletionNewParamsFunction struct {
	// The name of the function to be called. Must be a-z, A-Z, 0-9, or contain
	// underscores and dashes, with a maximum length of 64.
	Name param.String `json:"name,omitzero,required"`
	// A description of what the function does, used by the model to choose when and
	// how to call the function.
	Description param.String `json:"description,omitzero"`
	// The parameters the functions accepts, described as a JSON Schema object. See the
	// [guide](https://platform.openai.com/docs/guides/function-calling) for examples,
	// and the
	// [JSON Schema reference](https://json-schema.org/understanding-json-schema/) for
	// documentation about the format.
	//
	// Omitting `parameters` defines a function with an empty parameter list.
	Parameters shared.FunctionParameters `json:"parameters,omitzero"`
	apiobject
}

func (f ChatCompletionNewParamsFunction) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

func (r ChatCompletionNewParamsFunction) MarshalJSON() (data []byte, err error) {
	type shadow ChatCompletionNewParamsFunction
	return param.MarshalObject(r, (*shadow)(&r))
}

// Only one field can be non-zero
type ChatCompletionNewParamsResponseFormatUnion struct {
	OfResponseFormatText       *shared.ResponseFormatTextParam
	OfResponseFormatJSONObject *shared.ResponseFormatJSONObjectParam
	OfResponseFormatJSONSchema *shared.ResponseFormatJSONSchemaParam
	apiunion
}

func (u ChatCompletionNewParamsResponseFormatUnion) IsMissing() bool {
	return param.IsOmitted(u) || u.IsNull()
}

func (u ChatCompletionNewParamsResponseFormatUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion[ChatCompletionNewParamsResponseFormatUnion](u.OfResponseFormatText, u.OfResponseFormatJSONObject, u.OfResponseFormatJSONSchema)
}

func (u ChatCompletionNewParamsResponseFormatUnion) GetJSONSchema() *shared.ResponseFormatJSONSchemaJSONSchemaParam {
	if vt := u.OfResponseFormatJSONSchema; vt != nil {
		return &vt.JSONSchema
	}
	return nil
}

func (u ChatCompletionNewParamsResponseFormatUnion) GetType() *string {
	if vt := u.OfResponseFormatText; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfResponseFormatJSONObject; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfResponseFormatJSONSchema; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Specifies the latency tier to use for processing the request. This parameter is
// relevant for customers subscribed to the scale tier service:
//
//   - If set to 'auto', and the Project is Scale tier enabled, the system will
//     utilize scale tier credits until they are exhausted.
//   - If set to 'auto', and the Project is not Scale tier enabled, the request will
//     be processed using the default service tier with a lower uptime SLA and no
//     latency guarantee.
//   - If set to 'default', the request will be processed using the default service
//     tier with a lower uptime SLA and no latency guarantee.
//   - When not set, the default behavior is 'auto'.
type ChatCompletionNewParamsServiceTier string

const (
	ChatCompletionNewParamsServiceTierAuto    ChatCompletionNewParamsServiceTier = "auto"
	ChatCompletionNewParamsServiceTierDefault ChatCompletionNewParamsServiceTier = "default"
)

// Only one field can be non-zero
type ChatCompletionNewParamsStopUnion struct {
	OfString                      param.String
	OfChatCompletionNewsStopArray []string
	apiunion
}

func (u ChatCompletionNewParamsStopUnion) IsMissing() bool { return param.IsOmitted(u) || u.IsNull() }

func (u ChatCompletionNewParamsStopUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion[ChatCompletionNewParamsStopUnion](u.OfString, u.OfChatCompletionNewsStopArray)
}

type ChatCompletionUpdateParams struct {
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard.
	//
	// Keys are strings with a maximum length of 64 characters. Values are strings with
	// a maximum length of 512 characters.
	Metadata shared.MetadataParam `json:"metadata,omitzero,required"`
	apiobject
}

func (f ChatCompletionUpdateParams) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

func (r ChatCompletionUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow ChatCompletionUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}

type ChatCompletionListParams struct {
	// Identifier for the last chat completion from the previous pagination request.
	After param.String `query:"after,omitzero"`
	// Number of chat completions to retrieve.
	Limit param.Int `query:"limit,omitzero"`
	// A list of metadata keys to filter the chat completions by. Example:
	//
	// `metadata[key1]=value1&metadata[key2]=value2`
	Metadata shared.MetadataParam `query:"metadata,omitzero"`
	// The model used to generate the chat completions.
	Model param.String `query:"model,omitzero"`
	// Sort order for chat completions by timestamp. Use `asc` for ascending order or
	// `desc` for descending order. Defaults to `asc`.
	//
	// Any of "asc", "desc"
	Order ChatCompletionListParamsOrder `query:"order,omitzero"`
	apiobject
}

func (f ChatCompletionListParams) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

// URLQuery serializes [ChatCompletionListParams]'s query parameters as
// `url.Values`.
func (r ChatCompletionListParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Sort order for chat completions by timestamp. Use `asc` for ascending order or
// `desc` for descending order. Defaults to `asc`.
type ChatCompletionListParamsOrder string

const (
	ChatCompletionListParamsOrderAsc  ChatCompletionListParamsOrder = "asc"
	ChatCompletionListParamsOrderDesc ChatCompletionListParamsOrder = "desc"
)

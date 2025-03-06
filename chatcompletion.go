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
	"github.com/openai/openai-go/internal/param"
	"github.com/openai/openai-go/internal/requestconfig"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/pagination"
	"github.com/openai/openai-go/packages/ssestream"
	"github.com/openai/openai-go/shared"
	"github.com/tidwall/sjson"
)

func UserMessage(content string) ChatCompletionMessageParamUnion {
	return UserMessageParts(TextPart(content))
}

func UserMessageParts(parts ...ChatCompletionContentPartUnionParam) ChatCompletionUserMessageParam {
	return ChatCompletionUserMessageParam{
		Role:    F(ChatCompletionUserMessageParamRoleUser),
		Content: F(parts),
	}
}

func TextPart(content string) ChatCompletionContentPartTextParam {
	return ChatCompletionContentPartTextParam{
		Type: F(ChatCompletionContentPartTextTypeText),
		Text: F(content),
	}
}

func RefusalPart(refusal string) ChatCompletionContentPartRefusalParam {
	return ChatCompletionContentPartRefusalParam{
		Type:    F(ChatCompletionContentPartRefusalTypeRefusal),
		Refusal: F(refusal),
	}
}

func ImagePart(url string) ChatCompletionContentPartImageParam {
	return ChatCompletionContentPartImageParam{
		Type: F(ChatCompletionContentPartImageTypeImageURL),
		ImageURL: F(ChatCompletionContentPartImageImageURLParam{
			URL: F(url),
		}),
	}
}

func AssistantMessage(content string) ChatCompletionAssistantMessageParam {
	return ChatCompletionAssistantMessageParam{
		Role: F(ChatCompletionAssistantMessageParamRoleAssistant),
		Content: F([]ChatCompletionAssistantMessageParamContentUnion{
			TextPart(content),
		}),
	}
}

func ToolMessage(toolCallID, content string) ChatCompletionToolMessageParam {
	return ChatCompletionToolMessageParam{
		Role:       F(ChatCompletionToolMessageParamRoleTool),
		ToolCallID: F(toolCallID),
		Content: F([]ChatCompletionContentPartTextParam{
			TextPart(content),
		}),
	}
}

func SystemMessage(content string) ChatCompletionMessageParamUnion {
	return ChatCompletionSystemMessageParam{
		Role: F(ChatCompletionSystemMessageParamRoleSystem),
		Content: F([]ChatCompletionContentPartTextParam{
			TextPart(content),
		}),
	}
}

func FunctionMessage(name, content string) ChatCompletionMessageParamUnion {
	return ChatCompletionFunctionMessageParam{
		Role:    F(ChatCompletionFunctionMessageParamRoleFunction),
		Name:    F(name),
		Content: F(content),
	}
}

// ChatCompletionService contains methods and other services that help with
// interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewChatCompletionService] method instead.
type ChatCompletionService struct {
	Options  []option.RequestOption
	Messages *ChatCompletionMessageService
}

// NewChatCompletionService generates a new service that applies the given options
// to each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewChatCompletionService(opts ...option.RequestOption) (r *ChatCompletionService) {
	r = &ChatCompletionService{}
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
	ID string `json:"id,required"`
	// A list of chat completion choices. Can be more than one if `n` is greater
	// than 1.
	Choices []ChatCompletionChoice `json:"choices,required"`
	// The Unix timestamp (in seconds) of when the chat completion was created.
	Created int64 `json:"created,required"`
	// The model used for the chat completion.
	Model string `json:"model,required"`
	// The object type, which is always `chat.completion`.
	Object ChatCompletionObject `json:"object,required"`
	// The service tier used for processing the request.
	ServiceTier ChatCompletionServiceTier `json:"service_tier,nullable"`
	// This fingerprint represents the backend configuration that the model runs with.
	//
	// Can be used in conjunction with the `seed` request parameter to understand when
	// backend changes have been made that might impact determinism.
	SystemFingerprint string `json:"system_fingerprint"`
	// Usage statistics for the completion request.
	Usage CompletionUsage    `json:"usage"`
	JSON  chatCompletionJSON `json:"-"`
}

// chatCompletionJSON contains the JSON metadata for the struct [ChatCompletion]
type chatCompletionJSON struct {
	ID                apijson.Field
	Choices           apijson.Field
	Created           apijson.Field
	Model             apijson.Field
	Object            apijson.Field
	ServiceTier       apijson.Field
	SystemFingerprint apijson.Field
	Usage             apijson.Field
	raw               string
	ExtraFields       map[string]apijson.Field
}

func (r *ChatCompletion) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r chatCompletionJSON) RawJSON() string {
	return r.raw
}

type ChatCompletionChoice struct {
	// The reason the model stopped generating tokens. This will be `stop` if the model
	// hit a natural stop point or a provided stop sequence, `length` if the maximum
	// number of tokens specified in the request was reached, `content_filter` if
	// content was omitted due to a flag from our content filters, `tool_calls` if the
	// model called a tool, or `function_call` (deprecated) if the model called a
	// function.
	FinishReason ChatCompletionChoicesFinishReason `json:"finish_reason,required"`
	// The index of the choice in the list of choices.
	Index int64 `json:"index,required"`
	// Log probability information for the choice.
	Logprobs ChatCompletionChoicesLogprobs `json:"logprobs,required,nullable"`
	// A chat completion message generated by the model.
	Message ChatCompletionMessage    `json:"message,required"`
	JSON    chatCompletionChoiceJSON `json:"-"`
}

// chatCompletionChoiceJSON contains the JSON metadata for the struct
// [ChatCompletionChoice]
type chatCompletionChoiceJSON struct {
	FinishReason apijson.Field
	Index        apijson.Field
	Logprobs     apijson.Field
	Message      apijson.Field
	raw          string
	ExtraFields  map[string]apijson.Field
}

func (r *ChatCompletionChoice) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r chatCompletionChoiceJSON) RawJSON() string {
	return r.raw
}

// The reason the model stopped generating tokens. This will be `stop` if the model
// hit a natural stop point or a provided stop sequence, `length` if the maximum
// number of tokens specified in the request was reached, `content_filter` if
// content was omitted due to a flag from our content filters, `tool_calls` if the
// model called a tool, or `function_call` (deprecated) if the model called a
// function.
type ChatCompletionChoicesFinishReason string

const (
	ChatCompletionChoicesFinishReasonStop          ChatCompletionChoicesFinishReason = "stop"
	ChatCompletionChoicesFinishReasonLength        ChatCompletionChoicesFinishReason = "length"
	ChatCompletionChoicesFinishReasonToolCalls     ChatCompletionChoicesFinishReason = "tool_calls"
	ChatCompletionChoicesFinishReasonContentFilter ChatCompletionChoicesFinishReason = "content_filter"
	ChatCompletionChoicesFinishReasonFunctionCall  ChatCompletionChoicesFinishReason = "function_call"
)

func (r ChatCompletionChoicesFinishReason) IsKnown() bool {
	switch r {
	case ChatCompletionChoicesFinishReasonStop, ChatCompletionChoicesFinishReasonLength, ChatCompletionChoicesFinishReasonToolCalls, ChatCompletionChoicesFinishReasonContentFilter, ChatCompletionChoicesFinishReasonFunctionCall:
		return true
	}
	return false
}

// Log probability information for the choice.
type ChatCompletionChoicesLogprobs struct {
	// A list of message content tokens with log probability information.
	Content []ChatCompletionTokenLogprob `json:"content,required,nullable"`
	// A list of message refusal tokens with log probability information.
	Refusal []ChatCompletionTokenLogprob      `json:"refusal,required,nullable"`
	JSON    chatCompletionChoicesLogprobsJSON `json:"-"`
}

// chatCompletionChoicesLogprobsJSON contains the JSON metadata for the struct
// [ChatCompletionChoicesLogprobs]
type chatCompletionChoicesLogprobsJSON struct {
	Content     apijson.Field
	Refusal     apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ChatCompletionChoicesLogprobs) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r chatCompletionChoicesLogprobsJSON) RawJSON() string {
	return r.raw
}

// The object type, which is always `chat.completion`.
type ChatCompletionObject string

const (
	ChatCompletionObjectChatCompletion ChatCompletionObject = "chat.completion"
)

func (r ChatCompletionObject) IsKnown() bool {
	switch r {
	case ChatCompletionObjectChatCompletion:
		return true
	}
	return false
}

// The service tier used for processing the request.
type ChatCompletionServiceTier string

const (
	ChatCompletionServiceTierScale   ChatCompletionServiceTier = "scale"
	ChatCompletionServiceTierDefault ChatCompletionServiceTier = "default"
)

func (r ChatCompletionServiceTier) IsKnown() bool {
	switch r {
	case ChatCompletionServiceTierScale, ChatCompletionServiceTierDefault:
		return true
	}
	return false
}

// Messages sent by the model in response to user messages.
type ChatCompletionAssistantMessageParam struct {
	// The role of the messages author, in this case `assistant`.
	Role param.Field[ChatCompletionAssistantMessageParamRole] `json:"role,required"`
	// Data about a previous audio response from the model.
	// [Learn more](https://platform.openai.com/docs/guides/audio).
	Audio param.Field[ChatCompletionAssistantMessageParamAudio] `json:"audio"`
	// The contents of the assistant message. Required unless `tool_calls` or
	// `function_call` is specified.
	Content param.Field[[]ChatCompletionAssistantMessageParamContentUnion] `json:"content"`
	// Deprecated and replaced by `tool_calls`. The name and arguments of a function
	// that should be called, as generated by the model.
	//
	// Deprecated: deprecated
	FunctionCall param.Field[ChatCompletionAssistantMessageParamFunctionCall] `json:"function_call"`
	// An optional name for the participant. Provides the model information to
	// differentiate between participants of the same role.
	Name param.Field[string] `json:"name"`
	// The refusal message by the assistant.
	Refusal param.Field[string] `json:"refusal"`
	// The tool calls generated by the model, such as function calls.
	ToolCalls param.Field[[]ChatCompletionMessageToolCallParam] `json:"tool_calls"`
}

func (r ChatCompletionAssistantMessageParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ChatCompletionAssistantMessageParam) implementsChatCompletionMessageParamUnion() {}

// The role of the messages author, in this case `assistant`.
type ChatCompletionAssistantMessageParamRole string

const (
	ChatCompletionAssistantMessageParamRoleAssistant ChatCompletionAssistantMessageParamRole = "assistant"
)

func (r ChatCompletionAssistantMessageParamRole) IsKnown() bool {
	switch r {
	case ChatCompletionAssistantMessageParamRoleAssistant:
		return true
	}
	return false
}

// Data about a previous audio response from the model.
// [Learn more](https://platform.openai.com/docs/guides/audio).
type ChatCompletionAssistantMessageParamAudio struct {
	// Unique identifier for a previous audio response from the model.
	ID param.Field[string] `json:"id,required"`
}

func (r ChatCompletionAssistantMessageParamAudio) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Learn about
// [text inputs](https://platform.openai.com/docs/guides/text-generation).
type ChatCompletionAssistantMessageParamContent struct {
	// The type of the content part.
	Type param.Field[ChatCompletionAssistantMessageParamContentType] `json:"type,required"`
	// The refusal message generated by the model.
	Refusal param.Field[string] `json:"refusal"`
	// The text content.
	Text param.Field[string] `json:"text"`
}

func (r ChatCompletionAssistantMessageParamContent) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ChatCompletionAssistantMessageParamContent) implementsChatCompletionAssistantMessageParamContentUnion() {
}

// Learn about
// [text inputs](https://platform.openai.com/docs/guides/text-generation).
//
// Satisfied by [ChatCompletionContentPartTextParam],
// [ChatCompletionContentPartRefusalParam],
// [ChatCompletionAssistantMessageParamContent].
type ChatCompletionAssistantMessageParamContentUnion interface {
	implementsChatCompletionAssistantMessageParamContentUnion()
}

// The type of the content part.
type ChatCompletionAssistantMessageParamContentType string

const (
	ChatCompletionAssistantMessageParamContentTypeText    ChatCompletionAssistantMessageParamContentType = "text"
	ChatCompletionAssistantMessageParamContentTypeRefusal ChatCompletionAssistantMessageParamContentType = "refusal"
)

func (r ChatCompletionAssistantMessageParamContentType) IsKnown() bool {
	switch r {
	case ChatCompletionAssistantMessageParamContentTypeText, ChatCompletionAssistantMessageParamContentTypeRefusal:
		return true
	}
	return false
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
	Arguments param.Field[string] `json:"arguments,required"`
	// The name of the function to call.
	Name param.Field[string] `json:"name,required"`
}

func (r ChatCompletionAssistantMessageParamFunctionCall) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// If the audio output modality is requested, this object contains data about the
// audio response from the model.
// [Learn more](https://platform.openai.com/docs/guides/audio).
type ChatCompletionAudio struct {
	// Unique identifier for this audio response.
	ID string `json:"id,required"`
	// Base64 encoded audio bytes generated by the model, in the format specified in
	// the request.
	Data string `json:"data,required"`
	// The Unix timestamp (in seconds) for when this audio response will no longer be
	// accessible on the server for use in multi-turn conversations.
	ExpiresAt int64 `json:"expires_at,required"`
	// Transcript of the audio generated by the model.
	Transcript string                  `json:"transcript,required"`
	JSON       chatCompletionAudioJSON `json:"-"`
}

// chatCompletionAudioJSON contains the JSON metadata for the struct
// [ChatCompletionAudio]
type chatCompletionAudioJSON struct {
	ID          apijson.Field
	Data        apijson.Field
	ExpiresAt   apijson.Field
	Transcript  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ChatCompletionAudio) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r chatCompletionAudioJSON) RawJSON() string {
	return r.raw
}

// Parameters for audio output. Required when audio output is requested with
// `modalities: ["audio"]`.
// [Learn more](https://platform.openai.com/docs/guides/audio).
type ChatCompletionAudioParam struct {
	// Specifies the output audio format. Must be one of `wav`, `mp3`, `flac`, `opus`,
	// or `pcm16`.
	Format param.Field[ChatCompletionAudioParamFormat] `json:"format,required"`
	// The voice the model uses to respond. Supported voices are `ash`, `ballad`,
	// `coral`, `sage`, and `verse` (also supported but not recommended are `alloy`,
	// `echo`, and `shimmer`; these voices are less expressive).
	Voice param.Field[ChatCompletionAudioParamVoice] `json:"voice,required"`
}

func (r ChatCompletionAudioParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Specifies the output audio format. Must be one of `wav`, `mp3`, `flac`, `opus`,
// or `pcm16`.
type ChatCompletionAudioParamFormat string

const (
	ChatCompletionAudioParamFormatWAV   ChatCompletionAudioParamFormat = "wav"
	ChatCompletionAudioParamFormatMP3   ChatCompletionAudioParamFormat = "mp3"
	ChatCompletionAudioParamFormatFLAC  ChatCompletionAudioParamFormat = "flac"
	ChatCompletionAudioParamFormatOpus  ChatCompletionAudioParamFormat = "opus"
	ChatCompletionAudioParamFormatPcm16 ChatCompletionAudioParamFormat = "pcm16"
)

func (r ChatCompletionAudioParamFormat) IsKnown() bool {
	switch r {
	case ChatCompletionAudioParamFormatWAV, ChatCompletionAudioParamFormatMP3, ChatCompletionAudioParamFormatFLAC, ChatCompletionAudioParamFormatOpus, ChatCompletionAudioParamFormatPcm16:
		return true
	}
	return false
}

// The voice the model uses to respond. Supported voices are `ash`, `ballad`,
// `coral`, `sage`, and `verse` (also supported but not recommended are `alloy`,
// `echo`, and `shimmer`; these voices are less expressive).
type ChatCompletionAudioParamVoice string

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

func (r ChatCompletionAudioParamVoice) IsKnown() bool {
	switch r {
	case ChatCompletionAudioParamVoiceAlloy, ChatCompletionAudioParamVoiceAsh, ChatCompletionAudioParamVoiceBallad, ChatCompletionAudioParamVoiceCoral, ChatCompletionAudioParamVoiceEcho, ChatCompletionAudioParamVoiceSage, ChatCompletionAudioParamVoiceShimmer, ChatCompletionAudioParamVoiceVerse:
		return true
	}
	return false
}

// Represents a streamed chunk of a chat completion response returned by model,
// based on the provided input.
type ChatCompletionChunk struct {
	// A unique identifier for the chat completion. Each chunk has the same ID.
	ID string `json:"id,required"`
	// A list of chat completion choices. Can contain more than one elements if `n` is
	// greater than 1. Can also be empty for the last chunk if you set
	// `stream_options: {"include_usage": true}`.
	Choices []ChatCompletionChunkChoice `json:"choices,required"`
	// The Unix timestamp (in seconds) of when the chat completion was created. Each
	// chunk has the same timestamp.
	Created int64 `json:"created,required"`
	// The model to generate the completion.
	Model string `json:"model,required"`
	// The object type, which is always `chat.completion.chunk`.
	Object ChatCompletionChunkObject `json:"object,required"`
	// The service tier used for processing the request.
	ServiceTier ChatCompletionChunkServiceTier `json:"service_tier,nullable"`
	// This fingerprint represents the backend configuration that the model runs with.
	// Can be used in conjunction with the `seed` request parameter to understand when
	// backend changes have been made that might impact determinism.
	SystemFingerprint string `json:"system_fingerprint"`
	// An optional field that will only be present when you set
	// `stream_options: {"include_usage": true}` in your request. When present, it
	// contains a null value except for the last chunk which contains the token usage
	// statistics for the entire request.
	Usage CompletionUsage         `json:"usage,nullable"`
	JSON  chatCompletionChunkJSON `json:"-"`
}

// chatCompletionChunkJSON contains the JSON metadata for the struct
// [ChatCompletionChunk]
type chatCompletionChunkJSON struct {
	ID                apijson.Field
	Choices           apijson.Field
	Created           apijson.Field
	Model             apijson.Field
	Object            apijson.Field
	ServiceTier       apijson.Field
	SystemFingerprint apijson.Field
	Usage             apijson.Field
	raw               string
	ExtraFields       map[string]apijson.Field
}

func (r *ChatCompletionChunk) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r chatCompletionChunkJSON) RawJSON() string {
	return r.raw
}

type ChatCompletionChunkChoice struct {
	// A chat completion delta generated by streamed model responses.
	Delta ChatCompletionChunkChoicesDelta `json:"delta,required"`
	// The reason the model stopped generating tokens. This will be `stop` if the model
	// hit a natural stop point or a provided stop sequence, `length` if the maximum
	// number of tokens specified in the request was reached, `content_filter` if
	// content was omitted due to a flag from our content filters, `tool_calls` if the
	// model called a tool, or `function_call` (deprecated) if the model called a
	// function.
	FinishReason ChatCompletionChunkChoicesFinishReason `json:"finish_reason,required,nullable"`
	// The index of the choice in the list of choices.
	Index int64 `json:"index,required"`
	// Log probability information for the choice.
	Logprobs ChatCompletionChunkChoicesLogprobs `json:"logprobs,nullable"`
	JSON     chatCompletionChunkChoiceJSON      `json:"-"`
}

// chatCompletionChunkChoiceJSON contains the JSON metadata for the struct
// [ChatCompletionChunkChoice]
type chatCompletionChunkChoiceJSON struct {
	Delta        apijson.Field
	FinishReason apijson.Field
	Index        apijson.Field
	Logprobs     apijson.Field
	raw          string
	ExtraFields  map[string]apijson.Field
}

func (r *ChatCompletionChunkChoice) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r chatCompletionChunkChoiceJSON) RawJSON() string {
	return r.raw
}

// A chat completion delta generated by streamed model responses.
type ChatCompletionChunkChoicesDelta struct {
	// The contents of the chunk message.
	Content string `json:"content,nullable"`
	// Deprecated and replaced by `tool_calls`. The name and arguments of a function
	// that should be called, as generated by the model.
	//
	// Deprecated: deprecated
	FunctionCall ChatCompletionChunkChoicesDeltaFunctionCall `json:"function_call"`
	// The refusal message generated by the model.
	Refusal string `json:"refusal,nullable"`
	// The role of the author of this message.
	Role      ChatCompletionChunkChoicesDeltaRole       `json:"role"`
	ToolCalls []ChatCompletionChunkChoicesDeltaToolCall `json:"tool_calls"`
	JSON      chatCompletionChunkChoicesDeltaJSON       `json:"-"`
}

// chatCompletionChunkChoicesDeltaJSON contains the JSON metadata for the struct
// [ChatCompletionChunkChoicesDelta]
type chatCompletionChunkChoicesDeltaJSON struct {
	Content      apijson.Field
	FunctionCall apijson.Field
	Refusal      apijson.Field
	Role         apijson.Field
	ToolCalls    apijson.Field
	raw          string
	ExtraFields  map[string]apijson.Field
}

func (r *ChatCompletionChunkChoicesDelta) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r chatCompletionChunkChoicesDeltaJSON) RawJSON() string {
	return r.raw
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
	Arguments string `json:"arguments"`
	// The name of the function to call.
	Name string                                          `json:"name"`
	JSON chatCompletionChunkChoicesDeltaFunctionCallJSON `json:"-"`
}

// chatCompletionChunkChoicesDeltaFunctionCallJSON contains the JSON metadata for
// the struct [ChatCompletionChunkChoicesDeltaFunctionCall]
type chatCompletionChunkChoicesDeltaFunctionCallJSON struct {
	Arguments   apijson.Field
	Name        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ChatCompletionChunkChoicesDeltaFunctionCall) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r chatCompletionChunkChoicesDeltaFunctionCallJSON) RawJSON() string {
	return r.raw
}

// The role of the author of this message.
type ChatCompletionChunkChoicesDeltaRole string

const (
	ChatCompletionChunkChoicesDeltaRoleDeveloper ChatCompletionChunkChoicesDeltaRole = "developer"
	ChatCompletionChunkChoicesDeltaRoleSystem    ChatCompletionChunkChoicesDeltaRole = "system"
	ChatCompletionChunkChoicesDeltaRoleUser      ChatCompletionChunkChoicesDeltaRole = "user"
	ChatCompletionChunkChoicesDeltaRoleAssistant ChatCompletionChunkChoicesDeltaRole = "assistant"
	ChatCompletionChunkChoicesDeltaRoleTool      ChatCompletionChunkChoicesDeltaRole = "tool"
)

func (r ChatCompletionChunkChoicesDeltaRole) IsKnown() bool {
	switch r {
	case ChatCompletionChunkChoicesDeltaRoleDeveloper, ChatCompletionChunkChoicesDeltaRoleSystem, ChatCompletionChunkChoicesDeltaRoleUser, ChatCompletionChunkChoicesDeltaRoleAssistant, ChatCompletionChunkChoicesDeltaRoleTool:
		return true
	}
	return false
}

type ChatCompletionChunkChoicesDeltaToolCall struct {
	Index int64 `json:"index,required"`
	// The ID of the tool call.
	ID       string                                           `json:"id"`
	Function ChatCompletionChunkChoicesDeltaToolCallsFunction `json:"function"`
	// The type of the tool. Currently, only `function` is supported.
	Type ChatCompletionChunkChoicesDeltaToolCallsType `json:"type"`
	JSON chatCompletionChunkChoicesDeltaToolCallJSON  `json:"-"`
}

// chatCompletionChunkChoicesDeltaToolCallJSON contains the JSON metadata for the
// struct [ChatCompletionChunkChoicesDeltaToolCall]
type chatCompletionChunkChoicesDeltaToolCallJSON struct {
	Index       apijson.Field
	ID          apijson.Field
	Function    apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ChatCompletionChunkChoicesDeltaToolCall) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r chatCompletionChunkChoicesDeltaToolCallJSON) RawJSON() string {
	return r.raw
}

type ChatCompletionChunkChoicesDeltaToolCallsFunction struct {
	// The arguments to call the function with, as generated by the model in JSON
	// format. Note that the model does not always generate valid JSON, and may
	// hallucinate parameters not defined by your function schema. Validate the
	// arguments in your code before calling your function.
	Arguments string `json:"arguments"`
	// The name of the function to call.
	Name string                                               `json:"name"`
	JSON chatCompletionChunkChoicesDeltaToolCallsFunctionJSON `json:"-"`
}

// chatCompletionChunkChoicesDeltaToolCallsFunctionJSON contains the JSON metadata
// for the struct [ChatCompletionChunkChoicesDeltaToolCallsFunction]
type chatCompletionChunkChoicesDeltaToolCallsFunctionJSON struct {
	Arguments   apijson.Field
	Name        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ChatCompletionChunkChoicesDeltaToolCallsFunction) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r chatCompletionChunkChoicesDeltaToolCallsFunctionJSON) RawJSON() string {
	return r.raw
}

// The type of the tool. Currently, only `function` is supported.
type ChatCompletionChunkChoicesDeltaToolCallsType string

const (
	ChatCompletionChunkChoicesDeltaToolCallsTypeFunction ChatCompletionChunkChoicesDeltaToolCallsType = "function"
)

func (r ChatCompletionChunkChoicesDeltaToolCallsType) IsKnown() bool {
	switch r {
	case ChatCompletionChunkChoicesDeltaToolCallsTypeFunction:
		return true
	}
	return false
}

// The reason the model stopped generating tokens. This will be `stop` if the model
// hit a natural stop point or a provided stop sequence, `length` if the maximum
// number of tokens specified in the request was reached, `content_filter` if
// content was omitted due to a flag from our content filters, `tool_calls` if the
// model called a tool, or `function_call` (deprecated) if the model called a
// function.
type ChatCompletionChunkChoicesFinishReason string

const (
	ChatCompletionChunkChoicesFinishReasonStop          ChatCompletionChunkChoicesFinishReason = "stop"
	ChatCompletionChunkChoicesFinishReasonLength        ChatCompletionChunkChoicesFinishReason = "length"
	ChatCompletionChunkChoicesFinishReasonToolCalls     ChatCompletionChunkChoicesFinishReason = "tool_calls"
	ChatCompletionChunkChoicesFinishReasonContentFilter ChatCompletionChunkChoicesFinishReason = "content_filter"
	ChatCompletionChunkChoicesFinishReasonFunctionCall  ChatCompletionChunkChoicesFinishReason = "function_call"
)

func (r ChatCompletionChunkChoicesFinishReason) IsKnown() bool {
	switch r {
	case ChatCompletionChunkChoicesFinishReasonStop, ChatCompletionChunkChoicesFinishReasonLength, ChatCompletionChunkChoicesFinishReasonToolCalls, ChatCompletionChunkChoicesFinishReasonContentFilter, ChatCompletionChunkChoicesFinishReasonFunctionCall:
		return true
	}
	return false
}

// Log probability information for the choice.
type ChatCompletionChunkChoicesLogprobs struct {
	// A list of message content tokens with log probability information.
	Content []ChatCompletionTokenLogprob `json:"content,required,nullable"`
	// A list of message refusal tokens with log probability information.
	Refusal []ChatCompletionTokenLogprob           `json:"refusal,required,nullable"`
	JSON    chatCompletionChunkChoicesLogprobsJSON `json:"-"`
}

// chatCompletionChunkChoicesLogprobsJSON contains the JSON metadata for the struct
// [ChatCompletionChunkChoicesLogprobs]
type chatCompletionChunkChoicesLogprobsJSON struct {
	Content     apijson.Field
	Refusal     apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ChatCompletionChunkChoicesLogprobs) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r chatCompletionChunkChoicesLogprobsJSON) RawJSON() string {
	return r.raw
}

// The object type, which is always `chat.completion.chunk`.
type ChatCompletionChunkObject string

const (
	ChatCompletionChunkObjectChatCompletionChunk ChatCompletionChunkObject = "chat.completion.chunk"
)

func (r ChatCompletionChunkObject) IsKnown() bool {
	switch r {
	case ChatCompletionChunkObjectChatCompletionChunk:
		return true
	}
	return false
}

// The service tier used for processing the request.
type ChatCompletionChunkServiceTier string

const (
	ChatCompletionChunkServiceTierScale   ChatCompletionChunkServiceTier = "scale"
	ChatCompletionChunkServiceTierDefault ChatCompletionChunkServiceTier = "default"
)

func (r ChatCompletionChunkServiceTier) IsKnown() bool {
	switch r {
	case ChatCompletionChunkServiceTierScale, ChatCompletionChunkServiceTierDefault:
		return true
	}
	return false
}

// Learn about
// [text inputs](https://platform.openai.com/docs/guides/text-generation).
type ChatCompletionContentPartParam struct {
	// The type of the content part.
	Type       param.Field[ChatCompletionContentPartType] `json:"type,required"`
	ImageURL   param.Field[interface{}]                   `json:"image_url"`
	InputAudio param.Field[interface{}]                   `json:"input_audio"`
	// The text content.
	Text param.Field[string] `json:"text"`
}

func (r ChatCompletionContentPartParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ChatCompletionContentPartParam) implementsChatCompletionContentPartUnionParam() {}

// Learn about
// [text inputs](https://platform.openai.com/docs/guides/text-generation).
//
// Satisfied by [ChatCompletionContentPartTextParam],
// [ChatCompletionContentPartImageParam],
// [ChatCompletionContentPartInputAudioParam], [ChatCompletionContentPartParam].
type ChatCompletionContentPartUnionParam interface {
	implementsChatCompletionContentPartUnionParam()
}

// The type of the content part.
type ChatCompletionContentPartType string

const (
	ChatCompletionContentPartTypeText       ChatCompletionContentPartType = "text"
	ChatCompletionContentPartTypeImageURL   ChatCompletionContentPartType = "image_url"
	ChatCompletionContentPartTypeInputAudio ChatCompletionContentPartType = "input_audio"
)

func (r ChatCompletionContentPartType) IsKnown() bool {
	switch r {
	case ChatCompletionContentPartTypeText, ChatCompletionContentPartTypeImageURL, ChatCompletionContentPartTypeInputAudio:
		return true
	}
	return false
}

// Learn about [image inputs](https://platform.openai.com/docs/guides/vision).
type ChatCompletionContentPartImageParam struct {
	ImageURL param.Field[ChatCompletionContentPartImageImageURLParam] `json:"image_url,required"`
	// The type of the content part.
	Type param.Field[ChatCompletionContentPartImageType] `json:"type,required"`
}

func (r ChatCompletionContentPartImageParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ChatCompletionContentPartImageParam) implementsChatCompletionContentPartUnionParam() {}

type ChatCompletionContentPartImageImageURLParam struct {
	// Either a URL of the image or the base64 encoded image data.
	URL param.Field[string] `json:"url,required" format:"uri"`
	// Specifies the detail level of the image. Learn more in the
	// [Vision guide](https://platform.openai.com/docs/guides/vision#low-or-high-fidelity-image-understanding).
	Detail param.Field[ChatCompletionContentPartImageImageURLDetail] `json:"detail"`
}

func (r ChatCompletionContentPartImageImageURLParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Specifies the detail level of the image. Learn more in the
// [Vision guide](https://platform.openai.com/docs/guides/vision#low-or-high-fidelity-image-understanding).
type ChatCompletionContentPartImageImageURLDetail string

const (
	ChatCompletionContentPartImageImageURLDetailAuto ChatCompletionContentPartImageImageURLDetail = "auto"
	ChatCompletionContentPartImageImageURLDetailLow  ChatCompletionContentPartImageImageURLDetail = "low"
	ChatCompletionContentPartImageImageURLDetailHigh ChatCompletionContentPartImageImageURLDetail = "high"
)

func (r ChatCompletionContentPartImageImageURLDetail) IsKnown() bool {
	switch r {
	case ChatCompletionContentPartImageImageURLDetailAuto, ChatCompletionContentPartImageImageURLDetailLow, ChatCompletionContentPartImageImageURLDetailHigh:
		return true
	}
	return false
}

// The type of the content part.
type ChatCompletionContentPartImageType string

const (
	ChatCompletionContentPartImageTypeImageURL ChatCompletionContentPartImageType = "image_url"
)

func (r ChatCompletionContentPartImageType) IsKnown() bool {
	switch r {
	case ChatCompletionContentPartImageTypeImageURL:
		return true
	}
	return false
}

// Learn about [audio inputs](https://platform.openai.com/docs/guides/audio).
type ChatCompletionContentPartInputAudioParam struct {
	InputAudio param.Field[ChatCompletionContentPartInputAudioInputAudioParam] `json:"input_audio,required"`
	// The type of the content part. Always `input_audio`.
	Type param.Field[ChatCompletionContentPartInputAudioType] `json:"type,required"`
}

func (r ChatCompletionContentPartInputAudioParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ChatCompletionContentPartInputAudioParam) implementsChatCompletionContentPartUnionParam() {}

type ChatCompletionContentPartInputAudioInputAudioParam struct {
	// Base64 encoded audio data.
	Data param.Field[string] `json:"data,required"`
	// The format of the encoded audio data. Currently supports "wav" and "mp3".
	Format param.Field[ChatCompletionContentPartInputAudioInputAudioFormat] `json:"format,required"`
}

func (r ChatCompletionContentPartInputAudioInputAudioParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// The format of the encoded audio data. Currently supports "wav" and "mp3".
type ChatCompletionContentPartInputAudioInputAudioFormat string

const (
	ChatCompletionContentPartInputAudioInputAudioFormatWAV ChatCompletionContentPartInputAudioInputAudioFormat = "wav"
	ChatCompletionContentPartInputAudioInputAudioFormatMP3 ChatCompletionContentPartInputAudioInputAudioFormat = "mp3"
)

func (r ChatCompletionContentPartInputAudioInputAudioFormat) IsKnown() bool {
	switch r {
	case ChatCompletionContentPartInputAudioInputAudioFormatWAV, ChatCompletionContentPartInputAudioInputAudioFormatMP3:
		return true
	}
	return false
}

// The type of the content part. Always `input_audio`.
type ChatCompletionContentPartInputAudioType string

const (
	ChatCompletionContentPartInputAudioTypeInputAudio ChatCompletionContentPartInputAudioType = "input_audio"
)

func (r ChatCompletionContentPartInputAudioType) IsKnown() bool {
	switch r {
	case ChatCompletionContentPartInputAudioTypeInputAudio:
		return true
	}
	return false
}

type ChatCompletionContentPartRefusalParam struct {
	// The refusal message generated by the model.
	Refusal param.Field[string] `json:"refusal,required"`
	// The type of the content part.
	Type param.Field[ChatCompletionContentPartRefusalType] `json:"type,required"`
}

func (r ChatCompletionContentPartRefusalParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ChatCompletionContentPartRefusalParam) implementsChatCompletionAssistantMessageParamContentUnion() {
}

// The type of the content part.
type ChatCompletionContentPartRefusalType string

const (
	ChatCompletionContentPartRefusalTypeRefusal ChatCompletionContentPartRefusalType = "refusal"
)

func (r ChatCompletionContentPartRefusalType) IsKnown() bool {
	switch r {
	case ChatCompletionContentPartRefusalTypeRefusal:
		return true
	}
	return false
}

// Learn about
// [text inputs](https://platform.openai.com/docs/guides/text-generation).
type ChatCompletionContentPartTextParam struct {
	// The text content.
	Text param.Field[string] `json:"text,required"`
	// The type of the content part.
	Type param.Field[ChatCompletionContentPartTextType] `json:"type,required"`
}

func (r ChatCompletionContentPartTextParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ChatCompletionContentPartTextParam) implementsChatCompletionAssistantMessageParamContentUnion() {
}

func (r ChatCompletionContentPartTextParam) implementsChatCompletionContentPartUnionParam() {}

// The type of the content part.
type ChatCompletionContentPartTextType string

const (
	ChatCompletionContentPartTextTypeText ChatCompletionContentPartTextType = "text"
)

func (r ChatCompletionContentPartTextType) IsKnown() bool {
	switch r {
	case ChatCompletionContentPartTextTypeText:
		return true
	}
	return false
}

type ChatCompletionDeleted struct {
	// The ID of the chat completion that was deleted.
	ID string `json:"id,required"`
	// Whether the chat completion was deleted.
	Deleted bool `json:"deleted,required"`
	// The type of object being deleted.
	Object ChatCompletionDeletedObject `json:"object,required"`
	JSON   chatCompletionDeletedJSON   `json:"-"`
}

// chatCompletionDeletedJSON contains the JSON metadata for the struct
// [ChatCompletionDeleted]
type chatCompletionDeletedJSON struct {
	ID          apijson.Field
	Deleted     apijson.Field
	Object      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ChatCompletionDeleted) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r chatCompletionDeletedJSON) RawJSON() string {
	return r.raw
}

// The type of object being deleted.
type ChatCompletionDeletedObject string

const (
	ChatCompletionDeletedObjectChatCompletionDeleted ChatCompletionDeletedObject = "chat.completion.deleted"
)

func (r ChatCompletionDeletedObject) IsKnown() bool {
	switch r {
	case ChatCompletionDeletedObjectChatCompletionDeleted:
		return true
	}
	return false
}

// Developer-provided instructions that the model should follow, regardless of
// messages sent by the user. With o1 models and newer, `developer` messages
// replace the previous `system` messages.
type ChatCompletionDeveloperMessageParam struct {
	// The contents of the developer message.
	Content param.Field[[]ChatCompletionContentPartTextParam] `json:"content,required"`
	// The role of the messages author, in this case `developer`.
	Role param.Field[ChatCompletionDeveloperMessageParamRole] `json:"role,required"`
	// An optional name for the participant. Provides the model information to
	// differentiate between participants of the same role.
	Name param.Field[string] `json:"name"`
}

func (r ChatCompletionDeveloperMessageParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ChatCompletionDeveloperMessageParam) implementsChatCompletionMessageParamUnion() {}

// The role of the messages author, in this case `developer`.
type ChatCompletionDeveloperMessageParamRole string

const (
	ChatCompletionDeveloperMessageParamRoleDeveloper ChatCompletionDeveloperMessageParamRole = "developer"
)

func (r ChatCompletionDeveloperMessageParamRole) IsKnown() bool {
	switch r {
	case ChatCompletionDeveloperMessageParamRoleDeveloper:
		return true
	}
	return false
}

// Specifying a particular function via `{"name": "my_function"}` forces the model
// to call that function.
type ChatCompletionFunctionCallOptionParam struct {
	// The name of the function to call.
	Name param.Field[string] `json:"name,required"`
}

func (r ChatCompletionFunctionCallOptionParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ChatCompletionFunctionCallOptionParam) implementsChatCompletionNewParamsFunctionCallUnion() {}

// Deprecated: deprecated
type ChatCompletionFunctionMessageParam struct {
	// The contents of the function message.
	Content param.Field[string] `json:"content,required"`
	// The name of the function to call.
	Name param.Field[string] `json:"name,required"`
	// The role of the messages author, in this case `function`.
	Role param.Field[ChatCompletionFunctionMessageParamRole] `json:"role,required"`
}

func (r ChatCompletionFunctionMessageParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ChatCompletionFunctionMessageParam) implementsChatCompletionMessageParamUnion() {}

// The role of the messages author, in this case `function`.
type ChatCompletionFunctionMessageParamRole string

const (
	ChatCompletionFunctionMessageParamRoleFunction ChatCompletionFunctionMessageParamRole = "function"
)

func (r ChatCompletionFunctionMessageParamRole) IsKnown() bool {
	switch r {
	case ChatCompletionFunctionMessageParamRoleFunction:
		return true
	}
	return false
}

// A chat completion message generated by the model.
type ChatCompletionMessage struct {
	// The contents of the message.
	Content string `json:"content,required,nullable"`
	// The refusal message generated by the model.
	Refusal string `json:"refusal,required,nullable"`
	// The role of the author of this message.
	Role ChatCompletionMessageRole `json:"role,required"`
	// If the audio output modality is requested, this object contains data about the
	// audio response from the model.
	// [Learn more](https://platform.openai.com/docs/guides/audio).
	Audio ChatCompletionAudio `json:"audio,nullable"`
	// Deprecated and replaced by `tool_calls`. The name and arguments of a function
	// that should be called, as generated by the model.
	//
	// Deprecated: deprecated
	FunctionCall ChatCompletionMessageFunctionCall `json:"function_call"`
	// The tool calls generated by the model, such as function calls.
	ToolCalls []ChatCompletionMessageToolCall `json:"tool_calls"`
	JSON      chatCompletionMessageJSON       `json:"-"`
}

// chatCompletionMessageJSON contains the JSON metadata for the struct
// [ChatCompletionMessage]
type chatCompletionMessageJSON struct {
	Content      apijson.Field
	Refusal      apijson.Field
	Role         apijson.Field
	Audio        apijson.Field
	FunctionCall apijson.Field
	ToolCalls    apijson.Field
	raw          string
	ExtraFields  map[string]apijson.Field
}

func (r *ChatCompletionMessage) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r ChatCompletionMessage) MarshalJSON() (data []byte, err error) {
	s := ""
	s, _ = sjson.Set(s, "role", r.Role)

	if r.FunctionCall.Name != "" {
		b, err := apijson.Marshal(r.FunctionCall)
		if err != nil {
			return nil, err
		}
		s, _ = sjson.SetRaw(s, "function_call", string(b))
	} else if len(r.ToolCalls) > 0 {
		b, err := apijson.Marshal(r.ToolCalls)
		if err != nil {
			return nil, err
		}
		s, _ = sjson.SetRaw(s, "tool_calls", string(b))
	} else {
		s, _ = sjson.Set(s, "content", r.Content)
	}

	return []byte(s), nil
}

func (r chatCompletionMessageJSON) RawJSON() string {
	return r.raw
}

func (r ChatCompletionMessage) implementsChatCompletionMessageParamUnion() {}

// The role of the author of this message.
type ChatCompletionMessageRole string

const (
	ChatCompletionMessageRoleAssistant ChatCompletionMessageRole = "assistant"
)

func (r ChatCompletionMessageRole) IsKnown() bool {
	switch r {
	case ChatCompletionMessageRoleAssistant:
		return true
	}
	return false
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
	Arguments string `json:"arguments,required"`
	// The name of the function to call.
	Name string                                `json:"name,required"`
	JSON chatCompletionMessageFunctionCallJSON `json:"-"`
}

// chatCompletionMessageFunctionCallJSON contains the JSON metadata for the struct
// [ChatCompletionMessageFunctionCall]
type chatCompletionMessageFunctionCallJSON struct {
	Arguments   apijson.Field
	Name        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ChatCompletionMessageFunctionCall) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r chatCompletionMessageFunctionCallJSON) RawJSON() string {
	return r.raw
}

// Developer-provided instructions that the model should follow, regardless of
// messages sent by the user. With o1 models and newer, `developer` messages
// replace the previous `system` messages.
type ChatCompletionMessageParam struct {
	// The role of the messages author, in this case `developer`.
	Role         param.Field[ChatCompletionMessageParamRole] `json:"role,required"`
	Audio        param.Field[interface{}]                    `json:"audio"`
	Content      param.Field[interface{}]                    `json:"content"`
	FunctionCall param.Field[interface{}]                    `json:"function_call"`
	// An optional name for the participant. Provides the model information to
	// differentiate between participants of the same role.
	Name param.Field[string] `json:"name"`
	// The refusal message by the assistant.
	Refusal param.Field[string] `json:"refusal"`
	// Tool call that this message is responding to.
	ToolCallID param.Field[string]      `json:"tool_call_id"`
	ToolCalls  param.Field[interface{}] `json:"tool_calls"`
}

func (r ChatCompletionMessageParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ChatCompletionMessageParam) implementsChatCompletionMessageParamUnion() {}

// Developer-provided instructions that the model should follow, regardless of
// messages sent by the user. With o1 models and newer, `developer` messages
// replace the previous `system` messages.
//
// Satisfied by [ChatCompletionDeveloperMessageParam],
// [ChatCompletionSystemMessageParam], [ChatCompletionUserMessageParam],
// [ChatCompletionAssistantMessageParam], [ChatCompletionToolMessageParam],
// [ChatCompletionFunctionMessageParam], [ChatCompletionMessageParam].
//
// This union is additionally satisfied by the return types [ChatCompletionMessage]
type ChatCompletionMessageParamUnion interface {
	implementsChatCompletionMessageParamUnion()
}

// The role of the messages author, in this case `developer`.
type ChatCompletionMessageParamRole string

const (
	ChatCompletionMessageParamRoleDeveloper ChatCompletionMessageParamRole = "developer"
	ChatCompletionMessageParamRoleSystem    ChatCompletionMessageParamRole = "system"
	ChatCompletionMessageParamRoleUser      ChatCompletionMessageParamRole = "user"
	ChatCompletionMessageParamRoleAssistant ChatCompletionMessageParamRole = "assistant"
	ChatCompletionMessageParamRoleTool      ChatCompletionMessageParamRole = "tool"
	ChatCompletionMessageParamRoleFunction  ChatCompletionMessageParamRole = "function"
)

func (r ChatCompletionMessageParamRole) IsKnown() bool {
	switch r {
	case ChatCompletionMessageParamRoleDeveloper, ChatCompletionMessageParamRoleSystem, ChatCompletionMessageParamRoleUser, ChatCompletionMessageParamRoleAssistant, ChatCompletionMessageParamRoleTool, ChatCompletionMessageParamRoleFunction:
		return true
	}
	return false
}

type ChatCompletionMessageToolCall struct {
	// The ID of the tool call.
	ID string `json:"id,required"`
	// The function that the model called.
	Function ChatCompletionMessageToolCallFunction `json:"function,required"`
	// The type of the tool. Currently, only `function` is supported.
	Type ChatCompletionMessageToolCallType `json:"type,required"`
	JSON chatCompletionMessageToolCallJSON `json:"-"`
}

// chatCompletionMessageToolCallJSON contains the JSON metadata for the struct
// [ChatCompletionMessageToolCall]
type chatCompletionMessageToolCallJSON struct {
	ID          apijson.Field
	Function    apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ChatCompletionMessageToolCall) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r chatCompletionMessageToolCallJSON) RawJSON() string {
	return r.raw
}

// The function that the model called.
type ChatCompletionMessageToolCallFunction struct {
	// The arguments to call the function with, as generated by the model in JSON
	// format. Note that the model does not always generate valid JSON, and may
	// hallucinate parameters not defined by your function schema. Validate the
	// arguments in your code before calling your function.
	Arguments string `json:"arguments,required"`
	// The name of the function to call.
	Name string                                    `json:"name,required"`
	JSON chatCompletionMessageToolCallFunctionJSON `json:"-"`
}

// chatCompletionMessageToolCallFunctionJSON contains the JSON metadata for the
// struct [ChatCompletionMessageToolCallFunction]
type chatCompletionMessageToolCallFunctionJSON struct {
	Arguments   apijson.Field
	Name        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ChatCompletionMessageToolCallFunction) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r chatCompletionMessageToolCallFunctionJSON) RawJSON() string {
	return r.raw
}

// The type of the tool. Currently, only `function` is supported.
type ChatCompletionMessageToolCallType string

const (
	ChatCompletionMessageToolCallTypeFunction ChatCompletionMessageToolCallType = "function"
)

func (r ChatCompletionMessageToolCallType) IsKnown() bool {
	switch r {
	case ChatCompletionMessageToolCallTypeFunction:
		return true
	}
	return false
}

type ChatCompletionMessageToolCallParam struct {
	// The ID of the tool call.
	ID param.Field[string] `json:"id,required"`
	// The function that the model called.
	Function param.Field[ChatCompletionMessageToolCallFunctionParam] `json:"function,required"`
	// The type of the tool. Currently, only `function` is supported.
	Type param.Field[ChatCompletionMessageToolCallType] `json:"type,required"`
}

func (r ChatCompletionMessageToolCallParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// The function that the model called.
type ChatCompletionMessageToolCallFunctionParam struct {
	// The arguments to call the function with, as generated by the model in JSON
	// format. Note that the model does not always generate valid JSON, and may
	// hallucinate parameters not defined by your function schema. Validate the
	// arguments in your code before calling your function.
	Arguments param.Field[string] `json:"arguments,required"`
	// The name of the function to call.
	Name param.Field[string] `json:"name,required"`
}

func (r ChatCompletionMessageToolCallFunctionParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type ChatCompletionModality string

const (
	ChatCompletionModalityText  ChatCompletionModality = "text"
	ChatCompletionModalityAudio ChatCompletionModality = "audio"
)

func (r ChatCompletionModality) IsKnown() bool {
	switch r {
	case ChatCompletionModalityText, ChatCompletionModalityAudio:
		return true
	}
	return false
}

// Specifies a tool the model should use. Use to force the model to call a specific
// function.
type ChatCompletionNamedToolChoiceParam struct {
	Function param.Field[ChatCompletionNamedToolChoiceFunctionParam] `json:"function,required"`
	// The type of the tool. Currently, only `function` is supported.
	Type param.Field[ChatCompletionNamedToolChoiceType] `json:"type,required"`
}

func (r ChatCompletionNamedToolChoiceParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ChatCompletionNamedToolChoiceParam) implementsChatCompletionToolChoiceOptionUnionParam() {}

type ChatCompletionNamedToolChoiceFunctionParam struct {
	// The name of the function to call.
	Name param.Field[string] `json:"name,required"`
}

func (r ChatCompletionNamedToolChoiceFunctionParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// The type of the tool. Currently, only `function` is supported.
type ChatCompletionNamedToolChoiceType string

const (
	ChatCompletionNamedToolChoiceTypeFunction ChatCompletionNamedToolChoiceType = "function"
)

func (r ChatCompletionNamedToolChoiceType) IsKnown() bool {
	switch r {
	case ChatCompletionNamedToolChoiceTypeFunction:
		return true
	}
	return false
}

// Static predicted output content, such as the content of a text file that is
// being regenerated.
type ChatCompletionPredictionContentParam struct {
	// The content that should be matched when generating a model response. If
	// generated tokens would match this content, the entire model response can be
	// returned much more quickly.
	Content param.Field[[]ChatCompletionContentPartTextParam] `json:"content,required"`
	// The type of the predicted content you want to provide. This type is currently
	// always `content`.
	Type param.Field[ChatCompletionPredictionContentType] `json:"type,required"`
}

func (r ChatCompletionPredictionContentParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// The type of the predicted content you want to provide. This type is currently
// always `content`.
type ChatCompletionPredictionContentType string

const (
	ChatCompletionPredictionContentTypeContent ChatCompletionPredictionContentType = "content"
)

func (r ChatCompletionPredictionContentType) IsKnown() bool {
	switch r {
	case ChatCompletionPredictionContentTypeContent:
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
type ChatCompletionReasoningEffort string

const (
	ChatCompletionReasoningEffortLow    ChatCompletionReasoningEffort = "low"
	ChatCompletionReasoningEffortMedium ChatCompletionReasoningEffort = "medium"
	ChatCompletionReasoningEffortHigh   ChatCompletionReasoningEffort = "high"
)

func (r ChatCompletionReasoningEffort) IsKnown() bool {
	switch r {
	case ChatCompletionReasoningEffortLow, ChatCompletionReasoningEffortMedium, ChatCompletionReasoningEffortHigh:
		return true
	}
	return false
}

// A chat completion message generated by the model.
type ChatCompletionStoreMessage struct {
	// The identifier of the chat message.
	ID   string                         `json:"id,required"`
	JSON chatCompletionStoreMessageJSON `json:"-"`
	ChatCompletionMessage
}

// chatCompletionStoreMessageJSON contains the JSON metadata for the struct
// [ChatCompletionStoreMessage]
type chatCompletionStoreMessageJSON struct {
	ID          apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ChatCompletionStoreMessage) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r chatCompletionStoreMessageJSON) RawJSON() string {
	return r.raw
}

// Options for streaming response. Only set this when you set `stream: true`.
type ChatCompletionStreamOptionsParam struct {
	// If set, an additional chunk will be streamed before the `data: [DONE]` message.
	// The `usage` field on this chunk shows the token usage statistics for the entire
	// request, and the `choices` field will always be an empty array. All other chunks
	// will also include a `usage` field, but with a null value.
	IncludeUsage param.Field[bool] `json:"include_usage"`
}

func (r ChatCompletionStreamOptionsParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Developer-provided instructions that the model should follow, regardless of
// messages sent by the user. With o1 models and newer, use `developer` messages
// for this purpose instead.
type ChatCompletionSystemMessageParam struct {
	// The contents of the system message.
	Content param.Field[[]ChatCompletionContentPartTextParam] `json:"content,required"`
	// The role of the messages author, in this case `system`.
	Role param.Field[ChatCompletionSystemMessageParamRole] `json:"role,required"`
	// An optional name for the participant. Provides the model information to
	// differentiate between participants of the same role.
	Name param.Field[string] `json:"name"`
}

func (r ChatCompletionSystemMessageParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ChatCompletionSystemMessageParam) implementsChatCompletionMessageParamUnion() {}

// The role of the messages author, in this case `system`.
type ChatCompletionSystemMessageParamRole string

const (
	ChatCompletionSystemMessageParamRoleSystem ChatCompletionSystemMessageParamRole = "system"
)

func (r ChatCompletionSystemMessageParamRole) IsKnown() bool {
	switch r {
	case ChatCompletionSystemMessageParamRoleSystem:
		return true
	}
	return false
}

type ChatCompletionTokenLogprob struct {
	// The token.
	Token string `json:"token,required"`
	// A list of integers representing the UTF-8 bytes representation of the token.
	// Useful in instances where characters are represented by multiple tokens and
	// their byte representations must be combined to generate the correct text
	// representation. Can be `null` if there is no bytes representation for the token.
	Bytes []int64 `json:"bytes,required,nullable"`
	// The log probability of this token, if it is within the top 20 most likely
	// tokens. Otherwise, the value `-9999.0` is used to signify that the token is very
	// unlikely.
	Logprob float64 `json:"logprob,required"`
	// List of the most likely tokens and their log probability, at this token
	// position. In rare cases, there may be fewer than the number of requested
	// `top_logprobs` returned.
	TopLogprobs []ChatCompletionTokenLogprobTopLogprob `json:"top_logprobs,required"`
	JSON        chatCompletionTokenLogprobJSON         `json:"-"`
}

// chatCompletionTokenLogprobJSON contains the JSON metadata for the struct
// [ChatCompletionTokenLogprob]
type chatCompletionTokenLogprobJSON struct {
	Token       apijson.Field
	Bytes       apijson.Field
	Logprob     apijson.Field
	TopLogprobs apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ChatCompletionTokenLogprob) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r chatCompletionTokenLogprobJSON) RawJSON() string {
	return r.raw
}

type ChatCompletionTokenLogprobTopLogprob struct {
	// The token.
	Token string `json:"token,required"`
	// A list of integers representing the UTF-8 bytes representation of the token.
	// Useful in instances where characters are represented by multiple tokens and
	// their byte representations must be combined to generate the correct text
	// representation. Can be `null` if there is no bytes representation for the token.
	Bytes []int64 `json:"bytes,required,nullable"`
	// The log probability of this token, if it is within the top 20 most likely
	// tokens. Otherwise, the value `-9999.0` is used to signify that the token is very
	// unlikely.
	Logprob float64                                  `json:"logprob,required"`
	JSON    chatCompletionTokenLogprobTopLogprobJSON `json:"-"`
}

// chatCompletionTokenLogprobTopLogprobJSON contains the JSON metadata for the
// struct [ChatCompletionTokenLogprobTopLogprob]
type chatCompletionTokenLogprobTopLogprobJSON struct {
	Token       apijson.Field
	Bytes       apijson.Field
	Logprob     apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ChatCompletionTokenLogprobTopLogprob) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r chatCompletionTokenLogprobTopLogprobJSON) RawJSON() string {
	return r.raw
}

type ChatCompletionToolParam struct {
	Function param.Field[shared.FunctionDefinitionParam] `json:"function,required"`
	// The type of the tool. Currently, only `function` is supported.
	Type param.Field[ChatCompletionToolType] `json:"type,required"`
}

func (r ChatCompletionToolParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// The type of the tool. Currently, only `function` is supported.
type ChatCompletionToolType string

const (
	ChatCompletionToolTypeFunction ChatCompletionToolType = "function"
)

func (r ChatCompletionToolType) IsKnown() bool {
	switch r {
	case ChatCompletionToolTypeFunction:
		return true
	}
	return false
}

// Controls which (if any) tool is called by the model. `none` means the model will
// not call any tool and instead generates a message. `auto` means the model can
// pick between generating a message or calling one or more tools. `required` means
// the model must call one or more tools. Specifying a particular tool via
// `{"type": "function", "function": {"name": "my_function"}}` forces the model to
// call that tool.
//
// `none` is the default when no tools are present. `auto` is the default if tools
// are present.
//
// Satisfied by [ChatCompletionToolChoiceOptionAuto],
// [ChatCompletionNamedToolChoiceParam].
type ChatCompletionToolChoiceOptionUnionParam interface {
	implementsChatCompletionToolChoiceOptionUnionParam()
}

// `none` means the model will not call any tool and instead generates a message.
// `auto` means the model can pick between generating a message or calling one or
// more tools. `required` means the model must call one or more tools.
type ChatCompletionToolChoiceOptionAuto string

const (
	ChatCompletionToolChoiceOptionAutoNone     ChatCompletionToolChoiceOptionAuto = "none"
	ChatCompletionToolChoiceOptionAutoAuto     ChatCompletionToolChoiceOptionAuto = "auto"
	ChatCompletionToolChoiceOptionAutoRequired ChatCompletionToolChoiceOptionAuto = "required"
)

func (r ChatCompletionToolChoiceOptionAuto) IsKnown() bool {
	switch r {
	case ChatCompletionToolChoiceOptionAutoNone, ChatCompletionToolChoiceOptionAutoAuto, ChatCompletionToolChoiceOptionAutoRequired:
		return true
	}
	return false
}

func (r ChatCompletionToolChoiceOptionAuto) implementsChatCompletionToolChoiceOptionUnionParam() {}

type ChatCompletionToolMessageParam struct {
	// The contents of the tool message.
	Content param.Field[[]ChatCompletionContentPartTextParam] `json:"content,required"`
	// The role of the messages author, in this case `tool`.
	Role param.Field[ChatCompletionToolMessageParamRole] `json:"role,required"`
	// Tool call that this message is responding to.
	ToolCallID param.Field[string] `json:"tool_call_id,required"`
}

func (r ChatCompletionToolMessageParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ChatCompletionToolMessageParam) implementsChatCompletionMessageParamUnion() {}

// The role of the messages author, in this case `tool`.
type ChatCompletionToolMessageParamRole string

const (
	ChatCompletionToolMessageParamRoleTool ChatCompletionToolMessageParamRole = "tool"
)

func (r ChatCompletionToolMessageParamRole) IsKnown() bool {
	switch r {
	case ChatCompletionToolMessageParamRoleTool:
		return true
	}
	return false
}

// Messages sent by an end user, containing prompts or additional context
// information.
type ChatCompletionUserMessageParam struct {
	// The contents of the user message.
	Content param.Field[[]ChatCompletionContentPartUnionParam] `json:"content,required"`
	// The role of the messages author, in this case `user`.
	Role param.Field[ChatCompletionUserMessageParamRole] `json:"role,required"`
	// An optional name for the participant. Provides the model information to
	// differentiate between participants of the same role.
	Name param.Field[string] `json:"name"`
}

func (r ChatCompletionUserMessageParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ChatCompletionUserMessageParam) implementsChatCompletionMessageParamUnion() {}

// The role of the messages author, in this case `user`.
type ChatCompletionUserMessageParamRole string

const (
	ChatCompletionUserMessageParamRoleUser ChatCompletionUserMessageParamRole = "user"
)

func (r ChatCompletionUserMessageParamRole) IsKnown() bool {
	switch r {
	case ChatCompletionUserMessageParamRoleUser:
		return true
	}
	return false
}

type ChatCompletionNewParams struct {
	// A list of messages comprising the conversation so far. Depending on the
	// [model](https://platform.openai.com/docs/models) you use, different message
	// types (modalities) are supported, like
	// [text](https://platform.openai.com/docs/guides/text-generation),
	// [images](https://platform.openai.com/docs/guides/vision), and
	// [audio](https://platform.openai.com/docs/guides/audio).
	Messages param.Field[[]ChatCompletionMessageParamUnion] `json:"messages,required"`
	// ID of the model to use. See the
	// [model endpoint compatibility](https://platform.openai.com/docs/models#model-endpoint-compatibility)
	// table for details on which models work with the Chat API.
	Model param.Field[shared.ChatModel] `json:"model,required"`
	// Parameters for audio output. Required when audio output is requested with
	// `modalities: ["audio"]`.
	// [Learn more](https://platform.openai.com/docs/guides/audio).
	Audio param.Field[ChatCompletionAudioParam] `json:"audio"`
	// Number between -2.0 and 2.0. Positive values penalize new tokens based on their
	// existing frequency in the text so far, decreasing the model's likelihood to
	// repeat the same line verbatim.
	FrequencyPenalty param.Field[float64] `json:"frequency_penalty"`
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
	FunctionCall param.Field[ChatCompletionNewParamsFunctionCallUnion] `json:"function_call"`
	// Deprecated in favor of `tools`.
	//
	// A list of functions the model may generate JSON inputs for.
	Functions param.Field[[]ChatCompletionNewParamsFunction] `json:"functions"`
	// Modify the likelihood of specified tokens appearing in the completion.
	//
	// Accepts a JSON object that maps tokens (specified by their token ID in the
	// tokenizer) to an associated bias value from -100 to 100. Mathematically, the
	// bias is added to the logits generated by the model prior to sampling. The exact
	// effect will vary per model, but values between -1 and 1 should decrease or
	// increase likelihood of selection; values like -100 or 100 should result in a ban
	// or exclusive selection of the relevant token.
	LogitBias param.Field[map[string]int64] `json:"logit_bias"`
	// Whether to return log probabilities of the output tokens or not. If true,
	// returns the log probabilities of each output token returned in the `content` of
	// `message`.
	Logprobs param.Field[bool] `json:"logprobs"`
	// An upper bound for the number of tokens that can be generated for a completion,
	// including visible output tokens and
	// [reasoning tokens](https://platform.openai.com/docs/guides/reasoning).
	MaxCompletionTokens param.Field[int64] `json:"max_completion_tokens"`
	// The maximum number of [tokens](/tokenizer) that can be generated in the chat
	// completion. This value can be used to control
	// [costs](https://openai.com/api/pricing/) for text generated via API.
	//
	// This value is now deprecated in favor of `max_completion_tokens`, and is not
	// compatible with
	// [o1 series models](https://platform.openai.com/docs/guides/reasoning).
	MaxTokens param.Field[int64] `json:"max_tokens"`
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard.
	//
	// Keys are strings with a maximum length of 64 characters. Values are strings with
	// a maximum length of 512 characters.
	Metadata param.Field[shared.MetadataParam] `json:"metadata"`
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
	Modalities param.Field[[]ChatCompletionModality] `json:"modalities"`
	// How many chat completion choices to generate for each input message. Note that
	// you will be charged based on the number of generated tokens across all of the
	// choices. Keep `n` as `1` to minimize costs.
	N param.Field[int64] `json:"n"`
	// Whether to enable
	// [parallel function calling](https://platform.openai.com/docs/guides/function-calling#configuring-parallel-function-calling)
	// during tool use.
	ParallelToolCalls param.Field[bool] `json:"parallel_tool_calls"`
	// Static predicted output content, such as the content of a text file that is
	// being regenerated.
	Prediction param.Field[ChatCompletionPredictionContentParam] `json:"prediction"`
	// Number between -2.0 and 2.0. Positive values penalize new tokens based on
	// whether they appear in the text so far, increasing the model's likelihood to
	// talk about new topics.
	PresencePenalty param.Field[float64] `json:"presence_penalty"`
	// **o1 and o3-mini models only**
	//
	// Constrains effort on reasoning for
	// [reasoning models](https://platform.openai.com/docs/guides/reasoning). Currently
	// supported values are `low`, `medium`, and `high`. Reducing reasoning effort can
	// result in faster responses and fewer tokens used on reasoning in a response.
	ReasoningEffort param.Field[ChatCompletionReasoningEffort] `json:"reasoning_effort"`
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
	ResponseFormat param.Field[ChatCompletionNewParamsResponseFormatUnion] `json:"response_format"`
	// This feature is in Beta. If specified, our system will make a best effort to
	// sample deterministically, such that repeated requests with the same `seed` and
	// parameters should return the same result. Determinism is not guaranteed, and you
	// should refer to the `system_fingerprint` response parameter to monitor changes
	// in the backend.
	Seed param.Field[int64] `json:"seed"`
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
	ServiceTier param.Field[ChatCompletionNewParamsServiceTier] `json:"service_tier"`
	// Up to 4 sequences where the API will stop generating further tokens.
	Stop param.Field[ChatCompletionNewParamsStopUnion] `json:"stop"`
	// Whether or not to store the output of this chat completion request for use in
	// our [model distillation](https://platform.openai.com/docs/guides/distillation)
	// or [evals](https://platform.openai.com/docs/guides/evals) products.
	Store param.Field[bool] `json:"store"`
	// Options for streaming response. Only set this when you set `stream: true`.
	StreamOptions param.Field[ChatCompletionStreamOptionsParam] `json:"stream_options"`
	// What sampling temperature to use, between 0 and 2. Higher values like 0.8 will
	// make the output more random, while lower values like 0.2 will make it more
	// focused and deterministic. We generally recommend altering this or `top_p` but
	// not both.
	Temperature param.Field[float64] `json:"temperature"`
	// Controls which (if any) tool is called by the model. `none` means the model will
	// not call any tool and instead generates a message. `auto` means the model can
	// pick between generating a message or calling one or more tools. `required` means
	// the model must call one or more tools. Specifying a particular tool via
	// `{"type": "function", "function": {"name": "my_function"}}` forces the model to
	// call that tool.
	//
	// `none` is the default when no tools are present. `auto` is the default if tools
	// are present.
	ToolChoice param.Field[ChatCompletionToolChoiceOptionUnionParam] `json:"tool_choice"`
	// A list of tools the model may call. Currently, only functions are supported as a
	// tool. Use this to provide a list of functions the model may generate JSON inputs
	// for. A max of 128 functions are supported.
	Tools param.Field[[]ChatCompletionToolParam] `json:"tools"`
	// An integer between 0 and 20 specifying the number of most likely tokens to
	// return at each token position, each with an associated log probability.
	// `logprobs` must be set to `true` if this parameter is used.
	TopLogprobs param.Field[int64] `json:"top_logprobs"`
	// An alternative to sampling with temperature, called nucleus sampling, where the
	// model considers the results of the tokens with top_p probability mass. So 0.1
	// means only the tokens comprising the top 10% probability mass are considered.
	//
	// We generally recommend altering this or `temperature` but not both.
	TopP param.Field[float64] `json:"top_p"`
	// A unique identifier representing your end-user, which can help OpenAI to monitor
	// and detect abuse.
	// [Learn more](https://platform.openai.com/docs/guides/safety-best-practices#end-user-ids).
	User param.Field[string] `json:"user"`
}

func (r ChatCompletionNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

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
//
// Satisfied by [ChatCompletionNewParamsFunctionCallAuto],
// [ChatCompletionFunctionCallOptionParam].
//
// Deprecated: deprecated
type ChatCompletionNewParamsFunctionCallUnion interface {
	implementsChatCompletionNewParamsFunctionCallUnion()
}

// `none` means the model will not call a function and instead generates a message.
// `auto` means the model can pick between generating a message or calling a
// function.
type ChatCompletionNewParamsFunctionCallAuto string

const (
	ChatCompletionNewParamsFunctionCallAutoNone ChatCompletionNewParamsFunctionCallAuto = "none"
	ChatCompletionNewParamsFunctionCallAutoAuto ChatCompletionNewParamsFunctionCallAuto = "auto"
)

func (r ChatCompletionNewParamsFunctionCallAuto) IsKnown() bool {
	switch r {
	case ChatCompletionNewParamsFunctionCallAutoNone, ChatCompletionNewParamsFunctionCallAutoAuto:
		return true
	}
	return false
}

func (r ChatCompletionNewParamsFunctionCallAuto) implementsChatCompletionNewParamsFunctionCallUnion() {
}

// Deprecated: deprecated
type ChatCompletionNewParamsFunction struct {
	// The name of the function to be called. Must be a-z, A-Z, 0-9, or contain
	// underscores and dashes, with a maximum length of 64.
	Name param.Field[string] `json:"name,required"`
	// A description of what the function does, used by the model to choose when and
	// how to call the function.
	Description param.Field[string] `json:"description"`
	// The parameters the functions accepts, described as a JSON Schema object. See the
	// [guide](https://platform.openai.com/docs/guides/function-calling) for examples,
	// and the
	// [JSON Schema reference](https://json-schema.org/understanding-json-schema/) for
	// documentation about the format.
	//
	// Omitting `parameters` defines a function with an empty parameter list.
	Parameters param.Field[shared.FunctionParameters] `json:"parameters"`
}

func (r ChatCompletionNewParamsFunction) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

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
type ChatCompletionNewParamsResponseFormat struct {
	// The type of response format being defined: `text`
	Type       param.Field[ChatCompletionNewParamsResponseFormatType] `json:"type,required"`
	JSONSchema param.Field[interface{}]                               `json:"json_schema"`
}

func (r ChatCompletionNewParamsResponseFormat) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ChatCompletionNewParamsResponseFormat) ImplementsChatCompletionNewParamsResponseFormatUnion() {
}

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
//
// Satisfied by [shared.ResponseFormatTextParam],
// [shared.ResponseFormatJSONObjectParam], [shared.ResponseFormatJSONSchemaParam],
// [ChatCompletionNewParamsResponseFormat].
type ChatCompletionNewParamsResponseFormatUnion interface {
	ImplementsChatCompletionNewParamsResponseFormatUnion()
}

// The type of response format being defined: `text`
type ChatCompletionNewParamsResponseFormatType string

const (
	ChatCompletionNewParamsResponseFormatTypeText       ChatCompletionNewParamsResponseFormatType = "text"
	ChatCompletionNewParamsResponseFormatTypeJSONObject ChatCompletionNewParamsResponseFormatType = "json_object"
	ChatCompletionNewParamsResponseFormatTypeJSONSchema ChatCompletionNewParamsResponseFormatType = "json_schema"
)

func (r ChatCompletionNewParamsResponseFormatType) IsKnown() bool {
	switch r {
	case ChatCompletionNewParamsResponseFormatTypeText, ChatCompletionNewParamsResponseFormatTypeJSONObject, ChatCompletionNewParamsResponseFormatTypeJSONSchema:
		return true
	}
	return false
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

func (r ChatCompletionNewParamsServiceTier) IsKnown() bool {
	switch r {
	case ChatCompletionNewParamsServiceTierAuto, ChatCompletionNewParamsServiceTierDefault:
		return true
	}
	return false
}

// Up to 4 sequences where the API will stop generating further tokens.
//
// Satisfied by [shared.UnionString], [ChatCompletionNewParamsStopArray].
type ChatCompletionNewParamsStopUnion interface {
	ImplementsChatCompletionNewParamsStopUnion()
}

type ChatCompletionNewParamsStopArray []string

func (r ChatCompletionNewParamsStopArray) ImplementsChatCompletionNewParamsStopUnion() {}

type ChatCompletionUpdateParams struct {
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard.
	//
	// Keys are strings with a maximum length of 64 characters. Values are strings with
	// a maximum length of 512 characters.
	Metadata param.Field[shared.MetadataParam] `json:"metadata,required"`
}

func (r ChatCompletionUpdateParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type ChatCompletionListParams struct {
	// Identifier for the last chat completion from the previous pagination request.
	After param.Field[string] `query:"after"`
	// Number of chat completions to retrieve.
	Limit param.Field[int64] `query:"limit"`
	// A list of metadata keys to filter the chat completions by. Example:
	//
	// `metadata[key1]=value1&metadata[key2]=value2`
	Metadata param.Field[shared.MetadataParam] `query:"metadata"`
	// The model used to generate the chat completions.
	Model param.Field[string] `query:"model"`
	// Sort order for chat completions by timestamp. Use `asc` for ascending order or
	// `desc` for descending order. Defaults to `asc`.
	Order param.Field[ChatCompletionListParamsOrder] `query:"order"`
}

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

func (r ChatCompletionListParamsOrder) IsKnown() bool {
	switch r {
	case ChatCompletionListParamsOrderAsc, ChatCompletionListParamsOrderDesc:
		return true
	}
	return false
}

// File generated from our OpenAPI spec by Castiron. See CONTRIBUTING.md for details.

package realtime

import (
	"encoding/json"

	"github.com/openai/openai-go/v3/internal/apijson"
	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/packages/respjson"
	"github.com/openai/openai-go/v3/responses"
	"github.com/openai/openai-go/v3/shared"
	"github.com/openai/openai-go/v3/shared/constant"
	"github.com/tidwall/gjson"
)

// RealtimeService contains methods and other services that help with interacting
// with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewRealtimeService] method instead.
type RealtimeService struct {
	Options       []option.RequestOption
	ClientSecrets ClientSecretService
	Calls         CallService
	Translations  TranslationService
}

// NewRealtimeService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewRealtimeService(opts ...option.RequestOption) (r RealtimeService) {
	r = RealtimeService{}
	r.Options = requestconfig.InheritedOptions(opts...)
	r.ClientSecrets = NewClientSecretService(opts...)
	r.Calls = NewCallService(opts...)
	r.Translations = NewTranslationService(opts...)
	return
}

type AudioTranscription struct {
	// Controls how long the model waits before emitting transcription text. Higher
	// values can improve transcription accuracy at the cost of latency. Only supported
	// with `gpt-realtime-whisper` in GA Realtime sessions.
	//
	// Any of "minimal", "low", "medium", "high", "xhigh".
	Delay AudioTranscriptionDelay `json:"delay"`
	// Words or phrases to guide transcription of the input audio. Supported by
	// `gpt-transcribe` and `gpt-live-transcribe`.
	Keywords []string `json:"keywords"`
	// The language of the input audio. Supplying the input language in
	// [ISO-639-1](https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes) (e.g. `en`)
	// format will improve accuracy and latency.
	Language string `json:"language"`
	// Possible languages of the input audio, in
	// [ISO-639-1](https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes) format.
	// Supported by `gpt-transcribe` and `gpt-live-transcribe`.
	Languages []string `json:"languages"`
	// The model to use for transcription. Current options are `whisper-1`,
	// `gpt-transcribe`, `gpt-live-transcribe`, `gpt-4o-mini-transcribe`,
	// `gpt-4o-mini-transcribe-2025-12-15`, `gpt-4o-transcribe`,
	// `gpt-4o-transcribe-diarize`, and `gpt-realtime-whisper`. Use
	// `gpt-4o-transcribe-diarize` when you need diarization with speaker labels.
	Model AudioTranscriptionModel `json:"model"`
	// An optional text to guide the model's style or continue a previous audio
	// segment. For `whisper-1`, the
	// [prompt is a list of keywords](https://developers.openai.com/api/docs/guides/speech-to-text#prompting).
	// For `gpt-4o-transcribe` models (excluding `gpt-4o-transcribe-diarize`), the
	// prompt is a free text string, for example "expect words related to technology".
	// Prompt is not supported with `gpt-realtime-whisper` in GA Realtime sessions.
	Prompt string `json:"prompt"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Delay       respjson.Field
		Keywords    respjson.Field
		Language    respjson.Field
		Languages   respjson.Field
		Model       respjson.Field
		Prompt      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AudioTranscription) RawJSON() string { return r.JSON.raw }
func (r *AudioTranscription) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this AudioTranscription to a AudioTranscriptionParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// AudioTranscriptionParam.Overrides()
func (r AudioTranscription) ToParam() AudioTranscriptionParam {
	return param.Override[AudioTranscriptionParam](json.RawMessage(r.RawJSON()))
}

// Controls how long the model waits before emitting transcription text. Higher
// values can improve transcription accuracy at the cost of latency. Only supported
// with `gpt-realtime-whisper` in GA Realtime sessions.
type AudioTranscriptionDelay string

const (
	AudioTranscriptionDelayMinimal AudioTranscriptionDelay = "minimal"
	AudioTranscriptionDelayLow     AudioTranscriptionDelay = "low"
	AudioTranscriptionDelayMedium  AudioTranscriptionDelay = "medium"
	AudioTranscriptionDelayHigh    AudioTranscriptionDelay = "high"
	AudioTranscriptionDelayXhigh   AudioTranscriptionDelay = "xhigh"
)

// The model to use for transcription. Current options are `whisper-1`,
// `gpt-transcribe`, `gpt-live-transcribe`, `gpt-4o-mini-transcribe`,
// `gpt-4o-mini-transcribe-2025-12-15`, `gpt-4o-transcribe`,
// `gpt-4o-transcribe-diarize`, and `gpt-realtime-whisper`. Use
// `gpt-4o-transcribe-diarize` when you need diarization with speaker labels.
type AudioTranscriptionModel string

const (
	// Deprecated: Announced shutdown date: 2027-02-26. See
	// https://developers.openai.com/api/docs/deprecations for details and recommended
	// replacements.
	AudioTranscriptionModelWhisper1          AudioTranscriptionModel = "whisper-1"
	AudioTranscriptionModelGPTTranscribe     AudioTranscriptionModel = "gpt-transcribe"
	AudioTranscriptionModelGPTLiveTranscribe AudioTranscriptionModel = "gpt-live-transcribe"
	// Deprecated: Announced shutdown date: 2027-02-26. See
	// https://developers.openai.com/api/docs/deprecations for details and recommended
	// replacements.
	AudioTranscriptionModelGPT4oMiniTranscribe AudioTranscriptionModel = "gpt-4o-mini-transcribe"
	// Deprecated: Announced shutdown date: 2027-02-26. See
	// https://developers.openai.com/api/docs/deprecations for details and recommended
	// replacements.
	AudioTranscriptionModelGPT4oMiniTranscribe2025_12_15 AudioTranscriptionModel = "gpt-4o-mini-transcribe-2025-12-15"
	// Deprecated: Announced shutdown date: 2027-02-26. See
	// https://developers.openai.com/api/docs/deprecations for details and recommended
	// replacements.
	AudioTranscriptionModelGPT4oTranscribe AudioTranscriptionModel = "gpt-4o-transcribe"
	// Deprecated: Announced shutdown date: 2027-02-26. See
	// https://developers.openai.com/api/docs/deprecations for details and recommended
	// replacements.
	AudioTranscriptionModelGPT4oTranscribeDiarize AudioTranscriptionModel = "gpt-4o-transcribe-diarize"
	AudioTranscriptionModelGPTRealtimeWhisper     AudioTranscriptionModel = "gpt-realtime-whisper"
)

type AudioTranscriptionParam struct {
	// The language of the input audio. Supplying the input language in
	// [ISO-639-1](https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes) (e.g. `en`)
	// format will improve accuracy and latency.
	Language param.Opt[string] `json:"language,omitzero"`
	// An optional text to guide the model's style or continue a previous audio
	// segment. For `whisper-1`, the
	// [prompt is a list of keywords](https://developers.openai.com/api/docs/guides/speech-to-text#prompting).
	// For `gpt-4o-transcribe` models (excluding `gpt-4o-transcribe-diarize`), the
	// prompt is a free text string, for example "expect words related to technology".
	// Prompt is not supported with `gpt-realtime-whisper` in GA Realtime sessions.
	Prompt param.Opt[string] `json:"prompt,omitzero"`
	// Controls how long the model waits before emitting transcription text. Higher
	// values can improve transcription accuracy at the cost of latency. Only supported
	// with `gpt-realtime-whisper` in GA Realtime sessions.
	//
	// Any of "minimal", "low", "medium", "high", "xhigh".
	Delay AudioTranscriptionDelay `json:"delay,omitzero"`
	// Words or phrases to guide transcription of the input audio. Supported by
	// `gpt-transcribe` and `gpt-live-transcribe`.
	Keywords []string `json:"keywords,omitzero"`
	// Possible languages of the input audio, in
	// [ISO-639-1](https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes) format.
	// Supported by `gpt-transcribe` and `gpt-live-transcribe`.
	Languages []string `json:"languages,omitzero"`
	// The model to use for transcription. Current options are `whisper-1`,
	// `gpt-transcribe`, `gpt-live-transcribe`, `gpt-4o-mini-transcribe`,
	// `gpt-4o-mini-transcribe-2025-12-15`, `gpt-4o-transcribe`,
	// `gpt-4o-transcribe-diarize`, and `gpt-realtime-whisper`. Use
	// `gpt-4o-transcribe-diarize` when you need diarization with speaker labels.
	Model AudioTranscriptionModel `json:"model,omitzero"`
	paramObj
}

func (r AudioTranscriptionParam) MarshalJSON() (data []byte, err error) {
	type shadow AudioTranscriptionParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AudioTranscriptionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Returned when a conversation is created. Emitted right after session creation.
type ConversationCreatedEvent struct {
	// The conversation resource.
	Conversation ConversationCreatedEventConversation `json:"conversation" api:"required"`
	// The unique ID of the server event.
	EventID string `json:"event_id" api:"required"`
	// The event type, must be `conversation.created`.
	Type constant.ConversationCreated `json:"type" default:"conversation.created"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Conversation respjson.Field
		EventID      respjson.Field
		Type         respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ConversationCreatedEvent) RawJSON() string { return r.JSON.raw }
func (r *ConversationCreatedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The conversation resource.
type ConversationCreatedEventConversation struct {
	// The unique ID of the conversation.
	ID string `json:"id"`
	// The object type, must be `realtime.conversation`.
	//
	// Any of "realtime.conversation".
	Object string `json:"object"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Object      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ConversationCreatedEventConversation) RawJSON() string { return r.JSON.raw }
func (r *ConversationCreatedEventConversation) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ConversationItemUnion contains all possible properties and values from
// [RealtimeConversationItemSystemMessage], [RealtimeConversationItemUserMessage],
// [RealtimeConversationItemAssistantMessage],
// [RealtimeConversationItemFunctionCall],
// [RealtimeConversationItemFunctionCallOutput], [RealtimeMcpApprovalResponse],
// [RealtimeMcpListTools], [RealtimeMcpToolCall], [RealtimeMcpApprovalRequest].
//
// Use the [ConversationItemUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ConversationItemUnion struct {
	Type                                         string                                     `json:"type" api:"required"`
	OfRealtimeConversationItemSystemMessage      RealtimeConversationItemSystemMessage      `json:",inline"`
	OfRealtimeConversationItemUserMessage        RealtimeConversationItemUserMessage        `json:",inline"`
	OfRealtimeConversationItemAssistantMessage   RealtimeConversationItemAssistantMessage   `json:",inline"`
	OfRealtimeConversationItemFunctionCall       RealtimeConversationItemFunctionCall       `json:",inline"`
	OfRealtimeConversationItemFunctionCallOutput RealtimeConversationItemFunctionCallOutput `json:",inline"`
	OfRealtimeMcpApprovalResponse                RealtimeMcpApprovalResponse                `json:",inline"`
	OfRealtimeMcpListTools                       RealtimeMcpListTools                       `json:",inline"`
	OfRealtimeMcpToolCall                        RealtimeMcpToolCall                        `json:",inline"`
	OfRealtimeMcpApprovalRequest                 RealtimeMcpApprovalRequest                 `json:",inline"`
	JSON                                         struct {
		Type                                         respjson.Field
		OfRealtimeConversationItemSystemMessage      respjson.Field
		OfRealtimeConversationItemUserMessage        respjson.Field
		OfRealtimeConversationItemAssistantMessage   respjson.Field
		OfRealtimeConversationItemFunctionCall       respjson.Field
		OfRealtimeConversationItemFunctionCallOutput respjson.Field
		OfRealtimeMcpApprovalResponse                respjson.Field
		OfRealtimeMcpListTools                       respjson.Field
		OfRealtimeMcpToolCall                        respjson.Field
		OfRealtimeMcpApprovalRequest                 respjson.Field
		raw                                          string
	} `json:"-"`
}

// anyConversationItem is implemented by each variant of [ConversationItemUnion] to
// add type safety for the return type of [ConversationItemUnion.AsAny]
type anyConversationItem interface {
	implConversationItemUnion()
}

func (RealtimeConversationItemSystemMessage) implConversationItemUnion()      {}
func (RealtimeConversationItemUserMessage) implConversationItemUnion()        {}
func (RealtimeConversationItemAssistantMessage) implConversationItemUnion()   {}
func (RealtimeConversationItemFunctionCall) implConversationItemUnion()       {}
func (RealtimeConversationItemFunctionCallOutput) implConversationItemUnion() {}
func (RealtimeMcpApprovalResponse) implConversationItemUnion()                {}
func (RealtimeMcpListTools) implConversationItemUnion()                       {}
func (RealtimeMcpToolCall) implConversationItemUnion()                        {}
func (RealtimeMcpApprovalRequest) implConversationItemUnion()                 {}

// Use the following switch statement to find the correct variant
//
//	switch variant := ConversationItemUnion.AsAny().(type) {
//	case realtime.RealtimeConversationItemSystemMessage:
//	case realtime.RealtimeConversationItemUserMessage:
//	case realtime.RealtimeConversationItemAssistantMessage:
//	case realtime.RealtimeConversationItemFunctionCall:
//	case realtime.RealtimeConversationItemFunctionCallOutput:
//	case realtime.RealtimeMcpApprovalResponse:
//	case realtime.RealtimeMcpListTools:
//	case realtime.RealtimeMcpToolCall:
//	case realtime.RealtimeMcpApprovalRequest:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ConversationItemUnion) AsAny() anyConversationItem {
	switch u.Type {
	case "message":
		discriminator := gjson.Get(u.JSON.raw, apijson.EscapeSJSONKey("role"))
		if discriminator.Type != gjson.String {
			return nil
		}
		switch discriminator.Str {
		case "system":
			return u.AsRealtimeConversationItemSystemMessage()
		case "user":
			return u.AsRealtimeConversationItemUserMessage()
		case "assistant":
			return u.AsRealtimeConversationItemAssistantMessage()
		}
	case "function_call":
		return u.AsRealtimeConversationItemFunctionCall()
	case "function_call_output":
		return u.AsRealtimeConversationItemFunctionCallOutput()
	case "mcp_approval_response":
		return u.AsRealtimeMcpApprovalResponse()
	case "mcp_list_tools":
		return u.AsRealtimeMcpListTools()
	case "mcp_call":
		return u.AsRealtimeMcpToolCall()
	case "mcp_approval_request":
		return u.AsRealtimeMcpApprovalRequest()
	}
	return nil
}

func (u ConversationItemUnion) AsRealtimeConversationItemSystemMessage() (v RealtimeConversationItemSystemMessage) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ConversationItemUnion) AsRealtimeConversationItemUserMessage() (v RealtimeConversationItemUserMessage) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ConversationItemUnion) AsRealtimeConversationItemAssistantMessage() (v RealtimeConversationItemAssistantMessage) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ConversationItemUnion) AsRealtimeConversationItemFunctionCall() (v RealtimeConversationItemFunctionCall) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ConversationItemUnion) AsRealtimeConversationItemFunctionCallOutput() (v RealtimeConversationItemFunctionCallOutput) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ConversationItemUnion) AsRealtimeMcpApprovalResponse() (v RealtimeMcpApprovalResponse) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ConversationItemUnion) AsRealtimeMcpListTools() (v RealtimeMcpListTools) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ConversationItemUnion) AsRealtimeMcpToolCall() (v RealtimeMcpToolCall) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ConversationItemUnion) AsRealtimeMcpApprovalRequest() (v RealtimeMcpApprovalRequest) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ConversationItemUnion) RawJSON() string { return u.JSON.raw }

func (r *ConversationItemUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this ConversationItemUnion to a ConversationItemUnionParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// ConversationItemUnionParam.Overrides()
func (r ConversationItemUnion) ToParam() ConversationItemUnionParam {
	return param.Override[ConversationItemUnionParam](json.RawMessage(r.RawJSON()))
}

func ConversationItemParamOfRealtimeSystemMessageItem(content []RealtimeConversationItemSystemMessageContentParam) ConversationItemUnionParam {
	var system RealtimeConversationItemSystemMessageParam
	system.Content = content
	return ConversationItemUnionParam{OfRealtimeSystemMessageItem: &system}
}

func ConversationItemParamOfRealtimeUserMessageItem(content []RealtimeConversationItemUserMessageContentParam) ConversationItemUnionParam {
	var user RealtimeConversationItemUserMessageParam
	user.Content = content
	return ConversationItemUnionParam{OfRealtimeUserMessageItem: &user}
}

func ConversationItemParamOfRealtimeAssistantMessageItem(content []RealtimeConversationItemAssistantMessageContentParam) ConversationItemUnionParam {
	var assistant RealtimeConversationItemAssistantMessageParam
	assistant.Content = content
	return ConversationItemUnionParam{OfRealtimeAssistantMessageItem: &assistant}
}

func ConversationItemParamOfFunctionCall(arguments string, name string) ConversationItemUnionParam {
	var functionCall RealtimeConversationItemFunctionCallParam
	functionCall.Arguments = arguments
	functionCall.Name = name
	return ConversationItemUnionParam{OfFunctionCall: &functionCall}
}

func ConversationItemParamOfFunctionCallOutput(callID string, output string) ConversationItemUnionParam {
	var functionCallOutput RealtimeConversationItemFunctionCallOutputParam
	functionCallOutput.CallID = callID
	functionCallOutput.Output = output
	return ConversationItemUnionParam{OfFunctionCallOutput: &functionCallOutput}
}

func ConversationItemParamOfMcpApprovalResponse(approvalRequestID string, approve bool, id string) ConversationItemUnionParam {
	var mcpApprovalResponse RealtimeMcpApprovalResponseParam
	mcpApprovalResponse.ApprovalRequestID = approvalRequestID
	mcpApprovalResponse.Approve = approve
	mcpApprovalResponse.ID = id
	return ConversationItemUnionParam{OfMcpApprovalResponse: &mcpApprovalResponse}
}

func ConversationItemParamOfMcpListTools(serverLabel string, tools []RealtimeMcpListToolsToolParam) ConversationItemUnionParam {
	var mcpListTools RealtimeMcpListToolsParam
	mcpListTools.ServerLabel = serverLabel
	mcpListTools.Tools = tools
	return ConversationItemUnionParam{OfMcpListTools: &mcpListTools}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ConversationItemUnionParam struct {
	OfRealtimeSystemMessageItem    *RealtimeConversationItemSystemMessageParam      `json:",omitzero,inline"`
	OfRealtimeUserMessageItem      *RealtimeConversationItemUserMessageParam        `json:",omitzero,inline"`
	OfRealtimeAssistantMessageItem *RealtimeConversationItemAssistantMessageParam   `json:",omitzero,inline"`
	OfFunctionCall                 *RealtimeConversationItemFunctionCallParam       `json:",omitzero,inline"`
	OfFunctionCallOutput           *RealtimeConversationItemFunctionCallOutputParam `json:",omitzero,inline"`
	OfMcpApprovalResponse          *RealtimeMcpApprovalResponseParam                `json:",omitzero,inline"`
	OfMcpListTools                 *RealtimeMcpListToolsParam                       `json:",omitzero,inline"`
	OfMcpCall                      *RealtimeMcpToolCallParam                        `json:",omitzero,inline"`
	OfMcpApprovalRequest           *RealtimeMcpApprovalRequestParam                 `json:",omitzero,inline"`
	paramUnion
}

func (u ConversationItemUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfRealtimeSystemMessageItem,
		u.OfRealtimeUserMessageItem,
		u.OfRealtimeAssistantMessageItem,
		u.OfFunctionCall,
		u.OfFunctionCallOutput,
		u.OfMcpApprovalResponse,
		u.OfMcpListTools,
		u.OfMcpCall,
		u.OfMcpApprovalRequest)
}
func (u *ConversationItemUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u ConversationItemUnionParam) GetApprove() *bool {
	if vt := u.OfMcpApprovalResponse; vt != nil {
		return &vt.Approve
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ConversationItemUnionParam) GetReason() *string {
	if vt := u.OfMcpApprovalResponse; vt != nil && vt.Reason.Valid() {
		return &vt.Reason.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ConversationItemUnionParam) GetTools() []RealtimeMcpListToolsToolParam {
	if vt := u.OfMcpListTools; vt != nil {
		return vt.Tools
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ConversationItemUnionParam) GetError() *RealtimeMcpToolCallErrorUnionParam {
	if vt := u.OfMcpCall; vt != nil {
		return &vt.Error
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ConversationItemUnionParam) GetRole() *string {
	if vt := u.OfRealtimeSystemMessageItem; vt != nil {
		return (*string)(&vt.Role)
	} else if vt := u.OfRealtimeUserMessageItem; vt != nil {
		return (*string)(&vt.Role)
	} else if vt := u.OfRealtimeAssistantMessageItem; vt != nil {
		return (*string)(&vt.Role)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ConversationItemUnionParam) GetType() *string {
	if vt := u.OfRealtimeSystemMessageItem; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfRealtimeUserMessageItem; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfRealtimeAssistantMessageItem; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfFunctionCall; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfFunctionCallOutput; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfMcpApprovalResponse; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfMcpListTools; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfMcpCall; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfMcpApprovalRequest; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ConversationItemUnionParam) GetID() *string {
	if vt := u.OfRealtimeSystemMessageItem; vt != nil && vt.ID.Valid() {
		return &vt.ID.Value
	} else if vt := u.OfRealtimeUserMessageItem; vt != nil && vt.ID.Valid() {
		return &vt.ID.Value
	} else if vt := u.OfRealtimeAssistantMessageItem; vt != nil && vt.ID.Valid() {
		return &vt.ID.Value
	} else if vt := u.OfFunctionCall; vt != nil && vt.ID.Valid() {
		return &vt.ID.Value
	} else if vt := u.OfFunctionCallOutput; vt != nil && vt.ID.Valid() {
		return &vt.ID.Value
	} else if vt := u.OfMcpApprovalResponse; vt != nil {
		return (*string)(&vt.ID)
	} else if vt := u.OfMcpListTools; vt != nil && vt.ID.Valid() {
		return &vt.ID.Value
	} else if vt := u.OfMcpCall; vt != nil {
		return (*string)(&vt.ID)
	} else if vt := u.OfMcpApprovalRequest; vt != nil {
		return (*string)(&vt.ID)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ConversationItemUnionParam) GetObject() *string {
	if vt := u.OfRealtimeSystemMessageItem; vt != nil {
		return (*string)(&vt.Object)
	} else if vt := u.OfRealtimeUserMessageItem; vt != nil {
		return (*string)(&vt.Object)
	} else if vt := u.OfRealtimeAssistantMessageItem; vt != nil {
		return (*string)(&vt.Object)
	} else if vt := u.OfFunctionCall; vt != nil {
		return (*string)(&vt.Object)
	} else if vt := u.OfFunctionCallOutput; vt != nil {
		return (*string)(&vt.Object)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ConversationItemUnionParam) GetStatus() *string {
	if vt := u.OfRealtimeSystemMessageItem; vt != nil {
		return (*string)(&vt.Status)
	} else if vt := u.OfRealtimeUserMessageItem; vt != nil {
		return (*string)(&vt.Status)
	} else if vt := u.OfRealtimeAssistantMessageItem; vt != nil {
		return (*string)(&vt.Status)
	} else if vt := u.OfFunctionCall; vt != nil {
		return (*string)(&vt.Status)
	} else if vt := u.OfFunctionCallOutput; vt != nil {
		return (*string)(&vt.Status)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ConversationItemUnionParam) GetArguments() *string {
	if vt := u.OfFunctionCall; vt != nil {
		return (*string)(&vt.Arguments)
	} else if vt := u.OfMcpCall; vt != nil {
		return (*string)(&vt.Arguments)
	} else if vt := u.OfMcpApprovalRequest; vt != nil {
		return (*string)(&vt.Arguments)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ConversationItemUnionParam) GetName() *string {
	if vt := u.OfFunctionCall; vt != nil {
		return (*string)(&vt.Name)
	} else if vt := u.OfMcpCall; vt != nil {
		return (*string)(&vt.Name)
	} else if vt := u.OfMcpApprovalRequest; vt != nil {
		return (*string)(&vt.Name)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ConversationItemUnionParam) GetCallID() *string {
	if vt := u.OfFunctionCall; vt != nil && vt.CallID.Valid() {
		return &vt.CallID.Value
	} else if vt := u.OfFunctionCallOutput; vt != nil {
		return (*string)(&vt.CallID)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ConversationItemUnionParam) GetOutput() *string {
	if vt := u.OfFunctionCallOutput; vt != nil {
		return (*string)(&vt.Output)
	} else if vt := u.OfMcpCall; vt != nil && vt.Output.Valid() {
		return &vt.Output.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ConversationItemUnionParam) GetApprovalRequestID() *string {
	if vt := u.OfMcpApprovalResponse; vt != nil {
		return (*string)(&vt.ApprovalRequestID)
	} else if vt := u.OfMcpCall; vt != nil && vt.ApprovalRequestID.Valid() {
		return &vt.ApprovalRequestID.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ConversationItemUnionParam) GetServerLabel() *string {
	if vt := u.OfMcpListTools; vt != nil {
		return (*string)(&vt.ServerLabel)
	} else if vt := u.OfMcpCall; vt != nil {
		return (*string)(&vt.ServerLabel)
	} else if vt := u.OfMcpApprovalRequest; vt != nil {
		return (*string)(&vt.ServerLabel)
	}
	return nil
}

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u ConversationItemUnionParam) GetContent() (res conversationItemUnionParamContent) {
	if vt := u.OfRealtimeSystemMessageItem; vt != nil {
		res.any = &vt.Content
	} else if vt := u.OfRealtimeUserMessageItem; vt != nil {
		res.any = &vt.Content
	} else if vt := u.OfRealtimeAssistantMessageItem; vt != nil {
		res.any = &vt.Content
	}
	return
}

// Can have the runtime types
// [_[]RealtimeConversationItemSystemMessageContentParam],
// [_[]RealtimeConversationItemUserMessageContentParam],
// [\*[]RealtimeConversationItemAssistantMessageContentParam]
type conversationItemUnionParamContent struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *[]realtime.RealtimeConversationItemSystemMessageContentParam:
//	case *[]realtime.RealtimeConversationItemUserMessageContentParam:
//	case *[]realtime.RealtimeConversationItemAssistantMessageContentParam:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u conversationItemUnionParamContent) AsAny() any { return u.any }

func init() {
	apijson.RegisterUnion[ConversationItemUnionParam](
		"type",
		apijson.Discriminator[RealtimeConversationItemSystemMessageParam]("message"),
		apijson.Discriminator[RealtimeConversationItemUserMessageParam]("message"),
		apijson.Discriminator[RealtimeConversationItemAssistantMessageParam]("message"),
		apijson.Discriminator[RealtimeConversationItemFunctionCallParam]("function_call"),
		apijson.Discriminator[RealtimeConversationItemFunctionCallOutputParam]("function_call_output"),
		apijson.Discriminator[RealtimeMcpApprovalResponseParam]("mcp_approval_response"),
		apijson.Discriminator[RealtimeMcpListToolsParam]("mcp_list_tools"),
		apijson.Discriminator[RealtimeMcpToolCallParam]("mcp_call"),
		apijson.Discriminator[RealtimeMcpApprovalRequestParam]("mcp_approval_request"),
	)
}

// Sent by the server when an Item is added to the default Conversation. This can
// happen in several cases:
//
//   - When the client sends a `conversation.item.create` event.
//   - When the input audio buffer is committed. In this case the item will be a user
//     message containing the audio from the buffer.
//   - When the model is generating a Response. In this case the
//     `conversation.item.added` event will be sent when the model starts generating
//     a specific Item, and thus it will not yet have any content (and `status` will
//     be `in_progress`).
//
// The event will include the full content of the Item (except when model is
// generating a Response) except for audio data, which can be retrieved separately
// with a `conversation.item.retrieve` event if necessary.
type ConversationItemAdded struct {
	// The unique ID of the server event.
	EventID string `json:"event_id" api:"required"`
	// A single item within a Realtime conversation.
	Item ConversationItemUnion `json:"item" api:"required"`
	// The event type, must be `conversation.item.added`.
	Type constant.ConversationItemAdded `json:"type" default:"conversation.item.added"`
	// The ID of the item that precedes this one, if any. This is used to maintain
	// ordering when items are inserted.
	PreviousItemID string `json:"previous_item_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventID        respjson.Field
		Item           respjson.Field
		Type           respjson.Field
		PreviousItemID respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ConversationItemAdded) RawJSON() string { return r.JSON.raw }
func (r *ConversationItemAdded) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Add a new Item to the Conversation's context, including messages, function
// calls, and function call responses. This event can be used both to populate a
// "history" of the conversation and to add new items mid-stream, but has the
// current limitation that it cannot populate assistant audio messages.
//
// If successful, the server will emit a `conversation.item.added` event and, when
// the item is finalized, a `conversation.item.done` event. Otherwise, an `error`
// event will be sent.
//
// The properties Item, Type are required.
type ConversationItemCreateEventParam struct {
	// A single item within a Realtime conversation.
	Item ConversationItemUnionParam `json:"item,omitzero" api:"required"`
	// Optional client-generated ID used to identify this event.
	EventID param.Opt[string] `json:"event_id,omitzero"`
	// The ID of the preceding item after which the new item will be inserted. If not
	// set, the new item will be appended to the end of the conversation.
	//
	// If set to `root`, the new item will be added to the beginning of the
	// conversation.
	//
	// If set to an existing ID, it allows an item to be inserted mid-conversation. If
	// the ID cannot be found, an error will be returned and the item will not be
	// added.
	PreviousItemID param.Opt[string] `json:"previous_item_id,omitzero"`
	// The event type, must be `conversation.item.create`.
	//
	// This field can be elided, and will marshal its zero value as
	// "conversation.item.create".
	Type constant.ConversationItemCreate `json:"type" default:"conversation.item.create"`
	paramObj
}

func (r ConversationItemCreateEventParam) MarshalJSON() (data []byte, err error) {
	type shadow ConversationItemCreateEventParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ConversationItemCreateEventParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Returned when a conversation item is created. There are several scenarios that
// produce this event:
//
//   - The server is generating a Response, which if successful will produce either
//     one or two Items, which will be of type `message` (role `assistant`) or type
//     `function_call`.
//   - The input audio buffer has been committed, either by the client or the server
//     (in `server_vad` mode). The server will take the content of the input audio
//     buffer and add it to a new user message Item.
//   - The client has sent a `conversation.item.create` event to add a new Item to
//     the Conversation.
type ConversationItemCreatedEvent struct {
	// The unique ID of the server event.
	EventID string `json:"event_id" api:"required"`
	// A single item within a Realtime conversation.
	Item ConversationItemUnion `json:"item" api:"required"`
	// The event type, must be `conversation.item.created`.
	Type constant.ConversationItemCreated `json:"type" default:"conversation.item.created"`
	// The ID of the preceding item in the Conversation context, allows the client to
	// understand the order of the conversation. Can be `null` if the item has no
	// predecessor.
	PreviousItemID string `json:"previous_item_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventID        respjson.Field
		Item           respjson.Field
		Type           respjson.Field
		PreviousItemID respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ConversationItemCreatedEvent) RawJSON() string { return r.JSON.raw }
func (r *ConversationItemCreatedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Send this event when you want to remove any item from the conversation history.
// The server will respond with a `conversation.item.deleted` event, unless the
// item does not exist in the conversation history, in which case the server will
// respond with an error.
//
// The properties ItemID, Type are required.
type ConversationItemDeleteEventParam struct {
	// The ID of the item to delete.
	ItemID string `json:"item_id" api:"required"`
	// Optional client-generated ID used to identify this event.
	EventID param.Opt[string] `json:"event_id,omitzero"`
	// The event type, must be `conversation.item.delete`.
	//
	// This field can be elided, and will marshal its zero value as
	// "conversation.item.delete".
	Type constant.ConversationItemDelete `json:"type" default:"conversation.item.delete"`
	paramObj
}

func (r ConversationItemDeleteEventParam) MarshalJSON() (data []byte, err error) {
	type shadow ConversationItemDeleteEventParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ConversationItemDeleteEventParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Returned when an item in the conversation is deleted by the client with a
// `conversation.item.delete` event. This event is used to synchronize the server's
// understanding of the conversation history with the client's view.
type ConversationItemDeletedEvent struct {
	// The unique ID of the server event.
	EventID string `json:"event_id" api:"required"`
	// The ID of the item that was deleted.
	ItemID string `json:"item_id" api:"required"`
	// The event type, must be `conversation.item.deleted`.
	Type constant.ConversationItemDeleted `json:"type" default:"conversation.item.deleted"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventID     respjson.Field
		ItemID      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ConversationItemDeletedEvent) RawJSON() string { return r.JSON.raw }
func (r *ConversationItemDeletedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Returned when a conversation item is finalized.
//
// The event will include the full content of the Item except for audio data, which
// can be retrieved separately with a `conversation.item.retrieve` event if needed.
type ConversationItemDone struct {
	// The unique ID of the server event.
	EventID string `json:"event_id" api:"required"`
	// A single item within a Realtime conversation.
	Item ConversationItemUnion `json:"item" api:"required"`
	// The event type, must be `conversation.item.done`.
	Type constant.ConversationItemDone `json:"type" default:"conversation.item.done"`
	// The ID of the item that precedes this one, if any. This is used to maintain
	// ordering when items are inserted.
	PreviousItemID string `json:"previous_item_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventID        respjson.Field
		Item           respjson.Field
		Type           respjson.Field
		PreviousItemID respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ConversationItemDone) RawJSON() string { return r.JSON.raw }
func (r *ConversationItemDone) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// This event is the output of audio transcription for user audio written to the
// user audio buffer. Transcription begins when the input audio buffer is committed
// by the client or server (when VAD is enabled). Transcription runs asynchronously
// with Response creation, so this event may come before or after the Response
// events.
//
// Realtime API models accept audio natively, and thus input transcription is a
// separate process run on a separate ASR (Automatic Speech Recognition) model. The
// transcript may diverge somewhat from the model's interpretation, and should be
// treated as a rough guide.
type ConversationItemInputAudioTranscriptionCompletedEvent struct {
	// The index of the content part containing the audio.
	ContentIndex int64 `json:"content_index" api:"required"`
	// The unique ID of the server event.
	EventID string `json:"event_id" api:"required"`
	// The ID of the item containing the audio that is being transcribed.
	ItemID string `json:"item_id" api:"required"`
	// The transcribed text.
	Transcript string `json:"transcript" api:"required"`
	// The event type, must be `conversation.item.input_audio_transcription.completed`.
	Type constant.ConversationItemInputAudioTranscriptionCompleted `json:"type" default:"conversation.item.input_audio_transcription.completed"`
	// Usage statistics for the transcription, this is billed according to the ASR
	// model's pricing rather than the realtime model's pricing.
	Usage ConversationItemInputAudioTranscriptionCompletedEventUsageUnion `json:"usage" api:"required"`
	// The languages detected in the audio. Returned by `gpt-transcribe`. An empty
	// array indicates that no language could be reliably detected.
	Languages []TranscriptionLanguage `json:"languages"`
	// The log probabilities of the transcription.
	Logprobs []LogProbProperties `json:"logprobs" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ContentIndex respjson.Field
		EventID      respjson.Field
		ItemID       respjson.Field
		Transcript   respjson.Field
		Type         respjson.Field
		Usage        respjson.Field
		Languages    respjson.Field
		Logprobs     respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ConversationItemInputAudioTranscriptionCompletedEvent) RawJSON() string { return r.JSON.raw }
func (r *ConversationItemInputAudioTranscriptionCompletedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ConversationItemInputAudioTranscriptionCompletedEventUsageUnion contains all
// possible properties and values from
// [ConversationItemInputAudioTranscriptionCompletedEventUsageTranscriptTextUsageTokens],
// [ConversationItemInputAudioTranscriptionCompletedEventUsageTranscriptTextUsageDuration].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ConversationItemInputAudioTranscriptionCompletedEventUsageUnion struct {
	OfConversationItemInputAudioTranscriptionCompletedEventUsageTranscriptTextUsageTokens   ConversationItemInputAudioTranscriptionCompletedEventUsageTranscriptTextUsageTokens   `json:",inline"`
	OfConversationItemInputAudioTranscriptionCompletedEventUsageTranscriptTextUsageDuration ConversationItemInputAudioTranscriptionCompletedEventUsageTranscriptTextUsageDuration `json:",inline"`
	JSON                                                                                    struct {
		OfConversationItemInputAudioTranscriptionCompletedEventUsageTranscriptTextUsageTokens   respjson.Field
		OfConversationItemInputAudioTranscriptionCompletedEventUsageTranscriptTextUsageDuration respjson.Field
		raw                                                                                     string
	} `json:"-"`
}

func (u ConversationItemInputAudioTranscriptionCompletedEventUsageUnion) AsConversationItemInputAudioTranscriptionCompletedEventUsageTranscriptTextUsageTokens() (v ConversationItemInputAudioTranscriptionCompletedEventUsageTranscriptTextUsageTokens) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ConversationItemInputAudioTranscriptionCompletedEventUsageUnion) AsConversationItemInputAudioTranscriptionCompletedEventUsageTranscriptTextUsageDuration() (v ConversationItemInputAudioTranscriptionCompletedEventUsageTranscriptTextUsageDuration) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ConversationItemInputAudioTranscriptionCompletedEventUsageUnion) RawJSON() string {
	return u.JSON.raw
}

func (r *ConversationItemInputAudioTranscriptionCompletedEventUsageUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Usage statistics for models billed by token usage.
type ConversationItemInputAudioTranscriptionCompletedEventUsageTranscriptTextUsageTokens struct {
	// Details about the input tokens billed for this request.
	InputTokenDetails ConversationItemInputAudioTranscriptionCompletedEventUsageTranscriptTextUsageTokensInputTokenDetails `json:"input_token_details"`
	// Number of input tokens billed for this request.
	InputTokens int64 `json:"input_tokens" api:"required"`
	// Number of output tokens generated.
	OutputTokens int64 `json:"output_tokens" api:"required"`
	// Total number of tokens used (input + output).
	TotalTokens int64 `json:"total_tokens" api:"required"`
	// The type of the usage object. Always `tokens` for this variant.
	Type constant.Tokens `json:"type" default:"tokens"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		InputTokenDetails respjson.Field
		InputTokens       respjson.Field
		OutputTokens      respjson.Field
		TotalTokens       respjson.Field
		Type              respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ConversationItemInputAudioTranscriptionCompletedEventUsageTranscriptTextUsageTokens) RawJSON() string {
	return r.JSON.raw
}
func (r *ConversationItemInputAudioTranscriptionCompletedEventUsageTranscriptTextUsageTokens) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Details about the input tokens billed for this request.
type ConversationItemInputAudioTranscriptionCompletedEventUsageTranscriptTextUsageTokensInputTokenDetails struct {
	// Number of audio tokens billed for this request.
	AudioTokens int64 `json:"audio_tokens"`
	// Number of text tokens billed for this request.
	TextTokens int64 `json:"text_tokens"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		AudioTokens respjson.Field
		TextTokens  respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ConversationItemInputAudioTranscriptionCompletedEventUsageTranscriptTextUsageTokensInputTokenDetails) RawJSON() string {
	return r.JSON.raw
}
func (r *ConversationItemInputAudioTranscriptionCompletedEventUsageTranscriptTextUsageTokensInputTokenDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Usage statistics for models billed by audio input duration.
type ConversationItemInputAudioTranscriptionCompletedEventUsageTranscriptTextUsageDuration struct {
	// Duration of the input audio in seconds.
	Seconds float64 `json:"seconds" api:"required"`
	// The type of the usage object. Always `duration` for this variant.
	Type constant.Duration `json:"type" default:"duration"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Seconds     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ConversationItemInputAudioTranscriptionCompletedEventUsageTranscriptTextUsageDuration) RawJSON() string {
	return r.JSON.raw
}
func (r *ConversationItemInputAudioTranscriptionCompletedEventUsageTranscriptTextUsageDuration) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Returned when the text value of an input audio transcription content part is
// updated with incremental transcription results.
type ConversationItemInputAudioTranscriptionDeltaEvent struct {
	// The unique ID of the server event.
	EventID string `json:"event_id" api:"required"`
	// The ID of the item containing the audio that is being transcribed.
	ItemID string `json:"item_id" api:"required"`
	// The event type, must be `conversation.item.input_audio_transcription.delta`.
	Type constant.ConversationItemInputAudioTranscriptionDelta `json:"type" default:"conversation.item.input_audio_transcription.delta"`
	// The index of the content part in the item's content array.
	ContentIndex int64 `json:"content_index"`
	// The text delta.
	Delta string `json:"delta"`
	// The log probabilities of the transcription. These can be enabled by configuring
	// the session with `"include": ["item.input_audio_transcription.logprobs"]`. Each
	// entry in the array corresponds a log probability of which token would be
	// selected for this chunk of transcription. This can help to identify if it was
	// possible there were multiple valid options for a given chunk of transcription.
	Logprobs []LogProbProperties `json:"logprobs" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventID      respjson.Field
		ItemID       respjson.Field
		Type         respjson.Field
		ContentIndex respjson.Field
		Delta        respjson.Field
		Logprobs     respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ConversationItemInputAudioTranscriptionDeltaEvent) RawJSON() string { return r.JSON.raw }
func (r *ConversationItemInputAudioTranscriptionDeltaEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Returned when input audio transcription is configured, and a transcription
// request for a user message failed. These events are separate from other `error`
// events so that the client can identify the related Item.
type ConversationItemInputAudioTranscriptionFailedEvent struct {
	// The index of the content part containing the audio.
	ContentIndex int64 `json:"content_index" api:"required"`
	// Details of the transcription error.
	Error ConversationItemInputAudioTranscriptionFailedEventError `json:"error" api:"required"`
	// The unique ID of the server event.
	EventID string `json:"event_id" api:"required"`
	// The ID of the user message item.
	ItemID string `json:"item_id" api:"required"`
	// The event type, must be `conversation.item.input_audio_transcription.failed`.
	Type constant.ConversationItemInputAudioTranscriptionFailed `json:"type" default:"conversation.item.input_audio_transcription.failed"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ContentIndex respjson.Field
		Error        respjson.Field
		EventID      respjson.Field
		ItemID       respjson.Field
		Type         respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ConversationItemInputAudioTranscriptionFailedEvent) RawJSON() string { return r.JSON.raw }
func (r *ConversationItemInputAudioTranscriptionFailedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Details of the transcription error.
type ConversationItemInputAudioTranscriptionFailedEventError struct {
	// Error code, if any.
	Code string `json:"code"`
	// A human-readable error message.
	Message string `json:"message"`
	// Parameter related to the error, if any.
	Param string `json:"param"`
	// The type of error.
	Type string `json:"type"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Code        respjson.Field
		Message     respjson.Field
		Param       respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ConversationItemInputAudioTranscriptionFailedEventError) RawJSON() string { return r.JSON.raw }
func (r *ConversationItemInputAudioTranscriptionFailedEventError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Returned when an input audio transcription segment is identified for an item.
type ConversationItemInputAudioTranscriptionSegment struct {
	// The segment identifier.
	ID string `json:"id" api:"required"`
	// The index of the input audio content part within the item.
	ContentIndex int64 `json:"content_index" api:"required"`
	// End time of the segment in seconds.
	End float64 `json:"end" api:"required"`
	// The unique ID of the server event.
	EventID string `json:"event_id" api:"required"`
	// The ID of the item containing the input audio content.
	ItemID string `json:"item_id" api:"required"`
	// The detected speaker label for this segment.
	Speaker string `json:"speaker" api:"required"`
	// Start time of the segment in seconds.
	Start float64 `json:"start" api:"required"`
	// The text for this segment.
	Text string `json:"text" api:"required"`
	// The event type, must be `conversation.item.input_audio_transcription.segment`.
	Type constant.ConversationItemInputAudioTranscriptionSegment `json:"type" default:"conversation.item.input_audio_transcription.segment"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID           respjson.Field
		ContentIndex respjson.Field
		End          respjson.Field
		EventID      respjson.Field
		ItemID       respjson.Field
		Speaker      respjson.Field
		Start        respjson.Field
		Text         respjson.Field
		Type         respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ConversationItemInputAudioTranscriptionSegment) RawJSON() string { return r.JSON.raw }
func (r *ConversationItemInputAudioTranscriptionSegment) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Send this event when you want to retrieve the server's representation of a
// specific item in the conversation history. This is useful, for example, to
// inspect user audio after noise cancellation and VAD. The server will respond
// with a `conversation.item.retrieved` event, unless the item does not exist in
// the conversation history, in which case the server will respond with an error.
//
// The properties ItemID, Type are required.
type ConversationItemRetrieveEventParam struct {
	// The ID of the item to retrieve.
	ItemID string `json:"item_id" api:"required"`
	// Optional client-generated ID used to identify this event.
	EventID param.Opt[string] `json:"event_id,omitzero"`
	// The event type, must be `conversation.item.retrieve`.
	//
	// This field can be elided, and will marshal its zero value as
	// "conversation.item.retrieve".
	Type constant.ConversationItemRetrieve `json:"type" default:"conversation.item.retrieve"`
	paramObj
}

func (r ConversationItemRetrieveEventParam) MarshalJSON() (data []byte, err error) {
	type shadow ConversationItemRetrieveEventParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ConversationItemRetrieveEventParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Send this event to truncate a previous assistant message’s audio. The server
// will produce audio faster than realtime, so this event is useful when the user
// interrupts to truncate audio that has already been sent to the client but not
// yet played. This will synchronize the server's understanding of the audio with
// the client's playback.
//
// Truncating audio will delete the server-side text transcript to ensure there is
// not text in the context that hasn't been heard by the user.
//
// If successful, the server will respond with a `conversation.item.truncated`
// event.
//
// The properties AudioEndMs, ContentIndex, ItemID, Type are required.
type ConversationItemTruncateEventParam struct {
	// Inclusive duration up to which audio is truncated, in milliseconds. If the
	// audio_end_ms is greater than the actual audio duration, the server will respond
	// with an error.
	AudioEndMs int64 `json:"audio_end_ms" api:"required"`
	// The index of the content part to truncate. Set this to `0`.
	ContentIndex int64 `json:"content_index" api:"required"`
	// The ID of the assistant message item to truncate. Only assistant message items
	// can be truncated.
	ItemID string `json:"item_id" api:"required"`
	// Optional client-generated ID used to identify this event.
	EventID param.Opt[string] `json:"event_id,omitzero"`
	// The event type, must be `conversation.item.truncate`.
	//
	// This field can be elided, and will marshal its zero value as
	// "conversation.item.truncate".
	Type constant.ConversationItemTruncate `json:"type" default:"conversation.item.truncate"`
	paramObj
}

func (r ConversationItemTruncateEventParam) MarshalJSON() (data []byte, err error) {
	type shadow ConversationItemTruncateEventParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ConversationItemTruncateEventParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Returned when an earlier assistant audio message item is truncated by the client
// with a `conversation.item.truncate` event. This event is used to synchronize the
// server's understanding of the audio with the client's playback.
//
// This action will truncate the audio and remove the server-side text transcript
// to ensure there is no text in the context that hasn't been heard by the user.
type ConversationItemTruncatedEvent struct {
	// The duration up to which the audio was truncated, in milliseconds.
	AudioEndMs int64 `json:"audio_end_ms" api:"required"`
	// The index of the content part that was truncated.
	ContentIndex int64 `json:"content_index" api:"required"`
	// The unique ID of the server event.
	EventID string `json:"event_id" api:"required"`
	// The ID of the assistant message item that was truncated.
	ItemID string `json:"item_id" api:"required"`
	// The event type, must be `conversation.item.truncated`.
	Type constant.ConversationItemTruncated `json:"type" default:"conversation.item.truncated"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		AudioEndMs   respjson.Field
		ContentIndex respjson.Field
		EventID      respjson.Field
		ItemID       respjson.Field
		Type         respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ConversationItemTruncatedEvent) RawJSON() string { return r.JSON.raw }
func (r *ConversationItemTruncatedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Send this event to append audio bytes to the input audio buffer. The audio
// buffer is temporary storage you can write to and later commit. A "commit" will
// create a new user message item in the conversation history from the buffer
// content and clear the buffer. Input audio transcription (if enabled) will be
// generated when the buffer is committed.
//
// If VAD is enabled the audio buffer is used to detect speech and the server will
// decide when to commit. When Server VAD is disabled, you must commit the audio
// buffer manually. Input audio noise reduction operates on writes to the audio
// buffer.
//
// The client may choose how much audio to place in each event up to a maximum of
// 15 MiB, for example streaming smaller chunks from the client may allow the VAD
// to be more responsive. Unlike most other client events, the server will not send
// a confirmation response to this event.
//
// The properties Audio, Type are required.
type InputAudioBufferAppendEventParam struct {
	// Base64-encoded audio bytes. This must be in the format specified by the
	// `input_audio_format` field in the session configuration.
	Audio string `json:"audio" api:"required"`
	// Optional client-generated ID used to identify this event.
	EventID param.Opt[string] `json:"event_id,omitzero"`
	// The event type, must be `input_audio_buffer.append`.
	//
	// This field can be elided, and will marshal its zero value as
	// "input_audio_buffer.append".
	Type constant.InputAudioBufferAppend `json:"type" default:"input_audio_buffer.append"`
	paramObj
}

func (r InputAudioBufferAppendEventParam) MarshalJSON() (data []byte, err error) {
	type shadow InputAudioBufferAppendEventParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *InputAudioBufferAppendEventParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Send this event to clear the audio bytes in the buffer. The server will respond
// with an `input_audio_buffer.cleared` event.
//
// The property Type is required.
type InputAudioBufferClearEventParam struct {
	// Optional client-generated ID used to identify this event.
	EventID param.Opt[string] `json:"event_id,omitzero"`
	// The event type, must be `input_audio_buffer.clear`.
	//
	// This field can be elided, and will marshal its zero value as
	// "input_audio_buffer.clear".
	Type constant.InputAudioBufferClear `json:"type" default:"input_audio_buffer.clear"`
	paramObj
}

func (r InputAudioBufferClearEventParam) MarshalJSON() (data []byte, err error) {
	type shadow InputAudioBufferClearEventParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *InputAudioBufferClearEventParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Returned when the input audio buffer is cleared by the client with a
// `input_audio_buffer.clear` event.
type InputAudioBufferClearedEvent struct {
	// The unique ID of the server event.
	EventID string `json:"event_id" api:"required"`
	// The event type, must be `input_audio_buffer.cleared`.
	Type constant.InputAudioBufferCleared `json:"type" default:"input_audio_buffer.cleared"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventID     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r InputAudioBufferClearedEvent) RawJSON() string { return r.JSON.raw }
func (r *InputAudioBufferClearedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Send this event to commit the user input audio buffer, which will create a new
// user message item in the conversation. This event will produce an error if the
// input audio buffer is empty. When in Server VAD mode, the client does not need
// to send this event, the server will commit the audio buffer automatically.
//
// Committing the input audio buffer will trigger input audio transcription (if
// enabled in session configuration), but it will not create a response from the
// model. The server will respond with an `input_audio_buffer.committed` event.
//
// The property Type is required.
type InputAudioBufferCommitEventParam struct {
	// Optional client-generated ID used to identify this event.
	EventID param.Opt[string] `json:"event_id,omitzero"`
	// The event type, must be `input_audio_buffer.commit`.
	//
	// This field can be elided, and will marshal its zero value as
	// "input_audio_buffer.commit".
	Type constant.InputAudioBufferCommit `json:"type" default:"input_audio_buffer.commit"`
	paramObj
}

func (r InputAudioBufferCommitEventParam) MarshalJSON() (data []byte, err error) {
	type shadow InputAudioBufferCommitEventParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *InputAudioBufferCommitEventParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Returned when an input audio buffer is committed, either by the client or
// automatically in server VAD mode. The `item_id` property is the ID of the user
// message item that will be created, thus a `conversation.item.created` event will
// also be sent to the client.
type InputAudioBufferCommittedEvent struct {
	// The unique ID of the server event.
	EventID string `json:"event_id" api:"required"`
	// The ID of the user message item that will be created.
	ItemID string `json:"item_id" api:"required"`
	// The event type, must be `input_audio_buffer.committed`.
	Type constant.InputAudioBufferCommitted `json:"type" default:"input_audio_buffer.committed"`
	// The ID of the preceding item after which the new item will be inserted. Can be
	// `null` if the item has no predecessor.
	PreviousItemID string `json:"previous_item_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventID        respjson.Field
		ItemID         respjson.Field
		Type           respjson.Field
		PreviousItemID respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r InputAudioBufferCommittedEvent) RawJSON() string { return r.JSON.raw }
func (r *InputAudioBufferCommittedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// **SIP Only:** Returned when an DTMF event is received. A DTMF event is a message
// that represents a telephone keypad press (0–9, *, #, A–D). The `event` property
// is the keypad that the user press. The `received_at` is the UTC Unix Timestamp
// that the server received the event.
type InputAudioBufferDtmfEventReceivedEvent struct {
	// The telephone keypad that was pressed by the user.
	Event string `json:"event" api:"required"`
	// UTC Unix Timestamp when DTMF Event was received by server.
	ReceivedAt int64 `json:"received_at" api:"required"`
	// The event type, must be `input_audio_buffer.dtmf_event_received`.
	Type constant.InputAudioBufferDtmfEventReceived `json:"type" default:"input_audio_buffer.dtmf_event_received"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Event       respjson.Field
		ReceivedAt  respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r InputAudioBufferDtmfEventReceivedEvent) RawJSON() string { return r.JSON.raw }
func (r *InputAudioBufferDtmfEventReceivedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Sent by the server when in `server_vad` mode to indicate that speech has been
// detected in the audio buffer. This can happen any time audio is added to the
// buffer (unless speech is already detected). The client may want to use this
// event to interrupt audio playback or provide visual feedback to the user.
//
// The client should expect to receive a `input_audio_buffer.speech_stopped` event
// when speech stops. The `item_id` property is the ID of the user message item
// that will be created when speech stops and will also be included in the
// `input_audio_buffer.speech_stopped` event (unless the client manually commits
// the audio buffer during VAD activation).
type InputAudioBufferSpeechStartedEvent struct {
	// Milliseconds from the start of all audio written to the buffer during the
	// session when speech was first detected. This will correspond to the beginning of
	// audio sent to the model, and thus includes the `prefix_padding_ms` configured in
	// the Session.
	AudioStartMs int64 `json:"audio_start_ms" api:"required"`
	// The unique ID of the server event.
	EventID string `json:"event_id" api:"required"`
	// The ID of the user message item that will be created when speech stops.
	ItemID string `json:"item_id" api:"required"`
	// The event type, must be `input_audio_buffer.speech_started`.
	Type constant.InputAudioBufferSpeechStarted `json:"type" default:"input_audio_buffer.speech_started"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		AudioStartMs respjson.Field
		EventID      respjson.Field
		ItemID       respjson.Field
		Type         respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r InputAudioBufferSpeechStartedEvent) RawJSON() string { return r.JSON.raw }
func (r *InputAudioBufferSpeechStartedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Returned in `server_vad` mode when the server detects the end of speech in the
// audio buffer. The server will also send an `conversation.item.created` event
// with the user message item that is created from the audio buffer.
type InputAudioBufferSpeechStoppedEvent struct {
	// Milliseconds since the session started when speech stopped. This will correspond
	// to the end of audio sent to the model, and thus includes the
	// `min_silence_duration_ms` configured in the Session.
	AudioEndMs int64 `json:"audio_end_ms" api:"required"`
	// The unique ID of the server event.
	EventID string `json:"event_id" api:"required"`
	// The ID of the user message item that will be created.
	ItemID string `json:"item_id" api:"required"`
	// The event type, must be `input_audio_buffer.speech_stopped`.
	Type constant.InputAudioBufferSpeechStopped `json:"type" default:"input_audio_buffer.speech_stopped"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		AudioEndMs  respjson.Field
		EventID     respjson.Field
		ItemID      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r InputAudioBufferSpeechStoppedEvent) RawJSON() string { return r.JSON.raw }
func (r *InputAudioBufferSpeechStoppedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Returned when the Server VAD timeout is triggered for the input audio buffer.
// This is configured with `idle_timeout_ms` in the `turn_detection` settings of
// the session, and it indicates that there hasn't been any speech detected for the
// configured duration.
//
// The `audio_start_ms` and `audio_end_ms` fields indicate the segment of audio
// after the last model response up to the triggering time, as an offset from the
// beginning of audio written to the input audio buffer. This means it demarcates
// the segment of audio that was silent and the difference between the start and
// end values will roughly match the configured timeout.
//
// The empty audio will be committed to the conversation as an `input_audio` item
// (there will be a `input_audio_buffer.committed` event) and a model response will
// be generated. There may be speech that didn't trigger VAD but is still detected
// by the model, so the model may respond with something relevant to the
// conversation or a prompt to continue speaking.
type InputAudioBufferTimeoutTriggered struct {
	// Millisecond offset of audio written to the input audio buffer at the time the
	// timeout was triggered.
	AudioEndMs int64 `json:"audio_end_ms" api:"required"`
	// Millisecond offset of audio written to the input audio buffer that was after the
	// playback time of the last model response.
	AudioStartMs int64 `json:"audio_start_ms" api:"required"`
	// The unique ID of the server event.
	EventID string `json:"event_id" api:"required"`
	// The ID of the item associated with this segment.
	ItemID string `json:"item_id" api:"required"`
	// The event type, must be `input_audio_buffer.timeout_triggered`.
	Type constant.InputAudioBufferTimeoutTriggered `json:"type" default:"input_audio_buffer.timeout_triggered"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		AudioEndMs   respjson.Field
		AudioStartMs respjson.Field
		EventID      respjson.Field
		ItemID       respjson.Field
		Type         respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r InputAudioBufferTimeoutTriggered) RawJSON() string { return r.JSON.raw }
func (r *InputAudioBufferTimeoutTriggered) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A log probability object.
type LogProbProperties struct {
	// The token that was used to generate the log probability.
	Token string `json:"token" api:"required"`
	// The bytes that were used to generate the log probability.
	Bytes []int64 `json:"bytes" api:"required"`
	// The log probability of the token.
	Logprob float64 `json:"logprob" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Token       respjson.Field
		Bytes       respjson.Field
		Logprob     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r LogProbProperties) RawJSON() string { return r.JSON.raw }
func (r *LogProbProperties) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Returned when listing MCP tools has completed for an item.
type McpListToolsCompleted struct {
	// The unique ID of the server event.
	EventID string `json:"event_id" api:"required"`
	// The ID of the MCP list tools item.
	ItemID string `json:"item_id" api:"required"`
	// The event type, must be `mcp_list_tools.completed`.
	Type constant.McpListToolsCompleted `json:"type" default:"mcp_list_tools.completed"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventID     respjson.Field
		ItemID      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r McpListToolsCompleted) RawJSON() string { return r.JSON.raw }
func (r *McpListToolsCompleted) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Returned when listing MCP tools has failed for an item.
type McpListToolsFailed struct {
	// The unique ID of the server event.
	EventID string `json:"event_id" api:"required"`
	// The ID of the MCP list tools item.
	ItemID string `json:"item_id" api:"required"`
	// The event type, must be `mcp_list_tools.failed`.
	Type constant.McpListToolsFailed `json:"type" default:"mcp_list_tools.failed"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventID     respjson.Field
		ItemID      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r McpListToolsFailed) RawJSON() string { return r.JSON.raw }
func (r *McpListToolsFailed) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Returned when listing MCP tools is in progress for an item.
type McpListToolsInProgress struct {
	// The unique ID of the server event.
	EventID string `json:"event_id" api:"required"`
	// The ID of the MCP list tools item.
	ItemID string `json:"item_id" api:"required"`
	// The event type, must be `mcp_list_tools.in_progress`.
	Type constant.McpListToolsInProgress `json:"type" default:"mcp_list_tools.in_progress"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventID     respjson.Field
		ItemID      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r McpListToolsInProgress) RawJSON() string { return r.JSON.raw }
func (r *McpListToolsInProgress) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Type of noise reduction. `near_field` is for close-talking microphones such as
// headphones, `far_field` is for far-field microphones such as laptop or
// conference room microphones.
type NoiseReductionType string

const (
	NoiseReductionTypeNearField NoiseReductionType = "near_field"
	NoiseReductionTypeFarField  NoiseReductionType = "far_field"
)

// **WebRTC/SIP Only:** Emit to cut off the current audio response. This will
// trigger the server to stop generating audio and emit a
// `output_audio_buffer.cleared` event. This event should be preceded by a
// `response.cancel` client event to stop the generation of the current response.
// [Learn more](https://developers.openai.com/api/docs/guides/realtime-conversations#client-and-server-events-for-audio-in-webrtc).
//
// The property Type is required.
type OutputAudioBufferClearEventParam struct {
	// The unique ID of the client event used for error handling.
	EventID param.Opt[string] `json:"event_id,omitzero"`
	// The event type, must be `output_audio_buffer.clear`.
	//
	// This field can be elided, and will marshal its zero value as
	// "output_audio_buffer.clear".
	Type constant.OutputAudioBufferClear `json:"type" default:"output_audio_buffer.clear"`
	paramObj
}

func (r OutputAudioBufferClearEventParam) MarshalJSON() (data []byte, err error) {
	type shadow OutputAudioBufferClearEventParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *OutputAudioBufferClearEventParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Emitted at the beginning of a Response to indicate the updated rate limits. When
// a Response is created some tokens will be "reserved" for the output tokens, the
// rate limits shown here reflect that reservation, which is then adjusted
// accordingly once the Response is completed.
type RateLimitsUpdatedEvent struct {
	// The unique ID of the server event.
	EventID string `json:"event_id" api:"required"`
	// List of rate limit information.
	RateLimits []RateLimitsUpdatedEventRateLimit `json:"rate_limits" api:"required"`
	// The event type, must be `rate_limits.updated`.
	Type constant.RateLimitsUpdated `json:"type" default:"rate_limits.updated"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventID     respjson.Field
		RateLimits  respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RateLimitsUpdatedEvent) RawJSON() string { return r.JSON.raw }
func (r *RateLimitsUpdatedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type RateLimitsUpdatedEventRateLimit struct {
	// The maximum allowed value for the rate limit.
	Limit int64 `json:"limit"`
	// The name of the rate limit (`requests`, `tokens`).
	//
	// Any of "requests", "tokens".
	Name string `json:"name"`
	// The remaining value before the limit is reached.
	Remaining int64 `json:"remaining"`
	// Seconds until the rate limit resets.
	ResetSeconds float64 `json:"reset_seconds"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Limit        respjson.Field
		Name         respjson.Field
		Remaining    respjson.Field
		ResetSeconds respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RateLimitsUpdatedEventRateLimit) RawJSON() string { return r.JSON.raw }
func (r *RateLimitsUpdatedEventRateLimit) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for input and output audio.
type RealtimeAudioConfig struct {
	Input  RealtimeAudioConfigInput  `json:"input"`
	Output RealtimeAudioConfigOutput `json:"output"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Input       respjson.Field
		Output      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeAudioConfig) RawJSON() string { return r.JSON.raw }
func (r *RealtimeAudioConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this RealtimeAudioConfig to a RealtimeAudioConfigParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// RealtimeAudioConfigParam.Overrides()
func (r RealtimeAudioConfig) ToParam() RealtimeAudioConfigParam {
	return param.Override[RealtimeAudioConfigParam](json.RawMessage(r.RawJSON()))
}

// Configuration for input and output audio.
type RealtimeAudioConfigParam struct {
	Input  RealtimeAudioConfigInputParam  `json:"input,omitzero"`
	Output RealtimeAudioConfigOutputParam `json:"output,omitzero"`
	paramObj
}

func (r RealtimeAudioConfigParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeAudioConfigParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeAudioConfigParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type RealtimeAudioConfigInput struct {
	// The format of the input audio.
	Format RealtimeAudioFormatsUnion `json:"format"`
	// Configuration for input audio noise reduction. This can be set to `null` to turn
	// off. Noise reduction filters audio added to the input audio buffer before it is
	// sent to VAD and the model. Filtering the audio can improve VAD and turn
	// detection accuracy (reducing false positives) and model performance by improving
	// perception of the input audio.
	NoiseReduction RealtimeAudioConfigInputNoiseReduction `json:"noise_reduction"`
	// Configuration for input audio transcription, defaults to off and can be set to
	// `null` to turn off once on. Input audio transcription is not native to the
	// model, since the model consumes audio directly. Transcription runs
	// asynchronously through
	// [the /audio/transcriptions endpoint](https://developers.openai.com/api/reference/resources/audio/subresources/transcriptions/methods/create)
	// and should be treated as guidance of input audio content rather than precisely
	// what the model heard. The client can optionally set the language and prompt for
	// transcription, these offer additional guidance to the transcription service.
	Transcription AudioTranscription `json:"transcription"`
	// Configuration for turn detection, ether Server VAD or Semantic VAD. This can be
	// set to `null` to turn off, in which case the client must manually trigger model
	// response.
	//
	// Server VAD means that the model will detect the start and end of speech based on
	// audio volume and respond at the end of user speech.
	//
	// Semantic VAD is more advanced and uses a turn detection model (in conjunction
	// with VAD) to semantically estimate whether the user has finished speaking, then
	// dynamically sets a timeout based on this probability. For example, if user audio
	// trails off with "uhhm", the model will score a low probability of turn end and
	// wait longer for the user to continue speaking. This can be useful for more
	// natural conversations, but may have a higher latency.
	//
	// For `gpt-realtime-whisper` transcription sessions, turn detection must be set to
	// `null`; VAD is not supported.
	TurnDetection RealtimeAudioInputTurnDetectionUnion `json:"turn_detection" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Format         respjson.Field
		NoiseReduction respjson.Field
		Transcription  respjson.Field
		TurnDetection  respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeAudioConfigInput) RawJSON() string { return r.JSON.raw }
func (r *RealtimeAudioConfigInput) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this RealtimeAudioConfigInput to a
// RealtimeAudioConfigInputParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// RealtimeAudioConfigInputParam.Overrides()
func (r RealtimeAudioConfigInput) ToParam() RealtimeAudioConfigInputParam {
	return param.Override[RealtimeAudioConfigInputParam](json.RawMessage(r.RawJSON()))
}

// Configuration for input audio noise reduction. This can be set to `null` to turn
// off. Noise reduction filters audio added to the input audio buffer before it is
// sent to VAD and the model. Filtering the audio can improve VAD and turn
// detection accuracy (reducing false positives) and model performance by improving
// perception of the input audio.
type RealtimeAudioConfigInputNoiseReduction struct {
	// Type of noise reduction. `near_field` is for close-talking microphones such as
	// headphones, `far_field` is for far-field microphones such as laptop or
	// conference room microphones.
	//
	// Any of "near_field", "far_field".
	Type NoiseReductionType `json:"type"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeAudioConfigInputNoiseReduction) RawJSON() string { return r.JSON.raw }
func (r *RealtimeAudioConfigInputNoiseReduction) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type RealtimeAudioConfigInputParam struct {
	// Configuration for turn detection, ether Server VAD or Semantic VAD. This can be
	// set to `null` to turn off, in which case the client must manually trigger model
	// response.
	//
	// Server VAD means that the model will detect the start and end of speech based on
	// audio volume and respond at the end of user speech.
	//
	// Semantic VAD is more advanced and uses a turn detection model (in conjunction
	// with VAD) to semantically estimate whether the user has finished speaking, then
	// dynamically sets a timeout based on this probability. For example, if user audio
	// trails off with "uhhm", the model will score a low probability of turn end and
	// wait longer for the user to continue speaking. This can be useful for more
	// natural conversations, but may have a higher latency.
	//
	// For `gpt-realtime-whisper` transcription sessions, turn detection must be set to
	// `null`; VAD is not supported.
	TurnDetection RealtimeAudioInputTurnDetectionUnionParam `json:"turn_detection,omitzero"`
	// The format of the input audio.
	Format RealtimeAudioFormatsUnionParam `json:"format,omitzero"`
	// Configuration for input audio noise reduction. This can be set to `null` to turn
	// off. Noise reduction filters audio added to the input audio buffer before it is
	// sent to VAD and the model. Filtering the audio can improve VAD and turn
	// detection accuracy (reducing false positives) and model performance by improving
	// perception of the input audio.
	NoiseReduction RealtimeAudioConfigInputNoiseReductionParam `json:"noise_reduction,omitzero"`
	// Configuration for input audio transcription, defaults to off and can be set to
	// `null` to turn off once on. Input audio transcription is not native to the
	// model, since the model consumes audio directly. Transcription runs
	// asynchronously through
	// [the /audio/transcriptions endpoint](https://developers.openai.com/api/reference/resources/audio/subresources/transcriptions/methods/create)
	// and should be treated as guidance of input audio content rather than precisely
	// what the model heard. The client can optionally set the language and prompt for
	// transcription, these offer additional guidance to the transcription service.
	Transcription AudioTranscriptionParam `json:"transcription,omitzero"`
	paramObj
}

func (r RealtimeAudioConfigInputParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeAudioConfigInputParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeAudioConfigInputParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for input audio noise reduction. This can be set to `null` to turn
// off. Noise reduction filters audio added to the input audio buffer before it is
// sent to VAD and the model. Filtering the audio can improve VAD and turn
// detection accuracy (reducing false positives) and model performance by improving
// perception of the input audio.
type RealtimeAudioConfigInputNoiseReductionParam struct {
	// Type of noise reduction. `near_field` is for close-talking microphones such as
	// headphones, `far_field` is for far-field microphones such as laptop or
	// conference room microphones.
	//
	// Any of "near_field", "far_field".
	Type NoiseReductionType `json:"type,omitzero"`
	paramObj
}

func (r RealtimeAudioConfigInputNoiseReductionParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeAudioConfigInputNoiseReductionParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeAudioConfigInputNoiseReductionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type RealtimeAudioConfigOutput struct {
	// The format of the output audio.
	Format RealtimeAudioFormatsUnion `json:"format"`
	// The speed of the model's spoken response as a multiple of the original speed.
	// 1.0 is the default speed. 0.25 is the minimum speed. 1.5 is the maximum speed.
	// This value can only be changed in between model turns, not while a response is
	// in progress.
	//
	// This parameter is a post-processing adjustment to the audio after it is
	// generated, it's also possible to prompt the model to speak faster or slower.
	Speed float64 `json:"speed"`
	// The voice the model uses to respond. Supported built-in voices are `alloy`,
	// `ash`, `ballad`, `coral`, `echo`, `sage`, `shimmer`, `verse`, `marin`, and
	// `cedar`. You may also provide a custom voice object with an `id`, for example
	// `{ "id": "voice_1234" }`. Voice cannot be changed during the session once the
	// model has responded with audio at least once. Custom voices must be created from
	// audio samples. Voices created from text prompts are supported only in Live. We
	// recommend `marin` and `cedar` for best quality.
	Voice RealtimeAudioConfigOutputVoiceUnion `json:"voice"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Format      respjson.Field
		Speed       respjson.Field
		Voice       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeAudioConfigOutput) RawJSON() string { return r.JSON.raw }
func (r *RealtimeAudioConfigOutput) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this RealtimeAudioConfigOutput to a
// RealtimeAudioConfigOutputParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// RealtimeAudioConfigOutputParam.Overrides()
func (r RealtimeAudioConfigOutput) ToParam() RealtimeAudioConfigOutputParam {
	return param.Override[RealtimeAudioConfigOutputParam](json.RawMessage(r.RawJSON()))
}

// RealtimeAudioConfigOutputVoiceUnion contains all possible properties and values
// from [string], [RealtimeAudioConfigOutputVoiceID].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfString]
type RealtimeAudioConfigOutputVoiceUnion struct {
	// This field will be present if the value is a [string] instead of an object.
	OfString                           string                           `json:",inline"`
	OfRealtimeAudioConfigOutputVoiceID RealtimeAudioConfigOutputVoiceID `json:",inline"`
	JSON                               struct {
		OfString                           respjson.Field
		OfRealtimeAudioConfigOutputVoiceID respjson.Field
		raw                                string
	} `json:"-"`
}

func (u RealtimeAudioConfigOutputVoiceUnion) AsString() (v string) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeAudioConfigOutputVoiceUnion) AsRealtimeAudioConfigOutputVoiceID() (v RealtimeAudioConfigOutputVoiceID) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u RealtimeAudioConfigOutputVoiceUnion) RawJSON() string { return u.JSON.raw }

func (r *RealtimeAudioConfigOutputVoiceUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type RealtimeAudioConfigOutputVoiceString2 string

const (
	RealtimeAudioConfigOutputVoiceString2Alloy   RealtimeAudioConfigOutputVoiceString2 = "alloy"
	RealtimeAudioConfigOutputVoiceString2Ash     RealtimeAudioConfigOutputVoiceString2 = "ash"
	RealtimeAudioConfigOutputVoiceString2Ballad  RealtimeAudioConfigOutputVoiceString2 = "ballad"
	RealtimeAudioConfigOutputVoiceString2Coral   RealtimeAudioConfigOutputVoiceString2 = "coral"
	RealtimeAudioConfigOutputVoiceString2Echo    RealtimeAudioConfigOutputVoiceString2 = "echo"
	RealtimeAudioConfigOutputVoiceString2Sage    RealtimeAudioConfigOutputVoiceString2 = "sage"
	RealtimeAudioConfigOutputVoiceString2Shimmer RealtimeAudioConfigOutputVoiceString2 = "shimmer"
	RealtimeAudioConfigOutputVoiceString2Verse   RealtimeAudioConfigOutputVoiceString2 = "verse"
	RealtimeAudioConfigOutputVoiceString2Marin   RealtimeAudioConfigOutputVoiceString2 = "marin"
	RealtimeAudioConfigOutputVoiceString2Cedar   RealtimeAudioConfigOutputVoiceString2 = "cedar"
)

// Custom voice reference.
type RealtimeAudioConfigOutputVoiceID struct {
	// The custom voice ID, e.g. `voice_1234`.
	ID string `json:"id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeAudioConfigOutputVoiceID) RawJSON() string { return r.JSON.raw }
func (r *RealtimeAudioConfigOutputVoiceID) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type RealtimeAudioConfigOutputParam struct {
	// The speed of the model's spoken response as a multiple of the original speed.
	// 1.0 is the default speed. 0.25 is the minimum speed. 1.5 is the maximum speed.
	// This value can only be changed in between model turns, not while a response is
	// in progress.
	//
	// This parameter is a post-processing adjustment to the audio after it is
	// generated, it's also possible to prompt the model to speak faster or slower.
	Speed param.Opt[float64] `json:"speed,omitzero"`
	// The format of the output audio.
	Format RealtimeAudioFormatsUnionParam `json:"format,omitzero"`
	// The voice the model uses to respond. Supported built-in voices are `alloy`,
	// `ash`, `ballad`, `coral`, `echo`, `sage`, `shimmer`, `verse`, `marin`, and
	// `cedar`. You may also provide a custom voice object with an `id`, for example
	// `{ "id": "voice_1234" }`. Voice cannot be changed during the session once the
	// model has responded with audio at least once. Custom voices must be created from
	// audio samples. Voices created from text prompts are supported only in Live. We
	// recommend `marin` and `cedar` for best quality.
	Voice RealtimeAudioConfigOutputVoiceUnionParam `json:"voice,omitzero"`
	paramObj
}

func (r RealtimeAudioConfigOutputParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeAudioConfigOutputParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeAudioConfigOutputParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type RealtimeAudioConfigOutputVoiceUnionParam struct {
	OfString param.Opt[string] `json:",omitzero,inline"`
	// Check if union is this variant with
	// !param.IsOmitted(union.OfRealtimeAudioConfigOutputVoiceString2)
	OfRealtimeAudioConfigOutputVoiceString2 param.Opt[string]                      `json:",omitzero,inline"`
	OfRealtimeAudioConfigOutputVoiceID      *RealtimeAudioConfigOutputVoiceIDParam `json:",omitzero,inline"`
	paramUnion
}

func (u RealtimeAudioConfigOutputVoiceUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfString, u.OfRealtimeAudioConfigOutputVoiceString2, u.OfRealtimeAudioConfigOutputVoiceID)
}
func (u *RealtimeAudioConfigOutputVoiceUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Custom voice reference.
//
// The property ID is required.
type RealtimeAudioConfigOutputVoiceIDParam struct {
	// The custom voice ID, e.g. `voice_1234`.
	ID string `json:"id" api:"required"`
	paramObj
}

func (r RealtimeAudioConfigOutputVoiceIDParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeAudioConfigOutputVoiceIDParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeAudioConfigOutputVoiceIDParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// RealtimeAudioFormatsUnion contains all possible properties and values from
// [RealtimeAudioFormatsAudioPCM], [RealtimeAudioFormatsAudioPCMU],
// [RealtimeAudioFormatsAudioPCMA].
//
// Use the [RealtimeAudioFormatsUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type RealtimeAudioFormatsUnion struct {
	// This field is from variant [RealtimeAudioFormatsAudioPCM].
	Rate int64 `json:"rate"`
	// Any of "audio/pcm", "audio/pcmu", "audio/pcma".
	Type string `json:"type"`
	JSON struct {
		Rate respjson.Field
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyRealtimeAudioFormats is implemented by each variant of
// [RealtimeAudioFormatsUnion] to add type safety for the return type of
// [RealtimeAudioFormatsUnion.AsAny]
type anyRealtimeAudioFormats interface {
	implRealtimeAudioFormatsUnion()
}

func (RealtimeAudioFormatsAudioPCM) implRealtimeAudioFormatsUnion()  {}
func (RealtimeAudioFormatsAudioPCMU) implRealtimeAudioFormatsUnion() {}
func (RealtimeAudioFormatsAudioPCMA) implRealtimeAudioFormatsUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := RealtimeAudioFormatsUnion.AsAny().(type) {
//	case realtime.RealtimeAudioFormatsAudioPCM:
//	case realtime.RealtimeAudioFormatsAudioPCMU:
//	case realtime.RealtimeAudioFormatsAudioPCMA:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u RealtimeAudioFormatsUnion) AsAny() anyRealtimeAudioFormats {
	switch u.Type {
	case "audio/pcm":
		return u.AsAudioPCM()
	case "audio/pcmu":
		return u.AsAudioPCMU()
	case "audio/pcma":
		return u.AsAudioPCMA()
	}
	return nil
}

func (u RealtimeAudioFormatsUnion) AsAudioPCM() (v RealtimeAudioFormatsAudioPCM) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeAudioFormatsUnion) AsAudioPCMU() (v RealtimeAudioFormatsAudioPCMU) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeAudioFormatsUnion) AsAudioPCMA() (v RealtimeAudioFormatsAudioPCMA) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u RealtimeAudioFormatsUnion) RawJSON() string { return u.JSON.raw }

func (r *RealtimeAudioFormatsUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this RealtimeAudioFormatsUnion to a
// RealtimeAudioFormatsUnionParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// RealtimeAudioFormatsUnionParam.Overrides()
func (r RealtimeAudioFormatsUnion) ToParam() RealtimeAudioFormatsUnionParam {
	return param.Override[RealtimeAudioFormatsUnionParam](json.RawMessage(r.RawJSON()))
}

// The PCM audio format. Only a 24kHz sample rate is supported.
type RealtimeAudioFormatsAudioPCM struct {
	// The sample rate of the audio. Always `24000`.
	//
	// Any of 24000.
	Rate int64 `json:"rate"`
	// The audio format. Always `audio/pcm`.
	//
	// Any of "audio/pcm".
	Type string `json:"type"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Rate        respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeAudioFormatsAudioPCM) RawJSON() string { return r.JSON.raw }
func (r *RealtimeAudioFormatsAudioPCM) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The G.711 μ-law format.
type RealtimeAudioFormatsAudioPCMU struct {
	// The audio format. Always `audio/pcmu`.
	//
	// Any of "audio/pcmu".
	Type string `json:"type"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeAudioFormatsAudioPCMU) RawJSON() string { return r.JSON.raw }
func (r *RealtimeAudioFormatsAudioPCMU) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The G.711 A-law format.
type RealtimeAudioFormatsAudioPCMA struct {
	// The audio format. Always `audio/pcma`.
	//
	// Any of "audio/pcma".
	Type string `json:"type"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeAudioFormatsAudioPCMA) RawJSON() string { return r.JSON.raw }
func (r *RealtimeAudioFormatsAudioPCMA) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type RealtimeAudioFormatsUnionParam struct {
	OfAudioPCM  *RealtimeAudioFormatsAudioPCMParam  `json:",omitzero,inline"`
	OfAudioPCMU *RealtimeAudioFormatsAudioPCMUParam `json:",omitzero,inline"`
	OfAudioPCMA *RealtimeAudioFormatsAudioPCMAParam `json:",omitzero,inline"`
	paramUnion
}

func (u RealtimeAudioFormatsUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfAudioPCM, u.OfAudioPCMU, u.OfAudioPCMA)
}
func (u *RealtimeAudioFormatsUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeAudioFormatsUnionParam) GetRate() *int64 {
	if vt := u.OfAudioPCM; vt != nil {
		return &vt.Rate
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeAudioFormatsUnionParam) GetType() *string {
	if vt := u.OfAudioPCM; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfAudioPCMU; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfAudioPCMA; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[RealtimeAudioFormatsUnionParam](
		"type",
		apijson.Discriminator[RealtimeAudioFormatsAudioPCMParam]("audio/pcm"),
		apijson.Discriminator[RealtimeAudioFormatsAudioPCMUParam]("audio/pcmu"),
		apijson.Discriminator[RealtimeAudioFormatsAudioPCMAParam]("audio/pcma"),
	)
}

// The PCM audio format. Only a 24kHz sample rate is supported.
type RealtimeAudioFormatsAudioPCMParam struct {
	// The sample rate of the audio. Always `24000`.
	//
	// Any of 24000.
	Rate int64 `json:"rate,omitzero"`
	// The audio format. Always `audio/pcm`.
	//
	// Any of "audio/pcm".
	Type string `json:"type,omitzero"`
	paramObj
}

func (r RealtimeAudioFormatsAudioPCMParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeAudioFormatsAudioPCMParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeAudioFormatsAudioPCMParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[RealtimeAudioFormatsAudioPCMParam](
		"rate", 24000,
	)
	apijson.RegisterFieldValidator[RealtimeAudioFormatsAudioPCMParam](
		"type", "audio/pcm",
	)
}

// The G.711 μ-law format.
type RealtimeAudioFormatsAudioPCMUParam struct {
	// The audio format. Always `audio/pcmu`.
	//
	// Any of "audio/pcmu".
	Type string `json:"type,omitzero"`
	paramObj
}

func (r RealtimeAudioFormatsAudioPCMUParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeAudioFormatsAudioPCMUParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeAudioFormatsAudioPCMUParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[RealtimeAudioFormatsAudioPCMUParam](
		"type", "audio/pcmu",
	)
}

// The G.711 A-law format.
type RealtimeAudioFormatsAudioPCMAParam struct {
	// The audio format. Always `audio/pcma`.
	//
	// Any of "audio/pcma".
	Type string `json:"type,omitzero"`
	paramObj
}

func (r RealtimeAudioFormatsAudioPCMAParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeAudioFormatsAudioPCMAParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeAudioFormatsAudioPCMAParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[RealtimeAudioFormatsAudioPCMAParam](
		"type", "audio/pcma",
	)
}

// RealtimeAudioInputTurnDetectionUnion contains all possible properties and values
// from [RealtimeAudioInputTurnDetectionServerVad],
// [RealtimeAudioInputTurnDetectionSemanticVad].
//
// Use the [RealtimeAudioInputTurnDetectionUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type RealtimeAudioInputTurnDetectionUnion struct {
	OfRealtimeAudioInputTurnDetectionServerVad   RealtimeAudioInputTurnDetectionServerVad   `json:",inline"`
	OfRealtimeAudioInputTurnDetectionSemanticVad RealtimeAudioInputTurnDetectionSemanticVad `json:",inline"`
	Type                                         string                                     `json:"type"`
	JSON                                         struct {
		OfRealtimeAudioInputTurnDetectionServerVad   respjson.Field
		OfRealtimeAudioInputTurnDetectionSemanticVad respjson.Field
		Type                                         respjson.Field
		raw                                          string
	} `json:"-"`
}

// anyRealtimeAudioInputTurnDetection is implemented by each variant of
// [RealtimeAudioInputTurnDetectionUnion] to add type safety for the return type of
// [RealtimeAudioInputTurnDetectionUnion.AsAny]
type anyRealtimeAudioInputTurnDetection interface {
	implRealtimeAudioInputTurnDetectionUnion()
}

func (RealtimeAudioInputTurnDetectionServerVad) implRealtimeAudioInputTurnDetectionUnion()   {}
func (RealtimeAudioInputTurnDetectionSemanticVad) implRealtimeAudioInputTurnDetectionUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := RealtimeAudioInputTurnDetectionUnion.AsAny().(type) {
//	case realtime.RealtimeAudioInputTurnDetectionServerVad:
//	case realtime.RealtimeAudioInputTurnDetectionSemanticVad:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u RealtimeAudioInputTurnDetectionUnion) AsAny() anyRealtimeAudioInputTurnDetection {
	switch u.Type {
	case "server_vad":
		return u.AsRealtimeAudioInputTurnDetectionServerVad()
	case "semantic_vad":
		return u.AsRealtimeAudioInputTurnDetectionSemanticVad()
	}
	return nil
}

func (u RealtimeAudioInputTurnDetectionUnion) AsRealtimeAudioInputTurnDetectionServerVad() (v RealtimeAudioInputTurnDetectionServerVad) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeAudioInputTurnDetectionUnion) AsRealtimeAudioInputTurnDetectionSemanticVad() (v RealtimeAudioInputTurnDetectionSemanticVad) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u RealtimeAudioInputTurnDetectionUnion) RawJSON() string { return u.JSON.raw }

func (r *RealtimeAudioInputTurnDetectionUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this RealtimeAudioInputTurnDetectionUnion to a
// RealtimeAudioInputTurnDetectionUnionParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// RealtimeAudioInputTurnDetectionUnionParam.Overrides()
func (r RealtimeAudioInputTurnDetectionUnion) ToParam() RealtimeAudioInputTurnDetectionUnionParam {
	return param.Override[RealtimeAudioInputTurnDetectionUnionParam](json.RawMessage(r.RawJSON()))
}

// Server-side voice activity detection (VAD) which flips on when user speech is
// detected and off after a period of silence.
type RealtimeAudioInputTurnDetectionServerVad struct {
	// Whether or not to automatically generate a response when a VAD stop event
	// occurs. If `interrupt_response` is set to `false` this may fail to create a
	// response if the model is already responding.
	//
	// If both `create_response` and `interrupt_response` are set to `false`, the model
	// will never respond automatically but VAD events will still be emitted.
	CreateResponse bool `json:"create_response"`
	// Optional timeout after which a model response will be triggered automatically.
	// This is useful for situations in which a long pause from the user is unexpected,
	// such as a phone call. The model will effectively prompt the user to continue the
	// conversation based on the current context.
	//
	// The timeout value will be applied after the last model response's audio has
	// finished playing, i.e. it's set to the `response.done` time plus audio playback
	// duration.
	//
	// An `input_audio_buffer.timeout_triggered` event (plus events associated with the
	// Response) will be emitted when the timeout is reached. Idle timeout is currently
	// only supported for `server_vad` mode.
	IdleTimeoutMs int64 `json:"idle_timeout_ms" api:"nullable"`
	// Whether or not to automatically interrupt (cancel) any ongoing response with
	// output to the default conversation (i.e. `conversation` of `auto`) when a VAD
	// start event occurs. If `true` then the response will be cancelled, otherwise it
	// will continue until complete.
	//
	// If both `create_response` and `interrupt_response` are set to `false`, the model
	// will never respond automatically but VAD events will still be emitted.
	InterruptResponse bool `json:"interrupt_response"`
	// Used only for `server_vad` mode. Amount of audio to include before the VAD
	// detected speech (in milliseconds). Defaults to 300ms.
	PrefixPaddingMs int64 `json:"prefix_padding_ms"`
	// Used only for `server_vad` mode. Duration of silence to detect speech stop (in
	// milliseconds). Defaults to 500ms. With shorter values the model will respond
	// more quickly, but may jump in on short pauses from the user.
	SilenceDurationMs int64 `json:"silence_duration_ms"`
	// Used only for `server_vad` mode. Activation threshold for VAD (0.0 to 1.0), this
	// defaults to 0.5. A higher threshold will require louder audio to activate the
	// model, and thus might perform better in noisy environments.
	Threshold float64 `json:"threshold"`
	// Type of turn detection, `server_vad` to turn on simple Server VAD.
	Type constant.ServerVad `json:"type" default:"server_vad"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		CreateResponse    respjson.Field
		IdleTimeoutMs     respjson.Field
		InterruptResponse respjson.Field
		PrefixPaddingMs   respjson.Field
		SilenceDurationMs respjson.Field
		Threshold         respjson.Field
		Type              respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeAudioInputTurnDetectionServerVad) RawJSON() string { return r.JSON.raw }
func (r *RealtimeAudioInputTurnDetectionServerVad) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Server-side semantic turn detection which uses a model to determine when the
// user has finished speaking.
type RealtimeAudioInputTurnDetectionSemanticVad struct {
	// Whether or not to automatically generate a response when a VAD stop event
	// occurs.
	CreateResponse bool `json:"create_response"`
	// Used only for `semantic_vad` mode. The eagerness of the model to respond. `low`
	// will wait longer for the user to continue speaking, `high` will respond more
	// quickly. `auto` is the default and is equivalent to `medium`. `low`, `medium`,
	// and `high` have max timeouts of 8s, 4s, and 2s respectively.
	//
	// Any of "low", "medium", "high", "auto".
	Eagerness string `json:"eagerness"`
	// Whether or not to automatically interrupt any ongoing response with output to
	// the default conversation (i.e. `conversation` of `auto`) when a VAD start event
	// occurs.
	InterruptResponse bool `json:"interrupt_response"`
	// Type of turn detection, `semantic_vad` to turn on Semantic VAD.
	Type constant.SemanticVad `json:"type" default:"semantic_vad"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		CreateResponse    respjson.Field
		Eagerness         respjson.Field
		InterruptResponse respjson.Field
		Type              respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeAudioInputTurnDetectionSemanticVad) RawJSON() string { return r.JSON.raw }
func (r *RealtimeAudioInputTurnDetectionSemanticVad) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type RealtimeAudioInputTurnDetectionUnionParam struct {
	OfServerVad   *RealtimeAudioInputTurnDetectionServerVadParam   `json:",omitzero,inline"`
	OfSemanticVad *RealtimeAudioInputTurnDetectionSemanticVadParam `json:",omitzero,inline"`
	paramUnion
}

func (u RealtimeAudioInputTurnDetectionUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfServerVad, u.OfSemanticVad)
}
func (u *RealtimeAudioInputTurnDetectionUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeAudioInputTurnDetectionUnionParam) GetIdleTimeoutMs() *int64 {
	if vt := u.OfServerVad; vt != nil && vt.IdleTimeoutMs.Valid() {
		return &vt.IdleTimeoutMs.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeAudioInputTurnDetectionUnionParam) GetPrefixPaddingMs() *int64 {
	if vt := u.OfServerVad; vt != nil && vt.PrefixPaddingMs.Valid() {
		return &vt.PrefixPaddingMs.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeAudioInputTurnDetectionUnionParam) GetSilenceDurationMs() *int64 {
	if vt := u.OfServerVad; vt != nil && vt.SilenceDurationMs.Valid() {
		return &vt.SilenceDurationMs.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeAudioInputTurnDetectionUnionParam) GetThreshold() *float64 {
	if vt := u.OfServerVad; vt != nil && vt.Threshold.Valid() {
		return &vt.Threshold.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeAudioInputTurnDetectionUnionParam) GetEagerness() *string {
	if vt := u.OfSemanticVad; vt != nil {
		return &vt.Eagerness
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeAudioInputTurnDetectionUnionParam) GetType() *string {
	if vt := u.OfServerVad; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfSemanticVad; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeAudioInputTurnDetectionUnionParam) GetCreateResponse() *bool {
	if vt := u.OfServerVad; vt != nil && vt.CreateResponse.Valid() {
		return &vt.CreateResponse.Value
	} else if vt := u.OfSemanticVad; vt != nil && vt.CreateResponse.Valid() {
		return &vt.CreateResponse.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeAudioInputTurnDetectionUnionParam) GetInterruptResponse() *bool {
	if vt := u.OfServerVad; vt != nil && vt.InterruptResponse.Valid() {
		return &vt.InterruptResponse.Value
	} else if vt := u.OfSemanticVad; vt != nil && vt.InterruptResponse.Valid() {
		return &vt.InterruptResponse.Value
	}
	return nil
}

func init() {
	apijson.RegisterUnion[RealtimeAudioInputTurnDetectionUnionParam](
		"type",
		apijson.Discriminator[RealtimeAudioInputTurnDetectionServerVadParam]("server_vad"),
		apijson.Discriminator[RealtimeAudioInputTurnDetectionSemanticVadParam]("semantic_vad"),
	)
}

// Server-side voice activity detection (VAD) which flips on when user speech is
// detected and off after a period of silence.
//
// The property Type is required.
type RealtimeAudioInputTurnDetectionServerVadParam struct {
	// Optional timeout after which a model response will be triggered automatically.
	// This is useful for situations in which a long pause from the user is unexpected,
	// such as a phone call. The model will effectively prompt the user to continue the
	// conversation based on the current context.
	//
	// The timeout value will be applied after the last model response's audio has
	// finished playing, i.e. it's set to the `response.done` time plus audio playback
	// duration.
	//
	// An `input_audio_buffer.timeout_triggered` event (plus events associated with the
	// Response) will be emitted when the timeout is reached. Idle timeout is currently
	// only supported for `server_vad` mode.
	IdleTimeoutMs param.Opt[int64] `json:"idle_timeout_ms,omitzero"`
	// Whether or not to automatically generate a response when a VAD stop event
	// occurs. If `interrupt_response` is set to `false` this may fail to create a
	// response if the model is already responding.
	//
	// If both `create_response` and `interrupt_response` are set to `false`, the model
	// will never respond automatically but VAD events will still be emitted.
	CreateResponse param.Opt[bool] `json:"create_response,omitzero"`
	// Whether or not to automatically interrupt (cancel) any ongoing response with
	// output to the default conversation (i.e. `conversation` of `auto`) when a VAD
	// start event occurs. If `true` then the response will be cancelled, otherwise it
	// will continue until complete.
	//
	// If both `create_response` and `interrupt_response` are set to `false`, the model
	// will never respond automatically but VAD events will still be emitted.
	InterruptResponse param.Opt[bool] `json:"interrupt_response,omitzero"`
	// Used only for `server_vad` mode. Amount of audio to include before the VAD
	// detected speech (in milliseconds). Defaults to 300ms.
	PrefixPaddingMs param.Opt[int64] `json:"prefix_padding_ms,omitzero"`
	// Used only for `server_vad` mode. Duration of silence to detect speech stop (in
	// milliseconds). Defaults to 500ms. With shorter values the model will respond
	// more quickly, but may jump in on short pauses from the user.
	SilenceDurationMs param.Opt[int64] `json:"silence_duration_ms,omitzero"`
	// Used only for `server_vad` mode. Activation threshold for VAD (0.0 to 1.0), this
	// defaults to 0.5. A higher threshold will require louder audio to activate the
	// model, and thus might perform better in noisy environments.
	Threshold param.Opt[float64] `json:"threshold,omitzero"`
	// Type of turn detection, `server_vad` to turn on simple Server VAD.
	//
	// This field can be elided, and will marshal its zero value as "server_vad".
	Type constant.ServerVad `json:"type" default:"server_vad"`
	paramObj
}

func (r RealtimeAudioInputTurnDetectionServerVadParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeAudioInputTurnDetectionServerVadParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeAudioInputTurnDetectionServerVadParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Server-side semantic turn detection which uses a model to determine when the
// user has finished speaking.
//
// The property Type is required.
type RealtimeAudioInputTurnDetectionSemanticVadParam struct {
	// Whether or not to automatically generate a response when a VAD stop event
	// occurs.
	CreateResponse param.Opt[bool] `json:"create_response,omitzero"`
	// Whether or not to automatically interrupt any ongoing response with output to
	// the default conversation (i.e. `conversation` of `auto`) when a VAD start event
	// occurs.
	InterruptResponse param.Opt[bool] `json:"interrupt_response,omitzero"`
	// Used only for `semantic_vad` mode. The eagerness of the model to respond. `low`
	// will wait longer for the user to continue speaking, `high` will respond more
	// quickly. `auto` is the default and is equivalent to `medium`. `low`, `medium`,
	// and `high` have max timeouts of 8s, 4s, and 2s respectively.
	//
	// Any of "low", "medium", "high", "auto".
	Eagerness string `json:"eagerness,omitzero"`
	// Type of turn detection, `semantic_vad` to turn on Semantic VAD.
	//
	// This field can be elided, and will marshal its zero value as "semantic_vad".
	Type constant.SemanticVad `json:"type" default:"semantic_vad"`
	paramObj
}

func (r RealtimeAudioInputTurnDetectionSemanticVadParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeAudioInputTurnDetectionSemanticVadParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeAudioInputTurnDetectionSemanticVadParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[RealtimeAudioInputTurnDetectionSemanticVadParam](
		"eagerness", "low", "medium", "high", "auto",
	)
}

func RealtimeClientEventParamOfConversationItemCreate[
	T RealtimeConversationItemSystemMessageParam | RealtimeConversationItemUserMessageParam | RealtimeConversationItemAssistantMessageParam | RealtimeConversationItemFunctionCallParam | RealtimeConversationItemFunctionCallOutputParam | RealtimeMcpApprovalResponseParam | RealtimeMcpListToolsParam | RealtimeMcpToolCallParam | RealtimeMcpApprovalRequestParam,
](item T) RealtimeClientEventUnionParam {
	var conversationItemCreate ConversationItemCreateEventParam
	switch v := any(item).(type) {
	case RealtimeConversationItemSystemMessageParam:
		conversationItemCreate.Item.OfRealtimeSystemMessageItem = &v
	case RealtimeConversationItemUserMessageParam:
		conversationItemCreate.Item.OfRealtimeUserMessageItem = &v
	case RealtimeConversationItemAssistantMessageParam:
		conversationItemCreate.Item.OfRealtimeAssistantMessageItem = &v
	case RealtimeConversationItemFunctionCallParam:
		conversationItemCreate.Item.OfFunctionCall = &v
	case RealtimeConversationItemFunctionCallOutputParam:
		conversationItemCreate.Item.OfFunctionCallOutput = &v
	case RealtimeMcpApprovalResponseParam:
		conversationItemCreate.Item.OfMcpApprovalResponse = &v
	case RealtimeMcpListToolsParam:
		conversationItemCreate.Item.OfMcpListTools = &v
	case RealtimeMcpToolCallParam:
		conversationItemCreate.Item.OfMcpCall = &v
	case RealtimeMcpApprovalRequestParam:
		conversationItemCreate.Item.OfMcpApprovalRequest = &v
	}
	return RealtimeClientEventUnionParam{OfConversationItemCreate: &conversationItemCreate}
}

func RealtimeClientEventParamOfConversationItemDelete(itemID string) RealtimeClientEventUnionParam {
	var conversationItemDelete ConversationItemDeleteEventParam
	conversationItemDelete.ItemID = itemID
	return RealtimeClientEventUnionParam{OfConversationItemDelete: &conversationItemDelete}
}

func RealtimeClientEventParamOfConversationItemRetrieve(itemID string) RealtimeClientEventUnionParam {
	var conversationItemRetrieve ConversationItemRetrieveEventParam
	conversationItemRetrieve.ItemID = itemID
	return RealtimeClientEventUnionParam{OfConversationItemRetrieve: &conversationItemRetrieve}
}

func RealtimeClientEventParamOfConversationItemTruncate(audioEndMs int64, contentIndex int64, itemID string) RealtimeClientEventUnionParam {
	var conversationItemTruncate ConversationItemTruncateEventParam
	conversationItemTruncate.AudioEndMs = audioEndMs
	conversationItemTruncate.ContentIndex = contentIndex
	conversationItemTruncate.ItemID = itemID
	return RealtimeClientEventUnionParam{OfConversationItemTruncate: &conversationItemTruncate}
}

func RealtimeClientEventParamOfInputAudioBufferAppend(audio string) RealtimeClientEventUnionParam {
	var inputAudioBufferAppend InputAudioBufferAppendEventParam
	inputAudioBufferAppend.Audio = audio
	return RealtimeClientEventUnionParam{OfInputAudioBufferAppend: &inputAudioBufferAppend}
}

func RealtimeClientEventParamOfSessionUpdate[
	T RealtimeSessionCreateRequestParam | RealtimeTranscriptionSessionCreateRequestParam,
](session T) RealtimeClientEventUnionParam {
	var sessionUpdate SessionUpdateEventParam
	switch v := any(session).(type) {
	case RealtimeSessionCreateRequestParam:
		sessionUpdate.Session.OfRealtime = &v
	case RealtimeTranscriptionSessionCreateRequestParam:
		sessionUpdate.Session.OfTranscription = &v
	}
	return RealtimeClientEventUnionParam{OfSessionUpdate: &sessionUpdate}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type RealtimeClientEventUnionParam struct {
	OfConversationItemCreate   *ConversationItemCreateEventParam   `json:",omitzero,inline"`
	OfConversationItemDelete   *ConversationItemDeleteEventParam   `json:",omitzero,inline"`
	OfConversationItemRetrieve *ConversationItemRetrieveEventParam `json:",omitzero,inline"`
	OfConversationItemTruncate *ConversationItemTruncateEventParam `json:",omitzero,inline"`
	OfInputAudioBufferAppend   *InputAudioBufferAppendEventParam   `json:",omitzero,inline"`
	OfInputAudioBufferClear    *InputAudioBufferClearEventParam    `json:",omitzero,inline"`
	OfOutputAudioBufferClear   *OutputAudioBufferClearEventParam   `json:",omitzero,inline"`
	OfInputAudioBufferCommit   *InputAudioBufferCommitEventParam   `json:",omitzero,inline"`
	OfResponseCancel           *ResponseCancelEventParam           `json:",omitzero,inline"`
	OfResponseCreate           *ResponseCreateEventParam           `json:",omitzero,inline"`
	OfSessionUpdate            *SessionUpdateEventParam            `json:",omitzero,inline"`
	paramUnion
}

func (u RealtimeClientEventUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfConversationItemCreate,
		u.OfConversationItemDelete,
		u.OfConversationItemRetrieve,
		u.OfConversationItemTruncate,
		u.OfInputAudioBufferAppend,
		u.OfInputAudioBufferClear,
		u.OfOutputAudioBufferClear,
		u.OfInputAudioBufferCommit,
		u.OfResponseCancel,
		u.OfResponseCreate,
		u.OfSessionUpdate)
}
func (u *RealtimeClientEventUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeClientEventUnionParam) GetItem() *ConversationItemUnionParam {
	if vt := u.OfConversationItemCreate; vt != nil {
		return &vt.Item
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeClientEventUnionParam) GetPreviousItemID() *string {
	if vt := u.OfConversationItemCreate; vt != nil && vt.PreviousItemID.Valid() {
		return &vt.PreviousItemID.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeClientEventUnionParam) GetAudioEndMs() *int64 {
	if vt := u.OfConversationItemTruncate; vt != nil {
		return &vt.AudioEndMs
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeClientEventUnionParam) GetContentIndex() *int64 {
	if vt := u.OfConversationItemTruncate; vt != nil {
		return &vt.ContentIndex
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeClientEventUnionParam) GetAudio() *string {
	if vt := u.OfInputAudioBufferAppend; vt != nil {
		return &vt.Audio
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeClientEventUnionParam) GetResponseID() *string {
	if vt := u.OfResponseCancel; vt != nil && vt.ResponseID.Valid() {
		return &vt.ResponseID.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeClientEventUnionParam) GetResponse() *RealtimeResponseCreateParams {
	if vt := u.OfResponseCreate; vt != nil {
		return &vt.Response
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeClientEventUnionParam) GetSession() *SessionUpdateEventSessionUnionParam {
	if vt := u.OfSessionUpdate; vt != nil {
		return &vt.Session
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeClientEventUnionParam) GetType() *string {
	if vt := u.OfConversationItemCreate; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfConversationItemDelete; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfConversationItemRetrieve; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfConversationItemTruncate; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfInputAudioBufferAppend; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfInputAudioBufferClear; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfOutputAudioBufferClear; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfInputAudioBufferCommit; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfResponseCancel; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfResponseCreate; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfSessionUpdate; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeClientEventUnionParam) GetEventID() *string {
	if vt := u.OfConversationItemCreate; vt != nil && vt.EventID.Valid() {
		return &vt.EventID.Value
	} else if vt := u.OfConversationItemDelete; vt != nil && vt.EventID.Valid() {
		return &vt.EventID.Value
	} else if vt := u.OfConversationItemRetrieve; vt != nil && vt.EventID.Valid() {
		return &vt.EventID.Value
	} else if vt := u.OfConversationItemTruncate; vt != nil && vt.EventID.Valid() {
		return &vt.EventID.Value
	} else if vt := u.OfInputAudioBufferAppend; vt != nil && vt.EventID.Valid() {
		return &vt.EventID.Value
	} else if vt := u.OfInputAudioBufferClear; vt != nil && vt.EventID.Valid() {
		return &vt.EventID.Value
	} else if vt := u.OfOutputAudioBufferClear; vt != nil && vt.EventID.Valid() {
		return &vt.EventID.Value
	} else if vt := u.OfInputAudioBufferCommit; vt != nil && vt.EventID.Valid() {
		return &vt.EventID.Value
	} else if vt := u.OfResponseCancel; vt != nil && vt.EventID.Valid() {
		return &vt.EventID.Value
	} else if vt := u.OfResponseCreate; vt != nil && vt.EventID.Valid() {
		return &vt.EventID.Value
	} else if vt := u.OfSessionUpdate; vt != nil && vt.EventID.Valid() {
		return &vt.EventID.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeClientEventUnionParam) GetItemID() *string {
	if vt := u.OfConversationItemDelete; vt != nil {
		return (*string)(&vt.ItemID)
	} else if vt := u.OfConversationItemRetrieve; vt != nil {
		return (*string)(&vt.ItemID)
	} else if vt := u.OfConversationItemTruncate; vt != nil {
		return (*string)(&vt.ItemID)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[RealtimeClientEventUnionParam](
		"type",
		apijson.Discriminator[ConversationItemCreateEventParam]("conversation.item.create"),
		apijson.Discriminator[ConversationItemDeleteEventParam]("conversation.item.delete"),
		apijson.Discriminator[ConversationItemRetrieveEventParam]("conversation.item.retrieve"),
		apijson.Discriminator[ConversationItemTruncateEventParam]("conversation.item.truncate"),
		apijson.Discriminator[InputAudioBufferAppendEventParam]("input_audio_buffer.append"),
		apijson.Discriminator[InputAudioBufferClearEventParam]("input_audio_buffer.clear"),
		apijson.Discriminator[OutputAudioBufferClearEventParam]("output_audio_buffer.clear"),
		apijson.Discriminator[InputAudioBufferCommitEventParam]("input_audio_buffer.commit"),
		apijson.Discriminator[ResponseCancelEventParam]("response.cancel"),
		apijson.Discriminator[ResponseCreateEventParam]("response.create"),
		apijson.Discriminator[SessionUpdateEventParam]("session.update"),
	)
}

// An assistant message item in a Realtime conversation.
type RealtimeConversationItemAssistantMessage struct {
	// The content of the message.
	Content []RealtimeConversationItemAssistantMessageContent `json:"content" api:"required"`
	// The role of the message sender. Always `assistant`.
	Role constant.Assistant `json:"role" default:"assistant"`
	// The type of the item. Always `message`.
	Type constant.Message `json:"type" default:"message"`
	// The unique ID of the item. This may be provided by the client or generated by
	// the server.
	ID string `json:"id"`
	// Identifier for the API object being returned - always `realtime.item`. Optional
	// when creating a new item.
	//
	// Any of "realtime.item".
	Object RealtimeConversationItemAssistantMessageObject `json:"object"`
	// The status of the item. Has no effect on the conversation.
	//
	// Any of "completed", "incomplete", "in_progress".
	Status RealtimeConversationItemAssistantMessageStatus `json:"status"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Content     respjson.Field
		Role        respjson.Field
		Type        respjson.Field
		ID          respjson.Field
		Object      respjson.Field
		Status      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeConversationItemAssistantMessage) RawJSON() string { return r.JSON.raw }
func (r *RealtimeConversationItemAssistantMessage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this RealtimeConversationItemAssistantMessage to a
// RealtimeConversationItemAssistantMessageParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// RealtimeConversationItemAssistantMessageParam.Overrides()
func (r RealtimeConversationItemAssistantMessage) ToParam() RealtimeConversationItemAssistantMessageParam {
	return param.Override[RealtimeConversationItemAssistantMessageParam](json.RawMessage(r.RawJSON()))
}

type RealtimeConversationItemAssistantMessageContent struct {
	// Base64-encoded audio bytes, these will be parsed as the format specified in the
	// session output audio type configuration. This defaults to PCM 16-bit 24kHz mono
	// if not specified.
	Audio string `json:"audio"`
	// The text content.
	Text string `json:"text"`
	// The transcript of the audio content, this will always be present if the output
	// type is `audio`.
	Transcript string `json:"transcript"`
	// The content type, `output_text` or `output_audio` depending on the session
	// `output_modalities` configuration.
	//
	// Any of "output_text", "output_audio".
	Type string `json:"type"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Audio       respjson.Field
		Text        respjson.Field
		Transcript  respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeConversationItemAssistantMessageContent) RawJSON() string { return r.JSON.raw }
func (r *RealtimeConversationItemAssistantMessageContent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Identifier for the API object being returned - always `realtime.item`. Optional
// when creating a new item.
type RealtimeConversationItemAssistantMessageObject string

const (
	RealtimeConversationItemAssistantMessageObjectRealtimeItem RealtimeConversationItemAssistantMessageObject = "realtime.item"
)

// The status of the item. Has no effect on the conversation.
type RealtimeConversationItemAssistantMessageStatus string

const (
	RealtimeConversationItemAssistantMessageStatusCompleted  RealtimeConversationItemAssistantMessageStatus = "completed"
	RealtimeConversationItemAssistantMessageStatusIncomplete RealtimeConversationItemAssistantMessageStatus = "incomplete"
	RealtimeConversationItemAssistantMessageStatusInProgress RealtimeConversationItemAssistantMessageStatus = "in_progress"
)

// An assistant message item in a Realtime conversation.
//
// The properties Content, Role, Type are required.
type RealtimeConversationItemAssistantMessageParam struct {
	// The content of the message.
	Content []RealtimeConversationItemAssistantMessageContentParam `json:"content,omitzero" api:"required"`
	// The unique ID of the item. This may be provided by the client or generated by
	// the server.
	ID param.Opt[string] `json:"id,omitzero"`
	// Identifier for the API object being returned - always `realtime.item`. Optional
	// when creating a new item.
	//
	// Any of "realtime.item".
	Object RealtimeConversationItemAssistantMessageObject `json:"object,omitzero"`
	// The status of the item. Has no effect on the conversation.
	//
	// Any of "completed", "incomplete", "in_progress".
	Status RealtimeConversationItemAssistantMessageStatus `json:"status,omitzero"`
	// The role of the message sender. Always `assistant`.
	//
	// This field can be elided, and will marshal its zero value as "assistant".
	Role constant.Assistant `json:"role" default:"assistant"`
	// The type of the item. Always `message`.
	//
	// This field can be elided, and will marshal its zero value as "message".
	Type constant.Message `json:"type" default:"message"`
	paramObj
}

func (r RealtimeConversationItemAssistantMessageParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeConversationItemAssistantMessageParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeConversationItemAssistantMessageParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type RealtimeConversationItemAssistantMessageContentParam struct {
	// Base64-encoded audio bytes, these will be parsed as the format specified in the
	// session output audio type configuration. This defaults to PCM 16-bit 24kHz mono
	// if not specified.
	Audio param.Opt[string] `json:"audio,omitzero"`
	// The text content.
	Text param.Opt[string] `json:"text,omitzero"`
	// The transcript of the audio content, this will always be present if the output
	// type is `audio`.
	Transcript param.Opt[string] `json:"transcript,omitzero"`
	// The content type, `output_text` or `output_audio` depending on the session
	// `output_modalities` configuration.
	//
	// Any of "output_text", "output_audio".
	Type string `json:"type,omitzero"`
	paramObj
}

func (r RealtimeConversationItemAssistantMessageContentParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeConversationItemAssistantMessageContentParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeConversationItemAssistantMessageContentParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[RealtimeConversationItemAssistantMessageContentParam](
		"type", "output_text", "output_audio",
	)
}

// A function call item in a Realtime conversation.
type RealtimeConversationItemFunctionCall struct {
	// The arguments of the function call. This is a JSON-encoded string representing
	// the arguments passed to the function, for example
	// `{"arg1": "value1", "arg2": 42}`.
	Arguments string `json:"arguments" api:"required"`
	// The name of the function being called.
	Name string `json:"name" api:"required"`
	// The type of the item. Always `function_call`.
	Type constant.FunctionCall `json:"type" default:"function_call"`
	// The unique ID of the item. This may be provided by the client or generated by
	// the server.
	ID string `json:"id"`
	// The ID of the function call.
	CallID string `json:"call_id"`
	// Identifier for the API object being returned - always `realtime.item`. Optional
	// when creating a new item.
	//
	// Any of "realtime.item".
	Object RealtimeConversationItemFunctionCallObject `json:"object"`
	// The status of the item. Has no effect on the conversation.
	//
	// Any of "completed", "incomplete", "in_progress".
	Status RealtimeConversationItemFunctionCallStatus `json:"status"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Arguments   respjson.Field
		Name        respjson.Field
		Type        respjson.Field
		ID          respjson.Field
		CallID      respjson.Field
		Object      respjson.Field
		Status      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeConversationItemFunctionCall) RawJSON() string { return r.JSON.raw }
func (r *RealtimeConversationItemFunctionCall) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this RealtimeConversationItemFunctionCall to a
// RealtimeConversationItemFunctionCallParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// RealtimeConversationItemFunctionCallParam.Overrides()
func (r RealtimeConversationItemFunctionCall) ToParam() RealtimeConversationItemFunctionCallParam {
	return param.Override[RealtimeConversationItemFunctionCallParam](json.RawMessage(r.RawJSON()))
}

// Identifier for the API object being returned - always `realtime.item`. Optional
// when creating a new item.
type RealtimeConversationItemFunctionCallObject string

const (
	RealtimeConversationItemFunctionCallObjectRealtimeItem RealtimeConversationItemFunctionCallObject = "realtime.item"
)

// The status of the item. Has no effect on the conversation.
type RealtimeConversationItemFunctionCallStatus string

const (
	RealtimeConversationItemFunctionCallStatusCompleted  RealtimeConversationItemFunctionCallStatus = "completed"
	RealtimeConversationItemFunctionCallStatusIncomplete RealtimeConversationItemFunctionCallStatus = "incomplete"
	RealtimeConversationItemFunctionCallStatusInProgress RealtimeConversationItemFunctionCallStatus = "in_progress"
)

// A function call item in a Realtime conversation.
//
// The properties Arguments, Name, Type are required.
type RealtimeConversationItemFunctionCallParam struct {
	// The arguments of the function call. This is a JSON-encoded string representing
	// the arguments passed to the function, for example
	// `{"arg1": "value1", "arg2": 42}`.
	Arguments string `json:"arguments" api:"required"`
	// The name of the function being called.
	Name string `json:"name" api:"required"`
	// The unique ID of the item. This may be provided by the client or generated by
	// the server.
	ID param.Opt[string] `json:"id,omitzero"`
	// The ID of the function call.
	CallID param.Opt[string] `json:"call_id,omitzero"`
	// Identifier for the API object being returned - always `realtime.item`. Optional
	// when creating a new item.
	//
	// Any of "realtime.item".
	Object RealtimeConversationItemFunctionCallObject `json:"object,omitzero"`
	// The status of the item. Has no effect on the conversation.
	//
	// Any of "completed", "incomplete", "in_progress".
	Status RealtimeConversationItemFunctionCallStatus `json:"status,omitzero"`
	// The type of the item. Always `function_call`.
	//
	// This field can be elided, and will marshal its zero value as "function_call".
	Type constant.FunctionCall `json:"type" default:"function_call"`
	paramObj
}

func (r RealtimeConversationItemFunctionCallParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeConversationItemFunctionCallParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeConversationItemFunctionCallParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A function call output item in a Realtime conversation.
type RealtimeConversationItemFunctionCallOutput struct {
	// The ID of the function call this output is for.
	CallID string `json:"call_id" api:"required"`
	// The output of the function call, this is free text and can contain any
	// information or simply be empty.
	Output string `json:"output" api:"required"`
	// The type of the item. Always `function_call_output`.
	Type constant.FunctionCallOutput `json:"type" default:"function_call_output"`
	// The unique ID of the item. This may be provided by the client or generated by
	// the server.
	ID string `json:"id"`
	// Identifier for the API object being returned - always `realtime.item`. Optional
	// when creating a new item.
	//
	// Any of "realtime.item".
	Object RealtimeConversationItemFunctionCallOutputObject `json:"object"`
	// The status of the item. Has no effect on the conversation.
	//
	// Any of "completed", "incomplete", "in_progress".
	Status RealtimeConversationItemFunctionCallOutputStatus `json:"status"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		CallID      respjson.Field
		Output      respjson.Field
		Type        respjson.Field
		ID          respjson.Field
		Object      respjson.Field
		Status      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeConversationItemFunctionCallOutput) RawJSON() string { return r.JSON.raw }
func (r *RealtimeConversationItemFunctionCallOutput) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this RealtimeConversationItemFunctionCallOutput to a
// RealtimeConversationItemFunctionCallOutputParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// RealtimeConversationItemFunctionCallOutputParam.Overrides()
func (r RealtimeConversationItemFunctionCallOutput) ToParam() RealtimeConversationItemFunctionCallOutputParam {
	return param.Override[RealtimeConversationItemFunctionCallOutputParam](json.RawMessage(r.RawJSON()))
}

// Identifier for the API object being returned - always `realtime.item`. Optional
// when creating a new item.
type RealtimeConversationItemFunctionCallOutputObject string

const (
	RealtimeConversationItemFunctionCallOutputObjectRealtimeItem RealtimeConversationItemFunctionCallOutputObject = "realtime.item"
)

// The status of the item. Has no effect on the conversation.
type RealtimeConversationItemFunctionCallOutputStatus string

const (
	RealtimeConversationItemFunctionCallOutputStatusCompleted  RealtimeConversationItemFunctionCallOutputStatus = "completed"
	RealtimeConversationItemFunctionCallOutputStatusIncomplete RealtimeConversationItemFunctionCallOutputStatus = "incomplete"
	RealtimeConversationItemFunctionCallOutputStatusInProgress RealtimeConversationItemFunctionCallOutputStatus = "in_progress"
)

// A function call output item in a Realtime conversation.
//
// The properties CallID, Output, Type are required.
type RealtimeConversationItemFunctionCallOutputParam struct {
	// The ID of the function call this output is for.
	CallID string `json:"call_id" api:"required"`
	// The output of the function call, this is free text and can contain any
	// information or simply be empty.
	Output string `json:"output" api:"required"`
	// The unique ID of the item. This may be provided by the client or generated by
	// the server.
	ID param.Opt[string] `json:"id,omitzero"`
	// Identifier for the API object being returned - always `realtime.item`. Optional
	// when creating a new item.
	//
	// Any of "realtime.item".
	Object RealtimeConversationItemFunctionCallOutputObject `json:"object,omitzero"`
	// The status of the item. Has no effect on the conversation.
	//
	// Any of "completed", "incomplete", "in_progress".
	Status RealtimeConversationItemFunctionCallOutputStatus `json:"status,omitzero"`
	// The type of the item. Always `function_call_output`.
	//
	// This field can be elided, and will marshal its zero value as
	// "function_call_output".
	Type constant.FunctionCallOutput `json:"type" default:"function_call_output"`
	paramObj
}

func (r RealtimeConversationItemFunctionCallOutputParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeConversationItemFunctionCallOutputParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeConversationItemFunctionCallOutputParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A system message in a Realtime conversation can be used to provide additional
// context or instructions to the model. This is similar but distinct from the
// instruction prompt provided at the start of a conversation, as system messages
// can be added at any point in the conversation. For major changes to the
// conversation's behavior, use instructions, but for smaller updates (e.g. "the
// user is now asking about a different topic"), use system messages.
type RealtimeConversationItemSystemMessage struct {
	// The content of the message.
	Content []RealtimeConversationItemSystemMessageContent `json:"content" api:"required"`
	// The role of the message sender. Always `system`.
	Role constant.System `json:"role" default:"system"`
	// The type of the item. Always `message`.
	Type constant.Message `json:"type" default:"message"`
	// The unique ID of the item. This may be provided by the client or generated by
	// the server.
	ID string `json:"id"`
	// Identifier for the API object being returned - always `realtime.item`. Optional
	// when creating a new item.
	//
	// Any of "realtime.item".
	Object RealtimeConversationItemSystemMessageObject `json:"object"`
	// The status of the item. Has no effect on the conversation.
	//
	// Any of "completed", "incomplete", "in_progress".
	Status RealtimeConversationItemSystemMessageStatus `json:"status"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Content     respjson.Field
		Role        respjson.Field
		Type        respjson.Field
		ID          respjson.Field
		Object      respjson.Field
		Status      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeConversationItemSystemMessage) RawJSON() string { return r.JSON.raw }
func (r *RealtimeConversationItemSystemMessage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this RealtimeConversationItemSystemMessage to a
// RealtimeConversationItemSystemMessageParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// RealtimeConversationItemSystemMessageParam.Overrides()
func (r RealtimeConversationItemSystemMessage) ToParam() RealtimeConversationItemSystemMessageParam {
	return param.Override[RealtimeConversationItemSystemMessageParam](json.RawMessage(r.RawJSON()))
}

type RealtimeConversationItemSystemMessageContent struct {
	// The text content.
	Text string `json:"text"`
	// The content type. Always `input_text` for system messages.
	//
	// Any of "input_text".
	Type string `json:"type"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Text        respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeConversationItemSystemMessageContent) RawJSON() string { return r.JSON.raw }
func (r *RealtimeConversationItemSystemMessageContent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Identifier for the API object being returned - always `realtime.item`. Optional
// when creating a new item.
type RealtimeConversationItemSystemMessageObject string

const (
	RealtimeConversationItemSystemMessageObjectRealtimeItem RealtimeConversationItemSystemMessageObject = "realtime.item"
)

// The status of the item. Has no effect on the conversation.
type RealtimeConversationItemSystemMessageStatus string

const (
	RealtimeConversationItemSystemMessageStatusCompleted  RealtimeConversationItemSystemMessageStatus = "completed"
	RealtimeConversationItemSystemMessageStatusIncomplete RealtimeConversationItemSystemMessageStatus = "incomplete"
	RealtimeConversationItemSystemMessageStatusInProgress RealtimeConversationItemSystemMessageStatus = "in_progress"
)

// A system message in a Realtime conversation can be used to provide additional
// context or instructions to the model. This is similar but distinct from the
// instruction prompt provided at the start of a conversation, as system messages
// can be added at any point in the conversation. For major changes to the
// conversation's behavior, use instructions, but for smaller updates (e.g. "the
// user is now asking about a different topic"), use system messages.
//
// The properties Content, Role, Type are required.
type RealtimeConversationItemSystemMessageParam struct {
	// The content of the message.
	Content []RealtimeConversationItemSystemMessageContentParam `json:"content,omitzero" api:"required"`
	// The unique ID of the item. This may be provided by the client or generated by
	// the server.
	ID param.Opt[string] `json:"id,omitzero"`
	// Identifier for the API object being returned - always `realtime.item`. Optional
	// when creating a new item.
	//
	// Any of "realtime.item".
	Object RealtimeConversationItemSystemMessageObject `json:"object,omitzero"`
	// The status of the item. Has no effect on the conversation.
	//
	// Any of "completed", "incomplete", "in_progress".
	Status RealtimeConversationItemSystemMessageStatus `json:"status,omitzero"`
	// The role of the message sender. Always `system`.
	//
	// This field can be elided, and will marshal its zero value as "system".
	Role constant.System `json:"role" default:"system"`
	// The type of the item. Always `message`.
	//
	// This field can be elided, and will marshal its zero value as "message".
	Type constant.Message `json:"type" default:"message"`
	paramObj
}

func (r RealtimeConversationItemSystemMessageParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeConversationItemSystemMessageParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeConversationItemSystemMessageParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type RealtimeConversationItemSystemMessageContentParam struct {
	// The text content.
	Text param.Opt[string] `json:"text,omitzero"`
	// The content type. Always `input_text` for system messages.
	//
	// Any of "input_text".
	Type string `json:"type,omitzero"`
	paramObj
}

func (r RealtimeConversationItemSystemMessageContentParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeConversationItemSystemMessageContentParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeConversationItemSystemMessageContentParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[RealtimeConversationItemSystemMessageContentParam](
		"type", "input_text",
	)
}

// A user message item in a Realtime conversation.
type RealtimeConversationItemUserMessage struct {
	// The content of the message.
	Content []RealtimeConversationItemUserMessageContent `json:"content" api:"required"`
	// The role of the message sender. Always `user`.
	Role constant.User `json:"role" default:"user"`
	// The type of the item. Always `message`.
	Type constant.Message `json:"type" default:"message"`
	// The unique ID of the item. This may be provided by the client or generated by
	// the server.
	ID string `json:"id"`
	// Identifier for the API object being returned - always `realtime.item`. Optional
	// when creating a new item.
	//
	// Any of "realtime.item".
	Object RealtimeConversationItemUserMessageObject `json:"object"`
	// The status of the item. Has no effect on the conversation.
	//
	// Any of "completed", "incomplete", "in_progress".
	Status RealtimeConversationItemUserMessageStatus `json:"status"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Content     respjson.Field
		Role        respjson.Field
		Type        respjson.Field
		ID          respjson.Field
		Object      respjson.Field
		Status      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeConversationItemUserMessage) RawJSON() string { return r.JSON.raw }
func (r *RealtimeConversationItemUserMessage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this RealtimeConversationItemUserMessage to a
// RealtimeConversationItemUserMessageParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// RealtimeConversationItemUserMessageParam.Overrides()
func (r RealtimeConversationItemUserMessage) ToParam() RealtimeConversationItemUserMessageParam {
	return param.Override[RealtimeConversationItemUserMessageParam](json.RawMessage(r.RawJSON()))
}

type RealtimeConversationItemUserMessageContent struct {
	// Base64-encoded audio bytes (for `input_audio`), these will be parsed as the
	// format specified in the session input audio type configuration. This defaults to
	// PCM 16-bit 24kHz mono if not specified.
	Audio string `json:"audio"`
	// The detail level of the image (for `input_image`). `auto` will default to
	// `high`.
	//
	// Any of "auto", "low", "high".
	Detail string `json:"detail"`
	// Base64-encoded image bytes (for `input_image`) as a data URI. For example
	// `data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA...`. Supported formats are PNG
	// and JPEG.
	ImageURL string `json:"image_url" format:"uri"`
	// The text content (for `input_text`).
	Text string `json:"text"`
	// Transcript of the audio (for `input_audio`). This is not sent to the model, but
	// will be attached to the message item for reference.
	Transcript string `json:"transcript"`
	// The content type (`input_text`, `input_audio`, or `input_image`).
	//
	// Any of "input_text", "input_audio", "input_image".
	Type string `json:"type"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Audio       respjson.Field
		Detail      respjson.Field
		ImageURL    respjson.Field
		Text        respjson.Field
		Transcript  respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeConversationItemUserMessageContent) RawJSON() string { return r.JSON.raw }
func (r *RealtimeConversationItemUserMessageContent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Identifier for the API object being returned - always `realtime.item`. Optional
// when creating a new item.
type RealtimeConversationItemUserMessageObject string

const (
	RealtimeConversationItemUserMessageObjectRealtimeItem RealtimeConversationItemUserMessageObject = "realtime.item"
)

// The status of the item. Has no effect on the conversation.
type RealtimeConversationItemUserMessageStatus string

const (
	RealtimeConversationItemUserMessageStatusCompleted  RealtimeConversationItemUserMessageStatus = "completed"
	RealtimeConversationItemUserMessageStatusIncomplete RealtimeConversationItemUserMessageStatus = "incomplete"
	RealtimeConversationItemUserMessageStatusInProgress RealtimeConversationItemUserMessageStatus = "in_progress"
)

// A user message item in a Realtime conversation.
//
// The properties Content, Role, Type are required.
type RealtimeConversationItemUserMessageParam struct {
	// The content of the message.
	Content []RealtimeConversationItemUserMessageContentParam `json:"content,omitzero" api:"required"`
	// The unique ID of the item. This may be provided by the client or generated by
	// the server.
	ID param.Opt[string] `json:"id,omitzero"`
	// Identifier for the API object being returned - always `realtime.item`. Optional
	// when creating a new item.
	//
	// Any of "realtime.item".
	Object RealtimeConversationItemUserMessageObject `json:"object,omitzero"`
	// The status of the item. Has no effect on the conversation.
	//
	// Any of "completed", "incomplete", "in_progress".
	Status RealtimeConversationItemUserMessageStatus `json:"status,omitzero"`
	// The role of the message sender. Always `user`.
	//
	// This field can be elided, and will marshal its zero value as "user".
	Role constant.User `json:"role" default:"user"`
	// The type of the item. Always `message`.
	//
	// This field can be elided, and will marshal its zero value as "message".
	Type constant.Message `json:"type" default:"message"`
	paramObj
}

func (r RealtimeConversationItemUserMessageParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeConversationItemUserMessageParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeConversationItemUserMessageParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type RealtimeConversationItemUserMessageContentParam struct {
	// Base64-encoded audio bytes (for `input_audio`), these will be parsed as the
	// format specified in the session input audio type configuration. This defaults to
	// PCM 16-bit 24kHz mono if not specified.
	Audio param.Opt[string] `json:"audio,omitzero"`
	// Base64-encoded image bytes (for `input_image`) as a data URI. For example
	// `data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA...`. Supported formats are PNG
	// and JPEG.
	ImageURL param.Opt[string] `json:"image_url,omitzero" format:"uri"`
	// The text content (for `input_text`).
	Text param.Opt[string] `json:"text,omitzero"`
	// Transcript of the audio (for `input_audio`). This is not sent to the model, but
	// will be attached to the message item for reference.
	Transcript param.Opt[string] `json:"transcript,omitzero"`
	// The detail level of the image (for `input_image`). `auto` will default to
	// `high`.
	//
	// Any of "auto", "low", "high".
	Detail string `json:"detail,omitzero"`
	// The content type (`input_text`, `input_audio`, or `input_image`).
	//
	// Any of "input_text", "input_audio", "input_image".
	Type string `json:"type,omitzero"`
	paramObj
}

func (r RealtimeConversationItemUserMessageContentParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeConversationItemUserMessageContentParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeConversationItemUserMessageContentParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[RealtimeConversationItemUserMessageContentParam](
		"detail", "auto", "low", "high",
	)
	apijson.RegisterFieldValidator[RealtimeConversationItemUserMessageContentParam](
		"type", "input_text", "input_audio", "input_image",
	)
}

// Details of the error.
type RealtimeError struct {
	// A human-readable error message.
	Message string `json:"message" api:"required"`
	// The type of error (e.g., "invalid_request_error", "server_error").
	Type string `json:"type" api:"required"`
	// Error code, if any.
	Code string `json:"code" api:"nullable"`
	// The event_id of the client event that caused the error, if applicable.
	EventID string `json:"event_id" api:"nullable"`
	// Parameter related to the error, if any.
	Param string `json:"param" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		Type        respjson.Field
		Code        respjson.Field
		EventID     respjson.Field
		Param       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeError) RawJSON() string { return r.JSON.raw }
func (r *RealtimeError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Returned when an error occurs, which could be a client problem or a server
// problem. Most errors are recoverable and the session will stay open, we
// recommend to implementers to monitor and log error messages by default.
type RealtimeErrorEvent struct {
	// Details of the error.
	Error RealtimeError `json:"error" api:"required"`
	// The unique ID of the server event.
	EventID string `json:"event_id" api:"required"`
	// The event type, must be `error`.
	Type constant.Error `json:"type" default:"error"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Error       respjson.Field
		EventID     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeErrorEvent) RawJSON() string { return r.JSON.raw }
func (r *RealtimeErrorEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type RealtimeFunctionTool struct {
	// The description of the function, including guidance on when and how to call it,
	// and guidance about what to tell the user when calling (if anything).
	Description string `json:"description"`
	// The name of the function.
	Name string `json:"name"`
	// Parameters of the function in JSON Schema.
	Parameters any `json:"parameters"`
	// The type of the tool, i.e. `function`.
	//
	// Any of "function".
	Type RealtimeFunctionToolType `json:"type"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Description respjson.Field
		Name        respjson.Field
		Parameters  respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeFunctionTool) RawJSON() string { return r.JSON.raw }
func (r *RealtimeFunctionTool) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this RealtimeFunctionTool to a RealtimeFunctionToolParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// RealtimeFunctionToolParam.Overrides()
func (r RealtimeFunctionTool) ToParam() RealtimeFunctionToolParam {
	return param.Override[RealtimeFunctionToolParam](json.RawMessage(r.RawJSON()))
}

// The type of the tool, i.e. `function`.
type RealtimeFunctionToolType string

const (
	RealtimeFunctionToolTypeFunction RealtimeFunctionToolType = "function"
)

type RealtimeFunctionToolParam struct {
	// The description of the function, including guidance on when and how to call it,
	// and guidance about what to tell the user when calling (if anything).
	Description param.Opt[string] `json:"description,omitzero"`
	// The name of the function.
	Name param.Opt[string] `json:"name,omitzero"`
	// Parameters of the function in JSON Schema.
	Parameters any `json:"parameters,omitzero"`
	// The type of the tool, i.e. `function`.
	//
	// Any of "function".
	Type RealtimeFunctionToolType `json:"type,omitzero"`
	paramObj
}

func (r RealtimeFunctionToolParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeFunctionToolParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeFunctionToolParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A Realtime item requesting human approval of a tool invocation.
type RealtimeMcpApprovalRequest struct {
	// The unique ID of the approval request.
	ID string `json:"id" api:"required"`
	// A JSON string of arguments for the tool.
	Arguments string `json:"arguments" api:"required"`
	// The name of the tool to run.
	Name string `json:"name" api:"required"`
	// The label of the MCP server making the request.
	ServerLabel string `json:"server_label" api:"required"`
	// The type of the item. Always `mcp_approval_request`.
	Type constant.McpApprovalRequest `json:"type" default:"mcp_approval_request"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Arguments   respjson.Field
		Name        respjson.Field
		ServerLabel respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeMcpApprovalRequest) RawJSON() string { return r.JSON.raw }
func (r *RealtimeMcpApprovalRequest) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this RealtimeMcpApprovalRequest to a
// RealtimeMcpApprovalRequestParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// RealtimeMcpApprovalRequestParam.Overrides()
func (r RealtimeMcpApprovalRequest) ToParam() RealtimeMcpApprovalRequestParam {
	return param.Override[RealtimeMcpApprovalRequestParam](json.RawMessage(r.RawJSON()))
}

// A Realtime item requesting human approval of a tool invocation.
//
// The properties ID, Arguments, Name, ServerLabel, Type are required.
type RealtimeMcpApprovalRequestParam struct {
	// The unique ID of the approval request.
	ID string `json:"id" api:"required"`
	// A JSON string of arguments for the tool.
	Arguments string `json:"arguments" api:"required"`
	// The name of the tool to run.
	Name string `json:"name" api:"required"`
	// The label of the MCP server making the request.
	ServerLabel string `json:"server_label" api:"required"`
	// The type of the item. Always `mcp_approval_request`.
	//
	// This field can be elided, and will marshal its zero value as
	// "mcp_approval_request".
	Type constant.McpApprovalRequest `json:"type" default:"mcp_approval_request"`
	paramObj
}

func (r RealtimeMcpApprovalRequestParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeMcpApprovalRequestParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeMcpApprovalRequestParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A Realtime item responding to an MCP approval request.
type RealtimeMcpApprovalResponse struct {
	// The unique ID of the approval response.
	ID string `json:"id" api:"required"`
	// The ID of the approval request being answered.
	ApprovalRequestID string `json:"approval_request_id" api:"required"`
	// Whether the request was approved.
	Approve bool `json:"approve" api:"required"`
	// The type of the item. Always `mcp_approval_response`.
	Type constant.McpApprovalResponse `json:"type" default:"mcp_approval_response"`
	// Optional reason for the decision.
	Reason string `json:"reason" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                respjson.Field
		ApprovalRequestID respjson.Field
		Approve           respjson.Field
		Type              respjson.Field
		Reason            respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeMcpApprovalResponse) RawJSON() string { return r.JSON.raw }
func (r *RealtimeMcpApprovalResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this RealtimeMcpApprovalResponse to a
// RealtimeMcpApprovalResponseParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// RealtimeMcpApprovalResponseParam.Overrides()
func (r RealtimeMcpApprovalResponse) ToParam() RealtimeMcpApprovalResponseParam {
	return param.Override[RealtimeMcpApprovalResponseParam](json.RawMessage(r.RawJSON()))
}

// A Realtime item responding to an MCP approval request.
//
// The properties ID, ApprovalRequestID, Approve, Type are required.
type RealtimeMcpApprovalResponseParam struct {
	// The unique ID of the approval response.
	ID string `json:"id" api:"required"`
	// The ID of the approval request being answered.
	ApprovalRequestID string `json:"approval_request_id" api:"required"`
	// Whether the request was approved.
	Approve bool `json:"approve" api:"required"`
	// Optional reason for the decision.
	Reason param.Opt[string] `json:"reason,omitzero"`
	// The type of the item. Always `mcp_approval_response`.
	//
	// This field can be elided, and will marshal its zero value as
	// "mcp_approval_response".
	Type constant.McpApprovalResponse `json:"type" default:"mcp_approval_response"`
	paramObj
}

func (r RealtimeMcpApprovalResponseParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeMcpApprovalResponseParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeMcpApprovalResponseParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A Realtime item listing tools available on an MCP server.
type RealtimeMcpListTools struct {
	// The label of the MCP server.
	ServerLabel string `json:"server_label" api:"required"`
	// The tools available on the server.
	Tools []RealtimeMcpListToolsTool `json:"tools" api:"required"`
	// The type of the item. Always `mcp_list_tools`.
	Type constant.McpListTools `json:"type" default:"mcp_list_tools"`
	// The unique ID of the list.
	ID string `json:"id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ServerLabel respjson.Field
		Tools       respjson.Field
		Type        respjson.Field
		ID          respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeMcpListTools) RawJSON() string { return r.JSON.raw }
func (r *RealtimeMcpListTools) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this RealtimeMcpListTools to a RealtimeMcpListToolsParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// RealtimeMcpListToolsParam.Overrides()
func (r RealtimeMcpListTools) ToParam() RealtimeMcpListToolsParam {
	return param.Override[RealtimeMcpListToolsParam](json.RawMessage(r.RawJSON()))
}

// A tool available on an MCP server.
type RealtimeMcpListToolsTool struct {
	// The JSON schema describing the tool's input.
	InputSchema any `json:"input_schema" api:"required"`
	// The name of the tool.
	Name string `json:"name" api:"required"`
	// Additional annotations about the tool.
	Annotations any `json:"annotations" api:"nullable"`
	// The description of the tool.
	Description string `json:"description" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		InputSchema respjson.Field
		Name        respjson.Field
		Annotations respjson.Field
		Description respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeMcpListToolsTool) RawJSON() string { return r.JSON.raw }
func (r *RealtimeMcpListToolsTool) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A Realtime item listing tools available on an MCP server.
//
// The properties ServerLabel, Tools, Type are required.
type RealtimeMcpListToolsParam struct {
	// The label of the MCP server.
	ServerLabel string `json:"server_label" api:"required"`
	// The tools available on the server.
	Tools []RealtimeMcpListToolsToolParam `json:"tools,omitzero" api:"required"`
	// The unique ID of the list.
	ID param.Opt[string] `json:"id,omitzero"`
	// The type of the item. Always `mcp_list_tools`.
	//
	// This field can be elided, and will marshal its zero value as "mcp_list_tools".
	Type constant.McpListTools `json:"type" default:"mcp_list_tools"`
	paramObj
}

func (r RealtimeMcpListToolsParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeMcpListToolsParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeMcpListToolsParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A tool available on an MCP server.
//
// The properties InputSchema, Name are required.
type RealtimeMcpListToolsToolParam struct {
	// The JSON schema describing the tool's input.
	InputSchema any `json:"input_schema,omitzero" api:"required"`
	// The name of the tool.
	Name string `json:"name" api:"required"`
	// The description of the tool.
	Description param.Opt[string] `json:"description,omitzero"`
	// Additional annotations about the tool.
	Annotations any `json:"annotations,omitzero"`
	paramObj
}

func (r RealtimeMcpListToolsToolParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeMcpListToolsToolParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeMcpListToolsToolParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type RealtimeMcpProtocolError struct {
	Code    int64                  `json:"code" api:"required"`
	Message string                 `json:"message" api:"required"`
	Type    constant.ProtocolError `json:"type" default:"protocol_error"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Code        respjson.Field
		Message     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeMcpProtocolError) RawJSON() string { return r.JSON.raw }
func (r *RealtimeMcpProtocolError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this RealtimeMcpProtocolError to a
// RealtimeMcpProtocolErrorParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// RealtimeMcpProtocolErrorParam.Overrides()
func (r RealtimeMcpProtocolError) ToParam() RealtimeMcpProtocolErrorParam {
	return param.Override[RealtimeMcpProtocolErrorParam](json.RawMessage(r.RawJSON()))
}

// The properties Code, Message, Type are required.
type RealtimeMcpProtocolErrorParam struct {
	Code    int64  `json:"code" api:"required"`
	Message string `json:"message" api:"required"`
	// This field can be elided, and will marshal its zero value as "protocol_error".
	Type constant.ProtocolError `json:"type" default:"protocol_error"`
	paramObj
}

func (r RealtimeMcpProtocolErrorParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeMcpProtocolErrorParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeMcpProtocolErrorParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A Realtime item representing an invocation of a tool on an MCP server.
type RealtimeMcpToolCall struct {
	// The unique ID of the tool call.
	ID string `json:"id" api:"required"`
	// A JSON string of the arguments passed to the tool.
	Arguments string `json:"arguments" api:"required"`
	// The name of the tool that was run.
	Name string `json:"name" api:"required"`
	// The label of the MCP server running the tool.
	ServerLabel string `json:"server_label" api:"required"`
	// The type of the item. Always `mcp_call`.
	Type constant.McpCall `json:"type" default:"mcp_call"`
	// The ID of an associated approval request, if any.
	ApprovalRequestID string `json:"approval_request_id" api:"nullable"`
	// The error from the tool call, if any.
	Error RealtimeMcpToolCallErrorUnion `json:"error" api:"nullable"`
	// The output from the tool call.
	Output string `json:"output" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                respjson.Field
		Arguments         respjson.Field
		Name              respjson.Field
		ServerLabel       respjson.Field
		Type              respjson.Field
		ApprovalRequestID respjson.Field
		Error             respjson.Field
		Output            respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeMcpToolCall) RawJSON() string { return r.JSON.raw }
func (r *RealtimeMcpToolCall) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this RealtimeMcpToolCall to a RealtimeMcpToolCallParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// RealtimeMcpToolCallParam.Overrides()
func (r RealtimeMcpToolCall) ToParam() RealtimeMcpToolCallParam {
	return param.Override[RealtimeMcpToolCallParam](json.RawMessage(r.RawJSON()))
}

// RealtimeMcpToolCallErrorUnion contains all possible properties and values from
// [RealtimeMcpProtocolError], [RealtimeMcpToolExecutionError],
// [RealtimeMcphttpError].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type RealtimeMcpToolCallErrorUnion struct {
	OfRealtimeMcpProtocolError      RealtimeMcpProtocolError      `json:",inline"`
	OfRealtimeMcpToolExecutionError RealtimeMcpToolExecutionError `json:",inline"`
	OfRealtimeMcphttpError          RealtimeMcphttpError          `json:",inline"`
	JSON                            struct {
		OfRealtimeMcpProtocolError      respjson.Field
		OfRealtimeMcpToolExecutionError respjson.Field
		OfRealtimeMcphttpError          respjson.Field
		raw                             string
	} `json:"-"`
}

func (u RealtimeMcpToolCallErrorUnion) AsProtocolError() (v RealtimeMcpProtocolError) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeMcpToolCallErrorUnion) AsToolExecutionError() (v RealtimeMcpToolExecutionError) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeMcpToolCallErrorUnion) AsHTTPError() (v RealtimeMcphttpError) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u RealtimeMcpToolCallErrorUnion) RawJSON() string { return u.JSON.raw }

func (r *RealtimeMcpToolCallErrorUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A Realtime item representing an invocation of a tool on an MCP server.
//
// The properties ID, Arguments, Name, ServerLabel, Type are required.
type RealtimeMcpToolCallParam struct {
	// The unique ID of the tool call.
	ID string `json:"id" api:"required"`
	// A JSON string of the arguments passed to the tool.
	Arguments string `json:"arguments" api:"required"`
	// The name of the tool that was run.
	Name string `json:"name" api:"required"`
	// The label of the MCP server running the tool.
	ServerLabel string `json:"server_label" api:"required"`
	// The ID of an associated approval request, if any.
	ApprovalRequestID param.Opt[string] `json:"approval_request_id,omitzero"`
	// The output from the tool call.
	Output param.Opt[string] `json:"output,omitzero"`
	// The error from the tool call, if any.
	Error RealtimeMcpToolCallErrorUnionParam `json:"error,omitzero"`
	// The type of the item. Always `mcp_call`.
	//
	// This field can be elided, and will marshal its zero value as "mcp_call".
	Type constant.McpCall `json:"type" default:"mcp_call"`
	paramObj
}

func (r RealtimeMcpToolCallParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeMcpToolCallParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeMcpToolCallParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type RealtimeMcpToolCallErrorUnionParam struct {
	OfProtocolError      *RealtimeMcpProtocolErrorParam      `json:",omitzero,inline"`
	OfToolExecutionError *RealtimeMcpToolExecutionErrorParam `json:",omitzero,inline"`
	OfHTTPError          *RealtimeMcphttpErrorParam          `json:",omitzero,inline"`
	paramUnion
}

func (u RealtimeMcpToolCallErrorUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfProtocolError, u.OfToolExecutionError, u.OfHTTPError)
}
func (u *RealtimeMcpToolCallErrorUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeMcpToolCallErrorUnionParam) GetCode() *int64 {
	if vt := u.OfProtocolError; vt != nil {
		return (*int64)(&vt.Code)
	} else if vt := u.OfHTTPError; vt != nil {
		return (*int64)(&vt.Code)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeMcpToolCallErrorUnionParam) GetMessage() *string {
	if vt := u.OfProtocolError; vt != nil {
		return (*string)(&vt.Message)
	} else if vt := u.OfToolExecutionError; vt != nil {
		return (*string)(&vt.Message)
	} else if vt := u.OfHTTPError; vt != nil {
		return (*string)(&vt.Message)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeMcpToolCallErrorUnionParam) GetType() *string {
	if vt := u.OfProtocolError; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfToolExecutionError; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfHTTPError; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[RealtimeMcpToolCallErrorUnionParam](
		"type",
		apijson.Discriminator[RealtimeMcpProtocolErrorParam]("protocol_error"),
		apijson.Discriminator[RealtimeMcpToolExecutionErrorParam]("tool_execution_error"),
		apijson.Discriminator[RealtimeMcphttpErrorParam]("http_error"),
	)
}

type RealtimeMcpToolExecutionError struct {
	Message string                      `json:"message" api:"required"`
	Type    constant.ToolExecutionError `json:"type" default:"tool_execution_error"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeMcpToolExecutionError) RawJSON() string { return r.JSON.raw }
func (r *RealtimeMcpToolExecutionError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this RealtimeMcpToolExecutionError to a
// RealtimeMcpToolExecutionErrorParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// RealtimeMcpToolExecutionErrorParam.Overrides()
func (r RealtimeMcpToolExecutionError) ToParam() RealtimeMcpToolExecutionErrorParam {
	return param.Override[RealtimeMcpToolExecutionErrorParam](json.RawMessage(r.RawJSON()))
}

// The properties Message, Type are required.
type RealtimeMcpToolExecutionErrorParam struct {
	Message string `json:"message" api:"required"`
	// This field can be elided, and will marshal its zero value as
	// "tool_execution_error".
	Type constant.ToolExecutionError `json:"type" default:"tool_execution_error"`
	paramObj
}

func (r RealtimeMcpToolExecutionErrorParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeMcpToolExecutionErrorParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeMcpToolExecutionErrorParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type RealtimeMcphttpError struct {
	Code    int64              `json:"code" api:"required"`
	Message string             `json:"message" api:"required"`
	Type    constant.HTTPError `json:"type" default:"http_error"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Code        respjson.Field
		Message     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeMcphttpError) RawJSON() string { return r.JSON.raw }
func (r *RealtimeMcphttpError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this RealtimeMcphttpError to a RealtimeMcphttpErrorParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// RealtimeMcphttpErrorParam.Overrides()
func (r RealtimeMcphttpError) ToParam() RealtimeMcphttpErrorParam {
	return param.Override[RealtimeMcphttpErrorParam](json.RawMessage(r.RawJSON()))
}

// The properties Code, Message, Type are required.
type RealtimeMcphttpErrorParam struct {
	Code    int64  `json:"code" api:"required"`
	Message string `json:"message" api:"required"`
	// This field can be elided, and will marshal its zero value as "http_error".
	Type constant.HTTPError `json:"type" default:"http_error"`
	paramObj
}

func (r RealtimeMcphttpErrorParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeMcphttpErrorParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeMcphttpErrorParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for reasoning-capable Realtime models such as `gpt-realtime-2`.
type RealtimeReasoning struct {
	// Constrains effort on reasoning for reasoning-capable Realtime models such as
	// `gpt-realtime-2`.
	//
	// Any of "minimal", "low", "medium", "high", "xhigh".
	Effort RealtimeReasoningEffort `json:"effort"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Effort      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeReasoning) RawJSON() string { return r.JSON.raw }
func (r *RealtimeReasoning) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this RealtimeReasoning to a RealtimeReasoningParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// RealtimeReasoningParam.Overrides()
func (r RealtimeReasoning) ToParam() RealtimeReasoningParam {
	return param.Override[RealtimeReasoningParam](json.RawMessage(r.RawJSON()))
}

// Configuration for reasoning-capable Realtime models such as `gpt-realtime-2`.
type RealtimeReasoningParam struct {
	// Constrains effort on reasoning for reasoning-capable Realtime models such as
	// `gpt-realtime-2`.
	//
	// Any of "minimal", "low", "medium", "high", "xhigh".
	Effort RealtimeReasoningEffort `json:"effort,omitzero"`
	paramObj
}

func (r RealtimeReasoningParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeReasoningParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeReasoningParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Constrains effort on reasoning for reasoning-capable Realtime models such as
// `gpt-realtime-2`.
type RealtimeReasoningEffort string

const (
	RealtimeReasoningEffortMinimal RealtimeReasoningEffort = "minimal"
	RealtimeReasoningEffortLow     RealtimeReasoningEffort = "low"
	RealtimeReasoningEffortMedium  RealtimeReasoningEffort = "medium"
	RealtimeReasoningEffortHigh    RealtimeReasoningEffort = "high"
	RealtimeReasoningEffortXhigh   RealtimeReasoningEffort = "xhigh"
)

// The response resource.
type RealtimeResponse struct {
	// The unique ID of the response, will look like `resp_1234`.
	ID string `json:"id"`
	// Configuration for audio output.
	Audio RealtimeResponseAudio `json:"audio"`
	// Which conversation the response is added to, determined by the `conversation`
	// field in the `response.create` event. If `auto`, the response will be added to
	// the default conversation and the value of `conversation_id` will be an id like
	// `conv_1234`. If `none`, the response will not be added to any conversation and
	// the value of `conversation_id` will be `null`. If responses are being triggered
	// automatically by VAD the response will be added to the default conversation
	ConversationID string `json:"conversation_id"`
	// Maximum number of output tokens for a single assistant response, inclusive of
	// tool calls, that was used in this response.
	MaxOutputTokens RealtimeResponseMaxOutputTokensUnion `json:"max_output_tokens"`
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard.
	//
	// Keys are strings with a maximum length of 64 characters. Values are strings with
	// a maximum length of 512 characters.
	Metadata shared.Metadata `json:"metadata" api:"nullable"`
	// The object type, must be `realtime.response`.
	//
	// Any of "realtime.response".
	Object RealtimeResponseObject `json:"object"`
	// The list of output items generated by the response.
	Output []ConversationItemUnion `json:"output"`
	// The set of modalities the model used to respond, currently the only possible
	// values are `[\"audio\"]`, `[\"text\"]`. Audio output always include a text
	// transcript. Setting the output to mode `text` will disable audio output from the
	// model.
	//
	// Any of "text", "audio".
	OutputModalities []string `json:"output_modalities"`
	// The final status of the response (`completed`, `cancelled`, `failed`, or
	// `incomplete`, `in_progress`).
	//
	// Any of "completed", "cancelled", "failed", "incomplete", "in_progress".
	Status RealtimeResponseStatus `json:"status"`
	// Additional details about the status.
	StatusDetails RealtimeResponseStatus `json:"status_details"`
	// Usage statistics for the Response, this will correspond to billing. A Realtime
	// API session will maintain a conversation context and append new Items to the
	// Conversation, thus output from previous turns (text and audio tokens) will
	// become the input for later turns.
	Usage RealtimeResponseUsage `json:"usage"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID               respjson.Field
		Audio            respjson.Field
		ConversationID   respjson.Field
		MaxOutputTokens  respjson.Field
		Metadata         respjson.Field
		Object           respjson.Field
		Output           respjson.Field
		OutputModalities respjson.Field
		Status           respjson.Field
		StatusDetails    respjson.Field
		Usage            respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeResponse) RawJSON() string { return r.JSON.raw }
func (r *RealtimeResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for audio output.
type RealtimeResponseAudio struct {
	Output RealtimeResponseAudioOutput `json:"output"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Output      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeResponseAudio) RawJSON() string { return r.JSON.raw }
func (r *RealtimeResponseAudio) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type RealtimeResponseAudioOutput struct {
	// The format of the output audio.
	Format RealtimeAudioFormatsUnion `json:"format"`
	// The voice the model uses to respond. Voice cannot be changed during the session
	// once the model has responded with audio at least once. Current voice options are
	// `alloy`, `ash`, `ballad`, `coral`, `echo`, `sage`, `shimmer`, `verse`, `marin`,
	// and `cedar`. We recommend `marin` and `cedar` for best quality.
	Voice string `json:"voice"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Format      respjson.Field
		Voice       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeResponseAudioOutput) RawJSON() string { return r.JSON.raw }
func (r *RealtimeResponseAudioOutput) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// RealtimeResponseMaxOutputTokensUnion contains all possible properties and values
// from [int64], [constant.Inf].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfInt OfInf]
type RealtimeResponseMaxOutputTokensUnion struct {
	// This field will be present if the value is a [int64] instead of an object.
	OfInt int64 `json:",inline"`
	// This field will be present if the value is a [constant.Inf] instead of an
	// object.
	OfInf constant.Inf `json:",inline"`
	JSON  struct {
		OfInt respjson.Field
		OfInf respjson.Field
		raw   string
	} `json:"-"`
}

func (u RealtimeResponseMaxOutputTokensUnion) AsInt() (v int64) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeResponseMaxOutputTokensUnion) AsInf() (v constant.Inf) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u RealtimeResponseMaxOutputTokensUnion) RawJSON() string { return u.JSON.raw }

func (r *RealtimeResponseMaxOutputTokensUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The object type, must be `realtime.response`.
type RealtimeResponseObject string

const (
	RealtimeResponseObjectRealtimeResponse RealtimeResponseObject = "realtime.response"
)

// The final status of the response (`completed`, `cancelled`, `failed`, or
// `incomplete`, `in_progress`).
type RealtimeResponseStatus string

const (
	RealtimeResponseStatusCompleted  RealtimeResponseStatus = "completed"
	RealtimeResponseStatusCancelled  RealtimeResponseStatus = "cancelled"
	RealtimeResponseStatusFailed     RealtimeResponseStatus = "failed"
	RealtimeResponseStatusIncomplete RealtimeResponseStatus = "incomplete"
	RealtimeResponseStatusInProgress RealtimeResponseStatus = "in_progress"
)

// Configuration for audio input and output.
type RealtimeResponseCreateAudioOutputParam struct {
	Output RealtimeResponseCreateAudioOutputOutputParam `json:"output,omitzero"`
	paramObj
}

func (r RealtimeResponseCreateAudioOutputParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeResponseCreateAudioOutputParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeResponseCreateAudioOutputParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type RealtimeResponseCreateAudioOutputOutputParam struct {
	// The format of the output audio.
	Format RealtimeAudioFormatsUnionParam `json:"format,omitzero"`
	// The voice the model uses to respond. Supported built-in voices are `alloy`,
	// `ash`, `ballad`, `coral`, `echo`, `sage`, `shimmer`, `verse`, `marin`, and
	// `cedar`. You may also provide a custom voice object with an `id`, for example
	// `{ "id": "voice_1234" }`. Voice cannot be changed during the session once the
	// model has responded with audio at least once. Custom voices must be created from
	// audio samples. Voices created from text prompts are supported only in Live. We
	// recommend `marin` and `cedar` for best quality.
	Voice RealtimeResponseCreateAudioOutputOutputVoiceUnionParam `json:"voice,omitzero"`
	paramObj
}

func (r RealtimeResponseCreateAudioOutputOutputParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeResponseCreateAudioOutputOutputParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeResponseCreateAudioOutputOutputParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type RealtimeResponseCreateAudioOutputOutputVoiceUnionParam struct {
	OfString param.Opt[string] `json:",omitzero,inline"`
	// Check if union is this variant with
	// !param.IsOmitted(union.OfRealtimeResponseCreateAudioOutputOutputVoiceString2)
	OfRealtimeResponseCreateAudioOutputOutputVoiceString2 param.Opt[string]                                    `json:",omitzero,inline"`
	OfRealtimeResponseCreateAudioOutputOutputVoiceID      *RealtimeResponseCreateAudioOutputOutputVoiceIDParam `json:",omitzero,inline"`
	paramUnion
}

func (u RealtimeResponseCreateAudioOutputOutputVoiceUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfString, u.OfRealtimeResponseCreateAudioOutputOutputVoiceString2, u.OfRealtimeResponseCreateAudioOutputOutputVoiceID)
}
func (u *RealtimeResponseCreateAudioOutputOutputVoiceUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

type RealtimeResponseCreateAudioOutputOutputVoiceString2 string

const (
	RealtimeResponseCreateAudioOutputOutputVoiceString2Alloy   RealtimeResponseCreateAudioOutputOutputVoiceString2 = "alloy"
	RealtimeResponseCreateAudioOutputOutputVoiceString2Ash     RealtimeResponseCreateAudioOutputOutputVoiceString2 = "ash"
	RealtimeResponseCreateAudioOutputOutputVoiceString2Ballad  RealtimeResponseCreateAudioOutputOutputVoiceString2 = "ballad"
	RealtimeResponseCreateAudioOutputOutputVoiceString2Coral   RealtimeResponseCreateAudioOutputOutputVoiceString2 = "coral"
	RealtimeResponseCreateAudioOutputOutputVoiceString2Echo    RealtimeResponseCreateAudioOutputOutputVoiceString2 = "echo"
	RealtimeResponseCreateAudioOutputOutputVoiceString2Sage    RealtimeResponseCreateAudioOutputOutputVoiceString2 = "sage"
	RealtimeResponseCreateAudioOutputOutputVoiceString2Shimmer RealtimeResponseCreateAudioOutputOutputVoiceString2 = "shimmer"
	RealtimeResponseCreateAudioOutputOutputVoiceString2Verse   RealtimeResponseCreateAudioOutputOutputVoiceString2 = "verse"
	RealtimeResponseCreateAudioOutputOutputVoiceString2Marin   RealtimeResponseCreateAudioOutputOutputVoiceString2 = "marin"
	RealtimeResponseCreateAudioOutputOutputVoiceString2Cedar   RealtimeResponseCreateAudioOutputOutputVoiceString2 = "cedar"
)

// Custom voice reference.
//
// The property ID is required.
type RealtimeResponseCreateAudioOutputOutputVoiceIDParam struct {
	// The custom voice ID, e.g. `voice_1234`.
	ID string `json:"id" api:"required"`
	paramObj
}

func (r RealtimeResponseCreateAudioOutputOutputVoiceIDParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeResponseCreateAudioOutputOutputVoiceIDParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeResponseCreateAudioOutputOutputVoiceIDParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Give the model access to additional tools via remote Model Context Protocol
// (MCP) servers.
// [Learn more about MCP](https://developers.openai.com/api/docs/guides/tools-connectors-mcp).
//
// The properties ServerLabel, Type are required.
type RealtimeResponseCreateMcpToolParam struct {
	// A label for this MCP server, used to identify it in tool calls.
	ServerLabel string `json:"server_label" api:"required"`
	// An OAuth access token that can be used with a remote MCP server, either with a
	// custom MCP server URL or a service connector. Your application must handle the
	// OAuth authorization flow and provide the token here.
	Authorization param.Opt[string] `json:"authorization,omitzero"`
	// Whether this MCP tool is deferred and discovered via tool search.
	DeferLoading param.Opt[bool] `json:"defer_loading,omitzero"`
	// Optional description of the MCP server, used to provide more context.
	ServerDescription param.Opt[string] `json:"server_description,omitzero"`
	// The URL for the MCP server. One of `server_url`, `connector_id`, or `tunnel_id`
	// must be provided.
	ServerURL param.Opt[string] `json:"server_url,omitzero" format:"uri"`
	// The Secure MCP Tunnel ID to use instead of a direct server URL. One of
	// `server_url`, `connector_id`, or `tunnel_id` must be provided.
	TunnelID param.Opt[string] `json:"tunnel_id,omitzero"`
	// The tool invocation context(s).
	//
	// Any of "direct", "programmatic".
	AllowedCallers []string `json:"allowed_callers,omitzero"`
	// List of allowed tool names or a filter object.
	AllowedTools RealtimeResponseCreateMcpToolAllowedToolsUnionParam `json:"allowed_tools,omitzero"`
	// Optional HTTP headers to send to the MCP server. Use for authentication or other
	// purposes.
	Headers map[string]string `json:"headers,omitzero"`
	// Specify which of the MCP server's tools require approval.
	RequireApproval RealtimeResponseCreateMcpToolRequireApprovalUnionParam `json:"require_approval,omitzero"`
	// Identifier for service connectors, like those available in ChatGPT. One of
	// `server_url`, `connector_id`, or `tunnel_id` must be provided. Learn more about
	// service connectors
	// [here](https://developers.openai.com/api/docs/guides/tools-connectors-mcp#connectors).
	//
	// This field is deprecated for models released after September 1, 2026. Use
	// `server_url` to connect to a remote MCP server, or `tunnel_id` to connect
	// through a Secure MCP Tunnel.
	//
	// Currently supported `connector_id` values are:
	//
	// - Dropbox: `connector_dropbox`
	// - Gmail: `connector_gmail`
	// - Google Calendar: `connector_googlecalendar`
	// - Google Drive: `connector_googledrive`
	// - Microsoft Teams: `connector_microsoftteams`
	// - Outlook Calendar: `connector_outlookcalendar`
	// - Outlook Email: `connector_outlookemail`
	// - SharePoint: `connector_sharepoint`
	//
	// Any of "connector_dropbox", "connector_gmail", "connector_googlecalendar",
	// "connector_googledrive", "connector_microsoftteams",
	// "connector_outlookcalendar", "connector_outlookemail", "connector_sharepoint".
	//
	// Deprecated: deprecated
	ConnectorID RealtimeResponseCreateMcpToolConnectorID `json:"connector_id,omitzero"`
	// The type of the MCP tool. Always `mcp`.
	//
	// This field can be elided, and will marshal its zero value as "mcp".
	Type constant.Mcp `json:"type" default:"mcp"`
	paramObj
}

func (r RealtimeResponseCreateMcpToolParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeResponseCreateMcpToolParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeResponseCreateMcpToolParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type RealtimeResponseCreateMcpToolAllowedToolsUnionParam struct {
	OfMcpAllowedTools []string                                                     `json:",omitzero,inline"`
	OfMcpToolFilter   *RealtimeResponseCreateMcpToolAllowedToolsMcpToolFilterParam `json:",omitzero,inline"`
	paramUnion
}

func (u RealtimeResponseCreateMcpToolAllowedToolsUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfMcpAllowedTools, u.OfMcpToolFilter)
}
func (u *RealtimeResponseCreateMcpToolAllowedToolsUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// A filter object to specify which tools are allowed.
type RealtimeResponseCreateMcpToolAllowedToolsMcpToolFilterParam struct {
	// Indicates whether or not a tool modifies data or is read-only. If an MCP server
	// is
	// [annotated with `readOnlyHint`](https://modelcontextprotocol.io/specification/2025-06-18/schema#toolannotations-readonlyhint),
	// it will match this filter.
	ReadOnly param.Opt[bool] `json:"read_only,omitzero"`
	// List of allowed tool names.
	ToolNames []string `json:"tool_names,omitzero"`
	paramObj
}

func (r RealtimeResponseCreateMcpToolAllowedToolsMcpToolFilterParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeResponseCreateMcpToolAllowedToolsMcpToolFilterParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeResponseCreateMcpToolAllowedToolsMcpToolFilterParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Identifier for service connectors, like those available in ChatGPT. One of
// `server_url`, `connector_id`, or `tunnel_id` must be provided. Learn more about
// service connectors
// [here](https://developers.openai.com/api/docs/guides/tools-connectors-mcp#connectors).
//
// This field is deprecated for models released after September 1, 2026. Use
// `server_url` to connect to a remote MCP server, or `tunnel_id` to connect
// through a Secure MCP Tunnel.
//
// Currently supported `connector_id` values are:
//
// - Dropbox: `connector_dropbox`
// - Gmail: `connector_gmail`
// - Google Calendar: `connector_googlecalendar`
// - Google Drive: `connector_googledrive`
// - Microsoft Teams: `connector_microsoftteams`
// - Outlook Calendar: `connector_outlookcalendar`
// - Outlook Email: `connector_outlookemail`
// - SharePoint: `connector_sharepoint`
type RealtimeResponseCreateMcpToolConnectorID string

const (
	RealtimeResponseCreateMcpToolConnectorIDConnectorDropbox         RealtimeResponseCreateMcpToolConnectorID = "connector_dropbox"
	RealtimeResponseCreateMcpToolConnectorIDConnectorGmail           RealtimeResponseCreateMcpToolConnectorID = "connector_gmail"
	RealtimeResponseCreateMcpToolConnectorIDConnectorGooglecalendar  RealtimeResponseCreateMcpToolConnectorID = "connector_googlecalendar"
	RealtimeResponseCreateMcpToolConnectorIDConnectorGoogledrive     RealtimeResponseCreateMcpToolConnectorID = "connector_googledrive"
	RealtimeResponseCreateMcpToolConnectorIDConnectorMicrosoftteams  RealtimeResponseCreateMcpToolConnectorID = "connector_microsoftteams"
	RealtimeResponseCreateMcpToolConnectorIDConnectorOutlookcalendar RealtimeResponseCreateMcpToolConnectorID = "connector_outlookcalendar"
	RealtimeResponseCreateMcpToolConnectorIDConnectorOutlookemail    RealtimeResponseCreateMcpToolConnectorID = "connector_outlookemail"
	RealtimeResponseCreateMcpToolConnectorIDConnectorSharepoint      RealtimeResponseCreateMcpToolConnectorID = "connector_sharepoint"
)

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type RealtimeResponseCreateMcpToolRequireApprovalUnionParam struct {
	OfRealtimeResponseCreateMcpToolRequireApprovalMcpToolApprovalFilter *RealtimeResponseCreateMcpToolRequireApprovalMcpToolApprovalFilterParam `json:",omitzero,inline"`
	// Check if union is this variant with
	// !param.IsOmitted(union.OfMcpToolApprovalSetting)
	OfMcpToolApprovalSetting param.Opt[string] `json:",omitzero,inline"`
	paramUnion
}

func (u RealtimeResponseCreateMcpToolRequireApprovalUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfRealtimeResponseCreateMcpToolRequireApprovalMcpToolApprovalFilter, u.OfMcpToolApprovalSetting)
}
func (u *RealtimeResponseCreateMcpToolRequireApprovalUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Specify which of the MCP server's tools require approval. Can be `always`,
// `never`, or a filter object associated with tools that require approval.
type RealtimeResponseCreateMcpToolRequireApprovalMcpToolApprovalFilterParam struct {
	// A filter object to specify which tools are allowed.
	Always RealtimeResponseCreateMcpToolRequireApprovalMcpToolApprovalFilterAlwaysParam `json:"always,omitzero"`
	// A filter object to specify which tools are allowed.
	Never RealtimeResponseCreateMcpToolRequireApprovalMcpToolApprovalFilterNeverParam `json:"never,omitzero"`
	paramObj
}

func (r RealtimeResponseCreateMcpToolRequireApprovalMcpToolApprovalFilterParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeResponseCreateMcpToolRequireApprovalMcpToolApprovalFilterParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeResponseCreateMcpToolRequireApprovalMcpToolApprovalFilterParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A filter object to specify which tools are allowed.
type RealtimeResponseCreateMcpToolRequireApprovalMcpToolApprovalFilterAlwaysParam struct {
	// Indicates whether or not a tool modifies data or is read-only. If an MCP server
	// is
	// [annotated with `readOnlyHint`](https://modelcontextprotocol.io/specification/2025-06-18/schema#toolannotations-readonlyhint),
	// it will match this filter.
	ReadOnly param.Opt[bool] `json:"read_only,omitzero"`
	// List of allowed tool names.
	ToolNames []string `json:"tool_names,omitzero"`
	paramObj
}

func (r RealtimeResponseCreateMcpToolRequireApprovalMcpToolApprovalFilterAlwaysParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeResponseCreateMcpToolRequireApprovalMcpToolApprovalFilterAlwaysParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeResponseCreateMcpToolRequireApprovalMcpToolApprovalFilterAlwaysParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A filter object to specify which tools are allowed.
type RealtimeResponseCreateMcpToolRequireApprovalMcpToolApprovalFilterNeverParam struct {
	// Indicates whether or not a tool modifies data or is read-only. If an MCP server
	// is
	// [annotated with `readOnlyHint`](https://modelcontextprotocol.io/specification/2025-06-18/schema#toolannotations-readonlyhint),
	// it will match this filter.
	ReadOnly param.Opt[bool] `json:"read_only,omitzero"`
	// List of allowed tool names.
	ToolNames []string `json:"tool_names,omitzero"`
	paramObj
}

func (r RealtimeResponseCreateMcpToolRequireApprovalMcpToolApprovalFilterNeverParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeResponseCreateMcpToolRequireApprovalMcpToolApprovalFilterNeverParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeResponseCreateMcpToolRequireApprovalMcpToolApprovalFilterNeverParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Specify a single approval policy for all tools. One of `always` or `never`. When
// set to `always`, all tools will require approval. When set to `never`, all tools
// will not require approval.
type RealtimeResponseCreateMcpToolRequireApprovalMcpToolApprovalSetting string

const (
	RealtimeResponseCreateMcpToolRequireApprovalMcpToolApprovalSettingAlways RealtimeResponseCreateMcpToolRequireApprovalMcpToolApprovalSetting = "always"
	RealtimeResponseCreateMcpToolRequireApprovalMcpToolApprovalSettingNever  RealtimeResponseCreateMcpToolRequireApprovalMcpToolApprovalSetting = "never"
)

// Create a new Realtime response with these parameters
type RealtimeResponseCreateParams struct {
	// The default system instructions (i.e. system message) prepended to model calls.
	// This field allows the client to guide the model on desired responses. The model
	// can be instructed on response content and format, (e.g. "be extremely succinct",
	// "act friendly", "here are examples of good responses") and on audio behavior
	// (e.g. "talk quickly", "inject emotion into your voice", "laugh frequently"). The
	// instructions are not guaranteed to be followed by the model, but they provide
	// guidance to the model on the desired behavior. Note that the server sets default
	// instructions which will be used if this field is not set and are visible in the
	// `session.created` event at the start of the session.
	Instructions param.Opt[string] `json:"instructions,omitzero"`
	// Whether the model may call multiple tools in parallel. Only supported by
	// reasoning Realtime models such as `gpt-realtime-2`.
	ParallelToolCalls param.Opt[bool] `json:"parallel_tool_calls,omitzero"`
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard.
	//
	// Keys are strings with a maximum length of 64 characters. Values are strings with
	// a maximum length of 512 characters.
	Metadata shared.Metadata `json:"metadata,omitzero"`
	// Reference to a prompt template and its variables.
	// [Learn more](https://developers.openai.com/api/docs/guides/text?api-mode=responses#version-prompts-in-code).
	Prompt responses.ResponsePromptParam `json:"prompt,omitzero"`
	// Configuration for audio input and output.
	Audio RealtimeResponseCreateAudioOutputParam `json:"audio,omitzero"`
	// Controls which conversation the response is added to. Currently supports `auto`
	// and `none`, with `auto` as the default value. The `auto` value means that the
	// contents of the response will be added to the default conversation. Set this to
	// `none` to create an out-of-band response which will not add items to default
	// conversation.
	Conversation RealtimeResponseCreateParamsConversation `json:"conversation,omitzero"`
	// Input items to include in the prompt for the model. Using this field creates a
	// new context for this Response instead of using the default conversation. An
	// empty array `[]` will clear the context for this Response. Note that this can
	// include references to items that previously appeared in the session using their
	// id.
	Input []ConversationItemUnionParam `json:"input,omitzero"`
	// Maximum number of output tokens for a single assistant response, inclusive of
	// tool calls. Provide an integer between 1 and 4096 to limit output tokens, or
	// `inf` for the maximum available tokens for a given model. Defaults to `inf`.
	MaxOutputTokens RealtimeResponseCreateParamsMaxOutputTokensUnion `json:"max_output_tokens,omitzero"`
	// The set of modalities the model used to respond, currently the only possible
	// values are `[\"audio\"]`, `[\"text\"]`. Audio output always include a text
	// transcript. Setting the output to mode `text` will disable audio output from the
	// model.
	//
	// Any of "text", "audio".
	OutputModalities []string `json:"output_modalities,omitzero"`
	// Configuration for reasoning-capable Realtime models such as `gpt-realtime-2`.
	Reasoning RealtimeReasoningParam `json:"reasoning,omitzero"`
	// How the model chooses tools. Provide one of the string modes or force a specific
	// function/MCP tool.
	ToolChoice RealtimeResponseCreateParamsToolChoiceUnion `json:"tool_choice,omitzero"`
	// Tools available to the model.
	Tools []RealtimeResponseCreateParamsToolUnion `json:"tools,omitzero"`
	paramObj
}

func (r RealtimeResponseCreateParams) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeResponseCreateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeResponseCreateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Controls which conversation the response is added to. Currently supports `auto`
// and `none`, with `auto` as the default value. The `auto` value means that the
// contents of the response will be added to the default conversation. Set this to
// `none` to create an out-of-band response which will not add items to default
// conversation.
type RealtimeResponseCreateParamsConversation = string

const (
	RealtimeResponseCreateParamsConversationAuto RealtimeResponseCreateParamsConversation = "auto"
	RealtimeResponseCreateParamsConversationNone RealtimeResponseCreateParamsConversation = "none"
)

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type RealtimeResponseCreateParamsMaxOutputTokensUnion struct {
	OfInt param.Opt[int64] `json:",omitzero,inline"`
	// Construct this variant with constant.ValueOf[constant.Inf]()
	OfInf constant.Inf `json:",omitzero,inline"`
	paramUnion
}

func (u RealtimeResponseCreateParamsMaxOutputTokensUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfInt, u.OfInf)
}
func (u *RealtimeResponseCreateParamsMaxOutputTokensUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type RealtimeResponseCreateParamsToolChoiceUnion struct {
	// Check if union is this variant with !param.IsOmitted(union.OfToolChoiceMode)
	OfToolChoiceMode param.Opt[responses.ToolChoiceOptions] `json:",omitzero,inline"`
	OfFunctionTool   *responses.ToolChoiceFunctionParam     `json:",omitzero,inline"`
	OfMcpTool        *responses.ToolChoiceMcpParam          `json:",omitzero,inline"`
	paramUnion
}

func (u RealtimeResponseCreateParamsToolChoiceUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfToolChoiceMode, u.OfFunctionTool, u.OfMcpTool)
}
func (u *RealtimeResponseCreateParamsToolChoiceUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeResponseCreateParamsToolChoiceUnion) GetServerLabel() *string {
	if vt := u.OfMcpTool; vt != nil {
		return &vt.ServerLabel
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeResponseCreateParamsToolChoiceUnion) GetName() *string {
	if vt := u.OfFunctionTool; vt != nil {
		return (*string)(&vt.Name)
	} else if vt := u.OfMcpTool; vt != nil && vt.Name.Valid() {
		return &vt.Name.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeResponseCreateParamsToolChoiceUnion) GetType() *string {
	if vt := u.OfFunctionTool; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfMcpTool; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type RealtimeResponseCreateParamsToolUnion struct {
	OfRealtimeFunctionTool *RealtimeFunctionToolParam          `json:",omitzero,inline"`
	OfMcp                  *RealtimeResponseCreateMcpToolParam `json:",omitzero,inline"`
	paramUnion
}

func (u RealtimeResponseCreateParamsToolUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfRealtimeFunctionTool, u.OfMcp)
}
func (u *RealtimeResponseCreateParamsToolUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeResponseCreateParamsToolUnion) GetDescription() *string {
	if vt := u.OfRealtimeFunctionTool; vt != nil && vt.Description.Valid() {
		return &vt.Description.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeResponseCreateParamsToolUnion) GetName() *string {
	if vt := u.OfRealtimeFunctionTool; vt != nil && vt.Name.Valid() {
		return &vt.Name.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeResponseCreateParamsToolUnion) GetParameters() *any {
	if vt := u.OfRealtimeFunctionTool; vt != nil {
		return &vt.Parameters
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeResponseCreateParamsToolUnion) GetServerLabel() *string {
	if vt := u.OfMcp; vt != nil {
		return &vt.ServerLabel
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeResponseCreateParamsToolUnion) GetAllowedCallers() []string {
	if vt := u.OfMcp; vt != nil {
		return vt.AllowedCallers
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeResponseCreateParamsToolUnion) GetAllowedTools() *RealtimeResponseCreateMcpToolAllowedToolsUnionParam {
	if vt := u.OfMcp; vt != nil {
		return &vt.AllowedTools
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeResponseCreateParamsToolUnion) GetAuthorization() *string {
	if vt := u.OfMcp; vt != nil && vt.Authorization.Valid() {
		return &vt.Authorization.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeResponseCreateParamsToolUnion) GetConnectorID() *string {
	if vt := u.OfMcp; vt != nil {
		return (*string)(&vt.ConnectorID)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeResponseCreateParamsToolUnion) GetDeferLoading() *bool {
	if vt := u.OfMcp; vt != nil && vt.DeferLoading.Valid() {
		return &vt.DeferLoading.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeResponseCreateParamsToolUnion) GetHeaders() map[string]string {
	if vt := u.OfMcp; vt != nil {
		return vt.Headers
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeResponseCreateParamsToolUnion) GetRequireApproval() *RealtimeResponseCreateMcpToolRequireApprovalUnionParam {
	if vt := u.OfMcp; vt != nil {
		return &vt.RequireApproval
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeResponseCreateParamsToolUnion) GetServerDescription() *string {
	if vt := u.OfMcp; vt != nil && vt.ServerDescription.Valid() {
		return &vt.ServerDescription.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeResponseCreateParamsToolUnion) GetServerURL() *string {
	if vt := u.OfMcp; vt != nil && vt.ServerURL.Valid() {
		return &vt.ServerURL.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeResponseCreateParamsToolUnion) GetTunnelID() *string {
	if vt := u.OfMcp; vt != nil && vt.TunnelID.Valid() {
		return &vt.TunnelID.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeResponseCreateParamsToolUnion) GetType() *string {
	if vt := u.OfRealtimeFunctionTool; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfMcp; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Usage statistics for the Response, this will correspond to billing. A Realtime
// API session will maintain a conversation context and append new Items to the
// Conversation, thus output from previous turns (text and audio tokens) will
// become the input for later turns.
type RealtimeResponseUsage struct {
	// Details about the input tokens used in the Response. Cached tokens are tokens
	// from previous turns in the conversation that are included as context for the
	// current response. Cached tokens here are counted as a subset of input tokens,
	// meaning input tokens will include cached and uncached tokens.
	InputTokenDetails RealtimeResponseUsageInputTokenDetails `json:"input_token_details"`
	// The number of input tokens used in the Response, including text and audio
	// tokens.
	InputTokens int64 `json:"input_tokens"`
	// Details about the output tokens used in the Response.
	OutputTokenDetails RealtimeResponseUsageOutputTokenDetails `json:"output_token_details"`
	// The number of output tokens sent in the Response, including text and audio
	// tokens.
	OutputTokens int64 `json:"output_tokens"`
	// The total number of tokens in the Response including input and output text and
	// audio tokens.
	TotalTokens int64 `json:"total_tokens"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		InputTokenDetails  respjson.Field
		InputTokens        respjson.Field
		OutputTokenDetails respjson.Field
		OutputTokens       respjson.Field
		TotalTokens        respjson.Field
		ExtraFields        map[string]respjson.Field
		raw                string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeResponseUsage) RawJSON() string { return r.JSON.raw }
func (r *RealtimeResponseUsage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Details about the input tokens used in the Response. Cached tokens are tokens
// from previous turns in the conversation that are included as context for the
// current response. Cached tokens here are counted as a subset of input tokens,
// meaning input tokens will include cached and uncached tokens.
type RealtimeResponseUsageInputTokenDetails struct {
	// The number of audio tokens used as input for the Response.
	AudioTokens int64 `json:"audio_tokens"`
	// The number of cached tokens used as input for the Response.
	CachedTokens int64 `json:"cached_tokens"`
	// Details about the cached tokens used as input for the Response.
	CachedTokensDetails RealtimeResponseUsageInputTokenDetailsCachedTokensDetails `json:"cached_tokens_details"`
	// The number of image tokens used as input for the Response.
	ImageTokens int64 `json:"image_tokens"`
	// The number of text tokens used as input for the Response.
	TextTokens int64 `json:"text_tokens"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		AudioTokens         respjson.Field
		CachedTokens        respjson.Field
		CachedTokensDetails respjson.Field
		ImageTokens         respjson.Field
		TextTokens          respjson.Field
		ExtraFields         map[string]respjson.Field
		raw                 string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeResponseUsageInputTokenDetails) RawJSON() string { return r.JSON.raw }
func (r *RealtimeResponseUsageInputTokenDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Details about the cached tokens used as input for the Response.
type RealtimeResponseUsageInputTokenDetailsCachedTokensDetails struct {
	// The number of cached audio tokens used as input for the Response.
	AudioTokens int64 `json:"audio_tokens"`
	// The number of cached image tokens used as input for the Response.
	ImageTokens int64 `json:"image_tokens"`
	// The number of cached text tokens used as input for the Response.
	TextTokens int64 `json:"text_tokens"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		AudioTokens respjson.Field
		ImageTokens respjson.Field
		TextTokens  respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeResponseUsageInputTokenDetailsCachedTokensDetails) RawJSON() string {
	return r.JSON.raw
}
func (r *RealtimeResponseUsageInputTokenDetailsCachedTokensDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Details about the output tokens used in the Response.
type RealtimeResponseUsageOutputTokenDetails struct {
	// The number of audio tokens used in the Response.
	AudioTokens int64 `json:"audio_tokens"`
	// The number of text tokens used in the Response.
	TextTokens int64 `json:"text_tokens"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		AudioTokens respjson.Field
		TextTokens  respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeResponseUsageOutputTokenDetails) RawJSON() string { return r.JSON.raw }
func (r *RealtimeResponseUsageOutputTokenDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// RealtimeServerEventUnion contains all possible properties and values from
// [ConversationCreatedEvent], [ConversationItemCreatedEvent],
// [ConversationItemDeletedEvent],
// [ConversationItemInputAudioTranscriptionCompletedEvent],
// [ConversationItemInputAudioTranscriptionDeltaEvent],
// [ConversationItemInputAudioTranscriptionFailedEvent],
// [RealtimeServerEventConversationItemRetrieved],
// [ConversationItemTruncatedEvent], [RealtimeErrorEvent],
// [InputAudioBufferClearedEvent], [InputAudioBufferCommittedEvent],
// [InputAudioBufferDtmfEventReceivedEvent], [InputAudioBufferSpeechStartedEvent],
// [InputAudioBufferSpeechStoppedEvent], [RateLimitsUpdatedEvent],
// [ResponseAudioDeltaEvent], [ResponseAudioDoneEvent],
// [ResponseAudioTranscriptDeltaEvent], [ResponseAudioTranscriptDoneEvent],
// [ResponseContentPartAddedEvent], [ResponseContentPartDoneEvent],
// [ResponseCreatedEvent], [ResponseDoneEvent],
// [ResponseFunctionCallArgumentsDeltaEvent],
// [ResponseFunctionCallArgumentsDoneEvent], [ResponseOutputItemAddedEvent],
// [ResponseOutputItemDoneEvent], [ResponseTextDeltaEvent],
// [ResponseTextDoneEvent], [SessionCreatedEvent], [SessionUpdatedEvent],
// [RealtimeServerEventOutputAudioBufferStarted],
// [RealtimeServerEventOutputAudioBufferStopped],
// [RealtimeServerEventOutputAudioBufferCleared], [ConversationItemAdded],
// [ConversationItemDone], [InputAudioBufferTimeoutTriggered],
// [ConversationItemInputAudioTranscriptionSegment], [McpListToolsInProgress],
// [McpListToolsCompleted], [McpListToolsFailed], [ResponseMcpCallArgumentsDelta],
// [ResponseMcpCallArgumentsDone], [ResponseMcpCallInProgress],
// [ResponseMcpCallCompleted], [ResponseMcpCallFailed].
//
// Use the [RealtimeServerEventUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type RealtimeServerEventUnion struct {
	OfConversationCreatedEvent                              ConversationCreatedEvent                              `json:",inline"`
	OfConversationItemCreatedEvent                          ConversationItemCreatedEvent                          `json:",inline"`
	OfConversationItemDeletedEvent                          ConversationItemDeletedEvent                          `json:",inline"`
	OfConversationItemInputAudioTranscriptionCompletedEvent ConversationItemInputAudioTranscriptionCompletedEvent `json:",inline"`
	OfConversationItemInputAudioTranscriptionDeltaEvent     ConversationItemInputAudioTranscriptionDeltaEvent     `json:",inline"`
	OfConversationItemInputAudioTranscriptionFailedEvent    ConversationItemInputAudioTranscriptionFailedEvent    `json:",inline"`
	OfRealtimeServerEventConversationItemRetrieved          RealtimeServerEventConversationItemRetrieved          `json:",inline"`
	OfConversationItemTruncatedEvent                        ConversationItemTruncatedEvent                        `json:",inline"`
	OfRealtimeErrorEvent                                    RealtimeErrorEvent                                    `json:",inline"`
	OfInputAudioBufferClearedEvent                          InputAudioBufferClearedEvent                          `json:",inline"`
	OfInputAudioBufferCommittedEvent                        InputAudioBufferCommittedEvent                        `json:",inline"`
	OfInputAudioBufferDtmfEventReceivedEvent                InputAudioBufferDtmfEventReceivedEvent                `json:",inline"`
	OfInputAudioBufferSpeechStartedEvent                    InputAudioBufferSpeechStartedEvent                    `json:",inline"`
	OfInputAudioBufferSpeechStoppedEvent                    InputAudioBufferSpeechStoppedEvent                    `json:",inline"`
	OfRateLimitsUpdatedEvent                                RateLimitsUpdatedEvent                                `json:",inline"`
	OfResponseAudioDeltaEvent                               ResponseAudioDeltaEvent                               `json:",inline"`
	OfResponseAudioDoneEvent                                ResponseAudioDoneEvent                                `json:",inline"`
	OfResponseAudioTranscriptDeltaEvent                     ResponseAudioTranscriptDeltaEvent                     `json:",inline"`
	OfResponseAudioTranscriptDoneEvent                      ResponseAudioTranscriptDoneEvent                      `json:",inline"`
	OfResponseContentPartAddedEvent                         ResponseContentPartAddedEvent                         `json:",inline"`
	OfResponseContentPartDoneEvent                          ResponseContentPartDoneEvent                          `json:",inline"`
	OfResponseCreatedEvent                                  ResponseCreatedEvent                                  `json:",inline"`
	OfResponseDoneEvent                                     ResponseDoneEvent                                     `json:",inline"`
	OfResponseFunctionCallArgumentsDeltaEvent               ResponseFunctionCallArgumentsDeltaEvent               `json:",inline"`
	OfResponseFunctionCallArgumentsDoneEvent                ResponseFunctionCallArgumentsDoneEvent                `json:",inline"`
	OfResponseOutputItemAddedEvent                          ResponseOutputItemAddedEvent                          `json:",inline"`
	OfResponseOutputItemDoneEvent                           ResponseOutputItemDoneEvent                           `json:",inline"`
	OfResponseTextDeltaEvent                                ResponseTextDeltaEvent                                `json:",inline"`
	OfResponseTextDoneEvent                                 ResponseTextDoneEvent                                 `json:",inline"`
	OfSessionCreatedEvent                                   SessionCreatedEvent                                   `json:",inline"`
	OfSessionUpdatedEvent                                   SessionUpdatedEvent                                   `json:",inline"`
	OfRealtimeServerEventOutputAudioBufferStarted           RealtimeServerEventOutputAudioBufferStarted           `json:",inline"`
	OfRealtimeServerEventOutputAudioBufferStopped           RealtimeServerEventOutputAudioBufferStopped           `json:",inline"`
	OfRealtimeServerEventOutputAudioBufferCleared           RealtimeServerEventOutputAudioBufferCleared           `json:",inline"`
	OfConversationItemAdded                                 ConversationItemAdded                                 `json:",inline"`
	OfConversationItemDone                                  ConversationItemDone                                  `json:",inline"`
	OfInputAudioBufferTimeoutTriggered                      InputAudioBufferTimeoutTriggered                      `json:",inline"`
	OfConversationItemInputAudioTranscriptionSegment        ConversationItemInputAudioTranscriptionSegment        `json:",inline"`
	OfMcpListToolsInProgress                                McpListToolsInProgress                                `json:",inline"`
	OfMcpListToolsCompleted                                 McpListToolsCompleted                                 `json:",inline"`
	OfMcpListToolsFailed                                    McpListToolsFailed                                    `json:",inline"`
	OfResponseMcpCallArgumentsDelta                         ResponseMcpCallArgumentsDelta                         `json:",inline"`
	OfResponseMcpCallArgumentsDone                          ResponseMcpCallArgumentsDone                          `json:",inline"`
	OfResponseMcpCallInProgress                             ResponseMcpCallInProgress                             `json:",inline"`
	OfResponseMcpCallCompleted                              ResponseMcpCallCompleted                              `json:",inline"`
	OfResponseMcpCallFailed                                 ResponseMcpCallFailed                                 `json:",inline"`
	Type                                                    string                                                `json:"type"`
	JSON                                                    struct {
		OfConversationCreatedEvent                              respjson.Field
		OfConversationItemCreatedEvent                          respjson.Field
		OfConversationItemDeletedEvent                          respjson.Field
		OfConversationItemInputAudioTranscriptionCompletedEvent respjson.Field
		OfConversationItemInputAudioTranscriptionDeltaEvent     respjson.Field
		OfConversationItemInputAudioTranscriptionFailedEvent    respjson.Field
		OfRealtimeServerEventConversationItemRetrieved          respjson.Field
		OfConversationItemTruncatedEvent                        respjson.Field
		OfRealtimeErrorEvent                                    respjson.Field
		OfInputAudioBufferClearedEvent                          respjson.Field
		OfInputAudioBufferCommittedEvent                        respjson.Field
		OfInputAudioBufferDtmfEventReceivedEvent                respjson.Field
		OfInputAudioBufferSpeechStartedEvent                    respjson.Field
		OfInputAudioBufferSpeechStoppedEvent                    respjson.Field
		OfRateLimitsUpdatedEvent                                respjson.Field
		OfResponseAudioDeltaEvent                               respjson.Field
		OfResponseAudioDoneEvent                                respjson.Field
		OfResponseAudioTranscriptDeltaEvent                     respjson.Field
		OfResponseAudioTranscriptDoneEvent                      respjson.Field
		OfResponseContentPartAddedEvent                         respjson.Field
		OfResponseContentPartDoneEvent                          respjson.Field
		OfResponseCreatedEvent                                  respjson.Field
		OfResponseDoneEvent                                     respjson.Field
		OfResponseFunctionCallArgumentsDeltaEvent               respjson.Field
		OfResponseFunctionCallArgumentsDoneEvent                respjson.Field
		OfResponseOutputItemAddedEvent                          respjson.Field
		OfResponseOutputItemDoneEvent                           respjson.Field
		OfResponseTextDeltaEvent                                respjson.Field
		OfResponseTextDoneEvent                                 respjson.Field
		OfSessionCreatedEvent                                   respjson.Field
		OfSessionUpdatedEvent                                   respjson.Field
		OfRealtimeServerEventOutputAudioBufferStarted           respjson.Field
		OfRealtimeServerEventOutputAudioBufferStopped           respjson.Field
		OfRealtimeServerEventOutputAudioBufferCleared           respjson.Field
		OfConversationItemAdded                                 respjson.Field
		OfConversationItemDone                                  respjson.Field
		OfInputAudioBufferTimeoutTriggered                      respjson.Field
		OfConversationItemInputAudioTranscriptionSegment        respjson.Field
		OfMcpListToolsInProgress                                respjson.Field
		OfMcpListToolsCompleted                                 respjson.Field
		OfMcpListToolsFailed                                    respjson.Field
		OfResponseMcpCallArgumentsDelta                         respjson.Field
		OfResponseMcpCallArgumentsDone                          respjson.Field
		OfResponseMcpCallInProgress                             respjson.Field
		OfResponseMcpCallCompleted                              respjson.Field
		OfResponseMcpCallFailed                                 respjson.Field
		Type                                                    respjson.Field
		raw                                                     string
	} `json:"-"`
}

// anyRealtimeServerEvent is implemented by each variant of
// [RealtimeServerEventUnion] to add type safety for the return type of
// [RealtimeServerEventUnion.AsAny]
type anyRealtimeServerEvent interface {
	implRealtimeServerEventUnion()
}

func (ConversationCreatedEvent) implRealtimeServerEventUnion()                              {}
func (ConversationItemCreatedEvent) implRealtimeServerEventUnion()                          {}
func (ConversationItemDeletedEvent) implRealtimeServerEventUnion()                          {}
func (ConversationItemInputAudioTranscriptionCompletedEvent) implRealtimeServerEventUnion() {}
func (ConversationItemInputAudioTranscriptionDeltaEvent) implRealtimeServerEventUnion()     {}
func (ConversationItemInputAudioTranscriptionFailedEvent) implRealtimeServerEventUnion()    {}
func (RealtimeServerEventConversationItemRetrieved) implRealtimeServerEventUnion()          {}
func (ConversationItemTruncatedEvent) implRealtimeServerEventUnion()                        {}
func (RealtimeErrorEvent) implRealtimeServerEventUnion()                                    {}
func (InputAudioBufferClearedEvent) implRealtimeServerEventUnion()                          {}
func (InputAudioBufferCommittedEvent) implRealtimeServerEventUnion()                        {}
func (InputAudioBufferDtmfEventReceivedEvent) implRealtimeServerEventUnion()                {}
func (InputAudioBufferSpeechStartedEvent) implRealtimeServerEventUnion()                    {}
func (InputAudioBufferSpeechStoppedEvent) implRealtimeServerEventUnion()                    {}
func (RateLimitsUpdatedEvent) implRealtimeServerEventUnion()                                {}
func (ResponseAudioDeltaEvent) implRealtimeServerEventUnion()                               {}
func (ResponseAudioDoneEvent) implRealtimeServerEventUnion()                                {}
func (ResponseAudioTranscriptDeltaEvent) implRealtimeServerEventUnion()                     {}
func (ResponseAudioTranscriptDoneEvent) implRealtimeServerEventUnion()                      {}
func (ResponseContentPartAddedEvent) implRealtimeServerEventUnion()                         {}
func (ResponseContentPartDoneEvent) implRealtimeServerEventUnion()                          {}
func (ResponseCreatedEvent) implRealtimeServerEventUnion()                                  {}
func (ResponseDoneEvent) implRealtimeServerEventUnion()                                     {}
func (ResponseFunctionCallArgumentsDeltaEvent) implRealtimeServerEventUnion()               {}
func (ResponseFunctionCallArgumentsDoneEvent) implRealtimeServerEventUnion()                {}
func (ResponseOutputItemAddedEvent) implRealtimeServerEventUnion()                          {}
func (ResponseOutputItemDoneEvent) implRealtimeServerEventUnion()                           {}
func (ResponseTextDeltaEvent) implRealtimeServerEventUnion()                                {}
func (ResponseTextDoneEvent) implRealtimeServerEventUnion()                                 {}
func (SessionCreatedEvent) implRealtimeServerEventUnion()                                   {}
func (SessionUpdatedEvent) implRealtimeServerEventUnion()                                   {}
func (RealtimeServerEventOutputAudioBufferStarted) implRealtimeServerEventUnion()           {}
func (RealtimeServerEventOutputAudioBufferStopped) implRealtimeServerEventUnion()           {}
func (RealtimeServerEventOutputAudioBufferCleared) implRealtimeServerEventUnion()           {}
func (ConversationItemAdded) implRealtimeServerEventUnion()                                 {}
func (ConversationItemDone) implRealtimeServerEventUnion()                                  {}
func (InputAudioBufferTimeoutTriggered) implRealtimeServerEventUnion()                      {}
func (ConversationItemInputAudioTranscriptionSegment) implRealtimeServerEventUnion()        {}
func (McpListToolsInProgress) implRealtimeServerEventUnion()                                {}
func (McpListToolsCompleted) implRealtimeServerEventUnion()                                 {}
func (McpListToolsFailed) implRealtimeServerEventUnion()                                    {}
func (ResponseMcpCallArgumentsDelta) implRealtimeServerEventUnion()                         {}
func (ResponseMcpCallArgumentsDone) implRealtimeServerEventUnion()                          {}
func (ResponseMcpCallInProgress) implRealtimeServerEventUnion()                             {}
func (ResponseMcpCallCompleted) implRealtimeServerEventUnion()                              {}
func (ResponseMcpCallFailed) implRealtimeServerEventUnion()                                 {}

// Use the following switch statement to find the correct variant
//
//	switch variant := RealtimeServerEventUnion.AsAny().(type) {
//	case realtime.ConversationCreatedEvent:
//	case realtime.ConversationItemCreatedEvent:
//	case realtime.ConversationItemDeletedEvent:
//	case realtime.ConversationItemInputAudioTranscriptionCompletedEvent:
//	case realtime.ConversationItemInputAudioTranscriptionDeltaEvent:
//	case realtime.ConversationItemInputAudioTranscriptionFailedEvent:
//	case realtime.RealtimeServerEventConversationItemRetrieved:
//	case realtime.ConversationItemTruncatedEvent:
//	case realtime.RealtimeErrorEvent:
//	case realtime.InputAudioBufferClearedEvent:
//	case realtime.InputAudioBufferCommittedEvent:
//	case realtime.InputAudioBufferDtmfEventReceivedEvent:
//	case realtime.InputAudioBufferSpeechStartedEvent:
//	case realtime.InputAudioBufferSpeechStoppedEvent:
//	case realtime.RateLimitsUpdatedEvent:
//	case realtime.ResponseAudioDeltaEvent:
//	case realtime.ResponseAudioDoneEvent:
//	case realtime.ResponseAudioTranscriptDeltaEvent:
//	case realtime.ResponseAudioTranscriptDoneEvent:
//	case realtime.ResponseContentPartAddedEvent:
//	case realtime.ResponseContentPartDoneEvent:
//	case realtime.ResponseCreatedEvent:
//	case realtime.ResponseDoneEvent:
//	case realtime.ResponseFunctionCallArgumentsDeltaEvent:
//	case realtime.ResponseFunctionCallArgumentsDoneEvent:
//	case realtime.ResponseOutputItemAddedEvent:
//	case realtime.ResponseOutputItemDoneEvent:
//	case realtime.ResponseTextDeltaEvent:
//	case realtime.ResponseTextDoneEvent:
//	case realtime.SessionCreatedEvent:
//	case realtime.SessionUpdatedEvent:
//	case realtime.RealtimeServerEventOutputAudioBufferStarted:
//	case realtime.RealtimeServerEventOutputAudioBufferStopped:
//	case realtime.RealtimeServerEventOutputAudioBufferCleared:
//	case realtime.ConversationItemAdded:
//	case realtime.ConversationItemDone:
//	case realtime.InputAudioBufferTimeoutTriggered:
//	case realtime.ConversationItemInputAudioTranscriptionSegment:
//	case realtime.McpListToolsInProgress:
//	case realtime.McpListToolsCompleted:
//	case realtime.McpListToolsFailed:
//	case realtime.ResponseMcpCallArgumentsDelta:
//	case realtime.ResponseMcpCallArgumentsDone:
//	case realtime.ResponseMcpCallInProgress:
//	case realtime.ResponseMcpCallCompleted:
//	case realtime.ResponseMcpCallFailed:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u RealtimeServerEventUnion) AsAny() anyRealtimeServerEvent {
	switch u.Type {
	case "conversation.created":
		return u.AsConversationCreatedEvent()
	case "conversation.item.created":
		return u.AsConversationItemCreatedEvent()
	case "conversation.item.deleted":
		return u.AsConversationItemDeletedEvent()
	case "conversation.item.input_audio_transcription.completed":
		return u.AsConversationItemInputAudioTranscriptionCompletedEvent()
	case "conversation.item.input_audio_transcription.delta":
		return u.AsConversationItemInputAudioTranscriptionDeltaEvent()
	case "conversation.item.input_audio_transcription.failed":
		return u.AsConversationItemInputAudioTranscriptionFailedEvent()
	case "conversation.item.retrieved":
		return u.AsRealtimeServerEventConversationItemRetrieved()
	case "conversation.item.truncated":
		return u.AsConversationItemTruncatedEvent()
	case "error":
		return u.AsRealtimeErrorEvent()
	case "input_audio_buffer.cleared":
		return u.AsInputAudioBufferClearedEvent()
	case "input_audio_buffer.committed":
		return u.AsInputAudioBufferCommittedEvent()
	case "input_audio_buffer.dtmf_event_received":
		return u.AsInputAudioBufferDtmfEventReceivedEvent()
	case "input_audio_buffer.speech_started":
		return u.AsInputAudioBufferSpeechStartedEvent()
	case "input_audio_buffer.speech_stopped":
		return u.AsInputAudioBufferSpeechStoppedEvent()
	case "rate_limits.updated":
		return u.AsRateLimitsUpdatedEvent()
	case "response.output_audio.delta":
		return u.AsResponseAudioDeltaEvent()
	case "response.output_audio.done":
		return u.AsResponseAudioDoneEvent()
	case "response.output_audio_transcript.delta":
		return u.AsResponseAudioTranscriptDeltaEvent()
	case "response.output_audio_transcript.done":
		return u.AsResponseAudioTranscriptDoneEvent()
	case "response.content_part.added":
		return u.AsResponseContentPartAddedEvent()
	case "response.content_part.done":
		return u.AsResponseContentPartDoneEvent()
	case "response.created":
		return u.AsResponseCreatedEvent()
	case "response.done":
		return u.AsResponseDoneEvent()
	case "response.function_call_arguments.delta":
		return u.AsResponseFunctionCallArgumentsDeltaEvent()
	case "response.function_call_arguments.done":
		return u.AsResponseFunctionCallArgumentsDoneEvent()
	case "response.output_item.added":
		return u.AsResponseOutputItemAddedEvent()
	case "response.output_item.done":
		return u.AsResponseOutputItemDoneEvent()
	case "response.output_text.delta":
		return u.AsResponseTextDeltaEvent()
	case "response.output_text.done":
		return u.AsResponseTextDoneEvent()
	case "session.created":
		return u.AsSessionCreatedEvent()
	case "session.updated":
		return u.AsSessionUpdatedEvent()
	case "output_audio_buffer.started":
		return u.AsRealtimeServerEventOutputAudioBufferStarted()
	case "output_audio_buffer.stopped":
		return u.AsRealtimeServerEventOutputAudioBufferStopped()
	case "output_audio_buffer.cleared":
		return u.AsRealtimeServerEventOutputAudioBufferCleared()
	case "conversation.item.added":
		return u.AsConversationItemAdded()
	case "conversation.item.done":
		return u.AsConversationItemDone()
	case "input_audio_buffer.timeout_triggered":
		return u.AsInputAudioBufferTimeoutTriggered()
	case "conversation.item.input_audio_transcription.segment":
		return u.AsConversationItemInputAudioTranscriptionSegment()
	case "mcp_list_tools.in_progress":
		return u.AsMcpListToolsInProgress()
	case "mcp_list_tools.completed":
		return u.AsMcpListToolsCompleted()
	case "mcp_list_tools.failed":
		return u.AsMcpListToolsFailed()
	case "response.mcp_call_arguments.delta":
		return u.AsResponseMcpCallArgumentsDelta()
	case "response.mcp_call_arguments.done":
		return u.AsResponseMcpCallArgumentsDone()
	case "response.mcp_call.in_progress":
		return u.AsResponseMcpCallInProgress()
	case "response.mcp_call.completed":
		return u.AsResponseMcpCallCompleted()
	case "response.mcp_call.failed":
		return u.AsResponseMcpCallFailed()
	}
	return nil
}

func (u RealtimeServerEventUnion) AsConversationCreatedEvent() (v ConversationCreatedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeServerEventUnion) AsConversationItemCreatedEvent() (v ConversationItemCreatedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeServerEventUnion) AsConversationItemDeletedEvent() (v ConversationItemDeletedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeServerEventUnion) AsConversationItemInputAudioTranscriptionCompletedEvent() (v ConversationItemInputAudioTranscriptionCompletedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeServerEventUnion) AsConversationItemInputAudioTranscriptionDeltaEvent() (v ConversationItemInputAudioTranscriptionDeltaEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeServerEventUnion) AsConversationItemInputAudioTranscriptionFailedEvent() (v ConversationItemInputAudioTranscriptionFailedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeServerEventUnion) AsRealtimeServerEventConversationItemRetrieved() (v RealtimeServerEventConversationItemRetrieved) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeServerEventUnion) AsConversationItemTruncatedEvent() (v ConversationItemTruncatedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeServerEventUnion) AsRealtimeErrorEvent() (v RealtimeErrorEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeServerEventUnion) AsInputAudioBufferClearedEvent() (v InputAudioBufferClearedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeServerEventUnion) AsInputAudioBufferCommittedEvent() (v InputAudioBufferCommittedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeServerEventUnion) AsInputAudioBufferDtmfEventReceivedEvent() (v InputAudioBufferDtmfEventReceivedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeServerEventUnion) AsInputAudioBufferSpeechStartedEvent() (v InputAudioBufferSpeechStartedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeServerEventUnion) AsInputAudioBufferSpeechStoppedEvent() (v InputAudioBufferSpeechStoppedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeServerEventUnion) AsRateLimitsUpdatedEvent() (v RateLimitsUpdatedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeServerEventUnion) AsResponseAudioDeltaEvent() (v ResponseAudioDeltaEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeServerEventUnion) AsResponseAudioDoneEvent() (v ResponseAudioDoneEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeServerEventUnion) AsResponseAudioTranscriptDeltaEvent() (v ResponseAudioTranscriptDeltaEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeServerEventUnion) AsResponseAudioTranscriptDoneEvent() (v ResponseAudioTranscriptDoneEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeServerEventUnion) AsResponseContentPartAddedEvent() (v ResponseContentPartAddedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeServerEventUnion) AsResponseContentPartDoneEvent() (v ResponseContentPartDoneEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeServerEventUnion) AsResponseCreatedEvent() (v ResponseCreatedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeServerEventUnion) AsResponseDoneEvent() (v ResponseDoneEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeServerEventUnion) AsResponseFunctionCallArgumentsDeltaEvent() (v ResponseFunctionCallArgumentsDeltaEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeServerEventUnion) AsResponseFunctionCallArgumentsDoneEvent() (v ResponseFunctionCallArgumentsDoneEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeServerEventUnion) AsResponseOutputItemAddedEvent() (v ResponseOutputItemAddedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeServerEventUnion) AsResponseOutputItemDoneEvent() (v ResponseOutputItemDoneEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeServerEventUnion) AsResponseTextDeltaEvent() (v ResponseTextDeltaEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeServerEventUnion) AsResponseTextDoneEvent() (v ResponseTextDoneEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeServerEventUnion) AsSessionCreatedEvent() (v SessionCreatedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeServerEventUnion) AsSessionUpdatedEvent() (v SessionUpdatedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeServerEventUnion) AsRealtimeServerEventOutputAudioBufferStarted() (v RealtimeServerEventOutputAudioBufferStarted) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeServerEventUnion) AsRealtimeServerEventOutputAudioBufferStopped() (v RealtimeServerEventOutputAudioBufferStopped) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeServerEventUnion) AsRealtimeServerEventOutputAudioBufferCleared() (v RealtimeServerEventOutputAudioBufferCleared) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeServerEventUnion) AsConversationItemAdded() (v ConversationItemAdded) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeServerEventUnion) AsConversationItemDone() (v ConversationItemDone) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeServerEventUnion) AsInputAudioBufferTimeoutTriggered() (v InputAudioBufferTimeoutTriggered) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeServerEventUnion) AsConversationItemInputAudioTranscriptionSegment() (v ConversationItemInputAudioTranscriptionSegment) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeServerEventUnion) AsMcpListToolsInProgress() (v McpListToolsInProgress) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeServerEventUnion) AsMcpListToolsCompleted() (v McpListToolsCompleted) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeServerEventUnion) AsMcpListToolsFailed() (v McpListToolsFailed) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeServerEventUnion) AsResponseMcpCallArgumentsDelta() (v ResponseMcpCallArgumentsDelta) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeServerEventUnion) AsResponseMcpCallArgumentsDone() (v ResponseMcpCallArgumentsDone) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeServerEventUnion) AsResponseMcpCallInProgress() (v ResponseMcpCallInProgress) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeServerEventUnion) AsResponseMcpCallCompleted() (v ResponseMcpCallCompleted) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeServerEventUnion) AsResponseMcpCallFailed() (v ResponseMcpCallFailed) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u RealtimeServerEventUnion) RawJSON() string { return u.JSON.raw }

func (r *RealtimeServerEventUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// unmarshalWebsocketEvent decodes only the event selected by its type. Public
// UnmarshalJSON keeps the existing general-purpose union decoding behavior.
func (r *RealtimeServerEventUnion) unmarshalWebsocketEvent(data []byte) error {
	var envelope struct {
		Type string `json:"type"`
		JSON struct {
			Type respjson.Field
		} `json:"-"`
	}
	if err := apijson.UnmarshalRoot(data, &envelope); err != nil {
		return err
	}
	eventType := envelope.Type
	*r = RealtimeServerEventUnion{Type: eventType}
	r.JSON.raw = string(data)
	r.JSON.Type = envelope.JSON.Type
	switch eventType {
	case "conversation.created":
		if err := apijson.UnmarshalRoot(data, &r.OfConversationCreatedEvent); err != nil {
			return err
		}
		r.JSON.OfConversationCreatedEvent = respjson.NewField(r.JSON.raw)
	case "conversation.item.created":
		if err := apijson.UnmarshalRoot(data, &r.OfConversationItemCreatedEvent); err != nil {
			return err
		}
		r.JSON.OfConversationItemCreatedEvent = respjson.NewField(r.JSON.raw)
	case "conversation.item.deleted":
		if err := apijson.UnmarshalRoot(data, &r.OfConversationItemDeletedEvent); err != nil {
			return err
		}
		r.JSON.OfConversationItemDeletedEvent = respjson.NewField(r.JSON.raw)
	case "conversation.item.input_audio_transcription.completed":
		if err := apijson.UnmarshalRoot(data, &r.OfConversationItemInputAudioTranscriptionCompletedEvent); err != nil {
			return err
		}
		r.JSON.OfConversationItemInputAudioTranscriptionCompletedEvent = respjson.NewField(r.JSON.raw)
	case "conversation.item.input_audio_transcription.delta":
		if err := apijson.UnmarshalRoot(data, &r.OfConversationItemInputAudioTranscriptionDeltaEvent); err != nil {
			return err
		}
		r.JSON.OfConversationItemInputAudioTranscriptionDeltaEvent = respjson.NewField(r.JSON.raw)
	case "conversation.item.input_audio_transcription.failed":
		if err := apijson.UnmarshalRoot(data, &r.OfConversationItemInputAudioTranscriptionFailedEvent); err != nil {
			return err
		}
		r.JSON.OfConversationItemInputAudioTranscriptionFailedEvent = respjson.NewField(r.JSON.raw)
	case "conversation.item.retrieved":
		if err := apijson.UnmarshalRoot(data, &r.OfRealtimeServerEventConversationItemRetrieved); err != nil {
			return err
		}
		r.JSON.OfRealtimeServerEventConversationItemRetrieved = respjson.NewField(r.JSON.raw)
	case "conversation.item.truncated":
		if err := apijson.UnmarshalRoot(data, &r.OfConversationItemTruncatedEvent); err != nil {
			return err
		}
		r.JSON.OfConversationItemTruncatedEvent = respjson.NewField(r.JSON.raw)
	case "error":
		if err := apijson.UnmarshalRoot(data, &r.OfRealtimeErrorEvent); err != nil {
			return err
		}
		r.JSON.OfRealtimeErrorEvent = respjson.NewField(r.JSON.raw)
	case "input_audio_buffer.cleared":
		if err := apijson.UnmarshalRoot(data, &r.OfInputAudioBufferClearedEvent); err != nil {
			return err
		}
		r.JSON.OfInputAudioBufferClearedEvent = respjson.NewField(r.JSON.raw)
	case "input_audio_buffer.committed":
		if err := apijson.UnmarshalRoot(data, &r.OfInputAudioBufferCommittedEvent); err != nil {
			return err
		}
		r.JSON.OfInputAudioBufferCommittedEvent = respjson.NewField(r.JSON.raw)
	case "input_audio_buffer.dtmf_event_received":
		if err := apijson.UnmarshalRoot(data, &r.OfInputAudioBufferDtmfEventReceivedEvent); err != nil {
			return err
		}
		r.JSON.OfInputAudioBufferDtmfEventReceivedEvent = respjson.NewField(r.JSON.raw)
	case "input_audio_buffer.speech_started":
		if err := apijson.UnmarshalRoot(data, &r.OfInputAudioBufferSpeechStartedEvent); err != nil {
			return err
		}
		r.JSON.OfInputAudioBufferSpeechStartedEvent = respjson.NewField(r.JSON.raw)
	case "input_audio_buffer.speech_stopped":
		if err := apijson.UnmarshalRoot(data, &r.OfInputAudioBufferSpeechStoppedEvent); err != nil {
			return err
		}
		r.JSON.OfInputAudioBufferSpeechStoppedEvent = respjson.NewField(r.JSON.raw)
	case "rate_limits.updated":
		if err := apijson.UnmarshalRoot(data, &r.OfRateLimitsUpdatedEvent); err != nil {
			return err
		}
		r.JSON.OfRateLimitsUpdatedEvent = respjson.NewField(r.JSON.raw)
	case "response.output_audio.delta":
		if err := apijson.UnmarshalRoot(data, &r.OfResponseAudioDeltaEvent); err != nil {
			return err
		}
		r.JSON.OfResponseAudioDeltaEvent = respjson.NewField(r.JSON.raw)
	case "response.output_audio.done":
		if err := apijson.UnmarshalRoot(data, &r.OfResponseAudioDoneEvent); err != nil {
			return err
		}
		r.JSON.OfResponseAudioDoneEvent = respjson.NewField(r.JSON.raw)
	case "response.output_audio_transcript.delta":
		if err := apijson.UnmarshalRoot(data, &r.OfResponseAudioTranscriptDeltaEvent); err != nil {
			return err
		}
		r.JSON.OfResponseAudioTranscriptDeltaEvent = respjson.NewField(r.JSON.raw)
	case "response.output_audio_transcript.done":
		if err := apijson.UnmarshalRoot(data, &r.OfResponseAudioTranscriptDoneEvent); err != nil {
			return err
		}
		r.JSON.OfResponseAudioTranscriptDoneEvent = respjson.NewField(r.JSON.raw)
	case "response.content_part.added":
		if err := apijson.UnmarshalRoot(data, &r.OfResponseContentPartAddedEvent); err != nil {
			return err
		}
		r.JSON.OfResponseContentPartAddedEvent = respjson.NewField(r.JSON.raw)
	case "response.content_part.done":
		if err := apijson.UnmarshalRoot(data, &r.OfResponseContentPartDoneEvent); err != nil {
			return err
		}
		r.JSON.OfResponseContentPartDoneEvent = respjson.NewField(r.JSON.raw)
	case "response.created":
		if err := apijson.UnmarshalRoot(data, &r.OfResponseCreatedEvent); err != nil {
			return err
		}
		r.JSON.OfResponseCreatedEvent = respjson.NewField(r.JSON.raw)
	case "response.done":
		if err := apijson.UnmarshalRoot(data, &r.OfResponseDoneEvent); err != nil {
			return err
		}
		r.JSON.OfResponseDoneEvent = respjson.NewField(r.JSON.raw)
	case "response.function_call_arguments.delta":
		if err := apijson.UnmarshalRoot(data, &r.OfResponseFunctionCallArgumentsDeltaEvent); err != nil {
			return err
		}
		r.JSON.OfResponseFunctionCallArgumentsDeltaEvent = respjson.NewField(r.JSON.raw)
	case "response.function_call_arguments.done":
		if err := apijson.UnmarshalRoot(data, &r.OfResponseFunctionCallArgumentsDoneEvent); err != nil {
			return err
		}
		r.JSON.OfResponseFunctionCallArgumentsDoneEvent = respjson.NewField(r.JSON.raw)
	case "response.output_item.added":
		if err := apijson.UnmarshalRoot(data, &r.OfResponseOutputItemAddedEvent); err != nil {
			return err
		}
		r.JSON.OfResponseOutputItemAddedEvent = respjson.NewField(r.JSON.raw)
	case "response.output_item.done":
		if err := apijson.UnmarshalRoot(data, &r.OfResponseOutputItemDoneEvent); err != nil {
			return err
		}
		r.JSON.OfResponseOutputItemDoneEvent = respjson.NewField(r.JSON.raw)
	case "response.output_text.delta":
		if err := apijson.UnmarshalRoot(data, &r.OfResponseTextDeltaEvent); err != nil {
			return err
		}
		r.JSON.OfResponseTextDeltaEvent = respjson.NewField(r.JSON.raw)
	case "response.output_text.done":
		if err := apijson.UnmarshalRoot(data, &r.OfResponseTextDoneEvent); err != nil {
			return err
		}
		r.JSON.OfResponseTextDoneEvent = respjson.NewField(r.JSON.raw)
	case "session.created":
		if err := apijson.UnmarshalRoot(data, &r.OfSessionCreatedEvent); err != nil {
			return err
		}
		r.JSON.OfSessionCreatedEvent = respjson.NewField(r.JSON.raw)
	case "session.updated":
		if err := apijson.UnmarshalRoot(data, &r.OfSessionUpdatedEvent); err != nil {
			return err
		}
		r.JSON.OfSessionUpdatedEvent = respjson.NewField(r.JSON.raw)
	case "output_audio_buffer.started":
		if err := apijson.UnmarshalRoot(data, &r.OfRealtimeServerEventOutputAudioBufferStarted); err != nil {
			return err
		}
		r.JSON.OfRealtimeServerEventOutputAudioBufferStarted = respjson.NewField(r.JSON.raw)
	case "output_audio_buffer.stopped":
		if err := apijson.UnmarshalRoot(data, &r.OfRealtimeServerEventOutputAudioBufferStopped); err != nil {
			return err
		}
		r.JSON.OfRealtimeServerEventOutputAudioBufferStopped = respjson.NewField(r.JSON.raw)
	case "output_audio_buffer.cleared":
		if err := apijson.UnmarshalRoot(data, &r.OfRealtimeServerEventOutputAudioBufferCleared); err != nil {
			return err
		}
		r.JSON.OfRealtimeServerEventOutputAudioBufferCleared = respjson.NewField(r.JSON.raw)
	case "conversation.item.added":
		if err := apijson.UnmarshalRoot(data, &r.OfConversationItemAdded); err != nil {
			return err
		}
		r.JSON.OfConversationItemAdded = respjson.NewField(r.JSON.raw)
	case "conversation.item.done":
		if err := apijson.UnmarshalRoot(data, &r.OfConversationItemDone); err != nil {
			return err
		}
		r.JSON.OfConversationItemDone = respjson.NewField(r.JSON.raw)
	case "input_audio_buffer.timeout_triggered":
		if err := apijson.UnmarshalRoot(data, &r.OfInputAudioBufferTimeoutTriggered); err != nil {
			return err
		}
		r.JSON.OfInputAudioBufferTimeoutTriggered = respjson.NewField(r.JSON.raw)
	case "conversation.item.input_audio_transcription.segment":
		if err := apijson.UnmarshalRoot(data, &r.OfConversationItemInputAudioTranscriptionSegment); err != nil {
			return err
		}
		r.JSON.OfConversationItemInputAudioTranscriptionSegment = respjson.NewField(r.JSON.raw)
	case "mcp_list_tools.in_progress":
		if err := apijson.UnmarshalRoot(data, &r.OfMcpListToolsInProgress); err != nil {
			return err
		}
		r.JSON.OfMcpListToolsInProgress = respjson.NewField(r.JSON.raw)
	case "mcp_list_tools.completed":
		if err := apijson.UnmarshalRoot(data, &r.OfMcpListToolsCompleted); err != nil {
			return err
		}
		r.JSON.OfMcpListToolsCompleted = respjson.NewField(r.JSON.raw)
	case "mcp_list_tools.failed":
		if err := apijson.UnmarshalRoot(data, &r.OfMcpListToolsFailed); err != nil {
			return err
		}
		r.JSON.OfMcpListToolsFailed = respjson.NewField(r.JSON.raw)
	case "response.mcp_call_arguments.delta":
		if err := apijson.UnmarshalRoot(data, &r.OfResponseMcpCallArgumentsDelta); err != nil {
			return err
		}
		r.JSON.OfResponseMcpCallArgumentsDelta = respjson.NewField(r.JSON.raw)
	case "response.mcp_call_arguments.done":
		if err := apijson.UnmarshalRoot(data, &r.OfResponseMcpCallArgumentsDone); err != nil {
			return err
		}
		r.JSON.OfResponseMcpCallArgumentsDone = respjson.NewField(r.JSON.raw)
	case "response.mcp_call.in_progress":
		if err := apijson.UnmarshalRoot(data, &r.OfResponseMcpCallInProgress); err != nil {
			return err
		}
		r.JSON.OfResponseMcpCallInProgress = respjson.NewField(r.JSON.raw)
	case "response.mcp_call.completed":
		if err := apijson.UnmarshalRoot(data, &r.OfResponseMcpCallCompleted); err != nil {
			return err
		}
		r.JSON.OfResponseMcpCallCompleted = respjson.NewField(r.JSON.raw)
	case "response.mcp_call.failed":
		if err := apijson.UnmarshalRoot(data, &r.OfResponseMcpCallFailed); err != nil {
			return err
		}
		r.JSON.OfResponseMcpCallFailed = respjson.NewField(r.JSON.raw)
	}
	return nil
}

// Returned when a conversation item is retrieved with
// `conversation.item.retrieve`. This is provided as a way to fetch the server's
// representation of an item, for example to get access to the post-processed audio
// data after noise cancellation and VAD. It includes the full content of the Item,
// including audio data.
type RealtimeServerEventConversationItemRetrieved struct {
	// The unique ID of the server event.
	EventID string `json:"event_id" api:"required"`
	// A single item within a Realtime conversation.
	Item ConversationItemUnion `json:"item" api:"required"`
	// The event type, must be `conversation.item.retrieved`.
	Type constant.ConversationItemRetrieved `json:"type" default:"conversation.item.retrieved"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventID     respjson.Field
		Item        respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeServerEventConversationItemRetrieved) RawJSON() string { return r.JSON.raw }
func (r *RealtimeServerEventConversationItemRetrieved) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// **WebRTC/SIP Only:** Emitted when the server begins streaming audio to the
// client. This event is emitted after an audio content part has been added
// (`response.content_part.added`) to the response.
// [Learn more](https://developers.openai.com/api/docs/guides/realtime-conversations#client-and-server-events-for-audio-in-webrtc).
type RealtimeServerEventOutputAudioBufferStarted struct {
	// The unique ID of the server event.
	EventID string `json:"event_id" api:"required"`
	// The unique ID of the response that produced the audio.
	ResponseID string `json:"response_id" api:"required"`
	// The event type, must be `output_audio_buffer.started`.
	Type constant.OutputAudioBufferStarted `json:"type" default:"output_audio_buffer.started"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventID     respjson.Field
		ResponseID  respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeServerEventOutputAudioBufferStarted) RawJSON() string { return r.JSON.raw }
func (r *RealtimeServerEventOutputAudioBufferStarted) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// **WebRTC/SIP Only:** Emitted when the output audio buffer has been completely
// drained on the server, and no more audio is forthcoming. This event is emitted
// after the full response data has been sent to the client (`response.done`).
// [Learn more](https://developers.openai.com/api/docs/guides/realtime-conversations#client-and-server-events-for-audio-in-webrtc).
type RealtimeServerEventOutputAudioBufferStopped struct {
	// The unique ID of the server event.
	EventID string `json:"event_id" api:"required"`
	// The unique ID of the response that produced the audio.
	ResponseID string `json:"response_id" api:"required"`
	// The event type, must be `output_audio_buffer.stopped`.
	Type constant.OutputAudioBufferStopped `json:"type" default:"output_audio_buffer.stopped"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventID     respjson.Field
		ResponseID  respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeServerEventOutputAudioBufferStopped) RawJSON() string { return r.JSON.raw }
func (r *RealtimeServerEventOutputAudioBufferStopped) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// **WebRTC/SIP Only:** Emitted when the output audio buffer is cleared. This
// happens either in VAD mode when the user has interrupted
// (`input_audio_buffer.speech_started`), or when the client has emitted the
// `output_audio_buffer.clear` event to manually cut off the current audio
// response.
// [Learn more](https://developers.openai.com/api/docs/guides/realtime-conversations#client-and-server-events-for-audio-in-webrtc).
type RealtimeServerEventOutputAudioBufferCleared struct {
	// The unique ID of the server event.
	EventID string `json:"event_id" api:"required"`
	// The unique ID of the response that produced the audio.
	ResponseID string `json:"response_id" api:"required"`
	// The event type, must be `output_audio_buffer.cleared`.
	Type constant.OutputAudioBufferCleared `json:"type" default:"output_audio_buffer.cleared"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventID     respjson.Field
		ResponseID  respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeServerEventOutputAudioBufferCleared) RawJSON() string { return r.JSON.raw }
func (r *RealtimeServerEventOutputAudioBufferCleared) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Realtime session object configuration.
type RealtimeSessionCreateRequest struct {
	// The type of session to create. Always `realtime` for the Realtime API.
	Type constant.Realtime `json:"type" default:"realtime"`
	// Configuration for input and output audio.
	Audio RealtimeAudioConfig `json:"audio"`
	// Additional fields to include in server outputs.
	//
	// `item.input_audio_transcription.logprobs`: Include logprobs for input audio
	// transcription.
	//
	// Any of "item.input_audio_transcription.logprobs".
	Include []string `json:"include"`
	// The default system instructions (i.e. system message) prepended to model calls.
	// This field allows the client to guide the model on desired responses. The model
	// can be instructed on response content and format, (e.g. "be extremely succinct",
	// "act friendly", "here are examples of good responses") and on audio behavior
	// (e.g. "talk quickly", "inject emotion into your voice", "laugh frequently"). The
	// instructions are not guaranteed to be followed by the model, but they provide
	// guidance to the model on the desired behavior.
	//
	// Note that the server sets default instructions which will be used if this field
	// is not set and are visible in the `session.created` event at the start of the
	// session.
	Instructions string `json:"instructions"`
	// Maximum number of output tokens for a single assistant response, inclusive of
	// tool calls. Provide an integer between 1 and 4096 to limit output tokens, or
	// `inf` for the maximum available tokens for a given model. Defaults to `inf`.
	MaxOutputTokens RealtimeSessionCreateRequestMaxOutputTokensUnion `json:"max_output_tokens"`
	// The Realtime model used for this session.
	Model RealtimeSessionCreateRequestModel `json:"model"`
	// The set of modalities the model can respond with. It defaults to `["audio"]`,
	// indicating that the model will respond with audio plus a transcript. `["text"]`
	// can be used to make the model respond with text only. It is not possible to
	// request both `text` and `audio` at the same time.
	//
	// Any of "text", "audio".
	OutputModalities []string `json:"output_modalities"`
	// Whether the model may call multiple tools in parallel. Only supported by
	// reasoning Realtime models such as `gpt-realtime-2`.
	ParallelToolCalls bool `json:"parallel_tool_calls"`
	// Reference to a prompt template and its variables.
	// [Learn more](https://developers.openai.com/api/docs/guides/text?api-mode=responses#version-prompts-in-code).
	Prompt responses.ResponsePrompt `json:"prompt" api:"nullable"`
	// Configuration for reasoning-capable Realtime models such as `gpt-realtime-2`.
	Reasoning RealtimeReasoning `json:"reasoning"`
	// How the model chooses tools. Provide one of the string modes or force a specific
	// function/MCP tool.
	ToolChoice RealtimeToolChoiceConfigUnion `json:"tool_choice"`
	// Tools available to the model.
	Tools RealtimeToolsConfig `json:"tools"`
	// Realtime API can write session traces to the
	// [Traces Dashboard](https://platform.openai.com/logs?api=traces). Set to null to
	// disable tracing. Once tracing is enabled for a session, the configuration cannot
	// be modified.
	//
	// `auto` will create a trace for the session with default values for the workflow
	// name, group id, and metadata.
	Tracing RealtimeTracingConfigUnion `json:"tracing" api:"nullable"`
	// When the number of tokens in a conversation exceeds the model's input token
	// limit, the conversation be truncated, meaning messages (starting from the
	// oldest) will not be included in the model's context. A 32k context model with
	// 4,096 max output tokens can only include 28,224 tokens in the context before
	// truncation occurs.
	//
	// Clients can configure truncation behavior to truncate with a lower max token
	// limit, which is an effective way to control token usage and cost.
	//
	// Truncation will reduce the number of cached tokens on the next turn (busting the
	// cache), since messages are dropped from the beginning of the context. However,
	// clients can also configure truncation to retain messages up to a fraction of the
	// maximum context size, which will reduce the need for future truncations and thus
	// improve the cache rate.
	//
	// Truncation can be disabled entirely, which means the server will never truncate
	// but would instead return an error if the conversation exceeds the model's input
	// token limit.
	Truncation RealtimeTruncationUnion `json:"truncation"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type              respjson.Field
		Audio             respjson.Field
		Include           respjson.Field
		Instructions      respjson.Field
		MaxOutputTokens   respjson.Field
		Model             respjson.Field
		OutputModalities  respjson.Field
		ParallelToolCalls respjson.Field
		Prompt            respjson.Field
		Reasoning         respjson.Field
		ToolChoice        respjson.Field
		Tools             respjson.Field
		Tracing           respjson.Field
		Truncation        respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeSessionCreateRequest) RawJSON() string { return r.JSON.raw }
func (r *RealtimeSessionCreateRequest) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this RealtimeSessionCreateRequest to a
// RealtimeSessionCreateRequestParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// RealtimeSessionCreateRequestParam.Overrides()
func (r RealtimeSessionCreateRequest) ToParam() RealtimeSessionCreateRequestParam {
	return param.Override[RealtimeSessionCreateRequestParam](json.RawMessage(r.RawJSON()))
}

// RealtimeSessionCreateRequestMaxOutputTokensUnion contains all possible
// properties and values from [int64], [constant.Inf].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfInt OfInf]
type RealtimeSessionCreateRequestMaxOutputTokensUnion struct {
	// This field will be present if the value is a [int64] instead of an object.
	OfInt int64 `json:",inline"`
	// This field will be present if the value is a [constant.Inf] instead of an
	// object.
	OfInf constant.Inf `json:",inline"`
	JSON  struct {
		OfInt respjson.Field
		OfInf respjson.Field
		raw   string
	} `json:"-"`
}

func (u RealtimeSessionCreateRequestMaxOutputTokensUnion) AsInt() (v int64) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeSessionCreateRequestMaxOutputTokensUnion) AsInf() (v constant.Inf) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u RealtimeSessionCreateRequestMaxOutputTokensUnion) RawJSON() string { return u.JSON.raw }

func (r *RealtimeSessionCreateRequestMaxOutputTokensUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The Realtime model used for this session.
type RealtimeSessionCreateRequestModel = string

const (
	// Deprecated: Announced shutdown date: 2027-01-20. See
	// https://developers.openai.com/api/docs/deprecations for details and recommended
	// replacements.
	RealtimeSessionCreateRequestModelGPTRealtime        RealtimeSessionCreateRequestModel = "gpt-realtime"
	RealtimeSessionCreateRequestModelGPTRealtime1_5     RealtimeSessionCreateRequestModel = "gpt-realtime-1.5"
	RealtimeSessionCreateRequestModelGPTRealtime2       RealtimeSessionCreateRequestModel = "gpt-realtime-2"
	RealtimeSessionCreateRequestModelGPTRealtime2_1     RealtimeSessionCreateRequestModel = "gpt-realtime-2.1"
	RealtimeSessionCreateRequestModelGPTRealtime2_1Mini RealtimeSessionCreateRequestModel = "gpt-realtime-2.1-mini"
	// Deprecated: Announced shutdown date: 2027-01-20. See
	// https://developers.openai.com/api/docs/deprecations for details and recommended
	// replacements.
	RealtimeSessionCreateRequestModelGPTRealtime2025_08_28              RealtimeSessionCreateRequestModel = "gpt-realtime-2025-08-28"
	RealtimeSessionCreateRequestModelGPT4oRealtimePreview               RealtimeSessionCreateRequestModel = "gpt-4o-realtime-preview"
	RealtimeSessionCreateRequestModelGPT4oRealtimePreview2024_10_01     RealtimeSessionCreateRequestModel = "gpt-4o-realtime-preview-2024-10-01"
	RealtimeSessionCreateRequestModelGPT4oRealtimePreview2024_12_17     RealtimeSessionCreateRequestModel = "gpt-4o-realtime-preview-2024-12-17"
	RealtimeSessionCreateRequestModelGPT4oRealtimePreview2025_06_03     RealtimeSessionCreateRequestModel = "gpt-4o-realtime-preview-2025-06-03"
	RealtimeSessionCreateRequestModelGPT4oMiniRealtimePreview           RealtimeSessionCreateRequestModel = "gpt-4o-mini-realtime-preview"
	RealtimeSessionCreateRequestModelGPT4oMiniRealtimePreview2024_12_17 RealtimeSessionCreateRequestModel = "gpt-4o-mini-realtime-preview-2024-12-17"
	// Deprecated: Announced shutdown date: 2027-01-20. See
	// https://developers.openai.com/api/docs/deprecations for details and recommended
	// replacements.
	RealtimeSessionCreateRequestModelGPTRealtimeMini           RealtimeSessionCreateRequestModel = "gpt-realtime-mini"
	RealtimeSessionCreateRequestModelGPTRealtimeMini2025_10_06 RealtimeSessionCreateRequestModel = "gpt-realtime-mini-2025-10-06"
	// Deprecated: Announced shutdown date: 2027-01-20. See
	// https://developers.openai.com/api/docs/deprecations for details and recommended
	// replacements.
	RealtimeSessionCreateRequestModelGPTRealtimeMini2025_12_15 RealtimeSessionCreateRequestModel = "gpt-realtime-mini-2025-12-15"
	RealtimeSessionCreateRequestModelGPTAudio1_5               RealtimeSessionCreateRequestModel = "gpt-audio-1.5"
	// Deprecated: Announced shutdown date: 2027-01-20. See
	// https://developers.openai.com/api/docs/deprecations for details and recommended
	// replacements.
	RealtimeSessionCreateRequestModelGPTAudioMini           RealtimeSessionCreateRequestModel = "gpt-audio-mini"
	RealtimeSessionCreateRequestModelGPTAudioMini2025_10_06 RealtimeSessionCreateRequestModel = "gpt-audio-mini-2025-10-06"
	// Deprecated: Announced shutdown date: 2027-01-20. See
	// https://developers.openai.com/api/docs/deprecations for details and recommended
	// replacements.
	RealtimeSessionCreateRequestModelGPTAudioMini2025_12_15 RealtimeSessionCreateRequestModel = "gpt-audio-mini-2025-12-15"
)

// Realtime session object configuration.
//
// The property Type is required.
type RealtimeSessionCreateRequestParam struct {
	// The default system instructions (i.e. system message) prepended to model calls.
	// This field allows the client to guide the model on desired responses. The model
	// can be instructed on response content and format, (e.g. "be extremely succinct",
	// "act friendly", "here are examples of good responses") and on audio behavior
	// (e.g. "talk quickly", "inject emotion into your voice", "laugh frequently"). The
	// instructions are not guaranteed to be followed by the model, but they provide
	// guidance to the model on the desired behavior.
	//
	// Note that the server sets default instructions which will be used if this field
	// is not set and are visible in the `session.created` event at the start of the
	// session.
	Instructions param.Opt[string] `json:"instructions,omitzero"`
	// Whether the model may call multiple tools in parallel. Only supported by
	// reasoning Realtime models such as `gpt-realtime-2`.
	ParallelToolCalls param.Opt[bool] `json:"parallel_tool_calls,omitzero"`
	// Reference to a prompt template and its variables.
	// [Learn more](https://developers.openai.com/api/docs/guides/text?api-mode=responses#version-prompts-in-code).
	Prompt responses.ResponsePromptParam `json:"prompt,omitzero"`
	// Realtime API can write session traces to the
	// [Traces Dashboard](https://platform.openai.com/logs?api=traces). Set to null to
	// disable tracing. Once tracing is enabled for a session, the configuration cannot
	// be modified.
	//
	// `auto` will create a trace for the session with default values for the workflow
	// name, group id, and metadata.
	Tracing RealtimeTracingConfigUnionParam `json:"tracing,omitzero"`
	// Configuration for input and output audio.
	Audio RealtimeAudioConfigParam `json:"audio,omitzero"`
	// Additional fields to include in server outputs.
	//
	// `item.input_audio_transcription.logprobs`: Include logprobs for input audio
	// transcription.
	//
	// Any of "item.input_audio_transcription.logprobs".
	Include []string `json:"include,omitzero"`
	// Maximum number of output tokens for a single assistant response, inclusive of
	// tool calls. Provide an integer between 1 and 4096 to limit output tokens, or
	// `inf` for the maximum available tokens for a given model. Defaults to `inf`.
	MaxOutputTokens RealtimeSessionCreateRequestMaxOutputTokensUnionParam `json:"max_output_tokens,omitzero"`
	// The Realtime model used for this session.
	Model RealtimeSessionCreateRequestModel `json:"model,omitzero"`
	// The set of modalities the model can respond with. It defaults to `["audio"]`,
	// indicating that the model will respond with audio plus a transcript. `["text"]`
	// can be used to make the model respond with text only. It is not possible to
	// request both `text` and `audio` at the same time.
	//
	// Any of "text", "audio".
	OutputModalities []string `json:"output_modalities,omitzero"`
	// Configuration for reasoning-capable Realtime models such as `gpt-realtime-2`.
	Reasoning RealtimeReasoningParam `json:"reasoning,omitzero"`
	// How the model chooses tools. Provide one of the string modes or force a specific
	// function/MCP tool.
	ToolChoice RealtimeToolChoiceConfigUnionParam `json:"tool_choice,omitzero"`
	// Tools available to the model.
	Tools RealtimeToolsConfigParam `json:"tools,omitzero"`
	// When the number of tokens in a conversation exceeds the model's input token
	// limit, the conversation be truncated, meaning messages (starting from the
	// oldest) will not be included in the model's context. A 32k context model with
	// 4,096 max output tokens can only include 28,224 tokens in the context before
	// truncation occurs.
	//
	// Clients can configure truncation behavior to truncate with a lower max token
	// limit, which is an effective way to control token usage and cost.
	//
	// Truncation will reduce the number of cached tokens on the next turn (busting the
	// cache), since messages are dropped from the beginning of the context. However,
	// clients can also configure truncation to retain messages up to a fraction of the
	// maximum context size, which will reduce the need for future truncations and thus
	// improve the cache rate.
	//
	// Truncation can be disabled entirely, which means the server will never truncate
	// but would instead return an error if the conversation exceeds the model's input
	// token limit.
	Truncation RealtimeTruncationUnionParam `json:"truncation,omitzero"`
	// The type of session to create. Always `realtime` for the Realtime API.
	//
	// This field can be elided, and will marshal its zero value as "realtime".
	Type constant.Realtime `json:"type" default:"realtime"`
	paramObj
}

func (r RealtimeSessionCreateRequestParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeSessionCreateRequestParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeSessionCreateRequestParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type RealtimeSessionCreateRequestMaxOutputTokensUnionParam struct {
	OfInt param.Opt[int64] `json:",omitzero,inline"`
	// Construct this variant with constant.ValueOf[constant.Inf]()
	OfInf constant.Inf `json:",omitzero,inline"`
	paramUnion
}

func (u RealtimeSessionCreateRequestMaxOutputTokensUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfInt, u.OfInf)
}
func (u *RealtimeSessionCreateRequestMaxOutputTokensUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// RealtimeToolChoiceConfigUnion contains all possible properties and values from
// [responses.ToolChoiceOptions], [responses.ToolChoiceFunction],
// [responses.ToolChoiceMcp].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfToolChoiceMode]
type RealtimeToolChoiceConfigUnion struct {
	// This field will be present if the value is a [responses.ToolChoiceOptions]
	// instead of an object.
	OfToolChoiceMode responses.ToolChoiceOptions `json:",inline"`
	Name             string                      `json:"name"`
	Type             string                      `json:"type"`
	// This field is from variant [responses.ToolChoiceMcp].
	ServerLabel string `json:"server_label"`
	JSON        struct {
		OfToolChoiceMode respjson.Field
		Name             respjson.Field
		Type             respjson.Field
		ServerLabel      respjson.Field
		raw              string
	} `json:"-"`
}

func (u RealtimeToolChoiceConfigUnion) AsToolChoiceMode() (v responses.ToolChoiceOptions) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeToolChoiceConfigUnion) AsFunctionTool() (v responses.ToolChoiceFunction) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeToolChoiceConfigUnion) AsMcpTool() (v responses.ToolChoiceMcp) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u RealtimeToolChoiceConfigUnion) RawJSON() string { return u.JSON.raw }

func (r *RealtimeToolChoiceConfigUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this RealtimeToolChoiceConfigUnion to a
// RealtimeToolChoiceConfigUnionParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// RealtimeToolChoiceConfigUnionParam.Overrides()
func (r RealtimeToolChoiceConfigUnion) ToParam() RealtimeToolChoiceConfigUnionParam {
	return param.Override[RealtimeToolChoiceConfigUnionParam](json.RawMessage(r.RawJSON()))
}

func RealtimeToolChoiceConfigParamOfFunctionTool(name string) RealtimeToolChoiceConfigUnionParam {
	var variant responses.ToolChoiceFunctionParam
	variant.Name = name
	return RealtimeToolChoiceConfigUnionParam{OfFunctionTool: &variant}
}

func RealtimeToolChoiceConfigParamOfMcpTool(serverLabel string) RealtimeToolChoiceConfigUnionParam {
	var variant responses.ToolChoiceMcpParam
	variant.ServerLabel = serverLabel
	return RealtimeToolChoiceConfigUnionParam{OfMcpTool: &variant}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type RealtimeToolChoiceConfigUnionParam struct {
	// Check if union is this variant with !param.IsOmitted(union.OfToolChoiceMode)
	OfToolChoiceMode param.Opt[responses.ToolChoiceOptions] `json:",omitzero,inline"`
	OfFunctionTool   *responses.ToolChoiceFunctionParam     `json:",omitzero,inline"`
	OfMcpTool        *responses.ToolChoiceMcpParam          `json:",omitzero,inline"`
	paramUnion
}

func (u RealtimeToolChoiceConfigUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfToolChoiceMode, u.OfFunctionTool, u.OfMcpTool)
}
func (u *RealtimeToolChoiceConfigUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeToolChoiceConfigUnionParam) GetServerLabel() *string {
	if vt := u.OfMcpTool; vt != nil {
		return &vt.ServerLabel
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeToolChoiceConfigUnionParam) GetName() *string {
	if vt := u.OfFunctionTool; vt != nil {
		return (*string)(&vt.Name)
	} else if vt := u.OfMcpTool; vt != nil && vt.Name.Valid() {
		return &vt.Name.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeToolChoiceConfigUnionParam) GetType() *string {
	if vt := u.OfFunctionTool; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfMcpTool; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Tools available to the model.
type RealtimeToolsConfig []RealtimeToolsConfigUnion

type RealtimeToolsConfigParam []RealtimeToolsConfigUnionParam

// RealtimeToolsConfigUnion contains all possible properties and values from
// [RealtimeFunctionTool], [RealtimeToolsConfigUnionMcp].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type RealtimeToolsConfigUnion struct {
	// This field is from variant [RealtimeFunctionTool].
	Description string `json:"description"`
	// This field is from variant [RealtimeFunctionTool].
	Name string `json:"name"`
	// This field is from variant [RealtimeFunctionTool].
	Parameters any `json:"parameters"`
	// This field is from variant [RealtimeFunctionTool].
	Type                          RealtimeFunctionToolType    `json:"type"`
	OfRealtimeToolsConfigUnionMcp RealtimeToolsConfigUnionMcp `json:",inline"`
	JSON                          struct {
		Description                   respjson.Field
		Name                          respjson.Field
		Parameters                    respjson.Field
		Type                          respjson.Field
		OfRealtimeToolsConfigUnionMcp respjson.Field
		raw                           string
	} `json:"-"`
}

func (u RealtimeToolsConfigUnion) AsRealtimeFunctionTool() (v RealtimeFunctionTool) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeToolsConfigUnion) AsRealtimeToolsConfigUnionMcp() (v RealtimeToolsConfigUnionMcp) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u RealtimeToolsConfigUnion) RawJSON() string { return u.JSON.raw }

func (r *RealtimeToolsConfigUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this RealtimeToolsConfigUnion to a
// RealtimeToolsConfigUnionParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// RealtimeToolsConfigUnionParam.Overrides()
func (r RealtimeToolsConfigUnion) ToParam() RealtimeToolsConfigUnionParam {
	return param.Override[RealtimeToolsConfigUnionParam](json.RawMessage(r.RawJSON()))
}

// Give the model access to additional tools via remote Model Context Protocol
// (MCP) servers.
// [Learn more about MCP](https://developers.openai.com/api/docs/guides/tools-connectors-mcp).
type RealtimeToolsConfigUnionMcp struct {
	// The tool invocation context(s).
	//
	// Any of "direct", "programmatic".
	AllowedCallers []string `json:"allowed_callers" api:"nullable"`
	// List of allowed tool names or a filter object.
	AllowedTools RealtimeToolsConfigUnionMcpAllowedTools `json:"allowed_tools" api:"nullable"`
	// An OAuth access token that can be used with a remote MCP server, either with a
	// custom MCP server URL or a service connector. Your application must handle the
	// OAuth authorization flow and provide the token here.
	Authorization string `json:"authorization"`
	// Identifier for service connectors, like those available in ChatGPT. One of
	// `server_url`, `connector_id`, or `tunnel_id` must be provided. Learn more about
	// service connectors
	// [here](https://developers.openai.com/api/docs/guides/tools-connectors-mcp#connectors).
	//
	// This field is deprecated for models released after September 1, 2026. Use
	// `server_url` to connect to a remote MCP server, or `tunnel_id` to connect
	// through a Secure MCP Tunnel.
	//
	// Currently supported `connector_id` values are:
	//
	// - Dropbox: `connector_dropbox`
	// - Gmail: `connector_gmail`
	// - Google Calendar: `connector_googlecalendar`
	// - Google Drive: `connector_googledrive`
	// - Microsoft Teams: `connector_microsoftteams`
	// - Outlook Calendar: `connector_outlookcalendar`
	// - Outlook Email: `connector_outlookemail`
	// - SharePoint: `connector_sharepoint`
	//
	// Any of "connector_dropbox", "connector_gmail", "connector_googlecalendar",
	// "connector_googledrive", "connector_microsoftteams",
	// "connector_outlookcalendar", "connector_outlookemail", "connector_sharepoint".
	//
	// Deprecated: deprecated
	ConnectorID string `json:"connector_id"`
	// Whether this MCP tool is deferred and discovered via tool search.
	DeferLoading bool `json:"defer_loading"`
	// Optional HTTP headers to send to the MCP server. Use for authentication or other
	// purposes.
	Headers map[string]string `json:"headers" api:"nullable"`
	// Specify which of the MCP server's tools require approval.
	RequireApproval RealtimeToolsConfigUnionMcpRequireApproval `json:"require_approval" api:"nullable"`
	// Optional description of the MCP server, used to provide more context.
	ServerDescription string `json:"server_description"`
	// A label for this MCP server, used to identify it in tool calls.
	ServerLabel string `json:"server_label" api:"required"`
	// The URL for the MCP server. One of `server_url`, `connector_id`, or `tunnel_id`
	// must be provided.
	ServerURL string `json:"server_url" format:"uri"`
	// The Secure MCP Tunnel ID to use instead of a direct server URL. One of
	// `server_url`, `connector_id`, or `tunnel_id` must be provided.
	TunnelID string `json:"tunnel_id"`
	// The type of the MCP tool. Always `mcp`.
	Type constant.Mcp `json:"type" default:"mcp"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		AllowedCallers    respjson.Field
		AllowedTools      respjson.Field
		Authorization     respjson.Field
		ConnectorID       respjson.Field
		DeferLoading      respjson.Field
		Headers           respjson.Field
		RequireApproval   respjson.Field
		ServerDescription respjson.Field
		ServerLabel       respjson.Field
		ServerURL         respjson.Field
		TunnelID          respjson.Field
		Type              respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeToolsConfigUnionMcp) RawJSON() string { return r.JSON.raw }
func (r *RealtimeToolsConfigUnionMcp) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// RealtimeToolsConfigUnionMcpAllowedTools contains all possible properties and
// values from [[]string], [RealtimeToolsConfigUnionMcpAllowedToolsMcpToolFilter].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfMcpAllowedTools]
type RealtimeToolsConfigUnionMcpAllowedTools struct {
	// This field will be present if the value is a [[]string] instead of an object.
	OfMcpAllowedTools                                      []string                                             `json:",inline"`
	OfRealtimeToolsConfigUnionMcpAllowedToolsMcpToolFilter RealtimeToolsConfigUnionMcpAllowedToolsMcpToolFilter `json:",inline"`
	JSON                                                   struct {
		OfMcpAllowedTools                                      respjson.Field
		OfRealtimeToolsConfigUnionMcpAllowedToolsMcpToolFilter respjson.Field
		raw                                                    string
	} `json:"-"`
}

func (u RealtimeToolsConfigUnionMcpAllowedTools) AsMcpAllowedTools() (v []string) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeToolsConfigUnionMcpAllowedTools) AsRealtimeToolsConfigUnionMcpAllowedToolsMcpToolFilter() (v RealtimeToolsConfigUnionMcpAllowedToolsMcpToolFilter) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u RealtimeToolsConfigUnionMcpAllowedTools) RawJSON() string { return u.JSON.raw }

func (r *RealtimeToolsConfigUnionMcpAllowedTools) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A filter object to specify which tools are allowed.
type RealtimeToolsConfigUnionMcpAllowedToolsMcpToolFilter struct {
	// Indicates whether or not a tool modifies data or is read-only. If an MCP server
	// is
	// [annotated with `readOnlyHint`](https://modelcontextprotocol.io/specification/2025-06-18/schema#toolannotations-readonlyhint),
	// it will match this filter.
	ReadOnly bool `json:"read_only"`
	// List of allowed tool names.
	ToolNames []string `json:"tool_names"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ReadOnly    respjson.Field
		ToolNames   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeToolsConfigUnionMcpAllowedToolsMcpToolFilter) RawJSON() string { return r.JSON.raw }
func (r *RealtimeToolsConfigUnionMcpAllowedToolsMcpToolFilter) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// RealtimeToolsConfigUnionMcpRequireApproval contains all possible properties and
// values from [RealtimeToolsConfigUnionMcpRequireApprovalMcpToolApprovalFilter],
// [string].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfString]
type RealtimeToolsConfigUnionMcpRequireApproval struct {
	// This field will be present if the value is a [string] instead of an object.
	OfString                                                          string                                                          `json:",inline"`
	OfRealtimeToolsConfigUnionMcpRequireApprovalMcpToolApprovalFilter RealtimeToolsConfigUnionMcpRequireApprovalMcpToolApprovalFilter `json:",inline"`
	JSON                                                              struct {
		OfString                                                          respjson.Field
		OfRealtimeToolsConfigUnionMcpRequireApprovalMcpToolApprovalFilter respjson.Field
		raw                                                               string
	} `json:"-"`
}

func (u RealtimeToolsConfigUnionMcpRequireApproval) AsRealtimeToolsConfigUnionMcpRequireApprovalMcpToolApprovalFilter() (v RealtimeToolsConfigUnionMcpRequireApprovalMcpToolApprovalFilter) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeToolsConfigUnionMcpRequireApproval) AsString() (v string) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u RealtimeToolsConfigUnionMcpRequireApproval) RawJSON() string { return u.JSON.raw }

func (r *RealtimeToolsConfigUnionMcpRequireApproval) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Specify which of the MCP server's tools require approval. Can be `always`,
// `never`, or a filter object associated with tools that require approval.
type RealtimeToolsConfigUnionMcpRequireApprovalMcpToolApprovalFilter struct {
	// A filter object to specify which tools are allowed.
	Always RealtimeToolsConfigUnionMcpRequireApprovalMcpToolApprovalFilterAlways `json:"always"`
	// A filter object to specify which tools are allowed.
	Never RealtimeToolsConfigUnionMcpRequireApprovalMcpToolApprovalFilterNever `json:"never"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Always      respjson.Field
		Never       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeToolsConfigUnionMcpRequireApprovalMcpToolApprovalFilter) RawJSON() string {
	return r.JSON.raw
}
func (r *RealtimeToolsConfigUnionMcpRequireApprovalMcpToolApprovalFilter) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A filter object to specify which tools are allowed.
type RealtimeToolsConfigUnionMcpRequireApprovalMcpToolApprovalFilterAlways struct {
	// Indicates whether or not a tool modifies data or is read-only. If an MCP server
	// is
	// [annotated with `readOnlyHint`](https://modelcontextprotocol.io/specification/2025-06-18/schema#toolannotations-readonlyhint),
	// it will match this filter.
	ReadOnly bool `json:"read_only"`
	// List of allowed tool names.
	ToolNames []string `json:"tool_names"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ReadOnly    respjson.Field
		ToolNames   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeToolsConfigUnionMcpRequireApprovalMcpToolApprovalFilterAlways) RawJSON() string {
	return r.JSON.raw
}
func (r *RealtimeToolsConfigUnionMcpRequireApprovalMcpToolApprovalFilterAlways) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A filter object to specify which tools are allowed.
type RealtimeToolsConfigUnionMcpRequireApprovalMcpToolApprovalFilterNever struct {
	// Indicates whether or not a tool modifies data or is read-only. If an MCP server
	// is
	// [annotated with `readOnlyHint`](https://modelcontextprotocol.io/specification/2025-06-18/schema#toolannotations-readonlyhint),
	// it will match this filter.
	ReadOnly bool `json:"read_only"`
	// List of allowed tool names.
	ToolNames []string `json:"tool_names"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ReadOnly    respjson.Field
		ToolNames   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeToolsConfigUnionMcpRequireApprovalMcpToolApprovalFilterNever) RawJSON() string {
	return r.JSON.raw
}
func (r *RealtimeToolsConfigUnionMcpRequireApprovalMcpToolApprovalFilterNever) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Specify a single approval policy for all tools. One of `always` or `never`. When
// set to `always`, all tools will require approval. When set to `never`, all tools
// will not require approval.
type RealtimeToolsConfigUnionMcpRequireApprovalMcpToolApprovalSetting string

const (
	RealtimeToolsConfigUnionMcpRequireApprovalMcpToolApprovalSettingAlways RealtimeToolsConfigUnionMcpRequireApprovalMcpToolApprovalSetting = "always"
	RealtimeToolsConfigUnionMcpRequireApprovalMcpToolApprovalSettingNever  RealtimeToolsConfigUnionMcpRequireApprovalMcpToolApprovalSetting = "never"
)

func RealtimeToolsConfigUnionParamOfMcp(serverLabel string) RealtimeToolsConfigUnionParam {
	var mcp RealtimeToolsConfigUnionMcpParam
	mcp.ServerLabel = serverLabel
	return RealtimeToolsConfigUnionParam{OfMcp: &mcp}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type RealtimeToolsConfigUnionParam struct {
	OfFunction *RealtimeFunctionToolParam        `json:",omitzero,inline"`
	OfMcp      *RealtimeToolsConfigUnionMcpParam `json:",omitzero,inline"`
	paramUnion
}

func (u RealtimeToolsConfigUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfFunction, u.OfMcp)
}
func (u *RealtimeToolsConfigUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeToolsConfigUnionParam) GetDescription() *string {
	if vt := u.OfFunction; vt != nil && vt.Description.Valid() {
		return &vt.Description.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeToolsConfigUnionParam) GetName() *string {
	if vt := u.OfFunction; vt != nil && vt.Name.Valid() {
		return &vt.Name.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeToolsConfigUnionParam) GetParameters() *any {
	if vt := u.OfFunction; vt != nil {
		return &vt.Parameters
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeToolsConfigUnionParam) GetServerLabel() *string {
	if vt := u.OfMcp; vt != nil {
		return &vt.ServerLabel
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeToolsConfigUnionParam) GetAllowedCallers() []string {
	if vt := u.OfMcp; vt != nil {
		return vt.AllowedCallers
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeToolsConfigUnionParam) GetAllowedTools() *RealtimeToolsConfigUnionMcpAllowedToolsParam {
	if vt := u.OfMcp; vt != nil {
		return &vt.AllowedTools
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeToolsConfigUnionParam) GetAuthorization() *string {
	if vt := u.OfMcp; vt != nil && vt.Authorization.Valid() {
		return &vt.Authorization.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeToolsConfigUnionParam) GetConnectorID() *string {
	if vt := u.OfMcp; vt != nil {
		return &vt.ConnectorID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeToolsConfigUnionParam) GetDeferLoading() *bool {
	if vt := u.OfMcp; vt != nil && vt.DeferLoading.Valid() {
		return &vt.DeferLoading.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeToolsConfigUnionParam) GetHeaders() map[string]string {
	if vt := u.OfMcp; vt != nil {
		return vt.Headers
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeToolsConfigUnionParam) GetRequireApproval() *RealtimeToolsConfigUnionMcpRequireApprovalParam {
	if vt := u.OfMcp; vt != nil {
		return &vt.RequireApproval
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeToolsConfigUnionParam) GetServerDescription() *string {
	if vt := u.OfMcp; vt != nil && vt.ServerDescription.Valid() {
		return &vt.ServerDescription.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeToolsConfigUnionParam) GetServerURL() *string {
	if vt := u.OfMcp; vt != nil && vt.ServerURL.Valid() {
		return &vt.ServerURL.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeToolsConfigUnionParam) GetTunnelID() *string {
	if vt := u.OfMcp; vt != nil && vt.TunnelID.Valid() {
		return &vt.TunnelID.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeToolsConfigUnionParam) GetType() *string {
	if vt := u.OfFunction; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfMcp; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[RealtimeToolsConfigUnionParam](
		"type",
		apijson.Discriminator[RealtimeFunctionToolParam]("function"),
		apijson.Discriminator[RealtimeToolsConfigUnionMcpParam]("mcp"),
	)
}

// Give the model access to additional tools via remote Model Context Protocol
// (MCP) servers.
// [Learn more about MCP](https://developers.openai.com/api/docs/guides/tools-connectors-mcp).
//
// The properties ServerLabel, Type are required.
type RealtimeToolsConfigUnionMcpParam struct {
	// A label for this MCP server, used to identify it in tool calls.
	ServerLabel string `json:"server_label" api:"required"`
	// An OAuth access token that can be used with a remote MCP server, either with a
	// custom MCP server URL or a service connector. Your application must handle the
	// OAuth authorization flow and provide the token here.
	Authorization param.Opt[string] `json:"authorization,omitzero"`
	// Whether this MCP tool is deferred and discovered via tool search.
	DeferLoading param.Opt[bool] `json:"defer_loading,omitzero"`
	// Optional description of the MCP server, used to provide more context.
	ServerDescription param.Opt[string] `json:"server_description,omitzero"`
	// The URL for the MCP server. One of `server_url`, `connector_id`, or `tunnel_id`
	// must be provided.
	ServerURL param.Opt[string] `json:"server_url,omitzero" format:"uri"`
	// The Secure MCP Tunnel ID to use instead of a direct server URL. One of
	// `server_url`, `connector_id`, or `tunnel_id` must be provided.
	TunnelID param.Opt[string] `json:"tunnel_id,omitzero"`
	// The tool invocation context(s).
	//
	// Any of "direct", "programmatic".
	AllowedCallers []string `json:"allowed_callers,omitzero"`
	// List of allowed tool names or a filter object.
	AllowedTools RealtimeToolsConfigUnionMcpAllowedToolsParam `json:"allowed_tools,omitzero"`
	// Optional HTTP headers to send to the MCP server. Use for authentication or other
	// purposes.
	Headers map[string]string `json:"headers,omitzero"`
	// Specify which of the MCP server's tools require approval.
	RequireApproval RealtimeToolsConfigUnionMcpRequireApprovalParam `json:"require_approval,omitzero"`
	// Identifier for service connectors, like those available in ChatGPT. One of
	// `server_url`, `connector_id`, or `tunnel_id` must be provided. Learn more about
	// service connectors
	// [here](https://developers.openai.com/api/docs/guides/tools-connectors-mcp#connectors).
	//
	// This field is deprecated for models released after September 1, 2026. Use
	// `server_url` to connect to a remote MCP server, or `tunnel_id` to connect
	// through a Secure MCP Tunnel.
	//
	// Currently supported `connector_id` values are:
	//
	// - Dropbox: `connector_dropbox`
	// - Gmail: `connector_gmail`
	// - Google Calendar: `connector_googlecalendar`
	// - Google Drive: `connector_googledrive`
	// - Microsoft Teams: `connector_microsoftteams`
	// - Outlook Calendar: `connector_outlookcalendar`
	// - Outlook Email: `connector_outlookemail`
	// - SharePoint: `connector_sharepoint`
	//
	// Any of "connector_dropbox", "connector_gmail", "connector_googlecalendar",
	// "connector_googledrive", "connector_microsoftteams",
	// "connector_outlookcalendar", "connector_outlookemail", "connector_sharepoint".
	//
	// Deprecated: deprecated
	ConnectorID string `json:"connector_id,omitzero"`
	// The type of the MCP tool. Always `mcp`.
	//
	// This field can be elided, and will marshal its zero value as "mcp".
	Type constant.Mcp `json:"type" default:"mcp"`
	paramObj
}

func (r RealtimeToolsConfigUnionMcpParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeToolsConfigUnionMcpParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeToolsConfigUnionMcpParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[RealtimeToolsConfigUnionMcpParam](
		"connector_id", "connector_dropbox", "connector_gmail", "connector_googlecalendar", "connector_googledrive", "connector_microsoftteams", "connector_outlookcalendar", "connector_outlookemail", "connector_sharepoint",
	)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type RealtimeToolsConfigUnionMcpAllowedToolsParam struct {
	OfMcpAllowedTools []string                                                   `json:",omitzero,inline"`
	OfMcpToolFilter   *RealtimeToolsConfigUnionMcpAllowedToolsMcpToolFilterParam `json:",omitzero,inline"`
	paramUnion
}

func (u RealtimeToolsConfigUnionMcpAllowedToolsParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfMcpAllowedTools, u.OfMcpToolFilter)
}
func (u *RealtimeToolsConfigUnionMcpAllowedToolsParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// A filter object to specify which tools are allowed.
type RealtimeToolsConfigUnionMcpAllowedToolsMcpToolFilterParam struct {
	// Indicates whether or not a tool modifies data or is read-only. If an MCP server
	// is
	// [annotated with `readOnlyHint`](https://modelcontextprotocol.io/specification/2025-06-18/schema#toolannotations-readonlyhint),
	// it will match this filter.
	ReadOnly param.Opt[bool] `json:"read_only,omitzero"`
	// List of allowed tool names.
	ToolNames []string `json:"tool_names,omitzero"`
	paramObj
}

func (r RealtimeToolsConfigUnionMcpAllowedToolsMcpToolFilterParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeToolsConfigUnionMcpAllowedToolsMcpToolFilterParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeToolsConfigUnionMcpAllowedToolsMcpToolFilterParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type RealtimeToolsConfigUnionMcpRequireApprovalParam struct {
	OfMcpToolApprovalFilter *RealtimeToolsConfigUnionMcpRequireApprovalMcpToolApprovalFilterParam `json:",omitzero,inline"`
	// Check if union is this variant with
	// !param.IsOmitted(union.OfMcpToolApprovalSetting)
	OfMcpToolApprovalSetting param.Opt[string] `json:",omitzero,inline"`
	paramUnion
}

func (u RealtimeToolsConfigUnionMcpRequireApprovalParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfMcpToolApprovalFilter, u.OfMcpToolApprovalSetting)
}
func (u *RealtimeToolsConfigUnionMcpRequireApprovalParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Specify which of the MCP server's tools require approval. Can be `always`,
// `never`, or a filter object associated with tools that require approval.
type RealtimeToolsConfigUnionMcpRequireApprovalMcpToolApprovalFilterParam struct {
	// A filter object to specify which tools are allowed.
	Always RealtimeToolsConfigUnionMcpRequireApprovalMcpToolApprovalFilterAlwaysParam `json:"always,omitzero"`
	// A filter object to specify which tools are allowed.
	Never RealtimeToolsConfigUnionMcpRequireApprovalMcpToolApprovalFilterNeverParam `json:"never,omitzero"`
	paramObj
}

func (r RealtimeToolsConfigUnionMcpRequireApprovalMcpToolApprovalFilterParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeToolsConfigUnionMcpRequireApprovalMcpToolApprovalFilterParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeToolsConfigUnionMcpRequireApprovalMcpToolApprovalFilterParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A filter object to specify which tools are allowed.
type RealtimeToolsConfigUnionMcpRequireApprovalMcpToolApprovalFilterAlwaysParam struct {
	// Indicates whether or not a tool modifies data or is read-only. If an MCP server
	// is
	// [annotated with `readOnlyHint`](https://modelcontextprotocol.io/specification/2025-06-18/schema#toolannotations-readonlyhint),
	// it will match this filter.
	ReadOnly param.Opt[bool] `json:"read_only,omitzero"`
	// List of allowed tool names.
	ToolNames []string `json:"tool_names,omitzero"`
	paramObj
}

func (r RealtimeToolsConfigUnionMcpRequireApprovalMcpToolApprovalFilterAlwaysParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeToolsConfigUnionMcpRequireApprovalMcpToolApprovalFilterAlwaysParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeToolsConfigUnionMcpRequireApprovalMcpToolApprovalFilterAlwaysParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A filter object to specify which tools are allowed.
type RealtimeToolsConfigUnionMcpRequireApprovalMcpToolApprovalFilterNeverParam struct {
	// Indicates whether or not a tool modifies data or is read-only. If an MCP server
	// is
	// [annotated with `readOnlyHint`](https://modelcontextprotocol.io/specification/2025-06-18/schema#toolannotations-readonlyhint),
	// it will match this filter.
	ReadOnly param.Opt[bool] `json:"read_only,omitzero"`
	// List of allowed tool names.
	ToolNames []string `json:"tool_names,omitzero"`
	paramObj
}

func (r RealtimeToolsConfigUnionMcpRequireApprovalMcpToolApprovalFilterNeverParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeToolsConfigUnionMcpRequireApprovalMcpToolApprovalFilterNeverParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeToolsConfigUnionMcpRequireApprovalMcpToolApprovalFilterNeverParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// RealtimeTracingConfigUnion contains all possible properties and values from
// [constant.Auto], [RealtimeTracingConfigTracingConfiguration].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfAuto]
type RealtimeTracingConfigUnion struct {
	// This field will be present if the value is a [constant.Auto] instead of an
	// object.
	OfAuto                                      constant.Auto                             `json:",inline"`
	OfRealtimeTracingConfigTracingConfiguration RealtimeTracingConfigTracingConfiguration `json:",inline"`
	JSON                                        struct {
		OfAuto                                      respjson.Field
		OfRealtimeTracingConfigTracingConfiguration respjson.Field
		raw                                         string
	} `json:"-"`
}

func (u RealtimeTracingConfigUnion) AsAuto() (v constant.Auto) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeTracingConfigUnion) AsRealtimeTracingConfigTracingConfiguration() (v RealtimeTracingConfigTracingConfiguration) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u RealtimeTracingConfigUnion) RawJSON() string { return u.JSON.raw }

func (r *RealtimeTracingConfigUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this RealtimeTracingConfigUnion to a
// RealtimeTracingConfigUnionParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// RealtimeTracingConfigUnionParam.Overrides()
func (r RealtimeTracingConfigUnion) ToParam() RealtimeTracingConfigUnionParam {
	return param.Override[RealtimeTracingConfigUnionParam](json.RawMessage(r.RawJSON()))
}

// Granular configuration for tracing.
type RealtimeTracingConfigTracingConfiguration struct {
	// The group id to attach to this trace to enable filtering and grouping in the
	// Traces Dashboard.
	GroupID string `json:"group_id"`
	// The arbitrary metadata to attach to this trace to enable filtering in the Traces
	// Dashboard.
	Metadata any `json:"metadata"`
	// The name of the workflow to attach to this trace. This is used to name the trace
	// in the Traces Dashboard.
	WorkflowName string `json:"workflow_name"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		GroupID      respjson.Field
		Metadata     respjson.Field
		WorkflowName respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeTracingConfigTracingConfiguration) RawJSON() string { return r.JSON.raw }
func (r *RealtimeTracingConfigTracingConfiguration) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func RealtimeTracingConfigParamOfAuto() RealtimeTracingConfigUnionParam {
	return RealtimeTracingConfigUnionParam{OfAuto: constant.ValueOf[constant.Auto]()}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type RealtimeTracingConfigUnionParam struct {
	// Construct this variant with constant.ValueOf[constant.Auto]()
	OfAuto                 constant.Auto                                   `json:",omitzero,inline"`
	OfTracingConfiguration *RealtimeTracingConfigTracingConfigurationParam `json:",omitzero,inline"`
	paramUnion
}

func (u RealtimeTracingConfigUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfAuto, u.OfTracingConfiguration)
}
func (u *RealtimeTracingConfigUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Granular configuration for tracing.
type RealtimeTracingConfigTracingConfigurationParam struct {
	// The group id to attach to this trace to enable filtering and grouping in the
	// Traces Dashboard.
	GroupID param.Opt[string] `json:"group_id,omitzero"`
	// The name of the workflow to attach to this trace. This is used to name the trace
	// in the Traces Dashboard.
	WorkflowName param.Opt[string] `json:"workflow_name,omitzero"`
	// The arbitrary metadata to attach to this trace to enable filtering in the Traces
	// Dashboard.
	Metadata any `json:"metadata,omitzero"`
	paramObj
}

func (r RealtimeTracingConfigTracingConfigurationParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeTracingConfigTracingConfigurationParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeTracingConfigTracingConfigurationParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for input and output audio.
type RealtimeTranscriptionSessionAudio struct {
	Input RealtimeTranscriptionSessionAudioInput `json:"input"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Input       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeTranscriptionSessionAudio) RawJSON() string { return r.JSON.raw }
func (r *RealtimeTranscriptionSessionAudio) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this RealtimeTranscriptionSessionAudio to a
// RealtimeTranscriptionSessionAudioParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// RealtimeTranscriptionSessionAudioParam.Overrides()
func (r RealtimeTranscriptionSessionAudio) ToParam() RealtimeTranscriptionSessionAudioParam {
	return param.Override[RealtimeTranscriptionSessionAudioParam](json.RawMessage(r.RawJSON()))
}

// Configuration for input and output audio.
type RealtimeTranscriptionSessionAudioParam struct {
	Input RealtimeTranscriptionSessionAudioInputParam `json:"input,omitzero"`
	paramObj
}

func (r RealtimeTranscriptionSessionAudioParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeTranscriptionSessionAudioParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeTranscriptionSessionAudioParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type RealtimeTranscriptionSessionAudioInput struct {
	// The PCM audio format. Only a 24kHz sample rate is supported.
	Format RealtimeAudioFormatsUnion `json:"format"`
	// Configuration for input audio noise reduction. This can be set to `null` to turn
	// off. Noise reduction filters audio added to the input audio buffer before it is
	// sent to VAD and the model. Filtering the audio can improve VAD and turn
	// detection accuracy (reducing false positives) and model performance by improving
	// perception of the input audio.
	NoiseReduction RealtimeTranscriptionSessionAudioInputNoiseReduction `json:"noise_reduction"`
	// Configuration for input audio transcription, defaults to off and can be set to
	// `null` to turn off once on. Input audio transcription is not native to the
	// model, since the model consumes audio directly. Transcription runs
	// asynchronously through
	// [the /audio/transcriptions endpoint](https://developers.openai.com/api/reference/resources/audio/subresources/transcriptions/methods/create)
	// and should be treated as guidance of input audio content rather than precisely
	// what the model heard. The client can optionally set the language and prompt for
	// transcription, these offer additional guidance to the transcription service.
	Transcription AudioTranscription `json:"transcription"`
	// Configuration for turn detection, ether Server VAD or Semantic VAD. This can be
	// set to `null` to turn off, in which case the client must manually trigger model
	// response.
	//
	// Server VAD means that the model will detect the start and end of speech based on
	// audio volume and respond at the end of user speech.
	//
	// Semantic VAD is more advanced and uses a turn detection model (in conjunction
	// with VAD) to semantically estimate whether the user has finished speaking, then
	// dynamically sets a timeout based on this probability. For example, if user audio
	// trails off with "uhhm", the model will score a low probability of turn end and
	// wait longer for the user to continue speaking. This can be useful for more
	// natural conversations, but may have a higher latency.
	//
	// For `gpt-realtime-whisper` transcription sessions, turn detection must be set to
	// `null`; VAD is not supported.
	TurnDetection RealtimeTranscriptionSessionAudioInputTurnDetectionUnion `json:"turn_detection" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Format         respjson.Field
		NoiseReduction respjson.Field
		Transcription  respjson.Field
		TurnDetection  respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeTranscriptionSessionAudioInput) RawJSON() string { return r.JSON.raw }
func (r *RealtimeTranscriptionSessionAudioInput) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this RealtimeTranscriptionSessionAudioInput to a
// RealtimeTranscriptionSessionAudioInputParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// RealtimeTranscriptionSessionAudioInputParam.Overrides()
func (r RealtimeTranscriptionSessionAudioInput) ToParam() RealtimeTranscriptionSessionAudioInputParam {
	return param.Override[RealtimeTranscriptionSessionAudioInputParam](json.RawMessage(r.RawJSON()))
}

// Configuration for input audio noise reduction. This can be set to `null` to turn
// off. Noise reduction filters audio added to the input audio buffer before it is
// sent to VAD and the model. Filtering the audio can improve VAD and turn
// detection accuracy (reducing false positives) and model performance by improving
// perception of the input audio.
type RealtimeTranscriptionSessionAudioInputNoiseReduction struct {
	// Type of noise reduction. `near_field` is for close-talking microphones such as
	// headphones, `far_field` is for far-field microphones such as laptop or
	// conference room microphones.
	//
	// Any of "near_field", "far_field".
	Type NoiseReductionType `json:"type"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeTranscriptionSessionAudioInputNoiseReduction) RawJSON() string { return r.JSON.raw }
func (r *RealtimeTranscriptionSessionAudioInputNoiseReduction) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type RealtimeTranscriptionSessionAudioInputParam struct {
	// Configuration for turn detection, ether Server VAD or Semantic VAD. This can be
	// set to `null` to turn off, in which case the client must manually trigger model
	// response.
	//
	// Server VAD means that the model will detect the start and end of speech based on
	// audio volume and respond at the end of user speech.
	//
	// Semantic VAD is more advanced and uses a turn detection model (in conjunction
	// with VAD) to semantically estimate whether the user has finished speaking, then
	// dynamically sets a timeout based on this probability. For example, if user audio
	// trails off with "uhhm", the model will score a low probability of turn end and
	// wait longer for the user to continue speaking. This can be useful for more
	// natural conversations, but may have a higher latency.
	//
	// For `gpt-realtime-whisper` transcription sessions, turn detection must be set to
	// `null`; VAD is not supported.
	TurnDetection RealtimeTranscriptionSessionAudioInputTurnDetectionUnionParam `json:"turn_detection,omitzero"`
	// The PCM audio format. Only a 24kHz sample rate is supported.
	Format RealtimeAudioFormatsUnionParam `json:"format,omitzero"`
	// Configuration for input audio noise reduction. This can be set to `null` to turn
	// off. Noise reduction filters audio added to the input audio buffer before it is
	// sent to VAD and the model. Filtering the audio can improve VAD and turn
	// detection accuracy (reducing false positives) and model performance by improving
	// perception of the input audio.
	NoiseReduction RealtimeTranscriptionSessionAudioInputNoiseReductionParam `json:"noise_reduction,omitzero"`
	// Configuration for input audio transcription, defaults to off and can be set to
	// `null` to turn off once on. Input audio transcription is not native to the
	// model, since the model consumes audio directly. Transcription runs
	// asynchronously through
	// [the /audio/transcriptions endpoint](https://developers.openai.com/api/reference/resources/audio/subresources/transcriptions/methods/create)
	// and should be treated as guidance of input audio content rather than precisely
	// what the model heard. The client can optionally set the language and prompt for
	// transcription, these offer additional guidance to the transcription service.
	Transcription AudioTranscriptionParam `json:"transcription,omitzero"`
	paramObj
}

func (r RealtimeTranscriptionSessionAudioInputParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeTranscriptionSessionAudioInputParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeTranscriptionSessionAudioInputParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for input audio noise reduction. This can be set to `null` to turn
// off. Noise reduction filters audio added to the input audio buffer before it is
// sent to VAD and the model. Filtering the audio can improve VAD and turn
// detection accuracy (reducing false positives) and model performance by improving
// perception of the input audio.
type RealtimeTranscriptionSessionAudioInputNoiseReductionParam struct {
	// Type of noise reduction. `near_field` is for close-talking microphones such as
	// headphones, `far_field` is for far-field microphones such as laptop or
	// conference room microphones.
	//
	// Any of "near_field", "far_field".
	Type NoiseReductionType `json:"type,omitzero"`
	paramObj
}

func (r RealtimeTranscriptionSessionAudioInputNoiseReductionParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeTranscriptionSessionAudioInputNoiseReductionParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeTranscriptionSessionAudioInputNoiseReductionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// RealtimeTranscriptionSessionAudioInputTurnDetectionUnion contains all possible
// properties and values from
// [RealtimeTranscriptionSessionAudioInputTurnDetectionServerVad],
// [RealtimeTranscriptionSessionAudioInputTurnDetectionSemanticVad].
//
// Use the [RealtimeTranscriptionSessionAudioInputTurnDetectionUnion.AsAny] method
// to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type RealtimeTranscriptionSessionAudioInputTurnDetectionUnion struct {
	OfRealtimeTranscriptionSessionAudioInputTurnDetectionServerVad   RealtimeTranscriptionSessionAudioInputTurnDetectionServerVad   `json:",inline"`
	OfRealtimeTranscriptionSessionAudioInputTurnDetectionSemanticVad RealtimeTranscriptionSessionAudioInputTurnDetectionSemanticVad `json:",inline"`
	Type                                                             string                                                         `json:"type"`
	JSON                                                             struct {
		OfRealtimeTranscriptionSessionAudioInputTurnDetectionServerVad   respjson.Field
		OfRealtimeTranscriptionSessionAudioInputTurnDetectionSemanticVad respjson.Field
		Type                                                             respjson.Field
		raw                                                              string
	} `json:"-"`
}

// anyRealtimeTranscriptionSessionAudioInputTurnDetection is implemented by each
// variant of [RealtimeTranscriptionSessionAudioInputTurnDetectionUnion] to add
// type safety for the return type of
// [RealtimeTranscriptionSessionAudioInputTurnDetectionUnion.AsAny]
type anyRealtimeTranscriptionSessionAudioInputTurnDetection interface {
	implRealtimeTranscriptionSessionAudioInputTurnDetectionUnion()
}

func (RealtimeTranscriptionSessionAudioInputTurnDetectionServerVad) implRealtimeTranscriptionSessionAudioInputTurnDetectionUnion() {
}
func (RealtimeTranscriptionSessionAudioInputTurnDetectionSemanticVad) implRealtimeTranscriptionSessionAudioInputTurnDetectionUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := RealtimeTranscriptionSessionAudioInputTurnDetectionUnion.AsAny().(type) {
//	case realtime.RealtimeTranscriptionSessionAudioInputTurnDetectionServerVad:
//	case realtime.RealtimeTranscriptionSessionAudioInputTurnDetectionSemanticVad:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u RealtimeTranscriptionSessionAudioInputTurnDetectionUnion) AsAny() anyRealtimeTranscriptionSessionAudioInputTurnDetection {
	switch u.Type {
	case "server_vad":
		return u.AsRealtimeTranscriptionSessionAudioInputTurnDetectionServerVad()
	case "semantic_vad":
		return u.AsRealtimeTranscriptionSessionAudioInputTurnDetectionSemanticVad()
	}
	return nil
}

func (u RealtimeTranscriptionSessionAudioInputTurnDetectionUnion) AsRealtimeTranscriptionSessionAudioInputTurnDetectionServerVad() (v RealtimeTranscriptionSessionAudioInputTurnDetectionServerVad) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeTranscriptionSessionAudioInputTurnDetectionUnion) AsRealtimeTranscriptionSessionAudioInputTurnDetectionSemanticVad() (v RealtimeTranscriptionSessionAudioInputTurnDetectionSemanticVad) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u RealtimeTranscriptionSessionAudioInputTurnDetectionUnion) RawJSON() string { return u.JSON.raw }

func (r *RealtimeTranscriptionSessionAudioInputTurnDetectionUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this RealtimeTranscriptionSessionAudioInputTurnDetectionUnion
// to a RealtimeTranscriptionSessionAudioInputTurnDetectionUnionParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// RealtimeTranscriptionSessionAudioInputTurnDetectionUnionParam.Overrides()
func (r RealtimeTranscriptionSessionAudioInputTurnDetectionUnion) ToParam() RealtimeTranscriptionSessionAudioInputTurnDetectionUnionParam {
	return param.Override[RealtimeTranscriptionSessionAudioInputTurnDetectionUnionParam](json.RawMessage(r.RawJSON()))
}

// Server-side voice activity detection (VAD) which flips on when user speech is
// detected and off after a period of silence.
type RealtimeTranscriptionSessionAudioInputTurnDetectionServerVad struct {
	// Whether or not to automatically generate a response when a VAD stop event
	// occurs. If `interrupt_response` is set to `false` this may fail to create a
	// response if the model is already responding.
	//
	// If both `create_response` and `interrupt_response` are set to `false`, the model
	// will never respond automatically but VAD events will still be emitted.
	CreateResponse bool `json:"create_response"`
	// Optional timeout after which a model response will be triggered automatically.
	// This is useful for situations in which a long pause from the user is unexpected,
	// such as a phone call. The model will effectively prompt the user to continue the
	// conversation based on the current context.
	//
	// The timeout value will be applied after the last model response's audio has
	// finished playing, i.e. it's set to the `response.done` time plus audio playback
	// duration.
	//
	// An `input_audio_buffer.timeout_triggered` event (plus events associated with the
	// Response) will be emitted when the timeout is reached. Idle timeout is currently
	// only supported for `server_vad` mode.
	IdleTimeoutMs int64 `json:"idle_timeout_ms" api:"nullable"`
	// Whether or not to automatically interrupt (cancel) any ongoing response with
	// output to the default conversation (i.e. `conversation` of `auto`) when a VAD
	// start event occurs. If `true` then the response will be cancelled, otherwise it
	// will continue until complete.
	//
	// If both `create_response` and `interrupt_response` are set to `false`, the model
	// will never respond automatically but VAD events will still be emitted.
	InterruptResponse bool `json:"interrupt_response"`
	// Used only for `server_vad` mode. Amount of audio to include before the VAD
	// detected speech (in milliseconds). Defaults to 300ms.
	PrefixPaddingMs int64 `json:"prefix_padding_ms"`
	// Used only for `server_vad` mode. Duration of silence to detect speech stop (in
	// milliseconds). Defaults to 500ms. With shorter values the model will respond
	// more quickly, but may jump in on short pauses from the user.
	SilenceDurationMs int64 `json:"silence_duration_ms"`
	// Used only for `server_vad` mode. Activation threshold for VAD (0.0 to 1.0), this
	// defaults to 0.5. A higher threshold will require louder audio to activate the
	// model, and thus might perform better in noisy environments.
	Threshold float64 `json:"threshold"`
	// Type of turn detection, `server_vad` to turn on simple Server VAD.
	Type constant.ServerVad `json:"type" default:"server_vad"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		CreateResponse    respjson.Field
		IdleTimeoutMs     respjson.Field
		InterruptResponse respjson.Field
		PrefixPaddingMs   respjson.Field
		SilenceDurationMs respjson.Field
		Threshold         respjson.Field
		Type              respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeTranscriptionSessionAudioInputTurnDetectionServerVad) RawJSON() string {
	return r.JSON.raw
}
func (r *RealtimeTranscriptionSessionAudioInputTurnDetectionServerVad) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Server-side semantic turn detection which uses a model to determine when the
// user has finished speaking.
type RealtimeTranscriptionSessionAudioInputTurnDetectionSemanticVad struct {
	// Whether or not to automatically generate a response when a VAD stop event
	// occurs.
	CreateResponse bool `json:"create_response"`
	// Used only for `semantic_vad` mode. The eagerness of the model to respond. `low`
	// will wait longer for the user to continue speaking, `high` will respond more
	// quickly. `auto` is the default and is equivalent to `medium`. `low`, `medium`,
	// and `high` have max timeouts of 8s, 4s, and 2s respectively.
	//
	// Any of "low", "medium", "high", "auto".
	Eagerness string `json:"eagerness"`
	// Whether or not to automatically interrupt any ongoing response with output to
	// the default conversation (i.e. `conversation` of `auto`) when a VAD start event
	// occurs.
	InterruptResponse bool `json:"interrupt_response"`
	// Type of turn detection, `semantic_vad` to turn on Semantic VAD.
	Type constant.SemanticVad `json:"type" default:"semantic_vad"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		CreateResponse    respjson.Field
		Eagerness         respjson.Field
		InterruptResponse respjson.Field
		Type              respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeTranscriptionSessionAudioInputTurnDetectionSemanticVad) RawJSON() string {
	return r.JSON.raw
}
func (r *RealtimeTranscriptionSessionAudioInputTurnDetectionSemanticVad) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type RealtimeTranscriptionSessionAudioInputTurnDetectionUnionParam struct {
	OfServerVad   *RealtimeTranscriptionSessionAudioInputTurnDetectionServerVadParam   `json:",omitzero,inline"`
	OfSemanticVad *RealtimeTranscriptionSessionAudioInputTurnDetectionSemanticVadParam `json:",omitzero,inline"`
	paramUnion
}

func (u RealtimeTranscriptionSessionAudioInputTurnDetectionUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfServerVad, u.OfSemanticVad)
}
func (u *RealtimeTranscriptionSessionAudioInputTurnDetectionUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeTranscriptionSessionAudioInputTurnDetectionUnionParam) GetIdleTimeoutMs() *int64 {
	if vt := u.OfServerVad; vt != nil && vt.IdleTimeoutMs.Valid() {
		return &vt.IdleTimeoutMs.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeTranscriptionSessionAudioInputTurnDetectionUnionParam) GetPrefixPaddingMs() *int64 {
	if vt := u.OfServerVad; vt != nil && vt.PrefixPaddingMs.Valid() {
		return &vt.PrefixPaddingMs.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeTranscriptionSessionAudioInputTurnDetectionUnionParam) GetSilenceDurationMs() *int64 {
	if vt := u.OfServerVad; vt != nil && vt.SilenceDurationMs.Valid() {
		return &vt.SilenceDurationMs.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeTranscriptionSessionAudioInputTurnDetectionUnionParam) GetThreshold() *float64 {
	if vt := u.OfServerVad; vt != nil && vt.Threshold.Valid() {
		return &vt.Threshold.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeTranscriptionSessionAudioInputTurnDetectionUnionParam) GetEagerness() *string {
	if vt := u.OfSemanticVad; vt != nil {
		return &vt.Eagerness
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeTranscriptionSessionAudioInputTurnDetectionUnionParam) GetType() *string {
	if vt := u.OfServerVad; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfSemanticVad; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeTranscriptionSessionAudioInputTurnDetectionUnionParam) GetCreateResponse() *bool {
	if vt := u.OfServerVad; vt != nil && vt.CreateResponse.Valid() {
		return &vt.CreateResponse.Value
	} else if vt := u.OfSemanticVad; vt != nil && vt.CreateResponse.Valid() {
		return &vt.CreateResponse.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeTranscriptionSessionAudioInputTurnDetectionUnionParam) GetInterruptResponse() *bool {
	if vt := u.OfServerVad; vt != nil && vt.InterruptResponse.Valid() {
		return &vt.InterruptResponse.Value
	} else if vt := u.OfSemanticVad; vt != nil && vt.InterruptResponse.Valid() {
		return &vt.InterruptResponse.Value
	}
	return nil
}

func init() {
	apijson.RegisterUnion[RealtimeTranscriptionSessionAudioInputTurnDetectionUnionParam](
		"type",
		apijson.Discriminator[RealtimeTranscriptionSessionAudioInputTurnDetectionServerVadParam]("server_vad"),
		apijson.Discriminator[RealtimeTranscriptionSessionAudioInputTurnDetectionSemanticVadParam]("semantic_vad"),
	)
}

// Server-side voice activity detection (VAD) which flips on when user speech is
// detected and off after a period of silence.
//
// The property Type is required.
type RealtimeTranscriptionSessionAudioInputTurnDetectionServerVadParam struct {
	// Optional timeout after which a model response will be triggered automatically.
	// This is useful for situations in which a long pause from the user is unexpected,
	// such as a phone call. The model will effectively prompt the user to continue the
	// conversation based on the current context.
	//
	// The timeout value will be applied after the last model response's audio has
	// finished playing, i.e. it's set to the `response.done` time plus audio playback
	// duration.
	//
	// An `input_audio_buffer.timeout_triggered` event (plus events associated with the
	// Response) will be emitted when the timeout is reached. Idle timeout is currently
	// only supported for `server_vad` mode.
	IdleTimeoutMs param.Opt[int64] `json:"idle_timeout_ms,omitzero"`
	// Whether or not to automatically generate a response when a VAD stop event
	// occurs. If `interrupt_response` is set to `false` this may fail to create a
	// response if the model is already responding.
	//
	// If both `create_response` and `interrupt_response` are set to `false`, the model
	// will never respond automatically but VAD events will still be emitted.
	CreateResponse param.Opt[bool] `json:"create_response,omitzero"`
	// Whether or not to automatically interrupt (cancel) any ongoing response with
	// output to the default conversation (i.e. `conversation` of `auto`) when a VAD
	// start event occurs. If `true` then the response will be cancelled, otherwise it
	// will continue until complete.
	//
	// If both `create_response` and `interrupt_response` are set to `false`, the model
	// will never respond automatically but VAD events will still be emitted.
	InterruptResponse param.Opt[bool] `json:"interrupt_response,omitzero"`
	// Used only for `server_vad` mode. Amount of audio to include before the VAD
	// detected speech (in milliseconds). Defaults to 300ms.
	PrefixPaddingMs param.Opt[int64] `json:"prefix_padding_ms,omitzero"`
	// Used only for `server_vad` mode. Duration of silence to detect speech stop (in
	// milliseconds). Defaults to 500ms. With shorter values the model will respond
	// more quickly, but may jump in on short pauses from the user.
	SilenceDurationMs param.Opt[int64] `json:"silence_duration_ms,omitzero"`
	// Used only for `server_vad` mode. Activation threshold for VAD (0.0 to 1.0), this
	// defaults to 0.5. A higher threshold will require louder audio to activate the
	// model, and thus might perform better in noisy environments.
	Threshold param.Opt[float64] `json:"threshold,omitzero"`
	// Type of turn detection, `server_vad` to turn on simple Server VAD.
	//
	// This field can be elided, and will marshal its zero value as "server_vad".
	Type constant.ServerVad `json:"type" default:"server_vad"`
	paramObj
}

func (r RealtimeTranscriptionSessionAudioInputTurnDetectionServerVadParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeTranscriptionSessionAudioInputTurnDetectionServerVadParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeTranscriptionSessionAudioInputTurnDetectionServerVadParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Server-side semantic turn detection which uses a model to determine when the
// user has finished speaking.
//
// The property Type is required.
type RealtimeTranscriptionSessionAudioInputTurnDetectionSemanticVadParam struct {
	// Whether or not to automatically generate a response when a VAD stop event
	// occurs.
	CreateResponse param.Opt[bool] `json:"create_response,omitzero"`
	// Whether or not to automatically interrupt any ongoing response with output to
	// the default conversation (i.e. `conversation` of `auto`) when a VAD start event
	// occurs.
	InterruptResponse param.Opt[bool] `json:"interrupt_response,omitzero"`
	// Used only for `semantic_vad` mode. The eagerness of the model to respond. `low`
	// will wait longer for the user to continue speaking, `high` will respond more
	// quickly. `auto` is the default and is equivalent to `medium`. `low`, `medium`,
	// and `high` have max timeouts of 8s, 4s, and 2s respectively.
	//
	// Any of "low", "medium", "high", "auto".
	Eagerness string `json:"eagerness,omitzero"`
	// Type of turn detection, `semantic_vad` to turn on Semantic VAD.
	//
	// This field can be elided, and will marshal its zero value as "semantic_vad".
	Type constant.SemanticVad `json:"type" default:"semantic_vad"`
	paramObj
}

func (r RealtimeTranscriptionSessionAudioInputTurnDetectionSemanticVadParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeTranscriptionSessionAudioInputTurnDetectionSemanticVadParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeTranscriptionSessionAudioInputTurnDetectionSemanticVadParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[RealtimeTranscriptionSessionAudioInputTurnDetectionSemanticVadParam](
		"eagerness", "low", "medium", "high", "auto",
	)
}

// Realtime transcription session object configuration.
type RealtimeTranscriptionSessionCreateRequest struct {
	// The type of session to create. Always `transcription` for transcription
	// sessions.
	Type constant.Transcription `json:"type" default:"transcription"`
	// Configuration for input and output audio.
	Audio RealtimeTranscriptionSessionAudio `json:"audio"`
	// Additional fields to include in server outputs.
	//
	// `item.input_audio_transcription.logprobs`: Include logprobs for input audio
	// transcription.
	//
	// Any of "item.input_audio_transcription.logprobs".
	Include []string `json:"include"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		Audio       respjson.Field
		Include     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeTranscriptionSessionCreateRequest) RawJSON() string { return r.JSON.raw }
func (r *RealtimeTranscriptionSessionCreateRequest) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this RealtimeTranscriptionSessionCreateRequest to a
// RealtimeTranscriptionSessionCreateRequestParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// RealtimeTranscriptionSessionCreateRequestParam.Overrides()
func (r RealtimeTranscriptionSessionCreateRequest) ToParam() RealtimeTranscriptionSessionCreateRequestParam {
	return param.Override[RealtimeTranscriptionSessionCreateRequestParam](json.RawMessage(r.RawJSON()))
}

// Realtime transcription session object configuration.
//
// The property Type is required.
type RealtimeTranscriptionSessionCreateRequestParam struct {
	// Configuration for input and output audio.
	Audio RealtimeTranscriptionSessionAudioParam `json:"audio,omitzero"`
	// Additional fields to include in server outputs.
	//
	// `item.input_audio_transcription.logprobs`: Include logprobs for input audio
	// transcription.
	//
	// Any of "item.input_audio_transcription.logprobs".
	Include []string `json:"include,omitzero"`
	// The type of session to create. Always `transcription` for transcription
	// sessions.
	//
	// This field can be elided, and will marshal its zero value as "transcription".
	Type constant.Transcription `json:"type" default:"transcription"`
	paramObj
}

func (r RealtimeTranscriptionSessionCreateRequestParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeTranscriptionSessionCreateRequestParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeTranscriptionSessionCreateRequestParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func RealtimeTranslationClientEventParamOfSessionUpdate(session RealtimeTranslationSessionUpdateRequestParam) RealtimeTranslationClientEventUnionParam {
	var sessionUpdate RealtimeTranslationSessionUpdateEventParam
	sessionUpdate.Session = session
	return RealtimeTranslationClientEventUnionParam{OfSessionUpdate: &sessionUpdate}
}

func RealtimeTranslationClientEventParamOfSessionInputAudioBufferAppend(audio string) RealtimeTranslationClientEventUnionParam {
	var sessionInputAudioBufferAppend RealtimeTranslationInputAudioBufferAppendEventParam
	sessionInputAudioBufferAppend.Audio = audio
	return RealtimeTranslationClientEventUnionParam{OfSessionInputAudioBufferAppend: &sessionInputAudioBufferAppend}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type RealtimeTranslationClientEventUnionParam struct {
	OfSessionUpdate                 *RealtimeTranslationSessionUpdateEventParam          `json:",omitzero,inline"`
	OfSessionInputAudioBufferAppend *RealtimeTranslationInputAudioBufferAppendEventParam `json:",omitzero,inline"`
	OfSessionClose                  *RealtimeTranslationSessionCloseEventParam           `json:",omitzero,inline"`
	paramUnion
}

func (u RealtimeTranslationClientEventUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfSessionUpdate, u.OfSessionInputAudioBufferAppend, u.OfSessionClose)
}
func (u *RealtimeTranslationClientEventUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeTranslationClientEventUnionParam) GetSession() *RealtimeTranslationSessionUpdateRequestParam {
	if vt := u.OfSessionUpdate; vt != nil {
		return &vt.Session
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeTranslationClientEventUnionParam) GetAudio() *string {
	if vt := u.OfSessionInputAudioBufferAppend; vt != nil {
		return &vt.Audio
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeTranslationClientEventUnionParam) GetType() *string {
	if vt := u.OfSessionUpdate; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfSessionInputAudioBufferAppend; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfSessionClose; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RealtimeTranslationClientEventUnionParam) GetEventID() *string {
	if vt := u.OfSessionUpdate; vt != nil && vt.EventID.Valid() {
		return &vt.EventID.Value
	} else if vt := u.OfSessionInputAudioBufferAppend; vt != nil && vt.EventID.Valid() {
		return &vt.EventID.Value
	} else if vt := u.OfSessionClose; vt != nil && vt.EventID.Valid() {
		return &vt.EventID.Value
	}
	return nil
}

func init() {
	apijson.RegisterUnion[RealtimeTranslationClientEventUnionParam](
		"type",
		apijson.Discriminator[RealtimeTranslationSessionUpdateEventParam]("session.update"),
		apijson.Discriminator[RealtimeTranslationInputAudioBufferAppendEventParam]("session.input_audio_buffer.append"),
		apijson.Discriminator[RealtimeTranslationSessionCloseEventParam]("session.close"),
	)
}

// Create a translation session and client secret for the Realtime API.
//
// The property Session is required.
type RealtimeTranslationClientSecretCreateRequestParam struct {
	// Realtime translation session configuration. Translation sessions stream source
	// audio in and translated audio plus transcript deltas out continuously.
	Session RealtimeTranslationSessionCreateRequestParam `json:"session,omitzero" api:"required"`
	// Configuration for the client secret expiration. Expiration refers to the time
	// after which a client secret will no longer be valid for creating sessions. The
	// session itself may continue after that time once started. A secret can be used
	// to create multiple sessions until it expires.
	ExpiresAfter RealtimeTranslationClientSecretCreateRequestExpiresAfterParam `json:"expires_after,omitzero"`
	paramObj
}

func (r RealtimeTranslationClientSecretCreateRequestParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeTranslationClientSecretCreateRequestParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeTranslationClientSecretCreateRequestParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for the client secret expiration. Expiration refers to the time
// after which a client secret will no longer be valid for creating sessions. The
// session itself may continue after that time once started. A secret can be used
// to create multiple sessions until it expires.
type RealtimeTranslationClientSecretCreateRequestExpiresAfterParam struct {
	// The number of seconds from the anchor point to the expiration. Select a value
	// between `10` and `7200` (2 hours). This default to 600 seconds (10 minutes) if
	// not specified.
	Seconds param.Opt[int64] `json:"seconds,omitzero"`
	// The anchor point for the client secret expiration, meaning that `seconds` will
	// be added to the `created_at` time of the client secret to produce an expiration
	// timestamp. Only `created_at` is currently supported.
	//
	// Any of "created_at".
	Anchor string `json:"anchor,omitzero"`
	paramObj
}

func (r RealtimeTranslationClientSecretCreateRequestExpiresAfterParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeTranslationClientSecretCreateRequestExpiresAfterParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeTranslationClientSecretCreateRequestExpiresAfterParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[RealtimeTranslationClientSecretCreateRequestExpiresAfterParam](
		"anchor", "created_at",
	)
}

// Response from creating a translation session and client secret for the Realtime
// API.
type RealtimeTranslationClientSecretCreateResponse struct {
	// Expiration timestamp for the client secret, in seconds since epoch.
	ExpiresAt int64 `json:"expires_at" api:"required" format:"unixtime"`
	// A Realtime translation session. Translation sessions continuously translate
	// input audio into the configured output language.
	Session RealtimeTranslationSession `json:"session" api:"required"`
	// The generated client secret value.
	Value string `json:"value" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ExpiresAt   respjson.Field
		Session     respjson.Field
		Value       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeTranslationClientSecretCreateResponse) RawJSON() string { return r.JSON.raw }
func (r *RealtimeTranslationClientSecretCreateResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Send this event to append audio bytes to the translation session input audio
// buffer.
//
// WebSocket translation sessions accept base64-encoded 24 kHz PCM16 mono
// little-endian raw audio bytes. Unsupported websocket audio formats return a
// validation error because lower-quality audio materially degrades translation
// quality.
//
// Translation consumes 200 ms engine frames. For best realtime behavior, append
// audio in 200 ms chunks. If a chunk is shorter, the server buffers it until it
// has enough audio for one frame. If a chunk is longer, the server splits it into
// 200 ms frames and enqueues them back-to-back.
//
// Keep appending silence while the session is active. If a client stops sending
// audio and later resumes, model time treats the resumed audio as contiguous with
// the previous audio rather than as a real-world pause.
//
// The properties Audio, Type are required.
type RealtimeTranslationInputAudioBufferAppendEventParam struct {
	// Base64-encoded 24 kHz PCM16 mono audio bytes.
	Audio string `json:"audio" api:"required"`
	// Optional client-generated ID used to identify this event.
	EventID param.Opt[string] `json:"event_id,omitzero"`
	// The event type, must be `session.input_audio_buffer.append`.
	//
	// This field can be elided, and will marshal its zero value as
	// "session.input_audio_buffer.append".
	Type constant.SessionInputAudioBufferAppend `json:"type" default:"session.input_audio_buffer.append"`
	paramObj
}

func (r RealtimeTranslationInputAudioBufferAppendEventParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeTranslationInputAudioBufferAppendEventParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeTranslationInputAudioBufferAppendEventParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Returned when optional source-language transcript text is available. This event
// is emitted only when `audio.input.transcription` is configured.
//
// Transcript deltas are append-only text fragments. Clients should not insert
// unconditional spaces between deltas.
type RealtimeTranslationInputTranscriptDeltaEvent struct {
	// Append-only source-language transcript text.
	Delta string `json:"delta" api:"required"`
	// The unique ID of the server event.
	EventID string `json:"event_id" api:"required"`
	// The event type, must be `session.input_transcript.delta`.
	Type constant.SessionInputTranscriptDelta `json:"type" default:"session.input_transcript.delta"`
	// Timing metadata for stream alignment, derived from the translation frame when
	// available. It advances in 200 ms increments, but multiple transcript deltas may
	// share the same `elapsed_ms`. Treat it as alignment metadata, not a unique
	// transcript-delta identifier.
	ElapsedMs int64 `json:"elapsed_ms" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Delta       respjson.Field
		EventID     respjson.Field
		Type        respjson.Field
		ElapsedMs   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeTranslationInputTranscriptDeltaEvent) RawJSON() string { return r.JSON.raw }
func (r *RealtimeTranslationInputTranscriptDeltaEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Returned when translated output audio is available. The `delta` contains a PCM16
// audio chunk whose length can vary. Clients should decode and queue the complete
// delta instead of assuming a fixed byte or sample count.
type RealtimeTranslationOutputAudioDeltaEvent struct {
	// Base64-encoded translated audio data.
	Delta string `json:"delta" api:"required"`
	// The unique ID of the server event.
	EventID string `json:"event_id" api:"required"`
	// The event type, must be `session.output_audio.delta`.
	Type constant.SessionOutputAudioDelta `json:"type" default:"session.output_audio.delta"`
	// Number of audio channels.
	Channels int64 `json:"channels"`
	// Timing metadata for stream alignment, derived from the translation frame when
	// available. Treat `elapsed_ms` as alignment metadata, not a unique event
	// identifier.
	ElapsedMs int64 `json:"elapsed_ms" api:"nullable"`
	// Audio encoding for `delta`.
	//
	// Any of "pcm16".
	Format RealtimeTranslationOutputAudioDeltaEventFormat `json:"format"`
	// Sample rate of the audio delta.
	SampleRate int64 `json:"sample_rate"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Delta       respjson.Field
		EventID     respjson.Field
		Type        respjson.Field
		Channels    respjson.Field
		ElapsedMs   respjson.Field
		Format      respjson.Field
		SampleRate  respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeTranslationOutputAudioDeltaEvent) RawJSON() string { return r.JSON.raw }
func (r *RealtimeTranslationOutputAudioDeltaEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Audio encoding for `delta`.
type RealtimeTranslationOutputAudioDeltaEventFormat string

const (
	RealtimeTranslationOutputAudioDeltaEventFormatPcm16 RealtimeTranslationOutputAudioDeltaEventFormat = "pcm16"
)

// Returned when translated transcript text is available.
//
// Transcript deltas are append-only text fragments. Clients should not insert
// unconditional spaces between deltas.
type RealtimeTranslationOutputTranscriptDeltaEvent struct {
	// Append-only transcript text for the translated output audio.
	Delta string `json:"delta" api:"required"`
	// The unique ID of the server event.
	EventID string `json:"event_id" api:"required"`
	// The event type, must be `session.output_transcript.delta`.
	Type constant.SessionOutputTranscriptDelta `json:"type" default:"session.output_transcript.delta"`
	// Timing metadata for stream alignment, derived from the translation frame when
	// available. It advances in 200 ms increments, but multiple transcript deltas may
	// share the same `elapsed_ms`. Treat it as alignment metadata, not a unique
	// transcript-delta identifier.
	ElapsedMs int64 `json:"elapsed_ms" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Delta       respjson.Field
		EventID     respjson.Field
		Type        respjson.Field
		ElapsedMs   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeTranslationOutputTranscriptDeltaEvent) RawJSON() string { return r.JSON.raw }
func (r *RealtimeTranslationOutputTranscriptDeltaEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// RealtimeTranslationServerEventUnion contains all possible properties and values
// from [RealtimeErrorEvent], [RealtimeTranslationSessionCreatedEvent],
// [RealtimeTranslationSessionUpdatedEvent],
// [RealtimeTranslationSessionClosedEvent],
// [RealtimeTranslationInputTranscriptDeltaEvent],
// [RealtimeTranslationOutputTranscriptDeltaEvent],
// [RealtimeTranslationOutputAudioDeltaEvent].
//
// Use the [RealtimeTranslationServerEventUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type RealtimeTranslationServerEventUnion struct {
	Type                                            string                                        `json:"type" api:"required"`
	OfRealtimeErrorEvent                            RealtimeErrorEvent                            `json:",inline"`
	OfRealtimeTranslationSessionCreatedEvent        RealtimeTranslationSessionCreatedEvent        `json:",inline"`
	OfRealtimeTranslationSessionUpdatedEvent        RealtimeTranslationSessionUpdatedEvent        `json:",inline"`
	OfRealtimeTranslationSessionClosedEvent         RealtimeTranslationSessionClosedEvent         `json:",inline"`
	OfRealtimeTranslationInputTranscriptDeltaEvent  RealtimeTranslationInputTranscriptDeltaEvent  `json:",inline"`
	OfRealtimeTranslationOutputTranscriptDeltaEvent RealtimeTranslationOutputTranscriptDeltaEvent `json:",inline"`
	OfRealtimeTranslationOutputAudioDeltaEvent      RealtimeTranslationOutputAudioDeltaEvent      `json:",inline"`
	JSON                                            struct {
		Type                                            respjson.Field
		OfRealtimeErrorEvent                            respjson.Field
		OfRealtimeTranslationSessionCreatedEvent        respjson.Field
		OfRealtimeTranslationSessionUpdatedEvent        respjson.Field
		OfRealtimeTranslationSessionClosedEvent         respjson.Field
		OfRealtimeTranslationInputTranscriptDeltaEvent  respjson.Field
		OfRealtimeTranslationOutputTranscriptDeltaEvent respjson.Field
		OfRealtimeTranslationOutputAudioDeltaEvent      respjson.Field
		raw                                             string
	} `json:"-"`
}

// anyRealtimeTranslationServerEvent is implemented by each variant of
// [RealtimeTranslationServerEventUnion] to add type safety for the return type of
// [RealtimeTranslationServerEventUnion.AsAny]
type anyRealtimeTranslationServerEvent interface {
	implRealtimeTranslationServerEventUnion()
}

func (RealtimeErrorEvent) implRealtimeTranslationServerEventUnion()                            {}
func (RealtimeTranslationSessionCreatedEvent) implRealtimeTranslationServerEventUnion()        {}
func (RealtimeTranslationSessionUpdatedEvent) implRealtimeTranslationServerEventUnion()        {}
func (RealtimeTranslationSessionClosedEvent) implRealtimeTranslationServerEventUnion()         {}
func (RealtimeTranslationInputTranscriptDeltaEvent) implRealtimeTranslationServerEventUnion()  {}
func (RealtimeTranslationOutputTranscriptDeltaEvent) implRealtimeTranslationServerEventUnion() {}
func (RealtimeTranslationOutputAudioDeltaEvent) implRealtimeTranslationServerEventUnion()      {}

// Use the following switch statement to find the correct variant
//
//	switch variant := RealtimeTranslationServerEventUnion.AsAny().(type) {
//	case realtime.RealtimeErrorEvent:
//	case realtime.RealtimeTranslationSessionCreatedEvent:
//	case realtime.RealtimeTranslationSessionUpdatedEvent:
//	case realtime.RealtimeTranslationSessionClosedEvent:
//	case realtime.RealtimeTranslationInputTranscriptDeltaEvent:
//	case realtime.RealtimeTranslationOutputTranscriptDeltaEvent:
//	case realtime.RealtimeTranslationOutputAudioDeltaEvent:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u RealtimeTranslationServerEventUnion) AsAny() anyRealtimeTranslationServerEvent {
	switch u.Type {
	case "error":
		return u.AsRealtimeErrorEvent()
	case "session.created":
		return u.AsRealtimeTranslationSessionCreatedEvent()
	case "session.updated":
		return u.AsRealtimeTranslationSessionUpdatedEvent()
	case "session.closed":
		return u.AsRealtimeTranslationSessionClosedEvent()
	case "session.input_transcript.delta":
		return u.AsRealtimeTranslationInputTranscriptDeltaEvent()
	case "session.output_transcript.delta":
		return u.AsRealtimeTranslationOutputTranscriptDeltaEvent()
	case "session.output_audio.delta":
		return u.AsRealtimeTranslationOutputAudioDeltaEvent()
	}
	return nil
}

func (u RealtimeTranslationServerEventUnion) AsRealtimeErrorEvent() (v RealtimeErrorEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeTranslationServerEventUnion) AsRealtimeTranslationSessionCreatedEvent() (v RealtimeTranslationSessionCreatedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeTranslationServerEventUnion) AsRealtimeTranslationSessionUpdatedEvent() (v RealtimeTranslationSessionUpdatedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeTranslationServerEventUnion) AsRealtimeTranslationSessionClosedEvent() (v RealtimeTranslationSessionClosedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeTranslationServerEventUnion) AsRealtimeTranslationInputTranscriptDeltaEvent() (v RealtimeTranslationInputTranscriptDeltaEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeTranslationServerEventUnion) AsRealtimeTranslationOutputTranscriptDeltaEvent() (v RealtimeTranslationOutputTranscriptDeltaEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeTranslationServerEventUnion) AsRealtimeTranslationOutputAudioDeltaEvent() (v RealtimeTranslationOutputAudioDeltaEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u RealtimeTranslationServerEventUnion) RawJSON() string { return u.JSON.raw }

func (r *RealtimeTranslationServerEventUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// unmarshalWebsocketEvent decodes only the event selected by its type. Public
// UnmarshalJSON keeps the existing general-purpose union decoding behavior.
func (r *RealtimeTranslationServerEventUnion) unmarshalWebsocketEvent(data []byte) error {
	var envelope struct {
		Type string `json:"type"`
		JSON struct {
			Type respjson.Field
		} `json:"-"`
	}
	if err := apijson.UnmarshalRoot(data, &envelope); err != nil {
		return err
	}
	eventType := envelope.Type
	*r = RealtimeTranslationServerEventUnion{Type: eventType}
	r.JSON.raw = string(data)
	r.JSON.Type = envelope.JSON.Type
	switch eventType {
	case "error":
		if err := apijson.UnmarshalRoot(data, &r.OfRealtimeErrorEvent); err != nil {
			return err
		}
		r.JSON.OfRealtimeErrorEvent = respjson.NewField(r.JSON.raw)
	case "session.created":
		if err := apijson.UnmarshalRoot(data, &r.OfRealtimeTranslationSessionCreatedEvent); err != nil {
			return err
		}
		r.JSON.OfRealtimeTranslationSessionCreatedEvent = respjson.NewField(r.JSON.raw)
	case "session.updated":
		if err := apijson.UnmarshalRoot(data, &r.OfRealtimeTranslationSessionUpdatedEvent); err != nil {
			return err
		}
		r.JSON.OfRealtimeTranslationSessionUpdatedEvent = respjson.NewField(r.JSON.raw)
	case "session.closed":
		if err := apijson.UnmarshalRoot(data, &r.OfRealtimeTranslationSessionClosedEvent); err != nil {
			return err
		}
		r.JSON.OfRealtimeTranslationSessionClosedEvent = respjson.NewField(r.JSON.raw)
	case "session.input_transcript.delta":
		if err := apijson.UnmarshalRoot(data, &r.OfRealtimeTranslationInputTranscriptDeltaEvent); err != nil {
			return err
		}
		r.JSON.OfRealtimeTranslationInputTranscriptDeltaEvent = respjson.NewField(r.JSON.raw)
	case "session.output_transcript.delta":
		if err := apijson.UnmarshalRoot(data, &r.OfRealtimeTranslationOutputTranscriptDeltaEvent); err != nil {
			return err
		}
		r.JSON.OfRealtimeTranslationOutputTranscriptDeltaEvent = respjson.NewField(r.JSON.raw)
	case "session.output_audio.delta":
		if err := apijson.UnmarshalRoot(data, &r.OfRealtimeTranslationOutputAudioDeltaEvent); err != nil {
			return err
		}
		r.JSON.OfRealtimeTranslationOutputAudioDeltaEvent = respjson.NewField(r.JSON.raw)
	}
	return nil
}

// A Realtime translation session. Translation sessions continuously translate
// input audio into the configured output language.
type RealtimeTranslationSession struct {
	// Unique identifier for the session that looks like `sess_1234567890abcdef`.
	ID string `json:"id" api:"required"`
	// Configuration for translation input and output audio.
	Audio RealtimeTranslationSessionAudio `json:"audio" api:"required"`
	// Expiration timestamp for the session, in seconds since epoch.
	ExpiresAt int64 `json:"expires_at" api:"required" format:"unixtime"`
	// The Realtime translation model used for this session. This field is set at
	// session creation and cannot be changed with `session.update`.
	Model string `json:"model" api:"required"`
	// The session type. Always `translation` for Realtime translation sessions.
	Type constant.Translation `json:"type" default:"translation"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Audio       respjson.Field
		ExpiresAt   respjson.Field
		Model       respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeTranslationSession) RawJSON() string { return r.JSON.raw }
func (r *RealtimeTranslationSession) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for translation input and output audio.
type RealtimeTranslationSessionAudio struct {
	Input  RealtimeTranslationSessionAudioInput  `json:"input"`
	Output RealtimeTranslationSessionAudioOutput `json:"output"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Input       respjson.Field
		Output      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeTranslationSessionAudio) RawJSON() string { return r.JSON.raw }
func (r *RealtimeTranslationSessionAudio) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type RealtimeTranslationSessionAudioInput struct {
	// Optional input noise reduction.
	NoiseReduction RealtimeTranslationSessionAudioInputNoiseReduction `json:"noise_reduction" api:"nullable"`
	// Optional source-language transcription. When configured, the server emits
	// `session.input_transcript.delta` events. Translation itself still runs from the
	// input audio stream.
	Transcription RealtimeTranslationSessionAudioInputTranscription `json:"transcription" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		NoiseReduction respjson.Field
		Transcription  respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeTranslationSessionAudioInput) RawJSON() string { return r.JSON.raw }
func (r *RealtimeTranslationSessionAudioInput) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Optional input noise reduction.
type RealtimeTranslationSessionAudioInputNoiseReduction struct {
	// Type of noise reduction. `near_field` is for close-talking microphones such as
	// headphones, `far_field` is for far-field microphones such as laptop or
	// conference room microphones.
	//
	// Any of "near_field", "far_field".
	Type NoiseReductionType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeTranslationSessionAudioInputNoiseReduction) RawJSON() string { return r.JSON.raw }
func (r *RealtimeTranslationSessionAudioInputNoiseReduction) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Optional source-language transcription. When configured, the server emits
// `session.input_transcript.delta` events. Translation itself still runs from the
// input audio stream.
type RealtimeTranslationSessionAudioInputTranscription struct {
	// The transcription model used for source transcript deltas.
	Model string `json:"model" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Model       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeTranslationSessionAudioInputTranscription) RawJSON() string { return r.JSON.raw }
func (r *RealtimeTranslationSessionAudioInputTranscription) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type RealtimeTranslationSessionAudioOutput struct {
	// Target language for translated output audio and transcript deltas.
	Language string `json:"language"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Language    respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeTranslationSessionAudioOutput) RawJSON() string { return r.JSON.raw }
func (r *RealtimeTranslationSessionAudioOutput) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Gracefully close the realtime translation session. The server flushes pending
// input audio and emits any remaining translated output before closing the
// session.
//
// The property Type is required.
type RealtimeTranslationSessionCloseEventParam struct {
	// Optional client-generated ID used to identify this event.
	EventID param.Opt[string] `json:"event_id,omitzero"`
	// The event type, must be `session.close`.
	//
	// This field can be elided, and will marshal its zero value as "session.close".
	Type constant.SessionClose `json:"type" default:"session.close"`
	paramObj
}

func (r RealtimeTranslationSessionCloseEventParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeTranslationSessionCloseEventParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeTranslationSessionCloseEventParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Returned when a realtime translation session is closed.
type RealtimeTranslationSessionClosedEvent struct {
	// The unique ID of the server event.
	EventID string `json:"event_id" api:"required"`
	// The event type, must be `session.closed`.
	Type constant.SessionClosed `json:"type" default:"session.closed"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventID     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeTranslationSessionClosedEvent) RawJSON() string { return r.JSON.raw }
func (r *RealtimeTranslationSessionClosedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Realtime translation session configuration. Translation sessions stream source
// audio in and translated audio plus transcript deltas out continuously.
//
// The property Model is required.
type RealtimeTranslationSessionCreateRequestParam struct {
	// The Realtime translation model used for this session.
	Model string `json:"model" api:"required"`
	// Configuration for translation input and output audio.
	Audio RealtimeTranslationSessionCreateRequestAudioParam `json:"audio,omitzero"`
	paramObj
}

func (r RealtimeTranslationSessionCreateRequestParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeTranslationSessionCreateRequestParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeTranslationSessionCreateRequestParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for translation input and output audio.
type RealtimeTranslationSessionCreateRequestAudioParam struct {
	Input  RealtimeTranslationSessionCreateRequestAudioInputParam  `json:"input,omitzero"`
	Output RealtimeTranslationSessionCreateRequestAudioOutputParam `json:"output,omitzero"`
	paramObj
}

func (r RealtimeTranslationSessionCreateRequestAudioParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeTranslationSessionCreateRequestAudioParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeTranslationSessionCreateRequestAudioParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type RealtimeTranslationSessionCreateRequestAudioInputParam struct {
	// Optional input noise reduction. Set to `null` to disable it.
	NoiseReduction RealtimeTranslationSessionCreateRequestAudioInputNoiseReductionParam `json:"noise_reduction,omitzero"`
	// Optional source-language transcription. When configured, the server emits
	// `session.input_transcript.delta` events. Translation itself still runs from the
	// input audio stream.
	Transcription RealtimeTranslationSessionCreateRequestAudioInputTranscriptionParam `json:"transcription,omitzero"`
	paramObj
}

func (r RealtimeTranslationSessionCreateRequestAudioInputParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeTranslationSessionCreateRequestAudioInputParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeTranslationSessionCreateRequestAudioInputParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Optional input noise reduction. Set to `null` to disable it.
//
// The property Type is required.
type RealtimeTranslationSessionCreateRequestAudioInputNoiseReductionParam struct {
	// Type of noise reduction. `near_field` is for close-talking microphones such as
	// headphones, `far_field` is for far-field microphones such as laptop or
	// conference room microphones.
	//
	// Any of "near_field", "far_field".
	Type NoiseReductionType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r RealtimeTranslationSessionCreateRequestAudioInputNoiseReductionParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeTranslationSessionCreateRequestAudioInputNoiseReductionParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeTranslationSessionCreateRequestAudioInputNoiseReductionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Optional source-language transcription. When configured, the server emits
// `session.input_transcript.delta` events. Translation itself still runs from the
// input audio stream.
//
// The property Model is required.
type RealtimeTranslationSessionCreateRequestAudioInputTranscriptionParam struct {
	// The transcription model to use for source transcript deltas.
	Model string `json:"model" api:"required"`
	paramObj
}

func (r RealtimeTranslationSessionCreateRequestAudioInputTranscriptionParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeTranslationSessionCreateRequestAudioInputTranscriptionParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeTranslationSessionCreateRequestAudioInputTranscriptionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type RealtimeTranslationSessionCreateRequestAudioOutputParam struct {
	// Target language for translated output audio and transcript deltas.
	Language param.Opt[string] `json:"language,omitzero"`
	paramObj
}

func (r RealtimeTranslationSessionCreateRequestAudioOutputParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeTranslationSessionCreateRequestAudioOutputParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeTranslationSessionCreateRequestAudioOutputParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Returned when a translation session is created. Emitted automatically when a new
// connection is established as the first server event. This event contains the
// default translation session configuration.
type RealtimeTranslationSessionCreatedEvent struct {
	// The unique ID of the server event.
	EventID string `json:"event_id" api:"required"`
	// The translation session configuration.
	Session RealtimeTranslationSession `json:"session" api:"required"`
	// The event type, must be `session.created`.
	Type constant.SessionCreated `json:"type" default:"session.created"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventID     respjson.Field
		Session     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeTranslationSessionCreatedEvent) RawJSON() string { return r.JSON.raw }
func (r *RealtimeTranslationSessionCreatedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Send this event to update the translation session configuration. Translation
// sessions support updates to `audio.output.language`,
// `audio.input.transcription`, and `audio.input.noise_reduction`.
//
// The properties Session, Type are required.
type RealtimeTranslationSessionUpdateEventParam struct {
	// Translation session fields to update. The session `type` and `model` are set at
	// creation and cannot be changed with `session.update`.
	Session RealtimeTranslationSessionUpdateRequestParam `json:"session" api:"required"`
	// Optional client-generated ID used to identify this event.
	EventID param.Opt[string] `json:"event_id,omitzero"`
	// The event type, must be `session.update`.
	//
	// This field can be elided, and will marshal its zero value as "session.update".
	Type constant.SessionUpdate `json:"type" default:"session.update"`
	paramObj
}

func (r RealtimeTranslationSessionUpdateEventParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeTranslationSessionUpdateEventParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeTranslationSessionUpdateEventParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Realtime translation session fields that can be updated with `session.update`.
type RealtimeTranslationSessionUpdateRequestParam struct {
	// Configuration for translation input and output audio.
	Audio RealtimeTranslationSessionUpdateRequestAudioParam `json:"audio,omitzero"`
	paramObj
}

func (r RealtimeTranslationSessionUpdateRequestParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeTranslationSessionUpdateRequestParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeTranslationSessionUpdateRequestParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for translation input and output audio.
type RealtimeTranslationSessionUpdateRequestAudioParam struct {
	Input  RealtimeTranslationSessionUpdateRequestAudioInputParam  `json:"input,omitzero"`
	Output RealtimeTranslationSessionUpdateRequestAudioOutputParam `json:"output,omitzero"`
	paramObj
}

func (r RealtimeTranslationSessionUpdateRequestAudioParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeTranslationSessionUpdateRequestAudioParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeTranslationSessionUpdateRequestAudioParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type RealtimeTranslationSessionUpdateRequestAudioInputParam struct {
	// Optional input noise reduction. Set to `null` to disable it.
	NoiseReduction RealtimeTranslationSessionUpdateRequestAudioInputNoiseReductionParam `json:"noise_reduction,omitzero"`
	// Optional source-language transcription. When configured, the server emits
	// `session.input_transcript.delta` events. Translation itself still runs from the
	// input audio stream.
	Transcription RealtimeTranslationSessionUpdateRequestAudioInputTranscriptionParam `json:"transcription,omitzero"`
	paramObj
}

func (r RealtimeTranslationSessionUpdateRequestAudioInputParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeTranslationSessionUpdateRequestAudioInputParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeTranslationSessionUpdateRequestAudioInputParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Optional input noise reduction. Set to `null` to disable it.
//
// The property Type is required.
type RealtimeTranslationSessionUpdateRequestAudioInputNoiseReductionParam struct {
	// Type of noise reduction. `near_field` is for close-talking microphones such as
	// headphones, `far_field` is for far-field microphones such as laptop or
	// conference room microphones.
	//
	// Any of "near_field", "far_field".
	Type NoiseReductionType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r RealtimeTranslationSessionUpdateRequestAudioInputNoiseReductionParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeTranslationSessionUpdateRequestAudioInputNoiseReductionParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeTranslationSessionUpdateRequestAudioInputNoiseReductionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Optional source-language transcription. When configured, the server emits
// `session.input_transcript.delta` events. Translation itself still runs from the
// input audio stream.
//
// The property Model is required.
type RealtimeTranslationSessionUpdateRequestAudioInputTranscriptionParam struct {
	// The transcription model to use for source transcript deltas.
	Model string `json:"model" api:"required"`
	paramObj
}

func (r RealtimeTranslationSessionUpdateRequestAudioInputTranscriptionParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeTranslationSessionUpdateRequestAudioInputTranscriptionParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeTranslationSessionUpdateRequestAudioInputTranscriptionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type RealtimeTranslationSessionUpdateRequestAudioOutputParam struct {
	// Target language for translated output audio and transcript deltas.
	Language param.Opt[string] `json:"language,omitzero"`
	paramObj
}

func (r RealtimeTranslationSessionUpdateRequestAudioOutputParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeTranslationSessionUpdateRequestAudioOutputParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeTranslationSessionUpdateRequestAudioOutputParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Returned when a translation session is updated with a `session.update` event,
// unless there is an error.
type RealtimeTranslationSessionUpdatedEvent struct {
	// The unique ID of the server event.
	EventID string `json:"event_id" api:"required"`
	// The translation session configuration.
	Session RealtimeTranslationSession `json:"session" api:"required"`
	// The event type, must be `session.updated`.
	Type constant.SessionUpdated `json:"type" default:"session.updated"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventID     respjson.Field
		Session     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeTranslationSessionUpdatedEvent) RawJSON() string { return r.JSON.raw }
func (r *RealtimeTranslationSessionUpdatedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// RealtimeTruncationUnion contains all possible properties and values from
// [string], [RealtimeTruncationRetentionRatio].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfRealtimeTruncationStrategy]
type RealtimeTruncationUnion struct {
	// This field will be present if the value is a [string] instead of an object.
	OfRealtimeTruncationStrategy string `json:",inline"`
	// This field is from variant [RealtimeTruncationRetentionRatio].
	RetentionRatio float64 `json:"retention_ratio"`
	// This field is from variant [RealtimeTruncationRetentionRatio].
	Type constant.RetentionRatio `json:"type"`
	// This field is from variant [RealtimeTruncationRetentionRatio].
	TokenLimits RealtimeTruncationRetentionRatioTokenLimits `json:"token_limits"`
	JSON        struct {
		OfRealtimeTruncationStrategy respjson.Field
		RetentionRatio               respjson.Field
		Type                         respjson.Field
		TokenLimits                  respjson.Field
		raw                          string
	} `json:"-"`
}

func (u RealtimeTruncationUnion) AsRealtimeTruncationStrategy() (v string) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeTruncationUnion) AsRetentionRatioTruncation() (v RealtimeTruncationRetentionRatio) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u RealtimeTruncationUnion) RawJSON() string { return u.JSON.raw }

func (r *RealtimeTruncationUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this RealtimeTruncationUnion to a RealtimeTruncationUnionParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// RealtimeTruncationUnionParam.Overrides()
func (r RealtimeTruncationUnion) ToParam() RealtimeTruncationUnionParam {
	return param.Override[RealtimeTruncationUnionParam](json.RawMessage(r.RawJSON()))
}

// The truncation strategy to use for the session. `auto` is the default truncation
// strategy. `disabled` will disable truncation and emit errors when the
// conversation exceeds the input token limit.
type RealtimeTruncationRealtimeTruncationStrategy string

const (
	RealtimeTruncationRealtimeTruncationStrategyAuto     RealtimeTruncationRealtimeTruncationStrategy = "auto"
	RealtimeTruncationRealtimeTruncationStrategyDisabled RealtimeTruncationRealtimeTruncationStrategy = "disabled"
)

func RealtimeTruncationParamOfRetentionRatioTruncation(retentionRatio float64) RealtimeTruncationUnionParam {
	var variant RealtimeTruncationRetentionRatioParam
	variant.RetentionRatio = retentionRatio
	return RealtimeTruncationUnionParam{OfRetentionRatioTruncation: &variant}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type RealtimeTruncationUnionParam struct {
	// Check if union is this variant with
	// !param.IsOmitted(union.OfRealtimeTruncationStrategy)
	OfRealtimeTruncationStrategy param.Opt[string]                      `json:",omitzero,inline"`
	OfRetentionRatioTruncation   *RealtimeTruncationRetentionRatioParam `json:",omitzero,inline"`
	paramUnion
}

func (u RealtimeTruncationUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfRealtimeTruncationStrategy, u.OfRetentionRatioTruncation)
}
func (u *RealtimeTruncationUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Retain a fraction of the conversation tokens when the conversation exceeds the
// input token limit. This allows you to amortize truncations across multiple
// turns, which can help improve cached token usage.
type RealtimeTruncationRetentionRatio struct {
	// Fraction of post-instruction conversation tokens to retain (`0.0` - `1.0`) when
	// the conversation exceeds the input token limit. Setting this to `0.8` means that
	// messages will be dropped until 80% of the maximum allowed tokens are used. This
	// helps reduce the frequency of truncations and improve cache rates.
	RetentionRatio float64 `json:"retention_ratio" api:"required"`
	// Use retention ratio truncation.
	Type constant.RetentionRatio `json:"type" default:"retention_ratio"`
	// Optional custom token limits for this truncation strategy. If not provided, the
	// model's default token limits will be used.
	TokenLimits RealtimeTruncationRetentionRatioTokenLimits `json:"token_limits"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		RetentionRatio respjson.Field
		Type           respjson.Field
		TokenLimits    respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeTruncationRetentionRatio) RawJSON() string { return r.JSON.raw }
func (r *RealtimeTruncationRetentionRatio) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this RealtimeTruncationRetentionRatio to a
// RealtimeTruncationRetentionRatioParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// RealtimeTruncationRetentionRatioParam.Overrides()
func (r RealtimeTruncationRetentionRatio) ToParam() RealtimeTruncationRetentionRatioParam {
	return param.Override[RealtimeTruncationRetentionRatioParam](json.RawMessage(r.RawJSON()))
}

// Optional custom token limits for this truncation strategy. If not provided, the
// model's default token limits will be used.
type RealtimeTruncationRetentionRatioTokenLimits struct {
	// Maximum tokens allowed in the conversation after instructions (which including
	// tool definitions). For example, setting this to 5,000 would mean that truncation
	// would occur when the conversation exceeds 5,000 tokens after instructions. This
	// cannot be higher than the model's context window size minus the maximum output
	// tokens.
	PostInstructions int64 `json:"post_instructions"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		PostInstructions respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeTruncationRetentionRatioTokenLimits) RawJSON() string { return r.JSON.raw }
func (r *RealtimeTruncationRetentionRatioTokenLimits) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Retain a fraction of the conversation tokens when the conversation exceeds the
// input token limit. This allows you to amortize truncations across multiple
// turns, which can help improve cached token usage.
//
// The properties RetentionRatio, Type are required.
type RealtimeTruncationRetentionRatioParam struct {
	// Fraction of post-instruction conversation tokens to retain (`0.0` - `1.0`) when
	// the conversation exceeds the input token limit. Setting this to `0.8` means that
	// messages will be dropped until 80% of the maximum allowed tokens are used. This
	// helps reduce the frequency of truncations and improve cache rates.
	RetentionRatio float64 `json:"retention_ratio" api:"required"`
	// Optional custom token limits for this truncation strategy. If not provided, the
	// model's default token limits will be used.
	TokenLimits RealtimeTruncationRetentionRatioTokenLimitsParam `json:"token_limits,omitzero"`
	// Use retention ratio truncation.
	//
	// This field can be elided, and will marshal its zero value as "retention_ratio".
	Type constant.RetentionRatio `json:"type" default:"retention_ratio"`
	paramObj
}

func (r RealtimeTruncationRetentionRatioParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeTruncationRetentionRatioParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeTruncationRetentionRatioParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Optional custom token limits for this truncation strategy. If not provided, the
// model's default token limits will be used.
type RealtimeTruncationRetentionRatioTokenLimitsParam struct {
	// Maximum tokens allowed in the conversation after instructions (which including
	// tool definitions). For example, setting this to 5,000 would mean that truncation
	// would occur when the conversation exceeds 5,000 tokens after instructions. This
	// cannot be higher than the model's context window size minus the maximum output
	// tokens.
	PostInstructions param.Opt[int64] `json:"post_instructions,omitzero"`
	paramObj
}

func (r RealtimeTruncationRetentionRatioTokenLimitsParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeTruncationRetentionRatioTokenLimitsParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeTruncationRetentionRatioTokenLimitsParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Returned when the model-generated audio is updated.
type ResponseAudioDeltaEvent struct {
	// The index of the content part in the item's content array.
	ContentIndex int64 `json:"content_index" api:"required"`
	// Base64-encoded audio data delta.
	Delta string `json:"delta" api:"required"`
	// The unique ID of the server event.
	EventID string `json:"event_id" api:"required"`
	// The ID of the item.
	ItemID string `json:"item_id" api:"required"`
	// The index of the output item in the response.
	OutputIndex int64 `json:"output_index" api:"required"`
	// The ID of the response.
	ResponseID string `json:"response_id" api:"required"`
	// The event type, must be `response.output_audio.delta`.
	Type constant.ResponseOutputAudioDelta `json:"type" default:"response.output_audio.delta"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ContentIndex respjson.Field
		Delta        respjson.Field
		EventID      respjson.Field
		ItemID       respjson.Field
		OutputIndex  respjson.Field
		ResponseID   respjson.Field
		Type         respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponseAudioDeltaEvent) RawJSON() string { return r.JSON.raw }
func (r *ResponseAudioDeltaEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Returned when the model-generated audio is done. Also emitted when a Response is
// interrupted, incomplete, or cancelled.
type ResponseAudioDoneEvent struct {
	// The index of the content part in the item's content array.
	ContentIndex int64 `json:"content_index" api:"required"`
	// The unique ID of the server event.
	EventID string `json:"event_id" api:"required"`
	// The ID of the item.
	ItemID string `json:"item_id" api:"required"`
	// The index of the output item in the response.
	OutputIndex int64 `json:"output_index" api:"required"`
	// The ID of the response.
	ResponseID string `json:"response_id" api:"required"`
	// The event type, must be `response.output_audio.done`.
	Type constant.ResponseOutputAudioDone `json:"type" default:"response.output_audio.done"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ContentIndex respjson.Field
		EventID      respjson.Field
		ItemID       respjson.Field
		OutputIndex  respjson.Field
		ResponseID   respjson.Field
		Type         respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponseAudioDoneEvent) RawJSON() string { return r.JSON.raw }
func (r *ResponseAudioDoneEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Returned when the model-generated transcription of audio output is updated.
type ResponseAudioTranscriptDeltaEvent struct {
	// The index of the content part in the item's content array.
	ContentIndex int64 `json:"content_index" api:"required"`
	// The transcript delta.
	Delta string `json:"delta" api:"required"`
	// The unique ID of the server event.
	EventID string `json:"event_id" api:"required"`
	// The ID of the item.
	ItemID string `json:"item_id" api:"required"`
	// The index of the output item in the response.
	OutputIndex int64 `json:"output_index" api:"required"`
	// The ID of the response.
	ResponseID string `json:"response_id" api:"required"`
	// The event type, must be `response.output_audio_transcript.delta`.
	Type constant.ResponseOutputAudioTranscriptDelta `json:"type" default:"response.output_audio_transcript.delta"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ContentIndex respjson.Field
		Delta        respjson.Field
		EventID      respjson.Field
		ItemID       respjson.Field
		OutputIndex  respjson.Field
		ResponseID   respjson.Field
		Type         respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponseAudioTranscriptDeltaEvent) RawJSON() string { return r.JSON.raw }
func (r *ResponseAudioTranscriptDeltaEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Returned when the model-generated transcription of audio output is done
// streaming. Also emitted when a Response is interrupted, incomplete, or
// cancelled.
type ResponseAudioTranscriptDoneEvent struct {
	// The index of the content part in the item's content array.
	ContentIndex int64 `json:"content_index" api:"required"`
	// The unique ID of the server event.
	EventID string `json:"event_id" api:"required"`
	// The ID of the item.
	ItemID string `json:"item_id" api:"required"`
	// The index of the output item in the response.
	OutputIndex int64 `json:"output_index" api:"required"`
	// The ID of the response.
	ResponseID string `json:"response_id" api:"required"`
	// The final transcript of the audio.
	Transcript string `json:"transcript" api:"required"`
	// The event type, must be `response.output_audio_transcript.done`.
	Type constant.ResponseOutputAudioTranscriptDone `json:"type" default:"response.output_audio_transcript.done"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ContentIndex respjson.Field
		EventID      respjson.Field
		ItemID       respjson.Field
		OutputIndex  respjson.Field
		ResponseID   respjson.Field
		Transcript   respjson.Field
		Type         respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponseAudioTranscriptDoneEvent) RawJSON() string { return r.JSON.raw }
func (r *ResponseAudioTranscriptDoneEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Send this event to cancel an in-progress response. The server will respond with
// a `response.done` event with a status of `response.status=cancelled`. If there
// is no response to cancel, the server will respond with an error. It's safe to
// call `response.cancel` even if no response is in progress, an error will be
// returned the session will remain unaffected.
//
// The property Type is required.
type ResponseCancelEventParam struct {
	// Optional client-generated ID used to identify this event.
	EventID param.Opt[string] `json:"event_id,omitzero"`
	// A specific response ID to cancel - if not provided, will cancel an in-progress
	// response in the default conversation.
	ResponseID param.Opt[string] `json:"response_id,omitzero"`
	// The event type, must be `response.cancel`.
	//
	// This field can be elided, and will marshal its zero value as "response.cancel".
	Type constant.ResponseCancel `json:"type" default:"response.cancel"`
	paramObj
}

func (r ResponseCancelEventParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponseCancelEventParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponseCancelEventParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Returned when a new content part is added to an assistant message item during
// response generation.
type ResponseContentPartAddedEvent struct {
	// The index of the content part in the item's content array.
	ContentIndex int64 `json:"content_index" api:"required"`
	// The unique ID of the server event.
	EventID string `json:"event_id" api:"required"`
	// The ID of the item to which the content part was added.
	ItemID string `json:"item_id" api:"required"`
	// The index of the output item in the response.
	OutputIndex int64 `json:"output_index" api:"required"`
	// The content part that was added.
	Part ResponseContentPartAddedEventPart `json:"part" api:"required"`
	// The ID of the response.
	ResponseID string `json:"response_id" api:"required"`
	// The event type, must be `response.content_part.added`.
	Type constant.ResponseContentPartAdded `json:"type" default:"response.content_part.added"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ContentIndex respjson.Field
		EventID      respjson.Field
		ItemID       respjson.Field
		OutputIndex  respjson.Field
		Part         respjson.Field
		ResponseID   respjson.Field
		Type         respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponseContentPartAddedEvent) RawJSON() string { return r.JSON.raw }
func (r *ResponseContentPartAddedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The content part that was added.
type ResponseContentPartAddedEventPart struct {
	// Base64-encoded audio data (if type is "audio").
	Audio string `json:"audio"`
	// The text content (if type is "text").
	Text string `json:"text"`
	// The transcript of the audio (if type is "audio").
	Transcript string `json:"transcript"`
	// The content type ("text", "audio").
	//
	// Any of "text", "audio".
	Type string `json:"type"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Audio       respjson.Field
		Text        respjson.Field
		Transcript  respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponseContentPartAddedEventPart) RawJSON() string { return r.JSON.raw }
func (r *ResponseContentPartAddedEventPart) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Returned when a content part is done streaming in an assistant message item.
// Also emitted when a Response is interrupted, incomplete, or cancelled.
type ResponseContentPartDoneEvent struct {
	// The index of the content part in the item's content array.
	ContentIndex int64 `json:"content_index" api:"required"`
	// The unique ID of the server event.
	EventID string `json:"event_id" api:"required"`
	// The ID of the item.
	ItemID string `json:"item_id" api:"required"`
	// The index of the output item in the response.
	OutputIndex int64 `json:"output_index" api:"required"`
	// The content part that is done.
	Part ResponseContentPartDoneEventPart `json:"part" api:"required"`
	// The ID of the response.
	ResponseID string `json:"response_id" api:"required"`
	// The event type, must be `response.content_part.done`.
	Type constant.ResponseContentPartDone `json:"type" default:"response.content_part.done"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ContentIndex respjson.Field
		EventID      respjson.Field
		ItemID       respjson.Field
		OutputIndex  respjson.Field
		Part         respjson.Field
		ResponseID   respjson.Field
		Type         respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponseContentPartDoneEvent) RawJSON() string { return r.JSON.raw }
func (r *ResponseContentPartDoneEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The content part that is done.
type ResponseContentPartDoneEventPart struct {
	// Base64-encoded audio data (if type is "audio").
	Audio string `json:"audio"`
	// The text content (if type is "text").
	Text string `json:"text"`
	// The transcript of the audio (if type is "audio").
	Transcript string `json:"transcript"`
	// The content type ("text", "audio").
	//
	// Any of "text", "audio".
	Type string `json:"type"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Audio       respjson.Field
		Text        respjson.Field
		Transcript  respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponseContentPartDoneEventPart) RawJSON() string { return r.JSON.raw }
func (r *ResponseContentPartDoneEventPart) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// This event instructs the server to create a Response, which means triggering
// model inference. When in Server VAD mode, the server will create Responses
// automatically.
//
// A Response will include at least one Item, and may have two, in which case the
// second will be a function call. These Items will be appended to the conversation
// history by default.
//
// The server will respond with a `response.created` event, events for Items and
// content created, and finally a `response.done` event to indicate the Response is
// complete.
//
// The `response.create` event includes inference configuration like `instructions`
// and `tools`. If these are set, they will override the Session's configuration
// for this Response only.
//
// Responses can be created out-of-band of the default Conversation, meaning that
// they can have arbitrary input, and it's possible to disable writing the output
// to the Conversation. Only one Response can write to the default Conversation at
// a time, but otherwise multiple Responses can be created in parallel. The
// `metadata` field is a good way to disambiguate multiple simultaneous Responses.
//
// Clients can set `conversation` to `none` to create a Response that does not
// write to the default Conversation. Arbitrary input can be provided with the
// `input` field, which is an array accepting raw Items and references to existing
// Items.
//
// The property Type is required.
type ResponseCreateEventParam struct {
	// Optional client-generated ID used to identify this event.
	EventID param.Opt[string] `json:"event_id,omitzero"`
	// Create a new Realtime response with these parameters
	Response RealtimeResponseCreateParams `json:"response,omitzero"`
	// The event type, must be `response.create`.
	//
	// This field can be elided, and will marshal its zero value as "response.create".
	Type constant.ResponseCreate `json:"type" default:"response.create"`
	paramObj
}

func (r ResponseCreateEventParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponseCreateEventParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponseCreateEventParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Returned when a new Response is created. The first event of response creation,
// where the response is in an initial state of `in_progress`.
type ResponseCreatedEvent struct {
	// The unique ID of the server event.
	EventID string `json:"event_id" api:"required"`
	// The response resource.
	Response RealtimeResponse `json:"response" api:"required"`
	// The event type, must be `response.created`.
	Type constant.ResponseCreated `json:"type" default:"response.created"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventID     respjson.Field
		Response    respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponseCreatedEvent) RawJSON() string { return r.JSON.raw }
func (r *ResponseCreatedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Returned when a Response is done streaming. Always emitted, no matter the final
// state. The Response object included in the `response.done` event will include
// all output Items in the Response but will omit the raw audio data.
//
// Clients should check the `status` field of the Response to determine if it was
// successful (`completed`) or if there was another outcome: `cancelled`, `failed`,
// or `incomplete`.
//
// A response will contain all output items that were generated during the
// response, excluding any audio content.
type ResponseDoneEvent struct {
	// The unique ID of the server event.
	EventID string `json:"event_id" api:"required"`
	// The response resource.
	Response RealtimeResponse `json:"response" api:"required"`
	// The event type, must be `response.done`.
	Type constant.ResponseDone `json:"type" default:"response.done"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventID     respjson.Field
		Response    respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponseDoneEvent) RawJSON() string { return r.JSON.raw }
func (r *ResponseDoneEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Returned when the model-generated function call arguments are updated.
type ResponseFunctionCallArgumentsDeltaEvent struct {
	// The ID of the function call.
	CallID string `json:"call_id" api:"required"`
	// The arguments delta as a JSON string.
	Delta string `json:"delta" api:"required"`
	// The unique ID of the server event.
	EventID string `json:"event_id" api:"required"`
	// The ID of the function call item.
	ItemID string `json:"item_id" api:"required"`
	// The index of the output item in the response.
	OutputIndex int64 `json:"output_index" api:"required"`
	// The ID of the response.
	ResponseID string `json:"response_id" api:"required"`
	// The event type, must be `response.function_call_arguments.delta`.
	Type constant.ResponseFunctionCallArgumentsDelta `json:"type" default:"response.function_call_arguments.delta"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		CallID      respjson.Field
		Delta       respjson.Field
		EventID     respjson.Field
		ItemID      respjson.Field
		OutputIndex respjson.Field
		ResponseID  respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponseFunctionCallArgumentsDeltaEvent) RawJSON() string { return r.JSON.raw }
func (r *ResponseFunctionCallArgumentsDeltaEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Returned when the model-generated function call arguments are done streaming.
// Also emitted when a Response is interrupted, incomplete, or cancelled.
type ResponseFunctionCallArgumentsDoneEvent struct {
	// The final arguments as a JSON string.
	Arguments string `json:"arguments" api:"required"`
	// The ID of the function call.
	CallID string `json:"call_id" api:"required"`
	// The unique ID of the server event.
	EventID string `json:"event_id" api:"required"`
	// The ID of the function call item.
	ItemID string `json:"item_id" api:"required"`
	// The name of the function that was called.
	Name string `json:"name" api:"required"`
	// The index of the output item in the response.
	OutputIndex int64 `json:"output_index" api:"required"`
	// The ID of the response.
	ResponseID string `json:"response_id" api:"required"`
	// The event type, must be `response.function_call_arguments.done`.
	Type constant.ResponseFunctionCallArgumentsDone `json:"type" default:"response.function_call_arguments.done"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Arguments   respjson.Field
		CallID      respjson.Field
		EventID     respjson.Field
		ItemID      respjson.Field
		Name        respjson.Field
		OutputIndex respjson.Field
		ResponseID  respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponseFunctionCallArgumentsDoneEvent) RawJSON() string { return r.JSON.raw }
func (r *ResponseFunctionCallArgumentsDoneEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Returned when MCP tool call arguments are updated during response generation.
type ResponseMcpCallArgumentsDelta struct {
	// The JSON-encoded arguments delta.
	Delta string `json:"delta" api:"required"`
	// The unique ID of the server event.
	EventID string `json:"event_id" api:"required"`
	// The ID of the MCP tool call item.
	ItemID string `json:"item_id" api:"required"`
	// The index of the output item in the response.
	OutputIndex int64 `json:"output_index" api:"required"`
	// The ID of the response.
	ResponseID string `json:"response_id" api:"required"`
	// The event type, must be `response.mcp_call_arguments.delta`.
	Type constant.ResponseMcpCallArgumentsDelta `json:"type" default:"response.mcp_call_arguments.delta"`
	// If present, indicates the delta text was obfuscated.
	Obfuscation string `json:"obfuscation" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Delta       respjson.Field
		EventID     respjson.Field
		ItemID      respjson.Field
		OutputIndex respjson.Field
		ResponseID  respjson.Field
		Type        respjson.Field
		Obfuscation respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponseMcpCallArgumentsDelta) RawJSON() string { return r.JSON.raw }
func (r *ResponseMcpCallArgumentsDelta) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Returned when MCP tool call arguments are finalized during response generation.
type ResponseMcpCallArgumentsDone struct {
	// The final JSON-encoded arguments string.
	Arguments string `json:"arguments" api:"required"`
	// The unique ID of the server event.
	EventID string `json:"event_id" api:"required"`
	// The ID of the MCP tool call item.
	ItemID string `json:"item_id" api:"required"`
	// The index of the output item in the response.
	OutputIndex int64 `json:"output_index" api:"required"`
	// The ID of the response.
	ResponseID string `json:"response_id" api:"required"`
	// The event type, must be `response.mcp_call_arguments.done`.
	Type constant.ResponseMcpCallArgumentsDone `json:"type" default:"response.mcp_call_arguments.done"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Arguments   respjson.Field
		EventID     respjson.Field
		ItemID      respjson.Field
		OutputIndex respjson.Field
		ResponseID  respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponseMcpCallArgumentsDone) RawJSON() string { return r.JSON.raw }
func (r *ResponseMcpCallArgumentsDone) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Returned when an MCP tool call has completed successfully.
type ResponseMcpCallCompleted struct {
	// The unique ID of the server event.
	EventID string `json:"event_id" api:"required"`
	// The ID of the MCP tool call item.
	ItemID string `json:"item_id" api:"required"`
	// The index of the output item in the response.
	OutputIndex int64 `json:"output_index" api:"required"`
	// The event type, must be `response.mcp_call.completed`.
	Type constant.ResponseMcpCallCompleted `json:"type" default:"response.mcp_call.completed"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventID     respjson.Field
		ItemID      respjson.Field
		OutputIndex respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponseMcpCallCompleted) RawJSON() string { return r.JSON.raw }
func (r *ResponseMcpCallCompleted) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Returned when an MCP tool call has failed.
type ResponseMcpCallFailed struct {
	// The unique ID of the server event.
	EventID string `json:"event_id" api:"required"`
	// The ID of the MCP tool call item.
	ItemID string `json:"item_id" api:"required"`
	// The index of the output item in the response.
	OutputIndex int64 `json:"output_index" api:"required"`
	// The event type, must be `response.mcp_call.failed`.
	Type constant.ResponseMcpCallFailed `json:"type" default:"response.mcp_call.failed"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventID     respjson.Field
		ItemID      respjson.Field
		OutputIndex respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponseMcpCallFailed) RawJSON() string { return r.JSON.raw }
func (r *ResponseMcpCallFailed) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Returned when an MCP tool call has started and is in progress.
type ResponseMcpCallInProgress struct {
	// The unique ID of the server event.
	EventID string `json:"event_id" api:"required"`
	// The ID of the MCP tool call item.
	ItemID string `json:"item_id" api:"required"`
	// The index of the output item in the response.
	OutputIndex int64 `json:"output_index" api:"required"`
	// The event type, must be `response.mcp_call.in_progress`.
	Type constant.ResponseMcpCallInProgress `json:"type" default:"response.mcp_call.in_progress"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventID     respjson.Field
		ItemID      respjson.Field
		OutputIndex respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponseMcpCallInProgress) RawJSON() string { return r.JSON.raw }
func (r *ResponseMcpCallInProgress) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Returned when a new Item is created during Response generation.
type ResponseOutputItemAddedEvent struct {
	// The unique ID of the server event.
	EventID string `json:"event_id" api:"required"`
	// A single item within a Realtime conversation.
	Item ConversationItemUnion `json:"item" api:"required"`
	// The index of the output item in the Response.
	OutputIndex int64 `json:"output_index" api:"required"`
	// The ID of the Response to which the item belongs.
	ResponseID string `json:"response_id" api:"required"`
	// The event type, must be `response.output_item.added`.
	Type constant.ResponseOutputItemAdded `json:"type" default:"response.output_item.added"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventID     respjson.Field
		Item        respjson.Field
		OutputIndex respjson.Field
		ResponseID  respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponseOutputItemAddedEvent) RawJSON() string { return r.JSON.raw }
func (r *ResponseOutputItemAddedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Returned when an Item is done streaming. Also emitted when a Response is
// interrupted, incomplete, or cancelled.
type ResponseOutputItemDoneEvent struct {
	// The unique ID of the server event.
	EventID string `json:"event_id" api:"required"`
	// A single item within a Realtime conversation.
	Item ConversationItemUnion `json:"item" api:"required"`
	// The index of the output item in the Response.
	OutputIndex int64 `json:"output_index" api:"required"`
	// The ID of the Response to which the item belongs.
	ResponseID string `json:"response_id" api:"required"`
	// The event type, must be `response.output_item.done`.
	Type constant.ResponseOutputItemDone `json:"type" default:"response.output_item.done"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventID     respjson.Field
		Item        respjson.Field
		OutputIndex respjson.Field
		ResponseID  respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponseOutputItemDoneEvent) RawJSON() string { return r.JSON.raw }
func (r *ResponseOutputItemDoneEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Returned when the text value of an "output_text" content part is updated.
type ResponseTextDeltaEvent struct {
	// The index of the content part in the item's content array.
	ContentIndex int64 `json:"content_index" api:"required"`
	// The text delta.
	Delta string `json:"delta" api:"required"`
	// The unique ID of the server event.
	EventID string `json:"event_id" api:"required"`
	// The ID of the item.
	ItemID string `json:"item_id" api:"required"`
	// The index of the output item in the response.
	OutputIndex int64 `json:"output_index" api:"required"`
	// The ID of the response.
	ResponseID string `json:"response_id" api:"required"`
	// The event type, must be `response.output_text.delta`.
	Type constant.ResponseOutputTextDelta `json:"type" default:"response.output_text.delta"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ContentIndex respjson.Field
		Delta        respjson.Field
		EventID      respjson.Field
		ItemID       respjson.Field
		OutputIndex  respjson.Field
		ResponseID   respjson.Field
		Type         respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponseTextDeltaEvent) RawJSON() string { return r.JSON.raw }
func (r *ResponseTextDeltaEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Returned when the text value of an "output_text" content part is done streaming.
// Also emitted when a Response is interrupted, incomplete, or cancelled.
type ResponseTextDoneEvent struct {
	// The index of the content part in the item's content array.
	ContentIndex int64 `json:"content_index" api:"required"`
	// The unique ID of the server event.
	EventID string `json:"event_id" api:"required"`
	// The ID of the item.
	ItemID string `json:"item_id" api:"required"`
	// The index of the output item in the response.
	OutputIndex int64 `json:"output_index" api:"required"`
	// The ID of the response.
	ResponseID string `json:"response_id" api:"required"`
	// The final text content.
	Text string `json:"text" api:"required"`
	// The event type, must be `response.output_text.done`.
	Type constant.ResponseOutputTextDone `json:"type" default:"response.output_text.done"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ContentIndex respjson.Field
		EventID      respjson.Field
		ItemID       respjson.Field
		OutputIndex  respjson.Field
		ResponseID   respjson.Field
		Text         respjson.Field
		Type         respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponseTextDoneEvent) RawJSON() string { return r.JSON.raw }
func (r *ResponseTextDoneEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Returned when a Session is created. Emitted automatically when a new connection
// is established as the first server event. This event will contain the default
// Session configuration.
type SessionCreatedEvent struct {
	// The unique ID of the server event.
	EventID string `json:"event_id" api:"required"`
	// The session configuration.
	Session SessionCreatedEventSessionUnion `json:"session" api:"required"`
	// The event type, must be `session.created`.
	Type constant.SessionCreated `json:"type" default:"session.created"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventID     respjson.Field
		Session     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r SessionCreatedEvent) RawJSON() string { return r.JSON.raw }
func (r *SessionCreatedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// SessionCreatedEventSessionUnion contains all possible properties and values from
// [RealtimeSessionCreateRequest], [RealtimeTranscriptionSessionCreateRequest].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type SessionCreatedEventSessionUnion struct {
	OfRealtimeSessionCreateRequest              RealtimeSessionCreateRequest              `json:",inline"`
	OfRealtimeTranscriptionSessionCreateRequest RealtimeTranscriptionSessionCreateRequest `json:",inline"`
	JSON                                        struct {
		OfRealtimeSessionCreateRequest              respjson.Field
		OfRealtimeTranscriptionSessionCreateRequest respjson.Field
		raw                                         string
	} `json:"-"`
}

func (u SessionCreatedEventSessionUnion) AsRealtimeSessionCreateRequest() (v RealtimeSessionCreateRequest) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u SessionCreatedEventSessionUnion) AsRealtimeTranscriptionSessionCreateRequest() (v RealtimeTranscriptionSessionCreateRequest) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u SessionCreatedEventSessionUnion) RawJSON() string { return u.JSON.raw }

func (r *SessionCreatedEventSessionUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Send this event to update the session’s configuration. The client may send this
// event at any time to update any field except for `voice` and `model`. `voice`
// can be updated only if there have been no other audio outputs yet.
//
// When the server receives a `session.update`, it will respond with a
// `session.updated` event showing the full, effective configuration. Only the
// fields that are present in the `session.update` are updated. To clear a field
// like `instructions`, pass an empty string. To clear a field like `tools`, pass
// an empty array. To clear a field like `turn_detection`, pass `null`.
//
// To turn off input audio noise reduction, send this Realtime event:
//
// ```json
// {
//
//	"type": "session.update",
//	"session": { "type": "realtime", "audio": { "input": { "noise_reduction": null } } }
//
// }
// ```
//
// For a transcription session, use `"type":"transcription"` inside `session`.
// Omitting `audio.input.noise_reduction` from an update leaves its current setting
// unchanged.
//
// The properties Session, Type are required.
type SessionUpdateEventParam struct {
	// Update the Realtime session. Choose either a realtime session or a transcription
	// session.
	Session SessionUpdateEventSessionUnionParam `json:"session,omitzero" api:"required"`
	// Optional client-generated ID used to identify this event. This is an arbitrary
	// string that a client may assign. It will be passed back if there is an error
	// with the event, but the corresponding `session.updated` event will not include
	// it.
	EventID param.Opt[string] `json:"event_id,omitzero"`
	// The event type, must be `session.update`.
	//
	// This field can be elided, and will marshal its zero value as "session.update".
	Type constant.SessionUpdate `json:"type" default:"session.update"`
	paramObj
}

func (r SessionUpdateEventParam) MarshalJSON() (data []byte, err error) {
	type shadow SessionUpdateEventParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *SessionUpdateEventParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type SessionUpdateEventSessionUnionParam struct {
	OfRealtime      *RealtimeSessionCreateRequestParam              `json:",omitzero,inline"`
	OfTranscription *RealtimeTranscriptionSessionCreateRequestParam `json:",omitzero,inline"`
	paramUnion
}

func (u SessionUpdateEventSessionUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfRealtime, u.OfTranscription)
}
func (u *SessionUpdateEventSessionUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u SessionUpdateEventSessionUnionParam) GetInstructions() *string {
	if vt := u.OfRealtime; vt != nil && vt.Instructions.Valid() {
		return &vt.Instructions.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u SessionUpdateEventSessionUnionParam) GetMaxOutputTokens() *RealtimeSessionCreateRequestMaxOutputTokensUnionParam {
	if vt := u.OfRealtime; vt != nil {
		return &vt.MaxOutputTokens
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u SessionUpdateEventSessionUnionParam) GetModel() *RealtimeSessionCreateRequestModel {
	if vt := u.OfRealtime; vt != nil {
		return &vt.Model
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u SessionUpdateEventSessionUnionParam) GetOutputModalities() []string {
	if vt := u.OfRealtime; vt != nil {
		return vt.OutputModalities
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u SessionUpdateEventSessionUnionParam) GetParallelToolCalls() *bool {
	if vt := u.OfRealtime; vt != nil && vt.ParallelToolCalls.Valid() {
		return &vt.ParallelToolCalls.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u SessionUpdateEventSessionUnionParam) GetPrompt() *responses.ResponsePromptParam {
	if vt := u.OfRealtime; vt != nil {
		return &vt.Prompt
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u SessionUpdateEventSessionUnionParam) GetReasoning() *RealtimeReasoningParam {
	if vt := u.OfRealtime; vt != nil {
		return &vt.Reasoning
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u SessionUpdateEventSessionUnionParam) GetToolChoice() *RealtimeToolChoiceConfigUnionParam {
	if vt := u.OfRealtime; vt != nil {
		return &vt.ToolChoice
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u SessionUpdateEventSessionUnionParam) GetTools() RealtimeToolsConfigParam {
	if vt := u.OfRealtime; vt != nil {
		return vt.Tools
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u SessionUpdateEventSessionUnionParam) GetTracing() *RealtimeTracingConfigUnionParam {
	if vt := u.OfRealtime; vt != nil {
		return &vt.Tracing
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u SessionUpdateEventSessionUnionParam) GetTruncation() *RealtimeTruncationUnionParam {
	if vt := u.OfRealtime; vt != nil {
		return &vt.Truncation
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u SessionUpdateEventSessionUnionParam) GetType() *string {
	if vt := u.OfRealtime; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfTranscription; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u SessionUpdateEventSessionUnionParam) GetAudio() (res sessionUpdateEventSessionUnionParamAudio) {
	if vt := u.OfRealtime; vt != nil {
		res.any = &vt.Audio
	} else if vt := u.OfTranscription; vt != nil {
		res.any = &vt.Audio
	}
	return
}

// Can have the runtime types [*RealtimeAudioConfigParam],
// [*RealtimeTranscriptionSessionAudioParam]
type sessionUpdateEventSessionUnionParamAudio struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *realtime.RealtimeAudioConfigParam:
//	case *realtime.RealtimeTranscriptionSessionAudioParam:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u sessionUpdateEventSessionUnionParamAudio) AsAny() any { return u.any }

// Returns a pointer to the underlying variant's property, if present.
func (u sessionUpdateEventSessionUnionParamAudio) GetOutput() *RealtimeAudioConfigOutputParam {
	switch vt := u.any.(type) {
	case *RealtimeAudioConfigParam:
		return &vt.Output
	}
	return nil
}

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u sessionUpdateEventSessionUnionParamAudio) GetInput() (res sessionUpdateEventSessionUnionParamAudioInput) {
	switch vt := u.any.(type) {
	case *RealtimeAudioConfigParam:
		res.any = &vt.Input
	case *RealtimeTranscriptionSessionAudioParam:
		res.any = &vt.Input
	}
	return res
}

// Can have the runtime types [*RealtimeAudioConfigInputParam],
// [*RealtimeTranscriptionSessionAudioInputParam]
type sessionUpdateEventSessionUnionParamAudioInput struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *realtime.RealtimeAudioConfigInputParam:
//	case *realtime.RealtimeTranscriptionSessionAudioInputParam:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u sessionUpdateEventSessionUnionParamAudioInput) AsAny() any { return u.any }

// Returns a pointer to the underlying variant's Format property, if present.
func (u sessionUpdateEventSessionUnionParamAudioInput) GetFormat() *RealtimeAudioFormatsUnionParam {
	switch vt := u.any.(type) {
	case *RealtimeAudioConfigInputParam:
		return &vt.Format
	case *RealtimeTranscriptionSessionAudioInputParam:
		return &vt.Format
	}
	return nil
}

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u sessionUpdateEventSessionUnionParamAudioInput) GetNoiseReduction() (res sessionUpdateEventSessionUnionParamAudioInputNoiseReduction) {
	switch vt := u.any.(type) {
	case *RealtimeAudioConfigInputParam:
		res.any = &vt.NoiseReduction
	case *RealtimeTranscriptionSessionAudioInputParam:
		res.any = &vt.NoiseReduction
	}
	return res
}

// Can have the runtime types [*RealtimeAudioConfigInputNoiseReductionParam],
// [*RealtimeTranscriptionSessionAudioInputNoiseReductionParam]
type sessionUpdateEventSessionUnionParamAudioInputNoiseReduction struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *realtime.RealtimeAudioConfigInputNoiseReductionParam:
//	case *realtime.RealtimeTranscriptionSessionAudioInputNoiseReductionParam:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u sessionUpdateEventSessionUnionParamAudioInputNoiseReduction) AsAny() any { return u.any }

// Returns a pointer to the underlying variant's property, if present.
func (u sessionUpdateEventSessionUnionParamAudioInputNoiseReduction) GetType() *string {
	switch vt := u.any.(type) {
	case *RealtimeAudioConfigInputNoiseReductionParam:
		return (*string)(&vt.Type)
	case *RealtimeTranscriptionSessionAudioInputNoiseReductionParam:
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's Transcription property, if
// present.
func (u sessionUpdateEventSessionUnionParamAudioInput) GetTranscription() *AudioTranscriptionParam {
	switch vt := u.any.(type) {
	case *RealtimeAudioConfigInputParam:
		return &vt.Transcription
	case *RealtimeTranscriptionSessionAudioInputParam:
		return &vt.Transcription
	}
	return nil
}

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u sessionUpdateEventSessionUnionParamAudioInput) GetTurnDetection() (res sessionUpdateEventSessionUnionParamAudioInputTurnDetection) {
	switch vt := u.any.(type) {
	case *RealtimeAudioConfigInputParam:
		res.any = vt.TurnDetection
	case *RealtimeTranscriptionSessionAudioInputParam:
		res.any = vt.TurnDetection
	}
	return res
}

// Can have the runtime types [RealtimeAudioInputTurnDetectionUnionParam],
// [RealtimeTranscriptionSessionAudioInputTurnDetectionUnionParam]
type sessionUpdateEventSessionUnionParamAudioInputTurnDetection struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case realtime.RealtimeAudioInputTurnDetectionUnionParam:
//	case realtime.RealtimeTranscriptionSessionAudioInputTurnDetectionUnionParam:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u sessionUpdateEventSessionUnionParamAudioInputTurnDetection) AsAny() any { return u.any }

// Returns a pointer to the underlying variant's property, if present.
func (u sessionUpdateEventSessionUnionParamAudioInputTurnDetection) GetType() *string {
	switch vt := u.any.(type) {
	case RealtimeAudioInputTurnDetectionUnionParam:
		return vt.GetType()
	case RealtimeTranscriptionSessionAudioInputTurnDetectionUnionParam:
		return vt.GetType()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u sessionUpdateEventSessionUnionParamAudioInputTurnDetection) GetCreateResponse() *bool {
	switch vt := u.any.(type) {
	case RealtimeAudioInputTurnDetectionUnionParam:
		return vt.GetCreateResponse()
	case RealtimeTranscriptionSessionAudioInputTurnDetectionUnionParam:
		return vt.GetCreateResponse()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u sessionUpdateEventSessionUnionParamAudioInputTurnDetection) GetIdleTimeoutMs() *int64 {
	switch vt := u.any.(type) {
	case RealtimeAudioInputTurnDetectionUnionParam:
		return vt.GetIdleTimeoutMs()
	case RealtimeTranscriptionSessionAudioInputTurnDetectionUnionParam:
		return vt.GetIdleTimeoutMs()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u sessionUpdateEventSessionUnionParamAudioInputTurnDetection) GetInterruptResponse() *bool {
	switch vt := u.any.(type) {
	case RealtimeAudioInputTurnDetectionUnionParam:
		return vt.GetInterruptResponse()
	case RealtimeTranscriptionSessionAudioInputTurnDetectionUnionParam:
		return vt.GetInterruptResponse()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u sessionUpdateEventSessionUnionParamAudioInputTurnDetection) GetPrefixPaddingMs() *int64 {
	switch vt := u.any.(type) {
	case RealtimeAudioInputTurnDetectionUnionParam:
		return vt.GetPrefixPaddingMs()
	case RealtimeTranscriptionSessionAudioInputTurnDetectionUnionParam:
		return vt.GetPrefixPaddingMs()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u sessionUpdateEventSessionUnionParamAudioInputTurnDetection) GetSilenceDurationMs() *int64 {
	switch vt := u.any.(type) {
	case RealtimeAudioInputTurnDetectionUnionParam:
		return vt.GetSilenceDurationMs()
	case RealtimeTranscriptionSessionAudioInputTurnDetectionUnionParam:
		return vt.GetSilenceDurationMs()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u sessionUpdateEventSessionUnionParamAudioInputTurnDetection) GetThreshold() *float64 {
	switch vt := u.any.(type) {
	case RealtimeAudioInputTurnDetectionUnionParam:
		return vt.GetThreshold()
	case RealtimeTranscriptionSessionAudioInputTurnDetectionUnionParam:
		return vt.GetThreshold()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u sessionUpdateEventSessionUnionParamAudioInputTurnDetection) GetEagerness() *string {
	switch vt := u.any.(type) {
	case RealtimeAudioInputTurnDetectionUnionParam:
		return vt.GetEagerness()
	case RealtimeTranscriptionSessionAudioInputTurnDetectionUnionParam:
		return vt.GetEagerness()
	}
	return nil
}

// Returns a pointer to the underlying variant's Include property, if present.
func (u SessionUpdateEventSessionUnionParam) GetInclude() []string {
	if vt := u.OfRealtime; vt != nil {
		return vt.Include
	} else if vt := u.OfTranscription; vt != nil {
		return vt.Include
	}
	return nil
}

// Returned when a session is updated with a `session.update` event, unless there
// is an error.
type SessionUpdatedEvent struct {
	// The unique ID of the server event.
	EventID string `json:"event_id" api:"required"`
	// The session configuration.
	Session SessionUpdatedEventSessionUnion `json:"session" api:"required"`
	// The event type, must be `session.updated`.
	Type constant.SessionUpdated `json:"type" default:"session.updated"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventID     respjson.Field
		Session     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r SessionUpdatedEvent) RawJSON() string { return r.JSON.raw }
func (r *SessionUpdatedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// SessionUpdatedEventSessionUnion contains all possible properties and values from
// [RealtimeSessionCreateRequest], [RealtimeTranscriptionSessionCreateRequest].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type SessionUpdatedEventSessionUnion struct {
	OfRealtimeSessionCreateRequest              RealtimeSessionCreateRequest              `json:",inline"`
	OfRealtimeTranscriptionSessionCreateRequest RealtimeTranscriptionSessionCreateRequest `json:",inline"`
	JSON                                        struct {
		OfRealtimeSessionCreateRequest              respjson.Field
		OfRealtimeTranscriptionSessionCreateRequest respjson.Field
		raw                                         string
	} `json:"-"`
}

func (u SessionUpdatedEventSessionUnion) AsRealtimeSessionCreateRequest() (v RealtimeSessionCreateRequest) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u SessionUpdatedEventSessionUnion) AsRealtimeTranscriptionSessionCreateRequest() (v RealtimeTranscriptionSessionCreateRequest) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u SessionUpdatedEventSessionUnion) RawJSON() string { return u.JSON.raw }

func (r *SessionUpdatedEventSessionUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A language detected in transcribed audio.
type TranscriptionLanguage struct {
	// The code of a language detected in the audio.
	Code string `json:"code" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Code        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r TranscriptionLanguage) RawJSON() string { return r.JSON.raw }
func (r *TranscriptionLanguage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// File generated from our OpenAPI spec by Castiron. See CONTRIBUTING.md for details.

package live

import (
	"context"
	"encoding/json"
	"net/http"
	"slices"

	"github.com/openai/openai-go/v4/internal/apijson"
	"github.com/openai/openai-go/v4/internal/paramutil"
	"github.com/openai/openai-go/v4/internal/requestconfig"
	"github.com/openai/openai-go/v4/option"
	"github.com/openai/openai-go/v4/packages/param"
	"github.com/openai/openai-go/v4/packages/respjson"
	"github.com/openai/openai-go/v4/responses"
	"github.com/openai/openai-go/v4/shared/constant"
)

// LiveService contains methods and other services that help with interacting with
// the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewLiveService] method instead.
type LiveService struct {
	Options  []option.RequestOption
	Sideband SidebandService
	Forks    ForkService
	Sessions SessionService
}

// NewLiveService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewLiveService(opts ...option.RequestOption) (r LiveService) {
	r = LiveService{}
	r.Options = requestconfig.InheritedOptions(opts...)
	r.Sideband = NewSidebandService(opts...)
	r.Forks = NewForkService(opts...)
	r.Sessions = NewSessionService(opts...)
	return
}

// Create a Live WebRTC session. Start with the
// [Live prompting guide](https://developers.openai.com/api/docs/guides/live-prompting).
func (r *LiveService) New(ctx context.Context, body LiveNewParams, opts ...option.RequestOption) (res *LiveNewResponse, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	path := "live/sessions"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// AudioFormatUnion contains all possible properties and values from
// [AudioFormatAudioPCM], [AudioFormatAudioPCMU], [AudioFormatAudioPCMA].
//
// Use the [AudioFormatUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type AudioFormatUnion struct {
	OfAudioFormatAudioPCM  AudioFormatAudioPCM  `json:",inline"`
	OfAudioFormatAudioPCMU AudioFormatAudioPCMU `json:",inline"`
	OfAudioFormatAudioPCMA AudioFormatAudioPCMA `json:",inline"`
	Type                   string               `json:"type"`
	JSON                   struct {
		OfAudioFormatAudioPCM  respjson.Field
		OfAudioFormatAudioPCMU respjson.Field
		OfAudioFormatAudioPCMA respjson.Field
		Type                   respjson.Field
		raw                    string
	} `json:"-"`
}

// anyAudioFormat is implemented by each variant of [AudioFormatUnion] to add type
// safety for the return type of [AudioFormatUnion.AsAny]
type anyAudioFormat interface {
	implAudioFormatUnion()
}

func (AudioFormatAudioPCM) implAudioFormatUnion()  {}
func (AudioFormatAudioPCMU) implAudioFormatUnion() {}
func (AudioFormatAudioPCMA) implAudioFormatUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := AudioFormatUnion.AsAny().(type) {
//	case live.AudioFormatAudioPCM:
//	case live.AudioFormatAudioPCMU:
//	case live.AudioFormatAudioPCMA:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u AudioFormatUnion) AsAny() anyAudioFormat {
	switch u.Type {
	case "audio/pcm":
		return u.AsAudioFormatAudioPCM()
	case "audio/pcmu":
		return u.AsAudioFormatAudioPCMU()
	case "audio/pcma":
		return u.AsAudioFormatAudioPCMA()
	}
	return nil
}

func (u AudioFormatUnion) AsAudioFormatAudioPCM() (v AudioFormatAudioPCM) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AudioFormatUnion) AsAudioFormatAudioPCMU() (v AudioFormatAudioPCMU) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AudioFormatUnion) AsAudioFormatAudioPCMA() (v AudioFormatAudioPCMA) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u AudioFormatUnion) RawJSON() string { return u.JSON.raw }

func (r *AudioFormatUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this AudioFormatUnion to a AudioFormatUnionParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// AudioFormatUnionParam.Overrides()
func (r AudioFormatUnion) ToParam() AudioFormatUnionParam {
	return param.Override[AudioFormatUnionParam](json.RawMessage(r.RawJSON()))
}

// Raw, mono 16-bit little-endian PCM audio for a Live WebSocket connection.
type AudioFormatAudioPCM struct {
	// Audio sample rate in hertz. Live WebSocket PCM audio supports 16000 or 24000 Hz.
	//
	// Any of 16000, 24000.
	Rate int64 `json:"rate" api:"required"`
	// The audio encoding. Always `audio/pcm`.
	Type constant.AudioPCM `json:"type" default:"audio/pcm"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Rate        respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AudioFormatAudioPCM) RawJSON() string { return r.JSON.raw }
func (r *AudioFormatAudioPCM) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Raw, mono G.711 μ-law audio for a Live WebSocket connection.
type AudioFormatAudioPCMU struct {
	// Audio sample rate in hertz. G.711 audio uses 8000 Hz.
	Rate int64 `json:"rate" api:"required"`
	// The audio encoding. Always `audio/pcmu`.
	Type constant.AudioPCMU `json:"type" default:"audio/pcmu"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Rate        respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AudioFormatAudioPCMU) RawJSON() string { return r.JSON.raw }
func (r *AudioFormatAudioPCMU) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Raw, mono G.711 A-law audio for a Live WebSocket connection.
type AudioFormatAudioPCMA struct {
	// Audio sample rate in hertz. G.711 audio uses 8000 Hz.
	Rate int64 `json:"rate" api:"required"`
	// The audio encoding. Always `audio/pcma`.
	Type constant.AudioPCMA `json:"type" default:"audio/pcma"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Rate        respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AudioFormatAudioPCMA) RawJSON() string { return r.JSON.raw }
func (r *AudioFormatAudioPCMA) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func AudioFormatParamOfAudioPCM(rate int64) AudioFormatUnionParam {
	var audioPCM AudioFormatAudioPCMParam
	audioPCM.Rate = rate
	return AudioFormatUnionParam{OfAudioPCM: &audioPCM}
}

func AudioFormatParamOfAudioPCMU(rate int64) AudioFormatUnionParam {
	var audioPCMU AudioFormatAudioPCMUParam
	audioPCMU.Rate = rate
	return AudioFormatUnionParam{OfAudioPCMU: &audioPCMU}
}

func AudioFormatParamOfAudioPCMA(rate int64) AudioFormatUnionParam {
	var audioPCMA AudioFormatAudioPCMAParam
	audioPCMA.Rate = rate
	return AudioFormatUnionParam{OfAudioPCMA: &audioPCMA}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type AudioFormatUnionParam struct {
	OfAudioPCM  *AudioFormatAudioPCMParam  `json:",omitzero,inline"`
	OfAudioPCMU *AudioFormatAudioPCMUParam `json:",omitzero,inline"`
	OfAudioPCMA *AudioFormatAudioPCMAParam `json:",omitzero,inline"`
	paramUnion
}

func (u AudioFormatUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfAudioPCM, u.OfAudioPCMU, u.OfAudioPCMA)
}
func (u *AudioFormatUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u AudioFormatUnionParam) GetRate() *int64 {
	if vt := u.OfAudioPCM; vt != nil {
		return (*int64)(&vt.Rate)
	} else if vt := u.OfAudioPCMU; vt != nil {
		return (*int64)(&vt.Rate)
	} else if vt := u.OfAudioPCMA; vt != nil {
		return (*int64)(&vt.Rate)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AudioFormatUnionParam) GetType() *string {
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
	apijson.RegisterUnion[AudioFormatUnionParam](
		"type",
		apijson.Discriminator[AudioFormatAudioPCMParam]("audio/pcm"),
		apijson.Discriminator[AudioFormatAudioPCMUParam]("audio/pcmu"),
		apijson.Discriminator[AudioFormatAudioPCMAParam]("audio/pcma"),
	)
}

// Raw, mono 16-bit little-endian PCM audio for a Live WebSocket connection.
//
// The properties Rate, Type are required.
type AudioFormatAudioPCMParam struct {
	// Audio sample rate in hertz. Live WebSocket PCM audio supports 16000 or 24000 Hz.
	//
	// Any of 16000, 24000.
	Rate int64 `json:"rate,omitzero" api:"required"`
	// The audio encoding. Always `audio/pcm`.
	//
	// This field can be elided, and will marshal its zero value as "audio/pcm".
	Type constant.AudioPCM `json:"type" default:"audio/pcm"`
	paramObj
}

func (r AudioFormatAudioPCMParam) MarshalJSON() (data []byte, err error) {
	type shadow AudioFormatAudioPCMParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AudioFormatAudioPCMParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[AudioFormatAudioPCMParam](
		"rate", 16000, 24000,
	)
}

// Raw, mono G.711 μ-law audio for a Live WebSocket connection.
//
// The properties Rate, Type are required.
type AudioFormatAudioPCMUParam struct {
	// Audio sample rate in hertz. G.711 audio uses 8000 Hz.
	Rate int64 `json:"rate" api:"required"`
	// The audio encoding. Always `audio/pcmu`.
	//
	// This field can be elided, and will marshal its zero value as "audio/pcmu".
	Type constant.AudioPCMU `json:"type" default:"audio/pcmu"`
	paramObj
}

func (r AudioFormatAudioPCMUParam) MarshalJSON() (data []byte, err error) {
	type shadow AudioFormatAudioPCMUParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AudioFormatAudioPCMUParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Raw, mono G.711 A-law audio for a Live WebSocket connection.
//
// The properties Rate, Type are required.
type AudioFormatAudioPCMAParam struct {
	// Audio sample rate in hertz. G.711 audio uses 8000 Hz.
	Rate int64 `json:"rate" api:"required"`
	// The audio encoding. Always `audio/pcma`.
	//
	// This field can be elided, and will marshal its zero value as "audio/pcma".
	Type constant.AudioPCMA `json:"type" default:"audio/pcma"`
	paramObj
}

func (r AudioFormatAudioPCMAParam) MarshalJSON() (data []byte, err error) {
	type shadow AudioFormatAudioPCMAParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AudioFormatAudioPCMAParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A built-in voice available for Live speech.
type BuiltInVoice string

const (
	BuiltInVoiceAlloy    BuiltInVoice = "alloy"
	BuiltInVoiceAsh      BuiltInVoice = "ash"
	BuiltInVoiceBallad   BuiltInVoice = "ballad"
	BuiltInVoiceBeacon   BuiltInVoice = "beacon"
	BuiltInVoiceBossa    BuiltInVoice = "bossa"
	BuiltInVoiceCedar    BuiltInVoice = "cedar"
	BuiltInVoiceCinder   BuiltInVoice = "cinder"
	BuiltInVoiceCoral    BuiltInVoice = "coral"
	BuiltInVoiceDelta    BuiltInVoice = "delta"
	BuiltInVoiceEcho     BuiltInVoice = "echo"
	BuiltInVoiceGleam    BuiltInVoice = "gleam"
	BuiltInVoiceMarin    BuiltInVoice = "marin"
	BuiltInVoiceMeridian BuiltInVoice = "meridian"
	BuiltInVoiceQuartz   BuiltInVoice = "quartz"
	BuiltInVoiceRipple   BuiltInVoice = "ripple"
	BuiltInVoiceSage     BuiltInVoice = "sage"
	BuiltInVoiceShida    BuiltInVoice = "shida"
	BuiltInVoiceShimmer  BuiltInVoice = "shimmer"
	BuiltInVoiceStone    BuiltInVoice = "stone"
	BuiltInVoiceTempo    BuiltInVoice = "tempo"
	BuiltInVoiceVerse    BuiltInVoice = "verse"
	BuiltInVoiceVesper   BuiltInVoice = "vesper"
	BuiltInVoiceWillow   BuiltInVoice = "willow"
)

// Startup-only capabilities for an untrusted frontend attached to a unified WebRTC
// session. Trusted sideband connections are unaffected.
type ClientConfig struct {
	// Client and server event permissions for the WebRTC frontend data channel.
	DataChannel DataChannelConfig `json:"data_channel" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		DataChannel respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ClientConfig) RawJSON() string { return r.JSON.raw }
func (r *ClientConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this ClientConfig to a ClientConfigParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// ClientConfigParam.Overrides()
func (r ClientConfig) ToParam() ClientConfigParam {
	return param.Override[ClientConfigParam](json.RawMessage(r.RawJSON()))
}

// Startup-only capabilities for an untrusted frontend attached to a unified WebRTC
// session. Trusted sideband connections are unaffected.
//
// The property DataChannel is required.
type ClientConfigParam struct {
	// Client and server event permissions for the WebRTC frontend data channel.
	DataChannel DataChannelConfigParam `json:"data_channel,omitzero" api:"required"`
	paramObj
}

func (r ClientConfigParam) MarshalJSON() (data []byte, err error) {
	type shadow ClientConfigParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ClientConfigParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Delegate tasks to your application. The Live session emits delegation events
// that your backend handles.
type ClientDelegation struct {
	// The delegation owner. Always `client` for tasks handled by your application.
	Type constant.Client `json:"type" default:"client"`
	// Milliseconds to wait before requesting an automatic spoken acknowledgement for a
	// delegation. The model may still speak before this delay. Any commentary or
	// thinking appended to the session cancels all pending acknowledgements, including
	// those for other delegations. Defaults to 400 for Decisions sessions and 0
	// (request immediately) for other sessions. Omission in a session update preserves
	// the current value. Updates apply to subsequent delegations.
	AcknowledgementDelayMs int64 `json:"acknowledgement_delay_ms"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type                   respjson.Field
		AcknowledgementDelayMs respjson.Field
		ExtraFields            map[string]respjson.Field
		raw                    string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ClientDelegation) RawJSON() string { return r.JSON.raw }
func (r *ClientDelegation) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this ClientDelegation to a ClientDelegationParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// ClientDelegationParam.Overrides()
func (r ClientDelegation) ToParam() ClientDelegationParam {
	return param.Override[ClientDelegationParam](json.RawMessage(r.RawJSON()))
}

func NewClientDelegationParam() ClientDelegationParam {
	return ClientDelegationParam{
		Type: "client",
	}
}

// Delegate tasks to your application. The Live session emits delegation events
// that your backend handles.
//
// Use [NewClientDelegationParam] to initialize the constant fields.
type ClientDelegationParam struct {
	// Milliseconds to wait before requesting an automatic spoken acknowledgement for a
	// delegation. The model may still speak before this delay. Any commentary or
	// thinking appended to the session cancels all pending acknowledgements, including
	// those for other delegations. Defaults to 400 for Decisions sessions and 0
	// (request immediately) for other sessions. Omission in a session update preserves
	// the current value. Updates apply to subsequent delegations.
	AcknowledgementDelayMs param.Opt[int64] `json:"acknowledgement_delay_ms,omitzero"`
	// The delegation owner. Always `client` for tasks handled by your application.
	Type constant.Client `json:"type" default:"client"`
	paramObj
}

func (r ClientDelegationParam) MarshalJSON() (data []byte, err error) {
	type shadow ClientDelegationParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ClientDelegationParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func ClientEventParamOfSessionStart(session SessionConfigParam) ClientEventUnionParam {
	var sessionStart SessionStartEventParam
	sessionStart.Session = session
	return ClientEventUnionParam{OfSessionStart: &sessionStart}
}

func ClientEventParamOfSessionUpdate(session SessionUpdateConfigParam) ClientEventUnionParam {
	var sessionUpdate SessionUpdateEventParam
	sessionUpdate.Session = session
	return ClientEventUnionParam{OfSessionUpdate: &sessionUpdate}
}

func ClientEventParamOfSessionInputAudioAppend(audio string) ClientEventUnionParam {
	var sessionInputAudioAppend InputAudioAppendEventParam
	sessionInputAudioAppend.Audio = audio
	return ClientEventUnionParam{OfSessionInputAudioAppend: &sessionInputAudioAppend}
}

func ClientEventParamOfSessionOutputAudioFixturePlay(audio string) ClientEventUnionParam {
	var sessionOutputAudioFixturePlay OutputAudioFixturePlayEventParam
	sessionOutputAudioFixturePlay.Audio = audio
	return ClientEventUnionParam{OfSessionOutputAudioFixturePlay: &sessionOutputAudioFixturePlay}
}

func ClientEventParamOfSessionInstructionsAppend(content string, delegationID string) ClientEventUnionParam {
	var sessionInstructionsAppend InstructionsAppendEventParam
	sessionInstructionsAppend.Content = content
	sessionInstructionsAppend.DelegationID = param.NewOpt(delegationID)
	return ClientEventUnionParam{OfSessionInstructionsAppend: &sessionInstructionsAppend}
}

func ClientEventParamOfSessionThinkingAppend(content string, delegationID string) ClientEventUnionParam {
	var sessionThinkingAppend ThinkingAppendEventParam
	sessionThinkingAppend.Content = content
	sessionThinkingAppend.DelegationID = param.NewOpt(delegationID)
	return ClientEventUnionParam{OfSessionThinkingAppend: &sessionThinkingAppend}
}

func ClientEventParamOfSessionCommentaryAppend(content string, delegationID string) ClientEventUnionParam {
	var sessionCommentaryAppend CommentaryAppendEventParam
	sessionCommentaryAppend.Content = content
	sessionCommentaryAppend.DelegationID = param.NewOpt(delegationID)
	return ClientEventUnionParam{OfSessionCommentaryAppend: &sessionCommentaryAppend}
}

func ClientEventParamOfDecisionCreate(eventID string) ClientEventUnionParam {
	var decisionCreate DecisionCreateEventParam
	decisionCreate.EventID = eventID
	return ClientEventUnionParam{OfDecisionCreate: &decisionCreate}
}

func ClientEventParamOfDecisionItemCreate[
	T DecisionUserMessageParam | DecisionDeveloperMessageParam | DecisionSystemMessageParam | DecisionAssistantMessageParam | DecisionItemCreateEventItemFunctionCallParam | DecisionItemCreateEventItemFunctionCallOutputParam | DecisionItemCreateEventItemCustomToolCallParam | DecisionItemCreateEventItemCustomToolCallOutputParam,
](item T) ClientEventUnionParam {
	var decisionItemCreate DecisionItemCreateEventParam
	switch v := any(item).(type) {
	case DecisionUserMessageParam:
		decisionItemCreate.Item.OfUserMessage = &v
	case DecisionDeveloperMessageParam:
		decisionItemCreate.Item.OfDeveloperMessage = &v
	case DecisionSystemMessageParam:
		decisionItemCreate.Item.OfSystemMessage = &v
	case DecisionAssistantMessageParam:
		decisionItemCreate.Item.OfAssistantMessage = &v
	case DecisionItemCreateEventItemFunctionCallParam:
		decisionItemCreate.Item.OfFunctionCall = &v
	case DecisionItemCreateEventItemFunctionCallOutputParam:
		decisionItemCreate.Item.OfFunctionCallOutput = &v
	case DecisionItemCreateEventItemCustomToolCallParam:
		decisionItemCreate.Item.OfCustomToolCall = &v
	case DecisionItemCreateEventItemCustomToolCallOutputParam:
		decisionItemCreate.Item.OfCustomToolCallOutput = &v
	}
	return ClientEventUnionParam{OfDecisionItemCreate: &decisionItemCreate}
}

func ClientEventParamOfDecisionItemDelete(itemID string) ClientEventUnionParam {
	var decisionItemDelete DecisionItemDeleteEventParam
	decisionItemDelete.ItemID = itemID
	return ClientEventUnionParam{OfDecisionItemDelete: &decisionItemDelete}
}

func ClientEventParamOfResponseItemCreate[
	T responses.EasyInputMessageParam | responses.ResponseInputItemMessageParam | responses.ResponseOutputMessageParam | responses.ResponseFileSearchToolCallParam | responses.ResponseComputerToolCallParam | responses.ResponseInputItemComputerCallOutputParam | responses.ResponseFunctionWebSearchParam | responses.ResponseFunctionToolCallParam | responses.ResponseInputItemFunctionCallOutputParam | responses.ResponseInputItemToolSearchCallParam | responses.ResponseToolSearchOutputItemParam | responses.ResponseInputItemAdditionalToolsParam | responses.ResponseConfigurationUpdateItemParam | responses.ResponseReasoningItemParam | responses.ResponseCompactionItemParam | responses.ResponseInputItemImageGenerationCallParam | responses.ResponseCodeInterpreterToolCallParam | responses.ResponseInputItemLocalShellCallParam | responses.ResponseInputItemLocalShellCallOutputParam | responses.ResponseInputItemShellCallParam | responses.ResponseInputItemShellCallOutputParam | responses.ResponseInputItemApplyPatchCallParam | responses.ResponseInputItemApplyPatchCallOutputParam | responses.ResponseInputItemMcpListToolsParam | responses.ResponseInputItemMcpApprovalRequestParam | responses.ResponseInputItemMcpApprovalResponseParam | responses.ResponseInputItemMcpCallParam | responses.ResponseCustomToolCallOutputParam | responses.ResponseCustomToolCallParam | responses.ResponseInputItemCompactionTriggerParam | responses.ResponseInputItemItemReferenceParam | responses.ResponseInputItemProgramParam | responses.ResponseInputItemProgramOutputParam,
](item T) ClientEventUnionParam {
	var responseItemCreate ResponseItemCreateEventParam
	switch v := any(item).(type) {
	case responses.EasyInputMessageParam:
		responseItemCreate.Item.OfMessage = &v
	case responses.ResponseInputItemMessageParam:
		responseItemCreate.Item.OfInputMessage = &v
	case responses.ResponseOutputMessageParam:
		responseItemCreate.Item.OfOutputMessage = &v
	case responses.ResponseFileSearchToolCallParam:
		responseItemCreate.Item.OfFileSearchCall = &v
	case responses.ResponseComputerToolCallParam:
		responseItemCreate.Item.OfComputerCall = &v
	case responses.ResponseInputItemComputerCallOutputParam:
		responseItemCreate.Item.OfComputerCallOutput = &v
	case responses.ResponseFunctionWebSearchParam:
		responseItemCreate.Item.OfWebSearchCall = &v
	case responses.ResponseFunctionToolCallParam:
		responseItemCreate.Item.OfFunctionCall = &v
	case responses.ResponseInputItemFunctionCallOutputParam:
		responseItemCreate.Item.OfFunctionCallOutput = &v
	case responses.ResponseInputItemToolSearchCallParam:
		responseItemCreate.Item.OfToolSearchCall = &v
	case responses.ResponseToolSearchOutputItemParam:
		responseItemCreate.Item.OfToolSearchOutput = &v
	case responses.ResponseInputItemAdditionalToolsParam:
		responseItemCreate.Item.OfAdditionalTools = &v
	case responses.ResponseConfigurationUpdateItemParam:
		responseItemCreate.Item.OfConfigurationUpdate = &v
	case responses.ResponseReasoningItemParam:
		responseItemCreate.Item.OfReasoning = &v
	case responses.ResponseCompactionItemParam:
		responseItemCreate.Item.OfCompaction = &v
	case responses.ResponseInputItemImageGenerationCallParam:
		responseItemCreate.Item.OfImageGenerationCall = &v
	case responses.ResponseCodeInterpreterToolCallParam:
		responseItemCreate.Item.OfCodeInterpreterCall = &v
	case responses.ResponseInputItemLocalShellCallParam:
		responseItemCreate.Item.OfLocalShellCall = &v
	case responses.ResponseInputItemLocalShellCallOutputParam:
		responseItemCreate.Item.OfLocalShellCallOutput = &v
	case responses.ResponseInputItemShellCallParam:
		responseItemCreate.Item.OfShellCall = &v
	case responses.ResponseInputItemShellCallOutputParam:
		responseItemCreate.Item.OfShellCallOutput = &v
	case responses.ResponseInputItemApplyPatchCallParam:
		responseItemCreate.Item.OfApplyPatchCall = &v
	case responses.ResponseInputItemApplyPatchCallOutputParam:
		responseItemCreate.Item.OfApplyPatchCallOutput = &v
	case responses.ResponseInputItemMcpListToolsParam:
		responseItemCreate.Item.OfMcpListTools = &v
	case responses.ResponseInputItemMcpApprovalRequestParam:
		responseItemCreate.Item.OfMcpApprovalRequest = &v
	case responses.ResponseInputItemMcpApprovalResponseParam:
		responseItemCreate.Item.OfMcpApprovalResponse = &v
	case responses.ResponseInputItemMcpCallParam:
		responseItemCreate.Item.OfMcpCall = &v
	case responses.ResponseCustomToolCallOutputParam:
		responseItemCreate.Item.OfCustomToolCallOutput = &v
	case responses.ResponseCustomToolCallParam:
		responseItemCreate.Item.OfCustomToolCall = &v
	case responses.ResponseInputItemCompactionTriggerParam:
		responseItemCreate.Item.OfCompactionTrigger = &v
	case responses.ResponseInputItemItemReferenceParam:
		responseItemCreate.Item.OfItemReference = &v
	case responses.ResponseInputItemProgramParam:
		responseItemCreate.Item.OfProgram = &v
	case responses.ResponseInputItemProgramOutputParam:
		responseItemCreate.Item.OfProgramOutput = &v
	}
	return ClientEventUnionParam{OfResponseItemCreate: &responseItemCreate}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ClientEventUnionParam struct {
	OfSessionStart                  *SessionStartEventParam           `json:",omitzero,inline"`
	OfSessionUpdate                 *SessionUpdateEventParam          `json:",omitzero,inline"`
	OfSessionInputAudioAppend       *InputAudioAppendEventParam       `json:",omitzero,inline"`
	OfSessionInputAudioMute         *InputAudioMuteEventParam         `json:",omitzero,inline"`
	OfSessionInputAudioUnmute       *InputAudioUnmuteEventParam       `json:",omitzero,inline"`
	OfSessionOutputAudioFixturePlay *OutputAudioFixturePlayEventParam `json:",omitzero,inline"`
	OfSessionInstructionsAppend     *InstructionsAppendEventParam     `json:",omitzero,inline"`
	OfSessionThinkingAppend         *ThinkingAppendEventParam         `json:",omitzero,inline"`
	OfSessionCommentaryAppend       *CommentaryAppendEventParam       `json:",omitzero,inline"`
	OfDecisionCreate                *DecisionCreateEventParam         `json:",omitzero,inline"`
	OfDecisionItemCreate            *DecisionItemCreateEventParam     `json:",omitzero,inline"`
	OfDecisionItemDelete            *DecisionItemDeleteEventParam     `json:",omitzero,inline"`
	OfResponseItemCreate            *ResponseItemCreateEventParam     `json:",omitzero,inline"`
	OfResponseCreate                *ResponseCreateEventParam         `json:",omitzero,inline"`
	OfSessionClose                  *SessionCloseEventParam           `json:",omitzero,inline"`
	paramUnion
}

func (u ClientEventUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfSessionStart,
		u.OfSessionUpdate,
		u.OfSessionInputAudioAppend,
		u.OfSessionInputAudioMute,
		u.OfSessionInputAudioUnmute,
		u.OfSessionOutputAudioFixturePlay,
		u.OfSessionInstructionsAppend,
		u.OfSessionThinkingAppend,
		u.OfSessionCommentaryAppend,
		u.OfDecisionCreate,
		u.OfDecisionItemCreate,
		u.OfDecisionItemDelete,
		u.OfResponseItemCreate,
		u.OfResponseCreate,
		u.OfSessionClose)
}
func (u *ClientEventUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u ClientEventUnionParam) GetItemID() *string {
	if vt := u.OfDecisionItemDelete; vt != nil {
		return &vt.ItemID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ClientEventUnionParam) GetItem() *responses.ResponseInputItemUnionParam {
	if vt := u.OfResponseItemCreate; vt != nil {
		return &vt.Item
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ClientEventUnionParam) GetType() *string {
	if vt := u.OfSessionStart; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfSessionUpdate; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfSessionInputAudioAppend; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfSessionInputAudioMute; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfSessionInputAudioUnmute; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfSessionOutputAudioFixturePlay; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfSessionInstructionsAppend; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfSessionThinkingAppend; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfSessionCommentaryAppend; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfDecisionCreate; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfDecisionItemCreate; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfDecisionItemDelete; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfResponseItemCreate; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfResponseCreate; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfSessionClose; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ClientEventUnionParam) GetEventID() *string {
	if vt := u.OfSessionStart; vt != nil && vt.EventID.Valid() {
		return &vt.EventID.Value
	} else if vt := u.OfSessionUpdate; vt != nil && vt.EventID.Valid() {
		return &vt.EventID.Value
	} else if vt := u.OfSessionInputAudioAppend; vt != nil && vt.EventID.Valid() {
		return &vt.EventID.Value
	} else if vt := u.OfSessionInputAudioMute; vt != nil && vt.EventID.Valid() {
		return &vt.EventID.Value
	} else if vt := u.OfSessionInputAudioUnmute; vt != nil && vt.EventID.Valid() {
		return &vt.EventID.Value
	} else if vt := u.OfSessionOutputAudioFixturePlay; vt != nil && vt.EventID.Valid() {
		return &vt.EventID.Value
	} else if vt := u.OfSessionInstructionsAppend; vt != nil && vt.EventID.Valid() {
		return &vt.EventID.Value
	} else if vt := u.OfSessionThinkingAppend; vt != nil && vt.EventID.Valid() {
		return &vt.EventID.Value
	} else if vt := u.OfSessionCommentaryAppend; vt != nil && vt.EventID.Valid() {
		return &vt.EventID.Value
	} else if vt := u.OfDecisionCreate; vt != nil {
		return (*string)(&vt.EventID)
	} else if vt := u.OfDecisionItemCreate; vt != nil && vt.EventID.Valid() {
		return &vt.EventID.Value
	} else if vt := u.OfDecisionItemDelete; vt != nil && vt.EventID.Valid() {
		return &vt.EventID.Value
	} else if vt := u.OfResponseItemCreate; vt != nil && vt.EventID.Valid() {
		return &vt.EventID.Value
	} else if vt := u.OfResponseCreate; vt != nil && vt.EventID.Valid() {
		return &vt.EventID.Value
	} else if vt := u.OfSessionClose; vt != nil && vt.EventID.Valid() {
		return &vt.EventID.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ClientEventUnionParam) GetAudio() *string {
	if vt := u.OfSessionInputAudioAppend; vt != nil {
		return (*string)(&vt.Audio)
	} else if vt := u.OfSessionOutputAudioFixturePlay; vt != nil {
		return (*string)(&vt.Audio)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ClientEventUnionParam) GetContent() *string {
	if vt := u.OfSessionInstructionsAppend; vt != nil {
		return (*string)(&vt.Content)
	} else if vt := u.OfSessionThinkingAppend; vt != nil {
		return (*string)(&vt.Content)
	} else if vt := u.OfSessionCommentaryAppend; vt != nil {
		return (*string)(&vt.Content)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ClientEventUnionParam) GetDelegationID() *string {
	if vt := u.OfSessionInstructionsAppend; vt != nil && vt.DelegationID.Valid() {
		return &vt.DelegationID.Value
	} else if vt := u.OfSessionThinkingAppend; vt != nil && vt.DelegationID.Valid() {
		return &vt.DelegationID.Value
	} else if vt := u.OfSessionCommentaryAppend; vt != nil && vt.DelegationID.Valid() {
		return &vt.DelegationID.Value
	}
	return nil
}

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u ClientEventUnionParam) GetSession() (res clientEventUnionParamSession) {
	if vt := u.OfSessionStart; vt != nil {
		res.any = &vt.Session
	} else if vt := u.OfSessionUpdate; vt != nil {
		res.any = &vt.Session
	}
	return
}

// Can have the runtime types [*SessionConfigParam], [*SessionUpdateConfigParam]
type clientEventUnionParamSession struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *live.SessionConfigParam:
//	case *live.SessionUpdateConfigParam:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u clientEventUnionParamSession) AsAny() any { return u.any }

// Returns a pointer to the underlying variant's property, if present.
func (u clientEventUnionParamSession) GetModel() *SessionConfigModel {
	switch vt := u.any.(type) {
	case *SessionConfigParam:
		return &vt.Model
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u clientEventUnionParamSession) GetAudio() *SessionConfigAudioParam {
	switch vt := u.any.(type) {
	case *SessionConfigParam:
		return &vt.Audio
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u clientEventUnionParamSession) GetClient() *ClientConfigParam {
	switch vt := u.any.(type) {
	case *SessionConfigParam:
		return &vt.Client
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u clientEventUnionParamSession) GetInput() []InitialItemUnionParam {
	switch vt := u.any.(type) {
	case *SessionConfigParam:
		return vt.Input
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u clientEventUnionParamSession) GetInstructions() *string {
	switch vt := u.any.(type) {
	case *SessionConfigParam:
		return paramutil.AddrIfPresent(vt.Instructions)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u clientEventUnionParamSession) GetStore() *bool {
	switch vt := u.any.(type) {
	case *SessionConfigParam:
		return paramutil.AddrIfPresent(vt.Store)
	}
	return nil
}

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u clientEventUnionParamSession) GetDelegation() (res clientEventUnionParamSessionDelegation) {
	switch vt := u.any.(type) {
	case *SessionConfigParam:
		res.any = vt.Delegation.asAny()
	case *SessionUpdateConfigParam:
		res.any = vt.Delegation.asAny()
	}
	return res
}

// Can have the runtime types [*ClientDelegationParam],
// [*SessionConfigDelegationResponsesParam],
// [*SessionConfigDelegationDecisionsParam],
// [*SessionUpdateConfigDelegationResponsesParam],
// [*SessionUpdateConfigDelegationDecisionsParam]
type clientEventUnionParamSessionDelegation struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *live.ClientDelegationParam:
//	case *live.SessionConfigDelegationResponsesParam:
//	case *live.SessionConfigDelegationDecisionsParam:
//	case *live.SessionUpdateConfigDelegationResponsesParam:
//	case *live.SessionUpdateConfigDelegationDecisionsParam:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u clientEventUnionParamSessionDelegation) AsAny() any { return u.any }

// Returns a pointer to the underlying variant's property, if present.
func (u clientEventUnionParamSessionDelegation) GetType() *string {
	switch vt := u.any.(type) {
	case *ClientDelegationParam:
		return (*string)(&vt.Type)
	case *SessionConfigDelegationResponsesParam:
		return (*string)(&vt.Type)
	case *SessionConfigDelegationDecisionsParam:
		return (*string)(&vt.Type)
	case *SessionUpdateConfigDelegationResponsesParam:
		return (*string)(&vt.Type)
	case *SessionUpdateConfigDelegationDecisionsParam:
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u clientEventUnionParamSessionDelegation) GetAcknowledgementDelayMs() *int64 {
	switch vt := u.any.(type) {
	case *ClientDelegationParam:
		return paramutil.AddrIfPresent(vt.AcknowledgementDelayMs)
	case *SessionConfigDelegationResponsesParam:
		return paramutil.AddrIfPresent(vt.AcknowledgementDelayMs)
	case *SessionConfigDelegationDecisionsParam:
		return paramutil.AddrIfPresent(vt.AcknowledgementDelayMs)
	case *SessionUpdateConfigDelegationResponsesParam:
		return paramutil.AddrIfPresent(vt.AcknowledgementDelayMs)
	case *SessionUpdateConfigDelegationDecisionsParam:
		return paramutil.AddrIfPresent(vt.AcknowledgementDelayMs)
	}
	return nil
}

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u clientEventUnionParamSessionDelegation) GetResponses() (res clientEventUnionParamSessionDelegationResponses) {
	switch vt := u.any.(type) {
	case *SessionConfigDelegationResponsesParam:
		res.any = &vt.Responses
	case *SessionUpdateConfigDelegationResponsesParam:
		res.any = &vt.Responses
	}
	return res
}

// Can have the runtime types [*ResponsesDelegationConfigParam],
// [*ResponsesDelegationUpdateConfigParam]
type clientEventUnionParamSessionDelegationResponses struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *live.ResponsesDelegationConfigParam:
//	case *live.ResponsesDelegationUpdateConfigParam:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u clientEventUnionParamSessionDelegationResponses) AsAny() any { return u.any }

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u clientEventUnionParamSessionDelegation) GetDecisions() (res clientEventUnionParamSessionDelegationDecisions) {
	switch vt := u.any.(type) {
	case *SessionConfigDelegationDecisionsParam:
		res.any = &vt.Decisions
	case *SessionUpdateConfigDelegationDecisionsParam:
		res.any = &vt.Decisions
	}
	return res
}

// Can have the runtime types [*DecisionsDelegationConfigParam],
// [*DecisionsDelegationUpdateConfigParam]
type clientEventUnionParamSessionDelegationDecisions struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *live.DecisionsDelegationConfigParam:
//	case *live.DecisionsDelegationUpdateConfigParam:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u clientEventUnionParamSessionDelegationDecisions) AsAny() any { return u.any }

// Returns a pointer to the underlying variant's Triggers property, if present.
func (u clientEventUnionParamSession) GetTriggers() []DecisionTriggerParam {
	switch vt := u.any.(type) {
	case *SessionConfigParam:
		return vt.Triggers
	case *SessionUpdateConfigParam:
		return vt.Triggers
	}
	return nil
}

func init() {
	apijson.RegisterUnion[ClientEventUnionParam](
		"type",
		apijson.Discriminator[SessionStartEventParam]("session.start"),
		apijson.Discriminator[SessionUpdateEventParam]("session.update"),
		apijson.Discriminator[InputAudioAppendEventParam]("session.input_audio.append"),
		apijson.Discriminator[InputAudioMuteEventParam]("session.input_audio.mute"),
		apijson.Discriminator[InputAudioUnmuteEventParam]("session.input_audio.unmute"),
		apijson.Discriminator[OutputAudioFixturePlayEventParam]("session.output_audio.fixture.play"),
		apijson.Discriminator[InstructionsAppendEventParam]("session.instructions.append"),
		apijson.Discriminator[ThinkingAppendEventParam]("session.thinking.append"),
		apijson.Discriminator[CommentaryAppendEventParam]("session.commentary.append"),
		apijson.Discriminator[DecisionCreateEventParam]("decision.create"),
		apijson.Discriminator[DecisionItemCreateEventParam]("decision.item.create"),
		apijson.Discriminator[DecisionItemDeleteEventParam]("decision.item.delete"),
		apijson.Discriminator[ResponseItemCreateEventParam]("response.item.create"),
		apijson.Discriminator[ResponseCreateEventParam]("response.create"),
		apijson.Discriminator[SessionCloseEventParam]("session.close"),
	)
}

// Provide context the Live model can communicate to the user, optionally for an
// existing client delegation. Rejected while an audio fixture is active.
//
// The properties Content, DelegationID, Type are required.
type CommentaryAppendEventParam struct {
	// Required, nullable. Set null for general session context, or use the ID from
	// session.delegation.created for an existing client delegation. Non-null IDs are
	// not accepted with Responses delegation.
	DelegationID param.Opt[string] `json:"delegation_id,omitzero" api:"required"`
	// Speakable context for the Live model, limited to 500 tokens. Use this for a
	// result the model should communicate; use session.thinking.append for silent
	// context.
	Content string `json:"content" api:"required"`
	// Optional client identifier for correlating this command with a server event's
	// client_event_id or error.client_event_id.
	EventID param.Opt[string] `json:"event_id,omitzero"`
	// The Live client event type. Always `session.commentary.append`.
	//
	// This field can be elided, and will marshal its zero value as
	// "session.commentary.append".
	Type constant.SessionCommentaryAppend `json:"type" default:"session.commentary.append"`
	paramObj
}

func (r CommentaryAppendEventParam) MarshalJSON() (data []byte, err error) {
	type shadow CommentaryAppendEventParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *CommentaryAppendEventParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Returned when a session.commentary.append command is accepted into the Live
// session timeline. Acknowledges the added commentary without guaranteeing exact
// wording or completed audio playback.
type CommentaryAppendedEvent struct {
	// The end of this event on the Live session timeline, in milliseconds from the
	// beginning of the session. For appended context, this can equal start_ms.
	EndMs int64 `json:"end_ms" api:"required"`
	// The unique ID of the Live server event.
	EventID string `json:"event_id" api:"required"`
	// The start of this event on the Live session timeline, in milliseconds from the
	// beginning of the session.
	StartMs int64 `json:"start_ms" api:"required"`
	// The event type, always `session.commentary.appended`.
	Type constant.SessionCommentaryAppended `json:"type" default:"session.commentary.appended"`
	// The event_id of the client command associated with this server event, when
	// supplied.
	ClientEventID string `json:"client_event_id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EndMs         respjson.Field
		EventID       respjson.Field
		StartMs       respjson.Field
		Type          respjson.Field
		ClientEventID respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CommentaryAppendedEvent) RawJSON() string { return r.JSON.raw }
func (r *CommentaryAppendedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type CustomVoice struct {
	ID string `json:"id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CustomVoice) RawJSON() string { return r.JSON.raw }
func (r *CustomVoice) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this CustomVoice to a CustomVoiceParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// CustomVoiceParam.Overrides()
func (r CustomVoice) ToParam() CustomVoiceParam {
	return param.Override[CustomVoiceParam](json.RawMessage(r.RawJSON()))
}

// The property ID is required.
type CustomVoiceParam struct {
	ID string `json:"id" api:"required"`
	paramObj
}

func (r CustomVoiceParam) MarshalJSON() (data []byte, err error) {
	type shadow CustomVoiceParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *CustomVoiceParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Control which Live events an untrusted WebRTC frontend can send and receive over
// its data channel. These restrictions do not apply to trusted sideband
// connections.
type DataChannelConfig struct {
	// Client event types that the frontend data channel may send. Use 'all' to allow
	// every client event; an empty array allows none. Omission preserves the existing
	// allow-all behavior.
	AllowedClientEvents DataChannelConfigAllowedClientEventsUnion `json:"allowed_client_events"`
	// Server events that may be sent to the frontend data channel. Use 'all' to allow
	// every server event; an empty array allows none. Omission preserves the existing
	// allow-all behavior. Responses events use an object with type 'response.event'
	// and a response_event selector.
	AllowedServerEvents DataChannelConfigAllowedServerEventsUnion `json:"allowed_server_events"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		AllowedClientEvents respjson.Field
		AllowedServerEvents respjson.Field
		ExtraFields         map[string]respjson.Field
		raw                 string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DataChannelConfig) RawJSON() string { return r.JSON.raw }
func (r *DataChannelConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this DataChannelConfig to a DataChannelConfigParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// DataChannelConfigParam.Overrides()
func (r DataChannelConfig) ToParam() DataChannelConfigParam {
	return param.Override[DataChannelConfigParam](json.RawMessage(r.RawJSON()))
}

// DataChannelConfigAllowedClientEventsUnion contains all possible properties and
// values from [constant.All], [[]string].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfAll OfArrayOfString]
type DataChannelConfigAllowedClientEventsUnion struct {
	// This field will be present if the value is a [constant.All] instead of an
	// object.
	OfAll constant.All `json:",inline"`
	// This field will be present if the value is a [[]string] instead of an object.
	OfArrayOfString []string `json:",inline"`
	JSON            struct {
		OfAll           respjson.Field
		OfArrayOfString respjson.Field
		raw             string
	} `json:"-"`
}

func (u DataChannelConfigAllowedClientEventsUnion) AsAll() (v constant.All) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u DataChannelConfigAllowedClientEventsUnion) AsArrayOfString() (v []string) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u DataChannelConfigAllowedClientEventsUnion) RawJSON() string { return u.JSON.raw }

func (r *DataChannelConfigAllowedClientEventsUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// DataChannelConfigAllowedServerEventsUnion contains all possible properties and
// values from [constant.All], [[]ServerEventSelector].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfAll OfEventSelectors]
type DataChannelConfigAllowedServerEventsUnion struct {
	// This field will be present if the value is a [constant.All] instead of an
	// object.
	OfAll constant.All `json:",inline"`
	// This field will be present if the value is a [[]ServerEventSelector] instead of
	// an object.
	OfEventSelectors []ServerEventSelector `json:",inline"`
	JSON             struct {
		OfAll            respjson.Field
		OfEventSelectors respjson.Field
		raw              string
	} `json:"-"`
}

func (u DataChannelConfigAllowedServerEventsUnion) AsAll() (v constant.All) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u DataChannelConfigAllowedServerEventsUnion) AsEventSelectors() (v []ServerEventSelector) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u DataChannelConfigAllowedServerEventsUnion) RawJSON() string { return u.JSON.raw }

func (r *DataChannelConfigAllowedServerEventsUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Control which Live events an untrusted WebRTC frontend can send and receive over
// its data channel. These restrictions do not apply to trusted sideband
// connections.
type DataChannelConfigParam struct {
	// Client event types that the frontend data channel may send. Use 'all' to allow
	// every client event; an empty array allows none. Omission preserves the existing
	// allow-all behavior.
	AllowedClientEvents DataChannelConfigAllowedClientEventsUnionParam `json:"allowed_client_events,omitzero"`
	// Server events that may be sent to the frontend data channel. Use 'all' to allow
	// every server event; an empty array allows none. Omission preserves the existing
	// allow-all behavior. Responses events use an object with type 'response.event'
	// and a response_event selector.
	AllowedServerEvents DataChannelConfigAllowedServerEventsUnionParam `json:"allowed_server_events,omitzero"`
	paramObj
}

func (r DataChannelConfigParam) MarshalJSON() (data []byte, err error) {
	type shadow DataChannelConfigParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DataChannelConfigParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DataChannelConfigAllowedClientEventsUnionParam struct {
	// Construct this variant with constant.ValueOf[constant.All]()
	OfAll            constant.All `json:",omitzero,inline"`
	OfArrayOfStrings []string     `json:",omitzero,inline"`
	paramUnion
}

func (u DataChannelConfigAllowedClientEventsUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfAll, u.OfArrayOfStrings)
}
func (u *DataChannelConfigAllowedClientEventsUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DataChannelConfigAllowedServerEventsUnionParam struct {
	// Construct this variant with constant.ValueOf[constant.All]()
	OfAll            constant.All               `json:",omitzero,inline"`
	OfEventSelectors []ServerEventSelectorParam `json:",omitzero,inline"`
	paramUnion
}

func (u DataChannelConfigAllowedServerEventsUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfAll, u.OfEventSelectors)
}
func (u *DataChannelConfigAllowedServerEventsUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

type Decision struct {
	// Answers in question order. Any question can return a refusal instead of an
	// answer; other questions can still receive answers. A refusal has
	// `type: "refusal"` and the question's `name`, or `null` if it was unnamed.
	Answers []DecisionAnswerUnion `json:"answers" api:"required"`
	// The model used.
	Model string `json:"model" api:"required"`
	// Token usage for the request.
	Usage DecisionUsage `json:"usage" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Answers     respjson.Field
		Model       respjson.Field
		Usage       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r Decision) RawJSON() string { return r.JSON.raw }
func (r *Decision) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// DecisionAnswerUnion contains all possible properties and values from
// [DecisionAnswerPredicate], [DecisionAnswerChoice], [DecisionAnswerScore],
// [DecisionAnswerRefusal].
//
// Use the [DecisionAnswerUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type DecisionAnswerUnion struct {
	OfDecisionAnswerPredicate DecisionAnswerPredicate `json:",inline"`
	OfDecisionAnswerChoice    DecisionAnswerChoice    `json:",inline"`
	OfDecisionAnswerScore     DecisionAnswerScore     `json:",inline"`
	OfDecisionAnswerRefusal   DecisionAnswerRefusal   `json:",inline"`
	Type                      string                  `json:"type"`
	JSON                      struct {
		OfDecisionAnswerPredicate respjson.Field
		OfDecisionAnswerChoice    respjson.Field
		OfDecisionAnswerScore     respjson.Field
		OfDecisionAnswerRefusal   respjson.Field
		Type                      respjson.Field
		raw                       string
	} `json:"-"`
}

// anyDecisionAnswer is implemented by each variant of [DecisionAnswerUnion] to add
// type safety for the return type of [DecisionAnswerUnion.AsAny]
type anyDecisionAnswer interface {
	implDecisionAnswerUnion()
}

func (DecisionAnswerPredicate) implDecisionAnswerUnion() {}
func (DecisionAnswerChoice) implDecisionAnswerUnion()    {}
func (DecisionAnswerScore) implDecisionAnswerUnion()     {}
func (DecisionAnswerRefusal) implDecisionAnswerUnion()   {}

// Use the following switch statement to find the correct variant
//
//	switch variant := DecisionAnswerUnion.AsAny().(type) {
//	case live.DecisionAnswerPredicate:
//	case live.DecisionAnswerChoice:
//	case live.DecisionAnswerScore:
//	case live.DecisionAnswerRefusal:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u DecisionAnswerUnion) AsAny() anyDecisionAnswer {
	switch u.Type {
	case "predicate":
		return u.AsDecisionAnswerPredicate()
	case "choice":
		return u.AsDecisionAnswerChoice()
	case "score":
		return u.AsDecisionAnswerScore()
	case "refusal":
		return u.AsDecisionAnswerRefusal()
	}
	return nil
}

func (u DecisionAnswerUnion) AsDecisionAnswerPredicate() (v DecisionAnswerPredicate) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u DecisionAnswerUnion) AsDecisionAnswerChoice() (v DecisionAnswerChoice) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u DecisionAnswerUnion) AsDecisionAnswerScore() (v DecisionAnswerScore) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u DecisionAnswerUnion) AsDecisionAnswerRefusal() (v DecisionAnswerRefusal) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u DecisionAnswerUnion) RawJSON() string { return u.JSON.raw }

func (r *DecisionAnswerUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type DecisionAnswerPredicate struct {
	Name        string  `json:"name" api:"required"`
	Probability float64 `json:"probability" api:"required"`
	// The type of the object. Always `predicate`.
	Type constant.Predicate `json:"type" default:"predicate"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Name        respjson.Field
		Probability respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DecisionAnswerPredicate) RawJSON() string { return r.JSON.raw }
func (r *DecisionAnswerPredicate) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type DecisionAnswerChoice struct {
	// Choice values are typed: a string and a boolean with the same text are distinct.
	Choice        DecisionAnswerChoiceChoiceUnion   `json:"choice" api:"required"`
	Confidence    float64                           `json:"confidence" api:"required"`
	Name          string                            `json:"name" api:"required"`
	Probabilities []DecisionAnswerChoiceProbability `json:"probabilities" api:"required"`
	// The type of the object. Always `choice`.
	Type constant.Choice `json:"type" default:"choice"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Choice        respjson.Field
		Confidence    respjson.Field
		Name          respjson.Field
		Probabilities respjson.Field
		Type          respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DecisionAnswerChoice) RawJSON() string { return r.JSON.raw }
func (r *DecisionAnswerChoice) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// DecisionAnswerChoiceChoiceUnion contains all possible properties and values from
// [string], [bool].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfString OfBool]
type DecisionAnswerChoiceChoiceUnion struct {
	// This field will be present if the value is a [string] instead of an object.
	OfString string `json:",inline"`
	// This field will be present if the value is a [bool] instead of an object.
	OfBool bool `json:",inline"`
	JSON   struct {
		OfString respjson.Field
		OfBool   respjson.Field
		raw      string
	} `json:"-"`
}

func (u DecisionAnswerChoiceChoiceUnion) AsString() (v string) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u DecisionAnswerChoiceChoiceUnion) AsBool() (v bool) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u DecisionAnswerChoiceChoiceUnion) RawJSON() string { return u.JSON.raw }

func (r *DecisionAnswerChoiceChoiceUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type DecisionAnswerChoiceProbability struct {
	Probability float64 `json:"probability" api:"required"`
	// Choice values are typed: a string and a boolean with the same text are distinct.
	Value DecisionAnswerChoiceProbabilityValueUnion `json:"value" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Probability respjson.Field
		Value       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DecisionAnswerChoiceProbability) RawJSON() string { return r.JSON.raw }
func (r *DecisionAnswerChoiceProbability) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// DecisionAnswerChoiceProbabilityValueUnion contains all possible properties and
// values from [string], [bool].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfString OfBool]
type DecisionAnswerChoiceProbabilityValueUnion struct {
	// This field will be present if the value is a [string] instead of an object.
	OfString string `json:",inline"`
	// This field will be present if the value is a [bool] instead of an object.
	OfBool bool `json:",inline"`
	JSON   struct {
		OfString respjson.Field
		OfBool   respjson.Field
		raw      string
	} `json:"-"`
}

func (u DecisionAnswerChoiceProbabilityValueUnion) AsString() (v string) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u DecisionAnswerChoiceProbabilityValueUnion) AsBool() (v bool) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u DecisionAnswerChoiceProbabilityValueUnion) RawJSON() string { return u.JSON.raw }

func (r *DecisionAnswerChoiceProbabilityValueUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type DecisionAnswerScore struct {
	Confidence    float64                          `json:"confidence" api:"required"`
	Name          string                           `json:"name" api:"required"`
	Probabilities []DecisionAnswerScoreProbability `json:"probabilities" api:"required"`
	Score         float64                          `json:"score" api:"required"`
	// The type of the object. Always `score`.
	Type constant.Score `json:"type" default:"score"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Confidence    respjson.Field
		Name          respjson.Field
		Probabilities respjson.Field
		Score         respjson.Field
		Type          respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DecisionAnswerScore) RawJSON() string { return r.JSON.raw }
func (r *DecisionAnswerScore) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type DecisionAnswerScoreProbability struct {
	Label       string  `json:"label" api:"required"`
	Probability float64 `json:"probability" api:"required"`
	Value       int64   `json:"value" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Label       respjson.Field
		Probability respjson.Field
		Value       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DecisionAnswerScoreProbability) RawJSON() string { return r.JSON.raw }
func (r *DecisionAnswerScoreProbability) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The model declined to answer this question. Other questions in the same request
// can still receive answers.
type DecisionAnswerRefusal struct {
	Name string `json:"name" api:"required"`
	// The type of the object. Always `refusal`.
	Type constant.Refusal `json:"type" default:"refusal"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Name        respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DecisionAnswerRefusal) RawJSON() string { return r.JSON.raw }
func (r *DecisionAnswerRefusal) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Token usage for the request.
type DecisionUsage struct {
	InputTokens         int64                            `json:"input_tokens" api:"required"`
	InputTokensDetails  DecisionUsageInputTokensDetails  `json:"input_tokens_details" api:"required"`
	OutputTokens        int64                            `json:"output_tokens" api:"required"`
	OutputTokensDetails DecisionUsageOutputTokensDetails `json:"output_tokens_details" api:"required"`
	TotalTokens         int64                            `json:"total_tokens" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		InputTokens         respjson.Field
		InputTokensDetails  respjson.Field
		OutputTokens        respjson.Field
		OutputTokensDetails respjson.Field
		TotalTokens         respjson.Field
		ExtraFields         map[string]respjson.Field
		raw                 string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DecisionUsage) RawJSON() string { return r.JSON.raw }
func (r *DecisionUsage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type DecisionUsageInputTokensDetails struct {
	CacheWriteTokens int64 `json:"cache_write_tokens" api:"required"`
	CachedTokens     int64 `json:"cached_tokens" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		CacheWriteTokens respjson.Field
		CachedTokens     respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DecisionUsageInputTokensDetails) RawJSON() string { return r.JSON.raw }
func (r *DecisionUsageInputTokensDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type DecisionUsageOutputTokensDetails struct {
	ReasoningTokens int64 `json:"reasoning_tokens" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ReasoningTokens respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DecisionUsageOutputTokensDetails) RawJSON() string { return r.JSON.raw }
func (r *DecisionUsageOutputTokensDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Content, Role are required.
type DecisionAssistantMessageParam struct {
	// The text input to the model.
	Content DecisionAssistantMessageContentUnionParam `json:"content,omitzero" api:"required"`
	// The unique ID of this message item.
	ID param.Opt[string] `json:"id,omitzero"`
	// Labels an `assistant` message as intermediate commentary (`commentary`) or the
	// final answer (`final_answer`). For models like `gpt-5.3-codex` and beyond, when
	// sending follow-up requests, preserve and resend phase on all assistant messages
	// — dropping it can degrade performance. Not used for user messages.
	//
	// Any of "commentary", "final_answer".
	Phase DecisionAssistantMessagePhase `json:"phase,omitzero"`
	// Any of "message".
	Type DecisionAssistantMessageType `json:"type,omitzero"`
	// This field can be elided, and will marshal its zero value as "assistant".
	Role constant.Assistant `json:"role" default:"assistant"`
	paramObj
}

func (r DecisionAssistantMessageParam) MarshalJSON() (data []byte, err error) {
	type shadow DecisionAssistantMessageParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionAssistantMessageParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DecisionAssistantMessageContentUnionParam struct {
	OfString param.Opt[string]                               `json:",omitzero,inline"`
	OfParts  []DecisionAssistantMessageContentPartUnionParam `json:",omitzero,inline"`
	paramUnion
}

func (u DecisionAssistantMessageContentUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfString, u.OfParts)
}
func (u *DecisionAssistantMessageContentUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *DecisionAssistantMessageContentUnionParam) asAny() any {
	if !param.IsOmitted(u.OfString) {
		return &u.OfString.Value
	} else if !param.IsOmitted(u.OfParts) {
		return &u.OfParts
	}
	return nil
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DecisionAssistantMessageContentPartUnionParam struct {
	OfOutputText *DecisionAssistantMessageContentPartOutputTextParam `json:",omitzero,inline"`
	OfRefusal    *DecisionAssistantMessageContentPartRefusalParam    `json:",omitzero,inline"`
	paramUnion
}

func (u DecisionAssistantMessageContentPartUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfOutputText, u.OfRefusal)
}
func (u *DecisionAssistantMessageContentPartUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionAssistantMessageContentPartUnionParam) GetText() *string {
	if vt := u.OfOutputText; vt != nil {
		return &vt.Text
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionAssistantMessageContentPartUnionParam) GetRefusal() *string {
	if vt := u.OfRefusal; vt != nil {
		return &vt.Refusal
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionAssistantMessageContentPartUnionParam) GetType() *string {
	if vt := u.OfOutputText; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfRefusal; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[DecisionAssistantMessageContentPartUnionParam](
		"type",
		apijson.Discriminator[DecisionAssistantMessageContentPartOutputTextParam]("output_text"),
		apijson.Discriminator[DecisionAssistantMessageContentPartRefusalParam]("refusal"),
	)
}

// The properties Text, Type are required.
type DecisionAssistantMessageContentPartOutputTextParam struct {
	// The text content.
	Text string `json:"text" api:"required"`
	// This field can be elided, and will marshal its zero value as "output_text".
	Type constant.OutputText `json:"type" default:"output_text"`
	paramObj
}

func (r DecisionAssistantMessageContentPartOutputTextParam) MarshalJSON() (data []byte, err error) {
	type shadow DecisionAssistantMessageContentPartOutputTextParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionAssistantMessageContentPartOutputTextParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Refusal, Type are required.
type DecisionAssistantMessageContentPartRefusalParam struct {
	// The refusal text.
	Refusal string `json:"refusal" api:"required"`
	// This field can be elided, and will marshal its zero value as "refusal".
	Type constant.Refusal `json:"type" default:"refusal"`
	paramObj
}

func (r DecisionAssistantMessageContentPartRefusalParam) MarshalJSON() (data []byte, err error) {
	type shadow DecisionAssistantMessageContentPartRefusalParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionAssistantMessageContentPartRefusalParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Labels an `assistant` message as intermediate commentary (`commentary`) or the
// final answer (`final_answer`). For models like `gpt-5.3-codex` and beyond, when
// sending follow-up requests, preserve and resend phase on all assistant messages
// — dropping it can degrade performance. Not used for user messages.
type DecisionAssistantMessagePhase string

const (
	DecisionAssistantMessagePhaseCommentary  DecisionAssistantMessagePhase = "commentary"
	DecisionAssistantMessagePhaseFinalAnswer DecisionAssistantMessagePhase = "final_answer"
)

type DecisionAssistantMessageType string

const (
	DecisionAssistantMessageTypeMessage DecisionAssistantMessageType = "message"
)

type DecisionChoiceOption struct {
	Value       DecisionChoiceOptionValueUnion `json:"value" api:"required"`
	Description string                         `json:"description"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Value       respjson.Field
		Description respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DecisionChoiceOption) RawJSON() string { return r.JSON.raw }
func (r *DecisionChoiceOption) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this DecisionChoiceOption to a DecisionChoiceOptionParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// DecisionChoiceOptionParam.Overrides()
func (r DecisionChoiceOption) ToParam() DecisionChoiceOptionParam {
	return param.Override[DecisionChoiceOptionParam](json.RawMessage(r.RawJSON()))
}

// DecisionChoiceOptionValueUnion contains all possible properties and values from
// [string], [bool].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfString OfBool]
type DecisionChoiceOptionValueUnion struct {
	// This field will be present if the value is a [string] instead of an object.
	OfString string `json:",inline"`
	// This field will be present if the value is a [bool] instead of an object.
	OfBool bool `json:",inline"`
	JSON   struct {
		OfString respjson.Field
		OfBool   respjson.Field
		raw      string
	} `json:"-"`
}

func (u DecisionChoiceOptionValueUnion) AsString() (v string) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u DecisionChoiceOptionValueUnion) AsBool() (v bool) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u DecisionChoiceOptionValueUnion) RawJSON() string { return u.JSON.raw }

func (r *DecisionChoiceOptionValueUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The property Value is required.
type DecisionChoiceOptionParam struct {
	Value       DecisionChoiceOptionValueUnionParam `json:"value,omitzero" api:"required"`
	Description param.Opt[string]                   `json:"description,omitzero"`
	paramObj
}

func (r DecisionChoiceOptionParam) MarshalJSON() (data []byte, err error) {
	type shadow DecisionChoiceOptionParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionChoiceOptionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DecisionChoiceOptionValueUnionParam struct {
	OfString param.Opt[string] `json:",omitzero,inline"`
	OfBool   param.Opt[bool]   `json:",omitzero,inline"`
	paramUnion
}

func (u DecisionChoiceOptionValueUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfString, u.OfBool)
}
func (u *DecisionChoiceOptionValueUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Evaluate the configured questions using the current conversation and Decision
// inputs. Reusing an event_id is rejected for the lifetime of the session.
//
// The properties EventID, Type are required.
type DecisionCreateEventParam struct {
	EventID string `json:"event_id" api:"required"`
	// This field can be elided, and will marshal its zero value as "decision.create".
	Type constant.DecisionCreate `json:"type" default:"decision.create"`
	paramObj
}

func (r DecisionCreateEventParam) MarshalJSON() (data []byte, err error) {
	type shadow DecisionCreateEventParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionCreateEventParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Content, Role are required.
type DecisionDeveloperMessageParam struct {
	// The text input to the model.
	Content DecisionDeveloperMessageContentUnionParam `json:"content,omitzero" api:"required"`
	// The unique ID of this message item.
	ID param.Opt[string] `json:"id,omitzero"`
	// Any of "message".
	Type DecisionDeveloperMessageType `json:"type,omitzero"`
	// This field can be elided, and will marshal its zero value as "developer".
	Role constant.Developer `json:"role" default:"developer"`
	paramObj
}

func (r DecisionDeveloperMessageParam) MarshalJSON() (data []byte, err error) {
	type shadow DecisionDeveloperMessageParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionDeveloperMessageParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DecisionDeveloperMessageContentUnionParam struct {
	OfString param.Opt[string]        `json:",omitzero,inline"`
	OfParts  []DecisionInputTextParam `json:",omitzero,inline"`
	paramUnion
}

func (u DecisionDeveloperMessageContentUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfString, u.OfParts)
}
func (u *DecisionDeveloperMessageContentUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *DecisionDeveloperMessageContentUnionParam) asAny() any {
	if !param.IsOmitted(u.OfString) {
		return &u.OfString.Value
	} else if !param.IsOmitted(u.OfParts) {
		return &u.OfParts
	}
	return nil
}

type DecisionDeveloperMessageType string

const (
	DecisionDeveloperMessageTypeMessage DecisionDeveloperMessageType = "message"
)

// A Decision API event for a Live delegation or transcript trigger. For a trigger,
// source_event_id identifies the transcript delta that caused this evaluation.
type DecisionEvent struct {
	// The Live delegation ID, or null for a transcript trigger.
	DelegationID string             `json:"delegation_id" api:"required"`
	Event        DecisionEventEvent `json:"event" api:"required"`
	// The unique ID of this Live server event.
	EventID string `json:"event_id" api:"required"`
	// The event_id of the transcript delta that fired this trigger, or null for a
	// delegation.
	SourceEventID   string                             `json:"source_event_id" api:"required"`
	SourceEventType DecisionEventSourceEventTypeString `json:"source_event_type" api:"required"`
	// The Live-assigned ID of this trigger firing, or null for a delegation.
	TriggerID string                 `json:"trigger_id" api:"required"`
	Type      constant.DecisionEvent `json:"type" default:"decision.event"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		DelegationID    respjson.Field
		Event           respjson.Field
		EventID         respjson.Field
		SourceEventID   respjson.Field
		SourceEventType respjson.Field
		TriggerID       respjson.Field
		Type            respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DecisionEvent) RawJSON() string { return r.JSON.raw }
func (r *DecisionEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type DecisionEventEvent struct {
	// The native Decision API response, including its model, ordered answers, and
	// token usage.
	Decision Decision                   `json:"decision" api:"required"`
	Type     constant.DecisionCompleted `json:"type" default:"decision.completed"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Decision    respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DecisionEventEvent) RawJSON() string { return r.JSON.raw }
func (r *DecisionEventEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type DecisionEventSourceEventTypeString string

type DecisionInputText struct {
	// The text input to the model.
	Text string             `json:"text" api:"required"`
	Type constant.InputText `json:"type" default:"input_text"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Text        respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DecisionInputText) RawJSON() string { return r.JSON.raw }
func (r *DecisionInputText) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this DecisionInputText to a DecisionInputTextParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// DecisionInputTextParam.Overrides()
func (r DecisionInputText) ToParam() DecisionInputTextParam {
	return param.Override[DecisionInputTextParam](json.RawMessage(r.RawJSON()))
}

// The properties Text, Type are required.
type DecisionInputTextParam struct {
	// The text input to the model.
	Text string `json:"text" api:"required"`
	// This field can be elided, and will marshal its zero value as "input_text".
	Type constant.InputText `json:"type" default:"input_text"`
	paramObj
}

func (r DecisionInputTextParam) MarshalJSON() (data []byte, err error) {
	type shadow DecisionInputTextParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionInputTextParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Append an input item for future decisions without starting an evaluation or
// changing Live voice context.
//
// The properties Item, Type are required.
type DecisionItemCreateEventParam struct {
	Item DecisionItemCreateEventItemUnionParam `json:"item,omitzero" api:"required"`
	// Optional client identifier for correlating this command with a server event's
	// client_event_id or error.client_event_id.
	EventID param.Opt[string] `json:"event_id,omitzero"`
	// This field can be elided, and will marshal its zero value as
	// "decision.item.create".
	Type constant.DecisionItemCreate `json:"type" default:"decision.item.create"`
	paramObj
}

func (r DecisionItemCreateEventParam) MarshalJSON() (data []byte, err error) {
	type shadow DecisionItemCreateEventParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionItemCreateEventParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DecisionItemCreateEventItemUnionParam struct {
	OfUserMessage          *DecisionUserMessageParam                             `json:",omitzero,inline"`
	OfDeveloperMessage     *DecisionDeveloperMessageParam                        `json:",omitzero,inline"`
	OfSystemMessage        *DecisionSystemMessageParam                           `json:",omitzero,inline"`
	OfAssistantMessage     *DecisionAssistantMessageParam                        `json:",omitzero,inline"`
	OfFunctionCall         *DecisionItemCreateEventItemFunctionCallParam         `json:",omitzero,inline"`
	OfFunctionCallOutput   *DecisionItemCreateEventItemFunctionCallOutputParam   `json:",omitzero,inline"`
	OfCustomToolCall       *DecisionItemCreateEventItemCustomToolCallParam       `json:",omitzero,inline"`
	OfCustomToolCallOutput *DecisionItemCreateEventItemCustomToolCallOutputParam `json:",omitzero,inline"`
	paramUnion             `default:"message"`
}

func (u DecisionItemCreateEventItemUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfUserMessage,
		u.OfDeveloperMessage,
		u.OfSystemMessage,
		u.OfAssistantMessage,
		u.OfFunctionCall,
		u.OfFunctionCallOutput,
		u.OfCustomToolCall,
		u.OfCustomToolCallOutput)
}
func (u *DecisionItemCreateEventItemUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionItemCreateEventItemUnionParam) GetPhase() *string {
	if vt := u.OfAssistantMessage; vt != nil {
		return (*string)(&vt.Phase)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionItemCreateEventItemUnionParam) GetArguments() *string {
	if vt := u.OfFunctionCall; vt != nil {
		return &vt.Arguments
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionItemCreateEventItemUnionParam) GetInput() *string {
	if vt := u.OfCustomToolCall; vt != nil {
		return &vt.Input
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionItemCreateEventItemUnionParam) GetRole() *string {
	if vt := u.OfUserMessage; vt != nil {
		return (*string)(&vt.Role)
	} else if vt := u.OfDeveloperMessage; vt != nil {
		return (*string)(&vt.Role)
	} else if vt := u.OfSystemMessage; vt != nil {
		return (*string)(&vt.Role)
	} else if vt := u.OfAssistantMessage; vt != nil {
		return (*string)(&vt.Role)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionItemCreateEventItemUnionParam) GetID() *string {
	if vt := u.OfUserMessage; vt != nil && vt.ID.Valid() {
		return &vt.ID.Value
	} else if vt := u.OfDeveloperMessage; vt != nil && vt.ID.Valid() {
		return &vt.ID.Value
	} else if vt := u.OfSystemMessage; vt != nil && vt.ID.Valid() {
		return &vt.ID.Value
	} else if vt := u.OfAssistantMessage; vt != nil && vt.ID.Valid() {
		return &vt.ID.Value
	} else if vt := u.OfFunctionCall; vt != nil && vt.ID.Valid() {
		return &vt.ID.Value
	} else if vt := u.OfFunctionCallOutput; vt != nil && vt.ID.Valid() {
		return &vt.ID.Value
	} else if vt := u.OfCustomToolCall; vt != nil && vt.ID.Valid() {
		return &vt.ID.Value
	} else if vt := u.OfCustomToolCallOutput; vt != nil && vt.ID.Valid() {
		return &vt.ID.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionItemCreateEventItemUnionParam) GetType() *string {
	if vt := u.OfUserMessage; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfDeveloperMessage; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfSystemMessage; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfAssistantMessage; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfFunctionCall; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfFunctionCallOutput; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfCustomToolCall; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfCustomToolCallOutput; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionItemCreateEventItemUnionParam) GetCallID() *string {
	if vt := u.OfFunctionCall; vt != nil {
		return (*string)(&vt.CallID)
	} else if vt := u.OfFunctionCallOutput; vt != nil {
		return (*string)(&vt.CallID)
	} else if vt := u.OfCustomToolCall; vt != nil {
		return (*string)(&vt.CallID)
	} else if vt := u.OfCustomToolCallOutput; vt != nil {
		return (*string)(&vt.CallID)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionItemCreateEventItemUnionParam) GetName() *string {
	if vt := u.OfFunctionCall; vt != nil {
		return (*string)(&vt.Name)
	} else if vt := u.OfCustomToolCall; vt != nil {
		return (*string)(&vt.Name)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionItemCreateEventItemUnionParam) GetNamespace() *string {
	if vt := u.OfFunctionCall; vt != nil && vt.Namespace.Valid() {
		return &vt.Namespace.Value
	} else if vt := u.OfCustomToolCall; vt != nil && vt.Namespace.Valid() {
		return &vt.Namespace.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionItemCreateEventItemUnionParam) GetStatus() *string {
	if vt := u.OfFunctionCall; vt != nil {
		return (*string)(&vt.Status)
	} else if vt := u.OfFunctionCallOutput; vt != nil {
		return (*string)(&vt.Status)
	} else if vt := u.OfCustomToolCall; vt != nil {
		return (*string)(&vt.Status)
	} else if vt := u.OfCustomToolCallOutput; vt != nil {
		return (*string)(&vt.Status)
	}
	return nil
}

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u DecisionItemCreateEventItemUnionParam) GetContent() (res decisionItemCreateEventItemUnionParamContent) {
	if vt := u.OfUserMessage; vt != nil {
		res.any = vt.Content.asAny()
	} else if vt := u.OfDeveloperMessage; vt != nil {
		res.any = vt.Content.asAny()
	} else if vt := u.OfSystemMessage; vt != nil {
		res.any = vt.Content.asAny()
	} else if vt := u.OfAssistantMessage; vt != nil {
		res.any = vt.Content.asAny()
	}
	return
}

// Can have the runtime types [*string],
// [_[]DecisionUserMessageContentPartUnionParam], [_[]DecisionInputTextParam],
// [\*[]DecisionAssistantMessageContentPartUnionParam]
type decisionItemCreateEventItemUnionParamContent struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *string:
//	case *[]live.DecisionUserMessageContentPartUnionParam:
//	case *[]live.DecisionInputTextParam:
//	case *[]live.DecisionAssistantMessageContentPartUnionParam:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u decisionItemCreateEventItemUnionParamContent) AsAny() any { return u.any }

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u DecisionItemCreateEventItemUnionParam) GetOutput() (res decisionItemCreateEventItemUnionParamOutput) {
	if vt := u.OfFunctionCallOutput; vt != nil {
		res.any = vt.Output.asAny()
	} else if vt := u.OfCustomToolCallOutput; vt != nil {
		res.any = vt.Output.asAny()
	}
	return
}

// Can have the runtime types [*string],
// [_[]DecisionItemCreateEventItemFunctionCallOutputOutputPartsItemUnionParam],
// [_[]DecisionItemCreateEventItemCustomToolCallOutputOutputPartsItemUnionParam]
type decisionItemCreateEventItemUnionParamOutput struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *string:
//	case *[]live.DecisionItemCreateEventItemFunctionCallOutputOutputPartsItemUnionParam:
//	case *[]live.DecisionItemCreateEventItemCustomToolCallOutputOutputPartsItemUnionParam:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u decisionItemCreateEventItemUnionParamOutput) AsAny() any { return u.any }

func init() {
	apijson.RegisterUnion[DecisionItemCreateEventItemUnionParam](
		"type",
		apijson.Discriminator[DecisionUserMessageParam]("message"),
		apijson.Discriminator[DecisionDeveloperMessageParam]("message"),
		apijson.Discriminator[DecisionSystemMessageParam]("message"),
		apijson.Discriminator[DecisionAssistantMessageParam]("message"),
		apijson.Discriminator[DecisionItemCreateEventItemFunctionCallParam]("function_call"),
		apijson.Discriminator[DecisionItemCreateEventItemFunctionCallOutputParam]("function_call_output"),
		apijson.Discriminator[DecisionItemCreateEventItemCustomToolCallParam]("custom_tool_call"),
		apijson.Discriminator[DecisionItemCreateEventItemCustomToolCallOutputParam]("custom_tool_call_output"),
		apijson.Discriminator[DecisionUserMessageParam](nil),
		apijson.Discriminator[DecisionDeveloperMessageParam](nil),
		apijson.Discriminator[DecisionSystemMessageParam](nil),
		apijson.Discriminator[DecisionAssistantMessageParam](nil),
	)
}

// The properties Arguments, CallID, Name, Type are required.
type DecisionItemCreateEventItemFunctionCallParam struct {
	// The function arguments as a JSON string.
	Arguments string `json:"arguments" api:"required"`
	// The unique ID of the function tool call generated by the model.
	CallID string `json:"call_id" api:"required"`
	Name   string `json:"name" api:"required"`
	// The local Decision item ID.
	ID param.Opt[string] `json:"id,omitzero"`
	// The namespace of the function to call.
	Namespace param.Opt[string] `json:"namespace,omitzero"`
	// Any of "in_progress", "completed", "incomplete".
	Status string `json:"status,omitzero"`
	// The item type. Always `function_call`.
	//
	// This field can be elided, and will marshal its zero value as "function_call".
	Type constant.FunctionCall `json:"type" default:"function_call"`
	paramObj
}

func (r DecisionItemCreateEventItemFunctionCallParam) MarshalJSON() (data []byte, err error) {
	type shadow DecisionItemCreateEventItemFunctionCallParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionItemCreateEventItemFunctionCallParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[DecisionItemCreateEventItemFunctionCallParam](
		"status", "in_progress", "completed", "incomplete",
	)
}

// The properties CallID, Output, Type are required.
type DecisionItemCreateEventItemFunctionCallOutputParam struct {
	// The unique ID of the function tool call generated by the model.
	CallID string `json:"call_id" api:"required"`
	// The text input to the model.
	Output DecisionItemCreateEventItemFunctionCallOutputOutputUnionParam `json:"output,omitzero" api:"required"`
	// The local Decision item ID.
	ID param.Opt[string] `json:"id,omitzero"`
	// Any of "in_progress", "completed", "incomplete".
	Status string `json:"status,omitzero"`
	// The type of the function tool call output. Always `function_call_output`.
	//
	// This field can be elided, and will marshal its zero value as
	// "function_call_output".
	Type constant.FunctionCallOutput `json:"type" default:"function_call_output"`
	paramObj
}

func (r DecisionItemCreateEventItemFunctionCallOutputParam) MarshalJSON() (data []byte, err error) {
	type shadow DecisionItemCreateEventItemFunctionCallOutputParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionItemCreateEventItemFunctionCallOutputParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[DecisionItemCreateEventItemFunctionCallOutputParam](
		"status", "in_progress", "completed", "incomplete",
	)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DecisionItemCreateEventItemFunctionCallOutputOutputUnionParam struct {
	OfString param.Opt[string]                                                        `json:",omitzero,inline"`
	OfParts  []DecisionItemCreateEventItemFunctionCallOutputOutputPartsItemUnionParam `json:",omitzero,inline"`
	paramUnion
}

func (u DecisionItemCreateEventItemFunctionCallOutputOutputUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfString, u.OfParts)
}
func (u *DecisionItemCreateEventItemFunctionCallOutputOutputUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *DecisionItemCreateEventItemFunctionCallOutputOutputUnionParam) asAny() any {
	if !param.IsOmitted(u.OfString) {
		return &u.OfString.Value
	} else if !param.IsOmitted(u.OfParts) {
		return &u.OfParts
	}
	return nil
}

func DecisionItemCreateEventItemFunctionCallOutputOutputPartsItemParamOfInputText(text string) DecisionItemCreateEventItemFunctionCallOutputOutputPartsItemUnionParam {
	var inputText DecisionInputTextParam
	inputText.Text = text
	return DecisionItemCreateEventItemFunctionCallOutputOutputPartsItemUnionParam{OfInputText: &inputText}
}

func DecisionItemCreateEventItemFunctionCallOutputOutputPartsItemParamOfDecisionItemCreateEventItemFunctionCallOutputOutputPartInputImage(imageURL string) DecisionItemCreateEventItemFunctionCallOutputOutputPartsItemUnionParam {
	var decisionItemCreateEventItemFunctionCallOutputOutputPartInputImage DecisionItemCreateEventItemFunctionCallOutputOutputPartInputImageParam
	decisionItemCreateEventItemFunctionCallOutputOutputPartInputImage.ImageURL = imageURL
	return DecisionItemCreateEventItemFunctionCallOutputOutputPartsItemUnionParam{OfDecisionItemCreateEventItemFunctionCallOutputOutputPartInputImage: &decisionItemCreateEventItemFunctionCallOutputOutputPartInputImage}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DecisionItemCreateEventItemFunctionCallOutputOutputPartsItemUnionParam struct {
	OfInputText                                                         *DecisionInputTextParam                                                 `json:",omitzero,inline"`
	OfDecisionItemCreateEventItemFunctionCallOutputOutputPartInputImage *DecisionItemCreateEventItemFunctionCallOutputOutputPartInputImageParam `json:",omitzero,inline"`
	paramUnion
}

func (u DecisionItemCreateEventItemFunctionCallOutputOutputPartsItemUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfInputText, u.OfDecisionItemCreateEventItemFunctionCallOutputOutputPartInputImage)
}
func (u *DecisionItemCreateEventItemFunctionCallOutputOutputPartsItemUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionItemCreateEventItemFunctionCallOutputOutputPartsItemUnionParam) GetText() *string {
	if vt := u.OfInputText; vt != nil {
		return &vt.Text
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionItemCreateEventItemFunctionCallOutputOutputPartsItemUnionParam) GetImageURL() *string {
	if vt := u.OfDecisionItemCreateEventItemFunctionCallOutputOutputPartInputImage; vt != nil {
		return &vt.ImageURL
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionItemCreateEventItemFunctionCallOutputOutputPartsItemUnionParam) GetDetail() *string {
	if vt := u.OfDecisionItemCreateEventItemFunctionCallOutputOutputPartInputImage; vt != nil {
		return &vt.Detail
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionItemCreateEventItemFunctionCallOutputOutputPartsItemUnionParam) GetType() *string {
	if vt := u.OfInputText; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfDecisionItemCreateEventItemFunctionCallOutputOutputPartInputImage; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[DecisionItemCreateEventItemFunctionCallOutputOutputPartsItemUnionParam](
		"type",
		apijson.Discriminator[DecisionInputTextParam]("input_text"),
		apijson.Discriminator[DecisionItemCreateEventItemFunctionCallOutputOutputPartInputImageParam]("input_image"),
	)
}

// The properties ImageURL, Type are required.
type DecisionItemCreateEventItemFunctionCallOutputOutputPartInputImageParam struct {
	// A base64-encoded image data URL or a publicly accessible HTTP(S) image URL.
	ImageURL string `json:"image_url" api:"required"`
	// The detail level of the image to be sent to the model. One of `high`, `low`,
	// `auto`, or `original`. Defaults to `auto`.
	//
	// Any of "low", "high", "auto", "original".
	Detail string `json:"detail,omitzero"`
	// This field can be elided, and will marshal its zero value as "input_image".
	Type constant.InputImage `json:"type" default:"input_image"`
	paramObj
}

func (r DecisionItemCreateEventItemFunctionCallOutputOutputPartInputImageParam) MarshalJSON() (data []byte, err error) {
	type shadow DecisionItemCreateEventItemFunctionCallOutputOutputPartInputImageParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionItemCreateEventItemFunctionCallOutputOutputPartInputImageParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[DecisionItemCreateEventItemFunctionCallOutputOutputPartInputImageParam](
		"detail", "low", "high", "auto", "original",
	)
}

// The properties CallID, Input, Name, Type are required.
type DecisionItemCreateEventItemCustomToolCallParam struct {
	// The unique ID of the function tool call generated by the model.
	CallID string `json:"call_id" api:"required"`
	// The input to the custom tool call.
	Input string `json:"input" api:"required"`
	// The name of the custom tool to call.
	Name string `json:"name" api:"required"`
	// The local Decision item ID.
	ID param.Opt[string] `json:"id,omitzero"`
	// The namespace of the function to call.
	Namespace param.Opt[string] `json:"namespace,omitzero"`
	// Any of "in_progress", "completed", "incomplete".
	Status string `json:"status,omitzero"`
	// The item type. Always `custom_tool_call`.
	//
	// This field can be elided, and will marshal its zero value as "custom_tool_call".
	Type constant.CustomToolCall `json:"type" default:"custom_tool_call"`
	paramObj
}

func (r DecisionItemCreateEventItemCustomToolCallParam) MarshalJSON() (data []byte, err error) {
	type shadow DecisionItemCreateEventItemCustomToolCallParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionItemCreateEventItemCustomToolCallParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[DecisionItemCreateEventItemCustomToolCallParam](
		"status", "in_progress", "completed", "incomplete",
	)
}

// The properties CallID, Output, Type are required.
type DecisionItemCreateEventItemCustomToolCallOutputParam struct {
	// The unique ID of the function tool call generated by the model.
	CallID string `json:"call_id" api:"required"`
	// The text input to the model.
	Output DecisionItemCreateEventItemCustomToolCallOutputOutputUnionParam `json:"output,omitzero" api:"required"`
	// The local Decision item ID.
	ID param.Opt[string] `json:"id,omitzero"`
	// Any of "in_progress", "completed", "incomplete".
	Status string `json:"status,omitzero"`
	// This field can be elided, and will marshal its zero value as
	// "custom_tool_call_output".
	Type constant.CustomToolCallOutput `json:"type" default:"custom_tool_call_output"`
	paramObj
}

func (r DecisionItemCreateEventItemCustomToolCallOutputParam) MarshalJSON() (data []byte, err error) {
	type shadow DecisionItemCreateEventItemCustomToolCallOutputParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionItemCreateEventItemCustomToolCallOutputParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[DecisionItemCreateEventItemCustomToolCallOutputParam](
		"status", "in_progress", "completed", "incomplete",
	)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DecisionItemCreateEventItemCustomToolCallOutputOutputUnionParam struct {
	OfString param.Opt[string]                                                          `json:",omitzero,inline"`
	OfParts  []DecisionItemCreateEventItemCustomToolCallOutputOutputPartsItemUnionParam `json:",omitzero,inline"`
	paramUnion
}

func (u DecisionItemCreateEventItemCustomToolCallOutputOutputUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfString, u.OfParts)
}
func (u *DecisionItemCreateEventItemCustomToolCallOutputOutputUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *DecisionItemCreateEventItemCustomToolCallOutputOutputUnionParam) asAny() any {
	if !param.IsOmitted(u.OfString) {
		return &u.OfString.Value
	} else if !param.IsOmitted(u.OfParts) {
		return &u.OfParts
	}
	return nil
}

func DecisionItemCreateEventItemCustomToolCallOutputOutputPartsItemParamOfInputText(text string) DecisionItemCreateEventItemCustomToolCallOutputOutputPartsItemUnionParam {
	var inputText DecisionInputTextParam
	inputText.Text = text
	return DecisionItemCreateEventItemCustomToolCallOutputOutputPartsItemUnionParam{OfInputText: &inputText}
}

func DecisionItemCreateEventItemCustomToolCallOutputOutputPartsItemParamOfDecisionItemCreateEventItemCustomToolCallOutputOutputPartInputImage(imageURL string) DecisionItemCreateEventItemCustomToolCallOutputOutputPartsItemUnionParam {
	var decisionItemCreateEventItemCustomToolCallOutputOutputPartInputImage DecisionItemCreateEventItemCustomToolCallOutputOutputPartInputImageParam
	decisionItemCreateEventItemCustomToolCallOutputOutputPartInputImage.ImageURL = imageURL
	return DecisionItemCreateEventItemCustomToolCallOutputOutputPartsItemUnionParam{OfDecisionItemCreateEventItemCustomToolCallOutputOutputPartInputImage: &decisionItemCreateEventItemCustomToolCallOutputOutputPartInputImage}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DecisionItemCreateEventItemCustomToolCallOutputOutputPartsItemUnionParam struct {
	OfInputText                                                           *DecisionInputTextParam                                                   `json:",omitzero,inline"`
	OfDecisionItemCreateEventItemCustomToolCallOutputOutputPartInputImage *DecisionItemCreateEventItemCustomToolCallOutputOutputPartInputImageParam `json:",omitzero,inline"`
	paramUnion
}

func (u DecisionItemCreateEventItemCustomToolCallOutputOutputPartsItemUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfInputText, u.OfDecisionItemCreateEventItemCustomToolCallOutputOutputPartInputImage)
}
func (u *DecisionItemCreateEventItemCustomToolCallOutputOutputPartsItemUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionItemCreateEventItemCustomToolCallOutputOutputPartsItemUnionParam) GetText() *string {
	if vt := u.OfInputText; vt != nil {
		return &vt.Text
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionItemCreateEventItemCustomToolCallOutputOutputPartsItemUnionParam) GetImageURL() *string {
	if vt := u.OfDecisionItemCreateEventItemCustomToolCallOutputOutputPartInputImage; vt != nil {
		return &vt.ImageURL
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionItemCreateEventItemCustomToolCallOutputOutputPartsItemUnionParam) GetDetail() *string {
	if vt := u.OfDecisionItemCreateEventItemCustomToolCallOutputOutputPartInputImage; vt != nil {
		return &vt.Detail
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionItemCreateEventItemCustomToolCallOutputOutputPartsItemUnionParam) GetType() *string {
	if vt := u.OfInputText; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfDecisionItemCreateEventItemCustomToolCallOutputOutputPartInputImage; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[DecisionItemCreateEventItemCustomToolCallOutputOutputPartsItemUnionParam](
		"type",
		apijson.Discriminator[DecisionInputTextParam]("input_text"),
		apijson.Discriminator[DecisionItemCreateEventItemCustomToolCallOutputOutputPartInputImageParam]("input_image"),
	)
}

// The properties ImageURL, Type are required.
type DecisionItemCreateEventItemCustomToolCallOutputOutputPartInputImageParam struct {
	// A base64-encoded image data URL or a publicly accessible HTTP(S) image URL.
	ImageURL string `json:"image_url" api:"required"`
	// The detail level of the image to be sent to the model. One of `high`, `low`,
	// `auto`, or `original`. Defaults to `auto`.
	//
	// Any of "low", "high", "auto", "original".
	Detail string `json:"detail,omitzero"`
	// This field can be elided, and will marshal its zero value as "input_image".
	Type constant.InputImage `json:"type" default:"input_image"`
	paramObj
}

func (r DecisionItemCreateEventItemCustomToolCallOutputOutputPartInputImageParam) MarshalJSON() (data []byte, err error) {
	type shadow DecisionItemCreateEventItemCustomToolCallOutputOutputPartInputImageParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionItemCreateEventItemCustomToolCallOutputOutputPartInputImageParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[DecisionItemCreateEventItemCustomToolCallOutputOutputPartInputImageParam](
		"detail", "low", "high", "auto", "original",
	)
}

// An input item was appended for future decisions.
type DecisionItemCreatedEvent struct {
	// The unique ID of the Live server event.
	EventID string                            `json:"event_id" api:"required"`
	Item    DecisionItemCreatedEventItemUnion `json:"item" api:"required"`
	Type    constant.DecisionItemCreated      `json:"type" default:"decision.item.created"`
	// The event_id of the client command associated with this server event, when
	// supplied.
	ClientEventID string `json:"client_event_id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventID       respjson.Field
		Item          respjson.Field
		Type          respjson.Field
		ClientEventID respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DecisionItemCreatedEvent) RawJSON() string { return r.JSON.raw }
func (r *DecisionItemCreatedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// DecisionItemCreatedEventItemUnion contains all possible properties and values
// from [DecisionStoredUserMessage], [DecisionStoredDeveloperMessage],
// [DecisionStoredSystemMessage], [DecisionStoredAssistantMessage],
// [DecisionItemCreatedEventItemFunctionCall],
// [DecisionItemCreatedEventItemFunctionCallOutput],
// [DecisionItemCreatedEventItemCustomToolCall],
// [DecisionItemCreatedEventItemCustomToolCallOutput].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type DecisionItemCreatedEventItemUnion struct {
	OfDecisionStoredUserMessage DecisionStoredUserMessage `json:",inline"`
	ID                          string                    `json:"id"`
	// This field is a union of [DecisionStoredUserMessageContentUnion],
	// [DecisionStoredDeveloperMessageContentUnion],
	// [DecisionStoredSystemMessageContentUnion],
	// [DecisionStoredAssistantMessageContentUnion]
	Content DecisionItemCreatedEventItemUnionContent `json:"content"`
	Role    string                                   `json:"role"`
	// Any of "message", "message", "message", "message", "function_call",
	// "function_call_output", "custom_tool_call", "custom_tool_call_output".
	Type                             string                         `json:"type"`
	OfDecisionStoredDeveloperMessage DecisionStoredDeveloperMessage `json:",inline"`
	OfDecisionStoredSystemMessage    DecisionStoredSystemMessage    `json:",inline"`
	OfDecisionStoredAssistantMessage DecisionStoredAssistantMessage `json:",inline"`
	// This field is from variant [DecisionStoredAssistantMessage].
	Phase                                      DecisionStoredAssistantMessagePhase      `json:"phase"`
	OfDecisionItemCreatedEventItemFunctionCall DecisionItemCreatedEventItemFunctionCall `json:",inline"`
	// This field is from variant [DecisionItemCreatedEventItemFunctionCall].
	Arguments                                        string                                         `json:"arguments"`
	CallID                                           string                                         `json:"call_id"`
	Name                                             string                                         `json:"name"`
	Namespace                                        string                                         `json:"namespace"`
	Status                                           string                                         `json:"status"`
	OfDecisionItemCreatedEventItemFunctionCallOutput DecisionItemCreatedEventItemFunctionCallOutput `json:",inline"`
	// This field is a union of
	// [DecisionItemCreatedEventItemFunctionCallOutputOutputUnion],
	// [DecisionItemCreatedEventItemCustomToolCallOutputOutputUnion]
	Output                                       DecisionItemCreatedEventItemUnionOutput    `json:"output"`
	OfDecisionItemCreatedEventItemCustomToolCall DecisionItemCreatedEventItemCustomToolCall `json:",inline"`
	// This field is from variant [DecisionItemCreatedEventItemCustomToolCall].
	Input                                              string                                           `json:"input"`
	OfDecisionItemCreatedEventItemCustomToolCallOutput DecisionItemCreatedEventItemCustomToolCallOutput `json:",inline"`
	JSON                                               struct {
		OfDecisionStoredUserMessage                        respjson.Field
		ID                                                 respjson.Field
		Content                                            respjson.Field
		Role                                               respjson.Field
		Type                                               respjson.Field
		OfDecisionStoredDeveloperMessage                   respjson.Field
		OfDecisionStoredSystemMessage                      respjson.Field
		OfDecisionStoredAssistantMessage                   respjson.Field
		Phase                                              respjson.Field
		OfDecisionItemCreatedEventItemFunctionCall         respjson.Field
		Arguments                                          respjson.Field
		CallID                                             respjson.Field
		Name                                               respjson.Field
		Namespace                                          respjson.Field
		Status                                             respjson.Field
		OfDecisionItemCreatedEventItemFunctionCallOutput   respjson.Field
		Output                                             respjson.Field
		OfDecisionItemCreatedEventItemCustomToolCall       respjson.Field
		Input                                              respjson.Field
		OfDecisionItemCreatedEventItemCustomToolCallOutput respjson.Field
		raw                                                string
	} `json:"-"`
}

func (u DecisionItemCreatedEventItemUnion) AsDecisionStoredUserMessage() (v DecisionStoredUserMessage) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u DecisionItemCreatedEventItemUnion) AsDecisionStoredDeveloperMessage() (v DecisionStoredDeveloperMessage) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u DecisionItemCreatedEventItemUnion) AsDecisionStoredSystemMessage() (v DecisionStoredSystemMessage) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u DecisionItemCreatedEventItemUnion) AsDecisionStoredAssistantMessage() (v DecisionStoredAssistantMessage) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u DecisionItemCreatedEventItemUnion) AsDecisionItemCreatedEventItemFunctionCall() (v DecisionItemCreatedEventItemFunctionCall) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u DecisionItemCreatedEventItemUnion) AsDecisionItemCreatedEventItemFunctionCallOutput() (v DecisionItemCreatedEventItemFunctionCallOutput) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u DecisionItemCreatedEventItemUnion) AsDecisionItemCreatedEventItemCustomToolCall() (v DecisionItemCreatedEventItemCustomToolCall) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u DecisionItemCreatedEventItemUnion) AsDecisionItemCreatedEventItemCustomToolCallOutput() (v DecisionItemCreatedEventItemCustomToolCallOutput) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u DecisionItemCreatedEventItemUnion) RawJSON() string { return u.JSON.raw }

func (r *DecisionItemCreatedEventItemUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// DecisionItemCreatedEventItemUnionContent is an implicit subunion of
// [DecisionItemCreatedEventItemUnion]. DecisionItemCreatedEventItemUnionContent
// provides convenient access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [DecisionItemCreatedEventItemUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfString OfParts OfPartsDecisionInputText
// OfPartsDecisionStoredAssistantMessageContentPartUnion]
type DecisionItemCreatedEventItemUnionContent struct {
	// This field will be present if the value is a [string] instead of an object.
	OfString string `json:",inline"`
	// This field will be present if the value is a
	// [[]DecisionStoredUserMessageContentPartUnion] instead of an object.
	OfParts []DecisionStoredUserMessageContentPartUnion `json:",inline"`
	// This field will be present if the value is a [[]DecisionInputText] instead of an
	// object.
	OfPartsDecisionInputText []DecisionInputText `json:",inline"`
	// This field will be present if the value is a
	// [[]DecisionStoredAssistantMessageContentPartUnion] instead of an object.
	OfPartsDecisionStoredAssistantMessageContentPartUnion []DecisionStoredAssistantMessageContentPartUnion `json:",inline"`
	JSON                                                  struct {
		OfString                                              respjson.Field
		OfParts                                               respjson.Field
		OfPartsDecisionInputText                              respjson.Field
		OfPartsDecisionStoredAssistantMessageContentPartUnion respjson.Field
		raw                                                   string
	} `json:"-"`
}

func (r *DecisionItemCreatedEventItemUnionContent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// DecisionItemCreatedEventItemUnionOutput is an implicit subunion of
// [DecisionItemCreatedEventItemUnion]. DecisionItemCreatedEventItemUnionOutput
// provides convenient access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [DecisionItemCreatedEventItemUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfString OfParts
// OfPartsDecisionItemCreatedEventItemCustomToolCallOutputOutputPartsItemUnion]
type DecisionItemCreatedEventItemUnionOutput struct {
	// This field will be present if the value is a [string] instead of an object.
	OfString string `json:",inline"`
	// This field will be present if the value is a
	// [[]DecisionItemCreatedEventItemFunctionCallOutputOutputPartsItemUnion] instead
	// of an object.
	OfParts []DecisionItemCreatedEventItemFunctionCallOutputOutputPartsItemUnion `json:",inline"`
	// This field will be present if the value is a
	// [[]DecisionItemCreatedEventItemCustomToolCallOutputOutputPartsItemUnion] instead
	// of an object.
	OfPartsDecisionItemCreatedEventItemCustomToolCallOutputOutputPartsItemUnion []DecisionItemCreatedEventItemCustomToolCallOutputOutputPartsItemUnion `json:",inline"`
	JSON                                                                        struct {
		OfString                                                                    respjson.Field
		OfParts                                                                     respjson.Field
		OfPartsDecisionItemCreatedEventItemCustomToolCallOutputOutputPartsItemUnion respjson.Field
		raw                                                                         string
	} `json:"-"`
}

func (r *DecisionItemCreatedEventItemUnionOutput) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type DecisionItemCreatedEventItemFunctionCall struct {
	// The function arguments as a JSON string.
	Arguments string `json:"arguments" api:"required"`
	// The unique ID of the function tool call generated by the model.
	CallID string `json:"call_id" api:"required"`
	ID     string `json:"id" api:"required"`
	Name   string `json:"name" api:"required"`
	// The namespace of the function to call.
	Namespace string `json:"namespace" api:"nullable"`
	// Any of "in_progress", "completed", "incomplete".
	Status string `json:"status" api:"nullable"`
	// The item type. Always `function_call`.
	Type constant.FunctionCall `json:"type" default:"function_call"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Arguments   respjson.Field
		CallID      respjson.Field
		ID          respjson.Field
		Name        respjson.Field
		Namespace   respjson.Field
		Status      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DecisionItemCreatedEventItemFunctionCall) RawJSON() string { return r.JSON.raw }
func (r *DecisionItemCreatedEventItemFunctionCall) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type DecisionItemCreatedEventItemFunctionCallOutput struct {
	// The unique ID of the function tool call generated by the model.
	CallID string `json:"call_id" api:"required"`
	ID     string `json:"id" api:"required"`
	// The text input to the model.
	Output DecisionItemCreatedEventItemFunctionCallOutputOutputUnion `json:"output" api:"required"`
	// Any of "in_progress", "completed", "incomplete".
	Status string `json:"status" api:"nullable"`
	// The type of the function tool call output. Always `function_call_output`.
	Type constant.FunctionCallOutput `json:"type" default:"function_call_output"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		CallID      respjson.Field
		ID          respjson.Field
		Output      respjson.Field
		Status      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DecisionItemCreatedEventItemFunctionCallOutput) RawJSON() string { return r.JSON.raw }
func (r *DecisionItemCreatedEventItemFunctionCallOutput) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// DecisionItemCreatedEventItemFunctionCallOutputOutputUnion contains all possible
// properties and values from [string],
// [[]DecisionItemCreatedEventItemFunctionCallOutputOutputPartsItemUnion].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfString OfParts]
type DecisionItemCreatedEventItemFunctionCallOutputOutputUnion struct {
	// This field will be present if the value is a [string] instead of an object.
	OfString string `json:",inline"`
	// This field will be present if the value is a
	// [[]DecisionItemCreatedEventItemFunctionCallOutputOutputPartsItemUnion] instead
	// of an object.
	OfParts []DecisionItemCreatedEventItemFunctionCallOutputOutputPartsItemUnion `json:",inline"`
	JSON    struct {
		OfString respjson.Field
		OfParts  respjson.Field
		raw      string
	} `json:"-"`
}

func (u DecisionItemCreatedEventItemFunctionCallOutputOutputUnion) AsString() (v string) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u DecisionItemCreatedEventItemFunctionCallOutputOutputUnion) AsParts() (v []DecisionItemCreatedEventItemFunctionCallOutputOutputPartsItemUnion) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u DecisionItemCreatedEventItemFunctionCallOutputOutputUnion) RawJSON() string {
	return u.JSON.raw
}

func (r *DecisionItemCreatedEventItemFunctionCallOutputOutputUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// DecisionItemCreatedEventItemFunctionCallOutputOutputPartsItemUnion contains all
// possible properties and values from [DecisionInputText],
// [DecisionItemCreatedEventItemFunctionCallOutputOutputPartInputImage].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type DecisionItemCreatedEventItemFunctionCallOutputOutputPartsItemUnion struct {
	OfDecisionInputText                                                  DecisionInputText                                                  `json:",inline"`
	OfDecisionItemCreatedEventItemFunctionCallOutputOutputPartInputImage DecisionItemCreatedEventItemFunctionCallOutputOutputPartInputImage `json:",inline"`
	JSON                                                                 struct {
		OfDecisionInputText                                                  respjson.Field
		OfDecisionItemCreatedEventItemFunctionCallOutputOutputPartInputImage respjson.Field
		raw                                                                  string
	} `json:"-"`
}

func (u DecisionItemCreatedEventItemFunctionCallOutputOutputPartsItemUnion) AsDecisionInputText() (v DecisionInputText) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u DecisionItemCreatedEventItemFunctionCallOutputOutputPartsItemUnion) AsDecisionItemCreatedEventItemFunctionCallOutputOutputPartInputImage() (v DecisionItemCreatedEventItemFunctionCallOutputOutputPartInputImage) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u DecisionItemCreatedEventItemFunctionCallOutputOutputPartsItemUnion) RawJSON() string {
	return u.JSON.raw
}

func (r *DecisionItemCreatedEventItemFunctionCallOutputOutputPartsItemUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type DecisionItemCreatedEventItemFunctionCallOutputOutputPartInputImage struct {
	// The detail level of the image to be sent to the model. One of `high`, `low`,
	// `auto`, or `original`. Defaults to `auto`.
	//
	// Any of "low", "high", "auto", "original".
	Detail string `json:"detail" api:"nullable"`
	// A base64-encoded image data URL or a publicly accessible HTTP(S) image URL.
	ImageURL string              `json:"image_url" api:"required"`
	Type     constant.InputImage `json:"type" default:"input_image"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Detail      respjson.Field
		ImageURL    respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DecisionItemCreatedEventItemFunctionCallOutputOutputPartInputImage) RawJSON() string {
	return r.JSON.raw
}
func (r *DecisionItemCreatedEventItemFunctionCallOutputOutputPartInputImage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type DecisionItemCreatedEventItemCustomToolCall struct {
	// The unique ID of the function tool call generated by the model.
	CallID string `json:"call_id" api:"required"`
	ID     string `json:"id" api:"required"`
	// The input to the custom tool call.
	Input string `json:"input" api:"required"`
	// The name of the custom tool to call.
	Name string `json:"name" api:"required"`
	// The namespace of the function to call.
	Namespace string `json:"namespace" api:"nullable"`
	// Any of "in_progress", "completed", "incomplete".
	Status string `json:"status" api:"nullable"`
	// The item type. Always `custom_tool_call`.
	Type constant.CustomToolCall `json:"type" default:"custom_tool_call"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		CallID      respjson.Field
		ID          respjson.Field
		Input       respjson.Field
		Name        respjson.Field
		Namespace   respjson.Field
		Status      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DecisionItemCreatedEventItemCustomToolCall) RawJSON() string { return r.JSON.raw }
func (r *DecisionItemCreatedEventItemCustomToolCall) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type DecisionItemCreatedEventItemCustomToolCallOutput struct {
	// The unique ID of the function tool call generated by the model.
	CallID string `json:"call_id" api:"required"`
	ID     string `json:"id" api:"required"`
	// The text input to the model.
	Output DecisionItemCreatedEventItemCustomToolCallOutputOutputUnion `json:"output" api:"required"`
	// Any of "in_progress", "completed", "incomplete".
	Status string                        `json:"status" api:"nullable"`
	Type   constant.CustomToolCallOutput `json:"type" default:"custom_tool_call_output"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		CallID      respjson.Field
		ID          respjson.Field
		Output      respjson.Field
		Status      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DecisionItemCreatedEventItemCustomToolCallOutput) RawJSON() string { return r.JSON.raw }
func (r *DecisionItemCreatedEventItemCustomToolCallOutput) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// DecisionItemCreatedEventItemCustomToolCallOutputOutputUnion contains all
// possible properties and values from [string],
// [[]DecisionItemCreatedEventItemCustomToolCallOutputOutputPartsItemUnion].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfString OfParts]
type DecisionItemCreatedEventItemCustomToolCallOutputOutputUnion struct {
	// This field will be present if the value is a [string] instead of an object.
	OfString string `json:",inline"`
	// This field will be present if the value is a
	// [[]DecisionItemCreatedEventItemCustomToolCallOutputOutputPartsItemUnion] instead
	// of an object.
	OfParts []DecisionItemCreatedEventItemCustomToolCallOutputOutputPartsItemUnion `json:",inline"`
	JSON    struct {
		OfString respjson.Field
		OfParts  respjson.Field
		raw      string
	} `json:"-"`
}

func (u DecisionItemCreatedEventItemCustomToolCallOutputOutputUnion) AsString() (v string) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u DecisionItemCreatedEventItemCustomToolCallOutputOutputUnion) AsParts() (v []DecisionItemCreatedEventItemCustomToolCallOutputOutputPartsItemUnion) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u DecisionItemCreatedEventItemCustomToolCallOutputOutputUnion) RawJSON() string {
	return u.JSON.raw
}

func (r *DecisionItemCreatedEventItemCustomToolCallOutputOutputUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// DecisionItemCreatedEventItemCustomToolCallOutputOutputPartsItemUnion contains
// all possible properties and values from [DecisionInputText],
// [DecisionItemCreatedEventItemCustomToolCallOutputOutputPartInputImage].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type DecisionItemCreatedEventItemCustomToolCallOutputOutputPartsItemUnion struct {
	OfDecisionInputText                                                    DecisionInputText                                                    `json:",inline"`
	OfDecisionItemCreatedEventItemCustomToolCallOutputOutputPartInputImage DecisionItemCreatedEventItemCustomToolCallOutputOutputPartInputImage `json:",inline"`
	JSON                                                                   struct {
		OfDecisionInputText                                                    respjson.Field
		OfDecisionItemCreatedEventItemCustomToolCallOutputOutputPartInputImage respjson.Field
		raw                                                                    string
	} `json:"-"`
}

func (u DecisionItemCreatedEventItemCustomToolCallOutputOutputPartsItemUnion) AsDecisionInputText() (v DecisionInputText) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u DecisionItemCreatedEventItemCustomToolCallOutputOutputPartsItemUnion) AsDecisionItemCreatedEventItemCustomToolCallOutputOutputPartInputImage() (v DecisionItemCreatedEventItemCustomToolCallOutputOutputPartInputImage) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u DecisionItemCreatedEventItemCustomToolCallOutputOutputPartsItemUnion) RawJSON() string {
	return u.JSON.raw
}

func (r *DecisionItemCreatedEventItemCustomToolCallOutputOutputPartsItemUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type DecisionItemCreatedEventItemCustomToolCallOutputOutputPartInputImage struct {
	// The detail level of the image to be sent to the model. One of `high`, `low`,
	// `auto`, or `original`. Defaults to `auto`.
	//
	// Any of "low", "high", "auto", "original".
	Detail string `json:"detail" api:"nullable"`
	// A base64-encoded image data URL or a publicly accessible HTTP(S) image URL.
	ImageURL string              `json:"image_url" api:"required"`
	Type     constant.InputImage `json:"type" default:"input_image"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Detail      respjson.Field
		ImageURL    respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DecisionItemCreatedEventItemCustomToolCallOutputOutputPartInputImage) RawJSON() string {
	return r.JSON.raw
}
func (r *DecisionItemCreatedEventItemCustomToolCallOutputOutputPartInputImage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Delete an input item from future decisions without starting an evaluation or
// changing Live voice context.
//
// The properties ItemID, Type are required.
type DecisionItemDeleteEventParam struct {
	ItemID string `json:"item_id" api:"required"`
	// Optional client identifier for correlating this command with a server event's
	// client_event_id or error.client_event_id.
	EventID param.Opt[string] `json:"event_id,omitzero"`
	// This field can be elided, and will marshal its zero value as
	// "decision.item.delete".
	Type constant.DecisionItemDelete `json:"type" default:"decision.item.delete"`
	paramObj
}

func (r DecisionItemDeleteEventParam) MarshalJSON() (data []byte, err error) {
	type shadow DecisionItemDeleteEventParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionItemDeleteEventParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// An input item was removed from future decisions.
type DecisionItemDeletedEvent struct {
	// The unique ID of the Live server event.
	EventID string                       `json:"event_id" api:"required"`
	ItemID  string                       `json:"item_id" api:"required"`
	Type    constant.DecisionItemDeleted `json:"type" default:"decision.item.deleted"`
	// The event_id of the client command associated with this server event, when
	// supplied.
	ClientEventID string `json:"client_event_id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventID       respjson.Field
		ItemID        respjson.Field
		Type          respjson.Field
		ClientEventID respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DecisionItemDeletedEvent) RawJSON() string { return r.JSON.raw }
func (r *DecisionItemDeletedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type DecisionJSONObject map[string]any

// DecisionQuestionUnion contains all possible properties and values from
// [DecisionQuestionChoice], [DecisionQuestionPredicate], [DecisionQuestionScore].
//
// Use the [DecisionQuestionUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type DecisionQuestionUnion struct {
	OfDecisionQuestionChoice    DecisionQuestionChoice    `json:",inline"`
	OfDecisionQuestionPredicate DecisionQuestionPredicate `json:",inline"`
	OfDecisionQuestionScore     DecisionQuestionScore     `json:",inline"`
	Type                        string                    `json:"type"`
	JSON                        struct {
		OfDecisionQuestionChoice    respjson.Field
		OfDecisionQuestionPredicate respjson.Field
		OfDecisionQuestionScore     respjson.Field
		Type                        respjson.Field
		raw                         string
	} `json:"-"`
}

// anyDecisionQuestion is implemented by each variant of [DecisionQuestionUnion] to
// add type safety for the return type of [DecisionQuestionUnion.AsAny]
type anyDecisionQuestion interface {
	implDecisionQuestionUnion()
}

func (DecisionQuestionChoice) implDecisionQuestionUnion()    {}
func (DecisionQuestionPredicate) implDecisionQuestionUnion() {}
func (DecisionQuestionScore) implDecisionQuestionUnion()     {}

// Use the following switch statement to find the correct variant
//
//	switch variant := DecisionQuestionUnion.AsAny().(type) {
//	case live.DecisionQuestionChoice:
//	case live.DecisionQuestionPredicate:
//	case live.DecisionQuestionScore:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u DecisionQuestionUnion) AsAny() anyDecisionQuestion {
	switch u.Type {
	case "choice":
		return u.AsDecisionQuestionChoice()
	case "predicate":
		return u.AsDecisionQuestionPredicate()
	case "score":
		return u.AsDecisionQuestionScore()
	}
	return nil
}

func (u DecisionQuestionUnion) AsDecisionQuestionChoice() (v DecisionQuestionChoice) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u DecisionQuestionUnion) AsDecisionQuestionPredicate() (v DecisionQuestionPredicate) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u DecisionQuestionUnion) AsDecisionQuestionScore() (v DecisionQuestionScore) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u DecisionQuestionUnion) RawJSON() string { return u.JSON.raw }

func (r *DecisionQuestionUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this DecisionQuestionUnion to a DecisionQuestionUnionParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// DecisionQuestionUnionParam.Overrides()
func (r DecisionQuestionUnion) ToParam() DecisionQuestionUnionParam {
	return param.Override[DecisionQuestionUnionParam](json.RawMessage(r.RawJSON()))
}

type DecisionQuestionChoice struct {
	Choices      []DecisionChoiceOption                  `json:"choices" api:"required"`
	Instructions DecisionQuestionChoiceInstructionsUnion `json:"instructions" api:"required"`
	Name         string                                  `json:"name"`
	Type         constant.Choice                         `json:"type" default:"choice"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Choices      respjson.Field
		Instructions respjson.Field
		Name         respjson.Field
		Type         respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DecisionQuestionChoice) RawJSON() string { return r.JSON.raw }
func (r *DecisionQuestionChoice) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// DecisionQuestionChoiceInstructionsUnion contains all possible properties and
// values from [string], [[]DecisionQuestionChoiceInstructionsPartsItemUnion].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfString OfParts]
type DecisionQuestionChoiceInstructionsUnion struct {
	// This field will be present if the value is a [string] instead of an object.
	OfString string `json:",inline"`
	// This field will be present if the value is a
	// [[]DecisionQuestionChoiceInstructionsPartsItemUnion] instead of an object.
	OfParts []DecisionQuestionChoiceInstructionsPartsItemUnion `json:",inline"`
	JSON    struct {
		OfString respjson.Field
		OfParts  respjson.Field
		raw      string
	} `json:"-"`
}

func (u DecisionQuestionChoiceInstructionsUnion) AsString() (v string) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u DecisionQuestionChoiceInstructionsUnion) AsParts() (v []DecisionQuestionChoiceInstructionsPartsItemUnion) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u DecisionQuestionChoiceInstructionsUnion) RawJSON() string { return u.JSON.raw }

func (r *DecisionQuestionChoiceInstructionsUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// DecisionQuestionChoiceInstructionsPartsItemUnion contains all possible
// properties and values from [DecisionInputText],
// [DecisionQuestionChoiceInstructionsPartInputImage],
// [DecisionQuestionChoiceInstructionsPartInputJSON].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type DecisionQuestionChoiceInstructionsPartsItemUnion struct {
	OfDecisionInputText                                DecisionInputText                                `json:",inline"`
	OfDecisionQuestionChoiceInstructionsPartInputImage DecisionQuestionChoiceInstructionsPartInputImage `json:",inline"`
	OfDecisionQuestionChoiceInstructionsPartInputJSON  DecisionQuestionChoiceInstructionsPartInputJSON  `json:",inline"`
	JSON                                               struct {
		OfDecisionInputText                                respjson.Field
		OfDecisionQuestionChoiceInstructionsPartInputImage respjson.Field
		OfDecisionQuestionChoiceInstructionsPartInputJSON  respjson.Field
		raw                                                string
	} `json:"-"`
}

func (u DecisionQuestionChoiceInstructionsPartsItemUnion) AsDecisionInputText() (v DecisionInputText) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u DecisionQuestionChoiceInstructionsPartsItemUnion) AsDecisionQuestionChoiceInstructionsPartInputImage() (v DecisionQuestionChoiceInstructionsPartInputImage) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u DecisionQuestionChoiceInstructionsPartsItemUnion) AsDecisionQuestionChoiceInstructionsPartInputJSON() (v DecisionQuestionChoiceInstructionsPartInputJSON) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u DecisionQuestionChoiceInstructionsPartsItemUnion) RawJSON() string { return u.JSON.raw }

func (r *DecisionQuestionChoiceInstructionsPartsItemUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type DecisionQuestionChoiceInstructionsPartInputImage struct {
	// The detail level of the image to be sent to the model. One of `high`, `low`,
	// `auto`, or `original`. Defaults to `auto`.
	//
	// Any of "low", "high", "auto", "original".
	Detail string `json:"detail" api:"nullable"`
	// A base64-encoded image data URL or a publicly accessible HTTP(S) image URL.
	ImageURL string              `json:"image_url" api:"required"`
	Type     constant.InputImage `json:"type" default:"input_image"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Detail      respjson.Field
		ImageURL    respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DecisionQuestionChoiceInstructionsPartInputImage) RawJSON() string { return r.JSON.raw }
func (r *DecisionQuestionChoiceInstructionsPartInputImage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type DecisionQuestionChoiceInstructionsPartInputJSON struct {
	Data DecisionQuestionChoiceInstructionsPartInputJSONDataUnion `json:"data" api:"required"`
	Type constant.InputJSON                                       `json:"type" default:"input_json"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Data        respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DecisionQuestionChoiceInstructionsPartInputJSON) RawJSON() string { return r.JSON.raw }
func (r *DecisionQuestionChoiceInstructionsPartInputJSON) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// DecisionQuestionChoiceInstructionsPartInputJSONDataUnion contains all possible
// properties and values from [DecisionJSONObject], [[]DecisionJSONObject].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfArrayOfDecisionJSONObject]
type DecisionQuestionChoiceInstructionsPartInputJSONDataUnion struct {
	OfLiveDecisionJSONObject DecisionJSONObject `json:",inline"`
	// This field will be present if the value is a [[]DecisionJSONObject] instead of
	// an object.
	OfArrayOfDecisionJSONObject []DecisionJSONObject `json:",inline"`
	JSON                        struct {
		OfLiveDecisionJSONObject    respjson.Field
		OfArrayOfDecisionJSONObject respjson.Field
		raw                         string
	} `json:"-"`
}

func (u DecisionQuestionChoiceInstructionsPartInputJSONDataUnion) AsLiveDecisionJSONObject() (v DecisionJSONObject) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u DecisionQuestionChoiceInstructionsPartInputJSONDataUnion) AsArrayOfDecisionJSONObject() (v []DecisionJSONObject) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u DecisionQuestionChoiceInstructionsPartInputJSONDataUnion) RawJSON() string { return u.JSON.raw }

func (r *DecisionQuestionChoiceInstructionsPartInputJSONDataUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type LiveDecisionJSONObject DecisionJSONObject

type DecisionQuestionPredicate struct {
	Instructions DecisionQuestionPredicateInstructionsUnion `json:"instructions" api:"required"`
	Name         string                                     `json:"name"`
	Type         constant.Predicate                         `json:"type" default:"predicate"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Instructions respjson.Field
		Name         respjson.Field
		Type         respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DecisionQuestionPredicate) RawJSON() string { return r.JSON.raw }
func (r *DecisionQuestionPredicate) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// DecisionQuestionPredicateInstructionsUnion contains all possible properties and
// values from [string], [[]DecisionQuestionPredicateInstructionsPartsItemUnion].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfString OfParts]
type DecisionQuestionPredicateInstructionsUnion struct {
	// This field will be present if the value is a [string] instead of an object.
	OfString string `json:",inline"`
	// This field will be present if the value is a
	// [[]DecisionQuestionPredicateInstructionsPartsItemUnion] instead of an object.
	OfParts []DecisionQuestionPredicateInstructionsPartsItemUnion `json:",inline"`
	JSON    struct {
		OfString respjson.Field
		OfParts  respjson.Field
		raw      string
	} `json:"-"`
}

func (u DecisionQuestionPredicateInstructionsUnion) AsString() (v string) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u DecisionQuestionPredicateInstructionsUnion) AsParts() (v []DecisionQuestionPredicateInstructionsPartsItemUnion) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u DecisionQuestionPredicateInstructionsUnion) RawJSON() string { return u.JSON.raw }

func (r *DecisionQuestionPredicateInstructionsUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// DecisionQuestionPredicateInstructionsPartsItemUnion contains all possible
// properties and values from [DecisionInputText],
// [DecisionQuestionPredicateInstructionsPartInputImage],
// [DecisionQuestionPredicateInstructionsPartInputJSON].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type DecisionQuestionPredicateInstructionsPartsItemUnion struct {
	OfDecisionInputText                                   DecisionInputText                                   `json:",inline"`
	OfDecisionQuestionPredicateInstructionsPartInputImage DecisionQuestionPredicateInstructionsPartInputImage `json:",inline"`
	OfDecisionQuestionPredicateInstructionsPartInputJSON  DecisionQuestionPredicateInstructionsPartInputJSON  `json:",inline"`
	JSON                                                  struct {
		OfDecisionInputText                                   respjson.Field
		OfDecisionQuestionPredicateInstructionsPartInputImage respjson.Field
		OfDecisionQuestionPredicateInstructionsPartInputJSON  respjson.Field
		raw                                                   string
	} `json:"-"`
}

func (u DecisionQuestionPredicateInstructionsPartsItemUnion) AsDecisionInputText() (v DecisionInputText) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u DecisionQuestionPredicateInstructionsPartsItemUnion) AsDecisionQuestionPredicateInstructionsPartInputImage() (v DecisionQuestionPredicateInstructionsPartInputImage) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u DecisionQuestionPredicateInstructionsPartsItemUnion) AsDecisionQuestionPredicateInstructionsPartInputJSON() (v DecisionQuestionPredicateInstructionsPartInputJSON) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u DecisionQuestionPredicateInstructionsPartsItemUnion) RawJSON() string { return u.JSON.raw }

func (r *DecisionQuestionPredicateInstructionsPartsItemUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type DecisionQuestionPredicateInstructionsPartInputImage struct {
	// The detail level of the image to be sent to the model. One of `high`, `low`,
	// `auto`, or `original`. Defaults to `auto`.
	//
	// Any of "low", "high", "auto", "original".
	Detail string `json:"detail" api:"nullable"`
	// A base64-encoded image data URL or a publicly accessible HTTP(S) image URL.
	ImageURL string              `json:"image_url" api:"required"`
	Type     constant.InputImage `json:"type" default:"input_image"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Detail      respjson.Field
		ImageURL    respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DecisionQuestionPredicateInstructionsPartInputImage) RawJSON() string { return r.JSON.raw }
func (r *DecisionQuestionPredicateInstructionsPartInputImage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type DecisionQuestionPredicateInstructionsPartInputJSON struct {
	Data DecisionQuestionPredicateInstructionsPartInputJSONDataUnion `json:"data" api:"required"`
	Type constant.InputJSON                                          `json:"type" default:"input_json"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Data        respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DecisionQuestionPredicateInstructionsPartInputJSON) RawJSON() string { return r.JSON.raw }
func (r *DecisionQuestionPredicateInstructionsPartInputJSON) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// DecisionQuestionPredicateInstructionsPartInputJSONDataUnion contains all
// possible properties and values from [DecisionJSONObject],
// [[]DecisionJSONObject].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfArrayOfDecisionJSONObject]
type DecisionQuestionPredicateInstructionsPartInputJSONDataUnion struct {
	OfLiveDecisionJSONObject DecisionJSONObject `json:",inline"`
	// This field will be present if the value is a [[]DecisionJSONObject] instead of
	// an object.
	OfArrayOfDecisionJSONObject []DecisionJSONObject `json:",inline"`
	JSON                        struct {
		OfLiveDecisionJSONObject    respjson.Field
		OfArrayOfDecisionJSONObject respjson.Field
		raw                         string
	} `json:"-"`
}

func (u DecisionQuestionPredicateInstructionsPartInputJSONDataUnion) AsLiveDecisionJSONObject() (v DecisionJSONObject) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u DecisionQuestionPredicateInstructionsPartInputJSONDataUnion) AsArrayOfDecisionJSONObject() (v []DecisionJSONObject) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u DecisionQuestionPredicateInstructionsPartInputJSONDataUnion) RawJSON() string {
	return u.JSON.raw
}

func (r *DecisionQuestionPredicateInstructionsPartInputJSONDataUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type DecisionQuestionScore struct {
	Instructions DecisionQuestionScoreInstructionsUnion `json:"instructions" api:"required"`
	Levels       []DecisionQuestionScoreLevel           `json:"levels" api:"required"`
	Name         string                                 `json:"name"`
	Type         constant.Score                         `json:"type" default:"score"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Instructions respjson.Field
		Levels       respjson.Field
		Name         respjson.Field
		Type         respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DecisionQuestionScore) RawJSON() string { return r.JSON.raw }
func (r *DecisionQuestionScore) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// DecisionQuestionScoreInstructionsUnion contains all possible properties and
// values from [string], [[]DecisionQuestionScoreInstructionsPartsItemUnion].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfString OfParts]
type DecisionQuestionScoreInstructionsUnion struct {
	// This field will be present if the value is a [string] instead of an object.
	OfString string `json:",inline"`
	// This field will be present if the value is a
	// [[]DecisionQuestionScoreInstructionsPartsItemUnion] instead of an object.
	OfParts []DecisionQuestionScoreInstructionsPartsItemUnion `json:",inline"`
	JSON    struct {
		OfString respjson.Field
		OfParts  respjson.Field
		raw      string
	} `json:"-"`
}

func (u DecisionQuestionScoreInstructionsUnion) AsString() (v string) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u DecisionQuestionScoreInstructionsUnion) AsParts() (v []DecisionQuestionScoreInstructionsPartsItemUnion) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u DecisionQuestionScoreInstructionsUnion) RawJSON() string { return u.JSON.raw }

func (r *DecisionQuestionScoreInstructionsUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// DecisionQuestionScoreInstructionsPartsItemUnion contains all possible properties
// and values from [DecisionInputText],
// [DecisionQuestionScoreInstructionsPartInputImage],
// [DecisionQuestionScoreInstructionsPartInputJSON].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type DecisionQuestionScoreInstructionsPartsItemUnion struct {
	OfDecisionInputText                               DecisionInputText                               `json:",inline"`
	OfDecisionQuestionScoreInstructionsPartInputImage DecisionQuestionScoreInstructionsPartInputImage `json:",inline"`
	OfDecisionQuestionScoreInstructionsPartInputJSON  DecisionQuestionScoreInstructionsPartInputJSON  `json:",inline"`
	JSON                                              struct {
		OfDecisionInputText                               respjson.Field
		OfDecisionQuestionScoreInstructionsPartInputImage respjson.Field
		OfDecisionQuestionScoreInstructionsPartInputJSON  respjson.Field
		raw                                               string
	} `json:"-"`
}

func (u DecisionQuestionScoreInstructionsPartsItemUnion) AsDecisionInputText() (v DecisionInputText) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u DecisionQuestionScoreInstructionsPartsItemUnion) AsDecisionQuestionScoreInstructionsPartInputImage() (v DecisionQuestionScoreInstructionsPartInputImage) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u DecisionQuestionScoreInstructionsPartsItemUnion) AsDecisionQuestionScoreInstructionsPartInputJSON() (v DecisionQuestionScoreInstructionsPartInputJSON) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u DecisionQuestionScoreInstructionsPartsItemUnion) RawJSON() string { return u.JSON.raw }

func (r *DecisionQuestionScoreInstructionsPartsItemUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type DecisionQuestionScoreInstructionsPartInputImage struct {
	// The detail level of the image to be sent to the model. One of `high`, `low`,
	// `auto`, or `original`. Defaults to `auto`.
	//
	// Any of "low", "high", "auto", "original".
	Detail string `json:"detail" api:"nullable"`
	// A base64-encoded image data URL or a publicly accessible HTTP(S) image URL.
	ImageURL string              `json:"image_url" api:"required"`
	Type     constant.InputImage `json:"type" default:"input_image"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Detail      respjson.Field
		ImageURL    respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DecisionQuestionScoreInstructionsPartInputImage) RawJSON() string { return r.JSON.raw }
func (r *DecisionQuestionScoreInstructionsPartInputImage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type DecisionQuestionScoreInstructionsPartInputJSON struct {
	Data DecisionQuestionScoreInstructionsPartInputJSONDataUnion `json:"data" api:"required"`
	Type constant.InputJSON                                      `json:"type" default:"input_json"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Data        respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DecisionQuestionScoreInstructionsPartInputJSON) RawJSON() string { return r.JSON.raw }
func (r *DecisionQuestionScoreInstructionsPartInputJSON) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// DecisionQuestionScoreInstructionsPartInputJSONDataUnion contains all possible
// properties and values from [DecisionJSONObject], [[]DecisionJSONObject].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfArrayOfDecisionJSONObject]
type DecisionQuestionScoreInstructionsPartInputJSONDataUnion struct {
	OfLiveDecisionJSONObject DecisionJSONObject `json:",inline"`
	// This field will be present if the value is a [[]DecisionJSONObject] instead of
	// an object.
	OfArrayOfDecisionJSONObject []DecisionJSONObject `json:",inline"`
	JSON                        struct {
		OfLiveDecisionJSONObject    respjson.Field
		OfArrayOfDecisionJSONObject respjson.Field
		raw                         string
	} `json:"-"`
}

func (u DecisionQuestionScoreInstructionsPartInputJSONDataUnion) AsLiveDecisionJSONObject() (v DecisionJSONObject) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u DecisionQuestionScoreInstructionsPartInputJSONDataUnion) AsArrayOfDecisionJSONObject() (v []DecisionJSONObject) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u DecisionQuestionScoreInstructionsPartInputJSONDataUnion) RawJSON() string { return u.JSON.raw }

func (r *DecisionQuestionScoreInstructionsPartInputJSONDataUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type DecisionQuestionScoreLevel struct {
	Label       string `json:"label" api:"required"`
	Description string `json:"description"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Label       respjson.Field
		Description respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DecisionQuestionScoreLevel) RawJSON() string { return r.JSON.raw }
func (r *DecisionQuestionScoreLevel) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func DecisionQuestionParamOfChoice[
	T string | []DecisionQuestionChoiceInstructionsPartsItemUnionParam,
](choices []DecisionChoiceOptionParam, instructions T) DecisionQuestionUnionParam {
	var choice DecisionQuestionChoiceParam
	choice.Choices = choices
	switch v := any(instructions).(type) {
	case string:
		choice.Instructions.OfString = param.NewOpt(v)
	case []DecisionQuestionChoiceInstructionsPartsItemUnionParam:
		choice.Instructions.OfParts = v
	}
	return DecisionQuestionUnionParam{OfChoice: &choice}
}

func DecisionQuestionParamOfPredicate[
	T string | []DecisionQuestionPredicateInstructionsPartsItemUnionParam,
](instructions T) DecisionQuestionUnionParam {
	var predicate DecisionQuestionPredicateParam
	switch v := any(instructions).(type) {
	case string:
		predicate.Instructions.OfString = param.NewOpt(v)
	case []DecisionQuestionPredicateInstructionsPartsItemUnionParam:
		predicate.Instructions.OfParts = v
	}
	return DecisionQuestionUnionParam{OfPredicate: &predicate}
}

func DecisionQuestionParamOfScore[
	T string | []DecisionQuestionScoreInstructionsPartsItemUnionParam,
](instructions T, levels []DecisionQuestionScoreLevelParam) DecisionQuestionUnionParam {
	var score DecisionQuestionScoreParam
	switch v := any(instructions).(type) {
	case string:
		score.Instructions.OfString = param.NewOpt(v)
	case []DecisionQuestionScoreInstructionsPartsItemUnionParam:
		score.Instructions.OfParts = v
	}
	score.Levels = levels
	return DecisionQuestionUnionParam{OfScore: &score}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DecisionQuestionUnionParam struct {
	OfChoice    *DecisionQuestionChoiceParam    `json:",omitzero,inline"`
	OfPredicate *DecisionQuestionPredicateParam `json:",omitzero,inline"`
	OfScore     *DecisionQuestionScoreParam     `json:",omitzero,inline"`
	paramUnion
}

func (u DecisionQuestionUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfChoice, u.OfPredicate, u.OfScore)
}
func (u *DecisionQuestionUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionQuestionUnionParam) GetChoices() []DecisionChoiceOptionParam {
	if vt := u.OfChoice; vt != nil {
		return vt.Choices
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionQuestionUnionParam) GetLevels() []DecisionQuestionScoreLevelParam {
	if vt := u.OfScore; vt != nil {
		return vt.Levels
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionQuestionUnionParam) GetType() *string {
	if vt := u.OfChoice; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfPredicate; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfScore; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionQuestionUnionParam) GetName() *string {
	if vt := u.OfChoice; vt != nil && vt.Name.Valid() {
		return &vt.Name.Value
	} else if vt := u.OfPredicate; vt != nil && vt.Name.Valid() {
		return &vt.Name.Value
	} else if vt := u.OfScore; vt != nil && vt.Name.Valid() {
		return &vt.Name.Value
	}
	return nil
}

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u DecisionQuestionUnionParam) GetInstructions() (res decisionQuestionUnionParamInstructions) {
	if vt := u.OfChoice; vt != nil {
		res.any = vt.Instructions.asAny()
	} else if vt := u.OfPredicate; vt != nil {
		res.any = vt.Instructions.asAny()
	} else if vt := u.OfScore; vt != nil {
		res.any = vt.Instructions.asAny()
	}
	return
}

// Can have the runtime types [*string],
// [_[]DecisionQuestionChoiceInstructionsPartsItemUnionParam],
// [_[]DecisionQuestionPredicateInstructionsPartsItemUnionParam],
// [\*[]DecisionQuestionScoreInstructionsPartsItemUnionParam]
type decisionQuestionUnionParamInstructions struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *string:
//	case *[]live.DecisionQuestionChoiceInstructionsPartsItemUnionParam:
//	case *[]live.DecisionQuestionPredicateInstructionsPartsItemUnionParam:
//	case *[]live.DecisionQuestionScoreInstructionsPartsItemUnionParam:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u decisionQuestionUnionParamInstructions) AsAny() any { return u.any }

func init() {
	apijson.RegisterUnion[DecisionQuestionUnionParam](
		"type",
		apijson.Discriminator[DecisionQuestionChoiceParam]("choice"),
		apijson.Discriminator[DecisionQuestionPredicateParam]("predicate"),
		apijson.Discriminator[DecisionQuestionScoreParam]("score"),
	)
}

// The properties Choices, Instructions, Type are required.
type DecisionQuestionChoiceParam struct {
	Choices      []DecisionChoiceOptionParam                  `json:"choices,omitzero" api:"required"`
	Instructions DecisionQuestionChoiceInstructionsUnionParam `json:"instructions,omitzero" api:"required"`
	Name         param.Opt[string]                            `json:"name,omitzero"`
	// This field can be elided, and will marshal its zero value as "choice".
	Type constant.Choice `json:"type" default:"choice"`
	paramObj
}

func (r DecisionQuestionChoiceParam) MarshalJSON() (data []byte, err error) {
	type shadow DecisionQuestionChoiceParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionQuestionChoiceParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DecisionQuestionChoiceInstructionsUnionParam struct {
	OfString param.Opt[string]                                       `json:",omitzero,inline"`
	OfParts  []DecisionQuestionChoiceInstructionsPartsItemUnionParam `json:",omitzero,inline"`
	paramUnion
}

func (u DecisionQuestionChoiceInstructionsUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfString, u.OfParts)
}
func (u *DecisionQuestionChoiceInstructionsUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *DecisionQuestionChoiceInstructionsUnionParam) asAny() any {
	if !param.IsOmitted(u.OfString) {
		return &u.OfString.Value
	} else if !param.IsOmitted(u.OfParts) {
		return &u.OfParts
	}
	return nil
}

func DecisionQuestionChoiceInstructionsPartsItemParamOfInputText(text string) DecisionQuestionChoiceInstructionsPartsItemUnionParam {
	var inputText DecisionInputTextParam
	inputText.Text = text
	return DecisionQuestionChoiceInstructionsPartsItemUnionParam{OfInputText: &inputText}
}

func DecisionQuestionChoiceInstructionsPartsItemParamOfDecisionQuestionChoiceInstructionsPartInputImage(imageURL string) DecisionQuestionChoiceInstructionsPartsItemUnionParam {
	var decisionQuestionChoiceInstructionsPartInputImage DecisionQuestionChoiceInstructionsPartInputImageParam
	decisionQuestionChoiceInstructionsPartInputImage.ImageURL = imageURL
	return DecisionQuestionChoiceInstructionsPartsItemUnionParam{OfDecisionQuestionChoiceInstructionsPartInputImage: &decisionQuestionChoiceInstructionsPartInputImage}
}

func DecisionQuestionChoiceInstructionsPartsItemParamOfDecisionQuestionChoiceInstructionsPartInputJSON[
	T DecisionJSONObject | []DecisionJSONObject,
](data T) DecisionQuestionChoiceInstructionsPartsItemUnionParam {
	var decisionQuestionChoiceInstructionsPartInputJSON DecisionQuestionChoiceInstructionsPartInputJSONParam
	switch v := any(data).(type) {
	case DecisionJSONObject:
		decisionQuestionChoiceInstructionsPartInputJSON.Data.OfLiveDecisionJSONObject = &v
	case []DecisionJSONObject:
		decisionQuestionChoiceInstructionsPartInputJSON.Data.OfArrayOfJSONObjects = v
	}
	return DecisionQuestionChoiceInstructionsPartsItemUnionParam{OfDecisionQuestionChoiceInstructionsPartInputJSON: &decisionQuestionChoiceInstructionsPartInputJSON}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DecisionQuestionChoiceInstructionsPartsItemUnionParam struct {
	OfInputText                                        *DecisionInputTextParam                                `json:",omitzero,inline"`
	OfDecisionQuestionChoiceInstructionsPartInputImage *DecisionQuestionChoiceInstructionsPartInputImageParam `json:",omitzero,inline"`
	OfDecisionQuestionChoiceInstructionsPartInputJSON  *DecisionQuestionChoiceInstructionsPartInputJSONParam  `json:",omitzero,inline"`
	paramUnion
}

func (u DecisionQuestionChoiceInstructionsPartsItemUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfInputText, u.OfDecisionQuestionChoiceInstructionsPartInputImage, u.OfDecisionQuestionChoiceInstructionsPartInputJSON)
}
func (u *DecisionQuestionChoiceInstructionsPartsItemUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionQuestionChoiceInstructionsPartsItemUnionParam) GetText() *string {
	if vt := u.OfInputText; vt != nil {
		return &vt.Text
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionQuestionChoiceInstructionsPartsItemUnionParam) GetImageURL() *string {
	if vt := u.OfDecisionQuestionChoiceInstructionsPartInputImage; vt != nil {
		return &vt.ImageURL
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionQuestionChoiceInstructionsPartsItemUnionParam) GetDetail() *string {
	if vt := u.OfDecisionQuestionChoiceInstructionsPartInputImage; vt != nil {
		return &vt.Detail
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionQuestionChoiceInstructionsPartsItemUnionParam) GetData() *DecisionQuestionChoiceInstructionsPartInputJSONDataUnionParam {
	if vt := u.OfDecisionQuestionChoiceInstructionsPartInputJSON; vt != nil {
		return &vt.Data
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionQuestionChoiceInstructionsPartsItemUnionParam) GetType() *string {
	if vt := u.OfInputText; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfDecisionQuestionChoiceInstructionsPartInputImage; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfDecisionQuestionChoiceInstructionsPartInputJSON; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[DecisionQuestionChoiceInstructionsPartsItemUnionParam](
		"type",
		apijson.Discriminator[DecisionInputTextParam]("input_text"),
		apijson.Discriminator[DecisionQuestionChoiceInstructionsPartInputImageParam]("input_image"),
		apijson.Discriminator[DecisionQuestionChoiceInstructionsPartInputJSONParam]("input_json"),
	)
}

// The properties ImageURL, Type are required.
type DecisionQuestionChoiceInstructionsPartInputImageParam struct {
	// A base64-encoded image data URL or a publicly accessible HTTP(S) image URL.
	ImageURL string `json:"image_url" api:"required"`
	// The detail level of the image to be sent to the model. One of `high`, `low`,
	// `auto`, or `original`. Defaults to `auto`.
	//
	// Any of "low", "high", "auto", "original".
	Detail string `json:"detail,omitzero"`
	// This field can be elided, and will marshal its zero value as "input_image".
	Type constant.InputImage `json:"type" default:"input_image"`
	paramObj
}

func (r DecisionQuestionChoiceInstructionsPartInputImageParam) MarshalJSON() (data []byte, err error) {
	type shadow DecisionQuestionChoiceInstructionsPartInputImageParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionQuestionChoiceInstructionsPartInputImageParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[DecisionQuestionChoiceInstructionsPartInputImageParam](
		"detail", "low", "high", "auto", "original",
	)
}

// The properties Data, Type are required.
type DecisionQuestionChoiceInstructionsPartInputJSONParam struct {
	Data DecisionQuestionChoiceInstructionsPartInputJSONDataUnionParam `json:"data,omitzero" api:"required"`
	// This field can be elided, and will marshal its zero value as "input_json".
	Type constant.InputJSON `json:"type" default:"input_json"`
	paramObj
}

func (r DecisionQuestionChoiceInstructionsPartInputJSONParam) MarshalJSON() (data []byte, err error) {
	type shadow DecisionQuestionChoiceInstructionsPartInputJSONParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionQuestionChoiceInstructionsPartInputJSONParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DecisionQuestionChoiceInstructionsPartInputJSONDataUnionParam struct {
	OfLiveDecisionJSONObject *DecisionJSONObject  `json:",omitzero,inline"`
	OfArrayOfJSONObjects     []DecisionJSONObject `json:",omitzero,inline"`
	paramUnion
}

func (u DecisionQuestionChoiceInstructionsPartInputJSONDataUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfLiveDecisionJSONObject, u.OfArrayOfJSONObjects)
}
func (u *DecisionQuestionChoiceInstructionsPartInputJSONDataUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// The properties Instructions, Type are required.
type DecisionQuestionPredicateParam struct {
	Instructions DecisionQuestionPredicateInstructionsUnionParam `json:"instructions,omitzero" api:"required"`
	Name         param.Opt[string]                               `json:"name,omitzero"`
	// This field can be elided, and will marshal its zero value as "predicate".
	Type constant.Predicate `json:"type" default:"predicate"`
	paramObj
}

func (r DecisionQuestionPredicateParam) MarshalJSON() (data []byte, err error) {
	type shadow DecisionQuestionPredicateParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionQuestionPredicateParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DecisionQuestionPredicateInstructionsUnionParam struct {
	OfString param.Opt[string]                                          `json:",omitzero,inline"`
	OfParts  []DecisionQuestionPredicateInstructionsPartsItemUnionParam `json:",omitzero,inline"`
	paramUnion
}

func (u DecisionQuestionPredicateInstructionsUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfString, u.OfParts)
}
func (u *DecisionQuestionPredicateInstructionsUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *DecisionQuestionPredicateInstructionsUnionParam) asAny() any {
	if !param.IsOmitted(u.OfString) {
		return &u.OfString.Value
	} else if !param.IsOmitted(u.OfParts) {
		return &u.OfParts
	}
	return nil
}

func DecisionQuestionPredicateInstructionsPartsItemParamOfInputText(text string) DecisionQuestionPredicateInstructionsPartsItemUnionParam {
	var inputText DecisionInputTextParam
	inputText.Text = text
	return DecisionQuestionPredicateInstructionsPartsItemUnionParam{OfInputText: &inputText}
}

func DecisionQuestionPredicateInstructionsPartsItemParamOfDecisionQuestionPredicateInstructionsPartInputImage(imageURL string) DecisionQuestionPredicateInstructionsPartsItemUnionParam {
	var decisionQuestionPredicateInstructionsPartInputImage DecisionQuestionPredicateInstructionsPartInputImageParam
	decisionQuestionPredicateInstructionsPartInputImage.ImageURL = imageURL
	return DecisionQuestionPredicateInstructionsPartsItemUnionParam{OfDecisionQuestionPredicateInstructionsPartInputImage: &decisionQuestionPredicateInstructionsPartInputImage}
}

func DecisionQuestionPredicateInstructionsPartsItemParamOfDecisionQuestionPredicateInstructionsPartInputJSON[
	T DecisionJSONObject | []DecisionJSONObject,
](data T) DecisionQuestionPredicateInstructionsPartsItemUnionParam {
	var decisionQuestionPredicateInstructionsPartInputJSON DecisionQuestionPredicateInstructionsPartInputJSONParam
	switch v := any(data).(type) {
	case DecisionJSONObject:
		decisionQuestionPredicateInstructionsPartInputJSON.Data.OfLiveDecisionJSONObject = &v
	case []DecisionJSONObject:
		decisionQuestionPredicateInstructionsPartInputJSON.Data.OfArrayOfJSONObjects = v
	}
	return DecisionQuestionPredicateInstructionsPartsItemUnionParam{OfDecisionQuestionPredicateInstructionsPartInputJSON: &decisionQuestionPredicateInstructionsPartInputJSON}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DecisionQuestionPredicateInstructionsPartsItemUnionParam struct {
	OfInputText                                           *DecisionInputTextParam                                   `json:",omitzero,inline"`
	OfDecisionQuestionPredicateInstructionsPartInputImage *DecisionQuestionPredicateInstructionsPartInputImageParam `json:",omitzero,inline"`
	OfDecisionQuestionPredicateInstructionsPartInputJSON  *DecisionQuestionPredicateInstructionsPartInputJSONParam  `json:",omitzero,inline"`
	paramUnion
}

func (u DecisionQuestionPredicateInstructionsPartsItemUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfInputText, u.OfDecisionQuestionPredicateInstructionsPartInputImage, u.OfDecisionQuestionPredicateInstructionsPartInputJSON)
}
func (u *DecisionQuestionPredicateInstructionsPartsItemUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionQuestionPredicateInstructionsPartsItemUnionParam) GetText() *string {
	if vt := u.OfInputText; vt != nil {
		return &vt.Text
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionQuestionPredicateInstructionsPartsItemUnionParam) GetImageURL() *string {
	if vt := u.OfDecisionQuestionPredicateInstructionsPartInputImage; vt != nil {
		return &vt.ImageURL
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionQuestionPredicateInstructionsPartsItemUnionParam) GetDetail() *string {
	if vt := u.OfDecisionQuestionPredicateInstructionsPartInputImage; vt != nil {
		return &vt.Detail
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionQuestionPredicateInstructionsPartsItemUnionParam) GetData() *DecisionQuestionPredicateInstructionsPartInputJSONDataUnionParam {
	if vt := u.OfDecisionQuestionPredicateInstructionsPartInputJSON; vt != nil {
		return &vt.Data
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionQuestionPredicateInstructionsPartsItemUnionParam) GetType() *string {
	if vt := u.OfInputText; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfDecisionQuestionPredicateInstructionsPartInputImage; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfDecisionQuestionPredicateInstructionsPartInputJSON; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[DecisionQuestionPredicateInstructionsPartsItemUnionParam](
		"type",
		apijson.Discriminator[DecisionInputTextParam]("input_text"),
		apijson.Discriminator[DecisionQuestionPredicateInstructionsPartInputImageParam]("input_image"),
		apijson.Discriminator[DecisionQuestionPredicateInstructionsPartInputJSONParam]("input_json"),
	)
}

// The properties ImageURL, Type are required.
type DecisionQuestionPredicateInstructionsPartInputImageParam struct {
	// A base64-encoded image data URL or a publicly accessible HTTP(S) image URL.
	ImageURL string `json:"image_url" api:"required"`
	// The detail level of the image to be sent to the model. One of `high`, `low`,
	// `auto`, or `original`. Defaults to `auto`.
	//
	// Any of "low", "high", "auto", "original".
	Detail string `json:"detail,omitzero"`
	// This field can be elided, and will marshal its zero value as "input_image".
	Type constant.InputImage `json:"type" default:"input_image"`
	paramObj
}

func (r DecisionQuestionPredicateInstructionsPartInputImageParam) MarshalJSON() (data []byte, err error) {
	type shadow DecisionQuestionPredicateInstructionsPartInputImageParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionQuestionPredicateInstructionsPartInputImageParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[DecisionQuestionPredicateInstructionsPartInputImageParam](
		"detail", "low", "high", "auto", "original",
	)
}

// The properties Data, Type are required.
type DecisionQuestionPredicateInstructionsPartInputJSONParam struct {
	Data DecisionQuestionPredicateInstructionsPartInputJSONDataUnionParam `json:"data,omitzero" api:"required"`
	// This field can be elided, and will marshal its zero value as "input_json".
	Type constant.InputJSON `json:"type" default:"input_json"`
	paramObj
}

func (r DecisionQuestionPredicateInstructionsPartInputJSONParam) MarshalJSON() (data []byte, err error) {
	type shadow DecisionQuestionPredicateInstructionsPartInputJSONParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionQuestionPredicateInstructionsPartInputJSONParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DecisionQuestionPredicateInstructionsPartInputJSONDataUnionParam struct {
	OfLiveDecisionJSONObject *DecisionJSONObject  `json:",omitzero,inline"`
	OfArrayOfJSONObjects     []DecisionJSONObject `json:",omitzero,inline"`
	paramUnion
}

func (u DecisionQuestionPredicateInstructionsPartInputJSONDataUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfLiveDecisionJSONObject, u.OfArrayOfJSONObjects)
}
func (u *DecisionQuestionPredicateInstructionsPartInputJSONDataUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// The properties Instructions, Levels, Type are required.
type DecisionQuestionScoreParam struct {
	Instructions DecisionQuestionScoreInstructionsUnionParam `json:"instructions,omitzero" api:"required"`
	Levels       []DecisionQuestionScoreLevelParam           `json:"levels,omitzero" api:"required"`
	Name         param.Opt[string]                           `json:"name,omitzero"`
	// This field can be elided, and will marshal its zero value as "score".
	Type constant.Score `json:"type" default:"score"`
	paramObj
}

func (r DecisionQuestionScoreParam) MarshalJSON() (data []byte, err error) {
	type shadow DecisionQuestionScoreParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionQuestionScoreParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DecisionQuestionScoreInstructionsUnionParam struct {
	OfString param.Opt[string]                                      `json:",omitzero,inline"`
	OfParts  []DecisionQuestionScoreInstructionsPartsItemUnionParam `json:",omitzero,inline"`
	paramUnion
}

func (u DecisionQuestionScoreInstructionsUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfString, u.OfParts)
}
func (u *DecisionQuestionScoreInstructionsUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *DecisionQuestionScoreInstructionsUnionParam) asAny() any {
	if !param.IsOmitted(u.OfString) {
		return &u.OfString.Value
	} else if !param.IsOmitted(u.OfParts) {
		return &u.OfParts
	}
	return nil
}

func DecisionQuestionScoreInstructionsPartsItemParamOfInputText(text string) DecisionQuestionScoreInstructionsPartsItemUnionParam {
	var inputText DecisionInputTextParam
	inputText.Text = text
	return DecisionQuestionScoreInstructionsPartsItemUnionParam{OfInputText: &inputText}
}

func DecisionQuestionScoreInstructionsPartsItemParamOfDecisionQuestionScoreInstructionsPartInputImage(imageURL string) DecisionQuestionScoreInstructionsPartsItemUnionParam {
	var decisionQuestionScoreInstructionsPartInputImage DecisionQuestionScoreInstructionsPartInputImageParam
	decisionQuestionScoreInstructionsPartInputImage.ImageURL = imageURL
	return DecisionQuestionScoreInstructionsPartsItemUnionParam{OfDecisionQuestionScoreInstructionsPartInputImage: &decisionQuestionScoreInstructionsPartInputImage}
}

func DecisionQuestionScoreInstructionsPartsItemParamOfDecisionQuestionScoreInstructionsPartInputJSON[
	T DecisionJSONObject | []DecisionJSONObject,
](data T) DecisionQuestionScoreInstructionsPartsItemUnionParam {
	var decisionQuestionScoreInstructionsPartInputJSON DecisionQuestionScoreInstructionsPartInputJSONParam
	switch v := any(data).(type) {
	case DecisionJSONObject:
		decisionQuestionScoreInstructionsPartInputJSON.Data.OfLiveDecisionJSONObject = &v
	case []DecisionJSONObject:
		decisionQuestionScoreInstructionsPartInputJSON.Data.OfArrayOfJSONObjects = v
	}
	return DecisionQuestionScoreInstructionsPartsItemUnionParam{OfDecisionQuestionScoreInstructionsPartInputJSON: &decisionQuestionScoreInstructionsPartInputJSON}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DecisionQuestionScoreInstructionsPartsItemUnionParam struct {
	OfInputText                                       *DecisionInputTextParam                               `json:",omitzero,inline"`
	OfDecisionQuestionScoreInstructionsPartInputImage *DecisionQuestionScoreInstructionsPartInputImageParam `json:",omitzero,inline"`
	OfDecisionQuestionScoreInstructionsPartInputJSON  *DecisionQuestionScoreInstructionsPartInputJSONParam  `json:",omitzero,inline"`
	paramUnion
}

func (u DecisionQuestionScoreInstructionsPartsItemUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfInputText, u.OfDecisionQuestionScoreInstructionsPartInputImage, u.OfDecisionQuestionScoreInstructionsPartInputJSON)
}
func (u *DecisionQuestionScoreInstructionsPartsItemUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionQuestionScoreInstructionsPartsItemUnionParam) GetText() *string {
	if vt := u.OfInputText; vt != nil {
		return &vt.Text
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionQuestionScoreInstructionsPartsItemUnionParam) GetImageURL() *string {
	if vt := u.OfDecisionQuestionScoreInstructionsPartInputImage; vt != nil {
		return &vt.ImageURL
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionQuestionScoreInstructionsPartsItemUnionParam) GetDetail() *string {
	if vt := u.OfDecisionQuestionScoreInstructionsPartInputImage; vt != nil {
		return &vt.Detail
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionQuestionScoreInstructionsPartsItemUnionParam) GetData() *DecisionQuestionScoreInstructionsPartInputJSONDataUnionParam {
	if vt := u.OfDecisionQuestionScoreInstructionsPartInputJSON; vt != nil {
		return &vt.Data
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionQuestionScoreInstructionsPartsItemUnionParam) GetType() *string {
	if vt := u.OfInputText; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfDecisionQuestionScoreInstructionsPartInputImage; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfDecisionQuestionScoreInstructionsPartInputJSON; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[DecisionQuestionScoreInstructionsPartsItemUnionParam](
		"type",
		apijson.Discriminator[DecisionInputTextParam]("input_text"),
		apijson.Discriminator[DecisionQuestionScoreInstructionsPartInputImageParam]("input_image"),
		apijson.Discriminator[DecisionQuestionScoreInstructionsPartInputJSONParam]("input_json"),
	)
}

// The properties ImageURL, Type are required.
type DecisionQuestionScoreInstructionsPartInputImageParam struct {
	// A base64-encoded image data URL or a publicly accessible HTTP(S) image URL.
	ImageURL string `json:"image_url" api:"required"`
	// The detail level of the image to be sent to the model. One of `high`, `low`,
	// `auto`, or `original`. Defaults to `auto`.
	//
	// Any of "low", "high", "auto", "original".
	Detail string `json:"detail,omitzero"`
	// This field can be elided, and will marshal its zero value as "input_image".
	Type constant.InputImage `json:"type" default:"input_image"`
	paramObj
}

func (r DecisionQuestionScoreInstructionsPartInputImageParam) MarshalJSON() (data []byte, err error) {
	type shadow DecisionQuestionScoreInstructionsPartInputImageParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionQuestionScoreInstructionsPartInputImageParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[DecisionQuestionScoreInstructionsPartInputImageParam](
		"detail", "low", "high", "auto", "original",
	)
}

// The properties Data, Type are required.
type DecisionQuestionScoreInstructionsPartInputJSONParam struct {
	Data DecisionQuestionScoreInstructionsPartInputJSONDataUnionParam `json:"data,omitzero" api:"required"`
	// This field can be elided, and will marshal its zero value as "input_json".
	Type constant.InputJSON `json:"type" default:"input_json"`
	paramObj
}

func (r DecisionQuestionScoreInstructionsPartInputJSONParam) MarshalJSON() (data []byte, err error) {
	type shadow DecisionQuestionScoreInstructionsPartInputJSONParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionQuestionScoreInstructionsPartInputJSONParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DecisionQuestionScoreInstructionsPartInputJSONDataUnionParam struct {
	OfLiveDecisionJSONObject *DecisionJSONObject  `json:",omitzero,inline"`
	OfArrayOfJSONObjects     []DecisionJSONObject `json:",omitzero,inline"`
	paramUnion
}

func (u DecisionQuestionScoreInstructionsPartInputJSONDataUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfLiveDecisionJSONObject, u.OfArrayOfJSONObjects)
}
func (u *DecisionQuestionScoreInstructionsPartInputJSONDataUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// The property Label is required.
type DecisionQuestionScoreLevelParam struct {
	Label       string            `json:"label" api:"required"`
	Description param.Opt[string] `json:"description,omitzero"`
	paramObj
}

func (r DecisionQuestionScoreLevelParam) MarshalJSON() (data []byte, err error) {
	type shadow DecisionQuestionScoreLevelParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionQuestionScoreLevelParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type DecisionStoredAssistantMessage struct {
	ID string `json:"id" api:"required"`
	// The text input to the model.
	Content DecisionStoredAssistantMessageContentUnion `json:"content" api:"required"`
	Role    constant.Assistant                         `json:"role" default:"assistant"`
	// Labels an `assistant` message as intermediate commentary (`commentary`) or the
	// final answer (`final_answer`). For models like `gpt-5.3-codex` and beyond, when
	// sending follow-up requests, preserve and resend phase on all assistant messages
	// — dropping it can degrade performance. Not used for user messages.
	//
	// Any of "commentary", "final_answer".
	Phase DecisionStoredAssistantMessagePhase `json:"phase" api:"nullable"`
	// Any of "message".
	Type DecisionStoredAssistantMessageType `json:"type"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Content     respjson.Field
		Role        respjson.Field
		Phase       respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DecisionStoredAssistantMessage) RawJSON() string { return r.JSON.raw }
func (r *DecisionStoredAssistantMessage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// DecisionStoredAssistantMessageContentUnion contains all possible properties and
// values from [string], [[]DecisionStoredAssistantMessageContentPartUnion].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfString OfParts]
type DecisionStoredAssistantMessageContentUnion struct {
	// This field will be present if the value is a [string] instead of an object.
	OfString string `json:",inline"`
	// This field will be present if the value is a
	// [[]DecisionStoredAssistantMessageContentPartUnion] instead of an object.
	OfParts []DecisionStoredAssistantMessageContentPartUnion `json:",inline"`
	JSON    struct {
		OfString respjson.Field
		OfParts  respjson.Field
		raw      string
	} `json:"-"`
}

func (u DecisionStoredAssistantMessageContentUnion) AsString() (v string) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u DecisionStoredAssistantMessageContentUnion) AsParts() (v []DecisionStoredAssistantMessageContentPartUnion) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u DecisionStoredAssistantMessageContentUnion) RawJSON() string { return u.JSON.raw }

func (r *DecisionStoredAssistantMessageContentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// DecisionStoredAssistantMessageContentPartUnion contains all possible properties
// and values from [DecisionStoredAssistantMessageContentPartOutputText],
// [DecisionStoredAssistantMessageContentPartRefusal].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type DecisionStoredAssistantMessageContentPartUnion struct {
	OfDecisionStoredAssistantMessageContentPartOutputText DecisionStoredAssistantMessageContentPartOutputText `json:",inline"`
	OfDecisionStoredAssistantMessageContentPartRefusal    DecisionStoredAssistantMessageContentPartRefusal    `json:",inline"`
	JSON                                                  struct {
		OfDecisionStoredAssistantMessageContentPartOutputText respjson.Field
		OfDecisionStoredAssistantMessageContentPartRefusal    respjson.Field
		raw                                                   string
	} `json:"-"`
}

func (u DecisionStoredAssistantMessageContentPartUnion) AsDecisionStoredAssistantMessageContentPartOutputText() (v DecisionStoredAssistantMessageContentPartOutputText) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u DecisionStoredAssistantMessageContentPartUnion) AsDecisionStoredAssistantMessageContentPartRefusal() (v DecisionStoredAssistantMessageContentPartRefusal) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u DecisionStoredAssistantMessageContentPartUnion) RawJSON() string { return u.JSON.raw }

func (r *DecisionStoredAssistantMessageContentPartUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type DecisionStoredAssistantMessageContentPartOutputText struct {
	// The text content.
	Text string              `json:"text" api:"required"`
	Type constant.OutputText `json:"type" default:"output_text"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Text        respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DecisionStoredAssistantMessageContentPartOutputText) RawJSON() string { return r.JSON.raw }
func (r *DecisionStoredAssistantMessageContentPartOutputText) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type DecisionStoredAssistantMessageContentPartRefusal struct {
	// The refusal text.
	Refusal string           `json:"refusal" api:"required"`
	Type    constant.Refusal `json:"type" default:"refusal"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Refusal     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DecisionStoredAssistantMessageContentPartRefusal) RawJSON() string { return r.JSON.raw }
func (r *DecisionStoredAssistantMessageContentPartRefusal) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Labels an `assistant` message as intermediate commentary (`commentary`) or the
// final answer (`final_answer`). For models like `gpt-5.3-codex` and beyond, when
// sending follow-up requests, preserve and resend phase on all assistant messages
// — dropping it can degrade performance. Not used for user messages.
type DecisionStoredAssistantMessagePhase string

const (
	DecisionStoredAssistantMessagePhaseCommentary  DecisionStoredAssistantMessagePhase = "commentary"
	DecisionStoredAssistantMessagePhaseFinalAnswer DecisionStoredAssistantMessagePhase = "final_answer"
)

type DecisionStoredAssistantMessageType string

const (
	DecisionStoredAssistantMessageTypeMessage DecisionStoredAssistantMessageType = "message"
)

type DecisionStoredDeveloperMessage struct {
	ID string `json:"id" api:"required"`
	// The text input to the model.
	Content DecisionStoredDeveloperMessageContentUnion `json:"content" api:"required"`
	Role    constant.Developer                         `json:"role" default:"developer"`
	// Any of "message".
	Type DecisionStoredDeveloperMessageType `json:"type"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Content     respjson.Field
		Role        respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DecisionStoredDeveloperMessage) RawJSON() string { return r.JSON.raw }
func (r *DecisionStoredDeveloperMessage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// DecisionStoredDeveloperMessageContentUnion contains all possible properties and
// values from [string], [[]DecisionInputText].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfString OfParts]
type DecisionStoredDeveloperMessageContentUnion struct {
	// This field will be present if the value is a [string] instead of an object.
	OfString string `json:",inline"`
	// This field will be present if the value is a [[]DecisionInputText] instead of an
	// object.
	OfParts []DecisionInputText `json:",inline"`
	JSON    struct {
		OfString respjson.Field
		OfParts  respjson.Field
		raw      string
	} `json:"-"`
}

func (u DecisionStoredDeveloperMessageContentUnion) AsString() (v string) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u DecisionStoredDeveloperMessageContentUnion) AsParts() (v []DecisionInputText) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u DecisionStoredDeveloperMessageContentUnion) RawJSON() string { return u.JSON.raw }

func (r *DecisionStoredDeveloperMessageContentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type DecisionStoredDeveloperMessageType string

const (
	DecisionStoredDeveloperMessageTypeMessage DecisionStoredDeveloperMessageType = "message"
)

type DecisionStoredSystemMessage struct {
	ID string `json:"id" api:"required"`
	// The text input to the model.
	Content DecisionStoredSystemMessageContentUnion `json:"content" api:"required"`
	Role    constant.System                         `json:"role" default:"system"`
	// Any of "message".
	Type DecisionStoredSystemMessageType `json:"type"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Content     respjson.Field
		Role        respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DecisionStoredSystemMessage) RawJSON() string { return r.JSON.raw }
func (r *DecisionStoredSystemMessage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// DecisionStoredSystemMessageContentUnion contains all possible properties and
// values from [string], [[]DecisionInputText].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfString OfParts]
type DecisionStoredSystemMessageContentUnion struct {
	// This field will be present if the value is a [string] instead of an object.
	OfString string `json:",inline"`
	// This field will be present if the value is a [[]DecisionInputText] instead of an
	// object.
	OfParts []DecisionInputText `json:",inline"`
	JSON    struct {
		OfString respjson.Field
		OfParts  respjson.Field
		raw      string
	} `json:"-"`
}

func (u DecisionStoredSystemMessageContentUnion) AsString() (v string) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u DecisionStoredSystemMessageContentUnion) AsParts() (v []DecisionInputText) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u DecisionStoredSystemMessageContentUnion) RawJSON() string { return u.JSON.raw }

func (r *DecisionStoredSystemMessageContentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type DecisionStoredSystemMessageType string

const (
	DecisionStoredSystemMessageTypeMessage DecisionStoredSystemMessageType = "message"
)

type DecisionStoredUserMessage struct {
	ID string `json:"id" api:"required"`
	// The text input to the model.
	Content DecisionStoredUserMessageContentUnion `json:"content" api:"required"`
	// The message role. Always `user`.
	Role constant.User `json:"role" default:"user"`
	// Any of "message".
	Type DecisionStoredUserMessageType `json:"type"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Content     respjson.Field
		Role        respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DecisionStoredUserMessage) RawJSON() string { return r.JSON.raw }
func (r *DecisionStoredUserMessage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// DecisionStoredUserMessageContentUnion contains all possible properties and
// values from [string], [[]DecisionStoredUserMessageContentPartUnion].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfString OfParts]
type DecisionStoredUserMessageContentUnion struct {
	// This field will be present if the value is a [string] instead of an object.
	OfString string `json:",inline"`
	// This field will be present if the value is a
	// [[]DecisionStoredUserMessageContentPartUnion] instead of an object.
	OfParts []DecisionStoredUserMessageContentPartUnion `json:",inline"`
	JSON    struct {
		OfString respjson.Field
		OfParts  respjson.Field
		raw      string
	} `json:"-"`
}

func (u DecisionStoredUserMessageContentUnion) AsString() (v string) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u DecisionStoredUserMessageContentUnion) AsParts() (v []DecisionStoredUserMessageContentPartUnion) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u DecisionStoredUserMessageContentUnion) RawJSON() string { return u.JSON.raw }

func (r *DecisionStoredUserMessageContentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// DecisionStoredUserMessageContentPartUnion contains all possible properties and
// values from [DecisionStoredUserMessageContentPartInputText],
// [DecisionStoredUserMessageContentPartInputImage],
// [DecisionStoredUserMessageContentPartInputJSON].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type DecisionStoredUserMessageContentPartUnion struct {
	OfDecisionStoredUserMessageContentPartInputText  DecisionStoredUserMessageContentPartInputText  `json:",inline"`
	OfDecisionStoredUserMessageContentPartInputImage DecisionStoredUserMessageContentPartInputImage `json:",inline"`
	OfDecisionStoredUserMessageContentPartInputJSON  DecisionStoredUserMessageContentPartInputJSON  `json:",inline"`
	JSON                                             struct {
		OfDecisionStoredUserMessageContentPartInputText  respjson.Field
		OfDecisionStoredUserMessageContentPartInputImage respjson.Field
		OfDecisionStoredUserMessageContentPartInputJSON  respjson.Field
		raw                                              string
	} `json:"-"`
}

func (u DecisionStoredUserMessageContentPartUnion) AsDecisionStoredUserMessageContentPartInputText() (v DecisionStoredUserMessageContentPartInputText) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u DecisionStoredUserMessageContentPartUnion) AsDecisionStoredUserMessageContentPartInputImage() (v DecisionStoredUserMessageContentPartInputImage) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u DecisionStoredUserMessageContentPartUnion) AsDecisionStoredUserMessageContentPartInputJSON() (v DecisionStoredUserMessageContentPartInputJSON) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u DecisionStoredUserMessageContentPartUnion) RawJSON() string { return u.JSON.raw }

func (r *DecisionStoredUserMessageContentPartUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type DecisionStoredUserMessageContentPartInputText struct {
	// The text input to the model.
	Text string             `json:"text" api:"required"`
	Type constant.InputText `json:"type" default:"input_text"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Text        respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DecisionStoredUserMessageContentPartInputText) RawJSON() string { return r.JSON.raw }
func (r *DecisionStoredUserMessageContentPartInputText) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type DecisionStoredUserMessageContentPartInputImage struct {
	// The detail level of the image to be sent to the model. One of `high`, `low`,
	// `auto`, or `original`. Defaults to `auto`.
	//
	// Any of "low", "high", "auto", "original".
	Detail string `json:"detail" api:"nullable"`
	// A base64-encoded image data URL or a publicly accessible HTTP(S) image URL.
	ImageURL string              `json:"image_url" api:"required"`
	Type     constant.InputImage `json:"type" default:"input_image"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Detail      respjson.Field
		ImageURL    respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DecisionStoredUserMessageContentPartInputImage) RawJSON() string { return r.JSON.raw }
func (r *DecisionStoredUserMessageContentPartInputImage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type DecisionStoredUserMessageContentPartInputJSON struct {
	Data DecisionStoredUserMessageContentPartInputJSONDataUnion `json:"data" api:"required"`
	Type constant.InputJSON                                     `json:"type" default:"input_json"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Data        respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DecisionStoredUserMessageContentPartInputJSON) RawJSON() string { return r.JSON.raw }
func (r *DecisionStoredUserMessageContentPartInputJSON) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// DecisionStoredUserMessageContentPartInputJSONDataUnion contains all possible
// properties and values from [DecisionJSONObject], [[]DecisionJSONObject].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfArrayOfDecisionJSONObject]
type DecisionStoredUserMessageContentPartInputJSONDataUnion struct {
	OfLiveDecisionJSONObject DecisionJSONObject `json:",inline"`
	// This field will be present if the value is a [[]DecisionJSONObject] instead of
	// an object.
	OfArrayOfDecisionJSONObject []DecisionJSONObject `json:",inline"`
	JSON                        struct {
		OfLiveDecisionJSONObject    respjson.Field
		OfArrayOfDecisionJSONObject respjson.Field
		raw                         string
	} `json:"-"`
}

func (u DecisionStoredUserMessageContentPartInputJSONDataUnion) AsLiveDecisionJSONObject() (v DecisionJSONObject) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u DecisionStoredUserMessageContentPartInputJSONDataUnion) AsArrayOfDecisionJSONObject() (v []DecisionJSONObject) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u DecisionStoredUserMessageContentPartInputJSONDataUnion) RawJSON() string { return u.JSON.raw }

func (r *DecisionStoredUserMessageContentPartInputJSONDataUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type DecisionStoredUserMessageType string

const (
	DecisionStoredUserMessageTypeMessage DecisionStoredUserMessageType = "message"
)

// The properties Content, Role are required.
type DecisionSystemMessageParam struct {
	// The text input to the model.
	Content DecisionSystemMessageContentUnionParam `json:"content,omitzero" api:"required"`
	// The unique ID of this message item.
	ID param.Opt[string] `json:"id,omitzero"`
	// Any of "message".
	Type DecisionSystemMessageType `json:"type,omitzero"`
	// This field can be elided, and will marshal its zero value as "system".
	Role constant.System `json:"role" default:"system"`
	paramObj
}

func (r DecisionSystemMessageParam) MarshalJSON() (data []byte, err error) {
	type shadow DecisionSystemMessageParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionSystemMessageParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DecisionSystemMessageContentUnionParam struct {
	OfString param.Opt[string]        `json:",omitzero,inline"`
	OfParts  []DecisionInputTextParam `json:",omitzero,inline"`
	paramUnion
}

func (u DecisionSystemMessageContentUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfString, u.OfParts)
}
func (u *DecisionSystemMessageContentUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *DecisionSystemMessageContentUnionParam) asAny() any {
	if !param.IsOmitted(u.OfString) {
		return &u.OfString.Value
	} else if !param.IsOmitted(u.OfParts) {
		return &u.OfParts
	}
	return nil
}

type DecisionSystemMessageType string

const (
	DecisionSystemMessageTypeMessage DecisionSystemMessageType = "message"
)

type DecisionTrigger struct {
	// Questions evaluated using the current transcript and the session's shared
	// Decision items.
	Decisions DecisionTriggerDecisions `json:"decisions" api:"required"`
	// Any of "session.input_transcript.delta", "session.output_transcript.delta".
	EventType DecisionTriggerEventType `json:"event_type" api:"required"`
	Type      constant.Decisions       `json:"type" default:"decisions"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Decisions   respjson.Field
		EventType   respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DecisionTrigger) RawJSON() string { return r.JSON.raw }
func (r *DecisionTrigger) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this DecisionTrigger to a DecisionTriggerParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// DecisionTriggerParam.Overrides()
func (r DecisionTrigger) ToParam() DecisionTriggerParam {
	return param.Override[DecisionTriggerParam](json.RawMessage(r.RawJSON()))
}

// Questions evaluated using the current transcript and the session's shared
// Decision items.
type DecisionTriggerDecisions struct {
	// Native Decision API questions evaluated together, in order.
	Questions []DecisionQuestionUnion `json:"questions" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Questions   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DecisionTriggerDecisions) RawJSON() string { return r.JSON.raw }
func (r *DecisionTriggerDecisions) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type DecisionTriggerEventType string

const (
	DecisionTriggerEventTypeSessionInputTranscriptDelta  DecisionTriggerEventType = "session.input_transcript.delta"
	DecisionTriggerEventTypeSessionOutputTranscriptDelta DecisionTriggerEventType = "session.output_transcript.delta"
)

// The properties Decisions, EventType, Type are required.
type DecisionTriggerParam struct {
	// Questions evaluated using the current transcript and the session's shared
	// Decision items.
	Decisions DecisionTriggerDecisionsParam `json:"decisions" api:"required"`
	// Any of "session.input_transcript.delta", "session.output_transcript.delta".
	EventType DecisionTriggerEventType `json:"event_type,omitzero" api:"required"`
	// This field can be elided, and will marshal its zero value as "decisions".
	Type constant.Decisions `json:"type" default:"decisions"`
	paramObj
}

func (r DecisionTriggerParam) MarshalJSON() (data []byte, err error) {
	type shadow DecisionTriggerParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionTriggerParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Questions evaluated using the current transcript and the session's shared
// Decision items.
//
// The property Questions is required.
type DecisionTriggerDecisionsParam struct {
	// Native Decision API questions evaluated together, in order.
	Questions []DecisionQuestionUnionParam `json:"questions,omitzero" api:"required"`
	paramObj
}

func (r DecisionTriggerDecisionsParam) MarshalJSON() (data []byte, err error) {
	type shadow DecisionTriggerDecisionsParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionTriggerDecisionsParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Content, Role are required.
type DecisionUserMessageParam struct {
	// The text input to the model.
	Content DecisionUserMessageContentUnionParam `json:"content,omitzero" api:"required"`
	// The unique ID of this message item.
	ID param.Opt[string] `json:"id,omitzero"`
	// Any of "message".
	Type DecisionUserMessageType `json:"type,omitzero"`
	// The message role. Always `user`.
	//
	// This field can be elided, and will marshal its zero value as "user".
	Role constant.User `json:"role" default:"user"`
	paramObj
}

func (r DecisionUserMessageParam) MarshalJSON() (data []byte, err error) {
	type shadow DecisionUserMessageParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionUserMessageParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DecisionUserMessageContentUnionParam struct {
	OfString param.Opt[string]                          `json:",omitzero,inline"`
	OfParts  []DecisionUserMessageContentPartUnionParam `json:",omitzero,inline"`
	paramUnion
}

func (u DecisionUserMessageContentUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfString, u.OfParts)
}
func (u *DecisionUserMessageContentUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *DecisionUserMessageContentUnionParam) asAny() any {
	if !param.IsOmitted(u.OfString) {
		return &u.OfString.Value
	} else if !param.IsOmitted(u.OfParts) {
		return &u.OfParts
	}
	return nil
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DecisionUserMessageContentPartUnionParam struct {
	OfInputText  *DecisionUserMessageContentPartInputTextParam  `json:",omitzero,inline"`
	OfInputImage *DecisionUserMessageContentPartInputImageParam `json:",omitzero,inline"`
	OfInputJSON  *DecisionUserMessageContentPartInputJSONParam  `json:",omitzero,inline"`
	paramUnion
}

func (u DecisionUserMessageContentPartUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfInputText, u.OfInputImage, u.OfInputJSON)
}
func (u *DecisionUserMessageContentPartUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionUserMessageContentPartUnionParam) GetText() *string {
	if vt := u.OfInputText; vt != nil {
		return &vt.Text
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionUserMessageContentPartUnionParam) GetImageURL() *string {
	if vt := u.OfInputImage; vt != nil {
		return &vt.ImageURL
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionUserMessageContentPartUnionParam) GetDetail() *string {
	if vt := u.OfInputImage; vt != nil {
		return &vt.Detail
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionUserMessageContentPartUnionParam) GetData() *DecisionUserMessageContentPartInputJSONDataUnionParam {
	if vt := u.OfInputJSON; vt != nil {
		return &vt.Data
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionUserMessageContentPartUnionParam) GetType() *string {
	if vt := u.OfInputText; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfInputImage; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfInputJSON; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[DecisionUserMessageContentPartUnionParam](
		"type",
		apijson.Discriminator[DecisionUserMessageContentPartInputTextParam]("input_text"),
		apijson.Discriminator[DecisionUserMessageContentPartInputImageParam]("input_image"),
		apijson.Discriminator[DecisionUserMessageContentPartInputJSONParam]("input_json"),
	)
}

// The properties Text, Type are required.
type DecisionUserMessageContentPartInputTextParam struct {
	// The text input to the model.
	Text string `json:"text" api:"required"`
	// This field can be elided, and will marshal its zero value as "input_text".
	Type constant.InputText `json:"type" default:"input_text"`
	paramObj
}

func (r DecisionUserMessageContentPartInputTextParam) MarshalJSON() (data []byte, err error) {
	type shadow DecisionUserMessageContentPartInputTextParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionUserMessageContentPartInputTextParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties ImageURL, Type are required.
type DecisionUserMessageContentPartInputImageParam struct {
	// A base64-encoded image data URL or a publicly accessible HTTP(S) image URL.
	ImageURL string `json:"image_url" api:"required"`
	// The detail level of the image to be sent to the model. One of `high`, `low`,
	// `auto`, or `original`. Defaults to `auto`.
	//
	// Any of "low", "high", "auto", "original".
	Detail string `json:"detail,omitzero"`
	// This field can be elided, and will marshal its zero value as "input_image".
	Type constant.InputImage `json:"type" default:"input_image"`
	paramObj
}

func (r DecisionUserMessageContentPartInputImageParam) MarshalJSON() (data []byte, err error) {
	type shadow DecisionUserMessageContentPartInputImageParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionUserMessageContentPartInputImageParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[DecisionUserMessageContentPartInputImageParam](
		"detail", "low", "high", "auto", "original",
	)
}

// The properties Data, Type are required.
type DecisionUserMessageContentPartInputJSONParam struct {
	Data DecisionUserMessageContentPartInputJSONDataUnionParam `json:"data,omitzero" api:"required"`
	// This field can be elided, and will marshal its zero value as "input_json".
	Type constant.InputJSON `json:"type" default:"input_json"`
	paramObj
}

func (r DecisionUserMessageContentPartInputJSONParam) MarshalJSON() (data []byte, err error) {
	type shadow DecisionUserMessageContentPartInputJSONParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionUserMessageContentPartInputJSONParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DecisionUserMessageContentPartInputJSONDataUnionParam struct {
	OfLiveDecisionJSONObject *DecisionJSONObject  `json:",omitzero,inline"`
	OfArrayOfJSONObjects     []DecisionJSONObject `json:",omitzero,inline"`
	paramUnion
}

func (u DecisionUserMessageContentPartInputJSONDataUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfLiveDecisionJSONObject, u.OfArrayOfJSONObjects)
}
func (u *DecisionUserMessageContentPartInputJSONDataUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

type DecisionUserMessageType string

const (
	DecisionUserMessageTypeMessage DecisionUserMessageType = "message"
)

type DecisionsDelegationConfig struct {
	Model string `json:"model" api:"required"`
	// Native Decision API questions evaluated together, in order.
	Questions []DecisionQuestionUnion `json:"questions" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Model       respjson.Field
		Questions   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DecisionsDelegationConfig) RawJSON() string { return r.JSON.raw }
func (r *DecisionsDelegationConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this DecisionsDelegationConfig to a
// DecisionsDelegationConfigParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// DecisionsDelegationConfigParam.Overrides()
func (r DecisionsDelegationConfig) ToParam() DecisionsDelegationConfigParam {
	return param.Override[DecisionsDelegationConfigParam](json.RawMessage(r.RawJSON()))
}

// The properties Model, Questions are required.
type DecisionsDelegationConfigParam struct {
	// Native Decision API questions evaluated together, in order.
	Questions []DecisionQuestionUnionParam `json:"questions,omitzero" api:"required"`
	Model     string                       `json:"model" api:"required"`
	paramObj
}

func (r DecisionsDelegationConfigParam) MarshalJSON() (data []byte, err error) {
	type shadow DecisionsDelegationConfigParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionsDelegationConfigParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type DecisionsDelegationUpdateConfigParam struct {
	Model param.Opt[string] `json:"model,omitzero"`
	// Native Decision API questions evaluated together, in order.
	Questions []DecisionQuestionUnionParam `json:"questions,omitzero"`
	paramObj
}

func (r DecisionsDelegationUpdateConfigParam) MarshalJSON() (data []byte, err error) {
	type shadow DecisionsDelegationUpdateConfigParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionsDelegationUpdateConfigParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Returned when the Live model delegates work to your application, Responses, or
// Decisions. Contains delegation metadata and the position on the session timeline
// where the work was delegated.
type DelegationCreatedEvent struct {
	// The delegated work identifier and destination. This object contains metadata,
	// not the task text.
	Delegation DelegationCreatedEventDelegation `json:"delegation" api:"required"`
	// The unique ID of the Live server event.
	EventID string `json:"event_id" api:"required"`
	// The position on the Live session timeline where the delegation was created, in
	// milliseconds from the beginning of the session.
	OffsetMs int64 `json:"offset_ms" api:"required"`
	// The event type, always `session.delegation.created`.
	Type constant.SessionDelegationCreated `json:"type" default:"session.delegation.created"`
	// The event_id of the client command associated with this server event, when
	// supplied.
	ClientEventID string `json:"client_event_id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Delegation    respjson.Field
		EventID       respjson.Field
		OffsetMs      respjson.Field
		Type          respjson.Field
		ClientEventID respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DelegationCreatedEvent) RawJSON() string { return r.JSON.raw }
func (r *DelegationCreatedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The delegated work identifier and destination. This object contains metadata,
// not the task text.
type DelegationCreatedEventDelegation struct {
	// The unique ID of the delegation. Use this as delegation_id when replying to
	// client-owned work or correlating backend events.
	ID string `json:"id" api:"required"`
	// Where the Live model delegated the work: `client`, `responses`, or `decisions`.
	Target string `json:"target" api:"required"`
	// The object type, always `delegation`.
	Type constant.Delegation `json:"type" default:"delegation"`
	// The ID of the Responses API response associated with a Responses delegation.
	// Omitted for client delegations.
	ResponseID string `json:"response_id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Target      respjson.Field
		Type        respjson.Field
		ResponseID  respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DelegationCreatedEventDelegation) RawJSON() string { return r.JSON.raw }
func (r *DelegationCreatedEventDelegation) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Details of an error encountered by the Live session, including the affected
// parameter or client command when available.
type LiveError struct {
	// A machine-readable code identifying the Live error, such as `unknown_parameter`.
	Code string `json:"code" api:"required"`
	// A human-readable explanation of the Live error.
	Message string `json:"message" api:"required"`
	// The category of error, such as `invalid_request_error` for an invalid Live
	// client command.
	Type string `json:"type" api:"required"`
	// The event_id of the client command that caused the error, when supplied.
	ClientEventID string `json:"client_event_id"`
	// The parameter that caused the error, when applicable, such as `session.voice`.
	Param string `json:"param"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Code          respjson.Field
		Message       respjson.Field
		Type          respjson.Field
		ClientEventID respjson.Field
		Param         respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r LiveError) RawJSON() string { return r.JSON.raw }
func (r *LiveError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Reports an error in the Live session, such as an invalid client command. Use
// error.client_event_id, when present, to identify the command that caused the
// error.
type ErrorEvent struct {
	// Details of the Live error and the client command that caused it, when known.
	Error LiveError `json:"error" api:"required"`
	// The unique ID of the Live server event.
	EventID string `json:"event_id" api:"required"`
	// The event type, always `error`.
	Type constant.Error `json:"type" default:"error"`
	// The event_id of the client command associated with this server event, when
	// supplied.
	ClientEventID string `json:"client_event_id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Error         respjson.Field
		EventID       respjson.Field
		Type          respjson.Field
		ClientEventID respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ErrorEvent) RawJSON() string { return r.JSON.raw }
func (r *ErrorEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Overrides for a stored session after connecting to the fork WebSocket. An empty
// object inherits the stored configuration; do not supply a new model.
// audio.format applies only to the new WebSocket connection. client overrides are
// only supported for WebRTC forks.
type ForkSessionConfigParam struct {
	// Whether to store the forked session. Omission inherits the stored session's
	// setting.
	Store param.Opt[bool] `json:"store,omitzero"`
	// Audio format for a WebSocket fork. WebRTC forks negotiate their audio format and
	// must omit this field.
	Audio ForkSessionConfigAudioParam `json:"audio,omitzero"`
	// Frontend data-channel permissions for a WebRTC fork. Omitted permissions inherit
	// the stored values. Not supported for WebSocket forks.
	Client ClientConfigParam `json:"client,omitzero"`
	// Overrides for the stored session’s Responses backend. Only supported when the
	// stored session already uses Responses delegation; the delegation type cannot
	// change.
	Delegation ForkSessionConfigDelegationParam `json:"delegation,omitzero"`
	paramObj
}

func (r ForkSessionConfigParam) MarshalJSON() (data []byte, err error) {
	type shadow ForkSessionConfigParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ForkSessionConfigParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Audio format for a WebSocket fork. WebRTC forks negotiate their audio format and
// must omit this field.
type ForkSessionConfigAudioParam struct {
	// Audio encoding and sample rate for audio sent and received over a Live WebSocket
	// connection. WebRTC and SIP negotiate their media format separately.
	Format AudioFormatUnionParam `json:"format,omitzero"`
	paramObj
}

func (r ForkSessionConfigAudioParam) MarshalJSON() (data []byte, err error) {
	type shadow ForkSessionConfigAudioParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ForkSessionConfigAudioParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Overrides for the stored session’s Responses backend. Only supported when the
// stored session already uses Responses delegation; the delegation type cannot
// change.
//
// The property Type is required.
type ForkSessionConfigDelegationParam struct {
	// Milliseconds to wait before requesting an automatic spoken acknowledgement for a
	// delegation. The model may still speak before this delay. Any commentary or
	// thinking appended to the session cancels all pending acknowledgements, including
	// those for other delegations. Defaults to 400 for Decisions sessions and 0
	// (request immediately) for other sessions. Omission in a session update preserves
	// the current value. Updates apply to subsequent delegations.
	AcknowledgementDelayMs param.Opt[int64] `json:"acknowledgement_delay_ms,omitzero"`
	// Responses backend settings to update. Omitted settings keep their existing
	// values.
	Responses ResponsesDelegationUpdateConfigParam `json:"responses,omitzero"`
	// The delegation owner. Always `responses` for tasks handled by the Responses API.
	//
	// This field can be elided, and will marshal its zero value as "responses".
	Type constant.Responses `json:"type" default:"responses"`
	paramObj
}

func (r ForkSessionConfigDelegationParam) MarshalJSON() (data []byte, err error) {
	type shadow ForkSessionConfigDelegationParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ForkSessionConfigDelegationParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Start a Live session after connecting to a stored session’s fork WebSocket. Send
// an empty `session` object to use the stored configuration.
//
// The properties Session, Type are required.
type ForkSessionStartEventParam struct {
	// Overrides for a stored session after connecting to the fork WebSocket. An empty
	// object inherits the stored configuration; do not supply a new model.
	// audio.format applies only to the new WebSocket connection. client overrides are
	// only supported for WebRTC forks.
	Session ForkSessionConfigParam `json:"session" api:"required"`
	// Optional client identifier for correlating this command with a server event's
	// client_event_id or error.client_event_id.
	EventID param.Opt[string] `json:"event_id,omitzero"`
	// The Live client event type. Always `session.start`.
	//
	// This field can be elided, and will marshal its zero value as "session.start".
	Type constant.SessionStart `json:"type" default:"session.start"`
	paramObj
}

func (r ForkSessionStartEventParam) MarshalJSON() (data []byte, err error) {
	type shadow ForkSessionStartEventParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ForkSessionStartEventParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A function tool available to the Responses backend when the Live model delegates
// a task.
type FunctionTool struct {
	// The name the delegated Responses model uses when calling this function.
	Name string `json:"name" api:"required"`
	// The tool type. Always `function`.
	Type constant.Function `json:"type" default:"function"`
	// What the function does and when the delegated Responses model should call it.
	Description string `json:"description" api:"nullable"`
	// A JSON Schema object describing the arguments accepted by the function.
	Parameters map[string]any `json:"parameters" api:"nullable"`
	// Whether the delegated Responses model must follow the function’s parameter
	// schema exactly.
	Strict bool `json:"strict" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Name        respjson.Field
		Type        respjson.Field
		Description respjson.Field
		Parameters  respjson.Field
		Strict      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r FunctionTool) RawJSON() string { return r.JSON.raw }
func (r *FunctionTool) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this FunctionTool to a FunctionToolParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// FunctionToolParam.Overrides()
func (r FunctionTool) ToParam() FunctionToolParam {
	return param.Override[FunctionToolParam](json.RawMessage(r.RawJSON()))
}

// A function tool available to the Responses backend when the Live model delegates
// a task.
//
// The properties Name, Type are required.
type FunctionToolParam struct {
	// The name the delegated Responses model uses when calling this function.
	Name string `json:"name" api:"required"`
	// What the function does and when the delegated Responses model should call it.
	Description param.Opt[string] `json:"description,omitzero"`
	// Whether the delegated Responses model must follow the function’s parameter
	// schema exactly.
	Strict param.Opt[bool] `json:"strict,omitzero"`
	// A JSON Schema object describing the arguments accepted by the function.
	Parameters map[string]any `json:"parameters,omitzero"`
	// The tool type. Always `function`.
	//
	// This field can be elided, and will marshal its zero value as "function".
	Type constant.Function `json:"type" default:"function"`
	paramObj
}

func (r FunctionToolParam) MarshalJSON() (data []byte, err error) {
	type shadow FunctionToolParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *FunctionToolParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// An informational notice about the Live session, such as the event permissions
// applied to a frontend data channel.
type InfoEvent struct {
	// A machine-readable code for the notice, such as `data_channel_permissions`.
	Code string `json:"code" api:"required"`
	// The unique ID of the Live server event.
	EventID string `json:"event_id" api:"required"`
	// A human-readable explanation of the Live session notice.
	Message string `json:"message" api:"required"`
	// The event type, always `info`.
	Type constant.Info `json:"type" default:"info"`
	// The event_id of the client command associated with this server event, when
	// supplied.
	ClientEventID string `json:"client_event_id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Code          respjson.Field
		EventID       respjson.Field
		Message       respjson.Field
		Type          respjson.Field
		ClientEventID respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r InfoEvent) RawJSON() string { return r.JSON.raw }
func (r *InfoEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// InitialItemUnion contains all possible properties and values from
// [InitialItemDeveloper], [InitialItemUser], [InitialItemAssistant].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type InitialItemUnion struct {
	OfInitialItemDeveloper InitialItemDeveloper `json:",inline"`
	OfInitialItemUser      InitialItemUser      `json:",inline"`
	OfInitialItemAssistant InitialItemAssistant `json:",inline"`
	JSON                   struct {
		OfInitialItemDeveloper respjson.Field
		OfInitialItemUser      respjson.Field
		OfInitialItemAssistant respjson.Field
		raw                    string
	} `json:"-"`
}

func (u InitialItemUnion) AsInitialItemDeveloper() (v InitialItemDeveloper) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u InitialItemUnion) AsInitialItemUser() (v InitialItemUser) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u InitialItemUnion) AsInitialItemAssistant() (v InitialItemAssistant) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u InitialItemUnion) RawJSON() string { return u.JSON.raw }

func (r *InitialItemUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this InitialItemUnion to a InitialItemUnionParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// InitialItemUnionParam.Overrides()
func (r InitialItemUnion) ToParam() InitialItemUnionParam {
	return param.Override[InitialItemUnionParam](json.RawMessage(r.RawJSON()))
}

// A developer message included in the initial text history of a Live session.
type InitialItemDeveloper struct {
	// The message content. Supply exactly one text part for the initial Live
	// conversation history.
	Content []InitialItemDeveloperContent `json:"content" api:"required"`
	// An optional identifier for the supplied history message. Live uses the message’s
	// role and text to initialize the conversation.
	ID string `json:"id" api:"nullable"`
	// The author of this history message. Always `developer`.
	Role constant.Developer `json:"role" default:"developer"`
	// The supplied message’s status. Live uses its text as history and does not resume
	// an incomplete message.
	//
	// Any of "incomplete", "completed".
	Status string `json:"status" api:"nullable"`
	// The history item type. Always `message`.
	//
	// Any of "message".
	Type string `json:"type"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Content     respjson.Field
		ID          respjson.Field
		Role        respjson.Field
		Status      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r InitialItemDeveloper) RawJSON() string { return r.JSON.raw }
func (r *InitialItemDeveloper) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Text supplied in a developer or user message when starting a Live session.
type InitialItemDeveloperContent struct {
	// The message text to include in the Live session’s initial conversation history.
	Text string `json:"text" api:"required"`
	// The text content type. Always `input_text`.
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
func (r InitialItemDeveloperContent) RawJSON() string { return r.JSON.raw }
func (r *InitialItemDeveloperContent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A user message included in the initial text history of a Live session.
type InitialItemUser struct {
	// The message content. Supply exactly one text part for the initial Live
	// conversation history.
	Content []InitialItemUserContent `json:"content" api:"required"`
	// An optional identifier for the supplied history message. Live uses the message’s
	// role and text to initialize the conversation.
	ID string `json:"id" api:"nullable"`
	// The author of this history message. Always `user`.
	Role constant.User `json:"role" default:"user"`
	// The supplied message’s status. Live uses its text as history and does not resume
	// an incomplete message.
	//
	// Any of "incomplete", "completed".
	Status string `json:"status" api:"nullable"`
	// The history item type. Always `message`.
	//
	// Any of "message".
	Type string `json:"type"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Content     respjson.Field
		ID          respjson.Field
		Role        respjson.Field
		Status      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r InitialItemUser) RawJSON() string { return r.JSON.raw }
func (r *InitialItemUser) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Text supplied in a developer or user message when starting a Live session.
type InitialItemUserContent struct {
	// The message text to include in the Live session’s initial conversation history.
	Text string `json:"text" api:"required"`
	// The text content type. Always `input_text`.
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
func (r InitialItemUserContent) RawJSON() string { return r.JSON.raw }
func (r *InitialItemUserContent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// An assistant message included in the initial text history of a Live session.
type InitialItemAssistant struct {
	// The message content. Supply exactly one text part for the initial Live
	// conversation history.
	Content []InitialItemAssistantContentUnion `json:"content" api:"required"`
	// An optional identifier for the supplied history message. Live uses the message’s
	// role and text to initialize the conversation.
	ID string `json:"id" api:"nullable"`
	// The author of this history message. Always `assistant`.
	Role constant.Assistant `json:"role" default:"assistant"`
	// The supplied message’s status. Live uses its text as history and does not resume
	// an incomplete message.
	//
	// Any of "incomplete", "completed".
	Status string `json:"status" api:"nullable"`
	// The history item type. Always `message`.
	//
	// Any of "message".
	Type string `json:"type"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Content     respjson.Field
		ID          respjson.Field
		Role        respjson.Field
		Status      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r InitialItemAssistant) RawJSON() string { return r.JSON.raw }
func (r *InitialItemAssistant) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// InitialItemAssistantContentUnion contains all possible properties and values
// from [InitialItemAssistantContentText], [InitialItemAssistantContentOutputText].
//
// Use the [InitialItemAssistantContentUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type InitialItemAssistantContentUnion struct {
	OfInitialItemAssistantContentText       InitialItemAssistantContentText       `json:",inline"`
	OfInitialItemAssistantContentOutputText InitialItemAssistantContentOutputText `json:",inline"`
	Type                                    string                                `json:"type"`
	JSON                                    struct {
		OfInitialItemAssistantContentText       respjson.Field
		OfInitialItemAssistantContentOutputText respjson.Field
		Type                                    respjson.Field
		raw                                     string
	} `json:"-"`
}

// anyInitialItemAssistantContent is implemented by each variant of
// [InitialItemAssistantContentUnion] to add type safety for the return type of
// [InitialItemAssistantContentUnion.AsAny]
type anyInitialItemAssistantContent interface {
	implInitialItemAssistantContentUnion()
}

func (InitialItemAssistantContentText) implInitialItemAssistantContentUnion()       {}
func (InitialItemAssistantContentOutputText) implInitialItemAssistantContentUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := InitialItemAssistantContentUnion.AsAny().(type) {
//	case live.InitialItemAssistantContentText:
//	case live.InitialItemAssistantContentOutputText:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u InitialItemAssistantContentUnion) AsAny() anyInitialItemAssistantContent {
	switch u.Type {
	case "text":
		return u.AsInitialItemAssistantContentText()
	case "output_text":
		return u.AsInitialItemAssistantContentOutputText()
	}
	return nil
}

func (u InitialItemAssistantContentUnion) AsInitialItemAssistantContentText() (v InitialItemAssistantContentText) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u InitialItemAssistantContentUnion) AsInitialItemAssistantContentOutputText() (v InitialItemAssistantContentOutputText) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u InitialItemAssistantContentUnion) RawJSON() string { return u.JSON.raw }

func (r *InitialItemAssistantContentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Assistant text supplied as conversation history when starting a Live session.
type InitialItemAssistantContentText struct {
	// The message text to include in the Live session’s initial conversation history.
	Text string `json:"text" api:"required"`
	// The text content type. Always `text`.
	//
	// Any of "text".
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
func (r InitialItemAssistantContentText) RawJSON() string { return r.JSON.raw }
func (r *InitialItemAssistantContentText) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Assistant output text supplied as conversation history when starting a Live
// session.
type InitialItemAssistantContentOutputText struct {
	// The message text to include in the Live session’s initial conversation history.
	Text string `json:"text" api:"required"`
	// The text content type. Always `output_text`.
	Type constant.OutputText `json:"type" default:"output_text"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Text        respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r InitialItemAssistantContentOutputText) RawJSON() string { return r.JSON.raw }
func (r *InitialItemAssistantContentOutputText) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func InitialItemParamOfDeveloper(content []InitialItemDeveloperContentParam) InitialItemUnionParam {
	var developer InitialItemDeveloperParam
	developer.Content = content
	return InitialItemUnionParam{OfDeveloper: &developer}
}

func InitialItemParamOfUser(content []InitialItemUserContentParam) InitialItemUnionParam {
	var user InitialItemUserParam
	user.Content = content
	return InitialItemUnionParam{OfUser: &user}
}

func InitialItemParamOfAssistant(content []InitialItemAssistantContentUnionParam) InitialItemUnionParam {
	var assistant InitialItemAssistantParam
	assistant.Content = content
	return InitialItemUnionParam{OfAssistant: &assistant}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type InitialItemUnionParam struct {
	OfDeveloper *InitialItemDeveloperParam `json:",omitzero,inline"`
	OfUser      *InitialItemUserParam      `json:",omitzero,inline"`
	OfAssistant *InitialItemAssistantParam `json:",omitzero,inline"`
	paramUnion
}

func (u InitialItemUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfDeveloper, u.OfUser, u.OfAssistant)
}
func (u *InitialItemUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u InitialItemUnionParam) GetRole() *string {
	if vt := u.OfDeveloper; vt != nil {
		return (*string)(&vt.Role)
	} else if vt := u.OfUser; vt != nil {
		return (*string)(&vt.Role)
	} else if vt := u.OfAssistant; vt != nil {
		return (*string)(&vt.Role)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u InitialItemUnionParam) GetID() *string {
	if vt := u.OfDeveloper; vt != nil && vt.ID.Valid() {
		return &vt.ID.Value
	} else if vt := u.OfUser; vt != nil && vt.ID.Valid() {
		return &vt.ID.Value
	} else if vt := u.OfAssistant; vt != nil && vt.ID.Valid() {
		return &vt.ID.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u InitialItemUnionParam) GetStatus() *string {
	if vt := u.OfDeveloper; vt != nil {
		return (*string)(&vt.Status)
	} else if vt := u.OfUser; vt != nil {
		return (*string)(&vt.Status)
	} else if vt := u.OfAssistant; vt != nil {
		return (*string)(&vt.Status)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u InitialItemUnionParam) GetType() *string {
	if vt := u.OfDeveloper; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfUser; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfAssistant; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u InitialItemUnionParam) GetContent() (res initialItemUnionParamContent) {
	if vt := u.OfDeveloper; vt != nil {
		res.any = &vt.Content
	} else if vt := u.OfUser; vt != nil {
		res.any = &vt.Content
	} else if vt := u.OfAssistant; vt != nil {
		res.any = &vt.Content
	}
	return
}

// Can have the runtime types [_[]InitialItemDeveloperContentParam],
// [_[]InitialItemUserContentParam], [\*[]InitialItemAssistantContentUnionParam]
type initialItemUnionParamContent struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *[]live.InitialItemDeveloperContentParam:
//	case *[]live.InitialItemUserContentParam:
//	case *[]live.InitialItemAssistantContentUnionParam:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u initialItemUnionParamContent) AsAny() any { return u.any }

func init() {
	apijson.RegisterUnion[InitialItemUnionParam](
		"role",
		apijson.Discriminator[InitialItemDeveloperParam]("developer"),
		apijson.Discriminator[InitialItemUserParam]("user"),
		apijson.Discriminator[InitialItemAssistantParam]("assistant"),
	)
}

// A developer message included in the initial text history of a Live session.
//
// The properties Content, Role are required.
type InitialItemDeveloperParam struct {
	// The message content. Supply exactly one text part for the initial Live
	// conversation history.
	Content []InitialItemDeveloperContentParam `json:"content,omitzero" api:"required"`
	// An optional identifier for the supplied history message. Live uses the message’s
	// role and text to initialize the conversation.
	ID param.Opt[string] `json:"id,omitzero"`
	// The supplied message’s status. Live uses its text as history and does not resume
	// an incomplete message.
	//
	// Any of "incomplete", "completed".
	Status string `json:"status,omitzero"`
	// The history item type. Always `message`.
	//
	// Any of "message".
	Type string `json:"type,omitzero"`
	// The author of this history message. Always `developer`.
	//
	// This field can be elided, and will marshal its zero value as "developer".
	Role constant.Developer `json:"role" default:"developer"`
	paramObj
}

func (r InitialItemDeveloperParam) MarshalJSON() (data []byte, err error) {
	type shadow InitialItemDeveloperParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *InitialItemDeveloperParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[InitialItemDeveloperParam](
		"status", "incomplete", "completed",
	)
	apijson.RegisterFieldValidator[InitialItemDeveloperParam](
		"type", "message",
	)
}

// Text supplied in a developer or user message when starting a Live session.
//
// The property Text is required.
type InitialItemDeveloperContentParam struct {
	// The message text to include in the Live session’s initial conversation history.
	Text string `json:"text" api:"required"`
	// The text content type. Always `input_text`.
	//
	// Any of "input_text".
	Type string `json:"type,omitzero"`
	paramObj
}

func (r InitialItemDeveloperContentParam) MarshalJSON() (data []byte, err error) {
	type shadow InitialItemDeveloperContentParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *InitialItemDeveloperContentParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[InitialItemDeveloperContentParam](
		"type", "input_text",
	)
}

// A user message included in the initial text history of a Live session.
//
// The properties Content, Role are required.
type InitialItemUserParam struct {
	// The message content. Supply exactly one text part for the initial Live
	// conversation history.
	Content []InitialItemUserContentParam `json:"content,omitzero" api:"required"`
	// An optional identifier for the supplied history message. Live uses the message’s
	// role and text to initialize the conversation.
	ID param.Opt[string] `json:"id,omitzero"`
	// The supplied message’s status. Live uses its text as history and does not resume
	// an incomplete message.
	//
	// Any of "incomplete", "completed".
	Status string `json:"status,omitzero"`
	// The history item type. Always `message`.
	//
	// Any of "message".
	Type string `json:"type,omitzero"`
	// The author of this history message. Always `user`.
	//
	// This field can be elided, and will marshal its zero value as "user".
	Role constant.User `json:"role" default:"user"`
	paramObj
}

func (r InitialItemUserParam) MarshalJSON() (data []byte, err error) {
	type shadow InitialItemUserParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *InitialItemUserParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[InitialItemUserParam](
		"status", "incomplete", "completed",
	)
	apijson.RegisterFieldValidator[InitialItemUserParam](
		"type", "message",
	)
}

// Text supplied in a developer or user message when starting a Live session.
//
// The property Text is required.
type InitialItemUserContentParam struct {
	// The message text to include in the Live session’s initial conversation history.
	Text string `json:"text" api:"required"`
	// The text content type. Always `input_text`.
	//
	// Any of "input_text".
	Type string `json:"type,omitzero"`
	paramObj
}

func (r InitialItemUserContentParam) MarshalJSON() (data []byte, err error) {
	type shadow InitialItemUserContentParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *InitialItemUserContentParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[InitialItemUserContentParam](
		"type", "input_text",
	)
}

// An assistant message included in the initial text history of a Live session.
//
// The properties Content, Role are required.
type InitialItemAssistantParam struct {
	// The message content. Supply exactly one text part for the initial Live
	// conversation history.
	Content []InitialItemAssistantContentUnionParam `json:"content,omitzero" api:"required"`
	// An optional identifier for the supplied history message. Live uses the message’s
	// role and text to initialize the conversation.
	ID param.Opt[string] `json:"id,omitzero"`
	// The supplied message’s status. Live uses its text as history and does not resume
	// an incomplete message.
	//
	// Any of "incomplete", "completed".
	Status string `json:"status,omitzero"`
	// The history item type. Always `message`.
	//
	// Any of "message".
	Type string `json:"type,omitzero"`
	// The author of this history message. Always `assistant`.
	//
	// This field can be elided, and will marshal its zero value as "assistant".
	Role constant.Assistant `json:"role" default:"assistant"`
	paramObj
}

func (r InitialItemAssistantParam) MarshalJSON() (data []byte, err error) {
	type shadow InitialItemAssistantParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *InitialItemAssistantParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[InitialItemAssistantParam](
		"status", "incomplete", "completed",
	)
	apijson.RegisterFieldValidator[InitialItemAssistantParam](
		"type", "message",
	)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type InitialItemAssistantContentUnionParam struct {
	OfText       *InitialItemAssistantContentTextParam       `json:",omitzero,inline"`
	OfOutputText *InitialItemAssistantContentOutputTextParam `json:",omitzero,inline"`
	paramUnion   `default:"text"`
}

func (u InitialItemAssistantContentUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfText, u.OfOutputText)
}
func (u *InitialItemAssistantContentUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u InitialItemAssistantContentUnionParam) GetText() *string {
	if vt := u.OfText; vt != nil {
		return (*string)(&vt.Text)
	} else if vt := u.OfOutputText; vt != nil {
		return (*string)(&vt.Text)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u InitialItemAssistantContentUnionParam) GetType() *string {
	if vt := u.OfText; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfOutputText; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[InitialItemAssistantContentUnionParam](
		"type",
		apijson.Discriminator[InitialItemAssistantContentTextParam]("text"),
		apijson.Discriminator[InitialItemAssistantContentOutputTextParam]("output_text"),
		apijson.Discriminator[InitialItemAssistantContentTextParam](nil),
	)
}

// Assistant text supplied as conversation history when starting a Live session.
//
// The property Text is required.
type InitialItemAssistantContentTextParam struct {
	// The message text to include in the Live session’s initial conversation history.
	Text string `json:"text" api:"required"`
	// The text content type. Always `text`.
	//
	// Any of "text".
	Type string `json:"type,omitzero"`
	paramObj
}

func (r InitialItemAssistantContentTextParam) MarshalJSON() (data []byte, err error) {
	type shadow InitialItemAssistantContentTextParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *InitialItemAssistantContentTextParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[InitialItemAssistantContentTextParam](
		"type", "text",
	)
}

// Assistant output text supplied as conversation history when starting a Live
// session.
//
// The properties Text, Type are required.
type InitialItemAssistantContentOutputTextParam struct {
	// The message text to include in the Live session’s initial conversation history.
	Text string `json:"text" api:"required"`
	// The text content type. Always `output_text`.
	//
	// This field can be elided, and will marshal its zero value as "output_text".
	Type constant.OutputText `json:"type" default:"output_text"`
	paramObj
}

func (r InitialItemAssistantContentOutputTextParam) MarshalJSON() (data []byte, err error) {
	type shadow InitialItemAssistantContentOutputTextParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *InitialItemAssistantContentOutputTextParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Send audio to a Live session over its primary WebSocket. WebRTC and SIP sessions
// send audio over their media transport.
//
// The properties Audio, Type are required.
type InputAudioAppendEventParam struct {
	// Base64-encoded raw audio in the startup-selected format, without a WAV or other
	// container header. Primary WebSocket only; media transports use their audio
	// track. Audio appends have no acknowledgment. Reflected sideband server events
	// reuse this event type and audio key, with no timestamps or event_id; their audio
	// is always mono PCM16LE at 24 kHz.
	Audio string `json:"audio" api:"required"`
	// Optional client identifier for correlating this command with a server event's
	// client_event_id or error.client_event_id.
	EventID param.Opt[string] `json:"event_id,omitzero"`
	// The Live client event type. Always `session.input_audio.append`.
	//
	// This field can be elided, and will marshal its zero value as
	// "session.input_audio.append".
	Type constant.SessionInputAudioAppend `json:"type" default:"session.input_audio.append"`
	paramObj
}

func (r InputAudioAppendEventParam) MarshalJSON() (data []byte, err error) {
	type shadow InputAudioAppendEventParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *InputAudioAppendEventParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Mute audio input to the Live model without closing the session. The server
// acknowledges with `session.input_audio.muted`.
//
// The property Type is required.
type InputAudioMuteEventParam struct {
	// Optional client identifier for correlating this command with a server event's
	// client_event_id or error.client_event_id.
	EventID param.Opt[string] `json:"event_id,omitzero"`
	// The Live client event type. Always `session.input_audio.mute`.
	//
	// This field can be elided, and will marshal its zero value as
	// "session.input_audio.mute".
	Type constant.SessionInputAudioMute `json:"type" default:"session.input_audio.mute"`
	paramObj
}

func (r InputAudioMuteEventParam) MarshalJSON() (data []byte, err error) {
	type shadow InputAudioMuteEventParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *InputAudioMuteEventParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Acknowledges a session.input_audio.mute command. An independent input mute sends
// silence to the model and preserves sideband audio reflection. Fixture playback
// does not emit this event.
type InputAudioMutedEvent struct {
	// The unique ID of the Live server event.
	EventID string `json:"event_id" api:"required"`
	// The event type, always `session.input_audio.muted`.
	Type constant.SessionInputAudioMuted `json:"type" default:"session.input_audio.muted"`
	// The event_id of the client command associated with this server event, when
	// supplied.
	ClientEventID string `json:"client_event_id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventID       respjson.Field
		Type          respjson.Field
		ClientEventID respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r InputAudioMutedEvent) RawJSON() string { return r.JSON.raw }
func (r *InputAudioMutedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Resume audio input to a Live model after muting it. The server acknowledges with
// `session.input_audio.unmuted`. An active audio fixture still pauses inference
// until playback finishes.
//
// The property Type is required.
type InputAudioUnmuteEventParam struct {
	// Optional client identifier for correlating this command with a server event's
	// client_event_id or error.client_event_id.
	EventID param.Opt[string] `json:"event_id,omitzero"`
	// The Live client event type. Always `session.input_audio.unmute`.
	//
	// This field can be elided, and will marshal its zero value as
	// "session.input_audio.unmute".
	Type constant.SessionInputAudioUnmute `json:"type" default:"session.input_audio.unmute"`
	paramObj
}

func (r InputAudioUnmuteEventParam) MarshalJSON() (data []byte, err error) {
	type shadow InputAudioUnmuteEventParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *InputAudioUnmuteEventParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Acknowledges a session.input_audio.unmute command, removing the independent
// input mute. An active audio fixture still pauses inference until playback
// finishes. Fixture playback does not emit this event.
type InputAudioUnmutedEvent struct {
	// The unique ID of the Live server event.
	EventID string `json:"event_id" api:"required"`
	// The event type, always `session.input_audio.unmuted`.
	Type constant.SessionInputAudioUnmuted `json:"type" default:"session.input_audio.unmuted"`
	// The event_id of the client command associated with this server event, when
	// supplied.
	ClientEventID string `json:"client_event_id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventID       respjson.Field
		Type          respjson.Field
		ClientEventID respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r InputAudioUnmutedEvent) RawJSON() string { return r.JSON.raw }
func (r *InputAudioUnmutedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A transcript fragment for user input audio in the Live session. Accumulate
// fragments in delivery order; these events do not define complete turns or
// include a transcript-done event.
type InputTranscriptDeltaEvent struct {
	// The transcript text fragment for the audio in this time range. Append fragments
	// in delivery order to build the transcript.
	Delta string `json:"delta" api:"required"`
	// The end of this event on the Live session timeline, in milliseconds from the
	// beginning of the session. For appended context, this can equal start_ms.
	EndMs int64 `json:"end_ms" api:"required"`
	// The unique ID of the Live server event.
	EventID string `json:"event_id" api:"required"`
	// The start of this event on the Live session timeline, in milliseconds from the
	// beginning of the session.
	StartMs int64 `json:"start_ms" api:"required"`
	// The event type, always `session.input_transcript.delta`.
	Type constant.SessionInputTranscriptDelta `json:"type" default:"session.input_transcript.delta"`
	// The event_id of the client command associated with this server event, when
	// supplied.
	ClientEventID string `json:"client_event_id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Delta         respjson.Field
		EndMs         respjson.Field
		EventID       respjson.Field
		StartMs       respjson.Field
		Type          respjson.Field
		ClientEventID respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r InputTranscriptDeltaEvent) RawJSON() string { return r.JSON.raw }
func (r *InputTranscriptDeltaEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Append developer instructions to the Live conversation while it is running,
// optionally associating them with an existing client delegation. Rejected while
// an audio fixture is active. These instructions remain in conversation history
// through compaction and are included in future backend delegations. For external
// audio, describe content already delivered so the assistant does not repeat it.
//
// The properties Content, DelegationID, Type are required.
type InstructionsAppendEventParam struct {
	// Required, nullable. Set null for general session context, or use the ID from
	// session.delegation.created for an existing client delegation. Non-null IDs are
	// not accepted with Responses delegation.
	DelegationID param.Opt[string] `json:"delegation_id,omitzero" api:"required"`
	// Instruction text to append, limited to 500 tokens. This is a plain string, not
	// an array of content parts.
	Content string `json:"content" api:"required"`
	// Optional client identifier for correlating this command with a server event's
	// client_event_id or error.client_event_id.
	EventID param.Opt[string] `json:"event_id,omitzero"`
	// The Live client event type. Always `session.instructions.append`.
	//
	// This field can be elided, and will marshal its zero value as
	// "session.instructions.append".
	Type constant.SessionInstructionsAppend `json:"type" default:"session.instructions.append"`
	paramObj
}

func (r InstructionsAppendEventParam) MarshalJSON() (data []byte, err error) {
	type shadow InstructionsAppendEventParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *InstructionsAppendEventParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Returned when a session.instructions.append command is accepted into the Live
// session timeline. Acknowledges the appended instructions without guaranteeing
// that the model has acted on them.
type InstructionsAppendedEvent struct {
	// The end of this event on the Live session timeline, in milliseconds from the
	// beginning of the session. For appended context, this can equal start_ms.
	EndMs int64 `json:"end_ms" api:"required"`
	// The unique ID of the Live server event.
	EventID string `json:"event_id" api:"required"`
	// The start of this event on the Live session timeline, in milliseconds from the
	// beginning of the session.
	StartMs int64 `json:"start_ms" api:"required"`
	// The event type, always `session.instructions.appended`.
	Type constant.SessionInstructionsAppended `json:"type" default:"session.instructions.appended"`
	// The event_id of the client command associated with this server event, when
	// supplied.
	ClientEventID string `json:"client_event_id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EndMs         respjson.Field
		EventID       respjson.Field
		StartMs       respjson.Field
		Type          respjson.Field
		ClientEventID respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r InstructionsAppendedEvent) RawJSON() string { return r.JSON.raw }
func (r *InstructionsAppendedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Startup configuration for a Live media session. Follow the
// [Live prompting guide](https://developers.openai.com/api/docs/guides/live-prompting)
// when writing frontend instructions and the backend prompt under
// delegation.responses.instructions.
//
// The property Model is required.
type MediaSessionConfigParam struct {
	// The Live model. Required in the session configuration for every transport; do
	// not pass it as a URL query parameter.
	Model MediaSessionConfigModel `json:"model,omitzero" api:"required"`
	// Frontend instructions for voice, conversation, interruptions, and when to
	// delegate. Start with the
	// [Live prompting guide](https://developers.openai.com/api/docs/guides/live-prompting);
	// put business rules and tool workflows in a separate
	// [backend prompt](https://developers.openai.com/api/docs/guides/live-delegation#start-with-your-existing-backend-prompt).
	// Limited to 16,384 client-supplied tokens. Omitted or blank instructions use
	// server defaults. Immutable after startup.
	Instructions param.Opt[string] `json:"instructions,omitzero"`
	// Whether to store the session for later forking and recording download. Defaults
	// to false for new sessions.
	Store param.Opt[bool] `json:"store,omitzero"`
	// Who handles tasks delegated by the Live model. Omitted or null selects your
	// application; use `responses` or `decisions` to delegate to the corresponding
	// API.
	Delegation MediaSessionConfigDelegationUnionParam `json:"delegation,omitzero"`
	// Startup audio configuration. WebRTC and SIP negotiate their audio format on the
	// media transport.
	Audio MediaSessionConfigAudioParam `json:"audio,omitzero"`
	// Startup-only capabilities for an untrusted frontend attached to a unified WebRTC
	// session. Trusted sideband connections are unaffected.
	Client ClientConfigParam `json:"client,omitzero"`
	// Ordered text-only history supplied before startup. Supports developer, user, and
	// assistant messages with one text part each; at most 128 messages and a
	// configured rendered-token limit (8,192 by default).
	Input []InitialItemUnionParam `json:"input,omitzero"`
	// At most one Decisions trigger per transcript event type. Triggers run
	// independently of delegation. On update, omit triggers to keep the current
	// configuration, send a list to replace it, or send an empty list to disable
	// triggers.
	Triggers []DecisionTriggerParam `json:"triggers,omitzero"`
	paramObj
}

func (r MediaSessionConfigParam) MarshalJSON() (data []byte, err error) {
	type shadow MediaSessionConfigParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *MediaSessionConfigParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The Live model. Required in the session configuration for every transport; do
// not pass it as a URL query parameter.
type MediaSessionConfigModel = string

const (
	MediaSessionConfigModelGPTLive1 MediaSessionConfigModel = "gpt-live-1"
)

// Startup audio configuration. WebRTC and SIP negotiate their audio format on the
// media transport.
type MediaSessionConfigAudioParam struct {
	// Settings for speech generated by the Live model. Choose the voice before
	// starting the session.
	Output MediaSessionConfigAudioOutputParam `json:"output,omitzero"`
	paramObj
}

func (r MediaSessionConfigAudioParam) MarshalJSON() (data []byte, err error) {
	type shadow MediaSessionConfigAudioParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *MediaSessionConfigAudioParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Settings for speech generated by the Live model. Choose the voice before
// starting the session.
type MediaSessionConfigAudioOutputParam struct {
	// The voice used for Live speech, as a built-in voice name or a custom voice
	// object containing its ID. Defaults to `marin` and cannot change after startup.
	Voice MediaSessionConfigAudioOutputVoiceUnionParam `json:"voice,omitzero"`
	paramObj
}

func (r MediaSessionConfigAudioOutputParam) MarshalJSON() (data []byte, err error) {
	type shadow MediaSessionConfigAudioOutputParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *MediaSessionConfigAudioOutputParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type MediaSessionConfigAudioOutputVoiceUnionParam struct {
	OfString param.Opt[string] `json:",omitzero,inline"`
	// Check if union is this variant with !param.IsOmitted(union.OfBuiltIn)
	OfBuiltIn     param.Opt[BuiltInVoice] `json:",omitzero,inline"`
	OfCustomVoice *CustomVoiceParam       `json:",omitzero,inline"`
	paramUnion
}

func (u MediaSessionConfigAudioOutputVoiceUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfString, u.OfBuiltIn, u.OfCustomVoice)
}
func (u *MediaSessionConfigAudioOutputVoiceUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type MediaSessionConfigDelegationUnionParam struct {
	OfClient    *ClientDelegationParam                      `json:",omitzero,inline"`
	OfResponses *MediaSessionConfigDelegationResponsesParam `json:",omitzero,inline"`
	OfDecisions *MediaSessionConfigDelegationDecisionsParam `json:",omitzero,inline"`
	paramUnion
}

func (u MediaSessionConfigDelegationUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfClient, u.OfResponses, u.OfDecisions)
}
func (u *MediaSessionConfigDelegationUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u MediaSessionConfigDelegationUnionParam) GetResponses() *ResponsesDelegationConfigParam {
	if vt := u.OfResponses; vt != nil {
		return &vt.Responses
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u MediaSessionConfigDelegationUnionParam) GetDecisions() *DecisionsDelegationConfigParam {
	if vt := u.OfDecisions; vt != nil {
		return &vt.Decisions
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u MediaSessionConfigDelegationUnionParam) GetType() *string {
	if vt := u.OfClient; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfResponses; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfDecisions; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u MediaSessionConfigDelegationUnionParam) GetAcknowledgementDelayMs() *int64 {
	if vt := u.OfClient; vt != nil && vt.AcknowledgementDelayMs.Valid() {
		return &vt.AcknowledgementDelayMs.Value
	} else if vt := u.OfResponses; vt != nil && vt.AcknowledgementDelayMs.Valid() {
		return &vt.AcknowledgementDelayMs.Value
	} else if vt := u.OfDecisions; vt != nil && vt.AcknowledgementDelayMs.Valid() {
		return &vt.AcknowledgementDelayMs.Value
	}
	return nil
}

func init() {
	apijson.RegisterUnion[MediaSessionConfigDelegationUnionParam](
		"type",
		apijson.Discriminator[ClientDelegationParam]("client"),
		apijson.Discriminator[MediaSessionConfigDelegationResponsesParam]("responses"),
		apijson.Discriminator[MediaSessionConfigDelegationDecisionsParam]("decisions"),
	)
}

// Delegate tasks to a Responses model managed by the Live session.
//
// The properties Responses, Type are required.
type MediaSessionConfigDelegationResponsesParam struct {
	// Backend model, prompt, and tools used when the Live session delegates a task to
	// Responses.
	Responses ResponsesDelegationConfigParam `json:"responses,omitzero" api:"required"`
	// Milliseconds to wait before requesting an automatic spoken acknowledgement for a
	// delegation. The model may still speak before this delay. Any commentary or
	// thinking appended to the session cancels all pending acknowledgements, including
	// those for other delegations. Defaults to 400 for Decisions sessions and 0
	// (request immediately) for other sessions. Omission in a session update preserves
	// the current value. Updates apply to subsequent delegations.
	AcknowledgementDelayMs param.Opt[int64] `json:"acknowledgement_delay_ms,omitzero"`
	// The delegation owner. Always `responses` for tasks handled by the Responses API.
	//
	// This field can be elided, and will marshal its zero value as "responses".
	Type constant.Responses `json:"type" default:"responses"`
	paramObj
}

func (r MediaSessionConfigDelegationResponsesParam) MarshalJSON() (data []byte, err error) {
	type shadow MediaSessionConfigDelegationResponsesParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *MediaSessionConfigDelegationResponsesParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Decisions, Type are required.
type MediaSessionConfigDelegationDecisionsParam struct {
	Decisions DecisionsDelegationConfigParam `json:"decisions" api:"required"`
	// Milliseconds to wait before requesting an automatic spoken acknowledgement for a
	// delegation. The model may still speak before this delay. Any commentary or
	// thinking appended to the session cancels all pending acknowledgements, including
	// those for other delegations. Defaults to 400. Updates apply to subsequent
	// delegations.
	AcknowledgementDelayMs param.Opt[int64] `json:"acknowledgement_delay_ms,omitzero"`
	// This field can be elided, and will marshal its zero value as "decisions".
	Type constant.Decisions `json:"type" default:"decisions"`
	paramObj
}

func (r MediaSessionConfigDelegationDecisionsParam) MarshalJSON() (data []byte, err error) {
	type shadow MediaSessionConfigDelegationDecisionsParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *MediaSessionConfigDelegationDecisionsParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Optional overrides for a stored Live session. Omitted settings are inherited.
// The model, voice, frontend instructions, and prior conversation come from the
// stored session. WebRTC negotiates its audio format; audio.format is only
// supported on WebSocket forks.
type MediaSessionForkConfigParam struct {
	// Whether to store the forked session. Omission inherits the stored session's
	// setting.
	Store param.Opt[bool] `json:"store,omitzero"`
	// Startup-only capabilities for an untrusted frontend attached to a unified WebRTC
	// session. Trusted sideband connections are unaffected.
	Client ClientConfigParam `json:"client,omitzero"`
	// Update the Responses backend for an existing Live session without changing
	// delegation ownership.
	Delegation MediaSessionForkConfigDelegationParam `json:"delegation,omitzero"`
	paramObj
}

func (r MediaSessionForkConfigParam) MarshalJSON() (data []byte, err error) {
	type shadow MediaSessionForkConfigParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *MediaSessionForkConfigParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Update the Responses backend for an existing Live session without changing
// delegation ownership.
//
// The property Type is required.
type MediaSessionForkConfigDelegationParam struct {
	// Milliseconds to wait before requesting an automatic spoken acknowledgement for a
	// delegation. The model may still speak before this delay. Any commentary or
	// thinking appended to the session cancels all pending acknowledgements, including
	// those for other delegations. Defaults to 400 for Decisions sessions and 0
	// (request immediately) for other sessions. Omission in a session update preserves
	// the current value. Updates apply to subsequent delegations.
	AcknowledgementDelayMs param.Opt[int64] `json:"acknowledgement_delay_ms,omitzero"`
	// Responses backend settings to update. Omitted settings keep their existing
	// values.
	Responses ResponsesDelegationUpdateConfigParam `json:"responses,omitzero"`
	// The delegation owner. Always `responses` for tasks handled by the Responses API.
	//
	// This field can be elided, and will marshal its zero value as "responses".
	Type constant.Responses `json:"type" default:"responses"`
	paramObj
}

func (r MediaSessionForkConfigDelegationParam) MarshalJSON() (data []byte, err error) {
	type shadow MediaSessionForkConfigDelegationParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *MediaSessionForkConfigDelegationParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// An audio chunk generated by the Live model. Decode and play primary WebSocket
// chunks in delivery order using the configured session audio format. Sideband
// connections receive reflected output audio with timestamps.
type OutputAudioDeltaEvent struct {
	// Base64-encoded raw audio. Primary WebSocket events use the session's configured
	// format; reflected sideband events use mono PCM16LE at 24 kHz.
	Delta string `json:"delta" api:"required"`
	// The event type, always `session.output_audio.delta`.
	Type constant.SessionOutputAudioDelta `json:"type" default:"session.output_audio.delta"`
	// Exclusive session-relative end in milliseconds. Required on reflected sideband
	// events; omitted on the primary WebSocket. Dropped output frames leave gaps
	// between reflected ranges.
	EndMs int64 `json:"end_ms"`
	// Inclusive session-relative start in milliseconds. Required on reflected sideband
	// events; omitted on the primary WebSocket.
	StartMs int64 `json:"start_ms"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Delta       respjson.Field
		Type        respjson.Field
		EndMs       respjson.Field
		StartMs     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r OutputAudioDeltaEvent) RawJSON() string { return r.JSON.raw }
func (r *OutputAudioDeltaEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Play one complete audio fixture through the session output. Send this event over
// the primary WebSocket or a server sideband. Audio must be raw mono PCM16
// little-endian at 24 kHz, regardless of the session audio format, and at most 30
// seconds. The entire payload is validated before model speech is interrupted and
// input is muted. Only one fixture can be active. fixture.playing acknowledges
// acceptance and fixture.played reports server mixer drain. These notifications
// are not exact transport-media boundaries. Input arriving during playback is
// dropped and not replayed.
//
// The properties Audio, Type are required.
type OutputAudioFixturePlayEventParam struct {
	// Base64-encoded raw mono PCM16 audio, little-endian, at 24 kHz. Do not include a
	// WAV or other container header.
	Audio string `json:"audio" api:"required"`
	// Optional client identifier for correlating this command with a server event's
	// client_event_id or error.client_event_id.
	EventID param.Opt[string] `json:"event_id,omitzero"`
	// This field can be elided, and will marshal its zero value as
	// "session.output_audio.fixture.play".
	Type constant.SessionOutputAudioFixturePlay `json:"type" default:"session.output_audio.fixture.play"`
	paramObj
}

func (r OutputAudioFixturePlayEventParam) MarshalJSON() (data []byte, err error) {
	type shadow OutputAudioFixturePlayEventParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *OutputAudioFixturePlayEventParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The complete audio fixture has drained from the server mixer and normal
// inference has resumed. Transport buffers may still contain audio; this event
// does not prove listener-side completion. The optional client_event_id echoes the
// play command. This event is not emitted for canceled or failed playback.
type OutputAudioFixturePlayedEvent struct {
	// The unique ID of the Live server event.
	EventID string                                   `json:"event_id" api:"required"`
	Type    constant.SessionOutputAudioFixturePlayed `json:"type" default:"session.output_audio.fixture.played"`
	// The event_id of the client command associated with this server event, when
	// supplied.
	ClientEventID string `json:"client_event_id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventID       respjson.Field
		Type          respjson.Field
		ClientEventID respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r OutputAudioFixturePlayedEvent) RawJSON() string { return r.JSON.raw }
func (r *OutputAudioFixturePlayedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The existing playback path has accepted the validated audio fixture. This
// acknowledgment does not identify when the listener receives the first sample.
// The optional client_event_id echoes the play command.
type OutputAudioFixturePlayingEvent struct {
	// The unique ID of the Live server event.
	EventID string                                    `json:"event_id" api:"required"`
	Type    constant.SessionOutputAudioFixturePlaying `json:"type" default:"session.output_audio.fixture.playing"`
	// The event_id of the client command associated with this server event, when
	// supplied.
	ClientEventID string `json:"client_event_id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventID       respjson.Field
		Type          respjson.Field
		ClientEventID respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r OutputAudioFixturePlayingEvent) RawJSON() string { return r.JSON.raw }
func (r *OutputAudioFixturePlayingEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A transcript fragment for assistant output audio in the Live session. Accumulate
// fragments in delivery order; these events do not define complete turns or
// include a transcript-done event.
type OutputTranscriptDeltaEvent struct {
	// The transcript text fragment for the audio in this time range. Append fragments
	// in delivery order to build the transcript.
	Delta string `json:"delta" api:"required"`
	// The end of this event on the Live session timeline, in milliseconds from the
	// beginning of the session. For appended context, this can equal start_ms.
	EndMs int64 `json:"end_ms" api:"required"`
	// The unique ID of the Live server event.
	EventID string `json:"event_id" api:"required"`
	// The start of this event on the Live session timeline, in milliseconds from the
	// beginning of the session.
	StartMs int64 `json:"start_ms" api:"required"`
	// The event type, always `session.output_transcript.delta`.
	Type constant.SessionOutputTranscriptDelta `json:"type" default:"session.output_transcript.delta"`
	// The event_id of the client command associated with this server event, when
	// supplied.
	ClientEventID string `json:"client_event_id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Delta         respjson.Field
		EndMs         respjson.Field
		EventID       respjson.Field
		StartMs       respjson.Field
		Type          respjson.Field
		ClientEventID respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r OutputTranscriptDeltaEvent) RawJSON() string { return r.JSON.raw }
func (r *OutputTranscriptDeltaEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Request a response from the Live session’s Responses backend, or continue a
// delegated response waiting for tool results. Requires Responses delegation.
//
// The property Type is required.
type ResponseCreateEventParam struct {
	// Optional client identifier for correlating this command with a server event's
	// client_event_id or error.client_event_id.
	EventID param.Opt[string] `json:"event_id,omitzero"`
	// The Live client event type. Always `response.create`.
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

// A streaming Responses API event from a backend delegated to by the Live session.
// Use the outer delegation_id to associate the nested stream with its Live
// delegation.
type ResponseEvent struct {
	// The nested Responses streaming event. Dispatch on its type field. Response
	// lifecycle snapshots omit input and clear instructions, tools, and output to keep
	// messages small; consume granular output events for the generated content.
	Event map[string]any `json:"event" api:"required"`
	// The unique ID of the Live server event.
	EventID string `json:"event_id" api:"required"`
	// The event type, always `response.event`.
	Type constant.ResponseEvent `json:"type" default:"response.event"`
	// The event_id of the client command associated with this server event, when
	// supplied.
	ClientEventID string `json:"client_event_id"`
	// The Live delegation associated with the nested Responses event. May be null or
	// omitted when the event cannot be correlated with a delegation.
	DelegationID string `json:"delegation_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Event         respjson.Field
		EventID       respjson.Field
		Type          respjson.Field
		ClientEventID respjson.Field
		DelegationID  respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponseEvent) RawJSON() string { return r.JSON.raw }
func (r *ResponseEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Add an input item to the Live session’s Responses backend. Requires Responses
// delegation; use `response.create` to request a response.
//
// The properties Item, Type are required.
type ResponseItemCreateEventParam struct {
	// An input item to append to the Responses backend conversation, such as a user
	// message or a function tool result.
	Item responses.ResponseInputItemUnionParam `json:"item,omitzero" api:"required"`
	// Optional client identifier for correlating this command with a server event's
	// client_event_id or error.client_event_id.
	EventID param.Opt[string] `json:"event_id,omitzero"`
	// The Live client event type. Always `response.item.create`.
	//
	// This field can be elided, and will marshal its zero value as
	// "response.item.create".
	Type constant.ResponseItemCreate `json:"type" default:"response.item.create"`
	paramObj
}

func (r ResponseItemCreateEventParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponseItemCreateEventParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponseItemCreateEventParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Model, prompt, and tool settings for tasks delegated by the Live session to a
// Responses backend.
type ResponsesDelegationConfig struct {
	// The model used for server-owned Responses delegations.
	Model string `json:"model" api:"required"`
	// Instructions for the delegated Responses model, separate from Live instructions.
	// See
	// [backend prompting](https://developers.openai.com/api/docs/guides/live-delegation#start-with-your-existing-backend-prompt).
	Instructions string `json:"instructions" api:"nullable"`
	// Maximum number of output tokens for each delegated response.
	MaxOutputTokens int64 `json:"max_output_tokens" api:"nullable"`
	// Whether the delegated Responses model may request multiple tool calls in a
	// single response.
	ParallelToolCalls bool `json:"parallel_tool_calls" api:"nullable"`
	// Reasoning settings passed to each delegated Responses request.
	Reasoning ResponsesDelegationConfigReasoning `json:"reasoning" api:"nullable"`
	// Service tier for delegated Responses requests.
	//
	// Any of "auto", "default", "fast_tier_temp_pilot", "flex", "priority",
	// "ultrafast".
	ServiceTier ResponsesDelegationConfigServiceTier `json:"service_tier" api:"nullable"`
	// Text generation settings passed to each delegated Responses request.
	Text ResponsesDelegationConfigText `json:"text" api:"nullable"`
	// Controls which tool the Responses backend uses when handling a task delegated by
	// the Live model.
	ToolChoice ResponsesDelegationConfigToolChoiceUnion `json:"tool_choice"`
	// Tools available to the Responses backend while it handles tasks delegated by the
	// Live model.
	Tools []ResponsesDelegationConfigToolUnion `json:"tools"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Model             respjson.Field
		Instructions      respjson.Field
		MaxOutputTokens   respjson.Field
		ParallelToolCalls respjson.Field
		Reasoning         respjson.Field
		ServiceTier       respjson.Field
		Text              respjson.Field
		ToolChoice        respjson.Field
		Tools             respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponsesDelegationConfig) RawJSON() string { return r.JSON.raw }
func (r *ResponsesDelegationConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this ResponsesDelegationConfig to a
// ResponsesDelegationConfigParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// ResponsesDelegationConfigParam.Overrides()
func (r ResponsesDelegationConfig) ToParam() ResponsesDelegationConfigParam {
	return param.Override[ResponsesDelegationConfigParam](json.RawMessage(r.RawJSON()))
}

// Reasoning settings passed to each delegated Responses request.
type ResponsesDelegationConfigReasoning struct {
	// How much reasoning effort the delegated Responses model should use. Supported
	// values depend on the backend model.
	//
	// Any of "none", "minimal", "low", "medium", "high", "xhigh".
	Effort string `json:"effort" api:"nullable"`
	// The reasoning summary to request from the delegated Responses model, when
	// supported.
	//
	// Any of "concise", "detailed", "auto".
	Summary string `json:"summary" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Effort      respjson.Field
		Summary     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponsesDelegationConfigReasoning) RawJSON() string { return r.JSON.raw }
func (r *ResponsesDelegationConfigReasoning) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Service tier for delegated Responses requests.
type ResponsesDelegationConfigServiceTier string

const (
	ResponsesDelegationConfigServiceTierAuto              ResponsesDelegationConfigServiceTier = "auto"
	ResponsesDelegationConfigServiceTierDefault           ResponsesDelegationConfigServiceTier = "default"
	ResponsesDelegationConfigServiceTierFastTierTempPilot ResponsesDelegationConfigServiceTier = "fast_tier_temp_pilot"
	ResponsesDelegationConfigServiceTierFlex              ResponsesDelegationConfigServiceTier = "flex"
	ResponsesDelegationConfigServiceTierPriority          ResponsesDelegationConfigServiceTier = "priority"
	ResponsesDelegationConfigServiceTierUltrafast         ResponsesDelegationConfigServiceTier = "ultrafast"
)

// Text generation settings passed to each delegated Responses request.
type ResponsesDelegationConfigText struct {
	// The amount of detail in text generated by the Responses backend. This does not
	// configure the Live model’s spoken delivery.
	//
	// Any of "low", "medium", "high".
	Verbosity string `json:"verbosity" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Verbosity   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponsesDelegationConfigText) RawJSON() string { return r.JSON.raw }
func (r *ResponsesDelegationConfigText) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ResponsesDelegationConfigToolChoiceUnion contains all possible properties and
// values from [string], [ResponsesDelegationConfigToolChoiceFunction],
// [ResponsesDelegationConfigToolChoiceMcp],
// [ResponsesDelegationConfigToolChoiceFileSearch],
// [ResponsesDelegationConfigToolChoiceWebSearch],
// [ResponsesDelegationConfigToolChoiceWebSearchPreview],
// [ResponsesDelegationConfigToolChoiceImageGeneration],
// [ResponsesDelegationConfigToolChoiceComputer],
// [ResponsesDelegationConfigToolChoiceCodeInterpreter],
// [ResponsesDelegationConfigToolChoiceProgrammaticToolCalling],
// [ResponsesDelegationConfigToolChoiceShell],
// [ResponsesDelegationConfigToolChoiceCustom],
// [ResponsesDelegationConfigToolChoiceApplyPatch],
// [ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParam].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfString]
type ResponsesDelegationConfigToolChoiceUnion struct {
	// This field will be present if the value is a [string] instead of an object.
	OfString                                                         string                                                         `json:",inline"`
	OfResponsesDelegationConfigToolChoiceFunction                    ResponsesDelegationConfigToolChoiceFunction                    `json:",inline"`
	OfResponsesDelegationConfigToolChoiceMcp                         ResponsesDelegationConfigToolChoiceMcp                         `json:",inline"`
	OfResponsesDelegationConfigToolChoiceFileSearch                  ResponsesDelegationConfigToolChoiceFileSearch                  `json:",inline"`
	OfResponsesDelegationConfigToolChoiceWebSearch                   ResponsesDelegationConfigToolChoiceWebSearch                   `json:",inline"`
	OfResponsesDelegationConfigToolChoiceWebSearchPreview            ResponsesDelegationConfigToolChoiceWebSearchPreview            `json:",inline"`
	OfResponsesDelegationConfigToolChoiceImageGeneration             ResponsesDelegationConfigToolChoiceImageGeneration             `json:",inline"`
	OfResponsesDelegationConfigToolChoiceComputer                    ResponsesDelegationConfigToolChoiceComputer                    `json:",inline"`
	OfResponsesDelegationConfigToolChoiceCodeInterpreter             ResponsesDelegationConfigToolChoiceCodeInterpreter             `json:",inline"`
	OfResponsesDelegationConfigToolChoiceProgrammaticToolCalling     ResponsesDelegationConfigToolChoiceProgrammaticToolCalling     `json:",inline"`
	OfResponsesDelegationConfigToolChoiceShell                       ResponsesDelegationConfigToolChoiceShell                       `json:",inline"`
	OfResponsesDelegationConfigToolChoiceCustom                      ResponsesDelegationConfigToolChoiceCustom                      `json:",inline"`
	OfResponsesDelegationConfigToolChoiceApplyPatch                  ResponsesDelegationConfigToolChoiceApplyPatch                  `json:",inline"`
	OfResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParam ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParam `json:",inline"`
	JSON                                                             struct {
		OfString                                                         respjson.Field
		OfResponsesDelegationConfigToolChoiceFunction                    respjson.Field
		OfResponsesDelegationConfigToolChoiceMcp                         respjson.Field
		OfResponsesDelegationConfigToolChoiceFileSearch                  respjson.Field
		OfResponsesDelegationConfigToolChoiceWebSearch                   respjson.Field
		OfResponsesDelegationConfigToolChoiceWebSearchPreview            respjson.Field
		OfResponsesDelegationConfigToolChoiceImageGeneration             respjson.Field
		OfResponsesDelegationConfigToolChoiceComputer                    respjson.Field
		OfResponsesDelegationConfigToolChoiceCodeInterpreter             respjson.Field
		OfResponsesDelegationConfigToolChoiceProgrammaticToolCalling     respjson.Field
		OfResponsesDelegationConfigToolChoiceShell                       respjson.Field
		OfResponsesDelegationConfigToolChoiceCustom                      respjson.Field
		OfResponsesDelegationConfigToolChoiceApplyPatch                  respjson.Field
		OfResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParam respjson.Field
		raw                                                              string
	} `json:"-"`
}

func (u ResponsesDelegationConfigToolChoiceUnion) AsString() (v string) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ResponsesDelegationConfigToolChoiceUnion) AsResponsesDelegationConfigToolChoiceFunction() (v ResponsesDelegationConfigToolChoiceFunction) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ResponsesDelegationConfigToolChoiceUnion) AsResponsesDelegationConfigToolChoiceMcp() (v ResponsesDelegationConfigToolChoiceMcp) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ResponsesDelegationConfigToolChoiceUnion) AsResponsesDelegationConfigToolChoiceFileSearch() (v ResponsesDelegationConfigToolChoiceFileSearch) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ResponsesDelegationConfigToolChoiceUnion) AsResponsesDelegationConfigToolChoiceWebSearch() (v ResponsesDelegationConfigToolChoiceWebSearch) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ResponsesDelegationConfigToolChoiceUnion) AsResponsesDelegationConfigToolChoiceWebSearchPreview() (v ResponsesDelegationConfigToolChoiceWebSearchPreview) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ResponsesDelegationConfigToolChoiceUnion) AsResponsesDelegationConfigToolChoiceImageGeneration() (v ResponsesDelegationConfigToolChoiceImageGeneration) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ResponsesDelegationConfigToolChoiceUnion) AsResponsesDelegationConfigToolChoiceComputer() (v ResponsesDelegationConfigToolChoiceComputer) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ResponsesDelegationConfigToolChoiceUnion) AsResponsesDelegationConfigToolChoiceCodeInterpreter() (v ResponsesDelegationConfigToolChoiceCodeInterpreter) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ResponsesDelegationConfigToolChoiceUnion) AsResponsesDelegationConfigToolChoiceProgrammaticToolCalling() (v ResponsesDelegationConfigToolChoiceProgrammaticToolCalling) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ResponsesDelegationConfigToolChoiceUnion) AsResponsesDelegationConfigToolChoiceShell() (v ResponsesDelegationConfigToolChoiceShell) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ResponsesDelegationConfigToolChoiceUnion) AsResponsesDelegationConfigToolChoiceCustom() (v ResponsesDelegationConfigToolChoiceCustom) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ResponsesDelegationConfigToolChoiceUnion) AsResponsesDelegationConfigToolChoiceApplyPatch() (v ResponsesDelegationConfigToolChoiceApplyPatch) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ResponsesDelegationConfigToolChoiceUnion) AsResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParam() (v ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParam) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ResponsesDelegationConfigToolChoiceUnion) RawJSON() string { return u.JSON.raw }

func (r *ResponsesDelegationConfigToolChoiceUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ResponsesDelegationConfigToolChoiceLiveToolChoiceEnum string

const (
	ResponsesDelegationConfigToolChoiceLiveToolChoiceEnumAuto     ResponsesDelegationConfigToolChoiceLiveToolChoiceEnum = "auto"
	ResponsesDelegationConfigToolChoiceLiveToolChoiceEnumNone     ResponsesDelegationConfigToolChoiceLiveToolChoiceEnum = "none"
	ResponsesDelegationConfigToolChoiceLiveToolChoiceEnumRequired ResponsesDelegationConfigToolChoiceLiveToolChoiceEnum = "required"
)

// Forces the model to call the shell tool when a tool call is required.
type ResponsesDelegationConfigToolChoiceFunction struct {
	// The name of the function tool to call.
	Name string `json:"name" api:"required"`
	// The tool to call. Always `function`.
	Type constant.Function `json:"type" default:"function"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Name        respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponsesDelegationConfigToolChoiceFunction) RawJSON() string { return r.JSON.raw }
func (r *ResponsesDelegationConfigToolChoiceFunction) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ResponsesDelegationConfigToolChoiceMcp struct {
	// The name of the MCP tool to call. If omitted, the server may choose a default.
	Name string `json:"name" api:"nullable"`
	// The label of the MCP server to call.
	ServerLabel string `json:"server_label" api:"required"`
	// The tool to call. Always `mcp`.
	Type constant.Mcp `json:"type" default:"mcp"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Name        respjson.Field
		ServerLabel respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponsesDelegationConfigToolChoiceMcp) RawJSON() string { return r.JSON.raw }
func (r *ResponsesDelegationConfigToolChoiceMcp) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ResponsesDelegationConfigToolChoiceFileSearch struct {
	// The tool to call. Always `file_search`.
	Type constant.FileSearch `json:"type" default:"file_search"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponsesDelegationConfigToolChoiceFileSearch) RawJSON() string { return r.JSON.raw }
func (r *ResponsesDelegationConfigToolChoiceFileSearch) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ResponsesDelegationConfigToolChoiceWebSearch struct {
	// The tool to call. Always `web_search`.
	Type constant.WebSearch `json:"type" default:"web_search"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponsesDelegationConfigToolChoiceWebSearch) RawJSON() string { return r.JSON.raw }
func (r *ResponsesDelegationConfigToolChoiceWebSearch) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ResponsesDelegationConfigToolChoiceWebSearchPreview struct {
	// The tool to call. Always `web_search_preview`.
	Type constant.WebSearchPreview `json:"type" default:"web_search_preview"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponsesDelegationConfigToolChoiceWebSearchPreview) RawJSON() string { return r.JSON.raw }
func (r *ResponsesDelegationConfigToolChoiceWebSearchPreview) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ResponsesDelegationConfigToolChoiceImageGeneration struct {
	// The tool to call. Always `image_generation`.
	Type constant.ImageGeneration `json:"type" default:"image_generation"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponsesDelegationConfigToolChoiceImageGeneration) RawJSON() string { return r.JSON.raw }
func (r *ResponsesDelegationConfigToolChoiceImageGeneration) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ResponsesDelegationConfigToolChoiceComputer struct {
	// The tool to call. Always `computer`.
	Type constant.Computer `json:"type" default:"computer"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponsesDelegationConfigToolChoiceComputer) RawJSON() string { return r.JSON.raw }
func (r *ResponsesDelegationConfigToolChoiceComputer) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ResponsesDelegationConfigToolChoiceCodeInterpreter struct {
	// The tool to call. Always `code_interpreter`.
	Type constant.CodeInterpreter `json:"type" default:"code_interpreter"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponsesDelegationConfigToolChoiceCodeInterpreter) RawJSON() string { return r.JSON.raw }
func (r *ResponsesDelegationConfigToolChoiceCodeInterpreter) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ResponsesDelegationConfigToolChoiceProgrammaticToolCalling struct {
	// The tool to call. Always `programmatic_tool_calling`.
	Type constant.ProgrammaticToolCalling `json:"type" default:"programmatic_tool_calling"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponsesDelegationConfigToolChoiceProgrammaticToolCalling) RawJSON() string {
	return r.JSON.raw
}
func (r *ResponsesDelegationConfigToolChoiceProgrammaticToolCalling) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Forces the model to call the shell tool when a tool call is required.
type ResponsesDelegationConfigToolChoiceShell struct {
	// The tool to call. Always `shell`.
	Type constant.Shell `json:"type" default:"shell"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponsesDelegationConfigToolChoiceShell) RawJSON() string { return r.JSON.raw }
func (r *ResponsesDelegationConfigToolChoiceShell) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ResponsesDelegationConfigToolChoiceCustom struct {
	// The name of the custom tool to call.
	Name string `json:"name" api:"required"`
	// The tool to call. Always `custom`.
	Type constant.Custom `json:"type" default:"custom"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Name        respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponsesDelegationConfigToolChoiceCustom) RawJSON() string { return r.JSON.raw }
func (r *ResponsesDelegationConfigToolChoiceCustom) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Forces the model to call the apply_patch tool when executing a tool call.
type ResponsesDelegationConfigToolChoiceApplyPatch struct {
	// The tool to call. Always `apply_patch`.
	Type constant.ApplyPatch `json:"type" default:"apply_patch"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponsesDelegationConfigToolChoiceApplyPatch) RawJSON() string { return r.JSON.raw }
func (r *ResponsesDelegationConfigToolChoiceApplyPatch) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolUnion contains
// all possible properties and values from
// [ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolFunction],
// [ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolMcp],
// [ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolFileSearch],
// [ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearch],
// [ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearchPreview],
// [ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolImageGeneration],
// [ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolComputer],
// [ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolCodeInterpreter],
// [ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolProgrammaticToolCalling],
// [ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolShell],
// [ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolCustom],
// [ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolApplyPatch].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolUnion struct {
	OfResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolFunction                ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolFunction                `json:",inline"`
	OfResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolMcp                     ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolMcp                     `json:",inline"`
	OfResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolFileSearch              ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolFileSearch              `json:",inline"`
	OfResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearch               ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearch               `json:",inline"`
	OfResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearchPreview        ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearchPreview        `json:",inline"`
	OfResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolImageGeneration         ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolImageGeneration         `json:",inline"`
	OfResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolComputer                ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolComputer                `json:",inline"`
	OfResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolCodeInterpreter         ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolCodeInterpreter         `json:",inline"`
	OfResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolProgrammaticToolCalling ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolProgrammaticToolCalling `json:",inline"`
	OfResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolShell                   ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolShell                   `json:",inline"`
	OfResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolCustom                  ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolCustom                  `json:",inline"`
	OfResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolApplyPatch              ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolApplyPatch              `json:",inline"`
	JSON                                                                                        struct {
		OfResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolFunction                respjson.Field
		OfResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolMcp                     respjson.Field
		OfResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolFileSearch              respjson.Field
		OfResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearch               respjson.Field
		OfResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearchPreview        respjson.Field
		OfResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolImageGeneration         respjson.Field
		OfResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolComputer                respjson.Field
		OfResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolCodeInterpreter         respjson.Field
		OfResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolProgrammaticToolCalling respjson.Field
		OfResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolShell                   respjson.Field
		OfResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolCustom                  respjson.Field
		OfResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolApplyPatch              respjson.Field
		raw                                                                                         string
	} `json:"-"`
}

func (u ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolUnion) AsResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolFunction() (v ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolFunction) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolUnion) AsResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolMcp() (v ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolMcp) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolUnion) AsResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolFileSearch() (v ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolFileSearch) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolUnion) AsResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearch() (v ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearch) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolUnion) AsResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearchPreview() (v ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearchPreview) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolUnion) AsResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolImageGeneration() (v ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolImageGeneration) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolUnion) AsResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolComputer() (v ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolComputer) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolUnion) AsResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolCodeInterpreter() (v ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolCodeInterpreter) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolUnion) AsResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolProgrammaticToolCalling() (v ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolProgrammaticToolCalling) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolUnion) AsResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolShell() (v ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolShell) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolUnion) AsResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolCustom() (v ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolCustom) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolUnion) AsResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolApplyPatch() (v ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolApplyPatch) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolUnion) RawJSON() string {
	return u.JSON.raw
}

func (r *ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolFunction struct {
	// The name of the function tool to call.
	Name string `json:"name" api:"required"`
	// The tool to call. Always `function`.
	Type constant.Function `json:"type" default:"function"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Name        respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolFunction) RawJSON() string {
	return r.JSON.raw
}
func (r *ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolFunction) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolMcp struct {
	// The name of the MCP tool to call. If omitted, the server may choose a default.
	Name string `json:"name" api:"nullable"`
	// The label of the MCP server to call.
	ServerLabel string `json:"server_label" api:"required"`
	// The tool to call. Always `mcp`.
	Type constant.Mcp `json:"type" default:"mcp"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Name        respjson.Field
		ServerLabel respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolMcp) RawJSON() string {
	return r.JSON.raw
}
func (r *ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolMcp) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolFileSearch struct {
	// The tool to call. Always `file_search`.
	Type constant.FileSearch `json:"type" default:"file_search"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolFileSearch) RawJSON() string {
	return r.JSON.raw
}
func (r *ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolFileSearch) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearch struct {
	// The tool to call. Always `web_search`.
	Type constant.WebSearch `json:"type" default:"web_search"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearch) RawJSON() string {
	return r.JSON.raw
}
func (r *ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearch) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearchPreview struct {
	// The tool to call. Always `web_search_preview`.
	Type constant.WebSearchPreview `json:"type" default:"web_search_preview"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearchPreview) RawJSON() string {
	return r.JSON.raw
}
func (r *ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearchPreview) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolImageGeneration struct {
	// The tool to call. Always `image_generation`.
	Type constant.ImageGeneration `json:"type" default:"image_generation"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolImageGeneration) RawJSON() string {
	return r.JSON.raw
}
func (r *ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolImageGeneration) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolComputer struct {
	// The tool to call. Always `computer`.
	Type constant.Computer `json:"type" default:"computer"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolComputer) RawJSON() string {
	return r.JSON.raw
}
func (r *ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolComputer) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolCodeInterpreter struct {
	// The tool to call. Always `code_interpreter`.
	Type constant.CodeInterpreter `json:"type" default:"code_interpreter"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolCodeInterpreter) RawJSON() string {
	return r.JSON.raw
}
func (r *ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolCodeInterpreter) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolProgrammaticToolCalling struct {
	// The tool to call. Always `programmatic_tool_calling`.
	Type constant.ProgrammaticToolCalling `json:"type" default:"programmatic_tool_calling"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolProgrammaticToolCalling) RawJSON() string {
	return r.JSON.raw
}
func (r *ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolProgrammaticToolCalling) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Forces the model to call the shell tool when a tool call is required.
type ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolShell struct {
	// The tool to call. Always `shell`.
	Type constant.Shell `json:"type" default:"shell"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolShell) RawJSON() string {
	return r.JSON.raw
}
func (r *ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolShell) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolCustom struct {
	// The name of the custom tool to call.
	Name string `json:"name" api:"required"`
	// The tool to call. Always `custom`.
	Type constant.Custom `json:"type" default:"custom"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Name        respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolCustom) RawJSON() string {
	return r.JSON.raw
}
func (r *ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolCustom) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Forces the model to call the apply_patch tool when executing a tool call.
type ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolApplyPatch struct {
	// The tool to call. Always `apply_patch`.
	Type constant.ApplyPatch `json:"type" default:"apply_patch"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolApplyPatch) RawJSON() string {
	return r.JSON.raw
}
func (r *ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolApplyPatch) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ResponsesDelegationConfigToolUnion contains all possible properties and values
// from [FunctionTool], [ResponsesDelegationConfigToolWebSearch],
// [ResponsesDelegationConfigToolFileSearch],
// [ResponsesDelegationConfigToolCodeInterpreter],
// [ResponsesDelegationConfigToolShell],
// [ResponsesDelegationConfigToolImageGeneration],
// [ResponsesDelegationConfigToolMcp], [ResponsesDelegationConfigToolCustom],
// [ResponsesDelegationConfigToolNamespace],
// [ResponsesDelegationConfigToolToolSearch],
// [ResponsesDelegationConfigToolProgrammaticToolCalling],
// [ResponsesDelegationConfigToolComputer],
// [ResponsesDelegationConfigToolApplyPatch].
//
// Use the [ResponsesDelegationConfigToolUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ResponsesDelegationConfigToolUnion struct {
	OfFunctionTool                                         FunctionTool                                         `json:",inline"`
	OfResponsesDelegationConfigToolWebSearch               ResponsesDelegationConfigToolWebSearch               `json:",inline"`
	OfResponsesDelegationConfigToolFileSearch              ResponsesDelegationConfigToolFileSearch              `json:",inline"`
	OfResponsesDelegationConfigToolCodeInterpreter         ResponsesDelegationConfigToolCodeInterpreter         `json:",inline"`
	OfResponsesDelegationConfigToolShell                   ResponsesDelegationConfigToolShell                   `json:",inline"`
	OfResponsesDelegationConfigToolImageGeneration         ResponsesDelegationConfigToolImageGeneration         `json:",inline"`
	OfResponsesDelegationConfigToolMcp                     ResponsesDelegationConfigToolMcp                     `json:",inline"`
	OfResponsesDelegationConfigToolCustom                  ResponsesDelegationConfigToolCustom                  `json:",inline"`
	OfResponsesDelegationConfigToolNamespace               ResponsesDelegationConfigToolNamespace               `json:",inline"`
	OfResponsesDelegationConfigToolToolSearch              ResponsesDelegationConfigToolToolSearch              `json:",inline"`
	OfResponsesDelegationConfigToolProgrammaticToolCalling ResponsesDelegationConfigToolProgrammaticToolCalling `json:",inline"`
	OfResponsesDelegationConfigToolComputer                ResponsesDelegationConfigToolComputer                `json:",inline"`
	OfResponsesDelegationConfigToolApplyPatch              ResponsesDelegationConfigToolApplyPatch              `json:",inline"`
	Type                                                   string                                               `json:"type"`
	JSON                                                   struct {
		OfFunctionTool                                         respjson.Field
		OfResponsesDelegationConfigToolWebSearch               respjson.Field
		OfResponsesDelegationConfigToolFileSearch              respjson.Field
		OfResponsesDelegationConfigToolCodeInterpreter         respjson.Field
		OfResponsesDelegationConfigToolShell                   respjson.Field
		OfResponsesDelegationConfigToolImageGeneration         respjson.Field
		OfResponsesDelegationConfigToolMcp                     respjson.Field
		OfResponsesDelegationConfigToolCustom                  respjson.Field
		OfResponsesDelegationConfigToolNamespace               respjson.Field
		OfResponsesDelegationConfigToolToolSearch              respjson.Field
		OfResponsesDelegationConfigToolProgrammaticToolCalling respjson.Field
		OfResponsesDelegationConfigToolComputer                respjson.Field
		OfResponsesDelegationConfigToolApplyPatch              respjson.Field
		Type                                                   respjson.Field
		raw                                                    string
	} `json:"-"`
}

// anyResponsesDelegationConfigTool is implemented by each variant of
// [ResponsesDelegationConfigToolUnion] to add type safety for the return type of
// [ResponsesDelegationConfigToolUnion.AsAny]
type anyResponsesDelegationConfigTool interface {
	implResponsesDelegationConfigToolUnion()
}

func (FunctionTool) implResponsesDelegationConfigToolUnion()                                 {}
func (ResponsesDelegationConfigToolWebSearch) implResponsesDelegationConfigToolUnion()       {}
func (ResponsesDelegationConfigToolFileSearch) implResponsesDelegationConfigToolUnion()      {}
func (ResponsesDelegationConfigToolCodeInterpreter) implResponsesDelegationConfigToolUnion() {}
func (ResponsesDelegationConfigToolShell) implResponsesDelegationConfigToolUnion()           {}
func (ResponsesDelegationConfigToolImageGeneration) implResponsesDelegationConfigToolUnion() {}
func (ResponsesDelegationConfigToolMcp) implResponsesDelegationConfigToolUnion()             {}
func (ResponsesDelegationConfigToolCustom) implResponsesDelegationConfigToolUnion()          {}
func (ResponsesDelegationConfigToolNamespace) implResponsesDelegationConfigToolUnion()       {}
func (ResponsesDelegationConfigToolToolSearch) implResponsesDelegationConfigToolUnion()      {}
func (ResponsesDelegationConfigToolProgrammaticToolCalling) implResponsesDelegationConfigToolUnion() {
}
func (ResponsesDelegationConfigToolComputer) implResponsesDelegationConfigToolUnion()   {}
func (ResponsesDelegationConfigToolApplyPatch) implResponsesDelegationConfigToolUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := ResponsesDelegationConfigToolUnion.AsAny().(type) {
//	case live.FunctionTool:
//	case live.ResponsesDelegationConfigToolWebSearch:
//	case live.ResponsesDelegationConfigToolFileSearch:
//	case live.ResponsesDelegationConfigToolCodeInterpreter:
//	case live.ResponsesDelegationConfigToolShell:
//	case live.ResponsesDelegationConfigToolImageGeneration:
//	case live.ResponsesDelegationConfigToolMcp:
//	case live.ResponsesDelegationConfigToolCustom:
//	case live.ResponsesDelegationConfigToolNamespace:
//	case live.ResponsesDelegationConfigToolToolSearch:
//	case live.ResponsesDelegationConfigToolProgrammaticToolCalling:
//	case live.ResponsesDelegationConfigToolComputer:
//	case live.ResponsesDelegationConfigToolApplyPatch:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ResponsesDelegationConfigToolUnion) AsAny() anyResponsesDelegationConfigTool {
	switch u.Type {
	case "function":
		return u.AsFunctionTool()
	case "web_search":
		return u.AsResponsesDelegationConfigToolWebSearch()
	case "file_search":
		return u.AsResponsesDelegationConfigToolFileSearch()
	case "code_interpreter":
		return u.AsResponsesDelegationConfigToolCodeInterpreter()
	case "shell":
		return u.AsResponsesDelegationConfigToolShell()
	case "image_generation":
		return u.AsResponsesDelegationConfigToolImageGeneration()
	case "mcp":
		return u.AsResponsesDelegationConfigToolMcp()
	case "custom":
		return u.AsResponsesDelegationConfigToolCustom()
	case "namespace":
		return u.AsResponsesDelegationConfigToolNamespace()
	case "tool_search":
		return u.AsResponsesDelegationConfigToolToolSearch()
	case "programmatic_tool_calling":
		return u.AsResponsesDelegationConfigToolProgrammaticToolCalling()
	case "computer":
		return u.AsResponsesDelegationConfigToolComputer()
	case "apply_patch":
		return u.AsResponsesDelegationConfigToolApplyPatch()
	}
	return nil
}

func (u ResponsesDelegationConfigToolUnion) AsFunctionTool() (v FunctionTool) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ResponsesDelegationConfigToolUnion) AsResponsesDelegationConfigToolWebSearch() (v ResponsesDelegationConfigToolWebSearch) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ResponsesDelegationConfigToolUnion) AsResponsesDelegationConfigToolFileSearch() (v ResponsesDelegationConfigToolFileSearch) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ResponsesDelegationConfigToolUnion) AsResponsesDelegationConfigToolCodeInterpreter() (v ResponsesDelegationConfigToolCodeInterpreter) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ResponsesDelegationConfigToolUnion) AsResponsesDelegationConfigToolShell() (v ResponsesDelegationConfigToolShell) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ResponsesDelegationConfigToolUnion) AsResponsesDelegationConfigToolImageGeneration() (v ResponsesDelegationConfigToolImageGeneration) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ResponsesDelegationConfigToolUnion) AsResponsesDelegationConfigToolMcp() (v ResponsesDelegationConfigToolMcp) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ResponsesDelegationConfigToolUnion) AsResponsesDelegationConfigToolCustom() (v ResponsesDelegationConfigToolCustom) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ResponsesDelegationConfigToolUnion) AsResponsesDelegationConfigToolNamespace() (v ResponsesDelegationConfigToolNamespace) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ResponsesDelegationConfigToolUnion) AsResponsesDelegationConfigToolToolSearch() (v ResponsesDelegationConfigToolToolSearch) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ResponsesDelegationConfigToolUnion) AsResponsesDelegationConfigToolProgrammaticToolCalling() (v ResponsesDelegationConfigToolProgrammaticToolCalling) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ResponsesDelegationConfigToolUnion) AsResponsesDelegationConfigToolComputer() (v ResponsesDelegationConfigToolComputer) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ResponsesDelegationConfigToolUnion) AsResponsesDelegationConfigToolApplyPatch() (v ResponsesDelegationConfigToolApplyPatch) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ResponsesDelegationConfigToolUnion) RawJSON() string { return u.JSON.raw }

func (r *ResponsesDelegationConfigToolUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A web search tool available to the Live session’s Responses backend.
type ResponsesDelegationConfigToolWebSearch struct {
	// The tool type. Always `web_search`.
	Type constant.WebSearch `json:"type" default:"web_search"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponsesDelegationConfigToolWebSearch) RawJSON() string { return r.JSON.raw }
func (r *ResponsesDelegationConfigToolWebSearch) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ResponsesDelegationConfigToolFileSearch struct {
	Type constant.FileSearch `json:"type" default:"file_search"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponsesDelegationConfigToolFileSearch) RawJSON() string { return r.JSON.raw }
func (r *ResponsesDelegationConfigToolFileSearch) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ResponsesDelegationConfigToolCodeInterpreter struct {
	Type constant.CodeInterpreter `json:"type" default:"code_interpreter"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponsesDelegationConfigToolCodeInterpreter) RawJSON() string { return r.JSON.raw }
func (r *ResponsesDelegationConfigToolCodeInterpreter) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A Responses shell tool. Use a hosted container or return local shell results
// with response.item.create. Domain secrets are not supported.
type ResponsesDelegationConfigToolShell struct {
	Environment ResponsesDelegationConfigToolShellEnvironmentUnion `json:"environment" api:"nullable"`
	Type        constant.Shell                                     `json:"type" default:"shell"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Environment respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponsesDelegationConfigToolShell) RawJSON() string { return r.JSON.raw }
func (r *ResponsesDelegationConfigToolShell) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ResponsesDelegationConfigToolShellEnvironmentUnion contains all possible
// properties and values from
// [ResponsesDelegationConfigToolShellEnvironmentContainerAuto],
// [ResponsesDelegationConfigToolShellEnvironmentContainerReference],
// [ResponsesDelegationConfigToolShellEnvironmentLocal].
//
// Use the [ResponsesDelegationConfigToolShellEnvironmentUnion.AsAny] method to
// switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ResponsesDelegationConfigToolShellEnvironmentUnion struct {
	OfResponsesDelegationConfigToolShellEnvironmentContainerAuto ResponsesDelegationConfigToolShellEnvironmentContainerAuto `json:",inline"`
	// Any of "container_auto", "container_reference", "local".
	Type string `json:"type"`
	// This field is from variant
	// [ResponsesDelegationConfigToolShellEnvironmentContainerAuto].
	FileIDs []string `json:"file_ids"`
	// This field is from variant
	// [ResponsesDelegationConfigToolShellEnvironmentContainerAuto].
	MemoryLimit string `json:"memory_limit"`
	// This field is from variant
	// [ResponsesDelegationConfigToolShellEnvironmentContainerAuto].
	NetworkPolicy ResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyUnion `json:"network_policy"`
	// This field is a union of
	// [[]ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillUnion],
	// [[]ResponsesDelegationConfigToolShellEnvironmentLocalSkill]
	Skills                                                            ResponsesDelegationConfigToolShellEnvironmentUnionSkills        `json:"skills"`
	OfResponsesDelegationConfigToolShellEnvironmentContainerReference ResponsesDelegationConfigToolShellEnvironmentContainerReference `json:",inline"`
	// This field is from variant
	// [ResponsesDelegationConfigToolShellEnvironmentContainerReference].
	ContainerID                                          string                                             `json:"container_id"`
	OfResponsesDelegationConfigToolShellEnvironmentLocal ResponsesDelegationConfigToolShellEnvironmentLocal `json:",inline"`
	JSON                                                 struct {
		OfResponsesDelegationConfigToolShellEnvironmentContainerAuto      respjson.Field
		Type                                                              respjson.Field
		FileIDs                                                           respjson.Field
		MemoryLimit                                                       respjson.Field
		NetworkPolicy                                                     respjson.Field
		Skills                                                            respjson.Field
		OfResponsesDelegationConfigToolShellEnvironmentContainerReference respjson.Field
		ContainerID                                                       respjson.Field
		OfResponsesDelegationConfigToolShellEnvironmentLocal              respjson.Field
		raw                                                               string
	} `json:"-"`
}

// anyResponsesDelegationConfigToolShellEnvironment is implemented by each variant
// of [ResponsesDelegationConfigToolShellEnvironmentUnion] to add type safety for
// the return type of [ResponsesDelegationConfigToolShellEnvironmentUnion.AsAny]
type anyResponsesDelegationConfigToolShellEnvironment interface {
	implResponsesDelegationConfigToolShellEnvironmentUnion()
}

func (ResponsesDelegationConfigToolShellEnvironmentContainerAuto) implResponsesDelegationConfigToolShellEnvironmentUnion() {
}
func (ResponsesDelegationConfigToolShellEnvironmentContainerReference) implResponsesDelegationConfigToolShellEnvironmentUnion() {
}
func (ResponsesDelegationConfigToolShellEnvironmentLocal) implResponsesDelegationConfigToolShellEnvironmentUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := ResponsesDelegationConfigToolShellEnvironmentUnion.AsAny().(type) {
//	case live.ResponsesDelegationConfigToolShellEnvironmentContainerAuto:
//	case live.ResponsesDelegationConfigToolShellEnvironmentContainerReference:
//	case live.ResponsesDelegationConfigToolShellEnvironmentLocal:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ResponsesDelegationConfigToolShellEnvironmentUnion) AsAny() anyResponsesDelegationConfigToolShellEnvironment {
	switch u.Type {
	case "container_auto":
		return u.AsResponsesDelegationConfigToolShellEnvironmentContainerAuto()
	case "container_reference":
		return u.AsResponsesDelegationConfigToolShellEnvironmentContainerReference()
	case "local":
		return u.AsResponsesDelegationConfigToolShellEnvironmentLocal()
	}
	return nil
}

func (u ResponsesDelegationConfigToolShellEnvironmentUnion) AsResponsesDelegationConfigToolShellEnvironmentContainerAuto() (v ResponsesDelegationConfigToolShellEnvironmentContainerAuto) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ResponsesDelegationConfigToolShellEnvironmentUnion) AsResponsesDelegationConfigToolShellEnvironmentContainerReference() (v ResponsesDelegationConfigToolShellEnvironmentContainerReference) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ResponsesDelegationConfigToolShellEnvironmentUnion) AsResponsesDelegationConfigToolShellEnvironmentLocal() (v ResponsesDelegationConfigToolShellEnvironmentLocal) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ResponsesDelegationConfigToolShellEnvironmentUnion) RawJSON() string { return u.JSON.raw }

func (r *ResponsesDelegationConfigToolShellEnvironmentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ResponsesDelegationConfigToolShellEnvironmentUnionSkills is an implicit subunion
// of [ResponsesDelegationConfigToolShellEnvironmentUnion].
// ResponsesDelegationConfigToolShellEnvironmentUnionSkills provides convenient
// access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [ResponsesDelegationConfigToolShellEnvironmentUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid:
// OfResponsesDelegationConfigToolShellEnvironmentContainerAutoSkills
// OfResponsesDelegationConfigToolShellEnvironmentLocalSkills]
type ResponsesDelegationConfigToolShellEnvironmentUnionSkills struct {
	// This field will be present if the value is a
	// [[]ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillUnion] instead
	// of an object.
	OfResponsesDelegationConfigToolShellEnvironmentContainerAutoSkills []ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillUnion `json:",inline"`
	// This field will be present if the value is a
	// [[]ResponsesDelegationConfigToolShellEnvironmentLocalSkill] instead of an
	// object.
	OfResponsesDelegationConfigToolShellEnvironmentLocalSkills []ResponsesDelegationConfigToolShellEnvironmentLocalSkill `json:",inline"`
	JSON                                                       struct {
		OfResponsesDelegationConfigToolShellEnvironmentContainerAutoSkills respjson.Field
		OfResponsesDelegationConfigToolShellEnvironmentLocalSkills         respjson.Field
		raw                                                                string
	} `json:"-"`
}

func (r *ResponsesDelegationConfigToolShellEnvironmentUnionSkills) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ResponsesDelegationConfigToolShellEnvironmentContainerAuto struct {
	// An optional list of uploaded files to make available to your code.
	FileIDs []string `json:"file_ids" api:"nullable"`
	// The memory limit for the container.
	//
	// Any of "1g", "4g", "16g", "64g".
	MemoryLimit string `json:"memory_limit" api:"nullable"`
	// Network access policy for the container.
	NetworkPolicy ResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyUnion `json:"network_policy" api:"nullable"`
	// An optional list of skills referenced by id or inline data.
	Skills []ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillUnion `json:"skills" api:"nullable"`
	// Automatically creates a container for this request
	Type constant.ContainerAuto `json:"type" default:"container_auto"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		FileIDs       respjson.Field
		MemoryLimit   respjson.Field
		NetworkPolicy respjson.Field
		Skills        respjson.Field
		Type          respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponsesDelegationConfigToolShellEnvironmentContainerAuto) RawJSON() string {
	return r.JSON.raw
}
func (r *ResponsesDelegationConfigToolShellEnvironmentContainerAuto) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyUnion
// contains all possible properties and values from
// [ResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyDisabled],
// [ResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyAllowlist].
//
// Use the
// [ResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyUnion.AsAny]
// method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyUnion struct {
	OfResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyDisabled ResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyDisabled `json:",inline"`
	// Any of "disabled", "allowlist".
	Type                                                                               string                                                                           `json:"type"`
	OfResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyAllowlist ResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyAllowlist `json:",inline"`
	// This field is from variant
	// [ResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyAllowlist].
	AllowedDomains []string `json:"allowed_domains"`
	JSON           struct {
		OfResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyDisabled  respjson.Field
		Type                                                                               respjson.Field
		OfResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyAllowlist respjson.Field
		AllowedDomains                                                                     respjson.Field
		raw                                                                                string
	} `json:"-"`
}

// anyResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicy is
// implemented by each variant of
// [ResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyUnion]
// to add type safety for the return type of
// [ResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyUnion.AsAny]
type anyResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicy interface {
	implResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyUnion()
}

func (ResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyDisabled) implResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyUnion() {
}
func (ResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyAllowlist) implResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := ResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyUnion.AsAny().(type) {
//	case live.ResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyDisabled:
//	case live.ResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyAllowlist:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyUnion) AsAny() anyResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicy {
	switch u.Type {
	case "disabled":
		return u.AsResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyDisabled()
	case "allowlist":
		return u.AsResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyAllowlist()
	}
	return nil
}

func (u ResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyUnion) AsResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyDisabled() (v ResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyDisabled) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyUnion) AsResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyAllowlist() (v ResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyAllowlist) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyUnion) RawJSON() string {
	return u.JSON.raw
}

func (r *ResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyDisabled struct {
	// Disable outbound network access. Always `disabled`.
	Type constant.Disabled `json:"type" default:"disabled"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyDisabled) RawJSON() string {
	return r.JSON.raw
}
func (r *ResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyDisabled) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyAllowlist struct {
	// A list of allowed domains when type is `allowlist`.
	AllowedDomains []string `json:"allowed_domains" api:"required"`
	// Allow outbound network access only to specified domains. Always `allowlist`.
	Type constant.Allowlist `json:"type" default:"allowlist"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		AllowedDomains respjson.Field
		Type           respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyAllowlist) RawJSON() string {
	return r.JSON.raw
}
func (r *ResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyAllowlist) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillUnion contains
// all possible properties and values from
// [ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillSkillReference],
// [ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillInline].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillUnion struct {
	OfResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillSkillReference ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillSkillReference `json:",inline"`
	OfResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillInline         ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillInline         `json:",inline"`
	JSON                                                                            struct {
		OfResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillSkillReference respjson.Field
		OfResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillInline         respjson.Field
		raw                                                                             string
	} `json:"-"`
}

func (u ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillUnion) AsResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillSkillReference() (v ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillSkillReference) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillUnion) AsResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillInline() (v ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillInline) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillUnion) RawJSON() string {
	return u.JSON.raw
}

func (r *ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillSkillReference struct {
	// The ID of the referenced skill.
	SkillID string `json:"skill_id" api:"required"`
	// References a skill created with the /v1/skills endpoint.
	Type constant.SkillReference `json:"type" default:"skill_reference"`
	// Optional skill version. Use a positive integer or 'latest'. Omit for default.
	Version string `json:"version" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		SkillID     respjson.Field
		Type        respjson.Field
		Version     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillSkillReference) RawJSON() string {
	return r.JSON.raw
}
func (r *ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillSkillReference) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillInline struct {
	// The description of the skill.
	Description string `json:"description" api:"required"`
	// The name of the skill.
	Name string `json:"name" api:"required"`
	// Inline skill payload
	Source ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillInlineSource `json:"source" api:"required"`
	// Defines an inline skill for this request.
	Type constant.Inline `json:"type" default:"inline"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Description respjson.Field
		Name        respjson.Field
		Source      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillInline) RawJSON() string {
	return r.JSON.raw
}
func (r *ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillInline) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Inline skill payload
type ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillInlineSource struct {
	// Base64-encoded skill zip bundle.
	Data string `json:"data" api:"required"`
	// The media type of the inline skill payload. Must be `application/zip`.
	MediaType constant.ApplicationZip `json:"media_type" default:"application/zip"`
	// The type of the inline skill source. Must be `base64`.
	Type constant.Base64 `json:"type" default:"base64"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Data        respjson.Field
		MediaType   respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillInlineSource) RawJSON() string {
	return r.JSON.raw
}
func (r *ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillInlineSource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ResponsesDelegationConfigToolShellEnvironmentContainerReference struct {
	// The ID of the referenced container.
	ContainerID string `json:"container_id" api:"required"`
	// References a container created with the /v1/containers endpoint
	Type constant.ContainerReference `json:"type" default:"container_reference"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ContainerID respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponsesDelegationConfigToolShellEnvironmentContainerReference) RawJSON() string {
	return r.JSON.raw
}
func (r *ResponsesDelegationConfigToolShellEnvironmentContainerReference) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ResponsesDelegationConfigToolShellEnvironmentLocal struct {
	// An optional list of skills.
	Skills []ResponsesDelegationConfigToolShellEnvironmentLocalSkill `json:"skills" api:"nullable"`
	// Use a local computer environment.
	Type constant.Local `json:"type" default:"local"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Skills      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponsesDelegationConfigToolShellEnvironmentLocal) RawJSON() string { return r.JSON.raw }
func (r *ResponsesDelegationConfigToolShellEnvironmentLocal) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ResponsesDelegationConfigToolShellEnvironmentLocalSkill struct {
	// The description of the skill.
	Description string `json:"description" api:"required"`
	// The name of the skill.
	Name string `json:"name" api:"required"`
	// The path to the directory containing the skill.
	Path string `json:"path" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Description respjson.Field
		Name        respjson.Field
		Path        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponsesDelegationConfigToolShellEnvironmentLocalSkill) RawJSON() string { return r.JSON.raw }
func (r *ResponsesDelegationConfigToolShellEnvironmentLocalSkill) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ResponsesDelegationConfigToolImageGeneration struct {
	Type constant.ImageGeneration `json:"type" default:"image_generation"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponsesDelegationConfigToolImageGeneration) RawJSON() string { return r.JSON.raw }
func (r *ResponsesDelegationConfigToolImageGeneration) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ResponsesDelegationConfigToolMcp struct {
	Type constant.Mcp `json:"type" default:"mcp"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponsesDelegationConfigToolMcp) RawJSON() string { return r.JSON.raw }
func (r *ResponsesDelegationConfigToolMcp) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ResponsesDelegationConfigToolCustom struct {
	Type constant.Custom `json:"type" default:"custom"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponsesDelegationConfigToolCustom) RawJSON() string { return r.JSON.raw }
func (r *ResponsesDelegationConfigToolCustom) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ResponsesDelegationConfigToolNamespace struct {
	Type constant.Namespace `json:"type" default:"namespace"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponsesDelegationConfigToolNamespace) RawJSON() string { return r.JSON.raw }
func (r *ResponsesDelegationConfigToolNamespace) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ResponsesDelegationConfigToolToolSearch struct {
	Type constant.ToolSearch `json:"type" default:"tool_search"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponsesDelegationConfigToolToolSearch) RawJSON() string { return r.JSON.raw }
func (r *ResponsesDelegationConfigToolToolSearch) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ResponsesDelegationConfigToolProgrammaticToolCalling struct {
	Type constant.ProgrammaticToolCalling `json:"type" default:"programmatic_tool_calling"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponsesDelegationConfigToolProgrammaticToolCalling) RawJSON() string { return r.JSON.raw }
func (r *ResponsesDelegationConfigToolProgrammaticToolCalling) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ResponsesDelegationConfigToolComputer struct {
	Type constant.Computer `json:"type" default:"computer"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponsesDelegationConfigToolComputer) RawJSON() string { return r.JSON.raw }
func (r *ResponsesDelegationConfigToolComputer) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ResponsesDelegationConfigToolApplyPatch struct {
	Type constant.ApplyPatch `json:"type" default:"apply_patch"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponsesDelegationConfigToolApplyPatch) RawJSON() string { return r.JSON.raw }
func (r *ResponsesDelegationConfigToolApplyPatch) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Model, prompt, and tool settings for tasks delegated by the Live session to a
// Responses backend.
//
// The property Model is required.
type ResponsesDelegationConfigParam struct {
	// The model used for server-owned Responses delegations.
	Model string `json:"model" api:"required"`
	// Instructions for the delegated Responses model, separate from Live instructions.
	// See
	// [backend prompting](https://developers.openai.com/api/docs/guides/live-delegation#start-with-your-existing-backend-prompt).
	Instructions param.Opt[string] `json:"instructions,omitzero"`
	// Maximum number of output tokens for each delegated response.
	MaxOutputTokens param.Opt[int64] `json:"max_output_tokens,omitzero"`
	// Whether the delegated Responses model may request multiple tool calls in a
	// single response.
	ParallelToolCalls param.Opt[bool] `json:"parallel_tool_calls,omitzero"`
	// Reasoning settings passed to each delegated Responses request.
	Reasoning ResponsesDelegationConfigReasoningParam `json:"reasoning,omitzero"`
	// Service tier for delegated Responses requests.
	//
	// Any of "auto", "default", "fast_tier_temp_pilot", "flex", "priority",
	// "ultrafast".
	ServiceTier ResponsesDelegationConfigServiceTier `json:"service_tier,omitzero"`
	// Text generation settings passed to each delegated Responses request.
	Text ResponsesDelegationConfigTextParam `json:"text,omitzero"`
	// Controls which tool the Responses backend uses when handling a task delegated by
	// the Live model.
	ToolChoice ResponsesDelegationConfigToolChoiceUnionParam `json:"tool_choice,omitzero"`
	// Tools available to the Responses backend while it handles tasks delegated by the
	// Live model.
	Tools []ResponsesDelegationConfigToolUnionParam `json:"tools,omitzero"`
	paramObj
}

func (r ResponsesDelegationConfigParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationConfigParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationConfigParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Reasoning settings passed to each delegated Responses request.
type ResponsesDelegationConfigReasoningParam struct {
	// How much reasoning effort the delegated Responses model should use. Supported
	// values depend on the backend model.
	//
	// Any of "none", "minimal", "low", "medium", "high", "xhigh".
	Effort string `json:"effort,omitzero"`
	// The reasoning summary to request from the delegated Responses model, when
	// supported.
	//
	// Any of "concise", "detailed", "auto".
	Summary string `json:"summary,omitzero"`
	paramObj
}

func (r ResponsesDelegationConfigReasoningParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationConfigReasoningParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationConfigReasoningParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[ResponsesDelegationConfigReasoningParam](
		"effort", "none", "minimal", "low", "medium", "high", "xhigh",
	)
	apijson.RegisterFieldValidator[ResponsesDelegationConfigReasoningParam](
		"summary", "concise", "detailed", "auto",
	)
}

// Text generation settings passed to each delegated Responses request.
type ResponsesDelegationConfigTextParam struct {
	// The amount of detail in text generated by the Responses backend. This does not
	// configure the Live model’s spoken delivery.
	//
	// Any of "low", "medium", "high".
	Verbosity string `json:"verbosity,omitzero"`
	paramObj
}

func (r ResponsesDelegationConfigTextParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationConfigTextParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationConfigTextParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[ResponsesDelegationConfigTextParam](
		"verbosity", "low", "medium", "high",
	)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ResponsesDelegationConfigToolChoiceUnionParam struct {
	// Check if union is this variant with !param.IsOmitted(union.OfLiveToolChoiceEnum)
	OfLiveToolChoiceEnum                                             param.Opt[string]                                                `json:",omitzero,inline"`
	OfResponsesDelegationConfigToolChoiceFunction                    *ResponsesDelegationConfigToolChoiceFunctionParam                `json:",omitzero,inline"`
	OfResponsesDelegationConfigToolChoiceMcp                         *ResponsesDelegationConfigToolChoiceMcpParam                     `json:",omitzero,inline"`
	OfResponsesDelegationConfigToolChoiceFileSearch                  *ResponsesDelegationConfigToolChoiceFileSearchParam              `json:",omitzero,inline"`
	OfResponsesDelegationConfigToolChoiceWebSearch                   *ResponsesDelegationConfigToolChoiceWebSearchParam               `json:",omitzero,inline"`
	OfResponsesDelegationConfigToolChoiceWebSearchPreview            *ResponsesDelegationConfigToolChoiceWebSearchPreviewParam        `json:",omitzero,inline"`
	OfResponsesDelegationConfigToolChoiceImageGeneration             *ResponsesDelegationConfigToolChoiceImageGenerationParam         `json:",omitzero,inline"`
	OfResponsesDelegationConfigToolChoiceComputer                    *ResponsesDelegationConfigToolChoiceComputerParam                `json:",omitzero,inline"`
	OfResponsesDelegationConfigToolChoiceCodeInterpreter             *ResponsesDelegationConfigToolChoiceCodeInterpreterParam         `json:",omitzero,inline"`
	OfResponsesDelegationConfigToolChoiceProgrammaticToolCalling     *ResponsesDelegationConfigToolChoiceProgrammaticToolCallingParam `json:",omitzero,inline"`
	OfResponsesDelegationConfigToolChoiceShell                       *ResponsesDelegationConfigToolChoiceShellParam                   `json:",omitzero,inline"`
	OfResponsesDelegationConfigToolChoiceCustom                      *ResponsesDelegationConfigToolChoiceCustomParam                  `json:",omitzero,inline"`
	OfResponsesDelegationConfigToolChoiceApplyPatch                  *ResponsesDelegationConfigToolChoiceApplyPatchParam              `json:",omitzero,inline"`
	OfResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParam *ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParam  `json:",omitzero,inline"`
	paramUnion
}

func (u ResponsesDelegationConfigToolChoiceUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfLiveToolChoiceEnum,
		u.OfResponsesDelegationConfigToolChoiceFunction,
		u.OfResponsesDelegationConfigToolChoiceMcp,
		u.OfResponsesDelegationConfigToolChoiceFileSearch,
		u.OfResponsesDelegationConfigToolChoiceWebSearch,
		u.OfResponsesDelegationConfigToolChoiceWebSearchPreview,
		u.OfResponsesDelegationConfigToolChoiceImageGeneration,
		u.OfResponsesDelegationConfigToolChoiceComputer,
		u.OfResponsesDelegationConfigToolChoiceCodeInterpreter,
		u.OfResponsesDelegationConfigToolChoiceProgrammaticToolCalling,
		u.OfResponsesDelegationConfigToolChoiceShell,
		u.OfResponsesDelegationConfigToolChoiceCustom,
		u.OfResponsesDelegationConfigToolChoiceApplyPatch,
		u.OfResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParam)
}
func (u *ResponsesDelegationConfigToolChoiceUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationConfigToolChoiceUnionParam) GetServerLabel() *string {
	if vt := u.OfResponsesDelegationConfigToolChoiceMcp; vt != nil {
		return &vt.ServerLabel
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationConfigToolChoiceUnionParam) GetTools() []ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolUnionParam {
	if vt := u.OfResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParam; vt != nil {
		return vt.Tools
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationConfigToolChoiceUnionParam) GetMode() *string {
	if vt := u.OfResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParam; vt != nil {
		return &vt.Mode
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationConfigToolChoiceUnionParam) GetName() *string {
	if vt := u.OfResponsesDelegationConfigToolChoiceFunction; vt != nil {
		return (*string)(&vt.Name)
	} else if vt := u.OfResponsesDelegationConfigToolChoiceMcp; vt != nil && vt.Name.Valid() {
		return &vt.Name.Value
	} else if vt := u.OfResponsesDelegationConfigToolChoiceCustom; vt != nil {
		return (*string)(&vt.Name)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationConfigToolChoiceUnionParam) GetType() *string {
	if vt := u.OfResponsesDelegationConfigToolChoiceFunction; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfResponsesDelegationConfigToolChoiceMcp; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfResponsesDelegationConfigToolChoiceFileSearch; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfResponsesDelegationConfigToolChoiceWebSearch; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfResponsesDelegationConfigToolChoiceWebSearchPreview; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfResponsesDelegationConfigToolChoiceImageGeneration; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfResponsesDelegationConfigToolChoiceComputer; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfResponsesDelegationConfigToolChoiceCodeInterpreter; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfResponsesDelegationConfigToolChoiceProgrammaticToolCalling; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfResponsesDelegationConfigToolChoiceShell; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfResponsesDelegationConfigToolChoiceCustom; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfResponsesDelegationConfigToolChoiceApplyPatch; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParam; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Forces the model to call the shell tool when a tool call is required.
//
// The properties Name, Type are required.
type ResponsesDelegationConfigToolChoiceFunctionParam struct {
	// The name of the function tool to call.
	Name string `json:"name" api:"required"`
	// The tool to call. Always `function`.
	//
	// This field can be elided, and will marshal its zero value as "function".
	Type constant.Function `json:"type" default:"function"`
	paramObj
}

func (r ResponsesDelegationConfigToolChoiceFunctionParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationConfigToolChoiceFunctionParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationConfigToolChoiceFunctionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties ServerLabel, Type are required.
type ResponsesDelegationConfigToolChoiceMcpParam struct {
	// The label of the MCP server to call.
	ServerLabel string `json:"server_label" api:"required"`
	// The name of the MCP tool to call. If omitted, the server may choose a default.
	Name param.Opt[string] `json:"name,omitzero"`
	// The tool to call. Always `mcp`.
	//
	// This field can be elided, and will marshal its zero value as "mcp".
	Type constant.Mcp `json:"type" default:"mcp"`
	paramObj
}

func (r ResponsesDelegationConfigToolChoiceMcpParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationConfigToolChoiceMcpParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationConfigToolChoiceMcpParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewResponsesDelegationConfigToolChoiceFileSearchParam() ResponsesDelegationConfigToolChoiceFileSearchParam {
	return ResponsesDelegationConfigToolChoiceFileSearchParam{
		Type: "file_search",
	}
}

// This struct has a constant value, construct it with
// [NewResponsesDelegationConfigToolChoiceFileSearchParam].
type ResponsesDelegationConfigToolChoiceFileSearchParam struct {
	// The tool to call. Always `file_search`.
	Type constant.FileSearch `json:"type" default:"file_search"`
	paramObj
}

func (r ResponsesDelegationConfigToolChoiceFileSearchParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationConfigToolChoiceFileSearchParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationConfigToolChoiceFileSearchParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewResponsesDelegationConfigToolChoiceWebSearchParam() ResponsesDelegationConfigToolChoiceWebSearchParam {
	return ResponsesDelegationConfigToolChoiceWebSearchParam{
		Type: "web_search",
	}
}

// This struct has a constant value, construct it with
// [NewResponsesDelegationConfigToolChoiceWebSearchParam].
type ResponsesDelegationConfigToolChoiceWebSearchParam struct {
	// The tool to call. Always `web_search`.
	Type constant.WebSearch `json:"type" default:"web_search"`
	paramObj
}

func (r ResponsesDelegationConfigToolChoiceWebSearchParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationConfigToolChoiceWebSearchParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationConfigToolChoiceWebSearchParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewResponsesDelegationConfigToolChoiceWebSearchPreviewParam() ResponsesDelegationConfigToolChoiceWebSearchPreviewParam {
	return ResponsesDelegationConfigToolChoiceWebSearchPreviewParam{
		Type: "web_search_preview",
	}
}

// This struct has a constant value, construct it with
// [NewResponsesDelegationConfigToolChoiceWebSearchPreviewParam].
type ResponsesDelegationConfigToolChoiceWebSearchPreviewParam struct {
	// The tool to call. Always `web_search_preview`.
	Type constant.WebSearchPreview `json:"type" default:"web_search_preview"`
	paramObj
}

func (r ResponsesDelegationConfigToolChoiceWebSearchPreviewParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationConfigToolChoiceWebSearchPreviewParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationConfigToolChoiceWebSearchPreviewParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewResponsesDelegationConfigToolChoiceImageGenerationParam() ResponsesDelegationConfigToolChoiceImageGenerationParam {
	return ResponsesDelegationConfigToolChoiceImageGenerationParam{
		Type: "image_generation",
	}
}

// This struct has a constant value, construct it with
// [NewResponsesDelegationConfigToolChoiceImageGenerationParam].
type ResponsesDelegationConfigToolChoiceImageGenerationParam struct {
	// The tool to call. Always `image_generation`.
	Type constant.ImageGeneration `json:"type" default:"image_generation"`
	paramObj
}

func (r ResponsesDelegationConfigToolChoiceImageGenerationParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationConfigToolChoiceImageGenerationParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationConfigToolChoiceImageGenerationParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewResponsesDelegationConfigToolChoiceComputerParam() ResponsesDelegationConfigToolChoiceComputerParam {
	return ResponsesDelegationConfigToolChoiceComputerParam{
		Type: "computer",
	}
}

// This struct has a constant value, construct it with
// [NewResponsesDelegationConfigToolChoiceComputerParam].
type ResponsesDelegationConfigToolChoiceComputerParam struct {
	// The tool to call. Always `computer`.
	Type constant.Computer `json:"type" default:"computer"`
	paramObj
}

func (r ResponsesDelegationConfigToolChoiceComputerParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationConfigToolChoiceComputerParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationConfigToolChoiceComputerParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewResponsesDelegationConfigToolChoiceCodeInterpreterParam() ResponsesDelegationConfigToolChoiceCodeInterpreterParam {
	return ResponsesDelegationConfigToolChoiceCodeInterpreterParam{
		Type: "code_interpreter",
	}
}

// This struct has a constant value, construct it with
// [NewResponsesDelegationConfigToolChoiceCodeInterpreterParam].
type ResponsesDelegationConfigToolChoiceCodeInterpreterParam struct {
	// The tool to call. Always `code_interpreter`.
	Type constant.CodeInterpreter `json:"type" default:"code_interpreter"`
	paramObj
}

func (r ResponsesDelegationConfigToolChoiceCodeInterpreterParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationConfigToolChoiceCodeInterpreterParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationConfigToolChoiceCodeInterpreterParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewResponsesDelegationConfigToolChoiceProgrammaticToolCallingParam() ResponsesDelegationConfigToolChoiceProgrammaticToolCallingParam {
	return ResponsesDelegationConfigToolChoiceProgrammaticToolCallingParam{
		Type: "programmatic_tool_calling",
	}
}

// This struct has a constant value, construct it with
// [NewResponsesDelegationConfigToolChoiceProgrammaticToolCallingParam].
type ResponsesDelegationConfigToolChoiceProgrammaticToolCallingParam struct {
	// The tool to call. Always `programmatic_tool_calling`.
	Type constant.ProgrammaticToolCalling `json:"type" default:"programmatic_tool_calling"`
	paramObj
}

func (r ResponsesDelegationConfigToolChoiceProgrammaticToolCallingParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationConfigToolChoiceProgrammaticToolCallingParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationConfigToolChoiceProgrammaticToolCallingParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewResponsesDelegationConfigToolChoiceShellParam() ResponsesDelegationConfigToolChoiceShellParam {
	return ResponsesDelegationConfigToolChoiceShellParam{
		Type: "shell",
	}
}

// Forces the model to call the shell tool when a tool call is required.
//
// This struct has a constant value, construct it with
// [NewResponsesDelegationConfigToolChoiceShellParam].
type ResponsesDelegationConfigToolChoiceShellParam struct {
	// The tool to call. Always `shell`.
	Type constant.Shell `json:"type" default:"shell"`
	paramObj
}

func (r ResponsesDelegationConfigToolChoiceShellParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationConfigToolChoiceShellParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationConfigToolChoiceShellParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Name, Type are required.
type ResponsesDelegationConfigToolChoiceCustomParam struct {
	// The name of the custom tool to call.
	Name string `json:"name" api:"required"`
	// The tool to call. Always `custom`.
	//
	// This field can be elided, and will marshal its zero value as "custom".
	Type constant.Custom `json:"type" default:"custom"`
	paramObj
}

func (r ResponsesDelegationConfigToolChoiceCustomParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationConfigToolChoiceCustomParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationConfigToolChoiceCustomParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewResponsesDelegationConfigToolChoiceApplyPatchParam() ResponsesDelegationConfigToolChoiceApplyPatchParam {
	return ResponsesDelegationConfigToolChoiceApplyPatchParam{
		Type: "apply_patch",
	}
}

// Forces the model to call the apply_patch tool when executing a tool call.
//
// This struct has a constant value, construct it with
// [NewResponsesDelegationConfigToolChoiceApplyPatchParam].
type ResponsesDelegationConfigToolChoiceApplyPatchParam struct {
	// The tool to call. Always `apply_patch`.
	Type constant.ApplyPatch `json:"type" default:"apply_patch"`
	paramObj
}

func (r ResponsesDelegationConfigToolChoiceApplyPatchParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationConfigToolChoiceApplyPatchParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationConfigToolChoiceApplyPatchParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Tools, Type are required.
type ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParam struct {
	// The tools that the delegated Responses model may call.
	Tools []ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolUnionParam `json:"tools,omitzero" api:"required"`
	// How to select a tool from the allowed set.
	//
	// Any of "none", "auto", "required".
	Mode string `json:"mode,omitzero"`
	// The tool choice type. Always `allowed_tools`.
	//
	// This field can be elided, and will marshal its zero value as "allowed_tools".
	Type constant.AllowedTools `json:"type" default:"allowed_tools"`
	paramObj
}

func (r ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParam](
		"mode", "none", "auto", "required",
	)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolUnionParam struct {
	OfFunction                *ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolFunctionParam                `json:",omitzero,inline"`
	OfMcp                     *ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolMcpParam                     `json:",omitzero,inline"`
	OfFileSearch              *ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolFileSearchParam              `json:",omitzero,inline"`
	OfWebSearch               *ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearchParam               `json:",omitzero,inline"`
	OfWebSearchPreview        *ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearchPreviewParam        `json:",omitzero,inline"`
	OfImageGeneration         *ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolImageGenerationParam         `json:",omitzero,inline"`
	OfComputer                *ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolComputerParam                `json:",omitzero,inline"`
	OfCodeInterpreter         *ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolCodeInterpreterParam         `json:",omitzero,inline"`
	OfProgrammaticToolCalling *ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolProgrammaticToolCallingParam `json:",omitzero,inline"`
	OfShell                   *ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolShellParam                   `json:",omitzero,inline"`
	OfCustom                  *ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolCustomParam                  `json:",omitzero,inline"`
	OfApplyPatch              *ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolApplyPatchParam              `json:",omitzero,inline"`
	paramUnion
}

func (u ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfFunction,
		u.OfMcp,
		u.OfFileSearch,
		u.OfWebSearch,
		u.OfWebSearchPreview,
		u.OfImageGeneration,
		u.OfComputer,
		u.OfCodeInterpreter,
		u.OfProgrammaticToolCalling,
		u.OfShell,
		u.OfCustom,
		u.OfApplyPatch)
}
func (u *ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolUnionParam) GetServerLabel() *string {
	if vt := u.OfMcp; vt != nil {
		return &vt.ServerLabel
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolUnionParam) GetName() *string {
	if vt := u.OfFunction; vt != nil {
		return (*string)(&vt.Name)
	} else if vt := u.OfMcp; vt != nil && vt.Name.Valid() {
		return &vt.Name.Value
	} else if vt := u.OfCustom; vt != nil {
		return (*string)(&vt.Name)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolUnionParam) GetType() *string {
	if vt := u.OfFunction; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfMcp; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfFileSearch; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfWebSearch; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfWebSearchPreview; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfImageGeneration; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfComputer; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfCodeInterpreter; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfProgrammaticToolCalling; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfShell; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfCustom; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfApplyPatch; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolUnionParam](
		"type",
		apijson.Discriminator[ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolFunctionParam]("function"),
		apijson.Discriminator[ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolMcpParam]("mcp"),
		apijson.Discriminator[ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolFileSearchParam]("file_search"),
		apijson.Discriminator[ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearchParam]("web_search"),
		apijson.Discriminator[ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearchPreviewParam]("web_search_preview"),
		apijson.Discriminator[ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolImageGenerationParam]("image_generation"),
		apijson.Discriminator[ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolComputerParam]("computer"),
		apijson.Discriminator[ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolCodeInterpreterParam]("code_interpreter"),
		apijson.Discriminator[ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolProgrammaticToolCallingParam]("programmatic_tool_calling"),
		apijson.Discriminator[ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolShellParam]("shell"),
		apijson.Discriminator[ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolCustomParam]("custom"),
		apijson.Discriminator[ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolApplyPatchParam]("apply_patch"),
	)
}

// The properties Name, Type are required.
type ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolFunctionParam struct {
	// The name of the function tool to call.
	Name string `json:"name" api:"required"`
	// The tool to call. Always `function`.
	//
	// This field can be elided, and will marshal its zero value as "function".
	Type constant.Function `json:"type" default:"function"`
	paramObj
}

func (r ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolFunctionParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolFunctionParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolFunctionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties ServerLabel, Type are required.
type ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolMcpParam struct {
	// The label of the MCP server to call.
	ServerLabel string `json:"server_label" api:"required"`
	// The name of the MCP tool to call. If omitted, the server may choose a default.
	Name param.Opt[string] `json:"name,omitzero"`
	// The tool to call. Always `mcp`.
	//
	// This field can be elided, and will marshal its zero value as "mcp".
	Type constant.Mcp `json:"type" default:"mcp"`
	paramObj
}

func (r ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolMcpParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolMcpParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolMcpParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolFileSearchParam() ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolFileSearchParam {
	return ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolFileSearchParam{
		Type: "file_search",
	}
}

// This struct has a constant value, construct it with
// [NewResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolFileSearchParam].
type ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolFileSearchParam struct {
	// The tool to call. Always `file_search`.
	Type constant.FileSearch `json:"type" default:"file_search"`
	paramObj
}

func (r ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolFileSearchParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolFileSearchParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolFileSearchParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearchParam() ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearchParam {
	return ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearchParam{
		Type: "web_search",
	}
}

// This struct has a constant value, construct it with
// [NewResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearchParam].
type ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearchParam struct {
	// The tool to call. Always `web_search`.
	Type constant.WebSearch `json:"type" default:"web_search"`
	paramObj
}

func (r ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearchParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearchParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearchParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearchPreviewParam() ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearchPreviewParam {
	return ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearchPreviewParam{
		Type: "web_search_preview",
	}
}

// This struct has a constant value, construct it with
// [NewResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearchPreviewParam].
type ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearchPreviewParam struct {
	// The tool to call. Always `web_search_preview`.
	Type constant.WebSearchPreview `json:"type" default:"web_search_preview"`
	paramObj
}

func (r ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearchPreviewParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearchPreviewParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearchPreviewParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolImageGenerationParam() ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolImageGenerationParam {
	return ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolImageGenerationParam{
		Type: "image_generation",
	}
}

// This struct has a constant value, construct it with
// [NewResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolImageGenerationParam].
type ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolImageGenerationParam struct {
	// The tool to call. Always `image_generation`.
	Type constant.ImageGeneration `json:"type" default:"image_generation"`
	paramObj
}

func (r ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolImageGenerationParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolImageGenerationParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolImageGenerationParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolComputerParam() ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolComputerParam {
	return ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolComputerParam{
		Type: "computer",
	}
}

// This struct has a constant value, construct it with
// [NewResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolComputerParam].
type ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolComputerParam struct {
	// The tool to call. Always `computer`.
	Type constant.Computer `json:"type" default:"computer"`
	paramObj
}

func (r ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolComputerParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolComputerParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolComputerParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolCodeInterpreterParam() ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolCodeInterpreterParam {
	return ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolCodeInterpreterParam{
		Type: "code_interpreter",
	}
}

// This struct has a constant value, construct it with
// [NewResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolCodeInterpreterParam].
type ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolCodeInterpreterParam struct {
	// The tool to call. Always `code_interpreter`.
	Type constant.CodeInterpreter `json:"type" default:"code_interpreter"`
	paramObj
}

func (r ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolCodeInterpreterParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolCodeInterpreterParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolCodeInterpreterParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolProgrammaticToolCallingParam() ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolProgrammaticToolCallingParam {
	return ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolProgrammaticToolCallingParam{
		Type: "programmatic_tool_calling",
	}
}

// This struct has a constant value, construct it with
// [NewResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolProgrammaticToolCallingParam].
type ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolProgrammaticToolCallingParam struct {
	// The tool to call. Always `programmatic_tool_calling`.
	Type constant.ProgrammaticToolCalling `json:"type" default:"programmatic_tool_calling"`
	paramObj
}

func (r ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolProgrammaticToolCallingParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolProgrammaticToolCallingParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolProgrammaticToolCallingParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolShellParam() ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolShellParam {
	return ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolShellParam{
		Type: "shell",
	}
}

// Forces the model to call the shell tool when a tool call is required.
//
// This struct has a constant value, construct it with
// [NewResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolShellParam].
type ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolShellParam struct {
	// The tool to call. Always `shell`.
	Type constant.Shell `json:"type" default:"shell"`
	paramObj
}

func (r ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolShellParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolShellParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolShellParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Name, Type are required.
type ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolCustomParam struct {
	// The name of the custom tool to call.
	Name string `json:"name" api:"required"`
	// The tool to call. Always `custom`.
	//
	// This field can be elided, and will marshal its zero value as "custom".
	Type constant.Custom `json:"type" default:"custom"`
	paramObj
}

func (r ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolCustomParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolCustomParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolCustomParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolApplyPatchParam() ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolApplyPatchParam {
	return ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolApplyPatchParam{
		Type: "apply_patch",
	}
}

// Forces the model to call the apply_patch tool when executing a tool call.
//
// This struct has a constant value, construct it with
// [NewResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolApplyPatchParam].
type ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolApplyPatchParam struct {
	// The tool to call. Always `apply_patch`.
	Type constant.ApplyPatch `json:"type" default:"apply_patch"`
	paramObj
}

func (r ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolApplyPatchParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolApplyPatchParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolApplyPatchParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ResponsesDelegationConfigToolUnionParam struct {
	OfFunction                *FunctionToolParam                                         `json:",omitzero,inline"`
	OfWebSearch               *ResponsesDelegationConfigToolWebSearchParam               `json:",omitzero,inline"`
	OfFileSearch              *ResponsesDelegationConfigToolFileSearchParam              `json:",omitzero,inline"`
	OfCodeInterpreter         *ResponsesDelegationConfigToolCodeInterpreterParam         `json:",omitzero,inline"`
	OfShell                   *ResponsesDelegationConfigToolShellParam                   `json:",omitzero,inline"`
	OfImageGeneration         *ResponsesDelegationConfigToolImageGenerationParam         `json:",omitzero,inline"`
	OfMcp                     *ResponsesDelegationConfigToolMcpParam                     `json:",omitzero,inline"`
	OfCustom                  *ResponsesDelegationConfigToolCustomParam                  `json:",omitzero,inline"`
	OfNamespace               *ResponsesDelegationConfigToolNamespaceParam               `json:",omitzero,inline"`
	OfToolSearch              *ResponsesDelegationConfigToolToolSearchParam              `json:",omitzero,inline"`
	OfProgrammaticToolCalling *ResponsesDelegationConfigToolProgrammaticToolCallingParam `json:",omitzero,inline"`
	OfComputer                *ResponsesDelegationConfigToolComputerParam                `json:",omitzero,inline"`
	OfApplyPatch              *ResponsesDelegationConfigToolApplyPatchParam              `json:",omitzero,inline"`
	paramUnion
}

func (u ResponsesDelegationConfigToolUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfFunction,
		u.OfWebSearch,
		u.OfFileSearch,
		u.OfCodeInterpreter,
		u.OfShell,
		u.OfImageGeneration,
		u.OfMcp,
		u.OfCustom,
		u.OfNamespace,
		u.OfToolSearch,
		u.OfProgrammaticToolCalling,
		u.OfComputer,
		u.OfApplyPatch)
}
func (u *ResponsesDelegationConfigToolUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationConfigToolUnionParam) GetName() *string {
	if vt := u.OfFunction; vt != nil {
		return &vt.Name
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationConfigToolUnionParam) GetDescription() *string {
	if vt := u.OfFunction; vt != nil && vt.Description.Valid() {
		return &vt.Description.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationConfigToolUnionParam) GetParameters() map[string]any {
	if vt := u.OfFunction; vt != nil {
		return vt.Parameters
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationConfigToolUnionParam) GetStrict() *bool {
	if vt := u.OfFunction; vt != nil && vt.Strict.Valid() {
		return &vt.Strict.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationConfigToolUnionParam) GetEnvironment() *ResponsesDelegationConfigToolShellEnvironmentUnionParam {
	if vt := u.OfShell; vt != nil {
		return &vt.Environment
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationConfigToolUnionParam) GetType() *string {
	if vt := u.OfFunction; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfWebSearch; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfFileSearch; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfCodeInterpreter; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfShell; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfImageGeneration; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfMcp; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfCustom; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfNamespace; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfToolSearch; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfProgrammaticToolCalling; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfComputer; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfApplyPatch; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[ResponsesDelegationConfigToolUnionParam](
		"type",
		apijson.Discriminator[FunctionToolParam]("function"),
		apijson.Discriminator[ResponsesDelegationConfigToolWebSearchParam]("web_search"),
		apijson.Discriminator[ResponsesDelegationConfigToolFileSearchParam]("file_search"),
		apijson.Discriminator[ResponsesDelegationConfigToolCodeInterpreterParam]("code_interpreter"),
		apijson.Discriminator[ResponsesDelegationConfigToolShellParam]("shell"),
		apijson.Discriminator[ResponsesDelegationConfigToolImageGenerationParam]("image_generation"),
		apijson.Discriminator[ResponsesDelegationConfigToolMcpParam]("mcp"),
		apijson.Discriminator[ResponsesDelegationConfigToolCustomParam]("custom"),
		apijson.Discriminator[ResponsesDelegationConfigToolNamespaceParam]("namespace"),
		apijson.Discriminator[ResponsesDelegationConfigToolToolSearchParam]("tool_search"),
		apijson.Discriminator[ResponsesDelegationConfigToolProgrammaticToolCallingParam]("programmatic_tool_calling"),
		apijson.Discriminator[ResponsesDelegationConfigToolComputerParam]("computer"),
		apijson.Discriminator[ResponsesDelegationConfigToolApplyPatchParam]("apply_patch"),
	)
}

func NewResponsesDelegationConfigToolWebSearchParam() ResponsesDelegationConfigToolWebSearchParam {
	return ResponsesDelegationConfigToolWebSearchParam{
		Type: "web_search",
	}
}

// A web search tool available to the Live session’s Responses backend.
//
// This struct has a constant value, construct it with
// [NewResponsesDelegationConfigToolWebSearchParam].
type ResponsesDelegationConfigToolWebSearchParam struct {
	// The tool type. Always `web_search`.
	Type constant.WebSearch `json:"type" default:"web_search"`
	paramObj
}

func (r ResponsesDelegationConfigToolWebSearchParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationConfigToolWebSearchParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationConfigToolWebSearchParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewResponsesDelegationConfigToolFileSearchParam() ResponsesDelegationConfigToolFileSearchParam {
	return ResponsesDelegationConfigToolFileSearchParam{
		Type: "file_search",
	}
}

// This struct has a constant value, construct it with
// [NewResponsesDelegationConfigToolFileSearchParam].
type ResponsesDelegationConfigToolFileSearchParam struct {
	Type constant.FileSearch `json:"type" default:"file_search"`
	paramObj
}

func (r ResponsesDelegationConfigToolFileSearchParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationConfigToolFileSearchParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationConfigToolFileSearchParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewResponsesDelegationConfigToolCodeInterpreterParam() ResponsesDelegationConfigToolCodeInterpreterParam {
	return ResponsesDelegationConfigToolCodeInterpreterParam{
		Type: "code_interpreter",
	}
}

// This struct has a constant value, construct it with
// [NewResponsesDelegationConfigToolCodeInterpreterParam].
type ResponsesDelegationConfigToolCodeInterpreterParam struct {
	Type constant.CodeInterpreter `json:"type" default:"code_interpreter"`
	paramObj
}

func (r ResponsesDelegationConfigToolCodeInterpreterParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationConfigToolCodeInterpreterParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationConfigToolCodeInterpreterParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A Responses shell tool. Use a hosted container or return local shell results
// with response.item.create. Domain secrets are not supported.
//
// The property Type is required.
type ResponsesDelegationConfigToolShellParam struct {
	Environment ResponsesDelegationConfigToolShellEnvironmentUnionParam `json:"environment,omitzero"`
	// This field can be elided, and will marshal its zero value as "shell".
	Type constant.Shell `json:"type" default:"shell"`
	paramObj
}

func (r ResponsesDelegationConfigToolShellParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationConfigToolShellParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationConfigToolShellParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ResponsesDelegationConfigToolShellEnvironmentUnionParam struct {
	OfContainerAuto      *ResponsesDelegationConfigToolShellEnvironmentContainerAutoParam      `json:",omitzero,inline"`
	OfContainerReference *ResponsesDelegationConfigToolShellEnvironmentContainerReferenceParam `json:",omitzero,inline"`
	OfLocal              *ResponsesDelegationConfigToolShellEnvironmentLocalParam              `json:",omitzero,inline"`
	paramUnion
}

func (u ResponsesDelegationConfigToolShellEnvironmentUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfContainerAuto, u.OfContainerReference, u.OfLocal)
}
func (u *ResponsesDelegationConfigToolShellEnvironmentUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationConfigToolShellEnvironmentUnionParam) GetFileIDs() []string {
	if vt := u.OfContainerAuto; vt != nil {
		return vt.FileIDs
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationConfigToolShellEnvironmentUnionParam) GetMemoryLimit() *string {
	if vt := u.OfContainerAuto; vt != nil {
		return &vt.MemoryLimit
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationConfigToolShellEnvironmentUnionParam) GetNetworkPolicy() *ResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyUnionParam {
	if vt := u.OfContainerAuto; vt != nil {
		return &vt.NetworkPolicy
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationConfigToolShellEnvironmentUnionParam) GetContainerID() *string {
	if vt := u.OfContainerReference; vt != nil {
		return &vt.ContainerID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationConfigToolShellEnvironmentUnionParam) GetType() *string {
	if vt := u.OfContainerAuto; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfContainerReference; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfLocal; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u ResponsesDelegationConfigToolShellEnvironmentUnionParam) GetSkills() (res responsesDelegationConfigToolShellEnvironmentUnionParamSkills) {
	if vt := u.OfContainerAuto; vt != nil {
		res.any = &vt.Skills
	} else if vt := u.OfLocal; vt != nil {
		res.any = &vt.Skills
	}
	return
}

// Can have the runtime types
// [_[]ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillUnionParam],
// [_[]ResponsesDelegationConfigToolShellEnvironmentLocalSkillParam]
type responsesDelegationConfigToolShellEnvironmentUnionParamSkills struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *[]live.ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillUnionParam:
//	case *[]live.ResponsesDelegationConfigToolShellEnvironmentLocalSkillParam:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u responsesDelegationConfigToolShellEnvironmentUnionParamSkills) AsAny() any { return u.any }

func init() {
	apijson.RegisterUnion[ResponsesDelegationConfigToolShellEnvironmentUnionParam](
		"type",
		apijson.Discriminator[ResponsesDelegationConfigToolShellEnvironmentContainerAutoParam]("container_auto"),
		apijson.Discriminator[ResponsesDelegationConfigToolShellEnvironmentContainerReferenceParam]("container_reference"),
		apijson.Discriminator[ResponsesDelegationConfigToolShellEnvironmentLocalParam]("local"),
	)
}

// The property Type is required.
type ResponsesDelegationConfigToolShellEnvironmentContainerAutoParam struct {
	// An optional list of uploaded files to make available to your code.
	FileIDs []string `json:"file_ids,omitzero"`
	// The memory limit for the container.
	//
	// Any of "1g", "4g", "16g", "64g".
	MemoryLimit string `json:"memory_limit,omitzero"`
	// Network access policy for the container.
	NetworkPolicy ResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyUnionParam `json:"network_policy,omitzero"`
	// An optional list of skills referenced by id or inline data.
	Skills []ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillUnionParam `json:"skills,omitzero"`
	// Automatically creates a container for this request
	//
	// This field can be elided, and will marshal its zero value as "container_auto".
	Type constant.ContainerAuto `json:"type" default:"container_auto"`
	paramObj
}

func (r ResponsesDelegationConfigToolShellEnvironmentContainerAutoParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationConfigToolShellEnvironmentContainerAutoParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationConfigToolShellEnvironmentContainerAutoParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[ResponsesDelegationConfigToolShellEnvironmentContainerAutoParam](
		"memory_limit", "1g", "4g", "16g", "64g",
	)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyUnionParam struct {
	OfDisabled  *ResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyDisabledParam  `json:",omitzero,inline"`
	OfAllowlist *ResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyAllowlistParam `json:",omitzero,inline"`
	paramUnion
}

func (u ResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfDisabled, u.OfAllowlist)
}
func (u *ResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyUnionParam) GetAllowedDomains() []string {
	if vt := u.OfAllowlist; vt != nil {
		return vt.AllowedDomains
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyUnionParam) GetType() *string {
	if vt := u.OfDisabled; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfAllowlist; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[ResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyUnionParam](
		"type",
		apijson.Discriminator[ResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyDisabledParam]("disabled"),
		apijson.Discriminator[ResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyAllowlistParam]("allowlist"),
	)
}

func NewResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyDisabledParam() ResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyDisabledParam {
	return ResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyDisabledParam{
		Type: "disabled",
	}
}

// This struct has a constant value, construct it with
// [NewResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyDisabledParam].
type ResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyDisabledParam struct {
	// Disable outbound network access. Always `disabled`.
	Type constant.Disabled `json:"type" default:"disabled"`
	paramObj
}

func (r ResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyDisabledParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyDisabledParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyDisabledParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties AllowedDomains, Type are required.
type ResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyAllowlistParam struct {
	// A list of allowed domains when type is `allowlist`.
	AllowedDomains []string `json:"allowed_domains,omitzero" api:"required"`
	// Allow outbound network access only to specified domains. Always `allowlist`.
	//
	// This field can be elided, and will marshal its zero value as "allowlist".
	Type constant.Allowlist `json:"type" default:"allowlist"`
	paramObj
}

func (r ResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyAllowlistParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyAllowlistParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationConfigToolShellEnvironmentContainerAutoNetworkPolicyAllowlistParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillUnionParam struct {
	OfSkillReference *ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillSkillReferenceParam `json:",omitzero,inline"`
	OfInline         *ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillInlineParam         `json:",omitzero,inline"`
	paramUnion
}

func (u ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfSkillReference, u.OfInline)
}
func (u *ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillUnionParam) GetSkillID() *string {
	if vt := u.OfSkillReference; vt != nil {
		return &vt.SkillID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillUnionParam) GetVersion() *string {
	if vt := u.OfSkillReference; vt != nil && vt.Version.Valid() {
		return &vt.Version.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillUnionParam) GetDescription() *string {
	if vt := u.OfInline; vt != nil {
		return &vt.Description
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillUnionParam) GetName() *string {
	if vt := u.OfInline; vt != nil {
		return &vt.Name
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillUnionParam) GetSource() *ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillInlineSourceParam {
	if vt := u.OfInline; vt != nil {
		return &vt.Source
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillUnionParam) GetType() *string {
	if vt := u.OfSkillReference; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfInline; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillUnionParam](
		"type",
		apijson.Discriminator[ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillSkillReferenceParam]("skill_reference"),
		apijson.Discriminator[ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillInlineParam]("inline"),
	)
}

// The properties SkillID, Type are required.
type ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillSkillReferenceParam struct {
	// The ID of the referenced skill.
	SkillID string `json:"skill_id" api:"required"`
	// Optional skill version. Use a positive integer or 'latest'. Omit for default.
	Version param.Opt[string] `json:"version,omitzero"`
	// References a skill created with the /v1/skills endpoint.
	//
	// This field can be elided, and will marshal its zero value as "skill_reference".
	Type constant.SkillReference `json:"type" default:"skill_reference"`
	paramObj
}

func (r ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillSkillReferenceParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillSkillReferenceParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillSkillReferenceParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Description, Name, Source, Type are required.
type ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillInlineParam struct {
	// The description of the skill.
	Description string `json:"description" api:"required"`
	// The name of the skill.
	Name string `json:"name" api:"required"`
	// Inline skill payload
	Source ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillInlineSourceParam `json:"source" api:"required"`
	// Defines an inline skill for this request.
	//
	// This field can be elided, and will marshal its zero value as "inline".
	Type constant.Inline `json:"type" default:"inline"`
	paramObj
}

func (r ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillInlineParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillInlineParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillInlineParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Inline skill payload
//
// The properties Data, MediaType, Type are required.
type ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillInlineSourceParam struct {
	// Base64-encoded skill zip bundle.
	Data string `json:"data" api:"required"`
	// The media type of the inline skill payload. Must be `application/zip`.
	//
	// This field can be elided, and will marshal its zero value as "application/zip".
	MediaType constant.ApplicationZip `json:"media_type" default:"application/zip"`
	// The type of the inline skill source. Must be `base64`.
	//
	// This field can be elided, and will marshal its zero value as "base64".
	Type constant.Base64 `json:"type" default:"base64"`
	paramObj
}

func (r ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillInlineSourceParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillInlineSourceParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillInlineSourceParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties ContainerID, Type are required.
type ResponsesDelegationConfigToolShellEnvironmentContainerReferenceParam struct {
	// The ID of the referenced container.
	ContainerID string `json:"container_id" api:"required"`
	// References a container created with the /v1/containers endpoint
	//
	// This field can be elided, and will marshal its zero value as
	// "container_reference".
	Type constant.ContainerReference `json:"type" default:"container_reference"`
	paramObj
}

func (r ResponsesDelegationConfigToolShellEnvironmentContainerReferenceParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationConfigToolShellEnvironmentContainerReferenceParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationConfigToolShellEnvironmentContainerReferenceParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The property Type is required.
type ResponsesDelegationConfigToolShellEnvironmentLocalParam struct {
	// An optional list of skills.
	Skills []ResponsesDelegationConfigToolShellEnvironmentLocalSkillParam `json:"skills,omitzero"`
	// Use a local computer environment.
	//
	// This field can be elided, and will marshal its zero value as "local".
	Type constant.Local `json:"type" default:"local"`
	paramObj
}

func (r ResponsesDelegationConfigToolShellEnvironmentLocalParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationConfigToolShellEnvironmentLocalParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationConfigToolShellEnvironmentLocalParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Description, Name, Path are required.
type ResponsesDelegationConfigToolShellEnvironmentLocalSkillParam struct {
	// The description of the skill.
	Description string `json:"description" api:"required"`
	// The name of the skill.
	Name string `json:"name" api:"required"`
	// The path to the directory containing the skill.
	Path string `json:"path" api:"required"`
	paramObj
}

func (r ResponsesDelegationConfigToolShellEnvironmentLocalSkillParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationConfigToolShellEnvironmentLocalSkillParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationConfigToolShellEnvironmentLocalSkillParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewResponsesDelegationConfigToolImageGenerationParam() ResponsesDelegationConfigToolImageGenerationParam {
	return ResponsesDelegationConfigToolImageGenerationParam{
		Type: "image_generation",
	}
}

// This struct has a constant value, construct it with
// [NewResponsesDelegationConfigToolImageGenerationParam].
type ResponsesDelegationConfigToolImageGenerationParam struct {
	Type constant.ImageGeneration `json:"type" default:"image_generation"`
	paramObj
}

func (r ResponsesDelegationConfigToolImageGenerationParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationConfigToolImageGenerationParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationConfigToolImageGenerationParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewResponsesDelegationConfigToolMcpParam() ResponsesDelegationConfigToolMcpParam {
	return ResponsesDelegationConfigToolMcpParam{
		Type: "mcp",
	}
}

// This struct has a constant value, construct it with
// [NewResponsesDelegationConfigToolMcpParam].
type ResponsesDelegationConfigToolMcpParam struct {
	Type constant.Mcp `json:"type" default:"mcp"`
	paramObj
}

func (r ResponsesDelegationConfigToolMcpParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationConfigToolMcpParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationConfigToolMcpParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewResponsesDelegationConfigToolCustomParam() ResponsesDelegationConfigToolCustomParam {
	return ResponsesDelegationConfigToolCustomParam{
		Type: "custom",
	}
}

// This struct has a constant value, construct it with
// [NewResponsesDelegationConfigToolCustomParam].
type ResponsesDelegationConfigToolCustomParam struct {
	Type constant.Custom `json:"type" default:"custom"`
	paramObj
}

func (r ResponsesDelegationConfigToolCustomParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationConfigToolCustomParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationConfigToolCustomParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewResponsesDelegationConfigToolNamespaceParam() ResponsesDelegationConfigToolNamespaceParam {
	return ResponsesDelegationConfigToolNamespaceParam{
		Type: "namespace",
	}
}

// This struct has a constant value, construct it with
// [NewResponsesDelegationConfigToolNamespaceParam].
type ResponsesDelegationConfigToolNamespaceParam struct {
	Type constant.Namespace `json:"type" default:"namespace"`
	paramObj
}

func (r ResponsesDelegationConfigToolNamespaceParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationConfigToolNamespaceParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationConfigToolNamespaceParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewResponsesDelegationConfigToolToolSearchParam() ResponsesDelegationConfigToolToolSearchParam {
	return ResponsesDelegationConfigToolToolSearchParam{
		Type: "tool_search",
	}
}

// This struct has a constant value, construct it with
// [NewResponsesDelegationConfigToolToolSearchParam].
type ResponsesDelegationConfigToolToolSearchParam struct {
	Type constant.ToolSearch `json:"type" default:"tool_search"`
	paramObj
}

func (r ResponsesDelegationConfigToolToolSearchParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationConfigToolToolSearchParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationConfigToolToolSearchParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewResponsesDelegationConfigToolProgrammaticToolCallingParam() ResponsesDelegationConfigToolProgrammaticToolCallingParam {
	return ResponsesDelegationConfigToolProgrammaticToolCallingParam{
		Type: "programmatic_tool_calling",
	}
}

// This struct has a constant value, construct it with
// [NewResponsesDelegationConfigToolProgrammaticToolCallingParam].
type ResponsesDelegationConfigToolProgrammaticToolCallingParam struct {
	Type constant.ProgrammaticToolCalling `json:"type" default:"programmatic_tool_calling"`
	paramObj
}

func (r ResponsesDelegationConfigToolProgrammaticToolCallingParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationConfigToolProgrammaticToolCallingParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationConfigToolProgrammaticToolCallingParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewResponsesDelegationConfigToolComputerParam() ResponsesDelegationConfigToolComputerParam {
	return ResponsesDelegationConfigToolComputerParam{
		Type: "computer",
	}
}

// This struct has a constant value, construct it with
// [NewResponsesDelegationConfigToolComputerParam].
type ResponsesDelegationConfigToolComputerParam struct {
	Type constant.Computer `json:"type" default:"computer"`
	paramObj
}

func (r ResponsesDelegationConfigToolComputerParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationConfigToolComputerParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationConfigToolComputerParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewResponsesDelegationConfigToolApplyPatchParam() ResponsesDelegationConfigToolApplyPatchParam {
	return ResponsesDelegationConfigToolApplyPatchParam{
		Type: "apply_patch",
	}
}

// This struct has a constant value, construct it with
// [NewResponsesDelegationConfigToolApplyPatchParam].
type ResponsesDelegationConfigToolApplyPatchParam struct {
	Type constant.ApplyPatch `json:"type" default:"apply_patch"`
	paramObj
}

func (r ResponsesDelegationConfigToolApplyPatchParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationConfigToolApplyPatchParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationConfigToolApplyPatchParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Updates to the Responses backend of an existing Live session. Omitted settings
// retain their current values.
type ResponsesDelegationUpdateConfigParam struct {
	// Instructions for the delegated Responses model, separate from Live instructions.
	// See
	// [backend prompting](https://developers.openai.com/api/docs/guides/live-delegation#start-with-your-existing-backend-prompt).
	Instructions param.Opt[string] `json:"instructions,omitzero"`
	// Maximum number of output tokens for each delegated response.
	MaxOutputTokens param.Opt[int64] `json:"max_output_tokens,omitzero"`
	// Whether the delegated Responses model may request multiple tool calls in a
	// single response.
	ParallelToolCalls param.Opt[bool] `json:"parallel_tool_calls,omitzero"`
	// The Responses backend model to use for subsequent delegated requests. Omit to
	// keep the current backend model.
	Model param.Opt[string] `json:"model,omitzero"`
	// Reasoning settings passed to each delegated Responses request.
	Reasoning ResponsesDelegationUpdateConfigReasoningParam `json:"reasoning,omitzero"`
	// Service tier for delegated Responses requests.
	//
	// Any of "auto", "default", "fast_tier_temp_pilot", "flex", "priority",
	// "ultrafast".
	ServiceTier ResponsesDelegationUpdateConfigServiceTier `json:"service_tier,omitzero"`
	// Text generation settings passed to each delegated Responses request.
	Text ResponsesDelegationUpdateConfigTextParam `json:"text,omitzero"`
	// Controls which tool the Responses backend uses when handling a task delegated by
	// the Live model.
	ToolChoice ResponsesDelegationUpdateConfigToolChoiceUnionParam `json:"tool_choice,omitzero"`
	// Tools available to the Responses backend while it handles tasks delegated by the
	// Live model.
	Tools []ResponsesDelegationUpdateConfigToolUnionParam `json:"tools,omitzero"`
	paramObj
}

func (r ResponsesDelegationUpdateConfigParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationUpdateConfigParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationUpdateConfigParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Reasoning settings passed to each delegated Responses request.
type ResponsesDelegationUpdateConfigReasoningParam struct {
	// How much reasoning effort the delegated Responses model should use. Supported
	// values depend on the backend model.
	//
	// Any of "none", "minimal", "low", "medium", "high", "xhigh".
	Effort string `json:"effort,omitzero"`
	// The reasoning summary to request from the delegated Responses model, when
	// supported.
	//
	// Any of "concise", "detailed", "auto".
	Summary string `json:"summary,omitzero"`
	paramObj
}

func (r ResponsesDelegationUpdateConfigReasoningParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationUpdateConfigReasoningParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationUpdateConfigReasoningParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[ResponsesDelegationUpdateConfigReasoningParam](
		"effort", "none", "minimal", "low", "medium", "high", "xhigh",
	)
	apijson.RegisterFieldValidator[ResponsesDelegationUpdateConfigReasoningParam](
		"summary", "concise", "detailed", "auto",
	)
}

// Service tier for delegated Responses requests.
type ResponsesDelegationUpdateConfigServiceTier string

const (
	ResponsesDelegationUpdateConfigServiceTierAuto              ResponsesDelegationUpdateConfigServiceTier = "auto"
	ResponsesDelegationUpdateConfigServiceTierDefault           ResponsesDelegationUpdateConfigServiceTier = "default"
	ResponsesDelegationUpdateConfigServiceTierFastTierTempPilot ResponsesDelegationUpdateConfigServiceTier = "fast_tier_temp_pilot"
	ResponsesDelegationUpdateConfigServiceTierFlex              ResponsesDelegationUpdateConfigServiceTier = "flex"
	ResponsesDelegationUpdateConfigServiceTierPriority          ResponsesDelegationUpdateConfigServiceTier = "priority"
	ResponsesDelegationUpdateConfigServiceTierUltrafast         ResponsesDelegationUpdateConfigServiceTier = "ultrafast"
)

// Text generation settings passed to each delegated Responses request.
type ResponsesDelegationUpdateConfigTextParam struct {
	// The amount of detail in text generated by the Responses backend. This does not
	// configure the Live model’s spoken delivery.
	//
	// Any of "low", "medium", "high".
	Verbosity string `json:"verbosity,omitzero"`
	paramObj
}

func (r ResponsesDelegationUpdateConfigTextParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationUpdateConfigTextParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationUpdateConfigTextParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[ResponsesDelegationUpdateConfigTextParam](
		"verbosity", "low", "medium", "high",
	)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ResponsesDelegationUpdateConfigToolChoiceUnionParam struct {
	// Check if union is this variant with !param.IsOmitted(union.OfLiveToolChoiceEnum)
	OfLiveToolChoiceEnum                                                   param.Opt[string]                                                      `json:",omitzero,inline"`
	OfResponsesDelegationUpdateConfigToolChoiceFunction                    *ResponsesDelegationUpdateConfigToolChoiceFunctionParam                `json:",omitzero,inline"`
	OfResponsesDelegationUpdateConfigToolChoiceMcp                         *ResponsesDelegationUpdateConfigToolChoiceMcpParam                     `json:",omitzero,inline"`
	OfResponsesDelegationUpdateConfigToolChoiceFileSearch                  *ResponsesDelegationUpdateConfigToolChoiceFileSearchParam              `json:",omitzero,inline"`
	OfResponsesDelegationUpdateConfigToolChoiceWebSearch                   *ResponsesDelegationUpdateConfigToolChoiceWebSearchParam               `json:",omitzero,inline"`
	OfResponsesDelegationUpdateConfigToolChoiceWebSearchPreview            *ResponsesDelegationUpdateConfigToolChoiceWebSearchPreviewParam        `json:",omitzero,inline"`
	OfResponsesDelegationUpdateConfigToolChoiceImageGeneration             *ResponsesDelegationUpdateConfigToolChoiceImageGenerationParam         `json:",omitzero,inline"`
	OfResponsesDelegationUpdateConfigToolChoiceComputer                    *ResponsesDelegationUpdateConfigToolChoiceComputerParam                `json:",omitzero,inline"`
	OfResponsesDelegationUpdateConfigToolChoiceCodeInterpreter             *ResponsesDelegationUpdateConfigToolChoiceCodeInterpreterParam         `json:",omitzero,inline"`
	OfResponsesDelegationUpdateConfigToolChoiceProgrammaticToolCalling     *ResponsesDelegationUpdateConfigToolChoiceProgrammaticToolCallingParam `json:",omitzero,inline"`
	OfResponsesDelegationUpdateConfigToolChoiceShell                       *ResponsesDelegationUpdateConfigToolChoiceShellParam                   `json:",omitzero,inline"`
	OfResponsesDelegationUpdateConfigToolChoiceCustom                      *ResponsesDelegationUpdateConfigToolChoiceCustomParam                  `json:",omitzero,inline"`
	OfResponsesDelegationUpdateConfigToolChoiceApplyPatch                  *ResponsesDelegationUpdateConfigToolChoiceApplyPatchParam              `json:",omitzero,inline"`
	OfResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParam *ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParam  `json:",omitzero,inline"`
	paramUnion
}

func (u ResponsesDelegationUpdateConfigToolChoiceUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfLiveToolChoiceEnum,
		u.OfResponsesDelegationUpdateConfigToolChoiceFunction,
		u.OfResponsesDelegationUpdateConfigToolChoiceMcp,
		u.OfResponsesDelegationUpdateConfigToolChoiceFileSearch,
		u.OfResponsesDelegationUpdateConfigToolChoiceWebSearch,
		u.OfResponsesDelegationUpdateConfigToolChoiceWebSearchPreview,
		u.OfResponsesDelegationUpdateConfigToolChoiceImageGeneration,
		u.OfResponsesDelegationUpdateConfigToolChoiceComputer,
		u.OfResponsesDelegationUpdateConfigToolChoiceCodeInterpreter,
		u.OfResponsesDelegationUpdateConfigToolChoiceProgrammaticToolCalling,
		u.OfResponsesDelegationUpdateConfigToolChoiceShell,
		u.OfResponsesDelegationUpdateConfigToolChoiceCustom,
		u.OfResponsesDelegationUpdateConfigToolChoiceApplyPatch,
		u.OfResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParam)
}
func (u *ResponsesDelegationUpdateConfigToolChoiceUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationUpdateConfigToolChoiceUnionParam) GetServerLabel() *string {
	if vt := u.OfResponsesDelegationUpdateConfigToolChoiceMcp; vt != nil {
		return &vt.ServerLabel
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationUpdateConfigToolChoiceUnionParam) GetTools() []ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolUnionParam {
	if vt := u.OfResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParam; vt != nil {
		return vt.Tools
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationUpdateConfigToolChoiceUnionParam) GetMode() *string {
	if vt := u.OfResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParam; vt != nil {
		return &vt.Mode
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationUpdateConfigToolChoiceUnionParam) GetName() *string {
	if vt := u.OfResponsesDelegationUpdateConfigToolChoiceFunction; vt != nil {
		return (*string)(&vt.Name)
	} else if vt := u.OfResponsesDelegationUpdateConfigToolChoiceMcp; vt != nil && vt.Name.Valid() {
		return &vt.Name.Value
	} else if vt := u.OfResponsesDelegationUpdateConfigToolChoiceCustom; vt != nil {
		return (*string)(&vt.Name)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationUpdateConfigToolChoiceUnionParam) GetType() *string {
	if vt := u.OfResponsesDelegationUpdateConfigToolChoiceFunction; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfResponsesDelegationUpdateConfigToolChoiceMcp; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfResponsesDelegationUpdateConfigToolChoiceFileSearch; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfResponsesDelegationUpdateConfigToolChoiceWebSearch; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfResponsesDelegationUpdateConfigToolChoiceWebSearchPreview; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfResponsesDelegationUpdateConfigToolChoiceImageGeneration; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfResponsesDelegationUpdateConfigToolChoiceComputer; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfResponsesDelegationUpdateConfigToolChoiceCodeInterpreter; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfResponsesDelegationUpdateConfigToolChoiceProgrammaticToolCalling; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfResponsesDelegationUpdateConfigToolChoiceShell; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfResponsesDelegationUpdateConfigToolChoiceCustom; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfResponsesDelegationUpdateConfigToolChoiceApplyPatch; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParam; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

type ResponsesDelegationUpdateConfigToolChoiceLiveToolChoiceEnum string

const (
	ResponsesDelegationUpdateConfigToolChoiceLiveToolChoiceEnumAuto     ResponsesDelegationUpdateConfigToolChoiceLiveToolChoiceEnum = "auto"
	ResponsesDelegationUpdateConfigToolChoiceLiveToolChoiceEnumNone     ResponsesDelegationUpdateConfigToolChoiceLiveToolChoiceEnum = "none"
	ResponsesDelegationUpdateConfigToolChoiceLiveToolChoiceEnumRequired ResponsesDelegationUpdateConfigToolChoiceLiveToolChoiceEnum = "required"
)

// Forces the model to call the shell tool when a tool call is required.
//
// The properties Name, Type are required.
type ResponsesDelegationUpdateConfigToolChoiceFunctionParam struct {
	// The name of the function tool to call.
	Name string `json:"name" api:"required"`
	// The tool to call. Always `function`.
	//
	// This field can be elided, and will marshal its zero value as "function".
	Type constant.Function `json:"type" default:"function"`
	paramObj
}

func (r ResponsesDelegationUpdateConfigToolChoiceFunctionParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationUpdateConfigToolChoiceFunctionParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationUpdateConfigToolChoiceFunctionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties ServerLabel, Type are required.
type ResponsesDelegationUpdateConfigToolChoiceMcpParam struct {
	// The label of the MCP server to call.
	ServerLabel string `json:"server_label" api:"required"`
	// The name of the MCP tool to call. If omitted, the server may choose a default.
	Name param.Opt[string] `json:"name,omitzero"`
	// The tool to call. Always `mcp`.
	//
	// This field can be elided, and will marshal its zero value as "mcp".
	Type constant.Mcp `json:"type" default:"mcp"`
	paramObj
}

func (r ResponsesDelegationUpdateConfigToolChoiceMcpParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationUpdateConfigToolChoiceMcpParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationUpdateConfigToolChoiceMcpParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewResponsesDelegationUpdateConfigToolChoiceFileSearchParam() ResponsesDelegationUpdateConfigToolChoiceFileSearchParam {
	return ResponsesDelegationUpdateConfigToolChoiceFileSearchParam{
		Type: "file_search",
	}
}

// This struct has a constant value, construct it with
// [NewResponsesDelegationUpdateConfigToolChoiceFileSearchParam].
type ResponsesDelegationUpdateConfigToolChoiceFileSearchParam struct {
	// The tool to call. Always `file_search`.
	Type constant.FileSearch `json:"type" default:"file_search"`
	paramObj
}

func (r ResponsesDelegationUpdateConfigToolChoiceFileSearchParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationUpdateConfigToolChoiceFileSearchParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationUpdateConfigToolChoiceFileSearchParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewResponsesDelegationUpdateConfigToolChoiceWebSearchParam() ResponsesDelegationUpdateConfigToolChoiceWebSearchParam {
	return ResponsesDelegationUpdateConfigToolChoiceWebSearchParam{
		Type: "web_search",
	}
}

// This struct has a constant value, construct it with
// [NewResponsesDelegationUpdateConfigToolChoiceWebSearchParam].
type ResponsesDelegationUpdateConfigToolChoiceWebSearchParam struct {
	// The tool to call. Always `web_search`.
	Type constant.WebSearch `json:"type" default:"web_search"`
	paramObj
}

func (r ResponsesDelegationUpdateConfigToolChoiceWebSearchParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationUpdateConfigToolChoiceWebSearchParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationUpdateConfigToolChoiceWebSearchParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewResponsesDelegationUpdateConfigToolChoiceWebSearchPreviewParam() ResponsesDelegationUpdateConfigToolChoiceWebSearchPreviewParam {
	return ResponsesDelegationUpdateConfigToolChoiceWebSearchPreviewParam{
		Type: "web_search_preview",
	}
}

// This struct has a constant value, construct it with
// [NewResponsesDelegationUpdateConfigToolChoiceWebSearchPreviewParam].
type ResponsesDelegationUpdateConfigToolChoiceWebSearchPreviewParam struct {
	// The tool to call. Always `web_search_preview`.
	Type constant.WebSearchPreview `json:"type" default:"web_search_preview"`
	paramObj
}

func (r ResponsesDelegationUpdateConfigToolChoiceWebSearchPreviewParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationUpdateConfigToolChoiceWebSearchPreviewParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationUpdateConfigToolChoiceWebSearchPreviewParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewResponsesDelegationUpdateConfigToolChoiceImageGenerationParam() ResponsesDelegationUpdateConfigToolChoiceImageGenerationParam {
	return ResponsesDelegationUpdateConfigToolChoiceImageGenerationParam{
		Type: "image_generation",
	}
}

// This struct has a constant value, construct it with
// [NewResponsesDelegationUpdateConfigToolChoiceImageGenerationParam].
type ResponsesDelegationUpdateConfigToolChoiceImageGenerationParam struct {
	// The tool to call. Always `image_generation`.
	Type constant.ImageGeneration `json:"type" default:"image_generation"`
	paramObj
}

func (r ResponsesDelegationUpdateConfigToolChoiceImageGenerationParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationUpdateConfigToolChoiceImageGenerationParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationUpdateConfigToolChoiceImageGenerationParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewResponsesDelegationUpdateConfigToolChoiceComputerParam() ResponsesDelegationUpdateConfigToolChoiceComputerParam {
	return ResponsesDelegationUpdateConfigToolChoiceComputerParam{
		Type: "computer",
	}
}

// This struct has a constant value, construct it with
// [NewResponsesDelegationUpdateConfigToolChoiceComputerParam].
type ResponsesDelegationUpdateConfigToolChoiceComputerParam struct {
	// The tool to call. Always `computer`.
	Type constant.Computer `json:"type" default:"computer"`
	paramObj
}

func (r ResponsesDelegationUpdateConfigToolChoiceComputerParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationUpdateConfigToolChoiceComputerParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationUpdateConfigToolChoiceComputerParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewResponsesDelegationUpdateConfigToolChoiceCodeInterpreterParam() ResponsesDelegationUpdateConfigToolChoiceCodeInterpreterParam {
	return ResponsesDelegationUpdateConfigToolChoiceCodeInterpreterParam{
		Type: "code_interpreter",
	}
}

// This struct has a constant value, construct it with
// [NewResponsesDelegationUpdateConfigToolChoiceCodeInterpreterParam].
type ResponsesDelegationUpdateConfigToolChoiceCodeInterpreterParam struct {
	// The tool to call. Always `code_interpreter`.
	Type constant.CodeInterpreter `json:"type" default:"code_interpreter"`
	paramObj
}

func (r ResponsesDelegationUpdateConfigToolChoiceCodeInterpreterParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationUpdateConfigToolChoiceCodeInterpreterParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationUpdateConfigToolChoiceCodeInterpreterParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewResponsesDelegationUpdateConfigToolChoiceProgrammaticToolCallingParam() ResponsesDelegationUpdateConfigToolChoiceProgrammaticToolCallingParam {
	return ResponsesDelegationUpdateConfigToolChoiceProgrammaticToolCallingParam{
		Type: "programmatic_tool_calling",
	}
}

// This struct has a constant value, construct it with
// [NewResponsesDelegationUpdateConfigToolChoiceProgrammaticToolCallingParam].
type ResponsesDelegationUpdateConfigToolChoiceProgrammaticToolCallingParam struct {
	// The tool to call. Always `programmatic_tool_calling`.
	Type constant.ProgrammaticToolCalling `json:"type" default:"programmatic_tool_calling"`
	paramObj
}

func (r ResponsesDelegationUpdateConfigToolChoiceProgrammaticToolCallingParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationUpdateConfigToolChoiceProgrammaticToolCallingParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationUpdateConfigToolChoiceProgrammaticToolCallingParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewResponsesDelegationUpdateConfigToolChoiceShellParam() ResponsesDelegationUpdateConfigToolChoiceShellParam {
	return ResponsesDelegationUpdateConfigToolChoiceShellParam{
		Type: "shell",
	}
}

// Forces the model to call the shell tool when a tool call is required.
//
// This struct has a constant value, construct it with
// [NewResponsesDelegationUpdateConfigToolChoiceShellParam].
type ResponsesDelegationUpdateConfigToolChoiceShellParam struct {
	// The tool to call. Always `shell`.
	Type constant.Shell `json:"type" default:"shell"`
	paramObj
}

func (r ResponsesDelegationUpdateConfigToolChoiceShellParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationUpdateConfigToolChoiceShellParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationUpdateConfigToolChoiceShellParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Name, Type are required.
type ResponsesDelegationUpdateConfigToolChoiceCustomParam struct {
	// The name of the custom tool to call.
	Name string `json:"name" api:"required"`
	// The tool to call. Always `custom`.
	//
	// This field can be elided, and will marshal its zero value as "custom".
	Type constant.Custom `json:"type" default:"custom"`
	paramObj
}

func (r ResponsesDelegationUpdateConfigToolChoiceCustomParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationUpdateConfigToolChoiceCustomParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationUpdateConfigToolChoiceCustomParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewResponsesDelegationUpdateConfigToolChoiceApplyPatchParam() ResponsesDelegationUpdateConfigToolChoiceApplyPatchParam {
	return ResponsesDelegationUpdateConfigToolChoiceApplyPatchParam{
		Type: "apply_patch",
	}
}

// Forces the model to call the apply_patch tool when executing a tool call.
//
// This struct has a constant value, construct it with
// [NewResponsesDelegationUpdateConfigToolChoiceApplyPatchParam].
type ResponsesDelegationUpdateConfigToolChoiceApplyPatchParam struct {
	// The tool to call. Always `apply_patch`.
	Type constant.ApplyPatch `json:"type" default:"apply_patch"`
	paramObj
}

func (r ResponsesDelegationUpdateConfigToolChoiceApplyPatchParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationUpdateConfigToolChoiceApplyPatchParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationUpdateConfigToolChoiceApplyPatchParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Tools, Type are required.
type ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParam struct {
	// The tools that the delegated Responses model may call.
	Tools []ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolUnionParam `json:"tools,omitzero" api:"required"`
	// How to select a tool from the allowed set.
	//
	// Any of "none", "auto", "required".
	Mode string `json:"mode,omitzero"`
	// The tool choice type. Always `allowed_tools`.
	//
	// This field can be elided, and will marshal its zero value as "allowed_tools".
	Type constant.AllowedTools `json:"type" default:"allowed_tools"`
	paramObj
}

func (r ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParam](
		"mode", "none", "auto", "required",
	)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolUnionParam struct {
	OfFunction                *ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolFunctionParam                `json:",omitzero,inline"`
	OfMcp                     *ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolMcpParam                     `json:",omitzero,inline"`
	OfFileSearch              *ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolFileSearchParam              `json:",omitzero,inline"`
	OfWebSearch               *ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearchParam               `json:",omitzero,inline"`
	OfWebSearchPreview        *ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearchPreviewParam        `json:",omitzero,inline"`
	OfImageGeneration         *ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolImageGenerationParam         `json:",omitzero,inline"`
	OfComputer                *ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolComputerParam                `json:",omitzero,inline"`
	OfCodeInterpreter         *ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolCodeInterpreterParam         `json:",omitzero,inline"`
	OfProgrammaticToolCalling *ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolProgrammaticToolCallingParam `json:",omitzero,inline"`
	OfShell                   *ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolShellParam                   `json:",omitzero,inline"`
	OfCustom                  *ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolCustomParam                  `json:",omitzero,inline"`
	OfApplyPatch              *ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolApplyPatchParam              `json:",omitzero,inline"`
	paramUnion
}

func (u ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfFunction,
		u.OfMcp,
		u.OfFileSearch,
		u.OfWebSearch,
		u.OfWebSearchPreview,
		u.OfImageGeneration,
		u.OfComputer,
		u.OfCodeInterpreter,
		u.OfProgrammaticToolCalling,
		u.OfShell,
		u.OfCustom,
		u.OfApplyPatch)
}
func (u *ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolUnionParam) GetServerLabel() *string {
	if vt := u.OfMcp; vt != nil {
		return &vt.ServerLabel
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolUnionParam) GetName() *string {
	if vt := u.OfFunction; vt != nil {
		return (*string)(&vt.Name)
	} else if vt := u.OfMcp; vt != nil && vt.Name.Valid() {
		return &vt.Name.Value
	} else if vt := u.OfCustom; vt != nil {
		return (*string)(&vt.Name)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolUnionParam) GetType() *string {
	if vt := u.OfFunction; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfMcp; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfFileSearch; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfWebSearch; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfWebSearchPreview; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfImageGeneration; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfComputer; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfCodeInterpreter; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfProgrammaticToolCalling; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfShell; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfCustom; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfApplyPatch; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolUnionParam](
		"type",
		apijson.Discriminator[ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolFunctionParam]("function"),
		apijson.Discriminator[ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolMcpParam]("mcp"),
		apijson.Discriminator[ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolFileSearchParam]("file_search"),
		apijson.Discriminator[ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearchParam]("web_search"),
		apijson.Discriminator[ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearchPreviewParam]("web_search_preview"),
		apijson.Discriminator[ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolImageGenerationParam]("image_generation"),
		apijson.Discriminator[ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolComputerParam]("computer"),
		apijson.Discriminator[ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolCodeInterpreterParam]("code_interpreter"),
		apijson.Discriminator[ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolProgrammaticToolCallingParam]("programmatic_tool_calling"),
		apijson.Discriminator[ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolShellParam]("shell"),
		apijson.Discriminator[ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolCustomParam]("custom"),
		apijson.Discriminator[ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolApplyPatchParam]("apply_patch"),
	)
}

// The properties Name, Type are required.
type ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolFunctionParam struct {
	// The name of the function tool to call.
	Name string `json:"name" api:"required"`
	// The tool to call. Always `function`.
	//
	// This field can be elided, and will marshal its zero value as "function".
	Type constant.Function `json:"type" default:"function"`
	paramObj
}

func (r ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolFunctionParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolFunctionParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolFunctionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties ServerLabel, Type are required.
type ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolMcpParam struct {
	// The label of the MCP server to call.
	ServerLabel string `json:"server_label" api:"required"`
	// The name of the MCP tool to call. If omitted, the server may choose a default.
	Name param.Opt[string] `json:"name,omitzero"`
	// The tool to call. Always `mcp`.
	//
	// This field can be elided, and will marshal its zero value as "mcp".
	Type constant.Mcp `json:"type" default:"mcp"`
	paramObj
}

func (r ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolMcpParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolMcpParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolMcpParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolFileSearchParam() ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolFileSearchParam {
	return ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolFileSearchParam{
		Type: "file_search",
	}
}

// This struct has a constant value, construct it with
// [NewResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolFileSearchParam].
type ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolFileSearchParam struct {
	// The tool to call. Always `file_search`.
	Type constant.FileSearch `json:"type" default:"file_search"`
	paramObj
}

func (r ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolFileSearchParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolFileSearchParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolFileSearchParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearchParam() ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearchParam {
	return ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearchParam{
		Type: "web_search",
	}
}

// This struct has a constant value, construct it with
// [NewResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearchParam].
type ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearchParam struct {
	// The tool to call. Always `web_search`.
	Type constant.WebSearch `json:"type" default:"web_search"`
	paramObj
}

func (r ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearchParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearchParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearchParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearchPreviewParam() ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearchPreviewParam {
	return ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearchPreviewParam{
		Type: "web_search_preview",
	}
}

// This struct has a constant value, construct it with
// [NewResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearchPreviewParam].
type ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearchPreviewParam struct {
	// The tool to call. Always `web_search_preview`.
	Type constant.WebSearchPreview `json:"type" default:"web_search_preview"`
	paramObj
}

func (r ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearchPreviewParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearchPreviewParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolWebSearchPreviewParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolImageGenerationParam() ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolImageGenerationParam {
	return ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolImageGenerationParam{
		Type: "image_generation",
	}
}

// This struct has a constant value, construct it with
// [NewResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolImageGenerationParam].
type ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolImageGenerationParam struct {
	// The tool to call. Always `image_generation`.
	Type constant.ImageGeneration `json:"type" default:"image_generation"`
	paramObj
}

func (r ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolImageGenerationParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolImageGenerationParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolImageGenerationParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolComputerParam() ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolComputerParam {
	return ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolComputerParam{
		Type: "computer",
	}
}

// This struct has a constant value, construct it with
// [NewResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolComputerParam].
type ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolComputerParam struct {
	// The tool to call. Always `computer`.
	Type constant.Computer `json:"type" default:"computer"`
	paramObj
}

func (r ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolComputerParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolComputerParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolComputerParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolCodeInterpreterParam() ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolCodeInterpreterParam {
	return ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolCodeInterpreterParam{
		Type: "code_interpreter",
	}
}

// This struct has a constant value, construct it with
// [NewResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolCodeInterpreterParam].
type ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolCodeInterpreterParam struct {
	// The tool to call. Always `code_interpreter`.
	Type constant.CodeInterpreter `json:"type" default:"code_interpreter"`
	paramObj
}

func (r ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolCodeInterpreterParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolCodeInterpreterParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolCodeInterpreterParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolProgrammaticToolCallingParam() ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolProgrammaticToolCallingParam {
	return ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolProgrammaticToolCallingParam{
		Type: "programmatic_tool_calling",
	}
}

// This struct has a constant value, construct it with
// [NewResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolProgrammaticToolCallingParam].
type ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolProgrammaticToolCallingParam struct {
	// The tool to call. Always `programmatic_tool_calling`.
	Type constant.ProgrammaticToolCalling `json:"type" default:"programmatic_tool_calling"`
	paramObj
}

func (r ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolProgrammaticToolCallingParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolProgrammaticToolCallingParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolProgrammaticToolCallingParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolShellParam() ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolShellParam {
	return ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolShellParam{
		Type: "shell",
	}
}

// Forces the model to call the shell tool when a tool call is required.
//
// This struct has a constant value, construct it with
// [NewResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolShellParam].
type ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolShellParam struct {
	// The tool to call. Always `shell`.
	Type constant.Shell `json:"type" default:"shell"`
	paramObj
}

func (r ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolShellParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolShellParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolShellParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Name, Type are required.
type ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolCustomParam struct {
	// The name of the custom tool to call.
	Name string `json:"name" api:"required"`
	// The tool to call. Always `custom`.
	//
	// This field can be elided, and will marshal its zero value as "custom".
	Type constant.Custom `json:"type" default:"custom"`
	paramObj
}

func (r ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolCustomParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolCustomParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolCustomParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolApplyPatchParam() ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolApplyPatchParam {
	return ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolApplyPatchParam{
		Type: "apply_patch",
	}
}

// Forces the model to call the apply_patch tool when executing a tool call.
//
// This struct has a constant value, construct it with
// [NewResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolApplyPatchParam].
type ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolApplyPatchParam struct {
	// The tool to call. Always `apply_patch`.
	Type constant.ApplyPatch `json:"type" default:"apply_patch"`
	paramObj
}

func (r ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolApplyPatchParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolApplyPatchParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolApplyPatchParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ResponsesDelegationUpdateConfigToolUnionParam struct {
	OfFunction                *FunctionToolParam                                               `json:",omitzero,inline"`
	OfWebSearch               *ResponsesDelegationUpdateConfigToolWebSearchParam               `json:",omitzero,inline"`
	OfFileSearch              *ResponsesDelegationUpdateConfigToolFileSearchParam              `json:",omitzero,inline"`
	OfCodeInterpreter         *ResponsesDelegationUpdateConfigToolCodeInterpreterParam         `json:",omitzero,inline"`
	OfShell                   *ResponsesDelegationUpdateConfigToolShellParam                   `json:",omitzero,inline"`
	OfImageGeneration         *ResponsesDelegationUpdateConfigToolImageGenerationParam         `json:",omitzero,inline"`
	OfMcp                     *ResponsesDelegationUpdateConfigToolMcpParam                     `json:",omitzero,inline"`
	OfCustom                  *ResponsesDelegationUpdateConfigToolCustomParam                  `json:",omitzero,inline"`
	OfNamespace               *ResponsesDelegationUpdateConfigToolNamespaceParam               `json:",omitzero,inline"`
	OfToolSearch              *ResponsesDelegationUpdateConfigToolToolSearchParam              `json:",omitzero,inline"`
	OfProgrammaticToolCalling *ResponsesDelegationUpdateConfigToolProgrammaticToolCallingParam `json:",omitzero,inline"`
	OfComputer                *ResponsesDelegationUpdateConfigToolComputerParam                `json:",omitzero,inline"`
	OfApplyPatch              *ResponsesDelegationUpdateConfigToolApplyPatchParam              `json:",omitzero,inline"`
	paramUnion
}

func (u ResponsesDelegationUpdateConfigToolUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfFunction,
		u.OfWebSearch,
		u.OfFileSearch,
		u.OfCodeInterpreter,
		u.OfShell,
		u.OfImageGeneration,
		u.OfMcp,
		u.OfCustom,
		u.OfNamespace,
		u.OfToolSearch,
		u.OfProgrammaticToolCalling,
		u.OfComputer,
		u.OfApplyPatch)
}
func (u *ResponsesDelegationUpdateConfigToolUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationUpdateConfigToolUnionParam) GetName() *string {
	if vt := u.OfFunction; vt != nil {
		return &vt.Name
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationUpdateConfigToolUnionParam) GetDescription() *string {
	if vt := u.OfFunction; vt != nil && vt.Description.Valid() {
		return &vt.Description.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationUpdateConfigToolUnionParam) GetParameters() map[string]any {
	if vt := u.OfFunction; vt != nil {
		return vt.Parameters
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationUpdateConfigToolUnionParam) GetStrict() *bool {
	if vt := u.OfFunction; vt != nil && vt.Strict.Valid() {
		return &vt.Strict.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationUpdateConfigToolUnionParam) GetEnvironment() *ResponsesDelegationUpdateConfigToolShellEnvironmentUnionParam {
	if vt := u.OfShell; vt != nil {
		return &vt.Environment
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationUpdateConfigToolUnionParam) GetType() *string {
	if vt := u.OfFunction; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfWebSearch; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfFileSearch; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfCodeInterpreter; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfShell; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfImageGeneration; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfMcp; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfCustom; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfNamespace; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfToolSearch; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfProgrammaticToolCalling; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfComputer; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfApplyPatch; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[ResponsesDelegationUpdateConfigToolUnionParam](
		"type",
		apijson.Discriminator[FunctionToolParam]("function"),
		apijson.Discriminator[ResponsesDelegationUpdateConfigToolWebSearchParam]("web_search"),
		apijson.Discriminator[ResponsesDelegationUpdateConfigToolFileSearchParam]("file_search"),
		apijson.Discriminator[ResponsesDelegationUpdateConfigToolCodeInterpreterParam]("code_interpreter"),
		apijson.Discriminator[ResponsesDelegationUpdateConfigToolShellParam]("shell"),
		apijson.Discriminator[ResponsesDelegationUpdateConfigToolImageGenerationParam]("image_generation"),
		apijson.Discriminator[ResponsesDelegationUpdateConfigToolMcpParam]("mcp"),
		apijson.Discriminator[ResponsesDelegationUpdateConfigToolCustomParam]("custom"),
		apijson.Discriminator[ResponsesDelegationUpdateConfigToolNamespaceParam]("namespace"),
		apijson.Discriminator[ResponsesDelegationUpdateConfigToolToolSearchParam]("tool_search"),
		apijson.Discriminator[ResponsesDelegationUpdateConfigToolProgrammaticToolCallingParam]("programmatic_tool_calling"),
		apijson.Discriminator[ResponsesDelegationUpdateConfigToolComputerParam]("computer"),
		apijson.Discriminator[ResponsesDelegationUpdateConfigToolApplyPatchParam]("apply_patch"),
	)
}

func NewResponsesDelegationUpdateConfigToolWebSearchParam() ResponsesDelegationUpdateConfigToolWebSearchParam {
	return ResponsesDelegationUpdateConfigToolWebSearchParam{
		Type: "web_search",
	}
}

// A web search tool available to the Live session’s Responses backend.
//
// This struct has a constant value, construct it with
// [NewResponsesDelegationUpdateConfigToolWebSearchParam].
type ResponsesDelegationUpdateConfigToolWebSearchParam struct {
	// The tool type. Always `web_search`.
	Type constant.WebSearch `json:"type" default:"web_search"`
	paramObj
}

func (r ResponsesDelegationUpdateConfigToolWebSearchParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationUpdateConfigToolWebSearchParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationUpdateConfigToolWebSearchParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewResponsesDelegationUpdateConfigToolFileSearchParam() ResponsesDelegationUpdateConfigToolFileSearchParam {
	return ResponsesDelegationUpdateConfigToolFileSearchParam{
		Type: "file_search",
	}
}

// This struct has a constant value, construct it with
// [NewResponsesDelegationUpdateConfigToolFileSearchParam].
type ResponsesDelegationUpdateConfigToolFileSearchParam struct {
	Type constant.FileSearch `json:"type" default:"file_search"`
	paramObj
}

func (r ResponsesDelegationUpdateConfigToolFileSearchParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationUpdateConfigToolFileSearchParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationUpdateConfigToolFileSearchParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewResponsesDelegationUpdateConfigToolCodeInterpreterParam() ResponsesDelegationUpdateConfigToolCodeInterpreterParam {
	return ResponsesDelegationUpdateConfigToolCodeInterpreterParam{
		Type: "code_interpreter",
	}
}

// This struct has a constant value, construct it with
// [NewResponsesDelegationUpdateConfigToolCodeInterpreterParam].
type ResponsesDelegationUpdateConfigToolCodeInterpreterParam struct {
	Type constant.CodeInterpreter `json:"type" default:"code_interpreter"`
	paramObj
}

func (r ResponsesDelegationUpdateConfigToolCodeInterpreterParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationUpdateConfigToolCodeInterpreterParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationUpdateConfigToolCodeInterpreterParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A Responses shell tool. Use a hosted container or return local shell results
// with response.item.create. Domain secrets are not supported.
//
// The property Type is required.
type ResponsesDelegationUpdateConfigToolShellParam struct {
	Environment ResponsesDelegationUpdateConfigToolShellEnvironmentUnionParam `json:"environment,omitzero"`
	// This field can be elided, and will marshal its zero value as "shell".
	Type constant.Shell `json:"type" default:"shell"`
	paramObj
}

func (r ResponsesDelegationUpdateConfigToolShellParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationUpdateConfigToolShellParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationUpdateConfigToolShellParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ResponsesDelegationUpdateConfigToolShellEnvironmentUnionParam struct {
	OfContainerAuto      *ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoParam      `json:",omitzero,inline"`
	OfContainerReference *ResponsesDelegationUpdateConfigToolShellEnvironmentContainerReferenceParam `json:",omitzero,inline"`
	OfLocal              *ResponsesDelegationUpdateConfigToolShellEnvironmentLocalParam              `json:",omitzero,inline"`
	paramUnion
}

func (u ResponsesDelegationUpdateConfigToolShellEnvironmentUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfContainerAuto, u.OfContainerReference, u.OfLocal)
}
func (u *ResponsesDelegationUpdateConfigToolShellEnvironmentUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationUpdateConfigToolShellEnvironmentUnionParam) GetFileIDs() []string {
	if vt := u.OfContainerAuto; vt != nil {
		return vt.FileIDs
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationUpdateConfigToolShellEnvironmentUnionParam) GetMemoryLimit() *string {
	if vt := u.OfContainerAuto; vt != nil {
		return &vt.MemoryLimit
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationUpdateConfigToolShellEnvironmentUnionParam) GetNetworkPolicy() *ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoNetworkPolicyUnionParam {
	if vt := u.OfContainerAuto; vt != nil {
		return &vt.NetworkPolicy
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationUpdateConfigToolShellEnvironmentUnionParam) GetContainerID() *string {
	if vt := u.OfContainerReference; vt != nil {
		return &vt.ContainerID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationUpdateConfigToolShellEnvironmentUnionParam) GetType() *string {
	if vt := u.OfContainerAuto; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfContainerReference; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfLocal; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u ResponsesDelegationUpdateConfigToolShellEnvironmentUnionParam) GetSkills() (res responsesDelegationUpdateConfigToolShellEnvironmentUnionParamSkills) {
	if vt := u.OfContainerAuto; vt != nil {
		res.any = &vt.Skills
	} else if vt := u.OfLocal; vt != nil {
		res.any = &vt.Skills
	}
	return
}

// Can have the runtime types
// [_[]ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoSkillUnionParam],
// [_[]ResponsesDelegationUpdateConfigToolShellEnvironmentLocalSkillParam]
type responsesDelegationUpdateConfigToolShellEnvironmentUnionParamSkills struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *[]live.ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoSkillUnionParam:
//	case *[]live.ResponsesDelegationUpdateConfigToolShellEnvironmentLocalSkillParam:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u responsesDelegationUpdateConfigToolShellEnvironmentUnionParamSkills) AsAny() any {
	return u.any
}

func init() {
	apijson.RegisterUnion[ResponsesDelegationUpdateConfigToolShellEnvironmentUnionParam](
		"type",
		apijson.Discriminator[ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoParam]("container_auto"),
		apijson.Discriminator[ResponsesDelegationUpdateConfigToolShellEnvironmentContainerReferenceParam]("container_reference"),
		apijson.Discriminator[ResponsesDelegationUpdateConfigToolShellEnvironmentLocalParam]("local"),
	)
}

// The property Type is required.
type ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoParam struct {
	// An optional list of uploaded files to make available to your code.
	FileIDs []string `json:"file_ids,omitzero"`
	// The memory limit for the container.
	//
	// Any of "1g", "4g", "16g", "64g".
	MemoryLimit string `json:"memory_limit,omitzero"`
	// Network access policy for the container.
	NetworkPolicy ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoNetworkPolicyUnionParam `json:"network_policy,omitzero"`
	// An optional list of skills referenced by id or inline data.
	Skills []ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoSkillUnionParam `json:"skills,omitzero"`
	// Automatically creates a container for this request
	//
	// This field can be elided, and will marshal its zero value as "container_auto".
	Type constant.ContainerAuto `json:"type" default:"container_auto"`
	paramObj
}

func (r ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoParam](
		"memory_limit", "1g", "4g", "16g", "64g",
	)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoNetworkPolicyUnionParam struct {
	OfDisabled  *ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoNetworkPolicyDisabledParam  `json:",omitzero,inline"`
	OfAllowlist *ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoNetworkPolicyAllowlistParam `json:",omitzero,inline"`
	paramUnion
}

func (u ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoNetworkPolicyUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfDisabled, u.OfAllowlist)
}
func (u *ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoNetworkPolicyUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoNetworkPolicyUnionParam) GetAllowedDomains() []string {
	if vt := u.OfAllowlist; vt != nil {
		return vt.AllowedDomains
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoNetworkPolicyUnionParam) GetType() *string {
	if vt := u.OfDisabled; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfAllowlist; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoNetworkPolicyUnionParam](
		"type",
		apijson.Discriminator[ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoNetworkPolicyDisabledParam]("disabled"),
		apijson.Discriminator[ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoNetworkPolicyAllowlistParam]("allowlist"),
	)
}

func NewResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoNetworkPolicyDisabledParam() ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoNetworkPolicyDisabledParam {
	return ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoNetworkPolicyDisabledParam{
		Type: "disabled",
	}
}

// This struct has a constant value, construct it with
// [NewResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoNetworkPolicyDisabledParam].
type ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoNetworkPolicyDisabledParam struct {
	// Disable outbound network access. Always `disabled`.
	Type constant.Disabled `json:"type" default:"disabled"`
	paramObj
}

func (r ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoNetworkPolicyDisabledParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoNetworkPolicyDisabledParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoNetworkPolicyDisabledParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties AllowedDomains, Type are required.
type ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoNetworkPolicyAllowlistParam struct {
	// A list of allowed domains when type is `allowlist`.
	AllowedDomains []string `json:"allowed_domains,omitzero" api:"required"`
	// Allow outbound network access only to specified domains. Always `allowlist`.
	//
	// This field can be elided, and will marshal its zero value as "allowlist".
	Type constant.Allowlist `json:"type" default:"allowlist"`
	paramObj
}

func (r ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoNetworkPolicyAllowlistParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoNetworkPolicyAllowlistParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoNetworkPolicyAllowlistParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoSkillUnionParam struct {
	OfSkillReference *ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoSkillSkillReferenceParam `json:",omitzero,inline"`
	OfInline         *ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoSkillInlineParam         `json:",omitzero,inline"`
	paramUnion
}

func (u ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoSkillUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfSkillReference, u.OfInline)
}
func (u *ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoSkillUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoSkillUnionParam) GetSkillID() *string {
	if vt := u.OfSkillReference; vt != nil {
		return &vt.SkillID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoSkillUnionParam) GetVersion() *string {
	if vt := u.OfSkillReference; vt != nil && vt.Version.Valid() {
		return &vt.Version.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoSkillUnionParam) GetDescription() *string {
	if vt := u.OfInline; vt != nil {
		return &vt.Description
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoSkillUnionParam) GetName() *string {
	if vt := u.OfInline; vt != nil {
		return &vt.Name
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoSkillUnionParam) GetSource() *ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoSkillInlineSourceParam {
	if vt := u.OfInline; vt != nil {
		return &vt.Source
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoSkillUnionParam) GetType() *string {
	if vt := u.OfSkillReference; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfInline; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoSkillUnionParam](
		"type",
		apijson.Discriminator[ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoSkillSkillReferenceParam]("skill_reference"),
		apijson.Discriminator[ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoSkillInlineParam]("inline"),
	)
}

// The properties SkillID, Type are required.
type ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoSkillSkillReferenceParam struct {
	// The ID of the referenced skill.
	SkillID string `json:"skill_id" api:"required"`
	// Optional skill version. Use a positive integer or 'latest'. Omit for default.
	Version param.Opt[string] `json:"version,omitzero"`
	// References a skill created with the /v1/skills endpoint.
	//
	// This field can be elided, and will marshal its zero value as "skill_reference".
	Type constant.SkillReference `json:"type" default:"skill_reference"`
	paramObj
}

func (r ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoSkillSkillReferenceParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoSkillSkillReferenceParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoSkillSkillReferenceParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Description, Name, Source, Type are required.
type ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoSkillInlineParam struct {
	// The description of the skill.
	Description string `json:"description" api:"required"`
	// The name of the skill.
	Name string `json:"name" api:"required"`
	// Inline skill payload
	Source ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoSkillInlineSourceParam `json:"source" api:"required"`
	// Defines an inline skill for this request.
	//
	// This field can be elided, and will marshal its zero value as "inline".
	Type constant.Inline `json:"type" default:"inline"`
	paramObj
}

func (r ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoSkillInlineParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoSkillInlineParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoSkillInlineParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Inline skill payload
//
// The properties Data, MediaType, Type are required.
type ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoSkillInlineSourceParam struct {
	// Base64-encoded skill zip bundle.
	Data string `json:"data" api:"required"`
	// The media type of the inline skill payload. Must be `application/zip`.
	//
	// This field can be elided, and will marshal its zero value as "application/zip".
	MediaType constant.ApplicationZip `json:"media_type" default:"application/zip"`
	// The type of the inline skill source. Must be `base64`.
	//
	// This field can be elided, and will marshal its zero value as "base64".
	Type constant.Base64 `json:"type" default:"base64"`
	paramObj
}

func (r ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoSkillInlineSourceParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoSkillInlineSourceParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoSkillInlineSourceParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties ContainerID, Type are required.
type ResponsesDelegationUpdateConfigToolShellEnvironmentContainerReferenceParam struct {
	// The ID of the referenced container.
	ContainerID string `json:"container_id" api:"required"`
	// References a container created with the /v1/containers endpoint
	//
	// This field can be elided, and will marshal its zero value as
	// "container_reference".
	Type constant.ContainerReference `json:"type" default:"container_reference"`
	paramObj
}

func (r ResponsesDelegationUpdateConfigToolShellEnvironmentContainerReferenceParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationUpdateConfigToolShellEnvironmentContainerReferenceParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationUpdateConfigToolShellEnvironmentContainerReferenceParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The property Type is required.
type ResponsesDelegationUpdateConfigToolShellEnvironmentLocalParam struct {
	// An optional list of skills.
	Skills []ResponsesDelegationUpdateConfigToolShellEnvironmentLocalSkillParam `json:"skills,omitzero"`
	// Use a local computer environment.
	//
	// This field can be elided, and will marshal its zero value as "local".
	Type constant.Local `json:"type" default:"local"`
	paramObj
}

func (r ResponsesDelegationUpdateConfigToolShellEnvironmentLocalParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationUpdateConfigToolShellEnvironmentLocalParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationUpdateConfigToolShellEnvironmentLocalParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Description, Name, Path are required.
type ResponsesDelegationUpdateConfigToolShellEnvironmentLocalSkillParam struct {
	// The description of the skill.
	Description string `json:"description" api:"required"`
	// The name of the skill.
	Name string `json:"name" api:"required"`
	// The path to the directory containing the skill.
	Path string `json:"path" api:"required"`
	paramObj
}

func (r ResponsesDelegationUpdateConfigToolShellEnvironmentLocalSkillParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationUpdateConfigToolShellEnvironmentLocalSkillParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationUpdateConfigToolShellEnvironmentLocalSkillParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewResponsesDelegationUpdateConfigToolImageGenerationParam() ResponsesDelegationUpdateConfigToolImageGenerationParam {
	return ResponsesDelegationUpdateConfigToolImageGenerationParam{
		Type: "image_generation",
	}
}

// This struct has a constant value, construct it with
// [NewResponsesDelegationUpdateConfigToolImageGenerationParam].
type ResponsesDelegationUpdateConfigToolImageGenerationParam struct {
	Type constant.ImageGeneration `json:"type" default:"image_generation"`
	paramObj
}

func (r ResponsesDelegationUpdateConfigToolImageGenerationParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationUpdateConfigToolImageGenerationParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationUpdateConfigToolImageGenerationParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewResponsesDelegationUpdateConfigToolMcpParam() ResponsesDelegationUpdateConfigToolMcpParam {
	return ResponsesDelegationUpdateConfigToolMcpParam{
		Type: "mcp",
	}
}

// This struct has a constant value, construct it with
// [NewResponsesDelegationUpdateConfigToolMcpParam].
type ResponsesDelegationUpdateConfigToolMcpParam struct {
	Type constant.Mcp `json:"type" default:"mcp"`
	paramObj
}

func (r ResponsesDelegationUpdateConfigToolMcpParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationUpdateConfigToolMcpParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationUpdateConfigToolMcpParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewResponsesDelegationUpdateConfigToolCustomParam() ResponsesDelegationUpdateConfigToolCustomParam {
	return ResponsesDelegationUpdateConfigToolCustomParam{
		Type: "custom",
	}
}

// This struct has a constant value, construct it with
// [NewResponsesDelegationUpdateConfigToolCustomParam].
type ResponsesDelegationUpdateConfigToolCustomParam struct {
	Type constant.Custom `json:"type" default:"custom"`
	paramObj
}

func (r ResponsesDelegationUpdateConfigToolCustomParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationUpdateConfigToolCustomParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationUpdateConfigToolCustomParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewResponsesDelegationUpdateConfigToolNamespaceParam() ResponsesDelegationUpdateConfigToolNamespaceParam {
	return ResponsesDelegationUpdateConfigToolNamespaceParam{
		Type: "namespace",
	}
}

// This struct has a constant value, construct it with
// [NewResponsesDelegationUpdateConfigToolNamespaceParam].
type ResponsesDelegationUpdateConfigToolNamespaceParam struct {
	Type constant.Namespace `json:"type" default:"namespace"`
	paramObj
}

func (r ResponsesDelegationUpdateConfigToolNamespaceParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationUpdateConfigToolNamespaceParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationUpdateConfigToolNamespaceParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewResponsesDelegationUpdateConfigToolToolSearchParam() ResponsesDelegationUpdateConfigToolToolSearchParam {
	return ResponsesDelegationUpdateConfigToolToolSearchParam{
		Type: "tool_search",
	}
}

// This struct has a constant value, construct it with
// [NewResponsesDelegationUpdateConfigToolToolSearchParam].
type ResponsesDelegationUpdateConfigToolToolSearchParam struct {
	Type constant.ToolSearch `json:"type" default:"tool_search"`
	paramObj
}

func (r ResponsesDelegationUpdateConfigToolToolSearchParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationUpdateConfigToolToolSearchParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationUpdateConfigToolToolSearchParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewResponsesDelegationUpdateConfigToolProgrammaticToolCallingParam() ResponsesDelegationUpdateConfigToolProgrammaticToolCallingParam {
	return ResponsesDelegationUpdateConfigToolProgrammaticToolCallingParam{
		Type: "programmatic_tool_calling",
	}
}

// This struct has a constant value, construct it with
// [NewResponsesDelegationUpdateConfigToolProgrammaticToolCallingParam].
type ResponsesDelegationUpdateConfigToolProgrammaticToolCallingParam struct {
	Type constant.ProgrammaticToolCalling `json:"type" default:"programmatic_tool_calling"`
	paramObj
}

func (r ResponsesDelegationUpdateConfigToolProgrammaticToolCallingParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationUpdateConfigToolProgrammaticToolCallingParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationUpdateConfigToolProgrammaticToolCallingParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewResponsesDelegationUpdateConfigToolComputerParam() ResponsesDelegationUpdateConfigToolComputerParam {
	return ResponsesDelegationUpdateConfigToolComputerParam{
		Type: "computer",
	}
}

// This struct has a constant value, construct it with
// [NewResponsesDelegationUpdateConfigToolComputerParam].
type ResponsesDelegationUpdateConfigToolComputerParam struct {
	Type constant.Computer `json:"type" default:"computer"`
	paramObj
}

func (r ResponsesDelegationUpdateConfigToolComputerParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationUpdateConfigToolComputerParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationUpdateConfigToolComputerParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewResponsesDelegationUpdateConfigToolApplyPatchParam() ResponsesDelegationUpdateConfigToolApplyPatchParam {
	return ResponsesDelegationUpdateConfigToolApplyPatchParam{
		Type: "apply_patch",
	}
}

// This struct has a constant value, construct it with
// [NewResponsesDelegationUpdateConfigToolApplyPatchParam].
type ResponsesDelegationUpdateConfigToolApplyPatchParam struct {
	Type constant.ApplyPatch `json:"type" default:"apply_patch"`
	paramObj
}

func (r ResponsesDelegationUpdateConfigToolApplyPatchParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationUpdateConfigToolApplyPatchParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationUpdateConfigToolApplyPatchParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ServerEventUnion contains all possible properties and values from
// [SessionStartedEvent], [SessionUpdatedEvent], [InputAudioMutedEvent],
// [InputAudioUnmutedEvent], [OutputAudioFixturePlayingEvent],
// [OutputAudioFixturePlayedEvent], [InstructionsAppendedEvent],
// [ThinkingAppendedEvent], [CommentaryAppendedEvent],
// [ServerEventSessionInputAudioAppend], [OutputAudioDeltaEvent],
// [InputTranscriptDeltaEvent], [OutputTranscriptDeltaEvent],
// [DelegationCreatedEvent], [DecisionEvent], [DecisionItemCreatedEvent],
// [DecisionItemDeletedEvent], [ResponseEvent], [SessionUsageUpdatedEvent],
// [SessionClosedEvent], [ErrorEvent], [InfoEvent],
// [ServerEventTransportDtmfReceived], [ServerEventTransportDtmfSend],
// [ServerEventTransportRinging], [ServerEventTransportAnswered],
// [ServerEventTransportFailed].
//
// Use the [ServerEventUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ServerEventUnion struct {
	OfSessionStartedEvent                SessionStartedEvent                `json:",inline"`
	OfSessionUpdatedEvent                SessionUpdatedEvent                `json:",inline"`
	OfInputAudioMutedEvent               InputAudioMutedEvent               `json:",inline"`
	OfInputAudioUnmutedEvent             InputAudioUnmutedEvent             `json:",inline"`
	OfOutputAudioFixturePlayingEvent     OutputAudioFixturePlayingEvent     `json:",inline"`
	OfOutputAudioFixturePlayedEvent      OutputAudioFixturePlayedEvent      `json:",inline"`
	OfInstructionsAppendedEvent          InstructionsAppendedEvent          `json:",inline"`
	OfThinkingAppendedEvent              ThinkingAppendedEvent              `json:",inline"`
	OfCommentaryAppendedEvent            CommentaryAppendedEvent            `json:",inline"`
	OfServerEventSessionInputAudioAppend ServerEventSessionInputAudioAppend `json:",inline"`
	OfOutputAudioDeltaEvent              OutputAudioDeltaEvent              `json:",inline"`
	OfInputTranscriptDeltaEvent          InputTranscriptDeltaEvent          `json:",inline"`
	OfOutputTranscriptDeltaEvent         OutputTranscriptDeltaEvent         `json:",inline"`
	OfDelegationCreatedEvent             DelegationCreatedEvent             `json:",inline"`
	OfDecisionEvent                      DecisionEvent                      `json:",inline"`
	OfDecisionItemCreatedEvent           DecisionItemCreatedEvent           `json:",inline"`
	OfDecisionItemDeletedEvent           DecisionItemDeletedEvent           `json:",inline"`
	OfResponseEvent                      ResponseEvent                      `json:",inline"`
	OfSessionUsageUpdatedEvent           SessionUsageUpdatedEvent           `json:",inline"`
	OfSessionClosedEvent                 SessionClosedEvent                 `json:",inline"`
	OfErrorEvent                         ErrorEvent                         `json:",inline"`
	OfInfoEvent                          InfoEvent                          `json:",inline"`
	OfServerEventTransportDtmfReceived   ServerEventTransportDtmfReceived   `json:",inline"`
	OfServerEventTransportDtmfSend       ServerEventTransportDtmfSend       `json:",inline"`
	OfServerEventTransportRinging        ServerEventTransportRinging        `json:",inline"`
	OfServerEventTransportAnswered       ServerEventTransportAnswered       `json:",inline"`
	OfServerEventTransportFailed         ServerEventTransportFailed         `json:",inline"`
	Type                                 string                             `json:"type"`
	JSON                                 struct {
		OfSessionStartedEvent                respjson.Field
		OfSessionUpdatedEvent                respjson.Field
		OfInputAudioMutedEvent               respjson.Field
		OfInputAudioUnmutedEvent             respjson.Field
		OfOutputAudioFixturePlayingEvent     respjson.Field
		OfOutputAudioFixturePlayedEvent      respjson.Field
		OfInstructionsAppendedEvent          respjson.Field
		OfThinkingAppendedEvent              respjson.Field
		OfCommentaryAppendedEvent            respjson.Field
		OfServerEventSessionInputAudioAppend respjson.Field
		OfOutputAudioDeltaEvent              respjson.Field
		OfInputTranscriptDeltaEvent          respjson.Field
		OfOutputTranscriptDeltaEvent         respjson.Field
		OfDelegationCreatedEvent             respjson.Field
		OfDecisionEvent                      respjson.Field
		OfDecisionItemCreatedEvent           respjson.Field
		OfDecisionItemDeletedEvent           respjson.Field
		OfResponseEvent                      respjson.Field
		OfSessionUsageUpdatedEvent           respjson.Field
		OfSessionClosedEvent                 respjson.Field
		OfErrorEvent                         respjson.Field
		OfInfoEvent                          respjson.Field
		OfServerEventTransportDtmfReceived   respjson.Field
		OfServerEventTransportDtmfSend       respjson.Field
		OfServerEventTransportRinging        respjson.Field
		OfServerEventTransportAnswered       respjson.Field
		OfServerEventTransportFailed         respjson.Field
		Type                                 respjson.Field
		raw                                  string
	} `json:"-"`
}

// anyServerEvent is implemented by each variant of [ServerEventUnion] to add type
// safety for the return type of [ServerEventUnion.AsAny]
type anyServerEvent interface {
	implServerEventUnion()
}

func (SessionStartedEvent) implServerEventUnion()                {}
func (SessionUpdatedEvent) implServerEventUnion()                {}
func (InputAudioMutedEvent) implServerEventUnion()               {}
func (InputAudioUnmutedEvent) implServerEventUnion()             {}
func (OutputAudioFixturePlayingEvent) implServerEventUnion()     {}
func (OutputAudioFixturePlayedEvent) implServerEventUnion()      {}
func (InstructionsAppendedEvent) implServerEventUnion()          {}
func (ThinkingAppendedEvent) implServerEventUnion()              {}
func (CommentaryAppendedEvent) implServerEventUnion()            {}
func (ServerEventSessionInputAudioAppend) implServerEventUnion() {}
func (OutputAudioDeltaEvent) implServerEventUnion()              {}
func (InputTranscriptDeltaEvent) implServerEventUnion()          {}
func (OutputTranscriptDeltaEvent) implServerEventUnion()         {}
func (DelegationCreatedEvent) implServerEventUnion()             {}
func (DecisionEvent) implServerEventUnion()                      {}
func (DecisionItemCreatedEvent) implServerEventUnion()           {}
func (DecisionItemDeletedEvent) implServerEventUnion()           {}
func (ResponseEvent) implServerEventUnion()                      {}
func (SessionUsageUpdatedEvent) implServerEventUnion()           {}
func (SessionClosedEvent) implServerEventUnion()                 {}
func (ErrorEvent) implServerEventUnion()                         {}
func (InfoEvent) implServerEventUnion()                          {}
func (ServerEventTransportDtmfReceived) implServerEventUnion()   {}
func (ServerEventTransportDtmfSend) implServerEventUnion()       {}
func (ServerEventTransportRinging) implServerEventUnion()        {}
func (ServerEventTransportAnswered) implServerEventUnion()       {}
func (ServerEventTransportFailed) implServerEventUnion()         {}

// Use the following switch statement to find the correct variant
//
//	switch variant := ServerEventUnion.AsAny().(type) {
//	case live.SessionStartedEvent:
//	case live.SessionUpdatedEvent:
//	case live.InputAudioMutedEvent:
//	case live.InputAudioUnmutedEvent:
//	case live.OutputAudioFixturePlayingEvent:
//	case live.OutputAudioFixturePlayedEvent:
//	case live.InstructionsAppendedEvent:
//	case live.ThinkingAppendedEvent:
//	case live.CommentaryAppendedEvent:
//	case live.ServerEventSessionInputAudioAppend:
//	case live.OutputAudioDeltaEvent:
//	case live.InputTranscriptDeltaEvent:
//	case live.OutputTranscriptDeltaEvent:
//	case live.DelegationCreatedEvent:
//	case live.DecisionEvent:
//	case live.DecisionItemCreatedEvent:
//	case live.DecisionItemDeletedEvent:
//	case live.ResponseEvent:
//	case live.SessionUsageUpdatedEvent:
//	case live.SessionClosedEvent:
//	case live.ErrorEvent:
//	case live.InfoEvent:
//	case live.ServerEventTransportDtmfReceived:
//	case live.ServerEventTransportDtmfSend:
//	case live.ServerEventTransportRinging:
//	case live.ServerEventTransportAnswered:
//	case live.ServerEventTransportFailed:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ServerEventUnion) AsAny() anyServerEvent {
	switch u.Type {
	case "session.started":
		return u.AsSessionStartedEvent()
	case "session.updated":
		return u.AsSessionUpdatedEvent()
	case "session.input_audio.muted":
		return u.AsInputAudioMutedEvent()
	case "session.input_audio.unmuted":
		return u.AsInputAudioUnmutedEvent()
	case "session.output_audio.fixture.playing":
		return u.AsOutputAudioFixturePlayingEvent()
	case "session.output_audio.fixture.played":
		return u.AsOutputAudioFixturePlayedEvent()
	case "session.instructions.appended":
		return u.AsInstructionsAppendedEvent()
	case "session.thinking.appended":
		return u.AsThinkingAppendedEvent()
	case "session.commentary.appended":
		return u.AsCommentaryAppendedEvent()
	case "session.input_audio.append":
		return u.AsServerEventSessionInputAudioAppend()
	case "session.output_audio.delta":
		return u.AsOutputAudioDeltaEvent()
	case "session.input_transcript.delta":
		return u.AsInputTranscriptDeltaEvent()
	case "session.output_transcript.delta":
		return u.AsOutputTranscriptDeltaEvent()
	case "session.delegation.created":
		return u.AsDelegationCreatedEvent()
	case "decision.event":
		return u.AsDecisionEvent()
	case "decision.item.created":
		return u.AsDecisionItemCreatedEvent()
	case "decision.item.deleted":
		return u.AsDecisionItemDeletedEvent()
	case "response.event":
		return u.AsResponseEvent()
	case "session.usage.updated":
		return u.AsSessionUsageUpdatedEvent()
	case "session.closed":
		return u.AsSessionClosedEvent()
	case "error":
		return u.AsErrorEvent()
	case "info":
		return u.AsInfoEvent()
	case "transport.dtmf.received":
		return u.AsServerEventTransportDtmfReceived()
	case "transport.dtmf.send":
		return u.AsServerEventTransportDtmfSend()
	case "transport.ringing":
		return u.AsServerEventTransportRinging()
	case "transport.answered":
		return u.AsServerEventTransportAnswered()
	case "transport.failed":
		return u.AsServerEventTransportFailed()
	}
	return nil
}

func (u ServerEventUnion) AsSessionStartedEvent() (v SessionStartedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ServerEventUnion) AsSessionUpdatedEvent() (v SessionUpdatedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ServerEventUnion) AsInputAudioMutedEvent() (v InputAudioMutedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ServerEventUnion) AsInputAudioUnmutedEvent() (v InputAudioUnmutedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ServerEventUnion) AsOutputAudioFixturePlayingEvent() (v OutputAudioFixturePlayingEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ServerEventUnion) AsOutputAudioFixturePlayedEvent() (v OutputAudioFixturePlayedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ServerEventUnion) AsInstructionsAppendedEvent() (v InstructionsAppendedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ServerEventUnion) AsThinkingAppendedEvent() (v ThinkingAppendedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ServerEventUnion) AsCommentaryAppendedEvent() (v CommentaryAppendedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ServerEventUnion) AsServerEventSessionInputAudioAppend() (v ServerEventSessionInputAudioAppend) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ServerEventUnion) AsOutputAudioDeltaEvent() (v OutputAudioDeltaEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ServerEventUnion) AsInputTranscriptDeltaEvent() (v InputTranscriptDeltaEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ServerEventUnion) AsOutputTranscriptDeltaEvent() (v OutputTranscriptDeltaEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ServerEventUnion) AsDelegationCreatedEvent() (v DelegationCreatedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ServerEventUnion) AsDecisionEvent() (v DecisionEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ServerEventUnion) AsDecisionItemCreatedEvent() (v DecisionItemCreatedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ServerEventUnion) AsDecisionItemDeletedEvent() (v DecisionItemDeletedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ServerEventUnion) AsResponseEvent() (v ResponseEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ServerEventUnion) AsSessionUsageUpdatedEvent() (v SessionUsageUpdatedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ServerEventUnion) AsSessionClosedEvent() (v SessionClosedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ServerEventUnion) AsErrorEvent() (v ErrorEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ServerEventUnion) AsInfoEvent() (v InfoEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ServerEventUnion) AsServerEventTransportDtmfReceived() (v ServerEventTransportDtmfReceived) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ServerEventUnion) AsServerEventTransportDtmfSend() (v ServerEventTransportDtmfSend) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ServerEventUnion) AsServerEventTransportRinging() (v ServerEventTransportRinging) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ServerEventUnion) AsServerEventTransportAnswered() (v ServerEventTransportAnswered) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ServerEventUnion) AsServerEventTransportFailed() (v ServerEventTransportFailed) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ServerEventUnion) RawJSON() string { return u.JSON.raw }

func (r *ServerEventUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// unmarshalWebsocketEventWithRawType decodes only the event selected by its type. Public
// UnmarshalJSON keeps the existing general-purpose union decoding behavior.
func (r *ServerEventUnion) unmarshalWebsocketEventWithRawType(data []byte) error {
	var envelope struct {
		Type json.RawMessage `json:"type"`
	}
	if err := json.Unmarshal(data, &envelope); err != nil {
		return err
	}
	var eventType string
	if err := json.Unmarshal(envelope.Type, &eventType); err != nil {
		return err
	}
	*r = ServerEventUnion{Type: eventType}
	r.JSON.raw = string(data)
	r.JSON.Type = respjson.NewField(string(envelope.Type))
	switch eventType {
	case "session.started":
		if err := apijson.UnmarshalRoot(data, &r.OfSessionStartedEvent); err != nil {
			return err
		}
		r.JSON.OfSessionStartedEvent = respjson.NewField(r.JSON.raw)
	case "session.updated":
		if err := apijson.UnmarshalRoot(data, &r.OfSessionUpdatedEvent); err != nil {
			return err
		}
		r.JSON.OfSessionUpdatedEvent = respjson.NewField(r.JSON.raw)
	case "session.input_audio.muted":
		if err := apijson.UnmarshalRoot(data, &r.OfInputAudioMutedEvent); err != nil {
			return err
		}
		r.JSON.OfInputAudioMutedEvent = respjson.NewField(r.JSON.raw)
	case "session.input_audio.unmuted":
		if err := apijson.UnmarshalRoot(data, &r.OfInputAudioUnmutedEvent); err != nil {
			return err
		}
		r.JSON.OfInputAudioUnmutedEvent = respjson.NewField(r.JSON.raw)
	case "session.output_audio.fixture.playing":
		if err := apijson.UnmarshalRoot(data, &r.OfOutputAudioFixturePlayingEvent); err != nil {
			return err
		}
		r.JSON.OfOutputAudioFixturePlayingEvent = respjson.NewField(r.JSON.raw)
	case "session.output_audio.fixture.played":
		if err := apijson.UnmarshalRoot(data, &r.OfOutputAudioFixturePlayedEvent); err != nil {
			return err
		}
		r.JSON.OfOutputAudioFixturePlayedEvent = respjson.NewField(r.JSON.raw)
	case "session.instructions.appended":
		if err := apijson.UnmarshalRoot(data, &r.OfInstructionsAppendedEvent); err != nil {
			return err
		}
		r.JSON.OfInstructionsAppendedEvent = respjson.NewField(r.JSON.raw)
	case "session.thinking.appended":
		if err := apijson.UnmarshalRoot(data, &r.OfThinkingAppendedEvent); err != nil {
			return err
		}
		r.JSON.OfThinkingAppendedEvent = respjson.NewField(r.JSON.raw)
	case "session.commentary.appended":
		if err := apijson.UnmarshalRoot(data, &r.OfCommentaryAppendedEvent); err != nil {
			return err
		}
		r.JSON.OfCommentaryAppendedEvent = respjson.NewField(r.JSON.raw)
	case "session.input_audio.append":
		if err := apijson.UnmarshalRoot(data, &r.OfServerEventSessionInputAudioAppend); err != nil {
			return err
		}
		r.JSON.OfServerEventSessionInputAudioAppend = respjson.NewField(r.JSON.raw)
	case "session.output_audio.delta":
		if err := apijson.UnmarshalRoot(data, &r.OfOutputAudioDeltaEvent); err != nil {
			return err
		}
		r.JSON.OfOutputAudioDeltaEvent = respjson.NewField(r.JSON.raw)
	case "session.input_transcript.delta":
		if err := apijson.UnmarshalRoot(data, &r.OfInputTranscriptDeltaEvent); err != nil {
			return err
		}
		r.JSON.OfInputTranscriptDeltaEvent = respjson.NewField(r.JSON.raw)
	case "session.output_transcript.delta":
		if err := apijson.UnmarshalRoot(data, &r.OfOutputTranscriptDeltaEvent); err != nil {
			return err
		}
		r.JSON.OfOutputTranscriptDeltaEvent = respjson.NewField(r.JSON.raw)
	case "session.delegation.created":
		if err := apijson.UnmarshalRoot(data, &r.OfDelegationCreatedEvent); err != nil {
			return err
		}
		r.JSON.OfDelegationCreatedEvent = respjson.NewField(r.JSON.raw)
	case "decision.event":
		if err := apijson.UnmarshalRoot(data, &r.OfDecisionEvent); err != nil {
			return err
		}
		r.JSON.OfDecisionEvent = respjson.NewField(r.JSON.raw)
	case "decision.item.created":
		if err := apijson.UnmarshalRoot(data, &r.OfDecisionItemCreatedEvent); err != nil {
			return err
		}
		r.JSON.OfDecisionItemCreatedEvent = respjson.NewField(r.JSON.raw)
	case "decision.item.deleted":
		if err := apijson.UnmarshalRoot(data, &r.OfDecisionItemDeletedEvent); err != nil {
			return err
		}
		r.JSON.OfDecisionItemDeletedEvent = respjson.NewField(r.JSON.raw)
	case "response.event":
		if err := apijson.UnmarshalRoot(data, &r.OfResponseEvent); err != nil {
			return err
		}
		r.JSON.OfResponseEvent = respjson.NewField(r.JSON.raw)
	case "session.usage.updated":
		if err := apijson.UnmarshalRoot(data, &r.OfSessionUsageUpdatedEvent); err != nil {
			return err
		}
		r.JSON.OfSessionUsageUpdatedEvent = respjson.NewField(r.JSON.raw)
	case "session.closed":
		if err := apijson.UnmarshalRoot(data, &r.OfSessionClosedEvent); err != nil {
			return err
		}
		r.JSON.OfSessionClosedEvent = respjson.NewField(r.JSON.raw)
	case "error":
		if err := apijson.UnmarshalRoot(data, &r.OfErrorEvent); err != nil {
			return err
		}
		r.JSON.OfErrorEvent = respjson.NewField(r.JSON.raw)
	case "info":
		if err := apijson.UnmarshalRoot(data, &r.OfInfoEvent); err != nil {
			return err
		}
		r.JSON.OfInfoEvent = respjson.NewField(r.JSON.raw)
	case "transport.dtmf.received":
		if err := apijson.UnmarshalRoot(data, &r.OfServerEventTransportDtmfReceived); err != nil {
			return err
		}
		r.JSON.OfServerEventTransportDtmfReceived = respjson.NewField(r.JSON.raw)
	case "transport.dtmf.send":
		if err := apijson.UnmarshalRoot(data, &r.OfServerEventTransportDtmfSend); err != nil {
			return err
		}
		r.JSON.OfServerEventTransportDtmfSend = respjson.NewField(r.JSON.raw)
	case "transport.ringing":
		if err := apijson.UnmarshalRoot(data, &r.OfServerEventTransportRinging); err != nil {
			return err
		}
		r.JSON.OfServerEventTransportRinging = respjson.NewField(r.JSON.raw)
	case "transport.answered":
		if err := apijson.UnmarshalRoot(data, &r.OfServerEventTransportAnswered); err != nil {
			return err
		}
		r.JSON.OfServerEventTransportAnswered = respjson.NewField(r.JSON.raw)
	case "transport.failed":
		if err := apijson.UnmarshalRoot(data, &r.OfServerEventTransportFailed); err != nil {
			return err
		}
		r.JSON.OfServerEventTransportFailed = respjson.NewField(r.JSON.raw)
	}
	return nil
}

// Input audio received from the primary transport and reflected to a Live sideband
// connection before model-input muting.
type ServerEventSessionInputAudioAppend struct {
	// Base64-encoded raw mono PCM16LE at 24 kHz received from the primary transport,
	// reflected to the sideband before model-input muting. This server event uses the
	// same audio key as the client command, but is not an acknowledgment of it.
	Audio string `json:"audio" api:"required"`
	// The event type, always `session.input_audio.append`.
	Type constant.SessionInputAudioAppend `json:"type" default:"session.input_audio.append"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Audio       respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ServerEventSessionInputAudioAppend) RawJSON() string { return r.JSON.raw }
func (r *ServerEventSessionInputAudioAppend) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A SIP DTMF keypress received from the caller. Delivered only to sideband
// observers.
type ServerEventTransportDtmfReceived struct {
	Event   string                         `json:"event" api:"required"`
	EventID string                         `json:"event_id" api:"required"`
	Type    constant.TransportDtmfReceived `json:"type" default:"transport.dtmf.received"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Event       respjson.Field
		EventID     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ServerEventTransportDtmfReceived) RawJSON() string { return r.JSON.raw }
func (r *ServerEventTransportDtmfReceived) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A SIP DTMF keypress successfully sent to the SIP trunk. Delivered only to
// sideband observers; this is not a client command.
type ServerEventTransportDtmfSend struct {
	// The event_id of the client command, when supplied.
	ClientEventID string                     `json:"client_event_id"`
	Event         string                     `json:"event" api:"required"`
	EventID       string                     `json:"event_id" api:"required"`
	Type          constant.TransportDtmfSend `json:"type" default:"transport.dtmf.send"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ClientEventID respjson.Field
		Event         respjson.Field
		EventID       respjson.Field
		Type          respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ServerEventTransportDtmfSend) RawJSON() string { return r.JSON.raw }
func (r *ServerEventTransportDtmfSend) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The outbound SIP provider leg is ringing or providing early media. Delivered
// only to sideband observers.
type ServerEventTransportRinging struct {
	EventID string `json:"event_id" api:"required"`
	// The canonical Live session ID.
	SessionID string                    `json:"session_id" api:"required"`
	Type      constant.TransportRinging `json:"type" default:"transport.ringing"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventID     respjson.Field
		SessionID   respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ServerEventTransportRinging) RawJSON() string { return r.JSON.raw }
func (r *ServerEventTransportRinging) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The outbound SIP provider leg answered and media is established. Delivered only
// to sideband observers.
type ServerEventTransportAnswered struct {
	EventID string `json:"event_id" api:"required"`
	// The canonical Live session ID.
	SessionID string                     `json:"session_id" api:"required"`
	Type      constant.TransportAnswered `json:"type" default:"transport.answered"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventID     respjson.Field
		SessionID   respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ServerEventTransportAnswered) RawJSON() string { return r.JSON.raw }
func (r *ServerEventTransportAnswered) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// An asynchronous outbound SIP setup failure. Delivered only to sideband
// observers.
type ServerEventTransportFailed struct {
	Error   ServerEventTransportFailedError `json:"error" api:"required"`
	EventID string                          `json:"event_id" api:"required"`
	// The canonical Live session ID.
	SessionID string                   `json:"session_id" api:"required"`
	Type      constant.TransportFailed `json:"type" default:"transport.failed"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Error       respjson.Field
		EventID     respjson.Field
		SessionID   respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ServerEventTransportFailed) RawJSON() string { return r.JSON.raw }
func (r *ServerEventTransportFailed) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ServerEventTransportFailedError struct {
	// The call setup failure code.
	Code    string             `json:"code" api:"required"`
	Message string             `json:"message" api:"required"`
	Type    constant.CallError `json:"type" default:"call_error"`
	// The parameter related to the error, if any. Empty when no parameter applies.
	Param string `json:"param"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Code        respjson.Field
		Message     respjson.Field
		Type        respjson.Field
		Param       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ServerEventTransportFailedError) RawJSON() string { return r.JSON.raw }
func (r *ServerEventTransportFailedError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A Live server event selector for the WebRTC frontend data channel.
type ServerEventSelector struct {
	// The outer Live server event type. Use 'response.event' for Responses events.
	Type string `json:"type" api:"required"`
	// The nested Responses event type. Required when type is 'response.event';
	// forbidden for other event types.
	ResponseEvent string `json:"response_event"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type          respjson.Field
		ResponseEvent respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ServerEventSelector) RawJSON() string { return r.JSON.raw }
func (r *ServerEventSelector) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this ServerEventSelector to a ServerEventSelectorParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// ServerEventSelectorParam.Overrides()
func (r ServerEventSelector) ToParam() ServerEventSelectorParam {
	return param.Override[ServerEventSelectorParam](json.RawMessage(r.RawJSON()))
}

// A Live server event selector for the WebRTC frontend data channel.
//
// The property Type is required.
type ServerEventSelectorParam struct {
	// The outer Live server event type. Use 'response.event' for Responses events.
	Type string `json:"type" api:"required"`
	// The nested Responses event type. Required when type is 'response.event';
	// forbidden for other event types.
	ResponseEvent param.Opt[string] `json:"response_event,omitzero"`
	paramObj
}

func (r ServerEventSelectorParam) MarshalJSON() (data []byte, err error) {
	type shadow ServerEventSelectorParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ServerEventSelectorParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Request that the Live session close. The terminal `session.closed` event
// contains the close reason and final usage.
//
// The property Type is required.
type SessionCloseEventParam struct {
	// Optional client identifier for correlating this command with a server event's
	// client_event_id or error.client_event_id.
	EventID param.Opt[string] `json:"event_id,omitzero"`
	// The Live client event type. Always `session.close`.
	//
	// This field can be elided, and will marshal its zero value as "session.close".
	Type constant.SessionClose `json:"type" default:"session.close"`
	paramObj
}

func (r SessionCloseEventParam) MarshalJSON() (data []byte, err error) {
	type shadow SessionCloseEventParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *SessionCloseEventParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Returned after the Live session finishes finalizing, with the close reason,
// final session snapshot, and cumulative audio usage. A connection closing without
// this event does not confirm successful finalization.
type SessionClosedEvent struct {
	// The unique ID of the Live server event.
	EventID string `json:"event_id" api:"required"`
	// Why the Live session ended: `close_requested` for an application close or hangup
	// request, `expired` for the session duration limit, `content` for a safety
	// filter, `remote_hangup` for a graceful remote disconnect, or `connection_lost`
	// for an unexpected primary or upstream disconnection.
	Reason SessionClosedEventReasonString `json:"reason" api:"required"`
	// The resolved Live session configuration and server-assigned session metadata.
	Session SessionResource `json:"session" api:"required"`
	// The event type, always `session.closed`.
	Type constant.SessionClosed `json:"type" default:"session.closed"`
	// The final cumulative Live audio usage after session finalization.
	Usage SessionUsage `json:"usage" api:"required"`
	// The event_id of the client command associated with this server event, when
	// supplied.
	ClientEventID string `json:"client_event_id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventID       respjson.Field
		Reason        respjson.Field
		Session       respjson.Field
		Type          respjson.Field
		Usage         respjson.Field
		ClientEventID respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r SessionClosedEvent) RawJSON() string { return r.JSON.raw }
func (r *SessionClosedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Why the Live session ended: `close_requested` for an application close or hangup
// request, `expired` for the session duration limit, `content` for a safety
// filter, `remote_hangup` for a graceful remote disconnect, or `connection_lost`
// for an unexpected primary or upstream disconnection.
type SessionClosedEventReasonString string

// Initial configuration for a Live session, including its model, conversation
// instructions, audio, and delegated task handling.
//
// The property Model is required.
type SessionConfigParam struct {
	// The Live model. Required in the session configuration for every transport; do
	// not pass it as a URL query parameter.
	Model SessionConfigModel `json:"model,omitzero" api:"required"`
	// Frontend instructions for voice, conversation, interruptions, and when to
	// delegate. Start with the
	// [Live prompting guide](https://developers.openai.com/api/docs/guides/live-prompting);
	// put business rules and tool workflows in a separate
	// [backend prompt](https://developers.openai.com/api/docs/guides/live-delegation#start-with-your-existing-backend-prompt).
	// Limited to 16,384 client-supplied tokens. Omitted or blank instructions use
	// server defaults. Immutable after startup.
	Instructions param.Opt[string] `json:"instructions,omitzero"`
	// Whether to store the session for later forking and recording download. Defaults
	// to false for new sessions.
	Store param.Opt[bool] `json:"store,omitzero"`
	// Who handles tasks delegated by the Live model. Omitted or null selects your
	// application; use `responses` or `decisions` to delegate to the corresponding
	// API.
	Delegation SessionConfigDelegationUnionParam `json:"delegation,omitzero"`
	// Startup audio configuration. Only primary WebSockets accept audio.format; WebRTC
	// and SIP negotiate their media format. Voice and format are immutable after
	// startup.
	Audio SessionConfigAudioParam `json:"audio,omitzero"`
	// Startup-only capabilities for an untrusted frontend attached to a unified WebRTC
	// session. Trusted sideband connections are unaffected.
	Client ClientConfigParam `json:"client,omitzero"`
	// Ordered text-only history supplied before startup. Supports developer, user, and
	// assistant messages with one text part each; at most 128 messages and a
	// configured rendered-token limit (8,192 by default).
	Input []InitialItemUnionParam `json:"input,omitzero"`
	// At most one Decisions trigger per transcript event type. Triggers run
	// independently of delegation. On update, omit triggers to keep the current
	// configuration, send a list to replace it, or send an empty list to disable
	// triggers.
	Triggers []DecisionTriggerParam `json:"triggers,omitzero"`
	paramObj
}

func (r SessionConfigParam) MarshalJSON() (data []byte, err error) {
	type shadow SessionConfigParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *SessionConfigParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The Live model. Required in the session configuration for every transport; do
// not pass it as a URL query parameter.
type SessionConfigModel = string

const (
	SessionConfigModelGPTLive1 SessionConfigModel = "gpt-live-1"
)

// Startup audio configuration. Only primary WebSockets accept audio.format; WebRTC
// and SIP negotiate their media format. Voice and format are immutable after
// startup.
type SessionConfigAudioParam struct {
	// Audio encoding and sample rate for audio sent and received over a Live WebSocket
	// connection. WebRTC and SIP negotiate their media format separately.
	Format AudioFormatUnionParam `json:"format,omitzero"`
	// The voice used for speech generated by the Live model.
	Output SessionConfigAudioOutputParam `json:"output,omitzero"`
	paramObj
}

func (r SessionConfigAudioParam) MarshalJSON() (data []byte, err error) {
	type shadow SessionConfigAudioParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *SessionConfigAudioParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The voice used for speech generated by the Live model.
type SessionConfigAudioOutputParam struct {
	// The voice used for Live speech, as a built-in voice name or a custom voice
	// object containing its ID. Defaults to `marin` and cannot change after startup.
	Voice SessionConfigAudioOutputVoiceUnionParam `json:"voice,omitzero"`
	paramObj
}

func (r SessionConfigAudioOutputParam) MarshalJSON() (data []byte, err error) {
	type shadow SessionConfigAudioOutputParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *SessionConfigAudioOutputParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type SessionConfigAudioOutputVoiceUnionParam struct {
	OfString param.Opt[string] `json:",omitzero,inline"`
	// Check if union is this variant with !param.IsOmitted(union.OfBuiltIn)
	OfBuiltIn     param.Opt[BuiltInVoice] `json:",omitzero,inline"`
	OfCustomVoice *CustomVoiceParam       `json:",omitzero,inline"`
	paramUnion
}

func (u SessionConfigAudioOutputVoiceUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfString, u.OfBuiltIn, u.OfCustomVoice)
}
func (u *SessionConfigAudioOutputVoiceUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type SessionConfigDelegationUnionParam struct {
	OfClient    *ClientDelegationParam                 `json:",omitzero,inline"`
	OfResponses *SessionConfigDelegationResponsesParam `json:",omitzero,inline"`
	OfDecisions *SessionConfigDelegationDecisionsParam `json:",omitzero,inline"`
	paramUnion
}

func (u SessionConfigDelegationUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfClient, u.OfResponses, u.OfDecisions)
}
func (u *SessionConfigDelegationUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *SessionConfigDelegationUnionParam) asAny() any {
	if !param.IsOmitted(u.OfClient) {
		return u.OfClient
	} else if !param.IsOmitted(u.OfResponses) {
		return u.OfResponses
	} else if !param.IsOmitted(u.OfDecisions) {
		return u.OfDecisions
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u SessionConfigDelegationUnionParam) GetResponses() *ResponsesDelegationConfigParam {
	if vt := u.OfResponses; vt != nil {
		return &vt.Responses
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u SessionConfigDelegationUnionParam) GetDecisions() *DecisionsDelegationConfigParam {
	if vt := u.OfDecisions; vt != nil {
		return &vt.Decisions
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u SessionConfigDelegationUnionParam) GetType() *string {
	if vt := u.OfClient; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfResponses; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfDecisions; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u SessionConfigDelegationUnionParam) GetAcknowledgementDelayMs() *int64 {
	if vt := u.OfClient; vt != nil && vt.AcknowledgementDelayMs.Valid() {
		return &vt.AcknowledgementDelayMs.Value
	} else if vt := u.OfResponses; vt != nil && vt.AcknowledgementDelayMs.Valid() {
		return &vt.AcknowledgementDelayMs.Value
	} else if vt := u.OfDecisions; vt != nil && vt.AcknowledgementDelayMs.Valid() {
		return &vt.AcknowledgementDelayMs.Value
	}
	return nil
}

func init() {
	apijson.RegisterUnion[SessionConfigDelegationUnionParam](
		"type",
		apijson.Discriminator[ClientDelegationParam]("client"),
		apijson.Discriminator[SessionConfigDelegationResponsesParam]("responses"),
		apijson.Discriminator[SessionConfigDelegationDecisionsParam]("decisions"),
	)
}

// Delegate tasks to a Responses model managed by the Live session.
//
// The properties Responses, Type are required.
type SessionConfigDelegationResponsesParam struct {
	// Backend model, prompt, and tools used when the Live session delegates a task to
	// Responses.
	Responses ResponsesDelegationConfigParam `json:"responses,omitzero" api:"required"`
	// Milliseconds to wait before requesting an automatic spoken acknowledgement for a
	// delegation. The model may still speak before this delay. Any commentary or
	// thinking appended to the session cancels all pending acknowledgements, including
	// those for other delegations. Defaults to 400 for Decisions sessions and 0
	// (request immediately) for other sessions. Omission in a session update preserves
	// the current value. Updates apply to subsequent delegations.
	AcknowledgementDelayMs param.Opt[int64] `json:"acknowledgement_delay_ms,omitzero"`
	// The delegation owner. Always `responses` for tasks handled by the Responses API.
	//
	// This field can be elided, and will marshal its zero value as "responses".
	Type constant.Responses `json:"type" default:"responses"`
	paramObj
}

func (r SessionConfigDelegationResponsesParam) MarshalJSON() (data []byte, err error) {
	type shadow SessionConfigDelegationResponsesParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *SessionConfigDelegationResponsesParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Decisions, Type are required.
type SessionConfigDelegationDecisionsParam struct {
	Decisions DecisionsDelegationConfigParam `json:"decisions" api:"required"`
	// Milliseconds to wait before requesting an automatic spoken acknowledgement for a
	// delegation. The model may still speak before this delay. Any commentary or
	// thinking appended to the session cancels all pending acknowledgements, including
	// those for other delegations. Defaults to 400. Updates apply to subsequent
	// delegations.
	AcknowledgementDelayMs param.Opt[int64] `json:"acknowledgement_delay_ms,omitzero"`
	// This field can be elided, and will marshal its zero value as "decisions".
	Type constant.Decisions `json:"type" default:"decisions"`
	paramObj
}

func (r SessionConfigDelegationDecisionsParam) MarshalJSON() (data []byte, err error) {
	type shadow SessionConfigDelegationDecisionsParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *SessionConfigDelegationDecisionsParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The resolved Live session configuration and server-assigned session metadata.
type SessionResource struct {
	// The unique ID of the Live session. Use this ID for sideband connections,
	// forking, and recording download.
	ID string `json:"id" api:"required"`
	// The Unix timestamp, in seconds, at which the Live session expires.
	ExpiresAt int64 `json:"expires_at" api:"required"`
	// The Live model. Required in the session configuration for every transport; do
	// not pass it as a URL query parameter.
	Model SessionResourceModel `json:"model" api:"required"`
	// The status of the session snapshot. Always `active`, including the final
	// snapshot in session.closed; use the event type to determine that the session has
	// closed.
	Status constant.Active `json:"status" default:"active"`
	// Startup audio configuration. Only primary WebSockets accept audio.format; WebRTC
	// and SIP negotiate their media format. Voice and format are immutable after
	// startup.
	Audio SessionResourceAudio `json:"audio"`
	// Startup-only capabilities for an untrusted frontend attached to a unified WebRTC
	// session. Trusted sideband connections are unaffected.
	Client ClientConfig `json:"client"`
	// Who handles tasks delegated by the Live model. Omitted or null selects your
	// application; use `responses` or `decisions` to delegate to the corresponding
	// API.
	Delegation SessionResourceDelegationUnion `json:"delegation" api:"nullable"`
	// Ordered text-only history supplied before startup. Supports developer, user, and
	// assistant messages with one text part each; at most 128 messages and a
	// configured rendered-token limit (8,192 by default).
	Input []InitialItemUnion `json:"input"`
	// Frontend instructions for voice, conversation, interruptions, and when to
	// delegate. Start with the
	// [Live prompting guide](https://developers.openai.com/api/docs/guides/live-prompting);
	// put business rules and tool workflows in a separate
	// [backend prompt](https://developers.openai.com/api/docs/guides/live-delegation#start-with-your-existing-backend-prompt).
	// Limited to 16,384 client-supplied tokens. Omitted or blank instructions use
	// server defaults. Immutable after startup.
	Instructions string `json:"instructions" api:"nullable"`
	// Whether to store the session for later forking and recording download. Defaults
	// to false for new sessions.
	Store bool `json:"store"`
	// At most one Decisions trigger per transcript event type. Triggers run
	// independently of delegation. On update, omit triggers to keep the current
	// configuration, send a list to replace it, or send an empty list to disable
	// triggers.
	Triggers []DecisionTrigger `json:"triggers"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID           respjson.Field
		ExpiresAt    respjson.Field
		Model        respjson.Field
		Status       respjson.Field
		Audio        respjson.Field
		Client       respjson.Field
		Delegation   respjson.Field
		Input        respjson.Field
		Instructions respjson.Field
		Store        respjson.Field
		Triggers     respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r SessionResource) RawJSON() string { return r.JSON.raw }
func (r *SessionResource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The Live model. Required in the session configuration for every transport; do
// not pass it as a URL query parameter.
type SessionResourceModel = string

const (
	SessionResourceModelGPTLive1 SessionResourceModel = "gpt-live-1"
)

// Startup audio configuration. Only primary WebSockets accept audio.format; WebRTC
// and SIP negotiate their media format. Voice and format are immutable after
// startup.
type SessionResourceAudio struct {
	// Audio encoding and sample rate for audio sent and received over a Live WebSocket
	// connection. WebRTC and SIP negotiate their media format separately.
	Format AudioFormatUnion `json:"format"`
	// The voice used for speech generated by the Live model.
	Output SessionResourceAudioOutput `json:"output"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Format      respjson.Field
		Output      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r SessionResourceAudio) RawJSON() string { return r.JSON.raw }
func (r *SessionResourceAudio) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The voice used for speech generated by the Live model.
type SessionResourceAudioOutput struct {
	// The voice used for Live speech, as a built-in voice name or a custom voice
	// object containing its ID. Defaults to `marin` and cannot change after startup.
	Voice SessionResourceAudioOutputVoiceUnion `json:"voice"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Voice       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r SessionResourceAudioOutput) RawJSON() string { return r.JSON.raw }
func (r *SessionResourceAudioOutput) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// SessionResourceAudioOutputVoiceUnion contains all possible properties and values
// from [string], [BuiltInVoice], [CustomVoice].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfString OfBuiltIn]
type SessionResourceAudioOutputVoiceUnion struct {
	// This field will be present if the value is a [string] instead of an object.
	OfString string `json:",inline"`
	// This field will be present if the value is a [BuiltInVoice] instead of an
	// object.
	OfBuiltIn     BuiltInVoice `json:",inline"`
	OfCustomVoice CustomVoice  `json:",inline"`
	JSON          struct {
		OfString      respjson.Field
		OfBuiltIn     respjson.Field
		OfCustomVoice respjson.Field
		raw           string
	} `json:"-"`
}

func (u SessionResourceAudioOutputVoiceUnion) AsString() (v string) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u SessionResourceAudioOutputVoiceUnion) AsBuiltIn() (v BuiltInVoice) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u SessionResourceAudioOutputVoiceUnion) AsCustomVoice() (v CustomVoice) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u SessionResourceAudioOutputVoiceUnion) RawJSON() string { return u.JSON.raw }

func (r *SessionResourceAudioOutputVoiceUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// SessionResourceDelegationUnion contains all possible properties and values from
// [ClientDelegation], [SessionResourceDelegationResponses],
// [SessionResourceDelegationDecisions].
//
// Use the [SessionResourceDelegationUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type SessionResourceDelegationUnion struct {
	OfClientDelegation ClientDelegation `json:",inline"`
	// Any of "client", "responses", "decisions".
	Type                                 string                             `json:"type"`
	AcknowledgementDelayMs               int64                              `json:"acknowledgement_delay_ms"`
	OfSessionResourceDelegationResponses SessionResourceDelegationResponses `json:",inline"`
	// This field is from variant [SessionResourceDelegationResponses].
	Responses                            ResponsesDelegationConfig          `json:"responses"`
	OfSessionResourceDelegationDecisions SessionResourceDelegationDecisions `json:",inline"`
	// This field is from variant [SessionResourceDelegationDecisions].
	Decisions DecisionsDelegationConfig `json:"decisions"`
	JSON      struct {
		OfClientDelegation                   respjson.Field
		Type                                 respjson.Field
		AcknowledgementDelayMs               respjson.Field
		OfSessionResourceDelegationResponses respjson.Field
		Responses                            respjson.Field
		OfSessionResourceDelegationDecisions respjson.Field
		Decisions                            respjson.Field
		raw                                  string
	} `json:"-"`
}

// anySessionResourceDelegation is implemented by each variant of
// [SessionResourceDelegationUnion] to add type safety for the return type of
// [SessionResourceDelegationUnion.AsAny]
type anySessionResourceDelegation interface {
	implSessionResourceDelegationUnion()
}

func (ClientDelegation) implSessionResourceDelegationUnion()                   {}
func (SessionResourceDelegationResponses) implSessionResourceDelegationUnion() {}
func (SessionResourceDelegationDecisions) implSessionResourceDelegationUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := SessionResourceDelegationUnion.AsAny().(type) {
//	case live.ClientDelegation:
//	case live.SessionResourceDelegationResponses:
//	case live.SessionResourceDelegationDecisions:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u SessionResourceDelegationUnion) AsAny() anySessionResourceDelegation {
	switch u.Type {
	case "client":
		return u.AsClientDelegation()
	case "responses":
		return u.AsSessionResourceDelegationResponses()
	case "decisions":
		return u.AsSessionResourceDelegationDecisions()
	}
	return nil
}

func (u SessionResourceDelegationUnion) AsClientDelegation() (v ClientDelegation) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u SessionResourceDelegationUnion) AsSessionResourceDelegationResponses() (v SessionResourceDelegationResponses) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u SessionResourceDelegationUnion) AsSessionResourceDelegationDecisions() (v SessionResourceDelegationDecisions) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u SessionResourceDelegationUnion) RawJSON() string { return u.JSON.raw }

func (r *SessionResourceDelegationUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Delegate tasks to a Responses model managed by the Live session.
type SessionResourceDelegationResponses struct {
	// Milliseconds to wait before requesting an automatic spoken acknowledgement for a
	// delegation. The model may still speak before this delay. Any commentary or
	// thinking appended to the session cancels all pending acknowledgements, including
	// those for other delegations. Defaults to 400 for Decisions sessions and 0
	// (request immediately) for other sessions. Omission in a session update preserves
	// the current value. Updates apply to subsequent delegations.
	AcknowledgementDelayMs int64 `json:"acknowledgement_delay_ms"`
	// Backend model, prompt, and tools used when the Live session delegates a task to
	// Responses.
	Responses ResponsesDelegationConfig `json:"responses" api:"required"`
	// The delegation owner. Always `responses` for tasks handled by the Responses API.
	Type constant.Responses `json:"type" default:"responses"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		AcknowledgementDelayMs respjson.Field
		Responses              respjson.Field
		Type                   respjson.Field
		ExtraFields            map[string]respjson.Field
		raw                    string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r SessionResourceDelegationResponses) RawJSON() string { return r.JSON.raw }
func (r *SessionResourceDelegationResponses) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type SessionResourceDelegationDecisions struct {
	// Milliseconds to wait before requesting an automatic spoken acknowledgement for a
	// delegation. The model may still speak before this delay. Any commentary or
	// thinking appended to the session cancels all pending acknowledgements, including
	// those for other delegations. Defaults to 400. Updates apply to subsequent
	// delegations.
	AcknowledgementDelayMs int64                     `json:"acknowledgement_delay_ms"`
	Decisions              DecisionsDelegationConfig `json:"decisions" api:"required"`
	Type                   constant.Decisions        `json:"type" default:"decisions"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		AcknowledgementDelayMs respjson.Field
		Decisions              respjson.Field
		Type                   respjson.Field
		ExtraFields            map[string]respjson.Field
		raw                    string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r SessionResourceDelegationDecisions) RawJSON() string { return r.JSON.raw }
func (r *SessionResourceDelegationDecisions) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Start a Live session on a primary WebSocket. Send this event before other
// commands and wait for `session.started`.
//
// The properties Session, Type are required.
type SessionStartEventParam struct {
	// Initial configuration for a primary WebSocket. Send session.start first and wait
	// for session.started before application commands. WebRTC creation already starts
	// the session; do not send this event again on its data channel.
	Session SessionConfigParam `json:"session" api:"required"`
	// Optional client identifier for correlating this command with a server event's
	// client_event_id or error.client_event_id.
	EventID param.Opt[string] `json:"event_id,omitzero"`
	// The Live client event type. Always `session.start`.
	//
	// This field can be elided, and will marshal its zero value as "session.start".
	Type constant.SessionStart `json:"type" default:"session.start"`
	paramObj
}

func (r SessionStartEventParam) MarshalJSON() (data []byte, err error) {
	type shadow SessionStartEventParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *SessionStartEventParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Returned when a Live session has started. Contains the resolved session
// configuration, including server defaults.
type SessionStartedEvent struct {
	// The unique ID of the Live server event.
	EventID string `json:"event_id" api:"required"`
	// The resolved Live session configuration and server-assigned session metadata.
	Session SessionResource `json:"session" api:"required"`
	// The event type, always `session.started`.
	Type constant.SessionStarted `json:"type" default:"session.started"`
	// The event_id of the client command associated with this server event, when
	// supplied.
	ClientEventID string `json:"client_event_id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventID       respjson.Field
		Session       respjson.Field
		Type          respjson.Field
		ClientEventID respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r SessionStartedEvent) RawJSON() string { return r.JSON.raw }
func (r *SessionStartedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Changes to an active Live session's delegation settings and transcript triggers.
type SessionUpdateConfigParam struct {
	// Delegation settings to update. The delegation type must match the current
	// session; omitted settings retain their values.
	Delegation SessionUpdateConfigDelegationUnionParam `json:"delegation,omitzero"`
	// At most one Decisions trigger per transcript event type. Triggers run
	// independently of delegation. On update, omit triggers to keep the current
	// configuration, send a list to replace it, or send an empty list to disable
	// triggers.
	Triggers []DecisionTriggerParam `json:"triggers,omitzero"`
	paramObj
}

func (r SessionUpdateConfigParam) MarshalJSON() (data []byte, err error) {
	type shadow SessionUpdateConfigParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *SessionUpdateConfigParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type SessionUpdateConfigDelegationUnionParam struct {
	OfClient    *ClientDelegationParam                       `json:",omitzero,inline"`
	OfResponses *SessionUpdateConfigDelegationResponsesParam `json:",omitzero,inline"`
	OfDecisions *SessionUpdateConfigDelegationDecisionsParam `json:",omitzero,inline"`
	paramUnion
}

func (u SessionUpdateConfigDelegationUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfClient, u.OfResponses, u.OfDecisions)
}
func (u *SessionUpdateConfigDelegationUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *SessionUpdateConfigDelegationUnionParam) asAny() any {
	if !param.IsOmitted(u.OfClient) {
		return u.OfClient
	} else if !param.IsOmitted(u.OfResponses) {
		return u.OfResponses
	} else if !param.IsOmitted(u.OfDecisions) {
		return u.OfDecisions
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u SessionUpdateConfigDelegationUnionParam) GetResponses() *ResponsesDelegationUpdateConfigParam {
	if vt := u.OfResponses; vt != nil {
		return &vt.Responses
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u SessionUpdateConfigDelegationUnionParam) GetDecisions() *DecisionsDelegationUpdateConfigParam {
	if vt := u.OfDecisions; vt != nil {
		return &vt.Decisions
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u SessionUpdateConfigDelegationUnionParam) GetType() *string {
	if vt := u.OfClient; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfResponses; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfDecisions; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u SessionUpdateConfigDelegationUnionParam) GetAcknowledgementDelayMs() *int64 {
	if vt := u.OfClient; vt != nil && vt.AcknowledgementDelayMs.Valid() {
		return &vt.AcknowledgementDelayMs.Value
	} else if vt := u.OfResponses; vt != nil && vt.AcknowledgementDelayMs.Valid() {
		return &vt.AcknowledgementDelayMs.Value
	} else if vt := u.OfDecisions; vt != nil && vt.AcknowledgementDelayMs.Valid() {
		return &vt.AcknowledgementDelayMs.Value
	}
	return nil
}

func init() {
	apijson.RegisterUnion[SessionUpdateConfigDelegationUnionParam](
		"type",
		apijson.Discriminator[ClientDelegationParam]("client"),
		apijson.Discriminator[SessionUpdateConfigDelegationResponsesParam]("responses"),
		apijson.Discriminator[SessionUpdateConfigDelegationDecisionsParam]("decisions"),
	)
}

// Update the Responses backend for an existing Live session without changing
// delegation ownership.
//
// The property Type is required.
type SessionUpdateConfigDelegationResponsesParam struct {
	// Milliseconds to wait before requesting an automatic spoken acknowledgement for a
	// delegation. The model may still speak before this delay. Any commentary or
	// thinking appended to the session cancels all pending acknowledgements, including
	// those for other delegations. Defaults to 400 for Decisions sessions and 0
	// (request immediately) for other sessions. Omission in a session update preserves
	// the current value. Updates apply to subsequent delegations.
	AcknowledgementDelayMs param.Opt[int64] `json:"acknowledgement_delay_ms,omitzero"`
	// Responses backend settings to update. Omitted settings keep their existing
	// values.
	Responses ResponsesDelegationUpdateConfigParam `json:"responses,omitzero"`
	// The delegation owner. Always `responses` for tasks handled by the Responses API.
	//
	// This field can be elided, and will marshal its zero value as "responses".
	Type constant.Responses `json:"type" default:"responses"`
	paramObj
}

func (r SessionUpdateConfigDelegationResponsesParam) MarshalJSON() (data []byte, err error) {
	type shadow SessionUpdateConfigDelegationResponsesParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *SessionUpdateConfigDelegationResponsesParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The property Type is required.
type SessionUpdateConfigDelegationDecisionsParam struct {
	// Milliseconds to wait before requesting an automatic spoken acknowledgement for a
	// delegation. The model may still speak before this delay. Any commentary or
	// thinking appended to the session cancels all pending acknowledgements, including
	// those for other delegations. Defaults to 400 for Decisions sessions and 0
	// (request immediately) for other sessions. Omission in a session update preserves
	// the current value. Updates apply to subsequent delegations.
	AcknowledgementDelayMs param.Opt[int64]                     `json:"acknowledgement_delay_ms,omitzero"`
	Decisions              DecisionsDelegationUpdateConfigParam `json:"decisions,omitzero"`
	// This field can be elided, and will marshal its zero value as "decisions".
	Type constant.Decisions `json:"type" default:"decisions"`
	paramObj
}

func (r SessionUpdateConfigDelegationDecisionsParam) MarshalJSON() (data []byte, err error) {
	type shadow SessionUpdateConfigDelegationDecisionsParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *SessionUpdateConfigDelegationDecisionsParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Update the delegation settings or transcript triggers of an active Live session.
// The server acknowledges accepted changes with `session.updated`.
//
// The properties Session, Type are required.
type SessionUpdateEventParam struct {
	// Sparse delegation and trigger updates. Omitted settings retain their values. The
	// delegation type cannot change. Model, frontend instructions, audio, and startup
	// input are immutable.
	Session SessionUpdateConfigParam `json:"session" api:"required"`
	// Optional client identifier for correlating this command with a server event's
	// client_event_id or error.client_event_id.
	EventID param.Opt[string] `json:"event_id,omitzero"`
	// The Live client event type. Always `session.update`.
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

// Returned when a Live session update is accepted. Contains the resolved session
// configuration after the update.
type SessionUpdatedEvent struct {
	// The unique ID of the Live server event.
	EventID string `json:"event_id" api:"required"`
	// The resolved Live session configuration and server-assigned session metadata.
	Session SessionResource `json:"session" api:"required"`
	// The event type, always `session.updated`.
	Type constant.SessionUpdated `json:"type" default:"session.updated"`
	// The event_id of the client command associated with this server event, when
	// supplied.
	ClientEventID string `json:"client_event_id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventID       respjson.Field
		Session       respjson.Field
		Type          respjson.Field
		ClientEventID respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r SessionUpdatedEvent) RawJSON() string { return r.JSON.raw }
func (r *SessionUpdatedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Cumulative audio duration for a Live session. Values are totals for the session,
// not increments to sum across usage events.
type SessionUsage struct {
	// The cumulative Live audio duration in seconds. Do not sum this value across
	// usage events.
	Seconds float64 `json:"seconds" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Seconds     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r SessionUsage) RawJSON() string { return r.JSON.raw }
func (r *SessionUsage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Reports cumulative Live audio usage and, when available, the most recent
// context-window usage. Delegated Responses token usage is reported separately in
// response.event events.
type SessionUsageUpdatedEvent struct {
	// The unique ID of the Live server event.
	EventID string `json:"event_id" api:"required"`
	// The event type, always `session.usage.updated`.
	Type constant.SessionUsageUpdated `json:"type" default:"session.usage.updated"`
	// The cumulative Live audio usage so far.
	Usage SessionUsage `json:"usage" api:"required"`
	// The event_id of the client command associated with this server event, when
	// supplied.
	ClientEventID string `json:"client_event_id"`
	// The latest measured Live context-window usage. Omitted when the context limit is
	// unknown.
	ContextWindow SessionUsageUpdatedEventContextWindow `json:"context_window"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventID       respjson.Field
		Type          respjson.Field
		Usage         respjson.Field
		ClientEventID respjson.Field
		ContextWindow respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r SessionUsageUpdatedEvent) RawJSON() string { return r.JSON.raw }
func (r *SessionUsageUpdatedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The latest measured Live context-window usage. Omitted when the context limit is
// unknown.
type SessionUsageUpdatedEventContextWindow struct {
	// The latest active context token count divided by the Live model context limit.
	// Can decrease after compaction and may lag between measured audio frames.
	UsageRatio float64 `json:"usage_ratio" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		UsageRatio  respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r SessionUsageUpdatedEventContextWindow) RawJSON() string { return r.JSON.raw }
func (r *SessionUsageUpdatedEventContextWindow) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Provide silent reasoning or progress context to the Live model, optionally for
// an existing client delegation. Rejected while an audio fixture is active.
//
// The properties Content, DelegationID, Type are required.
type ThinkingAppendEventParam struct {
	// Required, nullable. Set null for general session context, or use the ID from
	// session.delegation.created for an existing client delegation. Non-null IDs are
	// not accepted with Responses delegation.
	DelegationID param.Opt[string] `json:"delegation_id,omitzero" api:"required"`
	// Silent reasoning or progress context, limited to 500 tokens. It does not
	// directly request speech, but can influence later speech and is not a secrecy
	// boundary.
	Content string `json:"content" api:"required"`
	// Optional client identifier for correlating this command with a server event's
	// client_event_id or error.client_event_id.
	EventID param.Opt[string] `json:"event_id,omitzero"`
	// The Live client event type. Always `session.thinking.append`.
	//
	// This field can be elided, and will marshal its zero value as
	// "session.thinking.append".
	Type constant.SessionThinkingAppend `json:"type" default:"session.thinking.append"`
	paramObj
}

func (r ThinkingAppendEventParam) MarshalJSON() (data []byte, err error) {
	type shadow ThinkingAppendEventParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ThinkingAppendEventParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Returned when a session.thinking.append command is accepted into the Live
// session timeline. Acknowledges the added reasoning context without guaranteeing
// any spoken output.
type ThinkingAppendedEvent struct {
	// The end of this event on the Live session timeline, in milliseconds from the
	// beginning of the session. For appended context, this can equal start_ms.
	EndMs int64 `json:"end_ms" api:"required"`
	// The unique ID of the Live server event.
	EventID string `json:"event_id" api:"required"`
	// The start of this event on the Live session timeline, in milliseconds from the
	// beginning of the session.
	StartMs int64 `json:"start_ms" api:"required"`
	// The event type, always `session.thinking.appended`.
	Type constant.SessionThinkingAppended `json:"type" default:"session.thinking.appended"`
	// The event_id of the client command associated with this server event, when
	// supplied.
	ClientEventID string `json:"client_event_id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EndMs         respjson.Field
		EventID       respjson.Field
		StartMs       respjson.Field
		Type          respjson.Field
		ClientEventID respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ThinkingAppendedEvent) RawJSON() string { return r.JSON.raw }
func (r *ThinkingAppendedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The created Live session identifier and WebRTC answer. Apply transport.sdp as
// the peer's remote answer and wait for session.started on the data channel before
// sending commands.
type LiveNewResponse struct {
	// The newly created Live session. Use its ID for session controls and sideband
	// connections.
	Session LiveNewResponseSession `json:"session" api:"required"`
	// WebRTC transport with the SDP answer.
	Transport LiveNewResponseTransport `json:"transport" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Session     respjson.Field
		Transport   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r LiveNewResponse) RawJSON() string { return r.JSON.raw }
func (r *LiveNewResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The newly created Live session. Use its ID for session controls and sideband
// connections.
type LiveNewResponseSession struct {
	// Opaque session identifier. Preserve the returned value unchanged, including its
	// prefix.
	ID string `json:"id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r LiveNewResponseSession) RawJSON() string { return r.JSON.raw }
func (r *LiveNewResponseSession) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// WebRTC transport with the SDP answer.
type LiveNewResponseTransport struct {
	// Session Description Protocol message for the WebRTC connection.
	Sdp string `json:"sdp" api:"required"`
	// The transport used for the Live session. Always `webrtc`.
	Type constant.Webrtc `json:"type" default:"webrtc"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Sdp         respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r LiveNewResponseTransport) RawJSON() string { return r.JSON.raw }
func (r *LiveNewResponseTransport) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type LiveNewParams struct {
	// Startup configuration for the Live session.
	Session MediaSessionConfigParam `json:"session,omitzero" api:"required"`
	// WebRTC transport with the browser's SDP offer.
	Transport LiveNewParamsTransport `json:"transport,omitzero" api:"required"`
	paramObj
}

func (r LiveNewParams) MarshalJSON() (data []byte, err error) {
	type shadow LiveNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *LiveNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// WebRTC transport with the browser's SDP offer.
//
// The properties Sdp, Type are required.
type LiveNewParamsTransport struct {
	// Session Description Protocol message for the WebRTC connection.
	Sdp string `json:"sdp" api:"required"`
	// The transport used for the Live session. Always `webrtc`.
	//
	// This field can be elided, and will marshal its zero value as "webrtc".
	Type constant.Webrtc `json:"type" default:"webrtc"`
	paramObj
}

func (r LiveNewParamsTransport) MarshalJSON() (data []byte, err error) {
	type shadow LiveNewParamsTransport
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *LiveNewParamsTransport) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

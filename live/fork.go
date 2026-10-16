// File generated from our OpenAPI spec by Castiron. See CONTRIBUTING.md for details.

package live

import (
	"encoding/json"

	"github.com/openai/openai-go/v3/internal/apijson"
	"github.com/openai/openai-go/v3/internal/paramutil"
	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/packages/respjson"
	"github.com/openai/openai-go/v3/responses"
	"github.com/openai/openai-go/v3/shared/constant"
)

// ForkService contains methods and other services that help with interacting with
// the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewForkService] method instead.
type ForkService struct {
	Options []option.RequestOption
}

// NewForkService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewForkService(opts ...option.RequestOption) (r ForkService) {
	r = ForkService{}
	r.Options = requestconfig.InheritedOptions(opts...)
	return
}

func ForkClientEventParamOfSessionStart(session ForkSessionConfigParam) ForkClientEventUnionParam {
	var sessionStart ForkSessionStartEventParam
	sessionStart.Session = session
	return ForkClientEventUnionParam{OfSessionStart: &sessionStart}
}

func ForkClientEventParamOfSessionUpdate(session SessionUpdateConfigParam) ForkClientEventUnionParam {
	var sessionUpdate SessionUpdateEventParam
	sessionUpdate.Session = session
	return ForkClientEventUnionParam{OfSessionUpdate: &sessionUpdate}
}

func ForkClientEventParamOfSessionInputAudioAppend(audio string) ForkClientEventUnionParam {
	var sessionInputAudioAppend InputAudioAppendEventParam
	sessionInputAudioAppend.Audio = audio
	return ForkClientEventUnionParam{OfSessionInputAudioAppend: &sessionInputAudioAppend}
}

func ForkClientEventParamOfSessionOutputAudioFixturePlay(audio string) ForkClientEventUnionParam {
	var sessionOutputAudioFixturePlay OutputAudioFixturePlayEventParam
	sessionOutputAudioFixturePlay.Audio = audio
	return ForkClientEventUnionParam{OfSessionOutputAudioFixturePlay: &sessionOutputAudioFixturePlay}
}

func ForkClientEventParamOfSessionInstructionsAppend(content string, delegationID string) ForkClientEventUnionParam {
	var sessionInstructionsAppend InstructionsAppendEventParam
	sessionInstructionsAppend.Content = content
	sessionInstructionsAppend.DelegationID = param.NewOpt(delegationID)
	return ForkClientEventUnionParam{OfSessionInstructionsAppend: &sessionInstructionsAppend}
}

func ForkClientEventParamOfSessionThinkingAppend(content string, delegationID string) ForkClientEventUnionParam {
	var sessionThinkingAppend ThinkingAppendEventParam
	sessionThinkingAppend.Content = content
	sessionThinkingAppend.DelegationID = param.NewOpt(delegationID)
	return ForkClientEventUnionParam{OfSessionThinkingAppend: &sessionThinkingAppend}
}

func ForkClientEventParamOfSessionCommentaryAppend(content string, delegationID string) ForkClientEventUnionParam {
	var sessionCommentaryAppend CommentaryAppendEventParam
	sessionCommentaryAppend.Content = content
	sessionCommentaryAppend.DelegationID = param.NewOpt(delegationID)
	return ForkClientEventUnionParam{OfSessionCommentaryAppend: &sessionCommentaryAppend}
}

func ForkClientEventParamOfResponseItemCreate[
	T responses.EasyInputMessageParam | responses.ResponseInputItemMessageParam | responses.ResponseOutputMessageParam | responses.ResponseFileSearchToolCallParam | responses.ResponseComputerToolCallParam | responses.ResponseInputItemComputerCallOutputParam | responses.ResponseFunctionWebSearchParam | responses.ResponseFunctionToolCallParam | responses.ResponseInputItemFunctionCallOutputParam | responses.ResponseInputItemToolSearchCallParam | responses.ResponseToolSearchOutputItemParam | responses.ResponseInputItemAdditionalToolsParam | responses.ResponseConfigurationUpdateItemParam | responses.ResponseReasoningItemParam | responses.ResponseCompactionItemParam | responses.ResponseInputItemImageGenerationCallParam | responses.ResponseCodeInterpreterToolCallParam | responses.ResponseInputItemLocalShellCallParam | responses.ResponseInputItemLocalShellCallOutputParam | responses.ResponseInputItemShellCallParam | responses.ResponseInputItemShellCallOutputParam | responses.ResponseInputItemApplyPatchCallParam | responses.ResponseInputItemApplyPatchCallOutputParam | responses.ResponseInputItemMcpListToolsParam | responses.ResponseInputItemMcpApprovalRequestParam | responses.ResponseInputItemMcpApprovalResponseParam | responses.ResponseInputItemMcpCallParam | responses.ResponseCustomToolCallOutputParam | responses.ResponseCustomToolCallParam | responses.ResponseInputItemCompactionTriggerParam | responses.ResponseInputItemItemReferenceParam | responses.ResponseInputItemProgramParam | responses.ResponseInputItemProgramOutputParam,
](item T) ForkClientEventUnionParam {
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
	return ForkClientEventUnionParam{OfResponseItemCreate: &responseItemCreate}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ForkClientEventUnionParam struct {
	OfSessionStart                  *ForkSessionStartEventParam       `json:",omitzero,inline"`
	OfSessionUpdate                 *SessionUpdateEventParam          `json:",omitzero,inline"`
	OfSessionInputAudioAppend       *InputAudioAppendEventParam       `json:",omitzero,inline"`
	OfSessionInputAudioMute         *InputAudioMuteEventParam         `json:",omitzero,inline"`
	OfSessionInputAudioUnmute       *InputAudioUnmuteEventParam       `json:",omitzero,inline"`
	OfSessionOutputAudioFixturePlay *OutputAudioFixturePlayEventParam `json:",omitzero,inline"`
	OfSessionInstructionsAppend     *InstructionsAppendEventParam     `json:",omitzero,inline"`
	OfSessionThinkingAppend         *ThinkingAppendEventParam         `json:",omitzero,inline"`
	OfSessionCommentaryAppend       *CommentaryAppendEventParam       `json:",omitzero,inline"`
	OfResponseItemCreate            *ResponseItemCreateEventParam     `json:",omitzero,inline"`
	OfResponseCreate                *ResponseCreateEventParam         `json:",omitzero,inline"`
	OfSessionClose                  *SessionCloseEventParam           `json:",omitzero,inline"`
	paramUnion
}

func (u ForkClientEventUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfSessionStart,
		u.OfSessionUpdate,
		u.OfSessionInputAudioAppend,
		u.OfSessionInputAudioMute,
		u.OfSessionInputAudioUnmute,
		u.OfSessionOutputAudioFixturePlay,
		u.OfSessionInstructionsAppend,
		u.OfSessionThinkingAppend,
		u.OfSessionCommentaryAppend,
		u.OfResponseItemCreate,
		u.OfResponseCreate,
		u.OfSessionClose)
}
func (u *ForkClientEventUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u ForkClientEventUnionParam) GetItem() *responses.ResponseInputItemUnionParam {
	if vt := u.OfResponseItemCreate; vt != nil {
		return &vt.Item
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ForkClientEventUnionParam) GetType() *string {
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
func (u ForkClientEventUnionParam) GetEventID() *string {
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
func (u ForkClientEventUnionParam) GetAudio() *string {
	if vt := u.OfSessionInputAudioAppend; vt != nil {
		return (*string)(&vt.Audio)
	} else if vt := u.OfSessionOutputAudioFixturePlay; vt != nil {
		return (*string)(&vt.Audio)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ForkClientEventUnionParam) GetContent() *string {
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
func (u ForkClientEventUnionParam) GetDelegationID() *string {
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
func (u ForkClientEventUnionParam) GetSession() (res forkClientEventUnionParamSession) {
	if vt := u.OfSessionStart; vt != nil {
		res.any = &vt.Session
	} else if vt := u.OfSessionUpdate; vt != nil {
		res.any = &vt.Session
	}
	return
}

// Can have the runtime types [*ForkSessionConfigParam],
// [*SessionUpdateConfigParam]
type forkClientEventUnionParamSession struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *live.ForkSessionConfigParam:
//	case *live.SessionUpdateConfigParam:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u forkClientEventUnionParamSession) AsAny() any { return u.any }

// Returns a pointer to the underlying variant's property, if present.
func (u forkClientEventUnionParamSession) GetAudio() *ForkSessionConfigAudioParam {
	switch vt := u.any.(type) {
	case *ForkSessionConfigParam:
		return &vt.Audio
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u forkClientEventUnionParamSession) GetClient() *ClientConfigParam {
	switch vt := u.any.(type) {
	case *ForkSessionConfigParam:
		return &vt.Client
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u forkClientEventUnionParamSession) GetStore() *bool {
	switch vt := u.any.(type) {
	case *ForkSessionConfigParam:
		return paramutil.AddrIfPresent(vt.Store)
	}
	return nil
}

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u forkClientEventUnionParamSession) GetDelegation() (res forkClientEventUnionParamSessionDelegation) {
	switch vt := u.any.(type) {
	case *ForkSessionConfigParam:
		res.any = &vt.Delegation
	case *SessionUpdateConfigParam:
		res.any = vt.Delegation.asAny()
	}
	return res
}

// Can have the runtime types [*ForkSessionConfigDelegationParam],
// [*ClientDelegationParam], [*SessionUpdateConfigDelegationResponsesParam]
type forkClientEventUnionParamSessionDelegation struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *live.ForkSessionConfigDelegationParam:
//	case *live.ClientDelegationParam:
//	case *live.SessionUpdateConfigDelegationResponsesParam:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u forkClientEventUnionParamSessionDelegation) AsAny() any { return u.any }

// Returns a pointer to the underlying variant's property, if present.
func (u forkClientEventUnionParamSessionDelegation) GetType() *string {
	switch vt := u.any.(type) {
	case *ForkSessionConfigDelegationParam:
		return (*string)(&vt.Type)
	case *ClientDelegationParam:
		return (*string)(&vt.Type)
	case *SessionUpdateConfigDelegationResponsesParam:
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's Responses property, if present.
func (u forkClientEventUnionParamSessionDelegation) GetResponses() *ResponsesDelegationUpdateConfigParam {
	switch vt := u.any.(type) {
	case *ForkSessionConfigDelegationParam:
		return &vt.Responses
	case *SessionUpdateConfigDelegationResponsesParam:
		return &vt.Responses
	}
	return nil
}

func init() {
	apijson.RegisterUnion[ForkClientEventUnionParam](
		"type",
		apijson.Discriminator[ForkSessionStartEventParam]("session.start"),
		apijson.Discriminator[SessionUpdateEventParam]("session.update"),
		apijson.Discriminator[InputAudioAppendEventParam]("session.input_audio.append"),
		apijson.Discriminator[InputAudioMuteEventParam]("session.input_audio.mute"),
		apijson.Discriminator[InputAudioUnmuteEventParam]("session.input_audio.unmute"),
		apijson.Discriminator[OutputAudioFixturePlayEventParam]("session.output_audio.fixture.play"),
		apijson.Discriminator[InstructionsAppendEventParam]("session.instructions.append"),
		apijson.Discriminator[ThinkingAppendEventParam]("session.thinking.append"),
		apijson.Discriminator[CommentaryAppendEventParam]("session.commentary.append"),
		apijson.Discriminator[ResponseItemCreateEventParam]("response.item.create"),
		apijson.Discriminator[ResponseCreateEventParam]("response.create"),
		apijson.Discriminator[SessionCloseEventParam]("session.close"),
	)
}

// ForkServerEventUnion contains all possible properties and values from
// [SessionStartedEvent], [SessionUpdatedEvent], [InputAudioMutedEvent],
// [InputAudioUnmutedEvent], [OutputAudioFixturePlayingEvent],
// [OutputAudioFixturePlayedEvent], [InstructionsAppendedEvent],
// [ThinkingAppendedEvent], [CommentaryAppendedEvent],
// [ForkServerEventSessionInputAudioAppend], [OutputAudioDeltaEvent],
// [InputTranscriptDeltaEvent], [OutputTranscriptDeltaEvent],
// [DelegationCreatedEvent], [ResponseEvent], [SessionUsageUpdatedEvent],
// [SessionClosedEvent], [ErrorEvent], [InfoEvent],
// [ForkServerEventTransportDtmfReceived], [ForkServerEventTransportDtmfSend],
// [ForkServerEventTransportRinging], [ForkServerEventTransportAnswered],
// [ForkServerEventTransportFailed].
//
// Use the [ForkServerEventUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ForkServerEventUnion struct {
	OfSessionStartedEvent                    SessionStartedEvent                    `json:",inline"`
	OfSessionUpdatedEvent                    SessionUpdatedEvent                    `json:",inline"`
	OfInputAudioMutedEvent                   InputAudioMutedEvent                   `json:",inline"`
	OfInputAudioUnmutedEvent                 InputAudioUnmutedEvent                 `json:",inline"`
	OfOutputAudioFixturePlayingEvent         OutputAudioFixturePlayingEvent         `json:",inline"`
	OfOutputAudioFixturePlayedEvent          OutputAudioFixturePlayedEvent          `json:",inline"`
	OfInstructionsAppendedEvent              InstructionsAppendedEvent              `json:",inline"`
	OfThinkingAppendedEvent                  ThinkingAppendedEvent                  `json:",inline"`
	OfCommentaryAppendedEvent                CommentaryAppendedEvent                `json:",inline"`
	OfForkServerEventSessionInputAudioAppend ForkServerEventSessionInputAudioAppend `json:",inline"`
	OfOutputAudioDeltaEvent                  OutputAudioDeltaEvent                  `json:",inline"`
	OfInputTranscriptDeltaEvent              InputTranscriptDeltaEvent              `json:",inline"`
	OfOutputTranscriptDeltaEvent             OutputTranscriptDeltaEvent             `json:",inline"`
	OfDelegationCreatedEvent                 DelegationCreatedEvent                 `json:",inline"`
	OfResponseEvent                          ResponseEvent                          `json:",inline"`
	OfSessionUsageUpdatedEvent               SessionUsageUpdatedEvent               `json:",inline"`
	OfSessionClosedEvent                     SessionClosedEvent                     `json:",inline"`
	OfErrorEvent                             ErrorEvent                             `json:",inline"`
	OfInfoEvent                              InfoEvent                              `json:",inline"`
	OfForkServerEventTransportDtmfReceived   ForkServerEventTransportDtmfReceived   `json:",inline"`
	OfForkServerEventTransportDtmfSend       ForkServerEventTransportDtmfSend       `json:",inline"`
	OfForkServerEventTransportRinging        ForkServerEventTransportRinging        `json:",inline"`
	OfForkServerEventTransportAnswered       ForkServerEventTransportAnswered       `json:",inline"`
	OfForkServerEventTransportFailed         ForkServerEventTransportFailed         `json:",inline"`
	Type                                     string                                 `json:"type"`
	JSON                                     struct {
		OfSessionStartedEvent                    respjson.Field
		OfSessionUpdatedEvent                    respjson.Field
		OfInputAudioMutedEvent                   respjson.Field
		OfInputAudioUnmutedEvent                 respjson.Field
		OfOutputAudioFixturePlayingEvent         respjson.Field
		OfOutputAudioFixturePlayedEvent          respjson.Field
		OfInstructionsAppendedEvent              respjson.Field
		OfThinkingAppendedEvent                  respjson.Field
		OfCommentaryAppendedEvent                respjson.Field
		OfForkServerEventSessionInputAudioAppend respjson.Field
		OfOutputAudioDeltaEvent                  respjson.Field
		OfInputTranscriptDeltaEvent              respjson.Field
		OfOutputTranscriptDeltaEvent             respjson.Field
		OfDelegationCreatedEvent                 respjson.Field
		OfResponseEvent                          respjson.Field
		OfSessionUsageUpdatedEvent               respjson.Field
		OfSessionClosedEvent                     respjson.Field
		OfErrorEvent                             respjson.Field
		OfInfoEvent                              respjson.Field
		OfForkServerEventTransportDtmfReceived   respjson.Field
		OfForkServerEventTransportDtmfSend       respjson.Field
		OfForkServerEventTransportRinging        respjson.Field
		OfForkServerEventTransportAnswered       respjson.Field
		OfForkServerEventTransportFailed         respjson.Field
		Type                                     respjson.Field
		raw                                      string
	} `json:"-"`
}

// anyForkServerEvent is implemented by each variant of [ForkServerEventUnion] to
// add type safety for the return type of [ForkServerEventUnion.AsAny]
type anyForkServerEvent interface {
	implForkServerEventUnion()
}

func (SessionStartedEvent) implForkServerEventUnion()                    {}
func (SessionUpdatedEvent) implForkServerEventUnion()                    {}
func (InputAudioMutedEvent) implForkServerEventUnion()                   {}
func (InputAudioUnmutedEvent) implForkServerEventUnion()                 {}
func (OutputAudioFixturePlayingEvent) implForkServerEventUnion()         {}
func (OutputAudioFixturePlayedEvent) implForkServerEventUnion()          {}
func (InstructionsAppendedEvent) implForkServerEventUnion()              {}
func (ThinkingAppendedEvent) implForkServerEventUnion()                  {}
func (CommentaryAppendedEvent) implForkServerEventUnion()                {}
func (ForkServerEventSessionInputAudioAppend) implForkServerEventUnion() {}
func (OutputAudioDeltaEvent) implForkServerEventUnion()                  {}
func (InputTranscriptDeltaEvent) implForkServerEventUnion()              {}
func (OutputTranscriptDeltaEvent) implForkServerEventUnion()             {}
func (DelegationCreatedEvent) implForkServerEventUnion()                 {}
func (ResponseEvent) implForkServerEventUnion()                          {}
func (SessionUsageUpdatedEvent) implForkServerEventUnion()               {}
func (SessionClosedEvent) implForkServerEventUnion()                     {}
func (ErrorEvent) implForkServerEventUnion()                             {}
func (InfoEvent) implForkServerEventUnion()                              {}
func (ForkServerEventTransportDtmfReceived) implForkServerEventUnion()   {}
func (ForkServerEventTransportDtmfSend) implForkServerEventUnion()       {}
func (ForkServerEventTransportRinging) implForkServerEventUnion()        {}
func (ForkServerEventTransportAnswered) implForkServerEventUnion()       {}
func (ForkServerEventTransportFailed) implForkServerEventUnion()         {}

// Use the following switch statement to find the correct variant
//
//	switch variant := ForkServerEventUnion.AsAny().(type) {
//	case live.SessionStartedEvent:
//	case live.SessionUpdatedEvent:
//	case live.InputAudioMutedEvent:
//	case live.InputAudioUnmutedEvent:
//	case live.OutputAudioFixturePlayingEvent:
//	case live.OutputAudioFixturePlayedEvent:
//	case live.InstructionsAppendedEvent:
//	case live.ThinkingAppendedEvent:
//	case live.CommentaryAppendedEvent:
//	case live.ForkServerEventSessionInputAudioAppend:
//	case live.OutputAudioDeltaEvent:
//	case live.InputTranscriptDeltaEvent:
//	case live.OutputTranscriptDeltaEvent:
//	case live.DelegationCreatedEvent:
//	case live.ResponseEvent:
//	case live.SessionUsageUpdatedEvent:
//	case live.SessionClosedEvent:
//	case live.ErrorEvent:
//	case live.InfoEvent:
//	case live.ForkServerEventTransportDtmfReceived:
//	case live.ForkServerEventTransportDtmfSend:
//	case live.ForkServerEventTransportRinging:
//	case live.ForkServerEventTransportAnswered:
//	case live.ForkServerEventTransportFailed:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ForkServerEventUnion) AsAny() anyForkServerEvent {
	switch u.Type {
	case "session.started":
		return u.AsSessionStarted()
	case "session.updated":
		return u.AsSessionUpdated()
	case "session.input_audio.muted":
		return u.AsSessionInputAudioMuted()
	case "session.input_audio.unmuted":
		return u.AsSessionInputAudioUnmuted()
	case "session.output_audio.fixture.playing":
		return u.AsSessionOutputAudioFixturePlaying()
	case "session.output_audio.fixture.played":
		return u.AsSessionOutputAudioFixturePlayed()
	case "session.instructions.appended":
		return u.AsSessionInstructionsAppended()
	case "session.thinking.appended":
		return u.AsSessionThinkingAppended()
	case "session.commentary.appended":
		return u.AsSessionCommentaryAppended()
	case "session.input_audio.append":
		return u.AsSessionInputAudioAppend()
	case "session.output_audio.delta":
		return u.AsSessionOutputAudioDelta()
	case "session.input_transcript.delta":
		return u.AsSessionInputTranscriptDelta()
	case "session.output_transcript.delta":
		return u.AsSessionOutputTranscriptDelta()
	case "session.delegation.created":
		return u.AsSessionDelegationCreated()
	case "response.event":
		return u.AsResponseEvent()
	case "session.usage.updated":
		return u.AsSessionUsageUpdated()
	case "session.closed":
		return u.AsSessionClosed()
	case "error":
		return u.AsError()
	case "info":
		return u.AsInfo()
	case "transport.dtmf.received":
		return u.AsTransportDtmfReceived()
	case "transport.dtmf.send":
		return u.AsTransportDtmfSend()
	case "transport.ringing":
		return u.AsTransportRinging()
	case "transport.answered":
		return u.AsTransportAnswered()
	case "transport.failed":
		return u.AsTransportFailed()
	}
	return nil
}

func (u ForkServerEventUnion) AsSessionStarted() (v SessionStartedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ForkServerEventUnion) AsSessionUpdated() (v SessionUpdatedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ForkServerEventUnion) AsSessionInputAudioMuted() (v InputAudioMutedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ForkServerEventUnion) AsSessionInputAudioUnmuted() (v InputAudioUnmutedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ForkServerEventUnion) AsSessionOutputAudioFixturePlaying() (v OutputAudioFixturePlayingEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ForkServerEventUnion) AsSessionOutputAudioFixturePlayed() (v OutputAudioFixturePlayedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ForkServerEventUnion) AsSessionInstructionsAppended() (v InstructionsAppendedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ForkServerEventUnion) AsSessionThinkingAppended() (v ThinkingAppendedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ForkServerEventUnion) AsSessionCommentaryAppended() (v CommentaryAppendedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ForkServerEventUnion) AsSessionInputAudioAppend() (v ForkServerEventSessionInputAudioAppend) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ForkServerEventUnion) AsSessionOutputAudioDelta() (v OutputAudioDeltaEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ForkServerEventUnion) AsSessionInputTranscriptDelta() (v InputTranscriptDeltaEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ForkServerEventUnion) AsSessionOutputTranscriptDelta() (v OutputTranscriptDeltaEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ForkServerEventUnion) AsSessionDelegationCreated() (v DelegationCreatedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ForkServerEventUnion) AsResponseEvent() (v ResponseEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ForkServerEventUnion) AsSessionUsageUpdated() (v SessionUsageUpdatedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ForkServerEventUnion) AsSessionClosed() (v SessionClosedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ForkServerEventUnion) AsError() (v ErrorEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ForkServerEventUnion) AsInfo() (v InfoEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ForkServerEventUnion) AsTransportDtmfReceived() (v ForkServerEventTransportDtmfReceived) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ForkServerEventUnion) AsTransportDtmfSend() (v ForkServerEventTransportDtmfSend) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ForkServerEventUnion) AsTransportRinging() (v ForkServerEventTransportRinging) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ForkServerEventUnion) AsTransportAnswered() (v ForkServerEventTransportAnswered) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ForkServerEventUnion) AsTransportFailed() (v ForkServerEventTransportFailed) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ForkServerEventUnion) RawJSON() string { return u.JSON.raw }

func (r *ForkServerEventUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// unmarshalWebsocketEvent decodes only the event selected by its type. Public
// UnmarshalJSON keeps the existing general-purpose union decoding behavior.
func (r *ForkServerEventUnion) unmarshalWebsocketEvent(data []byte) error {
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
	*r = ForkServerEventUnion{Type: eventType}
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
		if err := apijson.UnmarshalRoot(data, &r.OfForkServerEventSessionInputAudioAppend); err != nil {
			return err
		}
		r.JSON.OfForkServerEventSessionInputAudioAppend = respjson.NewField(r.JSON.raw)
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
		if err := apijson.UnmarshalRoot(data, &r.OfForkServerEventTransportDtmfReceived); err != nil {
			return err
		}
		r.JSON.OfForkServerEventTransportDtmfReceived = respjson.NewField(r.JSON.raw)
	case "transport.dtmf.send":
		if err := apijson.UnmarshalRoot(data, &r.OfForkServerEventTransportDtmfSend); err != nil {
			return err
		}
		r.JSON.OfForkServerEventTransportDtmfSend = respjson.NewField(r.JSON.raw)
	case "transport.ringing":
		if err := apijson.UnmarshalRoot(data, &r.OfForkServerEventTransportRinging); err != nil {
			return err
		}
		r.JSON.OfForkServerEventTransportRinging = respjson.NewField(r.JSON.raw)
	case "transport.answered":
		if err := apijson.UnmarshalRoot(data, &r.OfForkServerEventTransportAnswered); err != nil {
			return err
		}
		r.JSON.OfForkServerEventTransportAnswered = respjson.NewField(r.JSON.raw)
	case "transport.failed":
		if err := apijson.UnmarshalRoot(data, &r.OfForkServerEventTransportFailed); err != nil {
			return err
		}
		r.JSON.OfForkServerEventTransportFailed = respjson.NewField(r.JSON.raw)
	}
	return nil
}

// Input audio received from the primary transport and reflected to a Live sideband
// connection before model-input muting.
type ForkServerEventSessionInputAudioAppend struct {
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
func (r ForkServerEventSessionInputAudioAppend) RawJSON() string { return r.JSON.raw }
func (r *ForkServerEventSessionInputAudioAppend) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A SIP DTMF keypress received from the caller. Delivered only to sideband
// observers.
type ForkServerEventTransportDtmfReceived struct {
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
func (r ForkServerEventTransportDtmfReceived) RawJSON() string { return r.JSON.raw }
func (r *ForkServerEventTransportDtmfReceived) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A SIP DTMF keypress successfully sent to the SIP trunk. Delivered only to
// sideband observers; this is not a client command.
type ForkServerEventTransportDtmfSend struct {
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
func (r ForkServerEventTransportDtmfSend) RawJSON() string { return r.JSON.raw }
func (r *ForkServerEventTransportDtmfSend) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The outbound SIP provider leg is ringing or providing early media. Delivered
// only to sideband observers.
type ForkServerEventTransportRinging struct {
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
func (r ForkServerEventTransportRinging) RawJSON() string { return r.JSON.raw }
func (r *ForkServerEventTransportRinging) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The outbound SIP provider leg answered and media is established. Delivered only
// to sideband observers.
type ForkServerEventTransportAnswered struct {
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
func (r ForkServerEventTransportAnswered) RawJSON() string { return r.JSON.raw }
func (r *ForkServerEventTransportAnswered) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// An asynchronous outbound SIP setup failure. Delivered only to sideband
// observers.
type ForkServerEventTransportFailed struct {
	Error   ForkServerEventTransportFailedError `json:"error" api:"required"`
	EventID string                              `json:"event_id" api:"required"`
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
func (r ForkServerEventTransportFailed) RawJSON() string { return r.JSON.raw }
func (r *ForkServerEventTransportFailed) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ForkServerEventTransportFailedError struct {
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
func (r ForkServerEventTransportFailedError) RawJSON() string { return r.JSON.raw }
func (r *ForkServerEventTransportFailedError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

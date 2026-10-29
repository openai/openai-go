// File generated from our OpenAPI spec by Castiron. See CONTRIBUTING.md for details.

package live

import (
	"encoding/json"

	"github.com/openai/openai-go/v3/internal/apijson"
	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/packages/respjson"
	"github.com/openai/openai-go/v3/responses"
	"github.com/openai/openai-go/v3/shared/constant"
)

// SidebandService contains methods and other services that help with interacting
// with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewSidebandService] method instead.
type SidebandService struct {
	Options []option.RequestOption
}

// NewSidebandService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewSidebandService(opts ...option.RequestOption) (r SidebandService) {
	r = SidebandService{}
	r.Options = requestconfig.InheritedOptions(opts...)
	return
}

func ConnectClientEventParamOfSessionUpdate(session SessionUpdateConfigParam) ConnectClientEventUnionParam {
	var sessionUpdate SessionUpdateEventParam
	sessionUpdate.Session = session
	return ConnectClientEventUnionParam{OfSessionUpdate: &sessionUpdate}
}

func ConnectClientEventParamOfSessionOutputAudioFixturePlay(audio string) ConnectClientEventUnionParam {
	var sessionOutputAudioFixturePlay OutputAudioFixturePlayEventParam
	sessionOutputAudioFixturePlay.Audio = audio
	return ConnectClientEventUnionParam{OfSessionOutputAudioFixturePlay: &sessionOutputAudioFixturePlay}
}

func ConnectClientEventParamOfTransportDtmfSend(event string) ConnectClientEventUnionParam {
	var transportDtmfSend ConnectClientEventTransportDtmfSendParam
	transportDtmfSend.Event = event
	return ConnectClientEventUnionParam{OfTransportDtmfSend: &transportDtmfSend}
}

func ConnectClientEventParamOfSessionInstructionsAppend(content string, delegationID string) ConnectClientEventUnionParam {
	var sessionInstructionsAppend InstructionsAppendEventParam
	sessionInstructionsAppend.Content = content
	sessionInstructionsAppend.DelegationID = param.NewOpt(delegationID)
	return ConnectClientEventUnionParam{OfSessionInstructionsAppend: &sessionInstructionsAppend}
}

func ConnectClientEventParamOfSessionThinkingAppend(content string, delegationID string) ConnectClientEventUnionParam {
	var sessionThinkingAppend ThinkingAppendEventParam
	sessionThinkingAppend.Content = content
	sessionThinkingAppend.DelegationID = param.NewOpt(delegationID)
	return ConnectClientEventUnionParam{OfSessionThinkingAppend: &sessionThinkingAppend}
}

func ConnectClientEventParamOfSessionCommentaryAppend(content string, delegationID string) ConnectClientEventUnionParam {
	var sessionCommentaryAppend CommentaryAppendEventParam
	sessionCommentaryAppend.Content = content
	sessionCommentaryAppend.DelegationID = param.NewOpt(delegationID)
	return ConnectClientEventUnionParam{OfSessionCommentaryAppend: &sessionCommentaryAppend}
}

func ConnectClientEventParamOfDecisionCreate(eventID string) ConnectClientEventUnionParam {
	var decisionCreate DecisionCreateEventParam
	decisionCreate.EventID = eventID
	return ConnectClientEventUnionParam{OfDecisionCreate: &decisionCreate}
}

func ConnectClientEventParamOfDecisionItemCreate(item DecisionItemCreateEventItemParam) ConnectClientEventUnionParam {
	var decisionItemCreate DecisionItemCreateEventParam
	decisionItemCreate.Item = item
	return ConnectClientEventUnionParam{OfDecisionItemCreate: &decisionItemCreate}
}

func ConnectClientEventParamOfDecisionItemDelete(itemID string) ConnectClientEventUnionParam {
	var decisionItemDelete DecisionItemDeleteEventParam
	decisionItemDelete.ItemID = itemID
	return ConnectClientEventUnionParam{OfDecisionItemDelete: &decisionItemDelete}
}

func ConnectClientEventParamOfResponseItemCreate[
	T responses.EasyInputMessageParam | responses.ResponseInputItemMessageParam | responses.ResponseOutputMessageParam | responses.ResponseFileSearchToolCallParam | responses.ResponseComputerToolCallParam | responses.ResponseInputItemComputerCallOutputParam | responses.ResponseFunctionWebSearchParam | responses.ResponseFunctionToolCallParam | responses.ResponseInputItemFunctionCallOutputParam | responses.ResponseInputItemToolSearchCallParam | responses.ResponseToolSearchOutputItemParam | responses.ResponseInputItemAdditionalToolsParam | responses.ResponseConfigurationUpdateItemParam | responses.ResponseReasoningItemParam | responses.ResponseCompactionItemParam | responses.ResponseInputItemImageGenerationCallParam | responses.ResponseCodeInterpreterToolCallParam | responses.ResponseInputItemLocalShellCallParam | responses.ResponseInputItemLocalShellCallOutputParam | responses.ResponseInputItemShellCallParam | responses.ResponseInputItemShellCallOutputParam | responses.ResponseInputItemApplyPatchCallParam | responses.ResponseInputItemApplyPatchCallOutputParam | responses.ResponseInputItemMcpListToolsParam | responses.ResponseInputItemMcpApprovalRequestParam | responses.ResponseInputItemMcpApprovalResponseParam | responses.ResponseInputItemMcpCallParam | responses.ResponseCustomToolCallOutputParam | responses.ResponseCustomToolCallParam | responses.ResponseInputItemCompactionTriggerParam | responses.ResponseInputItemItemReferenceParam | responses.ResponseInputItemProgramParam | responses.ResponseInputItemProgramOutputParam,
](item T) ConnectClientEventUnionParam {
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
	return ConnectClientEventUnionParam{OfResponseItemCreate: &responseItemCreate}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ConnectClientEventUnionParam struct {
	OfSessionUpdate                 *SessionUpdateEventParam                  `json:",omitzero,inline"`
	OfSessionInputAudioMute         *InputAudioMuteEventParam                 `json:",omitzero,inline"`
	OfSessionInputAudioUnmute       *InputAudioUnmuteEventParam               `json:",omitzero,inline"`
	OfSessionOutputAudioFixturePlay *OutputAudioFixturePlayEventParam         `json:",omitzero,inline"`
	OfTransportDtmfSend             *ConnectClientEventTransportDtmfSendParam `json:",omitzero,inline"`
	OfSessionInstructionsAppend     *InstructionsAppendEventParam             `json:",omitzero,inline"`
	OfSessionThinkingAppend         *ThinkingAppendEventParam                 `json:",omitzero,inline"`
	OfSessionCommentaryAppend       *CommentaryAppendEventParam               `json:",omitzero,inline"`
	OfDecisionCreate                *DecisionCreateEventParam                 `json:",omitzero,inline"`
	OfDecisionItemCreate            *DecisionItemCreateEventParam             `json:",omitzero,inline"`
	OfDecisionItemDelete            *DecisionItemDeleteEventParam             `json:",omitzero,inline"`
	OfResponseItemCreate            *ResponseItemCreateEventParam             `json:",omitzero,inline"`
	OfResponseCreate                *ResponseCreateEventParam                 `json:",omitzero,inline"`
	OfSessionClose                  *SessionCloseEventParam                   `json:",omitzero,inline"`
	paramUnion
}

func (u ConnectClientEventUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfSessionUpdate,
		u.OfSessionInputAudioMute,
		u.OfSessionInputAudioUnmute,
		u.OfSessionOutputAudioFixturePlay,
		u.OfTransportDtmfSend,
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
func (u *ConnectClientEventUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u ConnectClientEventUnionParam) GetSession() *SessionUpdateConfigParam {
	if vt := u.OfSessionUpdate; vt != nil {
		return &vt.Session
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ConnectClientEventUnionParam) GetAudio() *string {
	if vt := u.OfSessionOutputAudioFixturePlay; vt != nil {
		return &vt.Audio
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ConnectClientEventUnionParam) GetEvent() *string {
	if vt := u.OfTransportDtmfSend; vt != nil {
		return &vt.Event
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ConnectClientEventUnionParam) GetItemID() *string {
	if vt := u.OfDecisionItemDelete; vt != nil {
		return &vt.ItemID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ConnectClientEventUnionParam) GetItem() *responses.ResponseInputItemUnionParam {
	if vt := u.OfResponseItemCreate; vt != nil {
		return &vt.Item
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ConnectClientEventUnionParam) GetType() *string {
	if vt := u.OfSessionUpdate; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfSessionInputAudioMute; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfSessionInputAudioUnmute; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfSessionOutputAudioFixturePlay; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfTransportDtmfSend; vt != nil {
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
func (u ConnectClientEventUnionParam) GetEventID() *string {
	if vt := u.OfSessionUpdate; vt != nil && vt.EventID.Valid() {
		return &vt.EventID.Value
	} else if vt := u.OfSessionInputAudioMute; vt != nil && vt.EventID.Valid() {
		return &vt.EventID.Value
	} else if vt := u.OfSessionInputAudioUnmute; vt != nil && vt.EventID.Valid() {
		return &vt.EventID.Value
	} else if vt := u.OfSessionOutputAudioFixturePlay; vt != nil && vt.EventID.Valid() {
		return &vt.EventID.Value
	} else if vt := u.OfTransportDtmfSend; vt != nil && vt.EventID.Valid() {
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
func (u ConnectClientEventUnionParam) GetContent() *string {
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
func (u ConnectClientEventUnionParam) GetDelegationID() *string {
	if vt := u.OfSessionInstructionsAppend; vt != nil && vt.DelegationID.Valid() {
		return &vt.DelegationID.Value
	} else if vt := u.OfSessionThinkingAppend; vt != nil && vt.DelegationID.Valid() {
		return &vt.DelegationID.Value
	} else if vt := u.OfSessionCommentaryAppend; vt != nil && vt.DelegationID.Valid() {
		return &vt.DelegationID.Value
	}
	return nil
}

func init() {
	apijson.RegisterUnion[ConnectClientEventUnionParam](
		"type",
		apijson.Discriminator[SessionUpdateEventParam]("session.update"),
		apijson.Discriminator[InputAudioMuteEventParam]("session.input_audio.mute"),
		apijson.Discriminator[InputAudioUnmuteEventParam]("session.input_audio.unmute"),
		apijson.Discriminator[OutputAudioFixturePlayEventParam]("session.output_audio.fixture.play"),
		apijson.Discriminator[ConnectClientEventTransportDtmfSendParam]("transport.dtmf.send"),
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

// Send one DTMF keypress to the SIP trunk through a Live sideband connection. This
// command is supported only for SIP calls. A successful send emits
// `transport.dtmf.sent` on sideband. If you provide `event_id`, the success
// notification includes it as `client_event_id`; a failure includes it as
// `error.client_event_id`. Send separate events for multiple keypresses.
//
// The properties Event, Type are required.
type ConnectClientEventTransportDtmfSendParam struct {
	// The DTMF key to send: one of `0`–`9`, `A`–`D`, `*`, or `#`.
	Event string `json:"event" api:"required"`
	// Optional client identifier for correlating this command with a server event's
	// client_event_id or error.client_event_id.
	EventID param.Opt[string] `json:"event_id,omitzero"`
	// This field can be elided, and will marshal its zero value as
	// "transport.dtmf.send".
	Type constant.TransportDtmfSend `json:"type" default:"transport.dtmf.send"`
	paramObj
}

func (r ConnectClientEventTransportDtmfSendParam) MarshalJSON() (data []byte, err error) {
	type shadow ConnectClientEventTransportDtmfSendParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ConnectClientEventTransportDtmfSendParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ConnectServerEventUnion contains all possible properties and values from
// [SessionStartedEvent], [SessionUpdatedEvent], [InputAudioMutedEvent],
// [InputAudioUnmutedEvent], [OutputAudioFixturePlayingEvent],
// [OutputAudioFixturePlayedEvent], [InstructionsAppendedEvent],
// [ThinkingAppendedEvent], [CommentaryAppendedEvent], [InputTranscriptDeltaEvent],
// [OutputTranscriptDeltaEvent], [DelegationCreatedEvent], [DecisionEvent],
// [DecisionItemCreatedEvent], [DecisionItemDeletedEvent], [ResponseEvent],
// [SessionUsageUpdatedEvent], [SessionClosedEvent], [ErrorEvent], [InfoEvent],
// [ConnectServerEventSessionInputAudioAppend], [OutputAudioDeltaEvent],
// [ConnectServerEventTransportDtmfReceived],
// [ConnectServerEventTransportDtmfSend], [ConnectServerEventTransportDtmfSent],
// [ConnectServerEventTransportRinging], [ConnectServerEventTransportAnswered],
// [ConnectServerEventTransportFailed].
//
// Use the [ConnectServerEventUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ConnectServerEventUnion struct {
	OfSessionStartedEvent                       SessionStartedEvent                       `json:",inline"`
	OfSessionUpdatedEvent                       SessionUpdatedEvent                       `json:",inline"`
	OfInputAudioMutedEvent                      InputAudioMutedEvent                      `json:",inline"`
	OfInputAudioUnmutedEvent                    InputAudioUnmutedEvent                    `json:",inline"`
	OfOutputAudioFixturePlayingEvent            OutputAudioFixturePlayingEvent            `json:",inline"`
	OfOutputAudioFixturePlayedEvent             OutputAudioFixturePlayedEvent             `json:",inline"`
	OfInstructionsAppendedEvent                 InstructionsAppendedEvent                 `json:",inline"`
	OfThinkingAppendedEvent                     ThinkingAppendedEvent                     `json:",inline"`
	OfCommentaryAppendedEvent                   CommentaryAppendedEvent                   `json:",inline"`
	OfInputTranscriptDeltaEvent                 InputTranscriptDeltaEvent                 `json:",inline"`
	OfOutputTranscriptDeltaEvent                OutputTranscriptDeltaEvent                `json:",inline"`
	OfDelegationCreatedEvent                    DelegationCreatedEvent                    `json:",inline"`
	OfDecisionEvent                             DecisionEvent                             `json:",inline"`
	OfDecisionItemCreatedEvent                  DecisionItemCreatedEvent                  `json:",inline"`
	OfDecisionItemDeletedEvent                  DecisionItemDeletedEvent                  `json:",inline"`
	OfResponseEvent                             ResponseEvent                             `json:",inline"`
	OfSessionUsageUpdatedEvent                  SessionUsageUpdatedEvent                  `json:",inline"`
	OfSessionClosedEvent                        SessionClosedEvent                        `json:",inline"`
	OfErrorEvent                                ErrorEvent                                `json:",inline"`
	OfInfoEvent                                 InfoEvent                                 `json:",inline"`
	OfConnectServerEventSessionInputAudioAppend ConnectServerEventSessionInputAudioAppend `json:",inline"`
	OfOutputAudioDeltaEvent                     OutputAudioDeltaEvent                     `json:",inline"`
	OfConnectServerEventTransportDtmfReceived   ConnectServerEventTransportDtmfReceived   `json:",inline"`
	OfConnectServerEventTransportDtmfSend       ConnectServerEventTransportDtmfSend       `json:",inline"`
	OfConnectServerEventTransportDtmfSent       ConnectServerEventTransportDtmfSent       `json:",inline"`
	OfConnectServerEventTransportRinging        ConnectServerEventTransportRinging        `json:",inline"`
	OfConnectServerEventTransportAnswered       ConnectServerEventTransportAnswered       `json:",inline"`
	OfConnectServerEventTransportFailed         ConnectServerEventTransportFailed         `json:",inline"`
	Type                                        string                                    `json:"type"`
	JSON                                        struct {
		OfSessionStartedEvent                       respjson.Field
		OfSessionUpdatedEvent                       respjson.Field
		OfInputAudioMutedEvent                      respjson.Field
		OfInputAudioUnmutedEvent                    respjson.Field
		OfOutputAudioFixturePlayingEvent            respjson.Field
		OfOutputAudioFixturePlayedEvent             respjson.Field
		OfInstructionsAppendedEvent                 respjson.Field
		OfThinkingAppendedEvent                     respjson.Field
		OfCommentaryAppendedEvent                   respjson.Field
		OfInputTranscriptDeltaEvent                 respjson.Field
		OfOutputTranscriptDeltaEvent                respjson.Field
		OfDelegationCreatedEvent                    respjson.Field
		OfDecisionEvent                             respjson.Field
		OfDecisionItemCreatedEvent                  respjson.Field
		OfDecisionItemDeletedEvent                  respjson.Field
		OfResponseEvent                             respjson.Field
		OfSessionUsageUpdatedEvent                  respjson.Field
		OfSessionClosedEvent                        respjson.Field
		OfErrorEvent                                respjson.Field
		OfInfoEvent                                 respjson.Field
		OfConnectServerEventSessionInputAudioAppend respjson.Field
		OfOutputAudioDeltaEvent                     respjson.Field
		OfConnectServerEventTransportDtmfReceived   respjson.Field
		OfConnectServerEventTransportDtmfSend       respjson.Field
		OfConnectServerEventTransportDtmfSent       respjson.Field
		OfConnectServerEventTransportRinging        respjson.Field
		OfConnectServerEventTransportAnswered       respjson.Field
		OfConnectServerEventTransportFailed         respjson.Field
		Type                                        respjson.Field
		raw                                         string
	} `json:"-"`
}

// anyConnectServerEvent is implemented by each variant of
// [ConnectServerEventUnion] to add type safety for the return type of
// [ConnectServerEventUnion.AsAny]
type anyConnectServerEvent interface {
	implConnectServerEventUnion()
}

func (SessionStartedEvent) implConnectServerEventUnion()                       {}
func (SessionUpdatedEvent) implConnectServerEventUnion()                       {}
func (InputAudioMutedEvent) implConnectServerEventUnion()                      {}
func (InputAudioUnmutedEvent) implConnectServerEventUnion()                    {}
func (OutputAudioFixturePlayingEvent) implConnectServerEventUnion()            {}
func (OutputAudioFixturePlayedEvent) implConnectServerEventUnion()             {}
func (InstructionsAppendedEvent) implConnectServerEventUnion()                 {}
func (ThinkingAppendedEvent) implConnectServerEventUnion()                     {}
func (CommentaryAppendedEvent) implConnectServerEventUnion()                   {}
func (InputTranscriptDeltaEvent) implConnectServerEventUnion()                 {}
func (OutputTranscriptDeltaEvent) implConnectServerEventUnion()                {}
func (DelegationCreatedEvent) implConnectServerEventUnion()                    {}
func (DecisionEvent) implConnectServerEventUnion()                             {}
func (DecisionItemCreatedEvent) implConnectServerEventUnion()                  {}
func (DecisionItemDeletedEvent) implConnectServerEventUnion()                  {}
func (ResponseEvent) implConnectServerEventUnion()                             {}
func (SessionUsageUpdatedEvent) implConnectServerEventUnion()                  {}
func (SessionClosedEvent) implConnectServerEventUnion()                        {}
func (ErrorEvent) implConnectServerEventUnion()                                {}
func (InfoEvent) implConnectServerEventUnion()                                 {}
func (ConnectServerEventSessionInputAudioAppend) implConnectServerEventUnion() {}
func (OutputAudioDeltaEvent) implConnectServerEventUnion()                     {}
func (ConnectServerEventTransportDtmfReceived) implConnectServerEventUnion()   {}
func (ConnectServerEventTransportDtmfSend) implConnectServerEventUnion()       {}
func (ConnectServerEventTransportDtmfSent) implConnectServerEventUnion()       {}
func (ConnectServerEventTransportRinging) implConnectServerEventUnion()        {}
func (ConnectServerEventTransportAnswered) implConnectServerEventUnion()       {}
func (ConnectServerEventTransportFailed) implConnectServerEventUnion()         {}

// Use the following switch statement to find the correct variant
//
//	switch variant := ConnectServerEventUnion.AsAny().(type) {
//	case live.SessionStartedEvent:
//	case live.SessionUpdatedEvent:
//	case live.InputAudioMutedEvent:
//	case live.InputAudioUnmutedEvent:
//	case live.OutputAudioFixturePlayingEvent:
//	case live.OutputAudioFixturePlayedEvent:
//	case live.InstructionsAppendedEvent:
//	case live.ThinkingAppendedEvent:
//	case live.CommentaryAppendedEvent:
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
//	case live.ConnectServerEventSessionInputAudioAppend:
//	case live.OutputAudioDeltaEvent:
//	case live.ConnectServerEventTransportDtmfReceived:
//	case live.ConnectServerEventTransportDtmfSend:
//	case live.ConnectServerEventTransportDtmfSent:
//	case live.ConnectServerEventTransportRinging:
//	case live.ConnectServerEventTransportAnswered:
//	case live.ConnectServerEventTransportFailed:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ConnectServerEventUnion) AsAny() anyConnectServerEvent {
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
	case "session.input_audio.append":
		return u.AsConnectServerEventSessionInputAudioAppend()
	case "session.output_audio.delta":
		return u.AsOutputAudioDeltaEvent()
	case "transport.dtmf.received":
		return u.AsConnectServerEventTransportDtmfReceived()
	case "transport.dtmf.send":
		return u.AsConnectServerEventTransportDtmfSend()
	case "transport.dtmf.sent":
		return u.AsConnectServerEventTransportDtmfSent()
	case "transport.ringing":
		return u.AsConnectServerEventTransportRinging()
	case "transport.answered":
		return u.AsConnectServerEventTransportAnswered()
	case "transport.failed":
		return u.AsConnectServerEventTransportFailed()
	}
	return nil
}

func (u ConnectServerEventUnion) AsSessionStartedEvent() (v SessionStartedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ConnectServerEventUnion) AsSessionUpdatedEvent() (v SessionUpdatedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ConnectServerEventUnion) AsInputAudioMutedEvent() (v InputAudioMutedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ConnectServerEventUnion) AsInputAudioUnmutedEvent() (v InputAudioUnmutedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ConnectServerEventUnion) AsOutputAudioFixturePlayingEvent() (v OutputAudioFixturePlayingEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ConnectServerEventUnion) AsOutputAudioFixturePlayedEvent() (v OutputAudioFixturePlayedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ConnectServerEventUnion) AsInstructionsAppendedEvent() (v InstructionsAppendedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ConnectServerEventUnion) AsThinkingAppendedEvent() (v ThinkingAppendedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ConnectServerEventUnion) AsCommentaryAppendedEvent() (v CommentaryAppendedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ConnectServerEventUnion) AsInputTranscriptDeltaEvent() (v InputTranscriptDeltaEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ConnectServerEventUnion) AsOutputTranscriptDeltaEvent() (v OutputTranscriptDeltaEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ConnectServerEventUnion) AsDelegationCreatedEvent() (v DelegationCreatedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ConnectServerEventUnion) AsDecisionEvent() (v DecisionEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ConnectServerEventUnion) AsDecisionItemCreatedEvent() (v DecisionItemCreatedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ConnectServerEventUnion) AsDecisionItemDeletedEvent() (v DecisionItemDeletedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ConnectServerEventUnion) AsResponseEvent() (v ResponseEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ConnectServerEventUnion) AsSessionUsageUpdatedEvent() (v SessionUsageUpdatedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ConnectServerEventUnion) AsSessionClosedEvent() (v SessionClosedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ConnectServerEventUnion) AsErrorEvent() (v ErrorEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ConnectServerEventUnion) AsInfoEvent() (v InfoEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ConnectServerEventUnion) AsConnectServerEventSessionInputAudioAppend() (v ConnectServerEventSessionInputAudioAppend) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ConnectServerEventUnion) AsOutputAudioDeltaEvent() (v OutputAudioDeltaEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ConnectServerEventUnion) AsConnectServerEventTransportDtmfReceived() (v ConnectServerEventTransportDtmfReceived) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ConnectServerEventUnion) AsConnectServerEventTransportDtmfSend() (v ConnectServerEventTransportDtmfSend) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ConnectServerEventUnion) AsConnectServerEventTransportDtmfSent() (v ConnectServerEventTransportDtmfSent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ConnectServerEventUnion) AsConnectServerEventTransportRinging() (v ConnectServerEventTransportRinging) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ConnectServerEventUnion) AsConnectServerEventTransportAnswered() (v ConnectServerEventTransportAnswered) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ConnectServerEventUnion) AsConnectServerEventTransportFailed() (v ConnectServerEventTransportFailed) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ConnectServerEventUnion) RawJSON() string { return u.JSON.raw }

func (r *ConnectServerEventUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// unmarshalWebsocketEventWithRawType decodes only the event selected by its type. Public
// UnmarshalJSON keeps the existing general-purpose union decoding behavior.
func (r *ConnectServerEventUnion) unmarshalWebsocketEventWithRawType(data []byte) error {
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
	*r = ConnectServerEventUnion{Type: eventType}
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
	case "session.input_audio.append":
		if err := apijson.UnmarshalRoot(data, &r.OfConnectServerEventSessionInputAudioAppend); err != nil {
			return err
		}
		r.JSON.OfConnectServerEventSessionInputAudioAppend = respjson.NewField(r.JSON.raw)
	case "session.output_audio.delta":
		if err := apijson.UnmarshalRoot(data, &r.OfOutputAudioDeltaEvent); err != nil {
			return err
		}
		r.JSON.OfOutputAudioDeltaEvent = respjson.NewField(r.JSON.raw)
	case "transport.dtmf.received":
		if err := apijson.UnmarshalRoot(data, &r.OfConnectServerEventTransportDtmfReceived); err != nil {
			return err
		}
		r.JSON.OfConnectServerEventTransportDtmfReceived = respjson.NewField(r.JSON.raw)
	case "transport.dtmf.send":
		if err := apijson.UnmarshalRoot(data, &r.OfConnectServerEventTransportDtmfSend); err != nil {
			return err
		}
		r.JSON.OfConnectServerEventTransportDtmfSend = respjson.NewField(r.JSON.raw)
	case "transport.dtmf.sent":
		if err := apijson.UnmarshalRoot(data, &r.OfConnectServerEventTransportDtmfSent); err != nil {
			return err
		}
		r.JSON.OfConnectServerEventTransportDtmfSent = respjson.NewField(r.JSON.raw)
	case "transport.ringing":
		if err := apijson.UnmarshalRoot(data, &r.OfConnectServerEventTransportRinging); err != nil {
			return err
		}
		r.JSON.OfConnectServerEventTransportRinging = respjson.NewField(r.JSON.raw)
	case "transport.answered":
		if err := apijson.UnmarshalRoot(data, &r.OfConnectServerEventTransportAnswered); err != nil {
			return err
		}
		r.JSON.OfConnectServerEventTransportAnswered = respjson.NewField(r.JSON.raw)
	case "transport.failed":
		if err := apijson.UnmarshalRoot(data, &r.OfConnectServerEventTransportFailed); err != nil {
			return err
		}
		r.JSON.OfConnectServerEventTransportFailed = respjson.NewField(r.JSON.raw)
	}
	return nil
}

// Input audio received from the primary transport and reflected to a Live sideband
// connection before model-input muting.
type ConnectServerEventSessionInputAudioAppend struct {
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
func (r ConnectServerEventSessionInputAudioAppend) RawJSON() string { return r.JSON.raw }
func (r *ConnectServerEventSessionInputAudioAppend) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A SIP DTMF keypress received from the caller. Delivered only to sideband
// observers.
type ConnectServerEventTransportDtmfReceived struct {
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
func (r ConnectServerEventTransportDtmfReceived) RawJSON() string { return r.JSON.raw }
func (r *ConnectServerEventTransportDtmfReceived) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Legacy SIP DTMF delivery notification retained for SDK compatibility and saved
// events. Current servers confirm delivery with transport.dtmf.sent;
// transport.dtmf.send is now the client command.
type ConnectServerEventTransportDtmfSend struct {
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
func (r ConnectServerEventTransportDtmfSend) RawJSON() string { return r.JSON.raw }
func (r *ConnectServerEventTransportDtmfSend) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A SIP DTMF keypress successfully sent to the SIP trunk. Delivered only to
// sideband observers; this is not a client command.
type ConnectServerEventTransportDtmfSent struct {
	// The event_id of the client command, when supplied.
	ClientEventID string                     `json:"client_event_id"`
	Event         string                     `json:"event" api:"required"`
	EventID       string                     `json:"event_id" api:"required"`
	Type          constant.TransportDtmfSent `json:"type" default:"transport.dtmf.sent"`
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
func (r ConnectServerEventTransportDtmfSent) RawJSON() string { return r.JSON.raw }
func (r *ConnectServerEventTransportDtmfSent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The outbound SIP provider leg is ringing or providing early media. Delivered
// only to sideband observers.
type ConnectServerEventTransportRinging struct {
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
func (r ConnectServerEventTransportRinging) RawJSON() string { return r.JSON.raw }
func (r *ConnectServerEventTransportRinging) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The outbound SIP provider leg answered and media is established. Delivered only
// to sideband observers.
type ConnectServerEventTransportAnswered struct {
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
func (r ConnectServerEventTransportAnswered) RawJSON() string { return r.JSON.raw }
func (r *ConnectServerEventTransportAnswered) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// An asynchronous outbound SIP setup failure. Delivered only to sideband
// observers.
type ConnectServerEventTransportFailed struct {
	Error   ConnectServerEventTransportFailedError `json:"error" api:"required"`
	EventID string                                 `json:"event_id" api:"required"`
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
func (r ConnectServerEventTransportFailed) RawJSON() string { return r.JSON.raw }
func (r *ConnectServerEventTransportFailed) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ConnectServerEventTransportFailedError struct {
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
func (r ConnectServerEventTransportFailedError) RawJSON() string { return r.JSON.raw }
func (r *ConnectServerEventTransportFailedError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

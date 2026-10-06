// File generated from our OpenAPI spec by Castiron. See CONTRIBUTING.md for details.

package live

import (
	"context"
	"net/http"
	"slices"

	"github.com/openai/openai-go/v3/internal/apijson"
	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/packages/respjson"
	"github.com/openai/openai-go/v3/shared/constant"
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
	BuiltInVoiceShimmer  BuiltInVoice = "shimmer"
	BuiltInVoiceStone    BuiltInVoice = "stone"
	BuiltInVoiceTempo    BuiltInVoice = "tempo"
	BuiltInVoiceVerse    BuiltInVoice = "verse"
	BuiltInVoiceVesper   BuiltInVoice = "vesper"
	BuiltInVoiceWillow   BuiltInVoice = "willow"
)

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

func NewClientDelegationParam() ClientDelegationParam {
	return ClientDelegationParam{
		Type: "client",
	}
}

// Delegate tasks to your application. The Live session emits delegation events
// that your backend handles.
//
// This struct has a constant value, construct it with [NewClientDelegationParam].
type ClientDelegationParam struct {
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
	paramUnion
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
	// application; use `responses` to let the API manage a Responses backend.
	Delegation MediaSessionConfigDelegationUnionParam `json:"delegation,omitzero"`
	// Startup audio configuration. WebRTC and SIP negotiate their audio format on the
	// media transport.
	Audio MediaSessionConfigAudioParam `json:"audio,omitzero"`
	// Startup-only capabilities for an untrusted frontend attached to a unified WebRTC
	// session. Trusted sideband connections are unaffected.
	Client ClientConfigParam `json:"client,omitzero"`
	// Ordered text-only history supplied before startup. Supports developer, user, and
	// assistant messages with one text part each; at most 128 messages and 8,192
	// rendered tokens in total.
	Input []InitialItemUnionParam `json:"input,omitzero"`
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
	paramUnion
}

func (u MediaSessionConfigDelegationUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfClient, u.OfResponses)
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
func (u MediaSessionConfigDelegationUnionParam) GetType() *string {
	if vt := u.OfClient; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfResponses; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[MediaSessionConfigDelegationUnionParam](
		"type",
		apijson.Discriminator[ClientDelegationParam]("client"),
		apijson.Discriminator[MediaSessionConfigDelegationResponsesParam]("responses"),
	)
}

// Delegate tasks to a Responses model managed by the Live session.
//
// The properties Responses, Type are required.
type MediaSessionConfigDelegationResponsesParam struct {
	// Backend model, prompt, and tools used when the Live session delegates a task to
	// Responses.
	Responses ResponsesDelegationConfigParam `json:"responses,omitzero" api:"required"`
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
	OfLiveToolChoiceEnum                                             param.Opt[string]                                               `json:",omitzero,inline"`
	OfResponsesDelegationConfigToolChoiceLiveFunctionToolChoiceParam *ResponsesDelegationConfigToolChoiceLiveFunctionToolChoiceParam `json:",omitzero,inline"`
	OfResponsesDelegationConfigToolChoiceLiveMcpToolChoiceParam      *ResponsesDelegationConfigToolChoiceLiveMcpToolChoiceParam      `json:",omitzero,inline"`
	paramUnion
}

func (u ResponsesDelegationConfigToolChoiceUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfLiveToolChoiceEnum, u.OfResponsesDelegationConfigToolChoiceLiveFunctionToolChoiceParam, u.OfResponsesDelegationConfigToolChoiceLiveMcpToolChoiceParam)
}
func (u *ResponsesDelegationConfigToolChoiceUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationConfigToolChoiceUnionParam) GetServerLabel() *string {
	if vt := u.OfResponsesDelegationConfigToolChoiceLiveMcpToolChoiceParam; vt != nil {
		return &vt.ServerLabel
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationConfigToolChoiceUnionParam) GetName() *string {
	if vt := u.OfResponsesDelegationConfigToolChoiceLiveFunctionToolChoiceParam; vt != nil {
		return (*string)(&vt.Name)
	} else if vt := u.OfResponsesDelegationConfigToolChoiceLiveMcpToolChoiceParam; vt != nil {
		return (*string)(&vt.Name)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationConfigToolChoiceUnionParam) GetType() *string {
	if vt := u.OfResponsesDelegationConfigToolChoiceLiveFunctionToolChoiceParam; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfResponsesDelegationConfigToolChoiceLiveMcpToolChoiceParam; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

type ResponsesDelegationConfigToolChoiceLiveToolChoiceEnum string

const (
	ResponsesDelegationConfigToolChoiceLiveToolChoiceEnumAuto     ResponsesDelegationConfigToolChoiceLiveToolChoiceEnum = "auto"
	ResponsesDelegationConfigToolChoiceLiveToolChoiceEnumNone     ResponsesDelegationConfigToolChoiceLiveToolChoiceEnum = "none"
	ResponsesDelegationConfigToolChoiceLiveToolChoiceEnumRequired ResponsesDelegationConfigToolChoiceLiveToolChoiceEnum = "required"
)

// The properties Name, Type are required.
type ResponsesDelegationConfigToolChoiceLiveFunctionToolChoiceParam struct {
	Name string `json:"name" api:"required"`
	// This field can be elided, and will marshal its zero value as "function".
	Type constant.Function `json:"type" default:"function"`
	paramObj
}

func (r ResponsesDelegationConfigToolChoiceLiveFunctionToolChoiceParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationConfigToolChoiceLiveFunctionToolChoiceParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationConfigToolChoiceLiveFunctionToolChoiceParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Name, ServerLabel, Type are required.
type ResponsesDelegationConfigToolChoiceLiveMcpToolChoiceParam struct {
	Name        string `json:"name" api:"required"`
	ServerLabel string `json:"server_label" api:"required"`
	// This field can be elided, and will marshal its zero value as "mcp".
	Type constant.Mcp `json:"type" default:"mcp"`
	paramObj
}

func (r ResponsesDelegationConfigToolChoiceLiveMcpToolChoiceParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationConfigToolChoiceLiveMcpToolChoiceParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationConfigToolChoiceLiveMcpToolChoiceParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ResponsesDelegationConfigToolUnionParam struct {
	OfFunction  *FunctionToolParam                           `json:",omitzero,inline"`
	OfWebSearch *ResponsesDelegationConfigToolWebSearchParam `json:",omitzero,inline"`
	paramUnion
}

func (u ResponsesDelegationConfigToolUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfFunction, u.OfWebSearch)
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
func (u ResponsesDelegationConfigToolUnionParam) GetType() *string {
	if vt := u.OfFunction; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfWebSearch; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[ResponsesDelegationConfigToolUnionParam](
		"type",
		apijson.Discriminator[FunctionToolParam]("function"),
		apijson.Discriminator[ResponsesDelegationConfigToolWebSearchParam]("web_search"),
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
	OfLiveToolChoiceEnum                                                   param.Opt[string]                                                     `json:",omitzero,inline"`
	OfResponsesDelegationUpdateConfigToolChoiceLiveFunctionToolChoiceParam *ResponsesDelegationUpdateConfigToolChoiceLiveFunctionToolChoiceParam `json:",omitzero,inline"`
	OfResponsesDelegationUpdateConfigToolChoiceLiveMcpToolChoiceParam      *ResponsesDelegationUpdateConfigToolChoiceLiveMcpToolChoiceParam      `json:",omitzero,inline"`
	paramUnion
}

func (u ResponsesDelegationUpdateConfigToolChoiceUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfLiveToolChoiceEnum, u.OfResponsesDelegationUpdateConfigToolChoiceLiveFunctionToolChoiceParam, u.OfResponsesDelegationUpdateConfigToolChoiceLiveMcpToolChoiceParam)
}
func (u *ResponsesDelegationUpdateConfigToolChoiceUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationUpdateConfigToolChoiceUnionParam) GetServerLabel() *string {
	if vt := u.OfResponsesDelegationUpdateConfigToolChoiceLiveMcpToolChoiceParam; vt != nil {
		return &vt.ServerLabel
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationUpdateConfigToolChoiceUnionParam) GetName() *string {
	if vt := u.OfResponsesDelegationUpdateConfigToolChoiceLiveFunctionToolChoiceParam; vt != nil {
		return (*string)(&vt.Name)
	} else if vt := u.OfResponsesDelegationUpdateConfigToolChoiceLiveMcpToolChoiceParam; vt != nil {
		return (*string)(&vt.Name)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ResponsesDelegationUpdateConfigToolChoiceUnionParam) GetType() *string {
	if vt := u.OfResponsesDelegationUpdateConfigToolChoiceLiveFunctionToolChoiceParam; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfResponsesDelegationUpdateConfigToolChoiceLiveMcpToolChoiceParam; vt != nil {
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

// The properties Name, Type are required.
type ResponsesDelegationUpdateConfigToolChoiceLiveFunctionToolChoiceParam struct {
	Name string `json:"name" api:"required"`
	// This field can be elided, and will marshal its zero value as "function".
	Type constant.Function `json:"type" default:"function"`
	paramObj
}

func (r ResponsesDelegationUpdateConfigToolChoiceLiveFunctionToolChoiceParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationUpdateConfigToolChoiceLiveFunctionToolChoiceParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationUpdateConfigToolChoiceLiveFunctionToolChoiceParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Name, ServerLabel, Type are required.
type ResponsesDelegationUpdateConfigToolChoiceLiveMcpToolChoiceParam struct {
	Name        string `json:"name" api:"required"`
	ServerLabel string `json:"server_label" api:"required"`
	// This field can be elided, and will marshal its zero value as "mcp".
	Type constant.Mcp `json:"type" default:"mcp"`
	paramObj
}

func (r ResponsesDelegationUpdateConfigToolChoiceLiveMcpToolChoiceParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponsesDelegationUpdateConfigToolChoiceLiveMcpToolChoiceParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResponsesDelegationUpdateConfigToolChoiceLiveMcpToolChoiceParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ResponsesDelegationUpdateConfigToolUnionParam struct {
	OfFunction  *FunctionToolParam                                 `json:",omitzero,inline"`
	OfWebSearch *ResponsesDelegationUpdateConfigToolWebSearchParam `json:",omitzero,inline"`
	paramUnion
}

func (u ResponsesDelegationUpdateConfigToolUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfFunction, u.OfWebSearch)
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
func (u ResponsesDelegationUpdateConfigToolUnionParam) GetType() *string {
	if vt := u.OfFunction; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfWebSearch; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[ResponsesDelegationUpdateConfigToolUnionParam](
		"type",
		apijson.Discriminator[FunctionToolParam]("function"),
		apijson.Discriminator[ResponsesDelegationUpdateConfigToolWebSearchParam]("web_search"),
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

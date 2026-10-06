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
	BuiltInVoiceBrise    BuiltInVoice = "brise"
	BuiltInVoiceCedar    BuiltInVoice = "cedar"
	BuiltInVoiceCinder   BuiltInVoice = "cinder"
	BuiltInVoiceCoral    BuiltInVoice = "coral"
	BuiltInVoiceDelta    BuiltInVoice = "delta"
	BuiltInVoiceEcho     BuiltInVoice = "echo"
	BuiltInVoiceFlitz    BuiltInVoice = "flitz"
	BuiltInVoiceGleam    BuiltInVoice = "gleam"
	BuiltInVoiceHarema   BuiltInVoice = "harema"
	BuiltInVoiceJuni     BuiltInVoice = "juni"
	BuiltInVoiceMarin    BuiltInVoice = "marin"
	BuiltInVoiceMeridian BuiltInVoice = "meridian"
	BuiltInVoiceNira     BuiltInVoice = "nira"
	BuiltInVoiceNoeul    BuiltInVoice = "noeul"
	BuiltInVoiceNuri     BuiltInVoice = "nuri"
	BuiltInVoiceQuartz   BuiltInVoice = "quartz"
	BuiltInVoiceRipple   BuiltInVoice = "ripple"
	BuiltInVoiceSage     BuiltInVoice = "sage"
	BuiltInVoiceShimmer  BuiltInVoice = "shimmer"
	BuiltInVoiceShitan   BuiltInVoice = "shitan"
	BuiltInVoiceSillage  BuiltInVoice = "sillage"
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
	OfLiveToolChoiceEnum param.Opt[string] `json:",omitzero,inline"`
	OfAnyMap             map[string]any    `json:",omitzero,inline"`
	paramUnion
}

func (u ResponsesDelegationConfigToolChoiceUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfLiveToolChoiceEnum, u.OfAnyMap)
}
func (u *ResponsesDelegationConfigToolChoiceUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

type ResponsesDelegationConfigToolChoiceLiveToolChoiceEnum string

const (
	ResponsesDelegationConfigToolChoiceLiveToolChoiceEnumAuto     ResponsesDelegationConfigToolChoiceLiveToolChoiceEnum = "auto"
	ResponsesDelegationConfigToolChoiceLiveToolChoiceEnumNone     ResponsesDelegationConfigToolChoiceLiveToolChoiceEnum = "none"
	ResponsesDelegationConfigToolChoiceLiveToolChoiceEnumRequired ResponsesDelegationConfigToolChoiceLiveToolChoiceEnum = "required"
)

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
	Source ResponsesDelegationConfigToolShellEnvironmentContainerAutoSkillInlineSourceParam `json:"source,omitzero" api:"required"`
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
	OfLiveToolChoiceEnum param.Opt[string] `json:",omitzero,inline"`
	OfAnyMap             map[string]any    `json:",omitzero,inline"`
	paramUnion
}

func (u ResponsesDelegationUpdateConfigToolChoiceUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfLiveToolChoiceEnum, u.OfAnyMap)
}
func (u *ResponsesDelegationUpdateConfigToolChoiceUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

type ResponsesDelegationUpdateConfigToolChoiceLiveToolChoiceEnum string

const (
	ResponsesDelegationUpdateConfigToolChoiceLiveToolChoiceEnumAuto     ResponsesDelegationUpdateConfigToolChoiceLiveToolChoiceEnum = "auto"
	ResponsesDelegationUpdateConfigToolChoiceLiveToolChoiceEnumNone     ResponsesDelegationUpdateConfigToolChoiceLiveToolChoiceEnum = "none"
	ResponsesDelegationUpdateConfigToolChoiceLiveToolChoiceEnumRequired ResponsesDelegationUpdateConfigToolChoiceLiveToolChoiceEnum = "required"
)

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
	Source ResponsesDelegationUpdateConfigToolShellEnvironmentContainerAutoSkillInlineSourceParam `json:"source,omitzero" api:"required"`
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

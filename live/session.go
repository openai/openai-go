// File generated from our OpenAPI spec by Castiron. See CONTRIBUTING.md for details.

package live

import (
	"context"
	"errors"
	"net/http"
	"slices"

	"github.com/openai/openai-go/v3/internal/apijson"
	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/packages/respjson"
	"github.com/openai/openai-go/v3/shared/constant"
)

// SessionService contains methods and other services that help with interacting
// with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewSessionService] method instead.
type SessionService struct {
	Options []option.RequestOption
}

// NewSessionService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewSessionService(opts ...option.RequestOption) (r SessionService) {
	r = SessionService{}
	r.Options = requestconfig.InheritedOptions(opts...)
	return
}

// Accept an incoming SIP call. Supply session with type live, the model, and
// startup configuration. Before accepting calls, follow the
// [Live prompting guide](https://developers.openai.com/api/docs/guides/live-prompting)
// to write frontend conversation instructions and a separate backend prompt. SIP
// media format is negotiated; omit audio.format.
func (r *SessionService) Accept(ctx context.Context, sessionID string, body SessionAcceptParams, opts ...option.RequestOption) (err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return err
	}
	path := requestconfig.FormatPath("live/sessions/%s/accept", sessionID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, nil, opts...)
	return err
}

// Get Live session content
func (r *SessionService) DownloadRecording(ctx context.Context, sessionID string, opts ...option.RequestOption) (res *http.Response, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "application/binary")}, opts...)
	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("live/sessions/%s/content", sessionID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Fork a stored Live session onto a new WebRTC connection.
func (r *SessionService) Fork(ctx context.Context, sessionID string, body SessionForkParams, opts ...option.RequestOption) (res *SessionForkResponse, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("live/sessions/%s/fork", sessionID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// End a SIP call identified by session_id.
func (r *SessionService) Hangup(ctx context.Context, sessionID string, opts ...option.RequestOption) (err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return err
	}
	path := requestconfig.FormatPath("live/sessions/%s/hangup", sessionID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, nil, opts...)
	return err
}

// Transfer a SIP call to another destination. Supply a nonblank target_uri for the
// SIP Refer-To header.
func (r *SessionService) Refer(ctx context.Context, sessionID string, body SessionReferParams, opts ...option.RequestOption) (err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return err
	}
	path := requestconfig.FormatPath("live/sessions/%s/refer", sessionID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, nil, opts...)
	return err
}

// Reject an incoming SIP call. Send a required SIP rejection status_code between
// 300 and 699.
func (r *SessionService) Reject(ctx context.Context, sessionID string, body SessionRejectParams, opts ...option.RequestOption) (err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return err
	}
	path := requestconfig.FormatPath("live/sessions/%s/reject", sessionID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, nil, opts...)
	return err
}

// The created Live session identifier and WebRTC answer. Apply transport.sdp as
// the peer's remote answer and wait for session.started on the data channel before
// sending commands.
type SessionForkResponse struct {
	// The newly created Live session. Use its ID for session controls and sideband
	// connections.
	Session SessionForkResponseSession `json:"session" api:"required"`
	// WebRTC transport with the SDP answer.
	Transport SessionForkResponseTransport `json:"transport" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Session     respjson.Field
		Transport   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r SessionForkResponse) RawJSON() string { return r.JSON.raw }
func (r *SessionForkResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The newly created Live session. Use its ID for session controls and sideband
// connections.
type SessionForkResponseSession struct {
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
func (r SessionForkResponseSession) RawJSON() string { return r.JSON.raw }
func (r *SessionForkResponseSession) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// WebRTC transport with the SDP answer.
type SessionForkResponseTransport struct {
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
func (r SessionForkResponseTransport) RawJSON() string { return r.JSON.raw }
func (r *SessionForkResponseTransport) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type SessionAcceptParams struct {
	// Model and startup configuration for the Live session that answers the incoming
	// SIP call.
	Session SessionAcceptParamsSession `json:"session,omitzero" api:"required"`
	paramObj
}

func (r SessionAcceptParams) MarshalJSON() (data []byte, err error) {
	type shadow SessionAcceptParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *SessionAcceptParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Model and startup configuration for the Live session that answers the incoming
// SIP call.
//
// The properties Model, Type are required.
type SessionAcceptParamsSession struct {
	// The Live model to use for the accepted call.
	Model string `json:"model" api:"required"`
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
	Delegation SessionAcceptParamsSessionDelegationUnion `json:"delegation,omitzero"`
	// Startup audio output configuration. SIP negotiates the media format;
	// audio.format is only accepted for primary WebSockets. Voice cannot change after
	// startup.
	Audio SessionAcceptParamsSessionAudio `json:"audio,omitzero"`
	// Ordered text-only history supplied before startup. Supports developer, user, and
	// assistant messages with one text part each; at most 128 messages and 8,192
	// rendered tokens in total.
	Input []InitialItemUnionParam `json:"input,omitzero"`
	// The session type. Always `live`.
	//
	// This field can be elided, and will marshal its zero value as "live".
	Type constant.Live `json:"type" default:"live"`
	paramObj
}

func (r SessionAcceptParamsSession) MarshalJSON() (data []byte, err error) {
	type shadow SessionAcceptParamsSession
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *SessionAcceptParamsSession) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The Live model. Required in the session configuration for every transport; do
// not pass it as a URL query parameter.
type SessionAcceptParamsSessionModel string

const (
	SessionAcceptParamsSessionModelGPTLive1 SessionAcceptParamsSessionModel = "gpt-live-1"
)

// Startup audio output configuration. SIP negotiates the media format;
// audio.format is only accepted for primary WebSockets. Voice cannot change after
// startup.
type SessionAcceptParamsSessionAudio struct {
	// Settings for speech generated by the Live model. Choose the voice before
	// starting the session.
	Output SessionAcceptParamsSessionAudioOutput `json:"output,omitzero"`
	paramObj
}

func (r SessionAcceptParamsSessionAudio) MarshalJSON() (data []byte, err error) {
	type shadow SessionAcceptParamsSessionAudio
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *SessionAcceptParamsSessionAudio) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Settings for speech generated by the Live model. Choose the voice before
// starting the session.
type SessionAcceptParamsSessionAudioOutput struct {
	// The voice used for Live speech, as a built-in voice name or a custom voice
	// object containing its ID. Defaults to `marin` and cannot change after startup.
	Voice SessionAcceptParamsSessionAudioOutputVoiceUnion `json:"voice,omitzero"`
	paramObj
}

func (r SessionAcceptParamsSessionAudioOutput) MarshalJSON() (data []byte, err error) {
	type shadow SessionAcceptParamsSessionAudioOutput
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *SessionAcceptParamsSessionAudioOutput) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type SessionAcceptParamsSessionAudioOutputVoiceUnion struct {
	OfString param.Opt[string] `json:",omitzero,inline"`
	// Check if union is this variant with !param.IsOmitted(union.OfBuiltIn)
	OfBuiltIn     param.Opt[BuiltInVoice] `json:",omitzero,inline"`
	OfCustomVoice *CustomVoiceParam       `json:",omitzero,inline"`
	paramUnion
}

func (u SessionAcceptParamsSessionAudioOutputVoiceUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfString, u.OfBuiltIn, u.OfCustomVoice)
}
func (u *SessionAcceptParamsSessionAudioOutputVoiceUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type SessionAcceptParamsSessionDelegationUnion struct {
	OfClient    *ClientDelegationParam                         `json:",omitzero,inline"`
	OfResponses *SessionAcceptParamsSessionDelegationResponses `json:",omitzero,inline"`
	paramUnion
}

func (u SessionAcceptParamsSessionDelegationUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfClient, u.OfResponses)
}
func (u *SessionAcceptParamsSessionDelegationUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u SessionAcceptParamsSessionDelegationUnion) GetResponses() *ResponsesDelegationConfigParam {
	if vt := u.OfResponses; vt != nil {
		return &vt.Responses
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u SessionAcceptParamsSessionDelegationUnion) GetType() *string {
	if vt := u.OfClient; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfResponses; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[SessionAcceptParamsSessionDelegationUnion](
		"type",
		apijson.Discriminator[ClientDelegationParam]("client"),
		apijson.Discriminator[SessionAcceptParamsSessionDelegationResponses]("responses"),
	)
}

// Delegate tasks to a Responses model managed by the Live session.
//
// The properties Responses, Type are required.
type SessionAcceptParamsSessionDelegationResponses struct {
	// Backend model, prompt, and tools used when the Live session delegates a task to
	// Responses.
	Responses ResponsesDelegationConfigParam `json:"responses,omitzero" api:"required"`
	// The delegation owner. Always `responses` for tasks handled by the Responses API.
	//
	// This field can be elided, and will marshal its zero value as "responses".
	Type constant.Responses `json:"type" default:"responses"`
	paramObj
}

func (r SessionAcceptParamsSessionDelegationResponses) MarshalJSON() (data []byte, err error) {
	type shadow SessionAcceptParamsSessionDelegationResponses
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *SessionAcceptParamsSessionDelegationResponses) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type SessionForkParams struct {
	// WebRTC transport with an SDP offer for the new connection to the forked session.
	Transport SessionForkParamsTransport `json:"transport,omitzero" api:"required"`
	// Optional configuration overrides for the new Live session. Omit this object or
	// send an empty object to inherit the stored session's settings.
	Session MediaSessionForkConfigParam `json:"session,omitzero"`
	paramObj
}

func (r SessionForkParams) MarshalJSON() (data []byte, err error) {
	type shadow SessionForkParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *SessionForkParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// WebRTC transport with an SDP offer for the new connection to the forked session.
//
// The properties Sdp, Type are required.
type SessionForkParamsTransport struct {
	// Session Description Protocol message for the WebRTC connection.
	Sdp string `json:"sdp" api:"required"`
	// The transport used for the Live session. Always `webrtc`.
	//
	// This field can be elided, and will marshal its zero value as "webrtc".
	Type constant.Webrtc `json:"type" default:"webrtc"`
	paramObj
}

func (r SessionForkParamsTransport) MarshalJSON() (data []byte, err error) {
	type shadow SessionForkParamsTransport
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *SessionForkParamsTransport) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type SessionReferParams struct {
	// Nonblank URI for the SIP Refer-To header, such as tel:+14155550123 or
	// sip:agent@example.com.
	TargetUri string `json:"target_uri" api:"required"`
	paramObj
}

func (r SessionReferParams) MarshalJSON() (data []byte, err error) {
	type shadow SessionReferParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *SessionReferParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type SessionRejectParams struct {
	// SIP rejection status sent to the caller. This field is required.
	StatusCode int64 `json:"status_code" api:"required"`
	paramObj
}

func (r SessionRejectParams) MarshalJSON() (data []byte, err error) {
	type shadow SessionRejectParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *SessionRejectParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

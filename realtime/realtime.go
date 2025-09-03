// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package realtime

import (
	"github.com/openai/openai-go/v2/internal/apijson"
	"github.com/openai/openai-go/v2/option"
	"github.com/openai/openai-go/v2/packages/param"
	"github.com/openai/openai-go/v2/responses"
	"github.com/openai/openai-go/v2/shared/constant"
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
}

// NewRealtimeService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewRealtimeService(opts ...option.RequestOption) (r RealtimeService) {
	r = RealtimeService{}
	r.Options = opts
	r.ClientSecrets = NewClientSecretService(opts...)
	return
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

type RealtimeAudioConfigInputParam struct {
	// The format of input audio. Options are `pcm16`, `g711_ulaw`, or `g711_alaw`. For
	// `pcm16`, input audio must be 16-bit PCM at a 24kHz sample rate, single channel
	// (mono), and little-endian byte order.
	//
	// Any of "pcm16", "g711_ulaw", "g711_alaw".
	Format string `json:"format,omitzero"`
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
	// [the /audio/transcriptions endpoint](https://platform.openai.com/docs/api-reference/audio/createTranscription)
	// and should be treated as guidance of input audio content rather than precisely
	// what the model heard. The client can optionally set the language and prompt for
	// transcription, these offer additional guidance to the transcription service.
	Transcription RealtimeAudioConfigInputTranscriptionParam `json:"transcription,omitzero"`
	// Configuration for turn detection, ether Server VAD or Semantic VAD. This can be
	// set to `null` to turn off, in which case the client must manually trigger model
	// response. Server VAD means that the model will detect the start and end of
	// speech based on audio volume and respond at the end of user speech. Semantic VAD
	// is more advanced and uses a turn detection model (in conjunction with VAD) to
	// semantically estimate whether the user has finished speaking, then dynamically
	// sets a timeout based on this probability. For example, if user audio trails off
	// with "uhhm", the model will score a low probability of turn end and wait longer
	// for the user to continue speaking. This can be useful for more natural
	// conversations, but may have a higher latency.
	TurnDetection RealtimeAudioConfigInputTurnDetectionParam `json:"turn_detection,omitzero"`
	paramObj
}

func (r RealtimeAudioConfigInputParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeAudioConfigInputParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeAudioConfigInputParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[RealtimeAudioConfigInputParam](
		"format", "pcm16", "g711_ulaw", "g711_alaw",
	)
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
	Type string `json:"type,omitzero"`
	paramObj
}

func (r RealtimeAudioConfigInputNoiseReductionParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeAudioConfigInputNoiseReductionParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeAudioConfigInputNoiseReductionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[RealtimeAudioConfigInputNoiseReductionParam](
		"type", "near_field", "far_field",
	)
}

// Configuration for input audio transcription, defaults to off and can be set to
// `null` to turn off once on. Input audio transcription is not native to the
// model, since the model consumes audio directly. Transcription runs
// asynchronously through
// [the /audio/transcriptions endpoint](https://platform.openai.com/docs/api-reference/audio/createTranscription)
// and should be treated as guidance of input audio content rather than precisely
// what the model heard. The client can optionally set the language and prompt for
// transcription, these offer additional guidance to the transcription service.
type RealtimeAudioConfigInputTranscriptionParam struct {
	// The language of the input audio. Supplying the input language in
	// [ISO-639-1](https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes) (e.g. `en`)
	// format will improve accuracy and latency.
	Language param.Opt[string] `json:"language,omitzero"`
	// An optional text to guide the model's style or continue a previous audio
	// segment. For `whisper-1`, the
	// [prompt is a list of keywords](https://platform.openai.com/docs/guides/speech-to-text#prompting).
	// For `gpt-4o-transcribe` models, the prompt is a free text string, for example
	// "expect words related to technology".
	Prompt param.Opt[string] `json:"prompt,omitzero"`
	// The model to use for transcription. Current options are `whisper-1`,
	// `gpt-4o-transcribe-latest`, `gpt-4o-mini-transcribe`, `gpt-4o-transcribe`, and
	// `gpt-4o-transcribe-diarize`.
	//
	// Any of "whisper-1", "gpt-4o-transcribe-latest", "gpt-4o-mini-transcribe",
	// "gpt-4o-transcribe", "gpt-4o-transcribe-diarize".
	Model string `json:"model,omitzero"`
	paramObj
}

func (r RealtimeAudioConfigInputTranscriptionParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeAudioConfigInputTranscriptionParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeAudioConfigInputTranscriptionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[RealtimeAudioConfigInputTranscriptionParam](
		"model", "whisper-1", "gpt-4o-transcribe-latest", "gpt-4o-mini-transcribe", "gpt-4o-transcribe", "gpt-4o-transcribe-diarize",
	)
}

// Configuration for turn detection, ether Server VAD or Semantic VAD. This can be
// set to `null` to turn off, in which case the client must manually trigger model
// response. Server VAD means that the model will detect the start and end of
// speech based on audio volume and respond at the end of user speech. Semantic VAD
// is more advanced and uses a turn detection model (in conjunction with VAD) to
// semantically estimate whether the user has finished speaking, then dynamically
// sets a timeout based on this probability. For example, if user audio trails off
// with "uhhm", the model will score a low probability of turn end and wait longer
// for the user to continue speaking. This can be useful for more natural
// conversations, but may have a higher latency.
type RealtimeAudioConfigInputTurnDetectionParam struct {
	// Optional idle timeout after which turn detection will auto-timeout when no
	// additional audio is received.
	IdleTimeoutMs param.Opt[int64] `json:"idle_timeout_ms,omitzero"`
	// Whether or not to automatically generate a response when a VAD stop event
	// occurs.
	CreateResponse param.Opt[bool] `json:"create_response,omitzero"`
	// Whether or not to automatically interrupt any ongoing response with output to
	// the default conversation (i.e. `conversation` of `auto`) when a VAD start event
	// occurs.
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
	// Used only for `semantic_vad` mode. The eagerness of the model to respond. `low`
	// will wait longer for the user to continue speaking, `high` will respond more
	// quickly. `auto` is the default and is equivalent to `medium`.
	//
	// Any of "low", "medium", "high", "auto".
	Eagerness string `json:"eagerness,omitzero"`
	// Type of turn detection.
	//
	// Any of "server_vad", "semantic_vad".
	Type string `json:"type,omitzero"`
	paramObj
}

func (r RealtimeAudioConfigInputTurnDetectionParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeAudioConfigInputTurnDetectionParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeAudioConfigInputTurnDetectionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[RealtimeAudioConfigInputTurnDetectionParam](
		"eagerness", "low", "medium", "high", "auto",
	)
	apijson.RegisterFieldValidator[RealtimeAudioConfigInputTurnDetectionParam](
		"type", "server_vad", "semantic_vad",
	)
}

type RealtimeAudioConfigOutputParam struct {
	// The speed of the model's spoken response. 1.0 is the default speed. 0.25 is the
	// minimum speed. 1.5 is the maximum speed. This value can only be changed in
	// between model turns, not while a response is in progress.
	Speed param.Opt[float64] `json:"speed,omitzero"`
	// The format of output audio. Options are `pcm16`, `g711_ulaw`, or `g711_alaw`.
	// For `pcm16`, output audio is sampled at a rate of 24kHz.
	//
	// Any of "pcm16", "g711_ulaw", "g711_alaw".
	Format string `json:"format,omitzero"`
	// The voice the model uses to respond. Voice cannot be changed during the session
	// once the model has responded with audio at least once. Current voice options are
	// `alloy`, `ash`, `ballad`, `coral`, `echo`, `sage`, `shimmer`, `verse`, `marin`,
	// and `cedar`.
	Voice string `json:"voice,omitzero"`
	paramObj
}

func (r RealtimeAudioConfigOutputParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeAudioConfigOutputParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeAudioConfigOutputParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[RealtimeAudioConfigOutputParam](
		"format", "pcm16", "g711_ulaw", "g711_alaw",
	)
}

// Configuration options for the generated client secret.
type RealtimeClientSecretConfigParam struct {
	// Configuration for the ephemeral token expiration.
	ExpiresAfter RealtimeClientSecretConfigExpiresAfterParam `json:"expires_after,omitzero"`
	paramObj
}

func (r RealtimeClientSecretConfigParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeClientSecretConfigParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeClientSecretConfigParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for the ephemeral token expiration.
//
// The property Anchor is required.
type RealtimeClientSecretConfigExpiresAfterParam struct {
	// The anchor point for the ephemeral token expiration. Only `created_at` is
	// currently supported.
	//
	// Any of "created_at".
	Anchor string `json:"anchor,omitzero,required"`
	// The number of seconds from the anchor point to the expiration. Select a value
	// between `10` and `7200`.
	Seconds param.Opt[int64] `json:"seconds,omitzero"`
	paramObj
}

func (r RealtimeClientSecretConfigExpiresAfterParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeClientSecretConfigExpiresAfterParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeClientSecretConfigExpiresAfterParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[RealtimeClientSecretConfigExpiresAfterParam](
		"anchor", "created_at",
	)
}

// Realtime session object configuration.
//
// The properties Model, Type are required.
type RealtimeSessionCreateRequestParam struct {
	// The Realtime model used for this session.
	Model RealtimeSessionCreateRequestModel `json:"model,omitzero,required"`
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
	// Sampling temperature for the model, limited to [0.6, 1.2]. For audio models a
	// temperature of 0.8 is highly recommended for best performance.
	Temperature param.Opt[float64] `json:"temperature,omitzero"`
	// Reference to a prompt template and its variables.
	// [Learn more](https://platform.openai.com/docs/guides/text?api-mode=responses#reusable-prompts).
	Prompt responses.ResponsePromptParam `json:"prompt,omitzero"`
	// Configuration options for tracing. Set to null to disable tracing. Once tracing
	// is enabled for a session, the configuration cannot be modified.
	//
	// `auto` will create a trace for the session with default values for the workflow
	// name, group id, and metadata.
	Tracing RealtimeTracingConfigUnionParam `json:"tracing,omitzero"`
	// Configuration for input and output audio.
	Audio RealtimeAudioConfigParam `json:"audio,omitzero"`
	// Configuration options for the generated client secret.
	ClientSecret RealtimeClientSecretConfigParam `json:"client_secret,omitzero"`
	// Additional fields to include in server outputs.
	//
	//   - `item.input_audio_transcription.logprobs`: Include logprobs for input audio
	//     transcription.
	//
	// Any of "item.input_audio_transcription.logprobs".
	Include []string `json:"include,omitzero"`
	// Maximum number of output tokens for a single assistant response, inclusive of
	// tool calls. Provide an integer between 1 and 4096 to limit output tokens, or
	// `inf` for the maximum available tokens for a given model. Defaults to `inf`.
	MaxOutputTokens RealtimeSessionCreateRequestMaxOutputTokensUnionParam `json:"max_output_tokens,omitzero"`
	// The set of modalities the model can respond with. To disable audio, set this to
	// ["text"].
	//
	// Any of "text", "audio".
	OutputModalities []string `json:"output_modalities,omitzero"`
	// How the model chooses tools. Provide one of the string modes or force a specific
	// function/MCP tool.
	ToolChoice RealtimeToolChoiceConfigUnionParam `json:"tool_choice,omitzero"`
	// Tools available to the model.
	Tools RealtimeToolsConfigParam `json:"tools,omitzero"`
	// Controls how the realtime conversation is truncated prior to model inference.
	// The default is `auto`. When set to `retention_ratio`, the server retains a
	// fraction of the conversation tokens prior to the instructions.
	Truncation RealtimeTruncationUnionParam `json:"truncation,omitzero"`
	// The type of session to create. Always `realtime` for the Realtime API.
	//
	// This field can be elided, and will marshal its zero value as "realtime".
	Type constant.Realtime `json:"type,required"`
	paramObj
}

func (r RealtimeSessionCreateRequestParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeSessionCreateRequestParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeSessionCreateRequestParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The Realtime model used for this session.
type RealtimeSessionCreateRequestModel = string

const (
	RealtimeSessionCreateRequestModelGPTRealtime                        RealtimeSessionCreateRequestModel = "gpt-realtime"
	RealtimeSessionCreateRequestModelGPTRealtime2025_08_28              RealtimeSessionCreateRequestModel = "gpt-realtime-2025-08-28"
	RealtimeSessionCreateRequestModelGPT4oRealtime                      RealtimeSessionCreateRequestModel = "gpt-4o-realtime"
	RealtimeSessionCreateRequestModelGPT4oMiniRealtime                  RealtimeSessionCreateRequestModel = "gpt-4o-mini-realtime"
	RealtimeSessionCreateRequestModelGPT4oRealtimePreview               RealtimeSessionCreateRequestModel = "gpt-4o-realtime-preview"
	RealtimeSessionCreateRequestModelGPT4oRealtimePreview2024_10_01     RealtimeSessionCreateRequestModel = "gpt-4o-realtime-preview-2024-10-01"
	RealtimeSessionCreateRequestModelGPT4oRealtimePreview2024_12_17     RealtimeSessionCreateRequestModel = "gpt-4o-realtime-preview-2024-12-17"
	RealtimeSessionCreateRequestModelGPT4oRealtimePreview2025_06_03     RealtimeSessionCreateRequestModel = "gpt-4o-realtime-preview-2025-06-03"
	RealtimeSessionCreateRequestModelGPT4oMiniRealtimePreview           RealtimeSessionCreateRequestModel = "gpt-4o-mini-realtime-preview"
	RealtimeSessionCreateRequestModelGPT4oMiniRealtimePreview2024_12_17 RealtimeSessionCreateRequestModel = "gpt-4o-mini-realtime-preview-2024-12-17"
)

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

func (u *RealtimeSessionCreateRequestMaxOutputTokensUnionParam) asAny() any {
	if !param.IsOmitted(u.OfInt) {
		return &u.OfInt.Value
	} else if !param.IsOmitted(u.OfInf) {
		return &u.OfInf
	}
	return nil
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

func (u *RealtimeToolChoiceConfigUnionParam) asAny() any {
	if !param.IsOmitted(u.OfToolChoiceMode) {
		return &u.OfToolChoiceMode
	} else if !param.IsOmitted(u.OfFunctionTool) {
		return u.OfFunctionTool
	} else if !param.IsOmitted(u.OfMcpTool) {
		return u.OfMcpTool
	}
	return nil
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

type RealtimeToolsConfigParam []RealtimeToolsConfigUnionParam

func RealtimeToolsConfigUnionParamOfMcp(serverLabel string) RealtimeToolsConfigUnionParam {
	var mcp RealtimeToolsConfigUnionMcpParam
	mcp.ServerLabel = serverLabel
	return RealtimeToolsConfigUnionParam{OfMcp: &mcp}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type RealtimeToolsConfigUnionParam struct {
	OfFunction *RealtimeToolsConfigUnionFunctionParam `json:",omitzero,inline"`
	OfMcp      *RealtimeToolsConfigUnionMcpParam      `json:",omitzero,inline"`
	paramUnion
}

func (u RealtimeToolsConfigUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfFunction, u.OfMcp)
}
func (u *RealtimeToolsConfigUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *RealtimeToolsConfigUnionParam) asAny() any {
	if !param.IsOmitted(u.OfFunction) {
		return u.OfFunction
	} else if !param.IsOmitted(u.OfMcp) {
		return u.OfMcp
	}
	return nil
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
		apijson.Discriminator[RealtimeToolsConfigUnionFunctionParam]("function"),
		apijson.Discriminator[RealtimeToolsConfigUnionMcpParam]("mcp"),
	)
}

type RealtimeToolsConfigUnionFunctionParam struct {
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
	Type string `json:"type,omitzero"`
	paramObj
}

func (r RealtimeToolsConfigUnionFunctionParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeToolsConfigUnionFunctionParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeToolsConfigUnionFunctionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[RealtimeToolsConfigUnionFunctionParam](
		"type", "function",
	)
}

// Give the model access to additional tools via remote Model Context Protocol
// (MCP) servers.
// [Learn more about MCP](https://platform.openai.com/docs/guides/tools-remote-mcp).
//
// The properties ServerLabel, Type are required.
type RealtimeToolsConfigUnionMcpParam struct {
	// A label for this MCP server, used to identify it in tool calls.
	ServerLabel string `json:"server_label,required"`
	// An OAuth access token that can be used with a remote MCP server, either with a
	// custom MCP server URL or a service connector. Your application must handle the
	// OAuth authorization flow and provide the token here.
	Authorization param.Opt[string] `json:"authorization,omitzero"`
	// Optional description of the MCP server, used to provide more context.
	ServerDescription param.Opt[string] `json:"server_description,omitzero"`
	// The URL for the MCP server. One of `server_url` or `connector_id` must be
	// provided.
	ServerURL param.Opt[string] `json:"server_url,omitzero"`
	// List of allowed tool names or a filter object.
	AllowedTools RealtimeToolsConfigUnionMcpAllowedToolsParam `json:"allowed_tools,omitzero"`
	// Optional HTTP headers to send to the MCP server. Use for authentication or other
	// purposes.
	Headers map[string]string `json:"headers,omitzero"`
	// Specify which of the MCP server's tools require approval.
	RequireApproval RealtimeToolsConfigUnionMcpRequireApprovalParam `json:"require_approval,omitzero"`
	// Identifier for service connectors, like those available in ChatGPT. One of
	// `server_url` or `connector_id` must be provided. Learn more about service
	// connectors
	// [here](https://platform.openai.com/docs/guides/tools-remote-mcp#connectors).
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
	ConnectorID string `json:"connector_id,omitzero"`
	// The type of the MCP tool. Always `mcp`.
	//
	// This field can be elided, and will marshal its zero value as "mcp".
	Type constant.Mcp `json:"type,required"`
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

func (u *RealtimeToolsConfigUnionMcpAllowedToolsParam) asAny() any {
	if !param.IsOmitted(u.OfMcpAllowedTools) {
		return &u.OfMcpAllowedTools
	} else if !param.IsOmitted(u.OfMcpToolFilter) {
		return u.OfMcpToolFilter
	}
	return nil
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

func (u *RealtimeToolsConfigUnionMcpRequireApprovalParam) asAny() any {
	if !param.IsOmitted(u.OfMcpToolApprovalFilter) {
		return u.OfMcpToolApprovalFilter
	} else if !param.IsOmitted(u.OfMcpToolApprovalSetting) {
		return &u.OfMcpToolApprovalSetting
	}
	return nil
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

// Specify a single approval policy for all tools. One of `always` or `never`. When
// set to `always`, all tools will require approval. When set to `never`, all tools
// will not require approval.
type RealtimeToolsConfigUnionMcpRequireApprovalMcpToolApprovalSetting string

const (
	RealtimeToolsConfigUnionMcpRequireApprovalMcpToolApprovalSettingAlways RealtimeToolsConfigUnionMcpRequireApprovalMcpToolApprovalSetting = "always"
	RealtimeToolsConfigUnionMcpRequireApprovalMcpToolApprovalSettingNever  RealtimeToolsConfigUnionMcpRequireApprovalMcpToolApprovalSetting = "never"
)

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

func (u *RealtimeTracingConfigUnionParam) asAny() any {
	if !param.IsOmitted(u.OfAuto) {
		return &u.OfAuto
	} else if !param.IsOmitted(u.OfTracingConfiguration) {
		return u.OfTracingConfiguration
	}
	return nil
}

// Granular configuration for tracing.
type RealtimeTracingConfigTracingConfigurationParam struct {
	// The group id to attach to this trace to enable filtering and grouping in the
	// traces dashboard.
	GroupID param.Opt[string] `json:"group_id,omitzero"`
	// The name of the workflow to attach to this trace. This is used to name the trace
	// in the traces dashboard.
	WorkflowName param.Opt[string] `json:"workflow_name,omitzero"`
	// The arbitrary metadata to attach to this trace to enable filtering in the traces
	// dashboard.
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

// Realtime transcription session object configuration.
//
// The properties Model, Type are required.
type RealtimeTranscriptionSessionCreateRequestParam struct {
	// ID of the model to use. The options are `gpt-4o-transcribe`,
	// `gpt-4o-mini-transcribe`, and `whisper-1` (which is powered by our open source
	// Whisper V2 model).
	Model RealtimeTranscriptionSessionCreateRequestModel `json:"model,omitzero,required"`
	// The set of items to include in the transcription. Current available items are:
	//
	// - `item.input_audio_transcription.logprobs`
	//
	// Any of "item.input_audio_transcription.logprobs".
	Include []string `json:"include,omitzero"`
	// The format of input audio. Options are `pcm16`, `g711_ulaw`, or `g711_alaw`. For
	// `pcm16`, input audio must be 16-bit PCM at a 24kHz sample rate, single channel
	// (mono), and little-endian byte order.
	//
	// Any of "pcm16", "g711_ulaw", "g711_alaw".
	InputAudioFormat RealtimeTranscriptionSessionCreateRequestInputAudioFormat `json:"input_audio_format,omitzero"`
	// Configuration for input audio noise reduction. This can be set to `null` to turn
	// off. Noise reduction filters audio added to the input audio buffer before it is
	// sent to VAD and the model. Filtering the audio can improve VAD and turn
	// detection accuracy (reducing false positives) and model performance by improving
	// perception of the input audio.
	InputAudioNoiseReduction RealtimeTranscriptionSessionCreateRequestInputAudioNoiseReductionParam `json:"input_audio_noise_reduction,omitzero"`
	// Configuration for input audio transcription. The client can optionally set the
	// language and prompt for transcription, these offer additional guidance to the
	// transcription service.
	InputAudioTranscription RealtimeTranscriptionSessionCreateRequestInputAudioTranscriptionParam `json:"input_audio_transcription,omitzero"`
	// Configuration for turn detection. Can be set to `null` to turn off. Server VAD
	// means that the model will detect the start and end of speech based on audio
	// volume and respond at the end of user speech.
	TurnDetection RealtimeTranscriptionSessionCreateRequestTurnDetectionParam `json:"turn_detection,omitzero"`
	// The type of session to create. Always `transcription` for transcription
	// sessions.
	//
	// This field can be elided, and will marshal its zero value as "transcription".
	Type constant.Transcription `json:"type,required"`
	paramObj
}

func (r RealtimeTranscriptionSessionCreateRequestParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeTranscriptionSessionCreateRequestParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeTranscriptionSessionCreateRequestParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ID of the model to use. The options are `gpt-4o-transcribe`,
// `gpt-4o-mini-transcribe`, and `whisper-1` (which is powered by our open source
// Whisper V2 model).
type RealtimeTranscriptionSessionCreateRequestModel = string

const (
	RealtimeTranscriptionSessionCreateRequestModelWhisper1            RealtimeTranscriptionSessionCreateRequestModel = "whisper-1"
	RealtimeTranscriptionSessionCreateRequestModelGPT4oTranscribe     RealtimeTranscriptionSessionCreateRequestModel = "gpt-4o-transcribe"
	RealtimeTranscriptionSessionCreateRequestModelGPT4oMiniTranscribe RealtimeTranscriptionSessionCreateRequestModel = "gpt-4o-mini-transcribe"
)

// The format of input audio. Options are `pcm16`, `g711_ulaw`, or `g711_alaw`. For
// `pcm16`, input audio must be 16-bit PCM at a 24kHz sample rate, single channel
// (mono), and little-endian byte order.
type RealtimeTranscriptionSessionCreateRequestInputAudioFormat string

const (
	RealtimeTranscriptionSessionCreateRequestInputAudioFormatPcm16    RealtimeTranscriptionSessionCreateRequestInputAudioFormat = "pcm16"
	RealtimeTranscriptionSessionCreateRequestInputAudioFormatG711Ulaw RealtimeTranscriptionSessionCreateRequestInputAudioFormat = "g711_ulaw"
	RealtimeTranscriptionSessionCreateRequestInputAudioFormatG711Alaw RealtimeTranscriptionSessionCreateRequestInputAudioFormat = "g711_alaw"
)

// Configuration for input audio noise reduction. This can be set to `null` to turn
// off. Noise reduction filters audio added to the input audio buffer before it is
// sent to VAD and the model. Filtering the audio can improve VAD and turn
// detection accuracy (reducing false positives) and model performance by improving
// perception of the input audio.
type RealtimeTranscriptionSessionCreateRequestInputAudioNoiseReductionParam struct {
	// Type of noise reduction. `near_field` is for close-talking microphones such as
	// headphones, `far_field` is for far-field microphones such as laptop or
	// conference room microphones.
	//
	// Any of "near_field", "far_field".
	Type string `json:"type,omitzero"`
	paramObj
}

func (r RealtimeTranscriptionSessionCreateRequestInputAudioNoiseReductionParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeTranscriptionSessionCreateRequestInputAudioNoiseReductionParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeTranscriptionSessionCreateRequestInputAudioNoiseReductionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[RealtimeTranscriptionSessionCreateRequestInputAudioNoiseReductionParam](
		"type", "near_field", "far_field",
	)
}

// Configuration for input audio transcription. The client can optionally set the
// language and prompt for transcription, these offer additional guidance to the
// transcription service.
type RealtimeTranscriptionSessionCreateRequestInputAudioTranscriptionParam struct {
	// The language of the input audio. Supplying the input language in
	// [ISO-639-1](https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes) (e.g. `en`)
	// format will improve accuracy and latency.
	Language param.Opt[string] `json:"language,omitzero"`
	// An optional text to guide the model's style or continue a previous audio
	// segment. For `whisper-1`, the
	// [prompt is a list of keywords](https://platform.openai.com/docs/guides/speech-to-text#prompting).
	// For `gpt-4o-transcribe` models, the prompt is a free text string, for example
	// "expect words related to technology".
	Prompt param.Opt[string] `json:"prompt,omitzero"`
	// The model to use for transcription, current options are `gpt-4o-transcribe`,
	// `gpt-4o-mini-transcribe`, and `whisper-1`.
	//
	// Any of "gpt-4o-transcribe", "gpt-4o-mini-transcribe", "whisper-1".
	Model string `json:"model,omitzero"`
	paramObj
}

func (r RealtimeTranscriptionSessionCreateRequestInputAudioTranscriptionParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeTranscriptionSessionCreateRequestInputAudioTranscriptionParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeTranscriptionSessionCreateRequestInputAudioTranscriptionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[RealtimeTranscriptionSessionCreateRequestInputAudioTranscriptionParam](
		"model", "gpt-4o-transcribe", "gpt-4o-mini-transcribe", "whisper-1",
	)
}

// Configuration for turn detection. Can be set to `null` to turn off. Server VAD
// means that the model will detect the start and end of speech based on audio
// volume and respond at the end of user speech.
type RealtimeTranscriptionSessionCreateRequestTurnDetectionParam struct {
	// Amount of audio to include before the VAD detected speech (in milliseconds).
	// Defaults to 300ms.
	PrefixPaddingMs param.Opt[int64] `json:"prefix_padding_ms,omitzero"`
	// Duration of silence to detect speech stop (in milliseconds). Defaults to 500ms.
	// With shorter values the model will respond more quickly, but may jump in on
	// short pauses from the user.
	SilenceDurationMs param.Opt[int64] `json:"silence_duration_ms,omitzero"`
	// Activation threshold for VAD (0.0 to 1.0), this defaults to 0.5. A higher
	// threshold will require louder audio to activate the model, and thus might
	// perform better in noisy environments.
	Threshold param.Opt[float64] `json:"threshold,omitzero"`
	// Type of turn detection. Only `server_vad` is currently supported for
	// transcription sessions.
	//
	// Any of "server_vad".
	Type string `json:"type,omitzero"`
	paramObj
}

func (r RealtimeTranscriptionSessionCreateRequestTurnDetectionParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeTranscriptionSessionCreateRequestTurnDetectionParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeTranscriptionSessionCreateRequestTurnDetectionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[RealtimeTranscriptionSessionCreateRequestTurnDetectionParam](
		"type", "server_vad",
	)
}

func RealtimeTruncationParamOfRetentionRatioTruncation(retentionRatio float64) RealtimeTruncationUnionParam {
	var variant RealtimeTruncationRetentionRatioTruncationParam
	variant.RetentionRatio = retentionRatio
	return RealtimeTruncationUnionParam{OfRetentionRatioTruncation: &variant}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type RealtimeTruncationUnionParam struct {
	// Check if union is this variant with
	// !param.IsOmitted(union.OfRealtimeTruncationStrategy)
	OfRealtimeTruncationStrategy param.Opt[string]                                `json:",omitzero,inline"`
	OfRetentionRatioTruncation   *RealtimeTruncationRetentionRatioTruncationParam `json:",omitzero,inline"`
	paramUnion
}

func (u RealtimeTruncationUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfRealtimeTruncationStrategy, u.OfRetentionRatioTruncation)
}
func (u *RealtimeTruncationUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *RealtimeTruncationUnionParam) asAny() any {
	if !param.IsOmitted(u.OfRealtimeTruncationStrategy) {
		return &u.OfRealtimeTruncationStrategy
	} else if !param.IsOmitted(u.OfRetentionRatioTruncation) {
		return u.OfRetentionRatioTruncation
	}
	return nil
}

// The truncation strategy to use for the session.
type RealtimeTruncationRealtimeTruncationStrategy string

const (
	RealtimeTruncationRealtimeTruncationStrategyAuto     RealtimeTruncationRealtimeTruncationStrategy = "auto"
	RealtimeTruncationRealtimeTruncationStrategyDisabled RealtimeTruncationRealtimeTruncationStrategy = "disabled"
)

// Retain a fraction of the conversation tokens.
//
// The properties RetentionRatio, Type are required.
type RealtimeTruncationRetentionRatioTruncationParam struct {
	// Fraction of pre-instruction conversation tokens to retain (0.0 - 1.0).
	RetentionRatio float64 `json:"retention_ratio,required"`
	// Optional cap on tokens allowed after the instructions.
	PostInstructionsTokenLimit param.Opt[int64] `json:"post_instructions_token_limit,omitzero"`
	// Use retention ratio truncation.
	//
	// This field can be elided, and will marshal its zero value as "retention_ratio".
	Type constant.RetentionRatio `json:"type,required"`
	paramObj
}

func (r RealtimeTruncationRetentionRatioTruncationParam) MarshalJSON() (data []byte, err error) {
	type shadow RealtimeTruncationRetentionRatioTruncationParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RealtimeTruncationRetentionRatioTruncationParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

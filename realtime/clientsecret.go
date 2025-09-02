// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package realtime

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/openai/openai-go/v2/internal/apijson"
	"github.com/openai/openai-go/v2/internal/requestconfig"
	"github.com/openai/openai-go/v2/option"
	"github.com/openai/openai-go/v2/packages/param"
	"github.com/openai/openai-go/v2/packages/respjson"
	"github.com/openai/openai-go/v2/responses"
	"github.com/openai/openai-go/v2/shared/constant"
)

// ClientSecretService contains methods and other services that help with
// interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewClientSecretService] method instead.
type ClientSecretService struct {
	Options []option.RequestOption
}

// NewClientSecretService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewClientSecretService(opts ...option.RequestOption) (r ClientSecretService) {
	r = ClientSecretService{}
	r.Options = opts
	return
}

// Create a Realtime session and client secret for either realtime or
// transcription.
func (r *ClientSecretService) New(ctx context.Context, body ClientSecretNewParams, opts ...option.RequestOption) (res *ClientSecretNewResponse, err error) {
	opts = append(r.Options[:], opts...)
	path := "realtime/client_secrets"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// A Realtime session configuration object.
type RealtimeSessionCreateResponse struct {
	// Unique identifier for the session that looks like `sess_1234567890abcdef`.
	ID string `json:"id"`
	// Configuration for input and output audio for the session.
	Audio RealtimeSessionCreateResponseAudio `json:"audio"`
	// Expiration timestamp for the session, in seconds since epoch.
	ExpiresAt int64 `json:"expires_at"`
	// Additional fields to include in server outputs.
	//
	//   - `item.input_audio_transcription.logprobs`: Include logprobs for input audio
	//     transcription.
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
	MaxOutputTokens RealtimeSessionCreateResponseMaxOutputTokensUnion `json:"max_output_tokens"`
	// The Realtime model used for this session.
	Model string `json:"model"`
	// The object type. Always `realtime.session`.
	Object string `json:"object"`
	// The set of modalities the model can respond with. To disable audio, set this to
	// ["text"].
	//
	// Any of "text", "audio".
	OutputModalities []string `json:"output_modalities"`
	// How the model chooses tools. Options are `auto`, `none`, `required`, or specify
	// a function.
	ToolChoice string `json:"tool_choice"`
	// Tools (functions) available to the model.
	Tools []RealtimeSessionCreateResponseTool `json:"tools"`
	// Configuration options for tracing. Set to null to disable tracing. Once tracing
	// is enabled for a session, the configuration cannot be modified.
	//
	// `auto` will create a trace for the session with default values for the workflow
	// name, group id, and metadata.
	Tracing RealtimeSessionCreateResponseTracingUnion `json:"tracing"`
	// Configuration for turn detection. Can be set to `null` to turn off. Server VAD
	// means that the model will detect the start and end of speech based on audio
	// volume and respond at the end of user speech.
	TurnDetection RealtimeSessionCreateResponseTurnDetection `json:"turn_detection"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID               respjson.Field
		Audio            respjson.Field
		ExpiresAt        respjson.Field
		Include          respjson.Field
		Instructions     respjson.Field
		MaxOutputTokens  respjson.Field
		Model            respjson.Field
		Object           respjson.Field
		OutputModalities respjson.Field
		ToolChoice       respjson.Field
		Tools            respjson.Field
		Tracing          respjson.Field
		TurnDetection    respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeSessionCreateResponse) RawJSON() string { return r.JSON.raw }
func (r *RealtimeSessionCreateResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for input and output audio for the session.
type RealtimeSessionCreateResponseAudio struct {
	Input  RealtimeSessionCreateResponseAudioInput  `json:"input"`
	Output RealtimeSessionCreateResponseAudioOutput `json:"output"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Input       respjson.Field
		Output      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeSessionCreateResponseAudio) RawJSON() string { return r.JSON.raw }
func (r *RealtimeSessionCreateResponseAudio) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type RealtimeSessionCreateResponseAudioInput struct {
	// The format of input audio. Options are `pcm16`, `g711_ulaw`, or `g711_alaw`.
	Format string `json:"format"`
	// Configuration for input audio noise reduction.
	NoiseReduction RealtimeSessionCreateResponseAudioInputNoiseReduction `json:"noise_reduction"`
	// Configuration for input audio transcription.
	Transcription RealtimeSessionCreateResponseAudioInputTranscription `json:"transcription"`
	// Configuration for turn detection.
	TurnDetection RealtimeSessionCreateResponseAudioInputTurnDetection `json:"turn_detection"`
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
func (r RealtimeSessionCreateResponseAudioInput) RawJSON() string { return r.JSON.raw }
func (r *RealtimeSessionCreateResponseAudioInput) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for input audio noise reduction.
type RealtimeSessionCreateResponseAudioInputNoiseReduction struct {
	// Any of "near_field", "far_field".
	Type string `json:"type"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeSessionCreateResponseAudioInputNoiseReduction) RawJSON() string { return r.JSON.raw }
func (r *RealtimeSessionCreateResponseAudioInputNoiseReduction) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for input audio transcription.
type RealtimeSessionCreateResponseAudioInputTranscription struct {
	// The language of the input audio.
	Language string `json:"language"`
	// The model to use for transcription.
	Model string `json:"model"`
	// Optional text to guide the model's style or continue a previous audio segment.
	Prompt string `json:"prompt"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Language    respjson.Field
		Model       respjson.Field
		Prompt      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeSessionCreateResponseAudioInputTranscription) RawJSON() string { return r.JSON.raw }
func (r *RealtimeSessionCreateResponseAudioInputTranscription) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for turn detection.
type RealtimeSessionCreateResponseAudioInputTurnDetection struct {
	PrefixPaddingMs   int64   `json:"prefix_padding_ms"`
	SilenceDurationMs int64   `json:"silence_duration_ms"`
	Threshold         float64 `json:"threshold"`
	// Type of turn detection, only `server_vad` is currently supported.
	Type string `json:"type"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		PrefixPaddingMs   respjson.Field
		SilenceDurationMs respjson.Field
		Threshold         respjson.Field
		Type              respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeSessionCreateResponseAudioInputTurnDetection) RawJSON() string { return r.JSON.raw }
func (r *RealtimeSessionCreateResponseAudioInputTurnDetection) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type RealtimeSessionCreateResponseAudioOutput struct {
	// The format of output audio. Options are `pcm16`, `g711_ulaw`, or `g711_alaw`.
	Format string  `json:"format"`
	Speed  float64 `json:"speed"`
	Voice  string  `json:"voice"`
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
func (r RealtimeSessionCreateResponseAudioOutput) RawJSON() string { return r.JSON.raw }
func (r *RealtimeSessionCreateResponseAudioOutput) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// RealtimeSessionCreateResponseMaxOutputTokensUnion contains all possible
// properties and values from [int64], [constant.Inf].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfInt OfInf]
type RealtimeSessionCreateResponseMaxOutputTokensUnion struct {
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

func (u RealtimeSessionCreateResponseMaxOutputTokensUnion) AsInt() (v int64) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeSessionCreateResponseMaxOutputTokensUnion) AsInf() (v constant.Inf) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u RealtimeSessionCreateResponseMaxOutputTokensUnion) RawJSON() string { return u.JSON.raw }

func (r *RealtimeSessionCreateResponseMaxOutputTokensUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type RealtimeSessionCreateResponseTool struct {
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
	Type string `json:"type"`
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
func (r RealtimeSessionCreateResponseTool) RawJSON() string { return r.JSON.raw }
func (r *RealtimeSessionCreateResponseTool) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// RealtimeSessionCreateResponseTracingUnion contains all possible properties and
// values from [constant.Auto],
// [RealtimeSessionCreateResponseTracingTracingConfiguration].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfAuto]
type RealtimeSessionCreateResponseTracingUnion struct {
	// This field will be present if the value is a [constant.Auto] instead of an
	// object.
	OfAuto constant.Auto `json:",inline"`
	// This field is from variant
	// [RealtimeSessionCreateResponseTracingTracingConfiguration].
	GroupID string `json:"group_id"`
	// This field is from variant
	// [RealtimeSessionCreateResponseTracingTracingConfiguration].
	Metadata any `json:"metadata"`
	// This field is from variant
	// [RealtimeSessionCreateResponseTracingTracingConfiguration].
	WorkflowName string `json:"workflow_name"`
	JSON         struct {
		OfAuto       respjson.Field
		GroupID      respjson.Field
		Metadata     respjson.Field
		WorkflowName respjson.Field
		raw          string
	} `json:"-"`
}

func (u RealtimeSessionCreateResponseTracingUnion) AsAuto() (v constant.Auto) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RealtimeSessionCreateResponseTracingUnion) AsTracingConfiguration() (v RealtimeSessionCreateResponseTracingTracingConfiguration) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u RealtimeSessionCreateResponseTracingUnion) RawJSON() string { return u.JSON.raw }

func (r *RealtimeSessionCreateResponseTracingUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Granular configuration for tracing.
type RealtimeSessionCreateResponseTracingTracingConfiguration struct {
	// The group id to attach to this trace to enable filtering and grouping in the
	// traces dashboard.
	GroupID string `json:"group_id"`
	// The arbitrary metadata to attach to this trace to enable filtering in the traces
	// dashboard.
	Metadata any `json:"metadata"`
	// The name of the workflow to attach to this trace. This is used to name the trace
	// in the traces dashboard.
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
func (r RealtimeSessionCreateResponseTracingTracingConfiguration) RawJSON() string { return r.JSON.raw }
func (r *RealtimeSessionCreateResponseTracingTracingConfiguration) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for turn detection. Can be set to `null` to turn off. Server VAD
// means that the model will detect the start and end of speech based on audio
// volume and respond at the end of user speech.
type RealtimeSessionCreateResponseTurnDetection struct {
	// Amount of audio to include before the VAD detected speech (in milliseconds).
	// Defaults to 300ms.
	PrefixPaddingMs int64 `json:"prefix_padding_ms"`
	// Duration of silence to detect speech stop (in milliseconds). Defaults to 500ms.
	// With shorter values the model will respond more quickly, but may jump in on
	// short pauses from the user.
	SilenceDurationMs int64 `json:"silence_duration_ms"`
	// Activation threshold for VAD (0.0 to 1.0), this defaults to 0.5. A higher
	// threshold will require louder audio to activate the model, and thus might
	// perform better in noisy environments.
	Threshold float64 `json:"threshold"`
	// Type of turn detection, only `server_vad` is currently supported.
	Type string `json:"type"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		PrefixPaddingMs   respjson.Field
		SilenceDurationMs respjson.Field
		Threshold         respjson.Field
		Type              respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeSessionCreateResponseTurnDetection) RawJSON() string { return r.JSON.raw }
func (r *RealtimeSessionCreateResponseTurnDetection) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Response from creating a session and client secret for the Realtime API.
type ClientSecretNewResponse struct {
	// Expiration timestamp for the client secret, in seconds since epoch.
	ExpiresAt int64 `json:"expires_at,required"`
	// The session configuration for either a realtime or transcription session.
	Session ClientSecretNewResponseSessionUnion `json:"session,required"`
	// The generated client secret value.
	Value string `json:"value,required"`
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
func (r ClientSecretNewResponse) RawJSON() string { return r.JSON.raw }
func (r *ClientSecretNewResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ClientSecretNewResponseSessionUnion contains all possible properties and values
// from [RealtimeSessionCreateResponse],
// [ClientSecretNewResponseSessionRealtimeTranscriptionSessionConfigurationObject].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ClientSecretNewResponseSessionUnion struct {
	ID string `json:"id"`
	// This field is a union of [RealtimeSessionCreateResponseAudio],
	// [ClientSecretNewResponseSessionRealtimeTranscriptionSessionConfigurationObjectAudio]
	Audio     ClientSecretNewResponseSessionUnionAudio `json:"audio"`
	ExpiresAt int64                                    `json:"expires_at"`
	Include   []string                                 `json:"include"`
	// This field is from variant [RealtimeSessionCreateResponse].
	Instructions string `json:"instructions"`
	// This field is from variant [RealtimeSessionCreateResponse].
	MaxOutputTokens RealtimeSessionCreateResponseMaxOutputTokensUnion `json:"max_output_tokens"`
	// This field is from variant [RealtimeSessionCreateResponse].
	Model  string `json:"model"`
	Object string `json:"object"`
	// This field is from variant [RealtimeSessionCreateResponse].
	OutputModalities []string `json:"output_modalities"`
	// This field is from variant [RealtimeSessionCreateResponse].
	ToolChoice string `json:"tool_choice"`
	// This field is from variant [RealtimeSessionCreateResponse].
	Tools []RealtimeSessionCreateResponseTool `json:"tools"`
	// This field is from variant [RealtimeSessionCreateResponse].
	Tracing RealtimeSessionCreateResponseTracingUnion `json:"tracing"`
	// This field is from variant [RealtimeSessionCreateResponse].
	TurnDetection RealtimeSessionCreateResponseTurnDetection `json:"turn_detection"`
	JSON          struct {
		ID               respjson.Field
		Audio            respjson.Field
		ExpiresAt        respjson.Field
		Include          respjson.Field
		Instructions     respjson.Field
		MaxOutputTokens  respjson.Field
		Model            respjson.Field
		Object           respjson.Field
		OutputModalities respjson.Field
		ToolChoice       respjson.Field
		Tools            respjson.Field
		Tracing          respjson.Field
		TurnDetection    respjson.Field
		raw              string
	} `json:"-"`
}

func (u ClientSecretNewResponseSessionUnion) AsRealtimeSessionConfigurationObject() (v RealtimeSessionCreateResponse) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ClientSecretNewResponseSessionUnion) AsRealtimeTranscriptionSessionConfigurationObject() (v ClientSecretNewResponseSessionRealtimeTranscriptionSessionConfigurationObject) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ClientSecretNewResponseSessionUnion) RawJSON() string { return u.JSON.raw }

func (r *ClientSecretNewResponseSessionUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ClientSecretNewResponseSessionUnionAudio is an implicit subunion of
// [ClientSecretNewResponseSessionUnion]. ClientSecretNewResponseSessionUnionAudio
// provides convenient access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [ClientSecretNewResponseSessionUnion].
type ClientSecretNewResponseSessionUnionAudio struct {
	// This field is a union of [RealtimeSessionCreateResponseAudioInput],
	// [ClientSecretNewResponseSessionRealtimeTranscriptionSessionConfigurationObjectAudioInput]
	Input ClientSecretNewResponseSessionUnionAudioInput `json:"input"`
	// This field is from variant [RealtimeSessionCreateResponseAudio].
	Output RealtimeSessionCreateResponseAudioOutput `json:"output"`
	JSON   struct {
		Input  respjson.Field
		Output respjson.Field
		raw    string
	} `json:"-"`
}

func (r *ClientSecretNewResponseSessionUnionAudio) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ClientSecretNewResponseSessionUnionAudioInput is an implicit subunion of
// [ClientSecretNewResponseSessionUnion].
// ClientSecretNewResponseSessionUnionAudioInput provides convenient access to the
// sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [ClientSecretNewResponseSessionUnion].
type ClientSecretNewResponseSessionUnionAudioInput struct {
	Format string `json:"format"`
	// This field is a union of
	// [RealtimeSessionCreateResponseAudioInputNoiseReduction],
	// [ClientSecretNewResponseSessionRealtimeTranscriptionSessionConfigurationObjectAudioInputNoiseReduction]
	NoiseReduction ClientSecretNewResponseSessionUnionAudioInputNoiseReduction `json:"noise_reduction"`
	// This field is a union of [RealtimeSessionCreateResponseAudioInputTranscription],
	// [ClientSecretNewResponseSessionRealtimeTranscriptionSessionConfigurationObjectAudioInputTranscription]
	Transcription ClientSecretNewResponseSessionUnionAudioInputTranscription `json:"transcription"`
	// This field is a union of [RealtimeSessionCreateResponseAudioInputTurnDetection],
	// [ClientSecretNewResponseSessionRealtimeTranscriptionSessionConfigurationObjectAudioInputTurnDetection]
	TurnDetection ClientSecretNewResponseSessionUnionAudioInputTurnDetection `json:"turn_detection"`
	JSON          struct {
		Format         respjson.Field
		NoiseReduction respjson.Field
		Transcription  respjson.Field
		TurnDetection  respjson.Field
		raw            string
	} `json:"-"`
}

func (r *ClientSecretNewResponseSessionUnionAudioInput) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ClientSecretNewResponseSessionUnionAudioInputNoiseReduction is an implicit
// subunion of [ClientSecretNewResponseSessionUnion].
// ClientSecretNewResponseSessionUnionAudioInputNoiseReduction provides convenient
// access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [ClientSecretNewResponseSessionUnion].
type ClientSecretNewResponseSessionUnionAudioInputNoiseReduction struct {
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

func (r *ClientSecretNewResponseSessionUnionAudioInputNoiseReduction) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ClientSecretNewResponseSessionUnionAudioInputTranscription is an implicit
// subunion of [ClientSecretNewResponseSessionUnion].
// ClientSecretNewResponseSessionUnionAudioInputTranscription provides convenient
// access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [ClientSecretNewResponseSessionUnion].
type ClientSecretNewResponseSessionUnionAudioInputTranscription struct {
	Language string `json:"language"`
	Model    string `json:"model"`
	Prompt   string `json:"prompt"`
	JSON     struct {
		Language respjson.Field
		Model    respjson.Field
		Prompt   respjson.Field
		raw      string
	} `json:"-"`
}

func (r *ClientSecretNewResponseSessionUnionAudioInputTranscription) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ClientSecretNewResponseSessionUnionAudioInputTurnDetection is an implicit
// subunion of [ClientSecretNewResponseSessionUnion].
// ClientSecretNewResponseSessionUnionAudioInputTurnDetection provides convenient
// access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [ClientSecretNewResponseSessionUnion].
type ClientSecretNewResponseSessionUnionAudioInputTurnDetection struct {
	PrefixPaddingMs   int64   `json:"prefix_padding_ms"`
	SilenceDurationMs int64   `json:"silence_duration_ms"`
	Threshold         float64 `json:"threshold"`
	Type              string  `json:"type"`
	JSON              struct {
		PrefixPaddingMs   respjson.Field
		SilenceDurationMs respjson.Field
		Threshold         respjson.Field
		Type              respjson.Field
		raw               string
	} `json:"-"`
}

func (r *ClientSecretNewResponseSessionUnionAudioInputTurnDetection) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A Realtime transcription session configuration object.
type ClientSecretNewResponseSessionRealtimeTranscriptionSessionConfigurationObject struct {
	// Unique identifier for the session that looks like `sess_1234567890abcdef`.
	ID string `json:"id"`
	// Configuration for input audio for the session.
	Audio ClientSecretNewResponseSessionRealtimeTranscriptionSessionConfigurationObjectAudio `json:"audio"`
	// Expiration timestamp for the session, in seconds since epoch.
	ExpiresAt int64 `json:"expires_at"`
	// Additional fields to include in server outputs.
	//
	//   - `item.input_audio_transcription.logprobs`: Include logprobs for input audio
	//     transcription.
	//
	// Any of "item.input_audio_transcription.logprobs".
	Include []string `json:"include"`
	// The object type. Always `realtime.transcription_session`.
	Object string `json:"object"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Audio       respjson.Field
		ExpiresAt   respjson.Field
		Include     respjson.Field
		Object      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ClientSecretNewResponseSessionRealtimeTranscriptionSessionConfigurationObject) RawJSON() string {
	return r.JSON.raw
}
func (r *ClientSecretNewResponseSessionRealtimeTranscriptionSessionConfigurationObject) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for input audio for the session.
type ClientSecretNewResponseSessionRealtimeTranscriptionSessionConfigurationObjectAudio struct {
	Input ClientSecretNewResponseSessionRealtimeTranscriptionSessionConfigurationObjectAudioInput `json:"input"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Input       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ClientSecretNewResponseSessionRealtimeTranscriptionSessionConfigurationObjectAudio) RawJSON() string {
	return r.JSON.raw
}
func (r *ClientSecretNewResponseSessionRealtimeTranscriptionSessionConfigurationObjectAudio) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ClientSecretNewResponseSessionRealtimeTranscriptionSessionConfigurationObjectAudioInput struct {
	// The format of input audio. Options are `pcm16`, `g711_ulaw`, or `g711_alaw`.
	Format string `json:"format"`
	// Configuration for input audio noise reduction.
	NoiseReduction ClientSecretNewResponseSessionRealtimeTranscriptionSessionConfigurationObjectAudioInputNoiseReduction `json:"noise_reduction"`
	// Configuration of the transcription model.
	Transcription ClientSecretNewResponseSessionRealtimeTranscriptionSessionConfigurationObjectAudioInputTranscription `json:"transcription"`
	// Configuration for turn detection.
	TurnDetection ClientSecretNewResponseSessionRealtimeTranscriptionSessionConfigurationObjectAudioInputTurnDetection `json:"turn_detection"`
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
func (r ClientSecretNewResponseSessionRealtimeTranscriptionSessionConfigurationObjectAudioInput) RawJSON() string {
	return r.JSON.raw
}
func (r *ClientSecretNewResponseSessionRealtimeTranscriptionSessionConfigurationObjectAudioInput) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for input audio noise reduction.
type ClientSecretNewResponseSessionRealtimeTranscriptionSessionConfigurationObjectAudioInputNoiseReduction struct {
	// Any of "near_field", "far_field".
	Type string `json:"type"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ClientSecretNewResponseSessionRealtimeTranscriptionSessionConfigurationObjectAudioInputNoiseReduction) RawJSON() string {
	return r.JSON.raw
}
func (r *ClientSecretNewResponseSessionRealtimeTranscriptionSessionConfigurationObjectAudioInputNoiseReduction) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration of the transcription model.
type ClientSecretNewResponseSessionRealtimeTranscriptionSessionConfigurationObjectAudioInputTranscription struct {
	// The language of the input audio. Supplying the input language in
	// [ISO-639-1](https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes) (e.g. `en`)
	// format will improve accuracy and latency.
	Language string `json:"language"`
	// The model to use for transcription. Can be `gpt-4o-transcribe`,
	// `gpt-4o-mini-transcribe`, or `whisper-1`.
	//
	// Any of "gpt-4o-transcribe", "gpt-4o-mini-transcribe", "whisper-1".
	Model string `json:"model"`
	// An optional text to guide the model's style or continue a previous audio
	// segment. The
	// [prompt](https://platform.openai.com/docs/guides/speech-to-text#prompting)
	// should match the audio language.
	Prompt string `json:"prompt"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Language    respjson.Field
		Model       respjson.Field
		Prompt      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ClientSecretNewResponseSessionRealtimeTranscriptionSessionConfigurationObjectAudioInputTranscription) RawJSON() string {
	return r.JSON.raw
}
func (r *ClientSecretNewResponseSessionRealtimeTranscriptionSessionConfigurationObjectAudioInputTranscription) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for turn detection.
type ClientSecretNewResponseSessionRealtimeTranscriptionSessionConfigurationObjectAudioInputTurnDetection struct {
	PrefixPaddingMs   int64   `json:"prefix_padding_ms"`
	SilenceDurationMs int64   `json:"silence_duration_ms"`
	Threshold         float64 `json:"threshold"`
	// Type of turn detection, only `server_vad` is currently supported.
	Type string `json:"type"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		PrefixPaddingMs   respjson.Field
		SilenceDurationMs respjson.Field
		Threshold         respjson.Field
		Type              respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ClientSecretNewResponseSessionRealtimeTranscriptionSessionConfigurationObjectAudioInputTurnDetection) RawJSON() string {
	return r.JSON.raw
}
func (r *ClientSecretNewResponseSessionRealtimeTranscriptionSessionConfigurationObjectAudioInputTurnDetection) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ClientSecretNewParams struct {
	// Configuration for the ephemeral token expiration.
	ExpiresAfter ClientSecretNewParamsExpiresAfter `json:"expires_after,omitzero"`
	// Session configuration to use for the client secret. Choose either a realtime
	// session or a transcription session.
	Session ClientSecretNewParamsSessionUnion `json:"session,omitzero"`
	paramObj
}

func (r ClientSecretNewParams) MarshalJSON() (data []byte, err error) {
	type shadow ClientSecretNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ClientSecretNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for the ephemeral token expiration.
type ClientSecretNewParamsExpiresAfter struct {
	// The number of seconds from the anchor point to the expiration. Select a value
	// between `10` and `7200`.
	Seconds param.Opt[int64] `json:"seconds,omitzero"`
	// The anchor point for the ephemeral token expiration. Only `created_at` is
	// currently supported.
	//
	// Any of "created_at".
	Anchor string `json:"anchor,omitzero"`
	paramObj
}

func (r ClientSecretNewParamsExpiresAfter) MarshalJSON() (data []byte, err error) {
	type shadow ClientSecretNewParamsExpiresAfter
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ClientSecretNewParamsExpiresAfter) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[ClientSecretNewParamsExpiresAfter](
		"anchor", "created_at",
	)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ClientSecretNewParamsSessionUnion struct {
	OfRealtime      *RealtimeSessionCreateRequestParam              `json:",omitzero,inline"`
	OfTranscription *RealtimeTranscriptionSessionCreateRequestParam `json:",omitzero,inline"`
	paramUnion
}

func (u ClientSecretNewParamsSessionUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfRealtime, u.OfTranscription)
}
func (u *ClientSecretNewParamsSessionUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *ClientSecretNewParamsSessionUnion) asAny() any {
	if !param.IsOmitted(u.OfRealtime) {
		return u.OfRealtime
	} else if !param.IsOmitted(u.OfTranscription) {
		return u.OfTranscription
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ClientSecretNewParamsSessionUnion) GetAudio() *RealtimeAudioConfigParam {
	if vt := u.OfRealtime; vt != nil {
		return &vt.Audio
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ClientSecretNewParamsSessionUnion) GetClientSecret() *RealtimeClientSecretConfigParam {
	if vt := u.OfRealtime; vt != nil {
		return &vt.ClientSecret
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ClientSecretNewParamsSessionUnion) GetInstructions() *string {
	if vt := u.OfRealtime; vt != nil && vt.Instructions.Valid() {
		return &vt.Instructions.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ClientSecretNewParamsSessionUnion) GetMaxOutputTokens() *RealtimeSessionCreateRequestMaxOutputTokensUnionParam {
	if vt := u.OfRealtime; vt != nil {
		return &vt.MaxOutputTokens
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ClientSecretNewParamsSessionUnion) GetOutputModalities() []string {
	if vt := u.OfRealtime; vt != nil {
		return vt.OutputModalities
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ClientSecretNewParamsSessionUnion) GetPrompt() *responses.ResponsePromptParam {
	if vt := u.OfRealtime; vt != nil {
		return &vt.Prompt
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ClientSecretNewParamsSessionUnion) GetTemperature() *float64 {
	if vt := u.OfRealtime; vt != nil && vt.Temperature.Valid() {
		return &vt.Temperature.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ClientSecretNewParamsSessionUnion) GetToolChoice() *RealtimeToolChoiceConfigUnionParam {
	if vt := u.OfRealtime; vt != nil {
		return &vt.ToolChoice
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ClientSecretNewParamsSessionUnion) GetTools() RealtimeToolsConfigParam {
	if vt := u.OfRealtime; vt != nil {
		return vt.Tools
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ClientSecretNewParamsSessionUnion) GetTracing() *RealtimeTracingConfigUnionParam {
	if vt := u.OfRealtime; vt != nil {
		return &vt.Tracing
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ClientSecretNewParamsSessionUnion) GetTruncation() *RealtimeTruncationUnionParam {
	if vt := u.OfRealtime; vt != nil {
		return &vt.Truncation
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ClientSecretNewParamsSessionUnion) GetInputAudioFormat() *string {
	if vt := u.OfTranscription; vt != nil {
		return (*string)(&vt.InputAudioFormat)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ClientSecretNewParamsSessionUnion) GetInputAudioNoiseReduction() *RealtimeTranscriptionSessionCreateRequestInputAudioNoiseReductionParam {
	if vt := u.OfTranscription; vt != nil {
		return &vt.InputAudioNoiseReduction
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ClientSecretNewParamsSessionUnion) GetInputAudioTranscription() *RealtimeTranscriptionSessionCreateRequestInputAudioTranscriptionParam {
	if vt := u.OfTranscription; vt != nil {
		return &vt.InputAudioTranscription
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ClientSecretNewParamsSessionUnion) GetTurnDetection() *RealtimeTranscriptionSessionCreateRequestTurnDetectionParam {
	if vt := u.OfTranscription; vt != nil {
		return &vt.TurnDetection
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ClientSecretNewParamsSessionUnion) GetType() *string {
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
func (u ClientSecretNewParamsSessionUnion) GetModel() (res clientSecretNewParamsSessionUnionModel) {
	if vt := u.OfRealtime; vt != nil {
		res.any = &vt.Model
	} else if vt := u.OfTranscription; vt != nil {
		res.any = &vt.Model
	}
	return
}

// Can have the runtime types [*RealtimeSessionCreateRequestModel],
// [*RealtimeTranscriptionSessionCreateRequestModel]
type clientSecretNewParamsSessionUnionModel struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *realtime.RealtimeSessionCreateRequestModel:
//	case *realtime.RealtimeTranscriptionSessionCreateRequestModel:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u clientSecretNewParamsSessionUnionModel) AsAny() any { return u.any }

// Returns a pointer to the underlying variant's Include property, if present.
func (u ClientSecretNewParamsSessionUnion) GetInclude() []string {
	if vt := u.OfRealtime; vt != nil {
		return vt.Include
	} else if vt := u.OfTranscription; vt != nil {
		return vt.Include
	}
	return nil
}

func init() {
	apijson.RegisterUnion[ClientSecretNewParamsSessionUnion](
		"type",
		apijson.Discriminator[RealtimeSessionCreateRequestParam]("realtime"),
		apijson.Discriminator[RealtimeTranscriptionSessionCreateRequestParam]("transcription"),
	)
}

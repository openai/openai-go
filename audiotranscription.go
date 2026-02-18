// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"slices"

	"github.com/openai/openai-go/v3/internal/apiform"
	"github.com/openai/openai-go/v3/internal/apijson"
	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/packages/respjson"
	"github.com/openai/openai-go/v3/packages/ssestream"
	"github.com/openai/openai-go/v3/shared/constant"
)

// AudioTranscriptionService contains methods and other services that help with
// interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAudioTranscriptionService] method instead.
type AudioTranscriptionService struct {
	Options []option.RequestOption
}

// NewAudioTranscriptionService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewAudioTranscriptionService(opts ...option.RequestOption) (r AudioTranscriptionService) {
	r = AudioTranscriptionService{}
	r.Options = opts
	return
}

// Transcribes audio into the input language.
//
// Returns a transcription object in `json`, `diarized_json`, or `verbose_json`
// format, or a stream of transcript events.
func (r *AudioTranscriptionService) New(ctx context.Context, body AudioTranscriptionNewParams, opts ...option.RequestOption) (res *AudioTranscriptionNewResponseUnion, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "audio/transcriptions"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Transcribes audio into the input language.
//
// Returns a transcription object in `json`, `diarized_json`, or `verbose_json`
// format, or a stream of transcript events.
func (r *AudioTranscriptionService) NewStreaming(ctx context.Context, body AudioTranscriptionNewParams, opts ...option.RequestOption) (stream *ssestream.Stream[TranscriptionStreamEventUnion]) {
	var (
		raw *http.Response
		err error
	)
	opts = slices.Concat(r.Options, opts)
	body.SetExtraFields(map[string]any{
		"stream": "true",
	})
	path := "audio/transcriptions"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &raw, opts...)
	return ssestream.NewStream[TranscriptionStreamEventUnion](ssestream.NewDecoder(raw), err)
}

// Represents a transcription response returned by model, based on the provided
// input.
type Transcription struct {
	// The transcribed text.
	Text string `json:"text,required"`
	// The log probabilities of the tokens in the transcription. Only returned with the
	// models `gpt-4o-transcribe` and `gpt-4o-mini-transcribe` if `logprobs` is added
	// to the `include` array.
	Logprobs []TranscriptionLogprob `json:"logprobs"`
	// Token usage statistics for the request.
	Usage TranscriptionUsageUnion `json:"usage"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Text        respjson.Field
		Logprobs    respjson.Field
		Usage       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r Transcription) RawJSON() string { return r.JSON.raw }
func (r *Transcription) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type TranscriptionLogprob struct {
	// The token in the transcription.
	Token string `json:"token"`
	// The bytes of the token.
	Bytes []float64 `json:"bytes"`
	// The log probability of the token.
	Logprob float64 `json:"logprob"`
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
func (r TranscriptionLogprob) RawJSON() string { return r.JSON.raw }
func (r *TranscriptionLogprob) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// TranscriptionUsageUnion contains all possible properties and values from
// [TranscriptionUsageTokens], [TranscriptionUsageDuration].
//
// Use the [TranscriptionUsageUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type TranscriptionUsageUnion struct {
	// This field is from variant [TranscriptionUsageTokens].
	InputTokens int64 `json:"input_tokens"`
	// This field is from variant [TranscriptionUsageTokens].
	OutputTokens int64 `json:"output_tokens"`
	// This field is from variant [TranscriptionUsageTokens].
	TotalTokens int64 `json:"total_tokens"`
	// Any of "tokens", "duration".
	Type string `json:"type"`
	// This field is from variant [TranscriptionUsageTokens].
	InputTokenDetails TranscriptionUsageTokensInputTokenDetails `json:"input_token_details"`
	// This field is from variant [TranscriptionUsageDuration].
	Seconds float64 `json:"seconds"`
	JSON    struct {
		InputTokens       respjson.Field
		OutputTokens      respjson.Field
		TotalTokens       respjson.Field
		Type              respjson.Field
		InputTokenDetails respjson.Field
		Seconds           respjson.Field
		raw               string
	} `json:"-"`
}

// anyTranscriptionUsage is implemented by each variant of
// [TranscriptionUsageUnion] to add type safety for the return type of
// [TranscriptionUsageUnion.AsAny]
type anyTranscriptionUsage interface {
	implTranscriptionUsageUnion()
}

func (TranscriptionUsageTokens) implTranscriptionUsageUnion()   {}
func (TranscriptionUsageDuration) implTranscriptionUsageUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := TranscriptionUsageUnion.AsAny().(type) {
//	case openai.TranscriptionUsageTokens:
//	case openai.TranscriptionUsageDuration:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u TranscriptionUsageUnion) AsAny() anyTranscriptionUsage {
	switch u.Type {
	case "tokens":
		return u.AsTokens()
	case "duration":
		return u.AsDuration()
	}
	return nil
}

func (u TranscriptionUsageUnion) AsTokens() (v TranscriptionUsageTokens) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u TranscriptionUsageUnion) AsDuration() (v TranscriptionUsageDuration) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u TranscriptionUsageUnion) RawJSON() string { return u.JSON.raw }

func (r *TranscriptionUsageUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Usage statistics for models billed by token usage.
type TranscriptionUsageTokens struct {
	// Number of input tokens billed for this request.
	InputTokens int64 `json:"input_tokens,required"`
	// Number of output tokens generated.
	OutputTokens int64 `json:"output_tokens,required"`
	// Total number of tokens used (input + output).
	TotalTokens int64 `json:"total_tokens,required"`
	// The type of the usage object. Always `tokens` for this variant.
	Type constant.Tokens `json:"type,required"`
	// Details about the input tokens billed for this request.
	InputTokenDetails TranscriptionUsageTokensInputTokenDetails `json:"input_token_details"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		InputTokens       respjson.Field
		OutputTokens      respjson.Field
		TotalTokens       respjson.Field
		Type              respjson.Field
		InputTokenDetails respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r TranscriptionUsageTokens) RawJSON() string { return r.JSON.raw }
func (r *TranscriptionUsageTokens) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Details about the input tokens billed for this request.
type TranscriptionUsageTokensInputTokenDetails struct {
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
func (r TranscriptionUsageTokensInputTokenDetails) RawJSON() string { return r.JSON.raw }
func (r *TranscriptionUsageTokensInputTokenDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Usage statistics for models billed by audio input duration.
type TranscriptionUsageDuration struct {
	// Duration of the input audio in seconds.
	Seconds float64 `json:"seconds,required"`
	// The type of the usage object. Always `duration` for this variant.
	Type constant.Duration `json:"type,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Seconds     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r TranscriptionUsageDuration) RawJSON() string { return r.JSON.raw }
func (r *TranscriptionUsageDuration) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type TranscriptionInclude string

const (
	TranscriptionIncludeLogprobs TranscriptionInclude = "logprobs"
)

type TranscriptionSegment struct {
	// Unique identifier of the segment.
	ID int64 `json:"id,required"`
	// Average logprob of the segment. If the value is lower than -1, consider the
	// logprobs failed.
	AvgLogprob float64 `json:"avg_logprob,required"`
	// Compression ratio of the segment. If the value is greater than 2.4, consider the
	// compression failed.
	CompressionRatio float64 `json:"compression_ratio,required"`
	// End time of the segment in seconds.
	End float64 `json:"end,required"`
	// Probability of no speech in the segment. If the value is higher than 1.0 and the
	// `avg_logprob` is below -1, consider this segment silent.
	NoSpeechProb float64 `json:"no_speech_prob,required"`
	// Seek offset of the segment.
	Seek int64 `json:"seek,required"`
	// Start time of the segment in seconds.
	Start float64 `json:"start,required"`
	// Temperature parameter used for generating the segment.
	Temperature float64 `json:"temperature,required"`
	// Text content of the segment.
	Text string `json:"text,required"`
	// Array of token IDs for the text content.
	Tokens []int64 `json:"tokens,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID               respjson.Field
		AvgLogprob       respjson.Field
		CompressionRatio respjson.Field
		End              respjson.Field
		NoSpeechProb     respjson.Field
		Seek             respjson.Field
		Start            respjson.Field
		Temperature      respjson.Field
		Text             respjson.Field
		Tokens           respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r TranscriptionSegment) RawJSON() string { return r.JSON.raw }
func (r *TranscriptionSegment) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// TranscriptionStreamEventUnion contains all possible properties and values from
// [TranscriptionTextSegmentEvent], [TranscriptionTextDeltaEvent],
// [TranscriptionTextDoneEvent].
//
// Use the [TranscriptionStreamEventUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type TranscriptionStreamEventUnion struct {
	// This field is from variant [TranscriptionTextSegmentEvent].
	ID string `json:"id"`
	// This field is from variant [TranscriptionTextSegmentEvent].
	End float64 `json:"end"`
	// This field is from variant [TranscriptionTextSegmentEvent].
	Speaker string `json:"speaker"`
	// This field is from variant [TranscriptionTextSegmentEvent].
	Start float64 `json:"start"`
	Text  string  `json:"text"`
	// Any of "transcript.text.segment", "transcript.text.delta",
	// "transcript.text.done".
	Type string `json:"type"`
	// This field is from variant [TranscriptionTextDeltaEvent].
	Delta string `json:"delta"`
	// This field is a union of [[]TranscriptionTextDeltaEventLogprob],
	// [[]TranscriptionTextDoneEventLogprob]
	Logprobs TranscriptionStreamEventUnionLogprobs `json:"logprobs"`
	// This field is from variant [TranscriptionTextDeltaEvent].
	SegmentID string `json:"segment_id"`
	// This field is from variant [TranscriptionTextDoneEvent].
	Usage TranscriptionTextDoneEventUsage `json:"usage"`
	JSON  struct {
		ID        respjson.Field
		End       respjson.Field
		Speaker   respjson.Field
		Start     respjson.Field
		Text      respjson.Field
		Type      respjson.Field
		Delta     respjson.Field
		Logprobs  respjson.Field
		SegmentID respjson.Field
		Usage     respjson.Field
		raw       string
	} `json:"-"`
}

// anyTranscriptionStreamEvent is implemented by each variant of
// [TranscriptionStreamEventUnion] to add type safety for the return type of
// [TranscriptionStreamEventUnion.AsAny]
type anyTranscriptionStreamEvent interface {
	implTranscriptionStreamEventUnion()
}

func (TranscriptionTextSegmentEvent) implTranscriptionStreamEventUnion() {}
func (TranscriptionTextDeltaEvent) implTranscriptionStreamEventUnion()   {}
func (TranscriptionTextDoneEvent) implTranscriptionStreamEventUnion()    {}

// Use the following switch statement to find the correct variant
//
//	switch variant := TranscriptionStreamEventUnion.AsAny().(type) {
//	case openai.TranscriptionTextSegmentEvent:
//	case openai.TranscriptionTextDeltaEvent:
//	case openai.TranscriptionTextDoneEvent:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u TranscriptionStreamEventUnion) AsAny() anyTranscriptionStreamEvent {
	switch u.Type {
	case "transcript.text.segment":
		return u.AsTranscriptTextSegment()
	case "transcript.text.delta":
		return u.AsTranscriptTextDelta()
	case "transcript.text.done":
		return u.AsTranscriptTextDone()
	}
	return nil
}

func (u TranscriptionStreamEventUnion) AsTranscriptTextSegment() (v TranscriptionTextSegmentEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u TranscriptionStreamEventUnion) AsTranscriptTextDelta() (v TranscriptionTextDeltaEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u TranscriptionStreamEventUnion) AsTranscriptTextDone() (v TranscriptionTextDoneEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u TranscriptionStreamEventUnion) RawJSON() string { return u.JSON.raw }

func (r *TranscriptionStreamEventUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// TranscriptionStreamEventUnionLogprobs is an implicit subunion of
// [TranscriptionStreamEventUnion]. TranscriptionStreamEventUnionLogprobs provides
// convenient access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [TranscriptionStreamEventUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfTranscriptionTextDeltaEventLogprobs
// OfTranscriptionTextDoneEventLogprobs]
type TranscriptionStreamEventUnionLogprobs struct {
	// This field will be present if the value is a
	// [[]TranscriptionTextDeltaEventLogprob] instead of an object.
	OfTranscriptionTextDeltaEventLogprobs []TranscriptionTextDeltaEventLogprob `json:",inline"`
	// This field will be present if the value is a
	// [[]TranscriptionTextDoneEventLogprob] instead of an object.
	OfTranscriptionTextDoneEventLogprobs []TranscriptionTextDoneEventLogprob `json:",inline"`
	JSON                                 struct {
		OfTranscriptionTextDeltaEventLogprobs respjson.Field
		OfTranscriptionTextDoneEventLogprobs  respjson.Field
		raw                                   string
	} `json:"-"`
}

func (r *TranscriptionStreamEventUnionLogprobs) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Emitted when there is an additional text delta. This is also the first event
// emitted when the transcription starts. Only emitted when you
// [create a transcription](https://platform.openai.com/docs/api-reference/audio/create-transcription)
// with the `Stream` parameter set to `true`.
type TranscriptionTextDeltaEvent struct {
	// The text delta that was additionally transcribed.
	Delta string `json:"delta,required"`
	// The type of the event. Always `transcript.text.delta`.
	Type constant.TranscriptTextDelta `json:"type,required"`
	// The log probabilities of the delta. Only included if you
	// [create a transcription](https://platform.openai.com/docs/api-reference/audio/create-transcription)
	// with the `include[]` parameter set to `logprobs`.
	Logprobs []TranscriptionTextDeltaEventLogprob `json:"logprobs"`
	// Identifier of the diarized segment that this delta belongs to. Only present when
	// using `gpt-4o-transcribe-diarize`.
	SegmentID string `json:"segment_id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Delta       respjson.Field
		Type        respjson.Field
		Logprobs    respjson.Field
		SegmentID   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r TranscriptionTextDeltaEvent) RawJSON() string { return r.JSON.raw }
func (r *TranscriptionTextDeltaEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type TranscriptionTextDeltaEventLogprob struct {
	// The token that was used to generate the log probability.
	Token string `json:"token"`
	// The bytes that were used to generate the log probability.
	Bytes []int64 `json:"bytes"`
	// The log probability of the token.
	Logprob float64 `json:"logprob"`
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
func (r TranscriptionTextDeltaEventLogprob) RawJSON() string { return r.JSON.raw }
func (r *TranscriptionTextDeltaEventLogprob) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Emitted when the transcription is complete. Contains the complete transcription
// text. Only emitted when you
// [create a transcription](https://platform.openai.com/docs/api-reference/audio/create-transcription)
// with the `Stream` parameter set to `true`.
type TranscriptionTextDoneEvent struct {
	// The text that was transcribed.
	Text string `json:"text,required"`
	// The type of the event. Always `transcript.text.done`.
	Type constant.TranscriptTextDone `json:"type,required"`
	// The log probabilities of the individual tokens in the transcription. Only
	// included if you
	// [create a transcription](https://platform.openai.com/docs/api-reference/audio/create-transcription)
	// with the `include[]` parameter set to `logprobs`.
	Logprobs []TranscriptionTextDoneEventLogprob `json:"logprobs"`
	// Usage statistics for models billed by token usage.
	Usage TranscriptionTextDoneEventUsage `json:"usage"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Text        respjson.Field
		Type        respjson.Field
		Logprobs    respjson.Field
		Usage       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r TranscriptionTextDoneEvent) RawJSON() string { return r.JSON.raw }
func (r *TranscriptionTextDoneEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type TranscriptionTextDoneEventLogprob struct {
	// The token that was used to generate the log probability.
	Token string `json:"token"`
	// The bytes that were used to generate the log probability.
	Bytes []int64 `json:"bytes"`
	// The log probability of the token.
	Logprob float64 `json:"logprob"`
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
func (r TranscriptionTextDoneEventLogprob) RawJSON() string { return r.JSON.raw }
func (r *TranscriptionTextDoneEventLogprob) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Usage statistics for models billed by token usage.
type TranscriptionTextDoneEventUsage struct {
	// Number of input tokens billed for this request.
	InputTokens int64 `json:"input_tokens,required"`
	// Number of output tokens generated.
	OutputTokens int64 `json:"output_tokens,required"`
	// Total number of tokens used (input + output).
	TotalTokens int64 `json:"total_tokens,required"`
	// The type of the usage object. Always `tokens` for this variant.
	Type constant.Tokens `json:"type,required"`
	// Details about the input tokens billed for this request.
	InputTokenDetails TranscriptionTextDoneEventUsageInputTokenDetails `json:"input_token_details"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		InputTokens       respjson.Field
		OutputTokens      respjson.Field
		TotalTokens       respjson.Field
		Type              respjson.Field
		InputTokenDetails respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r TranscriptionTextDoneEventUsage) RawJSON() string { return r.JSON.raw }
func (r *TranscriptionTextDoneEventUsage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Details about the input tokens billed for this request.
type TranscriptionTextDoneEventUsageInputTokenDetails struct {
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
func (r TranscriptionTextDoneEventUsageInputTokenDetails) RawJSON() string { return r.JSON.raw }
func (r *TranscriptionTextDoneEventUsageInputTokenDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Emitted when a diarized transcription returns a completed segment with speaker
// information. Only emitted when you
// [create a transcription](https://platform.openai.com/docs/api-reference/audio/create-transcription)
// with `stream` set to `true` and `response_format` set to `diarized_json`.
type TranscriptionTextSegmentEvent struct {
	// Unique identifier for the segment.
	ID string `json:"id,required"`
	// End timestamp of the segment in seconds.
	End float64 `json:"end,required"`
	// Speaker label for this segment.
	Speaker string `json:"speaker,required"`
	// Start timestamp of the segment in seconds.
	Start float64 `json:"start,required"`
	// Transcript text for this segment.
	Text string `json:"text,required"`
	// The type of the event. Always `transcript.text.segment`.
	Type constant.TranscriptTextSegment `json:"type,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		End         respjson.Field
		Speaker     respjson.Field
		Start       respjson.Field
		Text        respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r TranscriptionTextSegmentEvent) RawJSON() string { return r.JSON.raw }
func (r *TranscriptionTextSegmentEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Represents a verbose json transcription response returned by model, based on the
// provided input.
type TranscriptionVerbose struct {
	// The duration of the input audio.
	Duration float64 `json:"duration,required"`
	// The language of the input audio.
	Language string `json:"language,required"`
	// The transcribed text.
	Text string `json:"text,required"`
	// Segments of the transcribed text and their corresponding details.
	Segments []TranscriptionSegment `json:"segments"`
	// Usage statistics for models billed by audio input duration.
	Usage TranscriptionVerboseUsage `json:"usage"`
	// Extracted words and their corresponding timestamps.
	Words []TranscriptionWord `json:"words"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Duration    respjson.Field
		Language    respjson.Field
		Text        respjson.Field
		Segments    respjson.Field
		Usage       respjson.Field
		Words       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r TranscriptionVerbose) RawJSON() string { return r.JSON.raw }
func (r *TranscriptionVerbose) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Usage statistics for models billed by audio input duration.
type TranscriptionVerboseUsage struct {
	// Duration of the input audio in seconds.
	Seconds float64 `json:"seconds,required"`
	// The type of the usage object. Always `duration` for this variant.
	Type constant.Duration `json:"type,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Seconds     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r TranscriptionVerboseUsage) RawJSON() string { return r.JSON.raw }
func (r *TranscriptionVerboseUsage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type TranscriptionWord struct {
	// End time of the word in seconds.
	End float64 `json:"end,required"`
	// Start time of the word in seconds.
	Start float64 `json:"start,required"`
	// The text content of the word.
	Word string `json:"word,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		End         respjson.Field
		Start       respjson.Field
		Word        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r TranscriptionWord) RawJSON() string { return r.JSON.raw }
func (r *TranscriptionWord) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// AudioTranscriptionNewResponseUnion contains all possible properties and values
// from [Transcription], [TranscriptionVerbose].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type AudioTranscriptionNewResponseUnion struct {
	Text string `json:"text"`
	// This field is from variant [Transcription].
	Logprobs []TranscriptionLogprob `json:"logprobs"`
	// This field is a union of [TranscriptionUsageUnion], [TranscriptionVerboseUsage]
	Usage AudioTranscriptionNewResponseUnionUsage `json:"usage"`
	// This field is from variant [TranscriptionVerbose].
	Duration float64 `json:"duration"`
	// This field is from variant [TranscriptionVerbose].
	Language string `json:"language"`
	// This field is from variant [TranscriptionVerbose].
	Segments []TranscriptionSegment `json:"segments"`
	// This field is from variant [TranscriptionVerbose].
	Words []TranscriptionWord `json:"words"`
	JSON  struct {
		Text     respjson.Field
		Logprobs respjson.Field
		Usage    respjson.Field
		Duration respjson.Field
		Language respjson.Field
		Segments respjson.Field
		Words    respjson.Field
		raw      string
	} `json:"-"`
}

func (u AudioTranscriptionNewResponseUnion) AsTranscription() (v Transcription) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AudioTranscriptionNewResponseUnion) AsTranscriptionVerbose() (v TranscriptionVerbose) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u AudioTranscriptionNewResponseUnion) RawJSON() string { return u.JSON.raw }

func (r *AudioTranscriptionNewResponseUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// AudioTranscriptionNewResponseUnionUsage is an implicit subunion of
// [AudioTranscriptionNewResponseUnion]. AudioTranscriptionNewResponseUnionUsage
// provides convenient access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [AudioTranscriptionNewResponseUnion].
type AudioTranscriptionNewResponseUnionUsage struct {
	// This field is from variant [TranscriptionUsageUnion].
	InputTokens int64 `json:"input_tokens"`
	// This field is from variant [TranscriptionUsageUnion].
	OutputTokens int64 `json:"output_tokens"`
	// This field is from variant [TranscriptionUsageUnion].
	TotalTokens int64  `json:"total_tokens"`
	Type        string `json:"type"`
	// This field is from variant [TranscriptionUsageUnion].
	InputTokenDetails TranscriptionUsageTokensInputTokenDetails `json:"input_token_details"`
	Seconds           float64                                   `json:"seconds"`
	JSON              struct {
		InputTokens       respjson.Field
		OutputTokens      respjson.Field
		TotalTokens       respjson.Field
		Type              respjson.Field
		InputTokenDetails respjson.Field
		Seconds           respjson.Field
		raw               string
	} `json:"-"`
}

func (r *AudioTranscriptionNewResponseUnionUsage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AudioTranscriptionNewParams struct {
	// The audio file object (not file name) to transcribe, in one of these formats:
	// flac, mp3, mp4, mpeg, mpga, m4a, ogg, wav, or webm.
	File io.Reader `json:"file,omitzero,required" format:"binary"`
	// ID of the model to use. The options are `gpt-4o-transcribe`,
	// `gpt-4o-mini-transcribe`, `gpt-4o-mini-transcribe-2025-12-15`, `whisper-1`
	// (which is powered by our open source Whisper V2 model), and
	// `gpt-4o-transcribe-diarize`.
	Model AudioModel `json:"model,omitzero,required"`
	// The language of the input audio. Supplying the input language in
	// [ISO-639-1](https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes) (e.g. `en`)
	// format will improve accuracy and latency.
	Language param.Opt[string] `json:"language,omitzero"`
	// An optional text to guide the model's style or continue a previous audio
	// segment. The
	// [prompt](https://platform.openai.com/docs/guides/speech-to-text#prompting)
	// should match the audio language. This field is not supported when using
	// `gpt-4o-transcribe-diarize`.
	Prompt param.Opt[string] `json:"prompt,omitzero"`
	// The sampling temperature, between 0 and 1. Higher values like 0.8 will make the
	// output more random, while lower values like 0.2 will make it more focused and
	// deterministic. If set to 0, the model will use
	// [log probability](https://en.wikipedia.org/wiki/Log_probability) to
	// automatically increase the temperature until certain thresholds are hit.
	Temperature param.Opt[float64] `json:"temperature,omitzero"`
	// Controls how the audio is cut into chunks. When set to `"auto"`, the server
	// first normalizes loudness and then uses voice activity detection (VAD) to choose
	// boundaries. `server_vad` object can be provided to tweak VAD detection
	// parameters manually. If unset, the audio is transcribed as a single block.
	// Required when using `gpt-4o-transcribe-diarize` for inputs longer than 30
	// seconds.
	ChunkingStrategy AudioTranscriptionNewParamsChunkingStrategyUnion `json:"chunking_strategy,omitzero"`
	// Additional information to include in the transcription response. `logprobs` will
	// return the log probabilities of the tokens in the response to understand the
	// model's confidence in the transcription. `logprobs` only works with
	// response_format set to `json` and only with the models `gpt-4o-transcribe`,
	// `gpt-4o-mini-transcribe`, and `gpt-4o-mini-transcribe-2025-12-15`. This field is
	// not supported when using `gpt-4o-transcribe-diarize`.
	Include []TranscriptionInclude `json:"include,omitzero"`
	// Optional list of speaker names that correspond to the audio samples provided in
	// `known_speaker_references[]`. Each entry should be a short identifier (for
	// example `customer` or `agent`). Up to 4 speakers are supported.
	KnownSpeakerNames []string `json:"known_speaker_names,omitzero"`
	// Optional list of audio samples (as
	// [data URLs](https://developer.mozilla.org/en-US/docs/Web/HTTP/Basics_of_HTTP/Data_URLs))
	// that contain known speaker references matching `known_speaker_names[]`. Each
	// sample must be between 2 and 10 seconds, and can use any of the same input audio
	// formats supported by `file`.
	KnownSpeakerReferences []string `json:"known_speaker_references,omitzero"`
	// The format of the output, in one of these options: `json`, `text`, `srt`,
	// `verbose_json`, `vtt`, or `diarized_json`. For `gpt-4o-transcribe` and
	// `gpt-4o-mini-transcribe`, the only supported format is `json`. For
	// `gpt-4o-transcribe-diarize`, the supported formats are `json`, `text`, and
	// `diarized_json`, with `diarized_json` required to receive speaker annotations.
	//
	// Any of "json", "text", "srt", "verbose_json", "vtt", "diarized_json".
	ResponseFormat AudioResponseFormat `json:"response_format,omitzero"`
	// The timestamp granularities to populate for this transcription.
	// `response_format` must be set `verbose_json` to use timestamp granularities.
	// Either or both of these options are supported: `word`, or `segment`. Note: There
	// is no additional latency for segment timestamps, but generating word timestamps
	// incurs additional latency. This option is not available for
	// `gpt-4o-transcribe-diarize`.
	//
	// Any of "word", "segment".
	TimestampGranularities []string `json:"timestamp_granularities,omitzero"`
	paramObj
}

func (r AudioTranscriptionNewParams) MarshalMultipart() (data []byte, contentType string, err error) {
	buf := bytes.NewBuffer(nil)
	writer := multipart.NewWriter(buf)
	err = apiform.MarshalRoot(r, writer)
	if err == nil {
		err = apiform.WriteExtras(writer, r.ExtraFields())
	}
	if err != nil {
		writer.Close()
		return nil, "", err
	}
	err = writer.Close()
	if err != nil {
		return nil, "", err
	}
	return buf.Bytes(), writer.FormDataContentType(), nil
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type AudioTranscriptionNewParamsChunkingStrategyUnion struct {
	// Construct this variant with constant.ValueOf[constant.Auto]()
	OfAuto                                            constant.Auto                                         `json:",omitzero,inline"`
	OfAudioTranscriptionNewsChunkingStrategyVadConfig *AudioTranscriptionNewParamsChunkingStrategyVadConfig `json:",omitzero,inline"`
	paramUnion
}

func (u AudioTranscriptionNewParamsChunkingStrategyUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfAuto, u.OfAudioTranscriptionNewsChunkingStrategyVadConfig)
}
func (u *AudioTranscriptionNewParamsChunkingStrategyUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *AudioTranscriptionNewParamsChunkingStrategyUnion) asAny() any {
	if !param.IsOmitted(u.OfAuto) {
		return &u.OfAuto
	} else if !param.IsOmitted(u.OfAudioTranscriptionNewsChunkingStrategyVadConfig) {
		return u.OfAudioTranscriptionNewsChunkingStrategyVadConfig
	}
	return nil
}

// The property Type is required.
type AudioTranscriptionNewParamsChunkingStrategyVadConfig struct {
	// Must be set to `server_vad` to enable manual chunking using server side VAD.
	//
	// Any of "server_vad".
	Type string `json:"type,omitzero,required"`
	// Amount of audio to include before the VAD detected speech (in milliseconds).
	PrefixPaddingMs param.Opt[int64] `json:"prefix_padding_ms,omitzero"`
	// Duration of silence to detect speech stop (in milliseconds). With shorter values
	// the model will respond more quickly, but may jump in on short pauses from the
	// user.
	SilenceDurationMs param.Opt[int64] `json:"silence_duration_ms,omitzero"`
	// Sensitivity threshold (0.0 to 1.0) for voice activity detection. A higher
	// threshold will require louder audio to activate the model, and thus might
	// perform better in noisy environments.
	Threshold param.Opt[float64] `json:"threshold,omitzero"`
	paramObj
}

func (r AudioTranscriptionNewParamsChunkingStrategyVadConfig) MarshalJSON() (data []byte, err error) {
	type shadow AudioTranscriptionNewParamsChunkingStrategyVadConfig
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AudioTranscriptionNewParamsChunkingStrategyVadConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[AudioTranscriptionNewParamsChunkingStrategyVadConfig](
		"type", "server_vad",
	)
}

// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"reflect"

	"github.com/openai/openai-go/internal/apiform"
	"github.com/openai/openai-go/internal/apijson"
	"github.com/openai/openai-go/internal/param"
	"github.com/openai/openai-go/internal/requestconfig"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/ssestream"
	"github.com/tidwall/gjson"
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
func NewAudioTranscriptionService(opts ...option.RequestOption) (r *AudioTranscriptionService) {
	r = &AudioTranscriptionService{}
	r.Options = opts
	return
}

// Transcribes audio into the input language.
func (r *AudioTranscriptionService) New(ctx context.Context, body AudioTranscriptionNewParams, opts ...option.RequestOption) (res *Transcription, err error) {
	opts = append(r.Options[:], opts...)
	path := "audio/transcriptions"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Transcribes audio into the input language.
func (r *AudioTranscriptionService) NewStreaming(ctx context.Context, body AudioTranscriptionNewParams, opts ...option.RequestOption) (stream *ssestream.Stream[TranscriptionStreamEvent]) {
	var (
		raw *http.Response
		err error
	)
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithJSONSet("stream", true)}, opts...)
	path := "audio/transcriptions"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &raw, opts...)
	return ssestream.NewStream[TranscriptionStreamEvent](ssestream.NewDecoder(raw), err)
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
	JSON     transcriptionJSON      `json:"-"`
}

// transcriptionJSON contains the JSON metadata for the struct [Transcription]
type transcriptionJSON struct {
	Text        apijson.Field
	Logprobs    apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *Transcription) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r transcriptionJSON) RawJSON() string {
	return r.raw
}

type TranscriptionLogprob struct {
	// The token in the transcription.
	Token string `json:"token"`
	// The bytes of the token.
	Bytes []float64 `json:"bytes"`
	// The log probability of the token.
	Logprob float64                  `json:"logprob"`
	JSON    transcriptionLogprobJSON `json:"-"`
}

// transcriptionLogprobJSON contains the JSON metadata for the struct
// [TranscriptionLogprob]
type transcriptionLogprobJSON struct {
	Token       apijson.Field
	Bytes       apijson.Field
	Logprob     apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *TranscriptionLogprob) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r transcriptionLogprobJSON) RawJSON() string {
	return r.raw
}

type TranscriptionInclude string

const (
	TranscriptionIncludeLogprobs TranscriptionInclude = "logprobs"
)

func (r TranscriptionInclude) IsKnown() bool {
	switch r {
	case TranscriptionIncludeLogprobs:
		return true
	}
	return false
}

// Emitted when there is an additional text delta. This is also the first event
// emitted when the transcription starts. Only emitted when you
// [create a transcription](https://platform.openai.com/docs/api-reference/audio/create-transcription)
// with the `Stream` parameter set to `true`.
type TranscriptionStreamEvent struct {
	// The type of the event. Always `transcript.text.delta`.
	Type TranscriptionStreamEventType `json:"type,required"`
	// The text delta that was additionally transcribed.
	Delta string `json:"delta"`
	// This field can have the runtime type of [[]TranscriptionTextDeltaEventLogprob],
	// [[]TranscriptionTextDoneEventLogprob].
	Logprobs interface{} `json:"logprobs"`
	// The text that was transcribed.
	Text  string                       `json:"text"`
	JSON  transcriptionStreamEventJSON `json:"-"`
	union TranscriptionStreamEventUnion
}

// transcriptionStreamEventJSON contains the JSON metadata for the struct
// [TranscriptionStreamEvent]
type transcriptionStreamEventJSON struct {
	Type        apijson.Field
	Delta       apijson.Field
	Logprobs    apijson.Field
	Text        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r transcriptionStreamEventJSON) RawJSON() string {
	return r.raw
}

func (r *TranscriptionStreamEvent) UnmarshalJSON(data []byte) (err error) {
	*r = TranscriptionStreamEvent{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [TranscriptionStreamEventUnion] interface which you can cast
// to the specific types for more type safety.
//
// Possible runtime types of the union are [TranscriptionTextDeltaEvent],
// [TranscriptionTextDoneEvent].
func (r TranscriptionStreamEvent) AsUnion() TranscriptionStreamEventUnion {
	return r.union
}

// Emitted when there is an additional text delta. This is also the first event
// emitted when the transcription starts. Only emitted when you
// [create a transcription](https://platform.openai.com/docs/api-reference/audio/create-transcription)
// with the `Stream` parameter set to `true`.
//
// Union satisfied by [TranscriptionTextDeltaEvent] or
// [TranscriptionTextDoneEvent].
type TranscriptionStreamEventUnion interface {
	implementsTranscriptionStreamEvent()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*TranscriptionStreamEventUnion)(nil)).Elem(),
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(TranscriptionTextDeltaEvent{}),
			DiscriminatorValue: "transcript.text.delta",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(TranscriptionTextDoneEvent{}),
			DiscriminatorValue: "transcript.text.done",
		},
	)
}

// The type of the event. Always `transcript.text.delta`.
type TranscriptionStreamEventType string

const (
	TranscriptionStreamEventTypeTranscriptTextDelta TranscriptionStreamEventType = "transcript.text.delta"
	TranscriptionStreamEventTypeTranscriptTextDone  TranscriptionStreamEventType = "transcript.text.done"
)

func (r TranscriptionStreamEventType) IsKnown() bool {
	switch r {
	case TranscriptionStreamEventTypeTranscriptTextDelta, TranscriptionStreamEventTypeTranscriptTextDone:
		return true
	}
	return false
}

// Emitted when there is an additional text delta. This is also the first event
// emitted when the transcription starts. Only emitted when you
// [create a transcription](https://platform.openai.com/docs/api-reference/audio/create-transcription)
// with the `Stream` parameter set to `true`.
type TranscriptionTextDeltaEvent struct {
	// The text delta that was additionally transcribed.
	Delta string `json:"delta,required"`
	// The type of the event. Always `transcript.text.delta`.
	Type TranscriptionTextDeltaEventType `json:"type,required"`
	// The log probabilities of the delta. Only included if you
	// [create a transcription](https://platform.openai.com/docs/api-reference/audio/create-transcription)
	// with the `include[]` parameter set to `logprobs`.
	Logprobs []TranscriptionTextDeltaEventLogprob `json:"logprobs"`
	JSON     transcriptionTextDeltaEventJSON      `json:"-"`
}

// transcriptionTextDeltaEventJSON contains the JSON metadata for the struct
// [TranscriptionTextDeltaEvent]
type transcriptionTextDeltaEventJSON struct {
	Delta       apijson.Field
	Type        apijson.Field
	Logprobs    apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *TranscriptionTextDeltaEvent) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r transcriptionTextDeltaEventJSON) RawJSON() string {
	return r.raw
}

func (r TranscriptionTextDeltaEvent) implementsTranscriptionStreamEvent() {}

// The type of the event. Always `transcript.text.delta`.
type TranscriptionTextDeltaEventType string

const (
	TranscriptionTextDeltaEventTypeTranscriptTextDelta TranscriptionTextDeltaEventType = "transcript.text.delta"
)

func (r TranscriptionTextDeltaEventType) IsKnown() bool {
	switch r {
	case TranscriptionTextDeltaEventTypeTranscriptTextDelta:
		return true
	}
	return false
}

type TranscriptionTextDeltaEventLogprob struct {
	// The token that was used to generate the log probability.
	Token string `json:"token"`
	// The bytes that were used to generate the log probability.
	Bytes []interface{} `json:"bytes"`
	// The log probability of the token.
	Logprob float64                                `json:"logprob"`
	JSON    transcriptionTextDeltaEventLogprobJSON `json:"-"`
}

// transcriptionTextDeltaEventLogprobJSON contains the JSON metadata for the struct
// [TranscriptionTextDeltaEventLogprob]
type transcriptionTextDeltaEventLogprobJSON struct {
	Token       apijson.Field
	Bytes       apijson.Field
	Logprob     apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *TranscriptionTextDeltaEventLogprob) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r transcriptionTextDeltaEventLogprobJSON) RawJSON() string {
	return r.raw
}

// Emitted when the transcription is complete. Contains the complete transcription
// text. Only emitted when you
// [create a transcription](https://platform.openai.com/docs/api-reference/audio/create-transcription)
// with the `Stream` parameter set to `true`.
type TranscriptionTextDoneEvent struct {
	// The text that was transcribed.
	Text string `json:"text,required"`
	// The type of the event. Always `transcript.text.done`.
	Type TranscriptionTextDoneEventType `json:"type,required"`
	// The log probabilities of the individual tokens in the transcription. Only
	// included if you
	// [create a transcription](https://platform.openai.com/docs/api-reference/audio/create-transcription)
	// with the `include[]` parameter set to `logprobs`.
	Logprobs []TranscriptionTextDoneEventLogprob `json:"logprobs"`
	JSON     transcriptionTextDoneEventJSON      `json:"-"`
}

// transcriptionTextDoneEventJSON contains the JSON metadata for the struct
// [TranscriptionTextDoneEvent]
type transcriptionTextDoneEventJSON struct {
	Text        apijson.Field
	Type        apijson.Field
	Logprobs    apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *TranscriptionTextDoneEvent) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r transcriptionTextDoneEventJSON) RawJSON() string {
	return r.raw
}

func (r TranscriptionTextDoneEvent) implementsTranscriptionStreamEvent() {}

// The type of the event. Always `transcript.text.done`.
type TranscriptionTextDoneEventType string

const (
	TranscriptionTextDoneEventTypeTranscriptTextDone TranscriptionTextDoneEventType = "transcript.text.done"
)

func (r TranscriptionTextDoneEventType) IsKnown() bool {
	switch r {
	case TranscriptionTextDoneEventTypeTranscriptTextDone:
		return true
	}
	return false
}

type TranscriptionTextDoneEventLogprob struct {
	// The token that was used to generate the log probability.
	Token string `json:"token"`
	// The bytes that were used to generate the log probability.
	Bytes []interface{} `json:"bytes"`
	// The log probability of the token.
	Logprob float64                               `json:"logprob"`
	JSON    transcriptionTextDoneEventLogprobJSON `json:"-"`
}

// transcriptionTextDoneEventLogprobJSON contains the JSON metadata for the struct
// [TranscriptionTextDoneEventLogprob]
type transcriptionTextDoneEventLogprobJSON struct {
	Token       apijson.Field
	Bytes       apijson.Field
	Logprob     apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *TranscriptionTextDoneEventLogprob) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r transcriptionTextDoneEventLogprobJSON) RawJSON() string {
	return r.raw
}

type AudioTranscriptionNewParams struct {
	// The audio file object (not file name) to transcribe, in one of these formats:
	// flac, mp3, mp4, mpeg, mpga, m4a, ogg, wav, or webm.
	File param.Field[io.Reader] `json:"file,required" format:"binary"`
	// ID of the model to use. The options are `gpt-4o-transcribe`,
	// `gpt-4o-mini-transcribe`, and `whisper-1` (which is powered by our open source
	// Whisper V2 model).
	Model param.Field[AudioModel] `json:"model,required"`
	// Additional information to include in the transcription response. `logprobs` will
	// return the log probabilities of the tokens in the response to understand the
	// model's confidence in the transcription. `logprobs` only works with
	// response_format set to `json` and only with the models `gpt-4o-transcribe` and
	// `gpt-4o-mini-transcribe`.
	Include param.Field[[]TranscriptionInclude] `json:"include"`
	// The language of the input audio. Supplying the input language in
	// [ISO-639-1](https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes) (e.g. `en`)
	// format will improve accuracy and latency.
	Language param.Field[string] `json:"language"`
	// An optional text to guide the model's style or continue a previous audio
	// segment. The
	// [prompt](https://platform.openai.com/docs/guides/speech-to-text#prompting)
	// should match the audio language.
	Prompt param.Field[string] `json:"prompt"`
	// The format of the output, in one of these options: `json`, `text`, `srt`,
	// `verbose_json`, or `vtt`. For `gpt-4o-transcribe` and `gpt-4o-mini-transcribe`,
	// the only supported format is `json`.
	ResponseFormat param.Field[AudioResponseFormat] `json:"response_format"`
	// The sampling temperature, between 0 and 1. Higher values like 0.8 will make the
	// output more random, while lower values like 0.2 will make it more focused and
	// deterministic. If set to 0, the model will use
	// [log probability](https://en.wikipedia.org/wiki/Log_probability) to
	// automatically increase the temperature until certain thresholds are hit.
	Temperature param.Field[float64] `json:"temperature"`
	// The timestamp granularities to populate for this transcription.
	// `response_format` must be set `verbose_json` to use timestamp granularities.
	// Either or both of these options are supported: `word`, or `segment`. Note: There
	// is no additional latency for segment timestamps, but generating word timestamps
	// incurs additional latency.
	TimestampGranularities param.Field[[]AudioTranscriptionNewParamsTimestampGranularity] `json:"timestamp_granularities"`
}

func (r AudioTranscriptionNewParams) MarshalMultipart() (data []byte, contentType string, err error) {
	buf := bytes.NewBuffer(nil)
	writer := multipart.NewWriter(buf)
	err = apiform.MarshalRoot(r, writer)
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

type AudioTranscriptionNewParamsTimestampGranularity string

const (
	AudioTranscriptionNewParamsTimestampGranularityWord    AudioTranscriptionNewParamsTimestampGranularity = "word"
	AudioTranscriptionNewParamsTimestampGranularitySegment AudioTranscriptionNewParamsTimestampGranularity = "segment"
)

func (r AudioTranscriptionNewParamsTimestampGranularity) IsKnown() bool {
	switch r {
	case AudioTranscriptionNewParamsTimestampGranularityWord, AudioTranscriptionNewParamsTimestampGranularitySegment:
		return true
	}
	return false
}

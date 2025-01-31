// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/openai/openai-go/internal/apiform"
	"github.com/openai/openai-go/internal/apijson"
	"github.com/openai/openai-go/internal/param"
	"github.com/openai/openai-go/internal/requestconfig"
	"github.com/openai/openai-go/option"
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

// Represents a transcription response returned by model, based on the provided
// input.
type Transcription struct {
	// The transcribed text.
	Text string            `json:"text,required"`
	JSON transcriptionJSON `json:"-"`
}

// transcriptionJSON contains the JSON metadata for the struct [Transcription]
type transcriptionJSON struct {
	Text        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *Transcription) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r transcriptionJSON) RawJSON() string {
	return r.raw
}

type AudioTranscriptionNewParams struct {
	// The audio file object (not file name) to transcribe, in one of these formats:
	// flac, mp3, mp4, mpeg, mpga, m4a, ogg, wav, or webm.
	File param.Field[io.Reader] `json:"file,required" format:"binary"`
	// ID of the model to use. Only `whisper-1` (which is powered by our open source
	// Whisper V2 model) is currently available.
	Model param.Field[AudioModel] `json:"model,required"`
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
	// `verbose_json`, or `vtt`.
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

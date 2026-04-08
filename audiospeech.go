// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai

import (
	"context"
	"net/http"
	"slices"

	"github.com/openai/openai-go/v3/internal/apijson"
	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/param"
)

// Turn audio into text or text into audio.
//
// AudioSpeechService contains methods and other services that help with
// interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAudioSpeechService] method instead.
type AudioSpeechService struct {
	Options []option.RequestOption
}

// NewAudioSpeechService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewAudioSpeechService(opts ...option.RequestOption) (r AudioSpeechService) {
	r = AudioSpeechService{}
	r.Options = opts
	return
}

// Generates audio from the input text.
//
// Returns the audio file content, or a stream of audio events.
func (r *AudioSpeechService) New(ctx context.Context, body AudioSpeechNewParams, opts ...option.RequestOption) (res *http.Response, err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "application/octet-stream")}, opts...)
	path := "audio/speech"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

type SpeechModel = string

const (
	SpeechModelTTS1                   SpeechModel = "tts-1"
	SpeechModelTTS1HD                 SpeechModel = "tts-1-hd"
	SpeechModelGPT4oMiniTTS           SpeechModel = "gpt-4o-mini-tts"
	SpeechModelGPT4oMiniTTS2025_12_15 SpeechModel = "gpt-4o-mini-tts-2025-12-15"
)

type AudioSpeechNewParams struct {
	// The text to generate audio for. The maximum length is 4096 characters.
	Input string `json:"input" api:"required"`
	// One of the available [TTS models](https://platform.openai.com/docs/models#tts):
	// `tts-1`, `tts-1-hd`, `gpt-4o-mini-tts`, or `gpt-4o-mini-tts-2025-12-15`.
	Model SpeechModel `json:"model,omitzero" api:"required"`
	// The voice to use when generating the audio. Supported built-in voices are
	// `alloy`, `ash`, `ballad`, `coral`, `echo`, `fable`, `onyx`, `nova`, `sage`,
	// `shimmer`, `verse`, `marin`, and `cedar`. You may also provide a custom voice
	// object with an `id`, for example `{ "id": "voice_1234" }`. Previews of the
	// voices are available in the
	// [Text to speech guide](https://platform.openai.com/docs/guides/text-to-speech#voice-options).
	Voice AudioSpeechNewParamsVoiceUnion `json:"voice,omitzero" api:"required"`
	// Control the voice of your generated audio with additional instructions. Does not
	// work with `tts-1` or `tts-1-hd`.
	Instructions param.Opt[string] `json:"instructions,omitzero"`
	// The speed of the generated audio. Select a value from `0.25` to `4.0`. `1.0` is
	// the default.
	Speed param.Opt[float64] `json:"speed,omitzero"`
	// The format to audio in. Supported formats are `mp3`, `opus`, `aac`, `flac`,
	// `wav`, and `pcm`.
	//
	// Any of "mp3", "opus", "aac", "flac", "wav", "pcm".
	ResponseFormat AudioSpeechNewParamsResponseFormat `json:"response_format,omitzero"`
	// The format to stream the audio in. Supported formats are `sse` and `audio`.
	// `sse` is not supported for `tts-1` or `tts-1-hd`.
	//
	// Any of "sse", "audio".
	StreamFormat AudioSpeechNewParamsStreamFormat `json:"stream_format,omitzero"`
	paramObj
}

func (r AudioSpeechNewParams) MarshalJSON() (data []byte, err error) {
	type shadow AudioSpeechNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AudioSpeechNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type AudioSpeechNewParamsVoiceUnion struct {
	OfString param.Opt[string] `json:",omitzero,inline"`
	// Check if union is this variant with
	// !param.IsOmitted(union.OfAudioSpeechNewsVoiceString2)
	OfAudioSpeechNewsVoiceString2 param.Opt[string]            `json:",omitzero,inline"`
	OfAudioSpeechNewsVoiceID      *AudioSpeechNewParamsVoiceID `json:",omitzero,inline"`
	paramUnion
}

func (u AudioSpeechNewParamsVoiceUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfString, u.OfAudioSpeechNewsVoiceString2, u.OfAudioSpeechNewsVoiceID)
}
func (u *AudioSpeechNewParamsVoiceUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *AudioSpeechNewParamsVoiceUnion) asAny() any {
	if !param.IsOmitted(u.OfString) {
		return &u.OfString.Value
	} else if !param.IsOmitted(u.OfAudioSpeechNewsVoiceString2) {
		return &u.OfAudioSpeechNewsVoiceString2
	} else if !param.IsOmitted(u.OfAudioSpeechNewsVoiceID) {
		return u.OfAudioSpeechNewsVoiceID
	}
	return nil
}

type AudioSpeechNewParamsVoiceString2 string

const (
	AudioSpeechNewParamsVoiceString2Alloy   AudioSpeechNewParamsVoiceString2 = "alloy"
	AudioSpeechNewParamsVoiceString2Ash     AudioSpeechNewParamsVoiceString2 = "ash"
	AudioSpeechNewParamsVoiceString2Ballad  AudioSpeechNewParamsVoiceString2 = "ballad"
	AudioSpeechNewParamsVoiceString2Coral   AudioSpeechNewParamsVoiceString2 = "coral"
	AudioSpeechNewParamsVoiceString2Echo    AudioSpeechNewParamsVoiceString2 = "echo"
	AudioSpeechNewParamsVoiceString2Sage    AudioSpeechNewParamsVoiceString2 = "sage"
	AudioSpeechNewParamsVoiceString2Shimmer AudioSpeechNewParamsVoiceString2 = "shimmer"
	AudioSpeechNewParamsVoiceString2Verse   AudioSpeechNewParamsVoiceString2 = "verse"
	AudioSpeechNewParamsVoiceString2Marin   AudioSpeechNewParamsVoiceString2 = "marin"
	AudioSpeechNewParamsVoiceString2Cedar   AudioSpeechNewParamsVoiceString2 = "cedar"
)

// Custom voice reference.
//
// The property ID is required.
type AudioSpeechNewParamsVoiceID struct {
	// The custom voice ID, e.g. `voice_1234`.
	ID string `json:"id" api:"required"`
	paramObj
}

func (r AudioSpeechNewParamsVoiceID) MarshalJSON() (data []byte, err error) {
	type shadow AudioSpeechNewParamsVoiceID
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AudioSpeechNewParamsVoiceID) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The format to audio in. Supported formats are `mp3`, `opus`, `aac`, `flac`,
// `wav`, and `pcm`.
type AudioSpeechNewParamsResponseFormat string

const (
	AudioSpeechNewParamsResponseFormatMP3  AudioSpeechNewParamsResponseFormat = "mp3"
	AudioSpeechNewParamsResponseFormatOpus AudioSpeechNewParamsResponseFormat = "opus"
	AudioSpeechNewParamsResponseFormatAAC  AudioSpeechNewParamsResponseFormat = "aac"
	AudioSpeechNewParamsResponseFormatFLAC AudioSpeechNewParamsResponseFormat = "flac"
	AudioSpeechNewParamsResponseFormatWAV  AudioSpeechNewParamsResponseFormat = "wav"
	AudioSpeechNewParamsResponseFormatPCM  AudioSpeechNewParamsResponseFormat = "pcm"
)

// The format to stream the audio in. Supported formats are `sse` and `audio`.
// `sse` is not supported for `tts-1` or `tts-1-hd`.
type AudioSpeechNewParamsStreamFormat string

const (
	AudioSpeechNewParamsStreamFormatSSE   AudioSpeechNewParamsStreamFormat = "sse"
	AudioSpeechNewParamsStreamFormatAudio AudioSpeechNewParamsStreamFormat = "audio"
)

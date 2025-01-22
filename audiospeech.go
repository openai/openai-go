// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai

import (
	"context"
	"net/http"

	"github.com/openai/openai-go/internal/apijson"
	"github.com/openai/openai-go/internal/param"
	"github.com/openai/openai-go/internal/requestconfig"
	"github.com/openai/openai-go/option"
)

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
func NewAudioSpeechService(opts ...option.RequestOption) (r *AudioSpeechService) {
	r = &AudioSpeechService{}
	r.Options = opts
	return
}

// Generates audio from the input text.
func (r *AudioSpeechService) New(ctx context.Context, body AudioSpeechNewParams, opts ...option.RequestOption) (res *http.Response, err error) {
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "application/octet-stream")}, opts...)
	path := "audio/speech"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

type SpeechModel = string

const (
	SpeechModelTTS1   SpeechModel = "tts-1"
	SpeechModelTTS1HD SpeechModel = "tts-1-hd"
)

type AudioSpeechNewParams struct {
	// The text to generate audio for. The maximum length is 4096 characters.
	Input param.Field[string] `json:"input,required"`
	// One of the available [TTS models](https://platform.openai.com/docs/models#tts):
	// `tts-1` or `tts-1-hd`
	Model param.Field[SpeechModel] `json:"model,required"`
	// The voice to use when generating the audio. Supported voices are `alloy`, `ash`,
	// `coral`, `echo`, `fable`, `onyx`, `nova`, `sage` and `shimmer`. Previews of the
	// voices are available in the
	// [Text to speech guide](https://platform.openai.com/docs/guides/text-to-speech#voice-options).
	Voice param.Field[AudioSpeechNewParamsVoice] `json:"voice,required"`
	// The format to audio in. Supported formats are `mp3`, `opus`, `aac`, `flac`,
	// `wav`, and `pcm`.
	ResponseFormat param.Field[AudioSpeechNewParamsResponseFormat] `json:"response_format"`
	// The speed of the generated audio. Select a value from `0.25` to `4.0`. `1.0` is
	// the default.
	Speed param.Field[float64] `json:"speed"`
}

func (r AudioSpeechNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// The voice to use when generating the audio. Supported voices are `alloy`, `ash`,
// `coral`, `echo`, `fable`, `onyx`, `nova`, `sage` and `shimmer`. Previews of the
// voices are available in the
// [Text to speech guide](https://platform.openai.com/docs/guides/text-to-speech#voice-options).
type AudioSpeechNewParamsVoice string

const (
	AudioSpeechNewParamsVoiceAlloy   AudioSpeechNewParamsVoice = "alloy"
	AudioSpeechNewParamsVoiceAsh     AudioSpeechNewParamsVoice = "ash"
	AudioSpeechNewParamsVoiceCoral   AudioSpeechNewParamsVoice = "coral"
	AudioSpeechNewParamsVoiceEcho    AudioSpeechNewParamsVoice = "echo"
	AudioSpeechNewParamsVoiceFable   AudioSpeechNewParamsVoice = "fable"
	AudioSpeechNewParamsVoiceOnyx    AudioSpeechNewParamsVoice = "onyx"
	AudioSpeechNewParamsVoiceNova    AudioSpeechNewParamsVoice = "nova"
	AudioSpeechNewParamsVoiceSage    AudioSpeechNewParamsVoice = "sage"
	AudioSpeechNewParamsVoiceShimmer AudioSpeechNewParamsVoice = "shimmer"
)

func (r AudioSpeechNewParamsVoice) IsKnown() bool {
	switch r {
	case AudioSpeechNewParamsVoiceAlloy, AudioSpeechNewParamsVoiceAsh, AudioSpeechNewParamsVoiceCoral, AudioSpeechNewParamsVoiceEcho, AudioSpeechNewParamsVoiceFable, AudioSpeechNewParamsVoiceOnyx, AudioSpeechNewParamsVoiceNova, AudioSpeechNewParamsVoiceSage, AudioSpeechNewParamsVoiceShimmer:
		return true
	}
	return false
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

func (r AudioSpeechNewParamsResponseFormat) IsKnown() bool {
	switch r {
	case AudioSpeechNewParamsResponseFormatMP3, AudioSpeechNewParamsResponseFormatOpus, AudioSpeechNewParamsResponseFormatAAC, AudioSpeechNewParamsResponseFormatFLAC, AudioSpeechNewParamsResponseFormatWAV, AudioSpeechNewParamsResponseFormatPCM:
		return true
	}
	return false
}

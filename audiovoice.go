// File generated from our OpenAPI spec by Castiron. See CONTRIBUTING.md for details.

package openai

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"maps"
	"mime/multipart"
	"net/http"
	"slices"

	"github.com/openai/openai-go/v3/internal/apiform"
	"github.com/openai/openai-go/v3/internal/apijson"
	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/packages/respjson"
	"github.com/openai/openai-go/v3/shared/constant"
)

// Turn audio into text or text into audio.
//
// AudioVoiceService contains methods and other services that help with interacting
// with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAudioVoiceService] method instead.
type AudioVoiceService struct {
	Options []option.RequestOption
}

// NewAudioVoiceService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewAudioVoiceService(opts ...option.RequestOption) (r AudioVoiceService) {
	r = AudioVoiceService{}
	r.Options = requestconfig.InheritedOptions(opts...)
	return
}

// Creates a voice from a text prompt or from a consent recording and an audio
// sample.
//
// For prompt-based creation, send `type: "prompt"` with a `name` and `prompt` as
// JSON or multipart form data. For creation from an audio sample, send
// `type: "audio_sample"` with a `name`, `audio_sample`, and `consent` recording ID
// as multipart form data. The type defaults to `audio_sample` when omitted.
//
// Returns the saved voice's metadata. Voices created from text prompts are
// supported only in Live, not in Realtime or the speech endpoint. The response
// does not include preview audio.
func (r *AudioVoiceService) New(ctx context.Context, body AudioVoiceNewParams, opts ...option.RequestOption) (res *Voice, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	path := "audio/voices"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// A custom voice that can be used for audio output. Voices created from text
// prompts are supported only in Live.
type Voice struct {
	// The voice identifier, which can be referenced in API endpoints.
	ID string `json:"id" api:"required"`
	// The Unix timestamp (in seconds) for when the voice was created.
	CreatedAt int64 `json:"created_at" api:"required" format:"unixtime"`
	// The name of the voice.
	Name string `json:"name" api:"required"`
	// The object type, which is always `audio.voice`.
	Object constant.AudioVoice `json:"object" default:"audio.voice"`
	// How the voice was created. Voices created from text prompts are supported only
	// in Live.
	//
	// Any of "audio_sample", "prompt".
	Type VoiceType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		Name        respjson.Field
		Object      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r Voice) RawJSON() string { return r.JSON.raw }
func (r *Voice) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// How the voice was created. Voices created from text prompts are supported only
// in Live.
type VoiceType string

const (
	VoiceTypeAudioSample VoiceType = "audio_sample"
	VoiceTypePrompt      VoiceType = "prompt"
)

type AudioVoiceNewParams struct {

	//
	// Request body variants
	//

	// This field is a request body variant, only one variant field can be set.
	OfAudioSample *AudioVoiceNewParamsBodyAudioSample `json:",inline"`
	// This field is a request body variant, only one variant field can be set.
	OfPrompt *AudioVoiceNewParamsBodyPrompt `json:",inline"`

	paramObj
}

func (r *AudioVoiceNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (u AudioVoiceNewParams) MarshalMultipart() (data []byte, contentType string, err error) {
	// Overrides replace the entire request, including any already selected body.
	if _, overridden := u.Overrides(); overridden {
		data, err = param.MarshalUnion(u)
		return data, "application/json", err
	}
	populatedVariants := 0

	if u.OfAudioSample != nil {
		populatedVariants++
	}

	if u.OfPrompt != nil {
		populatedVariants++
	}

	if populatedVariants > 1 {
		return nil, "", fmt.Errorf("expected union to have only one present variant, got %d", populatedVariants)
	}

	var selected interface {
		apiform.Marshaler
		param.ParamStruct
		ExtraFields() map[string]any
	}
	if u.OfAudioSample != nil {
		selected = u.OfAudioSample
	} else if u.OfPrompt != nil {
		selected = u.OfPrompt
	}
	if selected != nil {
		if replacement, ok := selected.Overrides(); ok {
			data, err = param.MarshalWithExtras(selected, replacement, u.ExtraFields())
			return data, "application/json", err
		}
		if len(u.ExtraFields()) == 0 && len(selected.ExtraFields()) == 0 {
			return selected.MarshalMultipart()
		}
		// Copy metadata before applying root fields so callers can reuse a body.
		extras := make(map[string]any, len(selected.ExtraFields())+len(u.ExtraFields()))
		maps.Copy(extras, selected.ExtraFields())
		maps.Copy(extras, u.ExtraFields())
		buf := bytes.NewBuffer(nil)
		writer := multipart.NewWriter(buf)
		if err = apiform.MarshalEncodedRoot(selected, writer, extras, nil); err != nil {
			_ = writer.Close()
			return nil, "", err
		}
		if err = writer.Close(); err != nil {
			return nil, "", err
		}
		return buf.Bytes(), writer.FormDataContentType(), nil
	}

	if len(u.ExtraFields()) > 0 {
		data, err = param.MarshalObject(u, struct{}{})
		return data, "application/json", err
	}
	data, err = param.MarshalUnion(u)
	return data, "application/json", err
}

// Creates a voice from a consent recording and an audio sample. Requires
// multipart/form-data.
//
// The properties AudioSample, Consent, Name are required.
type AudioVoiceNewParamsBodyAudioSample struct {
	// The sample audio recording file. Maximum size is 10 MiB.
	//
	// Supported MIME types: `audio/mpeg`, `audio/wav`, `audio/x-wav`, `audio/ogg`,
	// `audio/aac`, `audio/flac`, `audio/webm`, `audio/mp4`.
	AudioSample io.Reader `json:"audio_sample" format:"binary" api:"required"`
	// The consent recording ID (for example, `cons_1234`).
	Consent string `json:"consent" api:"required"`
	// The name of the new voice.
	Name string `json:"name" api:"required"`
	// The voice creation method. Defaults to `audio_sample` when omitted.
	//
	// Any of "audio_sample".
	Type string `json:"type,omitzero"`
	paramObj
}

func (r AudioVoiceNewParamsBodyAudioSample) MarshalMultipart() (data []byte, contentType string, err error) {
	buf := bytes.NewBuffer(nil)
	writer := multipart.NewWriter(buf)
	err = apiform.MarshalRoot(r, writer)
	if err == nil {
		err = apiform.WriteExtras(writer, r.ExtraFields())
	}
	if err != nil {
		_ = writer.Close()
		return nil, "", err
	}
	err = writer.Close()
	if err != nil {
		return nil, "", err
	}
	return buf.Bytes(), writer.FormDataContentType(), nil
}

// MarshalMultipartTo writes multipart fields without buffering file contents.
// The caller owns writer and must close it to finish the multipart body.
func (r AudioVoiceNewParamsBodyAudioSample) MarshalMultipartTo(writer *multipart.Writer) error {
	if err := apiform.MarshalRoot(r, writer); err != nil {
		return err
	}
	return apiform.WriteExtras(writer, r.ExtraFields())
}

func init() {
	apijson.RegisterFieldValidator[AudioVoiceNewParamsBodyAudioSample](
		"type", "audio_sample",
	)
}

// Creates a synthetic voice from a text description. Supports application/json or
// multipart/form-data.
//
// The properties Name, Prompt, Type are required.
type AudioVoiceNewParamsBodyPrompt struct {
	// The name of the new voice.
	Name string `json:"name" api:"required"`
	// A description of the desired voice. Must not contain only whitespace.
	Prompt string `json:"prompt" api:"required"`
	// Set to `prompt` to create a voice from a text description.
	//
	// Any of "prompt".
	Type string `json:"type" api:"required"`
	// The voice creation model to use. Defaults to `auto`.
	//
	// Any of "auto", "2026-10-01".
	Model string `json:"model,omitzero"`
	// Optional text for the voice to speak during creation. If omitted, a script is
	// generated from the prompt. Must not be blank after trimming whitespace; scripts
	// that are too short are rejected.
	ScriptHint param.Opt[string] `json:"script_hint,omitzero"`
	paramObj
}

func (r AudioVoiceNewParamsBodyPrompt) MarshalMultipart() (data []byte, contentType string, err error) {
	buf := bytes.NewBuffer(nil)
	writer := multipart.NewWriter(buf)
	err = apiform.MarshalRoot(r, writer)
	if err == nil {
		err = apiform.WriteExtras(writer, r.ExtraFields())
	}
	if err != nil {
		_ = writer.Close()
		return nil, "", err
	}
	err = writer.Close()
	if err != nil {
		return nil, "", err
	}
	return buf.Bytes(), writer.FormDataContentType(), nil
}

// MarshalMultipartTo writes multipart fields without buffering file contents.
// The caller owns writer and must close it to finish the multipart body.
func (r AudioVoiceNewParamsBodyPrompt) MarshalMultipartTo(writer *multipart.Writer) error {
	if err := apiform.MarshalRoot(r, writer); err != nil {
		return err
	}
	return apiform.WriteExtras(writer, r.ExtraFields())
}

func init() {
	apijson.RegisterFieldValidator[AudioVoiceNewParamsBodyPrompt](
		"type", "prompt",
	)
	apijson.RegisterFieldValidator[AudioVoiceNewParamsBodyPrompt](
		"model", "auto", "2026-10-01",
	)
}

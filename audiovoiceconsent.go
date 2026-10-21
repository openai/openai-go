// File generated from our OpenAPI spec by Castiron. See CONTRIBUTING.md for details.

package openai

import (
	"bytes"
	"context"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"slices"

	"github.com/openai/openai-go/v3/internal/apiform"
	"github.com/openai/openai-go/v3/internal/apijson"
	"github.com/openai/openai-go/v3/internal/apiquery"
	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/pagination"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/packages/respjson"
	"github.com/openai/openai-go/v3/shared/constant"
)

// Turn audio into text or text into audio.
//
// AudioVoiceConsentService contains methods and other services that help with
// interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAudioVoiceConsentService] method instead.
type AudioVoiceConsentService struct {
	Options []option.RequestOption
}

// NewAudioVoiceConsentService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewAudioVoiceConsentService(opts ...option.RequestOption) (r AudioVoiceConsentService) {
	r = AudioVoiceConsentService{}
	r.Options = requestconfig.InheritedOptions(opts...)
	return
}

// Upload a voice consent recording.
func (r *AudioVoiceConsentService) New(ctx context.Context, body AudioVoiceConsentNewParams, opts ...option.RequestOption) (res *VoiceConsent, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	path := "audio/voice_consents"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, requestconfig.MultipartBody(body), &res, opts...)
	return res, err
}

// Retrieves a voice consent recording.
func (r *AudioVoiceConsentService) Get(ctx context.Context, consentID string, opts ...option.RequestOption) (res *VoiceConsent, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if consentID == "" {
		err = errors.New("missing required consent_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("audio/voice_consents/%s", consentID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Updates a voice consent recording (metadata only).
func (r *AudioVoiceConsentService) Update(ctx context.Context, consentID string, body AudioVoiceConsentUpdateParams, opts ...option.RequestOption) (res *VoiceConsent, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if consentID == "" {
		err = errors.New("missing required consent_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("audio/voice_consents/%s", consentID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Returns a list of voice consent recordings.
func (r *AudioVoiceConsentService) List(ctx context.Context, query AudioVoiceConsentListParams, opts ...option.RequestOption) (res *pagination.CursorPage[VoiceConsent], err error) {
	var raw *http.Response
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "audio/voice_consents"
	cfg, err := requestconfig.NewRequestConfig(ctx, http.MethodGet, path, query, &res, opts...)
	if err != nil {
		return nil, err
	}
	err = cfg.Execute()
	if err != nil {
		return nil, err
	}
	res.SetPageConfig(cfg, raw)
	return res, nil
}

// Returns a list of voice consent recordings.
func (r *AudioVoiceConsentService) ListAutoPaging(ctx context.Context, query AudioVoiceConsentListParams, opts ...option.RequestOption) *pagination.CursorPageAutoPager[VoiceConsent] {
	return pagination.NewCursorPageAutoPager(r.List(ctx, query, opts...))
}

// Deletes a voice consent recording.
func (r *AudioVoiceConsentService) Delete(ctx context.Context, consentID string, opts ...option.RequestOption) (res *VoiceConsentDeleted, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if consentID == "" {
		err = errors.New("missing required consent_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("audio/voice_consents/%s", consentID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// A consent recording used to authorize creation of a custom voice.
type VoiceConsent struct {
	// The consent recording identifier.
	ID string `json:"id" api:"required"`
	// The Unix timestamp (in seconds) for when the consent recording was created.
	CreatedAt int64 `json:"created_at" api:"required" format:"unixtime"`
	// The BCP 47 language tag for the consent phrase (for example, `en-US`).
	Language string `json:"language" api:"required"`
	// The label provided when the consent recording was uploaded.
	Name string `json:"name" api:"required"`
	// The object type, which is always `audio.voice_consent`.
	Object constant.AudioVoiceConsent `json:"object" default:"audio.voice_consent"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		Language    respjson.Field
		Name        respjson.Field
		Object      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r VoiceConsent) RawJSON() string { return r.JSON.raw }
func (r *VoiceConsent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type VoiceConsentDeleted struct {
	// The consent recording identifier.
	ID      string                     `json:"id" api:"required"`
	Deleted bool                       `json:"deleted" api:"required"`
	Object  constant.AudioVoiceConsent `json:"object" default:"audio.voice_consent"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Deleted     respjson.Field
		Object      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r VoiceConsentDeleted) RawJSON() string { return r.JSON.raw }
func (r *VoiceConsentDeleted) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AudioVoiceConsentNewParams struct {
	// The BCP 47 language tag for the consent phrase (for example, `en-US`).
	Language string `json:"language" api:"required"`
	// The label to use for this consent recording.
	Name string `json:"name" api:"required"`
	// The consent audio recording file. Maximum size is 10 MiB.
	//
	// Supported MIME types: `audio/mpeg`, `audio/wav`, `audio/x-wav`, `audio/ogg`,
	// `audio/aac`, `audio/flac`, `audio/webm`, `audio/mp4`.
	Recording io.Reader `json:"recording,omitzero" api:"required" format:"binary"`
	paramObj
}

func (r AudioVoiceConsentNewParams) MarshalMultipart() (data []byte, contentType string, err error) {
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
func (r AudioVoiceConsentNewParams) MarshalMultipartTo(writer *multipart.Writer) error {
	if err := apiform.MarshalRoot(r, writer); err != nil {
		return err
	}
	return apiform.WriteExtras(writer, r.ExtraFields())
}

type AudioVoiceConsentUpdateParams struct {
	// The updated label for this consent recording.
	Name string `json:"name" api:"required"`
	paramObj
}

func (r AudioVoiceConsentUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow AudioVoiceConsentUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AudioVoiceConsentUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AudioVoiceConsentListParams struct {
	// A cursor for use in pagination. `after` is an object ID that defines your place
	// in the list. For instance, if you make a list request and receive 100 objects,
	// ending with obj_foo, your subsequent call can include after=obj_foo in order to
	// fetch the next page of the list.
	After param.Opt[string] `query:"after,omitzero" json:"-"`
	// A limit on the number of objects to be returned. Limit can range between 1 and
	// 100, and the default is 20.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [AudioVoiceConsentListParams]'s query parameters as
// `url.Values`.
func (r AudioVoiceConsentListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai

import (
	"bytes"
	"context"
	"errors"
	"fmt"
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

// VideoService contains methods and other services that help with interacting with
// the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewVideoService] method instead.
type VideoService struct {
	Options []option.RequestOption
}

// NewVideoService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewVideoService(opts ...option.RequestOption) (r VideoService) {
	r = VideoService{}
	r.Options = opts
	return
}

// Create a new video generation job from a prompt and optional reference assets.
func (r *VideoService) New(ctx context.Context, body VideoNewParams, opts ...option.RequestOption) (res *Video, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "videos"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Create Video and Poll for Completion
//
// Polls the API and blocks until the task is complete.
// Default polling interval is 1 second.
func (r *VideoService) NewAndPoll(ctx context.Context, body VideoNewParams, pollIntervalMs int, opts ...option.RequestOption) (res *Video, err error) {
	video, err := r.New(ctx, body, opts...)
	if err != nil {
		return nil, err
	}
	return r.PollStatus(ctx, video.ID, pollIntervalMs, opts...)
}

// Fetch the latest metadata for a generated video.
func (r *VideoService) Get(ctx context.Context, videoID string, opts ...option.RequestOption) (res *Video, err error) {
	opts = slices.Concat(r.Options, opts)
	if videoID == "" {
		err = errors.New("missing required video_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("videos/%s", videoID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// List recently generated videos for the current project.
func (r *VideoService) List(ctx context.Context, query VideoListParams, opts ...option.RequestOption) (res *pagination.ConversationCursorPage[Video], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "videos"
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

// List recently generated videos for the current project.
func (r *VideoService) ListAutoPaging(ctx context.Context, query VideoListParams, opts ...option.RequestOption) *pagination.ConversationCursorPageAutoPager[Video] {
	return pagination.NewConversationCursorPageAutoPager(r.List(ctx, query, opts...))
}

// Permanently delete a completed or failed video and its stored assets.
func (r *VideoService) Delete(ctx context.Context, videoID string, opts ...option.RequestOption) (res *VideoDeleteResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if videoID == "" {
		err = errors.New("missing required video_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("videos/%s", videoID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Create a character from an uploaded video.
func (r *VideoService) NewCharacter(ctx context.Context, body VideoNewCharacterParams, opts ...option.RequestOption) (res *VideoNewCharacterResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "videos/characters"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Download the generated video bytes or a derived preview asset.
//
// Streams the rendered video content for the specified video job.
func (r *VideoService) DownloadContent(ctx context.Context, videoID string, query VideoDownloadContentParams, opts ...option.RequestOption) (res *http.Response, err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "application/binary")}, opts...)
	if videoID == "" {
		err = errors.New("missing required video_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("videos/%s/content", videoID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

// Create a new video generation job by editing a source video or existing
// generated video.
func (r *VideoService) Edit(ctx context.Context, body VideoEditParams, opts ...option.RequestOption) (res *Video, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "videos/edits"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Create an extension of a completed video.
func (r *VideoService) Extend(ctx context.Context, body VideoExtendParams, opts ...option.RequestOption) (res *Video, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "videos/extensions"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Fetch a character.
func (r *VideoService) GetCharacter(ctx context.Context, characterID string, opts ...option.RequestOption) (res *VideoGetCharacterResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if characterID == "" {
		err = errors.New("missing required character_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("videos/characters/%s", characterID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Create a remix of a completed video using a refreshed prompt.
func (r *VideoService) Remix(ctx context.Context, videoID string, body VideoRemixParams, opts ...option.RequestOption) (res *Video, err error) {
	opts = slices.Concat(r.Options, opts)
	if videoID == "" {
		err = errors.New("missing required video_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("videos/%s/remix", videoID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

type ImageInputReferenceParam struct {
	FileID param.Opt[string] `json:"file_id,omitzero"`
	// A fully qualified URL or base64-encoded data URL.
	ImageURL param.Opt[string] `json:"image_url,omitzero"`
	paramObj
}

func (r ImageInputReferenceParam) MarshalJSON() (data []byte, err error) {
	type shadow ImageInputReferenceParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ImageInputReferenceParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Structured information describing a generated video job.
type Video struct {
	// Unique identifier for the video job.
	ID string `json:"id" api:"required"`
	// Unix timestamp (seconds) for when the job completed, if finished.
	CompletedAt int64 `json:"completed_at" api:"required"`
	// Unix timestamp (seconds) for when the job was created.
	CreatedAt int64 `json:"created_at" api:"required"`
	// Error payload that explains why generation failed, if applicable.
	Error VideoCreateError `json:"error" api:"required"`
	// Unix timestamp (seconds) for when the downloadable assets expire, if set.
	ExpiresAt int64 `json:"expires_at" api:"required"`
	// The video generation model that produced the job.
	Model VideoModel `json:"model" api:"required"`
	// The object type, which is always `video`.
	Object constant.Video `json:"object" default:"video"`
	// Approximate completion percentage for the generation task.
	Progress int64 `json:"progress" api:"required"`
	// The prompt that was used to generate the video.
	Prompt string `json:"prompt" api:"required"`
	// Identifier of the source video if this video is a remix.
	RemixedFromVideoID string `json:"remixed_from_video_id" api:"required"`
	// Duration of the generated clip in seconds. For extensions, this is the stitched
	// total duration.
	Seconds VideoSeconds `json:"seconds" api:"required"`
	// The resolution of the generated video.
	//
	// Any of "720x1280", "1280x720", "1024x1792", "1792x1024".
	Size VideoSize `json:"size" api:"required"`
	// Current lifecycle status of the video job.
	//
	// Any of "queued", "in_progress", "completed", "failed".
	Status VideoStatus `json:"status" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                 respjson.Field
		CompletedAt        respjson.Field
		CreatedAt          respjson.Field
		Error              respjson.Field
		ExpiresAt          respjson.Field
		Model              respjson.Field
		Object             respjson.Field
		Progress           respjson.Field
		Prompt             respjson.Field
		RemixedFromVideoID respjson.Field
		Seconds            respjson.Field
		Size               respjson.Field
		Status             respjson.Field
		ExtraFields        map[string]respjson.Field
		raw                string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r Video) RawJSON() string { return r.JSON.raw }
func (r *Video) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Current lifecycle status of the video job.
type VideoStatus string

const (
	VideoStatusQueued     VideoStatus = "queued"
	VideoStatusInProgress VideoStatus = "in_progress"
	VideoStatusCompleted  VideoStatus = "completed"
	VideoStatusFailed     VideoStatus = "failed"
)

// An error that occurred while generating the response.
type VideoCreateError struct {
	// A machine-readable error code that was returned.
	Code string `json:"code" api:"required"`
	// A human-readable description of the error that was returned.
	Message string `json:"message" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Code        respjson.Field
		Message     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r VideoCreateError) RawJSON() string { return r.JSON.raw }
func (r *VideoCreateError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type VideoModel string

const (
	VideoModelSora2              VideoModel = "sora-2"
	VideoModelSora2Pro           VideoModel = "sora-2-pro"
	VideoModelSora2_2025_10_06   VideoModel = "sora-2-2025-10-06"
	VideoModelSora2Pro2025_10_06 VideoModel = "sora-2-pro-2025-10-06"
	VideoModelSora2_2025_12_08   VideoModel = "sora-2-2025-12-08"
)

type VideoSeconds string

const (
	VideoSeconds4  VideoSeconds = "4"
	VideoSeconds8  VideoSeconds = "8"
	VideoSeconds12 VideoSeconds = "12"
)

type VideoSize string

const (
	VideoSize720x1280  VideoSize = "720x1280"
	VideoSize1280x720  VideoSize = "1280x720"
	VideoSize1024x1792 VideoSize = "1024x1792"
	VideoSize1792x1024 VideoSize = "1792x1024"
)

// Confirmation payload returned after deleting a video.
type VideoDeleteResponse struct {
	// Identifier of the deleted video.
	ID string `json:"id" api:"required"`
	// Indicates that the video resource was deleted.
	Deleted bool `json:"deleted" api:"required"`
	// The object type that signals the deletion response.
	Object constant.VideoDeleted `json:"object" default:"video.deleted"`
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
func (r VideoDeleteResponse) RawJSON() string { return r.JSON.raw }
func (r *VideoDeleteResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type VideoNewCharacterResponse struct {
	// Identifier for the character creation cameo.
	ID string `json:"id" api:"required"`
	// Unix timestamp (in seconds) when the character was created.
	CreatedAt int64 `json:"created_at" api:"required"`
	// Display name for the character.
	Name string `json:"name" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		Name        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r VideoNewCharacterResponse) RawJSON() string { return r.JSON.raw }
func (r *VideoNewCharacterResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type VideoGetCharacterResponse struct {
	// Identifier for the character creation cameo.
	ID string `json:"id" api:"required"`
	// Unix timestamp (in seconds) when the character was created.
	CreatedAt int64 `json:"created_at" api:"required"`
	// Display name for the character.
	Name string `json:"name" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		Name        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r VideoGetCharacterResponse) RawJSON() string { return r.JSON.raw }
func (r *VideoGetCharacterResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type VideoNewParams struct {
	// Text prompt that describes the video to generate.
	Prompt string `json:"prompt" api:"required"`
	// Optional reference asset upload or reference object that guides generation.
	InputReference VideoNewParamsInputReferenceUnion `json:"input_reference,omitzero" format:"binary"`
	// The video generation model to use (allowed values: sora-2, sora-2-pro). Defaults
	// to `sora-2`.
	Model VideoModel `json:"model,omitzero"`
	// Clip duration in seconds (allowed values: 4, 8, 12). Defaults to 4 seconds.
	//
	// Any of "4", "8", "12".
	Seconds VideoSeconds `json:"seconds,omitzero"`
	// Output resolution formatted as width x height (allowed values: 720x1280,
	// 1280x720, 1024x1792, 1792x1024). Defaults to 720x1280.
	//
	// Any of "720x1280", "1280x720", "1024x1792", "1792x1024".
	Size VideoSize `json:"size,omitzero"`
	paramObj
}

func (r VideoNewParams) MarshalMultipart() (data []byte, contentType string, err error) {
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
type VideoNewParamsInputReferenceUnion struct {
	OfFile                io.Reader                 `json:",omitzero,inline"`
	OfImageInputReference *ImageInputReferenceParam `json:",omitzero,inline"`
	paramUnion
}

func (u VideoNewParamsInputReferenceUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfFile, u.OfImageInputReference)
}
func (u *VideoNewParamsInputReferenceUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *VideoNewParamsInputReferenceUnion) asAny() any {
	if !param.IsOmitted(u.OfFile) {
		return &u.OfFile
	} else if !param.IsOmitted(u.OfImageInputReference) {
		return u.OfImageInputReference
	}
	return nil
}

type VideoListParams struct {
	// Identifier for the last item from the previous pagination request
	After param.Opt[string] `query:"after,omitzero" json:"-"`
	// Number of items to retrieve
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Sort order of results by timestamp. Use `asc` for ascending order or `desc` for
	// descending order.
	//
	// Any of "asc", "desc".
	Order VideoListParamsOrder `query:"order,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [VideoListParams]'s query parameters as `url.Values`.
func (r VideoListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Sort order of results by timestamp. Use `asc` for ascending order or `desc` for
// descending order.
type VideoListParamsOrder string

const (
	VideoListParamsOrderAsc  VideoListParamsOrder = "asc"
	VideoListParamsOrderDesc VideoListParamsOrder = "desc"
)

type VideoNewCharacterParams struct {
	// Display name for this API character.
	Name string `json:"name" api:"required"`
	// Video file used to create a character.
	Video io.Reader `json:"video,omitzero" api:"required" format:"binary"`
	paramObj
}

func (r VideoNewCharacterParams) MarshalMultipart() (data []byte, contentType string, err error) {
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

type VideoDownloadContentParams struct {
	// Which downloadable asset to return. Defaults to the MP4 video.
	//
	// Any of "video", "thumbnail", "spritesheet".
	Variant VideoDownloadContentParamsVariant `query:"variant,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [VideoDownloadContentParams]'s query parameters as
// `url.Values`.
func (r VideoDownloadContentParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Which downloadable asset to return. Defaults to the MP4 video.
type VideoDownloadContentParamsVariant string

const (
	VideoDownloadContentParamsVariantVideo       VideoDownloadContentParamsVariant = "video"
	VideoDownloadContentParamsVariantThumbnail   VideoDownloadContentParamsVariant = "thumbnail"
	VideoDownloadContentParamsVariantSpritesheet VideoDownloadContentParamsVariant = "spritesheet"
)

type VideoEditParams struct {
	// Text prompt that describes how to edit the source video.
	Prompt string `json:"prompt" api:"required"`
	// Reference to the completed video to edit.
	Video VideoEditParamsVideoUnion `json:"video,omitzero" api:"required" format:"binary"`
	paramObj
}

func (r VideoEditParams) MarshalMultipart() (data []byte, contentType string, err error) {
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
type VideoEditParamsVideoUnion struct {
	OfFile                                    io.Reader                                     `json:",omitzero,inline"`
	OfVideoEditsVideoVideoReferenceInputParam *VideoEditParamsVideoVideoReferenceInputParam `json:",omitzero,inline"`
	paramUnion
}

func (u VideoEditParamsVideoUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfFile, u.OfVideoEditsVideoVideoReferenceInputParam)
}
func (u *VideoEditParamsVideoUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *VideoEditParamsVideoUnion) asAny() any {
	if !param.IsOmitted(u.OfFile) {
		return &u.OfFile
	} else if !param.IsOmitted(u.OfVideoEditsVideoVideoReferenceInputParam) {
		return u.OfVideoEditsVideoVideoReferenceInputParam
	}
	return nil
}

// Reference to the completed video.
//
// The property ID is required.
type VideoEditParamsVideoVideoReferenceInputParam struct {
	// The identifier of the completed video.
	ID string `json:"id" api:"required"`
	paramObj
}

func (r VideoEditParamsVideoVideoReferenceInputParam) MarshalJSON() (data []byte, err error) {
	type shadow VideoEditParamsVideoVideoReferenceInputParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *VideoEditParamsVideoVideoReferenceInputParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type VideoExtendParams struct {
	// Updated text prompt that directs the extension generation.
	Prompt string `json:"prompt" api:"required"`
	// Length of the newly generated extension segment in seconds (allowed values: 4,
	// 8, 12, 16, 20).
	//
	// Any of "4", "8", "12".
	Seconds VideoSeconds `json:"seconds,omitzero" api:"required"`
	// Reference to the completed video to extend.
	Video VideoExtendParamsVideoUnion `json:"video,omitzero" api:"required" format:"binary"`
	paramObj
}

func (r VideoExtendParams) MarshalMultipart() (data []byte, contentType string, err error) {
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
type VideoExtendParamsVideoUnion struct {
	OfFile                                      io.Reader                                       `json:",omitzero,inline"`
	OfVideoExtendsVideoVideoReferenceInputParam *VideoExtendParamsVideoVideoReferenceInputParam `json:",omitzero,inline"`
	paramUnion
}

func (u VideoExtendParamsVideoUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfFile, u.OfVideoExtendsVideoVideoReferenceInputParam)
}
func (u *VideoExtendParamsVideoUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *VideoExtendParamsVideoUnion) asAny() any {
	if !param.IsOmitted(u.OfFile) {
		return &u.OfFile
	} else if !param.IsOmitted(u.OfVideoExtendsVideoVideoReferenceInputParam) {
		return u.OfVideoExtendsVideoVideoReferenceInputParam
	}
	return nil
}

// Reference to the completed video.
//
// The property ID is required.
type VideoExtendParamsVideoVideoReferenceInputParam struct {
	// The identifier of the completed video.
	ID string `json:"id" api:"required"`
	paramObj
}

func (r VideoExtendParamsVideoVideoReferenceInputParam) MarshalJSON() (data []byte, err error) {
	type shadow VideoExtendParamsVideoVideoReferenceInputParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *VideoExtendParamsVideoVideoReferenceInputParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type VideoRemixParams struct {
	// Updated text prompt that directs the remix generation.
	Prompt string `json:"prompt" api:"required"`
	paramObj
}

func (r VideoRemixParams) MarshalJSON() (data []byte, err error) {
	type shadow VideoRemixParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *VideoRemixParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

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
	return
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
		return
	}
	path := fmt.Sprintf("videos/%s", videoID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
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
		return
	}
	path := fmt.Sprintf("videos/%s", videoID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return
}

// Download the generated video bytes or a derived preview asset.
//
// Streams the rendered video content for the specified video job.
func (r *VideoService) DownloadContent(ctx context.Context, videoID string, query VideoDownloadContentParams, opts ...option.RequestOption) (res *http.Response, err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "application/binary")}, opts...)
	if videoID == "" {
		err = errors.New("missing required video_id parameter")
		return
	}
	path := fmt.Sprintf("videos/%s/content", videoID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return
}

// Create a remix of a completed video using a refreshed prompt.
func (r *VideoService) Remix(ctx context.Context, videoID string, body VideoRemixParams, opts ...option.RequestOption) (res *Video, err error) {
	opts = slices.Concat(r.Options, opts)
	if videoID == "" {
		err = errors.New("missing required video_id parameter")
		return
	}
	path := fmt.Sprintf("videos/%s/remix", videoID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Structured information describing a generated video job.
type Video struct {
	// Unique identifier for the video job.
	ID string `json:"id,required"`
	// Unix timestamp (seconds) for when the job completed, if finished.
	CompletedAt int64 `json:"completed_at,required"`
	// Unix timestamp (seconds) for when the job was created.
	CreatedAt int64 `json:"created_at,required"`
	// Error payload that explains why generation failed, if applicable.
	Error VideoCreateError `json:"error,required"`
	// Unix timestamp (seconds) for when the downloadable assets expire, if set.
	ExpiresAt int64 `json:"expires_at,required"`
	// The video generation model that produced the job.
	Model VideoModel `json:"model,required"`
	// The object type, which is always `video`.
	Object constant.Video `json:"object,required"`
	// Approximate completion percentage for the generation task.
	Progress int64 `json:"progress,required"`
	// The prompt that was used to generate the video.
	Prompt string `json:"prompt,required"`
	// Identifier of the source video if this video is a remix.
	RemixedFromVideoID string `json:"remixed_from_video_id,required"`
	// Duration of the generated clip in seconds.
	//
	// Any of "4", "8", "12".
	Seconds VideoSeconds `json:"seconds,required"`
	// The resolution of the generated video.
	//
	// Any of "720x1280", "1280x720", "1024x1792", "1792x1024".
	Size VideoSize `json:"size,required"`
	// Current lifecycle status of the video job.
	//
	// Any of "queued", "in_progress", "completed", "failed".
	Status VideoStatus `json:"status,required"`
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
	Code string `json:"code,required"`
	// A human-readable description of the error that was returned.
	Message string `json:"message,required"`
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
	ID string `json:"id,required"`
	// Indicates that the video resource was deleted.
	Deleted bool `json:"deleted,required"`
	// The object type that signals the deletion response.
	Object constant.VideoDeleted `json:"object,required"`
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

type VideoNewParams struct {
	// Text prompt that describes the video to generate.
	Prompt string `json:"prompt,required"`
	// Optional image reference that guides generation.
	InputReference io.Reader `json:"input_reference,omitzero" format:"binary"`
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

type VideoRemixParams struct {
	// Updated text prompt that directs the remix generation.
	Prompt string `json:"prompt,required"`
	paramObj
}

func (r VideoRemixParams) MarshalJSON() (data []byte, err error) {
	type shadow VideoRemixParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *VideoRemixParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

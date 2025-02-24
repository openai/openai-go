// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/openai/openai-go/internal/apijson"
	"github.com/openai/openai-go/internal/apiquery"
	"github.com/openai/openai-go/internal/requestconfig"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/pagination"
	"github.com/openai/openai-go/packages/param"
	"github.com/openai/openai-go/packages/resp"
	"github.com/openai/openai-go/shared"
	"github.com/openai/openai-go/shared/constant"
)

// BetaThreadMessageService contains methods and other services that help with
// interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaThreadMessageService] method instead.
type BetaThreadMessageService struct {
	Options []option.RequestOption
}

// NewBetaThreadMessageService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewBetaThreadMessageService(opts ...option.RequestOption) (r BetaThreadMessageService) {
	r = BetaThreadMessageService{}
	r.Options = opts
	return
}

// Create a message.
func (r *BetaThreadMessageService) New(ctx context.Context, threadID string, body BetaThreadMessageNewParams, opts ...option.RequestOption) (res *Message, err error) {
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "assistants=v2")}, opts...)
	if threadID == "" {
		err = errors.New("missing required thread_id parameter")
		return
	}
	path := fmt.Sprintf("threads/%s/messages", threadID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Retrieve a message.
func (r *BetaThreadMessageService) Get(ctx context.Context, threadID string, messageID string, opts ...option.RequestOption) (res *Message, err error) {
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "assistants=v2")}, opts...)
	if threadID == "" {
		err = errors.New("missing required thread_id parameter")
		return
	}
	if messageID == "" {
		err = errors.New("missing required message_id parameter")
		return
	}
	path := fmt.Sprintf("threads/%s/messages/%s", threadID, messageID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// Modifies a message.
func (r *BetaThreadMessageService) Update(ctx context.Context, threadID string, messageID string, body BetaThreadMessageUpdateParams, opts ...option.RequestOption) (res *Message, err error) {
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "assistants=v2")}, opts...)
	if threadID == "" {
		err = errors.New("missing required thread_id parameter")
		return
	}
	if messageID == "" {
		err = errors.New("missing required message_id parameter")
		return
	}
	path := fmt.Sprintf("threads/%s/messages/%s", threadID, messageID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Returns a list of messages for a given thread.
func (r *BetaThreadMessageService) List(ctx context.Context, threadID string, query BetaThreadMessageListParams, opts ...option.RequestOption) (res *pagination.CursorPage[Message], err error) {
	var raw *http.Response
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "assistants=v2"), option.WithResponseInto(&raw)}, opts...)
	if threadID == "" {
		err = errors.New("missing required thread_id parameter")
		return
	}
	path := fmt.Sprintf("threads/%s/messages", threadID)
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

// Returns a list of messages for a given thread.
func (r *BetaThreadMessageService) ListAutoPaging(ctx context.Context, threadID string, query BetaThreadMessageListParams, opts ...option.RequestOption) *pagination.CursorPageAutoPager[Message] {
	return pagination.NewCursorPageAutoPager(r.List(ctx, threadID, query, opts...))
}

// Deletes a message.
func (r *BetaThreadMessageService) Delete(ctx context.Context, threadID string, messageID string, opts ...option.RequestOption) (res *MessageDeleted, err error) {
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "assistants=v2")}, opts...)
	if threadID == "" {
		err = errors.New("missing required thread_id parameter")
		return
	}
	if messageID == "" {
		err = errors.New("missing required message_id parameter")
		return
	}
	path := fmt.Sprintf("threads/%s/messages/%s", threadID, messageID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return
}

type AnnotationUnion struct {
	EndIndex     int64                              `json:"end_index"`
	FileCitation FileCitationAnnotationFileCitation `json:"file_citation"`
	StartIndex   int64                              `json:"start_index"`
	Text         string                             `json:"text"`
	Type         string                             `json:"type"`
	FilePath     FilePathAnnotationFilePath         `json:"file_path"`
	JSON         struct {
		EndIndex     resp.Field
		FileCitation resp.Field
		StartIndex   resp.Field
		Text         resp.Field
		Type         resp.Field
		FilePath     resp.Field
		raw          string
	} `json:"-"`
}

// note: this function is generated only for discriminated unions
func (u AnnotationUnion) Variant() (res struct {
	OfFileCitation *FileCitationAnnotation
	OfFilePath     *FilePathAnnotation
}) {
	switch u.Type {
	case "file_citation":
		v := u.AsFileCitation()
		res.OfFileCitation = &v
	case "file_path":
		v := u.AsFilePath()
		res.OfFilePath = &v
	}
	return
}

func (u AnnotationUnion) WhichKind() string {
	return u.Type
}

func (u AnnotationUnion) AsFileCitation() (v FileCitationAnnotation) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AnnotationUnion) AsFilePath() (v FilePathAnnotation) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AnnotationUnion) RawJSON() string { return u.JSON.raw }

func (r *AnnotationUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AnnotationDeltaUnion struct {
	Index        int64                                   `json:"index"`
	Type         string                                  `json:"type"`
	EndIndex     int64                                   `json:"end_index"`
	FileCitation FileCitationDeltaAnnotationFileCitation `json:"file_citation"`
	StartIndex   int64                                   `json:"start_index"`
	Text         string                                  `json:"text"`
	FilePath     FilePathDeltaAnnotationFilePath         `json:"file_path"`
	JSON         struct {
		Index        resp.Field
		Type         resp.Field
		EndIndex     resp.Field
		FileCitation resp.Field
		StartIndex   resp.Field
		Text         resp.Field
		FilePath     resp.Field
		raw          string
	} `json:"-"`
}

// note: this function is generated only for discriminated unions
func (u AnnotationDeltaUnion) Variant() (res struct {
	OfFileCitation *FileCitationDeltaAnnotation
	OfFilePath     *FilePathDeltaAnnotation
}) {
	switch u.Type {
	case "file_citation":
		v := u.AsFileCitation()
		res.OfFileCitation = &v
	case "file_path":
		v := u.AsFilePath()
		res.OfFilePath = &v
	}
	return
}

func (u AnnotationDeltaUnion) WhichKind() string {
	return u.Type
}

func (u AnnotationDeltaUnion) AsFileCitation() (v FileCitationDeltaAnnotation) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AnnotationDeltaUnion) AsFilePath() (v FilePathDeltaAnnotation) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AnnotationDeltaUnion) RawJSON() string { return u.JSON.raw }

func (r *AnnotationDeltaUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A citation within the message that points to a specific quote from a specific
// File associated with the assistant or the message. Generated when the assistant
// uses the "file_search" tool to search files.
type FileCitationAnnotation struct {
	EndIndex     int64                              `json:"end_index,omitzero,required"`
	FileCitation FileCitationAnnotationFileCitation `json:"file_citation,omitzero,required"`
	StartIndex   int64                              `json:"start_index,omitzero,required"`
	// The text in the message content that needs to be replaced.
	Text string `json:"text,omitzero,required"`
	// Always `file_citation`.
	//
	// This field can be elided, and will be automatically set as "file_citation".
	Type constant.FileCitation `json:"type,required"`
	JSON struct {
		EndIndex     resp.Field
		FileCitation resp.Field
		StartIndex   resp.Field
		Text         resp.Field
		Type         resp.Field
		raw          string
	} `json:"-"`
}

func (r FileCitationAnnotation) RawJSON() string { return r.JSON.raw }
func (r *FileCitationAnnotation) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type FileCitationAnnotationFileCitation struct {
	// The ID of the specific File the citation is from.
	FileID string `json:"file_id,omitzero,required"`
	JSON   struct {
		FileID resp.Field
		raw    string
	} `json:"-"`
}

func (r FileCitationAnnotationFileCitation) RawJSON() string { return r.JSON.raw }
func (r *FileCitationAnnotationFileCitation) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A citation within the message that points to a specific quote from a specific
// File associated with the assistant or the message. Generated when the assistant
// uses the "file_search" tool to search files.
type FileCitationDeltaAnnotation struct {
	// The index of the annotation in the text content part.
	Index int64 `json:"index,omitzero,required"`
	// Always `file_citation`.
	//
	// This field can be elided, and will be automatically set as "file_citation".
	Type         constant.FileCitation                   `json:"type,required"`
	EndIndex     int64                                   `json:"end_index,omitzero"`
	FileCitation FileCitationDeltaAnnotationFileCitation `json:"file_citation,omitzero"`
	StartIndex   int64                                   `json:"start_index,omitzero"`
	// The text in the message content that needs to be replaced.
	Text string `json:"text,omitzero"`
	JSON struct {
		Index        resp.Field
		Type         resp.Field
		EndIndex     resp.Field
		FileCitation resp.Field
		StartIndex   resp.Field
		Text         resp.Field
		raw          string
	} `json:"-"`
}

func (r FileCitationDeltaAnnotation) RawJSON() string { return r.JSON.raw }
func (r *FileCitationDeltaAnnotation) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type FileCitationDeltaAnnotationFileCitation struct {
	// The ID of the specific File the citation is from.
	FileID string `json:"file_id,omitzero"`
	// The specific quote in the file.
	Quote string `json:"quote,omitzero"`
	JSON  struct {
		FileID resp.Field
		Quote  resp.Field
		raw    string
	} `json:"-"`
}

func (r FileCitationDeltaAnnotationFileCitation) RawJSON() string { return r.JSON.raw }
func (r *FileCitationDeltaAnnotationFileCitation) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A URL for the file that's generated when the assistant used the
// `code_interpreter` tool to generate a file.
type FilePathAnnotation struct {
	EndIndex   int64                      `json:"end_index,omitzero,required"`
	FilePath   FilePathAnnotationFilePath `json:"file_path,omitzero,required"`
	StartIndex int64                      `json:"start_index,omitzero,required"`
	// The text in the message content that needs to be replaced.
	Text string `json:"text,omitzero,required"`
	// Always `file_path`.
	//
	// This field can be elided, and will be automatically set as "file_path".
	Type constant.FilePath `json:"type,required"`
	JSON struct {
		EndIndex   resp.Field
		FilePath   resp.Field
		StartIndex resp.Field
		Text       resp.Field
		Type       resp.Field
		raw        string
	} `json:"-"`
}

func (r FilePathAnnotation) RawJSON() string { return r.JSON.raw }
func (r *FilePathAnnotation) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type FilePathAnnotationFilePath struct {
	// The ID of the file that was generated.
	FileID string `json:"file_id,omitzero,required"`
	JSON   struct {
		FileID resp.Field
		raw    string
	} `json:"-"`
}

func (r FilePathAnnotationFilePath) RawJSON() string { return r.JSON.raw }
func (r *FilePathAnnotationFilePath) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A URL for the file that's generated when the assistant used the
// `code_interpreter` tool to generate a file.
type FilePathDeltaAnnotation struct {
	// The index of the annotation in the text content part.
	Index int64 `json:"index,omitzero,required"`
	// Always `file_path`.
	//
	// This field can be elided, and will be automatically set as "file_path".
	Type       constant.FilePath               `json:"type,required"`
	EndIndex   int64                           `json:"end_index,omitzero"`
	FilePath   FilePathDeltaAnnotationFilePath `json:"file_path,omitzero"`
	StartIndex int64                           `json:"start_index,omitzero"`
	// The text in the message content that needs to be replaced.
	Text string `json:"text,omitzero"`
	JSON struct {
		Index      resp.Field
		Type       resp.Field
		EndIndex   resp.Field
		FilePath   resp.Field
		StartIndex resp.Field
		Text       resp.Field
		raw        string
	} `json:"-"`
}

func (r FilePathDeltaAnnotation) RawJSON() string { return r.JSON.raw }
func (r *FilePathDeltaAnnotation) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type FilePathDeltaAnnotationFilePath struct {
	// The ID of the file that was generated.
	FileID string `json:"file_id,omitzero"`
	JSON   struct {
		FileID resp.Field
		raw    string
	} `json:"-"`
}

func (r FilePathDeltaAnnotationFilePath) RawJSON() string { return r.JSON.raw }
func (r *FilePathDeltaAnnotationFilePath) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ImageFile struct {
	// The [File](https://platform.openai.com/docs/api-reference/files) ID of the image
	// in the message content. Set `purpose="vision"` when uploading the File if you
	// need to later display the file content.
	FileID string `json:"file_id,omitzero,required"`
	// Specifies the detail level of the image if specified by the user. `low` uses
	// fewer tokens, you can opt in to high resolution using `high`.
	//
	// Any of "auto", "low", "high"
	Detail string `json:"detail,omitzero"`
	JSON   struct {
		FileID resp.Field
		Detail resp.Field
		raw    string
	} `json:"-"`
}

func (r ImageFile) RawJSON() string { return r.JSON.raw }
func (r *ImageFile) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this ImageFile to a ImageFileParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// ImageFileParam.IsOverridden()
func (r ImageFile) ToParam() ImageFileParam {
	return param.Override[ImageFileParam](r.RawJSON())
}

// Specifies the detail level of the image if specified by the user. `low` uses
// fewer tokens, you can opt in to high resolution using `high`.
type ImageFileDetail = string

const (
	ImageFileDetailAuto ImageFileDetail = "auto"
	ImageFileDetailLow  ImageFileDetail = "low"
	ImageFileDetailHigh ImageFileDetail = "high"
)

type ImageFileParam struct {
	// The [File](https://platform.openai.com/docs/api-reference/files) ID of the image
	// in the message content. Set `purpose="vision"` when uploading the File if you
	// need to later display the file content.
	FileID param.String `json:"file_id,omitzero,required"`
	// Specifies the detail level of the image if specified by the user. `low` uses
	// fewer tokens, you can opt in to high resolution using `high`.
	//
	// Any of "auto", "low", "high"
	Detail string `json:"detail,omitzero"`
	apiobject
}

func (f ImageFileParam) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

func (r ImageFileParam) MarshalJSON() (data []byte, err error) {
	type shadow ImageFileParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// References an image [File](https://platform.openai.com/docs/api-reference/files)
// in the content of a message.
type ImageFileContentBlock struct {
	ImageFile ImageFile `json:"image_file,omitzero,required"`
	// Always `image_file`.
	//
	// This field can be elided, and will be automatically set as "image_file".
	Type constant.ImageFile `json:"type,required"`
	JSON struct {
		ImageFile resp.Field
		Type      resp.Field
		raw       string
	} `json:"-"`
}

func (r ImageFileContentBlock) RawJSON() string { return r.JSON.raw }
func (r *ImageFileContentBlock) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this ImageFileContentBlock to a ImageFileContentBlockParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// ImageFileContentBlockParam.IsOverridden()
func (r ImageFileContentBlock) ToParam() ImageFileContentBlockParam {
	return param.Override[ImageFileContentBlockParam](r.RawJSON())
}

// References an image [File](https://platform.openai.com/docs/api-reference/files)
// in the content of a message.
type ImageFileContentBlockParam struct {
	ImageFile ImageFileParam `json:"image_file,omitzero,required"`
	// Always `image_file`.
	//
	// This field can be elided, and will be automatically set as "image_file".
	Type constant.ImageFile `json:"type,required"`
	apiobject
}

func (f ImageFileContentBlockParam) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

func (r ImageFileContentBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow ImageFileContentBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}

type ImageFileDelta struct {
	// Specifies the detail level of the image if specified by the user. `low` uses
	// fewer tokens, you can opt in to high resolution using `high`.
	//
	// Any of "auto", "low", "high"
	Detail string `json:"detail,omitzero"`
	// The [File](https://platform.openai.com/docs/api-reference/files) ID of the image
	// in the message content. Set `purpose="vision"` when uploading the File if you
	// need to later display the file content.
	FileID string `json:"file_id,omitzero"`
	JSON   struct {
		Detail resp.Field
		FileID resp.Field
		raw    string
	} `json:"-"`
}

func (r ImageFileDelta) RawJSON() string { return r.JSON.raw }
func (r *ImageFileDelta) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Specifies the detail level of the image if specified by the user. `low` uses
// fewer tokens, you can opt in to high resolution using `high`.
type ImageFileDeltaDetail = string

const (
	ImageFileDeltaDetailAuto ImageFileDeltaDetail = "auto"
	ImageFileDeltaDetailLow  ImageFileDeltaDetail = "low"
	ImageFileDeltaDetailHigh ImageFileDeltaDetail = "high"
)

// References an image [File](https://platform.openai.com/docs/api-reference/files)
// in the content of a message.
type ImageFileDeltaBlock struct {
	// The index of the content part in the message.
	Index int64 `json:"index,omitzero,required"`
	// Always `image_file`.
	//
	// This field can be elided, and will be automatically set as "image_file".
	Type      constant.ImageFile `json:"type,required"`
	ImageFile ImageFileDelta     `json:"image_file,omitzero"`
	JSON      struct {
		Index     resp.Field
		Type      resp.Field
		ImageFile resp.Field
		raw       string
	} `json:"-"`
}

func (r ImageFileDeltaBlock) RawJSON() string { return r.JSON.raw }
func (r *ImageFileDeltaBlock) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ImageURL struct {
	// The external URL of the image, must be a supported image types: jpeg, jpg, png,
	// gif, webp.
	URL string `json:"url,omitzero,required" format:"uri"`
	// Specifies the detail level of the image. `low` uses fewer tokens, you can opt in
	// to high resolution using `high`. Default value is `auto`
	//
	// Any of "auto", "low", "high"
	Detail string `json:"detail,omitzero"`
	JSON   struct {
		URL    resp.Field
		Detail resp.Field
		raw    string
	} `json:"-"`
}

func (r ImageURL) RawJSON() string { return r.JSON.raw }
func (r *ImageURL) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this ImageURL to a ImageURLParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// ImageURLParam.IsOverridden()
func (r ImageURL) ToParam() ImageURLParam {
	return param.Override[ImageURLParam](r.RawJSON())
}

// Specifies the detail level of the image. `low` uses fewer tokens, you can opt in
// to high resolution using `high`. Default value is `auto`
type ImageURLDetail = string

const (
	ImageURLDetailAuto ImageURLDetail = "auto"
	ImageURLDetailLow  ImageURLDetail = "low"
	ImageURLDetailHigh ImageURLDetail = "high"
)

type ImageURLParam struct {
	// The external URL of the image, must be a supported image types: jpeg, jpg, png,
	// gif, webp.
	URL param.String `json:"url,omitzero,required" format:"uri"`
	// Specifies the detail level of the image. `low` uses fewer tokens, you can opt in
	// to high resolution using `high`. Default value is `auto`
	//
	// Any of "auto", "low", "high"
	Detail string `json:"detail,omitzero"`
	apiobject
}

func (f ImageURLParam) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

func (r ImageURLParam) MarshalJSON() (data []byte, err error) {
	type shadow ImageURLParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// References an image URL in the content of a message.
type ImageURLContentBlock struct {
	ImageURL ImageURL `json:"image_url,omitzero,required"`
	// The type of the content part.
	//
	// This field can be elided, and will be automatically set as "image_url".
	Type constant.ImageURL `json:"type,required"`
	JSON struct {
		ImageURL resp.Field
		Type     resp.Field
		raw      string
	} `json:"-"`
}

func (r ImageURLContentBlock) RawJSON() string { return r.JSON.raw }
func (r *ImageURLContentBlock) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this ImageURLContentBlock to a ImageURLContentBlockParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// ImageURLContentBlockParam.IsOverridden()
func (r ImageURLContentBlock) ToParam() ImageURLContentBlockParam {
	return param.Override[ImageURLContentBlockParam](r.RawJSON())
}

// References an image URL in the content of a message.
type ImageURLContentBlockParam struct {
	ImageURL ImageURLParam `json:"image_url,omitzero,required"`
	// The type of the content part.
	//
	// This field can be elided, and will be automatically set as "image_url".
	Type constant.ImageURL `json:"type,required"`
	apiobject
}

func (f ImageURLContentBlockParam) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

func (r ImageURLContentBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow ImageURLContentBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}

type ImageURLDelta struct {
	// Specifies the detail level of the image. `low` uses fewer tokens, you can opt in
	// to high resolution using `high`.
	//
	// Any of "auto", "low", "high"
	Detail string `json:"detail,omitzero"`
	// The URL of the image, must be a supported image types: jpeg, jpg, png, gif,
	// webp.
	URL  string `json:"url,omitzero"`
	JSON struct {
		Detail resp.Field
		URL    resp.Field
		raw    string
	} `json:"-"`
}

func (r ImageURLDelta) RawJSON() string { return r.JSON.raw }
func (r *ImageURLDelta) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Specifies the detail level of the image. `low` uses fewer tokens, you can opt in
// to high resolution using `high`.
type ImageURLDeltaDetail = string

const (
	ImageURLDeltaDetailAuto ImageURLDeltaDetail = "auto"
	ImageURLDeltaDetailLow  ImageURLDeltaDetail = "low"
	ImageURLDeltaDetailHigh ImageURLDeltaDetail = "high"
)

// References an image URL in the content of a message.
type ImageURLDeltaBlock struct {
	// The index of the content part in the message.
	Index int64 `json:"index,omitzero,required"`
	// Always `image_url`.
	//
	// This field can be elided, and will be automatically set as "image_url".
	Type     constant.ImageURL `json:"type,required"`
	ImageURL ImageURLDelta     `json:"image_url,omitzero"`
	JSON     struct {
		Index    resp.Field
		Type     resp.Field
		ImageURL resp.Field
		raw      string
	} `json:"-"`
}

func (r ImageURLDeltaBlock) RawJSON() string { return r.JSON.raw }
func (r *ImageURLDeltaBlock) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Represents a message within a
// [thread](https://platform.openai.com/docs/api-reference/threads).
type Message struct {
	// The identifier, which can be referenced in API endpoints.
	ID string `json:"id,omitzero,required"`
	// If applicable, the ID of the
	// [assistant](https://platform.openai.com/docs/api-reference/assistants) that
	// authored this message.
	AssistantID string `json:"assistant_id,omitzero,required,nullable"`
	// A list of files attached to the message, and the tools they were added to.
	Attachments []MessageAttachment `json:"attachments,omitzero,required,nullable"`
	// The Unix timestamp (in seconds) for when the message was completed.
	CompletedAt int64 `json:"completed_at,omitzero,required,nullable"`
	// The content of the message in array of text and/or images.
	Content []MessageContentUnion `json:"content,omitzero,required"`
	// The Unix timestamp (in seconds) for when the message was created.
	CreatedAt int64 `json:"created_at,omitzero,required"`
	// The Unix timestamp (in seconds) for when the message was marked as incomplete.
	IncompleteAt int64 `json:"incomplete_at,omitzero,required,nullable"`
	// On an incomplete message, details about why the message is incomplete.
	IncompleteDetails MessageIncompleteDetails `json:"incomplete_details,omitzero,required,nullable"`
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard.
	//
	// Keys are strings with a maximum length of 64 characters. Values are strings with
	// a maximum length of 512 characters.
	Metadata shared.Metadata `json:"metadata,omitzero,required,nullable"`
	// The object type, which is always `thread.message`.
	//
	// This field can be elided, and will be automatically set as "thread.message".
	Object constant.ThreadMessage `json:"object,required"`
	// The entity that produced the message. One of `user` or `assistant`.
	//
	// Any of "user", "assistant"
	Role string `json:"role,omitzero,required"`
	// The ID of the [run](https://platform.openai.com/docs/api-reference/runs)
	// associated with the creation of this message. Value is `null` when messages are
	// created manually using the create message or create thread endpoints.
	RunID string `json:"run_id,omitzero,required,nullable"`
	// The status of the message, which can be either `in_progress`, `incomplete`, or
	// `completed`.
	//
	// Any of "in_progress", "incomplete", "completed"
	Status string `json:"status,omitzero,required"`
	// The [thread](https://platform.openai.com/docs/api-reference/threads) ID that
	// this message belongs to.
	ThreadID string `json:"thread_id,omitzero,required"`
	JSON     struct {
		ID                resp.Field
		AssistantID       resp.Field
		Attachments       resp.Field
		CompletedAt       resp.Field
		Content           resp.Field
		CreatedAt         resp.Field
		IncompleteAt      resp.Field
		IncompleteDetails resp.Field
		Metadata          resp.Field
		Object            resp.Field
		Role              resp.Field
		RunID             resp.Field
		Status            resp.Field
		ThreadID          resp.Field
		raw               string
	} `json:"-"`
}

func (r Message) RawJSON() string { return r.JSON.raw }
func (r *Message) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type MessageAttachment struct {
	// The ID of the file to attach to the message.
	FileID string `json:"file_id,omitzero"`
	// The tools to add this file to.
	Tools []MessageAttachmentsToolsUnion `json:"tools,omitzero"`
	JSON  struct {
		FileID resp.Field
		Tools  resp.Field
		raw    string
	} `json:"-"`
}

func (r MessageAttachment) RawJSON() string { return r.JSON.raw }
func (r *MessageAttachment) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type MessageAttachmentsToolsUnion struct {
	Type string `json:"type"`
	JSON struct {
		Type resp.Field
		raw  string
	} `json:"-"`
}

func (u MessageAttachmentsToolsUnion) AsCodeInterpreterTool() (v CodeInterpreterTool) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u MessageAttachmentsToolsUnion) AsFileSearchTool() (v MessageAttachmentsToolsAssistantToolsFileSearchTypeOnly) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u MessageAttachmentsToolsUnion) RawJSON() string { return u.JSON.raw }

func (r *MessageAttachmentsToolsUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type MessageAttachmentsToolsAssistantToolsFileSearchTypeOnly struct {
	// The type of tool being defined: `file_search`
	//
	// This field can be elided, and will be automatically set as "file_search".
	Type constant.FileSearch `json:"type,required"`
	JSON struct {
		Type resp.Field
		raw  string
	} `json:"-"`
}

func (r MessageAttachmentsToolsAssistantToolsFileSearchTypeOnly) RawJSON() string { return r.JSON.raw }
func (r *MessageAttachmentsToolsAssistantToolsFileSearchTypeOnly) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// On an incomplete message, details about why the message is incomplete.
type MessageIncompleteDetails struct {
	// The reason the message is incomplete.
	//
	// Any of "content_filter", "max_tokens", "run_cancelled", "run_expired",
	// "run_failed"
	Reason string `json:"reason,omitzero,required"`
	JSON   struct {
		Reason resp.Field
		raw    string
	} `json:"-"`
}

func (r MessageIncompleteDetails) RawJSON() string { return r.JSON.raw }
func (r *MessageIncompleteDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The reason the message is incomplete.
type MessageIncompleteDetailsReason = string

const (
	MessageIncompleteDetailsReasonContentFilter MessageIncompleteDetailsReason = "content_filter"
	MessageIncompleteDetailsReasonMaxTokens     MessageIncompleteDetailsReason = "max_tokens"
	MessageIncompleteDetailsReasonRunCancelled  MessageIncompleteDetailsReason = "run_cancelled"
	MessageIncompleteDetailsReasonRunExpired    MessageIncompleteDetailsReason = "run_expired"
	MessageIncompleteDetailsReasonRunFailed     MessageIncompleteDetailsReason = "run_failed"
)

// The entity that produced the message. One of `user` or `assistant`.
type MessageRole = string

const (
	MessageRoleUser      MessageRole = "user"
	MessageRoleAssistant MessageRole = "assistant"
)

// The status of the message, which can be either `in_progress`, `incomplete`, or
// `completed`.
type MessageStatus = string

const (
	MessageStatusInProgress MessageStatus = "in_progress"
	MessageStatusIncomplete MessageStatus = "incomplete"
	MessageStatusCompleted  MessageStatus = "completed"
)

type MessageContentUnion struct {
	ImageFile ImageFile `json:"image_file"`
	Type      string    `json:"type"`
	ImageURL  ImageURL  `json:"image_url"`
	Text      Text      `json:"text"`
	Refusal   string    `json:"refusal"`
	JSON      struct {
		ImageFile resp.Field
		Type      resp.Field
		ImageURL  resp.Field
		Text      resp.Field
		Refusal   resp.Field
		raw       string
	} `json:"-"`
}

// note: this function is generated only for discriminated unions
func (u MessageContentUnion) Variant() (res struct {
	OfImageFile *ImageFileContentBlock
	OfImageURL  *ImageURLContentBlock
	OfText      *TextContentBlock
	OfRefusal   *RefusalContentBlock
}) {
	switch u.Type {
	case "image_file":
		v := u.AsImageFile()
		res.OfImageFile = &v
	case "image_url":
		v := u.AsImageURL()
		res.OfImageURL = &v
	case "text":
		v := u.AsText()
		res.OfText = &v
	case "refusal":
		v := u.AsRefusal()
		res.OfRefusal = &v
	}
	return
}

func (u MessageContentUnion) WhichKind() string {
	return u.Type
}

func (u MessageContentUnion) AsImageFile() (v ImageFileContentBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u MessageContentUnion) AsImageURL() (v ImageURLContentBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u MessageContentUnion) AsText() (v TextContentBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u MessageContentUnion) AsRefusal() (v RefusalContentBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u MessageContentUnion) RawJSON() string { return u.JSON.raw }

func (r *MessageContentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type MessageContentDeltaUnion struct {
	Index     int64          `json:"index"`
	Type      string         `json:"type"`
	ImageFile ImageFileDelta `json:"image_file"`
	Text      TextDelta      `json:"text"`
	Refusal   string         `json:"refusal"`
	ImageURL  ImageURLDelta  `json:"image_url"`
	JSON      struct {
		Index     resp.Field
		Type      resp.Field
		ImageFile resp.Field
		Text      resp.Field
		Refusal   resp.Field
		ImageURL  resp.Field
		raw       string
	} `json:"-"`
}

// note: this function is generated only for discriminated unions
func (u MessageContentDeltaUnion) Variant() (res struct {
	OfImageFile *ImageFileDeltaBlock
	OfText      *TextDeltaBlock
	OfRefusal   *RefusalDeltaBlock
	OfImageURL  *ImageURLDeltaBlock
}) {
	switch u.Type {
	case "image_file":
		v := u.AsImageFile()
		res.OfImageFile = &v
	case "text":
		v := u.AsText()
		res.OfText = &v
	case "refusal":
		v := u.AsRefusal()
		res.OfRefusal = &v
	case "image_url":
		v := u.AsImageURL()
		res.OfImageURL = &v
	}
	return
}

func (u MessageContentDeltaUnion) WhichKind() string {
	return u.Type
}

func (u MessageContentDeltaUnion) AsImageFile() (v ImageFileDeltaBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u MessageContentDeltaUnion) AsText() (v TextDeltaBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u MessageContentDeltaUnion) AsRefusal() (v RefusalDeltaBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u MessageContentDeltaUnion) AsImageURL() (v ImageURLDeltaBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u MessageContentDeltaUnion) RawJSON() string { return u.JSON.raw }

func (r *MessageContentDeltaUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewMessageContentPartParamOfImageFile(imageFile ImageFileParam) MessageContentPartParamUnion {
	var image_file ImageFileContentBlockParam
	image_file.ImageFile = imageFile
	return MessageContentPartParamUnion{OfImageFile: &image_file}
}

func NewMessageContentPartParamOfImageURL(imageURL ImageURLParam) MessageContentPartParamUnion {
	var image_url ImageURLContentBlockParam
	image_url.ImageURL = imageURL
	return MessageContentPartParamUnion{OfImageURL: &image_url}
}

func NewMessageContentPartParamOfText(text string) MessageContentPartParamUnion {
	var variant TextContentBlockParam
	variant.Text = newString(text)
	return MessageContentPartParamUnion{OfText: &variant}
}

// Only one field can be non-zero
type MessageContentPartParamUnion struct {
	OfImageFile *ImageFileContentBlockParam
	OfImageURL  *ImageURLContentBlockParam
	OfText      *TextContentBlockParam
	apiunion
}

func (u MessageContentPartParamUnion) IsMissing() bool { return param.IsOmitted(u) || u.IsNull() }

func (u MessageContentPartParamUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion[MessageContentPartParamUnion](u.OfImageFile, u.OfImageURL, u.OfText)
}

func (u MessageContentPartParamUnion) GetImageFile() *ImageFileParam {
	if vt := u.OfImageFile; vt != nil {
		return &vt.ImageFile
	}
	return nil
}

func (u MessageContentPartParamUnion) GetImageURL() *ImageURLParam {
	if vt := u.OfImageURL; vt != nil {
		return &vt.ImageURL
	}
	return nil
}

func (u MessageContentPartParamUnion) GetText() *string {
	if vt := u.OfText; vt != nil && !vt.Text.IsOmitted() {
		return &vt.Text.V
	}
	return nil
}

func (u MessageContentPartParamUnion) GetType() *string {
	if vt := u.OfImageFile; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfImageURL; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfText; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

type MessageDeleted struct {
	ID      string `json:"id,omitzero,required"`
	Deleted bool   `json:"deleted,omitzero,required"`
	// This field can be elided, and will be automatically set as
	// "thread.message.deleted".
	Object constant.ThreadMessageDeleted `json:"object,required"`
	JSON   struct {
		ID      resp.Field
		Deleted resp.Field
		Object  resp.Field
		raw     string
	} `json:"-"`
}

func (r MessageDeleted) RawJSON() string { return r.JSON.raw }
func (r *MessageDeleted) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The delta containing the fields that have changed on the Message.
type MessageDelta struct {
	// The content of the message in array of text and/or images.
	Content []MessageContentDeltaUnion `json:"content,omitzero"`
	// The entity that produced the message. One of `user` or `assistant`.
	//
	// Any of "user", "assistant"
	Role string `json:"role,omitzero"`
	JSON struct {
		Content resp.Field
		Role    resp.Field
		raw     string
	} `json:"-"`
}

func (r MessageDelta) RawJSON() string { return r.JSON.raw }
func (r *MessageDelta) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The entity that produced the message. One of `user` or `assistant`.
type MessageDeltaRole = string

const (
	MessageDeltaRoleUser      MessageDeltaRole = "user"
	MessageDeltaRoleAssistant MessageDeltaRole = "assistant"
)

// Represents a message delta i.e. any changed fields on a message during
// streaming.
type MessageDeltaEvent struct {
	// The identifier of the message, which can be referenced in API endpoints.
	ID string `json:"id,omitzero,required"`
	// The delta containing the fields that have changed on the Message.
	Delta MessageDelta `json:"delta,omitzero,required"`
	// The object type, which is always `thread.message.delta`.
	//
	// This field can be elided, and will be automatically set as
	// "thread.message.delta".
	Object constant.ThreadMessageDelta `json:"object,required"`
	JSON   struct {
		ID     resp.Field
		Delta  resp.Field
		Object resp.Field
		raw    string
	} `json:"-"`
}

func (r MessageDeltaEvent) RawJSON() string { return r.JSON.raw }
func (r *MessageDeltaEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The refusal content generated by the assistant.
type RefusalContentBlock struct {
	Refusal string `json:"refusal,omitzero,required"`
	// Always `refusal`.
	//
	// This field can be elided, and will be automatically set as "refusal".
	Type constant.Refusal `json:"type,required"`
	JSON struct {
		Refusal resp.Field
		Type    resp.Field
		raw     string
	} `json:"-"`
}

func (r RefusalContentBlock) RawJSON() string { return r.JSON.raw }
func (r *RefusalContentBlock) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The refusal content that is part of a message.
type RefusalDeltaBlock struct {
	// The index of the refusal part in the message.
	Index int64 `json:"index,omitzero,required"`
	// Always `refusal`.
	//
	// This field can be elided, and will be automatically set as "refusal".
	Type    constant.Refusal `json:"type,required"`
	Refusal string           `json:"refusal,omitzero"`
	JSON    struct {
		Index   resp.Field
		Type    resp.Field
		Refusal resp.Field
		raw     string
	} `json:"-"`
}

func (r RefusalDeltaBlock) RawJSON() string { return r.JSON.raw }
func (r *RefusalDeltaBlock) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type Text struct {
	Annotations []AnnotationUnion `json:"annotations,omitzero,required"`
	// The data that makes up the text.
	Value string `json:"value,omitzero,required"`
	JSON  struct {
		Annotations resp.Field
		Value       resp.Field
		raw         string
	} `json:"-"`
}

func (r Text) RawJSON() string { return r.JSON.raw }
func (r *Text) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The text content that is part of a message.
type TextContentBlock struct {
	Text Text `json:"text,omitzero,required"`
	// Always `text`.
	//
	// This field can be elided, and will be automatically set as "text".
	Type constant.Text `json:"type,required"`
	JSON struct {
		Text resp.Field
		Type resp.Field
		raw  string
	} `json:"-"`
}

func (r TextContentBlock) RawJSON() string { return r.JSON.raw }
func (r *TextContentBlock) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The text content that is part of a message.
type TextContentBlockParam struct {
	// Text content to be sent to the model
	Text param.String `json:"text,omitzero,required"`
	// Always `text`.
	//
	// This field can be elided, and will be automatically set as "text".
	Type constant.Text `json:"type,required"`
	apiobject
}

func (f TextContentBlockParam) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

func (r TextContentBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow TextContentBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}

type TextDelta struct {
	Annotations []AnnotationDeltaUnion `json:"annotations,omitzero"`
	// The data that makes up the text.
	Value string `json:"value,omitzero"`
	JSON  struct {
		Annotations resp.Field
		Value       resp.Field
		raw         string
	} `json:"-"`
}

func (r TextDelta) RawJSON() string { return r.JSON.raw }
func (r *TextDelta) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The text content that is part of a message.
type TextDeltaBlock struct {
	// The index of the content part in the message.
	Index int64 `json:"index,omitzero,required"`
	// Always `text`.
	//
	// This field can be elided, and will be automatically set as "text".
	Type constant.Text `json:"type,required"`
	Text TextDelta     `json:"text,omitzero"`
	JSON struct {
		Index resp.Field
		Type  resp.Field
		Text  resp.Field
		raw   string
	} `json:"-"`
}

func (r TextDeltaBlock) RawJSON() string { return r.JSON.raw }
func (r *TextDeltaBlock) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaThreadMessageNewParams struct {
	// An array of content parts with a defined type, each can be of type `text` or
	// images can be passed with `image_url` or `image_file`. Image types are only
	// supported on
	// [Vision-compatible models](https://platform.openai.com/docs/models).
	Content []MessageContentPartParamUnion `json:"content,omitzero,required"`
	// The role of the entity that is creating the message. Allowed values include:
	//
	//   - `user`: Indicates the message is sent by an actual user and should be used in
	//     most cases to represent user-generated messages.
	//   - `assistant`: Indicates the message is generated by the assistant. Use this
	//     value to insert messages from the assistant into the conversation.
	//
	// Any of "user", "assistant"
	Role BetaThreadMessageNewParamsRole `json:"role,omitzero,required"`
	// A list of files attached to the message, and the tools they should be added to.
	Attachments []BetaThreadMessageNewParamsAttachment `json:"attachments,omitzero"`
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard.
	//
	// Keys are strings with a maximum length of 64 characters. Values are strings with
	// a maximum length of 512 characters.
	Metadata shared.MetadataParam `json:"metadata,omitzero"`
	apiobject
}

func (f BetaThreadMessageNewParams) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

func (r BetaThreadMessageNewParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaThreadMessageNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}

// The role of the entity that is creating the message. Allowed values include:
//
//   - `user`: Indicates the message is sent by an actual user and should be used in
//     most cases to represent user-generated messages.
//   - `assistant`: Indicates the message is generated by the assistant. Use this
//     value to insert messages from the assistant into the conversation.
type BetaThreadMessageNewParamsRole string

const (
	BetaThreadMessageNewParamsRoleUser      BetaThreadMessageNewParamsRole = "user"
	BetaThreadMessageNewParamsRoleAssistant BetaThreadMessageNewParamsRole = "assistant"
)

type BetaThreadMessageNewParamsAttachment struct {
	// The ID of the file to attach to the message.
	FileID param.String `json:"file_id,omitzero"`
	// The tools to add this file to.
	Tools []BetaThreadMessageNewParamsAttachmentsToolUnion `json:"tools,omitzero"`
	apiobject
}

func (f BetaThreadMessageNewParamsAttachment) IsMissing() bool {
	return param.IsOmitted(f) || f.IsNull()
}

func (r BetaThreadMessageNewParamsAttachment) MarshalJSON() (data []byte, err error) {
	type shadow BetaThreadMessageNewParamsAttachment
	return param.MarshalObject(r, (*shadow)(&r))
}

// Only one field can be non-zero
type BetaThreadMessageNewParamsAttachmentsToolUnion struct {
	OfCodeInterpreter *CodeInterpreterToolParam
	OfFileSearch      *BetaThreadMessageNewParamsAttachmentsToolsFileSearch
	apiunion
}

func (u BetaThreadMessageNewParamsAttachmentsToolUnion) IsMissing() bool {
	return param.IsOmitted(u) || u.IsNull()
}

func (u BetaThreadMessageNewParamsAttachmentsToolUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion[BetaThreadMessageNewParamsAttachmentsToolUnion](u.OfCodeInterpreter, u.OfFileSearch)
}

func (u BetaThreadMessageNewParamsAttachmentsToolUnion) GetType() *string {
	if vt := u.OfCodeInterpreter; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfFileSearch; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

type BetaThreadMessageNewParamsAttachmentsToolsFileSearch struct {
	// The type of tool being defined: `file_search`
	//
	// This field can be elided, and will be automatically set as "file_search".
	Type constant.FileSearch `json:"type,required"`
	apiobject
}

func (f BetaThreadMessageNewParamsAttachmentsToolsFileSearch) IsMissing() bool {
	return param.IsOmitted(f) || f.IsNull()
}

func (r BetaThreadMessageNewParamsAttachmentsToolsFileSearch) MarshalJSON() (data []byte, err error) {
	type shadow BetaThreadMessageNewParamsAttachmentsToolsFileSearch
	return param.MarshalObject(r, (*shadow)(&r))
}

type BetaThreadMessageUpdateParams struct {
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard.
	//
	// Keys are strings with a maximum length of 64 characters. Values are strings with
	// a maximum length of 512 characters.
	Metadata shared.MetadataParam `json:"metadata,omitzero"`
	apiobject
}

func (f BetaThreadMessageUpdateParams) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

func (r BetaThreadMessageUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaThreadMessageUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}

type BetaThreadMessageListParams struct {
	// A cursor for use in pagination. `after` is an object ID that defines your place
	// in the list. For instance, if you make a list request and receive 100 objects,
	// ending with obj_foo, your subsequent call can include after=obj_foo in order to
	// fetch the next page of the list.
	After param.String `query:"after,omitzero"`
	// A cursor for use in pagination. `before` is an object ID that defines your place
	// in the list. For instance, if you make a list request and receive 100 objects,
	// starting with obj_foo, your subsequent call can include before=obj_foo in order
	// to fetch the previous page of the list.
	Before param.String `query:"before,omitzero"`
	// A limit on the number of objects to be returned. Limit can range between 1 and
	// 100, and the default is 20.
	Limit param.Int `query:"limit,omitzero"`
	// Sort order by the `created_at` timestamp of the objects. `asc` for ascending
	// order and `desc` for descending order.
	//
	// Any of "asc", "desc"
	Order BetaThreadMessageListParamsOrder `query:"order,omitzero"`
	// Filter messages by the run ID that generated them.
	RunID param.String `query:"run_id,omitzero"`
	apiobject
}

func (f BetaThreadMessageListParams) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

// URLQuery serializes [BetaThreadMessageListParams]'s query parameters as
// `url.Values`.
func (r BetaThreadMessageListParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Sort order by the `created_at` timestamp of the objects. `asc` for ascending
// order and `desc` for descending order.
type BetaThreadMessageListParamsOrder string

const (
	BetaThreadMessageListParamsOrderAsc  BetaThreadMessageListParamsOrder = "asc"
	BetaThreadMessageListParamsOrderDesc BetaThreadMessageListParamsOrder = "desc"
)

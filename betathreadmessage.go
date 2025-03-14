// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"

	"github.com/openai/openai-go/internal/apijson"
	"github.com/openai/openai-go/internal/apiquery"
	"github.com/openai/openai-go/internal/param"
	"github.com/openai/openai-go/internal/requestconfig"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/pagination"
	"github.com/openai/openai-go/shared"
	"github.com/tidwall/gjson"
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
func NewBetaThreadMessageService(opts ...option.RequestOption) (r *BetaThreadMessageService) {
	r = &BetaThreadMessageService{}
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

// A citation within the message that points to a specific quote from a specific
// File associated with the assistant or the message. Generated when the assistant
// uses the "file_search" tool to search files.
type Annotation struct {
	EndIndex   int64 `json:"end_index,required"`
	StartIndex int64 `json:"start_index,required"`
	// The text in the message content that needs to be replaced.
	Text string `json:"text,required"`
	// Always `file_citation`.
	Type AnnotationType `json:"type,required"`
	// This field can have the runtime type of [FileCitationAnnotationFileCitation].
	FileCitation interface{} `json:"file_citation"`
	// This field can have the runtime type of [FilePathAnnotationFilePath].
	FilePath interface{}    `json:"file_path"`
	JSON     annotationJSON `json:"-"`
	union    AnnotationUnion
}

// annotationJSON contains the JSON metadata for the struct [Annotation]
type annotationJSON struct {
	EndIndex     apijson.Field
	StartIndex   apijson.Field
	Text         apijson.Field
	Type         apijson.Field
	FileCitation apijson.Field
	FilePath     apijson.Field
	raw          string
	ExtraFields  map[string]apijson.Field
}

func (r annotationJSON) RawJSON() string {
	return r.raw
}

func (r *Annotation) UnmarshalJSON(data []byte) (err error) {
	*r = Annotation{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [AnnotationUnion] interface which you can cast to the specific
// types for more type safety.
//
// Possible runtime types of the union are [FileCitationAnnotation],
// [FilePathAnnotation].
func (r Annotation) AsUnion() AnnotationUnion {
	return r.union
}

// A citation within the message that points to a specific quote from a specific
// File associated with the assistant or the message. Generated when the assistant
// uses the "file_search" tool to search files.
//
// Union satisfied by [FileCitationAnnotation] or [FilePathAnnotation].
type AnnotationUnion interface {
	implementsAnnotation()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*AnnotationUnion)(nil)).Elem(),
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(FileCitationAnnotation{}),
			DiscriminatorValue: "file_citation",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(FilePathAnnotation{}),
			DiscriminatorValue: "file_path",
		},
	)
}

// Always `file_citation`.
type AnnotationType string

const (
	AnnotationTypeFileCitation AnnotationType = "file_citation"
	AnnotationTypeFilePath     AnnotationType = "file_path"
)

func (r AnnotationType) IsKnown() bool {
	switch r {
	case AnnotationTypeFileCitation, AnnotationTypeFilePath:
		return true
	}
	return false
}

// A citation within the message that points to a specific quote from a specific
// File associated with the assistant or the message. Generated when the assistant
// uses the "file_search" tool to search files.
type AnnotationDelta struct {
	// The index of the annotation in the text content part.
	Index int64 `json:"index,required"`
	// Always `file_citation`.
	Type     AnnotationDeltaType `json:"type,required"`
	EndIndex int64               `json:"end_index"`
	// This field can have the runtime type of
	// [FileCitationDeltaAnnotationFileCitation].
	FileCitation interface{} `json:"file_citation"`
	// This field can have the runtime type of [FilePathDeltaAnnotationFilePath].
	FilePath   interface{} `json:"file_path"`
	StartIndex int64       `json:"start_index"`
	// The text in the message content that needs to be replaced.
	Text  string              `json:"text"`
	JSON  annotationDeltaJSON `json:"-"`
	union AnnotationDeltaUnion
}

// annotationDeltaJSON contains the JSON metadata for the struct [AnnotationDelta]
type annotationDeltaJSON struct {
	Index        apijson.Field
	Type         apijson.Field
	EndIndex     apijson.Field
	FileCitation apijson.Field
	FilePath     apijson.Field
	StartIndex   apijson.Field
	Text         apijson.Field
	raw          string
	ExtraFields  map[string]apijson.Field
}

func (r annotationDeltaJSON) RawJSON() string {
	return r.raw
}

func (r *AnnotationDelta) UnmarshalJSON(data []byte) (err error) {
	*r = AnnotationDelta{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [AnnotationDeltaUnion] interface which you can cast to the
// specific types for more type safety.
//
// Possible runtime types of the union are [FileCitationDeltaAnnotation],
// [FilePathDeltaAnnotation].
func (r AnnotationDelta) AsUnion() AnnotationDeltaUnion {
	return r.union
}

// A citation within the message that points to a specific quote from a specific
// File associated with the assistant or the message. Generated when the assistant
// uses the "file_search" tool to search files.
//
// Union satisfied by [FileCitationDeltaAnnotation] or [FilePathDeltaAnnotation].
type AnnotationDeltaUnion interface {
	implementsAnnotationDelta()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*AnnotationDeltaUnion)(nil)).Elem(),
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(FileCitationDeltaAnnotation{}),
			DiscriminatorValue: "file_citation",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(FilePathDeltaAnnotation{}),
			DiscriminatorValue: "file_path",
		},
	)
}

// Always `file_citation`.
type AnnotationDeltaType string

const (
	AnnotationDeltaTypeFileCitation AnnotationDeltaType = "file_citation"
	AnnotationDeltaTypeFilePath     AnnotationDeltaType = "file_path"
)

func (r AnnotationDeltaType) IsKnown() bool {
	switch r {
	case AnnotationDeltaTypeFileCitation, AnnotationDeltaTypeFilePath:
		return true
	}
	return false
}

// A citation within the message that points to a specific quote from a specific
// File associated with the assistant or the message. Generated when the assistant
// uses the "file_search" tool to search files.
type FileCitationAnnotation struct {
	EndIndex     int64                              `json:"end_index,required"`
	FileCitation FileCitationAnnotationFileCitation `json:"file_citation,required"`
	StartIndex   int64                              `json:"start_index,required"`
	// The text in the message content that needs to be replaced.
	Text string `json:"text,required"`
	// Always `file_citation`.
	Type FileCitationAnnotationType `json:"type,required"`
	JSON fileCitationAnnotationJSON `json:"-"`
}

// fileCitationAnnotationJSON contains the JSON metadata for the struct
// [FileCitationAnnotation]
type fileCitationAnnotationJSON struct {
	EndIndex     apijson.Field
	FileCitation apijson.Field
	StartIndex   apijson.Field
	Text         apijson.Field
	Type         apijson.Field
	raw          string
	ExtraFields  map[string]apijson.Field
}

func (r *FileCitationAnnotation) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r fileCitationAnnotationJSON) RawJSON() string {
	return r.raw
}

func (r FileCitationAnnotation) implementsAnnotation() {}

type FileCitationAnnotationFileCitation struct {
	// The ID of the specific File the citation is from.
	FileID string                                 `json:"file_id,required"`
	JSON   fileCitationAnnotationFileCitationJSON `json:"-"`
}

// fileCitationAnnotationFileCitationJSON contains the JSON metadata for the struct
// [FileCitationAnnotationFileCitation]
type fileCitationAnnotationFileCitationJSON struct {
	FileID      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *FileCitationAnnotationFileCitation) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r fileCitationAnnotationFileCitationJSON) RawJSON() string {
	return r.raw
}

// Always `file_citation`.
type FileCitationAnnotationType string

const (
	FileCitationAnnotationTypeFileCitation FileCitationAnnotationType = "file_citation"
)

func (r FileCitationAnnotationType) IsKnown() bool {
	switch r {
	case FileCitationAnnotationTypeFileCitation:
		return true
	}
	return false
}

// A citation within the message that points to a specific quote from a specific
// File associated with the assistant or the message. Generated when the assistant
// uses the "file_search" tool to search files.
type FileCitationDeltaAnnotation struct {
	// The index of the annotation in the text content part.
	Index int64 `json:"index,required"`
	// Always `file_citation`.
	Type         FileCitationDeltaAnnotationType         `json:"type,required"`
	EndIndex     int64                                   `json:"end_index"`
	FileCitation FileCitationDeltaAnnotationFileCitation `json:"file_citation"`
	StartIndex   int64                                   `json:"start_index"`
	// The text in the message content that needs to be replaced.
	Text string                          `json:"text"`
	JSON fileCitationDeltaAnnotationJSON `json:"-"`
}

// fileCitationDeltaAnnotationJSON contains the JSON metadata for the struct
// [FileCitationDeltaAnnotation]
type fileCitationDeltaAnnotationJSON struct {
	Index        apijson.Field
	Type         apijson.Field
	EndIndex     apijson.Field
	FileCitation apijson.Field
	StartIndex   apijson.Field
	Text         apijson.Field
	raw          string
	ExtraFields  map[string]apijson.Field
}

func (r *FileCitationDeltaAnnotation) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r fileCitationDeltaAnnotationJSON) RawJSON() string {
	return r.raw
}

func (r FileCitationDeltaAnnotation) implementsAnnotationDelta() {}

// Always `file_citation`.
type FileCitationDeltaAnnotationType string

const (
	FileCitationDeltaAnnotationTypeFileCitation FileCitationDeltaAnnotationType = "file_citation"
)

func (r FileCitationDeltaAnnotationType) IsKnown() bool {
	switch r {
	case FileCitationDeltaAnnotationTypeFileCitation:
		return true
	}
	return false
}

type FileCitationDeltaAnnotationFileCitation struct {
	// The ID of the specific File the citation is from.
	FileID string `json:"file_id"`
	// The specific quote in the file.
	Quote string                                      `json:"quote"`
	JSON  fileCitationDeltaAnnotationFileCitationJSON `json:"-"`
}

// fileCitationDeltaAnnotationFileCitationJSON contains the JSON metadata for the
// struct [FileCitationDeltaAnnotationFileCitation]
type fileCitationDeltaAnnotationFileCitationJSON struct {
	FileID      apijson.Field
	Quote       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *FileCitationDeltaAnnotationFileCitation) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r fileCitationDeltaAnnotationFileCitationJSON) RawJSON() string {
	return r.raw
}

// A URL for the file that's generated when the assistant used the
// `code_interpreter` tool to generate a file.
type FilePathAnnotation struct {
	EndIndex   int64                      `json:"end_index,required"`
	FilePath   FilePathAnnotationFilePath `json:"file_path,required"`
	StartIndex int64                      `json:"start_index,required"`
	// The text in the message content that needs to be replaced.
	Text string `json:"text,required"`
	// Always `file_path`.
	Type FilePathAnnotationType `json:"type,required"`
	JSON filePathAnnotationJSON `json:"-"`
}

// filePathAnnotationJSON contains the JSON metadata for the struct
// [FilePathAnnotation]
type filePathAnnotationJSON struct {
	EndIndex    apijson.Field
	FilePath    apijson.Field
	StartIndex  apijson.Field
	Text        apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *FilePathAnnotation) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r filePathAnnotationJSON) RawJSON() string {
	return r.raw
}

func (r FilePathAnnotation) implementsAnnotation() {}

type FilePathAnnotationFilePath struct {
	// The ID of the file that was generated.
	FileID string                         `json:"file_id,required"`
	JSON   filePathAnnotationFilePathJSON `json:"-"`
}

// filePathAnnotationFilePathJSON contains the JSON metadata for the struct
// [FilePathAnnotationFilePath]
type filePathAnnotationFilePathJSON struct {
	FileID      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *FilePathAnnotationFilePath) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r filePathAnnotationFilePathJSON) RawJSON() string {
	return r.raw
}

// Always `file_path`.
type FilePathAnnotationType string

const (
	FilePathAnnotationTypeFilePath FilePathAnnotationType = "file_path"
)

func (r FilePathAnnotationType) IsKnown() bool {
	switch r {
	case FilePathAnnotationTypeFilePath:
		return true
	}
	return false
}

// A URL for the file that's generated when the assistant used the
// `code_interpreter` tool to generate a file.
type FilePathDeltaAnnotation struct {
	// The index of the annotation in the text content part.
	Index int64 `json:"index,required"`
	// Always `file_path`.
	Type       FilePathDeltaAnnotationType     `json:"type,required"`
	EndIndex   int64                           `json:"end_index"`
	FilePath   FilePathDeltaAnnotationFilePath `json:"file_path"`
	StartIndex int64                           `json:"start_index"`
	// The text in the message content that needs to be replaced.
	Text string                      `json:"text"`
	JSON filePathDeltaAnnotationJSON `json:"-"`
}

// filePathDeltaAnnotationJSON contains the JSON metadata for the struct
// [FilePathDeltaAnnotation]
type filePathDeltaAnnotationJSON struct {
	Index       apijson.Field
	Type        apijson.Field
	EndIndex    apijson.Field
	FilePath    apijson.Field
	StartIndex  apijson.Field
	Text        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *FilePathDeltaAnnotation) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r filePathDeltaAnnotationJSON) RawJSON() string {
	return r.raw
}

func (r FilePathDeltaAnnotation) implementsAnnotationDelta() {}

// Always `file_path`.
type FilePathDeltaAnnotationType string

const (
	FilePathDeltaAnnotationTypeFilePath FilePathDeltaAnnotationType = "file_path"
)

func (r FilePathDeltaAnnotationType) IsKnown() bool {
	switch r {
	case FilePathDeltaAnnotationTypeFilePath:
		return true
	}
	return false
}

type FilePathDeltaAnnotationFilePath struct {
	// The ID of the file that was generated.
	FileID string                              `json:"file_id"`
	JSON   filePathDeltaAnnotationFilePathJSON `json:"-"`
}

// filePathDeltaAnnotationFilePathJSON contains the JSON metadata for the struct
// [FilePathDeltaAnnotationFilePath]
type filePathDeltaAnnotationFilePathJSON struct {
	FileID      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *FilePathDeltaAnnotationFilePath) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r filePathDeltaAnnotationFilePathJSON) RawJSON() string {
	return r.raw
}

type ImageFile struct {
	// The [File](https://platform.openai.com/docs/api-reference/files) ID of the image
	// in the message content. Set `purpose="vision"` when uploading the File if you
	// need to later display the file content.
	FileID string `json:"file_id,required"`
	// Specifies the detail level of the image if specified by the user. `low` uses
	// fewer tokens, you can opt in to high resolution using `high`.
	Detail ImageFileDetail `json:"detail"`
	JSON   imageFileJSON   `json:"-"`
}

// imageFileJSON contains the JSON metadata for the struct [ImageFile]
type imageFileJSON struct {
	FileID      apijson.Field
	Detail      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ImageFile) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r imageFileJSON) RawJSON() string {
	return r.raw
}

// Specifies the detail level of the image if specified by the user. `low` uses
// fewer tokens, you can opt in to high resolution using `high`.
type ImageFileDetail string

const (
	ImageFileDetailAuto ImageFileDetail = "auto"
	ImageFileDetailLow  ImageFileDetail = "low"
	ImageFileDetailHigh ImageFileDetail = "high"
)

func (r ImageFileDetail) IsKnown() bool {
	switch r {
	case ImageFileDetailAuto, ImageFileDetailLow, ImageFileDetailHigh:
		return true
	}
	return false
}

type ImageFileParam struct {
	// The [File](https://platform.openai.com/docs/api-reference/files) ID of the image
	// in the message content. Set `purpose="vision"` when uploading the File if you
	// need to later display the file content.
	FileID param.Field[string] `json:"file_id,required"`
	// Specifies the detail level of the image if specified by the user. `low` uses
	// fewer tokens, you can opt in to high resolution using `high`.
	Detail param.Field[ImageFileDetail] `json:"detail"`
}

func (r ImageFileParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// References an image [File](https://platform.openai.com/docs/api-reference/files)
// in the content of a message.
type ImageFileContentBlock struct {
	ImageFile ImageFile `json:"image_file,required"`
	// Always `image_file`.
	Type ImageFileContentBlockType `json:"type,required"`
	JSON imageFileContentBlockJSON `json:"-"`
}

// imageFileContentBlockJSON contains the JSON metadata for the struct
// [ImageFileContentBlock]
type imageFileContentBlockJSON struct {
	ImageFile   apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ImageFileContentBlock) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r imageFileContentBlockJSON) RawJSON() string {
	return r.raw
}

func (r ImageFileContentBlock) implementsMessageContent() {}

// Always `image_file`.
type ImageFileContentBlockType string

const (
	ImageFileContentBlockTypeImageFile ImageFileContentBlockType = "image_file"
)

func (r ImageFileContentBlockType) IsKnown() bool {
	switch r {
	case ImageFileContentBlockTypeImageFile:
		return true
	}
	return false
}

// References an image [File](https://platform.openai.com/docs/api-reference/files)
// in the content of a message.
type ImageFileContentBlockParam struct {
	ImageFile param.Field[ImageFileParam] `json:"image_file,required"`
	// Always `image_file`.
	Type param.Field[ImageFileContentBlockType] `json:"type,required"`
}

func (r ImageFileContentBlockParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ImageFileContentBlockParam) implementsMessageContentPartParamUnion() {}

type ImageFileDelta struct {
	// Specifies the detail level of the image if specified by the user. `low` uses
	// fewer tokens, you can opt in to high resolution using `high`.
	Detail ImageFileDeltaDetail `json:"detail"`
	// The [File](https://platform.openai.com/docs/api-reference/files) ID of the image
	// in the message content. Set `purpose="vision"` when uploading the File if you
	// need to later display the file content.
	FileID string             `json:"file_id"`
	JSON   imageFileDeltaJSON `json:"-"`
}

// imageFileDeltaJSON contains the JSON metadata for the struct [ImageFileDelta]
type imageFileDeltaJSON struct {
	Detail      apijson.Field
	FileID      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ImageFileDelta) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r imageFileDeltaJSON) RawJSON() string {
	return r.raw
}

// Specifies the detail level of the image if specified by the user. `low` uses
// fewer tokens, you can opt in to high resolution using `high`.
type ImageFileDeltaDetail string

const (
	ImageFileDeltaDetailAuto ImageFileDeltaDetail = "auto"
	ImageFileDeltaDetailLow  ImageFileDeltaDetail = "low"
	ImageFileDeltaDetailHigh ImageFileDeltaDetail = "high"
)

func (r ImageFileDeltaDetail) IsKnown() bool {
	switch r {
	case ImageFileDeltaDetailAuto, ImageFileDeltaDetailLow, ImageFileDeltaDetailHigh:
		return true
	}
	return false
}

// References an image [File](https://platform.openai.com/docs/api-reference/files)
// in the content of a message.
type ImageFileDeltaBlock struct {
	// The index of the content part in the message.
	Index int64 `json:"index,required"`
	// Always `image_file`.
	Type      ImageFileDeltaBlockType `json:"type,required"`
	ImageFile ImageFileDelta          `json:"image_file"`
	JSON      imageFileDeltaBlockJSON `json:"-"`
}

// imageFileDeltaBlockJSON contains the JSON metadata for the struct
// [ImageFileDeltaBlock]
type imageFileDeltaBlockJSON struct {
	Index       apijson.Field
	Type        apijson.Field
	ImageFile   apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ImageFileDeltaBlock) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r imageFileDeltaBlockJSON) RawJSON() string {
	return r.raw
}

func (r ImageFileDeltaBlock) implementsMessageContentDelta() {}

// Always `image_file`.
type ImageFileDeltaBlockType string

const (
	ImageFileDeltaBlockTypeImageFile ImageFileDeltaBlockType = "image_file"
)

func (r ImageFileDeltaBlockType) IsKnown() bool {
	switch r {
	case ImageFileDeltaBlockTypeImageFile:
		return true
	}
	return false
}

type ImageURL struct {
	// The external URL of the image, must be a supported image types: jpeg, jpg, png,
	// gif, webp.
	URL string `json:"url,required" format:"uri"`
	// Specifies the detail level of the image. `low` uses fewer tokens, you can opt in
	// to high resolution using `high`. Default value is `auto`
	Detail ImageURLDetail `json:"detail"`
	JSON   imageURLJSON   `json:"-"`
}

// imageURLJSON contains the JSON metadata for the struct [ImageURL]
type imageURLJSON struct {
	URL         apijson.Field
	Detail      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ImageURL) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r imageURLJSON) RawJSON() string {
	return r.raw
}

// Specifies the detail level of the image. `low` uses fewer tokens, you can opt in
// to high resolution using `high`. Default value is `auto`
type ImageURLDetail string

const (
	ImageURLDetailAuto ImageURLDetail = "auto"
	ImageURLDetailLow  ImageURLDetail = "low"
	ImageURLDetailHigh ImageURLDetail = "high"
)

func (r ImageURLDetail) IsKnown() bool {
	switch r {
	case ImageURLDetailAuto, ImageURLDetailLow, ImageURLDetailHigh:
		return true
	}
	return false
}

type ImageURLParam struct {
	// The external URL of the image, must be a supported image types: jpeg, jpg, png,
	// gif, webp.
	URL param.Field[string] `json:"url,required" format:"uri"`
	// Specifies the detail level of the image. `low` uses fewer tokens, you can opt in
	// to high resolution using `high`. Default value is `auto`
	Detail param.Field[ImageURLDetail] `json:"detail"`
}

func (r ImageURLParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// References an image URL in the content of a message.
type ImageURLContentBlock struct {
	ImageURL ImageURL `json:"image_url,required"`
	// The type of the content part.
	Type ImageURLContentBlockType `json:"type,required"`
	JSON imageURLContentBlockJSON `json:"-"`
}

// imageURLContentBlockJSON contains the JSON metadata for the struct
// [ImageURLContentBlock]
type imageURLContentBlockJSON struct {
	ImageURL    apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ImageURLContentBlock) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r imageURLContentBlockJSON) RawJSON() string {
	return r.raw
}

func (r ImageURLContentBlock) implementsMessageContent() {}

// The type of the content part.
type ImageURLContentBlockType string

const (
	ImageURLContentBlockTypeImageURL ImageURLContentBlockType = "image_url"
)

func (r ImageURLContentBlockType) IsKnown() bool {
	switch r {
	case ImageURLContentBlockTypeImageURL:
		return true
	}
	return false
}

// References an image URL in the content of a message.
type ImageURLContentBlockParam struct {
	ImageURL param.Field[ImageURLParam] `json:"image_url,required"`
	// The type of the content part.
	Type param.Field[ImageURLContentBlockType] `json:"type,required"`
}

func (r ImageURLContentBlockParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ImageURLContentBlockParam) implementsMessageContentPartParamUnion() {}

type ImageURLDelta struct {
	// Specifies the detail level of the image. `low` uses fewer tokens, you can opt in
	// to high resolution using `high`.
	Detail ImageURLDeltaDetail `json:"detail"`
	// The URL of the image, must be a supported image types: jpeg, jpg, png, gif,
	// webp.
	URL  string            `json:"url"`
	JSON imageURLDeltaJSON `json:"-"`
}

// imageURLDeltaJSON contains the JSON metadata for the struct [ImageURLDelta]
type imageURLDeltaJSON struct {
	Detail      apijson.Field
	URL         apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ImageURLDelta) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r imageURLDeltaJSON) RawJSON() string {
	return r.raw
}

// Specifies the detail level of the image. `low` uses fewer tokens, you can opt in
// to high resolution using `high`.
type ImageURLDeltaDetail string

const (
	ImageURLDeltaDetailAuto ImageURLDeltaDetail = "auto"
	ImageURLDeltaDetailLow  ImageURLDeltaDetail = "low"
	ImageURLDeltaDetailHigh ImageURLDeltaDetail = "high"
)

func (r ImageURLDeltaDetail) IsKnown() bool {
	switch r {
	case ImageURLDeltaDetailAuto, ImageURLDeltaDetailLow, ImageURLDeltaDetailHigh:
		return true
	}
	return false
}

// References an image URL in the content of a message.
type ImageURLDeltaBlock struct {
	// The index of the content part in the message.
	Index int64 `json:"index,required"`
	// Always `image_url`.
	Type     ImageURLDeltaBlockType `json:"type,required"`
	ImageURL ImageURLDelta          `json:"image_url"`
	JSON     imageURLDeltaBlockJSON `json:"-"`
}

// imageURLDeltaBlockJSON contains the JSON metadata for the struct
// [ImageURLDeltaBlock]
type imageURLDeltaBlockJSON struct {
	Index       apijson.Field
	Type        apijson.Field
	ImageURL    apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ImageURLDeltaBlock) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r imageURLDeltaBlockJSON) RawJSON() string {
	return r.raw
}

func (r ImageURLDeltaBlock) implementsMessageContentDelta() {}

// Always `image_url`.
type ImageURLDeltaBlockType string

const (
	ImageURLDeltaBlockTypeImageURL ImageURLDeltaBlockType = "image_url"
)

func (r ImageURLDeltaBlockType) IsKnown() bool {
	switch r {
	case ImageURLDeltaBlockTypeImageURL:
		return true
	}
	return false
}

// Represents a message within a
// [thread](https://platform.openai.com/docs/api-reference/threads).
type Message struct {
	// The identifier, which can be referenced in API endpoints.
	ID string `json:"id,required"`
	// If applicable, the ID of the
	// [assistant](https://platform.openai.com/docs/api-reference/assistants) that
	// authored this message.
	AssistantID string `json:"assistant_id,required,nullable"`
	// A list of files attached to the message, and the tools they were added to.
	Attachments []MessageAttachment `json:"attachments,required,nullable"`
	// The Unix timestamp (in seconds) for when the message was completed.
	CompletedAt int64 `json:"completed_at,required,nullable"`
	// The content of the message in array of text and/or images.
	Content []MessageContent `json:"content,required"`
	// The Unix timestamp (in seconds) for when the message was created.
	CreatedAt int64 `json:"created_at,required"`
	// The Unix timestamp (in seconds) for when the message was marked as incomplete.
	IncompleteAt int64 `json:"incomplete_at,required,nullable"`
	// On an incomplete message, details about why the message is incomplete.
	IncompleteDetails MessageIncompleteDetails `json:"incomplete_details,required,nullable"`
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard.
	//
	// Keys are strings with a maximum length of 64 characters. Values are strings with
	// a maximum length of 512 characters.
	Metadata shared.Metadata `json:"metadata,required,nullable"`
	// The object type, which is always `thread.message`.
	Object MessageObject `json:"object,required"`
	// The entity that produced the message. One of `user` or `assistant`.
	Role MessageRole `json:"role,required"`
	// The ID of the [run](https://platform.openai.com/docs/api-reference/runs)
	// associated with the creation of this message. Value is `null` when messages are
	// created manually using the create message or create thread endpoints.
	RunID string `json:"run_id,required,nullable"`
	// The status of the message, which can be either `in_progress`, `incomplete`, or
	// `completed`.
	Status MessageStatus `json:"status,required"`
	// The [thread](https://platform.openai.com/docs/api-reference/threads) ID that
	// this message belongs to.
	ThreadID string      `json:"thread_id,required"`
	JSON     messageJSON `json:"-"`
}

// messageJSON contains the JSON metadata for the struct [Message]
type messageJSON struct {
	ID                apijson.Field
	AssistantID       apijson.Field
	Attachments       apijson.Field
	CompletedAt       apijson.Field
	Content           apijson.Field
	CreatedAt         apijson.Field
	IncompleteAt      apijson.Field
	IncompleteDetails apijson.Field
	Metadata          apijson.Field
	Object            apijson.Field
	Role              apijson.Field
	RunID             apijson.Field
	Status            apijson.Field
	ThreadID          apijson.Field
	raw               string
	ExtraFields       map[string]apijson.Field
}

func (r *Message) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r messageJSON) RawJSON() string {
	return r.raw
}

type MessageAttachment struct {
	// The ID of the file to attach to the message.
	FileID string `json:"file_id"`
	// The tools to add this file to.
	Tools []MessageAttachmentsTool `json:"tools"`
	JSON  messageAttachmentJSON    `json:"-"`
}

// messageAttachmentJSON contains the JSON metadata for the struct
// [MessageAttachment]
type messageAttachmentJSON struct {
	FileID      apijson.Field
	Tools       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *MessageAttachment) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r messageAttachmentJSON) RawJSON() string {
	return r.raw
}

type MessageAttachmentsTool struct {
	// The type of tool being defined: `code_interpreter`
	Type  MessageAttachmentsToolsType `json:"type,required"`
	JSON  messageAttachmentsToolJSON  `json:"-"`
	union MessageAttachmentsToolsUnion
}

// messageAttachmentsToolJSON contains the JSON metadata for the struct
// [MessageAttachmentsTool]
type messageAttachmentsToolJSON struct {
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r messageAttachmentsToolJSON) RawJSON() string {
	return r.raw
}

func (r *MessageAttachmentsTool) UnmarshalJSON(data []byte) (err error) {
	*r = MessageAttachmentsTool{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [MessageAttachmentsToolsUnion] interface which you can cast to
// the specific types for more type safety.
//
// Possible runtime types of the union are [CodeInterpreterTool],
// [MessageAttachmentsToolsAssistantToolsFileSearchTypeOnly].
func (r MessageAttachmentsTool) AsUnion() MessageAttachmentsToolsUnion {
	return r.union
}

// Union satisfied by [CodeInterpreterTool] or
// [MessageAttachmentsToolsAssistantToolsFileSearchTypeOnly].
type MessageAttachmentsToolsUnion interface {
	implementsMessageAttachmentsTool()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*MessageAttachmentsToolsUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CodeInterpreterTool{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(MessageAttachmentsToolsAssistantToolsFileSearchTypeOnly{}),
		},
	)
}

type MessageAttachmentsToolsAssistantToolsFileSearchTypeOnly struct {
	// The type of tool being defined: `file_search`
	Type MessageAttachmentsToolsAssistantToolsFileSearchTypeOnlyType `json:"type,required"`
	JSON messageAttachmentsToolsAssistantToolsFileSearchTypeOnlyJSON `json:"-"`
}

// messageAttachmentsToolsAssistantToolsFileSearchTypeOnlyJSON contains the JSON
// metadata for the struct
// [MessageAttachmentsToolsAssistantToolsFileSearchTypeOnly]
type messageAttachmentsToolsAssistantToolsFileSearchTypeOnlyJSON struct {
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *MessageAttachmentsToolsAssistantToolsFileSearchTypeOnly) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r messageAttachmentsToolsAssistantToolsFileSearchTypeOnlyJSON) RawJSON() string {
	return r.raw
}

func (r MessageAttachmentsToolsAssistantToolsFileSearchTypeOnly) implementsMessageAttachmentsTool() {}

// The type of tool being defined: `file_search`
type MessageAttachmentsToolsAssistantToolsFileSearchTypeOnlyType string

const (
	MessageAttachmentsToolsAssistantToolsFileSearchTypeOnlyTypeFileSearch MessageAttachmentsToolsAssistantToolsFileSearchTypeOnlyType = "file_search"
)

func (r MessageAttachmentsToolsAssistantToolsFileSearchTypeOnlyType) IsKnown() bool {
	switch r {
	case MessageAttachmentsToolsAssistantToolsFileSearchTypeOnlyTypeFileSearch:
		return true
	}
	return false
}

// The type of tool being defined: `code_interpreter`
type MessageAttachmentsToolsType string

const (
	MessageAttachmentsToolsTypeCodeInterpreter MessageAttachmentsToolsType = "code_interpreter"
	MessageAttachmentsToolsTypeFileSearch      MessageAttachmentsToolsType = "file_search"
)

func (r MessageAttachmentsToolsType) IsKnown() bool {
	switch r {
	case MessageAttachmentsToolsTypeCodeInterpreter, MessageAttachmentsToolsTypeFileSearch:
		return true
	}
	return false
}

// On an incomplete message, details about why the message is incomplete.
type MessageIncompleteDetails struct {
	// The reason the message is incomplete.
	Reason MessageIncompleteDetailsReason `json:"reason,required"`
	JSON   messageIncompleteDetailsJSON   `json:"-"`
}

// messageIncompleteDetailsJSON contains the JSON metadata for the struct
// [MessageIncompleteDetails]
type messageIncompleteDetailsJSON struct {
	Reason      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *MessageIncompleteDetails) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r messageIncompleteDetailsJSON) RawJSON() string {
	return r.raw
}

// The reason the message is incomplete.
type MessageIncompleteDetailsReason string

const (
	MessageIncompleteDetailsReasonContentFilter MessageIncompleteDetailsReason = "content_filter"
	MessageIncompleteDetailsReasonMaxTokens     MessageIncompleteDetailsReason = "max_tokens"
	MessageIncompleteDetailsReasonRunCancelled  MessageIncompleteDetailsReason = "run_cancelled"
	MessageIncompleteDetailsReasonRunExpired    MessageIncompleteDetailsReason = "run_expired"
	MessageIncompleteDetailsReasonRunFailed     MessageIncompleteDetailsReason = "run_failed"
)

func (r MessageIncompleteDetailsReason) IsKnown() bool {
	switch r {
	case MessageIncompleteDetailsReasonContentFilter, MessageIncompleteDetailsReasonMaxTokens, MessageIncompleteDetailsReasonRunCancelled, MessageIncompleteDetailsReasonRunExpired, MessageIncompleteDetailsReasonRunFailed:
		return true
	}
	return false
}

// The object type, which is always `thread.message`.
type MessageObject string

const (
	MessageObjectThreadMessage MessageObject = "thread.message"
)

func (r MessageObject) IsKnown() bool {
	switch r {
	case MessageObjectThreadMessage:
		return true
	}
	return false
}

// The entity that produced the message. One of `user` or `assistant`.
type MessageRole string

const (
	MessageRoleUser      MessageRole = "user"
	MessageRoleAssistant MessageRole = "assistant"
)

func (r MessageRole) IsKnown() bool {
	switch r {
	case MessageRoleUser, MessageRoleAssistant:
		return true
	}
	return false
}

// The status of the message, which can be either `in_progress`, `incomplete`, or
// `completed`.
type MessageStatus string

const (
	MessageStatusInProgress MessageStatus = "in_progress"
	MessageStatusIncomplete MessageStatus = "incomplete"
	MessageStatusCompleted  MessageStatus = "completed"
)

func (r MessageStatus) IsKnown() bool {
	switch r {
	case MessageStatusInProgress, MessageStatusIncomplete, MessageStatusCompleted:
		return true
	}
	return false
}

// References an image [File](https://platform.openai.com/docs/api-reference/files)
// in the content of a message.
type MessageContent struct {
	// Always `image_file`.
	Type      MessageContentType `json:"type,required"`
	ImageFile ImageFile          `json:"image_file"`
	ImageURL  ImageURL           `json:"image_url"`
	Refusal   string             `json:"refusal"`
	Text      Text               `json:"text"`
	JSON      messageContentJSON `json:"-"`
	union     MessageContentUnion
}

// messageContentJSON contains the JSON metadata for the struct [MessageContent]
type messageContentJSON struct {
	Type        apijson.Field
	ImageFile   apijson.Field
	ImageURL    apijson.Field
	Refusal     apijson.Field
	Text        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r messageContentJSON) RawJSON() string {
	return r.raw
}

func (r *MessageContent) UnmarshalJSON(data []byte) (err error) {
	*r = MessageContent{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [MessageContentUnion] interface which you can cast to the
// specific types for more type safety.
//
// Possible runtime types of the union are [ImageFileContentBlock],
// [ImageURLContentBlock], [TextContentBlock], [RefusalContentBlock].
func (r MessageContent) AsUnion() MessageContentUnion {
	return r.union
}

// References an image [File](https://platform.openai.com/docs/api-reference/files)
// in the content of a message.
//
// Union satisfied by [ImageFileContentBlock], [ImageURLContentBlock],
// [TextContentBlock] or [RefusalContentBlock].
type MessageContentUnion interface {
	implementsMessageContent()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*MessageContentUnion)(nil)).Elem(),
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(ImageFileContentBlock{}),
			DiscriminatorValue: "image_file",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(ImageURLContentBlock{}),
			DiscriminatorValue: "image_url",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(TextContentBlock{}),
			DiscriminatorValue: "text",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(RefusalContentBlock{}),
			DiscriminatorValue: "refusal",
		},
	)
}

// Always `image_file`.
type MessageContentType string

const (
	MessageContentTypeImageFile MessageContentType = "image_file"
	MessageContentTypeImageURL  MessageContentType = "image_url"
	MessageContentTypeText      MessageContentType = "text"
	MessageContentTypeRefusal   MessageContentType = "refusal"
)

func (r MessageContentType) IsKnown() bool {
	switch r {
	case MessageContentTypeImageFile, MessageContentTypeImageURL, MessageContentTypeText, MessageContentTypeRefusal:
		return true
	}
	return false
}

// References an image [File](https://platform.openai.com/docs/api-reference/files)
// in the content of a message.
type MessageContentDelta struct {
	// The index of the content part in the message.
	Index int64 `json:"index,required"`
	// Always `image_file`.
	Type      MessageContentDeltaType `json:"type,required"`
	ImageFile ImageFileDelta          `json:"image_file"`
	ImageURL  ImageURLDelta           `json:"image_url"`
	Refusal   string                  `json:"refusal"`
	Text      TextDelta               `json:"text"`
	JSON      messageContentDeltaJSON `json:"-"`
	union     MessageContentDeltaUnion
}

// messageContentDeltaJSON contains the JSON metadata for the struct
// [MessageContentDelta]
type messageContentDeltaJSON struct {
	Index       apijson.Field
	Type        apijson.Field
	ImageFile   apijson.Field
	ImageURL    apijson.Field
	Refusal     apijson.Field
	Text        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r messageContentDeltaJSON) RawJSON() string {
	return r.raw
}

func (r *MessageContentDelta) UnmarshalJSON(data []byte) (err error) {
	*r = MessageContentDelta{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [MessageContentDeltaUnion] interface which you can cast to the
// specific types for more type safety.
//
// Possible runtime types of the union are [ImageFileDeltaBlock], [TextDeltaBlock],
// [RefusalDeltaBlock], [ImageURLDeltaBlock].
func (r MessageContentDelta) AsUnion() MessageContentDeltaUnion {
	return r.union
}

// References an image [File](https://platform.openai.com/docs/api-reference/files)
// in the content of a message.
//
// Union satisfied by [ImageFileDeltaBlock], [TextDeltaBlock], [RefusalDeltaBlock]
// or [ImageURLDeltaBlock].
type MessageContentDeltaUnion interface {
	implementsMessageContentDelta()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*MessageContentDeltaUnion)(nil)).Elem(),
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(ImageFileDeltaBlock{}),
			DiscriminatorValue: "image_file",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(TextDeltaBlock{}),
			DiscriminatorValue: "text",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(RefusalDeltaBlock{}),
			DiscriminatorValue: "refusal",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(ImageURLDeltaBlock{}),
			DiscriminatorValue: "image_url",
		},
	)
}

// Always `image_file`.
type MessageContentDeltaType string

const (
	MessageContentDeltaTypeImageFile MessageContentDeltaType = "image_file"
	MessageContentDeltaTypeText      MessageContentDeltaType = "text"
	MessageContentDeltaTypeRefusal   MessageContentDeltaType = "refusal"
	MessageContentDeltaTypeImageURL  MessageContentDeltaType = "image_url"
)

func (r MessageContentDeltaType) IsKnown() bool {
	switch r {
	case MessageContentDeltaTypeImageFile, MessageContentDeltaTypeText, MessageContentDeltaTypeRefusal, MessageContentDeltaTypeImageURL:
		return true
	}
	return false
}

// References an image [File](https://platform.openai.com/docs/api-reference/files)
// in the content of a message.
type MessageContentPartParam struct {
	// Always `image_file`.
	Type      param.Field[MessageContentPartParamType] `json:"type,required"`
	ImageFile param.Field[ImageFileParam]              `json:"image_file"`
	ImageURL  param.Field[ImageURLParam]               `json:"image_url"`
	// Text content to be sent to the model
	Text param.Field[string] `json:"text"`
}

func (r MessageContentPartParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r MessageContentPartParam) implementsMessageContentPartParamUnion() {}

// References an image [File](https://platform.openai.com/docs/api-reference/files)
// in the content of a message.
//
// Satisfied by [ImageFileContentBlockParam], [ImageURLContentBlockParam],
// [TextContentBlockParam], [MessageContentPartParam].
type MessageContentPartParamUnion interface {
	implementsMessageContentPartParamUnion()
}

// Always `image_file`.
type MessageContentPartParamType string

const (
	MessageContentPartParamTypeImageFile MessageContentPartParamType = "image_file"
	MessageContentPartParamTypeImageURL  MessageContentPartParamType = "image_url"
	MessageContentPartParamTypeText      MessageContentPartParamType = "text"
)

func (r MessageContentPartParamType) IsKnown() bool {
	switch r {
	case MessageContentPartParamTypeImageFile, MessageContentPartParamTypeImageURL, MessageContentPartParamTypeText:
		return true
	}
	return false
}

type MessageDeleted struct {
	ID      string               `json:"id,required"`
	Deleted bool                 `json:"deleted,required"`
	Object  MessageDeletedObject `json:"object,required"`
	JSON    messageDeletedJSON   `json:"-"`
}

// messageDeletedJSON contains the JSON metadata for the struct [MessageDeleted]
type messageDeletedJSON struct {
	ID          apijson.Field
	Deleted     apijson.Field
	Object      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *MessageDeleted) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r messageDeletedJSON) RawJSON() string {
	return r.raw
}

type MessageDeletedObject string

const (
	MessageDeletedObjectThreadMessageDeleted MessageDeletedObject = "thread.message.deleted"
)

func (r MessageDeletedObject) IsKnown() bool {
	switch r {
	case MessageDeletedObjectThreadMessageDeleted:
		return true
	}
	return false
}

// The delta containing the fields that have changed on the Message.
type MessageDelta struct {
	// The content of the message in array of text and/or images.
	Content []MessageContentDelta `json:"content"`
	// The entity that produced the message. One of `user` or `assistant`.
	Role MessageDeltaRole `json:"role"`
	JSON messageDeltaJSON `json:"-"`
}

// messageDeltaJSON contains the JSON metadata for the struct [MessageDelta]
type messageDeltaJSON struct {
	Content     apijson.Field
	Role        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *MessageDelta) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r messageDeltaJSON) RawJSON() string {
	return r.raw
}

// The entity that produced the message. One of `user` or `assistant`.
type MessageDeltaRole string

const (
	MessageDeltaRoleUser      MessageDeltaRole = "user"
	MessageDeltaRoleAssistant MessageDeltaRole = "assistant"
)

func (r MessageDeltaRole) IsKnown() bool {
	switch r {
	case MessageDeltaRoleUser, MessageDeltaRoleAssistant:
		return true
	}
	return false
}

// Represents a message delta i.e. any changed fields on a message during
// streaming.
type MessageDeltaEvent struct {
	// The identifier of the message, which can be referenced in API endpoints.
	ID string `json:"id,required"`
	// The delta containing the fields that have changed on the Message.
	Delta MessageDelta `json:"delta,required"`
	// The object type, which is always `thread.message.delta`.
	Object MessageDeltaEventObject `json:"object,required"`
	JSON   messageDeltaEventJSON   `json:"-"`
}

// messageDeltaEventJSON contains the JSON metadata for the struct
// [MessageDeltaEvent]
type messageDeltaEventJSON struct {
	ID          apijson.Field
	Delta       apijson.Field
	Object      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *MessageDeltaEvent) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r messageDeltaEventJSON) RawJSON() string {
	return r.raw
}

// The object type, which is always `thread.message.delta`.
type MessageDeltaEventObject string

const (
	MessageDeltaEventObjectThreadMessageDelta MessageDeltaEventObject = "thread.message.delta"
)

func (r MessageDeltaEventObject) IsKnown() bool {
	switch r {
	case MessageDeltaEventObjectThreadMessageDelta:
		return true
	}
	return false
}

// The refusal content generated by the assistant.
type RefusalContentBlock struct {
	Refusal string `json:"refusal,required"`
	// Always `refusal`.
	Type RefusalContentBlockType `json:"type,required"`
	JSON refusalContentBlockJSON `json:"-"`
}

// refusalContentBlockJSON contains the JSON metadata for the struct
// [RefusalContentBlock]
type refusalContentBlockJSON struct {
	Refusal     apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *RefusalContentBlock) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r refusalContentBlockJSON) RawJSON() string {
	return r.raw
}

func (r RefusalContentBlock) implementsMessageContent() {}

// Always `refusal`.
type RefusalContentBlockType string

const (
	RefusalContentBlockTypeRefusal RefusalContentBlockType = "refusal"
)

func (r RefusalContentBlockType) IsKnown() bool {
	switch r {
	case RefusalContentBlockTypeRefusal:
		return true
	}
	return false
}

// The refusal content that is part of a message.
type RefusalDeltaBlock struct {
	// The index of the refusal part in the message.
	Index int64 `json:"index,required"`
	// Always `refusal`.
	Type    RefusalDeltaBlockType `json:"type,required"`
	Refusal string                `json:"refusal"`
	JSON    refusalDeltaBlockJSON `json:"-"`
}

// refusalDeltaBlockJSON contains the JSON metadata for the struct
// [RefusalDeltaBlock]
type refusalDeltaBlockJSON struct {
	Index       apijson.Field
	Type        apijson.Field
	Refusal     apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *RefusalDeltaBlock) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r refusalDeltaBlockJSON) RawJSON() string {
	return r.raw
}

func (r RefusalDeltaBlock) implementsMessageContentDelta() {}

// Always `refusal`.
type RefusalDeltaBlockType string

const (
	RefusalDeltaBlockTypeRefusal RefusalDeltaBlockType = "refusal"
)

func (r RefusalDeltaBlockType) IsKnown() bool {
	switch r {
	case RefusalDeltaBlockTypeRefusal:
		return true
	}
	return false
}

type Text struct {
	Annotations []Annotation `json:"annotations,required"`
	// The data that makes up the text.
	Value string   `json:"value,required"`
	JSON  textJSON `json:"-"`
}

// textJSON contains the JSON metadata for the struct [Text]
type textJSON struct {
	Annotations apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *Text) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r textJSON) RawJSON() string {
	return r.raw
}

// The text content that is part of a message.
type TextContentBlock struct {
	Text Text `json:"text,required"`
	// Always `text`.
	Type TextContentBlockType `json:"type,required"`
	JSON textContentBlockJSON `json:"-"`
}

// textContentBlockJSON contains the JSON metadata for the struct
// [TextContentBlock]
type textContentBlockJSON struct {
	Text        apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *TextContentBlock) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r textContentBlockJSON) RawJSON() string {
	return r.raw
}

func (r TextContentBlock) implementsMessageContent() {}

// Always `text`.
type TextContentBlockType string

const (
	TextContentBlockTypeText TextContentBlockType = "text"
)

func (r TextContentBlockType) IsKnown() bool {
	switch r {
	case TextContentBlockTypeText:
		return true
	}
	return false
}

// The text content that is part of a message.
type TextContentBlockParam struct {
	// Text content to be sent to the model
	Text param.Field[string] `json:"text,required"`
	// Always `text`.
	Type param.Field[TextContentBlockParamType] `json:"type,required"`
}

func (r TextContentBlockParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r TextContentBlockParam) implementsMessageContentPartParamUnion() {}

// Always `text`.
type TextContentBlockParamType string

const (
	TextContentBlockParamTypeText TextContentBlockParamType = "text"
)

func (r TextContentBlockParamType) IsKnown() bool {
	switch r {
	case TextContentBlockParamTypeText:
		return true
	}
	return false
}

type TextDelta struct {
	Annotations []AnnotationDelta `json:"annotations"`
	// The data that makes up the text.
	Value string        `json:"value"`
	JSON  textDeltaJSON `json:"-"`
}

// textDeltaJSON contains the JSON metadata for the struct [TextDelta]
type textDeltaJSON struct {
	Annotations apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *TextDelta) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r textDeltaJSON) RawJSON() string {
	return r.raw
}

// The text content that is part of a message.
type TextDeltaBlock struct {
	// The index of the content part in the message.
	Index int64 `json:"index,required"`
	// Always `text`.
	Type TextDeltaBlockType `json:"type,required"`
	Text TextDelta          `json:"text"`
	JSON textDeltaBlockJSON `json:"-"`
}

// textDeltaBlockJSON contains the JSON metadata for the struct [TextDeltaBlock]
type textDeltaBlockJSON struct {
	Index       apijson.Field
	Type        apijson.Field
	Text        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *TextDeltaBlock) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r textDeltaBlockJSON) RawJSON() string {
	return r.raw
}

func (r TextDeltaBlock) implementsMessageContentDelta() {}

// Always `text`.
type TextDeltaBlockType string

const (
	TextDeltaBlockTypeText TextDeltaBlockType = "text"
)

func (r TextDeltaBlockType) IsKnown() bool {
	switch r {
	case TextDeltaBlockTypeText:
		return true
	}
	return false
}

type BetaThreadMessageNewParams struct {
	// An array of content parts with a defined type, each can be of type `text` or
	// images can be passed with `image_url` or `image_file`. Image types are only
	// supported on
	// [Vision-compatible models](https://platform.openai.com/docs/models).
	Content param.Field[[]MessageContentPartParamUnion] `json:"content,required"`
	// The role of the entity that is creating the message. Allowed values include:
	//
	//   - `user`: Indicates the message is sent by an actual user and should be used in
	//     most cases to represent user-generated messages.
	//   - `assistant`: Indicates the message is generated by the assistant. Use this
	//     value to insert messages from the assistant into the conversation.
	Role param.Field[BetaThreadMessageNewParamsRole] `json:"role,required"`
	// A list of files attached to the message, and the tools they should be added to.
	Attachments param.Field[[]BetaThreadMessageNewParamsAttachment] `json:"attachments"`
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard.
	//
	// Keys are strings with a maximum length of 64 characters. Values are strings with
	// a maximum length of 512 characters.
	Metadata param.Field[shared.MetadataParam] `json:"metadata"`
}

func (r BetaThreadMessageNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
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

func (r BetaThreadMessageNewParamsRole) IsKnown() bool {
	switch r {
	case BetaThreadMessageNewParamsRoleUser, BetaThreadMessageNewParamsRoleAssistant:
		return true
	}
	return false
}

type BetaThreadMessageNewParamsAttachment struct {
	// The ID of the file to attach to the message.
	FileID param.Field[string] `json:"file_id"`
	// The tools to add this file to.
	Tools param.Field[[]BetaThreadMessageNewParamsAttachmentsToolUnion] `json:"tools"`
}

func (r BetaThreadMessageNewParamsAttachment) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type BetaThreadMessageNewParamsAttachmentsTool struct {
	// The type of tool being defined: `code_interpreter`
	Type param.Field[BetaThreadMessageNewParamsAttachmentsToolsType] `json:"type,required"`
}

func (r BetaThreadMessageNewParamsAttachmentsTool) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r BetaThreadMessageNewParamsAttachmentsTool) implementsBetaThreadMessageNewParamsAttachmentsToolUnion() {
}

// Satisfied by [CodeInterpreterToolParam],
// [BetaThreadMessageNewParamsAttachmentsToolsFileSearch],
// [BetaThreadMessageNewParamsAttachmentsTool].
type BetaThreadMessageNewParamsAttachmentsToolUnion interface {
	implementsBetaThreadMessageNewParamsAttachmentsToolUnion()
}

type BetaThreadMessageNewParamsAttachmentsToolsFileSearch struct {
	// The type of tool being defined: `file_search`
	Type param.Field[BetaThreadMessageNewParamsAttachmentsToolsFileSearchType] `json:"type,required"`
}

func (r BetaThreadMessageNewParamsAttachmentsToolsFileSearch) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r BetaThreadMessageNewParamsAttachmentsToolsFileSearch) implementsBetaThreadMessageNewParamsAttachmentsToolUnion() {
}

// The type of tool being defined: `file_search`
type BetaThreadMessageNewParamsAttachmentsToolsFileSearchType string

const (
	BetaThreadMessageNewParamsAttachmentsToolsFileSearchTypeFileSearch BetaThreadMessageNewParamsAttachmentsToolsFileSearchType = "file_search"
)

func (r BetaThreadMessageNewParamsAttachmentsToolsFileSearchType) IsKnown() bool {
	switch r {
	case BetaThreadMessageNewParamsAttachmentsToolsFileSearchTypeFileSearch:
		return true
	}
	return false
}

// The type of tool being defined: `code_interpreter`
type BetaThreadMessageNewParamsAttachmentsToolsType string

const (
	BetaThreadMessageNewParamsAttachmentsToolsTypeCodeInterpreter BetaThreadMessageNewParamsAttachmentsToolsType = "code_interpreter"
	BetaThreadMessageNewParamsAttachmentsToolsTypeFileSearch      BetaThreadMessageNewParamsAttachmentsToolsType = "file_search"
)

func (r BetaThreadMessageNewParamsAttachmentsToolsType) IsKnown() bool {
	switch r {
	case BetaThreadMessageNewParamsAttachmentsToolsTypeCodeInterpreter, BetaThreadMessageNewParamsAttachmentsToolsTypeFileSearch:
		return true
	}
	return false
}

type BetaThreadMessageUpdateParams struct {
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard.
	//
	// Keys are strings with a maximum length of 64 characters. Values are strings with
	// a maximum length of 512 characters.
	Metadata param.Field[shared.MetadataParam] `json:"metadata"`
}

func (r BetaThreadMessageUpdateParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type BetaThreadMessageListParams struct {
	// A cursor for use in pagination. `after` is an object ID that defines your place
	// in the list. For instance, if you make a list request and receive 100 objects,
	// ending with obj_foo, your subsequent call can include after=obj_foo in order to
	// fetch the next page of the list.
	After param.Field[string] `query:"after"`
	// A cursor for use in pagination. `before` is an object ID that defines your place
	// in the list. For instance, if you make a list request and receive 100 objects,
	// starting with obj_foo, your subsequent call can include before=obj_foo in order
	// to fetch the previous page of the list.
	Before param.Field[string] `query:"before"`
	// A limit on the number of objects to be returned. Limit can range between 1 and
	// 100, and the default is 20.
	Limit param.Field[int64] `query:"limit"`
	// Sort order by the `created_at` timestamp of the objects. `asc` for ascending
	// order and `desc` for descending order.
	Order param.Field[BetaThreadMessageListParamsOrder] `query:"order"`
	// Filter messages by the run ID that generated them.
	RunID param.Field[string] `query:"run_id"`
}

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

func (r BetaThreadMessageListParamsOrder) IsKnown() bool {
	switch r {
	case BetaThreadMessageListParamsOrderAsc, BetaThreadMessageListParamsOrderDesc:
		return true
	}
	return false
}

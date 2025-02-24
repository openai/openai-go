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

	"github.com/openai/openai-go/internal/apiform"
	"github.com/openai/openai-go/internal/apijson"
	"github.com/openai/openai-go/internal/apiquery"
	"github.com/openai/openai-go/internal/requestconfig"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/pagination"
	"github.com/openai/openai-go/packages/param"
	"github.com/openai/openai-go/packages/resp"
	"github.com/openai/openai-go/shared/constant"
)

// FileService contains methods and other services that help with interacting with
// the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewFileService] method instead.
type FileService struct {
	Options []option.RequestOption
}

// NewFileService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewFileService(opts ...option.RequestOption) (r FileService) {
	r = FileService{}
	r.Options = opts
	return
}

// Upload a file that can be used across various endpoints. Individual files can be
// up to 512 MB, and the size of all files uploaded by one organization can be up
// to 100 GB.
//
// The Assistants API supports files up to 2 million tokens and of specific file
// types. See the
// [Assistants Tools guide](https://platform.openai.com/docs/assistants/tools) for
// details.
//
// The Fine-tuning API only supports `.jsonl` files. The input also has certain
// required formats for fine-tuning
// [chat](https://platform.openai.com/docs/api-reference/fine-tuning/chat-input) or
// [completions](https://platform.openai.com/docs/api-reference/fine-tuning/completions-input)
// models.
//
// The Batch API only supports `.jsonl` files up to 200 MB in size. The input also
// has a specific required
// [format](https://platform.openai.com/docs/api-reference/batch/request-input).
//
// Please [contact us](https://help.openai.com/) if you need to increase these
// storage limits.
func (r *FileService) New(ctx context.Context, body FileNewParams, opts ...option.RequestOption) (res *FileObject, err error) {
	opts = append(r.Options[:], opts...)
	path := "files"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Returns information about a specific file.
func (r *FileService) Get(ctx context.Context, fileID string, opts ...option.RequestOption) (res *FileObject, err error) {
	opts = append(r.Options[:], opts...)
	if fileID == "" {
		err = errors.New("missing required file_id parameter")
		return
	}
	path := fmt.Sprintf("files/%s", fileID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// Returns a list of files.
func (r *FileService) List(ctx context.Context, query FileListParams, opts ...option.RequestOption) (res *pagination.CursorPage[FileObject], err error) {
	var raw *http.Response
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "files"
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

// Returns a list of files.
func (r *FileService) ListAutoPaging(ctx context.Context, query FileListParams, opts ...option.RequestOption) *pagination.CursorPageAutoPager[FileObject] {
	return pagination.NewCursorPageAutoPager(r.List(ctx, query, opts...))
}

// Delete a file.
func (r *FileService) Delete(ctx context.Context, fileID string, opts ...option.RequestOption) (res *FileDeleted, err error) {
	opts = append(r.Options[:], opts...)
	if fileID == "" {
		err = errors.New("missing required file_id parameter")
		return
	}
	path := fmt.Sprintf("files/%s", fileID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return
}

// Returns the contents of the specified file.
func (r *FileService) Content(ctx context.Context, fileID string, opts ...option.RequestOption) (res *http.Response, err error) {
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "application/binary")}, opts...)
	if fileID == "" {
		err = errors.New("missing required file_id parameter")
		return
	}
	path := fmt.Sprintf("files/%s/content", fileID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

type FileDeleted struct {
	ID      string `json:"id,omitzero,required"`
	Deleted bool   `json:"deleted,omitzero,required"`
	// This field can be elided, and will be automatically set as "file".
	Object constant.File `json:"object,required"`
	JSON   struct {
		ID      resp.Field
		Deleted resp.Field
		Object  resp.Field
		raw     string
	} `json:"-"`
}

func (r FileDeleted) RawJSON() string { return r.JSON.raw }
func (r *FileDeleted) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The `File` object represents a document that has been uploaded to OpenAI.
type FileObject struct {
	// The file identifier, which can be referenced in the API endpoints.
	ID string `json:"id,omitzero,required"`
	// The size of the file, in bytes.
	Bytes int64 `json:"bytes,omitzero,required"`
	// The Unix timestamp (in seconds) for when the file was created.
	CreatedAt int64 `json:"created_at,omitzero,required"`
	// The name of the file.
	Filename string `json:"filename,omitzero,required"`
	// The object type, which is always `file`.
	//
	// This field can be elided, and will be automatically set as "file".
	Object constant.File `json:"object,required"`
	// The intended purpose of the file. Supported values are `assistants`,
	// `assistants_output`, `batch`, `batch_output`, `fine-tune`, `fine-tune-results`
	// and `vision`.
	//
	// Any of "assistants", "assistants_output", "batch", "batch_output", "fine-tune",
	// "fine-tune-results", "vision"
	Purpose string `json:"purpose,omitzero,required"`
	// Deprecated. The current status of the file, which can be either `uploaded`,
	// `processed`, or `error`.
	//
	// Any of "uploaded", "processed", "error"
	//
	// Deprecated: deprecated
	Status string `json:"status,omitzero,required"`
	// Deprecated. For details on why a fine-tuning training file failed validation,
	// see the `error` field on `fine_tuning.job`.
	//
	// Deprecated: deprecated
	StatusDetails string `json:"status_details,omitzero"`
	JSON          struct {
		ID            resp.Field
		Bytes         resp.Field
		CreatedAt     resp.Field
		Filename      resp.Field
		Object        resp.Field
		Purpose       resp.Field
		Status        resp.Field
		StatusDetails resp.Field
		raw           string
	} `json:"-"`
}

func (r FileObject) RawJSON() string { return r.JSON.raw }
func (r *FileObject) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The intended purpose of the file. Supported values are `assistants`,
// `assistants_output`, `batch`, `batch_output`, `fine-tune`, `fine-tune-results`
// and `vision`.
type FileObjectPurpose = string

const (
	FileObjectPurposeAssistants       FileObjectPurpose = "assistants"
	FileObjectPurposeAssistantsOutput FileObjectPurpose = "assistants_output"
	FileObjectPurposeBatch            FileObjectPurpose = "batch"
	FileObjectPurposeBatchOutput      FileObjectPurpose = "batch_output"
	FileObjectPurposeFineTune         FileObjectPurpose = "fine-tune"
	FileObjectPurposeFineTuneResults  FileObjectPurpose = "fine-tune-results"
	FileObjectPurposeVision           FileObjectPurpose = "vision"
)

// Deprecated. The current status of the file, which can be either `uploaded`,
// `processed`, or `error`.
type FileObjectStatus = string

const (
	FileObjectStatusUploaded  FileObjectStatus = "uploaded"
	FileObjectStatusProcessed FileObjectStatus = "processed"
	FileObjectStatusError     FileObjectStatus = "error"
)

// The intended purpose of the uploaded file.
//
// Use "assistants" for
// [Assistants](https://platform.openai.com/docs/api-reference/assistants) and
// [Message](https://platform.openai.com/docs/api-reference/messages) files,
// "vision" for Assistants image file inputs, "batch" for
// [Batch API](https://platform.openai.com/docs/guides/batch), and "fine-tune" for
// [Fine-tuning](https://platform.openai.com/docs/api-reference/fine-tuning).
type FilePurpose string

const (
	FilePurposeAssistants FilePurpose = "assistants"
	FilePurposeBatch      FilePurpose = "batch"
	FilePurposeFineTune   FilePurpose = "fine-tune"
	FilePurposeVision     FilePurpose = "vision"
)

type FileNewParams struct {
	// The File object (not file name) to be uploaded.
	File io.Reader `json:"file,omitzero,required" format:"binary"`
	// The intended purpose of the uploaded file.
	//
	// Use "assistants" for
	// [Assistants](https://platform.openai.com/docs/api-reference/assistants) and
	// [Message](https://platform.openai.com/docs/api-reference/messages) files,
	// "vision" for Assistants image file inputs, "batch" for
	// [Batch API](https://platform.openai.com/docs/guides/batch), and "fine-tune" for
	// [Fine-tuning](https://platform.openai.com/docs/api-reference/fine-tuning).
	//
	// Any of "assistants", "batch", "fine-tune", "vision"
	Purpose FilePurpose `json:"purpose,omitzero,required"`
	apiobject
}

func (f FileNewParams) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

func (r FileNewParams) MarshalMultipart() (data []byte, contentType string, err error) {
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

type FileListParams struct {
	// A cursor for use in pagination. `after` is an object ID that defines your place
	// in the list. For instance, if you make a list request and receive 100 objects,
	// ending with obj_foo, your subsequent call can include after=obj_foo in order to
	// fetch the next page of the list.
	After param.String `query:"after,omitzero"`
	// A limit on the number of objects to be returned. Limit can range between 1 and
	// 10,000, and the default is 10,000.
	Limit param.Int `query:"limit,omitzero"`
	// Sort order by the `created_at` timestamp of the objects. `asc` for ascending
	// order and `desc` for descending order.
	//
	// Any of "asc", "desc"
	Order FileListParamsOrder `query:"order,omitzero"`
	// Only return files with the given purpose.
	Purpose param.String `query:"purpose,omitzero"`
	apiobject
}

func (f FileListParams) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

// URLQuery serializes [FileListParams]'s query parameters as `url.Values`.
func (r FileListParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Sort order by the `created_at` timestamp of the objects. `asc` for ascending
// order and `desc` for descending order.
type FileListParamsOrder string

const (
	FileListParamsOrderAsc  FileListParamsOrder = "asc"
	FileListParamsOrderDesc FileListParamsOrder = "desc"
)

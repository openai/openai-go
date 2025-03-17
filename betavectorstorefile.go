// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/openai/openai-go/internal/apijson"
	"github.com/openai/openai-go/internal/apiquery"
	"github.com/openai/openai-go/internal/param"
	"github.com/openai/openai-go/internal/requestconfig"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/pagination"
)

// BetaVectorStoreFileService contains methods and other services that help with
// interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaVectorStoreFileService] method instead.
type BetaVectorStoreFileService struct {
	Options []option.RequestOption
}

// NewBetaVectorStoreFileService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewBetaVectorStoreFileService(opts ...option.RequestOption) (r *BetaVectorStoreFileService) {
	r = &BetaVectorStoreFileService{}
	r.Options = opts
	return
}

// Create a vector store file by attaching a
// [File](https://platform.openai.com/docs/api-reference/files) to a
// [vector store](https://platform.openai.com/docs/api-reference/vector-stores/object).
func (r *BetaVectorStoreFileService) New(ctx context.Context, vectorStoreID string, body BetaVectorStoreFileNewParams, opts ...option.RequestOption) (res *VectorStoreFile, err error) {
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "assistants=v2")}, opts...)
	if vectorStoreID == "" {
		err = errors.New("missing required vector_store_id parameter")
		return
	}
	path := fmt.Sprintf("vector_stores/%s/files", vectorStoreID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Create a vector store file by attaching a
// [File](https://platform.openai.com/docs/api-reference/files) to a
// [vector store](https://platform.openai.com/docs/api-reference/vector-stores/object).
//
// Polls the API and blocks until the task is complete.
// Default polling interval is 1 second.
func (r *BetaVectorStoreFileService) NewAndPoll(ctx context.Context, vectorStoreId string, body BetaVectorStoreFileNewParams, pollIntervalMs int, opts ...option.RequestOption) (res *VectorStoreFile, err error) {
	file, err := r.New(ctx, vectorStoreId, body, opts...)
	if err != nil {
		return nil, err
	}
	return r.PollStatus(ctx, vectorStoreId, file.ID, pollIntervalMs, opts...)
}

// Upload a file to the `files` API and then attach it to the given vector store.
//
// Note the file will be asynchronously processed (you can use the alternative
// polling helper method to wait for processing to complete).
func (r *BetaVectorStoreFileService) Upload(ctx context.Context, vectorStoreID string, body FileNewParams, opts ...option.RequestOption) (*VectorStoreFile, error) {
	filesService := NewFileService(r.Options...)
	fileObj, err := filesService.New(ctx, body, opts...)
	if err != nil {
		return nil, err
	}

	return r.New(ctx, vectorStoreID, BetaVectorStoreFileNewParams{
		FileID: F(fileObj.ID),
	}, opts...)
}

// Add a file to a vector store and poll until processing is complete.
// Default polling interval is 1 second.
func (r *BetaVectorStoreFileService) UploadAndPoll(ctx context.Context, vectorStoreID string, body FileNewParams, pollIntervalMs int, opts ...option.RequestOption) (*VectorStoreFile, error) {
	res, err := r.Upload(ctx, vectorStoreID, body, opts...)
	if err != nil {
		return nil, err
	}
	return r.PollStatus(ctx, vectorStoreID, res.ID, pollIntervalMs, opts...)
}

// Retrieves a vector store file.
func (r *BetaVectorStoreFileService) Get(ctx context.Context, vectorStoreID string, fileID string, opts ...option.RequestOption) (res *VectorStoreFile, err error) {
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "assistants=v2")}, opts...)
	if vectorStoreID == "" {
		err = errors.New("missing required vector_store_id parameter")
		return
	}
	if fileID == "" {
		err = errors.New("missing required file_id parameter")
		return
	}
	path := fmt.Sprintf("vector_stores/%s/files/%s", vectorStoreID, fileID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// Returns a list of vector store files.
func (r *BetaVectorStoreFileService) List(ctx context.Context, vectorStoreID string, query BetaVectorStoreFileListParams, opts ...option.RequestOption) (res *pagination.CursorPage[VectorStoreFile], err error) {
	var raw *http.Response
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "assistants=v2"), option.WithResponseInto(&raw)}, opts...)
	if vectorStoreID == "" {
		err = errors.New("missing required vector_store_id parameter")
		return
	}
	path := fmt.Sprintf("vector_stores/%s/files", vectorStoreID)
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

// Returns a list of vector store files.
func (r *BetaVectorStoreFileService) ListAutoPaging(ctx context.Context, vectorStoreID string, query BetaVectorStoreFileListParams, opts ...option.RequestOption) *pagination.CursorPageAutoPager[VectorStoreFile] {
	return pagination.NewCursorPageAutoPager(r.List(ctx, vectorStoreID, query, opts...))
}

// Delete a vector store file. This will remove the file from the vector store but
// the file itself will not be deleted. To delete the file, use the
// [delete file](https://platform.openai.com/docs/api-reference/files/delete)
// endpoint.
func (r *BetaVectorStoreFileService) Delete(ctx context.Context, vectorStoreID string, fileID string, opts ...option.RequestOption) (res *VectorStoreFileDeleted, err error) {
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "assistants=v2")}, opts...)
	if vectorStoreID == "" {
		err = errors.New("missing required vector_store_id parameter")
		return
	}
	if fileID == "" {
		err = errors.New("missing required file_id parameter")
		return
	}
	path := fmt.Sprintf("vector_stores/%s/files/%s", vectorStoreID, fileID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return
}

// A list of files attached to a vector store.
type VectorStoreFile struct {
	// The identifier, which can be referenced in API endpoints.
	ID string `json:"id,required"`
	// The Unix timestamp (in seconds) for when the vector store file was created.
	CreatedAt int64 `json:"created_at,required"`
	// The last error associated with this vector store file. Will be `null` if there
	// are no errors.
	LastError VectorStoreFileLastError `json:"last_error,required,nullable"`
	// The object type, which is always `vector_store.file`.
	Object VectorStoreFileObject `json:"object,required"`
	// The status of the vector store file, which can be either `in_progress`,
	// `completed`, `cancelled`, or `failed`. The status `completed` indicates that the
	// vector store file is ready for use.
	Status VectorStoreFileStatus `json:"status,required"`
	// The total vector store usage in bytes. Note that this may be different from the
	// original file size.
	UsageBytes int64 `json:"usage_bytes,required"`
	// The ID of the
	// [vector store](https://platform.openai.com/docs/api-reference/vector-stores/object)
	// that the [File](https://platform.openai.com/docs/api-reference/files) is
	// attached to.
	VectorStoreID string `json:"vector_store_id,required"`
	// The strategy used to chunk the file.
	ChunkingStrategy FileChunkingStrategy `json:"chunking_strategy"`
	JSON             vectorStoreFileJSON  `json:"-"`
}

// vectorStoreFileJSON contains the JSON metadata for the struct [VectorStoreFile]
type vectorStoreFileJSON struct {
	ID               apijson.Field
	CreatedAt        apijson.Field
	LastError        apijson.Field
	Object           apijson.Field
	Status           apijson.Field
	UsageBytes       apijson.Field
	VectorStoreID    apijson.Field
	ChunkingStrategy apijson.Field
	raw              string
	ExtraFields      map[string]apijson.Field
}

func (r *VectorStoreFile) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r vectorStoreFileJSON) RawJSON() string {
	return r.raw
}

// The last error associated with this vector store file. Will be `null` if there
// are no errors.
type VectorStoreFileLastError struct {
	// One of `server_error` or `rate_limit_exceeded`.
	Code VectorStoreFileLastErrorCode `json:"code,required"`
	// A human-readable description of the error.
	Message string                       `json:"message,required"`
	JSON    vectorStoreFileLastErrorJSON `json:"-"`
}

// vectorStoreFileLastErrorJSON contains the JSON metadata for the struct
// [VectorStoreFileLastError]
type vectorStoreFileLastErrorJSON struct {
	Code        apijson.Field
	Message     apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *VectorStoreFileLastError) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r vectorStoreFileLastErrorJSON) RawJSON() string {
	return r.raw
}

// One of `server_error` or `rate_limit_exceeded`.
type VectorStoreFileLastErrorCode string

const (
	VectorStoreFileLastErrorCodeServerError     VectorStoreFileLastErrorCode = "server_error"
	VectorStoreFileLastErrorCodeUnsupportedFile VectorStoreFileLastErrorCode = "unsupported_file"
	VectorStoreFileLastErrorCodeInvalidFile     VectorStoreFileLastErrorCode = "invalid_file"
)

func (r VectorStoreFileLastErrorCode) IsKnown() bool {
	switch r {
	case VectorStoreFileLastErrorCodeServerError, VectorStoreFileLastErrorCodeUnsupportedFile, VectorStoreFileLastErrorCodeInvalidFile:
		return true
	}
	return false
}

// The object type, which is always `vector_store.file`.
type VectorStoreFileObject string

const (
	VectorStoreFileObjectVectorStoreFile VectorStoreFileObject = "vector_store.file"
)

func (r VectorStoreFileObject) IsKnown() bool {
	switch r {
	case VectorStoreFileObjectVectorStoreFile:
		return true
	}
	return false
}

// The status of the vector store file, which can be either `in_progress`,
// `completed`, `cancelled`, or `failed`. The status `completed` indicates that the
// vector store file is ready for use.
type VectorStoreFileStatus string

const (
	VectorStoreFileStatusInProgress VectorStoreFileStatus = "in_progress"
	VectorStoreFileStatusCompleted  VectorStoreFileStatus = "completed"
	VectorStoreFileStatusCancelled  VectorStoreFileStatus = "cancelled"
	VectorStoreFileStatusFailed     VectorStoreFileStatus = "failed"
)

func (r VectorStoreFileStatus) IsKnown() bool {
	switch r {
	case VectorStoreFileStatusInProgress, VectorStoreFileStatusCompleted, VectorStoreFileStatusCancelled, VectorStoreFileStatusFailed:
		return true
	}
	return false
}

type VectorStoreFileDeleted struct {
	ID      string                       `json:"id,required"`
	Deleted bool                         `json:"deleted,required"`
	Object  VectorStoreFileDeletedObject `json:"object,required"`
	JSON    vectorStoreFileDeletedJSON   `json:"-"`
}

// vectorStoreFileDeletedJSON contains the JSON metadata for the struct
// [VectorStoreFileDeleted]
type vectorStoreFileDeletedJSON struct {
	ID          apijson.Field
	Deleted     apijson.Field
	Object      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *VectorStoreFileDeleted) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r vectorStoreFileDeletedJSON) RawJSON() string {
	return r.raw
}

type VectorStoreFileDeletedObject string

const (
	VectorStoreFileDeletedObjectVectorStoreFileDeleted VectorStoreFileDeletedObject = "vector_store.file.deleted"
)

func (r VectorStoreFileDeletedObject) IsKnown() bool {
	switch r {
	case VectorStoreFileDeletedObjectVectorStoreFileDeleted:
		return true
	}
	return false
}

type BetaVectorStoreFileNewParams struct {
	// A [File](https://platform.openai.com/docs/api-reference/files) ID that the
	// vector store should use. Useful for tools like `file_search` that can access
	// files.
	FileID param.Field[string] `json:"file_id,required"`
	// The chunking strategy used to chunk the file(s). If not set, will use the `auto`
	// strategy. Only applicable if `file_ids` is non-empty.
	ChunkingStrategy param.Field[FileChunkingStrategyParamUnion] `json:"chunking_strategy"`
}

func (r BetaVectorStoreFileNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type BetaVectorStoreFileListParams struct {
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
	// Filter by file status. One of `in_progress`, `completed`, `failed`, `cancelled`.
	Filter param.Field[BetaVectorStoreFileListParamsFilter] `query:"filter"`
	// A limit on the number of objects to be returned. Limit can range between 1 and
	// 100, and the default is 20.
	Limit param.Field[int64] `query:"limit"`
	// Sort order by the `created_at` timestamp of the objects. `asc` for ascending
	// order and `desc` for descending order.
	Order param.Field[BetaVectorStoreFileListParamsOrder] `query:"order"`
}

// URLQuery serializes [BetaVectorStoreFileListParams]'s query parameters as
// `url.Values`.
func (r BetaVectorStoreFileListParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Filter by file status. One of `in_progress`, `completed`, `failed`, `cancelled`.
type BetaVectorStoreFileListParamsFilter string

const (
	BetaVectorStoreFileListParamsFilterInProgress BetaVectorStoreFileListParamsFilter = "in_progress"
	BetaVectorStoreFileListParamsFilterCompleted  BetaVectorStoreFileListParamsFilter = "completed"
	BetaVectorStoreFileListParamsFilterFailed     BetaVectorStoreFileListParamsFilter = "failed"
	BetaVectorStoreFileListParamsFilterCancelled  BetaVectorStoreFileListParamsFilter = "cancelled"
)

func (r BetaVectorStoreFileListParamsFilter) IsKnown() bool {
	switch r {
	case BetaVectorStoreFileListParamsFilterInProgress, BetaVectorStoreFileListParamsFilterCompleted, BetaVectorStoreFileListParamsFilterFailed, BetaVectorStoreFileListParamsFilterCancelled:
		return true
	}
	return false
}

// Sort order by the `created_at` timestamp of the objects. `asc` for ascending
// order and `desc` for descending order.
type BetaVectorStoreFileListParamsOrder string

const (
	BetaVectorStoreFileListParamsOrderAsc  BetaVectorStoreFileListParamsOrder = "asc"
	BetaVectorStoreFileListParamsOrderDesc BetaVectorStoreFileListParamsOrder = "desc"
)

func (r BetaVectorStoreFileListParamsOrder) IsKnown() bool {
	switch r {
	case BetaVectorStoreFileListParamsOrderAsc, BetaVectorStoreFileListParamsOrderDesc:
		return true
	}
	return false
}

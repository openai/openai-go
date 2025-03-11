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

// VectorStoreFileBatchService contains methods and other services that help with
// interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewVectorStoreFileBatchService] method instead.
type VectorStoreFileBatchService struct {
	Options []option.RequestOption
}

// NewVectorStoreFileBatchService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewVectorStoreFileBatchService(opts ...option.RequestOption) (r *VectorStoreFileBatchService) {
	r = &VectorStoreFileBatchService{}
	r.Options = opts
	return
}

// Create a vector store file batch.
func (r *VectorStoreFileBatchService) New(ctx context.Context, vectorStoreID string, body VectorStoreFileBatchNewParams, opts ...option.RequestOption) (res *VectorStoreFileBatch, err error) {
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "assistants=v2")}, opts...)
	if vectorStoreID == "" {
		err = errors.New("missing required vector_store_id parameter")
		return
	}
	path := fmt.Sprintf("vector_stores/%s/file_batches", vectorStoreID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Retrieves a vector store file batch.
func (r *VectorStoreFileBatchService) Get(ctx context.Context, vectorStoreID string, batchID string, opts ...option.RequestOption) (res *VectorStoreFileBatch, err error) {
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "assistants=v2")}, opts...)
	if vectorStoreID == "" {
		err = errors.New("missing required vector_store_id parameter")
		return
	}
	if batchID == "" {
		err = errors.New("missing required batch_id parameter")
		return
	}
	path := fmt.Sprintf("vector_stores/%s/file_batches/%s", vectorStoreID, batchID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// Cancel a vector store file batch. This attempts to cancel the processing of
// files in this batch as soon as possible.
func (r *VectorStoreFileBatchService) Cancel(ctx context.Context, vectorStoreID string, batchID string, opts ...option.RequestOption) (res *VectorStoreFileBatch, err error) {
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "assistants=v2")}, opts...)
	if vectorStoreID == "" {
		err = errors.New("missing required vector_store_id parameter")
		return
	}
	if batchID == "" {
		err = errors.New("missing required batch_id parameter")
		return
	}
	path := fmt.Sprintf("vector_stores/%s/file_batches/%s/cancel", vectorStoreID, batchID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return
}

// Returns a list of vector store files in a batch.
func (r *VectorStoreFileBatchService) ListFiles(ctx context.Context, vectorStoreID string, batchID string, query VectorStoreFileBatchListFilesParams, opts ...option.RequestOption) (res *pagination.CursorPage[VectorStoreFile], err error) {
	var raw *http.Response
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "assistants=v2"), option.WithResponseInto(&raw)}, opts...)
	if vectorStoreID == "" {
		err = errors.New("missing required vector_store_id parameter")
		return
	}
	if batchID == "" {
		err = errors.New("missing required batch_id parameter")
		return
	}
	path := fmt.Sprintf("vector_stores/%s/file_batches/%s/files", vectorStoreID, batchID)
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

// Returns a list of vector store files in a batch.
func (r *VectorStoreFileBatchService) ListFilesAutoPaging(ctx context.Context, vectorStoreID string, batchID string, query VectorStoreFileBatchListFilesParams, opts ...option.RequestOption) *pagination.CursorPageAutoPager[VectorStoreFile] {
	return pagination.NewCursorPageAutoPager(r.ListFiles(ctx, vectorStoreID, batchID, query, opts...))
}

// A batch of files attached to a vector store.
type VectorStoreFileBatch struct {
	// The identifier, which can be referenced in API endpoints.
	ID string `json:"id,required"`
	// The Unix timestamp (in seconds) for when the vector store files batch was
	// created.
	CreatedAt  int64                          `json:"created_at,required"`
	FileCounts VectorStoreFileBatchFileCounts `json:"file_counts,required"`
	// The object type, which is always `vector_store.file_batch`.
	Object VectorStoreFileBatchObject `json:"object,required"`
	// The status of the vector store files batch, which can be either `in_progress`,
	// `completed`, `cancelled` or `failed`.
	Status VectorStoreFileBatchStatus `json:"status,required"`
	// The ID of the
	// [vector store](https://platform.openai.com/docs/api-reference/vector-stores/object)
	// that the [File](https://platform.openai.com/docs/api-reference/files) is
	// attached to.
	VectorStoreID string                   `json:"vector_store_id,required"`
	JSON          vectorStoreFileBatchJSON `json:"-"`
}

// vectorStoreFileBatchJSON contains the JSON metadata for the struct
// [VectorStoreFileBatch]
type vectorStoreFileBatchJSON struct {
	ID            apijson.Field
	CreatedAt     apijson.Field
	FileCounts    apijson.Field
	Object        apijson.Field
	Status        apijson.Field
	VectorStoreID apijson.Field
	raw           string
	ExtraFields   map[string]apijson.Field
}

func (r *VectorStoreFileBatch) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r vectorStoreFileBatchJSON) RawJSON() string {
	return r.raw
}

type VectorStoreFileBatchFileCounts struct {
	// The number of files that where cancelled.
	Cancelled int64 `json:"cancelled,required"`
	// The number of files that have been processed.
	Completed int64 `json:"completed,required"`
	// The number of files that have failed to process.
	Failed int64 `json:"failed,required"`
	// The number of files that are currently being processed.
	InProgress int64 `json:"in_progress,required"`
	// The total number of files.
	Total int64                              `json:"total,required"`
	JSON  vectorStoreFileBatchFileCountsJSON `json:"-"`
}

// vectorStoreFileBatchFileCountsJSON contains the JSON metadata for the struct
// [VectorStoreFileBatchFileCounts]
type vectorStoreFileBatchFileCountsJSON struct {
	Cancelled   apijson.Field
	Completed   apijson.Field
	Failed      apijson.Field
	InProgress  apijson.Field
	Total       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *VectorStoreFileBatchFileCounts) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r vectorStoreFileBatchFileCountsJSON) RawJSON() string {
	return r.raw
}

// The object type, which is always `vector_store.file_batch`.
type VectorStoreFileBatchObject string

const (
	VectorStoreFileBatchObjectVectorStoreFilesBatch VectorStoreFileBatchObject = "vector_store.files_batch"
)

func (r VectorStoreFileBatchObject) IsKnown() bool {
	switch r {
	case VectorStoreFileBatchObjectVectorStoreFilesBatch:
		return true
	}
	return false
}

// The status of the vector store files batch, which can be either `in_progress`,
// `completed`, `cancelled` or `failed`.
type VectorStoreFileBatchStatus string

const (
	VectorStoreFileBatchStatusInProgress VectorStoreFileBatchStatus = "in_progress"
	VectorStoreFileBatchStatusCompleted  VectorStoreFileBatchStatus = "completed"
	VectorStoreFileBatchStatusCancelled  VectorStoreFileBatchStatus = "cancelled"
	VectorStoreFileBatchStatusFailed     VectorStoreFileBatchStatus = "failed"
)

func (r VectorStoreFileBatchStatus) IsKnown() bool {
	switch r {
	case VectorStoreFileBatchStatusInProgress, VectorStoreFileBatchStatusCompleted, VectorStoreFileBatchStatusCancelled, VectorStoreFileBatchStatusFailed:
		return true
	}
	return false
}

type VectorStoreFileBatchNewParams struct {
	// A list of [File](https://platform.openai.com/docs/api-reference/files) IDs that
	// the vector store should use. Useful for tools like `file_search` that can access
	// files.
	FileIDs param.Field[[]string] `json:"file_ids,required"`
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard. Keys are strings with a maximum
	// length of 64 characters. Values are strings with a maximum length of 512
	// characters, booleans, or numbers.
	Attributes param.Field[map[string]VectorStoreFileBatchNewParamsAttributesUnion] `json:"attributes"`
	// The chunking strategy used to chunk the file(s). If not set, will use the `auto`
	// strategy. Only applicable if `file_ids` is non-empty.
	ChunkingStrategy param.Field[FileChunkingStrategyParamUnion] `json:"chunking_strategy"`
}

func (r VectorStoreFileBatchNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Satisfied by [shared.UnionString], [shared.UnionFloat], [shared.UnionBool].
type VectorStoreFileBatchNewParamsAttributesUnion interface {
	ImplementsVectorStoreFileBatchNewParamsAttributesUnion()
}

type VectorStoreFileBatchListFilesParams struct {
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
	Filter param.Field[VectorStoreFileBatchListFilesParamsFilter] `query:"filter"`
	// A limit on the number of objects to be returned. Limit can range between 1 and
	// 100, and the default is 20.
	Limit param.Field[int64] `query:"limit"`
	// Sort order by the `created_at` timestamp of the objects. `asc` for ascending
	// order and `desc` for descending order.
	Order param.Field[VectorStoreFileBatchListFilesParamsOrder] `query:"order"`
}

// URLQuery serializes [VectorStoreFileBatchListFilesParams]'s query parameters as
// `url.Values`.
func (r VectorStoreFileBatchListFilesParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Filter by file status. One of `in_progress`, `completed`, `failed`, `cancelled`.
type VectorStoreFileBatchListFilesParamsFilter string

const (
	VectorStoreFileBatchListFilesParamsFilterInProgress VectorStoreFileBatchListFilesParamsFilter = "in_progress"
	VectorStoreFileBatchListFilesParamsFilterCompleted  VectorStoreFileBatchListFilesParamsFilter = "completed"
	VectorStoreFileBatchListFilesParamsFilterFailed     VectorStoreFileBatchListFilesParamsFilter = "failed"
	VectorStoreFileBatchListFilesParamsFilterCancelled  VectorStoreFileBatchListFilesParamsFilter = "cancelled"
)

func (r VectorStoreFileBatchListFilesParamsFilter) IsKnown() bool {
	switch r {
	case VectorStoreFileBatchListFilesParamsFilterInProgress, VectorStoreFileBatchListFilesParamsFilterCompleted, VectorStoreFileBatchListFilesParamsFilterFailed, VectorStoreFileBatchListFilesParamsFilterCancelled:
		return true
	}
	return false
}

// Sort order by the `created_at` timestamp of the objects. `asc` for ascending
// order and `desc` for descending order.
type VectorStoreFileBatchListFilesParamsOrder string

const (
	VectorStoreFileBatchListFilesParamsOrderAsc  VectorStoreFileBatchListFilesParamsOrder = "asc"
	VectorStoreFileBatchListFilesParamsOrderDesc VectorStoreFileBatchListFilesParamsOrder = "desc"
)

func (r VectorStoreFileBatchListFilesParamsOrder) IsKnown() bool {
	switch r {
	case VectorStoreFileBatchListFilesParamsOrderAsc, VectorStoreFileBatchListFilesParamsOrderDesc:
		return true
	}
	return false
}

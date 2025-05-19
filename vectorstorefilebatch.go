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
	"github.com/openai/openai-go/internal/requestconfig"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/pagination"
	"github.com/openai/openai-go/packages/param"
	"github.com/openai/openai-go/packages/respjson"
	"github.com/openai/openai-go/shared/constant"
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
func NewVectorStoreFileBatchService(opts ...option.RequestOption) (r VectorStoreFileBatchService) {
	r = VectorStoreFileBatchService{}
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
func (r *VectorStoreFileBatchService) Get(ctx context.Context, batchID string, query VectorStoreFileBatchGetParams, opts ...option.RequestOption) (res *VectorStoreFileBatch, err error) {
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "assistants=v2")}, opts...)
	if query.VectorStoreID == "" {
		err = errors.New("missing required vector_store_id parameter")
		return
	}
	if batchID == "" {
		err = errors.New("missing required batch_id parameter")
		return
	}
	path := fmt.Sprintf("vector_stores/%s/file_batches/%s", query.VectorStoreID, batchID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// Cancel a vector store file batch. This attempts to cancel the processing of
// files in this batch as soon as possible.
func (r *VectorStoreFileBatchService) Cancel(ctx context.Context, batchID string, body VectorStoreFileBatchCancelParams, opts ...option.RequestOption) (res *VectorStoreFileBatch, err error) {
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "assistants=v2")}, opts...)
	if body.VectorStoreID == "" {
		err = errors.New("missing required vector_store_id parameter")
		return
	}
	if batchID == "" {
		err = errors.New("missing required batch_id parameter")
		return
	}
	path := fmt.Sprintf("vector_stores/%s/file_batches/%s/cancel", body.VectorStoreID, batchID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return
}

// Returns a list of vector store files in a batch.
func (r *VectorStoreFileBatchService) ListFiles(ctx context.Context, batchID string, params VectorStoreFileBatchListFilesParams, opts ...option.RequestOption) (res *pagination.CursorPage[VectorStoreFile], err error) {
	var raw *http.Response
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "assistants=v2"), option.WithResponseInto(&raw)}, opts...)
	if params.VectorStoreID == "" {
		err = errors.New("missing required vector_store_id parameter")
		return
	}
	if batchID == "" {
		err = errors.New("missing required batch_id parameter")
		return
	}
	path := fmt.Sprintf("vector_stores/%s/file_batches/%s/files", params.VectorStoreID, batchID)
	cfg, err := requestconfig.NewRequestConfig(ctx, http.MethodGet, path, params, &res, opts...)
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
func (r *VectorStoreFileBatchService) ListFilesAutoPaging(ctx context.Context, batchID string, params VectorStoreFileBatchListFilesParams, opts ...option.RequestOption) *pagination.CursorPageAutoPager[VectorStoreFile] {
	return pagination.NewCursorPageAutoPager(r.ListFiles(ctx, batchID, params, opts...))
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
	Object constant.VectorStoreFilesBatch `json:"object,required"`
	// The status of the vector store files batch, which can be either `in_progress`,
	// `completed`, `cancelled` or `failed`.
	//
	// Any of "in_progress", "completed", "cancelled", "failed".
	Status VectorStoreFileBatchStatus `json:"status,required"`
	// The ID of the
	// [vector store](https://platform.openai.com/docs/api-reference/vector-stores/object)
	// that the [File](https://platform.openai.com/docs/api-reference/files) is
	// attached to.
	VectorStoreID string `json:"vector_store_id,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID            respjson.Field
		CreatedAt     respjson.Field
		FileCounts    respjson.Field
		Object        respjson.Field
		Status        respjson.Field
		VectorStoreID respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r VectorStoreFileBatch) RawJSON() string { return r.JSON.raw }
func (r *VectorStoreFileBatch) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
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
	Total int64 `json:"total,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Cancelled   respjson.Field
		Completed   respjson.Field
		Failed      respjson.Field
		InProgress  respjson.Field
		Total       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r VectorStoreFileBatchFileCounts) RawJSON() string { return r.JSON.raw }
func (r *VectorStoreFileBatchFileCounts) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
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

type VectorStoreFileBatchNewParams struct {
	// A list of [File](https://platform.openai.com/docs/api-reference/files) IDs that
	// the vector store should use. Useful for tools like `file_search` that can access
	// files.
	FileIDs []string `json:"file_ids,omitzero,required"`
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard. Keys are strings with a maximum
	// length of 64 characters. Values are strings with a maximum length of 512
	// characters, booleans, or numbers.
	Attributes map[string]VectorStoreFileBatchNewParamsAttributeUnion `json:"attributes,omitzero"`
	// The chunking strategy used to chunk the file(s). If not set, will use the `auto`
	// strategy. Only applicable if `file_ids` is non-empty.
	ChunkingStrategy FileChunkingStrategyParamUnion `json:"chunking_strategy,omitzero"`
	paramObj
}

func (r VectorStoreFileBatchNewParams) MarshalJSON() (data []byte, err error) {
	type shadow VectorStoreFileBatchNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *VectorStoreFileBatchNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type VectorStoreFileBatchNewParamsAttributeUnion struct {
	OfString param.Opt[string]  `json:",omitzero,inline"`
	OfFloat  param.Opt[float64] `json:",omitzero,inline"`
	OfBool   param.Opt[bool]    `json:",omitzero,inline"`
	paramUnion
}

func (u VectorStoreFileBatchNewParamsAttributeUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion[VectorStoreFileBatchNewParamsAttributeUnion](u.OfString, u.OfFloat, u.OfBool)
}
func (u *VectorStoreFileBatchNewParamsAttributeUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *VectorStoreFileBatchNewParamsAttributeUnion) asAny() any {
	if !param.IsOmitted(u.OfString) {
		return &u.OfString.Value
	} else if !param.IsOmitted(u.OfFloat) {
		return &u.OfFloat.Value
	} else if !param.IsOmitted(u.OfBool) {
		return &u.OfBool.Value
	}
	return nil
}

type VectorStoreFileBatchGetParams struct {
	VectorStoreID string `path:"vector_store_id,required" json:"-"`
	paramObj
}

type VectorStoreFileBatchCancelParams struct {
	VectorStoreID string `path:"vector_store_id,required" json:"-"`
	paramObj
}

type VectorStoreFileBatchListFilesParams struct {
	VectorStoreID string `path:"vector_store_id,required" json:"-"`
	// A cursor for use in pagination. `after` is an object ID that defines your place
	// in the list. For instance, if you make a list request and receive 100 objects,
	// ending with obj_foo, your subsequent call can include after=obj_foo in order to
	// fetch the next page of the list.
	After param.Opt[string] `query:"after,omitzero" json:"-"`
	// A cursor for use in pagination. `before` is an object ID that defines your place
	// in the list. For instance, if you make a list request and receive 100 objects,
	// starting with obj_foo, your subsequent call can include before=obj_foo in order
	// to fetch the previous page of the list.
	Before param.Opt[string] `query:"before,omitzero" json:"-"`
	// A limit on the number of objects to be returned. Limit can range between 1 and
	// 100, and the default is 20.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Filter by file status. One of `in_progress`, `completed`, `failed`, `cancelled`.
	//
	// Any of "in_progress", "completed", "failed", "cancelled".
	Filter VectorStoreFileBatchListFilesParamsFilter `query:"filter,omitzero" json:"-"`
	// Sort order by the `created_at` timestamp of the objects. `asc` for ascending
	// order and `desc` for descending order.
	//
	// Any of "asc", "desc".
	Order VectorStoreFileBatchListFilesParamsOrder `query:"order,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [VectorStoreFileBatchListFilesParams]'s query parameters as
// `url.Values`.
func (r VectorStoreFileBatchListFilesParams) URLQuery() (v url.Values, err error) {
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

// Sort order by the `created_at` timestamp of the objects. `asc` for ascending
// order and `desc` for descending order.
type VectorStoreFileBatchListFilesParamsOrder string

const (
	VectorStoreFileBatchListFilesParamsOrderAsc  VectorStoreFileBatchListFilesParamsOrder = "asc"
	VectorStoreFileBatchListFilesParamsOrderDesc VectorStoreFileBatchListFilesParamsOrder = "desc"
)

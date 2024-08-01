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
	"github.com/openai/openai-go/internal/pagination"
	"github.com/openai/openai-go/internal/param"
	"github.com/openai/openai-go/internal/requestconfig"
	"github.com/openai/openai-go/option"
)

// BetaVectorStoreService contains methods and other services that help with
// interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaVectorStoreService] method instead.
type BetaVectorStoreService struct {
	Options     []option.RequestOption
	Files       *BetaVectorStoreFileService
	FileBatches *BetaVectorStoreFileBatchService
}

// NewBetaVectorStoreService generates a new service that applies the given options
// to each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewBetaVectorStoreService(opts ...option.RequestOption) (r *BetaVectorStoreService) {
	r = &BetaVectorStoreService{}
	r.Options = opts
	r.Files = NewBetaVectorStoreFileService(opts...)
	r.FileBatches = NewBetaVectorStoreFileBatchService(opts...)
	return
}

// Create a vector store.
func (r *BetaVectorStoreService) New(ctx context.Context, body BetaVectorStoreNewParams, opts ...option.RequestOption) (res *VectorStore, err error) {
	opts = append(r.Options[:], opts...)
	path := "vector_stores"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Retrieves a vector store.
func (r *BetaVectorStoreService) Get(ctx context.Context, vectorStoreID string, opts ...option.RequestOption) (res *VectorStore, err error) {
	opts = append(r.Options[:], opts...)
	if vectorStoreID == "" {
		err = errors.New("missing required vector_store_id parameter")
		return
	}
	path := fmt.Sprintf("vector_stores/%s", vectorStoreID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// Modifies a vector store.
func (r *BetaVectorStoreService) Update(ctx context.Context, vectorStoreID string, body BetaVectorStoreUpdateParams, opts ...option.RequestOption) (res *VectorStore, err error) {
	opts = append(r.Options[:], opts...)
	if vectorStoreID == "" {
		err = errors.New("missing required vector_store_id parameter")
		return
	}
	path := fmt.Sprintf("vector_stores/%s", vectorStoreID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Returns a list of vector stores.
func (r *BetaVectorStoreService) List(ctx context.Context, query BetaVectorStoreListParams, opts ...option.RequestOption) (res *pagination.CursorPage[VectorStore], err error) {
	var raw *http.Response
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "vector_stores"
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

// Returns a list of vector stores.
func (r *BetaVectorStoreService) ListAutoPaging(ctx context.Context, query BetaVectorStoreListParams, opts ...option.RequestOption) *pagination.CursorPageAutoPager[VectorStore] {
	return pagination.NewCursorPageAutoPager(r.List(ctx, query, opts...))
}

// Delete a vector store.
func (r *BetaVectorStoreService) Delete(ctx context.Context, vectorStoreID string, opts ...option.RequestOption) (res *VectorStoreDeleted, err error) {
	opts = append(r.Options[:], opts...)
	if vectorStoreID == "" {
		err = errors.New("missing required vector_store_id parameter")
		return
	}
	path := fmt.Sprintf("vector_stores/%s", vectorStoreID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return
}

// A vector store is a collection of processed files can be used by the
// `file_search` tool.
type VectorStore struct {
	// The identifier, which can be referenced in API endpoints.
	ID string `json:"id,required"`
	// The Unix timestamp (in seconds) for when the vector store was created.
	CreatedAt  int64                 `json:"created_at,required"`
	FileCounts VectorStoreFileCounts `json:"file_counts,required"`
	// The Unix timestamp (in seconds) for when the vector store was last active.
	LastActiveAt int64 `json:"last_active_at,required,nullable"`
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format. Keys
	// can be a maximum of 64 characters long and values can be a maxium of 512
	// characters long.
	Metadata interface{} `json:"metadata,required,nullable"`
	// The name of the vector store.
	Name string `json:"name,required"`
	// The object type, which is always `vector_store`.
	Object VectorStoreObject `json:"object,required"`
	// The status of the vector store, which can be either `expired`, `in_progress`, or
	// `completed`. A status of `completed` indicates that the vector store is ready
	// for use.
	Status VectorStoreStatus `json:"status,required"`
	// The total number of bytes used by the files in the vector store.
	UsageBytes int64 `json:"usage_bytes,required"`
	// The expiration policy for a vector store.
	ExpiresAfter VectorStoreExpiresAfter `json:"expires_after"`
	// The Unix timestamp (in seconds) for when the vector store will expire.
	ExpiresAt int64           `json:"expires_at,nullable"`
	JSON      vectorStoreJSON `json:"-"`
}

// vectorStoreJSON contains the JSON metadata for the struct [VectorStore]
type vectorStoreJSON struct {
	ID           apijson.Field
	CreatedAt    apijson.Field
	FileCounts   apijson.Field
	LastActiveAt apijson.Field
	Metadata     apijson.Field
	Name         apijson.Field
	Object       apijson.Field
	Status       apijson.Field
	UsageBytes   apijson.Field
	ExpiresAfter apijson.Field
	ExpiresAt    apijson.Field
	raw          string
	ExtraFields  map[string]apijson.Field
}

func (r *VectorStore) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r vectorStoreJSON) RawJSON() string {
	return r.raw
}

type VectorStoreFileCounts struct {
	// The number of files that were cancelled.
	Cancelled int64 `json:"cancelled,required"`
	// The number of files that have been successfully processed.
	Completed int64 `json:"completed,required"`
	// The number of files that have failed to process.
	Failed int64 `json:"failed,required"`
	// The number of files that are currently being processed.
	InProgress int64 `json:"in_progress,required"`
	// The total number of files.
	Total int64                     `json:"total,required"`
	JSON  vectorStoreFileCountsJSON `json:"-"`
}

// vectorStoreFileCountsJSON contains the JSON metadata for the struct
// [VectorStoreFileCounts]
type vectorStoreFileCountsJSON struct {
	Cancelled   apijson.Field
	Completed   apijson.Field
	Failed      apijson.Field
	InProgress  apijson.Field
	Total       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *VectorStoreFileCounts) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r vectorStoreFileCountsJSON) RawJSON() string {
	return r.raw
}

// The object type, which is always `vector_store`.
type VectorStoreObject string

const (
	VectorStoreObjectVectorStore VectorStoreObject = "vector_store"
)

func (r VectorStoreObject) IsKnown() bool {
	switch r {
	case VectorStoreObjectVectorStore:
		return true
	}
	return false
}

// The status of the vector store, which can be either `expired`, `in_progress`, or
// `completed`. A status of `completed` indicates that the vector store is ready
// for use.
type VectorStoreStatus string

const (
	VectorStoreStatusExpired    VectorStoreStatus = "expired"
	VectorStoreStatusInProgress VectorStoreStatus = "in_progress"
	VectorStoreStatusCompleted  VectorStoreStatus = "completed"
)

func (r VectorStoreStatus) IsKnown() bool {
	switch r {
	case VectorStoreStatusExpired, VectorStoreStatusInProgress, VectorStoreStatusCompleted:
		return true
	}
	return false
}

// The expiration policy for a vector store.
type VectorStoreExpiresAfter struct {
	// Anchor timestamp after which the expiration policy applies. Supported anchors:
	// `last_active_at`.
	Anchor VectorStoreExpiresAfterAnchor `json:"anchor,required"`
	// The number of days after the anchor time that the vector store will expire.
	Days int64                       `json:"days,required"`
	JSON vectorStoreExpiresAfterJSON `json:"-"`
}

// vectorStoreExpiresAfterJSON contains the JSON metadata for the struct
// [VectorStoreExpiresAfter]
type vectorStoreExpiresAfterJSON struct {
	Anchor      apijson.Field
	Days        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *VectorStoreExpiresAfter) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r vectorStoreExpiresAfterJSON) RawJSON() string {
	return r.raw
}

// Anchor timestamp after which the expiration policy applies. Supported anchors:
// `last_active_at`.
type VectorStoreExpiresAfterAnchor string

const (
	VectorStoreExpiresAfterAnchorLastActiveAt VectorStoreExpiresAfterAnchor = "last_active_at"
)

func (r VectorStoreExpiresAfterAnchor) IsKnown() bool {
	switch r {
	case VectorStoreExpiresAfterAnchorLastActiveAt:
		return true
	}
	return false
}

type VectorStoreDeleted struct {
	ID      string                   `json:"id,required"`
	Deleted bool                     `json:"deleted,required"`
	Object  VectorStoreDeletedObject `json:"object,required"`
	JSON    vectorStoreDeletedJSON   `json:"-"`
}

// vectorStoreDeletedJSON contains the JSON metadata for the struct
// [VectorStoreDeleted]
type vectorStoreDeletedJSON struct {
	ID          apijson.Field
	Deleted     apijson.Field
	Object      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *VectorStoreDeleted) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r vectorStoreDeletedJSON) RawJSON() string {
	return r.raw
}

type VectorStoreDeletedObject string

const (
	VectorStoreDeletedObjectVectorStoreDeleted VectorStoreDeletedObject = "vector_store.deleted"
)

func (r VectorStoreDeletedObject) IsKnown() bool {
	switch r {
	case VectorStoreDeletedObjectVectorStoreDeleted:
		return true
	}
	return false
}

type BetaVectorStoreNewParams struct {
	// The chunking strategy used to chunk the file(s). If not set, will use the `auto`
	// strategy. Only applicable if `file_ids` is non-empty.
	ChunkingStrategy param.Field[BetaVectorStoreNewParamsChunkingStrategyUnion] `json:"chunking_strategy"`
	// The expiration policy for a vector store.
	ExpiresAfter param.Field[BetaVectorStoreNewParamsExpiresAfter] `json:"expires_after"`
	// A list of [File](https://platform.openai.com/docs/api-reference/files) IDs that
	// the vector store should use. Useful for tools like `file_search` that can access
	// files.
	FileIDs param.Field[[]string] `json:"file_ids"`
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format. Keys
	// can be a maximum of 64 characters long and values can be a maxium of 512
	// characters long.
	Metadata param.Field[interface{}] `json:"metadata"`
	// The name of the vector store.
	Name param.Field[string] `json:"name"`
}

func (r BetaVectorStoreNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// The chunking strategy used to chunk the file(s). If not set, will use the `auto`
// strategy. Only applicable if `file_ids` is non-empty.
type BetaVectorStoreNewParamsChunkingStrategy struct {
	// Always `auto`.
	Type   param.Field[BetaVectorStoreNewParamsChunkingStrategyType] `json:"type,required"`
	Static param.Field[interface{}]                                  `json:"static,required"`
}

func (r BetaVectorStoreNewParamsChunkingStrategy) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r BetaVectorStoreNewParamsChunkingStrategy) implementsBetaVectorStoreNewParamsChunkingStrategyUnion() {
}

// The chunking strategy used to chunk the file(s). If not set, will use the `auto`
// strategy. Only applicable if `file_ids` is non-empty.
//
// Satisfied by [BetaVectorStoreNewParamsChunkingStrategyAuto],
// [BetaVectorStoreNewParamsChunkingStrategyStatic],
// [BetaVectorStoreNewParamsChunkingStrategy].
type BetaVectorStoreNewParamsChunkingStrategyUnion interface {
	implementsBetaVectorStoreNewParamsChunkingStrategyUnion()
}

// The default strategy. This strategy currently uses a `max_chunk_size_tokens` of
// `800` and `chunk_overlap_tokens` of `400`.
type BetaVectorStoreNewParamsChunkingStrategyAuto struct {
	// Always `auto`.
	Type param.Field[BetaVectorStoreNewParamsChunkingStrategyAutoType] `json:"type,required"`
}

func (r BetaVectorStoreNewParamsChunkingStrategyAuto) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r BetaVectorStoreNewParamsChunkingStrategyAuto) implementsBetaVectorStoreNewParamsChunkingStrategyUnion() {
}

// Always `auto`.
type BetaVectorStoreNewParamsChunkingStrategyAutoType string

const (
	BetaVectorStoreNewParamsChunkingStrategyAutoTypeAuto BetaVectorStoreNewParamsChunkingStrategyAutoType = "auto"
)

func (r BetaVectorStoreNewParamsChunkingStrategyAutoType) IsKnown() bool {
	switch r {
	case BetaVectorStoreNewParamsChunkingStrategyAutoTypeAuto:
		return true
	}
	return false
}

type BetaVectorStoreNewParamsChunkingStrategyStatic struct {
	Static param.Field[BetaVectorStoreNewParamsChunkingStrategyStaticStatic] `json:"static,required"`
	// Always `static`.
	Type param.Field[BetaVectorStoreNewParamsChunkingStrategyStaticType] `json:"type,required"`
}

func (r BetaVectorStoreNewParamsChunkingStrategyStatic) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r BetaVectorStoreNewParamsChunkingStrategyStatic) implementsBetaVectorStoreNewParamsChunkingStrategyUnion() {
}

type BetaVectorStoreNewParamsChunkingStrategyStaticStatic struct {
	// The number of tokens that overlap between chunks. The default value is `400`.
	//
	// Note that the overlap must not exceed half of `max_chunk_size_tokens`.
	ChunkOverlapTokens param.Field[int64] `json:"chunk_overlap_tokens,required"`
	// The maximum number of tokens in each chunk. The default value is `800`. The
	// minimum value is `100` and the maximum value is `4096`.
	MaxChunkSizeTokens param.Field[int64] `json:"max_chunk_size_tokens,required"`
}

func (r BetaVectorStoreNewParamsChunkingStrategyStaticStatic) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Always `static`.
type BetaVectorStoreNewParamsChunkingStrategyStaticType string

const (
	BetaVectorStoreNewParamsChunkingStrategyStaticTypeStatic BetaVectorStoreNewParamsChunkingStrategyStaticType = "static"
)

func (r BetaVectorStoreNewParamsChunkingStrategyStaticType) IsKnown() bool {
	switch r {
	case BetaVectorStoreNewParamsChunkingStrategyStaticTypeStatic:
		return true
	}
	return false
}

// Always `auto`.
type BetaVectorStoreNewParamsChunkingStrategyType string

const (
	BetaVectorStoreNewParamsChunkingStrategyTypeAuto   BetaVectorStoreNewParamsChunkingStrategyType = "auto"
	BetaVectorStoreNewParamsChunkingStrategyTypeStatic BetaVectorStoreNewParamsChunkingStrategyType = "static"
)

func (r BetaVectorStoreNewParamsChunkingStrategyType) IsKnown() bool {
	switch r {
	case BetaVectorStoreNewParamsChunkingStrategyTypeAuto, BetaVectorStoreNewParamsChunkingStrategyTypeStatic:
		return true
	}
	return false
}

// The expiration policy for a vector store.
type BetaVectorStoreNewParamsExpiresAfter struct {
	// Anchor timestamp after which the expiration policy applies. Supported anchors:
	// `last_active_at`.
	Anchor param.Field[BetaVectorStoreNewParamsExpiresAfterAnchor] `json:"anchor,required"`
	// The number of days after the anchor time that the vector store will expire.
	Days param.Field[int64] `json:"days,required"`
}

func (r BetaVectorStoreNewParamsExpiresAfter) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Anchor timestamp after which the expiration policy applies. Supported anchors:
// `last_active_at`.
type BetaVectorStoreNewParamsExpiresAfterAnchor string

const (
	BetaVectorStoreNewParamsExpiresAfterAnchorLastActiveAt BetaVectorStoreNewParamsExpiresAfterAnchor = "last_active_at"
)

func (r BetaVectorStoreNewParamsExpiresAfterAnchor) IsKnown() bool {
	switch r {
	case BetaVectorStoreNewParamsExpiresAfterAnchorLastActiveAt:
		return true
	}
	return false
}

type BetaVectorStoreUpdateParams struct {
	// The expiration policy for a vector store.
	ExpiresAfter param.Field[BetaVectorStoreUpdateParamsExpiresAfter] `json:"expires_after"`
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format. Keys
	// can be a maximum of 64 characters long and values can be a maxium of 512
	// characters long.
	Metadata param.Field[interface{}] `json:"metadata"`
	// The name of the vector store.
	Name param.Field[string] `json:"name"`
}

func (r BetaVectorStoreUpdateParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// The expiration policy for a vector store.
type BetaVectorStoreUpdateParamsExpiresAfter struct {
	// Anchor timestamp after which the expiration policy applies. Supported anchors:
	// `last_active_at`.
	Anchor param.Field[BetaVectorStoreUpdateParamsExpiresAfterAnchor] `json:"anchor,required"`
	// The number of days after the anchor time that the vector store will expire.
	Days param.Field[int64] `json:"days,required"`
}

func (r BetaVectorStoreUpdateParamsExpiresAfter) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Anchor timestamp after which the expiration policy applies. Supported anchors:
// `last_active_at`.
type BetaVectorStoreUpdateParamsExpiresAfterAnchor string

const (
	BetaVectorStoreUpdateParamsExpiresAfterAnchorLastActiveAt BetaVectorStoreUpdateParamsExpiresAfterAnchor = "last_active_at"
)

func (r BetaVectorStoreUpdateParamsExpiresAfterAnchor) IsKnown() bool {
	switch r {
	case BetaVectorStoreUpdateParamsExpiresAfterAnchorLastActiveAt:
		return true
	}
	return false
}

type BetaVectorStoreListParams struct {
	// A cursor for use in pagination. `after` is an object ID that defines your place
	// in the list. For instance, if you make a list request and receive 100 objects,
	// ending with obj_foo, your subsequent call can include after=obj_foo in order to
	// fetch the next page of the list.
	After param.Field[string] `query:"after"`
	// A cursor for use in pagination. `before` is an object ID that defines your place
	// in the list. For instance, if you make a list request and receive 100 objects,
	// ending with obj_foo, your subsequent call can include before=obj_foo in order to
	// fetch the previous page of the list.
	Before param.Field[string] `query:"before"`
	// A limit on the number of objects to be returned. Limit can range between 1 and
	// 100, and the default is 20.
	Limit param.Field[int64] `query:"limit"`
	// Sort order by the `created_at` timestamp of the objects. `asc` for ascending
	// order and `desc` for descending order.
	Order param.Field[BetaVectorStoreListParamsOrder] `query:"order"`
}

// URLQuery serializes [BetaVectorStoreListParams]'s query parameters as
// `url.Values`.
func (r BetaVectorStoreListParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Sort order by the `created_at` timestamp of the objects. `asc` for ascending
// order and `desc` for descending order.
type BetaVectorStoreListParamsOrder string

const (
	BetaVectorStoreListParamsOrderAsc  BetaVectorStoreListParamsOrder = "asc"
	BetaVectorStoreListParamsOrderDesc BetaVectorStoreListParamsOrder = "desc"
)

func (r BetaVectorStoreListParamsOrder) IsKnown() bool {
	switch r {
	case BetaVectorStoreListParamsOrderAsc, BetaVectorStoreListParamsOrderDesc:
		return true
	}
	return false
}

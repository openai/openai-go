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
	"github.com/openai/openai-go/internal/pagination"
	"github.com/openai/openai-go/internal/param"
	"github.com/openai/openai-go/internal/requestconfig"
	"github.com/openai/openai-go/option"
	"github.com/tidwall/gjson"
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
	if vectorStoreID == "" {
		err = errors.New("missing required vector_store_id parameter")
		return
	}
	path := fmt.Sprintf("vector_stores/%s/files", vectorStoreID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Retrieves a vector store file.
func (r *BetaVectorStoreFileService) Get(ctx context.Context, vectorStoreID string, fileID string, opts ...option.RequestOption) (res *VectorStoreFile, err error) {
	opts = append(r.Options[:], opts...)
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
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
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
	ChunkingStrategy VectorStoreFileChunkingStrategy `json:"chunking_strategy"`
	JSON             vectorStoreFileJSON             `json:"-"`
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

// The strategy used to chunk the file.
type VectorStoreFileChunkingStrategy struct {
	// Always `static`.
	Type VectorStoreFileChunkingStrategyType `json:"type,required"`
	// This field can have the runtime type of
	// [VectorStoreFileChunkingStrategyStaticStatic].
	Static interface{}                         `json:"static,required"`
	JSON   vectorStoreFileChunkingStrategyJSON `json:"-"`
	union  VectorStoreFileChunkingStrategyUnion
}

// vectorStoreFileChunkingStrategyJSON contains the JSON metadata for the struct
// [VectorStoreFileChunkingStrategy]
type vectorStoreFileChunkingStrategyJSON struct {
	Type        apijson.Field
	Static      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r vectorStoreFileChunkingStrategyJSON) RawJSON() string {
	return r.raw
}

func (r *VectorStoreFileChunkingStrategy) UnmarshalJSON(data []byte) (err error) {
	*r = VectorStoreFileChunkingStrategy{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [VectorStoreFileChunkingStrategyUnion] interface which you can
// cast to the specific types for more type safety.
//
// Possible runtime types of the union are [VectorStoreFileChunkingStrategyStatic],
// [VectorStoreFileChunkingStrategyOther].
func (r VectorStoreFileChunkingStrategy) AsUnion() VectorStoreFileChunkingStrategyUnion {
	return r.union
}

// The strategy used to chunk the file.
//
// Union satisfied by [VectorStoreFileChunkingStrategyStatic] or
// [VectorStoreFileChunkingStrategyOther].
type VectorStoreFileChunkingStrategyUnion interface {
	implementsVectorStoreFileChunkingStrategy()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*VectorStoreFileChunkingStrategyUnion)(nil)).Elem(),
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(VectorStoreFileChunkingStrategyStatic{}),
			DiscriminatorValue: "static",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(VectorStoreFileChunkingStrategyOther{}),
			DiscriminatorValue: "other",
		},
	)
}

type VectorStoreFileChunkingStrategyStatic struct {
	Static VectorStoreFileChunkingStrategyStaticStatic `json:"static,required"`
	// Always `static`.
	Type VectorStoreFileChunkingStrategyStaticType `json:"type,required"`
	JSON vectorStoreFileChunkingStrategyStaticJSON `json:"-"`
}

// vectorStoreFileChunkingStrategyStaticJSON contains the JSON metadata for the
// struct [VectorStoreFileChunkingStrategyStatic]
type vectorStoreFileChunkingStrategyStaticJSON struct {
	Static      apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *VectorStoreFileChunkingStrategyStatic) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r vectorStoreFileChunkingStrategyStaticJSON) RawJSON() string {
	return r.raw
}

func (r VectorStoreFileChunkingStrategyStatic) implementsVectorStoreFileChunkingStrategy() {}

type VectorStoreFileChunkingStrategyStaticStatic struct {
	// The number of tokens that overlap between chunks. The default value is `400`.
	//
	// Note that the overlap must not exceed half of `max_chunk_size_tokens`.
	ChunkOverlapTokens int64 `json:"chunk_overlap_tokens,required"`
	// The maximum number of tokens in each chunk. The default value is `800`. The
	// minimum value is `100` and the maximum value is `4096`.
	MaxChunkSizeTokens int64                                           `json:"max_chunk_size_tokens,required"`
	JSON               vectorStoreFileChunkingStrategyStaticStaticJSON `json:"-"`
}

// vectorStoreFileChunkingStrategyStaticStaticJSON contains the JSON metadata for
// the struct [VectorStoreFileChunkingStrategyStaticStatic]
type vectorStoreFileChunkingStrategyStaticStaticJSON struct {
	ChunkOverlapTokens apijson.Field
	MaxChunkSizeTokens apijson.Field
	raw                string
	ExtraFields        map[string]apijson.Field
}

func (r *VectorStoreFileChunkingStrategyStaticStatic) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r vectorStoreFileChunkingStrategyStaticStaticJSON) RawJSON() string {
	return r.raw
}

// Always `static`.
type VectorStoreFileChunkingStrategyStaticType string

const (
	VectorStoreFileChunkingStrategyStaticTypeStatic VectorStoreFileChunkingStrategyStaticType = "static"
)

func (r VectorStoreFileChunkingStrategyStaticType) IsKnown() bool {
	switch r {
	case VectorStoreFileChunkingStrategyStaticTypeStatic:
		return true
	}
	return false
}

// This is returned when the chunking strategy is unknown. Typically, this is
// because the file was indexed before the `chunking_strategy` concept was
// introduced in the API.
type VectorStoreFileChunkingStrategyOther struct {
	// Always `other`.
	Type VectorStoreFileChunkingStrategyOtherType `json:"type,required"`
	JSON vectorStoreFileChunkingStrategyOtherJSON `json:"-"`
}

// vectorStoreFileChunkingStrategyOtherJSON contains the JSON metadata for the
// struct [VectorStoreFileChunkingStrategyOther]
type vectorStoreFileChunkingStrategyOtherJSON struct {
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *VectorStoreFileChunkingStrategyOther) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r vectorStoreFileChunkingStrategyOtherJSON) RawJSON() string {
	return r.raw
}

func (r VectorStoreFileChunkingStrategyOther) implementsVectorStoreFileChunkingStrategy() {}

// Always `other`.
type VectorStoreFileChunkingStrategyOtherType string

const (
	VectorStoreFileChunkingStrategyOtherTypeOther VectorStoreFileChunkingStrategyOtherType = "other"
)

func (r VectorStoreFileChunkingStrategyOtherType) IsKnown() bool {
	switch r {
	case VectorStoreFileChunkingStrategyOtherTypeOther:
		return true
	}
	return false
}

// Always `static`.
type VectorStoreFileChunkingStrategyType string

const (
	VectorStoreFileChunkingStrategyTypeStatic VectorStoreFileChunkingStrategyType = "static"
	VectorStoreFileChunkingStrategyTypeOther  VectorStoreFileChunkingStrategyType = "other"
)

func (r VectorStoreFileChunkingStrategyType) IsKnown() bool {
	switch r {
	case VectorStoreFileChunkingStrategyTypeStatic, VectorStoreFileChunkingStrategyTypeOther:
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
	// strategy.
	ChunkingStrategy param.Field[BetaVectorStoreFileNewParamsChunkingStrategyUnion] `json:"chunking_strategy"`
}

func (r BetaVectorStoreFileNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// The chunking strategy used to chunk the file(s). If not set, will use the `auto`
// strategy.
type BetaVectorStoreFileNewParamsChunkingStrategy struct {
	// Always `auto`.
	Type   param.Field[BetaVectorStoreFileNewParamsChunkingStrategyType] `json:"type,required"`
	Static param.Field[interface{}]                                      `json:"static,required"`
}

func (r BetaVectorStoreFileNewParamsChunkingStrategy) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r BetaVectorStoreFileNewParamsChunkingStrategy) implementsBetaVectorStoreFileNewParamsChunkingStrategyUnion() {
}

// The chunking strategy used to chunk the file(s). If not set, will use the `auto`
// strategy.
//
// Satisfied by
// [BetaVectorStoreFileNewParamsChunkingStrategyAutoChunkingStrategyRequestParam],
// [BetaVectorStoreFileNewParamsChunkingStrategyStaticChunkingStrategyRequestParam],
// [BetaVectorStoreFileNewParamsChunkingStrategy].
type BetaVectorStoreFileNewParamsChunkingStrategyUnion interface {
	implementsBetaVectorStoreFileNewParamsChunkingStrategyUnion()
}

// The default strategy. This strategy currently uses a `max_chunk_size_tokens` of
// `800` and `chunk_overlap_tokens` of `400`.
type BetaVectorStoreFileNewParamsChunkingStrategyAutoChunkingStrategyRequestParam struct {
	// Always `auto`.
	Type param.Field[BetaVectorStoreFileNewParamsChunkingStrategyAutoChunkingStrategyRequestParamType] `json:"type,required"`
}

func (r BetaVectorStoreFileNewParamsChunkingStrategyAutoChunkingStrategyRequestParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r BetaVectorStoreFileNewParamsChunkingStrategyAutoChunkingStrategyRequestParam) implementsBetaVectorStoreFileNewParamsChunkingStrategyUnion() {
}

// Always `auto`.
type BetaVectorStoreFileNewParamsChunkingStrategyAutoChunkingStrategyRequestParamType string

const (
	BetaVectorStoreFileNewParamsChunkingStrategyAutoChunkingStrategyRequestParamTypeAuto BetaVectorStoreFileNewParamsChunkingStrategyAutoChunkingStrategyRequestParamType = "auto"
)

func (r BetaVectorStoreFileNewParamsChunkingStrategyAutoChunkingStrategyRequestParamType) IsKnown() bool {
	switch r {
	case BetaVectorStoreFileNewParamsChunkingStrategyAutoChunkingStrategyRequestParamTypeAuto:
		return true
	}
	return false
}

type BetaVectorStoreFileNewParamsChunkingStrategyStaticChunkingStrategyRequestParam struct {
	Static param.Field[BetaVectorStoreFileNewParamsChunkingStrategyStaticChunkingStrategyRequestParamStatic] `json:"static,required"`
	// Always `static`.
	Type param.Field[BetaVectorStoreFileNewParamsChunkingStrategyStaticChunkingStrategyRequestParamType] `json:"type,required"`
}

func (r BetaVectorStoreFileNewParamsChunkingStrategyStaticChunkingStrategyRequestParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r BetaVectorStoreFileNewParamsChunkingStrategyStaticChunkingStrategyRequestParam) implementsBetaVectorStoreFileNewParamsChunkingStrategyUnion() {
}

type BetaVectorStoreFileNewParamsChunkingStrategyStaticChunkingStrategyRequestParamStatic struct {
	// The number of tokens that overlap between chunks. The default value is `400`.
	//
	// Note that the overlap must not exceed half of `max_chunk_size_tokens`.
	ChunkOverlapTokens param.Field[int64] `json:"chunk_overlap_tokens,required"`
	// The maximum number of tokens in each chunk. The default value is `800`. The
	// minimum value is `100` and the maximum value is `4096`.
	MaxChunkSizeTokens param.Field[int64] `json:"max_chunk_size_tokens,required"`
}

func (r BetaVectorStoreFileNewParamsChunkingStrategyStaticChunkingStrategyRequestParamStatic) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Always `static`.
type BetaVectorStoreFileNewParamsChunkingStrategyStaticChunkingStrategyRequestParamType string

const (
	BetaVectorStoreFileNewParamsChunkingStrategyStaticChunkingStrategyRequestParamTypeStatic BetaVectorStoreFileNewParamsChunkingStrategyStaticChunkingStrategyRequestParamType = "static"
)

func (r BetaVectorStoreFileNewParamsChunkingStrategyStaticChunkingStrategyRequestParamType) IsKnown() bool {
	switch r {
	case BetaVectorStoreFileNewParamsChunkingStrategyStaticChunkingStrategyRequestParamTypeStatic:
		return true
	}
	return false
}

// Always `auto`.
type BetaVectorStoreFileNewParamsChunkingStrategyType string

const (
	BetaVectorStoreFileNewParamsChunkingStrategyTypeAuto   BetaVectorStoreFileNewParamsChunkingStrategyType = "auto"
	BetaVectorStoreFileNewParamsChunkingStrategyTypeStatic BetaVectorStoreFileNewParamsChunkingStrategyType = "static"
)

func (r BetaVectorStoreFileNewParamsChunkingStrategyType) IsKnown() bool {
	switch r {
	case BetaVectorStoreFileNewParamsChunkingStrategyTypeAuto, BetaVectorStoreFileNewParamsChunkingStrategyTypeStatic:
		return true
	}
	return false
}

type BetaVectorStoreFileListParams struct {
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

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

// BetaVectorStoreService contains methods and other services that help with
// interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaVectorStoreService] method instead.
type BetaVectorStoreService struct {
	Options     []option.RequestOption
	Files       BetaVectorStoreFileService
	FileBatches BetaVectorStoreFileBatchService
}

// NewBetaVectorStoreService generates a new service that applies the given options
// to each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewBetaVectorStoreService(opts ...option.RequestOption) (r BetaVectorStoreService) {
	r = BetaVectorStoreService{}
	r.Options = opts
	r.Files = NewBetaVectorStoreFileService(opts...)
	r.FileBatches = NewBetaVectorStoreFileBatchService(opts...)
	return
}

// Create a vector store.
func (r *BetaVectorStoreService) New(ctx context.Context, body BetaVectorStoreNewParams, opts ...option.RequestOption) (res *VectorStore, err error) {
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "assistants=v2")}, opts...)
	path := "vector_stores"
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
		FileID: String(fileObj.ID),
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

// Retrieves a vector store.
func (r *BetaVectorStoreService) Get(ctx context.Context, vectorStoreID string, opts ...option.RequestOption) (res *VectorStore, err error) {
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "assistants=v2")}, opts...)
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
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "assistants=v2")}, opts...)
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
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "assistants=v2"), option.WithResponseInto(&raw)}, opts...)
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
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "assistants=v2")}, opts...)
	if vectorStoreID == "" {
		err = errors.New("missing required vector_store_id parameter")
		return
	}
	path := fmt.Sprintf("vector_stores/%s", vectorStoreID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return
}

// The default strategy. This strategy currently uses a `max_chunk_size_tokens` of
// `800` and `chunk_overlap_tokens` of `400`.
type AutoFileChunkingStrategyParam struct {
	// Always `auto`.
	//
	// This field can be elided, and will be automatically set as "auto".
	Type constant.Auto `json:"type,required"`
	apiobject
}

func (f AutoFileChunkingStrategyParam) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

func (r AutoFileChunkingStrategyParam) MarshalJSON() (data []byte, err error) {
	type shadow AutoFileChunkingStrategyParam
	return param.MarshalObject(r, (*shadow)(&r))
}

type FileChunkingStrategyUnion struct {
	Static StaticFileChunkingStrategy `json:"static"`
	Type   string                     `json:"type"`
	JSON   struct {
		Static resp.Field
		Type   resp.Field
		raw    string
	} `json:"-"`
}

// note: this function is generated only for discriminated unions
func (u FileChunkingStrategyUnion) Variant() (res struct {
	OfStatic *StaticFileChunkingStrategyObject
	OfOther  *OtherFileChunkingStrategyObject
}) {
	switch u.Type {
	case "static":
		v := u.AsStatic()
		res.OfStatic = &v
	case "other":
		v := u.AsOther()
		res.OfOther = &v
	}
	return
}

func (u FileChunkingStrategyUnion) WhichKind() string {
	return u.Type
}

func (u FileChunkingStrategyUnion) AsStatic() (v StaticFileChunkingStrategyObject) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u FileChunkingStrategyUnion) AsOther() (v OtherFileChunkingStrategyObject) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u FileChunkingStrategyUnion) RawJSON() string { return u.JSON.raw }

func (r *FileChunkingStrategyUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewFileChunkingStrategyParamOfStatic(static StaticFileChunkingStrategyParam) FileChunkingStrategyParamUnion {
	var variant StaticFileChunkingStrategyObjectParam
	variant.Static = static
	return FileChunkingStrategyParamUnion{OfStatic: &variant}
}

// Only one field can be non-zero
type FileChunkingStrategyParamUnion struct {
	OfAuto   *AutoFileChunkingStrategyParam
	OfStatic *StaticFileChunkingStrategyObjectParam
	apiunion
}

func (u FileChunkingStrategyParamUnion) IsMissing() bool { return param.IsOmitted(u) || u.IsNull() }

func (u FileChunkingStrategyParamUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion[FileChunkingStrategyParamUnion](u.OfAuto, u.OfStatic)
}

func (u FileChunkingStrategyParamUnion) GetStatic() *StaticFileChunkingStrategyParam {
	if vt := u.OfStatic; vt != nil {
		return &vt.Static
	}
	return nil
}

func (u FileChunkingStrategyParamUnion) GetType() *string {
	if vt := u.OfAuto; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfStatic; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// This is returned when the chunking strategy is unknown. Typically, this is
// because the file was indexed before the `chunking_strategy` concept was
// introduced in the API.
type OtherFileChunkingStrategyObject struct {
	// Always `other`.
	//
	// This field can be elided, and will be automatically set as "other".
	Type constant.Other `json:"type,required"`
	JSON struct {
		Type resp.Field
		raw  string
	} `json:"-"`
}

func (r OtherFileChunkingStrategyObject) RawJSON() string { return r.JSON.raw }
func (r *OtherFileChunkingStrategyObject) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type StaticFileChunkingStrategy struct {
	// The number of tokens that overlap between chunks. The default value is `400`.
	//
	// Note that the overlap must not exceed half of `max_chunk_size_tokens`.
	ChunkOverlapTokens int64 `json:"chunk_overlap_tokens,omitzero,required"`
	// The maximum number of tokens in each chunk. The default value is `800`. The
	// minimum value is `100` and the maximum value is `4096`.
	MaxChunkSizeTokens int64 `json:"max_chunk_size_tokens,omitzero,required"`
	JSON               struct {
		ChunkOverlapTokens resp.Field
		MaxChunkSizeTokens resp.Field
		raw                string
	} `json:"-"`
}

func (r StaticFileChunkingStrategy) RawJSON() string { return r.JSON.raw }
func (r *StaticFileChunkingStrategy) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this StaticFileChunkingStrategy to a
// StaticFileChunkingStrategyParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// StaticFileChunkingStrategyParam.IsOverridden()
func (r StaticFileChunkingStrategy) ToParam() StaticFileChunkingStrategyParam {
	return param.Override[StaticFileChunkingStrategyParam](r.RawJSON())
}

type StaticFileChunkingStrategyParam struct {
	// The number of tokens that overlap between chunks. The default value is `400`.
	//
	// Note that the overlap must not exceed half of `max_chunk_size_tokens`.
	ChunkOverlapTokens param.Int `json:"chunk_overlap_tokens,omitzero,required"`
	// The maximum number of tokens in each chunk. The default value is `800`. The
	// minimum value is `100` and the maximum value is `4096`.
	MaxChunkSizeTokens param.Int `json:"max_chunk_size_tokens,omitzero,required"`
	apiobject
}

func (f StaticFileChunkingStrategyParam) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

func (r StaticFileChunkingStrategyParam) MarshalJSON() (data []byte, err error) {
	type shadow StaticFileChunkingStrategyParam
	return param.MarshalObject(r, (*shadow)(&r))
}

type StaticFileChunkingStrategyObject struct {
	Static StaticFileChunkingStrategy `json:"static,omitzero,required"`
	// Always `static`.
	//
	// This field can be elided, and will be automatically set as "static".
	Type constant.Static `json:"type,required"`
	JSON struct {
		Static resp.Field
		Type   resp.Field
		raw    string
	} `json:"-"`
}

func (r StaticFileChunkingStrategyObject) RawJSON() string { return r.JSON.raw }
func (r *StaticFileChunkingStrategyObject) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type StaticFileChunkingStrategyObjectParam struct {
	Static StaticFileChunkingStrategyParam `json:"static,omitzero,required"`
	// Always `static`.
	//
	// This field can be elided, and will be automatically set as "static".
	Type constant.Static `json:"type,required"`
	apiobject
}

func (f StaticFileChunkingStrategyObjectParam) IsMissing() bool {
	return param.IsOmitted(f) || f.IsNull()
}

func (r StaticFileChunkingStrategyObjectParam) MarshalJSON() (data []byte, err error) {
	type shadow StaticFileChunkingStrategyObjectParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// A vector store is a collection of processed files can be used by the
// `file_search` tool.
type VectorStore struct {
	// The identifier, which can be referenced in API endpoints.
	ID string `json:"id,omitzero,required"`
	// The Unix timestamp (in seconds) for when the vector store was created.
	CreatedAt  int64                 `json:"created_at,omitzero,required"`
	FileCounts VectorStoreFileCounts `json:"file_counts,omitzero,required"`
	// The Unix timestamp (in seconds) for when the vector store was last active.
	LastActiveAt int64 `json:"last_active_at,omitzero,required,nullable"`
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard.
	//
	// Keys are strings with a maximum length of 64 characters. Values are strings with
	// a maximum length of 512 characters.
	Metadata shared.Metadata `json:"metadata,omitzero,required,nullable"`
	// The name of the vector store.
	Name string `json:"name,omitzero,required"`
	// The object type, which is always `vector_store`.
	//
	// This field can be elided, and will be automatically set as "vector_store".
	Object constant.VectorStore `json:"object,required"`
	// The status of the vector store, which can be either `expired`, `in_progress`, or
	// `completed`. A status of `completed` indicates that the vector store is ready
	// for use.
	//
	// Any of "expired", "in_progress", "completed"
	Status string `json:"status,omitzero,required"`
	// The total number of bytes used by the files in the vector store.
	UsageBytes int64 `json:"usage_bytes,omitzero,required"`
	// The expiration policy for a vector store.
	ExpiresAfter VectorStoreExpiresAfter `json:"expires_after,omitzero"`
	// The Unix timestamp (in seconds) for when the vector store will expire.
	ExpiresAt int64 `json:"expires_at,omitzero,nullable"`
	JSON      struct {
		ID           resp.Field
		CreatedAt    resp.Field
		FileCounts   resp.Field
		LastActiveAt resp.Field
		Metadata     resp.Field
		Name         resp.Field
		Object       resp.Field
		Status       resp.Field
		UsageBytes   resp.Field
		ExpiresAfter resp.Field
		ExpiresAt    resp.Field
		raw          string
	} `json:"-"`
}

func (r VectorStore) RawJSON() string { return r.JSON.raw }
func (r *VectorStore) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type VectorStoreFileCounts struct {
	// The number of files that were cancelled.
	Cancelled int64 `json:"cancelled,omitzero,required"`
	// The number of files that have been successfully processed.
	Completed int64 `json:"completed,omitzero,required"`
	// The number of files that have failed to process.
	Failed int64 `json:"failed,omitzero,required"`
	// The number of files that are currently being processed.
	InProgress int64 `json:"in_progress,omitzero,required"`
	// The total number of files.
	Total int64 `json:"total,omitzero,required"`
	JSON  struct {
		Cancelled  resp.Field
		Completed  resp.Field
		Failed     resp.Field
		InProgress resp.Field
		Total      resp.Field
		raw        string
	} `json:"-"`
}

func (r VectorStoreFileCounts) RawJSON() string { return r.JSON.raw }
func (r *VectorStoreFileCounts) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The status of the vector store, which can be either `expired`, `in_progress`, or
// `completed`. A status of `completed` indicates that the vector store is ready
// for use.
type VectorStoreStatus = string

const (
	VectorStoreStatusExpired    VectorStoreStatus = "expired"
	VectorStoreStatusInProgress VectorStoreStatus = "in_progress"
	VectorStoreStatusCompleted  VectorStoreStatus = "completed"
)

// The expiration policy for a vector store.
type VectorStoreExpiresAfter struct {
	// Anchor timestamp after which the expiration policy applies. Supported anchors:
	// `last_active_at`.
	//
	// This field can be elided, and will be automatically set as "last_active_at".
	Anchor constant.LastActiveAt `json:"anchor,required"`
	// The number of days after the anchor time that the vector store will expire.
	Days int64 `json:"days,omitzero,required"`
	JSON struct {
		Anchor resp.Field
		Days   resp.Field
		raw    string
	} `json:"-"`
}

func (r VectorStoreExpiresAfter) RawJSON() string { return r.JSON.raw }
func (r *VectorStoreExpiresAfter) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type VectorStoreDeleted struct {
	ID      string `json:"id,omitzero,required"`
	Deleted bool   `json:"deleted,omitzero,required"`
	// This field can be elided, and will be automatically set as
	// "vector_store.deleted".
	Object constant.VectorStoreDeleted `json:"object,required"`
	JSON   struct {
		ID      resp.Field
		Deleted resp.Field
		Object  resp.Field
		raw     string
	} `json:"-"`
}

func (r VectorStoreDeleted) RawJSON() string { return r.JSON.raw }
func (r *VectorStoreDeleted) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaVectorStoreNewParams struct {
	// The chunking strategy used to chunk the file(s). If not set, will use the `auto`
	// strategy. Only applicable if `file_ids` is non-empty.
	ChunkingStrategy FileChunkingStrategyParamUnion `json:"chunking_strategy,omitzero"`
	// The expiration policy for a vector store.
	ExpiresAfter BetaVectorStoreNewParamsExpiresAfter `json:"expires_after,omitzero"`
	// A list of [File](https://platform.openai.com/docs/api-reference/files) IDs that
	// the vector store should use. Useful for tools like `file_search` that can access
	// files.
	FileIDs []string `json:"file_ids,omitzero"`
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard.
	//
	// Keys are strings with a maximum length of 64 characters. Values are strings with
	// a maximum length of 512 characters.
	Metadata shared.MetadataParam `json:"metadata,omitzero"`
	// The name of the vector store.
	Name param.String `json:"name,omitzero"`
	apiobject
}

func (f BetaVectorStoreNewParams) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

func (r BetaVectorStoreNewParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaVectorStoreNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}

// The expiration policy for a vector store.
type BetaVectorStoreNewParamsExpiresAfter struct {
	// Anchor timestamp after which the expiration policy applies. Supported anchors:
	// `last_active_at`.
	//
	// This field can be elided, and will be automatically set as "last_active_at".
	Anchor constant.LastActiveAt `json:"anchor,required"`
	// The number of days after the anchor time that the vector store will expire.
	Days param.Int `json:"days,omitzero,required"`
	apiobject
}

func (f BetaVectorStoreNewParamsExpiresAfter) IsMissing() bool {
	return param.IsOmitted(f) || f.IsNull()
}

func (r BetaVectorStoreNewParamsExpiresAfter) MarshalJSON() (data []byte, err error) {
	type shadow BetaVectorStoreNewParamsExpiresAfter
	return param.MarshalObject(r, (*shadow)(&r))
}

type BetaVectorStoreUpdateParams struct {
	// The expiration policy for a vector store.
	ExpiresAfter BetaVectorStoreUpdateParamsExpiresAfter `json:"expires_after,omitzero"`
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard.
	//
	// Keys are strings with a maximum length of 64 characters. Values are strings with
	// a maximum length of 512 characters.
	Metadata shared.MetadataParam `json:"metadata,omitzero"`
	// The name of the vector store.
	Name param.String `json:"name,omitzero"`
	apiobject
}

func (f BetaVectorStoreUpdateParams) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

func (r BetaVectorStoreUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaVectorStoreUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}

// The expiration policy for a vector store.
type BetaVectorStoreUpdateParamsExpiresAfter struct {
	// Anchor timestamp after which the expiration policy applies. Supported anchors:
	// `last_active_at`.
	//
	// This field can be elided, and will be automatically set as "last_active_at".
	Anchor constant.LastActiveAt `json:"anchor,required"`
	// The number of days after the anchor time that the vector store will expire.
	Days param.Int `json:"days,omitzero,required"`
	apiobject
}

func (f BetaVectorStoreUpdateParamsExpiresAfter) IsMissing() bool {
	return param.IsOmitted(f) || f.IsNull()
}

func (r BetaVectorStoreUpdateParamsExpiresAfter) MarshalJSON() (data []byte, err error) {
	type shadow BetaVectorStoreUpdateParamsExpiresAfter
	return param.MarshalObject(r, (*shadow)(&r))
}

type BetaVectorStoreListParams struct {
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
	Order BetaVectorStoreListParamsOrder `query:"order,omitzero"`
	apiobject
}

func (f BetaVectorStoreListParams) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

// URLQuery serializes [BetaVectorStoreListParams]'s query parameters as
// `url.Values`.
func (r BetaVectorStoreListParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
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

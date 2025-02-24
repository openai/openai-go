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
	"github.com/openai/openai-go/packages/resp"
	"github.com/openai/openai-go/shared/constant"
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
func NewBetaVectorStoreFileService(opts ...option.RequestOption) (r BetaVectorStoreFileService) {
	r = BetaVectorStoreFileService{}
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
	ID string `json:"id,omitzero,required"`
	// The Unix timestamp (in seconds) for when the vector store file was created.
	CreatedAt int64 `json:"created_at,omitzero,required"`
	// The last error associated with this vector store file. Will be `null` if there
	// are no errors.
	LastError VectorStoreFileLastError `json:"last_error,omitzero,required,nullable"`
	// The object type, which is always `vector_store.file`.
	//
	// This field can be elided, and will be automatically set as "vector_store.file".
	Object constant.VectorStoreFile `json:"object,required"`
	// The status of the vector store file, which can be either `in_progress`,
	// `completed`, `cancelled`, or `failed`. The status `completed` indicates that the
	// vector store file is ready for use.
	//
	// Any of "in_progress", "completed", "cancelled", "failed"
	Status string `json:"status,omitzero,required"`
	// The total vector store usage in bytes. Note that this may be different from the
	// original file size.
	UsageBytes int64 `json:"usage_bytes,omitzero,required"`
	// The ID of the
	// [vector store](https://platform.openai.com/docs/api-reference/vector-stores/object)
	// that the [File](https://platform.openai.com/docs/api-reference/files) is
	// attached to.
	VectorStoreID string `json:"vector_store_id,omitzero,required"`
	// The strategy used to chunk the file.
	ChunkingStrategy FileChunkingStrategyUnion `json:"chunking_strategy,omitzero"`
	JSON             struct {
		ID               resp.Field
		CreatedAt        resp.Field
		LastError        resp.Field
		Object           resp.Field
		Status           resp.Field
		UsageBytes       resp.Field
		VectorStoreID    resp.Field
		ChunkingStrategy resp.Field
		raw              string
	} `json:"-"`
}

func (r VectorStoreFile) RawJSON() string { return r.JSON.raw }
func (r *VectorStoreFile) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The last error associated with this vector store file. Will be `null` if there
// are no errors.
type VectorStoreFileLastError struct {
	// One of `server_error` or `rate_limit_exceeded`.
	//
	// Any of "server_error", "unsupported_file", "invalid_file"
	Code string `json:"code,omitzero,required"`
	// A human-readable description of the error.
	Message string `json:"message,omitzero,required"`
	JSON    struct {
		Code    resp.Field
		Message resp.Field
		raw     string
	} `json:"-"`
}

func (r VectorStoreFileLastError) RawJSON() string { return r.JSON.raw }
func (r *VectorStoreFileLastError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// One of `server_error` or `rate_limit_exceeded`.
type VectorStoreFileLastErrorCode = string

const (
	VectorStoreFileLastErrorCodeServerError     VectorStoreFileLastErrorCode = "server_error"
	VectorStoreFileLastErrorCodeUnsupportedFile VectorStoreFileLastErrorCode = "unsupported_file"
	VectorStoreFileLastErrorCodeInvalidFile     VectorStoreFileLastErrorCode = "invalid_file"
)

// The status of the vector store file, which can be either `in_progress`,
// `completed`, `cancelled`, or `failed`. The status `completed` indicates that the
// vector store file is ready for use.
type VectorStoreFileStatus = string

const (
	VectorStoreFileStatusInProgress VectorStoreFileStatus = "in_progress"
	VectorStoreFileStatusCompleted  VectorStoreFileStatus = "completed"
	VectorStoreFileStatusCancelled  VectorStoreFileStatus = "cancelled"
	VectorStoreFileStatusFailed     VectorStoreFileStatus = "failed"
)

type VectorStoreFileDeleted struct {
	ID      string `json:"id,omitzero,required"`
	Deleted bool   `json:"deleted,omitzero,required"`
	// This field can be elided, and will be automatically set as
	// "vector_store.file.deleted".
	Object constant.VectorStoreFileDeleted `json:"object,required"`
	JSON   struct {
		ID      resp.Field
		Deleted resp.Field
		Object  resp.Field
		raw     string
	} `json:"-"`
}

func (r VectorStoreFileDeleted) RawJSON() string { return r.JSON.raw }
func (r *VectorStoreFileDeleted) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaVectorStoreFileNewParams struct {
	// A [File](https://platform.openai.com/docs/api-reference/files) ID that the
	// vector store should use. Useful for tools like `file_search` that can access
	// files.
	FileID param.String `json:"file_id,omitzero,required"`
	// The chunking strategy used to chunk the file(s). If not set, will use the `auto`
	// strategy. Only applicable if `file_ids` is non-empty.
	ChunkingStrategy FileChunkingStrategyParamUnion `json:"chunking_strategy,omitzero"`
	apiobject
}

func (f BetaVectorStoreFileNewParams) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

func (r BetaVectorStoreFileNewParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaVectorStoreFileNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}

type BetaVectorStoreFileListParams struct {
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
	// Filter by file status. One of `in_progress`, `completed`, `failed`, `cancelled`.
	//
	// Any of "in_progress", "completed", "failed", "cancelled"
	Filter BetaVectorStoreFileListParamsFilter `query:"filter,omitzero"`
	// A limit on the number of objects to be returned. Limit can range between 1 and
	// 100, and the default is 20.
	Limit param.Int `query:"limit,omitzero"`
	// Sort order by the `created_at` timestamp of the objects. `asc` for ascending
	// order and `desc` for descending order.
	//
	// Any of "asc", "desc"
	Order BetaVectorStoreFileListParamsOrder `query:"order,omitzero"`
	apiobject
}

func (f BetaVectorStoreFileListParams) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

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

// Sort order by the `created_at` timestamp of the objects. `asc` for ascending
// order and `desc` for descending order.
type BetaVectorStoreFileListParamsOrder string

const (
	BetaVectorStoreFileListParamsOrderAsc  BetaVectorStoreFileListParamsOrder = "asc"
	BetaVectorStoreFileListParamsOrderDesc BetaVectorStoreFileListParamsOrder = "desc"
)

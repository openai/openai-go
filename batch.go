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
	"github.com/openai/openai-go/shared"
	"github.com/openai/openai-go/shared/constant"
)

// BatchService contains methods and other services that help with interacting with
// the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBatchService] method instead.
type BatchService struct {
	Options []option.RequestOption
}

// NewBatchService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewBatchService(opts ...option.RequestOption) (r BatchService) {
	r = BatchService{}
	r.Options = opts
	return
}

// Creates and executes a batch from an uploaded file of requests
func (r *BatchService) New(ctx context.Context, body BatchNewParams, opts ...option.RequestOption) (res *Batch, err error) {
	opts = append(r.Options[:], opts...)
	path := "batches"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Retrieves a batch.
func (r *BatchService) Get(ctx context.Context, batchID string, opts ...option.RequestOption) (res *Batch, err error) {
	opts = append(r.Options[:], opts...)
	if batchID == "" {
		err = errors.New("missing required batch_id parameter")
		return
	}
	path := fmt.Sprintf("batches/%s", batchID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// List your organization's batches.
func (r *BatchService) List(ctx context.Context, query BatchListParams, opts ...option.RequestOption) (res *pagination.CursorPage[Batch], err error) {
	var raw *http.Response
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "batches"
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

// List your organization's batches.
func (r *BatchService) ListAutoPaging(ctx context.Context, query BatchListParams, opts ...option.RequestOption) *pagination.CursorPageAutoPager[Batch] {
	return pagination.NewCursorPageAutoPager(r.List(ctx, query, opts...))
}

// Cancels an in-progress batch. The batch will be in status `cancelling` for up to
// 10 minutes, before changing to `cancelled`, where it will have partial results
// (if any) available in the output file.
func (r *BatchService) Cancel(ctx context.Context, batchID string, opts ...option.RequestOption) (res *Batch, err error) {
	opts = append(r.Options[:], opts...)
	if batchID == "" {
		err = errors.New("missing required batch_id parameter")
		return
	}
	path := fmt.Sprintf("batches/%s/cancel", batchID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return
}

type Batch struct {
	ID string `json:"id,omitzero,required"`
	// The time frame within which the batch should be processed.
	CompletionWindow string `json:"completion_window,omitzero,required"`
	// The Unix timestamp (in seconds) for when the batch was created.
	CreatedAt int64 `json:"created_at,omitzero,required"`
	// The OpenAI API endpoint used by the batch.
	Endpoint string `json:"endpoint,omitzero,required"`
	// The ID of the input file for the batch.
	InputFileID string `json:"input_file_id,omitzero,required"`
	// The object type, which is always `batch`.
	//
	// This field can be elided, and will be automatically set as "batch".
	Object constant.Batch `json:"object,required"`
	// The current status of the batch.
	//
	// Any of "validating", "failed", "in_progress", "finalizing", "completed",
	// "expired", "cancelling", "cancelled"
	Status string `json:"status,omitzero,required"`
	// The Unix timestamp (in seconds) for when the batch was cancelled.
	CancelledAt int64 `json:"cancelled_at,omitzero"`
	// The Unix timestamp (in seconds) for when the batch started cancelling.
	CancellingAt int64 `json:"cancelling_at,omitzero"`
	// The Unix timestamp (in seconds) for when the batch was completed.
	CompletedAt int64 `json:"completed_at,omitzero"`
	// The ID of the file containing the outputs of requests with errors.
	ErrorFileID string      `json:"error_file_id,omitzero"`
	Errors      BatchErrors `json:"errors,omitzero"`
	// The Unix timestamp (in seconds) for when the batch expired.
	ExpiredAt int64 `json:"expired_at,omitzero"`
	// The Unix timestamp (in seconds) for when the batch will expire.
	ExpiresAt int64 `json:"expires_at,omitzero"`
	// The Unix timestamp (in seconds) for when the batch failed.
	FailedAt int64 `json:"failed_at,omitzero"`
	// The Unix timestamp (in seconds) for when the batch started finalizing.
	FinalizingAt int64 `json:"finalizing_at,omitzero"`
	// The Unix timestamp (in seconds) for when the batch started processing.
	InProgressAt int64 `json:"in_progress_at,omitzero"`
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard.
	//
	// Keys are strings with a maximum length of 64 characters. Values are strings with
	// a maximum length of 512 characters.
	Metadata shared.Metadata `json:"metadata,omitzero,nullable"`
	// The ID of the file containing the outputs of successfully executed requests.
	OutputFileID string `json:"output_file_id,omitzero"`
	// The request counts for different statuses within the batch.
	RequestCounts BatchRequestCounts `json:"request_counts,omitzero"`
	JSON          struct {
		ID               resp.Field
		CompletionWindow resp.Field
		CreatedAt        resp.Field
		Endpoint         resp.Field
		InputFileID      resp.Field
		Object           resp.Field
		Status           resp.Field
		CancelledAt      resp.Field
		CancellingAt     resp.Field
		CompletedAt      resp.Field
		ErrorFileID      resp.Field
		Errors           resp.Field
		ExpiredAt        resp.Field
		ExpiresAt        resp.Field
		FailedAt         resp.Field
		FinalizingAt     resp.Field
		InProgressAt     resp.Field
		Metadata         resp.Field
		OutputFileID     resp.Field
		RequestCounts    resp.Field
		raw              string
	} `json:"-"`
}

func (r Batch) RawJSON() string { return r.JSON.raw }
func (r *Batch) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The current status of the batch.
type BatchStatus = string

const (
	BatchStatusValidating BatchStatus = "validating"
	BatchStatusFailed     BatchStatus = "failed"
	BatchStatusInProgress BatchStatus = "in_progress"
	BatchStatusFinalizing BatchStatus = "finalizing"
	BatchStatusCompleted  BatchStatus = "completed"
	BatchStatusExpired    BatchStatus = "expired"
	BatchStatusCancelling BatchStatus = "cancelling"
	BatchStatusCancelled  BatchStatus = "cancelled"
)

type BatchErrors struct {
	Data []BatchError `json:"data,omitzero"`
	// The object type, which is always `list`.
	Object string `json:"object,omitzero"`
	JSON   struct {
		Data   resp.Field
		Object resp.Field
		raw    string
	} `json:"-"`
}

func (r BatchErrors) RawJSON() string { return r.JSON.raw }
func (r *BatchErrors) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BatchError struct {
	// An error code identifying the error type.
	Code string `json:"code,omitzero"`
	// The line number of the input file where the error occurred, if applicable.
	Line int64 `json:"line,omitzero,nullable"`
	// A human-readable message providing more details about the error.
	Message string `json:"message,omitzero"`
	// The name of the parameter that caused the error, if applicable.
	Param string `json:"param,omitzero,nullable"`
	JSON  struct {
		Code    resp.Field
		Line    resp.Field
		Message resp.Field
		Param   resp.Field
		raw     string
	} `json:"-"`
}

func (r BatchError) RawJSON() string { return r.JSON.raw }
func (r *BatchError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The request counts for different statuses within the batch.
type BatchRequestCounts struct {
	// Number of requests that have been completed successfully.
	Completed int64 `json:"completed,omitzero,required"`
	// Number of requests that have failed.
	Failed int64 `json:"failed,omitzero,required"`
	// Total number of requests in the batch.
	Total int64 `json:"total,omitzero,required"`
	JSON  struct {
		Completed resp.Field
		Failed    resp.Field
		Total     resp.Field
		raw       string
	} `json:"-"`
}

func (r BatchRequestCounts) RawJSON() string { return r.JSON.raw }
func (r *BatchRequestCounts) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BatchNewParams struct {
	// The time frame within which the batch should be processed. Currently only `24h`
	// is supported.
	//
	// Any of "24h"
	CompletionWindow BatchNewParamsCompletionWindow `json:"completion_window,required"`
	// The endpoint to be used for all requests in the batch. Currently
	// `/v1/chat/completions`, `/v1/embeddings`, and `/v1/completions` are supported.
	// Note that `/v1/embeddings` batches are also restricted to a maximum of 50,000
	// embedding inputs across all requests in the batch.
	//
	// Any of "/v1/chat/completions", "/v1/embeddings", "/v1/completions"
	Endpoint BatchNewParamsEndpoint `json:"endpoint,omitzero,required"`
	// The ID of an uploaded file that contains requests for the new batch.
	//
	// See [upload file](https://platform.openai.com/docs/api-reference/files/create)
	// for how to upload a file.
	//
	// Your input file must be formatted as a
	// [JSONL file](https://platform.openai.com/docs/api-reference/batch/request-input),
	// and must be uploaded with the purpose `batch`. The file can contain up to 50,000
	// requests, and can be up to 200 MB in size.
	InputFileID param.String `json:"input_file_id,omitzero,required"`
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard.
	//
	// Keys are strings with a maximum length of 64 characters. Values are strings with
	// a maximum length of 512 characters.
	Metadata shared.MetadataParam `json:"metadata,omitzero"`
	apiobject
}

func (f BatchNewParams) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

func (r BatchNewParams) MarshalJSON() (data []byte, err error) {
	type shadow BatchNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}

// The time frame within which the batch should be processed. Currently only `24h`
// is supported.
type BatchNewParamsCompletionWindow string

const (
	BatchNewParamsCompletionWindow24h BatchNewParamsCompletionWindow = "24h"
)

// The endpoint to be used for all requests in the batch. Currently
// `/v1/chat/completions`, `/v1/embeddings`, and `/v1/completions` are supported.
// Note that `/v1/embeddings` batches are also restricted to a maximum of 50,000
// embedding inputs across all requests in the batch.
type BatchNewParamsEndpoint string

const (
	BatchNewParamsEndpointV1ChatCompletions BatchNewParamsEndpoint = "/v1/chat/completions"
	BatchNewParamsEndpointV1Embeddings      BatchNewParamsEndpoint = "/v1/embeddings"
	BatchNewParamsEndpointV1Completions     BatchNewParamsEndpoint = "/v1/completions"
)

type BatchListParams struct {
	// A cursor for use in pagination. `after` is an object ID that defines your place
	// in the list. For instance, if you make a list request and receive 100 objects,
	// ending with obj_foo, your subsequent call can include after=obj_foo in order to
	// fetch the next page of the list.
	After param.String `query:"after,omitzero"`
	// A limit on the number of objects to be returned. Limit can range between 1 and
	// 100, and the default is 20.
	Limit param.Int `query:"limit,omitzero"`
	apiobject
}

func (f BatchListParams) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

// URLQuery serializes [BatchListParams]'s query parameters as `url.Values`.
func (r BatchListParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

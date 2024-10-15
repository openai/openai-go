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

// FineTuningJobCheckpointService contains methods and other services that help
// with interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewFineTuningJobCheckpointService] method instead.
type FineTuningJobCheckpointService struct {
	Options []option.RequestOption
}

// NewFineTuningJobCheckpointService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewFineTuningJobCheckpointService(opts ...option.RequestOption) (r *FineTuningJobCheckpointService) {
	r = &FineTuningJobCheckpointService{}
	r.Options = opts
	return
}

// List checkpoints for a fine-tuning job.
func (r *FineTuningJobCheckpointService) List(ctx context.Context, fineTuningJobID string, query FineTuningJobCheckpointListParams, opts ...option.RequestOption) (res *pagination.CursorPage[FineTuningJobCheckpoint], err error) {
	var raw *http.Response
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	if fineTuningJobID == "" {
		err = errors.New("missing required fine_tuning_job_id parameter")
		return
	}
	path := fmt.Sprintf("fine_tuning/jobs/%s/checkpoints", fineTuningJobID)
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

// List checkpoints for a fine-tuning job.
func (r *FineTuningJobCheckpointService) ListAutoPaging(ctx context.Context, fineTuningJobID string, query FineTuningJobCheckpointListParams, opts ...option.RequestOption) *pagination.CursorPageAutoPager[FineTuningJobCheckpoint] {
	return pagination.NewCursorPageAutoPager(r.List(ctx, fineTuningJobID, query, opts...))
}

// The `fine_tuning.job.checkpoint` object represents a model checkpoint for a
// fine-tuning job that is ready to use.
type FineTuningJobCheckpoint struct {
	// The checkpoint identifier, which can be referenced in the API endpoints.
	ID string `json:"id,required"`
	// The Unix timestamp (in seconds) for when the checkpoint was created.
	CreatedAt int64 `json:"created_at,required"`
	// The name of the fine-tuned checkpoint model that is created.
	FineTunedModelCheckpoint string `json:"fine_tuned_model_checkpoint,required"`
	// The name of the fine-tuning job that this checkpoint was created from.
	FineTuningJobID string `json:"fine_tuning_job_id,required"`
	// Metrics at the step number during the fine-tuning job.
	Metrics FineTuningJobCheckpointMetrics `json:"metrics,required"`
	// The object type, which is always "fine_tuning.job.checkpoint".
	Object FineTuningJobCheckpointObject `json:"object,required"`
	// The step number that the checkpoint was created at.
	StepNumber int64                       `json:"step_number,required"`
	JSON       fineTuningJobCheckpointJSON `json:"-"`
}

// fineTuningJobCheckpointJSON contains the JSON metadata for the struct
// [FineTuningJobCheckpoint]
type fineTuningJobCheckpointJSON struct {
	ID                       apijson.Field
	CreatedAt                apijson.Field
	FineTunedModelCheckpoint apijson.Field
	FineTuningJobID          apijson.Field
	Metrics                  apijson.Field
	Object                   apijson.Field
	StepNumber               apijson.Field
	raw                      string
	ExtraFields              map[string]apijson.Field
}

func (r *FineTuningJobCheckpoint) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r fineTuningJobCheckpointJSON) RawJSON() string {
	return r.raw
}

// Metrics at the step number during the fine-tuning job.
type FineTuningJobCheckpointMetrics struct {
	FullValidLoss              float64                            `json:"full_valid_loss"`
	FullValidMeanTokenAccuracy float64                            `json:"full_valid_mean_token_accuracy"`
	Step                       float64                            `json:"step"`
	TrainLoss                  float64                            `json:"train_loss"`
	TrainMeanTokenAccuracy     float64                            `json:"train_mean_token_accuracy"`
	ValidLoss                  float64                            `json:"valid_loss"`
	ValidMeanTokenAccuracy     float64                            `json:"valid_mean_token_accuracy"`
	JSON                       fineTuningJobCheckpointMetricsJSON `json:"-"`
}

// fineTuningJobCheckpointMetricsJSON contains the JSON metadata for the struct
// [FineTuningJobCheckpointMetrics]
type fineTuningJobCheckpointMetricsJSON struct {
	FullValidLoss              apijson.Field
	FullValidMeanTokenAccuracy apijson.Field
	Step                       apijson.Field
	TrainLoss                  apijson.Field
	TrainMeanTokenAccuracy     apijson.Field
	ValidLoss                  apijson.Field
	ValidMeanTokenAccuracy     apijson.Field
	raw                        string
	ExtraFields                map[string]apijson.Field
}

func (r *FineTuningJobCheckpointMetrics) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r fineTuningJobCheckpointMetricsJSON) RawJSON() string {
	return r.raw
}

// The object type, which is always "fine_tuning.job.checkpoint".
type FineTuningJobCheckpointObject string

const (
	FineTuningJobCheckpointObjectFineTuningJobCheckpoint FineTuningJobCheckpointObject = "fine_tuning.job.checkpoint"
)

func (r FineTuningJobCheckpointObject) IsKnown() bool {
	switch r {
	case FineTuningJobCheckpointObjectFineTuningJobCheckpoint:
		return true
	}
	return false
}

type FineTuningJobCheckpointListParams struct {
	// Identifier for the last checkpoint ID from the previous pagination request.
	After param.Field[string] `query:"after"`
	// Number of checkpoints to retrieve.
	Limit param.Field[int64] `query:"limit"`
}

// URLQuery serializes [FineTuningJobCheckpointListParams]'s query parameters as
// `url.Values`.
func (r FineTuningJobCheckpointListParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

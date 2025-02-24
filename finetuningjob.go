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
	"github.com/openai/openai-go/shared/constant"
)

// FineTuningJobService contains methods and other services that help with
// interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewFineTuningJobService] method instead.
type FineTuningJobService struct {
	Options     []option.RequestOption
	Checkpoints FineTuningJobCheckpointService
}

// NewFineTuningJobService generates a new service that applies the given options
// to each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewFineTuningJobService(opts ...option.RequestOption) (r FineTuningJobService) {
	r = FineTuningJobService{}
	r.Options = opts
	r.Checkpoints = NewFineTuningJobCheckpointService(opts...)
	return
}

// Creates a fine-tuning job which begins the process of creating a new model from
// a given dataset.
//
// Response includes details of the enqueued job including job status and the name
// of the fine-tuned models once complete.
//
// [Learn more about fine-tuning](https://platform.openai.com/docs/guides/fine-tuning)
func (r *FineTuningJobService) New(ctx context.Context, body FineTuningJobNewParams, opts ...option.RequestOption) (res *FineTuningJob, err error) {
	opts = append(r.Options[:], opts...)
	path := "fine_tuning/jobs"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Get info about a fine-tuning job.
//
// [Learn more about fine-tuning](https://platform.openai.com/docs/guides/fine-tuning)
func (r *FineTuningJobService) Get(ctx context.Context, fineTuningJobID string, opts ...option.RequestOption) (res *FineTuningJob, err error) {
	opts = append(r.Options[:], opts...)
	if fineTuningJobID == "" {
		err = errors.New("missing required fine_tuning_job_id parameter")
		return
	}
	path := fmt.Sprintf("fine_tuning/jobs/%s", fineTuningJobID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// List your organization's fine-tuning jobs
func (r *FineTuningJobService) List(ctx context.Context, query FineTuningJobListParams, opts ...option.RequestOption) (res *pagination.CursorPage[FineTuningJob], err error) {
	var raw *http.Response
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "fine_tuning/jobs"
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

// List your organization's fine-tuning jobs
func (r *FineTuningJobService) ListAutoPaging(ctx context.Context, query FineTuningJobListParams, opts ...option.RequestOption) *pagination.CursorPageAutoPager[FineTuningJob] {
	return pagination.NewCursorPageAutoPager(r.List(ctx, query, opts...))
}

// Immediately cancel a fine-tune job.
func (r *FineTuningJobService) Cancel(ctx context.Context, fineTuningJobID string, opts ...option.RequestOption) (res *FineTuningJob, err error) {
	opts = append(r.Options[:], opts...)
	if fineTuningJobID == "" {
		err = errors.New("missing required fine_tuning_job_id parameter")
		return
	}
	path := fmt.Sprintf("fine_tuning/jobs/%s/cancel", fineTuningJobID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return
}

// Get status updates for a fine-tuning job.
func (r *FineTuningJobService) ListEvents(ctx context.Context, fineTuningJobID string, query FineTuningJobListEventsParams, opts ...option.RequestOption) (res *pagination.CursorPage[FineTuningJobEvent], err error) {
	var raw *http.Response
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	if fineTuningJobID == "" {
		err = errors.New("missing required fine_tuning_job_id parameter")
		return
	}
	path := fmt.Sprintf("fine_tuning/jobs/%s/events", fineTuningJobID)
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

// Get status updates for a fine-tuning job.
func (r *FineTuningJobService) ListEventsAutoPaging(ctx context.Context, fineTuningJobID string, query FineTuningJobListEventsParams, opts ...option.RequestOption) *pagination.CursorPageAutoPager[FineTuningJobEvent] {
	return pagination.NewCursorPageAutoPager(r.ListEvents(ctx, fineTuningJobID, query, opts...))
}

// The `fine_tuning.job` object represents a fine-tuning job that has been created
// through the API.
type FineTuningJob struct {
	// The object identifier, which can be referenced in the API endpoints.
	ID string `json:"id,omitzero,required"`
	// The Unix timestamp (in seconds) for when the fine-tuning job was created.
	CreatedAt int64 `json:"created_at,omitzero,required"`
	// For fine-tuning jobs that have `failed`, this will contain more information on
	// the cause of the failure.
	Error FineTuningJobError `json:"error,omitzero,required,nullable"`
	// The name of the fine-tuned model that is being created. The value will be null
	// if the fine-tuning job is still running.
	FineTunedModel string `json:"fine_tuned_model,omitzero,required,nullable"`
	// The Unix timestamp (in seconds) for when the fine-tuning job was finished. The
	// value will be null if the fine-tuning job is still running.
	FinishedAt int64 `json:"finished_at,omitzero,required,nullable"`
	// The hyperparameters used for the fine-tuning job. This value will only be
	// returned when running `supervised` jobs.
	Hyperparameters FineTuningJobHyperparameters `json:"hyperparameters,omitzero,required"`
	// The base model that is being fine-tuned.
	Model string `json:"model,omitzero,required"`
	// The object type, which is always "fine_tuning.job".
	//
	// This field can be elided, and will be automatically set as "fine_tuning.job".
	Object constant.FineTuningJob `json:"object,required"`
	// The organization that owns the fine-tuning job.
	OrganizationID string `json:"organization_id,omitzero,required"`
	// The compiled results file ID(s) for the fine-tuning job. You can retrieve the
	// results with the
	// [Files API](https://platform.openai.com/docs/api-reference/files/retrieve-contents).
	ResultFiles []string `json:"result_files,omitzero,required"`
	// The seed used for the fine-tuning job.
	Seed int64 `json:"seed,omitzero,required"`
	// The current status of the fine-tuning job, which can be either
	// `validating_files`, `queued`, `running`, `succeeded`, `failed`, or `cancelled`.
	//
	// Any of "validating_files", "queued", "running", "succeeded", "failed",
	// "cancelled"
	Status string `json:"status,omitzero,required"`
	// The total number of billable tokens processed by this fine-tuning job. The value
	// will be null if the fine-tuning job is still running.
	TrainedTokens int64 `json:"trained_tokens,omitzero,required,nullable"`
	// The file ID used for training. You can retrieve the training data with the
	// [Files API](https://platform.openai.com/docs/api-reference/files/retrieve-contents).
	TrainingFile string `json:"training_file,omitzero,required"`
	// The file ID used for validation. You can retrieve the validation results with
	// the
	// [Files API](https://platform.openai.com/docs/api-reference/files/retrieve-contents).
	ValidationFile string `json:"validation_file,omitzero,required,nullable"`
	// The Unix timestamp (in seconds) for when the fine-tuning job is estimated to
	// finish. The value will be null if the fine-tuning job is not running.
	EstimatedFinish int64 `json:"estimated_finish,omitzero,nullable"`
	// A list of integrations to enable for this fine-tuning job.
	Integrations []FineTuningJobWandbIntegrationObject `json:"integrations,omitzero,nullable"`
	// The method used for fine-tuning.
	Method FineTuningJobMethod `json:"method,omitzero"`
	JSON   struct {
		ID              resp.Field
		CreatedAt       resp.Field
		Error           resp.Field
		FineTunedModel  resp.Field
		FinishedAt      resp.Field
		Hyperparameters resp.Field
		Model           resp.Field
		Object          resp.Field
		OrganizationID  resp.Field
		ResultFiles     resp.Field
		Seed            resp.Field
		Status          resp.Field
		TrainedTokens   resp.Field
		TrainingFile    resp.Field
		ValidationFile  resp.Field
		EstimatedFinish resp.Field
		Integrations    resp.Field
		Method          resp.Field
		raw             string
	} `json:"-"`
}

func (r FineTuningJob) RawJSON() string { return r.JSON.raw }
func (r *FineTuningJob) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// For fine-tuning jobs that have `failed`, this will contain more information on
// the cause of the failure.
type FineTuningJobError struct {
	// A machine-readable error code.
	Code string `json:"code,omitzero,required"`
	// A human-readable error message.
	Message string `json:"message,omitzero,required"`
	// The parameter that was invalid, usually `training_file` or `validation_file`.
	// This field will be null if the failure was not parameter-specific.
	Param string `json:"param,omitzero,required,nullable"`
	JSON  struct {
		Code    resp.Field
		Message resp.Field
		Param   resp.Field
		raw     string
	} `json:"-"`
}

func (r FineTuningJobError) RawJSON() string { return r.JSON.raw }
func (r *FineTuningJobError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The hyperparameters used for the fine-tuning job. This value will only be
// returned when running `supervised` jobs.
type FineTuningJobHyperparameters struct {
	// Number of examples in each batch. A larger batch size means that model
	// parameters are updated less frequently, but with lower variance.
	BatchSize FineTuningJobHyperparametersBatchSizeUnion `json:"batch_size,omitzero"`
	// Scaling factor for the learning rate. A smaller learning rate may be useful to
	// avoid overfitting.
	LearningRateMultiplier FineTuningJobHyperparametersLearningRateMultiplierUnion `json:"learning_rate_multiplier,omitzero"`
	// The number of epochs to train the model for. An epoch refers to one full cycle
	// through the training dataset.
	NEpochs FineTuningJobHyperparametersNEpochsUnion `json:"n_epochs,omitzero"`
	JSON    struct {
		BatchSize              resp.Field
		LearningRateMultiplier resp.Field
		NEpochs                resp.Field
		raw                    string
	} `json:"-"`
}

func (r FineTuningJobHyperparameters) RawJSON() string { return r.JSON.raw }
func (r *FineTuningJobHyperparameters) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type FineTuningJobHyperparametersBatchSizeUnion struct {
	OfAuto  constant.Auto `json:",inline"`
	OfInt64 int64         `json:",inline"`
	JSON    struct {
		OfAuto  resp.Field
		OfInt64 resp.Field
		raw     string
	} `json:"-"`
}

func (u FineTuningJobHyperparametersBatchSizeUnion) AsAuto() (v constant.Auto) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u FineTuningJobHyperparametersBatchSizeUnion) AsInt() (v int64) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u FineTuningJobHyperparametersBatchSizeUnion) RawJSON() string { return u.JSON.raw }

func (r *FineTuningJobHyperparametersBatchSizeUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type FineTuningJobHyperparametersLearningRateMultiplierUnion struct {
	OfAuto    constant.Auto `json:",inline"`
	OfFloat64 float64       `json:",inline"`
	JSON      struct {
		OfAuto    resp.Field
		OfFloat64 resp.Field
		raw       string
	} `json:"-"`
}

func (u FineTuningJobHyperparametersLearningRateMultiplierUnion) AsAuto() (v constant.Auto) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u FineTuningJobHyperparametersLearningRateMultiplierUnion) AsFloat() (v float64) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u FineTuningJobHyperparametersLearningRateMultiplierUnion) RawJSON() string { return u.JSON.raw }

func (r *FineTuningJobHyperparametersLearningRateMultiplierUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type FineTuningJobHyperparametersNEpochsUnion struct {
	OfAuto  constant.Auto `json:",inline"`
	OfInt64 int64         `json:",inline"`
	JSON    struct {
		OfAuto  resp.Field
		OfInt64 resp.Field
		raw     string
	} `json:"-"`
}

func (u FineTuningJobHyperparametersNEpochsUnion) AsAuto() (v constant.Auto) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u FineTuningJobHyperparametersNEpochsUnion) AsInt() (v int64) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u FineTuningJobHyperparametersNEpochsUnion) RawJSON() string { return u.JSON.raw }

func (r *FineTuningJobHyperparametersNEpochsUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The current status of the fine-tuning job, which can be either
// `validating_files`, `queued`, `running`, `succeeded`, `failed`, or `cancelled`.
type FineTuningJobStatus = string

const (
	FineTuningJobStatusValidatingFiles FineTuningJobStatus = "validating_files"
	FineTuningJobStatusQueued          FineTuningJobStatus = "queued"
	FineTuningJobStatusRunning         FineTuningJobStatus = "running"
	FineTuningJobStatusSucceeded       FineTuningJobStatus = "succeeded"
	FineTuningJobStatusFailed          FineTuningJobStatus = "failed"
	FineTuningJobStatusCancelled       FineTuningJobStatus = "cancelled"
)

// The method used for fine-tuning.
type FineTuningJobMethod struct {
	// Configuration for the DPO fine-tuning method.
	Dpo FineTuningJobMethodDpo `json:"dpo,omitzero"`
	// Configuration for the supervised fine-tuning method.
	Supervised FineTuningJobMethodSupervised `json:"supervised,omitzero"`
	// The type of method. Is either `supervised` or `dpo`.
	//
	// Any of "supervised", "dpo"
	Type string `json:"type,omitzero"`
	JSON struct {
		Dpo        resp.Field
		Supervised resp.Field
		Type       resp.Field
		raw        string
	} `json:"-"`
}

func (r FineTuningJobMethod) RawJSON() string { return r.JSON.raw }
func (r *FineTuningJobMethod) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for the DPO fine-tuning method.
type FineTuningJobMethodDpo struct {
	// The hyperparameters used for the fine-tuning job.
	Hyperparameters FineTuningJobMethodDpoHyperparameters `json:"hyperparameters,omitzero"`
	JSON            struct {
		Hyperparameters resp.Field
		raw             string
	} `json:"-"`
}

func (r FineTuningJobMethodDpo) RawJSON() string { return r.JSON.raw }
func (r *FineTuningJobMethodDpo) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The hyperparameters used for the fine-tuning job.
type FineTuningJobMethodDpoHyperparameters struct {
	// Number of examples in each batch. A larger batch size means that model
	// parameters are updated less frequently, but with lower variance.
	BatchSize FineTuningJobMethodDpoHyperparametersBatchSizeUnion `json:"batch_size,omitzero"`
	// The beta value for the DPO method. A higher beta value will increase the weight
	// of the penalty between the policy and reference model.
	Beta FineTuningJobMethodDpoHyperparametersBetaUnion `json:"beta,omitzero"`
	// Scaling factor for the learning rate. A smaller learning rate may be useful to
	// avoid overfitting.
	LearningRateMultiplier FineTuningJobMethodDpoHyperparametersLearningRateMultiplierUnion `json:"learning_rate_multiplier,omitzero"`
	// The number of epochs to train the model for. An epoch refers to one full cycle
	// through the training dataset.
	NEpochs FineTuningJobMethodDpoHyperparametersNEpochsUnion `json:"n_epochs,omitzero"`
	JSON    struct {
		BatchSize              resp.Field
		Beta                   resp.Field
		LearningRateMultiplier resp.Field
		NEpochs                resp.Field
		raw                    string
	} `json:"-"`
}

func (r FineTuningJobMethodDpoHyperparameters) RawJSON() string { return r.JSON.raw }
func (r *FineTuningJobMethodDpoHyperparameters) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type FineTuningJobMethodDpoHyperparametersBatchSizeUnion struct {
	OfAuto  constant.Auto `json:",inline"`
	OfInt64 int64         `json:",inline"`
	JSON    struct {
		OfAuto  resp.Field
		OfInt64 resp.Field
		raw     string
	} `json:"-"`
}

func (u FineTuningJobMethodDpoHyperparametersBatchSizeUnion) AsAuto() (v constant.Auto) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u FineTuningJobMethodDpoHyperparametersBatchSizeUnion) AsInt() (v int64) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u FineTuningJobMethodDpoHyperparametersBatchSizeUnion) RawJSON() string { return u.JSON.raw }

func (r *FineTuningJobMethodDpoHyperparametersBatchSizeUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type FineTuningJobMethodDpoHyperparametersBetaUnion struct {
	OfAuto    constant.Auto `json:",inline"`
	OfFloat64 float64       `json:",inline"`
	JSON      struct {
		OfAuto    resp.Field
		OfFloat64 resp.Field
		raw       string
	} `json:"-"`
}

func (u FineTuningJobMethodDpoHyperparametersBetaUnion) AsAuto() (v constant.Auto) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u FineTuningJobMethodDpoHyperparametersBetaUnion) AsFloat() (v float64) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u FineTuningJobMethodDpoHyperparametersBetaUnion) RawJSON() string { return u.JSON.raw }

func (r *FineTuningJobMethodDpoHyperparametersBetaUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type FineTuningJobMethodDpoHyperparametersLearningRateMultiplierUnion struct {
	OfAuto    constant.Auto `json:",inline"`
	OfFloat64 float64       `json:",inline"`
	JSON      struct {
		OfAuto    resp.Field
		OfFloat64 resp.Field
		raw       string
	} `json:"-"`
}

func (u FineTuningJobMethodDpoHyperparametersLearningRateMultiplierUnion) AsAuto() (v constant.Auto) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u FineTuningJobMethodDpoHyperparametersLearningRateMultiplierUnion) AsFloat() (v float64) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u FineTuningJobMethodDpoHyperparametersLearningRateMultiplierUnion) RawJSON() string {
	return u.JSON.raw
}

func (r *FineTuningJobMethodDpoHyperparametersLearningRateMultiplierUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type FineTuningJobMethodDpoHyperparametersNEpochsUnion struct {
	OfAuto  constant.Auto `json:",inline"`
	OfInt64 int64         `json:",inline"`
	JSON    struct {
		OfAuto  resp.Field
		OfInt64 resp.Field
		raw     string
	} `json:"-"`
}

func (u FineTuningJobMethodDpoHyperparametersNEpochsUnion) AsAuto() (v constant.Auto) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u FineTuningJobMethodDpoHyperparametersNEpochsUnion) AsInt() (v int64) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u FineTuningJobMethodDpoHyperparametersNEpochsUnion) RawJSON() string { return u.JSON.raw }

func (r *FineTuningJobMethodDpoHyperparametersNEpochsUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for the supervised fine-tuning method.
type FineTuningJobMethodSupervised struct {
	// The hyperparameters used for the fine-tuning job.
	Hyperparameters FineTuningJobMethodSupervisedHyperparameters `json:"hyperparameters,omitzero"`
	JSON            struct {
		Hyperparameters resp.Field
		raw             string
	} `json:"-"`
}

func (r FineTuningJobMethodSupervised) RawJSON() string { return r.JSON.raw }
func (r *FineTuningJobMethodSupervised) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The hyperparameters used for the fine-tuning job.
type FineTuningJobMethodSupervisedHyperparameters struct {
	// Number of examples in each batch. A larger batch size means that model
	// parameters are updated less frequently, but with lower variance.
	BatchSize FineTuningJobMethodSupervisedHyperparametersBatchSizeUnion `json:"batch_size,omitzero"`
	// Scaling factor for the learning rate. A smaller learning rate may be useful to
	// avoid overfitting.
	LearningRateMultiplier FineTuningJobMethodSupervisedHyperparametersLearningRateMultiplierUnion `json:"learning_rate_multiplier,omitzero"`
	// The number of epochs to train the model for. An epoch refers to one full cycle
	// through the training dataset.
	NEpochs FineTuningJobMethodSupervisedHyperparametersNEpochsUnion `json:"n_epochs,omitzero"`
	JSON    struct {
		BatchSize              resp.Field
		LearningRateMultiplier resp.Field
		NEpochs                resp.Field
		raw                    string
	} `json:"-"`
}

func (r FineTuningJobMethodSupervisedHyperparameters) RawJSON() string { return r.JSON.raw }
func (r *FineTuningJobMethodSupervisedHyperparameters) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type FineTuningJobMethodSupervisedHyperparametersBatchSizeUnion struct {
	OfAuto  constant.Auto `json:",inline"`
	OfInt64 int64         `json:",inline"`
	JSON    struct {
		OfAuto  resp.Field
		OfInt64 resp.Field
		raw     string
	} `json:"-"`
}

func (u FineTuningJobMethodSupervisedHyperparametersBatchSizeUnion) AsAuto() (v constant.Auto) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u FineTuningJobMethodSupervisedHyperparametersBatchSizeUnion) AsInt() (v int64) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u FineTuningJobMethodSupervisedHyperparametersBatchSizeUnion) RawJSON() string {
	return u.JSON.raw
}

func (r *FineTuningJobMethodSupervisedHyperparametersBatchSizeUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type FineTuningJobMethodSupervisedHyperparametersLearningRateMultiplierUnion struct {
	OfAuto    constant.Auto `json:",inline"`
	OfFloat64 float64       `json:",inline"`
	JSON      struct {
		OfAuto    resp.Field
		OfFloat64 resp.Field
		raw       string
	} `json:"-"`
}

func (u FineTuningJobMethodSupervisedHyperparametersLearningRateMultiplierUnion) AsAuto() (v constant.Auto) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u FineTuningJobMethodSupervisedHyperparametersLearningRateMultiplierUnion) AsFloat() (v float64) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u FineTuningJobMethodSupervisedHyperparametersLearningRateMultiplierUnion) RawJSON() string {
	return u.JSON.raw
}

func (r *FineTuningJobMethodSupervisedHyperparametersLearningRateMultiplierUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type FineTuningJobMethodSupervisedHyperparametersNEpochsUnion struct {
	OfAuto  constant.Auto `json:",inline"`
	OfInt64 int64         `json:",inline"`
	JSON    struct {
		OfAuto  resp.Field
		OfInt64 resp.Field
		raw     string
	} `json:"-"`
}

func (u FineTuningJobMethodSupervisedHyperparametersNEpochsUnion) AsAuto() (v constant.Auto) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u FineTuningJobMethodSupervisedHyperparametersNEpochsUnion) AsInt() (v int64) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u FineTuningJobMethodSupervisedHyperparametersNEpochsUnion) RawJSON() string { return u.JSON.raw }

func (r *FineTuningJobMethodSupervisedHyperparametersNEpochsUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The type of method. Is either `supervised` or `dpo`.
type FineTuningJobMethodType = string

const (
	FineTuningJobMethodTypeSupervised FineTuningJobMethodType = "supervised"
	FineTuningJobMethodTypeDpo        FineTuningJobMethodType = "dpo"
)

// Fine-tuning job event object
type FineTuningJobEvent struct {
	// The object identifier.
	ID string `json:"id,omitzero,required"`
	// The Unix timestamp (in seconds) for when the fine-tuning job was created.
	CreatedAt int64 `json:"created_at,omitzero,required"`
	// The log level of the event.
	//
	// Any of "info", "warn", "error"
	Level string `json:"level,omitzero,required"`
	// The message of the event.
	Message string `json:"message,omitzero,required"`
	// The object type, which is always "fine_tuning.job.event".
	//
	// This field can be elided, and will be automatically set as
	// "fine_tuning.job.event".
	Object constant.FineTuningJobEvent `json:"object,required"`
	// The data associated with the event.
	Data interface{} `json:"data,omitzero"`
	// The type of event.
	//
	// Any of "message", "metrics"
	Type string `json:"type,omitzero"`
	JSON struct {
		ID        resp.Field
		CreatedAt resp.Field
		Level     resp.Field
		Message   resp.Field
		Object    resp.Field
		Data      resp.Field
		Type      resp.Field
		raw       string
	} `json:"-"`
}

func (r FineTuningJobEvent) RawJSON() string { return r.JSON.raw }
func (r *FineTuningJobEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The log level of the event.
type FineTuningJobEventLevel = string

const (
	FineTuningJobEventLevelInfo  FineTuningJobEventLevel = "info"
	FineTuningJobEventLevelWarn  FineTuningJobEventLevel = "warn"
	FineTuningJobEventLevelError FineTuningJobEventLevel = "error"
)

// The type of event.
type FineTuningJobEventType = string

const (
	FineTuningJobEventTypeMessage FineTuningJobEventType = "message"
	FineTuningJobEventTypeMetrics FineTuningJobEventType = "metrics"
)

type FineTuningJobWandbIntegrationObject struct {
	// The type of the integration being enabled for the fine-tuning job
	//
	// This field can be elided, and will be automatically set as "wandb".
	Type constant.Wandb `json:"type,required"`
	// The settings for your integration with Weights and Biases. This payload
	// specifies the project that metrics will be sent to. Optionally, you can set an
	// explicit display name for your run, add tags to your run, and set a default
	// entity (team, username, etc) to be associated with your run.
	Wandb FineTuningJobWandbIntegration `json:"wandb,omitzero,required"`
	JSON  struct {
		Type  resp.Field
		Wandb resp.Field
		raw   string
	} `json:"-"`
}

func (r FineTuningJobWandbIntegrationObject) RawJSON() string { return r.JSON.raw }
func (r *FineTuningJobWandbIntegrationObject) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The settings for your integration with Weights and Biases. This payload
// specifies the project that metrics will be sent to. Optionally, you can set an
// explicit display name for your run, add tags to your run, and set a default
// entity (team, username, etc) to be associated with your run.
type FineTuningJobWandbIntegration struct {
	// The name of the project that the new run will be created under.
	Project string `json:"project,omitzero,required"`
	// The entity to use for the run. This allows you to set the team or username of
	// the WandB user that you would like associated with the run. If not set, the
	// default entity for the registered WandB API key is used.
	Entity string `json:"entity,omitzero,nullable"`
	// A display name to set for the run. If not set, we will use the Job ID as the
	// name.
	Name string `json:"name,omitzero,nullable"`
	// A list of tags to be attached to the newly created run. These tags are passed
	// through directly to WandB. Some default tags are generated by OpenAI:
	// "openai/finetune", "openai/{base-model}", "openai/{ftjob-abcdef}".
	Tags []string `json:"tags,omitzero"`
	JSON struct {
		Project resp.Field
		Entity  resp.Field
		Name    resp.Field
		Tags    resp.Field
		raw     string
	} `json:"-"`
}

func (r FineTuningJobWandbIntegration) RawJSON() string { return r.JSON.raw }
func (r *FineTuningJobWandbIntegration) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type FineTuningJobNewParams struct {
	// The name of the model to fine-tune. You can select one of the
	// [supported models](https://platform.openai.com/docs/guides/fine-tuning#which-models-can-be-fine-tuned).
	Model string `json:"model,omitzero,required"`
	// The ID of an uploaded file that contains training data.
	//
	// See [upload file](https://platform.openai.com/docs/api-reference/files/create)
	// for how to upload a file.
	//
	// Your dataset must be formatted as a JSONL file. Additionally, you must upload
	// your file with the purpose `fine-tune`.
	//
	// The contents of the file should differ depending on if the model uses the
	// [chat](https://platform.openai.com/docs/api-reference/fine-tuning/chat-input),
	// [completions](https://platform.openai.com/docs/api-reference/fine-tuning/completions-input)
	// format, or if the fine-tuning method uses the
	// [preference](https://platform.openai.com/docs/api-reference/fine-tuning/preference-input)
	// format.
	//
	// See the [fine-tuning guide](https://platform.openai.com/docs/guides/fine-tuning)
	// for more details.
	TrainingFile param.String `json:"training_file,omitzero,required"`
	// The hyperparameters used for the fine-tuning job. This value is now deprecated
	// in favor of `method`, and should be passed in under the `method` parameter.
	Hyperparameters FineTuningJobNewParamsHyperparameters `json:"hyperparameters,omitzero"`
	// A list of integrations to enable for your fine-tuning job.
	Integrations []FineTuningJobNewParamsIntegration `json:"integrations,omitzero"`
	// The method used for fine-tuning.
	Method FineTuningJobNewParamsMethod `json:"method,omitzero"`
	// The seed controls the reproducibility of the job. Passing in the same seed and
	// job parameters should produce the same results, but may differ in rare cases. If
	// a seed is not specified, one will be generated for you.
	Seed param.Int `json:"seed,omitzero"`
	// A string of up to 64 characters that will be added to your fine-tuned model
	// name.
	//
	// For example, a `suffix` of "custom-model-name" would produce a model name like
	// `ft:gpt-4o-mini:openai:custom-model-name:7p4lURel`.
	Suffix param.String `json:"suffix,omitzero"`
	// The ID of an uploaded file that contains validation data.
	//
	// If you provide this file, the data is used to generate validation metrics
	// periodically during fine-tuning. These metrics can be viewed in the fine-tuning
	// results file. The same data should not be present in both train and validation
	// files.
	//
	// Your dataset must be formatted as a JSONL file. You must upload your file with
	// the purpose `fine-tune`.
	//
	// See the [fine-tuning guide](https://platform.openai.com/docs/guides/fine-tuning)
	// for more details.
	ValidationFile param.String `json:"validation_file,omitzero"`
	apiobject
}

func (f FineTuningJobNewParams) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

func (r FineTuningJobNewParams) MarshalJSON() (data []byte, err error) {
	type shadow FineTuningJobNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}

// The name of the model to fine-tune. You can select one of the
// [supported models](https://platform.openai.com/docs/guides/fine-tuning#which-models-can-be-fine-tuned).
type FineTuningJobNewParamsModel = string

const (
	FineTuningJobNewParamsModelBabbage002  FineTuningJobNewParamsModel = "babbage-002"
	FineTuningJobNewParamsModelDavinci002  FineTuningJobNewParamsModel = "davinci-002"
	FineTuningJobNewParamsModelGPT3_5Turbo FineTuningJobNewParamsModel = "gpt-3.5-turbo"
	FineTuningJobNewParamsModelGPT4oMini   FineTuningJobNewParamsModel = "gpt-4o-mini"
)

// The hyperparameters used for the fine-tuning job. This value is now deprecated
// in favor of `method`, and should be passed in under the `method` parameter.
//
// Deprecated: deprecated
type FineTuningJobNewParamsHyperparameters struct {
	// Number of examples in each batch. A larger batch size means that model
	// parameters are updated less frequently, but with lower variance.
	BatchSize FineTuningJobNewParamsHyperparametersBatchSizeUnion `json:"batch_size,omitzero"`
	// Scaling factor for the learning rate. A smaller learning rate may be useful to
	// avoid overfitting.
	LearningRateMultiplier FineTuningJobNewParamsHyperparametersLearningRateMultiplierUnion `json:"learning_rate_multiplier,omitzero"`
	// The number of epochs to train the model for. An epoch refers to one full cycle
	// through the training dataset.
	NEpochs FineTuningJobNewParamsHyperparametersNEpochsUnion `json:"n_epochs,omitzero"`
	apiobject
}

func (f FineTuningJobNewParamsHyperparameters) IsMissing() bool {
	return param.IsOmitted(f) || f.IsNull()
}

func (r FineTuningJobNewParamsHyperparameters) MarshalJSON() (data []byte, err error) {
	type shadow FineTuningJobNewParamsHyperparameters
	return param.MarshalObject(r, (*shadow)(&r))
}

// Only one field can be non-zero
type FineTuningJobNewParamsHyperparametersBatchSizeUnion struct {
	// Construct this variant with constant.ValueOf[constant.Auto]() Check if union is
	// this variant with !param.IsOmitted(union.OfAuto)
	OfAuto constant.Auto
	OfInt  param.Int
	apiunion
}

func (u FineTuningJobNewParamsHyperparametersBatchSizeUnion) IsMissing() bool {
	return param.IsOmitted(u) || u.IsNull()
}

func (u FineTuningJobNewParamsHyperparametersBatchSizeUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion[FineTuningJobNewParamsHyperparametersBatchSizeUnion](u.OfAuto, u.OfInt)
}

// Only one field can be non-zero
type FineTuningJobNewParamsHyperparametersLearningRateMultiplierUnion struct {
	// Construct this variant with constant.ValueOf[constant.Auto]() Check if union is
	// this variant with !param.IsOmitted(union.OfAuto)
	OfAuto  constant.Auto
	OfFloat param.Float
	apiunion
}

func (u FineTuningJobNewParamsHyperparametersLearningRateMultiplierUnion) IsMissing() bool {
	return param.IsOmitted(u) || u.IsNull()
}

func (u FineTuningJobNewParamsHyperparametersLearningRateMultiplierUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion[FineTuningJobNewParamsHyperparametersLearningRateMultiplierUnion](u.OfAuto, u.OfFloat)
}

// Only one field can be non-zero
type FineTuningJobNewParamsHyperparametersNEpochsUnion struct {
	// Construct this variant with constant.ValueOf[constant.Auto]() Check if union is
	// this variant with !param.IsOmitted(union.OfAuto)
	OfAuto constant.Auto
	OfInt  param.Int
	apiunion
}

func (u FineTuningJobNewParamsHyperparametersNEpochsUnion) IsMissing() bool {
	return param.IsOmitted(u) || u.IsNull()
}

func (u FineTuningJobNewParamsHyperparametersNEpochsUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion[FineTuningJobNewParamsHyperparametersNEpochsUnion](u.OfAuto, u.OfInt)
}

type FineTuningJobNewParamsIntegration struct {
	// The type of integration to enable. Currently, only "wandb" (Weights and Biases)
	// is supported.
	//
	// This field can be elided, and will be automatically set as "wandb".
	Type constant.Wandb `json:"type,required"`
	// The settings for your integration with Weights and Biases. This payload
	// specifies the project that metrics will be sent to. Optionally, you can set an
	// explicit display name for your run, add tags to your run, and set a default
	// entity (team, username, etc) to be associated with your run.
	Wandb FineTuningJobNewParamsIntegrationsWandb `json:"wandb,omitzero,required"`
	apiobject
}

func (f FineTuningJobNewParamsIntegration) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

func (r FineTuningJobNewParamsIntegration) MarshalJSON() (data []byte, err error) {
	type shadow FineTuningJobNewParamsIntegration
	return param.MarshalObject(r, (*shadow)(&r))
}

// The settings for your integration with Weights and Biases. This payload
// specifies the project that metrics will be sent to. Optionally, you can set an
// explicit display name for your run, add tags to your run, and set a default
// entity (team, username, etc) to be associated with your run.
type FineTuningJobNewParamsIntegrationsWandb struct {
	// The name of the project that the new run will be created under.
	Project param.String `json:"project,omitzero,required"`
	// The entity to use for the run. This allows you to set the team or username of
	// the WandB user that you would like associated with the run. If not set, the
	// default entity for the registered WandB API key is used.
	Entity param.String `json:"entity,omitzero"`
	// A display name to set for the run. If not set, we will use the Job ID as the
	// name.
	Name param.String `json:"name,omitzero"`
	// A list of tags to be attached to the newly created run. These tags are passed
	// through directly to WandB. Some default tags are generated by OpenAI:
	// "openai/finetune", "openai/{base-model}", "openai/{ftjob-abcdef}".
	Tags []string `json:"tags,omitzero"`
	apiobject
}

func (f FineTuningJobNewParamsIntegrationsWandb) IsMissing() bool {
	return param.IsOmitted(f) || f.IsNull()
}

func (r FineTuningJobNewParamsIntegrationsWandb) MarshalJSON() (data []byte, err error) {
	type shadow FineTuningJobNewParamsIntegrationsWandb
	return param.MarshalObject(r, (*shadow)(&r))
}

// The method used for fine-tuning.
type FineTuningJobNewParamsMethod struct {
	// Configuration for the DPO fine-tuning method.
	Dpo FineTuningJobNewParamsMethodDpo `json:"dpo,omitzero"`
	// Configuration for the supervised fine-tuning method.
	Supervised FineTuningJobNewParamsMethodSupervised `json:"supervised,omitzero"`
	// The type of method. Is either `supervised` or `dpo`.
	//
	// Any of "supervised", "dpo"
	Type string `json:"type,omitzero"`
	apiobject
}

func (f FineTuningJobNewParamsMethod) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

func (r FineTuningJobNewParamsMethod) MarshalJSON() (data []byte, err error) {
	type shadow FineTuningJobNewParamsMethod
	return param.MarshalObject(r, (*shadow)(&r))
}

// Configuration for the DPO fine-tuning method.
type FineTuningJobNewParamsMethodDpo struct {
	// The hyperparameters used for the fine-tuning job.
	Hyperparameters FineTuningJobNewParamsMethodDpoHyperparameters `json:"hyperparameters,omitzero"`
	apiobject
}

func (f FineTuningJobNewParamsMethodDpo) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

func (r FineTuningJobNewParamsMethodDpo) MarshalJSON() (data []byte, err error) {
	type shadow FineTuningJobNewParamsMethodDpo
	return param.MarshalObject(r, (*shadow)(&r))
}

// The hyperparameters used for the fine-tuning job.
type FineTuningJobNewParamsMethodDpoHyperparameters struct {
	// Number of examples in each batch. A larger batch size means that model
	// parameters are updated less frequently, but with lower variance.
	BatchSize FineTuningJobNewParamsMethodDpoHyperparametersBatchSizeUnion `json:"batch_size,omitzero"`
	// The beta value for the DPO method. A higher beta value will increase the weight
	// of the penalty between the policy and reference model.
	Beta FineTuningJobNewParamsMethodDpoHyperparametersBetaUnion `json:"beta,omitzero"`
	// Scaling factor for the learning rate. A smaller learning rate may be useful to
	// avoid overfitting.
	LearningRateMultiplier FineTuningJobNewParamsMethodDpoHyperparametersLearningRateMultiplierUnion `json:"learning_rate_multiplier,omitzero"`
	// The number of epochs to train the model for. An epoch refers to one full cycle
	// through the training dataset.
	NEpochs FineTuningJobNewParamsMethodDpoHyperparametersNEpochsUnion `json:"n_epochs,omitzero"`
	apiobject
}

func (f FineTuningJobNewParamsMethodDpoHyperparameters) IsMissing() bool {
	return param.IsOmitted(f) || f.IsNull()
}

func (r FineTuningJobNewParamsMethodDpoHyperparameters) MarshalJSON() (data []byte, err error) {
	type shadow FineTuningJobNewParamsMethodDpoHyperparameters
	return param.MarshalObject(r, (*shadow)(&r))
}

// Only one field can be non-zero
type FineTuningJobNewParamsMethodDpoHyperparametersBatchSizeUnion struct {
	// Construct this variant with constant.ValueOf[constant.Auto]() Check if union is
	// this variant with !param.IsOmitted(union.OfAuto)
	OfAuto constant.Auto
	OfInt  param.Int
	apiunion
}

func (u FineTuningJobNewParamsMethodDpoHyperparametersBatchSizeUnion) IsMissing() bool {
	return param.IsOmitted(u) || u.IsNull()
}

func (u FineTuningJobNewParamsMethodDpoHyperparametersBatchSizeUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion[FineTuningJobNewParamsMethodDpoHyperparametersBatchSizeUnion](u.OfAuto, u.OfInt)
}

// Only one field can be non-zero
type FineTuningJobNewParamsMethodDpoHyperparametersBetaUnion struct {
	// Construct this variant with constant.ValueOf[constant.Auto]() Check if union is
	// this variant with !param.IsOmitted(union.OfAuto)
	OfAuto  constant.Auto
	OfFloat param.Float
	apiunion
}

func (u FineTuningJobNewParamsMethodDpoHyperparametersBetaUnion) IsMissing() bool {
	return param.IsOmitted(u) || u.IsNull()
}

func (u FineTuningJobNewParamsMethodDpoHyperparametersBetaUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion[FineTuningJobNewParamsMethodDpoHyperparametersBetaUnion](u.OfAuto, u.OfFloat)
}

// Only one field can be non-zero
type FineTuningJobNewParamsMethodDpoHyperparametersLearningRateMultiplierUnion struct {
	// Construct this variant with constant.ValueOf[constant.Auto]() Check if union is
	// this variant with !param.IsOmitted(union.OfAuto)
	OfAuto  constant.Auto
	OfFloat param.Float
	apiunion
}

func (u FineTuningJobNewParamsMethodDpoHyperparametersLearningRateMultiplierUnion) IsMissing() bool {
	return param.IsOmitted(u) || u.IsNull()
}

func (u FineTuningJobNewParamsMethodDpoHyperparametersLearningRateMultiplierUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion[FineTuningJobNewParamsMethodDpoHyperparametersLearningRateMultiplierUnion](u.OfAuto, u.OfFloat)
}

// Only one field can be non-zero
type FineTuningJobNewParamsMethodDpoHyperparametersNEpochsUnion struct {
	// Construct this variant with constant.ValueOf[constant.Auto]() Check if union is
	// this variant with !param.IsOmitted(union.OfAuto)
	OfAuto constant.Auto
	OfInt  param.Int
	apiunion
}

func (u FineTuningJobNewParamsMethodDpoHyperparametersNEpochsUnion) IsMissing() bool {
	return param.IsOmitted(u) || u.IsNull()
}

func (u FineTuningJobNewParamsMethodDpoHyperparametersNEpochsUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion[FineTuningJobNewParamsMethodDpoHyperparametersNEpochsUnion](u.OfAuto, u.OfInt)
}

// Configuration for the supervised fine-tuning method.
type FineTuningJobNewParamsMethodSupervised struct {
	// The hyperparameters used for the fine-tuning job.
	Hyperparameters FineTuningJobNewParamsMethodSupervisedHyperparameters `json:"hyperparameters,omitzero"`
	apiobject
}

func (f FineTuningJobNewParamsMethodSupervised) IsMissing() bool {
	return param.IsOmitted(f) || f.IsNull()
}

func (r FineTuningJobNewParamsMethodSupervised) MarshalJSON() (data []byte, err error) {
	type shadow FineTuningJobNewParamsMethodSupervised
	return param.MarshalObject(r, (*shadow)(&r))
}

// The hyperparameters used for the fine-tuning job.
type FineTuningJobNewParamsMethodSupervisedHyperparameters struct {
	// Number of examples in each batch. A larger batch size means that model
	// parameters are updated less frequently, but with lower variance.
	BatchSize FineTuningJobNewParamsMethodSupervisedHyperparametersBatchSizeUnion `json:"batch_size,omitzero"`
	// Scaling factor for the learning rate. A smaller learning rate may be useful to
	// avoid overfitting.
	LearningRateMultiplier FineTuningJobNewParamsMethodSupervisedHyperparametersLearningRateMultiplierUnion `json:"learning_rate_multiplier,omitzero"`
	// The number of epochs to train the model for. An epoch refers to one full cycle
	// through the training dataset.
	NEpochs FineTuningJobNewParamsMethodSupervisedHyperparametersNEpochsUnion `json:"n_epochs,omitzero"`
	apiobject
}

func (f FineTuningJobNewParamsMethodSupervisedHyperparameters) IsMissing() bool {
	return param.IsOmitted(f) || f.IsNull()
}

func (r FineTuningJobNewParamsMethodSupervisedHyperparameters) MarshalJSON() (data []byte, err error) {
	type shadow FineTuningJobNewParamsMethodSupervisedHyperparameters
	return param.MarshalObject(r, (*shadow)(&r))
}

// Only one field can be non-zero
type FineTuningJobNewParamsMethodSupervisedHyperparametersBatchSizeUnion struct {
	// Construct this variant with constant.ValueOf[constant.Auto]() Check if union is
	// this variant with !param.IsOmitted(union.OfAuto)
	OfAuto constant.Auto
	OfInt  param.Int
	apiunion
}

func (u FineTuningJobNewParamsMethodSupervisedHyperparametersBatchSizeUnion) IsMissing() bool {
	return param.IsOmitted(u) || u.IsNull()
}

func (u FineTuningJobNewParamsMethodSupervisedHyperparametersBatchSizeUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion[FineTuningJobNewParamsMethodSupervisedHyperparametersBatchSizeUnion](u.OfAuto, u.OfInt)
}

// Only one field can be non-zero
type FineTuningJobNewParamsMethodSupervisedHyperparametersLearningRateMultiplierUnion struct {
	// Construct this variant with constant.ValueOf[constant.Auto]() Check if union is
	// this variant with !param.IsOmitted(union.OfAuto)
	OfAuto  constant.Auto
	OfFloat param.Float
	apiunion
}

func (u FineTuningJobNewParamsMethodSupervisedHyperparametersLearningRateMultiplierUnion) IsMissing() bool {
	return param.IsOmitted(u) || u.IsNull()
}

func (u FineTuningJobNewParamsMethodSupervisedHyperparametersLearningRateMultiplierUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion[FineTuningJobNewParamsMethodSupervisedHyperparametersLearningRateMultiplierUnion](u.OfAuto, u.OfFloat)
}

// Only one field can be non-zero
type FineTuningJobNewParamsMethodSupervisedHyperparametersNEpochsUnion struct {
	// Construct this variant with constant.ValueOf[constant.Auto]() Check if union is
	// this variant with !param.IsOmitted(union.OfAuto)
	OfAuto constant.Auto
	OfInt  param.Int
	apiunion
}

func (u FineTuningJobNewParamsMethodSupervisedHyperparametersNEpochsUnion) IsMissing() bool {
	return param.IsOmitted(u) || u.IsNull()
}

func (u FineTuningJobNewParamsMethodSupervisedHyperparametersNEpochsUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion[FineTuningJobNewParamsMethodSupervisedHyperparametersNEpochsUnion](u.OfAuto, u.OfInt)
}

// The type of method. Is either `supervised` or `dpo`.
type FineTuningJobNewParamsMethodType = string

const (
	FineTuningJobNewParamsMethodTypeSupervised FineTuningJobNewParamsMethodType = "supervised"
	FineTuningJobNewParamsMethodTypeDpo        FineTuningJobNewParamsMethodType = "dpo"
)

type FineTuningJobListParams struct {
	// Identifier for the last job from the previous pagination request.
	After param.String `query:"after,omitzero"`
	// Number of fine-tuning jobs to retrieve.
	Limit param.Int `query:"limit,omitzero"`
	apiobject
}

func (f FineTuningJobListParams) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

// URLQuery serializes [FineTuningJobListParams]'s query parameters as
// `url.Values`.
func (r FineTuningJobListParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type FineTuningJobListEventsParams struct {
	// Identifier for the last event from the previous pagination request.
	After param.String `query:"after,omitzero"`
	// Number of events to retrieve.
	Limit param.Int `query:"limit,omitzero"`
	apiobject
}

func (f FineTuningJobListEventsParams) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

// URLQuery serializes [FineTuningJobListEventsParams]'s query parameters as
// `url.Values`.
func (r FineTuningJobListEventsParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

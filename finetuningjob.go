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
	"github.com/openai/openai-go/internal/param"
	"github.com/openai/openai-go/internal/requestconfig"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/pagination"
	"github.com/openai/openai-go/shared"
	"github.com/tidwall/gjson"
)

// FineTuningJobService contains methods and other services that help with
// interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewFineTuningJobService] method instead.
type FineTuningJobService struct {
	Options     []option.RequestOption
	Checkpoints *FineTuningJobCheckpointService
}

// NewFineTuningJobService generates a new service that applies the given options
// to each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewFineTuningJobService(opts ...option.RequestOption) (r *FineTuningJobService) {
	r = &FineTuningJobService{}
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
	ID string `json:"id,required"`
	// The Unix timestamp (in seconds) for when the fine-tuning job was created.
	CreatedAt int64 `json:"created_at,required"`
	// For fine-tuning jobs that have `failed`, this will contain more information on
	// the cause of the failure.
	Error FineTuningJobError `json:"error,required,nullable"`
	// The name of the fine-tuned model that is being created. The value will be null
	// if the fine-tuning job is still running.
	FineTunedModel string `json:"fine_tuned_model,required,nullable"`
	// The Unix timestamp (in seconds) for when the fine-tuning job was finished. The
	// value will be null if the fine-tuning job is still running.
	FinishedAt int64 `json:"finished_at,required,nullable"`
	// The hyperparameters used for the fine-tuning job. This value will only be
	// returned when running `supervised` jobs.
	Hyperparameters FineTuningJobHyperparameters `json:"hyperparameters,required"`
	// The base model that is being fine-tuned.
	Model string `json:"model,required"`
	// The object type, which is always "fine_tuning.job".
	Object FineTuningJobObject `json:"object,required"`
	// The organization that owns the fine-tuning job.
	OrganizationID string `json:"organization_id,required"`
	// The compiled results file ID(s) for the fine-tuning job. You can retrieve the
	// results with the
	// [Files API](https://platform.openai.com/docs/api-reference/files/retrieve-contents).
	ResultFiles []string `json:"result_files,required"`
	// The seed used for the fine-tuning job.
	Seed int64 `json:"seed,required"`
	// The current status of the fine-tuning job, which can be either
	// `validating_files`, `queued`, `running`, `succeeded`, `failed`, or `cancelled`.
	Status FineTuningJobStatus `json:"status,required"`
	// The total number of billable tokens processed by this fine-tuning job. The value
	// will be null if the fine-tuning job is still running.
	TrainedTokens int64 `json:"trained_tokens,required,nullable"`
	// The file ID used for training. You can retrieve the training data with the
	// [Files API](https://platform.openai.com/docs/api-reference/files/retrieve-contents).
	TrainingFile string `json:"training_file,required"`
	// The file ID used for validation. You can retrieve the validation results with
	// the
	// [Files API](https://platform.openai.com/docs/api-reference/files/retrieve-contents).
	ValidationFile string `json:"validation_file,required,nullable"`
	// The Unix timestamp (in seconds) for when the fine-tuning job is estimated to
	// finish. The value will be null if the fine-tuning job is not running.
	EstimatedFinish int64 `json:"estimated_finish,nullable"`
	// A list of integrations to enable for this fine-tuning job.
	Integrations []FineTuningJobWandbIntegrationObject `json:"integrations,nullable"`
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard.
	//
	// Keys are strings with a maximum length of 64 characters. Values are strings with
	// a maximum length of 512 characters.
	Metadata shared.Metadata `json:"metadata,nullable"`
	// The method used for fine-tuning.
	Method FineTuningJobMethod `json:"method"`
	JSON   fineTuningJobJSON   `json:"-"`
}

// fineTuningJobJSON contains the JSON metadata for the struct [FineTuningJob]
type fineTuningJobJSON struct {
	ID              apijson.Field
	CreatedAt       apijson.Field
	Error           apijson.Field
	FineTunedModel  apijson.Field
	FinishedAt      apijson.Field
	Hyperparameters apijson.Field
	Model           apijson.Field
	Object          apijson.Field
	OrganizationID  apijson.Field
	ResultFiles     apijson.Field
	Seed            apijson.Field
	Status          apijson.Field
	TrainedTokens   apijson.Field
	TrainingFile    apijson.Field
	ValidationFile  apijson.Field
	EstimatedFinish apijson.Field
	Integrations    apijson.Field
	Metadata        apijson.Field
	Method          apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r *FineTuningJob) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r fineTuningJobJSON) RawJSON() string {
	return r.raw
}

// For fine-tuning jobs that have `failed`, this will contain more information on
// the cause of the failure.
type FineTuningJobError struct {
	// A machine-readable error code.
	Code string `json:"code,required"`
	// A human-readable error message.
	Message string `json:"message,required"`
	// The parameter that was invalid, usually `training_file` or `validation_file`.
	// This field will be null if the failure was not parameter-specific.
	Param string                 `json:"param,required,nullable"`
	JSON  fineTuningJobErrorJSON `json:"-"`
}

// fineTuningJobErrorJSON contains the JSON metadata for the struct
// [FineTuningJobError]
type fineTuningJobErrorJSON struct {
	Code        apijson.Field
	Message     apijson.Field
	Param       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *FineTuningJobError) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r fineTuningJobErrorJSON) RawJSON() string {
	return r.raw
}

// The hyperparameters used for the fine-tuning job. This value will only be
// returned when running `supervised` jobs.
type FineTuningJobHyperparameters struct {
	// Number of examples in each batch. A larger batch size means that model
	// parameters are updated less frequently, but with lower variance.
	BatchSize FineTuningJobHyperparametersBatchSizeUnion `json:"batch_size"`
	// Scaling factor for the learning rate. A smaller learning rate may be useful to
	// avoid overfitting.
	LearningRateMultiplier FineTuningJobHyperparametersLearningRateMultiplierUnion `json:"learning_rate_multiplier"`
	// The number of epochs to train the model for. An epoch refers to one full cycle
	// through the training dataset.
	NEpochs FineTuningJobHyperparametersNEpochsUnion `json:"n_epochs"`
	JSON    fineTuningJobHyperparametersJSON         `json:"-"`
}

// fineTuningJobHyperparametersJSON contains the JSON metadata for the struct
// [FineTuningJobHyperparameters]
type fineTuningJobHyperparametersJSON struct {
	BatchSize              apijson.Field
	LearningRateMultiplier apijson.Field
	NEpochs                apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *FineTuningJobHyperparameters) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r fineTuningJobHyperparametersJSON) RawJSON() string {
	return r.raw
}

// Number of examples in each batch. A larger batch size means that model
// parameters are updated less frequently, but with lower variance.
//
// Union satisfied by [FineTuningJobHyperparametersBatchSizeAuto] or
// [shared.UnionInt].
type FineTuningJobHyperparametersBatchSizeUnion interface {
	ImplementsFineTuningJobHyperparametersBatchSizeUnion()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*FineTuningJobHyperparametersBatchSizeUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.String,
			Type:       reflect.TypeOf(FineTuningJobHyperparametersBatchSizeAuto("")),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.Number,
			Type:       reflect.TypeOf(shared.UnionInt(0)),
		},
	)
}

type FineTuningJobHyperparametersBatchSizeAuto string

const (
	FineTuningJobHyperparametersBatchSizeAutoAuto FineTuningJobHyperparametersBatchSizeAuto = "auto"
)

func (r FineTuningJobHyperparametersBatchSizeAuto) IsKnown() bool {
	switch r {
	case FineTuningJobHyperparametersBatchSizeAutoAuto:
		return true
	}
	return false
}

func (r FineTuningJobHyperparametersBatchSizeAuto) ImplementsFineTuningJobHyperparametersBatchSizeUnion() {
}

// Scaling factor for the learning rate. A smaller learning rate may be useful to
// avoid overfitting.
//
// Union satisfied by [FineTuningJobHyperparametersLearningRateMultiplierAuto] or
// [shared.UnionFloat].
type FineTuningJobHyperparametersLearningRateMultiplierUnion interface {
	ImplementsFineTuningJobHyperparametersLearningRateMultiplierUnion()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*FineTuningJobHyperparametersLearningRateMultiplierUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.String,
			Type:       reflect.TypeOf(FineTuningJobHyperparametersLearningRateMultiplierAuto("")),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.Number,
			Type:       reflect.TypeOf(shared.UnionFloat(0)),
		},
	)
}

type FineTuningJobHyperparametersLearningRateMultiplierAuto string

const (
	FineTuningJobHyperparametersLearningRateMultiplierAutoAuto FineTuningJobHyperparametersLearningRateMultiplierAuto = "auto"
)

func (r FineTuningJobHyperparametersLearningRateMultiplierAuto) IsKnown() bool {
	switch r {
	case FineTuningJobHyperparametersLearningRateMultiplierAutoAuto:
		return true
	}
	return false
}

func (r FineTuningJobHyperparametersLearningRateMultiplierAuto) ImplementsFineTuningJobHyperparametersLearningRateMultiplierUnion() {
}

// The number of epochs to train the model for. An epoch refers to one full cycle
// through the training dataset.
//
// Union satisfied by [FineTuningJobHyperparametersNEpochsAuto] or
// [shared.UnionInt].
type FineTuningJobHyperparametersNEpochsUnion interface {
	ImplementsFineTuningJobHyperparametersNEpochsUnion()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*FineTuningJobHyperparametersNEpochsUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.String,
			Type:       reflect.TypeOf(FineTuningJobHyperparametersNEpochsAuto("")),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.Number,
			Type:       reflect.TypeOf(shared.UnionInt(0)),
		},
	)
}

type FineTuningJobHyperparametersNEpochsAuto string

const (
	FineTuningJobHyperparametersNEpochsAutoAuto FineTuningJobHyperparametersNEpochsAuto = "auto"
)

func (r FineTuningJobHyperparametersNEpochsAuto) IsKnown() bool {
	switch r {
	case FineTuningJobHyperparametersNEpochsAutoAuto:
		return true
	}
	return false
}

func (r FineTuningJobHyperparametersNEpochsAuto) ImplementsFineTuningJobHyperparametersNEpochsUnion() {
}

// The object type, which is always "fine_tuning.job".
type FineTuningJobObject string

const (
	FineTuningJobObjectFineTuningJob FineTuningJobObject = "fine_tuning.job"
)

func (r FineTuningJobObject) IsKnown() bool {
	switch r {
	case FineTuningJobObjectFineTuningJob:
		return true
	}
	return false
}

// The current status of the fine-tuning job, which can be either
// `validating_files`, `queued`, `running`, `succeeded`, `failed`, or `cancelled`.
type FineTuningJobStatus string

const (
	FineTuningJobStatusValidatingFiles FineTuningJobStatus = "validating_files"
	FineTuningJobStatusQueued          FineTuningJobStatus = "queued"
	FineTuningJobStatusRunning         FineTuningJobStatus = "running"
	FineTuningJobStatusSucceeded       FineTuningJobStatus = "succeeded"
	FineTuningJobStatusFailed          FineTuningJobStatus = "failed"
	FineTuningJobStatusCancelled       FineTuningJobStatus = "cancelled"
)

func (r FineTuningJobStatus) IsKnown() bool {
	switch r {
	case FineTuningJobStatusValidatingFiles, FineTuningJobStatusQueued, FineTuningJobStatusRunning, FineTuningJobStatusSucceeded, FineTuningJobStatusFailed, FineTuningJobStatusCancelled:
		return true
	}
	return false
}

// The method used for fine-tuning.
type FineTuningJobMethod struct {
	// Configuration for the DPO fine-tuning method.
	Dpo FineTuningJobMethodDpo `json:"dpo"`
	// Configuration for the supervised fine-tuning method.
	Supervised FineTuningJobMethodSupervised `json:"supervised"`
	// The type of method. Is either `supervised` or `dpo`.
	Type FineTuningJobMethodType `json:"type"`
	JSON fineTuningJobMethodJSON `json:"-"`
}

// fineTuningJobMethodJSON contains the JSON metadata for the struct
// [FineTuningJobMethod]
type fineTuningJobMethodJSON struct {
	Dpo         apijson.Field
	Supervised  apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *FineTuningJobMethod) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r fineTuningJobMethodJSON) RawJSON() string {
	return r.raw
}

// Configuration for the DPO fine-tuning method.
type FineTuningJobMethodDpo struct {
	// The hyperparameters used for the fine-tuning job.
	Hyperparameters FineTuningJobMethodDpoHyperparameters `json:"hyperparameters"`
	JSON            fineTuningJobMethodDpoJSON            `json:"-"`
}

// fineTuningJobMethodDpoJSON contains the JSON metadata for the struct
// [FineTuningJobMethodDpo]
type fineTuningJobMethodDpoJSON struct {
	Hyperparameters apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r *FineTuningJobMethodDpo) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r fineTuningJobMethodDpoJSON) RawJSON() string {
	return r.raw
}

// The hyperparameters used for the fine-tuning job.
type FineTuningJobMethodDpoHyperparameters struct {
	// Number of examples in each batch. A larger batch size means that model
	// parameters are updated less frequently, but with lower variance.
	BatchSize FineTuningJobMethodDpoHyperparametersBatchSizeUnion `json:"batch_size"`
	// The beta value for the DPO method. A higher beta value will increase the weight
	// of the penalty between the policy and reference model.
	Beta FineTuningJobMethodDpoHyperparametersBetaUnion `json:"beta"`
	// Scaling factor for the learning rate. A smaller learning rate may be useful to
	// avoid overfitting.
	LearningRateMultiplier FineTuningJobMethodDpoHyperparametersLearningRateMultiplierUnion `json:"learning_rate_multiplier"`
	// The number of epochs to train the model for. An epoch refers to one full cycle
	// through the training dataset.
	NEpochs FineTuningJobMethodDpoHyperparametersNEpochsUnion `json:"n_epochs"`
	JSON    fineTuningJobMethodDpoHyperparametersJSON         `json:"-"`
}

// fineTuningJobMethodDpoHyperparametersJSON contains the JSON metadata for the
// struct [FineTuningJobMethodDpoHyperparameters]
type fineTuningJobMethodDpoHyperparametersJSON struct {
	BatchSize              apijson.Field
	Beta                   apijson.Field
	LearningRateMultiplier apijson.Field
	NEpochs                apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *FineTuningJobMethodDpoHyperparameters) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r fineTuningJobMethodDpoHyperparametersJSON) RawJSON() string {
	return r.raw
}

// Number of examples in each batch. A larger batch size means that model
// parameters are updated less frequently, but with lower variance.
//
// Union satisfied by [FineTuningJobMethodDpoHyperparametersBatchSizeAuto] or
// [shared.UnionInt].
type FineTuningJobMethodDpoHyperparametersBatchSizeUnion interface {
	ImplementsFineTuningJobMethodDpoHyperparametersBatchSizeUnion()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*FineTuningJobMethodDpoHyperparametersBatchSizeUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.String,
			Type:       reflect.TypeOf(FineTuningJobMethodDpoHyperparametersBatchSizeAuto("")),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.Number,
			Type:       reflect.TypeOf(shared.UnionInt(0)),
		},
	)
}

type FineTuningJobMethodDpoHyperparametersBatchSizeAuto string

const (
	FineTuningJobMethodDpoHyperparametersBatchSizeAutoAuto FineTuningJobMethodDpoHyperparametersBatchSizeAuto = "auto"
)

func (r FineTuningJobMethodDpoHyperparametersBatchSizeAuto) IsKnown() bool {
	switch r {
	case FineTuningJobMethodDpoHyperparametersBatchSizeAutoAuto:
		return true
	}
	return false
}

func (r FineTuningJobMethodDpoHyperparametersBatchSizeAuto) ImplementsFineTuningJobMethodDpoHyperparametersBatchSizeUnion() {
}

// The beta value for the DPO method. A higher beta value will increase the weight
// of the penalty between the policy and reference model.
//
// Union satisfied by [FineTuningJobMethodDpoHyperparametersBetaAuto] or
// [shared.UnionFloat].
type FineTuningJobMethodDpoHyperparametersBetaUnion interface {
	ImplementsFineTuningJobMethodDpoHyperparametersBetaUnion()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*FineTuningJobMethodDpoHyperparametersBetaUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.String,
			Type:       reflect.TypeOf(FineTuningJobMethodDpoHyperparametersBetaAuto("")),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.Number,
			Type:       reflect.TypeOf(shared.UnionFloat(0)),
		},
	)
}

type FineTuningJobMethodDpoHyperparametersBetaAuto string

const (
	FineTuningJobMethodDpoHyperparametersBetaAutoAuto FineTuningJobMethodDpoHyperparametersBetaAuto = "auto"
)

func (r FineTuningJobMethodDpoHyperparametersBetaAuto) IsKnown() bool {
	switch r {
	case FineTuningJobMethodDpoHyperparametersBetaAutoAuto:
		return true
	}
	return false
}

func (r FineTuningJobMethodDpoHyperparametersBetaAuto) ImplementsFineTuningJobMethodDpoHyperparametersBetaUnion() {
}

// Scaling factor for the learning rate. A smaller learning rate may be useful to
// avoid overfitting.
//
// Union satisfied by
// [FineTuningJobMethodDpoHyperparametersLearningRateMultiplierAuto] or
// [shared.UnionFloat].
type FineTuningJobMethodDpoHyperparametersLearningRateMultiplierUnion interface {
	ImplementsFineTuningJobMethodDpoHyperparametersLearningRateMultiplierUnion()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*FineTuningJobMethodDpoHyperparametersLearningRateMultiplierUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.String,
			Type:       reflect.TypeOf(FineTuningJobMethodDpoHyperparametersLearningRateMultiplierAuto("")),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.Number,
			Type:       reflect.TypeOf(shared.UnionFloat(0)),
		},
	)
}

type FineTuningJobMethodDpoHyperparametersLearningRateMultiplierAuto string

const (
	FineTuningJobMethodDpoHyperparametersLearningRateMultiplierAutoAuto FineTuningJobMethodDpoHyperparametersLearningRateMultiplierAuto = "auto"
)

func (r FineTuningJobMethodDpoHyperparametersLearningRateMultiplierAuto) IsKnown() bool {
	switch r {
	case FineTuningJobMethodDpoHyperparametersLearningRateMultiplierAutoAuto:
		return true
	}
	return false
}

func (r FineTuningJobMethodDpoHyperparametersLearningRateMultiplierAuto) ImplementsFineTuningJobMethodDpoHyperparametersLearningRateMultiplierUnion() {
}

// The number of epochs to train the model for. An epoch refers to one full cycle
// through the training dataset.
//
// Union satisfied by [FineTuningJobMethodDpoHyperparametersNEpochsAuto] or
// [shared.UnionInt].
type FineTuningJobMethodDpoHyperparametersNEpochsUnion interface {
	ImplementsFineTuningJobMethodDpoHyperparametersNEpochsUnion()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*FineTuningJobMethodDpoHyperparametersNEpochsUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.String,
			Type:       reflect.TypeOf(FineTuningJobMethodDpoHyperparametersNEpochsAuto("")),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.Number,
			Type:       reflect.TypeOf(shared.UnionInt(0)),
		},
	)
}

type FineTuningJobMethodDpoHyperparametersNEpochsAuto string

const (
	FineTuningJobMethodDpoHyperparametersNEpochsAutoAuto FineTuningJobMethodDpoHyperparametersNEpochsAuto = "auto"
)

func (r FineTuningJobMethodDpoHyperparametersNEpochsAuto) IsKnown() bool {
	switch r {
	case FineTuningJobMethodDpoHyperparametersNEpochsAutoAuto:
		return true
	}
	return false
}

func (r FineTuningJobMethodDpoHyperparametersNEpochsAuto) ImplementsFineTuningJobMethodDpoHyperparametersNEpochsUnion() {
}

// Configuration for the supervised fine-tuning method.
type FineTuningJobMethodSupervised struct {
	// The hyperparameters used for the fine-tuning job.
	Hyperparameters FineTuningJobMethodSupervisedHyperparameters `json:"hyperparameters"`
	JSON            fineTuningJobMethodSupervisedJSON            `json:"-"`
}

// fineTuningJobMethodSupervisedJSON contains the JSON metadata for the struct
// [FineTuningJobMethodSupervised]
type fineTuningJobMethodSupervisedJSON struct {
	Hyperparameters apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r *FineTuningJobMethodSupervised) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r fineTuningJobMethodSupervisedJSON) RawJSON() string {
	return r.raw
}

// The hyperparameters used for the fine-tuning job.
type FineTuningJobMethodSupervisedHyperparameters struct {
	// Number of examples in each batch. A larger batch size means that model
	// parameters are updated less frequently, but with lower variance.
	BatchSize FineTuningJobMethodSupervisedHyperparametersBatchSizeUnion `json:"batch_size"`
	// Scaling factor for the learning rate. A smaller learning rate may be useful to
	// avoid overfitting.
	LearningRateMultiplier FineTuningJobMethodSupervisedHyperparametersLearningRateMultiplierUnion `json:"learning_rate_multiplier"`
	// The number of epochs to train the model for. An epoch refers to one full cycle
	// through the training dataset.
	NEpochs FineTuningJobMethodSupervisedHyperparametersNEpochsUnion `json:"n_epochs"`
	JSON    fineTuningJobMethodSupervisedHyperparametersJSON         `json:"-"`
}

// fineTuningJobMethodSupervisedHyperparametersJSON contains the JSON metadata for
// the struct [FineTuningJobMethodSupervisedHyperparameters]
type fineTuningJobMethodSupervisedHyperparametersJSON struct {
	BatchSize              apijson.Field
	LearningRateMultiplier apijson.Field
	NEpochs                apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *FineTuningJobMethodSupervisedHyperparameters) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r fineTuningJobMethodSupervisedHyperparametersJSON) RawJSON() string {
	return r.raw
}

// Number of examples in each batch. A larger batch size means that model
// parameters are updated less frequently, but with lower variance.
//
// Union satisfied by [FineTuningJobMethodSupervisedHyperparametersBatchSizeAuto]
// or [shared.UnionInt].
type FineTuningJobMethodSupervisedHyperparametersBatchSizeUnion interface {
	ImplementsFineTuningJobMethodSupervisedHyperparametersBatchSizeUnion()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*FineTuningJobMethodSupervisedHyperparametersBatchSizeUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.String,
			Type:       reflect.TypeOf(FineTuningJobMethodSupervisedHyperparametersBatchSizeAuto("")),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.Number,
			Type:       reflect.TypeOf(shared.UnionInt(0)),
		},
	)
}

type FineTuningJobMethodSupervisedHyperparametersBatchSizeAuto string

const (
	FineTuningJobMethodSupervisedHyperparametersBatchSizeAutoAuto FineTuningJobMethodSupervisedHyperparametersBatchSizeAuto = "auto"
)

func (r FineTuningJobMethodSupervisedHyperparametersBatchSizeAuto) IsKnown() bool {
	switch r {
	case FineTuningJobMethodSupervisedHyperparametersBatchSizeAutoAuto:
		return true
	}
	return false
}

func (r FineTuningJobMethodSupervisedHyperparametersBatchSizeAuto) ImplementsFineTuningJobMethodSupervisedHyperparametersBatchSizeUnion() {
}

// Scaling factor for the learning rate. A smaller learning rate may be useful to
// avoid overfitting.
//
// Union satisfied by
// [FineTuningJobMethodSupervisedHyperparametersLearningRateMultiplierAuto] or
// [shared.UnionFloat].
type FineTuningJobMethodSupervisedHyperparametersLearningRateMultiplierUnion interface {
	ImplementsFineTuningJobMethodSupervisedHyperparametersLearningRateMultiplierUnion()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*FineTuningJobMethodSupervisedHyperparametersLearningRateMultiplierUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.String,
			Type:       reflect.TypeOf(FineTuningJobMethodSupervisedHyperparametersLearningRateMultiplierAuto("")),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.Number,
			Type:       reflect.TypeOf(shared.UnionFloat(0)),
		},
	)
}

type FineTuningJobMethodSupervisedHyperparametersLearningRateMultiplierAuto string

const (
	FineTuningJobMethodSupervisedHyperparametersLearningRateMultiplierAutoAuto FineTuningJobMethodSupervisedHyperparametersLearningRateMultiplierAuto = "auto"
)

func (r FineTuningJobMethodSupervisedHyperparametersLearningRateMultiplierAuto) IsKnown() bool {
	switch r {
	case FineTuningJobMethodSupervisedHyperparametersLearningRateMultiplierAutoAuto:
		return true
	}
	return false
}

func (r FineTuningJobMethodSupervisedHyperparametersLearningRateMultiplierAuto) ImplementsFineTuningJobMethodSupervisedHyperparametersLearningRateMultiplierUnion() {
}

// The number of epochs to train the model for. An epoch refers to one full cycle
// through the training dataset.
//
// Union satisfied by [FineTuningJobMethodSupervisedHyperparametersNEpochsAuto] or
// [shared.UnionInt].
type FineTuningJobMethodSupervisedHyperparametersNEpochsUnion interface {
	ImplementsFineTuningJobMethodSupervisedHyperparametersNEpochsUnion()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*FineTuningJobMethodSupervisedHyperparametersNEpochsUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.String,
			Type:       reflect.TypeOf(FineTuningJobMethodSupervisedHyperparametersNEpochsAuto("")),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.Number,
			Type:       reflect.TypeOf(shared.UnionInt(0)),
		},
	)
}

type FineTuningJobMethodSupervisedHyperparametersNEpochsAuto string

const (
	FineTuningJobMethodSupervisedHyperparametersNEpochsAutoAuto FineTuningJobMethodSupervisedHyperparametersNEpochsAuto = "auto"
)

func (r FineTuningJobMethodSupervisedHyperparametersNEpochsAuto) IsKnown() bool {
	switch r {
	case FineTuningJobMethodSupervisedHyperparametersNEpochsAutoAuto:
		return true
	}
	return false
}

func (r FineTuningJobMethodSupervisedHyperparametersNEpochsAuto) ImplementsFineTuningJobMethodSupervisedHyperparametersNEpochsUnion() {
}

// The type of method. Is either `supervised` or `dpo`.
type FineTuningJobMethodType string

const (
	FineTuningJobMethodTypeSupervised FineTuningJobMethodType = "supervised"
	FineTuningJobMethodTypeDpo        FineTuningJobMethodType = "dpo"
)

func (r FineTuningJobMethodType) IsKnown() bool {
	switch r {
	case FineTuningJobMethodTypeSupervised, FineTuningJobMethodTypeDpo:
		return true
	}
	return false
}

// Fine-tuning job event object
type FineTuningJobEvent struct {
	// The object identifier.
	ID string `json:"id,required"`
	// The Unix timestamp (in seconds) for when the fine-tuning job was created.
	CreatedAt int64 `json:"created_at,required"`
	// The log level of the event.
	Level FineTuningJobEventLevel `json:"level,required"`
	// The message of the event.
	Message string `json:"message,required"`
	// The object type, which is always "fine_tuning.job.event".
	Object FineTuningJobEventObject `json:"object,required"`
	// The data associated with the event.
	Data interface{} `json:"data"`
	// The type of event.
	Type FineTuningJobEventType `json:"type"`
	JSON fineTuningJobEventJSON `json:"-"`
}

// fineTuningJobEventJSON contains the JSON metadata for the struct
// [FineTuningJobEvent]
type fineTuningJobEventJSON struct {
	ID          apijson.Field
	CreatedAt   apijson.Field
	Level       apijson.Field
	Message     apijson.Field
	Object      apijson.Field
	Data        apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *FineTuningJobEvent) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r fineTuningJobEventJSON) RawJSON() string {
	return r.raw
}

// The log level of the event.
type FineTuningJobEventLevel string

const (
	FineTuningJobEventLevelInfo  FineTuningJobEventLevel = "info"
	FineTuningJobEventLevelWarn  FineTuningJobEventLevel = "warn"
	FineTuningJobEventLevelError FineTuningJobEventLevel = "error"
)

func (r FineTuningJobEventLevel) IsKnown() bool {
	switch r {
	case FineTuningJobEventLevelInfo, FineTuningJobEventLevelWarn, FineTuningJobEventLevelError:
		return true
	}
	return false
}

// The object type, which is always "fine_tuning.job.event".
type FineTuningJobEventObject string

const (
	FineTuningJobEventObjectFineTuningJobEvent FineTuningJobEventObject = "fine_tuning.job.event"
)

func (r FineTuningJobEventObject) IsKnown() bool {
	switch r {
	case FineTuningJobEventObjectFineTuningJobEvent:
		return true
	}
	return false
}

// The type of event.
type FineTuningJobEventType string

const (
	FineTuningJobEventTypeMessage FineTuningJobEventType = "message"
	FineTuningJobEventTypeMetrics FineTuningJobEventType = "metrics"
)

func (r FineTuningJobEventType) IsKnown() bool {
	switch r {
	case FineTuningJobEventTypeMessage, FineTuningJobEventTypeMetrics:
		return true
	}
	return false
}

// The settings for your integration with Weights and Biases. This payload
// specifies the project that metrics will be sent to. Optionally, you can set an
// explicit display name for your run, add tags to your run, and set a default
// entity (team, username, etc) to be associated with your run.
type FineTuningJobWandbIntegration struct {
	// The name of the project that the new run will be created under.
	Project string `json:"project,required"`
	// The entity to use for the run. This allows you to set the team or username of
	// the WandB user that you would like associated with the run. If not set, the
	// default entity for the registered WandB API key is used.
	Entity string `json:"entity,nullable"`
	// A display name to set for the run. If not set, we will use the Job ID as the
	// name.
	Name string `json:"name,nullable"`
	// A list of tags to be attached to the newly created run. These tags are passed
	// through directly to WandB. Some default tags are generated by OpenAI:
	// "openai/finetune", "openai/{base-model}", "openai/{ftjob-abcdef}".
	Tags []string                          `json:"tags"`
	JSON fineTuningJobWandbIntegrationJSON `json:"-"`
}

// fineTuningJobWandbIntegrationJSON contains the JSON metadata for the struct
// [FineTuningJobWandbIntegration]
type fineTuningJobWandbIntegrationJSON struct {
	Project     apijson.Field
	Entity      apijson.Field
	Name        apijson.Field
	Tags        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *FineTuningJobWandbIntegration) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r fineTuningJobWandbIntegrationJSON) RawJSON() string {
	return r.raw
}

type FineTuningJobWandbIntegrationObject struct {
	// The type of the integration being enabled for the fine-tuning job
	Type FineTuningJobWandbIntegrationObjectType `json:"type,required"`
	// The settings for your integration with Weights and Biases. This payload
	// specifies the project that metrics will be sent to. Optionally, you can set an
	// explicit display name for your run, add tags to your run, and set a default
	// entity (team, username, etc) to be associated with your run.
	Wandb FineTuningJobWandbIntegration           `json:"wandb,required"`
	JSON  fineTuningJobWandbIntegrationObjectJSON `json:"-"`
}

// fineTuningJobWandbIntegrationObjectJSON contains the JSON metadata for the
// struct [FineTuningJobWandbIntegrationObject]
type fineTuningJobWandbIntegrationObjectJSON struct {
	Type        apijson.Field
	Wandb       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *FineTuningJobWandbIntegrationObject) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r fineTuningJobWandbIntegrationObjectJSON) RawJSON() string {
	return r.raw
}

// The type of the integration being enabled for the fine-tuning job
type FineTuningJobWandbIntegrationObjectType string

const (
	FineTuningJobWandbIntegrationObjectTypeWandb FineTuningJobWandbIntegrationObjectType = "wandb"
)

func (r FineTuningJobWandbIntegrationObjectType) IsKnown() bool {
	switch r {
	case FineTuningJobWandbIntegrationObjectTypeWandb:
		return true
	}
	return false
}

type FineTuningJobNewParams struct {
	// The name of the model to fine-tune. You can select one of the
	// [supported models](https://platform.openai.com/docs/guides/fine-tuning#which-models-can-be-fine-tuned).
	Model param.Field[FineTuningJobNewParamsModel] `json:"model,required"`
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
	TrainingFile param.Field[string] `json:"training_file,required"`
	// The hyperparameters used for the fine-tuning job. This value is now deprecated
	// in favor of `method`, and should be passed in under the `method` parameter.
	Hyperparameters param.Field[FineTuningJobNewParamsHyperparameters] `json:"hyperparameters"`
	// A list of integrations to enable for your fine-tuning job.
	Integrations param.Field[[]FineTuningJobNewParamsIntegration] `json:"integrations"`
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard.
	//
	// Keys are strings with a maximum length of 64 characters. Values are strings with
	// a maximum length of 512 characters.
	Metadata param.Field[shared.MetadataParam] `json:"metadata"`
	// The method used for fine-tuning.
	Method param.Field[FineTuningJobNewParamsMethod] `json:"method"`
	// The seed controls the reproducibility of the job. Passing in the same seed and
	// job parameters should produce the same results, but may differ in rare cases. If
	// a seed is not specified, one will be generated for you.
	Seed param.Field[int64] `json:"seed"`
	// A string of up to 64 characters that will be added to your fine-tuned model
	// name.
	//
	// For example, a `suffix` of "custom-model-name" would produce a model name like
	// `ft:gpt-4o-mini:openai:custom-model-name:7p4lURel`.
	Suffix param.Field[string] `json:"suffix"`
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
	ValidationFile param.Field[string] `json:"validation_file"`
}

func (r FineTuningJobNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// The name of the model to fine-tune. You can select one of the
// [supported models](https://platform.openai.com/docs/guides/fine-tuning#which-models-can-be-fine-tuned).
type FineTuningJobNewParamsModel string

const (
	FineTuningJobNewParamsModelBabbage002  FineTuningJobNewParamsModel = "babbage-002"
	FineTuningJobNewParamsModelDavinci002  FineTuningJobNewParamsModel = "davinci-002"
	FineTuningJobNewParamsModelGPT3_5Turbo FineTuningJobNewParamsModel = "gpt-3.5-turbo"
	FineTuningJobNewParamsModelGPT4oMini   FineTuningJobNewParamsModel = "gpt-4o-mini"
)

func (r FineTuningJobNewParamsModel) IsKnown() bool {
	switch r {
	case FineTuningJobNewParamsModelBabbage002, FineTuningJobNewParamsModelDavinci002, FineTuningJobNewParamsModelGPT3_5Turbo, FineTuningJobNewParamsModelGPT4oMini:
		return true
	}
	return false
}

// The hyperparameters used for the fine-tuning job. This value is now deprecated
// in favor of `method`, and should be passed in under the `method` parameter.
//
// Deprecated: deprecated
type FineTuningJobNewParamsHyperparameters struct {
	// Number of examples in each batch. A larger batch size means that model
	// parameters are updated less frequently, but with lower variance.
	BatchSize param.Field[FineTuningJobNewParamsHyperparametersBatchSizeUnion] `json:"batch_size"`
	// Scaling factor for the learning rate. A smaller learning rate may be useful to
	// avoid overfitting.
	LearningRateMultiplier param.Field[FineTuningJobNewParamsHyperparametersLearningRateMultiplierUnion] `json:"learning_rate_multiplier"`
	// The number of epochs to train the model for. An epoch refers to one full cycle
	// through the training dataset.
	NEpochs param.Field[FineTuningJobNewParamsHyperparametersNEpochsUnion] `json:"n_epochs"`
}

func (r FineTuningJobNewParamsHyperparameters) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Number of examples in each batch. A larger batch size means that model
// parameters are updated less frequently, but with lower variance.
//
// Satisfied by [FineTuningJobNewParamsHyperparametersBatchSizeAuto],
// [shared.UnionInt].
type FineTuningJobNewParamsHyperparametersBatchSizeUnion interface {
	ImplementsFineTuningJobNewParamsHyperparametersBatchSizeUnion()
}

type FineTuningJobNewParamsHyperparametersBatchSizeAuto string

const (
	FineTuningJobNewParamsHyperparametersBatchSizeAutoAuto FineTuningJobNewParamsHyperparametersBatchSizeAuto = "auto"
)

func (r FineTuningJobNewParamsHyperparametersBatchSizeAuto) IsKnown() bool {
	switch r {
	case FineTuningJobNewParamsHyperparametersBatchSizeAutoAuto:
		return true
	}
	return false
}

func (r FineTuningJobNewParamsHyperparametersBatchSizeAuto) ImplementsFineTuningJobNewParamsHyperparametersBatchSizeUnion() {
}

// Scaling factor for the learning rate. A smaller learning rate may be useful to
// avoid overfitting.
//
// Satisfied by [FineTuningJobNewParamsHyperparametersLearningRateMultiplierAuto],
// [shared.UnionFloat].
type FineTuningJobNewParamsHyperparametersLearningRateMultiplierUnion interface {
	ImplementsFineTuningJobNewParamsHyperparametersLearningRateMultiplierUnion()
}

type FineTuningJobNewParamsHyperparametersLearningRateMultiplierAuto string

const (
	FineTuningJobNewParamsHyperparametersLearningRateMultiplierAutoAuto FineTuningJobNewParamsHyperparametersLearningRateMultiplierAuto = "auto"
)

func (r FineTuningJobNewParamsHyperparametersLearningRateMultiplierAuto) IsKnown() bool {
	switch r {
	case FineTuningJobNewParamsHyperparametersLearningRateMultiplierAutoAuto:
		return true
	}
	return false
}

func (r FineTuningJobNewParamsHyperparametersLearningRateMultiplierAuto) ImplementsFineTuningJobNewParamsHyperparametersLearningRateMultiplierUnion() {
}

// The number of epochs to train the model for. An epoch refers to one full cycle
// through the training dataset.
//
// Satisfied by [FineTuningJobNewParamsHyperparametersNEpochsAuto],
// [shared.UnionInt].
type FineTuningJobNewParamsHyperparametersNEpochsUnion interface {
	ImplementsFineTuningJobNewParamsHyperparametersNEpochsUnion()
}

type FineTuningJobNewParamsHyperparametersNEpochsAuto string

const (
	FineTuningJobNewParamsHyperparametersNEpochsAutoAuto FineTuningJobNewParamsHyperparametersNEpochsAuto = "auto"
)

func (r FineTuningJobNewParamsHyperparametersNEpochsAuto) IsKnown() bool {
	switch r {
	case FineTuningJobNewParamsHyperparametersNEpochsAutoAuto:
		return true
	}
	return false
}

func (r FineTuningJobNewParamsHyperparametersNEpochsAuto) ImplementsFineTuningJobNewParamsHyperparametersNEpochsUnion() {
}

type FineTuningJobNewParamsIntegration struct {
	// The type of integration to enable. Currently, only "wandb" (Weights and Biases)
	// is supported.
	Type param.Field[FineTuningJobNewParamsIntegrationsType] `json:"type,required"`
	// The settings for your integration with Weights and Biases. This payload
	// specifies the project that metrics will be sent to. Optionally, you can set an
	// explicit display name for your run, add tags to your run, and set a default
	// entity (team, username, etc) to be associated with your run.
	Wandb param.Field[FineTuningJobNewParamsIntegrationsWandb] `json:"wandb,required"`
}

func (r FineTuningJobNewParamsIntegration) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// The type of integration to enable. Currently, only "wandb" (Weights and Biases)
// is supported.
type FineTuningJobNewParamsIntegrationsType string

const (
	FineTuningJobNewParamsIntegrationsTypeWandb FineTuningJobNewParamsIntegrationsType = "wandb"
)

func (r FineTuningJobNewParamsIntegrationsType) IsKnown() bool {
	switch r {
	case FineTuningJobNewParamsIntegrationsTypeWandb:
		return true
	}
	return false
}

// The settings for your integration with Weights and Biases. This payload
// specifies the project that metrics will be sent to. Optionally, you can set an
// explicit display name for your run, add tags to your run, and set a default
// entity (team, username, etc) to be associated with your run.
type FineTuningJobNewParamsIntegrationsWandb struct {
	// The name of the project that the new run will be created under.
	Project param.Field[string] `json:"project,required"`
	// The entity to use for the run. This allows you to set the team or username of
	// the WandB user that you would like associated with the run. If not set, the
	// default entity for the registered WandB API key is used.
	Entity param.Field[string] `json:"entity"`
	// A display name to set for the run. If not set, we will use the Job ID as the
	// name.
	Name param.Field[string] `json:"name"`
	// A list of tags to be attached to the newly created run. These tags are passed
	// through directly to WandB. Some default tags are generated by OpenAI:
	// "openai/finetune", "openai/{base-model}", "openai/{ftjob-abcdef}".
	Tags param.Field[[]string] `json:"tags"`
}

func (r FineTuningJobNewParamsIntegrationsWandb) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// The method used for fine-tuning.
type FineTuningJobNewParamsMethod struct {
	// Configuration for the DPO fine-tuning method.
	Dpo param.Field[FineTuningJobNewParamsMethodDpo] `json:"dpo"`
	// Configuration for the supervised fine-tuning method.
	Supervised param.Field[FineTuningJobNewParamsMethodSupervised] `json:"supervised"`
	// The type of method. Is either `supervised` or `dpo`.
	Type param.Field[FineTuningJobNewParamsMethodType] `json:"type"`
}

func (r FineTuningJobNewParamsMethod) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Configuration for the DPO fine-tuning method.
type FineTuningJobNewParamsMethodDpo struct {
	// The hyperparameters used for the fine-tuning job.
	Hyperparameters param.Field[FineTuningJobNewParamsMethodDpoHyperparameters] `json:"hyperparameters"`
}

func (r FineTuningJobNewParamsMethodDpo) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// The hyperparameters used for the fine-tuning job.
type FineTuningJobNewParamsMethodDpoHyperparameters struct {
	// Number of examples in each batch. A larger batch size means that model
	// parameters are updated less frequently, but with lower variance.
	BatchSize param.Field[FineTuningJobNewParamsMethodDpoHyperparametersBatchSizeUnion] `json:"batch_size"`
	// The beta value for the DPO method. A higher beta value will increase the weight
	// of the penalty between the policy and reference model.
	Beta param.Field[FineTuningJobNewParamsMethodDpoHyperparametersBetaUnion] `json:"beta"`
	// Scaling factor for the learning rate. A smaller learning rate may be useful to
	// avoid overfitting.
	LearningRateMultiplier param.Field[FineTuningJobNewParamsMethodDpoHyperparametersLearningRateMultiplierUnion] `json:"learning_rate_multiplier"`
	// The number of epochs to train the model for. An epoch refers to one full cycle
	// through the training dataset.
	NEpochs param.Field[FineTuningJobNewParamsMethodDpoHyperparametersNEpochsUnion] `json:"n_epochs"`
}

func (r FineTuningJobNewParamsMethodDpoHyperparameters) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Number of examples in each batch. A larger batch size means that model
// parameters are updated less frequently, but with lower variance.
//
// Satisfied by [FineTuningJobNewParamsMethodDpoHyperparametersBatchSizeAuto],
// [shared.UnionInt].
type FineTuningJobNewParamsMethodDpoHyperparametersBatchSizeUnion interface {
	ImplementsFineTuningJobNewParamsMethodDpoHyperparametersBatchSizeUnion()
}

type FineTuningJobNewParamsMethodDpoHyperparametersBatchSizeAuto string

const (
	FineTuningJobNewParamsMethodDpoHyperparametersBatchSizeAutoAuto FineTuningJobNewParamsMethodDpoHyperparametersBatchSizeAuto = "auto"
)

func (r FineTuningJobNewParamsMethodDpoHyperparametersBatchSizeAuto) IsKnown() bool {
	switch r {
	case FineTuningJobNewParamsMethodDpoHyperparametersBatchSizeAutoAuto:
		return true
	}
	return false
}

func (r FineTuningJobNewParamsMethodDpoHyperparametersBatchSizeAuto) ImplementsFineTuningJobNewParamsMethodDpoHyperparametersBatchSizeUnion() {
}

// The beta value for the DPO method. A higher beta value will increase the weight
// of the penalty between the policy and reference model.
//
// Satisfied by [FineTuningJobNewParamsMethodDpoHyperparametersBetaAuto],
// [shared.UnionFloat].
type FineTuningJobNewParamsMethodDpoHyperparametersBetaUnion interface {
	ImplementsFineTuningJobNewParamsMethodDpoHyperparametersBetaUnion()
}

type FineTuningJobNewParamsMethodDpoHyperparametersBetaAuto string

const (
	FineTuningJobNewParamsMethodDpoHyperparametersBetaAutoAuto FineTuningJobNewParamsMethodDpoHyperparametersBetaAuto = "auto"
)

func (r FineTuningJobNewParamsMethodDpoHyperparametersBetaAuto) IsKnown() bool {
	switch r {
	case FineTuningJobNewParamsMethodDpoHyperparametersBetaAutoAuto:
		return true
	}
	return false
}

func (r FineTuningJobNewParamsMethodDpoHyperparametersBetaAuto) ImplementsFineTuningJobNewParamsMethodDpoHyperparametersBetaUnion() {
}

// Scaling factor for the learning rate. A smaller learning rate may be useful to
// avoid overfitting.
//
// Satisfied by
// [FineTuningJobNewParamsMethodDpoHyperparametersLearningRateMultiplierAuto],
// [shared.UnionFloat].
type FineTuningJobNewParamsMethodDpoHyperparametersLearningRateMultiplierUnion interface {
	ImplementsFineTuningJobNewParamsMethodDpoHyperparametersLearningRateMultiplierUnion()
}

type FineTuningJobNewParamsMethodDpoHyperparametersLearningRateMultiplierAuto string

const (
	FineTuningJobNewParamsMethodDpoHyperparametersLearningRateMultiplierAutoAuto FineTuningJobNewParamsMethodDpoHyperparametersLearningRateMultiplierAuto = "auto"
)

func (r FineTuningJobNewParamsMethodDpoHyperparametersLearningRateMultiplierAuto) IsKnown() bool {
	switch r {
	case FineTuningJobNewParamsMethodDpoHyperparametersLearningRateMultiplierAutoAuto:
		return true
	}
	return false
}

func (r FineTuningJobNewParamsMethodDpoHyperparametersLearningRateMultiplierAuto) ImplementsFineTuningJobNewParamsMethodDpoHyperparametersLearningRateMultiplierUnion() {
}

// The number of epochs to train the model for. An epoch refers to one full cycle
// through the training dataset.
//
// Satisfied by [FineTuningJobNewParamsMethodDpoHyperparametersNEpochsAuto],
// [shared.UnionInt].
type FineTuningJobNewParamsMethodDpoHyperparametersNEpochsUnion interface {
	ImplementsFineTuningJobNewParamsMethodDpoHyperparametersNEpochsUnion()
}

type FineTuningJobNewParamsMethodDpoHyperparametersNEpochsAuto string

const (
	FineTuningJobNewParamsMethodDpoHyperparametersNEpochsAutoAuto FineTuningJobNewParamsMethodDpoHyperparametersNEpochsAuto = "auto"
)

func (r FineTuningJobNewParamsMethodDpoHyperparametersNEpochsAuto) IsKnown() bool {
	switch r {
	case FineTuningJobNewParamsMethodDpoHyperparametersNEpochsAutoAuto:
		return true
	}
	return false
}

func (r FineTuningJobNewParamsMethodDpoHyperparametersNEpochsAuto) ImplementsFineTuningJobNewParamsMethodDpoHyperparametersNEpochsUnion() {
}

// Configuration for the supervised fine-tuning method.
type FineTuningJobNewParamsMethodSupervised struct {
	// The hyperparameters used for the fine-tuning job.
	Hyperparameters param.Field[FineTuningJobNewParamsMethodSupervisedHyperparameters] `json:"hyperparameters"`
}

func (r FineTuningJobNewParamsMethodSupervised) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// The hyperparameters used for the fine-tuning job.
type FineTuningJobNewParamsMethodSupervisedHyperparameters struct {
	// Number of examples in each batch. A larger batch size means that model
	// parameters are updated less frequently, but with lower variance.
	BatchSize param.Field[FineTuningJobNewParamsMethodSupervisedHyperparametersBatchSizeUnion] `json:"batch_size"`
	// Scaling factor for the learning rate. A smaller learning rate may be useful to
	// avoid overfitting.
	LearningRateMultiplier param.Field[FineTuningJobNewParamsMethodSupervisedHyperparametersLearningRateMultiplierUnion] `json:"learning_rate_multiplier"`
	// The number of epochs to train the model for. An epoch refers to one full cycle
	// through the training dataset.
	NEpochs param.Field[FineTuningJobNewParamsMethodSupervisedHyperparametersNEpochsUnion] `json:"n_epochs"`
}

func (r FineTuningJobNewParamsMethodSupervisedHyperparameters) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Number of examples in each batch. A larger batch size means that model
// parameters are updated less frequently, but with lower variance.
//
// Satisfied by
// [FineTuningJobNewParamsMethodSupervisedHyperparametersBatchSizeAuto],
// [shared.UnionInt].
type FineTuningJobNewParamsMethodSupervisedHyperparametersBatchSizeUnion interface {
	ImplementsFineTuningJobNewParamsMethodSupervisedHyperparametersBatchSizeUnion()
}

type FineTuningJobNewParamsMethodSupervisedHyperparametersBatchSizeAuto string

const (
	FineTuningJobNewParamsMethodSupervisedHyperparametersBatchSizeAutoAuto FineTuningJobNewParamsMethodSupervisedHyperparametersBatchSizeAuto = "auto"
)

func (r FineTuningJobNewParamsMethodSupervisedHyperparametersBatchSizeAuto) IsKnown() bool {
	switch r {
	case FineTuningJobNewParamsMethodSupervisedHyperparametersBatchSizeAutoAuto:
		return true
	}
	return false
}

func (r FineTuningJobNewParamsMethodSupervisedHyperparametersBatchSizeAuto) ImplementsFineTuningJobNewParamsMethodSupervisedHyperparametersBatchSizeUnion() {
}

// Scaling factor for the learning rate. A smaller learning rate may be useful to
// avoid overfitting.
//
// Satisfied by
// [FineTuningJobNewParamsMethodSupervisedHyperparametersLearningRateMultiplierAuto],
// [shared.UnionFloat].
type FineTuningJobNewParamsMethodSupervisedHyperparametersLearningRateMultiplierUnion interface {
	ImplementsFineTuningJobNewParamsMethodSupervisedHyperparametersLearningRateMultiplierUnion()
}

type FineTuningJobNewParamsMethodSupervisedHyperparametersLearningRateMultiplierAuto string

const (
	FineTuningJobNewParamsMethodSupervisedHyperparametersLearningRateMultiplierAutoAuto FineTuningJobNewParamsMethodSupervisedHyperparametersLearningRateMultiplierAuto = "auto"
)

func (r FineTuningJobNewParamsMethodSupervisedHyperparametersLearningRateMultiplierAuto) IsKnown() bool {
	switch r {
	case FineTuningJobNewParamsMethodSupervisedHyperparametersLearningRateMultiplierAutoAuto:
		return true
	}
	return false
}

func (r FineTuningJobNewParamsMethodSupervisedHyperparametersLearningRateMultiplierAuto) ImplementsFineTuningJobNewParamsMethodSupervisedHyperparametersLearningRateMultiplierUnion() {
}

// The number of epochs to train the model for. An epoch refers to one full cycle
// through the training dataset.
//
// Satisfied by [FineTuningJobNewParamsMethodSupervisedHyperparametersNEpochsAuto],
// [shared.UnionInt].
type FineTuningJobNewParamsMethodSupervisedHyperparametersNEpochsUnion interface {
	ImplementsFineTuningJobNewParamsMethodSupervisedHyperparametersNEpochsUnion()
}

type FineTuningJobNewParamsMethodSupervisedHyperparametersNEpochsAuto string

const (
	FineTuningJobNewParamsMethodSupervisedHyperparametersNEpochsAutoAuto FineTuningJobNewParamsMethodSupervisedHyperparametersNEpochsAuto = "auto"
)

func (r FineTuningJobNewParamsMethodSupervisedHyperparametersNEpochsAuto) IsKnown() bool {
	switch r {
	case FineTuningJobNewParamsMethodSupervisedHyperparametersNEpochsAutoAuto:
		return true
	}
	return false
}

func (r FineTuningJobNewParamsMethodSupervisedHyperparametersNEpochsAuto) ImplementsFineTuningJobNewParamsMethodSupervisedHyperparametersNEpochsUnion() {
}

// The type of method. Is either `supervised` or `dpo`.
type FineTuningJobNewParamsMethodType string

const (
	FineTuningJobNewParamsMethodTypeSupervised FineTuningJobNewParamsMethodType = "supervised"
	FineTuningJobNewParamsMethodTypeDpo        FineTuningJobNewParamsMethodType = "dpo"
)

func (r FineTuningJobNewParamsMethodType) IsKnown() bool {
	switch r {
	case FineTuningJobNewParamsMethodTypeSupervised, FineTuningJobNewParamsMethodTypeDpo:
		return true
	}
	return false
}

type FineTuningJobListParams struct {
	// Identifier for the last job from the previous pagination request.
	After param.Field[string] `query:"after"`
	// Number of fine-tuning jobs to retrieve.
	Limit param.Field[int64] `query:"limit"`
	// Optional metadata filter. To filter, use the syntax `metadata[k]=v`.
	// Alternatively, set `metadata=null` to indicate no metadata.
	Metadata param.Field[map[string]string] `query:"metadata"`
}

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
	After param.Field[string] `query:"after"`
	// Number of events to retrieve.
	Limit param.Field[int64] `query:"limit"`
}

// URLQuery serializes [FineTuningJobListEventsParams]'s query parameters as
// `url.Values`.
func (r FineTuningJobListEventsParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

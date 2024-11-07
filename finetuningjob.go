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
	// The hyperparameters used for the fine-tuning job. See the
	// [fine-tuning guide](https://platform.openai.com/docs/guides/fine-tuning) for
	// more details.
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
	JSON         fineTuningJobJSON                     `json:"-"`
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

// The hyperparameters used for the fine-tuning job. See the
// [fine-tuning guide](https://platform.openai.com/docs/guides/fine-tuning) for
// more details.
type FineTuningJobHyperparameters struct {
	// The number of epochs to train the model for. An epoch refers to one full cycle
	// through the training dataset. "auto" decides the optimal number of epochs based
	// on the size of the dataset. If setting the number manually, we support any
	// number between 1 and 50 epochs.
	NEpochs FineTuningJobHyperparametersNEpochsUnion `json:"n_epochs,required"`
	JSON    fineTuningJobHyperparametersJSON         `json:"-"`
}

// fineTuningJobHyperparametersJSON contains the JSON metadata for the struct
// [FineTuningJobHyperparameters]
type fineTuningJobHyperparametersJSON struct {
	NEpochs     apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *FineTuningJobHyperparameters) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r fineTuningJobHyperparametersJSON) RawJSON() string {
	return r.raw
}

// The number of epochs to train the model for. An epoch refers to one full cycle
// through the training dataset. "auto" decides the optimal number of epochs based
// on the size of the dataset. If setting the number manually, we support any
// number between 1 and 50 epochs.
//
// Union satisfied by [FineTuningJobHyperparametersNEpochsBehavior] or
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
			Type:       reflect.TypeOf(FineTuningJobHyperparametersNEpochsBehavior("")),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.Number,
			Type:       reflect.TypeOf(shared.UnionInt(0)),
		},
	)
}

type FineTuningJobHyperparametersNEpochsBehavior string

const (
	FineTuningJobHyperparametersNEpochsBehaviorAuto FineTuningJobHyperparametersNEpochsBehavior = "auto"
)

func (r FineTuningJobHyperparametersNEpochsBehavior) IsKnown() bool {
	switch r {
	case FineTuningJobHyperparametersNEpochsBehaviorAuto:
		return true
	}
	return false
}

func (r FineTuningJobHyperparametersNEpochsBehavior) ImplementsFineTuningJobHyperparametersNEpochsUnion() {
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

// Fine-tuning job event object
type FineTuningJobEvent struct {
	ID        string                   `json:"id,required"`
	CreatedAt int64                    `json:"created_at,required"`
	Level     FineTuningJobEventLevel  `json:"level,required"`
	Message   string                   `json:"message,required"`
	Object    FineTuningJobEventObject `json:"object,required"`
	JSON      fineTuningJobEventJSON   `json:"-"`
}

// fineTuningJobEventJSON contains the JSON metadata for the struct
// [FineTuningJobEvent]
type fineTuningJobEventJSON struct {
	ID          apijson.Field
	CreatedAt   apijson.Field
	Level       apijson.Field
	Message     apijson.Field
	Object      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *FineTuningJobEvent) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r fineTuningJobEventJSON) RawJSON() string {
	return r.raw
}

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
	// [chat](https://platform.openai.com/docs/api-reference/fine-tuning/chat-input) or
	// [completions](https://platform.openai.com/docs/api-reference/fine-tuning/completions-input)
	// format.
	//
	// See the [fine-tuning guide](https://platform.openai.com/docs/guides/fine-tuning)
	// for more details.
	TrainingFile param.Field[string] `json:"training_file,required"`
	// The hyperparameters used for the fine-tuning job.
	Hyperparameters param.Field[FineTuningJobNewParamsHyperparameters] `json:"hyperparameters"`
	// A list of integrations to enable for your fine-tuning job.
	Integrations param.Field[[]FineTuningJobNewParamsIntegration] `json:"integrations"`
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

// The hyperparameters used for the fine-tuning job.
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
// Satisfied by [FineTuningJobNewParamsHyperparametersBatchSizeBehavior],
// [shared.UnionInt].
type FineTuningJobNewParamsHyperparametersBatchSizeUnion interface {
	ImplementsFineTuningJobNewParamsHyperparametersBatchSizeUnion()
}

type FineTuningJobNewParamsHyperparametersBatchSizeBehavior string

const (
	FineTuningJobNewParamsHyperparametersBatchSizeBehaviorAuto FineTuningJobNewParamsHyperparametersBatchSizeBehavior = "auto"
)

func (r FineTuningJobNewParamsHyperparametersBatchSizeBehavior) IsKnown() bool {
	switch r {
	case FineTuningJobNewParamsHyperparametersBatchSizeBehaviorAuto:
		return true
	}
	return false
}

func (r FineTuningJobNewParamsHyperparametersBatchSizeBehavior) ImplementsFineTuningJobNewParamsHyperparametersBatchSizeUnion() {
}

// Scaling factor for the learning rate. A smaller learning rate may be useful to
// avoid overfitting.
//
// Satisfied by
// [FineTuningJobNewParamsHyperparametersLearningRateMultiplierBehavior],
// [shared.UnionFloat].
type FineTuningJobNewParamsHyperparametersLearningRateMultiplierUnion interface {
	ImplementsFineTuningJobNewParamsHyperparametersLearningRateMultiplierUnion()
}

type FineTuningJobNewParamsHyperparametersLearningRateMultiplierBehavior string

const (
	FineTuningJobNewParamsHyperparametersLearningRateMultiplierBehaviorAuto FineTuningJobNewParamsHyperparametersLearningRateMultiplierBehavior = "auto"
)

func (r FineTuningJobNewParamsHyperparametersLearningRateMultiplierBehavior) IsKnown() bool {
	switch r {
	case FineTuningJobNewParamsHyperparametersLearningRateMultiplierBehaviorAuto:
		return true
	}
	return false
}

func (r FineTuningJobNewParamsHyperparametersLearningRateMultiplierBehavior) ImplementsFineTuningJobNewParamsHyperparametersLearningRateMultiplierUnion() {
}

// The number of epochs to train the model for. An epoch refers to one full cycle
// through the training dataset.
//
// Satisfied by [FineTuningJobNewParamsHyperparametersNEpochsBehavior],
// [shared.UnionInt].
type FineTuningJobNewParamsHyperparametersNEpochsUnion interface {
	ImplementsFineTuningJobNewParamsHyperparametersNEpochsUnion()
}

type FineTuningJobNewParamsHyperparametersNEpochsBehavior string

const (
	FineTuningJobNewParamsHyperparametersNEpochsBehaviorAuto FineTuningJobNewParamsHyperparametersNEpochsBehavior = "auto"
)

func (r FineTuningJobNewParamsHyperparametersNEpochsBehavior) IsKnown() bool {
	switch r {
	case FineTuningJobNewParamsHyperparametersNEpochsBehaviorAuto:
		return true
	}
	return false
}

func (r FineTuningJobNewParamsHyperparametersNEpochsBehavior) ImplementsFineTuningJobNewParamsHyperparametersNEpochsUnion() {
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

type FineTuningJobListParams struct {
	// Identifier for the last job from the previous pagination request.
	After param.Field[string] `query:"after"`
	// Number of fine-tuning jobs to retrieve.
	Limit param.Field[int64] `query:"limit"`
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

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
	"github.com/openai/openai-go/packages/ssestream"
	"github.com/openai/openai-go/shared"
)

// BetaThreadRunService contains methods and other services that help with
// interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaThreadRunService] method instead.
type BetaThreadRunService struct {
	Options []option.RequestOption
	Steps   *BetaThreadRunStepService
}

// NewBetaThreadRunService generates a new service that applies the given options
// to each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewBetaThreadRunService(opts ...option.RequestOption) (r *BetaThreadRunService) {
	r = &BetaThreadRunService{}
	r.Options = opts
	r.Steps = NewBetaThreadRunStepService(opts...)
	return
}

// Create a run.
func (r *BetaThreadRunService) New(ctx context.Context, threadID string, params BetaThreadRunNewParams, opts ...option.RequestOption) (res *Run, err error) {
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "assistants=v2")}, opts...)
	if threadID == "" {
		err = errors.New("missing required thread_id parameter")
		return
	}
	path := fmt.Sprintf("threads/%s/runs", threadID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return
}

// Create a run.
func (r *BetaThreadRunService) NewStreaming(ctx context.Context, threadID string, params BetaThreadRunNewParams, opts ...option.RequestOption) (stream *ssestream.Stream[AssistantStreamEvent]) {
	var (
		raw *http.Response
		err error
	)
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "assistants=v2"), option.WithJSONSet("stream", true)}, opts...)
	if threadID == "" {
		err = errors.New("missing required thread_id parameter")
		return
	}
	path := fmt.Sprintf("threads/%s/runs", threadID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &raw, opts...)
	return ssestream.NewStream[AssistantStreamEvent](ssestream.NewDecoder(raw), err)
}

// Create a run and poll until task is completed.
// Pass 0 to pollIntervalMs to use the default polling interval.
func (r *BetaThreadRunService) NewAndPoll(ctx context.Context, threadID string, params BetaThreadRunNewParams, pollIntervalMs int, opts ...option.RequestOption) (res *Run, err error) {
	run, err := r.New(ctx, threadID, params, opts...)
	if err != nil {
		return nil, err
	}
	return r.PollStatus(ctx, threadID, run.ID, pollIntervalMs, opts...)
}

// Retrieves a run.
func (r *BetaThreadRunService) Get(ctx context.Context, threadID string, runID string, opts ...option.RequestOption) (res *Run, err error) {
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "assistants=v2")}, opts...)
	if threadID == "" {
		err = errors.New("missing required thread_id parameter")
		return
	}
	if runID == "" {
		err = errors.New("missing required run_id parameter")
		return
	}
	path := fmt.Sprintf("threads/%s/runs/%s", threadID, runID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// Modifies a run.
func (r *BetaThreadRunService) Update(ctx context.Context, threadID string, runID string, body BetaThreadRunUpdateParams, opts ...option.RequestOption) (res *Run, err error) {
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "assistants=v2")}, opts...)
	if threadID == "" {
		err = errors.New("missing required thread_id parameter")
		return
	}
	if runID == "" {
		err = errors.New("missing required run_id parameter")
		return
	}
	path := fmt.Sprintf("threads/%s/runs/%s", threadID, runID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Returns a list of runs belonging to a thread.
func (r *BetaThreadRunService) List(ctx context.Context, threadID string, query BetaThreadRunListParams, opts ...option.RequestOption) (res *pagination.CursorPage[Run], err error) {
	var raw *http.Response
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "assistants=v2"), option.WithResponseInto(&raw)}, opts...)
	if threadID == "" {
		err = errors.New("missing required thread_id parameter")
		return
	}
	path := fmt.Sprintf("threads/%s/runs", threadID)
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

// Returns a list of runs belonging to a thread.
func (r *BetaThreadRunService) ListAutoPaging(ctx context.Context, threadID string, query BetaThreadRunListParams, opts ...option.RequestOption) *pagination.CursorPageAutoPager[Run] {
	return pagination.NewCursorPageAutoPager(r.List(ctx, threadID, query, opts...))
}

// Cancels a run that is `in_progress`.
func (r *BetaThreadRunService) Cancel(ctx context.Context, threadID string, runID string, opts ...option.RequestOption) (res *Run, err error) {
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "assistants=v2")}, opts...)
	if threadID == "" {
		err = errors.New("missing required thread_id parameter")
		return
	}
	if runID == "" {
		err = errors.New("missing required run_id parameter")
		return
	}
	path := fmt.Sprintf("threads/%s/runs/%s/cancel", threadID, runID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return
}

// When a run has the `status: "requires_action"` and `required_action.type` is
// `submit_tool_outputs`, this endpoint can be used to submit the outputs from the
// tool calls once they're all completed. All outputs must be submitted in a single
// request.
func (r *BetaThreadRunService) SubmitToolOutputs(ctx context.Context, threadID string, runID string, body BetaThreadRunSubmitToolOutputsParams, opts ...option.RequestOption) (res *Run, err error) {
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "assistants=v2")}, opts...)
	if threadID == "" {
		err = errors.New("missing required thread_id parameter")
		return
	}
	if runID == "" {
		err = errors.New("missing required run_id parameter")
		return
	}
	path := fmt.Sprintf("threads/%s/runs/%s/submit_tool_outputs", threadID, runID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// A helper to submit a tool output to a run and poll for a terminal run state.
// Pass 0 to pollIntervalMs to use the default polling interval.
// More information on Run lifecycles can be found here:
// https://platform.openai.com/docs/assistants/how-it-works/runs-and-run-steps
func (r *BetaThreadRunService) SubmitToolOutputsAndPoll(ctx context.Context, threadID string, runID string, body BetaThreadRunSubmitToolOutputsParams, pollIntervalMs int, opts ...option.RequestOption) (*Run, error) {
	run, err := r.SubmitToolOutputs(ctx, threadID, runID, body, opts...)
	if err != nil {
		return nil, err
	}
	return r.PollStatus(ctx, threadID, run.ID, pollIntervalMs, opts...)
}

// When a run has the `status: "requires_action"` and `required_action.type` is
// `submit_tool_outputs`, this endpoint can be used to submit the outputs from the
// tool calls once they're all completed. All outputs must be submitted in a single
// request.
func (r *BetaThreadRunService) SubmitToolOutputsStreaming(ctx context.Context, threadID string, runID string, body BetaThreadRunSubmitToolOutputsParams, opts ...option.RequestOption) (stream *ssestream.Stream[AssistantStreamEvent]) {
	var (
		raw *http.Response
		err error
	)
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "assistants=v2"), option.WithJSONSet("stream", true)}, opts...)
	if threadID == "" {
		err = errors.New("missing required thread_id parameter")
		return
	}
	if runID == "" {
		err = errors.New("missing required run_id parameter")
		return
	}
	path := fmt.Sprintf("threads/%s/runs/%s/submit_tool_outputs", threadID, runID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &raw, opts...)
	return ssestream.NewStream[AssistantStreamEvent](ssestream.NewDecoder(raw), err)
}

// Tool call objects
type RequiredActionFunctionToolCall struct {
	// The ID of the tool call. This ID must be referenced when you submit the tool
	// outputs in using the
	// [Submit tool outputs to run](https://platform.openai.com/docs/api-reference/runs/submitToolOutputs)
	// endpoint.
	ID string `json:"id,required"`
	// The function definition.
	Function RequiredActionFunctionToolCallFunction `json:"function,required"`
	// The type of tool call the output is required for. For now, this is always
	// `function`.
	Type RequiredActionFunctionToolCallType `json:"type,required"`
	JSON requiredActionFunctionToolCallJSON `json:"-"`
}

// requiredActionFunctionToolCallJSON contains the JSON metadata for the struct
// [RequiredActionFunctionToolCall]
type requiredActionFunctionToolCallJSON struct {
	ID          apijson.Field
	Function    apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *RequiredActionFunctionToolCall) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r requiredActionFunctionToolCallJSON) RawJSON() string {
	return r.raw
}

// The function definition.
type RequiredActionFunctionToolCallFunction struct {
	// The arguments that the model expects you to pass to the function.
	Arguments string `json:"arguments,required"`
	// The name of the function.
	Name string                                     `json:"name,required"`
	JSON requiredActionFunctionToolCallFunctionJSON `json:"-"`
}

// requiredActionFunctionToolCallFunctionJSON contains the JSON metadata for the
// struct [RequiredActionFunctionToolCallFunction]
type requiredActionFunctionToolCallFunctionJSON struct {
	Arguments   apijson.Field
	Name        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *RequiredActionFunctionToolCallFunction) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r requiredActionFunctionToolCallFunctionJSON) RawJSON() string {
	return r.raw
}

// The type of tool call the output is required for. For now, this is always
// `function`.
type RequiredActionFunctionToolCallType string

const (
	RequiredActionFunctionToolCallTypeFunction RequiredActionFunctionToolCallType = "function"
)

func (r RequiredActionFunctionToolCallType) IsKnown() bool {
	switch r {
	case RequiredActionFunctionToolCallTypeFunction:
		return true
	}
	return false
}

// Represents an execution run on a
// [thread](https://platform.openai.com/docs/api-reference/threads).
type Run struct {
	// The identifier, which can be referenced in API endpoints.
	ID string `json:"id,required"`
	// The ID of the
	// [assistant](https://platform.openai.com/docs/api-reference/assistants) used for
	// execution of this run.
	AssistantID string `json:"assistant_id,required"`
	// The Unix timestamp (in seconds) for when the run was cancelled.
	CancelledAt int64 `json:"cancelled_at,required,nullable"`
	// The Unix timestamp (in seconds) for when the run was completed.
	CompletedAt int64 `json:"completed_at,required,nullable"`
	// The Unix timestamp (in seconds) for when the run was created.
	CreatedAt int64 `json:"created_at,required"`
	// The Unix timestamp (in seconds) for when the run will expire.
	ExpiresAt int64 `json:"expires_at,required,nullable"`
	// The Unix timestamp (in seconds) for when the run failed.
	FailedAt int64 `json:"failed_at,required,nullable"`
	// Details on why the run is incomplete. Will be `null` if the run is not
	// incomplete.
	IncompleteDetails RunIncompleteDetails `json:"incomplete_details,required,nullable"`
	// The instructions that the
	// [assistant](https://platform.openai.com/docs/api-reference/assistants) used for
	// this run.
	Instructions string `json:"instructions,required"`
	// The last error associated with this run. Will be `null` if there are no errors.
	LastError RunLastError `json:"last_error,required,nullable"`
	// The maximum number of completion tokens specified to have been used over the
	// course of the run.
	MaxCompletionTokens int64 `json:"max_completion_tokens,required,nullable"`
	// The maximum number of prompt tokens specified to have been used over the course
	// of the run.
	MaxPromptTokens int64 `json:"max_prompt_tokens,required,nullable"`
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard.
	//
	// Keys are strings with a maximum length of 64 characters. Values are strings with
	// a maximum length of 512 characters.
	Metadata shared.Metadata `json:"metadata,required,nullable"`
	// The model that the
	// [assistant](https://platform.openai.com/docs/api-reference/assistants) used for
	// this run.
	Model string `json:"model,required"`
	// The object type, which is always `thread.run`.
	Object RunObject `json:"object,required"`
	// Whether to enable
	// [parallel function calling](https://platform.openai.com/docs/guides/function-calling#configuring-parallel-function-calling)
	// during tool use.
	ParallelToolCalls bool `json:"parallel_tool_calls,required"`
	// Details on the action required to continue the run. Will be `null` if no action
	// is required.
	RequiredAction RunRequiredAction `json:"required_action,required,nullable"`
	// The Unix timestamp (in seconds) for when the run was started.
	StartedAt int64 `json:"started_at,required,nullable"`
	// The status of the run, which can be either `queued`, `in_progress`,
	// `requires_action`, `cancelling`, `cancelled`, `failed`, `completed`,
	// `incomplete`, or `expired`.
	Status RunStatus `json:"status,required"`
	// The ID of the [thread](https://platform.openai.com/docs/api-reference/threads)
	// that was executed on as a part of this run.
	ThreadID string `json:"thread_id,required"`
	// Controls which (if any) tool is called by the model. `none` means the model will
	// not call any tools and instead generates a message. `auto` is the default value
	// and means the model can pick between generating a message or calling one or more
	// tools. `required` means the model must call one or more tools before responding
	// to the user. Specifying a particular tool like `{"type": "file_search"}` or
	// `{"type": "function", "function": {"name": "my_function"}}` forces the model to
	// call that tool.
	ToolChoice AssistantToolChoiceOptionUnion `json:"tool_choice,required,nullable"`
	// The list of tools that the
	// [assistant](https://platform.openai.com/docs/api-reference/assistants) used for
	// this run.
	Tools []AssistantTool `json:"tools,required"`
	// Controls for how a thread will be truncated prior to the run. Use this to
	// control the intial context window of the run.
	TruncationStrategy RunTruncationStrategy `json:"truncation_strategy,required,nullable"`
	// Usage statistics related to the run. This value will be `null` if the run is not
	// in a terminal state (i.e. `in_progress`, `queued`, etc.).
	Usage RunUsage `json:"usage,required,nullable"`
	// The sampling temperature used for this run. If not set, defaults to 1.
	Temperature float64 `json:"temperature,nullable"`
	// The nucleus sampling value used for this run. If not set, defaults to 1.
	TopP float64 `json:"top_p,nullable"`
	JSON runJSON `json:"-"`
}

// runJSON contains the JSON metadata for the struct [Run]
type runJSON struct {
	ID                  apijson.Field
	AssistantID         apijson.Field
	CancelledAt         apijson.Field
	CompletedAt         apijson.Field
	CreatedAt           apijson.Field
	ExpiresAt           apijson.Field
	FailedAt            apijson.Field
	IncompleteDetails   apijson.Field
	Instructions        apijson.Field
	LastError           apijson.Field
	MaxCompletionTokens apijson.Field
	MaxPromptTokens     apijson.Field
	Metadata            apijson.Field
	Model               apijson.Field
	Object              apijson.Field
	ParallelToolCalls   apijson.Field
	RequiredAction      apijson.Field
	StartedAt           apijson.Field
	Status              apijson.Field
	ThreadID            apijson.Field
	ToolChoice          apijson.Field
	Tools               apijson.Field
	TruncationStrategy  apijson.Field
	Usage               apijson.Field
	Temperature         apijson.Field
	TopP                apijson.Field
	raw                 string
	ExtraFields         map[string]apijson.Field
}

func (r *Run) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r runJSON) RawJSON() string {
	return r.raw
}

// Details on why the run is incomplete. Will be `null` if the run is not
// incomplete.
type RunIncompleteDetails struct {
	// The reason why the run is incomplete. This will point to which specific token
	// limit was reached over the course of the run.
	Reason RunIncompleteDetailsReason `json:"reason"`
	JSON   runIncompleteDetailsJSON   `json:"-"`
}

// runIncompleteDetailsJSON contains the JSON metadata for the struct
// [RunIncompleteDetails]
type runIncompleteDetailsJSON struct {
	Reason      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *RunIncompleteDetails) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r runIncompleteDetailsJSON) RawJSON() string {
	return r.raw
}

// The reason why the run is incomplete. This will point to which specific token
// limit was reached over the course of the run.
type RunIncompleteDetailsReason string

const (
	RunIncompleteDetailsReasonMaxCompletionTokens RunIncompleteDetailsReason = "max_completion_tokens"
	RunIncompleteDetailsReasonMaxPromptTokens     RunIncompleteDetailsReason = "max_prompt_tokens"
)

func (r RunIncompleteDetailsReason) IsKnown() bool {
	switch r {
	case RunIncompleteDetailsReasonMaxCompletionTokens, RunIncompleteDetailsReasonMaxPromptTokens:
		return true
	}
	return false
}

// The last error associated with this run. Will be `null` if there are no errors.
type RunLastError struct {
	// One of `server_error`, `rate_limit_exceeded`, or `invalid_prompt`.
	Code RunLastErrorCode `json:"code,required"`
	// A human-readable description of the error.
	Message string           `json:"message,required"`
	JSON    runLastErrorJSON `json:"-"`
}

// runLastErrorJSON contains the JSON metadata for the struct [RunLastError]
type runLastErrorJSON struct {
	Code        apijson.Field
	Message     apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *RunLastError) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r runLastErrorJSON) RawJSON() string {
	return r.raw
}

// One of `server_error`, `rate_limit_exceeded`, or `invalid_prompt`.
type RunLastErrorCode string

const (
	RunLastErrorCodeServerError       RunLastErrorCode = "server_error"
	RunLastErrorCodeRateLimitExceeded RunLastErrorCode = "rate_limit_exceeded"
	RunLastErrorCodeInvalidPrompt     RunLastErrorCode = "invalid_prompt"
)

func (r RunLastErrorCode) IsKnown() bool {
	switch r {
	case RunLastErrorCodeServerError, RunLastErrorCodeRateLimitExceeded, RunLastErrorCodeInvalidPrompt:
		return true
	}
	return false
}

// The object type, which is always `thread.run`.
type RunObject string

const (
	RunObjectThreadRun RunObject = "thread.run"
)

func (r RunObject) IsKnown() bool {
	switch r {
	case RunObjectThreadRun:
		return true
	}
	return false
}

// Details on the action required to continue the run. Will be `null` if no action
// is required.
type RunRequiredAction struct {
	// Details on the tool outputs needed for this run to continue.
	SubmitToolOutputs RunRequiredActionSubmitToolOutputs `json:"submit_tool_outputs,required"`
	// For now, this is always `submit_tool_outputs`.
	Type RunRequiredActionType `json:"type,required"`
	JSON runRequiredActionJSON `json:"-"`
}

// runRequiredActionJSON contains the JSON metadata for the struct
// [RunRequiredAction]
type runRequiredActionJSON struct {
	SubmitToolOutputs apijson.Field
	Type              apijson.Field
	raw               string
	ExtraFields       map[string]apijson.Field
}

func (r *RunRequiredAction) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r runRequiredActionJSON) RawJSON() string {
	return r.raw
}

// Details on the tool outputs needed for this run to continue.
type RunRequiredActionSubmitToolOutputs struct {
	// A list of the relevant tool calls.
	ToolCalls []RequiredActionFunctionToolCall       `json:"tool_calls,required"`
	JSON      runRequiredActionSubmitToolOutputsJSON `json:"-"`
}

// runRequiredActionSubmitToolOutputsJSON contains the JSON metadata for the struct
// [RunRequiredActionSubmitToolOutputs]
type runRequiredActionSubmitToolOutputsJSON struct {
	ToolCalls   apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *RunRequiredActionSubmitToolOutputs) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r runRequiredActionSubmitToolOutputsJSON) RawJSON() string {
	return r.raw
}

// For now, this is always `submit_tool_outputs`.
type RunRequiredActionType string

const (
	RunRequiredActionTypeSubmitToolOutputs RunRequiredActionType = "submit_tool_outputs"
)

func (r RunRequiredActionType) IsKnown() bool {
	switch r {
	case RunRequiredActionTypeSubmitToolOutputs:
		return true
	}
	return false
}

// Controls for how a thread will be truncated prior to the run. Use this to
// control the intial context window of the run.
type RunTruncationStrategy struct {
	// The truncation strategy to use for the thread. The default is `auto`. If set to
	// `last_messages`, the thread will be truncated to the n most recent messages in
	// the thread. When set to `auto`, messages in the middle of the thread will be
	// dropped to fit the context length of the model, `max_prompt_tokens`.
	Type RunTruncationStrategyType `json:"type,required"`
	// The number of most recent messages from the thread when constructing the context
	// for the run.
	LastMessages int64                     `json:"last_messages,nullable"`
	JSON         runTruncationStrategyJSON `json:"-"`
}

// runTruncationStrategyJSON contains the JSON metadata for the struct
// [RunTruncationStrategy]
type runTruncationStrategyJSON struct {
	Type         apijson.Field
	LastMessages apijson.Field
	raw          string
	ExtraFields  map[string]apijson.Field
}

func (r *RunTruncationStrategy) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r runTruncationStrategyJSON) RawJSON() string {
	return r.raw
}

// The truncation strategy to use for the thread. The default is `auto`. If set to
// `last_messages`, the thread will be truncated to the n most recent messages in
// the thread. When set to `auto`, messages in the middle of the thread will be
// dropped to fit the context length of the model, `max_prompt_tokens`.
type RunTruncationStrategyType string

const (
	RunTruncationStrategyTypeAuto         RunTruncationStrategyType = "auto"
	RunTruncationStrategyTypeLastMessages RunTruncationStrategyType = "last_messages"
)

func (r RunTruncationStrategyType) IsKnown() bool {
	switch r {
	case RunTruncationStrategyTypeAuto, RunTruncationStrategyTypeLastMessages:
		return true
	}
	return false
}

// Usage statistics related to the run. This value will be `null` if the run is not
// in a terminal state (i.e. `in_progress`, `queued`, etc.).
type RunUsage struct {
	// Number of completion tokens used over the course of the run.
	CompletionTokens int64 `json:"completion_tokens,required"`
	// Number of prompt tokens used over the course of the run.
	PromptTokens int64 `json:"prompt_tokens,required"`
	// Total number of tokens used (prompt + completion).
	TotalTokens int64        `json:"total_tokens,required"`
	JSON        runUsageJSON `json:"-"`
}

// runUsageJSON contains the JSON metadata for the struct [RunUsage]
type runUsageJSON struct {
	CompletionTokens apijson.Field
	PromptTokens     apijson.Field
	TotalTokens      apijson.Field
	raw              string
	ExtraFields      map[string]apijson.Field
}

func (r *RunUsage) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r runUsageJSON) RawJSON() string {
	return r.raw
}

// The status of the run, which can be either `queued`, `in_progress`,
// `requires_action`, `cancelling`, `cancelled`, `failed`, `completed`,
// `incomplete`, or `expired`.
type RunStatus string

const (
	RunStatusQueued         RunStatus = "queued"
	RunStatusInProgress     RunStatus = "in_progress"
	RunStatusRequiresAction RunStatus = "requires_action"
	RunStatusCancelling     RunStatus = "cancelling"
	RunStatusCancelled      RunStatus = "cancelled"
	RunStatusFailed         RunStatus = "failed"
	RunStatusCompleted      RunStatus = "completed"
	RunStatusIncomplete     RunStatus = "incomplete"
	RunStatusExpired        RunStatus = "expired"
)

func (r RunStatus) IsKnown() bool {
	switch r {
	case RunStatusQueued, RunStatusInProgress, RunStatusRequiresAction, RunStatusCancelling, RunStatusCancelled, RunStatusFailed, RunStatusCompleted, RunStatusIncomplete, RunStatusExpired:
		return true
	}
	return false
}

type BetaThreadRunNewParams struct {
	// The ID of the
	// [assistant](https://platform.openai.com/docs/api-reference/assistants) to use to
	// execute this run.
	AssistantID param.Field[string] `json:"assistant_id,required"`
	// A list of additional fields to include in the response. Currently the only
	// supported value is `step_details.tool_calls[*].file_search.results[*].content`
	// to fetch the file search result content.
	//
	// See the
	// [file search tool documentation](https://platform.openai.com/docs/assistants/tools/file-search#customizing-file-search-settings)
	// for more information.
	Include param.Field[[]RunStepInclude] `query:"include"`
	// Appends additional instructions at the end of the instructions for the run. This
	// is useful for modifying the behavior on a per-run basis without overriding other
	// instructions.
	AdditionalInstructions param.Field[string] `json:"additional_instructions"`
	// Adds additional messages to the thread before creating the run.
	AdditionalMessages param.Field[[]BetaThreadRunNewParamsAdditionalMessage] `json:"additional_messages"`
	// Overrides the
	// [instructions](https://platform.openai.com/docs/api-reference/assistants/createAssistant)
	// of the assistant. This is useful for modifying the behavior on a per-run basis.
	Instructions param.Field[string] `json:"instructions"`
	// The maximum number of completion tokens that may be used over the course of the
	// run. The run will make a best effort to use only the number of completion tokens
	// specified, across multiple turns of the run. If the run exceeds the number of
	// completion tokens specified, the run will end with status `incomplete`. See
	// `incomplete_details` for more info.
	MaxCompletionTokens param.Field[int64] `json:"max_completion_tokens"`
	// The maximum number of prompt tokens that may be used over the course of the run.
	// The run will make a best effort to use only the number of prompt tokens
	// specified, across multiple turns of the run. If the run exceeds the number of
	// prompt tokens specified, the run will end with status `incomplete`. See
	// `incomplete_details` for more info.
	MaxPromptTokens param.Field[int64] `json:"max_prompt_tokens"`
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard.
	//
	// Keys are strings with a maximum length of 64 characters. Values are strings with
	// a maximum length of 512 characters.
	Metadata param.Field[shared.MetadataParam] `json:"metadata"`
	// The ID of the [Model](https://platform.openai.com/docs/api-reference/models) to
	// be used to execute this run. If a value is provided here, it will override the
	// model associated with the assistant. If not, the model associated with the
	// assistant will be used.
	Model param.Field[shared.ChatModel] `json:"model"`
	// Whether to enable
	// [parallel function calling](https://platform.openai.com/docs/guides/function-calling#configuring-parallel-function-calling)
	// during tool use.
	ParallelToolCalls param.Field[bool] `json:"parallel_tool_calls"`
	// **o1 and o3-mini models only**
	//
	// Constrains effort on reasoning for
	// [reasoning models](https://platform.openai.com/docs/guides/reasoning). Currently
	// supported values are `low`, `medium`, and `high`. Reducing reasoning effort can
	// result in faster responses and fewer tokens used on reasoning in a response.
	ReasoningEffort param.Field[BetaThreadRunNewParamsReasoningEffort] `json:"reasoning_effort"`
	// What sampling temperature to use, between 0 and 2. Higher values like 0.8 will
	// make the output more random, while lower values like 0.2 will make it more
	// focused and deterministic.
	Temperature param.Field[float64] `json:"temperature"`
	// Controls which (if any) tool is called by the model. `none` means the model will
	// not call any tools and instead generates a message. `auto` is the default value
	// and means the model can pick between generating a message or calling one or more
	// tools. `required` means the model must call one or more tools before responding
	// to the user. Specifying a particular tool like `{"type": "file_search"}` or
	// `{"type": "function", "function": {"name": "my_function"}}` forces the model to
	// call that tool.
	ToolChoice param.Field[AssistantToolChoiceOptionUnionParam] `json:"tool_choice"`
	// Override the tools the assistant can use for this run. This is useful for
	// modifying the behavior on a per-run basis.
	Tools param.Field[[]AssistantToolUnionParam] `json:"tools"`
	// An alternative to sampling with temperature, called nucleus sampling, where the
	// model considers the results of the tokens with top_p probability mass. So 0.1
	// means only the tokens comprising the top 10% probability mass are considered.
	//
	// We generally recommend altering this or temperature but not both.
	TopP param.Field[float64] `json:"top_p"`
	// Controls for how a thread will be truncated prior to the run. Use this to
	// control the intial context window of the run.
	TruncationStrategy param.Field[BetaThreadRunNewParamsTruncationStrategy] `json:"truncation_strategy"`
}

func (r BetaThreadRunNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// URLQuery serializes [BetaThreadRunNewParams]'s query parameters as `url.Values`.
func (r BetaThreadRunNewParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BetaThreadRunNewParamsAdditionalMessage struct {
	// An array of content parts with a defined type, each can be of type `text` or
	// images can be passed with `image_url` or `image_file`. Image types are only
	// supported on
	// [Vision-compatible models](https://platform.openai.com/docs/models).
	Content param.Field[[]MessageContentPartParamUnion] `json:"content,required"`
	// The role of the entity that is creating the message. Allowed values include:
	//
	//   - `user`: Indicates the message is sent by an actual user and should be used in
	//     most cases to represent user-generated messages.
	//   - `assistant`: Indicates the message is generated by the assistant. Use this
	//     value to insert messages from the assistant into the conversation.
	Role param.Field[BetaThreadRunNewParamsAdditionalMessagesRole] `json:"role,required"`
	// A list of files attached to the message, and the tools they should be added to.
	Attachments param.Field[[]BetaThreadRunNewParamsAdditionalMessagesAttachment] `json:"attachments"`
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard.
	//
	// Keys are strings with a maximum length of 64 characters. Values are strings with
	// a maximum length of 512 characters.
	Metadata param.Field[shared.MetadataParam] `json:"metadata"`
}

func (r BetaThreadRunNewParamsAdditionalMessage) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// The role of the entity that is creating the message. Allowed values include:
//
//   - `user`: Indicates the message is sent by an actual user and should be used in
//     most cases to represent user-generated messages.
//   - `assistant`: Indicates the message is generated by the assistant. Use this
//     value to insert messages from the assistant into the conversation.
type BetaThreadRunNewParamsAdditionalMessagesRole string

const (
	BetaThreadRunNewParamsAdditionalMessagesRoleUser      BetaThreadRunNewParamsAdditionalMessagesRole = "user"
	BetaThreadRunNewParamsAdditionalMessagesRoleAssistant BetaThreadRunNewParamsAdditionalMessagesRole = "assistant"
)

func (r BetaThreadRunNewParamsAdditionalMessagesRole) IsKnown() bool {
	switch r {
	case BetaThreadRunNewParamsAdditionalMessagesRoleUser, BetaThreadRunNewParamsAdditionalMessagesRoleAssistant:
		return true
	}
	return false
}

type BetaThreadRunNewParamsAdditionalMessagesAttachment struct {
	// The ID of the file to attach to the message.
	FileID param.Field[string] `json:"file_id"`
	// The tools to add this file to.
	Tools param.Field[[]BetaThreadRunNewParamsAdditionalMessagesAttachmentsToolUnion] `json:"tools"`
}

func (r BetaThreadRunNewParamsAdditionalMessagesAttachment) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type BetaThreadRunNewParamsAdditionalMessagesAttachmentsTool struct {
	// The type of tool being defined: `code_interpreter`
	Type param.Field[BetaThreadRunNewParamsAdditionalMessagesAttachmentsToolsType] `json:"type,required"`
}

func (r BetaThreadRunNewParamsAdditionalMessagesAttachmentsTool) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r BetaThreadRunNewParamsAdditionalMessagesAttachmentsTool) implementsBetaThreadRunNewParamsAdditionalMessagesAttachmentsToolUnion() {
}

// Satisfied by [CodeInterpreterToolParam],
// [BetaThreadRunNewParamsAdditionalMessagesAttachmentsToolsFileSearch],
// [BetaThreadRunNewParamsAdditionalMessagesAttachmentsTool].
type BetaThreadRunNewParamsAdditionalMessagesAttachmentsToolUnion interface {
	implementsBetaThreadRunNewParamsAdditionalMessagesAttachmentsToolUnion()
}

type BetaThreadRunNewParamsAdditionalMessagesAttachmentsToolsFileSearch struct {
	// The type of tool being defined: `file_search`
	Type param.Field[BetaThreadRunNewParamsAdditionalMessagesAttachmentsToolsFileSearchType] `json:"type,required"`
}

func (r BetaThreadRunNewParamsAdditionalMessagesAttachmentsToolsFileSearch) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r BetaThreadRunNewParamsAdditionalMessagesAttachmentsToolsFileSearch) implementsBetaThreadRunNewParamsAdditionalMessagesAttachmentsToolUnion() {
}

// The type of tool being defined: `file_search`
type BetaThreadRunNewParamsAdditionalMessagesAttachmentsToolsFileSearchType string

const (
	BetaThreadRunNewParamsAdditionalMessagesAttachmentsToolsFileSearchTypeFileSearch BetaThreadRunNewParamsAdditionalMessagesAttachmentsToolsFileSearchType = "file_search"
)

func (r BetaThreadRunNewParamsAdditionalMessagesAttachmentsToolsFileSearchType) IsKnown() bool {
	switch r {
	case BetaThreadRunNewParamsAdditionalMessagesAttachmentsToolsFileSearchTypeFileSearch:
		return true
	}
	return false
}

// The type of tool being defined: `code_interpreter`
type BetaThreadRunNewParamsAdditionalMessagesAttachmentsToolsType string

const (
	BetaThreadRunNewParamsAdditionalMessagesAttachmentsToolsTypeCodeInterpreter BetaThreadRunNewParamsAdditionalMessagesAttachmentsToolsType = "code_interpreter"
	BetaThreadRunNewParamsAdditionalMessagesAttachmentsToolsTypeFileSearch      BetaThreadRunNewParamsAdditionalMessagesAttachmentsToolsType = "file_search"
)

func (r BetaThreadRunNewParamsAdditionalMessagesAttachmentsToolsType) IsKnown() bool {
	switch r {
	case BetaThreadRunNewParamsAdditionalMessagesAttachmentsToolsTypeCodeInterpreter, BetaThreadRunNewParamsAdditionalMessagesAttachmentsToolsTypeFileSearch:
		return true
	}
	return false
}

// **o1 and o3-mini models only**
//
// Constrains effort on reasoning for
// [reasoning models](https://platform.openai.com/docs/guides/reasoning). Currently
// supported values are `low`, `medium`, and `high`. Reducing reasoning effort can
// result in faster responses and fewer tokens used on reasoning in a response.
type BetaThreadRunNewParamsReasoningEffort string

const (
	BetaThreadRunNewParamsReasoningEffortLow    BetaThreadRunNewParamsReasoningEffort = "low"
	BetaThreadRunNewParamsReasoningEffortMedium BetaThreadRunNewParamsReasoningEffort = "medium"
	BetaThreadRunNewParamsReasoningEffortHigh   BetaThreadRunNewParamsReasoningEffort = "high"
)

func (r BetaThreadRunNewParamsReasoningEffort) IsKnown() bool {
	switch r {
	case BetaThreadRunNewParamsReasoningEffortLow, BetaThreadRunNewParamsReasoningEffortMedium, BetaThreadRunNewParamsReasoningEffortHigh:
		return true
	}
	return false
}

// Controls for how a thread will be truncated prior to the run. Use this to
// control the intial context window of the run.
type BetaThreadRunNewParamsTruncationStrategy struct {
	// The truncation strategy to use for the thread. The default is `auto`. If set to
	// `last_messages`, the thread will be truncated to the n most recent messages in
	// the thread. When set to `auto`, messages in the middle of the thread will be
	// dropped to fit the context length of the model, `max_prompt_tokens`.
	Type param.Field[BetaThreadRunNewParamsTruncationStrategyType] `json:"type,required"`
	// The number of most recent messages from the thread when constructing the context
	// for the run.
	LastMessages param.Field[int64] `json:"last_messages"`
}

func (r BetaThreadRunNewParamsTruncationStrategy) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// The truncation strategy to use for the thread. The default is `auto`. If set to
// `last_messages`, the thread will be truncated to the n most recent messages in
// the thread. When set to `auto`, messages in the middle of the thread will be
// dropped to fit the context length of the model, `max_prompt_tokens`.
type BetaThreadRunNewParamsTruncationStrategyType string

const (
	BetaThreadRunNewParamsTruncationStrategyTypeAuto         BetaThreadRunNewParamsTruncationStrategyType = "auto"
	BetaThreadRunNewParamsTruncationStrategyTypeLastMessages BetaThreadRunNewParamsTruncationStrategyType = "last_messages"
)

func (r BetaThreadRunNewParamsTruncationStrategyType) IsKnown() bool {
	switch r {
	case BetaThreadRunNewParamsTruncationStrategyTypeAuto, BetaThreadRunNewParamsTruncationStrategyTypeLastMessages:
		return true
	}
	return false
}

type BetaThreadRunUpdateParams struct {
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard.
	//
	// Keys are strings with a maximum length of 64 characters. Values are strings with
	// a maximum length of 512 characters.
	Metadata param.Field[shared.MetadataParam] `json:"metadata"`
}

func (r BetaThreadRunUpdateParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type BetaThreadRunListParams struct {
	// A cursor for use in pagination. `after` is an object ID that defines your place
	// in the list. For instance, if you make a list request and receive 100 objects,
	// ending with obj_foo, your subsequent call can include after=obj_foo in order to
	// fetch the next page of the list.
	After param.Field[string] `query:"after"`
	// A cursor for use in pagination. `before` is an object ID that defines your place
	// in the list. For instance, if you make a list request and receive 100 objects,
	// starting with obj_foo, your subsequent call can include before=obj_foo in order
	// to fetch the previous page of the list.
	Before param.Field[string] `query:"before"`
	// A limit on the number of objects to be returned. Limit can range between 1 and
	// 100, and the default is 20.
	Limit param.Field[int64] `query:"limit"`
	// Sort order by the `created_at` timestamp of the objects. `asc` for ascending
	// order and `desc` for descending order.
	Order param.Field[BetaThreadRunListParamsOrder] `query:"order"`
}

// URLQuery serializes [BetaThreadRunListParams]'s query parameters as
// `url.Values`.
func (r BetaThreadRunListParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Sort order by the `created_at` timestamp of the objects. `asc` for ascending
// order and `desc` for descending order.
type BetaThreadRunListParamsOrder string

const (
	BetaThreadRunListParamsOrderAsc  BetaThreadRunListParamsOrder = "asc"
	BetaThreadRunListParamsOrderDesc BetaThreadRunListParamsOrder = "desc"
)

func (r BetaThreadRunListParamsOrder) IsKnown() bool {
	switch r {
	case BetaThreadRunListParamsOrderAsc, BetaThreadRunListParamsOrderDesc:
		return true
	}
	return false
}

type BetaThreadRunSubmitToolOutputsParams struct {
	// A list of tools for which the outputs are being submitted.
	ToolOutputs param.Field[[]BetaThreadRunSubmitToolOutputsParamsToolOutput] `json:"tool_outputs,required"`
}

func (r BetaThreadRunSubmitToolOutputsParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type BetaThreadRunSubmitToolOutputsParamsToolOutput struct {
	// The output of the tool call to be submitted to continue the run.
	Output param.Field[string] `json:"output"`
	// The ID of the tool call in the `required_action` object within the run object
	// the output is being submitted for.
	ToolCallID param.Field[string] `json:"tool_call_id"`
}

func (r BetaThreadRunSubmitToolOutputsParamsToolOutput) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

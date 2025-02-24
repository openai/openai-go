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

// BetaThreadRunStepService contains methods and other services that help with
// interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaThreadRunStepService] method instead.
type BetaThreadRunStepService struct {
	Options []option.RequestOption
}

// NewBetaThreadRunStepService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewBetaThreadRunStepService(opts ...option.RequestOption) (r BetaThreadRunStepService) {
	r = BetaThreadRunStepService{}
	r.Options = opts
	return
}

// Retrieves a run step.
func (r *BetaThreadRunStepService) Get(ctx context.Context, threadID string, runID string, stepID string, query BetaThreadRunStepGetParams, opts ...option.RequestOption) (res *RunStep, err error) {
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
	if stepID == "" {
		err = errors.New("missing required step_id parameter")
		return
	}
	path := fmt.Sprintf("threads/%s/runs/%s/steps/%s", threadID, runID, stepID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return
}

// Returns a list of run steps belonging to a run.
func (r *BetaThreadRunStepService) List(ctx context.Context, threadID string, runID string, query BetaThreadRunStepListParams, opts ...option.RequestOption) (res *pagination.CursorPage[RunStep], err error) {
	var raw *http.Response
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "assistants=v2"), option.WithResponseInto(&raw)}, opts...)
	if threadID == "" {
		err = errors.New("missing required thread_id parameter")
		return
	}
	if runID == "" {
		err = errors.New("missing required run_id parameter")
		return
	}
	path := fmt.Sprintf("threads/%s/runs/%s/steps", threadID, runID)
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

// Returns a list of run steps belonging to a run.
func (r *BetaThreadRunStepService) ListAutoPaging(ctx context.Context, threadID string, runID string, query BetaThreadRunStepListParams, opts ...option.RequestOption) *pagination.CursorPageAutoPager[RunStep] {
	return pagination.NewCursorPageAutoPager(r.List(ctx, threadID, runID, query, opts...))
}

// Text output from the Code Interpreter tool call as part of a run step.
type CodeInterpreterLogs struct {
	// The index of the output in the outputs array.
	Index int64 `json:"index,omitzero,required"`
	// Always `logs`.
	//
	// This field can be elided, and will be automatically set as "logs".
	Type constant.Logs `json:"type,required"`
	// The text output from the Code Interpreter tool call.
	Logs string `json:"logs,omitzero"`
	JSON struct {
		Index resp.Field
		Type  resp.Field
		Logs  resp.Field
		raw   string
	} `json:"-"`
}

func (r CodeInterpreterLogs) RawJSON() string { return r.JSON.raw }
func (r *CodeInterpreterLogs) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type CodeInterpreterOutputImage struct {
	// The index of the output in the outputs array.
	Index int64 `json:"index,omitzero,required"`
	// Always `image`.
	//
	// This field can be elided, and will be automatically set as "image".
	Type  constant.Image                  `json:"type,required"`
	Image CodeInterpreterOutputImageImage `json:"image,omitzero"`
	JSON  struct {
		Index resp.Field
		Type  resp.Field
		Image resp.Field
		raw   string
	} `json:"-"`
}

func (r CodeInterpreterOutputImage) RawJSON() string { return r.JSON.raw }
func (r *CodeInterpreterOutputImage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type CodeInterpreterOutputImageImage struct {
	// The [file](https://platform.openai.com/docs/api-reference/files) ID of the
	// image.
	FileID string `json:"file_id,omitzero"`
	JSON   struct {
		FileID resp.Field
		raw    string
	} `json:"-"`
}

func (r CodeInterpreterOutputImageImage) RawJSON() string { return r.JSON.raw }
func (r *CodeInterpreterOutputImageImage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Details of the Code Interpreter tool call the run step was involved in.
type CodeInterpreterToolCall struct {
	// The ID of the tool call.
	ID string `json:"id,omitzero,required"`
	// The Code Interpreter tool call definition.
	CodeInterpreter CodeInterpreterToolCallCodeInterpreter `json:"code_interpreter,omitzero,required"`
	// The type of tool call. This is always going to be `code_interpreter` for this
	// type of tool call.
	//
	// This field can be elided, and will be automatically set as "code_interpreter".
	Type constant.CodeInterpreter `json:"type,required"`
	JSON struct {
		ID              resp.Field
		CodeInterpreter resp.Field
		Type            resp.Field
		raw             string
	} `json:"-"`
}

func (r CodeInterpreterToolCall) RawJSON() string { return r.JSON.raw }
func (r *CodeInterpreterToolCall) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The Code Interpreter tool call definition.
type CodeInterpreterToolCallCodeInterpreter struct {
	// The input to the Code Interpreter tool call.
	Input string `json:"input,omitzero,required"`
	// The outputs from the Code Interpreter tool call. Code Interpreter can output one
	// or more items, including text (`logs`) or images (`image`). Each of these are
	// represented by a different object type.
	Outputs []CodeInterpreterToolCallCodeInterpreterOutputsUnion `json:"outputs,omitzero,required"`
	JSON    struct {
		Input   resp.Field
		Outputs resp.Field
		raw     string
	} `json:"-"`
}

func (r CodeInterpreterToolCallCodeInterpreter) RawJSON() string { return r.JSON.raw }
func (r *CodeInterpreterToolCallCodeInterpreter) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type CodeInterpreterToolCallCodeInterpreterOutputsUnion struct {
	Logs  string                                                  `json:"logs"`
	Type  string                                                  `json:"type"`
	Image CodeInterpreterToolCallCodeInterpreterOutputsImageImage `json:"image"`
	JSON  struct {
		Logs  resp.Field
		Type  resp.Field
		Image resp.Field
		raw   string
	} `json:"-"`
}

// note: this function is generated only for discriminated unions
func (u CodeInterpreterToolCallCodeInterpreterOutputsUnion) Variant() (res struct {
	OfLogs  *CodeInterpreterToolCallCodeInterpreterOutputsLogs
	OfImage *CodeInterpreterToolCallCodeInterpreterOutputsImage
}) {
	switch u.Type {
	case "logs":
		v := u.AsLogs()
		res.OfLogs = &v
	case "image":
		v := u.AsImage()
		res.OfImage = &v
	}
	return
}

func (u CodeInterpreterToolCallCodeInterpreterOutputsUnion) WhichKind() string {
	return u.Type
}

func (u CodeInterpreterToolCallCodeInterpreterOutputsUnion) AsLogs() (v CodeInterpreterToolCallCodeInterpreterOutputsLogs) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u CodeInterpreterToolCallCodeInterpreterOutputsUnion) AsImage() (v CodeInterpreterToolCallCodeInterpreterOutputsImage) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u CodeInterpreterToolCallCodeInterpreterOutputsUnion) RawJSON() string { return u.JSON.raw }

func (r *CodeInterpreterToolCallCodeInterpreterOutputsUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Text output from the Code Interpreter tool call as part of a run step.
type CodeInterpreterToolCallCodeInterpreterOutputsLogs struct {
	// The text output from the Code Interpreter tool call.
	Logs string `json:"logs,omitzero,required"`
	// Always `logs`.
	//
	// This field can be elided, and will be automatically set as "logs".
	Type constant.Logs `json:"type,required"`
	JSON struct {
		Logs resp.Field
		Type resp.Field
		raw  string
	} `json:"-"`
}

func (r CodeInterpreterToolCallCodeInterpreterOutputsLogs) RawJSON() string { return r.JSON.raw }
func (r *CodeInterpreterToolCallCodeInterpreterOutputsLogs) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type CodeInterpreterToolCallCodeInterpreterOutputsImage struct {
	Image CodeInterpreterToolCallCodeInterpreterOutputsImageImage `json:"image,omitzero,required"`
	// Always `image`.
	//
	// This field can be elided, and will be automatically set as "image".
	Type constant.Image `json:"type,required"`
	JSON struct {
		Image resp.Field
		Type  resp.Field
		raw   string
	} `json:"-"`
}

func (r CodeInterpreterToolCallCodeInterpreterOutputsImage) RawJSON() string { return r.JSON.raw }
func (r *CodeInterpreterToolCallCodeInterpreterOutputsImage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type CodeInterpreterToolCallCodeInterpreterOutputsImageImage struct {
	// The [file](https://platform.openai.com/docs/api-reference/files) ID of the
	// image.
	FileID string `json:"file_id,omitzero,required"`
	JSON   struct {
		FileID resp.Field
		raw    string
	} `json:"-"`
}

func (r CodeInterpreterToolCallCodeInterpreterOutputsImageImage) RawJSON() string { return r.JSON.raw }
func (r *CodeInterpreterToolCallCodeInterpreterOutputsImageImage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Details of the Code Interpreter tool call the run step was involved in.
type CodeInterpreterToolCallDelta struct {
	// The index of the tool call in the tool calls array.
	Index int64 `json:"index,omitzero,required"`
	// The type of tool call. This is always going to be `code_interpreter` for this
	// type of tool call.
	//
	// This field can be elided, and will be automatically set as "code_interpreter".
	Type constant.CodeInterpreter `json:"type,required"`
	// The ID of the tool call.
	ID string `json:"id,omitzero"`
	// The Code Interpreter tool call definition.
	CodeInterpreter CodeInterpreterToolCallDeltaCodeInterpreter `json:"code_interpreter,omitzero"`
	JSON            struct {
		Index           resp.Field
		Type            resp.Field
		ID              resp.Field
		CodeInterpreter resp.Field
		raw             string
	} `json:"-"`
}

func (r CodeInterpreterToolCallDelta) RawJSON() string { return r.JSON.raw }
func (r *CodeInterpreterToolCallDelta) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The Code Interpreter tool call definition.
type CodeInterpreterToolCallDeltaCodeInterpreter struct {
	// The input to the Code Interpreter tool call.
	Input string `json:"input,omitzero"`
	// The outputs from the Code Interpreter tool call. Code Interpreter can output one
	// or more items, including text (`logs`) or images (`image`). Each of these are
	// represented by a different object type.
	Outputs []CodeInterpreterToolCallDeltaCodeInterpreterOutputsUnion `json:"outputs,omitzero"`
	JSON    struct {
		Input   resp.Field
		Outputs resp.Field
		raw     string
	} `json:"-"`
}

func (r CodeInterpreterToolCallDeltaCodeInterpreter) RawJSON() string { return r.JSON.raw }
func (r *CodeInterpreterToolCallDeltaCodeInterpreter) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type CodeInterpreterToolCallDeltaCodeInterpreterOutputsUnion struct {
	Index int64                           `json:"index"`
	Type  string                          `json:"type"`
	Logs  string                          `json:"logs"`
	Image CodeInterpreterOutputImageImage `json:"image"`
	JSON  struct {
		Index resp.Field
		Type  resp.Field
		Logs  resp.Field
		Image resp.Field
		raw   string
	} `json:"-"`
}

// note: this function is generated only for discriminated unions
func (u CodeInterpreterToolCallDeltaCodeInterpreterOutputsUnion) Variant() (res struct {
	OfLogs  *CodeInterpreterLogs
	OfImage *CodeInterpreterOutputImage
}) {
	switch u.Type {
	case "logs":
		v := u.AsLogs()
		res.OfLogs = &v
	case "image":
		v := u.AsImage()
		res.OfImage = &v
	}
	return
}

func (u CodeInterpreterToolCallDeltaCodeInterpreterOutputsUnion) WhichKind() string {
	return u.Type
}

func (u CodeInterpreterToolCallDeltaCodeInterpreterOutputsUnion) AsLogs() (v CodeInterpreterLogs) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u CodeInterpreterToolCallDeltaCodeInterpreterOutputsUnion) AsImage() (v CodeInterpreterOutputImage) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u CodeInterpreterToolCallDeltaCodeInterpreterOutputsUnion) RawJSON() string { return u.JSON.raw }

func (r *CodeInterpreterToolCallDeltaCodeInterpreterOutputsUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type FileSearchToolCall struct {
	// The ID of the tool call object.
	ID string `json:"id,omitzero,required"`
	// For now, this is always going to be an empty object.
	FileSearch FileSearchToolCallFileSearch `json:"file_search,omitzero,required"`
	// The type of tool call. This is always going to be `file_search` for this type of
	// tool call.
	//
	// This field can be elided, and will be automatically set as "file_search".
	Type constant.FileSearch `json:"type,required"`
	JSON struct {
		ID         resp.Field
		FileSearch resp.Field
		Type       resp.Field
		raw        string
	} `json:"-"`
}

func (r FileSearchToolCall) RawJSON() string { return r.JSON.raw }
func (r *FileSearchToolCall) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// For now, this is always going to be an empty object.
type FileSearchToolCallFileSearch struct {
	// The ranking options for the file search.
	RankingOptions FileSearchToolCallFileSearchRankingOptions `json:"ranking_options,omitzero"`
	// The results of the file search.
	Results []FileSearchToolCallFileSearchResult `json:"results,omitzero"`
	JSON    struct {
		RankingOptions resp.Field
		Results        resp.Field
		raw            string
	} `json:"-"`
}

func (r FileSearchToolCallFileSearch) RawJSON() string { return r.JSON.raw }
func (r *FileSearchToolCallFileSearch) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The ranking options for the file search.
type FileSearchToolCallFileSearchRankingOptions struct {
	// The ranker used for the file search.
	//
	// This field can be elided, and will be automatically set as "default_2024_08_21".
	Ranker constant.Default2024_08_21 `json:"ranker,required"`
	// The score threshold for the file search. All values must be a floating point
	// number between 0 and 1.
	ScoreThreshold float64 `json:"score_threshold,omitzero,required"`
	JSON           struct {
		Ranker         resp.Field
		ScoreThreshold resp.Field
		raw            string
	} `json:"-"`
}

func (r FileSearchToolCallFileSearchRankingOptions) RawJSON() string { return r.JSON.raw }
func (r *FileSearchToolCallFileSearchRankingOptions) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A result instance of the file search.
type FileSearchToolCallFileSearchResult struct {
	// The ID of the file that result was found in.
	FileID string `json:"file_id,omitzero,required"`
	// The name of the file that result was found in.
	FileName string `json:"file_name,omitzero,required"`
	// The score of the result. All values must be a floating point number between 0
	// and 1.
	Score float64 `json:"score,omitzero,required"`
	// The content of the result that was found. The content is only included if
	// requested via the include query parameter.
	Content []FileSearchToolCallFileSearchResultsContent `json:"content,omitzero"`
	JSON    struct {
		FileID   resp.Field
		FileName resp.Field
		Score    resp.Field
		Content  resp.Field
		raw      string
	} `json:"-"`
}

func (r FileSearchToolCallFileSearchResult) RawJSON() string { return r.JSON.raw }
func (r *FileSearchToolCallFileSearchResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type FileSearchToolCallFileSearchResultsContent struct {
	// The text content of the file.
	Text string `json:"text,omitzero"`
	// The type of the content.
	//
	// Any of "text"
	Type string `json:"type"`
	JSON struct {
		Text resp.Field
		Type resp.Field
		raw  string
	} `json:"-"`
}

func (r FileSearchToolCallFileSearchResultsContent) RawJSON() string { return r.JSON.raw }
func (r *FileSearchToolCallFileSearchResultsContent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The type of the content.
type FileSearchToolCallFileSearchResultsContentType = string

const (
	FileSearchToolCallFileSearchResultsContentTypeText FileSearchToolCallFileSearchResultsContentType = "text"
)

type FileSearchToolCallDelta struct {
	// For now, this is always going to be an empty object.
	FileSearch interface{} `json:"file_search,omitzero,required"`
	// The index of the tool call in the tool calls array.
	Index int64 `json:"index,omitzero,required"`
	// The type of tool call. This is always going to be `file_search` for this type of
	// tool call.
	//
	// This field can be elided, and will be automatically set as "file_search".
	Type constant.FileSearch `json:"type,required"`
	// The ID of the tool call object.
	ID   string `json:"id,omitzero"`
	JSON struct {
		FileSearch resp.Field
		Index      resp.Field
		Type       resp.Field
		ID         resp.Field
		raw        string
	} `json:"-"`
}

func (r FileSearchToolCallDelta) RawJSON() string { return r.JSON.raw }
func (r *FileSearchToolCallDelta) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type FunctionToolCall struct {
	// The ID of the tool call object.
	ID string `json:"id,omitzero,required"`
	// The definition of the function that was called.
	Function FunctionToolCallFunction `json:"function,omitzero,required"`
	// The type of tool call. This is always going to be `function` for this type of
	// tool call.
	//
	// This field can be elided, and will be automatically set as "function".
	Type constant.Function `json:"type,required"`
	JSON struct {
		ID       resp.Field
		Function resp.Field
		Type     resp.Field
		raw      string
	} `json:"-"`
}

func (r FunctionToolCall) RawJSON() string { return r.JSON.raw }
func (r *FunctionToolCall) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The definition of the function that was called.
type FunctionToolCallFunction struct {
	// The arguments passed to the function.
	Arguments string `json:"arguments,omitzero,required"`
	// The name of the function.
	Name string `json:"name,omitzero,required"`
	// The output of the function. This will be `null` if the outputs have not been
	// [submitted](https://platform.openai.com/docs/api-reference/runs/submitToolOutputs)
	// yet.
	Output string `json:"output,omitzero,required,nullable"`
	JSON   struct {
		Arguments resp.Field
		Name      resp.Field
		Output    resp.Field
		raw       string
	} `json:"-"`
}

func (r FunctionToolCallFunction) RawJSON() string { return r.JSON.raw }
func (r *FunctionToolCallFunction) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type FunctionToolCallDelta struct {
	// The index of the tool call in the tool calls array.
	Index int64 `json:"index,omitzero,required"`
	// The type of tool call. This is always going to be `function` for this type of
	// tool call.
	//
	// This field can be elided, and will be automatically set as "function".
	Type constant.Function `json:"type,required"`
	// The ID of the tool call object.
	ID string `json:"id,omitzero"`
	// The definition of the function that was called.
	Function FunctionToolCallDeltaFunction `json:"function,omitzero"`
	JSON     struct {
		Index    resp.Field
		Type     resp.Field
		ID       resp.Field
		Function resp.Field
		raw      string
	} `json:"-"`
}

func (r FunctionToolCallDelta) RawJSON() string { return r.JSON.raw }
func (r *FunctionToolCallDelta) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The definition of the function that was called.
type FunctionToolCallDeltaFunction struct {
	// The arguments passed to the function.
	Arguments string `json:"arguments,omitzero"`
	// The name of the function.
	Name string `json:"name,omitzero"`
	// The output of the function. This will be `null` if the outputs have not been
	// [submitted](https://platform.openai.com/docs/api-reference/runs/submitToolOutputs)
	// yet.
	Output string `json:"output,omitzero,nullable"`
	JSON   struct {
		Arguments resp.Field
		Name      resp.Field
		Output    resp.Field
		raw       string
	} `json:"-"`
}

func (r FunctionToolCallDeltaFunction) RawJSON() string { return r.JSON.raw }
func (r *FunctionToolCallDeltaFunction) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Details of the message creation by the run step.
type MessageCreationStepDetails struct {
	MessageCreation MessageCreationStepDetailsMessageCreation `json:"message_creation,omitzero,required"`
	// Always `message_creation`.
	//
	// This field can be elided, and will be automatically set as "message_creation".
	Type constant.MessageCreation `json:"type,required"`
	JSON struct {
		MessageCreation resp.Field
		Type            resp.Field
		raw             string
	} `json:"-"`
}

func (r MessageCreationStepDetails) RawJSON() string { return r.JSON.raw }
func (r *MessageCreationStepDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type MessageCreationStepDetailsMessageCreation struct {
	// The ID of the message that was created by this run step.
	MessageID string `json:"message_id,omitzero,required"`
	JSON      struct {
		MessageID resp.Field
		raw       string
	} `json:"-"`
}

func (r MessageCreationStepDetailsMessageCreation) RawJSON() string { return r.JSON.raw }
func (r *MessageCreationStepDetailsMessageCreation) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Represents a step in execution of a run.
type RunStep struct {
	// The identifier of the run step, which can be referenced in API endpoints.
	ID string `json:"id,omitzero,required"`
	// The ID of the
	// [assistant](https://platform.openai.com/docs/api-reference/assistants)
	// associated with the run step.
	AssistantID string `json:"assistant_id,omitzero,required"`
	// The Unix timestamp (in seconds) for when the run step was cancelled.
	CancelledAt int64 `json:"cancelled_at,omitzero,required,nullable"`
	// The Unix timestamp (in seconds) for when the run step completed.
	CompletedAt int64 `json:"completed_at,omitzero,required,nullable"`
	// The Unix timestamp (in seconds) for when the run step was created.
	CreatedAt int64 `json:"created_at,omitzero,required"`
	// The Unix timestamp (in seconds) for when the run step expired. A step is
	// considered expired if the parent run is expired.
	ExpiredAt int64 `json:"expired_at,omitzero,required,nullable"`
	// The Unix timestamp (in seconds) for when the run step failed.
	FailedAt int64 `json:"failed_at,omitzero,required,nullable"`
	// The last error associated with this run step. Will be `null` if there are no
	// errors.
	LastError RunStepLastError `json:"last_error,omitzero,required,nullable"`
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard.
	//
	// Keys are strings with a maximum length of 64 characters. Values are strings with
	// a maximum length of 512 characters.
	Metadata shared.Metadata `json:"metadata,omitzero,required,nullable"`
	// The object type, which is always `thread.run.step`.
	//
	// This field can be elided, and will be automatically set as "thread.run.step".
	Object constant.ThreadRunStep `json:"object,required"`
	// The ID of the [run](https://platform.openai.com/docs/api-reference/runs) that
	// this run step is a part of.
	RunID string `json:"run_id,omitzero,required"`
	// The status of the run step, which can be either `in_progress`, `cancelled`,
	// `failed`, `completed`, or `expired`.
	//
	// Any of "in_progress", "cancelled", "failed", "completed", "expired"
	Status string `json:"status,omitzero,required"`
	// The details of the run step.
	StepDetails RunStepStepDetailsUnion `json:"step_details,omitzero,required"`
	// The ID of the [thread](https://platform.openai.com/docs/api-reference/threads)
	// that was run.
	ThreadID string `json:"thread_id,omitzero,required"`
	// The type of run step, which can be either `message_creation` or `tool_calls`.
	//
	// Any of "message_creation", "tool_calls"
	Type string `json:"type,omitzero,required"`
	// Usage statistics related to the run step. This value will be `null` while the
	// run step's status is `in_progress`.
	Usage RunStepUsage `json:"usage,omitzero,required,nullable"`
	JSON  struct {
		ID          resp.Field
		AssistantID resp.Field
		CancelledAt resp.Field
		CompletedAt resp.Field
		CreatedAt   resp.Field
		ExpiredAt   resp.Field
		FailedAt    resp.Field
		LastError   resp.Field
		Metadata    resp.Field
		Object      resp.Field
		RunID       resp.Field
		Status      resp.Field
		StepDetails resp.Field
		ThreadID    resp.Field
		Type        resp.Field
		Usage       resp.Field
		raw         string
	} `json:"-"`
}

func (r RunStep) RawJSON() string { return r.JSON.raw }
func (r *RunStep) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The last error associated with this run step. Will be `null` if there are no
// errors.
type RunStepLastError struct {
	// One of `server_error` or `rate_limit_exceeded`.
	//
	// Any of "server_error", "rate_limit_exceeded"
	Code string `json:"code,omitzero,required"`
	// A human-readable description of the error.
	Message string `json:"message,omitzero,required"`
	JSON    struct {
		Code    resp.Field
		Message resp.Field
		raw     string
	} `json:"-"`
}

func (r RunStepLastError) RawJSON() string { return r.JSON.raw }
func (r *RunStepLastError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// One of `server_error` or `rate_limit_exceeded`.
type RunStepLastErrorCode = string

const (
	RunStepLastErrorCodeServerError       RunStepLastErrorCode = "server_error"
	RunStepLastErrorCodeRateLimitExceeded RunStepLastErrorCode = "rate_limit_exceeded"
)

// The status of the run step, which can be either `in_progress`, `cancelled`,
// `failed`, `completed`, or `expired`.
type RunStepStatus = string

const (
	RunStepStatusInProgress RunStepStatus = "in_progress"
	RunStepStatusCancelled  RunStepStatus = "cancelled"
	RunStepStatusFailed     RunStepStatus = "failed"
	RunStepStatusCompleted  RunStepStatus = "completed"
	RunStepStatusExpired    RunStepStatus = "expired"
)

type RunStepStepDetailsUnion struct {
	MessageCreation MessageCreationStepDetailsMessageCreation `json:"message_creation"`
	Type            string                                    `json:"type"`
	ToolCalls       []ToolCallUnion                           `json:"tool_calls"`
	JSON            struct {
		MessageCreation resp.Field
		Type            resp.Field
		ToolCalls       resp.Field
		raw             string
	} `json:"-"`
}

// note: this function is generated only for discriminated unions
func (u RunStepStepDetailsUnion) Variant() (res struct {
	OfMessageCreation *MessageCreationStepDetails
	OfToolCalls       *ToolCallsStepDetails
}) {
	switch u.Type {
	case "message_creation":
		v := u.AsMessageCreation()
		res.OfMessageCreation = &v
	case "tool_calls":
		v := u.AsToolCalls()
		res.OfToolCalls = &v
	}
	return
}

func (u RunStepStepDetailsUnion) WhichKind() string {
	return u.Type
}

func (u RunStepStepDetailsUnion) AsMessageCreation() (v MessageCreationStepDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RunStepStepDetailsUnion) AsToolCalls() (v ToolCallsStepDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RunStepStepDetailsUnion) RawJSON() string { return u.JSON.raw }

func (r *RunStepStepDetailsUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The type of run step, which can be either `message_creation` or `tool_calls`.
type RunStepType = string

const (
	RunStepTypeMessageCreation RunStepType = "message_creation"
	RunStepTypeToolCalls       RunStepType = "tool_calls"
)

// Usage statistics related to the run step. This value will be `null` while the
// run step's status is `in_progress`.
type RunStepUsage struct {
	// Number of completion tokens used over the course of the run step.
	CompletionTokens int64 `json:"completion_tokens,omitzero,required"`
	// Number of prompt tokens used over the course of the run step.
	PromptTokens int64 `json:"prompt_tokens,omitzero,required"`
	// Total number of tokens used (prompt + completion).
	TotalTokens int64 `json:"total_tokens,omitzero,required"`
	JSON        struct {
		CompletionTokens resp.Field
		PromptTokens     resp.Field
		TotalTokens      resp.Field
		raw              string
	} `json:"-"`
}

func (r RunStepUsage) RawJSON() string { return r.JSON.raw }
func (r *RunStepUsage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The delta containing the fields that have changed on the run step.
type RunStepDelta struct {
	// The details of the run step.
	StepDetails RunStepDeltaStepDetailsUnion `json:"step_details,omitzero"`
	JSON        struct {
		StepDetails resp.Field
		raw         string
	} `json:"-"`
}

func (r RunStepDelta) RawJSON() string { return r.JSON.raw }
func (r *RunStepDelta) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type RunStepDeltaStepDetailsUnion struct {
	Type            string                                  `json:"type"`
	MessageCreation RunStepDeltaMessageDeltaMessageCreation `json:"message_creation"`
	ToolCalls       []ToolCallDeltaUnion                    `json:"tool_calls"`
	JSON            struct {
		Type            resp.Field
		MessageCreation resp.Field
		ToolCalls       resp.Field
		raw             string
	} `json:"-"`
}

// note: this function is generated only for discriminated unions
func (u RunStepDeltaStepDetailsUnion) Variant() (res struct {
	OfMessageCreation *RunStepDeltaMessageDelta
	OfToolCalls       *ToolCallDeltaObject
}) {
	switch u.Type {
	case "message_creation":
		v := u.AsMessageCreation()
		res.OfMessageCreation = &v
	case "tool_calls":
		v := u.AsToolCalls()
		res.OfToolCalls = &v
	}
	return
}

func (u RunStepDeltaStepDetailsUnion) WhichKind() string {
	return u.Type
}

func (u RunStepDeltaStepDetailsUnion) AsMessageCreation() (v RunStepDeltaMessageDelta) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RunStepDeltaStepDetailsUnion) AsToolCalls() (v ToolCallDeltaObject) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RunStepDeltaStepDetailsUnion) RawJSON() string { return u.JSON.raw }

func (r *RunStepDeltaStepDetailsUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Represents a run step delta i.e. any changed fields on a run step during
// streaming.
type RunStepDeltaEvent struct {
	// The identifier of the run step, which can be referenced in API endpoints.
	ID string `json:"id,omitzero,required"`
	// The delta containing the fields that have changed on the run step.
	Delta RunStepDelta `json:"delta,omitzero,required"`
	// The object type, which is always `thread.run.step.delta`.
	//
	// This field can be elided, and will be automatically set as
	// "thread.run.step.delta".
	Object constant.ThreadRunStepDelta `json:"object,required"`
	JSON   struct {
		ID     resp.Field
		Delta  resp.Field
		Object resp.Field
		raw    string
	} `json:"-"`
}

func (r RunStepDeltaEvent) RawJSON() string { return r.JSON.raw }
func (r *RunStepDeltaEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Details of the message creation by the run step.
type RunStepDeltaMessageDelta struct {
	// Always `message_creation`.
	//
	// This field can be elided, and will be automatically set as "message_creation".
	Type            constant.MessageCreation                `json:"type,required"`
	MessageCreation RunStepDeltaMessageDeltaMessageCreation `json:"message_creation,omitzero"`
	JSON            struct {
		Type            resp.Field
		MessageCreation resp.Field
		raw             string
	} `json:"-"`
}

func (r RunStepDeltaMessageDelta) RawJSON() string { return r.JSON.raw }
func (r *RunStepDeltaMessageDelta) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type RunStepDeltaMessageDeltaMessageCreation struct {
	// The ID of the message that was created by this run step.
	MessageID string `json:"message_id,omitzero"`
	JSON      struct {
		MessageID resp.Field
		raw       string
	} `json:"-"`
}

func (r RunStepDeltaMessageDeltaMessageCreation) RawJSON() string { return r.JSON.raw }
func (r *RunStepDeltaMessageDeltaMessageCreation) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type RunStepInclude string

const (
	RunStepIncludeStepDetailsToolCallsFileSearchResultsContent RunStepInclude = "step_details.tool_calls[*].file_search.results[*].content"
)

type ToolCallUnion struct {
	ID              string                                 `json:"id"`
	CodeInterpreter CodeInterpreterToolCallCodeInterpreter `json:"code_interpreter"`
	Type            string                                 `json:"type"`
	FileSearch      FileSearchToolCallFileSearch           `json:"file_search"`
	Function        FunctionToolCallFunction               `json:"function"`
	JSON            struct {
		ID              resp.Field
		CodeInterpreter resp.Field
		Type            resp.Field
		FileSearch      resp.Field
		Function        resp.Field
		raw             string
	} `json:"-"`
}

// note: this function is generated only for discriminated unions
func (u ToolCallUnion) Variant() (res struct {
	OfCodeInterpreter *CodeInterpreterToolCall
	OfFileSearch      *FileSearchToolCall
	OfFunction        *FunctionToolCall
}) {
	switch u.Type {
	case "code_interpreter":
		v := u.AsCodeInterpreter()
		res.OfCodeInterpreter = &v
	case "file_search":
		v := u.AsFileSearch()
		res.OfFileSearch = &v
	case "function":
		v := u.AsFunction()
		res.OfFunction = &v
	}
	return
}

func (u ToolCallUnion) WhichKind() string {
	return u.Type
}

func (u ToolCallUnion) AsCodeInterpreter() (v CodeInterpreterToolCall) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ToolCallUnion) AsFileSearch() (v FileSearchToolCall) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ToolCallUnion) AsFunction() (v FunctionToolCall) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ToolCallUnion) RawJSON() string { return u.JSON.raw }

func (r *ToolCallUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ToolCallDeltaUnion struct {
	Index           int64                                       `json:"index"`
	Type            string                                      `json:"type"`
	ID              string                                      `json:"id"`
	CodeInterpreter CodeInterpreterToolCallDeltaCodeInterpreter `json:"code_interpreter"`
	FileSearch      interface{}                                 `json:"file_search"`
	Function        FunctionToolCallDeltaFunction               `json:"function"`
	JSON            struct {
		Index           resp.Field
		Type            resp.Field
		ID              resp.Field
		CodeInterpreter resp.Field
		FileSearch      resp.Field
		Function        resp.Field
		raw             string
	} `json:"-"`
}

// note: this function is generated only for discriminated unions
func (u ToolCallDeltaUnion) Variant() (res struct {
	OfCodeInterpreter *CodeInterpreterToolCallDelta
	OfFileSearch      *FileSearchToolCallDelta
	OfFunction        *FunctionToolCallDelta
}) {
	switch u.Type {
	case "code_interpreter":
		v := u.AsCodeInterpreter()
		res.OfCodeInterpreter = &v
	case "file_search":
		v := u.AsFileSearch()
		res.OfFileSearch = &v
	case "function":
		v := u.AsFunction()
		res.OfFunction = &v
	}
	return
}

func (u ToolCallDeltaUnion) WhichKind() string {
	return u.Type
}

func (u ToolCallDeltaUnion) AsCodeInterpreter() (v CodeInterpreterToolCallDelta) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ToolCallDeltaUnion) AsFileSearch() (v FileSearchToolCallDelta) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ToolCallDeltaUnion) AsFunction() (v FunctionToolCallDelta) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ToolCallDeltaUnion) RawJSON() string { return u.JSON.raw }

func (r *ToolCallDeltaUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Details of the tool call.
type ToolCallDeltaObject struct {
	// Always `tool_calls`.
	//
	// This field can be elided, and will be automatically set as "tool_calls".
	Type constant.ToolCalls `json:"type,required"`
	// An array of tool calls the run step was involved in. These can be associated
	// with one of three types of tools: `code_interpreter`, `file_search`, or
	// `function`.
	ToolCalls []ToolCallDeltaUnion `json:"tool_calls,omitzero"`
	JSON      struct {
		Type      resp.Field
		ToolCalls resp.Field
		raw       string
	} `json:"-"`
}

func (r ToolCallDeltaObject) RawJSON() string { return r.JSON.raw }
func (r *ToolCallDeltaObject) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Details of the tool call.
type ToolCallsStepDetails struct {
	// An array of tool calls the run step was involved in. These can be associated
	// with one of three types of tools: `code_interpreter`, `file_search`, or
	// `function`.
	ToolCalls []ToolCallUnion `json:"tool_calls,omitzero,required"`
	// Always `tool_calls`.
	//
	// This field can be elided, and will be automatically set as "tool_calls".
	Type constant.ToolCalls `json:"type,required"`
	JSON struct {
		ToolCalls resp.Field
		Type      resp.Field
		raw       string
	} `json:"-"`
}

func (r ToolCallsStepDetails) RawJSON() string { return r.JSON.raw }
func (r *ToolCallsStepDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaThreadRunStepGetParams struct {
	// A list of additional fields to include in the response. Currently the only
	// supported value is `step_details.tool_calls[*].file_search.results[*].content`
	// to fetch the file search result content.
	//
	// See the
	// [file search tool documentation](https://platform.openai.com/docs/assistants/tools/file-search#customizing-file-search-settings)
	// for more information.
	Include []RunStepInclude `query:"include,omitzero"`
	apiobject
}

func (f BetaThreadRunStepGetParams) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

// URLQuery serializes [BetaThreadRunStepGetParams]'s query parameters as
// `url.Values`.
func (r BetaThreadRunStepGetParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BetaThreadRunStepListParams struct {
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
	// A list of additional fields to include in the response. Currently the only
	// supported value is `step_details.tool_calls[*].file_search.results[*].content`
	// to fetch the file search result content.
	//
	// See the
	// [file search tool documentation](https://platform.openai.com/docs/assistants/tools/file-search#customizing-file-search-settings)
	// for more information.
	Include []RunStepInclude `query:"include,omitzero"`
	// A limit on the number of objects to be returned. Limit can range between 1 and
	// 100, and the default is 20.
	Limit param.Int `query:"limit,omitzero"`
	// Sort order by the `created_at` timestamp of the objects. `asc` for ascending
	// order and `desc` for descending order.
	//
	// Any of "asc", "desc"
	Order BetaThreadRunStepListParamsOrder `query:"order,omitzero"`
	apiobject
}

func (f BetaThreadRunStepListParams) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

// URLQuery serializes [BetaThreadRunStepListParams]'s query parameters as
// `url.Values`.
func (r BetaThreadRunStepListParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Sort order by the `created_at` timestamp of the objects. `asc` for ascending
// order and `desc` for descending order.
type BetaThreadRunStepListParamsOrder string

const (
	BetaThreadRunStepListParamsOrderAsc  BetaThreadRunStepListParamsOrder = "asc"
	BetaThreadRunStepListParamsOrderDesc BetaThreadRunStepListParamsOrder = "desc"
)

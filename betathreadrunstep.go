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
func NewBetaThreadRunStepService(opts ...option.RequestOption) (r *BetaThreadRunStepService) {
	r = &BetaThreadRunStepService{}
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
	Index int64 `json:"index,required"`
	// Always `logs`.
	Type CodeInterpreterLogsType `json:"type,required"`
	// The text output from the Code Interpreter tool call.
	Logs string                  `json:"logs"`
	JSON codeInterpreterLogsJSON `json:"-"`
}

// codeInterpreterLogsJSON contains the JSON metadata for the struct
// [CodeInterpreterLogs]
type codeInterpreterLogsJSON struct {
	Index       apijson.Field
	Type        apijson.Field
	Logs        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CodeInterpreterLogs) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r codeInterpreterLogsJSON) RawJSON() string {
	return r.raw
}

func (r CodeInterpreterLogs) implementsCodeInterpreterToolCallDeltaCodeInterpreterOutput() {}

// Always `logs`.
type CodeInterpreterLogsType string

const (
	CodeInterpreterLogsTypeLogs CodeInterpreterLogsType = "logs"
)

func (r CodeInterpreterLogsType) IsKnown() bool {
	switch r {
	case CodeInterpreterLogsTypeLogs:
		return true
	}
	return false
}

type CodeInterpreterOutputImage struct {
	// The index of the output in the outputs array.
	Index int64 `json:"index,required"`
	// Always `image`.
	Type  CodeInterpreterOutputImageType  `json:"type,required"`
	Image CodeInterpreterOutputImageImage `json:"image"`
	JSON  codeInterpreterOutputImageJSON  `json:"-"`
}

// codeInterpreterOutputImageJSON contains the JSON metadata for the struct
// [CodeInterpreterOutputImage]
type codeInterpreterOutputImageJSON struct {
	Index       apijson.Field
	Type        apijson.Field
	Image       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CodeInterpreterOutputImage) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r codeInterpreterOutputImageJSON) RawJSON() string {
	return r.raw
}

func (r CodeInterpreterOutputImage) implementsCodeInterpreterToolCallDeltaCodeInterpreterOutput() {}

// Always `image`.
type CodeInterpreterOutputImageType string

const (
	CodeInterpreterOutputImageTypeImage CodeInterpreterOutputImageType = "image"
)

func (r CodeInterpreterOutputImageType) IsKnown() bool {
	switch r {
	case CodeInterpreterOutputImageTypeImage:
		return true
	}
	return false
}

type CodeInterpreterOutputImageImage struct {
	// The [file](https://platform.openai.com/docs/api-reference/files) ID of the
	// image.
	FileID string                              `json:"file_id"`
	JSON   codeInterpreterOutputImageImageJSON `json:"-"`
}

// codeInterpreterOutputImageImageJSON contains the JSON metadata for the struct
// [CodeInterpreterOutputImageImage]
type codeInterpreterOutputImageImageJSON struct {
	FileID      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CodeInterpreterOutputImageImage) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r codeInterpreterOutputImageImageJSON) RawJSON() string {
	return r.raw
}

// Details of the Code Interpreter tool call the run step was involved in.
type CodeInterpreterToolCall struct {
	// The ID of the tool call.
	ID string `json:"id,required"`
	// The Code Interpreter tool call definition.
	CodeInterpreter CodeInterpreterToolCallCodeInterpreter `json:"code_interpreter,required"`
	// The type of tool call. This is always going to be `code_interpreter` for this
	// type of tool call.
	Type CodeInterpreterToolCallType `json:"type,required"`
	JSON codeInterpreterToolCallJSON `json:"-"`
}

// codeInterpreterToolCallJSON contains the JSON metadata for the struct
// [CodeInterpreterToolCall]
type codeInterpreterToolCallJSON struct {
	ID              apijson.Field
	CodeInterpreter apijson.Field
	Type            apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r *CodeInterpreterToolCall) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r codeInterpreterToolCallJSON) RawJSON() string {
	return r.raw
}

func (r CodeInterpreterToolCall) implementsToolCall() {}

// The Code Interpreter tool call definition.
type CodeInterpreterToolCallCodeInterpreter struct {
	// The input to the Code Interpreter tool call.
	Input string `json:"input,required"`
	// The outputs from the Code Interpreter tool call. Code Interpreter can output one
	// or more items, including text (`logs`) or images (`image`). Each of these are
	// represented by a different object type.
	Outputs []CodeInterpreterToolCallCodeInterpreterOutput `json:"outputs,required"`
	JSON    codeInterpreterToolCallCodeInterpreterJSON     `json:"-"`
}

// codeInterpreterToolCallCodeInterpreterJSON contains the JSON metadata for the
// struct [CodeInterpreterToolCallCodeInterpreter]
type codeInterpreterToolCallCodeInterpreterJSON struct {
	Input       apijson.Field
	Outputs     apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CodeInterpreterToolCallCodeInterpreter) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r codeInterpreterToolCallCodeInterpreterJSON) RawJSON() string {
	return r.raw
}

// Text output from the Code Interpreter tool call as part of a run step.
type CodeInterpreterToolCallCodeInterpreterOutput struct {
	// Always `logs`.
	Type CodeInterpreterToolCallCodeInterpreterOutputsType `json:"type,required"`
	// This field can have the runtime type of
	// [CodeInterpreterToolCallCodeInterpreterOutputsImageImage].
	Image interface{} `json:"image"`
	// The text output from the Code Interpreter tool call.
	Logs  string                                           `json:"logs"`
	JSON  codeInterpreterToolCallCodeInterpreterOutputJSON `json:"-"`
	union CodeInterpreterToolCallCodeInterpreterOutputsUnion
}

// codeInterpreterToolCallCodeInterpreterOutputJSON contains the JSON metadata for
// the struct [CodeInterpreterToolCallCodeInterpreterOutput]
type codeInterpreterToolCallCodeInterpreterOutputJSON struct {
	Type        apijson.Field
	Image       apijson.Field
	Logs        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r codeInterpreterToolCallCodeInterpreterOutputJSON) RawJSON() string {
	return r.raw
}

func (r *CodeInterpreterToolCallCodeInterpreterOutput) UnmarshalJSON(data []byte) (err error) {
	*r = CodeInterpreterToolCallCodeInterpreterOutput{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [CodeInterpreterToolCallCodeInterpreterOutputsUnion] interface
// which you can cast to the specific types for more type safety.
//
// Possible runtime types of the union are
// [CodeInterpreterToolCallCodeInterpreterOutputsLogs],
// [CodeInterpreterToolCallCodeInterpreterOutputsImage].
func (r CodeInterpreterToolCallCodeInterpreterOutput) AsUnion() CodeInterpreterToolCallCodeInterpreterOutputsUnion {
	return r.union
}

// Text output from the Code Interpreter tool call as part of a run step.
//
// Union satisfied by [CodeInterpreterToolCallCodeInterpreterOutputsLogs] or
// [CodeInterpreterToolCallCodeInterpreterOutputsImage].
type CodeInterpreterToolCallCodeInterpreterOutputsUnion interface {
	implementsCodeInterpreterToolCallCodeInterpreterOutput()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*CodeInterpreterToolCallCodeInterpreterOutputsUnion)(nil)).Elem(),
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(CodeInterpreterToolCallCodeInterpreterOutputsLogs{}),
			DiscriminatorValue: "logs",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(CodeInterpreterToolCallCodeInterpreterOutputsImage{}),
			DiscriminatorValue: "image",
		},
	)
}

// Text output from the Code Interpreter tool call as part of a run step.
type CodeInterpreterToolCallCodeInterpreterOutputsLogs struct {
	// The text output from the Code Interpreter tool call.
	Logs string `json:"logs,required"`
	// Always `logs`.
	Type CodeInterpreterToolCallCodeInterpreterOutputsLogsType `json:"type,required"`
	JSON codeInterpreterToolCallCodeInterpreterOutputsLogsJSON `json:"-"`
}

// codeInterpreterToolCallCodeInterpreterOutputsLogsJSON contains the JSON metadata
// for the struct [CodeInterpreterToolCallCodeInterpreterOutputsLogs]
type codeInterpreterToolCallCodeInterpreterOutputsLogsJSON struct {
	Logs        apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CodeInterpreterToolCallCodeInterpreterOutputsLogs) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r codeInterpreterToolCallCodeInterpreterOutputsLogsJSON) RawJSON() string {
	return r.raw
}

func (r CodeInterpreterToolCallCodeInterpreterOutputsLogs) implementsCodeInterpreterToolCallCodeInterpreterOutput() {
}

// Always `logs`.
type CodeInterpreterToolCallCodeInterpreterOutputsLogsType string

const (
	CodeInterpreterToolCallCodeInterpreterOutputsLogsTypeLogs CodeInterpreterToolCallCodeInterpreterOutputsLogsType = "logs"
)

func (r CodeInterpreterToolCallCodeInterpreterOutputsLogsType) IsKnown() bool {
	switch r {
	case CodeInterpreterToolCallCodeInterpreterOutputsLogsTypeLogs:
		return true
	}
	return false
}

type CodeInterpreterToolCallCodeInterpreterOutputsImage struct {
	Image CodeInterpreterToolCallCodeInterpreterOutputsImageImage `json:"image,required"`
	// Always `image`.
	Type CodeInterpreterToolCallCodeInterpreterOutputsImageType `json:"type,required"`
	JSON codeInterpreterToolCallCodeInterpreterOutputsImageJSON `json:"-"`
}

// codeInterpreterToolCallCodeInterpreterOutputsImageJSON contains the JSON
// metadata for the struct [CodeInterpreterToolCallCodeInterpreterOutputsImage]
type codeInterpreterToolCallCodeInterpreterOutputsImageJSON struct {
	Image       apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CodeInterpreterToolCallCodeInterpreterOutputsImage) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r codeInterpreterToolCallCodeInterpreterOutputsImageJSON) RawJSON() string {
	return r.raw
}

func (r CodeInterpreterToolCallCodeInterpreterOutputsImage) implementsCodeInterpreterToolCallCodeInterpreterOutput() {
}

type CodeInterpreterToolCallCodeInterpreterOutputsImageImage struct {
	// The [file](https://platform.openai.com/docs/api-reference/files) ID of the
	// image.
	FileID string                                                      `json:"file_id,required"`
	JSON   codeInterpreterToolCallCodeInterpreterOutputsImageImageJSON `json:"-"`
}

// codeInterpreterToolCallCodeInterpreterOutputsImageImageJSON contains the JSON
// metadata for the struct
// [CodeInterpreterToolCallCodeInterpreterOutputsImageImage]
type codeInterpreterToolCallCodeInterpreterOutputsImageImageJSON struct {
	FileID      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CodeInterpreterToolCallCodeInterpreterOutputsImageImage) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r codeInterpreterToolCallCodeInterpreterOutputsImageImageJSON) RawJSON() string {
	return r.raw
}

// Always `image`.
type CodeInterpreterToolCallCodeInterpreterOutputsImageType string

const (
	CodeInterpreterToolCallCodeInterpreterOutputsImageTypeImage CodeInterpreterToolCallCodeInterpreterOutputsImageType = "image"
)

func (r CodeInterpreterToolCallCodeInterpreterOutputsImageType) IsKnown() bool {
	switch r {
	case CodeInterpreterToolCallCodeInterpreterOutputsImageTypeImage:
		return true
	}
	return false
}

// Always `logs`.
type CodeInterpreterToolCallCodeInterpreterOutputsType string

const (
	CodeInterpreterToolCallCodeInterpreterOutputsTypeLogs  CodeInterpreterToolCallCodeInterpreterOutputsType = "logs"
	CodeInterpreterToolCallCodeInterpreterOutputsTypeImage CodeInterpreterToolCallCodeInterpreterOutputsType = "image"
)

func (r CodeInterpreterToolCallCodeInterpreterOutputsType) IsKnown() bool {
	switch r {
	case CodeInterpreterToolCallCodeInterpreterOutputsTypeLogs, CodeInterpreterToolCallCodeInterpreterOutputsTypeImage:
		return true
	}
	return false
}

// The type of tool call. This is always going to be `code_interpreter` for this
// type of tool call.
type CodeInterpreterToolCallType string

const (
	CodeInterpreterToolCallTypeCodeInterpreter CodeInterpreterToolCallType = "code_interpreter"
)

func (r CodeInterpreterToolCallType) IsKnown() bool {
	switch r {
	case CodeInterpreterToolCallTypeCodeInterpreter:
		return true
	}
	return false
}

// Details of the Code Interpreter tool call the run step was involved in.
type CodeInterpreterToolCallDelta struct {
	// The index of the tool call in the tool calls array.
	Index int64 `json:"index,required"`
	// The type of tool call. This is always going to be `code_interpreter` for this
	// type of tool call.
	Type CodeInterpreterToolCallDeltaType `json:"type,required"`
	// The ID of the tool call.
	ID string `json:"id"`
	// The Code Interpreter tool call definition.
	CodeInterpreter CodeInterpreterToolCallDeltaCodeInterpreter `json:"code_interpreter"`
	JSON            codeInterpreterToolCallDeltaJSON            `json:"-"`
}

// codeInterpreterToolCallDeltaJSON contains the JSON metadata for the struct
// [CodeInterpreterToolCallDelta]
type codeInterpreterToolCallDeltaJSON struct {
	Index           apijson.Field
	Type            apijson.Field
	ID              apijson.Field
	CodeInterpreter apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r *CodeInterpreterToolCallDelta) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r codeInterpreterToolCallDeltaJSON) RawJSON() string {
	return r.raw
}

func (r CodeInterpreterToolCallDelta) implementsToolCallDelta() {}

// The type of tool call. This is always going to be `code_interpreter` for this
// type of tool call.
type CodeInterpreterToolCallDeltaType string

const (
	CodeInterpreterToolCallDeltaTypeCodeInterpreter CodeInterpreterToolCallDeltaType = "code_interpreter"
)

func (r CodeInterpreterToolCallDeltaType) IsKnown() bool {
	switch r {
	case CodeInterpreterToolCallDeltaTypeCodeInterpreter:
		return true
	}
	return false
}

// The Code Interpreter tool call definition.
type CodeInterpreterToolCallDeltaCodeInterpreter struct {
	// The input to the Code Interpreter tool call.
	Input string `json:"input"`
	// The outputs from the Code Interpreter tool call. Code Interpreter can output one
	// or more items, including text (`logs`) or images (`image`). Each of these are
	// represented by a different object type.
	Outputs []CodeInterpreterToolCallDeltaCodeInterpreterOutput `json:"outputs"`
	JSON    codeInterpreterToolCallDeltaCodeInterpreterJSON     `json:"-"`
}

// codeInterpreterToolCallDeltaCodeInterpreterJSON contains the JSON metadata for
// the struct [CodeInterpreterToolCallDeltaCodeInterpreter]
type codeInterpreterToolCallDeltaCodeInterpreterJSON struct {
	Input       apijson.Field
	Outputs     apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CodeInterpreterToolCallDeltaCodeInterpreter) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r codeInterpreterToolCallDeltaCodeInterpreterJSON) RawJSON() string {
	return r.raw
}

// Text output from the Code Interpreter tool call as part of a run step.
type CodeInterpreterToolCallDeltaCodeInterpreterOutput struct {
	// The index of the output in the outputs array.
	Index int64 `json:"index,required"`
	// Always `logs`.
	Type CodeInterpreterToolCallDeltaCodeInterpreterOutputsType `json:"type,required"`
	// This field can have the runtime type of [CodeInterpreterOutputImageImage].
	Image interface{} `json:"image"`
	// The text output from the Code Interpreter tool call.
	Logs  string                                                `json:"logs"`
	JSON  codeInterpreterToolCallDeltaCodeInterpreterOutputJSON `json:"-"`
	union CodeInterpreterToolCallDeltaCodeInterpreterOutputsUnion
}

// codeInterpreterToolCallDeltaCodeInterpreterOutputJSON contains the JSON metadata
// for the struct [CodeInterpreterToolCallDeltaCodeInterpreterOutput]
type codeInterpreterToolCallDeltaCodeInterpreterOutputJSON struct {
	Index       apijson.Field
	Type        apijson.Field
	Image       apijson.Field
	Logs        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r codeInterpreterToolCallDeltaCodeInterpreterOutputJSON) RawJSON() string {
	return r.raw
}

func (r *CodeInterpreterToolCallDeltaCodeInterpreterOutput) UnmarshalJSON(data []byte) (err error) {
	*r = CodeInterpreterToolCallDeltaCodeInterpreterOutput{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [CodeInterpreterToolCallDeltaCodeInterpreterOutputsUnion]
// interface which you can cast to the specific types for more type safety.
//
// Possible runtime types of the union are [CodeInterpreterLogs],
// [CodeInterpreterOutputImage].
func (r CodeInterpreterToolCallDeltaCodeInterpreterOutput) AsUnion() CodeInterpreterToolCallDeltaCodeInterpreterOutputsUnion {
	return r.union
}

// Text output from the Code Interpreter tool call as part of a run step.
//
// Union satisfied by [CodeInterpreterLogs] or [CodeInterpreterOutputImage].
type CodeInterpreterToolCallDeltaCodeInterpreterOutputsUnion interface {
	implementsCodeInterpreterToolCallDeltaCodeInterpreterOutput()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*CodeInterpreterToolCallDeltaCodeInterpreterOutputsUnion)(nil)).Elem(),
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(CodeInterpreterLogs{}),
			DiscriminatorValue: "logs",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(CodeInterpreterOutputImage{}),
			DiscriminatorValue: "image",
		},
	)
}

// Always `logs`.
type CodeInterpreterToolCallDeltaCodeInterpreterOutputsType string

const (
	CodeInterpreterToolCallDeltaCodeInterpreterOutputsTypeLogs  CodeInterpreterToolCallDeltaCodeInterpreterOutputsType = "logs"
	CodeInterpreterToolCallDeltaCodeInterpreterOutputsTypeImage CodeInterpreterToolCallDeltaCodeInterpreterOutputsType = "image"
)

func (r CodeInterpreterToolCallDeltaCodeInterpreterOutputsType) IsKnown() bool {
	switch r {
	case CodeInterpreterToolCallDeltaCodeInterpreterOutputsTypeLogs, CodeInterpreterToolCallDeltaCodeInterpreterOutputsTypeImage:
		return true
	}
	return false
}

type FileSearchToolCall struct {
	// The ID of the tool call object.
	ID string `json:"id,required"`
	// For now, this is always going to be an empty object.
	FileSearch FileSearchToolCallFileSearch `json:"file_search,required"`
	// The type of tool call. This is always going to be `file_search` for this type of
	// tool call.
	Type FileSearchToolCallType `json:"type,required"`
	JSON fileSearchToolCallJSON `json:"-"`
}

// fileSearchToolCallJSON contains the JSON metadata for the struct
// [FileSearchToolCall]
type fileSearchToolCallJSON struct {
	ID          apijson.Field
	FileSearch  apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *FileSearchToolCall) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r fileSearchToolCallJSON) RawJSON() string {
	return r.raw
}

func (r FileSearchToolCall) implementsToolCall() {}

// For now, this is always going to be an empty object.
type FileSearchToolCallFileSearch struct {
	// The ranking options for the file search.
	RankingOptions FileSearchToolCallFileSearchRankingOptions `json:"ranking_options"`
	// The results of the file search.
	Results []FileSearchToolCallFileSearchResult `json:"results"`
	JSON    fileSearchToolCallFileSearchJSON     `json:"-"`
}

// fileSearchToolCallFileSearchJSON contains the JSON metadata for the struct
// [FileSearchToolCallFileSearch]
type fileSearchToolCallFileSearchJSON struct {
	RankingOptions apijson.Field
	Results        apijson.Field
	raw            string
	ExtraFields    map[string]apijson.Field
}

func (r *FileSearchToolCallFileSearch) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r fileSearchToolCallFileSearchJSON) RawJSON() string {
	return r.raw
}

// The ranking options for the file search.
type FileSearchToolCallFileSearchRankingOptions struct {
	// The ranker to use for the file search. If not specified will use the `auto`
	// ranker.
	Ranker FileSearchToolCallFileSearchRankingOptionsRanker `json:"ranker,required"`
	// The score threshold for the file search. All values must be a floating point
	// number between 0 and 1.
	ScoreThreshold float64                                        `json:"score_threshold,required"`
	JSON           fileSearchToolCallFileSearchRankingOptionsJSON `json:"-"`
}

// fileSearchToolCallFileSearchRankingOptionsJSON contains the JSON metadata for
// the struct [FileSearchToolCallFileSearchRankingOptions]
type fileSearchToolCallFileSearchRankingOptionsJSON struct {
	Ranker         apijson.Field
	ScoreThreshold apijson.Field
	raw            string
	ExtraFields    map[string]apijson.Field
}

func (r *FileSearchToolCallFileSearchRankingOptions) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r fileSearchToolCallFileSearchRankingOptionsJSON) RawJSON() string {
	return r.raw
}

// The ranker to use for the file search. If not specified will use the `auto`
// ranker.
type FileSearchToolCallFileSearchRankingOptionsRanker string

const (
	FileSearchToolCallFileSearchRankingOptionsRankerAuto              FileSearchToolCallFileSearchRankingOptionsRanker = "auto"
	FileSearchToolCallFileSearchRankingOptionsRankerDefault2024_08_21 FileSearchToolCallFileSearchRankingOptionsRanker = "default_2024_08_21"
)

func (r FileSearchToolCallFileSearchRankingOptionsRanker) IsKnown() bool {
	switch r {
	case FileSearchToolCallFileSearchRankingOptionsRankerAuto, FileSearchToolCallFileSearchRankingOptionsRankerDefault2024_08_21:
		return true
	}
	return false
}

// A result instance of the file search.
type FileSearchToolCallFileSearchResult struct {
	// The ID of the file that result was found in.
	FileID string `json:"file_id,required"`
	// The name of the file that result was found in.
	FileName string `json:"file_name,required"`
	// The score of the result. All values must be a floating point number between 0
	// and 1.
	Score float64 `json:"score,required"`
	// The content of the result that was found. The content is only included if
	// requested via the include query parameter.
	Content []FileSearchToolCallFileSearchResultsContent `json:"content"`
	JSON    fileSearchToolCallFileSearchResultJSON       `json:"-"`
}

// fileSearchToolCallFileSearchResultJSON contains the JSON metadata for the struct
// [FileSearchToolCallFileSearchResult]
type fileSearchToolCallFileSearchResultJSON struct {
	FileID      apijson.Field
	FileName    apijson.Field
	Score       apijson.Field
	Content     apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *FileSearchToolCallFileSearchResult) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r fileSearchToolCallFileSearchResultJSON) RawJSON() string {
	return r.raw
}

type FileSearchToolCallFileSearchResultsContent struct {
	// The text content of the file.
	Text string `json:"text"`
	// The type of the content.
	Type FileSearchToolCallFileSearchResultsContentType `json:"type"`
	JSON fileSearchToolCallFileSearchResultsContentJSON `json:"-"`
}

// fileSearchToolCallFileSearchResultsContentJSON contains the JSON metadata for
// the struct [FileSearchToolCallFileSearchResultsContent]
type fileSearchToolCallFileSearchResultsContentJSON struct {
	Text        apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *FileSearchToolCallFileSearchResultsContent) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r fileSearchToolCallFileSearchResultsContentJSON) RawJSON() string {
	return r.raw
}

// The type of the content.
type FileSearchToolCallFileSearchResultsContentType string

const (
	FileSearchToolCallFileSearchResultsContentTypeText FileSearchToolCallFileSearchResultsContentType = "text"
)

func (r FileSearchToolCallFileSearchResultsContentType) IsKnown() bool {
	switch r {
	case FileSearchToolCallFileSearchResultsContentTypeText:
		return true
	}
	return false
}

// The type of tool call. This is always going to be `file_search` for this type of
// tool call.
type FileSearchToolCallType string

const (
	FileSearchToolCallTypeFileSearch FileSearchToolCallType = "file_search"
)

func (r FileSearchToolCallType) IsKnown() bool {
	switch r {
	case FileSearchToolCallTypeFileSearch:
		return true
	}
	return false
}

type FileSearchToolCallDelta struct {
	// For now, this is always going to be an empty object.
	FileSearch interface{} `json:"file_search,required"`
	// The index of the tool call in the tool calls array.
	Index int64 `json:"index,required"`
	// The type of tool call. This is always going to be `file_search` for this type of
	// tool call.
	Type FileSearchToolCallDeltaType `json:"type,required"`
	// The ID of the tool call object.
	ID   string                      `json:"id"`
	JSON fileSearchToolCallDeltaJSON `json:"-"`
}

// fileSearchToolCallDeltaJSON contains the JSON metadata for the struct
// [FileSearchToolCallDelta]
type fileSearchToolCallDeltaJSON struct {
	FileSearch  apijson.Field
	Index       apijson.Field
	Type        apijson.Field
	ID          apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *FileSearchToolCallDelta) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r fileSearchToolCallDeltaJSON) RawJSON() string {
	return r.raw
}

func (r FileSearchToolCallDelta) implementsToolCallDelta() {}

// The type of tool call. This is always going to be `file_search` for this type of
// tool call.
type FileSearchToolCallDeltaType string

const (
	FileSearchToolCallDeltaTypeFileSearch FileSearchToolCallDeltaType = "file_search"
)

func (r FileSearchToolCallDeltaType) IsKnown() bool {
	switch r {
	case FileSearchToolCallDeltaTypeFileSearch:
		return true
	}
	return false
}

type FunctionToolCall struct {
	// The ID of the tool call object.
	ID string `json:"id,required"`
	// The definition of the function that was called.
	Function FunctionToolCallFunction `json:"function,required"`
	// The type of tool call. This is always going to be `function` for this type of
	// tool call.
	Type FunctionToolCallType `json:"type,required"`
	JSON functionToolCallJSON `json:"-"`
}

// functionToolCallJSON contains the JSON metadata for the struct
// [FunctionToolCall]
type functionToolCallJSON struct {
	ID          apijson.Field
	Function    apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *FunctionToolCall) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r functionToolCallJSON) RawJSON() string {
	return r.raw
}

func (r FunctionToolCall) implementsToolCall() {}

// The definition of the function that was called.
type FunctionToolCallFunction struct {
	// The arguments passed to the function.
	Arguments string `json:"arguments,required"`
	// The name of the function.
	Name string `json:"name,required"`
	// The output of the function. This will be `null` if the outputs have not been
	// [submitted](https://platform.openai.com/docs/api-reference/runs/submitToolOutputs)
	// yet.
	Output string                       `json:"output,required,nullable"`
	JSON   functionToolCallFunctionJSON `json:"-"`
}

// functionToolCallFunctionJSON contains the JSON metadata for the struct
// [FunctionToolCallFunction]
type functionToolCallFunctionJSON struct {
	Arguments   apijson.Field
	Name        apijson.Field
	Output      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *FunctionToolCallFunction) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r functionToolCallFunctionJSON) RawJSON() string {
	return r.raw
}

// The type of tool call. This is always going to be `function` for this type of
// tool call.
type FunctionToolCallType string

const (
	FunctionToolCallTypeFunction FunctionToolCallType = "function"
)

func (r FunctionToolCallType) IsKnown() bool {
	switch r {
	case FunctionToolCallTypeFunction:
		return true
	}
	return false
}

type FunctionToolCallDelta struct {
	// The index of the tool call in the tool calls array.
	Index int64 `json:"index,required"`
	// The type of tool call. This is always going to be `function` for this type of
	// tool call.
	Type FunctionToolCallDeltaType `json:"type,required"`
	// The ID of the tool call object.
	ID string `json:"id"`
	// The definition of the function that was called.
	Function FunctionToolCallDeltaFunction `json:"function"`
	JSON     functionToolCallDeltaJSON     `json:"-"`
}

// functionToolCallDeltaJSON contains the JSON metadata for the struct
// [FunctionToolCallDelta]
type functionToolCallDeltaJSON struct {
	Index       apijson.Field
	Type        apijson.Field
	ID          apijson.Field
	Function    apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *FunctionToolCallDelta) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r functionToolCallDeltaJSON) RawJSON() string {
	return r.raw
}

func (r FunctionToolCallDelta) implementsToolCallDelta() {}

// The type of tool call. This is always going to be `function` for this type of
// tool call.
type FunctionToolCallDeltaType string

const (
	FunctionToolCallDeltaTypeFunction FunctionToolCallDeltaType = "function"
)

func (r FunctionToolCallDeltaType) IsKnown() bool {
	switch r {
	case FunctionToolCallDeltaTypeFunction:
		return true
	}
	return false
}

// The definition of the function that was called.
type FunctionToolCallDeltaFunction struct {
	// The arguments passed to the function.
	Arguments string `json:"arguments"`
	// The name of the function.
	Name string `json:"name"`
	// The output of the function. This will be `null` if the outputs have not been
	// [submitted](https://platform.openai.com/docs/api-reference/runs/submitToolOutputs)
	// yet.
	Output string                            `json:"output,nullable"`
	JSON   functionToolCallDeltaFunctionJSON `json:"-"`
}

// functionToolCallDeltaFunctionJSON contains the JSON metadata for the struct
// [FunctionToolCallDeltaFunction]
type functionToolCallDeltaFunctionJSON struct {
	Arguments   apijson.Field
	Name        apijson.Field
	Output      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *FunctionToolCallDeltaFunction) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r functionToolCallDeltaFunctionJSON) RawJSON() string {
	return r.raw
}

// Details of the message creation by the run step.
type MessageCreationStepDetails struct {
	MessageCreation MessageCreationStepDetailsMessageCreation `json:"message_creation,required"`
	// Always `message_creation`.
	Type MessageCreationStepDetailsType `json:"type,required"`
	JSON messageCreationStepDetailsJSON `json:"-"`
}

// messageCreationStepDetailsJSON contains the JSON metadata for the struct
// [MessageCreationStepDetails]
type messageCreationStepDetailsJSON struct {
	MessageCreation apijson.Field
	Type            apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r *MessageCreationStepDetails) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r messageCreationStepDetailsJSON) RawJSON() string {
	return r.raw
}

func (r MessageCreationStepDetails) implementsRunStepStepDetails() {}

type MessageCreationStepDetailsMessageCreation struct {
	// The ID of the message that was created by this run step.
	MessageID string                                        `json:"message_id,required"`
	JSON      messageCreationStepDetailsMessageCreationJSON `json:"-"`
}

// messageCreationStepDetailsMessageCreationJSON contains the JSON metadata for the
// struct [MessageCreationStepDetailsMessageCreation]
type messageCreationStepDetailsMessageCreationJSON struct {
	MessageID   apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *MessageCreationStepDetailsMessageCreation) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r messageCreationStepDetailsMessageCreationJSON) RawJSON() string {
	return r.raw
}

// Always `message_creation`.
type MessageCreationStepDetailsType string

const (
	MessageCreationStepDetailsTypeMessageCreation MessageCreationStepDetailsType = "message_creation"
)

func (r MessageCreationStepDetailsType) IsKnown() bool {
	switch r {
	case MessageCreationStepDetailsTypeMessageCreation:
		return true
	}
	return false
}

// Represents a step in execution of a run.
type RunStep struct {
	// The identifier of the run step, which can be referenced in API endpoints.
	ID string `json:"id,required"`
	// The ID of the
	// [assistant](https://platform.openai.com/docs/api-reference/assistants)
	// associated with the run step.
	AssistantID string `json:"assistant_id,required"`
	// The Unix timestamp (in seconds) for when the run step was cancelled.
	CancelledAt int64 `json:"cancelled_at,required,nullable"`
	// The Unix timestamp (in seconds) for when the run step completed.
	CompletedAt int64 `json:"completed_at,required,nullable"`
	// The Unix timestamp (in seconds) for when the run step was created.
	CreatedAt int64 `json:"created_at,required"`
	// The Unix timestamp (in seconds) for when the run step expired. A step is
	// considered expired if the parent run is expired.
	ExpiredAt int64 `json:"expired_at,required,nullable"`
	// The Unix timestamp (in seconds) for when the run step failed.
	FailedAt int64 `json:"failed_at,required,nullable"`
	// The last error associated with this run step. Will be `null` if there are no
	// errors.
	LastError RunStepLastError `json:"last_error,required,nullable"`
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard.
	//
	// Keys are strings with a maximum length of 64 characters. Values are strings with
	// a maximum length of 512 characters.
	Metadata shared.Metadata `json:"metadata,required,nullable"`
	// The object type, which is always `thread.run.step`.
	Object RunStepObject `json:"object,required"`
	// The ID of the [run](https://platform.openai.com/docs/api-reference/runs) that
	// this run step is a part of.
	RunID string `json:"run_id,required"`
	// The status of the run step, which can be either `in_progress`, `cancelled`,
	// `failed`, `completed`, or `expired`.
	Status RunStepStatus `json:"status,required"`
	// The details of the run step.
	StepDetails RunStepStepDetails `json:"step_details,required"`
	// The ID of the [thread](https://platform.openai.com/docs/api-reference/threads)
	// that was run.
	ThreadID string `json:"thread_id,required"`
	// The type of run step, which can be either `message_creation` or `tool_calls`.
	Type RunStepType `json:"type,required"`
	// Usage statistics related to the run step. This value will be `null` while the
	// run step's status is `in_progress`.
	Usage RunStepUsage `json:"usage,required,nullable"`
	JSON  runStepJSON  `json:"-"`
}

// runStepJSON contains the JSON metadata for the struct [RunStep]
type runStepJSON struct {
	ID          apijson.Field
	AssistantID apijson.Field
	CancelledAt apijson.Field
	CompletedAt apijson.Field
	CreatedAt   apijson.Field
	ExpiredAt   apijson.Field
	FailedAt    apijson.Field
	LastError   apijson.Field
	Metadata    apijson.Field
	Object      apijson.Field
	RunID       apijson.Field
	Status      apijson.Field
	StepDetails apijson.Field
	ThreadID    apijson.Field
	Type        apijson.Field
	Usage       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *RunStep) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r runStepJSON) RawJSON() string {
	return r.raw
}

// The last error associated with this run step. Will be `null` if there are no
// errors.
type RunStepLastError struct {
	// One of `server_error` or `rate_limit_exceeded`.
	Code RunStepLastErrorCode `json:"code,required"`
	// A human-readable description of the error.
	Message string               `json:"message,required"`
	JSON    runStepLastErrorJSON `json:"-"`
}

// runStepLastErrorJSON contains the JSON metadata for the struct
// [RunStepLastError]
type runStepLastErrorJSON struct {
	Code        apijson.Field
	Message     apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *RunStepLastError) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r runStepLastErrorJSON) RawJSON() string {
	return r.raw
}

// One of `server_error` or `rate_limit_exceeded`.
type RunStepLastErrorCode string

const (
	RunStepLastErrorCodeServerError       RunStepLastErrorCode = "server_error"
	RunStepLastErrorCodeRateLimitExceeded RunStepLastErrorCode = "rate_limit_exceeded"
)

func (r RunStepLastErrorCode) IsKnown() bool {
	switch r {
	case RunStepLastErrorCodeServerError, RunStepLastErrorCodeRateLimitExceeded:
		return true
	}
	return false
}

// The object type, which is always `thread.run.step`.
type RunStepObject string

const (
	RunStepObjectThreadRunStep RunStepObject = "thread.run.step"
)

func (r RunStepObject) IsKnown() bool {
	switch r {
	case RunStepObjectThreadRunStep:
		return true
	}
	return false
}

// The status of the run step, which can be either `in_progress`, `cancelled`,
// `failed`, `completed`, or `expired`.
type RunStepStatus string

const (
	RunStepStatusInProgress RunStepStatus = "in_progress"
	RunStepStatusCancelled  RunStepStatus = "cancelled"
	RunStepStatusFailed     RunStepStatus = "failed"
	RunStepStatusCompleted  RunStepStatus = "completed"
	RunStepStatusExpired    RunStepStatus = "expired"
)

func (r RunStepStatus) IsKnown() bool {
	switch r {
	case RunStepStatusInProgress, RunStepStatusCancelled, RunStepStatusFailed, RunStepStatusCompleted, RunStepStatusExpired:
		return true
	}
	return false
}

// The details of the run step.
type RunStepStepDetails struct {
	// Always `message_creation`.
	Type RunStepStepDetailsType `json:"type,required"`
	// This field can have the runtime type of
	// [MessageCreationStepDetailsMessageCreation].
	MessageCreation interface{} `json:"message_creation"`
	// This field can have the runtime type of [[]ToolCall].
	ToolCalls interface{}            `json:"tool_calls"`
	JSON      runStepStepDetailsJSON `json:"-"`
	union     RunStepStepDetailsUnion
}

// runStepStepDetailsJSON contains the JSON metadata for the struct
// [RunStepStepDetails]
type runStepStepDetailsJSON struct {
	Type            apijson.Field
	MessageCreation apijson.Field
	ToolCalls       apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r runStepStepDetailsJSON) RawJSON() string {
	return r.raw
}

func (r *RunStepStepDetails) UnmarshalJSON(data []byte) (err error) {
	*r = RunStepStepDetails{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [RunStepStepDetailsUnion] interface which you can cast to the
// specific types for more type safety.
//
// Possible runtime types of the union are [MessageCreationStepDetails],
// [ToolCallsStepDetails].
func (r RunStepStepDetails) AsUnion() RunStepStepDetailsUnion {
	return r.union
}

// The details of the run step.
//
// Union satisfied by [MessageCreationStepDetails] or [ToolCallsStepDetails].
type RunStepStepDetailsUnion interface {
	implementsRunStepStepDetails()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*RunStepStepDetailsUnion)(nil)).Elem(),
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(MessageCreationStepDetails{}),
			DiscriminatorValue: "message_creation",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(ToolCallsStepDetails{}),
			DiscriminatorValue: "tool_calls",
		},
	)
}

// Always `message_creation`.
type RunStepStepDetailsType string

const (
	RunStepStepDetailsTypeMessageCreation RunStepStepDetailsType = "message_creation"
	RunStepStepDetailsTypeToolCalls       RunStepStepDetailsType = "tool_calls"
)

func (r RunStepStepDetailsType) IsKnown() bool {
	switch r {
	case RunStepStepDetailsTypeMessageCreation, RunStepStepDetailsTypeToolCalls:
		return true
	}
	return false
}

// The type of run step, which can be either `message_creation` or `tool_calls`.
type RunStepType string

const (
	RunStepTypeMessageCreation RunStepType = "message_creation"
	RunStepTypeToolCalls       RunStepType = "tool_calls"
)

func (r RunStepType) IsKnown() bool {
	switch r {
	case RunStepTypeMessageCreation, RunStepTypeToolCalls:
		return true
	}
	return false
}

// Usage statistics related to the run step. This value will be `null` while the
// run step's status is `in_progress`.
type RunStepUsage struct {
	// Number of completion tokens used over the course of the run step.
	CompletionTokens int64 `json:"completion_tokens,required"`
	// Number of prompt tokens used over the course of the run step.
	PromptTokens int64 `json:"prompt_tokens,required"`
	// Total number of tokens used (prompt + completion).
	TotalTokens int64            `json:"total_tokens,required"`
	JSON        runStepUsageJSON `json:"-"`
}

// runStepUsageJSON contains the JSON metadata for the struct [RunStepUsage]
type runStepUsageJSON struct {
	CompletionTokens apijson.Field
	PromptTokens     apijson.Field
	TotalTokens      apijson.Field
	raw              string
	ExtraFields      map[string]apijson.Field
}

func (r *RunStepUsage) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r runStepUsageJSON) RawJSON() string {
	return r.raw
}

// The delta containing the fields that have changed on the run step.
type RunStepDelta struct {
	// The details of the run step.
	StepDetails RunStepDeltaStepDetails `json:"step_details"`
	JSON        runStepDeltaJSON        `json:"-"`
}

// runStepDeltaJSON contains the JSON metadata for the struct [RunStepDelta]
type runStepDeltaJSON struct {
	StepDetails apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *RunStepDelta) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r runStepDeltaJSON) RawJSON() string {
	return r.raw
}

// The details of the run step.
type RunStepDeltaStepDetails struct {
	// Always `message_creation`.
	Type RunStepDeltaStepDetailsType `json:"type,required"`
	// This field can have the runtime type of
	// [RunStepDeltaMessageDeltaMessageCreation].
	MessageCreation interface{} `json:"message_creation"`
	// This field can have the runtime type of [[]ToolCallDelta].
	ToolCalls interface{}                 `json:"tool_calls"`
	JSON      runStepDeltaStepDetailsJSON `json:"-"`
	union     RunStepDeltaStepDetailsUnion
}

// runStepDeltaStepDetailsJSON contains the JSON metadata for the struct
// [RunStepDeltaStepDetails]
type runStepDeltaStepDetailsJSON struct {
	Type            apijson.Field
	MessageCreation apijson.Field
	ToolCalls       apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r runStepDeltaStepDetailsJSON) RawJSON() string {
	return r.raw
}

func (r *RunStepDeltaStepDetails) UnmarshalJSON(data []byte) (err error) {
	*r = RunStepDeltaStepDetails{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [RunStepDeltaStepDetailsUnion] interface which you can cast to
// the specific types for more type safety.
//
// Possible runtime types of the union are [RunStepDeltaMessageDelta],
// [ToolCallDeltaObject].
func (r RunStepDeltaStepDetails) AsUnion() RunStepDeltaStepDetailsUnion {
	return r.union
}

// The details of the run step.
//
// Union satisfied by [RunStepDeltaMessageDelta] or [ToolCallDeltaObject].
type RunStepDeltaStepDetailsUnion interface {
	implementsRunStepDeltaStepDetails()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*RunStepDeltaStepDetailsUnion)(nil)).Elem(),
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(RunStepDeltaMessageDelta{}),
			DiscriminatorValue: "message_creation",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(ToolCallDeltaObject{}),
			DiscriminatorValue: "tool_calls",
		},
	)
}

// Always `message_creation`.
type RunStepDeltaStepDetailsType string

const (
	RunStepDeltaStepDetailsTypeMessageCreation RunStepDeltaStepDetailsType = "message_creation"
	RunStepDeltaStepDetailsTypeToolCalls       RunStepDeltaStepDetailsType = "tool_calls"
)

func (r RunStepDeltaStepDetailsType) IsKnown() bool {
	switch r {
	case RunStepDeltaStepDetailsTypeMessageCreation, RunStepDeltaStepDetailsTypeToolCalls:
		return true
	}
	return false
}

// Represents a run step delta i.e. any changed fields on a run step during
// streaming.
type RunStepDeltaEvent struct {
	// The identifier of the run step, which can be referenced in API endpoints.
	ID string `json:"id,required"`
	// The delta containing the fields that have changed on the run step.
	Delta RunStepDelta `json:"delta,required"`
	// The object type, which is always `thread.run.step.delta`.
	Object RunStepDeltaEventObject `json:"object,required"`
	JSON   runStepDeltaEventJSON   `json:"-"`
}

// runStepDeltaEventJSON contains the JSON metadata for the struct
// [RunStepDeltaEvent]
type runStepDeltaEventJSON struct {
	ID          apijson.Field
	Delta       apijson.Field
	Object      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *RunStepDeltaEvent) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r runStepDeltaEventJSON) RawJSON() string {
	return r.raw
}

// The object type, which is always `thread.run.step.delta`.
type RunStepDeltaEventObject string

const (
	RunStepDeltaEventObjectThreadRunStepDelta RunStepDeltaEventObject = "thread.run.step.delta"
)

func (r RunStepDeltaEventObject) IsKnown() bool {
	switch r {
	case RunStepDeltaEventObjectThreadRunStepDelta:
		return true
	}
	return false
}

// Details of the message creation by the run step.
type RunStepDeltaMessageDelta struct {
	// Always `message_creation`.
	Type            RunStepDeltaMessageDeltaType            `json:"type,required"`
	MessageCreation RunStepDeltaMessageDeltaMessageCreation `json:"message_creation"`
	JSON            runStepDeltaMessageDeltaJSON            `json:"-"`
}

// runStepDeltaMessageDeltaJSON contains the JSON metadata for the struct
// [RunStepDeltaMessageDelta]
type runStepDeltaMessageDeltaJSON struct {
	Type            apijson.Field
	MessageCreation apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r *RunStepDeltaMessageDelta) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r runStepDeltaMessageDeltaJSON) RawJSON() string {
	return r.raw
}

func (r RunStepDeltaMessageDelta) implementsRunStepDeltaStepDetails() {}

// Always `message_creation`.
type RunStepDeltaMessageDeltaType string

const (
	RunStepDeltaMessageDeltaTypeMessageCreation RunStepDeltaMessageDeltaType = "message_creation"
)

func (r RunStepDeltaMessageDeltaType) IsKnown() bool {
	switch r {
	case RunStepDeltaMessageDeltaTypeMessageCreation:
		return true
	}
	return false
}

type RunStepDeltaMessageDeltaMessageCreation struct {
	// The ID of the message that was created by this run step.
	MessageID string                                      `json:"message_id"`
	JSON      runStepDeltaMessageDeltaMessageCreationJSON `json:"-"`
}

// runStepDeltaMessageDeltaMessageCreationJSON contains the JSON metadata for the
// struct [RunStepDeltaMessageDeltaMessageCreation]
type runStepDeltaMessageDeltaMessageCreationJSON struct {
	MessageID   apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *RunStepDeltaMessageDeltaMessageCreation) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r runStepDeltaMessageDeltaMessageCreationJSON) RawJSON() string {
	return r.raw
}

type RunStepInclude string

const (
	RunStepIncludeStepDetailsToolCallsFileSearchResultsContent RunStepInclude = "step_details.tool_calls[*].file_search.results[*].content"
)

func (r RunStepInclude) IsKnown() bool {
	switch r {
	case RunStepIncludeStepDetailsToolCallsFileSearchResultsContent:
		return true
	}
	return false
}

// Details of the Code Interpreter tool call the run step was involved in.
type ToolCall struct {
	// The ID of the tool call.
	ID string `json:"id,required"`
	// The type of tool call. This is always going to be `code_interpreter` for this
	// type of tool call.
	Type ToolCallType `json:"type,required"`
	// This field can have the runtime type of
	// [CodeInterpreterToolCallCodeInterpreter].
	CodeInterpreter interface{} `json:"code_interpreter"`
	// This field can have the runtime type of [FileSearchToolCallFileSearch].
	FileSearch interface{} `json:"file_search"`
	// This field can have the runtime type of [FunctionToolCallFunction].
	Function interface{}  `json:"function"`
	JSON     toolCallJSON `json:"-"`
	union    ToolCallUnion
}

// toolCallJSON contains the JSON metadata for the struct [ToolCall]
type toolCallJSON struct {
	ID              apijson.Field
	Type            apijson.Field
	CodeInterpreter apijson.Field
	FileSearch      apijson.Field
	Function        apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r toolCallJSON) RawJSON() string {
	return r.raw
}

func (r *ToolCall) UnmarshalJSON(data []byte) (err error) {
	*r = ToolCall{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [ToolCallUnion] interface which you can cast to the specific
// types for more type safety.
//
// Possible runtime types of the union are [CodeInterpreterToolCall],
// [FileSearchToolCall], [FunctionToolCall].
func (r ToolCall) AsUnion() ToolCallUnion {
	return r.union
}

// Details of the Code Interpreter tool call the run step was involved in.
//
// Union satisfied by [CodeInterpreterToolCall], [FileSearchToolCall] or
// [FunctionToolCall].
type ToolCallUnion interface {
	implementsToolCall()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*ToolCallUnion)(nil)).Elem(),
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(CodeInterpreterToolCall{}),
			DiscriminatorValue: "code_interpreter",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(FileSearchToolCall{}),
			DiscriminatorValue: "file_search",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(FunctionToolCall{}),
			DiscriminatorValue: "function",
		},
	)
}

// The type of tool call. This is always going to be `code_interpreter` for this
// type of tool call.
type ToolCallType string

const (
	ToolCallTypeCodeInterpreter ToolCallType = "code_interpreter"
	ToolCallTypeFileSearch      ToolCallType = "file_search"
	ToolCallTypeFunction        ToolCallType = "function"
)

func (r ToolCallType) IsKnown() bool {
	switch r {
	case ToolCallTypeCodeInterpreter, ToolCallTypeFileSearch, ToolCallTypeFunction:
		return true
	}
	return false
}

// Details of the Code Interpreter tool call the run step was involved in.
type ToolCallDelta struct {
	// The index of the tool call in the tool calls array.
	Index int64 `json:"index,required"`
	// The type of tool call. This is always going to be `code_interpreter` for this
	// type of tool call.
	Type ToolCallDeltaType `json:"type,required"`
	// The ID of the tool call.
	ID string `json:"id"`
	// This field can have the runtime type of
	// [CodeInterpreterToolCallDeltaCodeInterpreter].
	CodeInterpreter interface{} `json:"code_interpreter"`
	// This field can have the runtime type of [interface{}].
	FileSearch interface{} `json:"file_search"`
	// This field can have the runtime type of [FunctionToolCallDeltaFunction].
	Function interface{}       `json:"function"`
	JSON     toolCallDeltaJSON `json:"-"`
	union    ToolCallDeltaUnion
}

// toolCallDeltaJSON contains the JSON metadata for the struct [ToolCallDelta]
type toolCallDeltaJSON struct {
	Index           apijson.Field
	Type            apijson.Field
	ID              apijson.Field
	CodeInterpreter apijson.Field
	FileSearch      apijson.Field
	Function        apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r toolCallDeltaJSON) RawJSON() string {
	return r.raw
}

func (r *ToolCallDelta) UnmarshalJSON(data []byte) (err error) {
	*r = ToolCallDelta{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [ToolCallDeltaUnion] interface which you can cast to the
// specific types for more type safety.
//
// Possible runtime types of the union are [CodeInterpreterToolCallDelta],
// [FileSearchToolCallDelta], [FunctionToolCallDelta].
func (r ToolCallDelta) AsUnion() ToolCallDeltaUnion {
	return r.union
}

// Details of the Code Interpreter tool call the run step was involved in.
//
// Union satisfied by [CodeInterpreterToolCallDelta], [FileSearchToolCallDelta] or
// [FunctionToolCallDelta].
type ToolCallDeltaUnion interface {
	implementsToolCallDelta()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*ToolCallDeltaUnion)(nil)).Elem(),
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(CodeInterpreterToolCallDelta{}),
			DiscriminatorValue: "code_interpreter",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(FileSearchToolCallDelta{}),
			DiscriminatorValue: "file_search",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(FunctionToolCallDelta{}),
			DiscriminatorValue: "function",
		},
	)
}

// The type of tool call. This is always going to be `code_interpreter` for this
// type of tool call.
type ToolCallDeltaType string

const (
	ToolCallDeltaTypeCodeInterpreter ToolCallDeltaType = "code_interpreter"
	ToolCallDeltaTypeFileSearch      ToolCallDeltaType = "file_search"
	ToolCallDeltaTypeFunction        ToolCallDeltaType = "function"
)

func (r ToolCallDeltaType) IsKnown() bool {
	switch r {
	case ToolCallDeltaTypeCodeInterpreter, ToolCallDeltaTypeFileSearch, ToolCallDeltaTypeFunction:
		return true
	}
	return false
}

// Details of the tool call.
type ToolCallDeltaObject struct {
	// Always `tool_calls`.
	Type ToolCallDeltaObjectType `json:"type,required"`
	// An array of tool calls the run step was involved in. These can be associated
	// with one of three types of tools: `code_interpreter`, `file_search`, or
	// `function`.
	ToolCalls []ToolCallDelta         `json:"tool_calls"`
	JSON      toolCallDeltaObjectJSON `json:"-"`
}

// toolCallDeltaObjectJSON contains the JSON metadata for the struct
// [ToolCallDeltaObject]
type toolCallDeltaObjectJSON struct {
	Type        apijson.Field
	ToolCalls   apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ToolCallDeltaObject) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r toolCallDeltaObjectJSON) RawJSON() string {
	return r.raw
}

func (r ToolCallDeltaObject) implementsRunStepDeltaStepDetails() {}

// Always `tool_calls`.
type ToolCallDeltaObjectType string

const (
	ToolCallDeltaObjectTypeToolCalls ToolCallDeltaObjectType = "tool_calls"
)

func (r ToolCallDeltaObjectType) IsKnown() bool {
	switch r {
	case ToolCallDeltaObjectTypeToolCalls:
		return true
	}
	return false
}

// Details of the tool call.
type ToolCallsStepDetails struct {
	// An array of tool calls the run step was involved in. These can be associated
	// with one of three types of tools: `code_interpreter`, `file_search`, or
	// `function`.
	ToolCalls []ToolCall `json:"tool_calls,required"`
	// Always `tool_calls`.
	Type ToolCallsStepDetailsType `json:"type,required"`
	JSON toolCallsStepDetailsJSON `json:"-"`
}

// toolCallsStepDetailsJSON contains the JSON metadata for the struct
// [ToolCallsStepDetails]
type toolCallsStepDetailsJSON struct {
	ToolCalls   apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ToolCallsStepDetails) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r toolCallsStepDetailsJSON) RawJSON() string {
	return r.raw
}

func (r ToolCallsStepDetails) implementsRunStepStepDetails() {}

// Always `tool_calls`.
type ToolCallsStepDetailsType string

const (
	ToolCallsStepDetailsTypeToolCalls ToolCallsStepDetailsType = "tool_calls"
)

func (r ToolCallsStepDetailsType) IsKnown() bool {
	switch r {
	case ToolCallsStepDetailsTypeToolCalls:
		return true
	}
	return false
}

type BetaThreadRunStepGetParams struct {
	// A list of additional fields to include in the response. Currently the only
	// supported value is `step_details.tool_calls[*].file_search.results[*].content`
	// to fetch the file search result content.
	//
	// See the
	// [file search tool documentation](https://platform.openai.com/docs/assistants/tools/file-search#customizing-file-search-settings)
	// for more information.
	Include param.Field[[]RunStepInclude] `query:"include"`
}

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
	After param.Field[string] `query:"after"`
	// A cursor for use in pagination. `before` is an object ID that defines your place
	// in the list. For instance, if you make a list request and receive 100 objects,
	// starting with obj_foo, your subsequent call can include before=obj_foo in order
	// to fetch the previous page of the list.
	Before param.Field[string] `query:"before"`
	// A list of additional fields to include in the response. Currently the only
	// supported value is `step_details.tool_calls[*].file_search.results[*].content`
	// to fetch the file search result content.
	//
	// See the
	// [file search tool documentation](https://platform.openai.com/docs/assistants/tools/file-search#customizing-file-search-settings)
	// for more information.
	Include param.Field[[]RunStepInclude] `query:"include"`
	// A limit on the number of objects to be returned. Limit can range between 1 and
	// 100, and the default is 20.
	Limit param.Field[int64] `query:"limit"`
	// Sort order by the `created_at` timestamp of the objects. `asc` for ascending
	// order and `desc` for descending order.
	Order param.Field[BetaThreadRunStepListParamsOrder] `query:"order"`
}

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

func (r BetaThreadRunStepListParamsOrder) IsKnown() bool {
	switch r {
	case BetaThreadRunStepListParamsOrderAsc, BetaThreadRunStepListParamsOrderDesc:
		return true
	}
	return false
}

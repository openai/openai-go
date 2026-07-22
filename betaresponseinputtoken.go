// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai

import (
	"context"
	"fmt"
	"net/http"
	"slices"

	"github.com/openai/openai-go/v3/internal/apijson"
	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/packages/respjson"
	"github.com/openai/openai-go/v3/shared/constant"
)

// BetaResponseInputTokenService contains methods and other services that help with
// interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaResponseInputTokenService] method instead.
type BetaResponseInputTokenService struct {
	Options []option.RequestOption
}

// NewBetaResponseInputTokenService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewBetaResponseInputTokenService(opts ...option.RequestOption) (r BetaResponseInputTokenService) {
	r = BetaResponseInputTokenService{}
	r.Options = opts
	return
}

// Returns input token counts of the request.
//
// Returns an object with `object` set to `response.input_tokens` and an
// `input_tokens` count.
func (r *BetaResponseInputTokenService) Count(ctx context.Context, params BetaResponseInputTokenCountParams, opts ...option.RequestOption) (res *BetaResponseInputTokenCountResponse, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("openai-beta", fmt.Sprintf("%v", v)))
	}
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	path := "responses/input_tokens?beta=true"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

type BetaResponseInputTokenCountResponse struct {
	InputTokens int64                        `json:"input_tokens" api:"required"`
	Object      constant.ResponseInputTokens `json:"object" default:"response.input_tokens"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		InputTokens respjson.Field
		Object      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaResponseInputTokenCountResponse) RawJSON() string { return r.JSON.raw }
func (r *BetaResponseInputTokenCountResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaResponseInputTokenCountParams struct {
	// A system (or developer) message inserted into the model's context. When used
	// along with `previous_response_id`, the instructions from a previous response
	// will not be carried over to the next response. This makes it simple to swap out
	// system (or developer) messages in new responses.
	Instructions param.Opt[string] `json:"instructions,omitzero"`
	// Model ID used to generate the response, like `gpt-4o` or `o3`. OpenAI offers a
	// wide range of models with different capabilities, performance characteristics,
	// and price points. Refer to the
	// [model guide](https://platform.openai.com/docs/models) to browse and compare
	// available models.
	Model param.Opt[string] `json:"model,omitzero"`
	// Whether to allow the model to run tool calls in parallel.
	ParallelToolCalls param.Opt[bool] `json:"parallel_tool_calls,omitzero"`
	// The unique ID of the previous response to the model. Use this to create
	// multi-turn conversations. Learn more about
	// [conversation state](https://platform.openai.com/docs/guides/conversation-state).
	// Cannot be used in conjunction with `conversation`.
	PreviousResponseID param.Opt[string] `json:"previous_response_id,omitzero"`
	// The conversation that this response belongs to. Items from this conversation are
	// prepended to `input_items` for this response request. Input items and output
	// items from this response are automatically added to this conversation after this
	// response completes.
	Conversation BetaResponseInputTokenCountParamsConversationUnion `json:"conversation,omitzero"`
	// Text, image, or file inputs to the model, used to generate a response
	Input BetaResponseInputTokenCountParamsInputUnion `json:"input,omitzero"`
	// **gpt-5 and o-series models only** Configuration options for
	// [reasoning models](https://platform.openai.com/docs/guides/reasoning).
	Reasoning BetaResponseInputTokenCountParamsReasoning `json:"reasoning,omitzero"`
	// Configuration options for a text response from the model. Can be plain text or
	// structured JSON data. Learn more:
	//
	// - [Text inputs and outputs](https://platform.openai.com/docs/guides/text)
	// - [Structured Outputs](https://platform.openai.com/docs/guides/structured-outputs)
	Text BetaResponseInputTokenCountParamsText `json:"text,omitzero"`
	// Controls which tool the model should use, if any.
	ToolChoice BetaResponseInputTokenCountParamsToolChoiceUnion `json:"tool_choice,omitzero"`
	// An array of tools the model may call while generating a response. You can
	// specify which tool to use by setting the `tool_choice` parameter.
	Tools []BetaToolUnionParam `json:"tools,omitzero"`
	// A model-owned style preset to apply to this request. Omit this parameter to use
	// the model's default style. Supported values may expand over time. Values must be
	// at most 64 characters.
	Personality BetaResponseInputTokenCountParamsPersonality `json:"personality,omitzero"`
	// The truncation strategy to use for the model response. - `auto`: If the input to
	// this Response exceeds the model's context window size, the model will truncate
	// the response to fit the context window by dropping items from the beginning of
	// the conversation. - `disabled` (default): If the input size will exceed the
	// context window size for a model, the request will fail with a 400 error.
	//
	// Any of "auto", "disabled".
	Truncation BetaResponseInputTokenCountParamsTruncation `json:"truncation,omitzero"`
	// Any of "responses_multi_agent=v1".
	Betas []string `header:"openai-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaResponseInputTokenCountParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaResponseInputTokenCountParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaResponseInputTokenCountParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaResponseInputTokenCountParamsConversationUnion struct {
	OfString             param.Opt[string]              `json:",omitzero,inline"`
	OfConversationObject *BetaResponseConversationParam `json:",omitzero,inline"`
	paramUnion
}

func (u BetaResponseInputTokenCountParamsConversationUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfString, u.OfConversationObject)
}
func (u *BetaResponseInputTokenCountParamsConversationUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaResponseInputTokenCountParamsConversationUnion) asAny() any {
	if !param.IsOmitted(u.OfString) {
		return &u.OfString.Value
	} else if !param.IsOmitted(u.OfConversationObject) {
		return u.OfConversationObject
	}
	return nil
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaResponseInputTokenCountParamsInputUnion struct {
	OfString                     param.Opt[string]                 `json:",omitzero,inline"`
	OfBetaResponseInputItemArray []BetaResponseInputItemUnionParam `json:",omitzero,inline"`
	paramUnion
}

func (u BetaResponseInputTokenCountParamsInputUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfString, u.OfBetaResponseInputItemArray)
}
func (u *BetaResponseInputTokenCountParamsInputUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaResponseInputTokenCountParamsInputUnion) asAny() any {
	if !param.IsOmitted(u.OfString) {
		return &u.OfString.Value
	} else if !param.IsOmitted(u.OfBetaResponseInputItemArray) {
		return &u.OfBetaResponseInputItemArray
	}
	return nil
}

// A model-owned style preset to apply to this request. Omit this parameter to use
// the model's default style. Supported values may expand over time. Values must be
// at most 64 characters.
type BetaResponseInputTokenCountParamsPersonality string

const (
	BetaResponseInputTokenCountParamsPersonalityFriendly  BetaResponseInputTokenCountParamsPersonality = "friendly"
	BetaResponseInputTokenCountParamsPersonalityPragmatic BetaResponseInputTokenCountParamsPersonality = "pragmatic"
)

// **gpt-5 and o-series models only** Configuration options for
// [reasoning models](https://platform.openai.com/docs/guides/reasoning).
type BetaResponseInputTokenCountParamsReasoning struct {
	// Controls which reasoning items are rendered back to the model on later turns. If
	// omitted or set to `auto`, the model determines the context mode. The `gpt-5.6`
	// model family defaults to `all_turns`; earlier models default to `current_turn`.
	//
	// When returned on a response, this is the effective reasoning context mode used
	// for the response.
	//
	// Any of "auto", "current_turn", "all_turns".
	Context string `json:"context,omitzero"`
	// Constrains effort on reasoning for reasoning models. Currently supported values
	// are `none`, `minimal`, `low`, `medium`, `high`, `xhigh`, and `max`. Reducing
	// reasoning effort can result in faster responses and fewer tokens used on
	// reasoning in a response. Not all reasoning models support every value. See the
	// [reasoning guide](https://platform.openai.com/docs/guides/reasoning) for
	// model-specific support.
	//
	// Any of "none", "minimal", "low", "medium", "high", "xhigh", "max".
	Effort string `json:"effort,omitzero"`
	// **Deprecated:** use `summary` instead.
	//
	// A summary of the reasoning performed by the model. This can be useful for
	// debugging and understanding the model's reasoning process. One of `auto`,
	// `concise`, or `detailed`.
	//
	// Any of "auto", "concise", "detailed".
	//
	// Deprecated: deprecated
	GenerateSummary string `json:"generate_summary,omitzero"`
	// A summary of the reasoning performed by the model. This can be useful for
	// debugging and understanding the model's reasoning process. One of `auto`,
	// `concise`, or `detailed`.
	//
	// `concise` is supported for `computer-use-preview` models and all reasoning
	// models after `gpt-5`.
	//
	// Any of "auto", "concise", "detailed".
	Summary string `json:"summary,omitzero"`
	// Controls the reasoning execution mode for the request.
	//
	// When returned on a response, this is the effective execution mode.
	Mode string `json:"mode,omitzero"`
	paramObj
}

func (r BetaResponseInputTokenCountParamsReasoning) MarshalJSON() (data []byte, err error) {
	type shadow BetaResponseInputTokenCountParamsReasoning
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaResponseInputTokenCountParamsReasoning) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[BetaResponseInputTokenCountParamsReasoning](
		"context", "auto", "current_turn", "all_turns",
	)
	apijson.RegisterFieldValidator[BetaResponseInputTokenCountParamsReasoning](
		"effort", "none", "minimal", "low", "medium", "high", "xhigh", "max",
	)
	apijson.RegisterFieldValidator[BetaResponseInputTokenCountParamsReasoning](
		"generate_summary", "auto", "concise", "detailed",
	)
	apijson.RegisterFieldValidator[BetaResponseInputTokenCountParamsReasoning](
		"summary", "auto", "concise", "detailed",
	)
}

// Configuration options for a text response from the model. Can be plain text or
// structured JSON data. Learn more:
//
// - [Text inputs and outputs](https://platform.openai.com/docs/guides/text)
// - [Structured Outputs](https://platform.openai.com/docs/guides/structured-outputs)
type BetaResponseInputTokenCountParamsText struct {
	// Constrains the verbosity of the model's response. Lower values will result in
	// more concise responses, while higher values will result in more verbose
	// responses. Currently supported values are `low`, `medium`, and `high`. The
	// default is `medium`.
	//
	// Any of "low", "medium", "high".
	Verbosity string `json:"verbosity,omitzero"`
	// An object specifying the format that the model must output.
	//
	// Configuring `{ "type": "json_schema" }` enables Structured Outputs, which
	// ensures the model will match your supplied JSON schema. Learn more in the
	// [Structured Outputs guide](https://platform.openai.com/docs/guides/structured-outputs).
	//
	// The default format is `{ "type": "text" }` with no additional options.
	//
	// **Not recommended for gpt-4o and newer models:**
	//
	// Setting to `{ "type": "json_object" }` enables the older JSON mode, which
	// ensures the message the model generates is valid JSON. Using `json_schema` is
	// preferred for models that support it.
	Format BetaResponseFormatTextConfigUnionParam `json:"format,omitzero"`
	paramObj
}

func (r BetaResponseInputTokenCountParamsText) MarshalJSON() (data []byte, err error) {
	type shadow BetaResponseInputTokenCountParamsText
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaResponseInputTokenCountParamsText) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[BetaResponseInputTokenCountParamsText](
		"verbosity", "low", "medium", "high",
	)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaResponseInputTokenCountParamsToolChoiceUnion struct {
	// Check if union is this variant with !param.IsOmitted(union.OfToolChoiceMode)
	OfToolChoiceMode                                                                 param.Opt[BetaToolChoiceOptions]                                                     `json:",omitzero,inline"`
	OfAllowedTools                                                                   *BetaToolChoiceAllowedParam                                                          `json:",omitzero,inline"`
	OfHostedTool                                                                     *BetaToolChoiceTypesParam                                                            `json:",omitzero,inline"`
	OfFunctionTool                                                                   *BetaToolChoiceFunctionParam                                                         `json:",omitzero,inline"`
	OfMcpTool                                                                        *BetaToolChoiceMcpParam                                                              `json:",omitzero,inline"`
	OfCustomTool                                                                     *BetaToolChoiceCustomParam                                                           `json:",omitzero,inline"`
	OfBetaResponseInputTokenCountsToolChoiceBetaSpecificProgrammaticToolCallingParam *BetaResponseInputTokenCountParamsToolChoiceBetaSpecificProgrammaticToolCallingParam `json:",omitzero,inline"`
	OfSpecificApplyPatchToolChoice                                                   *BetaToolChoiceApplyPatchParam                                                       `json:",omitzero,inline"`
	OfSpecificShellToolChoice                                                        *BetaToolChoiceShellParam                                                            `json:",omitzero,inline"`
	paramUnion
}

func (u BetaResponseInputTokenCountParamsToolChoiceUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfToolChoiceMode,
		u.OfAllowedTools,
		u.OfHostedTool,
		u.OfFunctionTool,
		u.OfMcpTool,
		u.OfCustomTool,
		u.OfBetaResponseInputTokenCountsToolChoiceBetaSpecificProgrammaticToolCallingParam,
		u.OfSpecificApplyPatchToolChoice,
		u.OfSpecificShellToolChoice)
}
func (u *BetaResponseInputTokenCountParamsToolChoiceUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaResponseInputTokenCountParamsToolChoiceUnion) asAny() any {
	if !param.IsOmitted(u.OfToolChoiceMode) {
		return &u.OfToolChoiceMode
	} else if !param.IsOmitted(u.OfAllowedTools) {
		return u.OfAllowedTools
	} else if !param.IsOmitted(u.OfHostedTool) {
		return u.OfHostedTool
	} else if !param.IsOmitted(u.OfFunctionTool) {
		return u.OfFunctionTool
	} else if !param.IsOmitted(u.OfMcpTool) {
		return u.OfMcpTool
	} else if !param.IsOmitted(u.OfCustomTool) {
		return u.OfCustomTool
	} else if !param.IsOmitted(u.OfBetaResponseInputTokenCountsToolChoiceBetaSpecificProgrammaticToolCallingParam) {
		return u.OfBetaResponseInputTokenCountsToolChoiceBetaSpecificProgrammaticToolCallingParam
	} else if !param.IsOmitted(u.OfSpecificApplyPatchToolChoice) {
		return u.OfSpecificApplyPatchToolChoice
	} else if !param.IsOmitted(u.OfSpecificShellToolChoice) {
		return u.OfSpecificShellToolChoice
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaResponseInputTokenCountParamsToolChoiceUnion) GetMode() *string {
	if vt := u.OfAllowedTools; vt != nil {
		return (*string)(&vt.Mode)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaResponseInputTokenCountParamsToolChoiceUnion) GetTools() []map[string]any {
	if vt := u.OfAllowedTools; vt != nil {
		return vt.Tools
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaResponseInputTokenCountParamsToolChoiceUnion) GetServerLabel() *string {
	if vt := u.OfMcpTool; vt != nil {
		return &vt.ServerLabel
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaResponseInputTokenCountParamsToolChoiceUnion) GetType() *string {
	if vt := u.OfAllowedTools; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfHostedTool; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfFunctionTool; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfMcpTool; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfCustomTool; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfBetaResponseInputTokenCountsToolChoiceBetaSpecificProgrammaticToolCallingParam; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfSpecificApplyPatchToolChoice; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfSpecificShellToolChoice; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaResponseInputTokenCountParamsToolChoiceUnion) GetName() *string {
	if vt := u.OfFunctionTool; vt != nil {
		return (*string)(&vt.Name)
	} else if vt := u.OfMcpTool; vt != nil && vt.Name.Valid() {
		return &vt.Name.Value
	} else if vt := u.OfCustomTool; vt != nil {
		return (*string)(&vt.Name)
	}
	return nil
}

func NewBetaResponseInputTokenCountParamsToolChoiceBetaSpecificProgrammaticToolCallingParam() BetaResponseInputTokenCountParamsToolChoiceBetaSpecificProgrammaticToolCallingParam {
	return BetaResponseInputTokenCountParamsToolChoiceBetaSpecificProgrammaticToolCallingParam{
		Type: "programmatic_tool_calling",
	}
}

// This struct has a constant value, construct it with
// [NewBetaResponseInputTokenCountParamsToolChoiceBetaSpecificProgrammaticToolCallingParam].
type BetaResponseInputTokenCountParamsToolChoiceBetaSpecificProgrammaticToolCallingParam struct {
	// The tool to call. Always `programmatic_tool_calling`.
	Type constant.ProgrammaticToolCalling `json:"type" default:"programmatic_tool_calling"`
	paramObj
}

func (r BetaResponseInputTokenCountParamsToolChoiceBetaSpecificProgrammaticToolCallingParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaResponseInputTokenCountParamsToolChoiceBetaSpecificProgrammaticToolCallingParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaResponseInputTokenCountParamsToolChoiceBetaSpecificProgrammaticToolCallingParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The truncation strategy to use for the model response. - `auto`: If the input to
// this Response exceeds the model's context window size, the model will truncate
// the response to fit the context window by dropping items from the beginning of
// the conversation. - `disabled` (default): If the input size will exceed the
// context window size for a model, the request will fail with a 400 error.
type BetaResponseInputTokenCountParamsTruncation string

const (
	BetaResponseInputTokenCountParamsTruncationAuto     BetaResponseInputTokenCountParamsTruncation = "auto"
	BetaResponseInputTokenCountParamsTruncationDisabled BetaResponseInputTokenCountParamsTruncation = "disabled"
)

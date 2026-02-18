// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package responses

import (
	"context"
	"net/http"
	"slices"

	"github.com/openai/openai-go/v3/internal/apijson"
	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/packages/respjson"
	"github.com/openai/openai-go/v3/shared"
	"github.com/openai/openai-go/v3/shared/constant"
)

// InputTokenService contains methods and other services that help with interacting
// with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewInputTokenService] method instead.
type InputTokenService struct {
	Options []option.RequestOption
}

// NewInputTokenService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewInputTokenService(opts ...option.RequestOption) (r InputTokenService) {
	r = InputTokenService{}
	r.Options = opts
	return
}

// Returns input token counts of the request.
//
// Returns an object with `object` set to `response.input_tokens` and an
// `input_tokens` count.
func (r *InputTokenService) Count(ctx context.Context, body InputTokenCountParams, opts ...option.RequestOption) (res *InputTokenCountResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "responses/input_tokens"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

type InputTokenCountResponse struct {
	InputTokens int64                        `json:"input_tokens,required"`
	Object      constant.ResponseInputTokens `json:"object,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		InputTokens respjson.Field
		Object      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r InputTokenCountResponse) RawJSON() string { return r.JSON.raw }
func (r *InputTokenCountResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type InputTokenCountParams struct {
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
	Conversation InputTokenCountParamsConversationUnion `json:"conversation,omitzero"`
	// Text, image, or file inputs to the model, used to generate a response
	Input InputTokenCountParamsInputUnion `json:"input,omitzero"`
	// Configuration options for a text response from the model. Can be plain text or
	// structured JSON data. Learn more:
	//
	// - [Text inputs and outputs](https://platform.openai.com/docs/guides/text)
	// - [Structured Outputs](https://platform.openai.com/docs/guides/structured-outputs)
	Text InputTokenCountParamsText `json:"text,omitzero"`
	// Controls which tool the model should use, if any.
	ToolChoice InputTokenCountParamsToolChoiceUnion `json:"tool_choice,omitzero"`
	// An array of tools the model may call while generating a response. You can
	// specify which tool to use by setting the `tool_choice` parameter.
	Tools []ToolUnionParam `json:"tools,omitzero"`
	// **gpt-5 and o-series models only** Configuration options for
	// [reasoning models](https://platform.openai.com/docs/guides/reasoning).
	Reasoning shared.ReasoningParam `json:"reasoning,omitzero"`
	// The truncation strategy to use for the model response. - `auto`: If the input to
	// this Response exceeds the model's context window size, the model will truncate
	// the response to fit the context window by dropping items from the beginning of
	// the conversation. - `disabled` (default): If the input size will exceed the
	// context window size for a model, the request will fail with a 400 error.
	//
	// Any of "auto", "disabled".
	Truncation InputTokenCountParamsTruncation `json:"truncation,omitzero"`
	paramObj
}

func (r InputTokenCountParams) MarshalJSON() (data []byte, err error) {
	type shadow InputTokenCountParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *InputTokenCountParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type InputTokenCountParamsConversationUnion struct {
	OfString             param.Opt[string]          `json:",omitzero,inline"`
	OfConversationObject *ResponseConversationParam `json:",omitzero,inline"`
	paramUnion
}

func (u InputTokenCountParamsConversationUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfString, u.OfConversationObject)
}
func (u *InputTokenCountParamsConversationUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *InputTokenCountParamsConversationUnion) asAny() any {
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
type InputTokenCountParamsInputUnion struct {
	OfString                 param.Opt[string]             `json:",omitzero,inline"`
	OfResponseInputItemArray []ResponseInputItemUnionParam `json:",omitzero,inline"`
	paramUnion
}

func (u InputTokenCountParamsInputUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfString, u.OfResponseInputItemArray)
}
func (u *InputTokenCountParamsInputUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *InputTokenCountParamsInputUnion) asAny() any {
	if !param.IsOmitted(u.OfString) {
		return &u.OfString.Value
	} else if !param.IsOmitted(u.OfResponseInputItemArray) {
		return &u.OfResponseInputItemArray
	}
	return nil
}

// Configuration options for a text response from the model. Can be plain text or
// structured JSON data. Learn more:
//
// - [Text inputs and outputs](https://platform.openai.com/docs/guides/text)
// - [Structured Outputs](https://platform.openai.com/docs/guides/structured-outputs)
type InputTokenCountParamsText struct {
	// Constrains the verbosity of the model's response. Lower values will result in
	// more concise responses, while higher values will result in more verbose
	// responses. Currently supported values are `low`, `medium`, and `high`.
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
	Format ResponseFormatTextConfigUnionParam `json:"format,omitzero"`
	paramObj
}

func (r InputTokenCountParamsText) MarshalJSON() (data []byte, err error) {
	type shadow InputTokenCountParamsText
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *InputTokenCountParamsText) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[InputTokenCountParamsText](
		"verbosity", "low", "medium", "high",
	)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type InputTokenCountParamsToolChoiceUnion struct {
	// Check if union is this variant with !param.IsOmitted(union.OfToolChoiceMode)
	OfToolChoiceMode               param.Opt[ToolChoiceOptions] `json:",omitzero,inline"`
	OfAllowedTools                 *ToolChoiceAllowedParam      `json:",omitzero,inline"`
	OfHostedTool                   *ToolChoiceTypesParam        `json:",omitzero,inline"`
	OfFunctionTool                 *ToolChoiceFunctionParam     `json:",omitzero,inline"`
	OfMcpTool                      *ToolChoiceMcpParam          `json:",omitzero,inline"`
	OfCustomTool                   *ToolChoiceCustomParam       `json:",omitzero,inline"`
	OfSpecificApplyPatchToolChoice *ToolChoiceApplyPatchParam   `json:",omitzero,inline"`
	OfSpecificShellToolChoice      *ToolChoiceShellParam        `json:",omitzero,inline"`
	paramUnion
}

func (u InputTokenCountParamsToolChoiceUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfToolChoiceMode,
		u.OfAllowedTools,
		u.OfHostedTool,
		u.OfFunctionTool,
		u.OfMcpTool,
		u.OfCustomTool,
		u.OfSpecificApplyPatchToolChoice,
		u.OfSpecificShellToolChoice)
}
func (u *InputTokenCountParamsToolChoiceUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *InputTokenCountParamsToolChoiceUnion) asAny() any {
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
	} else if !param.IsOmitted(u.OfSpecificApplyPatchToolChoice) {
		return u.OfSpecificApplyPatchToolChoice
	} else if !param.IsOmitted(u.OfSpecificShellToolChoice) {
		return u.OfSpecificShellToolChoice
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u InputTokenCountParamsToolChoiceUnion) GetMode() *string {
	if vt := u.OfAllowedTools; vt != nil {
		return (*string)(&vt.Mode)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u InputTokenCountParamsToolChoiceUnion) GetTools() []map[string]any {
	if vt := u.OfAllowedTools; vt != nil {
		return vt.Tools
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u InputTokenCountParamsToolChoiceUnion) GetServerLabel() *string {
	if vt := u.OfMcpTool; vt != nil {
		return &vt.ServerLabel
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u InputTokenCountParamsToolChoiceUnion) GetType() *string {
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
	} else if vt := u.OfSpecificApplyPatchToolChoice; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfSpecificShellToolChoice; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u InputTokenCountParamsToolChoiceUnion) GetName() *string {
	if vt := u.OfFunctionTool; vt != nil {
		return (*string)(&vt.Name)
	} else if vt := u.OfMcpTool; vt != nil && vt.Name.Valid() {
		return &vt.Name.Value
	} else if vt := u.OfCustomTool; vt != nil {
		return (*string)(&vt.Name)
	}
	return nil
}

// The truncation strategy to use for the model response. - `auto`: If the input to
// this Response exceeds the model's context window size, the model will truncate
// the response to fit the context window by dropping items from the beginning of
// the conversation. - `disabled` (default): If the input size will exceed the
// context window size for a model, the request will fail with a 400 error.
type InputTokenCountParamsTruncation string

const (
	InputTokenCountParamsTruncationAuto     InputTokenCountParamsTruncation = "auto"
	InputTokenCountParamsTruncationDisabled InputTokenCountParamsTruncation = "disabled"
)

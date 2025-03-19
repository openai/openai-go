// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package shared

import (
	"github.com/openai/openai-go/internal/apijson"
	"github.com/openai/openai-go/internal/param"
)

type ChatModel = string

const (
	ChatModelO3Mini                           ChatModel = "o3-mini"
	ChatModelO3Mini2025_01_31                 ChatModel = "o3-mini-2025-01-31"
	ChatModelO1                               ChatModel = "o1"
	ChatModelO1_2024_12_17                    ChatModel = "o1-2024-12-17"
	ChatModelO1Preview                        ChatModel = "o1-preview"
	ChatModelO1Preview2024_09_12              ChatModel = "o1-preview-2024-09-12"
	ChatModelO1Mini                           ChatModel = "o1-mini"
	ChatModelO1Mini2024_09_12                 ChatModel = "o1-mini-2024-09-12"
	ChatModelGPT4o                            ChatModel = "gpt-4o"
	ChatModelGPT4o2024_11_20                  ChatModel = "gpt-4o-2024-11-20"
	ChatModelGPT4o2024_08_06                  ChatModel = "gpt-4o-2024-08-06"
	ChatModelGPT4o2024_05_13                  ChatModel = "gpt-4o-2024-05-13"
	ChatModelGPT4oAudioPreview                ChatModel = "gpt-4o-audio-preview"
	ChatModelGPT4oAudioPreview2024_10_01      ChatModel = "gpt-4o-audio-preview-2024-10-01"
	ChatModelGPT4oAudioPreview2024_12_17      ChatModel = "gpt-4o-audio-preview-2024-12-17"
	ChatModelGPT4oMiniAudioPreview            ChatModel = "gpt-4o-mini-audio-preview"
	ChatModelGPT4oMiniAudioPreview2024_12_17  ChatModel = "gpt-4o-mini-audio-preview-2024-12-17"
	ChatModelGPT4oSearchPreview               ChatModel = "gpt-4o-search-preview"
	ChatModelGPT4oMiniSearchPreview           ChatModel = "gpt-4o-mini-search-preview"
	ChatModelGPT4oSearchPreview2025_03_11     ChatModel = "gpt-4o-search-preview-2025-03-11"
	ChatModelGPT4oMiniSearchPreview2025_03_11 ChatModel = "gpt-4o-mini-search-preview-2025-03-11"
	ChatModelChatgpt4oLatest                  ChatModel = "chatgpt-4o-latest"
	ChatModelGPT4oMini                        ChatModel = "gpt-4o-mini"
	ChatModelGPT4oMini2024_07_18              ChatModel = "gpt-4o-mini-2024-07-18"
	ChatModelGPT4Turbo                        ChatModel = "gpt-4-turbo"
	ChatModelGPT4Turbo2024_04_09              ChatModel = "gpt-4-turbo-2024-04-09"
	ChatModelGPT4_0125Preview                 ChatModel = "gpt-4-0125-preview"
	ChatModelGPT4TurboPreview                 ChatModel = "gpt-4-turbo-preview"
	ChatModelGPT4_1106Preview                 ChatModel = "gpt-4-1106-preview"
	ChatModelGPT4VisionPreview                ChatModel = "gpt-4-vision-preview"
	ChatModelGPT4                             ChatModel = "gpt-4"
	ChatModelGPT4_0314                        ChatModel = "gpt-4-0314"
	ChatModelGPT4_0613                        ChatModel = "gpt-4-0613"
	ChatModelGPT4_32k                         ChatModel = "gpt-4-32k"
	ChatModelGPT4_32k0314                     ChatModel = "gpt-4-32k-0314"
	ChatModelGPT4_32k0613                     ChatModel = "gpt-4-32k-0613"
	ChatModelGPT3_5Turbo                      ChatModel = "gpt-3.5-turbo"
	ChatModelGPT3_5Turbo16k                   ChatModel = "gpt-3.5-turbo-16k"
	ChatModelGPT3_5Turbo0301                  ChatModel = "gpt-3.5-turbo-0301"
	ChatModelGPT3_5Turbo0613                  ChatModel = "gpt-3.5-turbo-0613"
	ChatModelGPT3_5Turbo1106                  ChatModel = "gpt-3.5-turbo-1106"
	ChatModelGPT3_5Turbo0125                  ChatModel = "gpt-3.5-turbo-0125"
	ChatModelGPT3_5Turbo16k0613               ChatModel = "gpt-3.5-turbo-16k-0613"
)

// A filter used to compare a specified attribute key to a given value using a
// defined comparison operation.
type ComparisonFilterParam struct {
	// The key to compare against the value.
	Key param.Field[string] `json:"key,required"`
	// Specifies the comparison operator: `eq`, `ne`, `gt`, `gte`, `lt`, `lte`.
	//
	// - `eq`: equals
	// - `ne`: not equal
	// - `gt`: greater than
	// - `gte`: greater than or equal
	// - `lt`: less than
	// - `lte`: less than or equal
	Type param.Field[ComparisonFilterType] `json:"type,required"`
	// The value to compare against the attribute key; supports string, number, or
	// boolean types.
	Value param.Field[ComparisonFilterValueUnionParam] `json:"value,required"`
}

func (r ComparisonFilterParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ComparisonFilterParam) ImplementsVectorStoreSearchParamsFiltersUnion() {}

// Specifies the comparison operator: `eq`, `ne`, `gt`, `gte`, `lt`, `lte`.
//
// - `eq`: equals
// - `ne`: not equal
// - `gt`: greater than
// - `gte`: greater than or equal
// - `lt`: less than
// - `lte`: less than or equal
type ComparisonFilterType string

const (
	ComparisonFilterTypeEq  ComparisonFilterType = "eq"
	ComparisonFilterTypeNe  ComparisonFilterType = "ne"
	ComparisonFilterTypeGt  ComparisonFilterType = "gt"
	ComparisonFilterTypeGte ComparisonFilterType = "gte"
	ComparisonFilterTypeLt  ComparisonFilterType = "lt"
	ComparisonFilterTypeLte ComparisonFilterType = "lte"
)

func (r ComparisonFilterType) IsKnown() bool {
	switch r {
	case ComparisonFilterTypeEq, ComparisonFilterTypeNe, ComparisonFilterTypeGt, ComparisonFilterTypeGte, ComparisonFilterTypeLt, ComparisonFilterTypeLte:
		return true
	}
	return false
}

// The value to compare against the attribute key; supports string, number, or
// boolean types.
//
// Satisfied by [shared.UnionString], [shared.UnionFloat], [shared.UnionBool].
type ComparisonFilterValueUnionParam interface {
	ImplementsComparisonFilterValueUnionParam()
}

// Combine multiple filters using `and` or `or`.
type CompoundFilterParam struct {
	// Array of filters to combine. Items can be `ComparisonFilter` or
	// `CompoundFilter`.
	Filters param.Field[[]ComparisonFilterParam] `json:"filters,required"`
	// Type of operation: `and` or `or`.
	Type param.Field[CompoundFilterType] `json:"type,required"`
}

func (r CompoundFilterParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r CompoundFilterParam) ImplementsVectorStoreSearchParamsFiltersUnion() {}

// Type of operation: `and` or `or`.
type CompoundFilterType string

const (
	CompoundFilterTypeAnd CompoundFilterType = "and"
	CompoundFilterTypeOr  CompoundFilterType = "or"
)

func (r CompoundFilterType) IsKnown() bool {
	switch r {
	case CompoundFilterTypeAnd, CompoundFilterTypeOr:
		return true
	}
	return false
}

type ErrorObject struct {
	Code    string          `json:"code,required,nullable"`
	Message string          `json:"message,required"`
	Param   string          `json:"param,required,nullable"`
	Type    string          `json:"type,required"`
	JSON    errorObjectJSON `json:"-"`
}

// errorObjectJSON contains the JSON metadata for the struct [ErrorObject]
type errorObjectJSON struct {
	Code        apijson.Field
	Message     apijson.Field
	Param       apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ErrorObject) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r errorObjectJSON) RawJSON() string {
	return r.raw
}

type FunctionDefinition struct {
	// The name of the function to be called. Must be a-z, A-Z, 0-9, or contain
	// underscores and dashes, with a maximum length of 64.
	Name string `json:"name,required"`
	// A description of what the function does, used by the model to choose when and
	// how to call the function.
	Description string `json:"description"`
	// The parameters the functions accepts, described as a JSON Schema object. See the
	// [guide](https://platform.openai.com/docs/guides/function-calling) for examples,
	// and the
	// [JSON Schema reference](https://json-schema.org/understanding-json-schema/) for
	// documentation about the format.
	//
	// Omitting `parameters` defines a function with an empty parameter list.
	Parameters FunctionParameters `json:"parameters"`
	// Whether to enable strict schema adherence when generating the function call. If
	// set to true, the model will follow the exact schema defined in the `parameters`
	// field. Only a subset of JSON Schema is supported when `strict` is `true`. Learn
	// more about Structured Outputs in the
	// [function calling guide](docs/guides/function-calling).
	Strict bool                   `json:"strict,nullable"`
	JSON   functionDefinitionJSON `json:"-"`
}

// functionDefinitionJSON contains the JSON metadata for the struct
// [FunctionDefinition]
type functionDefinitionJSON struct {
	Name        apijson.Field
	Description apijson.Field
	Parameters  apijson.Field
	Strict      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *FunctionDefinition) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r functionDefinitionJSON) RawJSON() string {
	return r.raw
}

type FunctionDefinitionParam struct {
	// The name of the function to be called. Must be a-z, A-Z, 0-9, or contain
	// underscores and dashes, with a maximum length of 64.
	Name param.Field[string] `json:"name,required"`
	// A description of what the function does, used by the model to choose when and
	// how to call the function.
	Description param.Field[string] `json:"description"`
	// The parameters the functions accepts, described as a JSON Schema object. See the
	// [guide](https://platform.openai.com/docs/guides/function-calling) for examples,
	// and the
	// [JSON Schema reference](https://json-schema.org/understanding-json-schema/) for
	// documentation about the format.
	//
	// Omitting `parameters` defines a function with an empty parameter list.
	Parameters param.Field[FunctionParameters] `json:"parameters"`
	// Whether to enable strict schema adherence when generating the function call. If
	// set to true, the model will follow the exact schema defined in the `parameters`
	// field. Only a subset of JSON Schema is supported when `strict` is `true`. Learn
	// more about Structured Outputs in the
	// [function calling guide](docs/guides/function-calling).
	Strict param.Field[bool] `json:"strict"`
}

func (r FunctionDefinitionParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type FunctionParameters map[string]interface{}

type Metadata map[string]string

type MetadataParam map[string]string

// **o-series models only**
//
// Constrains effort on reasoning for
// [reasoning models](https://platform.openai.com/docs/guides/reasoning). Currently
// supported values are `low`, `medium`, and `high`. Reducing reasoning effort can
// result in faster responses and fewer tokens used on reasoning in a response.
type ReasoningEffort string

const (
	ReasoningEffortLow    ReasoningEffort = "low"
	ReasoningEffortMedium ReasoningEffort = "medium"
	ReasoningEffortHigh   ReasoningEffort = "high"
)

func (r ReasoningEffort) IsKnown() bool {
	switch r {
	case ReasoningEffortLow, ReasoningEffortMedium, ReasoningEffortHigh:
		return true
	}
	return false
}

// JSON object response format. An older method of generating JSON responses. Using
// `json_schema` is recommended for models that support it. Note that the model
// will not generate JSON without a system or user message instructing it to do so.
type ResponseFormatJSONObject struct {
	// The type of response format being defined. Always `json_object`.
	Type ResponseFormatJSONObjectType `json:"type,required"`
	JSON responseFormatJSONObjectJSON `json:"-"`
}

// responseFormatJSONObjectJSON contains the JSON metadata for the struct
// [ResponseFormatJSONObject]
type responseFormatJSONObjectJSON struct {
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ResponseFormatJSONObject) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r responseFormatJSONObjectJSON) RawJSON() string {
	return r.raw
}

func (r ResponseFormatJSONObject) ImplementsAssistantResponseFormatOptionUnion() {}

// The type of response format being defined. Always `json_object`.
type ResponseFormatJSONObjectType string

const (
	ResponseFormatJSONObjectTypeJSONObject ResponseFormatJSONObjectType = "json_object"
)

func (r ResponseFormatJSONObjectType) IsKnown() bool {
	switch r {
	case ResponseFormatJSONObjectTypeJSONObject:
		return true
	}
	return false
}

// JSON object response format. An older method of generating JSON responses. Using
// `json_schema` is recommended for models that support it. Note that the model
// will not generate JSON without a system or user message instructing it to do so.
type ResponseFormatJSONObjectParam struct {
	// The type of response format being defined. Always `json_object`.
	Type param.Field[ResponseFormatJSONObjectType] `json:"type,required"`
}

func (r ResponseFormatJSONObjectParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ResponseFormatJSONObjectParam) ImplementsChatCompletionNewParamsResponseFormatUnion() {}

func (r ResponseFormatJSONObjectParam) ImplementsAssistantResponseFormatOptionUnionParam() {}

// JSON Schema response format. Used to generate structured JSON responses. Learn
// more about
// [Structured Outputs](https://platform.openai.com/docs/guides/structured-outputs).
type ResponseFormatJSONSchema struct {
	// Structured Outputs configuration options, including a JSON Schema.
	JSONSchema ResponseFormatJSONSchemaJSONSchema `json:"json_schema,required"`
	// The type of response format being defined. Always `json_schema`.
	Type ResponseFormatJSONSchemaType `json:"type,required"`
	JSON responseFormatJSONSchemaJSON `json:"-"`
}

// responseFormatJSONSchemaJSON contains the JSON metadata for the struct
// [ResponseFormatJSONSchema]
type responseFormatJSONSchemaJSON struct {
	JSONSchema  apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ResponseFormatJSONSchema) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r responseFormatJSONSchemaJSON) RawJSON() string {
	return r.raw
}

func (r ResponseFormatJSONSchema) ImplementsAssistantResponseFormatOptionUnion() {}

// Structured Outputs configuration options, including a JSON Schema.
type ResponseFormatJSONSchemaJSONSchema struct {
	// The name of the response format. Must be a-z, A-Z, 0-9, or contain underscores
	// and dashes, with a maximum length of 64.
	Name string `json:"name,required"`
	// A description of what the response format is for, used by the model to determine
	// how to respond in the format.
	Description string `json:"description"`
	// The schema for the response format, described as a JSON Schema object. Learn how
	// to build JSON schemas [here](https://json-schema.org/).
	Schema map[string]interface{} `json:"schema"`
	// Whether to enable strict schema adherence when generating the output. If set to
	// true, the model will always follow the exact schema defined in the `schema`
	// field. Only a subset of JSON Schema is supported when `strict` is `true`. To
	// learn more, read the
	// [Structured Outputs guide](https://platform.openai.com/docs/guides/structured-outputs).
	Strict bool                                   `json:"strict,nullable"`
	JSON   responseFormatJSONSchemaJSONSchemaJSON `json:"-"`
}

// responseFormatJSONSchemaJSONSchemaJSON contains the JSON metadata for the struct
// [ResponseFormatJSONSchemaJSONSchema]
type responseFormatJSONSchemaJSONSchemaJSON struct {
	Name        apijson.Field
	Description apijson.Field
	Schema      apijson.Field
	Strict      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ResponseFormatJSONSchemaJSONSchema) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r responseFormatJSONSchemaJSONSchemaJSON) RawJSON() string {
	return r.raw
}

// The type of response format being defined. Always `json_schema`.
type ResponseFormatJSONSchemaType string

const (
	ResponseFormatJSONSchemaTypeJSONSchema ResponseFormatJSONSchemaType = "json_schema"
)

func (r ResponseFormatJSONSchemaType) IsKnown() bool {
	switch r {
	case ResponseFormatJSONSchemaTypeJSONSchema:
		return true
	}
	return false
}

// JSON Schema response format. Used to generate structured JSON responses. Learn
// more about
// [Structured Outputs](https://platform.openai.com/docs/guides/structured-outputs).
type ResponseFormatJSONSchemaParam struct {
	// Structured Outputs configuration options, including a JSON Schema.
	JSONSchema param.Field[ResponseFormatJSONSchemaJSONSchemaParam] `json:"json_schema,required"`
	// The type of response format being defined. Always `json_schema`.
	Type param.Field[ResponseFormatJSONSchemaType] `json:"type,required"`
}

func (r ResponseFormatJSONSchemaParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ResponseFormatJSONSchemaParam) ImplementsChatCompletionNewParamsResponseFormatUnion() {}

func (r ResponseFormatJSONSchemaParam) ImplementsAssistantResponseFormatOptionUnionParam() {}

// Structured Outputs configuration options, including a JSON Schema.
type ResponseFormatJSONSchemaJSONSchemaParam struct {
	// The name of the response format. Must be a-z, A-Z, 0-9, or contain underscores
	// and dashes, with a maximum length of 64.
	Name param.Field[string] `json:"name,required"`
	// A description of what the response format is for, used by the model to determine
	// how to respond in the format.
	Description param.Field[string] `json:"description"`
	// The schema for the response format, described as a JSON Schema object. Learn how
	// to build JSON schemas [here](https://json-schema.org/).
	Schema param.Field[interface{}] `json:"schema"`
	// Whether to enable strict schema adherence when generating the output. If set to
	// true, the model will always follow the exact schema defined in the `schema`
	// field. Only a subset of JSON Schema is supported when `strict` is `true`. To
	// learn more, read the
	// [Structured Outputs guide](https://platform.openai.com/docs/guides/structured-outputs).
	Strict param.Field[bool] `json:"strict"`
}

func (r ResponseFormatJSONSchemaJSONSchemaParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Default response format. Used to generate text responses.
type ResponseFormatText struct {
	// The type of response format being defined. Always `text`.
	Type ResponseFormatTextType `json:"type,required"`
	JSON responseFormatTextJSON `json:"-"`
}

// responseFormatTextJSON contains the JSON metadata for the struct
// [ResponseFormatText]
type responseFormatTextJSON struct {
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ResponseFormatText) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r responseFormatTextJSON) RawJSON() string {
	return r.raw
}

func (r ResponseFormatText) ImplementsAssistantResponseFormatOptionUnion() {}

// The type of response format being defined. Always `text`.
type ResponseFormatTextType string

const (
	ResponseFormatTextTypeText ResponseFormatTextType = "text"
)

func (r ResponseFormatTextType) IsKnown() bool {
	switch r {
	case ResponseFormatTextTypeText:
		return true
	}
	return false
}

// Default response format. Used to generate text responses.
type ResponseFormatTextParam struct {
	// The type of response format being defined. Always `text`.
	Type param.Field[ResponseFormatTextType] `json:"type,required"`
}

func (r ResponseFormatTextParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ResponseFormatTextParam) ImplementsChatCompletionNewParamsResponseFormatUnion() {}

func (r ResponseFormatTextParam) ImplementsAssistantResponseFormatOptionUnionParam() {}

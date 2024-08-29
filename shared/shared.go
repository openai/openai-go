// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package shared

import (
	"github.com/openai/openai-go/internal/apijson"
	"github.com/openai/openai-go/internal/param"
)

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

type ResponseFormatJSONObjectParam struct {
	// The type of response format being defined: `json_object`
	Type param.Field[ResponseFormatJSONObjectType] `json:"type,required"`
}

func (r ResponseFormatJSONObjectParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ResponseFormatJSONObjectParam) ImplementsChatCompletionNewParamsResponseFormatUnion() {}

// The type of response format being defined: `json_object`
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

type ResponseFormatJSONSchemaParam struct {
	JSONSchema param.Field[ResponseFormatJSONSchemaJSONSchemaParam] `json:"json_schema,required"`
	// The type of response format being defined: `json_schema`
	Type param.Field[ResponseFormatJSONSchemaType] `json:"type,required"`
}

func (r ResponseFormatJSONSchemaParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ResponseFormatJSONSchemaParam) ImplementsChatCompletionNewParamsResponseFormatUnion() {}

type ResponseFormatJSONSchemaJSONSchemaParam struct {
	// The name of the response format. Must be a-z, A-Z, 0-9, or contain underscores
	// and dashes, with a maximum length of 64.
	Name param.Field[string] `json:"name,required"`
	// A description of what the response format is for, used by the model to determine
	// how to respond in the format.
	Description param.Field[string] `json:"description"`
	// The schema for the response format, described as a JSON Schema object.
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

// The type of response format being defined: `json_schema`
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

type ResponseFormatTextParam struct {
	// The type of response format being defined: `text`
	Type param.Field[ResponseFormatTextType] `json:"type,required"`
}

func (r ResponseFormatTextParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ResponseFormatTextParam) ImplementsChatCompletionNewParamsResponseFormatUnion() {}

// The type of response format being defined: `text`
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

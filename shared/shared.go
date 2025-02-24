// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package shared

import (
	"github.com/openai/openai-go/internal/apijson"
	"github.com/openai/openai-go/packages/param"
	"github.com/openai/openai-go/packages/resp"
	"github.com/openai/openai-go/shared/constant"
)

// aliased to make param.APIUnion private when embedding
type apiunion = param.APIUnion

// aliased to make param.APIObject private when embedding
type apiobject = param.APIObject

type ErrorObject struct {
	Code    string `json:"code,omitzero,required,nullable"`
	Message string `json:"message,omitzero,required"`
	Param   string `json:"param,omitzero,required,nullable"`
	Type    string `json:"type,omitzero,required"`
	JSON    struct {
		Code    resp.Field
		Message resp.Field
		Param   resp.Field
		Type    resp.Field
		raw     string
	} `json:"-"`
}

func (r ErrorObject) RawJSON() string { return r.JSON.raw }
func (r *ErrorObject) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type FunctionDefinition struct {
	// The name of the function to be called. Must be a-z, A-Z, 0-9, or contain
	// underscores and dashes, with a maximum length of 64.
	Name string `json:"name,omitzero,required"`
	// A description of what the function does, used by the model to choose when and
	// how to call the function.
	Description string `json:"description,omitzero"`
	// The parameters the functions accepts, described as a JSON Schema object. See the
	// [guide](https://platform.openai.com/docs/guides/function-calling) for examples,
	// and the
	// [JSON Schema reference](https://json-schema.org/understanding-json-schema/) for
	// documentation about the format.
	//
	// Omitting `parameters` defines a function with an empty parameter list.
	Parameters FunctionParameters `json:"parameters,omitzero"`
	// Whether to enable strict schema adherence when generating the function call. If
	// set to true, the model will follow the exact schema defined in the `parameters`
	// field. Only a subset of JSON Schema is supported when `strict` is `true`. Learn
	// more about Structured Outputs in the
	// [function calling guide](docs/guides/function-calling).
	Strict bool `json:"strict,omitzero,nullable"`
	JSON   struct {
		Name        resp.Field
		Description resp.Field
		Parameters  resp.Field
		Strict      resp.Field
		raw         string
	} `json:"-"`
}

func (r FunctionDefinition) RawJSON() string { return r.JSON.raw }
func (r *FunctionDefinition) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this FunctionDefinition to a FunctionDefinitionParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// FunctionDefinitionParam.IsOverridden()
func (r FunctionDefinition) ToParam() FunctionDefinitionParam {
	return param.Override[FunctionDefinitionParam](r.RawJSON())
}

type FunctionDefinitionParam struct {
	// The name of the function to be called. Must be a-z, A-Z, 0-9, or contain
	// underscores and dashes, with a maximum length of 64.
	Name param.String `json:"name,omitzero,required"`
	// A description of what the function does, used by the model to choose when and
	// how to call the function.
	Description param.String `json:"description,omitzero"`
	// The parameters the functions accepts, described as a JSON Schema object. See the
	// [guide](https://platform.openai.com/docs/guides/function-calling) for examples,
	// and the
	// [JSON Schema reference](https://json-schema.org/understanding-json-schema/) for
	// documentation about the format.
	//
	// Omitting `parameters` defines a function with an empty parameter list.
	Parameters FunctionParameters `json:"parameters,omitzero"`
	// Whether to enable strict schema adherence when generating the function call. If
	// set to true, the model will follow the exact schema defined in the `parameters`
	// field. Only a subset of JSON Schema is supported when `strict` is `true`. Learn
	// more about Structured Outputs in the
	// [function calling guide](docs/guides/function-calling).
	Strict param.Bool `json:"strict,omitzero"`
	apiobject
}

func (f FunctionDefinitionParam) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

func (r FunctionDefinitionParam) MarshalJSON() (data []byte, err error) {
	type shadow FunctionDefinitionParam
	return param.MarshalObject(r, (*shadow)(&r))
}

type FunctionParameters map[string]interface{}

type Metadata map[string]string

type MetadataParam map[string]string

type ResponseFormatJSONObjectParam struct {
	// The type of response format being defined: `json_object`
	//
	// This field can be elided, and will be automatically set as "json_object".
	Type constant.JSONObject `json:"type,required"`
	apiobject
}

func (f ResponseFormatJSONObjectParam) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

func (r ResponseFormatJSONObjectParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponseFormatJSONObjectParam
	return param.MarshalObject(r, (*shadow)(&r))
}

type ResponseFormatJSONSchemaParam struct {
	JSONSchema ResponseFormatJSONSchemaJSONSchemaParam `json:"json_schema,omitzero,required"`
	// The type of response format being defined: `json_schema`
	//
	// This field can be elided, and will be automatically set as "json_schema".
	Type constant.JSONSchema `json:"type,required"`
	apiobject
}

func (f ResponseFormatJSONSchemaParam) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

func (r ResponseFormatJSONSchemaParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponseFormatJSONSchemaParam
	return param.MarshalObject(r, (*shadow)(&r))
}

type ResponseFormatJSONSchemaJSONSchemaParam struct {
	// The name of the response format. Must be a-z, A-Z, 0-9, or contain underscores
	// and dashes, with a maximum length of 64.
	Name param.String `json:"name,omitzero,required"`
	// A description of what the response format is for, used by the model to determine
	// how to respond in the format.
	Description param.String `json:"description,omitzero"`
	// The schema for the response format, described as a JSON Schema object.
	Schema map[string]interface{} `json:"schema,omitzero"`
	// Whether to enable strict schema adherence when generating the output. If set to
	// true, the model will always follow the exact schema defined in the `schema`
	// field. Only a subset of JSON Schema is supported when `strict` is `true`. To
	// learn more, read the
	// [Structured Outputs guide](https://platform.openai.com/docs/guides/structured-outputs).
	Strict param.Bool `json:"strict,omitzero"`
	apiobject
}

func (f ResponseFormatJSONSchemaJSONSchemaParam) IsMissing() bool {
	return param.IsOmitted(f) || f.IsNull()
}

func (r ResponseFormatJSONSchemaJSONSchemaParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponseFormatJSONSchemaJSONSchemaParam
	return param.MarshalObject(r, (*shadow)(&r))
}

type ResponseFormatTextParam struct {
	// The type of response format being defined: `text`
	//
	// This field can be elided, and will be automatically set as "text".
	Type constant.Text `json:"type,required"`
	apiobject
}

func (f ResponseFormatTextParam) IsMissing() bool { return param.IsOmitted(f) || f.IsNull() }

func (r ResponseFormatTextParam) MarshalJSON() (data []byte, err error) {
	type shadow ResponseFormatTextParam
	return param.MarshalObject(r, (*shadow)(&r))
}

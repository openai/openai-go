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
	Parameters FunctionParameters     `json:"parameters"`
	JSON       functionDefinitionJSON `json:"-"`
}

// functionDefinitionJSON contains the JSON metadata for the struct
// [FunctionDefinition]
type functionDefinitionJSON struct {
	Name        apijson.Field
	Description apijson.Field
	Parameters  apijson.Field
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
}

func (r FunctionDefinitionParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type FunctionParameters map[string]interface{}

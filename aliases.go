// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai

import (
	"github.com/openai/openai-go/internal/apierror"
	"github.com/openai/openai-go/shared"
)

type Error = apierror.Error

// This is an alias to an internal type.
type ErrorObject = shared.ErrorObject

// This is an alias to an internal type.
type FunctionDefinition = shared.FunctionDefinition

// This is an alias to an internal type.
type FunctionDefinitionParam = shared.FunctionDefinitionParam

// The parameters the functions accepts, described as a JSON Schema object. See the
// [guide](https://platform.openai.com/docs/guides/function-calling) for examples,
// and the
// [JSON Schema reference](https://json-schema.org/understanding-json-schema/) for
// documentation about the format.
//
// Omitting `parameters` defines a function with an empty parameter list.
//
// This is an alias to an internal type.
type FunctionParameters = shared.FunctionParameters

// This is an alias to an internal type.
type ResponseFormatJSONObjectParam = shared.ResponseFormatJSONObjectParam

// The type of response format being defined: `json_object`
//
// This is an alias to an internal type.
type ResponseFormatJSONObjectType = shared.ResponseFormatJSONObjectType

// This is an alias to an internal value.
const ResponseFormatJSONObjectTypeJSONObject = shared.ResponseFormatJSONObjectTypeJSONObject

// This is an alias to an internal type.
type ResponseFormatJSONSchemaParam = shared.ResponseFormatJSONSchemaParam

// This is an alias to an internal type.
type ResponseFormatJSONSchemaJSONSchemaParam = shared.ResponseFormatJSONSchemaJSONSchemaParam

// The type of response format being defined: `json_schema`
//
// This is an alias to an internal type.
type ResponseFormatJSONSchemaType = shared.ResponseFormatJSONSchemaType

// This is an alias to an internal value.
const ResponseFormatJSONSchemaTypeJSONSchema = shared.ResponseFormatJSONSchemaTypeJSONSchema

// This is an alias to an internal type.
type ResponseFormatTextParam = shared.ResponseFormatTextParam

// The type of response format being defined: `text`
//
// This is an alias to an internal type.
type ResponseFormatTextType = shared.ResponseFormatTextType

// This is an alias to an internal value.
const ResponseFormatTextTypeText = shared.ResponseFormatTextTypeText

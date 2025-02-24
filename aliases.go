// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai

import (
	"time"

	"github.com/openai/openai-go/internal/apierror"
	"github.com/openai/openai-go/packages/param"
	"github.com/openai/openai-go/packages/resp"
	"github.com/openai/openai-go/shared"
	"github.com/openai/openai-go/shared/constant"
)

// aliased to make param.APIUnion private when embedding
type apiunion = param.APIUnion

// aliased to make param.APIObject private when embedding
type apiobject = param.APIObject

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

// Set of 16 key-value pairs that can be attached to an object. This can be useful
// for storing additional information about the object in a structured format, and
// querying for objects via API or the dashboard.
//
// Keys are strings with a maximum length of 64 characters. Values are strings with
// a maximum length of 512 characters.
//
// This is an alias to an internal type.
type Metadata = shared.Metadata

// Set of 16 key-value pairs that can be attached to an object. This can be useful
// for storing additional information about the object in a structured format, and
// querying for objects via API or the dashboard.
//
// Keys are strings with a maximum length of 64 characters. Values are strings with
// a maximum length of 512 characters.
//
// This is an alias to an internal type.
type MetadataParam = shared.MetadataParam

// This is an alias to an internal type.
type ResponseFormatJSONObjectParam = shared.ResponseFormatJSONObjectParam

// This is an alias to an internal type.
type ResponseFormatJSONSchemaParam = shared.ResponseFormatJSONSchemaParam

// This is an alias to an internal type.
type ResponseFormatJSONSchemaJSONSchemaParam = shared.ResponseFormatJSONSchemaJSONSchemaParam

// This is an alias to an internal type.
type ResponseFormatTextParam = shared.ResponseFormatTextParam

// Internal helpers for converting from response to param types

func newString(value string) param.String {
	res := param.NeverOmitted[param.String]()
	res.V = value
	return res
}

func newInt(value int64) param.Int {
	res := param.NeverOmitted[param.Int]()
	res.V = value
	return res
}

func newBool(value bool) param.Bool {
	res := param.NeverOmitted[param.Bool]()
	res.V = value
	return res
}

func newFloat(value float64) param.Float {
	res := param.NeverOmitted[param.Float]()
	res.V = value
	return res
}

func newDatetime(value time.Time) param.Datetime {
	res := param.NeverOmitted[param.Datetime]()
	res.V = value
	return res
}

func newDate(value time.Time) param.Date {
	res := param.NeverOmitted[param.Date]()
	res.V = value
	return res
}

func toParamString(value string, meta resp.Field) param.String {
	if !meta.IsMissing() {
		return newString(value)
	}
	return param.String{}
}

func toParamInt(value int64, meta resp.Field) param.Int {
	if !meta.IsMissing() {
		return newInt(value)
	}
	return param.Int{}
}

func toParamBool(value bool, meta resp.Field) param.Bool {
	if !meta.IsMissing() {
		return newBool(value)
	}
	return param.Bool{}
}

func toParamFloat(value float64, meta resp.Field) param.Float {
	if !meta.IsMissing() {
		return newFloat(value)
	}
	return param.Float{}
}

func toParamDatetime(value time.Time, meta resp.Field) param.Datetime {
	if !meta.IsMissing() {
		return newDatetime(value)
	}
	return param.Datetime{}
}

func toParamDate(value time.Time, meta resp.Field) param.Date {
	if !meta.IsMissing() {
		return newDate(value)
	}
	return param.Date{}
}

func ptrToConstant[T constant.Constant[T]](c T) *T {
	if param.IsOmitted(c) {
		c = c.Default()
	}
	return &c
}

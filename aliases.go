// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai

import (
	"github.com/openai/openai-go/internal/apierror"
	"github.com/openai/openai-go/shared"
)

type Error = apierror.Error

// This is an alias to an internal type.
type ChatModel = shared.ChatModel

// This is an alias to an internal value.
const ChatModelO3Mini = shared.ChatModelO3Mini

// This is an alias to an internal value.
const ChatModelO3Mini2025_01_31 = shared.ChatModelO3Mini2025_01_31

// This is an alias to an internal value.
const ChatModelO1 = shared.ChatModelO1

// This is an alias to an internal value.
const ChatModelO1_2024_12_17 = shared.ChatModelO1_2024_12_17

// This is an alias to an internal value.
const ChatModelO1Preview = shared.ChatModelO1Preview

// This is an alias to an internal value.
const ChatModelO1Preview2024_09_12 = shared.ChatModelO1Preview2024_09_12

// This is an alias to an internal value.
const ChatModelO1Mini = shared.ChatModelO1Mini

// This is an alias to an internal value.
const ChatModelO1Mini2024_09_12 = shared.ChatModelO1Mini2024_09_12

// This is an alias to an internal value.
const ChatModelGPT4_5Preview = shared.ChatModelGPT4_5Preview

// This is an alias to an internal value.
const ChatModelGPT4_5Preview2025_02_27 = shared.ChatModelGPT4_5Preview2025_02_27

// This is an alias to an internal value.
const ChatModelGPT4o = shared.ChatModelGPT4o

// This is an alias to an internal value.
const ChatModelGPT4o2024_11_20 = shared.ChatModelGPT4o2024_11_20

// This is an alias to an internal value.
const ChatModelGPT4o2024_08_06 = shared.ChatModelGPT4o2024_08_06

// This is an alias to an internal value.
const ChatModelGPT4o2024_05_13 = shared.ChatModelGPT4o2024_05_13

// This is an alias to an internal value.
const ChatModelGPT4oAudioPreview = shared.ChatModelGPT4oAudioPreview

// This is an alias to an internal value.
const ChatModelGPT4oAudioPreview2024_10_01 = shared.ChatModelGPT4oAudioPreview2024_10_01

// This is an alias to an internal value.
const ChatModelGPT4oAudioPreview2024_12_17 = shared.ChatModelGPT4oAudioPreview2024_12_17

// This is an alias to an internal value.
const ChatModelGPT4oMiniAudioPreview = shared.ChatModelGPT4oMiniAudioPreview

// This is an alias to an internal value.
const ChatModelGPT4oMiniAudioPreview2024_12_17 = shared.ChatModelGPT4oMiniAudioPreview2024_12_17

// This is an alias to an internal value.
const ChatModelChatgpt4oLatest = shared.ChatModelChatgpt4oLatest

// This is an alias to an internal value.
const ChatModelGPT4oMini = shared.ChatModelGPT4oMini

// This is an alias to an internal value.
const ChatModelGPT4oMini2024_07_18 = shared.ChatModelGPT4oMini2024_07_18

// This is an alias to an internal value.
const ChatModelGPT4Turbo = shared.ChatModelGPT4Turbo

// This is an alias to an internal value.
const ChatModelGPT4Turbo2024_04_09 = shared.ChatModelGPT4Turbo2024_04_09

// This is an alias to an internal value.
const ChatModelGPT4_0125Preview = shared.ChatModelGPT4_0125Preview

// This is an alias to an internal value.
const ChatModelGPT4TurboPreview = shared.ChatModelGPT4TurboPreview

// This is an alias to an internal value.
const ChatModelGPT4_1106Preview = shared.ChatModelGPT4_1106Preview

// This is an alias to an internal value.
const ChatModelGPT4VisionPreview = shared.ChatModelGPT4VisionPreview

// This is an alias to an internal value.
const ChatModelGPT4 = shared.ChatModelGPT4

// This is an alias to an internal value.
const ChatModelGPT4_0314 = shared.ChatModelGPT4_0314

// This is an alias to an internal value.
const ChatModelGPT4_0613 = shared.ChatModelGPT4_0613

// This is an alias to an internal value.
const ChatModelGPT4_32k = shared.ChatModelGPT4_32k

// This is an alias to an internal value.
const ChatModelGPT4_32k0314 = shared.ChatModelGPT4_32k0314

// This is an alias to an internal value.
const ChatModelGPT4_32k0613 = shared.ChatModelGPT4_32k0613

// This is an alias to an internal value.
const ChatModelGPT3_5Turbo = shared.ChatModelGPT3_5Turbo

// This is an alias to an internal value.
const ChatModelGPT3_5Turbo16k = shared.ChatModelGPT3_5Turbo16k

// This is an alias to an internal value.
const ChatModelGPT3_5Turbo0301 = shared.ChatModelGPT3_5Turbo0301

// This is an alias to an internal value.
const ChatModelGPT3_5Turbo0613 = shared.ChatModelGPT3_5Turbo0613

// This is an alias to an internal value.
const ChatModelGPT3_5Turbo1106 = shared.ChatModelGPT3_5Turbo1106

// This is an alias to an internal value.
const ChatModelGPT3_5Turbo0125 = shared.ChatModelGPT3_5Turbo0125

// This is an alias to an internal value.
const ChatModelGPT3_5Turbo16k0613 = shared.ChatModelGPT3_5Turbo16k0613

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

# OpenAI Go API Library

<a href="https://pkg.go.dev/github.com/openai/openai-go"><img src="https://pkg.go.dev/badge/github.com/openai/openai-go.svg" alt="Go Reference"></a>

The OpenAI Go library provides convenient access to the [OpenAI REST API](https://platform.openai.com/docs)
from applications written in Go.

> [!WARNING]
> The latest version of this package uses a new design with significant breaking changes.
> Please refer to the [migration guide](./MIGRATION.md) for more information on how to update your code.

## Installation

<!-- x-release-please-start-version -->

```go
import (
	"github.com/openai/openai-go" // imported as openai
)
```

<!-- x-release-please-end -->

Or to pin the version:

<!-- x-release-please-start-version -->

```sh
go get -u 'github.com/openai/openai-go@v1.12.0'
```

<!-- x-release-please-end -->

## Requirements

This library requires Go 1.18+.

## Usage

The full API of this library can be found in [api.md](api.md).

```go
package main

import (
	"context"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/shared"
)

func main() {
	client := openai.NewClient(
		option.WithAPIKey("My API Key"), // defaults to os.LookupEnv("OPENAI_API_KEY")
	)
	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage("Say this is a test"),
		},
		Model: openai.ChatModelGPT4o,
	})
	if err != nil {
		panic(err.Error())
	}
	println(chatCompletion.Choices[0].Message.Content)
}

```

<details>
<summary>Conversations</summary>

```go
param := openai.ChatCompletionNewParams{
	Messages: []openai.ChatCompletionMessageParamUnion{
		openai.UserMessage("What kind of houseplant is easy to take care of?"),
	},
	Seed:     openai.Int(1),
	Model:    openai.ChatModelGPT4o,
}

completion, err := client.Chat.Completions.New(ctx, param)

param.Messages = append(param.Messages, completion.Choices[0].Message.ToParam())
param.Messages = append(param.Messages, openai.UserMessage("How big are those?"))

// continue the conversation
completion, err = client.Chat.Completions.New(ctx, param)
```

</details>

<details>
<summary>Streaming responses</summary>

```go
question := "Write an epic"

stream := client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
	Messages: []openai.ChatCompletionMessageParamUnion{
		openai.UserMessage(question),
	},
	Seed:  openai.Int(0),
	Model: openai.ChatModelGPT4o,
})

// optionally, an accumulator helper can be used
acc := openai.ChatCompletionAccumulator{}

for stream.Next() {
	chunk := stream.Current()
	acc.AddChunk(chunk)

	if content, ok := acc.JustFinishedContent(); ok {
		println("Content stream finished:", content)
	}

	// if using tool calls
	if tool, ok := acc.JustFinishedToolCall(); ok {
		println("Tool call stream finished:", tool.Index, tool.Name, tool.Arguments)
	}

	if refusal, ok := acc.JustFinishedRefusal(); ok {
		println("Refusal stream finished:", refusal)
	}

	// it's best to use chunks after handling JustFinished events
	if len(chunk.Choices) > 0 {
		println(chunk.Choices[0].Delta.Content)
	}
}

if stream.Err() != nil {
	panic(stream.Err())
}

// After the stream is finished, acc can be used like a ChatCompletion
_ = acc.Choices[0].Message.Content
```

> See the [full streaming and accumulation example](./examples/chat-completion-accumulating/main.go)

</details>

<details>
<summary>Tool calling</summary>

```go
import (
	"encoding/json"
	// ...
)

// ...

question := "What is the weather in New York City?"

params := openai.ChatCompletionNewParams{
	Messages: []openai.ChatCompletionMessageParamUnion{
		openai.UserMessage(question),
	},
	Tools: []openai.ChatCompletionToolParam{
		{
			Function: openai.FunctionDefinitionParam{
				Name:        "get_weather",
				Description: openai.String("Get weather at the given location"),
				Parameters: openai.FunctionParameters{
					"type": "object",
					"properties": map[string]interface{}{
						"location": map[string]string{
							"type": "string",
						},
					},
					"required": []string{"location"},
				},
			},
		},
	},
	Model: openai.ChatModelGPT4o,
}

// If there is a was a function call, continue the conversation
params.Messages = append(params.Messages, completion.Choices[0].Message.ToParam())
for _, toolCall := range toolCalls {
	if toolCall.Function.Name == "get_weather" {
		// Extract the location from the function call arguments
		var args map[string]interface{}
		err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args)
		if err != nil {
			panic(err)
		}
		location := args["location"].(string)

		// Simulate getting weather data
		weatherData := getWeather(location)

		// Print the weather data
		fmt.Printf("Weather in %s: %s\n", location, weatherData)

		params.Messages = append(params.Messages, openai.ToolMessage(weatherData, toolCall.ID))
	}
}

// ... continue the conversation with the information provided by the tool
```

> See the [full tool calling example](./examples/chat-completion-tool-calling/main.go)

</details>

<details>
<summary>Structured outputs</summary>

```go
import (
	"encoding/json"
	"github.com/invopop/jsonschema"
	// ...
)

// A struct that will be converted to a Structured Outputs response schema
type HistoricalComputer struct {
	Origin       Origin   `json:"origin" jsonschema_description:"The origin of the computer"`
	Name         string   `json:"full_name" jsonschema_description:"The name of the device model"`
	Legacy       string   `json:"legacy" jsonschema:"enum=positive,enum=neutral,enum=negative" jsonschema_description:"Its influence on the field of computing"`
	NotableFacts []string `json:"notable_facts" jsonschema_description:"A few key facts about the computer"`
}

type Origin struct {
	YearBuilt    int64  `json:"year_of_construction" jsonschema_description:"The year it was made"`
	Organization string `json:"organization" jsonschema_description:"The organization that was in charge of its development"`
}

func GenerateSchema[T any]() interface{} {
	// Structured Outputs uses a subset of JSON schema
	// These flags are necessary to comply with the subset
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	schema := reflector.Reflect(v)
	return schema
}

// Generate the JSON schema at initialization time
var HistoricalComputerResponseSchema = GenerateSchema[HistoricalComputer]()

func main() {

	// ...

	question := "What computer ran the first neural network?"

	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        "historical_computer",
		Description: openai.String("Notable information about a computer"),
		Schema:      HistoricalComputerResponseSchema,
		Strict:      openai.Bool(true),
	}

	chat, _ := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		// ...
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
				JSONSchema: schemaParam,
			},
		},
		// only certain models can perform structured outputs
		Model: openai.ChatModelGPT4o2024_08_06,
	})

	// extract into a well-typed struct
	var historicalComputer HistoricalComputer
	_ = json.Unmarshal([]byte(chat.Choices[0].Message.Content), &historicalComputer)

	historicalComputer.Name
	historicalComputer.Origin.YearBuilt
	historicalComputer.Origin.Organization
	for i, fact := range historicalComputer.NotableFacts {
		// ...
	}
}
```

> See the [full structured outputs example](./examples/structured-outputs/main.go)

</details>


### Request fields

The openai library uses the [`omitzero`](https://tip.golang.org/doc/go1.24#encodingjsonpkgencodingjson)
semantics from the Go 1.24+ `encoding/json` release for request fields.

Required primitive fields (`int64`, `string`, etc.) feature the tag <code>\`json:"...,required"\`</code>. These
fields are always serialized, even their zero values.

Optional primitive types are wrapped in a `param.Opt[T]`. These fields can be set with the provided constructors, `openai.String(string)`, `openai.Int(int64)`, etc.

Any `param.Opt[T]`, map, slice, struct or string enum uses the
tag <code>\`json:"...,omitzero"\`</code>. Its zero value is considered omitted.

The `param.IsOmitted(any)` function can confirm the presence of any `omitzero` field.

```go
p := openai.ExampleParams{
	ID:   "id_xxx",             // required property
	Name: openai.String("..."), // optional property

	Point: openai.Point{
		X: 0,             // required field will serialize as 0
		Y: openai.Int(1), // optional field will serialize as 1
		// ... omitted non-required fields will not be serialized
	},

	Origin: openai.Origin{}, // the zero value of [Origin] is considered omitted
}
```

To send `null` instead of a `param.Opt[T]`, use `param.Null[T]()`.
To send `null` instead of a struct `T`, use `param.NullStruct[T]()`.

```go
p.Name = param.Null[string]()       // 'null' instead of string
p.Point = param.NullStruct[Point]() // 'null' instead of struct

param.IsNull(p.Name)  // true
param.IsNull(p.Point) // true
```

Request structs contain a `.SetExtraFields(map[string]any)` method which can send non-conforming
fields in the request body. Extra fields overwrite any struct fields with a matching
key. For security reasons, only use `SetExtraFields` with trusted data.

To send a custom value instead of a struct, use `param.Override[T](value)`.

```go
// In cases where the API specifies a given type,
// but you want to send something else, use [SetExtraFields]:
p.SetExtraFields(map[string]any{
	"x": 0.01, // send "x" as a float instead of int
})

// Send a number instead of an object
custom := param.Override[openai.FooParams](12)
```

### Request unions

Unions are represented as a struct with fields prefixed by "Of" for each of it's variants,
only one field can be non-zero. The non-zero field will be serialized.

Sub-properties of the union can be accessed via methods on the union struct.
These methods return a mutable pointer to the underlying data, if present.

```go
// Only one field can be non-zero, use param.IsOmitted() to check if a field is set
type AnimalUnionParam struct {
	OfCat *Cat `json:",omitzero,inline`
	OfDog *Dog `json:",omitzero,inline`
}

animal := AnimalUnionParam{
	OfCat: &Cat{
		Name: "Whiskers",
		Owner: PersonParam{
			Address: AddressParam{Street: "3333 Coyote Hill Rd", Zip: 0},
		},
	},
}

// Mutating a field
if address := animal.GetOwner().GetAddress(); address != nil {
	address.ZipCode = 94304
}
```

### Response objects

All fields in response structs are ordinary value types (not pointers or wrappers).
Response structs also include a special `JSON` field containing metadata about
each property.

```go
type Animal struct {
	Name   string `json:"name,nullable"`
	Owners int    `json:"owners"`
	Age    int    `json:"age"`
	JSON   struct {
		Name        respjson.Field
		Owner       respjson.Field
		Age         respjson.Field
		ExtraFields map[string]respjson.Field
	} `json:"-"`
}
```

To handle optional data, use the `.Valid()` method on the JSON field.
`.Valid()` returns true if a field is not `null`, not present, or couldn't be marshaled.

If `.Valid()` is false, the corresponding field will simply be its zero value.

```go
raw := `{"owners": 1, "name": null}`

var res Animal
json.Unmarshal([]byte(raw), &res)

// Accessing regular fields

res.Owners // 1
res.Name   // ""
res.Age    // 0

// Optional field checks

res.JSON.Owners.Valid() // true
res.JSON.Name.Valid()   // false
res.JSON.Age.Valid()    // false

// Raw JSON values

res.JSON.Owners.Raw()                  // "1"
res.JSON.Name.Raw() == "null"          // true
res.JSON.Name.Raw() == respjson.Null   // true
res.JSON.Age.Raw() == ""               // true
res.JSON.Age.Raw() == respjson.Omitted // true
```

These `.JSON` structs also include an `ExtraFields` map containing
any properties in the json response that were not specified
in the struct. This can be useful for API features not yet
present in the SDK.

```go
body := res.JSON.ExtraFields["my_unexpected_field"].Raw()
```

### Response Unions

In responses, unions are represented by a flattened struct containing all possible fields from each of the
object variants.
To convert it to a variant use the `.AsFooVariant()` method or the `.AsAny()` method if present.

If a response value union contains primitive values, primitive fields will be alongside
the properties but prefixed with `Of` and feature the tag `json:"...,inline"`.

```go
type AnimalUnion struct {
	// From variants [Dog], [Cat]
	Owner Person `json:"owner"`
	// From variant [Dog]
	DogBreed string `json:"dog_breed"`
	// From variant [Cat]
	CatBreed string `json:"cat_breed"`
	// ...

	JSON struct {
		Owner respjson.Field
		// ...
	} `json:"-"`
}

// If animal variant
if animal.Owner.Address.ZipCode == "" {
	panic("missing zip code")
}

// Switch on the variant
switch variant := animal.AsAny().(type) {
case Dog:
case Cat:
default:
	panic("unexpected type")
}
```

### RequestOptions

This library uses the functional options pattern. Functions defined in the
`option` package return a `RequestOption`, which is a closure that mutates a
`RequestConfig`. These options can be supplied to the client or at individual
requests. For example:

```go
client := openai.NewClient(
	// Adds a header to every request made by the client
	option.WithHeader("X-Some-Header", "custom_header_info"),
)

client.Chat.Completions.New(context.TODO(), ...,
	// Override the header
	option.WithHeader("X-Some-Header", "some_other_custom_header_info"),
	// Add an undocumented field to the request body, using sjson syntax
	option.WithJSONSet("some.json.path", map[string]string{"my": "object"}),
)
```

The request option `option.WithDebugLog(nil)` may be helpful while debugging.

See the [full list of request options](https://pkg.go.dev/github.com/openai/openai-go/option).

### Pagination

This library provides some conveniences for working with paginated list endpoints.

You can use `.ListAutoPaging()` methods to iterate through items across all pages:

```go
iter := client.FineTuning.Jobs.ListAutoPaging(context.TODO(), openai.FineTuningJobListParams{
	Limit: openai.Int(20),
})
// Automatically fetches more pages as needed.
for iter.Next() {
	fineTuningJob := iter.Current()
	fmt.Printf("%+v\n", fineTuningJob)
}
if err := iter.Err(); err != nil {
	panic(err.Error())
}
```

Or you can use simple `.List()` methods to fetch a single page and receive a standard response object
with additional helper methods like `.GetNextPage()`, e.g.:

```go
page, err := client.FineTuning.Jobs.List(context.TODO(), openai.FineTuningJobListParams{
	Limit: openai.Int(20),
})
for page != nil {
	for _, job := range page.Data {
		fmt.Printf("%+v\n", job)
	}
	page, err = page.GetNextPage()
}
if err != nil {
	panic(err.Error())
}
```

### Errors

When the API returns a non-success status code, we return an error with type
`*openai.Error`. This contains the `StatusCode`, `*http.Request`, and
`*http.Response` values of the request, as well as the JSON of the error body
(much like other response objects in the SDK).

To handle errors, we recommend that you use the `errors.As` pattern:

```go
_, err := client.FineTuning.Jobs.New(context.TODO(), openai.FineTuningJobNewParams{
	Model:        openai.FineTuningJobNewParamsModelBabbage002,
	TrainingFile: "file-abc123",
})
if err != nil {
	var apierr *openai.Error
	if errors.As(err, &apierr) {
		println(string(apierr.DumpRequest(true)))  // Prints the serialized HTTP request
		println(string(apierr.DumpResponse(true))) // Prints the serialized HTTP response
	}
	panic(err.Error()) // GET "/fine_tuning/jobs": 400 Bad Request { ... }
}
```

When other errors occur, they are returned unwrapped; for example,
if HTTP transport fails, you might receive `*url.Error` wrapping `*net.OpError`.

### Timeouts

Requests do not time out by default; use context to configure a timeout for a request lifecycle.

Note that if a request is [retried](#retries), the context timeout does not start over.
To set a per-retry timeout, use `option.WithRequestTimeout()`.

```go
// This sets the timeout for the request, including all the retries.
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()
client.Chat.Completions.New(
	ctx,
	openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{{
			OfUser: &openai.ChatCompletionUserMessageParam{
				Content: openai.ChatCompletionUserMessageParamContentUnion{
					OfString: openai.String("How can I list all files in a directory using Python?"),
				},
			},
		}},
		Model: shared.ChatModelGPT4_1,
	},
	// This sets the per-retry timeout
	option.WithRequestTimeout(20*time.Second),
)
```

### File uploads

Request parameters that correspond to file uploads in multipart requests are typed as
`io.Reader`. The contents of the `io.Reader` will by default be sent as a multipart form
part with the file name of "anonymous_file" and content-type of "application/octet-stream".

The file name and content-type can be customized by implementing `Name() string` or `ContentType()
string` on the run-time type of `io.Reader`. Note that `os.File` implements `Name() string`, so a
file returned by `os.Open` will be sent with the file name on disk.

We also provide a helper `openai.File(reader io.Reader, filename string, contentType string)`
which can be used to wrap any `io.Reader` with the appropriate file name and content type.

```go
// A file from the file system
file, err := os.Open("input.jsonl")
openai.FileNewParams{
	File:    file,
	Purpose: openai.FilePurposeFineTune,
}

// A file from a string
openai.FileNewParams{
	File:    strings.NewReader("my file contents"),
	Purpose: openai.FilePurposeFineTune,
}

// With a custom filename and contentType
openai.FileNewParams{
	File:    openai.File(strings.NewReader(`{"hello": "foo"}`), "file.go", "application/json"),
	Purpose: openai.FilePurposeFineTune,
}
```

## Webhook Verification

Verifying webhook signatures is _optional but encouraged_.

For more information about webhooks, see [the API docs](https://platform.openai.com/docs/guides/webhooks).

### Parsing webhook payloads

For most use cases, you will likely want to verify the webhook and parse the payload at the same time. To achieve this, we provide the method `client.Webhooks.Unwrap()`, which parses a webhook request and verifies that it was sent by OpenAI. This method will return an error if the signature is invalid.

Note that the `body` parameter should be the raw JSON bytes sent from the server (do not parse it first). The `Unwrap()` method will parse this JSON for you into an event object after verifying the webhook was sent from OpenAI.

```go
package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/webhooks"
)

func main() {
	client := openai.NewClient(
		option.WithWebhookSecret(os.Getenv("OPENAI_WEBHOOK_SECRET")), // env var used by default; explicit here.
	)

	r := gin.Default()
	
	r.POST("/webhook", func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading request body"})
			return
		}
		defer c.Request.Body.Close()

		webhookEvent, err := client.Webhooks.Unwrap(body, c.Request.Header)
		if err != nil {
			log.Printf("Invalid webhook signature: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid signature"})
			return
		}

		switch event := webhookEvent.AsAny().(type) {
		case webhooks.ResponseCompletedWebhookEvent:
			log.Printf("Response completed: %+v", event.Data)
		case webhooks.ResponseFailedWebhookEvent:
			log.Printf("Response failed: %+v", event.Data)
		default:
			log.Printf("Unhandled event type: %T", event)
		}

		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	r.Run(":8000")
}
```

### Verifying webhook payloads directly

In some cases, you may want to verify the webhook separately from parsing the payload. If you prefer to handle these steps separately, we provide the method `client.Webhooks.VerifySignature()` to _only verify_ the signature of a webhook request. Like `Unwrap()`, this method will return an error if the signature is invalid.

Note that the `body` parameter should be the raw JSON bytes sent from the server (do not parse it first). You will then need to parse the body after verifying the signature.

```go
package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func main() {
	client := openai.NewClient(
		option.WithWebhookSecret(os.Getenv("OPENAI_WEBHOOK_SECRET")), // env var used by default; explicit here.
	)

	r := gin.Default()
	
	r.POST("/webhook", func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading request body"})
			return
		}
		defer c.Request.Body.Close()

		err = client.Webhooks.VerifySignature(body, c.Request.Header)
		if err != nil {
			log.Printf("Invalid webhook signature: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid signature"})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	r.Run(":8000")
}
```

### Retries

Certain errors will be automatically retried 2 times by default, with a short exponential backoff.
We retry by default all connection errors, 408 Request Timeout, 409 Conflict, 429 Rate Limit,
and >=500 Internal errors.

You can use the `WithMaxRetries` option to configure or disable this:

```go
// Configure the default for all requests:
client := openai.NewClient(
	option.WithMaxRetries(0), // default is 2
)

// Override per-request:
client.Chat.Completions.New(
	context.TODO(),
	openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{{
			OfUser: &openai.ChatCompletionUserMessageParam{
				Content: openai.ChatCompletionUserMessageParamContentUnion{
					OfString: openai.String("How can I get the name of the current day in JavaScript?"),
				},
			},
		}},
		Model: shared.ChatModelGPT4_1,
	},
	option.WithMaxRetries(5),
)
```

### Accessing raw response data (e.g. response headers)

You can access the raw HTTP response data by using the `option.WithResponseInto()` request option. This is useful when
you need to examine response headers, status codes, or other details.

```go
// Create a variable to store the HTTP response
var response *http.Response
chatCompletion, err := client.Chat.Completions.New(
	context.TODO(),
	openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{{
			OfUser: &openai.ChatCompletionUserMessageParam{
				Content: openai.ChatCompletionUserMessageParamContentUnion{
					OfString: openai.String("Say this is a test"),
				},
			},
		}},
		Model: shared.ChatModelGPT4_1,
	},
	option.WithResponseInto(&response),
)
if err != nil {
	// handle error
}
fmt.Printf("%+v\n", chatCompletion)

fmt.Printf("Status Code: %d\n", response.StatusCode)
fmt.Printf("Headers: %+#v\n", response.Header)
```

### Making custom/undocumented requests

This library is typed for convenient access to the documented API. If you need to access undocumented
endpoints, params, or response properties, the library can still be used.

#### Undocumented endpoints

To make requests to undocumented endpoints, you can use `client.Get`, `client.Post`, and other HTTP verbs.
`RequestOptions` on the client, such as retries, will be respected when making these requests.

```go
var (
    // params can be an io.Reader, a []byte, an encoding/json serializable object,
    // or a "…Params" struct defined in this library.
    params map[string]any

    // result can be an []byte, *http.Response, a encoding/json deserializable object,
    // or a model defined in this library.
    result *http.Response
)
err := client.Post(context.Background(), "/unspecified", params, &result)
if err != nil {
    …
}
```

#### Undocumented request params

To make requests using undocumented parameters, you may use either the `option.WithQuerySet()`
or the `option.WithJSONSet()` methods.

```go
params := FooNewParams{
    ID:   "id_xxxx",
    Data: FooNewParamsData{
        FirstName: openai.String("John"),
    },
}
client.Foo.New(context.Background(), params, option.WithJSONSet("data.last_name", "Doe"))
```

#### Undocumented response properties

To access undocumented response properties, you may either access the raw JSON of the response as a string
with `result.JSON.RawJSON()`, or get the raw JSON of a particular field on the result with
`result.JSON.Foo.Raw()`.

Any fields that are not present on the response struct will be saved and can be accessed by `result.JSON.ExtraFields()` which returns the extra fields as a `map[string]Field`.

### Middleware

We provide `option.WithMiddleware` which applies the given
middleware to requests.

```go
func Logger(req *http.Request, next option.MiddlewareNext) (res *http.Response, err error) {
	// Before the request
	start := time.Now()
	LogReq(req)

	// Forward the request to the next handler
	res, err = next(req)

	// Handle stuff after the request
	end := time.Now()
	LogRes(res, err, start - end)

    return res, err
}

client := openai.NewClient(
	option.WithMiddleware(Logger),
)
```

When multiple middlewares are provided as variadic arguments, the middlewares
are applied left to right. If `option.WithMiddleware` is given
multiple times, for example first in the client then the method, the
middleware in the client will run first and the middleware given in the method
will run next.

You may also replace the default `http.Client` with
`option.WithHTTPClient(client)`. Only one http client is
accepted (this overwrites any previous client) and receives requests after any
middleware has been applied.

## Microsoft Azure OpenAI

To use this library with [Azure OpenAI]https://learn.microsoft.com/azure/ai-services/openai/overview),
use the option.RequestOption functions in the `azure` package.

```go
package main

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/azure"
)

func main() {
	const azureOpenAIEndpoint = "https://<azure-openai-resource>.openai.azure.com"

	// The latest API versions, including previews, can be found here:
	// ttps://learn.microsoft.com/en-us/azure/ai-services/openai/reference#rest-api-versionng
	const azureOpenAIAPIVersion = "2024-06-01"

	tokenCredential, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		fmt.Printf("Failed to create the DefaultAzureCredential: %s", err)
		os.Exit(1)
	}

	client := openai.NewClient(
		azure.WithEndpoint(azureOpenAIEndpoint, azureOpenAIAPIVersion),

		// Choose between authenticating using a TokenCredential or an API Key
		azure.WithTokenCredential(tokenCredential),
		// or azure.WithAPIKey(azureOpenAIAPIKey),
	)
}
```


## Semantic versioning

This package generally follows [SemVer](https://semver.org/spec/v2.0.0.html) conventions, though certain backwards-incompatible changes may be released as minor versions:

1. Changes to library internals which are technically public but not intended or documented for external use. _(Please open a GitHub issue to let us know if you are relying on such internals.)_
2. Changes that we do not expect to impact the vast majority of users in practice.

We take backwards-compatibility seriously and work hard to ensure you can rely on a smooth upgrade experience.

We are keen for your feedback; please open an [issue](https://www.github.com/openai/openai-go/issues) with questions, bugs, or suggestions.

## Contributing

See [the contributing documentation](./CONTRIBUTING.md).

# OpenAI Go API Library

<!-- x-release-please-start-version -->

<a href="https://pkg.go.dev/github.com/openai/openai-go/v3"><img src="https://pkg.go.dev/badge/github.com/openai/openai-go.svg" alt="Go Reference"></a>

<!-- x-release-please-end -->

The OpenAI Go library provides convenient access to the [OpenAI REST API](https://platform.openai.com/docs)
from applications written in Go.

> [!WARNING]
> The latest version of this package has small and limited breaking changes.
> See the [changelog](CHANGELOG.md) for details.

## Installation

<!-- x-release-please-start-version -->

```go
import (
	"github.com/openai/openai-go/v3" // imported as openai
)
```

<!-- x-release-please-end -->

Or to pin an SDK version (see the Go compatibility note below):

<!-- x-release-please-start-version -->

```sh
go get -u 'github.com/openai/openai-go/v3@v3.54.0'
```

<!-- x-release-please-end -->

## Requirements

SDK v3.45.0 and later require Go 1.25 or later. If your application must
remain on Go 1.22–1.24, pin SDK v3.44.0, the final compatible release. Older
SDK releases receive no guaranteed fixes or security backports.

See the [Go version support policy](GO_VERSION_POLICY.md) for the supported
release window and upgrade guidance.

## Usage

The full API of this library can be found in [api.md](api.md).

The primary API for interacting with OpenAI models is the [Responses API](https://platform.openai.com/docs/api-reference/responses). You can generate text from the model with the code below.

```go
package main

import (
	"context"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

func main() {
	ctx := context.Background()
	client := openai.NewClient(
		option.WithAPIKey("My API Key"), // defaults to os.LookupEnv("OPENAI_API_KEY")
	)

	question := "Write me a haiku about computers"

	resp, err := client.Responses.New(ctx, responses.ResponseNewParams{
		Input: responses.ResponseNewParamsInputUnion{OfString: openai.String(question)},
		Model: openai.ChatModelGPT5_2,
	})

	if err != nil {
		panic(err)
	}

	println(resp.OutputText())
}
```

<details>
<summary>Multi-turn Responses</summary>

```go
response, err := client.Responses.New(ctx, responses.ResponseNewParams{
	Model: openai.ChatModelGPT5_2,
	Input: responses.ResponseNewParamsInputUnion{
		OfString: openai.String("What is the capital of France?"),
	},
})
if err != nil {
	panic(err)
}
fmt.Println("First response:", response.OutputText())

// Use PreviousResponseID to continue the conversation
response, err = client.Responses.New(ctx, responses.ResponseNewParams{
	Model:              openai.ChatModelGPT5_2,
	PreviousResponseID: openai.String(response.ID),
	Input: responses.ResponseNewParamsInputUnion{
		OfString: openai.String("And what is the population of that city?"),
	},
})
if err != nil {
	panic(err)
}
fmt.Println("Second response:", response.OutputText())
```
</details>

<details>
<summary>Conversations</summary>

```go
conv, err := client.Conversations.New(ctx, conversations.ConversationNewParams{})
if err != nil {
	panic(err)
}
fmt.Println("Created conversation:", conv.ID)

response, err := client.Responses.New(ctx, responses.ResponseNewParams{
	Model: openai.ChatModelGPT5_2,
	Input: responses.ResponseNewParamsInputUnion{
		OfString: openai.String("Hello! Remember that my favorite color is blue."),
	},
	Conversation: responses.ResponseNewParamsConversationUnion{
		OfConversationObject: &responses.ResponseConversationParam{
			ID: conv.ID,
		},
	},
})
if err != nil {
	panic(err)
}
fmt.Println("First response:", response.OutputText())

// Continue the conversation
response, err = client.Responses.New(ctx, responses.ResponseNewParams{
	Model: openai.ChatModelGPT5_2,
	Input: responses.ResponseNewParamsInputUnion{
		OfString: openai.String("What is my favorite color?"),
	},
	Conversation: responses.ResponseNewParamsConversationUnion{
		OfConversationObject: &responses.ResponseConversationParam{
			ID: conv.ID,
		},
	},
})
if err != nil {
	panic(err)
}
fmt.Println("Second response:", response.OutputText())

items, err := client.Conversations.Items.List(ctx, conv.ID, conversations.ItemListParams{})
if err != nil {
	panic(err)
}
fmt.Println("Conversation has", len(items.Data), "items")
```

</details>

<details>
<summary>Streaming responses</summary>

```go
ctx := context.Background()

stream := client.Responses.NewStreaming(ctx, responses.ResponseNewParams{
	Model: openai.ChatModelGPT5_2,
	Input: responses.ResponseNewParamsInputUnion{
		OfString: openai.String("Write a haiku about programming"),
	},
})
defer func() { _ = stream.Close() }()

for stream.Next() {
	event := stream.Current()
	print(event.Delta)
}

if stream.Err() != nil {
	panic(stream.Err())
}
```

> See the [full streaming example](./examples/responses-streaming/main.go)

</details>

<details>
<summary>Tool calling</summary>

```go
ctx := context.Background()

params := responses.ResponseNewParams{
	Model: openai.ChatModelGPT5_2,
	Input: responses.ResponseNewParamsInputUnion{
		OfString: openai.String("What is the weather in New York City?"),
	},
	Tools: []responses.ToolUnionParam{{
		OfFunction: &responses.FunctionToolParam{
			Name:        "get_weather",
			Description: openai.String("Get weather at the given location"),
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"location": map[string]string{
						"type": "string",
					},
				},
				"required": []string{"location"},
			},
		},
	}},
}

response, _ := client.Responses.New(ctx, params)

// Check for function calls in the response output
for _, item := range response.Output {
	if item.Type == "function_call" {
		toolCall := item.AsFunctionCall()
		if toolCall.Name == "get_weather" {
			// Extract arguments and call your function
			var args map[string]any
			json.Unmarshal([]byte(toolCall.Arguments), &args)
			location := args["location"].(string)

			// Simulate getting weather data
			weatherData := getWeather(location)
			fmt.Printf("Weather in %s: %s\n", location, weatherData)

			// Continue conversation with function result
			response, _ = client.Responses.New(ctx, responses.ResponseNewParams{
				Model:              openai.ChatModelGPT5_2,
				PreviousResponseID: openai.String(response.ID),
				Input: responses.ResponseNewParamsInputUnion{
					OfInputItemList: []responses.ResponseInputItemUnionParam{{
						OfFunctionCallOutput: &responses.ResponseInputItemFunctionCallOutputParam{
							CallID: toolCall.CallID,
							Output: responses.ResponseInputItemFunctionCallOutputOutputUnionParam{
								OfString: openai.String(weatherData),
							},
						},
					}},
				},
			})
		}
	}
}
```

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

// Structured Outputs uses a subset of JSON schema
// These flags are necessary to comply with the subset
func GenerateSchema[T any]() (map[string]any, error) {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	schema := reflector.Reflect(v)

	data, err := json.Marshal(schema)
	if err != nil {
		return nil, err
	}
	// Preserve schema values as raw JSON so integer constraints are not rounded
	// through float64 before the SDK serializes the request.
	var rawSchema map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawSchema); err != nil {
		return nil, err
	}
	result := make(map[string]any, len(rawSchema))
	for key, value := range rawSchema {
		result[key] = value
	}
	return result, nil
}

func main() {
	client := openai.NewClient()
	ctx := context.Background()
	schema, err := GenerateSchema[HistoricalComputer]()
	if err != nil {
		panic(err)
	}

	response, err := client.Responses.New(ctx, responses.ResponseNewParams{
		Model: openai.ChatModelGPT5_2,
		Input: responses.ResponseNewParamsInputUnion{
			OfString: openai.String("What computer ran the first neural network?"),
		},
		Text: responses.ResponseTextConfigParam{
			Format: responses.ResponseFormatTextConfigUnionParam{
				OfJSONSchema: &responses.ResponseFormatTextJSONSchemaConfigParam{
					Name:        "historical_computer",
					Description: openai.String("Notable information about a computer"),
					Schema:      schema,
					Strict:      openai.Bool(true),
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}

	// extract into a well-typed struct
	var historicalComputer HistoricalComputer
	if err := json.Unmarshal([]byte(response.OutputText()), &historicalComputer); err != nil {
		panic(err)
	}

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

### Chat Completions API

The previous standard (supported indefinitely) for generating text is the [Chat Completions API](https://platform.openai.com/docs/api-reference/chat). You can use that API to generate text from the model with the code below.

```go
package main

import (
	"context"

	"github.com/openai/openai-go/v3"
)

func main() {
	client := openai.NewClient()

	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.DeveloperMessage("You are a coding assistant that talks like a pirate."),
			openai.UserMessage("How do I check if a slice is empty in Go?"),
		},
		Model: openai.ChatModelGPT5_2,
	})
	if err != nil {
		panic(err)
	}

	println(chatCompletion.Choices[0].Message.Content)
}
```

### Request fields

The openai library uses the [`omitzero`](https://tip.golang.org/doc/go1.24#encodingjsonpkgencodingjson)
semantics from the Go 1.24+ `encoding/json` release for request fields.

Required primitive fields (`int64`, `string`, etc.) feature the tag <code>\`api:"required"\`</code>. These
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

Unions are represented as a struct with fields prefixed by "Of" for each of its variants,
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

client.Responses.New(context.TODO(), responses.ResponseNewParams{...},
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

> [!WARNING]
> `Error.DumpRequest`, `Error.DumpResponse`, and `Error.Error` expose raw
> diagnostics that may include authorization headers, credentials in URLs, and
> sensitive request or response bodies. The dump `body` option does not redact
> headers. Sanitize this output before logging, sharing, or storing it.

```go
_, err := client.FineTuning.Jobs.New(context.TODO(), openai.FineTuningJobNewParams{
	Model:        openai.FineTuningJobNewParamsModel("gpt-4o"),
	TrainingFile: "file-abc123",
})
if err != nil {
	var apierr *openai.Error
	if errors.As(err, &apierr) {
		fmt.Printf("OpenAI API error (status: %d)\n", apierr.StatusCode)
	}
	panic("OpenAI request failed")
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
client.Responses.New(
	ctx,
	responses.ResponseNewParams{
		Model: openai.ChatModelGPT5_2,
		Input: responses.ResponseNewParamsInputUnion{
			OfString: openai.String("How can I list all files in a directory using Python?"),
		},
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

Webhook payloads are unauthenticated until signature verification succeeds. The
example below limits each request to 1 MiB before buffering it and configures
server read timeouts. Enforce the same or a smaller limit at your reverse proxy
as defense in depth.

```go
package main

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/webhooks"
)

const maxWebhookBodySize = 1 << 20 // 1 MiB

func main() {
	client := openai.NewClient(
		option.WithWebhookSecret(os.Getenv("OPENAI_WEBHOOK_SECRET")), // env var used by default; explicit here.
	)

	r := gin.Default()

	r.POST("/webhook", func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxWebhookBodySize)
		defer func() { _ = c.Request.Body.Close() }()

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			var maxBytesError *http.MaxBytesError
			if errors.As(err, &maxBytesError) {
				c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "request body too large"})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": "error reading request body"})
			return
		}

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

	server := &http.Server{
		Addr:              ":8000",
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
```

### Verifying webhook payloads directly

In some cases, you may want to verify the webhook separately from parsing the payload. If you prefer to handle these steps separately, we provide the method `client.Webhooks.VerifySignature()` to _only verify_ the signature of a webhook request. Like `Unwrap()`, this method will return an error if the signature is invalid.

Note that the `body` parameter should be the raw JSON bytes sent from the server (do not parse it first). You will then need to parse the body after verifying the signature.

As above, bound the unauthenticated request body before reading it and configure
server read timeouts. This example also uses a 1 MiB maximum.

```go
package main

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

const maxWebhookBodySize = 1 << 20 // 1 MiB

func main() {
	client := openai.NewClient(
		option.WithWebhookSecret(os.Getenv("OPENAI_WEBHOOK_SECRET")), // env var used by default; explicit here.
	)

	r := gin.Default()

	r.POST("/webhook", func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxWebhookBodySize)
		defer func() { _ = c.Request.Body.Close() }()

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			var maxBytesError *http.MaxBytesError
			if errors.As(err, &maxBytesError) {
				c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "request body too large"})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": "error reading request body"})
			return
		}

		err = client.Webhooks.VerifySignature(body, c.Request.Header)
		if err != nil {
			log.Printf("Invalid webhook signature: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid signature"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	server := &http.Server{
		Addr:              ":8000",
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
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
client.Responses.New(
	context.TODO(),
	responses.ResponseNewParams{
		Model: openai.ChatModelGPT5_2,
		Input: responses.ResponseNewParamsInputUnion{
			OfString: openai.String("How can I get the name of the current day in JavaScript?"),
		},
	},
	option.WithMaxRetries(5),
)
```

### Accessing raw response data (e.g. response headers)

You can access the raw HTTP response data by using the `option.WithResponseInto()` request option. This is useful when
you need to examine response headers, status codes, or other details.

```go
// Create a variable to store the HTTP response
var httpResp *http.Response
response, err := client.Responses.New(
	context.TODO(),
	responses.ResponseNewParams{
		Model: openai.ChatModelGPT5_2,
		Input: responses.ResponseNewParamsInputUnion{
			OfString: openai.String("Say this is a test"),
		},
	},
	option.WithResponseInto(&httpResp),
)
if err != nil {
	// handle error
}
fmt.Printf("%+v\n", response)

fmt.Printf("Status Code: %d\n", httpResp.StatusCode)
fmt.Printf("Headers: %+#v\n", httpResp.Header)
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

When `client` is a native `*http.Client`, the SDK keeps its credential-origin
checks in the redirect path, including when the client uses a custom
`http.RoundTripper`. A bespoke implementation of `option.HTTPClient` owns any
redirects it performs inside `Do` and must keep credentialed requests on the
configured origin. Prefer a native `*http.Client` with a custom transport when
possible.

### Mutual TLS with a custom HTTP client

For API-key authenticated HTTP requests that require mutual TLS, configure a
native Go `*http.Client` and pass it through `option.WithHTTPClient`. The
certificate file must contain the client leaf followed by every required
intermediate. Presenting intermediates requires certificate-chain support to be
enabled for your organization; otherwise, the client certificate must be
signed directly by an active uploaded certificate. See the
[OpenAI mTLS setup requirements](https://help.openai.com/en/articles/10876024):

```go
certificate, err := tls.LoadX509KeyPair(
	"/secrets/openai/client-chain.pem",
	"/secrets/openai/client.key",
)
if err != nil {
	return err
}

defaultTransport, ok := http.DefaultTransport.(*http.Transport)
if !ok {
	return errors.New("http.DefaultTransport is not an *http.Transport")
}
transport := defaultTransport.Clone()
transport.Proxy = nil
transport.DialTLS = nil
transport.DialTLSContext = nil
transport.ResponseHeaderTimeout = 10 * time.Minute
transport.TLSClientConfig = &tls.Config{
	Certificates: []tls.Certificate{certificate},
	GetClientCertificate: func(*tls.CertificateRequestInfo) (*tls.Certificate, error) {
		return &certificate, nil
	},
}

httpClient := &http.Client{
	Transport: transport,
	CheckRedirect: func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

client := openai.NewClient(
	option.WithBaseURL("https://mtls.api.openai.com/v1"),
	option.WithHTTPClient(httpClient),
)

if _, err := client.Models.List(context.Background()); err != nil {
	return err
}
```

The SDK does not select an mTLS endpoint automatically when a custom HTTP
client is used. The explicit `option.WithBaseURL` above overrides
`OPENAI_BASE_URL`; replace it with `https://mtls-eu.api.openai.com/v1` for the
EU endpoint, or remove it to use `OPENAI_BASE_URL`. Keep server trust separate
by configuring `RootCAs` on the fresh `tls.Config` when custom roots are
required.

`tls.LoadX509KeyPair` fails for unreadable files and for malformed or mismatched
leaf/key material. It loads later `CERTIFICATE` blocks into the presented chain
without validating those intermediates. Certificate validity, intermediate
parsing, chain trust, and OpenAI product policy remain TLS-handshake/server
checks. Rebuild the transport and OpenAI client after rotating a certificate
because existing TLS connections cannot renegotiate client authentication.
When overriding the HTTP client, the application also owns redirect, proxy, and
timeout policy. This dedicated client bypasses proxies, retains the SDK's
10-minute response-header timeout, replaces inherited client-certificate
callbacks, TLS dial hooks, and TLS session state with a fresh TLS config, and
disables redirects so the client certificate is only offered to the configured
API endpoint. Its callback always returns the configured certificate because
Go's automatic selection can otherwise suppress it when a server's
acceptable-CA hint does not match the local chain. If a proxy is required, use
a transport that keeps the proxy TLS configuration separate from the origin
client certificate.

The complete tested recipe is in
[`examples/mutual-tls`](./examples/mutual-tls/main.go).

## Workload Identity Authentication

For cloud workloads (Kubernetes, Azure, Google Cloud Platform), you can use workload identity authentication instead of API keys. This provides short-lived tokens that are automatically refreshed.

### Kubernetes

```go
import (
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/auth"
	"github.com/openai/openai-go/v3/option"
)

client := openai.NewClient(
	option.WithWorkloadIdentity(auth.WorkloadIdentity{
		IdentityProviderID: "idp-123",
		ServiceAccountID:   "sa-456",
		Provider:           auth.K8sServiceAccountTokenProvider(""),
	}),
)
```

### Azure Managed Identity

```go
client := openai.NewClient(
	option.WithWorkloadIdentity(auth.WorkloadIdentity{
		IdentityProviderID: "idp-123",
		ServiceAccountID:   "sa-456",
		Provider:           auth.AzureManagedIdentityTokenProvider(nil),
	}),
)
```

### Google Cloud Compute Engine

```go
client := openai.NewClient(
	option.WithWorkloadIdentity(auth.WorkloadIdentity{
		IdentityProviderID: "idp-123",
		ServiceAccountID:   "sa-456",
		Provider:           auth.GCPIDTokenProvider(nil),
	}),
)
```

### Custom Subject Token Provider

You can implement your own subject token provider:

```go
import (
	"context"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/auth"
	"github.com/openai/openai-go/v3/option"
)

type customTokenProvider struct{}

func (p *customTokenProvider) TokenType() auth.SubjectTokenType {
	return auth.SubjectTokenTypeJWT
}

func (p *customTokenProvider) GetToken(ctx context.Context, httpClient auth.HTTPDoer) (string, error) {
	return "your-token", nil
}

client := openai.NewClient(
	option.WithWorkloadIdentity(auth.WorkloadIdentity{
		IdentityProviderID: "idp-123",
		ServiceAccountID:   "sa-456",
		Provider:           &customTokenProvider{},
	}),
)
```

### Customizing Refresh Buffer

By default, tokens are refreshed 20 minutes (1200 seconds) before expiry. You can customize this:

```go
client := openai.NewClient(
	option.WithWorkloadIdentity(auth.WorkloadIdentity{
		IdentityProviderID:   "idp-123",
		ServiceAccountID:     "sa-456",
		Provider:             auth.K8sServiceAccountTokenProvider(""),
		RefreshBufferSeconds: 600,
	}),
)
```

### X.509 Workload Identity

An enrolled workload certificate can also be exchanged directly for a
short-lived OpenAI bearer token. Load and own the certificate and private key in
your application, configure one static certificate on a native transport, and
attest that transport before creating the client:

```go
import (
	"context"
	"crypto/tls"
	"net/http"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/auth"
	"github.com/openai/openai-go/v3/option"
)

certificate, err := tls.LoadX509KeyPair("workload.crt", "workload.key")
if err != nil {
	return err
}

transport, err := auth.NewX509Transport(&http.Transport{
	TLSClientConfig: &tls.Config{
		MinVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{certificate},
	},
})
if err != nil {
	return err
}
defer transport.Close()

client := openai.NewClient(option.WithX509WorkloadIdentity(auth.X509WorkloadIdentity{
	IdentityProviderID: "idp-123",
	ServiceAccountID:   "sa-456",
	RefreshBuffer:      5 * time.Minute, // Optional; five minutes is the default.
	Transport:          transport,
}))

requestContext, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
defer cancel()

if _, err := client.Models.List(requestContext); err != nil {
	return err
}
```

The application retains ownership of its certificate, private key, trust roots,
and original transport. The attested capability creates its own isolated
connection pool, so call `Close` when it is no longer needed. `Close` rejects
new calls but does not cancel or wait for requests already in progress; an
in-progress request may finish after `Close` returns. Rotating the certificate
requires a newly attested transport and client.

Transport attestation validates the native transport and TLS policy and
snapshots the certificate bytes. It is not proof that the certificate is
currently valid or enrolled, or that an arbitrary manually assembled private
key matches it; the TLS handshake and issuer enforce those properties. The
capability retains the `crypto.Signer` supplied at construction, so a
hardware-backed or custom signer must remain usable and safe for concurrent
TLS handshakes. Keep the original transport and its TLS configuration unchanged
while the capability is in use; their attested policy is revalidated before
requests.

The capability always presents its one attested certificate, including when a
server advertises unrelated acceptable certificate-authority hints. It does not
accept caller-provided certificate-selection callbacks. When the transport
template does not specify its own values, the isolated connection pool applies
a 30-second TCP dial timeout, a 10-second TLS handshake timeout, and the SDK's
10-minute response-header timeout. Each `GetToken` call without a caller
deadline, including time spent waiting for a concurrent exchange and all issuer
retries, uses a 30-second default. An explicit caller deadline is authoritative,
even when later. Use `option.WithRequestTimeout` to bound each request attempt,
including token acquisition and response handling, or a caller context deadline
to bound the entire call across retries. The default exchange timeout does not
add a deadline to API response bodies, so long-running streams remain supported.

Token exchange is pinned to `https://mtls.auth.openai.com/oauth/token`; API
requests use `https://mtls.api.openai.com/v1/`. Existing `OPENAI_BASE_URL` and
explicit endpoint settings must match that global API endpoint; an explicit
default `:443` port is equivalent. Azure, Amazon Bedrock, regional endpoints,
HTTPS proxies, dynamic certificate selection, HTTP trace hooks, and separate
custom HTTP clients are not supported. A legacy `http.Transport.Dial` remains
supported for compatibility when `DialContext` is unset, but its in-progress
dial call cannot itself be interrupted by context cancellation. The capability
admits at most 32 concurrent dial calls; additional live requests wait for
admission, while canceled or completed requests release their waiters. This
prevents a non-cooperative custom dialer from growing without bound. The native
transport still applies `MaxConnsPerHost` independently to each host.
Organization and project metadata remain available on API requests but are not
sent to the token issuer.

Successful bearer tokens are cached per identity and transport generation.
Concurrent refreshes share the requesting caller's context, and the effective
refresh buffer is reduced automatically for short-lived tokens. Transient issuer
failures receive at most three bounded, cancellable attempts; permanent OAuth
failures are not retried. Issuer retries, ordinary API retries, and unauthorized
recovery share the client's single request retry budget. A rejected API token is
refreshed and replayed at most once, and only when the request body can be
recreated; body-bearing requests with caller middleware are not replayed because
that middleware may have transformed their bytes. During a temporary proactive
refresh failure, an unexpired cached token remains usable for a bounded cooldown.
Bearer generations rejected by the API remain tombstoned until their original
expiry, including late rejections that arrive after a replacement is cached.
An identity retains at most 1,024 simultaneously unexpired bearer generations
and fails closed before publishing a bearer whose expiry cannot be tracked.
The resulting access
token is an ordinary bearer token; it is not cryptographically bound to a
client certificate. Rotating or revoking a certificate does not by itself
invalidate bearer tokens already minted from that certificate; server-side
identity mapping, token revocation, and audit policy are outside the SDK. Treat
both the private key and minted bearers as credentials.

Issuer-provided OAuth error descriptions are intentionally redacted because
they may contain sensitive details. Use `OAuthError.StatusCode` and
`OAuthError.ErrorCode` for programmatic handling. `OAuthError.ErrorDescription`
is retained for compatibility but is normally empty.

## Amazon Bedrock

Use the `bedrock` package to call OpenAI models through Amazon Bedrock's
OpenAI-compatible API. The standard AWS SDK credential chain is used by
default, so existing environment credentials, `~/.aws/credentials`,
`~/.aws/config`, SSO or assume-role profiles, and workload credentials work
without custom signing code.

```go
package main

import (
	"context"
	"log"

	"github.com/openai/openai-go/v3/bedrock"
)

func main() {
	client, err := bedrock.NewClient(context.Background(), bedrock.Config{
		AWSRegion: "us-west-2",
	})
	if err != nil {
		log.Fatal(err)
	}

	_ = client
}
```

The region is resolved from `AWSRegion`, `AWS_REGION`, `AWS_DEFAULT_REGION`,
or the standard AWS config chain. The base URL is resolved from `BaseURL`,
`AWS_BEDROCK_BASE_URL`, or
`https://bedrock-mantle.{region}.api.aws/openai/v1`.
The `/openai/v1` prefix is intentional; Bedrock's generic `/v1` route is a
different API surface and is not interchangeable with the OpenAI-compatible
route.

To select a named profile:

```go
client, err := bedrock.NewClient(context.Background(), bedrock.Config{
	AWSProfile: "production",
})
```

You can also provide temporary static credentials or an AWS SDK v2
`aws.CredentialsProvider` through `bedrock.Config`. Prefer roles, profiles, and
other temporary credential sources over long-lived static credentials.

Bedrock bearer credentials remain supported through `APIKey`,
`BedrockTokenProvider`, or `AWS_BEARER_TOKEN_BEDROCK`:

```go
client, err := bedrock.NewClient(context.Background(), bedrock.Config{
	AWSRegion: "us-west-2",
	APIKey:    "bedrock-bearer-token",
})
```

Explicit bearer and AWS credential modes are mutually exclusive. AWS SigV4
authentication signs the fully serialized request again on every retry using
the `bedrock-mantle` service name. Request bodies must be replayable; response
streaming is supported. Authenticated Bedrock requests do not automatically
follow redirects.
Ambient `OPENAI_*` credentials, routing, and headers are not inherited by a
Bedrock client.

See [`examples/bedrock`](examples/bedrock) for a complete Responses API example.

## Azure OpenAI in Azure AI Foundry Models

To use this library with [Azure OpenAI in Azure AI Foundry Models](https://learn.microsoft.com/azure/ai-services/openai/overview),
use the option.RequestOption functions in the `azure` package.

```go
package main

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/azure"
)

func main() {
	const azureOpenAIEndpoint = "https://<azure-openai-resource>.openai.azure.com"

	// The latest API versions, including previews, can be found here:
	// https://learn.microsoft.com/en-us/azure/ai-services/openai/reference#rest-api-versioning
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

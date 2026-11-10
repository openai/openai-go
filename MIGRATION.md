# OpenAI Go Migration Guide

<a href="https://pkg.go.dev/github.com/openai/openai-go/v3"><img src="https://pkg.go.dev/badge/github.com/openai/openai-go.svg" alt="Go Reference"></a>

## Go v4 breaking change: Azure Identity version selection

If your application uses `azidentity`, declare the version it needs in your own
`go.mod`. The SDK no longer sets the minimum to `v1.14.1`. Your other
dependencies can select an older version. Run the following commands to keep the previous minimum.

```sh
go get github.com/Azure/azure-sdk-for-go/sdk/azidentity@v1.14.1
go mod tidy
```

To choose the latest release instead, use `@latest`. Continue passing your
credential to `azure.WithTokenCredential`. Use `azure.WithTokenCredentialScopes`
if your deployment requires a custom scope. Applications using an Azure API key
can continue using `azure.WithAPIKey` without installing Azure Identity.

## Retired chat model constants in v4

Go v4 removes three constants from `shared` and their aliases in `openai`,
`conversations`, `live`, `realtime`, `responses`, and `webhooks`.

| Removed constant | Replacement to evaluate |
| --- | --- |
| `ChatModelGPT4VisionPreview` | `openai.ChatModelGPT6Luna` for requests with image input |
| `ChatModelGPT3_5Turbo0613` | `openai.ChatModelGPT5_6Terra` |
| `ChatModelGPT3_5Turbo16k0613` | `openai.ChatModelGPT5_6Terra` |

The original retirement notices recommended GPT-4o and GPT-3.5 Turbo.
Subsequent notices list the replacements above. Test your application's
messages with the model you select. You can also set
`ChatCompletionNewParams.Model` from a string variable.
[API model deprecations](https://developers.openai.com/api/docs/deprecations)

## Go v4 removes the Videos API

The Sora API shut down on September 24, 2026. Remove the video-generation feature
from your application or choose another service. There is no direct replacement
in the OpenAI API. Installing an older SDK cannot restore access.

Remove calls to `client.Videos`, including `NewAndPoll`, `PollStatus` and
`ListAutoPaging`. Remove references to `openai.NewVideoService`, the video
request and response types, `openai.ImageInputReferenceParam`, and Sora model
constants. Remove batch submissions that use `openai.BatchNewParamsEndpointV1Videos`.

## Earlier SDK changes

This SDK includes breaking changes to improve the ergonomics of constructing parameters and accessing responses.

To reduce verbosity, the `openai.F(...)` and `param.Field[T]` have been removed.
All calls to `openai.F(...)` can be deleted.

The SDK now uses the `` `json:"...,omitzero"` `` struct tag to omit fields. Nested structs, arrays and maps
can be declared like normal.

The old SDK used interfaces for unions in requests, which required
a type assertion to access variants and fields. The new design uses
structs with a field for each variant, wherein only one field can be set.
These struct unions also expose 'Get' methods to access and mutate subfields
which may be shared by multiple variants.

## Typed headers and per-call options in v4

Go v4 applies service defaults first, typed header parameters second, and per-call options last. A call that supplies a typed `IdempotencyKey` of `"from-params"` and `option.WithHeader("Idempotency-Key", "from-option")` sends `"from-option"`. In v3 it sends `"from-params"`.

Choose one source for each idempotency key. To keep the typed key, remove the conflicting per-call option. Reuse the same key for retries of one operation and use a distinct key for each new operation.

Check shared options for `WithHeaderDel` and `WithHeaderAdd` as well. Deleting a typed header removes it from the request. Adding values to a repeated header such as `OpenAI-Beta` places them after the typed values, in the order you supplied them.

# Request parameters

## Required primitives parameters serialize their zero values (`string`, `int64`, etc.)

> [!CAUTION]
>
> **This change can cause new behavior in existing code, without compiler warnings.**

While migrating, ensure that all required fields are explicitly set. A required primitive
field `Age` will use the `` `json:"age,required"` `` struct tag without `omitzero`.

If a required primitive field is not set, the zero value will be serialized.
This was not the case in with `param.Field[T]`.

```diff
type FooParams struct {
-        Age  param.Field[int64]  `json:"age,required"`
-        Name param.Field[string] `json:"name"`
+        Age  int64               `json:"age,required"` // <== Notice no omitzero
+        Name param.Opt[string]   `json:"name,omitzero"`
}
```

<table>
<tr>
<th>Previous</th>
<th>New</th>
</tr>
<tr>
<td>

```go
_ = FooParams{
    Name: openai.String("Jerry")
}
`{"name": "Jerry"}` // (after serialization)
```

</td>
<td>

```go
_ = FooParams{
    Name: openai.String("Jerry")
}
`{"name": "Jerry", "age": 0}` // <== Notice the age field
```

</td>
</tr>
</table>

The required field `"age"` is now present as `0`. Fields without the `` `json:"...,omitzero"` `` struct tag
are always serialized, including their zero values.

## Transition from `param.Field[T]` to `omitzero`

The `openai.F(...)` function and `param.Field[T]` type are no longer present in the new SDK.

To represent omitted fields, the SDK uses <a href="https://pkg.go.dev/encoding/json#Marshal">`` `json:"...,omitzero"` `` semantics</a> from Go 1.24+ for JSON encoding[^1]. `omitzero` always omits fields
with zero values.

In all cases other than optional primitives, `openai.F()` can simply be removed.
For optional primitive types, such as `param.Opt[string]`, you can use `openai.String(string)` to construct the value.
Similar functions exist for other primitive types like `openai.Int(int)`, `openai.Bool(bool)`, etc.

`omitzero` is used for fields whose type is either a struct, slice, map, string enum,
or wrapped optional primitive (e.g. `param.Opt[T]`). Required primitive fields don't use `omitzero`.

**Example User Code: Constructing a request**

```diff
foo = FooParams{
-    RequiredString: openai.String("hello"),
+    RequiredString: "hello",

-    OptionalString: openai.String("hi"),
+    OptionalString: openai.String("hi"),

-    Array: openai.F([]BarParam{
-        BarParam{Prop: ... }
-    }),
+    Array: []BarParam{
+        BarParam{Prop: ... }
+    },

-    RequiredObject: openai.F(BarParam{ ... }),
+    RequiredObject: BarParam{ ... },

-    OptionalObject: openai.F(BarParam{ ... }),
+    OptionalObject: BarParam{ ... },

-    StringEnum: openai.F[BazEnum]("baz-ok"),
+    StringEnum: "baz-ok",
}
```

**Internal SDK Code: Fields of a request struct:**

```diff
type FooParams struct {
-    RequiredString param.Field[string]   `json:"required_string,required"`
+    RequiredString string                `json:"required_string,required"`

-    OptionalString param.Field[string]   `json:"optional_string"`
+    OptionalString param.Opt[string]     `json:"optional_string,omitzero"`

-    Array param.Field[[]BarParam]        `json"array"`
+    Array []BarParam                     `json"array,omitzero"`

-    Map param.Field[map[string]BarParam] `json"map"`
+    Map map[string]BarParam              `json"map,omitzero"`

-    RequiredObject param.Field[BarParam] `json:"required_object,required"`
+    RequiredObject BarParam              `json:"required_object,omitzero,required"`

-    OptionalObject param.Field[BarParam] `json:"optional_object"`
+    OptionalObject BarParam              `json:"optional_object,omitzero"`

-    StringEnum     param.Field[BazEnum]  `json:"string_enum"`
+    StringEnum     BazEnum               `json:"string_enum,omitzero"`
}
```

## Request Unions: Removing interfaces and moving to structs

For a type `AnimalUnionParam` which could be either a `CatParam | DogParam`.

<table>
<tr><th>Previous</th> <th>New</th></tr>
<tr>
<td>

```go
type AnimalParam interface {
	ImplAnimalParam()
}

func (Dog)         ImplAnimalParam() {}
func (Cat)         ImplAnimalParam() {}
```

</td>
<td>

```go
type AnimalUnionParam struct {
	OfCat 	 *Cat              `json:",omitzero,inline`
	OfDog    *Dog              `json:",omitzero,inline`
}
```

</td>
</tr>

<tr style="background:rgb(209, 217, 224)">
<td>

```go
var dog AnimalParam = DogParam{
	Name: "spot", ...
}
var cat AnimalParam = CatParam{
	Name: "whiskers", ...
}
```

</td>
<td>

```go
dog := AnimalUnionParam{
	OfDog: &DogParam{Name: "spot", ... },
}
cat := AnimalUnionParam{
	OfCat: &CatParam{Name: "whiskers", ... },
}
```

</td>
</tr>

<tr>
<td>

```go
var name string
switch v := animal.(type) {
case Dog:
	name = v.Name
case Cat:
	name = v.Name
}
```

</td>
<td>

```go
// Accessing fields
var name *string = animal.GetName()
```

</td>
</tr>
</table>

## Sending explicit `null` values

The old SDK had a function `param.Null[T]()` which could set `param.Field[T]` to `null`.

The new SDK uses `param.Null[T]()` for to set a `param.Opt[T]` to `null`,
but `param.NullStruct[T]()` to set a param struct `T` to `null`.

```diff
- var nullPrimitive param.Field[int64] = param.Null[int64]()
+ var nullPrimitive param.Opt[int64]   = param.Null[int64]()

- var nullStruct param.Field[BarParam] = param.Null[BarParam]()
+ var nullStruct BarParam              = param.NullStruct[BarParam]()
```

## Sending custom values

The `openai.Raw[T](any)` function has been removed. All request structs now support a
`.WithExtraField(map[string]any)` method to customize the fields.

```diff
foo := FooParams{
     A: param.String("hello"),
-    B: param.Raw[string](12) // sending `12` instead of a string
}
+ foo.SetExtraFields(map[string]any{
+    "B": 12,
+ })
```

# Response Properties

## Checking for presence of optional fields

The `.IsNull()` method has been changed to `.Valid()` to better reflect its behavior.

```diff
- if !resp.Foo.JSON.Bar.IsNull() {
+ if resp.Foo.JSON.Bar.Valid() {
    println("bar is present:", resp.Foo.Bar)
}
```

| Previous       | New                      | Returns true for values |
| -------------- | ------------------------ | ----------------------- |
| `.IsNull()`    | `!.Valid()`              | `null` or Omitted       |
| `.IsMissing()` | `.Raw() == resp.Omitted` | Omitted                 |
|                | `.Raw() == resp.Null`    |

## Checking Raw JSON of a response

The `.RawJSON()` method has moved to the parent of the `.JSON` property.

```diff
- resp.Foo.JSON.RawJSON()
+ resp.Foo.RawJSON()
```

[^1]: The SDK doesn't require Go 1.24, despite supporting the `omitzero` feature

## Workload identity token reuse

Reusing `option.WithWorkloadIdentity` across clients no longer shares one token
cache across different HTTP identities. Each native HTTP client and transport
pair keeps its own cache; request-level `option.WithHTTPClient` overrides select
the corresponding identity. A rejected token invalidates only that identity's
cache. The public method and option signatures are unchanged.

Applications that previously relied on token reuse across different clients may
make additional token exchanges. Reuse a stable client and transport when they
represent the same identity. Native HTTP clients also retain tokens when their
transport is a function or other non-comparable value; replacing that transport
selects a separate cache. Non-comparable custom HTTP doers remain supported,
but acquire tokens independently for each logical SDK request while retaining
authentication state across that request’s retry attempts. The cache retains at most 32 idle
or active identity entries; eviction or overflow may also require another
exchange without limiting the number of clients you can use. An overflowing
identity retains its authentication state across attempts of one SDK request.

Keep provider and custom-doer identities stable while in use. To rotate opaque
credentials, construct a new provider/authenticator or HTTP doer. To rotate
client certificates, use a new transport. Do not mutate a client or transport
concurrently with requests.

## Live Responses delegation tool choices

Live Responses delegation function and MCP tool choices have shorter Go names. The
request and response variants now have separate types. All names below are in the
`live` package.

| Previous type | Request type | Response type |
| --- | --- | --- |
| `ResponsesDelegationConfigToolChoiceLiveFunctionToolChoiceParam` | `ResponsesDelegationConfigToolChoiceFunctionParam` | `ResponsesDelegationConfigToolChoiceFunction` |
| `ResponsesDelegationConfigToolChoiceLiveMcpToolChoiceParam` | `ResponsesDelegationConfigToolChoiceMcpParam` | `ResponsesDelegationConfigToolChoiceMcp` |
| `ResponsesDelegationUpdateConfigToolChoiceLiveFunctionToolChoiceParam` | `ResponsesDelegationUpdateConfigToolChoiceFunctionParam` | — |
| `ResponsesDelegationUpdateConfigToolChoiceLiveMcpToolChoiceParam` | `ResponsesDelegationUpdateConfigToolChoiceMcpParam` | — |

The create/response rows below apply to `ResponsesDelegationConfigToolChoiceUnionParam`
and `ResponsesDelegationConfigToolChoiceUnion`, including the response union's `JSON`
metadata. The update rows apply to `ResponsesDelegationUpdateConfigToolChoiceUnionParam`.

| Union | Previous field | New field |
| --- | --- | --- |
| Create/response | `OfResponsesDelegationConfigToolChoiceLiveFunctionToolChoiceParam` | `OfResponsesDelegationConfigToolChoiceFunction` |
| Create/response | `OfResponsesDelegationConfigToolChoiceLiveMcpToolChoiceParam` | `OfResponsesDelegationConfigToolChoiceMcp` |
| Update | `OfResponsesDelegationUpdateConfigToolChoiceLiveFunctionToolChoiceParam` | `OfResponsesDelegationUpdateConfigToolChoiceFunction` |
| Update | `OfResponsesDelegationUpdateConfigToolChoiceLiveMcpToolChoiceParam` | `OfResponsesDelegationUpdateConfigToolChoiceMcp` |

On `ResponsesDelegationConfigToolChoiceUnion`, replace
`AsResponsesDelegationConfigToolChoiceLiveFunctionToolChoiceParam()` with
`AsResponsesDelegationConfigToolChoiceFunction()` and
`AsResponsesDelegationConfigToolChoiceLiveMcpToolChoiceParam()` with
`AsResponsesDelegationConfigToolChoiceMcp()`. These accessors return the response
types in the first table.

MCP tool-choice `Name` request fields now use `param.Opt[string]`. This applies to
the two MCP request types above and to
`ResponsesDelegationConfigToolChoiceLiveAllowedToolsChoiceParamToolMcpParam` and
`ResponsesDelegationUpdateConfigToolChoiceLiveAllowedToolsChoiceParamToolMcpParam`.
Use `openai.String("lookup")` for a string, `param.Null[string]()` for JSON `null`, or
leave `Name` unset to omit it. `ServerLabel` remains required. A request union's
`GetName()` returns `nil` for both an omitted and a null name.

**Before**

```go
choice := live.ResponsesDelegationConfigToolChoiceUnionParam{
	OfResponsesDelegationConfigToolChoiceLiveMcpToolChoiceParam: &live.ResponsesDelegationConfigToolChoiceLiveMcpToolChoiceParam{
		ServerLabel: "docs",
		Name:        "lookup",
	},
}
```

**After**

```go
choice := live.ResponsesDelegationConfigToolChoiceUnionParam{
	OfResponsesDelegationConfigToolChoiceMcp: &live.ResponsesDelegationConfigToolChoiceMcpParam{
		ServerLabel: "docs",
		Name:        openai.String("lookup"),
	},
}
```

Both requests select the `lookup` tool on the `docs` MCP server.

For responses, `Name` is a string. Check `JSON.Name.Valid()` for a value and
`JSON.Name.Raw()` against `respjson.Omitted` or `respjson.Null` to distinguish an
omitted field from JSON `null`. Use `RawJSON()` to retain the original wire form.

## Retained messages in compaction output

GA and beta compaction results can contain retained input messages. If you
declare the type of the output slice, change the GA type from
`[]responses.ResponseOutputItemUnion` to `[]responses.ResponseItemUnion`
and the beta type from `[]openai.BetaResponseOutputItemUnion` to
`[]openai.BetaResponseItemUnion`.

For a message, `AsAny()` selects the input-message type for the `user`,
`system`, and `developer` roles and the output-message type for
`assistant`. For GA, handle both variants when iterating over the output.

```go
for _, item := range compacted.Output {
    switch msg := item.AsAny().(type) {
    case responses.ResponseInputMessageItem:
        fmt.Println(msg.Role, msg.Content)
    case responses.ResponseOutputMessage:
        fmt.Println(msg.Role, msg.Content)
    }
}
```

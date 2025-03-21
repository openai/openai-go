# OpenAI Go Migration Guide

<a href="https://pkg.go.dev/github.com/openai/openai-go"><img src="https://pkg.go.dev/badge/github.com/openai/openai-go.svg" alt="Go Reference"></a>

This SDK includes breaking changes from the previous version to improve the ergonomics of constructing parameters and accessing responses.

# Request parameters

## Required primitives parameters serialize their zero values (`string`, `int64`, etc.)

> [!CAUTION] > **This change can cause new behavior in existing code, without compiler warnings.**

```diff
type FooParams struct {
-        Age  param.Field[int64]  `json:"age,required"`
-        Name param.Field[string] `json:"name"`
+        Age  int64               `json:"age,required"`
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

The required field `"age"` is now present as `0`. Required primitive fields without the <code>\`json:"...,omitzero"\`</code> struct tag
are always serialized, including their zero values.

## Transition from `param.Field[T]` to `omitzero`

The new SDK uses <a href="https://pkg.go.dev/encoding/json#Marshal"><code>\`json:"...,omitzero"\`</code> semantics</a> from Go 1.24+ for JSON encoding[^1].

`omitzero` is used for structs, slices, maps, string enums, and optional primitive types wrapped in `param.Opt[T]` (e.g. `param.Opt[string]`).

**Fields of a request struct:**

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

**Previous vs New SDK: Constructing a request**

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

`param.Opt[string]` can be constructed with `openai.String(string)`. Similar functions exist for other primitive
types like `openai.Int(int)`, `openai.Bool(bool)`, etc.

## Request Unions: Removing interfaces and moving to structs

For a type `AnimalUnionParam` which could be either a `string | CatParam | DogParam`.

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

The new SDK uses `param.NullOpt[T]()` for to set a `param.Opt[T]` to `null`,
and `param.NullObj[T]()` to set a param struct `T` to `null`.

```diff
- var nullObj param.Field[BarParam] = param.Null[BarParam]()
+ var nullObj BarParam              = param.NullObj[BarParam]()

- var nullPrimitive param.Field[int64] = param.Null[int64]()
+ var nullPrimitive param.Opt[int64]   = param.NullOpt[int64]()
```

## Sending custom values

The `openai.Raw[T](any)` function has been removed. All request structs now support a
`.WithExtraField(map[string]any)` method to customize the fields.

```diff
foo := FooParams{
     A: param.String("hello"),
-    B: param.Raw[string](12) // sending `12` instead of a string
}
+ foo.WithExtraFields(map[string]any{
+    "B": 12,
+ })
```

# Response Properties

## Checking for presence of optional fields

The `.IsNull()` method has been changed to `.IsPresent()` to better reflect its behavior.

```diff
- if !resp.Foo.JSON.Bar.IsNull() {
+ if resp.Foo.JSON.Bar.IsPresent() {
     println("bar is present:", resp.Foo.Bar)
  }
```

| Previous       | New                 | Returns true for values |
| -------------- | ------------------- | ----------------------- |
| `.IsNull()`    | `!.IsPresent()`     | `null` or Omitted       |
| `.IsMissing()` | `.Raw() == ""`      | Omitted                 |
| `.Invalid()`   | `.IsExplicitNull()` | `null`                  |

[^1]: The SDK doesn't require Go 1.24, despite supporting the `omitzero` feature

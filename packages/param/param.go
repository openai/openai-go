package param

import (
	"encoding/json"
	"github.com/openai/openai-go/internal/apifield"
	"reflect"
	"time"
)

func IsOmitted(v any) bool {
	if i, ok := v.(interface{ IsOmitted() bool }); ok {
		return i.IsOmitted()
	}
	return reflect.ValueOf(v).IsZero()
}

type Fielder interface {
	IsNull() bool
	IsOverridden() (any, bool)
	RawResponse() json.RawMessage
}

type ObjectFielder interface {
	Fielder
	IsFieldNull(string) bool
	IsFieldOverridden(string) (any, bool)
	IsFieldRawResponse(string) json.RawMessage
}

// This pattern allows mutable generics, no code should require that this type be provided.
type SettableFielder[T Fielder] interface {
	setMetadata(MetadataProvider)
	*T
}

// Override the field with a custom json value, the v parameter uses
// the same semantics as json.Marshal from encoding/json.
//
// The SettableFielder type parameter should never be provided, it is always inferred.
//
//	var f param.String = param.Override[param.String](12)
//	json.Marshal(f) == `12`
func Override[T Fielder, PT SettableFielder[T]](v any) T {
	var x T
	PT(&x).setMetadata(apifield.CustomValue{Override: v})
	return x
}

// Set the field to null
//
// The SettableFielder type parameter should never be provided, it is always inferred.
//
//	var f param.String = param.Null[param.String]()
//	json.Marshal(f) == `null`
func Null[T Fielder, PT SettableFielder[T]]() T {
	var x T
	PT(&x).setMetadata(apifield.ExplicitNull{})
	return x
}

// Constructs a field whose zero value is never omitted.
// This is useful for internal code, there are other preferred ways to construct non-omitted fields.
func NeverOmitted[T Fielder, PT SettableFielder[T]]() T {
	var x T
	PT(&x).setMetadata(apifield.NeverOmitted{})
	return x
}

// MarshalExtraField adds an additional field to be set with custom json. The v parameter
// uses the same semantics as json.Marshal from encoding/json.
// If any native field with a matching json field name will be zeroed and omitted.
// func InsertFields(obj MutableFieldLike, k string, v any) {}

// UnmarshalExtraField accesses an extra field and unmarshals the result in the value pointed to by v.
// UnmarshalExtraField uses similar semantics to json.Unmarshal from encoding/json. However,
// if v is nil or not a pointer, or the extra field cannot be unmarshaled into the the value v,
// UnmarshalExtraField returns false.
// func GetField(obj ObjectFieldLike, k string, v any) (exists bool) {
//   return apifield.UnmarshalExtraField(obj, k, v)
// }

type String struct {
	V string
	metadata
}

type Int struct {
	V int64
	metadata
}

type Bool struct {
	V bool
	metadata
}

type Float struct {
	V float64
	metadata
}

type Datetime struct {
	V time.Time
	metadata
}

type Date struct {
	V time.Time
	metadata
}

// Either null or omitted
func (f String) IsMissing() bool                 { return IsOmitted(f) || f.IsNull() }
func (f String) IsOmitted() bool                 { return f == String{} }
func (f String) MarshalJSON() ([]byte, error)    { return marshalField(f, f.V) }
func (f String) UnmarshalJSON(data []byte) error { return unmarshalField(&f.V, &f.metadata, data) }
func (f String) String() string                  { return stringifyField(f, f.V) }

// Either null or omitted
func (f Int) IsMissing() bool                 { return f.IsOmitted() || f.IsNull() }
func (f Int) IsOmitted() bool                 { return f == Int{} }
func (f Int) MarshalJSON() ([]byte, error)    { return marshalField(f, f.V) }
func (f Int) UnmarshalJSON(data []byte) error { return unmarshalField(&f.V, &f.metadata, data) }
func (f Int) String() string                  { return stringifyField(f, f.V) }

// Either null or omitted
func (f Bool) IsMissing() bool                 { return f.IsOmitted() || f.IsNull() }
func (f Bool) IsOmitted() bool                 { return f == Bool{} }
func (f Bool) MarshalJSON() ([]byte, error)    { return marshalField(f, f.V) }
func (f Bool) UnmarshalJSON(data []byte) error { return unmarshalField(&f.V, &f.metadata, data) }
func (f Bool) String() string                  { return stringifyField(f, f.V) }

// Either null or omitted
func (f Float) IsMissing() bool                 { return f.IsOmitted() || f.IsNull() }
func (f Float) IsOmitted() bool                 { return f == Float{} }
func (f Float) MarshalJSON() ([]byte, error)    { return marshalField(f, f.V) }
func (f Float) UnmarshalJSON(data []byte) error { return unmarshalField(&f.V, &f.metadata, data) }
func (f Float) String() string                  { return stringifyField(f, f.V) }

// Either null or omitted
func (f Datetime) IsMissing() bool                 { return f.IsOmitted() || f.IsNull() }
func (f Datetime) IsOmitted() bool                 { return f == Datetime{} }
func (f Datetime) MarshalJSON() ([]byte, error)    { return marshalField(f, f.V.Format(time.RFC3339)) }
func (f Datetime) UnmarshalJSON(data []byte) error { return unmarshalField(&f.V, &f.metadata, data) }
func (f Datetime) String() string                  { return stringifyField(f, f.V.Format(time.RFC3339)) }

// Either null or omitted
func (f Date) IsMissing() bool                 { return f.IsOmitted() || f.IsNull() }
func (f Date) IsOmitted() bool                 { return f == Date{} }
func (f Date) MarshalJSON() ([]byte, error)    { return marshalField(f, f.V.Format("2006-01-02")) }
func (f Date) UnmarshalJSON(data []byte) error { return unmarshalField(&f.V, &f.metadata, data) }
func (f Date) String() string                  { return stringifyField(f, f.V.Format("2006-01-02")) }

// APIObject should be embedded in api object fields, preferably using an alias to make private
type APIObject struct {
	metadata
}

// APIUnion should be embedded in all api unions fields, preferably using an alias to make private
type APIUnion struct {
	metadata
}

func (o APIObject) IsFieldNull(string) bool                   { return false }
func (o APIObject) IsFieldOverridden(string) (any, bool)      { return nil, false }
func (o APIObject) IsFieldRawResponse(string) json.RawMessage { return nil }

type metadata struct {
	// provider is an interface used to determine the status of the field.
	// As an optimization, we expect certain concrete types.
	//
	// While there are simpler ways to implement metadata, the primary incentive here is to
	// minimize the bytes in the struct, since it will be embedded in every field.
	provider MetadataProvider
}

func (m metadata) IsNull() bool {
	if m.provider == nil {
		return false
	}

	// avoid dynamic dispatch call for the most common cases
	if _, ok := m.provider.(apifield.NeverOmitted); ok {
		return false
	} else if _, ok := m.provider.(apifield.ExplicitNull); ok {
		return true
	}

	return m.provider.IsNull()
}

func (m metadata) RawResponse() json.RawMessage {
	if m.provider == nil {
		return nil
	}
	// avoid dynamic dispatch call for the most common case
	if r, ok := m.provider.(apifield.ResponseData); ok {
		return json.RawMessage(r)
	}
	return m.provider.RawResponse()
}

func (m metadata) IsOverridden() (any, bool) {
	if m.provider == nil {
		return nil, false
	}
	return m.provider.IsOverridden()
}

func (m *metadata) setMetadata(mp MetadataProvider) {
	m.provider = mp
}

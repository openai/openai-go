package param_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/Nordlys-Labs/openai-go/v3/packages/param"
)

type Struct struct {
	A string `json:"a"`
	B int64  `json:"b"`
	param.APIObject
}

func (r Struct) MarshalJSON() (data []byte, err error) {
	type shadow Struct
	return param.MarshalObject(r, (*shadow)(&r))
}

// Note that the order of fields affects the JSON
// key order. Changing the order of the fields in this struct
// will fail tests unnecessarily.
type FieldStruct struct {
	A param.Opt[string]    `json:"a,omitzero"`
	B param.Opt[int64]     `json:"b,omitzero"`
	C Struct               `json:"c,omitzero"`
	D time.Time            `json:"d,omitzero" format:"date"`
	E time.Time            `json:"e,omitzero"`
	F param.Opt[time.Time] `json:"f,omitzero" format:"date"`
	G param.Opt[time.Time] `json:"g,omitzero"`
	H param.Opt[time.Time] `json:"h,omitzero" format:"date-time"`
	param.APIObject
}

func (r FieldStruct) MarshalJSON() (data []byte, err error) {
	type shadow FieldStruct
	return param.MarshalObject(r, (*shadow)(&r))
}

type StructWithAdditionalProperties struct {
	First       string         `json:"first"`
	Second      int            `json:"second"`
	ExtraFields map[string]any `json:"-"`
	param.APIObject
}

func (s StructWithAdditionalProperties) MarshalJSON() ([]byte, error) {
	type shadow StructWithAdditionalProperties
	return param.MarshalWithExtras(s, (*shadow)(&s), s.ExtraFields)
}

func TestIsNullish(t *testing.T) {
	nullTests := map[string]param.ParamNullable{
		"null_string": param.Null[string](),
		"null_int64":  param.Null[int64](),
		"null_time":   param.Null[time.Time](),
		"null_struct": param.NullStruct[Struct](),
	}

	for name, test := range nullTests {
		t.Run(name, func(t *testing.T) {
			if !param.IsNull(test) {
				t.Fatalf("expected %s to be null", name)
			}
			if param.IsOmitted(test) {
				t.Fatalf("expected %s to not be omitted", name)
			}
		})
	}

	omitTests := map[string]param.ParamNullable{
		"omit_string": param.Opt[string]{},
		"omit_int64":  param.Opt[int64]{},
		"omit_time":   param.Opt[time.Time]{},
		"omit_struct": Struct{},
	}

	for name, test := range omitTests {
		t.Run(name, func(t *testing.T) {
			if param.IsNull(test) {
				t.Fatalf("expected %s to be null", name)
			}
			if !param.IsOmitted(test) {
				t.Fatalf("expected %s to not be omitted", name)
			}
		})
	}
}

func TestFieldMarshal(t *testing.T) {
	tests := map[string]struct {
		value    any
		expected string
	}{
		"null_string": {param.Null[string](), "null"},
		"null_int64":  {param.Null[int64](), "null"},
		"null_time":   {param.Null[time.Time](), "null"},
		"null_struct": {param.NullStruct[Struct](), "null"},

		"float_zero":  {param.NewOpt(float64(0.0)), "0"},
		"string_zero": {param.NewOpt(""), `""`},
		"time_zero":   {param.NewOpt(time.Time{}), `"0001-01-01T00:00:00Z"`},

		"string": {param.Opt[string]{Value: "string"}, `"string"`},
		"int":    {param.Opt[int64]{Value: 123}, "123"},
		"int64":  {param.Opt[int64]{Value: int64(123456789123456789)}, "123456789123456789"},
		"struct": {Struct{A: "yo", B: 123}, `{"a":"yo","b":123}`},
		"datetime": {
			param.Opt[time.Time]{Value: time.Date(2023, time.March, 18, 14, 47, 38, 0, time.UTC)},
			`"2023-03-18T14:47:38Z"`,
		},
		"optional_date": {
			FieldStruct{
				F: param.Opt[time.Time]{Value: time.Date(2023, time.March, 18, 14, 47, 38, 0, time.UTC)},
			},
			`{"f":"2023-03-18"}`,
		},
		"optional_time": {
			FieldStruct{
				G: param.Opt[time.Time]{Value: time.Date(2023, time.March, 18, 14, 47, 38, 0, time.UTC)},
			},
			`{"g":"2023-03-18T14:47:38Z"}`,
		},
		"optional_datetime_explicit_format": {
			FieldStruct{
				H: param.Opt[time.Time]{Value: time.Date(2023, time.March, 18, 14, 47, 38, 0, time.UTC)},
			},
			`{"h":"2023-03-18T14:47:38Z"}`,
		},
		"param_struct": {
			FieldStruct{
				A: param.Opt[string]{Value: "hello"},
				B: param.Opt[int64]{Value: int64(12)},
				D: time.Date(2023, time.March, 18, 14, 47, 38, 0, time.UTC),
				E: time.Date(2023, time.March, 18, 14, 47, 38, 0, time.UTC),
				F: param.Opt[time.Time]{Value: time.Date(2023, time.March, 18, 14, 47, 38, 0, time.UTC)},
				G: param.Opt[time.Time]{Value: time.Date(2023, time.March, 18, 14, 47, 38, 0, time.UTC)},
				H: param.Opt[time.Time]{Value: time.Date(2023, time.March, 18, 14, 47, 38, 0, time.UTC)},
			},
			`{"a":"hello","b":12,"d":"2023-03-18","e":"2023-03-18T14:47:38Z","f":"2023-03-18","g":"2023-03-18T14:47:38Z","h":"2023-03-18T14:47:38Z"}`,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			b, err := json.Marshal(test.value)
			if err != nil {
				t.Fatalf("didn't expect error %v, expected %s", err, test.expected)
			}
			if string(b) != test.expected {
				t.Fatalf("expected %s, received %s", test.expected, string(b))
			}
		})
	}
}

func TestAdditionalProperties(t *testing.T) {
	s := StructWithAdditionalProperties{
		First:  "hello",
		Second: 14,
		ExtraFields: map[string]any{
			"hi": "there",
		},
	}
	exp := `{"first":"hello","second":14,"hi":"there"}`

	bytes, err := json.Marshal(s)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	if string(bytes) != exp {
		t.Fatalf("expected %s, got %s", exp, string(bytes))
	}
}

func TestExtraFields(t *testing.T) {
	v := Struct{
		A: "hello",
		B: 123,
	}
	v.SetExtraFields(map[string]any{
		"extra": Struct{A: "recursive"},
		"b":     nil,
	})
	bytes, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}
	if string(bytes) != `{"a":"hello","b":null,"extra":{"a":"recursive","b":0}}` {
		t.Fatalf("failed to marshal: got %v", string(bytes))
	}
	if v.B != 123 {
		t.Fatalf("marshal modified field B: got %v", v.B)
	}
}

func TestExtraFieldsForceOmitted(t *testing.T) {
	v := Struct{
		// Testing with the zero value.
		// A: "",
		// B: 0,
	}
	v.SetExtraFields(map[string]any{
		"b": param.Omit,
	})
	bytes, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}
	if string(bytes) != `{"a":""}` {
		t.Fatalf("failed to marshal: got %v", string(bytes))
	}
}

type UnionWithDates struct {
	OfDate param.Opt[time.Time]
	OfTime param.Opt[time.Time]
	param.APIUnion
}

func (r UnionWithDates) MarshalJSON() (data []byte, err error) {
	return param.MarshalUnion(r, param.EncodedAsDate(r.OfDate), r.OfTime)
}

func TestUnionDateMarshal(t *testing.T) {
	tests := map[string]struct {
		value    UnionWithDates
		expected string
	}{
		"date_only": {
			UnionWithDates{
				OfDate: param.Opt[time.Time]{Value: time.Date(2023, time.March, 18, 0, 0, 0, 0, time.UTC)},
			},
			`"2023-03-18"`,
		},
		"datetime_only": {
			UnionWithDates{
				OfTime: param.Opt[time.Time]{Value: time.Date(2023, time.March, 18, 14, 47, 38, 0, time.UTC)},
			},
			`"2023-03-18T14:47:38Z"`,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			b, err := json.Marshal(test.value)
			if err != nil {
				t.Fatalf("didn't expect error %v, expected %s", err, test.expected)
			}
			if string(b) != test.expected {
				t.Fatalf("expected %s, received %s", test.expected, string(b))
			}
		})
	}
}

func TestOverride(t *testing.T) {
	tests := map[string]struct {
		value    param.ParamStruct
		expected string
	}{
		"param_struct": {
			param.Override[FieldStruct](map[string]any{
				"a": "hello",
				"b": 12,
				"c": nil,
			}),
			`{"a":"hello","b":12,"c":null}`,
		},
		"param_struct_primitive": {
			param.Override[FieldStruct](12),
			`12`,
		},
		"param_struct_null": {
			param.Override[FieldStruct](nil),
			`null`,
		},
	}

	f := FieldStruct{}

	f.SetExtraFields(map[string]any{
		"z": "ok",
	})

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			b, err := json.Marshal(test.value)
			if err != nil {
				t.Fatalf("didn't expect error %v, expected %s", err, test.expected)
			}
			if string(b) != test.expected {
				t.Fatalf("expected %s, received %s", test.expected, string(b))
			}
			if _, ok := test.value.Overrides(); !ok {
				t.Fatalf("expected to be overridden")
			}
		})
	}
}

// Despite implementing the interface, this struct is not an param.Optional
// since it was defined in a different package.
type almostOpt struct{}

func (almostOpt) Valid() bool  { return true }
func (almostOpt) Null() bool   { return false }
func (almostOpt) isZero() bool { return false }

func (almostOpt) implOpt() {}

func TestOptionalInterfaceAssignability(t *testing.T) {
	optInt := param.Opt[int]{}
	if _, ok := any(optInt).(param.Optional); !ok {
		t.Fatalf("failed to assign")
	}

	notOpt := almostOpt{}
	if _, ok := any(notOpt).(param.Optional); ok {
		t.Fatalf("unexpected successful assignment")
	}

	notOpt.implOpt() // silence the warning
}

type PrimitiveUnion struct {
	OfString param.Opt[string]
	OfInt    param.Opt[int]
	param.APIUnion
}

func (p PrimitiveUnion) MarshalJSON() (data []byte, err error) {
	return param.MarshalUnion(p, p.OfString, p.OfInt)
}

func TestOverriddenUnion(t *testing.T) {
	tests := map[string]struct {
		value    PrimitiveUnion
		expected string
	}{
		"string": {
			param.Override[PrimitiveUnion](json.RawMessage(`"hello"`)),
			`"hello"`,
		},
		"int": {
			param.Override[PrimitiveUnion](json.RawMessage(`42`)),
			`42`,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			b, err := json.Marshal(test.value)
			if err != nil {
				t.Fatalf("didn't expect error %v, expected %s", err, test.expected)
			}
			if string(b) != test.expected {
				t.Fatalf("expected %s, received %s", test.expected, string(b))
			}
		})
	}
}

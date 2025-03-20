package param_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/openai/openai-go/packages/param"
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

type FieldStruct struct {
	A param.Opt[string]    `json:"a,omitzero"`
	B param.Opt[int64]     `json:"b,omitzero"`
	C Struct               `json:"c,omitzero"`
	D time.Time            `json:"d,omitzero" format:"date"`
	E time.Time            `json:"e,omitzero"`
	F param.Opt[time.Time] `json:"f,omitzero" format:"date"`
	G param.Opt[time.Time] `json:"g,omitzero"`
	param.APIObject
}

func (r FieldStruct) MarshalJSON() (data []byte, err error) {
	type shadow FieldStruct
	return param.MarshalObject(r, (*shadow)(&r))
}

func TestFieldMarshal(t *testing.T) {
	tests := map[string]struct {
		value    interface{}
		expected string
	}{
		"null_string": {param.NullOpt[string](), "null"},
		"null_int64":  {param.NullOpt[int64](), "null"},
		"null_time":   {param.NullOpt[time.Time](), "null"},
		"null_struct": {param.NullObj[Struct](), "null"},

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

		"param_struct": {
			FieldStruct{
				A: param.Opt[string]{Value: "hello"},
				B: param.Opt[int64]{Value: int64(12)},
				D: time.Date(2023, time.March, 18, 14, 47, 38, 0, time.UTC),
				E: time.Date(2023, time.March, 18, 14, 47, 38, 0, time.UTC),
				F: param.Opt[time.Time]{Value: time.Date(2023, time.March, 18, 14, 47, 38, 0, time.UTC)},
				G: param.Opt[time.Time]{Value: time.Date(2023, time.March, 18, 14, 47, 38, 0, time.UTC)},
			},
			`{"a":"hello","b":12,"d":"2023-03-18","e":"2023-03-18T14:47:38Z","f":"2023-03-18","g":"2023-03-18T14:47:38Z"}`,
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

func TestExtraFields(t *testing.T) {
	v := Struct{
		A: "hello",
		B: 123,
	}
	v.WithExtraFields(map[string]any{
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

type UnionWithDates struct {
	OfDate param.Opt[time.Time]
	OfTime param.Opt[time.Time]
}

func (r UnionWithDates) MarshalJSON() (data []byte, err error) {
	return param.MarshalUnion[UnionWithDates](param.EncodedAsDate(r.OfDate), r.OfTime)
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

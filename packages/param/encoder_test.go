package param_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/openai/openai-go/packages/param"
)

type Struct struct {
	A string `json:"a,omitzero"`
	B int64  `json:"b,omitzero"`
	param.APIObject
}

func (r Struct) MarshalJSON() (data []byte, err error) {
	type shadow Struct
	return param.MarshalObject(r, (*shadow)(&r))
}

type FieldStruct struct {
	A param.String   `json:"a,omitzero"`
	B param.Int      `json:"b,omitzero"`
	C Struct         `json:"c,omitzero"`
	D param.Date     `json:"d,omitzero"`
	E param.Datetime `json:"e,omitzero"`
	F param.Int      `json:"f,omitzero"`
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
		"null_string": {param.Null[param.String](), "null"},
		"null_int":    {param.Null[param.Int](), "null"},
		"null_int64":  {param.Null[param.Int](), "null"},
		"null_struct": {param.Null[Struct](), "null"},

		"string": {param.String{V: "string"}, `"string"`},
		"int":    {param.Int{V: 123}, "123"},
		"int64":  {param.Int{V: int64(123456789123456789)}, "123456789123456789"},
		"struct": {Struct{A: "yo", B: 123}, `{"a":"yo","b":123}`},
		"date":   {param.Date{V: time.Date(2023, time.March, 18, 14, 47, 38, 0, time.UTC)}, `"2023-03-18"`},
		"datetime": {
			param.Datetime{V: time.Date(2023, time.March, 18, 14, 47, 38, 0, time.UTC)},
			`"2023-03-18T14:47:38Z"`,
		},

		"string_raw": {param.Override[param.Int]("string"), `"string"`},
		"int_raw":    {param.Override[param.Int](123), "123"},
		"int64_raw":  {param.Override[param.Int](int64(123456789123456789)), "123456789123456789"},
		"struct_raw": {param.Override[param.Int](Struct{A: "yo", B: 123}), `{"a":"yo","b":123}`},

		"param_struct": {
			FieldStruct{
				A: param.String{V: "hello"},
				B: param.Int{V: int64(12)},
				D: param.Date{V: time.Date(2023, time.March, 18, 14, 47, 38, 0, time.UTC)},
				E: param.Datetime{V: time.Date(2023, time.March, 18, 14, 47, 38, 0, time.UTC)},
			},
			`{"a":"hello","b":12,"d":"2023-03-18","e":"2023-03-18T14:47:38Z"}`,
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

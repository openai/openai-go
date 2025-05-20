package apijson

import (
	"reflect"
	"testing"
)

type EnumStruct struct {
	NormalString string        `json:"normal_string"`
	StringEnum   string        `json:"string_enum"`
	NamedEnum    NamedEnumType `json:"named_enum"`

	IntEnum  int  `json:"int_enum"`
	BoolEnum bool `json:"bool_enum"`

	WeirdBoolEnum bool `json:"weird_bool_enum"`
}

func (o *EnumStruct) UnmarshalJSON(data []byte) error {
	return UnmarshalRoot(data, o)
}

func init() {
	RegisterFieldValidator[EnumStruct]("string_enum", "one", "two", "three")
	RegisterFieldValidator[EnumStruct]("int_enum", 200, 404)
	RegisterFieldValidator[EnumStruct]("bool_enum", false)
	RegisterFieldValidator[EnumStruct]("weird_bool_enum", true, false)
}

type NamedEnumType string

const (
	NamedEnumOne   NamedEnumType = "one"
	NamedEnumTwo   NamedEnumType = "two"
	NamedEnumThree NamedEnumType = "three"
)

func (e NamedEnumType) IsKnown() bool {
	return e == NamedEnumOne || e == NamedEnumTwo || e == NamedEnumThree
}

func TestEnumStructStringValidator(t *testing.T) {
	cases := map[string]struct {
		exactness
		EnumStruct
	}{
		`{"string_enum":"one"}`:     {exact, EnumStruct{StringEnum: "one"}},
		`{"string_enum":"two"}`:     {exact, EnumStruct{StringEnum: "two"}},
		`{"string_enum":"three"}`:   {exact, EnumStruct{StringEnum: "three"}},
		`{"string_enum":"none"}`:    {loose, EnumStruct{StringEnum: "none"}},
		`{"int_enum":200}`:          {exact, EnumStruct{IntEnum: 200}},
		`{"int_enum":404}`:          {exact, EnumStruct{IntEnum: 404}},
		`{"int_enum":500}`:          {loose, EnumStruct{IntEnum: 500}},
		`{"bool_enum":false}`:       {exact, EnumStruct{BoolEnum: false}},
		`{"bool_enum":true}`:        {loose, EnumStruct{BoolEnum: true}},
		`{"weird_bool_enum":true}`:  {exact, EnumStruct{WeirdBoolEnum: true}},
		`{"weird_bool_enum":false}`: {exact, EnumStruct{WeirdBoolEnum: false}},

		`{"named_enum":"one"}`:  {exact, EnumStruct{NamedEnum: NamedEnumOne}},
		`{"named_enum":"none"}`: {loose, EnumStruct{NamedEnum: "none"}},

		`{"string_enum":"one","named_enum":"one"}`: {exact, EnumStruct{NamedEnum: "one", StringEnum: "one"}},
		`{"string_enum":"four","named_enum":"one"}`: {
			loose,
			EnumStruct{NamedEnum: "one", StringEnum: "four"},
		},
		`{"string_enum":"one","named_enum":"four"}`: {
			loose, EnumStruct{NamedEnum: "four", StringEnum: "one"},
		},
		`{"wrong_key":"one"}`: {extras, EnumStruct{StringEnum: ""}},
	}

	for raw, expected := range cases {
		var dst EnumStruct

		dec := decoderBuilder{root: true}
		exactness, _ := dec.unmarshalWithExactness([]byte(raw), &dst)

		if !reflect.DeepEqual(dst, expected.EnumStruct) {
			t.Fatalf("failed equality check %#v", dst)
		}

		if exactness != expected.exactness {
			t.Fatalf("exactness got %d expected %d %s", exactness, expected.exactness, raw)
		}
	}
}

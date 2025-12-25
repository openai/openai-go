package apijson_test

import (
	"encoding/json"
	"fmt"
	"github.com/Nordlys-Labs/openai-go/v3/internal/apijson"
	"github.com/Nordlys-Labs/openai-go/v3/packages/param"
	"reflect"
	"testing"
)

func TestOptionalDecoders(t *testing.T) {
	cases := map[string]struct {
		buf string
		val any
	}{

		"opt_string_present": {
			`"hello"`,
			param.NewOpt("hello"),
		},
		"opt_string_empty_present": {
			`""`,
			param.NewOpt(""),
		},
		"opt_string_null": {
			`null`,
			param.Null[string](),
		},
		"opt_string_null_with_whitespace": {
			`  null  `,
			param.Null[string](),
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			result := reflect.New(reflect.TypeOf(test.val))
			if err := json.Unmarshal([]byte(test.buf), result.Interface()); err != nil {
				t.Fatalf("deserialization of %v failed with error %v", result, err)
			}

			if !reflect.DeepEqual(result.Elem().Interface(), test.val) {
				t.Fatalf("expected '%s' to deserialize to \n%#v\nbut got\n%#v", test.buf, test.val, result.Elem().Interface())
			}
		})
	}
}

type paramObject = param.APIObject

type BasicObject struct {
	ReqInt    int64   `json:"req_int,required"`
	ReqFloat  float64 `json:"req_float,required"`
	ReqString string  `json:"req_string,required"`
	ReqBool   bool    `json:"req_bool,required"`

	OptInt    param.Opt[int64]   `json:"opt_int"`
	OptFloat  param.Opt[float64] `json:"opt_float"`
	OptString param.Opt[string]  `json:"opt_string"`
	OptBool   param.Opt[bool]    `json:"opt_bool"`

	paramObject
}

func (o *BasicObject) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, o) }

func TestBasicObjectWithNull(t *testing.T) {
	raw := `{"opt_int":null,"opt_string":null,"opt_bool":null}`
	var dst BasicObject
	target := BasicObject{
		OptInt: param.Null[int64](),
		// OptFloat:  param.Opt[float64]{},
		OptString: param.Null[string](),
		OptBool:   param.Null[bool](),
	}

	err := json.Unmarshal([]byte(raw), &dst)
	if err != nil {
		t.Fatalf("failed unmarshal")
	}

	if !reflect.DeepEqual(dst, target) {
		t.Fatalf("failed equality check %#v", dst)
	}
}

func TestBasicObject(t *testing.T) {
	raw := `{"req_int":1,"req_float":1.3,"req_string":"test","req_bool":true,"opt_int":2,"opt_float":2.0,"opt_string":"test","opt_bool":false}`
	var dst BasicObject
	target := BasicObject{
		ReqInt:    1,
		ReqFloat:  1.3,
		ReqString: "test",
		ReqBool:   true,
		OptInt:    param.NewOpt[int64](2),
		OptFloat:  param.NewOpt(2.0),
		OptString: param.NewOpt("test"),
		OptBool:   param.NewOpt(false),
	}

	err := json.Unmarshal([]byte(raw), &dst)
	if err != nil {
		t.Fatalf("failed unmarshal")
	}

	if !reflect.DeepEqual(dst, target) {
		t.Fatalf("failed equality check %#v", dst)
	}
}

type ComplexObject struct {
	Basic BasicObject `json:"basic,required"`
	Enum  string      `json:"enum"`
	paramObject
}

func (o *ComplexObject) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, o) }

func init() {
	apijson.RegisterFieldValidator[ComplexObject]("enum", "a", "b", "c")
}

func TestComplexObject(t *testing.T) {
	raw := `{"basic":{"req_int":1,"req_float":1.3,"req_string":"test","req_bool":true,"opt_int":2,"opt_float":2.0,"opt_string":"test","opt_bool":false},"enum":"a"}`
	var dst ComplexObject

	target := ComplexObject{
		Basic: BasicObject{
			ReqInt:    1,
			ReqFloat:  1.3,
			ReqString: "test",
			ReqBool:   true,
			OptInt:    param.NewOpt[int64](2),
			OptFloat:  param.NewOpt(2.0),
			OptString: param.NewOpt("test"),
			OptBool:   param.NewOpt(false),
		},
		Enum: "a",
	}

	err := json.Unmarshal([]byte(raw), &dst)
	if err != nil {
		t.Fatalf("failed unmarshal")
	}

	if !reflect.DeepEqual(dst, target) {
		t.Fatalf("failed equality check %#v", dst)
	}
}

type paramUnion = param.APIUnion

type MemberA struct {
	Name string `json:"name,required"`
	Age  int    `json:"age,required"`
}

type MemberB struct {
	Name string `json:"name,required"`
	Age  string `json:"age,required"`
}

type MemberC struct {
	Name   string `json:"name,required"`
	Age    int    `json:"age,required"`
	Status string `json:"status"`
}

type MemberD struct {
	Cost   int    `json:"cost,required"`
	Status string `json:"status,required"`
}

type MemberE struct {
	Cost   int    `json:"cost,required"`
	Status string `json:"status,required"`
}

type MemberF struct {
	D int            `json:"d"`
	E string         `json:"e"`
	F float64        `json:"f"`
	G param.Opt[int] `json:"g"`
}

type MemberG struct {
	D int             `json:"d"`
	E string          `json:"e"`
	F float64         `json:"f"`
	G param.Opt[bool] `json:"g"`
}

func init() {
	apijson.RegisterFieldValidator[MemberD]("status", "good", "ok", "bad")
	apijson.RegisterFieldValidator[MemberE]("status", "GOOD", "OK", "BAD")
}

type UnionStruct struct {
	OfMemberA *MemberA          `json:",inline"`
	OfMemberB *MemberB          `json:",inline"`
	OfMemberC *MemberC          `json:",inline"`
	OfMemberD *MemberD          `json:",inline"`
	OfMemberE *MemberE          `json:",inline"`
	OfMemberF *MemberF          `json:",inline"`
	OfMemberG *MemberG          `json:",inline"`
	OfString  param.Opt[string] `json:",inline"`

	paramUnion
}

func (union *UnionStruct) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, union)
}

func TestUnionStruct(t *testing.T) {
	tests := map[string]struct {
		raw        string
		target     UnionStruct
		shouldFail bool
	}{
		"fail": {
			raw:        `1200`,
			target:     UnionStruct{},
			shouldFail: true,
		},
		"easy": {
			raw:    `{"age":30}`,
			target: UnionStruct{OfMemberA: &MemberA{Age: 30}},
		},
		"less-easy": {
			raw:    `{"age":"thirty"}`,
			target: UnionStruct{OfMemberB: &MemberB{Age: "thirty"}},
		},
		"even-less-easy": {
			raw:    `{"age":"30"}`,
			target: UnionStruct{OfMemberB: &MemberB{Age: "30"}},
		},
		"medium": {
			raw: `{"name":"jacob","age":30}`,
			target: UnionStruct{OfMemberA: &MemberA{
				Age:  30,
				Name: "jacob",
			}},
		},
		"less-medium": {
			raw: `{"name":"jacob","age":"thirty"}`,
			target: UnionStruct{OfMemberB: &MemberB{
				Age:  "thirty",
				Name: "jacob",
			}},
		},
		"even-less-medium": {
			raw: `{"name":"jacob","age":"30"}`,
			target: UnionStruct{OfMemberB: &MemberB{
				Name: "jacob",
				Age:  "30",
			}},
		},
		"hard": {
			raw: `{"name":"jacob","age":30,"status":"active"}`,
			target: UnionStruct{OfMemberC: &MemberC{
				Name:   "jacob",
				Age:    30,
				Status: "active",
			}},
		},
		"inline-string": {
			raw:    `"hello there"`,
			target: UnionStruct{OfString: param.NewOpt("hello there")},
		},
		"enum-field": {
			raw:    `{"cost":100,"status":"ok"}`,
			target: UnionStruct{OfMemberD: &MemberD{Cost: 100, Status: "ok"}},
		},
		"other-enum-field": {
			raw:    `{"cost":100,"status":"GOOD"}`,
			target: UnionStruct{OfMemberE: &MemberE{Cost: 100, Status: "GOOD"}},
		},
		"tricky-extra-fields": {
			raw:    `{"d":12,"e":"hello","f":1.00}`,
			target: UnionStruct{OfMemberF: &MemberF{D: 12, E: "hello", F: 1.00}},
		},
		"optional-fields": {
			raw:    `{"d":12,"e":"hello","f":1.00,"g":12}`,
			target: UnionStruct{OfMemberF: &MemberF{D: 12, E: "hello", F: 1.00, G: param.NewOpt(12)}},
		},
		"optional-fields-2": {
			raw:    `{"d":12,"e":"hello","f":1.00,"g":false}`,
			target: UnionStruct{OfMemberG: &MemberG{D: 12, E: "hello", F: 1.00, G: param.NewOpt(false)}},
		},
	}

	for name, test := range tests {
		var dst UnionStruct
		t.Run(name, func(t *testing.T) {
			err := json.Unmarshal([]byte(test.raw), &dst)
			if err != nil && !test.shouldFail {
				t.Fatalf("failed unmarshal with err: %v %#v", err, dst)
			}

			if !reflect.DeepEqual(dst, test.target) {
				if dst.OfMemberA != nil {
					fmt.Printf("%#v", dst.OfMemberA)
				}
				t.Fatalf("failed equality, got %#v but expected %#v", dst, test.target)
			}
		})
	}
}

type ConstantA string
type ConstantB string
type ConstantC string

func (c ConstantA) Default() string { return "A" }
func (c ConstantB) Default() string { return "B" }
func (c ConstantC) Default() string { return "C" }

type DiscVariantA struct {
	Name string    `json:"name,required"`
	Age  int       `json:"age,required"`
	Type ConstantA `json:"type,required"`
}

type DiscVariantB struct {
	Name string    `json:"name,required"`
	Age  int       `json:"age,required"`
	Type ConstantB `json:"type,required"`
}

type DiscVariantC struct {
	Name string    `json:"name,required"`
	Age  float64   `json:"age,required"`
	Type ConstantC `json:"type,required"`
}

type DiscriminatedUnion struct {
	OfA *DiscVariantA `json:",inline"`
	OfB *DiscVariantB `json:",inline"`
	OfC *DiscVariantC `json:",inline"`

	paramUnion
}

func init() {
	apijson.RegisterDiscriminatedUnion[DiscriminatedUnion]("type", map[string]reflect.Type{
		"A": reflect.TypeOf(DiscVariantA{}),
		"B": reflect.TypeOf(DiscVariantB{}),
		"C": reflect.TypeOf(DiscVariantC{}),
	})
}

type FooVariant struct {
	Type  string `json:"type,required"`
	Value string `json:"value,required"`
}

type BarVariant struct {
	Type   string `json:"type,required"`
	Enable bool   `json:"enable,required"`
}

type MultiDiscriminatorUnion struct {
	OfFoo *FooVariant `json:",inline"`
	OfBar *BarVariant `json:",inline"`

	paramUnion
}

func init() {
	apijson.RegisterDiscriminatedUnion[MultiDiscriminatorUnion]("type", map[string]reflect.Type{
		"foo":        reflect.TypeOf(FooVariant{}),
		"foo_v2":     reflect.TypeOf(FooVariant{}),
		"bar":        reflect.TypeOf(BarVariant{}),
		"bar_legacy": reflect.TypeOf(BarVariant{}),
	})
}

func (m *MultiDiscriminatorUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, m)
}

func (d *DiscriminatedUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, d)
}

func TestDiscriminatedUnion(t *testing.T) {
	tests := map[string]struct {
		raw        string
		target     DiscriminatedUnion
		shouldFail bool
	}{
		"variant_a": {
			raw: `{"name":"Alice","age":25,"type":"A"}`,
			target: DiscriminatedUnion{OfA: &DiscVariantA{
				Name: "Alice",
				Age:  25,
				Type: "A",
			}},
		},
		"variant_b": {
			raw: `{"name":"Bob","age":30,"type":"B"}`,
			target: DiscriminatedUnion{OfB: &DiscVariantB{
				Name: "Bob",
				Age:  30,
				Type: "B",
			}},
		},
		"variant_c": {
			raw: `{"name":"Charlie","age":35.5,"type":"C"}`,
			target: DiscriminatedUnion{OfC: &DiscVariantC{
				Name: "Charlie",
				Age:  35.5,
				Type: "C",
			}},
		},
		"invalid_type": {
			raw:        `{"name":"Unknown","age":40,"type":"D"}`,
			target:     DiscriminatedUnion{},
			shouldFail: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var dst DiscriminatedUnion
			err := json.Unmarshal([]byte(test.raw), &dst)
			if err != nil && !test.shouldFail {
				t.Fatalf("failed unmarshal with err: %v", err)
			}
			if err == nil && test.shouldFail {
				t.Fatalf("expected unmarshal to fail but it succeeded")
			}
			if !reflect.DeepEqual(dst, test.target) {
				t.Fatalf("failed equality, got %#v but expected %#v", dst, test.target)
			}
		})
	}
}

func TestMultiDiscriminatorUnion(t *testing.T) {
	tests := map[string]struct {
		raw        string
		target     MultiDiscriminatorUnion
		shouldFail bool
	}{
		"foo_variant": {
			raw: `{"type":"foo","value":"test"}`,
			target: MultiDiscriminatorUnion{OfFoo: &FooVariant{
				Type:  "foo",
				Value: "test",
			}},
		},
		"foo_v2_variant": {
			raw: `{"type":"foo_v2","value":"test_v2"}`,
			target: MultiDiscriminatorUnion{OfFoo: &FooVariant{
				Type:  "foo_v2",
				Value: "test_v2",
			}},
		},
		"bar_variant": {
			raw: `{"type":"bar","enable":true}`,
			target: MultiDiscriminatorUnion{OfBar: &BarVariant{
				Type:   "bar",
				Enable: true,
			}},
		},
		"bar_legacy_variant": {
			raw: `{"type":"bar_legacy","enable":false}`,
			target: MultiDiscriminatorUnion{OfBar: &BarVariant{
				Type:   "bar_legacy",
				Enable: false,
			}},
		},
		"invalid_type": {
			raw:        `{"type":"unknown","value":"test"}`,
			target:     MultiDiscriminatorUnion{},
			shouldFail: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var dst MultiDiscriminatorUnion
			err := json.Unmarshal([]byte(test.raw), &dst)
			if err != nil && !test.shouldFail {
				t.Fatalf("failed unmarshal with err: %v", err)
			}
			if err == nil && test.shouldFail {
				t.Fatalf("expected unmarshal to fail but it succeeded")
			}
			if !reflect.DeepEqual(dst, test.target) {
				t.Fatalf("failed equality, got %#v but expected %#v", dst, test.target)
			}
		})
	}
}

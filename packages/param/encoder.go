package param

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/openai/openai-go/internal/apifield"
	shimjson "github.com/openai/openai-go/internal/encoding/json"
)

// This uses a shimmed 'encoding/json' from Go 1.24, to support the 'omitzero' tag
func MarshalObject[T ObjectFielder](f T, underlying any) ([]byte, error) {
	if f.IsNull() {
		return []byte("null"), nil
	} else if ovr, ok := f.IsOverridden(); ok {
		// TODO(v2): handle if ovr.(ExtraFields)
		return shimjson.Marshal(ovr)
	} else {
		return shimjson.Marshal(underlying)
	}
}

// This uses a shimmed 'encoding/json' from Go 1.24, to support the 'omitzero' tag
func MarshalUnion[T any](variants ...any) ([]byte, error) {
	nPresent := 0
	idx := -1
	for i, variant := range variants {
		if !IsOmitted(variant) {
			nPresent++
			idx = i
		}
	}
	if nPresent == 0 || idx == -1 {
		return []byte(`{}`), nil
	} else if nPresent > 1 {
		return nil, &json.MarshalerError{
			Type: reflect.TypeOf((*T)(nil)).Elem(),
			Err:  fmt.Errorf("expected union to have one present variant, got %d", nPresent),
		}
	}
	return shimjson.Marshal(variants[idx])
}

// This uses a shimmed stdlib 'encoding/json' from Go 1.24, to support omitzero
func marshalField[T interface {
	Fielder
	IsOmitted() bool
}](f T, happyPath any) ([]byte, error) {
	if f.IsNull() {
		return []byte("null"), nil
	} else if ovr, ok := f.IsOverridden(); ok {
		return shimjson.Marshal(ovr)
	} else {
		return shimjson.Marshal(happyPath)
	}
}

// This uses a shimmed 'encoding/json' from Go 1.24, to support the 'omitzero' tag
func unmarshalField[T any](underlying *T, meta *metadata, data []byte) error {
	if string(data) == `null` {
		meta.setMetadata(apifield.ExplicitNull{})
		return nil
	}
	if err := shimjson.Unmarshal(data, &underlying); err != nil {
		meta.setMetadata(apifield.ResponseData(data))
		return err
	}
	meta.setMetadata(apifield.NeverOmitted{})
	return nil
}

func stringifyField[T Fielder](f T, fallback any) string {
	if f.IsNull() {
		return "null"
	}
	if v, ok := f.IsOverridden(); ok {
		return fmt.Sprintf("%v", v)
	}
	return fmt.Sprintf("%v", fallback)
}

// shimmed from Go 1.23 "reflect" package
func TypeFor[T any]() reflect.Type {
	var v T
	if t := reflect.TypeOf(v); t != nil {
		return t // optimize for T being a non-interface kind
	}
	return reflect.TypeOf((*T)(nil)).Elem() // only for an interface kind
}

var richStringType = TypeFor[String]()
var richIntType = TypeFor[Int]()
var richFloatType = TypeFor[Float]()
var richBoolType = TypeFor[Bool]()
var richDateType = TypeFor[Date]()
var richDatetimeType = TypeFor[Datetime]()

// indexOfUnderlyingValueField must only be called at initialization time
func indexOfUnderlyingValueField(t reflect.Type) []int {
	field, ok := t.FieldByName("V")
	if !ok {
		panic("unreachable: initialization issue, underlying value field not found")
	}
	return field.Index
}

var RichPrimitiveTypes = map[reflect.Type][]int{
	richStringType:   indexOfUnderlyingValueField(richStringType),
	richIntType:      indexOfUnderlyingValueField(richIntType),
	richFloatType:    indexOfUnderlyingValueField(richFloatType),
	richBoolType:     indexOfUnderlyingValueField(richBoolType),
	richDateType:     indexOfUnderlyingValueField(richDateType),
	richDatetimeType: indexOfUnderlyingValueField(richDatetimeType),
}

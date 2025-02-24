package apiquery

import (
	"fmt"
	"github.com/openai/openai-go/packages/param"
	"reflect"
)

// TODO(v2): verify this is correct w.r.t. to null, override and omit handling
func (e *encoder) newRichFieldTypeEncoder(t reflect.Type, underlyingValueIdx []int) encoderFunc {
	underlying := t.FieldByIndex(underlyingValueIdx)
	primitiveEncoder := e.newPrimitiveTypeEncoder(underlying.Type)
	return func(key string, value reflect.Value) []Pair {
		if fielder, ok := value.Interface().(param.Fielder); ok {
			if fielder.IsNull() {
				return []Pair{{key, "null"}}
			} else if ovr, ok := fielder.IsOverridden(); ok {
				ovr := reflect.ValueOf(ovr)
				encode := e.newTypeEncoder(ovr.Type())
				return encode(key, ovr)
			} else if !param.IsOmitted(fielder) {
				res := primitiveEncoder(key, value.FieldByName("V"))
				fmt.Printf("%#v\n", res)
				return res
			}
		}
		return nil
	}
}

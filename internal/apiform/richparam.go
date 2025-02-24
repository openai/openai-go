package apiform

import (
	"github.com/openai/openai-go/packages/param"
	"mime/multipart"
	"reflect"
)

// TODO(v2): verify this is correct, w.r.t. to null, overrides and omit
func (e *encoder) newRichFieldTypeEncoder(t reflect.Type, underlyingValueIdx []int) encoderFunc {
	underlying := t.FieldByIndex(underlyingValueIdx)
	primitiveEncoder := e.newPrimitiveTypeEncoder(underlying.Type)
	return func(key string, value reflect.Value, writer *multipart.Writer) error {
		if fielder, ok := value.Interface().(param.Fielder); ok {
			if fielder.IsNull() {
				return writer.WriteField(key, "null")
			} else if ovr, ok := fielder.IsOverridden(); ok {
				ovr := reflect.ValueOf(ovr)
				encode := e.newTypeEncoder(ovr.Type())
				return encode(key, ovr, writer)
			} else if !param.IsOmitted(fielder) {
				return primitiveEncoder(key, value.FieldByName("V"), writer)
			}
		}
		return nil
	}
}

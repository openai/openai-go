package apijson

import (
	"github.com/openai/openai-go/packages/resp"
	"reflect"
)

func getSubField(root reflect.Value, index []int, name string) reflect.Value {
	strct := root.FieldByIndex(index[:len(index)-1])
	if !strct.IsValid() {
		panic("couldn't find encapsulating struct for field " + name)
	}
	meta := strct.FieldByName("JSON")
	if !meta.IsValid() {
		return reflect.Value{}
	}
	field := meta.FieldByName(name)
	if !field.IsValid() {
		return reflect.Value{}
	}
	return field
}

var respFieldType = reflect.TypeOf(resp.Field{})
var fieldType = reflect.TypeOf(Field{})

func setSubField(root reflect.Value, index []int, name string, meta Field) {
	if metadata := getSubField(root, index, name); metadata.IsValid() {
		if metadata.Type() == respFieldType {
			var rf resp.Field
			if meta.IsNull() {
				rf = resp.NewNullField()
			} else if meta.IsMissing() {
				_ = rf
			} else if meta.IsInvalid() {
				rf = resp.NewInvalidField(meta.raw)
			} else {
				rf = resp.NewValidField(meta.raw)
			}
			metadata.Set(reflect.ValueOf(rf))
		} else if metadata.Type() == fieldType {
			metadata.Set(reflect.ValueOf(meta))
		}
	}
}

package resp

import (
	"fmt"
	"reflect"
)

func VariantFromUnion(u reflect.Value) (any, error) {
	if u.Kind() == reflect.Ptr {
		u = u.Elem()
	}

	if u.Kind() != reflect.Struct {
		return nil, fmt.Errorf("resp: cannot extract variant from non-struct union")
	}

	nVariants := 0
	variantIdx := -1
	for i := 0; i < u.NumField(); i++ {
		if !u.Field(i).IsZero() {
			nVariants++
			variantIdx = i
		}
	}

	if nVariants > 1 {
		return nil, fmt.Errorf("resp: cannot extract variant from union with multiple variants")
	}

	if nVariants == 0 {
		return nil, fmt.Errorf("resp: cannot extract variant from union with no variants")
	}

	return u.Field(variantIdx).Interface(), nil
}

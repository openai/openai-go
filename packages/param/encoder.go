package param

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	shimjson "github.com/openai/openai-go/internal/encoding/json"

	"github.com/tidwall/sjson"
)

// This type will not be stable and shouldn't be relied upon
type EncodedAsDate Opt[time.Time]

func (m EncodedAsDate) MarshalJSON() ([]byte, error) {
	underlying := Opt[time.Time](m)
	bytes := underlying.MarshalJSONWithTimeLayout("2006-01-02")
	if len(bytes) > 0 {
		return bytes, nil
	}
	return underlying.MarshalJSON()
}

// This uses a shimmed 'encoding/json' from Go 1.24, to support the 'omitzero' tag
func MarshalObject[T OverridableObject](f T, underlying any) ([]byte, error) {
	if f.IsNull() {
		return []byte("null"), nil
	} else if extras := f.GetExtraFields(); extras != nil {
		bytes, err := shimjson.Marshal(underlying)
		if err != nil {
			return nil, err
		}
		for k, v := range extras {
			if v == Omit {
				// Errors handling ForceOmitted are ignored.
				if b, e := sjson.DeleteBytes(bytes, k); e == nil {
					bytes = b
				}
			} else {
				bytes, err = sjson.SetBytes(bytes, k, v)
			}
			if err != nil {
				return nil, err
			}
		}
		return bytes, nil
	} else if ovr, ok := f.IsOverridden(); ok {
		return shimjson.Marshal(ovr)
	} else {
		return shimjson.Marshal(underlying)
	}
}

// This uses a shimmed 'encoding/json' from Go 1.24, to support the 'omitzero' tag
func MarshalUnion[T any](variants ...any) ([]byte, error) {
	nPresent := 0
	presentIdx := -1
	for i, variant := range variants {
		if !IsOmitted(variant) {
			nPresent++
			presentIdx = i
		}
	}
	if nPresent == 0 || presentIdx == -1 {
		return []byte(`null`), nil
	} else if nPresent > 1 {
		return nil, &json.MarshalerError{
			Type: typeFor[T](),
			Err:  fmt.Errorf("expected union to have only one present variant, got %d", nPresent),
		}
	}
	return shimjson.Marshal(variants[presentIdx])
}

// shimmed from Go 1.23 "reflect" package
func typeFor[T any]() reflect.Type {
	var v T
	if t := reflect.TypeOf(v); t != nil {
		return t // optimize for T being a non-interface kind
	}
	return reflect.TypeOf((*T)(nil)).Elem() // only for an interface kind
}

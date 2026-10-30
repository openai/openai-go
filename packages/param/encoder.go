package param

import (
	"encoding/json"
	"fmt"
	"maps"
	"reflect"
	"slices"
	"strings"
	"time"

	shimjson "github.com/openai/openai-go/v3/internal/encoding/json"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// EncodedAsDate is not be stable and shouldn't be relied upon
type EncodedAsDate Opt[time.Time]

// If we want to set a literal key value into JSON using sjson, we need to make sure it doesn't have
// special characters that sjson interprets as a path.
var EscapeSJSONKey = strings.NewReplacer("\\", "\\\\", "|", "\\|", "#", "\\#", "@", "\\@", "*", "\\*", ".", "\\.", ":", "\\:", "?", "\\?").Replace

type forceOmit int

func (m EncodedAsDate) MarshalJSON() ([]byte, error) {
	underlying := Opt[time.Time](m)
	bytes := underlying.MarshalJSONWithTimeLayout("2006-01-02")
	if len(bytes) > 0 {
		return bytes, nil
	}
	return underlying.MarshalJSON()
}

// MarshalObject uses a shimmed 'encoding/json' from Go 1.24, to support the 'omitzero' tag
//
// Stability for the API of MarshalObject is not guaranteed.
func MarshalObject[T ParamStruct](f T, underlying any) ([]byte, error) {
	return MarshalWithExtras(f, underlying, f.extraFields())
}

// MarshalWithExtras is used to marshal a struct with additional properties.
//
// Stability for the API of MarshalWithExtras is not guaranteed.
func MarshalWithExtras[T ParamStruct, R any](f T, underlying any, extras map[string]R) ([]byte, error) {
	if f.null() {
		return []byte("null"), nil
	} else if len(extras) > 0 {
		bytes, err := shimjson.Marshal(underlying)
		if err != nil {
			return nil, err
		}
		keys := slices.Sorted(maps.Keys(extras))
		object := gjson.ParseBytes(bytes)
		legacy := !object.IsObject()
		for _, k := range keys {
			// Omit has always accepted an unescaped SJSON path. Preserve its
			// effects on earlier extras as well as on the underlying object.
			if any(extras[k]) == Omit && (k == "" || EscapeSJSONKey(k) != k) {
				legacy = true
				break
			}
		}
		if legacy {
			for _, k := range keys {
				v := extras[k]
				if any(v) == Omit {
					// Errors when handling ForceOmitted are ignored.
					if b, e := sjson.DeleteBytes(bytes, k); e == nil {
						bytes = b
					}
					continue
				}
				bytes, err = sjson.SetBytes(bytes, EscapeSJSONKey(k), v)
				if err != nil {
					return nil, err
				}
			}
			return bytes, nil
		}

		// Encode in the established sorted order before merging, so custom
		// marshalers and errors still occur in that order. SJSON encodes []byte
		// and numbers differently from encoding/json. Encoding against a fixed
		// small object avoids copying the growing output for every extra field.
		type encodedExtra struct {
			value  []byte
			member []byte
		}
		fields := make(map[string]encodedExtra, len(keys))
		for _, k := range keys {
			v := extras[k]
			if any(v) == Omit {
				fields[k] = encodedExtra{}
				continue
			}
			encoded, err := sjson.SetBytes([]byte("{}"), "v", v)
			if err != nil {
				return nil, err
			}
			value := encoded[len(`{"v":`) : len(encoded)-1]
			member, err := sjson.SetRawBytes([]byte("{}"), EscapeSJSONKey(k), value)
			if err != nil {
				return nil, err
			}
			fields[k] = encodedExtra{value: value, member: member[1 : len(member)-1]}
		}

		merged := make([]byte, 0, len(bytes))
		merged = append(merged, '{')
		object.ForEach(func(key, value gjson.Result) bool {
			field, replaced := fields[key.Str]
			if replaced {
				delete(fields, key.Str)
				if field.value == nil { // Omit
					return true
				}
			}
			if len(merged) > 1 {
				merged = append(merged, ',')
			}
			merged = append(merged, key.Raw...)
			merged = append(merged, ':')
			if replaced {
				merged = append(merged, field.value...)
			} else {
				merged = append(merged, value.Raw...)
			}
			return true
		})
		for _, k := range keys {
			if field, remaining := fields[k]; remaining && field.value != nil {
				if len(merged) > 1 {
					merged = append(merged, ',')
				}
				merged = append(merged, field.member...)
			}
		}
		return append(merged, '}'), nil
	} else if ovr, ok := f.Overrides(); ok {
		return shimjson.Marshal(ovr)
	} else {
		return shimjson.Marshal(underlying, shimjson.WithSkipCompaction(true))
	}
}

// MarshalUnion uses a shimmed 'encoding/json' from Go 1.24, to support the 'omitzero' tag
//
// Stability for the API of MarshalUnion is not guaranteed.
func MarshalUnion[T ParamStruct](metadata T, variants ...any) ([]byte, error) {
	nPresent := 0
	presentIdx := -1
	for i, variant := range variants {
		if !IsOmitted(variant) {
			nPresent++
			presentIdx = i
		}
	}
	if nPresent == 0 || presentIdx == -1 {
		if metadata.null() {
			return []byte("null"), nil
		}
		if ovr, ok := metadata.Overrides(); ok {
			return shimjson.Marshal(ovr)
		}
		return []byte(`null`), nil
	} else if nPresent > 1 {
		return nil, &json.MarshalerError{
			Type: typeFor[T](),
			Err:  fmt.Errorf("expected union to have only one present variant, got %d", nPresent),
		}
	}
	return shimjson.Marshal(variants[presentIdx], shimjson.WithSkipCompaction(true))
}

// typeFor is shimmed from Go 1.23 "reflect" package
func typeFor[T any]() reflect.Type {
	var v T
	if t := reflect.TypeOf(v); t != nil {
		return t // optimize for T being a non-interface kind
	}
	return reflect.TypeOf((*T)(nil)).Elem() // only for an interface kind
}

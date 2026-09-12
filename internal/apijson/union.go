package apijson

import (
	"errors"
	"reflect"

	"github.com/tidwall/gjson"

	"github.com/openai/openai-go/v3/packages/param"
)

var apiUnionType = reflect.TypeOf(param.APIUnion{})

// param.SetJSON accepts an unexported setter interface. The decoder only has a
// runtime reflect.Value, so invoke it reflectively after verifying the pointer.
var setParamJSON = reflect.ValueOf(param.SetJSON)

func preserveUnknownParamUnion(v reflect.Value, raw string) bool {
	if !v.CanAddr() || !v.Addr().Type().Implements(setParamJSON.Type().In(1)) {
		return false
	}
	setParamJSON.Call([]reflect.Value{reflect.ValueOf([]byte(raw)), v.Addr()})
	return true
}

func isStructUnion(t reflect.Type) bool {
	if t.Kind() != reflect.Struct {
		return false
	}
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Type == apiUnionType && t.Field(i).Anonymous {
			return true
		}
	}
	return false
}

func RegisterDiscriminatedUnion[T any](key string, mappings map[string]reflect.Type) {
	var t T
	entry := unionEntry{
		discriminatorKey: key,
		variants:         []UnionVariant{},
	}
	for k, typ := range mappings {
		entry.variants = append(entry.variants, UnionVariant{
			DiscriminatorValue: k,
			Type:               typ,
		})
	}
	unionRegistry[reflect.TypeOf(t)] = entry
}

func (d *decoderBuilder) newStructUnionDecoder(t reflect.Type) decoderFunc {
	type variantDecoder struct {
		decoder decoderFunc
		field   reflect.StructField
	}
	decoders := []variantDecoder{}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if field.Anonymous && field.Type == apiUnionType {
			continue
		}

		decoder := d.typeDecoder(field.Type)
		decoders = append(decoders, variantDecoder{
			decoder: decoder,
			field:   field,
		})
	}

	type discriminatedDecoder struct {
		variantDecoder
		discriminator any
	}
	discriminatedDecoders := []discriminatedDecoder{}
	unionEntry, discriminated := unionRegistry[t]
	for _, variant := range unionEntry.variants {
		// For each union variant, find a matching decoder and save it
		for _, decoder := range decoders {
			if decoder.field.Type.Elem() == variant.Type {
				discriminatedDecoders = append(discriminatedDecoders, discriminatedDecoder{
					decoder,
					variant.DiscriminatorValue,
				})
				break
			}
		}
	}

	return func(n gjson.Result, v reflect.Value, state *decoderState) error {
		if discriminated && n.Type == gjson.JSON && len(unionEntry.discriminatorKey) != 0 {
			discriminatorNode := n.Get(EscapeSJSONKey(unionEntry.discriminatorKey))
			discriminator := discriminatorNode.Value()
			compatibleDiscriminatorType := false
			for _, decoder := range discriminatedDecoders {
				if reflect.TypeOf(discriminator) == reflect.TypeOf(decoder.discriminator) {
					compatibleDiscriminatorType = true
				}
				if discriminator == decoder.discriminator {
					inner := v.FieldByIndex(decoder.field.Index)
					// Preservation only applies to the union that is itself an array
					// element. Nested arrays will opt their own elements in.
					preserveUnknown := state.preserveUnknownDiscriminatedUnion
					state.preserveUnknownDiscriminatedUnion = false
					err := decoder.decoder(n, inner, state)
					state.preserveUnknownDiscriminatedUnion = preserveUnknown
					return err
				}
			}
			if state.preserveUnknownDiscriminatedUnion && discriminatorNode.Exists() && compatibleDiscriminatorType && preserveUnknownParamUnion(v, n.Raw) {
				return nil
			}
			return errors.New("apijson: was not able to find discriminated union variant")
		}

		// Set bestExactness to worse than loose
		bestExactness := loose - 1
		bestVariant := -1
		for i, decoder := range decoders {
			// Pointers are used to discern JSON object variants from value variants
			if n.Type != gjson.JSON && decoder.field.Type.Kind() == reflect.Pointer {
				continue
			}

			sub := decoderState{strict: state.strict, exactness: exact}
			inner := v.FieldByIndex(decoder.field.Index)
			err := decoder.decoder(n, inner, &sub)
			if err != nil {
				continue
			}
			if sub.exactness == exact {
				bestExactness = exact
				bestVariant = i
				break
			}
			if sub.exactness > bestExactness {
				bestExactness = sub.exactness
				bestVariant = i
			}
		}

		if bestExactness < loose {
			return errors.New("apijson: was not able to coerce type as union")
		}

		if guardStrict(state, bestExactness != exact) {
			return errors.New("apijson: was not able to coerce type as union strictly")
		}

		for i := 0; i < len(decoders); i++ {
			if i == bestVariant {
				continue
			}
			v.FieldByIndex(decoders[i].field.Index).SetZero()
		}

		return nil
	}
}

// newUnionDecoder returns a decoderFunc that deserializes into a union using an
// algorithm roughly similar to Pydantic's [smart algorithm].
//
// Conceptually this is equivalent to choosing the best schema based on how 'exact'
// the deserialization is for each of the schemas.
//
// If there is a tie in the level of exactness, then the tie is broken
// left-to-right.
//
// [smart algorithm]: https://docs.pydantic.dev/latest/concepts/unions/#smart-mode
func (d *decoderBuilder) newUnionDecoder(t reflect.Type) decoderFunc {
	unionEntry, ok := unionRegistry[t]
	if !ok {
		panic("apijson: couldn't find union of type " + t.String() + " in union registry")
	}
	decoders := []decoderFunc{}
	for _, variant := range unionEntry.variants {
		decoder := d.typeDecoder(variant.Type)
		decoders = append(decoders, decoder)
	}
	return func(n gjson.Result, v reflect.Value, state *decoderState) error {
		// If there is a discriminator match, circumvent the exactness logic entirely
		for idx, variant := range unionEntry.variants {
			decoder := decoders[idx]
			if variant.TypeFilter != n.Type {
				continue
			}

			if len(unionEntry.discriminatorKey) != 0 {
				discriminatorValue := n.Get(EscapeSJSONKey(unionEntry.discriminatorKey)).Value()
				if discriminatorValue == variant.DiscriminatorValue {
					inner := reflect.New(variant.Type).Elem()
					err := decoder(n, inner, state)
					v.Set(inner)
					return err
				}
			}
		}

		// Set bestExactness to worse than loose
		bestExactness := loose - 1
		for idx, variant := range unionEntry.variants {
			decoder := decoders[idx]
			if variant.TypeFilter != n.Type {
				continue
			}
			sub := decoderState{strict: state.strict, exactness: exact}
			inner := reflect.New(variant.Type).Elem()
			err := decoder(n, inner, &sub)
			if err != nil {
				continue
			}
			if sub.exactness == exact {
				v.Set(inner)
				return nil
			}
			if sub.exactness > bestExactness {
				v.Set(inner)
				bestExactness = sub.exactness
			}
		}

		if bestExactness < loose {
			return errors.New("apijson: was not able to coerce type as union")
		}

		if guardStrict(state, bestExactness != exact) {
			return errors.New("apijson: was not able to coerce type as union strictly")
		}

		return nil
	}
}

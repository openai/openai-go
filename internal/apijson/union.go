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

func preserveParamJSON(v reflect.Value, raw string) bool {
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
	// Some variants share a primary tag (for example, messages with different
	// roles). A shared constant field can distinguish them without trying every
	// decoder. Keep first-match behavior when no such field or value is present.
	type secondaryDiscriminator struct {
		key      string
		variants map[string]int
	}
	secondaryDiscriminators := map[any]secondaryDiscriminator{}
	groups := map[any][]int{}
	for i, decoder := range discriminatedDecoders {
		groups[decoder.discriminator] = append(groups[decoder.discriminator], i)
	}
	for primary, indices := range groups {
		if len(indices) < 2 {
			continue
		}
		types := make([]reflect.Type, len(indices))
		for i, index := range indices {
			types[i] = discriminatedDecoders[index].field.Type.Elem()
		}
		key, values := sharedUnionConstant(types, unionEntry.discriminatorKey)
		if key == "" {
			continue
		}
		secondary := secondaryDiscriminator{key: EscapeSJSONKey(key), variants: map[string]int{}}
		for i, value := range values {
			secondary.variants[value] = indices[i]
		}
		secondaryDiscriminators[primary] = secondary
	}

	return func(n gjson.Result, v reflect.Value, state *decoderState) error {
		if discriminated && n.Type == gjson.JSON && len(unionEntry.discriminatorKey) != 0 {
			discriminatorNode := n.Get(EscapeSJSONKey(unionEntry.discriminatorKey))
			discriminator := discriminatorNode.Value()
			preferred := -1
			if discriminatorNode.Type != gjson.JSON {
				if secondary, ok := secondaryDiscriminators[discriminator]; ok {
					value := n.Get(secondary.key)
					if index, ok := secondary.variants[value.Str]; ok && value.Type == gjson.String {
						preferred = index
					}
				}
			}
			compatibleDiscriminatorType := false
			for i, decoder := range discriminatedDecoders {
				if reflect.TypeOf(discriminator) == reflect.TypeOf(decoder.discriminator) {
					compatibleDiscriminatorType = true
				}
				if discriminator == decoder.discriminator && (preferred < 0 || preferred == i) {
					inner := v.FieldByIndex(decoder.field.Index)
					// Preservation only applies to the union that is itself an array
					// element. Nested arrays will opt their own elements in.
					preserveUnknown := state.preserveUnknownDiscriminatedUnion
					state.preserveUnknownDiscriminatedUnion = false
					err := decoder.decoder(n, inner, state)
					state.preserveUnknownDiscriminatedUnion = preserveUnknown
					if err == nil {
						for _, other := range decoders {
							if other.field.Name != decoder.field.Name {
								v.FieldByIndex(other.field.Index).SetZero()
							}
						}
					}
					return err
				}
			}
			if state.preserveUnknownDiscriminatedUnion && discriminatorNode.Exists() && compatibleDiscriminatorType && preserveParamJSON(v, n.Raw) {
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

// sharedUnionConstant finds a constant field with a distinct value in every
// candidate. Ordinary string defaults are not discriminator evidence.
func sharedUnionConstant(types []reflect.Type, primary string) (string, []string) {
	for i := 0; i < types[0].NumField(); i++ {
		field := types[0].Field(i)
		name, _, ok := unionStringConstant(field)
		if !ok || name == primary {
			continue
		}
		values := make([]string, 0, len(types))
		seen := map[string]bool{}
		for _, typ := range types {
			candidate, found := typ.FieldByName(field.Name)
			if !found {
				break
			}
			key, value, ok := unionStringConstant(candidate)
			if !ok || key != name || seen[value] {
				break
			}
			seen[value] = true
			values = append(values, value)
		}
		if len(values) == len(types) {
			return name, values
		}
	}
	return "", nil
}

func unionStringConstant(field reflect.StructField) (string, string, bool) {
	tag, ok := parseJSONStructTag(field)
	value, hasDefault := tag.defaultValue.(string)
	method, hasMethod := field.Type.MethodByName("Default")
	if !ok || !hasDefault || tag.name == "" || tag.name == "-" || !hasMethod ||
		method.Type.NumIn() != 1 || method.Type.NumOut() != 1 || method.Type.Out(0) != field.Type {
		return "", "", false
	}
	return tag.name, value, true
}

// Inline response variants can share a primary tag. Select a single variant
// only when another declared constant distinguishes it; retain loose decoding
// for unknown or absent discriminators.
func inlineUnionSelector(t reflect.Type, fields []decoderField) func(gjson.Result) int {
	type constant struct{ key, value string }
	type selector struct {
		primary   constant
		secondary string
		variants  map[string]int
	}
	groups := map[constant][]int{}
	for i, field := range fields {
		typ := t.FieldByIndex(field.idx).Type
		if typ.Kind() != reflect.Struct {
			continue
		}
		for j := 0; j < typ.NumField(); j++ {
			if key, value, ok := unionStringConstant(typ.Field(j)); ok {
				c := constant{key, value}
				groups[c] = append(groups[c], i)
			}
		}
	}
	var selectors []selector
	for primary, indices := range groups {
		if len(indices) < 2 {
			continue
		}
		types := make([]reflect.Type, len(indices))
		for i, index := range indices {
			types[i] = t.FieldByIndex(fields[index].idx).Type
		}
		key, values := sharedUnionConstant(types, primary.key)
		if key == "" {
			continue
		}
		s := selector{constant{EscapeSJSONKey(primary.key), primary.value}, EscapeSJSONKey(key), map[string]int{}}
		for i, value := range values {
			s.variants[value] = indices[i]
		}
		selectors = append(selectors, s)
	}
	return func(node gjson.Result) int {
		selected := -1
		for _, s := range selectors {
			primary, secondary := node.Get(s.primary.key), node.Get(s.secondary)
			if primary.Type != gjson.String || primary.Str != s.primary.value || secondary.Type != gjson.String {
				continue
			}
			if index, ok := s.variants[secondary.Str]; ok {
				if selected >= 0 && selected != index {
					return -1
				}
				selected = index
			}
		}
		return selected
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

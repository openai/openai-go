package resp

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

var fieldsMap sync.Map

type structField struct {
	Index         []int
	metadataIndex []int
	inline        bool
	name          string
	format        string
}

func newObjectPreprocessor(t reflect.Type) ([]structField, error) {
	var structFields []structField

	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("resp: expected struct but got %v of kind %v", t.String(), t.Kind().String())
	}

	meta, ok := t.FieldByName("JSON")
	if !ok || meta.Type.Kind() != reflect.Struct {
		return nil, fmt.Errorf("resp: expected struct %v to have JSON metadata field", t.String())
	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		fmeta, hasMeta := meta.Type.FieldByName(f.Name)
		if !hasMeta {
			continue
		}

		var sf structField
		sf.Index = f.Index
		sf.metadataIndex = fmeta.Index
		sf.name, sf.inline = collectTags(f)
		if sf.name == "" && !sf.inline {
			continue
		}

		structFields = append(structFields, sf)
	}

	return structFields, nil
}

func collectTags(f reflect.StructField) (name string, inline bool) {
	tag, ok := f.Tag.Lookup("json")
	if !ok {
		return
	}

	tags := strings.Split(tag, ",")
	if len(tags) == 0 {
		return
	}

	if tags[0] != "-" {
		name = tags[0]
	}

	for _, tag := range tags[1:] {
		if tag == "inline" {
			inline = true
		}
	}

	return
}

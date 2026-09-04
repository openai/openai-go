// File generated from our OpenAPI spec by Castiron. See CONTRIBUTING.md for details.

package apiform

import (
	"encoding/json"
	"fmt"
	"mime"
	"mime/multipart"
	"net/textproto"
	"reflect"
	"sort"
	"time"

	"github.com/openai/openai-go/v3/packages/param"
)

type PartEncoding struct {
	ContentType string
	JSON        bool
}

// MarshalEncodedRoot preserves annotated fields as single typed parts. Ordinary
// fields still use the existing form encoder, including file and array handling.
func MarshalEncodedRoot(value any, writer *multipart.Writer, extras map[string]any, encodings map[string]PartEncoding) error {
	type fieldValue struct {
		value      any
		dateFormat string
	}
	fields := map[string]fieldValue{}
	var collect func(reflect.Value)
	collect = func(v reflect.Value) {
		if v.Kind() == reflect.Pointer {
			if v.IsNil() {
				return
			}
			v = v.Elem()
		}
		for i := 0; i < v.NumField(); i++ {
			field := v.Type().Field(i)
			if !field.IsExported() {
				continue
			}
			fv := v.Field(i)
			if field.Anonymous {
				collect(fv)
				continue
			}
			tag, ok := parseFormStructTag(field)
			if !ok || tag.name == "-" || tag.name == "" {
				continue
			}
			if tag.omitzero && fv.IsZero() {
				continue
			}
			fieldData := fv.Interface()
			if tag.defaultValue != nil && fv.IsZero() {
				fieldData = tag.defaultValue
			}
			dateFormat := time.RFC3339
			if format, ok := parseFormatStructTag(field); ok && format == "date" {
				dateFormat = "2006-01-02"
			}
			fields[tag.name] = fieldValue{fieldData, dateFormat}
		}
	}
	collect(reflect.ValueOf(value))
	for key, value := range extras {
		if value == param.Omit {
			delete(fields, key)
		} else {
			fields[key] = fieldValue{value, time.RFC3339}
		}
	}
	keys := make([]string, 0, len(fields))
	for key := range fields {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		field := fields[key]
		encoding, encoded := encodings[key]
		if !encoded {
			e := &encoder{dateFormat: field.dateFormat, arrayFmt: "brackets"}
			v := reflect.ValueOf(field.value)
			if !v.IsValid() {
				continue
			}
			if err := e.typeEncoder(v.Type())(key, v, writer); err != nil {
				return err
			}
			continue
		}
		data, err := json.Marshal(field.value)
		if err != nil {
			return err
		}
		if !encoding.JSON {
			var text *string
			if decodeErr := json.Unmarshal(data, &text); decodeErr != nil {
				return fmt.Errorf("multipart field %s must be a string: %w", key, decodeErr)
			}
			if text == nil {
				return fmt.Errorf("multipart field %s must be a string", key)
			}
			data = []byte(*text)
		}
		header := make(textproto.MIMEHeader)
		header.Set("Content-Disposition", mime.FormatMediaType("form-data", map[string]string{"name": key}))
		header.Set("Content-Type", encoding.ContentType)
		part, err := writer.CreatePart(header)
		if err != nil {
			return err
		}
		if _, writeErr := part.Write(data); writeErr != nil {
			return writeErr
		}
	}
	return nil
}

package apijson

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
)

type rootPointerShape struct {
	Value string `json:"value"`
}

func TestRootPointerStructRejectsNonObjectJSON(t *testing.T) {
	for _, input := range []string{`[]`, `true`, `42`, `"value"`} {
		t.Run(input, func(t *testing.T) {
			var got *rootPointerShape
			err := UnmarshalRoot([]byte(input), &got)
			if err == nil {
				t.Fatalf("UnmarshalRoot(%s) succeeded with %#v", input, got)
			}
			var typeErr *json.UnmarshalTypeError
			if !errors.As(err, &typeErr) {
				t.Fatalf("error type = %T, want *json.UnmarshalTypeError", err)
			}
			if typeErr.Type != reflect.TypeOf(rootPointerShape{}) {
				t.Fatalf("error target type = %v", typeErr.Type)
			}
		})
	}
}

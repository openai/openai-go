package apijson

import (
	"reflect"
	"testing"

	"github.com/openai/openai-go/v3/packages/respjson"
)

func TestResponseMapArrayUnionShape(t *testing.T) {
	for _, raw := range []string{`{}`, `{"ok":true}`, `[]`, `[{"ok":true}]`} {
		t.Run(raw, func(t *testing.T) {
			var got struct {
				OfMap   map[string]any   `json:",inline"`
				OfArray []map[string]any `json:",inline"`
				JSON    struct {
					OfMap   respjson.Field
					OfArray respjson.Field
				} `json:"-"`
			}
			if err := UnmarshalRoot([]byte(raw), &got); err != nil {
				t.Fatal(err)
			}
			object := raw[0] == '{'
			if got.JSON.OfMap.Valid() != object || got.JSON.OfArray.Valid() == object {
				t.Fatalf("wrong variant presence: map=%v array=%v", got.JSON.OfMap.Valid(), got.JSON.OfArray.Valid())
			}
			if object && got.OfMap == nil || !object && got.OfArray == nil {
				t.Fatal("selected variant is nil")
			}
			if raw == `{"ok":true}` && !reflect.DeepEqual(got.OfMap, map[string]any{"ok": true}) {
				t.Fatalf("object data lost: %#v", got.OfMap)
			}
			if raw == `[{"ok":true}]` && !reflect.DeepEqual(got.OfArray, []map[string]any{{"ok": true}}) {
				t.Fatalf("array data lost: %#v", got.OfArray)
			}
		})
	}
}

func TestMapDecoderRejectsNonObjects(t *testing.T) {
	for _, raw := range []string{`[]`, `[{}]`, `"secret-value"`, `42`, `true`, `false`} {
		t.Run(raw, func(t *testing.T) {
			got := map[string]any{"preserved": true}
			if err := UnmarshalRoot([]byte(raw), &got); err == nil {
				t.Fatal("non-object unexpectedly accepted as a map")
			}
			if !reflect.DeepEqual(got, map[string]any{"preserved": true}) {
				t.Fatal("failed decode changed destination")
			}
		})
	}
}

func TestMapDecoderPreservesNullField(t *testing.T) {
	var got struct {
		Value map[string]any                 `json:"value"`
		JSON  struct{ Value respjson.Field } `json:"-"`
	}
	if err := UnmarshalRoot([]byte(`{"value":null}`), &got); err != nil {
		t.Fatal(err)
	}
	if got.Value != nil || got.JSON.Value.Raw() != "null" || got.JSON.Value.Valid() {
		t.Fatal("null map field lost its null metadata")
	}
}

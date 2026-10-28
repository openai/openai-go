package openai_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/openai/openai-go/v3/shared"
)

func TestDecisionJSONResponsePreservesShape(t *testing.T) {
	for _, tc := range []struct {
		name      string
		data      string
		wantMap   map[string]any
		wantArray []map[string]any
	}{
		{"empty object", `{}`, map[string]any{}, nil},
		{"object", `{"ok":true}`, map[string]any{"ok": true}, nil},
		{"empty array", `[]`, nil, []map[string]any{}},
		{"array", `[{"ok":true}]`, nil, []map[string]any{{"ok": true}}},
		{"array of empty objects", `[{}]`, nil, []map[string]any{{}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var got shared.DecisionInputJSON
			raw := `{"type":"input_json","data":` + tc.data + `}`
			if err := json.Unmarshal([]byte(raw), &got); err != nil {
				t.Fatal(err)
			}
			if got.Data.JSON.OfDecisionInputJSONDataMapItem.Valid() != (tc.wantMap != nil) ||
				got.Data.JSON.OfArrayOfMap.Valid() != (tc.wantArray != nil) {
				t.Fatalf("wrong variant presence for %s", tc.data)
			}
			if !reflect.DeepEqual(got.Data.OfDecisionInputJSONDataMapItem, tc.wantMap) ||
				!reflect.DeepEqual(got.Data.AsMap(), tc.wantMap) {
				t.Fatalf("object data changed for %s", tc.data)
			}
			if !reflect.DeepEqual(got.Data.OfArrayOfMap, tc.wantArray) ||
				!reflect.DeepEqual(got.Data.AsArrayOfMap(), tc.wantArray) {
				t.Fatalf("array data changed for %s", tc.data)
			}
			if got.Data.RawJSON() != tc.data || got.RawJSON() != raw {
				t.Fatalf("raw JSON changed for %s", tc.data)
			}
		})
	}
}

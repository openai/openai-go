package openai_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/openai/openai-go/v4/live"
)

func TestLiveDecisionJSONResponsePreservesShape(t *testing.T) {
	for _, tc := range []struct {
		name      string
		data      string
		wantMap   live.DecisionJSONObject
		wantArray []live.DecisionJSONObject
	}{
		{"empty object", `{}`, live.DecisionJSONObject{}, nil},
		{"object", `{"ok":true}`, live.DecisionJSONObject{"ok": true}, nil},
		{"empty array", `[]`, nil, []live.DecisionJSONObject{}},
		{"array", `[{"ok":true}]`, nil, []live.DecisionJSONObject{{"ok": true}}},
		{"array of empty objects", `[{}]`, nil, []live.DecisionJSONObject{{}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var got live.DecisionQuestionChoiceInstructionsPartInputJSON
			raw := `{"type":"input_json","data":` + tc.data + `}`
			if err := json.Unmarshal([]byte(raw), &got); err != nil {
				t.Fatal(err)
			}
			if got.Data.JSON.OfLiveDecisionJSONObject.Valid() != (tc.wantMap != nil) ||
				got.Data.JSON.OfArrayOfDecisionJSONObject.Valid() != (tc.wantArray != nil) {
				t.Fatalf("wrong variant presence for %s", tc.data)
			}
			if !reflect.DeepEqual(got.Data.OfLiveDecisionJSONObject, tc.wantMap) ||
				!reflect.DeepEqual(got.Data.AsLiveDecisionJSONObject(), tc.wantMap) {
				t.Fatalf("object data changed for %s", tc.data)
			}
			if !reflect.DeepEqual(got.Data.OfArrayOfDecisionJSONObject, tc.wantArray) ||
				!reflect.DeepEqual(got.Data.AsArrayOfDecisionJSONObject(), tc.wantArray) {
				t.Fatalf("array data changed for %s", tc.data)
			}
			if got.Data.RawJSON() != tc.data || got.RawJSON() != raw {
				t.Fatalf("raw JSON changed for %s", tc.data)
			}
		})
	}
}

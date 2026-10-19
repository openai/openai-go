package openai_test

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/conversations"
	"github.com/openai/openai-go/v3/packages/respjson"
	"github.com/openai/openai-go/v3/responses"
)

// Keep malformed actions in field metadata without treating them as decoded
// variants. Reusing each receiver also catches stale inline metadata.
func TestImageGenerationUnionAction(t *testing.T) {
	for _, model := range []struct {
		name        string
		value       any
		scalarField string
	}{
		{"stable_input", &responses.ResponseInputItemUnion{}, "OfResponseInputItemImageGenerationCallAction"},
		{"stable_item", &responses.ResponseItemUnion{}, "OfResponseItemImageGenerationCallAction"},
		{"stable_output", &responses.ResponseOutputItemUnion{}, "OfResponseOutputItemImageGenerationCallAction"},
		{"beta_input", &openai.BetaResponseInputItemUnion{}, "OfBetaResponseInputItemMultiAgentCallOutputAction"},
		{"beta_item", &openai.BetaResponseItemUnion{}, "OfBetaResponseItemMultiAgentCallOutputAction"},
		{"beta_output", &openai.BetaResponseOutputItemUnion{}, "OfBetaResponseOutputItemMultiAgentCallOutputAction"},
		{"conversation_item", &conversations.ConversationItemUnion{}, "OfConversationItemImageGenerationCallAction"},
	} {
		t.Run(model.name, func(t *testing.T) {
			for _, tc := range []struct {
				raw    string
				valid  bool
				scalar *string
			}{
				{`"generate"`, true, imageActionStringPointer("generate")},
				{`"edit"`, true, imageActionStringPointer("edit")},
				{`"future_action"`, true, imageActionStringPointer("future_action")},
				{`""`, true, imageActionStringPointer("")},
				{`{"type":"search","query":"example"}`, true, nil},
				{`null`, false, nil},
				{`true`, false, nil},
				{`42`, false, nil},
				{`[]`, false, nil},
				{`"generate"`, true, imageActionStringPointer("generate")},
			} {
				t.Run(tc.raw, func(t *testing.T) {
					data := []byte(`{"type":"image_generation_call","action":` + tc.raw + `}`)
					if err := json.Unmarshal(data, model.value); err != nil {
						t.Fatalf("json.Unmarshal(%s) error = %v, want nil", tc.raw, err)
					}
					item := reflect.ValueOf(model.value).Elem()
					meta := item.FieldByName("JSON").FieldByName("Action").Interface().(respjson.Field)
					if meta.Valid() != tc.valid || meta.Raw() != tc.raw {
						t.Errorf("json.Unmarshal(%s) action metadata = valid %v, raw %q; want valid %v, raw %q", tc.raw, meta.Valid(), meta.Raw(), tc.valid, tc.raw)
					}
					action := item.FieldByName("Action")
					validArms := 0
					for i := 0; i < action.NumField(); i++ {
						field := action.Type().Field(i)
						if !strings.Contains(field.Tag.Get("json"), "inline") {
							continue
						}
						armMeta := action.FieldByName("JSON").FieldByName(field.Name).Interface().(respjson.Field)
						if !armMeta.Valid() {
							continue
						}
						validArms++
						if field.Name != model.scalarField {
							t.Errorf("json.Unmarshal(%s) valid inline arm = %s, want %s", tc.raw, field.Name, model.scalarField)
						}
						if tc.scalar == nil || action.Field(i).String() != *tc.scalar {
							t.Errorf("json.Unmarshal(%s) decoded unexpected inline arm %s = %v", tc.raw, field.Name, action.Field(i))
						}
					}
					wantArms := 0
					if tc.scalar != nil {
						wantArms = 1
					}
					if validArms != wantArms {
						t.Errorf("json.Unmarshal(%s) valid inline arms = %d, want %d", tc.raw, validArms, wantArms)
					}
					if !tc.valid && tc.raw != "null" {
						direct := reflect.New(action.Type()).Interface()
						if err := json.Unmarshal([]byte(`"generate"`), direct); err != nil {
							t.Fatalf("json.Unmarshal(generate) error = %v, want nil", err)
						}
						if err := json.Unmarshal([]byte(tc.raw), direct); err == nil {
							t.Errorf("direct json.Unmarshal(%s) error = nil, want rejection", tc.raw)
						}
					}
				})
			}
		})
	}
}

func imageActionStringPointer(value string) *string { return &value }

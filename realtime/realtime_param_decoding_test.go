package realtime

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/openai/openai-go/v3/internal/apijson"
)

func TestConversationItemDecodingUsesMessageRole(t *testing.T) {
	var reused ConversationItemUnionParam
	for _, tc := range []struct{ wire, field string }{
		{`{"type":"message","role":"system","content":[{"type":"input_text","text":"hello"}]}`, "OfRealtimeSystemMessageItem"},
		{`{"type":"message","role":"user","content":[{"type":"input_text","text":"hello"}]}`, "OfRealtimeUserMessageItem"},
		{`{"type":"message","role":"user","content":[{"type":"input_audio","audio":"AAAA"}]}`, "OfRealtimeUserMessageItem"},
		{`{"type":"message","role":"user","content":[{"type":"input_image","image_url":"https://example.invalid/image"}]}`, "OfRealtimeUserMessageItem"},
		{`{"type":"message","role":"assistant","content":[{"type":"output_text","text":"hello"}]}`, "OfRealtimeAssistantMessageItem"},
		{`{"type":"function_call","name":"hello","arguments":"{}"}`, "OfFunctionCall"},
	} {
		t.Run(tc.wire, func(t *testing.T) {
			var item ConversationItemUnionParam
			if err := json.Unmarshal([]byte(tc.wire), &item); err != nil {
				t.Fatal(err)
			}
			if err := json.Unmarshal([]byte(tc.wire), &reused); err != nil {
				t.Fatal(err)
			}
			// Array elements use a separate decoder path and must select the same variant.
			var items []ConversationItemUnionParam
			if err := apijson.UnmarshalRoot([]byte("["+tc.wire+"]"), &items); err != nil {
				t.Fatal(err)
			}
			var event ConversationItemCreateEventParam
			if err := json.Unmarshal([]byte(`{"type":"conversation.item.create","item":`+tc.wire+`}`), &event); err != nil {
				t.Fatal(err)
			}
			for _, got := range []ConversationItemUnionParam{item, items[0], event.Item, reused} {
				if reflect.ValueOf(got).FieldByName(tc.field).IsNil() {
					t.Fatalf("did not select %s: %#v", tc.field, got)
				}
				wire, err := json.Marshal(got)
				if err != nil {
					t.Fatal(err)
				}
				var actual, expected any
				if err := json.Unmarshal(wire, &actual); err != nil {
					t.Fatal(err)
				}
				if err := json.Unmarshal([]byte(tc.wire), &expected); err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(actual, expected) {
					t.Fatalf("roundtrip = %s, want %s", wire, tc.wire)
				}
			}
		})
	}
}

func TestConversationItemUnknownDiscriminatorsKeepExistingBehavior(t *testing.T) {
	for _, wire := range []string{
		`{"type":"message","content":[]}`,
		`{"type":"message","role":"future-role","content":[]}`,
	} {
		var item ConversationItemUnionParam
		if err := json.Unmarshal([]byte(wire), &item); err != nil || item.OfRealtimeSystemMessageItem == nil {
			t.Fatalf("fallback changed for %s: %#v, %v", wire, item, err)
		}
	}
	for _, wire := range []string{
		`{"type":"future-item"}`,
		`{"type":{}}`,
		`{"type":[]}`,
		`{}`,
	} {
		var item ConversationItemUnionParam
		if err := json.Unmarshal([]byte(wire), &item); err == nil {
			t.Errorf("accepted unknown discriminator in %s", wire)
		}
	}
}

func TestNestedTurnDetectionGetters(t *testing.T) {
	zero, half := 0.0, 0.5
	type turnDetection interface {
		AsAny() any
		GetType() *string
		GetThreshold() *float64
		GetEagerness() *string
	}
	for _, sessionType := range []string{"realtime", "transcription"} {
		for _, tc := range []struct {
			wire, kind, eagerness string
			threshold             *float64
		}{
			{wire: `{}`},
			{wire: `{"type":"server_vad","threshold":0}`, kind: "server_vad", threshold: &zero},
			{wire: `{"type":"server_vad","threshold":0.5}`, kind: "server_vad", threshold: &half},
			{wire: `{"type":"semantic_vad","eagerness":"high"}`, kind: "semantic_vad", eagerness: "high"},
		} {
			t.Run(sessionType+tc.wire, func(t *testing.T) {
				input := `"turn_detection":` + tc.wire
				if tc.kind == "" {
					input = ""
				}
				wire := []byte(`{"type":"` + sessionType + `","audio":{"input":{` + input + `}}}`)
				var session SessionUpdateEventSessionUnionParam
				var secret ClientSecretNewParamsSessionUnion
				if err := json.Unmarshal(wire, &session); err != nil {
					t.Fatal(err)
				}
				if err := json.Unmarshal(wire, &secret); err != nil {
					t.Fatal(err)
				}
				for _, view := range []turnDetection{
					session.GetAudio().GetInput().GetTurnDetection(),
					secret.GetAudio().GetInput().GetTurnDetection(),
				} {
					if reflect.TypeOf(view.AsAny()).Kind() != reflect.Struct {
						t.Fatalf("AsAny no longer returns the union value: %T", view.AsAny())
					}
					kind := view.GetType()
					if (tc.kind == "" && kind != nil) || (tc.kind != "" && (kind == nil || *kind != tc.kind)) {
						t.Errorf("GetType = %v, want %q", kind, tc.kind)
					}
					if !reflect.DeepEqual(view.GetThreshold(), tc.threshold) {
						t.Errorf("GetThreshold = %v, want %v", view.GetThreshold(), tc.threshold)
					}
					eagerness := view.GetEagerness()
					if (tc.eagerness == "" && eagerness != nil) || (tc.eagerness != "" && (eagerness == nil || *eagerness != tc.eagerness)) {
						t.Errorf("GetEagerness = %v, want %q", eagerness, tc.eagerness)
					}
				}
			})
		}
	}
}

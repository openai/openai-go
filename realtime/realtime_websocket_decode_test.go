package realtime

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

// A syntactically valid but malformed nested payload must keep the public
// decoder's acceptance and typed-field metadata. Only wire envelope errors
// are rejected by the existing reader.
func TestRealtimeWebsocketDecodeMalformedPayloadParity(t *testing.T) {
	for _, data := range []string{
		`{"type":"response.done","response":42,"event_id":["bad"]}`,
		`{"type":"session.created","session":[],"event_id":false}`,
		`{"type":"response.output_audio.delta","delta":{"bad":"type"},"content_index":"bad","event_id":"evt_fake"}`,
		`{"type":"error","error":"bad","event_id":9}`,
		`{"type":"error","error":{"type":7,"code":[],"message":null,"param":42}}`,
		`{"type":"future.event","event_id":false,"session":42}`,
		`{"Type":"response.done","response":42}`,
		`{"Type":"response.done","type":"error","error":{"message":"hello"}}`,
		`{"type":"error","Type":"response.done","error":{"message":"hello"}}`,
		`{"type":"error","type":"response.done","error":{"message":"hello"}}`,
	} {
		t.Run(data, func(t *testing.T) {
			frame, err := decodeRealtimeWebsocketEvent([]byte(data))
			if err != nil {
				t.Fatalf("existing reader rejected fixture: %v", err)
			}
			var general, actual RealtimeServerEventUnion
			wantErr := json.Unmarshal(frame.data, &general)
			gotErr := actual.unmarshalWebsocketEvent(frame.data)
			if wantErr != nil || gotErr != nil {
				t.Fatalf("valid wire with malformed payload should remain accepted: got %v, original decoder %v", gotErr, wantErr)
			}
			if actual.Type != general.Type || actual.RawJSON() != general.RawJSON() ||
				!reflect.DeepEqual(actual.JSON.Type, general.JSON.Type) ||
				!reflect.DeepEqual(actual.AsAny(), general.AsAny()) {
				t.Fatal("malformed payload lost public Type, raw JSON or selected field metadata")
			}
			if selected := reflect.ValueOf(general.AsAny()); selected.IsValid() {
				name := "Of" + selected.Type().Name()
				got, want := reflect.ValueOf(actual), reflect.ValueOf(general)
				if !reflect.DeepEqual(got.FieldByName(name).Interface(), want.FieldByName(name).Interface()) ||
					!reflect.DeepEqual(got.FieldByName("JSON").FieldByName(name).Interface(), want.FieldByName("JSON").FieldByName(name).Interface()) {
					t.Fatalf("%s malformed field or inline presence differs from public decoding", name)
				}
			}
		})
	}
}

func TestWebsocketRejectsInvalidFirstType(t *testing.T) {
	for name, decode := range map[string]func([]byte) error{
		"realtime": func(data []byte) error {
			_, err := decodeRealtimeWebsocketEvent(data)
			return err
		},
		"translation": func(data []byte) error {
			_, err := decodeTranslationWebsocketEvent(data)
			return err
		},
	} {
		t.Run(name, func(t *testing.T) {
			for _, data := range []string{
				`{"type":"","type":"session.created"}`,
				`{"type":null,"type":"session.created"}`,
			} {
				if err := decode([]byte(data)); err == nil {
					t.Errorf("accepted invalid first type in %s", data)
				}
			}
		})
	}
}

func BenchmarkRealtimeWebsocketConsumedEvent(b *testing.B) {
	for _, tc := range []struct{ name, wire string }{
		{"text-delta", `{"type":"response.output_text.delta","event_id":"evt_fake","delta":"hello","content_index":0,"output_index":0,"item_id":"item_fake","response_id":"resp_fake"}`},
		{"audio-64KiB", `{"type":"response.output_audio.delta","event_id":"evt_fake","delta":"` + strings.Repeat("a", 64<<10) + `","content_index":0,"output_index":0,"item_id":"item_fake","response_id":"resp_fake"}`},
		{"unknown-64KiB", `{"type":"future.event","future":"` + strings.Repeat("x", 64<<10) + `"}`},
	} {
		data := []byte(tc.wire)
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				var event RealtimeServerEventUnion
				if err := event.unmarshalWebsocketEvent(data); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

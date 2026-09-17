package responses

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/openai/openai-go/v3/packages/respjson"
)

func TestFinalWebsocketResponseTerminalShapes(t *testing.T) {
	for _, eventType := range []string{"response.completed", "response.failed", "response.incomplete"} {
		for _, test := range []struct {
			name    string
			fields  string
			wantID  string
			wantErr bool
		}{
			{name: "snapshot", fields: `"response":{"id":"resp_final","metadata":{"synthetic":"retained"}}`, wantID: "resp_final"},
			{name: "empty snapshot", fields: `"response":{}`},
			{name: "omitted", fields: `"sequence_number":1`, wantErr: true},
			{name: "null", fields: `"response":null`, wantErr: true},
			{name: "invalid snapshot", fields: `"response":42`, wantErr: true},
			{name: "invalid sequence", fields: `"sequence_number":"invalid","response":{"id":"resp_final"}`, wantID: "resp_final"},
			{name: "invalid snapshot field", fields: `"response":{"id":"resp_final","created_at":"invalid"}`, wantID: "resp_final"},
		} {
			t.Run(eventType+"/"+test.name, func(t *testing.T) {
				data := []byte(fmt.Sprintf(`{"type":%q,%s}`, eventType, test.fields))
				var event ResponsesServerEventUnion
				if err := json.Unmarshal(data, &event); err != nil {
					t.Fatalf("json.Unmarshal(%s) = %v, want success", data, err)
				}
				received := false
				response, err := finalWebsocketResponse(context.Background(), func(context.Context) (ResponsesServerEventUnion, error) {
					if received {
						return ResponsesServerEventUnion{}, io.EOF
					}
					received = true
					return event, nil
				})
				if (err != nil) != test.wantErr || errors.Is(err, io.EOF) {
					t.Fatalf("finalWebsocketResponse(%s) error = %v, want error presence %t before another receive", data, err, test.wantErr)
				}
				if !test.wantErr && (response == nil || response.ID != test.wantID) {
					t.Errorf("finalWebsocketResponse(%s) = %v, want response ID %q", data, response, test.wantID)
				}
			})
		}
	}
}

func TestResponseWebsocketWireValidation(t *testing.T) {
	for _, test := range []struct {
		name       string
		data       string
		wantStream string
		wantErr    bool
	}{
		{name: "malformed", data: `{"type":`, wantErr: true},
		{name: "missing type", data: `{"stream_id":"lane"}`, wantErr: true},
		{name: "invalid type", data: `{"type":42}`, wantErr: true},
		{name: "null", data: `null`, wantErr: true},
		{name: "unknown event", data: `{"type":"response.future","stream_id":"lane","extra":true}`, wantStream: "lane"},
		{name: "invalid stream", data: `{"type":"response.future","stream_id":42}`},
		{name: "null stream", data: `{"type":"response.future","stream_id":null}`},
		{name: "escaped stream", data: `{"type":"response.future","stream_id":"l\u0061ne"}`, wantStream: "lane"},
		{name: "invalid nested field", data: `{"type":"response.completed","response":42}`},
	} {
		t.Run(test.name, func(t *testing.T) {
			frame, err := decodeResponseWebsocketEvent([]byte(test.data))
			if (err != nil) != test.wantErr {
				t.Fatalf("decodeResponseWebsocketEvent(%s) error = %v, want error presence %t", test.data, err, test.wantErr)
			}
			if test.wantErr {
				return
			}
			if frame.streamID != test.wantStream {
				t.Errorf("decodeResponseWebsocketEvent(%s) stream = %q, want %q", test.data, frame.streamID, test.wantStream)
			}
			event, err := frame.decode()
			if err != nil || event.RawJSON() != test.data {
				t.Errorf("decode(%s) = raw %q, error %v, want unchanged raw JSON and nil error", test.data, event.RawJSON(), err)
			}
		})
	}
}

func BenchmarkFinalWebsocketResponse(b *testing.B) {
	text := strings.Repeat("x", 1<<20)
	data := []byte(`{"type":"response.completed","response":{"id":"resp_final","metadata":{"synthetic":"` + text + `"}}}`)
	var event ResponsesServerEventUnion
	if err := json.Unmarshal(data, &event); err != nil {
		b.Fatal(err)
	}
	receive := func(context.Context) (ResponsesServerEventUnion, error) { return event, nil }
	b.ReportAllocs()
	b.SetBytes(int64(len(text)))
	for b.Loop() {
		response, err := finalWebsocketResponse(context.Background(), receive)
		if err != nil || response.ID != "resp_final" {
			b.Fatalf("finalWebsocketResponse(1 MiB metadata) = %v, %v, want resp_final", response, err)
		}
	}
}

func BenchmarkDecodeResponseWebsocketEvent(b *testing.B) {
	data := []byte(`{"type":"response.output_text.delta","stream_id":"lane","delta":"x","output_index":0,"content_index":0,"sequence_number":1}`)
	event, err := decodeResponseWebsocketEvent(data)
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	for b.Loop() {
		if _, err := decodeResponseWebsocketEvent(data); err != nil {
			b.Fatal(err)
		}
	}
	b.ReportMetric(float64(reflect.TypeOf(event).Size()), "event-B")
	b.ReportMetric(float64(reflect.TypeFor[ResponsesServerEventUnion]().Size()), "public-event-B")
}

func TestResponseWebsocketDecodeSelectsVariant(t *testing.T) {
	for _, test := range []struct {
		name      string
		data      string
		eventType string
		field     string
	}{
		{name: "completed", data: `{"type":"response.completed","response":{"id":"resp_only"},"sequence_number":1}`, eventType: "response.completed", field: "OfResponsesServerEventResponseWsCompleted"},
		{name: "delta", data: `{"type":"response.output_text.delta","delta":"hello","sequence_number":1}`, eventType: "response.output_text.delta", field: "OfResponsesServerEventResponseTextWsDelta"},
		{name: "invalid nested field", data: `{"type":"response.completed","response":42}`, eventType: "response.completed", field: "OfResponsesServerEventResponseWsCompleted"},
		{name: "unknown", data: `{"type":"response.future","response":{"id":"resp_unknown"},"delta":"not a known delta"}`, eventType: "response.future"},
	} {
		t.Run(test.name, func(t *testing.T) {
			frame, err := decodeResponseWebsocketEvent([]byte(test.data))
			if err != nil {
				t.Fatal(err)
			}
			event, err := frame.decode()
			if err != nil {
				t.Fatal(err)
			}
			if event.Type != test.eventType || event.RawJSON() != test.data || !event.JSON.Type.Valid() {
				t.Fatalf("decode(%s) = type %q, raw %q, valid type %t", test.data, event.Type, event.RawJSON(), event.JSON.Type.Valid())
			}
			value := reflect.ValueOf(event)
			metadata := value.FieldByName("JSON")
			for i := 0; i < value.NumField(); i++ {
				field := value.Type().Field(i).Name
				if !strings.HasPrefix(field, "Of") {
					continue
				}
				present := !value.Field(i).IsZero()
				want := field == test.field
				if present != want {
					t.Errorf("decode(%s).%s present = %t, want %t", test.data, field, present, want)
				}
				meta := metadata.FieldByName(field).Interface().(respjson.Field)
				if meta.Valid() != want {
					t.Errorf("decode(%s).JSON.%s.Valid() = %t, want %t", test.data, field, meta.Valid(), want)
				}
				if want && meta.Raw() != test.data {
					t.Errorf("decode(%s).JSON.%s.Raw() = %q, want original event", test.data, field, meta.Raw())
				}
			}
			if test.name == "completed" && event.OfResponsesServerEventResponseWsCompleted.Response.ID != "resp_only" {
				t.Errorf("completed response ID = %q, want resp_only", event.OfResponsesServerEventResponseWsCompleted.Response.ID)
			}
			if test.name == "delta" && event.OfResponsesServerEventResponseTextWsDelta.Delta != "hello" {
				t.Errorf("delta text = %q, want hello", event.OfResponsesServerEventResponseTextWsDelta.Delta)
			}
			if test.name == "invalid nested field" {
				meta := event.OfResponsesServerEventResponseWsCompleted.JSON.Response
				if meta.Valid() || meta.Raw() != "42" {
					t.Errorf("invalid response metadata = valid %t, raw %q; want false, 42", meta.Valid(), meta.Raw())
				}
			}
			if test.name == "unknown" && event.AsAny() != nil {
				t.Errorf("unknown AsAny() = %T, want nil", event.AsAny())
			}
		})
	}
}

func BenchmarkResponseWebsocketConsumedEvent(b *testing.B) {
	for _, test := range []struct{ name, data string }{
		{name: "delta", data: `{"type":"response.output_text.delta","delta":"hello","sequence_number":1}`},
		{name: "terminal-1MiB", data: `{"type":"response.completed","response":{"id":"resp_large","metadata":{"text":"` + strings.Repeat("x", 1<<20) + `"}}}`},
	} {
		b.Run(test.name, func(b *testing.B) {
			frame, err := decodeResponseWebsocketEvent([]byte(test.data))
			if err != nil {
				b.Fatal(err)
			}
			b.ReportAllocs()
			for b.Loop() {
				if _, err := frame.decode(); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

package responses_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/openai/openai-go/v3/responses"
)

func TestInputMessageUnionRoundTrip(t *testing.T) {
	for _, test := range []struct {
		name string
		raw  string
	}{
		{"plain input", `{"type":"message","role":"user","content":"hi"}`},
		{"complete output", `{"id":"msg_1","type":"message","role":"assistant","status":"completed","content":[{"type":"output_text","text":"hello","annotations":[]}]}`},
		{"message without id", `{"type":"message","role":"assistant","status":"completed","content":[{"type":"output_text","text":"hello","annotations":[]}]}`},
	} {
		t.Run(test.name, func(t *testing.T) {
			var message responses.ResponseInputItemUnionParam
			if err := json.Unmarshal([]byte(test.raw), &message); err != nil {
				t.Fatal(err)
			}
			switch test.name {
			case "plain input":
				if message.OfMessage == nil || message.OfInputMessage != nil || message.OfOutputMessage != nil {
					t.Fatal("expected only OfMessage to be selected")
				}
			case "complete output":
				if message.OfOutputMessage == nil || message.OfMessage != nil || message.OfInputMessage != nil {
					t.Fatal("expected only OfOutputMessage to be selected")
				}
			case "message without id":
				if message.OfInputMessage == nil || message.OfMessage != nil || message.OfOutputMessage != nil {
					t.Fatal("expected only OfInputMessage to be selected")
				}
			}
			assertMessageUnionJSON(t, message, test.raw)
		})
	}
}

func TestInputMessageUnionReuse(t *testing.T) {
	var message responses.ResponseInputItemUnionParam
	output := `{"id":"msg_1","type":"message","role":"assistant","status":"completed","content":[{"type":"output_text","text":"hello","annotations":[]}]}`
	if err := json.Unmarshal([]byte(output), &message); err != nil {
		t.Fatal(err)
	}
	if message.OfOutputMessage == nil {
		t.Fatal("expected OfOutputMessage before reusing the union")
	}
	assertMessageUnionJSON(t, message, output)

	input := `{"type":"message","role":"user","content":"next message"}`
	if err := json.Unmarshal([]byte(input), &message); err != nil {
		t.Fatal(err)
	}
	if message.OfMessage == nil || message.OfInputMessage != nil || message.OfOutputMessage != nil {
		t.Fatal("expected only OfMessage after reusing the union")
	}
	assertMessageUnionJSON(t, message, input)
}

func assertMessageUnionJSON(t *testing.T, message responses.ResponseInputItemUnionParam, want string) {
	t.Helper()
	data, err := json.Marshal(message)
	if err != nil {
		t.Fatal(err)
	}
	var gotObject, wantObject any
	if err := json.Unmarshal(data, &gotObject); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal([]byte(want), &wantObject); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(gotObject, wantObject) {
		t.Fatalf("round trip lost message fields: got %s, want %s", data, want)
	}
}

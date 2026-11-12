package live_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	wire "github.com/coder/websocket"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/live"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/param"
)

func TestDecisionItemTypedSendPreservesWireItems(t *testing.T) {
	items := map[string]string{
		"user":                       `{"type":"message","role":"user","content":"hello"}`,
		"system":                     `{"type":"message","role":"system","content":"hello"}`,
		"developer":                  `{"type":"message","role":"developer","content":"hello"}`,
		"assistant":                  `{"type":"message","role":"assistant","content":"hello"}`,
		"implicit_message":           `{"role":"user","content":"hello"}`,
		"future_item":                `{"type":"future_item","payload":{"keep":true}}`,
		"json_content":               `{"type":"message","role":"user","content":[{"type":"input_json","json":{"ready":true}}]}`,
		"json_content_object":        `{"type":"message","role":"user","content":[{"type":"input_json","data":{"nested":{"id":9007199254740993}}}]}`,
		"json_content_array":         `{"type":"message","role":"user","content":[{"type":"input_json","data":[{"ids":[9007199254740993]}]}]}`,
		"function_call":              `{"type":"function_call","call_id":"call_test","name":"weather","arguments":"{}","future":{"keep":true}}`,
		"function_call_output":       `{"type":"function_call_output","call_id":"call_test","output":"done"}`,
		"custom_tool_call":           `{"type":"custom_tool_call","call_id":"call_test","name":"weather","input":"hello"}`,
		"custom_tool_output":         `{"type":"custom_tool_call_output","call_id":"call_test","output":"done"}`,
		"function_call_output_parts": `{"type":"function_call_output","call_id":"call_test","output":[{"type":"input_text","text":"done","future":true}]}`,
		"custom_tool_output_parts":   `{"type":"custom_tool_call_output","call_id":"call_test","output":[{"type":"input_text","text":"done","future":{"keep":true}}]}`,
	}
	for _, family := range []string{"primary", "sideband", "fork"} {
		for name, item := range items {
			t.Run(family+"/"+name, func(t *testing.T) {
				ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
				defer cancel()
				type received struct {
					data []byte
					err  error
				}
				frames := make(chan received, 1)
				server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					socket, err := wire.Accept(w, r, nil)
					if err != nil {
						frames <- received{err: err}
						return
					}
					defer func() { _ = socket.CloseNow() }()
					_, data, err := socket.Read(ctx)
					frames <- received{data, err}
				}))
				defer server.Close()
				client := openai.NewClient(option.WithAPIKey("synthetic-test-key"), option.WithBaseURL(server.URL), option.WithHTTPClient(server.Client()))
				payload := []byte(`{"type":"decision.item.create","item":` + item + `}`)
				if err := sendDecisionWireItem(ctx, client, family, payload); err != nil {
					t.Fatal(err)
				}
				select {
				case got := <-frames:
					if got.err != nil {
						t.Fatal(got.err)
					}
					assertDecisionWireJSON(t, got.data, payload)
				case <-ctx.Done():
					t.Fatal(ctx.Err())
				}
			})
		}
	}
}

func sendDecisionWireItem(ctx context.Context, client openai.Client, family string, payload []byte) error {
	switch family {
	case "primary":
		connection, err := client.Live.Connect(ctx, live.LiveConnectionOptions{})
		if err != nil {
			return err
		}
		defer connection.Abort()
		var event live.ClientEventUnionParam
		if err := json.Unmarshal(payload, &event); err != nil {
			return err
		}
		return connection.Send(ctx, event)
	case "sideband":
		connection, err := client.Live.Sideband.Connect(ctx, "session_test", live.SidebandConnectionOptions{})
		if err != nil {
			return err
		}
		defer connection.Abort()
		var event live.ConnectClientEventUnionParam
		if err := json.Unmarshal(payload, &event); err != nil {
			return err
		}
		return connection.Send(ctx, event)
	default:
		connection, err := client.Live.Forks.Connect(ctx, "session_test", live.ForkConnectionOptions{})
		if err != nil {
			return err
		}
		defer connection.Abort()
		var event live.ForkClientEventUnionParam
		if err := json.Unmarshal(payload, &event); err != nil {
			return err
		}
		return connection.Send(ctx, event)
	}
}

func TestDecisionItemDecodedKnownFieldsRemainMutable(t *testing.T) {
	var item live.DecisionItemCreateEventItemUnionParam
	if err := json.Unmarshal([]byte(`{"type":"function_call","call_id":"call_test","name":"weather","arguments":"{}","future":{"keep":true}}`), &item); err != nil {
		t.Fatal(err)
	}
	if item.OfFunctionCall == nil {
		t.Fatal("function call variant was not selected")
	}
	item.OfFunctionCall.ID = openai.String("item_changed")
	item.OfFunctionCall.Arguments = `{"city":"Boston"}`
	data, err := json.Marshal(item)
	if err != nil {
		t.Fatal(err)
	}
	assertDecisionWireJSON(t, data, []byte(`{"type":"function_call","id":"item_changed","call_id":"call_test","name":"weather","arguments":"{\"city\":\"Boston\"}","future":{"keep":true}}`))
}

func TestDecisionItemManualToolFieldsDoNotAddRole(t *testing.T) {
	item := live.DecisionItemCreateEventItemUnionParam{
		OfFunctionCall: &live.DecisionItemCreateEventItemFunctionCallParam{
			CallID: "call_test", Name: "weather", Arguments: "{}",
		},
	}
	data, err := json.Marshal(item)
	if err != nil {
		t.Fatal(err)
	}
	assertDecisionWireJSON(t, data, []byte(`{"type":"function_call","call_id":"call_test","name":"weather","arguments":"{}"}`))
}

func TestDecisionItemExplicitOverridesKeepPrecedence(t *testing.T) {
	var item live.DecisionItemCreateEventItemFunctionCallParam
	if err := json.Unmarshal([]byte(`{"type":"function_call","call_id":"call_test","name":"weather","arguments":"{}","future":true}`), &item); err != nil {
		t.Fatal(err)
	}
	item.SetExtraFields(map[string]any{"name": "forecast", "future": param.Omit})
	data, err := json.Marshal(item)
	if err != nil {
		t.Fatal(err)
	}
	assertDecisionWireJSON(t, data, []byte(`{"type":"function_call","call_id":"call_test","name":"forecast","arguments":"{}"}`))

	raw := json.RawMessage(`{"type":"future_item","role":null,"keep":true}`)
	union := param.Override[live.DecisionItemCreateEventItemUnionParam](raw)
	data, err = json.Marshal(union)
	if err != nil {
		t.Fatal(err)
	}
	assertDecisionWireJSON(t, data, raw)

	data, err = json.Marshal(param.NullStruct[live.DecisionItemCreateEventItemUnionParam]())
	if err != nil || string(data) != "null" {
		t.Fatalf("null override changed: %s, %v", data, err)
	}
}

func TestDecisionItemReusedMetadata(t *testing.T) {
	type itemParam = live.DecisionUserMessageParam
	tests := []struct {
		name         string
		seed         func(map[string]any) itemParam
		want         string
		wantIncoming string
		wantNull     bool
		wantOverride bool
	}{
		{
			name: "null-struct",
			seed: func(map[string]any) itemParam { return param.NullStruct[itemParam]() },
			want: "null", wantIncoming: "null", wantNull: true, wantOverride: true,
		},
		{
			name: "raw-override",
			seed: func(map[string]any) itemParam {
				return param.Override[itemParam](json.RawMessage(`{"type":"future_item","kept":true}`))
			},
			want:         `{"type":"future_item","kept":true}`,
			wantIncoming: `{"type":"future_item","kept":true}`,
			wantOverride: true,
		},
		{
			name: "extra-fields",
			seed: func(extras map[string]any) itemParam {
				p := itemParam{Type: "message", ID: openai.String("original-id")}
				p.SetExtraFields(extras)
				return p
			},
			want:         `{"type":"message","id":"extra-id","role":"user","content":"edited","extension":"retained"}`,
			wantIncoming: `{"type":"message","id":"extra-id","role":"user","content":"edited","extension":"retained","incoming":true}`,
		},
	}
	inputs := []struct {
		name       string
		wire       string
		incoming   bool
		knownIDSet bool
	}{
		{name: "null", wire: `null`},
		{name: "empty-object", wire: `{}`},
		{name: "known-id", wire: `{"id":"new-id"}`, knownIDSet: true},
		{name: "unknown-and-collision", wire: `{"extension":"replace","incoming":true,"omit":"must-hide"}`, incoming: true},
	}
	for _, tc := range tests {
		for _, input := range inputs {
			t.Run(tc.name+"/"+input.name, func(t *testing.T) {
				// Callers can keep their extras map or a copy of a parameter.
				// Decoding one parameter must not change either of those values.
				originalExtras := map[string]any{"id": "extra-id", "extension": "retained", "omit": param.Omit}
				item := tc.seed(originalExtras)
				original := item
				originalWire, err := json.Marshal(original)
				if err != nil {
					t.Fatal(err)
				}
				if decodeErr := json.Unmarshal([]byte(input.wire), &item); decodeErr != nil {
					t.Fatal(decodeErr)
				}
				if input.knownIDSet && (!item.ID.Valid() || item.ID.Value != "new-id") {
					t.Fatalf("typed ID did not decode: %#v", item.ID)
				}
				item.Content.OfString = openai.String("edited")
				data, err := json.Marshal(item)
				if err != nil {
					t.Fatal(err)
				}
				want := tc.want
				if input.incoming {
					want = tc.wantIncoming
				}
				assertDecisionWireJSON(t, data, []byte(want))
				if got := param.IsNull(item); got != tc.wantNull {
					t.Fatalf("IsNull = %v; want %v", got, tc.wantNull)
				}
				if _, got := item.Overrides(); got != tc.wantOverride {
					t.Fatalf("Overrides = %v; want %v", got, tc.wantOverride)
				}
				if !reflect.DeepEqual(originalExtras, map[string]any{"id": "extra-id", "extension": "retained", "omit": param.Omit}) {
					t.Fatalf("changed caller's extras map: %#v", originalExtras)
				}
				afterWire, err := json.Marshal(original)
				if err != nil {
					t.Fatal(err)
				}
				assertDecisionWireJSON(t, afterWire, originalWire)
			})
		}
	}
}

func TestDecisionItemEmptyDecodePreservesOmission(t *testing.T) {
	for _, input := range []string{`null`, `{}`} {
		t.Run(input, func(t *testing.T) {
			var item live.DecisionItemCreateEventItemUnionParam
			if err := json.Unmarshal([]byte(input), &item); err != nil {
				t.Fatal(err)
			}
			if !param.IsOmitted(item) {
				t.Fatalf("empty item after %s is no longer omitted", input)
			}
		})
	}
}

func TestDecisionItemNilReceiverPreservesError(t *testing.T) {
	var item *live.DecisionItemCreateEventItemUnionParam
	if err := item.UnmarshalJSON([]byte(`null`)); err == nil {
		t.Fatal("nil receiver must return an error")
	}
}

func TestDecisionItemRepeatedDecodeKeepsLatestToolCall(t *testing.T) {
	for _, override := range []string{"none", "before-decode", "after-decode"} {
		t.Run(override, func(t *testing.T) {
			var item live.DecisionItemCreateEventItemFunctionCallParam
			// RawMessage is valid caller data too. It must never be mistaken
			// for an automatically decoded field that later inputs can replace.
			caller := map[string]any{"call_id": json.RawMessage(`"caller"`), "name": param.Omit}
			if override == "before-decode" {
				item.SetExtraFields(caller)
			}
			if err := json.Unmarshal([]byte(`{"type":"function_call","call_id":"call1","name":"first","arguments":"{\"first\":1}","keep":"first-only"}`), &item); err != nil {
				t.Fatal(err)
			}
			if override == "after-decode" {
				item.SetExtraFields(caller)
			}
			original := item
			originalWire, err := json.Marshal(original)
			if err != nil {
				t.Fatal(err)
			}

			if decodeErr := json.Unmarshal([]byte(`{"id":"item2","call_id":"call2","name":"second","arguments":"{\"second\":2}"}`), &item); decodeErr != nil {
				t.Fatal(decodeErr)
			}
			// Repeated decoding updates wire fields and retains caller extras.
			got, err := json.Marshal(item)
			if err != nil {
				t.Fatal(err)
			}
			want := `{"id":"item2","type":"function_call","call_id":"call2","name":"second","arguments":"{\"second\":2}","keep":"first-only"}`
			if override != "none" {
				want = `{"id":"item2","type":"function_call","call_id":"caller","arguments":"{\"second\":2}","keep":"first-only"}`
			}
			assertDecisionWireJSON(t, got, []byte(want))

			// ExtraFields exposes caller fields as it did in the released SDK.
			// It must not conflate those fields with decoded wire properties.
			if override == "none" {
				if item.ExtraFields() != nil {
					t.Fatalf("decoded fields appeared in caller extras: %#v", item.ExtraFields())
				}
			} else if !reflect.DeepEqual(item.ExtraFields(), caller) {
				t.Fatalf("caller extras changed: got %#v; want %#v", item.ExtraFields(), caller)
			}
			if !reflect.DeepEqual(caller, map[string]any{"call_id": json.RawMessage(`"caller"`), "name": param.Omit}) {
				t.Fatalf("changed caller's map: %#v", caller)
			}
			afterWire, err := json.Marshal(original)
			if err != nil {
				t.Fatal(err)
			}
			assertDecisionWireJSON(t, afterWire, originalWire)
		})
	}
}

func TestDecisionJSONContentEditsPreserveNumbers(t *testing.T) {
	t.Run("message object", func(t *testing.T) {
		var part live.DecisionUserMessageContentPartInputJSONParam
		if err := json.Unmarshal([]byte(`{"type":"input_json","data":{"id":9007199254740993,"label":"before"}}`), &part); err != nil {
			t.Fatal(err)
		}
		if part.Data.OfLiveDecisionJSONObject == nil {
			t.Fatal("object data variant was not selected")
		}
		(*part.Data.OfLiveDecisionJSONObject)["label"] = "after"
		(*part.Data.OfLiveDecisionJSONObject)["nullable"] = param.NullMap[live.DecisionJSONObject]()
		wire, err := json.Marshal(part)
		if err != nil {
			t.Fatal(err)
		}
		assertDecisionWireJSON(t, wire, []byte(`{"type":"input_json","data":{"id":9007199254740993,"label":"after","nullable":null}}`))
	})
	t.Run("instruction array", func(t *testing.T) {
		var part live.DecisionQuestionInstructionPartInputJSONParam
		if err := json.Unmarshal([]byte(`{"type":"input_json","data":[{"id":9007199254740993,"label":"before"}]}`), &part); err != nil {
			t.Fatal(err)
		}
		if len(part.Data.OfArrayOfJSONObjects) != 1 {
			t.Fatal("array data variant was not selected")
		}
		part.Data.OfArrayOfJSONObjects[0]["label"] = "after"
		wire, err := json.Marshal(part)
		if err != nil {
			t.Fatal(err)
		}
		assertDecisionWireJSON(t, wire, []byte(`{"type":"input_json","data":[{"id":9007199254740993,"label":"after"}]}`))
	})
}

func assertDecisionWireJSON(t *testing.T, got, want []byte) {
	t.Helper()
	var gotValue, wantValue any
	gotDecoder := json.NewDecoder(bytes.NewReader(got))
	gotDecoder.UseNumber()
	if err := gotDecoder.Decode(&gotValue); err != nil {
		t.Fatal(err)
	}
	wantDecoder := json.NewDecoder(bytes.NewReader(want))
	wantDecoder.UseNumber()
	if err := wantDecoder.Decode(&wantValue); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(gotValue, wantValue) {
		t.Fatalf("wire item changed: got %s; want %s", got, want)
	}
}

func TestDecisionItemRoleDefaultsOnlyForMessages(t *testing.T) {
	cases := []struct {
		name string
		raw  string
		want string
	}{
		{"future_decoded", `{"type":"future_item","payload":{"keep":true}}`, `{"type":"future_item","payload":{"keep":true}}`},
		{"absent_type", `{"content":"hello"}`, `{"content":"hello","role":"user"}`},
		{"empty_type", `{"type":"","content":"hello"}`, `{"content":"hello","role":"user"}`},
		{"explicit_message", `{"type":"message","content":"hello"}`, `{"type":"message","content":"hello","role":"user"}`},
		{"future_explicit_role", `{"type":"future_item","role":"assistant","payload":1}`, `{"type":"future_item","role":"assistant","payload":1}`},
		{"known_tool", `{"type":"function_call","call_id":"call_test","name":"f","arguments":"{}"}`, `{"type":"function_call","call_id":"call_test","name":"f","arguments":"{}"}`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var item live.DecisionItemCreateEventItemUnionParam
			if err := json.Unmarshal([]byte(tc.raw), &item); err != nil {
				t.Fatal(err)
			}
			data, err := json.Marshal(item)
			if err != nil {
				t.Fatal(err)
			}
			assertDecisionWireJSON(t, data, []byte(tc.want))
		})
	}
}

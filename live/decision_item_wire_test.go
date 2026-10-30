package live_test

import (
	"context"
	"encoding/json"
	"errors"
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
		"user":                 `{"type":"message","role":"user","content":"hello"}`,
		"system":               `{"type":"message","role":"system","content":"hello"}`,
		"developer":            `{"type":"message","role":"developer","content":"hello"}`,
		"assistant":            `{"type":"message","role":"assistant","content":"hello"}`,
		"implicit_message":     `{"role":"user","content":"hello"}`,
		"json_content":         `{"type":"message","role":"user","content":[{"type":"input_json","json":{"ready":true}}]}`,
		"function_call":        `{"type":"function_call","call_id":"call_test","name":"weather","arguments":"{}","future":{"keep":true}}`,
		"function_call_output": `{"type":"function_call_output","call_id":"call_test","output":"done"}`,
		"custom_tool_call":     `{"type":"custom_tool_call","call_id":"call_test","name":"weather","input":"hello"}`,
		"custom_tool_output":   `{"type":"custom_tool_call_output","call_id":"call_test","output":"done"}`,
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
	var item live.DecisionItemCreateEventItemParam
	if err := json.Unmarshal([]byte(`{"type":"function_call","call_id":"call_test","name":"weather","arguments":"{}"}`), &item); err != nil {
		t.Fatal(err)
	}
	item.ID = openai.String("item_changed")
	item.Role = "assistant"
	data, err := json.Marshal(item)
	if err != nil {
		t.Fatal(err)
	}
	assertDecisionWireJSON(t, data, []byte(`{"type":"function_call","id":"item_changed","role":"assistant","call_id":"call_test","name":"weather","arguments":"{}"}`))
	item.Type = "message"
	item.Content.OfString = openai.String("changed")
	data, err = json.Marshal(item)
	if err != nil {
		t.Fatal(err)
	}
	assertDecisionWireJSON(t, data, []byte(`{"type":"message","id":"item_changed","role":"assistant","content":"changed","call_id":"call_test","name":"weather","arguments":"{}"}`))
}

func TestDecisionItemManualToolFieldsDoNotAddRole(t *testing.T) {
	item := live.DecisionItemCreateEventItemParam{Type: "function_call"}
	item.SetExtraFields(map[string]any{"call_id": "call_test", "name": "weather", "arguments": "{}"})
	data, err := json.Marshal(item)
	if err != nil {
		t.Fatal(err)
	}
	assertDecisionWireJSON(t, data, []byte(`{"type":"function_call","call_id":"call_test","name":"weather","arguments":"{}"}`))
}

func TestDecisionItemExplicitOverridesKeepPrecedence(t *testing.T) {
	item := live.DecisionItemCreateEventItemParam{Type: "function_call"}
	item.SetExtraFields(map[string]any{"role": "system", "call_id": "call_test", "name": "weather", "arguments": "{}"})
	data, err := json.Marshal(item)
	if err != nil {
		t.Fatal(err)
	}
	assertDecisionWireJSON(t, data, []byte(`{"type":"function_call","role":"system","call_id":"call_test","name":"weather","arguments":"{}"}`))
	item.Role = "assistant"
	item.SetExtraFields(map[string]any{"role": param.Omit})
	data, err = json.Marshal(item)
	if err != nil {
		t.Fatal(err)
	}
	assertDecisionWireJSON(t, data, []byte(`{"type":"function_call"}`))

	raw := json.RawMessage(`{"type":"future_item","role":null,"keep":true}`)
	item = param.Override[live.DecisionItemCreateEventItemParam](raw)
	data, err = json.Marshal(item)
	if err != nil {
		t.Fatal(err)
	}
	assertDecisionWireJSON(t, data, raw)

	item = param.NullStruct[live.DecisionItemCreateEventItemParam]()
	data, err = json.Marshal(item)
	if err != nil || string(data) != "null" {
		t.Fatalf("null override changed: %s, %v", data, err)
	}
}

func TestDecisionItemReusedMetadata(t *testing.T) {
	type itemParam = live.DecisionItemCreateEventItemParam
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
			var item live.DecisionItemCreateEventItemParam
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
	var item *live.DecisionItemCreateEventItemParam
	if err := item.UnmarshalJSON([]byte(`null`)); err == nil {
		t.Fatal("nil receiver must return an error")
	}
}

func TestDecisionItemRepeatedDecodeKeepsLatestToolCall(t *testing.T) {
	for _, override := range []string{"none", "before-decode", "after-decode"} {
		t.Run(override, func(t *testing.T) {
			var item live.DecisionItemCreateEventItemParam
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
			// Typed fields remain mutable while repeated decoding updates extras.
			item.Role = "assistant"
			got, err := json.Marshal(item)
			if err != nil {
				t.Fatal(err)
			}
			want := `{"id":"item2","type":"function_call","role":"assistant","call_id":"call2","name":"second","arguments":"{\"second\":2}","keep":"first-only"}`
			if override != "none" {
				want = `{"id":"item2","type":"function_call","role":"assistant","call_id":"caller","arguments":"{\"second\":2}","keep":"first-only"}`
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

func assertDecisionWireJSON(t *testing.T, got, want []byte) {
	t.Helper()
	var gotValue, wantValue any
	if err := json.Unmarshal(got, &gotValue); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(want, &wantValue); err != nil {
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
			var item live.DecisionItemCreateEventItemParam
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
	t.Run("future_manual_and_caller_overrides", func(t *testing.T) {
		item := live.DecisionItemCreateEventItemParam{Type: "future_item"}
		item.SetExtraFields(map[string]any{"payload": true})
		data, err := json.Marshal(item)
		if err != nil {
			t.Fatal(err)
		}
		assertDecisionWireJSON(t, data, []byte(`{"type":"future_item","payload":true}`))
		item.SetExtraFields(map[string]any{"payload": true, "role": "developer"})
		data, err = json.Marshal(item)
		if err != nil {
			t.Fatal(err)
		}
		assertDecisionWireJSON(t, data, []byte(`{"type":"future_item","payload":true,"role":"developer"}`))
		item.SetExtraFields(map[string]any{"payload": true, "role": nil})
		data, err = json.Marshal(item)
		if err != nil {
			t.Fatal(err)
		}
		assertDecisionWireJSON(t, data, []byte(`{"type":"future_item","payload":true,"role":null}`))
		item.Role = "assistant"
		item.SetExtraFields(map[string]any{"payload": true, "role": param.Omit})
		data, err = json.Marshal(item)
		if err != nil {
			t.Fatal(err)
		}
		assertDecisionWireJSON(t, data, []byte(`{"type":"future_item","payload":true}`))
	})
}

type decisionRawType string

type decisionTypeMarshaler struct {
	calls *int
	data  json.RawMessage
	err   error
}

func (v decisionTypeMarshaler) MarshalJSON() ([]byte, error) {
	*v.calls++
	return v.data, v.err
}

func TestDecisionItemRawTypeOverridePreservesReleasedDefaults(t *testing.T) {
	for _, tc := range []struct {
		name   string
		typed  string
		extras map[string]any
		want   string
	}{
		{"future_to_message", "future_item", map[string]any{"type": "message", "content": "hello"}, `{"type":"message","content":"hello","role":"user"}`},
		{"tool_to_message", "function_call", map[string]any{"type": "message", "content": "hello"}, `{"type":"message","content":"hello","role":"user"}`},
		{"message_to_tool_raw_defaults", "message", map[string]any{"type": "function_call", "call_id": "c", "name": "f", "arguments": "{}"}, `{"type":"function_call","call_id":"c","name":"f","arguments":"{}","role":"user"}`},
		{"extra_only_tool_omit_role", "", map[string]any{"type": "function_call", "call_id": "c", "name": "f", "arguments": "{}", "role": param.Omit}, `{"type":"function_call","call_id":"c","name":"f","arguments":"{}"}`},
		{"future_to_raw_message", "future_item", map[string]any{"type": json.RawMessage(`"message"`)}, `{"type":"message","role":"user"}`},
		{"future_to_named_message", "future_item", map[string]any{"type": decisionRawType("message")}, `{"type":"message","role":"user"}`},
		{"future_to_null_type", "future_item", map[string]any{"type": nil}, `{"type":null,"role":"user"}`},
		{"future_to_omitted_type", "future_item", map[string]any{"type": param.Omit}, `{"role":"user"}`},
		{"future_to_empty_type", "future_item", map[string]any{"type": ""}, `{"type":"","role":"user"}`},
		{"future_to_tool_raw_defaults", "future_item", map[string]any{"type": "custom_tool_call", "call_id": "c", "name": "f", "input": "hello"}, `{"type":"custom_tool_call","call_id":"c","name":"f","input":"hello","role":"user"}`},
		{"future_to_message_explicit_role", "future_item", map[string]any{"type": "message", "role": "developer"}, `{"type":"message","role":"developer"}`},
		{"future_to_message_explicit_null_role", "future_item", map[string]any{"type": "message", "role": nil}, `{"type":"message","role":null}`},
		{"future_to_message_explicit_omitted_role", "future_item", map[string]any{"type": "message", "role": param.Omit}, `{"type":"message"}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			item := live.DecisionItemCreateEventItemParam{Type: tc.typed}
			item.SetExtraFields(tc.extras)
			data, err := json.Marshal(item)
			if err != nil {
				t.Fatal(err)
			}
			assertDecisionWireJSON(t, data, []byte(tc.want))
		})
	}
}

func TestDecisionItemRawTypeMarshalOnce(t *testing.T) {
	for _, tc := range []struct {
		name string
		data json.RawMessage
		err  error
		want string
	}{
		{"message", json.RawMessage(`"message"`), nil, `{"type":"message","role":"user"}`},
		{"custom_error", nil, errors.New("synthetic type error"), ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			item := live.DecisionItemCreateEventItemParam{Type: "future_item"}
			item.SetExtraFields(map[string]any{"type": decisionTypeMarshaler{calls: &calls, data: tc.data, err: tc.err}})
			data, err := json.Marshal(item)
			if !errors.Is(err, tc.err) {
				t.Fatalf("marshal error = %v; want %v", err, tc.err)
			}
			if calls != 1 {
				t.Fatalf("type marshaler called %d times; want 1", calls)
			}
			if err == nil {
				assertDecisionWireJSON(t, data, []byte(tc.want))
			}
		})
	}
}

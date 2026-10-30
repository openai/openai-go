package live_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	wire "github.com/coder/websocket"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/live"
	"github.com/openai/openai-go/v3/option"
)

func TestDecisionItemTypedReceivePreservesStoredWireItems(t *testing.T) {
	items := []struct {
		name   string
		fields string
		kind   string
		role   string
	}{
		{"user", `"type":"message","role":"user","content":"hello"`, "message", "user"},
		{"system", `"type":"message","role":"system","content":"hello"`, "message", "system"},
		{"developer", `"type":"message","role":"developer","content":"hello"`, "message", "developer"},
		{"assistant", `"type":"message","role":"assistant","content":"hello"`, "message", "assistant"},
		{"implicit_message", `"role":"user","content":"hello"`, "", "user"},
		{"json_content", `"type":"message","role":"user","content":[{"type":"input_json","json":{"ready":true,"data":[1,2,3]}}]`, "message", "user"},
		{"function_call", `"type":"function_call","call_id":"call_test","name":"weather","arguments":"{\"city\":\"Boston\"}"`, "function_call", ""},
		{"function_call_output", `"type":"function_call_output","call_id":"call_test","output":"done"`, "function_call_output", ""},
		{"custom_tool_call", `"type":"custom_tool_call","call_id":"call_test","name":"weather","input":"hello"`, "custom_tool_call", ""},
		{"custom_tool_call_output", `"type":"custom_tool_call_output","call_id":"call_test","output":"done"`, "custom_tool_call_output", ""},
	}

	for _, family := range []string{"primary", "sideband", "fork"} {
		t.Run(family, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
			defer cancel()
			// Keep unusual whitespace and unknown fields so each RawJSON check
			// also detects a lossy rebuild of either the item or its envelope.
			itemJSON := make([]string, len(items))
			frames := make([]string, len(items))
			for i, item := range items {
				itemJSON[i] = fmt.Sprintf(`{ "id" : %q, %s, "future_item" : {"keep":[true,{"n":2}]} }`, "item_"+item.name, item.fields)
				frames[i] = fmt.Sprintf(`{ "type" : "decision.item.created", "event_id" : %q, "client_event_id": %q, "item" : %s, "future_event" : {"keep":"yes"} }`,
					"evt_"+item.name, "client_"+item.name, itemJSON[i])
			}
			server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				socket, err := wire.Accept(w, r, nil)
				if err != nil {
					t.Errorf("accept %s: %v", family, err)
					return
				}
				defer func() { _ = socket.CloseNow() }()
				for _, frame := range frames {
					if err := socket.Write(ctx, wire.MessageText, []byte(frame)); err != nil {
						t.Errorf("write %s: %v", family, err)
						return
					}
				}
				_, _, _ = socket.Read(ctx)
			}))
			defer server.Close()
			client := openai.NewClient(option.WithAPIKey("synthetic-test-key"), option.WithBaseURL(server.URL), option.WithHTTPClient(server.Client()))

			var recv func() (string, string, any, error)
			switch family {
			case "primary":
				connection, err := client.Live.Connect(ctx, live.LiveConnectionOptions{})
				if err != nil {
					t.Fatal(err)
				}
				defer connection.Abort()
				recv = func() (string, string, any, error) {
					event, err := connection.Recv(ctx)
					return event.Type, event.RawJSON(), event.AsAny(), err
				}
			case "sideband":
				connection, err := client.Live.Sideband.Connect(ctx, "session_test", live.SidebandConnectionOptions{})
				if err != nil {
					t.Fatal(err)
				}
				defer connection.Abort()
				recv = func() (string, string, any, error) {
					event, err := connection.Recv(ctx)
					return event.Type, event.RawJSON(), event.AsAny(), err
				}
			case "fork":
				connection, err := client.Live.Forks.Connect(ctx, "session_test", live.ForkConnectionOptions{})
				if err != nil {
					t.Fatal(err)
				}
				defer connection.Abort()
				recv = func() (string, string, any, error) {
					event, err := connection.Recv(ctx)
					return event.Type, event.RawJSON(), event.AsAny(), err
				}
			}

			for i, item := range items {
				t.Run(item.name, func(t *testing.T) {
					kind, raw, variant, err := recv()
					if err != nil {
						t.Fatal(err)
					}
					if kind != "decision.item.created" || raw != frames[i] {
						t.Fatalf("received event type %q, raw %q; want decision.item.created, %q", kind, raw, frames[i])
					}
					created, ok := variant.(live.DecisionItemCreatedEvent)
					if !ok {
						t.Fatalf("received typed variant %T; want live.DecisionItemCreatedEvent", variant)
					}
					if string(created.Type) != "decision.item.created" || created.EventID != "evt_"+item.name ||
						created.ClientEventID != "client_"+item.name || created.RawJSON() != frames[i] {
						t.Fatalf("created envelope fields or raw JSON changed: %+v, raw %q", created, created.RawJSON())
					}
					if created.Item.ID != "item_"+item.name || (item.kind != "" && created.Item.Type != item.kind) ||
						(item.role != "" && string(created.Item.Role) != item.role) || created.Item.RawJSON() != itemJSON[i] {
						t.Fatalf("created item fields or raw JSON changed: %+v, raw %q; want id %q, type %q, role %q, raw %q",
							created.Item, created.Item.RawJSON(), "item_"+item.name, item.kind, item.role, itemJSON[i])
					}
					if item.role == "" && (created.Item.JSON.Role.Valid() || created.Item.JSON.Content.Valid()) {
						t.Fatalf("tool item gained message-only role/content fields: %q", created.Item.RawJSON())
					}
					if item.kind == "" && created.Item.JSON.Type.Valid() {
						t.Fatalf("implicit message gained an explicit type field: %q", created.Item.RawJSON())
					}
				})
			}
		})
	}
}

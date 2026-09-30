package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type recordingWallet struct {
	args  []balanceArguments
	cause error
}

func (w *recordingWallet) Balance(_ context.Context, args balanceArguments) (balanceReceipt, error) {
	w.args = append(w.args, args)
	return balanceReceipt{Wallet: "bound-wallet", Asset: args.Asset, Balance: "12.50"}, w.cause
}

func TestBalanceActionRejectsInvalidArguments(t *testing.T) {
	for name, args := range map[string]map[string]any{
		"nil": nil, "missing": {}, "null": {"asset": nil},
		"wrong type": {"asset": 123}, "unsupported asset": {"asset": "BTC"},
		"unexpected field": {"asset": "USDC", "wallet": "other-wallet"},
		"case mismatch":    {"Asset": "USDC"},
	} {
		t.Run(name, func(t *testing.T) {
			wallet := &recordingWallet{}
			_, err := (balanceAction{wallet: wallet}).handle(context.Background(), args)
			if err == nil || len(wallet.args) != 0 {
				t.Fatalf("invalid arguments reached the wallet: calls=%d, error=%v", len(wallet.args), err)
			}
		})
	}
}

func TestBalanceActionStream(t *testing.T) {
	for _, fail := range []bool{false, true} {
		t.Run(fmt.Sprintf("provider_error=%t", fail), func(t *testing.T) {
			wallet := &recordingWallet{}
			if fail {
				wallet.cause = errors.New("private provider detail")
			}
			action := balanceAction{wallet: wallet}
			tool, err := action.tool()
			if err != nil {
				t.Fatal(err)
			}
			definition := tool.OfParamFunction
			schema := definition.Parameters
			properties := schema["properties"].(map[string]any)
			if schema["type"] != "object" || schema["additionalProperties"] != false || len(properties) != 1 ||
				!reflect.DeepEqual(schema["required"], []any{"asset"}) ||
				!reflect.DeepEqual(properties["asset"].(map[string]any)["enum"], []any{"USDC", "ETH"}) {
				t.Fatalf("schema does not describe the typed arguments: %v", schema)
			}

			posts := make(chan map[string]any, 2)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.Method + " " + r.URL.Path {
				case "GET /agents/sessions/session":
					w.Header().Set("Content-Type", "application/json")
					_, _ = fmt.Fprint(w, `{"id":"session","status":"idle"}`)
				case "GET /agents/sessions/session/events":
					w.Header().Set("Content-Type", "text/event-stream")
					for _, event := range []string{
						`{"type":"agent.session.turn.created","event_id":"created","turn_id":"turn","turn":{"id":"turn","subagent_id":null}}`,
						`{"type":"agent.session.turn.item.added","event_id":"call","turn_id":"turn","item":{"id":"item","type":"function_call","call_id":"call","name":"wallet_balance","arguments":{"asset":"USDC"},"status":"in_progress"}}`,
						`{"type":"agent.session.turn.completed","event_id":"done","turn_id":"turn"}`,
						`{"type":"agent.session.idle","event_id":"idle","session":{"status":"idle"}}`,
					} {
						_, _ = fmt.Fprintf(w, "data: %s\n\n", event)
					}
				case "POST /agents/sessions/session/events":
					var body map[string]any
					if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
						t.Errorf("decode request: %v", err)
					}
					posts <- body
					w.WriteHeader(http.StatusNoContent)
				default:
					t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
					w.WriteHeader(http.StatusNotFound)
				}
			}))
			defer server.Close()
			client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("synthetic"), option.WithMaxRetries(0))
			stream := client.Beta.Agents.Sessions.Stream(context.Background(), "session", openai.AgentSessionStreamParams{
				Input:        "What is my USDC balance?",
				ToolHandlers: map[string]openai.AgentToolHandler{definition.Name: action.handle},
			})
			defer func() { _ = stream.Close() }()
			for stream.Next() {
			}
			if err := stream.Err(); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(wallet.args, []balanceArguments{{Asset: "USDC"}}) || len(posts) != 2 {
				t.Fatalf("unexpected dispatch: calls=%v, posts=%d", wallet.args, len(posts))
			}
			<-posts // The user input.
			result := (<-posts)["events"].([]any)[0].(map[string]any)
			if fail {
				if result["error"] != "Tool handler failed." || result["success"] == true {
					t.Fatalf("provider failure not safely reported: %v", result)
				}
			} else {
				var receipt balanceReceipt
				if err := json.Unmarshal([]byte(result["output"].(string)), &receipt); err != nil {
					t.Fatal(err)
				}
				if result["success"] != true || receipt != (balanceReceipt{Wallet: "bound-wallet", Asset: "USDC", Balance: "12.50"}) {
					t.Fatalf("unexpected typed receipt: %v", result)
				}
			}
		})
	}
}

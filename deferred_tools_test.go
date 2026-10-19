package openai_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/responses"
)

func deferredToolSchema() map[string]any {
	return map[string]any{
		"type": "object", "additionalProperties": false,
		"properties": map[string]any{"item_id": map[string]any{"type": "string"}},
		"required":   []string{"item_id"},
	}
}

func checkDeferredDefinition(t *testing.T, tools []any, flag param.Opt[bool]) {
	t.Helper()
	if len(tools) != 2 || tools[0].(map[string]any)["type"] != "tool_search" {
		t.Errorf("unexpected tools: %#v", tools)
		return
	}
	function := tools[1].(map[string]any)
	value, present := function["defer_loading"]
	if present != flag.Valid() || (present && value != flag.Value) {
		t.Errorf("defer_loading=%v (present=%t), want %v (present=%t)", value, present, flag.Value, flag.Valid())
	}
	if function["type"] != "function" || function["name"] != "lookup_item" || function["parameters"].(map[string]any)["additionalProperties"] != false {
		t.Errorf("function definition changed: %#v", function)
	}
}

func TestDeferredFunctionTools(t *testing.T) {
	for _, flag := range []struct {
		name  string
		value param.Opt[bool]
	}{{"omitted", param.Opt[bool]{}}, {"false", openai.Bool(false)}, {"true", openai.Bool(true)}} {
		t.Run(flag.name, func(t *testing.T) {
			t.Run("agents", func(t *testing.T) {
				m, client := newAgentHelperServer(t, append([]string{
					agentCreated("created", "turn", "null"),
					agentCall("call", "turn", "call", "lookup_item", `{"item_id":"A123"}`),
				}, agentEnd("turn")...)...)
				session, err := client.Beta.Agents.Sessions.New(context.Background(), openai.BetaAgentSessionNewParams{
					Environment: openai.EnvironmentParamUnion{OfParamOpenAIHosted: &openai.EnvironmentParamOpenAIHosted{}},
					Agent: openai.BetaAgentSessionNewParamsAgent{Model: openai.String("test-model"), Tools: []openai.AgentToolParamUnion{
						{OfParamToolSearch: &openai.AgentToolParamToolSearch{}},
						{OfParamFunction: &openai.AgentToolParamFunction{
							Name: "lookup_item", Description: "Look up a catalog item", Parameters: deferredToolSchema(), DeferLoading: flag.value,
						}},
					}},
				}, option.WithMiddleware(func(r *http.Request, _ option.MiddlewareNext) (*http.Response, error) {
					var body map[string]any
					if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
						return nil, err
					}
					checkDeferredDefinition(t, body["agent"].(map[string]any)["tools"].([]any), flag.value)
					return &http.Response{StatusCode: http.StatusOK, Header: http.Header{"Content-Type": {"application/json"}}, Body: io.NopCloser(strings.NewReader(`{"id":"session","status":"idle"}`)), Request: r}, nil
				}))
				if err != nil {
					t.Fatal(err)
				}
				calls := 0
				consumeAgentStream(t, client.Beta.Agents.Sessions.Stream(context.Background(), session.ID, openai.AgentSessionStreamParams{
					Input: "Look up item A123.", ToolHandlers: map[string]openai.AgentToolHandler{
						"lookup_item": func(_ context.Context, args map[string]any) (any, error) {
							calls++
							if args["item_id"] != "A123" {
								t.Fatalf("unexpected arguments: %#v", args)
							}
							return "in stock", nil
						},
					},
				}))
				posts := m.submissions()
				if calls != 1 || len(posts) != 2 {
					t.Fatalf("calls=%d posts=%d", calls, len(posts))
				}
				result := posts[1].body["events"].([]any)[0].(map[string]any)
				if result["call_id"] != "call" || result["success"] != true || result["output"] != "in stock" {
					t.Fatalf("unexpected tool result: %#v", result)
				}
			})
			for _, streaming := range []bool{false, true} {
				t.Run(fmt.Sprintf("responses/stream=%t", streaming), func(t *testing.T) {
					server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						var body map[string]any
						if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
							t.Error(err)
							w.WriteHeader(http.StatusBadRequest)
							return
						}
						checkDeferredDefinition(t, body["tools"].([]any), flag.value)
						if body["stream"] == true != streaming {
							t.Error("stream option did not reach request")
						}
						const response = `{"id":"resp_test","output":[{"type":"function_call","id":"fc_test","call_id":"call","name":"lookup_item","arguments":"{\"item_id\":\"A123\"}","status":"completed"}]}`
						if streaming {
							w.Header().Set("Content-Type", "text/event-stream")
							_, _ = fmt.Fprintf(w, "event: response.completed\ndata: {\"type\":\"response.completed\",\"sequence_number\":1,\"response\":%s}\n\n", response)
						} else {
							w.Header().Set("Content-Type", "application/json")
							_, _ = io.WriteString(w, response)
						}
					}))
					defer server.Close()
					client := openai.NewClient(option.WithBaseURL(server.URL), option.WithHTTPClient(server.Client()), option.WithAPIKey("synthetic"), option.WithMaxRetries(0))
					params := responses.ResponseNewParams{
						Model: "test-model", Input: responses.ResponseNewParamsInputUnion{OfString: openai.String("Look up item A123.")},
						Tools: []responses.ToolUnionParam{
							{OfToolSearch: &responses.ToolSearchToolParam{}},
							{OfFunction: &responses.FunctionToolParam{Name: "lookup_item", Parameters: deferredToolSchema(), Strict: openai.Bool(true), DeferLoading: flag.value}},
						},
					}
					var result *responses.Response
					if streaming {
						stream := client.Responses.NewStreaming(context.Background(), params)
						defer func() { _ = stream.Close() }()
						for stream.Next() {
							if event := stream.Current(); event.Type == "response.completed" {
								completed := event.AsResponseCompleted().Response
								result = &completed
							}
						}
						if err := stream.Err(); err != nil {
							t.Fatal(err)
						}
					} else {
						var err error
						result, err = client.Responses.New(context.Background(), params)
						if err != nil {
							t.Fatal(err)
						}
					}
					if result == nil || len(result.Output) != 1 || result.Output[0].Type != "function_call" {
						t.Fatal("missing function call")
					}
					call := result.Output[0].AsFunctionCall()
					var args struct {
						ItemID string `json:"item_id"`
					}
					if err := json.Unmarshal([]byte(call.Arguments), &args); err != nil {
						t.Fatal(err)
					}
					if call.Name != "lookup_item" || call.CallID != "call" || args.ItemID != "A123" {
						t.Fatalf("unexpected function call: %#v", call)
					}
				})
			}
		})
	}
}

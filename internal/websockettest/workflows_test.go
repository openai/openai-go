package websockettest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	wire "github.com/coder/websocket"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/responses"
)

// Exercise the README recipes through the public client, checking commands at
// the peer rather than asserting how the service implements its cache or queue.
func TestResponsesWorkflowRecipes(t *testing.T) {
	const window = `[{"type":"message","role":"user","content":[{"type":"input_text","text":"Keep this context"}]},{"type":"compaction","id":"cmp_item","encrypted_content":"opaque","future_field":"retain"}]`
	completed := func(id, output, stream string) string {
		lane := ""
		if stream != "" {
			lane = `,"stream_id":"` + stream + `"`
		}
		return fmt.Sprintf(`{"type":"response.completed","sequence_number":1,"response":{"id":%q,"object":"response","status":"completed","output":%s}%s}`, id, output, lane)
	}
	steps := []struct {
		request string
		events  []string
	}{
		{`{"type":"response.create","model":"gpt-4o-mini","input":"Seed","generate":false}`, []string{completed("resp_warm", `[]`, "")}},
		{`{"type":"response.create","model":"gpt-4o-mini","input":"Look up the result","previous_response_id":"resp_warm"}`, []string{completed("resp_tool", `[{"type":"function_call","id":"fc_1","call_id":"call_1","name":"lookup","arguments":"{}","status":"completed"}]`, "")}},
		{`{"type":"response.create","model":"gpt-4o-mini","previous_response_id":"resp_tool","input":[{"type":"function_call_output","call_id":"call_1","output":"42"},{"role":"user","content":"Explain the result."}],"context_management":[{"type":"compaction","compact_threshold":20000}]}`, []string{completed("resp_parent", `[]`, "")}},
		{`{"type":"response.create","model":"gpt-4o-mini","input":"Explore another approach.","previous_response_id":"resp_parent","stream_id":"fork","store":false}`, []string{`{"type":"response.created","stream_id":"fork","response":{"id":"resp_fork"}}`, `{"type":"response.in_progress","stream_id":"fork","response":{"id":"resp_fork"}}`}},
		{`{"type":"response.create","model":"gpt-4o-mini","input":"Advance source","previous_response_id":"resp_parent"}`, []string{completed("resp_source", `[]`, ""), completed("resp_fork", `[]`, "fork")}},
		{`{"type":"response.create","model":"gpt-4o-mini","input":` + window + `}`, []string{completed("resp_fresh", `[]`, "")}},
	}
	done := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer close(done)
		socket, err := wire.Accept(w, r, nil)
		if err != nil {
			t.Errorf("accept workflow connection: %v", err)
			return
		}
		defer func() { _ = socket.CloseNow() }()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		for i, step := range steps {
			_, data, err := socket.Read(ctx)
			if err != nil {
				t.Errorf("read workflow step %d: %v", i, err)
				return
			}
			var got, want any
			if err := json.Unmarshal(data, &got); err != nil {
				t.Errorf("decode workflow step %d: %v", i, err)
				return
			}
			if err := json.Unmarshal([]byte(step.request), &want); err != nil {
				t.Errorf("decode expected workflow step %d: %v", i, err)
				return
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("workflow step %d = %s, want %s", i, data, step.request)
				return
			}
			for _, event := range step.events {
				if err := socket.Write(ctx, wire.MessageText, []byte(event)); err != nil {
					t.Errorf("send workflow step %d: %v", i, err)
					return
				}
			}
		}
		_, _, _ = socket.Read(ctx)
	}))
	defer server.Close()
	client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("fixture-key"))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	connection, err := client.Responses.Connect(ctx, responses.ResponseConnectionOptions{})
	if err != nil {
		t.Fatal(err)
	}
	defer connection.Abort()
	create := func(command responses.ResponsesClientEventResponseCreateParam) *responses.Response {
		t.Helper()
		if createErr := connection.Create(ctx, command); createErr != nil {
			t.Fatal(createErr)
		}
		response, finalErr := connection.FinalResponse(ctx)
		if finalErr != nil {
			t.Fatal(finalErr)
		}
		return response
	}
	command := responses.ResponsesClientEventResponseCreateParam{
		Model: "gpt-4o-mini",
		Input: responses.ResponsesClientEventResponseCreateInputUnionParam{OfString: openai.String("Seed")},
	}
	command.SetExtraFields(map[string]any{"generate": false})
	warmup := create(command)
	command.SetExtraFields(nil)
	command.PreviousResponseID = openai.String(warmup.ID)
	command.Input.OfString = openai.String("Look up the result")
	toolResponse := create(command)
	output := responses.ResponseInputItemParamOfFunctionCallOutput("42")
	output.OfFunctionCallOutput.CallID = openai.String(toolResponse.Output[0].AsFunctionCall().CallID)
	input := responses.ResponseInputParam{output, responses.ResponseInputItemParamOfMessage("Explain the result.", responses.EasyInputMessageRoleUser)}
	command.PreviousResponseID = openai.String(toolResponse.ID)
	command.Input = responses.ResponsesClientEventResponseCreateInputUnionParam{OfResponse: &input}
	command.ContextManagement = []responses.ResponsesClientEventResponseCreateContextManagementParam{{Type: "compaction", CompactThreshold: openai.Int(20000)}}
	parent := create(command)
	fork, err := connection.Lane("fork")
	if err != nil {
		t.Fatal(err)
	}
	defer fork.Close()
	if err := connection.Create(ctx, responses.ResponsesClientEventResponseCreateParam{
		Model: "gpt-4o-mini", Store: openai.Bool(false), StreamID: openai.String("fork"), PreviousResponseID: openai.String(parent.ID),
		Input: responses.ResponsesClientEventResponseCreateInputUnionParam{OfString: openai.String("Explore another approach.")},
	}); err != nil {
		t.Fatal(err)
	}
	for {
		event, err := fork.Recv(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if event.Type == "response.in_progress" {
			break
		}
		if event.Type != "response.created" {
			t.Fatalf("fork readiness event = %q, want created or in_progress", event.Type)
		}
	}
	create(responses.ResponsesClientEventResponseCreateParam{
		Model: "gpt-4o-mini", PreviousResponseID: openai.String(parent.ID),
		Input: responses.ResponsesClientEventResponseCreateInputUnionParam{OfString: openai.String("Advance source")},
	})
	if _, err := fork.FinalResponse(ctx); err != nil {
		t.Fatal(err)
	}
	var compacted responses.CompactedResponse
	if err := json.Unmarshal([]byte(`{"id":"cmp_resource","object":"response.compaction","output":`+window+`}`), &compacted); err != nil {
		t.Fatal(err)
	}
	var compactedInput responses.ResponseInputParam
	for _, item := range compacted.Output {
		compactedInput = append(compactedInput, param.Override[responses.ResponseInputItemUnionParam](json.RawMessage(item.RawJSON())))
	}
	fresh := create(responses.ResponsesClientEventResponseCreateParam{
		Model: "gpt-4o-mini", Input: responses.ResponsesClientEventResponseCreateInputUnionParam{OfResponse: &compactedInput},
	})
	if fresh.ID != "resp_fresh" {
		t.Errorf("compacted restart response = %q, want resp_fresh", fresh.ID)
	}
	_ = connection.Close()
	<-done
}

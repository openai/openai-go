package websockettest

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/openai/openai-go/v3/responses"
	"github.com/tidwall/gjson"
)

func richEvent(t testing.TB, a *responses.ResponseAccumulator, event map[string]any) {
	t.Helper()
	data, err := json.Marshal(event)
	if err != nil {
		t.Fatal(err)
	}
	addMCPProgressEvent(t, a, string(data))
}

func richCreated(t testing.TB, a *responses.ResponseAccumulator) {
	t.Helper()
	richEvent(t, a, map[string]any{"type": "response.created", "sequence_number": 0, "response": richLifecycleResponse("in_progress", []any{})})
}

func richLifecycleResponse(status string, output []any) map[string]any {
	return map[string]any{"id": "resp_fixture", "object": "response", "created_at": 0, "status": status,
		"error": nil, "incomplete_details": nil, "instructions": nil, "metadata": map[string]any{"fixture": "kept"},
		"model": "fixture-model", "output": output, "parallel_tool_calls": true, "temperature": 1,
		"tool_choice": "auto", "tools": []any{}, "top_p": 1}
}

func richValue(t *testing.T, a *responses.ResponseAccumulator, path string) string {
	t.Helper()
	snapshot := a.DetailedSnapshot()
	if len(snapshot.Output) != 1 {
		t.Fatalf("want one observed item, got %d", len(snapshot.Output))
	}
	data, err := json.Marshal(snapshot.Output[0])
	if err != nil {
		t.Fatal(err)
	}
	return gjson.GetBytes(data, path).String()
}

func TestResponsesAccumulatorRichProgress(t *testing.T) {
	reasoning := func() map[string]any {
		return map[string]any{"id": "item_fixture", "type": "reasoning", "summary": []any{map[string]any{"type": "summary_text", "text": "seed "}}, "content": []any{map[string]any{"type": "reasoning_text", "text": "seed "}}, "status": "in_progress", "provider_metadata": map[string]any{"keep": true}}
	}
	code := func() map[string]any {
		return map[string]any{"id": "item_fixture", "type": "code_interpreter_call", "code": "seed ", "container_id": "container_fixture", "outputs": []any{}, "status": "in_progress"}
	}
	shell := func() map[string]any {
		return map[string]any{"id": "item_fixture", "type": "shell_call", "call_id": "call_fixture", "status": "in_progress", "action": map[string]any{"commands": []any{"seed "}, "timeout_ms": 1000, "max_output_length": 1000}, "environment": map[string]any{"type": "local"}}
	}
	shellOutput := func() map[string]any {
		return map[string]any{"id": "item_fixture", "type": "shell_call_output", "status": "in_progress", "max_output_length": 1000, "call_id": "call_fixture", "output": []any{map[string]any{"stdout": "seed ", "stderr": "old ", "outcome": map[string]any{"type": "exit", "exit_code": 0}, "provider_metadata": true}}}
	}
	type progressCase struct {
		event      string
		item       func() map[string]any
		fields     map[string]any
		path, want string
	}
	cases := []progressCase{
		{"response.reasoning_summary_part.added", reasoning, map[string]any{"summary_index": 1, "part": map[string]any{"type": "summary_text", "text": "added"}}, "Item.summary.1.text", "added"},
		{"response.reasoning_summary_part.done", reasoning, map[string]any{"summary_index": 0, "part": map[string]any{"type": "summary_text", "text": "corrected"}, "status": "completed"}, "Item.summary.0.text", "corrected"},
		{"response.reasoning_summary_text.delta", reasoning, map[string]any{"summary_index": 0, "delta": "kept"}, "Item.summary.0.text", "seed kept"},
		{"response.reasoning_summary_text.done", reasoning, map[string]any{"summary_index": 0, "text": "corrected"}, "Item.summary.0.text", "corrected"},
		{"response.reasoning_text.delta", reasoning, map[string]any{"content_index": 0, "delta": "kept"}, "Content.0.text", "seed kept"},
		{"response.reasoning_text.done", reasoning, map[string]any{"content_index": 0, "text": "corrected"}, "Content.0.text", "corrected"},
		{"response.shell_call_command.added", shell, map[string]any{"command_index": 1, "command": "added"}, "Item.action.commands.1", "added"},
		{"response.shell_call_command.delta", shell, map[string]any{"command_index": 0, "delta": "kept"}, "Item.action.commands.0", "seed kept"},
		{"response.shell_call_command.done", shell, map[string]any{"command_index": 0, "command": "corrected"}, "Item.action.commands.0", "corrected"},
		{"response.shell_call_output_content.delta", shellOutput, map[string]any{"command_index": 0, "delta": map[string]any{"stdout": "kept", "stderr": "err"}}, "Content.0.stdout", "kept"},
		{"response.shell_call_output_content.done", shellOutput, map[string]any{"command_index": 0, "output": []any{map[string]any{"stdout": "corrected", "stderr": "final", "outcome": map[string]any{"type": "exit", "exit_code": 7}, "provider_metadata": true}}}, "Content.0.output.0.stdout", "corrected"},
		{"response.code_interpreter_call_code.delta", code, map[string]any{"delta": "kept"}, "Item.code", "seed kept"},
		{"response.code_interpreter_call_code.done", code, map[string]any{"code": "corrected"}, "Item.code", "corrected"},
	}
	for _, group := range []struct {
		prefix   string
		item     func() map[string]any
		statuses []string
	}{
		{"response.code_interpreter_call", code, []string{"in_progress", "interpreting", "completed"}},
		{"response.file_search_call", func() map[string]any {
			return map[string]any{"id": "item_fixture", "type": "file_search_call", "queries": []any{}, "status": "in_progress", "results": nil}
		}, []string{"in_progress", "searching", "completed"}},
		{"response.web_search_call", func() map[string]any {
			return map[string]any{"id": "item_fixture", "type": "web_search_call", "status": "in_progress", "action": map[string]any{"type": "search", "query": "fixture"}}
		}, []string{"in_progress", "searching", "completed"}},
		{"response.image_generation_call", func() map[string]any {
			return map[string]any{"id": "item_fixture", "type": "image_generation_call", "status": "in_progress", "result": nil}
		}, []string{"in_progress", "generating", "completed"}},
		{"response.mcp_call", func() map[string]any {
			return map[string]any{"id": "item_fixture", "type": "mcp_call", "name": "lookup", "server_label": "fixture", "arguments": "{}", "approval_request_id": nil, "error": nil, "output": nil}
		}, []string{"in_progress", "completed", "failed"}},
	} {
		for _, status := range group.statuses {
			cases = append(cases, progressCase{group.prefix + "." + status, group.item, nil, "Item.status", status})
		}
	}
	if len(cases) != 28 {
		t.Fatalf("want all28 rich event cases, got%d", len(cases))
	}
	for _, tc := range cases {
		t.Run(tc.event, func(t *testing.T) {
			var acc responses.ResponseAccumulator
			richCreated(t, &acc)
			richEvent(t, &acc, map[string]any{"type": "response.output_item.added", "sequence_number": 1, "output_index": 0, "item": tc.item()})
			prior := acc.DetailedSnapshot()
			priorJSON, _ := json.Marshal(prior)
			event := map[string]any{"type": tc.event, "sequence_number": 2, "output_index": 0, "item_id": "item_fixture"}
			for key, value := range tc.fields {
				event[key] = value
			}
			if tc.event == "response.shell_call_command.added" || tc.event == "response.shell_call_command.delta" || tc.event == "response.shell_call_command.done" {
				delete(event, "item_id")
			}
			richEvent(t, &acc, event)
			if got := richValue(t, &acc, tc.path); got != tc.want {
				t.Errorf("%s=%q, want%q", tc.path, got, tc.want)
			}
			if strings.HasPrefix(tc.event, "response.reasoning_text.") && richValue(t, &acc, "Item.content.0.text") != tc.want {
				t.Fatal("reasoning content views disagree")
			}
			if acc.Snapshot().OutputText() != "" {
				t.Error("rich progress became user-facing output text")
			}
			afterJSON, _ := json.Marshal(prior)
			if !reflect.DeepEqual(priorJSON, afterJSON) {
				t.Error("prior snapshot changed")
			}

			observed := acc.DetailedSnapshot()
			for _, invalidIndex := range []any{-1, "0", 0.5} {
				malformed := make(map[string]any)
				for key, value := range event {
					malformed[key] = value
				}
				malformed["output_index"] = invalidIndex
				raw, _ := json.Marshal(malformed)
				var decoded responses.ResponsesServerEventUnion
				if err := json.Unmarshal(raw, &decoded); err != nil {
					t.Fatal(err)
				}
				if err := acc.AddEvent(decoded); err == nil {
					t.Errorf("accepted malformed output index %v", invalidIndex)
				}
				if !reflect.DeepEqual(observed, acc.DetailedSnapshot()) {
					t.Fatal("malformed event changed progress")
				}
			}
			if _, present := event["item_id"]; present {
				event["item_id"] = "other_item"
				richEvent(t, &acc, event)
				if !reflect.DeepEqual(observed, acc.DetailedSnapshot()) {
					t.Fatal("event for another item changed progress")
				}
				event["item_id"] = "item_fixture"
			}
			copy := acc.DetailedSnapshot()
			copy.Output[0].Item["caller_change"] = json.RawMessage(`true`)
			if !reflect.DeepEqual(observed, acc.DetailedSnapshot()) {
				t.Fatal("caller mutation reached accumulator")
			}
			finalItem := tc.item()
			if _, present := finalItem["status"]; present {
				finalItem["status"] = "completed"
				prior.Output[0].Item["status"] = json.RawMessage(`"completed"`)
			}
			richEvent(t, &acc, map[string]any{"type": "response.output_item.done", "sequence_number": 3, "output_index": 0, "item": finalItem})
			if !reflect.DeepEqual(prior, acc.DetailedSnapshot()) {
				t.Fatal("full item did not replace provisional progress")
			}
			richEvent(t, &acc, map[string]any{"type": "response.completed", "sequence_number": 4, "response": richLifecycleResponse("completed", []any{})})
			if len(acc.DetailedSnapshot().Output) != 0 {
				t.Fatal("explicit empty terminal output did not clear progress")
			}
		})
	}
}

func TestResponsesAccumulatorReasoningContentParts(t *testing.T) {
	for _, kind := range []string{"added", "done"} {
		t.Run(kind, func(t *testing.T) {
			var acc responses.ResponseAccumulator
			richCreated(t, &acc)
			richEvent(t, &acc, map[string]any{"type": "response.output_item.added", "sequence_number": 1, "output_index": 0, "item": map[string]any{"id": "reasoning_fixture", "type": "reasoning", "summary": []any{}}})
			richEvent(t, &acc, map[string]any{"type": "response.content_part." + kind, "sequence_number": 2, "output_index": 0, "content_index": 3, "item_id": "reasoning_fixture", "part": map[string]any{"type": "reasoning_text", "text": "observed", "provider_metadata": true}})
			if got := richValue(t, &acc, "Content.3.text"); got != "observed" {
				t.Errorf("reasoning part=%q", got)
			}
			if acc.Snapshot().OutputText() != "" {
				t.Error("reasoning part became output text")
			}
		})
	}
}

func TestResponsesAccumulatorQueuedMetadata(t *testing.T) {
	var acc responses.ResponseAccumulator
	richCreated(t, &acc)
	richEvent(t, &acc, map[string]any{"type": "response.queued", "sequence_number": 1, "response": map[string]any{"id": "resp_fixture", "object": "response", "created_at": 0, "status": "queued", "error": nil, "incomplete_details": nil, "instructions": nil, "metadata": map[string]any{"fixture": "queued"}, "model": "fixture-model", "output": []any{}, "parallel_tool_calls": true, "temperature": 1, "tool_choice": "auto", "tools": []any{}, "top_p": 1}})
	got := acc.DetailedSnapshot()
	if string(got.Response["status"]) != `"queued"` || string(got.Response["metadata"]) != `{"fixture":"queued"}` {
		t.Errorf("queued metadata not refreshed: status=%s metadata=%s", got.Response["status"], got.Response["metadata"])
	}
}

func TestResponsesAccumulatorSparseMCPProgress(t *testing.T) {
	var acc responses.ResponseAccumulator
	richEvent(t, &acc, map[string]any{"type": "response.mcp_call_arguments.delta", "sequence_number": 0, "output_index": int64(1) << 40, "item_id": "mcp_sparse", "delta": "{"})
	richEvent(t, &acc, map[string]any{"type": "response.mcp_call_arguments.done", "sequence_number": 1, "output_index": int64(1) << 40, "item_id": "mcp_sparse", "arguments": "{\"observed\":true}"})
	snapshot := acc.Snapshot()
	if len(snapshot.Output) != 1 {
		t.Fatalf("missing sparse MCP progress: %#v", snapshot)
	}
	item := snapshot.Output[0]
	if item.OutputIndex != int64(1)<<40 || item.ItemID != "mcp_sparse" || item.Arguments != "{\"observed\":true}" {
		t.Errorf("unexpected sparse MCP observation: %#v", item)
	}
	detail := acc.DetailedSnapshot().Output[0].Item
	for _, key := range []string{"type", "name", "server_label", "status", "call_id"} {
		if _, ok := detail[key]; ok {
			t.Errorf("fabricated MCP field%s=%s", key, detail[key])
		}
	}
	if _, err := acc.ReconstructedResponse(); err == nil {
		t.Fatal("sparse progress invented terminal response")
	}
}

func TestResponsesAccumulatorRichProgressLargeFragments(t *testing.T) {
	var acc responses.ResponseAccumulator
	richCreated(t, &acc)
	richEvent(t, &acc, map[string]any{"type": "response.output_item.added", "sequence_number": 1, "output_index": 0, "item": map[string]any{"id": "reasoning_fixture", "type": "reasoning", "summary": []any{map[string]any{"type": "summary_text", "text": ""}}, "content": []any{map[string]any{"type": "reasoning_text", "text": ""}}}})
	fragment := strings.Repeat("x", 1024)
	for i := 0; i < 4096; i++ {
		richEvent(t, &acc, map[string]any{"type": "response.reasoning_summary_text.delta", "sequence_number": i + 2, "output_index": 0, "summary_index": 0, "item_id": "reasoning_fixture", "delta": fragment})
	}
	if got := richValue(t, &acc, "Item.summary.0.text"); len(got) != 4096*len(fragment) {
		t.Errorf("large reasoning text length=%d", len(got))
	}
	old := acc.DetailedSnapshot()
	richEvent(t, &acc, map[string]any{"type": "response.reasoning_summary_text.done", "sequence_number": 4098, "output_index": 0, "summary_index": 0, "item_id": "reasoning_fixture", "text": "short correction"})
	if got := richValue(t, &acc, "Item.summary.0.text"); got != "short correction" {
		t.Errorf("done correction=%q", got)
	}
	if got := gjson.GetBytes(old.Output[0].Item["summary"], "0.text").String(); len(got) != 4096*len(fragment) {
		t.Errorf("old large snapshot changed length=%d", len(got))
	}
}

func TestResponsesAccumulatorRichProgressDoesNotBindWithoutItem(t *testing.T) {
	var acc responses.ResponseAccumulator
	richEvent(t, &acc, map[string]any{"type": "response.code_interpreter_call_code.delta", "sequence_number": 0, "stream_id": "unrelated", "output_index": 0, "item_id": "unknown", "delta": "ignored"})
	addMCPProgressEvent(t, &acc, `{"type":"response.output_text.delta","sequence_number":1,"stream_id":"active","output_index":0,"content_index":0,"item_id":"message","delta":"kept"}`)
	if got := acc.Snapshot(); got.StreamID != "active" || got.OutputText() != "kept" {
		t.Errorf("unmatched rich event bound accumulator: %#v", got)
	}
}

func TestResponsesAccumulatorShellCommandExtraMetadata(t *testing.T) {
	for _, kind := range []string{"added", "delta", "done"} {
		t.Run(kind, func(t *testing.T) {
			var acc responses.ResponseAccumulator
			richCreated(t, &acc)
			richEvent(t, &acc, map[string]any{"type": "response.output_item.added", "sequence_number": 1, "output_index": 0, "item": map[string]any{"id": "shell", "type": "shell_call", "call_id": "call", "status": "in_progress", "action": map[string]any{"commands": []any{"seed"}, "timeout_ms": 1000, "max_output_length": 1000}, "environment": map[string]any{"type": "local"}}})
			index := 0
			if kind == "added" {
				index = 1
			}
			event := map[string]any{"type": "response.shell_call_command." + kind, "sequence_number": 2, "output_index": 0, "command_index": index, "item_id": map[string]any{"future": true}, "command": "observed", "delta": "observed"}
			richEvent(t, &acc, event)
			want := "observed"
			if kind == "delta" {
				want = "seedobserved"
			}
			data, _ := json.Marshal(acc.DetailedSnapshot().Output[0].Item["action"])
			commands := gjson.GetBytes(data, "commands").Array()
			if len(commands) <= index || commands[index].Str != want {
				t.Fatalf("extra metadata blocked command: %s", data)
			}
		})
	}
}

func TestResponsesAccumulatorSparseShellOutput(t *testing.T) {
	for _, kind := range []string{"delta", "done"} {
		for _, index := range []int{0, 1} {
			t.Run(kind+string(rune('0'+index)), func(t *testing.T) {
				var acc responses.ResponseAccumulator
				richCreated(t, &acc)
				richEvent(t, &acc, map[string]any{"type": "response.output_item.added", "sequence_number": 1, "output_index": 0, "item": map[string]any{"id": "shell", "type": "shell_call", "call_id": "call", "status": "in_progress", "action": map[string]any{"commands": []any{"one", "two"}, "timeout_ms": 1000, "max_output_length": 1000}, "environment": map[string]any{"type": "local"}}})
				richEvent(t, &acc, map[string]any{"type": "response.output_item.added", "sequence_number": 2, "output_index": 1, "item": map[string]any{"id": "shell_output", "type": "shell_call_output", "status": "in_progress", "max_output_length": 1000, "call_id": "call", "output": []any{}}})
				event := map[string]any{"type": "response.shell_call_output_content." + kind, "sequence_number": 3, "output_index": 1, "item_id": "shell_output", "command_index": index, "delta": map[string]any{"stdout": "observed"}, "output": []any{map[string]any{"stdout": "observed", "stderr": "", "outcome": map[string]any{"type": "exit", "exit_code": 7}, "future": true}}}
				richEvent(t, &acc, event)
				item := acc.DetailedSnapshot().Output[1]
				row := item.Content[int64(index)]
				if kind == "delta" {
					if string(row["stdout"]) != `"observed"` {
						t.Fatalf("missing indexed shell progress: %#v", item)
					}
					for _, key := range []string{"stderr", "outcome"} {
						if _, ok := row[key]; ok {
							t.Fatalf("invented %s: %s", key, row[key])
						}
					}
				} else if gjson.GetBytes(row["output"], "0.future").Bool() != true || gjson.GetBytes(row["output"], "0.outcome.exit_code").Int() != 7 {
					t.Fatalf("lost full output metadata: %#v", row)
				}
				values := gjson.ParseBytes(item.Item["output"]).Array()
				if len(values) != 0 {
					t.Fatalf("guessed flat output placement: %s", item.Item["output"])
				}
				before := acc.DetailedSnapshot()
				for _, bad := range []int64{-1, 2, 1 << 40} {
					event["command_index"] = bad
					data, _ := json.Marshal(event)
					var decoded responses.ResponsesServerEventUnion
					if err := json.Unmarshal(data, &decoded); err != nil {
						t.Fatal(err)
					}
					if err := acc.AddEvent(decoded); err == nil {
						t.Fatalf("accepted out-of-bounds command %d", bad)
					}
					if !reflect.DeepEqual(before, acc.DetailedSnapshot()) {
						t.Fatal("invalid command changed observed content")
					}
				}
			})
		}
	}
}

func TestResponsesAccumulatorShellOutputWithoutKnownBound(t *testing.T) {
	var acc responses.ResponseAccumulator
	richCreated(t, &acc)
	richEvent(t, &acc, map[string]any{"type": "response.output_item.added", "sequence_number": 1, "output_index": 0, "item": map[string]any{"id": "shell_output", "type": "shell_call_output", "status": "in_progress", "max_output_length": 1000, "call_id": "unobserved_call", "output": []any{}}})
	for _, index := range []int{0, 7} {
		richEvent(t, &acc, map[string]any{"type": "response.shell_call_output_content.delta", "sequence_number": index + 2, "output_index": 0, "item_id": "shell_output", "command_index": index, "delta": map[string]any{"stderr": "observed"}})
	}
	item := acc.DetailedSnapshot().Output[0]
	if len(item.Content) != 2 || len(gjson.ParseBytes(item.Item["output"]).Array()) != 0 {
		t.Fatalf("missing sparse observations or guessed flat output: %#v", item)
	}
	for _, row := range item.Content {
		if len(row) != 1 || string(row["stderr"]) != `"observed"` {
			t.Fatalf("invented shell output fields: %#v", row)
		}
	}
}

func BenchmarkResponsesAccumulatorShellFragments(b *testing.B) {
	for _, count := range []int{1, 4096} {
		b.Run(strconv.Itoa(count), func(b *testing.B) {
			var acc responses.ResponseAccumulator
			richCreated(b, &acc)
			commands := make([]string, count)
			output := make([]any, count)
			for i := range commands {
				commands[i] = "echo observed"
				output[i] = map[string]any{"stdout": "", "stderr": "", "outcome": map[string]any{"type": "exit", "exit_code": 0}}
			}
			richEvent(b, &acc, map[string]any{"type": "response.output_item.added", "sequence_number": 1, "output_index": 0, "item": map[string]any{"id": "shell", "type": "shell_call", "call_id": "call", "status": "in_progress", "action": map[string]any{"commands": commands, "timeout_ms": 1000, "max_output_length": 1000}, "environment": map[string]any{"type": "local"}}})
			richEvent(b, &acc, map[string]any{"type": "response.output_item.added", "sequence_number": 2, "output_index": 1, "item": map[string]any{"id": "shell_output", "type": "shell_call_output", "status": "in_progress", "max_output_length": 1000, "call_id": "call", "output": output}})
			raw := `{"type":"response.shell_call_output_content.delta","sequence_number":3,"output_index":1,"item_id":"shell_output","command_index":0,"delta":{"stdout":"x"}}`
			var event responses.ResponsesServerEventUnion
			if err := json.Unmarshal([]byte(raw), &event); err != nil {
				b.Fatal(err)
			}
			if err := acc.AddEvent(event); err != nil {
				b.Fatal(err)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if err := acc.AddEvent(event); err != nil {
					b.Fatal(err)
				}
			}
			b.StopTimer()
			if got := acc.DetailedSnapshot().Output[1].Content[0]["stdout"]; len(got) != b.N+3 {
				b.Fatalf("lost shell fragments: %d encoded bytes", len(got))
			}
		})
	}
}

func TestResponsesAccumulatorShellOutputRefreshesCommandBounds(t *testing.T) {
	var acc responses.ResponseAccumulator
	richCreated(t, &acc)
	addCall := func(commands []string) {
		richEvent(t, &acc, map[string]any{"type": "response.output_item.added", "sequence_number": 1, "output_index": 0, "item": map[string]any{"id": "shell", "type": "shell_call", "call_id": "call", "status": "in_progress", "action": map[string]any{"commands": commands, "timeout_ms": 1000, "max_output_length": 1000}, "environment": map[string]any{"type": "local"}}})
	}
	addCall([]string{"one", "two"})
	richEvent(t, &acc, map[string]any{"type": "response.output_item.added", "sequence_number": 2, "output_index": 1, "item": map[string]any{"id": "shell_output", "type": "shell_call_output", "status": "in_progress", "max_output_length": 1000, "call_id": "call", "output": []any{}}})
	raw := `{"type":"response.shell_call_output_content.delta","sequence_number":3,"output_index":1,"item_id":"shell_output","command_index":1,"delta":{"stdout":"x"}}`
	addMCPProgressEvent(t, &acc, raw)
	addCall([]string{"one"})
	before := acc.DetailedSnapshot()
	var event responses.ResponsesServerEventUnion
	if err := json.Unmarshal([]byte(raw), &event); err != nil {
		t.Fatal(err)
	}
	if err := acc.AddEvent(event); err == nil {
		t.Fatal("retained obsolete command bound after full replacement")
	}
	if !reflect.DeepEqual(before, acc.DetailedSnapshot()) {
		t.Fatal("rejected command changed output")
	}
	richEvent(t, &acc, map[string]any{"type": "response.shell_call_command.added", "sequence_number": 4, "output_index": 0, "command_index": 1, "command": "two"})
	addMCPProgressEvent(t, &acc, raw)
	if got := acc.DetailedSnapshot().Output[1].Content[1]["stdout"]; string(got) != `"xx"` {
		t.Fatalf("new command bound not observed: %s", got)
	}
	acc.Reset()
	richCreated(t, &acc)
	richEvent(t, &acc, map[string]any{"type": "response.output_item.added", "sequence_number": 1, "output_index": 0, "item": map[string]any{"id": "shell_output", "type": "shell_call_output", "status": "in_progress", "max_output_length": 1000, "call_id": "call", "output": []any{}}})
	richEvent(t, &acc, map[string]any{"type": "response.shell_call_output_content.delta", "sequence_number": 2, "output_index": 0, "item_id": "shell_output", "command_index": 0, "delta": map[string]any{"stdout": "reset"}})
	if got := richValue(t, &acc, "Content.0.stdout"); got != "reset" {
		t.Fatalf("stale bounds after Reset: %q", got)
	}
}

func TestResponsesAccumulatorShellOutputCommandContents(t *testing.T) {
	var acc responses.ResponseAccumulator
	richCreated(t, &acc)
	richEvent(t, &acc, map[string]any{"type": "response.output_item.added", "sequence_number": 1, "output_index": 0, "item": map[string]any{"id": "shell", "type": "shell_call", "call_id": "call", "status": "in_progress", "action": map[string]any{"commands": []any{"one", "two"}, "timeout_ms": 1000, "max_output_length": 1000}, "environment": map[string]any{"type": "local"}}})
	item := map[string]any{"id": "shell_output", "type": "shell_call_output", "status": "in_progress", "max_output_length": 1000, "call_id": "call", "output": []any{map[string]any{"stdout": "ungrouped", "stderr": "", "outcome": map[string]any{"type": "exit", "exit_code": 0}}}}
	richEvent(t, &acc, map[string]any{"type": "response.output_item.added", "sequence_number": 2, "output_index": 1, "item": item})
	original := acc.DetailedSnapshot().Output[1].Item["output"]
	for i, step := range []struct {
		index int
		text  string
	}{{1, "first"}, {0, "other"}, {1, "+tail"}} {
		richEvent(t, &acc, map[string]any{"type": "response.shell_call_output_content.delta", "sequence_number": i + 3, "output_index": 1, "item_id": "shell_output", "command_index": step.index, "delta": map[string]any{"stdout": step.text}})
	}
	prior := acc.DetailedSnapshot()
	if got := string(prior.Output[1].Content[1]["stdout"]); got != `"first+tail"` {
		t.Fatalf("invented flat-to-command grouping: %s", got)
	}
	rows := []any{
		map[string]any{"stdout": "chunk-one", "stderr": "", "outcome": map[string]any{"type": "exit", "exit_code": 7}, "created_by": "actor-one", "future": map[string]any{"kept": true}},
		map[string]any{"stdout": "chunk-two", "stderr": "timeout", "outcome": map[string]any{"type": "timeout"}, "created_by": "actor-two"},
	}
	richEvent(t, &acc, map[string]any{"type": "response.shell_call_output_content.done", "sequence_number": 6, "output_index": 1, "item_id": "shell_output", "command_index": 1, "output": rows})
	got := acc.DetailedSnapshot().Output[1]
	want, err := json.Marshal(rows)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Content) != 2 || len(got.Content[1]) != 1 || string(got.Content[1]["output"]) != string(want) {
		t.Fatalf("lost per-command chunks/metadata: %#v", got.Content)
	}
	if string(got.Content[0]["stdout"]) != `"other"` {
		t.Fatal("one command replaced another command's progress")
	}
	if !reflect.DeepEqual(original, got.Item["output"]) {
		t.Fatal("guessed placement in ungrouped flat output")
	}
	if string(prior.Output[1].Content[1]["stdout"]) != `"first+tail"` {
		t.Fatal("done mutated earlier returned snapshot")
	}
	got.Content[1]["output"][0] = '!'
	if string(acc.DetailedSnapshot().Output[1].Content[1]["output"]) != string(want) {
		t.Fatal("caller mutation changed retained chunks")
	}
	richEvent(t, &acc, map[string]any{"type": "response.shell_call_output_content.done", "sequence_number": 7, "output_index": 1, "item_id": "shell_output", "command_index": 1, "output": []any{}})
	if got := acc.DetailedSnapshot().Output[1].Content[1]; len(got) != 1 || string(got["output"]) != "[]" {
		t.Fatalf("empty done did not replace command contents: %#v", got)
	}
	item["status"] = "completed"
	item["output"] = rows
	richEvent(t, &acc, map[string]any{"type": "response.output_item.done", "sequence_number": 8, "output_index": 1, "item": item})
	if got := acc.DetailedSnapshot().Output[1]; len(got.Content) != 0 || string(got.Item["output"]) != string(want) {
		t.Fatal("full item did not replace all command facts")
	}
	richEvent(t, &acc, map[string]any{"type": "response.completed", "sequence_number": 9, "response": richLifecycleResponse("completed", []any{})})
	if len(acc.DetailedSnapshot().Output) != 0 {
		t.Fatal("terminal empty output did not clear command facts")
	}
}

func TestResponsesAccumulatorShellOutputUnknownBounds(t *testing.T) {
	var acc responses.ResponseAccumulator
	richCreated(t, &acc)
	richEvent(t, &acc, map[string]any{"type": "response.output_item.added", "sequence_number": 1, "output_index": 0, "item": map[string]any{"id": "shell_output", "type": "shell_call_output", "status": "in_progress", "max_output_length": 1000, "call_id": "unobserved_call", "output": []any{}}})
	index := int64(1) << 40
	richEvent(t, &acc, map[string]any{"type": "response.shell_call_output_content.delta", "sequence_number": 2, "output_index": 0, "item_id": "shell_output", "command_index": index, "delta": map[string]any{"stderr": "observed"}})
	got := acc.DetailedSnapshot().Output[0]
	if len(got.Content) != 1 || len(got.Content[index]) != 1 || string(got.Content[index]["stderr"]) != `"observed"` || string(got.Item["output"]) != "[]" {
		t.Fatalf("unknown bounds lost sparse facts or invented output: %#v", got)
	}
	richEvent(t, &acc, map[string]any{"type": "response.shell_call_output_content.done", "sequence_number": 3, "output_index": 0, "item_id": "shell_output", "command_index": index, "output": []any{}})
	got = acc.DetailedSnapshot().Output[0]
	if len(got.Content) != 1 || len(got.Content[index]) != 1 || string(got.Content[index]["output"]) != "[]" {
		t.Fatalf("sparse empty command contents not retained: %#v", got)
	}
}

func TestResponsesAccumulatorShellOutputKnownEmptyBound(t *testing.T) {
	var acc responses.ResponseAccumulator
	richCreated(t, &acc)
	richEvent(t, &acc, map[string]any{"type": "response.output_item.added", "sequence_number": 1, "output_index": 0, "item": map[string]any{"id": "shell", "type": "shell_call", "call_id": "call", "status": "in_progress", "action": map[string]any{"commands": []any{}, "timeout_ms": 1000, "max_output_length": 1000}, "environment": map[string]any{"type": "local"}}})
	richEvent(t, &acc, map[string]any{"type": "response.output_item.added", "sequence_number": 2, "output_index": 1, "item": map[string]any{"id": "shell_output", "type": "shell_call_output", "status": "in_progress", "max_output_length": 1000, "call_id": "call", "output": []any{}}})
	prior := acc.DetailedSnapshot()
	var event responses.ResponsesServerEventUnion
	if err := json.Unmarshal([]byte(`{"type":"response.shell_call_output_content.delta","sequence_number":3,"output_index":1,"item_id":"shell_output","command_index":0,"delta":{"stdout":"rejected"}}`), &event); err != nil {
		t.Fatal(err)
	}
	if err := acc.AddEvent(event); err == nil {
		t.Fatal("explicit empty command array was treated as unknown")
	}
	if !reflect.DeepEqual(prior, acc.DetailedSnapshot()) {
		t.Fatal("rejected command changed snapshot")
	}
}

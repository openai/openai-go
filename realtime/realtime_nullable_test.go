package realtime_test

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"testing"

	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/respjson"
	"github.com/openai/openai-go/v3/realtime"
)

func TestRealtimeNullableFieldsThroughWebsocket(t *testing.T) {
	const (
		added         = `{"type":"conversation.item.added","event_id":"evt_added","previous_item_id":"prev_added","item":{"type":"mcp_approval_response","id":"approval_result","approval_request_id":"approval_1","approve":true,"reason":"approved synthetic"}}`
		created       = `{"type":"conversation.item.created","event_id":"evt_created","previous_item_id":"prev_created","item":{"type":"mcp_list_tools","id":"tools_1","server_label":"local","tools":[{"name":"lookup","input_schema":{"type":"object"},"description":"synthetic lookup","annotations":{"readOnlyHint":true}}]}}`
		done          = `{"type":"conversation.item.done","event_id":"evt_done","previous_item_id":"prev_done","item":{"type":"mcp_call","id":"call_1","server_label":"local","name":"lookup","arguments":"{}","approval_request_id":"approval_9","output":"synthetic result","error":{"type":"tool_execution_error","message":"synthetic tool error"}}}`
		committed     = `{"type":"input_audio_buffer.committed","event_id":"evt_audio","item_id":"audio_1","previous_item_id":"prev_audio"}`
		completed     = `{"type":"conversation.item.input_audio_transcription.completed","event_id":"evt_completed","item_id":"audio_1","content_index":0,"transcript":"hello","usage":{"type":"duration","seconds":1.25},"logprobs":[{"token":"hello","bytes":[104,101,108,108,111],"logprob":-0.25}]}`
		delta         = `{"type":"conversation.item.input_audio_transcription.delta","event_id":"evt_delta","item_id":"audio_1","content_index":0,"delta":"hi","logprobs":[{"token":"hi","bytes":[104,105],"logprob":-0.5}]}`
		errorEvent    = `{"type":"error","event_id":"evt_error","error":{"type":"invalid_request_error","message":"synthetic error","code":"invalid_event","event_id":"client_17","param":"audio.format"}}`
		response      = `{"type":"response.done","event_id":"evt_response","response":{"id":"resp_1","object":"realtime.response","status":"completed","output":[],"metadata":{"audit":"response_17"}}}`
		session       = `{"type":"session.updated","event_id":"evt_session","session":{"type":"realtime","prompt":{"id":"pmpt_audit"},"tracing":"auto","audio":{"input":{"turn_detection":{"type":"server_vad","idle_timeout_ms":6000}}},"tools":[{"type":"mcp","server_label":"local","server_url":"https://example.invalid/mcp","allowed_callers":["direct"],"allowed_tools":["lookup"],"headers":{"X-Synthetic":"audit"},"require_approval":"always"}]}}`
		transcription = `{"type":"session.updated","event_id":"evt_transcription","session":{"type":"transcription","audio":{"input":{"turn_detection":{"type":"server_vad","idle_timeout_ms":7000}}}}}`
		arguments     = `{"type":"response.mcp_call_arguments.delta","event_id":"evt_args","response_id":"resp_1","item_id":"call_1","output_index":0,"delta":"{}","obfuscation":"mask_audit"}`
	)
	cases := []struct {
		name, wire, path string
		want             any
		read             func(*realtime.RealtimeServerEventUnion) (any, respjson.Field)
	}{
		{"added/previous_item_id", added, "previous_item_id", "prev_added", func(e *realtime.RealtimeServerEventUnion) (any, respjson.Field) {
			v := e.OfConversationItemAdded
			return v.PreviousItemID, v.JSON.PreviousItemID
		}},
		{"added/item.reason", added, "item.reason", "approved synthetic", func(e *realtime.RealtimeServerEventUnion) (any, respjson.Field) {
			v := e.OfConversationItemAdded.Item.OfRealtimeMcpApprovalResponse
			return v.Reason, v.JSON.Reason
		}},
		{"created/previous_item_id", created, "previous_item_id", "prev_created", func(e *realtime.RealtimeServerEventUnion) (any, respjson.Field) {
			v := e.OfConversationItemCreatedEvent
			return v.PreviousItemID, v.JSON.PreviousItemID
		}},
		{"created/item.tools.0.description", created, "item.tools.0.description", "synthetic lookup", func(e *realtime.RealtimeServerEventUnion) (any, respjson.Field) {
			v := e.OfConversationItemCreatedEvent.Item.OfRealtimeMcpListTools.Tools[0]
			return v.Description, v.JSON.Description
		}},
		{"created/item.tools.0.annotations", created, "item.tools.0.annotations", map[string]any{"readOnlyHint": true}, func(e *realtime.RealtimeServerEventUnion) (any, respjson.Field) {
			v := e.OfConversationItemCreatedEvent.Item.OfRealtimeMcpListTools.Tools[0]
			return v.Annotations, v.JSON.Annotations
		}},
		{"done/previous_item_id", done, "previous_item_id", "prev_done", func(e *realtime.RealtimeServerEventUnion) (any, respjson.Field) {
			v := e.OfConversationItemDone
			return v.PreviousItemID, v.JSON.PreviousItemID
		}},
		{"done/item.approval_request_id", done, "item.approval_request_id", "approval_9", func(e *realtime.RealtimeServerEventUnion) (any, respjson.Field) {
			v := e.OfConversationItemDone.Item.OfRealtimeMcpToolCall
			return v.ApprovalRequestID, v.JSON.ApprovalRequestID
		}},
		{"done/item.output", done, "item.output", "synthetic result", func(e *realtime.RealtimeServerEventUnion) (any, respjson.Field) {
			v := e.OfConversationItemDone.Item.OfRealtimeMcpToolCall
			return v.Output, v.JSON.Output
		}},
		{"done/item.error", done, "item.error", "synthetic tool error", func(e *realtime.RealtimeServerEventUnion) (any, respjson.Field) {
			v := e.OfConversationItemDone.Item.OfRealtimeMcpToolCall
			return v.Error.OfRealtimeMcpToolExecutionError.Message, v.JSON.Error
		}},
		{"committed/previous_item_id", committed, "previous_item_id", "prev_audio", func(e *realtime.RealtimeServerEventUnion) (any, respjson.Field) {
			v := e.OfInputAudioBufferCommittedEvent
			return v.PreviousItemID, v.JSON.PreviousItemID
		}},
		{"completed/logprobs", completed, "logprobs", []any{[]any{"hello", []int64{104, 101, 108, 108, 111}, -0.25}}, func(e *realtime.RealtimeServerEventUnion) (any, respjson.Field) {
			v := e.OfConversationItemInputAudioTranscriptionCompletedEvent
			return nullableLogprobs(v.Logprobs), v.JSON.Logprobs
		}},
		{"delta/logprobs", delta, "logprobs", []any{[]any{"hi", []int64{104, 105}, -0.5}}, func(e *realtime.RealtimeServerEventUnion) (any, respjson.Field) {
			v := e.OfConversationItemInputAudioTranscriptionDeltaEvent
			return nullableLogprobs(v.Logprobs), v.JSON.Logprobs
		}},
		{"errorEvent/error.code", errorEvent, "error.code", "invalid_event", func(e *realtime.RealtimeServerEventUnion) (any, respjson.Field) {
			v := e.OfRealtimeErrorEvent.Error
			return v.Code, v.JSON.Code
		}},
		{"errorEvent/error.event_id", errorEvent, "error.event_id", "client_17", func(e *realtime.RealtimeServerEventUnion) (any, respjson.Field) {
			v := e.OfRealtimeErrorEvent.Error
			return v.EventID, v.JSON.EventID
		}},
		{"errorEvent/error.param", errorEvent, "error.param", "audio.format", func(e *realtime.RealtimeServerEventUnion) (any, respjson.Field) {
			v := e.OfRealtimeErrorEvent.Error
			return v.Param, v.JSON.Param
		}},
		{"response/response.metadata", response, "response.metadata", map[string]string{"audit": "response_17"}, func(e *realtime.RealtimeServerEventUnion) (any, respjson.Field) {
			v := e.OfResponseDoneEvent.Response
			return map[string]string(v.Metadata), v.JSON.Metadata
		}},
		{"session/session.prompt", session, "session.prompt", "pmpt_audit", func(e *realtime.RealtimeServerEventUnion) (any, respjson.Field) {
			v := e.OfSessionUpdatedEvent.Session.OfRealtimeSessionCreateRequest
			return v.Prompt.ID, v.JSON.Prompt
		}},
		{"session/session.tracing", session, "session.tracing", "auto", func(e *realtime.RealtimeServerEventUnion) (any, respjson.Field) {
			v := e.OfSessionUpdatedEvent.Session.OfRealtimeSessionCreateRequest
			return string(v.Tracing.OfAuto), v.JSON.Tracing
		}},
		{"session/session.audio.input.turn_detection", session, "session.audio.input.turn_detection", "server_vad", func(e *realtime.RealtimeServerEventUnion) (any, respjson.Field) {
			v := e.OfSessionUpdatedEvent.Session.OfRealtimeSessionCreateRequest.Audio.Input
			return v.TurnDetection.Type, v.JSON.TurnDetection
		}},
		{"session/session.audio.input.turn_detection.idle_timeout_ms", session, "session.audio.input.turn_detection.idle_timeout_ms", int64(6000), func(e *realtime.RealtimeServerEventUnion) (any, respjson.Field) {
			v := e.OfSessionUpdatedEvent.Session.OfRealtimeSessionCreateRequest.Audio.Input.TurnDetection.OfRealtimeAudioInputTurnDetectionServerVad
			return v.IdleTimeoutMs, v.JSON.IdleTimeoutMs
		}},
		{"session/session.tools.0.allowed_callers", session, "session.tools.0.allowed_callers", []string{"direct"}, func(e *realtime.RealtimeServerEventUnion) (any, respjson.Field) {
			v := e.OfSessionUpdatedEvent.Session.OfRealtimeSessionCreateRequest.Tools[0].OfRealtimeToolsConfigUnionMcp
			return v.AllowedCallers, v.JSON.AllowedCallers
		}},
		{"session/session.tools.0.allowed_tools", session, "session.tools.0.allowed_tools", []string{"lookup"}, func(e *realtime.RealtimeServerEventUnion) (any, respjson.Field) {
			v := e.OfSessionUpdatedEvent.Session.OfRealtimeSessionCreateRequest.Tools[0].OfRealtimeToolsConfigUnionMcp
			return v.AllowedTools.OfMcpAllowedTools, v.JSON.AllowedTools
		}},
		{"session/session.tools.0.headers", session, "session.tools.0.headers", map[string]string{"X-Synthetic": "audit"}, func(e *realtime.RealtimeServerEventUnion) (any, respjson.Field) {
			v := e.OfSessionUpdatedEvent.Session.OfRealtimeSessionCreateRequest.Tools[0].OfRealtimeToolsConfigUnionMcp
			return v.Headers, v.JSON.Headers
		}},
		{"session/session.tools.0.require_approval", session, "session.tools.0.require_approval", "always", func(e *realtime.RealtimeServerEventUnion) (any, respjson.Field) {
			v := e.OfSessionUpdatedEvent.Session.OfRealtimeSessionCreateRequest.Tools[0].OfRealtimeToolsConfigUnionMcp
			return v.RequireApproval.OfString, v.JSON.RequireApproval
		}},
		{"transcription/session.audio.input.turn_detection", transcription, "session.audio.input.turn_detection", "server_vad", func(e *realtime.RealtimeServerEventUnion) (any, respjson.Field) {
			v := e.OfSessionUpdatedEvent.Session.OfRealtimeTranscriptionSessionCreateRequest.Audio.Input
			return v.TurnDetection.Type, v.JSON.TurnDetection
		}},
		{"transcription/session.audio.input.turn_detection.idle_timeout_ms", transcription, "session.audio.input.turn_detection.idle_timeout_ms", int64(7000), func(e *realtime.RealtimeServerEventUnion) (any, respjson.Field) {
			v := e.OfSessionUpdatedEvent.Session.OfRealtimeTranscriptionSessionCreateRequest.Audio.Input.TurnDetection.OfRealtimeTranscriptionSessionAudioInputTurnDetectionServerVad
			return v.IdleTimeoutMs, v.JSON.IdleTimeoutMs
		}},
		{"arguments/obfuscation", arguments, "obfuscation", "mask_audit", func(e *realtime.RealtimeServerEventUnion) (any, respjson.Field) {
			v := e.OfResponseMcpCallArgumentsDelta
			return v.Obfuscation, v.JSON.Obfuscation
		}},
	}
	p := newRealtimePeer(t, true)
	client := openai.NewClient(option.WithHTTPClient(p.server.Client()), option.WithBaseURL(p.server.URL+"/v1"), option.WithAPIKey("synthetic-key"))
	conn, err := client.Realtime.Connect(p.ctx, realtime.RealtimeConnectionOptions{})
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Abort()
	req := peerValue(t, p.ctx, p.requests)
	if req.Method != http.MethodGet || req.URL.Path != "/v1/realtime" || req.Header.Get("Authorization") != "Bearer synthetic-key" {
		t.Fatal("synthetic Realtime handshake was not exercised")
	}
	ws := peerValue(t, p.ctx, p.opened)
	for _, tc := range cases {
		for _, state := range []string{"value", "null", "omitted"} {
			t.Run(tc.name+"/"+state, func(t *testing.T) {
				writeRealtime(t, p, ws, nullableWire(t, tc.wire, tc.path, state))
				event, err := conn.Recv(p.ctx)
				if err != nil {
					t.Fatal(err)
				}
				got, meta := tc.read(&event)
				assertNullableValue(t, state, got, tc.want, meta)
			})
		}
	}
	if err := conn.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestTranslationNullableElapsedThroughWebsocket(t *testing.T) {
	p := newRealtimePeer(t, true)
	client := openai.NewClient(option.WithHTTPClient(p.server.Client()), option.WithBaseURL(p.server.URL+"/v1"), option.WithAPIKey("synthetic-key"))
	conn, err := client.Realtime.ConnectTranslation(p.ctx, "gpt-realtime-translate", realtime.TranslationConnectionOptions{})
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Abort()
	req := peerValue(t, p.ctx, p.requests)
	if req.Method != http.MethodGet || req.URL.Path != "/v1/realtime/translations" || req.Header.Get("Authorization") != "Bearer synthetic-key" {
		t.Fatal("synthetic Translation handshake was not exercised")
	}
	ws := peerValue(t, p.ctx, p.opened)
	for _, tc := range []struct {
		kind, delta string
		elapsed     int64
		read        func(*realtime.RealtimeTranslationServerEventUnion) (int64, respjson.Field)
	}{
		{"session.input_transcript.delta", "hello", 600, func(e *realtime.RealtimeTranslationServerEventUnion) (int64, respjson.Field) {
			v := e.OfRealtimeTranslationInputTranscriptDeltaEvent
			return v.ElapsedMs, v.JSON.ElapsedMs
		}},
		{"session.output_transcript.delta", "bonjour", 800, func(e *realtime.RealtimeTranslationServerEventUnion) (int64, respjson.Field) {
			v := e.OfRealtimeTranslationOutputTranscriptDeltaEvent
			return v.ElapsedMs, v.JSON.ElapsedMs
		}},
		{"session.output_audio.delta", "AAAA", 1000, func(e *realtime.RealtimeTranslationServerEventUnion) (int64, respjson.Field) {
			v := e.OfRealtimeTranslationOutputAudioDeltaEvent
			return v.ElapsedMs, v.JSON.ElapsedMs
		}},
	} {
		payload, err := json.Marshal(map[string]any{"type": tc.kind, "event_id": "evt_translation", "delta": tc.delta, "elapsed_ms": tc.elapsed})
		if err != nil {
			t.Fatal(err)
		}
		for _, state := range []string{"value", "null", "omitted"} {
			t.Run(tc.kind+"/"+state, func(t *testing.T) {
				writeRealtime(t, p, ws, nullableWire(t, string(payload), "elapsed_ms", state))
				event, err := conn.Recv(p.ctx)
				if err != nil {
					t.Fatal(err)
				}
				got, meta := tc.read(&event)
				assertNullableValue(t, state, got, tc.elapsed, meta)
			})
		}
	}
	if err := conn.Close(); err != nil {
		t.Fatal(err)
	}
}

func nullableLogprobs(values []realtime.LogProbProperties) []any {
	var result []any
	for _, v := range values {
		result = append(result, []any{v.Token, v.Bytes, v.Logprob})
	}
	return result
}

// Change only the field under test, retaining its enclosing objects.
func nullableWire(t *testing.T, wire, path, state string) string {
	t.Helper()
	if state == "value" {
		return wire
	}
	var root any
	if err := json.Unmarshal([]byte(wire), &root); err != nil {
		t.Fatal(err)
	}
	parts := strings.Split(path, ".")
	parent := root
	for _, part := range parts[:len(parts)-1] {
		switch node := parent.(type) {
		case map[string]any:
			parent = node[part]
		case []any:
			index, err := strconv.Atoi(part)
			if err != nil {
				t.Fatal(err)
			}
			parent = node[index]
		default:
			t.Fatalf("invalid fixture path %q", path)
		}
	}
	object := parent.(map[string]any)
	key := parts[len(parts)-1]
	if state == "null" {
		object[key] = nil
	} else {
		delete(object, key)
	}
	result, err := json.Marshal(root)
	if err != nil {
		t.Fatal(err)
	}
	return string(result)
}

func assertNullableValue(t *testing.T, state string, got, want any, meta respjson.Field) {
	t.Helper()
	if state == "value" {
		if !reflect.DeepEqual(got, want) {
			t.Errorf("typed value = %#v, want %#v", got, want)
		}
		if !meta.Valid() {
			t.Error("non-null field is not valid")
		}
		return
	}
	value := reflect.ValueOf(got)
	if value.IsValid() && !value.IsZero() {
		t.Errorf("nullish typed value = %#v, want zero", got)
	}
	expectedRaw := respjson.Null
	if state == "omitted" {
		expectedRaw = respjson.Omitted
	}
	if meta.Valid() || meta.Raw() != expectedRaw {
		t.Errorf("presence = (valid %t, raw %q), want (false, %q)", meta.Valid(), meta.Raw(), expectedRaw)
	}
}

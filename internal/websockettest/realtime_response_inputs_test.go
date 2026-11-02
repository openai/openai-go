package websockettest

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	wire "github.com/coder/websocket"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/packages/websocket"
	"github.com/openai/openai-go/v3/realtime"
)

func connectRealtimeInputPeer(t *testing.T, options realtime.RealtimeConnectionOptions) (context.Context, *realtime.RealtimeConnection, <-chan []byte) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	frames := make(chan []byte, 4)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		peer, err := wire.Accept(w, r, nil)
		if err != nil {
			t.Errorf("Accept() = %v", err)
			return
		}
		defer func() { _ = peer.CloseNow() }()
		for {
			_, data, err := peer.Read(ctx)
			if err != nil {
				return
			}
			select {
			case frames <- data:
			case <-ctx.Done():
				return
			}
		}
	}))
	t.Cleanup(func() { cancel(); server.Close() })
	client := openai.NewClient(option.WithAPIKey("synthetic-key"), option.WithUnsafeAllowHTTP(), option.WithBaseURL(server.URL))
	connection, err := client.Realtime.Connect(ctx, options)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(connection.Abort)
	return ctx, connection, frames
}

func readRealtimeInputFrame(t *testing.T, ctx context.Context, frames <-chan []byte) map[string]any {
	t.Helper()
	select {
	case data := <-frames:
		var event map[string]any
		if err := json.Unmarshal(data, &event); err != nil {
			t.Fatal(err)
		}
		return event
	case <-ctx.Done():
		t.Fatal("peer did not receive response.create")
		return nil
	}
}

func TestRealtimeResponseCreateInputsPreserveEffectiveEvent(t *testing.T) {
	ctx, connection, frames := connectRealtimeInputPeer(t, realtime.RealtimeConnectionOptions{})
	item := realtime.ConversationItemParamOfFunctionCallOutput("call_1", "tool result")
	legacy := realtime.RealtimeResponseCreateParams{Input: []realtime.ConversationItemUnionParam{item}}
	if _, err := json.Marshal(legacy); err != nil {
		t.Fatal(err)
	}
	original := map[string]any{
		"type": "response.create", "event_id": "evt_mixed", "future": map[string]any{"enabled": true},
		"response": map[string]any{
			"instructions": "retain", "future": map[string]any{"count": 3},
			"input": []any{map[string]any{"type": "function_call_output", "call_id": "old_call", "output": "old output"}},
		},
	}
	event := param.Override[realtime.ResponseCreateEventParam](original)
	before, err := json.Marshal(event)
	if err != nil {
		t.Fatal(err)
	}
	inputs := []realtime.ResponseInput{realtime.ResponseInputReference("item_prior-1"), realtime.ResponseInputItem(item)}
	if err = connection.SendResponseCreate(ctx, event, inputs); err != nil {
		t.Fatal(err)
	}
	actual := readRealtimeInputFrame(t, ctx, frames)
	var expected map[string]any
	if err = json.Unmarshal(before, &expected); err != nil {
		t.Fatal(err)
	}
	expected["response"].(map[string]any)["input"] = []any{
		map[string]any{"type": "item_reference", "id": "item_prior-1"},
		map[string]any{"type": "function_call_output", "call_id": "call_1", "output": "tool result"},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("response.create changed effective fields: got %#v, want %#v", actual, expected)
	}
	after, err := json.Marshal(event)
	if err != nil || string(before) != string(after) {
		t.Fatal("SendResponseCreate mutated the caller's event")
	}
}

func TestRealtimeResponseCreateEmptyInputsClearContext(t *testing.T) {
	for _, inputs := range [][]realtime.ResponseInput{nil, {}} {
		ctx, connection, frames := connectRealtimeInputPeer(t, realtime.RealtimeConnectionOptions{})
		if err := connection.SendResponseCreate(ctx, realtime.ResponseCreateEventParam{}, inputs); err != nil {
			t.Fatal(err)
		}
		event := readRealtimeInputFrame(t, ctx, frames)
		if event["type"] != "response.create" || !reflect.DeepEqual(event["response"], map[string]any{"input": []any{}}) {
			t.Fatalf("empty input did not clear context with an array: %#v", event)
		}
	}
}

func requireRealtimeInputNotAttempted(t *testing.T, err error) {
	t.Helper()
	var delivery *websocket.DeliveryError
	if !errors.As(err, &delivery) || delivery.MayHaveBeenSent {
		t.Fatalf("error did not report a never-attempted write: %v", err)
	}
}

func TestRealtimeResponseCreateRejectsMalformedValuesWithoutWriting(t *testing.T) {
	checkRejected := func(t *testing.T, event realtime.ResponseCreateEventParam) {
		t.Helper()
		ctx, connection, frames := connectRealtimeInputPeer(t, realtime.RealtimeConnectionOptions{})
		err := connection.SendResponseCreate(ctx, event, nil)
		requireRealtimeInputNotAttempted(t, err)
		// The next valid command must be the peer's first frame.
		if err := connection.SendResponseCreate(ctx, realtime.ResponseCreateEventParam{}, nil); err != nil {
			t.Fatal(err)
		}
		first := readRealtimeInputFrame(t, ctx, frames)
		if first["type"] != "response.create" || !reflect.DeepEqual(first["response"], map[string]any{"input": []any{}}) {
			t.Fatalf("rejected value reached the peer: %#v", first)
		}
	}
	for _, raw := range []json.RawMessage{json.RawMessage(`null`), json.RawMessage(`[]`), json.RawMessage(`{"type":"response.create","response":null}`), json.RawMessage(`{"type":"response.create","response":[]}`)} {
		t.Run(string(raw), func(t *testing.T) {
			checkRejected(t, param.Override[realtime.ResponseCreateEventParam](raw))
		})
	}
	for _, invalid := range []struct {
		name  string
		event realtime.ResponseCreateEventParam
	}{
		{"nil event map", param.Override[realtime.ResponseCreateEventParam]((*map[string]any)(nil))},
		{"nil response map", param.Override[realtime.ResponseCreateEventParam](map[string]any{"type": "response.create", "response": (*map[string]any)(nil)})},
		{"nil response params", param.Override[realtime.ResponseCreateEventParam](map[string]any{"type": "response.create", "response": (*realtime.RealtimeResponseCreateParams)(nil)})},
	} {
		t.Run(invalid.name, func(t *testing.T) { checkRejected(t, invalid.event) })
	}
	ctx, connection, frames := connectRealtimeInputPeer(t, realtime.RealtimeConnectionOptions{})
	err := connection.SendResponseCreate(ctx, realtime.ResponseCreateEventParam{}, []realtime.ResponseInput{{}})
	requireRealtimeInputNotAttempted(t, err)
	if err := connection.SendResponseCreate(ctx, realtime.ResponseCreateEventParam{}, nil); err != nil {
		t.Fatal(err)
	}
	if first := readRealtimeInputFrame(t, ctx, frames); !reflect.DeepEqual(first["response"], map[string]any{"input": []any{}}) {
		t.Fatalf("uninitialized input reached the peer: %#v", first)
	}
}

type realtimeInputEncodingProbe struct {
	calls   atomic.Int32
	started chan struct{}
	release <-chan struct{}
	value   json.RawMessage
	err     error
}

func (probe *realtimeInputEncodingProbe) MarshalJSON() ([]byte, error) {
	if probe.calls.Add(1) == 1 && probe.started != nil {
		close(probe.started)
	}
	if probe.release != nil {
		<-probe.release
	}
	return probe.value, probe.err
}

func TestRealtimeResponseCreateChecksAdmissionBeforeEncoding(t *testing.T) {
	for _, closed := range []bool{false, true} {
		ctx, connection, _ := connectRealtimeInputPeer(t, realtime.RealtimeConnectionOptions{})
		probe := &realtimeInputEncodingProbe{value: json.RawMessage(`{"type":"response.create"}`)}
		callCtx, cancel := context.WithCancel(ctx)
		want := context.Canceled
		if closed {
			connection.Abort()
			want = websocket.ErrClosed
		} else {
			cancel()
		}
		err := connection.SendResponseCreate(callCtx, param.Override[realtime.ResponseCreateEventParam](probe), nil)
		cancel()
		if !errors.Is(err, want) || probe.calls.Load() != 0 {
			t.Fatalf("rejected send encoded its event: closed=%t, calls=%d, err=%v", closed, probe.calls.Load(), err)
		}
		requireRealtimeInputNotAttempted(t, err)
	}
}

func TestRealtimeResponseCreateRetainsPendingSendLimit(t *testing.T) {
	ctx, connection, frames := connectRealtimeInputPeer(t, realtime.RealtimeConnectionOptions{MaxPendingSends: 1})
	release := make(chan struct{})
	var once sync.Once
	unblock := func() { once.Do(func() { close(release) }) }
	t.Cleanup(unblock)
	first := &realtimeInputEncodingProbe{started: make(chan struct{}), release: release, value: json.RawMessage(`{"type":"response.create","event_id":"evt_first"}`)}
	done := make(chan error, 1)
	go func() {
		done <- connection.SendResponseCreate(ctx, param.Override[realtime.ResponseCreateEventParam](first), nil)
	}()
	select {
	case <-first.started:
	case <-ctx.Done():
		t.Fatal("first command did not begin encoding")
	}
	second := &realtimeInputEncodingProbe{value: json.RawMessage(`{"type":"response.create","event_id":"evt_rejected"}`)}
	err := connection.SendResponseCreate(ctx, param.Override[realtime.ResponseCreateEventParam](second), nil)
	if !errors.Is(err, websocket.ErrPendingLimit) || second.calls.Load() != 0 {
		t.Fatalf("full pending queue admitted encoding: calls=%d, err=%v", second.calls.Load(), err)
	}
	requireRealtimeInputNotAttempted(t, err)
	unblock()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-ctx.Done():
		t.Fatal("admitted command did not complete")
	}
	if first.calls.Load() != 1 || readRealtimeInputFrame(t, ctx, frames)["event_id"] != "evt_first" {
		t.Fatal("admitted command was changed or encoded repeatedly")
	}
}

type realtimeInputOpaqueMap map[string]*realtimeInputEncodingProbe

func (value realtimeInputOpaqueMap) MarshalJSON() ([]byte, error) {
	return value["input"].MarshalJSON()
}

func TestRealtimeResponseCreateEncodingFailureLeavesConnectionUsable(t *testing.T) {
	for _, carrier := range []string{"event", "response", "unrelated response extra", "defined response map"} {
		ctx, connection, frames := connectRealtimeInputPeer(t, realtime.RealtimeConnectionOptions{})
		cause := errors.New("synthetic input encoder failure")
		probe := &realtimeInputEncodingProbe{err: cause}
		event := param.Override[realtime.ResponseCreateEventParam](probe)
		if carrier == "response" {
			event = realtime.ResponseCreateEventParam{Response: param.Override[realtime.RealtimeResponseCreateParams](probe)}
		} else if carrier == "defined response map" {
			event = realtime.ResponseCreateEventParam{Response: param.Override[realtime.RealtimeResponseCreateParams](realtimeInputOpaqueMap{"input": probe})}
		} else if carrier == "unrelated response extra" {
			event = realtime.ResponseCreateEventParam{}
			event.Response.SetExtraFields(map[string]any{"future": probe})
		}
		err := connection.SendResponseCreate(ctx, event, nil)
		if !errors.Is(err, cause) || probe.calls.Load() != 1 {
			t.Fatalf("serialization cause was lost: carrier=%s, calls=%d, err=%v", carrier, probe.calls.Load(), err)
		}
		requireRealtimeInputNotAttempted(t, err)
		if err := connection.SendResponseCreate(ctx, realtime.ResponseCreateEventParam{}, nil); err != nil {
			t.Fatal(err)
		}
		if first := readRealtimeInputFrame(t, ctx, frames); !reflect.DeepEqual(first["response"], map[string]any{"input": []any{}}) {
			t.Fatalf("failed serialization wrote a command: %#v", first)
		}
	}
}

func TestRealtimeResponseCreateCancellationDuringEncodingDoesNotWrite(t *testing.T) {
	ctx, connection, frames := connectRealtimeInputPeer(t, realtime.RealtimeConnectionOptions{})
	callCtx, cancel := context.WithCancel(ctx)
	t.Cleanup(cancel)
	release := make(chan struct{})
	var once sync.Once
	unblock := func() { once.Do(func() { close(release) }) }
	t.Cleanup(unblock)
	probe := &realtimeInputEncodingProbe{started: make(chan struct{}), release: release, value: json.RawMessage(`{"type":"response.create","event_id":"evt_canceled"}`)}
	done := make(chan error, 1)
	go func() {
		done <- connection.SendResponseCreate(callCtx, param.Override[realtime.ResponseCreateEventParam](probe), nil)
	}()
	select {
	case <-probe.started:
	case <-ctx.Done():
		t.Fatal("command did not begin encoding")
	}
	cancel()
	unblock()
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("cancellation was lost: %v", err)
		}
		requireRealtimeInputNotAttempted(t, err)
	case <-ctx.Done():
		t.Fatal("canceled command did not return")
	}
	if err := connection.SendResponseCreate(ctx, realtime.ResponseCreateEventParam{}, nil); err != nil {
		t.Fatal(err)
	}
	if first := readRealtimeInputFrame(t, ctx, frames); first["event_id"] != nil {
		t.Fatalf("canceled command reached the peer: %#v", first)
	}
}

func TestRealtimeResponseCreateSkipsSupersededTypedInput(t *testing.T) {
	for _, name := range []string{
		"multiple legacy variants", "unused input encoder", "response extra input",
		"response override map", "response override map pointer", "event override map",
		"event override map pointer", "event extra response map", "event extra response map pointer",
		"event extra typed response", "event extra typed response pointer",
	} {
		t.Run(name, func(t *testing.T) {
			ctx, connection, frames := connectRealtimeInputPeer(t, realtime.RealtimeConnectionOptions{})
			probe := &realtimeInputEncodingProbe{err: errors.New("unused input encoder")}
			makeEvent := func() realtime.ResponseCreateEventParam {
				oldItem := realtime.ConversationItemParamOfFunctionCallOutput("old_call", "old output")
				oldItem.OfRealtimeUserMessageItem = &realtime.RealtimeConversationItemUserMessageParam{}
				if name == "unused input encoder" {
					oldItem = param.Override[realtime.ConversationItemUnionParam](probe)
				}
				response := realtime.RealtimeResponseCreateParams{
					Instructions: param.NewOpt("retain"), Input: []realtime.ConversationItemUnionParam{oldItem},
				}
				response.SetExtraFields(map[string]any{"future": map[string]any{"count": 3}})
				responseMap := map[string]any{"instructions": "retain", "input": probe, "future": map[string]any{"count": 3}}
				eventMap := map[string]any{"type": "response.create", "event_id": "evt_replaced", "response": responseMap, "future": map[string]any{"enabled": true}}
				event := realtime.ResponseCreateEventParam{EventID: param.NewOpt("evt_replaced"), Response: response}
				event.SetExtraFields(map[string]any{"future": map[string]any{"enabled": true}})
				switch name {
				case "response extra input":
					event.Response.SetExtraFields(map[string]any{"input": probe, "future": map[string]any{"count": 3}, "omitted": param.Omit})
				case "response override map":
					event.Response = param.Override[realtime.RealtimeResponseCreateParams](responseMap)
				case "response override map pointer":
					event.Response = param.Override[realtime.RealtimeResponseCreateParams](&responseMap)
				case "event override map":
					event = param.Override[realtime.ResponseCreateEventParam](eventMap)
				case "event override map pointer":
					event = param.Override[realtime.ResponseCreateEventParam](&eventMap)
				case "event extra response map":
					event.SetExtraFields(map[string]any{"response": responseMap, "future": map[string]any{"enabled": true}})
				case "event extra response map pointer":
					event.SetExtraFields(map[string]any{"response": &responseMap, "future": map[string]any{"enabled": true}})
				case "event extra typed response":
					event.SetExtraFields(map[string]any{"response": response, "future": map[string]any{"enabled": true}})
				case "event extra typed response pointer":
					event.SetExtraFields(map[string]any{"response": &response, "future": map[string]any{"enabled": true}})
				}
				return event
			}
			event, expectedCaller := makeEvent(), makeEvent()
			if sendErr := connection.SendResponseCreate(ctx, event, []realtime.ResponseInput{realtime.ResponseInputReference("item_prior-1")}); sendErr != nil {
				t.Fatalf("superseded input blocked a valid replacement: %v", sendErr)
			}
			actual := readRealtimeInputFrame(t, ctx, frames)
			want := map[string]any{
				"type": "response.create", "event_id": "evt_replaced", "future": map[string]any{"enabled": true},
				"response": map[string]any{
					"instructions": "retain", "future": map[string]any{"count": float64(3)},
					"input": []any{map[string]any{"type": "item_reference", "id": "item_prior-1"}},
				},
			}
			if !reflect.DeepEqual(actual, want) || probe.calls.Load() != 0 {
				t.Fatalf("wrong replacement or unused encoding: event=%#v, calls=%d", actual, probe.calls.Load())
			}
			if !reflect.DeepEqual(event, expectedCaller) {
				t.Fatal("helper changed caller input or metadata")
			}
		})
	}
}

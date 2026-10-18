package realtime_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	wire "github.com/coder/websocket"
	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/azure"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/respjson"
	"github.com/openai/openai-go/v3/packages/websocket"
	"github.com/openai/openai-go/v3/realtime"
)

type realtimePeer struct {
	server   *httptest.Server
	ctx      context.Context
	opened   chan *wire.Conn
	requests chan *http.Request
	messages chan []byte
}

func newRealtimePeer(t *testing.T, tls bool) *realtimePeer {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	p := &realtimePeer{
		ctx: ctx, opened: make(chan *wire.Conn, 4),
		requests: make(chan *http.Request, 4), messages: make(chan []byte, 16),
	}
	handle := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p.requests <- r.Clone(context.Background())
		ws, err := wire.Accept(w, r, nil)
		if err != nil {
			t.Errorf("accept: %v", err)
			return
		}
		defer func() { _ = ws.CloseNow() }()
		p.opened <- ws
		for {
			_, data, err := ws.Read(ctx)
			if err != nil {
				return
			}
			select {
			case p.messages <- data:
			case <-ctx.Done():
				return
			}
		}
	})
	if tls {
		p.server = httptest.NewTLSServer(handle)
	} else {
		p.server = httptest.NewServer(handle)
	}
	t.Cleanup(p.server.Close)
	t.Cleanup(cancel)
	return p
}

func peerValue[T any](t *testing.T, ctx context.Context, ch <-chan T) T {
	t.Helper()
	select {
	case value := <-ch:
		return value
	case <-ctx.Done():
		t.Fatalf("waiting for peer: %v", ctx.Err())
		var zero T
		return zero
	}
}

func writeRealtime(t *testing.T, p *realtimePeer, ws *wire.Conn, text string) {
	t.Helper()
	if err := ws.Write(p.ctx, wire.MessageText, []byte(text)); err != nil {
		t.Fatalf("peer write: %v", err)
	}
}

func TestRealtimeWebsocketRecoveryRejectsLateSessionSelector(t *testing.T) {
	for _, tc := range []struct {
		name      string
		reconnect bool
	}{
		{name: "middleware"},
		{name: "same-origin redirect on reconnect", reconnect: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			var requests, selected atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Query().Has("intent") || r.URL.Query().Has("call_id") {
					selected.Add(1)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				if requests.Add(1) == 1 {
					socket, err := wire.Accept(w, r, nil)
					if err != nil {
						t.Errorf("accept: %v", err)
						return
					}
					defer func() { _ = socket.CloseNow() }()
					_, _, _ = socket.Read(ctx)
					return
				}
				http.Redirect(w, r, "/realtime", http.StatusTemporaryRedirect)
			}))
			defer server.Close()
			httpClient := server.Client()
			opts := []option.RequestOption{option.WithBaseURL(server.URL), option.WithAPIKey("fake-key"), option.WithUnsafeAllowHTTP(), option.WithHTTPClient(httpClient), option.WithWebSocketRecovery(websocket.RecoveryOptions{Enabled: true})}
			if tc.reconnect {
				httpClient.CheckRedirect = func(r *http.Request, _ []*http.Request) error {
					r.URL.RawQuery = "call_id=selected"
					return nil
				}
			} else {
				opts = append(opts, option.WithMiddleware(func(r *http.Request, next option.MiddlewareNext) (*http.Response, error) {
					r.URL.RawQuery = "intent=transcription"
					return next(r)
				}))
			}
			client := openai.NewClient(opts...)
			conn, err := client.Realtime.Connect(ctx, realtime.RealtimeConnectionOptions{})
			if tc.reconnect {
				if err != nil {
					t.Fatalf("initial connection: %v", err)
				}
				defer conn.Abort()
				conn, err = conn.Reconnect(ctx)
			}
			if conn != nil {
				conn.Abort()
			}
			if err == nil || !strings.Contains(err.Error(), "automatic recovery does not support this session selector") || selected.Load() != 0 {
				t.Fatalf("connection error = %v, selected endpoint requests = %d", err, selected.Load())
			}
		})
	}
}

func TestRealtimeWebsocketSelectedHandshakeAndTypedSession(t *testing.T) {
	for _, mode := range []struct{ name, key, value string }{
		{"model", "model", "gpt-realtime"},
		{"transcription", "intent", "transcription"},
		{"call", "call_id", "rtc_fake + /?"},
	} {
		t.Run(mode.name, func(t *testing.T) {
			p := newRealtimePeer(t, true)
			client := openai.NewClient(
				option.WithHTTPClient(p.server.Client()),
				option.WithBaseURL(strings.Replace(p.server.URL, "https:", "wss:", 1)+"/proxy/v1/?keep=a%2Bb&multi=base1&multi=base2&"+mode.key+"=old"),
				option.WithAPIKey("fake-client-key"),
				option.WithHeader("X-Selection", "client"),
			)
			client.Realtime.Options = append(client.Realtime.Options, option.WithHeader("X-Selection", "service"))
			openCtx, cancel := context.WithCancel(p.ctx)
			conn, err := client.Realtime.Connect(openCtx, realtime.RealtimeConnectionOptions{},
				option.WithQuery(mode.key, mode.value),
				option.WithQueryAdd("multi", "request"),
				option.WithHeader("X-Selection", "request"),
				option.WithHeader("Authorization", "Bearer fake-selected"))
			if err != nil {
				cancel()
				t.Fatal(err)
			}
			defer func() { _ = conn.Close() }()
			cancel() // Opening context cannot end an established Realtime session.
			req := peerValue(t, p.ctx, p.requests)
			ws := peerValue(t, p.ctx, p.opened)
			if req.Method != http.MethodGet || req.URL.Path != "/proxy/v1/realtime" ||
				req.URL.Query().Get(mode.key) != mode.value || req.URL.Query().Get("keep") != "a+b" ||
				!reflect.DeepEqual(req.URL.Query()["multi"], []string{"base1", "base2", "request"}) ||
				req.Header.Get("Authorization") != "Bearer fake-selected" || req.Header.Get("X-Selection") != "request" {
				t.Fatalf("unexpected synthetic handshake: method=%s path=%s query=%v selected=%s",
					req.Method, req.URL.Path, req.URL.Query(), req.Header.Get("X-Selection"))
			}
			// The helper does not manufacture a session update at connect.
			select {
			case msg := <-p.messages:
				t.Fatalf("unexpected automatic command: %q", msg)
			default:
			}
			sessionType := "realtime"
			sessionUpdate := realtime.SessionUpdateEventSessionUnionParam{
				OfRealtime: &realtime.RealtimeSessionCreateRequestParam{Instructions: openai.String("synthetic")},
			}
			if mode.name == "transcription" {
				sessionType = "transcription"
				sessionUpdate = realtime.SessionUpdateEventSessionUnionParam{
					OfTranscription: &realtime.RealtimeTranscriptionSessionCreateRequestParam{},
				}
			}
			writeRealtime(t, p, ws, fmt.Sprintf(`{"type":"session.created","event_id":"evt_created","session":{"id":"sess_test","type":%q,"model":"gpt-realtime"}}`, sessionType))
			event, err := conn.Recv(p.ctx)
			if err != nil || event.Type != "session.created" || event.AsSessionCreatedEvent().EventID != "evt_created" {
				t.Fatalf("typed session.created: type=%s error=%v", event.Type, err)
			}
			createdSession := event.AsSessionCreatedEvent().Session
			if sessionType == "transcription" {
				if string(createdSession.AsRealtimeTranscriptionSessionCreateRequest().Type) != sessionType {
					t.Fatalf("lost typed transcription session %s", createdSession.RawJSON())
				}
			} else if string(createdSession.AsRealtimeSessionCreateRequest().Type) != sessionType {
				t.Fatalf("lost typed Realtime session %s", createdSession.RawJSON())
			}
			if err := conn.Send(p.ctx, realtime.RealtimeClientEventUnionParam{
				OfSessionUpdate: &realtime.SessionUpdateEventParam{Session: sessionUpdate},
			}); err != nil {
				t.Fatal(err)
			}
			var update struct {
				Type    string `json:"type"`
				Session struct {
					Type string `json:"type"`
				} `json:"session"`
			}
			if err := json.Unmarshal(peerValue(t, p.ctx, p.messages), &update); err != nil ||
				update.Type != "session.update" || update.Session.Type != sessionType {
				t.Fatalf("typed session update: %#v error=%v", update, err)
			}
			appendAudio := realtime.RealtimeClientEventUnionParam{
				OfInputAudioBufferAppend: &realtime.InputAudioBufferAppendEventParam{Audio: "AQID", EventID: openai.String("client_1")},
			}
			if err := conn.Send(p.ctx, appendAudio); err != nil {
				t.Fatal(err)
			}
			var sent map[string]any
			if err := json.Unmarshal(peerValue(t, p.ctx, p.messages), &sent); err != nil {
				t.Fatal(err)
			}
			if sent["type"] != "input_audio_buffer.append" || sent["audio"] != "AQID" || sent["event_id"] != "client_1" {
				t.Fatalf("typed input audio: %v", sent)
			}
			// Protocol errors remain in order and a response terminal never ends a session.
			for _, frame := range []struct{ eventType, text string }{
				{"error", `{"type":"error","event_id":"evt_error","error":{"type":"invalid_request_error","message":"synthetic recoverable error","code":"invalid_event","event_id":"client_1"}}`},
				{"response.done", `{"type":"response.done","event_id":"evt_done","response":{"id":"resp_1","status":"completed"}}`},
				{"session.updated", fmt.Sprintf(`{"type":"session.updated","event_id":"evt_updated","session":{"id":"sess_test","type":%q,"model":"gpt-realtime"}}`, sessionType)},
				{"new.realtime.event", `{"type":"new.realtime.event","event_id":"evt_future","nested":{"missing":null},"count":0}`},
			} {
				writeRealtime(t, p, ws, frame.text)
				got, err := conn.Recv(p.ctx)
				if err != nil || got.Type != frame.eventType || got.RawJSON() != frame.text {
					t.Fatalf("wire event changed: %q, error=%v", got.RawJSON(), err)
				}
				switch got.Type {
				case "error":
					if got.AsRealtimeErrorEvent().Error.Message != "synthetic recoverable error" {
						t.Fatalf("lost typed error %s", got.RawJSON())
					}
				case "response.done":
					if got.AsResponseDoneEvent().EventID != "evt_done" {
						t.Fatalf("lost typed terminal %s", got.RawJSON())
					}
				case "session.updated":
					if got.AsSessionUpdatedEvent().EventID != "evt_updated" {
						t.Fatalf("lost typed session %s", got.RawJSON())
					}
				}
			}
			if err := conn.Send(p.ctx, realtime.RealtimeClientEventUnionParam{OfInputAudioBufferCommit: &realtime.InputAudioBufferCommitEventParam{}}); err != nil {
				t.Fatalf("session must remain usable after response.done: %v", err)
			}
			if string(peerValue(t, p.ctx, p.messages)) != `{"type":"input_audio_buffer.commit"}` {
				t.Error("commit was not serialized as the typed event")
			}
		})
	}
}

func TestRealtimeWebsocketCancellationAndTerminalDrain(t *testing.T) {
	p := newRealtimePeer(t, false)
	client := openai.NewClient(option.WithBaseURL(p.server.URL+"/v1/"), option.WithAPIKey("fake-key"), option.WithUnsafeAllowHTTP())
	conn, err := client.Realtime.Connect(p.ctx, realtime.RealtimeConnectionOptions{})
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()
	ws := peerValue(t, p.ctx, p.opened)
	cancelCtx, cancel := context.WithCancel(p.ctx)
	read := make(chan error, 1)
	go func() { _, recvErr := conn.Recv(cancelCtx); read <- recvErr }()
	cancel()
	if err = peerValue(t, p.ctx, read); !errors.Is(err, context.Canceled) {
		t.Fatalf("canceled wait: %v", err)
	}
	large := strings.Repeat("a", 5<<20)
	text := `{"type":"response.audio.delta","event_id":"evt_audio","delta":"` + large + `","response_id":"resp_1","item_id":"item_1","output_index":0,"content_index":0}`
	writeRealtime(t, p, ws, text)
	writeRealtime(t, p, ws, `{"type":"response.done","event_id":"evt_tail","response":{"id":"resp_1","status":"completed"}}`)
	if err = ws.Close(wire.StatusNormalClosure, "fixture end"); err != nil {
		t.Fatal(err)
	}
	audio, err := conn.Recv(p.ctx)
	if err != nil || audio.AsResponseAudioDeltaEvent().Delta != large || audio.RawJSON() != text {
		t.Fatalf("cancellation, default audio size or final drain: type=%s error=%v", audio.Type, err)
	}
	tail, err := conn.Recv(p.ctx)
	if err != nil || tail.AsResponseDoneEvent().EventID != "evt_tail" {
		t.Fatalf("terminal lost after close: type=%s error=%v", tail.Type, err)
	}
	if _, err := conn.Recv(p.ctx); err == nil {
		t.Fatal("expected EOF after wire close")
	}
	select {
	case <-p.requests:
	default:
		t.Fatal("opening not observed")
	}
	select {
	case <-p.requests:
		t.Fatal("auto reconnected after the peer ended")
	default:
	}
}

func TestRealtimeWebsocketInvalidEnvelopesAndAbort(t *testing.T) {
	for _, invalid := range []string{`null`, `{"type":`, `{"data":"synthetic"}`, `{"type":42}`, `{"type":null}`} {
		t.Run(invalid, func(t *testing.T) {
			p := newRealtimePeer(t, false)
			client := openai.NewClient(option.WithBaseURL(p.server.URL+"/v1/"), option.WithAPIKey("fake-key"), option.WithUnsafeAllowHTTP())
			conn, err := client.Realtime.Connect(p.ctx, realtime.RealtimeConnectionOptions{})
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = conn.Close() }()
			writeRealtime(t, p, peerValue(t, p.ctx, p.opened), invalid)
			if _, err := conn.Recv(p.ctx); err == nil {
				t.Fatal("accepted invalid envelope")
			}
		})
	}
	t.Run("abort pending wait", func(t *testing.T) {
		p := newRealtimePeer(t, false)
		client := openai.NewClient(option.WithBaseURL(p.server.URL+"/v1/"), option.WithAPIKey("fake-key"), option.WithUnsafeAllowHTTP())
		conn, err := client.Realtime.Connect(p.ctx, realtime.RealtimeConnectionOptions{})
		if err != nil {
			t.Fatal(err)
		}
		read := make(chan error, 1)
		go func() { _, err := conn.Recv(p.ctx); read <- err }()
		conn.Abort()
		if err := peerValue(t, p.ctx, read); !errors.Is(err, websocket.ErrClosed) {
			t.Fatalf("abort did not release pending read: %v", err)
		}
		if err := conn.Send(p.ctx, realtime.RealtimeClientEventUnionParam{OfInputAudioBufferCommit: &realtime.InputAudioBufferCommitEventParam{}}); err == nil {
			t.Fatal("send succeeded after Abort")
		}
	})
}

func TestRealtimeWebsocketRejectsUnsupportedProviderBeforeNetwork(t *testing.T) {
	var touched atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		touched.Add(1)
		w.WriteHeader(http.StatusForbidden)
	}))
	defer srv.Close()
	client := openai.NewClient(
		option.WithAPIKey("fake-ambient-openai"),
		azure.WithEndpoint(srv.URL, "fake-version"),
	)
	conn, err := client.Realtime.Connect(context.Background(), realtime.RealtimeConnectionOptions{})
	if err == nil || conn != nil || touched.Load() != 0 || strings.Contains(err.Error(), "fake-ambient-openai") {
		t.Fatalf("unsupported provider leaked request or credential: conn=%v err=%v requests=%d", conn != nil, err, touched.Load())
	}
}

// Every generated GA variant must select its existing Of field. The payload is
// intentionally forward-compatible: preserve unknown fields and invalid typed
// nested data, including inline presence metadata, as the general decoder does.
func TestRealtimeWebsocketReceivesEveryGeneratedVariant(t *testing.T) {
	p := newRealtimePeer(t, true)
	client := openai.NewClient(option.WithHTTPClient(p.server.Client()), option.WithBaseURL(p.server.URL), option.WithAPIKey("fake"))
	conn, err := client.Realtime.Connect(p.ctx, realtime.RealtimeConnectionOptions{})
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Abort()
	ws := peerValue(t, p.ctx, p.opened)
	uType := reflect.TypeFor[realtime.RealtimeServerEventUnion]()
	for i := range uType.NumField() {
		f := uType.Field(i)
		if !strings.HasPrefix(f.Name, "Of") {
			continue
		}
		t.Run(f.Name, func(t *testing.T) {
			typeField, ok := f.Type.FieldByName("Type")
			if !ok || typeField.Tag.Get("default") == "" {
				t.Fatal("Realtime variant has no type discriminator")
			}
			eventType := typeField.Tag.Get("default")
			data := fmt.Sprintf(`{"type":%q,"event_id":"evt_fake","delta":"hello","content_index":0,"output_index":0,"item_id":"item_fake","response_id":"resp_fake","session":{"id":"sess_fake","type":"realtime","future_session":true},"error":{"type":"invalid_request_error","code":"fake_code","message":"fake_message","param":null,"future_error":[1]},"future_field":{"n":1},"nullable_future":null}`, eventType)
			var general realtime.RealtimeServerEventUnion
			if err := json.Unmarshal([]byte(data), &general); err != nil {
				t.Fatal(err)
			}
			writeRealtime(t, p, ws, data)
			event, err := conn.Recv(p.ctx)
			if err != nil {
				t.Fatal(err)
			}
			if event.Type != eventType || event.RawJSON() != data || event.JSON.Type.Raw() != general.JSON.Type.Raw() {
				t.Fatalf("wire envelope not preserved for %s", eventType)
			}
			value := reflect.ValueOf(event)
			metadata := value.FieldByName("JSON")
			wantValue := reflect.ValueOf(general)
			for j := range uType.NumField() {
				field := uType.Field(j)
				if !strings.HasPrefix(field.Name, "Of") {
					continue
				}
				actual := value.Field(j)
				if j == i {
					if !reflect.DeepEqual(actual.Interface(), wantValue.Field(j).Interface()) {
						t.Fatalf("%s lost general decoder values or JSON metadata/ExtraFields", field.Name)
					}
					meta := metadata.FieldByName(field.Name).Interface().(respjson.Field)
					wantMeta := wantValue.FieldByName("JSON").FieldByName(field.Name).Interface().(respjson.Field)
					if !reflect.DeepEqual(meta, wantMeta) {
						t.Fatalf("%s lost inline presence metadata", field.Name)
					}
				} else if !actual.IsZero() {
					t.Fatalf("unselected %s was populated", field.Name)
				}
			}
		})
	}
}

func TestRealtimeWebsocketUnknownAndMalformedEventMetadata(t *testing.T) {
	p := newRealtimePeer(t, true)
	client := openai.NewClient(option.WithHTTPClient(p.server.Client()), option.WithBaseURL(p.server.URL), option.WithAPIKey("fake"))
	conn, err := client.Realtime.Connect(p.ctx, realtime.RealtimeConnectionOptions{})
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Abort()
	ws := peerValue(t, p.ctx, p.opened)
	for _, data := range []string{
		`{"type":"future.event","delta":"unrelated","future":{"nested":[false,null,3]}}`,
		`{"type":"response.done","event_id":"evt_invalid","response":42,"future_field":null}`,
		`{"type":"error","event_id":"evt_error","error":{"type":"invalid_request_error","code":"fake_code","message":"fake_message","param":null,"future_error":{"n":1}},"future_field":true}`,
	} {
		var general realtime.RealtimeServerEventUnion
		if err := json.Unmarshal([]byte(data), &general); err != nil {
			t.Fatal(err)
		}
		writeRealtime(t, p, ws, data)
		event, err := conn.Recv(p.ctx)
		if err != nil || event.RawJSON() != data || event.Type != general.Type || !reflect.DeepEqual(event.AsAny(), general.AsAny()) {
			t.Fatalf("event or malformed-field metadata changed: %v", err)
		}
	}
	writeRealtime(t, p, ws, `{"type":42}`)
	if _, err := conn.Recv(p.ctx); err == nil {
		t.Fatal("non-string event type accepted")
	}
}

func TestRealtimeWebsocketConversationItemRoles(t *testing.T) {
	p := newRealtimePeer(t, true)
	client := openai.NewClient(option.WithHTTPClient(p.server.Client()), option.WithBaseURL(p.server.URL), option.WithAPIKey("fake"))
	conn, err := client.Realtime.Connect(p.ctx, realtime.RealtimeConnectionOptions{})
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Abort()
	ws := peerValue(t, p.ctx, p.opened)
	var reused realtime.ConversationItemUnion
	for _, role := range []string{"system", "user", "assistant"} {
		for _, empty := range []bool{false, true} {
			contentType := "input_text"
			if role == "assistant" {
				contentType = "output_text"
			}
			content := `[{"type":"` + contentType + `","text":"synthetic"}]`
			if empty {
				content = `[]`
			}
			item := `{"type":"message","role":"` + role + `","content":` + content + `}`
			if err := json.Unmarshal([]byte(item), &reused); err != nil {
				t.Fatal(err)
			}
			writeRealtime(t, p, ws, `{"type":"conversation.item.created","item":`+item+`}`)
			event, err := conn.Recv(p.ctx)
			if err != nil {
				t.Fatal(err)
			}
			writeRealtime(t, p, ws, `{"type":"response.done","response":{"output":[`+item+`]}}`)
			response, err := conn.Recv(p.ctx)
			if err != nil {
				t.Fatal(err)
			}
			output := response.AsResponseDoneEvent().Response.Output
			if len(output) != 1 {
				t.Fatalf("response output has %d items, want one", len(output))
			}
			want := map[string]string{
				"system":    "OfRealtimeConversationItemSystemMessage",
				"user":      "OfRealtimeConversationItemUserMessage",
				"assistant": "OfRealtimeConversationItemAssistantMessage",
			}[role]
			for _, got := range []realtime.ConversationItemUnion{reused, event.AsConversationItemCreatedEvent().Item, output[0]} {
				if got.RawJSON() != item {
					t.Errorf("%s raw item changed", role)
				}
				if raw, ok := got.ToParam().Overrides(); !ok || !reflect.DeepEqual(raw, json.RawMessage(item)) {
					t.Errorf("%s ToParam lost the original item", role)
				}
				if got.AsRealtimeConversationItemUserMessage().RawJSON() != item {
					t.Errorf("%s explicit projection lost the original item", role)
				}
				value, metadata := reflect.ValueOf(got), reflect.ValueOf(got.JSON)
				for i := 0; i < value.NumField(); i++ {
					name := value.Type().Field(i).Name
					if !strings.HasPrefix(name, "Of") {
						continue
					}
					present := metadata.FieldByName(name).Interface().(respjson.Field).Valid()
					if present != (name == want) || (name != want && !value.Field(i).IsZero()) {
						t.Errorf("%s empty=%v: %s presence=%v, want only %s", role, empty, name, present, want)
					}
				}
			}
		}
	}
}

func TestRealtimeWebsocketInvalidJSONPreservesSyntaxError(t *testing.T) {
	for _, tc := range []struct{ name, wire string }{
		{"truncated", `{"type":"response.done","response":`},
		{"trailing", `{"type":"response.done"} {"trailing":true}`},
		{"bad nested value", `{"type":"error","error":{"message":undefined}}`},
		{"bad escape", `{"type":"response.output_text.delta","delta":"\q"}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			p := newRealtimePeer(t, true)
			client := openai.NewClient(option.WithHTTPClient(p.server.Client()), option.WithBaseURL(p.server.URL), option.WithAPIKey("fake"))
			conn, err := client.Realtime.Connect(p.ctx, realtime.RealtimeConnectionOptions{})
			if err != nil {
				t.Fatal(err)
			}
			defer conn.Abort()
			ws := peerValue(t, p.ctx, p.opened)
			var before realtime.RealtimeServerEventUnion
			want := json.Unmarshal([]byte(tc.wire), &before)
			writeRealtime(t, p, ws, tc.wire)
			_, err = conn.Recv(p.ctx)
			var wantSyntax, gotSyntax *json.SyntaxError
			if !errors.As(want, &wantSyntax) || !errors.As(err, &gotSyntax) ||
				wantSyntax.Offset != gotSyntax.Offset || wantSyntax.Error() != gotSyntax.Error() {
				t.Fatalf("Recv lost existing SyntaxError: got %v, want %v", err, want)
			}
		})
	}
}

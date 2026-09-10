package openai_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/ssestream"
)

type agentHelperPost struct {
	body   map[string]any
	header http.Header
}
type agentHelperServer struct {
	server     *httptest.Server
	mu         sync.Mutex
	posts      []agentHelperPost
	subscribed bool
	status     string
	events     []string
	eof        bool
	failInput  bool
	postStatus int
	postError  string
	failures   int
	reads      int
	closed     atomic.Int32
	drop       atomic.Bool
}
type agentHelperTransport struct {
	base http.RoundTripper
	mock *agentHelperServer
}
type agentHelperBody struct {
	io.ReadCloser
	closed *atomic.Int32
}

func (b *agentHelperBody) Close() error { b.closed.Add(1); return b.ReadCloser.Close() }
func (tr agentHelperTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	res, err := tr.base.RoundTrip(r)
	if err != nil {
		return nil, err
	}
	if r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/events") {
		res.Body = &agentHelperBody{ReadCloser: res.Body, closed: &tr.mock.closed}
	}
	if r.Method == http.MethodPost && tr.mock.drop.CompareAndSwap(true, false) {
		_, _ = io.Copy(io.Discard, res.Body)
		_ = res.Body.Close()
		return nil, errors.New("synthetic accepted response lost")
	}
	return res, nil
}
func newAgentHelperServer(t *testing.T, events ...string) (*agentHelperServer, openai.Client) {
	t.Helper()
	m := &agentHelperServer{status: "idle", events: events}
	m.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && !strings.HasSuffix(r.URL.Path, "/events") {
			m.mu.Lock()
			m.reads++
			status := m.status
			m.mu.Unlock()
			w.Header().Set("Content-Type", "application/json")
			_, _ = fmt.Fprintf(w, `{"id":"session","status":%q}`, status)
			return
		}
		if r.Method == http.MethodGet {
			m.mu.Lock()
			m.subscribed = true
			events := m.events
			eof := m.eof
			m.mu.Unlock()
			w.Header().Set("Content-Type", "text/event-stream")
			w.WriteHeader(http.StatusOK)
			w.(http.Flusher).Flush()
			for _, event := range events {
				_, _ = fmt.Fprintf(w, "data: %s\n\n", event)
				w.(http.Flusher).Flush()
			}
			if !eof {
				<-r.Context().Done()
			}
			return
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode synthetic request: %v", err)
		}
		m.mu.Lock()
		if !m.subscribed {
			t.Error("input submitted before subscription")
		}
		m.posts = append(m.posts, agentHelperPost{body: body, header: r.Header.Clone()})
		status, message := m.postStatus, m.postError
		if m.failures > 0 && (len(m.posts) > 1 || m.failInput) {
			m.failures--
		} else {
			status = 0
		}
		m.mu.Unlock()
		if status != 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(status)
			_, _ = io.WriteString(w, message)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	t.Cleanup(m.server.Close)
	hc := m.server.Client()
	hc.Transport = agentHelperTransport{base: hc.Transport, mock: m}
	c := openai.NewClient(option.WithBaseURL(m.server.URL), option.WithAPIKey("synthetic"), option.WithHTTPClient(hc), option.WithMaxRetries(0))
	return m, c
}
func (m *agentHelperServer) submissions() []agentHelperPost {
	m.mu.Lock()
	defer m.mu.Unlock()
	return append([]agentHelperPost(nil), m.posts...)
}
func agentEvent(kind, id, extra string) string {
	return fmt.Sprintf(`{"type":%q,"event_id":%q%s}`, "agent.session."+kind, id, extra)
}
func agentCreated(id, turn, subagent string) string {
	return agentEvent("turn.created", id, fmt.Sprintf(`,"turn_id":%q,"turn":{"id":%q,"subagent_id":%s}`, turn, turn, subagent))
}
func agentCall(id, turn, call, name, args string) string {
	return agentEvent("turn.item.added", id, fmt.Sprintf(`,"turn_id":%q,"item":{"id":%q,"type":"function_call","turn_id":%q,"call_id":%q,"name":%q,"arguments":%s,"status":"in_progress"}`, turn, id, turn, call, name, args))
}
func agentEnd(turn string) []string {
	return []string{agentEvent("turn.completed", "done-"+turn, `,"turn_id":"`+turn+`"`), agentEvent("idle", "idle-"+turn, `,"session":{"status":"idle"}`)}
}
func consumeAgentStream(t *testing.T, stream *openai.AgentSessionStream) []openai.AgentSessionEventUnion {
	t.Helper()
	defer func() { _ = stream.Close() }()
	var result []openai.AgentSessionEventUnion
	for stream.Next() {
		result = append(result, stream.Current())
	}
	if err := stream.Err(); err != nil {
		t.Fatal(err)
	}
	return result
}

func TestAgentsHelperLifecycle(t *testing.T) {
	for _, terminal := range []string{"completed", "failed", "cancelled"} {
		t.Run(terminal, func(t *testing.T) {
			events := []string{agentEvent("idle", "initial", ""), agentCreated("sub", "subturn", `"worker"`), agentEvent("turn.completed", "subdone", `,"turn_id":"subturn"`), agentEvent("idle", "subidle", ""), agentCreated("create", "target", "null"), agentCreated("later", "other", "null"), agentEvent("turn.completed", "otherdone", `,"turn_id":"other"`), agentEvent("idle", "otheridle", ""), agentEvent("turn.output_text.delta", "text", `,"turn_id":"target","delta":"hello"`), agentEvent("turn."+terminal, "done", `,"turn_id":"target"`), agentEvent("idle", "final", ""), agentEvent("idle", "unread", "")}
			m, c := newAgentHelperServer(t, events...)
			stream := c.Beta.Agents.Sessions.Stream(context.Background(), "session", openai.AgentSessionStreamParams{Input: "hello"})
			got := consumeAgentStream(t, stream)
			if len(got) != 11 || got[8].AsAgentSessionTurnOutputTextDelta().Delta != "hello" || got[10].EventID != "final" {
				t.Fatalf("unexpected events: %d", len(got))
			}
			if m.closed.Load() != 1 {
				t.Fatalf("body closed %d times", m.closed.Load())
			}
			posts := m.submissions()
			if len(posts) != 1 {
				t.Fatalf("posts: %d", len(posts))
			}
			want := map[string]any{"events": []any{map[string]any{"type": "agent.session.input.message", "input": []any{map[string]any{"role": "user", "content": []any{map[string]any{"type": "input_text", "text": "hello"}}}}}}}
			if !reflect.DeepEqual(posts[0].body, want) {
				t.Fatalf("input: %#v", posts[0].body)
			}
		})
	}
}

func TestAgentsHelperFailuresAndClose(t *testing.T) {
	t.Run("session failure visible", func(t *testing.T) {
		m, c := newAgentHelperServer(t, agentEvent("failed", "failed", `,"session":{"status":"failed"}`))
		got := consumeAgentStream(t, c.Beta.Agents.Sessions.Stream(context.Background(), "session", openai.AgentSessionStreamParams{Input: "hi"}))
		if len(got) != 1 || got[0].AsAgentSessionFailed().Session.Status != "failed" || m.closed.Load() != 1 {
			t.Fatal("missing typed failure or cleanup")
		}
	})
	t.Run("unexpected EOF", func(t *testing.T) {
		m, c := newAgentHelperServer(t, agentEvent("idle", "initial", ""))
		m.eof = true
		s := c.Beta.Agents.Sessions.Stream(context.Background(), "session", openai.AgentSessionStreamParams{Input: "hi"})
		defer func() { _ = s.Close() }()
		for s.Next() {
		}
		if !errors.Is(s.Err(), io.ErrUnexpectedEOF) || m.closed.Load() != 1 {
			t.Fatalf("err=%v closed=%d", s.Err(), m.closed.Load())
		}
	})
	for _, input := range []any{"", []openai.AgentSessionInputMessageParam{}, 42} {
		t.Run(fmt.Sprintf("invalid %T", input), func(t *testing.T) {
			m, c := newAgentHelperServer(t)
			s := c.Beta.Agents.Sessions.Stream(context.Background(), "session", openai.AgentSessionStreamParams{Input: input})
			defer func() { _ = s.Close() }()
			if s.Err() == nil || s.Next() || m.reads != 0 {
				t.Fatal("invalid input sent")
			}
		})
	}
	t.Run("active session rejected", func(t *testing.T) {
		m, c := newAgentHelperServer(t)
		m.status = "in_progress"
		s := c.Beta.Agents.Sessions.Stream(context.Background(), "session", openai.AgentSessionStreamParams{Input: "hi"})
		defer func() { _ = s.Close() }()
		if s.Err() == nil || len(m.submissions()) != 0 || m.subscribed {
			t.Fatal("active session accepted")
		}
	})
	for _, cancelContext := range []bool{false, true} {
		t.Run(fmt.Sprintf("cancelContext=%t", cancelContext), func(t *testing.T) {
			m, c := newAgentHelperServer(t)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			s := c.Beta.Agents.Sessions.Stream(ctx, "session", openai.AgentSessionStreamParams{Input: "hi"})
			defer func() { _ = s.Close() }()
			done := make(chan bool, 1)
			go func() { done <- s.Next() }()
			if cancelContext {
				cancel()
			} else {
				_ = s.Close()
			}
			select {
			case ok := <-done:
				if ok {
					t.Error("unexpected event")
				}
			case <-time.After(time.Second):
				t.Fatal("Next blocked after cancellation")
			}
			if cancelContext && !errors.Is(s.Err(), context.Canceled) {
				t.Fatalf("err=%v", s.Err())
			}
			if m.closed.Load() != 1 || len(m.submissions()) != 1 {
				t.Fatal("cleanup or backend cancellation error")
			}
		})
	}
}

func TestAgentsHelperTools(t *testing.T) {
	for _, tc := range []struct {
		name, args string
		output     any
		handlerErr error
		success    bool
		want       any
	}{
		{name: "object", args: `{"nested":{"value":1}}`, output: map[string]any{"answer": 42}, success: true, want: `{"answer":42}`},
		{name: "encoded object", args: `"{\"nested\":{\"value\":1}}"`, output: "answer", success: true, want: "answer"},
		{name: "nil", args: `{}`, output: nil, success: true, want: nil},
		{name: "empty text", args: `{}`, output: "", success: true, want: ""},
		{name: "content", args: `{}`, output: []openai.InputContentParamUnion{openai.InputContentParamOfParamInputText("answer")}, success: true, want: []any{map[string]any{"type": "input_text", "text": "answer"}}},
		{name: "empty content", args: `{}`, output: []openai.InputContentParamUnion{}, success: true, want: []any{}},
		{name: "handler error", args: `{}`, handlerErr: errors.New("secret-token")},
		{name: "invalid JSON", args: `"{"`},
		{name: "array args", args: `[]`},
		{name: "null args", args: `null`},
		{name: "scalar args", args: `42`},
		{name: "unserializable", args: `{}`, output: map[string]any{"secret-token": make(chan int)}},
		{name: "invalid output", args: `{}`, output: 42},
	} {
		t.Run(tc.name, func(t *testing.T) {
			events := append([]string{agentCreated("created", "turn", "null"), agentCall("call", "turn", "call", "tool", tc.args)}, agentEnd("turn")...)
			m, c := newAgentHelperServer(t, events...)
			calls := 0
			s := c.Beta.Agents.Sessions.Stream(context.Background(), "session", openai.AgentSessionStreamParams{Input: "hi", ToolHandlers: map[string]openai.AgentToolHandler{"tool": func(_ context.Context, args map[string]any) (any, error) {
				calls++
				if nested, ok := args["nested"].(map[string]any); ok {
					nested["value"] = "mutated"
				}
				return tc.output, tc.handlerErr
			}}})
			got := consumeAgentStream(t, s)
			if len(got) != 4 {
				t.Fatalf("events=%d", len(got))
			}
			if tc.name == "object" && got[1].Item.Arguments.(map[string]any)["nested"].(map[string]any)["value"] != float64(1) {
				t.Fatal("event mutated")
			}
			if strings.Contains(got[1].RawJSON(), "mutated") {
				t.Fatal("raw event mutated")
			}
			posts := m.submissions()
			if len(posts) != 2 {
				t.Fatalf("posts=%d", len(posts))
			}
			result := posts[1].body["events"].([]any)[0].(map[string]any)
			if result["type"] != "agent.session.input.tool_result" || result["turn_id"] != "turn" || result["call_id"] != "call" || result["success"] != tc.success {
				t.Fatalf("result=%#v", result)
			}
			if tc.success {
				if output, exists := result["output"]; !exists || !reflect.DeepEqual(output, tc.want) {
					t.Fatalf("output=%#v want=%#v", result, tc.want)
				}
			} else if result["error"] != "Tool handler failed." {
				t.Fatalf("failure=%#v", result)
			}
			if tc.success && calls != 1 {
				t.Fatalf("calls=%d", calls)
			}
			encoded, err := json.Marshal(posts[1].body)
			if err != nil {
				t.Fatal(err)
			}
			if strings.Contains(string(encoded), "secret-token") {
				t.Fatal("secret leaked")
			}
		})
	}
}

func TestAgentsHelperToolDedupAndOptions(t *testing.T) {
	events := []string{agentCreated("created", "turn", "null"), agentCall("call", "turn", "same", "tool", `{}`), agentCall("call", "turn", "same", "tool", `{}`)}
	for i := 0; i < 1030; i++ {
		events = append(events, agentEvent("turn.output_text.delta", fmt.Sprintf("delta%d", i), `,"delta":"x"`))
	}
	events = append(events, agentCall("call", "turn", "same", "tool", `{}`), agentCall("other-event", "turn", "same", "tool", `{}`), agentCall("worker-call", "worker", "same", "tool", `{}`), agentCall("manual", "turn", "manual", "unknown", `{}`))
	events = append(events, agentEnd("turn")...)
	m, c := newAgentHelperServer(t, events...)
	calls := 0
	s := c.Beta.Agents.Sessions.Stream(context.Background(), "session", openai.AgentSessionStreamParams{Input: []openai.AgentSessionInputMessageParam{{Content: []openai.InputContentParamUnion{openai.InputContentParamOfParamInputText("hi")}}}, IdempotencyKey: openai.String("named"), ToolHandlers: map[string]openai.AgentToolHandler{"tool": func(context.Context, map[string]any) (any, error) { calls++; return nil, nil }}}, option.WithHeader("iDeMpOtEnCy-KeY", "override"), option.WithHeader("X-Helper-Test", "preserved"), option.WithQuery("test", "value"))
	if !s.Next() {
		t.Fatal("missing turn event")
	}
	if !s.Next() || calls != 0 {
		t.Fatal("handler ran before call was yielded")
	}
	consumeAgentStream(t, s)
	posts := m.submissions()
	if len(posts) != 3 || calls != 2 {
		t.Fatalf("posts=%d calls=%d", len(posts), calls)
	}
	keys := map[string]bool{}
	for i, post := range posts {
		key := post.header.Get("Idempotency-Key")
		if key == "" || keys[key] || len(post.header.Values("Idempotency-Key")) != 1 || post.header.Get("X-Helper-Test") != "preserved" {
			t.Fatal("invalid key or options")
		}
		keys[key] = true
		if i == 0 && key != "override" {
			t.Fatalf("input key=%s", key)
		}
	}
}

func TestAgentsHelperIdempotencyAndRegistrationRace(t *testing.T) {
	for _, exact := range []bool{true, false} {
		t.Run(fmt.Sprintf("exact=%t", exact), func(t *testing.T) {
			m, c := newAgentHelperServer(t, append([]string{agentCreated("created", "turn", "null"), agentCall("call", "turn", "call", "tool", `{}`)}, agentEnd("turn")...)...)
			m.postStatus = 400
			m.failures = 1
			message := "Unknown pending tool call: call"
			if !exact {
				message += " different"
			}
			m.postError = fmt.Sprintf(`{"error":{"type":"invalid_request_error","code":"invalid_request_error","message":%q,"param":null}}`, message)
			s := c.Beta.Agents.Sessions.Stream(context.Background(), "session", openai.AgentSessionStreamParams{Input: "hi", ToolHandlers: map[string]openai.AgentToolHandler{"tool": func(context.Context, map[string]any) (any, error) { return "ok", nil }}})
			defer func() { _ = s.Close() }()
			for s.Next() {
			}
			posts := m.submissions()
			if exact {
				if s.Err() != nil || len(posts) != 3 || posts[1].header.Get("Idempotency-Key") != posts[2].header.Get("Idempotency-Key") {
					t.Fatalf("race retry failed: posts=%d err=%v", len(posts), s.Err())
				}
			} else if s.Err() == nil || len(posts) != 2 {
				t.Fatal("unrelated error retried")
			}
			if m.closed.Load() != 1 {
				t.Fatal("body not closed")
			}
		})
	}
	m, c := newAgentHelperServer(t, append([]string{agentCreated("created", "turn", "null")}, agentEnd("turn")...)...)
	m.drop.Store(true)
	consumeAgentStream(t, c.Beta.Agents.Sessions.Stream(context.Background(), "session", openai.AgentSessionStreamParams{Input: "hi"}, option.WithMaxRetries(1)))
	posts := m.submissions()
	if len(posts) != 2 || posts[0].header.Get("Idempotency-Key") == "" || posts[0].header.Get("Idempotency-Key") != posts[1].header.Get("Idempotency-Key") || !reflect.DeepEqual(posts[0].body, posts[1].body) {
		t.Fatal("accepted input not retried idempotently")
	}
}

func TestAgentsHelperMessageOutputText(t *testing.T) {
	for _, phase := range []string{"commentary", "final_answer"} {
		var message openai.AgentSessionMessage
		raw := fmt.Sprintf(`{"phase":%q,"content":[{"type":"input_text","text":"ignored"},{"type":"output_text","text":"one"},{"type":"input_image","image_url":"synthetic"},{"type":"output_text","text":" two"}]}`, phase)
		if err := json.Unmarshal([]byte(raw), &message); err != nil {
			t.Fatal(err)
		}
		if message.OutputText() != "one two" || message.RawJSON() != raw {
			t.Fatal("unexpected text or mutation")
		}
	}
	if (openai.AgentSessionMessage{}).OutputText() != "" {
		t.Fatal("empty message")
	}
}

func TestAgentsHelperInputFailureClosesSubscription(t *testing.T) {
	m, c := newAgentHelperServer(t)
	m.failInput = true
	m.failures = 1
	m.postStatus = 400
	m.postError = `{"error":{"code":"invalid_request_error","message":"synthetic input failure"}}`
	s := c.Beta.Agents.Sessions.Stream(context.Background(), "session", openai.AgentSessionStreamParams{Input: "hi"})
	defer func() { _ = s.Close() }()
	if s.Err() == nil || s.Next() || m.closed.Load() != 1 || len(m.submissions()) != 1 {
		t.Fatal("input failure did not clean up")
	}
}

func TestAgentsHelperToolResponseLost(t *testing.T) {
	m, c := newAgentHelperServer(t, append([]string{agentCreated("created", "turn", "null"), agentCall("call", "turn", "call", "tool", `{}`)}, agentEnd("turn")...)...)
	calls := 0
	consumeAgentStream(t, c.Beta.Agents.Sessions.Stream(context.Background(), "session", openai.AgentSessionStreamParams{Input: "hi", ToolHandlers: map[string]openai.AgentToolHandler{"tool": func(context.Context, map[string]any) (any, error) { calls++; m.drop.Store(true); return "answer", nil }}}, option.WithMaxRetries(1)))
	posts := m.submissions()
	if calls != 1 || len(posts) != 3 || posts[0].header.Get("Idempotency-Key") == posts[1].header.Get("Idempotency-Key") || posts[1].header.Get("Idempotency-Key") != posts[2].header.Get("Idempotency-Key") || !reflect.DeepEqual(posts[1].body, posts[2].body) {
		t.Fatal("tool result transport retry was not idempotent")
	}
}

func TestAgentsHelperCancelledHandlerAndPanicCleanup(t *testing.T) {
	for _, panics := range []bool{false, true} {
		t.Run(fmt.Sprintf("panic=%t", panics), func(t *testing.T) {
			m, c := newAgentHelperServer(t, agentCreated("created", "turn", "null"), agentCall("call", "turn", "call", "tool", `{}`))
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			s := c.Beta.Agents.Sessions.Stream(ctx, "session", openai.AgentSessionStreamParams{Input: "hi", ToolHandlers: map[string]openai.AgentToolHandler{"tool": func(handlerCtx context.Context, _ map[string]any) (any, error) {
				if panics {
					panic("synthetic")
				}
				cancel()
				<-handlerCtx.Done()
				return nil, handlerCtx.Err()
			}}})
			defer func() { _ = s.Close() }()
			if !s.Next() {
				t.Fatal("missing turn event")
			}
			if !s.Next() {
				t.Fatal("missing call")
			}
			var recovered any
			func() {
				defer func() { recovered = recover() }()
				if s.Next() {
					t.Error("unexpected event")
				}
			}()
			if panics && recovered == nil {
				t.Fatal("panic not propagated")
			}
			if !panics && !errors.Is(s.Err(), context.Canceled) {
				t.Fatalf("err=%v", s.Err())
			}
			if m.closed.Load() != 1 || len(m.submissions()) != 1 {
				t.Fatal("handler failure did not close without a result")
			}
		})
	}
}

func TestAgentsHelperListedMessageOutputText(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/agents/sessions/session/items" {
			t.Errorf("path=%s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"data":[{"id":"message","type":"message","role":"assistant","phase":"commentary","content":[{"type":"output_text","text":"hello"},{"type":"input_image","image_url":"synthetic"},{"type":"output_text","text":" world"}]}],"has_more":false}`)
	}))
	defer server.Close()
	client := openai.NewClient(option.WithAPIKey("synthetic"), option.WithBaseURL(server.URL), option.WithMaxRetries(0))
	page, err := client.Beta.Agents.Sessions.Items.List(context.Background(), "session", openai.BetaAgentSessionItemListParams{})
	if err != nil {
		t.Fatal(err)
	}
	if len(page.Data) != 1 {
		t.Fatalf("items=%d", len(page.Data))
	}
	item := page.Data[0]
	if item.AsMessage().OutputText() != "hello world" || item.AsMessage().Phase != "commentary" {
		t.Fatal("listed message text or phase was lost")
	}
}

func TestAgentsHelperConsumerMutationCannotChangeDispatch(t *testing.T) {
	m, c := newAgentHelperServer(t, append([]string{agentCreated("created", "turn", "null"), agentCall("call", "turn", "call", "tool", `{"nested":{"value":"original"}}`)}, agentEnd("turn")...)...)
	calls := 0
	handlers := map[string]openai.AgentToolHandler{"tool": func(_ context.Context, args map[string]any) (any, error) {
		calls++
		if args["nested"].(map[string]any)["value"] != "original" {
			t.Fatal("consumer changed handler arguments")
		}
		return "answer", nil
	}, "wrong": func(context.Context, map[string]any) (any, error) {
		t.Fatal("consumer redirected handler")
		return nil, nil
	}}
	s := c.Beta.Agents.Sessions.Stream(context.Background(), "session", openai.AgentSessionStreamParams{Input: "hi", ToolHandlers: handlers})
	defer func() { _ = s.Close() }()
	if !s.Next() {
		t.Fatal("missing turn")
	}
	if !s.Next() {
		t.Fatal("missing call")
	}
	event := s.Current()
	event.Item.Arguments.(map[string]any)["nested"].(map[string]any)["value"] = "mutated"
	event.Item.Name = "wrong"
	event.Item.TurnID = "wrong-turn"
	event.Item.CallID = "wrong-call"
	handlers["tool"] = handlers["wrong"]
	if calls != 0 {
		t.Fatal("handler called before yield")
	}
	consumeAgentStream(t, s)
	posts := m.submissions()
	if calls != 1 || len(posts) != 2 {
		t.Fatalf("calls=%d posts=%d", calls, len(posts))
	}
	result := posts[1].body["events"].([]any)[0].(map[string]any)
	if result["turn_id"] != "turn" || result["call_id"] != "call" || result["output"] != "answer" {
		t.Fatalf("consumer changed dispatch: %#v", result)
	}
}

func TestAgentsHelperRejectsResponseBodyOverrides(t *testing.T) {
	for _, level := range []string{"request", "inherited", "events"} {
		t.Run(level, func(t *testing.T) {
			for _, target := range []string{"map", "typed session", "raw response"} {
				t.Run(target, func(t *testing.T) {
					m, c := newAgentHelperServer(t)
					var body map[string]any
					var session *openai.AgentSession
					var raw *http.Response
					var dst any = &body
					if target == "typed session" {
						dst = &session
					}
					if target == "raw response" {
						dst = &raw
					}
					override := option.WithResponseBodyInto(dst)
					var opts []option.RequestOption
					switch level {
					case "request":
						opts = []option.RequestOption{override}
					case "inherited":
						c.Beta.Agents.Sessions.Options = append(c.Beta.Agents.Sessions.Options, override)
					case "events":
						c.Beta.Agents.Sessions.Events.Options = append(c.Beta.Agents.Sessions.Events.Options, override)
					}
					s := c.Beta.Agents.Sessions.Stream(context.Background(), "session", openai.AgentSessionStreamParams{Input: "hi"}, opts...)
					defer func() { _ = s.Close() }()
					if s.Err() == nil || !strings.Contains(s.Err().Error(), "WithResponseBodyInto") || s.Next() {
						t.Fatalf("override error=%v", s.Err())
					}
					if m.subscribed || len(m.submissions()) != 0 || raw != nil || session != nil || body != nil {
						t.Fatal("response override reached transport")
					}
					if level != "events" && m.reads != 0 {
						t.Fatal("idle probe sent before rejecting override")
					}
				})
			}
		})
	}
}

func TestAgentsHelperProtocolErrorSurfacesAndCloses(t *testing.T) {
	m, c := newAgentHelperServer(t, agentCreated("created", "turn", "null"), `{"type":"error","event_id":"protocol-error","session_id":"session","error":{"code":"internal_error","message":"synthetic protocol failure"}}`, agentEvent("failed", "failed", `,"session":{"status":"failed"}`))
	s := c.Beta.Agents.Sessions.Stream(context.Background(), "session", openai.AgentSessionStreamParams{Input: "hi"})
	defer func() { _ = s.Close() }()
	if !s.Next() || s.Current().Type != "agent.session.turn.created" {
		t.Fatal("missing typed event before protocol error")
	}
	if s.Next() {
		t.Fatal("protocol error should follow the SDK Err path")
	}
	var streamError *ssestream.StreamError
	if !errors.As(s.Err(), &streamError) || !strings.Contains(string(streamError.Event.Data), "protocol-error") {
		t.Fatalf("missing original protocol error: %v", s.Err())
	}
	if m.closed.Load() != 1 || len(m.submissions()) != 1 {
		t.Fatal("protocol error did not close without cancelling the backend")
	}
}

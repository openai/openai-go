package live_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	wire "github.com/coder/websocket"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/live"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

type responseItemGetter interface {
	GetItem() *responses.ResponseInputItemUnionParam
}

var (
	_ responseItemGetter = live.ClientEventUnionParam{}
	_ responseItemGetter = live.ForkClientEventUnionParam{}
	_ responseItemGetter = live.ConnectClientEventUnionParam{}
)

func TestDecisionsPreserveResponseItemGetters(t *testing.T) {
	item := &live.ResponseItemCreateEventParam{}
	for _, event := range []responseItemGetter{
		live.ClientEventUnionParam{OfResponseItemCreate: item},
		live.ForkClientEventUnionParam{OfResponseItemCreate: item},
		live.ConnectClientEventUnionParam{OfResponseItemCreate: item},
	} {
		if event.GetItem() != &item.Item {
			t.Errorf("%T changed its existing Responses item getter", event)
		}
	}
}

func TestDecisionsConcreteItemsRetainWireBody(t *testing.T) {
	item := &live.DecisionItemCreateEventParam{
		Item: live.DecisionItemCreateEventItemParam{
			Content: live.DecisionItemCreateEventItemContentUnionParam{OfString: openai.String("synthetic input")},
		},
	}
	for _, event := range []responseItemGetter{
		live.ClientEventUnionParam{OfDecisionItemCreate: item},
		live.ForkClientEventUnionParam{OfDecisionItemCreate: item},
		live.ConnectClientEventUnionParam{OfDecisionItemCreate: item},
	} {
		if event.GetItem() != nil {
			t.Errorf("%T exposed a Decision item as a Responses item", event)
		}
		body, err := json.Marshal(event)
		if err != nil {
			t.Fatal(err)
		}
		var got struct {
			Type string `json:"type"`
			Item struct {
				Content string `json:"content"`
				Role    string `json:"role"`
			} `json:"item"`
		}
		if err := json.Unmarshal(body, &got); err != nil {
			t.Fatal(err)
		}
		if got.Type != "decision.item.create" || got.Item.Content != "synthetic input" || got.Item.Role != "user" {
			t.Errorf("%T lost its concrete Decision item: %s", event, body)
		}
	}
}

func TestDecisionsSessionHTTPBodies(t *testing.T) {
	const sessionJSON = `{"model":"synthetic-live","delegation":{"type":"decisions","decisions":{"model":"synthetic-decisions","questions":[{"type":"choice","instructions":"Ready?","choices":[{"value":false},{"value":true}]}]}},"triggers":[{"type":"decisions","event_type":"session.input_transcript.delta","decisions":{"questions":[{"type":"predicate","instructions":"Input ready?"}]}},{"type":"decisions","event_type":"session.output_transcript.delta","decisions":{"questions":[{"type":"predicate","instructions":"Output ready?"}]}}]}`
	for _, endpoint := range []string{"/live/sessions", "/live/sessions/session_test/accept"} {
		t.Run(endpoint, func(t *testing.T) {
			type request struct {
				method, path string
				body         []byte
				err          error
			}
			requests := make(chan request, 1)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				requests <- request{r.Method, r.URL.Path, body, err}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"id":"session_test"}`))
			}))
			defer server.Close()
			client := openai.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("synthetic-test-key"), option.WithUnsafeAllowHTTP(), option.WithMaxRetries(0))
			if endpoint == "/live/sessions" {
				var session live.MediaSessionConfigParam
				if err := json.Unmarshal([]byte(sessionJSON), &session); err != nil {
					t.Fatal(err)
				}
				_, err := client.Live.New(t.Context(), live.LiveNewParams{Session: session, Transport: live.LiveNewParamsTransport{Sdp: "synthetic-sdp"}})
				if err != nil {
					t.Fatal(err)
				}
			} else {
				var session live.SessionAcceptParamsSession
				if err := json.Unmarshal([]byte(sessionJSON), &session); err != nil {
					t.Fatal(err)
				}
				if err := client.Live.Sessions.Accept(t.Context(), "session_test", live.SessionAcceptParams{Session: session}); err != nil {
					t.Fatal(err)
				}
			}
			got := <-requests
			if got.err != nil || got.method != http.MethodPost || got.path != endpoint {
				t.Fatalf("unexpected request: method=%s path=%s err=%v", got.method, got.path, got.err)
			}
			var body struct{ Session map[string]any }
			var want map[string]any
			if err := json.Unmarshal(got.body, &body); err != nil {
				t.Fatal(err)
			}
			if err := json.Unmarshal([]byte(sessionJSON), &want); err != nil {
				t.Fatal(err)
			}
			if endpoint == "/live/sessions/session_test/accept" {
				want["type"] = "live"
			}
			if !reflect.DeepEqual(body.Session, want) {
				t.Fatalf("session body changed: got %s; want %s", got.body, sessionJSON)
			}
		})
	}
}

func TestDecisionsTriggerUpdateOmission(t *testing.T) {
	for _, tc := range []struct {
		name     string
		triggers []live.DecisionTriggerParam
		want     string
	}{
		{name: "omitted", want: ""},
		{name: "clear", triggers: []live.DecisionTriggerParam{}, want: "[]"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			body, err := json.Marshal(live.SessionUpdateConfigParam{Triggers: tc.triggers})
			if err != nil {
				t.Fatal(err)
			}
			var fields map[string]json.RawMessage
			if err := json.Unmarshal(body, &fields); err != nil {
				t.Fatal(err)
			}
			if string(fields["triggers"]) != tc.want {
				t.Fatalf("unexpected trigger update: %s", body)
			}
		})
	}
}

func TestDecisionsWebSocketResultFields(t *testing.T) {
	for _, nulls := range []bool{false, true} {
		name := "values"
		if nulls {
			name = "nulls"
		}
		t.Run(name, func(t *testing.T) {
			fields := map[string]any{
				"type": "decision.event", "event_id": "event_synthetic",
				"delegation_id": "delegation_synthetic", "trigger_id": "trigger_synthetic",
				"source_event_id": "source_synthetic", "source_event_type": "session.input_transcript.delta",
				"event": map[string]any{"type": "decision.completed", "decision": map[string]any{
					"model": "synthetic-decisions", "answers": []any{map[string]any{"type": "predicate", "name": "ready", "probability": 0.75}},
					"usage": map[string]any{"input_tokens": 8, "output_tokens": 2, "total_tokens": 10,
						"input_tokens_details":  map[string]any{"cached_tokens": 0, "cache_write_tokens": 0},
						"output_tokens_details": map[string]any{"reasoning_tokens": 0}},
				}},
			}
			if nulls {
				for _, key := range []string{"delegation_id", "trigger_id", "source_event_id", "source_event_type"} {
					fields[key] = nil
				}
			}
			frame, err := json.Marshal(fields)
			if err != nil {
				t.Fatal(err)
			}
			server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				socket, acceptErr := wire.Accept(w, r, nil)
				if acceptErr != nil {
					t.Error(acceptErr)
					return
				}
				defer func() { _ = socket.CloseNow() }()
				if writeErr := socket.Write(r.Context(), wire.MessageText, frame); writeErr != nil {
					t.Error(writeErr)
					return
				}
				_, _, _ = socket.Read(r.Context())
			}))
			defer server.Close()
			client := openai.NewClient(option.WithAPIKey("synthetic"), option.WithBaseURL(server.URL), option.WithHTTPClient(server.Client()))
			ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
			defer cancel()
			connection, err := client.Live.Sideband.Connect(ctx, "live_synthetic", live.SidebandConnectionOptions{})
			if err != nil {
				t.Fatal(err)
			}
			defer connection.Abort()
			event, err := connection.Recv(ctx)
			if err != nil {
				t.Fatal(err)
			}
			result := event.AsDecisionEvent()
			if result.EventID != "event_synthetic" || result.Event.Decision.Model != "synthetic-decisions" || result.Event.Decision.Usage.TotalTokens != 10 {
				t.Fatalf("Decision result lost fields: %+v", result)
			}
			got := []string{result.DelegationID, result.TriggerID, result.SourceEventID, string(result.SourceEventType)}
			want := []string{"delegation_synthetic", "trigger_synthetic", "source_synthetic", "session.input_transcript.delta"}
			if nulls {
				want = []string{"", "", "", ""}
			}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("Decision correlation fields: got %q; want %q", got, want)
			}
			for _, isNull := range []bool{result.JSON.DelegationID.Raw() == "null", result.JSON.TriggerID.Raw() == "null", result.JSON.SourceEventID.Raw() == "null", result.JSON.SourceEventType.Raw() == "null"} {
				if isNull != nulls {
					t.Fatal("Decision correlation null metadata changed")
				}
			}
		})
	}
}

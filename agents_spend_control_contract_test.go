package openai_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/param"
)

func TestSessionSpendControlUpdateBody(t *testing.T) {
	for _, tt := range []struct {
		name    string
		control openai.BetaAgentSessionUpdateParamsSpendControl
		want    map[string]any
	}{
		{"omitted", openai.BetaAgentSessionUpdateParamsSpendControl{}, map[string]any{}},
		{"null object", param.NullStruct[openai.BetaAgentSessionUpdateParamsSpendControl](), map[string]any{"spend_control": nil}},
		{"null limit", openai.BetaAgentSessionUpdateParamsSpendControl{Limit: param.Null[int64]()}, map[string]any{"spend_control": map[string]any{"limit": nil}}},
		{"positive limit", openai.BetaAgentSessionUpdateParamsSpendControl{Limit: openai.Int(500)}, map[string]any{"spend_control": map[string]any{"limit": float64(500)}}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost || r.URL.Path != "/agents/sessions/session_test" {
					t.Errorf("request = %s %s, want POST /agents/sessions/session_test", r.Method, r.URL.Path)
				}
				var body map[string]any
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Errorf("decode request: %v", err)
					http.Error(w, "invalid JSON", http.StatusBadRequest)
					return
				}
				if !reflect.DeepEqual(body, tt.want) {
					t.Errorf("body = %#v, want %#v", body, tt.want)
				}
				w.Header().Set("Content-Type", "application/json")
				if _, err := w.Write([]byte(`{"id":"session_test","object":"agent.session"}`)); err != nil {
					t.Errorf("write response: %v", err)
				}
			}))
			t.Cleanup(server.Close)
			client := openai.NewClient(option.WithAPIKey("test-key"), option.WithBaseURL(server.URL), option.WithUnsafeAllowHTTP(), option.WithHTTPClient(server.Client()), option.WithMaxRetries(0))
			session, err := client.Beta.Agents.Sessions.Update(t.Context(), "session_test", openai.BetaAgentSessionUpdateParams{SpendControl: tt.control})
			if err != nil {
				t.Fatal(err)
			}
			if session.ID != "session_test" {
				t.Errorf("session ID = %q, want session_test", session.ID)
			}
		})
	}
}

func TestSessionSpendControlConsumedPresence(t *testing.T) {
	for _, tt := range []struct {
		name  string
		raw   string
		want  int64
		valid bool
	}{
		{"unavailable", "null", 0, false},
		{"zero spend", "0", 0, true},
		{"recorded spend", "42", 42, true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet || r.URL.Path != "/agents/sessions/session_test" {
					t.Errorf("request = %s %s, want GET /agents/sessions/session_test", r.Method, r.URL.Path)
				}
				w.Header().Set("Content-Type", "application/json")
				if _, err := w.Write([]byte(`{"id":"session_test","object":"agent.session","spend_control":{"limit":500,"consumed":` + tt.raw + `}}`)); err != nil {
					t.Errorf("write response: %v", err)
				}
			}))
			t.Cleanup(server.Close)
			client := openai.NewClient(option.WithAPIKey("test-key"), option.WithBaseURL(server.URL), option.WithUnsafeAllowHTTP(), option.WithHTTPClient(server.Client()), option.WithMaxRetries(0))
			session, err := client.Beta.Agents.Sessions.Get(t.Context(), "session_test")
			if err != nil {
				t.Fatal(err)
			}
			control := session.SpendControl
			if !session.JSON.SpendControl.Valid() || control.Limit != 500 || !control.JSON.Limit.Valid() {
				t.Errorf("spend control = %#v, want present limit 500", control)
			}
			if control.Consumed != tt.want || control.JSON.Consumed.Valid() != tt.valid || control.JSON.Consumed.Raw() != tt.raw {
				t.Errorf("consumed = %d, valid = %t, raw = %q; want %d, %t, %q", control.Consumed, control.JSON.Consumed.Valid(), control.JSON.Consumed.Raw(), tt.want, tt.valid, tt.raw)
			}
		})
	}
}

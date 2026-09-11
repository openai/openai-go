package live_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/live"
	"github.com/openai/openai-go/v3/option"
)

func TestLiveUnifiedSDPContract(t *testing.T) {
	const offer = "v=0\r\ns=synthetic offer\r\n"
	const answer = "v=0\r\ns=synthetic answer\r\n"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/live/sessions" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		if r.URL.RawQuery != "" {
			t.Error("session configuration must be sent in the body")
		}
		for _, header := range []string{"OpenAI-Alpha", "OpenAI-Beta"} {
			if r.Header.Get(header) != "" {
				t.Errorf("GA request unexpectedly sets %s", header)
			}
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode request: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		want := map[string]any{
			"session":   map[string]any{"model": "gpt-live-1"},
			"transport": map[string]any{"type": "webrtc", "sdp": offer},
		}
		if !reflect.DeepEqual(body, want) {
			t.Errorf("body = %#v, want %#v", body, want)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]any{
			"session":   map[string]any{"id": "live_opaque-session-id"},
			"transport": map[string]any{"type": "webrtc", "sdp": answer},
		}); err != nil {
			t.Errorf("encode response: %v", err)
		}
	}))
	defer server.Close()
	client := openai.NewClient(option.WithAPIKey("test-key"), option.WithBaseURL(server.URL))
	result, err := client.Live.New(context.Background(), live.LiveNewParams{
		Session:   live.MediaSessionConfigParam{Model: "gpt-live-1"},
		Transport: live.LiveNewParamsTransport{Sdp: offer},
	})
	if err != nil {
		t.Fatalf("create session: %v", err)
	}
	if result.Session.ID != "live_opaque-session-id" || result.Transport.Sdp != answer {
		t.Fatalf("session ID or SDP answer was not preserved: %#v", result)
	}
}

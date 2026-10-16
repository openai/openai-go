package openai_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func TestSessionAttachesPrewarmedEnvironment(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/agents/sessions" {
			t.Errorf("request = %s %s, want POST /agents/sessions", r.Method, r.URL.Path)
		}
		if r.Header.Get("OpenAI-Beta") != "agents=v1" {
			t.Error("missing Managed Agents beta header")
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode request: %v", err)
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		want := map[string]any{"type": "openai_hosted", "environment_id": "ccarenv_test"}
		if !reflect.DeepEqual(body["environment"], want) {
			t.Errorf("environment = %#v, want %#v", body["environment"], want)
		}
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(`{"id":"session_test","object":"agent.session"}`)); err != nil {
			t.Errorf("write response: %v", err)
		}
	}))
	t.Cleanup(server.Close)
	client := openai.NewClient(option.WithAPIKey("test-key"), option.WithBaseURL(server.URL), option.WithUnsafeAllowHTTP(), option.WithHTTPClient(server.Client()), option.WithMaxRetries(0))
	environment := openai.EnvironmentParamUnion{OfParamOpenAIHosted: &openai.EnvironmentParamOpenAIHosted{EnvironmentID: openai.String("ccarenv_test")}}
	if got := environment.GetEnvironmentID(); got == nil || *got != "ccarenv_test" {
		t.Fatalf("GetEnvironmentID() = %v, want ccarenv_test", got)
	}
	session, err := client.Beta.Agents.Sessions.New(t.Context(), openai.BetaAgentSessionNewParams{Environment: environment})
	if err != nil {
		t.Fatal(err)
	}
	if session.ID != "session_test" {
		t.Errorf("session ID = %q, want session_test", session.ID)
	}
	omitted := openai.EnvironmentParamUnion{OfParamOpenAIHosted: &openai.EnvironmentParamOpenAIHosted{}}
	if got := omitted.GetEnvironmentID(); got != nil {
		t.Errorf("omitted GetEnvironmentID() = %q, want nil", *got)
	}
}

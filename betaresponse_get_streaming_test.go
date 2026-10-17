package openai_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func TestBetaResponseGetStreamingSendsStreamQueryParameter(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.Method, http.MethodGet; got != want {
			t.Errorf("request method = %q, want %q", got, want)
		}

		if got, want := r.URL.Path, "/responses/resp_123"; got != want {
			t.Errorf("request path = %q, want %q", got, want)
		}

		if got, want := r.URL.Query().Get("beta"), "true"; got != want {
			t.Errorf("beta query parameter = %q, want %q", got, want)
		}

		if got, want := r.URL.Query().Get("stream"), "true"; got != want {
			t.Errorf("stream query parameter = %q, want %q", got, want)
		}

		w.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
		if _, err := w.Write([]byte(": done\n\n")); err != nil {
			t.Errorf("write response body: %v", err)
		}
	}))
	defer server.Close()

	client := openai.NewClient(
		option.WithBaseURL(server.URL),
		option.WithAPIKey("My API Key"),
		option.WithMaxRetries(0),
	)

	stream := client.Beta.Responses.GetStreaming(
		context.TODO(),
		"resp_123",
		openai.BetaResponseGetParams{},
	)
	defer func() { _ = stream.Close() }()
}

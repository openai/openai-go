package openai_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/ssestream"
)

type eventBudgetDoer struct {
	body string
}

func (d eventBudgetDoer) Do(*http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": {"text/event-stream"}},
		Body:       io.NopCloser(strings.NewReader(d.body)),
	}, nil
}

func TestStreamingEventBudgetIsOptIn(t *testing.T) {
	payload := `{"id":"chatcmpl_test","object":"chat.completion.chunk","choices":[]}`
	body := "data: " + payload + "\n\n"
	eventBytes := len("data: " + payload + "\n")

	tests := map[string]struct {
		clientBudget       int
		methodBudget       int
		disableClientLimit bool
		wantErr            bool
	}{
		"default remains unlimited": {},
		"exact explicit budget": {
			methodBudget: eventBytes,
		},
		"explicit overflow": {
			methodBudget: eventBytes - 1,
			wantErr:      true,
		},
		"client-level overflow": {
			clientBudget: eventBytes - 1,
			wantErr:      true,
		},
		"method disables client budget": {
			clientBudget:       eventBytes - 1,
			disableClientLimit: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			clientOptions := []option.RequestOption{
				option.WithAPIKey("test-key"),
				option.WithBaseURL("https://example.com"),
				option.WithMaxRetries(0),
				option.WithHTTPClient(eventBudgetDoer{body: body}),
			}
			if test.clientBudget > 0 {
				clientOptions = append(clientOptions, option.WithSSEMaxEventBytes(test.clientBudget))
			}
			client := openai.NewClient(clientOptions...)

			var methodOptions []option.RequestOption
			switch {
			case test.methodBudget > 0:
				methodOptions = append(methodOptions, option.WithSSEMaxEventBytes(test.methodBudget))
			case test.disableClientLimit:
				methodOptions = append(methodOptions, option.WithSSEMaxEventBytes(0))
			}

			stream := client.Chat.Completions.NewStreaming(
				context.Background(),
				openai.ChatCompletionNewParams{Model: "gpt-4o-mini"},
				methodOptions...,
			)
			defer func() { _ = stream.Close() }()

			if test.wantErr {
				if stream.Next() {
					t.Fatal("oversized event unexpectedly produced a value")
				}
				if !errors.Is(stream.Err(), ssestream.ErrEventTooLarge) {
					t.Fatalf("stream error = %v, want ErrEventTooLarge", stream.Err())
				}
				return
			}

			if !stream.Next() {
				t.Fatalf("event was not delivered: %v", stream.Err())
			}
			if got := stream.Current().ID; got != "chatcmpl_test" {
				t.Fatalf("event ID = %q, want chatcmpl_test", got)
			}
		})
	}
}

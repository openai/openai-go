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
	"github.com/openai/openai-go/v3/shared"
)

type sseLimitTransport func(*http.Request) (*http.Response, error)

func (t sseLimitTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t(req)
}

func TestStreamingRejectsOversizedSSEEvent(t *testing.T) {
	client := openai.NewClient(
		option.WithAPIKey("My API Key"),
		option.WithHTTPClient(&http.Client{
			Transport: sseLimitTransport(func(*http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Status:     "200 OK",
					Header: http.Header{
						"Content-Type": {"text/event-stream"},
					},
					Body: io.NopCloser(strings.NewReader(
						strings.Repeat("data: {}\n", ssestream.DefaultMaxEventLines+1) + "\n",
					)),
				}, nil
			}),
		}),
	)
	stream := client.Chat.Completions.NewStreaming(context.Background(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{{
			OfUser: &openai.ChatCompletionUserMessageParam{
				Content: openai.ChatCompletionUserMessageParamContentUnion{
					OfString: openai.String("hello"),
				},
			},
		}},
		Model: shared.ChatModelGPT4o,
	})

	if stream.Next() {
		t.Fatalf("oversized event unexpectedly decoded: %+v", stream.Current())
	}
	if !errors.Is(stream.Err(), ssestream.ErrEventTooLarge) {
		t.Fatalf("stream error = %v, want %v", stream.Err(), ssestream.ErrEventTooLarge)
	}
}

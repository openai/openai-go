package testutil_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

// High memory use is intentional: these historically supported payloads must
// not regress under new body, event, or accumulation limits. Do not shrink the
// payloads or raise client limits to make a new cap pass. These are regression
// probes, not API maxima. Generate data in memory and run cases sequentially.
const largePayloadSize = 32*1024*1024 + 1

// The SSE decoder has long had a 32 MiB line limit. Preserve it, reserving 1 KiB
// for the data prefix and JSON envelope instead of treating its removal as a
// regression fix. JSON bodies and accumulated multi-event text are independent.
const largeStreamingPayloadSize = 32*1024*1024 - 1024

func TestLargeResponsesPayloadContract(t *testing.T) {
	for _, streaming := range []bool{false, true} {
		name := "blocking JSON"
		if streaming {
			name = "streaming event"
		}
		t.Run(name, func(t *testing.T) {
			payloadSize := largePayloadSize
			if streaming {
				payloadSize = largeStreamingPayloadSize
			}
			text := strings.Repeat("x", payloadSize)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost || r.URL.Path != "/responses" {
					t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
					http.NotFound(w, r)
					return
				}
				if streaming {
					w.Header().Set("Content-Type", "text/event-stream")
					writeLargePayload(t, w, "event: response.output_text.done\ndata: {\"type\":\"response.output_text.done\",\"text\":\"")
					writeLargePayload(t, w, text)
					writeLargePayload(t, w, "\",\"item_id\":\"msg_test\",\"output_index\":0,\"content_index\":0,\"sequence_number\":1,\"logprobs\":[]}\n\ndata: [DONE]\n\n")
				} else {
					w.Header().Set("Content-Type", "application/json")
					writeLargePayload(t, w, `{"id":"resp_test","object":"response","status":"completed","output":[{"type":"message","id":"msg_test","role":"assistant","status":"completed","content":[{"type":"output_text","text":"`)
					writeLargePayload(t, w, text)
					writeLargePayload(t, w, `","annotations":[]}]}]}`)
				}
			}))
			defer server.Close()
			client := openai.NewClient(option.WithAPIKey("test-key"), option.WithBaseURL(server.URL), option.WithMaxRetries(0))
			params := responses.ResponseNewParams{Model: "gpt-4o-mini", Input: responses.ResponseNewParamsInputUnion{OfString: openai.String("Hello")}}
			if streaming {
				stream := client.Responses.NewStreaming(context.Background(), params)
				defer func() { _ = stream.Close() }()
				if !stream.Next() {
					t.Fatalf("large event was not delivered: %v", stream.Err())
				}
				event := stream.Current()
				if event.Type != "response.output_text.done" || event.AsResponseOutputTextDone().Text != text {
					t.Fatal("large streaming event was truncated or changed")
				}
				if stream.Next() || stream.Err() != nil {
					t.Fatalf("unexpected end of stream: %v", stream.Err())
				}
			} else {
				response, err := client.Responses.New(context.Background(), params)
				if err != nil {
					t.Fatal(err)
				}
				if response.OutputText() != text {
					t.Fatal("large blocking JSON message was truncated or changed")
				}
			}
		})
	}
}

func TestLargeChatCompletionPayloadContract(t *testing.T) {
	for _, streaming := range []bool{false, true} {
		name := "blocking JSON"
		if streaming {
			name = "streaming accumulation"
		}
		t.Run(name, func(t *testing.T) {
			text := strings.Repeat("x", largePayloadSize)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost || r.URL.Path != "/chat/completions" {
					t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
					http.NotFound(w, r)
					return
				}
				if streaming {
					w.Header().Set("Content-Type", "text/event-stream")
					// Accumulation is a separate contract from single-event decoding:
					// each chunk fits below 32 MiB, but the final content does not.
					for _, part := range []string{text[:len(text)/2], text[len(text)/2:]} {
						writeLargePayload(t, w, `data: {"id":"chatcmpl_test","object":"chat.completion.chunk","choices":[{"index":0,"delta":{"content":"`)
						writeLargePayload(t, w, part)
						writeLargePayload(t, w, "\"}}]}\n\n")
					}
					writeLargePayload(t, w, "data: {\"id\":\"chatcmpl_test\",\"object\":\"chat.completion.chunk\",\"choices\":[{\"index\":0,\"delta\":{},\"finish_reason\":\"stop\"}]}\n\ndata: [DONE]\n\n")
				} else {
					w.Header().Set("Content-Type", "application/json")
					writeLargePayload(t, w, `{"id":"chatcmpl_test","object":"chat.completion","choices":[{"index":0,"message":{"role":"assistant","content":"`)
					writeLargePayload(t, w, text)
					writeLargePayload(t, w, `"},"finish_reason":"stop"}]}`)
				}
			}))
			defer server.Close()
			client := openai.NewClient(option.WithAPIKey("test-key"), option.WithBaseURL(server.URL), option.WithMaxRetries(0))
			params := openai.ChatCompletionNewParams{Model: "gpt-4o-mini", Messages: []openai.ChatCompletionMessageParamUnion{openai.UserMessage("Hello")}}
			if streaming {
				stream := client.Chat.Completions.NewStreaming(context.Background(), params)
				defer func() { _ = stream.Close() }()
				acc := openai.ChatCompletionAccumulator{}
				finished := false
				for stream.Next() {
					if !acc.AddChunk(stream.Current()) {
						t.Fatal("large content was rejected by the accumulator")
					}
					if content, ok := acc.JustFinishedContent(); ok {
						finished = true
						if content != text {
							t.Fatal("large finished content was truncated or changed")
						}
					}
				}
				if err := stream.Err(); err != nil {
					t.Fatal(err)
				}
				if !finished || len(acc.Choices) != 1 || acc.Choices[0].Message.Content != text {
					t.Fatal("large accumulated message was not completed intact")
				}
			} else {
				completion, err := client.Chat.Completions.New(context.Background(), params)
				if err != nil {
					t.Fatal(err)
				}
				if len(completion.Choices) != 1 || completion.Choices[0].Message.Content != text {
					t.Fatal("large blocking JSON message was truncated or changed")
				}
			}
		})
	}
}

func writeLargePayload(t *testing.T, w io.Writer, part string) {
	t.Helper()
	if _, err := io.WriteString(w, part); err != nil {
		t.Errorf("write synthetic response: %v", err)
	}
}

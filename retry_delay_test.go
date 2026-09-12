package openai_test

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/shared"
)

type retryDelayRoundTripperFunc func(*http.Request) (*http.Response, error)

func (f retryDelayRoundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestRetryAfterClampsToConfiguredMaximum(t *testing.T) {
	attempts := 0
	client := openai.NewClient(
		option.WithAPIKey("My API Key"),
		option.WithMaxRetries(1),
		option.WithMaxRetryDelay(10*time.Millisecond),
		option.WithHTTPClient(&http.Client{
			Transport: retryDelayRoundTripperFunc(func(req *http.Request) (*http.Response, error) {
				attempts++
				return &http.Response{
					StatusCode: http.StatusTooManyRequests,
					Header: http.Header{
						http.CanonicalHeaderKey("Retry-After"): []string{"31536000"},
					},
					Body: http.NoBody,
				}, nil
			}),
		}),
	)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{{
			OfUser: &openai.ChatCompletionUserMessageParam{
				Content: openai.ChatCompletionUserMessageParamContentUnion{
					OfString: openai.String("Say this is a test"),
				},
			},
		}},
		Model: shared.ChatModelGPT4o,
	})
	if err == nil {
		t.Fatal("expected retry response to return an error")
	}
	if attempts != 2 {
		t.Fatalf("attempts = %d, want 2", attempts)
	}
}

func TestRetryAfterDelayEdgeCases(t *testing.T) {
	tests := map[string]struct {
		header        http.Header
		maxRetryDelay time.Duration
		timeout       time.Duration
		wantAttempts  int
	}{
		"zero seconds": {
			header:       http.Header{"Retry-After": {"0"}},
			timeout:      250 * time.Millisecond,
			wantAttempts: 2,
		},
		"zero milliseconds": {
			header:       http.Header{"Retry-After-Ms": {"0"}},
			timeout:      250 * time.Millisecond,
			wantAttempts: 2,
		},
		"elapsed HTTP date": {
			header:       http.Header{"Retry-After": {time.Now().Add(-time.Minute).UTC().Format(time.RFC1123)}},
			timeout:      250 * time.Millisecond,
			wantAttempts: 2,
		},
		"finite scaling overflow": {
			header:        http.Header{"Retry-After": {"2" + strings.Repeat("0", 299)}},
			maxRetryDelay: 700 * time.Millisecond,
			timeout:       600 * time.Millisecond,
			wantAttempts:  1,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			attempts := 0
			opts := []option.RequestOption{
				option.WithAPIKey("My API Key"),
				option.WithMaxRetries(1),
				option.WithHTTPClient(&http.Client{Transport: retryDelayRoundTripperFunc(func(req *http.Request) (*http.Response, error) {
					attempts++
					return &http.Response{
						StatusCode: http.StatusTooManyRequests,
						Header:     test.header,
						Body:       http.NoBody,
					}, nil
				})}),
			}
			if test.maxRetryDelay > 0 {
				opts = append(opts, option.WithMaxRetryDelay(test.maxRetryDelay))
			}
			client := openai.NewClient(opts...)
			ctx, cancel := context.WithTimeout(context.Background(), test.timeout)
			t.Cleanup(cancel)

			_, _ = client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
				Messages: []openai.ChatCompletionMessageParamUnion{{
					OfUser: &openai.ChatCompletionUserMessageParam{
						Content: openai.ChatCompletionUserMessageParamContentUnion{
							OfString: openai.String("Say this is a test"),
						},
					},
				}},
				Model: shared.ChatModelGPT4o,
			})
			if attempts != test.wantAttempts {
				t.Fatalf("attempts = %d, want %d", attempts, test.wantAttempts)
			}
		})
	}
}

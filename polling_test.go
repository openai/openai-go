package openai_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func pollResponse(headerValue string, body string) *http.Response {
	header := http.Header{"Content-Type": {"application/json"}}
	if headerValue != "" {
		header.Set("openai-poll-after-ms", headerValue)
	}
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     header,
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func TestPollStatusRemoteDelayObservesContext(t *testing.T) {
	tests := map[string]struct {
		header string
		body   string
		poll   func(context.Context, *openai.Client) error
	}{
		"vector store file positive delay": {
			header: "1000",
			body:   `{"id":"file_123","status":"in_progress","vector_store_id":"vs_123"}`,
			poll: func(ctx context.Context, client *openai.Client) error {
				_, err := client.VectorStores.Files.PollStatus(ctx, "vs_123", "file_123", 0)
				return err
			},
		},
		"vector store batch negative delay": {
			header: "-1",
			body:   `{"id":"batch_123","status":"in_progress","vector_store_id":"vs_123"}`,
			poll: func(ctx context.Context, client *openai.Client) error {
				_, err := client.VectorStores.FileBatches.PollStatus(ctx, "vs_123", "batch_123", 0)
				return err
			},
		},
		"video overflowing delay": {
			header: "999999999999999999999999",
			body:   `{"id":"video_123","status":"in_progress"}`,
			poll: func(ctx context.Context, client *openai.Client) error {
				_, err := client.Videos.PollStatus(ctx, "video_123", 0)
				return err
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			attempts := 0
			client := openai.NewClient(
				option.WithAPIKey("test-key"),
				option.WithBaseURL("https://example.com/"),
				option.WithMaxRetries(0),
				option.WithHTTPClient(&http.Client{Transport: &closureTransport{
					fn: func(*http.Request) (*http.Response, error) {
						attempts++
						return pollResponse(test.header, test.body), nil
					},
				}}),
			)
			ctx, cancel := context.WithTimeout(context.Background(), 25*time.Millisecond)
			defer cancel()
			started := time.Now()

			err := test.poll(ctx, &client)
			if !errors.Is(err, context.DeadlineExceeded) {
				t.Fatalf("PollStatus() error = %v, want %v", err, context.DeadlineExceeded)
			}
			if elapsed := time.Since(started); elapsed >= 500*time.Millisecond {
				t.Fatalf("PollStatus() returned after %s, want context cancellation before 500ms", elapsed)
			}
			if attempts != 1 {
				t.Fatalf("attempts = %d, want 1", attempts)
			}
		})
	}
}

func TestPollStatusUsesValidServerDelay(t *testing.T) {
	attempts := 0
	client := openai.NewClient(
		option.WithAPIKey("test-key"),
		option.WithBaseURL("https://example.com/"),
		option.WithMaxRetries(0),
		option.WithHTTPClient(&http.Client{Transport: &closureTransport{
			fn: func(*http.Request) (*http.Response, error) {
				attempts++
				status := "in_progress"
				if attempts == 2 {
					status = "completed"
				}
				return pollResponse("1", `{"id":"file_123","status":"`+status+`","vector_store_id":"vs_123"}`), nil
			},
		}}),
	)

	file, err := client.VectorStores.Files.PollStatus(context.Background(), "vs_123", "file_123", 0)
	if err != nil {
		t.Fatal(err)
	}
	if file.Status != openai.VectorStoreFileStatusCompleted {
		t.Fatalf("status = %q, want %q", file.Status, openai.VectorStoreFileStatusCompleted)
	}
	if attempts != 2 {
		t.Fatalf("attempts = %d, want 2", attempts)
	}
}

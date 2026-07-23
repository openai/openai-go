package openai

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/openai/openai-go/v3/option"
)

type pollingHTTPDoerFunc func(*http.Request) (*http.Response, error)

func (f pollingHTTPDoerFunc) Do(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestVideoPollStatusReturnsWhenContextCancelled(t *testing.T) {
	client := NewClient(
		option.WithAPIKey("test-key"),
		option.WithHTTPClient(pollingHTTPDoerFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       io.NopCloser(strings.NewReader(`{"id":"video_123","status":"in_progress"}`)),
				Request:    req,
			}, nil
		})),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	t.Cleanup(cancel)

	start := time.Now()
	_, err := client.Videos.PollStatus(ctx, "video_123", 1000)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("PollStatus() error = %v, want %v", err, context.DeadlineExceeded)
	}
	if elapsed := time.Since(start); elapsed > 500*time.Millisecond {
		t.Fatalf("PollStatus() took %s after context cancellation", elapsed)
	}
}

package openai_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func TestBytesReaderRequestPreservesRemainingContent(t *testing.T) {
	for _, test := range []struct {
		name   string
		offset int64
		want   string
	}{
		{name: "partial reader", offset: 7, want: "payload"},
		{name: "consumed reader", offset: 14, want: ""},
	} {
		t.Run(test.name, func(t *testing.T) {
			var received []string
			var mu sync.Mutex
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Error(err)
				}
				mu.Lock()
				received = append(received, string(body))
				attempt := len(received)
				mu.Unlock()
				w.Header().Set("Content-Type", "application/json")
				if attempt == 1 {
					w.Header().Set("Retry-After", "0")
					w.WriteHeader(http.StatusTooManyRequests)
					_, _ = io.WriteString(w, `{"error":{"message":"retry"}}`)
					return
				}
				_, _ = io.WriteString(w, `{"id":"completion-test","object":"chat.completion","choices":[]}`)
			}))
			defer server.Close()
			reader := bytes.NewReader([]byte("prefix:payload"))
			if _, err := reader.Seek(test.offset, io.SeekStart); err != nil {
				t.Fatal(err)
			}
			client := openai.NewClient(
				option.WithBaseURL(server.URL), option.WithAPIKey("test-key"), option.WithMaxRetries(1),
			)
			_, err := client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{},
				option.WithRequestBody("application/octet-stream", reader))
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}
			mu.Lock()
			defer mu.Unlock()
			if len(received) != 2 {
				t.Fatalf("received %d requests, want 2", len(received))
			}
			for i, body := range received {
				if body != test.want {
					t.Errorf("request %d body = %q, want %q", i, body, test.want)
				}
			}
		})
	}
}

func TestBytesReaderReplayBodiesHaveIndependentPositions(t *testing.T) {
	reader := bytes.NewReader([]byte("payload"))
	var received string
	var mu sync.Mutex
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
		}
		mu.Lock()
		received = string(body)
		mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"id":"completion-test","object":"chat.completion","choices":[]}`)
	}))
	defer server.Close()
	client := openai.NewClient(
		option.WithBaseURL(server.URL), option.WithAPIKey("test-key"), option.WithMaxRetries(0),
		option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
			if req.GetBody == nil {
				t.Fatal("byte reader request has no replay factory")
			}
			first, err := req.GetBody()
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				if closeErr := first.Close(); closeErr != nil {
					t.Error(closeErr)
				}
			}()
			second, err := req.GetBody()
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				if closeErr := second.Close(); closeErr != nil {
					t.Error(closeErr)
				}
			}()
			prefix := make([]byte, 3)
			if _, readErr := io.ReadFull(first, prefix); readErr != nil {
				t.Fatal(readErr)
			}
			if string(prefix) != "pay" {
				t.Errorf("first replay prefix = %q", prefix)
			}
			all, err := io.ReadAll(second)
			if err != nil {
				t.Fatal(err)
			}
			if string(all) != "payload" {
				t.Errorf("second replay = %q, want payload", all)
			}
			rest, err := io.ReadAll(first)
			if err != nil {
				t.Fatal(err)
			}
			if string(rest) != "load" {
				t.Errorf("first replay suffix = %q, want load", rest)
			}
			return next(req)
		}),
	)
	_, err := client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{},
		option.WithRequestBody("application/octet-stream", reader))
	if err != nil {
		t.Fatalf("request failed after opening replay readers: %v", err)
	}
	if reader.Len() != 0 {
		t.Errorf("initial request did not consume the supplied reader")
	}
	mu.Lock()
	defer mu.Unlock()
	if received != "payload" {
		t.Errorf("original request body = %q, want payload", received)
	}
}

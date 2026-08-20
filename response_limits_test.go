package openai_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type responseDoerFunc func(*http.Request) (*http.Response, error)

func (f responseDoerFunc) Do(req *http.Request) (*http.Response, error) {
	return f(req)
}

type countingResponseBody struct {
	reader  io.Reader
	endless bool
	reads   int
	closed  bool
}

func (b *countingResponseBody) Read(p []byte) (int, error) {
	var n int
	var err error
	if b.endless {
		for i := range p {
			p[i] = 'x'
		}
		n = len(p)
	} else {
		n, err = b.reader.Read(p)
	}
	b.reads += n
	return n, err
}

func (b *countingResponseBody) Close() error {
	b.closed = true
	return nil
}

type blockingResponseBody struct {
	closed chan struct{}
	once   sync.Once
	mu     sync.Mutex
	closes int
}

func newBlockingResponseBody() *blockingResponseBody {
	return &blockingResponseBody{closed: make(chan struct{})}
}

func (b *blockingResponseBody) Read([]byte) (int, error) {
	<-b.closed
	return 0, errors.New("response body closed")
}

func (b *blockingResponseBody) Close() error {
	b.mu.Lock()
	b.closes++
	b.mu.Unlock()
	b.once.Do(func() { close(b.closed) })
	return nil
}

func (b *blockingResponseBody) closeCount() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.closes
}

func newResponseLimitClient(body io.ReadCloser, status int, contentType string, opts ...option.RequestOption) openai.Client {
	opts = append([]option.RequestOption{
		option.WithAPIKey("test-key"),
		option.WithMaxRetries(0),
		option.WithHTTPClient(responseDoerFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: status,
				Header:     http.Header{"Content-Type": {contentType}},
				Body:       body,
				Request:    req,
			}, nil
		})),
	}, opts...)
	return openai.NewClient(opts...)
}

func TestExecuteBoundsTypedSuccessResponse(t *testing.T) {
	body := &countingResponseBody{endless: true}
	client := newResponseLimitClient(body, http.StatusOK, "application/json")
	var response struct {
		OK bool `json:"ok"`
	}

	err := client.Execute(
		context.Background(),
		http.MethodGet,
		"test",
		nil,
		&response,
		option.WithMaxResponseBodyBytes(8),
	)
	if err == nil || !strings.Contains(err.Error(), "exceeded configured limit of 8 bytes") {
		t.Fatalf("Execute() error = %v, want response body limit error", err)
	}
	if body.reads != 9 {
		t.Fatalf("response bytes read = %d, want 9", body.reads)
	}
	if !body.closed {
		t.Fatal("response body was not closed")
	}
}

func TestExecuteBoundsErrorResponseRequestedAsRaw(t *testing.T) {
	errorJSON := `{"error":{"message":"bad","type":"invalid_request_error","param":"p","code":"c"}}`
	body := &countingResponseBody{reader: strings.NewReader(errorJSON + strings.Repeat(" ", 1024))}
	client := newResponseLimitClient(body, http.StatusBadRequest, "application/json")
	var raw *http.Response

	err := client.Execute(
		context.Background(),
		http.MethodGet,
		"test",
		nil,
		&raw,
		option.WithMaxErrorResponseBodyBytes(int64(len(errorJSON))),
	)
	if err == nil || !strings.Contains(err.Error(), "exceeded configured limit") {
		t.Fatalf("Execute() error = %v, want error response body limit error", err)
	}
	var apiErr *openai.Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("errors.As(%T, *openai.Error) = false", err)
	}
	if apiErr.StatusCode != http.StatusBadRequest {
		t.Fatalf("API error status = %d, want %d", apiErr.StatusCode, http.StatusBadRequest)
	}
	if raw != apiErr.Response {
		t.Fatal("raw response and API error response differ")
	}
	retained, readErr := io.ReadAll(raw.Body)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if string(retained) != errorJSON {
		t.Fatalf("retained error body = %q, want %q", retained, errorJSON)
	}
	if body.reads != len(errorJSON)+1 {
		t.Fatalf("response bytes read = %d, want %d", body.reads, len(errorJSON)+1)
	}
}

func TestExecutePreservesOrdinaryAPIError(t *testing.T) {
	errorJSON := `{"error":{"message":"bad","type":"invalid_request_error","param":"p","code":"c"}}`
	body := &countingResponseBody{reader: strings.NewReader(errorJSON)}
	client := newResponseLimitClient(body, http.StatusBadRequest, "application/json")
	var response map[string]any

	err := client.Execute(context.Background(), http.MethodGet, "test", nil, &response)
	var apiErr *openai.Error
	if reflect.TypeOf(err) != reflect.TypeOf(apiErr) {
		t.Fatalf("Execute() error type = %T, want *openai.Error", err)
	}
	if !errors.As(err, &apiErr) {
		t.Fatalf("errors.As(%T, *openai.Error) = false", err)
	}
	if apiErr.StatusCode != http.StatusBadRequest {
		t.Fatalf("API error status = %d, want %d", apiErr.StatusCode, http.StatusBadRequest)
	}
	retained, readErr := io.ReadAll(apiErr.Response.Body)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if string(retained) != errorJSON {
		t.Fatalf("retained error body = %q, want %q", retained, errorJSON)
	}
}

func TestExecuteAllowsMethodResponseLimitOverride(t *testing.T) {
	body := &countingResponseBody{reader: strings.NewReader(`{"ok":true}`)}
	client := newResponseLimitClient(
		body,
		http.StatusOK,
		"application/json",
		option.WithMaxResponseBodyBytes(4),
	)
	var response struct {
		OK bool `json:"ok"`
	}

	err := client.Execute(
		context.Background(),
		http.MethodGet,
		"test",
		nil,
		&response,
		option.WithMaxResponseBodyBytes(64),
	)
	if err != nil {
		t.Fatal(err)
	}
	if !response.OK {
		t.Fatal("typed response was not decoded")
	}
}

func TestExecuteLeavesRawSuccessResponseStreaming(t *testing.T) {
	body := &countingResponseBody{endless: true}
	client := newResponseLimitClient(
		body,
		http.StatusOK,
		"application/octet-stream",
		option.WithMaxResponseBodyBytes(1),
		option.WithResponseBodyTimeout(time.Nanosecond),
	)
	var raw *http.Response

	if err := client.Execute(context.Background(), http.MethodGet, "test", nil, &raw); err != nil {
		t.Fatal(err)
	}
	if body.reads != 0 {
		t.Fatalf("response bytes read = %d, want 0", body.reads)
	}
	if body.closed {
		t.Fatal("raw response body was closed before handoff")
	}
	if err := raw.Body.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestExecuteTimesOutNonStreamingResponseBodies(t *testing.T) {
	for _, status := range []int{http.StatusOK, http.StatusBadRequest} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			body := newBlockingResponseBody()
			client := newResponseLimitClient(
				body,
				status,
				"application/json",
				option.WithResponseBodyTimeout(10*time.Millisecond),
			)
			var response map[string]any
			done := make(chan error, 1)
			go func() {
				done <- client.Execute(context.Background(), http.MethodGet, "test", nil, &response)
			}()

			select {
			case err := <-done:
				if !errors.Is(err, context.DeadlineExceeded) {
					t.Fatalf("Execute() error = %v, want context deadline exceeded", err)
				}
				if body.closeCount() != 1 {
					t.Fatalf("response body close count = %d, want 1", body.closeCount())
				}
			case <-time.After(time.Second):
				_ = body.Close()
				t.Fatal("Execute() did not enforce the response body timeout")
			}
		})
	}
}

func TestExecuteTimeoutInterruptsSlowHTTPResponse(t *testing.T) {
	headersSent := make(chan struct{})
	releaseBody := make(chan struct{})
	var releaseOnce sync.Once
	release := func() { releaseOnce.Do(func() { close(releaseBody) }) }

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.(http.Flusher).Flush()
		close(headersSent)
		<-releaseBody
		_, _ = io.WriteString(w, `{"ok":true}`)
	}))
	defer func() {
		release()
		server.Close()
	}()

	client := openai.NewClient(
		option.WithAPIKey("test-key"),
		option.WithBaseURL(server.URL+"/"),
		option.WithMaxRetries(0),
		option.WithResponseBodyTimeout(10*time.Millisecond),
	)
	var response struct {
		OK bool `json:"ok"`
	}
	done := make(chan error, 1)
	go func() {
		done <- client.Get(context.Background(), "slow", nil, &response)
	}()
	<-headersSent

	select {
	case err := <-done:
		release()
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("Get() error = %v, want context deadline exceeded", err)
		}
	case <-time.After(time.Second):
		release()
		<-done
		t.Fatal("response body timeout did not interrupt a stalled standard HTTP response")
	}
}

package openai_test

import (
	"bytes"
	"compress/gzip"
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

type responseRoundTripperFunc func(*http.Request) (*http.Response, error)

func (f responseRoundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

type delegatingResponseDoer struct {
	*http.Client
}

type countingResponseBody struct {
	reader        io.Reader
	endless       bool
	reads         int
	closed        bool
	closeFinished chan struct{}
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
	if b.closeFinished != nil {
		close(b.closeFinished)
	}
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

type contextResponseBody struct {
	ctx context.Context
	mu  sync.Mutex

	closes    int
	closed    chan struct{}
	closeOnce sync.Once
}

func (b *contextResponseBody) Read([]byte) (int, error) {
	<-b.ctx.Done()
	return 0, b.ctx.Err()
}

func (b *contextResponseBody) Close() error {
	b.mu.Lock()
	b.closes++
	b.mu.Unlock()
	b.closeOnce.Do(func() { close(b.closed) })
	return nil
}

func (b *contextResponseBody) closeCount() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.closes
}

type closeContextResponseBody struct {
	ctx    context.Context
	reader io.Reader
	err    error

	closed             bool
	contextDoneOnClose bool
}

func (b *closeContextResponseBody) Read(p []byte) (int, error) {
	if b.reader != nil {
		return b.reader.Read(p)
	}
	return 0, io.EOF
}

func (b *closeContextResponseBody) Close() error {
	b.closed = true
	b.contextDoneOnClose = b.ctx.Err() != nil
	return b.err
}

type closeContextResponseDoer struct {
	ctx      context.Context
	body     *closeContextResponseBody
	closeErr error
	status   int
	contents string
}

func (d *closeContextResponseDoer) Do(req *http.Request) (*http.Response, error) {
	d.ctx = req.Context()
	d.body = &closeContextResponseBody{ctx: d.ctx, reader: strings.NewReader(d.contents), err: d.closeErr}
	status := d.status
	if status == 0 {
		status = http.StatusOK
	}
	return &http.Response{
		StatusCode: status,
		Header:     http.Header{"Content-Type": {"application/json"}},
		Body:       d.body,
		Request:    req,
	}, nil
}

func newCloseContextResponseClient(closeErr error) (openai.Client, *closeContextResponseDoer) {
	doer := &closeContextResponseDoer{closeErr: closeErr}
	client := openai.NewClient(
		option.WithAPIKey("test-key"),
		option.WithMaxRetries(0),
		option.WithHTTPClient(doer),
	)
	return client, doer
}

type blockingCloseContextResponseBody struct {
	ctx           context.Context
	releaseClose  chan struct{}
	closeFinished chan struct{}
	closeOnce     sync.Once

	contextDoneOnClose bool
}

func newBlockingCloseContextResponseBody(ctx context.Context) *blockingCloseContextResponseBody {
	return &blockingCloseContextResponseBody{
		ctx:           ctx,
		releaseClose:  make(chan struct{}),
		closeFinished: make(chan struct{}),
	}
}

func (b *blockingCloseContextResponseBody) Read([]byte) (int, error) {
	<-b.ctx.Done()
	return 0, b.ctx.Err()
}

func (b *blockingCloseContextResponseBody) Close() error {
	b.closeOnce.Do(func() {
		b.contextDoneOnClose = b.ctx.Err() != nil
		<-b.releaseClose
		close(b.closeFinished)
	})
	return nil
}

type repeatedByteReader byte

func (b repeatedByteReader) Read(p []byte) (int, error) {
	for i := range p {
		p[i] = byte(b)
	}
	return len(p), nil
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
	body := &countingResponseBody{endless: true, closeFinished: make(chan struct{})}
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
	select {
	case <-body.closeFinished:
	case <-time.After(time.Second):
		t.Fatal("response body was not closed")
	}
}

func TestExecuteBoundsCompressedWireResponse(t *testing.T) {
	const wireLimit = 64
	var compressed bytes.Buffer
	for i := 0; i < 32; i++ {
		zw := gzip.NewWriter(&compressed)
		if i == 0 {
			if _, err := io.WriteString(zw, "{}"); err != nil {
				t.Fatal(err)
			}
		}
		if err := zw.Close(); err != nil {
			t.Fatal(err)
		}
	}
	if compressed.Len() <= wireLimit {
		t.Fatalf("compressed fixture length = %d, want more than wire limit", compressed.Len())
	}

	for _, status := range []int{http.StatusOK, http.StatusBadRequest} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if got := r.Header.Get("Accept-Encoding"); got != "gzip" {
					t.Errorf("Accept-Encoding = %q, want gzip", got)
				}
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Content-Encoding", "gzip")
				w.WriteHeader(status)
				_, _ = w.Write(compressed.Bytes())
			}))
			defer server.Close()

			clients := []struct {
				name string
				opts []option.RequestOption
			}{
				{name: "native client"},
				{
					name: "delegating doer",
					opts: []option.RequestOption{
						option.WithHTTPClient(&delegatingResponseDoer{Client: &http.Client{}}),
					},
				},
				{
					name: "delegating doer replaces compression-disabled client",
					opts: []option.RequestOption{
						option.WithHTTPClient(&http.Client{
							Transport: &http.Transport{DisableCompression: true},
						}),
						option.WithHTTPClient(&delegatingResponseDoer{Client: &http.Client{}}),
					},
				},
			}
			for _, test := range clients {
				t.Run(test.name, func(t *testing.T) {
					opts := []option.RequestOption{
						option.WithAPIKey("test-key"),
						option.WithBaseURL(server.URL + "/"),
						option.WithMaxRetries(0),
						option.WithMaxResponseBodyBytes(wireLimit),
						option.WithMaxErrorResponseBodyBytes(wireLimit),
					}
					opts = append(opts, test.opts...)

					client := openai.NewClient(opts...)
					var response map[string]any
					err := client.Get(context.Background(), "compressed", nil, &response)
					if err == nil || !strings.Contains(err.Error(), "compressed response body exceeded configured limit of 64 bytes") {
						t.Fatalf("Get() error = %v, want compressed response body limit error", err)
					}
					if status >= http.StatusBadRequest {
						var apiErr *openai.Error
						if !errors.As(err, &apiErr) {
							t.Fatalf("errors.As(%T, *openai.Error) = false", err)
						}
					}
				})
			}
		})
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

func TestExecutePreservesTransparentGzipForRawResponse(t *testing.T) {
	var compressed bytes.Buffer
	zw := gzip.NewWriter(&compressed)
	if _, err := io.WriteString(zw, `{"ok":true}`); err != nil {
		t.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Encoding", "gzip")
		_, _ = w.Write(compressed.Bytes())
	}))
	defer server.Close()

	client := openai.NewClient(
		option.WithAPIKey("test-key"),
		option.WithBaseURL(server.URL+"/"),
		option.WithMaxRetries(0),
		option.WithMaxResponseBodyBytes(1),
	)
	var raw *http.Response
	if err := client.Get(context.Background(), "compressed", nil, &raw); err != nil {
		t.Fatal(err)
	}
	if !raw.Uncompressed {
		t.Fatal("raw response was not marked as transparently decompressed")
	}
	if got := raw.Header.Get("Content-Encoding"); got != "" {
		t.Fatalf("Content-Encoding = %q, want removed after transparent decompression", got)
	}
	if raw.ContentLength != -1 {
		t.Fatalf("ContentLength = %d, want -1 after transparent decompression", raw.ContentLength)
	}
	contents, err := io.ReadAll(raw.Body)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := string(contents), `{"ok":true}`; got != want {
		t.Fatalf("raw body = %q, want %q", got, want)
	}
	if err := raw.Body.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestExecuteRawResponseOwnsAttemptContextLifecycle(t *testing.T) {
	var attemptCtx context.Context
	body := &countingResponseBody{reader: strings.NewReader("{}")}
	client := openai.NewClient(
		option.WithAPIKey("test-key"),
		option.WithMaxRetries(0),
		option.WithHTTPClient(responseDoerFunc(func(req *http.Request) (*http.Response, error) {
			attemptCtx = req.Context()
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": {"application/json"}},
				Body:       body,
				Request:    req,
			}, nil
		})),
	)
	var raw *http.Response
	if err := client.Get(context.Background(), "test", nil, &raw); err != nil {
		t.Fatal(err)
	}
	select {
	case <-attemptCtx.Done():
		t.Fatal("attempt context canceled before raw response ownership ended")
	default:
	}
	if err := raw.Body.Close(); err != nil {
		t.Fatal(err)
	}
	select {
	case <-attemptCtx.Done():
	case <-time.After(time.Second):
		t.Fatal("raw response close did not end the attempt context lifecycle")
	}
}

func TestExecuteRawResponseClosesBodyBeforeCancelingAttemptContext(t *testing.T) {
	closeErr := errors.New("close failed")
	client, doer := newCloseContextResponseClient(closeErr)

	var raw *http.Response
	if err := client.Get(context.Background(), "test", nil, &raw); err != nil {
		t.Fatal(err)
	}
	if err := raw.Body.Close(); !errors.Is(err, closeErr) {
		t.Fatalf("raw response Close() error = %v, want %v", err, closeErr)
	}
	if !doer.body.closed {
		t.Fatal("raw response Close() did not close the underlying body")
	}
	if doer.body.contextDoneOnClose {
		t.Fatal("attempt context was canceled before the underlying body was closed")
	}
	select {
	case <-doer.ctx.Done():
	case <-time.After(time.Second):
		t.Fatal("attempt context was not canceled after a failed underlying Close")
	}
}

func TestExecuteRawReadAllRetainsAttemptContextUntilClose(t *testing.T) {
	client, doer := newCloseContextResponseClient(nil)
	doer.contents = `{}`

	var raw *http.Response
	if err := client.Get(context.Background(), "test", nil, &raw); err != nil {
		t.Fatal(err)
	}
	if _, err := io.ReadAll(raw.Body); err != nil {
		t.Fatal(err)
	}
	select {
	case <-doer.ctx.Done():
		t.Fatal("attempt context canceled at EOF before raw body Close")
	default:
	}
	if err := raw.Body.Close(); err != nil {
		t.Fatal(err)
	}
	if doer.body.contextDoneOnClose {
		t.Fatal("attempt context was canceled before the underlying body was closed")
	}
	select {
	case <-doer.ctx.Done():
	case <-time.After(time.Second):
		t.Fatal("attempt context was not canceled after raw body Close")
	}
}

func TestExecuteClosesConsumedBodyBeforeCancelingAttemptContext(t *testing.T) {
	for _, status := range []int{http.StatusOK, http.StatusBadRequest} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			closeErr := errors.New("close failed")
			client, doer := newCloseContextResponseClient(closeErr)
			doer.status = status
			if status == http.StatusOK {
				doer.contents = `{}`
			} else {
				doer.contents = `{"error":{"message":"bad"}}`
			}

			var response map[string]any
			err := client.Get(context.Background(), "test", nil, &response)
			if status == http.StatusOK && err != nil {
				t.Fatal(err)
			}
			if status != http.StatusOK && err == nil {
				t.Fatal("Get() error = nil, want API error")
			}
			if !doer.body.closed {
				t.Fatal("consumed response body was not closed")
			}
			if doer.body.contextDoneOnClose {
				t.Fatal("attempt context was canceled before the consumed body was closed")
			}
			select {
			case <-doer.ctx.Done():
			case <-time.After(time.Second):
				t.Fatal("attempt context was not canceled after a failed body Close")
			}
		})
	}
}

func TestExecuteClosesSuccessResponseWithoutOwner(t *testing.T) {
	client, doer := newCloseContextResponseClient(nil)

	if err := client.Get(context.Background(), "test", nil, nil); err != nil {
		t.Fatal(err)
	}
	if !doer.body.closed {
		t.Fatal("unowned success response body was not closed")
	}
	if doer.body.contextDoneOnClose {
		t.Fatal("attempt context was canceled before the unowned body was closed")
	}
	select {
	case <-doer.ctx.Done():
	case <-time.After(time.Second):
		t.Fatal("unowned success response did not end the attempt context lifecycle")
	}
}

func TestExecuteResponseIntoOwnsSuccessBody(t *testing.T) {
	client, doer := newCloseContextResponseClient(nil)

	var raw *http.Response
	if err := client.Get(
		context.Background(),
		"test",
		nil,
		nil,
		option.WithResponseInto(&raw),
	); err != nil {
		t.Fatal(err)
	}
	if doer.body.closed {
		t.Fatal("caller-owned response body was closed before handoff")
	}
	select {
	case <-doer.ctx.Done():
		t.Fatal("attempt context canceled before ResponseInto ownership ended")
	default:
	}
	if err := raw.Body.Close(); err != nil {
		t.Fatal(err)
	}
	if !doer.body.closed {
		t.Fatal("ResponseInto body Close() did not close the underlying body")
	}
	if doer.body.contextDoneOnClose {
		t.Fatal("attempt context was canceled before the ResponseInto body was closed")
	}
	select {
	case <-doer.ctx.Done():
	case <-time.After(time.Second):
		t.Fatal("ResponseInto body close did not end the attempt context lifecycle")
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

func TestExecuteBodyTimeoutCancelsCustomTransportContext(t *testing.T) {
	for _, status := range []int{http.StatusOK, http.StatusBadRequest} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			callerCtx, cancelCaller := context.WithCancel(context.Background())
			defer cancelCaller()

			var body *contextResponseBody
			client := openai.NewClient(
				option.WithAPIKey("test-key"),
				option.WithMaxRetries(0),
				option.WithResponseBodyTimeout(10*time.Millisecond),
				option.WithHTTPClient(&http.Client{
					Transport: responseRoundTripperFunc(func(req *http.Request) (*http.Response, error) {
						body = &contextResponseBody{ctx: req.Context(), closed: make(chan struct{})}
						return &http.Response{
							StatusCode: status,
							Header:     http.Header{"Content-Type": {"application/json"}},
							Body:       body,
							Request:    req,
						}, nil
					}),
				}),
			)
			var response map[string]any
			done := make(chan error, 1)
			go func() {
				done <- client.Get(callerCtx, "test", nil, &response)
			}()

			select {
			case err := <-done:
				if !errors.Is(err, context.DeadlineExceeded) {
					t.Fatalf("Get() error = %v, want context deadline exceeded", err)
				}
				select {
				case <-body.closed:
				case <-time.After(time.Second):
					t.Fatal("timed out waiting for response body cleanup")
				}
				if body.closeCount() != 1 {
					t.Fatalf("response body close count = %d, want 1", body.closeCount())
				}
			case <-time.After(250 * time.Millisecond):
				cancelCaller()
				<-done
				t.Fatal("response body timeout did not cancel the custom transport request context")
			}
		})
	}
}

func TestExecuteBodyTimeoutDoesNotWaitForBlockingClose(t *testing.T) {
	bodyReady := make(chan *blockingCloseContextResponseBody, 1)
	client := openai.NewClient(
		option.WithAPIKey("test-key"),
		option.WithMaxRetries(0),
		option.WithResponseBodyTimeout(10*time.Millisecond),
		option.WithHTTPClient(&http.Client{
			Transport: responseRoundTripperFunc(func(req *http.Request) (*http.Response, error) {
				body := newBlockingCloseContextResponseBody(req.Context())
				bodyReady <- body
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     http.Header{"Content-Type": {"application/json"}},
					Body:       body,
					Request:    req,
				}, nil
			}),
		}),
	)
	var response map[string]any
	done := make(chan error, 1)
	go func() {
		done <- client.Get(context.Background(), "test", nil, &response)
	}()
	body := <-bodyReady

	var err error
	returnedBeforeRelease := false
	select {
	case err = <-done:
		returnedBeforeRelease = true
	case <-time.After(250 * time.Millisecond):
	}
	close(body.releaseClose)
	if !returnedBeforeRelease {
		err = <-done
	}
	<-body.closeFinished

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Get() error = %v, want context deadline exceeded", err)
	}
	if !returnedBeforeRelease {
		t.Fatal("response body timeout waited for blocking Close")
	}
	if !body.contextDoneOnClose {
		t.Fatal("timeout Close began before the attempt context was canceled")
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

func TestImagesGenerateAllowsDocumented4KBase64Response(t *testing.T) {
	const (
		imageCount     = 3
		rawImageBytes  = int64(3840 * 2160 * 3)
		base64Bytes    = 4 * ((rawImageBytes + 2) / 3)
		responseHeader = `{"created":0,"data":[`
		imagePrefix    = `{"b64_json":"`
		imageSuffix    = `"}`
		responseFooter = `]}`
	)

	parts := []io.Reader{strings.NewReader(responseHeader)}
	expectedBytes := int64(len(responseHeader) + len(responseFooter))
	for i := 0; i < imageCount; i++ {
		if i != 0 {
			parts = append(parts, strings.NewReader(","))
			expectedBytes++
		}
		parts = append(
			parts,
			strings.NewReader(imagePrefix),
			io.LimitReader(repeatedByteReader('A'), base64Bytes),
			strings.NewReader(imageSuffix),
		)
		expectedBytes += int64(len(imagePrefix)+len(imageSuffix)) + base64Bytes
	}
	parts = append(parts, strings.NewReader(responseFooter))

	client := openai.NewClient(
		option.WithAPIKey("test-key"),
		option.WithMaxRetries(0),
		option.WithHTTPClient(responseDoerFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": {"application/json"}},
				Body:       io.NopCloser(io.MultiReader(parts...)),
				Request:    req,
			}, nil
		})),
	)
	var rawResponse []byte
	_, err := client.Images.Generate(
		context.Background(),
		openai.ImageGenerateParams{
			Prompt: "synthetic compatibility fixture",
			Model:  openai.ImageModelGPTImage2,
			N:      openai.Int(imageCount),
			Size:   openai.ImageGenerateParamsSize("3840x2160"),
		},
		option.WithResponseBodyInto(&rawResponse),
	)
	if err != nil {
		t.Fatalf("Images.Generate() error = %v, want documented multi-image response to fit the endpoint policy", err)
	}
	if int64(len(rawResponse)) != expectedBytes {
		t.Fatalf("response bytes = %d, want %d", len(rawResponse), expectedBytes)
	}
}

func TestImagesGenerateRespectsClientResponseLimitOverride(t *testing.T) {
	client := openai.NewClient(
		option.WithAPIKey("test-key"),
		option.WithMaxRetries(0),
		option.WithMaxResponseBodyBytes(4),
		option.WithHTTPClient(responseDoerFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": {"application/json"}},
				Body:       io.NopCloser(strings.NewReader(`{"data":[]}`)),
				Request:    req,
			}, nil
		})),
	)
	_, err := client.Images.Generate(context.Background(), openai.ImageGenerateParams{Prompt: "test"})
	if err == nil || !strings.Contains(err.Error(), "configured limit of 4 bytes") {
		t.Fatalf("Images.Generate() error = %v, want client response limit override", err)
	}
}

package openai_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type stalledHTTP2ResponseBody struct {
	io.ReadCloser
	attemptCtx    context.Context
	closeStarted  chan struct{}
	releaseClose  <-chan struct{}
	closeFinished chan struct{}
}

type canceledHTTP2ResponseBody struct {
	*stalledHTTP2ResponseBody
	endAttempt func()
}

func (b *canceledHTTP2ResponseBody) Read([]byte) (int, error) {
	b.endAttempt()
	return 0, b.attemptCtx.Err()
}

func (b *stalledHTTP2ResponseBody) Close() error {
	close(b.closeStarted)
	<-b.releaseClose
	err := b.ReadCloser.Close()
	close(b.closeFinished)
	return err
}

func wrapStalledHTTP2ResponseBody(
	req *http.Request,
	res *http.Response,
	release <-chan struct{},
) (*stalledHTTP2ResponseBody, error) {
	if res.ProtoMajor != 2 {
		_ = res.Body.Close()
		return nil, fmt.Errorf("response protocol = %s, want HTTP/2", res.Proto)
	}
	body := &stalledHTTP2ResponseBody{
		ReadCloser:    res.Body,
		attemptCtx:    req.Context(),
		closeStarted:  make(chan struct{}),
		releaseClose:  release,
		closeFinished: make(chan struct{}),
	}
	res.Body = body
	return body, nil
}

func newStalledHTTP2Server(t *testing.T, status int, release <-chan struct{}) *httptest.Server {
	t.Helper()
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.ProtoMajor != 2 {
			t.Errorf("request protocol = %s, want HTTP/2", req.Proto)
		}
		w.Header().Set("Content-Type", "application/json")
		if status == http.StatusTooManyRequests {
			w.Header().Set("Retry-After-Ms", "2000")
		}
		w.WriteHeader(status)
		if _, err := io.WriteString(w, `{"ok":true}`+strings.Repeat(" ", 256)); err != nil {
			t.Errorf("write response: %v", err)
			return
		}
		w.(http.Flusher).Flush()
		<-release
	}))
	server.EnableHTTP2 = true
	server.StartTLS()
	return server
}

func TestExecuteOverflowDoesNotWaitForStalledHTTP2Close(t *testing.T) {
	for _, status := range []int{http.StatusOK, http.StatusBadRequest} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			releaseClose := make(chan struct{})
			var releaseOnce sync.Once
			release := func() { releaseOnce.Do(func() { close(releaseClose) }) }

			server := newStalledHTTP2Server(t, status, releaseClose)
			defer func() {
				release()
				server.Close()
			}()

			bodyReady := make(chan *stalledHTTP2ResponseBody, 1)
			nativeTransport := server.Client().Transport
			client := openai.NewClient(
				option.WithAPIKey("test-key"),
				option.WithBaseURL(server.URL+"/"),
				option.WithMaxRetries(0),
				option.WithMaxResponseBodyBytes(8),
				option.WithMaxErrorResponseBodyBytes(8),
				option.WithResponseBodyTimeout(20*time.Millisecond),
				option.WithHTTPClient(&http.Client{
					Transport: responseRoundTripperFunc(func(req *http.Request) (*http.Response, error) {
						res, err := nativeTransport.RoundTrip(req)
						if err != nil {
							return nil, err
						}
						body, err := wrapStalledHTTP2ResponseBody(req, res, releaseClose)
						if err != nil {
							return nil, err
						}
						bodyReady <- body
						return res, nil
					}),
				}),
			)

			callerCtx, cancelCaller := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancelCaller()
			var response map[string]any
			done := make(chan error, 1)
			go func() {
				done <- client.Get(callerCtx, "overflow", nil, &response)
			}()

			var body *stalledHTTP2ResponseBody
			select {
			case body = <-bodyReady:
			case <-callerCtx.Done():
				t.Fatal("HTTP/2 response was not received")
			}
			select {
			case <-body.closeStarted:
			case <-callerCtx.Done():
				t.Fatal("overflow response cleanup was not started")
			}

			var err error
			returnedBeforeRelease := false
			select {
			case err = <-done:
				returnedBeforeRelease = true
			case <-time.After(250 * time.Millisecond):
			}
			if returnedBeforeRelease {
				select {
				case <-body.attemptCtx.Done():
				default:
					t.Fatal("overflow returned before canceling the HTTP/2 attempt")
				}
			}
			release()
			if !returnedBeforeRelease {
				err = <-done
			}
			<-body.closeFinished

			if err == nil || !strings.Contains(err.Error(), "exceeded configured limit of 8 bytes") {
				t.Fatalf("Get() error = %v, want response body limit error", err)
			}
			if !returnedBeforeRelease {
				t.Fatal("HTTP/2 response overflow waited for a stalled body Close")
			}
			if callerCtx.Err() != nil {
				t.Fatalf("caller context ended before overflow returned: %v", callerCtx.Err())
			}
		})
	}
}

func TestExecuteExpiredContextDoesNotWaitForStalledHTTP2Close(t *testing.T) {
	tests := []struct {
		name           string
		requestTimeout bool
		duringRead     bool
		wantErr        error
	}{
		{name: "caller cancellation after response", wantErr: context.Canceled},
		{name: "caller cancellation during read", duringRead: true, wantErr: context.Canceled},
		{name: "request timeout after response", requestTimeout: true, wantErr: context.DeadlineExceeded},
		{name: "request timeout during read", requestTimeout: true, duringRead: true, wantErr: context.DeadlineExceeded},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			releaseClose := make(chan struct{})
			var releaseOnce sync.Once
			release := func() { releaseOnce.Do(func() { close(releaseClose) }) }

			server := newStalledHTTP2Server(t, http.StatusOK, releaseClose)
			defer func() {
				release()
				server.Close()
			}()

			callerCtx, cancelCaller := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancelCaller()
			bodyReady := make(chan *stalledHTTP2ResponseBody, 1)
			opts := []option.RequestOption{
				option.WithAPIKey("test-key"),
				option.WithBaseURL(server.URL + "/"),
				option.WithMaxRetries(0),
				option.WithHTTPClient(responseDoerFunc(func(req *http.Request) (*http.Response, error) {
					res, err := server.Client().Transport.RoundTrip(req)
					if err != nil {
						return nil, err
					}
					body, err := wrapStalledHTTP2ResponseBody(req, res, releaseClose)
					if err != nil {
						return nil, err
					}
					bodyReady <- body
					endAttempt := cancelCaller
					if test.requestTimeout {
						endAttempt = func() { <-req.Context().Done() }
					}
					if test.duringRead {
						res.Body = &canceledHTTP2ResponseBody{stalledHTTP2ResponseBody: body, endAttempt: endAttempt}
					} else {
						endAttempt()
					}
					return res, nil
				})),
			}
			if test.requestTimeout {
				opts = append(opts, option.WithRequestTimeout(500*time.Millisecond))
			}
			client := openai.NewClient(opts...)

			var response map[string]any
			done := make(chan error, 1)
			go func() {
				done <- client.Get(callerCtx, "canceled", nil, &response)
			}()

			var body *stalledHTTP2ResponseBody
			select {
			case body = <-bodyReady:
			case <-time.After(5 * time.Second):
				t.Fatal("HTTP/2 response was not received")
			}
			select {
			case <-body.closeStarted:
			case <-time.After(5 * time.Second):
				t.Fatal("expired request cleanup was not started")
			}

			var err error
			returnedBeforeRelease := false
			select {
			case err = <-done:
				returnedBeforeRelease = true
			case <-time.After(250 * time.Millisecond):
			}
			release()
			if !returnedBeforeRelease {
				err = <-done
			}
			<-body.closeFinished

			if !errors.Is(err, test.wantErr) {
				t.Fatalf("Get() error = %v, want %v", err, test.wantErr)
			}
			if !returnedBeforeRelease {
				t.Fatal("expired HTTP/2 request waited for a stalled body Close")
			}
			if test.requestTimeout && callerCtx.Err() != nil {
				t.Fatalf("caller safety deadline expired before request timeout returned: %v", callerCtx.Err())
			}
		})
	}
}

func TestExecuteUnownedResponseDoesNotWaitForStalledHTTP2Close(t *testing.T) {
	releaseClose := make(chan struct{})
	var releaseOnce sync.Once
	release := func() { releaseOnce.Do(func() { close(releaseClose) }) }
	server := newStalledHTTP2Server(t, http.StatusOK, releaseClose)
	defer func() {
		release()
		server.Close()
	}()

	bodyReady := make(chan *stalledHTTP2ResponseBody, 1)
	client := openai.NewClient(
		option.WithAPIKey("test-key"),
		option.WithBaseURL(server.URL+"/"),
		option.WithMaxRetries(0),
		option.WithHTTPClient(responseDoerFunc(func(req *http.Request) (*http.Response, error) {
			res, err := server.Client().Transport.RoundTrip(req)
			if err != nil {
				return nil, err
			}
			body, err := wrapStalledHTTP2ResponseBody(req, res, releaseClose)
			if err != nil {
				return nil, err
			}
			bodyReady <- body
			return res, nil
		})),
	)

	done := make(chan error, 1)
	go func() { done <- client.Get(context.Background(), "discard", nil, nil) }()

	var body *stalledHTTP2ResponseBody
	select {
	case body = <-bodyReady:
	case <-time.After(5 * time.Second):
		t.Fatal("HTTP/2 response was not received")
	}
	select {
	case <-body.closeStarted:
	case <-time.After(5 * time.Second):
		t.Fatal("unowned response cleanup was not started")
	}

	var err error
	returnedBeforeRelease := false
	select {
	case err = <-done:
		returnedBeforeRelease = true
	case <-time.After(250 * time.Millisecond):
	}
	if returnedBeforeRelease && body.attemptCtx.Err() == nil {
		t.Fatal("unowned response returned before canceling its attempt context")
	}
	release()
	if !returnedBeforeRelease {
		err = <-done
	}
	<-body.closeFinished

	if err != nil {
		t.Fatalf("Get() error = %v, want nil", err)
	}
	if !returnedBeforeRelease {
		t.Fatal("unowned HTTP/2 response waited for a stalled body Close")
	}
}

func TestExecuteRetryDoesNotWaitForStalledHTTP2Close(t *testing.T) {
	tests := []struct {
		name           string
		requestTimeout bool
	}{
		{name: "caller deadline"},
		{name: "request attempt deadline", requestTimeout: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			releaseClose := make(chan struct{})
			var releaseOnce sync.Once
			release := func() { releaseOnce.Do(func() { close(releaseClose) }) }
			server := newStalledHTTP2Server(t, http.StatusTooManyRequests, releaseClose)
			defer func() {
				release()
				server.Close()
			}()

			deadline := 500 * time.Millisecond
			callerTimeout := deadline
			if test.requestTimeout {
				callerTimeout = 5 * time.Second
			}
			callerCtx, cancelCaller := context.WithTimeout(context.Background(), callerTimeout)
			defer cancelCaller()
			bodyReady := make(chan *stalledHTTP2ResponseBody, 1)
			opts := []option.RequestOption{
				option.WithAPIKey("test-key"),
				option.WithBaseURL(server.URL + "/"),
				option.WithMaxRetries(1),
				option.WithMaxRetryDelay(2 * time.Second),
				option.WithHTTPClient(responseDoerFunc(func(req *http.Request) (*http.Response, error) {
					res, err := server.Client().Transport.RoundTrip(req)
					if err != nil {
						return nil, err
					}
					body, err := wrapStalledHTTP2ResponseBody(req, res, releaseClose)
					if err != nil {
						return nil, err
					}
					bodyReady <- body
					return res, nil
				})),
			}
			if test.requestTimeout {
				opts = append(opts, option.WithRequestTimeout(deadline))
			}
			client := openai.NewClient(opts...)
			done := make(chan error, 1)
			go func() {
				var response map[string]any
				done <- client.Get(callerCtx, "retry", nil, &response)
			}()

			var body *stalledHTTP2ResponseBody
			select {
			case body = <-bodyReady:
			case <-time.After(5 * time.Second):
				t.Fatal("retryable HTTP/2 response was not received")
			}
			select {
			case <-body.closeStarted:
			case <-time.After(5 * time.Second):
				t.Fatal("retry response cleanup was not started")
			}

			var err error
			returnedBeforeRelease := false
			select {
			case err = <-done:
				returnedBeforeRelease = true
			case <-time.After(2 * deadline):
			}
			release()
			if !returnedBeforeRelease {
				err = <-done
			}
			<-body.closeFinished
			if !errors.Is(err, context.DeadlineExceeded) {
				t.Fatalf("Get() error = %v, want deadline exceeded", err)
			}
			if !returnedBeforeRelease {
				t.Fatal("retry waited for a stalled HTTP/2 body Close")
			}
			if test.requestTimeout && callerCtx.Err() != nil {
				t.Fatalf("retry wait ignored attempt deadline: %v", callerCtx.Err())
			}
		})
	}
}

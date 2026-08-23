package openai_test

import (
	"context"
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

func (b *stalledHTTP2ResponseBody) Close() error {
	close(b.closeStarted)
	<-b.releaseClose
	err := b.ReadCloser.Close()
	close(b.closeFinished)
	return err
}

func TestExecuteOverflowDoesNotWaitForStalledHTTP2Close(t *testing.T) {
	for _, status := range []int{http.StatusOK, http.StatusBadRequest} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			releaseClose := make(chan struct{})
			var releaseOnce sync.Once
			release := func() { releaseOnce.Do(func() { close(releaseClose) }) }

			server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				if req.ProtoMajor != 2 {
					t.Errorf("request protocol = %s, want HTTP/2", req.Proto)
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(status)
				if _, err := io.WriteString(w, `{"ok":true}`+strings.Repeat(" ", 256)); err != nil {
					t.Errorf("write response: %v", err)
					return
				}
				w.(http.Flusher).Flush()
				<-releaseClose
			}))
			server.EnableHTTP2 = true
			server.StartTLS()
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
						if res.ProtoMajor != 2 {
							_ = res.Body.Close()
							return nil, fmt.Errorf("response protocol = %s, want HTTP/2", res.Proto)
						}
						body := &stalledHTTP2ResponseBody{
							ReadCloser:    res.Body,
							attemptCtx:    req.Context(),
							closeStarted:  make(chan struct{}),
							releaseClose:  releaseClose,
							closeFinished: make(chan struct{}),
						}
						res.Body = body
						bodyReady <- body
						return res, nil
					}),
				}),
			)

			callerCtx, cancelCaller := context.WithTimeout(context.Background(), time.Second)
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

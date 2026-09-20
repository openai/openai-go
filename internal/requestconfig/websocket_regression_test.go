package requestconfig_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	wire "github.com/coder/websocket"
	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
	transport "github.com/openai/openai-go/v3/packages/websocket"
)

func TestWebSocketQueryOperationsIncludeBaseQuery(t *testing.T) {
	for _, tc := range []struct {
		name string
		opts []option.RequestOption
		want []string
	}{
		{"add", []option.RequestOption{option.WithQueryAdd("tag", "extra")}, []string{"base", "extra"}},
		{"delete", []option.RequestOption{option.WithQueryDel("tag")}, nil},
		{"set", []option.RequestOption{option.WithQuery("tag", "new")}, []string{"new"}},
		{"delete then add", []option.RequestOption{option.WithQueryDel("tag"), option.WithQueryAdd("tag", "new")}, []string{"new"}},
		{"add then delete", []option.RequestOption{option.WithQueryAdd("tag", "extra"), option.WithQueryDel("tag")}, nil},
		{"add then set", []option.RequestOption{option.WithQueryAdd("tag", "extra"), option.WithQuery("tag", "new")}, []string{"new"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := make(chan url.Values, 1)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				got <- r.URL.Query()
				socket, err := wire.Accept(w, r, nil)
				if err == nil {
					defer func() { _ = socket.CloseNow() }()
					_, _, _ = socket.Read(r.Context())
				}
			}))
			defer server.Close()
			opts := append([]option.RequestOption{option.WithBaseURL(server.URL + "/v1?tag=base&keep=yes")}, tc.opts...)
			connection, err := dialPrepared(t, opts...)
			if err != nil {
				t.Fatal(err)
			}
			defer connection.Abort()
			query := <-got
			if !reflect.DeepEqual(query["tag"], tc.want) || query.Get("keep") != "yes" {
				t.Fatalf("handshake query = %v, want tag=%v and keep=yes", query, tc.want)
			}
		})
	}
}

func TestWebSocketRedirectsUseOneLogicalAttempt(t *testing.T) {
	for _, retries := range []int{0, 2} {
		t.Run(strconv.Itoa(retries), func(t *testing.T) {
			var middlewareCalls, requests atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				n := requests.Add(1)
				if r.Header.Get("X-Handshake") != "preserved" {
					t.Error("lost middleware header")
				}
				if n < 4 {
					http.Redirect(w, r, "/next", http.StatusTemporaryRedirect)
					return
				}
				socket, err := wire.Accept(w, r, nil)
				if err == nil {
					defer func() { _ = socket.CloseNow() }()
					_, _, _ = socket.Read(r.Context())
				}
			}))
			defer server.Close()
			connection, err := dialPrepared(t, option.WithBaseURL(server.URL), option.WithMaxRetries(retries),
				requestconfig.RequestOptionFunc(func(cfg *requestconfig.RequestConfig) error {
					cfg.InstallRequestRetryScope(true)
					cfg.InstallRequestAttemptMiddleware()
					return nil
				}), option.WithMiddleware(func(r *http.Request, next option.MiddlewareNext) (*http.Response, error) {
					middlewareCalls.Add(1)
					r.Header.Set("X-Handshake", "preserved")
					return next(r)
				}))
			if err != nil {
				t.Fatal(err)
			}
			defer connection.Abort()
			if middlewareCalls.Load() != 1 || requests.Load() != 4 {
				t.Fatalf("middleware calls=%d, HTTP requests=%d", middlewareCalls.Load(), requests.Load())
			}
		})
	}
}

type legacyUpgradeTransport struct {
	canceled chan struct{}
	once     sync.Once
	calls    atomic.Int32
	requests atomic.Int32
	redirect bool
}

func (tr *legacyUpgradeTransport) RoundTrip(*http.Request) (*http.Response, error) {
	if tr.redirect && tr.requests.Add(1) == 1 {
		return &http.Response{StatusCode: http.StatusTemporaryRedirect,
			Header: http.Header{"Location": []string{"/next"}},
			Body:   &legacyUpgradeBody{canceled: tr.canceled}, ContentLength: -1}, nil
	}
	<-tr.canceled // A legacy transport does not observe Request.Context.
	return nil, context.DeadlineExceeded
}

type legacyUpgradeBody struct{ canceled <-chan struct{} }

func (b *legacyUpgradeBody) Read([]byte) (int, error) {
	<-b.canceled
	return 0, context.DeadlineExceeded
}

func (b *legacyUpgradeBody) Close() error { return nil }

func (tr *legacyUpgradeTransport) CancelRequest(*http.Request) {
	tr.calls.Add(1)
	tr.once.Do(func() { close(tr.canceled) })
}

func TestWebSocketTimeoutCancelsLegacyTransport(t *testing.T) {
	for _, tc := range []struct {
		name               string
		prepared, redirect bool
	}{
		{name: "direct"}, {name: "prepared", prepared: true},
		{name: "direct redirect body", redirect: true},
		{name: "prepared redirect body", prepared: true, redirect: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tr := &legacyUpgradeTransport{canceled: make(chan struct{}), redirect: tc.redirect}
			defer tr.once.Do(func() { close(tr.canceled) })
			client := &http.Client{Transport: tr, Timeout: 50 * time.Millisecond}
			done := make(chan error, 1)
			go func() {
				request, _ := http.NewRequest(http.MethodGet, "http://example.invalid/responses", nil)
				if tc.prepared {
					cfg, err := requestconfig.NewRequestConfig(context.Background(), http.MethodGet, "responses", nil, nil, option.WithBaseURL("http://example.invalid/"), option.WithHTTPClient(client))
					if err != nil {
						done <- err
						return
					}
					request, client, err = cfg.PrepareWebSocket()
					if err != nil {
						done <- err
						return
					}
				}
				_, err := transport.Dial(context.Background(), request, client, transport.Options{}, func([]byte) (string, error) { return "", nil }, func(string) string { return "" })
				done <- err
			}()
			select {
			case err := <-done:
				if !errors.Is(err, context.DeadlineExceeded) || tr.calls.Load() != 1 {
					t.Fatalf("timeout = %v, legacy cancels = %d", err, tr.calls.Load())
				}
			case <-time.After(2 * time.Second):
				tr.once.Do(func() { close(tr.canceled) })
				<-done
				t.Fatal("handshake timeout did not cancel legacy transport")
			}
		})
	}
}

type failingUpgradeTransport struct{ cause error }

func (tr failingUpgradeTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, tr.cause
}

func TestWebSocketDialErrorDoesNotExposeQuery(t *testing.T) {
	want := errors.New("transport failure")
	_, err := dialPrepared(t, option.WithBaseURL("https://example.invalid/v1?api_key=fake-query-secret"), option.WithHTTPClient(&http.Client{Transport: failingUpgradeTransport{cause: want}}))
	if err == nil || strings.Contains(err.Error(), "fake-query-secret") {
		t.Fatalf("unsafe error: %v", err)
	}
	var urlError *url.Error
	if !errors.Is(err, want) || !errors.As(err, &urlError) || !strings.Contains(err.Error(), want.Error()) {
		t.Fatalf("transport error lost its cause: %v", err)
	}
}

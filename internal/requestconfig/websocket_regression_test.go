package requestconfig_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"

	wire "github.com/coder/websocket"
	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
	transport "github.com/openai/openai-go/v3/packages/websocket"
)

func TestWebSocketRetryAfterExceedingBudgetPreservesResponse(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		attempts := 0
		cfg, err := requestconfig.NewRequestConfig(context.Background(), http.MethodGet, "responses", nil, nil,
			option.WithBaseURL("https://example.invalid/v1"),
			option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
				attempts++
				return &http.Response{
					StatusCode: http.StatusUnauthorized,
					Header:     http.Header{"Retry-After": {"121"}, "X-Should-Retry": {"true"}},
					Body:       io.NopCloser(strings.NewReader("retry later")),
					Request:    req,
				}, nil
			}),
		)
		if err != nil {
			t.Fatal(err)
		}
		cfg.InstallRequestRetryScope(true)
		req, client, err := cfg.PrepareWebSocket()
		if err != nil {
			t.Fatal(err)
		}
		start := time.Now()
		response, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer func() { _ = response.Body.Close() }()
		body, err := io.ReadAll(response.Body)
		if err != nil || string(body) != "retry later" || response.StatusCode != http.StatusUnauthorized {
			t.Fatalf("response = %d %q, error = %v; want original handshake rejection", response.StatusCode, body, err)
		}
		if attempts != 1 || time.Since(start) != 0 {
			t.Fatalf("attempts = %d, elapsed = %s; want one attempt without waiting", attempts, time.Since(start))
		}
	})
}

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
		requestTimeout     bool
	}{
		{name: "direct"}, {name: "prepared", prepared: true},
		{name: "direct redirect body", redirect: true},
		{name: "prepared redirect body", prepared: true, redirect: true},
		{name: "request timeout", prepared: true, requestTimeout: true},
		{name: "request timeout redirect body", prepared: true, redirect: true, requestTimeout: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tr := &legacyUpgradeTransport{canceled: make(chan struct{}), redirect: tc.redirect}
			defer tr.once.Do(func() { close(tr.canceled) })
			client := &http.Client{Transport: tr, Timeout: 50 * time.Millisecond}
			var opts []option.RequestOption
			if tc.requestTimeout {
				client.Timeout = 0
				opts = append(opts, option.WithRequestTimeout(50*time.Millisecond))
			}
			done := make(chan error, 1)
			go func() {
				request, _ := http.NewRequest(http.MethodGet, "http://example.invalid/responses", nil)
				if tc.prepared {
					cfg, err := requestconfig.NewRequestConfig(context.Background(), http.MethodGet, "responses", nil, nil, append(opts, option.WithBaseURL("http://example.invalid/"), option.WithHTTPClient(client))...)
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

func TestPreparedWebSocketAttemptTimeoutDoesNotCloseLiveSocket(t *testing.T) {
	const timeout = 100 * time.Millisecond
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		socket, err := wire.Accept(w, r, nil)
		if err != nil {
			t.Errorf("upgrade: %v", err)
			return
		}
		defer func() { _ = socket.CloseNow() }()
		timer := time.NewTimer(2 * timeout)
		defer timer.Stop()
		select {
		case <-r.Context().Done():
			return
		case <-timer.C:
		}
		if err := socket.Write(r.Context(), wire.MessageText, []byte("after opening deadline")); err != nil {
			t.Errorf("write after opening deadline: %v", err)
		}
		_, _, _ = socket.Read(r.Context())
	}))
	defer server.Close()
	connection, err := dialPrepared(t, option.WithBaseURL(server.URL), option.WithRequestTimeout(timeout))
	if err != nil {
		t.Fatal(err)
	}
	defer connection.Abort()
	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
	defer cancel()
	got, err := connection.Recv(ctx)
	if err != nil || got != "after opening deadline" {
		t.Fatalf("receive after opening timeout = %q, %v", got, err)
	}
}

type openingResponseBody struct {
	io.Reader
	ctx    context.Context
	closed bool
	err    error
}

func (b *openingResponseBody) Close() error {
	b.closed = true
	b.err = b.ctx.Err()
	return nil
}

func TestPreparedWebSocketRejectionPreservesBodyOwnership(t *testing.T) {
	for _, status := range []int{http.StatusUnauthorized, http.StatusSwitchingProtocols} {
		t.Run(strconv.Itoa(status), func(t *testing.T) {
			var body *openingResponseBody
			cfg, err := requestconfig.NewRequestConfig(t.Context(), http.MethodGet, "responses", nil, nil,
				option.WithBaseURL("https://example.invalid/v1"), option.WithRequestTimeout(time.Hour),
				option.WithMiddleware(func(req *http.Request, _ option.MiddlewareNext) (*http.Response, error) {
					body = &openingResponseBody{Reader: strings.NewReader("server rejection"), ctx: req.Context()}
					return &http.Response{StatusCode: status, Header: make(http.Header), Body: body, Request: req}, nil
				}),
			)
			if err != nil {
				t.Fatal(err)
			}
			req, client, err := cfg.PrepareWebSocket()
			if err != nil {
				t.Fatal(err)
			}
			resp, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			if _, writable := resp.Body.(io.ReadWriteCloser); writable {
				t.Fatal("timer wrapper added unsupported write capability")
			}
			if _, ok := body.ctx.Deadline(); !ok {
				t.Error("opening attempt deadline not visible to middleware")
			}
			data, err := io.ReadAll(resp.Body)
			if err != nil || string(data) != "server rejection" || body.closed {
				t.Fatalf("response = %q, %v, closed=%t", data, err, body.closed)
			}
			if err := resp.Body.Close(); err != nil {
				t.Fatal(err)
			}
			if !body.closed || body.err != nil || !errors.Is(body.ctx.Err(), context.Canceled) {
				t.Fatalf("close = %t, context during/after close = %v / %v", body.closed, body.err, body.ctx.Err())
			}
		})
	}
}

func TestPreparedWebSocketMiddlewareCannotOutliveAttempt(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var body *openingResponseBody
		var captured *http.Response
		client := &http.Client{Transport: responseTransport(func(req *http.Request) (*http.Response, error) {
			body = &openingResponseBody{Reader: strings.NewReader("server rejection"), ctx: req.Context()}
			return &http.Response{StatusCode: http.StatusForbidden, Header: make(http.Header), Body: body, Request: req}, nil
		})}
		cfg, err := requestconfig.NewRequestConfig(t.Context(), http.MethodGet, "responses", nil, nil,
			option.WithBaseURL("https://example.invalid/v1"), option.WithRequestTimeout(time.Second),
			option.WithHTTPClient(client), option.WithResponseInto(&captured),
			option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
				resp, err := next(req)
				<-req.Context().Done()
				return resp, err
			}),
		)
		if err != nil {
			t.Fatal(err)
		}
		req, client, err := cfg.PrepareWebSocket()
		if err != nil {
			t.Fatal(err)
		}
		conn, err := transport.Dial(t.Context(), req, client, transport.Options{}, func([]byte) (string, error) {
			return "", nil
		}, func(string) string { return "" })
		if conn != nil {
			defer conn.Abort()
		}
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Errorf("opening after middleware deadline = %v, want context deadline", err)
		}
		if captured == nil || captured.StatusCode != http.StatusForbidden {
			t.Errorf("captured response = %v, want original rejection", captured)
		}
		if body == nil || !body.closed {
			t.Error("rejected response body was not closed")
		}
	})
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

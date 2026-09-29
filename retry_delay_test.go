package openai_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"testing/synctest"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/shared"
)

type retryDelayRoundTripperFunc func(*http.Request) (*http.Response, error)

func (f retryDelayRoundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

type retryDelayResponseBody struct {
	io.Reader
	closes int
}

func (b *retryDelayResponseBody) Read(p []byte) (int, error) {
	if b.closes != 0 {
		return 0, io.ErrClosedPipe
	}
	return b.Reader.Read(p)
}

func (b *retryDelayResponseBody) Close() error {
	b.closes++
	return nil
}

func TestRetryAfterExceedingConfiguredMaximumDoesNotRetry(t *testing.T) {
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
	if attempts != 1 {
		t.Fatalf("attempts = %d, want 1", attempts)
	}
}

func TestRetryAfterServerMinimum(t *testing.T) {
	for _, tc := range []struct {
		name          string
		header        http.Header
		maxRetryDelay time.Duration
		wantDelay     time.Duration
		wantRetry     bool
	}{
		{name: "seconds", header: http.Header{"Retry-After": {"30"}}, wantDelay: 30 * time.Second, wantRetry: true},
		{name: "milliseconds take priority", header: http.Header{"Retry-After-Ms": {"30000"}, "Retry-After": {"121"}}, wantDelay: 30 * time.Second, wantRetry: true},
		{name: "default boundary", header: http.Header{"Retry-After": {"120"}}, wantDelay: 2 * time.Minute, wantRetry: true},
		{name: "over default", header: http.Header{"Retry-After": {"120.001"}}},
		{name: "over default milliseconds", header: http.Header{"Retry-After-Ms": {"120001"}}},
		{name: "explicit retry still bounded", header: http.Header{"Retry-After": {"121"}, "X-Should-Retry": {"true"}}},
		{name: "finite duration overflow", header: http.Header{"Retry-After": {"1e100"}}},
		{name: "numeric overflow", header: http.Header{"Retry-After": {"1e1000"}}},
		{name: "millisecond overflow does not fall back", header: http.Header{"Retry-After-Ms": {"1e1000"}, "Retry-After": {"0"}}},
		{name: "configured boundary", header: http.Header{"Retry-After": {"0.01"}}, maxRetryDelay: 10 * time.Millisecond, wantDelay: 10 * time.Millisecond, wantRetry: true},
		{name: "over configured", header: http.Header{"Retry-After": {"0.011"}}, maxRetryDelay: 10 * time.Millisecond},
		{name: "configured above default", header: http.Header{"Retry-After": {"180"}}, maxRetryDelay: 3 * time.Minute, wantDelay: 3 * time.Minute, wantRetry: true},
		{name: "explicit no retry", header: http.Header{"Retry-After": {"30"}, "X-Should-Retry": {"false"}}},
		{name: "HTTP date", wantDelay: 30 * time.Second, wantRetry: true},
		{name: "over default HTTP date"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				header := tc.header.Clone()
				if header == nil {
					delay := 30 * time.Second
					if !tc.wantRetry {
						delay = 121 * time.Second
					}
					header = http.Header{"Retry-After": {time.Now().Add(delay).UTC().Format(http.TimeFormat)}}
				}
				header.Set("Content-Type", "application/json")
				attempts := 0
				var bodies []*retryDelayResponseBody
				opts := []option.RequestOption{
					option.WithAPIKey("fake-key"),
					option.WithBaseURL("https://example.invalid/v1"),
					option.WithMaxRetries(1),
					option.WithHTTPClient(&http.Client{Transport: retryDelayRoundTripperFunc(func(req *http.Request) (*http.Response, error) {
						attempts++
						status, body := http.StatusServiceUnavailable, `{"error":{"message":"overloaded","type":"server_error"}}`
						if attempts > 1 {
							status, body = http.StatusOK, `{}`
						}
						responseBody := &retryDelayResponseBody{Reader: strings.NewReader(body)}
						bodies = append(bodies, responseBody)
						return &http.Response{StatusCode: status, Header: header, Body: responseBody, Request: req}, nil
					})}),
				}
				if tc.maxRetryDelay > 0 {
					opts = append(opts, option.WithMaxRetryDelay(tc.maxRetryDelay))
				}
				client := openai.NewClient(opts...)
				var response map[string]any
				start := time.Now()
				err := client.Get(context.Background(), "models", nil, &response)
				if elapsed := time.Since(start); elapsed != tc.wantDelay {
					t.Errorf("elapsed = %s, want %s", elapsed, tc.wantDelay)
				}
				wantAttempts := 1
				if tc.wantRetry {
					wantAttempts = 2
					if err != nil {
						t.Errorf("retry failed: %v", err)
					}
				} else {
					var apiErr *openai.Error
					if !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusServiceUnavailable || apiErr.Message != "overloaded" {
						t.Fatalf("error = %v, want original overload error", err)
					}
					if apiErr.Response.Header.Get("Retry-After") != header.Get("Retry-After") || apiErr.Response.Header.Get("Retry-After-Ms") != header.Get("Retry-After-Ms") {
						t.Error("original retry headers were not preserved")
					}
				}
				if attempts != wantAttempts {
					t.Errorf("attempts = %d, want %d", attempts, wantAttempts)
				}
				for i, body := range bodies {
					if body.closes != 1 {
						t.Errorf("response body %d closed %d times, want 1", i, body.closes)
					}
				}
			})
		})
	}
}

func TestRetryAfterServerDelayCancellation(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		attempts := 0
		body := &retryDelayResponseBody{Reader: strings.NewReader(`{"error":{"message":"overloaded"}}`)}
		client := openai.NewClient(
			option.WithAPIKey("fake-key"),
			option.WithBaseURL("https://example.invalid/v1"),
			option.WithHTTPClient(&http.Client{Transport: retryDelayRoundTripperFunc(func(req *http.Request) (*http.Response, error) {
				attempts++
				return &http.Response{StatusCode: http.StatusServiceUnavailable, Header: http.Header{"Retry-After": {"30"}}, Body: body, Request: req}, nil
			})}),
		)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		start := time.Now()
		var response map[string]any
		err := client.Get(ctx, "models", nil, &response)
		if !errors.Is(err, context.DeadlineExceeded) || attempts != 1 || time.Since(start) != 5*time.Second {
			t.Fatalf("error = %v, attempts = %d, elapsed = %s; want deadline exceeded after 5s with one attempt", err, attempts, time.Since(start))
		}
		if body.closes != 1 {
			t.Fatalf("response body closed %d times, want 1", body.closes)
		}
	})
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

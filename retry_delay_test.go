package openai_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"testing/synctest"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type retryDelayRoundTripperFunc func(*http.Request) (*http.Response, error)

func (f retryDelayRoundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestRetryAfterExceedingMaximumPreservesAPIError(t *testing.T) {
	for _, test := range []struct {
		status   int
		typeName string
		code     string
	}{
		{http.StatusTooManyRequests, "rate_limit_error", "slow_down"},
		{http.StatusServiceUnavailable, "service_unavailable_error", "server_is_overloaded"},
	} {
		t.Run(http.StatusText(test.status), func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				attempts := 0
				body := fmt.Sprintf(`{"error":{"type":%q,"code":%q,"message":"Please try again later.","param":null}}`, test.typeName, test.code)
				headers := http.Header{"Retry-After": {"90"}, "X-Request-Id": {"request-test"}}
				client := openai.NewClient(
					option.WithAPIKey("test-key"),
					option.WithMaxRetries(1),
					option.WithHTTPClient(&http.Client{Transport: retryDelayRoundTripperFunc(func(req *http.Request) (*http.Response, error) {
						attempts++
						return &http.Response{StatusCode: test.status, Header: headers.Clone(), Body: io.NopCloser(strings.NewReader(body)), Request: req}, nil
					})}),
				)
				var response *http.Response
				started := time.Now()
				err := client.Get(t.Context(), "/models/test", nil, nil, option.WithResponseInto(&response))
				var apiErr *openai.Error
				if !errors.As(err, &apiErr) {
					t.Fatalf("Get() error = %v, want *openai.Error", err)
				}
				if attempts != 1 || time.Since(started) != 0 {
					t.Errorf("Get(Retry-After: 90) attempts = %d, elapsed = %s; want 1, 0s", attempts, time.Since(started))
				}
				if apiErr.StatusCode != test.status || apiErr.Type != test.typeName || apiErr.Code != test.code {
					t.Errorf("Get() error fields = (%d, %q, %q), want (%d, %q, %q)", apiErr.StatusCode, apiErr.Type, apiErr.Code, test.status, test.typeName, test.code)
				}
				if response != apiErr.Response || response.Header.Get("Retry-After") != "90" || response.Header.Get("X-Request-Id") != "request-test" {
					t.Errorf("Get() did not preserve the response and retry/request-ID headers")
				}
				defer func() { _ = response.Body.Close() }()
				gotBody, readErr := io.ReadAll(response.Body)
				if readErr != nil || string(gotBody) != body {
					t.Errorf("Get() response body = %q, error = %v; want %q, nil", gotBody, readErr, body)
				}
			})
		})
	}
}

func TestRetryAfterDelayEdgeCases(t *testing.T) {
	tests := map[string]struct {
		header         http.Header
		maxRetryDelay  time.Duration
		wantAttempts   int
		minimumWait    time.Duration
		maximumWait    time.Duration
		httpDateOffset time.Duration
	}{
		"allowed HTTP date":                         {httpDateOffset: 3 * time.Second, wantAttempts: 2, minimumWait: 3 * time.Second, maximumWait: 3 * time.Second},
		"excessive HTTP date":                       {httpDateOffset: 90 * time.Second, wantAttempts: 1},
		"explicit retry cannot shorten the minimum": {header: http.Header{"Retry-After": {"90"}, "X-Should-Retry": {"true"}}, wantAttempts: 1},
		"explicit retry opt-out":                    {header: http.Header{"Retry-After": {"3"}, "X-Should-Retry": {"false"}}, wantAttempts: 1},
		"zero seconds":                              {header: http.Header{"Retry-After": {"0"}}, wantAttempts: 2},
		"zero milliseconds":                         {header: http.Header{"Retry-After-Ms": {"0"}}, wantAttempts: 2},
		"elapsed HTTP date":                         {httpDateOffset: -time.Minute, wantAttempts: 2},
		"finite scaling overflow":                   {header: http.Header{"Retry-After": {"2" + strings.Repeat("0", 299)}}, wantAttempts: 1},
		"above configured maximum":                  {header: http.Header{"Retry-After-Ms": {"11"}}, maxRetryDelay: 10 * time.Millisecond, wantAttempts: 1},
		"exact configured maximum":                  {header: http.Header{"Retry-After-Ms": {"10"}}, maxRetryDelay: 10 * time.Millisecond, wantAttempts: 2, minimumWait: 10 * time.Millisecond, maximumWait: 10 * time.Millisecond},
		"exact default maximum":                     {header: http.Header{"Retry-After": {"8"}}, wantAttempts: 2, minimumWait: 8 * time.Second, maximumWait: 8 * time.Second},
		"long delay allowed by configuration":       {header: http.Header{"Retry-After": {"90"}}, maxRetryDelay: 90 * time.Second, wantAttempts: 2, minimumWait: 90 * time.Second, maximumWait: 90 * time.Second},
		"milliseconds take precedence":              {header: http.Header{"Retry-After-Ms": {"1"}, "Retry-After": {"90"}}, wantAttempts: 2, minimumWait: time.Millisecond, maximumWait: time.Millisecond},
		"excessive milliseconds take precedence":    {header: http.Header{"Retry-After-Ms": {"9000"}, "Retry-After": {"0"}}, wantAttempts: 1},
		"invalid milliseconds fall back to seconds": {header: http.Header{"Retry-After-Ms": {"invalid"}, "Retry-After": {"0.25"}}, wantAttempts: 2, minimumWait: 250 * time.Millisecond, maximumWait: 250 * time.Millisecond},
		"invalid header uses bounded backoff":       {header: http.Header{"Retry-After": {"NaN"}}, maxRetryDelay: 10 * time.Millisecond, wantAttempts: 2, minimumWait: 7500 * time.Microsecond, maximumWait: 10 * time.Millisecond},
		"missing header uses bounded backoff":       {maxRetryDelay: 10 * time.Millisecond, wantAttempts: 2, minimumWait: 7500 * time.Microsecond, maximumWait: 10 * time.Millisecond},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				if test.httpDateOffset != 0 {
					test.header = http.Header{"Retry-After": {time.Now().Add(test.httpDateOffset).UTC().Format(time.RFC1123)}}
				}
				attempts := 0
				opts := []option.RequestOption{
					option.WithAPIKey("test-key"),
					option.WithMaxRetries(1),
					option.WithHTTPClient(&http.Client{Transport: retryDelayRoundTripperFunc(func(req *http.Request) (*http.Response, error) {
						attempts++
						return &http.Response{StatusCode: http.StatusTooManyRequests, Header: test.header, Body: http.NoBody, Request: req}, nil
					})}),
				}
				if test.maxRetryDelay > 0 {
					opts = append(opts, option.WithMaxRetryDelay(test.maxRetryDelay))
				}
				client := openai.NewClient(opts...)
				started := time.Now()
				err := client.Get(t.Context(), "/models/test", nil, nil)
				var apiErr *openai.Error
				if !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusTooManyRequests {
					t.Fatalf("Get(%v) error = %v, want rate-limit API error", test.header, err)
				}
				if attempts != test.wantAttempts {
					t.Errorf("Get(%v) attempts = %d, want %d", test.header, attempts, test.wantAttempts)
				}
				if elapsed := time.Since(started); elapsed < test.minimumWait || elapsed > test.maximumWait {
					t.Errorf("Get(%v) elapsed = %s, want [%s, %s]", test.header, elapsed, test.minimumWait, test.maximumWait)
				}
			})
		})
	}
}

func TestRetryAfterContextErrors(t *testing.T) {
	for _, test := range []struct {
		name                 string
		hint                 string
		cancelDuringResponse bool
		want                 error
	}{
		{name: "cancelled before excessive delay", hint: "90", cancelDuringResponse: true, want: context.Canceled},
		{name: "deadline during eligible delay", hint: "3", want: context.DeadlineExceeded},
	} {
		t.Run(test.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				ctx, cancel := context.WithTimeout(t.Context(), time.Second)
				defer cancel()
				attempts := 0
				client := openai.NewClient(option.WithAPIKey("test-key"), option.WithHTTPClient(&http.Client{
					Transport: retryDelayRoundTripperFunc(func(req *http.Request) (*http.Response, error) {
						attempts++
						if test.cancelDuringResponse {
							cancel()
						}
						return &http.Response{StatusCode: http.StatusTooManyRequests, Header: http.Header{"Retry-After": {test.hint}}, Body: http.NoBody, Request: req}, nil
					}),
				}))
				err := client.Get(ctx, "/models/test", nil, nil)
				if !errors.Is(err, test.want) || attempts != 1 {
					t.Errorf("Get(Retry-After: %s) error = %v, attempts = %d; want %v, 1", test.hint, err, attempts, test.want)
				}
			})
		})
	}
}

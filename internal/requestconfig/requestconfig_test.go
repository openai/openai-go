package requestconfig

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"
)

type closeTrackingReadCloser struct {
	io.ReadCloser
	closes int
}

func (b *closeTrackingReadCloser) Close() error {
	b.closes++
	return b.ReadCloser.Close()
}

type httpDoerFunc func(*http.Request) (*http.Response, error)

func (f httpDoerFunc) Do(req *http.Request) (*http.Response, error) {
	return f(req)
}

func newTrackedFileBody(t *testing.T) *closeTrackingReadCloser {
	t.Helper()

	f, err := os.CreateTemp(t.TempDir(), "request-body-*")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString("body"); err != nil {
		t.Fatal(err)
	}
	if _, err := f.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	return &closeTrackingReadCloser{ReadCloser: f}
}

func newBodyCloseRequestConfig(t *testing.T, body io.ReadCloser) *RequestConfig {
	t.Helper()

	cfg, err := NewRequestConfig(context.Background(), http.MethodPost, "/models", body, nil)
	if err != nil {
		t.Fatal(err)
	}
	cfg.BaseURL, err = url.Parse("https://example.com/")
	if err != nil {
		t.Fatal(err)
	}
	return cfg
}

func TestFormatPathEscapesPathParams(t *testing.T) {
	tests := map[string]struct {
		format string
		params []string
		want   string
	}{
		"slash": {
			format: "vector_stores/%s",
			params: []string{"../videos/vid_123"},
			want:   "vector_stores/..%2Fvideos%2Fvid_123",
		},
		"query and fragment": {
			format: "vector_stores/%s",
			params: []string{"vs_123/files/file_456?limit=1#frag"},
			want:   "vector_stores/vs_123%2Ffiles%2Ffile_456%3Flimit=1%23frag",
		},
		"encoded dot segments": {
			format: "vector_stores/%s",
			params: []string{"%2e%2e/videos/vid_123"},
			want:   "vector_stores/%252e%252e%2Fvideos%2Fvid_123",
		},
		"bare dot": {
			format: "vector_stores/%s",
			params: []string{"."},
			want:   "vector_stores/%2E",
		},
		"bare dot dot": {
			format: "vector_stores/%s",
			params: []string{".."},
			want:   "vector_stores/%2E%2E",
		},
		"multiple params": {
			format: "organization/projects/%s/api_keys/%s",
			params: []string{"proj_123/../../admin_api_keys/key_456?", "ignored"},
			want:   "organization/projects/proj_123%2F..%2F..%2Fadmin_api_keys%2Fkey_456%3F/api_keys/ignored",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if got := FormatPath(test.format, test.params...); got != test.want {
				t.Fatalf("FormatPath() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestNormalizeOS(t *testing.T) {
	tests := map[string]string{
		"android": "Android",
		"darwin":  "MacOS",
		"freebsd": "FreeBSD",
		"ios":     "iOS",
		"linux":   "Linux",
		"openbsd": "OpenBSD",
		"solaris": "Other:solaris",
		"windows": "Windows",
	}

	for input, want := range tests {
		t.Run(input, func(t *testing.T) {
			if got := normalizeOS(input); got != want {
				t.Fatalf("normalizeOS(%q) = %q, want %q", input, got, want)
			}
		})
	}
}

func TestRequestFinalizerComposesThroughApply(t *testing.T) {
	finalized := false
	wrapped := RequestOptionFunc(func(cfg *RequestConfig) error {
		return WithRequestFinalizer(func(cfg *RequestConfig) error {
			finalized = true
			if got := cfg.Request.Header.Get("X-Late-Option"); got != "present" {
				t.Fatalf("late option header = %q", got)
			}
			return nil
		}).Apply(cfg)
	})
	lateOption := RequestOptionFunc(func(cfg *RequestConfig) error {
		cfg.Request.Header.Set("X-Late-Option", "present")
		return nil
	})

	_, err := NewRequestConfig(
		context.Background(),
		"GET",
		"/models",
		nil,
		nil,
		wrapped,
		lateOption,
	)
	if err != nil {
		t.Fatal(err)
	}
	if !finalized {
		t.Fatal("request finalizer did not run")
	}
}

func TestCloneDoesNotAliasMiddlewareSlice(t *testing.T) {
	cfg, err := NewRequestConfig(context.Background(), http.MethodGet, "/models", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	cfg.Middlewares = []middleware{func(*http.Request, middlewareNext) (*http.Response, error) {
		return nil, nil
	}}

	clone := cfg.Clone(context.Background())
	if clone == nil {
		t.Fatal("Clone() = nil")
	}
	if &clone.Middlewares[0] == &cfg.Middlewares[0] {
		t.Fatal("clone middleware aliases the original configuration")
	}
}

func TestParseRetryAfterHeaderBoundsRemoteDelays(t *testing.T) {
	tests := map[string]struct {
		header http.Header
		want   time.Duration
		ok     bool
	}{
		"milliseconds": {
			header: http.Header{"Retry-After-Ms": {"125"}},
			want:   125 * time.Millisecond,
			ok:     true,
		},
		"fractional seconds": {
			header: http.Header{"Retry-After": {"0.25"}},
			want:   250 * time.Millisecond,
			ok:     true,
		},
		"huge value": {
			header: http.Header{"Retry-After": {"1e100"}},
			want:   DefaultMaxServerDelay,
			ok:     true,
		},
		"finite scaling overflow": {
			header: http.Header{"Retry-After": {"2" + strings.Repeat("0", 299)}},
			want:   DefaultMaxServerDelay,
			ok:     true,
		},
		"far future date": {
			header: http.Header{"Retry-After": {time.Now().Add(time.Hour).UTC().Format(time.RFC1123)}},
			want:   DefaultMaxServerDelay,
			ok:     true,
		},
		"invalid preferred header falls back": {
			header: http.Header{"Retry-After-Ms": {"-1"}, "Retry-After": {"0.5"}},
			want:   500 * time.Millisecond,
			ok:     true,
		},
		"zero": {
			header: http.Header{"Retry-After": {"0"}},
			ok:     true,
		},
		"zero milliseconds": {
			header: http.Header{"Retry-After-Ms": {"0"}},
			ok:     true,
		},
		"negative": {
			header: http.Header{"Retry-After-Ms": {"-100"}},
		},
		"not a number": {
			header: http.Header{"Retry-After": {"NaN"}},
		},
		"infinite": {
			header: http.Header{"Retry-After": {"+Inf"}},
		},
		"past date": {
			header: http.Header{"Retry-After": {time.Now().Add(-time.Hour).UTC().Format(time.RFC1123)}},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got, ok := parseRetryAfterHeader(&http.Response{Header: test.header}, DefaultMaxServerDelay)
			if ok != test.ok || got != test.want {
				t.Fatalf("parseRetryAfterHeader() = (%s, %t), want (%s, %t)", got, ok, test.want, test.ok)
			}
		})
	}
}

func TestWaitForDelayObservesCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := WaitForDelay(ctx, 0); !errors.Is(err, context.Canceled) {
		t.Fatalf("WaitForDelay() error = %v, want %v", err, context.Canceled)
	}
}

func TestExecuteClosesAttemptBodyOnHandlerError(t *testing.T) {
	t.Run("no retry", func(t *testing.T) {
		body := newTrackedFileBody(t)
		cfg := newBodyCloseRequestConfig(t, body)
		cfg.Request.GetBody = func() (io.ReadCloser, error) {
			t.Fatal("GetBody called for no-retry error")
			return nil, nil
		}

		attempts := 0
		cfg.Middlewares = []middleware{func(*http.Request, middlewareNext) (*http.Response, error) {
			attempts++
			return nil, WithNoRetryError(errors.New("blocked"))
		}}

		err := cfg.Execute()
		if err == nil || err.Error() != "blocked" {
			t.Fatalf("Execute() error = %v, want blocked", err)
		}
		if attempts != 1 {
			t.Fatalf("attempts = %d, want 1", attempts)
		}
		if body.closes != 1 {
			t.Fatalf("body closes = %d, want 1", body.closes)
		}
	})

	t.Run("retry", func(t *testing.T) {
		firstBody := newTrackedFileBody(t)
		bodies := []*closeTrackingReadCloser{firstBody}
		cfg := newBodyCloseRequestConfig(t, firstBody)
		cfg.MaxRetries = 1
		cfg.Request.GetBody = func() (io.ReadCloser, error) {
			body := newTrackedFileBody(t)
			bodies = append(bodies, body)
			return body, nil
		}

		attempts := 0
		cfg.Middlewares = []middleware{func(*http.Request, middlewareNext) (*http.Response, error) {
			attempts++
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				Header:     http.Header{"Retry-After-Ms": {"0"}},
				Body:       http.NoBody,
			}, errors.New("transient")
		}}

		err := cfg.Execute()
		if err == nil || err.Error() != "transient" {
			t.Fatalf("Execute() error = %v, want transient", err)
		}
		if attempts != 2 {
			t.Fatalf("attempts = %d, want 2", attempts)
		}
		if len(bodies) != 2 {
			t.Fatalf("bodies = %d, want 2", len(bodies))
		}
		for i, body := range bodies {
			if body.closes != 1 {
				t.Fatalf("body %d closes = %d, want 1", i, body.closes)
			}
		}
	})

	t.Run("successful transport", func(t *testing.T) {
		body := newTrackedFileBody(t)
		cfg := newBodyCloseRequestConfig(t, body)
		cfg.CustomHTTPDoer = httpDoerFunc(func(req *http.Request) (*http.Response, error) {
			if err := req.Body.Close(); err != nil {
				t.Fatal(err)
			}
			return &http.Response{
				StatusCode: http.StatusNoContent,
				Header:     make(http.Header),
				Body:       http.NoBody,
			}, nil
		})

		if err := cfg.Execute(); err != nil {
			t.Fatal(err)
		}
		if body.closes != 1 {
			t.Fatalf("body closes = %d, want 1", body.closes)
		}
	})
}

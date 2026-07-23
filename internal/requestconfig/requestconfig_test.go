package requestconfig

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"testing"
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

func TestClonePreservesCustomHTTPDoer(t *testing.T) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://example.com", nil)
	if err != nil {
		t.Fatal(err)
	}

	called := false
	doer := httpDoerFunc(func(*http.Request) (*http.Response, error) {
		called = true
		return nil, nil
	})
	cfg := &RequestConfig{Request: req, CustomHTTPDoer: doer}

	clone := cfg.Clone(context.Background())
	if clone.CustomHTTPDoer == nil {
		t.Fatal("CustomHTTPDoer is nil after cloning")
	}
	if _, err := clone.CustomHTTPDoer.Do(req); err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Fatal("cloned CustomHTTPDoer did not call the configured client")
	}
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

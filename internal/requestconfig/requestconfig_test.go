package requestconfig

import (
	"context"
	"testing"
)

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

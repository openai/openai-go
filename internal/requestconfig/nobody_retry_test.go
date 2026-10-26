package requestconfig

import (
	"io"
	"net/http"
	"testing"
)

func TestShouldRetryTreatsHTTPNoBodyAsReplayable(t *testing.T) {
	request, err := http.NewRequestWithContext(t.Context(), http.MethodPost, "https://example.test/v1/models", nil)
	if err != nil {
		t.Fatalf("construct empty request: %v", err)
	}
	request.Body = http.NoBody
	response := &http.Response{StatusCode: http.StatusInternalServerError, Header: make(http.Header)}
	if !shouldRetry(request, response, nil) {
		t.Error("http.NoBody was incorrectly classified as a non-replayable request payload")
	}
}

func TestExecuteDoesNotRecreateHTTPNoBodyFromStaleFactory(t *testing.T) {
	cfg := newBodyCloseRequestConfig(t, http.NoBody)
	cfg.MaxRetries = 1
	cfg.Request.GetBody = func() (io.ReadCloser, error) {
		t.Fatal("stale GetBody factory recreated a request explicitly configured without a body")
		return nil, nil
	}
	attempts := 0
	cfg.Middlewares = []middleware{func(request *http.Request, _ middlewareNext) (*http.Response, error) {
		if request.Body != http.NoBody {
			t.Errorf("empty request attempt did not retain its canonical http.NoBody sentinel: %T", request.Body)
		}
		attempts++
		status := http.StatusNoContent
		if attempts == 1 {
			status = http.StatusInternalServerError
		}
		return &http.Response{
			StatusCode: status,
			Header:     http.Header{"Retry-After-Ms": {"0"}},
			Body:       http.NoBody,
			Request:    request,
		}, nil
	}}
	if err := cfg.Execute(); err != nil {
		t.Fatalf("retry configured empty request: %v", err)
	}
	if attempts != 2 {
		t.Errorf("empty request attempts = %d, want two", attempts)
	}
}

func TestClonePreservesHTTPNoBodyWithoutCallingItsFactory(t *testing.T) {
	for _, test := range []struct {
		name         string
		staleFactory bool
	}{
		{name: "nil factory"},
		{name: "stale factory", staleFactory: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			cfg := newBodyCloseRequestConfig(t, http.NoBody)
			cfg.Request.Body = http.NoBody
			cfg.Request.GetBody = nil
			if test.staleFactory {
				cfg.Request.GetBody = func() (io.ReadCloser, error) {
					t.Fatal("stale GetBody factory recreated a cloned empty request body")
					return nil, nil
				}
			}
			clone := cfg.Clone(t.Context())
			if clone == nil || clone.Request.Body != http.NoBody {
				t.Fatalf("cloned canonical empty body = %#v, want http.NoBody", clone)
			}
		})
	}
}

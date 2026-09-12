package openai_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type paginationHTTPDoerFunc func(*http.Request) (*http.Response, error)

func (f paginationHTTPDoerFunc) Do(req *http.Request) (*http.Response, error) {
	return f(req)
}

type paginationRoundTripperFunc func(*http.Request) (*http.Response, error)

func (f paginationRoundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestPaginationPreservesCustomHTTPClient(t *testing.T) {
	customCalls := 0
	customClient := paginationHTTPDoerFunc(func(req *http.Request) (*http.Response, error) {
		customCalls++
		var body string
		switch req.URL.Query().Get("after") {
		case "":
			body = `{"data":[{"id":"job-1"}],"has_more":true}`
		case "job-1":
			body = `{"data":[{"id":"job-2"}],"has_more":false}`
		default:
			return nil, errors.New("unexpected pagination cursor")
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": {"application/json"}},
			Body:       io.NopCloser(strings.NewReader(body)),
			Request:    req,
		}, nil
	})

	fallbackCalls := 0
	fallbackClient := &http.Client{Transport: paginationRoundTripperFunc(func(*http.Request) (*http.Response, error) {
		fallbackCalls++
		return nil, errors.New("fallback HTTP client called")
	})}

	client := openai.NewClient(
		option.WithBaseURL("https://example.com/v1"),
		option.WithAPIKey("test-key"),
		option.WithMaxRetries(0),
		// Install a deterministic native fallback before the bespoke doer. Every
		// page must continue through the latter.
		option.WithHTTPClient(fallbackClient),
		option.WithHTTPClient(customClient),
	)
	pager := client.FineTuning.Jobs.ListAutoPaging(context.Background(), openai.FineTuningJobListParams{})

	var jobIDs []string
	for pager.Next() {
		jobIDs = append(jobIDs, pager.Current().ID)
	}
	if err := pager.Err(); err != nil {
		t.Fatal(err)
	}
	if want := []string{"job-1", "job-2"}; !reflect.DeepEqual(jobIDs, want) {
		t.Fatalf("job IDs = %v, want %v", jobIDs, want)
	}
	if customCalls != 2 {
		t.Fatalf("custom HTTP client calls = %d, want 2", customCalls)
	}
	if fallbackCalls != 0 {
		t.Fatalf("fallback HTTP client calls = %d, want 0", fallbackCalls)
	}
}

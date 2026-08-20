package openai_test

import (
	"context"
	"errors"
	"io"
	"net/http"
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
	calls := 0
	customClient := paginationHTTPDoerFunc(func(req *http.Request) (*http.Response, error) {
		calls++

		var body string
		switch calls {
		case 1:
			if after := req.URL.Query().Get("after"); after != "" {
				t.Fatalf("first request after = %q, want empty", after)
			}
			body = `{"data":[{"id":"file_1"}],"has_more":true}`
		case 2:
			if after := req.URL.Query().Get("after"); after != "file_1" {
				t.Fatalf("second request after = %q, want %q", after, "file_1")
			}
			body = `{"data":[{"id":"file_2"}],"has_more":false}`
		default:
			t.Fatalf("unexpected request %d", calls)
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(body)),
			Request:    req,
		}, nil
	})
	fallbackClient := &http.Client{
		Transport: paginationRoundTripperFunc(func(*http.Request) (*http.Response, error) {
			return nil, errors.New("fallback HTTP client used")
		}),
	}
	client := openai.NewClient(
		option.WithAPIKey("test-key"),
		option.WithBaseURL("https://example.com/v1/"),
		option.WithHTTPClient(fallbackClient),
		option.WithHTTPClient(customClient),
		option.WithMaxRetries(0),
	)

	page, err := client.Files.List(context.Background(), openai.FileListParams{})
	if err != nil {
		t.Fatal(err)
	}
	nextPage, err := page.GetNextPage()
	if err != nil {
		t.Fatal(err)
	}
	if nextPage == nil {
		t.Fatal("next page is nil")
	}
	if len(nextPage.Data) != 1 || nextPage.Data[0].ID != "file_2" {
		t.Fatalf("next page data = %#v, want file_2", nextPage.Data)
	}
	if calls != 2 {
		t.Fatalf("custom client calls = %d, want 2", calls)
	}
}

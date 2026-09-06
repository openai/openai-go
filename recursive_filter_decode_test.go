package openai_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

func TestResponsesNewDecodesDeepRecursiveFilter(t *testing.T) {
	const depth = 1000
	var filter strings.Builder
	for range depth {
		filter.WriteString(`{"type":"and","filters":[`)
	}
	filter.WriteString(`{"type":"eq","key":"category","value":"reference"}`)
	for range depth {
		filter.WriteString(`]}`)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintf(w, `{"tools":[{"type":"file_search","filters":%s}]}`, filter.String())
	}))
	t.Cleanup(server.Close)

	client := openai.NewClient(
		option.WithAPIKey("test-api-key"),
		option.WithBaseURL(server.URL),
	)
	response, err := client.Responses.New(context.Background(), responses.ResponseNewParams{
		Input: responses.ResponseNewParamsInputUnion{OfString: openai.String("test")},
		Model: "gpt-4o",
	})
	if err != nil {
		t.Fatalf("Responses.New() error = %v", err)
	}

	filters := response.Tools[0].Filters.Filters
	for level := 0; level < depth; level++ {
		if len(filters) != 1 {
			t.Fatalf("level %d filters = %d, want 1", level, len(filters))
		}
		if level == depth-1 {
			if got := filters[0].Key; got != "category" {
				t.Fatalf("comparison key = %q, want category", got)
			}
			return
		}
		filters = filters[0].OfCompoundFilter.Filters
	}
}

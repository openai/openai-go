package openai_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
	"github.com/openai/openai-go/v3/shared"
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

func TestCompoundFilterRetainsInvalidNestedObjectMetadata(t *testing.T) {
	for _, raw := range []string{`42`, `[]`, `"filter"`, `true`} {
		t.Run(raw, func(t *testing.T) {
			var filter shared.CompoundFilter
			if err := json.Unmarshal([]byte(`{"type":"and","filters":[`+raw+`]}`), &filter); err != nil {
				t.Fatalf("Unmarshal() error = %v", err)
			}
			if len(filter.Filters) != 1 {
				t.Fatalf("filters = %d, want 1", len(filter.Filters))
			}
			child := filter.Filters[0]
			if child.JSON.OfCompoundFilter.Valid() {
				t.Fatal("non-object filter was marked as a valid compound filter")
			}
			if got := child.RawJSON(); got != raw {
				t.Fatalf("filter raw JSON = %q, want %q", got, raw)
			}
		})
	}
}

func TestCompoundFilterRetainsNestedFieldMetadata(t *testing.T) {
	var filter shared.CompoundFilter
	if err := json.Unmarshal([]byte(`{"type":"and","filters":[{"type":42,"filters":[]}]}`), &filter); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if len(filter.Filters) != 1 {
		t.Fatalf("filters = %d, want 1", len(filter.Filters))
	}
	child := filter.Filters[0].OfCompoundFilter
	if child.Type != "42" || !child.JSON.Type.Valid() {
		t.Fatalf("nested type = %q, valid = %v; want 42 and valid", child.Type, child.JSON.Type.Valid())
	}
	if got := filter.Filters[0].AsCompoundFilter(); got.Type != child.Type || got.JSON.Type.Valid() != child.JSON.Type.Valid() {
		t.Fatal("nested field metadata differs from AsCompoundFilter")
	}
}

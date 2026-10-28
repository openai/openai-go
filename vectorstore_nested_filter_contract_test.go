package openai_test

import (
	"context"
	"os"
	"testing"

	"github.com/openai/openai-go/v4"
	"github.com/openai/openai-go/v4/internal/testutil"
	"github.com/openai/openai-go/v4/option"
	"github.com/openai/openai-go/v4/shared"
)

func TestVectorStoreSearchWithNestedCompoundFilter(t *testing.T) {
	baseURL := os.Getenv("TEST_API_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:4010"
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}

	leaf := shared.ComparisonFilterParam{
		Key:   "category",
		Type:  shared.ComparisonFilterTypeEq,
		Value: shared.ComparisonFilterValueUnionParam{OfString: openai.String("test")},
	}
	nested := shared.CompoundFilterParam{
		Type:    shared.CompoundFilterTypeOr,
		Filters: []shared.CompoundFilterFilterUnionParam{{OfComparison: &leaf}},
	}
	params := openai.VectorStoreSearchParams{
		Query: openai.VectorStoreSearchParamsQueryUnion{OfString: openai.String("test")},
		Filters: openai.VectorStoreSearchParamsFiltersUnion{OfCompoundFilter: &shared.CompoundFilterParam{
			Type:    shared.CompoundFilterTypeAnd,
			Filters: []shared.CompoundFilterFilterUnionParam{{OfFilter: &nested}},
		}},
	}
	client := openai.NewClient(option.WithBaseURL(baseURL), option.WithUnsafeAllowHTTP(), option.WithAPIKey("test"), option.WithMaxRetries(0))
	if _, err := client.VectorStores.Search(context.Background(), "vs_abc123", params); err != nil {
		t.Fatalf("search with nested compound filter: %v", err)
	}
}

// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/internal/testutil"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/shared"
)

func TestVectorStoreNewWithOptionalParams(t *testing.T) {
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := openai.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
	)
	_, err := client.VectorStores.New(context.TODO(), openai.VectorStoreNewParams{
		ChunkingStrategy: openai.F[openai.FileChunkingStrategyParamUnion](openai.AutoFileChunkingStrategyParam{
			Type: openai.F(openai.AutoFileChunkingStrategyParamTypeAuto),
		}),
		ExpiresAfter: openai.F(openai.VectorStoreNewParamsExpiresAfter{
			Anchor: openai.F(openai.VectorStoreNewParamsExpiresAfterAnchorLastActiveAt),
			Days:   openai.F(int64(1)),
		}),
		FileIDs: openai.F([]string{"string"}),
		Metadata: openai.F(shared.MetadataParam{
			"foo": "string",
		}),
		Name: openai.F("name"),
	})
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestVectorStoreGet(t *testing.T) {
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := openai.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
	)
	_, err := client.VectorStores.Get(context.TODO(), "vector_store_id")
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestVectorStoreUpdateWithOptionalParams(t *testing.T) {
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := openai.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
	)
	_, err := client.VectorStores.Update(
		context.TODO(),
		"vector_store_id",
		openai.VectorStoreUpdateParams{
			ExpiresAfter: openai.F(openai.VectorStoreUpdateParamsExpiresAfter{
				Anchor: openai.F(openai.VectorStoreUpdateParamsExpiresAfterAnchorLastActiveAt),
				Days:   openai.F(int64(1)),
			}),
			Metadata: openai.F(shared.MetadataParam{
				"foo": "string",
			}),
			Name: openai.F("name"),
		},
	)
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestVectorStoreListWithOptionalParams(t *testing.T) {
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := openai.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
	)
	_, err := client.VectorStores.List(context.TODO(), openai.VectorStoreListParams{
		After:  openai.F("after"),
		Before: openai.F("before"),
		Limit:  openai.F(int64(0)),
		Order:  openai.F(openai.VectorStoreListParamsOrderAsc),
	})
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestVectorStoreDelete(t *testing.T) {
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := openai.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
	)
	_, err := client.VectorStores.Delete(context.TODO(), "vector_store_id")
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestVectorStoreSearchWithOptionalParams(t *testing.T) {
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := openai.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
	)
	_, err := client.VectorStores.Search(
		context.TODO(),
		"vs_abc123",
		openai.VectorStoreSearchParams{
			Query: openai.F[openai.VectorStoreSearchParamsQueryUnion](shared.UnionString("string")),
			Filters: openai.F[openai.VectorStoreSearchParamsFiltersUnion](shared.ComparisonFilterParam{
				Key:   openai.F("key"),
				Type:  openai.F(shared.ComparisonFilterTypeEq),
				Value: openai.F[shared.ComparisonFilterValueUnionParam](shared.UnionString("string")),
			}),
			MaxNumResults: openai.F(int64(1)),
			RankingOptions: openai.F(openai.VectorStoreSearchParamsRankingOptions{
				Ranker:         openai.F(openai.VectorStoreSearchParamsRankingOptionsRankerAuto),
				ScoreThreshold: openai.F(0.000000),
			}),
			RewriteQuery: openai.F(true),
		},
	)
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

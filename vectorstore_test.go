// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/Nordlys-Labs/openai-go/v3"
	"github.com/Nordlys-Labs/openai-go/v3/internal/testutil"
	"github.com/Nordlys-Labs/openai-go/v3/option"
	"github.com/Nordlys-Labs/openai-go/v3/shared"
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
		ChunkingStrategy: openai.FileChunkingStrategyParamUnion{
			OfAuto: &openai.AutoFileChunkingStrategyParam{},
		},
		Description: openai.String("description"),
		ExpiresAfter: openai.VectorStoreNewParamsExpiresAfter{
			Days: 1,
		},
		FileIDs: []string{"string"},
		Metadata: shared.Metadata{
			"foo": "string",
		},
		Name: openai.String("name"),
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
			ExpiresAfter: openai.VectorStoreUpdateParamsExpiresAfter{
				Days: 1,
			},
			Metadata: shared.Metadata{
				"foo": "string",
			},
			Name: openai.String("name"),
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
		After:  openai.String("after"),
		Before: openai.String("before"),
		Limit:  openai.Int(0),
		Order:  openai.VectorStoreListParamsOrderAsc,
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
			Query: openai.VectorStoreSearchParamsQueryUnion{
				OfString: openai.String("string"),
			},
			Filters: openai.VectorStoreSearchParamsFiltersUnion{
				OfComparisonFilter: &shared.ComparisonFilterParam{
					Key:  "key",
					Type: shared.ComparisonFilterTypeEq,
					Value: shared.ComparisonFilterValueUnionParam{
						OfString: openai.String("string"),
					},
				},
			},
			MaxNumResults: openai.Int(1),
			RankingOptions: openai.VectorStoreSearchParamsRankingOptions{
				Ranker:         "none",
				ScoreThreshold: openai.Float(0),
			},
			RewriteQuery: openai.Bool(true),
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

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
)

func TestVectorStoreFileBatchNewWithOptionalParams(t *testing.T) {
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
	_, err := client.VectorStores.FileBatches.New(
		context.TODO(),
		"vs_abc123",
		openai.VectorStoreFileBatchNewParams{
			Attributes: map[string]openai.VectorStoreFileBatchNewParamsAttributeUnion{
				"foo": {
					OfString: openai.String("string"),
				},
			},
			ChunkingStrategy: openai.FileChunkingStrategyParamUnion{
				OfAuto: &openai.AutoFileChunkingStrategyParam{},
			},
			FileIDs: []string{"string"},
			Files: []openai.VectorStoreFileBatchNewParamsFile{{
				FileID: "file_id",
				Attributes: map[string]openai.VectorStoreFileBatchNewParamsFileAttributeUnion{
					"foo": {
						OfString: openai.String("string"),
					},
				},
				ChunkingStrategy: openai.FileChunkingStrategyParamUnion{
					OfAuto: &openai.AutoFileChunkingStrategyParam{},
				},
			}},
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

func TestVectorStoreFileBatchGet(t *testing.T) {
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
	_, err := client.VectorStores.FileBatches.Get(
		context.TODO(),
		"vs_abc123",
		"vsfb_abc123",
	)
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestVectorStoreFileBatchCancel(t *testing.T) {
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
	_, err := client.VectorStores.FileBatches.Cancel(
		context.TODO(),
		"vector_store_id",
		"batch_id",
	)
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestVectorStoreFileBatchListFilesWithOptionalParams(t *testing.T) {
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
	_, err := client.VectorStores.FileBatches.ListFiles(
		context.TODO(),
		"vector_store_id",
		"batch_id",
		openai.VectorStoreFileBatchListFilesParams{
			After:  openai.String("after"),
			Before: openai.String("before"),
			Filter: openai.VectorStoreFileBatchListFilesParamsFilterInProgress,
			Limit:  openai.Int(0),
			Order:  openai.VectorStoreFileBatchListFilesParamsOrderAsc,
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

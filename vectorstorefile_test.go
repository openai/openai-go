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
)

func TestVectorStoreFileNewWithOptionalParams(t *testing.T) {
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
	_, err := client.VectorStores.Files.New(
		context.TODO(),
		"vs_abc123",
		openai.VectorStoreFileNewParams{
			FileID: "file_id",
			Attributes: map[string]openai.VectorStoreFileNewParamsAttributeUnion{
				"foo": {
					OfString: openai.String("string"),
				},
			},
			ChunkingStrategy: openai.FileChunkingStrategyParamUnion{
				OfAuto: &openai.AutoFileChunkingStrategyParam{},
			},
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

func TestVectorStoreFileGet(t *testing.T) {
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
	_, err := client.VectorStores.Files.Get(
		context.TODO(),
		"file-abc123",
		openai.VectorStoreFileGetParams{
			VectorStoreID: "vs_abc123",
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

func TestVectorStoreFileUpdate(t *testing.T) {
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
	_, err := client.VectorStores.Files.Update(
		context.TODO(),
		"file-abc123",
		openai.VectorStoreFileUpdateParams{
			VectorStoreID: "vs_abc123",
			Attributes: map[string]openai.VectorStoreFileUpdateParamsAttributeUnion{
				"foo": {
					OfString: openai.String("string"),
				},
			},
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

func TestVectorStoreFileListWithOptionalParams(t *testing.T) {
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
	_, err := client.VectorStores.Files.List(
		context.TODO(),
		"vector_store_id",
		openai.VectorStoreFileListParams{
			After:  openai.String("after"),
			Before: openai.String("before"),
			Filter: openai.VectorStoreFileListParamsFilterInProgress,
			Limit:  openai.Int(0),
			Order:  openai.VectorStoreFileListParamsOrderAsc,
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

func TestVectorStoreFileDelete(t *testing.T) {
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
	_, err := client.VectorStores.Files.Delete(
		context.TODO(),
		"file_id",
		openai.VectorStoreFileDeleteParams{
			VectorStoreID: "vector_store_id",
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

func TestVectorStoreFileContent(t *testing.T) {
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
	_, err := client.VectorStores.Files.Content(
		context.TODO(),
		"file-abc123",
		openai.VectorStoreFileContentParams{
			VectorStoreID: "vs_abc123",
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

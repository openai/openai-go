// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/internal/testutil"
	"github.com/openai/openai-go/v3/option"
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
		"vs_abc123",
		"file-abc123",
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
		"vs_abc123",
		"file-abc123",
		openai.VectorStoreFileUpdateParams{
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
		"vector_store_id",
		"file_id",
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
		"vs_abc123",
		"file-abc123",
	)
	if err != nil {
		var apierr *openai.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestVectorStoreFilePollStatus(t *testing.T) {
	var capturedURLs []string
	callCount := 0

	client := openai.NewClient(
		option.WithAPIKey("My API Key"),
		option.WithHTTPClient(&http.Client{
			Transport: &vectorStoreFileClosureTransport{
				fn: func(req *http.Request) (*http.Response, error) {
					capturedURLs = append(capturedURLs, req.URL.String())
					callCount++

					var status openai.VectorStoreFileStatus
					if callCount < 3 {
						status = openai.VectorStoreFileStatusInProgress
					} else {
						status = openai.VectorStoreFileStatusCompleted
					}

					responseBody := fmt.Sprintf(`{
						"id": "file-abc123",
						"object": "vector_store.file",
						"status": "%s",
						"vector_store_id": "vs_abc123",
						"created_at": 1234567890,
						"usage_bytes": 1024
					}`, status)

					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(responseBody)),
						Header: http.Header{
							"Content-Type":         []string{"application/json"},
							"openai-poll-after-ms": []string{"10"},
						},
					}, nil
				},
			},
		}),
	)

	file, err := client.VectorStores.Files.PollStatus(
		context.Background(),
		"vs_abc123",
		"file-abc123",
		10,
	)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if file.Status != openai.VectorStoreFileStatusCompleted {
		t.Errorf("expected status to be completed, got: %s", file.Status)
	}

	expectedURL := "https://api.openai.com/v1/vector_stores/vs_abc123/files/file-abc123"
	for _, url := range capturedURLs {
		if url != expectedURL {
			t.Errorf("expected URL %s, got: %s", expectedURL, url)
		}
	}

	if callCount != 3 {
		t.Errorf("expected 3 calls, got: %d", callCount)
	}
}

type vectorStoreFileClosureTransport struct {
	fn func(req *http.Request) (*http.Response, error)
}

func (t *vectorStoreFileClosureTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.fn(req)
}

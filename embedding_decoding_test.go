package openai_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func TestEmbeddingResponseDecoding(t *testing.T) {
	cases := []struct {
		name   string
		format openai.EmbeddingNewParamsEncodingFormat
		value  []float64
		valid  bool
		raw    string
	}{
		{"float", openai.EmbeddingNewParamsEncodingFormatFloat, []float64{0.5, -0.25}, true, `[0.5,-0.25]`},
		{"base64", openai.EmbeddingNewParamsEncodingFormatBase64, []float64{0.5, -0.25}, true, `"AAAAPwAAgL4="`},
		{"invalid base64", openai.EmbeddingNewParamsEncodingFormatBase64, nil, false, `"invalid base64"`},
		{"incomplete float32", openai.EmbeddingNewParamsEncodingFormatBase64, nil, false, `"AQ=="`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = fmt.Fprintf(w, `{
					"object":"list","model":"text-embedding-3-small",
					"usage":{"prompt_tokens":1,"total_tokens":1},
					"data":[{"object":"embedding","index":0,"embedding":%s}]
				}`, tc.raw)
			}))
			defer server.Close()

			client := openai.NewClient(option.WithBaseURL(server.URL), option.WithUnsafeAllowHTTP(), option.WithAPIKey("test-key"))
			response, err := client.Embeddings.New(context.Background(), openai.EmbeddingNewParams{
				Input:          openai.EmbeddingNewParamsInputUnion{OfString: openai.String("hello")},
				Model:          openai.EmbeddingModelTextEmbedding3Small,
				EncodingFormat: tc.format,
			})
			if err != nil {
				t.Fatal(err)
			}
			if len(response.Data) != 1 {
				t.Fatalf("got %d embeddings, want 1", len(response.Data))
			}
			got := response.Data[0]
			if !slices.Equal(got.Embedding, tc.value) || got.JSON.Embedding.Valid() != tc.valid || got.JSON.Embedding.Raw() != tc.raw {
				t.Errorf("embedding = %v, valid = %t, raw = %q; want %v, %t, %q", got.Embedding, got.JSON.Embedding.Valid(), got.JSON.Embedding.Raw(), tc.value, tc.valid, tc.raw)
			}
		})
	}
}

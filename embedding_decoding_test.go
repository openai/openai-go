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
		name    string
		format  openai.EmbeddingNewParamsEncodingFormat
		value   []float64
		encoded string
		raw     string
	}{
		{"float", openai.EmbeddingNewParamsEncodingFormatFloat, []float64{0.5, -0.25}, "", `[0.5,-0.25]`},
		{"base64", openai.EmbeddingNewParamsEncodingFormatBase64, nil, "AAAAPwAAgL4=", `"AAAAPwAAgL4="`},
		{"invalid base64", openai.EmbeddingNewParamsEncodingFormatBase64, nil, "invalid base64", `"invalid base64"`},
		{"incomplete float32", openai.EmbeddingNewParamsEncodingFormatBase64, nil, "AQ==", `"AQ=="`},
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
			if !slices.Equal(got.Embedding.OfEmbeddingEmbeddingArray, tc.value) || got.Embedding.OfString != tc.encoded {
				t.Errorf("embedding union = %#v; want floats %v or encoded %q", got.Embedding, tc.value, tc.encoded)
			}
			if tc.value != nil {
				if !slices.Equal(got.Embedding.AsEmbeddingEmbeddingArray(), tc.value) || !got.Embedding.JSON.OfEmbeddingEmbeddingArray.Valid() {
					t.Fatal("float variant did not preserve its values")
				}
			} else if got.Embedding.AsString() != tc.encoded || !got.Embedding.JSON.OfString.Valid() {
				t.Fatal("string variant did not preserve its encoded payload")
			}
			if !got.JSON.Embedding.Valid() || got.JSON.Embedding.Raw() != tc.raw || got.Embedding.RawJSON() != tc.raw {
				t.Errorf("embedding raw = %q, valid = %t; want %q", got.Embedding.RawJSON(), got.JSON.Embedding.Valid(), tc.raw)
			}

		})
	}
}

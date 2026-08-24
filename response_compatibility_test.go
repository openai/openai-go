package openai_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

func TestResponsesPreservesLargeImageGenerationOutputsByDefault(t *testing.T) {
	const imageCount = 3
	const encodedImageBytes int64 = 24 << 20

	parts := []io.Reader{strings.NewReader(`{"id":"resp_test","object":"response","output":[`)}
	for i := 0; i < imageCount; i++ {
		if i != 0 {
			parts = append(parts, strings.NewReader(","))
		}
		parts = append(parts,
			strings.NewReader(`{"id":"ig_test","type":"image_generation_call","status":"completed","result":"`),
			io.LimitReader(repeatedByteReader('A'), encodedImageBytes),
			strings.NewReader(`"}`),
		)
	}
	parts = append(parts, strings.NewReader(`]}`))

	client := newResponseLimitClient(
		io.NopCloser(io.MultiReader(parts...)),
		http.StatusOK,
		"application/json",
	)
	var responseBody []byte
	_, err := client.Responses.New(
		context.Background(),
		responses.ResponseNewParams{Model: "gpt-4o-mini"},
		option.WithResponseBodyInto(&responseBody),
	)
	if err != nil {
		t.Fatalf("Responses.New() error = %v, want large image-generation output preserved", err)
	}
	if int64(len(responseBody)) <= 64<<20 {
		t.Fatalf("response length = %d, want more than the former 64 MiB default", len(responseBody))
	}
}

func TestExecutePreservesLargeBinaryDownloadsByDefault(t *testing.T) {
	const downloadBytes int64 = 64<<20 + 1
	client := newResponseLimitClient(
		io.NopCloser(io.LimitReader(repeatedByteReader('x'), downloadBytes)),
		http.StatusOK,
		"application/octet-stream",
	)
	var data []byte
	err := client.Get(
		context.Background(),
		"files/file_test/content",
		nil,
		&data,
		option.WithHeader("Accept", "application/binary"),
	)
	if err != nil {
		t.Fatalf("Get() error = %v, want large binary download preserved", err)
	}
	if int64(len(data)) != downloadBytes {
		t.Fatalf("download length = %d, want %d", len(data), downloadBytes)
	}
}

func TestExecutePreservesDisabledCompressionBehindOpaqueTransport(t *testing.T) {
	acceptEncoding := make(chan string, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		acceptEncoding <- req.Header.Get("Accept-Encoding")
		w.Header().Set("Content-Type", "application/json")
		if _, err := io.WriteString(w, "{}"); err != nil {
			t.Errorf("write response: %v", err)
		}
	}))
	defer server.Close()

	transport := server.Client().Transport.(*http.Transport).Clone()
	transport.DisableCompression = true
	client := openai.NewClient(
		option.WithAPIKey("test-key"),
		option.WithBaseURL(server.URL+"/"),
		option.WithMaxRetries(0),
		option.WithMaxResponseBodyBytes(2),
		option.WithHTTPClient(&http.Client{
			Transport: responseRoundTripperFunc(func(req *http.Request) (*http.Response, error) {
				return transport.RoundTrip(req)
			}),
		}),
	)
	var response map[string]any
	if err := client.Get(context.Background(), "wrapped", nil, &response); err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if got := <-acceptEncoding; got != "" {
		t.Fatalf("Accept-Encoding = %q, want wrapped disabled-compression policy preserved", got)
	}
}

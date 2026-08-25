package openai_test

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

type managedCompressionRoundTripper struct {
	http.RoundTripper
}

func (managedCompressionRoundTripper) CompressionDisabled() bool { return false }

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
	response, err := client.Responses.New(
		context.Background(),
		responses.ResponseNewParams{Model: "gpt-4o-mini"},
	)
	if err != nil {
		t.Fatalf("Responses.New() error = %v, want large image-generation output preserved", err)
	}
	if len(response.Output) != imageCount {
		t.Fatalf("image generation calls = %d, want %d", len(response.Output), imageCount)
	}
	for i, item := range response.Output {
		image := item.AsImageGenerationCall()
		if int64(len(image.Result)) != encodedImageBytes {
			t.Fatalf("image %d result length = %d, want %d", i, len(image.Result), encodedImageBytes)
		}
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

func TestExecutePreservesLargeStructuredErrorsByDefault(t *testing.T) {
	const detailBytes = 96 << 10
	detail := strings.Repeat("x", detailBytes)
	errorJSON := `{"error":{"message":"request rejected","type":"invalid_request_error","diagnostic":"` + detail + `"}}`
	client := newResponseLimitClient(
		io.NopCloser(strings.NewReader(errorJSON)),
		http.StatusBadRequest,
		"application/json",
	)

	var response map[string]any
	err := client.Get(context.Background(), "large-error", nil, &response)
	var apiErr *openai.Error
	if reflect.TypeOf(err) != reflect.TypeOf(apiErr) || !errors.As(err, &apiErr) {
		t.Fatalf("Get() error type = %T, want unwrapped *openai.Error", err)
	}
	if apiErr.Message != "request rejected" {
		t.Fatalf("API error message = %q, want request rejected", apiErr.Message)
	}
	if got := apiErr.JSON.ExtraFields["diagnostic"].Raw(); got != `"`+detail+`"` {
		t.Fatalf("diagnostic raw length = %d, want %d", len(got), len(detail)+2)
	}
	if !strings.Contains(apiErr.RawJSON(), detail) {
		t.Fatal("RawJSON lost the complete oversized diagnostic")
	}
	if !bytes.Contains(apiErr.DumpResponse(true), []byte(errorJSON)) {
		t.Fatal("DumpResponse lost the complete oversized error response")
	}
}

func TestExecutePreservesOpaqueTransportCompressionWithoutApplicableLimit(t *testing.T) {
	tests := []struct {
		name         string
		customDoer   bool
		raw          bool
		successLimit bool
	}{
		{name: "HTTP client typed response"},
		{name: "custom doer typed response", customDoer: true},
		{name: "HTTP client raw response", raw: true},
		{name: "custom doer raw response", customDoer: true, raw: true},
		{name: "HTTP client raw response with inapplicable success limit", raw: true, successLimit: true},
		{name: "custom doer raw response with inapplicable success limit", customDoer: true, raw: true, successLimit: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			acceptEncoding := make(chan string, 1)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				acceptEncoding <- req.Header.Get("Accept-Encoding")
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Content-Encoding", "gzip")
				compressed := gzip.NewWriter(w)
				if _, err := io.WriteString(compressed, "{}"); err != nil {
					t.Errorf("write gzip response: %v", err)
				}
				if err := compressed.Close(); err != nil {
					t.Errorf("close gzip response: %v", err)
				}
			}))
			defer server.Close()

			native := server.Client().Transport
			httpClient := &http.Client{Transport: responseRoundTripperFunc(native.RoundTrip)}
			var selectedClient option.HTTPClient = httpClient
			if test.customDoer {
				selectedClient = responseDoerFunc(httpClient.Do)
			}
			opts := []option.RequestOption{
				option.WithAPIKey("test-key"),
				option.WithBaseURL(server.URL + "/"),
				option.WithMaxRetries(0),
				option.WithHTTPClient(selectedClient),
			}
			if test.successLimit {
				opts = append(opts, option.WithMaxResponseBodyBytes(1))
			}
			client := openai.NewClient(opts...)

			if test.raw {
				var raw *http.Response
				if err := client.Get(context.Background(), "default-gzip", nil, &raw); err != nil {
					t.Fatalf("Get() raw error = %v", err)
				}
				decoded, err := io.ReadAll(raw.Body)
				if err != nil {
					t.Fatalf("read raw response: %v", err)
				}
				if err := raw.Body.Close(); err != nil {
					t.Fatalf("close raw response: %v", err)
				}
				if string(decoded) != "{}" {
					t.Fatalf("raw decoded response = %q, want {}", decoded)
				}
			} else {
				var response map[string]any
				if err := client.Get(context.Background(), "default-gzip", nil, &response); err != nil {
					t.Fatalf("Get() typed error = %v", err)
				}
			}
			if got := <-acceptEncoding; got != "gzip" {
				t.Fatalf("Accept-Encoding = %q, want native gzip preserved", got)
			}
		})
	}
}

func TestExecutePreservesDisabledCompressionBehindOpaqueTransport(t *testing.T) {
	for _, customDoer := range []bool{false, true} {
		name := "HTTP client transport"
		if customDoer {
			name = "custom HTTP doer"
		}
		t.Run(name, func(t *testing.T) {
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
			httpClient := &http.Client{
				Transport: responseRoundTripperFunc(transport.RoundTrip),
			}
			var selectedClient option.HTTPClient = httpClient
			if customDoer {
				selectedClient = responseDoerFunc(httpClient.Do)
			}
			client := openai.NewClient(
				option.WithAPIKey("test-key"),
				option.WithBaseURL(server.URL+"/"),
				option.WithMaxRetries(0),
				option.WithMaxResponseBodyBytes(2),
				option.WithHTTPClient(selectedClient),
			)
			var response map[string]any
			if err := client.Get(context.Background(), "wrapped", nil, &response); err != nil {
				t.Fatalf("Get() error = %v", err)
			}
			if got := <-acceptEncoding; got != "identity" {
				t.Fatalf("Accept-Encoding = %q, want opaque transport compression safely disabled", got)
			}
		})
	}
}

func TestExecuteOpaqueTransportCannotBypassCompressedWireLimit(t *testing.T) {
	const wireLimit int64 = 64
	var compressed bytes.Buffer
	for i := 0; i < 8; i++ {
		writer := gzip.NewWriter(&compressed)
		if i == 0 {
			if _, err := io.WriteString(writer, "{}"); err != nil {
				t.Fatal(err)
			}
		}
		if err := writer.Close(); err != nil {
			t.Fatal(err)
		}
	}
	if int64(compressed.Len()) <= wireLimit {
		t.Fatalf("compressed fixture length = %d, want more than %d", compressed.Len(), wireLimit)
	}

	for _, status := range []int{http.StatusOK, http.StatusBadRequest} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			acceptEncoding := make(chan string, 1)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				encoding := req.Header.Get("Accept-Encoding")
				acceptEncoding <- encoding
				w.Header().Set("Content-Type", "application/json")
				if strings.Contains(encoding, "gzip") {
					w.Header().Set("Content-Encoding", "gzip")
					w.WriteHeader(status)
					if _, err := w.Write(compressed.Bytes()); err != nil {
						t.Errorf("write compressed response: %v", err)
					}
					return
				}
				w.WriteHeader(status)
				body := "{}"
				if status >= http.StatusBadRequest {
					body = `{"error":{"message":"bounded"}}`
				}
				if _, err := io.WriteString(w, body); err != nil {
					t.Errorf("write response: %v", err)
				}
			}))
			defer server.Close()

			base := server.Client().Transport
			client := openai.NewClient(
				option.WithAPIKey("test-key"),
				option.WithBaseURL(server.URL+"/"),
				option.WithMaxRetries(0),
				option.WithMaxResponseBodyBytes(wireLimit),
				option.WithMaxErrorResponseBodyBytes(wireLimit),
				option.WithHTTPClient(&http.Client{Transport: responseRoundTripperFunc(base.RoundTrip)}),
			)
			var response map[string]any
			err := client.Get(context.Background(), "opaque", nil, &response)
			if status >= http.StatusBadRequest {
				var apiErr *openai.Error
				if !errors.As(err, &apiErr) {
					t.Fatalf("Get() error = %v, want preserved API error", err)
				}
			} else if err != nil {
				t.Fatalf("Get() error = %v", err)
			}
			if got := <-acceptEncoding; got != "identity" {
				t.Fatalf("Accept-Encoding = %q, want identity to prevent unaccounted native gzip", got)
			}
		})
	}
}

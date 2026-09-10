// Handwritten test, not generated. See CONTRIBUTING.md for details.

package openai_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/shared"
)

func TestAPIErrorWithNonObjectErrorBody(t *testing.T) {
	client := openai.NewClient(
		option.WithAPIKey("My API Key"),
		option.WithHTTPClient(&http.Client{
			Transport: &closureTransport{
				fn: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusBadRequest,
						Header:     http.Header{"Content-Type": []string{"application/json"}},
						Body:       io.NopCloser(strings.NewReader(`{"error": "you must provide a model parameter"}`)),
					}, nil
				},
			},
		}),
	)
	_, err := client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{{
			OfUser: &openai.ChatCompletionUserMessageParam{
				Content: openai.ChatCompletionUserMessageParamContentUnion{
					OfString: openai.String("Say this is a test"),
				},
			},
		}},
		Model: shared.ChatModelGPT4o,
	})
	if err == nil {
		t.Fatal("Expected an API error")
	}

	var apiErr *openai.Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("Expected error to be *openai.Error, got %T: %v", err, err)
	}
	if apiErr.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, apiErr.StatusCode)
	}
	if !strings.Contains(apiErr.Message, "you must provide a model parameter") {
		t.Errorf("Expected message to contain the raw body, got %q", apiErr.Message)
	}
	if apiErr.RawJSON() != `{"error": "you must provide a model parameter"}` {
		t.Errorf("Expected RawJSON to return the raw body, got %q", apiErr.RawJSON())
	}
	if !strings.Contains(err.Error(), "you must provide a model parameter") {
		t.Errorf("Expected Error() to contain the raw body, got %q", err.Error())
	}
	if apiErr.Response == nil {
		t.Error("Expected response to be populated")
	}
}

func TestAPIErrorWithNullErrorBody(t *testing.T) {
	client := openai.NewClient(
		option.WithAPIKey("My API Key"),
		option.WithHTTPClient(&http.Client{
			Transport: &closureTransport{
				fn: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusInternalServerError,
						Header:     http.Header{"Content-Type": []string{"application/json"}},
						Body:       io.NopCloser(strings.NewReader(`{"error": null}`)),
					}, nil
				},
			},
		}),
	)
	_, err := client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{{
			OfUser: &openai.ChatCompletionUserMessageParam{
				Content: openai.ChatCompletionUserMessageParamContentUnion{
					OfString: openai.String("Say this is a test"),
				},
			},
		}},
		Model: shared.ChatModelGPT4o,
	})
	if err == nil {
		t.Fatal("Expected an API error")
	}

	var apiErr *openai.Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("Expected error to be *openai.Error, got %T: %v", err, err)
	}
	if apiErr.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, apiErr.StatusCode)
	}
	if apiErr.Message != `{"error": null}` {
		t.Errorf("Expected message to contain the raw body, got %q", apiErr.Message)
	}
	if apiErr.RawJSON() != `{"error": null}` {
		t.Errorf("Expected RawJSON to return the raw body, got %q", apiErr.RawJSON())
	}
	if !strings.Contains(err.Error(), `{"error": null}`) {
		t.Errorf("Expected Error() to contain the raw body, got %q", err.Error())
	}
}

func TestAPIErrorWithNonJSONErrorBody(t *testing.T) {
	client := openai.NewClient(
		option.WithAPIKey("My API Key"),
		option.WithHTTPClient(&http.Client{
			Transport: &closureTransport{
				fn: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusBadGateway,
						Header:     http.Header{"Content-Type": []string{"text/html"}},
						Body:       io.NopCloser(strings.NewReader("<html><body>Bad Gateway</body></html>")),
					}, nil
				},
			},
		}),
	)
	_, err := client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{{
			OfUser: &openai.ChatCompletionUserMessageParam{
				Content: openai.ChatCompletionUserMessageParamContentUnion{
					OfString: openai.String("Say this is a test"),
				},
			},
		}},
		Model: shared.ChatModelGPT4o,
	})
	if err == nil {
		t.Fatal("Expected an API error")
	}

	var apiErr *openai.Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("Expected error to be *openai.Error, got %T: %v", err, err)
	}
	if apiErr.StatusCode != http.StatusBadGateway {
		t.Errorf("Expected status code %d, got %d", http.StatusBadGateway, apiErr.StatusCode)
	}
	if !strings.Contains(apiErr.Message, "Bad Gateway") {
		t.Errorf("Expected message to contain the raw body, got %q", apiErr.Message)
	}
	if !strings.Contains(apiErr.RawJSON(), "Bad Gateway") {
		t.Errorf("Expected RawJSON to contain the raw body, got %q", apiErr.RawJSON())
	}
	if !strings.Contains(err.Error(), "Bad Gateway") {
		t.Errorf("Expected Error() to contain the raw body, got %q", err.Error())
	}
}

func TestAPIErrorWithObjectErrorBody(t *testing.T) {
	client := openai.NewClient(
		option.WithAPIKey("My API Key"),
		option.WithHTTPClient(&http.Client{
			Transport: &closureTransport{
				fn: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusTooManyRequests,
						Header:     http.Header{"Content-Type": []string{"application/json"}},
						Body: io.NopCloser(strings.NewReader(
							`{"error": {"message": "rate limit exceeded", "type": "rate_limit_error", "param": null, "code": "rate_limit_exceeded"}}`,
						)),
					}, nil
				},
			},
		}),
	)
	_, err := client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{{
			OfUser: &openai.ChatCompletionUserMessageParam{
				Content: openai.ChatCompletionUserMessageParamContentUnion{
					OfString: openai.String("Say this is a test"),
				},
			},
		}},
		Model: shared.ChatModelGPT4o,
	})
	if err == nil {
		t.Fatal("Expected an API error")
	}

	var apiErr *openai.Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("Expected error to be *openai.Error, got %T: %v", err, err)
	}
	if apiErr.StatusCode != http.StatusTooManyRequests {
		t.Errorf("Expected status code %d, got %d", http.StatusTooManyRequests, apiErr.StatusCode)
	}
	if apiErr.Message != "rate limit exceeded" {
		t.Errorf("Expected message %q, got %q", "rate limit exceeded", apiErr.Message)
	}
	if apiErr.Type != "rate_limit_error" {
		t.Errorf("Expected type %q, got %q", "rate_limit_error", apiErr.Type)
	}
	if !strings.Contains(apiErr.RawJSON(), "rate limit exceeded") {
		t.Errorf("Expected RawJSON to contain the parsed error payload, got %q", apiErr.RawJSON())
	}
	if !strings.Contains(err.Error(), "rate limit exceeded") {
		t.Errorf("Expected Error() to contain the parsed error payload, got %q", err.Error())
	}
}

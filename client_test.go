// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai_test

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/internal"
	"github.com/openai/openai-go/option"
)

type closureTransport struct {
	fn func(req *http.Request) (*http.Response, error)
}

func (t *closureTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.fn(req)
}

func TestUserAgentHeader(t *testing.T) {
	var userAgent string
	client := openai.NewClient(
		option.WithHTTPClient(&http.Client{
			Transport: &closureTransport{
				fn: func(req *http.Request) (*http.Response, error) {
					userAgent = req.Header.Get("User-Agent")
					return &http.Response{
						StatusCode: http.StatusOK,
					}, nil
				},
			},
		}),
	)
	client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{openai.ChatCompletionUserMessageParam{
			Role:    openai.F(openai.ChatCompletionUserMessageParamRoleUser),
			Content: openai.F([]openai.ChatCompletionContentPartUnionParam{openai.ChatCompletionContentPartTextParam{Text: openai.F("text"), Type: openai.F(openai.ChatCompletionContentPartTextTypeText)}}),
		}}),
		Model: openai.F(openai.ChatModelO1Preview),
	})
	if userAgent != fmt.Sprintf("OpenAI/Go %s", internal.PackageVersion) {
		t.Errorf("Expected User-Agent to be correct, but got: %#v", userAgent)
	}
}

func TestRetryAfter(t *testing.T) {
	retryCountHeaders := make([]string, 0)
	client := openai.NewClient(
		option.WithHTTPClient(&http.Client{
			Transport: &closureTransport{
				fn: func(req *http.Request) (*http.Response, error) {
					retryCountHeaders = append(retryCountHeaders, req.Header.Get("X-Stainless-Retry-Count"))
					return &http.Response{
						StatusCode: http.StatusTooManyRequests,
						Header: http.Header{
							http.CanonicalHeaderKey("Retry-After"): []string{"0.1"},
						},
					}, nil
				},
			},
		}),
	)
	res, err := client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{openai.ChatCompletionUserMessageParam{
			Role:    openai.F(openai.ChatCompletionUserMessageParamRoleUser),
			Content: openai.F([]openai.ChatCompletionContentPartUnionParam{openai.ChatCompletionContentPartTextParam{Text: openai.F("text"), Type: openai.F(openai.ChatCompletionContentPartTextTypeText)}}),
		}}),
		Model: openai.F(openai.ChatModelO1Preview),
	})
	if err == nil || res != nil {
		t.Error("Expected there to be a cancel error and for the response to be nil")
	}

	attempts := len(retryCountHeaders)
	if attempts != 3 {
		t.Errorf("Expected %d attempts, got %d", 3, attempts)
	}

	expectedRetryCountHeaders := []string{"0", "1", "2"}
	if !reflect.DeepEqual(retryCountHeaders, expectedRetryCountHeaders) {
		t.Errorf("Expected %v retry count headers, got %v", expectedRetryCountHeaders, retryCountHeaders)
	}
}

func TestDeleteRetryCountHeader(t *testing.T) {
	retryCountHeaders := make([]string, 0)
	client := openai.NewClient(
		option.WithHTTPClient(&http.Client{
			Transport: &closureTransport{
				fn: func(req *http.Request) (*http.Response, error) {
					retryCountHeaders = append(retryCountHeaders, req.Header.Get("X-Stainless-Retry-Count"))
					return &http.Response{
						StatusCode: http.StatusTooManyRequests,
						Header: http.Header{
							http.CanonicalHeaderKey("Retry-After"): []string{"0.1"},
						},
					}, nil
				},
			},
		}),
		option.WithHeaderDel("X-Stainless-Retry-Count"),
	)
	res, err := client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{openai.ChatCompletionUserMessageParam{
			Role:    openai.F(openai.ChatCompletionUserMessageParamRoleUser),
			Content: openai.F([]openai.ChatCompletionContentPartUnionParam{openai.ChatCompletionContentPartTextParam{Text: openai.F("text"), Type: openai.F(openai.ChatCompletionContentPartTextTypeText)}}),
		}}),
		Model: openai.F(openai.ChatModelO1Preview),
	})
	if err == nil || res != nil {
		t.Error("Expected there to be a cancel error and for the response to be nil")
	}

	expectedRetryCountHeaders := []string{"", "", ""}
	if !reflect.DeepEqual(retryCountHeaders, expectedRetryCountHeaders) {
		t.Errorf("Expected %v retry count headers, got %v", expectedRetryCountHeaders, retryCountHeaders)
	}
}

func TestOverwriteRetryCountHeader(t *testing.T) {
	retryCountHeaders := make([]string, 0)
	client := openai.NewClient(
		option.WithHTTPClient(&http.Client{
			Transport: &closureTransport{
				fn: func(req *http.Request) (*http.Response, error) {
					retryCountHeaders = append(retryCountHeaders, req.Header.Get("X-Stainless-Retry-Count"))
					return &http.Response{
						StatusCode: http.StatusTooManyRequests,
						Header: http.Header{
							http.CanonicalHeaderKey("Retry-After"): []string{"0.1"},
						},
					}, nil
				},
			},
		}),
		option.WithHeader("X-Stainless-Retry-Count", "42"),
	)
	res, err := client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{openai.ChatCompletionUserMessageParam{
			Role:    openai.F(openai.ChatCompletionUserMessageParamRoleUser),
			Content: openai.F([]openai.ChatCompletionContentPartUnionParam{openai.ChatCompletionContentPartTextParam{Text: openai.F("text"), Type: openai.F(openai.ChatCompletionContentPartTextTypeText)}}),
		}}),
		Model: openai.F(openai.ChatModelO1Preview),
	})
	if err == nil || res != nil {
		t.Error("Expected there to be a cancel error and for the response to be nil")
	}

	expectedRetryCountHeaders := []string{"42", "42", "42"}
	if !reflect.DeepEqual(retryCountHeaders, expectedRetryCountHeaders) {
		t.Errorf("Expected %v retry count headers, got %v", expectedRetryCountHeaders, retryCountHeaders)
	}
}

func TestRetryAfterMs(t *testing.T) {
	attempts := 0
	client := openai.NewClient(
		option.WithHTTPClient(&http.Client{
			Transport: &closureTransport{
				fn: func(req *http.Request) (*http.Response, error) {
					attempts++
					return &http.Response{
						StatusCode: http.StatusTooManyRequests,
						Header: http.Header{
							http.CanonicalHeaderKey("Retry-After-Ms"): []string{"100"},
						},
					}, nil
				},
			},
		}),
	)
	res, err := client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{openai.ChatCompletionUserMessageParam{
			Role:    openai.F(openai.ChatCompletionUserMessageParamRoleUser),
			Content: openai.F([]openai.ChatCompletionContentPartUnionParam{openai.ChatCompletionContentPartTextParam{Text: openai.F("text"), Type: openai.F(openai.ChatCompletionContentPartTextTypeText)}}),
		}}),
		Model: openai.F(openai.ChatModelO1Preview),
	})
	if err == nil || res != nil {
		t.Error("Expected there to be a cancel error and for the response to be nil")
	}
	if want := 3; attempts != want {
		t.Errorf("Expected %d attempts, got %d", want, attempts)
	}
}

func TestContextCancel(t *testing.T) {
	client := openai.NewClient(
		option.WithHTTPClient(&http.Client{
			Transport: &closureTransport{
				fn: func(req *http.Request) (*http.Response, error) {
					<-req.Context().Done()
					return nil, req.Context().Err()
				},
			},
		}),
	)
	cancelCtx, cancel := context.WithCancel(context.Background())
	cancel()
	res, err := client.Chat.Completions.New(cancelCtx, openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{openai.ChatCompletionUserMessageParam{
			Role:    openai.F(openai.ChatCompletionUserMessageParamRoleUser),
			Content: openai.F([]openai.ChatCompletionContentPartUnionParam{openai.ChatCompletionContentPartTextParam{Text: openai.F("text"), Type: openai.F(openai.ChatCompletionContentPartTextTypeText)}}),
		}}),
		Model: openai.F(openai.ChatModelO1Preview),
	})
	if err == nil || res != nil {
		t.Error("Expected there to be a cancel error and for the response to be nil")
	}
}

func TestContextCancelDelay(t *testing.T) {
	client := openai.NewClient(
		option.WithHTTPClient(&http.Client{
			Transport: &closureTransport{
				fn: func(req *http.Request) (*http.Response, error) {
					<-req.Context().Done()
					return nil, req.Context().Err()
				},
			},
		}),
	)
	cancelCtx, cancel := context.WithTimeout(context.Background(), 2*time.Millisecond)
	defer cancel()
	res, err := client.Chat.Completions.New(cancelCtx, openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{openai.ChatCompletionUserMessageParam{
			Role:    openai.F(openai.ChatCompletionUserMessageParamRoleUser),
			Content: openai.F([]openai.ChatCompletionContentPartUnionParam{openai.ChatCompletionContentPartTextParam{Text: openai.F("text"), Type: openai.F(openai.ChatCompletionContentPartTextTypeText)}}),
		}}),
		Model: openai.F(openai.ChatModelO1Preview),
	})
	if err == nil || res != nil {
		t.Error("expected there to be a cancel error and for the response to be nil")
	}
}

func TestContextDeadline(t *testing.T) {
	testTimeout := time.After(3 * time.Second)
	testDone := make(chan struct{})

	deadline := time.Now().Add(100 * time.Millisecond)
	deadlineCtx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	go func() {
		client := openai.NewClient(
			option.WithHTTPClient(&http.Client{
				Transport: &closureTransport{
					fn: func(req *http.Request) (*http.Response, error) {
						<-req.Context().Done()
						return nil, req.Context().Err()
					},
				},
			}),
		)
		res, err := client.Chat.Completions.New(deadlineCtx, openai.ChatCompletionNewParams{
			Messages: openai.F([]openai.ChatCompletionMessageParamUnion{openai.ChatCompletionUserMessageParam{
				Role:    openai.F(openai.ChatCompletionUserMessageParamRoleUser),
				Content: openai.F([]openai.ChatCompletionContentPartUnionParam{openai.ChatCompletionContentPartTextParam{Text: openai.F("text"), Type: openai.F(openai.ChatCompletionContentPartTextTypeText)}}),
			}}),
			Model: openai.F(openai.ChatModelO1Preview),
		})
		if err == nil || res != nil {
			t.Error("expected there to be a deadline error and for the response to be nil")
		}
		close(testDone)
	}()

	select {
	case <-testTimeout:
		t.Fatal("client didn't finish in time")
	case <-testDone:
		if diff := time.Since(deadline); diff < -30*time.Millisecond || 30*time.Millisecond < diff {
			t.Fatalf("client did not return within 30ms of context deadline, got %s", diff)
		}
	}
}

// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/auth"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/shared"
)

type mockSubjectTokenProvider struct {
	token     string
	tokenType auth.SubjectTokenType
	callCount int
	mu        sync.Mutex
}

func (m *mockSubjectTokenProvider) TokenType() auth.SubjectTokenType {
	return m.tokenType
}

func (m *mockSubjectTokenProvider) GetToken(ctx context.Context, doer auth.HTTPDoer) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.callCount++
	return m.token, nil
}

func (m *mockSubjectTokenProvider) GetCallCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.callCount
}

func testWorkloadIdentity(provider auth.SubjectTokenProvider) auth.WorkloadIdentity {
	return auth.WorkloadIdentity{
		ClientID:           "test-client-id",
		IdentityProviderID: "test-idp-id",
		ServiceAccountID:   "test-sa-id",
		Provider:           provider,
	}
}

func TestClientWorkloadIdentityInitialization(t *testing.T) {
	provider := &mockSubjectTokenProvider{
		token:     "test-subject-token",
		tokenType: auth.SubjectTokenTypeJWT,
	}

	var capturedAuthHeader string
	oauthCallCount := 0
	apiCallCount := 0

	mockHTTPClient := &http.Client{
		Transport: &closureTransport{
			fn: func(req *http.Request) (*http.Response, error) {
				urlStr := req.URL.String()
				if strings.Contains(urlStr, "/oauth/token") || strings.Contains(urlStr, "auth.openai.com") {
					oauthCallCount++
					tokenResp := map[string]interface{}{
						"access_token": "exchanged-token-123",
						"expires_in":   3600,
					}
					body, _ := json.Marshal(tokenResp)
					headers := make(http.Header)
					headers.Set("Content-Type", "application/json")
					return &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(strings.NewReader(string(body))),
						Header:     headers,
					}, nil
				}

				apiCallCount++
				capturedAuthHeader = req.Header.Get("Authorization")
				headers := make(http.Header)
				headers.Set("Content-Type", "application/json")
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader(`{"id": "chatcmpl-123", "object": "chat.completion", "created": 1234567890, "model": "gpt-4", "choices": [{"index": 0, "message": {"role": "assistant", "content": "test"}, "finish_reason": "stop"}]}`)),
					Header:     headers,
				}, nil
			},
		},
	}

	client := openai.NewClient(
		option.WithOrganization("org-123"),
		option.WithProject("proj-456"),
		option.WithWorkloadIdentity(testWorkloadIdentity(provider)),
		option.WithHTTPClient(mockHTTPClient),
	)

	_, err := client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{{
			OfUser: &openai.ChatCompletionUserMessageParam{
				Content: openai.ChatCompletionUserMessageParamContentUnion{
					OfString: openai.String("test"),
				},
			},
		}},
		Model: shared.ChatModelGPT4o,
	})

	if err != nil {
		t.Fatalf("API call failed: %v", err)
	}

	if oauthCallCount != 1 {
		t.Errorf("OAuth call count = %d, want 1", oauthCallCount)
	}

	if apiCallCount != 1 {
		t.Errorf("API call count = %d, want 1", apiCallCount)
	}

	expectedAuthHeader := "Bearer exchanged-token-123"
	if capturedAuthHeader != expectedAuthHeader {
		t.Errorf("Authorization header = %q, want %q", capturedAuthHeader, expectedAuthHeader)
	}
}

func TestWorkloadIdentity401Retry(t *testing.T) {
	provider := &mockSubjectTokenProvider{
		token:     "test-subject-token",
		tokenType: auth.SubjectTokenTypeJWT,
	}

	oauthCallCount := 0
	apiCallCount := 0

	mockHTTPClient := &http.Client{
		Transport: &closureTransport{
			fn: func(req *http.Request) (*http.Response, error) {
				if strings.Contains(req.URL.String(), "/oauth/token") {
					oauthCallCount++
					tokenID := fmt.Sprintf("token-%d", oauthCallCount)
					tokenResp := map[string]interface{}{
						"access_token": tokenID,
						"expires_in":   3600,
					}
					body, _ := json.Marshal(tokenResp)
					return &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(strings.NewReader(string(body))),
						Header:     http.Header{"Content-Type": []string{"application/json"}},
					}, nil
				}

				apiCallCount++
				if apiCallCount == 1 {
					return &http.Response{
						StatusCode: 401,
						Body:       io.NopCloser(strings.NewReader(`{"error": {"message": "Unauthorized"}}`)),
						Header:     http.Header{"Content-Type": []string{"application/json"}},
					}, nil
				}

				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader(`{"id": "chatcmpl-123", "object": "chat.completion", "created": 1234567890, "model": "gpt-4", "choices": [{"index": 0, "message": {"role": "assistant", "content": "test"}, "finish_reason": "stop"}]}`)),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}, nil
			},
		},
	}

	client := openai.NewClient(
		option.WithOrganization("org-123"),
		option.WithProject("proj-456"),
		option.WithWorkloadIdentity(testWorkloadIdentity(provider)),
		option.WithHTTPClient(mockHTTPClient),
	)

	_, err := client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{{
			OfUser: &openai.ChatCompletionUserMessageParam{
				Content: openai.ChatCompletionUserMessageParamContentUnion{
					OfString: openai.String("test"),
				},
			},
		}},
		Model: shared.ChatModelGPT4o,
	})

	if err != nil {
		t.Fatalf("API call failed after retry: %v", err)
	}

	if oauthCallCount != 2 {
		t.Errorf("OAuth call count = %d, want 2 (initial + retry)", oauthCallCount)
	}

	if apiCallCount != 2 {
		t.Errorf("API call count = %d, want 2 (401 + retry)", apiCallCount)
	}

	if provider.GetCallCount() != 2 {
		t.Errorf("Provider call count = %d, want 2", provider.GetCallCount())
	}
}

func TestWorkloadIdentitySingleRetry(t *testing.T) {
	provider := &mockSubjectTokenProvider{
		token:     "test-subject-token",
		tokenType: auth.SubjectTokenTypeJWT,
	}

	oauthCallCount := 0
	apiCallCount := 0

	mockHTTPClient := &http.Client{
		Transport: &closureTransport{
			fn: func(req *http.Request) (*http.Response, error) {
				if strings.Contains(req.URL.String(), "/oauth/token") {
					oauthCallCount++
					tokenResp := map[string]interface{}{
						"access_token": "test-token",
						"expires_in":   3600,
					}
					body, _ := json.Marshal(tokenResp)
					return &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(strings.NewReader(string(body))),
						Header:     http.Header{"Content-Type": []string{"application/json"}},
					}, nil
				}

				apiCallCount++
				return &http.Response{
					StatusCode: 401,
					Body:       io.NopCloser(strings.NewReader(`{"error": {"message": "Unauthorized"}}`)),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}, nil
			},
		},
	}

	client := openai.NewClient(
		option.WithOrganization("org-123"),
		option.WithProject("proj-456"),
		option.WithWorkloadIdentity(testWorkloadIdentity(provider)),
		option.WithHTTPClient(mockHTTPClient),
	)

	_, err := client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{{
			OfUser: &openai.ChatCompletionUserMessageParam{
				Content: openai.ChatCompletionUserMessageParamContentUnion{
					OfString: openai.String("test"),
				},
			},
		}},
		Model: shared.ChatModelGPT4o,
	})

	if err == nil {
		t.Fatal("Expected error after failed retry, got nil")
	}

	if oauthCallCount != 2 {
		t.Errorf("OAuth call count = %d, want 2 (initial + one retry only)", oauthCallCount)
	}

	if apiCallCount != 2 {
		t.Errorf("API call count = %d, want 2 (initial 401 + one retry)", apiCallCount)
	}
}

func TestWorkloadIdentityTokenInjection(t *testing.T) {
	provider := &mockSubjectTokenProvider{
		token:     "test-subject-token-456",
		tokenType: auth.SubjectTokenTypeJWT,
	}

	var capturedAuthHeader string

	mockHTTPClient := &http.Client{
		Transport: &closureTransport{
			fn: func(req *http.Request) (*http.Response, error) {
				if strings.Contains(req.URL.String(), "/oauth/token") {
					tokenResp := map[string]interface{}{
						"access_token": "exchanged-token-789",
						"expires_in":   3600,
					}
					body, _ := json.Marshal(tokenResp)
					return &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(strings.NewReader(string(body))),
						Header:     http.Header{"Content-Type": []string{"application/json"}},
					}, nil
				}

				capturedAuthHeader = req.Header.Get("Authorization")
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader(`{"id": "chatcmpl-123", "object": "chat.completion", "created": 1234567890, "model": "gpt-4", "choices": [{"index": 0, "message": {"role": "assistant", "content": "test"}, "finish_reason": "stop"}]}`)),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}, nil
			},
		},
	}

	client := openai.NewClient(
		option.WithOrganization("org-123"),
		option.WithProject("proj-456"),
		option.WithWorkloadIdentity(testWorkloadIdentity(provider)),
		option.WithHTTPClient(mockHTTPClient),
	)

	_, err := client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{{
			OfUser: &openai.ChatCompletionUserMessageParam{
				Content: openai.ChatCompletionUserMessageParamContentUnion{
					OfString: openai.String("test"),
				},
			},
		}},
		Model: shared.ChatModelGPT4o,
	})

	if err != nil {
		t.Fatalf("API call failed: %v", err)
	}

	expectedAuthHeader := "Bearer exchanged-token-789"
	if capturedAuthHeader != expectedAuthHeader {
		t.Errorf("Authorization header = %q, want %q", capturedAuthHeader, expectedAuthHeader)
	}

	if strings.Contains(capturedAuthHeader, "workload-identity-auth") {
		t.Error("Authorization header should not contain placeholder 'workload-identity-auth'")
	}
}

func TestWorkloadIdentity401RetryReusesReplayableBody(t *testing.T) {
	provider := &mockSubjectTokenProvider{
		token:     "test-subject-token",
		tokenType: auth.SubjectTokenTypeJWT,
	}

	oauthCallCount := 0
	apiBodies := []string{}
	mockHTTPClient := &http.Client{
		Transport: &closureTransport{
			fn: func(req *http.Request) (*http.Response, error) {
				if strings.Contains(req.URL.String(), "/oauth/token") {
					oauthCallCount++
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(fmt.Sprintf(`{"access_token":"token-%d","expires_in":3600}`, oauthCallCount))),
						Header:     http.Header{"Content-Type": []string{"application/json"}},
					}, nil
				}

				payload, err := io.ReadAll(req.Body)
				if err != nil {
					t.Fatalf("failed reading request body: %v", err)
				}
				apiBodies = append(apiBodies, string(payload))

				if len(apiBodies) == 1 {
					return &http.Response{
						StatusCode: http.StatusUnauthorized,
						Body:       io.NopCloser(strings.NewReader(`{"error": {"message": "Unauthorized"}}`)),
						Header:     http.Header{"Content-Type": []string{"application/json"}},
					}, nil
				}

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`{"ok":true}`)),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}, nil
			},
		},
	}

	client := openai.NewClient(
		option.WithOrganization("org-123"),
		option.WithProject("proj-456"),
		option.WithWorkloadIdentity(testWorkloadIdentity(provider)),
		option.WithHTTPClient(mockHTTPClient),
	)

	var res map[string]any
	err := client.Execute(
		context.Background(),
		http.MethodPost,
		"/custom",
		nil,
		&res,
		option.WithRequestBody("application/json", []byte(`{"hello":"world"}`)),
	)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if oauthCallCount != 2 {
		t.Errorf("OAuth call count = %d, want 2", oauthCallCount)
	}
	if len(apiBodies) != 2 {
		t.Fatalf("API call count = %d, want 2", len(apiBodies))
	}
	if apiBodies[0] != `{"hello":"world"}` {
		t.Errorf("first request body = %q, want %q", apiBodies[0], `{"hello":"world"}`)
	}
	if apiBodies[1] != apiBodies[0] {
		t.Errorf("second request body = %q, want %q", apiBodies[1], apiBodies[0])
	}
}

func TestWorkloadIdentity401RetrySkipsNonReplayableBody(t *testing.T) {
	provider := &mockSubjectTokenProvider{
		token:     "test-subject-token",
		tokenType: auth.SubjectTokenTypeJWT,
	}

	oauthCallCount := 0
	apiCallCount := 0
	mockHTTPClient := &http.Client{
		Transport: &closureTransport{
			fn: func(req *http.Request) (*http.Response, error) {
				if strings.Contains(req.URL.String(), "/oauth/token") {
					oauthCallCount++
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(`{"access_token":"token-1","expires_in":3600}`)),
						Header:     http.Header{"Content-Type": []string{"application/json"}},
					}, nil
				}

				apiCallCount++
				return &http.Response{
					StatusCode: http.StatusUnauthorized,
					Body:       io.NopCloser(strings.NewReader(`{"error": {"message": "Unauthorized"}}`)),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}, nil
			},
		},
	}

	client := openai.NewClient(
		option.WithOrganization("org-123"),
		option.WithProject("proj-456"),
		option.WithWorkloadIdentity(testWorkloadIdentity(provider)),
		option.WithHTTPClient(mockHTTPClient),
	)

	err := client.Execute(
		context.Background(),
		http.MethodPost,
		"/custom",
		nil,
		nil,
		option.WithRequestBody("application/json", strings.NewReader(`{"hello":"world"}`)),
	)
	if err == nil {
		t.Fatal("Expected Execute() to return an error")
	}

	if oauthCallCount != 1 {
		t.Errorf("OAuth call count = %d, want 1", oauthCallCount)
	}
	if apiCallCount != 1 {
		t.Errorf("API call count = %d, want 1", apiCallCount)
	}
}

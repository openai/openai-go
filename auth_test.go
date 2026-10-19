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
					body, err := json.Marshal(tokenResp)
					if err != nil {
						return nil, fmt.Errorf("marshal OAuth token response: %w", err)
					}
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
		Model: shared.ChatModel("gpt-4o"),
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
					body, err := json.Marshal(tokenResp)
					if err != nil {
						return nil, fmt.Errorf("marshal OAuth token response: %w", err)
					}
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
		Model: shared.ChatModel("gpt-4o"),
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
					body, err := json.Marshal(tokenResp)
					if err != nil {
						return nil, fmt.Errorf("marshal OAuth token response: %w", err)
					}
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
		Model: shared.ChatModel("gpt-4o"),
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
					body, err := json.Marshal(tokenResp)
					if err != nil {
						return nil, fmt.Errorf("marshal OAuth token response: %w", err)
					}
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
		Model: shared.ChatModel("gpt-4o"),
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

func TestFineTuningCheckpointPermissionsUseAdminAPIKey(t *testing.T) {
	const permissionsPath = "/v1/fine_tuning/checkpoints/ft-checkpoint/permissions"
	const permissionList = `{"object":"list","data":[{"id":"perm_123","object":"checkpoint.permission","created_at":1,"project_id":"proj_coverage"}],"has_more":false,"first_id":"perm_123","last_id":"perm_123"}`
	tests := []struct {
		name     string
		method   string
		path     string
		query    string
		body     string
		response string
		call     func(*openai.Client) error
	}{
		{
			name: "New", method: http.MethodPost, path: permissionsPath,
			body: `{"project_ids":["proj_coverage"]}`, response: permissionList,
			call: func(client *openai.Client) error {
				_, err := client.FineTuning.Checkpoints.Permissions.New(context.Background(), "ft-checkpoint", openai.FineTuningCheckpointPermissionNewParams{ProjectIDs: []string{"proj_coverage"}})
				return err
			},
		},
		{
			name: "Get", method: http.MethodGet, path: permissionsPath,
			query: "limit=2&project_id=proj_coverage", response: permissionList,
			call: func(client *openai.Client) error {
				_, err := client.FineTuning.Checkpoints.Permissions.Get(context.Background(), "ft-checkpoint", openai.FineTuningCheckpointPermissionGetParams{Limit: openai.Int(2), ProjectID: openai.String("proj_coverage")})
				return err
			},
		},
		{
			name: "List", method: http.MethodGet, path: permissionsPath,
			query: "limit=2&project_id=proj_coverage", response: permissionList,
			call: func(client *openai.Client) error {
				_, err := client.FineTuning.Checkpoints.Permissions.List(context.Background(), "ft-checkpoint", openai.FineTuningCheckpointPermissionListParams{Limit: openai.Int(2), ProjectID: openai.String("proj_coverage")})
				return err
			},
		},
		{
			name: "Delete", method: http.MethodDelete, path: permissionsPath + "/perm_123",
			response: `{"id":"perm_123","object":"checkpoint.permission","deleted":true}`,
			call: func(client *openai.Client) error {
				_, err := client.FineTuning.Checkpoints.Permissions.Delete(context.Background(), "ft-checkpoint", "perm_123")
				return err
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			client := openai.NewClient(
				option.WithAPIKey("fake-project-key"),
				option.WithAdminAPIKey("fake-admin-key"),
				option.WithBaseURL("https://api.openai.test/v1/"),
				option.WithMaxRetries(0),
				option.WithHTTPClient(&http.Client{Transport: &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
					calls++
					if req.Header.Get("Authorization") != "Bearer fake-admin-key" {
						t.Error("request did not select the configured admin API key")
					}
					if req.Method != tt.method || req.URL.Host != "api.openai.test" || req.URL.Path != tt.path || req.URL.RawQuery != tt.query {
						t.Errorf("request = %s %s, want %s https://api.openai.test%s with query %q", req.Method, req.URL, tt.method, tt.path, tt.query)
					}
					var body []byte
					if req.Body != nil {
						var err error
						body, err = io.ReadAll(req.Body)
						closeErr := req.Body.Close()
						if err != nil {
							return nil, err
						}
						if closeErr != nil {
							return nil, closeErr
						}
					}
					if string(body) != tt.body {
						t.Errorf("request body = %q, want %q", body, tt.body)
					}
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(tt.response)),
						Header:     http.Header{"Content-Type": []string{"application/json"}},
					}, nil
				}}}),
			)
			if err := tt.call(&client); err != nil {
				t.Fatalf("%s() error = %v", tt.name, err)
			}
			if calls != 1 {
				t.Errorf("request count = %d, want 1", calls)
			}
		})
	}
}

func TestExecuteDefaultsToAPIKeyWhenBothKeysConfigured(t *testing.T) {
	var capturedAuthHeader string

	client := openai.NewClient(
		option.WithAPIKey("My API Key"),
		option.WithAdminAPIKey("My Admin API Key"),
		option.WithHTTPClient(&http.Client{
			Transport: &closureTransport{
				fn: func(req *http.Request) (*http.Response, error) {
					capturedAuthHeader = req.Header.Get("Authorization")
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(`{"ok":true}`)),
						Header:     http.Header{"Content-Type": []string{"application/json"}},
					}, nil
				},
			},
		}),
	)

	var res map[string]any
	err := client.Execute(context.Background(), http.MethodGet, "/admin", nil, &res)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if capturedAuthHeader != "Bearer My API Key" {
		t.Errorf("Authorization header = %q, want %q", capturedAuthHeader, "Bearer My API Key")
	}
}

func TestExecuteDefaultsToAdminAPIKeyWhenOnlyAdminAPIKeyConfigured(t *testing.T) {
	var capturedAuthHeader string

	client := openai.NewClient(
		option.WithAPIKey(""),
		option.WithAdminAPIKey("My Admin API Key"),
		option.WithHTTPClient(&http.Client{
			Transport: &closureTransport{
				fn: func(req *http.Request) (*http.Response, error) {
					capturedAuthHeader = req.Header.Get("Authorization")
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(`{"ok":true}`)),
						Header:     http.Header{"Content-Type": []string{"application/json"}},
					}, nil
				},
			},
		}),
	)

	var res map[string]any
	err := client.Execute(context.Background(), http.MethodGet, "/admin", nil, &res)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if capturedAuthHeader != "Bearer My Admin API Key" {
		t.Errorf("Authorization header = %q, want %q", capturedAuthHeader, "Bearer My Admin API Key")
	}
}

func TestExecuteRequestAPIKeyOverridesClientAdminAPIKey(t *testing.T) {
	var capturedAuthHeader string

	client := openai.NewClient(
		option.WithAPIKey("My API Key"),
		option.WithAdminAPIKey("My Admin API Key"),
		option.WithHTTPClient(&http.Client{
			Transport: &closureTransport{
				fn: func(req *http.Request) (*http.Response, error) {
					capturedAuthHeader = req.Header.Get("Authorization")
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(`{"ok":true}`)),
						Header:     http.Header{"Content-Type": []string{"application/json"}},
					}, nil
				},
			},
		}),
	)

	var res map[string]any
	err := client.Execute(
		context.Background(),
		http.MethodGet,
		"/custom",
		nil,
		&res,
		option.WithAPIKey("My Request API Key"),
	)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if capturedAuthHeader != "Bearer My Request API Key" {
		t.Errorf("Authorization header = %q, want %q", capturedAuthHeader, "Bearer My Request API Key")
	}
}

func TestExecuteRequestAdminAPIKeyOverridesClientAPIKey(t *testing.T) {
	var capturedAuthHeader string

	client := openai.NewClient(
		option.WithAPIKey("My API Key"),
		option.WithHTTPClient(&http.Client{
			Transport: &closureTransport{
				fn: func(req *http.Request) (*http.Response, error) {
					capturedAuthHeader = req.Header.Get("Authorization")
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(`{"ok":true}`)),
						Header:     http.Header{"Content-Type": []string{"application/json"}},
					}, nil
				},
			},
		}),
	)

	var res map[string]any
	err := client.Execute(
		context.Background(),
		http.MethodGet,
		"/admin",
		nil,
		&res,
		option.WithAdminAPIKey("My Request Admin API Key"),
	)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if capturedAuthHeader != "Bearer My Request Admin API Key" {
		t.Errorf("Authorization header = %q, want %q", capturedAuthHeader, "Bearer My Request Admin API Key")
	}
}

func TestRequestAPIKeyOverridesClientAuthorizationHeader(t *testing.T) {
	var capturedAuthHeader string

	client := openai.NewClient(
		option.WithAPIKey("My API Key"),
		option.WithHeader("Authorization", "Bearer Stale API Key"),
		option.WithHTTPClient(&http.Client{
			Transport: &closureTransport{
				fn: func(req *http.Request) (*http.Response, error) {
					capturedAuthHeader = req.Header.Get("Authorization")
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(`{"id":"chatcmpl-123","object":"chat.completion","created":1234567890,"model":"gpt-4","choices":[{"index":0,"message":{"role":"assistant","content":"test"},"finish_reason":"stop"}]}`)),
						Header:     http.Header{"Content-Type": []string{"application/json"}},
					}, nil
				},
			},
		}),
	)

	_, err := client.Chat.Completions.New(
		context.Background(),
		openai.ChatCompletionNewParams{
			Messages: []openai.ChatCompletionMessageParamUnion{{
				OfUser: &openai.ChatCompletionUserMessageParam{
					Content: openai.ChatCompletionUserMessageParamContentUnion{
						OfString: openai.String("test"),
					},
				},
			}},
			Model: shared.ChatModel("gpt-4o"),
		},
		option.WithAPIKey("My Request API Key"),
	)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if capturedAuthHeader != "Bearer My Request API Key" {
		t.Errorf("Authorization header = %q, want %q", capturedAuthHeader, "Bearer My Request API Key")
	}
}

package auth_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/openai/openai-go/v3/auth"
)

type mockProvider struct {
	token     string
	tokenType auth.SubjectTokenType
	callCount int
	delay     time.Duration
	err       error
	mu        sync.Mutex
}

func (m *mockProvider) TokenType() auth.SubjectTokenType {
	return m.tokenType
}

func (m *mockProvider) GetToken(ctx context.Context, _ auth.HTTPDoer) (string, error) {
	m.mu.Lock()
	m.callCount++
	m.mu.Unlock()

	if m.delay > 0 {
		time.Sleep(m.delay)
	}

	if m.err != nil {
		return "", m.err
	}

	return m.token, nil
}

func (m *mockProvider) GetCallCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.callCount
}

type closureTransport struct {
	fn func(req *http.Request) (*http.Response, error)
}

func (t *closureTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.fn(req)
}

func mockOAuthServer(responseBody string, statusCode int) *http.Client {
	return &http.Client{
		Transport: &closureTransport{
			fn: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: statusCode,
					Body:       io.NopCloser(strings.NewReader(responseBody)),
					Header:     make(http.Header),
				}, nil
			},
		},
	}
}

func TestTokenCaching(t *testing.T) {
	provider := &mockProvider{
		token:     "test-subject-token",
		tokenType: auth.SubjectTokenTypeJWT,
	}

	responseBody := `{"access_token": "exchanged-token-123", "expires_in": 3600}`
	httpClient := mockOAuthServer(responseBody, 200)

	config := auth.WorkloadIdentity{ClientID: "client-id", IdentityProviderID: "idp-id", ServiceAccountID: "sa-id", Provider: provider}
	wa, err := auth.NewWorkloadIdentityAuth(config)
	if err != nil {
		t.Fatalf("NewWorkloadIdentityAuth() error = %v", err)
	}

	token1, err := wa.GetToken(context.Background(), httpClient)
	if err != nil {
		t.Fatalf("First GetToken() error = %v", err)
	}

	token2, err := wa.GetToken(context.Background(), httpClient)
	if err != nil {
		t.Fatalf("Second GetToken() error = %v", err)
	}

	if token1 != token2 {
		t.Errorf("Tokens don't match: %q != %q", token1, token2)
	}

	if provider.GetCallCount() != 1 {
		t.Errorf("Provider call count = %d, want 1", provider.GetCallCount())
	}
}

func TestTokenExpiration(t *testing.T) {
	callCount := 0
	mu := sync.Mutex{}

	transport := &closureTransport{
		fn: func(req *http.Request) (*http.Response, error) {
			mu.Lock()
			callCount++
			mu.Unlock()

			responseBody := `{"access_token": "token-` + string(rune('0'+callCount)) + `", "expires_in": 0}`
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(responseBody)),
				Header:     make(http.Header),
			}, nil
		},
	}
	httpClient := &http.Client{Transport: transport}

	provider := &mockProvider{
		token:     "test-subject-token",
		tokenType: auth.SubjectTokenTypeJWT,
	}

	config := auth.WorkloadIdentity{ClientID: "client-id", IdentityProviderID: "idp-id", ServiceAccountID: "sa-id", Provider: provider}
	wa, err := auth.NewWorkloadIdentityAuth(config)
	if err != nil {
		t.Fatalf("NewWorkloadIdentityAuth() error = %v", err)
	}

	token1, err := wa.GetToken(context.Background(), httpClient)
	if err != nil {
		t.Fatalf("First GetToken() error = %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	token2, err := wa.GetToken(context.Background(), httpClient)
	if err != nil {
		t.Fatalf("Second GetToken() error = %v", err)
	}

	if token1 == token2 {
		t.Errorf("Expected different tokens after expiration, got same token: %q", token1)
	}

	if provider.GetCallCount() < 2 {
		t.Errorf("Provider call count = %d, want at least 2", provider.GetCallCount())
	}
}

func TestConcurrentDeduplication(t *testing.T) {
	provider := &mockProvider{
		token:     "test-subject-token",
		tokenType: auth.SubjectTokenTypeJWT,
		delay:     100 * time.Millisecond,
	}

	oauthCallCount := 0
	var oauthMu sync.Mutex

	transport := &closureTransport{
		fn: func(req *http.Request) (*http.Response, error) {
			oauthMu.Lock()
			oauthCallCount++
			oauthMu.Unlock()

			responseBody := `{"access_token": "exchanged-token-123", "expires_in": 3600}`
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(responseBody)),
				Header:     make(http.Header),
			}, nil
		},
	}
	httpClient := &http.Client{Transport: transport}

	config := auth.WorkloadIdentity{ClientID: "client-id", IdentityProviderID: "idp-id", ServiceAccountID: "sa-id", Provider: provider}
	wa, err := auth.NewWorkloadIdentityAuth(config)
	if err != nil {
		t.Fatalf("NewWorkloadIdentityAuth() error = %v", err)
	}

	const numGoroutines = 5
	type result struct {
		token string
		err   error
	}
	resultsChan := make(chan result, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			token, err := wa.GetToken(context.Background(), httpClient)
			resultsChan <- result{token: token, err: err}
		}()
	}

	results := make([]result, 0, numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		results = append(results, <-resultsChan)
	}

	for i, r := range results {
		if r.err != nil {
			t.Errorf("Goroutine %d error = %v", i, r.err)
		}
		if r.token == "" {
			t.Errorf("Goroutine %d got empty token (error: %v)", i, r.err)
		}
	}

	expectedToken := "exchanged-token-123"
	for i, r := range results {
		if r.token != expectedToken {
			t.Errorf("Goroutine %d got token %q, want %q", i, r.token, expectedToken)
		}
	}

	if provider.GetCallCount() != 1 {
		t.Errorf("Provider call count = %d, want 1", provider.GetCallCount())
	}

	oauthMu.Lock()
	finalOAuthCallCount := oauthCallCount
	oauthMu.Unlock()

	if finalOAuthCallCount != 1 {
		t.Errorf("OAuth call count = %d, want 1 (deduplication should prevent multiple calls)", finalOAuthCallCount)
	}
}

func TestProactiveRefresh(t *testing.T) {
	provider := &mockProvider{
		token:     "test-subject-token",
		tokenType: auth.SubjectTokenTypeJWT,
	}

	callCount := 0
	mu := sync.Mutex{}

	refreshSignal := make(chan struct{}, 10)
	transport := &closureTransport{
		fn: func(req *http.Request) (*http.Response, error) {
			mu.Lock()
			callCount++
			currentCount := callCount
			mu.Unlock()

			select {
			case refreshSignal <- struct{}{}:
			default:
			}

			responseBody := `{"access_token": "token-` + string(rune('0'+currentCount)) + `", "expires_in": 1}`
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(responseBody)),
				Header:     make(http.Header),
			}, nil
		},
	}
	httpClient := &http.Client{Transport: transport}

	bufferSeconds := 0
	config := auth.WorkloadIdentity{
		ClientID:             "client-id",
		IdentityProviderID:   "idp-id",
		ServiceAccountID:     "sa-id",
		Provider:             provider,
		RefreshBufferSeconds: bufferSeconds,
	}
	wa, err := auth.NewWorkloadIdentityAuth(config)
	if err != nil {
		t.Fatalf("NewWorkloadIdentityAuth() error = %v", err)
	}

	token1, err := wa.GetToken(context.Background(), httpClient)
	if err != nil {
		t.Fatalf("First GetToken() error = %v", err)
	}
	<-refreshSignal

	token2, err := wa.GetToken(context.Background(), httpClient)
	if err != nil {
		t.Fatalf("Second GetToken() error = %v", err)
	}

	if token1 != token2 {
		t.Log("Tokens differ, background refresh may have completed")
	}

	select {
	case <-refreshSignal:
	case <-time.After(5 * time.Second):
		t.Error("timed out waiting for background refresh")
	}

	mu.Lock()
	finalCallCount := callCount
	mu.Unlock()

	if finalCallCount < 2 {
		t.Errorf("OAuth call count = %d, want at least 2 (initial + background refresh)", finalCallCount)
	}

	if provider.GetCallCount() < 2 {
		t.Errorf("Provider call count = %d, want at least 2", provider.GetCallCount())
	}
}

func TestOAuthErrorHandling(t *testing.T) {
	provider := &mockProvider{
		token:     "test-subject-token",
		tokenType: auth.SubjectTokenTypeJWT,
	}

	testCases := []struct {
		statusCode         int
		shouldBeOAuthError bool
	}{
		{400, true},
		{401, true},
		{403, true},
		{500, false},
	}

	for _, tc := range testCases {
		errorBody := `{"error": "invalid_grant", "error_description": "Token exchange failed"}`
		httpClient := mockOAuthServer(errorBody, tc.statusCode)

		config := auth.WorkloadIdentity{ClientID: "client-id", IdentityProviderID: "idp-id", ServiceAccountID: "sa-id", Provider: provider}
		wa, err := auth.NewWorkloadIdentityAuth(config)
		if err != nil {
			t.Fatalf("NewWorkloadIdentityAuth() error = %v", err)
		}

		_, err = wa.GetToken(context.Background(), httpClient)
		if err == nil {
			t.Errorf("Status %d: expected error, got nil", tc.statusCode)
			continue
		}

		_, isOAuthError := err.(*auth.OAuthError)
		if isOAuthError != tc.shouldBeOAuthError {
			t.Errorf("Status %d: isOAuthError = %v, want %v", tc.statusCode, isOAuthError, tc.shouldBeOAuthError)
		}

		if tc.shouldBeOAuthError {
			oauthErr := err.(*auth.OAuthError)
			if oauthErr.StatusCode != tc.statusCode {
				t.Errorf("StatusCode = %d, want %d", oauthErr.StatusCode, tc.statusCode)
			}
			if oauthErr.ErrorCode != "invalid_grant" {
				t.Errorf("ErrorCode = %q, want %q", oauthErr.ErrorCode, "invalid_grant")
			}
			if oauthErr.ErrorDescription != "Token exchange failed" {
				t.Errorf("ErrorDescription = %q, want %q", oauthErr.ErrorDescription, "Token exchange failed")
			}
		}
	}
}

func TestDefaultValues(t *testing.T) {
	provider := &mockProvider{
		token:     "test-subject-token",
		tokenType: auth.SubjectTokenTypeJWT,
	}

	responseBody := `{"access_token": "test-token"}`
	httpClient := mockOAuthServer(responseBody, 200)

	config := auth.WorkloadIdentity{ClientID: "client-id", IdentityProviderID: "idp-id", ServiceAccountID: "sa-id", Provider: provider}
	wa, err := auth.NewWorkloadIdentityAuth(config)
	if err != nil {
		t.Fatalf("NewWorkloadIdentityAuth() error = %v", err)
	}

	_, err = wa.GetToken(context.Background(), httpClient)
	if err != nil {
		t.Fatalf("GetToken() error = %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	_, err = wa.GetToken(context.Background(), httpClient)
	if err != nil {
		t.Fatalf("Second GetToken() error = %v", err)
	}

	if provider.GetCallCount() != 1 {
		t.Errorf("Provider call count = %d, want 1 (token should be cached with 3600s default)", provider.GetCallCount())
	}

	customBuffer := 300
	config2 := auth.WorkloadIdentity{
		ClientID:             "client-id",
		IdentityProviderID:   "idp-id",
		ServiceAccountID:     "sa-id",
		Provider:             provider,
		RefreshBufferSeconds: customBuffer,
	}
	wa2, err2 := auth.NewWorkloadIdentityAuth(config2)
	if err2 != nil {
		t.Fatalf("NewWorkloadIdentityAuth() error = %v", err2)
	}

	_, err = wa2.GetToken(context.Background(), httpClient)
	if err != nil {
		t.Fatalf("GetToken() with custom buffer error = %v", err)
	}
}

package auth_test

import (
	"context"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/openai/openai-go/v3/auth"
)

func TestK8sProviderFileReading(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "k8s-token-*")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	tokenContent := "  test-jwt-token-123  \n"
	if _, err := tmpFile.WriteString(tokenContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	provider := auth.K8sServiceAccountTokenProvider(tmpFile.Name())

	token, err := provider.GetToken(context.Background(), nil)
	if err != nil {
		t.Fatalf("GetToken() error = %v", err)
	}

	expectedToken := "test-jwt-token-123"
	if token != expectedToken {
		t.Errorf("GetToken() = %q, want %q", token, expectedToken)
	}

	if provider.TokenType() != auth.SubjectTokenTypeJWT {
		t.Errorf("TokenType() = %v, want %v", provider.TokenType(), auth.SubjectTokenTypeJWT)
	}
}

func TestK8sProviderDefaultPath(t *testing.T) {
	provider := auth.K8sServiceAccountTokenProvider("")

	defaultPath := "/var/run/secrets/kubernetes.io/serviceaccount/token"

	_, err := provider.GetToken(context.Background(), nil)
	if err == nil {
		t.Log("Default path file exists, skipping validation")
		return
	}

	providerErr, ok := err.(*auth.SubjectTokenProviderError)
	if !ok {
		t.Fatalf("Expected *SubjectTokenProviderError, got %T", err)
	}

	if providerErr.Provider != "kubernetes" {
		t.Errorf("Provider = %q, want %q", providerErr.Provider, "kubernetes")
	}

	if !strings.Contains(providerErr.Error(), defaultPath) {
		t.Errorf("Error should reference default path %q: %v", defaultPath, providerErr)
	}
}

func TestK8sProviderErrorHandling(t *testing.T) {
	provider := auth.K8sServiceAccountTokenProvider("/nonexistent/path/to/token")

	_, err := provider.GetToken(context.Background(), nil)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	providerErr, ok := err.(*auth.SubjectTokenProviderError)
	if !ok {
		t.Fatalf("Expected *SubjectTokenProviderError, got %T", err)
	}

	if providerErr.Provider != "kubernetes" {
		t.Errorf("Provider = %q, want %q", providerErr.Provider, "kubernetes")
	}

	if providerErr.Cause == nil {
		t.Error("Expected Cause to be set")
	}
}

func TestK8sProviderEmptyToken(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "k8s-token-empty-*")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString("   \n   "); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	provider := auth.K8sServiceAccountTokenProvider(tmpFile.Name())

	_, err = provider.GetToken(context.Background(), nil)
	if err == nil {
		t.Fatal("Expected error for empty token, got nil")
	}

	providerErr, ok := err.(*auth.SubjectTokenProviderError)
	if !ok {
		t.Fatalf("Expected *SubjectTokenProviderError, got %T", err)
	}

	if providerErr.Provider != "kubernetes" {
		t.Errorf("Provider = %q, want %q", providerErr.Provider, "kubernetes")
	}

	if !strings.Contains(err.Error(), "empty") {
		t.Errorf("Error should mention 'empty': %v", err)
	}
}

func TestAzureProviderTokenType(t *testing.T) {
	provider := auth.AzureManagedIdentityTokenProvider(nil)

	if provider.TokenType() != auth.SubjectTokenTypeJWT {
		t.Errorf("TokenType() = %v, want %v", provider.TokenType(), auth.SubjectTokenTypeJWT)
	}
}

func TestAzureProviderGetToken(t *testing.T) {
	mockClient := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("Metadata") != "true" {
				t.Fatalf("Metadata header = %q, want true", req.Header.Get("Metadata"))
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"access_token":"azure-token-123"}`)),
				Header:     make(http.Header),
			}, nil
		}),
	}

	provider := auth.AzureManagedIdentityTokenProvider(nil)
	token, err := provider.GetToken(context.Background(), mockClient)
	if err != nil {
		t.Fatalf("GetToken() error = %v", err)
	}
	if token != "azure-token-123" {
		t.Errorf("GetToken() = %q, want %q", token, "azure-token-123")
	}
}

func TestGCPProviderTokenType(t *testing.T) {
	provider := auth.GCPIDTokenProvider(nil)

	if provider.TokenType() != auth.SubjectTokenTypeID {
		t.Errorf("TokenType() = %v, want %v", provider.TokenType(), auth.SubjectTokenTypeID)
	}
}

func TestAzureProviderCustomResource(t *testing.T) {
	mockClient := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if !strings.Contains(req.URL.RawQuery, "resource=https%3A%2F%2Fcustom.openai.com") {
				t.Fatalf("Expected custom resource in URL, got %q", req.URL.RawQuery)
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"access_token":"azure-custom-token"}`)),
				Header:     make(http.Header),
			}, nil
		}),
	}

	provider := auth.AzureManagedIdentityTokenProvider(&auth.AzureManagedIdentityTokenProviderConfig{
		Resource: "https://custom.openai.com",
	})
	token, err := provider.GetToken(context.Background(), mockClient)
	if err != nil {
		t.Fatalf("GetToken() error = %v", err)
	}
	if token != "azure-custom-token" {
		t.Errorf("GetToken() = %q, want %q", token, "azure-custom-token")
	}
}

func TestGCPProviderCustomAudience(t *testing.T) {
	mockClient := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if !strings.Contains(req.URL.RawQuery, "audience=https%3A%2F%2Fcustom.openai.com") {
				t.Fatalf("Expected custom audience in URL, got %q", req.URL.RawQuery)
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader("gcp-custom-token-jwt")),
				Header:     make(http.Header),
			}, nil
		}),
	}

	provider := auth.GCPIDTokenProvider(&auth.GCPIDTokenProviderConfig{
		Audience: "https://custom.openai.com",
	})
	token, err := provider.GetToken(context.Background(), mockClient)
	if err != nil {
		t.Fatalf("GetToken() error = %v", err)
	}
	if token != "gcp-custom-token-jwt" {
		t.Errorf("GetToken() = %q, want %q", token, "gcp-custom-token-jwt")
	}
}

func TestAzureProviderNilConfigBackwardCompatibility(t *testing.T) {
	mockClient := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if !strings.Contains(req.URL.RawQuery, "resource=https%3A%2F%2Fmanagement.azure.com%2F") {
				t.Fatalf("Expected default resource in URL, got %q", req.URL.RawQuery)
			}
			if !strings.Contains(req.URL.RawQuery, "api-version=2018-02-01") {
				t.Fatalf("Expected default API version in URL, got %q", req.URL.RawQuery)
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"access_token":"default-token"}`)),
				Header:     make(http.Header),
			}, nil
		}),
	}

	provider := auth.AzureManagedIdentityTokenProvider(nil)
	_, err := provider.GetToken(context.Background(), mockClient)
	if err != nil {
		t.Fatalf("GetToken() error = %v", err)
	}
}

func TestGCPProviderNilConfigBackwardCompatibility(t *testing.T) {
	mockClient := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if !strings.Contains(req.URL.RawQuery, "audience=https%3A%2F%2Fapi.openai.com") {
				t.Fatalf("Expected default audience in URL, got %q", req.URL.RawQuery)
			}
			expectedPath := "/computeMetadata/v1/instance/service-accounts/default/identity"
			if req.URL.Path != expectedPath {
				t.Fatalf("Expected path %q, got %q", expectedPath, req.URL.Path)
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader("default-gcp-token")),
				Header:     make(http.Header),
			}, nil
		}),
	}

	provider := auth.GCPIDTokenProvider(nil)
	_, err := provider.GetToken(context.Background(), mockClient)
	if err != nil {
		t.Fatalf("GetToken() error = %v", err)
	}
}

func TestAzureProviderResourceURLEncoding(t *testing.T) {
	mockClient := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if !strings.Contains(req.URL.RawQuery, "resource=https%3A%2F%2Fapi.openai.com%2Fv1%2Fspecial") {
				t.Fatalf("Expected URL-encoded resource with path, got %q", req.URL.RawQuery)
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"access_token":"encoded-token"}`)),
				Header:     make(http.Header),
			}, nil
		}),
	}

	provider := auth.AzureManagedIdentityTokenProvider(&auth.AzureManagedIdentityTokenProviderConfig{
		Resource: "https://api.openai.com/v1/special",
	})
	_, err := provider.GetToken(context.Background(), mockClient)
	if err != nil {
		t.Fatalf("GetToken() error = %v", err)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

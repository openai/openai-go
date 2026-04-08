package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	DefaultK8STokenPath    = "/var/run/secrets/kubernetes.io/serviceaccount/token"
	DefaultAudience        = "https://api.openai.com/v1"
	DefaultAzureResource   = "https://management.azure.com/"
	DefaultAzureAPIVersion = "2018-02-01"
)

type k8sServiceAccountTokenProvider struct {
	tokenPath string
}

func K8sServiceAccountTokenProvider(tokenPath string) SubjectTokenProvider {
	if tokenPath == "" {
		tokenPath = DefaultK8STokenPath
	}
	return &k8sServiceAccountTokenProvider{tokenPath: tokenPath}
}

func (p *k8sServiceAccountTokenProvider) TokenType() SubjectTokenType {
	return SubjectTokenTypeJWT
}

func (p *k8sServiceAccountTokenProvider) GetToken(ctx context.Context, _ HTTPDoer) (string, error) {
	data, err := os.ReadFile(p.tokenPath)
	if err != nil {
		return "", &SubjectTokenProviderError{
			Provider: "kubernetes",
			Message:  fmt.Sprintf("failed to read service account token from %s", p.tokenPath),
			Cause:    err,
		}
	}
	token := strings.TrimSpace(string(data))
	if token == "" {
		return "", &SubjectTokenProviderError{
			Provider: "kubernetes",
			Message:  "service account token is empty",
		}
	}
	return token, nil
}

type AzureManagedIdentityTokenProviderConfig struct {
	Resource   string
	ObjectID   string
	ClientID   string
	MSIResID   string
	APIVersion string
}

type azureManagedIdentityTokenProvider struct {
	config AzureManagedIdentityTokenProviderConfig
}

func AzureManagedIdentityTokenProvider(config *AzureManagedIdentityTokenProviderConfig) SubjectTokenProvider {
	if config == nil {
		config = &AzureManagedIdentityTokenProviderConfig{}
	}
	cfg := *config
	if cfg.Resource == "" {
		cfg.Resource = DefaultAzureResource
	}
	if cfg.APIVersion == "" {
		cfg.APIVersion = DefaultAzureAPIVersion
	}
	return &azureManagedIdentityTokenProvider{config: cfg}
}

func (p *azureManagedIdentityTokenProvider) TokenType() SubjectTokenType {
	return SubjectTokenTypeJWT
}

func (p *azureManagedIdentityTokenProvider) GetToken(ctx context.Context, httpClient HTTPDoer) (string, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	params := url.Values{}
	params.Set("api-version", p.config.APIVersion)
	params.Set("resource", p.config.Resource)
	if p.config.ObjectID != "" {
		params.Set("object_id", p.config.ObjectID)
	}
	if p.config.ClientID != "" {
		params.Set("client_id", p.config.ClientID)
	}
	if p.config.MSIResID != "" {
		params.Set("msi_res_id", p.config.MSIResID)
	}

	endpoint := "http://169.254.169.254/metadata/identity/oauth2/token?" + params.Encode()
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return "", &SubjectTokenProviderError{
			Provider: "azure-imds",
			Message:  "failed to create request",
			Cause:    err,
		}
	}
	req.Header.Set("Metadata", "true")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", &SubjectTokenProviderError{
			Provider: "azure-imds",
			Message:  "failed to fetch token from IMDS",
			Cause:    err,
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, readErr := io.ReadAll(resp.Body)
		msg := fmt.Sprintf("IMDS returned status %d", resp.StatusCode)

		if readErr != nil {
			msg = fmt.Sprintf("%s (failed to read body: %v)", msg, readErr)
		} else if len(body) > 0 {
			msg = fmt.Sprintf("%s: %s", msg, string(body))
		}

		return "", &SubjectTokenProviderError{
			Provider: "azure-imds",
			Message:  msg,
		}
	}

	var result struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", &SubjectTokenProviderError{
			Provider: "azure-imds",
			Message:  "failed to decode IMDS response",
			Cause:    err,
		}
	}

	if result.AccessToken == "" {
		return "", &SubjectTokenProviderError{
			Provider: "azure-imds",
			Message:  "IMDS response missing 'access_token' field",
		}
	}

	return result.AccessToken, nil
}

type GCPIDTokenProviderConfig struct {
	Audience string
}

type gcpIDTokenProvider struct {
	config GCPIDTokenProviderConfig
}

func GCPIDTokenProvider(config *GCPIDTokenProviderConfig) SubjectTokenProvider {
	if config == nil {
		config = &GCPIDTokenProviderConfig{}
	}
	cfg := *config
	if cfg.Audience == "" {
		cfg.Audience = DefaultAudience
	}
	return &gcpIDTokenProvider{config: cfg}
}

func (p *gcpIDTokenProvider) TokenType() SubjectTokenType {
	return SubjectTokenTypeID
}

func (p *gcpIDTokenProvider) GetToken(ctx context.Context, httpClient HTTPDoer) (string, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	endpoint := "http://metadata.google.internal/computeMetadata/v1/instance/service-accounts/default/identity"
	params := url.Values{}
	params.Set("audience", p.config.Audience)
	endpoint = endpoint + "?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return "", &SubjectTokenProviderError{
			Provider: "gcp-metadata",
			Message:  "failed to create request",
			Cause:    err,
		}
	}
	req.Header.Set("Metadata-Flavor", "Google")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", &SubjectTokenProviderError{
			Provider: "gcp-metadata",
			Message:  "failed to fetch token from metadata server",
			Cause:    err,
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", &SubjectTokenProviderError{
			Provider: "gcp-metadata",
			Message:  fmt.Sprintf("metadata server returned status %d: %s", resp.StatusCode, string(body)),
		}
	}

	token, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", &SubjectTokenProviderError{
			Provider: "gcp-metadata",
			Message:  "failed to read response body",
			Cause:    err,
		}
	}

	tokenStr := strings.TrimSpace(string(token))
	if tokenStr == "" {
		return "", &SubjectTokenProviderError{
			Provider: "gcp-metadata",
			Message:  "metadata server returned empty token",
		}
	}

	return tokenStr, nil
}

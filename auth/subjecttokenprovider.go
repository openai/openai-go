package auth

import (
	"bytes"
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
	config         AzureManagedIdentityTokenProviderConfig
	metadataClient HTTPDoer
}

// AzureManagedIdentityTokenProvider retrieves tokens directly from Azure IMDS.
// Metadata requests bypass the HTTP client passed to GetToken and use a five-second
// timeout, or the caller context deadline if earlier.
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
	return &azureManagedIdentityTokenProvider{config: cfg, metadataClient: newMetadataHTTPClient()}
}

func (p *azureManagedIdentityTokenProvider) TokenType() SubjectTokenType {
	return SubjectTokenTypeJWT
}

func (p *azureManagedIdentityTokenProvider) GetToken(ctx context.Context, _ HTTPDoer) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, metadataRequestTimeout)
	defer cancel()

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

	endpoint := azureMetadataURL + "?" + params.Encode()
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return "", &SubjectTokenProviderError{
			Provider: "azure-imds",
			Message:  "failed to create request",
			Cause:    err,
		}
	}
	req.Header.Set("Metadata", "true")

	resp, err := workloadIdentityDo(p.metadataClient, req)
	if err != nil {
		return "", &SubjectTokenProviderError{
			Provider: "azure-imds",
			Message:  "failed to fetch token from IMDS",
			Cause:    err,
		}
	}
	if resp == nil || resp.Body == nil {
		return "", &SubjectTokenProviderError{
			Provider: "azure-imds",
			Message:  "IMDS returned an invalid response",
		}
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := readSubjectTokenProviderResponse(ctx, resp, "azure-imds", "IMDS")
	if err != nil {
		return "", err
	}

	var result struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&result); err != nil {
		return "", &SubjectTokenProviderError{
			Provider: "azure-imds",
			Message:  "failed to decode IMDS response",
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
	config         GCPIDTokenProviderConfig
	metadataClient HTTPDoer
}

// GCPIDTokenProvider retrieves tokens directly from the Compute Engine metadata
// service. Metadata requests bypass the HTTP client passed to GetToken and use a
// five-second timeout, or the caller context deadline if earlier.
func GCPIDTokenProvider(config *GCPIDTokenProviderConfig) SubjectTokenProvider {
	if config == nil {
		config = &GCPIDTokenProviderConfig{}
	}
	cfg := *config
	if cfg.Audience == "" {
		cfg.Audience = DefaultAudience
	}
	return &gcpIDTokenProvider{config: cfg, metadataClient: newMetadataHTTPClient()}
}

func (p *gcpIDTokenProvider) TokenType() SubjectTokenType {
	return SubjectTokenTypeID
}

func (p *gcpIDTokenProvider) GetToken(ctx context.Context, _ HTTPDoer) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, metadataRequestTimeout)
	defer cancel()

	endpoint := gcpMetadataURL
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

	resp, err := workloadIdentityDo(p.metadataClient, req)
	if err != nil {
		return "", &SubjectTokenProviderError{
			Provider: "gcp-metadata",
			Message:  "failed to fetch token from metadata server",
			Cause:    err,
		}
	}
	if resp == nil || resp.Body == nil {
		return "", &SubjectTokenProviderError{
			Provider: "gcp-metadata",
			Message:  "metadata server returned an invalid response",
		}
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode == http.StatusOK {
		if flavors := resp.Header.Values("Metadata-Flavor"); len(flavors) != 1 || flavors[0] != "Google" {
			return "", &SubjectTokenProviderError{
				Provider: "gcp-metadata",
				Message:  "metadata server returned an invalid Metadata-Flavor header",
			}
		}

	}
	token, err := readSubjectTokenProviderResponse(ctx, resp, "gcp-metadata", "metadata server")
	if err != nil {
		return "", err
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

func readSubjectTokenProviderResponse(
	ctx context.Context,
	response *http.Response,
	provider, source string,
) ([]byte, error) {
	maximum := int64(x509SuccessResponseMaximum)
	if response.StatusCode != http.StatusOK {
		maximum = x509ErrorResponseMaximum
	}
	if response.ContentLength > maximum {
		return nil, &SubjectTokenProviderError{
			Provider: provider,
			Message:  fmt.Sprintf("%s response exceeds its size limit", source),
		}
	}
	body, err := io.ReadAll(io.LimitReader(response.Body, maximum+1))
	if err != nil {
		return nil, &SubjectTokenProviderError{
			Provider: provider,
			Message:  fmt.Sprintf("failed to read %s response", source),
			Cause:    ctx.Err(),
		}
	}
	if int64(len(body)) > maximum {
		return nil, &SubjectTokenProviderError{
			Provider: provider,
			Message:  fmt.Sprintf("%s response exceeds its size limit", source),
		}
	}
	if response.StatusCode != http.StatusOK {
		return nil, &SubjectTokenProviderError{
			Provider: provider,
			Message:  fmt.Sprintf("%s returned status %d", source, response.StatusCode),
		}
	}
	return body, nil
}

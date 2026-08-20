// Package azure provides configuration options so you can connect and use Azure OpenAI using the [openai.Client].
//
// Typical usage of this package will look like this:
//
//	client := openai.NewClient(
//		azure.WithEndpoint(azureOpenAIEndpoint, azureOpenAIAPIVersion),
//		azure.WithTokenCredential(azureIdentityTokenCredential),
//		// or azure.WithAPIKey(azureOpenAIAPIKey),
//	)
//
// Or, if you want to construct a specific service:
//
//	client := openai.NewChatCompletionService(
//		azure.WithEndpoint(azureOpenAIEndpoint, azureOpenAIAPIVersion),
//		azure.WithTokenCredential(azureIdentityTokenCredential),
//		// or azure.WithAPIKey(azureOpenAIAPIKey),
//	)
package azure

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
)

const (
	azureProvider            = "Azure"
	azureAPIKeyAuth          = "azure.WithAPIKey"
	azureTokenCredentialAuth = "azure.WithTokenCredential"
)

// WithEndpoint configures this client to connect to an Azure OpenAI endpoint.
//
//   - endpoint - the Azure OpenAI endpoint to connect to. Ex: https://<azure-openai-resource>.openai.azure.com
//   - apiVersion - the Azure OpenAI API version to target (ex: 2024-06-01). See [Azure OpenAI apiversions] for current API versions. This value cannot be empty.
//
// This function should be paired with a call to authenticate, like [azure.WithAPIKey] or [azure.WithTokenCredential], similar to this:
//
//	client := openai.NewClient(
//		azure.WithEndpoint(azureOpenAIEndpoint, azureOpenAIAPIVersion),
//		azure.WithTokenCredential(azureIdentityTokenCredential),
//		// or azure.WithAPIKey(azureOpenAIAPIKey),
//	)
//
// [Azure OpenAI apiversions]: https://learn.microsoft.com/en-us/azure/ai-services/openai/reference#rest-api-versioning
func WithEndpoint(endpoint string, apiVersion string) option.RequestOption {
	if !strings.HasSuffix(endpoint, "/") {
		endpoint += "/"
	}

	withQueryAdd := option.WithQueryAdd("api-version", apiVersion)
	withEndpoint := option.WithBaseURL(endpoint)

	withModelMiddleware := option.WithMiddleware(func(r *http.Request, mn option.MiddlewareNext) (*http.Response, error) {
		replacementPath, err := getReplacementPathWithDeployment(r)

		if err != nil {
			return nil, err
		}

		if err := setEscapedPath(r.URL, replacementPath); err != nil {
			return nil, err
		}
		return mn(r)
	})

	endpointOption := requestconfig.RequestOptionFunc(func(rc *requestconfig.RequestConfig) error {
		if apiVersion == "" {
			return errors.New("apiVersion is an empty string, but needs to be set. See https://learn.microsoft.com/en-us/azure/ai-services/openai/reference#rest-api-versioning for details.")
		}

		return rc.Apply(
			requestconfig.WithProviderEndpoint(azureProvider),
			withQueryAdd,
			withEndpoint,
			withModelMiddleware,
			requestconfig.WithRequestFinalizer(finalizeAzureProvider),
		)
	})

	return requestconfig.WithEnvironmentDefaultsDisabled(endpointOption)
}

type tokenCredentialConfig struct {
	Scopes []string
}

// TokenCredentialOption is the type for any options that can be used to customize
// [WithTokenCredential], including things like using custom scopes.
type TokenCredentialOption func(*tokenCredentialConfig) error

// WithTokenCredentialScopes overrides the default scope used when requesting access tokens.
func WithTokenCredentialScopes(scopes []string) func(*tokenCredentialConfig) error {
	return func(tc *tokenCredentialConfig) error {
		tc.Scopes = scopes
		return nil
	}
}

// WithTokenCredential configures this client to authenticate using an [Azure Identity] TokenCredential.
// This function should be paired with a call to [WithEndpoint] to point to your Azure OpenAI instance.
//
// [Azure Identity]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity
func WithTokenCredential(tokenCredential azcore.TokenCredential, options ...TokenCredentialOption) option.RequestOption {
	auth := requestconfig.NewProviderAuthOption(azureProvider, azureTokenCredentialAuth)
	return requestconfig.RequestOptionFunc(func(rc *requestconfig.RequestConfig) error {
		if err := rc.Apply(requestconfig.WithEndpointProvider(azureProvider), auth); err != nil {
			return err
		}
		tc := &tokenCredentialConfig{
			Scopes: []string{"https://cognitiveservices.azure.com/.default"},
		}

		for _, option := range options {
			if err := option(tc); err != nil {
				return err
			}
		}

		bearerTokenPolicy := runtime.NewBearerTokenPolicy(tokenCredential, tc.Scopes, nil)

		// add in a middleware that uses the bearer token generated from the token credential
		middlewareOption := option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
			if !auth.Selected(rc) {
				return next(req)
			}

			pipeline := runtime.NewPipeline("azopenai-extensions", version, runtime.PipelineOptions{}, &policy.ClientOptions{
				InsecureAllowCredentialWithHTTP: true, // allow for plain HTTP proxies, etc..
				PerRetryPolicies: []policy.Policy{
					bearerTokenPolicy,
					policyAdapter(next),
				},
			})

			req2, err := runtime.NewRequestFromRequest(req)

			if err != nil {
				return nil, err
			}

			return pipeline.Do(req2)
		})

		return rc.Apply(
			option.WithHeaderDel("Api-Key"),
			middlewareOption,
			requestconfig.WithRequestFinalizer(finalizeAzureProvider),
		)
	})
}

// WithAPIKey configures this client to authenticate using an API key.
// This function should be paired with a call to [WithEndpoint] to point to your Azure OpenAI instance.
func WithAPIKey(apiKey string) option.RequestOption {
	auth := requestconfig.NewProviderAuthOption(azureProvider, azureAPIKeyAuth)
	// NOTE: option.WithAPIKey() uses the Authorization header. Azure expects
	// Api-Key instead. Deleting Authorization also prevents request security from
	// automatically injecting environment-derived client credentials.
	return requestconfig.RequestOptionFunc(func(rc *requestconfig.RequestConfig) error {
		if err := rc.Apply(requestconfig.WithEndpointProvider(azureProvider), auth); err != nil {
			return err
		}
		return rc.Apply(
			option.WithHeaderDel("Authorization"),
			option.WithHeader("Api-Key", apiKey),
			requestconfig.WithRequestFinalizer(finalizeAzureProvider),
		)
	})
}

func finalizeAzureProvider(rc *requestconfig.RequestConfig) error {
	if !rc.ProviderEndpointIs(azureProvider) {
		return errors.New("azure: authentication requires azure.WithEndpoint")
	}
	if rc.APIKey != "" || rc.AdminAPIKey != "" {
		return errors.New("azure: Azure authentication cannot be combined with option.WithAPIKey or option.WithAdminAPIKey")
	}

	auth, ok := rc.ProviderAuth(azureProvider)
	if !ok {
		return errors.New("azure: authentication is required; configure exactly one of azure.WithAPIKey or azure.WithTokenCredential")
	}
	if rc.Request.Header.Get("Authorization") != "" {
		return errors.New("azure: Azure authentication cannot be combined with a custom Authorization header")
	}

	switch auth {
	case azureAPIKeyAuth:
		if strings.TrimSpace(rc.Request.Header.Get("Api-Key")) == "" {
			return errors.New("azure: API key must not be empty")
		}
	case azureTokenCredentialAuth:
		if rc.Request.Header.Get("Api-Key") != "" {
			return errors.New("azure: token credential authentication cannot be combined with an Api-Key header")
		}
	default:
		return errors.New("azure: invalid authentication mode")
	}
	return nil
}

// jsonRoutes have JSON payloads - we'll deserialize looking for a .model field in there
// so we won't have to worry about individual types for completions vs embeddings, etc...
var jsonRoutes = map[string]bool{
	"/completions":        true,
	"/chat/completions":   true,
	"/embeddings":         true,
	"/audio/speech":       true,
	"/images/generations": true,
}

// multipartRoutes have mime/multipart payloads. These are less generic - we're very much
// expecting a transcription or translation payload for these.
var multipartRoutes = map[string]bool{
	"/audio/transcriptions": true,
	"/audio/translations":   true,
	"/images/edits":         true,
}

// getReplacementPathWithDeployment parses the request body to extract out the Model parameter (or equivalent)
// (note, the req.Body is fully read as part of this, and is replaced with a bytes.Reader)
func getReplacementPathWithDeployment(req *http.Request) (string, error) {
	if jsonRoutes[req.URL.Path] {
		return getJSONRoute(req)
	}

	if multipartRoutes[req.URL.Path] {
		return getMultipartRoute(req)
	}

	// If route doesn't require deployment ID substitution, just return path with prefix.
	return "/openai" + req.URL.EscapedPath(), nil
}

func setEscapedPath(u *url.URL, escapedPath string) error {
	parsed, err := url.Parse(escapedPath)
	if err != nil {
		return err
	}
	u.Path = parsed.Path
	u.RawPath = parsed.RawPath
	return nil
}

func getJSONRoute(req *http.Request) (string, error) {
	// we need to deserialize the body, partly, in order to read out the model field.
	jsonBytes, err := io.ReadAll(req.Body)

	if err != nil {
		return "", err
	}

	// make sure we restore the body so it can be used in later middlewares.
	req.Body = io.NopCloser(bytes.NewReader(jsonBytes))

	var v *struct {
		Model string `json:"model"`
	}

	if err := json.Unmarshal(jsonBytes, &v); err != nil {
		return "", err
	}

	// Convert path from /chat/completions to /openai/deployments/{deployment-id}/chat/completions
	return requestconfig.FormatPath("/openai/deployments/%s", v.Model) + req.URL.EscapedPath(), nil
}

func getMultipartRoute(req *http.Request) (string, error) {
	// body is a multipart/mime body type instead.
	mimeBytes, err := io.ReadAll(req.Body)

	if err != nil {
		return "", err
	}

	// make sure we restore the body so it can be used in later middlewares.
	req.Body = io.NopCloser(bytes.NewReader(mimeBytes))

	_, mimeParams, err := mime.ParseMediaType(req.Header.Get("Content-Type"))

	if err != nil {
		return "", err
	}

	mimeReader := multipart.NewReader(
		io.NopCloser(bytes.NewReader(mimeBytes)),
		mimeParams["boundary"])

	for {
		mp, err := mimeReader.NextPart()

		if err != nil {
			if errors.Is(err, io.EOF) {
				return "", errors.New("unable to find the model part in multipart body")
			}

			return "", err
		}

		defer func() { _ = mp.Close() }()

		if mp.FormName() == "model" {
			modelBytes, err := io.ReadAll(mp)

			if err != nil {
				return "", err
			}

			// Convert path from /audio/transcriptions to /openai/deployments/{deployment-id}/audio/transcriptions
			return requestconfig.FormatPath("/openai/deployments/%s", string(modelBytes)) + req.URL.EscapedPath(), nil
		}
	}
}

type policyAdapter option.MiddlewareNext

func (mp policyAdapter) Do(req *policy.Request) (*http.Response, error) {
	return (option.MiddlewareNext)(mp)(req.Raw())
}

const version = "v.0.1.0"

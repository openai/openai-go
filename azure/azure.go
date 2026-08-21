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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"reflect"
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
// Authenticated endpoints must use HTTPS unless [WithUnsafeAllowHTTP] is also
// configured for a loopback-only local development endpoint.
// Azure authentication also requires custom networking to use an [*http.Client]
// with a custom [http.RoundTripper] so every redirect destination can be checked.
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
			return nil, requestconfig.WithNoRetryError(err)
		}

		if err := setEscapedPath(r.URL, replacementPath); err != nil {
			return nil, requestconfig.WithNoRetryError(err)
		}
		return mn(r)
	})

	endpointOption := requestconfig.RequestOptionFunc(func(rc *requestconfig.RequestConfig) error {
		if apiVersion == "" {
			return errors.New("apiVersion is an empty string, but needs to be set. See https://learn.microsoft.com/en-us/azure/ai-services/openai/reference#rest-api-versioning for details.")
		}

		return rc.Apply(
			requestconfig.WithProviderEndpointConfigured(azureProvider),
			withQueryAdd,
			withEndpoint,
			withModelMiddleware,
			requestconfig.WithRequestFinalizer(finalizeAzureProvider),
		)
	})

	return requestconfig.WithEnvironmentDefaultsDisabled(endpointOption)
}

type unsafeAllowHTTPContextKey struct{}

type azureCredentialOriginContextKey struct{}

type azureCredentialOrigin struct {
	scheme string
	host   string
	port   string
}

// WithUnsafeAllowHTTP permits Azure credentials to be sent over plaintext HTTP
// only when the final request destination is localhost or a loopback IP address.
// This option is intended exclusively for local development and testing. It
// should never be used in production.
func WithUnsafeAllowHTTP() option.RequestOption {
	return option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
		ctx := context.WithValue(req.Context(), unsafeAllowHTTPContextKey{}, true)
		return next(req.WithContext(ctx))
	})
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
	return requestconfig.RequestOptionFunc(func(rc *requestconfig.RequestConfig) error {
		if isNilTokenCredential(tokenCredential) {
			return errors.New("azure: token credential must not be nil")
		}
		auth := requestconfig.NewProviderAuthOption(azureProvider, azureTokenCredentialAuth)
		if err := rc.Apply(requestconfig.WithEndpointProvider(azureProvider), auth); err != nil {
			return err
		}
		rc.ClearInheritedAuthentication()
		tc := &tokenCredentialConfig{
			Scopes: []string{"https://cognitiveservices.azure.com/.default"},
		}

		for _, option := range options {
			if err := option(tc); err != nil {
				return err
			}
		}

		bearerTokenPolicy := runtime.NewBearerTokenPolicy(tokenCredential, tc.Scopes, nil)
		unsafeBearerTokenPolicy := runtime.NewBearerTokenPolicy(tokenCredential, tc.Scopes, &policy.BearerTokenOptions{
			InsecureAllowCredentialWithHTTP: true,
		})

		// add in a middleware that uses the bearer token generated from the token credential
		middlewareOption := withAzureCredentialMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
			if !auth.Selected(rc) {
				return next(req)
			}

			tokenPolicy := bearerTokenPolicy
			if azureCredentialHTTPAllowed(req) {
				tokenPolicy = unsafeBearerTokenPolicy
			}
			pipeline := runtime.NewPipeline("azopenai-extensions", version, runtime.PipelineOptions{}, &policy.ClientOptions{
				PerRetryPolicies: []policy.Policy{
					tokenPolicy,
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
			middlewareOption,
			requestconfig.WithRequestFinalizer(finalizeAzureProvider),
		)
	})
}

func isNilTokenCredential(tokenCredential azcore.TokenCredential) bool {
	if tokenCredential == nil {
		return true
	}
	value := reflect.ValueOf(tokenCredential)
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return value.IsNil()
	default:
		return false
	}
}

// WithAPIKey configures this client to authenticate using an API key.
// This function should be paired with a call to [WithEndpoint] to point to your Azure OpenAI instance.
func WithAPIKey(apiKey string) option.RequestOption {
	// NOTE: option.WithAPIKey() uses the Authorization header. Azure expects
	// Api-Key instead.
	return requestconfig.RequestOptionFunc(func(rc *requestconfig.RequestConfig) error {
		auth := requestconfig.NewProviderAuthOption(azureProvider, azureAPIKeyAuth)
		if err := rc.Apply(requestconfig.WithEndpointProvider(azureProvider), auth); err != nil {
			return err
		}
		rc.ClearInheritedAuthentication()
		return rc.Apply(
			option.WithHeader("Api-Key", apiKey),
			requestconfig.WithRequestFinalizer(finalizeAzureProvider),
			withAzureCredentialMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
				return next(req)
			}),
		)
	})
}

func finalizeAzureProvider(rc *requestconfig.RequestConfig) error {
	if !rc.ProviderEndpointConfigured(azureProvider) {
		return errors.New("azure: authentication requires azure.WithEndpoint")
	}
	if rc.APIKey != "" || rc.AdminAPIKey != "" {
		return errors.New("azure: Azure authentication cannot be combined with option.WithAPIKey or option.WithAdminAPIKey")
	}

	auth, ok := rc.ProviderAuth(azureProvider)
	if !ok {
		return errors.New("azure: authentication is required; configure exactly one of azure.WithAPIKey or azure.WithTokenCredential")
	}
	if nonEmptyHeaderValues(rc.Request.Header, "Authorization") != 0 {
		return errors.New("azure: Azure authentication cannot be combined with a custom Authorization header")
	}

	switch auth {
	case azureAPIKeyAuth:
		if nonEmptyHeaderValues(rc.Request.Header, "Api-Key") != 1 {
			return errors.New("azure: exactly one non-empty API key is required")
		}
	case azureTokenCredentialAuth:
		if nonEmptyHeaderValues(rc.Request.Header, "Api-Key") != 0 {
			return errors.New("azure: token credential authentication cannot be combined with an Api-Key header")
		}
	default:
		return errors.New("azure: invalid authentication mode")
	}
	return nil
}

func nonEmptyHeaderValues(header http.Header, name string) int {
	count := 0
	for _, value := range header.Values(name) {
		if strings.TrimSpace(value) != "" {
			count++
		}
	}
	return count
}

func withAzureCredentialMiddleware(authenticate option.Middleware) option.RequestOption {
	return requestconfig.WithRequestFinalizer(func(rc *requestconfig.RequestConfig) error {
		if rc.CustomHTTPDoer != nil {
			return errors.New("azure: custom HTTP clients must use *http.Client with a custom RoundTripper so redirects can be validated")
		}

		// Redirects run inside http.Client.Do and don't re-enter SDK middleware.
		// Clone the selected client so every redirect reaches this guard without
		// mutating the caller's client or replacing its CheckRedirect policy.
		client := *rc.HTTPClient
		transport := client.Transport
		if transport == nil {
			transport = http.DefaultTransport
		}
		client.Transport = azureCredentialTransport{base: transport}
		rc.HTTPClient = &client

		return option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
			origin := azureCredentialOriginFromURL(req.URL)
			ctx := context.WithValue(req.Context(), azureCredentialOriginContextKey{}, origin)
			req = req.WithContext(ctx)
			if err := validateAzureCredentialTransport(req); err != nil {
				return nil, err
			}
			return authenticate(req, next)
		}).Apply(rc)
	})
}

type azureCredentialTransport struct {
	base http.RoundTripper
}

func (t azureCredentialTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if err := validateAzureCredentialTransport(req); err != nil {
		if req.Body != nil {
			_ = req.Body.Close()
		}
		return nil, err
	}
	return t.base.RoundTrip(req)
}

func validateAzureCredentialTransport(req *http.Request) error {
	origin, ok := req.Context().Value(azureCredentialOriginContextKey{}).(azureCredentialOrigin)
	if !ok {
		return requestconfig.WithNoRetryError(&azureCredentialOriginError{})
	}

	destination := azureCredentialOriginFromURL(req.URL)
	if destination.scheme != "https" {
		// Unsafe mode may redirect between local development servers, but it
		// must never expand a remote credential origin to a loopback target.
		if azureCredentialOriginIsLoopback(origin) && azureCredentialHTTPAllowed(req) {
			return nil
		}
		return requestconfig.WithNoRetryError(&azureCredentialTransportError{})
	}

	if destination == origin {
		return nil
	}
	return requestconfig.WithNoRetryError(&azureCredentialOriginError{})
}

func azureCredentialOriginFromURL(u *url.URL) azureCredentialOrigin {
	scheme := strings.ToLower(u.Scheme)
	host := strings.ToLower(u.Hostname())
	if ip := net.ParseIP(host); ip != nil {
		host = ip.String()
	}
	port := u.Port()
	if port == "" {
		switch scheme {
		case "http":
			port = "80"
		case "https":
			port = "443"
		}
	}
	return azureCredentialOrigin{scheme: scheme, host: host, port: port}
}

func azureCredentialOriginIsLoopback(origin azureCredentialOrigin) bool {
	if origin.host == "localhost" {
		return true
	}
	ip := net.ParseIP(origin.host)
	return ip != nil && ip.IsLoopback()
}

type azureCredentialTransportError struct{}

func (*azureCredentialTransportError) Error() string {
	return "azure: authenticated requests require HTTPS; WithUnsafeAllowHTTP permits HTTP only for local development on loopback endpoints"
}

// NonRetriable implements the marker understood by the Azure pipeline's retry
// policy. The outer requestconfig marker independently stops SDK retries.
func (*azureCredentialTransportError) NonRetriable() {}

type azureCredentialOriginError struct{}

func (*azureCredentialOriginError) Error() string {
	return "azure: authenticated redirects must remain on the original origin"
}

func (*azureCredentialOriginError) NonRetriable() {}

func azureCredentialHTTPAllowed(req *http.Request) bool {
	if !strings.EqualFold(req.URL.Scheme, "http") {
		return false
	}
	allowed, _ := req.Context().Value(unsafeAllowHTTPContextKey{}).(bool)
	if !allowed {
		return false
	}

	// Keep this allowlist aligned with net/http's ProxyFromEnvironment bypass.
	// Other localhost spellings can resolve to loopback but still reach HTTP_PROXY.
	host := req.URL.Hostname()
	if host == "localhost" {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
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
	if req.Body == nil {
		return "", errors.New("azure: deployment routing requires a JSON request body")
	}

	// we need to deserialize the body, partly, in order to read out the model field.
	jsonBytes, err := io.ReadAll(req.Body)

	if err != nil {
		return "", fmt.Errorf("azure: could not read JSON request body for deployment routing: %w", err)
	}

	// make sure we restore the body so it can be used in later middlewares.
	req.Body = io.NopCloser(bytes.NewReader(jsonBytes))

	var v struct {
		Model string `json:"model"`
	}

	if err := json.Unmarshal(jsonBytes, &v); err != nil {
		return "", fmt.Errorf("azure: could not parse JSON request body for deployment routing: %w", err)
	}
	if v.Model == "" {
		return "", errors.New("azure: deployment routing requires a non-empty model field")
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

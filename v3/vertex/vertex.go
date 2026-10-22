// Package vertex configures the standard OpenAI client for Google Vertex AI.
package vertex

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"cloud.google.com/go/auth"
	"cloud.google.com/go/auth/credentials"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
)

// TokenProvider supplies an access token for each request attempt, including retries.
// Implementations may cache valid tokens and must be safe for concurrent calls.
type TokenProvider func(context.Context) (string, error)

// Config selects a Vertex project and location. BaseURL, if supplied, is the
// complete provider base URL up to the final /v1 (before /responses). Only this
// explicitly trusted origin can receive credentials. The OpenAI environment is ignored.
type Config struct {
	// ProjectID defaults to GOOGLE_CLOUD_PROJECT when BaseURL is not supplied.
	ProjectID string
	// Location defaults to GOOGLE_CLOUD_LOCATION when BaseURL is not supplied.
	Location string
	// BaseURL overrides the provider base URL.
	BaseURL string
	// TokenProvider overrides ADC, e.g. to use credentials managed by the application.
	TokenProvider TokenProvider
	// Credentials supplies native Google credentials, including their refresh
	// and quota project behavior. Mutually exclusive with TokenProvider.
	Credentials *auth.Credentials
}

// Legacy project IDs can include a domain prefix; locations also appear in hostnames.
var (
	projectSegment  = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.:-]*$`)
	locationSegment = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]*$`)
)

type credentialError struct{ cause error }

func (e *credentialError) Error() string { return "vertex: unable to resolve Google credentials" }
func (e *credentialError) Unwrap() error { return e.cause }

// NewClient creates a normal OpenAI client using Google Application Default
// Credentials unless Credentials or TokenProvider is supplied. Google token
// refresh and quota lookup receive the request context and its timeout.
func NewClient(ctx context.Context, cfg Config, userOpts ...option.RequestOption) (openai.Client, error) {
	if ctx == nil {
		return openai.Client{}, errors.New("vertex: nil context")
	}
	if err := requestconfig.ValidateEndpointOptions("Vertex", userOpts...); err != nil {
		return openai.Client{}, err
	}
	raw := cfg.BaseURL
	if raw == "" {
		if cfg.ProjectID == "" {
			cfg.ProjectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
		}
		if cfg.Location == "" {
			cfg.Location = os.Getenv("GOOGLE_CLOUD_LOCATION")
		}
		if !projectSegment.MatchString(cfg.ProjectID) || !locationSegment.MatchString(cfg.Location) {
			return openai.Client{}, errors.New("vertex: set ProjectID / GOOGLE_CLOUD_PROJECT and Location / GOOGLE_CLOUD_LOCATION to valid single path segments")
		}
		host := cfg.Location + "-aiplatform.googleapis.com"
		if cfg.Location == "global" {
			host = "aiplatform.googleapis.com"
		}
		raw = fmt.Sprintf("https://%s/v1/projects/%s/locations/%s/publishers/openai/v1/", host, cfg.ProjectID, cfg.Location)
	}
	// Normalize the escaped URL so encoded separators stay inside their path segments.
	base, err := url.Parse(strings.TrimRight(raw, "/") + "/")
	if err != nil || base == nil || base.Host == "" || (base.Scheme != "https" && base.Scheme != "http") ||
		base.User != nil || base.RawQuery != "" || base.Fragment != "" {
		return openai.Client{}, errors.New("vertex: BaseURL must be an absolute HTTP(S) URL without credentials, query or fragment")
	}
	if cfg.Credentials != nil && cfg.TokenProvider != nil {
		return openai.Client{}, errors.New("vertex: choose Credentials or TokenProvider, not both")
	}
	provider := cfg.TokenProvider
	googleCredentials := cfg.Credentials
	if provider == nil && googleCredentials == nil {
		googleCredentials, err = credentials.DetectDefault(&credentials.DetectOptions{
			Scopes:              []string{"https://www.googleapis.com/auth/cloud-platform"},
			DisableAsyncRefresh: true,
		})
		if err != nil {
			return openai.Client{}, &credentialError{cause: err}
		}
	}

	opts := []option.RequestOption{
		requestconfig.WithEnvironmentDefaultsDisabled(),
		requestconfig.WithEndpointProvider("Vertex"),
		option.WithBaseURL(base.String()),
	}
	opts = append(opts, userOpts...)
	opts = append(opts, requestconfig.WithRequestFinalizer(func(rc *requestconfig.RequestConfig) error {
		if rc.APIKey != "" || rc.AdminAPIKey != "" {
			return errors.New("vertex: configure authentication through vertex.Config; OpenAI credentials cannot be combined with Vertex")
		}
		if rc.BaseURL == nil || rc.BaseURL.String() != base.String() {
			return errors.New("vertex: configure routing through vertex.Config.BaseURL")
		}
		if rc.CustomHTTPDoer != nil || rc.HTTPClient == nil {
			return errors.New("vertex: authenticated requests require an *http.Client for redirect safety")
		}
		redirectErr := requestconfig.WithNoRetryError(errors.New("vertex: refusing to follow redirect"))
		client := *rc.HTTPClient
		client.CheckRedirect = func(*http.Request, []*http.Request) error { return redirectErr }
		rc.HTTPClient = &client
		return option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
			if !requestconfig.RequestHasOrigin(req, base) {
				return nil, requestconfig.WithNoRetryError(errors.New("vertex: refusing credentials to an origin other than the configured URL"))
			}
			if req.Header.Get("Authorization") != "" {
				return nil, requestconfig.WithNoRetryError(errors.New("vertex: custom Authorization conflicts with provider authentication"))
			}
			for name := range req.Header {
				if strings.EqualFold(name, "api-key") {
					return nil, requestconfig.WithNoRetryError(errors.New("vertex: custom api-key conflicts with provider authentication"))
				}
			}
			// OpenAI routing metadata is never meaningful for Vertex, including
			// values installed by method-level options and middleware.
			req.Header.Del("OpenAI-Organization")
			req.Header.Del("OpenAI-Project")
			var token string
			if googleCredentials != nil {
				access, authErr := googleCredentials.Token(req.Context())
				if authErr != nil {
					return nil, &credentialError{cause: authErr}
				}
				quota, quotaErr := googleCredentials.QuotaProjectID(req.Context())
				if quotaErr != nil {
					return nil, &credentialError{cause: quotaErr}
				}
				if quota != "" {
					req.Header.Set("X-Goog-User-Project", quota)
				}
				token = access.Value
			} else {
				var authErr error
				token, authErr = provider(req.Context())
				if authErr != nil {
					return nil, &credentialError{cause: authErr}
				}
			}
			token = strings.TrimSpace(token)
			if token == "" {
				return nil, requestconfig.WithNoRetryError(errors.New("vertex: token provider must return a non-empty token"))
			}
			req.Header.Set("Authorization", "Bearer "+token)
			response, err := next(req)
			if errors.Is(err, redirectErr) {
				// net/http closes the response body; discard its URL-bearing error wrapper.
				return nil, redirectErr
			}
			return response, err
		}).Apply(rc)
	}))
	return openai.NewClient(opts...), nil
}

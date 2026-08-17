package option

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"sync"

	"github.com/openai/openai-go/v3/auth"
	"github.com/openai/openai-go/v3/internal"
	"github.com/openai/openai-go/v3/internal/requestconfig"
)

const x509APIBaseURL = "https://mtls.api.openai.com/v1/"

// Keep enough transport-scoped auth states for a primary client plus a small
// set of method-level overrides without retaining every transport ever used by
// a long-lived client.
const x509WorkloadIdentityAuthCacheCapacity = 8

var (
	errX509RequestCredentialConflict = errors.New("X.509 workload identity cannot be combined with another request credential")
	errX509AuthorizationProvenance   = errors.New("X.509 workload identity request Authorization does not match SDK-selected authentication")
)

// WithX509WorkloadIdentity returns a RequestOption that configures X.509
// workload identity authentication. The configured HTTP client must present
// the client certificate for token exchange and API requests. Custom HTTPClient
// implementations are responsible for refusing redirects internally.
//
// When no base URL is configured explicitly or through OPENAI_BASE_URL, X.509
// workload identity clients use https://mtls.api.openai.com/v1. Explicit API
// bases must be absolute HTTPS URLs, and requests remain bound to that origin.
// Provider authentication options such as Azure and Bedrock are incompatible.
func WithX509WorkloadIdentity(config auth.X509WorkloadIdentity) RequestOption {
	cache := x509WorkloadIdentityAuthCache{}
	return requestconfig.RequestOptionFunc(func(r *requestconfig.RequestConfig) error {
		if err := r.UseWorkloadIdentityCredential(requestconfig.WorkloadIdentityCredentialSourceX509); err != nil {
			return err
		}
		r.SetAPIKey("")
		r.SetWorkloadIdentityFinalizer(func(r *requestconfig.RequestConfig) error {
			return configureX509Request(r, &cache, config)
		})
		return nil
	})
}

type x509WorkloadIdentityAuthCache struct {
	mu      sync.Mutex
	entries []x509WorkloadIdentityAuthCacheEntry
}

type x509WorkloadIdentityAuthCacheEntry struct {
	httpDoer      auth.HTTPDoer
	httpTransport http.RoundTripper
	auth          *auth.WorkloadIdentityAuth
}

func (c *x509WorkloadIdentityAuthCache) get(
	httpDoer auth.HTTPDoer,
	config auth.X509WorkloadIdentity,
) (*auth.WorkloadIdentityAuth, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	httpTransport := internal.NativeHTTPTransport(httpDoer)
	for i := range c.entries {
		if !sameX509CacheIdentity(c.entries[i].httpDoer, httpDoer) ||
			!sameX509CacheIdentity(c.entries[i].httpTransport, httpTransport) {
			continue
		}
		entry := c.entries[i]
		copy(c.entries[i:], c.entries[i+1:])
		c.entries[len(c.entries)-1] = entry
		return entry.auth, nil
	}
	wia, err := auth.NewX509WorkloadIdentityAuth(config)
	if err != nil {
		return nil, err
	}
	entry := x509WorkloadIdentityAuthCacheEntry{
		httpDoer:      httpDoer,
		httpTransport: httpTransport,
		auth:          wia,
	}
	if len(c.entries) == x509WorkloadIdentityAuthCacheCapacity {
		copy(c.entries, c.entries[1:])
		c.entries[len(c.entries)-1] = entry
	} else {
		c.entries = append(c.entries, entry)
	}
	return wia, nil
}

func sameX509CacheIdentity(left any, right any) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}
	leftValue := reflect.ValueOf(left)
	rightValue := reflect.ValueOf(right)
	if leftValue.Type() != rightValue.Type() || !leftValue.Comparable() || !rightValue.Comparable() {
		return false
	}
	return leftValue.Interface() == rightValue.Interface()
}

func configureX509Request(
	r *requestconfig.RequestConfig,
	cache *x509WorkloadIdentityAuthCache,
	config auth.X509WorkloadIdentity,
) error {
	if provider := r.ProviderAuthentication(); provider != "" {
		return fmt.Errorf("X.509 workload identity cannot be combined with %s provider authentication", provider)
	}
	if r.APIKey != "" || hasConflictingX509CredentialHeaders(r.Request.Header) ||
		!hasExactX509Authorization(r.Request.Header, "") {
		return errX509RequestCredentialConflict
	}

	var configuredBaseURL *url.URL
	if r.BaseURL == nil {
		if err := r.Apply(requestconfig.WithDefaultBaseURL(x509APIBaseURL)); err != nil {
			return err
		}
		configuredBaseURL = r.DefaultBaseURL
	} else {
		configuredBaseURL = r.BaseURL
	}
	apiOrigin, err := newX509APIOrigin(configuredBaseURL)
	if err != nil {
		return err
	}

	var configuredHTTPDoer auth.HTTPDoer = r.HTTPClient
	if r.CustomHTTPDoer != nil {
		configuredHTTPDoer = r.CustomHTTPDoer
	}
	if configuredHTTPDoer == nil {
		return errors.New("X.509 workload identity requires an HTTP client")
	}
	if r.CustomHTTPDoer == nil && r.HTTPClient != nil {
		clone := *r.HTTPClient
		clone.CheckRedirect = func(*http.Request, []*http.Request) error {
			return errors.New("X.509 workload identity API redirects are disabled")
		}
		r.HTTPClient = &clone
	}

	var authorizationMiddleware Middleware
	if r.Security.BearerAuth {
		wia, err := cache.get(configuredHTTPDoer, config)
		if err != nil {
			return err
		}
		authorizationMiddleware = func(req *http.Request, next MiddlewareNext) (*http.Response, error) {
			if err := apiOrigin.validateRequest(req); err != nil {
				return nil, requestconfig.WithNoRetryError(err)
			}
			return auth.WorkloadIdentityMiddleware(wia, configuredHTTPDoer, internal.WithX509RequestPolicy(req), next)
		}
	} else {
		expectedAuthorization := ""
		if r.Security.AdminAPIKeyAuth && r.AdminAPIKey != "" {
			expectedAuthorization = "Bearer " + r.AdminAPIKey
		}
		authorizationMiddleware = func(req *http.Request, next MiddlewareNext) (*http.Response, error) {
			if err := apiOrigin.validateRequest(req); err != nil {
				return nil, requestconfig.WithNoRetryError(err)
			}
			if !hasExactX509Authorization(req.Header, expectedAuthorization) {
				return nil, requestconfig.WithNoRetryError(errX509AuthorizationProvenance)
			}
			return next(internal.WithExpectedAuthorization(req, expectedAuthorization))
		}
	}
	r.Middlewares = append([]Middleware{authorizationMiddleware}, r.Middlewares...)
	r.Middlewares = append(r.Middlewares, func(req *http.Request, next MiddlewareNext) (*http.Response, error) {
		if err := apiOrigin.validateRequest(req); err != nil {
			return nil, requestconfig.WithNoRetryError(err)
		}
		if hasConflictingX509CredentialHeaders(req.Header) {
			return nil, requestconfig.WithNoRetryError(errX509RequestCredentialConflict)
		}
		expectedAuthorization, ok := internal.ExpectedAuthorization(req)
		if !ok {
			return nil, requestconfig.WithNoRetryError(errors.New("X.509 workload identity authorization state is missing"))
		}
		if !hasExactX509Authorization(req.Header, expectedAuthorization) {
			return nil, requestconfig.WithNoRetryError(errX509AuthorizationProvenance)
		}
		return next(req)
	})
	return nil
}

type x509APIOrigin struct {
	hostname string
	port     string
}

func newX509APIOrigin(value *url.URL) (x509APIOrigin, error) {
	if value == nil || !strings.EqualFold(value.Scheme, "https") || value.Opaque != "" ||
		value.User != nil || value.Host == "" || strings.ContainsAny(value.Host, `\%`) {
		return x509APIOrigin{}, errors.New("X.509 workload identity requires an absolute HTTPS API base URL without userinfo")
	}
	hostname := value.Hostname()
	if hostname == "" {
		return x509APIOrigin{}, errors.New("X.509 workload identity requires an absolute HTTPS API base URL without userinfo")
	}
	port := value.Port()
	if port == "" {
		port = "443"
	}
	return x509APIOrigin{hostname: strings.ToLower(hostname), port: port}, nil
}

func (origin x509APIOrigin) validate(value *url.URL) error {
	candidate, err := newX509APIOrigin(value)
	if err != nil {
		return err
	}
	if candidate != origin {
		return errors.New("X.509 workload identity cannot send credentials to an origin other than the configured API URL")
	}
	return nil
}

func (origin x509APIOrigin) validateRequest(req *http.Request) error {
	if err := origin.validate(req.URL); err != nil {
		return err
	}
	if req.Host == "" {
		return nil
	}
	if strings.ContainsAny(req.Host, `@\/%?#`) || strings.HasSuffix(req.Host, ":") {
		return errors.New("X.509 workload identity requires the HTTP request authority to match the configured API origin")
	}
	authority := &url.URL{Scheme: "https", Host: req.Host}
	if err := origin.validate(authority); err != nil {
		return errors.New("X.509 workload identity requires the HTTP request authority to match the configured API origin")
	}
	return nil
}

func hasConflictingX509CredentialHeaders(header http.Header) bool {
	for name, values := range header {
		normalizedName := strings.ToLower(strings.ReplaceAll(name, "_", "-"))
		if normalizedName != "api-key" && normalizedName != "x-api-key" && normalizedName != "proxy-authorization" {
			continue
		}
		for _, value := range values {
			if value != "" {
				return true
			}
		}
	}
	return false
}

func hasExactX509Authorization(header http.Header, expected string) bool {
	var values []string
	for name, headerValues := range header {
		if strings.EqualFold(name, "Authorization") {
			values = append(values, headerValues...)
		}
	}
	if expected == "" {
		return len(values) == 0
	}
	return len(values) == 1 && values[0] == expected
}

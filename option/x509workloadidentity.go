package option

import (
	"errors"
	"net/http"
	"net/url"
	"path"
	"strings"
	"sync"

	"github.com/openai/openai-go/v3/auth"
	"github.com/openai/openai-go/v3/internal/requestconfig"
)

const x509WorkloadAPIBaseURL = "https://mtls.api.openai.com/v1/"

// WithX509WorkloadIdentity authenticates OpenAI API requests by exchanging an
// attested workload certificate for short-lived bearer credentials. Only the
// global OpenAI mTLS endpoint and the configured transport are supported.
func WithX509WorkloadIdentity(config auth.X509WorkloadIdentity) RequestOption {
	var initialize sync.Once
	var identity *auth.X509WorkloadIdentityAuth
	var initializationError error

	return requestconfig.RequestOptionFunc(func(cfg *requestconfig.RequestConfig) error {
		selected := requestconfig.NewProviderAuthOption("OpenAI", "option.WithX509WorkloadIdentity")
		if err := selected.Apply(cfg); err != nil {
			return err
		}
		cfg.ClearInheritedAuthentication()
		return requestconfig.WithRequestFinalizer(func(final *requestconfig.RequestConfig) error {
			if !selected.Selected(final) {
				return nil
			}
			if final.HTTPClient == nil || final.CustomHTTPDoer != nil ||
				final.HTTPClientExplicitlySelected() {
				return errors.New("X.509 workload identity requires its attested transport without custom HTTP clients")
			}
			if final.EndpointProvider() != "" {
				return errors.New("X.509 workload identity cannot be combined with another API provider")
			}
			if !final.Security.BearerAuth {
				return errors.New("X.509 workload identity cannot authenticate an admin-only API operation")
			}
			if final.Request == nil || final.Request.Header == nil {
				return errors.New("X.509 workload identity requires a non-nil request and header map")
			}
			if final.APIKey != "" || final.AdminAPIKey != "" || final.AuthorizationHeaderOverridden() ||
				unsafeX509CredentialHeaders(final.Request.Header) {
				return errors.New("X.509 workload identity cannot be combined with other credentials")
			}
			if final.BaseURL == nil {
				if err := requestconfig.WithDefaultBaseURL(x509WorkloadAPIBaseURL).Apply(final); err != nil {
					return err
				}
			}
			base := final.BaseURL
			if base == nil {
				base = final.DefaultBaseURL
			}
			if !validX509WorkloadAPIBaseURL(base) {
				return errors.New("X.509 workload identity requires the global OpenAI mTLS API endpoint")
			}
			initialize.Do(func() {
				identity, initializationError = auth.NewX509WorkloadIdentityAuth(config)
			})
			if initializationError != nil {
				return initializationError
			}
			final.CustomHTTPDoer = config.Transport
			allowBodyReplay := len(final.Middlewares) == 0
			final.InstallRequestRetryScope(allowBodyReplay)
			final.InstallRequestAttemptMiddleware()
			return WithMiddleware(func(request *http.Request, next MiddlewareNext) (*http.Response, error) {
				if !validX509WorkloadAPIRequest(request) || unsafeX509CredentialHeaders(request.Header) {
					return nil, requestconfig.WithNoRetryError(
						errors.New("X.509 workload identity rejected an unsafe final API request"),
					)
				}
				if requestconfig.RequestRetryScopeFromContext(request.Context()) == nil {
					return nil, requestconfig.WithNoRetryError(
						errors.New("X.509 workload identity requires its request-owned retry scope"),
					)
				}
				authenticationFailed := true
				response, err := auth.X509WorkloadIdentityMiddleware(identity, config.Transport, request,
					func(authenticated *http.Request) (*http.Response, error) {
						response, dispatchErr := next(authenticated)
						authenticationFailed = dispatchErr == nil && response != nil &&
							response.StatusCode == http.StatusUnauthorized
						return response, dispatchErr
					})
				if err != nil && authenticationFailed {
					return nil, requestconfig.WithNoRetryError(err)
				}
				return redactX509Response(response), err
			}).Apply(final)
		}).Apply(cfg)
	})
}

func redactX509Response(response *http.Response) *http.Response {
	if response == nil || response.Request == nil {
		return response
	}
	redacted := *response
	redacted.Request = response.Request.Clone(response.Request.Context())
	for name := range redacted.Request.Header {
		switch strings.ToLower(strings.ReplaceAll(name, "_", "-")) {
		case "authorization", "proxy-authorization", "cookie", "set-cookie", "api-key", "x-api-key", "x-amz-security-token":
			delete(redacted.Request.Header, name)
		}
	}
	return &redacted
}

func validX509WorkloadAPIBaseURL(base *url.URL) bool {
	return base != nil && base.Scheme == "https" && validX509WorkloadAPIAuthority(base) &&
		base.Path == "/v1/" && base.EscapedPath() == "/v1/" && base.User == nil &&
		base.RawQuery == "" && !base.ForceQuery && base.Fragment == "" && base.RawFragment == "" && base.Opaque == ""
}

func validX509WorkloadAPIRequest(request *http.Request) bool {
	if request == nil || request.Header == nil || request.URL == nil || request.URL.Scheme != "https" ||
		!validX509WorkloadAPIAuthority(request.URL) || request.URL.User != nil || request.RequestURI != "" ||
		request.URL.Fragment != "" || request.URL.RawFragment != "" || request.URL.Opaque != "" ||
		(request.Host != "" && request.Host != request.URL.Host) || len(request.Trailer) != 0 ||
		len(request.TransferEncoding) != 0 {
		return false
	}
	//nolint:staticcheck // Legacy cancellation bypasses the issuer context and has no nondeprecated accessor.
	if request.Cancel != nil {
		return false
	}
	for index := 0; index < len(request.URL.RawQuery); index++ {
		value := request.URL.RawQuery[index]
		if value < ' ' || value == 0x7f {
			return false
		}
	}
	switch request.Method {
	case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch,
		http.MethodDelete, http.MethodHead, http.MethodOptions:
	default:
		return false
	}
	hasBody := request.Body != nil && request.Body != http.NoBody
	if request.ContentLength < -1 || (!hasBody && request.ContentLength != 0) {
		return false
	}
	for name, values := range request.Header {
		if !requestconfig.ValidHTTPHeaderName(name) {
			return false
		}
		for _, value := range values {
			if !requestconfig.ValidHTTPHeaderValue(value) {
				return false
			}
		}
		switch strings.ToLower(strings.ReplaceAll(name, "_", "-")) {
		case "transfer-encoding", "content-length", "connection", "upgrade", "trailer", "te",
			"proxy-connection", "keep-alive", "http2-settings":
			return false
		}
	}
	return strings.HasPrefix(request.URL.Path, "/v1/") && strings.HasPrefix(request.URL.EscapedPath(), "/v1/") &&
		path.Clean(request.URL.Path) == request.URL.Path
}

func validX509WorkloadAPIAuthority(target *url.URL) bool {
	if target == nil || target.Hostname() != "mtls.api.openai.com" {
		return false
	}
	port := target.Port()
	return (port == "" || port == "443") &&
		(target.Host == "mtls.api.openai.com" || target.Host == "mtls.api.openai.com:443")
}

func unsafeX509CredentialHeaders(headers http.Header) bool {
	for name := range headers {
		normalized := strings.ToLower(strings.ReplaceAll(name, "_", "-"))
		switch normalized {
		case "authorization", "api-key", "x-api-key", "proxy-authorization", "cookie", "set-cookie", "host", "x-amz-security-token":
			return true
		}
		if strings.HasPrefix(normalized, ":") {
			return true
		}
	}
	return false
}

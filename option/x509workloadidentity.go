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
		if err := requestconfig.WithDefaultBaseURL(x509WorkloadAPIBaseURL).Apply(cfg); err != nil {
			return err
		}
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
			return WithMiddleware(func(request *http.Request, next MiddlewareNext) (*http.Response, error) {
				if !validX509WorkloadAPIRequest(request) || unsafeX509CredentialHeaders(request.Header) {
					return nil, requestconfig.WithNoRetryError(
						errors.New("X.509 workload identity rejected an unsafe final API request"),
					)
				}
				token, err := identity.GetToken(request.Context(), config.Transport)
				if err != nil {
					return nil, requestconfig.WithNoRetryError(err)
				}
				authenticated := request.Clone(request.Context())
				authenticated.Header.Set("Authorization", "Bearer "+token)
				return next(authenticated)
			}).Apply(final)
		}).Apply(cfg)
	})
}

func validX509WorkloadAPIBaseURL(base *url.URL) bool {
	return base != nil && base.Scheme == "https" && base.Host == "mtls.api.openai.com" &&
		base.Path == "/v1/" && base.EscapedPath() == "/v1/" && base.User == nil &&
		base.RawQuery == "" && !base.ForceQuery && base.Fragment == "" && base.RawFragment == "" && base.Opaque == ""
}

func validX509WorkloadAPIRequest(request *http.Request) bool {
	if request == nil || request.Header == nil || request.URL == nil || request.URL.Scheme != "https" ||
		request.URL.Host != "mtls.api.openai.com" || request.URL.User != nil ||
		request.URL.Fragment != "" || request.URL.RawFragment != "" || request.URL.Opaque != "" ||
		(request.Host != "" && request.Host != request.URL.Host) || len(request.Trailer) != 0 ||
		len(request.TransferEncoding) != 0 {
		return false
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
	for name := range request.Header {
		switch strings.ToLower(strings.ReplaceAll(name, "_", "-")) {
		case "transfer-encoding", "content-length", "connection", "upgrade", "trailer", "te",
			"proxy-connection", "keep-alive", "http2-settings":
			return false
		}
	}
	return strings.HasPrefix(request.URL.Path, "/v1/") && strings.HasPrefix(request.URL.EscapedPath(), "/v1/") &&
		path.Clean(request.URL.Path) == request.URL.Path
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

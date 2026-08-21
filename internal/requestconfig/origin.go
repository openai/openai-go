package requestconfig

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
)

func validateRequestReference(value string) error {
	if strings.HasPrefix(value, "//") {
		return errors.New("requestconfig: request path must be a relative URL reference")
	}
	if err := validateRelativeRequestReference(value); err != nil {
		return err
	}
	normalized := strings.TrimPrefix(value, "/")
	if normalized == value {
		return nil
	}
	return validateRelativeRequestReference(normalized)
}

func validateRelativeRequestReference(value string) error {
	reference, err := url.Parse(value)
	if err != nil {
		return err
	}
	if reference.IsAbs() || reference.Host != "" {
		return errors.New("requestconfig: request path must be a relative URL reference")
	}
	return nil
}

// SameOrigin reports whether two URLs have the same scheme, host, and effective
// port. It is internal API shared by the generic request pipeline and provider
// authentication middleware.
func SameOrigin(left, right *url.URL) bool {
	if left == nil || right == nil || !strings.EqualFold(left.Scheme, right.Scheme) || !strings.EqualFold(left.Hostname(), right.Hostname()) {
		return false
	}
	return effectivePort(left) == effectivePort(right)
}

func effectivePort(value *url.URL) string {
	if port := value.Port(); port != "" {
		return port
	}
	if strings.EqualFold(value.Scheme, "https") {
		return "443"
	}
	if strings.EqualFold(value.Scheme, "http") {
		return "80"
	}
	return ""
}

type originTransport struct {
	origin *url.URL
	next   http.RoundTripper
}

func (t originTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req == nil || !SameOrigin(req.URL, t.origin) {
		return rejectRequestOrigin(req)
	}
	return t.next.RoundTrip(req)
}

func rejectRequestOrigin(req *http.Request) (*http.Response, error) {
	if req != nil && req.Body != nil {
		if _, ok := req.Body.(*closeOnceReadCloser); !ok {
			req.Body = &closeOnceReadCloser{ReadCloser: req.Body}
		}
		_ = req.Body.Close()
	}
	return nil, requestOriginError()
}

func requestOriginError() error {
	return WithNoRetryError(errors.New("requestconfig: request URL origin must match the configured base URL"))
}

func enforceRequestOrigin(origin *url.URL, next middlewareNext) middlewareNext {
	return func(req *http.Request) (*http.Response, error) {
		if req == nil || !SameOrigin(req.URL, origin) {
			return rejectRequestOrigin(req)
		}
		return next(req)
	}
}

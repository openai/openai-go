package internal

import (
	"context"
	"net/http"
)

type expectedAuthorizationContextKey struct{}
type x509RequestPolicyContextKey struct{}

// WithX509RequestPolicy marks a request whose X.509 destination and credential
// policy has been validated by the high-level client option.
func WithX509RequestPolicy(req *http.Request) *http.Request {
	ctx := context.WithValue(req.Context(), x509RequestPolicyContextKey{}, true)
	return req.WithContext(ctx)
}

// HasX509RequestPolicy reports whether the high-level X.509 request policy is
// bound to this request attempt.
func HasX509RequestPolicy(req *http.Request) bool {
	policyBound, _ := req.Context().Value(x509RequestPolicyContextKey{}).(bool)
	return policyBound
}

// WithExpectedAuthorization binds the SDK-selected Authorization value to one
// request attempt without exposing it through mutable middleware state.
func WithExpectedAuthorization(req *http.Request, authorization string) *http.Request {
	ctx := context.WithValue(req.Context(), expectedAuthorizationContextKey{}, authorization)
	return req.WithContext(ctx)
}

// ExpectedAuthorization returns the SDK-selected Authorization value bound to
// the current request attempt. An empty value is valid for headerless requests.
func ExpectedAuthorization(req *http.Request) (string, bool) {
	authorization, ok := req.Context().Value(expectedAuthorizationContextKey{}).(string)
	return authorization, ok
}

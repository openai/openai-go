package internal

import (
	"context"
	"net/http"
)

type expectedAuthorizationContextKey struct{}

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

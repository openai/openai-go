package auth

import (
	"fmt"

	"github.com/openai/openai-go/v3/shared"
)

// SubjectTokenProviderError is raised when failing to get the subject token from the cloud environment.
type SubjectTokenProviderError struct {
	Provider string
	Message  string
	Cause    error
}

func (e *SubjectTokenProviderError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s provider error: %s: %v", e.Provider, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s provider error: %s", e.Provider, e.Message)
}

func (e *SubjectTokenProviderError) Unwrap() error {
	return e.Cause
}

// OAuthError is raised when there is an error in the RFC-defined OAuth flow, such as failing to exchange the token.
// See https://datatracker.ietf.org/doc/html/rfc8693 for the OAuth 2.0 Token Exchange specification.
type OAuthError struct {
	StatusCode       int
	ErrorCode        shared.OAuthErrorCode
	ErrorDescription string
}

func (e *OAuthError) Error() string {
	return fmt.Sprintf("OAuth error (status %d): %s - %s", e.StatusCode, e.ErrorCode, e.ErrorDescription)
}

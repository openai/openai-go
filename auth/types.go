package auth

import (
	"context"
	"net/http"
)

type HTTPDoer interface {
	Do(*http.Request) (*http.Response, error)
}

type SubjectTokenType string

const (
	SubjectTokenTypeJWT SubjectTokenType = "jwt"
	SubjectTokenTypeID  SubjectTokenType = "id"
)

type SubjectTokenProvider interface {
	TokenType() SubjectTokenType
	GetToken(ctx context.Context, httpClient HTTPDoer) (string, error)
}

type WorkloadIdentity struct {
	ClientID             string
	IdentityProviderID   string
	ServiceAccountID     string
	Provider             SubjectTokenProvider
	RefreshBufferSeconds int
}

type TokenExchangeResponse struct {
	AccessToken     string `json:"access_token"`
	IssuedTokenType string `json:"issued_token_type"`
	TokenType       string `json:"token_type"`
	ExpiresIn       *int   `json:"expires_in,omitempty"`
}

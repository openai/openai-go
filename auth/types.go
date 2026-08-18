package auth

import (
	"context"
	"net/http"
	"time"
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

// X509WorkloadIdentity configures X.509 workload identity federation. Client
// certificate and TLS configuration remain the responsibility of the HTTP
// transport supplied to the OpenAI client. The client must use a native
// *http.Transport with exactly one static client certificate and no client TLS
// session cache.
type X509WorkloadIdentity struct {
	IdentityProviderID string
	ServiceAccountID   string
	RefreshBuffer      time.Duration
}

type TokenExchangeResponse struct {
	AccessToken     string `json:"access_token"`
	IssuedTokenType string `json:"issued_token_type"`
	TokenType       string `json:"token_type"`
	ExpiresIn       *int   `json:"expires_in,omitempty"`
}

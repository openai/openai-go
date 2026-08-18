package auth

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"
)

const (
	x509SubjectTokenType          = "urn:openai:params:oauth:token-type:x509"
	x509TokenExchangeURL          = "https://mtls.auth.openai.com/oauth/token"
	tokenExchangeMaxRetries       = 2
	tokenExchangeResponseBodySize = 1 << 20
	// Allow every accepted Retry-After delay plus a bounded aggregate budget
	// for the initial request and retries.
	x509TokenExchangeRequestBudget     = 30 * time.Second
	x509WorkloadIdentityRefreshTimeout = tokenExchangeMaxRetryDelay*tokenExchangeMaxRetries + x509TokenExchangeRequestBudget
)

// x509TokenExchangeRequest deliberately has no SubjectToken or ClientID field.
// The mutually authenticated TLS connection is the subject credential.
type x509TokenExchangeRequest struct {
	GrantType          string `json:"grant_type"`
	SubjectTokenType   string `json:"subject_token_type"`
	IdentityProviderID string `json:"identity_provider_id"`
	ServiceAccountID   string `json:"service_account_id"`
}

type x509CredentialSource struct {
	refreshBefore time.Duration
}

var errX509NativeHTTPTransport = errors.New("X.509 workload identity requires an *http.Client with a native *http.Transport so its client identity can be verified")

// x509TokenExchangeError prevents the generic API request loop from repeating
// an exchange that has already exhausted the X.509 exchange retry policy.
type x509TokenExchangeError struct {
	err error
}

func (e *x509TokenExchangeError) Error() string { return e.err.Error() }
func (e *x509TokenExchangeError) Unwrap() error { return e.err }
func (e *x509TokenExchangeError) NoRetry()      {}

func noRetryX509TokenExchangeError(err error) error {
	if err == nil {
		return nil
	}
	return &x509TokenExchangeError{err: err}
}

// NewX509WorkloadIdentityAuth creates the HTTP authentication state for an
// X.509 workload identity. Token exchange remains lazy until GetToken is called.
// The state binds to the first *http.Client used and rejects later client or
// transport changes so cached tokens cannot cross certificate-backed
// transports. Its native *http.Transport must use exactly one static client
// certificate with client session resumption disabled; replace the transport
// rather than mutating it when rotating the identity.
func NewX509WorkloadIdentityAuth(config X509WorkloadIdentity) (*WorkloadIdentityAuth, error) {
	if err := validateWorkloadIdentity("X509WorkloadIdentity", config.IdentityProviderID, config.ServiceAccountID); err != nil {
		return nil, err
	}
	if config.RefreshBuffer < 0 {
		return nil, fmt.Errorf("X509WorkloadIdentity: RefreshBuffer must be non-negative")
	}
	return newWorkloadIdentityAuth(
		config.IdentityProviderID,
		config.ServiceAccountID,
		x509CredentialSource{refreshBefore: config.RefreshBuffer},
	), nil
}

func x509HTTPTransportIdentity(transport *http.Transport) ([sha256.Size]byte, error) {
	legacyTLSDialConfigured := transport.DialTLS != nil //nolint:staticcheck // Reject a legacy hook that can select another client identity.
	tlsConfig := transport.TLSClientConfig
	if tlsConfig != nil && len(tlsConfig.Certificates) == 1 &&
		len(tlsConfig.Certificates[0].Certificate) != 0 &&
		tlsConfig.GetClientCertificate == nil &&
		tlsConfig.ClientSessionCache == nil &&
		!legacyTLSDialConfigured && transport.DialTLSContext == nil {
		return sha256.Sum256(tlsConfig.Certificates[0].Certificate[0]), nil
	}
	return [sha256.Size]byte{}, errors.New("X.509 workload identity requires one immutable client identity per native HTTP transport; configure exactly one static TLS certificate without TLS dial hooks, certificate-selection hooks, or client session caching, and replace the transport to rotate identities")
}

func (s x509CredentialSource) exchange(
	ctx context.Context,
	httpClient HTTPDoer,
	identityProviderID string,
	serviceAccountID string,
) (exchangedToken, error) {
	body, err := json.Marshal(x509TokenExchangeRequest{
		GrantType:          TokenExchangeGrantType,
		SubjectTokenType:   x509SubjectTokenType,
		IdentityProviderID: identityProviderID,
		ServiceAccountID:   serviceAccountID,
	})
	if err != nil {
		return exchangedToken{}, noRetryX509TokenExchangeError(fmt.Errorf("failed to marshal token exchange request: %w", err))
	}
	httpClient = withoutRedirects(httpClient)
	for retry := 0; ; retry++ {
		resp, exchangeErr := exchangeToken(ctx, httpClient, x509TokenExchangeURL, body, 0)
		if exchangeErr == nil {
			responseBody, readErr := readTokenExchangeResponse(resp.Body, tokenExchangeResponseBodySize)
			_ = resp.Body.Close()
			if readErr == nil {
				token, parseErr := parseX509TokenExchangeResponse(resp.StatusCode, responseBody)
				if parseErr == nil || !shouldRetryTokenExchange(resp, nil) || retry >= tokenExchangeMaxRetries {
					return token, noRetryX509TokenExchangeError(parseErr)
				}
			} else if retry >= tokenExchangeMaxRetries {
				return exchangedToken{}, noRetryX509TokenExchangeError(readErr)
			}
		} else if retry >= tokenExchangeMaxRetries {
			return exchangedToken{}, noRetryX509TokenExchangeError(exchangeErr)
		}

		if err := waitForTokenExchangeRetry(ctx, tokenExchangeRetryDelay(resp, retry)); err != nil {
			return exchangedToken{}, noRetryX509TokenExchangeError(err)
		}
	}
}

func (s x509CredentialSource) refreshBuffer(expiresIn time.Duration) time.Duration {
	refreshBefore := s.refreshBefore
	if refreshBefore == 0 {
		refreshBefore = DefaultRefreshBuffer
	}
	return min(refreshBefore, expiresIn/2)
}

func (x509CredentialSource) refreshTimeout() time.Duration {
	return x509WorkloadIdentityRefreshTimeout
}

func (x509CredentialSource) kind() workloadIdentityCredentialSourceKind {
	return workloadIdentityCredentialSourceX509
}

func parseX509TokenExchangeResponse(statusCode int, body []byte) (exchangedToken, error) {
	if statusCode == http.StatusBadRequest || statusCode == http.StatusUnauthorized || statusCode == http.StatusForbidden {
		// Preserve the typed OAuth status without retaining or surfacing an
		// untrusted response body that could contain credential material.
		return exchangedToken{}, &OAuthError{StatusCode: statusCode}
	}
	if statusCode != http.StatusOK {
		return exchangedToken{}, fmt.Errorf("token exchange failed with status %d", statusCode)
	}

	var tokenResp struct {
		AccessToken string          `json:"access_token"`
		ExpiresIn   *float64        `json:"expires_in"`
		TokenType   json.RawMessage `json:"token_type"`
	}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return exchangedToken{}, fmt.Errorf("failed to decode token exchange response: %w", err)
	}
	if tokenResp.ExpiresIn == nil {
		return exchangedToken{}, fmt.Errorf("token exchange response requires a positive numeric 'expires_in' field")
	}
	seconds := *tokenResp.ExpiresIn
	if math.IsNaN(seconds) || math.IsInf(seconds, 0) ||
		seconds > float64(math.MaxInt64)/float64(time.Second) {
		return exchangedToken{}, fmt.Errorf("token exchange response has invalid 'expires_in' field")
	}
	expiresIn := time.Duration(seconds * float64(time.Second))
	if expiresIn <= 0 {
		return exchangedToken{}, fmt.Errorf("token exchange response requires a positive numeric 'expires_in' field")
	}
	if len(tokenResp.TokenType) != 0 {
		var tokenType string
		if tokenResp.TokenType[0] != '"' || json.Unmarshal(tokenResp.TokenType, &tokenType) != nil ||
			!strings.EqualFold(tokenType, "Bearer") {
			return exchangedToken{}, fmt.Errorf("token exchange response has invalid 'token_type' field")
		}
	}
	if !validBearerAccessToken(tokenResp.AccessToken) {
		return exchangedToken{}, fmt.Errorf("token exchange response has invalid 'access_token' field")
	}
	return exchangedToken{accessToken: tokenResp.AccessToken, expiresIn: expiresIn}, nil
}

func validBearerAccessToken(token string) bool {
	if token == "" {
		return false
	}
	padding := false
	for i := range len(token) {
		character := token[i]
		switch {
		case character == '=':
			padding = true
		case character >= 'a' && character <= 'z',
			character >= 'A' && character <= 'Z',
			character >= '0' && character <= '9',
			character == '-', character == '.', character == '_', character == '~',
			character == '+', character == '/':
			if padding {
				return false
			}
		default:
			return false
		}
	}
	return token[0] != '='
}

func withoutRedirects(httpClient HTTPDoer) HTTPDoer {
	client, ok := httpClient.(*http.Client)
	if !ok {
		return httpClient
	}
	clone := *client
	clone.Jar = nil
	clone.CheckRedirect = func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	}
	return &clone
}

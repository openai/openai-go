package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/openai/openai-go/v3/shared"
)

const (
	TokenExchangeGrantType = "urn:ietf:params:oauth:grant-type:token-exchange"
	JWTTokenType           = "urn:ietf:params:oauth:token-type:jwt"
	IDTokenType            = "urn:ietf:params:oauth:token-type:id_token"
	DefaultTokenExpiry     = 60 * time.Minute
	DefaultRefreshBuffer   = 20 * time.Minute
	TokenExchangeURL       = "https://auth.openai.com/oauth/token"
)

type WorkloadIdentityAuth struct {
	identityProviderID string
	serviceAccountID   string
	source             workloadIdentityCredentialSource

	// Protects cachedToken, tokenExpiry, tokenRefreshAt, and refreshInFlight.
	mu              sync.Mutex
	cachedToken     string
	tokenExpiry     time.Time
	tokenRefreshAt  time.Time
	refreshInFlight *tokenRefreshState
}

type tokenRefreshResult struct {
	token     string
	expiresAt time.Time
	refreshAt time.Time
	err       error
}

// Coordinates concurrent access to a single in-flight refresh operation
// done channel signals completion to all waiting goroutines
type tokenRefreshState struct {
	done   chan struct{}
	result tokenRefreshResult
}

type subjectTokenExchangeRequest struct {
	GrantType          string `json:"grant_type"`
	ClientID           string `json:"client_id,omitempty"`
	SubjectToken       string `json:"subject_token"`
	SubjectTokenType   string `json:"subject_token_type"`
	IdentityProviderID string `json:"identity_provider_id"`
	ServiceAccountID   string `json:"service_account_id"`
}

type exchangedToken struct {
	accessToken string
	expiresIn   time.Duration
}

type workloadIdentityCredentialSource interface {
	exchange(context.Context, HTTPDoer, string, string) (exchangedToken, error)
	refreshBuffer(time.Duration) time.Duration
}

type subjectTokenCredentialSource struct {
	clientID      string
	provider      SubjectTokenProvider
	refreshBefore time.Duration
}

func (s subjectTokenCredentialSource) exchange(
	ctx context.Context,
	httpClient HTTPDoer,
	identityProviderID string,
	serviceAccountID string,
) (exchangedToken, error) {
	subjectToken, err := s.provider.GetToken(ctx, httpClient)
	if err != nil {
		return exchangedToken{}, err
	}

	providerTokenType := s.provider.TokenType()
	var subjectTokenType string
	switch providerTokenType {
	case SubjectTokenTypeJWT:
		subjectTokenType = JWTTokenType
	case SubjectTokenTypeID:
		subjectTokenType = IDTokenType
	default:
		return exchangedToken{}, fmt.Errorf("unsupported subject token type %q", providerTokenType)
	}

	body, err := json.Marshal(subjectTokenExchangeRequest{
		GrantType:          TokenExchangeGrantType,
		ClientID:           s.clientID,
		SubjectToken:       subjectToken,
		SubjectTokenType:   subjectTokenType,
		IdentityProviderID: identityProviderID,
		ServiceAccountID:   serviceAccountID,
	})
	if err != nil {
		return exchangedToken{}, fmt.Errorf("failed to marshal token exchange request: %w", err)
	}
	resp, err := exchangeToken(ctx, httpClient, TokenExchangeURL, body, 0)
	if err != nil {
		return exchangedToken{}, err
	}
	defer func() { _ = resp.Body.Close() }()
	responseBody, err := readTokenExchangeResponse(resp.Body, 0)
	if err != nil {
		return exchangedToken{}, err
	}
	return parseSubjectTokenExchangeResponse(resp.StatusCode, responseBody)
}

func (s subjectTokenCredentialSource) refreshBuffer(time.Duration) time.Duration {
	if s.refreshBefore == 0 {
		return DefaultRefreshBuffer
	}
	return s.refreshBefore
}

func NewWorkloadIdentityAuth(config WorkloadIdentity) (*WorkloadIdentityAuth, error) {
	if err := validateWorkloadIdentity("WorkloadIdentity", config.IdentityProviderID, config.ServiceAccountID); err != nil {
		return nil, err
	}
	if config.Provider == nil {
		return nil, fmt.Errorf("WorkloadIdentity: Provider is required")
	}
	if config.RefreshBufferSeconds < 0 {
		return nil, fmt.Errorf("WorkloadIdentity: RefreshBufferSeconds must be non-negative")
	}
	return newWorkloadIdentityAuth(
		config.IdentityProviderID,
		config.ServiceAccountID,
		subjectTokenCredentialSource{
			clientID:      config.ClientID,
			provider:      config.Provider,
			refreshBefore: time.Duration(config.RefreshBufferSeconds) * time.Second,
		},
	), nil
}

func validateWorkloadIdentity(configName string, identityProviderID string, serviceAccountID string) error {
	if identityProviderID == "" {
		return fmt.Errorf("%s: IdentityProviderID is required", configName)
	}
	if serviceAccountID == "" {
		return fmt.Errorf("%s: ServiceAccountID is required", configName)
	}
	return nil
}

func newWorkloadIdentityAuth(
	identityProviderID string,
	serviceAccountID string,
	source workloadIdentityCredentialSource,
) *WorkloadIdentityAuth {
	return &WorkloadIdentityAuth{
		identityProviderID: identityProviderID,
		serviceAccountID:   serviceAccountID,
		source:             source,
	}
}

func (w *WorkloadIdentityAuth) GetToken(ctx context.Context, httpClient HTTPDoer) (string, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	// Lock for entire decision: check cache, decide refresh strategy, potentially start background refresh
	w.mu.Lock()

	if w.cachedToken == "" {
		return w.handleLockedRefresh(ctx, httpClient)
	}

	now := time.Now()
	if !now.Before(w.tokenExpiry) {
		return w.handleLockedRefresh(ctx, httpClient)
	}

	// Proactive background refresh: start if within refresh window and no refresh active
	if !now.Before(w.tokenRefreshAt) && w.refreshInFlight == nil {
		state := w.beginRefreshLocked()
		// Background goroutine with independent context, lock released before spawn
		go func() {
			refreshCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			w.finishRefresh(state, w.refreshToken(refreshCtx, httpClient))
		}()
	}

	token := w.cachedToken
	w.mu.Unlock()
	return token, nil
}

// Single-flight pattern: ensures only one refresh runs, others wait for result
func (w *WorkloadIdentityAuth) handleLockedRefresh(ctx context.Context, httpClient HTTPDoer) (string, error) {
	if w.refreshInFlight == nil {
		// No refresh running: start foreground refresh, unlock before blocking operation
		state := w.beginRefreshLocked()
		w.mu.Unlock()
		result := w.refreshToken(ctx, httpClient)
		w.finishRefresh(state, result)
		return result.token, result.err
	}

	// Refresh already running: unlock and wait for its completion
	state := w.refreshInFlight
	w.mu.Unlock()
	return w.waitForRefresh(ctx, state)
}

func (w *WorkloadIdentityAuth) invalidateToken(token string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.cachedToken != token {
		return
	}
	w.cachedToken = ""
	w.tokenExpiry = time.Time{}
	w.tokenRefreshAt = time.Time{}
}

func (w *WorkloadIdentityAuth) beginRefreshLocked() *tokenRefreshState {
	w.refreshInFlight = &tokenRefreshState{done: make(chan struct{})}
	return w.refreshInFlight
}

// Atomically publishes refresh result and signals all waiting goroutines via channel close
func (w *WorkloadIdentityAuth) finishRefresh(state *tokenRefreshState, result tokenRefreshResult) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.refreshInFlight != state {
		return
	}
	state.result = result
	if result.err == nil {
		w.cachedToken = result.token
		w.tokenExpiry = result.expiresAt
		w.tokenRefreshAt = result.refreshAt
	}
	close(state.done) // Broadcasts completion to all waiters
	w.refreshInFlight = nil
}

// Blocks until refresh completes or context is canceled
func (w *WorkloadIdentityAuth) waitForRefresh(ctx context.Context, state *tokenRefreshState) (string, error) {
	select {
	case <-state.done: // Refresh completed
		return state.result.token, state.result.err
	case <-ctx.Done(): // Caller context canceled
		return "", ctx.Err()
	}
}

func (w *WorkloadIdentityAuth) refreshToken(ctx context.Context, httpClient HTTPDoer) tokenRefreshResult {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	token, err := w.source.exchange(ctx, httpClient, w.identityProviderID, w.serviceAccountID)
	if err != nil {
		return tokenRefreshResult{err: err}
	}

	now := time.Now()
	expiresAt := now.Add(token.expiresIn)

	return tokenRefreshResult{
		token:     token.accessToken,
		expiresAt: expiresAt,
		refreshAt: expiresAt.Add(-w.source.refreshBuffer(token.expiresIn)),
	}
}

func parseSubjectTokenExchangeResponse(statusCode int, body []byte) (exchangedToken, error) {
	if statusCode == http.StatusBadRequest || statusCode == http.StatusUnauthorized || statusCode == http.StatusForbidden {
		var oauthErr struct {
			Error            string `json:"error"`
			ErrorDescription string `json:"error_description"`
		}
		if json.Unmarshal(body, &oauthErr) == nil {
			return exchangedToken{}, &OAuthError{
				StatusCode:       statusCode,
				ErrorCode:        shared.OAuthErrorCode(oauthErr.Error),
				ErrorDescription: oauthErr.ErrorDescription,
			}
		}
	}
	if statusCode != http.StatusOK {
		return exchangedToken{}, fmt.Errorf("token exchange failed with status %d: %s", statusCode, string(body))
	}

	var tokenResp TokenExchangeResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return exchangedToken{}, fmt.Errorf("failed to decode token exchange response: %w", err)
	}
	if tokenResp.AccessToken == "" {
		return exchangedToken{}, fmt.Errorf("token exchange response missing 'access_token' field. Response: %s", string(body))
	}
	expiresIn := DefaultTokenExpiry
	if tokenResp.ExpiresIn != nil {
		expiresIn = time.Duration(*tokenResp.ExpiresIn) * time.Second
	}
	return exchangedToken{accessToken: tokenResp.AccessToken, expiresIn: expiresIn}, nil
}

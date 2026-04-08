package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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
	config WorkloadIdentity

	// Protects cachedToken, tokenExpiry, and refreshInFlight
	mu              sync.Mutex
	cachedToken     string
	tokenExpiry     time.Time
	refreshInFlight *tokenRefreshState
}

type tokenRefreshResult struct {
	token string
	err   error
}

// Coordinates concurrent access to a single in-flight refresh operation
// done channel signals completion to all waiting goroutines
type tokenRefreshState struct {
	done   chan struct{}
	result tokenRefreshResult
}

type tokenExchangeRequest struct {
	GrantType          string `json:"grant_type"`
	ClientID           string `json:"client_id"`
	SubjectToken       string `json:"subject_token"`
	SubjectTokenType   string `json:"subject_token_type"`
	IdentityProviderID string `json:"identity_provider_id"`
	ServiceAccountID   string `json:"service_account_id"`
}

func NewWorkloadIdentityAuth(config WorkloadIdentity) (*WorkloadIdentityAuth, error) {
	if config.ClientID == "" {
		return nil, fmt.Errorf("WorkloadIdentity: ClientID is required")
	}
	if config.IdentityProviderID == "" {
		return nil, fmt.Errorf("WorkloadIdentity: IdentityProviderID is required")
	}
	if config.ServiceAccountID == "" {
		return nil, fmt.Errorf("WorkloadIdentity: ServiceAccountID is required")
	}
	if config.Provider == nil {
		return nil, fmt.Errorf("WorkloadIdentity: Provider is required")
	}
	if config.RefreshBufferSeconds < 0 {
		return nil, fmt.Errorf("WorkloadIdentity: RefreshBufferSeconds must be non-negative")
	}
	return &WorkloadIdentityAuth{
		config: config,
	}, nil
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
	if now.After(w.tokenExpiry) {
		return w.handleLockedRefresh(ctx, httpClient)
	}

	refreshBuffer := w.config.RefreshBufferSeconds
	if refreshBuffer == 0 {
		refreshBuffer = int(DefaultRefreshBuffer / time.Second)
	}
	refreshTime := w.tokenExpiry.Add(-time.Duration(refreshBuffer) * time.Second)

	// Proactive background refresh: start if within refresh window and no refresh active
	if now.After(refreshTime) && w.refreshInFlight == nil {
		state := w.beginRefreshLocked()
		// Background goroutine with independent context, lock released before spawn
		go func() {
			refreshCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			token, err := w.refreshToken(refreshCtx, httpClient)
			w.finishRefresh(state, token, err)
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
		return w.completeForegroundRefresh(ctx, state, httpClient)
	}

	// Refresh already running: unlock and wait for its completion
	state := w.refreshInFlight
	w.mu.Unlock()
	return w.waitForRefresh(ctx, state)
}

func (w *WorkloadIdentityAuth) invalidateToken() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.cachedToken = ""
	w.tokenExpiry = time.Time{}
}

func (w *WorkloadIdentityAuth) beginRefreshLocked() *tokenRefreshState {
	w.refreshInFlight = &tokenRefreshState{done: make(chan struct{})}
	return w.refreshInFlight
}

func (w *WorkloadIdentityAuth) completeForegroundRefresh(ctx context.Context, state *tokenRefreshState, httpClient HTTPDoer) (string, error) {
	token, err := w.refreshToken(ctx, httpClient)
	w.finishRefresh(state, token, err)
	return token, err
}

// Atomically publishes refresh result and signals all waiting goroutines via channel close
func (w *WorkloadIdentityAuth) finishRefresh(state *tokenRefreshState, token string, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.refreshInFlight != state {
		return
	}
	state.result = tokenRefreshResult{token: token, err: err}
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

func (w *WorkloadIdentityAuth) refreshToken(ctx context.Context, httpClient HTTPDoer) (string, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	subjectToken, err := w.config.Provider.GetToken(ctx, httpClient)
	if err != nil {
		return "", err
	}

	subjectTokenType := w.config.Provider.TokenType()
	var subjectTokenTypeURN string
	switch subjectTokenType {
	case SubjectTokenTypeJWT:
		subjectTokenTypeURN = JWTTokenType
	case SubjectTokenTypeID:
		subjectTokenTypeURN = IDTokenType
	default:
		return "", fmt.Errorf("unsupported subject token type %q", subjectTokenType)
	}

	requestBody := tokenExchangeRequest{
		GrantType:          TokenExchangeGrantType,
		ClientID:           w.config.ClientID,
		SubjectToken:       subjectToken,
		SubjectTokenType:   subjectTokenTypeURN,
		IdentityProviderID: w.config.IdentityProviderID,
		ServiceAccountID:   w.config.ServiceAccountID,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal token exchange request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", TokenExchangeURL, bytes.NewReader(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create token exchange request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to exchange token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read token exchange response: %w", err)
	}

	if resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		var oauthErr struct {
			Error            string `json:"error"`
			ErrorDescription string `json:"error_description"`
		}
		if json.Unmarshal(body, &oauthErr) == nil {
			return "", &OAuthError{
				StatusCode:       resp.StatusCode,
				ErrorCode:        shared.OAuthErrorCode(oauthErr.Error),
				ErrorDescription: oauthErr.ErrorDescription,
			}
		}
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token exchange failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp TokenExchangeResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", fmt.Errorf("failed to decode token exchange response: %w", err)
	}

	if tokenResp.AccessToken == "" {
		return "", fmt.Errorf("token exchange response missing 'access_token' field. Response: %s", string(body))
	}

	expiresIn := int(DefaultTokenExpiry / time.Second)
	if tokenResp.ExpiresIn != nil {
		expiresIn = *tokenResp.ExpiresIn
	}

	// Atomically update cached token and expiry
	w.mu.Lock()
	w.cachedToken = tokenResp.AccessToken
	w.tokenExpiry = time.Now().Add(time.Duration(expiresIn) * time.Second)
	w.mu.Unlock()

	return tokenResp.AccessToken, nil
}

package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/openai/openai-go/v3/shared"
)

const (
	x509SubjectTokenType       = "urn:openai:params:oauth:token-type:x509"
	x509IssuedAccessTokenType  = "urn:ietf:params:oauth:token-type:access_token"
	x509SuccessResponseMaximum = 1 << 20
	x509ErrorResponseMaximum   = 4 << 10
	x509MaximumTokenLifetime   = 3600
	x509MaximumRetryAfter      = 8 * time.Second
)

type x509ExchangedToken struct {
	value     string
	expiresAt time.Time
}

type x509ExchangeHTTPError struct {
	statusCode    int
	retryAfter    time.Duration
	hasRetryAfter bool
}

func (err *x509ExchangeHTTPError) Error() string {
	return fmt.Sprintf("X.509 token exchange failed with HTTP status %d", err.statusCode)
}

func (err *x509ExchangeHTTPError) retryable() bool {
	return err.statusCode == http.StatusRequestTimeout ||
		err.statusCode == http.StatusConflict ||
		err.statusCode == http.StatusTooManyRequests ||
		(err.statusCode >= http.StatusInternalServerError && err.statusCode < 600)
}

type x509ExchangeReadError struct{}

func (*x509ExchangeReadError) Error() string {
	return "X.509 token exchange response could not be read"
}
func (*x509ExchangeReadError) retryable() bool { return true }

func x509Exchange(ctx context.Context, transport *X509Transport, identityProviderID, serviceAccountID string) (x509ExchangedToken, error) {
	if ctx == nil {
		return x509ExchangedToken{}, errors.New("X.509 token exchange requires a non-nil context")
	}
	if err := ctx.Err(); err != nil {
		return x509ExchangedToken{}, err
	}
	if strings.TrimSpace(identityProviderID) == "" || strings.TrimSpace(serviceAccountID) == "" {
		return x509ExchangedToken{}, errors.New("X.509 token exchange requires identity-provider and service-account IDs")
	}
	if err := transport.validateAttestation(); err != nil {
		return x509ExchangedToken{}, err
	}
	started := time.Now()
	requestBody, err := json.Marshal(struct {
		GrantType          string `json:"grant_type"`
		SubjectTokenType   string `json:"subject_token_type"`
		IdentityProviderID string `json:"identity_provider_id"`
		ServiceAccountID   string `json:"service_account_id"`
	}{
		GrantType:          TokenExchangeGrantType,
		SubjectTokenType:   x509SubjectTokenType,
		IdentityProviderID: identityProviderID,
		ServiceAccountID:   serviceAccountID,
	})
	if err != nil {
		return x509ExchangedToken{}, errors.New("X.509 token exchange could not encode its request")
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://"+x509AuthenticationHost+"/oauth/token", bytes.NewReader(requestBody))
	if err != nil {
		return x509ExchangedToken{}, errors.New("X.509 token exchange could not construct its request")
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := transport.Do(request)
	if err != nil {
		if contextErr := ctx.Err(); contextErr != nil {
			return x509ExchangedToken{}, contextErr
		}
		return x509ExchangedToken{}, err
	}
	defer func() { _ = response.Body.Close() }()

	if response.StatusCode != http.StatusOK {
		return x509ExchangedToken{}, x509ExchangeStatusError(ctx, response)
	}
	body, err := x509ReadExchangeResponse(ctx, response, x509SuccessResponseMaximum)
	if err != nil {
		return x509ExchangedToken{}, err
	}
	return x509DecodeExchangedToken(ctx, body, started)
}

func x509ReadExchangeResponse(ctx context.Context, response *http.Response, maximum int64) ([]byte, error) {
	if response.ContentLength > maximum {
		return nil, errors.New("X.509 token exchange response exceeds its size limit")
	}
	body, err := io.ReadAll(io.LimitReader(response.Body, maximum+1))
	if err != nil {
		if contextErr := ctx.Err(); contextErr != nil {
			return nil, contextErr
		}
		return nil, &x509ExchangeReadError{}
	}
	if int64(len(body)) > maximum {
		return nil, errors.New("X.509 token exchange response exceeds its size limit")
	}
	return body, nil
}

func x509DecodeObject(body []byte) (map[string]json.RawMessage, error) {
	decoder := json.NewDecoder(bytes.NewReader(body))
	opening, err := decoder.Token()
	if err != nil || opening != json.Delim('{') {
		return nil, errors.New("X.509 token exchange response must be a valid JSON object")
	}
	fields := make(map[string]json.RawMessage)
	for decoder.More() {
		key, keyErr := decoder.Token()
		if keyErr != nil {
			return nil, errors.New("X.509 token exchange response contains invalid JSON")
		}
		name, ok := key.(string)
		if !ok {
			return nil, errors.New("X.509 token exchange response contains an invalid field")
		}
		if _, exists := fields[name]; exists {
			return nil, errors.New("X.509 token exchange response contains duplicate fields")
		}
		var value json.RawMessage
		if decodeErr := decoder.Decode(&value); decodeErr != nil {
			return nil, errors.New("X.509 token exchange response contains invalid JSON")
		}
		fields[name] = value
	}
	closing, err := decoder.Token()
	if err != nil || closing != json.Delim('}') {
		return nil, errors.New("X.509 token exchange response contains invalid JSON")
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		return nil, errors.New("X.509 token exchange response contains trailing JSON")
	}
	return fields, nil
}

func x509DecodeExchangedToken(ctx context.Context, body []byte, started time.Time) (x509ExchangedToken, error) {
	fields, err := x509DecodeObject(body)
	if err != nil {
		return x509ExchangedToken{}, err
	}
	var accessToken, tokenType, issuedTokenType string
	if json.Unmarshal(fields["access_token"], &accessToken) != nil ||
		!validX509BearerHeader("Bearer "+accessToken) {
		return x509ExchangedToken{}, errors.New("X.509 token exchange returned an invalid access token")
	}
	if json.Unmarshal(fields["token_type"], &tokenType) != nil || !strings.EqualFold(tokenType, "Bearer") {
		return x509ExchangedToken{}, errors.New("X.509 token exchange returned an invalid token type")
	}
	if json.Unmarshal(fields["issued_token_type"], &issuedTokenType) != nil ||
		issuedTokenType != x509IssuedAccessTokenType {
		return x509ExchangedToken{}, errors.New("X.509 token exchange returned an invalid issued token type")
	}
	var lifetime int64
	if json.Unmarshal(fields["expires_in"], &lifetime) != nil || lifetime < 1 || lifetime > x509MaximumTokenLifetime {
		return x509ExchangedToken{}, errors.New("X.509 token exchange returned an invalid token lifetime")
	}
	if err := ctx.Err(); err != nil {
		return x509ExchangedToken{}, err
	}
	expiresAt := started.Add(time.Duration(lifetime) * time.Second)
	if !time.Now().Before(expiresAt) {
		return x509ExchangedToken{}, errors.New("X.509 token exchange returned an already expired token")
	}
	return x509ExchangedToken{value: accessToken, expiresAt: expiresAt}, nil
}

func x509ExchangeStatusError(ctx context.Context, response *http.Response) (result error) {
	if err := ctx.Err(); err != nil {
		return err
	}
	defer func() {
		if err := ctx.Err(); err != nil {
			result = err
		}
	}()
	if response.StatusCode != http.StatusBadRequest && response.StatusCode != http.StatusUnauthorized &&
		response.StatusCode != http.StatusForbidden {
		status := &x509ExchangeHTTPError{statusCode: response.StatusCode}
		if status.retryable() {
			status.retryAfter, status.hasRetryAfter = x509ParseRetryAfter(response.Header, time.Now())
			if err := x509DrainRetryableResponse(ctx, response); err != nil {
				return err
			}
		}
		return status
	}
	oauthError := &OAuthError{StatusCode: response.StatusCode}
	body, err := x509ReadExchangeResponse(ctx, response, x509ErrorResponseMaximum)
	if err != nil {
		if contextErr := ctx.Err(); contextErr != nil {
			return contextErr
		}
		return oauthError
	}
	fields, err := x509DecodeObject(body)
	if err != nil {
		return oauthError
	}
	var code string
	if json.Unmarshal(fields["error"], &code) != nil {
		envelope, envelopeErr := x509DecodeObject(fields["error"])
		if envelopeErr != nil || json.Unmarshal(envelope["code"], &code) != nil {
			return oauthError
		}
	}
	switch code {
	case "invalid_request", "invalid_client", "invalid_grant", "invalid_scope", "invalid_target", "unauthorized_client",
		"unsupported_grant_type", "invalid_subject_token", "access_denied", "temporarily_unavailable", "server_error":
		oauthError.ErrorCode = shared.OAuthErrorCode(code)
	}
	return oauthError
}

func x509DrainRetryableResponse(ctx context.Context, response *http.Response) error {
	if response.ContentLength > x509ErrorResponseMaximum {
		return nil
	}
	// Probe beyond the boundary so an exact-limit chunked body reaches its real EOF.
	if _, err := io.Copy(io.Discard, io.LimitReader(response.Body, x509ErrorResponseMaximum+1)); err != nil {
		if contextErr := ctx.Err(); contextErr != nil {
			return contextErr
		}
		return &x509ExchangeReadError{}
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	return nil
}

func x509ParseRetryAfter(headers http.Header, now time.Time) (time.Duration, bool) {
	for _, header := range []struct {
		name string
		unit time.Duration
	}{
		{name: "Retry-After-Ms", unit: time.Millisecond},
		{name: "Retry-After", unit: time.Second},
	} {
		value := headers.Get(header.name)
		if value == "" {
			continue
		}
		if parsed, err := strconv.ParseFloat(value, 64); err == nil {
			if math.IsNaN(parsed) || math.IsInf(parsed, 0) || parsed < 0 {
				continue
			}
			if parsed >= float64(x509MaximumRetryAfter)/float64(header.unit) {
				return x509MaximumRetryAfter, true
			}
			return time.Duration(parsed * float64(header.unit)), true
		}
		if header.name != "Retry-After" {
			continue
		}
		if deadline, err := http.ParseTime(value); err == nil {
			delay := deadline.Sub(now)
			if delay <= 0 {
				return 0, true
			}
			return min(delay, x509MaximumRetryAfter), true
		}
	}
	return 0, false
}

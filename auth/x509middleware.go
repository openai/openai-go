package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/openai/openai-go/v3/internal/requestconfig"
)

type x509IssuerRefusalContextKey struct{}

// X509WorkloadIdentityMiddleware authenticates a request using its attested
// workload identity. A rejected bearer is conditionally invalidated and the
// request is replayed at most once when its body can be recreated safely.
func X509WorkloadIdentityMiddleware(
	identity *X509WorkloadIdentityAuth,
	httpClient HTTPDoer,
	request *http.Request,
	next func(*http.Request) (*http.Response, error),
) (*http.Response, error) {
	if request == nil || request.Header == nil || next == nil {
		return nil, errors.New("X.509 workload identity requires a non-nil request and header map")
	}
	if err := validateX509Request(request); err != nil {
		return nil, err
	}
	if request.URL.Hostname() != x509APIHost {
		return nil, errors.New("X.509 workload identity requires the global OpenAI mTLS API origin")
	}
	for name := range request.Header {
		if strings.EqualFold(strings.ReplaceAll(name, "_", "-"), "authorization") {
			return nil, errors.New("X.509 workload identity cannot replace existing Authorization credentials")
		}
	}
	hadBody := request.Body != nil && request.Body != http.NoBody
	tokenContext := request.Context()
	if requestconfig.RequestRetryScopeFromContext(tokenContext) == nil {
		// Reuse only the refusal storage, without installing a retry budget or
		// changing standalone issuer retries and the single middleware replay.
		tokenContext = context.WithValue(tokenContext, x509IssuerRefusalContextKey{},
			requestconfig.NewRequestRetryScope(0, 0, false, nil))
	}
	token, err := identity.GetToken(tokenContext, httpClient)
	if err != nil {
		return nil, err
	}
	authenticated := request.Clone(request.Context())
	authenticated.Header.Set("Authorization", "Bearer "+token)
	response, err := next(authenticated)
	response = x509UnsignedResponse(response)
	if err != nil || response == nil || response.StatusCode != http.StatusUnauthorized {
		return response, err
	}
	identity.invalidateToken(token)
	scope := requestconfig.RequestRetryScopeFromContext(request.Context())
	if scope != nil {
		replayable := !hadBody || (request.GetBody != nil && scope.AllowBodyReplay())
		if replayable && scope.TryOuterReplay() {
			return unauthorizedRetryResponse(response, true), nil
		}
		if response.Header.Get("x-should-retry") == "true" {
			return unauthorizedRetryResponse(response, false), nil
		}
		return response, nil
	}
	if hadBody {
		return response, nil
	}
	if retry, waitErr := waitForUnauthorizedReplay(request.Context(), response); waitErr != nil {
		return nil, waitErr
	} else if !retry {
		return response, nil
	}
	if response.Body != nil {
		_ = response.Body.Close()
	}
	replay := request.Clone(request.Context())
	token, err = identity.GetToken(tokenContext, httpClient)
	if err != nil {
		return nil, err
	}
	replay.Header.Set("Authorization", "Bearer "+token)
	response, err = next(replay)
	response = x509UnsignedResponse(response)
	if err == nil && response != nil && response.StatusCode == http.StatusUnauthorized {
		identity.invalidateToken(token)
	}
	return response, err
}

func x509UnsignedResponse(response *http.Response) *http.Response {
	if response == nil || response.Request == nil || response.Request.Header == nil {
		return response
	}
	for name := range response.Request.Header {
		if !strings.EqualFold(strings.ReplaceAll(name, "_", "-"), "authorization") {
			continue
		}
		clone := *response
		clone.Request = response.Request.Clone(response.Request.Context())
		for credential := range clone.Request.Header {
			if strings.EqualFold(strings.ReplaceAll(credential, "_", "-"), "authorization") {
				delete(clone.Request.Header, credential)
			}
		}
		return &clone
	}
	return response
}

func unauthorizedRetryResponse(response *http.Response, retry bool) *http.Response {
	clone := *response
	clone.Header = response.Header.Clone()
	if clone.Header == nil {
		clone.Header = make(http.Header)
	}
	if retry {
		clone.Header.Set("x-should-retry", "true")
		if _, hasHint, _ := requestconfig.AuthenticationRetryDelay(response, 0); !hasHint {
			clone.Header.Set("Retry-After-Ms", "0")
		}
	} else {
		clone.Header.Set("x-should-retry", "false")
	}
	return &clone
}

func waitForUnauthorizedReplay(ctx context.Context, response *http.Response) (bool, error) {
	err := ctx.Err()
	if err == nil {
		delay, _, allowed := requestconfig.AuthenticationRetryDelay(response, 0)
		if !allowed {
			return false, nil
		}
		err = requestconfig.WaitForDelay(ctx, delay)
	}
	if err != nil {
		if response.Body != nil {
			_ = response.Body.Close()
		}
		return false, err
	}
	return true, nil
}

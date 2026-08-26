package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/openai/openai-go/v3/internal/requestconfig"
)

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
	token, err := identity.GetToken(request.Context(), httpClient)
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
	if response.Body != nil {
		_ = response.Body.Close()
	}
	replay := request.Clone(request.Context())
	token, err = identity.GetToken(request.Context(), httpClient)
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
		clone.Header.Set("Retry-After-Ms", "0")
	} else {
		clone.Header.Set("x-should-retry", "false")
	}
	return &clone
}

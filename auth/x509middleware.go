package auth

import (
	"errors"
	"net/http"

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
	if request == nil || request.Header == nil {
		return nil, errors.New("X.509 workload identity requires a non-nil request and header map")
	}
	hadBody := request.Body != nil
	token, err := identity.GetToken(request.Context(), httpClient)
	if err != nil {
		return nil, err
	}
	authenticated := request.Clone(request.Context())
	authenticated.Header.Set("Authorization", "Bearer "+token)
	response, err := next(authenticated)
	if err != nil || response == nil || response.StatusCode != http.StatusUnauthorized {
		return response, err
	}
	identity.invalidateToken(token)
	scope := requestconfig.RequestRetryScopeFromContext(request.Context())
	if hadBody && (request.GetBody == nil || scope == nil || !scope.AllowBodyReplay()) {
		return response, nil
	}
	if scope != nil && !scope.TryReplay() {
		return response, nil
	}
	_ = response.Body.Close()
	replay := request.Clone(request.Context())
	if hadBody {
		replay.Body, err = request.GetBody()
		if err != nil {
			return nil, err
		}
	}
	token, err = identity.GetToken(request.Context(), httpClient)
	if err != nil {
		if replay.Body != nil {
			_ = replay.Body.Close()
		}
		return nil, err
	}
	replay.Header.Set("Authorization", "Bearer "+token)
	response, err = next(replay)
	if err == nil && response != nil && response.StatusCode == http.StatusUnauthorized {
		identity.invalidateToken(token)
	}
	return response, err
}

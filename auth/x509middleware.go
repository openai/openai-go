package auth

import (
	"net/http"
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
	if request.Body != nil && request.GetBody == nil {
		return response, nil
	}
	_ = response.Body.Close()
	identity.invalidateToken(token)
	replay := request.Clone(request.Context())
	if request.GetBody != nil {
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
	return next(replay)
}

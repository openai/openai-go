package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/openai/openai-go/v3/internal/requestconfig"
)

func WorkloadIdentityMiddleware(
	wia *WorkloadIdentityAuth,
	httpClient HTTPDoer,
	req *http.Request,
	next func(*http.Request) (*http.Response, error),
) (*http.Response, error) {
	if req == nil || req.Header == nil || next == nil {
		return nil, errors.New("workload identity requires a non-nil request, header map, and handler")
	}
	hadBody := req.Body != nil && req.Body != http.NoBody
	token, err := wia.GetToken(req.Context(), httpClient)
	if err != nil {
		return nil, err
	}

	authenticated := req.Clone(req.Context())
	setWorkloadIdentityAuthorization(authenticated, token)

	resp, err := next(authenticated)
	resp = x509UnsignedResponse(resp)
	if err != nil || resp == nil || resp.StatusCode != http.StatusUnauthorized {
		return resp, err
	}

	wia.invalidateToken(token)

	if scope := requestconfig.RequestRetryScopeFromContext(req.Context()); scope != nil {
		replayable := !hadBody || (req.GetBody != nil && scope.AllowBodyReplay())
		if replayable && scope.TryOuterReplay() {
			return unauthorizedRetryResponse(resp, true), nil
		}
		if resp.Header.Get("x-should-retry") == "true" {
			return unauthorizedRetryResponse(resp, false), nil
		}
		return resp, nil
	}
	// Direct callers have no original SDK request from which to safely rebuild
	// a body after caller middleware may have transformed or removed it.
	if hadBody {
		return resp, nil
	}

	retryReq := req.Clone(req.Context())

	token, err = wia.GetToken(req.Context(), httpClient)
	if err != nil {
		if resp.Body != nil {
			_ = resp.Body.Close()
		}
		return nil, err
	}
	setWorkloadIdentityAuthorization(retryReq, token)

	if resp.Body != nil {
		_ = resp.Body.Close()
	}
	resp, err = next(retryReq)
	if err == nil && resp != nil && resp.StatusCode == http.StatusUnauthorized {
		wia.invalidateToken(token)
	}
	return x509UnsignedResponse(resp), err
}

func setWorkloadIdentityAuthorization(request *http.Request, token string) {
	for name := range request.Header {
		if strings.EqualFold(strings.ReplaceAll(name, "_", "-"), "authorization") {
			delete(request.Header, name)
		}
	}
	request.Header.Set("Authorization", "Bearer "+token)
}

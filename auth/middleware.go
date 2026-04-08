package auth

import (
	"net/http"
)

func WorkloadIdentityMiddleware(
	wia *WorkloadIdentityAuth,
	httpClient HTTPDoer,
	req *http.Request,
	next func(*http.Request) (*http.Response, error),
) (*http.Response, error) {
	token, err := wia.GetToken(req.Context(), httpClient)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := next(req)
	if err != nil || resp == nil || resp.StatusCode != http.StatusUnauthorized {
		return resp, err
	}

	wia.invalidateToken()

	if req.Body != nil && req.GetBody == nil {
		return resp, nil
	}

	retryReq := req.Clone(req.Context())

	token, err = wia.GetToken(req.Context(), httpClient)
	if err != nil {
		resp.Body.Close()
		return nil, err
	}
	retryReq.Header.Set("Authorization", "Bearer "+token)

	if req.GetBody != nil {
		retryReq.Body, err = req.GetBody()
		if err != nil {
			resp.Body.Close()
			return nil, err
		}
	}

	resp.Body.Close()
	return next(retryReq)
}

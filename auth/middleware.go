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
	var retryReq *http.Request
	if req.Body == nil || req.GetBody != nil {
		retryReq = req.Clone(req.Context())
	}

	resp, err := next(req)
	if err != nil || resp == nil || resp.StatusCode != http.StatusUnauthorized {
		return resp, err
	}

	wia.invalidateToken(token)

	if retryReq == nil {
		return resp, nil
	}

	token, err = wia.GetToken(req.Context(), httpClient)
	if err != nil {
		_ = resp.Body.Close()
		return nil, err
	}
	retryReq.Header.Set("Authorization", "Bearer "+token)

	if retryReq.GetBody != nil {
		retryReq.Body, err = retryReq.GetBody()
		if err != nil {
			_ = resp.Body.Close()
			return nil, err
		}
	}

	_ = resp.Body.Close()
	resp, err = next(retryReq)
	if err != nil && retryReq.Body != nil {
		_ = retryReq.Body.Close()
	}
	return resp, err
}

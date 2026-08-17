package auth

import (
	"net/http"

	"github.com/openai/openai-go/v3/internal"
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

	authorization := "Bearer " + token
	req = internal.WithExpectedAuthorization(req, authorization)
	req.Header.Set("Authorization", authorization)
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
	authorization = "Bearer " + token
	retryReq.Header.Set("Authorization", authorization)
	retryReq = internal.WithExpectedAuthorization(retryReq, authorization)

	if retryReq.GetBody != nil {
		retryReq.Body, err = retryReq.GetBody()
		if err != nil {
			_ = resp.Body.Close()
			return nil, err
		}
	}

	_ = resp.Body.Close()
	if retryReq.Body != nil {
		retryReq.Body = internal.NewCloseOnceReadCloser(retryReq.Body)
	}
	resp, err = next(retryReq)
	if err != nil && retryReq.Body != nil {
		_ = retryReq.Body.Close()
	}
	if resp != nil && resp.StatusCode == http.StatusUnauthorized {
		wia.invalidateToken(token)
	}
	return resp, err
}

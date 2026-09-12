package openai_test

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/openai/openai-go/v3/auth"
)

type originTestRoundTripper func(*http.Request) (*http.Response, error)

func (f originTestRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

type originTestSubjectTokenProvider struct{}

func (originTestSubjectTokenProvider) TokenType() auth.SubjectTokenType {
	return auth.SubjectTokenTypeJWT
}

func (originTestSubjectTokenProvider) GetToken(context.Context, auth.HTTPDoer) (string, error) {
	return "subject-token", nil
}

func originTestHTTPClient(fallback http.RoundTripper) *http.Client {
	return &http.Client{Transport: originTestRoundTripper(func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == auth.TokenExchangeURL {
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": {"application/json"}},
				Body:       io.NopCloser(strings.NewReader(`{"access_token":"workload-token","expires_in":3600}`)),
				Request:    req,
			}, nil
		}
		return fallback.RoundTrip(req)
	})}
}

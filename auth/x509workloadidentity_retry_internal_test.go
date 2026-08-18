package auth

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
)

type failingTokenResponseBody struct {
	read bool
}

func (b *failingTokenResponseBody) Read(p []byte) (int, error) {
	if b.read {
		return 0, io.ErrUnexpectedEOF
	}
	b.read = true
	return copy(p, `{"access_token":`), nil
}

func (*failingTokenResponseBody) Close() error { return nil }

func TestX509TokenExchangeRetriesTransientResponseReadFailure(t *testing.T) {
	var calls atomic.Int32
	doer := &internalHTTPDoer{do: func(*http.Request) (*http.Response, error) {
		if calls.Add(1) == 1 {
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Retry-After": []string{"0"}},
				Body:       &failingTokenResponseBody{},
			}, nil
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`{"access_token":"retry-token","expires_in":60}`)),
		}, nil
	}}

	token, err := (x509CredentialSource{}).exchange(t.Context(), doer, "idp-test", "svc-test")
	if err != nil {
		t.Fatalf("exchange() error = %v", err)
	}
	if token.accessToken != "retry-token" {
		t.Fatalf("exchange() token = %q, want retry-token", token.accessToken)
	}
	if got, want := calls.Load(), int32(2); got != want {
		t.Fatalf("exchange calls = %d, want %d", got, want)
	}
}

func TestX509TokenExchangeResponseReadRetriesAreBounded(t *testing.T) {
	var calls atomic.Int32
	doer := &internalHTTPDoer{do: func(*http.Request) (*http.Response, error) {
		calls.Add(1)
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Retry-After": []string{"0"}},
			Body:       &failingTokenResponseBody{},
		}, nil
	}}

	_, err := (x509CredentialSource{}).exchange(t.Context(), doer, "idp-test", "svc-test")
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("exchange() error = %v, want io.ErrUnexpectedEOF", err)
	}
	if got, want := calls.Load(), int32(tokenExchangeMaxRetries+1); got != want {
		t.Fatalf("exchange calls = %d, want %d", got, want)
	}
}

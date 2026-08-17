package auth_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/openai/openai-go/v3/auth"
	"github.com/openai/openai-go/v3/internal"
)

type closeTrackingReadCloser struct {
	io.ReadCloser
	closes atomic.Int32
}

func (b *closeTrackingReadCloser) Close() error {
	b.closes.Add(1)
	return b.ReadCloser.Close()
}

type dynamicallyNonComparableHTTPDoer struct {
	state any
	doer  auth.HTTPDoer
}

func (d dynamicallyNonComparableHTTPDoer) Do(req *http.Request) (*http.Response, error) {
	return d.doer.Do(req)
}

func testX509WorkloadIdentity() auth.X509WorkloadIdentity {
	return auth.X509WorkloadIdentity{
		IdentityProviderID: "idp-test",
		ServiceAccountID:   "svc-test",
	}
}

func tokenResponse(token string, expiresIn any) *http.Response {
	body, err := json.Marshal(map[string]any{
		"access_token": token,
		"expires_in":   expiresIn,
	})
	if err != nil {
		panic(err)
	}
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(string(body))),
	}
}

func TestX509WorkloadIdentityValidation(t *testing.T) {
	testCases := []struct {
		name   string
		config auth.X509WorkloadIdentity
	}{
		{name: "missing identity provider", config: auth.X509WorkloadIdentity{ServiceAccountID: "svc-test"}},
		{name: "missing service account", config: auth.X509WorkloadIdentity{IdentityProviderID: "idp-test"}},
		{name: "negative refresh buffer", config: auth.X509WorkloadIdentity{
			IdentityProviderID: "idp-test",
			ServiceAccountID:   "svc-test",
			RefreshBuffer:      -time.Second,
		}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if _, err := auth.NewX509WorkloadIdentityAuth(testCase.config); err == nil {
				t.Fatal("NewX509WorkloadIdentityAuth() error = nil")
			}
		})
	}
}

func TestX509WorkloadIdentityDoesNotOwnCertificateMaterial(t *testing.T) {
	typeOfConfig := reflect.TypeFor[auth.X509WorkloadIdentity]()
	gotFields := make([]string, 0, typeOfConfig.NumField())
	for i := range typeOfConfig.NumField() {
		gotFields = append(gotFields, typeOfConfig.Field(i).Name)
	}
	wantFields := []string{"IdentityProviderID", "ServiceAccountID", "RefreshBuffer"}
	if !reflect.DeepEqual(gotFields, wantFields) {
		t.Errorf("X509WorkloadIdentity fields = %v, want %v", gotFields, wantFields)
	}
}

func TestX509WorkloadIdentityMiddlewareRequiresRequestPolicyBeforeExchange(t *testing.T) {
	var exchangeCalls atomic.Int32
	var nextCalls atomic.Int32
	exchangeClient := &http.Client{Transport: &closureTransport{fn: func(*http.Request) (*http.Response, error) {
		exchangeCalls.Add(1)
		return tokenResponse("x509-token", 60), nil
	}}}
	wia, err := auth.NewX509WorkloadIdentityAuth(testX509WorkloadIdentity())
	if err != nil {
		t.Fatalf("NewX509WorkloadIdentityAuth() error = %v", err)
	}
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "http://attacker.invalid/collect", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := auth.WorkloadIdentityMiddleware(wia, exchangeClient, req, func(*http.Request) (*http.Response, error) {
		nextCalls.Add(1)
		return &http.Response{StatusCode: http.StatusOK, Body: http.NoBody}, nil
	})
	if err == nil {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
		t.Fatal("WorkloadIdentityMiddleware() error = nil")
	}
	if resp != nil {
		if resp.Body != nil {
			_ = resp.Body.Close()
		}
		t.Errorf("WorkloadIdentityMiddleware() response = %#v, want nil", resp)
	}
	if got := exchangeCalls.Load(); got != 0 {
		t.Errorf("token exchange calls = %d, want 0", got)
	}
	if got := nextCalls.Load(); got != 0 {
		t.Errorf("next calls = %d, want 0", got)
	}
}

func TestX509WorkloadIdentityAuthRejectsHTTPDoerChange(t *testing.T) {
	var callsA atomic.Int32
	var callsB atomic.Int32
	httpClientA := &http.Client{Transport: &closureTransport{fn: func(*http.Request) (*http.Response, error) {
		callsA.Add(1)
		return tokenResponse("token-a", 60), nil
	}}}
	httpClientB := &http.Client{Transport: &closureTransport{fn: func(*http.Request) (*http.Response, error) {
		callsB.Add(1)
		return tokenResponse("token-b", 60), nil
	}}}
	wia, err := auth.NewX509WorkloadIdentityAuth(testX509WorkloadIdentity())
	if err != nil {
		t.Fatalf("NewX509WorkloadIdentityAuth() error = %v", err)
	}

	token, err := wia.GetToken(t.Context(), httpClientA)
	if err != nil {
		t.Fatalf("first GetToken() error = %v", err)
	}
	if token != "token-a" {
		t.Fatalf("first GetToken() = %q, want token-a", token)
	}
	_, err = wia.GetToken(t.Context(), httpClientB)
	if err == nil {
		t.Fatal("GetToken() with different HTTP client error = nil")
	}
	token, err = wia.GetToken(t.Context(), httpClientA)
	if err != nil {
		t.Fatalf("cached GetToken() error = %v", err)
	}
	if token != "token-a" {
		t.Errorf("cached GetToken() = %q, want token-a", token)
	}
	if got, want := callsA.Load(), int32(1); got != want {
		t.Errorf("HTTP client A exchange calls = %d, want %d", got, want)
	}
	if got := callsB.Load(); got != 0 {
		t.Errorf("HTTP client B exchange calls = %d, want 0", got)
	}
}

func TestX509WorkloadIdentityAuthRejectsDynamicallyNonComparableHTTPDoer(t *testing.T) {
	httpDoer := dynamicallyNonComparableHTTPDoer{
		state: []string{"not-comparable"},
		doer: &http.Client{Transport: &closureTransport{fn: func(*http.Request) (*http.Response, error) {
			t.Fatal("dynamically non-comparable HTTP doer was called")
			return nil, nil
		}}},
	}
	wia, err := auth.NewX509WorkloadIdentityAuth(testX509WorkloadIdentity())
	if err != nil {
		t.Fatalf("NewX509WorkloadIdentityAuth() error = %v", err)
	}

	_, err = wia.GetToken(t.Context(), httpDoer)
	if err == nil || !strings.Contains(err.Error(), "requires a comparable HTTP client") {
		t.Fatalf("GetToken() error = %v, want comparable HTTP client error", err)
	}
}

func TestX509TokenExchangeRequest(t *testing.T) {
	var requestBody map[string]any
	httpClient := &http.Client{Transport: &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
		if got, want := req.Method, http.MethodPost; got != want {
			t.Errorf("request method = %q, want %q", got, want)
		}
		if got, want := req.URL.String(), "https://mtls.auth.openai.com/oauth/token"; got != want {
			t.Errorf("request URL = %q, want %q", got, want)
		}
		if got, want := req.Header.Get("Content-Type"), "application/json"; got != want {
			t.Errorf("Content-Type = %q, want %q", got, want)
		}
		body, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(body, &requestBody); err != nil {
			return nil, err
		}
		return tokenResponse("x509-token", 3600), nil
	}}}

	wia, err := auth.NewX509WorkloadIdentityAuth(testX509WorkloadIdentity())
	if err != nil {
		t.Fatalf("NewX509WorkloadIdentityAuth() error = %v", err)
	}
	token, err := wia.GetToken(t.Context(), httpClient)
	if err != nil {
		t.Fatalf("GetToken() error = %v", err)
	}
	if got, want := token, "x509-token"; got != want {
		t.Errorf("GetToken() = %q, want %q", got, want)
	}

	wantBody := map[string]any{
		"grant_type":           "urn:ietf:params:oauth:grant-type:token-exchange",
		"subject_token_type":   "urn:openai:params:oauth:token-type:x509",
		"identity_provider_id": "idp-test",
		"service_account_id":   "svc-test",
	}
	if !reflect.DeepEqual(requestBody, wantBody) {
		t.Errorf("request body = %#v, want %#v", requestBody, wantBody)
	}
	if _, ok := requestBody["subject_token"]; ok {
		t.Error("request body contains subject_token")
	}
	if _, ok := requestBody["client_id"]; ok {
		t.Error("request body contains client_id")
	}
}

func TestX509TokenExchangeDoesNotRetryOAuthErrors(t *testing.T) {
	for _, statusCode := range []int{http.StatusBadRequest, http.StatusUnauthorized, http.StatusForbidden} {
		t.Run(http.StatusText(statusCode), func(t *testing.T) {
			var calls atomic.Int32
			httpClient := &http.Client{Transport: &closureTransport{fn: func(*http.Request) (*http.Response, error) {
				calls.Add(1)
				return &http.Response{
					StatusCode: statusCode,
					Header:     http.Header{"Content-Type": []string{"application/json"}},
					Body: io.NopCloser(strings.NewReader(
						`{"error":"invalid_grant","error_description":"generic exchange failure"}`,
					)),
				}, nil
			}}}
			wia, err := auth.NewX509WorkloadIdentityAuth(testX509WorkloadIdentity())
			if err != nil {
				t.Fatalf("NewX509WorkloadIdentityAuth() error = %v", err)
			}
			_, err = wia.GetToken(t.Context(), httpClient)
			var oauthErr *auth.OAuthError
			if !errors.As(err, &oauthErr) {
				t.Fatalf("GetToken() error = %v, want *auth.OAuthError", err)
			}
			if got, want := oauthErr.StatusCode, statusCode; got != want {
				t.Errorf("OAuthError.StatusCode = %d, want %d", got, want)
			}
			if strings.Contains(err.Error(), "generic exchange failure") {
				t.Errorf("GetToken() error contains response body: %v", err)
			}
			if got, want := calls.Load(), int32(1); got != want {
				t.Errorf("token exchange calls = %d, want %d", got, want)
			}
		})
	}
}

func TestX509TokenExchangeRequiresPositiveExpiresIn(t *testing.T) {
	testCases := []struct {
		name string
		body string
	}{
		{name: "missing", body: `{"access_token":"secret-token"}`},
		{name: "zero", body: `{"access_token":"secret-token","expires_in":0}`},
		{name: "negative", body: `{"access_token":"secret-token","expires_in":-1}`},
		{name: "not numeric", body: `{"access_token":"secret-token","expires_in":"3600"}`},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			wia, err := auth.NewX509WorkloadIdentityAuth(testX509WorkloadIdentity())
			if err != nil {
				t.Fatalf("NewX509WorkloadIdentityAuth() error = %v", err)
			}
			httpClient := mockOAuthServer(testCase.body, http.StatusOK)
			if _, err := wia.GetToken(t.Context(), httpClient); err == nil {
				t.Fatal("GetToken() error = nil")
			} else if strings.Contains(err.Error(), "secret-token") {
				t.Errorf("GetToken() error contains access token: %v", err)
			}
		})
	}
}

func TestX509TokenExchangeValidatesBearerTokenTypeAndGrammar(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		for _, body := range []string{
			`{"access_token":"abc-._~+/==","expires_in":60}`,
			`{"access_token":"abc-._~+/==","expires_in":60,"token_type":"Bearer"}`,
			`{"access_token":"abc-._~+/==","expires_in":60,"token_type":"bearer"}`,
		} {
			wia, err := auth.NewX509WorkloadIdentityAuth(testX509WorkloadIdentity())
			if err != nil {
				t.Fatalf("NewX509WorkloadIdentityAuth() error = %v", err)
			}
			if _, err := wia.GetToken(t.Context(), mockOAuthServer(body, http.StatusOK)); err != nil {
				t.Errorf("GetToken() body %s error = %v", body, err)
			}
		}
	})

	t.Run("invalid token type", func(t *testing.T) {
		for _, tokenType := range []string{`"Basic"`, `"MAC"`, `"DPoP"`, `null`, `123`} {
			body := fmt.Sprintf(`{"access_token":"must-not-leak","expires_in":60,"token_type":%s}`, tokenType)
			wia, err := auth.NewX509WorkloadIdentityAuth(testX509WorkloadIdentity())
			if err != nil {
				t.Fatalf("NewX509WorkloadIdentityAuth() error = %v", err)
			}
			_, err = wia.GetToken(t.Context(), mockOAuthServer(body, http.StatusOK))
			if err == nil {
				t.Errorf("GetToken() token_type %s error = nil", tokenType)
			} else if strings.Contains(err.Error(), "must-not-leak") {
				t.Errorf("GetToken() error contains access token: %v", err)
			}
		}
	})

	t.Run("invalid access token", func(t *testing.T) {
		invalidTokens := []string{
			" leading", "trailing ", "two words", "line\r\nbreak", "line\nbreak",
			"nul\x00byte", "tab\tbyte", "nonascii-\u00e9", "bad!punctuation", "=padding-only", "abc=more",
		}
		for _, accessToken := range invalidTokens {
			body, err := json.Marshal(map[string]any{
				"access_token": accessToken,
				"expires_in":   60,
				"token_type":   "Bearer",
			})
			if err != nil {
				t.Fatal(err)
			}
			wia, err := auth.NewX509WorkloadIdentityAuth(testX509WorkloadIdentity())
			if err != nil {
				t.Fatalf("NewX509WorkloadIdentityAuth() error = %v", err)
			}
			_, err = wia.GetToken(t.Context(), mockOAuthServer(string(body), http.StatusOK))
			if err == nil {
				t.Errorf("GetToken() token %q error = nil", accessToken)
			} else if strings.Contains(err.Error(), accessToken) {
				t.Errorf("GetToken() error contains access token: %v", err)
			}
		}
	})

	t.Run("invalid token is not cached", func(t *testing.T) {
		var calls atomic.Int32
		httpClient := &http.Client{Transport: &closureTransport{fn: func(*http.Request) (*http.Response, error) {
			if calls.Add(1) == 1 {
				return tokenResponse("line\r\nbreak", 60), nil
			}
			return tokenResponse("safe-token", 60), nil
		}}}
		wia, err := auth.NewX509WorkloadIdentityAuth(testX509WorkloadIdentity())
		if err != nil {
			t.Fatalf("NewX509WorkloadIdentityAuth() error = %v", err)
		}
		_, err = wia.GetToken(t.Context(), httpClient)
		if err == nil {
			t.Fatal("first GetToken() error = nil")
		}
		token, err := wia.GetToken(t.Context(), httpClient)
		if err != nil {
			t.Fatalf("second GetToken() error = %v", err)
		}
		if got, want := token, "safe-token"; got != want {
			t.Fatalf("second GetToken() = %q, want %q", got, want)
		}
		if got, want := calls.Load(), int32(2); got != want {
			t.Fatalf("token exchange calls = %d, want %d", got, want)
		}
	})
}

func TestIDTokenWorkloadIdentityRegression(t *testing.T) {
	provider := &mockProvider{token: "id-subject-token", tokenType: auth.SubjectTokenTypeID}
	var requestURL string
	var requestBody map[string]any
	httpClient := &http.Client{Transport: &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
		requestURL = req.URL.String()
		body, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(body, &requestBody); err != nil {
			return nil, err
		}
		return tokenResponse("id-token-exchange", 3600), nil
	}}}
	wia, err := auth.NewWorkloadIdentityAuth(auth.WorkloadIdentity{
		IdentityProviderID: "idp-test",
		ServiceAccountID:   "svc-test",
		Provider:           provider,
	})
	if err != nil {
		t.Fatalf("NewWorkloadIdentityAuth() error = %v", err)
	}
	if _, err := wia.GetToken(t.Context(), httpClient); err != nil {
		t.Fatalf("GetToken() error = %v", err)
	}
	if got, want := requestURL, "https://auth.openai.com/oauth/token"; got != want {
		t.Errorf("request URL = %q, want %q", got, want)
	}
	if got, want := requestBody["subject_token"], "id-subject-token"; got != want {
		t.Errorf("subject_token = %v, want %q", got, want)
	}
	if got, want := requestBody["subject_token_type"], "urn:ietf:params:oauth:token-type:id_token"; got != want {
		t.Errorf("subject_token_type = %v, want %q", got, want)
	}
}

func TestX509TokenExchangeRetriesTransientFailures(t *testing.T) {
	var calls atomic.Int32
	httpClient := &http.Client{Transport: &closureTransport{fn: func(*http.Request) (*http.Response, error) {
		switch calls.Add(1) {
		case 1:
			return nil, errors.New("temporary connection failure")
		case 2:
			return &http.Response{
				StatusCode: http.StatusTooManyRequests,
				Header:     http.Header{"Retry-After": []string{"0"}},
				Body:       io.NopCloser(strings.NewReader(`{"error":"rate_limit"}`)),
			}, nil
		default:
			return tokenResponse("retried-token", 60), nil
		}
	}}}

	wia, err := auth.NewX509WorkloadIdentityAuth(testX509WorkloadIdentity())
	if err != nil {
		t.Fatalf("NewX509WorkloadIdentityAuth() error = %v", err)
	}
	if _, err := wia.GetToken(t.Context(), httpClient); err != nil {
		t.Fatalf("GetToken() error = %v", err)
	}
	if got, want := calls.Load(), int32(3); got != want {
		t.Errorf("token exchange calls = %d, want %d", got, want)
	}
}

func TestX509TokenExchangeRetriesAreBounded(t *testing.T) {
	var calls atomic.Int32
	httpClient := &http.Client{Transport: &closureTransport{fn: func(*http.Request) (*http.Response, error) {
		calls.Add(1)
		return &http.Response{
			StatusCode: http.StatusServiceUnavailable,
			Header:     http.Header{"Retry-After": []string{"0"}},
			Body:       io.NopCloser(strings.NewReader(`{"access_token":"must-not-leak"}`)),
		}, nil
	}}}

	wia, err := auth.NewX509WorkloadIdentityAuth(testX509WorkloadIdentity())
	if err != nil {
		t.Fatalf("NewX509WorkloadIdentityAuth() error = %v", err)
	}
	_, err = wia.GetToken(t.Context(), httpClient)
	if err == nil {
		t.Fatal("GetToken() error = nil")
	}
	if got, want := calls.Load(), int32(3); got != want {
		t.Errorf("token exchange calls = %d, want %d", got, want)
	}
	if strings.Contains(err.Error(), "must-not-leak") {
		t.Errorf("GetToken() error contains response body: %v", err)
	}
}

func TestX509TokenExchangeRetryWaitHonorsContext(t *testing.T) {
	requestStarted := make(chan struct{})
	var once sync.Once
	httpClient := &http.Client{Transport: &closureTransport{fn: func(*http.Request) (*http.Response, error) {
		once.Do(func() { close(requestStarted) })
		return &http.Response{
			StatusCode: http.StatusTooManyRequests,
			Header:     http.Header{"Retry-After": []string{"60"}},
			Body:       io.NopCloser(strings.NewReader(`{"error":"rate_limit"}`)),
		}, nil
	}}}

	wia, err := auth.NewX509WorkloadIdentityAuth(testX509WorkloadIdentity())
	if err != nil {
		t.Fatalf("NewX509WorkloadIdentityAuth() error = %v", err)
	}
	ctx, cancel := context.WithCancel(t.Context())
	done := make(chan error, 1)
	go func() {
		_, getTokenErr := wia.GetToken(ctx, httpClient)
		done <- getTokenErr
	}()
	<-requestStarted
	cancel()
	if err := <-done; !errors.Is(err, context.Canceled) {
		t.Errorf("GetToken() error = %v, want context.Canceled", err)
	}
}

func TestX509TokenExchangeRefusesRedirects(t *testing.T) {
	var redirectedRequests atomic.Int32
	redirectTarget := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		redirectedRequests.Add(1)
	}))
	t.Cleanup(redirectTarget.Close)

	httpClient := &http.Client{
		Transport: &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
			if req.URL.Host == "mtls.auth.openai.com" {
				return &http.Response{
					StatusCode: http.StatusFound,
					Header:     http.Header{"Location": []string{redirectTarget.URL}},
					Body:       io.NopCloser(strings.NewReader("redirect")),
				}, nil
			}
			return http.DefaultTransport.RoundTrip(req)
		}},
		CheckRedirect: func(*http.Request, []*http.Request) error { return nil },
	}

	wia, err := auth.NewX509WorkloadIdentityAuth(testX509WorkloadIdentity())
	if err != nil {
		t.Fatalf("NewX509WorkloadIdentityAuth() error = %v", err)
	}
	if _, err := wia.GetToken(t.Context(), httpClient); err == nil {
		t.Fatal("GetToken() error = nil")
	}
	if got := redirectedRequests.Load(); got != 0 {
		t.Errorf("redirect target requests = %d, want 0", got)
	}
}

func TestX509CanceledWaiterDoesNotCancelRefresh(t *testing.T) {
	exchangeStarted := make(chan struct{})
	releaseExchange := make(chan struct{})
	httpClient := &http.Client{Transport: &closureTransport{fn: func(*http.Request) (*http.Response, error) {
		close(exchangeStarted)
		<-releaseExchange
		return tokenResponse("shared-token", 60), nil
	}}}

	wia, err := auth.NewX509WorkloadIdentityAuth(testX509WorkloadIdentity())
	if err != nil {
		t.Fatalf("NewX509WorkloadIdentityAuth() error = %v", err)
	}
	leaderDone := make(chan error, 1)
	go func() {
		_, getTokenErr := wia.GetToken(t.Context(), httpClient)
		leaderDone <- getTokenErr
	}()
	<-exchangeStarted

	waiterCtx, cancelWaiter := context.WithCancel(t.Context())
	cancelWaiter()
	if _, err := wia.GetToken(waiterCtx, httpClient); !errors.Is(err, context.Canceled) {
		t.Errorf("waiter GetToken() error = %v, want context.Canceled", err)
	}
	select {
	case err := <-leaderDone:
		t.Fatalf("shared refresh completed before release with error %v", err)
	default:
	}

	close(releaseExchange)
	if err := <-leaderDone; err != nil {
		t.Fatalf("leader GetToken() error = %v", err)
	}
}

func TestX509CanceledOnlyWaiterCancelsSharedRefresh(t *testing.T) {
	exchangeStarted := make(chan struct{})
	exchangeCanceled := make(chan struct{})
	httpClient := &http.Client{Transport: &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
		close(exchangeStarted)
		<-req.Context().Done()
		close(exchangeCanceled)
		return nil, req.Context().Err()
	}}}
	wia, err := auth.NewX509WorkloadIdentityAuth(testX509WorkloadIdentity())
	if err != nil {
		t.Fatalf("NewX509WorkloadIdentityAuth() error = %v", err)
	}
	ctx, cancel := context.WithCancel(t.Context())
	done := make(chan error, 1)
	go func() {
		_, getTokenErr := wia.GetToken(ctx, httpClient)
		done <- getTokenErr
	}()
	<-exchangeStarted
	cancel()
	if err := <-done; !errors.Is(err, context.Canceled) {
		t.Errorf("GetToken() error = %v, want context.Canceled", err)
	}
	select {
	case <-exchangeCanceled:
	case <-time.After(time.Second):
		t.Fatal("shared exchange context was not canceled")
	}
}

func TestX509InitialRefreshIsSingleFlight(t *testing.T) {
	exchangeStarted := make(chan struct{})
	releaseExchange := make(chan struct{})
	var calls atomic.Int32
	var startOnce sync.Once
	httpClient := &http.Client{Transport: &closureTransport{fn: func(*http.Request) (*http.Response, error) {
		calls.Add(1)
		startOnce.Do(func() { close(exchangeStarted) })
		<-releaseExchange
		return tokenResponse("shared-token", 60), nil
	}}}
	wia, err := auth.NewX509WorkloadIdentityAuth(testX509WorkloadIdentity())
	if err != nil {
		t.Fatalf("NewX509WorkloadIdentityAuth() error = %v", err)
	}

	const goroutines = 16
	start := make(chan struct{})
	results := make(chan error, goroutines)
	for range goroutines {
		go func() {
			<-start
			token, getTokenErr := wia.GetToken(t.Context(), httpClient)
			if getTokenErr == nil && token != "shared-token" {
				getTokenErr = fmt.Errorf("GetToken() = %q, want shared-token", token)
			}
			results <- getTokenErr
		}()
	}
	close(start)
	<-exchangeStarted
	close(releaseExchange)
	for range goroutines {
		if err := <-results; err != nil {
			t.Error(err)
		}
	}
	if got, want := calls.Load(), int32(1); got != want {
		t.Errorf("token exchange calls = %d, want %d", got, want)
	}
}

func TestX509ExpiredTokenRefreshesSynchronously(t *testing.T) {
	var calls atomic.Int32
	httpClient := &http.Client{Transport: &closureTransport{fn: func(*http.Request) (*http.Response, error) {
		call := calls.Add(1)
		return tokenResponse(fmt.Sprintf("token-%d", call), 0.1), nil
	}}}
	config := testX509WorkloadIdentity()
	config.RefreshBuffer = time.Nanosecond
	wia, err := auth.NewX509WorkloadIdentityAuth(config)
	if err != nil {
		t.Fatalf("NewX509WorkloadIdentityAuth() error = %v", err)
	}

	if token, err := wia.GetToken(t.Context(), httpClient); err != nil || token != "token-1" {
		t.Fatalf("first GetToken() = %q, %v; want token-1, nil", token, err)
	}
	time.Sleep(150 * time.Millisecond)
	if token, err := wia.GetToken(t.Context(), httpClient); err != nil || token != "token-2" {
		t.Fatalf("second GetToken() = %q, %v; want token-2, nil", token, err)
	}
	if got, want := calls.Load(), int32(2); got != want {
		t.Errorf("token exchange calls = %d, want %d", got, want)
	}
}

func TestX509RefreshBufferClampedToHalfTTL(t *testing.T) {
	var calls atomic.Int32
	secondExchange := make(chan struct{})
	httpClient := &http.Client{Transport: &closureTransport{fn: func(*http.Request) (*http.Response, error) {
		call := calls.Add(1)
		if call == 2 {
			close(secondExchange)
		}
		return tokenResponse(fmt.Sprintf("token-%d", call), 1), nil
	}}}
	config := testX509WorkloadIdentity()
	config.RefreshBuffer = 10 * time.Second
	wia, err := auth.NewX509WorkloadIdentityAuth(config)
	if err != nil {
		t.Fatalf("NewX509WorkloadIdentityAuth() error = %v", err)
	}

	if _, err := wia.GetToken(t.Context(), httpClient); err != nil {
		t.Fatalf("first GetToken() error = %v", err)
	}
	if _, err := wia.GetToken(t.Context(), httpClient); err != nil {
		t.Fatalf("second GetToken() error = %v", err)
	}
	if got := calls.Load(); got != 1 {
		t.Fatalf("immediate token exchange calls = %d, want 1", got)
	}

	time.Sleep(600 * time.Millisecond)
	if _, err := wia.GetToken(t.Context(), httpClient); err != nil {
		t.Fatalf("proactive GetToken() error = %v", err)
	}
	select {
	case <-secondExchange:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for proactive refresh")
	}
}

func TestX509Concurrent401InvalidationKeepsNewToken(t *testing.T) {
	var exchangeCalls atomic.Int32
	exchangeClient := &http.Client{Transport: &closureTransport{fn: func(*http.Request) (*http.Response, error) {
		call := exchangeCalls.Add(1)
		return tokenResponse(fmt.Sprintf("token-%d", call), 60), nil
	}}}
	wia, err := auth.NewX509WorkloadIdentityAuth(testX509WorkloadIdentity())
	if err != nil {
		t.Fatalf("NewX509WorkloadIdentityAuth() error = %v", err)
	}
	if _, err := wia.GetToken(t.Context(), exchangeClient); err != nil {
		t.Fatalf("initial GetToken() error = %v", err)
	}

	bothOldRequestsStarted := make(chan struct{})
	newTokenUsed := make(chan struct{})
	var oldTokenCalls atomic.Int32
	var bothOnce sync.Once
	var newOnce sync.Once
	next := func(req *http.Request) (*http.Response, error) {
		switch req.Header.Get("Authorization") {
		case "Bearer token-1":
			call := oldTokenCalls.Add(1)
			if call == 2 {
				bothOnce.Do(func() { close(bothOldRequestsStarted) })
				<-newTokenUsed
			} else {
				<-bothOldRequestsStarted
			}
			return &http.Response{StatusCode: http.StatusUnauthorized, Body: http.NoBody, Header: make(http.Header)}, nil
		case "Bearer token-2":
			newOnce.Do(func() { close(newTokenUsed) })
			return &http.Response{StatusCode: http.StatusOK, Body: http.NoBody, Header: make(http.Header)}, nil
		default:
			return nil, fmt.Errorf("unexpected authorization header %q", req.Header.Get("Authorization"))
		}
	}

	results := make(chan error, 2)
	for range 2 {
		go func() {
			req, reqErr := http.NewRequestWithContext(t.Context(), http.MethodGet, "https://mtls.api.openai.com/v1/models", nil)
			if reqErr != nil {
				results <- reqErr
				return
			}
			resp, middlewareErr := auth.WorkloadIdentityMiddleware(wia, exchangeClient, internal.WithX509RequestPolicy(req), next)
			if resp != nil && resp.Body != nil {
				_ = resp.Body.Close()
			}
			results <- middlewareErr
		}()
	}
	for range 2 {
		if err := <-results; err != nil {
			t.Errorf("WorkloadIdentityMiddleware() error = %v", err)
		}
	}
	if got, want := exchangeCalls.Load(), int32(2); got != want {
		t.Errorf("token exchange calls = %d, want %d", got, want)
	}
}

func TestX509Second401InvalidatesReplayToken(t *testing.T) {
	var exchangeCalls atomic.Int32
	exchangeClient := &http.Client{Transport: &closureTransport{fn: func(*http.Request) (*http.Response, error) {
		call := exchangeCalls.Add(1)
		return tokenResponse(fmt.Sprintf("token-%d", call), 60), nil
	}}}
	wia, err := auth.NewX509WorkloadIdentityAuth(testX509WorkloadIdentity())
	if err != nil {
		t.Fatalf("NewX509WorkloadIdentityAuth() error = %v", err)
	}
	var authorizations []string
	next := func(req *http.Request) (*http.Response, error) {
		authorizations = append(authorizations, req.Header.Get("Authorization"))
		return &http.Response{
			StatusCode: http.StatusUnauthorized,
			Header:     make(http.Header),
			Body:       http.NoBody,
		}, nil
	}

	for range 2 {
		req, reqErr := http.NewRequestWithContext(t.Context(), http.MethodGet, "https://mtls.api.openai.com/v1/models", nil)
		if reqErr != nil {
			t.Fatal(reqErr)
		}
		resp, middlewareErr := auth.WorkloadIdentityMiddleware(wia, exchangeClient, internal.WithX509RequestPolicy(req), next)
		if middlewareErr != nil {
			t.Fatalf("WorkloadIdentityMiddleware() error = %v", middlewareErr)
		}
		_ = resp.Body.Close()
	}

	if got, want := strings.Join(authorizations, ","), "Bearer token-1,Bearer token-2,Bearer token-3,Bearer token-4"; got != want {
		t.Fatalf("Authorization sequence = %q, want %q", got, want)
	}
}

func TestX509TokenExchangeDoesNotInheritHTTPClientCookies(t *testing.T) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	authURL, err := url.Parse("https://mtls.auth.openai.com/")
	if err != nil {
		t.Fatal(err)
	}
	jar.SetCookies(authURL, []*http.Cookie{{Name: "api-session", Value: "secret", Path: "/"}})
	var exchangeCookie string
	httpClient := &http.Client{
		Jar: jar,
		Transport: &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
			exchangeCookie = req.Header.Get("Cookie")
			return tokenResponse("x509-token", 60), nil
		}},
	}
	wia, err := auth.NewX509WorkloadIdentityAuth(testX509WorkloadIdentity())
	if err != nil {
		t.Fatalf("NewX509WorkloadIdentityAuth() error = %v", err)
	}

	if _, err := wia.GetToken(t.Context(), httpClient); err != nil {
		t.Fatalf("GetToken() error = %v", err)
	}
	if exchangeCookie != "" {
		t.Fatalf("token exchange Cookie = %q, want empty", exchangeCookie)
	}
}

func TestX509WorkloadIdentity401ClosesReplayBodyWhenMiddlewareFails(t *testing.T) {
	testX509ReplayBodyClose(t, false)
}

func TestX509WorkloadIdentity401DoesNotDoubleCloseReplayBodyOnDownstreamError(t *testing.T) {
	testX509ReplayBodyClose(t, true)
}

func testX509ReplayBodyClose(t *testing.T, downstreamClosesBody bool) {
	t.Helper()

	var exchangeCalls atomic.Int32
	exchangeClient := &http.Client{Transport: &closureTransport{fn: func(*http.Request) (*http.Response, error) {
		call := exchangeCalls.Add(1)
		return tokenResponse(fmt.Sprintf("token-%d", call), 60), nil
	}}}
	wia, err := auth.NewX509WorkloadIdentityAuth(testX509WorkloadIdentity())
	if err != nil {
		t.Fatalf("NewX509WorkloadIdentityAuth() error = %v", err)
	}

	bodyPath := t.TempDir() + "/body"
	if writeErr := os.WriteFile(bodyPath, []byte("request body"), 0o600); writeErr != nil {
		t.Fatal(writeErr)
	}
	req, err := http.NewRequestWithContext(
		t.Context(),
		http.MethodPost,
		"https://mtls.api.openai.com/v1/custom",
		strings.NewReader("request body"),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = req.Body.Close() }()

	var replayBody *closeTrackingReadCloser
	req.GetBody = func() (io.ReadCloser, error) {
		body, openErr := os.Open(bodyPath)
		if openErr != nil {
			return nil, openErr
		}
		replayBody = &closeTrackingReadCloser{ReadCloser: body}
		return replayBody, nil
	}

	middlewareErr := errors.New("middleware rejected replay")
	var calls int
	next := func(replayReq *http.Request) (*http.Response, error) {
		calls++
		if calls == 1 {
			return &http.Response{
				StatusCode: http.StatusUnauthorized,
				Body:       http.NoBody,
				Header:     make(http.Header),
			}, nil
		}
		if downstreamClosesBody {
			_ = replayReq.Body.Close()
		}
		return nil, middlewareErr
	}

	resp, err := auth.WorkloadIdentityMiddleware(wia, exchangeClient, internal.WithX509RequestPolicy(req), next)
	if !errors.Is(err, middlewareErr) {
		t.Fatalf("WorkloadIdentityMiddleware() error = %v, want %v", err, middlewareErr)
	}
	if resp != nil {
		if resp.Body != nil {
			_ = resp.Body.Close()
		}
		t.Errorf("WorkloadIdentityMiddleware() response = %#v, want nil", resp)
	}
	if got, want := calls, 2; got != want {
		t.Errorf("middleware calls = %d, want %d", got, want)
	}
	if got, want := exchangeCalls.Load(), int32(2); got != want {
		t.Errorf("token exchange calls = %d, want %d", got, want)
	}
	if replayBody == nil {
		t.Fatal("GetBody was not called")
	}
	if got, want := replayBody.closes.Load(), int32(1); got != want {
		t.Errorf("replay body closes = %d, want %d", got, want)
	}
	if _, err := replayBody.Read(make([]byte, 1)); err == nil {
		t.Error("replay body remains readable after middleware error")
	}
}

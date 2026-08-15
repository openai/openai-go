package openai_test

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/auth"
	"github.com/openai/openai-go/v3/option"
)

type nonComparableHTTPDoer struct {
	do func(*http.Request) (*http.Response, error)
}

func (d nonComparableHTTPDoer) Do(req *http.Request) (*http.Response, error) {
	return d.do(req)
}

// dynamicallyNonComparableHTTPDoer has a comparable Go type, but values are
// not comparable when state contains a slice, map, or function.
type dynamicallyNonComparableHTTPDoer struct {
	state any
	doer  auth.HTTPDoer
}

func (d dynamicallyNonComparableHTTPDoer) Do(req *http.Request) (*http.Response, error) {
	return d.doer.Do(req)
}

func clientX509WorkloadIdentity() auth.X509WorkloadIdentity {
	return auth.X509WorkloadIdentity{
		IdentityProviderID: "idp-test",
		ServiceAccountID:   "svc-test",
	}
}

func modelsListResponse() *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(`{"object":"list","data":[]}`)),
	}
}

func TestClientX509WorkloadIdentityDefaultsMTLSAndReusesHTTPClient(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "api-key-must-not-be-used")
	var mu sync.Mutex
	requests := make([]*http.Request, 0, 2)
	httpClient := &http.Client{Transport: &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
		mu.Lock()
		requests = append(requests, req.Clone(req.Context()))
		mu.Unlock()
		if req.URL.Host == "mtls.auth.openai.com" {
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       io.NopCloser(strings.NewReader(`{"access_token":"x509-token","expires_in":3600}`)),
			}, nil
		}
		return modelsListResponse(), nil
	}}}

	client := openai.NewClient(
		option.WithX509WorkloadIdentity(clientX509WorkloadIdentity()),
		option.WithHTTPClient(httpClient),
	)
	mu.Lock()
	requestCount := len(requests)
	mu.Unlock()
	if requestCount != 0 {
		t.Fatalf("NewClient() made %d HTTP requests, want 0", requestCount)
	}

	if _, err := client.Models.List(t.Context()); err != nil {
		t.Fatalf("Models.List() error = %v", err)
	}
	mu.Lock()
	defer mu.Unlock()
	if got, want := len(requests), 2; got != want {
		t.Fatalf("HTTP requests = %d, want %d", got, want)
	}
	if got, want := requests[0].URL.String(), "https://mtls.auth.openai.com/oauth/token"; got != want {
		t.Errorf("exchange URL = %q, want %q", got, want)
	}
	if got, want := requests[1].URL.String(), "https://mtls.api.openai.com/v1/models"; got != want {
		t.Errorf("API URL = %q, want %q", got, want)
	}
	if got, want := requests[1].Header.Get("Authorization"), "Bearer x509-token"; got != want {
		t.Errorf("API Authorization = %q, want %q", got, want)
	}
	if strings.Contains(requests[1].Header.Get("Authorization"), "api-key-must-not-be-used") {
		t.Error("API Authorization contains ambient API key")
	}
}

func TestX509ChangesDoNotAffectAPIKeyClient(t *testing.T) {
	var requests []*http.Request
	httpClient := &http.Client{Transport: &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
		requests = append(requests, req.Clone(req.Context()))
		return modelsListResponse(), nil
	}}}
	client := openai.NewClient(
		option.WithAPIKey("api-key-token"),
		option.WithHTTPClient(httpClient),
	)
	if _, err := client.Models.List(t.Context()); err != nil {
		t.Fatalf("Models.List() error = %v", err)
	}
	if got, want := len(requests), 1; got != want {
		t.Fatalf("HTTP requests = %d, want %d", got, want)
	}
	if got, want := requests[0].URL.String(), "https://api.openai.com/v1/models"; got != want {
		t.Errorf("request URL = %q, want %q", got, want)
	}
	if got, want := requests[0].Header.Get("Authorization"), "Bearer api-key-token"; got != want {
		t.Errorf("Authorization = %q, want %q", got, want)
	}
}

func TestClientX509WorkloadIdentityBaseURLPrecedence(t *testing.T) {
	testCases := []struct {
		name    string
		opts    func(*http.Client) []option.RequestOption
		wantURL string
	}{
		{
			name: "explicit before x509",
			opts: func(httpClient *http.Client) []option.RequestOption {
				return []option.RequestOption{
					option.WithBaseURL("https://explicit.example/v1"),
					option.WithX509WorkloadIdentity(clientX509WorkloadIdentity()),
					option.WithHTTPClient(httpClient),
				}
			},
			wantURL: "https://explicit.example/v1/models",
		},
		{
			name: "explicit after x509",
			opts: func(httpClient *http.Client) []option.RequestOption {
				return []option.RequestOption{
					option.WithX509WorkloadIdentity(clientX509WorkloadIdentity()),
					option.WithBaseURL("https://explicit.example/v1"),
					option.WithHTTPClient(httpClient),
				}
			},
			wantURL: "https://explicit.example/v1/models",
		},
		{
			name: "production default after x509",
			opts: func(httpClient *http.Client) []option.RequestOption {
				return []option.RequestOption{
					option.WithX509WorkloadIdentity(clientX509WorkloadIdentity()),
					option.WithEnvironmentProduction(),
					option.WithHTTPClient(httpClient),
				}
			},
			wantURL: "https://mtls.api.openai.com/v1/models",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var apiURL string
			httpClient := &http.Client{Transport: &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
				if req.URL.Host == "mtls.auth.openai.com" {
					return &http.Response{
						StatusCode: http.StatusOK,
						Header:     http.Header{"Content-Type": []string{"application/json"}},
						Body:       io.NopCloser(strings.NewReader(`{"access_token":"x509-token","expires_in":3600}`)),
					}, nil
				}
				apiURL = req.URL.String()
				return modelsListResponse(), nil
			}}}
			client := openai.NewClient(testCase.opts(httpClient)...)
			if _, err := client.Models.List(t.Context()); err != nil {
				t.Fatalf("Models.List() error = %v", err)
			}
			if apiURL != testCase.wantURL {
				t.Errorf("API URL = %q, want %q", apiURL, testCase.wantURL)
			}
		})
	}
}

func TestClientX509WorkloadIdentityEnvironmentBaseURLWins(t *testing.T) {
	t.Setenv("OPENAI_BASE_URL", "https://environment.example/v1")
	var apiURL string
	httpClient := &http.Client{Transport: &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
		if req.URL.Host == "mtls.auth.openai.com" {
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       io.NopCloser(strings.NewReader(`{"access_token":"x509-token","expires_in":3600}`)),
			}, nil
		}
		apiURL = req.URL.String()
		return modelsListResponse(), nil
	}}}
	client := openai.NewClient(
		option.WithX509WorkloadIdentity(clientX509WorkloadIdentity()),
		option.WithHTTPClient(httpClient),
	)
	if _, err := client.Models.List(t.Context()); err != nil {
		t.Fatalf("Models.List() error = %v", err)
	}
	if got, want := apiURL, "https://environment.example/v1/models"; got != want {
		t.Errorf("API URL = %q, want %q", got, want)
	}
}

func TestClientRejectsMultipleWorkloadIdentityCredentialSources(t *testing.T) {
	provider := &mockSubjectTokenProvider{token: "subject-token", tokenType: auth.SubjectTokenTypeJWT}
	client := openai.NewClient(
		option.WithWorkloadIdentity(testWorkloadIdentity(provider)),
		option.WithX509WorkloadIdentity(clientX509WorkloadIdentity()),
	)
	if _, err := client.Models.List(t.Context()); err == nil {
		t.Fatal("Models.List() error = nil")
	}
}

func TestClientX509WorkloadIdentityRefusesAPIRedirects(t *testing.T) {
	var redirectedRequests atomic.Int32
	var apiRequests atomic.Int32
	redirectTarget := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		redirectedRequests.Add(1)
	}))
	t.Cleanup(redirectTarget.Close)

	httpClient := &http.Client{
		Transport: &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
			if req.URL.Host == "mtls.auth.openai.com" {
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     http.Header{"Content-Type": []string{"application/json"}},
					Body:       io.NopCloser(strings.NewReader(`{"access_token":"x509-token","expires_in":3600}`)),
				}, nil
			}
			if req.URL.Host == "api.example" {
				apiRequests.Add(1)
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
	client := openai.NewClient(
		option.WithX509WorkloadIdentity(clientX509WorkloadIdentity()),
		option.WithBaseURL("https://api.example/v1"),
		option.WithHTTPClient(httpClient),
	)
	if err := client.Execute(t.Context(), http.MethodDelete, "/custom", nil, nil); err == nil {
		t.Fatal("Execute() error = nil")
	}
	if got := redirectedRequests.Load(); got != 0 {
		t.Errorf("redirect target requests = %d, want 0", got)
	}
	if got := apiRequests.Load(); got != 1 {
		t.Errorf("API requests = %d, want 1", got)
	}
}

func TestClientX509WorkloadIdentityScopesTokenCacheToHTTPDoer(t *testing.T) {
	exchangeCalls := map[string]int{}
	apiAuth := map[string][]string{}
	newHTTPClient := func(name string) *http.Client {
		return &http.Client{Transport: &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
			if req.URL.Host == "mtls.auth.openai.com" {
				exchangeCalls[name]++
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     http.Header{"Content-Type": []string{"application/json"}},
					Body:       io.NopCloser(strings.NewReader(fmt.Sprintf(`{"access_token":"token-%s","expires_in":3600}`, name))),
				}, nil
			}
			apiAuth[name] = append(apiAuth[name], req.Header.Get("Authorization"))
			return modelsListResponse(), nil
		}}}
	}

	httpClientA := newHTTPClient("a")
	httpClientB := newHTTPClient("b")
	client := openai.NewClient(
		option.WithX509WorkloadIdentity(clientX509WorkloadIdentity()),
		option.WithHTTPClient(httpClientA),
	)

	if _, err := client.Models.List(t.Context()); err != nil {
		t.Fatalf("first Models.List() error = %v", err)
	}
	if _, err := client.Models.List(t.Context(), option.WithHTTPClient(httpClientB)); err != nil {
		t.Fatalf("overridden Models.List() error = %v", err)
	}
	if _, err := client.Models.List(t.Context()); err != nil {
		t.Fatalf("second Models.List() error = %v", err)
	}

	if got, want := exchangeCalls["a"], 1; got != want {
		t.Errorf("transport A exchange calls = %d, want %d", got, want)
	}
	if got, want := exchangeCalls["b"], 1; got != want {
		t.Errorf("transport B exchange calls = %d, want %d", got, want)
	}
	if got, want := strings.Join(apiAuth["a"], ","), "Bearer token-a,Bearer token-a"; got != want {
		t.Errorf("transport A Authorization = %q, want %q", got, want)
	}
	if got, want := strings.Join(apiAuth["b"], ","), "Bearer token-b"; got != want {
		t.Errorf("transport B Authorization = %q, want %q", got, want)
	}
}

func TestClientX509WorkloadIdentityCachesNonComparableHTTPDoer(t *testing.T) {
	var exchangeCalls atomic.Int32
	var apiCalls atomic.Int32
	httpDoer := nonComparableHTTPDoer{do: func(req *http.Request) (*http.Response, error) {
		if req.URL.Host == "mtls.auth.openai.com" {
			exchangeCalls.Add(1)
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       io.NopCloser(strings.NewReader(`{"access_token":"x509-token","expires_in":3600}`)),
			}, nil
		}
		apiCalls.Add(1)
		return modelsListResponse(), nil
	}}
	const requestCount = 16
	var ready sync.WaitGroup
	ready.Add(requestCount)
	releaseRequests := make(chan struct{})
	client := openai.NewClient(
		option.WithX509WorkloadIdentity(clientX509WorkloadIdentity()),
		option.WithHTTPClient(httpDoer),
		option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
			ready.Done()
			<-releaseRequests
			return next(req)
		}),
	)

	errs := make(chan error, requestCount)
	var requests sync.WaitGroup
	requests.Add(requestCount)
	for range requestCount {
		go func() {
			defer requests.Done()
			_, err := client.Models.List(t.Context())
			errs <- err
		}()
	}
	ready.Wait()
	close(releaseRequests)
	requests.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatalf("Models.List() error = %v", err)
		}
	}
	if got, want := exchangeCalls.Load(), int32(1); got != want {
		t.Errorf("token exchange calls = %d, want %d", got, want)
	}
	if got, want := apiCalls.Load(), int32(requestCount); got != want {
		t.Errorf("API calls = %d, want %d", got, want)
	}
}

func TestClientX509WorkloadIdentityCachesDynamicallyNonComparableHTTPDoer(t *testing.T) {
	var exchangeCalls atomic.Int32
	httpDoer := dynamicallyNonComparableHTTPDoer{
		state: []string{"not-comparable"},
		doer: nonComparableHTTPDoer{do: func(req *http.Request) (*http.Response, error) {
			if req.URL.Host == "mtls.auth.openai.com" {
				exchangeCalls.Add(1)
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     http.Header{"Content-Type": []string{"application/json"}},
					Body:       io.NopCloser(strings.NewReader(`{"access_token":"x509-token","expires_in":3600}`)),
				}, nil
			}
			return modelsListResponse(), nil
		}},
	}
	client := openai.NewClient(
		option.WithX509WorkloadIdentity(clientX509WorkloadIdentity()),
		option.WithHTTPClient(httpDoer),
	)

	for range 2 {
		if _, err := client.Models.List(t.Context()); err != nil {
			t.Fatalf("Models.List() error = %v", err)
		}
	}
	if got, want := exchangeCalls.Load(), int32(1); got != want {
		t.Errorf("token exchange calls = %d, want %d", got, want)
	}
}

func TestClientX509WorkloadIdentityCachesComparableHTTPDoerAcrossOptions(t *testing.T) {
	var exchangeCalls atomic.Int32
	var apiCalls atomic.Int32
	httpDoer := &nonComparableHTTPDoer{do: func(req *http.Request) (*http.Response, error) {
		if req.URL.Host == "mtls.auth.openai.com" {
			exchangeCalls.Add(1)
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       io.NopCloser(strings.NewReader(`{"access_token":"x509-token","expires_in":3600}`)),
			}, nil
		}
		apiCalls.Add(1)
		return modelsListResponse(), nil
	}}
	client := openai.NewClient(option.WithX509WorkloadIdentity(clientX509WorkloadIdentity()))

	if _, err := client.Models.List(t.Context(), option.WithHTTPClient(httpDoer)); err != nil {
		t.Fatalf("first Models.List() error = %v", err)
	}
	if _, err := client.Models.List(t.Context(), option.WithHTTPClient(httpDoer)); err != nil {
		t.Fatalf("second Models.List() error = %v", err)
	}

	if got, want := exchangeCalls.Load(), int32(1); got != want {
		t.Errorf("token exchange calls = %d, want %d", got, want)
	}
	if got, want := apiCalls.Load(), int32(2); got != want {
		t.Errorf("API calls = %d, want %d", got, want)
	}
}

func TestClientX509WorkloadIdentityLaterOptionReplacesCredential(t *testing.T) {
	var exchangedServiceAccounts []string
	var apiAuthorization string
	httpClient := &http.Client{Transport: &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
		if req.URL.Host == "mtls.auth.openai.com" {
			var body struct {
				ServiceAccountID string `json:"service_account_id"`
			}
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				return nil, err
			}
			exchangedServiceAccounts = append(exchangedServiceAccounts, body.ServiceAccountID)
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body: io.NopCloser(strings.NewReader(fmt.Sprintf(
					`{"access_token":"token-%s","expires_in":3600}`,
					body.ServiceAccountID,
				))),
			}, nil
		}
		apiAuthorization = req.Header.Get("Authorization")
		return modelsListResponse(), nil
	}}}
	first := clientX509WorkloadIdentity()
	first.ServiceAccountID = "svc-first"
	second := clientX509WorkloadIdentity()
	second.ServiceAccountID = "svc-second"
	client := openai.NewClient(
		option.WithX509WorkloadIdentity(first),
		option.WithHTTPClient(httpClient),
	)

	if _, err := client.Models.List(t.Context(), option.WithX509WorkloadIdentity(second)); err != nil {
		t.Fatalf("Models.List() error = %v", err)
	}
	if got, want := strings.Join(exchangedServiceAccounts, ","), "svc-second"; got != want {
		t.Errorf("exchanged service accounts = %q, want %q", got, want)
	}
	if got, want := apiAuthorization, "Bearer token-svc-second"; got != want {
		t.Errorf("API Authorization = %q, want %q", got, want)
	}
}

func TestClientX509WorkloadIdentityExchangeFailureDoesNotEnterAPIRetryLoop(t *testing.T) {
	testCases := []struct {
		name      string
		status    int
		wantCalls int
	}{
		{name: "OAuth failure", status: http.StatusForbidden, wantCalls: 1},
		{name: "exhausted transient failure", status: http.StatusServiceUnavailable, wantCalls: 3},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			exchangeCalls := 0
			httpClient := &http.Client{Transport: &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
				if req.URL.Host != "mtls.auth.openai.com" {
					t.Fatalf("unexpected API request %s", req.URL)
				}
				exchangeCalls++
				return &http.Response{
					StatusCode: testCase.status,
					Header:     http.Header{"Retry-After": []string{"0"}},
					Body:       io.NopCloser(strings.NewReader(`{"error":"exchange_failed"}`)),
				}, nil
			}}}
			client := openai.NewClient(
				option.WithX509WorkloadIdentity(clientX509WorkloadIdentity()),
				option.WithHTTPClient(httpClient),
			)

			_, err := client.Models.List(t.Context())
			if err == nil {
				t.Fatal("Models.List() error = nil")
			}
			if testCase.status == http.StatusForbidden {
				var oauthErr *auth.OAuthError
				if !errors.As(err, &oauthErr) {
					t.Fatalf("Models.List() error = %v, want *auth.OAuthError", err)
				}
			}
			if exchangeCalls != testCase.wantCalls {
				t.Errorf("token exchange calls = %d, want %d", exchangeCalls, testCase.wantCalls)
			}
		})
	}
}

func TestClientX509WorkloadIdentity401ReplaysReplayableBodyOnce(t *testing.T) {
	var exchangeCalls int
	var apiBodies []string
	var apiAuth []string
	httpClient := &http.Client{Transport: &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
		if req.URL.Host == "mtls.auth.openai.com" {
			exchangeCalls++
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body: io.NopCloser(strings.NewReader(fmt.Sprintf(
					`{"access_token":"token-%d","expires_in":3600}`,
					exchangeCalls,
				))),
			}, nil
		}
		body, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		apiBodies = append(apiBodies, string(body))
		apiAuth = append(apiAuth, req.Header.Get("Authorization"))
		if len(apiBodies) == 1 {
			return &http.Response{StatusCode: http.StatusUnauthorized, Header: make(http.Header), Body: http.NoBody}, nil
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(`{"ok":true}`)),
		}, nil
	}}}
	client := openai.NewClient(
		option.WithX509WorkloadIdentity(clientX509WorkloadIdentity()),
		option.WithBaseURL("https://api.example/v1"),
		option.WithHTTPClient(httpClient),
	)

	var result map[string]any
	err := client.Execute(
		t.Context(),
		http.MethodPost,
		"/custom",
		nil,
		&result,
		option.WithRequestBody("application/json", []byte(`{"hello":"world"}`)),
	)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if got, want := exchangeCalls, 2; got != want {
		t.Errorf("token exchange calls = %d, want %d", got, want)
	}
	if got, want := len(apiBodies), 2; got != want {
		t.Fatalf("API calls = %d, want %d", got, want)
	}
	if apiBodies[0] != apiBodies[1] {
		t.Errorf("replayed body = %q, want %q", apiBodies[1], apiBodies[0])
	}
	if got, want := strings.Join(apiAuth, ","), "Bearer token-1,Bearer token-2"; got != want {
		t.Errorf("API Authorization sequence = %q, want %q", got, want)
	}
}

func TestClientX509WorkloadIdentity401ReappliesBodyMiddleware(t *testing.T) {
	var exchangeCalls int
	var middlewareCalls int
	var apiBodies []string
	httpClient := &http.Client{Transport: &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
		if req.URL.Host == "mtls.auth.openai.com" {
			exchangeCalls++
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body: io.NopCloser(strings.NewReader(fmt.Sprintf(
					`{"access_token":"token-%d","expires_in":3600}`,
					exchangeCalls,
				))),
			}, nil
		}
		if got, want := req.Header.Get("Content-Encoding"), "gzip"; got != want {
			return nil, fmt.Errorf("Content-Encoding = %q, want %q", got, want)
		}
		reader, err := gzip.NewReader(req.Body)
		if err != nil {
			return nil, err
		}
		body, err := io.ReadAll(reader)
		if closeErr := reader.Close(); err == nil {
			err = closeErr
		}
		if err != nil {
			return nil, err
		}
		apiBodies = append(apiBodies, string(body))
		if len(apiBodies) == 1 {
			return &http.Response{StatusCode: http.StatusUnauthorized, Header: make(http.Header), Body: http.NoBody}, nil
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(`{"ok":true}`)),
		}, nil
	}}}
	client := openai.NewClient(
		option.WithX509WorkloadIdentity(clientX509WorkloadIdentity()),
		option.WithBaseURL("https://api.example/v1"),
		option.WithHTTPClient(httpClient),
		option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
			middlewareCalls++
			body, err := io.ReadAll(req.Body)
			_ = req.Body.Close()
			if err != nil {
				return nil, err
			}
			var compressed bytes.Buffer
			writer := gzip.NewWriter(&compressed)
			if _, err := writer.Write(body); err != nil {
				return nil, err
			}
			if err := writer.Close(); err != nil {
				return nil, err
			}
			req.Body = io.NopCloser(bytes.NewReader(compressed.Bytes()))
			req.ContentLength = int64(compressed.Len())
			req.Header.Set("Content-Encoding", "gzip")
			return next(req)
		}),
	)

	const requestBody = `{"hello":"world"}`
	var result map[string]any
	err := client.Execute(
		t.Context(),
		http.MethodPost,
		"/custom",
		nil,
		&result,
		option.WithRequestBody("application/json", []byte(requestBody)),
	)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if got, want := middlewareCalls, 2; got != want {
		t.Errorf("middleware calls = %d, want %d", got, want)
	}
	if got, want := exchangeCalls, 2; got != want {
		t.Errorf("token exchange calls = %d, want %d", got, want)
	}
	if got, want := strings.Join(apiBodies, ","), requestBody+","+requestBody; got != want {
		t.Errorf("API bodies = %q, want %q", got, want)
	}
}

func TestClientX509WorkloadIdentity401ReplaysFromPreMiddlewareState(t *testing.T) {
	var exchangeCalls int
	var middlewareCalls int
	var queryValueCounts []int
	var signatureCounts []int
	httpClient := &http.Client{Transport: &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
		if req.URL.Host == "mtls.auth.openai.com" {
			exchangeCalls++
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body: io.NopCloser(strings.NewReader(fmt.Sprintf(
					`{"access_token":"token-%d","expires_in":3600}`,
					exchangeCalls,
				))),
			}, nil
		}
		queryValueCounts = append(queryValueCounts, len(req.URL.Query()["signed"]))
		signatureCounts = append(signatureCounts, len(req.Header.Values("X-Test-Signature")))
		if len(queryValueCounts) == 1 {
			return &http.Response{StatusCode: http.StatusUnauthorized, Header: make(http.Header), Body: http.NoBody}, nil
		}
		return modelsListResponse(), nil
	}}}
	client := openai.NewClient(
		option.WithX509WorkloadIdentity(clientX509WorkloadIdentity()),
		option.WithHTTPClient(httpClient),
		option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
			middlewareCalls++
			query := req.URL.Query()
			query.Add("signed", "true")
			req.URL.RawQuery = query.Encode()
			req.Header.Add("X-Test-Signature", "signature")
			return next(req)
		}),
	)

	if _, err := client.Models.List(t.Context()); err != nil {
		t.Fatalf("Models.List() error = %v", err)
	}
	if got, want := middlewareCalls, 2; got != want {
		t.Errorf("middleware calls = %d, want %d", got, want)
	}
	if got, want := exchangeCalls, 2; got != want {
		t.Errorf("token exchange calls = %d, want %d", got, want)
	}
	if got, want := fmt.Sprint(queryValueCounts), "[1 1]"; got != want {
		t.Errorf("query value counts = %s, want %s", got, want)
	}
	if got, want := fmt.Sprint(signatureCounts), "[1 1]"; got != want {
		t.Errorf("signature counts = %s, want %s", got, want)
	}
}

func TestClientX509WorkloadIdentity401DoesNotReplayStreamingBody(t *testing.T) {
	var exchangeCalls int
	var apiCalls int
	httpClient := &http.Client{Transport: &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
		if req.URL.Host == "mtls.auth.openai.com" {
			exchangeCalls++
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       io.NopCloser(strings.NewReader(`{"access_token":"token-1","expires_in":3600}`)),
			}, nil
		}
		apiCalls++
		return &http.Response{StatusCode: http.StatusUnauthorized, Header: make(http.Header), Body: http.NoBody}, nil
	}}}
	client := openai.NewClient(
		option.WithX509WorkloadIdentity(clientX509WorkloadIdentity()),
		option.WithBaseURL("https://api.example/v1"),
		option.WithHTTPClient(httpClient),
	)

	err := client.Execute(
		context.Background(),
		http.MethodPost,
		"/custom",
		nil,
		nil,
		option.WithRequestBody("application/json", strings.NewReader(`{"hello":"world"}`)),
	)
	if err == nil {
		t.Fatal("Execute() error = nil")
	}
	if got, want := exchangeCalls, 1; got != want {
		t.Errorf("token exchange calls = %d, want %d", got, want)
	}
	if got, want := apiCalls, 1; got != want {
		t.Errorf("API calls = %d, want %d", got, want)
	}
}

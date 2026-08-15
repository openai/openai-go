package openai_test

import (
	"context"
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
	if _, err := client.Models.List(t.Context()); err == nil {
		t.Fatal("Models.List() error = nil")
	}
	if got := redirectedRequests.Load(); got != 0 {
		t.Errorf("redirect target requests = %d, want 0", got)
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

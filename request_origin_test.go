package openai_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	azpolicy "github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/auth"
	"github.com/openai/openai-go/v3/azure"
	"github.com/openai/openai-go/v3/bedrock"
	"github.com/openai/openai-go/v3/option"
)

type originTestRoundTripper func(*http.Request) (*http.Response, error)

func (f originTestRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

type originTestCloseTrackingBody struct {
	io.ReadCloser
	closes *atomic.Int64
}

func (b originTestCloseTrackingBody) Close() error {
	b.closes.Add(1)
	return b.ReadCloser.Close()
}

type originTestSubjectTokenProvider struct{}

func (originTestSubjectTokenProvider) TokenType() auth.SubjectTokenType {
	return auth.SubjectTokenTypeJWT
}

func (originTestSubjectTokenProvider) GetToken(context.Context, auth.HTTPDoer) (string, error) {
	return "subject-token", nil
}

type originTestAzureCredential struct{}

func (originTestAzureCredential) GetToken(context.Context, azpolicy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{Token: "azure-token", ExpiresOn: time.Now().Add(time.Hour)}, nil
}

func clearOriginTestEnvironment(t *testing.T) {
	t.Helper()
	for _, key := range []string{
		"OPENAI_API_KEY",
		"OPENAI_ADMIN_KEY",
		"OPENAI_BASE_URL",
		"OPENAI_CUSTOM_HEADERS",
		"OPENAI_ORG_ID",
		"OPENAI_PROJECT_ID",
		"OPENAI_WEBHOOK_SECRET",
	} {
		t.Setenv(key, "")
	}
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

type originTestClientFactory func(string, *http.Client) (openai.Client, error)

func originTestClientFactories() []struct {
	name             string
	credentialHeader string
	newClient        originTestClientFactory
} {
	return []struct {
		name             string
		credentialHeader string
		newClient        originTestClientFactory
	}{
		{
			name:             "API key",
			credentialHeader: "Authorization",
			newClient: func(baseURL string, httpClient *http.Client) (openai.Client, error) {
				return openai.NewClient(
					option.WithBaseURL(baseURL+"/v1"),
					option.WithAPIKey("api-key"),
					option.WithHTTPClient(httpClient),
					option.WithMaxRetries(0),
				), nil
			},
		},
		{
			name:             "admin API key",
			credentialHeader: "Authorization",
			newClient: func(baseURL string, httpClient *http.Client) (openai.Client, error) {
				return openai.NewClient(
					option.WithBaseURL(baseURL+"/v1"),
					option.WithAdminAPIKey("admin-key"),
					option.WithHTTPClient(httpClient),
					option.WithMaxRetries(0),
				), nil
			},
		},
		{
			name:             "workload identity",
			credentialHeader: "Authorization",
			newClient: func(baseURL string, httpClient *http.Client) (openai.Client, error) {
				return openai.NewClient(
					option.WithBaseURL(baseURL+"/v1"),
					option.WithWorkloadIdentity(auth.WorkloadIdentity{
						IdentityProviderID: "provider",
						ServiceAccountID:   "service-account",
						Provider:           originTestSubjectTokenProvider{},
					}),
					option.WithHTTPClient(httpClient),
					option.WithMaxRetries(0),
				), nil
			},
		},
		{
			name:             "Azure API key",
			credentialHeader: "Api-Key",
			newClient: func(baseURL string, httpClient *http.Client) (openai.Client, error) {
				return openai.NewClient(
					azure.WithEndpoint(baseURL, "2026-08-01"),
					azure.WithAPIKey("azure-key"),
					option.WithHTTPClient(httpClient),
					option.WithMaxRetries(0),
				), nil
			},
		},
		{
			name:             "Azure token credential",
			credentialHeader: "Authorization",
			newClient: func(baseURL string, httpClient *http.Client) (openai.Client, error) {
				return openai.NewClient(
					azure.WithEndpoint(baseURL, "2026-08-01"),
					azure.WithTokenCredential(originTestAzureCredential{}),
					option.WithHTTPClient(httpClient),
					option.WithMaxRetries(0),
				), nil
			},
		},
		{
			name:             "Bedrock gateway API key",
			credentialHeader: "Authorization",
			newClient: func(baseURL string, httpClient *http.Client) (openai.Client, error) {
				return bedrock.NewClient(context.Background(), bedrock.Config{
					BaseURL:  baseURL + "/openai/v1",
					SkipAuth: true,
				}, option.WithAPIKey("gateway-key"), option.WithHTTPClient(httpClient), option.WithMaxRetries(0))
			},
		},
		{
			name:             "Bedrock gateway admin API key",
			credentialHeader: "Authorization",
			newClient: func(baseURL string, httpClient *http.Client) (openai.Client, error) {
				return bedrock.NewClient(context.Background(), bedrock.Config{
					BaseURL:  baseURL + "/openai/v1",
					SkipAuth: true,
				}, option.WithAdminAPIKey("gateway-admin-key"), option.WithHTTPClient(httpClient), option.WithMaxRetries(0))
			},
		},
	}
}

func TestClientRejectsNonRelativeRequestReferences(t *testing.T) {
	clearOriginTestEnvironment(t)

	var transportCalls atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		transportCalls.Add(1)
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{}`)
	}))
	defer server.Close()

	client := openai.NewClient(
		option.WithBaseURL(server.URL+"/v1"),
		option.WithAPIKey("api-key"),
		option.WithHTTPClient(originTestHTTPClient(server.Client().Transport)),
		option.WithMaxRetries(0),
	)
	serverURL, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range []struct {
		name string
		path string
	}{
		{name: "cross-origin absolute", path: "https://other.example/v1/responses"},
		{name: "same-origin absolute", path: server.URL + "/v1/responses"},
		{name: "network path", path: "//" + serverURL.Host + "/v1/responses"},
		{name: "network path with user info", path: "//user@" + serverURL.Host + "/v1/responses"},
		{name: "network path without authority", path: "//"},
		{name: "single-slash-prefixed absolute", path: "/https://other.example/v1/responses"},
		{name: "slash-prefixed absolute", path: "///https://other.example/v1/responses"},
		{name: "opaque absolute", path: "https:responses"},
	} {
		t.Run(test.name, func(t *testing.T) {
			err := client.Post(context.Background(), test.path, strings.NewReader("private payload"), nil)
			if err == nil || !strings.Contains(err.Error(), "request path must be a relative URL reference") {
				t.Fatalf("Post() error = %v, want relative-reference error", err)
			}
		})
	}

	if got := transportCalls.Load(); got != 0 {
		t.Fatalf("transport calls = %d, want 0", got)
	}
}

func TestClientEnforcesCredentialOriginAfterRouting(t *testing.T) {
	clearOriginTestEnvironment(t)

	for _, test := range originTestClientFactories() {
		t.Run(test.name, func(t *testing.T) {
			t.Run("middleware reroute", func(t *testing.T) {
				var trustedCalls atomic.Int64
				trusted := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					trustedCalls.Add(1)
					w.Header().Set("Content-Type", "application/json")
					_, _ = io.WriteString(w, `{}`)
				}))
				defer trusted.Close()

				var otherOriginCalls atomic.Int64
				otherOrigin := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					otherOriginCalls.Add(1)
					_, _ = io.Copy(io.Discard, req.Body)
					w.Header().Set("Content-Type", "application/json")
					_, _ = io.WriteString(w, `{}`)
				}))
				defer otherOrigin.Close()

				client, err := test.newClient(trusted.URL, originTestHTTPClient(trusted.Client().Transport))
				if err != nil {
					t.Fatal(err)
				}
				target, err := url.Parse(otherOrigin.URL + "/capture")
				if err != nil {
					t.Fatal(err)
				}

				var sawCredential bool
				var sawBody bool
				err = client.Post(
					context.Background(),
					"responses",
					strings.NewReader("private payload"),
					nil,
					option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
						sawCredential = req.Header.Get(test.credentialHeader) != ""
						sawBody = req.Body != nil
						targetCopy := *target
						req.URL = &targetCopy
						return next(req)
					}),
				)
				if err == nil || !strings.Contains(err.Error(), "request URL origin must match the configured base URL") {
					t.Fatalf("Post() error = %v, want configured-origin error", err)
				}
				if !sawCredential {
					t.Fatalf("%s was not present before the transport origin check", test.credentialHeader)
				}
				if !sawBody {
					t.Fatal("request body was not present before the transport origin check")
				}
				if got := trustedCalls.Load(); got != 0 {
					t.Fatalf("trusted-origin calls = %d, want 0", got)
				}
				if got := otherOriginCalls.Load(); got != 0 {
					t.Fatalf("other-origin calls = %d, want 0", got)
				}
			})

			t.Run("redirect", func(t *testing.T) {
				var otherOriginCalls atomic.Int64
				otherOrigin := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					otherOriginCalls.Add(1)
					_, _ = io.Copy(io.Discard, req.Body)
					w.Header().Set("Content-Type", "application/json")
					_, _ = io.WriteString(w, `{}`)
				}))
				defer otherOrigin.Close()

				var trustedCalls atomic.Int64
				trusted := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					trustedCalls.Add(1)
					http.Redirect(w, req, otherOrigin.URL+"/capture", http.StatusTemporaryRedirect)
				}))
				defer trusted.Close()

				client, err := test.newClient(trusted.URL, originTestHTTPClient(trusted.Client().Transport))
				if err != nil {
					t.Fatal(err)
				}
				var redirectBodies, redirectBodyCloses atomic.Int64
				err = client.Post(
					context.Background(),
					"responses",
					map[string]string{"input": "private payload"},
					nil,
					option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
						getBody := req.GetBody
						req.GetBody = func() (io.ReadCloser, error) {
							body, getBodyErr := getBody()
							if getBodyErr != nil {
								return nil, getBodyErr
							}
							redirectBodies.Add(1)
							return originTestCloseTrackingBody{ReadCloser: body, closes: &redirectBodyCloses}, nil
						}
						return next(req)
					}),
				)
				if err == nil || !strings.Contains(err.Error(), "request URL origin must match the configured base URL") {
					t.Fatalf("Post() error = %v, want configured-origin error", err)
				}
				if got := trustedCalls.Load(); got != 1 {
					t.Fatalf("trusted-origin calls = %d, want 1", got)
				}
				if got := otherOriginCalls.Load(); got != 0 {
					t.Fatalf("other-origin calls = %d, want 0", got)
				}
				if got := redirectBodies.Load(); got != 1 {
					t.Fatalf("redirect bodies = %d, want 1", got)
				}
				if got := redirectBodyCloses.Load(); got != 1 {
					t.Fatalf("redirect body closes = %d, want 1", got)
				}
			})
		})
	}
}

func TestAzureTokenCredentialPreservesSDKRetries(t *testing.T) {
	clearOriginTestEnvironment(t)

	var calls atomic.Int64
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Header.Get("Authorization") == "" {
			t.Error("Authorization header is empty")
		}
		if calls.Add(1) == 1 {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = io.WriteString(w, `{}`)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{}`)
	}))
	defer server.Close()

	client := openai.NewClient(
		azure.WithEndpoint(server.URL, "2026-08-01"),
		azure.WithTokenCredential(originTestAzureCredential{}),
		option.WithHTTPClient(server.Client()),
		option.WithMaxRetries(1),
	)
	var response map[string]any
	if err := client.Get(context.Background(), "models", nil, &response); err != nil {
		t.Fatal(err)
	}
	if got := calls.Load(); got != 2 {
		t.Fatalf("transport calls = %d, want 2", got)
	}
}

func TestClientPreservesRelativeAndGeneratedPaths(t *testing.T) {
	clearOriginTestEnvironment(t)

	var paths []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		paths = append(paths, req.URL.EscapedPath())
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"id":"model","created":0,"object":"model","owned_by":"openai","ok":true}`)
	}))
	defer server.Close()

	client := openai.NewClient(
		option.WithBaseURL(server.URL+"/v1"),
		option.WithAPIKey("api-key"),
		option.WithHTTPClient(server.Client()),
		option.WithMaxRetries(0),
	)
	var response map[string]any
	if err := client.Post(context.Background(), "responses", map[string]string{"input": "hello"}, &response); err != nil {
		t.Fatal(err)
	}
	if err := client.Post(context.Background(), "/responses", map[string]string{"input": "hello"}, &response); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Models.Get(context.Background(), "https://other.example/model?revision=1#fragment"); err != nil {
		t.Fatal(err)
	}

	want := []string{
		"/v1/responses",
		"/v1/responses",
		"/v1/models/https:%2F%2Fother.example%2Fmodel%3Frevision=1%23fragment",
	}
	if len(paths) != len(want) {
		t.Fatalf("paths = %q, want %q", paths, want)
	}
	for i := range want {
		if paths[i] != want[i] {
			t.Fatalf("path %d = %q, want %q", i, paths[i], want[i])
		}
	}
}

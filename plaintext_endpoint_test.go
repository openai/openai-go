package openai_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/auth"
	"github.com/openai/openai-go/v3/option"
)

func TestCredentialsAllowConfiguredHTTPEndpoint(t *testing.T) {
	for _, source := range []string{"client", "environment", "request"} {
		for _, credential := range []string{"api", "admin", "workload"} {
			t.Run(source+"/"+credential, func(t *testing.T) {
				t.Setenv("OPENAI_API_KEY", "")
				t.Setenv("OPENAI_ADMIN_KEY", "")
				t.Setenv("OPENAI_BASE_URL", "https://api.openai.com/v1")
				provider := &mockSubjectTokenProvider{token: "synthetic-subject", tokenType: auth.SubjectTokenTypeJWT}
				apiCalls, exchangeCalls, credentialCalls := 0, 0, 0
				transport := originTestRoundTripper(func(req *http.Request) (*http.Response, error) {
					body := `{"data":[],"object":"list"}`
					if req.URL.String() == auth.TokenExchangeURL {
						exchangeCalls++
						body = `{"access_token":"synthetic-access","expires_in":3600}`
					} else {
						apiCalls++
						if req.Header.Get("Authorization") != "" {
							credentialCalls++
						}
					}
					return &http.Response{StatusCode: http.StatusOK, Header: http.Header{"Content-Type": {"application/json"}}, Body: io.NopCloser(strings.NewReader(body)), Request: req}, nil
				})
				opts := []option.RequestOption{option.WithHTTPClient(&http.Client{Transport: transport})}
				switch credential {
				case "api":
					opts = append(opts, option.WithAPIKey("synthetic-api"))
				case "admin":
					opts = append(opts, option.WithAdminAPIKey("synthetic-admin"))
				case "workload":
					opts = append(opts, option.WithWorkloadIdentity(testWorkloadIdentity(provider)))
				}
				const endpoint = "http://remote.invalid/v1"
				var requestOpts []option.RequestOption
				switch source {
				case "client":
					opts = append(opts, option.WithBaseURL(endpoint))
				case "environment":
					t.Setenv("OPENAI_BASE_URL", endpoint)
				case "request":
					requestOpts = append(requestOpts, option.WithBaseURL(endpoint))
				}
				client := openai.NewClient(opts...)
				var response map[string]any
				err := client.Get(t.Context(), "models", nil, &response, requestOpts...)
				if err != nil {
					t.Fatalf("configured HTTP endpoint failed: %v", err)
				}
				wantExchange := 0
				if credential == "workload" {
					wantExchange = 1
				}
				if apiCalls != 1 || credentialCalls != 1 || exchangeCalls != wantExchange || provider.GetCallCount() != wantExchange {
					t.Errorf("API=%d credential-bearing=%d exchange=%d provider=%d calls", apiCalls, credentialCalls, exchangeCalls, provider.GetCallCount())
				}
			})
		}
	}
}

func TestLocalHTTPUsesConfiguredClient(t *testing.T) {
	for _, legacyOptIn := range []bool{false, true} {
		name := "default"
		if legacyOptIn {
			name = "legacy opt-in"
		}
		t.Run(name, func(t *testing.T) {
			t.Setenv("OPENAI_API_KEY", "")
			t.Setenv("OPENAI_ADMIN_KEY", "")
			serverCalls := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				serverCalls++
				w.Header().Set("Content-Type", "application/json")
				_, _ = io.WriteString(w, `{"data":[],"object":"list"}`)
			}))
			defer server.Close()
			clientCalls := 0
			opts := []option.RequestOption{
				option.WithBaseURL(server.URL),
				option.WithAPIKey("synthetic-api"),
				option.WithHTTPClient(&http.Client{Transport: originTestRoundTripper(func(req *http.Request) (*http.Response, error) {
					clientCalls++
					return http.DefaultTransport.RoundTrip(req)
				})}),
			}
			if legacyOptIn {
				opts = append(opts, option.WithUnsafeAllowHTTP())
			}
			client := openai.NewClient(opts...)
			_, err := client.Models.List(t.Context())
			if err != nil {
				t.Fatal(err)
			}
			if clientCalls != 1 || serverCalls != 1 {
				t.Fatalf("client calls=%d, server calls=%d", clientCalls, serverCalls)
			}
		})
	}
}

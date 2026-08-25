package openai_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/auth"
	"github.com/openai/openai-go/v3/option"
)

func TestWorkloadIdentityUnauthorizedRecoveryRunsCompleteMiddlewareChain(t *testing.T) {
	for _, test := range []struct {
		name            string
		callerFirst     bool
		maximumRetries  int
		wantAttempts    int
		wantIssuerCalls int
	}{
		{name: "middleware before authentication", callerFirst: true,
			maximumRetries: 1, wantAttempts: 2, wantIssuerCalls: 2},
		{name: "middleware after authentication", maximumRetries: 1,
			wantAttempts: 2, wantIssuerCalls: 2},
		{name: "disabled retries", callerFirst: true,
			maximumRetries: 0, wantAttempts: 1, wantIssuerCalls: 1},
	} {
		t.Run(test.name, func(t *testing.T) {
			var apiAttempts, issuerCalls, middlewareCalls int
			httpClient := &http.Client{Transport: &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
				if req.URL.Host == "auth.openai.com" {
					issuerCalls++
					return rootWorkloadResponse(http.StatusOK,
						fmt.Sprintf(`{"access_token":"synthetic-bearer-%d","expires_in":3600}`, issuerCalls)), nil
				}
				apiAttempts++
				if got, want := req.Header.Get("X-Synthetic-Attempt"), fmt.Sprintf("attempt-%d", apiAttempts); got != want {
					t.Errorf("API attempt %d caller header=%q, want %q", apiAttempts, got, want)
				}
				if apiAttempts == 1 {
					return rootWorkloadResponse(http.StatusUnauthorized,
						`{"error":{"message":"synthetic rejected bearer"}}`), nil
				}
				return rootWorkloadResponse(http.StatusOK, `{"data":[]}`), nil
			}}}
			provider := &mockSubjectTokenProvider{
				token: "synthetic-subject", tokenType: auth.SubjectTokenTypeJWT,
			}
			caller := option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
				middlewareCalls++
				req.Header.Set("X-Synthetic-Attempt", fmt.Sprintf("attempt-%d", middlewareCalls))
				return next(req)
			})
			identity := option.WithWorkloadIdentity(testWorkloadIdentity(provider))
			opts := []option.RequestOption{
				option.WithHTTPClient(httpClient), option.WithMaxRetries(test.maximumRetries),
				option.WithMaxRetryDelay(time.Millisecond),
			}
			if test.callerFirst {
				opts = append(opts, caller, identity)
			} else {
				opts = append(opts, identity, caller)
			}
			client := openai.NewClient(opts...)
			_, err := client.Models.List(t.Context())
			if test.maximumRetries == 0 {
				if err == nil {
					t.Error("unauthorized request exceeded its zero retry budget")
				}
			} else if err != nil {
				t.Fatalf("complete unauthorized recovery: %v", err)
			}
			if apiAttempts != test.wantAttempts || middlewareCalls != test.wantAttempts ||
				issuerCalls != test.wantIssuerCalls {
				t.Errorf("API/middleware/issuer attempts=%d/%d/%d, want %d/%d/%d",
					apiAttempts, middlewareCalls, issuerCalls,
					test.wantAttempts, test.wantAttempts, test.wantIssuerCalls)
			}
		})
	}
}

func TestWorkloadIdentityRejectsRequestsWithoutOwnedRetryScope(t *testing.T) {
	for _, test := range []struct {
		name           string
		maximumRetries int
		nilRequest     bool
	}{
		{name: "caller removes retry scope with retries disabled"},
		{name: "caller removes retry scope with retries enabled", maximumRetries: 2},
		{name: "caller replaces request with nil", maximumRetries: 2, nilRequest: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			var issuerCalls, apiCalls, middlewareCalls int
			httpClient := &http.Client{Transport: &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
				if req.URL.Host == "auth.openai.com" {
					issuerCalls++
					return rootWorkloadResponse(http.StatusOK,
						`{"access_token":"synthetic-private-bearer","expires_in":3600}`), nil
				}
				apiCalls++
				return rootWorkloadResponse(http.StatusOK, `{"data":[]}`), nil
			}}}
			provider := &mockSubjectTokenProvider{
				token: "synthetic-subject", tokenType: auth.SubjectTokenTypeJWT,
			}
			client := openai.NewClient(
				option.WithHTTPClient(httpClient),
				option.WithMaxRetries(test.maximumRetries),
				option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
					middlewareCalls++
					if test.nilRequest {
						return next(nil)
					}
					return next(req.WithContext(context.Background()))
				}),
				option.WithWorkloadIdentity(testWorkloadIdentity(provider)),
			)

			_, err := client.Models.List(t.Context())
			if err == nil || !strings.Contains(err.Error(), "request-owned retry scope") {
				t.Fatalf("missing request-owned scope error=%v", err)
			}
			if issuerCalls != 0 || apiCalls != 0 || middlewareCalls != 1 {
				t.Errorf("issuer/API/middleware attempts=%d/%d/%d, want 0/0/1",
					issuerCalls, apiCalls, middlewareCalls)
			}
		})
	}
}

func TestWorkloadIdentitySendsExactlyOneTrustedAuthorizationHeader(t *testing.T) {
	var issuerCalls, apiCalls, middlewareCalls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if req.URL.Path == "/oauth/token" {
			call := issuerCalls.Add(1)
			_, _ = fmt.Fprintf(w, `{"access_token":"synthetic-trusted-bearer-%d","expires_in":3600}`, call)
			return
		}
		call := apiCalls.Add(1)
		values := req.Header.Values("Authorization")
		want := fmt.Sprintf("Bearer synthetic-trusted-bearer-%d", call)
		if len(values) != 1 || values[0] != want {
			t.Errorf("real-wire API attempt %d Authorization values=%q, want exactly %q", call, values, want)
		}
		if call == 1 {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = io.WriteString(w, `{"error":{"message":"synthetic rejected bearer"}}`)
			return
		}
		_, _ = io.WriteString(w, `{"data":[]}`)
	}))
	t.Cleanup(server.Close)
	httpClient := &http.Client{Transport: &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
		if req.URL.Host == "auth.openai.com" {
			req = req.Clone(req.Context())
			req.URL.Scheme = "http"
			req.URL.Host = server.Listener.Addr().String()
		}
		return http.DefaultTransport.RoundTrip(req)
	}}}
	provider := &mockSubjectTokenProvider{
		token: "synthetic-subject", tokenType: auth.SubjectTokenTypeJWT,
	}
	var captured *http.Response
	client := openai.NewClient(
		option.WithBaseURL(server.URL+"/v1/"),
		option.WithHTTPClient(httpClient),
		option.WithMaxRetries(1),
		option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
			middlewareCalls.Add(1)
			req.Header["Authorization"] = []string{"Bearer synthetic-canonical-attacker", "Bearer synthetic-extra-attacker"}
			req.Header["authorization"] = []string{"Bearer synthetic-lowercase-attacker"}
			req.Header["aUtHoRiZaTiOn"] = []string{"Bearer synthetic-mixed-case-attacker"}
			response, err := next(req)
			if len(req.Header["Authorization"]) != 2 || len(req.Header["authorization"]) != 1 ||
				len(req.Header["aUtHoRiZaTiOn"]) != 1 {
				t.Error("signer modified caller-owned request credentials")
			}
			if response != nil && response.Request != nil {
				for name := range response.Request.Header {
					if strings.EqualFold(strings.ReplaceAll(name, "_", "-"), "authorization") {
						t.Errorf("outer middleware observed response authorization alias %q", name)
					}
				}
			}
			return response, err
		}),
		option.WithWorkloadIdentity(testWorkloadIdentity(provider)),
	)

	_, err := client.Models.List(t.Context(), option.WithResponseInto(&captured))
	if err != nil {
		t.Fatalf("public workload authentication with ambiguous credentials: %v", err)
	}
	if issuerCalls.Load() != 2 || apiCalls.Load() != 2 || middlewareCalls.Load() != 2 {
		t.Errorf("issuer/API/middleware attempts=%d/%d/%d, want 2/2/2",
			issuerCalls.Load(), apiCalls.Load(), middlewareCalls.Load())
	}
	if captured == nil || captured.Request == nil {
		t.Fatal("public response did not preserve its unsigned request metadata")
	}
	for name := range captured.Request.Header {
		if strings.EqualFold(strings.ReplaceAll(name, "_", "-"), "authorization") {
			t.Errorf("WithResponseInto exposed authorization alias %q", name)
		}
	}
}

func TestWorkloadIdentityNeverReplaysCallerTransformedBody(t *testing.T) {
	for _, test := range []struct {
		name         string
		remove       bool
		wantAttempts int
	}{
		{name: "transformed body", wantAttempts: 1},
		{name: "explicitly removed body remains empty", remove: true, wantAttempts: 2},
	} {
		t.Run(test.name, func(t *testing.T) {
			var apiAttempts, middlewareCalls int
			httpClient := &http.Client{Transport: &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
				if req.URL.Host == "auth.openai.com" {
					return rootWorkloadResponse(http.StatusOK,
						`{"access_token":"synthetic-bearer","expires_in":3600}`), nil
				}
				apiAttempts++
				var payload []byte
				if req.Body != nil {
					var readErr error
					payload, readErr = io.ReadAll(req.Body)
					if readErr != nil {
						t.Fatalf("read caller-transformed request body: %v", readErr)
					}
				}
				if test.remove && len(payload) != 0 {
					t.Error("caller-removed payload was restored on an outer retry")
				}
				if !test.remove && string(payload) != "synthetic-transformed" {
					t.Errorf("caller-transformed payload=%q", payload)
				}
				if apiAttempts == 1 {
					return rootWorkloadResponse(http.StatusUnauthorized,
						`{"error":{"message":"synthetic rejected bearer"}}`), nil
				}
				return rootWorkloadResponse(http.StatusOK, `{"ok":true}`), nil
			}}}
			provider := &mockSubjectTokenProvider{
				token: "synthetic-subject", tokenType: auth.SubjectTokenTypeJWT,
			}
			client := openai.NewClient(
				option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
					middlewareCalls++
					if test.remove {
						req.Body = http.NoBody
						req.ContentLength = 0
					} else {
						req.Body = io.NopCloser(strings.NewReader("synthetic-transformed"))
						req.ContentLength = int64(len("synthetic-transformed"))
					}
					return next(req)
				}),
				option.WithWorkloadIdentity(testWorkloadIdentity(provider)),
				option.WithHTTPClient(httpClient),
				option.WithMaxRetries(1),
				option.WithMaxRetryDelay(time.Millisecond),
			)
			var result map[string]any
			err := client.Execute(t.Context(), http.MethodPost, "/responses", nil, &result,
				option.WithRequestBody("application/json", []byte("synthetic-original")))
			if test.remove && err != nil {
				t.Errorf("removed body did not recover through a fresh middleware attempt: %v", err)
			}
			if !test.remove && err == nil {
				t.Error("transformed caller body was retried without proving its replay safety")
			}
			if apiAttempts != test.wantAttempts || middlewareCalls != test.wantAttempts {
				t.Errorf("API/middleware attempts=%d/%d, want %d/%d",
					apiAttempts, middlewareCalls, test.wantAttempts, test.wantAttempts)
			}
		})
	}
}

func TestWorkloadIdentityNeverExposesSignedResponseMetadata(t *testing.T) {
	for _, test := range []struct {
		name          string
		status        int
		unauthorized  bool
		wantAPICalls  int32
		wantClientErr bool
	}{
		{name: "successful response", status: http.StatusOK, wantAPICalls: 1},
		{name: "API error response", status: http.StatusForbidden, wantAPICalls: 1, wantClientErr: true},
		{name: "unauthorized outer retry", status: http.StatusOK, unauthorized: true, wantAPICalls: 2},
	} {
		t.Run(test.name, func(t *testing.T) {
			var issuerCalls, apiCalls, observed atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				if req.URL.Path == "/oauth/token" {
					call := issuerCalls.Add(1)
					_, _ = fmt.Fprintf(w, `{"access_token":"synthetic-private-bearer-%d","expires_in":3600}`, call)
					return
				}
				call := apiCalls.Add(1)
				if got, want := req.Header.Get("Authorization"),
					fmt.Sprintf("Bearer synthetic-private-bearer-%d", call); got != want {
					t.Errorf("real-wire API attempt %d Authorization=%q, want %q", call, got, want)
				}
				status := test.status
				if test.unauthorized && call == 1 {
					status = http.StatusUnauthorized
				}
				w.WriteHeader(status)
				if status >= http.StatusBadRequest {
					_, _ = io.WriteString(w, `{"error":{"message":"synthetic API failure"}}`)
					return
				}
				_, _ = io.WriteString(w, `{"data":[]}`)
			}))
			t.Cleanup(server.Close)
			httpClient := &http.Client{Transport: &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
				if req.URL.Host == "auth.openai.com" {
					req = req.Clone(req.Context())
					req.URL.Scheme = "http"
					req.URL.Host = server.Listener.Addr().String()
				}
				return http.DefaultTransport.RoundTrip(req)
			}}}
			provider := &mockSubjectTokenProvider{
				token: "synthetic-subject", tokenType: auth.SubjectTokenTypeJWT,
			}
			var captured *http.Response
			client := openai.NewClient(
				option.WithBaseURL(server.URL+"/v1/"),
				option.WithHTTPClient(httpClient),
				option.WithMaxRetries(1),
				option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
					if got := req.Header.Get("Authorization"); got != "" {
						t.Errorf("outer caller request exposed credentials before authentication: %q", got)
					}
					response, err := next(req)
					if got := req.Header.Get("Authorization"); got != "" {
						t.Errorf("ordinary signer modified caller-owned request headers: %q", got)
					}
					if response != nil && response.Request != nil {
						observed.Add(1)
						if got := response.Request.Header.Get("Authorization"); got != "" {
							t.Errorf("outer caller observed response bearer credentials: %q", got)
						}
					}
					return response, err
				}),
				option.WithWorkloadIdentity(testWorkloadIdentity(provider)),
			)
			_, err := client.Models.List(t.Context(), option.WithResponseInto(&captured))
			if (err != nil) != test.wantClientErr {
				t.Errorf("public workload request error=%v, want error=%t", err, test.wantClientErr)
			}
			if captured == nil || captured.Request == nil {
				t.Fatal("WithResponseInto did not preserve its response and request metadata")
			}
			if got := captured.Request.Header.Get("Authorization"); got != "" {
				t.Errorf("WithResponseInto exposed private bearer credentials: %q", got)
			}
			if got := apiCalls.Load(); got != test.wantAPICalls {
				t.Errorf("real-wire API requests=%d, want %d", got, test.wantAPICalls)
			}
			if got := observed.Load(); got != test.wantAPICalls {
				t.Errorf("outer caller response observations=%d, want %d", got, test.wantAPICalls)
			}
		})
	}
}

func TestWorkloadIdentityRetriesOnlyTransientIssuerOAuthFailures(t *testing.T) {
	for _, test := range []struct {
		name            string
		status          int
		code            string
		transient       bool
		wantIssuerCalls int32
	}{
		{name: "permanent bad request", status: http.StatusBadRequest,
			code: "invalid_grant", wantIssuerCalls: 1},
		{name: "permanent unauthorized", status: http.StatusUnauthorized,
			code: "invalid_client", wantIssuerCalls: 1},
		{name: "permanent forbidden", status: http.StatusForbidden,
			code: "access_denied", wantIssuerCalls: 1},
		{name: "unknown OAuth rejection", status: http.StatusBadRequest,
			code: "synthetic-unknown-code", wantIssuerCalls: 1},
		{name: "temporarily unavailable", status: http.StatusBadRequest,
			code: "temporarily_unavailable", transient: true, wantIssuerCalls: 2},
		{name: "temporary server error", status: http.StatusForbidden,
			code: "server_error", transient: true, wantIssuerCalls: 2},
	} {
		t.Run(test.name, func(t *testing.T) {
			var issuerCalls, apiCalls atomic.Int32
			httpClient := &http.Client{Transport: &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
				if req.URL.Host == "auth.openai.com" {
					call := issuerCalls.Add(1)
					if !test.transient || call == 1 {
						return rootWorkloadResponse(test.status,
							fmt.Sprintf(`{"error":%q,"error_description":"synthetic-private-description"}`, test.code)), nil
					}
					return rootWorkloadResponse(http.StatusOK,
						`{"access_token":"synthetic-refreshed-bearer","expires_in":3600}`), nil
				}
				apiCalls.Add(1)
				return rootWorkloadResponse(http.StatusOK, `{"data":[]}`), nil
			}}}
			provider := &mockSubjectTokenProvider{
				token: "synthetic-subject", tokenType: auth.SubjectTokenTypeJWT,
			}
			client := openai.NewClient(
				option.WithWorkloadIdentity(testWorkloadIdentity(provider)),
				option.WithHTTPClient(httpClient),
				option.WithMaxRetries(2),
				option.WithMaxRetryDelay(time.Millisecond),
			)
			_, err := client.Models.List(t.Context())
			if test.transient {
				if err != nil || apiCalls.Load() != 1 {
					t.Errorf("transient OAuth recovery error=%v API requests=%d", err, apiCalls.Load())
				}
			} else {
				var oauthError *auth.OAuthError
				if !errors.As(err, &oauthError) || oauthError.StatusCode != test.status {
					t.Errorf("permanent OAuth error lost its typed public status: %v", err)
				}
				if err != nil && strings.Contains(err.Error(), "synthetic-private-description") {
					t.Error("permanent OAuth error exposed its issuer description")
				}
				if apiCalls.Load() != 0 {
					t.Errorf("permanently rejected workload dispatched %d API requests", apiCalls.Load())
				}
			}
			if got := issuerCalls.Load(); got != test.wantIssuerCalls {
				t.Errorf("issuer exchanges=%d, want %d", got, test.wantIssuerCalls)
			}
		})
	}
}

func rootWorkloadResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

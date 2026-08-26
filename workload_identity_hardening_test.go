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

func TestWorkloadIdentityPreservesExplicitCredentialOptionPrecedence(t *testing.T) {
	deleted := option.WithHeaderDel("Authorization")
	explicit := option.WithHeader("Authorization", "Bearer synthetic-explicit-credential")
	apiKey := option.WithAPIKey("synthetic-explicit-api-key")
	adminKey := option.WithAdminAPIKey("synthetic-explicit-admin-key")
	for _, test := range []struct {
		name           string
		clientBefore   []option.RequestOption
		clientAfter    []option.RequestOption
		requestBefore  []option.RequestOption
		requestAfter   []option.RequestOption
		methodIdentity bool
		rejected       bool
	}{
		{name: "ambient inherited credentials are replaced"},
		{name: "earlier same-layer API key is replaced", clientBefore: []option.RequestOption{apiKey}},
		{name: "later same-layer API key conflicts", clientAfter: []option.RequestOption{apiKey}, rejected: true},
		{name: "earlier same-layer admin key conflicts", clientBefore: []option.RequestOption{adminKey}, rejected: true},
		{name: "later same-layer admin key conflicts", clientAfter: []option.RequestOption{adminKey}, rejected: true},
		{name: "earlier same-layer header deletion conflicts", clientBefore: []option.RequestOption{deleted}, rejected: true},
		{name: "later same-layer header deletion conflicts", clientAfter: []option.RequestOption{deleted}, rejected: true},
		{name: "earlier same-layer explicit header conflicts", clientBefore: []option.RequestOption{explicit}, rejected: true},
		{name: "later same-layer explicit header conflicts", clientAfter: []option.RequestOption{explicit}, rejected: true},
		{name: "method header deletion conflicts", requestAfter: []option.RequestOption{deleted}, rejected: true},
		{name: "method explicit header conflicts", requestAfter: []option.RequestOption{explicit}, rejected: true},
		{name: "method API key conflicts", requestAfter: []option.RequestOption{apiKey}, rejected: true},
		{name: "method admin key conflicts", requestAfter: []option.RequestOption{adminKey}, rejected: true},
		{name: "method workload replaces inherited header deletion", clientBefore: []option.RequestOption{deleted}, methodIdentity: true},
		{name: "method workload replaces inherited explicit header", clientBefore: []option.RequestOption{explicit}, methodIdentity: true},
		{name: "method workload replaces inherited API key", clientBefore: []option.RequestOption{apiKey}, methodIdentity: true},
		{name: "method workload replaces inherited admin key", clientBefore: []option.RequestOption{adminKey}, methodIdentity: true},
		{name: "earlier method header deletion conflicts", requestBefore: []option.RequestOption{deleted}, methodIdentity: true, rejected: true},
		{name: "later method header deletion conflicts", requestAfter: []option.RequestOption{deleted}, methodIdentity: true, rejected: true},
		{name: "earlier method explicit header conflicts", requestBefore: []option.RequestOption{explicit}, methodIdentity: true, rejected: true},
		{name: "later method explicit header conflicts", requestAfter: []option.RequestOption{explicit}, methodIdentity: true, rejected: true},
		{name: "earlier method API key is replaced", requestBefore: []option.RequestOption{apiKey}, methodIdentity: true},
		{name: "later method API key conflicts", requestAfter: []option.RequestOption{apiKey}, methodIdentity: true, rejected: true},
		{name: "earlier method admin key conflicts", requestBefore: []option.RequestOption{adminKey}, methodIdentity: true, rejected: true},
		{name: "later method admin key conflicts", requestAfter: []option.RequestOption{adminKey}, methodIdentity: true, rejected: true},
		{name: "empty API key cannot erase explicit header deletion",
			requestAfter: []option.RequestOption{deleted, option.WithAPIKey("")}, rejected: true},
		{name: "unrelated explicit header remains compatible",
			requestAfter: []option.RequestOption{option.WithHeader("X-Synthetic-Metadata", "preserved")}},
		{name: "unrelated header deletion remains compatible",
			requestAfter: []option.RequestOption{option.WithHeaderDel("X-Synthetic-Metadata")}},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("OPENAI_API_KEY", "synthetic-ambient-api-key")
			t.Setenv("OPENAI_ADMIN_KEY", "synthetic-ambient-admin-key")
			t.Setenv("OPENAI_CUSTOM_HEADERS", "")
			var issuerCalls, apiCalls int
			httpClient := &http.Client{Transport: &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
				if req.URL.Host == "auth.openai.com" {
					issuerCalls++
					return rootWorkloadResponse(http.StatusOK,
						`{"access_token":"synthetic-workload-bearer","expires_in":3600}`), nil
				}
				apiCalls++
				if got := req.Header.Values("Authorization"); len(got) != 1 ||
					got[0] != "Bearer synthetic-workload-bearer" {
					t.Errorf("API Authorization values=%q, want only the selected workload bearer", got)
				}
				return rootWorkloadResponse(http.StatusOK, `{"data":[]}`), nil
			}}}
			provider := &mockSubjectTokenProvider{
				token: "synthetic-subject", tokenType: auth.SubjectTokenTypeJWT,
			}
			identity := option.WithWorkloadIdentity(testWorkloadIdentity(provider))
			clientOptions := []option.RequestOption{option.WithHTTPClient(httpClient), option.WithMaxRetries(0)}
			clientOptions = append(clientOptions, test.clientBefore...)
			if !test.methodIdentity {
				clientOptions = append(clientOptions, identity)
			}
			clientOptions = append(clientOptions, test.clientAfter...)
			client := openai.NewClient(clientOptions...)
			requestOptions := append([]option.RequestOption(nil), test.requestBefore...)
			if test.methodIdentity {
				requestOptions = append(requestOptions, identity)
			}
			requestOptions = append(requestOptions, test.requestAfter...)

			_, err := client.Models.List(t.Context(), requestOptions...)
			if test.rejected {
				if err == nil || !strings.Contains(err.Error(), "cannot be combined") {
					t.Errorf("conflicting workload credential error=%v", err)
				}
				if issuerCalls != 0 || apiCalls != 0 {
					t.Errorf("conflicting credentials reached issuer/API %d/%d times, want 0/0", issuerCalls, apiCalls)
				}
				return
			}
			if err != nil || issuerCalls != 1 || apiCalls != 1 {
				t.Errorf("selected workload authentication error=%v issuer/API=%d/%d, want nil and 1/1",
					err, issuerCalls, apiCalls)
			}
		})
	}
}

func TestWorkloadIdentityRejectsAdminOnlyOperationsBeforeTokenExchange(t *testing.T) {
	var issuerCalls, apiCalls int
	httpClient := &http.Client{Transport: &closureTransport{fn: func(request *http.Request) (*http.Response, error) {
		if request.URL.Host == "auth.openai.com" {
			issuerCalls++
			return rootWorkloadResponse(http.StatusOK,
				`{"access_token":"synthetic-workload-bearer","expires_in":3600}`), nil
		}
		apiCalls++
		return rootWorkloadResponse(http.StatusOK, `{}`), nil
	}}}
	provider := &mockSubjectTokenProvider{
		token: "synthetic-subject", tokenType: auth.SubjectTokenTypeJWT,
	}
	client := openai.NewClient(
		option.WithHTTPClient(httpClient),
		option.WithWorkloadIdentity(testWorkloadIdentity(provider)),
		option.WithMaxRetries(2),
	)

	_, err := client.Admin.Organization.DataRetention.Get(t.Context())
	if err == nil || !strings.Contains(err.Error(), "admin-only API operation") {
		t.Fatalf("ordinary workload identity authenticated an admin-only operation: %v", err)
	}
	if issuerCalls != 0 || apiCalls != 0 {
		t.Errorf("admin-only workload issuer/API requests=%d/%d, want 0/0", issuerCalls, apiCalls)
	}
}

func TestWorkloadIdentityCloudMetadataFailuresRemainPrivate(t *testing.T) {
	for _, provider := range []struct {
		name     string
		identity string
		new      func() auth.SubjectTokenProvider
	}{
		{name: "Azure", identity: "azure-imds", new: func() auth.SubjectTokenProvider {
			return auth.AzureManagedIdentityTokenProvider(nil)
		}},
		{name: "GCP", identity: "gcp-metadata", new: func() auth.SubjectTokenProvider {
			return auth.GCPIDTokenProvider(nil)
		}},
	} {
		for _, failure := range []struct {
			name      string
			redirect  bool
			oversized bool
		}{
			{name: "private unsuccessful metadata response"},
			{name: "oversized unsuccessful metadata response", oversized: true},
			{name: "cross-origin metadata redirect", redirect: true},
		} {
			t.Run(provider.name+" "+failure.name, func(t *testing.T) {
				var metadataCalls, issuerCalls, apiCalls, redirectedCalls atomic.Int32
				target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					redirectedCalls.Add(1)
					_, _ = io.WriteString(w, "synthetic-attacker-subject-token")
				}))
				t.Cleanup(target.Close)
				metadata := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					metadataCalls.Add(1)
					if failure.redirect {
						w.Header().Set("Location", target.URL+"/steal?credential=synthetic-private-metadata")
						w.WriteHeader(http.StatusTemporaryRedirect)
						return
					}
					w.WriteHeader(http.StatusUnauthorized)
					body := `{"access_token":"synthetic-private-metadata-token"}`
					if failure.oversized {
						body = strings.Repeat("a", (4<<10)+1)
					}
					_, _ = io.WriteString(w, body)
				}))
				t.Cleanup(metadata.Close)
				httpClient := &http.Client{Transport: &closureTransport{fn: func(request *http.Request) (*http.Response, error) {
					switch request.URL.Host {
					case "169.254.169.254", "metadata.google.internal":
						request = request.Clone(request.Context())
						request.URL.Scheme = "http"
						request.URL.Host = metadata.Listener.Addr().String()
						return http.DefaultTransport.RoundTrip(request)
					case "auth.openai.com":
						issuerCalls.Add(1)
						return rootWorkloadResponse(http.StatusOK,
							`{"access_token":"synthetic-workload-bearer","expires_in":3600}`), nil
					default:
						apiCalls.Add(1)
						return rootWorkloadResponse(http.StatusOK, `{"data":[]}`), nil
					}
				}}}
				retries := 0
				if failure.redirect {
					retries = 2
				}
				client := openai.NewClient(
					option.WithHTTPClient(httpClient),
					option.WithWorkloadIdentity(auth.WorkloadIdentity{
						IdentityProviderID: "synthetic-provider",
						ServiceAccountID:   "synthetic-account",
						Provider:           provider.new(),
					}),
					option.WithMaxRetries(retries),
				)

				_, err := client.Models.List(t.Context())
				var typed *auth.SubjectTokenProviderError
				if !errors.As(err, &typed) || typed.Provider != provider.identity {
					t.Fatalf("public metadata failure lost its provider identity: %v", err)
				}
				if strings.Contains(err.Error(), "synthetic-private") || strings.Contains(err.Error(), target.URL) {
					t.Errorf("public metadata failure exposed credentials: %q", err.Error())
				}
				if failure.redirect && !strings.Contains(err.Error(), "does not follow redirects") {
					t.Errorf("public metadata redirect error=%v", err)
				}
				if failure.oversized && !strings.Contains(err.Error(), "size limit") {
					t.Errorf("public oversized metadata error=%v", err)
				}
				if metadataCalls.Load() != 1 || issuerCalls.Load() != 0 || apiCalls.Load() != 0 ||
					redirectedCalls.Load() != 0 {
					t.Errorf("metadata/issuer/API/redirect requests=%d/%d/%d/%d, want 1/0/0/0",
						metadataCalls.Load(), issuerCalls.Load(), apiCalls.Load(), redirectedCalls.Load())
				}
			})
		}
	}
}

func TestWorkloadIdentityNeverFollowsIssuerRedirects(t *testing.T) {
	for _, status := range []int{
		http.StatusMovedPermanently,
		http.StatusFound,
		http.StatusSeeOther,
		http.StatusTemporaryRedirect,
		http.StatusPermanentRedirect,
	} {
		t.Run(fmt.Sprintf("HTTP %d", status), func(t *testing.T) {
			var issuerCalls, apiCalls, redirectedCalls, callerRedirectChecks atomic.Int32
			target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				redirectedCalls.Add(1)
				_, _ = io.WriteString(w, `{"access_token":"synthetic-attacker-bearer","expires_in":3600}`)
			}))
			t.Cleanup(target.Close)
			issuer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
				if request.URL.Path != "/oauth/token" {
					apiCalls.Add(1)
					_, _ = io.WriteString(w, `{"data":[]}`)
					return
				}
				issuerCalls.Add(1)
				body, err := io.ReadAll(request.Body)
				if err != nil || !strings.Contains(string(body), "synthetic-private-subject") {
					t.Errorf("original issuer subject token body=%q error=%v", body, err)
				}
				w.Header().Set("Location", target.URL+"/steal?credential=synthetic-private-destination")
				w.WriteHeader(status)
			}))
			t.Cleanup(issuer.Close)
			caller := &http.Client{
				Transport: &closureTransport{fn: func(request *http.Request) (*http.Response, error) {
					if request.URL.Host == "auth.openai.com" {
						request = request.Clone(request.Context())
						request.URL.Scheme = "http"
						request.URL.Host = issuer.Listener.Addr().String()
					}
					return http.DefaultTransport.RoundTrip(request)
				}},
				CheckRedirect: func(*http.Request, []*http.Request) error {
					callerRedirectChecks.Add(1)
					return nil
				},
			}
			provider := &mockSubjectTokenProvider{
				token: "synthetic-private-subject", tokenType: auth.SubjectTokenTypeJWT,
			}
			client := openai.NewClient(
				option.WithBaseURL(issuer.URL+"/v1/"),
				option.WithHTTPClient(caller),
				option.WithWorkloadIdentity(testWorkloadIdentity(provider)),
				option.WithMaxRetries(2),
			)

			_, err := client.Models.List(t.Context())
			if err == nil || !strings.Contains(err.Error(), "does not follow redirects") {
				t.Fatalf("redirected public token exchange error=%v", err)
			}
			if strings.Contains(err.Error(), "synthetic-private") || strings.Contains(err.Error(), target.URL) {
				t.Errorf("redirect error exposed its private destination or subject token: %q", err.Error())
			}
			if issuerCalls.Load() != 1 || apiCalls.Load() != 0 || redirectedCalls.Load() != 0 ||
				callerRedirectChecks.Load() != 0 {
				t.Errorf("issuer/API/redirect/caller-policy requests=%d/%d/%d/%d, want 1/0/0/0",
					issuerCalls.Load(), apiCalls.Load(), redirectedCalls.Load(), callerRedirectChecks.Load())
			}
			if caller.CheckRedirect == nil {
				t.Error("issuer exchange modified its caller-owned native HTTP client")
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
			var callerAuthorizationHeaders int
			for name, values := range req.Header {
				if !strings.EqualFold(strings.ReplaceAll(name, "_", "-"), "authorization") {
					continue
				}
				callerAuthorizationHeaders++
				want := 1
				if name == "Authorization" {
					want = 2
				}
				if len(values) != want {
					t.Errorf("caller-owned authorization header %q values=%q, want %d", name, values, want)
				}
			}
			if callerAuthorizationHeaders != 3 {
				t.Errorf("caller-owned authorization headers=%d, want three", callerAuthorizationHeaders)
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
			var apiAttempts, issuerCalls, middlewareCalls int
			httpClient := &http.Client{Transport: &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
				if req.URL.Host == "auth.openai.com" {
					issuerCalls++
					return rootWorkloadResponse(http.StatusOK,
						fmt.Sprintf(`{"access_token":"synthetic-bearer-%d","expires_in":3600}`, issuerCalls)), nil
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

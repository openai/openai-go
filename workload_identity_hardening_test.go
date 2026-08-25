package openai_test

import (
	"fmt"
	"io"
	"net/http"
	"strings"
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

func rootWorkloadResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

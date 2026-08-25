package auth_test

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/openai/openai-go/v3/auth"
)

func TestWorkloadIdentityMiddlewarePreservesBodylessRequests(t *testing.T) {
	for _, test := range []struct {
		name   string
		noBody bool
	}{
		{name: "nil body"},
		{name: "http.NoBody sentinel", noBody: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			identity := newOrdinaryWorkloadIdentity(t)
			var issuerCalls atomic.Int32
			httpClient := ordinaryWorkloadIssuer(t, func() string {
				return fmt.Sprintf("synthetic-bearer-%d", issuerCalls.Add(1))
			})
			request, err := http.NewRequestWithContext(t.Context(), http.MethodPost,
				"https://api.openai.com/v1/responses", nil)
			if err != nil {
				t.Fatalf("construct bodyless request: %v", err)
			}
			if test.noBody {
				request.Body = http.NoBody
			}
			var recreated atomic.Int32
			request.GetBody = func() (io.ReadCloser, error) {
				recreated.Add(1)
				return io.NopCloser(strings.NewReader("synthetic-removed-payload")), nil
			}

			var attempts int
			response, err := auth.WorkloadIdentityMiddleware(identity, httpClient, request,
				func(sent *http.Request) (*http.Response, error) {
					attempts++
					if sent.Body != nil && sent.Body != http.NoBody {
						t.Error("bodyless replay restored a removed payload")
					}
					status := http.StatusOK
					if attempts == 1 {
						status = http.StatusUnauthorized
					}
					return ordinaryWorkloadResponse(status, `{}`), nil
				})
			if err != nil || response == nil || response.StatusCode != http.StatusOK {
				t.Fatalf("bodyless authentication replay: response=%v error=%v", response, err)
			}
			if closeErr := response.Body.Close(); closeErr != nil {
				t.Fatalf("close bodyless authentication response: %v", closeErr)
			}
			if attempts != 2 || recreated.Load() != 0 || issuerCalls.Load() != 2 {
				t.Errorf("bodyless API attempts/recreated bodies/exchanges = %d/%d/%d, want 2/0/2",
					attempts, recreated.Load(), issuerCalls.Load())
			}
		})
	}
}

func TestWorkloadIdentityMiddlewareRejectsUnprovableBodyReplay(t *testing.T) {
	identity := newOrdinaryWorkloadIdentity(t)
	httpClient := ordinaryWorkloadIssuer(t, func() string { return "synthetic-bearer" })
	request, err := http.NewRequestWithContext(t.Context(), http.MethodPost,
		"https://api.openai.com/v1/responses", strings.NewReader("synthetic-original-payload"))
	if err != nil {
		t.Fatalf("construct transformed request: %v", err)
	}
	request.Body = io.NopCloser(strings.NewReader("synthetic-transformed-payload"))
	request.ContentLength = int64(len("synthetic-transformed-payload"))
	var attempts int
	response, err := auth.WorkloadIdentityMiddleware(identity, httpClient, request,
		func(sent *http.Request) (*http.Response, error) {
			attempts++
			body, readErr := io.ReadAll(sent.Body)
			if readErr != nil || string(body) != "synthetic-transformed-payload" {
				t.Errorf("transformed request body = %q, error = %v", body, readErr)
			}
			return ordinaryWorkloadResponse(http.StatusUnauthorized, `{}`), nil
		})
	if err != nil || response == nil || response.StatusCode != http.StatusUnauthorized {
		t.Fatalf("transformed body authentication response=%v error=%v", response, err)
	}
	if closeErr := response.Body.Close(); closeErr != nil {
		t.Fatalf("close transformed authentication response: %v", closeErr)
	}
	if attempts != 1 {
		t.Errorf("transformed request API attempts = %d, want exactly one", attempts)
	}
}

func TestWorkloadIdentityMiddlewareClosesOnlyPresentUnauthorizedBodies(t *testing.T) {
	identity := newOrdinaryWorkloadIdentity(t)
	httpClient := ordinaryWorkloadIssuer(t, func() string { return "synthetic-bearer" })
	request, err := http.NewRequestWithContext(t.Context(), http.MethodGet,
		"https://api.openai.com/v1/models", nil)
	if err != nil {
		t.Fatalf("construct synthetic unauthorized request: %v", err)
	}
	var attempts int
	response, err := auth.WorkloadIdentityMiddleware(identity, httpClient, request,
		func(*http.Request) (*http.Response, error) {
			attempts++
			if attempts == 1 {
				return &http.Response{StatusCode: http.StatusUnauthorized, Header: make(http.Header)}, nil
			}
			return ordinaryWorkloadResponse(http.StatusOK, `{}`), nil
		})
	if err != nil || response == nil || response.StatusCode != http.StatusOK || attempts != 2 {
		t.Fatalf("nil unauthorized response-body recovery: response=%v attempts=%d error=%v",
			response, attempts, err)
	}
	if closeErr := response.Body.Close(); closeErr != nil {
		t.Fatalf("close recovered synthetic response: %v", closeErr)
	}
}

func TestWorkloadIdentityIssuerErrorsRemainBoundedAndRedacted(t *testing.T) {
	const sensitive = "synthetic-private-issuer-bearer"
	for _, test := range []struct {
		name        string
		status      int
		body        string
		wantOAuth   bool
		wantCode    string
		wantLimited bool
	}{
		{name: "server error body", status: 500, body: `{"access_token":"` + sensitive + `"}`},
		{name: "OAuth description", status: 400,
			body:      `{"error":"invalid_grant","error_description":"` + sensitive + `"}`,
			wantOAuth: true, wantCode: "invalid_grant"},
		{name: "untrusted OAuth error code", status: 403,
			body: `{"error":"` + sensitive + `"}`, wantOAuth: true},
		{name: "missing bearer response", status: 200, body: `{"detail":"` + sensitive + `"}`},
		{name: "oversized successful response", status: 200,
			body: strings.Repeat("a", (1<<20)+1), wantLimited: true},
		{name: "oversized unsuccessful response", status: 500,
			body: strings.Repeat("a", (4<<10)+1), wantLimited: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			identity := newOrdinaryWorkloadIdentity(t)
			httpClient := mockOAuthServer(test.body, test.status)
			token, err := identity.GetToken(t.Context(), httpClient)
			if token != "" || err == nil {
				t.Fatalf("unsafe issuer response: token=%q error=%v", token, err)
			}
			if strings.Contains(err.Error(), sensitive) {
				t.Error("issuer response content escaped into the returned error")
			}
			if test.wantLimited && !strings.Contains(err.Error(), "size limit") {
				t.Errorf("oversized issuer response error = %v, want a safe limit error", err)
			}
			var oauthError *auth.OAuthError
			if got := errors.As(err, &oauthError); got != test.wantOAuth {
				t.Fatalf("typed OAuth error = %t, want %t", got, test.wantOAuth)
			}
			if oauthError != nil && (string(oauthError.ErrorCode) != test.wantCode ||
				oauthError.ErrorDescription != "") {
				t.Errorf("safe OAuth error code/description = %q/%q, want %q/empty",
					oauthError.ErrorCode, oauthError.ErrorDescription, test.wantCode)
			}
		})
	}
}

func TestWorkloadIdentityIssuerClosesOversizedAndUnreadableResponses(t *testing.T) {
	for _, test := range []struct {
		name          string
		contentLength int64
		readError     error
		wantReads     int32
	}{
		{name: "declared oversized issuer response", contentLength: (1 << 20) + 1},
		{name: "private response read failure", readError: errors.New("synthetic-private-read-cause"), wantReads: 1},
	} {
		t.Run(test.name, func(t *testing.T) {
			identity := newOrdinaryWorkloadIdentity(t)
			body := &ordinaryObservedBody{readError: test.readError}
			httpClient := &http.Client{Transport: &closureTransport{fn: func(*http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode:    http.StatusOK,
					Header:        make(http.Header),
					Body:          body,
					ContentLength: test.contentLength,
				}, nil
			}}}
			token, err := identity.GetToken(t.Context(), httpClient)
			if token != "" || err == nil || strings.Contains(err.Error(), "synthetic-private-read-cause") {
				t.Errorf("unsafe issuer read response: token=%q error=%v", token, err)
			}
			if got := body.reads.Load(); got != test.wantReads {
				t.Errorf("issuer response reads=%d, want %d", got, test.wantReads)
			}
			if got := body.closes.Load(); got != 1 {
				t.Errorf("issuer response body closes=%d, want exactly one", got)
			}
		})
	}
}

func TestWorkloadIdentityRefreshDoesNotRestoreRejectedGeneration(t *testing.T) {
	for _, test := range []struct {
		name          string
		alwaysReject  bool
		wantExchanges int32
	}{
		{name: "replacement generation", wantExchanges: 3},
		{name: "repeated rejected generation is bounded", alwaysReject: true, wantExchanges: 4},
	} {
		t.Run(test.name, func(t *testing.T) {
			identity := newOrdinaryWorkloadIdentity(t)
			const rejected = "synthetic-rejected-generation"
			refreshStarted := make(chan struct{})
			releaseRefresh := make(chan struct{})
			var exchanges atomic.Int32
			httpClient := &http.Client{Transport: &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
				call := exchanges.Add(1)
				if call == 2 {
					close(refreshStarted)
					select {
					case <-releaseRefresh:
					case <-req.Context().Done():
						return nil, req.Context().Err()
					}
				}
				token := rejected
				if call >= 3 && !test.alwaysReject {
					token = "synthetic-replacement-generation"
				}
				return ordinaryWorkloadResponse(http.StatusOK,
					fmt.Sprintf(`{"access_token":%q,"expires_in":600}`, token)), nil
			}}}
			if token, err := identity.GetToken(t.Context(), httpClient); token != rejected || err != nil {
				t.Fatalf("prime cached generation: token=%q error=%v", token, err)
			}
			if token, err := identity.GetToken(t.Context(), httpClient); token != rejected || err != nil {
				t.Fatalf("begin proactive refresh: token=%q error=%v", token, err)
			}
			select {
			case <-refreshStarted:
			case <-time.After(5 * time.Second):
				t.Fatal("proactive refresh never reached its synchronized issuer")
			}

			request, err := http.NewRequestWithContext(t.Context(), http.MethodGet,
				"https://api.openai.com/v1/models", nil)
			if err != nil {
				t.Fatalf("construct rejected-generation request: %v", err)
			}
			apiRejected := make(chan struct{})
			type outcome struct {
				header string
				err    error
			}
			completed := make(chan outcome, 1)
			go func() {
				var attempts int
				var replayed string
				response, authErr := auth.WorkloadIdentityMiddleware(identity, httpClient, request,
					func(sent *http.Request) (*http.Response, error) {
						attempts++
						if attempts == 1 {
							close(apiRejected)
							return ordinaryWorkloadResponse(http.StatusUnauthorized, `{}`), nil
						}
						replayed = sent.Header.Get("Authorization")
						return ordinaryWorkloadResponse(http.StatusOK, `{}`), nil
					})
				if response != nil && response.Body != nil {
					_ = response.Body.Close()
				}
				completed <- outcome{header: replayed, err: authErr}
			}()
			select {
			case <-apiRejected:
			case <-time.After(5 * time.Second):
				t.Fatal("API never rejected its cached generation")
			}
			close(releaseRefresh)
			select {
			case result := <-completed:
				if test.alwaysReject {
					if result.err == nil || strings.Contains(result.err.Error(), rejected) || result.header != "" {
						t.Errorf("repeated rejected bearer escaped: replay=%q error=%v", result.header, result.err)
					}
				} else if result.err != nil || result.header != "Bearer synthetic-replacement-generation" {
					t.Errorf("replacement generation replay=%q error=%v", result.header, result.err)
				}
			case <-time.After(5 * time.Second):
				t.Fatal("invalidated background refresh did not complete")
			}
			if got := exchanges.Load(); got != test.wantExchanges {
				t.Errorf("generation refresh exchanges=%d, want %d", got, test.wantExchanges)
			}
		})
	}
}

func newOrdinaryWorkloadIdentity(t *testing.T) *auth.WorkloadIdentityAuth {
	t.Helper()
	identity, err := auth.NewWorkloadIdentityAuth(auth.WorkloadIdentity{
		IdentityProviderID: "synthetic-provider",
		ServiceAccountID:   "synthetic-service-account",
		Provider:           &mockProvider{token: "synthetic-subject", tokenType: auth.SubjectTokenTypeJWT},
	})
	if err != nil {
		t.Fatalf("construct ordinary workload identity: %v", err)
	}
	return identity
}

func ordinaryWorkloadIssuer(t *testing.T, token func() string) *http.Client {
	t.Helper()
	return &http.Client{Transport: &closureTransport{fn: func(*http.Request) (*http.Response, error) {
		return ordinaryWorkloadResponse(http.StatusOK,
			fmt.Sprintf(`{"access_token":%q,"expires_in":3600}`, token())), nil
	}}}
}

func ordinaryWorkloadResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

type ordinaryObservedBody struct {
	reads     atomic.Int32
	closes    atomic.Int32
	readError error
}

func (body *ordinaryObservedBody) Read([]byte) (int, error) {
	body.reads.Add(1)
	if body.readError != nil {
		return 0, body.readError
	}
	return 0, io.EOF
}

func (body *ordinaryObservedBody) Close() error {
	body.closes.Add(1)
	return nil
}

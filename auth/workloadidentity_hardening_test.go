package auth_test

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
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

func TestWorkloadIdentityMiddlewareReplacesAllAuthorizationAliases(t *testing.T) {
	identity := newOrdinaryWorkloadIdentity(t)
	var exchanges atomic.Int32
	httpClient := ordinaryWorkloadIssuer(t, func() string {
		return fmt.Sprintf("synthetic-trusted-bearer-%d", exchanges.Add(1))
	})
	request, err := http.NewRequestWithContext(t.Context(), http.MethodGet,
		"https://api.openai.com/v1/models", nil)
	if err != nil {
		t.Fatalf("construct request with ambiguous credentials: %v", err)
	}
	request.Header["Authorization"] = []string{"Bearer synthetic-canonical-attacker", "Bearer synthetic-extra-attacker"}
	request.Header["authorization"] = []string{"Bearer synthetic-lowercase-attacker"}
	request.Header["aUtHoRiZaTiOn"] = []string{"Bearer synthetic-mixed-case-attacker"}
	request.Header.Set("X-Synthetic-Metadata", "preserved")

	var attempts int
	response, err := auth.WorkloadIdentityMiddleware(identity, httpClient, request,
		func(authenticated *http.Request) (*http.Response, error) {
			attempts++
			var authorizationHeaders int
			for name, values := range authenticated.Header {
				if !strings.EqualFold(strings.ReplaceAll(name, "_", "-"), "authorization") {
					continue
				}
				authorizationHeaders++
				want := fmt.Sprintf("Bearer synthetic-trusted-bearer-%d", attempts)
				if name != "Authorization" || len(values) != 1 || values[0] != want {
					t.Errorf("attempt %d authorization header %q=%q, want exactly %q",
						attempts, name, values, want)
				}
			}
			if authorizationHeaders != 1 {
				t.Errorf("attempt %d authorization headers=%d, want one", attempts, authorizationHeaders)
			}
			if got := authenticated.Header.Get("X-Synthetic-Metadata"); got != "preserved" {
				t.Errorf("attempt %d dropped unrelated header: %q", attempts, got)
			}
			status := http.StatusOK
			if attempts == 1 {
				status = http.StatusUnauthorized
			}
			response := ordinaryWorkloadResponse(status, `{}`)
			response.Request = authenticated
			return response, nil
		})
	if err != nil || response == nil || response.StatusCode != http.StatusOK {
		t.Fatalf("ambiguous-credential replay response=%v error=%v", response, err)
	}
	if closeErr := response.Body.Close(); closeErr != nil {
		t.Fatalf("close ambiguous-credential response: %v", closeErr)
	}
	if attempts != 2 || exchanges.Load() != 2 {
		t.Errorf("authenticated attempts/exchanges=%d/%d, want 2/2", attempts, exchanges.Load())
	}
	var callerAuthorizationHeaders int
	for name, values := range request.Header {
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
	for name := range response.Request.Header {
		if strings.EqualFold(strings.ReplaceAll(name, "_", "-"), "authorization") {
			t.Errorf("returned response exposed authorization alias %q", name)
		}
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

func TestWorkloadIdentityMiddlewareInvalidatesRejectedReplayBearer(t *testing.T) {
	identity := newOrdinaryWorkloadIdentity(t)
	var exchanges atomic.Int32
	httpClient := ordinaryWorkloadIssuer(t, func() string {
		return fmt.Sprintf("synthetic-bearer-%d", exchanges.Add(1))
	})
	request, err := http.NewRequestWithContext(t.Context(), http.MethodGet,
		"https://api.openai.com/v1/models", nil)
	if err != nil {
		t.Fatalf("construct repeatedly unauthorized request: %v", err)
	}
	var rejected int
	response, err := auth.WorkloadIdentityMiddleware(identity, httpClient, request,
		func(*http.Request) (*http.Response, error) {
			rejected++
			return ordinaryWorkloadResponse(http.StatusUnauthorized, `{}`), nil
		})
	if err != nil || response == nil || response.StatusCode != http.StatusUnauthorized || rejected != 2 {
		t.Fatalf("repeatedly unauthorized response=%v attempts=%d error=%v", response, rejected, err)
	}
	if closeErr := response.Body.Close(); closeErr != nil {
		t.Fatalf("close repeatedly unauthorized response: %v", closeErr)
	}

	response, err = auth.WorkloadIdentityMiddleware(identity, httpClient, request,
		func(sent *http.Request) (*http.Response, error) {
			if got := sent.Header.Get("Authorization"); got != "Bearer synthetic-bearer-3" {
				t.Errorf("new request reused a previously rejected bearer: %q", got)
			}
			return ordinaryWorkloadResponse(http.StatusOK, `{}`), nil
		})
	if err != nil || response == nil || response.StatusCode != http.StatusOK {
		t.Fatalf("request after rejected replay response=%v error=%v", response, err)
	}
	if closeErr := response.Body.Close(); closeErr != nil {
		t.Fatalf("close refreshed response: %v", closeErr)
	}
	if got := exchanges.Load(); got != 3 {
		t.Errorf("rejected-replay issuer exchanges=%d, want three", got)
	}
}

func TestWorkloadIdentityMiddlewareNeverReturnsSignedRequestMetadata(t *testing.T) {
	for _, test := range []struct {
		name          string
		status        int
		dispatchError bool
		unauthorized  bool
		nilRequest    bool
		nilHeader     bool
		wantAttempts  int
	}{
		{name: "successful response", status: http.StatusOK, wantAttempts: 1},
		{name: "API error response", status: http.StatusForbidden, wantAttempts: 1},
		{name: "response with dispatch error", status: http.StatusOK, dispatchError: true, wantAttempts: 1},
		{name: "unsigned unauthorized replay", status: http.StatusOK, unauthorized: true, wantAttempts: 2},
		{name: "response without request", status: http.StatusOK, nilRequest: true, wantAttempts: 1},
		{name: "response request without headers", status: http.StatusOK, nilHeader: true, wantAttempts: 1},
	} {
		t.Run(test.name, func(t *testing.T) {
			identity := newOrdinaryWorkloadIdentity(t)
			var exchanges atomic.Int32
			httpClient := ordinaryWorkloadIssuer(t, func() string {
				return fmt.Sprintf("synthetic-private-bearer-%d", exchanges.Add(1))
			})
			request, err := http.NewRequestWithContext(t.Context(), http.MethodGet,
				"https://api.openai.com/v1/models", nil)
			if err != nil {
				t.Fatalf("construct unsigned workload request: %v", err)
			}
			request.Header.Set("X-Synthetic-Metadata", "preserved")
			failure := errors.New("synthetic dispatch error")
			var attempts int
			var wire []*http.Response
			response, dispatchErr := auth.WorkloadIdentityMiddleware(identity, httpClient, request,
				func(authenticated *http.Request) (*http.Response, error) {
					attempts++
					if got, want := authenticated.Header.Get("Authorization"),
						fmt.Sprintf("Bearer synthetic-private-bearer-%d", attempts); got != want {
						t.Errorf("wire bearer on attempt %d = %q, want %q", attempts, got, want)
					}
					status := test.status
					if test.unauthorized && attempts == 1 {
						status = http.StatusUnauthorized
					}
					result := ordinaryWorkloadResponse(status, `{}`)
					result.Request = authenticated
					if test.nilRequest {
						result.Request = nil
					} else if test.nilHeader {
						result.Request = authenticated.Clone(authenticated.Context())
						result.Request.Header = nil
					} else {
						result.Request.Header["aUtHoRiZaTiOn"] = []string{"Bearer synthetic-private-alias"}
					}
					wire = append(wire, result)
					if test.dispatchError {
						return result, failure
					}
					return result, nil
				})
			if response == nil || response.StatusCode != test.status {
				t.Fatalf("unsigned response = %v, error = %v", response, dispatchErr)
			}
			if closeErr := response.Body.Close(); closeErr != nil {
				t.Fatalf("close unsigned workload response: %v", closeErr)
			}
			if got := errors.Is(dispatchErr, failure); got != test.dispatchError {
				t.Errorf("dispatch error preservation = %v, want %v", got, test.dispatchError)
			}
			if got := request.Header.Get("Authorization"); got != "" {
				t.Errorf("caller-owned request retained a private bearer: %q", got)
			}
			if response.Request != nil && response.Request.Header != nil {
				for name, values := range response.Request.Header {
					if strings.EqualFold(strings.ReplaceAll(name, "_", "-"), "authorization") {
						t.Errorf("returned response exposed credential header %q: %q", name, values)
					}
				}
				if got := response.Request.Header.Get("X-Synthetic-Metadata"); got != "preserved" {
					t.Errorf("returned response dropped unrelated metadata: %q", got)
				}
			}
			if !test.nilRequest && !test.nilHeader {
				for _, original := range wire {
					if original.Request.Header.Get("Authorization") == "" {
						t.Error("response redaction modified the live authenticated wire request")
					}
				}
			}
			if attempts != test.wantAttempts {
				t.Errorf("authenticated API attempts = %d, want %d", attempts, test.wantAttempts)
			}
		})
	}
}

func TestWorkloadIdentityNativeIssuerRejectsRedirectsWithoutMutatingItsClient(t *testing.T) {
	for _, status := range []int{
		http.StatusMovedPermanently,
		http.StatusFound,
		http.StatusSeeOther,
		http.StatusTemporaryRedirect,
		http.StatusPermanentRedirect,
	} {
		t.Run(fmt.Sprintf("HTTP %d", status), func(t *testing.T) {
			var issuerRequests, redirectedRequests, callerRedirectChecks atomic.Int32
			target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				redirectedRequests.Add(1)
				w.WriteHeader(http.StatusOK)
			}))
			t.Cleanup(target.Close)
			issuer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
				issuerRequests.Add(1)
				body, err := io.ReadAll(request.Body)
				if err != nil || !strings.Contains(string(body), "synthetic-subject") {
					t.Errorf("original issuer did not receive its expected subject token: body=%q error=%v", body, err)
				}
				w.Header().Set("Location", target.URL+"/steal?credential=synthetic-private-redirect-location")
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
			identity := newOrdinaryWorkloadIdentity(t)

			token, err := identity.GetToken(t.Context(), caller)
			if token != "" || err == nil || !strings.Contains(err.Error(), "does not follow redirects") {
				t.Fatalf("redirected token exchange returned token=%q error=%v", token, err)
			}
			if strings.Contains(err.Error(), "synthetic-private") || strings.Contains(err.Error(), target.URL) {
				t.Errorf("issuer redirect error disclosed its attacker-controlled destination: %q", err.Error())
			}
			if issuerRequests.Load() != 1 || redirectedRequests.Load() != 0 || callerRedirectChecks.Load() != 0 {
				t.Errorf("issuer/redirect/caller-policy requests=%d/%d/%d, want 1/0/0",
					issuerRequests.Load(), redirectedRequests.Load(), callerRedirectChecks.Load())
			}
			if caller.CheckRedirect == nil {
				t.Fatal("isolated issuer exchange removed the caller's redirect policy")
			}
			if err := caller.CheckRedirect(nil, nil); err != nil || callerRedirectChecks.Load() != 1 {
				t.Errorf("isolated issuer exchange modified its caller-owned redirect policy: %v", err)
			}
		})
	}
}

func TestWorkloadIdentityRejectsOpaqueIssuerRedirectResponses(t *testing.T) {
	for _, status := range []int{
		http.StatusMovedPermanently,
		http.StatusFound,
		http.StatusSeeOther,
		http.StatusTemporaryRedirect,
		http.StatusPermanentRedirect,
	} {
		t.Run(fmt.Sprintf("HTTP %d", status), func(t *testing.T) {
			body := &ordinaryObservedBody{}
			var requests atomic.Int32
			doer := ordinaryOpaqueHTTPDoer(func(*http.Request) (*http.Response, error) {
				requests.Add(1)
				return &http.Response{
					StatusCode: status,
					Header:     http.Header{"Location": []string{"https://synthetic.invalid/private"}},
					Body:       body,
				}, nil
			})
			token, err := newOrdinaryWorkloadIdentity(t).GetToken(t.Context(), doer)
			if token != "" || err == nil || !strings.Contains(err.Error(), "does not follow redirects") {
				t.Fatalf("opaque redirect response token=%q error=%v", token, err)
			}
			if requests.Load() != 1 || body.reads.Load() != 0 || body.closes.Load() != 1 {
				t.Errorf("opaque issuer requests/body reads/closes=%d/%d/%d, want 1/0/1",
					requests.Load(), body.reads.Load(), body.closes.Load())
			}
		})
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

type ordinaryOpaqueHTTPDoer func(*http.Request) (*http.Response, error)

func (do ordinaryOpaqueHTTPDoer) Do(request *http.Request) (*http.Response, error) {
	return do(request)
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

func TestWorkloadIdentityMiddlewareClosesUnauthorizedBeforeRefresh(t *testing.T) {
	identity := newOrdinaryWorkloadIdentity(t)
	body := &ordinaryObservedBody{}
	var exchanges atomic.Int32
	client := ordinaryWorkloadIssuer(t, func() string {
		attempt := exchanges.Add(1)
		if attempt == 2 && body.closes.Load() != 1 {
			t.Errorf("401 body closes before refresh = %d, want 1", body.closes.Load())
		}
		return fmt.Sprintf("synthetic-bearer-%d", attempt)
	})
	request, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "https://api.openai.com/v1/models", nil)
	if err != nil {
		t.Fatal(err)
	}
	calls := 0
	response, err := auth.WorkloadIdentityMiddleware(identity, client, request, func(*http.Request) (*http.Response, error) {
		calls++
		if calls == 1 {
			return &http.Response{StatusCode: 401, Body: body}, nil
		}
		return ordinaryWorkloadResponse(200, `{}`), nil
	})
	if response != nil {
		_ = response.Body.Close()
	}
	if err != nil || calls != 2 || exchanges.Load() != 2 || body.closes.Load() != 1 {
		t.Errorf("middleware replay = %v, calls=%d, exchanges=%d, closes=%d; want success, 2, 2, 1", err, calls, exchanges.Load(), body.closes.Load())
	}
}

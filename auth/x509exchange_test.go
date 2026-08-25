package auth

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"testing/iotest"
	"time"

	"github.com/openai/openai-go/v3/shared"
)

const x509ExchangeSyntheticToken = "synthetic-workload.token_123+/="

func TestX509ExchangeUsesExactPinnedMutualTLSWireContract(t *testing.T) {
	type observedRequest struct {
		method        string
		path          string
		host          string
		serverName    string
		peer          string
		contentType   string
		authorization string
		cookies       string
		fields        map[string]string
	}
	received := make(chan observedRequest, 1)
	fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		var fields map[string]string
		if json.NewDecoder(request.Body).Decode(&fields) != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		received <- observedRequest{
			method:        request.Method,
			path:          request.URL.RequestURI(),
			host:          request.Host,
			serverName:    request.TLS.ServerName,
			peer:          request.TLS.PeerCertificates[0].Subject.CommonName,
			contentType:   request.Header.Get("Content-Type"),
			authorization: request.Header.Get("Authorization"),
			cookies:       request.Header.Get("Cookie"),
			fields:        fields,
		}
		_, _ = io.WriteString(w, x509ValidExchangeResponse())
	}))
	started := time.Now()
	token, err := x509Exchange(t.Context(), fixture.capability, "synthetic-identity-provider", "synthetic-service-account")
	finished := time.Now()
	if err != nil {
		t.Fatalf("exchange over verified mutual TLS: %v", err)
	}
	if token.value != x509ExchangeSyntheticToken {
		t.Errorf("issued access token = %q, want synthetic workload token", token.value)
	}
	if token.expiresAt.Before(started.Add(60*time.Second)) || token.expiresAt.After(finished.Add(60*time.Second)) {
		t.Errorf("token expiry = %v, want exchange-start-based lifetime", token.expiresAt)
	}
	select {
	case request := <-received:
		if request.method != http.MethodPost || request.path != "/oauth/token" {
			t.Errorf("issuer request = %s %s", request.method, request.path)
		}
		if request.host != x509AuthenticationHost || request.serverName != x509AuthenticationHost {
			t.Errorf("issuer Host/SNI = %q/%q", request.host, request.serverName)
		}
		if request.peer != "synthetic-exchange-workload" || request.contentType != "application/json" {
			t.Errorf("issuer peer/content-type = %q/%q", request.peer, request.contentType)
		}
		if request.authorization != "" || request.cookies != "" {
			t.Errorf("issuer unexpectedly received authorization or cookies")
		}
		want := map[string]string{
			"grant_type":           TokenExchangeGrantType,
			"subject_token_type":   x509SubjectTokenType,
			"identity_provider_id": "synthetic-identity-provider",
			"service_account_id":   "synthetic-service-account",
		}
		if len(request.fields) != len(want) {
			t.Errorf("issuer received %d JSON fields, want exactly four", len(request.fields))
		}
		for name, value := range want {
			if request.fields[name] != value {
				t.Errorf("issuer field %q = %q, want %q", name, request.fields[name], value)
			}
		}
	default:
		t.Fatal("mutually authenticated issuer did not receive a request")
	}
}

func TestX509ExchangeRejectsMalformedOrAmbiguousResponses(t *testing.T) {
	valid := x509ValidExchangeResponse()
	replace := func(old, next string) string { return strings.Replace(valid, old, next, 1) }
	tests := []struct {
		name string
		body string
	}{
		{name: "empty", body: ""},
		{name: "null", body: "null"},
		{name: "array", body: "[]"},
		{name: "string", body: `"token"`},
		{name: "truncated object", body: `{"access_token":"private-malformed-token"`},
		{name: "duplicate token", body: replace(`"access_token":`, `"access_token":"attacker-token","access_token":`)},
		{name: "duplicate extension", body: strings.TrimSuffix(valid, "}") + `,"extension":1,"extension":2}`},
		{name: "trailing object", body: valid + `{"leaked":"private-trailing-token"}`},
		{name: "trailing scalar", body: valid + " true"},
		{name: "trailing malformed", body: valid + " private-trailing-token"},
		{name: "missing token", body: replace(`"access_token":"`+x509ExchangeSyntheticToken+`",`, "")},
		{name: "null token", body: replace(`"access_token":"`+x509ExchangeSyntheticToken+`"`, `"access_token":null`)},
		{name: "empty token", body: replace(x509ExchangeSyntheticToken, "")},
		{name: "token whitespace", body: replace(x509ExchangeSyntheticToken, "synthetic token")},
		{name: "token comma", body: replace(x509ExchangeSyntheticToken, "synthetic,token")},
		{name: "token control", body: replace(x509ExchangeSyntheticToken, `synthetic\ntoken`)},
		{name: "token leading padding", body: replace(x509ExchangeSyntheticToken, "=synthetic")},
		{name: "token interior padding", body: replace(x509ExchangeSyntheticToken, "synthetic=token")},
		{name: "missing token type", body: replace(`"token_type":"Bearer",`, "")},
		{name: "null token type", body: replace(`"token_type":"Bearer"`, `"token_type":null`)},
		{name: "lowercase token type", body: replace(`"token_type":"Bearer"`, `"token_type":"bearer"`)},
		{name: "missing issued type", body: replace(`"issued_token_type":"`+x509IssuedAccessTokenType+`",`, "")},
		{name: "incorrect issued type", body: replace(x509IssuedAccessTokenType, "urn:attacker:token")},
		{name: "missing expiry", body: replace(`,"expires_in":60`, "")},
		{name: "null expiry", body: replace(`"expires_in":60`, `"expires_in":null`)},
		{name: "zero expiry", body: replace(`"expires_in":60`, `"expires_in":0`)},
		{name: "negative expiry", body: replace(`"expires_in":60`, `"expires_in":-1`)},
		{name: "excessive expiry", body: replace(`"expires_in":60`, `"expires_in":3601`)},
		{name: "decimal expiry", body: replace(`"expires_in":60`, `"expires_in":1.5`)},
		{name: "exponent expiry", body: replace(`"expires_in":60`, `"expires_in":1e2`)},
		{name: "string expiry", body: replace(`"expires_in":60`, `"expires_in":"60"`)},
		{name: "overflow expiry", body: replace(`"expires_in":60`, `"expires_in":9223372036854775808`)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				_, _ = io.WriteString(w, test.body)
			}))
			token, err := x509Exchange(t.Context(), fixture.capability, "idp", "service-account")
			if err == nil || token != (x509ExchangedToken{}) {
				t.Fatalf("malformed exchange returned token=%v error=%v", token, err)
			}
			if strings.Contains(err.Error(), "private-") || strings.Contains(err.Error(), x509ExchangeSyntheticToken) {
				t.Errorf("malformed response error exposed sensitive issuer content: %q", err.Error())
			}
		})
	}
}

func TestX509ExchangeAcceptsValidLifetimesAndBearerTokens(t *testing.T) {
	for _, test := range []struct {
		name     string
		token    string
		lifetime int
	}{
		{name: "minimum lifetime", token: "a", lifetime: 1},
		{name: "maximum lifetime", token: "A-z.0_~+/===", lifetime: 3600},
	} {
		t.Run(test.name, func(t *testing.T) {
			body := strings.Replace(x509ValidExchangeResponse(), x509ExchangeSyntheticToken, test.token, 1)
			body = strings.Replace(body, `"expires_in":60`, `"expires_in":`+strconv.Itoa(test.lifetime), 1)
			fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				_, _ = io.WriteString(w, body)
			}))
			token, err := x509Exchange(t.Context(), fixture.capability, "idp", "service-account")
			if err != nil || token.value != test.token {
				t.Fatalf("valid exchange token=%v error=%v", token, err)
			}
		})
	}
}

func TestX509ExchangeAnchorsExpiryBeforeResponse(t *testing.T) {
	body := strings.Replace(x509ValidExchangeResponse(), `"expires_in":60`, `"expires_in":1`, 1)
	started := time.Now().Add(-2 * time.Second)
	token, err := x509DecodeExchangedToken(t.Context(), []byte(body), started)
	if err == nil || token != (x509ExchangedToken{}) || !strings.Contains(err.Error(), "expired") {
		t.Fatalf("already elapsed token lifetime returned token=%v error=%v", token, err)
	}
}

func TestX509ExchangeBoundsSuccessResponseSizes(t *testing.T) {
	for _, test := range []struct {
		name     string
		declared bool
	}{
		{name: "declared oversized body", declared: true},
		{name: "chunked oversized body", declared: false},
	} {
		t.Run(test.name, func(t *testing.T) {
			fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				if test.declared {
					w.Header().Set("Content-Length", strconv.Itoa(x509SuccessResponseMaximum+1))
				} else {
					w.(http.Flusher).Flush()
				}
				_, _ = io.CopyN(w, strings.NewReader(strings.Repeat("x", x509SuccessResponseMaximum+1)),
					int64(x509SuccessResponseMaximum+1))
			}))
			token, err := x509Exchange(t.Context(), fixture.capability, "idp", "service-account")
			if err == nil || token != (x509ExchangedToken{}) || !strings.Contains(err.Error(), "size limit") {
				t.Fatalf("oversized issuer response returned token=%v error=%v", token, err)
			}
		})
	}
}

func TestX509ExchangeRedactsAndClassifiesOAuthErrors(t *testing.T) {
	for _, test := range []struct {
		name      string
		status    int
		body      string
		code      shared.OAuthErrorCode
		wantTyped bool
		oversized bool
		retryable bool
	}{
		{name: "bad request nested issuer code", status: 400,
			body: `{"error":{"code":"invalid_grant","message":"private-nested-issuer-secret"}}`,
			code: shared.OAuthErrorCodeInvalidGrant, wantTyped: true},
		{name: "bad request nested invalid target", status: 400,
			body: `{"error":{"code":"invalid_target","message":"private-nested-issuer-secret"}}`,
			code: shared.OAuthErrorCode("invalid_target"), wantTyped: true},
		{name: "unauthorized nested subject code", status: 401,
			body: `{"error":{"code":"invalid_subject_token","description":"private-nested-issuer-secret"}}`,
			code: shared.OAuthErrorCodeInvalidSubjectToken, wantTyped: true},
		{name: "bad request known code", status: 400, body: `{"error":"invalid_grant","error_description":"private-issuer-secret"}`,
			code: shared.OAuthErrorCodeInvalidGrant, wantTyped: true},
		{name: "bad request flat invalid target", status: 400,
			body: `{"error":"invalid_target","error_description":"private-issuer-secret"}`,
			code: shared.OAuthErrorCode("invalid_target"), wantTyped: true},
		{name: "unauthorized known code", status: 401, body: `{"error":"invalid_subject_token","error_description":"private-issuer-secret"}`,
			code: shared.OAuthErrorCodeInvalidSubjectToken, wantTyped: true},
		{name: "forbidden unknown code", status: 403, body: `{"error":"private-attacker-error","error_description":"private-issuer-secret"}`,
			wantTyped: true},
		{name: "forbidden duplicate code", status: 403, body: `{"error":"invalid_grant","error":"private-attacker-error"}`,
			wantTyped: true},
		{name: "nested unknown code", status: 400, body: `{"error":{"code":"private-attacker-error"}}`,
			wantTyped: true},
		{name: "nested duplicate code", status: 400,
			body: `{"error":{"code":"invalid_grant","code":"private-attacker-error"}}`, wantTyped: true},
		{name: "nested duplicate extension", status: 400,
			body: `{"error":{"code":"invalid_grant","detail":1,"detail":2}}`, wantTyped: true},
		{name: "duplicate nested envelope", status: 400,
			body: `{"error":{"code":"invalid_grant"},"error":{"code":"invalid_subject_token"}}`, wantTyped: true},
		{name: "nested missing code", status: 400, body: `{"error":{"detail":"private-issuer-secret"}}`, wantTyped: true},
		{name: "nested non-string code", status: 400, body: `{"error":{"code":123}}`, wantTyped: true},
		{name: "invalid error envelope", status: 400, body: `{"error":["invalid_grant"]}`, wantTyped: true},
		{name: "bad request malformed", status: 400, body: `private-issuer-secret`, wantTyped: true},
		{name: "oversized typed error", status: 401, body: strings.Repeat("private-issuer-secret", x509ErrorResponseMaximum),
			wantTyped: true, oversized: true},
		{name: "request timeout", status: 408, body: `{"error":"private-issuer-secret"}`, retryable: true},
		{name: "temporary conflict", status: 409, body: `{"error":"private-issuer-secret"}`, retryable: true},
		{name: "rate limited", status: 429, body: `{"error":"private-issuer-secret"}`, retryable: true},
		{name: "server failure", status: 500, body: `{"error":"private-issuer-secret"}`, retryable: true},
		{name: "service unavailable", status: 503, body: `{"error":"private-issuer-secret"}`, retryable: true},
		{name: "not found permanent", status: 404, body: `{"error":"private-issuer-secret"}`},
	} {
		t.Run(test.name, func(t *testing.T) {
			fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				if test.oversized {
					w.Header().Set("Content-Length", strconv.Itoa(len(test.body)))
				}
				w.WriteHeader(test.status)
				_, _ = io.WriteString(w, test.body)
			}))
			token, err := x509Exchange(t.Context(), fixture.capability, "idp", "service-account")
			if err == nil || token != (x509ExchangedToken{}) {
				t.Fatalf("issuer failure returned token=%v error=%v", token, err)
			}
			for cause := err; cause != nil; cause = errors.Unwrap(cause) {
				if strings.Contains(cause.Error(), "private-") || strings.Contains(cause.Error(), "description") {
					t.Errorf("OAuth error chain exposed sensitive issuer content: %q", cause.Error())
				}
			}
			var oauthError *OAuthError
			if got := errors.As(err, &oauthError); got != test.wantTyped {
				t.Fatalf("typed OAuth error = %v, want %v", got, test.wantTyped)
			}
			if oauthError != nil && (oauthError.StatusCode != test.status || oauthError.ErrorCode != test.code ||
				oauthError.ErrorDescription != "") {
				t.Errorf("safe OAuth error = %+v, want status=%d code=%q without description", oauthError, test.status, test.code)
			}
			var statusError *x509ExchangeHTTPError
			if got := errors.As(err, &statusError); got == test.wantTyped {
				t.Fatalf("typed internal HTTP status error = %v, want %v", got, !test.wantTyped)
			}
			if statusError != nil && (statusError.statusCode != test.status || statusError.retryable() != test.retryable) {
				t.Errorf("safe HTTP status taxonomy = %d/retryable:%v, want %d/retryable:%v",
					statusError.statusCode, statusError.retryable(), test.status, test.retryable)
			}
		})
	}
}

func TestX509ExchangeRejectsRedirectBeforeHTTPStatusClassification(t *testing.T) {
	fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTemporaryRedirect)
		_, _ = io.WriteString(w, "private-redirect-response")
	}))
	token, err := x509Exchange(t.Context(), fixture.capability, "idp", "service-account")
	var transportError *x509TransportError
	var statusError *x509ExchangeHTTPError
	if token != (x509ExchangedToken{}) || !errors.As(err, &transportError) ||
		!errors.Is(err, errX509Redirect) || errors.As(err, &statusError) {
		t.Fatalf("issuer redirect = token:%v error:%v, want a safe transport redirect refusal", token, err)
	}
	if strings.Contains(err.Error(), "private-") {
		t.Errorf("redirect error disclosed sensitive issuer content: %q", err.Error())
	}
}

func TestX509ExchangePreservesOnlySafeBoundedRetryAfterDelays(t *testing.T) {
	for _, test := range []struct {
		name    string
		status  int
		headers http.Header
		want    time.Duration
		present bool
	}{
		{name: "preferred milliseconds", status: 429,
			headers: http.Header{"Retry-After-Ms": {"125"}, "Retry-After": {"3"}},
			want:    125 * time.Millisecond, present: true},
		{name: "fractional milliseconds", status: 503,
			headers: http.Header{"Retry-After-Ms": {"1.5"}}, want: 1500 * time.Microsecond, present: true},
		{name: "fractional seconds", status: 429,
			headers: http.Header{"Retry-After": {"0.25"}}, want: 250 * time.Millisecond, present: true},
		{name: "request timeout delay", status: 408,
			headers: http.Header{"Retry-After-Ms": {"20"}}, want: 20 * time.Millisecond, present: true},
		{name: "conflict delay", status: 409,
			headers: http.Header{"Retry-After": {"1"}}, want: time.Second, present: true},
		{name: "finite huge seconds", status: 503,
			headers: http.Header{"Retry-After": {"1e100"}}, want: x509MaximumRetryAfter, present: true},
		{name: "finite scaling overflow", status: 503,
			headers: http.Header{"Retry-After": {"2" + strings.Repeat("0", 299)}},
			want:    x509MaximumRetryAfter, present: true},
		{name: "finite huge milliseconds", status: 429,
			headers: http.Header{"Retry-After-Ms": {"1e100"}}, want: x509MaximumRetryAfter, present: true},
		{name: "future HTTP date bounded", status: 503,
			headers: http.Header{"Retry-After": {time.Now().Add(time.Hour).UTC().Format(http.TimeFormat)}},
			want:    x509MaximumRetryAfter, present: true},
		{name: "past HTTP date immediate", status: 503,
			headers: http.Header{"Retry-After": {time.Now().Add(-time.Hour).UTC().Format(http.TimeFormat)}}, present: true},
		{name: "invalid preferred header falls through", status: 429,
			headers: http.Header{"Retry-After-Ms": {"-1"}, "Retry-After": {"0.5"}},
			want:    500 * time.Millisecond, present: true},
		{name: "malformed preferred header falls through", status: 429,
			headers: http.Header{"Retry-After-Ms": {"private-header-secret"}, "Retry-After": {"0.5"}},
			want:    500 * time.Millisecond, present: true},
		{name: "zero seconds", status: 429, headers: http.Header{"Retry-After": {"0"}}, present: true},
		{name: "zero milliseconds", status: 429, headers: http.Header{"Retry-After-Ms": {"0"}}, present: true},
		{name: "negative milliseconds", status: 429, headers: http.Header{"Retry-After-Ms": {"-100"}}},
		{name: "negative seconds", status: 429, headers: http.Header{"Retry-After": {"-0.5"}}},
		{name: "not a number", status: 429, headers: http.Header{"Retry-After": {"NaN"}}},
		{name: "positive infinity", status: 429, headers: http.Header{"Retry-After": {"+Inf"}}},
		{name: "negative infinity", status: 429, headers: http.Header{"Retry-After": {"-Inf"}}},
		{name: "float overflow", status: 429, headers: http.Header{"Retry-After": {"1e999"}}},
		{name: "private invalid value", status: 503, headers: http.Header{"Retry-After": {"private-header-secret"}}},
		{name: "permanent status ignores header", status: 404,
			headers: http.Header{"Retry-After": {"3"}}},
		{name: "absent headers", status: 503},
	} {
		t.Run(test.name, func(t *testing.T) {
			fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				for name, values := range test.headers {
					for _, value := range values {
						w.Header().Add(name, value)
					}
				}
				w.WriteHeader(test.status)
				_, _ = io.WriteString(w, "private-retry-response-body")
			}))
			token, err := x509Exchange(t.Context(), fixture.capability, "idp", "service-account")
			var status *x509ExchangeHTTPError
			if token != (x509ExchangedToken{}) || !errors.As(err, &status) {
				t.Fatalf("issuer failure = token:%v error:%v, want a private HTTP status error", token, err)
			}
			if status.retryAfter != test.want || status.hasRetryAfter != test.present {
				t.Errorf("safe issuer-directed delay = (%s, %t), want (%s, %t)",
					status.retryAfter, status.hasRetryAfter, test.want, test.present)
			}
			for cause := err; cause != nil; cause = errors.Unwrap(cause) {
				if strings.Contains(cause.Error(), "private-") || strings.Contains(cause.Error(), "Retry-After") {
					t.Errorf("issuer retry error retained raw private headers or response: %q", cause.Error())
				}
			}
		})
	}
}

func TestX509ParseRetryAfterUsesProvidedClockForHTTPDates(t *testing.T) {
	now := time.Date(2026, time.August, 24, 12, 0, 0, 0, time.UTC)
	for _, test := range []struct {
		name   string
		format string
	}{
		{name: "IMF-fixdate", format: http.TimeFormat},
		{name: "obsolete RFC 850", format: time.RFC850},
		{name: "ANSI C asctime", format: time.ANSIC},
	} {
		t.Run(test.name, func(t *testing.T) {
			headers := http.Header{"Retry-After": {now.Add(3 * time.Second).Format(test.format)}}
			if delay, present := x509ParseRetryAfter(headers, now); !present || delay != 3*time.Second {
				t.Errorf("clock-relative Retry-After date = (%s, %t), want (3s, true)", delay, present)
			}
		})
	}
}

func TestX509ExchangeReusesMutualTLSConnectionAfterRetryableStatus(t *testing.T) {
	for _, status := range []int{
		http.StatusRequestTimeout,
		http.StatusConflict,
		http.StatusTooManyRequests,
		http.StatusInternalServerError,
		http.StatusServiceUnavailable,
	} {
		t.Run(strconv.Itoa(status), func(t *testing.T) {
			var requests, connections atomic.Int32
			fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				if requests.Add(1) == 1 {
					w.WriteHeader(status)
					_, _ = io.WriteString(w, "private-retry-response-body")
					return
				}
				_, _ = io.WriteString(w, x509ValidExchangeResponse())
			}))
			originalDial := fixture.template.DialContext
			fixture.template.DialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
				connections.Add(1)
				return originalDial(ctx, network, address)
			}
			transport, err := NewX509Transport(fixture.template)
			if err != nil {
				t.Fatalf("attest counted mutual-TLS transport: %v", err)
			}
			t.Cleanup(func() {
				if closeErr := transport.Close(); closeErr != nil {
					t.Errorf("close counted mutual-TLS transport: %v", closeErr)
				}
			})

			first, err := x509Exchange(t.Context(), transport, "idp", "service-account")
			var failure *x509ExchangeHTTPError
			if first != (x509ExchangedToken{}) || !errors.As(err, &failure) || failure.statusCode != status {
				t.Fatalf("retryable issuer response = token:%v error:%v", first, err)
			}
			if strings.Contains(err.Error(), "private-") {
				t.Errorf("retryable issuer response exposed discarded content: %v", err)
			}
			second, err := x509Exchange(t.Context(), transport, "idp", "service-account")
			if err != nil || second.value != x509ExchangeSyntheticToken {
				t.Fatalf("exchange after drained issuer response = token:%v error:%v", second, err)
			}
			if requests.Load() != 2 || connections.Load() != 1 {
				t.Errorf("retryable issuer requests/TLS handshakes = %d/%d, want 2/1",
					requests.Load(), connections.Load())
			}
		})
	}
}

func TestX509ExchangeReusesMutualTLSConnectionAtChunkedErrorBoundary(t *testing.T) {
	for _, test := range []struct {
		name            string
		bodyLength      int
		wantConnections int32
	}{
		{name: "below maximum error body", bodyLength: x509ErrorResponseMaximum - 1, wantConnections: 1},
		{name: "exactly maximum error body", bodyLength: x509ErrorResponseMaximum, wantConnections: 1},
		{name: "one-byte oversized error body", bodyLength: x509ErrorResponseMaximum + 1, wantConnections: 2},
		{name: "two-byte oversized error body", bodyLength: x509ErrorResponseMaximum + 2, wantConnections: 2},
	} {
		t.Run(test.name, func(t *testing.T) {
			var requests, connections atomic.Int32
			bodyWritten := make(chan struct{})
			releaseChunkedEOF := make(chan struct{})
			var releaseOnce sync.Once
			release := func() { releaseOnce.Do(func() { close(releaseChunkedEOF) }) }
			fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
				if request.ProtoMajor != 1 {
					t.Errorf("chunked retryable issuer used HTTP/%d, want HTTP/1", request.ProtoMajor)
				}
				if requests.Add(1) == 1 {
					w.WriteHeader(http.StatusTooManyRequests)
					if flush, ok := w.(http.Flusher); ok {
						flush.Flush()
					}
					_, _ = io.WriteString(w, strings.Repeat("x", test.bodyLength))
					if flush, ok := w.(http.Flusher); ok {
						flush.Flush()
					}
					close(bodyWritten)
					<-releaseChunkedEOF
					return
				}
				_, _ = io.WriteString(w, x509ValidExchangeResponse())
			}))
			t.Cleanup(release)
			originalDial := fixture.template.DialContext
			fixture.template.DialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
				connections.Add(1)
				return originalDial(ctx, network, address)
			}
			transport, err := NewX509Transport(fixture.template)
			if err != nil {
				t.Fatalf("attest chunked-response mutual-TLS transport: %v", err)
			}
			t.Cleanup(func() {
				if closeErr := transport.Close(); closeErr != nil {
					t.Errorf("close chunked-response mutual-TLS transport: %v", closeErr)
				}
			})
			type exchangeResult struct {
				token x509ExchangedToken
				err   error
			}
			finished := make(chan exchangeResult, 1)
			go func() {
				token, exchangeErr := x509Exchange(t.Context(), transport, "idp", "service-account")
				finished <- exchangeResult{token: token, err: exchangeErr}
			}()
			select {
			case <-bodyWritten:
			case <-time.After(5 * time.Second):
				t.Fatal("retryable issuer never flushed its chunked error payload")
			}
			var result exchangeResult
			completedBeforeEOF := false
			select {
			case result = <-finished:
				completedBeforeEOF = true
			case <-time.After(25 * time.Millisecond):
			}
			release()
			if !completedBeforeEOF {
				select {
				case result = <-finished:
				case <-time.After(5 * time.Second):
					t.Fatal("retryable exchange did not observe the released chunked EOF")
				}
			}
			if test.bodyLength == x509ErrorResponseMaximum && completedBeforeEOF {
				t.Error("exact-boundary retryable response returned before observing the chunked EOF")
			}
			var failure *x509ExchangeHTTPError
			if result.token != (x509ExchangedToken{}) || !errors.As(result.err, &failure) || !failure.retryable() {
				t.Fatalf("chunked retryable issuer response = token:%v error:%v", result.token, result.err)
			}
			second, err := x509Exchange(t.Context(), transport, "idp", "service-account")
			if err != nil || second.value != x509ExchangeSyntheticToken {
				t.Fatalf("exchange after chunked retryable response = token:%v error:%v", second, err)
			}
			if requests.Load() != 2 || connections.Load() != test.wantConnections {
				t.Errorf("chunked issuer requests/TLS handshakes = %d/%d, want 2/%d",
					requests.Load(), connections.Load(), test.wantConnections)
			}
		})
	}
}

func TestX509ExchangeBoundsRetryableErrorBodyDraining(t *testing.T) {
	for _, test := range []struct {
		name          string
		declared      bool
		wantRemaining int
	}{
		{name: "declared oversized error", declared: true, wantRemaining: x509ErrorResponseMaximum + 2},
		{name: "chunked oversized error", wantRemaining: 1},
	} {
		t.Run(test.name, func(t *testing.T) {
			body := strings.NewReader(strings.Repeat("x", x509ErrorResponseMaximum+2))
			response := &http.Response{
				StatusCode:    http.StatusTooManyRequests,
				ContentLength: -1,
				Header:        make(http.Header),
				Body:          io.NopCloser(body),
			}
			if test.declared {
				response.ContentLength = int64(body.Len())
			}
			err := x509ExchangeStatusError(t.Context(), response)
			var status *x509ExchangeHTTPError
			if !errors.As(err, &status) || !status.retryable() {
				t.Fatalf("oversized retryable error classification = %v", err)
			}
			if got := body.Len(); got != test.wantRemaining {
				t.Errorf("oversized retryable body retained %d unread bytes, want %d", got, test.wantRemaining)
			}
		})
	}
}

func TestX509ExchangeClosesOversizedRetryableMutualTLSResponses(t *testing.T) {
	for _, declared := range []bool{false, true} {
		name := "chunked oversized error"
		if declared {
			name = "declared oversized error"
		}
		t.Run(name, func(t *testing.T) {
			var requests, connections atomic.Int32
			fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				if requests.Add(1) == 1 {
					if declared {
						w.Header().Set("Content-Length", strconv.Itoa(x509ErrorResponseMaximum+2))
					}
					w.WriteHeader(http.StatusTooManyRequests)
					if !declared {
						if flush, ok := w.(http.Flusher); ok {
							flush.Flush()
						}
					}
					_, _ = io.WriteString(w, strings.Repeat("x", x509ErrorResponseMaximum+2))
					return
				}
				_, _ = io.WriteString(w, x509ValidExchangeResponse())
			}))
			originalDial := fixture.template.DialContext
			fixture.template.DialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
				connections.Add(1)
				return originalDial(ctx, network, address)
			}
			transport, err := NewX509Transport(fixture.template)
			if err != nil {
				t.Fatalf("attest oversized-response transport: %v", err)
			}
			t.Cleanup(func() {
				if closeErr := transport.Close(); closeErr != nil {
					t.Errorf("close oversized-response transport: %v", closeErr)
				}
			})
			if _, exchangeErr := x509Exchange(t.Context(), transport, "idp", "service-account"); exchangeErr == nil {
				t.Fatal("oversized retryable issuer response unexpectedly succeeded")
			}
			token, err := x509Exchange(t.Context(), transport, "idp", "service-account")
			if err != nil || token.value != x509ExchangeSyntheticToken {
				t.Fatalf("fresh connection after oversized response = token:%v error:%v", token, err)
			}
			if requests.Load() != 2 || connections.Load() != 2 {
				t.Errorf("oversized issuer requests/TLS handshakes = %d/%d, want 2/2",
					requests.Load(), connections.Load())
			}
		})
	}
}

func TestX509ExchangeRedactsRetryableErrorBodyReadFailures(t *testing.T) {
	sensitiveCause := errors.New("Authorization: Bearer private-retry-read-token")
	response := &http.Response{
		StatusCode:    http.StatusServiceUnavailable,
		ContentLength: -1,
		Header:        make(http.Header),
		Body:          io.NopCloser(iotest.ErrReader(sensitiveCause)),
	}
	err := x509ExchangeStatusError(t.Context(), response)
	var readError *x509ExchangeReadError
	if !errors.As(err, &readError) || !readError.retryable() {
		t.Fatalf("retryable error-body read failure = %v, want a safe retryable read error", err)
	}
	if errors.Is(err, sensitiveCause) || errors.Unwrap(err) != nil || strings.Contains(err.Error(), "private-") {
		t.Errorf("retryable error-body read retained sensitive content: %v", err)
	}

	fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Length", "1024")
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = io.WriteString(w, "private-truncated-retry-error")
	}))
	token, err := x509Exchange(t.Context(), fixture.capability, "idp", "service-account")
	readError = nil
	if token != (x509ExchangedToken{}) || !errors.As(err, &readError) || !readError.retryable() {
		t.Fatalf("truncated mutual-TLS retryable response = token:%v error:%v", token, err)
	}
	if strings.Contains(err.Error(), "private-") || errors.Unwrap(err) != nil {
		t.Errorf("truncated retryable response exposed sensitive content: %v", err)
	}
}

func TestX509ExchangeCancelsRetryableErrorBodyDraining(t *testing.T) {
	reached := make(chan struct{})
	fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		if flush, ok := w.(http.Flusher); ok {
			flush.Flush()
		}
		close(reached)
		<-request.Context().Done()
	}))
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	result := make(chan error, 1)
	go func() {
		_, err := x509Exchange(ctx, fixture.capability, "idp", "service-account")
		result <- err
	}()
	select {
	case <-reached:
		cancel()
	case <-time.After(5 * time.Second):
		t.Fatal("retryable issuer never began its response body")
	}
	select {
	case err := <-result:
		if !errors.Is(err, context.Canceled) {
			t.Errorf("canceled retryable response drain error = %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("retryable response body drain ignored request cancellation")
	}
}

func TestX509ExchangePreservesSanitizedReadFailureTaxonomy(t *testing.T) {
	sensitiveCause := errors.New("Authorization: Bearer private-response-read-token")
	response := &http.Response{
		ContentLength: -1,
		Body:          io.NopCloser(iotest.ErrReader(sensitiveCause)),
	}
	body, err := x509ReadExchangeResponse(t.Context(), response, x509SuccessResponseMaximum)
	var readError *x509ExchangeReadError
	if body != nil || !errors.As(err, &readError) || !readError.retryable() {
		t.Fatalf("issuer body-read failure = body:%q error:%v, want safe retryable typed error", body, err)
	}
	if errors.Is(err, sensitiveCause) || errors.Unwrap(err) != nil || strings.Contains(err.Error(), "private-") {
		t.Errorf("issuer body-read failure retained its sensitive original cause: %v", err)
	}

	fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Length", "1024")
		_, _ = io.WriteString(w, "private-truncated-response-token")
	}))
	token, err := x509Exchange(t.Context(), fixture.capability, "idp", "service-account")
	readError = nil
	if token != (x509ExchangedToken{}) || !errors.As(err, &readError) || !readError.retryable() {
		t.Fatalf("real mTLS body-read failure = token:%v error:%v", token, err)
	}
	if strings.Contains(err.Error(), "private-") || errors.Unwrap(err) != nil {
		t.Errorf("real mTLS body-read failure exposed sensitive response data: %v", err)
	}
}

func TestX509ExchangeCancellationPrecedesHTTPStatusClassification(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	for _, status := range []int{http.StatusBadRequest, http.StatusUnauthorized, http.StatusTooManyRequests,
		http.StatusInternalServerError, http.StatusServiceUnavailable} {
		t.Run(strconv.Itoa(status), func(t *testing.T) {
			response := &http.Response{
				StatusCode:    status,
				ContentLength: -1,
				Body:          io.NopCloser(strings.NewReader("private-canceled-response-body")),
			}
			err := x509ExchangeStatusError(ctx, response)
			if !errors.Is(err, context.Canceled) {
				t.Errorf("canceled HTTP %d classification = %v, want context.Canceled", status, err)
			}
			var statusError *x509ExchangeHTTPError
			var oauthError *OAuthError
			if errors.As(err, &statusError) || errors.As(err, &oauthError) {
				t.Errorf("canceled HTTP %d retained an OAuth/HTTP status classification", status)
			}
		})
	}
}

func TestX509ExchangeRejectsInvalidInputsBeforeNetwork(t *testing.T) {
	var dialed atomic.Int32
	fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	originalDial := fixture.template.DialContext
	fixture.template.DialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
		dialed.Add(1)
		return originalDial(ctx, network, address)
	}
	capability, err := NewX509Transport(fixture.template)
	if err != nil {
		t.Fatalf("attest counted native transport: %v", err)
	}
	t.Cleanup(func() { _ = capability.Close() })
	canceled, cancel := context.WithCancel(t.Context())
	cancel()
	tests := []struct {
		name       string
		ctx        context.Context
		transport  *X509Transport
		idp        string
		account    string
		wantCancel bool
	}{
		{name: "nil context", ctx: nil, transport: capability, idp: "idp", account: "account"},
		{name: "canceled context", ctx: canceled, transport: capability, idp: "idp", account: "account", wantCancel: true},
		{name: "nil transport", ctx: t.Context(), transport: nil, idp: "idp", account: "account"},
		{name: "empty identity provider", ctx: t.Context(), transport: capability, account: "account"},
		{name: "blank identity provider", ctx: t.Context(), transport: capability, idp: "  ", account: "account"},
		{name: "empty service account", ctx: t.Context(), transport: capability, idp: "idp"},
		{name: "blank service account", ctx: t.Context(), transport: capability, idp: "idp", account: " \t"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			token, exchangeErr := x509Exchange(test.ctx, test.transport, test.idp, test.account)
			if exchangeErr == nil || token != (x509ExchangedToken{}) {
				t.Fatalf("invalid exchange returned token=%v error=%v", token, exchangeErr)
			}
			if test.wantCancel && !errors.Is(exchangeErr, context.Canceled) {
				t.Errorf("canceled exchange error = %v, want context.Canceled", exchangeErr)
			}
		})
	}
	if got := dialed.Load(); got != 0 {
		t.Errorf("invalid exchange inputs caused %d network dials", got)
	}
}

func TestX509ExchangePreservesInFlightCancellation(t *testing.T) {
	reached := make(chan struct{})
	release := make(chan struct{})
	defer close(release)
	fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		close(reached)
		<-release
	}))
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	result := make(chan error, 1)
	go func() {
		_, err := x509Exchange(ctx, fixture.capability, "idp", "service-account")
		result <- err
	}()
	select {
	case <-reached:
		cancel()
	case <-time.After(5 * time.Second):
		t.Fatal("mutual-TLS issuer did not receive the cancellable exchange")
	}
	select {
	case err := <-result:
		if !errors.Is(err, context.Canceled) {
			t.Errorf("canceled in-flight exchange error = %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("in-flight token exchange did not observe cancellation")
	}
}

func x509ValidExchangeResponse() string {
	return `{"access_token":"` + x509ExchangeSyntheticToken +
		`","token_type":"Bearer","issued_token_type":"` + x509IssuedAccessTokenType + `","expires_in":60}`
}

type x509ExchangeFixture struct {
	template   *http.Transport
	capability *X509Transport
}

func newX509ExchangeFixture(t *testing.T, handler http.Handler) *x509ExchangeFixture {
	t.Helper()
	rootKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate ephemeral exchange root key: %v", err)
	}
	rootTemplate := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "synthetic exchange root"},
		NotBefore:             time.Now().Add(-time.Minute),
		NotAfter:              time.Now().Add(time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	rootDER, err := x509.CreateCertificate(rand.Reader, rootTemplate, rootTemplate, &rootKey.PublicKey, rootKey)
	if err != nil {
		t.Fatalf("issue ephemeral exchange root: %v", err)
	}
	root, err := x509.ParseCertificate(rootDER)
	if err != nil {
		t.Fatalf("parse ephemeral exchange root: %v", err)
	}
	trust := x509.NewCertPool()
	trust.AddCert(root)
	issue := func(name string, client bool) tls.Certificate {
		key, keyErr := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if keyErr != nil {
			t.Fatalf("generate ephemeral exchange certificate key: %v", keyErr)
		}
		publicKey, keyErr := x509.MarshalPKIXPublicKey(&key.PublicKey)
		if keyErr != nil {
			t.Fatalf("encode ephemeral exchange certificate key: %v", keyErr)
		}
		keyID := sha256.Sum256(publicKey)
		template := &x509.Certificate{
			SerialNumber:   big.NewInt(2),
			Subject:        pkix.Name{CommonName: name},
			NotBefore:      time.Now().Add(-time.Minute),
			NotAfter:       time.Now().Add(time.Hour),
			KeyUsage:       x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			DNSNames:       []string{name},
			ExtKeyUsage:    []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
			SubjectKeyId:   keyID[:],
			AuthorityKeyId: root.SubjectKeyId,
		}
		if client {
			template.DNSNames = []string{name + ".workload.example.test"}
			template.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}
		}
		encoded, issueErr := x509.CreateCertificate(rand.Reader, template, root, &key.PublicKey, rootKey)
		if issueErr != nil {
			t.Fatalf("issue ephemeral exchange certificate: %v", issueErr)
		}
		leaf, parseErr := x509.ParseCertificate(encoded)
		if parseErr != nil {
			t.Fatalf("parse ephemeral exchange certificate: %v", parseErr)
		}
		if client && !bytes.Equal(leaf.AuthorityKeyId, root.SubjectKeyId) {
			t.Fatal("ephemeral exchange workload certificate lacks its issuer key identity")
		}
		return tls.Certificate{Certificate: [][]byte{encoded}, PrivateKey: key, Leaf: leaf}
	}
	server := httptest.NewUnstartedServer(handler)
	server.Config.ErrorLog = log.New(io.Discard, "", 0)
	server.TLS = &tls.Config{
		MinVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{issue(x509AuthenticationHost, false)},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    trust,
	}
	server.StartTLS()
	t.Cleanup(server.Close)
	dialer := &net.Dialer{Timeout: 5 * time.Second}
	template := &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion:   tls.VersionTLS12,
			RootCAs:      trust,
			Certificates: []tls.Certificate{issue("synthetic-exchange-workload", true)},
		},
		DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
			if address != x509AuthenticationHost+":443" {
				return nil, fmt.Errorf("unexpected hermetic token-exchange address %q", address)
			}
			return dialer.DialContext(ctx, network, server.Listener.Addr().String())
		},
	}
	t.Cleanup(template.CloseIdleConnections)
	capability, err := NewX509Transport(template)
	if err != nil {
		t.Fatalf("attest ephemeral exchange transport: %v", err)
	}
	t.Cleanup(func() {
		if closeErr := capability.Close(); closeErr != nil {
			t.Errorf("close isolated exchange transport: %v", closeErr)
		}
	})
	return &x509ExchangeFixture{template: template, capability: capability}
}

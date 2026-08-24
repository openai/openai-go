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
	"sync/atomic"
	"testing"
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
		name       string
		status     int
		body       string
		code       shared.OAuthErrorCode
		wantTyped  bool
		oversized  bool
		wantCancel bool
	}{
		{name: "bad request known code", status: 400, body: `{"error":"invalid_grant","error_description":"private-issuer-secret"}`,
			code: shared.OAuthErrorCodeInvalidGrant, wantTyped: true},
		{name: "unauthorized known code", status: 401, body: `{"error":"invalid_subject_token","error_description":"private-issuer-secret"}`,
			code: shared.OAuthErrorCodeInvalidSubjectToken, wantTyped: true},
		{name: "forbidden unknown code", status: 403, body: `{"error":"private-attacker-error","error_description":"private-issuer-secret"}`,
			wantTyped: true},
		{name: "forbidden duplicate code", status: 403, body: `{"error":"invalid_grant","error":"private-attacker-error"}`,
			wantTyped: true},
		{name: "bad request malformed", status: 400, body: `private-issuer-secret`, wantTyped: true},
		{name: "oversized typed error", status: 401, body: strings.Repeat("private-issuer-secret", x509ErrorResponseMaximum),
			wantTyped: true, oversized: true},
		{name: "rate limited", status: 429, body: `{"error":"private-issuer-secret"}`},
		{name: "server failure", status: 500, body: `{"error":"private-issuer-secret"}`},
		{name: "unexpected redirect", status: 307, body: `{"error":"private-issuer-secret"}`},
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

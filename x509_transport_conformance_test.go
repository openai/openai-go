package openai_test

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
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/auth"
	"github.com/openai/openai-go/v3/option"
)

const (
	x509ConformanceIssuerHost = "mtls.auth.openai.com"
	x509ConformanceAPIHost    = "mtls.api.openai.com"
	x509ConformanceToken      = "synthetic-workload-access-token"
	x509ConformanceProxyAuth  = "Basic c3ludGhldGljLXByb3h5LWNyZWRlbnRpYWw="
	x509ConformanceTokenType  = "urn:openai:params:oauth:token-type:x509"
)

// X.509 workload exchanges issue ordinary bearer tokens. The token is not
// cryptographically bound to the certificate used during the issuing handshake.
func TestX509TransportConformanceBearerIsNotCertificateBound(t *testing.T) {
	lab := newX509ConformanceLab(t)
	issuer := lab.server(t, x509ConformanceIssuerHost, true)
	api := lab.server(t, x509ConformanceAPIHost, false)
	routes := x509ConformanceRoutes(issuer, api)

	issuerTransport := lab.transport(t, routes, lab.identity(t, "issuer-workload", true))
	firstAPITransport := lab.transport(t, routes, lab.identity(t, "first-api-workload", true))
	secondAPITransport := lab.transport(t, routes, lab.identity(t, "second-api-workload", true))

	httpClient := &http.Client{
		Transport: x509ConformanceHostTransport{
			x509ConformanceIssuerHost: issuerTransport,
			x509ConformanceAPIHost:    firstAPITransport,
		},
	}
	token, err := x509ConformanceExchange(t.Context(), httpClient)
	if err != nil {
		t.Fatalf("exchange the X.509 workload identity: %v", err)
	}
	client := x509ConformanceClient(t, httpClient, token)

	if _, err := client.Models.List(t.Context()); err != nil {
		t.Fatalf("request with independently selected issuer/API certificates: %v", err)
	}
	if _, err := client.Models.List(t.Context(), option.WithHTTPClient(&http.Client{
		Transport: x509ConformanceHostTransport{
			x509ConformanceIssuerHost: issuerTransport,
			x509ConformanceAPIHost:    secondAPITransport,
		},
	})); err != nil {
		t.Fatalf("request with the cached bearer and a different API certificate: %v", err)
	}

	issuerRequests := issuer.requests()
	apiRequests := api.requests()
	if len(issuerRequests) != 1 {
		t.Fatalf("issuer received %d requests, want exactly one cached exchange", len(issuerRequests))
	}
	if len(apiRequests) != 2 {
		t.Fatalf("API received %d requests, want two", len(apiRequests))
	}

	issuerRequest := issuerRequests[0]
	if issuerRequest.method != http.MethodPost || issuerRequest.path != "/oauth/token" {
		t.Errorf("issuer request = %s %s, want POST /oauth/token", issuerRequest.method, issuerRequest.path)
	}
	if issuerRequest.contentType != "application/json" {
		t.Errorf("issuer Content-Type = %q, want application/json", issuerRequest.contentType)
	}
	wantFields := map[string]string{
		"grant_type":           auth.TokenExchangeGrantType,
		"subject_token_type":   x509ConformanceTokenType,
		"identity_provider_id": "synthetic-identity-provider",
		"service_account_id":   "synthetic-service-account",
	}
	if len(issuerRequest.exchangeFields) != len(wantFields) {
		t.Errorf("issuer JSON fields = %v, want exactly %v", issuerRequest.exchangeFields, wantFields)
	}
	for name, want := range wantFields {
		if got := issuerRequest.exchangeFields[name]; got != want {
			t.Errorf("issuer JSON %q = %q, want %q", name, got, want)
		}
	}
	if issuerRequest.authorization != "" || issuerRequest.proxyAuthorization != "" {
		t.Errorf("issuer received origin/proxy credentials: authorization=%q proxy=%q", issuerRequest.authorization, issuerRequest.proxyAuthorization)
	}
	if issuerRequest.peerCommonName != "issuer-workload" {
		t.Errorf("issuer peer certificate = %q, want issuer-workload", issuerRequest.peerCommonName)
	}
	if issuerRequest.host != x509ConformanceIssuerHost || issuerRequest.serverName != x509ConformanceIssuerHost {
		t.Errorf("issuer Host/SNI = %q/%q, want %q", issuerRequest.host, issuerRequest.serverName, x509ConformanceIssuerHost)
	}

	for i, request := range apiRequests {
		if request.method != http.MethodGet || request.path != "/v1/models" {
			t.Errorf("API request %d = %s %s, want GET /v1/models", i, request.method, request.path)
		}
		if request.authorization != "Bearer "+x509ConformanceToken {
			t.Errorf("API request %d authorization = %q", i, request.authorization)
		}
		if request.proxyAuthorization != "" {
			t.Errorf("API request %d unexpectedly received proxy credentials", i)
		}
		if request.host != x509ConformanceAPIHost || request.serverName != x509ConformanceAPIHost {
			t.Errorf("API request %d Host/SNI = %q/%q, want %q", i, request.host, request.serverName, x509ConformanceAPIHost)
		}
	}
	if apiRequests[0].peerCommonName != "first-api-workload" || apiRequests[1].peerCommonName != "second-api-workload" {
		t.Errorf("API peer certificates = %q, %q; want first-api-workload, second-api-workload", apiRequests[0].peerCommonName, apiRequests[1].peerCommonName)
	}
}

func TestX509TransportConformanceRejectsInvalidTLSBeforeHTTP(t *testing.T) {
	tests := []struct {
		name           string
		issuerHostname string
		certificate    func(*testing.T, *x509ConformanceLab) []tls.Certificate
	}{
		{
			name:           "missing client certificate",
			issuerHostname: x509ConformanceIssuerHost,
			certificate:    func(*testing.T, *x509ConformanceLab) []tls.Certificate { return nil },
		},
		{
			name:           "incomplete client certificate chain",
			issuerHostname: x509ConformanceIssuerHost,
			certificate: func(t *testing.T, lab *x509ConformanceLab) []tls.Certificate {
				return []tls.Certificate{lab.identity(t, "incomplete-workload", false)}
			},
		},
		{
			name:           "server hostname mismatch",
			issuerHostname: "wrong-hostname.example.test",
			certificate: func(t *testing.T, lab *x509ConformanceLab) []tls.Certificate {
				return []tls.Certificate{lab.identity(t, "valid-workload", true)}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lab := newX509ConformanceLab(t)
			issuer := lab.server(t, test.issuerHostname, true)
			api := lab.server(t, x509ConformanceAPIHost, false)
			transport := lab.transport(t, x509ConformanceRoutes(issuer, api), test.certificate(t, lab)...)

			if _, err := x509ConformanceExchange(t.Context(), &http.Client{Transport: transport}); err == nil {
				t.Fatal("token exchange unexpectedly succeeded with invalid mutual TLS")
			}
			if requests := issuer.requests(); len(requests) != 0 {
				t.Errorf("issuer received %d HTTP requests after a rejected TLS handshake", len(requests))
			}
			if requests := api.requests(); len(requests) != 0 {
				t.Errorf("API received %d HTTP requests after a rejected issuer TLS handshake", len(requests))
			}
		})
	}
}

func TestX509TransportConformanceHTTPConnectKeepsCredentialsSeparated(t *testing.T) {
	lab := newX509ConformanceLab(t)
	issuer := lab.server(t, x509ConformanceIssuerHost, true)
	api := lab.server(t, x509ConformanceAPIHost, false)
	routes := x509ConformanceRoutes(issuer, api)
	proxy := newX509ConformanceProxy(t, routes)

	transport := lab.transport(t, routes, lab.identity(t, "proxied-workload", true))
	transport.Proxy = http.ProxyURL(proxy.url)
	transport.ProxyConnectHeader = http.Header{"Proxy-Authorization": {x509ConformanceProxyAuth}}

	httpClient := &http.Client{Transport: transport}
	token, err := x509ConformanceExchange(t.Context(), httpClient)
	if err != nil {
		t.Fatalf("X.509 token exchange through authenticated HTTP CONNECT proxy: %v", err)
	}
	client := x509ConformanceClient(t, httpClient, token)
	if _, err := client.Models.List(t.Context()); err != nil {
		t.Fatalf("SDK request through authenticated HTTP CONNECT proxy: %v", err)
	}

	connectRequests := proxy.requests()
	if len(connectRequests) != 2 {
		t.Fatalf("proxy received %d CONNECT requests, want issuer and API tunnels", len(connectRequests))
	}
	for i, request := range connectRequests {
		wantAuthority := []string{x509ConformanceIssuerHost + ":443", x509ConformanceAPIHost + ":443"}[i]
		if request.authority != wantAuthority {
			t.Errorf("CONNECT request %d authority = %q, want %q", i, request.authority, wantAuthority)
		}
		if request.proxyAuthorization != x509ConformanceProxyAuth {
			t.Errorf("CONNECT request %d proxy authorization = %q", i, request.proxyAuthorization)
		}
		if request.authorization != "" {
			t.Errorf("CONNECT request %d leaked workload bearer %q", i, request.authorization)
		}
	}

	issuerRequests, apiRequests := issuer.requests(), api.requests()
	for name, requests := range map[string][]x509ConformanceRequest{
		"issuer": issuerRequests,
		"API":    apiRequests,
	} {
		if len(requests) != 1 {
			t.Errorf("%s received %d requests, want one", name, len(requests))
			continue
		}
		if requests[0].proxyAuthorization != "" {
			t.Errorf("%s received CONNECT proxy credentials", name)
		}
		if requests[0].peerCommonName != "proxied-workload" {
			t.Errorf("%s peer certificate = %q, want proxied-workload", name, requests[0].peerCommonName)
		}
	}
	if len(issuerRequests) == 1 && len(apiRequests) == 1 &&
		(issuerRequests[0].authorization != "" || apiRequests[0].authorization != "Bearer "+x509ConformanceToken) {
		t.Error("issuer and API did not isolate the exchanged bearer credential")
	}
}

func TestX509TransportConformanceRejectsRedirectsBeforeCredentialDisclosure(t *testing.T) {
	for _, leg := range []string{"token exchange", "API request"} {
		t.Run(leg, func(t *testing.T) {
			lab := newX509ConformanceLab(t)
			issuer := lab.server(t, x509ConformanceIssuerHost, true)
			api := lab.server(t, x509ConformanceAPIHost, false)
			destination := lab.server(t, "redirect-target.x509.test", false)
			routes := x509ConformanceRoutes(issuer, api)
			routes["redirect-target.x509.test:443"] = destination.server.Listener.Addr().String()

			refused := errors.New("X.509 workload identity redirects are disabled")
			httpClient := &http.Client{
				Transport: lab.transport(t, routes, lab.identity(t, "redirect-workload", true)),
				CheckRedirect: func(*http.Request, []*http.Request) error {
					return refused
				},
			}
			location := "https://redirect-target.x509.test/credential-capture"

			var err error
			switch leg {
			case "token exchange":
				issuer.setRedirect(location)
				_, err = x509ConformanceExchange(t.Context(), httpClient)
			case "API request":
				var token string
				token, err = x509ConformanceExchange(t.Context(), httpClient)
				if err != nil {
					t.Fatalf("exchange token before API redirect: %v", err)
				}
				api.setRedirect(location)
				client := x509ConformanceClient(t, httpClient, token)
				_, err = client.Models.List(t.Context())
			}
			if !errors.Is(err, refused) {
				t.Fatalf("redirect error = %v, want the concrete redirect-refusal error", err)
			}
			if requests := destination.requests(); len(requests) != 0 {
				t.Errorf("redirect destination received %d HTTP requests", len(requests))
			}
			if handshakes := destination.handshakeCount(); handshakes != 0 {
				t.Errorf("redirect destination observed %d TLS handshakes/client certificates", handshakes)
			}
		})
	}
}

// net/http applies one TLSClientConfig to both the HTTPS proxy and the tunneled
// origin, so an HTTPS proxy requesting client authentication receives the
// caller's workload certificate before the intended issuer is contacted.
func TestX509TransportConformanceHTTPSProxyExposesWorkloadCertificate(t *testing.T) {
	workloadLab := newX509ConformanceLab(t)
	proxyLab := newX509ConformanceLab(t)
	issuer := workloadLab.server(t, x509ConformanceIssuerHost, true)
	api := workloadLab.server(t, x509ConformanceAPIHost, false)
	routes := x509ConformanceRoutes(issuer, api)

	type proxyObservation struct {
		method             string
		host               string
		serverName         string
		peerCommonName     string
		authorization      string
		proxyAuthorization string
	}
	observations := make(chan proxyObservation, 1)
	proxy := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		observations <- proxyObservation{
			method:             req.Method,
			host:               req.Host,
			serverName:         req.TLS.ServerName,
			peerCommonName:     req.TLS.PeerCertificates[0].Subject.CommonName,
			authorization:      req.Header.Get("Authorization"),
			proxyAuthorization: req.Header.Get("Proxy-Authorization"),
		}
		http.Error(w, "synthetic HTTPS proxy refuses CONNECT", http.StatusBadGateway)
	}))
	proxyLeaf, proxyKey := proxyLab.issue(
		t,
		"proxy.x509.test",
		proxyLab.intermediate,
		proxyLab.issuerKey,
		false,
		[]string{"proxy.x509.test"},
	)
	proxy.TLS = &tls.Config{
		MinVersion: tls.VersionTLS12,
		ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs:  workloadLab.trust,
		Certificates: []tls.Certificate{{
			Certificate: [][]byte{proxyLeaf.Raw, proxyLab.intermediate.Raw},
			PrivateKey:  proxyKey,
			Leaf:        proxyLeaf,
		}},
	}
	proxy.Config.ErrorLog = log.New(io.Discard, "", 0)
	proxy.StartTLS()
	t.Cleanup(proxy.Close)
	routes["proxy.x509.test:443"] = proxy.Listener.Addr().String()

	transport := workloadLab.transport(t, routes, workloadLab.identity(t, "disclosed-workload", true))
	trust := x509.NewCertPool()
	trust.AddCert(workloadLab.root)
	trust.AddCert(proxyLab.root)
	transport.TLSClientConfig.RootCAs = trust
	transport.Proxy = http.ProxyURL(&url.URL{Scheme: "https", Host: "proxy.x509.test"})
	transport.ProxyConnectHeader = http.Header{"Proxy-Authorization": {x509ConformanceProxyAuth}}

	if _, err := x509ConformanceExchange(t.Context(), &http.Client{Transport: transport}); err == nil {
		t.Fatal("token exchange unexpectedly succeeded through the rejecting HTTPS proxy")
	}

	select {
	case observed := <-observations:
		if observed.method != http.MethodConnect || observed.host != x509ConformanceIssuerHost+":443" {
			t.Errorf("HTTPS proxy request = %s %s, want CONNECT %s:443", observed.method, observed.host, x509ConformanceIssuerHost)
		}
		if observed.serverName != "proxy.x509.test" || observed.peerCommonName != "disclosed-workload" {
			t.Errorf("HTTPS proxy SNI/client certificate = %q/%q", observed.serverName, observed.peerCommonName)
		}
		if observed.authorization != "" || observed.proxyAuthorization != x509ConformanceProxyAuth {
			t.Errorf("HTTPS proxy origin/proxy authorization = %q/%q", observed.authorization, observed.proxyAuthorization)
		}
	default:
		t.Fatal("HTTPS proxy did not observe the workload certificate on a real TLS handshake")
	}
	if requests := issuer.requests(); len(requests) != 0 {
		t.Errorf("issuer received %d requests after the HTTPS proxy disclosed the client certificate", len(requests))
	}
	if requests := api.requests(); len(requests) != 0 {
		t.Errorf("API received %d requests after the HTTPS proxy disclosed the client certificate", len(requests))
	}
}

func x509ConformanceExchange(ctx context.Context, httpClient *http.Client) (string, error) {
	body, err := json.Marshal(struct {
		GrantType          string `json:"grant_type"`
		SubjectTokenType   string `json:"subject_token_type"`
		IdentityProviderID string `json:"identity_provider_id"`
		ServiceAccountID   string `json:"service_account_id"`
	}{
		GrantType:          auth.TokenExchangeGrantType,
		SubjectTokenType:   x509ConformanceTokenType,
		IdentityProviderID: "synthetic-identity-provider",
		ServiceAccountID:   "synthetic-service-account",
	})
	if err != nil {
		return "", fmt.Errorf("encode X.509 exchange request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://"+x509ConformanceIssuerHost+"/oauth/token", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create X.509 exchange request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	response, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("send X.509 exchange request: %w", err)
	}
	defer func() { _ = response.Body.Close() }()
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("X.509 token exchange returned status %d", response.StatusCode)
	}
	var token struct {
		AccessToken     string `json:"access_token"`
		TokenType       string `json:"token_type"`
		IssuedTokenType string `json:"issued_token_type"`
		ExpiresIn       int    `json:"expires_in"`
	}
	if err := json.NewDecoder(response.Body).Decode(&token); err != nil {
		return "", fmt.Errorf("decode X.509 exchange response: %w", err)
	}
	if token.AccessToken == "" || token.TokenType != "Bearer" ||
		token.IssuedTokenType != "urn:ietf:params:oauth:token-type:access_token" || token.ExpiresIn <= 0 {
		return "", fmt.Errorf("invalid X.509 token-exchange response")
	}
	return token.AccessToken, nil
}

func x509ConformanceClient(t *testing.T, httpClient *http.Client, token string) openai.Client {
	t.Helper()
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("OPENAI_ADMIN_KEY", "")
	t.Setenv("OPENAI_CUSTOM_HEADERS", "")
	return openai.NewClient(
		option.WithBaseURL("https://"+x509ConformanceAPIHost+"/v1/"),
		option.WithAPIKey(token),
		option.WithHTTPClient(httpClient),
		option.WithMaxRetries(0),
	)
}

type x509ConformanceLab struct {
	root         *x509.Certificate
	rootKey      *ecdsa.PrivateKey
	intermediate *x509.Certificate
	issuerKey    *ecdsa.PrivateKey
	trust        *x509.CertPool
}

func newX509ConformanceLab(t *testing.T) *x509ConformanceLab {
	t.Helper()
	lab := &x509ConformanceLab{}
	lab.root, lab.rootKey = lab.issue(t, "synthetic root", nil, nil, true, nil)
	lab.intermediate, lab.issuerKey = lab.issue(t, "synthetic intermediate", lab.root, lab.rootKey, true, nil)
	lab.trust = x509.NewCertPool()
	lab.trust.AddCert(lab.root)
	return lab
}

func (lab *x509ConformanceLab) issue(
	t *testing.T,
	name string,
	issuer *x509.Certificate,
	issuerKey *ecdsa.PrivateKey,
	certificateAuthority bool,
	hostnames []string,
) (*x509.Certificate, *ecdsa.PrivateKey) {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate ephemeral P-256 key: %v", err)
	}
	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		t.Fatalf("generate certificate serial: %v", err)
	}
	publicKey, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		t.Fatalf("encode ephemeral %q public key: %v", name, err)
	}
	keyID := sha256.Sum256(publicKey)

	template := &x509.Certificate{
		SerialNumber:          serial,
		Subject:               pkix.Name{CommonName: name},
		NotBefore:             time.Now().Add(-time.Minute),
		NotAfter:              time.Now().Add(time.Hour),
		BasicConstraintsValid: true,
		IsCA:                  certificateAuthority,
		DNSNames:              hostnames,
		KeyUsage:              x509.KeyUsageDigitalSignature,
		SubjectKeyId:          keyID[:],
	}
	if certificateAuthority {
		template.KeyUsage |= x509.KeyUsageCertSign
	} else if len(hostnames) != 0 {
		template.KeyUsage |= x509.KeyUsageKeyEncipherment
		template.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
	} else {
		template.KeyUsage |= x509.KeyUsageKeyEncipherment
		template.DNSNames = []string{name + ".workload.x509.test"}
		template.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}
	}
	if issuer == nil {
		issuer, issuerKey = template, key
	}
	template.AuthorityKeyId = issuer.SubjectKeyId

	encoded, err := x509.CreateCertificate(rand.Reader, template, issuer, &key.PublicKey, issuerKey)
	if err != nil {
		t.Fatalf("issue ephemeral %q certificate: %v", name, err)
	}
	certificate, err := x509.ParseCertificate(encoded)
	if err != nil {
		t.Fatalf("parse ephemeral %q certificate: %v", name, err)
	}
	return certificate, key
}

func (lab *x509ConformanceLab) identity(t *testing.T, name string, completeChain bool) tls.Certificate {
	t.Helper()
	certificate, key := lab.issue(t, name, lab.intermediate, lab.issuerKey, false, nil)
	requiredKeyUsage := x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment
	if len(certificate.DNSNames) == 0 || len(certificate.SubjectKeyId) == 0 ||
		!bytes.Equal(certificate.AuthorityKeyId, lab.intermediate.SubjectKeyId) ||
		certificate.KeyUsage&requiredKeyUsage != requiredKeyUsage ||
		len(certificate.ExtKeyUsage) != 1 || certificate.ExtKeyUsage[0] != x509.ExtKeyUsageClientAuth {
		t.Fatalf("generated workload certificate %q does not satisfy X.509 enrollment requirements", name)
	}
	chain := [][]byte{certificate.Raw}
	if completeChain {
		chain = append(chain, lab.intermediate.Raw)
	}
	return tls.Certificate{Certificate: chain, PrivateKey: key, Leaf: certificate}
}

func (lab *x509ConformanceLab) transport(
	t *testing.T,
	routes map[string]string,
	identities ...tls.Certificate,
) *http.Transport {
	t.Helper()
	dialer := &net.Dialer{Timeout: 5 * time.Second}
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion:   tls.VersionTLS12,
			RootCAs:      lab.trust,
			Certificates: identities,
		},
		DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
			if target, ok := routes[address]; ok {
				address = target
			}
			return dialer.DialContext(ctx, network, address)
		},
	}
	t.Cleanup(transport.CloseIdleConnections)
	return transport
}

type x509ConformanceRequest struct {
	method             string
	path               string
	host               string
	serverName         string
	peerCommonName     string
	contentType        string
	exchangeFields     map[string]string
	authorization      string
	proxyAuthorization string
}

type x509ConformanceServer struct {
	server     *httptest.Server
	mu         sync.Mutex
	records    []x509ConformanceRequest
	handshakes int
	redirectTo string
}

func (lab *x509ConformanceLab) server(t *testing.T, hostname string, exchange bool) *x509ConformanceServer {
	t.Helper()
	observed := &x509ConformanceServer{}
	observed.server = httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var exchangeFields map[string]string
		if exchange {
			if err := json.NewDecoder(req.Body).Decode(&exchangeFields); err != nil {
				http.Error(w, "invalid X.509 exchange request", http.StatusBadRequest)
				return
			}
		}
		record := x509ConformanceRequest{
			method:             req.Method,
			path:               req.URL.Path,
			host:               req.Host,
			serverName:         req.TLS.ServerName,
			peerCommonName:     req.TLS.PeerCertificates[0].Subject.CommonName,
			contentType:        req.Header.Get("Content-Type"),
			exchangeFields:     exchangeFields,
			authorization:      req.Header.Get("Authorization"),
			proxyAuthorization: req.Header.Get("Proxy-Authorization"),
		}
		observed.mu.Lock()
		observed.records = append(observed.records, record)
		redirectTo := observed.redirectTo
		observed.mu.Unlock()
		if redirectTo != "" {
			http.Redirect(w, req, redirectTo, http.StatusTemporaryRedirect)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		body := map[string]any{"data": []any{}}
		if exchange {
			body = map[string]any{
				"access_token":      x509ConformanceToken,
				"token_type":        "Bearer",
				"issued_token_type": "urn:ietf:params:oauth:token-type:access_token",
				"expires_in":        3600,
			}
		}
		if err := json.NewEncoder(w).Encode(body); err != nil {
			return
		}
	}))
	serverLeaf, serverKey := lab.issue(t, hostname, lab.intermediate, lab.issuerKey, false, []string{hostname})
	observed.server.TLS = &tls.Config{
		MinVersion: tls.VersionTLS12,
		ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs:  lab.trust,
		Certificates: []tls.Certificate{{
			Certificate: [][]byte{serverLeaf.Raw, lab.intermediate.Raw},
			PrivateKey:  serverKey,
			Leaf:        serverLeaf,
		}},
		GetConfigForClient: func(*tls.ClientHelloInfo) (*tls.Config, error) {
			observed.mu.Lock()
			observed.handshakes++
			observed.mu.Unlock()
			return nil, nil
		},
	}
	observed.server.Config.ErrorLog = log.New(io.Discard, "", 0)
	observed.server.StartTLS()
	t.Cleanup(observed.server.Close)
	return observed
}

func (server *x509ConformanceServer) requests() []x509ConformanceRequest {
	server.mu.Lock()
	defer server.mu.Unlock()
	return append([]x509ConformanceRequest(nil), server.records...)
}

func (server *x509ConformanceServer) setRedirect(destination string) {
	server.mu.Lock()
	defer server.mu.Unlock()
	server.redirectTo = destination
}

func (server *x509ConformanceServer) handshakeCount() int {
	server.mu.Lock()
	defer server.mu.Unlock()
	return server.handshakes
}

func x509ConformanceRoutes(issuer, api *x509ConformanceServer) map[string]string {
	return map[string]string{
		x509ConformanceIssuerHost + ":443": issuer.server.Listener.Addr().String(),
		x509ConformanceAPIHost + ":443":    api.server.Listener.Addr().String(),
	}
}

type x509ConformanceHostTransport map[string]*http.Transport

func (transports x509ConformanceHostTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	transport, ok := transports[req.URL.Hostname()]
	if !ok {
		return nil, fmt.Errorf("unexpected hermetic mTLS host %q", req.URL.Hostname())
	}
	return transport.RoundTrip(req)
}

type x509ConformanceConnectRequest struct {
	authority          string
	authorization      string
	proxyAuthorization string
}

type x509ConformanceProxy struct {
	server  *httptest.Server
	url     *url.URL
	mu      sync.Mutex
	records []x509ConformanceConnectRequest
}

func newX509ConformanceProxy(t *testing.T, routes map[string]string) *x509ConformanceProxy {
	t.Helper()
	proxy := &x509ConformanceProxy{}
	proxy.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodConnect {
			http.Error(w, "CONNECT required", http.StatusMethodNotAllowed)
			return
		}
		proxy.mu.Lock()
		proxy.records = append(proxy.records, x509ConformanceConnectRequest{
			authority:          req.Host,
			authorization:      req.Header.Get("Authorization"),
			proxyAuthorization: req.Header.Get("Proxy-Authorization"),
		})
		proxy.mu.Unlock()

		target, ok := routes[req.Host]
		if !ok {
			http.Error(w, "unexpected CONNECT authority", http.StatusBadGateway)
			return
		}
		dialer := &net.Dialer{Timeout: 5 * time.Second}
		upstream, err := dialer.DialContext(req.Context(), "tcp", target)
		if err != nil {
			http.Error(w, "upstream unavailable", http.StatusBadGateway)
			return
		}
		defer func() { _ = upstream.Close() }()

		hijacker, ok := w.(http.Hijacker)
		if !ok {
			http.Error(w, "CONNECT hijacking unavailable", http.StatusInternalServerError)
			return
		}
		connection, buffered, err := hijacker.Hijack()
		if err != nil {
			return
		}
		defer func() { _ = connection.Close() }()
		if _, err := buffered.WriteString("HTTP/1.1 200 Connection Established\r\n\r\n"); err != nil {
			return
		}
		if err := buffered.Flush(); err != nil {
			return
		}

		finished := make(chan struct{})
		go func() {
			defer close(finished)
			_, _ = io.Copy(upstream, buffered)
			_ = upstream.Close()
		}()
		_, _ = io.Copy(connection, upstream)
		_ = connection.Close()
		<-finished
	}))
	t.Cleanup(proxy.server.Close)

	parsed, err := url.Parse(proxy.server.URL)
	if err != nil {
		t.Fatalf("parse hermetic CONNECT proxy URL: %v", err)
	}
	proxy.url = parsed
	return proxy
}

func (proxy *x509ConformanceProxy) requests() []x509ConformanceConnectRequest {
	proxy.mu.Lock()
	defer proxy.mu.Unlock()
	return append([]x509ConformanceConnectRequest(nil), proxy.records...)
}

package openai_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
	"testing"

	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/auth"
	"github.com/openai/openai-go/v3/option"
)

// These fixtures handle every request in memory, including subject tokens and
// OAuth exchange. They never use live credentials or contact external services.
type isolatedSubjectTokenProvider struct{}

func (isolatedSubjectTokenProvider) TokenType() auth.SubjectTokenType { return auth.SubjectTokenTypeID }

func (isolatedSubjectTokenProvider) GetToken(ctx context.Context, client auth.HTTPDoer) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://subject.example.invalid/token", nil)
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()
	token, err := io.ReadAll(resp.Body)
	return string(token), err
}

type isolatedWorkloadDoer struct {
	t         *testing.T
	identity  string
	exchanges atomic.Int32
	metadata  atomic.Int32
	apiCalls  atomic.Int32
	reject    atomic.Bool
}

func (d *isolatedWorkloadDoer) Do(req *http.Request) (*http.Response, error) {
	switch req.URL.Host {
	case "subject.example.invalid":
		d.metadata.Add(1)
		return rootWorkloadResponse(http.StatusOK, "synthetic-subject-"+d.identity), nil
	case "auth.openai.com":
		d.exchanges.Add(1)
		var body struct {
			SubjectToken string `json:"subject_token"`
		}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return nil, err
		}
		if body.SubjectToken != "synthetic-subject-"+d.identity {
			d.t.Error("exchange used a different transport's subject token")
		}
		return rootWorkloadResponse(http.StatusOK,
			fmt.Sprintf(`{"access_token":"synthetic-bearer-%s","expires_in":3600}`, d.identity)), nil
	case "api.openai.com":
		d.apiCalls.Add(1)
		if req.Header.Get("Authorization") != "Bearer synthetic-bearer-"+d.identity {
			d.t.Error("API request did not use this transport's bearer")
		}
		if d.reject.Swap(false) {
			return rootWorkloadResponse(http.StatusUnauthorized, `{"error":{"message":"synthetic rejection"}}`), nil
		}
		return rootWorkloadResponse(http.StatusOK, `{"data":[]}`), nil
	default:
		return nil, fmt.Errorf("unexpected synthetic request host %q", req.URL.Host)
	}
}

func (d *isolatedWorkloadDoer) RoundTrip(req *http.Request) (*http.Response, error) {
	defer func() {
		if req.Body != nil {
			_ = req.Body.Close()
		}
	}()
	return d.Do(req)
}

type comparableWorkloadDoer struct{ *isolatedWorkloadDoer }
type sliceWorkloadDoer struct {
	*isolatedWorkloadDoer
	values []string
}
type interfaceWorkloadDoer struct {
	*isolatedWorkloadDoer
	value any
}
type functionWorkloadDoer func(*http.Request) (*http.Response, error)

func (d functionWorkloadDoer) Do(req *http.Request) (*http.Response, error) { return d(req) }

func TestWorkloadIdentityIsolatesReusedOptionsAndRequestClients(t *testing.T) {
	for _, test := range []struct {
		name   string
		wrap   func(*isolatedWorkloadDoer) option.HTTPClient
		cached bool
	}{
		{"native client", func(d *isolatedWorkloadDoer) option.HTTPClient { return &http.Client{Transport: d} }, true},
		{"custom pointer", func(d *isolatedWorkloadDoer) option.HTTPClient { return d }, true},
		{"comparable value", func(d *isolatedWorkloadDoer) option.HTTPClient { return comparableWorkloadDoer{d} }, true},
		{"slice value", func(d *isolatedWorkloadDoer) option.HTTPClient { return sliceWorkloadDoer{d, []string{"synthetic"}} }, false},
		{"interface containing slice", func(d *isolatedWorkloadDoer) option.HTTPClient { return interfaceWorkloadDoer{d, []int{1}} }, false},
		{"function value", func(d *isolatedWorkloadDoer) option.HTTPClient { return functionWorkloadDoer(d.Do) }, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			for _, methodOverride := range []bool{false, true} {
				t.Run(fmt.Sprintf("method override=%t", methodOverride), func(t *testing.T) {
					a := &isolatedWorkloadDoer{t: t, identity: "a"}
					b := &isolatedWorkloadDoer{t: t, identity: "b"}
					aHTTP, bHTTP := test.wrap(a), test.wrap(b)
					identity := option.WithWorkloadIdentity(auth.WorkloadIdentity{
						IdentityProviderID: "synthetic-provider", ServiceAccountID: "synthetic-account",
						Provider: isolatedSubjectTokenProvider{},
					})
					// Exercise both client-option orders as well as a method override.
					clientA := openai.NewClient(identity, option.WithHTTPClient(aHTTP), option.WithMaxRetries(0))
					clientB := openai.NewClient(option.WithHTTPClient(bHTTP), identity, option.WithMaxRetries(0))
					for range 2 {
						if _, err := clientA.Models.List(t.Context()); err != nil {
							t.Fatal(err)
						}
						var err error
						if methodOverride {
							_, err = clientA.Models.List(t.Context(), option.WithHTTPClient(bHTTP))
						} else {
							_, err = clientB.Models.List(t.Context())
						}
						if err != nil {
							t.Fatal(err)
						}
					}
					want := int32(2)
					if test.cached {
						want = 1
					}
					for _, d := range []*isolatedWorkloadDoer{a, b} {
						if d.exchanges.Load() != want || d.metadata.Load() != want || d.apiCalls.Load() != 2 {
							t.Errorf("identity %s exchanges/metadata/API=%d/%d/%d, want %d/%d/2",
								d.identity, d.exchanges.Load(), d.metadata.Load(), d.apiCalls.Load(), want, want)
						}
					}
				})
			}
		})
	}
}

func TestWorkloadIdentityIsolatesTransportReplacement(t *testing.T) {
	a := &isolatedWorkloadDoer{t: t, identity: "a"}
	b := &isolatedWorkloadDoer{t: t, identity: "b"}
	native := &http.Client{Transport: a}
	client := openai.NewClient(option.WithHTTPClient(native), option.WithMaxRetries(0),
		option.WithWorkloadIdentity(auth.WorkloadIdentity{
			IdentityProviderID: "synthetic-provider", ServiceAccountID: "synthetic-account",
			Provider: isolatedSubjectTokenProvider{},
		}))
	// Changes occur between completed calls, never concurrently with client use.
	for _, transport := range []*isolatedWorkloadDoer{a, b, a} {
		native.Transport = transport
		if _, err := client.Models.List(t.Context()); err != nil {
			t.Fatal(err)
		}
	}
	if a.exchanges.Load() != 1 || b.exchanges.Load() != 1 {
		t.Fatal("transport replacement did not keep independent reusable caches")
	}
}

func TestWorkloadIdentityUnauthorizedInvalidationStaysWithinIdentity(t *testing.T) {
	// Equal bearer bytes do not make two client identities interchangeable.
	a := &isolatedWorkloadDoer{t: t, identity: "same"}
	b := &isolatedWorkloadDoer{t: t, identity: "same"}
	identity := option.WithWorkloadIdentity(auth.WorkloadIdentity{
		IdentityProviderID: "synthetic-provider", ServiceAccountID: "synthetic-account",
		Provider: isolatedSubjectTokenProvider{},
	})
	client := openai.NewClient(identity, option.WithHTTPClient(a), option.WithMaxRetries(0))
	for _, doer := range []*isolatedWorkloadDoer{a, b} {
		if _, err := client.Models.List(t.Context(), option.WithHTTPClient(doer)); err != nil {
			t.Fatal(err)
		}
	}
	a.reject.Store(true)
	if _, err := client.Models.List(t.Context()); err == nil {
		t.Fatal("expected synthetic unauthorized response")
	}
	for _, doer := range []*isolatedWorkloadDoer{b, a} {
		if _, err := client.Models.List(t.Context(), option.WithHTTPClient(doer)); err != nil {
			t.Fatal(err)
		}
	}
	if a.exchanges.Load() != 2 || b.exchanges.Load() != 1 {
		t.Errorf("issuer exchanges A/B=%d/%d, want 2/1", a.exchanges.Load(), b.exchanges.Load())
	}
}

func TestWorkloadIdentityKeepsCertificateTransportsSeparate(t *testing.T) {
	lab := newX509ConformanceLab(t)
	issuer := lab.server(t, "auth.openai.com", true)
	api := lab.server(t, "api.openai.com", false)
	routes := map[string]string{
		"auth.openai.com:443": issuer.server.Listener.Addr().String(),
		"api.openai.com:443":  api.server.Listener.Addr().String(),
	}
	a := &http.Client{Transport: lab.transport(t, routes, lab.identity(t, "identity-a", true))}
	b := &http.Client{Transport: lab.transport(t, routes, lab.identity(t, "identity-b", true))}
	identity := option.WithWorkloadIdentity(testWorkloadIdentity(&mockSubjectTokenProvider{
		token: "synthetic-subject", tokenType: auth.SubjectTokenTypeJWT,
	}))
	client := openai.NewClient(identity, option.WithHTTPClient(a), option.WithMaxRetries(0))
	for _, transport := range []*http.Client{a, b, a, b} {
		if _, err := client.Models.List(t.Context(), option.WithHTTPClient(transport)); err != nil {
			t.Fatal(err)
		}
	}
	exchanges, requests := issuer.requests(), api.requests()
	if len(exchanges) != 2 || len(requests) != 4 {
		t.Fatalf("TLS exchanges/API calls=%d/%d, want 2/4", len(exchanges), len(requests))
	}
	for index, name := range []string{"identity-a", "identity-b"} {
		if exchanges[index].peerCommonName != name {
			t.Error("issuer did not observe the selected client certificate")
		}
		for _, request := range []x509ConformanceRequest{requests[index], requests[index+2]} {
			if request.peerCommonName != name || request.authorization != "Bearer "+x509ConformanceToken {
				t.Error("API did not observe the selected certificate and synthetic bearer")
			}
		}
	}
}

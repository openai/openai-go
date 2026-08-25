package auth

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewX509WorkloadIdentityAuthValidatesConfiguration(t *testing.T) {
	fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, x509ValidExchangeResponse())
	}))
	valid := X509WorkloadIdentity{
		IdentityProviderID: "synthetic-identity-provider",
		ServiceAccountID:   "synthetic-service-account",
		Transport:          fixture.capability,
	}
	for _, test := range []struct {
		name   string
		change func(*X509WorkloadIdentity)
	}{
		{name: "missing identity provider", change: func(config *X509WorkloadIdentity) { config.IdentityProviderID = "" }},
		{name: "blank identity provider", change: func(config *X509WorkloadIdentity) { config.IdentityProviderID = "  " }},
		{name: "missing service account", change: func(config *X509WorkloadIdentity) { config.ServiceAccountID = "" }},
		{name: "blank service account", change: func(config *X509WorkloadIdentity) { config.ServiceAccountID = "\t" }},
		{name: "missing transport", change: func(config *X509WorkloadIdentity) { config.Transport = nil }},
		{name: "negative refresh buffer", change: func(config *X509WorkloadIdentity) { config.RefreshBuffer = -time.Second }},
		{name: "excessive refresh buffer", change: func(config *X509WorkloadIdentity) { config.RefreshBuffer = time.Hour }},
	} {
		t.Run(test.name, func(t *testing.T) {
			config := valid
			test.change(&config)
			identity, err := NewX509WorkloadIdentityAuth(config)
			if identity != nil || err == nil {
				t.Fatalf("invalid X.509 workload configuration returned identity=%v error=%v", identity, err)
			}
		})
	}
	if _, err := NewX509WorkloadIdentityAuth(valid); err != nil {
		t.Fatalf("valid X.509 workload configuration: %v", err)
	}
}

func TestX509WorkloadIdentityAuthExchangesOnlyOnItsAttestedTransport(t *testing.T) {
	var exchanges atomic.Int32
	fixture := newX509ExchangeFixture(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		exchanges.Add(1)
		_, _ = io.WriteString(w, x509ValidExchangeResponse())
	}))
	identity, err := NewX509WorkloadIdentityAuth(X509WorkloadIdentity{
		IdentityProviderID: "synthetic-identity-provider",
		ServiceAccountID:   "synthetic-service-account",
		Transport:          fixture.capability,
	})
	if err != nil {
		t.Fatalf("construct X.509 workload identity: %v", err)
	}
	for range 2 {
		token, tokenErr := identity.GetToken(t.Context(), fixture.capability)
		if tokenErr != nil || token != x509ExchangeSyntheticToken {
			t.Fatalf("exchange on attested transport token=%q error=%v", token, tokenErr)
		}
	}
	if got := exchanges.Load(); got != 1 {
		t.Errorf("cached workload identity made %d exchanges, want one per transport generation", got)
	}
	for _, doer := range []HTTPDoer{nil, http.DefaultClient} {
		if token, tokenErr := identity.GetToken(t.Context(), doer); token != "" || tokenErr == nil {
			t.Errorf("unattested HTTP client returned token=%q error=%v", token, tokenErr)
		}
	}
	another, err := NewX509Transport(fixture.template)
	if err != nil {
		t.Fatalf("attest a separate certificate generation: %v", err)
	}
	t.Cleanup(func() { _ = another.Close() })
	if token, tokenErr := identity.GetToken(t.Context(), another); token != "" || tokenErr == nil {
		t.Errorf("different attested generation returned token=%q error=%v", token, tokenErr)
	}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	if token, tokenErr := identity.GetToken(ctx, fixture.capability); token != "" || !errors.Is(tokenErr, context.Canceled) {
		t.Errorf("canceled X.509 exchange returned token=%q error=%v", token, tokenErr)
	}
	var invalid *X509WorkloadIdentityAuth
	if token, tokenErr := invalid.GetToken(t.Context(), fixture.capability); token != "" ||
		tokenErr == nil || !strings.Contains(tokenErr.Error(), "invalid") {
		t.Errorf("nil X.509 identity returned token=%q error=%v", token, tokenErr)
	}
	if got := exchanges.Load(); got != 1 {
		t.Errorf("rejected transports/cancellation caused unexpected exchanges: %d", got)
	}
}

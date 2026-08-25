package auth

import (
	"context"
	"errors"
	"strings"
	"time"
)

// X509WorkloadIdentity configures an OpenAI workload authenticated with a
// caller-owned X.509 certificate and an explicitly attested native transport.
type X509WorkloadIdentity struct {
	IdentityProviderID string
	ServiceAccountID   string
	RefreshBuffer      time.Duration
	Transport          *X509Transport
}

// X509WorkloadIdentityAuth exchanges a static, attested X.509 workload identity
// for an ordinary OpenAI bearer token.
type X509WorkloadIdentityAuth struct {
	config X509WorkloadIdentity
}

// NewX509WorkloadIdentityAuth validates and binds a workload identity to one
// attested X.509 transport generation.
func NewX509WorkloadIdentityAuth(config X509WorkloadIdentity) (*X509WorkloadIdentityAuth, error) {
	if strings.TrimSpace(config.IdentityProviderID) == "" {
		return nil, errors.New("X.509 workload identity requires an identity-provider ID")
	}
	if strings.TrimSpace(config.ServiceAccountID) == "" {
		return nil, errors.New("X.509 workload identity requires a service-account ID")
	}
	if config.RefreshBuffer < 0 || config.RefreshBuffer >= x509MaximumTokenLifetime*time.Second {
		return nil, errors.New("X.509 workload identity requires a non-negative refresh buffer shorter than one hour")
	}
	if err := config.Transport.validateAttestation(); err != nil {
		return nil, err
	}
	return &X509WorkloadIdentityAuth{config: config}, nil
}

// GetToken exchanges the workload identity over its exact configured transport.
// A different HTTP client or transport cannot be substituted.
func (identity *X509WorkloadIdentityAuth) GetToken(ctx context.Context, doer HTTPDoer) (string, error) {
	if identity == nil {
		return "", errors.New("X.509 workload identity authentication is invalid")
	}
	transport, ok := doer.(*X509Transport)
	if !ok || transport != identity.config.Transport {
		return "", errors.New("X.509 workload identity requires its exact attested transport")
	}
	token, err := x509Exchange(ctx, transport, identity.config.IdentityProviderID, identity.config.ServiceAccountID)
	if err != nil {
		return "", err
	}
	return token.value, nil
}

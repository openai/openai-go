package auth

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"
)

const (
	x509DefaultRefreshBuffer = 5 * time.Minute
	x509MaximumAttempts      = 3
	x509InitialRetryDelay    = 25 * time.Millisecond
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
	config       X509WorkloadIdentity
	mu           sync.Mutex
	cached       x509ExchangedToken
	refreshAfter time.Time
	inFlight     *x509TokenRefresh
}

type x509TokenRefresh struct {
	done chan struct{}
	err  error
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
	if ctx == nil {
		return "", errors.New("X.509 workload identity requires a non-nil context")
	}
	if err := ctx.Err(); err != nil {
		return "", err
	}
	transport, ok := doer.(*X509Transport)
	if !ok || transport != identity.config.Transport {
		return "", errors.New("X.509 workload identity requires its exact attested transport")
	}
	if err := transport.validateAttestation(); err != nil {
		return "", err
	}
	for {
		if err := ctx.Err(); err != nil {
			return "", err
		}
		identity.mu.Lock()
		if identity.cached.value != "" && time.Now().Before(identity.refreshAfter) {
			token := identity.cached.value
			identity.mu.Unlock()
			return token, nil
		}
		if current := identity.inFlight; current != nil {
			identity.mu.Unlock()
			select {
			case <-ctx.Done():
				return "", ctx.Err()
			case <-current.done:
				if current.err == nil || errors.Is(current.err, context.Canceled) ||
					errors.Is(current.err, context.DeadlineExceeded) {
					continue
				}
				return "", current.err
			}
		}
		refresh := &x509TokenRefresh{done: make(chan struct{})}
		identity.inFlight = refresh
		identity.mu.Unlock()

		token, err := identity.exchangeWithRetry(ctx, transport)
		identity.mu.Lock()
		if err == nil {
			identity.cached = token
			buffer := identity.config.RefreshBuffer
			if buffer == 0 {
				buffer = x509DefaultRefreshBuffer
			}
			buffer = min(buffer, time.Until(token.expiresAt)/2)
			identity.refreshAfter = token.expiresAt.Add(-buffer)
		}
		refresh.err = err
		identity.inFlight = nil
		close(refresh.done)
		identity.mu.Unlock()
		if err != nil {
			return "", err
		}
		return token.value, nil
	}
}

func (identity *X509WorkloadIdentityAuth) exchangeWithRetry(ctx context.Context, transport *X509Transport) (x509ExchangedToken, error) {
	for attempt := 0; ; attempt++ {
		token, err := x509Exchange(ctx, transport, identity.config.IdentityProviderID, identity.config.ServiceAccountID)
		if err == nil || !retryableX509ExchangeError(err) || attempt+1 >= x509MaximumAttempts {
			return token, err
		}
		timer := time.NewTimer(x509InitialRetryDelay << attempt)
		select {
		case <-ctx.Done():
			if !timer.Stop() {
				<-timer.C
			}
			return x509ExchangedToken{}, ctx.Err()
		case <-timer.C:
		}
	}
}

func retryableX509ExchangeError(err error) bool {
	var status *x509ExchangeHTTPError
	if errors.As(err, &status) {
		return status.retryable()
	}
	var read *x509ExchangeReadError
	return errors.As(err, &read) && read.retryable()
}

func (identity *X509WorkloadIdentityAuth) invalidateToken(value string) {
	identity.mu.Lock()
	defer identity.mu.Unlock()
	if identity.cached.value == value {
		identity.cached = x509ExchangedToken{}
		identity.refreshAfter = time.Time{}
	}
}

package auth

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/openai/openai-go/v3/internal/requestconfig"
)

const (
	x509DefaultRefreshBuffer        = 5 * time.Minute
	x509DefaultTokenExchangeTimeout = 30 * time.Second
	x509MaximumAttempts             = 3
	x509InitialRetryDelay           = 25 * time.Millisecond
	x509MaximumRetryCooldown        = 5 * time.Second
)

var (
	errX509InvalidatedBearer    = errors.New("X.509 token exchange returned an invalidated bearer")
	errX509TokenExchangeTimeout = errors.New("X.509 token exchange timed out")
)

// X509WorkloadIdentity configures an OpenAI workload authenticated with a
// caller-owned X.509 certificate and an explicitly attested native transport.
type X509WorkloadIdentity struct {
	// IdentityProviderID identifies the enrolled workload identity provider.
	IdentityProviderID string
	// ServiceAccountID identifies the enrolled workload service account.
	ServiceAccountID string
	// RefreshBuffer controls how early cached bearer tokens are refreshed. Zero
	// uses a five-minute default, capped at half the token's remaining lifetime.
	RefreshBuffer time.Duration
	// Transport is the exact attested transport used for token exchange and API
	// requests.
	Transport *X509Transport
}

// X509WorkloadIdentityAuth exchanges a static, attested X.509 workload identity
// for an ordinary OpenAI bearer token.
type X509WorkloadIdentityAuth struct {
	config               X509WorkloadIdentity
	tokenExchangeTimeout time.Duration
	mu                   sync.Mutex
	cached               x509ExchangedToken
	bearers              *workloadBearerHistory
	refreshAfter         time.Time
	inFlight             *x509TokenRefresh
}

type x509TokenRefresh struct {
	done            chan struct{}
	wake            chan struct{}
	wakeOnce        sync.Once
	generation      string
	ownerContextErr error
	err             error
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
	return &X509WorkloadIdentityAuth{
		config:               config,
		tokenExchangeTimeout: x509DefaultTokenExchangeTimeout,
	}, nil
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
	if identity.mu.TryLock() {
		now := time.Now()
		if identity.cached.value != "" && now.Before(identity.refreshAfter) && now.Before(identity.cached.expiresAt) {
			token := identity.cached.value
			identity.mu.Unlock()
			return token, nil
		}
		identity.mu.Unlock()
	}
	callerCtx := ctx
	ctx, cancel := identity.exchangeContext(ctx)
	defer cancel()
	for {
		if ctx.Err() != nil {
			return identity.tokenAfterExchangeContextDone(ctx, callerCtx)
		}
		identity.mu.Lock()
		now := time.Now()
		if identity.cached.value != "" && now.Before(identity.refreshAfter) && now.Before(identity.cached.expiresAt) {
			token := identity.cached.value
			identity.mu.Unlock()
			return token, nil
		}
		if current := identity.inFlight; current != nil {
			identity.mu.Unlock()
			select {
			case <-ctx.Done():
				return identity.tokenAfterExchangeContextDone(ctx, callerCtx)
			case <-current.done:
				if current.err == nil || current.ownerContextErr != nil &&
					errors.Is(current.err, current.ownerContextErr) {
					continue
				}
				if retryableX509ExchangeError(current.err) {
					if scope := requestconfig.RequestRetryScopeFromContext(ctx); scope != nil && scope.TryRetry() {
						continue
					}
				}
				return "", current.err
			}
		}
		refresh := &x509TokenRefresh{
			done:       make(chan struct{}),
			wake:       make(chan struct{}),
			generation: identity.cached.value,
		}
		identity.inFlight = refresh
		identity.mu.Unlock()

		token, err := identity.exchangeWithRetry(ctx, transport, refresh)
		identity.mu.Lock()
		refresh.ownerContextErr = ctx.Err()
		now = time.Now()
		if err == nil && identity.bearers.isRejected(token.value, now) {
			err = errX509InvalidatedBearer
		}
		fallback := false
		if err != nil && x509CanFallBackToCachedToken(err, ctx, callerCtx) && identity.cached.value != "" &&
			time.Now().Before(identity.cached.expiresAt) {
			cooldown := min(x509MaximumRetryCooldown, time.Until(identity.cached.expiresAt)/2)
			identity.refreshAfter = time.Now().Add(cooldown)
			token, err = identity.cached, nil
			fallback = true
		}
		if err == nil && !fallback {
			remaining := time.Until(token.expiresAt)
			if remaining <= 0 {
				err = errors.New("X.509 token exchange returned an already expired token")
			} else {
				if historyErr := ensureWorkloadBearerHistory(&identity.bearers).
					recordIssued(token.value, token.expiresAt, now); historyErr != nil {
					err = historyErr
				} else {
					identity.cached = token
					buffer := identity.config.RefreshBuffer
					if buffer == 0 {
						buffer = x509DefaultRefreshBuffer
					}
					buffer = min(buffer, remaining/2)
					identity.refreshAfter = token.expiresAt.Add(-buffer)
				}
			}
		}
		if err == nil && !time.Now().Before(token.expiresAt) {
			err = errors.New("X.509 token exchange returned an already expired token")
			if identity.cached.value == token.value {
				identity.cached = x509ExchangedToken{}
				identity.refreshAfter = time.Time{}
			}
		}
		refresh.err = err
		identity.inFlight = nil
		close(refresh.done)
		identity.mu.Unlock()
		if err != nil {
			return "", err
		}
		if fallback {
			if callerErr := callerCtx.Err(); callerErr != nil {
				return "", callerErr
			}
		}
		return token.value, nil
	}
}

func (identity *X509WorkloadIdentityAuth) exchangeContext(ctx context.Context) (context.Context, context.CancelFunc) {
	if _, hasDeadline := ctx.Deadline(); hasDeadline {
		return context.WithCancel(ctx)
	}
	return context.WithTimeoutCause(ctx, identity.tokenExchangeTimeout, errX509TokenExchangeTimeout)
}

func (identity *X509WorkloadIdentityAuth) tokenAfterExchangeContextDone(
	exchangeCtx, callerCtx context.Context,
) (string, error) {
	err := exchangeCtx.Err()
	if callerCtx.Err() != nil || !errors.Is(context.Cause(exchangeCtx), errX509TokenExchangeTimeout) {
		return "", err
	}
	identity.mu.Lock()
	defer identity.mu.Unlock()
	if callerErr := callerCtx.Err(); callerErr != nil {
		return "", callerErr
	}
	if identity.cached.value == "" || identity.bearers.isRejected(identity.cached.value, time.Now()) ||
		!time.Now().Before(identity.cached.expiresAt) {
		return "", err
	}
	return identity.cached.value, nil
}

func x509CanFallBackToCachedToken(err error, exchangeCtx, callerCtx context.Context) bool {
	if callerCtx.Err() != nil {
		return false
	}
	if retryableX509ExchangeError(err) {
		return true
	}
	return errors.Is(err, context.DeadlineExceeded) &&
		errors.Is(context.Cause(exchangeCtx), errX509TokenExchangeTimeout)
}

func (identity *X509WorkloadIdentityAuth) exchangeWithRetry(
	ctx context.Context, transport *X509Transport, refresh *x509TokenRefresh,
) (x509ExchangedToken, error) {
	wake := refresh.wake
	for attempt := 0; ; attempt++ {
		token, err := x509Exchange(ctx, transport, identity.config.IdentityProviderID, identity.config.ServiceAccountID)
		if err == nil {
			identity.mu.Lock()
			if identity.bearers.isRejected(token.value, time.Now()) {
				err = errX509InvalidatedBearer
			}
			identity.mu.Unlock()
		}
		if err == nil || !retryableX509ExchangeError(err) || attempt+1 >= x509MaximumAttempts {
			return token, err
		}
		scope := requestconfig.RequestRetryScopeFromContext(ctx)
		if scope != nil && !scope.TryRetry() {
			return token, err
		}
		if errors.Is(err, errX509InvalidatedBearer) {
			continue
		}
		timer := time.NewTimer(x509ExchangeRetryDelay(err, attempt, scope))
		select {
		case <-ctx.Done():
			if !timer.Stop() {
				<-timer.C
			}
			return x509ExchangedToken{}, ctx.Err()
		case <-timer.C:
		case <-wake:
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
			wake = nil
		}
	}
}

func x509ExchangeRetryDelay(err error, attempt int, scope *requestconfig.RequestRetryScope) time.Duration {
	delay := x509InitialRetryDelay << attempt
	var status *x509ExchangeHTTPError
	if errors.As(err, &status) && status.hasRetryAfter {
		delay = status.retryAfter
	}
	maximum := x509MaximumRetryAfter
	if scope != nil {
		maximum = scope.MaxRetryDelay()
	}
	return min(delay, maximum)
}

func retryableX509ExchangeError(err error) bool {
	if errors.Is(err, errX509InvalidatedBearer) {
		return true
	}
	var oauth *OAuthError
	if errors.As(err, &oauth) {
		return oauth.ErrorCode == "temporarily_unavailable" || oauth.ErrorCode == "server_error"
	}
	var status *x509ExchangeHTTPError
	if errors.As(err, &status) {
		return status.retryable()
	}
	var read *x509ExchangeReadError
	if errors.As(err, &read) {
		return read.retryable()
	}
	var transport *x509TransportError
	return errors.As(err, &transport) && transport.retryable
}

func (identity *X509WorkloadIdentityAuth) invalidateToken(value string) {
	identity.mu.Lock()
	defer identity.mu.Unlock()
	identity.invalidateTokenLocked(value)
}

func (identity *X509WorkloadIdentityAuth) invalidateTokenLocked(value string) {
	now := time.Now()
	knownExpiry := time.Time{}
	if identity.cached.value == value {
		knownExpiry = identity.cached.expiresAt
	}
	ensureWorkloadBearerHistory(&identity.bearers).reject(value, knownExpiry, now)
	if identity.cached.value == value {
		identity.cached = x509ExchangedToken{}
		identity.refreshAfter = time.Time{}
		if current := identity.inFlight; current != nil && current.generation == value {
			current.wakeOnce.Do(func() {
				if current.wake != nil {
					close(current.wake)
				}
			})
		}
	}
}

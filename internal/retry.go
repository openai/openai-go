package internal

import "errors"

// noRetryError marks deterministic request setup or policy failures that
// cannot be fixed by replaying the same request.
type noRetryError struct {
	err error
}

func (e *noRetryError) Error() string { return e.err.Error() }
func (e *noRetryError) Unwrap() error { return e.err }

// WithNoRetryError marks err as deterministic for the generic request retry
// loop while preserving its message and unwrap chain.
func WithNoRetryError(err error) error {
	if err == nil {
		return nil
	}
	return &noRetryError{err: err}
}

// IsNoRetryError reports whether err contains the SDK's private retry marker.
func IsNoRetryError(err error) bool {
	var marker *noRetryError
	return errors.As(err, &marker)
}

package auth_test

import "github.com/openai/openai-go/v3/auth"

// Lock the released field order and arity against accidental changes.
var _ = auth.X509WorkloadIdentity{"identity-provider", "service-account", 0, nil}

// Private revocation state must not make the exported authentication handles
// unusable as map keys or in interface equality checks.
var (
	_ map[auth.WorkloadIdentityAuth]struct{}
	_ map[auth.X509WorkloadIdentityAuth]struct{}
)

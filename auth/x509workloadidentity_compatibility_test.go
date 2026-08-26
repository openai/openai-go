package auth_test

import "github.com/openai/openai-go/v3/auth"

// Preserve the original exported field order and arity for callers that used a
// positional literal before the X.509 configuration contract could be hardened.
var _ = auth.X509WorkloadIdentity{"identity-provider", "service-account", 0, nil}

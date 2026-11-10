package option

import "github.com/openai/openai-go/v3/internal/requestconfig"

// WithUnsafeAllowHTTP permits OpenAI credentials over plaintext HTTP only to
// localhost or a literal loopback IP address. Use it only for local development.
// These requests use a dedicated direct connection, bypassing proxies, custom
// transports, dialers, and HTTP doers. HTTPS keeps the configured HTTP client.
// Azure and Bedrock transport policies are unaffected.
func WithUnsafeAllowHTTP() RequestOption {
	return requestconfig.WithUnsafeAllowHTTP()
}

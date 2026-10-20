package option

import "github.com/openai/openai-go/v3/internal/requestconfig"

// WithUnsafeAllowHTTP was introduced in v3.69.0 with authenticated HTTP restrictions.
// This was reverted as it broke existing tests and integrations. The restrictions will
// return in the next major release. This option is currently a no-op retained for source
// compatibility. Azure and Bedrock retain their provider-specific policies
func WithUnsafeAllowHTTP() RequestOption {
	return requestconfig.RequestOptionFunc(func(*requestconfig.RequestConfig) error { return nil })
}

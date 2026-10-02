package option

import "github.com/openai/openai-go/v3/internal/requestconfig"

// WithUnsafeAllowHTTP is retained for source compatibility.
//
// Deprecated: OpenAI HTTP endpoints no longer require an opt-in. This option
// has no effect. Azure and Bedrock retain their provider-specific policies.
func WithUnsafeAllowHTTP() RequestOption {
	return requestconfig.RequestOptionFunc(func(*requestconfig.RequestConfig) error { return nil })
}

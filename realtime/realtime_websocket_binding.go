// File generated from our OpenAPI spec by Castiron. See CONTRIBUTING.md for details.

package realtime

import (
	"context"
	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
	"net/http"
	"slices"
)

// Capture request options once, then prepare a fresh request for every attempt.
func (r *RealtimeService) newWebSocketRequest(opts ...option.RequestOption) func(context.Context) (*requestconfig.RequestConfig, error) {
	requestOptions := slices.Concat([]option.RequestOption{requestconfig.WithBearerAuthSecurity()}, r.Options, opts)
	return func(ctx context.Context) (*requestconfig.RequestConfig, error) {
		return requestconfig.NewRequestConfig(ctx, http.MethodGet, "realtime", nil, nil, requestOptions...)
	}
}

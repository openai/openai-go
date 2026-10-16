// File generated from our OpenAPI spec by Castiron. See CONTRIBUTING.md for details.

package live

import (
	"context"
	"errors"
	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
	"net/http"
	"slices"
)

// Capture request options once, then prepare a fresh request for every attempt.
func (r *ForkService) newWebSocketRequest(sessionID string, opts ...option.RequestOption) func(context.Context) (*requestconfig.RequestConfig, error) {
	requestOptions := slices.Concat([]option.RequestOption{requestconfig.WithBearerAuthSecurity()}, r.Options, opts)
	return func(ctx context.Context) (*requestconfig.RequestConfig, error) {
		if sessionID == "" {
			return nil, errors.New("missing required session_id parameter")
		}
		return requestconfig.NewRequestConfig(ctx, http.MethodGet, requestconfig.FormatPath("live/sessions/%s/fork", sessionID), nil, nil, requestOptions...)
	}
}

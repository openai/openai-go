// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/openai/openai-go/v3/internal/apijson"
	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/param"
)

// BetaChatKitSessionService contains methods and other services that help with
// interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaChatKitSessionService] method instead.
type BetaChatKitSessionService struct {
	Options []option.RequestOption
}

// NewBetaChatKitSessionService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewBetaChatKitSessionService(opts ...option.RequestOption) (r BetaChatKitSessionService) {
	r = BetaChatKitSessionService{}
	r.Options = opts
	return
}

// Create a ChatKit session.
func (r *BetaChatKitSessionService) New(ctx context.Context, body BetaChatKitSessionNewParams, opts ...option.RequestOption) (res *ChatSession, err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "chatkit_beta=v1")}, opts...)
	path := "chatkit/sessions"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Cancel an active ChatKit session and return its most recent metadata.
//
// Cancelling prevents new requests from using the issued client secret.
func (r *BetaChatKitSessionService) Cancel(ctx context.Context, sessionID string, opts ...option.RequestOption) (res *ChatSession, err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "chatkit_beta=v1")}, opts...)
	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return
	}
	path := fmt.Sprintf("chatkit/sessions/%s/cancel", sessionID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return
}

type BetaChatKitSessionNewParams struct {
	// A free-form string that identifies your end user; ensures this Session can
	// access other objects that have the same `user` scope.
	User string `json:"user,required"`
	// Workflow that powers the session.
	Workflow ChatSessionWorkflowParam `json:"workflow,omitzero,required"`
	// Optional overrides for ChatKit runtime configuration features
	ChatKitConfiguration ChatSessionChatKitConfigurationParam `json:"chatkit_configuration,omitzero"`
	// Optional override for session expiration timing in seconds from creation.
	// Defaults to 10 minutes.
	ExpiresAfter ChatSessionExpiresAfterParam `json:"expires_after,omitzero"`
	// Optional override for per-minute request limits. When omitted, defaults to 10.
	RateLimits ChatSessionRateLimitsParam `json:"rate_limits,omitzero"`
	paramObj
}

func (r BetaChatKitSessionNewParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaChatKitSessionNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaChatKitSessionNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

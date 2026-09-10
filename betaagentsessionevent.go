// File generated from our OpenAPI spec by Castiron. See CONTRIBUTING.md for details.

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
	"github.com/openai/openai-go/v3/packages/ssestream"
)

// BetaAgentSessionEventService contains methods and other services that help with
// interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaAgentSessionEventService] method instead.
type BetaAgentSessionEventService struct {
	Options []option.RequestOption
}

// NewBetaAgentSessionEventService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewBetaAgentSessionEventService(opts ...option.RequestOption) (r BetaAgentSessionEventService) {
	r = BetaAgentSessionEventService{}
	r.Options = requestconfig.InheritedOptions(opts...)
	return
}

// Submits message, cancellation, or tool-result events to a managed agent session.
// See
// [session events](https://developers.openai.com/api/docs/guides/agents-api/sessions/events).
func (r *BetaAgentSessionEventService) New(ctx context.Context, sessionID string, params BetaAgentSessionEventNewParams, opts ...option.RequestOption) (err error) {
	if !param.IsOmitted(params.IdempotencyKey) {
		opts = append(opts, option.WithHeader("Idempotency-Key", fmt.Sprintf("%v", params.IdempotencyKey.Value)))
	}
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*"), option.WithHeader("OpenAI-Beta", "agents=v1")}, opts...)
	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return err
	}
	path := requestconfig.FormatPath("agents/sessions/%s/events", sessionID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, nil, opts...)
	return err
}

// Streams live events for an agent session. See
// [session events](https://developers.openai.com/api/docs/guides/agents-api/sessions/events).
func (r *BetaAgentSessionEventService) StreamStreaming(ctx context.Context, sessionID string, opts ...option.RequestOption) (stream *ssestream.Stream[AgentSessionEventUnion]) {
	var (
		raw *http.Response
		err error
	)
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "agents=v1"), option.WithHeader("Accept", "text/event-stream")}, opts...)
	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return ssestream.NewStream[AgentSessionEventUnion](nil, err)
	}
	path := requestconfig.FormatPath("agents/sessions/%s/events", sessionID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &raw, opts...)
	return ssestream.NewStream[AgentSessionEventUnion](ssestream.NewDecoder(raw), err)
}

type BetaAgentSessionEventNewParams struct {
	// The input events to submit to the session.
	Events         []AgentSessionInputParamUnion `json:"events,omitzero" api:"required"`
	IdempotencyKey param.Opt[string]             `header:"Idempotency-Key,omitzero" json:"-"`
	paramObj
}

func (r BetaAgentSessionEventNewParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaAgentSessionEventNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaAgentSessionEventNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

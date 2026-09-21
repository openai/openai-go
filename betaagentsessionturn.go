// File generated from our OpenAPI spec by Castiron. See CONTRIBUTING.md for details.

package openai

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"slices"

	"github.com/openai/openai-go/v3/internal/apijson"
	"github.com/openai/openai-go/v3/internal/apiquery"
	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/pagination"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/packages/respjson"
)

// BetaAgentSessionTurnService contains methods and other services that help with
// interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaAgentSessionTurnService] method instead.
type BetaAgentSessionTurnService struct {
	Options []option.RequestOption
}

// NewBetaAgentSessionTurnService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewBetaAgentSessionTurnService(opts ...option.RequestOption) (r BetaAgentSessionTurnService) {
	r = BetaAgentSessionTurnService{}
	r.Options = requestconfig.InheritedOptions(opts...)
	return
}

// Retrieves a turn's current status, timestamps, usage, and error. Returns 404 if
// the turn does not belong to the session. See
// [session turns](https://developers.openai.com/api/docs/guides/agents-api/sessions/manage#inspect-session-turns).
func (r *BetaAgentSessionTurnService) Get(ctx context.Context, sessionID string, turnID string, opts ...option.RequestOption) (res *Turn, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "agents=v1")}, opts...)
	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	if turnID == "" {
		err = errors.New("missing required turn_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("agents/sessions/%s/turns/%s", sessionID, turnID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Lists turns by creation time and turn ID. The after cursor is exclusive in the
// selected order. See
// [session turns](https://developers.openai.com/api/docs/guides/agents-api/sessions/manage#inspect-session-turns).
func (r *BetaAgentSessionTurnService) List(ctx context.Context, sessionID string, query BetaAgentSessionTurnListParams, opts ...option.RequestOption) (res *pagination.CursorPage[Turn], err error) {
	var raw *http.Response
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "agents=v1"), option.WithResponseInto(&raw)}, opts...)
	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("agents/sessions/%s/turns", sessionID)
	cfg, err := requestconfig.NewRequestConfig(ctx, http.MethodGet, path, query, &res, opts...)
	if err != nil {
		return nil, err
	}
	err = cfg.Execute()
	if err != nil {
		return nil, err
	}
	res.SetPageConfig(cfg, raw)
	return res, nil
}

// Lists turns by creation time and turn ID. The after cursor is exclusive in the
// selected order. See
// [session turns](https://developers.openai.com/api/docs/guides/agents-api/sessions/manage#inspect-session-turns).
func (r *BetaAgentSessionTurnService) ListAutoPaging(ctx context.Context, sessionID string, query BetaAgentSessionTurnListParams, opts ...option.RequestOption) *pagination.CursorPageAutoPager[Turn] {
	return pagination.NewCursorPageAutoPager(r.List(ctx, sessionID, query, opts...))
}

// The canonical public representation of a session turn.
type Turn struct {
	// The ID of the turn.
	ID string `json:"id" api:"required"`
	// The ID of the agent that ran the turn.
	AgentID string `json:"agent_id" api:"required"`
	// The Unix timestamp, in seconds, when the turn reached a terminal state.
	CompletedAt int64 `json:"completed_at" api:"required"`
	// The Unix timestamp, in seconds, used to order the turn by creation time.
	// Subagent turns use their start time, falling back to completion time or the
	// subagent opening time when the preceding timestamps are unavailable.
	CreatedAt int64 `json:"created_at" api:"required"`
	// A customer-safe error describing why a session request failed.
	Error SessionTurnError `json:"error" api:"required"`
	// The object type. Always `agent.session.turn`.
	//
	// Any of "agent.session.turn".
	Object TurnObject `json:"object" api:"required"`
	// The ID of the session that owns the turn.
	SessionID string `json:"session_id" api:"required"`
	// The Unix timestamp, in seconds, when the turn started.
	StartedAt int64 `json:"started_at" api:"required"`
	// The current status of the turn.
	//
	// Any of "queued", "in_progress", "waiting", "completed", "failed", "cancelled".
	Status TurnStatus `json:"status" api:"required"`
	// The ID of the subagent that ran the turn, if applicable.
	SubagentID string `json:"subagent_id" api:"required"`
	// Recorded token usage for a session or turn. Usage is best effort and may change.
	Usage TokenUsage `json:"usage" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		AgentID     respjson.Field
		CompletedAt respjson.Field
		CreatedAt   respjson.Field
		Error       respjson.Field
		Object      respjson.Field
		SessionID   respjson.Field
		StartedAt   respjson.Field
		Status      respjson.Field
		SubagentID  respjson.Field
		Usage       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r Turn) RawJSON() string { return r.JSON.raw }
func (r *Turn) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The object type. Always `agent.session.turn`.
type TurnObject string

const (
	TurnObjectAgentSessionTurn TurnObject = "agent.session.turn"
)

// The current status of the turn.
type TurnStatus string

const (
	TurnStatusQueued     TurnStatus = "queued"
	TurnStatusInProgress TurnStatus = "in_progress"
	TurnStatusWaiting    TurnStatus = "waiting"
	TurnStatusCompleted  TurnStatus = "completed"
	TurnStatusFailed     TurnStatus = "failed"
	TurnStatusCancelled  TurnStatus = "cancelled"
)

type BetaAgentSessionTurnListParams struct {
	// Return resources after this resource ID in the selected order.
	After param.Opt[string] `query:"after,omitzero" json:"-"`
	// The maximum number of resources to return, between 1 and 100. Defaults to 20.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// The order in which resources are returned. Defaults to `desc`.
	//
	// Any of "asc", "desc".
	Order BetaAgentSessionTurnListParamsOrder `query:"order,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaAgentSessionTurnListParams]'s query parameters as
// `url.Values`.
func (r BetaAgentSessionTurnListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// The order in which resources are returned. Defaults to `desc`.
type BetaAgentSessionTurnListParamsOrder string

const (
	BetaAgentSessionTurnListParamsOrderAsc  BetaAgentSessionTurnListParamsOrder = "asc"
	BetaAgentSessionTurnListParamsOrderDesc BetaAgentSessionTurnListParamsOrder = "desc"
)

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
	"github.com/openai/openai-go/v3/packages/ssestream"
)

// BetaAgentSessionService contains methods and other services that help with
// interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaAgentSessionService] method instead.
type BetaAgentSessionService struct {
	Options   []option.RequestOption
	Subagents BetaAgentSessionSubagentService
	Artifacts BetaAgentSessionArtifactService
	Items     BetaAgentSessionItemService
	Events    BetaAgentSessionEventService
	Turns     BetaAgentSessionTurnService
}

// NewBetaAgentSessionService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewBetaAgentSessionService(opts ...option.RequestOption) (r BetaAgentSessionService) {
	r = BetaAgentSessionService{}
	r.Options = requestconfig.InheritedOptions(opts...)
	r.Subagents = NewBetaAgentSessionSubagentService(opts...)
	r.Artifacts = NewBetaAgentSessionArtifactService(opts...)
	r.Items = NewBetaAgentSessionItemService(opts...)
	r.Events = NewBetaAgentSessionEventService(opts...)
	r.Turns = NewBetaAgentSessionTurnService(opts...)
	return
}

// Creates a managed agent session, optionally submits initial input, and returns
// the session or streams its events when stream is true. See
// [running sessions](https://developers.openai.com/api/docs/guides/agents-api/sessions).
func (r *BetaAgentSessionService) New(ctx context.Context, body BetaAgentSessionNewParams, opts ...option.RequestOption) (res *AgentSession, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "agents=v1")}, opts...)
	path := "agents/sessions"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Creates a managed agent session, optionally submits initial input, and returns
// the session or streams its events when stream is true. See
// [running sessions](https://developers.openai.com/api/docs/guides/agents-api/sessions).
func (r *BetaAgentSessionService) NewStreaming(ctx context.Context, body BetaAgentSessionNewParams, opts ...option.RequestOption) (stream *ssestream.Stream[AgentSessionEventUnion]) {
	var (
		raw *http.Response
		err error
	)
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "agents=v1")}, opts...)
	opts = append(opts, option.WithJSONSet("stream", true))
	path := "agents/sessions"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &raw, opts...)
	return ssestream.NewStream[AgentSessionEventUnion](ssestream.NewDecoder(raw), err)
}

// Retrieves the current state of a managed agent session. See
// [managing sessions](https://developers.openai.com/api/docs/guides/agents-api/sessions/manage).
func (r *BetaAgentSessionService) Get(ctx context.Context, sessionID string, opts ...option.RequestOption) (res *AgentSession, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "agents=v1")}, opts...)
	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("agents/sessions/%s", sessionID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Updates session metadata. Omitted fields are unchanged. See
// [managing sessions](https://developers.openai.com/api/docs/guides/agents-api/sessions/manage).
func (r *BetaAgentSessionService) Update(ctx context.Context, sessionID string, body BetaAgentSessionUpdateParams, opts ...option.RequestOption) (res *AgentSession, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "agents=v1")}, opts...)
	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("agents/sessions/%s", sessionID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Lists managed agent sessions using ID-based pagination and the requested sort
// order. See
// [managing sessions](https://developers.openai.com/api/docs/guides/agents-api/sessions/manage).
func (r *BetaAgentSessionService) List(ctx context.Context, query BetaAgentSessionListParams, opts ...option.RequestOption) (res *pagination.CursorPage[AgentSession], err error) {
	var raw *http.Response
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "agents=v1"), option.WithResponseInto(&raw)}, opts...)
	path := "agents/sessions"
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

// Lists managed agent sessions using ID-based pagination and the requested sort
// order. See
// [managing sessions](https://developers.openai.com/api/docs/guides/agents-api/sessions/manage).
func (r *BetaAgentSessionService) ListAutoPaging(ctx context.Context, query BetaAgentSessionListParams, opts ...option.RequestOption) *pagination.CursorPageAutoPager[AgentSession] {
	return pagination.NewCursorPageAutoPager(r.List(ctx, query, opts...))
}

// Removes a managed agent session from the public API and returns a deletion
// confirmation. Physical cleanup may continue asynchronously. See
// [managing sessions](https://developers.openai.com/api/docs/guides/agents-api/sessions/manage).
func (r *BetaAgentSessionService) Delete(ctx context.Context, sessionID string, opts ...option.RequestOption) (res *AgentSessionDeleted, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "agents=v1")}, opts...)
	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("agents/sessions/%s", sessionID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

type BetaAgentSessionNewParams struct {
	// An inline execution environment or a reference to an environment template.
	Environment EnvironmentParamUnion `json:"environment,omitzero" api:"required"`
	// The ID of a saved reusable agent. Omit `agent` to use its configuration
	// unchanged.
	AgentID param.Opt[string] `json:"agent_id,omitzero"`
	// Initial input submitted when creating a session.
	Input BetaAgentSessionNewParamsInputUnion `json:"input,omitzero"`
	// Up to 16 string key-value pairs, with keys up to 64 and values up to 512
	// characters. Omission or null defaults to an empty map.
	Metadata map[string]string `json:"metadata,omitzero"`
	// The IDs of vaults made available to the session.
	VaultIDs []string `json:"vault_ids,omitzero"`
	// Agent configuration. With `agent_id`, supplied fields override the saved agent
	// for this session. Without `agent_id`, `model` is required.
	Agent BetaAgentSessionNewParamsAgent `json:"agent,omitzero"`
	paramObj
}

func (r BetaAgentSessionNewParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaAgentSessionNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaAgentSessionNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Agent configuration. With `agent_id`, supplied fields override the saved agent
// for this session. Without `agent_id`, `model` is required.
type BetaAgentSessionNewParamsAgent struct {
	// Additional instructions appended to the agent's default base instructions. Omit
	// to leave unchanged.
	Instructions param.Opt[string] `json:"instructions,omitzero"`
	// The model to use for the agent. The requested model name is preserved.
	Model param.Opt[string] `json:"model,omitzero"`
	// The service tier used for model requests.
	//
	// Any of "auto", "default", "flex", "priority", "fast".
	ServiceTier string `json:"service_tier,omitzero"`
	// Tools available to the agent. Omit to inherit, or pass null to clear them.
	Tools []AgentToolParamUnion `json:"tools,omitzero"`
	// Explicit configuration for creating and coordinating subagents.
	MultiAgent MultiAgentConfigParam `json:"multi_agent,omitzero"`
	// Reasoning configuration for the agent.
	Reasoning AgentReasoningParam `json:"reasoning,omitzero"`
	// Configuration for text generated by the agent.
	Text AgentTextParam `json:"text,omitzero"`
	paramObj
}

func (r BetaAgentSessionNewParamsAgent) MarshalJSON() (data []byte, err error) {
	type shadow BetaAgentSessionNewParamsAgent
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaAgentSessionNewParamsAgent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[BetaAgentSessionNewParamsAgent](
		"service_tier", "auto", "default", "flex", "priority", "fast",
	)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaAgentSessionNewParamsInputUnion struct {
	OfString               param.Opt[string]               `json:",omitzero,inline"`
	OfArrayOfInputMessages []AgentSessionInputMessageParam `json:",omitzero,inline"`
	paramUnion
}

func (u BetaAgentSessionNewParamsInputUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfString, u.OfArrayOfInputMessages)
}
func (u *BetaAgentSessionNewParamsInputUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

type BetaAgentSessionUpdateParams struct {
	// Replaces all metadata. Omit to leave unchanged, or pass null or {} to clear it.
	// Up to 16 string key-value pairs, with keys up to 64 and values up to 512
	// characters.
	Metadata map[string]string `json:"metadata,omitzero"`
	paramObj
}

func (r BetaAgentSessionUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaAgentSessionUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaAgentSessionUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaAgentSessionListParams struct {
	// Return resources after this resource ID in the selected order.
	After param.Opt[string] `query:"after,omitzero" json:"-"`
	// Only return sessions whose root agent has this ID. Omit to return sessions for
	// all agents.
	AgentID param.Opt[string] `query:"agent_id,omitzero" json:"-"`
	// The maximum number of resources to return.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Sort order by the `created_at` timestamp. Use `asc` for ascending order or
	// `desc` for descending order. Defaults to `desc`.
	//
	// Any of "asc", "desc".
	Order BetaAgentSessionListParamsOrder `query:"order,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaAgentSessionListParams]'s query parameters as
// `url.Values`.
func (r BetaAgentSessionListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Sort order by the `created_at` timestamp. Use `asc` for ascending order or
// `desc` for descending order. Defaults to `desc`.
type BetaAgentSessionListParamsOrder string

const (
	BetaAgentSessionListParamsOrderAsc  BetaAgentSessionListParamsOrder = "asc"
	BetaAgentSessionListParamsOrderDesc BetaAgentSessionListParamsOrder = "desc"
)

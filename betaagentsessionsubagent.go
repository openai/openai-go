// File generated from our OpenAPI spec by Castiron. See CONTRIBUTING.md for details.

package openai

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"slices"

	"github.com/openai/openai-go/v3/internal/apiquery"
	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/pagination"
	"github.com/openai/openai-go/v3/packages/param"
)

// BetaAgentSessionSubagentService contains methods and other services that help
// with interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaAgentSessionSubagentService] method instead.
type BetaAgentSessionSubagentService struct {
	Options []option.RequestOption
	Items   BetaAgentSessionSubagentItemService
	Turns   BetaAgentSessionSubagentTurnService
}

// NewBetaAgentSessionSubagentService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewBetaAgentSessionSubagentService(opts ...option.RequestOption) (r BetaAgentSessionSubagentService) {
	r = BetaAgentSessionSubagentService{}
	r.Options = requestconfig.InheritedOptions(opts...)
	r.Items = NewBetaAgentSessionSubagentItemService(opts...)
	r.Turns = NewBetaAgentSessionSubagentTurnService(opts...)
	return
}

// Retrieves a subagent belonging to this session. See
// [subagent workflows](https://developers.openai.com/api/docs/guides/agents-api/multi-agent).
func (r *BetaAgentSessionSubagentService) Get(ctx context.Context, sessionID string, subagentID string, opts ...option.RequestOption) (res *Subagent, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "agents=v1")}, opts...)
	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	if subagentID == "" {
		err = errors.New("missing required subagent_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("agents/sessions/%s/subagents/%s", sessionID, subagentID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Lists subagents in a session, including nested and closed subagents. See
// [subagent workflows](https://developers.openai.com/api/docs/guides/agents-api/multi-agent).
func (r *BetaAgentSessionSubagentService) List(ctx context.Context, sessionID string, query BetaAgentSessionSubagentListParams, opts ...option.RequestOption) (res *pagination.CursorPage[Subagent], err error) {
	var raw *http.Response
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "agents=v1"), option.WithResponseInto(&raw)}, opts...)
	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("agents/sessions/%s/subagents", sessionID)
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

// Lists subagents in a session, including nested and closed subagents. See
// [subagent workflows](https://developers.openai.com/api/docs/guides/agents-api/multi-agent).
func (r *BetaAgentSessionSubagentService) ListAutoPaging(ctx context.Context, sessionID string, query BetaAgentSessionSubagentListParams, opts ...option.RequestOption) *pagination.CursorPageAutoPager[Subagent] {
	return pagination.NewCursorPageAutoPager(r.List(ctx, sessionID, query, opts...))
}

type BetaAgentSessionSubagentListParams struct {
	// Return resources after this resource ID in the selected order.
	After param.Opt[string] `query:"after,omitzero" json:"-"`
	// The maximum number of resources to return, between 1 and 100. Defaults to 20.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// The order in which resources are returned. Defaults to `desc`.
	//
	// Any of "asc", "desc".
	Order BetaAgentSessionSubagentListParamsOrder `query:"order,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaAgentSessionSubagentListParams]'s query parameters as
// `url.Values`.
func (r BetaAgentSessionSubagentListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// The order in which resources are returned. Defaults to `desc`.
type BetaAgentSessionSubagentListParamsOrder string

const (
	BetaAgentSessionSubagentListParamsOrderAsc  BetaAgentSessionSubagentListParamsOrder = "asc"
	BetaAgentSessionSubagentListParamsOrderDesc BetaAgentSessionSubagentListParamsOrder = "desc"
)

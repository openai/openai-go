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

// BetaAgentSessionSubagentTurnService contains methods and other services that
// help with interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaAgentSessionSubagentTurnService] method instead.
type BetaAgentSessionSubagentTurnService struct {
	Options []option.RequestOption
	Items   BetaAgentSessionSubagentTurnItemService
}

// NewBetaAgentSessionSubagentTurnService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewBetaAgentSessionSubagentTurnService(opts ...option.RequestOption) (r BetaAgentSessionSubagentTurnService) {
	r = BetaAgentSessionSubagentTurnService{}
	r.Options = requestconfig.InheritedOptions(opts...)
	r.Items = NewBetaAgentSessionSubagentTurnItemService(opts...)
	return
}

// Retrieves a turn belonging to this subagent. See
// [subagent workflows](https://developers.openai.com/api/docs/guides/agents-api/multi-agent).
func (r *BetaAgentSessionSubagentTurnService) Get(ctx context.Context, sessionID string, subagentID string, turnID string, opts ...option.RequestOption) (res *Turn, err error) {
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
	if turnID == "" {
		err = errors.New("missing required turn_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("agents/sessions/%s/subagents/%s/turns/%s", sessionID, subagentID, turnID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Lists all turns of this subagent, including turns after a resume. See
// [subagent workflows](https://developers.openai.com/api/docs/guides/agents-api/multi-agent).
func (r *BetaAgentSessionSubagentTurnService) List(ctx context.Context, sessionID string, subagentID string, query BetaAgentSessionSubagentTurnListParams, opts ...option.RequestOption) (res *pagination.CursorPage[Turn], err error) {
	var raw *http.Response
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "agents=v1"), option.WithResponseInto(&raw)}, opts...)
	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	if subagentID == "" {
		err = errors.New("missing required subagent_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("agents/sessions/%s/subagents/%s/turns", sessionID, subagentID)
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

// Lists all turns of this subagent, including turns after a resume. See
// [subagent workflows](https://developers.openai.com/api/docs/guides/agents-api/multi-agent).
func (r *BetaAgentSessionSubagentTurnService) ListAutoPaging(ctx context.Context, sessionID string, subagentID string, query BetaAgentSessionSubagentTurnListParams, opts ...option.RequestOption) *pagination.CursorPageAutoPager[Turn] {
	return pagination.NewCursorPageAutoPager(r.List(ctx, sessionID, subagentID, query, opts...))
}

type BetaAgentSessionSubagentTurnListParams struct {
	// Return resources after this resource ID in the selected order.
	After param.Opt[string] `query:"after,omitzero" json:"-"`
	// The maximum number of resources to return, between 1 and 100. Defaults to 20.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// The order in which resources are returned. Defaults to `desc`.
	//
	// Any of "asc", "desc".
	Order BetaAgentSessionSubagentTurnListParamsOrder `query:"order,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaAgentSessionSubagentTurnListParams]'s query parameters
// as `url.Values`.
func (r BetaAgentSessionSubagentTurnListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// The order in which resources are returned. Defaults to `desc`.
type BetaAgentSessionSubagentTurnListParamsOrder string

const (
	BetaAgentSessionSubagentTurnListParamsOrderAsc  BetaAgentSessionSubagentTurnListParamsOrder = "asc"
	BetaAgentSessionSubagentTurnListParamsOrderDesc BetaAgentSessionSubagentTurnListParamsOrder = "desc"
)

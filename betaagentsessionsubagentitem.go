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

// BetaAgentSessionSubagentItemService contains methods and other services that
// help with interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaAgentSessionSubagentItemService] method instead.
type BetaAgentSessionSubagentItemService struct {
	Options []option.RequestOption
}

// NewBetaAgentSessionSubagentItemService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewBetaAgentSessionSubagentItemService(opts ...option.RequestOption) (r BetaAgentSessionSubagentItemService) {
	r = BetaAgentSessionSubagentItemService{}
	r.Options = requestconfig.InheritedOptions(opts...)
	return
}

// Lists this subagent's own items across all of its turns. See
// [subagent workflows](https://developers.openai.com/api/docs/guides/agents-api/multi-agent).
func (r *BetaAgentSessionSubagentItemService) List(ctx context.Context, sessionID string, subagentID string, query BetaAgentSessionSubagentItemListParams, opts ...option.RequestOption) (res *pagination.CursorPage[AgentSessionItemUnion], err error) {
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
	path := requestconfig.FormatPath("agents/sessions/%s/subagents/%s/items", sessionID, subagentID)
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

// Lists this subagent's own items across all of its turns. See
// [subagent workflows](https://developers.openai.com/api/docs/guides/agents-api/multi-agent).
func (r *BetaAgentSessionSubagentItemService) ListAutoPaging(ctx context.Context, sessionID string, subagentID string, query BetaAgentSessionSubagentItemListParams, opts ...option.RequestOption) *pagination.CursorPageAutoPager[AgentSessionItemUnion] {
	return pagination.NewCursorPageAutoPager(r.List(ctx, sessionID, subagentID, query, opts...))
}

type BetaAgentSessionSubagentItemListParams struct {
	// Return resources after this resource ID in the selected order.
	After param.Opt[string] `query:"after,omitzero" json:"-"`
	// The maximum number of resources to return, between 1 and 100. Defaults to 20.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// The order in which resources are returned. Defaults to `desc`.
	//
	// Any of "asc", "desc".
	Order BetaAgentSessionSubagentItemListParamsOrder `query:"order,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaAgentSessionSubagentItemListParams]'s query parameters
// as `url.Values`.
func (r BetaAgentSessionSubagentItemListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// The order in which resources are returned. Defaults to `desc`.
type BetaAgentSessionSubagentItemListParamsOrder string

const (
	BetaAgentSessionSubagentItemListParamsOrderAsc  BetaAgentSessionSubagentItemListParamsOrder = "asc"
	BetaAgentSessionSubagentItemListParamsOrderDesc BetaAgentSessionSubagentItemListParamsOrder = "desc"
)

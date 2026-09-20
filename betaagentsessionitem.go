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

// BetaAgentSessionItemService contains methods and other services that help with
// interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaAgentSessionItemService] method instead.
type BetaAgentSessionItemService struct {
	Options []option.RequestOption
}

// NewBetaAgentSessionItemService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewBetaAgentSessionItemService(opts ...option.RequestOption) (r BetaAgentSessionItemService) {
	r = BetaAgentSessionItemService{}
	r.Options = requestconfig.InheritedOptions(opts...)
	return
}

// Lists items produced by the session's root agent, including its interactions
// with subagents. Each subagent has its own item history. See
// [inspecting agent output](https://developers.openai.com/api/docs/guides/agents-api/observability).
func (r *BetaAgentSessionItemService) List(ctx context.Context, sessionID string, query BetaAgentSessionItemListParams, opts ...option.RequestOption) (res *pagination.CursorPage[AgentSessionItemUnion], err error) {
	var raw *http.Response
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "agents=v1"), option.WithResponseInto(&raw)}, opts...)
	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("agents/sessions/%s/items", sessionID)
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

// Lists items produced by the session's root agent, including its interactions
// with subagents. Each subagent has its own item history. See
// [inspecting agent output](https://developers.openai.com/api/docs/guides/agents-api/observability).
func (r *BetaAgentSessionItemService) ListAutoPaging(ctx context.Context, sessionID string, query BetaAgentSessionItemListParams, opts ...option.RequestOption) *pagination.CursorPageAutoPager[AgentSessionItemUnion] {
	return pagination.NewCursorPageAutoPager(r.List(ctx, sessionID, query, opts...))
}

type BetaAgentSessionItemListParams struct {
	// Return resources after this resource ID in the selected order.
	After param.Opt[string] `query:"after,omitzero" json:"-"`
	// The maximum number of resources to return, between 1 and 100. Defaults to 20.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// The order in which resources are returned. Defaults to `desc`.
	//
	// Any of "asc", "desc".
	Order BetaAgentSessionItemListParamsOrder `query:"order,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaAgentSessionItemListParams]'s query parameters as
// `url.Values`.
func (r BetaAgentSessionItemListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// The order in which resources are returned. Defaults to `desc`.
type BetaAgentSessionItemListParamsOrder string

const (
	BetaAgentSessionItemListParamsOrderAsc  BetaAgentSessionItemListParamsOrder = "asc"
	BetaAgentSessionItemListParamsOrderDesc BetaAgentSessionItemListParamsOrder = "desc"
)

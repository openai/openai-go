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

// BetaAgentSessionTurnItemService contains methods and other services that help
// with interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaAgentSessionTurnItemService] method instead.
type BetaAgentSessionTurnItemService struct {
	Options []option.RequestOption
}

// NewBetaAgentSessionTurnItemService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewBetaAgentSessionTurnItemService(opts ...option.RequestOption) (r BetaAgentSessionTurnItemService) {
	r = BetaAgentSessionTurnItemService{}
	r.Options = requestconfig.InheritedOptions(opts...)
	return
}

// Lists items belonging to one root-agent turn, including its interactions with
// subagents. See
// [inspecting agent output](https://developers.openai.com/api/docs/guides/agents-api/observability).
func (r *BetaAgentSessionTurnItemService) List(ctx context.Context, sessionID string, turnID string, query BetaAgentSessionTurnItemListParams, opts ...option.RequestOption) (res *pagination.ConversationCursorPage[AgentSessionItemUnion], err error) {
	var raw *http.Response
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "agents=v1"), option.WithResponseInto(&raw)}, opts...)
	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	if turnID == "" {
		err = errors.New("missing required turn_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("agents/sessions/%s/turns/%s/items", sessionID, turnID)
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

// Lists items belonging to one root-agent turn, including its interactions with
// subagents. See
// [inspecting agent output](https://developers.openai.com/api/docs/guides/agents-api/observability).
func (r *BetaAgentSessionTurnItemService) ListAutoPaging(ctx context.Context, sessionID string, turnID string, query BetaAgentSessionTurnItemListParams, opts ...option.RequestOption) *pagination.ConversationCursorPageAutoPager[AgentSessionItemUnion] {
	return pagination.NewConversationCursorPageAutoPager(r.List(ctx, sessionID, turnID, query, opts...))
}

type BetaAgentSessionTurnItemListParams struct {
	// Return items after this cursor in the selected order. Pass the previous
	// response's last_id, which can differ from the last item's ID.
	After param.Opt[string] `query:"after,omitzero" json:"-"`
	// The maximum number of resources to return, between 1 and 100. Defaults to 20.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// The order in which resources are returned. Defaults to `desc`.
	//
	// Any of "asc", "desc".
	Order BetaAgentSessionTurnItemListParamsOrder `query:"order,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaAgentSessionTurnItemListParams]'s query parameters as
// `url.Values`.
func (r BetaAgentSessionTurnItemListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// The order in which resources are returned. Defaults to `desc`.
type BetaAgentSessionTurnItemListParamsOrder string

const (
	BetaAgentSessionTurnItemListParamsOrderAsc  BetaAgentSessionTurnItemListParamsOrder = "asc"
	BetaAgentSessionTurnItemListParamsOrderDesc BetaAgentSessionTurnItemListParamsOrder = "desc"
)

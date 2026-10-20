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
	"github.com/openai/openai-go/v3/shared/constant"
)

// BetaAgentSessionTraceService contains methods and other services that help with
// interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaAgentSessionTraceService] method instead.
type BetaAgentSessionTraceService struct {
	Options []option.RequestOption
}

// NewBetaAgentSessionTraceService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewBetaAgentSessionTraceService(opts ...option.RequestOption) (r BetaAgentSessionTraceService) {
	r = BetaAgentSessionTraceService{}
	r.Options = requestconfig.InheritedOptions(opts...)
	return
}

// Lists published root-turn traces as OTLP JSON, ordered by turn creation time and
// ID. Unpublished traces are skipped. Each page returns data available when read;
// it does not wait for late traces. Trace reads and the JSON response are limited
// to 16 MiB per request. If the limit is exceeded, request fewer traces.
func (r *BetaAgentSessionTraceService) List(ctx context.Context, sessionID string, query BetaAgentSessionTraceListParams, opts ...option.RequestOption) (res *pagination.CursorPage[SessionTrace], err error) {
	var raw *http.Response
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "agents=v1"), option.WithResponseInto(&raw)}, opts...)
	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("agents/sessions/%s/traces", sessionID)
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

// Lists published root-turn traces as OTLP JSON, ordered by turn creation time and
// ID. Unpublished traces are skipped. Each page returns data available when read;
// it does not wait for late traces. Trace reads and the JSON response are limited
// to 16 MiB per request. If the limit is exceeded, request fewer traces.
func (r *BetaAgentSessionTraceService) ListAutoPaging(ctx context.Context, sessionID string, query BetaAgentSessionTraceListParams, opts ...option.RequestOption) *pagination.CursorPageAutoPager[SessionTrace] {
	return pagination.NewCursorPageAutoPager(r.List(ctx, sessionID, query, opts...))
}

type SessionTrace struct {
	// The root turn ID. Use this ID as the pagination anchor.
	ID string `json:"id" api:"required"`
	// The Unix timestamp in seconds when the root turn was created.
	CreatedAt int64 `json:"created_at" api:"required"`
	// The object type, which is always `agent.session.trace`.
	Object constant.AgentSessionTrace `json:"object" default:"agent.session.trace"`
	// An OTLP JSON ExportTraceServiceRequest containing resourceSpans. Only currently
	// published data is returned; later trace updates are not awaited.
	Otlp map[string]any `json:"otlp" api:"required"`
	// The session that owns this trace.
	SessionID string `json:"session_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		Object      respjson.Field
		Otlp        respjson.Field
		SessionID   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r SessionTrace) RawJSON() string { return r.JSON.raw }
func (r *SessionTrace) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaAgentSessionTraceListParams struct {
	// Return resources after this resource ID in the selected order.
	After param.Opt[string] `query:"after,omitzero" json:"-"`
	// The maximum number of resources to return, between 1 and 100. Defaults to 20.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// The order in which resources are returned. Defaults to `desc`.
	//
	// Any of "asc", "desc".
	Order BetaAgentSessionTraceListParamsOrder `query:"order,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaAgentSessionTraceListParams]'s query parameters as
// `url.Values`.
func (r BetaAgentSessionTraceListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// The order in which resources are returned. Defaults to `desc`.
type BetaAgentSessionTraceListParamsOrder string

const (
	BetaAgentSessionTraceListParamsOrderAsc  BetaAgentSessionTraceListParamsOrder = "asc"
	BetaAgentSessionTraceListParamsOrderDesc BetaAgentSessionTraceListParamsOrder = "desc"
)

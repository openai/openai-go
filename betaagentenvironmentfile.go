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
	shimjson "github.com/openai/openai-go/v3/internal/encoding/json"
	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/pagination"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/packages/respjson"
	"github.com/openai/openai-go/v3/shared/constant"
)

// BetaAgentEnvironmentFileService contains methods and other services that help
// with interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaAgentEnvironmentFileService] method instead.
type BetaAgentEnvironmentFileService struct {
	Options []option.RequestOption
}

// NewBetaAgentEnvironmentFileService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewBetaAgentEnvironmentFileService(opts ...option.RequestOption) (r BetaAgentEnvironmentFileService) {
	r = BetaAgentEnvironmentFileService{}
	r.Options = requestconfig.InheritedOptions(opts...)
	return
}

// Copies inline bytes or a Files API file into a connected execution environment.
// See
// [environment files](https://developers.openai.com/api/docs/guides/agents-api/environments/files).
func (r *BetaAgentEnvironmentFileService) New(ctx context.Context, environmentID string, body BetaAgentEnvironmentFileNewParams, opts ...option.RequestOption) (res *EnvironmentFile, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "agents=v1")}, opts...)
	if environmentID == "" {
		err = errors.New("missing required environment_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("agents/environments/%s/files", environmentID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Lists live files on a connected execution environment with optional directory
// filtering and opaque cursor pagination. See
// [environment files](https://developers.openai.com/api/docs/guides/agents-api/environments/files).
func (r *BetaAgentEnvironmentFileService) List(ctx context.Context, environmentID string, query BetaAgentEnvironmentFileListParams, opts ...option.RequestOption) (res *pagination.TokenPage[EnvironmentFile], err error) {
	var raw *http.Response
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "agents=v1"), option.WithResponseInto(&raw)}, opts...)
	if environmentID == "" {
		err = errors.New("missing required environment_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("agents/environments/%s/files", environmentID)
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

// Lists live files on a connected execution environment with optional directory
// filtering and opaque cursor pagination. See
// [environment files](https://developers.openai.com/api/docs/guides/agents-api/environments/files).
func (r *BetaAgentEnvironmentFileService) ListAutoPaging(ctx context.Context, environmentID string, query BetaAgentEnvironmentFileListParams, opts ...option.RequestOption) *pagination.TokenPageAutoPager[EnvironmentFile] {
	return pagination.NewTokenPageAutoPager(r.List(ctx, environmentID, query, opts...))
}

// A live file in an execution environment.
type EnvironmentFile struct {
	// The ID of the environment containing this file.
	EnvironmentID string `json:"environment_id" api:"required"`
	// The object type. Always `agent.environment.file`.
	Object constant.AgentEnvironmentFile `json:"object" default:"agent.environment.file"`
	// The absolute file path inside the environment's workspace.
	Path string `json:"path" api:"required"`
	// The file size in bytes.
	SizeBytes int64 `json:"size_bytes" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EnvironmentID respjson.Field
		Object        respjson.Field
		Path          respjson.Field
		SizeBytes     respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EnvironmentFile) RawJSON() string { return r.JSON.raw }
func (r *EnvironmentFile) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaAgentEnvironmentFileNewParams struct {
	// A file materialized in an OpenAI-hosted execution environment.
	HostedEnvironmentFileParam HostedEnvironmentFileParamUnion
	paramObj
}

func (r BetaAgentEnvironmentFileNewParams) MarshalJSON() (data []byte, err error) {
	return shimjson.Marshal(r.HostedEnvironmentFileParam)
}
func (r *BetaAgentEnvironmentFileNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaAgentEnvironmentFileListParams struct {
	// The maximum number of files to return, between 1 and 100.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// The opaque token from the previous page. Keep the same path, order, and limit.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Restrict the listing to this absolute workspace directory.
	Path param.Opt[string] `query:"path,omitzero" json:"-"`
	// Sort by case-sensitive path components. Defaults to descending.
	//
	// Any of "asc", "desc".
	Order BetaAgentEnvironmentFileListParamsOrder `query:"order,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaAgentEnvironmentFileListParams]'s query parameters as
// `url.Values`.
func (r BetaAgentEnvironmentFileListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Sort by case-sensitive path components. Defaults to descending.
type BetaAgentEnvironmentFileListParamsOrder string

const (
	BetaAgentEnvironmentFileListParamsOrderAsc  BetaAgentEnvironmentFileListParamsOrder = "asc"
	BetaAgentEnvironmentFileListParamsOrderDesc BetaAgentEnvironmentFileListParamsOrder = "desc"
)

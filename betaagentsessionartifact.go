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

// BetaAgentSessionArtifactService contains methods and other services that help
// with interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaAgentSessionArtifactService] method instead.
type BetaAgentSessionArtifactService struct {
	Options []option.RequestOption
}

// NewBetaAgentSessionArtifactService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewBetaAgentSessionArtifactService(opts ...option.RequestOption) (r BetaAgentSessionArtifactService) {
	r = BetaAgentSessionArtifactService{}
	r.Options = requestconfig.InheritedOptions(opts...)
	return
}

// Retrieves immutable metadata for one durable session artifact. See
// [session artifacts](https://developers.openai.com/api/docs/guides/agents-api/environments/files#openai-hosted-artifacts).
func (r *BetaAgentSessionArtifactService) Get(ctx context.Context, sessionID string, artifactID string, opts ...option.RequestOption) (res *SessionArtifact, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "agents=v1")}, opts...)
	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	if artifactID == "" {
		err = errors.New("missing required artifact_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("agents/sessions/%s/artifacts/%s", sessionID, artifactID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Lists immutable artifacts published by completed hosted session turns. See
// [session artifacts](https://developers.openai.com/api/docs/guides/agents-api/environments/files#openai-hosted-artifacts).
func (r *BetaAgentSessionArtifactService) List(ctx context.Context, sessionID string, query BetaAgentSessionArtifactListParams, opts ...option.RequestOption) (res *pagination.CursorPage[SessionArtifact], err error) {
	var raw *http.Response
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "agents=v1"), option.WithResponseInto(&raw)}, opts...)
	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("agents/sessions/%s/artifacts", sessionID)
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

// Lists immutable artifacts published by completed hosted session turns. See
// [session artifacts](https://developers.openai.com/api/docs/guides/agents-api/environments/files#openai-hosted-artifacts).
func (r *BetaAgentSessionArtifactService) ListAutoPaging(ctx context.Context, sessionID string, query BetaAgentSessionArtifactListParams, opts ...option.RequestOption) *pagination.CursorPageAutoPager[SessionArtifact] {
	return pagination.NewCursorPageAutoPager(r.List(ctx, sessionID, query, opts...))
}

// Deletes an immutable session artifact without deleting its live environment file
// or original Files API object. See
// [session artifacts](https://developers.openai.com/api/docs/guides/agents-api/environments/files#openai-hosted-artifacts).
func (r *BetaAgentSessionArtifactService) Delete(ctx context.Context, sessionID string, artifactID string, opts ...option.RequestOption) (res *SessionArtifactDeleted, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "agents=v1")}, opts...)
	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	if artifactID == "" {
		err = errors.New("missing required artifact_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("agents/sessions/%s/artifacts/%s", sessionID, artifactID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Downloads immutable session artifact bytes after the execution environment
// expires. See
// [session artifacts](https://developers.openai.com/api/docs/guides/agents-api/environments/files#openai-hosted-artifacts).
func (r *BetaAgentSessionArtifactService) Content(ctx context.Context, sessionID string, artifactID string, opts ...option.RequestOption) (res *http.Response, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "agents=v1"), option.WithHeader("Accept", "application/octet-stream")}, opts...)
	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	if artifactID == "" {
		err = errors.New("missing required artifact_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("agents/sessions/%s/artifacts/%s/content", sessionID, artifactID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// An immutable file published by a completed hosted session turn.
type SessionArtifact struct {
	// The immutable artifact ID.
	ID string `json:"id" api:"required"`
	// The Unix timestamp, in seconds, when the artifact was published.
	CreatedAt int64 `json:"created_at" api:"required"`
	// The ID of the environment that produced the artifact.
	EnvironmentID string `json:"environment_id" api:"required"`
	// The object type. Always `agent.session.artifact`.
	Object constant.AgentSessionArtifact `json:"object" default:"agent.session.artifact"`
	// The original absolute file path in the execution environment.
	Path string `json:"path" api:"required"`
	// The ID of the session that owns the artifact.
	SessionID string `json:"session_id" api:"required"`
	// The immutable artifact size in bytes.
	SizeBytes int64 `json:"size_bytes" api:"required"`
	// The ID of the completed turn that published the artifact.
	TurnID string `json:"turn_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID            respjson.Field
		CreatedAt     respjson.Field
		EnvironmentID respjson.Field
		Object        respjson.Field
		Path          respjson.Field
		SessionID     respjson.Field
		SizeBytes     respjson.Field
		TurnID        respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r SessionArtifact) RawJSON() string { return r.JSON.raw }
func (r *SessionArtifact) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Confirmation that an immutable session artifact was deleted.
type SessionArtifactDeleted struct {
	// The ID of the deleted session artifact.
	ID string `json:"id" api:"required"`
	// Whether the session artifact was deleted. Always `true`.
	Deleted bool `json:"deleted" api:"required"`
	// The object type. Always `agent.session.artifact.deleted`.
	Object constant.AgentSessionArtifactDeleted `json:"object" default:"agent.session.artifact.deleted"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Deleted     respjson.Field
		Object      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r SessionArtifactDeleted) RawJSON() string { return r.JSON.raw }
func (r *SessionArtifactDeleted) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaAgentSessionArtifactListParams struct {
	// Return artifacts after this immutable artifact ID.
	After param.Opt[string] `query:"after,omitzero" json:"-"`
	// Restrict the listing to artifacts produced by this environment.
	EnvironmentID param.Opt[string] `query:"environment_id,omitzero" json:"-"`
	// The maximum number of artifacts to return, between 1 and 100.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Sort by creation time and ID. Defaults to descending.
	//
	// Any of "asc", "desc".
	Order BetaAgentSessionArtifactListParamsOrder `query:"order,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaAgentSessionArtifactListParams]'s query parameters as
// `url.Values`.
func (r BetaAgentSessionArtifactListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Sort by creation time and ID. Defaults to descending.
type BetaAgentSessionArtifactListParamsOrder string

const (
	BetaAgentSessionArtifactListParamsOrderAsc  BetaAgentSessionArtifactListParamsOrder = "asc"
	BetaAgentSessionArtifactListParamsOrderDesc BetaAgentSessionArtifactListParamsOrder = "desc"
)

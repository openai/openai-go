// File generated from our OpenAPI spec by Castiron. See CONTRIBUTING.md for details.

package openai

import (
	"context"
	"errors"
	"fmt"
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

// BetaAgentEnvironmentService contains methods and other services that help with
// interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaAgentEnvironmentService] method instead.
type BetaAgentEnvironmentService struct {
	Options   []option.RequestOption
	Files     BetaAgentEnvironmentFileService
	Templates BetaAgentEnvironmentTemplateService
}

// NewBetaAgentEnvironmentService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewBetaAgentEnvironmentService(opts ...option.RequestOption) (r BetaAgentEnvironmentService) {
	r = BetaAgentEnvironmentService{}
	r.Options = requestconfig.InheritedOptions(opts...)
	r.Files = NewBetaAgentEnvironmentFileService(opts...)
	r.Templates = NewBetaAgentEnvironmentTemplateService(opts...)
	return
}

// Creates an OpenAI-hosted environment before creating a session. Requires access
// to the prewarming beta.
func (r *BetaAgentEnvironmentService) New(ctx context.Context, params BetaAgentEnvironmentNewParams, opts ...option.RequestOption) (res *EnvironmentInfo, err error) {
	if !param.IsOmitted(params.IdempotencyKey) {
		opts = append(opts, option.WithHeader("Idempotency-Key", fmt.Sprintf("%v", params.IdempotencyKey.Value)))
	}
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "agents=v1")}, opts...)
	path := "agents/environments"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// Retrieves an execution environment's connection status and safe installed
// metadata. See
// [environment lifecycle](https://developers.openai.com/api/docs/guides/agents-api/environments/lifecycle).
func (r *BetaAgentEnvironmentService) Get(ctx context.Context, environmentID string, opts ...option.RequestOption) (res *EnvironmentInfo, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "agents=v1")}, opts...)
	if environmentID == "" {
		err = errors.New("missing required environment_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("agents/environments/%s", environmentID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Lists OpenAI-hosted environments owned by the authenticated principal. Requires
// access to the prewarming beta.
func (r *BetaAgentEnvironmentService) List(ctx context.Context, query BetaAgentEnvironmentListParams, opts ...option.RequestOption) (res *pagination.CursorPage[EnvironmentInfo], err error) {
	var raw *http.Response
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "agents=v1"), option.WithResponseInto(&raw)}, opts...)
	path := "agents/environments"
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

// Lists OpenAI-hosted environments owned by the authenticated principal. Requires
// access to the prewarming beta.
func (r *BetaAgentEnvironmentService) ListAutoPaging(ctx context.Context, query BetaAgentEnvironmentListParams, opts ...option.RequestOption) *pagination.CursorPageAutoPager[EnvironmentInfo] {
	return pagination.NewCursorPageAutoPager(r.List(ctx, query, opts...))
}

// Safe metadata for a first-class execution environment.
type EnvironmentInfo struct {
	// The ID of the environment.
	ID string `json:"id" api:"required"`
	// Files installed in the environment, without their contents.
	Files []HostedEnvironmentFileUnion `json:"files" api:"required"`
	// The object type. Always `agent.environment`.
	Object constant.AgentEnvironment `json:"object" default:"agent.environment"`
	// Plugins installed in the environment, without their archive contents.
	Plugins []HostedPlugin `json:"plugins" api:"required"`
	// Skills installed in the environment, without their archive contents.
	Skills []HostedSkillUnion `json:"skills" api:"required"`
	// The current environment connection status.
	//
	// Any of "pending", "ready", "connected", "disconnected", "suspended", "expired",
	// "failed".
	Status EnvironmentInfoStatus `json:"status" api:"required"`
	// Whether the environment is hosted by OpenAI or by the application.
	//
	// Any of "openai_hosted", "self_hosted".
	Type EnvironmentInfoType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Files       respjson.Field
		Object      respjson.Field
		Plugins     respjson.Field
		Skills      respjson.Field
		Status      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EnvironmentInfo) RawJSON() string { return r.JSON.raw }
func (r *EnvironmentInfo) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The current environment connection status.
type EnvironmentInfoStatus string

const (
	EnvironmentInfoStatusPending      EnvironmentInfoStatus = "pending"
	EnvironmentInfoStatusReady        EnvironmentInfoStatus = "ready"
	EnvironmentInfoStatusConnected    EnvironmentInfoStatus = "connected"
	EnvironmentInfoStatusDisconnected EnvironmentInfoStatus = "disconnected"
	EnvironmentInfoStatusSuspended    EnvironmentInfoStatus = "suspended"
	EnvironmentInfoStatusExpired      EnvironmentInfoStatus = "expired"
	EnvironmentInfoStatusFailed       EnvironmentInfoStatus = "failed"
)

// Whether the environment is hosted by OpenAI or by the application.
type EnvironmentInfoType string

const (
	EnvironmentInfoTypeOpenAIHosted EnvironmentInfoType = "openai_hosted"
	EnvironmentInfoTypeSelfHosted   EnvironmentInfoType = "self_hosted"
)

type BetaAgentEnvironmentNewParams struct {
	// The required hosting type and its configuration.
	Environment    BetaAgentEnvironmentNewParamsEnvironment `json:"environment" api:"required"`
	IdempotencyKey param.Opt[string]                        `header:"Idempotency-Key,omitzero" json:"-"`
	// The IDs of up to 10 vaults made available to an OpenAI-hosted environment.
	VaultIDs []string `json:"vault_ids,omitzero"`
	paramObj
}

func (r BetaAgentEnvironmentNewParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaAgentEnvironmentNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaAgentEnvironmentNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The required hosting type and its configuration.
//
// The property Type is required.
type BetaAgentEnvironmentNewParamsEnvironment struct {
	// A reusable hosted template applied before inline configuration. Omitted fields
	// inherit the template; network overrides cannot broaden its policy.
	EnvironmentTemplateID param.Opt[string] `json:"environment_template_id,omitzero"`
	// Directories that contain capabilities exposed to the agent. Defaults to an empty
	// list.
	CapabilityDirectories []string `json:"capability_directories,omitzero"`
	// Desktop provisioning. Omission or null inherits the template setting, or
	// defaults to disabled.
	Desktop BetaAgentEnvironmentNewParamsEnvironmentDesktop `json:"desktop,omitzero"`
	// Environment variables made available to the agent.
	Env map[string]string `json:"env,omitzero"`
	// Files available before the agent starts. Defaults to an empty list.
	Files []HostedEnvironmentFileParamUnion `json:"files,omitzero"`
	// Network access policy for the environment. If omitted, the API version
	// determines whether network access is enabled or disabled.
	Network BetaAgentEnvironmentNewParamsEnvironmentNetwork `json:"network,omitzero"`
	// Packages to install in the environment. Defaults to empty package lists.
	Packages BetaAgentEnvironmentNewParamsEnvironmentPackages `json:"packages,omitzero"`
	// Plugins provided as inline ZIP archives. Defaults to an empty list.
	Plugins []HostedPluginParam `json:"plugins,omitzero"`
	// Ordered, confidential setup commands. Command bodies are never returned.
	SetupCommands []SetupCommandParam `json:"setup_commands,omitzero"`
	// Skills referenced by ID or provided as inline ZIP archives. Defaults to an empty
	// list.
	Skills []HostedSkillParamUnion `json:"skills,omitzero"`
	// The type of the object. Always `openai_hosted`.
	//
	// This field can be elided, and will marshal its zero value as "openai_hosted".
	Type constant.OpenAIHosted `json:"type" default:"openai_hosted"`
	paramObj
}

func (r BetaAgentEnvironmentNewParamsEnvironment) MarshalJSON() (data []byte, err error) {
	type shadow BetaAgentEnvironmentNewParamsEnvironment
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaAgentEnvironmentNewParamsEnvironment) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Desktop provisioning. Omission or null inherits the template setting, or
// defaults to disabled.
//
// The property Enabled is required.
type BetaAgentEnvironmentNewParamsEnvironmentDesktop struct {
	// Whether to provision the desktop and its browser proxy.
	Enabled param.Opt[bool] `json:"enabled,omitzero" api:"required"`
	paramObj
}

func (r BetaAgentEnvironmentNewParamsEnvironmentDesktop) MarshalJSON() (data []byte, err error) {
	type shadow BetaAgentEnvironmentNewParamsEnvironmentDesktop
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaAgentEnvironmentNewParamsEnvironmentDesktop) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Network access policy for the environment. If omitted, the API version
// determines whether network access is enabled or disabled.
//
// The property Access is required.
type BetaAgentEnvironmentNewParamsEnvironmentNetwork struct {
	// The environment's network access mode.
	//
	// Any of "enabled", "disabled", "restricted".
	Access string `json:"access,omitzero" api:"required"`
	// Domains the environment may access when network access is restricted.
	AllowedDomains []string `json:"allowed_domains,omitzero"`
	// Domains blocked for both executor and browser when access is restricted. A
	// nonempty list requires `access: restricted` and cannot be combined with nonempty
	// `allowed_domains`. Wildcard domains are not supported.
	BlockedDomains []string `json:"blocked_domains,omitzero"`
	paramObj
}

func (r BetaAgentEnvironmentNewParamsEnvironmentNetwork) MarshalJSON() (data []byte, err error) {
	type shadow BetaAgentEnvironmentNewParamsEnvironmentNetwork
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaAgentEnvironmentNewParamsEnvironmentNetwork) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[BetaAgentEnvironmentNewParamsEnvironmentNetwork](
		"access", "enabled", "disabled", "restricted",
	)
}

// Packages to install in the environment. Defaults to empty package lists.
type BetaAgentEnvironmentNewParamsEnvironmentPackages struct {
	// npm packages to install globally. Defaults to an empty list.
	Npm []string `json:"npm,omitzero"`
	// Python packages to install. Defaults to an empty list.
	Python []string `json:"python,omitzero"`
	// System packages to install. Defaults to an empty list.
	System []string `json:"system,omitzero"`
	paramObj
}

func (r BetaAgentEnvironmentNewParamsEnvironmentPackages) MarshalJSON() (data []byte, err error) {
	type shadow BetaAgentEnvironmentNewParamsEnvironmentPackages
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaAgentEnvironmentNewParamsEnvironmentPackages) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaAgentEnvironmentListParams struct {
	// Return resources after this resource ID in the selected order.
	After param.Opt[string] `query:"after,omitzero" json:"-"`
	// The maximum number of resources to return, between 1 and 100. Defaults to 20.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// The order in which resources are returned. Defaults to `desc`.
	//
	// Any of "asc", "desc".
	Order BetaAgentEnvironmentListParamsOrder `query:"order,omitzero" json:"-"`
	// The hosting type to list. Defaults to `openai_hosted`.
	//
	// Any of "openai_hosted".
	Type BetaAgentEnvironmentListParamsType `query:"type,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaAgentEnvironmentListParams]'s query parameters as
// `url.Values`.
func (r BetaAgentEnvironmentListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// The order in which resources are returned. Defaults to `desc`.
type BetaAgentEnvironmentListParamsOrder string

const (
	BetaAgentEnvironmentListParamsOrderAsc  BetaAgentEnvironmentListParamsOrder = "asc"
	BetaAgentEnvironmentListParamsOrderDesc BetaAgentEnvironmentListParamsOrder = "desc"
)

// The hosting type to list. Defaults to `openai_hosted`.
type BetaAgentEnvironmentListParamsType string

const (
	BetaAgentEnvironmentListParamsTypeOpenAIHosted BetaAgentEnvironmentListParamsType = "openai_hosted"
)

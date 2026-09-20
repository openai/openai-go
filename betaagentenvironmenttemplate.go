// File generated from our OpenAPI spec by Castiron. See CONTRIBUTING.md for details.

package openai

import (
	"context"
	"encoding/json"
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

// BetaAgentEnvironmentTemplateService contains methods and other services that
// help with interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaAgentEnvironmentTemplateService] method instead.
type BetaAgentEnvironmentTemplateService struct {
	Options []option.RequestOption
}

// NewBetaAgentEnvironmentTemplateService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewBetaAgentEnvironmentTemplateService(opts ...option.RequestOption) (r BetaAgentEnvironmentTemplateService) {
	r = BetaAgentEnvironmentTemplateService{}
	r.Options = requestconfig.InheritedOptions(opts...)
	return
}

// Creates reusable environment configuration without returning confidential setup
// commands or environment values. See
// [reusing a hosted setup](https://developers.openai.com/api/docs/guides/agents-api/tools#reuse-a-hosted-plugin-setup).
func (r *BetaAgentEnvironmentTemplateService) New(ctx context.Context, body BetaAgentEnvironmentTemplateNewParams, opts ...option.RequestOption) (res *EnvironmentTemplate, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "agents=v1")}, opts...)
	path := "agents/environments/templates"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Retrieves reusable environment configuration without returning confidential
// values. See
// [reusing a hosted setup](https://developers.openai.com/api/docs/guides/agents-api/tools#reuse-a-hosted-plugin-setup).
func (r *BetaAgentEnvironmentTemplateService) Get(ctx context.Context, environmentTemplateID string, opts ...option.RequestOption) (res *EnvironmentTemplate, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "agents=v1")}, opts...)
	if environmentTemplateID == "" {
		err = errors.New("missing required environment_template_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("agents/environments/templates/%s", environmentTemplateID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Updates reusable environment configuration without returning confidential
// values. See
// [reusing a hosted setup](https://developers.openai.com/api/docs/guides/agents-api/tools#reuse-a-hosted-plugin-setup).
func (r *BetaAgentEnvironmentTemplateService) Update(ctx context.Context, environmentTemplateID string, body BetaAgentEnvironmentTemplateUpdateParams, opts ...option.RequestOption) (res *EnvironmentTemplate, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "agents=v1")}, opts...)
	if environmentTemplateID == "" {
		err = errors.New("missing required environment_template_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("agents/environments/templates/%s", environmentTemplateID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Lists reusable environment templates without returning confidential values. See
// [reusing a hosted setup](https://developers.openai.com/api/docs/guides/agents-api/tools#reuse-a-hosted-plugin-setup).
func (r *BetaAgentEnvironmentTemplateService) List(ctx context.Context, query BetaAgentEnvironmentTemplateListParams, opts ...option.RequestOption) (res *pagination.CursorPage[EnvironmentTemplate], err error) {
	var raw *http.Response
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "agents=v1"), option.WithResponseInto(&raw)}, opts...)
	path := "agents/environments/templates"
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

// Lists reusable environment templates without returning confidential values. See
// [reusing a hosted setup](https://developers.openai.com/api/docs/guides/agents-api/tools#reuse-a-hosted-plugin-setup).
func (r *BetaAgentEnvironmentTemplateService) ListAutoPaging(ctx context.Context, query BetaAgentEnvironmentTemplateListParams, opts ...option.RequestOption) *pagination.CursorPageAutoPager[EnvironmentTemplate] {
	return pagination.NewCursorPageAutoPager(r.List(ctx, query, opts...))
}

// Deletes reusable environment configuration and all confidential template inputs.
// See
// [reusing a hosted setup](https://developers.openai.com/api/docs/guides/agents-api/tools#reuse-a-hosted-plugin-setup).
func (r *BetaAgentEnvironmentTemplateService) Delete(ctx context.Context, environmentTemplateID string, opts ...option.RequestOption) (res *EnvironmentTemplateDeleted, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "agents=v1")}, opts...)
	if environmentTemplateID == "" {
		err = errors.New("missing required environment_template_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("agents/environments/templates/%s", environmentTemplateID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Reusable configuration that provisions a fresh OpenAI-hosted environment for
// each session.
type EnvironmentTemplate struct {
	// The ID of the reusable environment template.
	ID string `json:"id" api:"required"`
	// Directories that expose capabilities to the agent.
	CapabilityDirectories []string `json:"capability_directories" api:"required"`
	// The Unix timestamp, in seconds, when the template was created.
	CreatedAt int64 `json:"created_at" api:"required"`
	// Safe file metadata, excluding contents and session-scoped file IDs.
	Files []EnvironmentTemplateFileUnion `json:"files" api:"required"`
	// An optional human-readable display name for the template.
	Name string `json:"name" api:"required"`
	// Runtime network access for each OpenAI-hosted environment.
	Network EnvironmentTemplateNetwork `json:"network" api:"required"`
	// The object type. Always `agent.environment.template`.
	Object constant.AgentEnvironmentTemplate `json:"object" default:"agent.environment.template"`
	// Packages installed in each fresh OpenAI-hosted environment.
	Packages EnvironmentTemplatePackages `json:"packages" api:"required"`
	// Safe plugin metadata, excluding inline archive contents.
	Plugins []HostedPlugin `json:"plugins" api:"required"`
	// Safe skill metadata, preserving unresolved version selectors.
	Skills []EnvironmentTemplateSkillUnion `json:"skills" api:"required"`
	// The Unix timestamp, in seconds, when the template was last updated.
	UpdatedAt int64 `json:"updated_at" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                    respjson.Field
		CapabilityDirectories respjson.Field
		CreatedAt             respjson.Field
		Files                 respjson.Field
		Name                  respjson.Field
		Network               respjson.Field
		Object                respjson.Field
		Packages              respjson.Field
		Plugins               respjson.Field
		Skills                respjson.Field
		UpdatedAt             respjson.Field
		ExtraFields           map[string]respjson.Field
		raw                   string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EnvironmentTemplate) RawJSON() string { return r.JSON.raw }
func (r *EnvironmentTemplate) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// EnvironmentTemplateFileUnion contains all possible properties and values from
// [EnvironmentTemplateFileFileID], [EnvironmentTemplateFileInline].
//
// Use the [EnvironmentTemplateFileUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type EnvironmentTemplateFileUnion struct {
	// This field is from variant [EnvironmentTemplateFileFileID].
	FileID string `json:"file_id"`
	Path   string `json:"path"`
	// Any of "file_id", "inline".
	Type string `json:"type"`
	// This field is from variant [EnvironmentTemplateFileInline].
	SizeBytes int64 `json:"size_bytes"`
	JSON      struct {
		FileID    respjson.Field
		Path      respjson.Field
		Type      respjson.Field
		SizeBytes respjson.Field
		raw       string
	} `json:"-"`
}

// anyEnvironmentTemplateFile is implemented by each variant of
// [EnvironmentTemplateFileUnion] to add type safety for the return type of
// [EnvironmentTemplateFileUnion.AsAny]
type anyEnvironmentTemplateFile interface {
	implEnvironmentTemplateFileUnion()
}

func (EnvironmentTemplateFileFileID) implEnvironmentTemplateFileUnion() {}
func (EnvironmentTemplateFileInline) implEnvironmentTemplateFileUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := EnvironmentTemplateFileUnion.AsAny().(type) {
//	case openai.EnvironmentTemplateFileFileID:
//	case openai.EnvironmentTemplateFileInline:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u EnvironmentTemplateFileUnion) AsAny() anyEnvironmentTemplateFile {
	switch u.Type {
	case "file_id":
		return u.AsFileID()
	case "inline":
		return u.AsInline()
	}
	return nil
}

func (u EnvironmentTemplateFileUnion) AsFileID() (v EnvironmentTemplateFileFileID) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u EnvironmentTemplateFileUnion) AsInline() (v EnvironmentTemplateFileInline) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u EnvironmentTemplateFileUnion) RawJSON() string { return u.JSON.raw }

func (r *EnvironmentTemplateFileUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A project-scoped Files API reference resolved separately for each session.
type EnvironmentTemplateFileFileID struct {
	// The ID of the uploaded file.
	FileID string `json:"file_id" api:"required"`
	// The file's absolute path inside the environment.
	Path string `json:"path" api:"required"`
	// The type of the object. Always `file_id`.
	Type constant.FileID `json:"type" default:"file_id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		FileID      respjson.Field
		Path        respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EnvironmentTemplateFileFileID) RawJSON() string { return r.JSON.raw }
func (r *EnvironmentTemplateFileFileID) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Metadata for confidential inline file contents.
type EnvironmentTemplateFileInline struct {
	// The file's absolute path inside the environment.
	Path string `json:"path" api:"required"`
	// The decoded size of the inline file in bytes.
	SizeBytes int64 `json:"size_bytes" api:"required"`
	// The type of the object. Always `inline`.
	Type constant.Inline `json:"type" default:"inline"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Path        respjson.Field
		SizeBytes   respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EnvironmentTemplateFileInline) RawJSON() string { return r.JSON.raw }
func (r *EnvironmentTemplateFileInline) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Runtime network access for each OpenAI-hosted environment.
type EnvironmentTemplateNetwork struct {
	// The environment's network access mode.
	//
	// Any of "enabled", "disabled", "restricted".
	Access string `json:"access" api:"required"`
	// Domains the environment may access when network access is restricted.
	AllowedDomains []string `json:"allowed_domains" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Access         respjson.Field
		AllowedDomains respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EnvironmentTemplateNetwork) RawJSON() string { return r.JSON.raw }
func (r *EnvironmentTemplateNetwork) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Packages installed in each fresh OpenAI-hosted environment.
type EnvironmentTemplatePackages struct {
	// npm packages installed globally in the environment.
	Npm []string `json:"npm" api:"required"`
	// Python packages installed in the environment.
	Python []string `json:"python" api:"required"`
	// System packages installed in the environment.
	System []string `json:"system" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Npm         respjson.Field
		Python      respjson.Field
		System      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EnvironmentTemplatePackages) RawJSON() string { return r.JSON.raw }
func (r *EnvironmentTemplatePackages) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// EnvironmentTemplateSkillUnion contains all possible properties and values from
// [EnvironmentTemplateSkillSkillReference], [EnvironmentTemplateSkillInline].
//
// Use the [EnvironmentTemplateSkillUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type EnvironmentTemplateSkillUnion struct {
	// This field is from variant [EnvironmentTemplateSkillSkillReference].
	SkillID string `json:"skill_id"`
	// Any of "skill_reference", "inline".
	Type string `json:"type"`
	// This field is from variant [EnvironmentTemplateSkillSkillReference].
	Version string `json:"version"`
	// This field is from variant [EnvironmentTemplateSkillInline].
	Description string `json:"description"`
	// This field is from variant [EnvironmentTemplateSkillInline].
	Name string `json:"name"`
	JSON struct {
		SkillID     respjson.Field
		Type        respjson.Field
		Version     respjson.Field
		Description respjson.Field
		Name        respjson.Field
		raw         string
	} `json:"-"`
}

// anyEnvironmentTemplateSkill is implemented by each variant of
// [EnvironmentTemplateSkillUnion] to add type safety for the return type of
// [EnvironmentTemplateSkillUnion.AsAny]
type anyEnvironmentTemplateSkill interface {
	implEnvironmentTemplateSkillUnion()
}

func (EnvironmentTemplateSkillSkillReference) implEnvironmentTemplateSkillUnion() {}
func (EnvironmentTemplateSkillInline) implEnvironmentTemplateSkillUnion()         {}

// Use the following switch statement to find the correct variant
//
//	switch variant := EnvironmentTemplateSkillUnion.AsAny().(type) {
//	case openai.EnvironmentTemplateSkillSkillReference:
//	case openai.EnvironmentTemplateSkillInline:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u EnvironmentTemplateSkillUnion) AsAny() anyEnvironmentTemplateSkill {
	switch u.Type {
	case "skill_reference":
		return u.AsSkillReference()
	case "inline":
		return u.AsInline()
	}
	return nil
}

func (u EnvironmentTemplateSkillUnion) AsSkillReference() (v EnvironmentTemplateSkillSkillReference) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u EnvironmentTemplateSkillUnion) AsInline() (v EnvironmentTemplateSkillInline) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u EnvironmentTemplateSkillUnion) RawJSON() string { return u.JSON.raw }

func (r *EnvironmentTemplateSkillUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A skill resolved afresh from the Skills API whenever a session starts.
type EnvironmentTemplateSkillSkillReference struct {
	// The referenced skill ID.
	SkillID string `json:"skill_id" api:"required"`
	// The type of the object. Always `skill_reference`.
	Type constant.SkillReference `json:"type" default:"skill_reference"`
	// The requested version selector, including `latest`.
	Version string `json:"version" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		SkillID     respjson.Field
		Type        respjson.Field
		Version     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EnvironmentTemplateSkillSkillReference) RawJSON() string { return r.JSON.raw }
func (r *EnvironmentTemplateSkillSkillReference) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Safe metadata for an inline skill archive.
type EnvironmentTemplateSkillInline struct {
	// The skill description declared in `SKILL.md`.
	Description string `json:"description" api:"required"`
	// The skill name declared in `SKILL.md`.
	Name string `json:"name" api:"required"`
	// The type of the object. Always `inline`.
	Type constant.Inline `json:"type" default:"inline"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Description respjson.Field
		Name        respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EnvironmentTemplateSkillInline) RawJSON() string { return r.JSON.raw }
func (r *EnvironmentTemplateSkillInline) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A deleted reusable environment template.
type EnvironmentTemplateDeleted struct {
	// The ID of the deleted environment template.
	ID string `json:"id" api:"required"`
	// Whether the environment template was deleted. Always `true`.
	Deleted bool `json:"deleted" api:"required"`
	// The object type. Always `agent.environment.template.deleted`.
	Object constant.AgentEnvironmentTemplateDeleted `json:"object" default:"agent.environment.template.deleted"`
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
func (r EnvironmentTemplateDeleted) RawJSON() string { return r.JSON.raw }
func (r *EnvironmentTemplateDeleted) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaAgentEnvironmentTemplateNewParams struct {
	// An optional human-readable display name for the template.
	Name param.Opt[string] `json:"name,omitzero"`
	// Directories that contain capabilities exposed to the agent. Defaults to an empty
	// list.
	CapabilityDirectories []string `json:"capability_directories,omitzero"`
	// Environment variables made available to the agent.
	Env map[string]string `json:"env,omitzero"`
	// Files available before the agent starts. Defaults to an empty list.
	Files []HostedEnvironmentFileParamUnion `json:"files,omitzero"`
	// Network access for an OpenAI-hosted environment.
	Network BetaAgentEnvironmentTemplateNewParamsNetwork `json:"network,omitzero"`
	// Packages to install in an OpenAI-hosted environment.
	Packages BetaAgentEnvironmentTemplateNewParamsPackages `json:"packages,omitzero"`
	// Plugins provided as inline ZIP archives. Defaults to an empty list.
	Plugins []HostedPluginParam `json:"plugins,omitzero"`
	// Ordered, confidential setup commands. Command bodies are never returned.
	SetupCommands []SetupCommandParam `json:"setup_commands,omitzero"`
	// Skills referenced by ID or provided as inline ZIP archives. Defaults to an empty
	// list.
	Skills []HostedSkillParamUnion `json:"skills,omitzero"`
	paramObj
}

func (r BetaAgentEnvironmentTemplateNewParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaAgentEnvironmentTemplateNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaAgentEnvironmentTemplateNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Network access for an OpenAI-hosted environment.
//
// The property Access is required.
type BetaAgentEnvironmentTemplateNewParamsNetwork struct {
	// The environment's network access mode.
	//
	// Any of "enabled", "disabled", "restricted".
	Access string `json:"access,omitzero" api:"required"`
	// Domains the environment may access when network access is restricted.
	AllowedDomains []string `json:"allowed_domains,omitzero"`
	paramObj
}

func (r BetaAgentEnvironmentTemplateNewParamsNetwork) MarshalJSON() (data []byte, err error) {
	type shadow BetaAgentEnvironmentTemplateNewParamsNetwork
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaAgentEnvironmentTemplateNewParamsNetwork) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[BetaAgentEnvironmentTemplateNewParamsNetwork](
		"access", "enabled", "disabled", "restricted",
	)
}

// Packages to install in an OpenAI-hosted environment.
type BetaAgentEnvironmentTemplateNewParamsPackages struct {
	// npm packages to install globally. Defaults to an empty list.
	Npm []string `json:"npm,omitzero"`
	// Python packages to install. Defaults to an empty list.
	Python []string `json:"python,omitzero"`
	// System packages to install. Defaults to an empty list.
	System []string `json:"system,omitzero"`
	paramObj
}

func (r BetaAgentEnvironmentTemplateNewParamsPackages) MarshalJSON() (data []byte, err error) {
	type shadow BetaAgentEnvironmentTemplateNewParamsPackages
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaAgentEnvironmentTemplateNewParamsPackages) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaAgentEnvironmentTemplateUpdateParams struct {
	// A replacement human-readable display name, or `null` to clear the name.
	Name param.Opt[string] `json:"name,omitzero"`
	// Directories that expose capabilities to the agent.
	CapabilityDirectories []string `json:"capability_directories,omitzero"`
	// Replacement confidential environment values.
	Env map[string]string `json:"env,omitzero"`
	// Replacement file configuration materialized for each new session.
	Files []HostedEnvironmentFileParamUnion `json:"files,omitzero"`
	// Network access for an OpenAI-hosted environment.
	Network BetaAgentEnvironmentTemplateUpdateParamsNetwork `json:"network,omitzero"`
	// Packages to install in an OpenAI-hosted environment.
	Packages BetaAgentEnvironmentTemplateUpdateParamsPackages `json:"packages,omitzero"`
	// Replacement plugin configuration installed for each new session.
	Plugins []HostedPluginParam `json:"plugins,omitzero"`
	// Replacement confidential setup commands, never included in returned resources.
	SetupCommands []SetupCommandParam `json:"setup_commands,omitzero"`
	// Replacement skill configuration installed for each new session.
	Skills []HostedSkillParamUnion `json:"skills,omitzero"`
	paramObj
}

func (r BetaAgentEnvironmentTemplateUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaAgentEnvironmentTemplateUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaAgentEnvironmentTemplateUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Network access for an OpenAI-hosted environment.
//
// The property Access is required.
type BetaAgentEnvironmentTemplateUpdateParamsNetwork struct {
	// The environment's network access mode.
	//
	// Any of "enabled", "disabled", "restricted".
	Access string `json:"access,omitzero" api:"required"`
	// Domains the environment may access when network access is restricted.
	AllowedDomains []string `json:"allowed_domains,omitzero"`
	paramObj
}

func (r BetaAgentEnvironmentTemplateUpdateParamsNetwork) MarshalJSON() (data []byte, err error) {
	type shadow BetaAgentEnvironmentTemplateUpdateParamsNetwork
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaAgentEnvironmentTemplateUpdateParamsNetwork) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[BetaAgentEnvironmentTemplateUpdateParamsNetwork](
		"access", "enabled", "disabled", "restricted",
	)
}

// Packages to install in an OpenAI-hosted environment.
type BetaAgentEnvironmentTemplateUpdateParamsPackages struct {
	// npm packages to install globally. Defaults to an empty list.
	Npm []string `json:"npm,omitzero"`
	// Python packages to install. Defaults to an empty list.
	Python []string `json:"python,omitzero"`
	// System packages to install. Defaults to an empty list.
	System []string `json:"system,omitzero"`
	paramObj
}

func (r BetaAgentEnvironmentTemplateUpdateParamsPackages) MarshalJSON() (data []byte, err error) {
	type shadow BetaAgentEnvironmentTemplateUpdateParamsPackages
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaAgentEnvironmentTemplateUpdateParamsPackages) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaAgentEnvironmentTemplateListParams struct {
	// Return resources after this resource ID in the selected order.
	After param.Opt[string] `query:"after,omitzero" json:"-"`
	// The maximum number of resources to return, between 1 and 100. Defaults to 20.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// The order in which resources are returned. Defaults to `desc`.
	//
	// Any of "asc", "desc".
	Order BetaAgentEnvironmentTemplateListParamsOrder `query:"order,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaAgentEnvironmentTemplateListParams]'s query parameters
// as `url.Values`.
func (r BetaAgentEnvironmentTemplateListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// The order in which resources are returned. Defaults to `desc`.
type BetaAgentEnvironmentTemplateListParamsOrder string

const (
	BetaAgentEnvironmentTemplateListParamsOrderAsc  BetaAgentEnvironmentTemplateListParamsOrder = "asc"
	BetaAgentEnvironmentTemplateListParamsOrderDesc BetaAgentEnvironmentTemplateListParamsOrder = "desc"
)

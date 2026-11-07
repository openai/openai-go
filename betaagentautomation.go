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

// BetaAgentAutomationService contains methods and other services that help with
// interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaAgentAutomationService] method instead.
type BetaAgentAutomationService struct {
	Options     []option.RequestOption
	Invocations BetaAgentAutomationInvocationService
}

// NewBetaAgentAutomationService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewBetaAgentAutomationService(opts ...option.RequestOption) (r BetaAgentAutomationService) {
	r = BetaAgentAutomationService{}
	r.Options = requestconfig.InheritedOptions(opts...)
	r.Invocations = NewBetaAgentAutomationInvocationService(opts...)
	return
}

// Saves an active automation without invoking it immediately.
func (r *BetaAgentAutomationService) New(ctx context.Context, body BetaAgentAutomationNewParams, opts ...option.RequestOption) (res *AutomationUnion, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "agents=v1")}, opts...)
	path := "agents/automations"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Retrieves a project-owned automation.
func (r *BetaAgentAutomationService) Get(ctx context.Context, automationID string, opts ...option.RequestOption) (res *AutomationUnion, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "agents=v1")}, opts...)
	if automationID == "" {
		err = errors.New("missing required automation_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("agents/automations/%s", automationID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Updates supplied fields; target and metadata replace the whole field.
func (r *BetaAgentAutomationService) Update(ctx context.Context, automationID string, body BetaAgentAutomationUpdateParams, opts ...option.RequestOption) (res *AutomationUnion, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "agents=v1")}, opts...)
	if automationID == "" {
		err = errors.New("missing required automation_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("agents/automations/%s", automationID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Lists project-owned automations, optionally filtered by target session.
func (r *BetaAgentAutomationService) List(ctx context.Context, query BetaAgentAutomationListParams, opts ...option.RequestOption) (res *pagination.CursorPage[AutomationUnion], err error) {
	var raw *http.Response
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "agents=v1"), option.WithResponseInto(&raw)}, opts...)
	path := "agents/automations"
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

// Lists project-owned automations, optionally filtered by target session.
func (r *BetaAgentAutomationService) ListAutoPaging(ctx context.Context, query BetaAgentAutomationListParams, opts ...option.RequestOption) *pagination.CursorPageAutoPager[AutomationUnion] {
	return pagination.NewCursorPageAutoPager(r.List(ctx, query, opts...))
}

// Deletes the automation and prevents future invocations without deleting its
// target session.
func (r *BetaAgentAutomationService) Delete(ctx context.Context, automationID string, opts ...option.RequestOption) (res *AutomationDeleted, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "agents=v1")}, opts...)
	if automationID == "" {
		err = errors.New("missing required automation_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("agents/automations/%s", automationID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Pauses future scheduled delivery without stopping an active session.
func (r *BetaAgentAutomationService) Pause(ctx context.Context, automationID string, opts ...option.RequestOption) (res *AutomationUnion, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "agents=v1")}, opts...)
	if automationID == "" {
		err = errors.New("missing required automation_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("agents/automations/%s/pause", automationID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// Resumes scheduled delivery; one-time automations without a pending invocation
// require a future run_at.
func (r *BetaAgentAutomationService) Resume(ctx context.Context, automationID string, opts ...option.RequestOption) (res *AutomationUnion, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "agents=v1")}, opts...)
	if automationID == "" {
		err = errors.New("missing required automation_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("agents/automations/%s/resume", automationID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// AutomationUnion contains all possible properties and values from
// [PeriodicAutomation], [OneTimeAutomation].
//
// Use the [AutomationUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type AutomationUnion struct {
	ID        string `json:"id"`
	CreatedAt int64  `json:"created_at"`
	// This field is from variant [PeriodicAutomation].
	CreatedBy AutomationCreatedByUnion `json:"created_by"`
	// This field is from variant [PeriodicAutomation].
	Expression string            `json:"expression"`
	Metadata   map[string]string `json:"metadata"`
	Name       string            `json:"name"`
	// This field is from variant [PeriodicAutomation].
	Object constant.AgentAutomation `json:"object"`
	Status string                   `json:"status"`
	// This field is from variant [PeriodicAutomation].
	Target AutomationTargetUnion `json:"target"`
	// This field is from variant [PeriodicAutomation].
	Timezone string `json:"timezone"`
	// Any of "periodic", "once".
	Type string `json:"type"`
	// This field is from variant [OneTimeAutomation].
	RunAt int64 `json:"run_at"`
	JSON  struct {
		ID         respjson.Field
		CreatedAt  respjson.Field
		CreatedBy  respjson.Field
		Expression respjson.Field
		Metadata   respjson.Field
		Name       respjson.Field
		Object     respjson.Field
		Status     respjson.Field
		Target     respjson.Field
		Timezone   respjson.Field
		Type       respjson.Field
		RunAt      respjson.Field
		raw        string
	} `json:"-"`
}

// anyAutomation is implemented by each variant of [AutomationUnion] to add type
// safety for the return type of [AutomationUnion.AsAny]
type anyAutomation interface {
	implAutomationUnion()
}

func (PeriodicAutomation) implAutomationUnion() {}
func (OneTimeAutomation) implAutomationUnion()  {}

// Use the following switch statement to find the correct variant
//
//	switch variant := AutomationUnion.AsAny().(type) {
//	case openai.PeriodicAutomation:
//	case openai.OneTimeAutomation:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u AutomationUnion) AsAny() anyAutomation {
	switch u.Type {
	case "periodic":
		return u.AsPeriodic()
	case "once":
		return u.AsOnce()
	}
	return nil
}

func (u AutomationUnion) AsPeriodic() (v PeriodicAutomation) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AutomationUnion) AsOnce() (v OneTimeAutomation) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u AutomationUnion) RawJSON() string { return u.JSON.raw }

func (r *AutomationUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Inline agent configuration without a generated session-local identity.
type AutomationAgentConfig struct {
	// Custom instructions, or null when absent.
	Instructions string `json:"instructions" api:"required"`
	// The requested model name.
	Model string `json:"model" api:"required"`
	// The configuration for creating and coordinating subagents.
	MultiAgent MultiAgentConfig `json:"multi_agent" api:"required"`
	// The configured reasoning policy.
	Reasoning AgentReasoning `json:"reasoning" api:"required"`
	// The configured service-tier policy.
	//
	// Any of "auto", "default", "flex", "priority", "fast", "ultrafast".
	ServiceTier AutomationAgentConfigServiceTier `json:"service_tier" api:"required"`
	// The configured output format and verbosity.
	Text AgentText `json:"text" api:"required"`
	// Configured tools, excluding confidential credentials.
	Tools []AgentToolUnion `json:"tools" api:"required"`
	// The effective context compaction settings.
	Compaction AutomationAgentConfigCompaction `json:"compaction"`
	// The effective experimental settings for the agent.
	Experimental AutomationAgentConfigExperimental `json:"experimental"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Instructions respjson.Field
		Model        respjson.Field
		MultiAgent   respjson.Field
		Reasoning    respjson.Field
		ServiceTier  respjson.Field
		Text         respjson.Field
		Tools        respjson.Field
		Compaction   respjson.Field
		Experimental respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AutomationAgentConfig) RawJSON() string { return r.JSON.raw }
func (r *AutomationAgentConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The configured service-tier policy.
type AutomationAgentConfigServiceTier string

const (
	AutomationAgentConfigServiceTierAuto      AutomationAgentConfigServiceTier = "auto"
	AutomationAgentConfigServiceTierDefault   AutomationAgentConfigServiceTier = "default"
	AutomationAgentConfigServiceTierFlex      AutomationAgentConfigServiceTier = "flex"
	AutomationAgentConfigServiceTierPriority  AutomationAgentConfigServiceTier = "priority"
	AutomationAgentConfigServiceTierFast      AutomationAgentConfigServiceTier = "fast"
	AutomationAgentConfigServiceTierUltrafast AutomationAgentConfigServiceTier = "ultrafast"
)

// The effective context compaction settings.
type AutomationAgentConfigCompaction struct {
	// Whether context is compacted in the background. Defaults to true.
	Async bool `json:"async" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Async       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AutomationAgentConfigCompaction) RawJSON() string { return r.JSON.raw }
func (r *AutomationAgentConfigCompaction) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The effective experimental settings for the agent.
type AutomationAgentConfigExperimental struct {
	// Whether the model receives periodic current-time reminders.
	TimeAwareness bool `json:"time_awareness" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		TimeAwareness respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AutomationAgentConfigExperimental) RawJSON() string { return r.JSON.raw }
func (r *AutomationAgentConfigExperimental) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// AutomationCreatedByUnion contains all possible properties and values from
// [AutomationCreatedByDeveloper], [AutomationCreatedByAgent].
//
// Use the [AutomationCreatedByUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type AutomationCreatedByUnion struct {
	// Any of "developer", "agent".
	Type string `json:"type"`
	// This field is from variant [AutomationCreatedByAgent].
	SessionID string `json:"session_id"`
	// This field is from variant [AutomationCreatedByAgent].
	ToolCallID string `json:"tool_call_id"`
	JSON       struct {
		Type       respjson.Field
		SessionID  respjson.Field
		ToolCallID respjson.Field
		raw        string
	} `json:"-"`
}

// anyAutomationCreatedBy is implemented by each variant of
// [AutomationCreatedByUnion] to add type safety for the return type of
// [AutomationCreatedByUnion.AsAny]
type anyAutomationCreatedBy interface {
	implAutomationCreatedByUnion()
}

func (AutomationCreatedByDeveloper) implAutomationCreatedByUnion() {}
func (AutomationCreatedByAgent) implAutomationCreatedByUnion()     {}

// Use the following switch statement to find the correct variant
//
//	switch variant := AutomationCreatedByUnion.AsAny().(type) {
//	case openai.AutomationCreatedByDeveloper:
//	case openai.AutomationCreatedByAgent:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u AutomationCreatedByUnion) AsAny() anyAutomationCreatedBy {
	switch u.Type {
	case "developer":
		return u.AsDeveloper()
	case "agent":
		return u.AsAgent()
	}
	return nil
}

func (u AutomationCreatedByUnion) AsDeveloper() (v AutomationCreatedByDeveloper) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AutomationCreatedByUnion) AsAgent() (v AutomationCreatedByAgent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u AutomationCreatedByUnion) RawJSON() string { return u.JSON.raw }

func (r *AutomationCreatedByUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AutomationCreatedByDeveloper struct {
	// The type of the object. Always `developer`.
	Type constant.Developer `json:"type" default:"developer"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AutomationCreatedByDeveloper) RawJSON() string { return r.JSON.raw }
func (r *AutomationCreatedByDeveloper) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AutomationCreatedByAgent struct {
	// The creator session, which is also the automation's target.
	SessionID string `json:"session_id" api:"required"`
	// The tool call that created the automation.
	ToolCallID string `json:"tool_call_id" api:"required"`
	// The type of the object. Always `agent`.
	Type constant.Agent `json:"type" default:"agent"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		SessionID   respjson.Field
		ToolCallID  respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AutomationCreatedByAgent) RawJSON() string { return r.JSON.raw }
func (r *AutomationCreatedByAgent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AutomationDeleted struct {
	// The deleted automation ID.
	ID string `json:"id" api:"required"`
	// Always true.
	Deleted bool `json:"deleted" api:"required"`
	// Always `agent.automation.deleted`.
	Object constant.AgentAutomationDeleted `json:"object" default:"agent.automation.deleted"`
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
func (r AutomationDeleted) RawJSON() string { return r.JSON.raw }
func (r *AutomationDeleted) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// AutomationEnvironmentUnion contains all possible properties and values from
// [AutomationEnvironmentNone], [AutomationEnvironmentOpenAIHosted],
// [AutomationEnvironmentSelfHosted].
//
// Use the [AutomationEnvironmentUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type AutomationEnvironmentUnion struct {
	// Any of "none", "openai_hosted", "self_hosted".
	Type                  string   `json:"type"`
	CapabilityDirectories []string `json:"capability_directories"`
	// This field is from variant [AutomationEnvironmentOpenAIHosted].
	ContainerSize string `json:"container_size"`
	// This field is from variant [AutomationEnvironmentOpenAIHosted].
	Desktop       AutomationEnvironmentOpenAIHostedDesktop `json:"desktop"`
	EnvironmentID string                                   `json:"environment_id"`
	// This field is from variant [AutomationEnvironmentOpenAIHosted].
	EnvironmentTemplateID string `json:"environment_template_id"`
	// This field is from variant [AutomationEnvironmentOpenAIHosted].
	Files []AutomationEnvironmentOpenAIHostedFileUnion `json:"files"`
	// This field is from variant [AutomationEnvironmentOpenAIHosted].
	Network AutomationEnvironmentOpenAIHostedNetwork `json:"network"`
	// This field is from variant [AutomationEnvironmentOpenAIHosted].
	Packages AutomationEnvironmentOpenAIHostedPackages `json:"packages"`
	// This field is from variant [AutomationEnvironmentOpenAIHosted].
	Plugins []HostedPlugin `json:"plugins"`
	// This field is from variant [AutomationEnvironmentOpenAIHosted].
	Skills []AutomationEnvironmentOpenAIHostedSkillUnion `json:"skills"`
	// This field is from variant [AutomationEnvironmentSelfHosted].
	WorkspaceDirectory string `json:"workspace_directory"`
	JSON               struct {
		Type                  respjson.Field
		CapabilityDirectories respjson.Field
		ContainerSize         respjson.Field
		Desktop               respjson.Field
		EnvironmentID         respjson.Field
		EnvironmentTemplateID respjson.Field
		Files                 respjson.Field
		Network               respjson.Field
		Packages              respjson.Field
		Plugins               respjson.Field
		Skills                respjson.Field
		WorkspaceDirectory    respjson.Field
		raw                   string
	} `json:"-"`
}

// anyAutomationEnvironment is implemented by each variant of
// [AutomationEnvironmentUnion] to add type safety for the return type of
// [AutomationEnvironmentUnion.AsAny]
type anyAutomationEnvironment interface {
	implAutomationEnvironmentUnion()
}

func (AutomationEnvironmentNone) implAutomationEnvironmentUnion()         {}
func (AutomationEnvironmentOpenAIHosted) implAutomationEnvironmentUnion() {}
func (AutomationEnvironmentSelfHosted) implAutomationEnvironmentUnion()   {}

// Use the following switch statement to find the correct variant
//
//	switch variant := AutomationEnvironmentUnion.AsAny().(type) {
//	case openai.AutomationEnvironmentNone:
//	case openai.AutomationEnvironmentOpenAIHosted:
//	case openai.AutomationEnvironmentSelfHosted:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u AutomationEnvironmentUnion) AsAny() anyAutomationEnvironment {
	switch u.Type {
	case "none":
		return u.AsNone()
	case "openai_hosted":
		return u.AsOpenAIHosted()
	case "self_hosted":
		return u.AsSelfHosted()
	}
	return nil
}

func (u AutomationEnvironmentUnion) AsNone() (v AutomationEnvironmentNone) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AutomationEnvironmentUnion) AsOpenAIHosted() (v AutomationEnvironmentOpenAIHosted) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AutomationEnvironmentUnion) AsSelfHosted() (v AutomationEnvironmentSelfHosted) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u AutomationEnvironmentUnion) RawJSON() string { return u.JSON.raw }

func (r *AutomationEnvironmentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AutomationEnvironmentNone struct {
	// The type of the object. Always `none`.
	Type constant.None `json:"type" default:"none"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AutomationEnvironmentNone) RawJSON() string { return r.JSON.raw }
func (r *AutomationEnvironmentNone) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Hosted environment selection with confidential setup commands and variables
// omitted.
type AutomationEnvironmentOpenAIHosted struct {
	// Directory overrides, or null to inherit the template or creation default.
	CapabilityDirectories []string `json:"capability_directories" api:"required"`
	// The configured container size, or null to use the creation default.
	//
	// Any of "small", "medium", "large".
	ContainerSize string `json:"container_size" api:"required"`
	// Desktop overrides, or null to inherit the template or creation default.
	Desktop AutomationEnvironmentOpenAIHostedDesktop `json:"desktop" api:"required"`
	// An existing environment reference, or null when provisioning a new environment.
	EnvironmentID string `json:"environment_id" api:"required"`
	// A template reference, or null when no template is selected.
	EnvironmentTemplateID string `json:"environment_template_id" api:"required"`
	// File overrides, excluding contents, or null to inherit.
	Files []AutomationEnvironmentOpenAIHostedFileUnion `json:"files" api:"required"`
	// Network overrides, or null to inherit the template or creation default.
	Network AutomationEnvironmentOpenAIHostedNetwork `json:"network" api:"required"`
	// Package overrides, or null to inherit the template or creation default.
	Packages AutomationEnvironmentOpenAIHostedPackages `json:"packages" api:"required"`
	// Plugin overrides, excluding archive contents, or null to inherit.
	Plugins []HostedPlugin `json:"plugins" api:"required"`
	// Skill overrides, excluding archive contents, or null to inherit.
	Skills []AutomationEnvironmentOpenAIHostedSkillUnion `json:"skills" api:"required"`
	// The type of the object. Always `openai_hosted`.
	Type constant.OpenAIHosted `json:"type" default:"openai_hosted"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		CapabilityDirectories respjson.Field
		ContainerSize         respjson.Field
		Desktop               respjson.Field
		EnvironmentID         respjson.Field
		EnvironmentTemplateID respjson.Field
		Files                 respjson.Field
		Network               respjson.Field
		Packages              respjson.Field
		Plugins               respjson.Field
		Skills                respjson.Field
		Type                  respjson.Field
		ExtraFields           map[string]respjson.Field
		raw                   string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AutomationEnvironmentOpenAIHosted) RawJSON() string { return r.JSON.raw }
func (r *AutomationEnvironmentOpenAIHosted) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Desktop overrides, or null to inherit the template or creation default.
type AutomationEnvironmentOpenAIHostedDesktop struct {
	// Whether the environment provisions a desktop and browser proxy.
	Enabled bool `json:"enabled" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Enabled     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AutomationEnvironmentOpenAIHostedDesktop) RawJSON() string { return r.JSON.raw }
func (r *AutomationEnvironmentOpenAIHostedDesktop) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// AutomationEnvironmentOpenAIHostedFileUnion contains all possible properties and
// values from [AutomationEnvironmentOpenAIHostedFileFileID],
// [AutomationEnvironmentOpenAIHostedFileInline].
//
// Use the [AutomationEnvironmentOpenAIHostedFileUnion.AsAny] method to switch on
// the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type AutomationEnvironmentOpenAIHostedFileUnion struct {
	// This field is from variant [AutomationEnvironmentOpenAIHostedFileFileID].
	FileID string `json:"file_id"`
	Path   string `json:"path"`
	// Any of "file_id", "inline".
	Type string `json:"type"`
	// This field is from variant [AutomationEnvironmentOpenAIHostedFileInline].
	SizeBytes int64 `json:"size_bytes"`
	JSON      struct {
		FileID    respjson.Field
		Path      respjson.Field
		Type      respjson.Field
		SizeBytes respjson.Field
		raw       string
	} `json:"-"`
}

// anyAutomationEnvironmentOpenAIHostedFile is implemented by each variant of
// [AutomationEnvironmentOpenAIHostedFileUnion] to add type safety for the return
// type of [AutomationEnvironmentOpenAIHostedFileUnion.AsAny]
type anyAutomationEnvironmentOpenAIHostedFile interface {
	implAutomationEnvironmentOpenAIHostedFileUnion()
}

func (AutomationEnvironmentOpenAIHostedFileFileID) implAutomationEnvironmentOpenAIHostedFileUnion() {}
func (AutomationEnvironmentOpenAIHostedFileInline) implAutomationEnvironmentOpenAIHostedFileUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := AutomationEnvironmentOpenAIHostedFileUnion.AsAny().(type) {
//	case openai.AutomationEnvironmentOpenAIHostedFileFileID:
//	case openai.AutomationEnvironmentOpenAIHostedFileInline:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u AutomationEnvironmentOpenAIHostedFileUnion) AsAny() anyAutomationEnvironmentOpenAIHostedFile {
	switch u.Type {
	case "file_id":
		return u.AsFileID()
	case "inline":
		return u.AsInline()
	}
	return nil
}

func (u AutomationEnvironmentOpenAIHostedFileUnion) AsFileID() (v AutomationEnvironmentOpenAIHostedFileFileID) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AutomationEnvironmentOpenAIHostedFileUnion) AsInline() (v AutomationEnvironmentOpenAIHostedFileInline) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u AutomationEnvironmentOpenAIHostedFileUnion) RawJSON() string { return u.JSON.raw }

func (r *AutomationEnvironmentOpenAIHostedFileUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A project-scoped Files API reference resolved separately for each session.
type AutomationEnvironmentOpenAIHostedFileFileID struct {
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
func (r AutomationEnvironmentOpenAIHostedFileFileID) RawJSON() string { return r.JSON.raw }
func (r *AutomationEnvironmentOpenAIHostedFileFileID) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Metadata for confidential inline file contents.
type AutomationEnvironmentOpenAIHostedFileInline struct {
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
func (r AutomationEnvironmentOpenAIHostedFileInline) RawJSON() string { return r.JSON.raw }
func (r *AutomationEnvironmentOpenAIHostedFileInline) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Network overrides, or null to inherit the template or creation default.
type AutomationEnvironmentOpenAIHostedNetwork struct {
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
func (r AutomationEnvironmentOpenAIHostedNetwork) RawJSON() string { return r.JSON.raw }
func (r *AutomationEnvironmentOpenAIHostedNetwork) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Package overrides, or null to inherit the template or creation default.
type AutomationEnvironmentOpenAIHostedPackages struct {
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
func (r AutomationEnvironmentOpenAIHostedPackages) RawJSON() string { return r.JSON.raw }
func (r *AutomationEnvironmentOpenAIHostedPackages) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// AutomationEnvironmentOpenAIHostedSkillUnion contains all possible properties and
// values from [AutomationEnvironmentOpenAIHostedSkillSkillReference],
// [AutomationEnvironmentOpenAIHostedSkillInline].
//
// Use the [AutomationEnvironmentOpenAIHostedSkillUnion.AsAny] method to switch on
// the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type AutomationEnvironmentOpenAIHostedSkillUnion struct {
	// This field is from variant
	// [AutomationEnvironmentOpenAIHostedSkillSkillReference].
	SkillID string `json:"skill_id"`
	// Any of "skill_reference", "inline".
	Type string `json:"type"`
	// This field is from variant
	// [AutomationEnvironmentOpenAIHostedSkillSkillReference].
	Version string `json:"version"`
	// This field is from variant [AutomationEnvironmentOpenAIHostedSkillInline].
	Description string `json:"description"`
	// This field is from variant [AutomationEnvironmentOpenAIHostedSkillInline].
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

// anyAutomationEnvironmentOpenAIHostedSkill is implemented by each variant of
// [AutomationEnvironmentOpenAIHostedSkillUnion] to add type safety for the return
// type of [AutomationEnvironmentOpenAIHostedSkillUnion.AsAny]
type anyAutomationEnvironmentOpenAIHostedSkill interface {
	implAutomationEnvironmentOpenAIHostedSkillUnion()
}

func (AutomationEnvironmentOpenAIHostedSkillSkillReference) implAutomationEnvironmentOpenAIHostedSkillUnion() {
}
func (AutomationEnvironmentOpenAIHostedSkillInline) implAutomationEnvironmentOpenAIHostedSkillUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := AutomationEnvironmentOpenAIHostedSkillUnion.AsAny().(type) {
//	case openai.AutomationEnvironmentOpenAIHostedSkillSkillReference:
//	case openai.AutomationEnvironmentOpenAIHostedSkillInline:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u AutomationEnvironmentOpenAIHostedSkillUnion) AsAny() anyAutomationEnvironmentOpenAIHostedSkill {
	switch u.Type {
	case "skill_reference":
		return u.AsSkillReference()
	case "inline":
		return u.AsInline()
	}
	return nil
}

func (u AutomationEnvironmentOpenAIHostedSkillUnion) AsSkillReference() (v AutomationEnvironmentOpenAIHostedSkillSkillReference) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AutomationEnvironmentOpenAIHostedSkillUnion) AsInline() (v AutomationEnvironmentOpenAIHostedSkillInline) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u AutomationEnvironmentOpenAIHostedSkillUnion) RawJSON() string { return u.JSON.raw }

func (r *AutomationEnvironmentOpenAIHostedSkillUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A skill resolved afresh from the Skills API whenever a session starts.
type AutomationEnvironmentOpenAIHostedSkillSkillReference struct {
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
func (r AutomationEnvironmentOpenAIHostedSkillSkillReference) RawJSON() string { return r.JSON.raw }
func (r *AutomationEnvironmentOpenAIHostedSkillSkillReference) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Safe metadata for an inline skill archive.
type AutomationEnvironmentOpenAIHostedSkillInline struct {
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
func (r AutomationEnvironmentOpenAIHostedSkillInline) RawJSON() string { return r.JSON.raw }
func (r *AutomationEnvironmentOpenAIHostedSkillInline) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AutomationEnvironmentSelfHosted struct {
	// Directories that expose capabilities to the agent.
	CapabilityDirectories []string `json:"capability_directories" api:"required"`
	// An existing registration, or null to create a session-owned environment.
	EnvironmentID string `json:"environment_id" api:"required"`
	// The type of the object. Always `self_hosted`.
	Type constant.SelfHosted `json:"type" default:"self_hosted"`
	// The absolute project directory in the environment.
	WorkspaceDirectory string `json:"workspace_directory" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		CapabilityDirectories respjson.Field
		EnvironmentID         respjson.Field
		Type                  respjson.Field
		WorkspaceDirectory    respjson.Field
		ExtraFields           map[string]respjson.Field
		raw                   string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AutomationEnvironmentSelfHosted) RawJSON() string { return r.JSON.raw }
func (r *AutomationEnvironmentSelfHosted) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func AutomationScheduleParamOfParamOnce[
	T AutomationTargetParamNewSession | AutomationTargetParamExistingSession,
](name string, runAt int64, target T) AutomationScheduleParamUnion {
	var paramOnce AutomationScheduleParamOnce
	paramOnce.Name = name
	paramOnce.RunAt = runAt
	switch v := any(target).(type) {
	case AutomationTargetParamNewSession:
		paramOnce.Target.OfParamNewSession = &v
	case AutomationTargetParamExistingSession:
		paramOnce.Target.OfParamExistingSession = &v
	}
	return AutomationScheduleParamUnion{OfParamOnce: &paramOnce}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type AutomationScheduleParamUnion struct {
	OfParamPeriodic *AutomationScheduleParamPeriodic `json:",omitzero,inline"`
	OfParamOnce     *AutomationScheduleParamOnce     `json:",omitzero,inline"`
	paramUnion
}

func (u AutomationScheduleParamUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfParamPeriodic, u.OfParamOnce)
}
func (u *AutomationScheduleParamUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u AutomationScheduleParamUnion) GetExpression() *string {
	if vt := u.OfParamPeriodic; vt != nil {
		return &vt.Expression
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AutomationScheduleParamUnion) GetTimezone() *string {
	if vt := u.OfParamPeriodic; vt != nil {
		return &vt.Timezone
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AutomationScheduleParamUnion) GetRunAt() *int64 {
	if vt := u.OfParamOnce; vt != nil {
		return &vt.RunAt
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AutomationScheduleParamUnion) GetName() *string {
	if vt := u.OfParamPeriodic; vt != nil {
		return (*string)(&vt.Name)
	} else if vt := u.OfParamOnce; vt != nil {
		return (*string)(&vt.Name)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AutomationScheduleParamUnion) GetType() *string {
	if vt := u.OfParamPeriodic; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfParamOnce; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's Metadata property, if present.
func (u AutomationScheduleParamUnion) GetMetadata() map[string]string {
	if vt := u.OfParamPeriodic; vt != nil {
		return vt.Metadata
	} else if vt := u.OfParamOnce; vt != nil {
		return vt.Metadata
	}
	return nil
}

// Returns a pointer to the underlying variant's Target property, if present.
func (u AutomationScheduleParamUnion) GetTarget() *AutomationTargetParamUnion {
	if vt := u.OfParamPeriodic; vt != nil {
		return &vt.Target
	} else if vt := u.OfParamOnce; vt != nil {
		return &vt.Target
	}
	return nil
}

func init() {
	apijson.RegisterUnion[AutomationScheduleParamUnion](
		"type",
		apijson.Discriminator[AutomationScheduleParamPeriodic]("periodic"),
		apijson.Discriminator[AutomationScheduleParamOnce]("once"),
	)
}

// The properties Expression, Name, Target, Timezone, Type are required.
type AutomationScheduleParamPeriodic struct {
	// A standard five-field cron expression defining scheduled occurrences.
	Expression string `json:"expression" api:"required"`
	// A human-readable name for the automation.
	Name string `json:"name" api:"required"`
	// The session an automation invocation targets.
	Target AutomationTargetParamUnion `json:"target,omitzero" api:"required"`
	// The IANA timezone used to interpret the cron expression.
	Timezone string            `json:"timezone" api:"required"`
	Metadata map[string]string `json:"metadata,omitzero"`
	// The type of the object. Always `periodic`.
	//
	// This field can be elided, and will marshal its zero value as "periodic".
	Type constant.Periodic `json:"type" default:"periodic"`
	paramObj
}

func (r AutomationScheduleParamPeriodic) MarshalJSON() (data []byte, err error) {
	type shadow AutomationScheduleParamPeriodic
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AutomationScheduleParamPeriodic) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Name, RunAt, Target, Type are required.
type AutomationScheduleParamOnce struct {
	// A human-readable name for the automation.
	Name string `json:"name" api:"required"`
	// A future Unix timestamp in seconds.
	RunAt int64 `json:"run_at" api:"required"`
	// The session an automation invocation targets.
	Target   AutomationTargetParamUnion `json:"target,omitzero" api:"required"`
	Metadata map[string]string          `json:"metadata,omitzero"`
	// The type of the object. Always `once`.
	//
	// This field can be elided, and will marshal its zero value as "once".
	Type constant.Once `json:"type" default:"once"`
	paramObj
}

func (r AutomationScheduleParamOnce) MarshalJSON() (data []byte, err error) {
	type shadow AutomationScheduleParamOnce
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AutomationScheduleParamOnce) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// AutomationTargetUnion contains all possible properties and values from
// [AutomationTargetNewSession], [AutomationTargetExistingSession].
//
// Use the [AutomationTargetUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type AutomationTargetUnion struct {
	// This field is from variant [AutomationTargetNewSession].
	Agent AutomationAgentConfig `json:"agent"`
	// This field is from variant [AutomationTargetNewSession].
	AgentID string `json:"agent_id"`
	// This field is from variant [AutomationTargetNewSession].
	Environment AutomationEnvironmentUnion `json:"environment"`
	// This field is from variant [AutomationTargetNewSession].
	Input []AutomationTargetNewSessionInput `json:"input"`
	// This field is from variant [AutomationTargetNewSession].
	Metadata map[string]string `json:"metadata"`
	// Any of "new_session", "existing_session".
	Type string `json:"type"`
	// This field is from variant [AutomationTargetNewSession].
	VaultIDs []string `json:"vault_ids"`
	// This field is from variant [AutomationTargetExistingSession].
	Instructions string `json:"instructions"`
	// This field is from variant [AutomationTargetExistingSession].
	SessionID string `json:"session_id"`
	JSON      struct {
		Agent        respjson.Field
		AgentID      respjson.Field
		Environment  respjson.Field
		Input        respjson.Field
		Metadata     respjson.Field
		Type         respjson.Field
		VaultIDs     respjson.Field
		Instructions respjson.Field
		SessionID    respjson.Field
		raw          string
	} `json:"-"`
}

// anyAutomationTarget is implemented by each variant of [AutomationTargetUnion] to
// add type safety for the return type of [AutomationTargetUnion.AsAny]
type anyAutomationTarget interface {
	implAutomationTargetUnion()
}

func (AutomationTargetNewSession) implAutomationTargetUnion()      {}
func (AutomationTargetExistingSession) implAutomationTargetUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := AutomationTargetUnion.AsAny().(type) {
//	case openai.AutomationTargetNewSession:
//	case openai.AutomationTargetExistingSession:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u AutomationTargetUnion) AsAny() anyAutomationTarget {
	switch u.Type {
	case "new_session":
		return u.AsNewSession()
	case "existing_session":
		return u.AsExistingSession()
	}
	return nil
}

func (u AutomationTargetUnion) AsNewSession() (v AutomationTargetNewSession) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AutomationTargetUnion) AsExistingSession() (v AutomationTargetExistingSession) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u AutomationTargetUnion) RawJSON() string { return u.JSON.raw }

func (r *AutomationTargetUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AutomationTargetNewSession struct {
	// The inline agent configuration, or null when agent_id is configured.
	Agent AutomationAgentConfig `json:"agent" api:"required"`
	// The saved agent reference, or null when an inline agent is configured.
	AgentID string `json:"agent_id" api:"required"`
	// Configuration or references used to select the execution environment.
	Environment AutomationEnvironmentUnion `json:"environment" api:"required"`
	// Initial input delivered to each new session.
	Input []AutomationTargetNewSessionInput `json:"input" api:"required"`
	// Metadata attached to each new session.
	Metadata map[string]string `json:"metadata" api:"required"`
	// The type of the object. Always `new_session`.
	Type constant.NewSession `json:"type" default:"new_session"`
	// Vaults made available to each new session.
	VaultIDs []string `json:"vault_ids" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Agent       respjson.Field
		AgentID     respjson.Field
		Environment respjson.Field
		Input       respjson.Field
		Metadata    respjson.Field
		Type        respjson.Field
		VaultIDs    respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AutomationTargetNewSession) RawJSON() string { return r.JSON.raw }
func (r *AutomationTargetNewSession) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A user message recorded in a session.
type AutomationTargetNewSessionInput struct {
	// The content of the message.
	Content []InputContentUnion `json:"content" api:"required"`
	// The role of the message author. Always `user`.
	Role constant.User `json:"role" default:"user"`
	// The type of the input item. Always `message`.
	Type constant.Message `json:"type" default:"message"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Content     respjson.Field
		Role        respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AutomationTargetNewSessionInput) RawJSON() string { return r.JSON.raw }
func (r *AutomationTargetNewSessionInput) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AutomationTargetExistingSession struct {
	// Saved instructions delivered as an automation notification.
	Instructions string `json:"instructions" api:"required"`
	// The session receiving each invocation.
	SessionID string `json:"session_id" api:"required"`
	// The type of the object. Always `existing_session`.
	Type constant.ExistingSession `json:"type" default:"existing_session"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Instructions respjson.Field
		SessionID    respjson.Field
		Type         respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AutomationTargetExistingSession) RawJSON() string { return r.JSON.raw }
func (r *AutomationTargetExistingSession) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func AutomationTargetParamOfParamNewSession[
	T EnvironmentParamNone | EnvironmentParamOpenAIHosted | EnvironmentParamSelfHosted,
](environment T, input []AgentSessionInputMessageParam) AutomationTargetParamUnion {
	var paramNewSession AutomationTargetParamNewSession
	switch v := any(environment).(type) {
	case EnvironmentParamNone:
		paramNewSession.Environment.OfParamNone = &v
	case EnvironmentParamOpenAIHosted:
		paramNewSession.Environment.OfParamOpenAIHosted = &v
	case EnvironmentParamSelfHosted:
		paramNewSession.Environment.OfParamSelfHosted = &v
	}
	paramNewSession.Input = input
	return AutomationTargetParamUnion{OfParamNewSession: &paramNewSession}
}

func AutomationTargetParamOfParamExistingSession(instructions string, sessionID string) AutomationTargetParamUnion {
	var paramExistingSession AutomationTargetParamExistingSession
	paramExistingSession.Instructions = instructions
	paramExistingSession.SessionID = sessionID
	return AutomationTargetParamUnion{OfParamExistingSession: &paramExistingSession}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type AutomationTargetParamUnion struct {
	OfParamNewSession      *AutomationTargetParamNewSession      `json:",omitzero,inline"`
	OfParamExistingSession *AutomationTargetParamExistingSession `json:",omitzero,inline"`
	paramUnion
}

func (u AutomationTargetParamUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfParamNewSession, u.OfParamExistingSession)
}
func (u *AutomationTargetParamUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u AutomationTargetParamUnion) GetEnvironment() *EnvironmentParamUnion {
	if vt := u.OfParamNewSession; vt != nil {
		return &vt.Environment
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AutomationTargetParamUnion) GetInput() []AgentSessionInputMessageParam {
	if vt := u.OfParamNewSession; vt != nil {
		return vt.Input
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AutomationTargetParamUnion) GetAgent() *AutomationTargetParamNewSessionAgent {
	if vt := u.OfParamNewSession; vt != nil {
		return &vt.Agent
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AutomationTargetParamUnion) GetAgentID() *string {
	if vt := u.OfParamNewSession; vt != nil && vt.AgentID.Valid() {
		return &vt.AgentID.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AutomationTargetParamUnion) GetMetadata() map[string]string {
	if vt := u.OfParamNewSession; vt != nil {
		return vt.Metadata
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AutomationTargetParamUnion) GetVaultIDs() []string {
	if vt := u.OfParamNewSession; vt != nil {
		return vt.VaultIDs
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AutomationTargetParamUnion) GetInstructions() *string {
	if vt := u.OfParamExistingSession; vt != nil {
		return &vt.Instructions
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AutomationTargetParamUnion) GetSessionID() *string {
	if vt := u.OfParamExistingSession; vt != nil {
		return &vt.SessionID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AutomationTargetParamUnion) GetType() *string {
	if vt := u.OfParamNewSession; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfParamExistingSession; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[AutomationTargetParamUnion](
		"type",
		apijson.Discriminator[AutomationTargetParamNewSession]("new_session"),
		apijson.Discriminator[AutomationTargetParamExistingSession]("existing_session"),
	)
}

// Create a session for each invocation.
//
// The properties Environment, Input, Type are required.
type AutomationTargetParamNewSession struct {
	// The execution environment for each new session.
	Environment EnvironmentParamUnion `json:"environment,omitzero" api:"required"`
	// Initial user messages submitted with each new session.
	Input []AgentSessionInputMessageParam `json:"input,omitzero" api:"required"`
	// A saved agent ID. Provide exactly one of agent and agent_id.
	AgentID param.Opt[string] `json:"agent_id,omitzero"`
	// Metadata attached to each new session.
	Metadata map[string]string `json:"metadata,omitzero"`
	// The IDs of vaults made available to each new session.
	VaultIDs []string `json:"vault_ids,omitzero"`
	// Inline agent configuration. Provide exactly one of agent and agent_id.
	Agent AutomationTargetParamNewSessionAgent `json:"agent,omitzero"`
	// The type of the object. Always `new_session`.
	//
	// This field can be elided, and will marshal its zero value as "new_session".
	Type constant.NewSession `json:"type" default:"new_session"`
	paramObj
}

func (r AutomationTargetParamNewSession) MarshalJSON() (data []byte, err error) {
	type shadow AutomationTargetParamNewSession
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AutomationTargetParamNewSession) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Inline agent configuration. Provide exactly one of agent and agent_id.
type AutomationTargetParamNewSessionAgent struct {
	// Additional instructions appended to the agent's default base instructions. Omit
	// to leave unchanged.
	Instructions param.Opt[string] `json:"instructions,omitzero"`
	// The model to use for the agent. The requested model name is preserved.
	Model param.Opt[string] `json:"model,omitzero"`
	// Context compaction settings. Omit to keep the current settings. A supplied
	// object replaces the settings; null or an empty object enables background
	// compaction.
	Compaction AutomationTargetParamNewSessionAgentCompaction `json:"compaction,omitzero"`
	// The service tier used for model requests.
	//
	// Any of "auto", "default", "flex", "priority", "fast", "ultrafast".
	ServiceTier string `json:"service_tier,omitzero"`
	// Tools available to the agent. Omit to inherit, or pass null to clear them. The
	// resolved tool list must fit within 3 MiB (3,145,728 bytes) of compact UTF-8
	// JSON.
	Tools []AgentToolParamUnion `json:"tools,omitzero"`
	// Experimental settings to merge by key. Omitted settings keep their current
	// values.
	Experimental AutomationTargetParamNewSessionAgentExperimental `json:"experimental,omitzero"`
	// Configuration for creating and coordinating subagents.
	MultiAgent MultiAgentConfigParam `json:"multi_agent,omitzero"`
	// Configuration for model reasoning. Omit to keep the current settings; pass
	// `null` to reset to the model's default effort.
	Reasoning AgentReasoningParam `json:"reasoning,omitzero"`
	// Configuration for text generated by the agent.
	Text AgentTextParam `json:"text,omitzero"`
	paramObj
}

func (r AutomationTargetParamNewSessionAgent) MarshalJSON() (data []byte, err error) {
	type shadow AutomationTargetParamNewSessionAgent
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AutomationTargetParamNewSessionAgent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[AutomationTargetParamNewSessionAgent](
		"service_tier", "auto", "default", "flex", "priority", "fast", "ultrafast",
	)
}

// Context compaction settings. Omit to keep the current settings. A supplied
// object replaces the settings; null or an empty object enables background
// compaction.
type AutomationTargetParamNewSessionAgentCompaction struct {
	// Whether to compact context in the background. Omission or null defaults to true.
	// Set to false to use blocking compaction. Background compaction may increase
	// token usage.
	Async param.Opt[bool] `json:"async,omitzero"`
	paramObj
}

func (r AutomationTargetParamNewSessionAgentCompaction) MarshalJSON() (data []byte, err error) {
	type shadow AutomationTargetParamNewSessionAgentCompaction
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AutomationTargetParamNewSessionAgentCompaction) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Experimental settings to merge by key. Omitted settings keep their current
// values.
type AutomationTargetParamNewSessionAgentExperimental struct {
	// Whether to give the model periodic current-time reminders. Omit to keep the
	// current setting.
	TimeAwareness param.Opt[bool] `json:"time_awareness,omitzero"`
	paramObj
}

func (r AutomationTargetParamNewSessionAgentExperimental) MarshalJSON() (data []byte, err error) {
	type shadow AutomationTargetParamNewSessionAgentExperimental
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AutomationTargetParamNewSessionAgentExperimental) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Deliver instructions to an existing session.
//
// The properties Instructions, SessionID, Type are required.
type AutomationTargetParamExistingSession struct {
	// Instructions delivered as an automation notification.
	Instructions string `json:"instructions" api:"required"`
	// The session receiving each invocation.
	SessionID string `json:"session_id" api:"required"`
	// The type of the object. Always `existing_session`.
	//
	// This field can be elided, and will marshal its zero value as "existing_session".
	Type constant.ExistingSession `json:"type" default:"existing_session"`
	paramObj
}

func (r AutomationTargetParamExistingSession) MarshalJSON() (data []byte, err error) {
	type shadow AutomationTargetParamExistingSession
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AutomationTargetParamExistingSession) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Fields shared by periodic and one-time automations.
type OneTimeAutomation struct {
	// The automation ID.
	ID string `json:"id" api:"required"`
	// The creation time as a Unix timestamp in seconds.
	CreatedAt int64 `json:"created_at" api:"required"`
	// Service-recorded provenance; callers cannot supply this field.
	CreatedBy AutomationCreatedByUnion `json:"created_by" api:"required"`
	// Custom string key-value pairs attached to the automation.
	Metadata map[string]string `json:"metadata" api:"required"`
	// The automation's display name.
	Name string `json:"name" api:"required"`
	// Always `agent.automation`.
	Object constant.AgentAutomation `json:"object" default:"agent.automation"`
	// The Unix timestamp, in seconds, of the scheduled occurrence.
	RunAt int64 `json:"run_at" api:"required"`
	// Completion means the schedule has finished, not the agent's work.
	//
	// Any of "active", "paused", "completed".
	Status OneTimeAutomationStatus `json:"status" api:"required"`
	// The session configuration or existing session receiving invocations.
	Target AutomationTargetUnion `json:"target" api:"required"`
	// The type of the object. Always `once`.
	Type constant.Once `json:"type" default:"once"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		CreatedBy   respjson.Field
		Metadata    respjson.Field
		Name        respjson.Field
		Object      respjson.Field
		RunAt       respjson.Field
		Status      respjson.Field
		Target      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r OneTimeAutomation) RawJSON() string { return r.JSON.raw }
func (r *OneTimeAutomation) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Completion means the schedule has finished, not the agent's work.
type OneTimeAutomationStatus string

const (
	OneTimeAutomationStatusActive    OneTimeAutomationStatus = "active"
	OneTimeAutomationStatusPaused    OneTimeAutomationStatus = "paused"
	OneTimeAutomationStatusCompleted OneTimeAutomationStatus = "completed"
)

// Fields shared by periodic and one-time automations.
type PeriodicAutomation struct {
	// The automation ID.
	ID string `json:"id" api:"required"`
	// The creation time as a Unix timestamp in seconds.
	CreatedAt int64 `json:"created_at" api:"required"`
	// Service-recorded provenance; callers cannot supply this field.
	CreatedBy AutomationCreatedByUnion `json:"created_by" api:"required"`
	// The cron expression determining scheduled occurrences.
	Expression string `json:"expression" api:"required"`
	// Custom string key-value pairs attached to the automation.
	Metadata map[string]string `json:"metadata" api:"required"`
	// The automation's display name.
	Name string `json:"name" api:"required"`
	// Always `agent.automation`.
	Object constant.AgentAutomation `json:"object" default:"agent.automation"`
	// Whether scheduled delivery is active or paused.
	//
	// Any of "active", "paused".
	Status PeriodicAutomationStatus `json:"status" api:"required"`
	// The session configuration or existing session receiving invocations.
	Target AutomationTargetUnion `json:"target" api:"required"`
	// The timezone used to interpret the expression.
	Timezone string `json:"timezone" api:"required"`
	// The type of the object. Always `periodic`.
	Type constant.Periodic `json:"type" default:"periodic"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		CreatedBy   respjson.Field
		Expression  respjson.Field
		Metadata    respjson.Field
		Name        respjson.Field
		Object      respjson.Field
		Status      respjson.Field
		Target      respjson.Field
		Timezone    respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r PeriodicAutomation) RawJSON() string { return r.JSON.raw }
func (r *PeriodicAutomation) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Whether scheduled delivery is active or paused.
type PeriodicAutomationStatus string

const (
	PeriodicAutomationStatusActive PeriodicAutomationStatus = "active"
	PeriodicAutomationStatusPaused PeriodicAutomationStatus = "paused"
)

type BetaAgentAutomationNewParams struct {
	// Creates an automation without invoking it immediately.
	CreateAutomationParams AutomationScheduleParamUnion
	paramObj
}

func (r BetaAgentAutomationNewParams) MarshalJSON() (data []byte, err error) {
	return param.MarshalObject(r, r.CreateAutomationParams)
}
func (r *BetaAgentAutomationNewParams) UnmarshalJSON(data []byte) error {
	if r == nil {
		return apijson.UnmarshalRoot(data, r)
	}
	return apijson.UnmarshalWrappedObject(data, r, &r.CreateAutomationParams)
}

type BetaAgentAutomationUpdateParams struct {
	// Only accepted for periodic automations.
	Expression param.Opt[string] `json:"expression,omitzero"`
	Name       param.Opt[string] `json:"name,omitzero"`
	// A future Unix timestamp in seconds; accepted only before a one-time scheduled
	// invocation exists.
	RunAt param.Opt[int64] `json:"run_at,omitzero"`
	// Only accepted for periodic automations.
	Timezone param.Opt[string] `json:"timezone,omitzero"`
	// Replaces all metadata; omit to leave unchanged.
	Metadata map[string]string `json:"metadata,omitzero"`
	// Replaces the complete target, retaining its type and any creator session.
	Target AutomationTargetParamUnion `json:"target,omitzero"`
	paramObj
}

func (r BetaAgentAutomationUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaAgentAutomationUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaAgentAutomationUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaAgentAutomationListParams struct {
	// Return resources after this resource ID in the selected order.
	After param.Opt[string] `query:"after,omitzero" json:"-"`
	// The maximum number of resources to return, between 1 and 100. Defaults to 20.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Matches existing-session targets, including agent-created automations.
	SessionID param.Opt[string] `query:"session_id,omitzero" json:"-"`
	// The order in which resources are returned. Defaults to `desc`.
	//
	// Any of "asc", "desc".
	Order BetaAgentAutomationListParamsOrder `query:"order,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaAgentAutomationListParams]'s query parameters as
// `url.Values`.
func (r BetaAgentAutomationListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// The order in which resources are returned. Defaults to `desc`.
type BetaAgentAutomationListParamsOrder string

const (
	BetaAgentAutomationListParamsOrderAsc  BetaAgentAutomationListParamsOrder = "asc"
	BetaAgentAutomationListParamsOrderDesc BetaAgentAutomationListParamsOrder = "desc"
)

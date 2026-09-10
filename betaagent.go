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

// BetaAgentService contains methods and other services that help with interacting
// with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaAgentService] method instead.
type BetaAgentService struct {
	Options      []option.RequestOption
	Environments BetaAgentEnvironmentService
	Vaults       BetaAgentVaultService
	Sessions     BetaAgentSessionService
}

// NewBetaAgentService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewBetaAgentService(opts ...option.RequestOption) (r BetaAgentService) {
	r = BetaAgentService{}
	r.Options = requestconfig.InheritedOptions(opts...)
	r.Environments = NewBetaAgentEnvironmentService(opts...)
	r.Vaults = NewBetaAgentVaultService(opts...)
	r.Sessions = NewBetaAgentSessionService(opts...)
	return
}

// Creates a reusable agent without storing credentials. See
// [agent configuration](https://developers.openai.com/api/docs/guides/agents-api/configuration).
func (r *BetaAgentService) New(ctx context.Context, body BetaAgentNewParams, opts ...option.RequestOption) (res *Agent, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "agents=v1")}, opts...)
	path := "agents"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Retrieves a reusable agent by ID. See
// [agent configuration](https://developers.openai.com/api/docs/guides/agents-api/configuration).
func (r *BetaAgentService) Get(ctx context.Context, agentID string, opts ...option.RequestOption) (res *Agent, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "agents=v1")}, opts...)
	if agentID == "" {
		err = errors.New("missing required agent_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("agents/%s", agentID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Updates a reusable agent. See
// [agent configuration](https://developers.openai.com/api/docs/guides/agents-api/configuration).
func (r *BetaAgentService) Update(ctx context.Context, agentID string, body BetaAgentUpdateParams, opts ...option.RequestOption) (res *Agent, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "agents=v1")}, opts...)
	if agentID == "" {
		err = errors.New("missing required agent_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("agents/%s", agentID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Lists reusable agents in the current project. See
// [agent configuration](https://developers.openai.com/api/docs/guides/agents-api/configuration).
func (r *BetaAgentService) List(ctx context.Context, query BetaAgentListParams, opts ...option.RequestOption) (res *pagination.CursorPage[Agent], err error) {
	var raw *http.Response
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "agents=v1"), option.WithResponseInto(&raw)}, opts...)
	path := "agents"
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

// Lists reusable agents in the current project. See
// [agent configuration](https://developers.openai.com/api/docs/guides/agents-api/configuration).
func (r *BetaAgentService) ListAutoPaging(ctx context.Context, query BetaAgentListParams, opts ...option.RequestOption) *pagination.CursorPageAutoPager[Agent] {
	return pagination.NewCursorPageAutoPager(r.List(ctx, query, opts...))
}

// Deletes a reusable agent. See
// [agent configuration](https://developers.openai.com/api/docs/guides/agents-api/configuration).
func (r *BetaAgentService) Delete(ctx context.Context, agentID string, opts ...option.RequestOption) (res *AgentDeleted, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "agents=v1")}, opts...)
	if agentID == "" {
		err = errors.New("missing required agent_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("agents/%s", agentID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// A reusable agent scoped to the caller's project.
type Agent struct {
	// The ID of the reusable agent.
	ID string `json:"id" api:"required"`
	// The Unix timestamp, in seconds, when the agent was created.
	CreatedAt int64 `json:"created_at" api:"required"`
	// Custom instructions appended to the agent's default base instructions.
	Instructions string `json:"instructions" api:"required"`
	// Custom string key-value pairs attached to the agent.
	Metadata map[string]string `json:"metadata" api:"required"`
	// The requested model name used for inference.
	Model string `json:"model" api:"required"`
	// The resolved configuration for creating and coordinating subagents.
	MultiAgent MultiAgentConfig `json:"multi_agent" api:"required"`
	// A human-readable name for the agent, or null if it is unnamed.
	Name string `json:"name" api:"required"`
	// The object type. Always `agent`.
	Object constant.Agent `json:"object" default:"agent"`
	// The resolved reasoning configuration, including the model default for an omitted
	// effort.
	Reasoning AgentReasoning `json:"reasoning" api:"required"`
	// The resolved service-tier policy used for model requests.
	//
	// Any of "auto", "default", "flex", "priority", "fast".
	ServiceTier AgentServiceTier `json:"service_tier" api:"required"`
	// The resolved configuration for text generated by the agent.
	Text AgentText `json:"text" api:"required"`
	// Tools available to the agent.
	Tools []PersistedAgentToolUnion `json:"tools" api:"required"`
	// The Unix timestamp, in seconds, when the agent was last updated.
	UpdatedAt int64 `json:"updated_at" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID           respjson.Field
		CreatedAt    respjson.Field
		Instructions respjson.Field
		Metadata     respjson.Field
		Model        respjson.Field
		MultiAgent   respjson.Field
		Name         respjson.Field
		Object       respjson.Field
		Reasoning    respjson.Field
		ServiceTier  respjson.Field
		Text         respjson.Field
		Tools        respjson.Field
		UpdatedAt    respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r Agent) RawJSON() string { return r.JSON.raw }
func (r *Agent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The resolved service-tier policy used for model requests.
type AgentServiceTier string

const (
	AgentServiceTierAuto     AgentServiceTier = "auto"
	AgentServiceTierDefault  AgentServiceTier = "default"
	AgentServiceTierFlex     AgentServiceTier = "flex"
	AgentServiceTierPriority AgentServiceTier = "priority"
	AgentServiceTierFast     AgentServiceTier = "fast"
)

// A request to close a subagent.
type AgentCloseSubagentCallItem struct {
	// The ID of the tool call item.
	ID string `json:"id" api:"required"`
	// The ID of the agent to close.
	RecipientAgentID string `json:"recipient_agent_id" api:"required"`
	// The ID of the agent requesting the close.
	SenderAgentID string `json:"sender_agent_id" api:"required"`
	// The status of the tool call.
	//
	// Any of "in_progress", "completed", "failed", "incomplete".
	Status AgentFunctionCallStatus `json:"status" api:"required"`
	// The ID of the turn that contains this item.
	TurnID string `json:"turn_id" api:"required"`
	// The item type. Always `close_subagent_call`.
	Type constant.CloseSubagentCall `json:"type" default:"close_subagent_call"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID               respjson.Field
		RecipientAgentID respjson.Field
		SenderAgentID    respjson.Field
		Status           respjson.Field
		TurnID           respjson.Field
		Type             respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentCloseSubagentCallItem) RawJSON() string { return r.JSON.raw }
func (r *AgentCloseSubagentCallItem) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A command execution produced by the agent.
type AgentCommandExecutionItem struct {
	// The ID of the command execution item.
	ID string `json:"id" api:"required"`
	// The command that was executed.
	Command string `json:"command" api:"required"`
	// The working directory used to execute the command.
	Cwd string `json:"cwd" api:"required"`
	// The command duration in milliseconds.
	DurationMs int64 `json:"duration_ms" api:"required"`
	// The process exit code, if the command completed.
	ExitCode int64 `json:"exit_code" api:"required"`
	// The command output, if available.
	Output string `json:"output" api:"required"`
	// The status of the command execution.
	//
	// Any of "in_progress", "completed", "failed", "incomplete".
	Status AgentFunctionCallStatus `json:"status" api:"required"`
	// The ID of the turn that contains this item.
	TurnID string `json:"turn_id" api:"required"`
	// The item type. Always `command_execution`.
	Type constant.CommandExecution `json:"type" default:"command_execution"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Command     respjson.Field
		Cwd         respjson.Field
		DurationMs  respjson.Field
		ExitCode    respjson.Field
		Output      respjson.Field
		Status      respjson.Field
		TurnID      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentCommandExecutionItem) RawJSON() string { return r.JSON.raw }
func (r *AgentCommandExecutionItem) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// AgentContentUnion contains all possible properties and values from [OutputText],
// [AgentContentEncryptedContent].
//
// Use the [AgentContentUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type AgentContentUnion struct {
	// This field is from variant [OutputText].
	Text string `json:"text"`
	// Any of "output_text", "encrypted_content".
	Type string `json:"type"`
	// This field is from variant [AgentContentEncryptedContent].
	EncryptedContent string `json:"encrypted_content"`
	JSON             struct {
		Text             respjson.Field
		Type             respjson.Field
		EncryptedContent respjson.Field
		raw              string
	} `json:"-"`
}

// anyAgentContent is implemented by each variant of [AgentContentUnion] to add
// type safety for the return type of [AgentContentUnion.AsAny]
type anyAgentContent interface {
	implAgentContentUnion()
}

func (OutputText) implAgentContentUnion()                   {}
func (AgentContentEncryptedContent) implAgentContentUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := AgentContentUnion.AsAny().(type) {
//	case openai.OutputText:
//	case openai.AgentContentEncryptedContent:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u AgentContentUnion) AsAny() anyAgentContent {
	switch u.Type {
	case "output_text":
		return u.AsOutputText()
	case "encrypted_content":
		return u.AsEncryptedContent()
	}
	return nil
}

func (u AgentContentUnion) AsOutputText() (v OutputText) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentContentUnion) AsEncryptedContent() (v AgentContentEncryptedContent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u AgentContentUnion) RawJSON() string { return u.JSON.raw }

func (r *AgentContentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Encrypted content exchanged between agents.
type AgentContentEncryptedContent struct {
	// The encrypted content payload.
	EncryptedContent string `json:"encrypted_content" api:"required"`
	// The content type. Always `encrypted_content`.
	Type constant.EncryptedContent `json:"type" default:"encrypted_content"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EncryptedContent respjson.Field
		Type             respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentContentEncryptedContent) RawJSON() string { return r.JSON.raw }
func (r *AgentContentEncryptedContent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A request to spawn a subagent.
type AgentCreateSubagentCallItem struct {
	// The ID of the tool call item.
	ID string `json:"id" api:"required"`
	// The ID of the agent that requested the subagent.
	AgentID string `json:"agent_id" api:"required"`
	// The task given to the spawned agent.
	Content []AgentContentUnion `json:"content" api:"required"`
	// The model requested for the spawned agent.
	Model string `json:"model" api:"required"`
	// The reasoning effort requested for the spawned agent.
	ReasoningEffort string `json:"reasoning_effort" api:"required"`
	// The status of the tool call.
	//
	// Any of "in_progress", "completed", "failed", "incomplete".
	Status AgentFunctionCallStatus `json:"status" api:"required"`
	// The ID of the turn that contains this item.
	TurnID string `json:"turn_id" api:"required"`
	// The item type. Always `create_subagent_call`.
	Type constant.CreateSubagentCall `json:"type" default:"create_subagent_call"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID              respjson.Field
		AgentID         respjson.Field
		Content         respjson.Field
		Model           respjson.Field
		ReasoningEffort respjson.Field
		Status          respjson.Field
		TurnID          respjson.Field
		Type            respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentCreateSubagentCallItem) RawJSON() string { return r.JSON.raw }
func (r *AgentCreateSubagentCallItem) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A deleted reusable agent.
type AgentDeleted struct {
	// The ID of the deleted agent.
	ID string `json:"id" api:"required"`
	// Whether the agent was deleted. Always `true`.
	Deleted bool `json:"deleted" api:"required"`
	// The object type. Always `agent.deleted`.
	Object constant.AgentDeleted `json:"object" default:"agent.deleted"`
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
func (r AgentDeleted) RawJSON() string { return r.JSON.raw }
func (r *AgentDeleted) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A function call produced by the agent.
type AgentFunctionCallItem struct {
	// The ID of the function call item.
	ID string `json:"id" api:"required"`
	// The arguments to pass to the function.
	Arguments any `json:"arguments" api:"required"`
	// The ID used to submit the function result.
	CallID string `json:"call_id" api:"required"`
	// The name of the function to call.
	Name string `json:"name" api:"required"`
	// The status of the function call.
	//
	// Any of "in_progress", "completed", "failed", "incomplete".
	Status AgentFunctionCallStatus `json:"status" api:"required"`
	// The ID of the turn that contains this item.
	TurnID string `json:"turn_id" api:"required"`
	// The item type. Always `function_call`.
	Type constant.FunctionCall `json:"type" default:"function_call"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Arguments   respjson.Field
		CallID      respjson.Field
		Name        respjson.Field
		Status      respjson.Field
		TurnID      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentFunctionCallItem) RawJSON() string { return r.JSON.raw }
func (r *AgentFunctionCallItem) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// AgentFunctionCallOutputUnion contains all possible properties and values from
// [string], [[]InputContentUnion].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfString OfAgentFunctionCallOutputArray]
type AgentFunctionCallOutputUnion struct {
	// This field will be present if the value is a [string] instead of an object.
	OfString string `json:",inline"`
	// This field will be present if the value is a [[]InputContentUnion] instead of an
	// object.
	OfAgentFunctionCallOutputArray []InputContentUnion `json:",inline"`
	JSON                           struct {
		OfString                       respjson.Field
		OfAgentFunctionCallOutputArray respjson.Field
		raw                            string
	} `json:"-"`
}

func (u AgentFunctionCallOutputUnion) AsString() (v string) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentFunctionCallOutputUnion) AsAgentFunctionCallOutputArray() (v []InputContentUnion) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u AgentFunctionCallOutputUnion) RawJSON() string { return u.JSON.raw }

func (r *AgentFunctionCallOutputUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type AgentFunctionCallOutputParamUnion struct {
	OfString               param.Opt[string]        `json:",omitzero,inline"`
	OfArrayOfInputContents []InputContentParamUnion `json:",omitzero,inline"`
	paramUnion
}

func (u AgentFunctionCallOutputParamUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfString, u.OfArrayOfInputContents)
}
func (u *AgentFunctionCallOutputParamUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// The status of a tool call.
type AgentFunctionCallStatus string

const (
	AgentFunctionCallStatusInProgress AgentFunctionCallStatus = "in_progress"
	AgentFunctionCallStatusCompleted  AgentFunctionCallStatus = "completed"
	AgentFunctionCallStatusFailed     AgentFunctionCallStatus = "failed"
	AgentFunctionCallStatusIncomplete AgentFunctionCallStatus = "incomplete"
)

// A request to interrupt a subagent's current turn. The subagent remains
// available.
type AgentInterruptSubagentCallItem struct {
	// The ID of the tool call item.
	ID string `json:"id" api:"required"`
	// The ID of the agent to interrupt.
	RecipientAgentID string `json:"recipient_agent_id" api:"required"`
	// The ID of the agent requesting the interrupt.
	SenderAgentID string `json:"sender_agent_id" api:"required"`
	// The status of the tool call.
	//
	// Any of "in_progress", "completed", "failed", "incomplete".
	Status AgentFunctionCallStatus `json:"status" api:"required"`
	// The ID of the turn that contains this item.
	TurnID string `json:"turn_id" api:"required"`
	// The item type. Always `interrupt_subagent_call`.
	Type constant.InterruptSubagentCall `json:"type" default:"interrupt_subagent_call"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID               respjson.Field
		RecipientAgentID respjson.Field
		SenderAgentID    respjson.Field
		Status           respjson.Field
		TurnID           respjson.Field
		Type             respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentInterruptSubagentCallItem) RawJSON() string { return r.JSON.raw }
func (r *AgentInterruptSubagentCallItem) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A call to a tool on an MCP server.
type AgentMcpCallItem struct {
	// The ID of the MCP call item.
	ID string `json:"id" api:"required"`
	// The arguments passed to the MCP tool.
	Arguments any `json:"arguments" api:"required"`
	// The error returned by the MCP tool, if any.
	Error any `json:"error" api:"required"`
	// The name of the MCP tool.
	Name string `json:"name" api:"required"`
	// The output returned by the MCP tool, if any.
	Output any `json:"output" api:"required"`
	// The label of the MCP server.
	ServerLabel string `json:"server_label" api:"required"`
	// The status of the MCP tool call.
	//
	// Any of "in_progress", "completed", "failed", "incomplete".
	Status AgentFunctionCallStatus `json:"status" api:"required"`
	// The ID of the turn that contains this item.
	TurnID string `json:"turn_id" api:"required"`
	// The item type. Always `mcp_call`.
	Type constant.McpCall `json:"type" default:"mcp_call"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Arguments   respjson.Field
		Error       respjson.Field
		Name        respjson.Field
		Output      respjson.Field
		ServerLabel respjson.Field
		Status      respjson.Field
		TurnID      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentMcpCallItem) RawJSON() string { return r.JSON.raw }
func (r *AgentMcpCallItem) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Emitted when command execution produces an output delta.
type AgentOutputCommandExecutionOutputDeltaEvent struct {
	// The output text that was appended.
	Delta string `json:"delta" api:"required"`
	// The unique ID of the event.
	EventID string `json:"event_id" api:"required"`
	// The ID of the command execution item.
	ItemID string `json:"item_id" api:"required"`
	// The index of the item in the turn output.
	OutputIndex int64 `json:"output_index" api:"required"`
	// The ID of the session associated with the event.
	SessionID string `json:"session_id" api:"required"`
	// The ID of the turn associated with the event, when applicable.
	TurnID string `json:"turn_id" api:"required"`
	// The type of the object. Always `agent.output.command_execution_output.delta`.
	Type constant.AgentOutputCommandExecutionOutputDelta `json:"type" default:"agent.output.command_execution_output.delta"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Delta       respjson.Field
		EventID     respjson.Field
		ItemID      respjson.Field
		OutputIndex respjson.Field
		SessionID   respjson.Field
		TurnID      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentOutputCommandExecutionOutputDeltaEvent) RawJSON() string { return r.JSON.raw }
func (r *AgentOutputCommandExecutionOutputDeltaEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// AgentOutputItemUnion contains all possible properties and values from
// [AgentSessionAssistantMessage], [AgentReasoningItem], [AgentFunctionCallItem],
// [AgentMcpCallItem], [AgentWebSearchCallItem], [AgentCommandExecutionItem],
// [AgentCreateSubagentCallItem], [AgentSendSubagentInputCallItem],
// [AgentResumeSubagentCallItem], [AgentWaitForSubagentsCallItem],
// [AgentInterruptSubagentCallItem], [AgentCloseSubagentCallItem].
//
// Use the [AgentOutputItemUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type AgentOutputItemUnion struct {
	ID string `json:"id"`
	// This field is a union of [[]OutputText], [[]AgentContentUnion],
	// [[]AgentContentUnion]
	Content AgentOutputItemUnionContent `json:"content"`
	// This field is from variant [AgentSessionAssistantMessage].
	Phase AgentSessionAssistantMessagePhase `json:"phase"`
	// This field is from variant [AgentSessionAssistantMessage].
	Role   constant.Assistant `json:"role"`
	Status string             `json:"status"`
	TurnID string             `json:"turn_id"`
	// Any of "message", "reasoning", "function_call", "mcp_call", "web_search_call",
	// "command_execution", "create_subagent_call", "send_subagent_input_call",
	// "resume_subagent_call", "wait_for_subagents_call", "interrupt_subagent_call",
	// "close_subagent_call".
	Type string `json:"type"`
	// This field is from variant [AgentReasoningItem].
	Summary   []SummaryText `json:"summary"`
	Arguments any           `json:"arguments"`
	// This field is from variant [AgentFunctionCallItem].
	CallID string `json:"call_id"`
	Name   string `json:"name"`
	// This field is from variant [AgentMcpCallItem].
	Error any `json:"error"`
	// This field is a union of [any], [string]
	Output AgentOutputItemUnionOutput `json:"output"`
	// This field is from variant [AgentMcpCallItem].
	ServerLabel string `json:"server_label"`
	// This field is from variant [AgentWebSearchCallItem].
	Action WebSearchActionUnion `json:"action"`
	// This field is from variant [AgentCommandExecutionItem].
	Command string `json:"command"`
	// This field is from variant [AgentCommandExecutionItem].
	Cwd string `json:"cwd"`
	// This field is from variant [AgentCommandExecutionItem].
	DurationMs int64 `json:"duration_ms"`
	// This field is from variant [AgentCommandExecutionItem].
	ExitCode int64 `json:"exit_code"`
	// This field is from variant [AgentCreateSubagentCallItem].
	AgentID string `json:"agent_id"`
	// This field is from variant [AgentCreateSubagentCallItem].
	Model string `json:"model"`
	// This field is from variant [AgentCreateSubagentCallItem].
	ReasoningEffort  string `json:"reasoning_effort"`
	RecipientAgentID string `json:"recipient_agent_id"`
	SenderAgentID    string `json:"sender_agent_id"`
	// This field is from variant [AgentWaitForSubagentsCallItem].
	RecipientAgentIDs []string `json:"recipient_agent_ids"`
	JSON              struct {
		ID                respjson.Field
		Content           respjson.Field
		Phase             respjson.Field
		Role              respjson.Field
		Status            respjson.Field
		TurnID            respjson.Field
		Type              respjson.Field
		Summary           respjson.Field
		Arguments         respjson.Field
		CallID            respjson.Field
		Name              respjson.Field
		Error             respjson.Field
		Output            respjson.Field
		ServerLabel       respjson.Field
		Action            respjson.Field
		Command           respjson.Field
		Cwd               respjson.Field
		DurationMs        respjson.Field
		ExitCode          respjson.Field
		AgentID           respjson.Field
		Model             respjson.Field
		ReasoningEffort   respjson.Field
		RecipientAgentID  respjson.Field
		SenderAgentID     respjson.Field
		RecipientAgentIDs respjson.Field
		raw               string
	} `json:"-"`
}

// anyAgentOutputItem is implemented by each variant of [AgentOutputItemUnion] to
// add type safety for the return type of [AgentOutputItemUnion.AsAny]
type anyAgentOutputItem interface {
	implAgentOutputItemUnion()
}

func (AgentSessionAssistantMessage) implAgentOutputItemUnion()   {}
func (AgentReasoningItem) implAgentOutputItemUnion()             {}
func (AgentFunctionCallItem) implAgentOutputItemUnion()          {}
func (AgentMcpCallItem) implAgentOutputItemUnion()               {}
func (AgentWebSearchCallItem) implAgentOutputItemUnion()         {}
func (AgentCommandExecutionItem) implAgentOutputItemUnion()      {}
func (AgentCreateSubagentCallItem) implAgentOutputItemUnion()    {}
func (AgentSendSubagentInputCallItem) implAgentOutputItemUnion() {}
func (AgentResumeSubagentCallItem) implAgentOutputItemUnion()    {}
func (AgentWaitForSubagentsCallItem) implAgentOutputItemUnion()  {}
func (AgentInterruptSubagentCallItem) implAgentOutputItemUnion() {}
func (AgentCloseSubagentCallItem) implAgentOutputItemUnion()     {}

// Use the following switch statement to find the correct variant
//
//	switch variant := AgentOutputItemUnion.AsAny().(type) {
//	case openai.AgentSessionAssistantMessage:
//	case openai.AgentReasoningItem:
//	case openai.AgentFunctionCallItem:
//	case openai.AgentMcpCallItem:
//	case openai.AgentWebSearchCallItem:
//	case openai.AgentCommandExecutionItem:
//	case openai.AgentCreateSubagentCallItem:
//	case openai.AgentSendSubagentInputCallItem:
//	case openai.AgentResumeSubagentCallItem:
//	case openai.AgentWaitForSubagentsCallItem:
//	case openai.AgentInterruptSubagentCallItem:
//	case openai.AgentCloseSubagentCallItem:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u AgentOutputItemUnion) AsAny() anyAgentOutputItem {
	switch u.Type {
	case "message":
		return u.AsMessage()
	case "reasoning":
		return u.AsReasoning()
	case "function_call":
		return u.AsFunctionCall()
	case "mcp_call":
		return u.AsMcpCall()
	case "web_search_call":
		return u.AsWebSearchCall()
	case "command_execution":
		return u.AsCommandExecution()
	case "create_subagent_call":
		return u.AsCreateSubagentCall()
	case "send_subagent_input_call":
		return u.AsSendSubagentInputCall()
	case "resume_subagent_call":
		return u.AsResumeSubagentCall()
	case "wait_for_subagents_call":
		return u.AsWaitForSubagentsCall()
	case "interrupt_subagent_call":
		return u.AsInterruptSubagentCall()
	case "close_subagent_call":
		return u.AsCloseSubagentCall()
	}
	return nil
}

func (u AgentOutputItemUnion) AsMessage() (v AgentSessionAssistantMessage) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentOutputItemUnion) AsReasoning() (v AgentReasoningItem) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentOutputItemUnion) AsFunctionCall() (v AgentFunctionCallItem) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentOutputItemUnion) AsMcpCall() (v AgentMcpCallItem) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentOutputItemUnion) AsWebSearchCall() (v AgentWebSearchCallItem) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentOutputItemUnion) AsCommandExecution() (v AgentCommandExecutionItem) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentOutputItemUnion) AsCreateSubagentCall() (v AgentCreateSubagentCallItem) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentOutputItemUnion) AsSendSubagentInputCall() (v AgentSendSubagentInputCallItem) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentOutputItemUnion) AsResumeSubagentCall() (v AgentResumeSubagentCallItem) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentOutputItemUnion) AsWaitForSubagentsCall() (v AgentWaitForSubagentsCallItem) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentOutputItemUnion) AsInterruptSubagentCall() (v AgentInterruptSubagentCallItem) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentOutputItemUnion) AsCloseSubagentCall() (v AgentCloseSubagentCallItem) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u AgentOutputItemUnion) RawJSON() string { return u.JSON.raw }

func (r *AgentOutputItemUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// AgentOutputItemUnionContent is an implicit subunion of [AgentOutputItemUnion].
// AgentOutputItemUnionContent provides convenient access to the sub-properties of
// the union.
//
// For type safety it is recommended to directly use a variant of the
// [AgentOutputItemUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfOutputTextArray OfAgentContentArray]
type AgentOutputItemUnionContent struct {
	// This field will be present if the value is a [[]OutputText] instead of an
	// object.
	OfOutputTextArray []OutputText `json:",inline"`
	// This field will be present if the value is a [[]AgentContentUnion] instead of an
	// object.
	OfAgentContentArray []AgentContentUnion `json:",inline"`
	JSON                struct {
		OfOutputTextArray   respjson.Field
		OfAgentContentArray respjson.Field
		raw                 string
	} `json:"-"`
}

func (r *AgentOutputItemUnionContent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// AgentOutputItemUnionOutput is an implicit subunion of [AgentOutputItemUnion].
// AgentOutputItemUnionOutput provides convenient access to the sub-properties of
// the union.
//
// For type safety it is recommended to directly use a variant of the
// [AgentOutputItemUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfAgentMcpCallItemOutput OfString]
type AgentOutputItemUnionOutput struct {
	// This field will be present if the value is a [any] instead of an object.
	OfAgentMcpCallItemOutput any `json:",inline"`
	// This field will be present if the value is a [string] instead of an object.
	OfString string `json:",inline"`
	JSON     struct {
		OfAgentMcpCallItemOutput respjson.Field
		OfString                 respjson.Field
		raw                      string
	} `json:"-"`
}

func (r *AgentOutputItemUnionOutput) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The status of an agent output item.
type AgentOutputItemStatus string

const (
	AgentOutputItemStatusInProgress AgentOutputItemStatus = "in_progress"
	AgentOutputItemStatusCompleted  AgentOutputItemStatus = "completed"
	AgentOutputItemStatusIncomplete AgentOutputItemStatus = "incomplete"
)

// The reasoning configuration used by an agent.
type AgentReasoning struct {
	// The amount of reasoning effort used by an agent.
	//
	// Any of "none", "minimal", "low", "medium", "high", "xhigh", "max".
	Effort AgentReasoningEffort `json:"effort" api:"required"`
	// The reasoning summary format requested from an agent.
	//
	// Any of "concise", "detailed", "auto".
	Summary AgentReasoningSummary `json:"summary" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Effort      respjson.Field
		Summary     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentReasoning) RawJSON() string { return r.JSON.raw }
func (r *AgentReasoning) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The amount of reasoning effort used by an agent.
type AgentReasoningEffort string

const (
	AgentReasoningEffortNone    AgentReasoningEffort = "none"
	AgentReasoningEffortMinimal AgentReasoningEffort = "minimal"
	AgentReasoningEffortLow     AgentReasoningEffort = "low"
	AgentReasoningEffortMedium  AgentReasoningEffort = "medium"
	AgentReasoningEffortHigh    AgentReasoningEffort = "high"
	AgentReasoningEffortXhigh   AgentReasoningEffort = "xhigh"
	AgentReasoningEffortMax     AgentReasoningEffort = "max"
)

// The reasoning summary format requested from an agent.
type AgentReasoningSummary string

const (
	AgentReasoningSummaryConcise  AgentReasoningSummary = "concise"
	AgentReasoningSummaryDetailed AgentReasoningSummary = "detailed"
	AgentReasoningSummaryAuto     AgentReasoningSummary = "auto"
)

// A reasoning item produced by the agent.
type AgentReasoningItem struct {
	// The ID of the reasoning item.
	ID string `json:"id" api:"required"`
	// The status of an agent output item.
	//
	// Any of "in_progress", "completed", "incomplete".
	Status AgentOutputItemStatus `json:"status" api:"required"`
	// The reasoning summaries produced by the agent.
	Summary []SummaryText `json:"summary" api:"required"`
	// The ID of the turn that contains this item.
	TurnID string `json:"turn_id" api:"required"`
	// The item type. Always `reasoning`.
	Type constant.Reasoning `json:"type" default:"reasoning"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Status      respjson.Field
		Summary     respjson.Field
		TurnID      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentReasoningItem) RawJSON() string { return r.JSON.raw }
func (r *AgentReasoningItem) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Reasoning configuration for the agent.
type AgentReasoningParam struct {
	// The amount of reasoning effort the model should use.
	//
	// Any of "none", "minimal", "low", "medium", "high", "xhigh", "max".
	Effort AgentReasoningParamEffort `json:"effort,omitzero"`
	// The reasoning summary format requested from the model.
	//
	// Any of "concise", "detailed", "auto".
	Summary AgentReasoningParamSummary `json:"summary,omitzero"`
	paramObj
}

func (r AgentReasoningParam) MarshalJSON() (data []byte, err error) {
	type shadow AgentReasoningParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AgentReasoningParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The amount of reasoning effort the model should use.
type AgentReasoningParamEffort string

const (
	AgentReasoningParamEffortNone    AgentReasoningParamEffort = "none"
	AgentReasoningParamEffortMinimal AgentReasoningParamEffort = "minimal"
	AgentReasoningParamEffortLow     AgentReasoningParamEffort = "low"
	AgentReasoningParamEffortMedium  AgentReasoningParamEffort = "medium"
	AgentReasoningParamEffortHigh    AgentReasoningParamEffort = "high"
	AgentReasoningParamEffortXhigh   AgentReasoningParamEffort = "xhigh"
	AgentReasoningParamEffortMax     AgentReasoningParamEffort = "max"
)

// The reasoning summary format requested from the model.
type AgentReasoningParamSummary string

const (
	AgentReasoningParamSummaryConcise  AgentReasoningParamSummary = "concise"
	AgentReasoningParamSummaryDetailed AgentReasoningParamSummary = "detailed"
	AgentReasoningParamSummaryAuto     AgentReasoningParamSummary = "auto"
)

// A request to resume a subagent.
type AgentResumeSubagentCallItem struct {
	// The ID of the tool call item.
	ID string `json:"id" api:"required"`
	// The ID of the agent to resume.
	RecipientAgentID string `json:"recipient_agent_id" api:"required"`
	// The ID of the agent requesting the resume.
	SenderAgentID string `json:"sender_agent_id" api:"required"`
	// The status of the tool call.
	//
	// Any of "in_progress", "completed", "failed", "incomplete".
	Status AgentFunctionCallStatus `json:"status" api:"required"`
	// The ID of the turn that contains this item.
	TurnID string `json:"turn_id" api:"required"`
	// The item type. Always `resume_subagent_call`.
	Type constant.ResumeSubagentCall `json:"type" default:"resume_subagent_call"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID               respjson.Field
		RecipientAgentID respjson.Field
		SenderAgentID    respjson.Field
		Status           respjson.Field
		TurnID           respjson.Field
		Type             respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentResumeSubagentCallItem) RawJSON() string { return r.JSON.raw }
func (r *AgentResumeSubagentCallItem) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A request to send input to another agent.
type AgentSendSubagentInputCallItem struct {
	// The ID of the tool call item.
	ID string `json:"id" api:"required"`
	// The input sent to the receiving agent.
	Content []AgentContentUnion `json:"content" api:"required"`
	// The ID of the agent receiving the input.
	RecipientAgentID string `json:"recipient_agent_id" api:"required"`
	// The ID of the agent sending the input.
	SenderAgentID string `json:"sender_agent_id" api:"required"`
	// The status of the tool call.
	//
	// Any of "in_progress", "completed", "failed", "incomplete".
	Status AgentFunctionCallStatus `json:"status" api:"required"`
	// The ID of the turn that contains this item.
	TurnID string `json:"turn_id" api:"required"`
	// The item type. Always `send_subagent_input_call`.
	Type constant.SendSubagentInputCall `json:"type" default:"send_subagent_input_call"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID               respjson.Field
		Content          respjson.Field
		RecipientAgentID respjson.Field
		SenderAgentID    respjson.Field
		Status           respjson.Field
		TurnID           respjson.Field
		Type             respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentSendSubagentInputCallItem) RawJSON() string { return r.JSON.raw }
func (r *AgentSendSubagentInputCallItem) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A Managed Agents session.
type AgentSession struct {
	// The ID of the session.
	ID string `json:"id" api:"required"`
	// The agent running in the session.
	Agent AgentSessionAgent `json:"agent" api:"required"`
	// The Unix timestamp, in seconds, when the session was created.
	CreatedAt int64 `json:"created_at" api:"required"`
	// The execution environment for the session.
	Environment EnvironmentUnion `json:"environment" api:"required"`
	// The error that caused the session to fail, if any.
	Error string `json:"error" api:"required"`
	// The Unix timestamp, in seconds, when the session was last active.
	LastActiveAt int64 `json:"last_active_at" api:"required"`
	// Custom string key-value pairs attached to the session.
	Metadata map[string]string `json:"metadata" api:"required"`
	// The object type. Always `agent.session`.
	Object constant.AgentSession `json:"object" default:"agent.session"`
	// Actions that must be completed before the session can continue.
	RequiredActions []AgentSessionRequiredActionUnion `json:"required_actions" api:"required"`
	// The current status of the session.
	//
	// Any of "idle", "in_progress", "requires_action", "failed".
	Status AgentSessionStatus `json:"status" api:"required"`
	// Recorded token usage for a session or turn. Usage is best effort and may change.
	Usage TokenUsage `json:"usage" api:"required"`
	// The IDs of vaults made available to the session.
	VaultIDs []string `json:"vault_ids" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID              respjson.Field
		Agent           respjson.Field
		CreatedAt       respjson.Field
		Environment     respjson.Field
		Error           respjson.Field
		LastActiveAt    respjson.Field
		Metadata        respjson.Field
		Object          respjson.Field
		RequiredActions respjson.Field
		Status          respjson.Field
		Usage           respjson.Field
		VaultIDs        respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentSession) RawJSON() string { return r.JSON.raw }
func (r *AgentSession) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The agent running in the session.
type AgentSessionAgent struct {
	// The ID of the agent.
	ID string `json:"id" api:"required"`
	// Custom instructions appended to the agent's default base instructions.
	Instructions string `json:"instructions" api:"required"`
	// The model used by the agent.
	Model string `json:"model" api:"required"`
	// Configuration for creating and coordinating subagents.
	MultiAgent MultiAgentConfig `json:"multi_agent" api:"required"`
	// The reusable agent's name when the session was created, or null if no name was
	// saved. Later changes to the agent's name do not affect this value.
	Name string `json:"name" api:"required"`
	// The agent's reasoning configuration.
	Reasoning AgentReasoning `json:"reasoning" api:"required"`
	// The effective service-tier policy for model requests. Defaults to `auto`.
	//
	// Any of "auto", "default", "flex", "priority", "fast".
	ServiceTier string `json:"service_tier" api:"required"`
	// Configuration for text generated by the agent.
	Text AgentText `json:"text" api:"required"`
	// Tools available to the agent.
	Tools []AgentToolUnion `json:"tools" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID           respjson.Field
		Instructions respjson.Field
		Model        respjson.Field
		MultiAgent   respjson.Field
		Name         respjson.Field
		Reasoning    respjson.Field
		ServiceTier  respjson.Field
		Text         respjson.Field
		Tools        respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentSessionAgent) RawJSON() string { return r.JSON.raw }
func (r *AgentSessionAgent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// AgentSessionRequiredActionUnion contains all possible properties and values from
// [AgentSessionRequiredActionFunctionCall],
// [AgentSessionRequiredActionEnvironmentConnection].
//
// Use the [AgentSessionRequiredActionUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type AgentSessionRequiredActionUnion struct {
	// This field is from variant [AgentSessionRequiredActionFunctionCall].
	Arguments any `json:"arguments"`
	// This field is from variant [AgentSessionRequiredActionFunctionCall].
	CallID string `json:"call_id"`
	// This field is from variant [AgentSessionRequiredActionFunctionCall].
	Name string `json:"name"`
	// This field is from variant [AgentSessionRequiredActionFunctionCall].
	TurnID string `json:"turn_id"`
	// Any of "function_call", "environment_connection".
	Type string `json:"type"`
	// This field is from variant [AgentSessionRequiredActionEnvironmentConnection].
	EnvironmentID string `json:"environment_id"`
	JSON          struct {
		Arguments     respjson.Field
		CallID        respjson.Field
		Name          respjson.Field
		TurnID        respjson.Field
		Type          respjson.Field
		EnvironmentID respjson.Field
		raw           string
	} `json:"-"`
}

// anyAgentSessionRequiredAction is implemented by each variant of
// [AgentSessionRequiredActionUnion] to add type safety for the return type of
// [AgentSessionRequiredActionUnion.AsAny]
type anyAgentSessionRequiredAction interface {
	implAgentSessionRequiredActionUnion()
}

func (AgentSessionRequiredActionFunctionCall) implAgentSessionRequiredActionUnion()          {}
func (AgentSessionRequiredActionEnvironmentConnection) implAgentSessionRequiredActionUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := AgentSessionRequiredActionUnion.AsAny().(type) {
//	case openai.AgentSessionRequiredActionFunctionCall:
//	case openai.AgentSessionRequiredActionEnvironmentConnection:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u AgentSessionRequiredActionUnion) AsAny() anyAgentSessionRequiredAction {
	switch u.Type {
	case "function_call":
		return u.AsFunctionCall()
	case "environment_connection":
		return u.AsEnvironmentConnection()
	}
	return nil
}

func (u AgentSessionRequiredActionUnion) AsFunctionCall() (v AgentSessionRequiredActionFunctionCall) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentSessionRequiredActionUnion) AsEnvironmentConnection() (v AgentSessionRequiredActionEnvironmentConnection) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u AgentSessionRequiredActionUnion) RawJSON() string { return u.JSON.raw }

func (r *AgentSessionRequiredActionUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Run a function tool and submit its result.
type AgentSessionRequiredActionFunctionCall struct {
	// The arguments supplied by the model.
	Arguments any `json:"arguments" api:"required"`
	// The ID to include when submitting the function result.
	CallID string `json:"call_id" api:"required"`
	// The function name.
	Name string `json:"name" api:"required"`
	// The ID of the turn that requested the function call.
	TurnID string `json:"turn_id" api:"required"`
	// The type of the object. Always `function_call`.
	Type constant.FunctionCall `json:"type" default:"function_call"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Arguments   respjson.Field
		CallID      respjson.Field
		Name        respjson.Field
		TurnID      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentSessionRequiredActionFunctionCall) RawJSON() string { return r.JSON.raw }
func (r *AgentSessionRequiredActionFunctionCall) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Reconnect a session environment.
type AgentSessionRequiredActionEnvironmentConnection struct {
	// The ID of the environment to reconnect.
	EnvironmentID string `json:"environment_id" api:"required"`
	// The type of the object. Always `environment_connection`.
	Type constant.EnvironmentConnection `json:"type" default:"environment_connection"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EnvironmentID respjson.Field
		Type          respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentSessionRequiredActionEnvironmentConnection) RawJSON() string { return r.JSON.raw }
func (r *AgentSessionRequiredActionEnvironmentConnection) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The current status of the session.
type AgentSessionStatus string

const (
	AgentSessionStatusIdle           AgentSessionStatus = "idle"
	AgentSessionStatusInProgress     AgentSessionStatus = "in_progress"
	AgentSessionStatusRequiresAction AgentSessionStatus = "requires_action"
	AgentSessionStatusFailed         AgentSessionStatus = "failed"
)

// An assistant message produced by the agent.
type AgentSessionAssistantMessage struct {
	// The ID of the message.
	ID string `json:"id" api:"required"`
	// The content of the message.
	Content []OutputText `json:"content" api:"required"`
	// The phase of an assistant message.
	//
	// Any of "commentary", "final_answer".
	Phase AgentSessionAssistantMessagePhase `json:"phase" api:"required"`
	// The role of the message author. Always `assistant`.
	Role constant.Assistant `json:"role" default:"assistant"`
	// The status of the message.
	//
	// Any of "in_progress", "completed", "incomplete".
	Status AgentOutputItemStatus `json:"status" api:"required"`
	// The ID of the turn that contains this item.
	TurnID string `json:"turn_id" api:"required"`
	// The item type. Always `message`.
	Type constant.Message `json:"type" default:"message"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Content     respjson.Field
		Phase       respjson.Field
		Role        respjson.Field
		Status      respjson.Field
		TurnID      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentSessionAssistantMessage) RawJSON() string { return r.JSON.raw }
func (r *AgentSessionAssistantMessage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The phase of an assistant message.
type AgentSessionAssistantMessagePhase string

const (
	AgentSessionAssistantMessagePhaseCommentary  AgentSessionAssistantMessagePhase = "commentary"
	AgentSessionAssistantMessagePhaseFinalAnswer AgentSessionAssistantMessagePhase = "final_answer"
)

// Emitted when a session is created.
type AgentSessionCreatedEvent struct {
	// The unique ID of the event.
	EventID string `json:"event_id" api:"required"`
	// The session that was created.
	Session AgentSession `json:"session" api:"required"`
	// The type of the object. Always `agent.session.created`.
	Type constant.AgentSessionCreated `json:"type" default:"agent.session.created"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventID     respjson.Field
		Session     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentSessionCreatedEvent) RawJSON() string { return r.JSON.raw }
func (r *AgentSessionCreatedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A Managed Agents session removed from the public API. Physical cleanup may
// continue asynchronously.
type AgentSessionDeleted struct {
	// The ID of the deleted session.
	ID string `json:"id" api:"required"`
	// Whether the session has been removed from the public API. Always `true`.
	// Physical cleanup may still be in progress.
	Deleted bool `json:"deleted" api:"required"`
	// The object type. Always `agent.session.deleted`.
	Object constant.AgentSessionDeleted `json:"object" default:"agent.session.deleted"`
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
func (r AgentSessionDeleted) RawJSON() string { return r.JSON.raw }
func (r *AgentSessionDeleted) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Emitted when a session environment connects.
type AgentSessionEnvironmentConnectedEvent struct {
	// The current environment state.
	Environment AgentSessionEnvironmentState `json:"environment" api:"required"`
	// The unique ID of the event.
	EventID string `json:"event_id" api:"required"`
	// The ID of the session associated with the event.
	SessionID string `json:"session_id" api:"required"`
	// The ID of the turn associated with the event, when applicable.
	TurnID string `json:"turn_id" api:"required"`
	// The type of the object. Always `agent.session.environment.connected`.
	Type constant.AgentSessionEnvironmentConnected `json:"type" default:"agent.session.environment.connected"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Environment respjson.Field
		EventID     respjson.Field
		SessionID   respjson.Field
		TurnID      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentSessionEnvironmentConnectedEvent) RawJSON() string { return r.JSON.raw }
func (r *AgentSessionEnvironmentConnectedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Emitted when a session environment disconnects.
type AgentSessionEnvironmentDisconnectedEvent struct {
	// The current environment state.
	Environment AgentSessionEnvironmentState `json:"environment" api:"required"`
	// The unique ID of the event.
	EventID string `json:"event_id" api:"required"`
	// The ID of the session associated with the event.
	SessionID string `json:"session_id" api:"required"`
	// The ID of the turn associated with the event, when applicable.
	TurnID string `json:"turn_id" api:"required"`
	// The type of the object. Always `agent.session.environment.disconnected`.
	Type constant.AgentSessionEnvironmentDisconnected `json:"type" default:"agent.session.environment.disconnected"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Environment respjson.Field
		EventID     respjson.Field
		SessionID   respjson.Field
		TurnID      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentSessionEnvironmentDisconnectedEvent) RawJSON() string { return r.JSON.raw }
func (r *AgentSessionEnvironmentDisconnectedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Emitted when a session environment fails.
type AgentSessionEnvironmentFailedEvent struct {
	// The current environment state.
	Environment AgentSessionEnvironmentState `json:"environment" api:"required"`
	// The unique ID of the event.
	EventID string `json:"event_id" api:"required"`
	// The ID of the session associated with the event.
	SessionID string `json:"session_id" api:"required"`
	// The ID of the turn associated with the event, when applicable.
	TurnID string `json:"turn_id" api:"required"`
	// The type of the object. Always `agent.session.environment.failed`.
	Type constant.AgentSessionEnvironmentFailed `json:"type" default:"agent.session.environment.failed"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Environment respjson.Field
		EventID     respjson.Field
		SessionID   respjson.Field
		TurnID      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentSessionEnvironmentFailedEvent) RawJSON() string { return r.JSON.raw }
func (r *AgentSessionEnvironmentFailedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Emitted while a session environment is being prepared.
type AgentSessionEnvironmentPendingEvent struct {
	// The current environment state.
	Environment AgentSessionEnvironmentState `json:"environment" api:"required"`
	// The unique ID of the event.
	EventID string `json:"event_id" api:"required"`
	// The ID of the session associated with the event.
	SessionID string `json:"session_id" api:"required"`
	// The ID of the turn associated with the event, when applicable.
	TurnID string `json:"turn_id" api:"required"`
	// The type of the object. Always `agent.session.environment.pending`.
	Type constant.AgentSessionEnvironmentPending `json:"type" default:"agent.session.environment.pending"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Environment respjson.Field
		EventID     respjson.Field
		SessionID   respjson.Field
		TurnID      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentSessionEnvironmentPendingEvent) RawJSON() string { return r.JSON.raw }
func (r *AgentSessionEnvironmentPendingEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Emitted when a hosted session environment is ready to connect.
type AgentSessionEnvironmentReadyEvent struct {
	// The current environment state.
	Environment AgentSessionEnvironmentState `json:"environment" api:"required"`
	// The unique ID of the event.
	EventID string `json:"event_id" api:"required"`
	// The ID of the session associated with the event.
	SessionID string `json:"session_id" api:"required"`
	// The ID of the turn associated with the event, when applicable.
	TurnID string `json:"turn_id" api:"required"`
	// The type of the object. Always `agent.session.environment.ready`.
	Type constant.AgentSessionEnvironmentReady `json:"type" default:"agent.session.environment.ready"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Environment respjson.Field
		EventID     respjson.Field
		SessionID   respjson.Field
		TurnID      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentSessionEnvironmentReadyEvent) RawJSON() string { return r.JSON.raw }
func (r *AgentSessionEnvironmentReadyEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The current state of a session environment.
type AgentSessionEnvironmentState struct {
	// The public ID of the environment.
	ID string `json:"id" api:"required"`
	// An error reported while preparing a session environment.
	Error AgentSessionEnvironmentStateError `json:"error" api:"required"`
	// The environment's connection status.
	//
	// Any of "pending", "ready", "connected", "disconnected", "failed".
	Status AgentSessionEnvironmentStateStatus `json:"status" api:"required"`
	// The environment type.
	Type string `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Error       respjson.Field
		Status      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentSessionEnvironmentState) RawJSON() string { return r.JSON.raw }
func (r *AgentSessionEnvironmentState) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// An error reported while preparing a session environment.
type AgentSessionEnvironmentStateError struct {
	// A machine-readable error code.
	Code string `json:"code" api:"required"`
	// A human-readable error message.
	Message string `json:"message" api:"required"`
	// The error type.
	Type string `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Code        respjson.Field
		Message     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentSessionEnvironmentStateError) RawJSON() string { return r.JSON.raw }
func (r *AgentSessionEnvironmentStateError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The environment's connection status.
type AgentSessionEnvironmentStateStatus string

const (
	AgentSessionEnvironmentStateStatusPending      AgentSessionEnvironmentStateStatus = "pending"
	AgentSessionEnvironmentStateStatusReady        AgentSessionEnvironmentStateStatus = "ready"
	AgentSessionEnvironmentStateStatusConnected    AgentSessionEnvironmentStateStatus = "connected"
	AgentSessionEnvironmentStateStatusDisconnected AgentSessionEnvironmentStateStatus = "disconnected"
	AgentSessionEnvironmentStateStatusFailed       AgentSessionEnvironmentStateStatus = "failed"
)

// Emitted when a turn or session fails.
type AgentSessionErrorEvent struct {
	// The error that occurred.
	Error SessionError `json:"error" api:"required"`
	// The unique ID of the event.
	EventID string `json:"event_id" api:"required"`
	// The ID of the session associated with the event.
	SessionID string `json:"session_id" api:"required"`
	// The type of the object. Always `error`.
	Type constant.Error `json:"type" default:"error"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Error       respjson.Field
		EventID     respjson.Field
		SessionID   respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentSessionErrorEvent) RawJSON() string { return r.JSON.raw }
func (r *AgentSessionErrorEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// AgentSessionEventUnion contains all possible properties and values from
// [AgentSessionErrorEvent], [AgentSessionEnvironmentReadyEvent],
// [AgentOutputCommandExecutionOutputDeltaEvent], [AgentSessionCreatedEvent],
// [AgentSessionTurnCreatedEvent], [AgentSessionTurnInProgressEvent],
// [AgentSessionTurnCompletedEvent], [AgentSessionTurnFailedEvent],
// [AgentSessionTurnCancelledEvent], [AgentSessionTurnItemAddedEvent],
// [AgentSessionIdleEvent], [AgentSessionInProgressEvent],
// [AgentSessionRequiresActionEvent], [AgentSessionFailedEvent],
// [AgentSessionEnvironmentPendingEvent], [AgentSessionEnvironmentConnectedEvent],
// [AgentSessionEnvironmentDisconnectedEvent],
// [AgentSessionEnvironmentFailedEvent], [AgentSessionSubagentCreatedEvent],
// [AgentSessionSubagentActiveEvent], [AgentSessionSubagentClosedEvent],
// [AgentSessionTurnItemDoneEvent], [AgentSessionTurnContentPartAddedEvent],
// [AgentSessionTurnContentPartDoneEvent], [AgentSessionTurnOutputTextDeltaEvent],
// [AgentSessionTurnOutputTextDoneEvent],
// [AgentSessionTurnReasoningSummaryPartAddedEvent],
// [AgentSessionTurnReasoningSummaryPartDoneEvent],
// [AgentSessionTurnReasoningSummaryTextDeltaEvent],
// [AgentSessionTurnReasoningSummaryTextDoneEvent].
//
// Use the [AgentSessionEventUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type AgentSessionEventUnion struct {
	// This field is from variant [AgentSessionErrorEvent].
	Error     SessionError `json:"error"`
	EventID   string       `json:"event_id"`
	SessionID string       `json:"session_id"`
	// Any of "error", "agent.session.environment.ready",
	// "agent.output.command_execution_output.delta", "agent.session.created",
	// "agent.session.turn.created", "agent.session.turn.in_progress",
	// "agent.session.turn.completed", "agent.session.turn.failed",
	// "agent.session.turn.cancelled", "agent.session.turn.item.added",
	// "agent.session.idle", "agent.session.in_progress",
	// "agent.session.requires_action", "agent.session.failed",
	// "agent.session.environment.pending", "agent.session.environment.connected",
	// "agent.session.environment.disconnected", "agent.session.environment.failed",
	// "agent.session.subagent.created", "agent.session.subagent.active",
	// "agent.session.subagent.closed", "agent.session.turn.item.done",
	// "agent.session.turn.content_part.added", "agent.session.turn.content_part.done",
	// "agent.session.turn.output_text.delta", "agent.session.turn.output_text.done",
	// "agent.session.turn.reasoning_summary_part.added",
	// "agent.session.turn.reasoning_summary_part.done",
	// "agent.session.turn.reasoning_summary_text.delta",
	// "agent.session.turn.reasoning_summary_text.done".
	Type string `json:"type"`
	// This field is from variant [AgentSessionEnvironmentReadyEvent].
	Environment AgentSessionEnvironmentState `json:"environment"`
	TurnID      string                       `json:"turn_id"`
	Delta       string                       `json:"delta"`
	ItemID      string                       `json:"item_id"`
	OutputIndex int64                        `json:"output_index"`
	// This field is from variant [AgentSessionCreatedEvent].
	Session AgentSession `json:"session"`
	// This field is from variant [AgentSessionTurnCreatedEvent].
	Turn Turn `json:"turn"`
	// This field is from variant [AgentSessionTurnCompletedEvent].
	Usage TokenUsage `json:"usage"`
	// This field is a union of [AgentSessionItemUnion], [AgentOutputItemUnion]
	Item AgentSessionEventUnionItem `json:"item"`
	// This field is from variant [AgentSessionSubagentCreatedEvent].
	Subagent     Subagent `json:"subagent"`
	ContentIndex int64    `json:"content_index"`
	// This field is a union of [OutputText], [SummaryText]
	Part         AgentSessionEventUnionPart `json:"part"`
	Text         string                     `json:"text"`
	SummaryIndex int64                      `json:"summary_index"`
	// This field is from variant [AgentSessionTurnReasoningSummaryPartDoneEvent].
	Status constant.Incomplete `json:"status"`
	JSON   struct {
		Error        respjson.Field
		EventID      respjson.Field
		SessionID    respjson.Field
		Type         respjson.Field
		Environment  respjson.Field
		TurnID       respjson.Field
		Delta        respjson.Field
		ItemID       respjson.Field
		OutputIndex  respjson.Field
		Session      respjson.Field
		Turn         respjson.Field
		Usage        respjson.Field
		Item         respjson.Field
		Subagent     respjson.Field
		ContentIndex respjson.Field
		Part         respjson.Field
		Text         respjson.Field
		SummaryIndex respjson.Field
		Status       respjson.Field
		raw          string
	} `json:"-"`
}

// anyAgentSessionEvent is implemented by each variant of [AgentSessionEventUnion]
// to add type safety for the return type of [AgentSessionEventUnion.AsAny]
type anyAgentSessionEvent interface {
	implAgentSessionEventUnion()
}

func (AgentSessionErrorEvent) implAgentSessionEventUnion()                         {}
func (AgentSessionEnvironmentReadyEvent) implAgentSessionEventUnion()              {}
func (AgentOutputCommandExecutionOutputDeltaEvent) implAgentSessionEventUnion()    {}
func (AgentSessionCreatedEvent) implAgentSessionEventUnion()                       {}
func (AgentSessionTurnCreatedEvent) implAgentSessionEventUnion()                   {}
func (AgentSessionTurnInProgressEvent) implAgentSessionEventUnion()                {}
func (AgentSessionTurnCompletedEvent) implAgentSessionEventUnion()                 {}
func (AgentSessionTurnFailedEvent) implAgentSessionEventUnion()                    {}
func (AgentSessionTurnCancelledEvent) implAgentSessionEventUnion()                 {}
func (AgentSessionTurnItemAddedEvent) implAgentSessionEventUnion()                 {}
func (AgentSessionIdleEvent) implAgentSessionEventUnion()                          {}
func (AgentSessionInProgressEvent) implAgentSessionEventUnion()                    {}
func (AgentSessionRequiresActionEvent) implAgentSessionEventUnion()                {}
func (AgentSessionFailedEvent) implAgentSessionEventUnion()                        {}
func (AgentSessionEnvironmentPendingEvent) implAgentSessionEventUnion()            {}
func (AgentSessionEnvironmentConnectedEvent) implAgentSessionEventUnion()          {}
func (AgentSessionEnvironmentDisconnectedEvent) implAgentSessionEventUnion()       {}
func (AgentSessionEnvironmentFailedEvent) implAgentSessionEventUnion()             {}
func (AgentSessionSubagentCreatedEvent) implAgentSessionEventUnion()               {}
func (AgentSessionSubagentActiveEvent) implAgentSessionEventUnion()                {}
func (AgentSessionSubagentClosedEvent) implAgentSessionEventUnion()                {}
func (AgentSessionTurnItemDoneEvent) implAgentSessionEventUnion()                  {}
func (AgentSessionTurnContentPartAddedEvent) implAgentSessionEventUnion()          {}
func (AgentSessionTurnContentPartDoneEvent) implAgentSessionEventUnion()           {}
func (AgentSessionTurnOutputTextDeltaEvent) implAgentSessionEventUnion()           {}
func (AgentSessionTurnOutputTextDoneEvent) implAgentSessionEventUnion()            {}
func (AgentSessionTurnReasoningSummaryPartAddedEvent) implAgentSessionEventUnion() {}
func (AgentSessionTurnReasoningSummaryPartDoneEvent) implAgentSessionEventUnion()  {}
func (AgentSessionTurnReasoningSummaryTextDeltaEvent) implAgentSessionEventUnion() {}
func (AgentSessionTurnReasoningSummaryTextDoneEvent) implAgentSessionEventUnion()  {}

// Use the following switch statement to find the correct variant
//
//	switch variant := AgentSessionEventUnion.AsAny().(type) {
//	case openai.AgentSessionErrorEvent:
//	case openai.AgentSessionEnvironmentReadyEvent:
//	case openai.AgentOutputCommandExecutionOutputDeltaEvent:
//	case openai.AgentSessionCreatedEvent:
//	case openai.AgentSessionTurnCreatedEvent:
//	case openai.AgentSessionTurnInProgressEvent:
//	case openai.AgentSessionTurnCompletedEvent:
//	case openai.AgentSessionTurnFailedEvent:
//	case openai.AgentSessionTurnCancelledEvent:
//	case openai.AgentSessionTurnItemAddedEvent:
//	case openai.AgentSessionIdleEvent:
//	case openai.AgentSessionInProgressEvent:
//	case openai.AgentSessionRequiresActionEvent:
//	case openai.AgentSessionFailedEvent:
//	case openai.AgentSessionEnvironmentPendingEvent:
//	case openai.AgentSessionEnvironmentConnectedEvent:
//	case openai.AgentSessionEnvironmentDisconnectedEvent:
//	case openai.AgentSessionEnvironmentFailedEvent:
//	case openai.AgentSessionSubagentCreatedEvent:
//	case openai.AgentSessionSubagentActiveEvent:
//	case openai.AgentSessionSubagentClosedEvent:
//	case openai.AgentSessionTurnItemDoneEvent:
//	case openai.AgentSessionTurnContentPartAddedEvent:
//	case openai.AgentSessionTurnContentPartDoneEvent:
//	case openai.AgentSessionTurnOutputTextDeltaEvent:
//	case openai.AgentSessionTurnOutputTextDoneEvent:
//	case openai.AgentSessionTurnReasoningSummaryPartAddedEvent:
//	case openai.AgentSessionTurnReasoningSummaryPartDoneEvent:
//	case openai.AgentSessionTurnReasoningSummaryTextDeltaEvent:
//	case openai.AgentSessionTurnReasoningSummaryTextDoneEvent:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u AgentSessionEventUnion) AsAny() anyAgentSessionEvent {
	switch u.Type {
	case "error":
		return u.AsError()
	case "agent.session.environment.ready":
		return u.AsAgentSessionEnvironmentReady()
	case "agent.output.command_execution_output.delta":
		return u.AsAgentOutputCommandExecutionOutputDelta()
	case "agent.session.created":
		return u.AsAgentSessionCreated()
	case "agent.session.turn.created":
		return u.AsAgentSessionTurnCreated()
	case "agent.session.turn.in_progress":
		return u.AsAgentSessionTurnInProgress()
	case "agent.session.turn.completed":
		return u.AsAgentSessionTurnCompleted()
	case "agent.session.turn.failed":
		return u.AsAgentSessionTurnFailed()
	case "agent.session.turn.cancelled":
		return u.AsAgentSessionTurnCancelled()
	case "agent.session.turn.item.added":
		return u.AsAgentSessionTurnItemAdded()
	case "agent.session.idle":
		return u.AsAgentSessionIdle()
	case "agent.session.in_progress":
		return u.AsAgentSessionInProgress()
	case "agent.session.requires_action":
		return u.AsAgentSessionRequiresAction()
	case "agent.session.failed":
		return u.AsAgentSessionFailed()
	case "agent.session.environment.pending":
		return u.AsAgentSessionEnvironmentPending()
	case "agent.session.environment.connected":
		return u.AsAgentSessionEnvironmentConnected()
	case "agent.session.environment.disconnected":
		return u.AsAgentSessionEnvironmentDisconnected()
	case "agent.session.environment.failed":
		return u.AsAgentSessionEnvironmentFailed()
	case "agent.session.subagent.created":
		return u.AsAgentSessionSubagentCreated()
	case "agent.session.subagent.active":
		return u.AsAgentSessionSubagentActive()
	case "agent.session.subagent.closed":
		return u.AsAgentSessionSubagentClosed()
	case "agent.session.turn.item.done":
		return u.AsAgentSessionTurnItemDone()
	case "agent.session.turn.content_part.added":
		return u.AsAgentSessionTurnContentPartAdded()
	case "agent.session.turn.content_part.done":
		return u.AsAgentSessionTurnContentPartDone()
	case "agent.session.turn.output_text.delta":
		return u.AsAgentSessionTurnOutputTextDelta()
	case "agent.session.turn.output_text.done":
		return u.AsAgentSessionTurnOutputTextDone()
	case "agent.session.turn.reasoning_summary_part.added":
		return u.AsAgentSessionTurnReasoningSummaryPartAdded()
	case "agent.session.turn.reasoning_summary_part.done":
		return u.AsAgentSessionTurnReasoningSummaryPartDone()
	case "agent.session.turn.reasoning_summary_text.delta":
		return u.AsAgentSessionTurnReasoningSummaryTextDelta()
	case "agent.session.turn.reasoning_summary_text.done":
		return u.AsAgentSessionTurnReasoningSummaryTextDone()
	}
	return nil
}

func (u AgentSessionEventUnion) AsError() (v AgentSessionErrorEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentSessionEventUnion) AsAgentSessionEnvironmentReady() (v AgentSessionEnvironmentReadyEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentSessionEventUnion) AsAgentOutputCommandExecutionOutputDelta() (v AgentOutputCommandExecutionOutputDeltaEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentSessionEventUnion) AsAgentSessionCreated() (v AgentSessionCreatedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentSessionEventUnion) AsAgentSessionTurnCreated() (v AgentSessionTurnCreatedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentSessionEventUnion) AsAgentSessionTurnInProgress() (v AgentSessionTurnInProgressEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentSessionEventUnion) AsAgentSessionTurnCompleted() (v AgentSessionTurnCompletedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentSessionEventUnion) AsAgentSessionTurnFailed() (v AgentSessionTurnFailedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentSessionEventUnion) AsAgentSessionTurnCancelled() (v AgentSessionTurnCancelledEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentSessionEventUnion) AsAgentSessionTurnItemAdded() (v AgentSessionTurnItemAddedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentSessionEventUnion) AsAgentSessionIdle() (v AgentSessionIdleEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentSessionEventUnion) AsAgentSessionInProgress() (v AgentSessionInProgressEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentSessionEventUnion) AsAgentSessionRequiresAction() (v AgentSessionRequiresActionEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentSessionEventUnion) AsAgentSessionFailed() (v AgentSessionFailedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentSessionEventUnion) AsAgentSessionEnvironmentPending() (v AgentSessionEnvironmentPendingEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentSessionEventUnion) AsAgentSessionEnvironmentConnected() (v AgentSessionEnvironmentConnectedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentSessionEventUnion) AsAgentSessionEnvironmentDisconnected() (v AgentSessionEnvironmentDisconnectedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentSessionEventUnion) AsAgentSessionEnvironmentFailed() (v AgentSessionEnvironmentFailedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentSessionEventUnion) AsAgentSessionSubagentCreated() (v AgentSessionSubagentCreatedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentSessionEventUnion) AsAgentSessionSubagentActive() (v AgentSessionSubagentActiveEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentSessionEventUnion) AsAgentSessionSubagentClosed() (v AgentSessionSubagentClosedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentSessionEventUnion) AsAgentSessionTurnItemDone() (v AgentSessionTurnItemDoneEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentSessionEventUnion) AsAgentSessionTurnContentPartAdded() (v AgentSessionTurnContentPartAddedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentSessionEventUnion) AsAgentSessionTurnContentPartDone() (v AgentSessionTurnContentPartDoneEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentSessionEventUnion) AsAgentSessionTurnOutputTextDelta() (v AgentSessionTurnOutputTextDeltaEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentSessionEventUnion) AsAgentSessionTurnOutputTextDone() (v AgentSessionTurnOutputTextDoneEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentSessionEventUnion) AsAgentSessionTurnReasoningSummaryPartAdded() (v AgentSessionTurnReasoningSummaryPartAddedEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentSessionEventUnion) AsAgentSessionTurnReasoningSummaryPartDone() (v AgentSessionTurnReasoningSummaryPartDoneEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentSessionEventUnion) AsAgentSessionTurnReasoningSummaryTextDelta() (v AgentSessionTurnReasoningSummaryTextDeltaEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentSessionEventUnion) AsAgentSessionTurnReasoningSummaryTextDone() (v AgentSessionTurnReasoningSummaryTextDoneEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u AgentSessionEventUnion) RawJSON() string { return u.JSON.raw }

func (r *AgentSessionEventUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// AgentSessionEventUnionItem is an implicit subunion of [AgentSessionEventUnion].
// AgentSessionEventUnionItem provides convenient access to the sub-properties of
// the union.
//
// For type safety it is recommended to directly use a variant of the
// [AgentSessionEventUnion].
type AgentSessionEventUnionItem struct {
	ID string `json:"id"`
	// This field is a union of [[]AgentSessionMessageContentUnion],
	// [[]AgentContentUnion], [[]AgentContentUnion], [[]AgentContentUnion],
	// [[]OutputText]
	Content AgentSessionEventUnionItemContent `json:"content"`
	Phase   string                            `json:"phase"`
	Role    string                            `json:"role"`
	Status  string                            `json:"status"`
	TurnID  string                            `json:"turn_id"`
	Type    string                            `json:"type"`
	// This field is from variant [AgentSessionItemUnion], [AgentOutputItemUnion].
	Summary   []SummaryText `json:"summary"`
	Arguments any           `json:"arguments"`
	CallID    string        `json:"call_id"`
	Name      string        `json:"name"`
	// This field is a union of [string], [any]
	Error AgentSessionEventUnionItemError `json:"error"`
	// This field is a union of [AgentFunctionCallOutputUnion], [any], [string]
	Output           AgentSessionEventUnionItemOutput `json:"output"`
	RecipientAgentID string                           `json:"recipient_agent_id"`
	SenderAgentID    string                           `json:"sender_agent_id"`
	// This field is from variant [AgentSessionItemUnion], [AgentOutputItemUnion].
	ServerLabel string `json:"server_label"`
	// This field is from variant [AgentSessionItemUnion], [AgentOutputItemUnion].
	Action WebSearchActionUnion `json:"action"`
	// This field is from variant [AgentSessionItemUnion], [AgentOutputItemUnion].
	Command string `json:"command"`
	// This field is from variant [AgentSessionItemUnion], [AgentOutputItemUnion].
	Cwd string `json:"cwd"`
	// This field is from variant [AgentSessionItemUnion], [AgentOutputItemUnion].
	DurationMs int64 `json:"duration_ms"`
	// This field is from variant [AgentSessionItemUnion], [AgentOutputItemUnion].
	ExitCode int64 `json:"exit_code"`
	// This field is from variant [AgentSessionItemUnion], [AgentOutputItemUnion].
	AgentID string `json:"agent_id"`
	// This field is from variant [AgentSessionItemUnion], [AgentOutputItemUnion].
	Model string `json:"model"`
	// This field is from variant [AgentSessionItemUnion], [AgentOutputItemUnion].
	ReasoningEffort string `json:"reasoning_effort"`
	// This field is from variant [AgentSessionItemUnion], [AgentOutputItemUnion].
	RecipientAgentIDs []string `json:"recipient_agent_ids"`
	JSON              struct {
		ID                respjson.Field
		Content           respjson.Field
		Phase             respjson.Field
		Role              respjson.Field
		Status            respjson.Field
		TurnID            respjson.Field
		Type              respjson.Field
		Summary           respjson.Field
		Arguments         respjson.Field
		CallID            respjson.Field
		Name              respjson.Field
		Error             respjson.Field
		Output            respjson.Field
		RecipientAgentID  respjson.Field
		SenderAgentID     respjson.Field
		ServerLabel       respjson.Field
		Action            respjson.Field
		Command           respjson.Field
		Cwd               respjson.Field
		DurationMs        respjson.Field
		ExitCode          respjson.Field
		AgentID           respjson.Field
		Model             respjson.Field
		ReasoningEffort   respjson.Field
		RecipientAgentIDs respjson.Field
		raw               string
	} `json:"-"`
}

func (r *AgentSessionEventUnionItem) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// AgentSessionEventUnionItemContent is an implicit subunion of
// [AgentSessionEventUnion]. AgentSessionEventUnionItemContent provides convenient
// access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [AgentSessionEventUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfAgentSessionMessageContentArray OfAgentContentArray
// OfOutputTextArray]
type AgentSessionEventUnionItemContent struct {
	// This field will be present if the value is a [[]AgentSessionMessageContentUnion]
	// instead of an object.
	OfAgentSessionMessageContentArray []AgentSessionMessageContentUnion `json:",inline"`
	// This field will be present if the value is a [[]AgentContentUnion] instead of an
	// object.
	OfAgentContentArray []AgentContentUnion `json:",inline"`
	// This field will be present if the value is a [[]OutputText] instead of an
	// object.
	OfOutputTextArray []OutputText `json:",inline"`
	JSON              struct {
		OfAgentSessionMessageContentArray respjson.Field
		OfAgentContentArray               respjson.Field
		OfOutputTextArray                 respjson.Field
		raw                               string
	} `json:"-"`
}

func (r *AgentSessionEventUnionItemContent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// AgentSessionEventUnionItemError is an implicit subunion of
// [AgentSessionEventUnion]. AgentSessionEventUnionItemError provides convenient
// access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [AgentSessionEventUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfString OfAgentMcpCallItemError]
type AgentSessionEventUnionItemError struct {
	// This field will be present if the value is a [string] instead of an object.
	OfString string `json:",inline"`
	// This field will be present if the value is a [any] instead of an object.
	OfAgentMcpCallItemError any `json:",inline"`
	JSON                    struct {
		OfString                respjson.Field
		OfAgentMcpCallItemError respjson.Field
		raw                     string
	} `json:"-"`
}

func (r *AgentSessionEventUnionItemError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// AgentSessionEventUnionItemOutput is an implicit subunion of
// [AgentSessionEventUnion]. AgentSessionEventUnionItemOutput provides convenient
// access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [AgentSessionEventUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfString OfInputContentArray OfAgentMcpCallItemOutput]
type AgentSessionEventUnionItemOutput struct {
	// This field will be present if the value is a [string] instead of an object.
	OfString string `json:",inline"`
	// This field will be present if the value is a [[]InputContentUnion] instead of an
	// object.
	OfInputContentArray []InputContentUnion `json:",inline"`
	// This field will be present if the value is a [any] instead of an object.
	OfAgentMcpCallItemOutput any `json:",inline"`
	JSON                     struct {
		OfString                 respjson.Field
		OfInputContentArray      respjson.Field
		OfAgentMcpCallItemOutput respjson.Field
		raw                      string
	} `json:"-"`
}

func (r *AgentSessionEventUnionItemOutput) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// AgentSessionEventUnionPart is an implicit subunion of [AgentSessionEventUnion].
// AgentSessionEventUnionPart provides convenient access to the sub-properties of
// the union.
//
// For type safety it is recommended to directly use a variant of the
// [AgentSessionEventUnion].
type AgentSessionEventUnionPart struct {
	Text string `json:"text"`
	Type string `json:"type"`
	JSON struct {
		Text respjson.Field
		Type respjson.Field
		raw  string
	} `json:"-"`
}

func (r *AgentSessionEventUnionPart) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Emitted when a session fails.
type AgentSessionFailedEvent struct {
	// The unique ID of the event.
	EventID string `json:"event_id" api:"required"`
	// The failed session.
	Session AgentSession `json:"session" api:"required"`
	// The type of the object. Always `agent.session.failed`.
	Type constant.AgentSessionFailed `json:"type" default:"agent.session.failed"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventID     respjson.Field
		Session     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentSessionFailedEvent) RawJSON() string { return r.JSON.raw }
func (r *AgentSessionFailedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Emitted when a session becomes idle.
type AgentSessionIdleEvent struct {
	// The unique ID of the event.
	EventID string `json:"event_id" api:"required"`
	// The session that became idle.
	Session AgentSession `json:"session" api:"required"`
	// The type of the object. Always `agent.session.idle`.
	Type constant.AgentSessionIdle `json:"type" default:"agent.session.idle"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventID     respjson.Field
		Session     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentSessionIdleEvent) RawJSON() string { return r.JSON.raw }
func (r *AgentSessionIdleEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Emitted when a session starts processing a turn.
type AgentSessionInProgressEvent struct {
	// The unique ID of the event.
	EventID string `json:"event_id" api:"required"`
	// The session that started processing.
	Session AgentSession `json:"session" api:"required"`
	// The type of the object. Always `agent.session.in_progress`.
	Type constant.AgentSessionInProgress `json:"type" default:"agent.session.in_progress"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventID     respjson.Field
		Session     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentSessionInProgressEvent) RawJSON() string { return r.JSON.raw }
func (r *AgentSessionInProgressEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A user message submitted to a session.
//
// The properties Content, Role are required.
type AgentSessionInputMessageParam struct {
	// The content of the message.
	Content []InputContentParamUnion `json:"content,omitzero" api:"required"`
	// The type of the input item. Always `message`.
	//
	// Any of "message".
	Type AgentSessionInputMessageParamType `json:"type,omitzero"`
	// The role of the message author. Always `user`.
	//
	// This field can be elided, and will marshal its zero value as "user".
	Role constant.User `json:"role" default:"user"`
	paramObj
}

func (r AgentSessionInputMessageParam) MarshalJSON() (data []byte, err error) {
	type shadow AgentSessionInputMessageParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AgentSessionInputMessageParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The type of the input item. Always `message`.
type AgentSessionInputMessageParamType string

const (
	AgentSessionInputMessageParamTypeMessage AgentSessionInputMessageParamType = "message"
)

func AgentSessionInputParamOfParamAgentSessionInputMessage(input []AgentSessionInputMessageParam) AgentSessionInputParamUnion {
	var paramAgentSessionInputMessage AgentSessionInputParamAgentSessionInputMessage
	paramAgentSessionInputMessage.Input = input
	return AgentSessionInputParamUnion{OfParamAgentSessionInputMessage: &paramAgentSessionInputMessage}
}

func AgentSessionInputParamOfParamAgentSessionInputToolResult(callID string, success bool, turnID string) AgentSessionInputParamUnion {
	var paramAgentSessionInputToolResult AgentSessionInputParamAgentSessionInputToolResult
	paramAgentSessionInputToolResult.CallID = callID
	paramAgentSessionInputToolResult.Success = success
	paramAgentSessionInputToolResult.TurnID = turnID
	return AgentSessionInputParamUnion{OfParamAgentSessionInputToolResult: &paramAgentSessionInputToolResult}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type AgentSessionInputParamUnion struct {
	OfParamAgentSessionInputMessage    *AgentSessionInputParamAgentSessionInputMessage    `json:",omitzero,inline"`
	OfParamAgentSessionInputCancel     *AgentSessionInputParamAgentSessionInputCancel     `json:",omitzero,inline"`
	OfParamAgentSessionInputToolResult *AgentSessionInputParamAgentSessionInputToolResult `json:",omitzero,inline"`
	paramUnion
}

func (u AgentSessionInputParamUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfParamAgentSessionInputMessage, u.OfParamAgentSessionInputCancel, u.OfParamAgentSessionInputToolResult)
}
func (u *AgentSessionInputParamUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u AgentSessionInputParamUnion) GetInput() []AgentSessionInputMessageParam {
	if vt := u.OfParamAgentSessionInputMessage; vt != nil {
		return vt.Input
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AgentSessionInputParamUnion) GetCallID() *string {
	if vt := u.OfParamAgentSessionInputToolResult; vt != nil {
		return &vt.CallID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AgentSessionInputParamUnion) GetSuccess() *bool {
	if vt := u.OfParamAgentSessionInputToolResult; vt != nil {
		return &vt.Success
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AgentSessionInputParamUnion) GetTurnID() *string {
	if vt := u.OfParamAgentSessionInputToolResult; vt != nil {
		return &vt.TurnID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AgentSessionInputParamUnion) GetError() *string {
	if vt := u.OfParamAgentSessionInputToolResult; vt != nil && vt.Error.Valid() {
		return &vt.Error.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AgentSessionInputParamUnion) GetOutput() *AgentFunctionCallOutputParamUnion {
	if vt := u.OfParamAgentSessionInputToolResult; vt != nil {
		return &vt.Output
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AgentSessionInputParamUnion) GetType() *string {
	if vt := u.OfParamAgentSessionInputMessage; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfParamAgentSessionInputCancel; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfParamAgentSessionInputToolResult; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[AgentSessionInputParamUnion](
		"type",
		apijson.Discriminator[AgentSessionInputParamAgentSessionInputMessage]("agent.session.input.message"),
		apijson.Discriminator[AgentSessionInputParamAgentSessionInputCancel]("agent.session.input.cancel"),
		apijson.Discriminator[AgentSessionInputParamAgentSessionInputToolResult]("agent.session.input.tool_result"),
	)
}

// Adds one or more user messages and starts a turn.
//
// The properties Input, Type are required.
type AgentSessionInputParamAgentSessionInputMessage struct {
	// The user messages to add to the session.
	Input []AgentSessionInputMessageParam `json:"input,omitzero" api:"required"`
	// The type of the object. Always `agent.session.input.message`.
	//
	// This field can be elided, and will marshal its zero value as
	// "agent.session.input.message".
	Type constant.AgentSessionInputMessage `json:"type" default:"agent.session.input.message"`
	paramObj
}

func (r AgentSessionInputParamAgentSessionInputMessage) MarshalJSON() (data []byte, err error) {
	type shadow AgentSessionInputParamAgentSessionInputMessage
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AgentSessionInputParamAgentSessionInputMessage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewAgentSessionInputParamAgentSessionInputCancel() AgentSessionInputParamAgentSessionInputCancel {
	return AgentSessionInputParamAgentSessionInputCancel{
		Type: "agent.session.input.cancel",
	}
}

// Cancels the session's active turn.
//
// This struct has a constant value, construct it with
// [NewAgentSessionInputParamAgentSessionInputCancel].
type AgentSessionInputParamAgentSessionInputCancel struct {
	// The type of the object. Always `agent.session.input.cancel`.
	Type constant.AgentSessionInputCancel `json:"type" default:"agent.session.input.cancel"`
	paramObj
}

func (r AgentSessionInputParamAgentSessionInputCancel) MarshalJSON() (data []byte, err error) {
	type shadow AgentSessionInputParamAgentSessionInputCancel
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AgentSessionInputParamAgentSessionInputCancel) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Submits the result of a function call.
//
// The properties CallID, Success, TurnID, Type are required.
type AgentSessionInputParamAgentSessionInputToolResult struct {
	// The ID of the function call.
	CallID string `json:"call_id" api:"required"`
	// Whether the function call succeeded.
	Success bool `json:"success" api:"required"`
	// The ID of the turn that requested the function call.
	TurnID string `json:"turn_id" api:"required"`
	// The error message when the call failed.
	Error param.Opt[string] `json:"error,omitzero"`
	// A function result represented as text or supported model-input content.
	Output AgentFunctionCallOutputParamUnion `json:"output,omitzero"`
	// The type of the object. Always `agent.session.input.tool_result`.
	//
	// This field can be elided, and will marshal its zero value as
	// "agent.session.input.tool_result".
	Type constant.AgentSessionInputToolResult `json:"type" default:"agent.session.input.tool_result"`
	paramObj
}

func (r AgentSessionInputParamAgentSessionInputToolResult) MarshalJSON() (data []byte, err error) {
	type shadow AgentSessionInputParamAgentSessionInputToolResult
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AgentSessionInputParamAgentSessionInputToolResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// AgentSessionItemUnion contains all possible properties and values from
// [AgentSessionMessage], [AgentReasoningItem], [AgentFunctionCallItem],
// [AgentSessionItemFunctionCallOutput], [AgentSessionItemAgentMessage],
// [AgentMcpCallItem], [AgentWebSearchCallItem], [AgentCommandExecutionItem],
// [AgentCreateSubagentCallItem], [AgentSendSubagentInputCallItem],
// [AgentResumeSubagentCallItem], [AgentWaitForSubagentsCallItem],
// [AgentInterruptSubagentCallItem], [AgentCloseSubagentCallItem].
//
// Use the [AgentSessionItemUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type AgentSessionItemUnion struct {
	ID string `json:"id"`
	// This field is a union of [[]AgentSessionMessageContentUnion],
	// [[]AgentContentUnion], [[]AgentContentUnion], [[]AgentContentUnion]
	Content AgentSessionItemUnionContent `json:"content"`
	// This field is from variant [AgentSessionMessage].
	Phase AgentSessionMessagePhase `json:"phase"`
	// This field is from variant [AgentSessionMessage].
	Role   AgentSessionMessageRole `json:"role"`
	Status string                  `json:"status"`
	TurnID string                  `json:"turn_id"`
	// Any of "message", "reasoning", "function_call", "function_call_output",
	// "agent_message", "mcp_call", "web_search_call", "command_execution",
	// "create_subagent_call", "send_subagent_input_call", "resume_subagent_call",
	// "wait_for_subagents_call", "interrupt_subagent_call", "close_subagent_call".
	Type string `json:"type"`
	// This field is from variant [AgentReasoningItem].
	Summary   []SummaryText `json:"summary"`
	Arguments any           `json:"arguments"`
	CallID    string        `json:"call_id"`
	Name      string        `json:"name"`
	// This field is a union of [string], [any]
	Error AgentSessionItemUnionError `json:"error"`
	// This field is a union of [AgentFunctionCallOutputUnion], [any], [string]
	Output           AgentSessionItemUnionOutput `json:"output"`
	RecipientAgentID string                      `json:"recipient_agent_id"`
	SenderAgentID    string                      `json:"sender_agent_id"`
	// This field is from variant [AgentMcpCallItem].
	ServerLabel string `json:"server_label"`
	// This field is from variant [AgentWebSearchCallItem].
	Action WebSearchActionUnion `json:"action"`
	// This field is from variant [AgentCommandExecutionItem].
	Command string `json:"command"`
	// This field is from variant [AgentCommandExecutionItem].
	Cwd string `json:"cwd"`
	// This field is from variant [AgentCommandExecutionItem].
	DurationMs int64 `json:"duration_ms"`
	// This field is from variant [AgentCommandExecutionItem].
	ExitCode int64 `json:"exit_code"`
	// This field is from variant [AgentCreateSubagentCallItem].
	AgentID string `json:"agent_id"`
	// This field is from variant [AgentCreateSubagentCallItem].
	Model string `json:"model"`
	// This field is from variant [AgentCreateSubagentCallItem].
	ReasoningEffort string `json:"reasoning_effort"`
	// This field is from variant [AgentWaitForSubagentsCallItem].
	RecipientAgentIDs []string `json:"recipient_agent_ids"`
	JSON              struct {
		ID                respjson.Field
		Content           respjson.Field
		Phase             respjson.Field
		Role              respjson.Field
		Status            respjson.Field
		TurnID            respjson.Field
		Type              respjson.Field
		Summary           respjson.Field
		Arguments         respjson.Field
		CallID            respjson.Field
		Name              respjson.Field
		Error             respjson.Field
		Output            respjson.Field
		RecipientAgentID  respjson.Field
		SenderAgentID     respjson.Field
		ServerLabel       respjson.Field
		Action            respjson.Field
		Command           respjson.Field
		Cwd               respjson.Field
		DurationMs        respjson.Field
		ExitCode          respjson.Field
		AgentID           respjson.Field
		Model             respjson.Field
		ReasoningEffort   respjson.Field
		RecipientAgentIDs respjson.Field
		raw               string
	} `json:"-"`
}

// anyAgentSessionItem is implemented by each variant of [AgentSessionItemUnion] to
// add type safety for the return type of [AgentSessionItemUnion.AsAny]
type anyAgentSessionItem interface {
	implAgentSessionItemUnion()
}

func (AgentSessionMessage) implAgentSessionItemUnion()                {}
func (AgentReasoningItem) implAgentSessionItemUnion()                 {}
func (AgentFunctionCallItem) implAgentSessionItemUnion()              {}
func (AgentSessionItemFunctionCallOutput) implAgentSessionItemUnion() {}
func (AgentSessionItemAgentMessage) implAgentSessionItemUnion()       {}
func (AgentMcpCallItem) implAgentSessionItemUnion()                   {}
func (AgentWebSearchCallItem) implAgentSessionItemUnion()             {}
func (AgentCommandExecutionItem) implAgentSessionItemUnion()          {}
func (AgentCreateSubagentCallItem) implAgentSessionItemUnion()        {}
func (AgentSendSubagentInputCallItem) implAgentSessionItemUnion()     {}
func (AgentResumeSubagentCallItem) implAgentSessionItemUnion()        {}
func (AgentWaitForSubagentsCallItem) implAgentSessionItemUnion()      {}
func (AgentInterruptSubagentCallItem) implAgentSessionItemUnion()     {}
func (AgentCloseSubagentCallItem) implAgentSessionItemUnion()         {}

// Use the following switch statement to find the correct variant
//
//	switch variant := AgentSessionItemUnion.AsAny().(type) {
//	case openai.AgentSessionMessage:
//	case openai.AgentReasoningItem:
//	case openai.AgentFunctionCallItem:
//	case openai.AgentSessionItemFunctionCallOutput:
//	case openai.AgentSessionItemAgentMessage:
//	case openai.AgentMcpCallItem:
//	case openai.AgentWebSearchCallItem:
//	case openai.AgentCommandExecutionItem:
//	case openai.AgentCreateSubagentCallItem:
//	case openai.AgentSendSubagentInputCallItem:
//	case openai.AgentResumeSubagentCallItem:
//	case openai.AgentWaitForSubagentsCallItem:
//	case openai.AgentInterruptSubagentCallItem:
//	case openai.AgentCloseSubagentCallItem:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u AgentSessionItemUnion) AsAny() anyAgentSessionItem {
	switch u.Type {
	case "message":
		return u.AsMessage()
	case "reasoning":
		return u.AsReasoning()
	case "function_call":
		return u.AsFunctionCall()
	case "function_call_output":
		return u.AsFunctionCallOutput()
	case "agent_message":
		return u.AsAgentMessage()
	case "mcp_call":
		return u.AsMcpCall()
	case "web_search_call":
		return u.AsWebSearchCall()
	case "command_execution":
		return u.AsCommandExecution()
	case "create_subagent_call":
		return u.AsCreateSubagentCall()
	case "send_subagent_input_call":
		return u.AsSendSubagentInputCall()
	case "resume_subagent_call":
		return u.AsResumeSubagentCall()
	case "wait_for_subagents_call":
		return u.AsWaitForSubagentsCall()
	case "interrupt_subagent_call":
		return u.AsInterruptSubagentCall()
	case "close_subagent_call":
		return u.AsCloseSubagentCall()
	}
	return nil
}

func (u AgentSessionItemUnion) AsMessage() (v AgentSessionMessage) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentSessionItemUnion) AsReasoning() (v AgentReasoningItem) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentSessionItemUnion) AsFunctionCall() (v AgentFunctionCallItem) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentSessionItemUnion) AsFunctionCallOutput() (v AgentSessionItemFunctionCallOutput) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentSessionItemUnion) AsAgentMessage() (v AgentSessionItemAgentMessage) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentSessionItemUnion) AsMcpCall() (v AgentMcpCallItem) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentSessionItemUnion) AsWebSearchCall() (v AgentWebSearchCallItem) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentSessionItemUnion) AsCommandExecution() (v AgentCommandExecutionItem) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentSessionItemUnion) AsCreateSubagentCall() (v AgentCreateSubagentCallItem) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentSessionItemUnion) AsSendSubagentInputCall() (v AgentSendSubagentInputCallItem) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentSessionItemUnion) AsResumeSubagentCall() (v AgentResumeSubagentCallItem) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentSessionItemUnion) AsWaitForSubagentsCall() (v AgentWaitForSubagentsCallItem) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentSessionItemUnion) AsInterruptSubagentCall() (v AgentInterruptSubagentCallItem) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentSessionItemUnion) AsCloseSubagentCall() (v AgentCloseSubagentCallItem) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u AgentSessionItemUnion) RawJSON() string { return u.JSON.raw }

func (r *AgentSessionItemUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// AgentSessionItemUnionContent is an implicit subunion of [AgentSessionItemUnion].
// AgentSessionItemUnionContent provides convenient access to the sub-properties of
// the union.
//
// For type safety it is recommended to directly use a variant of the
// [AgentSessionItemUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfAgentSessionMessageContentArray OfAgentContentArray]
type AgentSessionItemUnionContent struct {
	// This field will be present if the value is a [[]AgentSessionMessageContentUnion]
	// instead of an object.
	OfAgentSessionMessageContentArray []AgentSessionMessageContentUnion `json:",inline"`
	// This field will be present if the value is a [[]AgentContentUnion] instead of an
	// object.
	OfAgentContentArray []AgentContentUnion `json:",inline"`
	JSON                struct {
		OfAgentSessionMessageContentArray respjson.Field
		OfAgentContentArray               respjson.Field
		raw                               string
	} `json:"-"`
}

func (r *AgentSessionItemUnionContent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// AgentSessionItemUnionError is an implicit subunion of [AgentSessionItemUnion].
// AgentSessionItemUnionError provides convenient access to the sub-properties of
// the union.
//
// For type safety it is recommended to directly use a variant of the
// [AgentSessionItemUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfString OfAgentMcpCallItemError]
type AgentSessionItemUnionError struct {
	// This field will be present if the value is a [string] instead of an object.
	OfString string `json:",inline"`
	// This field will be present if the value is a [any] instead of an object.
	OfAgentMcpCallItemError any `json:",inline"`
	JSON                    struct {
		OfString                respjson.Field
		OfAgentMcpCallItemError respjson.Field
		raw                     string
	} `json:"-"`
}

func (r *AgentSessionItemUnionError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// AgentSessionItemUnionOutput is an implicit subunion of [AgentSessionItemUnion].
// AgentSessionItemUnionOutput provides convenient access to the sub-properties of
// the union.
//
// For type safety it is recommended to directly use a variant of the
// [AgentSessionItemUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfString OfInputContentArray OfAgentMcpCallItemOutput]
type AgentSessionItemUnionOutput struct {
	// This field will be present if the value is a [string] instead of an object.
	OfString string `json:",inline"`
	// This field will be present if the value is a [[]InputContentUnion] instead of an
	// object.
	OfInputContentArray []InputContentUnion `json:",inline"`
	// This field will be present if the value is a [any] instead of an object.
	OfAgentMcpCallItemOutput any `json:",inline"`
	JSON                     struct {
		OfString                 respjson.Field
		OfInputContentArray      respjson.Field
		OfAgentMcpCallItemOutput respjson.Field
		raw                      string
	} `json:"-"`
}

func (r *AgentSessionItemUnionOutput) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The result supplied for a function call.
type AgentSessionItemFunctionCallOutput struct {
	// The ID of the function call output item.
	ID string `json:"id" api:"required"`
	// The ID of the function call that produced this output.
	CallID string `json:"call_id" api:"required"`
	// The error message, if the call failed.
	Error string `json:"error" api:"required"`
	// The text or model-input content supplied as a function result.
	Output AgentFunctionCallOutputUnion `json:"output" api:"required"`
	// The status of the function call.
	//
	// Any of "in_progress", "completed", "failed", "incomplete".
	Status AgentFunctionCallStatus `json:"status" api:"required"`
	// The ID of the turn that contains this item.
	TurnID string `json:"turn_id" api:"required"`
	// The item type. Always `function_call_output`.
	Type constant.FunctionCallOutput `json:"type" default:"function_call_output"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CallID      respjson.Field
		Error       respjson.Field
		Output      respjson.Field
		Status      respjson.Field
		TurnID      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentSessionItemFunctionCallOutput) RawJSON() string { return r.JSON.raw }
func (r *AgentSessionItemFunctionCallOutput) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A message exchanged between agent threads.
type AgentSessionItemAgentMessage struct {
	// The ID of the message.
	ID string `json:"id" api:"required"`
	// The content exchanged between the agents.
	Content []AgentContentUnion `json:"content" api:"required"`
	// The ID or name of the receiving agent.
	RecipientAgentID string `json:"recipient_agent_id" api:"required"`
	// The ID or name of the sending agent.
	SenderAgentID string `json:"sender_agent_id" api:"required"`
	// The ID of the turn that contains this item.
	TurnID string `json:"turn_id" api:"required"`
	// The item type. Always `agent_message`.
	Type constant.AgentMessage `json:"type" default:"agent_message"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID               respjson.Field
		Content          respjson.Field
		RecipientAgentID respjson.Field
		SenderAgentID    respjson.Field
		TurnID           respjson.Field
		Type             respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentSessionItemAgentMessage) RawJSON() string { return r.JSON.raw }
func (r *AgentSessionItemAgentMessage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A user or assistant message recorded in a session.
type AgentSessionMessage struct {
	// The ID of this item, or null for legacy user messages whose ID was not recorded.
	ID string `json:"id" api:"required"`
	// The content of the message. User messages contain input text or images;
	// assistant messages contain output text.
	Content []AgentSessionMessageContentUnion `json:"content" api:"required"`
	// The phase of an assistant message.
	//
	// Any of "commentary", "final_answer".
	Phase AgentSessionMessagePhase `json:"phase" api:"required"`
	// The role of the message author.
	//
	// Any of "user", "assistant".
	Role AgentSessionMessageRole `json:"role" api:"required"`
	// The status of the message. User messages are always `completed`.
	//
	// Any of "in_progress", "completed", "incomplete".
	Status AgentOutputItemStatus `json:"status" api:"required"`
	// The ID of the turn that contains this item.
	TurnID string `json:"turn_id" api:"required"`
	// The item type. Always `message`.
	Type constant.Message `json:"type" default:"message"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Content     respjson.Field
		Phase       respjson.Field
		Role        respjson.Field
		Status      respjson.Field
		TurnID      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentSessionMessage) RawJSON() string { return r.JSON.raw }
func (r *AgentSessionMessage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The phase of an assistant message.
type AgentSessionMessagePhase string

const (
	AgentSessionMessagePhaseCommentary  AgentSessionMessagePhase = "commentary"
	AgentSessionMessagePhaseFinalAnswer AgentSessionMessagePhase = "final_answer"
)

// The role of the message author.
type AgentSessionMessageRole string

const (
	AgentSessionMessageRoleUser      AgentSessionMessageRole = "user"
	AgentSessionMessageRoleAssistant AgentSessionMessageRole = "assistant"
)

// AgentSessionMessageContentUnion contains all possible properties and values from
// [AgentSessionMessageContentInputText], [AgentSessionMessageContentInputImage],
// [AgentSessionMessageContentOutputText].
//
// Use the [AgentSessionMessageContentUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type AgentSessionMessageContentUnion struct {
	Text string `json:"text"`
	// Any of "input_text", "input_image", "output_text".
	Type string `json:"type"`
	// This field is from variant [AgentSessionMessageContentInputImage].
	ImageURL string `json:"image_url"`
	JSON     struct {
		Text     respjson.Field
		Type     respjson.Field
		ImageURL respjson.Field
		raw      string
	} `json:"-"`
}

// anyAgentSessionMessageContent is implemented by each variant of
// [AgentSessionMessageContentUnion] to add type safety for the return type of
// [AgentSessionMessageContentUnion.AsAny]
type anyAgentSessionMessageContent interface {
	implAgentSessionMessageContentUnion()
}

func (AgentSessionMessageContentInputText) implAgentSessionMessageContentUnion()  {}
func (AgentSessionMessageContentInputImage) implAgentSessionMessageContentUnion() {}
func (AgentSessionMessageContentOutputText) implAgentSessionMessageContentUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := AgentSessionMessageContentUnion.AsAny().(type) {
//	case openai.AgentSessionMessageContentInputText:
//	case openai.AgentSessionMessageContentInputImage:
//	case openai.AgentSessionMessageContentOutputText:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u AgentSessionMessageContentUnion) AsAny() anyAgentSessionMessageContent {
	switch u.Type {
	case "input_text":
		return u.AsInputText()
	case "input_image":
		return u.AsInputImage()
	case "output_text":
		return u.AsOutputText()
	}
	return nil
}

func (u AgentSessionMessageContentUnion) AsInputText() (v AgentSessionMessageContentInputText) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentSessionMessageContentUnion) AsInputImage() (v AgentSessionMessageContentInputImage) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentSessionMessageContentUnion) AsOutputText() (v AgentSessionMessageContentOutputText) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u AgentSessionMessageContentUnion) RawJSON() string { return u.JSON.raw }

func (r *AgentSessionMessageContentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Text supplied by the user.
type AgentSessionMessageContentInputText struct {
	// The text supplied by the user.
	Text string `json:"text" api:"required"`
	// The type of the object. Always `input_text`.
	Type constant.InputText `json:"type" default:"input_text"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Text        respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentSessionMessageContentInputText) RawJSON() string { return r.JSON.raw }
func (r *AgentSessionMessageContentInputText) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// An image supplied by the user.
type AgentSessionMessageContentInputImage struct {
	// The URL of the image supplied by the user, which may be a base64-encoded data
	// URL.
	ImageURL string `json:"image_url" api:"required"`
	// The type of the object. Always `input_image`.
	Type constant.InputImage `json:"type" default:"input_image"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ImageURL    respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentSessionMessageContentInputImage) RawJSON() string { return r.JSON.raw }
func (r *AgentSessionMessageContentInputImage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Text produced by the assistant.
type AgentSessionMessageContentOutputText struct {
	// The text produced by the assistant.
	Text string `json:"text" api:"required"`
	// The type of the object. Always `output_text`.
	Type constant.OutputText `json:"type" default:"output_text"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Text        respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentSessionMessageContentOutputText) RawJSON() string { return r.JSON.raw }
func (r *AgentSessionMessageContentOutputText) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Emitted when a session is waiting for one or more required actions.
type AgentSessionRequiresActionEvent struct {
	// The unique ID of the event.
	EventID string `json:"event_id" api:"required"`
	// The session and its current required actions.
	Session AgentSession `json:"session" api:"required"`
	// The type of the object. Always `agent.session.requires_action`.
	Type constant.AgentSessionRequiresAction `json:"type" default:"agent.session.requires_action"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventID     respjson.Field
		Session     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentSessionRequiresActionEvent) RawJSON() string { return r.JSON.raw }
func (r *AgentSessionRequiresActionEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Emitted when a closed subagent successfully resumes.
type AgentSessionSubagentActiveEvent struct {
	// The unique ID of the event.
	EventID string `json:"event_id" api:"required"`
	// The subagent that resumed.
	Subagent Subagent `json:"subagent" api:"required"`
	// The type of the object. Always `agent.session.subagent.active`.
	Type constant.AgentSessionSubagentActive `json:"type" default:"agent.session.subagent.active"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventID     respjson.Field
		Subagent    respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentSessionSubagentActiveEvent) RawJSON() string { return r.JSON.raw }
func (r *AgentSessionSubagentActiveEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Emitted when a subagent is closed.
type AgentSessionSubagentClosedEvent struct {
	// The unique ID of the event.
	EventID string `json:"event_id" api:"required"`
	// The subagent that was closed.
	Subagent Subagent `json:"subagent" api:"required"`
	// The type of the object. Always `agent.session.subagent.closed`.
	Type constant.AgentSessionSubagentClosed `json:"type" default:"agent.session.subagent.closed"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventID     respjson.Field
		Subagent    respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentSessionSubagentClosedEvent) RawJSON() string { return r.JSON.raw }
func (r *AgentSessionSubagentClosedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Emitted when a subagent is created.
type AgentSessionSubagentCreatedEvent struct {
	// The unique ID of the event.
	EventID string `json:"event_id" api:"required"`
	// The subagent that was created.
	Subagent Subagent `json:"subagent" api:"required"`
	// The type of the object. Always `agent.session.subagent.created`.
	Type constant.AgentSessionSubagentCreated `json:"type" default:"agent.session.subagent.created"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventID     respjson.Field
		Subagent    respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentSessionSubagentCreatedEvent) RawJSON() string { return r.JSON.raw }
func (r *AgentSessionSubagentCreatedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Emitted when a turn is cancelled.
type AgentSessionTurnCancelledEvent struct {
	// The unique ID of the event.
	EventID string `json:"event_id" api:"required"`
	// The ID of the session associated with the event.
	SessionID string `json:"session_id" api:"required"`
	// The cancelled turn.
	Turn Turn `json:"turn" api:"required"`
	// The ID of the turn associated with the event.
	TurnID string `json:"turn_id" api:"required"`
	// The type of the object. Always `agent.session.turn.cancelled`.
	Type constant.AgentSessionTurnCancelled `json:"type" default:"agent.session.turn.cancelled"`
	// Recorded token usage for a session or turn. Usage is best effort and may change.
	Usage TokenUsage `json:"usage" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventID     respjson.Field
		SessionID   respjson.Field
		Turn        respjson.Field
		TurnID      respjson.Field
		Type        respjson.Field
		Usage       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentSessionTurnCancelledEvent) RawJSON() string { return r.JSON.raw }
func (r *AgentSessionTurnCancelledEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Emitted when a turn completes.
type AgentSessionTurnCompletedEvent struct {
	// The unique ID of the event.
	EventID string `json:"event_id" api:"required"`
	// The ID of the session associated with the event.
	SessionID string `json:"session_id" api:"required"`
	// The completed turn.
	Turn Turn `json:"turn" api:"required"`
	// The ID of the turn associated with the event.
	TurnID string `json:"turn_id" api:"required"`
	// The type of the object. Always `agent.session.turn.completed`.
	Type constant.AgentSessionTurnCompleted `json:"type" default:"agent.session.turn.completed"`
	// Recorded token usage for a session or turn. Usage is best effort and may change.
	Usage TokenUsage `json:"usage" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventID     respjson.Field
		SessionID   respjson.Field
		Turn        respjson.Field
		TurnID      respjson.Field
		Type        respjson.Field
		Usage       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentSessionTurnCompletedEvent) RawJSON() string { return r.JSON.raw }
func (r *AgentSessionTurnCompletedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Emitted when an output text content part is added.
type AgentSessionTurnContentPartAddedEvent struct {
	// The index of the content part in the message.
	ContentIndex int64 `json:"content_index" api:"required"`
	// The unique ID of the event.
	EventID string `json:"event_id" api:"required"`
	// The ID of the message item.
	ItemID string `json:"item_id" api:"required"`
	// The index of the item in the turn output.
	OutputIndex int64 `json:"output_index" api:"required"`
	// The initial content part.
	Part OutputText `json:"part" api:"required"`
	// The ID of the session associated with the event.
	SessionID string `json:"session_id" api:"required"`
	// The ID of the turn associated with the event, when applicable.
	TurnID string `json:"turn_id" api:"required"`
	// The type of the object. Always `agent.session.turn.content_part.added`.
	Type constant.AgentSessionTurnContentPartAdded `json:"type" default:"agent.session.turn.content_part.added"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ContentIndex respjson.Field
		EventID      respjson.Field
		ItemID       respjson.Field
		OutputIndex  respjson.Field
		Part         respjson.Field
		SessionID    respjson.Field
		TurnID       respjson.Field
		Type         respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentSessionTurnContentPartAddedEvent) RawJSON() string { return r.JSON.raw }
func (r *AgentSessionTurnContentPartAddedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Emitted when an output content part is complete.
type AgentSessionTurnContentPartDoneEvent struct {
	// The index of the content part in the message.
	ContentIndex int64 `json:"content_index" api:"required"`
	// The unique ID of the event.
	EventID string `json:"event_id" api:"required"`
	// The ID of the message item.
	ItemID string `json:"item_id" api:"required"`
	// The index of the item in the turn output.
	OutputIndex int64 `json:"output_index" api:"required"`
	// The completed content part.
	Part OutputText `json:"part" api:"required"`
	// The ID of the session associated with the event.
	SessionID string `json:"session_id" api:"required"`
	// The ID of the turn associated with the event, when applicable.
	TurnID string `json:"turn_id" api:"required"`
	// The type of the object. Always `agent.session.turn.content_part.done`.
	Type constant.AgentSessionTurnContentPartDone `json:"type" default:"agent.session.turn.content_part.done"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ContentIndex respjson.Field
		EventID      respjson.Field
		ItemID       respjson.Field
		OutputIndex  respjson.Field
		Part         respjson.Field
		SessionID    respjson.Field
		TurnID       respjson.Field
		Type         respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentSessionTurnContentPartDoneEvent) RawJSON() string { return r.JSON.raw }
func (r *AgentSessionTurnContentPartDoneEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Emitted when a turn is created.
type AgentSessionTurnCreatedEvent struct {
	// The unique ID of the event.
	EventID string `json:"event_id" api:"required"`
	// The ID of the session associated with the event.
	SessionID string `json:"session_id" api:"required"`
	// The turn at the time it was created.
	Turn Turn `json:"turn" api:"required"`
	// The ID of the turn associated with the event.
	TurnID string `json:"turn_id" api:"required"`
	// The type of the object. Always `agent.session.turn.created`.
	Type constant.AgentSessionTurnCreated `json:"type" default:"agent.session.turn.created"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventID     respjson.Field
		SessionID   respjson.Field
		Turn        respjson.Field
		TurnID      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentSessionTurnCreatedEvent) RawJSON() string { return r.JSON.raw }
func (r *AgentSessionTurnCreatedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Emitted when a turn fails.
type AgentSessionTurnFailedEvent struct {
	// The unique ID of the event.
	EventID string `json:"event_id" api:"required"`
	// The ID of the session associated with the event.
	SessionID string `json:"session_id" api:"required"`
	// The failed turn.
	Turn Turn `json:"turn" api:"required"`
	// The ID of the turn associated with the event.
	TurnID string `json:"turn_id" api:"required"`
	// The type of the object. Always `agent.session.turn.failed`.
	Type constant.AgentSessionTurnFailed `json:"type" default:"agent.session.turn.failed"`
	// Recorded token usage for a session or turn. Usage is best effort and may change.
	Usage TokenUsage `json:"usage" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventID     respjson.Field
		SessionID   respjson.Field
		Turn        respjson.Field
		TurnID      respjson.Field
		Type        respjson.Field
		Usage       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentSessionTurnFailedEvent) RawJSON() string { return r.JSON.raw }
func (r *AgentSessionTurnFailedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Emitted when a turn starts running.
type AgentSessionTurnInProgressEvent struct {
	// The unique ID of the event.
	EventID string `json:"event_id" api:"required"`
	// The ID of the session associated with the event.
	SessionID string `json:"session_id" api:"required"`
	// The turn at the time it started running.
	Turn Turn `json:"turn" api:"required"`
	// The ID of the turn associated with the event.
	TurnID string `json:"turn_id" api:"required"`
	// The type of the object. Always `agent.session.turn.in_progress`.
	Type constant.AgentSessionTurnInProgress `json:"type" default:"agent.session.turn.in_progress"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventID     respjson.Field
		SessionID   respjson.Field
		Turn        respjson.Field
		TurnID      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentSessionTurnInProgressEvent) RawJSON() string { return r.JSON.raw }
func (r *AgentSessionTurnInProgressEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Emitted when an item is added to a turn.
type AgentSessionTurnItemAddedEvent struct {
	// The unique ID of the event.
	EventID string `json:"event_id" api:"required"`
	// The item that was added.
	Item AgentSessionItemUnion `json:"item" api:"required"`
	// The index of the item in the turn output, when the item is agent output.
	OutputIndex int64 `json:"output_index" api:"required"`
	// The ID of the session associated with the event.
	SessionID string `json:"session_id" api:"required"`
	// The ID of the turn associated with the event, when applicable.
	TurnID string `json:"turn_id" api:"required"`
	// The type of the object. Always `agent.session.turn.item.added`.
	Type constant.AgentSessionTurnItemAdded `json:"type" default:"agent.session.turn.item.added"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventID     respjson.Field
		Item        respjson.Field
		OutputIndex respjson.Field
		SessionID   respjson.Field
		TurnID      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentSessionTurnItemAddedEvent) RawJSON() string { return r.JSON.raw }
func (r *AgentSessionTurnItemAddedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Emitted when an output item is complete.
type AgentSessionTurnItemDoneEvent struct {
	// The unique ID of the event.
	EventID string `json:"event_id" api:"required"`
	// The completed output item.
	Item AgentOutputItemUnion `json:"item" api:"required"`
	// The index of the output item in the turn output.
	OutputIndex int64 `json:"output_index" api:"required"`
	// The ID of the session associated with the event.
	SessionID string `json:"session_id" api:"required"`
	// The ID of the turn associated with the event, when applicable.
	TurnID string `json:"turn_id" api:"required"`
	// The type of the object. Always `agent.session.turn.item.done`.
	Type constant.AgentSessionTurnItemDone `json:"type" default:"agent.session.turn.item.done"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventID     respjson.Field
		Item        respjson.Field
		OutputIndex respjson.Field
		SessionID   respjson.Field
		TurnID      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentSessionTurnItemDoneEvent) RawJSON() string { return r.JSON.raw }
func (r *AgentSessionTurnItemDoneEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Emitted when text is appended to an output text content part.
type AgentSessionTurnOutputTextDeltaEvent struct {
	// The index of the content part in the message.
	ContentIndex int64 `json:"content_index" api:"required"`
	// The text that was appended.
	Delta string `json:"delta" api:"required"`
	// The unique ID of the event.
	EventID string `json:"event_id" api:"required"`
	// The ID of the message item.
	ItemID string `json:"item_id" api:"required"`
	// The index of the item in the turn output.
	OutputIndex int64 `json:"output_index" api:"required"`
	// The ID of the session associated with the event.
	SessionID string `json:"session_id" api:"required"`
	// The ID of the turn associated with the event, when applicable.
	TurnID string `json:"turn_id" api:"required"`
	// The type of the object. Always `agent.session.turn.output_text.delta`.
	Type constant.AgentSessionTurnOutputTextDelta `json:"type" default:"agent.session.turn.output_text.delta"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ContentIndex respjson.Field
		Delta        respjson.Field
		EventID      respjson.Field
		ItemID       respjson.Field
		OutputIndex  respjson.Field
		SessionID    respjson.Field
		TurnID       respjson.Field
		Type         respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentSessionTurnOutputTextDeltaEvent) RawJSON() string { return r.JSON.raw }
func (r *AgentSessionTurnOutputTextDeltaEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Emitted when an output text content part is complete.
type AgentSessionTurnOutputTextDoneEvent struct {
	// The index of the content part in the message.
	ContentIndex int64 `json:"content_index" api:"required"`
	// The unique ID of the event.
	EventID string `json:"event_id" api:"required"`
	// The ID of the message item.
	ItemID string `json:"item_id" api:"required"`
	// The index of the item in the turn output.
	OutputIndex int64 `json:"output_index" api:"required"`
	// The ID of the session associated with the event.
	SessionID string `json:"session_id" api:"required"`
	// The complete output text.
	Text string `json:"text" api:"required"`
	// The ID of the turn associated with the event, when applicable.
	TurnID string `json:"turn_id" api:"required"`
	// The type of the object. Always `agent.session.turn.output_text.done`.
	Type constant.AgentSessionTurnOutputTextDone `json:"type" default:"agent.session.turn.output_text.done"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ContentIndex respjson.Field
		EventID      respjson.Field
		ItemID       respjson.Field
		OutputIndex  respjson.Field
		SessionID    respjson.Field
		Text         respjson.Field
		TurnID       respjson.Field
		Type         respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentSessionTurnOutputTextDoneEvent) RawJSON() string { return r.JSON.raw }
func (r *AgentSessionTurnOutputTextDoneEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Emitted when a reasoning summary content part is added.
type AgentSessionTurnReasoningSummaryPartAddedEvent struct {
	// The unique ID of the event.
	EventID string `json:"event_id" api:"required"`
	// The ID of the reasoning item.
	ItemID string `json:"item_id" api:"required"`
	// The index of the item in the turn output.
	OutputIndex int64 `json:"output_index" api:"required"`
	// The initial summary part.
	Part SummaryText `json:"part" api:"required"`
	// The ID of the session associated with the event.
	SessionID string `json:"session_id" api:"required"`
	// The index of the summary content part.
	SummaryIndex int64 `json:"summary_index" api:"required"`
	// The ID of the turn associated with the event, when applicable.
	TurnID string `json:"turn_id" api:"required"`
	// The type of the object. Always
	// `agent.session.turn.reasoning_summary_part.added`.
	Type constant.AgentSessionTurnReasoningSummaryPartAdded `json:"type" default:"agent.session.turn.reasoning_summary_part.added"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventID      respjson.Field
		ItemID       respjson.Field
		OutputIndex  respjson.Field
		Part         respjson.Field
		SessionID    respjson.Field
		SummaryIndex respjson.Field
		TurnID       respjson.Field
		Type         respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentSessionTurnReasoningSummaryPartAddedEvent) RawJSON() string { return r.JSON.raw }
func (r *AgentSessionTurnReasoningSummaryPartAddedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Emitted when a reasoning summary part is complete.
type AgentSessionTurnReasoningSummaryPartDoneEvent struct {
	// The unique ID of the event.
	EventID string `json:"event_id" api:"required"`
	// The ID of the reasoning item.
	ItemID string `json:"item_id" api:"required"`
	// The index of the item in the turn output.
	OutputIndex int64 `json:"output_index" api:"required"`
	// The completed summary part.
	Part SummaryText `json:"part" api:"required"`
	// The ID of the session associated with the event.
	SessionID string `json:"session_id" api:"required"`
	// Present as `incomplete` when summary generation was interrupted.
	Status constant.Incomplete `json:"status" default:"incomplete"`
	// The index of the summary part.
	SummaryIndex int64 `json:"summary_index" api:"required"`
	// The ID of the turn associated with the event, when applicable.
	TurnID string `json:"turn_id" api:"required"`
	// The type of the object. Always `agent.session.turn.reasoning_summary_part.done`.
	Type constant.AgentSessionTurnReasoningSummaryPartDone `json:"type" default:"agent.session.turn.reasoning_summary_part.done"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventID      respjson.Field
		ItemID       respjson.Field
		OutputIndex  respjson.Field
		Part         respjson.Field
		SessionID    respjson.Field
		Status       respjson.Field
		SummaryIndex respjson.Field
		TurnID       respjson.Field
		Type         respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentSessionTurnReasoningSummaryPartDoneEvent) RawJSON() string { return r.JSON.raw }
func (r *AgentSessionTurnReasoningSummaryPartDoneEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Emitted when text is appended to a reasoning summary.
type AgentSessionTurnReasoningSummaryTextDeltaEvent struct {
	// The summary text that was appended.
	Delta string `json:"delta" api:"required"`
	// The unique ID of the event.
	EventID string `json:"event_id" api:"required"`
	// The ID of the reasoning item.
	ItemID string `json:"item_id" api:"required"`
	// The index of the item in the turn output.
	OutputIndex int64 `json:"output_index" api:"required"`
	// The ID of the session associated with the event.
	SessionID string `json:"session_id" api:"required"`
	// The index of the summary content part.
	SummaryIndex int64 `json:"summary_index" api:"required"`
	// The ID of the turn associated with the event, when applicable.
	TurnID string `json:"turn_id" api:"required"`
	// The type of the object. Always
	// `agent.session.turn.reasoning_summary_text.delta`.
	Type constant.AgentSessionTurnReasoningSummaryTextDelta `json:"type" default:"agent.session.turn.reasoning_summary_text.delta"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Delta        respjson.Field
		EventID      respjson.Field
		ItemID       respjson.Field
		OutputIndex  respjson.Field
		SessionID    respjson.Field
		SummaryIndex respjson.Field
		TurnID       respjson.Field
		Type         respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentSessionTurnReasoningSummaryTextDeltaEvent) RawJSON() string { return r.JSON.raw }
func (r *AgentSessionTurnReasoningSummaryTextDeltaEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Emitted when a reasoning summary content part is complete.
type AgentSessionTurnReasoningSummaryTextDoneEvent struct {
	// The unique ID of the event.
	EventID string `json:"event_id" api:"required"`
	// The ID of the reasoning item.
	ItemID string `json:"item_id" api:"required"`
	// The index of the item in the turn output.
	OutputIndex int64 `json:"output_index" api:"required"`
	// The ID of the session associated with the event.
	SessionID string `json:"session_id" api:"required"`
	// The index of the summary content part.
	SummaryIndex int64 `json:"summary_index" api:"required"`
	// The complete reasoning summary text.
	Text string `json:"text" api:"required"`
	// The ID of the turn associated with the event, when applicable.
	TurnID string `json:"turn_id" api:"required"`
	// The type of the object. Always `agent.session.turn.reasoning_summary_text.done`.
	Type constant.AgentSessionTurnReasoningSummaryTextDone `json:"type" default:"agent.session.turn.reasoning_summary_text.done"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventID      respjson.Field
		ItemID       respjson.Field
		OutputIndex  respjson.Field
		SessionID    respjson.Field
		SummaryIndex respjson.Field
		Text         respjson.Field
		TurnID       respjson.Field
		Type         respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentSessionTurnReasoningSummaryTextDoneEvent) RawJSON() string { return r.JSON.raw }
func (r *AgentSessionTurnReasoningSummaryTextDoneEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The text configuration used by an agent.
type AgentText struct {
	// The effective output format. Defaults to ordinary text.
	Format TextFormatUnion `json:"format" api:"required"`
	// The amount of text produced by the agent. Defaults to `medium`.
	//
	// Any of "low", "medium", "high".
	Verbosity AgentTextVerbosity `json:"verbosity" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Format      respjson.Field
		Verbosity   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentText) RawJSON() string { return r.JSON.raw }
func (r *AgentText) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The amount of text produced by the agent. Defaults to `medium`.
type AgentTextVerbosity string

const (
	AgentTextVerbosityLow    AgentTextVerbosity = "low"
	AgentTextVerbosityMedium AgentTextVerbosity = "medium"
	AgentTextVerbosityHigh   AgentTextVerbosity = "high"
)

// Configuration for text generated by the agent.
type AgentTextParam struct {
	// The output format for generated text.
	Format TextFormatParamUnion `json:"format,omitzero"`
	// The amount of text the model should produce.
	//
	// Any of "low", "medium", "high".
	Verbosity AgentTextParamVerbosity `json:"verbosity,omitzero"`
	paramObj
}

func (r AgentTextParam) MarshalJSON() (data []byte, err error) {
	type shadow AgentTextParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AgentTextParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The amount of text the model should produce.
type AgentTextParamVerbosity string

const (
	AgentTextParamVerbosityLow    AgentTextParamVerbosity = "low"
	AgentTextParamVerbosityMedium AgentTextParamVerbosity = "medium"
	AgentTextParamVerbosityHigh   AgentTextParamVerbosity = "high"
)

// AgentToolUnion contains all possible properties and values from
// [AgentToolFunction], [AgentToolProgrammaticToolCalling], [AgentToolMcp],
// [AgentToolWebSearch].
//
// Use the [AgentToolUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type AgentToolUnion struct {
	// This field is from variant [AgentToolFunction].
	DeferLoading bool `json:"defer_loading"`
	// This field is from variant [AgentToolFunction].
	Description string `json:"description"`
	// This field is from variant [AgentToolFunction].
	Name string `json:"name"`
	// This field is from variant [AgentToolFunction].
	Parameters map[string]any `json:"parameters"`
	// Any of "function", "programmatic_tool_calling", "mcp", "web_search".
	Type string `json:"type"`
	// This field is from variant [AgentToolProgrammaticToolCalling].
	Enabled bool `json:"enabled"`
	// This field is from variant [AgentToolMcp].
	AllowedTools []string `json:"allowed_tools"`
	// This field is from variant [AgentToolMcp].
	ConnectionOrigin string `json:"connection_origin"`
	// This field is from variant [AgentToolMcp].
	CredentialID string `json:"credential_id"`
	// This field is from variant [AgentToolMcp].
	RequestMetadata map[string]any `json:"request_metadata"`
	// This field is from variant [AgentToolMcp].
	Required bool `json:"required"`
	// This field is from variant [AgentToolMcp].
	ServerLabel string `json:"server_label"`
	// This field is from variant [AgentToolMcp].
	Transport McpTransportUnion `json:"transport"`
	// This field is from variant [AgentToolWebSearch].
	AllowedDomains []string `json:"allowed_domains"`
	// This field is from variant [AgentToolWebSearch].
	ContextSize string `json:"context_size"`
	// This field is from variant [AgentToolWebSearch].
	Location AgentToolWebSearchLocation `json:"location"`
	// This field is from variant [AgentToolWebSearch].
	Mode string `json:"mode"`
	JSON struct {
		DeferLoading     respjson.Field
		Description      respjson.Field
		Name             respjson.Field
		Parameters       respjson.Field
		Type             respjson.Field
		Enabled          respjson.Field
		AllowedTools     respjson.Field
		ConnectionOrigin respjson.Field
		CredentialID     respjson.Field
		RequestMetadata  respjson.Field
		Required         respjson.Field
		ServerLabel      respjson.Field
		Transport        respjson.Field
		AllowedDomains   respjson.Field
		ContextSize      respjson.Field
		Location         respjson.Field
		Mode             respjson.Field
		raw              string
	} `json:"-"`
}

// anyAgentTool is implemented by each variant of [AgentToolUnion] to add type
// safety for the return type of [AgentToolUnion.AsAny]
type anyAgentTool interface {
	implAgentToolUnion()
}

func (AgentToolFunction) implAgentToolUnion()                {}
func (AgentToolProgrammaticToolCalling) implAgentToolUnion() {}
func (AgentToolMcp) implAgentToolUnion()                     {}
func (AgentToolWebSearch) implAgentToolUnion()               {}

// Use the following switch statement to find the correct variant
//
//	switch variant := AgentToolUnion.AsAny().(type) {
//	case openai.AgentToolFunction:
//	case openai.AgentToolProgrammaticToolCalling:
//	case openai.AgentToolMcp:
//	case openai.AgentToolWebSearch:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u AgentToolUnion) AsAny() anyAgentTool {
	switch u.Type {
	case "function":
		return u.AsFunction()
	case "programmatic_tool_calling":
		return u.AsProgrammaticToolCalling()
	case "mcp":
		return u.AsMcp()
	case "web_search":
		return u.AsWebSearch()
	}
	return nil
}

func (u AgentToolUnion) AsFunction() (v AgentToolFunction) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentToolUnion) AsProgrammaticToolCalling() (v AgentToolProgrammaticToolCalling) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentToolUnion) AsMcp() (v AgentToolMcp) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AgentToolUnion) AsWebSearch() (v AgentToolWebSearch) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u AgentToolUnion) RawJSON() string { return u.JSON.raw }

func (r *AgentToolUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A function defined by the application.
type AgentToolFunction struct {
	// Whether the function is deferred and discovered through tool search.
	DeferLoading bool `json:"defer_loading" api:"required"`
	// A description of what the function does.
	Description string `json:"description" api:"required"`
	// The name of the function.
	Name string `json:"name" api:"required"`
	// A JSON Schema object describing the function's arguments.
	Parameters map[string]any `json:"parameters" api:"required"`
	// The type of the object. Always `function`.
	Type constant.Function `json:"type" default:"function"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		DeferLoading respjson.Field
		Description  respjson.Field
		Name         respjson.Field
		Parameters   respjson.Field
		Type         respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentToolFunction) RawJSON() string { return r.JSON.raw }
func (r *AgentToolFunction) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Enables calling tools from model-generated code.
type AgentToolProgrammaticToolCalling struct {
	// Whether tools can be called from model-generated code.
	Enabled bool `json:"enabled" api:"required"`
	// The type of the object. Always `programmatic_tool_calling`.
	Type constant.ProgrammaticToolCalling `json:"type" default:"programmatic_tool_calling"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Enabled     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentToolProgrammaticToolCalling) RawJSON() string { return r.JSON.raw }
func (r *AgentToolProgrammaticToolCalling) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Tools provided by a remote MCP server.
type AgentToolMcp struct {
	// The MCP tools the agent may call.
	AllowedTools []string `json:"allowed_tools" api:"required"`
	// Where outbound MCP HTTP connections originate.
	//
	// Any of "service", "environment".
	ConnectionOrigin string `json:"connection_origin" api:"required"`
	// The attached vault credential selected for this MCP server, if any. Optional
	// when exactly one attached credential matches the server URL.
	CredentialID string `json:"credential_id" api:"required"`
	// Metadata included with requests to this MCP server.
	RequestMetadata map[string]any `json:"request_metadata" api:"required"`
	// Whether this MCP server must initialize before the first turn.
	Required bool `json:"required" api:"required"`
	// A label used to identify the MCP server in tool calls.
	ServerLabel string `json:"server_label" api:"required"`
	// The transport used to connect to the MCP server.
	Transport McpTransportUnion `json:"transport" api:"required"`
	// The type of the object. Always `mcp`.
	Type constant.Mcp `json:"type" default:"mcp"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		AllowedTools     respjson.Field
		ConnectionOrigin respjson.Field
		CredentialID     respjson.Field
		RequestMetadata  respjson.Field
		Required         respjson.Field
		ServerLabel      respjson.Field
		Transport        respjson.Field
		Type             respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentToolMcp) RawJSON() string { return r.JSON.raw }
func (r *AgentToolMcp) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Web search.
type AgentToolWebSearch struct {
	// Allowed search domains, or `null` when the search is unrestricted.
	AllowedDomains []string `json:"allowed_domains" api:"required"`
	// The amount of search context made available to the model. Defaults to `medium`.
	//
	// Any of "low", "medium", "high".
	ContextSize string `json:"context_size" api:"required"`
	// Approximate user location used to localize web search results.
	Location AgentToolWebSearchLocation `json:"location" api:"required"`
	// The source used for web search results.
	//
	// Any of "disabled", "cached", "live".
	Mode string `json:"mode" api:"required"`
	// The type of the object. Always `web_search`.
	Type constant.WebSearch `json:"type" default:"web_search"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		AllowedDomains respjson.Field
		ContextSize    respjson.Field
		Location       respjson.Field
		Mode           respjson.Field
		Type           respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentToolWebSearch) RawJSON() string { return r.JSON.raw }
func (r *AgentToolWebSearch) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Approximate user location used to localize web search results.
type AgentToolWebSearchLocation struct {
	// The city name.
	City string `json:"city" api:"required"`
	// The two-letter ISO country code, such as `US`.
	Country string `json:"country" api:"required"`
	// The region or state name.
	Region string `json:"region" api:"required"`
	// The IANA timezone, such as `America/Los_Angeles`.
	Timezone string `json:"timezone" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		City        respjson.Field
		Country     respjson.Field
		Region      respjson.Field
		Timezone    respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentToolWebSearchLocation) RawJSON() string { return r.JSON.raw }
func (r *AgentToolWebSearchLocation) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func AgentToolParamOfParamFunction(description string, name string, parameters map[string]any) AgentToolParamUnion {
	var paramFunction AgentToolParamFunction
	paramFunction.Description = description
	paramFunction.Name = name
	paramFunction.Parameters = parameters
	return AgentToolParamUnion{OfParamFunction: &paramFunction}
}

func AgentToolParamOfParamMcp[
	T McpTransportParamHTTP | McpTransportParamStdio,
](serverLabel string, transport T) AgentToolParamUnion {
	var paramMcp AgentToolParamMcp
	paramMcp.ServerLabel = serverLabel
	switch v := any(transport).(type) {
	case McpTransportParamHTTP:
		paramMcp.Transport.OfParamHTTP = &v
	case McpTransportParamStdio:
		paramMcp.Transport.OfParamStdio = &v
	}
	return AgentToolParamUnion{OfParamMcp: &paramMcp}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type AgentToolParamUnion struct {
	OfParamFunction                *AgentToolParamFunction                `json:",omitzero,inline"`
	OfParamToolSearch              *AgentToolParamToolSearch              `json:",omitzero,inline"`
	OfParamProgrammaticToolCalling *AgentToolParamProgrammaticToolCalling `json:",omitzero,inline"`
	OfParamMcp                     *AgentToolParamMcp                     `json:",omitzero,inline"`
	OfParamWebSearch               *AgentToolParamWebSearch               `json:",omitzero,inline"`
	paramUnion
}

func (u AgentToolParamUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfParamFunction,
		u.OfParamToolSearch,
		u.OfParamProgrammaticToolCalling,
		u.OfParamMcp,
		u.OfParamWebSearch)
}
func (u *AgentToolParamUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u AgentToolParamUnion) GetDescription() *string {
	if vt := u.OfParamFunction; vt != nil {
		return &vt.Description
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AgentToolParamUnion) GetName() *string {
	if vt := u.OfParamFunction; vt != nil {
		return &vt.Name
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AgentToolParamUnion) GetParameters() map[string]any {
	if vt := u.OfParamFunction; vt != nil {
		return vt.Parameters
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AgentToolParamUnion) GetDeferLoading() *bool {
	if vt := u.OfParamFunction; vt != nil && vt.DeferLoading.Valid() {
		return &vt.DeferLoading.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AgentToolParamUnion) GetEnabled() *bool {
	if vt := u.OfParamProgrammaticToolCalling; vt != nil && vt.Enabled.Valid() {
		return &vt.Enabled.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AgentToolParamUnion) GetServerLabel() *string {
	if vt := u.OfParamMcp; vt != nil {
		return &vt.ServerLabel
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AgentToolParamUnion) GetTransport() *McpTransportParamUnion {
	if vt := u.OfParamMcp; vt != nil {
		return &vt.Transport
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AgentToolParamUnion) GetAllowedTools() []string {
	if vt := u.OfParamMcp; vt != nil {
		return vt.AllowedTools
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AgentToolParamUnion) GetConnectionOrigin() *string {
	if vt := u.OfParamMcp; vt != nil {
		return &vt.ConnectionOrigin
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AgentToolParamUnion) GetCredentialID() *string {
	if vt := u.OfParamMcp; vt != nil && vt.CredentialID.Valid() {
		return &vt.CredentialID.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AgentToolParamUnion) GetRequestMetadata() map[string]any {
	if vt := u.OfParamMcp; vt != nil {
		return vt.RequestMetadata
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AgentToolParamUnion) GetRequired() *bool {
	if vt := u.OfParamMcp; vt != nil && vt.Required.Valid() {
		return &vt.Required.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AgentToolParamUnion) GetAllowedDomains() []string {
	if vt := u.OfParamWebSearch; vt != nil {
		return vt.AllowedDomains
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AgentToolParamUnion) GetContextSize() *string {
	if vt := u.OfParamWebSearch; vt != nil {
		return &vt.ContextSize
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AgentToolParamUnion) GetLocation() *AgentToolParamWebSearchLocation {
	if vt := u.OfParamWebSearch; vt != nil {
		return &vt.Location
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AgentToolParamUnion) GetMode() *string {
	if vt := u.OfParamWebSearch; vt != nil {
		return &vt.Mode
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AgentToolParamUnion) GetType() *string {
	if vt := u.OfParamFunction; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfParamToolSearch; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfParamProgrammaticToolCalling; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfParamMcp; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfParamWebSearch; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[AgentToolParamUnion](
		"type",
		apijson.Discriminator[AgentToolParamFunction]("function"),
		apijson.Discriminator[AgentToolParamToolSearch]("tool_search"),
		apijson.Discriminator[AgentToolParamProgrammaticToolCalling]("programmatic_tool_calling"),
		apijson.Discriminator[AgentToolParamMcp]("mcp"),
		apijson.Discriminator[AgentToolParamWebSearch]("web_search"),
	)
}

// A function defined by the application.
//
// The properties Description, Name, Parameters, Type are required.
type AgentToolParamFunction struct {
	// A description of what the function does.
	Description string `json:"description" api:"required"`
	// The name of the function.
	Name string `json:"name" api:"required"`
	// A JSON Schema object describing the function's arguments.
	Parameters map[string]any `json:"parameters,omitzero" api:"required"`
	// Whether this function is deferred and discovered through tool search. Defaults
	// to `false`.
	DeferLoading param.Opt[bool] `json:"defer_loading,omitzero"`
	// The type of the object. Always `function`.
	//
	// This field can be elided, and will marshal its zero value as "function".
	Type constant.Function `json:"type" default:"function"`
	paramObj
}

func (r AgentToolParamFunction) MarshalJSON() (data []byte, err error) {
	type shadow AgentToolParamFunction
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AgentToolParamFunction) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewAgentToolParamToolSearch() AgentToolParamToolSearch {
	return AgentToolParamToolSearch{
		Type: "tool_search",
	}
}

// Discovers deferred function tools and loads them into the model context.
//
// This struct has a constant value, construct it with
// [NewAgentToolParamToolSearch].
type AgentToolParamToolSearch struct {
	// The type of the object. Always `tool_search`.
	Type constant.ToolSearch `json:"type" default:"tool_search"`
	paramObj
}

func (r AgentToolParamToolSearch) MarshalJSON() (data []byte, err error) {
	type shadow AgentToolParamToolSearch
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AgentToolParamToolSearch) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Enables calling tools from model-generated code.
//
// The property Type is required.
type AgentToolParamProgrammaticToolCalling struct {
	// Whether tools can be called from model-generated code. Defaults to `true`.
	Enabled param.Opt[bool] `json:"enabled,omitzero"`
	// The type of the object. Always `programmatic_tool_calling`.
	//
	// This field can be elided, and will marshal its zero value as
	// "programmatic_tool_calling".
	Type constant.ProgrammaticToolCalling `json:"type" default:"programmatic_tool_calling"`
	paramObj
}

func (r AgentToolParamProgrammaticToolCalling) MarshalJSON() (data []byte, err error) {
	type shadow AgentToolParamProgrammaticToolCalling
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AgentToolParamProgrammaticToolCalling) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Tools provided by a remote MCP server.
//
// The properties ServerLabel, Transport, Type are required.
type AgentToolParamMcp struct {
	// A label used to identify the MCP server in tool calls.
	ServerLabel string `json:"server_label" api:"required"`
	// The transport used to connect to the MCP server.
	Transport McpTransportParamUnion `json:"transport,omitzero" api:"required"`
	// The attached vault credential used to authenticate this MCP server. Optional
	// when exactly one attached credential matches the server URL.
	CredentialID param.Opt[string] `json:"credential_id,omitzero"`
	// Whether this MCP server must initialize before the first turn. Defaults to
	// `false`.
	Required param.Opt[bool] `json:"required,omitzero"`
	// The MCP tools the agent may call. All server tools are allowed when omitted.
	AllowedTools []string `json:"allowed_tools,omitzero"`
	// Where outbound MCP HTTP connections originate.
	//
	// Any of "service", "environment".
	ConnectionOrigin string `json:"connection_origin,omitzero"`
	// Metadata included with requests to this MCP server.
	RequestMetadata map[string]any `json:"request_metadata,omitzero"`
	// The type of the object. Always `mcp`.
	//
	// This field can be elided, and will marshal its zero value as "mcp".
	Type constant.Mcp `json:"type" default:"mcp"`
	paramObj
}

func (r AgentToolParamMcp) MarshalJSON() (data []byte, err error) {
	type shadow AgentToolParamMcp
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AgentToolParamMcp) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[AgentToolParamMcp](
		"connection_origin", "service", "environment",
	)
}

// Web search.
//
// The property Type is required.
type AgentToolParamWebSearch struct {
	// Domains the search may include.
	AllowedDomains []string `json:"allowed_domains,omitzero"`
	// The amount of web search context made available to the model.
	//
	// Any of "low", "medium", "high".
	ContextSize string `json:"context_size,omitzero"`
	// Approximate user location used to localize web search results.
	Location AgentToolParamWebSearchLocation `json:"location,omitzero"`
	// The source used for web search results.
	//
	// Any of "disabled", "cached", "live".
	Mode string `json:"mode,omitzero"`
	// The type of the object. Always `web_search`.
	//
	// This field can be elided, and will marshal its zero value as "web_search".
	Type constant.WebSearch `json:"type" default:"web_search"`
	paramObj
}

func (r AgentToolParamWebSearch) MarshalJSON() (data []byte, err error) {
	type shadow AgentToolParamWebSearch
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AgentToolParamWebSearch) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[AgentToolParamWebSearch](
		"context_size", "low", "medium", "high",
	)
	apijson.RegisterFieldValidator[AgentToolParamWebSearch](
		"mode", "disabled", "cached", "live",
	)
}

// Approximate user location used to localize web search results.
type AgentToolParamWebSearchLocation struct {
	// The city name.
	City param.Opt[string] `json:"city,omitzero"`
	// The two-letter ISO country code, such as `US`.
	Country param.Opt[string] `json:"country,omitzero"`
	// The region or state name.
	Region param.Opt[string] `json:"region,omitzero"`
	// The IANA timezone, such as `America/Los_Angeles`.
	Timezone param.Opt[string] `json:"timezone,omitzero"`
	paramObj
}

func (r AgentToolParamWebSearchLocation) MarshalJSON() (data []byte, err error) {
	type shadow AgentToolParamWebSearchLocation
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AgentToolParamWebSearchLocation) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A request to wait for one or more subagents.
type AgentWaitForSubagentsCallItem struct {
	// The ID of the tool call item.
	ID string `json:"id" api:"required"`
	// The IDs of the agents to wait for.
	RecipientAgentIDs []string `json:"recipient_agent_ids" api:"required"`
	// The ID of the agent waiting for results.
	SenderAgentID string `json:"sender_agent_id" api:"required"`
	// The status of the tool call.
	//
	// Any of "in_progress", "completed", "failed", "incomplete".
	Status AgentFunctionCallStatus `json:"status" api:"required"`
	// The ID of the turn that contains this item.
	TurnID string `json:"turn_id" api:"required"`
	// The item type. Always `wait_for_subagents_call`.
	Type constant.WaitForSubagentsCall `json:"type" default:"wait_for_subagents_call"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                respjson.Field
		RecipientAgentIDs respjson.Field
		SenderAgentID     respjson.Field
		Status            respjson.Field
		TurnID            respjson.Field
		Type              respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentWaitForSubagentsCallItem) RawJSON() string { return r.JSON.raw }
func (r *AgentWaitForSubagentsCallItem) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A web search call produced by the agent.
type AgentWebSearchCallItem struct {
	// The ID of the web search call.
	ID string `json:"id" api:"required"`
	// An action performed by the web search tool.
	Action WebSearchActionUnion `json:"action" api:"required"`
	// The status of the web search call.
	//
	// Any of "in_progress", "completed", "incomplete".
	Status AgentOutputItemStatus `json:"status" api:"required"`
	// The ID of the turn that contains this item.
	TurnID string `json:"turn_id" api:"required"`
	// The item type. Always `web_search_call`.
	Type constant.WebSearchCall `json:"type" default:"web_search_call"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Action      respjson.Field
		Status      respjson.Field
		TurnID      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AgentWebSearchCallItem) RawJSON() string { return r.JSON.raw }
func (r *AgentWebSearchCallItem) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// EnvironmentUnion contains all possible properties and values from
// [EnvironmentNone], [EnvironmentOpenAIHosted], [EnvironmentSelfHosted].
//
// Use the [EnvironmentUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type EnvironmentUnion struct {
	// Any of "none", "openai_hosted", "self_hosted".
	Type                  string   `json:"type"`
	ID                    string   `json:"id"`
	CapabilityDirectories []string `json:"capability_directories"`
	// This field is from variant [EnvironmentOpenAIHosted].
	Files []HostedEnvironmentFileUnion `json:"files"`
	// This field is from variant [EnvironmentOpenAIHosted].
	Network EnvironmentOpenAIHostedNetwork `json:"network"`
	// This field is from variant [EnvironmentOpenAIHosted].
	Packages EnvironmentOpenAIHostedPackages `json:"packages"`
	// This field is from variant [EnvironmentOpenAIHosted].
	Plugins []HostedPlugin `json:"plugins"`
	// This field is from variant [EnvironmentOpenAIHosted].
	Skills []HostedSkillUnion `json:"skills"`
	// This field is from variant [EnvironmentSelfHosted].
	RemoteURL string `json:"remote_url"`
	// This field is from variant [EnvironmentSelfHosted].
	WorkspaceDirectory string `json:"workspace_directory"`
	JSON               struct {
		Type                  respjson.Field
		ID                    respjson.Field
		CapabilityDirectories respjson.Field
		Files                 respjson.Field
		Network               respjson.Field
		Packages              respjson.Field
		Plugins               respjson.Field
		Skills                respjson.Field
		RemoteURL             respjson.Field
		WorkspaceDirectory    respjson.Field
		raw                   string
	} `json:"-"`
}

// anyEnvironment is implemented by each variant of [EnvironmentUnion] to add type
// safety for the return type of [EnvironmentUnion.AsAny]
type anyEnvironment interface {
	implEnvironmentUnion()
}

func (EnvironmentNone) implEnvironmentUnion()         {}
func (EnvironmentOpenAIHosted) implEnvironmentUnion() {}
func (EnvironmentSelfHosted) implEnvironmentUnion()   {}

// Use the following switch statement to find the correct variant
//
//	switch variant := EnvironmentUnion.AsAny().(type) {
//	case openai.EnvironmentNone:
//	case openai.EnvironmentOpenAIHosted:
//	case openai.EnvironmentSelfHosted:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u EnvironmentUnion) AsAny() anyEnvironment {
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

func (u EnvironmentUnion) AsNone() (v EnvironmentNone) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u EnvironmentUnion) AsOpenAIHosted() (v EnvironmentOpenAIHosted) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u EnvironmentUnion) AsSelfHosted() (v EnvironmentSelfHosted) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u EnvironmentUnion) RawJSON() string { return u.JSON.raw }

func (r *EnvironmentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The session talks to CCA without selecting or provisioning an execution
// environment.
type EnvironmentNone struct {
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
func (r EnvironmentNone) RawJSON() string { return r.JSON.raw }
func (r *EnvironmentNone) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// An environment hosted by OpenAI.
type EnvironmentOpenAIHosted struct {
	// The public ID of the environment.
	ID string `json:"id" api:"required"`
	// Directories that contain capabilities exposed to the agent.
	CapabilityDirectories []string `json:"capability_directories" api:"required"`
	// Files available in the environment, excluding their contents.
	Files []HostedEnvironmentFileUnion `json:"files" api:"required"`
	// The effective network access policy for the environment.
	Network EnvironmentOpenAIHostedNetwork `json:"network" api:"required"`
	// Packages installed in the environment.
	Packages EnvironmentOpenAIHostedPackages `json:"packages" api:"required"`
	// Plugins installed in the environment, excluding their archive contents.
	Plugins []HostedPlugin `json:"plugins" api:"required"`
	// Skills installed in the environment, excluding their archive contents.
	Skills []HostedSkillUnion `json:"skills" api:"required"`
	// The type of the object. Always `openai_hosted`.
	Type constant.OpenAIHosted `json:"type" default:"openai_hosted"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                    respjson.Field
		CapabilityDirectories respjson.Field
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
func (r EnvironmentOpenAIHosted) RawJSON() string { return r.JSON.raw }
func (r *EnvironmentOpenAIHosted) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The effective network access policy for the environment.
type EnvironmentOpenAIHostedNetwork struct {
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
func (r EnvironmentOpenAIHostedNetwork) RawJSON() string { return r.JSON.raw }
func (r *EnvironmentOpenAIHostedNetwork) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Packages installed in the environment.
type EnvironmentOpenAIHostedPackages struct {
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
func (r EnvironmentOpenAIHostedPackages) RawJSON() string { return r.JSON.raw }
func (r *EnvironmentOpenAIHostedPackages) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// An environment hosted by the application.
type EnvironmentSelfHosted struct {
	// The public ID of the environment.
	ID string `json:"id" api:"required"`
	// Directories that contain capabilities exposed to the agent.
	CapabilityDirectories []string `json:"capability_directories" api:"required"`
	// Pass this URL unchanged to `codex exec-server --remote` when connecting this
	// environment.
	RemoteURL string `json:"remote_url" api:"required"`
	// The type of the object. Always `self_hosted`.
	Type constant.SelfHosted `json:"type" default:"self_hosted"`
	// The absolute project directory inside the environment. Defaults to `/workspace`.
	WorkspaceDirectory string `json:"workspace_directory" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                    respjson.Field
		CapabilityDirectories respjson.Field
		RemoteURL             respjson.Field
		Type                  respjson.Field
		WorkspaceDirectory    respjson.Field
		ExtraFields           map[string]respjson.Field
		raw                   string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EnvironmentSelfHosted) RawJSON() string { return r.JSON.raw }
func (r *EnvironmentSelfHosted) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func EnvironmentParamOfParamSelfHosted(workspaceDirectory string) EnvironmentParamUnion {
	var paramSelfHosted EnvironmentParamSelfHosted
	paramSelfHosted.WorkspaceDirectory = workspaceDirectory
	return EnvironmentParamUnion{OfParamSelfHosted: &paramSelfHosted}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type EnvironmentParamUnion struct {
	OfParamNone         *EnvironmentParamNone         `json:",omitzero,inline"`
	OfParamOpenAIHosted *EnvironmentParamOpenAIHosted `json:",omitzero,inline"`
	OfParamSelfHosted   *EnvironmentParamSelfHosted   `json:",omitzero,inline"`
	paramUnion
}

func (u EnvironmentParamUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfParamNone, u.OfParamOpenAIHosted, u.OfParamSelfHosted)
}
func (u *EnvironmentParamUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u EnvironmentParamUnion) GetEnv() map[string]string {
	if vt := u.OfParamOpenAIHosted; vt != nil {
		return vt.Env
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u EnvironmentParamUnion) GetEnvironmentTemplateID() *string {
	if vt := u.OfParamOpenAIHosted; vt != nil && vt.EnvironmentTemplateID.Valid() {
		return &vt.EnvironmentTemplateID.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u EnvironmentParamUnion) GetFiles() []HostedEnvironmentFileParamUnion {
	if vt := u.OfParamOpenAIHosted; vt != nil {
		return vt.Files
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u EnvironmentParamUnion) GetNetwork() *EnvironmentParamOpenAIHostedNetwork {
	if vt := u.OfParamOpenAIHosted; vt != nil {
		return &vt.Network
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u EnvironmentParamUnion) GetPackages() *EnvironmentParamOpenAIHostedPackages {
	if vt := u.OfParamOpenAIHosted; vt != nil {
		return &vt.Packages
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u EnvironmentParamUnion) GetPlugins() []HostedPluginParam {
	if vt := u.OfParamOpenAIHosted; vt != nil {
		return vt.Plugins
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u EnvironmentParamUnion) GetSetupCommands() []SetupCommandParam {
	if vt := u.OfParamOpenAIHosted; vt != nil {
		return vt.SetupCommands
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u EnvironmentParamUnion) GetSkills() []HostedSkillParamUnion {
	if vt := u.OfParamOpenAIHosted; vt != nil {
		return vt.Skills
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u EnvironmentParamUnion) GetWorkspaceDirectory() *string {
	if vt := u.OfParamSelfHosted; vt != nil {
		return &vt.WorkspaceDirectory
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u EnvironmentParamUnion) GetType() *string {
	if vt := u.OfParamNone; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfParamOpenAIHosted; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfParamSelfHosted; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's CapabilityDirectories property, if
// present.
func (u EnvironmentParamUnion) GetCapabilityDirectories() []string {
	if vt := u.OfParamOpenAIHosted; vt != nil {
		return vt.CapabilityDirectories
	} else if vt := u.OfParamSelfHosted; vt != nil {
		return vt.CapabilityDirectories
	}
	return nil
}

func init() {
	apijson.RegisterUnion[EnvironmentParamUnion](
		"type",
		apijson.Discriminator[EnvironmentParamNone]("none"),
		apijson.Discriminator[EnvironmentParamOpenAIHosted]("openai_hosted"),
		apijson.Discriminator[EnvironmentParamSelfHosted]("self_hosted"),
	)
}

func NewEnvironmentParamNone() EnvironmentParamNone {
	return EnvironmentParamNone{
		Type: "none",
	}
}

// Runs the agent without an execution environment.
//
// This struct has a constant value, construct it with [NewEnvironmentParamNone].
type EnvironmentParamNone struct {
	// The type of the object. Always `none`.
	Type constant.None `json:"type" default:"none"`
	paramObj
}

func (r EnvironmentParamNone) MarshalJSON() (data []byte, err error) {
	type shadow EnvironmentParamNone
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *EnvironmentParamNone) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// An OpenAI-hosted environment, optionally based on a reusable template.
//
// The property Type is required.
type EnvironmentParamOpenAIHosted struct {
	// A reusable hosted template applied before inline session configuration. Omitted
	// fields inherit the template; network overrides cannot broaden its policy.
	EnvironmentTemplateID param.Opt[string] `json:"environment_template_id,omitzero"`
	// Files available before the agent starts. Defaults to an empty list.
	Files []HostedEnvironmentFileParamUnion `json:"files,omitzero"`
	// Plugins provided as inline ZIP archives. Defaults to an empty list.
	Plugins []HostedPluginParam `json:"plugins,omitzero"`
	// Ordered, confidential setup commands. Command bodies are never returned.
	SetupCommands []SetupCommandParam `json:"setup_commands,omitzero"`
	// Skills referenced by ID or provided as inline ZIP archives. Defaults to an empty
	// list.
	Skills []HostedSkillParamUnion `json:"skills,omitzero"`
	// Directories that contain capabilities exposed to the agent. Defaults to an empty
	// list.
	CapabilityDirectories []string `json:"capability_directories,omitzero"`
	// Environment variables made available to the agent.
	Env map[string]string `json:"env,omitzero"`
	// Network access for an OpenAI-hosted environment.
	Network EnvironmentParamOpenAIHostedNetwork `json:"network,omitzero"`
	// Packages to install in an OpenAI-hosted environment.
	Packages EnvironmentParamOpenAIHostedPackages `json:"packages,omitzero"`
	// The type of the object. Always `openai_hosted`.
	//
	// This field can be elided, and will marshal its zero value as "openai_hosted".
	Type constant.OpenAIHosted `json:"type" default:"openai_hosted"`
	paramObj
}

func (r EnvironmentParamOpenAIHosted) MarshalJSON() (data []byte, err error) {
	type shadow EnvironmentParamOpenAIHosted
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *EnvironmentParamOpenAIHosted) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Network access for an OpenAI-hosted environment.
//
// The property Access is required.
type EnvironmentParamOpenAIHostedNetwork struct {
	// The environment's network access mode.
	//
	// Any of "enabled", "disabled", "restricted".
	Access string `json:"access,omitzero" api:"required"`
	// Domains the environment may access when network access is restricted.
	AllowedDomains []string `json:"allowed_domains,omitzero"`
	paramObj
}

func (r EnvironmentParamOpenAIHostedNetwork) MarshalJSON() (data []byte, err error) {
	type shadow EnvironmentParamOpenAIHostedNetwork
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *EnvironmentParamOpenAIHostedNetwork) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[EnvironmentParamOpenAIHostedNetwork](
		"access", "enabled", "disabled", "restricted",
	)
}

// Packages to install in an OpenAI-hosted environment.
type EnvironmentParamOpenAIHostedPackages struct {
	// npm packages to install globally. Defaults to an empty list.
	Npm []string `json:"npm,omitzero"`
	// Python packages to install. Defaults to an empty list.
	Python []string `json:"python,omitzero"`
	// System packages to install. Defaults to an empty list.
	System []string `json:"system,omitzero"`
	paramObj
}

func (r EnvironmentParamOpenAIHostedPackages) MarshalJSON() (data []byte, err error) {
	type shadow EnvironmentParamOpenAIHostedPackages
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *EnvironmentParamOpenAIHostedPackages) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// An application-hosted environment configured inline.
//
// The properties Type, WorkspaceDirectory are required.
type EnvironmentParamSelfHosted struct {
	// Absolute project directory inside the self-hosted environment.
	WorkspaceDirectory string `json:"workspace_directory" api:"required"`
	// Directories that contain capabilities exposed to the agent. Defaults to an empty
	// list.
	CapabilityDirectories []string `json:"capability_directories,omitzero"`
	// The type of the object. Always `self_hosted`.
	//
	// This field can be elided, and will marshal its zero value as "self_hosted".
	Type constant.SelfHosted `json:"type" default:"self_hosted"`
	paramObj
}

func (r EnvironmentParamSelfHosted) MarshalJSON() (data []byte, err error) {
	type shadow EnvironmentParamSelfHosted
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *EnvironmentParamSelfHosted) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// HostedEnvironmentFileUnion contains all possible properties and values from
// [HostedEnvironmentFileID], [HostedEnvironmentFileInline].
//
// Use the [HostedEnvironmentFileUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type HostedEnvironmentFileUnion struct {
	ID string `json:"id"`
	// This field is from variant [HostedEnvironmentFileID].
	FileID    string `json:"file_id"`
	Path      string `json:"path"`
	SizeBytes int64  `json:"size_bytes"`
	// Any of "file_id", "inline".
	Type string `json:"type"`
	JSON struct {
		ID        respjson.Field
		FileID    respjson.Field
		Path      respjson.Field
		SizeBytes respjson.Field
		Type      respjson.Field
		raw       string
	} `json:"-"`
}

// anyHostedEnvironmentFile is implemented by each variant of
// [HostedEnvironmentFileUnion] to add type safety for the return type of
// [HostedEnvironmentFileUnion.AsAny]
type anyHostedEnvironmentFile interface {
	implHostedEnvironmentFileUnion()
}

func (HostedEnvironmentFileID) implHostedEnvironmentFileUnion()     {}
func (HostedEnvironmentFileInline) implHostedEnvironmentFileUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := HostedEnvironmentFileUnion.AsAny().(type) {
//	case openai.HostedEnvironmentFileID:
//	case openai.HostedEnvironmentFileInline:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u HostedEnvironmentFileUnion) AsAny() anyHostedEnvironmentFile {
	switch u.Type {
	case "file_id":
		return u.AsFileID()
	case "inline":
		return u.AsInline()
	}
	return nil
}

func (u HostedEnvironmentFileUnion) AsFileID() (v HostedEnvironmentFileID) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u HostedEnvironmentFileUnion) AsInline() (v HostedEnvironmentFileInline) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u HostedEnvironmentFileUnion) RawJSON() string { return u.JSON.raw }

func (r *HostedEnvironmentFileUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A file supplied inline when the session was created.
type HostedEnvironmentFileInline struct {
	// The session-scoped ID of the file in the execution environment.
	ID string `json:"id" api:"required"`
	// The file's absolute path inside the environment.
	Path string `json:"path" api:"required"`
	// The decoded file size in bytes.
	SizeBytes int64 `json:"size_bytes" api:"required"`
	// The type of the object. Always `inline`.
	Type constant.Inline `json:"type" default:"inline"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Path        respjson.Field
		SizeBytes   respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r HostedEnvironmentFileInline) RawJSON() string { return r.JSON.raw }
func (r *HostedEnvironmentFileInline) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A file copied from the OpenAI Files API.
type HostedEnvironmentFileID struct {
	// The session-scoped ID of the file in the execution environment.
	ID string `json:"id" api:"required"`
	// The ID of the uploaded file.
	FileID string `json:"file_id" api:"required"`
	// The file's absolute path inside the environment.
	Path string `json:"path" api:"required"`
	// The decoded file size in bytes.
	SizeBytes int64 `json:"size_bytes" api:"required"`
	// The type of the object. Always `file_id`.
	Type constant.FileID `json:"type" default:"file_id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		FileID      respjson.Field
		Path        respjson.Field
		SizeBytes   respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r HostedEnvironmentFileID) RawJSON() string { return r.JSON.raw }
func (r *HostedEnvironmentFileID) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func HostedEnvironmentFileParamOfParamFileID(fileID string, path string) HostedEnvironmentFileParamUnion {
	var paramFileID HostedEnvironmentFileParamFileID
	paramFileID.FileID = fileID
	paramFileID.Path = path
	return HostedEnvironmentFileParamUnion{OfParamFileID: &paramFileID}
}

func HostedEnvironmentFileParamOfParamInline(data string, path string) HostedEnvironmentFileParamUnion {
	var paramInline HostedEnvironmentFileParamInline
	paramInline.Data = data
	paramInline.Path = path
	return HostedEnvironmentFileParamUnion{OfParamInline: &paramInline}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type HostedEnvironmentFileParamUnion struct {
	OfParamFileID *HostedEnvironmentFileParamFileID `json:",omitzero,inline"`
	OfParamInline *HostedEnvironmentFileParamInline `json:",omitzero,inline"`
	paramUnion
}

func (u HostedEnvironmentFileParamUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfParamFileID, u.OfParamInline)
}
func (u *HostedEnvironmentFileParamUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u HostedEnvironmentFileParamUnion) GetFileID() *string {
	if vt := u.OfParamFileID; vt != nil {
		return &vt.FileID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u HostedEnvironmentFileParamUnion) GetData() *string {
	if vt := u.OfParamInline; vt != nil {
		return &vt.Data
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u HostedEnvironmentFileParamUnion) GetPath() *string {
	if vt := u.OfParamFileID; vt != nil {
		return (*string)(&vt.Path)
	} else if vt := u.OfParamInline; vt != nil {
		return (*string)(&vt.Path)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u HostedEnvironmentFileParamUnion) GetType() *string {
	if vt := u.OfParamFileID; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfParamInline; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[HostedEnvironmentFileParamUnion](
		"type",
		apijson.Discriminator[HostedEnvironmentFileParamFileID]("file_id"),
		apijson.Discriminator[HostedEnvironmentFileParamInline]("inline"),
	)
}

// A file previously uploaded through the OpenAI Files API.
//
// The properties FileID, Path, Type are required.
type HostedEnvironmentFileParamFileID struct {
	// The ID of the uploaded file.
	FileID string `json:"file_id" api:"required"`
	// The absolute destination path inside `/workspace`.
	Path string `json:"path" api:"required"`
	// The type of the object. Always `file_id`.
	//
	// This field can be elided, and will marshal its zero value as "file_id".
	Type constant.FileID `json:"type" default:"file_id"`
	paramObj
}

func (r HostedEnvironmentFileParamFileID) MarshalJSON() (data []byte, err error) {
	type shadow HostedEnvironmentFileParamFileID
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *HostedEnvironmentFileParamFileID) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A file supplied directly as standard-base64 data.
//
// The properties Data, Path, Type are required.
type HostedEnvironmentFileParamInline struct {
	// The standard-base64-encoded file contents.
	Data string `json:"data" api:"required"`
	// The absolute destination path inside `/workspace`.
	Path string `json:"path" api:"required"`
	// The type of the object. Always `inline`.
	//
	// This field can be elided, and will marshal its zero value as "inline".
	Type constant.Inline `json:"type" default:"inline"`
	paramObj
}

func (r HostedEnvironmentFileParamInline) MarshalJSON() (data []byte, err error) {
	type shadow HostedEnvironmentFileParamInline
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *HostedEnvironmentFileParamInline) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A plugin installed from an inline ZIP archive.
type HostedPlugin struct {
	// The installed plugin description.
	Description string `json:"description" api:"required"`
	// The installed plugin name.
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
func (r HostedPlugin) RawJSON() string { return r.JSON.raw }
func (r *HostedPlugin) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Supplies a plugin ZIP directly in the session request.
//
// The properties Description, Name, Source, Type are required.
type HostedPluginParam struct {
	// The plugin description declared in `.codex-plugin/plugin.json`.
	Description string `json:"description" api:"required"`
	// The plugin name declared in `.codex-plugin/plugin.json`.
	Name string `json:"name" api:"required"`
	// Provides ZIP bytes encoded with standard base64.
	Source InlineCapabilitySourceParam `json:"source,omitzero" api:"required"`
	// The type of the object. Always `inline`.
	//
	// This field can be elided, and will marshal its zero value as "inline".
	Type constant.Inline `json:"type" default:"inline"`
	paramObj
}

func (r HostedPluginParam) MarshalJSON() (data []byte, err error) {
	type shadow HostedPluginParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *HostedPluginParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// HostedSkillUnion contains all possible properties and values from
// [HostedSkillReference], [HostedSkillInline].
//
// Use the [HostedSkillUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type HostedSkillUnion struct {
	Description string `json:"description"`
	Name        string `json:"name"`
	// This field is from variant [HostedSkillReference].
	SkillID string `json:"skill_id"`
	// Any of "skill_reference", "inline".
	Type string `json:"type"`
	// This field is from variant [HostedSkillReference].
	Version string `json:"version"`
	JSON    struct {
		Description respjson.Field
		Name        respjson.Field
		SkillID     respjson.Field
		Type        respjson.Field
		Version     respjson.Field
		raw         string
	} `json:"-"`
}

// anyHostedSkill is implemented by each variant of [HostedSkillUnion] to add type
// safety for the return type of [HostedSkillUnion.AsAny]
type anyHostedSkill interface {
	implHostedSkillUnion()
}

func (HostedSkillReference) implHostedSkillUnion() {}
func (HostedSkillInline) implHostedSkillUnion()    {}

// Use the following switch statement to find the correct variant
//
//	switch variant := HostedSkillUnion.AsAny().(type) {
//	case openai.HostedSkillReference:
//	case openai.HostedSkillInline:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u HostedSkillUnion) AsAny() anyHostedSkill {
	switch u.Type {
	case "skill_reference":
		return u.AsSkillReference()
	case "inline":
		return u.AsInline()
	}
	return nil
}

func (u HostedSkillUnion) AsSkillReference() (v HostedSkillReference) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u HostedSkillUnion) AsInline() (v HostedSkillInline) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u HostedSkillUnion) RawJSON() string { return u.JSON.raw }

func (r *HostedSkillUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A skill installed from an inline ZIP archive.
type HostedSkillInline struct {
	// The installed skill description.
	Description string `json:"description" api:"required"`
	// The installed skill name.
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
func (r HostedSkillInline) RawJSON() string { return r.JSON.raw }
func (r *HostedSkillInline) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func HostedSkillParamOfParamSkillReference(skillID string) HostedSkillParamUnion {
	var paramSkillReference HostedSkillParamSkillReference
	paramSkillReference.SkillID = skillID
	return HostedSkillParamUnion{OfParamSkillReference: &paramSkillReference}
}

func HostedSkillParamOfParamInline(description string, name string, source InlineCapabilitySourceParam) HostedSkillParamUnion {
	var paramInline HostedSkillParamInline
	paramInline.Description = description
	paramInline.Name = name
	paramInline.Source = source
	return HostedSkillParamUnion{OfParamInline: &paramInline}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type HostedSkillParamUnion struct {
	OfParamSkillReference *HostedSkillParamSkillReference `json:",omitzero,inline"`
	OfParamInline         *HostedSkillParamInline         `json:",omitzero,inline"`
	paramUnion
}

func (u HostedSkillParamUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfParamSkillReference, u.OfParamInline)
}
func (u *HostedSkillParamUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u HostedSkillParamUnion) GetSkillID() *string {
	if vt := u.OfParamSkillReference; vt != nil {
		return &vt.SkillID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u HostedSkillParamUnion) GetVersion() *string {
	if vt := u.OfParamSkillReference; vt != nil && vt.Version.Valid() {
		return &vt.Version.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u HostedSkillParamUnion) GetDescription() *string {
	if vt := u.OfParamInline; vt != nil {
		return &vt.Description
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u HostedSkillParamUnion) GetName() *string {
	if vt := u.OfParamInline; vt != nil {
		return &vt.Name
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u HostedSkillParamUnion) GetSource() *InlineCapabilitySourceParam {
	if vt := u.OfParamInline; vt != nil {
		return &vt.Source
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u HostedSkillParamUnion) GetType() *string {
	if vt := u.OfParamSkillReference; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfParamInline; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[HostedSkillParamUnion](
		"type",
		apijson.Discriminator[HostedSkillParamSkillReference]("skill_reference"),
		apijson.Discriminator[HostedSkillParamInline]("inline"),
	)
}

// References a skill uploaded through the Skills API.
//
// The properties SkillID, Type are required.
type HostedSkillParamSkillReference struct {
	// The ID of the skill created through `/v1/skills`.
	SkillID string `json:"skill_id" api:"required"`
	// The skill version, a positive integer or `latest`; omission selects the default.
	Version param.Opt[string] `json:"version,omitzero"`
	// The type of the object. Always `skill_reference`.
	//
	// This field can be elided, and will marshal its zero value as "skill_reference".
	Type constant.SkillReference `json:"type" default:"skill_reference"`
	paramObj
}

func (r HostedSkillParamSkillReference) MarshalJSON() (data []byte, err error) {
	type shadow HostedSkillParamSkillReference
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *HostedSkillParamSkillReference) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Supplies a skill ZIP directly in the session request.
//
// The properties Description, Name, Source, Type are required.
type HostedSkillParamInline struct {
	// The skill description declared in `SKILL.md`.
	Description string `json:"description" api:"required"`
	// The skill name declared in `SKILL.md`.
	Name string `json:"name" api:"required"`
	// Provides ZIP bytes encoded with standard base64.
	Source InlineCapabilitySourceParam `json:"source,omitzero" api:"required"`
	// The type of the object. Always `inline`.
	//
	// This field can be elided, and will marshal its zero value as "inline".
	Type constant.Inline `json:"type" default:"inline"`
	paramObj
}

func (r HostedSkillParamInline) MarshalJSON() (data []byte, err error) {
	type shadow HostedSkillParamInline
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *HostedSkillParamInline) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A skill installed from the Skills API.
type HostedSkillReference struct {
	// The installed skill description.
	Description string `json:"description" api:"required"`
	// The installed skill name.
	Name string `json:"name" api:"required"`
	// The referenced skill ID.
	SkillID string `json:"skill_id" api:"required"`
	// The type of the object. Always `skill_reference`.
	Type constant.SkillReference `json:"type" default:"skill_reference"`
	// The concrete skill version installed for this session.
	Version string `json:"version" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Description respjson.Field
		Name        respjson.Field
		SkillID     respjson.Field
		Type        respjson.Field
		Version     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r HostedSkillReference) RawJSON() string { return r.JSON.raw }
func (r *HostedSkillReference) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Provides ZIP bytes encoded with standard base64.
//
// The properties Data, MediaType, Type are required.
type InlineCapabilitySourceParam struct {
	// Standard-base64 encoded ZIP archive bytes.
	Data string `json:"data" api:"required"`
	// The archive media type, always `application/zip`.
	//
	// This field can be elided, and will marshal its zero value as "application/zip".
	MediaType constant.ApplicationZip `json:"media_type" default:"application/zip"`
	// The type of the object. Always `base64`.
	//
	// This field can be elided, and will marshal its zero value as "base64".
	Type constant.Base64 `json:"type" default:"base64"`
	paramObj
}

func (r InlineCapabilitySourceParam) MarshalJSON() (data []byte, err error) {
	type shadow InlineCapabilitySourceParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *InlineCapabilitySourceParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// InputContentUnion contains all possible properties and values from
// [InputContentInputText], [InputContentInputImage].
//
// Use the [InputContentUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type InputContentUnion struct {
	// This field is from variant [InputContentInputText].
	Text string `json:"text"`
	// Any of "input_text", "input_image".
	Type string `json:"type"`
	// This field is from variant [InputContentInputImage].
	ImageURL string `json:"image_url"`
	JSON     struct {
		Text     respjson.Field
		Type     respjson.Field
		ImageURL respjson.Field
		raw      string
	} `json:"-"`
}

// anyInputContent is implemented by each variant of [InputContentUnion] to add
// type safety for the return type of [InputContentUnion.AsAny]
type anyInputContent interface {
	implInputContentUnion()
}

func (InputContentInputText) implInputContentUnion()  {}
func (InputContentInputImage) implInputContentUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := InputContentUnion.AsAny().(type) {
//	case openai.InputContentInputText:
//	case openai.InputContentInputImage:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u InputContentUnion) AsAny() anyInputContent {
	switch u.Type {
	case "input_text":
		return u.AsInputText()
	case "input_image":
		return u.AsInputImage()
	}
	return nil
}

func (u InputContentUnion) AsInputText() (v InputContentInputText) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u InputContentUnion) AsInputImage() (v InputContentInputImage) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u InputContentUnion) RawJSON() string { return u.JSON.raw }

func (r *InputContentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Text input recorded in a session item.
type InputContentInputText struct {
	// The text supplied to the agent.
	Text string `json:"text" api:"required"`
	// The type of the object. Always `input_text`.
	Type constant.InputText `json:"type" default:"input_text"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Text        respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r InputContentInputText) RawJSON() string { return r.JSON.raw }
func (r *InputContentInputText) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Image input recorded in a session item.
type InputContentInputImage struct {
	// The URL of the image supplied to the agent, which may be a base64-encoded data
	// URL.
	ImageURL string `json:"image_url" api:"required"`
	// The type of the object. Always `input_image`.
	Type constant.InputImage `json:"type" default:"input_image"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ImageURL    respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r InputContentInputImage) RawJSON() string { return r.JSON.raw }
func (r *InputContentInputImage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func InputContentParamOfParamInputText(text string) InputContentParamUnion {
	var paramInputText InputContentParamInputText
	paramInputText.Text = text
	return InputContentParamUnion{OfParamInputText: &paramInputText}
}

func InputContentParamOfParamInputImage(imageURL string) InputContentParamUnion {
	var paramInputImage InputContentParamInputImage
	paramInputImage.ImageURL = imageURL
	return InputContentParamUnion{OfParamInputImage: &paramInputImage}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type InputContentParamUnion struct {
	OfParamInputText  *InputContentParamInputText  `json:",omitzero,inline"`
	OfParamInputImage *InputContentParamInputImage `json:",omitzero,inline"`
	paramUnion
}

func (u InputContentParamUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfParamInputText, u.OfParamInputImage)
}
func (u *InputContentParamUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u InputContentParamUnion) GetText() *string {
	if vt := u.OfParamInputText; vt != nil {
		return &vt.Text
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u InputContentParamUnion) GetImageURL() *string {
	if vt := u.OfParamInputImage; vt != nil {
		return &vt.ImageURL
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u InputContentParamUnion) GetType() *string {
	if vt := u.OfParamInputText; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfParamInputImage; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[InputContentParamUnion](
		"type",
		apijson.Discriminator[InputContentParamInputText]("input_text"),
		apijson.Discriminator[InputContentParamInputImage]("input_image"),
	)
}

// Text input to the model.
//
// The properties Text, Type are required.
type InputContentParamInputText struct {
	// The text sent to the model.
	Text string `json:"text" api:"required"`
	// The type of the object. Always `input_text`.
	//
	// This field can be elided, and will marshal its zero value as "input_text".
	Type constant.InputText `json:"type" default:"input_text"`
	paramObj
}

func (r InputContentParamInputText) MarshalJSON() (data []byte, err error) {
	type shadow InputContentParamInputText
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *InputContentParamInputText) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Image input to the model.
//
// The properties ImageURL, Type are required.
type InputContentParamInputImage struct {
	// The URL of the image sent to the model.
	ImageURL string `json:"image_url" api:"required"`
	// The type of the object. Always `input_image`.
	//
	// This field can be elided, and will marshal its zero value as "input_image".
	Type constant.InputImage `json:"type" default:"input_image"`
	paramObj
}

func (r InputContentParamInputImage) MarshalJSON() (data []byte, err error) {
	type shadow InputContentParamInputImage
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *InputContentParamInputImage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// McpTransportUnion contains all possible properties and values from
// [McpTransportHTTP], [McpTransportStdio].
//
// Use the [McpTransportUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type McpTransportUnion struct {
	// This field is from variant [McpTransportHTTP].
	ServerURL string `json:"server_url"`
	// Any of "http", "stdio".
	Type string `json:"type"`
	// This field is from variant [McpTransportStdio].
	Args []string `json:"args"`
	// This field is from variant [McpTransportStdio].
	Command string `json:"command"`
	// This field is from variant [McpTransportStdio].
	Cwd string `json:"cwd"`
	// This field is from variant [McpTransportStdio].
	EnvVars []string `json:"env_vars"`
	JSON    struct {
		ServerURL respjson.Field
		Type      respjson.Field
		Args      respjson.Field
		Command   respjson.Field
		Cwd       respjson.Field
		EnvVars   respjson.Field
		raw       string
	} `json:"-"`
}

// anyMcpTransport is implemented by each variant of [McpTransportUnion] to add
// type safety for the return type of [McpTransportUnion.AsAny]
type anyMcpTransport interface {
	implMcpTransportUnion()
}

func (McpTransportHTTP) implMcpTransportUnion()  {}
func (McpTransportStdio) implMcpTransportUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := McpTransportUnion.AsAny().(type) {
//	case openai.McpTransportHTTP:
//	case openai.McpTransportStdio:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u McpTransportUnion) AsAny() anyMcpTransport {
	switch u.Type {
	case "http":
		return u.AsHTTP()
	case "stdio":
		return u.AsStdio()
	}
	return nil
}

func (u McpTransportUnion) AsHTTP() (v McpTransportHTTP) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u McpTransportUnion) AsStdio() (v McpTransportStdio) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u McpTransportUnion) RawJSON() string { return u.JSON.raw }

func (r *McpTransportUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Connects to an MCP server over HTTP.
type McpTransportHTTP struct {
	// The URL of the MCP server.
	ServerURL string `json:"server_url" api:"required"`
	// The type of the object. Always `http`.
	Type constant.HTTP `json:"type" default:"http"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ServerURL   respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r McpTransportHTTP) RawJSON() string { return r.JSON.raw }
func (r *McpTransportHTTP) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Starts an MCP server as a local process.
type McpTransportStdio struct {
	// Arguments passed to the MCP server command.
	Args []string `json:"args" api:"required"`
	// The command used to start the MCP server.
	Command string `json:"command" api:"required"`
	// The working directory used to start the MCP server.
	Cwd string `json:"cwd" api:"required"`
	// Environment variable names inherited from the execution environment.
	EnvVars []string `json:"env_vars" api:"required"`
	// The type of the object. Always `stdio`.
	Type constant.Stdio `json:"type" default:"stdio"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Args        respjson.Field
		Command     respjson.Field
		Cwd         respjson.Field
		EnvVars     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r McpTransportStdio) RawJSON() string { return r.JSON.raw }
func (r *McpTransportStdio) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func McpTransportParamOfParamHTTP(serverURL string) McpTransportParamUnion {
	var paramHTTP McpTransportParamHTTP
	paramHTTP.ServerURL = serverURL
	return McpTransportParamUnion{OfParamHTTP: &paramHTTP}
}

func McpTransportParamOfParamStdio(command string, cwd string) McpTransportParamUnion {
	var paramStdio McpTransportParamStdio
	paramStdio.Command = command
	paramStdio.Cwd = cwd
	return McpTransportParamUnion{OfParamStdio: &paramStdio}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type McpTransportParamUnion struct {
	OfParamHTTP  *McpTransportParamHTTP  `json:",omitzero,inline"`
	OfParamStdio *McpTransportParamStdio `json:",omitzero,inline"`
	paramUnion
}

func (u McpTransportParamUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfParamHTTP, u.OfParamStdio)
}
func (u *McpTransportParamUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u McpTransportParamUnion) GetServerURL() *string {
	if vt := u.OfParamHTTP; vt != nil {
		return &vt.ServerURL
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u McpTransportParamUnion) GetAuthorization() *string {
	if vt := u.OfParamHTTP; vt != nil && vt.Authorization.Valid() {
		return &vt.Authorization.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u McpTransportParamUnion) GetHeaders() map[string]string {
	if vt := u.OfParamHTTP; vt != nil {
		return vt.Headers
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u McpTransportParamUnion) GetCommand() *string {
	if vt := u.OfParamStdio; vt != nil {
		return &vt.Command
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u McpTransportParamUnion) GetCwd() *string {
	if vt := u.OfParamStdio; vt != nil {
		return &vt.Cwd
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u McpTransportParamUnion) GetArgs() []string {
	if vt := u.OfParamStdio; vt != nil {
		return vt.Args
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u McpTransportParamUnion) GetEnv() map[string]string {
	if vt := u.OfParamStdio; vt != nil {
		return vt.Env
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u McpTransportParamUnion) GetEnvVars() []string {
	if vt := u.OfParamStdio; vt != nil {
		return vt.EnvVars
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u McpTransportParamUnion) GetType() *string {
	if vt := u.OfParamHTTP; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfParamStdio; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[McpTransportParamUnion](
		"type",
		apijson.Discriminator[McpTransportParamHTTP]("http"),
		apijson.Discriminator[McpTransportParamStdio]("stdio"),
	)
}

// Connects to an MCP server over HTTP.
//
// The properties ServerURL, Type are required.
type McpTransportParamHTTP struct {
	// The URL of the MCP server.
	ServerURL string `json:"server_url" api:"required"`
	// The authorization value sent to the MCP server, if any.
	Authorization param.Opt[string] `json:"authorization,omitzero"`
	// Additional HTTP headers sent to the MCP server.
	Headers map[string]string `json:"headers,omitzero"`
	// The type of the object. Always `http`.
	//
	// This field can be elided, and will marshal its zero value as "http".
	Type constant.HTTP `json:"type" default:"http"`
	paramObj
}

func (r McpTransportParamHTTP) MarshalJSON() (data []byte, err error) {
	type shadow McpTransportParamHTTP
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *McpTransportParamHTTP) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Starts an MCP server as a local process.
//
// The properties Command, Cwd, Type are required.
type McpTransportParamStdio struct {
	// The command used to start the MCP server.
	Command string `json:"command" api:"required"`
	// The working directory used to start the MCP server.
	Cwd string `json:"cwd" api:"required"`
	// Arguments passed to the MCP server command.
	Args []string `json:"args,omitzero"`
	// Environment variables set for the MCP server process.
	Env map[string]string `json:"env,omitzero"`
	// Environment variable names to inherit from the selected execution environment.
	EnvVars []string `json:"env_vars,omitzero"`
	// The type of the object. Always `stdio`.
	//
	// This field can be elided, and will marshal its zero value as "stdio".
	Type constant.Stdio `json:"type" default:"stdio"`
	paramObj
}

func (r McpTransportParamStdio) MarshalJSON() (data []byte, err error) {
	type shadow McpTransportParamStdio
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *McpTransportParamStdio) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The resolved configuration for creating and coordinating subagents.
type MultiAgentConfig struct {
	// Whether subagent tools are enabled. Defaults to false.
	Enabled bool `json:"enabled" api:"required"`
	// Maximum number of subagents that may run concurrently, or null when disabled.
	// Defaults to 6 when enabled.
	MaxConcurrentSubagents int64 `json:"max_concurrent_subagents" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Enabled                respjson.Field
		MaxConcurrentSubagents respjson.Field
		ExtraFields            map[string]respjson.Field
		raw                    string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r MultiAgentConfig) RawJSON() string { return r.JSON.raw }
func (r *MultiAgentConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Explicit configuration for creating and coordinating subagents.
//
// The property Enabled is required.
type MultiAgentConfigParam struct {
	// Whether subagent tools are enabled.
	Enabled bool `json:"enabled" api:"required"`
	// Maximum number of subagents that may run concurrently. Defaults to 6.
	MaxConcurrentSubagents param.Opt[int64] `json:"max_concurrent_subagents,omitzero"`
	paramObj
}

func (r MultiAgentConfigParam) MarshalJSON() (data []byte, err error) {
	type shadow MultiAgentConfigParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *MultiAgentConfigParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A text content part produced by the agent.
type OutputText struct {
	// The text produced by the agent.
	Text string `json:"text" api:"required"`
	// The content type. Always `output_text`.
	Type constant.OutputText `json:"type" default:"output_text"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Text        respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r OutputText) RawJSON() string { return r.JSON.raw }
func (r *OutputText) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// PersistedAgentToolUnion contains all possible properties and values from
// [PersistedAgentToolFunction], [PersistedAgentToolToolSearch],
// [PersistedAgentToolProgrammaticToolCalling], [PersistedAgentToolMcp],
// [PersistedAgentToolWebSearch].
//
// Use the [PersistedAgentToolUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type PersistedAgentToolUnion struct {
	// This field is from variant [PersistedAgentToolFunction].
	DeferLoading bool `json:"defer_loading"`
	// This field is from variant [PersistedAgentToolFunction].
	Description string `json:"description"`
	// This field is from variant [PersistedAgentToolFunction].
	Name string `json:"name"`
	// This field is from variant [PersistedAgentToolFunction].
	Parameters map[string]any `json:"parameters"`
	// Any of "function", "tool_search", "programmatic_tool_calling", "mcp",
	// "web_search".
	Type string `json:"type"`
	// This field is from variant [PersistedAgentToolProgrammaticToolCalling].
	Enabled bool `json:"enabled"`
	// This field is from variant [PersistedAgentToolMcp].
	AllowedTools []string `json:"allowed_tools"`
	// This field is from variant [PersistedAgentToolMcp].
	ConnectionOrigin string `json:"connection_origin"`
	// This field is from variant [PersistedAgentToolMcp].
	CredentialID string `json:"credential_id"`
	// This field is from variant [PersistedAgentToolMcp].
	RequestMetadata map[string]any `json:"request_metadata"`
	// This field is from variant [PersistedAgentToolMcp].
	Required bool `json:"required"`
	// This field is from variant [PersistedAgentToolMcp].
	ServerLabel string `json:"server_label"`
	// This field is from variant [PersistedAgentToolMcp].
	Transport PersistedMcpTransportUnion `json:"transport"`
	// This field is from variant [PersistedAgentToolWebSearch].
	AllowedDomains []string `json:"allowed_domains"`
	// This field is from variant [PersistedAgentToolWebSearch].
	ContextSize string `json:"context_size"`
	// This field is from variant [PersistedAgentToolWebSearch].
	Location PersistedAgentToolWebSearchLocation `json:"location"`
	// This field is from variant [PersistedAgentToolWebSearch].
	Mode string `json:"mode"`
	JSON struct {
		DeferLoading     respjson.Field
		Description      respjson.Field
		Name             respjson.Field
		Parameters       respjson.Field
		Type             respjson.Field
		Enabled          respjson.Field
		AllowedTools     respjson.Field
		ConnectionOrigin respjson.Field
		CredentialID     respjson.Field
		RequestMetadata  respjson.Field
		Required         respjson.Field
		ServerLabel      respjson.Field
		Transport        respjson.Field
		AllowedDomains   respjson.Field
		ContextSize      respjson.Field
		Location         respjson.Field
		Mode             respjson.Field
		raw              string
	} `json:"-"`
}

// anyPersistedAgentTool is implemented by each variant of
// [PersistedAgentToolUnion] to add type safety for the return type of
// [PersistedAgentToolUnion.AsAny]
type anyPersistedAgentTool interface {
	implPersistedAgentToolUnion()
}

func (PersistedAgentToolFunction) implPersistedAgentToolUnion()                {}
func (PersistedAgentToolToolSearch) implPersistedAgentToolUnion()              {}
func (PersistedAgentToolProgrammaticToolCalling) implPersistedAgentToolUnion() {}
func (PersistedAgentToolMcp) implPersistedAgentToolUnion()                     {}
func (PersistedAgentToolWebSearch) implPersistedAgentToolUnion()               {}

// Use the following switch statement to find the correct variant
//
//	switch variant := PersistedAgentToolUnion.AsAny().(type) {
//	case openai.PersistedAgentToolFunction:
//	case openai.PersistedAgentToolToolSearch:
//	case openai.PersistedAgentToolProgrammaticToolCalling:
//	case openai.PersistedAgentToolMcp:
//	case openai.PersistedAgentToolWebSearch:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u PersistedAgentToolUnion) AsAny() anyPersistedAgentTool {
	switch u.Type {
	case "function":
		return u.AsFunction()
	case "tool_search":
		return u.AsToolSearch()
	case "programmatic_tool_calling":
		return u.AsProgrammaticToolCalling()
	case "mcp":
		return u.AsMcp()
	case "web_search":
		return u.AsWebSearch()
	}
	return nil
}

func (u PersistedAgentToolUnion) AsFunction() (v PersistedAgentToolFunction) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u PersistedAgentToolUnion) AsToolSearch() (v PersistedAgentToolToolSearch) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u PersistedAgentToolUnion) AsProgrammaticToolCalling() (v PersistedAgentToolProgrammaticToolCalling) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u PersistedAgentToolUnion) AsMcp() (v PersistedAgentToolMcp) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u PersistedAgentToolUnion) AsWebSearch() (v PersistedAgentToolWebSearch) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u PersistedAgentToolUnion) RawJSON() string { return u.JSON.raw }

func (r *PersistedAgentToolUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A function defined by the application.
type PersistedAgentToolFunction struct {
	// Whether the function is deferred and discovered through tool search.
	DeferLoading bool `json:"defer_loading" api:"required"`
	// A description of what the function does.
	Description string `json:"description" api:"required"`
	// The name of the function.
	Name string `json:"name" api:"required"`
	// A JSON Schema object describing the function's arguments.
	Parameters map[string]any `json:"parameters" api:"required"`
	// The type of the object. Always `function`.
	Type constant.Function `json:"type" default:"function"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		DeferLoading respjson.Field
		Description  respjson.Field
		Name         respjson.Field
		Parameters   respjson.Field
		Type         respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r PersistedAgentToolFunction) RawJSON() string { return r.JSON.raw }
func (r *PersistedAgentToolFunction) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Discovers deferred function tools and loads them into the model context.
type PersistedAgentToolToolSearch struct {
	// The type of the object. Always `tool_search`.
	Type constant.ToolSearch `json:"type" default:"tool_search"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r PersistedAgentToolToolSearch) RawJSON() string { return r.JSON.raw }
func (r *PersistedAgentToolToolSearch) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Enables calling tools from model-generated code.
type PersistedAgentToolProgrammaticToolCalling struct {
	// Whether tools can be called from model-generated code.
	Enabled bool `json:"enabled" api:"required"`
	// The type of the object. Always `programmatic_tool_calling`.
	Type constant.ProgrammaticToolCalling `json:"type" default:"programmatic_tool_calling"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Enabled     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r PersistedAgentToolProgrammaticToolCalling) RawJSON() string { return r.JSON.raw }
func (r *PersistedAgentToolProgrammaticToolCalling) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Tools provided by a remote MCP server without stored credentials.
type PersistedAgentToolMcp struct {
	// The MCP tools the agent may call, or null when all server tools are allowed.
	AllowedTools []string `json:"allowed_tools" api:"required"`
	// Where outbound MCP HTTP connections originate.
	//
	// Any of "service", "environment".
	ConnectionOrigin string `json:"connection_origin" api:"required"`
	// The vault credential selected for this MCP server, if any.
	CredentialID string `json:"credential_id" api:"required"`
	// Metadata included with requests to this MCP server.
	RequestMetadata map[string]any `json:"request_metadata" api:"required"`
	// Whether this MCP server must initialize before the first turn.
	Required bool `json:"required" api:"required"`
	// A label used to identify the MCP server in tool calls.
	ServerLabel string `json:"server_label" api:"required"`
	// The credential-free transport used to connect to the MCP server.
	Transport PersistedMcpTransportUnion `json:"transport" api:"required"`
	// The type of the object. Always `mcp`.
	Type constant.Mcp `json:"type" default:"mcp"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		AllowedTools     respjson.Field
		ConnectionOrigin respjson.Field
		CredentialID     respjson.Field
		RequestMetadata  respjson.Field
		Required         respjson.Field
		ServerLabel      respjson.Field
		Transport        respjson.Field
		Type             respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r PersistedAgentToolMcp) RawJSON() string { return r.JSON.raw }
func (r *PersistedAgentToolMcp) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Web search.
type PersistedAgentToolWebSearch struct {
	// Allowed search domains, or `null` when the search is unrestricted.
	AllowedDomains []string `json:"allowed_domains" api:"required"`
	// The amount of search context made available to the model. Defaults to `medium`.
	//
	// Any of "low", "medium", "high".
	ContextSize string `json:"context_size" api:"required"`
	// Approximate user location used to localize web search results.
	Location PersistedAgentToolWebSearchLocation `json:"location" api:"required"`
	// The source used for web search results.
	//
	// Any of "disabled", "cached", "live".
	Mode string `json:"mode" api:"required"`
	// The type of the object. Always `web_search`.
	Type constant.WebSearch `json:"type" default:"web_search"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		AllowedDomains respjson.Field
		ContextSize    respjson.Field
		Location       respjson.Field
		Mode           respjson.Field
		Type           respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r PersistedAgentToolWebSearch) RawJSON() string { return r.JSON.raw }
func (r *PersistedAgentToolWebSearch) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Approximate user location used to localize web search results.
type PersistedAgentToolWebSearchLocation struct {
	// The city name.
	City string `json:"city" api:"required"`
	// The two-letter ISO country code, such as `US`.
	Country string `json:"country" api:"required"`
	// The region or state name.
	Region string `json:"region" api:"required"`
	// The IANA timezone, such as `America/Los_Angeles`.
	Timezone string `json:"timezone" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		City        respjson.Field
		Country     respjson.Field
		Region      respjson.Field
		Timezone    respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r PersistedAgentToolWebSearchLocation) RawJSON() string { return r.JSON.raw }
func (r *PersistedAgentToolWebSearchLocation) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func PersistedAgentToolParamOfParamFunction(description string, name string, parameters map[string]any) PersistedAgentToolParamUnion {
	var paramFunction PersistedAgentToolParamFunction
	paramFunction.Description = description
	paramFunction.Name = name
	paramFunction.Parameters = parameters
	return PersistedAgentToolParamUnion{OfParamFunction: &paramFunction}
}

func PersistedAgentToolParamOfParamMcp[
	T PersistedMcpTransportParamHTTP | PersistedMcpTransportParamStdio,
](serverLabel string, transport T) PersistedAgentToolParamUnion {
	var paramMcp PersistedAgentToolParamMcp
	paramMcp.ServerLabel = serverLabel
	switch v := any(transport).(type) {
	case PersistedMcpTransportParamHTTP:
		paramMcp.Transport.OfParamHTTP = &v
	case PersistedMcpTransportParamStdio:
		paramMcp.Transport.OfParamStdio = &v
	}
	return PersistedAgentToolParamUnion{OfParamMcp: &paramMcp}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type PersistedAgentToolParamUnion struct {
	OfParamFunction                *PersistedAgentToolParamFunction                `json:",omitzero,inline"`
	OfParamToolSearch              *PersistedAgentToolParamToolSearch              `json:",omitzero,inline"`
	OfParamProgrammaticToolCalling *PersistedAgentToolParamProgrammaticToolCalling `json:",omitzero,inline"`
	OfParamMcp                     *PersistedAgentToolParamMcp                     `json:",omitzero,inline"`
	OfParamWebSearch               *PersistedAgentToolParamWebSearch               `json:",omitzero,inline"`
	paramUnion
}

func (u PersistedAgentToolParamUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfParamFunction,
		u.OfParamToolSearch,
		u.OfParamProgrammaticToolCalling,
		u.OfParamMcp,
		u.OfParamWebSearch)
}
func (u *PersistedAgentToolParamUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u PersistedAgentToolParamUnion) GetDescription() *string {
	if vt := u.OfParamFunction; vt != nil {
		return &vt.Description
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u PersistedAgentToolParamUnion) GetName() *string {
	if vt := u.OfParamFunction; vt != nil {
		return &vt.Name
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u PersistedAgentToolParamUnion) GetParameters() map[string]any {
	if vt := u.OfParamFunction; vt != nil {
		return vt.Parameters
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u PersistedAgentToolParamUnion) GetDeferLoading() *bool {
	if vt := u.OfParamFunction; vt != nil && vt.DeferLoading.Valid() {
		return &vt.DeferLoading.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u PersistedAgentToolParamUnion) GetEnabled() *bool {
	if vt := u.OfParamProgrammaticToolCalling; vt != nil && vt.Enabled.Valid() {
		return &vt.Enabled.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u PersistedAgentToolParamUnion) GetServerLabel() *string {
	if vt := u.OfParamMcp; vt != nil {
		return &vt.ServerLabel
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u PersistedAgentToolParamUnion) GetTransport() *PersistedMcpTransportParamUnion {
	if vt := u.OfParamMcp; vt != nil {
		return &vt.Transport
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u PersistedAgentToolParamUnion) GetAllowedTools() []string {
	if vt := u.OfParamMcp; vt != nil {
		return vt.AllowedTools
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u PersistedAgentToolParamUnion) GetConnectionOrigin() *string {
	if vt := u.OfParamMcp; vt != nil {
		return &vt.ConnectionOrigin
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u PersistedAgentToolParamUnion) GetCredentialID() *string {
	if vt := u.OfParamMcp; vt != nil && vt.CredentialID.Valid() {
		return &vt.CredentialID.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u PersistedAgentToolParamUnion) GetRequestMetadata() map[string]any {
	if vt := u.OfParamMcp; vt != nil {
		return vt.RequestMetadata
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u PersistedAgentToolParamUnion) GetRequired() *bool {
	if vt := u.OfParamMcp; vt != nil && vt.Required.Valid() {
		return &vt.Required.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u PersistedAgentToolParamUnion) GetAllowedDomains() []string {
	if vt := u.OfParamWebSearch; vt != nil {
		return vt.AllowedDomains
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u PersistedAgentToolParamUnion) GetContextSize() *string {
	if vt := u.OfParamWebSearch; vt != nil {
		return &vt.ContextSize
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u PersistedAgentToolParamUnion) GetLocation() *PersistedAgentToolParamWebSearchLocation {
	if vt := u.OfParamWebSearch; vt != nil {
		return &vt.Location
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u PersistedAgentToolParamUnion) GetMode() *string {
	if vt := u.OfParamWebSearch; vt != nil {
		return &vt.Mode
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u PersistedAgentToolParamUnion) GetType() *string {
	if vt := u.OfParamFunction; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfParamToolSearch; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfParamProgrammaticToolCalling; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfParamMcp; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfParamWebSearch; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[PersistedAgentToolParamUnion](
		"type",
		apijson.Discriminator[PersistedAgentToolParamFunction]("function"),
		apijson.Discriminator[PersistedAgentToolParamToolSearch]("tool_search"),
		apijson.Discriminator[PersistedAgentToolParamProgrammaticToolCalling]("programmatic_tool_calling"),
		apijson.Discriminator[PersistedAgentToolParamMcp]("mcp"),
		apijson.Discriminator[PersistedAgentToolParamWebSearch]("web_search"),
	)
}

// A function defined by the application.
//
// The properties Description, Name, Parameters, Type are required.
type PersistedAgentToolParamFunction struct {
	// A description of what the function does.
	Description string `json:"description" api:"required"`
	// The name of the function.
	Name string `json:"name" api:"required"`
	// A JSON Schema object describing the function's arguments.
	Parameters map[string]any `json:"parameters,omitzero" api:"required"`
	// Whether this function is deferred and discovered through tool search. Defaults
	// to `false`.
	DeferLoading param.Opt[bool] `json:"defer_loading,omitzero"`
	// The type of the object. Always `function`.
	//
	// This field can be elided, and will marshal its zero value as "function".
	Type constant.Function `json:"type" default:"function"`
	paramObj
}

func (r PersistedAgentToolParamFunction) MarshalJSON() (data []byte, err error) {
	type shadow PersistedAgentToolParamFunction
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *PersistedAgentToolParamFunction) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewPersistedAgentToolParamToolSearch() PersistedAgentToolParamToolSearch {
	return PersistedAgentToolParamToolSearch{
		Type: "tool_search",
	}
}

// Discovers deferred function tools and loads them into the model context.
//
// This struct has a constant value, construct it with
// [NewPersistedAgentToolParamToolSearch].
type PersistedAgentToolParamToolSearch struct {
	// The type of the object. Always `tool_search`.
	Type constant.ToolSearch `json:"type" default:"tool_search"`
	paramObj
}

func (r PersistedAgentToolParamToolSearch) MarshalJSON() (data []byte, err error) {
	type shadow PersistedAgentToolParamToolSearch
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *PersistedAgentToolParamToolSearch) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Enables calling tools from model-generated code.
//
// The property Type is required.
type PersistedAgentToolParamProgrammaticToolCalling struct {
	// Whether tools can be called from model-generated code. Defaults to `true`.
	Enabled param.Opt[bool] `json:"enabled,omitzero"`
	// The type of the object. Always `programmatic_tool_calling`.
	//
	// This field can be elided, and will marshal its zero value as
	// "programmatic_tool_calling".
	Type constant.ProgrammaticToolCalling `json:"type" default:"programmatic_tool_calling"`
	paramObj
}

func (r PersistedAgentToolParamProgrammaticToolCalling) MarshalJSON() (data []byte, err error) {
	type shadow PersistedAgentToolParamProgrammaticToolCalling
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *PersistedAgentToolParamProgrammaticToolCalling) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Tools provided by a remote MCP server without stored credentials.
//
// The properties ServerLabel, Transport, Type are required.
type PersistedAgentToolParamMcp struct {
	// A label used to identify the MCP server in tool calls.
	ServerLabel string `json:"server_label" api:"required"`
	// The credential-free transport used to connect to the MCP server.
	Transport PersistedMcpTransportParamUnion `json:"transport,omitzero" api:"required"`
	// The vault credential selected for this MCP server. Optional when exactly one
	// attached credential matches the server URL.
	CredentialID param.Opt[string] `json:"credential_id,omitzero"`
	// Whether this MCP server must initialize before the first turn. Defaults to
	// `false`.
	Required param.Opt[bool] `json:"required,omitzero"`
	// The MCP tools the agent may call. All server tools are allowed when omitted.
	AllowedTools []string `json:"allowed_tools,omitzero"`
	// Where outbound MCP HTTP connections originate.
	//
	// Any of "service", "environment".
	ConnectionOrigin string `json:"connection_origin,omitzero"`
	// Metadata included with requests to this MCP server.
	RequestMetadata map[string]any `json:"request_metadata,omitzero"`
	// The type of the object. Always `mcp`.
	//
	// This field can be elided, and will marshal its zero value as "mcp".
	Type constant.Mcp `json:"type" default:"mcp"`
	paramObj
}

func (r PersistedAgentToolParamMcp) MarshalJSON() (data []byte, err error) {
	type shadow PersistedAgentToolParamMcp
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *PersistedAgentToolParamMcp) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[PersistedAgentToolParamMcp](
		"connection_origin", "service", "environment",
	)
}

// Web search.
//
// The property Type is required.
type PersistedAgentToolParamWebSearch struct {
	// Domains the search may include.
	AllowedDomains []string `json:"allowed_domains,omitzero"`
	// The amount of web search context made available to the model.
	//
	// Any of "low", "medium", "high".
	ContextSize string `json:"context_size,omitzero"`
	// Approximate user location used to localize web search results.
	Location PersistedAgentToolParamWebSearchLocation `json:"location,omitzero"`
	// The source used for web search results.
	//
	// Any of "disabled", "cached", "live".
	Mode string `json:"mode,omitzero"`
	// The type of the object. Always `web_search`.
	//
	// This field can be elided, and will marshal its zero value as "web_search".
	Type constant.WebSearch `json:"type" default:"web_search"`
	paramObj
}

func (r PersistedAgentToolParamWebSearch) MarshalJSON() (data []byte, err error) {
	type shadow PersistedAgentToolParamWebSearch
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *PersistedAgentToolParamWebSearch) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[PersistedAgentToolParamWebSearch](
		"context_size", "low", "medium", "high",
	)
	apijson.RegisterFieldValidator[PersistedAgentToolParamWebSearch](
		"mode", "disabled", "cached", "live",
	)
}

// Approximate user location used to localize web search results.
type PersistedAgentToolParamWebSearchLocation struct {
	// The city name.
	City param.Opt[string] `json:"city,omitzero"`
	// The two-letter ISO country code, such as `US`.
	Country param.Opt[string] `json:"country,omitzero"`
	// The region or state name.
	Region param.Opt[string] `json:"region,omitzero"`
	// The IANA timezone, such as `America/Los_Angeles`.
	Timezone param.Opt[string] `json:"timezone,omitzero"`
	paramObj
}

func (r PersistedAgentToolParamWebSearchLocation) MarshalJSON() (data []byte, err error) {
	type shadow PersistedAgentToolParamWebSearchLocation
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *PersistedAgentToolParamWebSearchLocation) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// PersistedMcpTransportUnion contains all possible properties and values from
// [PersistedMcpTransportHTTP], [PersistedMcpTransportStdio].
//
// Use the [PersistedMcpTransportUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type PersistedMcpTransportUnion struct {
	// This field is from variant [PersistedMcpTransportHTTP].
	Headers map[string]string `json:"headers"`
	// This field is from variant [PersistedMcpTransportHTTP].
	ServerURL string `json:"server_url"`
	// Any of "http", "stdio".
	Type string `json:"type"`
	// This field is from variant [PersistedMcpTransportStdio].
	Args []string `json:"args"`
	// This field is from variant [PersistedMcpTransportStdio].
	Command string `json:"command"`
	// This field is from variant [PersistedMcpTransportStdio].
	Cwd string `json:"cwd"`
	// This field is from variant [PersistedMcpTransportStdio].
	EnvVars []string `json:"env_vars"`
	JSON    struct {
		Headers   respjson.Field
		ServerURL respjson.Field
		Type      respjson.Field
		Args      respjson.Field
		Command   respjson.Field
		Cwd       respjson.Field
		EnvVars   respjson.Field
		raw       string
	} `json:"-"`
}

// anyPersistedMcpTransport is implemented by each variant of
// [PersistedMcpTransportUnion] to add type safety for the return type of
// [PersistedMcpTransportUnion.AsAny]
type anyPersistedMcpTransport interface {
	implPersistedMcpTransportUnion()
}

func (PersistedMcpTransportHTTP) implPersistedMcpTransportUnion()  {}
func (PersistedMcpTransportStdio) implPersistedMcpTransportUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := PersistedMcpTransportUnion.AsAny().(type) {
//	case openai.PersistedMcpTransportHTTP:
//	case openai.PersistedMcpTransportStdio:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u PersistedMcpTransportUnion) AsAny() anyPersistedMcpTransport {
	switch u.Type {
	case "http":
		return u.AsHTTP()
	case "stdio":
		return u.AsStdio()
	}
	return nil
}

func (u PersistedMcpTransportUnion) AsHTTP() (v PersistedMcpTransportHTTP) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u PersistedMcpTransportUnion) AsStdio() (v PersistedMcpTransportStdio) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u PersistedMcpTransportUnion) RawJSON() string { return u.JSON.raw }

func (r *PersistedMcpTransportUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Connects to an MCP server over HTTP.
type PersistedMcpTransportHTTP struct {
	// Non-secret HTTP headers sent to the MCP server.
	Headers map[string]string `json:"headers" api:"required"`
	// The URL of the MCP server.
	ServerURL string `json:"server_url" api:"required"`
	// The type of the object. Always `http`.
	Type constant.HTTP `json:"type" default:"http"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Headers     respjson.Field
		ServerURL   respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r PersistedMcpTransportHTTP) RawJSON() string { return r.JSON.raw }
func (r *PersistedMcpTransportHTTP) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Starts an MCP server as a local process.
type PersistedMcpTransportStdio struct {
	// Arguments passed to the MCP server command.
	Args []string `json:"args" api:"required"`
	// The command used to start the MCP server.
	Command string `json:"command" api:"required"`
	// The working directory used to start the MCP server.
	Cwd string `json:"cwd" api:"required"`
	// Environment variable names inherited from the execution environment.
	EnvVars []string `json:"env_vars" api:"required"`
	// The type of the object. Always `stdio`.
	Type constant.Stdio `json:"type" default:"stdio"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Args        respjson.Field
		Command     respjson.Field
		Cwd         respjson.Field
		EnvVars     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r PersistedMcpTransportStdio) RawJSON() string { return r.JSON.raw }
func (r *PersistedMcpTransportStdio) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func PersistedMcpTransportParamOfParamHTTP(serverURL string) PersistedMcpTransportParamUnion {
	var paramHTTP PersistedMcpTransportParamHTTP
	paramHTTP.ServerURL = serverURL
	return PersistedMcpTransportParamUnion{OfParamHTTP: &paramHTTP}
}

func PersistedMcpTransportParamOfParamStdio(command string, cwd string) PersistedMcpTransportParamUnion {
	var paramStdio PersistedMcpTransportParamStdio
	paramStdio.Command = command
	paramStdio.Cwd = cwd
	return PersistedMcpTransportParamUnion{OfParamStdio: &paramStdio}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type PersistedMcpTransportParamUnion struct {
	OfParamHTTP  *PersistedMcpTransportParamHTTP  `json:",omitzero,inline"`
	OfParamStdio *PersistedMcpTransportParamStdio `json:",omitzero,inline"`
	paramUnion
}

func (u PersistedMcpTransportParamUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfParamHTTP, u.OfParamStdio)
}
func (u *PersistedMcpTransportParamUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u PersistedMcpTransportParamUnion) GetServerURL() *string {
	if vt := u.OfParamHTTP; vt != nil {
		return &vt.ServerURL
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u PersistedMcpTransportParamUnion) GetHeaders() map[string]string {
	if vt := u.OfParamHTTP; vt != nil {
		return vt.Headers
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u PersistedMcpTransportParamUnion) GetCommand() *string {
	if vt := u.OfParamStdio; vt != nil {
		return &vt.Command
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u PersistedMcpTransportParamUnion) GetCwd() *string {
	if vt := u.OfParamStdio; vt != nil {
		return &vt.Cwd
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u PersistedMcpTransportParamUnion) GetArgs() []string {
	if vt := u.OfParamStdio; vt != nil {
		return vt.Args
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u PersistedMcpTransportParamUnion) GetEnvVars() []string {
	if vt := u.OfParamStdio; vt != nil {
		return vt.EnvVars
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u PersistedMcpTransportParamUnion) GetType() *string {
	if vt := u.OfParamHTTP; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfParamStdio; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[PersistedMcpTransportParamUnion](
		"type",
		apijson.Discriminator[PersistedMcpTransportParamHTTP]("http"),
		apijson.Discriminator[PersistedMcpTransportParamStdio]("stdio"),
	)
}

// Connects to an MCP server over HTTP.
//
// The properties ServerURL, Type are required.
type PersistedMcpTransportParamHTTP struct {
	// The URL of the MCP server.
	ServerURL string `json:"server_url" api:"required"`
	// Non-secret HTTP headers sent to the MCP server.
	Headers map[string]string `json:"headers,omitzero"`
	// The type of the object. Always `http`.
	//
	// This field can be elided, and will marshal its zero value as "http".
	Type constant.HTTP `json:"type" default:"http"`
	paramObj
}

func (r PersistedMcpTransportParamHTTP) MarshalJSON() (data []byte, err error) {
	type shadow PersistedMcpTransportParamHTTP
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *PersistedMcpTransportParamHTTP) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Starts an MCP server as a local process.
//
// The properties Command, Cwd, Type are required.
type PersistedMcpTransportParamStdio struct {
	// The command used to start the MCP server.
	Command string `json:"command" api:"required"`
	// The working directory used to start the MCP server.
	Cwd string `json:"cwd" api:"required"`
	// Arguments passed to the MCP server command.
	Args []string `json:"args,omitzero"`
	// Environment variable names to inherit from the selected execution environment.
	EnvVars []string `json:"env_vars,omitzero"`
	// The type of the object. Always `stdio`.
	//
	// This field can be elided, and will marshal its zero value as "stdio".
	Type constant.Stdio `json:"type" default:"stdio"`
	paramObj
}

func (r PersistedMcpTransportParamStdio) MarshalJSON() (data []byte, err error) {
	type shadow PersistedMcpTransportParamStdio
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *PersistedMcpTransportParamStdio) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// An error payload with the same public fields as Responses API streaming errors.
type SessionError struct {
	// The machine-readable error code, if any.
	Code string `json:"code" api:"required"`
	// A customer-safe explanation of the error.
	Message string `json:"message" api:"required"`
	// The request parameter associated with the error, if any.
	Param string `json:"param" api:"required"`
	// The error type.
	Type string `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Code        respjson.Field
		Message     respjson.Field
		Param       respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r SessionError) RawJSON() string { return r.JSON.raw }
func (r *SessionError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A customer-safe error describing why a session request failed.
type SessionTurnError struct {
	// A stable, machine-readable failure category.
	//
	// Any of "context_length_exceeded", "session_budget_exceeded",
	// "usage_limit_exceeded", "rate_limit_exceeded", "server_overloaded",
	// "cyber_policy", "connection_failed", "server_error", "authentication_error",
	// "invalid_request", "resource_not_found", "sandbox_error",
	// "executor_version_incompatible", "active_turn_not_steerable", "request_timeout",
	// "internal_error".
	Code SessionTurnErrorCode `json:"code" api:"required"`
	// A customer-safe explanation of the failure.
	Message string `json:"message" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Code        respjson.Field
		Message     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r SessionTurnError) RawJSON() string { return r.JSON.raw }
func (r *SessionTurnError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A stable, machine-readable failure category.
type SessionTurnErrorCode string

const (
	SessionTurnErrorCodeContextLengthExceeded       SessionTurnErrorCode = "context_length_exceeded"
	SessionTurnErrorCodeSessionBudgetExceeded       SessionTurnErrorCode = "session_budget_exceeded"
	SessionTurnErrorCodeUsageLimitExceeded          SessionTurnErrorCode = "usage_limit_exceeded"
	SessionTurnErrorCodeRateLimitExceeded           SessionTurnErrorCode = "rate_limit_exceeded"
	SessionTurnErrorCodeServerOverloaded            SessionTurnErrorCode = "server_overloaded"
	SessionTurnErrorCodeCyberPolicy                 SessionTurnErrorCode = "cyber_policy"
	SessionTurnErrorCodeConnectionFailed            SessionTurnErrorCode = "connection_failed"
	SessionTurnErrorCodeServerError                 SessionTurnErrorCode = "server_error"
	SessionTurnErrorCodeAuthenticationError         SessionTurnErrorCode = "authentication_error"
	SessionTurnErrorCodeInvalidRequest              SessionTurnErrorCode = "invalid_request"
	SessionTurnErrorCodeResourceNotFound            SessionTurnErrorCode = "resource_not_found"
	SessionTurnErrorCodeSandboxError                SessionTurnErrorCode = "sandbox_error"
	SessionTurnErrorCodeExecutorVersionIncompatible SessionTurnErrorCode = "executor_version_incompatible"
	SessionTurnErrorCodeActiveTurnNotSteerable      SessionTurnErrorCode = "active_turn_not_steerable"
	SessionTurnErrorCodeRequestTimeout              SessionTurnErrorCode = "request_timeout"
	SessionTurnErrorCodeInternalError               SessionTurnErrorCode = "internal_error"
)

// A confidential setup command executed before the hosted agent starts.
//
// The property Command is required.
type SetupCommandParam struct {
	// The shell command to execute.
	Command string `json:"command" api:"required"`
	// The absolute working directory. Defaults to `/workspace`.
	Cwd param.Opt[string] `json:"cwd,omitzero"`
	paramObj
}

func (r SetupCommandParam) MarshalJSON() (data []byte, err error) {
	type shadow SetupCommandParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *SetupCommandParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A subagent created within a session.
type Subagent struct {
	// The ID of the subagent.
	ID string `json:"id" api:"required"`
	// The Unix timestamp, in seconds, when the subagent was closed. Null while active,
	// including after resume.
	ClosedAt int64 `json:"closed_at" api:"required"`
	// Initial task content, or null when unavailable. Text may contain placeholders
	// for images or audio when only a preview is available.
	Instructions []AgentContentUnion `json:"instructions" api:"required"`
	// The runner-assigned nickname, or null when unavailable.
	Name string `json:"name" api:"required"`
	// The object type. Always `agent.session.subagent`.
	//
	// Any of "agent.session.subagent".
	Object SubagentObject `json:"object" api:"required"`
	// The Unix timestamp, in seconds, when the subagent was first opened. Resuming
	// does not change it.
	OpenedAt int64 `json:"opened_at" api:"required"`
	// The ID of the agent that created this subagent.
	ParentAgentID string `json:"parent_agent_id" api:"required"`
	// The ID of the session that owns the subagent.
	SessionID string `json:"session_id" api:"required"`
	// The current status of the subagent.
	//
	// Any of "active", "closed".
	Status SubagentStatus `json:"status" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID            respjson.Field
		ClosedAt      respjson.Field
		Instructions  respjson.Field
		Name          respjson.Field
		Object        respjson.Field
		OpenedAt      respjson.Field
		ParentAgentID respjson.Field
		SessionID     respjson.Field
		Status        respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r Subagent) RawJSON() string { return r.JSON.raw }
func (r *Subagent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The object type. Always `agent.session.subagent`.
type SubagentObject string

const (
	SubagentObjectAgentSessionSubagent SubagentObject = "agent.session.subagent"
)

// The current status of the subagent.
type SubagentStatus string

const (
	SubagentStatusActive SubagentStatus = "active"
	SubagentStatusClosed SubagentStatus = "closed"
)

// A reasoning summary content part.
type SummaryText struct {
	// The reasoning summary text.
	Text string `json:"text" api:"required"`
	// The content type. Always `summary_text`.
	Type constant.SummaryText `json:"type" default:"summary_text"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Text        respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r SummaryText) RawJSON() string { return r.JSON.raw }
func (r *SummaryText) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// TextFormatUnion contains all possible properties and values from
// [TextFormatText], [TextFormatJSONSchema].
//
// Use the [TextFormatUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type TextFormatUnion struct {
	// Any of "text", "json_schema".
	Type string `json:"type"`
	// This field is from variant [TextFormatJSONSchema].
	Schema map[string]any `json:"schema"`
	JSON   struct {
		Type   respjson.Field
		Schema respjson.Field
		raw    string
	} `json:"-"`
}

// anyTextFormat is implemented by each variant of [TextFormatUnion] to add type
// safety for the return type of [TextFormatUnion.AsAny]
type anyTextFormat interface {
	implTextFormatUnion()
}

func (TextFormatText) implTextFormatUnion()       {}
func (TextFormatJSONSchema) implTextFormatUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := TextFormatUnion.AsAny().(type) {
//	case openai.TextFormatText:
//	case openai.TextFormatJSONSchema:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u TextFormatUnion) AsAny() anyTextFormat {
	switch u.Type {
	case "text":
		return u.AsText()
	case "json_schema":
		return u.AsJSONSchema()
	}
	return nil
}

func (u TextFormatUnion) AsText() (v TextFormatText) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u TextFormatUnion) AsJSONSchema() (v TextFormatJSONSchema) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u TextFormatUnion) RawJSON() string { return u.JSON.raw }

func (r *TextFormatUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Generates ordinary text without a structured-output constraint.
type TextFormatText struct {
	// The type of the object. Always `text`.
	Type constant.Text `json:"type" default:"text"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r TextFormatText) RawJSON() string { return r.JSON.raw }
func (r *TextFormatText) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Constrains generated text to a JSON Schema.
type TextFormatJSONSchema struct {
	// The JSON Schema that generated text must match.
	Schema map[string]any `json:"schema" api:"required"`
	// The type of the object. Always `json_schema`.
	Type constant.JSONSchema `json:"type" default:"json_schema"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Schema      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r TextFormatJSONSchema) RawJSON() string { return r.JSON.raw }
func (r *TextFormatJSONSchema) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func TextFormatParamOfParamJSONSchema(schema map[string]any) TextFormatParamUnion {
	var paramJSONSchema TextFormatParamJSONSchema
	paramJSONSchema.Schema = schema
	return TextFormatParamUnion{OfParamJSONSchema: &paramJSONSchema}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type TextFormatParamUnion struct {
	OfParamText       *TextFormatParamText       `json:",omitzero,inline"`
	OfParamJSONSchema *TextFormatParamJSONSchema `json:",omitzero,inline"`
	paramUnion
}

func (u TextFormatParamUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfParamText, u.OfParamJSONSchema)
}
func (u *TextFormatParamUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u TextFormatParamUnion) GetSchema() map[string]any {
	if vt := u.OfParamJSONSchema; vt != nil {
		return vt.Schema
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u TextFormatParamUnion) GetType() *string {
	if vt := u.OfParamText; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfParamJSONSchema; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[TextFormatParamUnion](
		"type",
		apijson.Discriminator[TextFormatParamText]("text"),
		apijson.Discriminator[TextFormatParamJSONSchema]("json_schema"),
	)
}

func NewTextFormatParamText() TextFormatParamText {
	return TextFormatParamText{
		Type: "text",
	}
}

// Generates ordinary text without a structured-output constraint.
//
// This struct has a constant value, construct it with [NewTextFormatParamText].
type TextFormatParamText struct {
	// The type of the object. Always `text`.
	Type constant.Text `json:"type" default:"text"`
	paramObj
}

func (r TextFormatParamText) MarshalJSON() (data []byte, err error) {
	type shadow TextFormatParamText
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *TextFormatParamText) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Constrains generated text to a JSON Schema.
//
// The properties Schema, Type are required.
type TextFormatParamJSONSchema struct {
	// The JSON Schema that generated text must match.
	Schema map[string]any `json:"schema,omitzero" api:"required"`
	// The type of the object. Always `json_schema`.
	//
	// This field can be elided, and will marshal its zero value as "json_schema".
	Type constant.JSONSchema `json:"type" default:"json_schema"`
	paramObj
}

func (r TextFormatParamJSONSchema) MarshalJSON() (data []byte, err error) {
	type shadow TextFormatParamJSONSchema
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *TextFormatParamJSONSchema) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Recorded token usage for a session or turn. Usage is best effort and may change.
type TokenUsage struct {
	// The number of input tokens used by the agent.
	InputTokens int64 `json:"input_tokens" api:"required"`
	// A breakdown of the agent's input token usage.
	InputTokensDetails TokenUsageInputTokensDetails `json:"input_tokens_details" api:"required"`
	// The number of output tokens generated by the agent.
	OutputTokens int64 `json:"output_tokens" api:"required"`
	// A breakdown of the agent's output token usage.
	OutputTokensDetails TokenUsageOutputTokensDetails `json:"output_tokens_details" api:"required"`
	// The total number of input and output tokens used by the agent.
	TotalTokens int64 `json:"total_tokens" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		InputTokens         respjson.Field
		InputTokensDetails  respjson.Field
		OutputTokens        respjson.Field
		OutputTokensDetails respjson.Field
		TotalTokens         respjson.Field
		ExtraFields         map[string]respjson.Field
		raw                 string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r TokenUsage) RawJSON() string { return r.JSON.raw }
func (r *TokenUsage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A breakdown of the agent's input token usage.
type TokenUsageInputTokensDetails struct {
	// The number of input tokens retrieved from the prompt cache.
	CachedTokens int64 `json:"cached_tokens" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		CachedTokens respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r TokenUsageInputTokensDetails) RawJSON() string { return r.JSON.raw }
func (r *TokenUsageInputTokensDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A breakdown of the agent's output token usage.
type TokenUsageOutputTokensDetails struct {
	// The number of output tokens used for reasoning.
	ReasoningTokens int64 `json:"reasoning_tokens" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ReasoningTokens respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r TokenUsageOutputTokensDetails) RawJSON() string { return r.JSON.raw }
func (r *TokenUsageOutputTokensDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// WebSearchActionUnion contains all possible properties and values from
// [WebSearchActionSearch], [WebSearchActionOpenPage], [WebSearchActionFindInPage],
// [WebSearchActionOther].
//
// Use the [WebSearchActionUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type WebSearchActionUnion struct {
	// This field is from variant [WebSearchActionSearch].
	Queries []string `json:"queries"`
	// This field is from variant [WebSearchActionSearch].
	Query string `json:"query"`
	// Any of "search", "open_page", "find_in_page", "other".
	Type string `json:"type"`
	URL  string `json:"url"`
	// This field is from variant [WebSearchActionFindInPage].
	Pattern string `json:"pattern"`
	JSON    struct {
		Queries respjson.Field
		Query   respjson.Field
		Type    respjson.Field
		URL     respjson.Field
		Pattern respjson.Field
		raw     string
	} `json:"-"`
}

// anyWebSearchAction is implemented by each variant of [WebSearchActionUnion] to
// add type safety for the return type of [WebSearchActionUnion.AsAny]
type anyWebSearchAction interface {
	implWebSearchActionUnion()
}

func (WebSearchActionSearch) implWebSearchActionUnion()     {}
func (WebSearchActionOpenPage) implWebSearchActionUnion()   {}
func (WebSearchActionFindInPage) implWebSearchActionUnion() {}
func (WebSearchActionOther) implWebSearchActionUnion()      {}

// Use the following switch statement to find the correct variant
//
//	switch variant := WebSearchActionUnion.AsAny().(type) {
//	case openai.WebSearchActionSearch:
//	case openai.WebSearchActionOpenPage:
//	case openai.WebSearchActionFindInPage:
//	case openai.WebSearchActionOther:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u WebSearchActionUnion) AsAny() anyWebSearchAction {
	switch u.Type {
	case "search":
		return u.AsSearch()
	case "open_page":
		return u.AsOpenPage()
	case "find_in_page":
		return u.AsFindInPage()
	case "other":
		return u.AsOther()
	}
	return nil
}

func (u WebSearchActionUnion) AsSearch() (v WebSearchActionSearch) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u WebSearchActionUnion) AsOpenPage() (v WebSearchActionOpenPage) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u WebSearchActionUnion) AsFindInPage() (v WebSearchActionFindInPage) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u WebSearchActionUnion) AsOther() (v WebSearchActionOther) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u WebSearchActionUnion) RawJSON() string { return u.JSON.raw }

func (r *WebSearchActionUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A search query or group of search queries.
type WebSearchActionSearch struct {
	// The search queries, when multiple queries were used.
	Queries []string `json:"queries" api:"required"`
	// The search query, when a single query was used.
	Query string `json:"query" api:"required"`
	// The type of the object. Always `search`.
	Type constant.Search `json:"type" default:"search"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Queries     respjson.Field
		Query       respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r WebSearchActionSearch) RawJSON() string { return r.JSON.raw }
func (r *WebSearchActionSearch) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Opens a web page.
type WebSearchActionOpenPage struct {
	// The type of the object. Always `open_page`.
	Type constant.OpenPage `json:"type" default:"open_page"`
	// The URL of the page that was opened.
	URL string `json:"url" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		URL         respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r WebSearchActionOpenPage) RawJSON() string { return r.JSON.raw }
func (r *WebSearchActionOpenPage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Finds text within a web page.
type WebSearchActionFindInPage struct {
	// The text pattern that was searched for.
	Pattern string `json:"pattern" api:"required"`
	// The type of the object. Always `find_in_page`.
	Type constant.FindInPage `json:"type" default:"find_in_page"`
	// The URL of the page that was searched.
	URL string `json:"url" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Pattern     respjson.Field
		Type        respjson.Field
		URL         respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r WebSearchActionFindInPage) RawJSON() string { return r.JSON.raw }
func (r *WebSearchActionFindInPage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Another web search action.
type WebSearchActionOther struct {
	// The type of the object. Always `other`.
	Type constant.Other `json:"type" default:"other"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r WebSearchActionOther) RawJSON() string { return r.JSON.raw }
func (r *WebSearchActionOther) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaAgentNewParams struct {
	// The model to use for the agent. The requested model name is preserved.
	Model string `json:"model" api:"required"`
	// Additional instructions appended to the agent's default base instructions. Omit
	// or set to null to add no custom instructions.
	Instructions param.Opt[string] `json:"instructions,omitzero"`
	// A human-readable name for the agent. Omission or null leaves the agent unnamed.
	Name param.Opt[string] `json:"name,omitzero"`
	// Up to 16 string key-value pairs, with keys up to 64 and values up to 512
	// characters. Omission or null defaults to an empty map.
	Metadata map[string]string `json:"metadata,omitzero"`
	// The service tier used for model requests.
	//
	// Any of "auto", "default", "flex", "priority", "fast".
	ServiceTier BetaAgentNewParamsServiceTier `json:"service_tier,omitzero"`
	// Tools available to the agent. Defaults to an empty list.
	Tools []PersistedAgentToolParamUnion `json:"tools,omitzero"`
	// Explicit configuration for creating and coordinating subagents.
	MultiAgent MultiAgentConfigParam `json:"multi_agent,omitzero"`
	// Reasoning configuration for the agent.
	Reasoning AgentReasoningParam `json:"reasoning,omitzero"`
	// Configuration for text generated by the agent.
	Text AgentTextParam `json:"text,omitzero"`
	paramObj
}

func (r BetaAgentNewParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaAgentNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaAgentNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The service tier used for model requests.
type BetaAgentNewParamsServiceTier string

const (
	BetaAgentNewParamsServiceTierAuto     BetaAgentNewParamsServiceTier = "auto"
	BetaAgentNewParamsServiceTierDefault  BetaAgentNewParamsServiceTier = "default"
	BetaAgentNewParamsServiceTierFlex     BetaAgentNewParamsServiceTier = "flex"
	BetaAgentNewParamsServiceTierPriority BetaAgentNewParamsServiceTier = "priority"
	BetaAgentNewParamsServiceTierFast     BetaAgentNewParamsServiceTier = "fast"
)

type BetaAgentUpdateParams struct {
	// Additional instructions appended to the agent's default base instructions. Omit
	// to leave unchanged.
	Instructions param.Opt[string] `json:"instructions,omitzero"`
	// A replacement name. Omit to leave unchanged, or pass null to clear it.
	Name param.Opt[string] `json:"name,omitzero"`
	// The model to use for the agent. The requested model name is preserved.
	Model param.Opt[string] `json:"model,omitzero"`
	// Replaces all metadata. Omit to leave unchanged, or pass null or {} to clear it.
	// Up to 16 string key-value pairs, with keys up to 64 and values up to 512
	// characters.
	Metadata map[string]string `json:"metadata,omitzero"`
	// The service tier used for model requests.
	//
	// Any of "auto", "default", "flex", "priority", "fast".
	ServiceTier BetaAgentUpdateParamsServiceTier `json:"service_tier,omitzero"`
	// Tools available to the agent.
	Tools []PersistedAgentToolParamUnion `json:"tools,omitzero"`
	// Explicit configuration for creating and coordinating subagents.
	MultiAgent MultiAgentConfigParam `json:"multi_agent,omitzero"`
	// Reasoning configuration for the agent.
	Reasoning AgentReasoningParam `json:"reasoning,omitzero"`
	// Configuration for text generated by the agent.
	Text AgentTextParam `json:"text,omitzero"`
	paramObj
}

func (r BetaAgentUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaAgentUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaAgentUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The service tier used for model requests.
type BetaAgentUpdateParamsServiceTier string

const (
	BetaAgentUpdateParamsServiceTierAuto     BetaAgentUpdateParamsServiceTier = "auto"
	BetaAgentUpdateParamsServiceTierDefault  BetaAgentUpdateParamsServiceTier = "default"
	BetaAgentUpdateParamsServiceTierFlex     BetaAgentUpdateParamsServiceTier = "flex"
	BetaAgentUpdateParamsServiceTierPriority BetaAgentUpdateParamsServiceTier = "priority"
	BetaAgentUpdateParamsServiceTierFast     BetaAgentUpdateParamsServiceTier = "fast"
)

type BetaAgentListParams struct {
	// Return resources after this resource ID in the selected order.
	After param.Opt[string] `query:"after,omitzero" json:"-"`
	// The maximum number of resources to return.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// The order in which resources are returned. Defaults to `desc`.
	//
	// Any of "asc", "desc".
	Order BetaAgentListParamsOrder `query:"order,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaAgentListParams]'s query parameters as `url.Values`.
func (r BetaAgentListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// The order in which resources are returned. Defaults to `desc`.
type BetaAgentListParamsOrder string

const (
	BetaAgentListParamsOrderAsc  BetaAgentListParamsOrder = "asc"
	BetaAgentListParamsOrderDesc BetaAgentListParamsOrder = "desc"
)

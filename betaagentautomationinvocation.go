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

// BetaAgentAutomationInvocationService contains methods and other services that
// help with interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaAgentAutomationInvocationService] method instead.
type BetaAgentAutomationInvocationService struct {
	Options []option.RequestOption
}

// NewBetaAgentAutomationInvocationService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewBetaAgentAutomationInvocationService(opts ...option.RequestOption) (r BetaAgentAutomationInvocationService) {
	r = BetaAgentAutomationInvocationService{}
	r.Options = requestconfig.InheritedOptions(opts...)
	return
}

// Queues a manual delivery without changing the saved schedule.
func (r *BetaAgentAutomationInvocationService) New(ctx context.Context, automationID string, opts ...option.RequestOption) (res *AutomationInvocation, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "agents=v1")}, opts...)
	if automationID == "" {
		err = errors.New("missing required automation_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("agents/automations/%s/invocations", automationID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// Retrieves a delivery record for a project-owned automation.
func (r *BetaAgentAutomationInvocationService) Get(ctx context.Context, automationID string, invocationID string, opts ...option.RequestOption) (res *AutomationInvocation, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "agents=v1")}, opts...)
	if automationID == "" {
		err = errors.New("missing required automation_id parameter")
		return nil, err
	}
	if invocationID == "" {
		err = errors.New("missing required invocation_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("agents/automations/%s/invocations/%s", automationID, invocationID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Lists delivery records for a project-owned automation.
func (r *BetaAgentAutomationInvocationService) List(ctx context.Context, automationID string, query BetaAgentAutomationInvocationListParams, opts ...option.RequestOption) (res *pagination.CursorPage[AutomationInvocation], err error) {
	var raw *http.Response
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "agents=v1"), option.WithResponseInto(&raw)}, opts...)
	if automationID == "" {
		err = errors.New("missing required automation_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("agents/automations/%s/invocations", automationID)
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

// Lists delivery records for a project-owned automation.
func (r *BetaAgentAutomationInvocationService) ListAutoPaging(ctx context.Context, automationID string, query BetaAgentAutomationInvocationListParams, opts ...option.RequestOption) *pagination.CursorPageAutoPager[AutomationInvocation] {
	return pagination.NewCursorPageAutoPager(r.List(ctx, automationID, query, opts...))
}

// One attempted delivery by an automation.
type AutomationInvocation struct {
	// The invocation ID.
	ID string `json:"id" api:"required"`
	// The automation that owns the invocation.
	AutomationID string `json:"automation_id" api:"required"`
	// The creation time as a Unix timestamp in seconds.
	CreatedAt int64 `json:"created_at" api:"required"`
	// The permanent delivery failure, or null when no failure occurred.
	Error AutomationInvocationError `json:"error" api:"required"`
	// Always `agent.automation.invocation`.
	Object constant.AgentAutomationInvocation `json:"object" default:"agent.automation.invocation"`
	// The intended occurrence as a Unix timestamp in seconds, or null for manual runs.
	ScheduledAt int64 `json:"scheduled_at" api:"required"`
	// The target or newly created session ID, or null before one is available.
	SessionID string `json:"session_id" api:"required"`
	// Why an occurrence was skipped, or null when it was not skipped.
	SkipReason constant.PendingDelivery `json:"skip_reason" default:"pending_delivery"`
	// Whether the invocation was manual or scheduled.
	Source AutomationInvocationSourceUnion `json:"source" api:"required"`
	// Delivery status; completed does not mean the agent's work has finished.
	//
	// Any of "queued", "completed", "failed", "skipped".
	Status AutomationInvocationStatus `json:"status" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID           respjson.Field
		AutomationID respjson.Field
		CreatedAt    respjson.Field
		Error        respjson.Field
		Object       respjson.Field
		ScheduledAt  respjson.Field
		SessionID    respjson.Field
		SkipReason   respjson.Field
		Source       respjson.Field
		Status       respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AutomationInvocation) RawJSON() string { return r.JSON.raw }
func (r *AutomationInvocation) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The permanent delivery failure, or null when no failure occurred.
type AutomationInvocationError struct {
	// A machine-readable error code.
	Code string `json:"code" api:"required"`
	// A public explanation of the delivery failure.
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
func (r AutomationInvocationError) RawJSON() string { return r.JSON.raw }
func (r *AutomationInvocationError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Delivery status; completed does not mean the agent's work has finished.
type AutomationInvocationStatus string

const (
	AutomationInvocationStatusQueued    AutomationInvocationStatus = "queued"
	AutomationInvocationStatusCompleted AutomationInvocationStatus = "completed"
	AutomationInvocationStatusFailed    AutomationInvocationStatus = "failed"
	AutomationInvocationStatusSkipped   AutomationInvocationStatus = "skipped"
)

// AutomationInvocationSourceUnion contains all possible properties and values from
// [AutomationInvocationSourceManual], [AutomationInvocationSourceAutomation].
//
// Use the [AutomationInvocationSourceUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type AutomationInvocationSourceUnion struct {
	// Any of "manual", "automation".
	Type string `json:"type"`
	// This field is from variant [AutomationInvocationSourceAutomation].
	Automation AutomationUnion `json:"automation"`
	JSON       struct {
		Type       respjson.Field
		Automation respjson.Field
		raw        string
	} `json:"-"`
}

// anyAutomationInvocationSource is implemented by each variant of
// [AutomationInvocationSourceUnion] to add type safety for the return type of
// [AutomationInvocationSourceUnion.AsAny]
type anyAutomationInvocationSource interface {
	implAutomationInvocationSourceUnion()
}

func (AutomationInvocationSourceManual) implAutomationInvocationSourceUnion()     {}
func (AutomationInvocationSourceAutomation) implAutomationInvocationSourceUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := AutomationInvocationSourceUnion.AsAny().(type) {
//	case openai.AutomationInvocationSourceManual:
//	case openai.AutomationInvocationSourceAutomation:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u AutomationInvocationSourceUnion) AsAny() anyAutomationInvocationSource {
	switch u.Type {
	case "manual":
		return u.AsManual()
	case "automation":
		return u.AsAutomation()
	}
	return nil
}

func (u AutomationInvocationSourceUnion) AsManual() (v AutomationInvocationSourceManual) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u AutomationInvocationSourceUnion) AsAutomation() (v AutomationInvocationSourceAutomation) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u AutomationInvocationSourceUnion) RawJSON() string { return u.JSON.raw }

func (r *AutomationInvocationSourceUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AutomationInvocationSourceManual struct {
	// The type of the object. Always `manual`.
	Type constant.Manual `json:"type" default:"manual"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AutomationInvocationSourceManual) RawJSON() string { return r.JSON.raw }
func (r *AutomationInvocationSourceManual) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AutomationInvocationSourceAutomation struct {
	// The automation snapshot for this occurrence, excluding confidential values.
	Automation AutomationUnion `json:"automation" api:"required"`
	// The type of the object. Always `automation`.
	Type constant.Automation `json:"type" default:"automation"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Automation  respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AutomationInvocationSourceAutomation) RawJSON() string { return r.JSON.raw }
func (r *AutomationInvocationSourceAutomation) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaAgentAutomationInvocationListParams struct {
	// Return resources after this resource ID in the selected order.
	After param.Opt[string] `query:"after,omitzero" json:"-"`
	// The maximum number of resources to return, between 1 and 100. Defaults to 20.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// The order in which resources are returned. Defaults to `desc`.
	//
	// Any of "asc", "desc".
	Order BetaAgentAutomationInvocationListParamsOrder `query:"order,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaAgentAutomationInvocationListParams]'s query parameters
// as `url.Values`.
func (r BetaAgentAutomationInvocationListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// The order in which resources are returned. Defaults to `desc`.
type BetaAgentAutomationInvocationListParamsOrder string

const (
	BetaAgentAutomationInvocationListParamsOrderAsc  BetaAgentAutomationInvocationListParamsOrder = "asc"
	BetaAgentAutomationInvocationListParamsOrderDesc BetaAgentAutomationInvocationListParamsOrder = "desc"
)

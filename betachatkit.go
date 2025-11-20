// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai

import (
	"encoding/json"

	"github.com/openai/openai-go/v3/internal/apijson"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/respjson"
)

// BetaChatKitService contains methods and other services that help with
// interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaChatKitService] method instead.
type BetaChatKitService struct {
	Options  []option.RequestOption
	Sessions BetaChatKitSessionService
	Threads  BetaChatKitThreadService
}

// NewBetaChatKitService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewBetaChatKitService(opts ...option.RequestOption) (r BetaChatKitService) {
	r = BetaChatKitService{}
	r.Options = opts
	r.Sessions = NewBetaChatKitSessionService(opts...)
	r.Threads = NewBetaChatKitThreadService(opts...)
	return
}

// Workflow metadata and state returned for the session.
type ChatKitWorkflow struct {
	// Identifier of the workflow backing the session.
	ID string `json:"id,required"`
	// State variable key-value pairs applied when invoking the workflow. Defaults to
	// null when no overrides were provided.
	StateVariables map[string]ChatKitWorkflowStateVariableUnion `json:"state_variables,required"`
	// Tracing settings applied to the workflow.
	Tracing ChatKitWorkflowTracing `json:"tracing,required"`
	// Specific workflow version used for the session. Defaults to null when using the
	// latest deployment.
	Version string `json:"version,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID             respjson.Field
		StateVariables respjson.Field
		Tracing        respjson.Field
		Version        respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ChatKitWorkflow) RawJSON() string { return r.JSON.raw }
func (r *ChatKitWorkflow) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ChatKitWorkflowStateVariableUnion contains all possible properties and values
// from [string], [bool], [float64].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfString OfBool OfFloat]
type ChatKitWorkflowStateVariableUnion struct {
	// This field will be present if the value is a [string] instead of an object.
	OfString string `json:",inline"`
	// This field will be present if the value is a [bool] instead of an object.
	OfBool bool `json:",inline"`
	// This field will be present if the value is a [float64] instead of an object.
	OfFloat float64 `json:",inline"`
	JSON    struct {
		OfString respjson.Field
		OfBool   respjson.Field
		OfFloat  respjson.Field
		raw      string
	} `json:"-"`
}

func (u ChatKitWorkflowStateVariableUnion) AsString() (v string) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ChatKitWorkflowStateVariableUnion) AsBool() (v bool) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ChatKitWorkflowStateVariableUnion) AsFloat() (v float64) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ChatKitWorkflowStateVariableUnion) RawJSON() string { return u.JSON.raw }

func (r *ChatKitWorkflowStateVariableUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Tracing settings applied to the workflow.
type ChatKitWorkflowTracing struct {
	// Indicates whether tracing is enabled.
	Enabled bool `json:"enabled,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Enabled     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ChatKitWorkflowTracing) RawJSON() string { return r.JSON.raw }
func (r *ChatKitWorkflowTracing) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

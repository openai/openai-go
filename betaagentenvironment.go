// File generated from our OpenAPI spec by Castiron. See CONTRIBUTING.md for details.

package openai

import (
	"context"
	"errors"
	"net/http"
	"slices"

	"github.com/openai/openai-go/v3/internal/apijson"
	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
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
	// Any of "pending", "connected", "disconnected", "expired", "failed".
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
	EnvironmentInfoStatusConnected    EnvironmentInfoStatus = "connected"
	EnvironmentInfoStatusDisconnected EnvironmentInfoStatus = "disconnected"
	EnvironmentInfoStatusExpired      EnvironmentInfoStatus = "expired"
	EnvironmentInfoStatusFailed       EnvironmentInfoStatus = "failed"
)

// Whether the environment is hosted by OpenAI or by the application.
type EnvironmentInfoType string

const (
	EnvironmentInfoTypeOpenAIHosted EnvironmentInfoType = "openai_hosted"
	EnvironmentInfoTypeSelfHosted   EnvironmentInfoType = "self_hosted"
)

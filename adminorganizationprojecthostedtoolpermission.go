// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai

import (
	"context"
	"errors"
	"net/http"
	"slices"

	"github.com/openai/openai-go/v3/internal/apijson"
	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/packages/respjson"
)

// AdminOrganizationProjectHostedToolPermissionService contains methods and other
// services that help with interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAdminOrganizationProjectHostedToolPermissionService] method instead.
type AdminOrganizationProjectHostedToolPermissionService struct {
	Options []option.RequestOption
}

// NewAdminOrganizationProjectHostedToolPermissionService generates a new service
// that applies the given options to each request. These options are applied after
// the parent client's options (if there is one), and before any request-specific
// options.
func NewAdminOrganizationProjectHostedToolPermissionService(opts ...option.RequestOption) (r AdminOrganizationProjectHostedToolPermissionService) {
	r = AdminOrganizationProjectHostedToolPermissionService{}
	r.Options = opts
	return
}

// Returns hosted tool permissions for a project.
func (r *AdminOrganizationProjectHostedToolPermissionService) Get(ctx context.Context, projectID string, opts ...option.RequestOption) (res *ProjectHostedToolPermissions, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if projectID == "" {
		err = errors.New("missing required project_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("organization/projects/%s/hosted_tool_permissions", projectID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Updates hosted tool permissions for a project.
func (r *AdminOrganizationProjectHostedToolPermissionService) Update(ctx context.Context, projectID string, body AdminOrganizationProjectHostedToolPermissionUpdateParams, opts ...option.RequestOption) (res *ProjectHostedToolPermissions, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if projectID == "" {
		err = errors.New("missing required project_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("organization/projects/%s/hosted_tool_permissions", projectID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Represents hosted tool permissions for a project.
type ProjectHostedToolPermissions struct {
	// Permission state for a single hosted tool on a project.
	CodeInterpreter ProjectHostedToolPermissionsCodeInterpreter `json:"code_interpreter" api:"required"`
	// Permission state for a single hosted tool on a project.
	FileSearch ProjectHostedToolPermissionsFileSearch `json:"file_search" api:"required"`
	// Permission state for a single hosted tool on a project.
	ImageGeneration ProjectHostedToolPermissionsImageGeneration `json:"image_generation" api:"required"`
	// Permission state for a single hosted tool on a project.
	Mcp ProjectHostedToolPermissionsMcp `json:"mcp" api:"required"`
	// Permission state for a single hosted tool on a project.
	WebSearch ProjectHostedToolPermissionsWebSearch `json:"web_search" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		CodeInterpreter respjson.Field
		FileSearch      respjson.Field
		ImageGeneration respjson.Field
		Mcp             respjson.Field
		WebSearch       respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProjectHostedToolPermissions) RawJSON() string { return r.JSON.raw }
func (r *ProjectHostedToolPermissions) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Permission state for a single hosted tool on a project.
type ProjectHostedToolPermissionsCodeInterpreter struct {
	// Whether the hosted tool is enabled for the project.
	Enabled bool `json:"enabled" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Enabled     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProjectHostedToolPermissionsCodeInterpreter) RawJSON() string { return r.JSON.raw }
func (r *ProjectHostedToolPermissionsCodeInterpreter) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Permission state for a single hosted tool on a project.
type ProjectHostedToolPermissionsFileSearch struct {
	// Whether the hosted tool is enabled for the project.
	Enabled bool `json:"enabled" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Enabled     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProjectHostedToolPermissionsFileSearch) RawJSON() string { return r.JSON.raw }
func (r *ProjectHostedToolPermissionsFileSearch) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Permission state for a single hosted tool on a project.
type ProjectHostedToolPermissionsImageGeneration struct {
	// Whether the hosted tool is enabled for the project.
	Enabled bool `json:"enabled" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Enabled     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProjectHostedToolPermissionsImageGeneration) RawJSON() string { return r.JSON.raw }
func (r *ProjectHostedToolPermissionsImageGeneration) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Permission state for a single hosted tool on a project.
type ProjectHostedToolPermissionsMcp struct {
	// Whether the hosted tool is enabled for the project.
	Enabled bool `json:"enabled" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Enabled     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProjectHostedToolPermissionsMcp) RawJSON() string { return r.JSON.raw }
func (r *ProjectHostedToolPermissionsMcp) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Permission state for a single hosted tool on a project.
type ProjectHostedToolPermissionsWebSearch struct {
	// Whether the hosted tool is enabled for the project.
	Enabled bool `json:"enabled" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Enabled     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProjectHostedToolPermissionsWebSearch) RawJSON() string { return r.JSON.raw }
func (r *ProjectHostedToolPermissionsWebSearch) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationProjectHostedToolPermissionUpdateParams struct {
	// The code interpreter permission update.
	CodeInterpreter AdminOrganizationProjectHostedToolPermissionUpdateParamsCodeInterpreter `json:"code_interpreter,omitzero"`
	// The file search permission update.
	FileSearch AdminOrganizationProjectHostedToolPermissionUpdateParamsFileSearch `json:"file_search,omitzero"`
	// The image generation permission update.
	ImageGeneration AdminOrganizationProjectHostedToolPermissionUpdateParamsImageGeneration `json:"image_generation,omitzero"`
	// The MCP permission update.
	Mcp AdminOrganizationProjectHostedToolPermissionUpdateParamsMcp `json:"mcp,omitzero"`
	// The web search permission update.
	WebSearch AdminOrganizationProjectHostedToolPermissionUpdateParamsWebSearch `json:"web_search,omitzero"`
	paramObj
}

func (r AdminOrganizationProjectHostedToolPermissionUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationProjectHostedToolPermissionUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationProjectHostedToolPermissionUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The code interpreter permission update.
//
// The property Enabled is required.
type AdminOrganizationProjectHostedToolPermissionUpdateParamsCodeInterpreter struct {
	// Whether to enable the hosted tool for the project.
	Enabled bool `json:"enabled" api:"required"`
	paramObj
}

func (r AdminOrganizationProjectHostedToolPermissionUpdateParamsCodeInterpreter) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationProjectHostedToolPermissionUpdateParamsCodeInterpreter
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationProjectHostedToolPermissionUpdateParamsCodeInterpreter) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The file search permission update.
//
// The property Enabled is required.
type AdminOrganizationProjectHostedToolPermissionUpdateParamsFileSearch struct {
	// Whether to enable the hosted tool for the project.
	Enabled bool `json:"enabled" api:"required"`
	paramObj
}

func (r AdminOrganizationProjectHostedToolPermissionUpdateParamsFileSearch) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationProjectHostedToolPermissionUpdateParamsFileSearch
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationProjectHostedToolPermissionUpdateParamsFileSearch) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The image generation permission update.
//
// The property Enabled is required.
type AdminOrganizationProjectHostedToolPermissionUpdateParamsImageGeneration struct {
	// Whether to enable the hosted tool for the project.
	Enabled bool `json:"enabled" api:"required"`
	paramObj
}

func (r AdminOrganizationProjectHostedToolPermissionUpdateParamsImageGeneration) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationProjectHostedToolPermissionUpdateParamsImageGeneration
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationProjectHostedToolPermissionUpdateParamsImageGeneration) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The MCP permission update.
//
// The property Enabled is required.
type AdminOrganizationProjectHostedToolPermissionUpdateParamsMcp struct {
	// Whether to enable the hosted tool for the project.
	Enabled bool `json:"enabled" api:"required"`
	paramObj
}

func (r AdminOrganizationProjectHostedToolPermissionUpdateParamsMcp) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationProjectHostedToolPermissionUpdateParamsMcp
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationProjectHostedToolPermissionUpdateParamsMcp) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The web search permission update.
//
// The property Enabled is required.
type AdminOrganizationProjectHostedToolPermissionUpdateParamsWebSearch struct {
	// Whether to enable the hosted tool for the project.
	Enabled bool `json:"enabled" api:"required"`
	paramObj
}

func (r AdminOrganizationProjectHostedToolPermissionUpdateParamsWebSearch) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationProjectHostedToolPermissionUpdateParamsWebSearch
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationProjectHostedToolPermissionUpdateParamsWebSearch) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

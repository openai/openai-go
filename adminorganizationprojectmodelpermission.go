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
	"github.com/openai/openai-go/v3/shared/constant"
)

// AdminOrganizationProjectModelPermissionService contains methods and other
// services that help with interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAdminOrganizationProjectModelPermissionService] method instead.
type AdminOrganizationProjectModelPermissionService struct {
	Options []option.RequestOption
}

// NewAdminOrganizationProjectModelPermissionService generates a new service that
// applies the given options to each request. These options are applied after the
// parent client's options (if there is one), and before any request-specific
// options.
func NewAdminOrganizationProjectModelPermissionService(opts ...option.RequestOption) (r AdminOrganizationProjectModelPermissionService) {
	r = AdminOrganizationProjectModelPermissionService{}
	r.Options = opts
	return
}

// Returns model permissions for a project.
func (r *AdminOrganizationProjectModelPermissionService) Get(ctx context.Context, projectID string, opts ...option.RequestOption) (res *ProjectModelPermissions, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if projectID == "" {
		err = errors.New("missing required project_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("organization/projects/%s/model_permissions", projectID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Updates model permissions for a project.
func (r *AdminOrganizationProjectModelPermissionService) Update(ctx context.Context, projectID string, body AdminOrganizationProjectModelPermissionUpdateParams, opts ...option.RequestOption) (res *ProjectModelPermissions, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if projectID == "" {
		err = errors.New("missing required project_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("organization/projects/%s/model_permissions", projectID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Deletes model permissions for a project.
func (r *AdminOrganizationProjectModelPermissionService) Delete(ctx context.Context, projectID string, opts ...option.RequestOption) (res *ProjectModelPermissionsDeleted, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if projectID == "" {
		err = errors.New("missing required project_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("organization/projects/%s/model_permissions", projectID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Represents the model allowlist or denylist policy for a project.
type ProjectModelPermissions struct {
	// Whether the project uses an allowlist or a denylist.
	//
	// Any of "allow_list", "deny_list".
	Mode ProjectModelPermissionsMode `json:"mode" api:"required"`
	// The model IDs included in the model permissions policy.
	ModelIDs []string `json:"model_ids" api:"required"`
	// The object type, which is always `project.model_permissions`.
	Object constant.ProjectModelPermissions `json:"object" default:"project.model_permissions"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Mode        respjson.Field
		ModelIDs    respjson.Field
		Object      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProjectModelPermissions) RawJSON() string { return r.JSON.raw }
func (r *ProjectModelPermissions) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Whether the project uses an allowlist or a denylist.
type ProjectModelPermissionsMode string

const (
	ProjectModelPermissionsModeAllowList ProjectModelPermissionsMode = "allow_list"
	ProjectModelPermissionsModeDenyList  ProjectModelPermissionsMode = "deny_list"
)

// Confirmation payload returned after deleting project model permissions.
type ProjectModelPermissionsDeleted struct {
	// Whether the project model permissions were deleted.
	Deleted bool `json:"deleted" api:"required"`
	// The object type, which is always `project.model_permissions.deleted`.
	Object constant.ProjectModelPermissionsDeleted `json:"object" default:"project.model_permissions.deleted"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Deleted     respjson.Field
		Object      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProjectModelPermissionsDeleted) RawJSON() string { return r.JSON.raw }
func (r *ProjectModelPermissionsDeleted) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationProjectModelPermissionUpdateParams struct {
	// The model permissions mode to apply.
	//
	// Any of "allow_list", "deny_list".
	Mode AdminOrganizationProjectModelPermissionUpdateParamsMode `json:"mode,omitzero" api:"required"`
	// The model IDs included in this permissions policy.
	ModelIDs []string `json:"model_ids,omitzero" api:"required"`
	paramObj
}

func (r AdminOrganizationProjectModelPermissionUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationProjectModelPermissionUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationProjectModelPermissionUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The model permissions mode to apply.
type AdminOrganizationProjectModelPermissionUpdateParamsMode string

const (
	AdminOrganizationProjectModelPermissionUpdateParamsModeAllowList AdminOrganizationProjectModelPermissionUpdateParamsMode = "allow_list"
	AdminOrganizationProjectModelPermissionUpdateParamsModeDenyList  AdminOrganizationProjectModelPermissionUpdateParamsMode = "deny_list"
)

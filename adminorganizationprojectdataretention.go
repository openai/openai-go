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

// AdminOrganizationProjectDataRetentionService contains methods and other services
// that help with interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAdminOrganizationProjectDataRetentionService] method instead.
type AdminOrganizationProjectDataRetentionService struct {
	Options []option.RequestOption
}

// NewAdminOrganizationProjectDataRetentionService generates a new service that
// applies the given options to each request. These options are applied after the
// parent client's options (if there is one), and before any request-specific
// options.
func NewAdminOrganizationProjectDataRetentionService(opts ...option.RequestOption) (r AdminOrganizationProjectDataRetentionService) {
	r = AdminOrganizationProjectDataRetentionService{}
	r.Options = opts
	return
}

// Retrieves project data retention controls.
func (r *AdminOrganizationProjectDataRetentionService) Get(ctx context.Context, projectID string, opts ...option.RequestOption) (res *ProjectDataRetention, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if projectID == "" {
		err = errors.New("missing required project_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("organization/projects/%s/data_retention", projectID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Updates project data retention controls.
func (r *AdminOrganizationProjectDataRetentionService) Update(ctx context.Context, projectID string, body AdminOrganizationProjectDataRetentionUpdateParams, opts ...option.RequestOption) (res *ProjectDataRetention, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if projectID == "" {
		err = errors.New("missing required project_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("organization/projects/%s/data_retention", projectID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Represents a project's data retention control setting.
type ProjectDataRetention struct {
	// The object type, which is always `project.data_retention`.
	Object constant.ProjectDataRetention `json:"object" default:"project.data_retention"`
	// The configured project data retention type.
	//
	// Any of "organization_default", "none", "zero_data_retention",
	// "modified_abuse_monitoring", "enhanced_zero_data_retention",
	// "enhanced_modified_abuse_monitoring".
	Type ProjectDataRetentionType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Object      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProjectDataRetention) RawJSON() string { return r.JSON.raw }
func (r *ProjectDataRetention) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The configured project data retention type.
type ProjectDataRetentionType string

const (
	ProjectDataRetentionTypeOrganizationDefault             ProjectDataRetentionType = "organization_default"
	ProjectDataRetentionTypeNone                            ProjectDataRetentionType = "none"
	ProjectDataRetentionTypeZeroDataRetention               ProjectDataRetentionType = "zero_data_retention"
	ProjectDataRetentionTypeModifiedAbuseMonitoring         ProjectDataRetentionType = "modified_abuse_monitoring"
	ProjectDataRetentionTypeEnhancedZeroDataRetention       ProjectDataRetentionType = "enhanced_zero_data_retention"
	ProjectDataRetentionTypeEnhancedModifiedAbuseMonitoring ProjectDataRetentionType = "enhanced_modified_abuse_monitoring"
)

type AdminOrganizationProjectDataRetentionUpdateParams struct {
	// The desired project data retention type.
	//
	// Any of "organization_default", "none", "zero_data_retention",
	// "modified_abuse_monitoring", "enhanced_zero_data_retention",
	// "enhanced_modified_abuse_monitoring".
	RetentionType AdminOrganizationProjectDataRetentionUpdateParamsRetentionType `json:"retention_type,omitzero" api:"required"`
	paramObj
}

func (r AdminOrganizationProjectDataRetentionUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationProjectDataRetentionUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationProjectDataRetentionUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The desired project data retention type.
type AdminOrganizationProjectDataRetentionUpdateParamsRetentionType string

const (
	AdminOrganizationProjectDataRetentionUpdateParamsRetentionTypeOrganizationDefault             AdminOrganizationProjectDataRetentionUpdateParamsRetentionType = "organization_default"
	AdminOrganizationProjectDataRetentionUpdateParamsRetentionTypeNone                            AdminOrganizationProjectDataRetentionUpdateParamsRetentionType = "none"
	AdminOrganizationProjectDataRetentionUpdateParamsRetentionTypeZeroDataRetention               AdminOrganizationProjectDataRetentionUpdateParamsRetentionType = "zero_data_retention"
	AdminOrganizationProjectDataRetentionUpdateParamsRetentionTypeModifiedAbuseMonitoring         AdminOrganizationProjectDataRetentionUpdateParamsRetentionType = "modified_abuse_monitoring"
	AdminOrganizationProjectDataRetentionUpdateParamsRetentionTypeEnhancedZeroDataRetention       AdminOrganizationProjectDataRetentionUpdateParamsRetentionType = "enhanced_zero_data_retention"
	AdminOrganizationProjectDataRetentionUpdateParamsRetentionTypeEnhancedModifiedAbuseMonitoring AdminOrganizationProjectDataRetentionUpdateParamsRetentionType = "enhanced_modified_abuse_monitoring"
)

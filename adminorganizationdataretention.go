// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai

import (
	"context"
	"net/http"
	"slices"

	"github.com/openai/openai-go/v3/internal/apijson"
	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/packages/respjson"
	"github.com/openai/openai-go/v3/shared/constant"
)

// AdminOrganizationDataRetentionService contains methods and other services that
// help with interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAdminOrganizationDataRetentionService] method instead.
type AdminOrganizationDataRetentionService struct {
	Options []option.RequestOption
}

// NewAdminOrganizationDataRetentionService generates a new service that applies
// the given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewAdminOrganizationDataRetentionService(opts ...option.RequestOption) (r AdminOrganizationDataRetentionService) {
	r = AdminOrganizationDataRetentionService{}
	r.Options = opts
	return
}

// Retrieves organization data retention controls.
func (r *AdminOrganizationDataRetentionService) Get(ctx context.Context, opts ...option.RequestOption) (res *OrganizationDataRetention, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	path := "organization/data_retention"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Updates organization data retention controls.
func (r *AdminOrganizationDataRetentionService) Update(ctx context.Context, body AdminOrganizationDataRetentionUpdateParams, opts ...option.RequestOption) (res *OrganizationDataRetention, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	path := "organization/data_retention"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Represents the organization's data retention control setting.
type OrganizationDataRetention struct {
	// The object type, which is always `organization.data_retention`.
	Object constant.OrganizationDataRetention `json:"object" default:"organization.data_retention"`
	// The configured organization data retention type.
	//
	// Any of "zero_data_retention", "modified_abuse_monitoring",
	// "enhanced_zero_data_retention", "enhanced_modified_abuse_monitoring".
	Type OrganizationDataRetentionType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Object      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r OrganizationDataRetention) RawJSON() string { return r.JSON.raw }
func (r *OrganizationDataRetention) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The configured organization data retention type.
type OrganizationDataRetentionType string

const (
	OrganizationDataRetentionTypeZeroDataRetention               OrganizationDataRetentionType = "zero_data_retention"
	OrganizationDataRetentionTypeModifiedAbuseMonitoring         OrganizationDataRetentionType = "modified_abuse_monitoring"
	OrganizationDataRetentionTypeEnhancedZeroDataRetention       OrganizationDataRetentionType = "enhanced_zero_data_retention"
	OrganizationDataRetentionTypeEnhancedModifiedAbuseMonitoring OrganizationDataRetentionType = "enhanced_modified_abuse_monitoring"
)

type AdminOrganizationDataRetentionUpdateParams struct {
	// The desired organization data retention type.
	//
	// Any of "zero_data_retention", "modified_abuse_monitoring",
	// "enhanced_zero_data_retention", "enhanced_modified_abuse_monitoring".
	RetentionType AdminOrganizationDataRetentionUpdateParamsRetentionType `json:"retention_type,omitzero" api:"required"`
	paramObj
}

func (r AdminOrganizationDataRetentionUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationDataRetentionUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationDataRetentionUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The desired organization data retention type.
type AdminOrganizationDataRetentionUpdateParamsRetentionType string

const (
	AdminOrganizationDataRetentionUpdateParamsRetentionTypeZeroDataRetention               AdminOrganizationDataRetentionUpdateParamsRetentionType = "zero_data_retention"
	AdminOrganizationDataRetentionUpdateParamsRetentionTypeModifiedAbuseMonitoring         AdminOrganizationDataRetentionUpdateParamsRetentionType = "modified_abuse_monitoring"
	AdminOrganizationDataRetentionUpdateParamsRetentionTypeEnhancedZeroDataRetention       AdminOrganizationDataRetentionUpdateParamsRetentionType = "enhanced_zero_data_retention"
	AdminOrganizationDataRetentionUpdateParamsRetentionTypeEnhancedModifiedAbuseMonitoring AdminOrganizationDataRetentionUpdateParamsRetentionType = "enhanced_modified_abuse_monitoring"
)

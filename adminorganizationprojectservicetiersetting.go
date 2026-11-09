// File generated from our OpenAPI spec by Castiron. See CONTRIBUTING.md for details.

package openai

import (
	"context"
	"encoding/json"
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

// AdminOrganizationProjectServiceTierSettingService contains methods and other
// services that help with interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAdminOrganizationProjectServiceTierSettingService] method instead.
type AdminOrganizationProjectServiceTierSettingService struct {
	Options []option.RequestOption
}

// NewAdminOrganizationProjectServiceTierSettingService generates a new service
// that applies the given options to each request. These options are applied after
// the parent client's options (if there is one), and before any request-specific
// options.
func NewAdminOrganizationProjectServiceTierSettingService(opts ...option.RequestOption) (r AdminOrganizationProjectServiceTierSettingService) {
	r = AdminOrganizationProjectServiceTierSettingService{}
	r.Options = requestconfig.InheritedOptions(opts...)
	return
}

// Retrieve the configured default and enabled service tiers for a project.
func (r *AdminOrganizationProjectServiceTierSettingService) Get(ctx context.Context, projectID string, opts ...option.RequestOption) (res *ProjectServiceTierSettings, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if projectID == "" {
		err = errors.New("missing required project_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("organization/projects/%s/service_tier_settings", projectID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update the default or enabled service tiers atomically. Supply at least one
// field. Omitted fields keep their current values.
func (r *AdminOrganizationProjectServiceTierSettingService) Update(ctx context.Context, projectID string, body AdminOrganizationProjectServiceTierSettingUpdateParams, opts ...option.RequestOption) (res *ProjectServiceTierSettings, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if projectID == "" {
		err = errors.New("missing required project_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("organization/projects/%s/service_tier_settings", projectID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

type ProjectOptionalServiceTier string

const (
	ProjectOptionalServiceTierFlex      ProjectOptionalServiceTier = "flex"
	ProjectOptionalServiceTierFast      ProjectOptionalServiceTier = "fast"
	ProjectOptionalServiceTierUltrafast ProjectOptionalServiceTier = "ultrafast"
)

type ProjectServiceTierSettings struct {
	// Default tier for requests that omit service_tier or use auto.
	//
	// Any of "standard", "fast", "scale".
	Default ProjectServiceTierSettingsDefault `json:"default" api:"required"`
	// Additional service tiers enabled for this project. "all" includes future tiers.
	// An empty array allows only Standard. Organization policy and eligibility still
	// apply.
	EnabledTiers ProjectServiceTierSettingsEnabledTiersUnion `json:"enabled_tiers" api:"required"`
	Object       constant.ProjectServiceTierSettings         `json:"object" default:"project.service_tier_settings"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Default      respjson.Field
		EnabledTiers respjson.Field
		Object       respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProjectServiceTierSettings) RawJSON() string { return r.JSON.raw }
func (r *ProjectServiceTierSettings) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Default tier for requests that omit service_tier or use auto.
type ProjectServiceTierSettingsDefault string

const (
	ProjectServiceTierSettingsDefaultStandard ProjectServiceTierSettingsDefault = "standard"
	ProjectServiceTierSettingsDefaultFast     ProjectServiceTierSettingsDefault = "fast"
	ProjectServiceTierSettingsDefaultScale    ProjectServiceTierSettingsDefault = "scale"
)

// ProjectServiceTierSettingsEnabledTiersUnion contains all possible properties and
// values from [constant.All], [[]ProjectOptionalServiceTier].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfAll OfProjectServiceTierSettingsEnabledTiers]
type ProjectServiceTierSettingsEnabledTiersUnion struct {
	// This field will be present if the value is a [constant.All] instead of an
	// object.
	OfAll constant.All `json:",inline"`
	// This field will be present if the value is a [[]ProjectOptionalServiceTier]
	// instead of an object.
	OfProjectServiceTierSettingsEnabledTiers []ProjectOptionalServiceTier `json:",inline"`
	JSON                                     struct {
		OfAll                                    respjson.Field
		OfProjectServiceTierSettingsEnabledTiers respjson.Field
		raw                                      string
	} `json:"-"`
}

func (u ProjectServiceTierSettingsEnabledTiersUnion) AsAll() (v constant.All) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ProjectServiceTierSettingsEnabledTiersUnion) AsProjectServiceTierSettingsEnabledTiers() (v []ProjectOptionalServiceTier) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ProjectServiceTierSettingsEnabledTiersUnion) RawJSON() string { return u.JSON.raw }

func (r *ProjectServiceTierSettingsEnabledTiersUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationProjectServiceTierSettingUpdateParams struct {
	// Default tier for requests that omit service_tier or use auto.
	//
	// Any of "standard", "fast", "scale".
	Default AdminOrganizationProjectServiceTierSettingUpdateParamsDefault `json:"default,omitzero"`
	// Additional service tiers enabled for this project. "all" includes future tiers.
	// An empty array allows only Standard. Organization policy and eligibility still
	// apply.
	EnabledTiers AdminOrganizationProjectServiceTierSettingUpdateParamsEnabledTiersUnion `json:"enabled_tiers,omitzero"`
	paramObj
}

func (r AdminOrganizationProjectServiceTierSettingUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationProjectServiceTierSettingUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationProjectServiceTierSettingUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Default tier for requests that omit service_tier or use auto.
type AdminOrganizationProjectServiceTierSettingUpdateParamsDefault string

const (
	AdminOrganizationProjectServiceTierSettingUpdateParamsDefaultStandard AdminOrganizationProjectServiceTierSettingUpdateParamsDefault = "standard"
	AdminOrganizationProjectServiceTierSettingUpdateParamsDefaultFast     AdminOrganizationProjectServiceTierSettingUpdateParamsDefault = "fast"
	AdminOrganizationProjectServiceTierSettingUpdateParamsDefaultScale    AdminOrganizationProjectServiceTierSettingUpdateParamsDefault = "scale"
)

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type AdminOrganizationProjectServiceTierSettingUpdateParamsEnabledTiersUnion struct {
	// Construct this variant with constant.ValueOf[constant.All]()
	OfAll                             constant.All                 `json:",omitzero,inline"`
	OfProjectOptionalServiceTierArray []ProjectOptionalServiceTier `json:",omitzero,inline"`
	paramUnion
}

func (u AdminOrganizationProjectServiceTierSettingUpdateParamsEnabledTiersUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfAll, u.OfProjectOptionalServiceTierArray)
}
func (u *AdminOrganizationProjectServiceTierSettingUpdateParamsEnabledTiersUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

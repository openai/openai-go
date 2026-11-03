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

// AdminOrganizationWorkloadIdentityProviderFederationRuleService contains methods
// and other services that help with interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAdminOrganizationWorkloadIdentityProviderFederationRuleService] method
// instead.
type AdminOrganizationWorkloadIdentityProviderFederationRuleService struct {
	Options []option.RequestOption
}

// NewAdminOrganizationWorkloadIdentityProviderFederationRuleService generates a
// new service that applies the given options to each request. These options are
// applied after the parent client's options (if there is one), and before any
// request-specific options.
func NewAdminOrganizationWorkloadIdentityProviderFederationRuleService(opts ...option.RequestOption) (r AdminOrganizationWorkloadIdentityProviderFederationRuleService) {
	r = AdminOrganizationWorkloadIdentityProviderFederationRuleService{}
	r.Options = requestconfig.InheritedOptions(opts...)
	return
}

// Bind verified token attributes to an existing project service account.
func (r *AdminOrganizationWorkloadIdentityProviderFederationRuleService) New(ctx context.Context, workloadIdentityProviderID string, body AdminOrganizationWorkloadIdentityProviderFederationRuleNewParams, opts ...option.RequestOption) (res *FederationRule, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if workloadIdentityProviderID == "" {
		err = errors.New("missing required workload_identity_provider_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("organization/workload_identity_providers/%s/federation_rules", workloadIdentityProviderID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Retrieve a federation rule owned by the authenticated organization and provider.
func (r *AdminOrganizationWorkloadIdentityProviderFederationRuleService) Get(ctx context.Context, workloadIdentityProviderID string, workloadIdentityFederationRuleID string, opts ...option.RequestOption) (res *FederationRule, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if workloadIdentityProviderID == "" {
		err = errors.New("missing required workload_identity_provider_id parameter")
		return nil, err
	}
	if workloadIdentityFederationRuleID == "" {
		err = errors.New("missing required workload_identity_federation_rule_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("organization/workload_identity_providers/%s/federation_rules/%s", workloadIdentityProviderID, workloadIdentityFederationRuleID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update supplied matcher, permission, metadata, or enabled fields.
func (r *AdminOrganizationWorkloadIdentityProviderFederationRuleService) Update(ctx context.Context, workloadIdentityProviderID string, workloadIdentityFederationRuleID string, body AdminOrganizationWorkloadIdentityProviderFederationRuleUpdateParams, opts ...option.RequestOption) (res *FederationRule, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if workloadIdentityProviderID == "" {
		err = errors.New("missing required workload_identity_provider_id parameter")
		return nil, err
	}
	if workloadIdentityFederationRuleID == "" {
		err = errors.New("missing required workload_identity_federation_rule_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("organization/workload_identity_providers/%s/federation_rules/%s", workloadIdentityProviderID, workloadIdentityFederationRuleID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// List the provider's federation rules to project service accounts.
func (r *AdminOrganizationWorkloadIdentityProviderFederationRuleService) List(ctx context.Context, workloadIdentityProviderID string, query AdminOrganizationWorkloadIdentityProviderFederationRuleListParams, opts ...option.RequestOption) (res *pagination.CursorPage[FederationRule], err error) {
	var raw *http.Response
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	if workloadIdentityProviderID == "" {
		err = errors.New("missing required workload_identity_provider_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("organization/workload_identity_providers/%s/federation_rules", workloadIdentityProviderID)
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

// List the provider's federation rules to project service accounts.
func (r *AdminOrganizationWorkloadIdentityProviderFederationRuleService) ListAutoPaging(ctx context.Context, workloadIdentityProviderID string, query AdminOrganizationWorkloadIdentityProviderFederationRuleListParams, opts ...option.RequestOption) *pagination.CursorPageAutoPager[FederationRule] {
	return pagination.NewCursorPageAutoPager(r.List(ctx, workloadIdentityProviderID, query, opts...))
}

// Delete a federation rule to a project service account.
func (r *AdminOrganizationWorkloadIdentityProviderFederationRuleService) Delete(ctx context.Context, workloadIdentityProviderID string, workloadIdentityFederationRuleID string, opts ...option.RequestOption) (res *FederationRuleDeleted, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if workloadIdentityProviderID == "" {
		err = errors.New("missing required workload_identity_provider_id parameter")
		return nil, err
	}
	if workloadIdentityFederationRuleID == "" {
		err = errors.New("missing required workload_identity_federation_rule_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("organization/workload_identity_providers/%s/federation_rules/%s", workloadIdentityProviderID, workloadIdentityFederationRuleID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

type FederationRule struct {
	ID                         string                                              `json:"id" api:"required"`
	AttributeMatchers          map[string]FederationRuleAttributeMatcherUnion      `json:"attribute_matchers" api:"required"`
	CreatedAt                  int64                                               `json:"created_at" api:"required"`
	Description                string                                              `json:"description" api:"required"`
	Enabled                    bool                                                `json:"enabled" api:"required"`
	Name                       string                                              `json:"name" api:"required"`
	Object                     constant.OrganizationWorkloadIdentityFederationRule `json:"object" default:"organization.workload_identity_federation_rule"`
	Permissions                FederationRulePermissionsUnion                      `json:"permissions" api:"required"`
	ProjectID                  string                                              `json:"project_id" api:"required"`
	ServiceAccountID           string                                              `json:"service_account_id" api:"required"`
	UpdatedAt                  int64                                               `json:"updated_at" api:"required"`
	WorkloadIdentityProviderID string                                              `json:"workload_identity_provider_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                         respjson.Field
		AttributeMatchers          respjson.Field
		CreatedAt                  respjson.Field
		Description                respjson.Field
		Enabled                    respjson.Field
		Name                       respjson.Field
		Object                     respjson.Field
		Permissions                respjson.Field
		ProjectID                  respjson.Field
		ServiceAccountID           respjson.Field
		UpdatedAt                  respjson.Field
		WorkloadIdentityProviderID respjson.Field
		ExtraFields                map[string]respjson.Field
		raw                        string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r FederationRule) RawJSON() string { return r.JSON.raw }
func (r *FederationRule) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// FederationRuleAttributeMatcherUnion contains all possible properties and values
// from [string], [bool].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfString OfBool]
type FederationRuleAttributeMatcherUnion struct {
	// This field will be present if the value is a [string] instead of an object.
	OfString string `json:",inline"`
	// This field will be present if the value is a [bool] instead of an object.
	OfBool bool `json:",inline"`
	JSON   struct {
		OfString respjson.Field
		OfBool   respjson.Field
		raw      string
	} `json:"-"`
}

func (u FederationRuleAttributeMatcherUnion) AsString() (v string) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u FederationRuleAttributeMatcherUnion) AsBool() (v bool) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u FederationRuleAttributeMatcherUnion) AsUnionMember2() (v UnionMember2) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u FederationRuleAttributeMatcherUnion) RawJSON() string { return u.JSON.raw }

func (r *FederationRuleAttributeMatcherUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// UnionMember2 contains all possible properties and values from [int64],
// [float64].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfInt OfFloat]
type UnionMember2 struct {
	// This field will be present if the value is a [int64] instead of an object.
	OfInt int64 `json:",inline"`
	// This field will be present if the value is a [float64] instead of an object.
	OfFloat float64 `json:",inline"`
	JSON    struct {
		OfInt   respjson.Field
		OfFloat respjson.Field
		raw     string
	} `json:"-"`
}

func (u UnionMember2) AsInt() (v int64) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u UnionMember2) AsFloat() (v float64) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u UnionMember2) RawJSON() string { return u.JSON.raw }

func (r *UnionMember2) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// FederationRulePermissionsUnion contains all possible properties and values from
// [FederationRulePermissionsServiceAccountDefault],
// [FederationRulePermissionsRestricted].
//
// Use the [FederationRulePermissionsUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type FederationRulePermissionsUnion struct {
	// Any of "service_account_default", "restricted".
	Type string `json:"type"`
	// This field is from variant [FederationRulePermissionsRestricted].
	Scopes []string `json:"scopes"`
	JSON   struct {
		Type   respjson.Field
		Scopes respjson.Field
		raw    string
	} `json:"-"`
}

// anyFederationRulePermissions is implemented by each variant of
// [FederationRulePermissionsUnion] to add type safety for the return type of
// [FederationRulePermissionsUnion.AsAny]
type anyFederationRulePermissions interface {
	implFederationRulePermissionsUnion()
}

func (FederationRulePermissionsServiceAccountDefault) implFederationRulePermissionsUnion() {}
func (FederationRulePermissionsRestricted) implFederationRulePermissionsUnion()            {}

// Use the following switch statement to find the correct variant
//
//	switch variant := FederationRulePermissionsUnion.AsAny().(type) {
//	case openai.FederationRulePermissionsServiceAccountDefault:
//	case openai.FederationRulePermissionsRestricted:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u FederationRulePermissionsUnion) AsAny() anyFederationRulePermissions {
	switch u.Type {
	case "service_account_default":
		return u.AsServiceAccountDefault()
	case "restricted":
		return u.AsRestricted()
	}
	return nil
}

func (u FederationRulePermissionsUnion) AsServiceAccountDefault() (v FederationRulePermissionsServiceAccountDefault) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u FederationRulePermissionsUnion) AsRestricted() (v FederationRulePermissionsRestricted) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u FederationRulePermissionsUnion) RawJSON() string { return u.JSON.raw }

func (r *FederationRulePermissionsUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type FederationRulePermissionsServiceAccountDefault struct {
	Type constant.ServiceAccountDefault `json:"type" default:"service_account_default"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r FederationRulePermissionsServiceAccountDefault) RawJSON() string { return r.JSON.raw }
func (r *FederationRulePermissionsServiceAccountDefault) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type FederationRulePermissionsRestricted struct {
	Scopes []string            `json:"scopes" api:"required"`
	Type   constant.Restricted `json:"type" default:"restricted"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Scopes      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r FederationRulePermissionsRestricted) RawJSON() string { return r.JSON.raw }
func (r *FederationRulePermissionsRestricted) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type FederationRuleDeleted struct {
	ID      string                                                     `json:"id" api:"required"`
	Deleted bool                                                       `json:"deleted" api:"required"`
	Object  constant.OrganizationWorkloadIdentityFederationRuleDeleted `json:"object" default:"organization.workload_identity_federation_rule.deleted"`
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
func (r FederationRuleDeleted) RawJSON() string { return r.JSON.raw }
func (r *FederationRuleDeleted) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationWorkloadIdentityProviderFederationRuleNewParams struct {
	AttributeMatchers map[string]AdminOrganizationWorkloadIdentityProviderFederationRuleNewParamsAttributeMatcherUnion `json:"attribute_matchers,omitzero" api:"required"`
	ProjectID         string                                                                                           `json:"project_id" api:"required"`
	ServiceAccountID  string                                                                                           `json:"service_account_id" api:"required"`
	Description       param.Opt[string]                                                                                `json:"description,omitzero"`
	Enabled           param.Opt[bool]                                                                                  `json:"enabled,omitzero"`
	Name              param.Opt[string]                                                                                `json:"name,omitzero"`
	Permissions       AdminOrganizationWorkloadIdentityProviderFederationRuleNewParamsPermissionsUnion                 `json:"permissions,omitzero"`
	paramObj
}

func (r AdminOrganizationWorkloadIdentityProviderFederationRuleNewParams) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationWorkloadIdentityProviderFederationRuleNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationWorkloadIdentityProviderFederationRuleNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type AdminOrganizationWorkloadIdentityProviderFederationRuleNewParamsAttributeMatcherUnion struct {
	OfString param.Opt[string]  `json:",omitzero,inline"`
	OfBool   param.Opt[bool]    `json:",omitzero,inline"`
	OfInt    param.Opt[int64]   `json:",omitzero,inline"`
	OfFloat  param.Opt[float64] `json:",omitzero,inline"`
	paramUnion
}

func (u AdminOrganizationWorkloadIdentityProviderFederationRuleNewParamsAttributeMatcherUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfString, u.OfBool, u.OfInt, u.OfFloat)
}
func (u *AdminOrganizationWorkloadIdentityProviderFederationRuleNewParamsAttributeMatcherUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type AdminOrganizationWorkloadIdentityProviderFederationRuleNewParamsPermissionsUnion struct {
	OfServiceAccountDefault *AdminOrganizationWorkloadIdentityProviderFederationRuleNewParamsPermissionsServiceAccountDefault `json:",omitzero,inline"`
	OfRestricted            *AdminOrganizationWorkloadIdentityProviderFederationRuleNewParamsPermissionsRestricted            `json:",omitzero,inline"`
	paramUnion
}

func (u AdminOrganizationWorkloadIdentityProviderFederationRuleNewParamsPermissionsUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfServiceAccountDefault, u.OfRestricted)
}
func (u *AdminOrganizationWorkloadIdentityProviderFederationRuleNewParamsPermissionsUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u AdminOrganizationWorkloadIdentityProviderFederationRuleNewParamsPermissionsUnion) GetScopes() []string {
	if vt := u.OfRestricted; vt != nil {
		return vt.Scopes
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AdminOrganizationWorkloadIdentityProviderFederationRuleNewParamsPermissionsUnion) GetType() *string {
	if vt := u.OfServiceAccountDefault; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfRestricted; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[AdminOrganizationWorkloadIdentityProviderFederationRuleNewParamsPermissionsUnion](
		"type",
		apijson.Discriminator[AdminOrganizationWorkloadIdentityProviderFederationRuleNewParamsPermissionsServiceAccountDefault]("service_account_default"),
		apijson.Discriminator[AdminOrganizationWorkloadIdentityProviderFederationRuleNewParamsPermissionsRestricted]("restricted"),
	)
}

func NewAdminOrganizationWorkloadIdentityProviderFederationRuleNewParamsPermissionsServiceAccountDefault() AdminOrganizationWorkloadIdentityProviderFederationRuleNewParamsPermissionsServiceAccountDefault {
	return AdminOrganizationWorkloadIdentityProviderFederationRuleNewParamsPermissionsServiceAccountDefault{
		Type: "service_account_default",
	}
}

// This struct has a constant value, construct it with
// [NewAdminOrganizationWorkloadIdentityProviderFederationRuleNewParamsPermissionsServiceAccountDefault].
type AdminOrganizationWorkloadIdentityProviderFederationRuleNewParamsPermissionsServiceAccountDefault struct {
	Type constant.ServiceAccountDefault `json:"type" default:"service_account_default"`
	paramObj
}

func (r AdminOrganizationWorkloadIdentityProviderFederationRuleNewParamsPermissionsServiceAccountDefault) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationWorkloadIdentityProviderFederationRuleNewParamsPermissionsServiceAccountDefault
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationWorkloadIdentityProviderFederationRuleNewParamsPermissionsServiceAccountDefault) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Scopes, Type are required.
type AdminOrganizationWorkloadIdentityProviderFederationRuleNewParamsPermissionsRestricted struct {
	Scopes []string `json:"scopes,omitzero" api:"required"`
	// This field can be elided, and will marshal its zero value as "restricted".
	Type constant.Restricted `json:"type" default:"restricted"`
	paramObj
}

func (r AdminOrganizationWorkloadIdentityProviderFederationRuleNewParamsPermissionsRestricted) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationWorkloadIdentityProviderFederationRuleNewParamsPermissionsRestricted
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationWorkloadIdentityProviderFederationRuleNewParamsPermissionsRestricted) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationWorkloadIdentityProviderFederationRuleUpdateParams struct {
	Description       param.Opt[string]                                                                                   `json:"description,omitzero"`
	Name              param.Opt[string]                                                                                   `json:"name,omitzero"`
	Enabled           param.Opt[bool]                                                                                     `json:"enabled,omitzero"`
	AttributeMatchers map[string]AdminOrganizationWorkloadIdentityProviderFederationRuleUpdateParamsAttributeMatcherUnion `json:"attribute_matchers,omitzero"`
	Permissions       AdminOrganizationWorkloadIdentityProviderFederationRuleUpdateParamsPermissionsUnion                 `json:"permissions,omitzero"`
	paramObj
}

func (r AdminOrganizationWorkloadIdentityProviderFederationRuleUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationWorkloadIdentityProviderFederationRuleUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationWorkloadIdentityProviderFederationRuleUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type AdminOrganizationWorkloadIdentityProviderFederationRuleUpdateParamsAttributeMatcherUnion struct {
	OfString param.Opt[string]  `json:",omitzero,inline"`
	OfBool   param.Opt[bool]    `json:",omitzero,inline"`
	OfInt    param.Opt[int64]   `json:",omitzero,inline"`
	OfFloat  param.Opt[float64] `json:",omitzero,inline"`
	paramUnion
}

func (u AdminOrganizationWorkloadIdentityProviderFederationRuleUpdateParamsAttributeMatcherUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfString, u.OfBool, u.OfInt, u.OfFloat)
}
func (u *AdminOrganizationWorkloadIdentityProviderFederationRuleUpdateParamsAttributeMatcherUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type AdminOrganizationWorkloadIdentityProviderFederationRuleUpdateParamsPermissionsUnion struct {
	OfServiceAccountDefault                                                                               *AdminOrganizationWorkloadIdentityProviderFederationRuleUpdateParamsPermissionsDefaultPermissionsParam    `json:",omitzero,inline"`
	OfAdminOrganizationWorkloadIdentityProviderFederationRuleUpdatesPermissionsRestrictedPermissionsParam *AdminOrganizationWorkloadIdentityProviderFederationRuleUpdateParamsPermissionsRestrictedPermissionsParam `json:",omitzero,inline"`
	paramUnion
}

func (u AdminOrganizationWorkloadIdentityProviderFederationRuleUpdateParamsPermissionsUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfServiceAccountDefault, u.OfAdminOrganizationWorkloadIdentityProviderFederationRuleUpdatesPermissionsRestrictedPermissionsParam)
}
func (u *AdminOrganizationWorkloadIdentityProviderFederationRuleUpdateParamsPermissionsUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u AdminOrganizationWorkloadIdentityProviderFederationRuleUpdateParamsPermissionsUnion) GetScopes() []string {
	if vt := u.OfAdminOrganizationWorkloadIdentityProviderFederationRuleUpdatesPermissionsRestrictedPermissionsParam; vt != nil {
		return vt.Scopes
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AdminOrganizationWorkloadIdentityProviderFederationRuleUpdateParamsPermissionsUnion) GetType() *string {
	if vt := u.OfServiceAccountDefault; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfAdminOrganizationWorkloadIdentityProviderFederationRuleUpdatesPermissionsRestrictedPermissionsParam; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[AdminOrganizationWorkloadIdentityProviderFederationRuleUpdateParamsPermissionsUnion](
		"type",
		apijson.Discriminator[AdminOrganizationWorkloadIdentityProviderFederationRuleUpdateParamsPermissionsDefaultPermissionsParam]("service_account_default"),
		apijson.Discriminator[AdminOrganizationWorkloadIdentityProviderFederationRuleUpdateParamsPermissionsRestrictedPermissionsParam]("restricted"),
	)
}

func NewAdminOrganizationWorkloadIdentityProviderFederationRuleUpdateParamsPermissionsDefaultPermissionsParam() AdminOrganizationWorkloadIdentityProviderFederationRuleUpdateParamsPermissionsDefaultPermissionsParam {
	return AdminOrganizationWorkloadIdentityProviderFederationRuleUpdateParamsPermissionsDefaultPermissionsParam{
		Type: "service_account_default",
	}
}

// This struct has a constant value, construct it with
// [NewAdminOrganizationWorkloadIdentityProviderFederationRuleUpdateParamsPermissionsDefaultPermissionsParam].
type AdminOrganizationWorkloadIdentityProviderFederationRuleUpdateParamsPermissionsDefaultPermissionsParam struct {
	Type constant.ServiceAccountDefault `json:"type" default:"service_account_default"`
	paramObj
}

func (r AdminOrganizationWorkloadIdentityProviderFederationRuleUpdateParamsPermissionsDefaultPermissionsParam) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationWorkloadIdentityProviderFederationRuleUpdateParamsPermissionsDefaultPermissionsParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationWorkloadIdentityProviderFederationRuleUpdateParamsPermissionsDefaultPermissionsParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Scopes, Type are required.
type AdminOrganizationWorkloadIdentityProviderFederationRuleUpdateParamsPermissionsRestrictedPermissionsParam struct {
	Scopes []string `json:"scopes,omitzero" api:"required"`
	// This field can be elided, and will marshal its zero value as "restricted".
	Type constant.Restricted `json:"type" default:"restricted"`
	paramObj
}

func (r AdminOrganizationWorkloadIdentityProviderFederationRuleUpdateParamsPermissionsRestrictedPermissionsParam) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationWorkloadIdentityProviderFederationRuleUpdateParamsPermissionsRestrictedPermissionsParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationWorkloadIdentityProviderFederationRuleUpdateParamsPermissionsRestrictedPermissionsParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationWorkloadIdentityProviderFederationRuleListParams struct {
	After            param.Opt[string]                                                            `query:"after,omitzero" json:"-"`
	Limit            param.Opt[int64]                                                             `query:"limit,omitzero" json:"-"`
	ServiceAccountID param.Opt[string]                                                            `query:"service_account_id,omitzero" json:"-"`
	Order            AdminOrganizationWorkloadIdentityProviderFederationRuleListParamsOrderString `query:"order,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes
// [AdminOrganizationWorkloadIdentityProviderFederationRuleListParams]'s query
// parameters as `url.Values`.
func (r AdminOrganizationWorkloadIdentityProviderFederationRuleListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type AdminOrganizationWorkloadIdentityProviderFederationRuleListParamsOrderString string

const (
	AdminOrganizationWorkloadIdentityProviderFederationRuleListParamsOrderStringAsc  AdminOrganizationWorkloadIdentityProviderFederationRuleListParamsOrderString = "asc"
	AdminOrganizationWorkloadIdentityProviderFederationRuleListParamsOrderStringDesc AdminOrganizationWorkloadIdentityProviderFederationRuleListParamsOrderString = "desc"
)

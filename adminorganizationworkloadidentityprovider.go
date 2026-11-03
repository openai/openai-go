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

// AdminOrganizationWorkloadIdentityProviderService contains methods and other
// services that help with interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAdminOrganizationWorkloadIdentityProviderService] method instead.
type AdminOrganizationWorkloadIdentityProviderService struct {
	Options         []option.RequestOption
	FederationRules AdminOrganizationWorkloadIdentityProviderFederationRuleService
}

// NewAdminOrganizationWorkloadIdentityProviderService generates a new service that
// applies the given options to each request. These options are applied after the
// parent client's options (if there is one), and before any request-specific
// options.
func NewAdminOrganizationWorkloadIdentityProviderService(opts ...option.RequestOption) (r AdminOrganizationWorkloadIdentityProviderService) {
	r = AdminOrganizationWorkloadIdentityProviderService{}
	r.Options = requestconfig.InheritedOptions(opts...)
	r.FederationRules = NewAdminOrganizationWorkloadIdentityProviderFederationRuleService(opts...)
	return
}

// Register an OIDC provider for keyless project service-account authentication.
func (r *AdminOrganizationWorkloadIdentityProviderService) New(ctx context.Context, body AdminOrganizationWorkloadIdentityProviderNewParams, opts ...option.RequestOption) (res *WorkloadIdentityProvider, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	path := "organization/workload_identity_providers"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Retrieve an OIDC provider in the authenticated organization.
func (r *AdminOrganizationWorkloadIdentityProviderService) Get(ctx context.Context, workloadIdentityProviderID string, opts ...option.RequestOption) (res *WorkloadIdentityProvider, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if workloadIdentityProviderID == "" {
		err = errors.New("missing required workload_identity_provider_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("organization/workload_identity_providers/%s", workloadIdentityProviderID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update supplied provider fields, preserving omitted fields.
func (r *AdminOrganizationWorkloadIdentityProviderService) Update(ctx context.Context, workloadIdentityProviderID string, body AdminOrganizationWorkloadIdentityProviderUpdateParams, opts ...option.RequestOption) (res *WorkloadIdentityProvider, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if workloadIdentityProviderID == "" {
		err = errors.New("missing required workload_identity_provider_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("organization/workload_identity_providers/%s", workloadIdentityProviderID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// List OIDC providers in the authenticated Admin API key's organization.
func (r *AdminOrganizationWorkloadIdentityProviderService) List(ctx context.Context, query AdminOrganizationWorkloadIdentityProviderListParams, opts ...option.RequestOption) (res *pagination.CursorPage[WorkloadIdentityProvider], err error) {
	var raw *http.Response
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "organization/workload_identity_providers"
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

// List OIDC providers in the authenticated Admin API key's organization.
func (r *AdminOrganizationWorkloadIdentityProviderService) ListAutoPaging(ctx context.Context, query AdminOrganizationWorkloadIdentityProviderListParams, opts ...option.RequestOption) *pagination.CursorPageAutoPager[WorkloadIdentityProvider] {
	return pagination.NewCursorPageAutoPager(r.List(ctx, query, opts...))
}

// Delete a provider after deleting all its federation rules.
func (r *AdminOrganizationWorkloadIdentityProviderService) Delete(ctx context.Context, workloadIdentityProviderID string, opts ...option.RequestOption) (res *WorkloadIdentityProviderDeleted, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if workloadIdentityProviderID == "" {
		err = errors.New("missing required workload_identity_provider_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("organization/workload_identity_providers/%s", workloadIdentityProviderID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

type WorkloadIdentityProvider struct {
	ID                       string                                            `json:"id" api:"required"`
	AttributeTransformations []WorkloadIdentityProviderAttributeTransformation `json:"attribute_transformations" api:"required"`
	Audience                 string                                            `json:"audience" api:"required"`
	CreatedAt                int64                                             `json:"created_at" api:"required"`
	Description              string                                            `json:"description" api:"required"`
	Enabled                  bool                                              `json:"enabled" api:"required"`
	Issuer                   string                                            `json:"issuer" api:"required"`
	Name                     string                                            `json:"name" api:"required"`
	Object                   constant.OrganizationWorkloadIdentityProvider     `json:"object" default:"organization.workload_identity_provider"`
	Type                     constant.Oidc                                     `json:"type" default:"oidc"`
	UpdatedAt                int64                                             `json:"updated_at" api:"required"`
	Verification             WorkloadIdentityProviderVerificationUnion         `json:"verification" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                       respjson.Field
		AttributeTransformations respjson.Field
		Audience                 respjson.Field
		CreatedAt                respjson.Field
		Description              respjson.Field
		Enabled                  respjson.Field
		Issuer                   respjson.Field
		Name                     respjson.Field
		Object                   respjson.Field
		Type                     respjson.Field
		UpdatedAt                respjson.Field
		Verification             respjson.Field
		ExtraFields              map[string]respjson.Field
		raw                      string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r WorkloadIdentityProvider) RawJSON() string { return r.JSON.raw }
func (r *WorkloadIdentityProvider) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type WorkloadIdentityProviderAttributeTransformation struct {
	Attribute  string `json:"attribute" api:"required"`
	Expression string `json:"expression" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Attribute   respjson.Field
		Expression  respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r WorkloadIdentityProviderAttributeTransformation) RawJSON() string { return r.JSON.raw }
func (r *WorkloadIdentityProviderAttributeTransformation) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// WorkloadIdentityProviderVerificationUnion contains all possible properties and
// values from [WorkloadIdentityProviderVerificationDiscovery],
// [WorkloadIdentityProviderVerificationUploadedJwks].
//
// Use the [WorkloadIdentityProviderVerificationUnion.AsAny] method to switch on
// the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type WorkloadIdentityProviderVerificationUnion struct {
	// This field is from variant [WorkloadIdentityProviderVerificationDiscovery].
	DiscoveryURL string `json:"discovery_url"`
	// Any of "discovery", "uploaded_jwks".
	Type string `json:"type"`
	// This field is from variant [WorkloadIdentityProviderVerificationUploadedJwks].
	Jwks map[string]any `json:"jwks"`
	JSON struct {
		DiscoveryURL respjson.Field
		Type         respjson.Field
		Jwks         respjson.Field
		raw          string
	} `json:"-"`
}

// anyWorkloadIdentityProviderVerification is implemented by each variant of
// [WorkloadIdentityProviderVerificationUnion] to add type safety for the return
// type of [WorkloadIdentityProviderVerificationUnion.AsAny]
type anyWorkloadIdentityProviderVerification interface {
	implWorkloadIdentityProviderVerificationUnion()
}

func (WorkloadIdentityProviderVerificationDiscovery) implWorkloadIdentityProviderVerificationUnion() {
}
func (WorkloadIdentityProviderVerificationUploadedJwks) implWorkloadIdentityProviderVerificationUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := WorkloadIdentityProviderVerificationUnion.AsAny().(type) {
//	case openai.WorkloadIdentityProviderVerificationDiscovery:
//	case openai.WorkloadIdentityProviderVerificationUploadedJwks:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u WorkloadIdentityProviderVerificationUnion) AsAny() anyWorkloadIdentityProviderVerification {
	switch u.Type {
	case "discovery":
		return u.AsDiscovery()
	case "uploaded_jwks":
		return u.AsUploadedJwks()
	}
	return nil
}

func (u WorkloadIdentityProviderVerificationUnion) AsDiscovery() (v WorkloadIdentityProviderVerificationDiscovery) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u WorkloadIdentityProviderVerificationUnion) AsUploadedJwks() (v WorkloadIdentityProviderVerificationUploadedJwks) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u WorkloadIdentityProviderVerificationUnion) RawJSON() string { return u.JSON.raw }

func (r *WorkloadIdentityProviderVerificationUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type WorkloadIdentityProviderVerificationDiscovery struct {
	DiscoveryURL string             `json:"discovery_url" api:"required"`
	Type         constant.Discovery `json:"type" default:"discovery"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		DiscoveryURL respjson.Field
		Type         respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r WorkloadIdentityProviderVerificationDiscovery) RawJSON() string { return r.JSON.raw }
func (r *WorkloadIdentityProviderVerificationDiscovery) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type WorkloadIdentityProviderVerificationUploadedJwks struct {
	Jwks map[string]any        `json:"jwks" api:"required"`
	Type constant.UploadedJwks `json:"type" default:"uploaded_jwks"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Jwks        respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r WorkloadIdentityProviderVerificationUploadedJwks) RawJSON() string { return r.JSON.raw }
func (r *WorkloadIdentityProviderVerificationUploadedJwks) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type WorkloadIdentityProviderDeleted struct {
	ID      string                                               `json:"id" api:"required"`
	Deleted bool                                                 `json:"deleted" api:"required"`
	Object  constant.OrganizationWorkloadIdentityProviderDeleted `json:"object" default:"organization.workload_identity_provider.deleted"`
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
func (r WorkloadIdentityProviderDeleted) RawJSON() string { return r.JSON.raw }
func (r *WorkloadIdentityProviderDeleted) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationWorkloadIdentityProviderNewParams struct {
	Audience                 string                                                                      `json:"audience" api:"required"`
	Issuer                   string                                                                      `json:"issuer" api:"required"`
	Name                     string                                                                      `json:"name" api:"required"`
	Description              param.Opt[string]                                                           `json:"description,omitzero"`
	Enabled                  param.Opt[bool]                                                             `json:"enabled,omitzero"`
	AttributeTransformations []AdminOrganizationWorkloadIdentityProviderNewParamsAttributeTransformation `json:"attribute_transformations,omitzero"`
	Verification             AdminOrganizationWorkloadIdentityProviderNewParamsVerificationUnion         `json:"verification,omitzero"`
	// This field can be elided, and will marshal its zero value as "oidc".
	Type constant.Oidc `json:"type" default:"oidc"`
	paramObj
}

func (r AdminOrganizationWorkloadIdentityProviderNewParams) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationWorkloadIdentityProviderNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationWorkloadIdentityProviderNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Attribute, Expression are required.
type AdminOrganizationWorkloadIdentityProviderNewParamsAttributeTransformation struct {
	Attribute  string `json:"attribute" api:"required"`
	Expression string `json:"expression" api:"required"`
	paramObj
}

func (r AdminOrganizationWorkloadIdentityProviderNewParamsAttributeTransformation) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationWorkloadIdentityProviderNewParamsAttributeTransformation
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationWorkloadIdentityProviderNewParamsAttributeTransformation) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type AdminOrganizationWorkloadIdentityProviderNewParamsVerificationUnion struct {
	OfDiscovery    *AdminOrganizationWorkloadIdentityProviderNewParamsVerificationDiscovery    `json:",omitzero,inline"`
	OfUploadedJwks *AdminOrganizationWorkloadIdentityProviderNewParamsVerificationUploadedJwks `json:",omitzero,inline"`
	paramUnion
}

func (u AdminOrganizationWorkloadIdentityProviderNewParamsVerificationUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfDiscovery, u.OfUploadedJwks)
}
func (u *AdminOrganizationWorkloadIdentityProviderNewParamsVerificationUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u AdminOrganizationWorkloadIdentityProviderNewParamsVerificationUnion) GetDiscoveryURL() *string {
	if vt := u.OfDiscovery; vt != nil && vt.DiscoveryURL.Valid() {
		return &vt.DiscoveryURL.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AdminOrganizationWorkloadIdentityProviderNewParamsVerificationUnion) GetJwks() map[string]any {
	if vt := u.OfUploadedJwks; vt != nil {
		return vt.Jwks
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AdminOrganizationWorkloadIdentityProviderNewParamsVerificationUnion) GetType() *string {
	if vt := u.OfDiscovery; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfUploadedJwks; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[AdminOrganizationWorkloadIdentityProviderNewParamsVerificationUnion](
		"type",
		apijson.Discriminator[AdminOrganizationWorkloadIdentityProviderNewParamsVerificationDiscovery]("discovery"),
		apijson.Discriminator[AdminOrganizationWorkloadIdentityProviderNewParamsVerificationUploadedJwks]("uploaded_jwks"),
	)
}

// The properties DiscoveryURL, Type are required.
type AdminOrganizationWorkloadIdentityProviderNewParamsVerificationDiscovery struct {
	DiscoveryURL param.Opt[string] `json:"discovery_url,omitzero" api:"required"`
	// This field can be elided, and will marshal its zero value as "discovery".
	Type constant.Discovery `json:"type" default:"discovery"`
	paramObj
}

func (r AdminOrganizationWorkloadIdentityProviderNewParamsVerificationDiscovery) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationWorkloadIdentityProviderNewParamsVerificationDiscovery
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationWorkloadIdentityProviderNewParamsVerificationDiscovery) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Jwks, Type are required.
type AdminOrganizationWorkloadIdentityProviderNewParamsVerificationUploadedJwks struct {
	Jwks map[string]any `json:"jwks,omitzero" api:"required"`
	// This field can be elided, and will marshal its zero value as "uploaded_jwks".
	Type constant.UploadedJwks `json:"type" default:"uploaded_jwks"`
	paramObj
}

func (r AdminOrganizationWorkloadIdentityProviderNewParamsVerificationUploadedJwks) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationWorkloadIdentityProviderNewParamsVerificationUploadedJwks
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationWorkloadIdentityProviderNewParamsVerificationUploadedJwks) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationWorkloadIdentityProviderUpdateParams struct {
	Description              param.Opt[string]                                                              `json:"description,omitzero"`
	Audience                 param.Opt[string]                                                              `json:"audience,omitzero"`
	Enabled                  param.Opt[bool]                                                                `json:"enabled,omitzero"`
	Issuer                   param.Opt[string]                                                              `json:"issuer,omitzero"`
	Name                     param.Opt[string]                                                              `json:"name,omitzero"`
	AttributeTransformations []AdminOrganizationWorkloadIdentityProviderUpdateParamsAttributeTransformation `json:"attribute_transformations,omitzero"`
	Verification             AdminOrganizationWorkloadIdentityProviderUpdateParamsVerificationUnion         `json:"verification,omitzero"`
	paramObj
}

func (r AdminOrganizationWorkloadIdentityProviderUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationWorkloadIdentityProviderUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationWorkloadIdentityProviderUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Attribute, Expression are required.
type AdminOrganizationWorkloadIdentityProviderUpdateParamsAttributeTransformation struct {
	Attribute  string `json:"attribute" api:"required"`
	Expression string `json:"expression" api:"required"`
	paramObj
}

func (r AdminOrganizationWorkloadIdentityProviderUpdateParamsAttributeTransformation) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationWorkloadIdentityProviderUpdateParamsAttributeTransformation
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationWorkloadIdentityProviderUpdateParamsAttributeTransformation) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type AdminOrganizationWorkloadIdentityProviderUpdateParamsVerificationUnion struct {
	OfAdminOrganizationWorkloadIdentityProviderUpdatesVerificationDiscoveryVerificationParam    *AdminOrganizationWorkloadIdentityProviderUpdateParamsVerificationDiscoveryVerificationParam    `json:",omitzero,inline"`
	OfAdminOrganizationWorkloadIdentityProviderUpdatesVerificationUploadedJwksVerificationParam *AdminOrganizationWorkloadIdentityProviderUpdateParamsVerificationUploadedJwksVerificationParam `json:",omitzero,inline"`
	paramUnion
}

func (u AdminOrganizationWorkloadIdentityProviderUpdateParamsVerificationUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfAdminOrganizationWorkloadIdentityProviderUpdatesVerificationDiscoveryVerificationParam, u.OfAdminOrganizationWorkloadIdentityProviderUpdatesVerificationUploadedJwksVerificationParam)
}
func (u *AdminOrganizationWorkloadIdentityProviderUpdateParamsVerificationUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u AdminOrganizationWorkloadIdentityProviderUpdateParamsVerificationUnion) GetDiscoveryURL() *string {
	if vt := u.OfAdminOrganizationWorkloadIdentityProviderUpdatesVerificationDiscoveryVerificationParam; vt != nil && vt.DiscoveryURL.Valid() {
		return &vt.DiscoveryURL.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AdminOrganizationWorkloadIdentityProviderUpdateParamsVerificationUnion) GetJwks() map[string]any {
	if vt := u.OfAdminOrganizationWorkloadIdentityProviderUpdatesVerificationUploadedJwksVerificationParam; vt != nil {
		return vt.Jwks
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AdminOrganizationWorkloadIdentityProviderUpdateParamsVerificationUnion) GetType() *string {
	if vt := u.OfAdminOrganizationWorkloadIdentityProviderUpdatesVerificationDiscoveryVerificationParam; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfAdminOrganizationWorkloadIdentityProviderUpdatesVerificationUploadedJwksVerificationParam; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[AdminOrganizationWorkloadIdentityProviderUpdateParamsVerificationUnion](
		"type",
		apijson.Discriminator[AdminOrganizationWorkloadIdentityProviderUpdateParamsVerificationDiscoveryVerificationParam]("discovery"),
		apijson.Discriminator[AdminOrganizationWorkloadIdentityProviderUpdateParamsVerificationUploadedJwksVerificationParam]("uploaded_jwks"),
	)
}

// The properties DiscoveryURL, Type are required.
type AdminOrganizationWorkloadIdentityProviderUpdateParamsVerificationDiscoveryVerificationParam struct {
	DiscoveryURL param.Opt[string] `json:"discovery_url,omitzero" api:"required"`
	// This field can be elided, and will marshal its zero value as "discovery".
	Type constant.Discovery `json:"type" default:"discovery"`
	paramObj
}

func (r AdminOrganizationWorkloadIdentityProviderUpdateParamsVerificationDiscoveryVerificationParam) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationWorkloadIdentityProviderUpdateParamsVerificationDiscoveryVerificationParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationWorkloadIdentityProviderUpdateParamsVerificationDiscoveryVerificationParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Jwks, Type are required.
type AdminOrganizationWorkloadIdentityProviderUpdateParamsVerificationUploadedJwksVerificationParam struct {
	Jwks map[string]any `json:"jwks,omitzero" api:"required"`
	// This field can be elided, and will marshal its zero value as "uploaded_jwks".
	Type constant.UploadedJwks `json:"type" default:"uploaded_jwks"`
	paramObj
}

func (r AdminOrganizationWorkloadIdentityProviderUpdateParamsVerificationUploadedJwksVerificationParam) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationWorkloadIdentityProviderUpdateParamsVerificationUploadedJwksVerificationParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationWorkloadIdentityProviderUpdateParamsVerificationUploadedJwksVerificationParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationWorkloadIdentityProviderListParams struct {
	After param.Opt[string]                                              `query:"after,omitzero" json:"-"`
	Limit param.Opt[int64]                                               `query:"limit,omitzero" json:"-"`
	Order AdminOrganizationWorkloadIdentityProviderListParamsOrderString `query:"order,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [AdminOrganizationWorkloadIdentityProviderListParams]'s
// query parameters as `url.Values`.
func (r AdminOrganizationWorkloadIdentityProviderListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type AdminOrganizationWorkloadIdentityProviderListParamsOrderString string

const (
	AdminOrganizationWorkloadIdentityProviderListParamsOrderStringAsc  AdminOrganizationWorkloadIdentityProviderListParamsOrderString = "asc"
	AdminOrganizationWorkloadIdentityProviderListParamsOrderStringDesc AdminOrganizationWorkloadIdentityProviderListParamsOrderString = "desc"
)

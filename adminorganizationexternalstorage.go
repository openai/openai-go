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

// AdminOrganizationExternalStorageService contains methods and other services that
// help with interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAdminOrganizationExternalStorageService] method instead.
type AdminOrganizationExternalStorageService struct {
	Options []option.RequestOption
}

// NewAdminOrganizationExternalStorageService generates a new service that applies
// the given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewAdminOrganizationExternalStorageService(opts ...option.RequestOption) (r AdminOrganizationExternalStorageService) {
	r = AdminOrganizationExternalStorageService{}
	r.Options = requestconfig.InheritedOptions(opts...)
	return
}

// Register one customer-managed external storage configuration.
func (r *AdminOrganizationExternalStorageService) New(ctx context.Context, body AdminOrganizationExternalStorageNewParams, opts ...option.RequestOption) (res *ExternalStorageConfiguration, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	path := "organization/external_storage"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Get one customer-managed external storage configuration.
func (r *AdminOrganizationExternalStorageService) Get(ctx context.Context, externalStorageID string, opts ...option.RequestOption) (res *ExternalStorageConfiguration, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if externalStorageID == "" {
		err = errors.New("missing required external_storage_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("organization/external_storage/%s", externalStorageID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// List the organization's customer-managed external storage configurations.
func (r *AdminOrganizationExternalStorageService) List(ctx context.Context, query AdminOrganizationExternalStorageListParams, opts ...option.RequestOption) (res *pagination.CursorPage[ExternalStorageConfiguration], err error) {
	var raw *http.Response
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "organization/external_storage"
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

// List the organization's customer-managed external storage configurations.
func (r *AdminOrganizationExternalStorageService) ListAutoPaging(ctx context.Context, query AdminOrganizationExternalStorageListParams, opts ...option.RequestOption) *pagination.CursorPageAutoPager[ExternalStorageConfiguration] {
	return pagination.NewCursorPageAutoPager(r.List(ctx, query, opts...))
}

// Soft-delete one customer-managed external storage configuration.
func (r *AdminOrganizationExternalStorageService) Delete(ctx context.Context, externalStorageID string, opts ...option.RequestOption) (res *ExternalStorageDeleted, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if externalStorageID == "" {
		err = errors.New("missing required external_storage_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("organization/external_storage/%s", externalStorageID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Validate one customer-managed external storage configuration.
func (r *AdminOrganizationExternalStorageService) Validate(ctx context.Context, externalStorageID string, opts ...option.RequestOption) (res *ExternalStorageConfiguration, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if externalStorageID == "" {
		err = errors.New("missing required external_storage_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("organization/external_storage/%s/validate", externalStorageID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

type AwsExternalStorageProvider struct {
	AccountID  string       `json:"account_id" api:"required"`
	Bucket     string       `json:"bucket" api:"required"`
	ExternalID string       `json:"external_id" api:"required"`
	Region     string       `json:"region" api:"required"`
	RoleArn    string       `json:"role_arn" api:"required"`
	Type       constant.Aws `json:"type" default:"aws"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		AccountID   respjson.Field
		Bucket      respjson.Field
		ExternalID  respjson.Field
		Region      respjson.Field
		RoleArn     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AwsExternalStorageProvider) RawJSON() string { return r.JSON.raw }
func (r *AwsExternalStorageProvider) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AzureExternalStorageProvider struct {
	AccountName    string         `json:"account_name" api:"required"`
	Container      string         `json:"container" api:"required"`
	Region         string         `json:"region" api:"required"`
	ResourceGroup  string         `json:"resource_group" api:"required"`
	SubscriptionID string         `json:"subscription_id" api:"required"`
	TenantID       string         `json:"tenant_id" api:"required"`
	Type           constant.Azure `json:"type" default:"azure"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		AccountName    respjson.Field
		Container      respjson.Field
		Region         respjson.Field
		ResourceGroup  respjson.Field
		SubscriptionID respjson.Field
		TenantID       respjson.Field
		Type           respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AzureExternalStorageProvider) RawJSON() string { return r.JSON.raw }
func (r *AzureExternalStorageProvider) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ExternalStorageConfiguration struct {
	ID        string                                    `json:"id" api:"required"`
	CreatedAt int64                                     `json:"created_at" api:"required" format:"unixtime"`
	Geography string                                    `json:"geography" api:"required"`
	Object    constant.OrganizationExternalStorage      `json:"object" default:"organization.external_storage"`
	ProjectID string                                    `json:"project_id" api:"required"`
	Provider  ExternalStorageConfigurationProviderUnion `json:"provider" api:"required"`
	// Any of "pending", "validated", "unhealthy".
	Status ExternalStorageConfigurationStatus `json:"status" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		Geography   respjson.Field
		Object      respjson.Field
		ProjectID   respjson.Field
		Provider    respjson.Field
		Status      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ExternalStorageConfiguration) RawJSON() string { return r.JSON.raw }
func (r *ExternalStorageConfiguration) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ExternalStorageConfigurationProviderUnion contains all possible properties and
// values from [AwsExternalStorageProvider], [AzureExternalStorageProvider].
//
// Use the [ExternalStorageConfigurationProviderUnion.AsAny] method to switch on
// the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ExternalStorageConfigurationProviderUnion struct {
	// This field is from variant [AwsExternalStorageProvider].
	AccountID string `json:"account_id"`
	// This field is from variant [AwsExternalStorageProvider].
	Bucket string `json:"bucket"`
	// This field is from variant [AwsExternalStorageProvider].
	ExternalID string `json:"external_id"`
	Region     string `json:"region"`
	// This field is from variant [AwsExternalStorageProvider].
	RoleArn string `json:"role_arn"`
	// Any of "aws", "azure".
	Type string `json:"type"`
	// This field is from variant [AzureExternalStorageProvider].
	AccountName string `json:"account_name"`
	// This field is from variant [AzureExternalStorageProvider].
	Container string `json:"container"`
	// This field is from variant [AzureExternalStorageProvider].
	ResourceGroup string `json:"resource_group"`
	// This field is from variant [AzureExternalStorageProvider].
	SubscriptionID string `json:"subscription_id"`
	// This field is from variant [AzureExternalStorageProvider].
	TenantID string `json:"tenant_id"`
	JSON     struct {
		AccountID      respjson.Field
		Bucket         respjson.Field
		ExternalID     respjson.Field
		Region         respjson.Field
		RoleArn        respjson.Field
		Type           respjson.Field
		AccountName    respjson.Field
		Container      respjson.Field
		ResourceGroup  respjson.Field
		SubscriptionID respjson.Field
		TenantID       respjson.Field
		raw            string
	} `json:"-"`
}

// anyExternalStorageConfigurationProvider is implemented by each variant of
// [ExternalStorageConfigurationProviderUnion] to add type safety for the return
// type of [ExternalStorageConfigurationProviderUnion.AsAny]
type anyExternalStorageConfigurationProvider interface {
	implExternalStorageConfigurationProviderUnion()
}

func (AwsExternalStorageProvider) implExternalStorageConfigurationProviderUnion()   {}
func (AzureExternalStorageProvider) implExternalStorageConfigurationProviderUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := ExternalStorageConfigurationProviderUnion.AsAny().(type) {
//	case openai.AwsExternalStorageProvider:
//	case openai.AzureExternalStorageProvider:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ExternalStorageConfigurationProviderUnion) AsAny() anyExternalStorageConfigurationProvider {
	switch u.Type {
	case "aws":
		return u.AsAws()
	case "azure":
		return u.AsAzure()
	}
	return nil
}

func (u ExternalStorageConfigurationProviderUnion) AsAws() (v AwsExternalStorageProvider) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ExternalStorageConfigurationProviderUnion) AsAzure() (v AzureExternalStorageProvider) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ExternalStorageConfigurationProviderUnion) RawJSON() string { return u.JSON.raw }

func (r *ExternalStorageConfigurationProviderUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ExternalStorageConfigurationStatus string

const (
	ExternalStorageConfigurationStatusPending   ExternalStorageConfigurationStatus = "pending"
	ExternalStorageConfigurationStatusValidated ExternalStorageConfigurationStatus = "validated"
	ExternalStorageConfigurationStatusUnhealthy ExternalStorageConfigurationStatus = "unhealthy"
)

type ExternalStorageDeleted struct {
	ID      string                                      `json:"id" api:"required"`
	Deleted bool                                        `json:"deleted" api:"required"`
	Object  constant.OrganizationExternalStorageDeleted `json:"object" default:"organization.external_storage.deleted"`
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
func (r ExternalStorageDeleted) RawJSON() string { return r.JSON.raw }
func (r *ExternalStorageDeleted) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationExternalStorageNewParams struct {
	ProjectID string                                                 `json:"project_id" api:"required"`
	Provider  AdminOrganizationExternalStorageNewParamsProviderUnion `json:"provider,omitzero" api:"required"`
	paramObj
}

func (r AdminOrganizationExternalStorageNewParams) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationExternalStorageNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationExternalStorageNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type AdminOrganizationExternalStorageNewParamsProviderUnion struct {
	OfAws   *AdminOrganizationExternalStorageNewParamsProviderAws   `json:",omitzero,inline"`
	OfAzure *AdminOrganizationExternalStorageNewParamsProviderAzure `json:",omitzero,inline"`
	paramUnion
}

func (u AdminOrganizationExternalStorageNewParamsProviderUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfAws, u.OfAzure)
}
func (u *AdminOrganizationExternalStorageNewParamsProviderUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u AdminOrganizationExternalStorageNewParamsProviderUnion) GetBucket() *string {
	if vt := u.OfAws; vt != nil {
		return &vt.Bucket
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AdminOrganizationExternalStorageNewParamsProviderUnion) GetRoleArn() *string {
	if vt := u.OfAws; vt != nil {
		return &vt.RoleArn
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AdminOrganizationExternalStorageNewParamsProviderUnion) GetAccountName() *string {
	if vt := u.OfAzure; vt != nil {
		return &vt.AccountName
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AdminOrganizationExternalStorageNewParamsProviderUnion) GetContainer() *string {
	if vt := u.OfAzure; vt != nil {
		return &vt.Container
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AdminOrganizationExternalStorageNewParamsProviderUnion) GetResourceGroup() *string {
	if vt := u.OfAzure; vt != nil {
		return &vt.ResourceGroup
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AdminOrganizationExternalStorageNewParamsProviderUnion) GetSubscriptionID() *string {
	if vt := u.OfAzure; vt != nil {
		return &vt.SubscriptionID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AdminOrganizationExternalStorageNewParamsProviderUnion) GetTenantID() *string {
	if vt := u.OfAzure; vt != nil {
		return &vt.TenantID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AdminOrganizationExternalStorageNewParamsProviderUnion) GetType() *string {
	if vt := u.OfAws; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfAzure; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[AdminOrganizationExternalStorageNewParamsProviderUnion](
		"type",
		apijson.Discriminator[AdminOrganizationExternalStorageNewParamsProviderAws]("aws"),
		apijson.Discriminator[AdminOrganizationExternalStorageNewParamsProviderAzure]("azure"),
	)
}

// The properties Bucket, RoleArn, Type are required.
type AdminOrganizationExternalStorageNewParamsProviderAws struct {
	Bucket  string `json:"bucket" api:"required"`
	RoleArn string `json:"role_arn" api:"required"`
	// This field can be elided, and will marshal its zero value as "aws".
	Type constant.Aws `json:"type" default:"aws"`
	paramObj
}

func (r AdminOrganizationExternalStorageNewParamsProviderAws) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationExternalStorageNewParamsProviderAws
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationExternalStorageNewParamsProviderAws) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties AccountName, Container, ResourceGroup, SubscriptionID, TenantID,
// Type are required.
type AdminOrganizationExternalStorageNewParamsProviderAzure struct {
	AccountName    string `json:"account_name" api:"required"`
	Container      string `json:"container" api:"required"`
	ResourceGroup  string `json:"resource_group" api:"required"`
	SubscriptionID string `json:"subscription_id" api:"required"`
	TenantID       string `json:"tenant_id" api:"required"`
	// This field can be elided, and will marshal its zero value as "azure".
	Type constant.Azure `json:"type" default:"azure"`
	paramObj
}

func (r AdminOrganizationExternalStorageNewParamsProviderAzure) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationExternalStorageNewParamsProviderAzure
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationExternalStorageNewParamsProviderAzure) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationExternalStorageListParams struct {
	// Return external storage configurations after this ID.
	After     param.Opt[string] `query:"after,omitzero" json:"-"`
	Limit     param.Opt[int64]  `query:"limit,omitzero" json:"-"`
	ProjectID param.Opt[string] `query:"project_id,omitzero" json:"-"`
	// Any of "asc", "desc".
	Order AdminOrganizationExternalStorageListParamsOrder `query:"order,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [AdminOrganizationExternalStorageListParams]'s query
// parameters as `url.Values`.
func (r AdminOrganizationExternalStorageListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type AdminOrganizationExternalStorageListParamsOrder string

const (
	AdminOrganizationExternalStorageListParamsOrderAsc  AdminOrganizationExternalStorageListParamsOrder = "asc"
	AdminOrganizationExternalStorageListParamsOrderDesc AdminOrganizationExternalStorageListParamsOrder = "desc"
)

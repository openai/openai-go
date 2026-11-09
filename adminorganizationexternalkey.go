// File generated from our OpenAPI spec by Castiron. See CONTRIBUTING.md for details.

package openai

import (
	"context"
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

// AdminOrganizationExternalKeyService contains methods and other services that
// help with interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAdminOrganizationExternalKeyService] method instead.
type AdminOrganizationExternalKeyService struct {
	Options []option.RequestOption
}

// NewAdminOrganizationExternalKeyService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewAdminOrganizationExternalKeyService(opts ...option.RequestOption) (r AdminOrganizationExternalKeyService) {
	r = AdminOrganizationExternalKeyService{}
	r.Options = requestconfig.InheritedOptions(opts...)
	return
}

// Add an external key config for an organization
func (r *AdminOrganizationExternalKeyService) New(ctx context.Context, body AdminOrganizationExternalKeyNewParams, opts ...option.RequestOption) (res *ExternalKey, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	path := "organization/external_keys"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Retrieve a stored external key registration for the user's organization.
func (r *AdminOrganizationExternalKeyService) Get(ctx context.Context, externalKeyID string, opts ...option.RequestOption) (res *ExternalKey, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if externalKeyID == "" {
		err = errors.New("missing required external_key_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("organization/external_keys/%s", externalKeyID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// List the external key configs registered on the user's organization
func (r *AdminOrganizationExternalKeyService) List(ctx context.Context, query AdminOrganizationExternalKeyListParams, opts ...option.RequestOption) (res *pagination.CursorPage[ExternalKey], err error) {
	var raw *http.Response
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "organization/external_keys"
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

// List the external key configs registered on the user's organization
func (r *AdminOrganizationExternalKeyService) ListAutoPaging(ctx context.Context, query AdminOrganizationExternalKeyListParams, opts ...option.RequestOption) *pagination.CursorPageAutoPager[ExternalKey] {
	return pagination.NewCursorPageAutoPager(r.List(ctx, query, opts...))
}

// Delete an external key config for the user's organization
func (r *AdminOrganizationExternalKeyService) Delete(ctx context.Context, externalKeyID string, opts ...option.RequestOption) (res *ExternalKey, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if externalKeyID == "" {
		err = errors.New("missing required external_key_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("organization/external_keys/%s", externalKeyID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Validate an external key config for the user's organization
func (r *AdminOrganizationExternalKeyService) Validate(ctx context.Context, externalKeyID string, opts ...option.RequestOption) (res *AdminOrganizationExternalKeyValidateResponse, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if externalKeyID == "" {
		err = errors.New("missing required external_key_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("organization/external_keys/%s/validate", externalKeyID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

type ExternalKey struct {
	ID             string   `json:"id" api:"required"`
	APIProjectIDs  []string `json:"api_project_ids" api:"required"`
	CreatedAt      int64    `json:"created_at" api:"required"`
	Name           string   `json:"name" api:"required"`
	Object         string   `json:"object" api:"required"`
	Type           string   `json:"type" api:"required"`
	WorkspaceID    string   `json:"workspace_id" api:"required"`
	Audience       string   `json:"audience"`
	ExternalID     string   `json:"external_id"`
	KeyName        string   `json:"key_name"`
	KeyVersion     string   `json:"key_version"`
	KmsArn         string   `json:"kms_arn"`
	KmsKeyLocation string   `json:"kms_key_location"`
	KmsKeyName     string   `json:"kms_key_name"`
	KmsKeyRingName string   `json:"kms_key_ring_name"`
	KmsKeyVersion  string   `json:"kms_key_version"`
	KmsProjectID   string   `json:"kms_project_id"`
	// Any of "ekm", "psp_decryption".
	Purpose                       ExternalKeyPurpose `json:"purpose"`
	RoleArn                       string             `json:"role_arn"`
	TenantID                      string             `json:"tenant_id"`
	VaultUri                      string             `json:"vault_uri"`
	WorkloadIdentityPoolID        string             `json:"workload_identity_pool_id"`
	WorkloadIdentityProjectNumber string             `json:"workload_identity_project_number"`
	WorkloadIdentityProviderID    string             `json:"workload_identity_provider_id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                            respjson.Field
		APIProjectIDs                 respjson.Field
		CreatedAt                     respjson.Field
		Name                          respjson.Field
		Object                        respjson.Field
		Type                          respjson.Field
		WorkspaceID                   respjson.Field
		Audience                      respjson.Field
		ExternalID                    respjson.Field
		KeyName                       respjson.Field
		KeyVersion                    respjson.Field
		KmsArn                        respjson.Field
		KmsKeyLocation                respjson.Field
		KmsKeyName                    respjson.Field
		KmsKeyRingName                respjson.Field
		KmsKeyVersion                 respjson.Field
		KmsProjectID                  respjson.Field
		Purpose                       respjson.Field
		RoleArn                       respjson.Field
		TenantID                      respjson.Field
		VaultUri                      respjson.Field
		WorkloadIdentityPoolID        respjson.Field
		WorkloadIdentityProjectNumber respjson.Field
		WorkloadIdentityProviderID    respjson.Field
		ExtraFields                   map[string]respjson.Field
		raw                           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ExternalKey) RawJSON() string { return r.JSON.raw }
func (r *ExternalKey) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ExternalKeyPurpose string

const (
	ExternalKeyPurposeEkm           ExternalKeyPurpose = "ekm"
	ExternalKeyPurposePspDecryption ExternalKeyPurpose = "psp_decryption"
)

type AdminOrganizationExternalKeyValidateResponse struct {
	Status string `json:"status" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Status      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationExternalKeyValidateResponse) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationExternalKeyValidateResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationExternalKeyNewParams struct {

	//
	// Request body variants
	//

	// This field is a request body variant, only one variant field can be set.
	OfAws *AdminOrganizationExternalKeyNewParamsBodyAws `json:",inline"`
	// This field is a request body variant, only one variant field can be set.
	OfGcp *AdminOrganizationExternalKeyNewParamsBodyGcp `json:",inline"`
	// This field is a request body variant, only one variant field can be set.
	OfAzure *AdminOrganizationExternalKeyNewParamsBodyAzure `json:",inline"`

	paramUnion
}

func (u AdminOrganizationExternalKeyNewParams) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfAws, u.OfGcp, u.OfAzure)
}
func (r *AdminOrganizationExternalKeyNewParams) UnmarshalJSON(data []byte) error {

	// Decode afresh so reusing a params value cannot leave another arm populated.
	var decoded AdminOrganizationExternalKeyNewParams
	if err := apijson.UnmarshalRoot(data, &decoded); err != nil {
		return err
	}
	*r = decoded
	return nil
}

func init() {
	apijson.RegisterUnion[AdminOrganizationExternalKeyNewParams]("type",
		apijson.Discriminator[AdminOrganizationExternalKeyNewParamsBodyAws]("aws"),
		apijson.Discriminator[AdminOrganizationExternalKeyNewParamsBodyGcp]("gcp"),
		apijson.Discriminator[AdminOrganizationExternalKeyNewParamsBodyAzure]("azure"),
	)
}

// The properties ExternalID, KmsArn, Name, RoleArn, Type are required.
type AdminOrganizationExternalKeyNewParamsBodyAws struct {
	// The ExternalID that OpenAI will pass in the AssumeRole call for the customer's
	// AWS account. This must be the OpenAI organization ID or an API project ID that
	// belongs to the organization.
	ExternalID string `json:"external_id" api:"required"`
	// The Key Management System ARN for the master key you manage
	KmsArn string `json:"kms_arn" api:"required"`
	// A friendly name to identify this config
	Name string `json:"name" api:"required"`
	// The Role ARN that OpenAI will assume to call your KMS
	RoleArn string `json:"role_arn" api:"required"`
	// Any of "ekm", "psp_decryption".
	Purpose param.Opt[string] `json:"purpose,omitzero"`
	// This field can be elided, and will marshal its zero value as "aws".
	Type constant.Aws `json:"type" default:"aws"`
	paramObj
}

func (r AdminOrganizationExternalKeyNewParamsBodyAws) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationExternalKeyNewParamsBodyAws
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationExternalKeyNewParamsBodyAws) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Audience, KmsKeyLocation, KmsKeyName, KmsKeyRingName,
// KmsProjectID, Name, WorkloadIdentityPoolID, WorkloadIdentityProjectNumber,
// WorkloadIdentityProviderID, Type are required.
type AdminOrganizationExternalKeyNewParamsBodyGcp struct {
	// The aud field in the identity token OpenAI will pass in the call to the
	// customer's GCP STS. This must be the OpenAI organization ID or an API project ID
	// that belongs to the organization.
	Audience string `json:"audience" api:"required"`
	// The region where your Key Management System master key is located
	KmsKeyLocation string `json:"kms_key_location" api:"required"`
	// The name of the Key Management System master key you manage
	KmsKeyName string `json:"kms_key_name" api:"required"`
	// The Key Management System key ring for the master key you manage
	KmsKeyRingName string `json:"kms_key_ring_name" api:"required"`
	// The string GCP project name where your KMS resides
	KmsProjectID string `json:"kms_project_id" api:"required"`
	// A friendly name to identify this config
	Name string `json:"name" api:"required"`
	// The pool containing the Workload Identity provider that you registered for
	// OpenAI
	WorkloadIdentityPoolID string `json:"workload_identity_pool_id" api:"required"`
	// The numeric GCP project number where you registered OpenAI's workload identity
	WorkloadIdentityProjectNumber string `json:"workload_identity_project_number" api:"required"`
	// The Workload Identity provider that you registered for OpenAI
	WorkloadIdentityProviderID string `json:"workload_identity_provider_id" api:"required"`
	// An exact CryptoKeyVersion number. Required for PSP decryption; not accepted for
	// EKM.
	KmsKeyVersion param.Opt[string] `json:"kms_key_version,omitzero"`
	// Any of "ekm", "psp_decryption".
	Purpose param.Opt[string] `json:"purpose,omitzero"`
	// This field can be elided, and will marshal its zero value as "gcp".
	Type constant.Gcp `json:"type" default:"gcp"`
	paramObj
}

func (r AdminOrganizationExternalKeyNewParamsBodyGcp) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationExternalKeyNewParamsBodyGcp
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationExternalKeyNewParamsBodyGcp) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties KeyName, Name, TenantID, VaultURI, Type are required.
type AdminOrganizationExternalKeyNewParamsBodyAzure struct {
	// The name of the Azure Key Vault master key that you manage
	KeyName string `json:"key_name" api:"required"`
	// A friendly name to identify this config
	Name string `json:"name" api:"required"`
	// Your Azure tenant UUID
	TenantID string `json:"tenant_id" api:"required"`
	// The URI of the Azure vault containing the master key you manage
	VaultURI string `json:"vault_uri" api:"required"`
	// An exact Azure Key Vault key version. Required for PSP decryption; not accepted
	// for EKM.
	KeyVersion param.Opt[string] `json:"key_version,omitzero"`
	// Any of "ekm", "psp_decryption".
	Purpose param.Opt[string] `json:"purpose,omitzero"`
	// This field can be elided, and will marshal its zero value as "azure".
	Type constant.Azure `json:"type" default:"azure"`
	paramObj
}

func (r AdminOrganizationExternalKeyNewParamsBodyAzure) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationExternalKeyNewParamsBodyAzure
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationExternalKeyNewParamsBodyAzure) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationExternalKeyListParams struct {
	After param.Opt[string] `query:"after,omitzero" json:"-"`
	Limit param.Opt[int64]  `query:"limit,omitzero" json:"-"`
	// Any of "asc", "desc".
	Order AdminOrganizationExternalKeyListParamsOrder `query:"order,omitzero" json:"-"`
	// Filter by key purpose. When omitted, returns keys of all purposes.
	//
	// Any of "ekm", "psp_decryption".
	Purpose AdminOrganizationExternalKeyListParamsPurpose `query:"purpose,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [AdminOrganizationExternalKeyListParams]'s query parameters
// as `url.Values`.
func (r AdminOrganizationExternalKeyListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type AdminOrganizationExternalKeyListParamsOrder string

const (
	AdminOrganizationExternalKeyListParamsOrderAsc  AdminOrganizationExternalKeyListParamsOrder = "asc"
	AdminOrganizationExternalKeyListParamsOrderDesc AdminOrganizationExternalKeyListParamsOrder = "desc"
)

// Filter by key purpose. When omitted, returns keys of all purposes.
type AdminOrganizationExternalKeyListParamsPurpose string

const (
	AdminOrganizationExternalKeyListParamsPurposeEkm           AdminOrganizationExternalKeyListParamsPurpose = "ekm"
	AdminOrganizationExternalKeyListParamsPurposePspDecryption AdminOrganizationExternalKeyListParamsPurpose = "psp_decryption"
)

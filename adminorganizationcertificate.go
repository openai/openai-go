// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai

import (
	"context"
	"errors"
	"fmt"
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

// AdminOrganizationCertificateService contains methods and other services that
// help with interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAdminOrganizationCertificateService] method instead.
type AdminOrganizationCertificateService struct {
	Options []option.RequestOption
}

// NewAdminOrganizationCertificateService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewAdminOrganizationCertificateService(opts ...option.RequestOption) (r AdminOrganizationCertificateService) {
	r = AdminOrganizationCertificateService{}
	r.Options = opts
	return
}

// Upload a certificate to the organization. This does **not** automatically
// activate the certificate.
//
// Organizations can upload up to 50 certificates.
func (r *AdminOrganizationCertificateService) New(ctx context.Context, body AdminOrganizationCertificateNewParams, opts ...option.RequestOption) (res *Certificate, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	path := "organization/certificates"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Get a certificate that has been uploaded to the organization.
//
// You can get a certificate regardless of whether it is active or not.
func (r *AdminOrganizationCertificateService) Get(ctx context.Context, certificateID string, query AdminOrganizationCertificateGetParams, opts ...option.RequestOption) (res *Certificate, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if certificateID == "" {
		err = errors.New("missing required certificate_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("organization/certificates/%s", certificateID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

// Modify a certificate. Note that only the name can be modified.
func (r *AdminOrganizationCertificateService) Update(ctx context.Context, certificateID string, body AdminOrganizationCertificateUpdateParams, opts ...option.RequestOption) (res *Certificate, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if certificateID == "" {
		err = errors.New("missing required certificate_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("organization/certificates/%s", certificateID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// List uploaded certificates for this organization.
func (r *AdminOrganizationCertificateService) List(ctx context.Context, query AdminOrganizationCertificateListParams, opts ...option.RequestOption) (res *pagination.ConversationCursorPage[AdminOrganizationCertificateListResponse], err error) {
	var raw *http.Response
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "organization/certificates"
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

// List uploaded certificates for this organization.
func (r *AdminOrganizationCertificateService) ListAutoPaging(ctx context.Context, query AdminOrganizationCertificateListParams, opts ...option.RequestOption) *pagination.ConversationCursorPageAutoPager[AdminOrganizationCertificateListResponse] {
	return pagination.NewConversationCursorPageAutoPager(r.List(ctx, query, opts...))
}

// Delete a certificate from the organization.
//
// The certificate must be inactive for the organization and all projects.
func (r *AdminOrganizationCertificateService) Delete(ctx context.Context, certificateID string, opts ...option.RequestOption) (res *AdminOrganizationCertificateDeleteResponse, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if certificateID == "" {
		err = errors.New("missing required certificate_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("organization/certificates/%s", certificateID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Activate certificates at the organization level.
//
// You can atomically and idempotently activate up to 10 certificates at a time.
func (r *AdminOrganizationCertificateService) Activate(ctx context.Context, body AdminOrganizationCertificateActivateParams, opts ...option.RequestOption) (res *pagination.Page[AdminOrganizationCertificateActivateResponse], err error) {
	var raw *http.Response
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "organization/certificates/activate"
	cfg, err := requestconfig.NewRequestConfig(ctx, http.MethodPost, path, body, &res, opts...)
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

// Activate certificates at the organization level.
//
// You can atomically and idempotently activate up to 10 certificates at a time.
func (r *AdminOrganizationCertificateService) ActivateAutoPaging(ctx context.Context, body AdminOrganizationCertificateActivateParams, opts ...option.RequestOption) *pagination.PageAutoPager[AdminOrganizationCertificateActivateResponse] {
	return pagination.NewPageAutoPager(r.Activate(ctx, body, opts...))
}

// Deactivate certificates at the organization level.
//
// You can atomically and idempotently deactivate up to 10 certificates at a time.
func (r *AdminOrganizationCertificateService) Deactivate(ctx context.Context, body AdminOrganizationCertificateDeactivateParams, opts ...option.RequestOption) (res *pagination.Page[AdminOrganizationCertificateDeactivateResponse], err error) {
	var raw *http.Response
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "organization/certificates/deactivate"
	cfg, err := requestconfig.NewRequestConfig(ctx, http.MethodPost, path, body, &res, opts...)
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

// Deactivate certificates at the organization level.
//
// You can atomically and idempotently deactivate up to 10 certificates at a time.
func (r *AdminOrganizationCertificateService) DeactivateAutoPaging(ctx context.Context, body AdminOrganizationCertificateDeactivateParams, opts ...option.RequestOption) *pagination.PageAutoPager[AdminOrganizationCertificateDeactivateResponse] {
	return pagination.NewPageAutoPager(r.Deactivate(ctx, body, opts...))
}

// Represents an individual `certificate` uploaded to the organization.
type Certificate struct {
	// The identifier, which can be referenced in API endpoints
	ID                 string                        `json:"id" api:"required"`
	CertificateDetails CertificateCertificateDetails `json:"certificate_details" api:"required"`
	// The Unix timestamp (in seconds) of when the certificate was uploaded.
	CreatedAt int64 `json:"created_at" api:"required" format:"unixtime"`
	// The name of the certificate.
	Name string `json:"name" api:"required"`
	// The object type.
	//
	//   - If creating, updating, or getting a specific certificate, the object type is
	//     `certificate`.
	//   - If listing, activating, or deactivating certificates for the organization, the
	//     object type is `organization.certificate`.
	//   - If listing, activating, or deactivating certificates for a project, the object
	//     type is `organization.project.certificate`.
	//
	// Any of "certificate", "organization.certificate",
	// "organization.project.certificate".
	Object CertificateObject `json:"object" api:"required"`
	// Whether the certificate is currently active at the specified scope. Not returned
	// when getting details for a specific certificate.
	Active bool `json:"active"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                 respjson.Field
		CertificateDetails respjson.Field
		CreatedAt          respjson.Field
		Name               respjson.Field
		Object             respjson.Field
		Active             respjson.Field
		ExtraFields        map[string]respjson.Field
		raw                string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r Certificate) RawJSON() string { return r.JSON.raw }
func (r *Certificate) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type CertificateCertificateDetails struct {
	// The content of the certificate in PEM format.
	Content string `json:"content"`
	// The Unix timestamp (in seconds) of when the certificate expires.
	ExpiresAt int64 `json:"expires_at" format:"unixtime"`
	// The Unix timestamp (in seconds) of when the certificate becomes valid.
	ValidAt int64 `json:"valid_at" format:"unixtime"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Content     respjson.Field
		ExpiresAt   respjson.Field
		ValidAt     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CertificateCertificateDetails) RawJSON() string { return r.JSON.raw }
func (r *CertificateCertificateDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The object type.
//
//   - If creating, updating, or getting a specific certificate, the object type is
//     `certificate`.
//   - If listing, activating, or deactivating certificates for the organization, the
//     object type is `organization.certificate`.
//   - If listing, activating, or deactivating certificates for a project, the object
//     type is `organization.project.certificate`.
type CertificateObject string

const (
	CertificateObjectCertificate                    CertificateObject = "certificate"
	CertificateObjectOrganizationCertificate        CertificateObject = "organization.certificate"
	CertificateObjectOrganizationProjectCertificate CertificateObject = "organization.project.certificate"
)

// Represents an individual certificate configured at the organization level.
type AdminOrganizationCertificateListResponse struct {
	// The identifier, which can be referenced in API endpoints
	ID string `json:"id" api:"required"`
	// Whether the certificate is currently active at the organization level.
	Active             bool                                                       `json:"active" api:"required"`
	CertificateDetails AdminOrganizationCertificateListResponseCertificateDetails `json:"certificate_details" api:"required"`
	// The Unix timestamp (in seconds) of when the certificate was uploaded.
	CreatedAt int64 `json:"created_at" api:"required" format:"unixtime"`
	// The name of the certificate.
	Name string `json:"name" api:"required"`
	// The object type, which is always `organization.certificate`.
	Object constant.OrganizationCertificate `json:"object" default:"organization.certificate"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                 respjson.Field
		Active             respjson.Field
		CertificateDetails respjson.Field
		CreatedAt          respjson.Field
		Name               respjson.Field
		Object             respjson.Field
		ExtraFields        map[string]respjson.Field
		raw                string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationCertificateListResponse) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationCertificateListResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationCertificateListResponseCertificateDetails struct {
	// The Unix timestamp (in seconds) of when the certificate expires.
	ExpiresAt int64 `json:"expires_at" format:"unixtime"`
	// The Unix timestamp (in seconds) of when the certificate becomes valid.
	ValidAt int64 `json:"valid_at" format:"unixtime"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ExpiresAt   respjson.Field
		ValidAt     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationCertificateListResponseCertificateDetails) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationCertificateListResponseCertificateDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationCertificateDeleteResponse struct {
	// The ID of the certificate that was deleted.
	ID string `json:"id" api:"required"`
	// The object type, must be `certificate.deleted`.
	Object constant.CertificateDeleted `json:"object" default:"certificate.deleted"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Object      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationCertificateDeleteResponse) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationCertificateDeleteResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Represents an individual certificate configured at the organization level.
type AdminOrganizationCertificateActivateResponse struct {
	// The identifier, which can be referenced in API endpoints
	ID string `json:"id" api:"required"`
	// Whether the certificate is currently active at the organization level.
	Active             bool                                                           `json:"active" api:"required"`
	CertificateDetails AdminOrganizationCertificateActivateResponseCertificateDetails `json:"certificate_details" api:"required"`
	// The Unix timestamp (in seconds) of when the certificate was uploaded.
	CreatedAt int64 `json:"created_at" api:"required" format:"unixtime"`
	// The name of the certificate.
	Name string `json:"name" api:"required"`
	// The object type, which is always `organization.certificate`.
	Object constant.OrganizationCertificate `json:"object" default:"organization.certificate"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                 respjson.Field
		Active             respjson.Field
		CertificateDetails respjson.Field
		CreatedAt          respjson.Field
		Name               respjson.Field
		Object             respjson.Field
		ExtraFields        map[string]respjson.Field
		raw                string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationCertificateActivateResponse) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationCertificateActivateResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationCertificateActivateResponseCertificateDetails struct {
	// The Unix timestamp (in seconds) of when the certificate expires.
	ExpiresAt int64 `json:"expires_at" format:"unixtime"`
	// The Unix timestamp (in seconds) of when the certificate becomes valid.
	ValidAt int64 `json:"valid_at" format:"unixtime"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ExpiresAt   respjson.Field
		ValidAt     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationCertificateActivateResponseCertificateDetails) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationCertificateActivateResponseCertificateDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Represents an individual certificate configured at the organization level.
type AdminOrganizationCertificateDeactivateResponse struct {
	// The identifier, which can be referenced in API endpoints
	ID string `json:"id" api:"required"`
	// Whether the certificate is currently active at the organization level.
	Active             bool                                                             `json:"active" api:"required"`
	CertificateDetails AdminOrganizationCertificateDeactivateResponseCertificateDetails `json:"certificate_details" api:"required"`
	// The Unix timestamp (in seconds) of when the certificate was uploaded.
	CreatedAt int64 `json:"created_at" api:"required" format:"unixtime"`
	// The name of the certificate.
	Name string `json:"name" api:"required"`
	// The object type, which is always `organization.certificate`.
	Object constant.OrganizationCertificate `json:"object" default:"organization.certificate"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                 respjson.Field
		Active             respjson.Field
		CertificateDetails respjson.Field
		CreatedAt          respjson.Field
		Name               respjson.Field
		Object             respjson.Field
		ExtraFields        map[string]respjson.Field
		raw                string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationCertificateDeactivateResponse) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationCertificateDeactivateResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationCertificateDeactivateResponseCertificateDetails struct {
	// The Unix timestamp (in seconds) of when the certificate expires.
	ExpiresAt int64 `json:"expires_at" format:"unixtime"`
	// The Unix timestamp (in seconds) of when the certificate becomes valid.
	ValidAt int64 `json:"valid_at" format:"unixtime"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ExpiresAt   respjson.Field
		ValidAt     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationCertificateDeactivateResponseCertificateDetails) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationCertificateDeactivateResponseCertificateDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationCertificateNewParams struct {
	// The certificate content in PEM format
	Certificate string `json:"certificate" api:"required"`
	// An optional name for the certificate
	Name param.Opt[string] `json:"name,omitzero"`
	paramObj
}

func (r AdminOrganizationCertificateNewParams) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationCertificateNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationCertificateNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationCertificateGetParams struct {
	// A list of additional fields to include in the response. Currently the only
	// supported value is `content` to fetch the PEM content of the certificate.
	//
	// Any of "content".
	Include []string `query:"include,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [AdminOrganizationCertificateGetParams]'s query parameters
// as `url.Values`.
func (r AdminOrganizationCertificateGetParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type AdminOrganizationCertificateUpdateParams struct {
	// The updated name for the certificate
	Name param.Opt[string] `json:"name,omitzero"`
	paramObj
}

func (r AdminOrganizationCertificateUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationCertificateUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationCertificateUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationCertificateListParams struct {
	// A cursor for use in pagination. `after` is an object ID that defines your place
	// in the list. For instance, if you make a list request and receive 100 objects,
	// ending with obj_foo, your subsequent call can include after=obj_foo in order to
	// fetch the next page of the list.
	After param.Opt[string] `query:"after,omitzero" json:"-"`
	// A limit on the number of objects to be returned. Limit can range between 1 and
	// 100, and the default is 20.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Sort order by the `created_at` timestamp of the objects. `asc` for ascending
	// order and `desc` for descending order.
	//
	// Any of "asc", "desc".
	Order AdminOrganizationCertificateListParamsOrder `query:"order,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [AdminOrganizationCertificateListParams]'s query parameters
// as `url.Values`.
func (r AdminOrganizationCertificateListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Sort order by the `created_at` timestamp of the objects. `asc` for ascending
// order and `desc` for descending order.
type AdminOrganizationCertificateListParamsOrder string

const (
	AdminOrganizationCertificateListParamsOrderAsc  AdminOrganizationCertificateListParamsOrder = "asc"
	AdminOrganizationCertificateListParamsOrderDesc AdminOrganizationCertificateListParamsOrder = "desc"
)

type AdminOrganizationCertificateActivateParams struct {
	CertificateIDs []string `json:"certificate_ids,omitzero" api:"required"`
	paramObj
}

func (r AdminOrganizationCertificateActivateParams) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationCertificateActivateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationCertificateActivateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationCertificateDeactivateParams struct {
	CertificateIDs []string `json:"certificate_ids,omitzero" api:"required"`
	paramObj
}

func (r AdminOrganizationCertificateDeactivateParams) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationCertificateDeactivateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationCertificateDeactivateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

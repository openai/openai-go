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
)

// AdminOrganizationProjectCertificateService contains methods and other services
// that help with interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAdminOrganizationProjectCertificateService] method instead.
type AdminOrganizationProjectCertificateService struct {
	Options []option.RequestOption
}

// NewAdminOrganizationProjectCertificateService generates a new service that
// applies the given options to each request. These options are applied after the
// parent client's options (if there is one), and before any request-specific
// options.
func NewAdminOrganizationProjectCertificateService(opts ...option.RequestOption) (r AdminOrganizationProjectCertificateService) {
	r = AdminOrganizationProjectCertificateService{}
	r.Options = opts
	return
}

// List certificates for this project.
func (r *AdminOrganizationProjectCertificateService) List(ctx context.Context, projectID string, query AdminOrganizationProjectCertificateListParams, opts ...option.RequestOption) (res *pagination.ConversationCursorPage[Certificate], err error) {
	var raw *http.Response
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	if projectID == "" {
		err = errors.New("missing required project_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("organization/projects/%s/certificates", projectID)
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

// List certificates for this project.
func (r *AdminOrganizationProjectCertificateService) ListAutoPaging(ctx context.Context, projectID string, query AdminOrganizationProjectCertificateListParams, opts ...option.RequestOption) *pagination.ConversationCursorPageAutoPager[Certificate] {
	return pagination.NewConversationCursorPageAutoPager(r.List(ctx, projectID, query, opts...))
}

// Activate certificates at the project level.
//
// You can atomically and idempotently activate up to 10 certificates at a time.
func (r *AdminOrganizationProjectCertificateService) Activate(ctx context.Context, projectID string, body AdminOrganizationProjectCertificateActivateParams, opts ...option.RequestOption) (res *pagination.Page[Certificate], err error) {
	var raw *http.Response
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	if projectID == "" {
		err = errors.New("missing required project_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("organization/projects/%s/certificates/activate", projectID)
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

// Activate certificates at the project level.
//
// You can atomically and idempotently activate up to 10 certificates at a time.
func (r *AdminOrganizationProjectCertificateService) ActivateAutoPaging(ctx context.Context, projectID string, body AdminOrganizationProjectCertificateActivateParams, opts ...option.RequestOption) *pagination.PageAutoPager[Certificate] {
	return pagination.NewPageAutoPager(r.Activate(ctx, projectID, body, opts...))
}

// Deactivate certificates at the project level. You can atomically and
// idempotently deactivate up to 10 certificates at a time.
func (r *AdminOrganizationProjectCertificateService) Deactivate(ctx context.Context, projectID string, body AdminOrganizationProjectCertificateDeactivateParams, opts ...option.RequestOption) (res *pagination.Page[Certificate], err error) {
	var raw *http.Response
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	if projectID == "" {
		err = errors.New("missing required project_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("organization/projects/%s/certificates/deactivate", projectID)
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

// Deactivate certificates at the project level. You can atomically and
// idempotently deactivate up to 10 certificates at a time.
func (r *AdminOrganizationProjectCertificateService) DeactivateAutoPaging(ctx context.Context, projectID string, body AdminOrganizationProjectCertificateDeactivateParams, opts ...option.RequestOption) *pagination.PageAutoPager[Certificate] {
	return pagination.NewPageAutoPager(r.Deactivate(ctx, projectID, body, opts...))
}

type AdminOrganizationProjectCertificateListParams struct {
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
	Order AdminOrganizationProjectCertificateListParamsOrder `query:"order,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [AdminOrganizationProjectCertificateListParams]'s query
// parameters as `url.Values`.
func (r AdminOrganizationProjectCertificateListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Sort order by the `created_at` timestamp of the objects. `asc` for ascending
// order and `desc` for descending order.
type AdminOrganizationProjectCertificateListParamsOrder string

const (
	AdminOrganizationProjectCertificateListParamsOrderAsc  AdminOrganizationProjectCertificateListParamsOrder = "asc"
	AdminOrganizationProjectCertificateListParamsOrderDesc AdminOrganizationProjectCertificateListParamsOrder = "desc"
)

type AdminOrganizationProjectCertificateActivateParams struct {
	CertificateIDs []string `json:"certificate_ids,omitzero" api:"required"`
	paramObj
}

func (r AdminOrganizationProjectCertificateActivateParams) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationProjectCertificateActivateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationProjectCertificateActivateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationProjectCertificateDeactivateParams struct {
	CertificateIDs []string `json:"certificate_ids,omitzero" api:"required"`
	paramObj
}

func (r AdminOrganizationProjectCertificateDeactivateParams) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationProjectCertificateDeactivateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationProjectCertificateDeactivateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

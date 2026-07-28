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

// AdminOrganizationSpendLimitService contains methods and other services that help
// with interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAdminOrganizationSpendLimitService] method instead.
type AdminOrganizationSpendLimitService struct {
	Options []option.RequestOption
}

// NewAdminOrganizationSpendLimitService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewAdminOrganizationSpendLimitService(opts ...option.RequestOption) (r AdminOrganizationSpendLimitService) {
	r = AdminOrganizationSpendLimitService{}
	r.Options = opts
	return
}

// Get the organization's hard spend limit.
func (r *AdminOrganizationSpendLimitService) Get(ctx context.Context, opts ...option.RequestOption) (res *OrganizationSpendLimit, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	path := "organization/spend_limit"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Create or replace the organization's hard spend limit.
func (r *AdminOrganizationSpendLimitService) Update(ctx context.Context, body AdminOrganizationSpendLimitUpdateParams, opts ...option.RequestOption) (res *OrganizationSpendLimit, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	path := "organization/spend_limit"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Delete the organization's hard spend limit.
func (r *AdminOrganizationSpendLimitService) Delete(ctx context.Context, opts ...option.RequestOption) (res *OrganizationSpendLimitDeleted, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	path := "organization/spend_limit"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Represents a hard spend limit configured at the organization level.
type OrganizationSpendLimit struct {
	// The currency for the threshold amount. Currently, only `USD` is supported.
	Currency OrganizationSpendLimitCurrency `json:"currency" api:"required"`
	// The current enforcement state of the hard spend limit.
	Enforcement OrganizationSpendLimitEnforcement `json:"enforcement" api:"required"`
	// The time interval for evaluating spend against the threshold. Currently, only
	// `month` is supported.
	Interval OrganizationSpendLimitInterval `json:"interval" api:"required"`
	// The object type, which is always `organization.spend_limit`.
	Object constant.OrganizationSpendLimit `json:"object" default:"organization.spend_limit"`
	// The hard spend limit amount, in cents.
	ThresholdAmount int64 `json:"threshold_amount" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Currency        respjson.Field
		Enforcement     respjson.Field
		Interval        respjson.Field
		Object          respjson.Field
		ThresholdAmount respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r OrganizationSpendLimit) RawJSON() string { return r.JSON.raw }
func (r *OrganizationSpendLimit) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The currency for the threshold amount. Currently, only `USD` is supported.
type OrganizationSpendLimitCurrency string

const (
	OrganizationSpendLimitCurrencyUsd OrganizationSpendLimitCurrency = "USD"
)

// The current enforcement state of the hard spend limit.
type OrganizationSpendLimitEnforcement struct {
	// Whether the hard spend limit is currently enforcing.
	Status string `json:"status" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Status      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r OrganizationSpendLimitEnforcement) RawJSON() string { return r.JSON.raw }
func (r *OrganizationSpendLimitEnforcement) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The time interval for evaluating spend against the threshold. Currently, only
// `month` is supported.
type OrganizationSpendLimitInterval string

const (
	OrganizationSpendLimitIntervalMonth OrganizationSpendLimitInterval = "month"
)

// Confirmation payload returned after deleting an organization hard spend limit.
type OrganizationSpendLimitDeleted struct {
	// Whether the hard spend limit was deleted.
	Deleted bool `json:"deleted" api:"required"`
	// The object type, which is always `organization.spend_limit.deleted`.
	Object constant.OrganizationSpendLimitDeleted `json:"object" default:"organization.spend_limit.deleted"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Deleted     respjson.Field
		Object      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r OrganizationSpendLimitDeleted) RawJSON() string { return r.JSON.raw }
func (r *OrganizationSpendLimitDeleted) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationSpendLimitUpdateParams struct {
	// The currency for the threshold amount. Currently, only `USD` is supported.
	//
	// Any of "USD".
	Currency AdminOrganizationSpendLimitUpdateParamsCurrency `json:"currency,omitzero" api:"required"`
	// The time interval for evaluating spend against the threshold. Currently, only
	// `month` is supported.
	//
	// Any of "month".
	Interval AdminOrganizationSpendLimitUpdateParamsInterval `json:"interval,omitzero" api:"required"`
	// The hard spend limit amount, in cents.
	ThresholdAmount int64 `json:"threshold_amount" api:"required"`
	paramObj
}

func (r AdminOrganizationSpendLimitUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationSpendLimitUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationSpendLimitUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The currency for the threshold amount. Currently, only `USD` is supported.
type AdminOrganizationSpendLimitUpdateParamsCurrency string

const (
	AdminOrganizationSpendLimitUpdateParamsCurrencyUsd AdminOrganizationSpendLimitUpdateParamsCurrency = "USD"
)

// The time interval for evaluating spend against the threshold. Currently, only
// `month` is supported.
type AdminOrganizationSpendLimitUpdateParamsInterval string

const (
	AdminOrganizationSpendLimitUpdateParamsIntervalMonth AdminOrganizationSpendLimitUpdateParamsInterval = "month"
)

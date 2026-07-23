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

// AdminOrganizationProjectSpendLimitService contains methods and other services
// that help with interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAdminOrganizationProjectSpendLimitService] method instead.
type AdminOrganizationProjectSpendLimitService struct {
	Options []option.RequestOption
}

// NewAdminOrganizationProjectSpendLimitService generates a new service that
// applies the given options to each request. These options are applied after the
// parent client's options (if there is one), and before any request-specific
// options.
func NewAdminOrganizationProjectSpendLimitService(opts ...option.RequestOption) (r AdminOrganizationProjectSpendLimitService) {
	r = AdminOrganizationProjectSpendLimitService{}
	r.Options = opts
	return
}

// Get a project's hard spend limit.
func (r *AdminOrganizationProjectSpendLimitService) Get(ctx context.Context, projectID string, opts ...option.RequestOption) (res *ProjectSpendLimit, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if projectID == "" {
		err = errors.New("missing required project_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("organization/projects/%s/spend_limit", projectID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Create or replace a project's hard spend limit.
func (r *AdminOrganizationProjectSpendLimitService) Update(ctx context.Context, projectID string, body AdminOrganizationProjectSpendLimitUpdateParams, opts ...option.RequestOption) (res *ProjectSpendLimit, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if projectID == "" {
		err = errors.New("missing required project_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("organization/projects/%s/spend_limit", projectID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Delete a project's hard spend limit.
func (r *AdminOrganizationProjectSpendLimitService) Delete(ctx context.Context, projectID string, opts ...option.RequestOption) (res *ProjectSpendLimitDeleted, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if projectID == "" {
		err = errors.New("missing required project_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("organization/projects/%s/spend_limit", projectID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Represents a hard spend limit configured at the project level.
type ProjectSpendLimit struct {
	// The currency for the threshold amount. Currently, only `USD` is supported.
	Currency ProjectSpendLimitCurrency `json:"currency" api:"required"`
	// The current enforcement state of the hard spend limit.
	Enforcement ProjectSpendLimitEnforcement `json:"enforcement" api:"required"`
	// The time interval for evaluating spend against the threshold. Currently, only
	// `month` is supported.
	Interval ProjectSpendLimitInterval `json:"interval" api:"required"`
	// The object type, which is always `project.spend_limit`.
	Object constant.ProjectSpendLimit `json:"object" default:"project.spend_limit"`
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
func (r ProjectSpendLimit) RawJSON() string { return r.JSON.raw }
func (r *ProjectSpendLimit) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The currency for the threshold amount. Currently, only `USD` is supported.
type ProjectSpendLimitCurrency string

const (
	ProjectSpendLimitCurrencyUsd ProjectSpendLimitCurrency = "USD"
)

// The current enforcement state of the hard spend limit.
type ProjectSpendLimitEnforcement struct {
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
func (r ProjectSpendLimitEnforcement) RawJSON() string { return r.JSON.raw }
func (r *ProjectSpendLimitEnforcement) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The time interval for evaluating spend against the threshold. Currently, only
// `month` is supported.
type ProjectSpendLimitInterval string

const (
	ProjectSpendLimitIntervalMonth ProjectSpendLimitInterval = "month"
)

// Confirmation payload returned after deleting a project hard spend limit.
type ProjectSpendLimitDeleted struct {
	// Whether the hard spend limit was deleted.
	Deleted bool `json:"deleted" api:"required"`
	// The object type, which is always `project.spend_limit.deleted`.
	Object constant.ProjectSpendLimitDeleted `json:"object" default:"project.spend_limit.deleted"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Deleted     respjson.Field
		Object      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ProjectSpendLimitDeleted) RawJSON() string { return r.JSON.raw }
func (r *ProjectSpendLimitDeleted) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationProjectSpendLimitUpdateParams struct {
	// The currency for the threshold amount. Currently, only `USD` is supported.
	//
	// Any of "USD".
	Currency AdminOrganizationProjectSpendLimitUpdateParamsCurrency `json:"currency,omitzero" api:"required"`
	// The time interval for evaluating spend against the threshold. Currently, only
	// `month` is supported.
	//
	// Any of "month".
	Interval AdminOrganizationProjectSpendLimitUpdateParamsInterval `json:"interval,omitzero" api:"required"`
	// The hard spend limit amount, in cents.
	ThresholdAmount int64 `json:"threshold_amount" api:"required"`
	paramObj
}

func (r AdminOrganizationProjectSpendLimitUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationProjectSpendLimitUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationProjectSpendLimitUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The currency for the threshold amount. Currently, only `USD` is supported.
type AdminOrganizationProjectSpendLimitUpdateParamsCurrency string

const (
	AdminOrganizationProjectSpendLimitUpdateParamsCurrencyUsd AdminOrganizationProjectSpendLimitUpdateParamsCurrency = "USD"
)

// The time interval for evaluating spend against the threshold. Currently, only
// `month` is supported.
type AdminOrganizationProjectSpendLimitUpdateParamsInterval string

const (
	AdminOrganizationProjectSpendLimitUpdateParamsIntervalMonth AdminOrganizationProjectSpendLimitUpdateParamsInterval = "month"
)

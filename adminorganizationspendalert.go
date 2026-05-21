// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

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

// AdminOrganizationSpendAlertService contains methods and other services that help
// with interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAdminOrganizationSpendAlertService] method instead.
type AdminOrganizationSpendAlertService struct {
	Options []option.RequestOption
}

// NewAdminOrganizationSpendAlertService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewAdminOrganizationSpendAlertService(opts ...option.RequestOption) (r AdminOrganizationSpendAlertService) {
	r = AdminOrganizationSpendAlertService{}
	r.Options = opts
	return
}

// Creates an organization spend alert.
func (r *AdminOrganizationSpendAlertService) New(ctx context.Context, body AdminOrganizationSpendAlertNewParams, opts ...option.RequestOption) (res *OrganizationSpendAlert, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	path := "organization/spend_alerts"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Updates an organization spend alert.
func (r *AdminOrganizationSpendAlertService) Update(ctx context.Context, alertID string, body AdminOrganizationSpendAlertUpdateParams, opts ...option.RequestOption) (res *OrganizationSpendAlert, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if alertID == "" {
		err = errors.New("missing required alert_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("organization/spend_alerts/%s", alertID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Lists organization spend alerts.
func (r *AdminOrganizationSpendAlertService) List(ctx context.Context, query AdminOrganizationSpendAlertListParams, opts ...option.RequestOption) (res *pagination.ConversationCursorPage[OrganizationSpendAlert], err error) {
	var raw *http.Response
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "organization/spend_alerts"
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

// Lists organization spend alerts.
func (r *AdminOrganizationSpendAlertService) ListAutoPaging(ctx context.Context, query AdminOrganizationSpendAlertListParams, opts ...option.RequestOption) *pagination.ConversationCursorPageAutoPager[OrganizationSpendAlert] {
	return pagination.NewConversationCursorPageAutoPager(r.List(ctx, query, opts...))
}

// Deletes an organization spend alert.
func (r *AdminOrganizationSpendAlertService) Delete(ctx context.Context, alertID string, opts ...option.RequestOption) (res *OrganizationSpendAlertDeleted, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if alertID == "" {
		err = errors.New("missing required alert_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("organization/spend_alerts/%s", alertID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Represents a spend alert configured at the organization level.
type OrganizationSpendAlert struct {
	// The identifier, which can be referenced in API endpoints.
	ID string `json:"id" api:"required"`
	// The currency for the threshold amount.
	//
	// Any of "USD".
	Currency OrganizationSpendAlertCurrency `json:"currency" api:"required"`
	// The time interval for evaluating spend against the threshold.
	//
	// Any of "month".
	Interval OrganizationSpendAlertInterval `json:"interval" api:"required"`
	// Email notification settings for a spend alert.
	NotificationChannel OrganizationSpendAlertNotificationChannel `json:"notification_channel" api:"required"`
	// The object type, which is always `organization.spend_alert`.
	Object constant.OrganizationSpendAlert `json:"object" default:"organization.spend_alert"`
	// The alert threshold amount, in cents.
	ThresholdAmount int64 `json:"threshold_amount" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                  respjson.Field
		Currency            respjson.Field
		Interval            respjson.Field
		NotificationChannel respjson.Field
		Object              respjson.Field
		ThresholdAmount     respjson.Field
		ExtraFields         map[string]respjson.Field
		raw                 string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r OrganizationSpendAlert) RawJSON() string { return r.JSON.raw }
func (r *OrganizationSpendAlert) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The currency for the threshold amount.
type OrganizationSpendAlertCurrency string

const (
	OrganizationSpendAlertCurrencyUsd OrganizationSpendAlertCurrency = "USD"
)

// The time interval for evaluating spend against the threshold.
type OrganizationSpendAlertInterval string

const (
	OrganizationSpendAlertIntervalMonth OrganizationSpendAlertInterval = "month"
)

// Email notification settings for a spend alert.
type OrganizationSpendAlertNotificationChannel struct {
	// Email addresses that receive the spend alert notification.
	Recipients []string `json:"recipients" api:"required"`
	// The notification channel type. Currently only `email` is supported.
	Type constant.Email `json:"type" default:"email"`
	// Optional subject prefix for alert emails.
	SubjectPrefix string `json:"subject_prefix" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Recipients    respjson.Field
		Type          respjson.Field
		SubjectPrefix respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r OrganizationSpendAlertNotificationChannel) RawJSON() string { return r.JSON.raw }
func (r *OrganizationSpendAlertNotificationChannel) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Confirmation payload returned after deleting an organization spend alert.
type OrganizationSpendAlertDeleted struct {
	// The deleted spend alert ID.
	ID string `json:"id" api:"required"`
	// Whether the spend alert was deleted.
	Deleted bool `json:"deleted" api:"required"`
	// Always `organization.spend_alert.deleted`.
	Object constant.OrganizationSpendAlertDeleted `json:"object" default:"organization.spend_alert.deleted"`
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
func (r OrganizationSpendAlertDeleted) RawJSON() string { return r.JSON.raw }
func (r *OrganizationSpendAlertDeleted) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationSpendAlertNewParams struct {
	// The currency for the threshold amount.
	//
	// Any of "USD".
	Currency AdminOrganizationSpendAlertNewParamsCurrency `json:"currency,omitzero" api:"required"`
	// The time interval for evaluating spend against the threshold.
	//
	// Any of "month".
	Interval AdminOrganizationSpendAlertNewParamsInterval `json:"interval,omitzero" api:"required"`
	// Email notification settings for a spend alert.
	NotificationChannel AdminOrganizationSpendAlertNewParamsNotificationChannel `json:"notification_channel,omitzero" api:"required"`
	// The alert threshold amount, in cents.
	ThresholdAmount int64 `json:"threshold_amount" api:"required"`
	paramObj
}

func (r AdminOrganizationSpendAlertNewParams) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationSpendAlertNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationSpendAlertNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The currency for the threshold amount.
type AdminOrganizationSpendAlertNewParamsCurrency string

const (
	AdminOrganizationSpendAlertNewParamsCurrencyUsd AdminOrganizationSpendAlertNewParamsCurrency = "USD"
)

// The time interval for evaluating spend against the threshold.
type AdminOrganizationSpendAlertNewParamsInterval string

const (
	AdminOrganizationSpendAlertNewParamsIntervalMonth AdminOrganizationSpendAlertNewParamsInterval = "month"
)

// Email notification settings for a spend alert.
//
// The properties Recipients, Type are required.
type AdminOrganizationSpendAlertNewParamsNotificationChannel struct {
	// Email addresses that receive the spend alert notification.
	Recipients []string `json:"recipients,omitzero" api:"required"`
	// Optional subject prefix for alert emails.
	SubjectPrefix param.Opt[string] `json:"subject_prefix,omitzero"`
	// The notification channel type. Currently only `email` is supported.
	//
	// This field can be elided, and will marshal its zero value as "email".
	Type constant.Email `json:"type" default:"email"`
	paramObj
}

func (r AdminOrganizationSpendAlertNewParamsNotificationChannel) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationSpendAlertNewParamsNotificationChannel
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationSpendAlertNewParamsNotificationChannel) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationSpendAlertUpdateParams struct {
	// The currency for the threshold amount.
	//
	// Any of "USD".
	Currency AdminOrganizationSpendAlertUpdateParamsCurrency `json:"currency,omitzero" api:"required"`
	// The time interval for evaluating spend against the threshold.
	//
	// Any of "month".
	Interval AdminOrganizationSpendAlertUpdateParamsInterval `json:"interval,omitzero" api:"required"`
	// Email notification settings for a spend alert.
	NotificationChannel AdminOrganizationSpendAlertUpdateParamsNotificationChannel `json:"notification_channel,omitzero" api:"required"`
	// The alert threshold amount, in cents.
	ThresholdAmount int64 `json:"threshold_amount" api:"required"`
	paramObj
}

func (r AdminOrganizationSpendAlertUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationSpendAlertUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationSpendAlertUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The currency for the threshold amount.
type AdminOrganizationSpendAlertUpdateParamsCurrency string

const (
	AdminOrganizationSpendAlertUpdateParamsCurrencyUsd AdminOrganizationSpendAlertUpdateParamsCurrency = "USD"
)

// The time interval for evaluating spend against the threshold.
type AdminOrganizationSpendAlertUpdateParamsInterval string

const (
	AdminOrganizationSpendAlertUpdateParamsIntervalMonth AdminOrganizationSpendAlertUpdateParamsInterval = "month"
)

// Email notification settings for a spend alert.
//
// The properties Recipients, Type are required.
type AdminOrganizationSpendAlertUpdateParamsNotificationChannel struct {
	// Email addresses that receive the spend alert notification.
	Recipients []string `json:"recipients,omitzero" api:"required"`
	// Optional subject prefix for alert emails.
	SubjectPrefix param.Opt[string] `json:"subject_prefix,omitzero"`
	// The notification channel type. Currently only `email` is supported.
	//
	// This field can be elided, and will marshal its zero value as "email".
	Type constant.Email `json:"type" default:"email"`
	paramObj
}

func (r AdminOrganizationSpendAlertUpdateParamsNotificationChannel) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationSpendAlertUpdateParamsNotificationChannel
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationSpendAlertUpdateParamsNotificationChannel) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationSpendAlertListParams struct {
	// Cursor for pagination. Provide the ID of the last spend alert from the previous
	// response to fetch the next page.
	After param.Opt[string] `query:"after,omitzero" json:"-"`
	// Cursor for pagination. Provide the ID of the first spend alert from the previous
	// response to fetch the previous page.
	Before param.Opt[string] `query:"before,omitzero" json:"-"`
	// A limit on the number of spend alerts to return. Defaults to 20.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Sort order for the returned spend alerts.
	//
	// Any of "asc", "desc".
	Order AdminOrganizationSpendAlertListParamsOrder `query:"order,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [AdminOrganizationSpendAlertListParams]'s query parameters
// as `url.Values`.
func (r AdminOrganizationSpendAlertListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Sort order for the returned spend alerts.
type AdminOrganizationSpendAlertListParamsOrder string

const (
	AdminOrganizationSpendAlertListParamsOrderAsc  AdminOrganizationSpendAlertListParamsOrder = "asc"
	AdminOrganizationSpendAlertListParamsOrderDesc AdminOrganizationSpendAlertListParamsOrder = "desc"
)

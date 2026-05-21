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

// AdminOrganizationProjectSpendAlertService contains methods and other services
// that help with interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAdminOrganizationProjectSpendAlertService] method instead.
type AdminOrganizationProjectSpendAlertService struct {
	Options []option.RequestOption
}

// NewAdminOrganizationProjectSpendAlertService generates a new service that
// applies the given options to each request. These options are applied after the
// parent client's options (if there is one), and before any request-specific
// options.
func NewAdminOrganizationProjectSpendAlertService(opts ...option.RequestOption) (r AdminOrganizationProjectSpendAlertService) {
	r = AdminOrganizationProjectSpendAlertService{}
	r.Options = opts
	return
}

// Creates a project spend alert.
func (r *AdminOrganizationProjectSpendAlertService) New(ctx context.Context, projectID string, body AdminOrganizationProjectSpendAlertNewParams, opts ...option.RequestOption) (res *ProjectSpendAlert, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if projectID == "" {
		err = errors.New("missing required project_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("organization/projects/%s/spend_alerts", projectID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Updates a project spend alert.
func (r *AdminOrganizationProjectSpendAlertService) Update(ctx context.Context, projectID string, alertID string, body AdminOrganizationProjectSpendAlertUpdateParams, opts ...option.RequestOption) (res *ProjectSpendAlert, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if projectID == "" {
		err = errors.New("missing required project_id parameter")
		return nil, err
	}
	if alertID == "" {
		err = errors.New("missing required alert_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("organization/projects/%s/spend_alerts/%s", projectID, alertID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Lists project spend alerts.
func (r *AdminOrganizationProjectSpendAlertService) List(ctx context.Context, projectID string, query AdminOrganizationProjectSpendAlertListParams, opts ...option.RequestOption) (res *pagination.ConversationCursorPage[ProjectSpendAlert], err error) {
	var raw *http.Response
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	if projectID == "" {
		err = errors.New("missing required project_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("organization/projects/%s/spend_alerts", projectID)
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

// Lists project spend alerts.
func (r *AdminOrganizationProjectSpendAlertService) ListAutoPaging(ctx context.Context, projectID string, query AdminOrganizationProjectSpendAlertListParams, opts ...option.RequestOption) *pagination.ConversationCursorPageAutoPager[ProjectSpendAlert] {
	return pagination.NewConversationCursorPageAutoPager(r.List(ctx, projectID, query, opts...))
}

// Deletes a project spend alert.
func (r *AdminOrganizationProjectSpendAlertService) Delete(ctx context.Context, projectID string, alertID string, opts ...option.RequestOption) (res *ProjectSpendAlertDeleted, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if projectID == "" {
		err = errors.New("missing required project_id parameter")
		return nil, err
	}
	if alertID == "" {
		err = errors.New("missing required alert_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("organization/projects/%s/spend_alerts/%s", projectID, alertID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Represents a spend alert configured at the project level.
type ProjectSpendAlert struct {
	// The identifier, which can be referenced in API endpoints.
	ID string `json:"id" api:"required"`
	// The currency for the threshold amount.
	//
	// Any of "USD".
	Currency ProjectSpendAlertCurrency `json:"currency" api:"required"`
	// The time interval for evaluating spend against the threshold.
	//
	// Any of "month".
	Interval ProjectSpendAlertInterval `json:"interval" api:"required"`
	// Email notification settings for a spend alert.
	NotificationChannel ProjectSpendAlertNotificationChannel `json:"notification_channel" api:"required"`
	// The object type, which is always `project.spend_alert`.
	Object constant.ProjectSpendAlert `json:"object" default:"project.spend_alert"`
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
func (r ProjectSpendAlert) RawJSON() string { return r.JSON.raw }
func (r *ProjectSpendAlert) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The currency for the threshold amount.
type ProjectSpendAlertCurrency string

const (
	ProjectSpendAlertCurrencyUsd ProjectSpendAlertCurrency = "USD"
)

// The time interval for evaluating spend against the threshold.
type ProjectSpendAlertInterval string

const (
	ProjectSpendAlertIntervalMonth ProjectSpendAlertInterval = "month"
)

// Email notification settings for a spend alert.
type ProjectSpendAlertNotificationChannel struct {
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
func (r ProjectSpendAlertNotificationChannel) RawJSON() string { return r.JSON.raw }
func (r *ProjectSpendAlertNotificationChannel) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Confirmation payload returned after deleting a project spend alert.
type ProjectSpendAlertDeleted struct {
	// The deleted spend alert ID.
	ID string `json:"id" api:"required"`
	// Whether the spend alert was deleted.
	Deleted bool `json:"deleted" api:"required"`
	// Always `project.spend_alert.deleted`.
	Object constant.ProjectSpendAlertDeleted `json:"object" default:"project.spend_alert.deleted"`
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
func (r ProjectSpendAlertDeleted) RawJSON() string { return r.JSON.raw }
func (r *ProjectSpendAlertDeleted) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationProjectSpendAlertNewParams struct {
	// The currency for the threshold amount.
	//
	// Any of "USD".
	Currency AdminOrganizationProjectSpendAlertNewParamsCurrency `json:"currency,omitzero" api:"required"`
	// The time interval for evaluating spend against the threshold.
	//
	// Any of "month".
	Interval AdminOrganizationProjectSpendAlertNewParamsInterval `json:"interval,omitzero" api:"required"`
	// Email notification settings for a spend alert.
	NotificationChannel AdminOrganizationProjectSpendAlertNewParamsNotificationChannel `json:"notification_channel,omitzero" api:"required"`
	// The alert threshold amount, in cents.
	ThresholdAmount int64 `json:"threshold_amount" api:"required"`
	paramObj
}

func (r AdminOrganizationProjectSpendAlertNewParams) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationProjectSpendAlertNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationProjectSpendAlertNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The currency for the threshold amount.
type AdminOrganizationProjectSpendAlertNewParamsCurrency string

const (
	AdminOrganizationProjectSpendAlertNewParamsCurrencyUsd AdminOrganizationProjectSpendAlertNewParamsCurrency = "USD"
)

// The time interval for evaluating spend against the threshold.
type AdminOrganizationProjectSpendAlertNewParamsInterval string

const (
	AdminOrganizationProjectSpendAlertNewParamsIntervalMonth AdminOrganizationProjectSpendAlertNewParamsInterval = "month"
)

// Email notification settings for a spend alert.
//
// The properties Recipients, Type are required.
type AdminOrganizationProjectSpendAlertNewParamsNotificationChannel struct {
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

func (r AdminOrganizationProjectSpendAlertNewParamsNotificationChannel) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationProjectSpendAlertNewParamsNotificationChannel
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationProjectSpendAlertNewParamsNotificationChannel) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationProjectSpendAlertUpdateParams struct {
	// The currency for the threshold amount.
	//
	// Any of "USD".
	Currency AdminOrganizationProjectSpendAlertUpdateParamsCurrency `json:"currency,omitzero" api:"required"`
	// The time interval for evaluating spend against the threshold.
	//
	// Any of "month".
	Interval AdminOrganizationProjectSpendAlertUpdateParamsInterval `json:"interval,omitzero" api:"required"`
	// Email notification settings for a spend alert.
	NotificationChannel AdminOrganizationProjectSpendAlertUpdateParamsNotificationChannel `json:"notification_channel,omitzero" api:"required"`
	// The alert threshold amount, in cents.
	ThresholdAmount int64 `json:"threshold_amount" api:"required"`
	paramObj
}

func (r AdminOrganizationProjectSpendAlertUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationProjectSpendAlertUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationProjectSpendAlertUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The currency for the threshold amount.
type AdminOrganizationProjectSpendAlertUpdateParamsCurrency string

const (
	AdminOrganizationProjectSpendAlertUpdateParamsCurrencyUsd AdminOrganizationProjectSpendAlertUpdateParamsCurrency = "USD"
)

// The time interval for evaluating spend against the threshold.
type AdminOrganizationProjectSpendAlertUpdateParamsInterval string

const (
	AdminOrganizationProjectSpendAlertUpdateParamsIntervalMonth AdminOrganizationProjectSpendAlertUpdateParamsInterval = "month"
)

// Email notification settings for a spend alert.
//
// The properties Recipients, Type are required.
type AdminOrganizationProjectSpendAlertUpdateParamsNotificationChannel struct {
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

func (r AdminOrganizationProjectSpendAlertUpdateParamsNotificationChannel) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationProjectSpendAlertUpdateParamsNotificationChannel
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationProjectSpendAlertUpdateParamsNotificationChannel) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationProjectSpendAlertListParams struct {
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
	Order AdminOrganizationProjectSpendAlertListParamsOrder `query:"order,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [AdminOrganizationProjectSpendAlertListParams]'s query
// parameters as `url.Values`.
func (r AdminOrganizationProjectSpendAlertListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Sort order for the returned spend alerts.
type AdminOrganizationProjectSpendAlertListParamsOrder string

const (
	AdminOrganizationProjectSpendAlertListParamsOrderAsc  AdminOrganizationProjectSpendAlertListParamsOrder = "asc"
	AdminOrganizationProjectSpendAlertListParamsOrderDesc AdminOrganizationProjectSpendAlertListParamsOrder = "desc"
)

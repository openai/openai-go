// File generated from our OpenAPI spec by Castiron. See CONTRIBUTING.md for details.

package openai

import (
	"context"
	"errors"
	"net/http"
	"slices"

	"github.com/openai/openai-go/v3/internal/apijson"
	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/respjson"
	"github.com/openai/openai-go/v3/shared/constant"
)

// SafetyAlertService contains methods and other services that help with
// interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewSafetyAlertService] method instead.
type SafetyAlertService struct {
	Options []option.RequestOption
}

// NewSafetyAlertService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewSafetyAlertService(opts ...option.RequestOption) (r SafetyAlertService) {
	r = SafetyAlertService{}
	r.Options = requestconfig.InheritedOptions(opts...)
	return
}

// Get a safety alert belonging to the authenticated API project.
func (r *SafetyAlertService) Get(ctx context.Context, id string, opts ...option.RequestOption) (res *SafetyAlert, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("safety/alerts/%s", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

type SafetyAlert struct {
	ID        string `json:"id" api:"required"`
	CreatedAt int64  `json:"created_at" api:"required" format:"unixtime"`
	// Any of "potentially_unintended_data_transfer",
	// "potentially_unintended_data_access",
	// "potentially_unintended_destructive_activity", "other".
	ErrorType SafetyAlertErrorType `json:"error_type" api:"required"`
	Model     string               `json:"model" api:"required"`
	Object    constant.SafetyAlert `json:"object" default:"safety.alert"`
	// A customer-safe description derived from error_type, or null for zero data
	// retention requests.
	Reason    string `json:"reason" api:"required"`
	RequestID string `json:"request_id" api:"required"`
	// Whether block registration succeeded for this request. This does not confirm
	// that response execution stopped.
	RequestPaused bool   `json:"request_paused" api:"required"`
	ResponseID    string `json:"response_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID            respjson.Field
		CreatedAt     respjson.Field
		ErrorType     respjson.Field
		Model         respjson.Field
		Object        respjson.Field
		Reason        respjson.Field
		RequestID     respjson.Field
		RequestPaused respjson.Field
		ResponseID    respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r SafetyAlert) RawJSON() string { return r.JSON.raw }
func (r *SafetyAlert) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type SafetyAlertErrorType string

const (
	SafetyAlertErrorTypePotentiallyUnintendedDataTransfer        SafetyAlertErrorType = "potentially_unintended_data_transfer"
	SafetyAlertErrorTypePotentiallyUnintendedDataAccess          SafetyAlertErrorType = "potentially_unintended_data_access"
	SafetyAlertErrorTypePotentiallyUnintendedDestructiveActivity SafetyAlertErrorType = "potentially_unintended_destructive_activity"
	SafetyAlertErrorTypeOther                                    SafetyAlertErrorType = "other"
)

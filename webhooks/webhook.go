// File generated from our OpenAPI spec by Castiron. See CONTRIBUTING.md for details.

package webhooks

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/openai/openai-go/v3/internal/apijson"
	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/respjson"
	"github.com/openai/openai-go/v3/shared/constant"
)

// WebhookService contains methods and other services that help with interacting
// with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewWebhookService] method instead.
type WebhookService struct {
	Options []option.RequestOption
}

// NewWebhookService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewWebhookService(opts ...option.RequestOption) (r WebhookService) {
	r = WebhookService{}
	r.Options = requestconfig.InheritedOptions(opts...)
	return
}

// Validates that the given payload was sent by OpenAI and parses the payload.
func (r *WebhookService) Unwrap(body []byte, headers http.Header, opts ...option.RequestOption) (*UnwrapWebhookEventUnion, error) {
	return unwrapWebhook(r, body, headers, opts...)
}

// UnwrapWithTolerance validates that the given payload was sent by OpenAI using custom tolerance, then parses the payload.
// tolerance specifies the maximum age of the webhook.
func (r *WebhookService) UnwrapWithTolerance(body []byte, headers http.Header, tolerance time.Duration, opts ...option.RequestOption) (*UnwrapWebhookEventUnion, error) {
	return unwrapWebhookWithTolerance(r, body, headers, tolerance, opts...)
}

// UnwrapWithToleranceAndTime validates that the given payload was sent by OpenAI using custom tolerance and time, then parses the payload.
// tolerance specifies the maximum age of the webhook.
// now allows specifying the current time for testing purposes.
func (r *WebhookService) UnwrapWithToleranceAndTime(body []byte, headers http.Header, tolerance time.Duration, now time.Time, opts ...option.RequestOption) (*UnwrapWebhookEventUnion, error) {
	return unwrapWebhookWithToleranceAndTime(r, body, headers, tolerance, now, opts...)
}

// VerifySignature validates whether or not the webhook payload was sent by OpenAI.
// An error will be raised if the webhook signature is invalid.
// tolerance specifies the maximum age of the webhook (default: 5 minutes).
func (r *WebhookService) VerifySignature(body []byte, headers http.Header, opts ...option.RequestOption) error {
	return r.VerifySignatureWithTolerance(body, headers, 5*time.Minute, opts...)
}

// VerifySignatureWithTolerance validates whether or not the webhook payload was sent by OpenAI.
// An error will be raised if the webhook signature is invalid.
// tolerance specifies the maximum age of the webhook.
func (r *WebhookService) VerifySignatureWithTolerance(body []byte, headers http.Header, tolerance time.Duration, opts ...option.RequestOption) error {
	return r.VerifySignatureWithToleranceAndTime(body, headers, tolerance, time.Now(), opts...)
}

// VerifySignatureWithToleranceAndTime validates whether or not the webhook payload was sent by OpenAI.
// An error will be raised if the webhook signature is invalid.
// tolerance specifies the maximum age of the webhook.
// now allows specifying the current time for testing purposes.
func (r *WebhookService) VerifySignatureWithToleranceAndTime(body []byte, headers http.Header, tolerance time.Duration, now time.Time, opts ...option.RequestOption) error {
	return verifyWebhookSignatureWithToleranceAndTime(r, body, headers, tolerance, now, opts...)
}

// Sent when a batch API request has been cancelled.
type BatchCancelledWebhookEvent struct {
	// The unique ID of the event.
	ID string `json:"id" api:"required"`
	// The Unix timestamp (in seconds) of when the batch API request was cancelled.
	CreatedAt int64 `json:"created_at" api:"required" format:"unixtime"`
	// Event data payload.
	Data BatchCancelledWebhookEventData `json:"data" api:"required"`
	// The type of the event. Always `batch.cancelled`.
	Type constant.BatchCancelled `json:"type" default:"batch.cancelled"`
	// The object of the event. Always `event`.
	//
	// Any of "event".
	Object BatchCancelledWebhookEventObject `json:"object"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		Data        respjson.Field
		Type        respjson.Field
		Object      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BatchCancelledWebhookEvent) RawJSON() string { return r.JSON.raw }
func (r *BatchCancelledWebhookEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Event data payload.
type BatchCancelledWebhookEventData struct {
	// The unique ID of the batch API request.
	ID string `json:"id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BatchCancelledWebhookEventData) RawJSON() string { return r.JSON.raw }
func (r *BatchCancelledWebhookEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The object of the event. Always `event`.
type BatchCancelledWebhookEventObject string

const (
	BatchCancelledWebhookEventObjectEvent BatchCancelledWebhookEventObject = "event"
)

// Sent when a batch API request has been completed.
type BatchCompletedWebhookEvent struct {
	// The unique ID of the event.
	ID string `json:"id" api:"required"`
	// The Unix timestamp (in seconds) of when the batch API request was completed.
	CreatedAt int64 `json:"created_at" api:"required" format:"unixtime"`
	// Event data payload.
	Data BatchCompletedWebhookEventData `json:"data" api:"required"`
	// The type of the event. Always `batch.completed`.
	Type constant.BatchCompleted `json:"type" default:"batch.completed"`
	// The object of the event. Always `event`.
	//
	// Any of "event".
	Object BatchCompletedWebhookEventObject `json:"object"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		Data        respjson.Field
		Type        respjson.Field
		Object      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BatchCompletedWebhookEvent) RawJSON() string { return r.JSON.raw }
func (r *BatchCompletedWebhookEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Event data payload.
type BatchCompletedWebhookEventData struct {
	// The unique ID of the batch API request.
	ID string `json:"id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BatchCompletedWebhookEventData) RawJSON() string { return r.JSON.raw }
func (r *BatchCompletedWebhookEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The object of the event. Always `event`.
type BatchCompletedWebhookEventObject string

const (
	BatchCompletedWebhookEventObjectEvent BatchCompletedWebhookEventObject = "event"
)

// Sent when a batch API request has expired.
type BatchExpiredWebhookEvent struct {
	// The unique ID of the event.
	ID string `json:"id" api:"required"`
	// The Unix timestamp (in seconds) of when the batch API request expired.
	CreatedAt int64 `json:"created_at" api:"required" format:"unixtime"`
	// Event data payload.
	Data BatchExpiredWebhookEventData `json:"data" api:"required"`
	// The type of the event. Always `batch.expired`.
	Type constant.BatchExpired `json:"type" default:"batch.expired"`
	// The object of the event. Always `event`.
	//
	// Any of "event".
	Object BatchExpiredWebhookEventObject `json:"object"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		Data        respjson.Field
		Type        respjson.Field
		Object      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BatchExpiredWebhookEvent) RawJSON() string { return r.JSON.raw }
func (r *BatchExpiredWebhookEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Event data payload.
type BatchExpiredWebhookEventData struct {
	// The unique ID of the batch API request.
	ID string `json:"id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BatchExpiredWebhookEventData) RawJSON() string { return r.JSON.raw }
func (r *BatchExpiredWebhookEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The object of the event. Always `event`.
type BatchExpiredWebhookEventObject string

const (
	BatchExpiredWebhookEventObjectEvent BatchExpiredWebhookEventObject = "event"
)

// Sent when a batch API request has failed.
type BatchFailedWebhookEvent struct {
	// The unique ID of the event.
	ID string `json:"id" api:"required"`
	// The Unix timestamp (in seconds) of when the batch API request failed.
	CreatedAt int64 `json:"created_at" api:"required" format:"unixtime"`
	// Event data payload.
	Data BatchFailedWebhookEventData `json:"data" api:"required"`
	// The type of the event. Always `batch.failed`.
	Type constant.BatchFailed `json:"type" default:"batch.failed"`
	// The object of the event. Always `event`.
	//
	// Any of "event".
	Object BatchFailedWebhookEventObject `json:"object"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		Data        respjson.Field
		Type        respjson.Field
		Object      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BatchFailedWebhookEvent) RawJSON() string { return r.JSON.raw }
func (r *BatchFailedWebhookEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Event data payload.
type BatchFailedWebhookEventData struct {
	// The unique ID of the batch API request.
	ID string `json:"id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BatchFailedWebhookEventData) RawJSON() string { return r.JSON.raw }
func (r *BatchFailedWebhookEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The object of the event. Always `event`.
type BatchFailedWebhookEventObject string

const (
	BatchFailedWebhookEventObjectEvent BatchFailedWebhookEventObject = "event"
)

// Sent when an eval run has been canceled.
type EvalRunCanceledWebhookEvent struct {
	// The unique ID of the event.
	ID string `json:"id" api:"required"`
	// The Unix timestamp (in seconds) of when the eval run was canceled.
	CreatedAt int64 `json:"created_at" api:"required" format:"unixtime"`
	// Event data payload.
	Data EvalRunCanceledWebhookEventData `json:"data" api:"required"`
	// The type of the event. Always `eval.run.canceled`.
	Type constant.EvalRunCanceled `json:"type" default:"eval.run.canceled"`
	// The object of the event. Always `event`.
	//
	// Any of "event".
	Object EvalRunCanceledWebhookEventObject `json:"object"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		Data        respjson.Field
		Type        respjson.Field
		Object      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EvalRunCanceledWebhookEvent) RawJSON() string { return r.JSON.raw }
func (r *EvalRunCanceledWebhookEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Event data payload.
type EvalRunCanceledWebhookEventData struct {
	// The unique ID of the eval run.
	ID string `json:"id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EvalRunCanceledWebhookEventData) RawJSON() string { return r.JSON.raw }
func (r *EvalRunCanceledWebhookEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The object of the event. Always `event`.
type EvalRunCanceledWebhookEventObject string

const (
	EvalRunCanceledWebhookEventObjectEvent EvalRunCanceledWebhookEventObject = "event"
)

// Sent when an eval run has failed.
type EvalRunFailedWebhookEvent struct {
	// The unique ID of the event.
	ID string `json:"id" api:"required"`
	// The Unix timestamp (in seconds) of when the eval run failed.
	CreatedAt int64 `json:"created_at" api:"required" format:"unixtime"`
	// Event data payload.
	Data EvalRunFailedWebhookEventData `json:"data" api:"required"`
	// The type of the event. Always `eval.run.failed`.
	Type constant.EvalRunFailed `json:"type" default:"eval.run.failed"`
	// The object of the event. Always `event`.
	//
	// Any of "event".
	Object EvalRunFailedWebhookEventObject `json:"object"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		Data        respjson.Field
		Type        respjson.Field
		Object      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EvalRunFailedWebhookEvent) RawJSON() string { return r.JSON.raw }
func (r *EvalRunFailedWebhookEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Event data payload.
type EvalRunFailedWebhookEventData struct {
	// The unique ID of the eval run.
	ID string `json:"id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EvalRunFailedWebhookEventData) RawJSON() string { return r.JSON.raw }
func (r *EvalRunFailedWebhookEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The object of the event. Always `event`.
type EvalRunFailedWebhookEventObject string

const (
	EvalRunFailedWebhookEventObjectEvent EvalRunFailedWebhookEventObject = "event"
)

// Sent when an eval run has succeeded.
type EvalRunSucceededWebhookEvent struct {
	// The unique ID of the event.
	ID string `json:"id" api:"required"`
	// The Unix timestamp (in seconds) of when the eval run succeeded.
	CreatedAt int64 `json:"created_at" api:"required" format:"unixtime"`
	// Event data payload.
	Data EvalRunSucceededWebhookEventData `json:"data" api:"required"`
	// The type of the event. Always `eval.run.succeeded`.
	Type constant.EvalRunSucceeded `json:"type" default:"eval.run.succeeded"`
	// The object of the event. Always `event`.
	//
	// Any of "event".
	Object EvalRunSucceededWebhookEventObject `json:"object"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		Data        respjson.Field
		Type        respjson.Field
		Object      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EvalRunSucceededWebhookEvent) RawJSON() string { return r.JSON.raw }
func (r *EvalRunSucceededWebhookEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Event data payload.
type EvalRunSucceededWebhookEventData struct {
	// The unique ID of the eval run.
	ID string `json:"id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EvalRunSucceededWebhookEventData) RawJSON() string { return r.JSON.raw }
func (r *EvalRunSucceededWebhookEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The object of the event. Always `event`.
type EvalRunSucceededWebhookEventObject string

const (
	EvalRunSucceededWebhookEventObjectEvent EvalRunSucceededWebhookEventObject = "event"
)

// Sent when a fine-tuning job has been cancelled.
type FineTuningJobCancelledWebhookEvent struct {
	// The unique ID of the event.
	ID string `json:"id" api:"required"`
	// The Unix timestamp (in seconds) of when the fine-tuning job was cancelled.
	CreatedAt int64 `json:"created_at" api:"required" format:"unixtime"`
	// Event data payload.
	Data FineTuningJobCancelledWebhookEventData `json:"data" api:"required"`
	// The type of the event. Always `fine_tuning.job.cancelled`.
	Type constant.FineTuningJobCancelled `json:"type" default:"fine_tuning.job.cancelled"`
	// The object of the event. Always `event`.
	//
	// Any of "event".
	Object FineTuningJobCancelledWebhookEventObject `json:"object"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		Data        respjson.Field
		Type        respjson.Field
		Object      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r FineTuningJobCancelledWebhookEvent) RawJSON() string { return r.JSON.raw }
func (r *FineTuningJobCancelledWebhookEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Event data payload.
type FineTuningJobCancelledWebhookEventData struct {
	// The unique ID of the fine-tuning job.
	ID string `json:"id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r FineTuningJobCancelledWebhookEventData) RawJSON() string { return r.JSON.raw }
func (r *FineTuningJobCancelledWebhookEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The object of the event. Always `event`.
type FineTuningJobCancelledWebhookEventObject string

const (
	FineTuningJobCancelledWebhookEventObjectEvent FineTuningJobCancelledWebhookEventObject = "event"
)

// Sent when a fine-tuning job has failed.
type FineTuningJobFailedWebhookEvent struct {
	// The unique ID of the event.
	ID string `json:"id" api:"required"`
	// The Unix timestamp (in seconds) of when the fine-tuning job failed.
	CreatedAt int64 `json:"created_at" api:"required" format:"unixtime"`
	// Event data payload.
	Data FineTuningJobFailedWebhookEventData `json:"data" api:"required"`
	// The type of the event. Always `fine_tuning.job.failed`.
	Type constant.FineTuningJobFailed `json:"type" default:"fine_tuning.job.failed"`
	// The object of the event. Always `event`.
	//
	// Any of "event".
	Object FineTuningJobFailedWebhookEventObject `json:"object"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		Data        respjson.Field
		Type        respjson.Field
		Object      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r FineTuningJobFailedWebhookEvent) RawJSON() string { return r.JSON.raw }
func (r *FineTuningJobFailedWebhookEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Event data payload.
type FineTuningJobFailedWebhookEventData struct {
	// The unique ID of the fine-tuning job.
	ID string `json:"id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r FineTuningJobFailedWebhookEventData) RawJSON() string { return r.JSON.raw }
func (r *FineTuningJobFailedWebhookEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The object of the event. Always `event`.
type FineTuningJobFailedWebhookEventObject string

const (
	FineTuningJobFailedWebhookEventObjectEvent FineTuningJobFailedWebhookEventObject = "event"
)

// Sent when a fine-tuning job has succeeded.
type FineTuningJobSucceededWebhookEvent struct {
	// The unique ID of the event.
	ID string `json:"id" api:"required"`
	// The Unix timestamp (in seconds) of when the fine-tuning job succeeded.
	CreatedAt int64 `json:"created_at" api:"required" format:"unixtime"`
	// Event data payload.
	Data FineTuningJobSucceededWebhookEventData `json:"data" api:"required"`
	// The type of the event. Always `fine_tuning.job.succeeded`.
	Type constant.FineTuningJobSucceeded `json:"type" default:"fine_tuning.job.succeeded"`
	// The object of the event. Always `event`.
	//
	// Any of "event".
	Object FineTuningJobSucceededWebhookEventObject `json:"object"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		Data        respjson.Field
		Type        respjson.Field
		Object      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r FineTuningJobSucceededWebhookEvent) RawJSON() string { return r.JSON.raw }
func (r *FineTuningJobSucceededWebhookEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Event data payload.
type FineTuningJobSucceededWebhookEventData struct {
	// The unique ID of the fine-tuning job.
	ID string `json:"id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r FineTuningJobSucceededWebhookEventData) RawJSON() string { return r.JSON.raw }
func (r *FineTuningJobSucceededWebhookEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The object of the event. Always `event`.
type FineTuningJobSucceededWebhookEventObject string

const (
	FineTuningJobSucceededWebhookEventObjectEvent FineTuningJobSucceededWebhookEventObject = "event"
)

// Sent when an incoming API SIP session is available for Live acceptance. The same
// pending session can also emit `realtime.call.incoming`; the first successful
// Realtime or Live accept endpoint selects the runtime surface.
type LiveCallIncomingWebhookEvent struct {
	// The unique ID of the event.
	ID string `json:"id" api:"required"`
	// The Unix timestamp (in seconds) of when the event was created.
	CreatedAt int64 `json:"created_at" api:"required" format:"unixtime"`
	// Event data payload.
	Data LiveCallIncomingWebhookEventData `json:"data" api:"required"`
	// The type of the event. Always `live.call.incoming`.
	Type constant.LiveCallIncoming `json:"type" default:"live.call.incoming"`
	// The object of the event. Always `event`.
	//
	// Any of "event".
	Object LiveCallIncomingWebhookEventObject `json:"object"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		Data        respjson.Field
		Type        respjson.Field
		Object      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r LiveCallIncomingWebhookEvent) RawJSON() string { return r.JSON.raw }
func (r *LiveCallIncomingWebhookEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Event data payload.
type LiveCallIncomingWebhookEventData struct {
	// The Transceiver `rtc_...` ID of the pending SIP session. The same value appears
	// as `call_id` in `realtime.call.incoming`.
	SessionID string `json:"session_id" api:"required"`
	// Headers from the SIP Invite.
	SipHeaders []LiveCallIncomingWebhookEventDataSipHeader `json:"sip_headers" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		SessionID   respjson.Field
		SipHeaders  respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r LiveCallIncomingWebhookEventData) RawJSON() string { return r.JSON.raw }
func (r *LiveCallIncomingWebhookEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A header from the SIP Invite.
type LiveCallIncomingWebhookEventDataSipHeader struct {
	// Name of the SIP Header.
	Name string `json:"name" api:"required"`
	// Value of the SIP Header.
	Value string `json:"value" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Name        respjson.Field
		Value       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r LiveCallIncomingWebhookEventDataSipHeader) RawJSON() string { return r.JSON.raw }
func (r *LiveCallIncomingWebhookEventDataSipHeader) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The object of the event. Always `event`.
type LiveCallIncomingWebhookEventObject string

const (
	LiveCallIncomingWebhookEventObjectEvent LiveCallIncomingWebhookEventObject = "event"
)

// Sent when an incoming API SIP session is available for Realtime acceptance. The
// same pending session can also emit `live.call.incoming`; the first successful
// Realtime or Live accept endpoint selects the runtime surface.
type RealtimeCallIncomingWebhookEvent struct {
	// The unique ID of the event.
	ID string `json:"id" api:"required"`
	// The Unix timestamp (in seconds) of when the model response was completed.
	CreatedAt int64 `json:"created_at" api:"required" format:"unixtime"`
	// Event data payload.
	Data RealtimeCallIncomingWebhookEventData `json:"data" api:"required"`
	// The type of the event. Always `realtime.call.incoming`.
	Type constant.RealtimeCallIncoming `json:"type" default:"realtime.call.incoming"`
	// The object of the event. Always `event`.
	//
	// Any of "event".
	Object RealtimeCallIncomingWebhookEventObject `json:"object"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		Data        respjson.Field
		Type        respjson.Field
		Object      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeCallIncomingWebhookEvent) RawJSON() string { return r.JSON.raw }
func (r *RealtimeCallIncomingWebhookEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Event data payload.
type RealtimeCallIncomingWebhookEventData struct {
	// The Transceiver `rtc_...` ID of the pending SIP session. The same value appears
	// as `session_id` in `live.call.incoming`.
	CallID string `json:"call_id" api:"required"`
	// Headers from the SIP Invite.
	SipHeaders []RealtimeCallIncomingWebhookEventDataSipHeader `json:"sip_headers" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		CallID      respjson.Field
		SipHeaders  respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeCallIncomingWebhookEventData) RawJSON() string { return r.JSON.raw }
func (r *RealtimeCallIncomingWebhookEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A header from the SIP Invite.
type RealtimeCallIncomingWebhookEventDataSipHeader struct {
	// Name of the SIP Header.
	Name string `json:"name" api:"required"`
	// Value of the SIP Header.
	Value string `json:"value" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Name        respjson.Field
		Value       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RealtimeCallIncomingWebhookEventDataSipHeader) RawJSON() string { return r.JSON.raw }
func (r *RealtimeCallIncomingWebhookEventDataSipHeader) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The object of the event. Always `event`.
type RealtimeCallIncomingWebhookEventObject string

const (
	RealtimeCallIncomingWebhookEventObjectEvent RealtimeCallIncomingWebhookEventObject = "event"
)

// Sent when a background response has been cancelled.
type ResponseCancelledWebhookEvent struct {
	// The unique ID of the event.
	ID string `json:"id" api:"required"`
	// The Unix timestamp (in seconds) of when the model response was cancelled.
	CreatedAt int64 `json:"created_at" api:"required" format:"unixtime"`
	// Event data payload.
	Data ResponseCancelledWebhookEventData `json:"data" api:"required"`
	// The type of the event. Always `response.cancelled`.
	Type constant.ResponseCancelled `json:"type" default:"response.cancelled"`
	// The object of the event. Always `event`.
	//
	// Any of "event".
	Object ResponseCancelledWebhookEventObject `json:"object"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		Data        respjson.Field
		Type        respjson.Field
		Object      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponseCancelledWebhookEvent) RawJSON() string { return r.JSON.raw }
func (r *ResponseCancelledWebhookEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Event data payload.
type ResponseCancelledWebhookEventData struct {
	// The unique ID of the model response.
	ID string `json:"id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponseCancelledWebhookEventData) RawJSON() string { return r.JSON.raw }
func (r *ResponseCancelledWebhookEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The object of the event. Always `event`.
type ResponseCancelledWebhookEventObject string

const (
	ResponseCancelledWebhookEventObjectEvent ResponseCancelledWebhookEventObject = "event"
)

// Sent when a background response has been completed.
type ResponseCompletedWebhookEvent struct {
	// The unique ID of the event.
	ID string `json:"id" api:"required"`
	// The Unix timestamp (in seconds) of when the model response was completed.
	CreatedAt int64 `json:"created_at" api:"required" format:"unixtime"`
	// Event data payload.
	Data ResponseCompletedWebhookEventData `json:"data" api:"required"`
	// The type of the event. Always `response.completed`.
	Type constant.ResponseCompleted `json:"type" default:"response.completed"`
	// The object of the event. Always `event`.
	//
	// Any of "event".
	Object ResponseCompletedWebhookEventObject `json:"object"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		Data        respjson.Field
		Type        respjson.Field
		Object      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponseCompletedWebhookEvent) RawJSON() string { return r.JSON.raw }
func (r *ResponseCompletedWebhookEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Event data payload.
type ResponseCompletedWebhookEventData struct {
	// The unique ID of the model response.
	ID string `json:"id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponseCompletedWebhookEventData) RawJSON() string { return r.JSON.raw }
func (r *ResponseCompletedWebhookEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The object of the event. Always `event`.
type ResponseCompletedWebhookEventObject string

const (
	ResponseCompletedWebhookEventObjectEvent ResponseCompletedWebhookEventObject = "event"
)

// Sent when a background response has failed.
type ResponseFailedWebhookEvent struct {
	// The unique ID of the event.
	ID string `json:"id" api:"required"`
	// The Unix timestamp (in seconds) of when the model response failed.
	CreatedAt int64 `json:"created_at" api:"required" format:"unixtime"`
	// Event data payload.
	Data ResponseFailedWebhookEventData `json:"data" api:"required"`
	// The type of the event. Always `response.failed`.
	Type constant.ResponseFailed `json:"type" default:"response.failed"`
	// The object of the event. Always `event`.
	//
	// Any of "event".
	Object ResponseFailedWebhookEventObject `json:"object"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		Data        respjson.Field
		Type        respjson.Field
		Object      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponseFailedWebhookEvent) RawJSON() string { return r.JSON.raw }
func (r *ResponseFailedWebhookEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Event data payload.
type ResponseFailedWebhookEventData struct {
	// The unique ID of the model response.
	ID string `json:"id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponseFailedWebhookEventData) RawJSON() string { return r.JSON.raw }
func (r *ResponseFailedWebhookEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The object of the event. Always `event`.
type ResponseFailedWebhookEventObject string

const (
	ResponseFailedWebhookEventObjectEvent ResponseFailedWebhookEventObject = "event"
)

// Sent when a background response has been interrupted.
type ResponseIncompleteWebhookEvent struct {
	// The unique ID of the event.
	ID string `json:"id" api:"required"`
	// The Unix timestamp (in seconds) of when the model response was interrupted.
	CreatedAt int64 `json:"created_at" api:"required" format:"unixtime"`
	// Event data payload.
	Data ResponseIncompleteWebhookEventData `json:"data" api:"required"`
	// The type of the event. Always `response.incomplete`.
	Type constant.ResponseIncomplete `json:"type" default:"response.incomplete"`
	// The object of the event. Always `event`.
	//
	// Any of "event".
	Object ResponseIncompleteWebhookEventObject `json:"object"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		Data        respjson.Field
		Type        respjson.Field
		Object      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponseIncompleteWebhookEvent) RawJSON() string { return r.JSON.raw }
func (r *ResponseIncompleteWebhookEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Event data payload.
type ResponseIncompleteWebhookEventData struct {
	// The unique ID of the model response.
	ID string `json:"id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ResponseIncompleteWebhookEventData) RawJSON() string { return r.JSON.raw }
func (r *ResponseIncompleteWebhookEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The object of the event. Always `event`.
type ResponseIncompleteWebhookEventObject string

const (
	ResponseIncompleteWebhookEventObjectEvent ResponseIncompleteWebhookEventObject = "event"
)

// Sent when an approved safety alert is available for an API project.
type SafetyAlertCreatedWebhookEvent struct {
	// The unique ID of the webhook event.
	ID string `json:"id" api:"required"`
	// The Unix timestamp in seconds when the event was created.
	CreatedAt int64                              `json:"created_at" api:"required" format:"unixtime"`
	Data      SafetyAlertCreatedWebhookEventData `json:"data" api:"required"`
	// Always `event`.
	Object constant.Event `json:"object" default:"event"`
	// Always `safety.alert.created`.
	Type constant.SafetyAlertCreated `json:"type" default:"safety.alert.created"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		Data        respjson.Field
		Object      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r SafetyAlertCreatedWebhookEvent) RawJSON() string { return r.JSON.raw }
func (r *SafetyAlertCreatedWebhookEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type SafetyAlertCreatedWebhookEventData struct {
	// The safety alert ID to pass to `GET /v1/safety/alerts/{id}`.
	ID string `json:"id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r SafetyAlertCreatedWebhookEventData) RawJSON() string { return r.JSON.raw }
func (r *SafetyAlertCreatedWebhookEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Sent when an approved safety alert is available for an enterprise workspace.
type SafetyOrgAlertCreatedWebhookEvent struct {
	// The unique ID of the webhook event.
	ID string `json:"id" api:"required"`
	// The Unix timestamp in seconds when the event was created.
	CreatedAt int64                                 `json:"created_at" api:"required" format:"unixtime"`
	Data      SafetyOrgAlertCreatedWebhookEventData `json:"data" api:"required"`
	// Always `event`.
	Object constant.Event `json:"object" default:"event"`
	// Always `safety.org_alert.created`.
	Type constant.SafetyOrgAlertCreated `json:"type" default:"safety.org_alert.created"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		Data        respjson.Field
		Object      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r SafetyOrgAlertCreatedWebhookEvent) RawJSON() string { return r.JSON.raw }
func (r *SafetyOrgAlertCreatedWebhookEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type SafetyOrgAlertCreatedWebhookEventData struct {
	// The safety alert ID to pass to `GET /v1/safety/alerts/{id}`.
	ID string `json:"id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r SafetyOrgAlertCreatedWebhookEventData) RawJSON() string { return r.JSON.raw }
func (r *SafetyOrgAlertCreatedWebhookEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// UnwrapWebhookEventUnion contains all possible properties and values from
// [BatchCancelledWebhookEvent], [BatchCompletedWebhookEvent],
// [BatchExpiredWebhookEvent], [BatchFailedWebhookEvent],
// [EvalRunCanceledWebhookEvent], [EvalRunFailedWebhookEvent],
// [EvalRunSucceededWebhookEvent], [FineTuningJobCancelledWebhookEvent],
// [FineTuningJobFailedWebhookEvent], [FineTuningJobSucceededWebhookEvent],
// [LiveCallIncomingWebhookEvent], [RealtimeCallIncomingWebhookEvent],
// [ResponseCancelledWebhookEvent], [ResponseCompletedWebhookEvent],
// [ResponseFailedWebhookEvent], [ResponseIncompleteWebhookEvent],
// [SafetyAlertCreatedWebhookEvent], [SafetyOrgAlertCreatedWebhookEvent].
//
// Use the [UnwrapWebhookEventUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type UnwrapWebhookEventUnion struct {
	ID        string `json:"id"`
	CreatedAt int64  `json:"created_at"`
	// This field is a union of [BatchCancelledWebhookEventData],
	// [BatchCompletedWebhookEventData], [BatchExpiredWebhookEventData],
	// [BatchFailedWebhookEventData], [EvalRunCanceledWebhookEventData],
	// [EvalRunFailedWebhookEventData], [EvalRunSucceededWebhookEventData],
	// [FineTuningJobCancelledWebhookEventData], [FineTuningJobFailedWebhookEventData],
	// [FineTuningJobSucceededWebhookEventData], [LiveCallIncomingWebhookEventData],
	// [RealtimeCallIncomingWebhookEventData], [ResponseCancelledWebhookEventData],
	// [ResponseCompletedWebhookEventData], [ResponseFailedWebhookEventData],
	// [ResponseIncompleteWebhookEventData], [SafetyAlertCreatedWebhookEventData],
	// [SafetyOrgAlertCreatedWebhookEventData]
	Data UnwrapWebhookEventUnionData `json:"data"`
	// Any of "batch.cancelled", "batch.completed", "batch.expired", "batch.failed",
	// "eval.run.canceled", "eval.run.failed", "eval.run.succeeded",
	// "fine_tuning.job.cancelled", "fine_tuning.job.failed",
	// "fine_tuning.job.succeeded", "live.call.incoming", "realtime.call.incoming",
	// "response.cancelled", "response.completed", "response.failed",
	// "response.incomplete", "safety.alert.created", "safety.org_alert.created".
	Type   string `json:"type"`
	Object string `json:"object"`
	JSON   struct {
		ID        respjson.Field
		CreatedAt respjson.Field
		Data      respjson.Field
		Type      respjson.Field
		Object    respjson.Field
		raw       string
	} `json:"-"`
}

// anyUnwrapWebhookEvent is implemented by each variant of
// [UnwrapWebhookEventUnion] to add type safety for the return type of
// [UnwrapWebhookEventUnion.AsAny]
type anyUnwrapWebhookEvent interface {
	implUnwrapWebhookEventUnion()
}

func (BatchCancelledWebhookEvent) implUnwrapWebhookEventUnion()         {}
func (BatchCompletedWebhookEvent) implUnwrapWebhookEventUnion()         {}
func (BatchExpiredWebhookEvent) implUnwrapWebhookEventUnion()           {}
func (BatchFailedWebhookEvent) implUnwrapWebhookEventUnion()            {}
func (EvalRunCanceledWebhookEvent) implUnwrapWebhookEventUnion()        {}
func (EvalRunFailedWebhookEvent) implUnwrapWebhookEventUnion()          {}
func (EvalRunSucceededWebhookEvent) implUnwrapWebhookEventUnion()       {}
func (FineTuningJobCancelledWebhookEvent) implUnwrapWebhookEventUnion() {}
func (FineTuningJobFailedWebhookEvent) implUnwrapWebhookEventUnion()    {}
func (FineTuningJobSucceededWebhookEvent) implUnwrapWebhookEventUnion() {}
func (LiveCallIncomingWebhookEvent) implUnwrapWebhookEventUnion()       {}
func (RealtimeCallIncomingWebhookEvent) implUnwrapWebhookEventUnion()   {}
func (ResponseCancelledWebhookEvent) implUnwrapWebhookEventUnion()      {}
func (ResponseCompletedWebhookEvent) implUnwrapWebhookEventUnion()      {}
func (ResponseFailedWebhookEvent) implUnwrapWebhookEventUnion()         {}
func (ResponseIncompleteWebhookEvent) implUnwrapWebhookEventUnion()     {}
func (SafetyAlertCreatedWebhookEvent) implUnwrapWebhookEventUnion()     {}
func (SafetyOrgAlertCreatedWebhookEvent) implUnwrapWebhookEventUnion()  {}

// Use the following switch statement to find the correct variant
//
//	switch variant := UnwrapWebhookEventUnion.AsAny().(type) {
//	case webhooks.BatchCancelledWebhookEvent:
//	case webhooks.BatchCompletedWebhookEvent:
//	case webhooks.BatchExpiredWebhookEvent:
//	case webhooks.BatchFailedWebhookEvent:
//	case webhooks.EvalRunCanceledWebhookEvent:
//	case webhooks.EvalRunFailedWebhookEvent:
//	case webhooks.EvalRunSucceededWebhookEvent:
//	case webhooks.FineTuningJobCancelledWebhookEvent:
//	case webhooks.FineTuningJobFailedWebhookEvent:
//	case webhooks.FineTuningJobSucceededWebhookEvent:
//	case webhooks.LiveCallIncomingWebhookEvent:
//	case webhooks.RealtimeCallIncomingWebhookEvent:
//	case webhooks.ResponseCancelledWebhookEvent:
//	case webhooks.ResponseCompletedWebhookEvent:
//	case webhooks.ResponseFailedWebhookEvent:
//	case webhooks.ResponseIncompleteWebhookEvent:
//	case webhooks.SafetyAlertCreatedWebhookEvent:
//	case webhooks.SafetyOrgAlertCreatedWebhookEvent:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u UnwrapWebhookEventUnion) AsAny() anyUnwrapWebhookEvent {
	switch u.Type {
	case "batch.cancelled":
		return u.AsBatchCancelled()
	case "batch.completed":
		return u.AsBatchCompleted()
	case "batch.expired":
		return u.AsBatchExpired()
	case "batch.failed":
		return u.AsBatchFailed()
	case "eval.run.canceled":
		return u.AsEvalRunCanceled()
	case "eval.run.failed":
		return u.AsEvalRunFailed()
	case "eval.run.succeeded":
		return u.AsEvalRunSucceeded()
	case "fine_tuning.job.cancelled":
		return u.AsFineTuningJobCancelled()
	case "fine_tuning.job.failed":
		return u.AsFineTuningJobFailed()
	case "fine_tuning.job.succeeded":
		return u.AsFineTuningJobSucceeded()
	case "live.call.incoming":
		return u.AsLiveCallIncoming()
	case "realtime.call.incoming":
		return u.AsRealtimeCallIncoming()
	case "response.cancelled":
		return u.AsResponseCancelled()
	case "response.completed":
		return u.AsResponseCompleted()
	case "response.failed":
		return u.AsResponseFailed()
	case "response.incomplete":
		return u.AsResponseIncomplete()
	case "safety.alert.created":
		return u.AsSafetyAlertCreated()
	case "safety.org_alert.created":
		return u.AsSafetyOrgAlertCreated()
	}
	return nil
}

func (u UnwrapWebhookEventUnion) AsBatchCancelled() (v BatchCancelledWebhookEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u UnwrapWebhookEventUnion) AsBatchCompleted() (v BatchCompletedWebhookEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u UnwrapWebhookEventUnion) AsBatchExpired() (v BatchExpiredWebhookEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u UnwrapWebhookEventUnion) AsBatchFailed() (v BatchFailedWebhookEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u UnwrapWebhookEventUnion) AsEvalRunCanceled() (v EvalRunCanceledWebhookEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u UnwrapWebhookEventUnion) AsEvalRunFailed() (v EvalRunFailedWebhookEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u UnwrapWebhookEventUnion) AsEvalRunSucceeded() (v EvalRunSucceededWebhookEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u UnwrapWebhookEventUnion) AsFineTuningJobCancelled() (v FineTuningJobCancelledWebhookEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u UnwrapWebhookEventUnion) AsFineTuningJobFailed() (v FineTuningJobFailedWebhookEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u UnwrapWebhookEventUnion) AsFineTuningJobSucceeded() (v FineTuningJobSucceededWebhookEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u UnwrapWebhookEventUnion) AsLiveCallIncoming() (v LiveCallIncomingWebhookEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u UnwrapWebhookEventUnion) AsRealtimeCallIncoming() (v RealtimeCallIncomingWebhookEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u UnwrapWebhookEventUnion) AsResponseCancelled() (v ResponseCancelledWebhookEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u UnwrapWebhookEventUnion) AsResponseCompleted() (v ResponseCompletedWebhookEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u UnwrapWebhookEventUnion) AsResponseFailed() (v ResponseFailedWebhookEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u UnwrapWebhookEventUnion) AsResponseIncomplete() (v ResponseIncompleteWebhookEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u UnwrapWebhookEventUnion) AsSafetyAlertCreated() (v SafetyAlertCreatedWebhookEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u UnwrapWebhookEventUnion) AsSafetyOrgAlertCreated() (v SafetyOrgAlertCreatedWebhookEvent) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u UnwrapWebhookEventUnion) RawJSON() string { return u.JSON.raw }

func (r *UnwrapWebhookEventUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// UnwrapWebhookEventUnionData is an implicit subunion of
// [UnwrapWebhookEventUnion]. UnwrapWebhookEventUnionData provides convenient
// access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [UnwrapWebhookEventUnion].
type UnwrapWebhookEventUnionData struct {
	ID string `json:"id"`
	// This field is from variant [LiveCallIncomingWebhookEventData].
	SessionID string `json:"session_id"`
	// This field is a union of [[]LiveCallIncomingWebhookEventDataSipHeader],
	// [[]RealtimeCallIncomingWebhookEventDataSipHeader]
	SipHeaders UnwrapWebhookEventUnionDataSipHeaders `json:"sip_headers"`
	// This field is from variant [RealtimeCallIncomingWebhookEventData].
	CallID string `json:"call_id"`
	JSON   struct {
		ID         respjson.Field
		SessionID  respjson.Field
		SipHeaders respjson.Field
		CallID     respjson.Field
		raw        string
	} `json:"-"`
}

func (r *UnwrapWebhookEventUnionData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// UnwrapWebhookEventUnionDataSipHeaders is an implicit subunion of
// [UnwrapWebhookEventUnion]. UnwrapWebhookEventUnionDataSipHeaders provides
// convenient access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [UnwrapWebhookEventUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfLiveCallIncomingWebhookEventDataSipHeaders
// OfRealtimeCallIncomingWebhookEventDataSipHeaders]
type UnwrapWebhookEventUnionDataSipHeaders struct {
	// This field will be present if the value is a
	// [[]LiveCallIncomingWebhookEventDataSipHeader] instead of an object.
	OfLiveCallIncomingWebhookEventDataSipHeaders []LiveCallIncomingWebhookEventDataSipHeader `json:",inline"`
	// This field will be present if the value is a
	// [[]RealtimeCallIncomingWebhookEventDataSipHeader] instead of an object.
	OfRealtimeCallIncomingWebhookEventDataSipHeaders []RealtimeCallIncomingWebhookEventDataSipHeader `json:",inline"`
	JSON                                             struct {
		OfLiveCallIncomingWebhookEventDataSipHeaders     respjson.Field
		OfRealtimeCallIncomingWebhookEventDataSipHeaders respjson.Field
		raw                                              string
	} `json:"-"`
}

func (r *UnwrapWebhookEventUnionDataSipHeaders) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

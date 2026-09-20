// File generated from our OpenAPI spec by Castiron. See CONTRIBUTING.md for details.

package webhooks

import (
	"context"
	"net/http"
	"slices"

	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
)

// EventTypeService contains methods and other services that help with interacting
// with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewEventTypeService] method instead.
type EventTypeService struct {
	Options []option.RequestOption
}

// NewEventTypeService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewEventTypeService(opts ...option.RequestOption) (r EventTypeService) {
	r = EventTypeService{}
	r.Options = requestconfig.InheritedOptions(opts...)
	return
}

// Returns webhook event types visible to the authenticated project.
func (r *EventTypeService) List(ctx context.Context, opts ...option.RequestOption) (res *WebhookEventTypeList, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	path := "webhook_event_types"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

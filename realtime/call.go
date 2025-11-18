// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package realtime

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/openai/openai-go/v3/internal/apijson"
	shimjson "github.com/openai/openai-go/v3/internal/encoding/json"
	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/param"
)

// CallService contains methods and other services that help with interacting with
// the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewCallService] method instead.
type CallService struct {
	Options []option.RequestOption
}

// NewCallService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewCallService(opts ...option.RequestOption) (r CallService) {
	r = CallService{}
	r.Options = opts
	return
}

// Accept an incoming SIP call and configure the realtime session that will handle
// it.
func (r *CallService) Accept(ctx context.Context, callID string, body CallAcceptParams, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if callID == "" {
		err = errors.New("missing required call_id parameter")
		return
	}
	path := fmt.Sprintf("realtime/calls/%s/accept", callID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, nil, opts...)
	return
}

// End an active Realtime API call, whether it was initiated over SIP or WebRTC.
func (r *CallService) Hangup(ctx context.Context, callID string, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if callID == "" {
		err = errors.New("missing required call_id parameter")
		return
	}
	path := fmt.Sprintf("realtime/calls/%s/hangup", callID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, nil, opts...)
	return
}

// Transfer an active SIP call to a new destination using the SIP REFER verb.
func (r *CallService) Refer(ctx context.Context, callID string, body CallReferParams, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if callID == "" {
		err = errors.New("missing required call_id parameter")
		return
	}
	path := fmt.Sprintf("realtime/calls/%s/refer", callID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, nil, opts...)
	return
}

// Decline an incoming SIP call by returning a SIP status code to the caller.
func (r *CallService) Reject(ctx context.Context, callID string, body CallRejectParams, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if callID == "" {
		err = errors.New("missing required call_id parameter")
		return
	}
	path := fmt.Sprintf("realtime/calls/%s/reject", callID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, nil, opts...)
	return
}

type CallAcceptParams struct {
	// Realtime session object configuration.
	RealtimeSessionCreateRequest RealtimeSessionCreateRequestParam
	paramObj
}

func (r CallAcceptParams) MarshalJSON() (data []byte, err error) {
	return shimjson.Marshal(r.RealtimeSessionCreateRequest)
}
func (r *CallAcceptParams) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &r.RealtimeSessionCreateRequest)
}

type CallReferParams struct {
	// URI that should appear in the SIP Refer-To header. Supports values like
	// `tel:+14155550123` or `sip:agent@example.com`.
	TargetUri string `json:"target_uri,required"`
	paramObj
}

func (r CallReferParams) MarshalJSON() (data []byte, err error) {
	type shadow CallReferParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *CallReferParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type CallRejectParams struct {
	// SIP response code to send back to the caller. Defaults to `603` (Decline) when
	// omitted.
	StatusCode param.Opt[int64] `json:"status_code,omitzero"`
	paramObj
}

func (r CallRejectParams) MarshalJSON() (data []byte, err error) {
	type shadow CallRejectParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *CallRejectParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

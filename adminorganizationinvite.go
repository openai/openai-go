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

// AdminOrganizationInviteService contains methods and other services that help
// with interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAdminOrganizationInviteService] method instead.
type AdminOrganizationInviteService struct {
	Options []option.RequestOption
}

// NewAdminOrganizationInviteService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewAdminOrganizationInviteService(opts ...option.RequestOption) (r AdminOrganizationInviteService) {
	r = AdminOrganizationInviteService{}
	r.Options = opts
	return
}

// Create an invite for a user to the organization. The invite must be accepted by
// the user before they have access to the organization.
func (r *AdminOrganizationInviteService) New(ctx context.Context, body AdminOrganizationInviteNewParams, opts ...option.RequestOption) (res *Invite, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	path := "organization/invites"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Retrieves an invite.
func (r *AdminOrganizationInviteService) Get(ctx context.Context, inviteID string, opts ...option.RequestOption) (res *Invite, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if inviteID == "" {
		err = errors.New("missing required invite_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("organization/invites/%s", inviteID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Returns a list of invites in the organization.
func (r *AdminOrganizationInviteService) List(ctx context.Context, query AdminOrganizationInviteListParams, opts ...option.RequestOption) (res *pagination.ConversationCursorPage[Invite], err error) {
	var raw *http.Response
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "organization/invites"
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

// Returns a list of invites in the organization.
func (r *AdminOrganizationInviteService) ListAutoPaging(ctx context.Context, query AdminOrganizationInviteListParams, opts ...option.RequestOption) *pagination.ConversationCursorPageAutoPager[Invite] {
	return pagination.NewConversationCursorPageAutoPager(r.List(ctx, query, opts...))
}

// Delete an invite. If the invite has already been accepted, it cannot be deleted.
func (r *AdminOrganizationInviteService) Delete(ctx context.Context, inviteID string, opts ...option.RequestOption) (res *AdminOrganizationInviteDeleteResponse, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	if inviteID == "" {
		err = errors.New("missing required invite_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("organization/invites/%s", inviteID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Represents an individual `invite` to the organization.
type Invite struct {
	// The identifier, which can be referenced in API endpoints
	ID string `json:"id" api:"required"`
	// The Unix timestamp (in seconds) of when the invite was sent.
	CreatedAt int64 `json:"created_at" api:"required" format:"unixtime"`
	// The email address of the individual to whom the invite was sent
	Email string `json:"email" api:"required"`
	// The Unix timestamp (in seconds) of when the invite expires.
	ExpiresAt int64 `json:"expires_at" api:"required" format:"unixtime"`
	// The object type, which is always `organization.invite`
	Object constant.OrganizationInvite `json:"object" default:"organization.invite"`
	// `owner` or `reader`
	//
	// Any of "owner", "reader".
	Role InviteRole `json:"role" api:"required"`
	// `accepted`,`expired`, or `pending`
	//
	// Any of "accepted", "expired", "pending".
	Status InviteStatus `json:"status" api:"required"`
	// The Unix timestamp (in seconds) of when the invite was accepted.
	AcceptedAt int64 `json:"accepted_at" api:"nullable" format:"unixtime"`
	// The projects that were granted membership upon acceptance of the invite.
	Projects []InviteProject `json:"projects"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		Email       respjson.Field
		ExpiresAt   respjson.Field
		Object      respjson.Field
		Role        respjson.Field
		Status      respjson.Field
		AcceptedAt  respjson.Field
		Projects    respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r Invite) RawJSON() string { return r.JSON.raw }
func (r *Invite) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// `owner` or `reader`
type InviteRole string

const (
	InviteRoleOwner  InviteRole = "owner"
	InviteRoleReader InviteRole = "reader"
)

// `accepted`,`expired`, or `pending`
type InviteStatus string

const (
	InviteStatusAccepted InviteStatus = "accepted"
	InviteStatusExpired  InviteStatus = "expired"
	InviteStatusPending  InviteStatus = "pending"
)

type InviteProject struct {
	// Project's public ID
	ID string `json:"id"`
	// Project membership role
	//
	// Any of "member", "owner".
	Role string `json:"role"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Role        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r InviteProject) RawJSON() string { return r.JSON.raw }
func (r *InviteProject) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationInviteDeleteResponse struct {
	ID      string `json:"id" api:"required"`
	Deleted bool   `json:"deleted" api:"required"`
	// The object type, which is always `organization.invite.deleted`
	Object constant.OrganizationInviteDeleted `json:"object" default:"organization.invite.deleted"`
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
func (r AdminOrganizationInviteDeleteResponse) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationInviteDeleteResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationInviteNewParams struct {
	// Send an email to this address
	Email string `json:"email" api:"required"`
	// `owner` or `reader`
	//
	// Any of "reader", "owner".
	Role AdminOrganizationInviteNewParamsRole `json:"role,omitzero" api:"required"`
	// An array of projects to which membership is granted at the same time the org
	// invite is accepted. If omitted, the user will be invited to the default project
	// for compatibility with legacy behavior.
	Projects []AdminOrganizationInviteNewParamsProject `json:"projects,omitzero"`
	paramObj
}

func (r AdminOrganizationInviteNewParams) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationInviteNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationInviteNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// `owner` or `reader`
type AdminOrganizationInviteNewParamsRole string

const (
	AdminOrganizationInviteNewParamsRoleReader AdminOrganizationInviteNewParamsRole = "reader"
	AdminOrganizationInviteNewParamsRoleOwner  AdminOrganizationInviteNewParamsRole = "owner"
)

// The properties ID, Role are required.
type AdminOrganizationInviteNewParamsProject struct {
	// Project's public ID
	ID string `json:"id" api:"required"`
	// Project membership role
	//
	// Any of "member", "owner".
	Role string `json:"role,omitzero" api:"required"`
	paramObj
}

func (r AdminOrganizationInviteNewParamsProject) MarshalJSON() (data []byte, err error) {
	type shadow AdminOrganizationInviteNewParamsProject
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdminOrganizationInviteNewParamsProject) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[AdminOrganizationInviteNewParamsProject](
		"role", "member", "owner",
	)
}

type AdminOrganizationInviteListParams struct {
	// A cursor for use in pagination. `after` is an object ID that defines your place
	// in the list. For instance, if you make a list request and receive 100 objects,
	// ending with obj_foo, your subsequent call can include after=obj_foo in order to
	// fetch the next page of the list.
	After param.Opt[string] `query:"after,omitzero" json:"-"`
	// A limit on the number of objects to be returned. Limit can range between 1 and
	// 100, and the default is 20.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [AdminOrganizationInviteListParams]'s query parameters as
// `url.Values`.
func (r AdminOrganizationInviteListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

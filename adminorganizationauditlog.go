// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai

import (
	"context"
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
)

// List user actions and configuration changes within this organization.
//
// AdminOrganizationAuditLogService contains methods and other services that help
// with interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAdminOrganizationAuditLogService] method instead.
type AdminOrganizationAuditLogService struct {
	Options []option.RequestOption
}

// NewAdminOrganizationAuditLogService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewAdminOrganizationAuditLogService(opts ...option.RequestOption) (r AdminOrganizationAuditLogService) {
	r = AdminOrganizationAuditLogService{}
	r.Options = opts
	return
}

// List user actions and configuration changes within this organization.
func (r *AdminOrganizationAuditLogService) List(ctx context.Context, query AdminOrganizationAuditLogListParams, opts ...option.RequestOption) (res *pagination.ConversationCursorPage[AdminOrganizationAuditLogListResponse], err error) {
	var raw *http.Response
	var preClientOpts = []option.RequestOption{requestconfig.WithAdminAPIKeyAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "organization/audit_logs"
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

// List user actions and configuration changes within this organization.
func (r *AdminOrganizationAuditLogService) ListAutoPaging(ctx context.Context, query AdminOrganizationAuditLogListParams, opts ...option.RequestOption) *pagination.ConversationCursorPageAutoPager[AdminOrganizationAuditLogListResponse] {
	return pagination.NewConversationCursorPageAutoPager(r.List(ctx, query, opts...))
}

// A log of a user action or configuration change within this organization.
type AdminOrganizationAuditLogListResponse struct {
	// The ID of this log.
	ID string `json:"id" api:"required"`
	// The Unix timestamp (in seconds) of the event.
	EffectiveAt int64 `json:"effective_at" api:"required" format:"unixtime"`
	// The event type.
	//
	// Any of "api_key.created", "api_key.updated", "api_key.deleted",
	// "certificate.created", "certificate.updated", "certificate.deleted",
	// "certificates.activated", "certificates.deactivated",
	// "checkpoint.permission.created", "checkpoint.permission.deleted",
	// "external_key.registered", "external_key.removed", "group.created",
	// "group.updated", "group.deleted", "invite.sent", "invite.accepted",
	// "invite.deleted", "ip_allowlist.created", "ip_allowlist.updated",
	// "ip_allowlist.deleted", "ip_allowlist.config.activated",
	// "ip_allowlist.config.deactivated", "login.succeeded", "login.failed",
	// "logout.succeeded", "logout.failed", "organization.updated", "project.created",
	// "project.updated", "project.archived", "project.deleted", "rate_limit.updated",
	// "rate_limit.deleted", "resource.deleted", "tunnel.created", "tunnel.updated",
	// "tunnel.deleted", "role.created", "role.updated", "role.deleted",
	// "role.assignment.created", "role.assignment.deleted", "scim.enabled",
	// "scim.disabled", "service_account.created", "service_account.updated",
	// "service_account.deleted", "user.added", "user.updated", "user.deleted".
	Type AdminOrganizationAuditLogListResponseType `json:"type" api:"required"`
	// The actor who performed the audit logged action.
	Actor AdminOrganizationAuditLogListResponseActor `json:"actor" api:"nullable"`
	// The details for events with this `type`.
	APIKeyCreated AdminOrganizationAuditLogListResponseAPIKeyCreated `json:"api_key.created"`
	// The details for events with this `type`.
	APIKeyDeleted AdminOrganizationAuditLogListResponseAPIKeyDeleted `json:"api_key.deleted"`
	// The details for events with this `type`.
	APIKeyUpdated AdminOrganizationAuditLogListResponseAPIKeyUpdated `json:"api_key.updated"`
	// The details for events with this `type`.
	CertificateCreated AdminOrganizationAuditLogListResponseCertificateCreated `json:"certificate.created"`
	// The details for events with this `type`.
	CertificateDeleted AdminOrganizationAuditLogListResponseCertificateDeleted `json:"certificate.deleted"`
	// The details for events with this `type`.
	CertificateUpdated AdminOrganizationAuditLogListResponseCertificateUpdated `json:"certificate.updated"`
	// The details for events with this `type`.
	CertificatesActivated AdminOrganizationAuditLogListResponseCertificatesActivated `json:"certificates.activated"`
	// The details for events with this `type`.
	CertificatesDeactivated AdminOrganizationAuditLogListResponseCertificatesDeactivated `json:"certificates.deactivated"`
	// The project and fine-tuned model checkpoint that the checkpoint permission was
	// created for.
	CheckpointPermissionCreated AdminOrganizationAuditLogListResponseCheckpointPermissionCreated `json:"checkpoint.permission.created"`
	// The details for events with this `type`.
	CheckpointPermissionDeleted AdminOrganizationAuditLogListResponseCheckpointPermissionDeleted `json:"checkpoint.permission.deleted"`
	// The details for events with this `type`.
	ExternalKeyRegistered AdminOrganizationAuditLogListResponseExternalKeyRegistered `json:"external_key.registered"`
	// The details for events with this `type`.
	ExternalKeyRemoved AdminOrganizationAuditLogListResponseExternalKeyRemoved `json:"external_key.removed"`
	// The details for events with this `type`.
	GroupCreated AdminOrganizationAuditLogListResponseGroupCreated `json:"group.created"`
	// The details for events with this `type`.
	GroupDeleted AdminOrganizationAuditLogListResponseGroupDeleted `json:"group.deleted"`
	// The details for events with this `type`.
	GroupUpdated AdminOrganizationAuditLogListResponseGroupUpdated `json:"group.updated"`
	// The details for events with this `type`.
	InviteAccepted AdminOrganizationAuditLogListResponseInviteAccepted `json:"invite.accepted"`
	// The details for events with this `type`.
	InviteDeleted AdminOrganizationAuditLogListResponseInviteDeleted `json:"invite.deleted"`
	// The details for events with this `type`.
	InviteSent AdminOrganizationAuditLogListResponseInviteSent `json:"invite.sent"`
	// The details for events with this `type`.
	IPAllowlistConfigActivated AdminOrganizationAuditLogListResponseIPAllowlistConfigActivated `json:"ip_allowlist.config.activated"`
	// The details for events with this `type`.
	IPAllowlistConfigDeactivated AdminOrganizationAuditLogListResponseIPAllowlistConfigDeactivated `json:"ip_allowlist.config.deactivated"`
	// The details for events with this `type`.
	IPAllowlistCreated AdminOrganizationAuditLogListResponseIPAllowlistCreated `json:"ip_allowlist.created"`
	// The details for events with this `type`.
	IPAllowlistDeleted AdminOrganizationAuditLogListResponseIPAllowlistDeleted `json:"ip_allowlist.deleted"`
	// The details for events with this `type`.
	IPAllowlistUpdated AdminOrganizationAuditLogListResponseIPAllowlistUpdated `json:"ip_allowlist.updated"`
	// The details for events with this `type`.
	LoginFailed AdminOrganizationAuditLogListResponseLoginFailed `json:"login.failed"`
	// This event has no additional fields beyond the standard audit log attributes.
	LoginSucceeded any `json:"login.succeeded"`
	// The details for events with this `type`.
	LogoutFailed AdminOrganizationAuditLogListResponseLogoutFailed `json:"logout.failed"`
	// This event has no additional fields beyond the standard audit log attributes.
	LogoutSucceeded any `json:"logout.succeeded"`
	// The details for events with this `type`.
	OrganizationUpdated AdminOrganizationAuditLogListResponseOrganizationUpdated `json:"organization.updated"`
	// The project that the action was scoped to. Absent for actions not scoped to
	// projects. Note that any admin actions taken via Admin API keys are associated
	// with the default project.
	Project AdminOrganizationAuditLogListResponseProject `json:"project"`
	// The details for events with this `type`.
	ProjectArchived AdminOrganizationAuditLogListResponseProjectArchived `json:"project.archived"`
	// The details for events with this `type`.
	ProjectCreated AdminOrganizationAuditLogListResponseProjectCreated `json:"project.created"`
	// The details for events with this `type`.
	ProjectDeleted AdminOrganizationAuditLogListResponseProjectDeleted `json:"project.deleted"`
	// The details for events with this `type`.
	ProjectUpdated AdminOrganizationAuditLogListResponseProjectUpdated `json:"project.updated"`
	// The details for events with this `type`.
	RateLimitDeleted AdminOrganizationAuditLogListResponseRateLimitDeleted `json:"rate_limit.deleted"`
	// The details for events with this `type`.
	RateLimitUpdated AdminOrganizationAuditLogListResponseRateLimitUpdated `json:"rate_limit.updated"`
	// The details for events with this `type`.
	RoleAssignmentCreated AdminOrganizationAuditLogListResponseRoleAssignmentCreated `json:"role.assignment.created"`
	// The details for events with this `type`.
	RoleAssignmentDeleted AdminOrganizationAuditLogListResponseRoleAssignmentDeleted `json:"role.assignment.deleted"`
	// The details for events with this `type`.
	RoleCreated AdminOrganizationAuditLogListResponseRoleCreated `json:"role.created"`
	// The details for events with this `type`.
	RoleDeleted AdminOrganizationAuditLogListResponseRoleDeleted `json:"role.deleted"`
	// The details for events with this `type`.
	RoleUpdated AdminOrganizationAuditLogListResponseRoleUpdated `json:"role.updated"`
	// The details for events with this `type`.
	ScimDisabled AdminOrganizationAuditLogListResponseScimDisabled `json:"scim.disabled"`
	// The details for events with this `type`.
	ScimEnabled AdminOrganizationAuditLogListResponseScimEnabled `json:"scim.enabled"`
	// The details for events with this `type`.
	ServiceAccountCreated AdminOrganizationAuditLogListResponseServiceAccountCreated `json:"service_account.created"`
	// The details for events with this `type`.
	ServiceAccountDeleted AdminOrganizationAuditLogListResponseServiceAccountDeleted `json:"service_account.deleted"`
	// The details for events with this `type`.
	ServiceAccountUpdated AdminOrganizationAuditLogListResponseServiceAccountUpdated `json:"service_account.updated"`
	// The details for events with this `type`.
	UserAdded AdminOrganizationAuditLogListResponseUserAdded `json:"user.added"`
	// The details for events with this `type`.
	UserDeleted AdminOrganizationAuditLogListResponseUserDeleted `json:"user.deleted"`
	// The details for events with this `type`.
	UserUpdated AdminOrganizationAuditLogListResponseUserUpdated `json:"user.updated"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                           respjson.Field
		EffectiveAt                  respjson.Field
		Type                         respjson.Field
		Actor                        respjson.Field
		APIKeyCreated                respjson.Field
		APIKeyDeleted                respjson.Field
		APIKeyUpdated                respjson.Field
		CertificateCreated           respjson.Field
		CertificateDeleted           respjson.Field
		CertificateUpdated           respjson.Field
		CertificatesActivated        respjson.Field
		CertificatesDeactivated      respjson.Field
		CheckpointPermissionCreated  respjson.Field
		CheckpointPermissionDeleted  respjson.Field
		ExternalKeyRegistered        respjson.Field
		ExternalKeyRemoved           respjson.Field
		GroupCreated                 respjson.Field
		GroupDeleted                 respjson.Field
		GroupUpdated                 respjson.Field
		InviteAccepted               respjson.Field
		InviteDeleted                respjson.Field
		InviteSent                   respjson.Field
		IPAllowlistConfigActivated   respjson.Field
		IPAllowlistConfigDeactivated respjson.Field
		IPAllowlistCreated           respjson.Field
		IPAllowlistDeleted           respjson.Field
		IPAllowlistUpdated           respjson.Field
		LoginFailed                  respjson.Field
		LoginSucceeded               respjson.Field
		LogoutFailed                 respjson.Field
		LogoutSucceeded              respjson.Field
		OrganizationUpdated          respjson.Field
		Project                      respjson.Field
		ProjectArchived              respjson.Field
		ProjectCreated               respjson.Field
		ProjectDeleted               respjson.Field
		ProjectUpdated               respjson.Field
		RateLimitDeleted             respjson.Field
		RateLimitUpdated             respjson.Field
		RoleAssignmentCreated        respjson.Field
		RoleAssignmentDeleted        respjson.Field
		RoleCreated                  respjson.Field
		RoleDeleted                  respjson.Field
		RoleUpdated                  respjson.Field
		ScimDisabled                 respjson.Field
		ScimEnabled                  respjson.Field
		ServiceAccountCreated        respjson.Field
		ServiceAccountDeleted        respjson.Field
		ServiceAccountUpdated        respjson.Field
		UserAdded                    respjson.Field
		UserDeleted                  respjson.Field
		UserUpdated                  respjson.Field
		ExtraFields                  map[string]respjson.Field
		raw                          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponse) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationAuditLogListResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The event type.
type AdminOrganizationAuditLogListResponseType string

const (
	AdminOrganizationAuditLogListResponseTypeAPIKeyCreated                AdminOrganizationAuditLogListResponseType = "api_key.created"
	AdminOrganizationAuditLogListResponseTypeAPIKeyUpdated                AdminOrganizationAuditLogListResponseType = "api_key.updated"
	AdminOrganizationAuditLogListResponseTypeAPIKeyDeleted                AdminOrganizationAuditLogListResponseType = "api_key.deleted"
	AdminOrganizationAuditLogListResponseTypeCertificateCreated           AdminOrganizationAuditLogListResponseType = "certificate.created"
	AdminOrganizationAuditLogListResponseTypeCertificateUpdated           AdminOrganizationAuditLogListResponseType = "certificate.updated"
	AdminOrganizationAuditLogListResponseTypeCertificateDeleted           AdminOrganizationAuditLogListResponseType = "certificate.deleted"
	AdminOrganizationAuditLogListResponseTypeCertificatesActivated        AdminOrganizationAuditLogListResponseType = "certificates.activated"
	AdminOrganizationAuditLogListResponseTypeCertificatesDeactivated      AdminOrganizationAuditLogListResponseType = "certificates.deactivated"
	AdminOrganizationAuditLogListResponseTypeCheckpointPermissionCreated  AdminOrganizationAuditLogListResponseType = "checkpoint.permission.created"
	AdminOrganizationAuditLogListResponseTypeCheckpointPermissionDeleted  AdminOrganizationAuditLogListResponseType = "checkpoint.permission.deleted"
	AdminOrganizationAuditLogListResponseTypeExternalKeyRegistered        AdminOrganizationAuditLogListResponseType = "external_key.registered"
	AdminOrganizationAuditLogListResponseTypeExternalKeyRemoved           AdminOrganizationAuditLogListResponseType = "external_key.removed"
	AdminOrganizationAuditLogListResponseTypeGroupCreated                 AdminOrganizationAuditLogListResponseType = "group.created"
	AdminOrganizationAuditLogListResponseTypeGroupUpdated                 AdminOrganizationAuditLogListResponseType = "group.updated"
	AdminOrganizationAuditLogListResponseTypeGroupDeleted                 AdminOrganizationAuditLogListResponseType = "group.deleted"
	AdminOrganizationAuditLogListResponseTypeInviteSent                   AdminOrganizationAuditLogListResponseType = "invite.sent"
	AdminOrganizationAuditLogListResponseTypeInviteAccepted               AdminOrganizationAuditLogListResponseType = "invite.accepted"
	AdminOrganizationAuditLogListResponseTypeInviteDeleted                AdminOrganizationAuditLogListResponseType = "invite.deleted"
	AdminOrganizationAuditLogListResponseTypeIPAllowlistCreated           AdminOrganizationAuditLogListResponseType = "ip_allowlist.created"
	AdminOrganizationAuditLogListResponseTypeIPAllowlistUpdated           AdminOrganizationAuditLogListResponseType = "ip_allowlist.updated"
	AdminOrganizationAuditLogListResponseTypeIPAllowlistDeleted           AdminOrganizationAuditLogListResponseType = "ip_allowlist.deleted"
	AdminOrganizationAuditLogListResponseTypeIPAllowlistConfigActivated   AdminOrganizationAuditLogListResponseType = "ip_allowlist.config.activated"
	AdminOrganizationAuditLogListResponseTypeIPAllowlistConfigDeactivated AdminOrganizationAuditLogListResponseType = "ip_allowlist.config.deactivated"
	AdminOrganizationAuditLogListResponseTypeLoginSucceeded               AdminOrganizationAuditLogListResponseType = "login.succeeded"
	AdminOrganizationAuditLogListResponseTypeLoginFailed                  AdminOrganizationAuditLogListResponseType = "login.failed"
	AdminOrganizationAuditLogListResponseTypeLogoutSucceeded              AdminOrganizationAuditLogListResponseType = "logout.succeeded"
	AdminOrganizationAuditLogListResponseTypeLogoutFailed                 AdminOrganizationAuditLogListResponseType = "logout.failed"
	AdminOrganizationAuditLogListResponseTypeOrganizationUpdated          AdminOrganizationAuditLogListResponseType = "organization.updated"
	AdminOrganizationAuditLogListResponseTypeProjectCreated               AdminOrganizationAuditLogListResponseType = "project.created"
	AdminOrganizationAuditLogListResponseTypeProjectUpdated               AdminOrganizationAuditLogListResponseType = "project.updated"
	AdminOrganizationAuditLogListResponseTypeProjectArchived              AdminOrganizationAuditLogListResponseType = "project.archived"
	AdminOrganizationAuditLogListResponseTypeProjectDeleted               AdminOrganizationAuditLogListResponseType = "project.deleted"
	AdminOrganizationAuditLogListResponseTypeRateLimitUpdated             AdminOrganizationAuditLogListResponseType = "rate_limit.updated"
	AdminOrganizationAuditLogListResponseTypeRateLimitDeleted             AdminOrganizationAuditLogListResponseType = "rate_limit.deleted"
	AdminOrganizationAuditLogListResponseTypeResourceDeleted              AdminOrganizationAuditLogListResponseType = "resource.deleted"
	AdminOrganizationAuditLogListResponseTypeTunnelCreated                AdminOrganizationAuditLogListResponseType = "tunnel.created"
	AdminOrganizationAuditLogListResponseTypeTunnelUpdated                AdminOrganizationAuditLogListResponseType = "tunnel.updated"
	AdminOrganizationAuditLogListResponseTypeTunnelDeleted                AdminOrganizationAuditLogListResponseType = "tunnel.deleted"
	AdminOrganizationAuditLogListResponseTypeRoleCreated                  AdminOrganizationAuditLogListResponseType = "role.created"
	AdminOrganizationAuditLogListResponseTypeRoleUpdated                  AdminOrganizationAuditLogListResponseType = "role.updated"
	AdminOrganizationAuditLogListResponseTypeRoleDeleted                  AdminOrganizationAuditLogListResponseType = "role.deleted"
	AdminOrganizationAuditLogListResponseTypeRoleAssignmentCreated        AdminOrganizationAuditLogListResponseType = "role.assignment.created"
	AdminOrganizationAuditLogListResponseTypeRoleAssignmentDeleted        AdminOrganizationAuditLogListResponseType = "role.assignment.deleted"
	AdminOrganizationAuditLogListResponseTypeScimEnabled                  AdminOrganizationAuditLogListResponseType = "scim.enabled"
	AdminOrganizationAuditLogListResponseTypeScimDisabled                 AdminOrganizationAuditLogListResponseType = "scim.disabled"
	AdminOrganizationAuditLogListResponseTypeServiceAccountCreated        AdminOrganizationAuditLogListResponseType = "service_account.created"
	AdminOrganizationAuditLogListResponseTypeServiceAccountUpdated        AdminOrganizationAuditLogListResponseType = "service_account.updated"
	AdminOrganizationAuditLogListResponseTypeServiceAccountDeleted        AdminOrganizationAuditLogListResponseType = "service_account.deleted"
	AdminOrganizationAuditLogListResponseTypeUserAdded                    AdminOrganizationAuditLogListResponseType = "user.added"
	AdminOrganizationAuditLogListResponseTypeUserUpdated                  AdminOrganizationAuditLogListResponseType = "user.updated"
	AdminOrganizationAuditLogListResponseTypeUserDeleted                  AdminOrganizationAuditLogListResponseType = "user.deleted"
)

// The actor who performed the audit logged action.
type AdminOrganizationAuditLogListResponseActor struct {
	// The API Key used to perform the audit logged action.
	APIKey AdminOrganizationAuditLogListResponseActorAPIKey `json:"api_key"`
	// The session in which the audit logged action was performed.
	Session AdminOrganizationAuditLogListResponseActorSession `json:"session"`
	// The type of actor. Is either `session` or `api_key`.
	//
	// Any of "session", "api_key".
	Type string `json:"type"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		APIKey      respjson.Field
		Session     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseActor) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationAuditLogListResponseActor) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The API Key used to perform the audit logged action.
type AdminOrganizationAuditLogListResponseActorAPIKey struct {
	// The tracking id of the API key.
	ID string `json:"id"`
	// The service account that performed the audit logged action.
	ServiceAccount AdminOrganizationAuditLogListResponseActorAPIKeyServiceAccount `json:"service_account"`
	// The type of API key. Can be either `user` or `service_account`.
	//
	// Any of "user", "service_account".
	Type string `json:"type"`
	// The user who performed the audit logged action.
	User AdminOrganizationAuditLogListResponseActorAPIKeyUser `json:"user"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID             respjson.Field
		ServiceAccount respjson.Field
		Type           respjson.Field
		User           respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseActorAPIKey) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationAuditLogListResponseActorAPIKey) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The service account that performed the audit logged action.
type AdminOrganizationAuditLogListResponseActorAPIKeyServiceAccount struct {
	// The service account id.
	ID string `json:"id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseActorAPIKeyServiceAccount) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationAuditLogListResponseActorAPIKeyServiceAccount) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The user who performed the audit logged action.
type AdminOrganizationAuditLogListResponseActorAPIKeyUser struct {
	// The user id.
	ID string `json:"id"`
	// The user email.
	Email string `json:"email"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Email       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseActorAPIKeyUser) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationAuditLogListResponseActorAPIKeyUser) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The session in which the audit logged action was performed.
type AdminOrganizationAuditLogListResponseActorSession struct {
	// The IP address from which the action was performed.
	IPAddress string `json:"ip_address"`
	// The user who performed the audit logged action.
	User AdminOrganizationAuditLogListResponseActorSessionUser `json:"user"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		IPAddress   respjson.Field
		User        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseActorSession) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationAuditLogListResponseActorSession) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The user who performed the audit logged action.
type AdminOrganizationAuditLogListResponseActorSessionUser struct {
	// The user id.
	ID string `json:"id"`
	// The user email.
	Email string `json:"email"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Email       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseActorSessionUser) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationAuditLogListResponseActorSessionUser) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The details for events with this `type`.
type AdminOrganizationAuditLogListResponseAPIKeyCreated struct {
	// The tracking ID of the API key.
	ID string `json:"id"`
	// The payload used to create the API key.
	Data AdminOrganizationAuditLogListResponseAPIKeyCreatedData `json:"data"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Data        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseAPIKeyCreated) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationAuditLogListResponseAPIKeyCreated) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The payload used to create the API key.
type AdminOrganizationAuditLogListResponseAPIKeyCreatedData struct {
	// A list of scopes allowed for the API key, e.g. `["api.model.request"]`
	Scopes []string `json:"scopes"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Scopes      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseAPIKeyCreatedData) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationAuditLogListResponseAPIKeyCreatedData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The details for events with this `type`.
type AdminOrganizationAuditLogListResponseAPIKeyDeleted struct {
	// The tracking ID of the API key.
	ID string `json:"id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseAPIKeyDeleted) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationAuditLogListResponseAPIKeyDeleted) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The details for events with this `type`.
type AdminOrganizationAuditLogListResponseAPIKeyUpdated struct {
	// The tracking ID of the API key.
	ID string `json:"id"`
	// The payload used to update the API key.
	ChangesRequested AdminOrganizationAuditLogListResponseAPIKeyUpdatedChangesRequested `json:"changes_requested"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID               respjson.Field
		ChangesRequested respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseAPIKeyUpdated) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationAuditLogListResponseAPIKeyUpdated) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The payload used to update the API key.
type AdminOrganizationAuditLogListResponseAPIKeyUpdatedChangesRequested struct {
	// A list of scopes allowed for the API key, e.g. `["api.model.request"]`
	Scopes []string `json:"scopes"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Scopes      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseAPIKeyUpdatedChangesRequested) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationAuditLogListResponseAPIKeyUpdatedChangesRequested) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The details for events with this `type`.
type AdminOrganizationAuditLogListResponseCertificateCreated struct {
	// The certificate ID.
	ID string `json:"id"`
	// The name of the certificate.
	Name string `json:"name"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Name        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseCertificateCreated) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationAuditLogListResponseCertificateCreated) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The details for events with this `type`.
type AdminOrganizationAuditLogListResponseCertificateDeleted struct {
	// The certificate ID.
	ID string `json:"id"`
	// The certificate content in PEM format.
	Certificate string `json:"certificate"`
	// The name of the certificate.
	Name string `json:"name"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Certificate respjson.Field
		Name        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseCertificateDeleted) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationAuditLogListResponseCertificateDeleted) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The details for events with this `type`.
type AdminOrganizationAuditLogListResponseCertificateUpdated struct {
	// The certificate ID.
	ID string `json:"id"`
	// The name of the certificate.
	Name string `json:"name"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Name        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseCertificateUpdated) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationAuditLogListResponseCertificateUpdated) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The details for events with this `type`.
type AdminOrganizationAuditLogListResponseCertificatesActivated struct {
	Certificates []AdminOrganizationAuditLogListResponseCertificatesActivatedCertificate `json:"certificates"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Certificates respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseCertificatesActivated) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationAuditLogListResponseCertificatesActivated) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationAuditLogListResponseCertificatesActivatedCertificate struct {
	// The certificate ID.
	ID string `json:"id"`
	// The name of the certificate.
	Name string `json:"name"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Name        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseCertificatesActivatedCertificate) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationAuditLogListResponseCertificatesActivatedCertificate) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The details for events with this `type`.
type AdminOrganizationAuditLogListResponseCertificatesDeactivated struct {
	Certificates []AdminOrganizationAuditLogListResponseCertificatesDeactivatedCertificate `json:"certificates"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Certificates respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseCertificatesDeactivated) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationAuditLogListResponseCertificatesDeactivated) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationAuditLogListResponseCertificatesDeactivatedCertificate struct {
	// The certificate ID.
	ID string `json:"id"`
	// The name of the certificate.
	Name string `json:"name"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Name        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseCertificatesDeactivatedCertificate) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationAuditLogListResponseCertificatesDeactivatedCertificate) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The project and fine-tuned model checkpoint that the checkpoint permission was
// created for.
type AdminOrganizationAuditLogListResponseCheckpointPermissionCreated struct {
	// The ID of the checkpoint permission.
	ID string `json:"id"`
	// The payload used to create the checkpoint permission.
	Data AdminOrganizationAuditLogListResponseCheckpointPermissionCreatedData `json:"data"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Data        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseCheckpointPermissionCreated) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationAuditLogListResponseCheckpointPermissionCreated) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The payload used to create the checkpoint permission.
type AdminOrganizationAuditLogListResponseCheckpointPermissionCreatedData struct {
	// The ID of the fine-tuned model checkpoint.
	FineTunedModelCheckpoint string `json:"fine_tuned_model_checkpoint"`
	// The ID of the project that the checkpoint permission was created for.
	ProjectID string `json:"project_id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		FineTunedModelCheckpoint respjson.Field
		ProjectID                respjson.Field
		ExtraFields              map[string]respjson.Field
		raw                      string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseCheckpointPermissionCreatedData) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationAuditLogListResponseCheckpointPermissionCreatedData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The details for events with this `type`.
type AdminOrganizationAuditLogListResponseCheckpointPermissionDeleted struct {
	// The ID of the checkpoint permission.
	ID string `json:"id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseCheckpointPermissionDeleted) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationAuditLogListResponseCheckpointPermissionDeleted) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The details for events with this `type`.
type AdminOrganizationAuditLogListResponseExternalKeyRegistered struct {
	// The ID of the external key configuration.
	ID string `json:"id"`
	// The configuration for the external key.
	Data any `json:"data"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Data        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseExternalKeyRegistered) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationAuditLogListResponseExternalKeyRegistered) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The details for events with this `type`.
type AdminOrganizationAuditLogListResponseExternalKeyRemoved struct {
	// The ID of the external key configuration.
	ID string `json:"id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseExternalKeyRemoved) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationAuditLogListResponseExternalKeyRemoved) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The details for events with this `type`.
type AdminOrganizationAuditLogListResponseGroupCreated struct {
	// The ID of the group.
	ID string `json:"id"`
	// Information about the created group.
	Data AdminOrganizationAuditLogListResponseGroupCreatedData `json:"data"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Data        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseGroupCreated) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationAuditLogListResponseGroupCreated) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Information about the created group.
type AdminOrganizationAuditLogListResponseGroupCreatedData struct {
	// The group name.
	GroupName string `json:"group_name"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		GroupName   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseGroupCreatedData) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationAuditLogListResponseGroupCreatedData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The details for events with this `type`.
type AdminOrganizationAuditLogListResponseGroupDeleted struct {
	// The ID of the group.
	ID string `json:"id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseGroupDeleted) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationAuditLogListResponseGroupDeleted) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The details for events with this `type`.
type AdminOrganizationAuditLogListResponseGroupUpdated struct {
	// The ID of the group.
	ID string `json:"id"`
	// The payload used to update the group.
	ChangesRequested AdminOrganizationAuditLogListResponseGroupUpdatedChangesRequested `json:"changes_requested"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID               respjson.Field
		ChangesRequested respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseGroupUpdated) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationAuditLogListResponseGroupUpdated) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The payload used to update the group.
type AdminOrganizationAuditLogListResponseGroupUpdatedChangesRequested struct {
	// The updated group name.
	GroupName string `json:"group_name"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		GroupName   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseGroupUpdatedChangesRequested) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationAuditLogListResponseGroupUpdatedChangesRequested) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The details for events with this `type`.
type AdminOrganizationAuditLogListResponseInviteAccepted struct {
	// The ID of the invite.
	ID string `json:"id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseInviteAccepted) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationAuditLogListResponseInviteAccepted) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The details for events with this `type`.
type AdminOrganizationAuditLogListResponseInviteDeleted struct {
	// The ID of the invite.
	ID string `json:"id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseInviteDeleted) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationAuditLogListResponseInviteDeleted) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The details for events with this `type`.
type AdminOrganizationAuditLogListResponseInviteSent struct {
	// The ID of the invite.
	ID string `json:"id"`
	// The payload used to create the invite.
	Data AdminOrganizationAuditLogListResponseInviteSentData `json:"data"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Data        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseInviteSent) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationAuditLogListResponseInviteSent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The payload used to create the invite.
type AdminOrganizationAuditLogListResponseInviteSentData struct {
	// The email invited to the organization.
	Email string `json:"email"`
	// The role the email was invited to be. Is either `owner` or `member`.
	Role string `json:"role"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Email       respjson.Field
		Role        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseInviteSentData) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationAuditLogListResponseInviteSentData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The details for events with this `type`.
type AdminOrganizationAuditLogListResponseIPAllowlistConfigActivated struct {
	// The configurations that were activated.
	Configs []AdminOrganizationAuditLogListResponseIPAllowlistConfigActivatedConfig `json:"configs"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Configs     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseIPAllowlistConfigActivated) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationAuditLogListResponseIPAllowlistConfigActivated) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationAuditLogListResponseIPAllowlistConfigActivatedConfig struct {
	// The ID of the IP allowlist configuration.
	ID string `json:"id"`
	// The name of the IP allowlist configuration.
	Name string `json:"name"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Name        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseIPAllowlistConfigActivatedConfig) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationAuditLogListResponseIPAllowlistConfigActivatedConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The details for events with this `type`.
type AdminOrganizationAuditLogListResponseIPAllowlistConfigDeactivated struct {
	// The configurations that were deactivated.
	Configs []AdminOrganizationAuditLogListResponseIPAllowlistConfigDeactivatedConfig `json:"configs"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Configs     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseIPAllowlistConfigDeactivated) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationAuditLogListResponseIPAllowlistConfigDeactivated) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationAuditLogListResponseIPAllowlistConfigDeactivatedConfig struct {
	// The ID of the IP allowlist configuration.
	ID string `json:"id"`
	// The name of the IP allowlist configuration.
	Name string `json:"name"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Name        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseIPAllowlistConfigDeactivatedConfig) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationAuditLogListResponseIPAllowlistConfigDeactivatedConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The details for events with this `type`.
type AdminOrganizationAuditLogListResponseIPAllowlistCreated struct {
	// The ID of the IP allowlist configuration.
	ID string `json:"id"`
	// The IP addresses or CIDR ranges included in the configuration.
	AllowedIPs []string `json:"allowed_ips"`
	// The name of the IP allowlist configuration.
	Name string `json:"name"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		AllowedIPs  respjson.Field
		Name        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseIPAllowlistCreated) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationAuditLogListResponseIPAllowlistCreated) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The details for events with this `type`.
type AdminOrganizationAuditLogListResponseIPAllowlistDeleted struct {
	// The ID of the IP allowlist configuration.
	ID string `json:"id"`
	// The IP addresses or CIDR ranges that were in the configuration.
	AllowedIPs []string `json:"allowed_ips"`
	// The name of the IP allowlist configuration.
	Name string `json:"name"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		AllowedIPs  respjson.Field
		Name        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseIPAllowlistDeleted) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationAuditLogListResponseIPAllowlistDeleted) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The details for events with this `type`.
type AdminOrganizationAuditLogListResponseIPAllowlistUpdated struct {
	// The ID of the IP allowlist configuration.
	ID string `json:"id"`
	// The updated set of IP addresses or CIDR ranges in the configuration.
	AllowedIPs []string `json:"allowed_ips"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		AllowedIPs  respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseIPAllowlistUpdated) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationAuditLogListResponseIPAllowlistUpdated) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The details for events with this `type`.
type AdminOrganizationAuditLogListResponseLoginFailed struct {
	// The error code of the failure.
	ErrorCode string `json:"error_code"`
	// The error message of the failure.
	ErrorMessage string `json:"error_message"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode    respjson.Field
		ErrorMessage respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseLoginFailed) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationAuditLogListResponseLoginFailed) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The details for events with this `type`.
type AdminOrganizationAuditLogListResponseLogoutFailed struct {
	// The error code of the failure.
	ErrorCode string `json:"error_code"`
	// The error message of the failure.
	ErrorMessage string `json:"error_message"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode    respjson.Field
		ErrorMessage respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseLogoutFailed) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationAuditLogListResponseLogoutFailed) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The details for events with this `type`.
type AdminOrganizationAuditLogListResponseOrganizationUpdated struct {
	// The organization ID.
	ID string `json:"id"`
	// The payload used to update the organization settings.
	ChangesRequested AdminOrganizationAuditLogListResponseOrganizationUpdatedChangesRequested `json:"changes_requested"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID               respjson.Field
		ChangesRequested respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseOrganizationUpdated) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationAuditLogListResponseOrganizationUpdated) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The payload used to update the organization settings.
type AdminOrganizationAuditLogListResponseOrganizationUpdatedChangesRequested struct {
	// How your organization logs data from supported API calls. One of `disabled`,
	// `enabled_per_call`, `enabled_for_all_projects`, or
	// `enabled_for_selected_projects`
	APICallLogging string `json:"api_call_logging"`
	// The list of project ids if api_call_logging is set to
	// `enabled_for_selected_projects`
	APICallLoggingProjectIDs string `json:"api_call_logging_project_ids"`
	// The organization description.
	Description string `json:"description"`
	// The organization name.
	Name string `json:"name"`
	// Visibility of the threads page which shows messages created with the Assistants
	// API and Playground. One of `ANY_ROLE`, `OWNERS`, or `NONE`.
	ThreadsUiVisibility string `json:"threads_ui_visibility"`
	// The organization title.
	Title string `json:"title"`
	// Visibility of the usage dashboard which shows activity and costs for your
	// organization. One of `ANY_ROLE` or `OWNERS`.
	UsageDashboardVisibility string `json:"usage_dashboard_visibility"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		APICallLogging           respjson.Field
		APICallLoggingProjectIDs respjson.Field
		Description              respjson.Field
		Name                     respjson.Field
		ThreadsUiVisibility      respjson.Field
		Title                    respjson.Field
		UsageDashboardVisibility respjson.Field
		ExtraFields              map[string]respjson.Field
		raw                      string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseOrganizationUpdatedChangesRequested) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationAuditLogListResponseOrganizationUpdatedChangesRequested) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The project that the action was scoped to. Absent for actions not scoped to
// projects. Note that any admin actions taken via Admin API keys are associated
// with the default project.
type AdminOrganizationAuditLogListResponseProject struct {
	// The project ID.
	ID string `json:"id"`
	// The project title.
	Name string `json:"name"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Name        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseProject) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationAuditLogListResponseProject) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The details for events with this `type`.
type AdminOrganizationAuditLogListResponseProjectArchived struct {
	// The project ID.
	ID string `json:"id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseProjectArchived) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationAuditLogListResponseProjectArchived) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The details for events with this `type`.
type AdminOrganizationAuditLogListResponseProjectCreated struct {
	// The project ID.
	ID string `json:"id"`
	// The payload used to create the project.
	Data AdminOrganizationAuditLogListResponseProjectCreatedData `json:"data"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Data        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseProjectCreated) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationAuditLogListResponseProjectCreated) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The payload used to create the project.
type AdminOrganizationAuditLogListResponseProjectCreatedData struct {
	// The project name.
	Name string `json:"name"`
	// The title of the project as seen on the dashboard.
	Title string `json:"title"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Name        respjson.Field
		Title       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseProjectCreatedData) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationAuditLogListResponseProjectCreatedData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The details for events with this `type`.
type AdminOrganizationAuditLogListResponseProjectDeleted struct {
	// The project ID.
	ID string `json:"id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseProjectDeleted) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationAuditLogListResponseProjectDeleted) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The details for events with this `type`.
type AdminOrganizationAuditLogListResponseProjectUpdated struct {
	// The project ID.
	ID string `json:"id"`
	// The payload used to update the project.
	ChangesRequested AdminOrganizationAuditLogListResponseProjectUpdatedChangesRequested `json:"changes_requested"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID               respjson.Field
		ChangesRequested respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseProjectUpdated) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationAuditLogListResponseProjectUpdated) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The payload used to update the project.
type AdminOrganizationAuditLogListResponseProjectUpdatedChangesRequested struct {
	// The title of the project as seen on the dashboard.
	Title string `json:"title"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Title       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseProjectUpdatedChangesRequested) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationAuditLogListResponseProjectUpdatedChangesRequested) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The details for events with this `type`.
type AdminOrganizationAuditLogListResponseRateLimitDeleted struct {
	// The rate limit ID
	ID string `json:"id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseRateLimitDeleted) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationAuditLogListResponseRateLimitDeleted) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The details for events with this `type`.
type AdminOrganizationAuditLogListResponseRateLimitUpdated struct {
	// The rate limit ID
	ID string `json:"id"`
	// The payload used to update the rate limits.
	ChangesRequested AdminOrganizationAuditLogListResponseRateLimitUpdatedChangesRequested `json:"changes_requested"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID               respjson.Field
		ChangesRequested respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseRateLimitUpdated) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationAuditLogListResponseRateLimitUpdated) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The payload used to update the rate limits.
type AdminOrganizationAuditLogListResponseRateLimitUpdatedChangesRequested struct {
	// The maximum batch input tokens per day. Only relevant for certain models.
	Batch1DayMaxInputTokens int64 `json:"batch_1_day_max_input_tokens"`
	// The maximum audio megabytes per minute. Only relevant for certain models.
	MaxAudioMegabytesPer1Minute int64 `json:"max_audio_megabytes_per_1_minute"`
	// The maximum images per minute. Only relevant for certain models.
	MaxImagesPer1Minute int64 `json:"max_images_per_1_minute"`
	// The maximum requests per day. Only relevant for certain models.
	MaxRequestsPer1Day int64 `json:"max_requests_per_1_day"`
	// The maximum requests per minute.
	MaxRequestsPer1Minute int64 `json:"max_requests_per_1_minute"`
	// The maximum tokens per minute.
	MaxTokensPer1Minute int64 `json:"max_tokens_per_1_minute"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Batch1DayMaxInputTokens     respjson.Field
		MaxAudioMegabytesPer1Minute respjson.Field
		MaxImagesPer1Minute         respjson.Field
		MaxRequestsPer1Day          respjson.Field
		MaxRequestsPer1Minute       respjson.Field
		MaxTokensPer1Minute         respjson.Field
		ExtraFields                 map[string]respjson.Field
		raw                         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseRateLimitUpdatedChangesRequested) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationAuditLogListResponseRateLimitUpdatedChangesRequested) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The details for events with this `type`.
type AdminOrganizationAuditLogListResponseRoleAssignmentCreated struct {
	// The identifier of the role assignment.
	ID string `json:"id"`
	// The principal (user or group) that received the role.
	PrincipalID string `json:"principal_id"`
	// The type of principal (user or group) that received the role.
	PrincipalType string `json:"principal_type"`
	// The resource the role assignment is scoped to.
	ResourceID string `json:"resource_id"`
	// The type of resource the role assignment is scoped to.
	ResourceType string `json:"resource_type"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID            respjson.Field
		PrincipalID   respjson.Field
		PrincipalType respjson.Field
		ResourceID    respjson.Field
		ResourceType  respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseRoleAssignmentCreated) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationAuditLogListResponseRoleAssignmentCreated) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The details for events with this `type`.
type AdminOrganizationAuditLogListResponseRoleAssignmentDeleted struct {
	// The identifier of the role assignment.
	ID string `json:"id"`
	// The principal (user or group) that had the role removed.
	PrincipalID string `json:"principal_id"`
	// The type of principal (user or group) that had the role removed.
	PrincipalType string `json:"principal_type"`
	// The resource the role assignment was scoped to.
	ResourceID string `json:"resource_id"`
	// The type of resource the role assignment was scoped to.
	ResourceType string `json:"resource_type"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID            respjson.Field
		PrincipalID   respjson.Field
		PrincipalType respjson.Field
		ResourceID    respjson.Field
		ResourceType  respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseRoleAssignmentDeleted) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationAuditLogListResponseRoleAssignmentDeleted) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The details for events with this `type`.
type AdminOrganizationAuditLogListResponseRoleCreated struct {
	// The role ID.
	ID string `json:"id"`
	// The permissions granted by the role.
	Permissions []string `json:"permissions"`
	// The resource the role is scoped to.
	ResourceID string `json:"resource_id"`
	// The type of resource the role belongs to.
	ResourceType string `json:"resource_type"`
	// The name of the role.
	RoleName string `json:"role_name"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID           respjson.Field
		Permissions  respjson.Field
		ResourceID   respjson.Field
		ResourceType respjson.Field
		RoleName     respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseRoleCreated) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationAuditLogListResponseRoleCreated) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The details for events with this `type`.
type AdminOrganizationAuditLogListResponseRoleDeleted struct {
	// The role ID.
	ID string `json:"id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseRoleDeleted) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationAuditLogListResponseRoleDeleted) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The details for events with this `type`.
type AdminOrganizationAuditLogListResponseRoleUpdated struct {
	// The role ID.
	ID string `json:"id"`
	// The payload used to update the role.
	ChangesRequested AdminOrganizationAuditLogListResponseRoleUpdatedChangesRequested `json:"changes_requested"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID               respjson.Field
		ChangesRequested respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseRoleUpdated) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationAuditLogListResponseRoleUpdated) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The payload used to update the role.
type AdminOrganizationAuditLogListResponseRoleUpdatedChangesRequested struct {
	// The updated role description, when provided.
	Description string `json:"description"`
	// Additional metadata stored on the role.
	Metadata any `json:"metadata"`
	// The permissions added to the role.
	PermissionsAdded []string `json:"permissions_added"`
	// The permissions removed from the role.
	PermissionsRemoved []string `json:"permissions_removed"`
	// The resource the role is scoped to.
	ResourceID string `json:"resource_id"`
	// The type of resource the role belongs to.
	ResourceType string `json:"resource_type"`
	// The updated role name, when provided.
	RoleName string `json:"role_name"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Description        respjson.Field
		Metadata           respjson.Field
		PermissionsAdded   respjson.Field
		PermissionsRemoved respjson.Field
		ResourceID         respjson.Field
		ResourceType       respjson.Field
		RoleName           respjson.Field
		ExtraFields        map[string]respjson.Field
		raw                string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseRoleUpdatedChangesRequested) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationAuditLogListResponseRoleUpdatedChangesRequested) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The details for events with this `type`.
type AdminOrganizationAuditLogListResponseScimDisabled struct {
	// The ID of the SCIM was disabled for.
	ID string `json:"id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseScimDisabled) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationAuditLogListResponseScimDisabled) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The details for events with this `type`.
type AdminOrganizationAuditLogListResponseScimEnabled struct {
	// The ID of the SCIM was enabled for.
	ID string `json:"id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseScimEnabled) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationAuditLogListResponseScimEnabled) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The details for events with this `type`.
type AdminOrganizationAuditLogListResponseServiceAccountCreated struct {
	// The service account ID.
	ID string `json:"id"`
	// The payload used to create the service account.
	Data AdminOrganizationAuditLogListResponseServiceAccountCreatedData `json:"data"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Data        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseServiceAccountCreated) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationAuditLogListResponseServiceAccountCreated) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The payload used to create the service account.
type AdminOrganizationAuditLogListResponseServiceAccountCreatedData struct {
	// The role of the service account. Is either `owner` or `member`.
	Role string `json:"role"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Role        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseServiceAccountCreatedData) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationAuditLogListResponseServiceAccountCreatedData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The details for events with this `type`.
type AdminOrganizationAuditLogListResponseServiceAccountDeleted struct {
	// The service account ID.
	ID string `json:"id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseServiceAccountDeleted) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationAuditLogListResponseServiceAccountDeleted) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The details for events with this `type`.
type AdminOrganizationAuditLogListResponseServiceAccountUpdated struct {
	// The service account ID.
	ID string `json:"id"`
	// The payload used to updated the service account.
	ChangesRequested AdminOrganizationAuditLogListResponseServiceAccountUpdatedChangesRequested `json:"changes_requested"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID               respjson.Field
		ChangesRequested respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseServiceAccountUpdated) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationAuditLogListResponseServiceAccountUpdated) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The payload used to updated the service account.
type AdminOrganizationAuditLogListResponseServiceAccountUpdatedChangesRequested struct {
	// The role of the service account. Is either `owner` or `member`.
	Role string `json:"role"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Role        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseServiceAccountUpdatedChangesRequested) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationAuditLogListResponseServiceAccountUpdatedChangesRequested) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The details for events with this `type`.
type AdminOrganizationAuditLogListResponseUserAdded struct {
	// The user ID.
	ID string `json:"id"`
	// The payload used to add the user to the project.
	Data AdminOrganizationAuditLogListResponseUserAddedData `json:"data"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Data        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseUserAdded) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationAuditLogListResponseUserAdded) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The payload used to add the user to the project.
type AdminOrganizationAuditLogListResponseUserAddedData struct {
	// The role of the user. Is either `owner` or `member`.
	Role string `json:"role"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Role        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseUserAddedData) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationAuditLogListResponseUserAddedData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The details for events with this `type`.
type AdminOrganizationAuditLogListResponseUserDeleted struct {
	// The user ID.
	ID string `json:"id"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseUserDeleted) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationAuditLogListResponseUserDeleted) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The details for events with this `type`.
type AdminOrganizationAuditLogListResponseUserUpdated struct {
	// The project ID.
	ID string `json:"id"`
	// The payload used to update the user.
	ChangesRequested AdminOrganizationAuditLogListResponseUserUpdatedChangesRequested `json:"changes_requested"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID               respjson.Field
		ChangesRequested respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseUserUpdated) RawJSON() string { return r.JSON.raw }
func (r *AdminOrganizationAuditLogListResponseUserUpdated) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The payload used to update the user.
type AdminOrganizationAuditLogListResponseUserUpdatedChangesRequested struct {
	// The role of the user. Is either `owner` or `member`.
	Role string `json:"role"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Role        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AdminOrganizationAuditLogListResponseUserUpdatedChangesRequested) RawJSON() string {
	return r.JSON.raw
}
func (r *AdminOrganizationAuditLogListResponseUserUpdatedChangesRequested) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdminOrganizationAuditLogListParams struct {
	// A cursor for use in pagination. `after` is an object ID that defines your place
	// in the list. For instance, if you make a list request and receive 100 objects,
	// ending with obj_foo, your subsequent call can include after=obj_foo in order to
	// fetch the next page of the list.
	After param.Opt[string] `query:"after,omitzero" json:"-"`
	// A cursor for use in pagination. `before` is an object ID that defines your place
	// in the list. For instance, if you make a list request and receive 100 objects,
	// starting with obj_foo, your subsequent call can include before=obj_foo in order
	// to fetch the previous page of the list.
	Before param.Opt[string] `query:"before,omitzero" json:"-"`
	// A limit on the number of objects to be returned. Limit can range between 1 and
	// 100, and the default is 20.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Return only events performed by users with these emails.
	ActorEmails []string `query:"actor_emails,omitzero" json:"-"`
	// Return only events performed by these actors. Can be a user ID, a service
	// account ID, or an api key tracking ID.
	ActorIDs []string `query:"actor_ids,omitzero" json:"-"`
	// Return only events whose `effective_at` (Unix seconds) is in this range.
	EffectiveAt AdminOrganizationAuditLogListParamsEffectiveAt `query:"effective_at,omitzero" json:"-"`
	// Return only events with a `type` in one of these values. For example,
	// `project.created`. For all options, see the documentation for the
	// [audit log object](https://platform.openai.com/docs/api-reference/audit-logs/object).
	//
	// Any of "api_key.created", "api_key.updated", "api_key.deleted",
	// "certificate.created", "certificate.updated", "certificate.deleted",
	// "certificates.activated", "certificates.deactivated",
	// "checkpoint.permission.created", "checkpoint.permission.deleted",
	// "external_key.registered", "external_key.removed", "group.created",
	// "group.updated", "group.deleted", "invite.sent", "invite.accepted",
	// "invite.deleted", "ip_allowlist.created", "ip_allowlist.updated",
	// "ip_allowlist.deleted", "ip_allowlist.config.activated",
	// "ip_allowlist.config.deactivated", "login.succeeded", "login.failed",
	// "logout.succeeded", "logout.failed", "organization.updated", "project.created",
	// "project.updated", "project.archived", "project.deleted", "rate_limit.updated",
	// "rate_limit.deleted", "resource.deleted", "tunnel.created", "tunnel.updated",
	// "tunnel.deleted", "role.created", "role.updated", "role.deleted",
	// "role.assignment.created", "role.assignment.deleted", "scim.enabled",
	// "scim.disabled", "service_account.created", "service_account.updated",
	// "service_account.deleted", "user.added", "user.updated", "user.deleted".
	EventTypes []string `query:"event_types,omitzero" json:"-"`
	// Return only events for these projects.
	ProjectIDs []string `query:"project_ids,omitzero" json:"-"`
	// Return only events performed on these targets. For example, a project ID
	// updated.
	ResourceIDs []string `query:"resource_ids,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [AdminOrganizationAuditLogListParams]'s query parameters as
// `url.Values`.
func (r AdminOrganizationAuditLogListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Return only events whose `effective_at` (Unix seconds) is in this range.
type AdminOrganizationAuditLogListParamsEffectiveAt struct {
	// Return only events whose `effective_at` (Unix seconds) is greater than this
	// value.
	Gt param.Opt[int64] `query:"gt,omitzero" json:"-"`
	// Return only events whose `effective_at` (Unix seconds) is greater than or equal
	// to this value.
	Gte param.Opt[int64] `query:"gte,omitzero" json:"-"`
	// Return only events whose `effective_at` (Unix seconds) is less than this value.
	Lt param.Opt[int64] `query:"lt,omitzero" json:"-"`
	// Return only events whose `effective_at` (Unix seconds) is less than or equal to
	// this value.
	Lte param.Opt[int64] `query:"lte,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [AdminOrganizationAuditLogListParamsEffectiveAt]'s query
// parameters as `url.Values`.
func (r AdminOrganizationAuditLogListParamsEffectiveAt) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

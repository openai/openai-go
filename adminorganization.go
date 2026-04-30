// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai

import (
	"github.com/openai/openai-go/v3/option"
)

// AdminOrganizationService contains methods and other services that help with
// interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAdminOrganizationService] method instead.
type AdminOrganizationService struct {
	Options []option.RequestOption
	// List user actions and configuration changes within this organization.
	AuditLogs    AdminOrganizationAuditLogService
	AdminAPIKeys AdminOrganizationAdminAPIKeyService
	Usage        AdminOrganizationUsageService
	Invites      AdminOrganizationInviteService
	Users        AdminOrganizationUserService
	Groups       AdminOrganizationGroupService
	Roles        AdminOrganizationRoleService
	Certificates AdminOrganizationCertificateService
	Projects     AdminOrganizationProjectService
}

// NewAdminOrganizationService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewAdminOrganizationService(opts ...option.RequestOption) (r AdminOrganizationService) {
	r = AdminOrganizationService{}
	r.Options = opts
	r.AuditLogs = NewAdminOrganizationAuditLogService(opts...)
	r.AdminAPIKeys = NewAdminOrganizationAdminAPIKeyService(opts...)
	r.Usage = NewAdminOrganizationUsageService(opts...)
	r.Invites = NewAdminOrganizationInviteService(opts...)
	r.Users = NewAdminOrganizationUserService(opts...)
	r.Groups = NewAdminOrganizationGroupService(opts...)
	r.Roles = NewAdminOrganizationRoleService(opts...)
	r.Certificates = NewAdminOrganizationCertificateService(opts...)
	r.Projects = NewAdminOrganizationProjectService(opts...)
	return
}

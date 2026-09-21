// File generated from our OpenAPI spec by Castiron. See CONTRIBUTING.md for details.

package openai

import (
	"context"
	"encoding/json"
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

// BetaAgentVaultCredentialService contains methods and other services that help
// with interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaAgentVaultCredentialService] method instead.
type BetaAgentVaultCredentialService struct {
	Options []option.RequestOption
}

// NewBetaAgentVaultCredentialService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewBetaAgentVaultCredentialService(opts ...option.RequestOption) (r BetaAgentVaultCredentialService) {
	r = BetaAgentVaultCredentialService{}
	r.Options = requestconfig.InheritedOptions(opts...)
	return
}

// Creates a vault credential. Secret values are write-only and are never returned.
// See
// [vaults](https://developers.openai.com/api/docs/guides/agents-api/tools/vaults).
func (r *BetaAgentVaultCredentialService) New(ctx context.Context, vaultID string, body BetaAgentVaultCredentialNewParams, opts ...option.RequestOption) (res *Credential, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "agents=v1")}, opts...)
	if vaultID == "" {
		err = errors.New("missing required vault_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("vaults/%s/credentials", vaultID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Retrieves vault credential metadata without returning secret values. See
// [vaults](https://developers.openai.com/api/docs/guides/agents-api/tools/vaults).
func (r *BetaAgentVaultCredentialService) Get(ctx context.Context, vaultID string, credentialID string, opts ...option.RequestOption) (res *Credential, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "agents=v1")}, opts...)
	if vaultID == "" {
		err = errors.New("missing required vault_id parameter")
		return nil, err
	}
	if credentialID == "" {
		err = errors.New("missing required credential_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("vaults/%s/credentials/%s", vaultID, credentialID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Rotates a vault credential's write-only secret and returns only credential
// metadata. See
// [vaults](https://developers.openai.com/api/docs/guides/agents-api/tools/vaults).
func (r *BetaAgentVaultCredentialService) Update(ctx context.Context, vaultID string, credentialID string, body BetaAgentVaultCredentialUpdateParams, opts ...option.RequestOption) (res *Credential, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "agents=v1")}, opts...)
	if vaultID == "" {
		err = errors.New("missing required vault_id parameter")
		return nil, err
	}
	if credentialID == "" {
		err = errors.New("missing required credential_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("vaults/%s/credentials/%s", vaultID, credentialID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Lists a vault's credentials using ID-based pagination without returning secret
// values. See
// [vaults](https://developers.openai.com/api/docs/guides/agents-api/tools/vaults).
func (r *BetaAgentVaultCredentialService) List(ctx context.Context, vaultID string, query BetaAgentVaultCredentialListParams, opts ...option.RequestOption) (res *pagination.CursorPage[Credential], err error) {
	var raw *http.Response
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "agents=v1"), option.WithResponseInto(&raw)}, opts...)
	if vaultID == "" {
		err = errors.New("missing required vault_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("vaults/%s/credentials", vaultID)
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

// Lists a vault's credentials using ID-based pagination without returning secret
// values. See
// [vaults](https://developers.openai.com/api/docs/guides/agents-api/tools/vaults).
func (r *BetaAgentVaultCredentialService) ListAutoPaging(ctx context.Context, vaultID string, query BetaAgentVaultCredentialListParams, opts ...option.RequestOption) *pagination.CursorPageAutoPager[Credential] {
	return pagination.NewCursorPageAutoPager(r.List(ctx, vaultID, query, opts...))
}

// Deletes a vault credential. See
// [vaults](https://developers.openai.com/api/docs/guides/agents-api/tools/vaults).
func (r *BetaAgentVaultCredentialService) Delete(ctx context.Context, vaultID string, credentialID string, opts ...option.RequestOption) (res *CredentialDeleted, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "agents=v1")}, opts...)
	if vaultID == "" {
		err = errors.New("missing required vault_id parameter")
		return nil, err
	}
	if credentialID == "" {
		err = errors.New("missing required credential_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("vaults/%s/credentials/%s", vaultID, credentialID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Metadata for a stored credential. Secret values are never returned.
type Credential struct {
	// The ID of the credential.
	ID string `json:"id" api:"required"`
	// The authentication method and non-secret configuration of the credential.
	Auth CredentialAuthUnion `json:"auth" api:"required"`
	// The Unix timestamp, in seconds, when the credential was created.
	CreatedAt int64 `json:"created_at" api:"required"`
	// The human-readable name of the credential.
	Name string `json:"name" api:"required"`
	// The object type. Always `vault.credential`.
	Object constant.VaultCredential `json:"object" default:"vault.credential"`
	// The Unix timestamp, in seconds, when the credential was last updated.
	UpdatedAt int64 `json:"updated_at" api:"required"`
	// The ID of the vault containing this credential.
	VaultID string `json:"vault_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Auth        respjson.Field
		CreatedAt   respjson.Field
		Name        respjson.Field
		Object      respjson.Field
		UpdatedAt   respjson.Field
		VaultID     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r Credential) RawJSON() string { return r.JSON.raw }
func (r *Credential) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// CredentialAuthUnion contains all possible properties and values from
// [CredentialAuthMcpOAuth], [CredentialAuthStaticBearer],
// [CredentialAuthEnvironmentVariable].
//
// Use the [CredentialAuthUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type CredentialAuthUnion struct {
	// This field is from variant [CredentialAuthMcpOAuth].
	ExpiresAt    string `json:"expires_at"`
	McpServerURL string `json:"mcp_server_url"`
	// This field is from variant [CredentialAuthMcpOAuth].
	Refresh CredentialAuthMcpOAuthRefresh `json:"refresh"`
	// Any of "mcp_oauth", "static_bearer", "environment_variable".
	Type string `json:"type"`
	// This field is from variant [CredentialAuthEnvironmentVariable].
	Networking CredentialNetworkingUnion `json:"networking"`
	// This field is from variant [CredentialAuthEnvironmentVariable].
	SecretName string `json:"secret_name"`
	JSON       struct {
		ExpiresAt    respjson.Field
		McpServerURL respjson.Field
		Refresh      respjson.Field
		Type         respjson.Field
		Networking   respjson.Field
		SecretName   respjson.Field
		raw          string
	} `json:"-"`
}

// anyCredentialAuth is implemented by each variant of [CredentialAuthUnion] to add
// type safety for the return type of [CredentialAuthUnion.AsAny]
type anyCredentialAuth interface {
	implCredentialAuthUnion()
}

func (CredentialAuthMcpOAuth) implCredentialAuthUnion()            {}
func (CredentialAuthStaticBearer) implCredentialAuthUnion()        {}
func (CredentialAuthEnvironmentVariable) implCredentialAuthUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := CredentialAuthUnion.AsAny().(type) {
//	case openai.CredentialAuthMcpOAuth:
//	case openai.CredentialAuthStaticBearer:
//	case openai.CredentialAuthEnvironmentVariable:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u CredentialAuthUnion) AsAny() anyCredentialAuth {
	switch u.Type {
	case "mcp_oauth":
		return u.AsMcpOAuth()
	case "static_bearer":
		return u.AsStaticBearer()
	case "environment_variable":
		return u.AsEnvironmentVariable()
	}
	return nil
}

func (u CredentialAuthUnion) AsMcpOAuth() (v CredentialAuthMcpOAuth) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u CredentialAuthUnion) AsStaticBearer() (v CredentialAuthStaticBearer) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u CredentialAuthUnion) AsEnvironmentVariable() (v CredentialAuthEnvironmentVariable) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u CredentialAuthUnion) RawJSON() string { return u.JSON.raw }

func (r *CredentialAuthUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Public metadata for an OAuth credential; tokens and client secrets are never
// returned.
type CredentialAuthMcpOAuth struct {
	// When the OAuth access token expires, as an RFC 3339 timestamp, if known.
	ExpiresAt string `json:"expires_at" api:"required"`
	// The HTTPS MCP server URL authorized by this credential.
	McpServerURL string `json:"mcp_server_url" api:"required"`
	// Configuration used to refresh an MCP OAuth access token, excluding secret
	// values.
	Refresh CredentialAuthMcpOAuthRefresh `json:"refresh" api:"required"`
	// The type of the object. Always `mcp_oauth`.
	Type constant.McpOAuth `json:"type" default:"mcp_oauth"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ExpiresAt    respjson.Field
		McpServerURL respjson.Field
		Refresh      respjson.Field
		Type         respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CredentialAuthMcpOAuth) RawJSON() string { return r.JSON.raw }
func (r *CredentialAuthMcpOAuth) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration used to refresh an MCP OAuth access token, excluding secret
// values.
type CredentialAuthMcpOAuthRefresh struct {
	// The OAuth client ID used when requesting a new access token.
	ClientID string `json:"client_id" api:"required"`
	// The resource URI sent to the OAuth token endpoint during refresh, if configured.
	Resource string `json:"resource" api:"required"`
	// Space-separated OAuth scopes requested during refresh, if configured.
	Scope string `json:"scope" api:"required"`
	// The HTTPS OAuth token endpoint used for refresh.
	TokenEndpoint string `json:"token_endpoint" api:"required"`
	// How the OAuth client authenticates to the token endpoint, excluding its client
	// secret.
	TokenEndpointAuth McpOAuthTokenEndpointAuthUnion `json:"token_endpoint_auth" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ClientID          respjson.Field
		Resource          respjson.Field
		Scope             respjson.Field
		TokenEndpoint     respjson.Field
		TokenEndpointAuth respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CredentialAuthMcpOAuthRefresh) RawJSON() string { return r.JSON.raw }
func (r *CredentialAuthMcpOAuthRefresh) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Metadata for a bearer-token credential, without automatic OAuth refresh.
type CredentialAuthStaticBearer struct {
	// The HTTPS MCP server URL authorized by this credential.
	McpServerURL string `json:"mcp_server_url" api:"required"`
	// The type of the object. Always `static_bearer`.
	Type constant.StaticBearer `json:"type" default:"static_bearer"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		McpServerURL respjson.Field
		Type         respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CredentialAuthStaticBearer) RawJSON() string { return r.JSON.raw }
func (r *CredentialAuthStaticBearer) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Metadata for an HTTP credential used only in OpenAI-hosted environments. Sandbox
// code receives a placeholder. The proxy substitutes the secret for allowed HTTPS
// destinations on ports 443 and 8443. The real secret is not available to sandbox
// code for local computation and is never returned in this resource.
type CredentialAuthEnvironmentVariable struct {
	// The destinations where the proxy can substitute the secret, subject to the
	// environment network policy.
	Networking CredentialNetworkingUnion `json:"networking" api:"required"`
	// The environment variable name that receives the placeholder in the sandbox.
	SecretName string `json:"secret_name" api:"required"`
	// The type of the object. Always `environment_variable`.
	Type constant.EnvironmentVariable `json:"type" default:"environment_variable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Networking  respjson.Field
		SecretName  respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CredentialAuthEnvironmentVariable) RawJSON() string { return r.JSON.raw }
func (r *CredentialAuthEnvironmentVariable) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func CredentialAuthCreateParamOfParamMcpOAuth(accessToken string, mcpServerURL string) CredentialAuthCreateParamUnion {
	var paramMcpOAuth CredentialAuthCreateParamMcpOAuth
	paramMcpOAuth.AccessToken = accessToken
	paramMcpOAuth.McpServerURL = mcpServerURL
	return CredentialAuthCreateParamUnion{OfParamMcpOAuth: &paramMcpOAuth}
}

func CredentialAuthCreateParamOfParamStaticBearer(mcpServerURL string, token string) CredentialAuthCreateParamUnion {
	var paramStaticBearer CredentialAuthCreateParamStaticBearer
	paramStaticBearer.McpServerURL = mcpServerURL
	paramStaticBearer.Token = token
	return CredentialAuthCreateParamUnion{OfParamStaticBearer: &paramStaticBearer}
}

func CredentialAuthCreateParamOfParamEnvironmentVariable[
	T CredentialNetworkingParamUnrestricted | CredentialNetworkingParamLimited,
](networking T, secretName string, secretValue string) CredentialAuthCreateParamUnion {
	var paramEnvironmentVariable CredentialAuthCreateParamEnvironmentVariable
	switch v := any(networking).(type) {
	case CredentialNetworkingParamUnrestricted:
		paramEnvironmentVariable.Networking.OfParamUnrestricted = &v
	case CredentialNetworkingParamLimited:
		paramEnvironmentVariable.Networking.OfParamLimited = &v
	}
	paramEnvironmentVariable.SecretName = secretName
	paramEnvironmentVariable.SecretValue = secretValue
	return CredentialAuthCreateParamUnion{OfParamEnvironmentVariable: &paramEnvironmentVariable}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type CredentialAuthCreateParamUnion struct {
	OfParamMcpOAuth            *CredentialAuthCreateParamMcpOAuth            `json:",omitzero,inline"`
	OfParamStaticBearer        *CredentialAuthCreateParamStaticBearer        `json:",omitzero,inline"`
	OfParamEnvironmentVariable *CredentialAuthCreateParamEnvironmentVariable `json:",omitzero,inline"`
	paramUnion
}

func (u CredentialAuthCreateParamUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfParamMcpOAuth, u.OfParamStaticBearer, u.OfParamEnvironmentVariable)
}
func (u *CredentialAuthCreateParamUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u CredentialAuthCreateParamUnion) GetAccessToken() *string {
	if vt := u.OfParamMcpOAuth; vt != nil {
		return &vt.AccessToken
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u CredentialAuthCreateParamUnion) GetExpiresAt() *string {
	if vt := u.OfParamMcpOAuth; vt != nil && vt.ExpiresAt.Valid() {
		return &vt.ExpiresAt.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u CredentialAuthCreateParamUnion) GetRefresh() *CredentialAuthCreateParamMcpOAuthRefresh {
	if vt := u.OfParamMcpOAuth; vt != nil {
		return &vt.Refresh
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u CredentialAuthCreateParamUnion) GetToken() *string {
	if vt := u.OfParamStaticBearer; vt != nil {
		return &vt.Token
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u CredentialAuthCreateParamUnion) GetNetworking() *CredentialNetworkingParamUnion {
	if vt := u.OfParamEnvironmentVariable; vt != nil {
		return &vt.Networking
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u CredentialAuthCreateParamUnion) GetSecretName() *string {
	if vt := u.OfParamEnvironmentVariable; vt != nil {
		return &vt.SecretName
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u CredentialAuthCreateParamUnion) GetSecretValue() *string {
	if vt := u.OfParamEnvironmentVariable; vt != nil {
		return &vt.SecretValue
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u CredentialAuthCreateParamUnion) GetMcpServerURL() *string {
	if vt := u.OfParamMcpOAuth; vt != nil {
		return (*string)(&vt.McpServerURL)
	} else if vt := u.OfParamStaticBearer; vt != nil {
		return (*string)(&vt.McpServerURL)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u CredentialAuthCreateParamUnion) GetType() *string {
	if vt := u.OfParamMcpOAuth; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfParamStaticBearer; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfParamEnvironmentVariable; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[CredentialAuthCreateParamUnion](
		"type",
		apijson.Discriminator[CredentialAuthCreateParamMcpOAuth]("mcp_oauth"),
		apijson.Discriminator[CredentialAuthCreateParamStaticBearer]("static_bearer"),
		apijson.Discriminator[CredentialAuthCreateParamEnvironmentVariable]("environment_variable"),
	)
}

// An OAuth credential for an HTTPS MCP destination.
//
// The properties AccessToken, McpServerURL, Type are required.
type CredentialAuthCreateParamMcpOAuth struct {
	// A write-only OAuth access token; never returned by credential resources.
	AccessToken string `json:"access_token" api:"required"`
	// The HTTPS MCP server URL authorized by this credential.
	McpServerURL string `json:"mcp_server_url" api:"required"`
	// When the OAuth access token expires, as an RFC 3339 timestamp, if known.
	ExpiresAt param.Opt[string] `json:"expires_at,omitzero"`
	// Configuration for refreshing the access token of an MCP OAuth credential.
	Refresh CredentialAuthCreateParamMcpOAuthRefresh `json:"refresh,omitzero"`
	// The type of the object. Always `mcp_oauth`.
	//
	// This field can be elided, and will marshal its zero value as "mcp_oauth".
	Type constant.McpOAuth `json:"type" default:"mcp_oauth"`
	paramObj
}

func (r CredentialAuthCreateParamMcpOAuth) MarshalJSON() (data []byte, err error) {
	type shadow CredentialAuthCreateParamMcpOAuth
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *CredentialAuthCreateParamMcpOAuth) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for refreshing the access token of an MCP OAuth credential.
//
// The properties ClientID, RefreshToken, TokenEndpoint, TokenEndpointAuth are
// required.
type CredentialAuthCreateParamMcpOAuthRefresh struct {
	// The OAuth client ID used when requesting a new access token.
	ClientID string `json:"client_id" api:"required"`
	// The refresh token to store. This secret is never returned in credential
	// resources.
	RefreshToken string `json:"refresh_token" api:"required"`
	// The HTTPS OAuth token endpoint used to exchange the refresh token for a new
	// access token.
	TokenEndpoint string `json:"token_endpoint" api:"required"`
	// How the OAuth client authenticates to the token endpoint.
	TokenEndpointAuth McpOAuthTokenEndpointAuthCreateParamUnion `json:"token_endpoint_auth,omitzero" api:"required"`
	// The resource URI to send to the OAuth token endpoint during refresh, if
	// required.
	Resource param.Opt[string] `json:"resource,omitzero"`
	// Space-separated OAuth scopes to request during refresh, if required.
	Scope param.Opt[string] `json:"scope,omitzero"`
	paramObj
}

func (r CredentialAuthCreateParamMcpOAuthRefresh) MarshalJSON() (data []byte, err error) {
	type shadow CredentialAuthCreateParamMcpOAuthRefresh
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *CredentialAuthCreateParamMcpOAuthRefresh) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A bearer token for an MCP server, without automatic OAuth refresh.
//
// The properties Token, McpServerURL, Type are required.
type CredentialAuthCreateParamStaticBearer struct {
	// The bearer token to store. This secret is never returned in credential
	// resources.
	Token string `json:"token" api:"required"`
	// The HTTPS MCP server URL authorized by this credential.
	McpServerURL string `json:"mcp_server_url" api:"required"`
	// The type of the object. Always `static_bearer`.
	//
	// This field can be elided, and will marshal its zero value as "static_bearer".
	Type constant.StaticBearer `json:"type" default:"static_bearer"`
	paramObj
}

func (r CredentialAuthCreateParamStaticBearer) MarshalJSON() (data []byte, err error) {
	type shadow CredentialAuthCreateParamStaticBearer
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *CredentialAuthCreateParamStaticBearer) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// An HTTP credential for OpenAI-hosted environments only. The sandbox receives an
// environment variable containing a placeholder, not the secret. Use the
// placeholder unchanged in outgoing requests. The egress proxy replaces the
// placeholder with the secret for allowed HTTPS destinations on ports 443 and
// 8443. Sandbox code cannot read the real secret or use it for local computation,
// such as signing a request.
//
// The properties Networking, SecretName, SecretValue, Type are required.
type CredentialAuthCreateParamEnvironmentVariable struct {
	// The destinations where the proxy can substitute this secret. The environment
	// network policy must also allow them.
	Networking CredentialNetworkingParamUnion `json:"networking,omitzero" api:"required"`
	// The environment variable name that receives the placeholder, such as
	// `SERVICE_API_KEY`. Use ASCII letters, digits, and underscores, starting with a
	// letter or underscore. Names starting with `CODEX_` and managed proxy or
	// certificate variable names are reserved.
	SecretName string `json:"secret_name" api:"required"`
	// The write-only secret to store. Never returned in credential resources or
	// supplied directly to sandbox code. Must be nonempty and must not contain
	// carriage returns, newlines, or NUL bytes.
	SecretValue string `json:"secret_value" api:"required"`
	// The type of the object. Always `environment_variable`.
	//
	// This field can be elided, and will marshal its zero value as
	// "environment_variable".
	Type constant.EnvironmentVariable `json:"type" default:"environment_variable"`
	paramObj
}

func (r CredentialAuthCreateParamEnvironmentVariable) MarshalJSON() (data []byte, err error) {
	type shadow CredentialAuthCreateParamEnvironmentVariable
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *CredentialAuthCreateParamEnvironmentVariable) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func CredentialAuthRotateParamOfParamStaticBearer(token string) CredentialAuthRotateParamUnion {
	var paramStaticBearer CredentialAuthRotateParamStaticBearer
	paramStaticBearer.Token = token
	return CredentialAuthRotateParamUnion{OfParamStaticBearer: &paramStaticBearer}
}

func CredentialAuthRotateParamOfParamEnvironmentVariable(secretValue string) CredentialAuthRotateParamUnion {
	var paramEnvironmentVariable CredentialAuthRotateParamEnvironmentVariable
	paramEnvironmentVariable.SecretValue = secretValue
	return CredentialAuthRotateParamUnion{OfParamEnvironmentVariable: &paramEnvironmentVariable}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type CredentialAuthRotateParamUnion struct {
	OfParamMcpOAuth            *CredentialAuthRotateParamMcpOAuth            `json:",omitzero,inline"`
	OfParamStaticBearer        *CredentialAuthRotateParamStaticBearer        `json:",omitzero,inline"`
	OfParamEnvironmentVariable *CredentialAuthRotateParamEnvironmentVariable `json:",omitzero,inline"`
	paramUnion
}

func (u CredentialAuthRotateParamUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfParamMcpOAuth, u.OfParamStaticBearer, u.OfParamEnvironmentVariable)
}
func (u *CredentialAuthRotateParamUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u CredentialAuthRotateParamUnion) GetAccessToken() *string {
	if vt := u.OfParamMcpOAuth; vt != nil && vt.AccessToken.Valid() {
		return &vt.AccessToken.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u CredentialAuthRotateParamUnion) GetExpiresAt() *string {
	if vt := u.OfParamMcpOAuth; vt != nil && vt.ExpiresAt.Valid() {
		return &vt.ExpiresAt.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u CredentialAuthRotateParamUnion) GetRefresh() *CredentialAuthRotateParamMcpOAuthRefresh {
	if vt := u.OfParamMcpOAuth; vt != nil {
		return &vt.Refresh
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u CredentialAuthRotateParamUnion) GetToken() *string {
	if vt := u.OfParamStaticBearer; vt != nil {
		return &vt.Token
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u CredentialAuthRotateParamUnion) GetSecretValue() *string {
	if vt := u.OfParamEnvironmentVariable; vt != nil {
		return &vt.SecretValue
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u CredentialAuthRotateParamUnion) GetType() *string {
	if vt := u.OfParamMcpOAuth; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfParamStaticBearer; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfParamEnvironmentVariable; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[CredentialAuthRotateParamUnion](
		"type",
		apijson.Discriminator[CredentialAuthRotateParamMcpOAuth]("mcp_oauth"),
		apijson.Discriminator[CredentialAuthRotateParamStaticBearer]("static_bearer"),
		apijson.Discriminator[CredentialAuthRotateParamEnvironmentVariable]("environment_variable"),
	)
}

// Rotate an OAuth credential for an HTTPS MCP destination.
//
// The property Type is required.
type CredentialAuthRotateParamMcpOAuth struct {
	// A write-only replacement OAuth access token.
	AccessToken param.Opt[string] `json:"access_token,omitzero"`
	// The replacement expiry as an RFC 3339 timestamp, or `null` to clear it. Omitting
	// this field preserves the expiry unless a new access token is supplied, in which
	// case the expiry is cleared.
	ExpiresAt param.Opt[string] `json:"expires_at,omitzero"`
	// Updates to an MCP credential's existing OAuth refresh configuration.
	Refresh CredentialAuthRotateParamMcpOAuthRefresh `json:"refresh,omitzero"`
	// The type of the object. Always `mcp_oauth`.
	//
	// This field can be elided, and will marshal its zero value as "mcp_oauth".
	Type constant.McpOAuth `json:"type" default:"mcp_oauth"`
	paramObj
}

func (r CredentialAuthRotateParamMcpOAuth) MarshalJSON() (data []byte, err error) {
	type shadow CredentialAuthRotateParamMcpOAuth
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *CredentialAuthRotateParamMcpOAuth) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Updates to an MCP credential's existing OAuth refresh configuration.
type CredentialAuthRotateParamMcpOAuthRefresh struct {
	// The replacement refresh token. Omit or pass `null` to keep the stored token.
	// This secret is never returned in resources.
	RefreshToken param.Opt[string] `json:"refresh_token,omitzero"`
	// Replacement space-separated OAuth scopes for refresh requests. Omit to keep the
	// scopes, or pass `null` to stop sending a scope parameter.
	Scope param.Opt[string] `json:"scope,omitzero"`
	// Client-secret updates that preserve the credential's OAuth authentication
	// method.
	TokenEndpointAuth McpOAuthTokenEndpointAuthRotateParamUnion `json:"token_endpoint_auth,omitzero"`
	paramObj
}

func (r CredentialAuthRotateParamMcpOAuthRefresh) MarshalJSON() (data []byte, err error) {
	type shadow CredentialAuthRotateParamMcpOAuthRefresh
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *CredentialAuthRotateParamMcpOAuthRefresh) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Replace the bearer token for the credential's MCP server.
//
// The properties Token, Type are required.
type CredentialAuthRotateParamStaticBearer struct {
	// The replacement bearer token. This secret is never returned in credential
	// resources.
	Token string `json:"token" api:"required"`
	// The type of the object. Always `static_bearer`.
	//
	// This field can be elided, and will marshal its zero value as "static_bearer".
	Type constant.StaticBearer `json:"type" default:"static_bearer"`
	paramObj
}

func (r CredentialAuthRotateParamStaticBearer) MarshalJSON() (data []byte, err error) {
	type shadow CredentialAuthRotateParamStaticBearer
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *CredentialAuthRotateParamStaticBearer) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Replace the secret for an OpenAI-hosted environment credential. The environment
// variable name and networking configuration remain unchanged.
//
// The properties SecretValue, Type are required.
type CredentialAuthRotateParamEnvironmentVariable struct {
	// The write-only replacement secret. Never returned in credential resources or
	// supplied directly to sandbox code. Must be nonempty and must not contain
	// carriage returns, newlines, or NUL bytes.
	SecretValue string `json:"secret_value" api:"required"`
	// The type of the object. Always `environment_variable`.
	//
	// This field can be elided, and will marshal its zero value as
	// "environment_variable".
	Type constant.EnvironmentVariable `json:"type" default:"environment_variable"`
	paramObj
}

func (r CredentialAuthRotateParamEnvironmentVariable) MarshalJSON() (data []byte, err error) {
	type shadow CredentialAuthRotateParamEnvironmentVariable
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *CredentialAuthRotateParamEnvironmentVariable) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Confirmation that a vault credential was deleted.
type CredentialDeleted struct {
	// The ID of the deleted credential.
	ID string `json:"id" api:"required"`
	// Whether the resource was deleted. Always `true`.
	Deleted bool `json:"deleted" api:"required"`
	// The object type. Always `vault.credential.deleted`.
	Object constant.VaultCredentialDeleted `json:"object" default:"vault.credential.deleted"`
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
func (r CredentialDeleted) RawJSON() string { return r.JSON.raw }
func (r *CredentialDeleted) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// CredentialNetworkingUnion contains all possible properties and values from
// [CredentialNetworkingUnrestricted], [CredentialNetworkingLimited].
//
// Use the [CredentialNetworkingUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type CredentialNetworkingUnion struct {
	// Any of "unrestricted", "limited".
	Type string `json:"type"`
	// This field is from variant [CredentialNetworkingLimited].
	AllowedHosts []string `json:"allowed_hosts"`
	JSON         struct {
		Type         respjson.Field
		AllowedHosts respjson.Field
		raw          string
	} `json:"-"`
}

// anyCredentialNetworking is implemented by each variant of
// [CredentialNetworkingUnion] to add type safety for the return type of
// [CredentialNetworkingUnion.AsAny]
type anyCredentialNetworking interface {
	implCredentialNetworkingUnion()
}

func (CredentialNetworkingUnrestricted) implCredentialNetworkingUnion() {}
func (CredentialNetworkingLimited) implCredentialNetworkingUnion()      {}

// Use the following switch statement to find the correct variant
//
//	switch variant := CredentialNetworkingUnion.AsAny().(type) {
//	case openai.CredentialNetworkingUnrestricted:
//	case openai.CredentialNetworkingLimited:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u CredentialNetworkingUnion) AsAny() anyCredentialNetworking {
	switch u.Type {
	case "unrestricted":
		return u.AsUnrestricted()
	case "limited":
		return u.AsLimited()
	}
	return nil
}

func (u CredentialNetworkingUnion) AsUnrestricted() (v CredentialNetworkingUnrestricted) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u CredentialNetworkingUnion) AsLimited() (v CredentialNetworkingLimited) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u CredentialNetworkingUnion) RawJSON() string { return u.JSON.raw }

func (r *CredentialNetworkingUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Allows substitution for destinations permitted by the environment network
// policy. Requires `environment.network.access` to be `restricted`, with explicit
// `allowed_domains`.
type CredentialNetworkingUnrestricted struct {
	// The type of the object. Always `unrestricted`.
	Type constant.Unrestricted `json:"type" default:"unrestricted"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CredentialNetworkingUnrestricted) RawJSON() string { return r.JSON.raw }
func (r *CredentialNetworkingUnrestricted) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Allows substitution only for the listed hosts. The environment network policy
// must also allow these hosts.
type CredentialNetworkingLimited struct {
	// The 1 to 16 distinct allowed hostnames or IPv4 addresses, normalized to
	// lowercase. Entries contain no scheme, path, port, or wildcard. IPv6 addresses
	// are not supported.
	AllowedHosts []string `json:"allowed_hosts" api:"required"`
	// The type of the object. Always `limited`.
	Type constant.Limited `json:"type" default:"limited"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		AllowedHosts respjson.Field
		Type         respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CredentialNetworkingLimited) RawJSON() string { return r.JSON.raw }
func (r *CredentialNetworkingLimited) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func CredentialNetworkingParamOfParamLimited(allowedHosts []string) CredentialNetworkingParamUnion {
	var paramLimited CredentialNetworkingParamLimited
	paramLimited.AllowedHosts = allowedHosts
	return CredentialNetworkingParamUnion{OfParamLimited: &paramLimited}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type CredentialNetworkingParamUnion struct {
	OfParamUnrestricted *CredentialNetworkingParamUnrestricted `json:",omitzero,inline"`
	OfParamLimited      *CredentialNetworkingParamLimited      `json:",omitzero,inline"`
	paramUnion
}

func (u CredentialNetworkingParamUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfParamUnrestricted, u.OfParamLimited)
}
func (u *CredentialNetworkingParamUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u CredentialNetworkingParamUnion) GetAllowedHosts() []string {
	if vt := u.OfParamLimited; vt != nil {
		return vt.AllowedHosts
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u CredentialNetworkingParamUnion) GetType() *string {
	if vt := u.OfParamUnrestricted; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfParamLimited; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[CredentialNetworkingParamUnion](
		"type",
		apijson.Discriminator[CredentialNetworkingParamUnrestricted]("unrestricted"),
		apijson.Discriminator[CredentialNetworkingParamLimited]("limited"),
	)
}

func NewCredentialNetworkingParamUnrestricted() CredentialNetworkingParamUnrestricted {
	return CredentialNetworkingParamUnrestricted{
		Type: "unrestricted",
	}
}

// Allows substitution for destinations permitted by the environment network
// policy. Requires `environment.network.access` to be `restricted`, with explicit
// `allowed_domains`.
//
// This struct has a constant value, construct it with
// [NewCredentialNetworkingParamUnrestricted].
type CredentialNetworkingParamUnrestricted struct {
	// The type of the object. Always `unrestricted`.
	Type constant.Unrestricted `json:"type" default:"unrestricted"`
	paramObj
}

func (r CredentialNetworkingParamUnrestricted) MarshalJSON() (data []byte, err error) {
	type shadow CredentialNetworkingParamUnrestricted
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *CredentialNetworkingParamUnrestricted) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Allows substitution only for the listed hosts. The environment network policy
// must also allow these hosts.
//
// The properties AllowedHosts, Type are required.
type CredentialNetworkingParamLimited struct {
	// The 1 to 16 distinct allowed hostnames or IPv4 addresses, normalized to
	// lowercase. Entries contain no scheme, path, port, or wildcard. IPv6 addresses
	// are not supported.
	AllowedHosts []string `json:"allowed_hosts,omitzero" api:"required"`
	// The type of the object. Always `limited`.
	//
	// This field can be elided, and will marshal its zero value as "limited".
	Type constant.Limited `json:"type" default:"limited"`
	paramObj
}

func (r CredentialNetworkingParamLimited) MarshalJSON() (data []byte, err error) {
	type shadow CredentialNetworkingParamLimited
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *CredentialNetworkingParamLimited) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// McpOAuthTokenEndpointAuthUnion contains all possible properties and values from
// [McpOAuthTokenEndpointAuthNone], [McpOAuthTokenEndpointAuthClientSecretBasic],
// [McpOAuthTokenEndpointAuthClientSecretPost].
//
// Use the [McpOAuthTokenEndpointAuthUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type McpOAuthTokenEndpointAuthUnion struct {
	// Any of "none", "client_secret_basic", "client_secret_post".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyMcpOAuthTokenEndpointAuth is implemented by each variant of
// [McpOAuthTokenEndpointAuthUnion] to add type safety for the return type of
// [McpOAuthTokenEndpointAuthUnion.AsAny]
type anyMcpOAuthTokenEndpointAuth interface {
	implMcpOAuthTokenEndpointAuthUnion()
}

func (McpOAuthTokenEndpointAuthNone) implMcpOAuthTokenEndpointAuthUnion()              {}
func (McpOAuthTokenEndpointAuthClientSecretBasic) implMcpOAuthTokenEndpointAuthUnion() {}
func (McpOAuthTokenEndpointAuthClientSecretPost) implMcpOAuthTokenEndpointAuthUnion()  {}

// Use the following switch statement to find the correct variant
//
//	switch variant := McpOAuthTokenEndpointAuthUnion.AsAny().(type) {
//	case openai.McpOAuthTokenEndpointAuthNone:
//	case openai.McpOAuthTokenEndpointAuthClientSecretBasic:
//	case openai.McpOAuthTokenEndpointAuthClientSecretPost:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u McpOAuthTokenEndpointAuthUnion) AsAny() anyMcpOAuthTokenEndpointAuth {
	switch u.Type {
	case "none":
		return u.AsNone()
	case "client_secret_basic":
		return u.AsClientSecretBasic()
	case "client_secret_post":
		return u.AsClientSecretPost()
	}
	return nil
}

func (u McpOAuthTokenEndpointAuthUnion) AsNone() (v McpOAuthTokenEndpointAuthNone) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u McpOAuthTokenEndpointAuthUnion) AsClientSecretBasic() (v McpOAuthTokenEndpointAuthClientSecretBasic) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u McpOAuthTokenEndpointAuthUnion) AsClientSecretPost() (v McpOAuthTokenEndpointAuthClientSecretPost) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u McpOAuthTokenEndpointAuthUnion) RawJSON() string { return u.JSON.raw }

func (r *McpOAuthTokenEndpointAuthUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Sends the client ID without a client secret.
type McpOAuthTokenEndpointAuthNone struct {
	// The type of the object. Always `none`.
	Type constant.None `json:"type" default:"none"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r McpOAuthTokenEndpointAuthNone) RawJSON() string { return r.JSON.raw }
func (r *McpOAuthTokenEndpointAuthNone) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Sends the client ID and secret using HTTP Basic authentication.
type McpOAuthTokenEndpointAuthClientSecretBasic struct {
	// The type of the object. Always `client_secret_basic`.
	Type constant.ClientSecretBasic `json:"type" default:"client_secret_basic"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r McpOAuthTokenEndpointAuthClientSecretBasic) RawJSON() string { return r.JSON.raw }
func (r *McpOAuthTokenEndpointAuthClientSecretBasic) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Sends the client ID and secret in the token request body.
type McpOAuthTokenEndpointAuthClientSecretPost struct {
	// The type of the object. Always `client_secret_post`.
	Type constant.ClientSecretPost `json:"type" default:"client_secret_post"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r McpOAuthTokenEndpointAuthClientSecretPost) RawJSON() string { return r.JSON.raw }
func (r *McpOAuthTokenEndpointAuthClientSecretPost) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func McpOAuthTokenEndpointAuthCreateParamOfParamClientSecretBasic(clientSecret string) McpOAuthTokenEndpointAuthCreateParamUnion {
	var paramClientSecretBasic McpOAuthTokenEndpointAuthCreateParamClientSecretBasic
	paramClientSecretBasic.ClientSecret = clientSecret
	return McpOAuthTokenEndpointAuthCreateParamUnion{OfParamClientSecretBasic: &paramClientSecretBasic}
}

func McpOAuthTokenEndpointAuthCreateParamOfParamClientSecretPost(clientSecret string) McpOAuthTokenEndpointAuthCreateParamUnion {
	var paramClientSecretPost McpOAuthTokenEndpointAuthCreateParamClientSecretPost
	paramClientSecretPost.ClientSecret = clientSecret
	return McpOAuthTokenEndpointAuthCreateParamUnion{OfParamClientSecretPost: &paramClientSecretPost}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type McpOAuthTokenEndpointAuthCreateParamUnion struct {
	OfParamNone              *McpOAuthTokenEndpointAuthCreateParamNone              `json:",omitzero,inline"`
	OfParamClientSecretBasic *McpOAuthTokenEndpointAuthCreateParamClientSecretBasic `json:",omitzero,inline"`
	OfParamClientSecretPost  *McpOAuthTokenEndpointAuthCreateParamClientSecretPost  `json:",omitzero,inline"`
	paramUnion
}

func (u McpOAuthTokenEndpointAuthCreateParamUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfParamNone, u.OfParamClientSecretBasic, u.OfParamClientSecretPost)
}
func (u *McpOAuthTokenEndpointAuthCreateParamUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u McpOAuthTokenEndpointAuthCreateParamUnion) GetType() *string {
	if vt := u.OfParamNone; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfParamClientSecretBasic; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfParamClientSecretPost; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u McpOAuthTokenEndpointAuthCreateParamUnion) GetClientSecret() *string {
	if vt := u.OfParamClientSecretBasic; vt != nil {
		return (*string)(&vt.ClientSecret)
	} else if vt := u.OfParamClientSecretPost; vt != nil {
		return (*string)(&vt.ClientSecret)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[McpOAuthTokenEndpointAuthCreateParamUnion](
		"type",
		apijson.Discriminator[McpOAuthTokenEndpointAuthCreateParamNone]("none"),
		apijson.Discriminator[McpOAuthTokenEndpointAuthCreateParamClientSecretBasic]("client_secret_basic"),
		apijson.Discriminator[McpOAuthTokenEndpointAuthCreateParamClientSecretPost]("client_secret_post"),
	)
}

func NewMcpOAuthTokenEndpointAuthCreateParamNone() McpOAuthTokenEndpointAuthCreateParamNone {
	return McpOAuthTokenEndpointAuthCreateParamNone{
		Type: "none",
	}
}

// Sends the client ID without a client secret.
//
// This struct has a constant value, construct it with
// [NewMcpOAuthTokenEndpointAuthCreateParamNone].
type McpOAuthTokenEndpointAuthCreateParamNone struct {
	// The type of the object. Always `none`.
	Type constant.None `json:"type" default:"none"`
	paramObj
}

func (r McpOAuthTokenEndpointAuthCreateParamNone) MarshalJSON() (data []byte, err error) {
	type shadow McpOAuthTokenEndpointAuthCreateParamNone
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *McpOAuthTokenEndpointAuthCreateParamNone) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Sends the client ID and secret using HTTP Basic authentication.
//
// The properties ClientSecret, Type are required.
type McpOAuthTokenEndpointAuthCreateParamClientSecretBasic struct {
	// The OAuth client secret to store. Never returned in credential resources.
	ClientSecret string `json:"client_secret" api:"required"`
	// The type of the object. Always `client_secret_basic`.
	//
	// This field can be elided, and will marshal its zero value as
	// "client_secret_basic".
	Type constant.ClientSecretBasic `json:"type" default:"client_secret_basic"`
	paramObj
}

func (r McpOAuthTokenEndpointAuthCreateParamClientSecretBasic) MarshalJSON() (data []byte, err error) {
	type shadow McpOAuthTokenEndpointAuthCreateParamClientSecretBasic
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *McpOAuthTokenEndpointAuthCreateParamClientSecretBasic) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Sends the client ID and secret in the token request body.
//
// The properties ClientSecret, Type are required.
type McpOAuthTokenEndpointAuthCreateParamClientSecretPost struct {
	// The OAuth client secret to store. Never returned in credential resources.
	ClientSecret string `json:"client_secret" api:"required"`
	// The type of the object. Always `client_secret_post`.
	//
	// This field can be elided, and will marshal its zero value as
	// "client_secret_post".
	Type constant.ClientSecretPost `json:"type" default:"client_secret_post"`
	paramObj
}

func (r McpOAuthTokenEndpointAuthCreateParamClientSecretPost) MarshalJSON() (data []byte, err error) {
	type shadow McpOAuthTokenEndpointAuthCreateParamClientSecretPost
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *McpOAuthTokenEndpointAuthCreateParamClientSecretPost) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type McpOAuthTokenEndpointAuthRotateParamUnion struct {
	OfParamClientSecretBasic *McpOAuthTokenEndpointAuthRotateParamClientSecretBasic `json:",omitzero,inline"`
	OfParamClientSecretPost  *McpOAuthTokenEndpointAuthRotateParamClientSecretPost  `json:",omitzero,inline"`
	paramUnion
}

func (u McpOAuthTokenEndpointAuthRotateParamUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfParamClientSecretBasic, u.OfParamClientSecretPost)
}
func (u *McpOAuthTokenEndpointAuthRotateParamUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u McpOAuthTokenEndpointAuthRotateParamUnion) GetType() *string {
	if vt := u.OfParamClientSecretBasic; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfParamClientSecretPost; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u McpOAuthTokenEndpointAuthRotateParamUnion) GetClientSecret() *string {
	if vt := u.OfParamClientSecretBasic; vt != nil && vt.ClientSecret.Valid() {
		return &vt.ClientSecret.Value
	} else if vt := u.OfParamClientSecretPost; vt != nil && vt.ClientSecret.Valid() {
		return &vt.ClientSecret.Value
	}
	return nil
}

func init() {
	apijson.RegisterUnion[McpOAuthTokenEndpointAuthRotateParamUnion](
		"type",
		apijson.Discriminator[McpOAuthTokenEndpointAuthRotateParamClientSecretBasic]("client_secret_basic"),
		apijson.Discriminator[McpOAuthTokenEndpointAuthRotateParamClientSecretPost]("client_secret_post"),
	)
}

// Updates credentials sent using HTTP Basic authentication.
//
// The property Type is required.
type McpOAuthTokenEndpointAuthRotateParamClientSecretBasic struct {
	// The replacement OAuth client secret. Omit or pass `null` to keep the stored
	// secret. This secret is never returned in resources.
	ClientSecret param.Opt[string] `json:"client_secret,omitzero"`
	// The type of the object. Always `client_secret_basic`.
	//
	// This field can be elided, and will marshal its zero value as
	// "client_secret_basic".
	Type constant.ClientSecretBasic `json:"type" default:"client_secret_basic"`
	paramObj
}

func (r McpOAuthTokenEndpointAuthRotateParamClientSecretBasic) MarshalJSON() (data []byte, err error) {
	type shadow McpOAuthTokenEndpointAuthRotateParamClientSecretBasic
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *McpOAuthTokenEndpointAuthRotateParamClientSecretBasic) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Updates credentials sent in the token request body.
//
// The property Type is required.
type McpOAuthTokenEndpointAuthRotateParamClientSecretPost struct {
	// The replacement OAuth client secret. Omit or pass `null` to keep the stored
	// secret. This secret is never returned in resources.
	ClientSecret param.Opt[string] `json:"client_secret,omitzero"`
	// The type of the object. Always `client_secret_post`.
	//
	// This field can be elided, and will marshal its zero value as
	// "client_secret_post".
	Type constant.ClientSecretPost `json:"type" default:"client_secret_post"`
	paramObj
}

func (r McpOAuthTokenEndpointAuthRotateParamClientSecretPost) MarshalJSON() (data []byte, err error) {
	type shadow McpOAuthTokenEndpointAuthRotateParamClientSecretPost
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *McpOAuthTokenEndpointAuthRotateParamClientSecretPost) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaAgentVaultCredentialNewParams struct {
	// The authentication method and write-only secret values to store.
	Auth CredentialAuthCreateParamUnion `json:"auth,omitzero" api:"required"`
	// The name is trimmed before storage. It must contain 1 to 256 UTF-8 bytes after
	// trimming.
	Name string `json:"name" api:"required"`
	paramObj
}

func (r BetaAgentVaultCredentialNewParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaAgentVaultCredentialNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaAgentVaultCredentialNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaAgentVaultCredentialUpdateParams struct {
	// Replacement values for the credential's existing authentication method.
	Auth CredentialAuthRotateParamUnion `json:"auth,omitzero" api:"required"`
	paramObj
}

func (r BetaAgentVaultCredentialUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaAgentVaultCredentialUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaAgentVaultCredentialUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaAgentVaultCredentialListParams struct {
	// Return resources after this resource ID in the selected order.
	After param.Opt[string] `query:"after,omitzero" json:"-"`
	// The maximum number of resources to return. Defaults to 20. Values are clamped
	// between 1 and 100.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Sort order by the `created_at` timestamp. Use `asc` for ascending order or
	// `desc` for descending order. Defaults to `desc`.
	//
	// Any of "asc", "desc".
	Order BetaAgentVaultCredentialListParamsOrder `query:"order,omitzero" json:"-"`
	// Filter by one status or a list, such as `status=active` or
	// `status[]=active&status[]=archived`. Both statuses are included by default.
	Status VaultStatusFilterUnionParam `query:"status,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaAgentVaultCredentialListParams]'s query parameters as
// `url.Values`.
func (r BetaAgentVaultCredentialListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Sort order by the `created_at` timestamp. Use `asc` for ascending order or
// `desc` for descending order. Defaults to `desc`.
type BetaAgentVaultCredentialListParamsOrder string

const (
	BetaAgentVaultCredentialListParamsOrderAsc  BetaAgentVaultCredentialListParamsOrder = "asc"
	BetaAgentVaultCredentialListParamsOrderDesc BetaAgentVaultCredentialListParamsOrder = "desc"
)

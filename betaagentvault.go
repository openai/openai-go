// File generated from our OpenAPI spec by Castiron. See CONTRIBUTING.md for details.

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

// BetaAgentVaultService contains methods and other services that help with
// interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaAgentVaultService] method instead.
type BetaAgentVaultService struct {
	Options     []option.RequestOption
	Credentials BetaAgentVaultCredentialService
}

// NewBetaAgentVaultService generates a new service that applies the given options
// to each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewBetaAgentVaultService(opts ...option.RequestOption) (r BetaAgentVaultService) {
	r = BetaAgentVaultService{}
	r.Options = requestconfig.InheritedOptions(opts...)
	r.Credentials = NewBetaAgentVaultCredentialService(opts...)
	return
}

// Creates a vault for the current project. See
// [vaults](https://developers.openai.com/api/docs/guides/agents-api/tools/vaults).
func (r *BetaAgentVaultService) New(ctx context.Context, body BetaAgentVaultNewParams, opts ...option.RequestOption) (res *Vault, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "agents=v1")}, opts...)
	path := "vaults"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Retrieves a vault by its ID. See
// [vaults](https://developers.openai.com/api/docs/guides/agents-api/tools/vaults).
func (r *BetaAgentVaultService) Get(ctx context.Context, vaultID string, opts ...option.RequestOption) (res *Vault, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "agents=v1")}, opts...)
	if vaultID == "" {
		err = errors.New("missing required vault_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("vaults/%s", vaultID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Lists vaults using ID-based pagination. See
// [vaults](https://developers.openai.com/api/docs/guides/agents-api/tools/vaults).
func (r *BetaAgentVaultService) List(ctx context.Context, query BetaAgentVaultListParams, opts ...option.RequestOption) (res *pagination.CursorPage[Vault], err error) {
	var raw *http.Response
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "agents=v1"), option.WithResponseInto(&raw)}, opts...)
	path := "vaults"
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

// Lists vaults using ID-based pagination. See
// [vaults](https://developers.openai.com/api/docs/guides/agents-api/tools/vaults).
func (r *BetaAgentVaultService) ListAutoPaging(ctx context.Context, query BetaAgentVaultListParams, opts ...option.RequestOption) *pagination.CursorPageAutoPager[Vault] {
	return pagination.NewCursorPageAutoPager(r.List(ctx, query, opts...))
}

// Deletes a vault and all its credentials. See
// [vaults](https://developers.openai.com/api/docs/guides/agents-api/tools/vaults).
func (r *BetaAgentVaultService) Delete(ctx context.Context, vaultID string, opts ...option.RequestOption) (res *VaultDeleted, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "agents=v1")}, opts...)
	if vaultID == "" {
		err = errors.New("missing required vault_id parameter")
		return nil, err
	}
	path := requestconfig.FormatPath("vaults/%s", vaultID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// A collection of credentials that agent tools can use to authenticate to MCP
// servers.
type Vault struct {
	// The ID of the vault.
	ID string `json:"id" api:"required"`
	// The Unix timestamp, in seconds, when the vault was created.
	CreatedAt int64 `json:"created_at" api:"required"`
	// Key-value pairs associated with the vault, such as an application or team
	// identifier.
	Metadata map[string]string `json:"metadata" api:"required"`
	// The human-readable name of the vault, if set.
	Name string `json:"name" api:"required"`
	// The object type. Always `vault`.
	Object constant.Vault `json:"object" default:"vault"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		Metadata    respjson.Field
		Name        respjson.Field
		Object      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r Vault) RawJSON() string { return r.JSON.raw }
func (r *Vault) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Confirmation that a vault was deleted.
type VaultDeleted struct {
	// The ID of the deleted vault.
	ID string `json:"id" api:"required"`
	// Whether the resource was deleted. Always `true`.
	Deleted bool `json:"deleted" api:"required"`
	// The object type. Always `vault.deleted`.
	Object constant.VaultDeleted `json:"object" default:"vault.deleted"`
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
func (r VaultDeleted) RawJSON() string { return r.JSON.raw }
func (r *VaultDeleted) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Whether a vault or credential is active or archived.
type VaultStatus string

const (
	VaultStatusActive   VaultStatus = "active"
	VaultStatusArchived VaultStatus = "archived"
)

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type VaultStatusFilterUnionParam struct {
	// Check if union is this variant with !param.IsOmitted(union.OfStatus)
	OfStatus          param.Opt[VaultStatus] `json:",omitzero,inline"`
	OfArrayOfStatuses []VaultStatus          `json:",omitzero,inline"`
	paramUnion
}

func (u VaultStatusFilterUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfStatus, u.OfArrayOfStatuses)
}
func (u *VaultStatusFilterUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

type BetaAgentVaultNewParams struct {
	// The name is trimmed before storage. It must contain 1 to 256 UTF-8 bytes after
	// trimming.
	Name param.Opt[string] `json:"name,omitzero"`
	// Key-value pairs to associate with the vault, such as an application or team
	// identifier.
	Metadata map[string]string `json:"metadata,omitzero"`
	paramObj
}

func (r BetaAgentVaultNewParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaAgentVaultNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaAgentVaultNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaAgentVaultListParams struct {
	// Return resources after this resource ID in the selected order.
	After param.Opt[string] `query:"after,omitzero" json:"-"`
	// The maximum number of resources to return. Defaults to 20. Values are clamped
	// between 1 and 100.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Sort order by the `created_at` timestamp. Use `asc` for ascending order or
	// `desc` for descending order. Defaults to `desc`.
	//
	// Any of "asc", "desc".
	Order BetaAgentVaultListParamsOrder `query:"order,omitzero" json:"-"`
	// Filter by one status or a list, such as `status=active` or
	// `status[]=active&status[]=archived`. Both statuses are included by default.
	Status VaultStatusFilterUnionParam `query:"status,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaAgentVaultListParams]'s query parameters as
// `url.Values`.
func (r BetaAgentVaultListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Sort order by the `created_at` timestamp. Use `asc` for ascending order or
// `desc` for descending order. Defaults to `desc`.
type BetaAgentVaultListParamsOrder string

const (
	BetaAgentVaultListParamsOrderAsc  BetaAgentVaultListParamsOrder = "asc"
	BetaAgentVaultListParamsOrderDesc BetaAgentVaultListParamsOrder = "desc"
)

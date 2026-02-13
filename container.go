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
	"github.com/openai/openai-go/v3/responses"
)

// ContainerService contains methods and other services that help with interacting
// with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewContainerService] method instead.
type ContainerService struct {
	Options []option.RequestOption
	Files   ContainerFileService
}

// NewContainerService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewContainerService(opts ...option.RequestOption) (r ContainerService) {
	r = ContainerService{}
	r.Options = opts
	r.Files = NewContainerFileService(opts...)
	return
}

// Create Container
func (r *ContainerService) New(ctx context.Context, body ContainerNewParams, opts ...option.RequestOption) (res *ContainerNewResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "containers"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Retrieve Container
func (r *ContainerService) Get(ctx context.Context, containerID string, opts ...option.RequestOption) (res *ContainerGetResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if containerID == "" {
		err = errors.New("missing required container_id parameter")
		return
	}
	path := fmt.Sprintf("containers/%s", containerID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// List Containers
func (r *ContainerService) List(ctx context.Context, query ContainerListParams, opts ...option.RequestOption) (res *pagination.CursorPage[ContainerListResponse], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "containers"
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

// List Containers
func (r *ContainerService) ListAutoPaging(ctx context.Context, query ContainerListParams, opts ...option.RequestOption) *pagination.CursorPageAutoPager[ContainerListResponse] {
	return pagination.NewCursorPageAutoPager(r.List(ctx, query, opts...))
}

// Delete Container
func (r *ContainerService) Delete(ctx context.Context, containerID string, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if containerID == "" {
		err = errors.New("missing required container_id parameter")
		return
	}
	path := fmt.Sprintf("containers/%s", containerID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, nil, opts...)
	return
}

type ContainerNewResponse struct {
	// Unique identifier for the container.
	ID string `json:"id,required"`
	// Unix timestamp (in seconds) when the container was created.
	CreatedAt int64 `json:"created_at,required"`
	// Name of the container.
	Name string `json:"name,required"`
	// The type of this object.
	Object string `json:"object,required"`
	// Status of the container (e.g., active, deleted).
	Status string `json:"status,required"`
	// The container will expire after this time period. The anchor is the reference
	// point for the expiration. The minutes is the number of minutes after the anchor
	// before the container expires.
	ExpiresAfter ContainerNewResponseExpiresAfter `json:"expires_after"`
	// Unix timestamp (in seconds) when the container was last active.
	LastActiveAt int64 `json:"last_active_at"`
	// The memory limit configured for the container.
	//
	// Any of "1g", "4g", "16g", "64g".
	MemoryLimit ContainerNewResponseMemoryLimit `json:"memory_limit"`
	// Network access policy for the container.
	NetworkPolicy ContainerNewResponseNetworkPolicy `json:"network_policy"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID            respjson.Field
		CreatedAt     respjson.Field
		Name          respjson.Field
		Object        respjson.Field
		Status        respjson.Field
		ExpiresAfter  respjson.Field
		LastActiveAt  respjson.Field
		MemoryLimit   respjson.Field
		NetworkPolicy respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ContainerNewResponse) RawJSON() string { return r.JSON.raw }
func (r *ContainerNewResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The container will expire after this time period. The anchor is the reference
// point for the expiration. The minutes is the number of minutes after the anchor
// before the container expires.
type ContainerNewResponseExpiresAfter struct {
	// The reference point for the expiration.
	//
	// Any of "last_active_at".
	Anchor string `json:"anchor"`
	// The number of minutes after the anchor before the container expires.
	Minutes int64 `json:"minutes"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Anchor      respjson.Field
		Minutes     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ContainerNewResponseExpiresAfter) RawJSON() string { return r.JSON.raw }
func (r *ContainerNewResponseExpiresAfter) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The memory limit configured for the container.
type ContainerNewResponseMemoryLimit string

const (
	ContainerNewResponseMemoryLimit1g  ContainerNewResponseMemoryLimit = "1g"
	ContainerNewResponseMemoryLimit4g  ContainerNewResponseMemoryLimit = "4g"
	ContainerNewResponseMemoryLimit16g ContainerNewResponseMemoryLimit = "16g"
	ContainerNewResponseMemoryLimit64g ContainerNewResponseMemoryLimit = "64g"
)

// Network access policy for the container.
type ContainerNewResponseNetworkPolicy struct {
	// The network policy mode.
	//
	// Any of "allowlist", "disabled".
	Type string `json:"type,required"`
	// Allowed outbound domains when `type` is `allowlist`.
	AllowedDomains []string `json:"allowed_domains"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type           respjson.Field
		AllowedDomains respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ContainerNewResponseNetworkPolicy) RawJSON() string { return r.JSON.raw }
func (r *ContainerNewResponseNetworkPolicy) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ContainerGetResponse struct {
	// Unique identifier for the container.
	ID string `json:"id,required"`
	// Unix timestamp (in seconds) when the container was created.
	CreatedAt int64 `json:"created_at,required"`
	// Name of the container.
	Name string `json:"name,required"`
	// The type of this object.
	Object string `json:"object,required"`
	// Status of the container (e.g., active, deleted).
	Status string `json:"status,required"`
	// The container will expire after this time period. The anchor is the reference
	// point for the expiration. The minutes is the number of minutes after the anchor
	// before the container expires.
	ExpiresAfter ContainerGetResponseExpiresAfter `json:"expires_after"`
	// Unix timestamp (in seconds) when the container was last active.
	LastActiveAt int64 `json:"last_active_at"`
	// The memory limit configured for the container.
	//
	// Any of "1g", "4g", "16g", "64g".
	MemoryLimit ContainerGetResponseMemoryLimit `json:"memory_limit"`
	// Network access policy for the container.
	NetworkPolicy ContainerGetResponseNetworkPolicy `json:"network_policy"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID            respjson.Field
		CreatedAt     respjson.Field
		Name          respjson.Field
		Object        respjson.Field
		Status        respjson.Field
		ExpiresAfter  respjson.Field
		LastActiveAt  respjson.Field
		MemoryLimit   respjson.Field
		NetworkPolicy respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ContainerGetResponse) RawJSON() string { return r.JSON.raw }
func (r *ContainerGetResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The container will expire after this time period. The anchor is the reference
// point for the expiration. The minutes is the number of minutes after the anchor
// before the container expires.
type ContainerGetResponseExpiresAfter struct {
	// The reference point for the expiration.
	//
	// Any of "last_active_at".
	Anchor string `json:"anchor"`
	// The number of minutes after the anchor before the container expires.
	Minutes int64 `json:"minutes"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Anchor      respjson.Field
		Minutes     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ContainerGetResponseExpiresAfter) RawJSON() string { return r.JSON.raw }
func (r *ContainerGetResponseExpiresAfter) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The memory limit configured for the container.
type ContainerGetResponseMemoryLimit string

const (
	ContainerGetResponseMemoryLimit1g  ContainerGetResponseMemoryLimit = "1g"
	ContainerGetResponseMemoryLimit4g  ContainerGetResponseMemoryLimit = "4g"
	ContainerGetResponseMemoryLimit16g ContainerGetResponseMemoryLimit = "16g"
	ContainerGetResponseMemoryLimit64g ContainerGetResponseMemoryLimit = "64g"
)

// Network access policy for the container.
type ContainerGetResponseNetworkPolicy struct {
	// The network policy mode.
	//
	// Any of "allowlist", "disabled".
	Type string `json:"type,required"`
	// Allowed outbound domains when `type` is `allowlist`.
	AllowedDomains []string `json:"allowed_domains"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type           respjson.Field
		AllowedDomains respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ContainerGetResponseNetworkPolicy) RawJSON() string { return r.JSON.raw }
func (r *ContainerGetResponseNetworkPolicy) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ContainerListResponse struct {
	// Unique identifier for the container.
	ID string `json:"id,required"`
	// Unix timestamp (in seconds) when the container was created.
	CreatedAt int64 `json:"created_at,required"`
	// Name of the container.
	Name string `json:"name,required"`
	// The type of this object.
	Object string `json:"object,required"`
	// Status of the container (e.g., active, deleted).
	Status string `json:"status,required"`
	// The container will expire after this time period. The anchor is the reference
	// point for the expiration. The minutes is the number of minutes after the anchor
	// before the container expires.
	ExpiresAfter ContainerListResponseExpiresAfter `json:"expires_after"`
	// Unix timestamp (in seconds) when the container was last active.
	LastActiveAt int64 `json:"last_active_at"`
	// The memory limit configured for the container.
	//
	// Any of "1g", "4g", "16g", "64g".
	MemoryLimit ContainerListResponseMemoryLimit `json:"memory_limit"`
	// Network access policy for the container.
	NetworkPolicy ContainerListResponseNetworkPolicy `json:"network_policy"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID            respjson.Field
		CreatedAt     respjson.Field
		Name          respjson.Field
		Object        respjson.Field
		Status        respjson.Field
		ExpiresAfter  respjson.Field
		LastActiveAt  respjson.Field
		MemoryLimit   respjson.Field
		NetworkPolicy respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ContainerListResponse) RawJSON() string { return r.JSON.raw }
func (r *ContainerListResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The container will expire after this time period. The anchor is the reference
// point for the expiration. The minutes is the number of minutes after the anchor
// before the container expires.
type ContainerListResponseExpiresAfter struct {
	// The reference point for the expiration.
	//
	// Any of "last_active_at".
	Anchor string `json:"anchor"`
	// The number of minutes after the anchor before the container expires.
	Minutes int64 `json:"minutes"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Anchor      respjson.Field
		Minutes     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ContainerListResponseExpiresAfter) RawJSON() string { return r.JSON.raw }
func (r *ContainerListResponseExpiresAfter) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The memory limit configured for the container.
type ContainerListResponseMemoryLimit string

const (
	ContainerListResponseMemoryLimit1g  ContainerListResponseMemoryLimit = "1g"
	ContainerListResponseMemoryLimit4g  ContainerListResponseMemoryLimit = "4g"
	ContainerListResponseMemoryLimit16g ContainerListResponseMemoryLimit = "16g"
	ContainerListResponseMemoryLimit64g ContainerListResponseMemoryLimit = "64g"
)

// Network access policy for the container.
type ContainerListResponseNetworkPolicy struct {
	// The network policy mode.
	//
	// Any of "allowlist", "disabled".
	Type string `json:"type,required"`
	// Allowed outbound domains when `type` is `allowlist`.
	AllowedDomains []string `json:"allowed_domains"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type           respjson.Field
		AllowedDomains respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ContainerListResponseNetworkPolicy) RawJSON() string { return r.JSON.raw }
func (r *ContainerListResponseNetworkPolicy) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ContainerNewParams struct {
	// Name of the container to create.
	Name string `json:"name,required"`
	// Container expiration time in seconds relative to the 'anchor' time.
	ExpiresAfter ContainerNewParamsExpiresAfter `json:"expires_after,omitzero"`
	// IDs of files to copy to the container.
	FileIDs []string `json:"file_ids,omitzero"`
	// Optional memory limit for the container. Defaults to "1g".
	//
	// Any of "1g", "4g", "16g", "64g".
	MemoryLimit ContainerNewParamsMemoryLimit `json:"memory_limit,omitzero"`
	// Network access policy for the container.
	NetworkPolicy ContainerNewParamsNetworkPolicyUnion `json:"network_policy,omitzero"`
	// An optional list of skills referenced by id or inline data.
	Skills []ContainerNewParamsSkillUnion `json:"skills,omitzero"`
	paramObj
}

func (r ContainerNewParams) MarshalJSON() (data []byte, err error) {
	type shadow ContainerNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ContainerNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Container expiration time in seconds relative to the 'anchor' time.
//
// The properties Anchor, Minutes are required.
type ContainerNewParamsExpiresAfter struct {
	// Time anchor for the expiration time. Currently only 'last_active_at' is
	// supported.
	//
	// Any of "last_active_at".
	Anchor  string `json:"anchor,omitzero,required"`
	Minutes int64  `json:"minutes,required"`
	paramObj
}

func (r ContainerNewParamsExpiresAfter) MarshalJSON() (data []byte, err error) {
	type shadow ContainerNewParamsExpiresAfter
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ContainerNewParamsExpiresAfter) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[ContainerNewParamsExpiresAfter](
		"anchor", "last_active_at",
	)
}

// Optional memory limit for the container. Defaults to "1g".
type ContainerNewParamsMemoryLimit string

const (
	ContainerNewParamsMemoryLimit1g  ContainerNewParamsMemoryLimit = "1g"
	ContainerNewParamsMemoryLimit4g  ContainerNewParamsMemoryLimit = "4g"
	ContainerNewParamsMemoryLimit16g ContainerNewParamsMemoryLimit = "16g"
	ContainerNewParamsMemoryLimit64g ContainerNewParamsMemoryLimit = "64g"
)

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ContainerNewParamsNetworkPolicyUnion struct {
	OfDisabled  *responses.ContainerNetworkPolicyDisabledParam  `json:",omitzero,inline"`
	OfAllowlist *responses.ContainerNetworkPolicyAllowlistParam `json:",omitzero,inline"`
	paramUnion
}

func (u ContainerNewParamsNetworkPolicyUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfDisabled, u.OfAllowlist)
}
func (u *ContainerNewParamsNetworkPolicyUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *ContainerNewParamsNetworkPolicyUnion) asAny() any {
	if !param.IsOmitted(u.OfDisabled) {
		return u.OfDisabled
	} else if !param.IsOmitted(u.OfAllowlist) {
		return u.OfAllowlist
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ContainerNewParamsNetworkPolicyUnion) GetAllowedDomains() []string {
	if vt := u.OfAllowlist; vt != nil {
		return vt.AllowedDomains
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ContainerNewParamsNetworkPolicyUnion) GetDomainSecrets() []responses.ContainerNetworkPolicyDomainSecretParam {
	if vt := u.OfAllowlist; vt != nil {
		return vt.DomainSecrets
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ContainerNewParamsNetworkPolicyUnion) GetType() *string {
	if vt := u.OfDisabled; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfAllowlist; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[ContainerNewParamsNetworkPolicyUnion](
		"type",
		apijson.Discriminator[responses.ContainerNetworkPolicyDisabledParam]("disabled"),
		apijson.Discriminator[responses.ContainerNetworkPolicyAllowlistParam]("allowlist"),
	)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ContainerNewParamsSkillUnion struct {
	OfSkillReference *responses.SkillReferenceParam `json:",omitzero,inline"`
	OfInline         *responses.InlineSkillParam    `json:",omitzero,inline"`
	paramUnion
}

func (u ContainerNewParamsSkillUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfSkillReference, u.OfInline)
}
func (u *ContainerNewParamsSkillUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *ContainerNewParamsSkillUnion) asAny() any {
	if !param.IsOmitted(u.OfSkillReference) {
		return u.OfSkillReference
	} else if !param.IsOmitted(u.OfInline) {
		return u.OfInline
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ContainerNewParamsSkillUnion) GetSkillID() *string {
	if vt := u.OfSkillReference; vt != nil {
		return &vt.SkillID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ContainerNewParamsSkillUnion) GetVersion() *string {
	if vt := u.OfSkillReference; vt != nil && vt.Version.Valid() {
		return &vt.Version.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ContainerNewParamsSkillUnion) GetDescription() *string {
	if vt := u.OfInline; vt != nil {
		return &vt.Description
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ContainerNewParamsSkillUnion) GetName() *string {
	if vt := u.OfInline; vt != nil {
		return &vt.Name
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ContainerNewParamsSkillUnion) GetSource() *responses.InlineSkillSourceParam {
	if vt := u.OfInline; vt != nil {
		return &vt.Source
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ContainerNewParamsSkillUnion) GetType() *string {
	if vt := u.OfSkillReference; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfInline; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[ContainerNewParamsSkillUnion](
		"type",
		apijson.Discriminator[responses.SkillReferenceParam]("skill_reference"),
		apijson.Discriminator[responses.InlineSkillParam]("inline"),
	)
}

type ContainerListParams struct {
	// A cursor for use in pagination. `after` is an object ID that defines your place
	// in the list. For instance, if you make a list request and receive 100 objects,
	// ending with obj_foo, your subsequent call can include after=obj_foo in order to
	// fetch the next page of the list.
	After param.Opt[string] `query:"after,omitzero" json:"-"`
	// A limit on the number of objects to be returned. Limit can range between 1 and
	// 100, and the default is 20.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Filter results by container name.
	Name param.Opt[string] `query:"name,omitzero" json:"-"`
	// Sort order by the `created_at` timestamp of the objects. `asc` for ascending
	// order and `desc` for descending order.
	//
	// Any of "asc", "desc".
	Order ContainerListParamsOrder `query:"order,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [ContainerListParams]'s query parameters as `url.Values`.
func (r ContainerListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Sort order by the `created_at` timestamp of the objects. `asc` for ascending
// order and `desc` for descending order.
type ContainerListParamsOrder string

const (
	ContainerListParamsOrderAsc  ContainerListParamsOrder = "asc"
	ContainerListParamsOrderDesc ContainerListParamsOrder = "desc"
)

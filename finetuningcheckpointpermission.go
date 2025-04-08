// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/openai/openai-go/internal/apijson"
	"github.com/openai/openai-go/internal/apiquery"
	"github.com/openai/openai-go/internal/requestconfig"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/pagination"
	"github.com/openai/openai-go/packages/param"
	"github.com/openai/openai-go/packages/resp"
	"github.com/openai/openai-go/shared/constant"
)

// FineTuningCheckpointPermissionService contains methods and other services that
// help with interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewFineTuningCheckpointPermissionService] method instead.
type FineTuningCheckpointPermissionService struct {
	Options []option.RequestOption
}

// NewFineTuningCheckpointPermissionService generates a new service that applies
// the given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewFineTuningCheckpointPermissionService(opts ...option.RequestOption) (r FineTuningCheckpointPermissionService) {
	r = FineTuningCheckpointPermissionService{}
	r.Options = opts
	return
}

// **NOTE:** Calling this endpoint requires an [admin API key](../admin-api-keys).
//
// This enables organization owners to share fine-tuned models with other projects
// in their organization.
func (r *FineTuningCheckpointPermissionService) New(ctx context.Context, fineTunedModelCheckpoint string, body FineTuningCheckpointPermissionNewParams, opts ...option.RequestOption) (res *pagination.Page[FineTuningCheckpointPermissionNewResponse], err error) {
	var raw *http.Response
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	if fineTunedModelCheckpoint == "" {
		err = errors.New("missing required fine_tuned_model_checkpoint parameter")
		return
	}
	path := fmt.Sprintf("fine_tuning/checkpoints/%s/permissions", fineTunedModelCheckpoint)
	cfg, err := requestconfig.NewRequestConfig(ctx, http.MethodPost, path, body, &res, opts...)
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

// **NOTE:** Calling this endpoint requires an [admin API key](../admin-api-keys).
//
// This enables organization owners to share fine-tuned models with other projects
// in their organization.
func (r *FineTuningCheckpointPermissionService) NewAutoPaging(ctx context.Context, fineTunedModelCheckpoint string, body FineTuningCheckpointPermissionNewParams, opts ...option.RequestOption) *pagination.PageAutoPager[FineTuningCheckpointPermissionNewResponse] {
	return pagination.NewPageAutoPager(r.New(ctx, fineTunedModelCheckpoint, body, opts...))
}

// **NOTE:** This endpoint requires an [admin API key](../admin-api-keys).
//
// Organization owners can use this endpoint to view all permissions for a
// fine-tuned model checkpoint.
func (r *FineTuningCheckpointPermissionService) Get(ctx context.Context, fineTunedModelCheckpoint string, query FineTuningCheckpointPermissionGetParams, opts ...option.RequestOption) (res *FineTuningCheckpointPermissionGetResponse, err error) {
	opts = append(r.Options[:], opts...)
	if fineTunedModelCheckpoint == "" {
		err = errors.New("missing required fine_tuned_model_checkpoint parameter")
		return
	}
	path := fmt.Sprintf("fine_tuning/checkpoints/%s/permissions", fineTunedModelCheckpoint)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return
}

// **NOTE:** This endpoint requires an [admin API key](../admin-api-keys).
//
// Organization owners can use this endpoint to delete a permission for a
// fine-tuned model checkpoint.
func (r *FineTuningCheckpointPermissionService) Delete(ctx context.Context, fineTunedModelCheckpoint string, opts ...option.RequestOption) (res *FineTuningCheckpointPermissionDeleteResponse, err error) {
	opts = append(r.Options[:], opts...)
	if fineTunedModelCheckpoint == "" {
		err = errors.New("missing required fine_tuned_model_checkpoint parameter")
		return
	}
	path := fmt.Sprintf("fine_tuning/checkpoints/%s/permissions", fineTunedModelCheckpoint)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return
}

// The `checkpoint.permission` object represents a permission for a fine-tuned
// model checkpoint.
type FineTuningCheckpointPermissionNewResponse struct {
	// The permission identifier, which can be referenced in the API endpoints.
	ID string `json:"id,required"`
	// The Unix timestamp (in seconds) for when the permission was created.
	CreatedAt int64 `json:"created_at,required"`
	// The object type, which is always "checkpoint.permission".
	Object constant.CheckpointPermission `json:"object,required"`
	// The project identifier that the permission is for.
	ProjectID string `json:"project_id,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		ID          resp.Field
		CreatedAt   resp.Field
		Object      resp.Field
		ProjectID   resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r FineTuningCheckpointPermissionNewResponse) RawJSON() string { return r.JSON.raw }
func (r *FineTuningCheckpointPermissionNewResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type FineTuningCheckpointPermissionGetResponse struct {
	Data    []FineTuningCheckpointPermissionGetResponseData `json:"data,required"`
	HasMore bool                                            `json:"has_more,required"`
	Object  constant.List                                   `json:"object,required"`
	FirstID string                                          `json:"first_id,nullable"`
	LastID  string                                          `json:"last_id,nullable"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		Data        resp.Field
		HasMore     resp.Field
		Object      resp.Field
		FirstID     resp.Field
		LastID      resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r FineTuningCheckpointPermissionGetResponse) RawJSON() string { return r.JSON.raw }
func (r *FineTuningCheckpointPermissionGetResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The `checkpoint.permission` object represents a permission for a fine-tuned
// model checkpoint.
type FineTuningCheckpointPermissionGetResponseData struct {
	// The permission identifier, which can be referenced in the API endpoints.
	ID string `json:"id,required"`
	// The Unix timestamp (in seconds) for when the permission was created.
	CreatedAt int64 `json:"created_at,required"`
	// The object type, which is always "checkpoint.permission".
	Object constant.CheckpointPermission `json:"object,required"`
	// The project identifier that the permission is for.
	ProjectID string `json:"project_id,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		ID          resp.Field
		CreatedAt   resp.Field
		Object      resp.Field
		ProjectID   resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r FineTuningCheckpointPermissionGetResponseData) RawJSON() string { return r.JSON.raw }
func (r *FineTuningCheckpointPermissionGetResponseData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type FineTuningCheckpointPermissionDeleteResponse struct {
	// The ID of the fine-tuned model checkpoint permission that was deleted.
	ID string `json:"id,required"`
	// Whether the fine-tuned model checkpoint permission was successfully deleted.
	Deleted bool `json:"deleted,required"`
	// The object type, which is always "checkpoint.permission".
	Object constant.CheckpointPermission `json:"object,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		ID          resp.Field
		Deleted     resp.Field
		Object      resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r FineTuningCheckpointPermissionDeleteResponse) RawJSON() string { return r.JSON.raw }
func (r *FineTuningCheckpointPermissionDeleteResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type FineTuningCheckpointPermissionNewParams struct {
	// The project identifiers to grant access to.
	ProjectIDs []string `json:"project_ids,omitzero,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f FineTuningCheckpointPermissionNewParams) IsPresent() bool {
	return !param.IsOmitted(f) && !f.IsNull()
}

func (r FineTuningCheckpointPermissionNewParams) MarshalJSON() (data []byte, err error) {
	type shadow FineTuningCheckpointPermissionNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}

type FineTuningCheckpointPermissionGetParams struct {
	// Identifier for the last permission ID from the previous pagination request.
	After param.Opt[string] `query:"after,omitzero" json:"-"`
	// Number of permissions to retrieve.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// The ID of the project to get permissions for.
	ProjectID param.Opt[string] `query:"project_id,omitzero" json:"-"`
	// The order in which to retrieve permissions.
	//
	// Any of "ascending", "descending".
	Order FineTuningCheckpointPermissionGetParamsOrder `query:"order,omitzero" json:"-"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f FineTuningCheckpointPermissionGetParams) IsPresent() bool {
	return !param.IsOmitted(f) && !f.IsNull()
}

// URLQuery serializes [FineTuningCheckpointPermissionGetParams]'s query parameters
// as `url.Values`.
func (r FineTuningCheckpointPermissionGetParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// The order in which to retrieve permissions.
type FineTuningCheckpointPermissionGetParamsOrder string

const (
	FineTuningCheckpointPermissionGetParamsOrderAscending  FineTuningCheckpointPermissionGetParamsOrder = "ascending"
	FineTuningCheckpointPermissionGetParamsOrderDescending FineTuningCheckpointPermissionGetParamsOrder = "descending"
)

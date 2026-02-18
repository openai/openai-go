// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"slices"

	"github.com/openai/openai-go/v3/internal/apiform"
	"github.com/openai/openai-go/v3/internal/apijson"
	"github.com/openai/openai-go/v3/internal/apiquery"
	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/pagination"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/packages/respjson"
	"github.com/openai/openai-go/v3/shared/constant"
)

// SkillVersionService contains methods and other services that help with
// interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewSkillVersionService] method instead.
type SkillVersionService struct {
	Options []option.RequestOption
	Content SkillVersionContentService
}

// NewSkillVersionService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewSkillVersionService(opts ...option.RequestOption) (r SkillVersionService) {
	r = SkillVersionService{}
	r.Options = opts
	r.Content = NewSkillVersionContentService(opts...)
	return
}

// Create a new immutable skill version.
func (r *SkillVersionService) New(ctx context.Context, skillID string, body SkillVersionNewParams, opts ...option.RequestOption) (res *SkillVersion, err error) {
	opts = slices.Concat(r.Options, opts)
	if skillID == "" {
		err = errors.New("missing required skill_id parameter")
		return
	}
	path := fmt.Sprintf("skills/%s/versions", skillID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Get a specific skill version.
func (r *SkillVersionService) Get(ctx context.Context, skillID string, version string, opts ...option.RequestOption) (res *SkillVersion, err error) {
	opts = slices.Concat(r.Options, opts)
	if skillID == "" {
		err = errors.New("missing required skill_id parameter")
		return
	}
	if version == "" {
		err = errors.New("missing required version parameter")
		return
	}
	path := fmt.Sprintf("skills/%s/versions/%s", skillID, version)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// List skill versions for a skill.
func (r *SkillVersionService) List(ctx context.Context, skillID string, query SkillVersionListParams, opts ...option.RequestOption) (res *pagination.CursorPage[SkillVersion], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	if skillID == "" {
		err = errors.New("missing required skill_id parameter")
		return
	}
	path := fmt.Sprintf("skills/%s/versions", skillID)
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

// List skill versions for a skill.
func (r *SkillVersionService) ListAutoPaging(ctx context.Context, skillID string, query SkillVersionListParams, opts ...option.RequestOption) *pagination.CursorPageAutoPager[SkillVersion] {
	return pagination.NewCursorPageAutoPager(r.List(ctx, skillID, query, opts...))
}

// Delete a skill version.
func (r *SkillVersionService) Delete(ctx context.Context, skillID string, version string, opts ...option.RequestOption) (res *DeletedSkillVersion, err error) {
	opts = slices.Concat(r.Options, opts)
	if skillID == "" {
		err = errors.New("missing required skill_id parameter")
		return
	}
	if version == "" {
		err = errors.New("missing required version parameter")
		return
	}
	path := fmt.Sprintf("skills/%s/versions/%s", skillID, version)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return
}

type DeletedSkillVersion struct {
	ID      string                       `json:"id,required"`
	Deleted bool                         `json:"deleted,required"`
	Object  constant.SkillVersionDeleted `json:"object,required"`
	// The deleted skill version.
	Version string `json:"version,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Deleted     respjson.Field
		Object      respjson.Field
		Version     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DeletedSkillVersion) RawJSON() string { return r.JSON.raw }
func (r *DeletedSkillVersion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type SkillVersion struct {
	// Unique identifier for the skill version.
	ID string `json:"id,required"`
	// Unix timestamp (seconds) for when the version was created.
	CreatedAt int64 `json:"created_at,required"`
	// Description of the skill version.
	Description string `json:"description,required"`
	// Name of the skill version.
	Name string `json:"name,required"`
	// The object type, which is `skill.version`.
	Object constant.SkillVersion `json:"object,required"`
	// Identifier of the skill for this version.
	SkillID string `json:"skill_id,required"`
	// Version number for this skill.
	Version string `json:"version,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		Description respjson.Field
		Name        respjson.Field
		Object      respjson.Field
		SkillID     respjson.Field
		Version     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r SkillVersion) RawJSON() string { return r.JSON.raw }
func (r *SkillVersion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type SkillVersionList struct {
	// A list of items
	Data []SkillVersion `json:"data,required"`
	// The ID of the first item in the list.
	FirstID string `json:"first_id,required"`
	// Whether there are more items available.
	HasMore bool `json:"has_more,required"`
	// The ID of the last item in the list.
	LastID string `json:"last_id,required"`
	// The type of object returned, must be `list`.
	Object constant.List `json:"object,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Data        respjson.Field
		FirstID     respjson.Field
		HasMore     respjson.Field
		LastID      respjson.Field
		Object      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r SkillVersionList) RawJSON() string { return r.JSON.raw }
func (r *SkillVersionList) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type SkillVersionNewParams struct {
	// Whether to set this version as the default.
	Default param.Opt[bool] `json:"default,omitzero"`
	// Skill files to upload (directory upload) or a single zip file.
	Files SkillVersionNewParamsFilesUnion `json:"files,omitzero" format:"binary"`
	paramObj
}

func (r SkillVersionNewParams) MarshalMultipart() (data []byte, contentType string, err error) {
	buf := bytes.NewBuffer(nil)
	writer := multipart.NewWriter(buf)
	err = apiform.MarshalRoot(r, writer)
	if err == nil {
		err = apiform.WriteExtras(writer, r.ExtraFields())
	}
	if err != nil {
		writer.Close()
		return nil, "", err
	}
	err = writer.Close()
	if err != nil {
		return nil, "", err
	}
	return buf.Bytes(), writer.FormDataContentType(), nil
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type SkillVersionNewParamsFilesUnion struct {
	OfFileArray []io.Reader `json:",omitzero,inline"`
	OfFile      io.Reader   `json:",omitzero,inline"`
	paramUnion
}

func (u SkillVersionNewParamsFilesUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfFileArray, u.OfFile)
}
func (u *SkillVersionNewParamsFilesUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *SkillVersionNewParamsFilesUnion) asAny() any {
	if !param.IsOmitted(u.OfFileArray) {
		return &u.OfFileArray
	} else if !param.IsOmitted(u.OfFile) {
		return &u.OfFile
	}
	return nil
}

type SkillVersionListParams struct {
	// The skill version ID to start after.
	After param.Opt[string] `query:"after,omitzero" json:"-"`
	// Number of versions to retrieve.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Sort order of results by version number.
	//
	// Any of "asc", "desc".
	Order SkillVersionListParamsOrder `query:"order,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [SkillVersionListParams]'s query parameters as `url.Values`.
func (r SkillVersionListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Sort order of results by version number.
type SkillVersionListParamsOrder string

const (
	SkillVersionListParamsOrderAsc  SkillVersionListParamsOrder = "asc"
	SkillVersionListParamsOrderDesc SkillVersionListParamsOrder = "desc"
)

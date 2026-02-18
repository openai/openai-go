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

// SkillService contains methods and other services that help with interacting with
// the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewSkillService] method instead.
type SkillService struct {
	Options  []option.RequestOption
	Content  SkillContentService
	Versions SkillVersionService
}

// NewSkillService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewSkillService(opts ...option.RequestOption) (r SkillService) {
	r = SkillService{}
	r.Options = opts
	r.Content = NewSkillContentService(opts...)
	r.Versions = NewSkillVersionService(opts...)
	return
}

// Create a new skill.
func (r *SkillService) New(ctx context.Context, body SkillNewParams, opts ...option.RequestOption) (res *Skill, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "skills"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Get a skill by its ID.
func (r *SkillService) Get(ctx context.Context, skillID string, opts ...option.RequestOption) (res *Skill, err error) {
	opts = slices.Concat(r.Options, opts)
	if skillID == "" {
		err = errors.New("missing required skill_id parameter")
		return
	}
	path := fmt.Sprintf("skills/%s", skillID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// Update the default version pointer for a skill.
func (r *SkillService) Update(ctx context.Context, skillID string, body SkillUpdateParams, opts ...option.RequestOption) (res *Skill, err error) {
	opts = slices.Concat(r.Options, opts)
	if skillID == "" {
		err = errors.New("missing required skill_id parameter")
		return
	}
	path := fmt.Sprintf("skills/%s", skillID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// List all skills for the current project.
func (r *SkillService) List(ctx context.Context, query SkillListParams, opts ...option.RequestOption) (res *pagination.CursorPage[Skill], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "skills"
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

// List all skills for the current project.
func (r *SkillService) ListAutoPaging(ctx context.Context, query SkillListParams, opts ...option.RequestOption) *pagination.CursorPageAutoPager[Skill] {
	return pagination.NewCursorPageAutoPager(r.List(ctx, query, opts...))
}

// Delete a skill by its ID.
func (r *SkillService) Delete(ctx context.Context, skillID string, opts ...option.RequestOption) (res *DeletedSkill, err error) {
	opts = slices.Concat(r.Options, opts)
	if skillID == "" {
		err = errors.New("missing required skill_id parameter")
		return
	}
	path := fmt.Sprintf("skills/%s", skillID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return
}

type DeletedSkill struct {
	ID      string                `json:"id,required"`
	Deleted bool                  `json:"deleted,required"`
	Object  constant.SkillDeleted `json:"object,required"`
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
func (r DeletedSkill) RawJSON() string { return r.JSON.raw }
func (r *DeletedSkill) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type Skill struct {
	// Unique identifier for the skill.
	ID string `json:"id,required"`
	// Unix timestamp (seconds) for when the skill was created.
	CreatedAt int64 `json:"created_at,required"`
	// Default version for the skill.
	DefaultVersion string `json:"default_version,required"`
	// Description of the skill.
	Description string `json:"description,required"`
	// Latest version for the skill.
	LatestVersion string `json:"latest_version,required"`
	// Name of the skill.
	Name string `json:"name,required"`
	// The object type, which is `skill`.
	Object constant.Skill `json:"object,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID             respjson.Field
		CreatedAt      respjson.Field
		DefaultVersion respjson.Field
		Description    respjson.Field
		LatestVersion  respjson.Field
		Name           respjson.Field
		Object         respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r Skill) RawJSON() string { return r.JSON.raw }
func (r *Skill) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type SkillList struct {
	// A list of items
	Data []Skill `json:"data,required"`
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
func (r SkillList) RawJSON() string { return r.JSON.raw }
func (r *SkillList) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type SkillNewParams struct {
	// Skill files to upload (directory upload) or a single zip file.
	Files SkillNewParamsFilesUnion `json:"files,omitzero" format:"binary"`
	paramObj
}

func (r SkillNewParams) MarshalMultipart() (data []byte, contentType string, err error) {
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
type SkillNewParamsFilesUnion struct {
	OfFileArray []io.Reader `json:",omitzero,inline"`
	OfFile      io.Reader   `json:",omitzero,inline"`
	paramUnion
}

func (u SkillNewParamsFilesUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfFileArray, u.OfFile)
}
func (u *SkillNewParamsFilesUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *SkillNewParamsFilesUnion) asAny() any {
	if !param.IsOmitted(u.OfFileArray) {
		return &u.OfFileArray
	} else if !param.IsOmitted(u.OfFile) {
		return &u.OfFile
	}
	return nil
}

type SkillUpdateParams struct {
	// The skill version number to set as default.
	DefaultVersion string `json:"default_version,required"`
	paramObj
}

func (r SkillUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow SkillUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *SkillUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type SkillListParams struct {
	// Identifier for the last item from the previous pagination request
	After param.Opt[string] `query:"after,omitzero" json:"-"`
	// Number of items to retrieve
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Sort order of results by timestamp. Use `asc` for ascending order or `desc` for
	// descending order.
	//
	// Any of "asc", "desc".
	Order SkillListParamsOrder `query:"order,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [SkillListParams]'s query parameters as `url.Values`.
func (r SkillListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Sort order of results by timestamp. Use `asc` for ascending order or `desc` for
// descending order.
type SkillListParamsOrder string

const (
	SkillListParamsOrderAsc  SkillListParamsOrder = "asc"
	SkillListParamsOrderDesc SkillListParamsOrder = "desc"
)

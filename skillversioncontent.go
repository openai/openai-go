// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
)

// SkillVersionContentService contains methods and other services that help with
// interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewSkillVersionContentService] method instead.
type SkillVersionContentService struct {
	Options []option.RequestOption
}

// NewSkillVersionContentService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewSkillVersionContentService(opts ...option.RequestOption) (r SkillVersionContentService) {
	r = SkillVersionContentService{}
	r.Options = opts
	return
}

// Download a skill version zip bundle.
func (r *SkillVersionContentService) Get(ctx context.Context, skillID string, version string, opts ...option.RequestOption) (res *http.Response, err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "application/binary")}, opts...)
	if skillID == "" {
		err = errors.New("missing required skill_id parameter")
		return
	}
	if version == "" {
		err = errors.New("missing required version parameter")
		return
	}
	path := fmt.Sprintf("skills/%s/versions/%s/content", skillID, version)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

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

// SkillContentService contains methods and other services that help with
// interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewSkillContentService] method instead.
type SkillContentService struct {
	Options []option.RequestOption
}

// NewSkillContentService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewSkillContentService(opts ...option.RequestOption) (r SkillContentService) {
	r = SkillContentService{}
	r.Options = opts
	return
}

// Download a skill zip bundle by its ID.
func (r *SkillContentService) Get(ctx context.Context, skillID string, opts ...option.RequestOption) (res *http.Response, err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "application/binary")}, opts...)
	if skillID == "" {
		err = errors.New("missing required skill_id parameter")
		return
	}
	path := fmt.Sprintf("skills/%s/content", skillID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

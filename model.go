// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/openai/openai-go/internal/apijson"
	"github.com/openai/openai-go/internal/requestconfig"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/pagination"
)

// ModelService contains methods and other services that help with interacting with
// the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewModelService] method instead.
type ModelService struct {
	Options []option.RequestOption
}

// NewModelService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewModelService(opts ...option.RequestOption) (r *ModelService) {
	r = &ModelService{}
	r.Options = opts
	return
}

// Retrieves a model instance, providing basic information about the model such as
// the owner and permissioning.
func (r *ModelService) Get(ctx context.Context, model string, opts ...option.RequestOption) (res *Model, err error) {
	opts = append(r.Options[:], opts...)
	if model == "" {
		err = errors.New("missing required model parameter")
		return
	}
	path := fmt.Sprintf("models/%s", model)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// Lists the currently available models, and provides basic information about each
// one such as the owner and availability.
func (r *ModelService) List(ctx context.Context, opts ...option.RequestOption) (res *pagination.Page[Model], err error) {
	var raw *http.Response
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "models"
	cfg, err := requestconfig.NewRequestConfig(ctx, http.MethodGet, path, nil, &res, opts...)
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

// Lists the currently available models, and provides basic information about each
// one such as the owner and availability.
func (r *ModelService) ListAutoPaging(ctx context.Context, opts ...option.RequestOption) *pagination.PageAutoPager[Model] {
	return pagination.NewPageAutoPager(r.List(ctx, opts...))
}

// Delete a fine-tuned model. You must have the Owner role in your organization to
// delete a model.
func (r *ModelService) Delete(ctx context.Context, model string, opts ...option.RequestOption) (res *ModelDeleted, err error) {
	opts = append(r.Options[:], opts...)
	if model == "" {
		err = errors.New("missing required model parameter")
		return
	}
	path := fmt.Sprintf("models/%s", model)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return
}

// Describes an OpenAI model offering that can be used with the API.
type Model struct {
	// The model identifier, which can be referenced in the API endpoints.
	ID string `json:"id,required"`
	// The Unix timestamp (in seconds) when the model was created.
	Created int64 `json:"created,required"`
	// The object type, which is always "model".
	Object ModelObject `json:"object,required"`
	// The organization that owns the model.
	OwnedBy string    `json:"owned_by,required"`
	JSON    modelJSON `json:"-"`
}

// modelJSON contains the JSON metadata for the struct [Model]
type modelJSON struct {
	ID          apijson.Field
	Created     apijson.Field
	Object      apijson.Field
	OwnedBy     apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *Model) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r modelJSON) RawJSON() string {
	return r.raw
}

// The object type, which is always "model".
type ModelObject string

const (
	ModelObjectModel ModelObject = "model"
)

func (r ModelObject) IsKnown() bool {
	switch r {
	case ModelObjectModel:
		return true
	}
	return false
}

type ModelDeleted struct {
	ID      string           `json:"id,required"`
	Deleted bool             `json:"deleted,required"`
	Object  string           `json:"object,required"`
	JSON    modelDeletedJSON `json:"-"`
}

// modelDeletedJSON contains the JSON metadata for the struct [ModelDeleted]
type modelDeletedJSON struct {
	ID          apijson.Field
	Deleted     apijson.Field
	Object      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ModelDeleted) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r modelDeletedJSON) RawJSON() string {
	return r.raw
}

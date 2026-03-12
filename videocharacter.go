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
	"slices"

	"github.com/openai/openai-go/v3/internal/apiform"
	"github.com/openai/openai-go/v3/internal/apijson"
	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/respjson"
)

// VideoCharacterService contains methods and other services that help with
// interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewVideoCharacterService] method instead.
type VideoCharacterService struct {
	Options []option.RequestOption
}

// NewVideoCharacterService generates a new service that applies the given options
// to each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewVideoCharacterService(opts ...option.RequestOption) (r VideoCharacterService) {
	r = VideoCharacterService{}
	r.Options = opts
	return
}

// Create a character from an uploaded video.
func (r *VideoCharacterService) New(ctx context.Context, body VideoCharacterNewParams, opts ...option.RequestOption) (res *VideoCharacterNewResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "videos/characters"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Fetch a character.
func (r *VideoCharacterService) Get(ctx context.Context, characterID string, opts ...option.RequestOption) (res *VideoCharacterGetResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if characterID == "" {
		err = errors.New("missing required character_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("videos/characters/%s", characterID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

type VideoCharacterNewResponse struct {
	// Identifier for the character creation cameo.
	ID string `json:"id" api:"required"`
	// Unix timestamp (in seconds) when the character was created.
	CreatedAt int64 `json:"created_at" api:"required"`
	// Display name for the character.
	Name string `json:"name" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		Name        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r VideoCharacterNewResponse) RawJSON() string { return r.JSON.raw }
func (r *VideoCharacterNewResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type VideoCharacterGetResponse struct {
	// Identifier for the character creation cameo.
	ID string `json:"id" api:"required"`
	// Unix timestamp (in seconds) when the character was created.
	CreatedAt int64 `json:"created_at" api:"required"`
	// Display name for the character.
	Name string `json:"name" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		Name        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r VideoCharacterGetResponse) RawJSON() string { return r.JSON.raw }
func (r *VideoCharacterGetResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type VideoCharacterNewParams struct {
	// Display name for this API character.
	Name string `json:"name" api:"required"`
	// Video file used to create a character.
	Video io.Reader `json:"video,omitzero" api:"required" format:"binary"`
	paramObj
}

func (r VideoCharacterNewParams) MarshalMultipart() (data []byte, contentType string, err error) {
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

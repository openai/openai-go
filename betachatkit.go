// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"slices"

	"github.com/openai/openai-go/v3/internal/apiform"
	"github.com/openai/openai-go/v3/internal/apijson"
	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/respjson"
	"github.com/openai/openai-go/v3/shared/constant"
)

// BetaChatKitService contains methods and other services that help with
// interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaChatKitService] method instead.
type BetaChatKitService struct {
	Options  []option.RequestOption
	Sessions BetaChatKitSessionService
	Threads  BetaChatKitThreadService
}

// NewBetaChatKitService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewBetaChatKitService(opts ...option.RequestOption) (r BetaChatKitService) {
	r = BetaChatKitService{}
	r.Options = opts
	r.Sessions = NewBetaChatKitSessionService(opts...)
	r.Threads = NewBetaChatKitThreadService(opts...)
	return
}

// Upload a ChatKit file
func (r *BetaChatKitService) UploadFile(ctx context.Context, body BetaChatKitUploadFileParams, opts ...option.RequestOption) (res *BetaChatKitUploadFileResponseUnion, err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "chatkit_beta=v1")}, opts...)
	path := "chatkit/files"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Workflow metadata and state returned for the session.
type ChatKitWorkflow struct {
	// Identifier of the workflow backing the session.
	ID string `json:"id,required"`
	// State variable key-value pairs applied when invoking the workflow. Defaults to
	// null when no overrides were provided.
	StateVariables map[string]ChatKitWorkflowStateVariableUnion `json:"state_variables,required"`
	// Tracing settings applied to the workflow.
	Tracing ChatKitWorkflowTracing `json:"tracing,required"`
	// Specific workflow version used for the session. Defaults to null when using the
	// latest deployment.
	Version string `json:"version,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID             respjson.Field
		StateVariables respjson.Field
		Tracing        respjson.Field
		Version        respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ChatKitWorkflow) RawJSON() string { return r.JSON.raw }
func (r *ChatKitWorkflow) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ChatKitWorkflowStateVariableUnion contains all possible properties and values
// from [string], [bool], [float64].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfString OfBool OfFloat]
type ChatKitWorkflowStateVariableUnion struct {
	// This field will be present if the value is a [string] instead of an object.
	OfString string `json:",inline"`
	// This field will be present if the value is a [bool] instead of an object.
	OfBool bool `json:",inline"`
	// This field will be present if the value is a [float64] instead of an object.
	OfFloat float64 `json:",inline"`
	JSON    struct {
		OfString respjson.Field
		OfBool   respjson.Field
		OfFloat  respjson.Field
		raw      string
	} `json:"-"`
}

func (u ChatKitWorkflowStateVariableUnion) AsString() (v string) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ChatKitWorkflowStateVariableUnion) AsBool() (v bool) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ChatKitWorkflowStateVariableUnion) AsFloat() (v float64) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ChatKitWorkflowStateVariableUnion) RawJSON() string { return u.JSON.raw }

func (r *ChatKitWorkflowStateVariableUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Tracing settings applied to the workflow.
type ChatKitWorkflowTracing struct {
	// Indicates whether tracing is enabled.
	Enabled bool `json:"enabled,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Enabled     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ChatKitWorkflowTracing) RawJSON() string { return r.JSON.raw }
func (r *ChatKitWorkflowTracing) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Metadata for a non-image file uploaded through ChatKit.
type FilePart struct {
	// Unique identifier for the uploaded file.
	ID string `json:"id,required"`
	// MIME type reported for the uploaded file. Defaults to null when unknown.
	MimeType string `json:"mime_type,required"`
	// Original filename supplied by the uploader. Defaults to null when unnamed.
	Name string `json:"name,required"`
	// Type discriminator that is always `file`.
	Type constant.File `json:"type,required"`
	// Signed URL for downloading the uploaded file. Defaults to null when no download
	// link is available.
	UploadURL string `json:"upload_url,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		MimeType    respjson.Field
		Name        respjson.Field
		Type        respjson.Field
		UploadURL   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r FilePart) RawJSON() string { return r.JSON.raw }
func (r *FilePart) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Metadata for an image uploaded through ChatKit.
type ImagePart struct {
	// Unique identifier for the uploaded image.
	ID string `json:"id,required"`
	// MIME type of the uploaded image.
	MimeType string `json:"mime_type,required"`
	// Original filename for the uploaded image. Defaults to null when unnamed.
	Name string `json:"name,required"`
	// Preview URL that can be rendered inline for the image.
	PreviewURL string `json:"preview_url,required"`
	// Type discriminator that is always `image`.
	Type constant.Image `json:"type,required"`
	// Signed URL for downloading the uploaded image. Defaults to null when no download
	// link is available.
	UploadURL string `json:"upload_url,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		MimeType    respjson.Field
		Name        respjson.Field
		PreviewURL  respjson.Field
		Type        respjson.Field
		UploadURL   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ImagePart) RawJSON() string { return r.JSON.raw }
func (r *ImagePart) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaChatKitUploadFileResponseUnion contains all possible properties and values
// from [FilePart], [ImagePart].
//
// Use the [BetaChatKitUploadFileResponseUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaChatKitUploadFileResponseUnion struct {
	ID       string `json:"id"`
	MimeType string `json:"mime_type"`
	Name     string `json:"name"`
	// Any of "file", "image".
	Type      string `json:"type"`
	UploadURL string `json:"upload_url"`
	// This field is from variant [ImagePart].
	PreviewURL string `json:"preview_url"`
	JSON       struct {
		ID         respjson.Field
		MimeType   respjson.Field
		Name       respjson.Field
		Type       respjson.Field
		UploadURL  respjson.Field
		PreviewURL respjson.Field
		raw        string
	} `json:"-"`
}

// anyBetaChatKitUploadFileResponse is implemented by each variant of
// [BetaChatKitUploadFileResponseUnion] to add type safety for the return type of
// [BetaChatKitUploadFileResponseUnion.AsAny]
type anyBetaChatKitUploadFileResponse interface {
	implBetaChatKitUploadFileResponseUnion()
}

func (FilePart) implBetaChatKitUploadFileResponseUnion()  {}
func (ImagePart) implBetaChatKitUploadFileResponseUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaChatKitUploadFileResponseUnion.AsAny().(type) {
//	case openai.FilePart:
//	case openai.ImagePart:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaChatKitUploadFileResponseUnion) AsAny() anyBetaChatKitUploadFileResponse {
	switch u.Type {
	case "file":
		return u.AsFile()
	case "image":
		return u.AsImage()
	}
	return nil
}

func (u BetaChatKitUploadFileResponseUnion) AsFile() (v FilePart) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaChatKitUploadFileResponseUnion) AsImage() (v ImagePart) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaChatKitUploadFileResponseUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaChatKitUploadFileResponseUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaChatKitUploadFileParams struct {
	// Binary file contents to store with the ChatKit session. Supports PDFs and PNG,
	// JPG, JPEG, GIF, or WEBP images.
	File io.Reader `json:"file,omitzero,required" format:"binary"`
	paramObj
}

func (r BetaChatKitUploadFileParams) MarshalMultipart() (data []byte, contentType string, err error) {
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

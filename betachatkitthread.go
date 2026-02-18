// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai

import (
	"context"
	"encoding/json"
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
	"github.com/openai/openai-go/v3/shared/constant"
)

// BetaChatKitThreadService contains methods and other services that help with
// interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaChatKitThreadService] method instead.
type BetaChatKitThreadService struct {
	Options []option.RequestOption
}

// NewBetaChatKitThreadService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewBetaChatKitThreadService(opts ...option.RequestOption) (r BetaChatKitThreadService) {
	r = BetaChatKitThreadService{}
	r.Options = opts
	return
}

// Retrieve a ChatKit thread by its identifier.
func (r *BetaChatKitThreadService) Get(ctx context.Context, threadID string, opts ...option.RequestOption) (res *ChatKitThread, err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "chatkit_beta=v1")}, opts...)
	if threadID == "" {
		err = errors.New("missing required thread_id parameter")
		return
	}
	path := fmt.Sprintf("chatkit/threads/%s", threadID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// List ChatKit threads with optional pagination and user filters.
func (r *BetaChatKitThreadService) List(ctx context.Context, query BetaChatKitThreadListParams, opts ...option.RequestOption) (res *pagination.ConversationCursorPage[ChatKitThread], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "chatkit_beta=v1"), option.WithResponseInto(&raw)}, opts...)
	path := "chatkit/threads"
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

// List ChatKit threads with optional pagination and user filters.
func (r *BetaChatKitThreadService) ListAutoPaging(ctx context.Context, query BetaChatKitThreadListParams, opts ...option.RequestOption) *pagination.ConversationCursorPageAutoPager[ChatKitThread] {
	return pagination.NewConversationCursorPageAutoPager(r.List(ctx, query, opts...))
}

// Delete a ChatKit thread along with its items and stored attachments.
func (r *BetaChatKitThreadService) Delete(ctx context.Context, threadID string, opts ...option.RequestOption) (res *BetaChatKitThreadDeleteResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "chatkit_beta=v1")}, opts...)
	if threadID == "" {
		err = errors.New("missing required thread_id parameter")
		return
	}
	path := fmt.Sprintf("chatkit/threads/%s", threadID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return
}

// List items that belong to a ChatKit thread.
func (r *BetaChatKitThreadService) ListItems(ctx context.Context, threadID string, query BetaChatKitThreadListItemsParams, opts ...option.RequestOption) (res *pagination.ConversationCursorPage[ChatKitThreadItemListDataUnion], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("OpenAI-Beta", "chatkit_beta=v1"), option.WithResponseInto(&raw)}, opts...)
	if threadID == "" {
		err = errors.New("missing required thread_id parameter")
		return
	}
	path := fmt.Sprintf("chatkit/threads/%s/items", threadID)
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

// List items that belong to a ChatKit thread.
func (r *BetaChatKitThreadService) ListItemsAutoPaging(ctx context.Context, threadID string, query BetaChatKitThreadListItemsParams, opts ...option.RequestOption) *pagination.ConversationCursorPageAutoPager[ChatKitThreadItemListDataUnion] {
	return pagination.NewConversationCursorPageAutoPager(r.ListItems(ctx, threadID, query, opts...))
}

// Represents a ChatKit session and its resolved configuration.
type ChatSession struct {
	// Identifier for the ChatKit session.
	ID string `json:"id,required"`
	// Resolved ChatKit feature configuration for the session.
	ChatKitConfiguration ChatSessionChatKitConfiguration `json:"chatkit_configuration,required"`
	// Ephemeral client secret that authenticates session requests.
	ClientSecret string `json:"client_secret,required"`
	// Unix timestamp (in seconds) for when the session expires.
	ExpiresAt int64 `json:"expires_at,required"`
	// Convenience copy of the per-minute request limit.
	MaxRequestsPer1Minute int64 `json:"max_requests_per_1_minute,required"`
	// Type discriminator that is always `chatkit.session`.
	Object constant.ChatKitSession `json:"object,required"`
	// Resolved rate limit values.
	RateLimits ChatSessionRateLimits `json:"rate_limits,required"`
	// Current lifecycle state of the session.
	//
	// Any of "active", "expired", "cancelled".
	Status ChatSessionStatus `json:"status,required"`
	// User identifier associated with the session.
	User string `json:"user,required"`
	// Workflow metadata for the session.
	Workflow ChatKitWorkflow `json:"workflow,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                    respjson.Field
		ChatKitConfiguration  respjson.Field
		ClientSecret          respjson.Field
		ExpiresAt             respjson.Field
		MaxRequestsPer1Minute respjson.Field
		Object                respjson.Field
		RateLimits            respjson.Field
		Status                respjson.Field
		User                  respjson.Field
		Workflow              respjson.Field
		ExtraFields           map[string]respjson.Field
		raw                   string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ChatSession) RawJSON() string { return r.JSON.raw }
func (r *ChatSession) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Automatic thread title preferences for the session.
type ChatSessionAutomaticThreadTitling struct {
	// Whether automatic thread titling is enabled.
	Enabled bool `json:"enabled,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Enabled     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ChatSessionAutomaticThreadTitling) RawJSON() string { return r.JSON.raw }
func (r *ChatSessionAutomaticThreadTitling) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ChatKit configuration for the session.
type ChatSessionChatKitConfiguration struct {
	// Automatic thread titling preferences.
	AutomaticThreadTitling ChatSessionAutomaticThreadTitling `json:"automatic_thread_titling,required"`
	// Upload settings for the session.
	FileUpload ChatSessionFileUpload `json:"file_upload,required"`
	// History retention configuration.
	History ChatSessionHistory `json:"history,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		AutomaticThreadTitling respjson.Field
		FileUpload             respjson.Field
		History                respjson.Field
		ExtraFields            map[string]respjson.Field
		raw                    string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ChatSessionChatKitConfiguration) RawJSON() string { return r.JSON.raw }
func (r *ChatSessionChatKitConfiguration) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Optional per-session configuration settings for ChatKit behavior.
type ChatSessionChatKitConfigurationParam struct {
	// Configuration for automatic thread titling. When omitted, automatic thread
	// titling is enabled by default.
	AutomaticThreadTitling ChatSessionChatKitConfigurationParamAutomaticThreadTitling `json:"automatic_thread_titling,omitzero"`
	// Configuration for upload enablement and limits. When omitted, uploads are
	// disabled by default (max_files 10, max_file_size 512 MB).
	FileUpload ChatSessionChatKitConfigurationParamFileUpload `json:"file_upload,omitzero"`
	// Configuration for chat history retention. When omitted, history is enabled by
	// default with no limit on recent_threads (null).
	History ChatSessionChatKitConfigurationParamHistory `json:"history,omitzero"`
	paramObj
}

func (r ChatSessionChatKitConfigurationParam) MarshalJSON() (data []byte, err error) {
	type shadow ChatSessionChatKitConfigurationParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ChatSessionChatKitConfigurationParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for automatic thread titling. When omitted, automatic thread
// titling is enabled by default.
type ChatSessionChatKitConfigurationParamAutomaticThreadTitling struct {
	// Enable automatic thread title generation. Defaults to true.
	Enabled param.Opt[bool] `json:"enabled,omitzero"`
	paramObj
}

func (r ChatSessionChatKitConfigurationParamAutomaticThreadTitling) MarshalJSON() (data []byte, err error) {
	type shadow ChatSessionChatKitConfigurationParamAutomaticThreadTitling
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ChatSessionChatKitConfigurationParamAutomaticThreadTitling) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for upload enablement and limits. When omitted, uploads are
// disabled by default (max_files 10, max_file_size 512 MB).
type ChatSessionChatKitConfigurationParamFileUpload struct {
	// Enable uploads for this session. Defaults to false.
	Enabled param.Opt[bool] `json:"enabled,omitzero"`
	// Maximum size in megabytes for each uploaded file. Defaults to 512 MB, which is
	// the maximum allowable size.
	MaxFileSize param.Opt[int64] `json:"max_file_size,omitzero"`
	// Maximum number of files that can be uploaded to the session. Defaults to 10.
	MaxFiles param.Opt[int64] `json:"max_files,omitzero"`
	paramObj
}

func (r ChatSessionChatKitConfigurationParamFileUpload) MarshalJSON() (data []byte, err error) {
	type shadow ChatSessionChatKitConfigurationParamFileUpload
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ChatSessionChatKitConfigurationParamFileUpload) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for chat history retention. When omitted, history is enabled by
// default with no limit on recent_threads (null).
type ChatSessionChatKitConfigurationParamHistory struct {
	// Enables chat users to access previous ChatKit threads. Defaults to true.
	Enabled param.Opt[bool] `json:"enabled,omitzero"`
	// Number of recent ChatKit threads users have access to. Defaults to unlimited
	// when unset.
	RecentThreads param.Opt[int64] `json:"recent_threads,omitzero"`
	paramObj
}

func (r ChatSessionChatKitConfigurationParamHistory) MarshalJSON() (data []byte, err error) {
	type shadow ChatSessionChatKitConfigurationParamHistory
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ChatSessionChatKitConfigurationParamHistory) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Controls when the session expires relative to an anchor timestamp.
//
// The properties Anchor, Seconds are required.
type ChatSessionExpiresAfterParam struct {
	// Number of seconds after the anchor when the session expires.
	Seconds int64 `json:"seconds,required"`
	// Base timestamp used to calculate expiration. Currently fixed to `created_at`.
	//
	// This field can be elided, and will marshal its zero value as "created_at".
	Anchor constant.CreatedAt `json:"anchor,required"`
	paramObj
}

func (r ChatSessionExpiresAfterParam) MarshalJSON() (data []byte, err error) {
	type shadow ChatSessionExpiresAfterParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ChatSessionExpiresAfterParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Upload permissions and limits applied to the session.
type ChatSessionFileUpload struct {
	// Indicates if uploads are enabled for the session.
	Enabled bool `json:"enabled,required"`
	// Maximum upload size in megabytes.
	MaxFileSize int64 `json:"max_file_size,required"`
	// Maximum number of uploads allowed during the session.
	MaxFiles int64 `json:"max_files,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Enabled     respjson.Field
		MaxFileSize respjson.Field
		MaxFiles    respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ChatSessionFileUpload) RawJSON() string { return r.JSON.raw }
func (r *ChatSessionFileUpload) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// History retention preferences returned for the session.
type ChatSessionHistory struct {
	// Indicates if chat history is persisted for the session.
	Enabled bool `json:"enabled,required"`
	// Number of prior threads surfaced in history views. Defaults to null when all
	// history is retained.
	RecentThreads int64 `json:"recent_threads,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Enabled       respjson.Field
		RecentThreads respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ChatSessionHistory) RawJSON() string { return r.JSON.raw }
func (r *ChatSessionHistory) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Active per-minute request limit for the session.
type ChatSessionRateLimits struct {
	// Maximum allowed requests per one-minute window.
	MaxRequestsPer1Minute int64 `json:"max_requests_per_1_minute,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		MaxRequestsPer1Minute respjson.Field
		ExtraFields           map[string]respjson.Field
		raw                   string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ChatSessionRateLimits) RawJSON() string { return r.JSON.raw }
func (r *ChatSessionRateLimits) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Controls request rate limits for the session.
type ChatSessionRateLimitsParam struct {
	// Maximum number of requests allowed per minute for the session. Defaults to 10.
	MaxRequestsPer1Minute param.Opt[int64] `json:"max_requests_per_1_minute,omitzero"`
	paramObj
}

func (r ChatSessionRateLimitsParam) MarshalJSON() (data []byte, err error) {
	type shadow ChatSessionRateLimitsParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ChatSessionRateLimitsParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ChatSessionStatus string

const (
	ChatSessionStatusActive    ChatSessionStatus = "active"
	ChatSessionStatusExpired   ChatSessionStatus = "expired"
	ChatSessionStatusCancelled ChatSessionStatus = "cancelled"
)

// Workflow reference and overrides applied to the chat session.
//
// The property ID is required.
type ChatSessionWorkflowParam struct {
	// Identifier for the workflow invoked by the session.
	ID string `json:"id,required"`
	// Specific workflow version to run. Defaults to the latest deployed version.
	Version param.Opt[string] `json:"version,omitzero"`
	// State variables forwarded to the workflow. Keys may be up to 64 characters,
	// values must be primitive types, and the map defaults to an empty object.
	StateVariables map[string]ChatSessionWorkflowParamStateVariableUnion `json:"state_variables,omitzero"`
	// Optional tracing overrides for the workflow invocation. When omitted, tracing is
	// enabled by default.
	Tracing ChatSessionWorkflowParamTracing `json:"tracing,omitzero"`
	paramObj
}

func (r ChatSessionWorkflowParam) MarshalJSON() (data []byte, err error) {
	type shadow ChatSessionWorkflowParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ChatSessionWorkflowParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ChatSessionWorkflowParamStateVariableUnion struct {
	OfString param.Opt[string]  `json:",omitzero,inline"`
	OfBool   param.Opt[bool]    `json:",omitzero,inline"`
	OfFloat  param.Opt[float64] `json:",omitzero,inline"`
	paramUnion
}

func (u ChatSessionWorkflowParamStateVariableUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfString, u.OfBool, u.OfFloat)
}
func (u *ChatSessionWorkflowParamStateVariableUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *ChatSessionWorkflowParamStateVariableUnion) asAny() any {
	if !param.IsOmitted(u.OfString) {
		return &u.OfString.Value
	} else if !param.IsOmitted(u.OfBool) {
		return &u.OfBool.Value
	} else if !param.IsOmitted(u.OfFloat) {
		return &u.OfFloat.Value
	}
	return nil
}

// Optional tracing overrides for the workflow invocation. When omitted, tracing is
// enabled by default.
type ChatSessionWorkflowParamTracing struct {
	// Whether tracing is enabled during the session. Defaults to true.
	Enabled param.Opt[bool] `json:"enabled,omitzero"`
	paramObj
}

func (r ChatSessionWorkflowParamTracing) MarshalJSON() (data []byte, err error) {
	type shadow ChatSessionWorkflowParamTracing
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ChatSessionWorkflowParamTracing) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Attachment metadata included on thread items.
type ChatKitAttachment struct {
	// Identifier for the attachment.
	ID string `json:"id,required"`
	// MIME type of the attachment.
	MimeType string `json:"mime_type,required"`
	// Original display name for the attachment.
	Name string `json:"name,required"`
	// Preview URL for rendering the attachment inline.
	PreviewURL string `json:"preview_url,required"`
	// Attachment discriminator.
	//
	// Any of "image", "file".
	Type ChatKitAttachmentType `json:"type,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		MimeType    respjson.Field
		Name        respjson.Field
		PreviewURL  respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ChatKitAttachment) RawJSON() string { return r.JSON.raw }
func (r *ChatKitAttachment) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Attachment discriminator.
type ChatKitAttachmentType string

const (
	ChatKitAttachmentTypeImage ChatKitAttachmentType = "image"
	ChatKitAttachmentTypeFile  ChatKitAttachmentType = "file"
)

// Assistant response text accompanied by optional annotations.
type ChatKitResponseOutputText struct {
	// Ordered list of annotations attached to the response text.
	Annotations []ChatKitResponseOutputTextAnnotationUnion `json:"annotations,required"`
	// Assistant generated text.
	Text string `json:"text,required"`
	// Type discriminator that is always `output_text`.
	Type constant.OutputText `json:"type,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Annotations respjson.Field
		Text        respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ChatKitResponseOutputText) RawJSON() string { return r.JSON.raw }
func (r *ChatKitResponseOutputText) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ChatKitResponseOutputTextAnnotationUnion contains all possible properties and
// values from [ChatKitResponseOutputTextAnnotationFile],
// [ChatKitResponseOutputTextAnnotationURL].
//
// Use the [ChatKitResponseOutputTextAnnotationUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ChatKitResponseOutputTextAnnotationUnion struct {
	// This field is a union of [ChatKitResponseOutputTextAnnotationFileSource],
	// [ChatKitResponseOutputTextAnnotationURLSource]
	Source ChatKitResponseOutputTextAnnotationUnionSource `json:"source"`
	// Any of "file", "url".
	Type string `json:"type"`
	JSON struct {
		Source respjson.Field
		Type   respjson.Field
		raw    string
	} `json:"-"`
}

// anyChatKitResponseOutputTextAnnotation is implemented by each variant of
// [ChatKitResponseOutputTextAnnotationUnion] to add type safety for the return
// type of [ChatKitResponseOutputTextAnnotationUnion.AsAny]
type anyChatKitResponseOutputTextAnnotation interface {
	implChatKitResponseOutputTextAnnotationUnion()
}

func (ChatKitResponseOutputTextAnnotationFile) implChatKitResponseOutputTextAnnotationUnion() {}
func (ChatKitResponseOutputTextAnnotationURL) implChatKitResponseOutputTextAnnotationUnion()  {}

// Use the following switch statement to find the correct variant
//
//	switch variant := ChatKitResponseOutputTextAnnotationUnion.AsAny().(type) {
//	case openai.ChatKitResponseOutputTextAnnotationFile:
//	case openai.ChatKitResponseOutputTextAnnotationURL:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ChatKitResponseOutputTextAnnotationUnion) AsAny() anyChatKitResponseOutputTextAnnotation {
	switch u.Type {
	case "file":
		return u.AsFile()
	case "url":
		return u.AsURL()
	}
	return nil
}

func (u ChatKitResponseOutputTextAnnotationUnion) AsFile() (v ChatKitResponseOutputTextAnnotationFile) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ChatKitResponseOutputTextAnnotationUnion) AsURL() (v ChatKitResponseOutputTextAnnotationURL) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ChatKitResponseOutputTextAnnotationUnion) RawJSON() string { return u.JSON.raw }

func (r *ChatKitResponseOutputTextAnnotationUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ChatKitResponseOutputTextAnnotationUnionSource is an implicit subunion of
// [ChatKitResponseOutputTextAnnotationUnion].
// ChatKitResponseOutputTextAnnotationUnionSource provides convenient access to the
// sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [ChatKitResponseOutputTextAnnotationUnion].
type ChatKitResponseOutputTextAnnotationUnionSource struct {
	// This field is from variant [ChatKitResponseOutputTextAnnotationFileSource].
	Filename string `json:"filename"`
	Type     string `json:"type"`
	// This field is from variant [ChatKitResponseOutputTextAnnotationURLSource].
	URL  string `json:"url"`
	JSON struct {
		Filename respjson.Field
		Type     respjson.Field
		URL      respjson.Field
		raw      string
	} `json:"-"`
}

func (r *ChatKitResponseOutputTextAnnotationUnionSource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Annotation that references an uploaded file.
type ChatKitResponseOutputTextAnnotationFile struct {
	// File attachment referenced by the annotation.
	Source ChatKitResponseOutputTextAnnotationFileSource `json:"source,required"`
	// Type discriminator that is always `file` for this annotation.
	Type constant.File `json:"type,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Source      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ChatKitResponseOutputTextAnnotationFile) RawJSON() string { return r.JSON.raw }
func (r *ChatKitResponseOutputTextAnnotationFile) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// File attachment referenced by the annotation.
type ChatKitResponseOutputTextAnnotationFileSource struct {
	// Filename referenced by the annotation.
	Filename string `json:"filename,required"`
	// Type discriminator that is always `file`.
	Type constant.File `json:"type,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Filename    respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ChatKitResponseOutputTextAnnotationFileSource) RawJSON() string { return r.JSON.raw }
func (r *ChatKitResponseOutputTextAnnotationFileSource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Annotation that references a URL.
type ChatKitResponseOutputTextAnnotationURL struct {
	// URL referenced by the annotation.
	Source ChatKitResponseOutputTextAnnotationURLSource `json:"source,required"`
	// Type discriminator that is always `url` for this annotation.
	Type constant.URL `json:"type,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Source      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ChatKitResponseOutputTextAnnotationURL) RawJSON() string { return r.JSON.raw }
func (r *ChatKitResponseOutputTextAnnotationURL) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// URL referenced by the annotation.
type ChatKitResponseOutputTextAnnotationURLSource struct {
	// Type discriminator that is always `url`.
	Type constant.URL `json:"type,required"`
	// URL referenced by the annotation.
	URL string `json:"url,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		URL         respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ChatKitResponseOutputTextAnnotationURLSource) RawJSON() string { return r.JSON.raw }
func (r *ChatKitResponseOutputTextAnnotationURLSource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Represents a ChatKit thread and its current status.
type ChatKitThread struct {
	// Identifier of the thread.
	ID string `json:"id,required"`
	// Unix timestamp (in seconds) for when the thread was created.
	CreatedAt int64 `json:"created_at,required"`
	// Type discriminator that is always `chatkit.thread`.
	Object constant.ChatKitThread `json:"object,required"`
	// Current status for the thread. Defaults to `active` for newly created threads.
	Status ChatKitThreadStatusUnion `json:"status,required"`
	// Optional human-readable title for the thread. Defaults to null when no title has
	// been generated.
	Title string `json:"title,required"`
	// Free-form string that identifies your end user who owns the thread.
	User string `json:"user,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		Object      respjson.Field
		Status      respjson.Field
		Title       respjson.Field
		User        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ChatKitThread) RawJSON() string { return r.JSON.raw }
func (r *ChatKitThread) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ChatKitThreadStatusUnion contains all possible properties and values from
// [ChatKitThreadStatusActive], [ChatKitThreadStatusLocked],
// [ChatKitThreadStatusClosed].
//
// Use the [ChatKitThreadStatusUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ChatKitThreadStatusUnion struct {
	// Any of "active", "locked", "closed".
	Type   string `json:"type"`
	Reason string `json:"reason"`
	JSON   struct {
		Type   respjson.Field
		Reason respjson.Field
		raw    string
	} `json:"-"`
}

// anyChatKitThreadStatus is implemented by each variant of
// [ChatKitThreadStatusUnion] to add type safety for the return type of
// [ChatKitThreadStatusUnion.AsAny]
type anyChatKitThreadStatus interface {
	implChatKitThreadStatusUnion()
}

func (ChatKitThreadStatusActive) implChatKitThreadStatusUnion() {}
func (ChatKitThreadStatusLocked) implChatKitThreadStatusUnion() {}
func (ChatKitThreadStatusClosed) implChatKitThreadStatusUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := ChatKitThreadStatusUnion.AsAny().(type) {
//	case openai.ChatKitThreadStatusActive:
//	case openai.ChatKitThreadStatusLocked:
//	case openai.ChatKitThreadStatusClosed:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ChatKitThreadStatusUnion) AsAny() anyChatKitThreadStatus {
	switch u.Type {
	case "active":
		return u.AsActive()
	case "locked":
		return u.AsLocked()
	case "closed":
		return u.AsClosed()
	}
	return nil
}

func (u ChatKitThreadStatusUnion) AsActive() (v ChatKitThreadStatusActive) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ChatKitThreadStatusUnion) AsLocked() (v ChatKitThreadStatusLocked) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ChatKitThreadStatusUnion) AsClosed() (v ChatKitThreadStatusClosed) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ChatKitThreadStatusUnion) RawJSON() string { return u.JSON.raw }

func (r *ChatKitThreadStatusUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Indicates that a thread is active.
type ChatKitThreadStatusActive struct {
	// Status discriminator that is always `active`.
	Type constant.Active `json:"type,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ChatKitThreadStatusActive) RawJSON() string { return r.JSON.raw }
func (r *ChatKitThreadStatusActive) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Indicates that a thread is locked and cannot accept new input.
type ChatKitThreadStatusLocked struct {
	// Reason that the thread was locked. Defaults to null when no reason is recorded.
	Reason string `json:"reason,required"`
	// Status discriminator that is always `locked`.
	Type constant.Locked `json:"type,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Reason      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ChatKitThreadStatusLocked) RawJSON() string { return r.JSON.raw }
func (r *ChatKitThreadStatusLocked) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Indicates that a thread has been closed.
type ChatKitThreadStatusClosed struct {
	// Reason that the thread was closed. Defaults to null when no reason is recorded.
	Reason string `json:"reason,required"`
	// Status discriminator that is always `closed`.
	Type constant.Closed `json:"type,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Reason      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ChatKitThreadStatusClosed) RawJSON() string { return r.JSON.raw }
func (r *ChatKitThreadStatusClosed) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Assistant-authored message within a thread.
type ChatKitThreadAssistantMessageItem struct {
	// Identifier of the thread item.
	ID string `json:"id,required"`
	// Ordered assistant response segments.
	Content []ChatKitResponseOutputText `json:"content,required"`
	// Unix timestamp (in seconds) for when the item was created.
	CreatedAt int64 `json:"created_at,required"`
	// Type discriminator that is always `chatkit.thread_item`.
	Object constant.ChatKitThreadItem `json:"object,required"`
	// Identifier of the parent thread.
	ThreadID string `json:"thread_id,required"`
	// Type discriminator that is always `chatkit.assistant_message`.
	Type constant.ChatKitAssistantMessage `json:"type,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Content     respjson.Field
		CreatedAt   respjson.Field
		Object      respjson.Field
		ThreadID    respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ChatKitThreadAssistantMessageItem) RawJSON() string { return r.JSON.raw }
func (r *ChatKitThreadAssistantMessageItem) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A paginated list of thread items rendered for the ChatKit API.
type ChatKitThreadItemList struct {
	// A list of items
	Data []ChatKitThreadItemListDataUnion `json:"data,required"`
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
func (r ChatKitThreadItemList) RawJSON() string { return r.JSON.raw }
func (r *ChatKitThreadItemList) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ChatKitThreadItemListDataUnion contains all possible properties and values from
// [ChatKitThreadUserMessageItem], [ChatKitThreadAssistantMessageItem],
// [ChatKitWidgetItem], [ChatKitThreadItemListDataChatKitClientToolCall],
// [ChatKitThreadItemListDataChatKitTask],
// [ChatKitThreadItemListDataChatKitTaskGroup].
//
// Use the [ChatKitThreadItemListDataUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ChatKitThreadItemListDataUnion struct {
	ID string `json:"id"`
	// This field is from variant [ChatKitThreadUserMessageItem].
	Attachments []ChatKitAttachment `json:"attachments"`
	// This field is a union of [[]ChatKitThreadUserMessageItemContentUnion],
	// [[]ChatKitResponseOutputText]
	Content   ChatKitThreadItemListDataUnionContent `json:"content"`
	CreatedAt int64                                 `json:"created_at"`
	// This field is from variant [ChatKitThreadUserMessageItem].
	InferenceOptions ChatKitThreadUserMessageItemInferenceOptions `json:"inference_options"`
	// This field is from variant [ChatKitThreadUserMessageItem].
	Object   constant.ChatKitThreadItem `json:"object"`
	ThreadID string                     `json:"thread_id"`
	// Any of "chatkit.user_message", "chatkit.assistant_message", "chatkit.widget",
	// "chatkit.client_tool_call", "chatkit.task", "chatkit.task_group".
	Type string `json:"type"`
	// This field is from variant [ChatKitWidgetItem].
	Widget string `json:"widget"`
	// This field is from variant [ChatKitThreadItemListDataChatKitClientToolCall].
	Arguments string `json:"arguments"`
	// This field is from variant [ChatKitThreadItemListDataChatKitClientToolCall].
	CallID string `json:"call_id"`
	// This field is from variant [ChatKitThreadItemListDataChatKitClientToolCall].
	Name string `json:"name"`
	// This field is from variant [ChatKitThreadItemListDataChatKitClientToolCall].
	Output string `json:"output"`
	// This field is from variant [ChatKitThreadItemListDataChatKitClientToolCall].
	Status string `json:"status"`
	// This field is from variant [ChatKitThreadItemListDataChatKitTask].
	Heading string `json:"heading"`
	// This field is from variant [ChatKitThreadItemListDataChatKitTask].
	Summary string `json:"summary"`
	// This field is from variant [ChatKitThreadItemListDataChatKitTask].
	TaskType string `json:"task_type"`
	// This field is from variant [ChatKitThreadItemListDataChatKitTaskGroup].
	Tasks []ChatKitThreadItemListDataChatKitTaskGroupTask `json:"tasks"`
	JSON  struct {
		ID               respjson.Field
		Attachments      respjson.Field
		Content          respjson.Field
		CreatedAt        respjson.Field
		InferenceOptions respjson.Field
		Object           respjson.Field
		ThreadID         respjson.Field
		Type             respjson.Field
		Widget           respjson.Field
		Arguments        respjson.Field
		CallID           respjson.Field
		Name             respjson.Field
		Output           respjson.Field
		Status           respjson.Field
		Heading          respjson.Field
		Summary          respjson.Field
		TaskType         respjson.Field
		Tasks            respjson.Field
		raw              string
	} `json:"-"`
}

// anyChatKitThreadItemListData is implemented by each variant of
// [ChatKitThreadItemListDataUnion] to add type safety for the return type of
// [ChatKitThreadItemListDataUnion.AsAny]
type anyChatKitThreadItemListData interface {
	implChatKitThreadItemListDataUnion()
}

func (ChatKitThreadUserMessageItem) implChatKitThreadItemListDataUnion()                   {}
func (ChatKitThreadAssistantMessageItem) implChatKitThreadItemListDataUnion()              {}
func (ChatKitWidgetItem) implChatKitThreadItemListDataUnion()                              {}
func (ChatKitThreadItemListDataChatKitClientToolCall) implChatKitThreadItemListDataUnion() {}
func (ChatKitThreadItemListDataChatKitTask) implChatKitThreadItemListDataUnion()           {}
func (ChatKitThreadItemListDataChatKitTaskGroup) implChatKitThreadItemListDataUnion()      {}

// Use the following switch statement to find the correct variant
//
//	switch variant := ChatKitThreadItemListDataUnion.AsAny().(type) {
//	case openai.ChatKitThreadUserMessageItem:
//	case openai.ChatKitThreadAssistantMessageItem:
//	case openai.ChatKitWidgetItem:
//	case openai.ChatKitThreadItemListDataChatKitClientToolCall:
//	case openai.ChatKitThreadItemListDataChatKitTask:
//	case openai.ChatKitThreadItemListDataChatKitTaskGroup:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ChatKitThreadItemListDataUnion) AsAny() anyChatKitThreadItemListData {
	switch u.Type {
	case "chatkit.user_message":
		return u.AsChatKitUserMessage()
	case "chatkit.assistant_message":
		return u.AsChatKitAssistantMessage()
	case "chatkit.widget":
		return u.AsChatKitWidget()
	case "chatkit.client_tool_call":
		return u.AsChatKitClientToolCall()
	case "chatkit.task":
		return u.AsChatKitTask()
	case "chatkit.task_group":
		return u.AsChatKitTaskGroup()
	}
	return nil
}

func (u ChatKitThreadItemListDataUnion) AsChatKitUserMessage() (v ChatKitThreadUserMessageItem) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ChatKitThreadItemListDataUnion) AsChatKitAssistantMessage() (v ChatKitThreadAssistantMessageItem) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ChatKitThreadItemListDataUnion) AsChatKitWidget() (v ChatKitWidgetItem) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ChatKitThreadItemListDataUnion) AsChatKitClientToolCall() (v ChatKitThreadItemListDataChatKitClientToolCall) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ChatKitThreadItemListDataUnion) AsChatKitTask() (v ChatKitThreadItemListDataChatKitTask) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ChatKitThreadItemListDataUnion) AsChatKitTaskGroup() (v ChatKitThreadItemListDataChatKitTaskGroup) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ChatKitThreadItemListDataUnion) RawJSON() string { return u.JSON.raw }

func (r *ChatKitThreadItemListDataUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ChatKitThreadItemListDataUnionContent is an implicit subunion of
// [ChatKitThreadItemListDataUnion]. ChatKitThreadItemListDataUnionContent provides
// convenient access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [ChatKitThreadItemListDataUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfChatKitThreadUserMessageItemContentArray
// OfChatKitResponseOutputTextArray]
type ChatKitThreadItemListDataUnionContent struct {
	// This field will be present if the value is a
	// [[]ChatKitThreadUserMessageItemContentUnion] instead of an object.
	OfChatKitThreadUserMessageItemContentArray []ChatKitThreadUserMessageItemContentUnion `json:",inline"`
	// This field will be present if the value is a [[]ChatKitResponseOutputText]
	// instead of an object.
	OfChatKitResponseOutputTextArray []ChatKitResponseOutputText `json:",inline"`
	JSON                             struct {
		OfChatKitThreadUserMessageItemContentArray respjson.Field
		OfChatKitResponseOutputTextArray           respjson.Field
		raw                                        string
	} `json:"-"`
}

func (r *ChatKitThreadItemListDataUnionContent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Record of a client side tool invocation initiated by the assistant.
type ChatKitThreadItemListDataChatKitClientToolCall struct {
	// Identifier of the thread item.
	ID string `json:"id,required"`
	// JSON-encoded arguments that were sent to the tool.
	Arguments string `json:"arguments,required"`
	// Identifier for the client tool call.
	CallID string `json:"call_id,required"`
	// Unix timestamp (in seconds) for when the item was created.
	CreatedAt int64 `json:"created_at,required"`
	// Tool name that was invoked.
	Name string `json:"name,required"`
	// Type discriminator that is always `chatkit.thread_item`.
	Object constant.ChatKitThreadItem `json:"object,required"`
	// JSON-encoded output captured from the tool. Defaults to null while execution is
	// in progress.
	Output string `json:"output,required"`
	// Execution status for the tool call.
	//
	// Any of "in_progress", "completed".
	Status string `json:"status,required"`
	// Identifier of the parent thread.
	ThreadID string `json:"thread_id,required"`
	// Type discriminator that is always `chatkit.client_tool_call`.
	Type constant.ChatKitClientToolCall `json:"type,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Arguments   respjson.Field
		CallID      respjson.Field
		CreatedAt   respjson.Field
		Name        respjson.Field
		Object      respjson.Field
		Output      respjson.Field
		Status      respjson.Field
		ThreadID    respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ChatKitThreadItemListDataChatKitClientToolCall) RawJSON() string { return r.JSON.raw }
func (r *ChatKitThreadItemListDataChatKitClientToolCall) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Task emitted by the workflow to show progress and status updates.
type ChatKitThreadItemListDataChatKitTask struct {
	// Identifier of the thread item.
	ID string `json:"id,required"`
	// Unix timestamp (in seconds) for when the item was created.
	CreatedAt int64 `json:"created_at,required"`
	// Optional heading for the task. Defaults to null when not provided.
	Heading string `json:"heading,required"`
	// Type discriminator that is always `chatkit.thread_item`.
	Object constant.ChatKitThreadItem `json:"object,required"`
	// Optional summary that describes the task. Defaults to null when omitted.
	Summary string `json:"summary,required"`
	// Subtype for the task.
	//
	// Any of "custom", "thought".
	TaskType string `json:"task_type,required"`
	// Identifier of the parent thread.
	ThreadID string `json:"thread_id,required"`
	// Type discriminator that is always `chatkit.task`.
	Type constant.ChatKitTask `json:"type,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		Heading     respjson.Field
		Object      respjson.Field
		Summary     respjson.Field
		TaskType    respjson.Field
		ThreadID    respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ChatKitThreadItemListDataChatKitTask) RawJSON() string { return r.JSON.raw }
func (r *ChatKitThreadItemListDataChatKitTask) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Collection of workflow tasks grouped together in the thread.
type ChatKitThreadItemListDataChatKitTaskGroup struct {
	// Identifier of the thread item.
	ID string `json:"id,required"`
	// Unix timestamp (in seconds) for when the item was created.
	CreatedAt int64 `json:"created_at,required"`
	// Type discriminator that is always `chatkit.thread_item`.
	Object constant.ChatKitThreadItem `json:"object,required"`
	// Tasks included in the group.
	Tasks []ChatKitThreadItemListDataChatKitTaskGroupTask `json:"tasks,required"`
	// Identifier of the parent thread.
	ThreadID string `json:"thread_id,required"`
	// Type discriminator that is always `chatkit.task_group`.
	Type constant.ChatKitTaskGroup `json:"type,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		Object      respjson.Field
		Tasks       respjson.Field
		ThreadID    respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ChatKitThreadItemListDataChatKitTaskGroup) RawJSON() string { return r.JSON.raw }
func (r *ChatKitThreadItemListDataChatKitTaskGroup) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Task entry that appears within a TaskGroup.
type ChatKitThreadItemListDataChatKitTaskGroupTask struct {
	// Optional heading for the grouped task. Defaults to null when not provided.
	Heading string `json:"heading,required"`
	// Optional summary that describes the grouped task. Defaults to null when omitted.
	Summary string `json:"summary,required"`
	// Subtype for the grouped task.
	//
	// Any of "custom", "thought".
	Type string `json:"type,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Heading     respjson.Field
		Summary     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ChatKitThreadItemListDataChatKitTaskGroupTask) RawJSON() string { return r.JSON.raw }
func (r *ChatKitThreadItemListDataChatKitTaskGroupTask) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// User-authored messages within a thread.
type ChatKitThreadUserMessageItem struct {
	// Identifier of the thread item.
	ID string `json:"id,required"`
	// Attachments associated with the user message. Defaults to an empty list.
	Attachments []ChatKitAttachment `json:"attachments,required"`
	// Ordered content elements supplied by the user.
	Content []ChatKitThreadUserMessageItemContentUnion `json:"content,required"`
	// Unix timestamp (in seconds) for when the item was created.
	CreatedAt int64 `json:"created_at,required"`
	// Inference overrides applied to the message. Defaults to null when unset.
	InferenceOptions ChatKitThreadUserMessageItemInferenceOptions `json:"inference_options,required"`
	// Type discriminator that is always `chatkit.thread_item`.
	Object constant.ChatKitThreadItem `json:"object,required"`
	// Identifier of the parent thread.
	ThreadID string                      `json:"thread_id,required"`
	Type     constant.ChatKitUserMessage `json:"type,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID               respjson.Field
		Attachments      respjson.Field
		Content          respjson.Field
		CreatedAt        respjson.Field
		InferenceOptions respjson.Field
		Object           respjson.Field
		ThreadID         respjson.Field
		Type             respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ChatKitThreadUserMessageItem) RawJSON() string { return r.JSON.raw }
func (r *ChatKitThreadUserMessageItem) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ChatKitThreadUserMessageItemContentUnion contains all possible properties and
// values from [ChatKitThreadUserMessageItemContentInputText],
// [ChatKitThreadUserMessageItemContentQuotedText].
//
// Use the [ChatKitThreadUserMessageItemContentUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ChatKitThreadUserMessageItemContentUnion struct {
	Text string `json:"text"`
	// Any of "input_text", "quoted_text".
	Type string `json:"type"`
	JSON struct {
		Text respjson.Field
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyChatKitThreadUserMessageItemContent is implemented by each variant of
// [ChatKitThreadUserMessageItemContentUnion] to add type safety for the return
// type of [ChatKitThreadUserMessageItemContentUnion.AsAny]
type anyChatKitThreadUserMessageItemContent interface {
	implChatKitThreadUserMessageItemContentUnion()
}

func (ChatKitThreadUserMessageItemContentInputText) implChatKitThreadUserMessageItemContentUnion()  {}
func (ChatKitThreadUserMessageItemContentQuotedText) implChatKitThreadUserMessageItemContentUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := ChatKitThreadUserMessageItemContentUnion.AsAny().(type) {
//	case openai.ChatKitThreadUserMessageItemContentInputText:
//	case openai.ChatKitThreadUserMessageItemContentQuotedText:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ChatKitThreadUserMessageItemContentUnion) AsAny() anyChatKitThreadUserMessageItemContent {
	switch u.Type {
	case "input_text":
		return u.AsInputText()
	case "quoted_text":
		return u.AsQuotedText()
	}
	return nil
}

func (u ChatKitThreadUserMessageItemContentUnion) AsInputText() (v ChatKitThreadUserMessageItemContentInputText) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ChatKitThreadUserMessageItemContentUnion) AsQuotedText() (v ChatKitThreadUserMessageItemContentQuotedText) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ChatKitThreadUserMessageItemContentUnion) RawJSON() string { return u.JSON.raw }

func (r *ChatKitThreadUserMessageItemContentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Text block that a user contributed to the thread.
type ChatKitThreadUserMessageItemContentInputText struct {
	// Plain-text content supplied by the user.
	Text string `json:"text,required"`
	// Type discriminator that is always `input_text`.
	Type constant.InputText `json:"type,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Text        respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ChatKitThreadUserMessageItemContentInputText) RawJSON() string { return r.JSON.raw }
func (r *ChatKitThreadUserMessageItemContentInputText) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Quoted snippet that the user referenced in their message.
type ChatKitThreadUserMessageItemContentQuotedText struct {
	// Quoted text content.
	Text string `json:"text,required"`
	// Type discriminator that is always `quoted_text`.
	Type constant.QuotedText `json:"type,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Text        respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ChatKitThreadUserMessageItemContentQuotedText) RawJSON() string { return r.JSON.raw }
func (r *ChatKitThreadUserMessageItemContentQuotedText) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Inference overrides applied to the message. Defaults to null when unset.
type ChatKitThreadUserMessageItemInferenceOptions struct {
	// Model name that generated the response. Defaults to null when using the session
	// default.
	Model string `json:"model,required"`
	// Preferred tool to invoke. Defaults to null when ChatKit should auto-select.
	ToolChoice ChatKitThreadUserMessageItemInferenceOptionsToolChoice `json:"tool_choice,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Model       respjson.Field
		ToolChoice  respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ChatKitThreadUserMessageItemInferenceOptions) RawJSON() string { return r.JSON.raw }
func (r *ChatKitThreadUserMessageItemInferenceOptions) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Preferred tool to invoke. Defaults to null when ChatKit should auto-select.
type ChatKitThreadUserMessageItemInferenceOptionsToolChoice struct {
	// Identifier of the requested tool.
	ID string `json:"id,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ChatKitThreadUserMessageItemInferenceOptionsToolChoice) RawJSON() string { return r.JSON.raw }
func (r *ChatKitThreadUserMessageItemInferenceOptionsToolChoice) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Thread item that renders a widget payload.
type ChatKitWidgetItem struct {
	// Identifier of the thread item.
	ID string `json:"id,required"`
	// Unix timestamp (in seconds) for when the item was created.
	CreatedAt int64 `json:"created_at,required"`
	// Type discriminator that is always `chatkit.thread_item`.
	Object constant.ChatKitThreadItem `json:"object,required"`
	// Identifier of the parent thread.
	ThreadID string `json:"thread_id,required"`
	// Type discriminator that is always `chatkit.widget`.
	Type constant.ChatKitWidget `json:"type,required"`
	// Serialized widget payload rendered in the UI.
	Widget string `json:"widget,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		Object      respjson.Field
		ThreadID    respjson.Field
		Type        respjson.Field
		Widget      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ChatKitWidgetItem) RawJSON() string { return r.JSON.raw }
func (r *ChatKitWidgetItem) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Confirmation payload returned after deleting a thread.
type BetaChatKitThreadDeleteResponse struct {
	// Identifier of the deleted thread.
	ID string `json:"id,required"`
	// Indicates that the thread has been deleted.
	Deleted bool `json:"deleted,required"`
	// Type discriminator that is always `chatkit.thread.deleted`.
	Object constant.ChatKitThreadDeleted `json:"object,required"`
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
func (r BetaChatKitThreadDeleteResponse) RawJSON() string { return r.JSON.raw }
func (r *BetaChatKitThreadDeleteResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaChatKitThreadListParams struct {
	// List items created after this thread item ID. Defaults to null for the first
	// page.
	After param.Opt[string] `query:"after,omitzero" json:"-"`
	// List items created before this thread item ID. Defaults to null for the newest
	// results.
	Before param.Opt[string] `query:"before,omitzero" json:"-"`
	// Maximum number of thread items to return. Defaults to 20.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Filter threads that belong to this user identifier. Defaults to null to return
	// all users.
	User param.Opt[string] `query:"user,omitzero" json:"-"`
	// Sort order for results by creation time. Defaults to `desc`.
	//
	// Any of "asc", "desc".
	Order BetaChatKitThreadListParamsOrder `query:"order,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaChatKitThreadListParams]'s query parameters as
// `url.Values`.
func (r BetaChatKitThreadListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Sort order for results by creation time. Defaults to `desc`.
type BetaChatKitThreadListParamsOrder string

const (
	BetaChatKitThreadListParamsOrderAsc  BetaChatKitThreadListParamsOrder = "asc"
	BetaChatKitThreadListParamsOrderDesc BetaChatKitThreadListParamsOrder = "desc"
)

type BetaChatKitThreadListItemsParams struct {
	// List items created after this thread item ID. Defaults to null for the first
	// page.
	After param.Opt[string] `query:"after,omitzero" json:"-"`
	// List items created before this thread item ID. Defaults to null for the newest
	// results.
	Before param.Opt[string] `query:"before,omitzero" json:"-"`
	// Maximum number of thread items to return. Defaults to 20.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Sort order for results by creation time. Defaults to `desc`.
	//
	// Any of "asc", "desc".
	Order BetaChatKitThreadListItemsParamsOrder `query:"order,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaChatKitThreadListItemsParams]'s query parameters as
// `url.Values`.
func (r BetaChatKitThreadListItemsParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Sort order for results by creation time. Defaults to `desc`.
type BetaChatKitThreadListItemsParamsOrder string

const (
	BetaChatKitThreadListItemsParamsOrderAsc  BetaChatKitThreadListItemsParamsOrder = "asc"
	BetaChatKitThreadListItemsParamsOrderDesc BetaChatKitThreadListItemsParamsOrder = "desc"
)

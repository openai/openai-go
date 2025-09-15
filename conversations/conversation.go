// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package conversations

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/openai/openai-go/v2/internal/apijson"
	"github.com/openai/openai-go/v2/internal/requestconfig"
	"github.com/openai/openai-go/v2/option"
	"github.com/openai/openai-go/v2/packages/param"
	"github.com/openai/openai-go/v2/packages/respjson"
	"github.com/openai/openai-go/v2/responses"
	"github.com/openai/openai-go/v2/shared"
	"github.com/openai/openai-go/v2/shared/constant"
)

// ConversationService contains methods and other services that help with
// interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewConversationService] method instead.
type ConversationService struct {
	Options []option.RequestOption
	Items   ItemService
}

// NewConversationService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewConversationService(opts ...option.RequestOption) (r ConversationService) {
	r = ConversationService{}
	r.Options = opts
	r.Items = NewItemService(opts...)
	return
}

// Create a conversation.
func (r *ConversationService) New(ctx context.Context, body ConversationNewParams, opts ...option.RequestOption) (res *Conversation, err error) {
	opts = append(r.Options[:], opts...)
	path := "conversations"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Get a conversation with the given ID.
func (r *ConversationService) Get(ctx context.Context, conversationID string, opts ...option.RequestOption) (res *Conversation, err error) {
	opts = append(r.Options[:], opts...)
	if conversationID == "" {
		err = errors.New("missing required conversation_id parameter")
		return
	}
	path := fmt.Sprintf("conversations/%s", conversationID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// Update a conversation's metadata with the given ID.
func (r *ConversationService) Update(ctx context.Context, conversationID string, body ConversationUpdateParams, opts ...option.RequestOption) (res *Conversation, err error) {
	opts = append(r.Options[:], opts...)
	if conversationID == "" {
		err = errors.New("missing required conversation_id parameter")
		return
	}
	path := fmt.Sprintf("conversations/%s", conversationID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Delete a conversation with the given ID.
func (r *ConversationService) Delete(ctx context.Context, conversationID string, opts ...option.RequestOption) (res *ConversationDeletedResource, err error) {
	opts = append(r.Options[:], opts...)
	if conversationID == "" {
		err = errors.New("missing required conversation_id parameter")
		return
	}
	path := fmt.Sprintf("conversations/%s", conversationID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return
}

type ComputerScreenshotContent struct {
	// The identifier of an uploaded file that contains the screenshot.
	FileID string `json:"file_id,required"`
	// The URL of the screenshot image.
	ImageURL string `json:"image_url,required"`
	// Specifies the event type. For a computer screenshot, this property is always set
	// to `computer_screenshot`.
	Type constant.ComputerScreenshot `json:"type,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		FileID      respjson.Field
		ImageURL    respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ComputerScreenshotContent) RawJSON() string { return r.JSON.raw }
func (r *ComputerScreenshotContent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ContainerFileCitationBody struct {
	// The ID of the container file.
	ContainerID string `json:"container_id,required"`
	// The index of the last character of the container file citation in the message.
	EndIndex int64 `json:"end_index,required"`
	// The ID of the file.
	FileID string `json:"file_id,required"`
	// The filename of the container file cited.
	Filename string `json:"filename,required"`
	// The index of the first character of the container file citation in the message.
	StartIndex int64 `json:"start_index,required"`
	// The type of the container file citation. Always `container_file_citation`.
	Type constant.ContainerFileCitation `json:"type,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ContainerID respjson.Field
		EndIndex    respjson.Field
		FileID      respjson.Field
		Filename    respjson.Field
		StartIndex  respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ContainerFileCitationBody) RawJSON() string { return r.JSON.raw }
func (r *ContainerFileCitationBody) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type Conversation struct {
	// The unique ID of the conversation.
	ID string `json:"id,required"`
	// The time at which the conversation was created, measured in seconds since the
	// Unix epoch.
	CreatedAt int64 `json:"created_at,required"`
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard. Keys are strings with a maximum
	// length of 64 characters. Values are strings with a maximum length of 512
	// characters.
	Metadata any `json:"metadata,required"`
	// The object type, which is always `conversation`.
	Object constant.Conversation `json:"object,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		Metadata    respjson.Field
		Object      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r Conversation) RawJSON() string { return r.JSON.raw }
func (r *Conversation) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ConversationDeletedResource struct {
	ID      string                       `json:"id,required"`
	Deleted bool                         `json:"deleted,required"`
	Object  constant.ConversationDeleted `json:"object,required"`
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
func (r ConversationDeletedResource) RawJSON() string { return r.JSON.raw }
func (r *ConversationDeletedResource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type FileCitationBody struct {
	// The ID of the file.
	FileID string `json:"file_id,required"`
	// The filename of the file cited.
	Filename string `json:"filename,required"`
	// The index of the file in the list of files.
	Index int64 `json:"index,required"`
	// The type of the file citation. Always `file_citation`.
	Type constant.FileCitation `json:"type,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		FileID      respjson.Field
		Filename    respjson.Field
		Index       respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r FileCitationBody) RawJSON() string { return r.JSON.raw }
func (r *FileCitationBody) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type InputFileContent struct {
	// The ID of the file to be sent to the model.
	FileID string `json:"file_id,required"`
	// The type of the input item. Always `input_file`.
	Type constant.InputFile `json:"type,required"`
	// The URL of the file to be sent to the model.
	FileURL string `json:"file_url"`
	// The name of the file to be sent to the model.
	Filename string `json:"filename"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		FileID      respjson.Field
		Type        respjson.Field
		FileURL     respjson.Field
		Filename    respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r InputFileContent) RawJSON() string { return r.JSON.raw }
func (r *InputFileContent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type InputImageContent struct {
	// The detail level of the image to be sent to the model. One of `high`, `low`, or
	// `auto`. Defaults to `auto`.
	//
	// Any of "low", "high", "auto".
	Detail InputImageContentDetail `json:"detail,required"`
	// The ID of the file to be sent to the model.
	FileID string `json:"file_id,required"`
	// The URL of the image to be sent to the model. A fully qualified URL or base64
	// encoded image in a data URL.
	ImageURL string `json:"image_url,required"`
	// The type of the input item. Always `input_image`.
	Type constant.InputImage `json:"type,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Detail      respjson.Field
		FileID      respjson.Field
		ImageURL    respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r InputImageContent) RawJSON() string { return r.JSON.raw }
func (r *InputImageContent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The detail level of the image to be sent to the model. One of `high`, `low`, or
// `auto`. Defaults to `auto`.
type InputImageContentDetail string

const (
	InputImageContentDetailLow  InputImageContentDetail = "low"
	InputImageContentDetailHigh InputImageContentDetail = "high"
	InputImageContentDetailAuto InputImageContentDetail = "auto"
)

type InputTextContent struct {
	// The text input to the model.
	Text string `json:"text,required"`
	// The type of the input item. Always `input_text`.
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
func (r InputTextContent) RawJSON() string { return r.JSON.raw }
func (r *InputTextContent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type LobProb struct {
	Token       string       `json:"token,required"`
	Bytes       []int64      `json:"bytes,required"`
	Logprob     float64      `json:"logprob,required"`
	TopLogprobs []TopLogProb `json:"top_logprobs,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Token       respjson.Field
		Bytes       respjson.Field
		Logprob     respjson.Field
		TopLogprobs respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r LobProb) RawJSON() string { return r.JSON.raw }
func (r *LobProb) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type Message struct {
	// The unique ID of the message.
	ID string `json:"id,required"`
	// The content of the message
	Content []MessageContentUnion `json:"content,required"`
	// The role of the message. One of `unknown`, `user`, `assistant`, `system`,
	// `critic`, `discriminator`, `developer`, or `tool`.
	//
	// Any of "unknown", "user", "assistant", "system", "critic", "discriminator",
	// "developer", "tool".
	Role MessageRole `json:"role,required"`
	// The status of item. One of `in_progress`, `completed`, or `incomplete`.
	// Populated when items are returned via API.
	//
	// Any of "in_progress", "completed", "incomplete".
	Status MessageStatus `json:"status,required"`
	// The type of the message. Always set to `message`.
	Type constant.Message `json:"type,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Content     respjson.Field
		Role        respjson.Field
		Status      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r Message) RawJSON() string { return r.JSON.raw }
func (r *Message) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// MessageContentUnion contains all possible properties and values from
// [InputTextContent], [OutputTextContent], [TextContent], [SummaryTextContent],
// [RefusalContent], [InputImageContent], [ComputerScreenshotContent],
// [InputFileContent].
//
// Use the [MessageContentUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type MessageContentUnion struct {
	Text string `json:"text"`
	// Any of "input_text", "output_text", "text", "summary_text", "refusal",
	// "input_image", "computer_screenshot", "input_file".
	Type string `json:"type"`
	// This field is from variant [OutputTextContent].
	Annotations []OutputTextContentAnnotationUnion `json:"annotations"`
	// This field is from variant [OutputTextContent].
	Logprobs []LobProb `json:"logprobs"`
	// This field is from variant [RefusalContent].
	Refusal string `json:"refusal"`
	// This field is from variant [InputImageContent].
	Detail   InputImageContentDetail `json:"detail"`
	FileID   string                  `json:"file_id"`
	ImageURL string                  `json:"image_url"`
	// This field is from variant [InputFileContent].
	FileURL string `json:"file_url"`
	// This field is from variant [InputFileContent].
	Filename string `json:"filename"`
	JSON     struct {
		Text        respjson.Field
		Type        respjson.Field
		Annotations respjson.Field
		Logprobs    respjson.Field
		Refusal     respjson.Field
		Detail      respjson.Field
		FileID      respjson.Field
		ImageURL    respjson.Field
		FileURL     respjson.Field
		Filename    respjson.Field
		raw         string
	} `json:"-"`
}

// anyMessageContent is implemented by each variant of [MessageContentUnion] to add
// type safety for the return type of [MessageContentUnion.AsAny]
type anyMessageContent interface {
	implMessageContentUnion()
}

func (InputTextContent) implMessageContentUnion()          {}
func (OutputTextContent) implMessageContentUnion()         {}
func (TextContent) implMessageContentUnion()               {}
func (SummaryTextContent) implMessageContentUnion()        {}
func (RefusalContent) implMessageContentUnion()            {}
func (InputImageContent) implMessageContentUnion()         {}
func (ComputerScreenshotContent) implMessageContentUnion() {}
func (InputFileContent) implMessageContentUnion()          {}

// Use the following switch statement to find the correct variant
//
//	switch variant := MessageContentUnion.AsAny().(type) {
//	case conversations.InputTextContent:
//	case conversations.OutputTextContent:
//	case conversations.TextContent:
//	case conversations.SummaryTextContent:
//	case conversations.RefusalContent:
//	case conversations.InputImageContent:
//	case conversations.ComputerScreenshotContent:
//	case conversations.InputFileContent:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u MessageContentUnion) AsAny() anyMessageContent {
	switch u.Type {
	case "input_text":
		return u.AsInputText()
	case "output_text":
		return u.AsOutputText()
	case "text":
		return u.AsText()
	case "summary_text":
		return u.AsSummaryText()
	case "refusal":
		return u.AsRefusal()
	case "input_image":
		return u.AsInputImage()
	case "computer_screenshot":
		return u.AsComputerScreenshot()
	case "input_file":
		return u.AsInputFile()
	}
	return nil
}

func (u MessageContentUnion) AsInputText() (v InputTextContent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u MessageContentUnion) AsOutputText() (v OutputTextContent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u MessageContentUnion) AsText() (v TextContent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u MessageContentUnion) AsSummaryText() (v SummaryTextContent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u MessageContentUnion) AsRefusal() (v RefusalContent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u MessageContentUnion) AsInputImage() (v InputImageContent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u MessageContentUnion) AsComputerScreenshot() (v ComputerScreenshotContent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u MessageContentUnion) AsInputFile() (v InputFileContent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u MessageContentUnion) RawJSON() string { return u.JSON.raw }

func (r *MessageContentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The role of the message. One of `unknown`, `user`, `assistant`, `system`,
// `critic`, `discriminator`, `developer`, or `tool`.
type MessageRole string

const (
	MessageRoleUnknown       MessageRole = "unknown"
	MessageRoleUser          MessageRole = "user"
	MessageRoleAssistant     MessageRole = "assistant"
	MessageRoleSystem        MessageRole = "system"
	MessageRoleCritic        MessageRole = "critic"
	MessageRoleDiscriminator MessageRole = "discriminator"
	MessageRoleDeveloper     MessageRole = "developer"
	MessageRoleTool          MessageRole = "tool"
)

// The status of item. One of `in_progress`, `completed`, or `incomplete`.
// Populated when items are returned via API.
type MessageStatus string

const (
	MessageStatusInProgress MessageStatus = "in_progress"
	MessageStatusCompleted  MessageStatus = "completed"
	MessageStatusIncomplete MessageStatus = "incomplete"
)

type OutputTextContent struct {
	// The annotations of the text output.
	Annotations []OutputTextContentAnnotationUnion `json:"annotations,required"`
	// The text output from the model.
	Text string `json:"text,required"`
	// The type of the output text. Always `output_text`.
	Type     constant.OutputText `json:"type,required"`
	Logprobs []LobProb           `json:"logprobs"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Annotations respjson.Field
		Text        respjson.Field
		Type        respjson.Field
		Logprobs    respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r OutputTextContent) RawJSON() string { return r.JSON.raw }
func (r *OutputTextContent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// OutputTextContentAnnotationUnion contains all possible properties and values
// from [FileCitationBody], [URLCitationBody], [ContainerFileCitationBody].
//
// Use the [OutputTextContentAnnotationUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type OutputTextContentAnnotationUnion struct {
	FileID   string `json:"file_id"`
	Filename string `json:"filename"`
	// This field is from variant [FileCitationBody].
	Index int64 `json:"index"`
	// Any of "file_citation", "url_citation", "container_file_citation".
	Type       string `json:"type"`
	EndIndex   int64  `json:"end_index"`
	StartIndex int64  `json:"start_index"`
	// This field is from variant [URLCitationBody].
	Title string `json:"title"`
	// This field is from variant [URLCitationBody].
	URL string `json:"url"`
	// This field is from variant [ContainerFileCitationBody].
	ContainerID string `json:"container_id"`
	JSON        struct {
		FileID      respjson.Field
		Filename    respjson.Field
		Index       respjson.Field
		Type        respjson.Field
		EndIndex    respjson.Field
		StartIndex  respjson.Field
		Title       respjson.Field
		URL         respjson.Field
		ContainerID respjson.Field
		raw         string
	} `json:"-"`
}

// anyOutputTextContentAnnotation is implemented by each variant of
// [OutputTextContentAnnotationUnion] to add type safety for the return type of
// [OutputTextContentAnnotationUnion.AsAny]
type anyOutputTextContentAnnotation interface {
	implOutputTextContentAnnotationUnion()
}

func (FileCitationBody) implOutputTextContentAnnotationUnion()          {}
func (URLCitationBody) implOutputTextContentAnnotationUnion()           {}
func (ContainerFileCitationBody) implOutputTextContentAnnotationUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := OutputTextContentAnnotationUnion.AsAny().(type) {
//	case conversations.FileCitationBody:
//	case conversations.URLCitationBody:
//	case conversations.ContainerFileCitationBody:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u OutputTextContentAnnotationUnion) AsAny() anyOutputTextContentAnnotation {
	switch u.Type {
	case "file_citation":
		return u.AsFileCitation()
	case "url_citation":
		return u.AsURLCitation()
	case "container_file_citation":
		return u.AsContainerFileCitation()
	}
	return nil
}

func (u OutputTextContentAnnotationUnion) AsFileCitation() (v FileCitationBody) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u OutputTextContentAnnotationUnion) AsURLCitation() (v URLCitationBody) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u OutputTextContentAnnotationUnion) AsContainerFileCitation() (v ContainerFileCitationBody) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u OutputTextContentAnnotationUnion) RawJSON() string { return u.JSON.raw }

func (r *OutputTextContentAnnotationUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type RefusalContent struct {
	// The refusal explanation from the model.
	Refusal string `json:"refusal,required"`
	// The type of the refusal. Always `refusal`.
	Type constant.Refusal `json:"type,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Refusal     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RefusalContent) RawJSON() string { return r.JSON.raw }
func (r *RefusalContent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type SummaryTextContent struct {
	Text string               `json:"text,required"`
	Type constant.SummaryText `json:"type,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Text        respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r SummaryTextContent) RawJSON() string { return r.JSON.raw }
func (r *SummaryTextContent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type TextContent struct {
	Text string        `json:"text,required"`
	Type constant.Text `json:"type,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Text        respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r TextContent) RawJSON() string { return r.JSON.raw }
func (r *TextContent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type TopLogProb struct {
	Token   string  `json:"token,required"`
	Bytes   []int64 `json:"bytes,required"`
	Logprob float64 `json:"logprob,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Token       respjson.Field
		Bytes       respjson.Field
		Logprob     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r TopLogProb) RawJSON() string { return r.JSON.raw }
func (r *TopLogProb) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type URLCitationBody struct {
	// The index of the last character of the URL citation in the message.
	EndIndex int64 `json:"end_index,required"`
	// The index of the first character of the URL citation in the message.
	StartIndex int64 `json:"start_index,required"`
	// The title of the web resource.
	Title string `json:"title,required"`
	// The type of the URL citation. Always `url_citation`.
	Type constant.URLCitation `json:"type,required"`
	// The URL of the web resource.
	URL string `json:"url,required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EndIndex    respjson.Field
		StartIndex  respjson.Field
		Title       respjson.Field
		Type        respjson.Field
		URL         respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r URLCitationBody) RawJSON() string { return r.JSON.raw }
func (r *URLCitationBody) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ConversationNewParams struct {
	// Initial items to include in the conversation context. You may add up to 20 items
	// at a time.
	Items []responses.ResponseInputItemUnionParam `json:"items,omitzero"`
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard.
	//
	// Keys are strings with a maximum length of 64 characters. Values are strings with
	// a maximum length of 512 characters.
	Metadata shared.Metadata `json:"metadata,omitzero"`
	paramObj
}

func (r ConversationNewParams) MarshalJSON() (data []byte, err error) {
	type shadow ConversationNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ConversationNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ConversationUpdateParams struct {
	// Set of 16 key-value pairs that can be attached to an object. This can be useful
	// for storing additional information about the object in a structured format, and
	// querying for objects via API or the dashboard. Keys are strings with a maximum
	// length of 64 characters. Values are strings with a maximum length of 512
	// characters.
	Metadata map[string]string `json:"metadata,omitzero,required"`
	paramObj
}

func (r ConversationUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow ConversationUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ConversationUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai

import (
	"context"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/openai/openai-go/v3/conversations"
	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/realtime"
	"github.com/openai/openai-go/v3/responses"
	"github.com/openai/openai-go/v3/webhooks"
)

// Client creates a struct with services and top level methods that help with
// interacting with the openai API. You should not instantiate this client
// directly, and instead use the [NewClient] method instead.
type Client struct {
	Options []option.RequestOption
	// Given a prompt, the model will return one or more predicted completions, and can
	// also return the probabilities of alternative tokens at each position.
	Completions CompletionService
	Chat        ChatService
	// Get a vector representation of a given input that can be easily consumed by
	// machine learning models and algorithms.
	Embeddings EmbeddingService
	// Files are used to upload documents that can be used with features like
	// Assistants and Fine-tuning.
	Files FileService
	// Given a prompt and/or an input image, the model will generate a new image.
	Images ImageService
	Audio  AudioService
	// Given text and/or image inputs, classifies if those inputs are potentially
	// harmful.
	Moderations ModerationService
	// List and describe the various models available in the API.
	Models       ModelService
	FineTuning   FineTuningService
	Graders      GraderService
	VectorStores VectorStoreService
	Webhooks     webhooks.WebhookService
	Beta         BetaService
	// Create large batches of API requests to run asynchronously.
	Batches BatchService
	// Use Uploads to upload large files in multiple parts.
	Uploads   UploadService
	Responses responses.ResponseService
	Realtime  realtime.RealtimeService
	// Manage conversations and conversation items.
	Conversations conversations.ConversationService
	Containers    ContainerService
	Skills        SkillService
	Videos        VideoService
}

// DefaultClientOptions read from the environment (OPENAI_API_KEY, OPENAI_ORG_ID,
// OPENAI_PROJECT_ID, OPENAI_WEBHOOK_SECRET, OPENAI_BASE_URL). This should be used
// to initialize new clients.
func DefaultClientOptions() []option.RequestOption {
	defaults := []option.RequestOption{option.WithHTTPClient(defaultHTTPClient()), option.WithEnvironmentProduction()}
	if o, ok := os.LookupEnv("OPENAI_BASE_URL"); ok {
		defaults = append(defaults, option.WithBaseURL(o))
	}
	if o, ok := os.LookupEnv("OPENAI_API_KEY"); ok {
		defaults = append(defaults, option.WithAPIKey(o))
	}
	if o, ok := os.LookupEnv("OPENAI_ORG_ID"); ok {
		defaults = append(defaults, option.WithOrganization(o))
	}
	if o, ok := os.LookupEnv("OPENAI_PROJECT_ID"); ok {
		defaults = append(defaults, option.WithProject(o))
	}
	if o, ok := os.LookupEnv("OPENAI_WEBHOOK_SECRET"); ok {
		defaults = append(defaults, option.WithWebhookSecret(o))
	}
	if o, ok := os.LookupEnv("OPENAI_CUSTOM_HEADERS"); ok {
		for _, line := range strings.Split(o, "\n") {
			colon := strings.Index(line, ":")
			if colon >= 0 {
				defaults = append(defaults, option.WithHeader(strings.TrimSpace(line[:colon]), strings.TrimSpace(line[colon+1:])))
			}
		}
	}
	return defaults
}

// NewClient generates a new client with the default option read from the
// environment (OPENAI_API_KEY, OPENAI_ORG_ID, OPENAI_PROJECT_ID,
// OPENAI_WEBHOOK_SECRET, OPENAI_BASE_URL). The option passed in as arguments are
// applied after these default arguments, and all option will be passed down to the
// services and requests that this client makes.
func NewClient(opts ...option.RequestOption) (r Client) {
	opts = append(DefaultClientOptions(), opts...)

	r = Client{Options: opts}

	r.Completions = NewCompletionService(opts...)
	r.Chat = NewChatService(opts...)
	r.Embeddings = NewEmbeddingService(opts...)
	r.Files = NewFileService(opts...)
	r.Images = NewImageService(opts...)
	r.Audio = NewAudioService(opts...)
	r.Moderations = NewModerationService(opts...)
	r.Models = NewModelService(opts...)
	r.FineTuning = NewFineTuningService(opts...)
	r.Graders = NewGraderService(opts...)
	r.VectorStores = NewVectorStoreService(opts...)
	r.Webhooks = webhooks.NewWebhookService(opts...)
	r.Beta = NewBetaService(opts...)
	r.Batches = NewBatchService(opts...)
	r.Uploads = NewUploadService(opts...)
	r.Responses = responses.NewResponseService(opts...)
	r.Realtime = realtime.NewRealtimeService(opts...)
	r.Conversations = conversations.NewConversationService(opts...)
	r.Containers = NewContainerService(opts...)
	r.Skills = NewSkillService(opts...)
	r.Videos = NewVideoService(opts...)

	return
}

// Execute makes a request with the given context, method, URL, request params,
// response, and request options. This is useful for hitting undocumented endpoints
// while retaining the base URL, auth, retries, and other options from the client.
//
// If a byte slice or an [io.Reader] is supplied to params, it will be used as-is
// for the request body.
//
// The params is by default serialized into the body using [encoding/json]. If your
// type implements a MarshalJSON function, it will be used instead to serialize the
// request. If a URLQuery method is implemented, the returned [url.Values] will be
// used as query strings to the url.
//
// If your params struct uses [param.Field], you must provide either [MarshalJSON],
// [URLQuery], and/or [MarshalForm] functions. It is undefined behavior to use a
// struct uses [param.Field] without specifying how it is serialized.
//
// Any "…Params" object defined in this library can be used as the request
// argument. Note that 'path' arguments will not be forwarded into the url.
//
// The response body will be deserialized into the res variable, depending on its
// type:
//
//   - A pointer to a [*http.Response] is populated by the raw response.
//   - A pointer to a byte array will be populated with the contents of the request
//     body.
//   - A pointer to any other type uses this library's default JSON decoding, which
//     respects UnmarshalJSON if it is defined on the type.
//   - A nil value will not read the response body.
//
// For even greater flexibility, see [option.WithResponseInto] and
// [option.WithResponseBodyInto].
func (r *Client) Execute(ctx context.Context, method string, path string, params any, res any, opts ...option.RequestOption) error {
	opts = slices.Concat(r.Options, opts)
	return requestconfig.ExecuteNewRequest(ctx, method, path, params, res, opts...)
}

// Get makes a GET request with the given URL, params, and optionally deserializes
// to a response. See [Execute] documentation on the params and response.
func (r *Client) Get(ctx context.Context, path string, params any, res any, opts ...option.RequestOption) error {
	return r.Execute(ctx, http.MethodGet, path, params, res, opts...)
}

// Post makes a POST request with the given URL, params, and optionally
// deserializes to a response. See [Execute] documentation on the params and
// response.
func (r *Client) Post(ctx context.Context, path string, params any, res any, opts ...option.RequestOption) error {
	return r.Execute(ctx, http.MethodPost, path, params, res, opts...)
}

// Put makes a PUT request with the given URL, params, and optionally deserializes
// to a response. See [Execute] documentation on the params and response.
func (r *Client) Put(ctx context.Context, path string, params any, res any, opts ...option.RequestOption) error {
	return r.Execute(ctx, http.MethodPut, path, params, res, opts...)
}

// Patch makes a PATCH request with the given URL, params, and optionally
// deserializes to a response. See [Execute] documentation on the params and
// response.
func (r *Client) Patch(ctx context.Context, path string, params any, res any, opts ...option.RequestOption) error {
	return r.Execute(ctx, http.MethodPatch, path, params, res, opts...)
}

// Delete makes a DELETE request with the given URL, params, and optionally
// deserializes to a response. See [Execute] documentation on the params and
// response.
func (r *Client) Delete(ctx context.Context, path string, params any, res any, opts ...option.RequestOption) error {
	return r.Execute(ctx, http.MethodDelete, path, params, res, opts...)
}

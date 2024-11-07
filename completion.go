// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai

import (
	"context"
	"net/http"

	"github.com/openai/openai-go/internal/apijson"
	"github.com/openai/openai-go/internal/param"
	"github.com/openai/openai-go/internal/requestconfig"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/ssestream"
)

// CompletionService contains methods and other services that help with interacting
// with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewCompletionService] method instead.
type CompletionService struct {
	Options []option.RequestOption
}

// NewCompletionService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewCompletionService(opts ...option.RequestOption) (r *CompletionService) {
	r = &CompletionService{}
	r.Options = opts
	return
}

// Creates a completion for the provided prompt and parameters.
func (r *CompletionService) New(ctx context.Context, body CompletionNewParams, opts ...option.RequestOption) (res *Completion, err error) {
	opts = append(r.Options[:], opts...)
	path := "completions"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Creates a completion for the provided prompt and parameters.
func (r *CompletionService) NewStreaming(ctx context.Context, body CompletionNewParams, opts ...option.RequestOption) (stream *ssestream.Stream[Completion]) {
	var (
		raw *http.Response
		err error
	)
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithJSONSet("stream", true)}, opts...)
	path := "completions"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &raw, opts...)
	return ssestream.NewStream[Completion](ssestream.NewDecoder(raw), err)
}

// Represents a completion response from the API. Note: both the streamed and
// non-streamed response objects share the same shape (unlike the chat endpoint).
type Completion struct {
	// A unique identifier for the completion.
	ID string `json:"id,required"`
	// The list of completion choices the model generated for the input prompt.
	Choices []CompletionChoice `json:"choices,required"`
	// The Unix timestamp (in seconds) of when the completion was created.
	Created int64 `json:"created,required"`
	// The model used for completion.
	Model string `json:"model,required"`
	// The object type, which is always "text_completion"
	Object CompletionObject `json:"object,required"`
	// This fingerprint represents the backend configuration that the model runs with.
	//
	// Can be used in conjunction with the `seed` request parameter to understand when
	// backend changes have been made that might impact determinism.
	SystemFingerprint string `json:"system_fingerprint"`
	// Usage statistics for the completion request.
	Usage CompletionUsage `json:"usage"`
	JSON  completionJSON  `json:"-"`
}

// completionJSON contains the JSON metadata for the struct [Completion]
type completionJSON struct {
	ID                apijson.Field
	Choices           apijson.Field
	Created           apijson.Field
	Model             apijson.Field
	Object            apijson.Field
	SystemFingerprint apijson.Field
	Usage             apijson.Field
	raw               string
	ExtraFields       map[string]apijson.Field
}

func (r *Completion) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r completionJSON) RawJSON() string {
	return r.raw
}

// The object type, which is always "text_completion"
type CompletionObject string

const (
	CompletionObjectTextCompletion CompletionObject = "text_completion"
)

func (r CompletionObject) IsKnown() bool {
	switch r {
	case CompletionObjectTextCompletion:
		return true
	}
	return false
}

type CompletionChoice struct {
	// The reason the model stopped generating tokens. This will be `stop` if the model
	// hit a natural stop point or a provided stop sequence, `length` if the maximum
	// number of tokens specified in the request was reached, or `content_filter` if
	// content was omitted due to a flag from our content filters.
	FinishReason CompletionChoiceFinishReason `json:"finish_reason,required"`
	Index        int64                        `json:"index,required"`
	Logprobs     CompletionChoiceLogprobs     `json:"logprobs,required,nullable"`
	Text         string                       `json:"text,required"`
	JSON         completionChoiceJSON         `json:"-"`
}

// completionChoiceJSON contains the JSON metadata for the struct
// [CompletionChoice]
type completionChoiceJSON struct {
	FinishReason apijson.Field
	Index        apijson.Field
	Logprobs     apijson.Field
	Text         apijson.Field
	raw          string
	ExtraFields  map[string]apijson.Field
}

func (r *CompletionChoice) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r completionChoiceJSON) RawJSON() string {
	return r.raw
}

// The reason the model stopped generating tokens. This will be `stop` if the model
// hit a natural stop point or a provided stop sequence, `length` if the maximum
// number of tokens specified in the request was reached, or `content_filter` if
// content was omitted due to a flag from our content filters.
type CompletionChoiceFinishReason string

const (
	CompletionChoiceFinishReasonStop          CompletionChoiceFinishReason = "stop"
	CompletionChoiceFinishReasonLength        CompletionChoiceFinishReason = "length"
	CompletionChoiceFinishReasonContentFilter CompletionChoiceFinishReason = "content_filter"
)

func (r CompletionChoiceFinishReason) IsKnown() bool {
	switch r {
	case CompletionChoiceFinishReasonStop, CompletionChoiceFinishReasonLength, CompletionChoiceFinishReasonContentFilter:
		return true
	}
	return false
}

type CompletionChoiceLogprobs struct {
	TextOffset    []int64                      `json:"text_offset"`
	TokenLogprobs []float64                    `json:"token_logprobs"`
	Tokens        []string                     `json:"tokens"`
	TopLogprobs   []map[string]float64         `json:"top_logprobs"`
	JSON          completionChoiceLogprobsJSON `json:"-"`
}

// completionChoiceLogprobsJSON contains the JSON metadata for the struct
// [CompletionChoiceLogprobs]
type completionChoiceLogprobsJSON struct {
	TextOffset    apijson.Field
	TokenLogprobs apijson.Field
	Tokens        apijson.Field
	TopLogprobs   apijson.Field
	raw           string
	ExtraFields   map[string]apijson.Field
}

func (r *CompletionChoiceLogprobs) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r completionChoiceLogprobsJSON) RawJSON() string {
	return r.raw
}

// Usage statistics for the completion request.
type CompletionUsage struct {
	// Number of tokens in the generated completion.
	CompletionTokens int64 `json:"completion_tokens,required"`
	// Number of tokens in the prompt.
	PromptTokens int64 `json:"prompt_tokens,required"`
	// Total number of tokens used in the request (prompt + completion).
	TotalTokens int64 `json:"total_tokens,required"`
	// Breakdown of tokens used in a completion.
	CompletionTokensDetails CompletionUsageCompletionTokensDetails `json:"completion_tokens_details"`
	// Breakdown of tokens used in the prompt.
	PromptTokensDetails CompletionUsagePromptTokensDetails `json:"prompt_tokens_details"`
	JSON                completionUsageJSON                `json:"-"`
}

// completionUsageJSON contains the JSON metadata for the struct [CompletionUsage]
type completionUsageJSON struct {
	CompletionTokens        apijson.Field
	PromptTokens            apijson.Field
	TotalTokens             apijson.Field
	CompletionTokensDetails apijson.Field
	PromptTokensDetails     apijson.Field
	raw                     string
	ExtraFields             map[string]apijson.Field
}

func (r *CompletionUsage) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r completionUsageJSON) RawJSON() string {
	return r.raw
}

// Breakdown of tokens used in a completion.
type CompletionUsageCompletionTokensDetails struct {
	// When using Predicted Outputs, the number of tokens in the prediction that
	// appeared in the completion.
	AcceptedPredictionTokens int64 `json:"accepted_prediction_tokens"`
	// Audio input tokens generated by the model.
	AudioTokens int64 `json:"audio_tokens"`
	// Tokens generated by the model for reasoning.
	ReasoningTokens int64 `json:"reasoning_tokens"`
	// When using Predicted Outputs, the number of tokens in the prediction that did
	// not appear in the completion. However, like reasoning tokens, these tokens are
	// still counted in the total completion tokens for purposes of billing, output,
	// and context window limits.
	RejectedPredictionTokens int64                                      `json:"rejected_prediction_tokens"`
	JSON                     completionUsageCompletionTokensDetailsJSON `json:"-"`
}

// completionUsageCompletionTokensDetailsJSON contains the JSON metadata for the
// struct [CompletionUsageCompletionTokensDetails]
type completionUsageCompletionTokensDetailsJSON struct {
	AcceptedPredictionTokens apijson.Field
	AudioTokens              apijson.Field
	ReasoningTokens          apijson.Field
	RejectedPredictionTokens apijson.Field
	raw                      string
	ExtraFields              map[string]apijson.Field
}

func (r *CompletionUsageCompletionTokensDetails) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r completionUsageCompletionTokensDetailsJSON) RawJSON() string {
	return r.raw
}

// Breakdown of tokens used in the prompt.
type CompletionUsagePromptTokensDetails struct {
	// Audio input tokens present in the prompt.
	AudioTokens int64 `json:"audio_tokens"`
	// Cached tokens present in the prompt.
	CachedTokens int64                                  `json:"cached_tokens"`
	JSON         completionUsagePromptTokensDetailsJSON `json:"-"`
}

// completionUsagePromptTokensDetailsJSON contains the JSON metadata for the struct
// [CompletionUsagePromptTokensDetails]
type completionUsagePromptTokensDetailsJSON struct {
	AudioTokens  apijson.Field
	CachedTokens apijson.Field
	raw          string
	ExtraFields  map[string]apijson.Field
}

func (r *CompletionUsagePromptTokensDetails) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r completionUsagePromptTokensDetailsJSON) RawJSON() string {
	return r.raw
}

type CompletionNewParams struct {
	// ID of the model to use. You can use the
	// [List models](https://platform.openai.com/docs/api-reference/models/list) API to
	// see all of your available models, or see our
	// [Model overview](https://platform.openai.com/docs/models) for descriptions of
	// them.
	Model param.Field[CompletionNewParamsModel] `json:"model,required"`
	// The prompt(s) to generate completions for, encoded as a string, array of
	// strings, array of tokens, or array of token arrays.
	//
	// Note that <|endoftext|> is the document separator that the model sees during
	// training, so if a prompt is not specified the model will generate as if from the
	// beginning of a new document.
	Prompt param.Field[CompletionNewParamsPromptUnion] `json:"prompt,required"`
	// Generates `best_of` completions server-side and returns the "best" (the one with
	// the highest log probability per token). Results cannot be streamed.
	//
	// When used with `n`, `best_of` controls the number of candidate completions and
	// `n` specifies how many to return – `best_of` must be greater than `n`.
	//
	// **Note:** Because this parameter generates many completions, it can quickly
	// consume your token quota. Use carefully and ensure that you have reasonable
	// settings for `max_tokens` and `stop`.
	BestOf param.Field[int64] `json:"best_of"`
	// Echo back the prompt in addition to the completion
	Echo param.Field[bool] `json:"echo"`
	// Number between -2.0 and 2.0. Positive values penalize new tokens based on their
	// existing frequency in the text so far, decreasing the model's likelihood to
	// repeat the same line verbatim.
	//
	// [See more information about frequency and presence penalties.](https://platform.openai.com/docs/guides/text-generation)
	FrequencyPenalty param.Field[float64] `json:"frequency_penalty"`
	// Modify the likelihood of specified tokens appearing in the completion.
	//
	// Accepts a JSON object that maps tokens (specified by their token ID in the GPT
	// tokenizer) to an associated bias value from -100 to 100. You can use this
	// [tokenizer tool](/tokenizer?view=bpe) to convert text to token IDs.
	// Mathematically, the bias is added to the logits generated by the model prior to
	// sampling. The exact effect will vary per model, but values between -1 and 1
	// should decrease or increase likelihood of selection; values like -100 or 100
	// should result in a ban or exclusive selection of the relevant token.
	//
	// As an example, you can pass `{"50256": -100}` to prevent the <|endoftext|> token
	// from being generated.
	LogitBias param.Field[map[string]int64] `json:"logit_bias"`
	// Include the log probabilities on the `logprobs` most likely output tokens, as
	// well the chosen tokens. For example, if `logprobs` is 5, the API will return a
	// list of the 5 most likely tokens. The API will always return the `logprob` of
	// the sampled token, so there may be up to `logprobs+1` elements in the response.
	//
	// The maximum value for `logprobs` is 5.
	Logprobs param.Field[int64] `json:"logprobs"`
	// The maximum number of [tokens](/tokenizer) that can be generated in the
	// completion.
	//
	// The token count of your prompt plus `max_tokens` cannot exceed the model's
	// context length.
	// [Example Python code](https://cookbook.openai.com/examples/how_to_count_tokens_with_tiktoken)
	// for counting tokens.
	MaxTokens param.Field[int64] `json:"max_tokens"`
	// How many completions to generate for each prompt.
	//
	// **Note:** Because this parameter generates many completions, it can quickly
	// consume your token quota. Use carefully and ensure that you have reasonable
	// settings for `max_tokens` and `stop`.
	N param.Field[int64] `json:"n"`
	// Number between -2.0 and 2.0. Positive values penalize new tokens based on
	// whether they appear in the text so far, increasing the model's likelihood to
	// talk about new topics.
	//
	// [See more information about frequency and presence penalties.](https://platform.openai.com/docs/guides/text-generation)
	PresencePenalty param.Field[float64] `json:"presence_penalty"`
	// If specified, our system will make a best effort to sample deterministically,
	// such that repeated requests with the same `seed` and parameters should return
	// the same result.
	//
	// Determinism is not guaranteed, and you should refer to the `system_fingerprint`
	// response parameter to monitor changes in the backend.
	Seed param.Field[int64] `json:"seed"`
	// Up to 4 sequences where the API will stop generating further tokens. The
	// returned text will not contain the stop sequence.
	Stop param.Field[CompletionNewParamsStopUnion] `json:"stop"`
	// Options for streaming response. Only set this when you set `stream: true`.
	StreamOptions param.Field[ChatCompletionStreamOptionsParam] `json:"stream_options"`
	// The suffix that comes after a completion of inserted text.
	//
	// This parameter is only supported for `gpt-3.5-turbo-instruct`.
	Suffix param.Field[string] `json:"suffix"`
	// What sampling temperature to use, between 0 and 2. Higher values like 0.8 will
	// make the output more random, while lower values like 0.2 will make it more
	// focused and deterministic.
	//
	// We generally recommend altering this or `top_p` but not both.
	Temperature param.Field[float64] `json:"temperature"`
	// An alternative to sampling with temperature, called nucleus sampling, where the
	// model considers the results of the tokens with top_p probability mass. So 0.1
	// means only the tokens comprising the top 10% probability mass are considered.
	//
	// We generally recommend altering this or `temperature` but not both.
	TopP param.Field[float64] `json:"top_p"`
	// A unique identifier representing your end-user, which can help OpenAI to monitor
	// and detect abuse.
	// [Learn more](https://platform.openai.com/docs/guides/safety-best-practices#end-user-ids).
	User param.Field[string] `json:"user"`
}

func (r CompletionNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// ID of the model to use. You can use the
// [List models](https://platform.openai.com/docs/api-reference/models/list) API to
// see all of your available models, or see our
// [Model overview](https://platform.openai.com/docs/models) for descriptions of
// them.
type CompletionNewParamsModel string

const (
	CompletionNewParamsModelGPT3_5TurboInstruct CompletionNewParamsModel = "gpt-3.5-turbo-instruct"
	CompletionNewParamsModelDavinci002          CompletionNewParamsModel = "davinci-002"
	CompletionNewParamsModelBabbage002          CompletionNewParamsModel = "babbage-002"
)

func (r CompletionNewParamsModel) IsKnown() bool {
	switch r {
	case CompletionNewParamsModelGPT3_5TurboInstruct, CompletionNewParamsModelDavinci002, CompletionNewParamsModelBabbage002:
		return true
	}
	return false
}

// The prompt(s) to generate completions for, encoded as a string, array of
// strings, array of tokens, or array of token arrays.
//
// Note that <|endoftext|> is the document separator that the model sees during
// training, so if a prompt is not specified the model will generate as if from the
// beginning of a new document.
//
// Satisfied by [shared.UnionString], [CompletionNewParamsPromptArrayOfStrings],
// [CompletionNewParamsPromptArrayOfTokens],
// [CompletionNewParamsPromptArrayOfTokenArrays].
type CompletionNewParamsPromptUnion interface {
	ImplementsCompletionNewParamsPromptUnion()
}

type CompletionNewParamsPromptArrayOfStrings []string

func (r CompletionNewParamsPromptArrayOfStrings) ImplementsCompletionNewParamsPromptUnion() {}

type CompletionNewParamsPromptArrayOfTokens []int64

func (r CompletionNewParamsPromptArrayOfTokens) ImplementsCompletionNewParamsPromptUnion() {}

type CompletionNewParamsPromptArrayOfTokenArrays [][]int64

func (r CompletionNewParamsPromptArrayOfTokenArrays) ImplementsCompletionNewParamsPromptUnion() {}

// Up to 4 sequences where the API will stop generating further tokens. The
// returned text will not contain the stop sequence.
//
// Satisfied by [shared.UnionString], [CompletionNewParamsStopArray].
type CompletionNewParamsStopUnion interface {
	ImplementsCompletionNewParamsStopUnion()
}

type CompletionNewParamsStopArray []string

func (r CompletionNewParamsStopArray) ImplementsCompletionNewParamsStopUnion() {}

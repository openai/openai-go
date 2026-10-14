// File generated from our OpenAPI spec by Castiron. See CONTRIBUTING.md for details.

package openai

import (
	"context"
	"encoding/json"
	"net/http"
	"slices"

	"github.com/openai/openai-go/v3/internal/apijson"
	"github.com/openai/openai-go/v3/internal/requestconfig"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/packages/respjson"
	"github.com/openai/openai-go/v3/shared/constant"
)

// DecisionService contains methods and other services that help with interacting
// with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewDecisionService] method instead.
type DecisionService struct {
	Options []option.RequestOption
}

// NewDecisionService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewDecisionService(opts ...option.RequestOption) (r DecisionService) {
	r = DecisionService{}
	r.Options = requestconfig.InheritedOptions(opts...)
	return
}

// Use this endpoint to ask classification or scoring questions about the same
// input. You’ll get the answers back in the order you asked the questions.
//
// For text, you can pass a string. You can also send user messages containing
// `input_text` and `input_image` parts, with up to 128 images per request. Images
// can be base64 data URLs or publicly accessible HTTP(S) URLs. File IDs aren’t
// accepted. Other message roles, function calls, files, audio, and item references
// aren’t supported.
//
// Sometimes a question returns a refusal instead of an answer. The result has type
// `refusal` and includes the question’s name, or `null` if you didn’t give it one.
func (r *DecisionService) New(ctx context.Context, body DecisionNewParams, opts ...option.RequestOption) (res *Decision, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	path := "decisions"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

type Decision struct {
	Answers []DecisionAnswerUnion `json:"answers" api:"required"`
	Model   string                `json:"model" api:"required"`
	Usage   DecisionUsage         `json:"usage" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Answers     respjson.Field
		Model       respjson.Field
		Usage       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r Decision) RawJSON() string { return r.JSON.raw }
func (r *Decision) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// DecisionAnswerUnion contains all possible properties and values from
// [DecisionAnswerPredicate], [DecisionAnswerChoice], [DecisionAnswerScore],
// [DecisionAnswerRefusal].
//
// Use the [DecisionAnswerUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type DecisionAnswerUnion struct {
	Name string `json:"name"`
	// This field is from variant [DecisionAnswerPredicate].
	Probability float64 `json:"probability"`
	// Any of "predicate", "choice", "score", "refusal".
	Type string `json:"type"`
	// This field is from variant [DecisionAnswerChoice].
	Choice     DecisionAnswerChoiceChoiceUnion `json:"choice"`
	Confidence float64                         `json:"confidence"`
	// This field is a union of [[]DecisionAnswerChoiceProbability],
	// [[]DecisionAnswerScoreProbability]
	Probabilities DecisionAnswerUnionProbabilities `json:"probabilities"`
	// This field is from variant [DecisionAnswerScore].
	Score float64 `json:"score"`
	JSON  struct {
		Name          respjson.Field
		Probability   respjson.Field
		Type          respjson.Field
		Choice        respjson.Field
		Confidence    respjson.Field
		Probabilities respjson.Field
		Score         respjson.Field
		raw           string
	} `json:"-"`
}

// anyDecisionAnswer is implemented by each variant of [DecisionAnswerUnion] to add
// type safety for the return type of [DecisionAnswerUnion.AsAny]
type anyDecisionAnswer interface {
	implDecisionAnswerUnion()
}

func (DecisionAnswerPredicate) implDecisionAnswerUnion() {}
func (DecisionAnswerChoice) implDecisionAnswerUnion()    {}
func (DecisionAnswerScore) implDecisionAnswerUnion()     {}
func (DecisionAnswerRefusal) implDecisionAnswerUnion()   {}

// Use the following switch statement to find the correct variant
//
//	switch variant := DecisionAnswerUnion.AsAny().(type) {
//	case openai.DecisionAnswerPredicate:
//	case openai.DecisionAnswerChoice:
//	case openai.DecisionAnswerScore:
//	case openai.DecisionAnswerRefusal:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u DecisionAnswerUnion) AsAny() anyDecisionAnswer {
	switch u.Type {
	case "predicate":
		return u.AsPredicate()
	case "choice":
		return u.AsChoice()
	case "score":
		return u.AsScore()
	case "refusal":
		return u.AsRefusal()
	}
	return nil
}

func (u DecisionAnswerUnion) AsPredicate() (v DecisionAnswerPredicate) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u DecisionAnswerUnion) AsChoice() (v DecisionAnswerChoice) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u DecisionAnswerUnion) AsScore() (v DecisionAnswerScore) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u DecisionAnswerUnion) AsRefusal() (v DecisionAnswerRefusal) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u DecisionAnswerUnion) RawJSON() string { return u.JSON.raw }

func (r *DecisionAnswerUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// DecisionAnswerUnionProbabilities is an implicit subunion of
// [DecisionAnswerUnion]. DecisionAnswerUnionProbabilities provides convenient
// access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [DecisionAnswerUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfDecisionAnswerChoiceProbabilities
// OfDecisionAnswerScoreProbabilities]
type DecisionAnswerUnionProbabilities struct {
	// This field will be present if the value is a [[]DecisionAnswerChoiceProbability]
	// instead of an object.
	OfDecisionAnswerChoiceProbabilities []DecisionAnswerChoiceProbability `json:",inline"`
	// This field will be present if the value is a [[]DecisionAnswerScoreProbability]
	// instead of an object.
	OfDecisionAnswerScoreProbabilities []DecisionAnswerScoreProbability `json:",inline"`
	JSON                               struct {
		OfDecisionAnswerChoiceProbabilities respjson.Field
		OfDecisionAnswerScoreProbabilities  respjson.Field
		raw                                 string
	} `json:"-"`
}

func (r *DecisionAnswerUnionProbabilities) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type DecisionAnswerPredicate struct {
	Name        string  `json:"name" api:"required"`
	Probability float64 `json:"probability" api:"required"`
	// The type of the object. Always `predicate`.
	Type constant.Predicate `json:"type" default:"predicate"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Name        respjson.Field
		Probability respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DecisionAnswerPredicate) RawJSON() string { return r.JSON.raw }
func (r *DecisionAnswerPredicate) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type DecisionAnswerChoice struct {
	// Choice values are typed: a string and a boolean with the same text are distinct.
	Choice        DecisionAnswerChoiceChoiceUnion   `json:"choice" api:"required"`
	Confidence    float64                           `json:"confidence" api:"required"`
	Name          string                            `json:"name" api:"required"`
	Probabilities []DecisionAnswerChoiceProbability `json:"probabilities" api:"required"`
	// The type of the object. Always `choice`.
	Type constant.Choice `json:"type" default:"choice"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Choice        respjson.Field
		Confidence    respjson.Field
		Name          respjson.Field
		Probabilities respjson.Field
		Type          respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DecisionAnswerChoice) RawJSON() string { return r.JSON.raw }
func (r *DecisionAnswerChoice) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// DecisionAnswerChoiceChoiceUnion contains all possible properties and values from
// [string], [bool].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfString OfBool]
type DecisionAnswerChoiceChoiceUnion struct {
	// This field will be present if the value is a [string] instead of an object.
	OfString string `json:",inline"`
	// This field will be present if the value is a [bool] instead of an object.
	OfBool bool `json:",inline"`
	JSON   struct {
		OfString respjson.Field
		OfBool   respjson.Field
		raw      string
	} `json:"-"`
}

func (u DecisionAnswerChoiceChoiceUnion) AsString() (v string) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u DecisionAnswerChoiceChoiceUnion) AsBool() (v bool) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u DecisionAnswerChoiceChoiceUnion) RawJSON() string { return u.JSON.raw }

func (r *DecisionAnswerChoiceChoiceUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type DecisionAnswerChoiceProbability struct {
	Probability float64 `json:"probability" api:"required"`
	// Choice values are typed: a string and a boolean with the same text are distinct.
	Value DecisionAnswerChoiceProbabilityValueUnion `json:"value" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Probability respjson.Field
		Value       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DecisionAnswerChoiceProbability) RawJSON() string { return r.JSON.raw }
func (r *DecisionAnswerChoiceProbability) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// DecisionAnswerChoiceProbabilityValueUnion contains all possible properties and
// values from [string], [bool].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfString OfBool]
type DecisionAnswerChoiceProbabilityValueUnion struct {
	// This field will be present if the value is a [string] instead of an object.
	OfString string `json:",inline"`
	// This field will be present if the value is a [bool] instead of an object.
	OfBool bool `json:",inline"`
	JSON   struct {
		OfString respjson.Field
		OfBool   respjson.Field
		raw      string
	} `json:"-"`
}

func (u DecisionAnswerChoiceProbabilityValueUnion) AsString() (v string) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u DecisionAnswerChoiceProbabilityValueUnion) AsBool() (v bool) {
	_ = apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u DecisionAnswerChoiceProbabilityValueUnion) RawJSON() string { return u.JSON.raw }

func (r *DecisionAnswerChoiceProbabilityValueUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type DecisionAnswerScore struct {
	Confidence    float64                          `json:"confidence" api:"required"`
	Name          string                           `json:"name" api:"required"`
	Probabilities []DecisionAnswerScoreProbability `json:"probabilities" api:"required"`
	Score         float64                          `json:"score" api:"required"`
	// The type of the object. Always `score`.
	Type constant.Score `json:"type" default:"score"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Confidence    respjson.Field
		Name          respjson.Field
		Probabilities respjson.Field
		Score         respjson.Field
		Type          respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DecisionAnswerScore) RawJSON() string { return r.JSON.raw }
func (r *DecisionAnswerScore) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type DecisionAnswerScoreProbability struct {
	Label       string  `json:"label" api:"required"`
	Probability float64 `json:"probability" api:"required"`
	Value       int64   `json:"value" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Label       respjson.Field
		Probability respjson.Field
		Value       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DecisionAnswerScoreProbability) RawJSON() string { return r.JSON.raw }
func (r *DecisionAnswerScoreProbability) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The model declined to answer this question. Other questions in the same request
// can still receive answers.
type DecisionAnswerRefusal struct {
	Name string `json:"name" api:"required"`
	// The type of the object. Always `refusal`.
	Type constant.Refusal `json:"type" default:"refusal"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Name        respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DecisionAnswerRefusal) RawJSON() string { return r.JSON.raw }
func (r *DecisionAnswerRefusal) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type DecisionUsage struct {
	InputTokens         int64                            `json:"input_tokens" api:"required"`
	InputTokensDetails  DecisionUsageInputTokensDetails  `json:"input_tokens_details" api:"required"`
	OutputTokens        int64                            `json:"output_tokens" api:"required"`
	OutputTokensDetails DecisionUsageOutputTokensDetails `json:"output_tokens_details" api:"required"`
	TotalTokens         int64                            `json:"total_tokens" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		InputTokens         respjson.Field
		InputTokensDetails  respjson.Field
		OutputTokens        respjson.Field
		OutputTokensDetails respjson.Field
		TotalTokens         respjson.Field
		ExtraFields         map[string]respjson.Field
		raw                 string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DecisionUsage) RawJSON() string { return r.JSON.raw }
func (r *DecisionUsage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type DecisionUsageInputTokensDetails struct {
	CacheWriteTokens int64 `json:"cache_write_tokens" api:"required"`
	CachedTokens     int64 `json:"cached_tokens" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		CacheWriteTokens respjson.Field
		CachedTokens     respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DecisionUsageInputTokensDetails) RawJSON() string { return r.JSON.raw }
func (r *DecisionUsageInputTokensDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type DecisionUsageOutputTokensDetails struct {
	ReasoningTokens int64 `json:"reasoning_tokens" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ReasoningTokens respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DecisionUsageOutputTokensDetails) RawJSON() string { return r.JSON.raw }
func (r *DecisionUsageOutputTokensDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// An image provided as a base64 data URL or a publicly accessible HTTP(S) URL.
// File IDs are not supported.
//
// The properties ImageURL, Type are required.
type DecisionInputImageParam struct {
	// A base64-encoded image in a data URL or a publicly accessible HTTP(S) image URL.
	ImageURL string `json:"image_url" api:"required"`
	// The image detail level, using the selected model's image profile. Defaults to
	// auto.
	//
	// Any of "low", "high", "auto", "original".
	Detail DecisionInputImageDetail `json:"detail,omitzero"`
	// This field can be elided, and will marshal its zero value as "input_image".
	Type constant.InputImage `json:"type" default:"input_image"`
	paramObj
}

func (r DecisionInputImageParam) MarshalJSON() (data []byte, err error) {
	type shadow DecisionInputImageParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionInputImageParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The image detail level, using the selected model's image profile. Defaults to
// auto.
type DecisionInputImageDetail string

const (
	DecisionInputImageDetailLow      DecisionInputImageDetail = "low"
	DecisionInputImageDetailHigh     DecisionInputImageDetail = "high"
	DecisionInputImageDetailAuto     DecisionInputImageDetail = "auto"
	DecisionInputImageDetailOriginal DecisionInputImageDetail = "original"
)

// A user message containing text or images.
//
// The properties Content, Role are required.
type DecisionInputMessageParam struct {
	// Text evidence or an ordered list of text and image parts.
	Content DecisionInputMessageContentUnionParam `json:"content,omitzero" api:"required"`
	// Any of "message".
	Type DecisionInputMessageType `json:"type,omitzero"`
	// This field can be elided, and will marshal its zero value as "user".
	Role constant.User `json:"role" default:"user"`
	paramObj
}

func (r DecisionInputMessageParam) MarshalJSON() (data []byte, err error) {
	type shadow DecisionInputMessageParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionInputMessageParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DecisionInputMessageContentUnionParam struct {
	OfString param.Opt[string]             `json:",omitzero,inline"`
	OfParts  []DecisionInputPartUnionParam `json:",omitzero,inline"`
	paramUnion
}

func (u DecisionInputMessageContentUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfString, u.OfParts)
}
func (u *DecisionInputMessageContentUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

type DecisionInputMessageType string

const (
	DecisionInputMessageTypeMessage DecisionInputMessageType = "message"
)

func DecisionInputPartParamOfInputText(text string) DecisionInputPartUnionParam {
	var inputText DecisionInputTextParam
	inputText.Text = text
	return DecisionInputPartUnionParam{OfInputText: &inputText}
}

func DecisionInputPartParamOfInputImage(imageURL string) DecisionInputPartUnionParam {
	var inputImage DecisionInputImageParam
	inputImage.ImageURL = imageURL
	return DecisionInputPartUnionParam{OfInputImage: &inputImage}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DecisionInputPartUnionParam struct {
	OfInputText  *DecisionInputTextParam  `json:",omitzero,inline"`
	OfInputImage *DecisionInputImageParam `json:",omitzero,inline"`
	paramUnion
}

func (u DecisionInputPartUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfInputText, u.OfInputImage)
}
func (u *DecisionInputPartUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionInputPartUnionParam) GetText() *string {
	if vt := u.OfInputText; vt != nil {
		return &vt.Text
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionInputPartUnionParam) GetImageURL() *string {
	if vt := u.OfInputImage; vt != nil {
		return &vt.ImageURL
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionInputPartUnionParam) GetDetail() *string {
	if vt := u.OfInputImage; vt != nil {
		return (*string)(&vt.Detail)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionInputPartUnionParam) GetType() *string {
	if vt := u.OfInputText; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfInputImage; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[DecisionInputPartUnionParam](
		"type",
		apijson.Discriminator[DecisionInputTextParam]("input_text"),
		apijson.Discriminator[DecisionInputImageParam]("input_image"),
	)
}

// The properties Text, Type are required.
type DecisionInputTextParam struct {
	Text string `json:"text" api:"required"`
	// This field can be elided, and will marshal its zero value as "input_text".
	Type constant.InputText `json:"type" default:"input_text"`
	paramObj
}

func (r DecisionInputTextParam) MarshalJSON() (data []byte, err error) {
	type shadow DecisionInputTextParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionInputTextParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type DecisionNewParams struct {
	// The text or images to evaluate for every question. Provide a text string or user
	// messages containing text and images. Images can be base64 data URLs or publicly
	// accessible HTTP(S) URLs; at most 128 images are allowed across all messages in
	// one request. Files, audio, tools, and item references are not supported.
	Input     DecisionNewParamsInputUnion      `json:"input,omitzero" api:"required"`
	Model     string                           `json:"model" api:"required"`
	Questions []DecisionNewParamsQuestionUnion `json:"questions,omitzero" api:"required"`
	// Opaque caller-provided end-user identifier, scoped by the verified org. Match
	// Responses' limit; this is never the authenticated user identity.
	SafetyIdentifier param.Opt[string] `json:"safety_identifier,omitzero"`
	paramObj
}

func (r DecisionNewParams) MarshalJSON() (data []byte, err error) {
	type shadow DecisionNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DecisionNewParamsInputUnion struct {
	OfString                    param.Opt[string]           `json:",omitzero,inline"`
	OfDecisionInputMessageArray []DecisionInputMessageParam `json:",omitzero,inline"`
	paramUnion
}

func (u DecisionNewParamsInputUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfString, u.OfDecisionInputMessageArray)
}
func (u *DecisionNewParamsInputUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DecisionNewParamsQuestionUnion struct {
	OfPredicate *DecisionNewParamsQuestionPredicate `json:",omitzero,inline"`
	OfChoice    *DecisionNewParamsQuestionChoice    `json:",omitzero,inline"`
	OfScore     *DecisionNewParamsQuestionScore     `json:",omitzero,inline"`
	paramUnion
}

func (u DecisionNewParamsQuestionUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfPredicate, u.OfChoice, u.OfScore)
}
func (u *DecisionNewParamsQuestionUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionNewParamsQuestionUnion) GetChoices() []DecisionNewParamsQuestionChoiceChoice {
	if vt := u.OfChoice; vt != nil {
		return vt.Choices
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionNewParamsQuestionUnion) GetLevels() []DecisionNewParamsQuestionScoreLevel {
	if vt := u.OfScore; vt != nil {
		return vt.Levels
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionNewParamsQuestionUnion) GetInstructions() *string {
	if vt := u.OfPredicate; vt != nil {
		return (*string)(&vt.Instructions)
	} else if vt := u.OfChoice; vt != nil {
		return (*string)(&vt.Instructions)
	} else if vt := u.OfScore; vt != nil {
		return (*string)(&vt.Instructions)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionNewParamsQuestionUnion) GetType() *string {
	if vt := u.OfPredicate; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfChoice; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfScore; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionNewParamsQuestionUnion) GetName() *string {
	if vt := u.OfPredicate; vt != nil && vt.Name.Valid() {
		return &vt.Name.Value
	} else if vt := u.OfChoice; vt != nil && vt.Name.Valid() {
		return &vt.Name.Value
	} else if vt := u.OfScore; vt != nil && vt.Name.Valid() {
		return &vt.Name.Value
	}
	return nil
}

func init() {
	apijson.RegisterUnion[DecisionNewParamsQuestionUnion](
		"type",
		apijson.Discriminator[DecisionNewParamsQuestionPredicate]("predicate"),
		apijson.Discriminator[DecisionNewParamsQuestionChoice]("choice"),
		apijson.Discriminator[DecisionNewParamsQuestionScore]("score"),
	)
}

// Estimate how likely it is that a statement about the input is true.
//
// The properties Instructions, Type are required.
type DecisionNewParamsQuestionPredicate struct {
	Instructions string            `json:"instructions" api:"required"`
	Name         param.Opt[string] `json:"name,omitzero"`
	// The type of the object. Always `predicate`.
	//
	// This field can be elided, and will marshal its zero value as "predicate".
	Type constant.Predicate `json:"type" default:"predicate"`
	paramObj
}

func (r DecisionNewParamsQuestionPredicate) MarshalJSON() (data []byte, err error) {
	type shadow DecisionNewParamsQuestionPredicate
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionNewParamsQuestionPredicate) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Choose from the supplied options based on the input.
//
// The properties Choices, Instructions, Type are required.
type DecisionNewParamsQuestionChoice struct {
	// Provide between 2 and 255 choices. Each choice must be unique.
	Choices      []DecisionNewParamsQuestionChoiceChoice `json:"choices,omitzero" api:"required"`
	Instructions string                                  `json:"instructions" api:"required"`
	Name         param.Opt[string]                       `json:"name,omitzero"`
	// The type of the object. Always `choice`.
	//
	// This field can be elided, and will marshal its zero value as "choice".
	Type constant.Choice `json:"type" default:"choice"`
	paramObj
}

func (r DecisionNewParamsQuestionChoice) MarshalJSON() (data []byte, err error) {
	type shadow DecisionNewParamsQuestionChoice
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionNewParamsQuestionChoice) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The property Value is required.
type DecisionNewParamsQuestionChoiceChoice struct {
	// Choice values are typed: a string and a boolean with the same text are distinct.
	Value       DecisionNewParamsQuestionChoiceChoiceValueUnion `json:"value,omitzero" api:"required"`
	Description param.Opt[string]                               `json:"description,omitzero"`
	paramObj
}

func (r DecisionNewParamsQuestionChoiceChoice) MarshalJSON() (data []byte, err error) {
	type shadow DecisionNewParamsQuestionChoiceChoice
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionNewParamsQuestionChoiceChoice) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DecisionNewParamsQuestionChoiceChoiceValueUnion struct {
	OfString param.Opt[string] `json:",omitzero,inline"`
	OfBool   param.Opt[bool]   `json:",omitzero,inline"`
	paramUnion
}

func (u DecisionNewParamsQuestionChoiceChoiceValueUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfString, u.OfBool)
}
func (u *DecisionNewParamsQuestionChoiceChoiceValueUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Rate the input against the supplied ordered levels.
//
// The properties Instructions, Levels, Type are required.
type DecisionNewParamsQuestionScore struct {
	Instructions string                                `json:"instructions" api:"required"`
	Levels       []DecisionNewParamsQuestionScoreLevel `json:"levels,omitzero" api:"required"`
	Name         param.Opt[string]                     `json:"name,omitzero"`
	// The type of the object. Always `score`.
	//
	// This field can be elided, and will marshal its zero value as "score".
	Type constant.Score `json:"type" default:"score"`
	paramObj
}

func (r DecisionNewParamsQuestionScore) MarshalJSON() (data []byte, err error) {
	type shadow DecisionNewParamsQuestionScore
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionNewParamsQuestionScore) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The property Label is required.
type DecisionNewParamsQuestionScoreLevel struct {
	Label       string            `json:"label" api:"required"`
	Description param.Opt[string] `json:"description,omitzero"`
	paramObj
}

func (r DecisionNewParamsQuestionScoreLevel) MarshalJSON() (data []byte, err error) {
	type shadow DecisionNewParamsQuestionScoreLevel
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionNewParamsQuestionScoreLevel) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

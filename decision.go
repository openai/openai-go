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
	"github.com/openai/openai-go/v3/shared"
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

// Use this endpoint to quickly answer classification or scoring questions about
// some input. You can ask multiple questions in one request.
func (r *DecisionService) New(ctx context.Context, body DecisionNewParams, opts ...option.RequestOption) (res *Decision, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithBearerAuthSecurity()}
	opts = slices.Concat(preClientOpts, r.Options, opts)
	path := "decisions"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

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

type Decision struct {
	// Answers in question order. Any question can return a refusal instead of an
	// answer; other questions can still receive answers. A refusal has
	// `type: "refusal"` and the question's `name`, or `null` if it was unnamed.
	Answers []DecisionAnswerUnion `json:"answers" api:"required"`
	// The model used.
	Model string `json:"model" api:"required"`
	// Token usage for the request.
	Usage DecisionUsage `json:"usage" api:"required"`
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

// Token usage for the request.
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
	// Marks the exact end of a reusable prompt prefix. The breakpoint inherits its TTL
	// from the request's `prompt_cache_options.ttl`; the boundary is not rounded to a
	// token block.
	PromptCacheBreakpoint DecisionPromptCacheBreakpointParam `json:"prompt_cache_breakpoint,omitzero"`
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

func DecisionInputMessageParamOfUser[T string | []DecisionInputPartUnionParam](content T) DecisionInputMessageUnionParam {
	var user DecisionInputMessageUserParam
	switch v := any(content).(type) {
	case string:
		user.Content.OfString = param.NewOpt(v)
	case []DecisionInputPartUnionParam:
		user.Content.OfParts = v
	}
	return DecisionInputMessageUnionParam{OfUser: &user}
}

func DecisionInputMessageParamOfDeveloper[T string | []DecisionInputTextParam](content T) DecisionInputMessageUnionParam {
	var developer DecisionInputMessageDeveloperParam
	switch v := any(content).(type) {
	case string:
		developer.Content.OfString = param.NewOpt(v)
	case []DecisionInputTextParam:
		developer.Content.OfParts = v
	}
	return DecisionInputMessageUnionParam{OfDeveloper: &developer}
}

func DecisionInputMessageParamOfSystem[T string | []DecisionInputTextParam](content T) DecisionInputMessageUnionParam {
	var system DecisionInputMessageSystemParam
	switch v := any(content).(type) {
	case string:
		system.Content.OfString = param.NewOpt(v)
	case []DecisionInputTextParam:
		system.Content.OfParts = v
	}
	return DecisionInputMessageUnionParam{OfSystem: &system}
}

func DecisionInputMessageParamOfAssistant[
	T string | []DecisionInputMessageAssistantContentPartsItemUnionParam,
](content T) DecisionInputMessageUnionParam {
	var assistant DecisionInputMessageAssistantParam
	switch v := any(content).(type) {
	case string:
		assistant.Content.OfString = param.NewOpt(v)
	case []DecisionInputMessageAssistantContentPartsItemUnionParam:
		assistant.Content.OfParts = v
	}
	return DecisionInputMessageUnionParam{OfAssistant: &assistant}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DecisionInputMessageUnionParam struct {
	OfUser      *DecisionInputMessageUserParam      `json:",omitzero,inline"`
	OfDeveloper *DecisionInputMessageDeveloperParam `json:",omitzero,inline"`
	OfSystem    *DecisionInputMessageSystemParam    `json:",omitzero,inline"`
	OfAssistant *DecisionInputMessageAssistantParam `json:",omitzero,inline"`
	paramUnion
}

func (u DecisionInputMessageUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfUser, u.OfDeveloper, u.OfSystem, u.OfAssistant)
}
func (u *DecisionInputMessageUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionInputMessageUnionParam) GetPhase() *string {
	if vt := u.OfAssistant; vt != nil {
		return &vt.Phase
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionInputMessageUnionParam) GetRole() *string {
	if vt := u.OfUser; vt != nil {
		return (*string)(&vt.Role)
	} else if vt := u.OfDeveloper; vt != nil {
		return (*string)(&vt.Role)
	} else if vt := u.OfSystem; vt != nil {
		return (*string)(&vt.Role)
	} else if vt := u.OfAssistant; vt != nil {
		return (*string)(&vt.Role)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionInputMessageUnionParam) GetType() *string {
	if vt := u.OfUser; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfDeveloper; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfSystem; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfAssistant; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u DecisionInputMessageUnionParam) GetContent() (res decisionInputMessageUnionParamContent) {
	if vt := u.OfUser; vt != nil {
		res.any = vt.Content.asAny()
	} else if vt := u.OfDeveloper; vt != nil {
		res.any = vt.Content.asAny()
	} else if vt := u.OfSystem; vt != nil {
		res.any = vt.Content.asAny()
	} else if vt := u.OfAssistant; vt != nil {
		res.any = vt.Content.asAny()
	}
	return
}

// Can have the runtime types [*string], [_[]DecisionInputPartUnionParam],
// [_[]DecisionInputTextParam],
// [\*[]DecisionInputMessageAssistantContentPartsItemUnionParam]
type decisionInputMessageUnionParamContent struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *string:
//	case *[]openai.DecisionInputPartUnionParam:
//	case *[]openai.DecisionInputTextParam:
//	case *[]openai.DecisionInputMessageAssistantContentPartsItemUnionParam:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u decisionInputMessageUnionParamContent) AsAny() any { return u.any }

func init() {
	apijson.RegisterUnion[DecisionInputMessageUnionParam](
		"role",
		apijson.Discriminator[DecisionInputMessageUserParam]("user"),
		apijson.Discriminator[DecisionInputMessageDeveloperParam]("developer"),
		apijson.Discriminator[DecisionInputMessageSystemParam]("system"),
		apijson.Discriminator[DecisionInputMessageAssistantParam]("assistant"),
	)
}

// A user message containing input text, images, or JSON.
//
// The properties Content, Role are required.
type DecisionInputMessageUserParam struct {
	// A text string or an ordered list of text, image, and JSON parts.
	Content DecisionInputMessageUserContentUnionParam `json:"content,omitzero" api:"required"`
	// Any of "message".
	Type string `json:"type,omitzero"`
	// This field can be elided, and will marshal its zero value as "user".
	Role constant.User `json:"role" default:"user"`
	paramObj
}

func (r DecisionInputMessageUserParam) MarshalJSON() (data []byte, err error) {
	type shadow DecisionInputMessageUserParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionInputMessageUserParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[DecisionInputMessageUserParam](
		"type", "message",
	)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DecisionInputMessageUserContentUnionParam struct {
	OfString param.Opt[string]             `json:",omitzero,inline"`
	OfParts  []DecisionInputPartUnionParam `json:",omitzero,inline"`
	paramUnion
}

func (u DecisionInputMessageUserContentUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfString, u.OfParts)
}
func (u *DecisionInputMessageUserContentUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *DecisionInputMessageUserContentUnionParam) asAny() any {
	if !param.IsOmitted(u.OfString) {
		return &u.OfString.Value
	} else if !param.IsOmitted(u.OfParts) {
		return &u.OfParts
	}
	return nil
}

// A developer message containing input text.
//
// The properties Content, Role are required.
type DecisionInputMessageDeveloperParam struct {
	// A text string or an ordered list of input text parts.
	Content DecisionInputMessageDeveloperContentUnionParam `json:"content,omitzero" api:"required"`
	// Any of "message".
	Type string `json:"type,omitzero"`
	// This field can be elided, and will marshal its zero value as "developer".
	Role constant.Developer `json:"role" default:"developer"`
	paramObj
}

func (r DecisionInputMessageDeveloperParam) MarshalJSON() (data []byte, err error) {
	type shadow DecisionInputMessageDeveloperParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionInputMessageDeveloperParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[DecisionInputMessageDeveloperParam](
		"type", "message",
	)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DecisionInputMessageDeveloperContentUnionParam struct {
	OfString param.Opt[string]        `json:",omitzero,inline"`
	OfParts  []DecisionInputTextParam `json:",omitzero,inline"`
	paramUnion
}

func (u DecisionInputMessageDeveloperContentUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfString, u.OfParts)
}
func (u *DecisionInputMessageDeveloperContentUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *DecisionInputMessageDeveloperContentUnionParam) asAny() any {
	if !param.IsOmitted(u.OfString) {
		return &u.OfString.Value
	} else if !param.IsOmitted(u.OfParts) {
		return &u.OfParts
	}
	return nil
}

// A system message containing input text.
//
// The properties Content, Role are required.
type DecisionInputMessageSystemParam struct {
	// A text string or an ordered list of input text parts.
	Content DecisionInputMessageSystemContentUnionParam `json:"content,omitzero" api:"required"`
	// Any of "message".
	Type string `json:"type,omitzero"`
	// This field can be elided, and will marshal its zero value as "system".
	Role constant.System `json:"role" default:"system"`
	paramObj
}

func (r DecisionInputMessageSystemParam) MarshalJSON() (data []byte, err error) {
	type shadow DecisionInputMessageSystemParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionInputMessageSystemParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[DecisionInputMessageSystemParam](
		"type", "message",
	)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DecisionInputMessageSystemContentUnionParam struct {
	OfString param.Opt[string]        `json:",omitzero,inline"`
	OfParts  []DecisionInputTextParam `json:",omitzero,inline"`
	paramUnion
}

func (u DecisionInputMessageSystemContentUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfString, u.OfParts)
}
func (u *DecisionInputMessageSystemContentUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *DecisionInputMessageSystemContentUnionParam) asAny() any {
	if !param.IsOmitted(u.OfString) {
		return &u.OfString.Value
	} else if !param.IsOmitted(u.OfParts) {
		return &u.OfParts
	}
	return nil
}

// An assistant message containing output text or refusals.
//
// The properties Content, Role are required.
type DecisionInputMessageAssistantParam struct {
	// A text string or an ordered list of output text and refusal parts.
	Content DecisionInputMessageAssistantContentUnionParam `json:"content,omitzero" api:"required"`
	// The phase of an assistant message.
	//
	// Any of "commentary", "final_answer".
	Phase string `json:"phase,omitzero"`
	// Any of "message".
	Type string `json:"type,omitzero"`
	// This field can be elided, and will marshal its zero value as "assistant".
	Role constant.Assistant `json:"role" default:"assistant"`
	paramObj
}

func (r DecisionInputMessageAssistantParam) MarshalJSON() (data []byte, err error) {
	type shadow DecisionInputMessageAssistantParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionInputMessageAssistantParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[DecisionInputMessageAssistantParam](
		"phase", "commentary", "final_answer",
	)
	apijson.RegisterFieldValidator[DecisionInputMessageAssistantParam](
		"type", "message",
	)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DecisionInputMessageAssistantContentUnionParam struct {
	OfString param.Opt[string]                                         `json:",omitzero,inline"`
	OfParts  []DecisionInputMessageAssistantContentPartsItemUnionParam `json:",omitzero,inline"`
	paramUnion
}

func (u DecisionInputMessageAssistantContentUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfString, u.OfParts)
}
func (u *DecisionInputMessageAssistantContentUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *DecisionInputMessageAssistantContentUnionParam) asAny() any {
	if !param.IsOmitted(u.OfString) {
		return &u.OfString.Value
	} else if !param.IsOmitted(u.OfParts) {
		return &u.OfParts
	}
	return nil
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DecisionInputMessageAssistantContentPartsItemUnionParam struct {
	OfOutputText *DecisionInputMessageAssistantContentPartOutputTextParam `json:",omitzero,inline"`
	OfRefusal    *DecisionInputMessageAssistantContentPartRefusalParam    `json:",omitzero,inline"`
	paramUnion
}

func (u DecisionInputMessageAssistantContentPartsItemUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfOutputText, u.OfRefusal)
}
func (u *DecisionInputMessageAssistantContentPartsItemUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionInputMessageAssistantContentPartsItemUnionParam) GetText() *string {
	if vt := u.OfOutputText; vt != nil {
		return &vt.Text
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionInputMessageAssistantContentPartsItemUnionParam) GetRefusal() *string {
	if vt := u.OfRefusal; vt != nil {
		return &vt.Refusal
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionInputMessageAssistantContentPartsItemUnionParam) GetType() *string {
	if vt := u.OfOutputText; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfRefusal; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[DecisionInputMessageAssistantContentPartsItemUnionParam](
		"type",
		apijson.Discriminator[DecisionInputMessageAssistantContentPartOutputTextParam]("output_text"),
		apijson.Discriminator[DecisionInputMessageAssistantContentPartRefusalParam]("refusal"),
	)
}

// The properties Text, Type are required.
type DecisionInputMessageAssistantContentPartOutputTextParam struct {
	Text string `json:"text" api:"required"`
	// This field can be elided, and will marshal its zero value as "output_text".
	Type constant.OutputText `json:"type" default:"output_text"`
	paramObj
}

func (r DecisionInputMessageAssistantContentPartOutputTextParam) MarshalJSON() (data []byte, err error) {
	type shadow DecisionInputMessageAssistantContentPartOutputTextParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionInputMessageAssistantContentPartOutputTextParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Refusal, Type are required.
type DecisionInputMessageAssistantContentPartRefusalParam struct {
	Refusal string `json:"refusal" api:"required"`
	// This field can be elided, and will marshal its zero value as "refusal".
	Type constant.Refusal `json:"type" default:"refusal"`
	paramObj
}

func (r DecisionInputMessageAssistantContentPartRefusalParam) MarshalJSON() (data []byte, err error) {
	type shadow DecisionInputMessageAssistantContentPartRefusalParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionInputMessageAssistantContentPartRefusalParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

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

func DecisionInputPartParamOfInputJSON[T map[string]any | []map[string]any](data T) DecisionInputPartUnionParam {
	var inputJSON shared.DecisionInputJSONParam
	switch v := any(data).(type) {
	case map[string]any:
		inputJSON.Data.OfAnyMap = v
	case []map[string]any:
		inputJSON.Data.OfArrayOfValues = v
	}
	return DecisionInputPartUnionParam{OfInputJSON: &inputJSON}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DecisionInputPartUnionParam struct {
	OfInputText  *DecisionInputTextParam        `json:",omitzero,inline"`
	OfInputImage *DecisionInputImageParam       `json:",omitzero,inline"`
	OfInputJSON  *shared.DecisionInputJSONParam `json:",omitzero,inline"`
	paramUnion
}

func (u DecisionInputPartUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfInputText, u.OfInputImage, u.OfInputJSON)
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
func (u DecisionInputPartUnionParam) GetData() *shared.DecisionInputJSONDataUnionParam {
	if vt := u.OfInputJSON; vt != nil {
		return &vt.Data
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionInputPartUnionParam) GetType() *string {
	if vt := u.OfInputText; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfInputImage; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfInputJSON; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's PromptCacheBreakpoint property, if
// present.
func (u DecisionInputPartUnionParam) GetPromptCacheBreakpoint() *DecisionPromptCacheBreakpointParam {
	if vt := u.OfInputText; vt != nil {
		return &vt.PromptCacheBreakpoint
	} else if vt := u.OfInputImage; vt != nil {
		return &vt.PromptCacheBreakpoint
	}
	return nil
}

func init() {
	apijson.RegisterUnion[DecisionInputPartUnionParam](
		"type",
		apijson.Discriminator[DecisionInputTextParam]("input_text"),
		apijson.Discriminator[DecisionInputImageParam]("input_image"),
		apijson.Discriminator[shared.DecisionInputJSONParam]("input_json"),
	)
}

// The properties Text, Type are required.
type DecisionInputTextParam struct {
	Text string `json:"text" api:"required"`
	// Marks the exact end of a reusable prompt prefix. The breakpoint inherits its TTL
	// from the request's `prompt_cache_options.ttl`; the boundary is not rounded to a
	// token block.
	PromptCacheBreakpoint DecisionPromptCacheBreakpointParam `json:"prompt_cache_breakpoint,omitzero"`
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

func NewDecisionPromptCacheBreakpointParam() DecisionPromptCacheBreakpointParam {
	return DecisionPromptCacheBreakpointParam{
		Mode: "explicit",
	}
}

// Marks the exact end of a reusable prompt prefix. The breakpoint inherits its TTL
// from the request's `prompt_cache_options.ttl`; the boundary is not rounded to a
// token block.
//
// This struct has a constant value, construct it with
// [NewDecisionPromptCacheBreakpointParam].
type DecisionPromptCacheBreakpointParam struct {
	// The breakpoint mode. Always `explicit`.
	Mode constant.Explicit `json:"mode" default:"explicit"`
	paramObj
}

func (r DecisionPromptCacheBreakpointParam) MarshalJSON() (data []byte, err error) {
	type shadow DecisionPromptCacheBreakpointParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionPromptCacheBreakpointParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type DecisionNewParams struct {
	// The input to evaluate for every question. Provide a text string or an ordered
	// array of messages, function calls and outputs, or custom tool calls and outputs.
	// Each tool output must identify a matching call by call ID.
	Input DecisionNewParamsInputUnion `json:"input,omitzero" api:"required"`
	// The model to use. Only `gpt-6-luna` is supported.
	Model string `json:"model" api:"required"`
	// Ask 1–400 questions about the input. Each question can be a predicate, choice,
	// or score. Answers come back in the same order.
	Questions []DecisionNewParamsQuestionUnion `json:"questions,omitzero" api:"required"`
	// Partition prompt cache reuse by caller-provided key.
	PromptCacheKey param.Opt[string] `json:"prompt_cache_key,omitzero"`
	// Opaque caller-provided end-user identifier, scoped by the verified org. Match
	// Responses' limit; this is never the authenticated user identity.
	SafetyIdentifier param.Opt[string] `json:"safety_identifier,omitzero"`
	// Options for prompt caching. Supported for `gpt-5.6` and later models. By
	// default, OpenAI automatically chooses one implicit cache breakpoint. You can add
	// explicit breakpoints to content blocks with `prompt_cache_breakpoint`. Each
	// request can write up to four breakpoints. For cache matching, OpenAI considers
	// up to the latest 80 breakpoints in the conversation, without a content-block
	// lookback limit. Set `mode` to `explicit` to disable the implicit breakpoint. The
	// `ttl` defaults to `30m`, which is currently the only supported value. See the
	// [prompt caching guide](https://developers.openai.com/api/docs/guides/prompt-caching)
	// for current details.
	PromptCacheOptions DecisionNewParamsPromptCacheOptions `json:"prompt_cache_options,omitzero"`
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
	OfString                    param.Opt[string]                      `json:",omitzero,inline"`
	OfDecisionInputMessageArray []DecisionInputMessageParam            `json:",omitzero,inline"`
	OfDecisionNewsInputArray    []DecisionNewParamsInputArrayItemUnion `json:",omitzero,inline"`
	paramUnion
}

func (u DecisionNewParamsInputUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfString, u.OfDecisionInputMessageArray, u.OfDecisionNewsInputArray)
}
func (u *DecisionNewParamsInputUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DecisionNewParamsInputArrayItemUnion struct {
	OfDecisionInputUserMessage      *DecisionNewParamsInputArrayItemDecisionInputUserMessage      `json:",omitzero,inline"`
	OfDecisionInputDeveloperMessage *DecisionNewParamsInputArrayItemDecisionInputDeveloperMessage `json:",omitzero,inline"`
	OfDecisionInputSystemMessage    *DecisionNewParamsInputArrayItemDecisionInputSystemMessage    `json:",omitzero,inline"`
	OfDecisionInputAssistantMessage *DecisionNewParamsInputArrayItemDecisionInputAssistantMessage `json:",omitzero,inline"`
	OfDecisionInputMessage          *DecisionInputMessageUnionParam                               `json:",omitzero,inline"`
	OfFunctionCall                  *DecisionNewParamsInputArrayItemFunctionCall                  `json:",omitzero,inline"`
	OfFunctionCallOutput            *DecisionNewParamsInputArrayItemFunctionCallOutput            `json:",omitzero,inline"`
	OfCustomToolCall                *DecisionNewParamsInputArrayItemCustomToolCall                `json:",omitzero,inline"`
	OfCustomToolCallOutput          *DecisionNewParamsInputArrayItemCustomToolCallOutput          `json:",omitzero,inline"`
	paramUnion                      `default:"message"`
}

func (u DecisionNewParamsInputArrayItemUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfDecisionInputUserMessage,
		u.OfDecisionInputDeveloperMessage,
		u.OfDecisionInputSystemMessage,
		u.OfDecisionInputAssistantMessage,
		u.OfDecisionInputMessage,
		u.OfFunctionCall,
		u.OfFunctionCallOutput,
		u.OfCustomToolCall,
		u.OfCustomToolCallOutput)
}
func (u *DecisionNewParamsInputArrayItemUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionNewParamsInputArrayItemUnion) GetPhase() *string {
	if vt := u.OfDecisionInputAssistantMessage; vt != nil {
		return &vt.Phase
	} else if vt := u.OfDecisionInputMessage; vt != nil {
		return vt.GetPhase()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionNewParamsInputArrayItemUnion) GetArguments() *string {
	if vt := u.OfFunctionCall; vt != nil {
		return &vt.Arguments
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionNewParamsInputArrayItemUnion) GetInput() *string {
	if vt := u.OfCustomToolCall; vt != nil {
		return &vt.Input
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionNewParamsInputArrayItemUnion) GetRole() *string {
	if vt := u.OfDecisionInputUserMessage; vt != nil {
		return (*string)(&vt.Role)
	} else if vt := u.OfDecisionInputDeveloperMessage; vt != nil {
		return (*string)(&vt.Role)
	} else if vt := u.OfDecisionInputSystemMessage; vt != nil {
		return (*string)(&vt.Role)
	} else if vt := u.OfDecisionInputAssistantMessage; vt != nil {
		return (*string)(&vt.Role)
	} else if vt := u.OfDecisionInputMessage; vt != nil {
		return vt.GetRole()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionNewParamsInputArrayItemUnion) GetType() *string {
	if vt := u.OfDecisionInputUserMessage; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfDecisionInputDeveloperMessage; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfDecisionInputSystemMessage; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfDecisionInputAssistantMessage; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfDecisionInputMessage; vt != nil {
		return vt.GetType()
	} else if vt := u.OfFunctionCall; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfFunctionCallOutput; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfCustomToolCall; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfCustomToolCallOutput; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionNewParamsInputArrayItemUnion) GetCallID() *string {
	if vt := u.OfFunctionCall; vt != nil {
		return (*string)(&vt.CallID)
	} else if vt := u.OfFunctionCallOutput; vt != nil {
		return (*string)(&vt.CallID)
	} else if vt := u.OfCustomToolCall; vt != nil {
		return (*string)(&vt.CallID)
	} else if vt := u.OfCustomToolCallOutput; vt != nil {
		return (*string)(&vt.CallID)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionNewParamsInputArrayItemUnion) GetName() *string {
	if vt := u.OfFunctionCall; vt != nil {
		return (*string)(&vt.Name)
	} else if vt := u.OfCustomToolCall; vt != nil {
		return (*string)(&vt.Name)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionNewParamsInputArrayItemUnion) GetID() *string {
	if vt := u.OfFunctionCall; vt != nil && vt.ID.Valid() {
		return &vt.ID.Value
	} else if vt := u.OfFunctionCallOutput; vt != nil && vt.ID.Valid() {
		return &vt.ID.Value
	} else if vt := u.OfCustomToolCall; vt != nil && vt.ID.Valid() {
		return &vt.ID.Value
	} else if vt := u.OfCustomToolCallOutput; vt != nil && vt.ID.Valid() {
		return &vt.ID.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionNewParamsInputArrayItemUnion) GetNamespace() *string {
	if vt := u.OfFunctionCall; vt != nil && vt.Namespace.Valid() {
		return &vt.Namespace.Value
	} else if vt := u.OfCustomToolCall; vt != nil && vt.Namespace.Valid() {
		return &vt.Namespace.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionNewParamsInputArrayItemUnion) GetStatus() *string {
	if vt := u.OfFunctionCall; vt != nil {
		return (*string)(&vt.Status)
	} else if vt := u.OfFunctionCallOutput; vt != nil {
		return (*string)(&vt.Status)
	} else if vt := u.OfCustomToolCall; vt != nil {
		return (*string)(&vt.Status)
	} else if vt := u.OfCustomToolCallOutput; vt != nil {
		return (*string)(&vt.Status)
	}
	return nil
}

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u DecisionNewParamsInputArrayItemUnion) GetContent() (res decisionNewParamsInputArrayItemUnionContent) {
	if vt := u.OfDecisionInputUserMessage; vt != nil {
		res.any = vt.Content.asAny()
	} else if vt := u.OfDecisionInputDeveloperMessage; vt != nil {
		res.any = vt.Content.asAny()
	} else if vt := u.OfDecisionInputSystemMessage; vt != nil {
		res.any = vt.Content.asAny()
	} else if vt := u.OfDecisionInputAssistantMessage; vt != nil {
		res.any = vt.Content.asAny()
	} else if vt := u.OfDecisionInputMessage; vt != nil {
		res.any = vt.GetContent().AsAny()
	}
	return
}

// Can have the runtime types [*string], [_[]DecisionInputPartUnionParam],
// [_[]DecisionNewParamsInputArrayItemDecisionInputDeveloperMessageContentParts],
// [_[]DecisionNewParamsInputArrayItemDecisionInputSystemMessageContentParts],
// [_[]DecisionNewParamsInputArrayItemDecisionInputAssistantMessageContentPartsItemUnion],
// [_[]DecisionInputTextParam],
// [_[]DecisionInputMessageAssistantContentPartsItemUnionParam]
type decisionNewParamsInputArrayItemUnionContent struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *string:
//	case *[]openai.DecisionInputPartUnionParam:
//	case *[]openai.DecisionNewParamsInputArrayItemDecisionInputDeveloperMessageContentParts:
//	case *[]openai.DecisionNewParamsInputArrayItemDecisionInputSystemMessageContentParts:
//	case *[]openai.DecisionNewParamsInputArrayItemDecisionInputAssistantMessageContentPartsItemUnion:
//	case *[]openai.DecisionInputTextParam:
//	case *[]openai.DecisionInputMessageAssistantContentPartsItemUnionParam:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u decisionNewParamsInputArrayItemUnionContent) AsAny() any { return u.any }

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u DecisionNewParamsInputArrayItemUnion) GetOutput() (res decisionNewParamsInputArrayItemUnionOutput) {
	if vt := u.OfFunctionCallOutput; vt != nil {
		res.any = vt.Output.asAny()
	} else if vt := u.OfCustomToolCallOutput; vt != nil {
		res.any = vt.Output.asAny()
	}
	return
}

// Can have the runtime types [*string],
// [_[]DecisionNewParamsInputArrayItemFunctionCallOutputOutputPartsItemUnion],
// [_[]DecisionNewParamsInputArrayItemCustomToolCallOutputOutputPartsItemUnion]
type decisionNewParamsInputArrayItemUnionOutput struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *string:
//	case *[]openai.DecisionNewParamsInputArrayItemFunctionCallOutputOutputPartsItemUnion:
//	case *[]openai.DecisionNewParamsInputArrayItemCustomToolCallOutputOutputPartsItemUnion:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u decisionNewParamsInputArrayItemUnionOutput) AsAny() any { return u.any }

func init() {
	apijson.RegisterUnion[DecisionNewParamsInputArrayItemUnion](
		"type",
		apijson.Discriminator[DecisionNewParamsInputArrayItemDecisionInputUserMessage]("message"),
		apijson.Discriminator[DecisionNewParamsInputArrayItemDecisionInputDeveloperMessage]("message"),
		apijson.Discriminator[DecisionNewParamsInputArrayItemDecisionInputSystemMessage]("message"),
		apijson.Discriminator[DecisionNewParamsInputArrayItemDecisionInputAssistantMessage]("message"),
		apijson.Discriminator[DecisionNewParamsInputArrayItemFunctionCall]("function_call"),
		apijson.Discriminator[DecisionNewParamsInputArrayItemFunctionCallOutput]("function_call_output"),
		apijson.Discriminator[DecisionNewParamsInputArrayItemCustomToolCall]("custom_tool_call"),
		apijson.Discriminator[DecisionNewParamsInputArrayItemCustomToolCallOutput]("custom_tool_call_output"),
		apijson.Discriminator[DecisionNewParamsInputArrayItemDecisionInputUserMessage](nil),
		apijson.Discriminator[DecisionNewParamsInputArrayItemDecisionInputDeveloperMessage](nil),
		apijson.Discriminator[DecisionNewParamsInputArrayItemDecisionInputSystemMessage](nil),
		apijson.Discriminator[DecisionNewParamsInputArrayItemDecisionInputAssistantMessage](nil),
	)
}

// A user message containing input text, images, or JSON.
//
// The properties Content, Role are required.
type DecisionNewParamsInputArrayItemDecisionInputUserMessage struct {
	// A text string or an ordered list of text, image, and JSON parts.
	Content DecisionNewParamsInputArrayItemDecisionInputUserMessageContentUnion `json:"content,omitzero" api:"required"`
	// Any of "message".
	Type string `json:"type,omitzero"`
	// This field can be elided, and will marshal its zero value as "user".
	Role constant.User `json:"role" default:"user"`
	paramObj
}

func (r DecisionNewParamsInputArrayItemDecisionInputUserMessage) MarshalJSON() (data []byte, err error) {
	type shadow DecisionNewParamsInputArrayItemDecisionInputUserMessage
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionNewParamsInputArrayItemDecisionInputUserMessage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[DecisionNewParamsInputArrayItemDecisionInputUserMessage](
		"type", "message",
	)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DecisionNewParamsInputArrayItemDecisionInputUserMessageContentUnion struct {
	OfString param.Opt[string]             `json:",omitzero,inline"`
	OfParts  []DecisionInputPartUnionParam `json:",omitzero,inline"`
	paramUnion
}

func (u DecisionNewParamsInputArrayItemDecisionInputUserMessageContentUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfString, u.OfParts)
}
func (u *DecisionNewParamsInputArrayItemDecisionInputUserMessageContentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *DecisionNewParamsInputArrayItemDecisionInputUserMessageContentUnion) asAny() any {
	if !param.IsOmitted(u.OfString) {
		return &u.OfString.Value
	} else if !param.IsOmitted(u.OfParts) {
		return &u.OfParts
	}
	return nil
}

// A developer message containing input text.
//
// The properties Content, Role are required.
type DecisionNewParamsInputArrayItemDecisionInputDeveloperMessage struct {
	// A text string or an ordered list of input text parts.
	Content DecisionNewParamsInputArrayItemDecisionInputDeveloperMessageContentUnion `json:"content,omitzero" api:"required"`
	// Any of "message".
	Type string `json:"type,omitzero"`
	// This field can be elided, and will marshal its zero value as "developer".
	Role constant.Developer `json:"role" default:"developer"`
	paramObj
}

func (r DecisionNewParamsInputArrayItemDecisionInputDeveloperMessage) MarshalJSON() (data []byte, err error) {
	type shadow DecisionNewParamsInputArrayItemDecisionInputDeveloperMessage
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionNewParamsInputArrayItemDecisionInputDeveloperMessage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[DecisionNewParamsInputArrayItemDecisionInputDeveloperMessage](
		"type", "message",
	)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DecisionNewParamsInputArrayItemDecisionInputDeveloperMessageContentUnion struct {
	OfString param.Opt[string]                                                          `json:",omitzero,inline"`
	OfParts  []DecisionNewParamsInputArrayItemDecisionInputDeveloperMessageContentParts `json:",omitzero,inline"`
	paramUnion
}

func (u DecisionNewParamsInputArrayItemDecisionInputDeveloperMessageContentUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfString, u.OfParts)
}
func (u *DecisionNewParamsInputArrayItemDecisionInputDeveloperMessageContentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *DecisionNewParamsInputArrayItemDecisionInputDeveloperMessageContentUnion) asAny() any {
	if !param.IsOmitted(u.OfString) {
		return &u.OfString.Value
	} else if !param.IsOmitted(u.OfParts) {
		return &u.OfParts
	}
	return nil
}

// A system message containing input text.
//
// The properties Content, Role are required.
type DecisionNewParamsInputArrayItemDecisionInputSystemMessage struct {
	// A text string or an ordered list of input text parts.
	Content DecisionNewParamsInputArrayItemDecisionInputSystemMessageContentUnion `json:"content,omitzero" api:"required"`
	// Any of "message".
	Type string `json:"type,omitzero"`
	// This field can be elided, and will marshal its zero value as "system".
	Role constant.System `json:"role" default:"system"`
	paramObj
}

func (r DecisionNewParamsInputArrayItemDecisionInputSystemMessage) MarshalJSON() (data []byte, err error) {
	type shadow DecisionNewParamsInputArrayItemDecisionInputSystemMessage
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionNewParamsInputArrayItemDecisionInputSystemMessage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[DecisionNewParamsInputArrayItemDecisionInputSystemMessage](
		"type", "message",
	)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DecisionNewParamsInputArrayItemDecisionInputSystemMessageContentUnion struct {
	OfString param.Opt[string]                                                       `json:",omitzero,inline"`
	OfParts  []DecisionNewParamsInputArrayItemDecisionInputSystemMessageContentParts `json:",omitzero,inline"`
	paramUnion
}

func (u DecisionNewParamsInputArrayItemDecisionInputSystemMessageContentUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfString, u.OfParts)
}
func (u *DecisionNewParamsInputArrayItemDecisionInputSystemMessageContentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *DecisionNewParamsInputArrayItemDecisionInputSystemMessageContentUnion) asAny() any {
	if !param.IsOmitted(u.OfString) {
		return &u.OfString.Value
	} else if !param.IsOmitted(u.OfParts) {
		return &u.OfParts
	}
	return nil
}

// An assistant message containing output text or refusals.
//
// The properties Content, Role are required.
type DecisionNewParamsInputArrayItemDecisionInputAssistantMessage struct {
	// A text string or an ordered list of output text and refusal parts.
	Content DecisionNewParamsInputArrayItemDecisionInputAssistantMessageContentUnion `json:"content,omitzero" api:"required"`
	// The phase of an assistant message.
	//
	// Any of "commentary", "final_answer".
	Phase string `json:"phase,omitzero"`
	// Any of "message".
	Type string `json:"type,omitzero"`
	// This field can be elided, and will marshal its zero value as "assistant".
	Role constant.Assistant `json:"role" default:"assistant"`
	paramObj
}

func (r DecisionNewParamsInputArrayItemDecisionInputAssistantMessage) MarshalJSON() (data []byte, err error) {
	type shadow DecisionNewParamsInputArrayItemDecisionInputAssistantMessage
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionNewParamsInputArrayItemDecisionInputAssistantMessage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[DecisionNewParamsInputArrayItemDecisionInputAssistantMessage](
		"phase", "commentary", "final_answer",
	)
	apijson.RegisterFieldValidator[DecisionNewParamsInputArrayItemDecisionInputAssistantMessage](
		"type", "message",
	)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DecisionNewParamsInputArrayItemDecisionInputAssistantMessageContentUnion struct {
	OfString param.Opt[string]                                                                   `json:",omitzero,inline"`
	OfParts  []DecisionNewParamsInputArrayItemDecisionInputAssistantMessageContentPartsItemUnion `json:",omitzero,inline"`
	paramUnion
}

func (u DecisionNewParamsInputArrayItemDecisionInputAssistantMessageContentUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfString, u.OfParts)
}
func (u *DecisionNewParamsInputArrayItemDecisionInputAssistantMessageContentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *DecisionNewParamsInputArrayItemDecisionInputAssistantMessageContentUnion) asAny() any {
	if !param.IsOmitted(u.OfString) {
		return &u.OfString.Value
	} else if !param.IsOmitted(u.OfParts) {
		return &u.OfParts
	}
	return nil
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DecisionNewParamsInputArrayItemDecisionInputAssistantMessageContentPartsItemUnion struct {
	OfOutputText *DecisionNewParamsInputArrayItemDecisionInputAssistantMessageContentPartOutputText `json:",omitzero,inline"`
	OfRefusal    *DecisionNewParamsInputArrayItemDecisionInputAssistantMessageContentPartRefusal    `json:",omitzero,inline"`
	paramUnion
}

func (u DecisionNewParamsInputArrayItemDecisionInputAssistantMessageContentPartsItemUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfOutputText, u.OfRefusal)
}
func (u *DecisionNewParamsInputArrayItemDecisionInputAssistantMessageContentPartsItemUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionNewParamsInputArrayItemDecisionInputAssistantMessageContentPartsItemUnion) GetText() *string {
	if vt := u.OfOutputText; vt != nil {
		return &vt.Text
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionNewParamsInputArrayItemDecisionInputAssistantMessageContentPartsItemUnion) GetRefusal() *string {
	if vt := u.OfRefusal; vt != nil {
		return &vt.Refusal
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionNewParamsInputArrayItemDecisionInputAssistantMessageContentPartsItemUnion) GetType() *string {
	if vt := u.OfOutputText; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfRefusal; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[DecisionNewParamsInputArrayItemDecisionInputAssistantMessageContentPartsItemUnion](
		"type",
		apijson.Discriminator[DecisionNewParamsInputArrayItemDecisionInputAssistantMessageContentPartOutputText]("output_text"),
		apijson.Discriminator[DecisionNewParamsInputArrayItemDecisionInputAssistantMessageContentPartRefusal]("refusal"),
	)
}

// The properties Text, Type are required.
type DecisionNewParamsInputArrayItemDecisionInputAssistantMessageContentPartOutputText struct {
	Text string `json:"text" api:"required"`
	// This field can be elided, and will marshal its zero value as "output_text".
	Type constant.OutputText `json:"type" default:"output_text"`
	paramObj
}

func (r DecisionNewParamsInputArrayItemDecisionInputAssistantMessageContentPartOutputText) MarshalJSON() (data []byte, err error) {
	type shadow DecisionNewParamsInputArrayItemDecisionInputAssistantMessageContentPartOutputText
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionNewParamsInputArrayItemDecisionInputAssistantMessageContentPartOutputText) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Refusal, Type are required.
type DecisionNewParamsInputArrayItemDecisionInputAssistantMessageContentPartRefusal struct {
	Refusal string `json:"refusal" api:"required"`
	// This field can be elided, and will marshal its zero value as "refusal".
	Type constant.Refusal `json:"type" default:"refusal"`
	paramObj
}

func (r DecisionNewParamsInputArrayItemDecisionInputAssistantMessageContentPartRefusal) MarshalJSON() (data []byte, err error) {
	type shadow DecisionNewParamsInputArrayItemDecisionInputAssistantMessageContentPartRefusal
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionNewParamsInputArrayItemDecisionInputAssistantMessageContentPartRefusal) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Arguments, CallID, Name, Type are required.
type DecisionNewParamsInputArrayItemFunctionCall struct {
	Arguments string            `json:"arguments" api:"required"`
	CallID    string            `json:"call_id" api:"required"`
	Name      string            `json:"name" api:"required"`
	ID        param.Opt[string] `json:"id,omitzero"`
	Namespace param.Opt[string] `json:"namespace,omitzero"`
	// Any of "in_progress", "completed", "incomplete".
	Status string `json:"status,omitzero"`
	// This field can be elided, and will marshal its zero value as "function_call".
	Type constant.FunctionCall `json:"type" default:"function_call"`
	paramObj
}

func (r DecisionNewParamsInputArrayItemFunctionCall) MarshalJSON() (data []byte, err error) {
	type shadow DecisionNewParamsInputArrayItemFunctionCall
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionNewParamsInputArrayItemFunctionCall) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[DecisionNewParamsInputArrayItemFunctionCall](
		"status", "in_progress", "completed", "incomplete",
	)
}

// The properties CallID, Output, Type are required.
type DecisionNewParamsInputArrayItemFunctionCallOutput struct {
	CallID string `json:"call_id" api:"required"`
	// Text evidence or an ordered list of text and image parts.
	Output DecisionNewParamsInputArrayItemFunctionCallOutputOutputUnion `json:"output,omitzero" api:"required"`
	ID     param.Opt[string]                                            `json:"id,omitzero"`
	// Any of "in_progress", "completed", "incomplete".
	Status string `json:"status,omitzero"`
	// This field can be elided, and will marshal its zero value as
	// "function_call_output".
	Type constant.FunctionCallOutput `json:"type" default:"function_call_output"`
	paramObj
}

func (r DecisionNewParamsInputArrayItemFunctionCallOutput) MarshalJSON() (data []byte, err error) {
	type shadow DecisionNewParamsInputArrayItemFunctionCallOutput
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionNewParamsInputArrayItemFunctionCallOutput) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[DecisionNewParamsInputArrayItemFunctionCallOutput](
		"status", "in_progress", "completed", "incomplete",
	)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DecisionNewParamsInputArrayItemFunctionCallOutputOutputUnion struct {
	OfString param.Opt[string]                                                       `json:",omitzero,inline"`
	OfParts  []DecisionNewParamsInputArrayItemFunctionCallOutputOutputPartsItemUnion `json:",omitzero,inline"`
	paramUnion
}

func (u DecisionNewParamsInputArrayItemFunctionCallOutputOutputUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfString, u.OfParts)
}
func (u *DecisionNewParamsInputArrayItemFunctionCallOutputOutputUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *DecisionNewParamsInputArrayItemFunctionCallOutputOutputUnion) asAny() any {
	if !param.IsOmitted(u.OfString) {
		return &u.OfString.Value
	} else if !param.IsOmitted(u.OfParts) {
		return &u.OfParts
	}
	return nil
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DecisionNewParamsInputArrayItemFunctionCallOutputOutputPartsItemUnion struct {
	OfInputText  *DecisionInputTextParam  `json:",omitzero,inline"`
	OfInputImage *DecisionInputImageParam `json:",omitzero,inline"`
	paramUnion
}

func (u DecisionNewParamsInputArrayItemFunctionCallOutputOutputPartsItemUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfInputText, u.OfInputImage)
}
func (u *DecisionNewParamsInputArrayItemFunctionCallOutputOutputPartsItemUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionNewParamsInputArrayItemFunctionCallOutputOutputPartsItemUnion) GetText() *string {
	if vt := u.OfInputText; vt != nil {
		return &vt.Text
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionNewParamsInputArrayItemFunctionCallOutputOutputPartsItemUnion) GetImageURL() *string {
	if vt := u.OfInputImage; vt != nil {
		return &vt.ImageURL
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionNewParamsInputArrayItemFunctionCallOutputOutputPartsItemUnion) GetDetail() *string {
	if vt := u.OfInputImage; vt != nil {
		return (*string)(&vt.Detail)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionNewParamsInputArrayItemFunctionCallOutputOutputPartsItemUnion) GetType() *string {
	if vt := u.OfInputText; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfInputImage; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's PromptCacheBreakpoint property, if
// present.
func (u DecisionNewParamsInputArrayItemFunctionCallOutputOutputPartsItemUnion) GetPromptCacheBreakpoint() *DecisionPromptCacheBreakpointParam {
	if vt := u.OfInputText; vt != nil {
		return &vt.PromptCacheBreakpoint
	} else if vt := u.OfInputImage; vt != nil {
		return &vt.PromptCacheBreakpoint
	}
	return nil
}

func init() {
	apijson.RegisterUnion[DecisionNewParamsInputArrayItemFunctionCallOutputOutputPartsItemUnion](
		"type",
		apijson.Discriminator[DecisionInputTextParam]("input_text"),
		apijson.Discriminator[DecisionInputImageParam]("input_image"),
	)
}

// The properties CallID, Input, Name, Type are required.
type DecisionNewParamsInputArrayItemCustomToolCall struct {
	CallID    string            `json:"call_id" api:"required"`
	Input     string            `json:"input" api:"required"`
	Name      string            `json:"name" api:"required"`
	ID        param.Opt[string] `json:"id,omitzero"`
	Namespace param.Opt[string] `json:"namespace,omitzero"`
	// Any of "in_progress", "completed", "incomplete".
	Status string `json:"status,omitzero"`
	// This field can be elided, and will marshal its zero value as "custom_tool_call".
	Type constant.CustomToolCall `json:"type" default:"custom_tool_call"`
	paramObj
}

func (r DecisionNewParamsInputArrayItemCustomToolCall) MarshalJSON() (data []byte, err error) {
	type shadow DecisionNewParamsInputArrayItemCustomToolCall
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionNewParamsInputArrayItemCustomToolCall) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[DecisionNewParamsInputArrayItemCustomToolCall](
		"status", "in_progress", "completed", "incomplete",
	)
}

// The properties CallID, Output, Type are required.
type DecisionNewParamsInputArrayItemCustomToolCallOutput struct {
	CallID string `json:"call_id" api:"required"`
	// Text evidence or an ordered list of text and image parts.
	Output DecisionNewParamsInputArrayItemCustomToolCallOutputOutputUnion `json:"output,omitzero" api:"required"`
	ID     param.Opt[string]                                              `json:"id,omitzero"`
	// Any of "in_progress", "completed", "incomplete".
	Status string `json:"status,omitzero"`
	// This field can be elided, and will marshal its zero value as
	// "custom_tool_call_output".
	Type constant.CustomToolCallOutput `json:"type" default:"custom_tool_call_output"`
	paramObj
}

func (r DecisionNewParamsInputArrayItemCustomToolCallOutput) MarshalJSON() (data []byte, err error) {
	type shadow DecisionNewParamsInputArrayItemCustomToolCallOutput
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionNewParamsInputArrayItemCustomToolCallOutput) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[DecisionNewParamsInputArrayItemCustomToolCallOutput](
		"status", "in_progress", "completed", "incomplete",
	)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DecisionNewParamsInputArrayItemCustomToolCallOutputOutputUnion struct {
	OfString param.Opt[string]                                                         `json:",omitzero,inline"`
	OfParts  []DecisionNewParamsInputArrayItemCustomToolCallOutputOutputPartsItemUnion `json:",omitzero,inline"`
	paramUnion
}

func (u DecisionNewParamsInputArrayItemCustomToolCallOutputOutputUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfString, u.OfParts)
}
func (u *DecisionNewParamsInputArrayItemCustomToolCallOutputOutputUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *DecisionNewParamsInputArrayItemCustomToolCallOutputOutputUnion) asAny() any {
	if !param.IsOmitted(u.OfString) {
		return &u.OfString.Value
	} else if !param.IsOmitted(u.OfParts) {
		return &u.OfParts
	}
	return nil
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DecisionNewParamsInputArrayItemCustomToolCallOutputOutputPartsItemUnion struct {
	OfInputText  *DecisionInputTextParam  `json:",omitzero,inline"`
	OfInputImage *DecisionInputImageParam `json:",omitzero,inline"`
	paramUnion
}

func (u DecisionNewParamsInputArrayItemCustomToolCallOutputOutputPartsItemUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfInputText, u.OfInputImage)
}
func (u *DecisionNewParamsInputArrayItemCustomToolCallOutputOutputPartsItemUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionNewParamsInputArrayItemCustomToolCallOutputOutputPartsItemUnion) GetText() *string {
	if vt := u.OfInputText; vt != nil {
		return &vt.Text
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionNewParamsInputArrayItemCustomToolCallOutputOutputPartsItemUnion) GetImageURL() *string {
	if vt := u.OfInputImage; vt != nil {
		return &vt.ImageURL
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionNewParamsInputArrayItemCustomToolCallOutputOutputPartsItemUnion) GetDetail() *string {
	if vt := u.OfInputImage; vt != nil {
		return (*string)(&vt.Detail)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DecisionNewParamsInputArrayItemCustomToolCallOutputOutputPartsItemUnion) GetType() *string {
	if vt := u.OfInputText; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfInputImage; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's PromptCacheBreakpoint property, if
// present.
func (u DecisionNewParamsInputArrayItemCustomToolCallOutputOutputPartsItemUnion) GetPromptCacheBreakpoint() *DecisionPromptCacheBreakpointParam {
	if vt := u.OfInputText; vt != nil {
		return &vt.PromptCacheBreakpoint
	} else if vt := u.OfInputImage; vt != nil {
		return &vt.PromptCacheBreakpoint
	}
	return nil
}

func init() {
	apijson.RegisterUnion[DecisionNewParamsInputArrayItemCustomToolCallOutputOutputPartsItemUnion](
		"type",
		apijson.Discriminator[DecisionInputTextParam]("input_text"),
		apijson.Discriminator[DecisionInputImageParam]("input_image"),
	)
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

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u DecisionNewParamsQuestionUnion) GetInstructions() (res decisionNewParamsQuestionUnionInstructions) {
	if vt := u.OfPredicate; vt != nil {
		res.any = vt.Instructions.asAny()
	} else if vt := u.OfChoice; vt != nil {
		res.any = vt.Instructions.asAny()
	} else if vt := u.OfScore; vt != nil {
		res.any = vt.Instructions.asAny()
	}
	return
}

// Can have the runtime types [*string],
// [\*[]shared.DecisionQuestionInputPartUnionParam]
type decisionNewParamsQuestionUnionInstructions struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *string:
//	case *[]shared.DecisionQuestionInputPartUnionParam:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u decisionNewParamsQuestionUnionInstructions) AsAny() any { return u.any }

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
	Instructions DecisionNewParamsQuestionPredicateInstructionsUnion `json:"instructions,omitzero" api:"required"`
	Name         param.Opt[string]                                   `json:"name,omitzero"`
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

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DecisionNewParamsQuestionPredicateInstructionsUnion struct {
	OfString                         param.Opt[string]                            `json:",omitzero,inline"`
	OfDecisionQuestionInputPartArray []shared.DecisionQuestionInputPartUnionParam `json:",omitzero,inline"`
	paramUnion
}

func (u DecisionNewParamsQuestionPredicateInstructionsUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfString, u.OfDecisionQuestionInputPartArray)
}
func (u *DecisionNewParamsQuestionPredicateInstructionsUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *DecisionNewParamsQuestionPredicateInstructionsUnion) asAny() any {
	if !param.IsOmitted(u.OfString) {
		return &u.OfString.Value
	} else if !param.IsOmitted(u.OfDecisionQuestionInputPartArray) {
		return &u.OfDecisionQuestionInputPartArray
	}
	return nil
}

// Choose from the supplied options based on the input.
//
// The properties Choices, Instructions, Type are required.
type DecisionNewParamsQuestionChoice struct {
	// Provide between 2 and 255 choices. Each choice must be unique.
	Choices      []DecisionNewParamsQuestionChoiceChoice          `json:"choices,omitzero" api:"required"`
	Instructions DecisionNewParamsQuestionChoiceInstructionsUnion `json:"instructions,omitzero" api:"required"`
	Name         param.Opt[string]                                `json:"name,omitzero"`
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

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DecisionNewParamsQuestionChoiceInstructionsUnion struct {
	OfString                         param.Opt[string]                            `json:",omitzero,inline"`
	OfDecisionQuestionInputPartArray []shared.DecisionQuestionInputPartUnionParam `json:",omitzero,inline"`
	paramUnion
}

func (u DecisionNewParamsQuestionChoiceInstructionsUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfString, u.OfDecisionQuestionInputPartArray)
}
func (u *DecisionNewParamsQuestionChoiceInstructionsUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *DecisionNewParamsQuestionChoiceInstructionsUnion) asAny() any {
	if !param.IsOmitted(u.OfString) {
		return &u.OfString.Value
	} else if !param.IsOmitted(u.OfDecisionQuestionInputPartArray) {
		return &u.OfDecisionQuestionInputPartArray
	}
	return nil
}

// Rate the input against the supplied ordered levels.
//
// The properties Instructions, Levels, Type are required.
type DecisionNewParamsQuestionScore struct {
	Instructions DecisionNewParamsQuestionScoreInstructionsUnion `json:"instructions,omitzero" api:"required"`
	Levels       []DecisionNewParamsQuestionScoreLevel           `json:"levels,omitzero" api:"required"`
	Name         param.Opt[string]                               `json:"name,omitzero"`
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

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DecisionNewParamsQuestionScoreInstructionsUnion struct {
	OfString                         param.Opt[string]                            `json:",omitzero,inline"`
	OfDecisionQuestionInputPartArray []shared.DecisionQuestionInputPartUnionParam `json:",omitzero,inline"`
	paramUnion
}

func (u DecisionNewParamsQuestionScoreInstructionsUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfString, u.OfDecisionQuestionInputPartArray)
}
func (u *DecisionNewParamsQuestionScoreInstructionsUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *DecisionNewParamsQuestionScoreInstructionsUnion) asAny() any {
	if !param.IsOmitted(u.OfString) {
		return &u.OfString.Value
	} else if !param.IsOmitted(u.OfDecisionQuestionInputPartArray) {
		return &u.OfDecisionQuestionInputPartArray
	}
	return nil
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

// Options for prompt caching. Supported for `gpt-5.6` and later models. By
// default, OpenAI automatically chooses one implicit cache breakpoint. You can add
// explicit breakpoints to content blocks with `prompt_cache_breakpoint`. Each
// request can write up to four breakpoints. For cache matching, OpenAI considers
// up to the latest 80 breakpoints in the conversation, without a content-block
// lookback limit. Set `mode` to `explicit` to disable the implicit breakpoint. The
// `ttl` defaults to `30m`, which is currently the only supported value. See the
// [prompt caching guide](https://developers.openai.com/api/docs/guides/prompt-caching)
// for current details.
type DecisionNewParamsPromptCacheOptions struct {
	// Controls whether OpenAI automatically creates an implicit cache breakpoint.
	// Defaults to `implicit`. With `implicit`, OpenAI creates one implicit breakpoint
	// and writes up to the latest three explicit breakpoints in the request. With
	// `explicit`, OpenAI does not create an implicit breakpoint and writes up to the
	// latest four explicit breakpoints. If there are no explicit breakpoints, the
	// request does not use prompt caching.
	//
	// Any of "implicit", "explicit".
	Mode param.Opt[string] `json:"mode,omitzero"`
	// The minimum lifetime applied to every implicit and explicit cache breakpoint
	// written by the request. Defaults to `30m`, which is currently the only supported
	// value. The backend may retain cache entries for longer.
	//
	// Any of "30m".
	Ttl string `json:"ttl,omitzero"`
	paramObj
}

func (r DecisionNewParamsPromptCacheOptions) MarshalJSON() (data []byte, err error) {
	type shadow DecisionNewParamsPromptCacheOptions
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionNewParamsPromptCacheOptions) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func init() {
	apijson.RegisterFieldValidator[DecisionNewParamsPromptCacheOptions](
		"ttl", "30m",
	)
}

// The properties Text, Type are required.
type DecisionNewParamsInputArrayItemDecisionInputDeveloperMessageContentParts struct {
	Text string `json:"text" api:"required"`
	// Marks the exact end of a reusable prompt prefix. The breakpoint inherits its TTL
	// from the request's `prompt_cache_options.ttl`; the boundary is not rounded to a
	// token block.
	PromptCacheBreakpoint DecisionPromptCacheBreakpointParam `json:"prompt_cache_breakpoint,omitzero"`
	// This field can be elided, and will marshal its zero value as "input_text".
	Type constant.InputText `json:"type" default:"input_text"`
	paramObj
}

func (r DecisionNewParamsInputArrayItemDecisionInputDeveloperMessageContentParts) MarshalJSON() (data []byte, err error) {
	type shadow DecisionNewParamsInputArrayItemDecisionInputDeveloperMessageContentParts
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionNewParamsInputArrayItemDecisionInputDeveloperMessageContentParts) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Text, Type are required.
type DecisionNewParamsInputArrayItemDecisionInputSystemMessageContentParts struct {
	Text string `json:"text" api:"required"`
	// Marks the exact end of a reusable prompt prefix. The breakpoint inherits its TTL
	// from the request's `prompt_cache_options.ttl`; the boundary is not rounded to a
	// token block.
	PromptCacheBreakpoint DecisionPromptCacheBreakpointParam `json:"prompt_cache_breakpoint,omitzero"`
	// This field can be elided, and will marshal its zero value as "input_text".
	Type constant.InputText `json:"type" default:"input_text"`
	paramObj
}

func (r DecisionNewParamsInputArrayItemDecisionInputSystemMessageContentParts) MarshalJSON() (data []byte, err error) {
	type shadow DecisionNewParamsInputArrayItemDecisionInputSystemMessageContentParts
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionNewParamsInputArrayItemDecisionInputSystemMessageContentParts) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

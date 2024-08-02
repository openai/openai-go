// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai

import (
	"context"
	"net/http"

	"github.com/openai/openai-go/internal/apijson"
	"github.com/openai/openai-go/internal/param"
	"github.com/openai/openai-go/internal/requestconfig"
	"github.com/openai/openai-go/option"
)

// ModerationService contains methods and other services that help with interacting
// with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewModerationService] method instead.
type ModerationService struct {
	Options []option.RequestOption
}

// NewModerationService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewModerationService(opts ...option.RequestOption) (r *ModerationService) {
	r = &ModerationService{}
	r.Options = opts
	return
}

// Classifies if text is potentially harmful.
func (r *ModerationService) New(ctx context.Context, body ModerationNewParams, opts ...option.RequestOption) (res *ModerationNewResponse, err error) {
	opts = append(r.Options[:], opts...)
	path := "moderations"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

type Moderation struct {
	// A list of the categories, and whether they are flagged or not.
	Categories ModerationCategories `json:"categories,required"`
	// A list of the categories along with their scores as predicted by model.
	CategoryScores ModerationCategoryScores `json:"category_scores,required"`
	// Whether any of the below categories are flagged.
	Flagged bool           `json:"flagged,required"`
	JSON    moderationJSON `json:"-"`
}

// moderationJSON contains the JSON metadata for the struct [Moderation]
type moderationJSON struct {
	Categories     apijson.Field
	CategoryScores apijson.Field
	Flagged        apijson.Field
	raw            string
	ExtraFields    map[string]apijson.Field
}

func (r *Moderation) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r moderationJSON) RawJSON() string {
	return r.raw
}

// A list of the categories, and whether they are flagged or not.
type ModerationCategories struct {
	// Content that expresses, incites, or promotes harassing language towards any
	// target.
	Harassment bool `json:"harassment,required"`
	// Harassment content that also includes violence or serious harm towards any
	// target.
	HarassmentThreatening bool `json:"harassment/threatening,required"`
	// Content that expresses, incites, or promotes hate based on race, gender,
	// ethnicity, religion, nationality, sexual orientation, disability status, or
	// caste. Hateful content aimed at non-protected groups (e.g., chess players) is
	// harassment.
	Hate bool `json:"hate,required"`
	// Hateful content that also includes violence or serious harm towards the targeted
	// group based on race, gender, ethnicity, religion, nationality, sexual
	// orientation, disability status, or caste.
	HateThreatening bool `json:"hate/threatening,required"`
	// Content that promotes, encourages, or depicts acts of self-harm, such as
	// suicide, cutting, and eating disorders.
	SelfHarm bool `json:"self-harm,required"`
	// Content that encourages performing acts of self-harm, such as suicide, cutting,
	// and eating disorders, or that gives instructions or advice on how to commit such
	// acts.
	SelfHarmInstructions bool `json:"self-harm/instructions,required"`
	// Content where the speaker expresses that they are engaging or intend to engage
	// in acts of self-harm, such as suicide, cutting, and eating disorders.
	SelfHarmIntent bool `json:"self-harm/intent,required"`
	// Content meant to arouse sexual excitement, such as the description of sexual
	// activity, or that promotes sexual services (excluding sex education and
	// wellness).
	Sexual bool `json:"sexual,required"`
	// Sexual content that includes an individual who is under 18 years old.
	SexualMinors bool `json:"sexual/minors,required"`
	// Content that depicts death, violence, or physical injury.
	Violence bool `json:"violence,required"`
	// Content that depicts death, violence, or physical injury in graphic detail.
	ViolenceGraphic bool                     `json:"violence/graphic,required"`
	JSON            moderationCategoriesJSON `json:"-"`
}

// moderationCategoriesJSON contains the JSON metadata for the struct
// [ModerationCategories]
type moderationCategoriesJSON struct {
	Harassment            apijson.Field
	HarassmentThreatening apijson.Field
	Hate                  apijson.Field
	HateThreatening       apijson.Field
	SelfHarm              apijson.Field
	SelfHarmInstructions  apijson.Field
	SelfHarmIntent        apijson.Field
	Sexual                apijson.Field
	SexualMinors          apijson.Field
	Violence              apijson.Field
	ViolenceGraphic       apijson.Field
	raw                   string
	ExtraFields           map[string]apijson.Field
}

func (r *ModerationCategories) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r moderationCategoriesJSON) RawJSON() string {
	return r.raw
}

// A list of the categories along with their scores as predicted by model.
type ModerationCategoryScores struct {
	// The score for the category 'harassment'.
	Harassment float64 `json:"harassment,required"`
	// The score for the category 'harassment/threatening'.
	HarassmentThreatening float64 `json:"harassment/threatening,required"`
	// The score for the category 'hate'.
	Hate float64 `json:"hate,required"`
	// The score for the category 'hate/threatening'.
	HateThreatening float64 `json:"hate/threatening,required"`
	// The score for the category 'self-harm'.
	SelfHarm float64 `json:"self-harm,required"`
	// The score for the category 'self-harm/instructions'.
	SelfHarmInstructions float64 `json:"self-harm/instructions,required"`
	// The score for the category 'self-harm/intent'.
	SelfHarmIntent float64 `json:"self-harm/intent,required"`
	// The score for the category 'sexual'.
	Sexual float64 `json:"sexual,required"`
	// The score for the category 'sexual/minors'.
	SexualMinors float64 `json:"sexual/minors,required"`
	// The score for the category 'violence'.
	Violence float64 `json:"violence,required"`
	// The score for the category 'violence/graphic'.
	ViolenceGraphic float64                      `json:"violence/graphic,required"`
	JSON            moderationCategoryScoresJSON `json:"-"`
}

// moderationCategoryScoresJSON contains the JSON metadata for the struct
// [ModerationCategoryScores]
type moderationCategoryScoresJSON struct {
	Harassment            apijson.Field
	HarassmentThreatening apijson.Field
	Hate                  apijson.Field
	HateThreatening       apijson.Field
	SelfHarm              apijson.Field
	SelfHarmInstructions  apijson.Field
	SelfHarmIntent        apijson.Field
	Sexual                apijson.Field
	SexualMinors          apijson.Field
	Violence              apijson.Field
	ViolenceGraphic       apijson.Field
	raw                   string
	ExtraFields           map[string]apijson.Field
}

func (r *ModerationCategoryScores) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r moderationCategoryScoresJSON) RawJSON() string {
	return r.raw
}

type ModerationModel = string

const (
	ModerationModelTextModerationLatest ModerationModel = "text-moderation-latest"
	ModerationModelTextModerationStable ModerationModel = "text-moderation-stable"
)

// Represents if a given text input is potentially harmful.
type ModerationNewResponse struct {
	// The unique identifier for the moderation request.
	ID string `json:"id,required"`
	// The model used to generate the moderation results.
	Model string `json:"model,required"`
	// A list of moderation objects.
	Results []Moderation              `json:"results,required"`
	JSON    moderationNewResponseJSON `json:"-"`
}

// moderationNewResponseJSON contains the JSON metadata for the struct
// [ModerationNewResponse]
type moderationNewResponseJSON struct {
	ID          apijson.Field
	Model       apijson.Field
	Results     apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ModerationNewResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r moderationNewResponseJSON) RawJSON() string {
	return r.raw
}

type ModerationNewParams struct {
	// The input text to classify
	Input param.Field[ModerationNewParamsInputUnion] `json:"input,required"`
	// Two content moderations models are available: `text-moderation-stable` and
	// `text-moderation-latest`.
	//
	// The default is `text-moderation-latest` which will be automatically upgraded
	// over time. This ensures you are always using our most accurate model. If you use
	// `text-moderation-stable`, we will provide advanced notice before updating the
	// model. Accuracy of `text-moderation-stable` may be slightly lower than for
	// `text-moderation-latest`.
	Model param.Field[ModerationModel] `json:"model"`
}

func (r ModerationNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// The input text to classify
//
// Satisfied by [shared.UnionString], [ModerationNewParamsInputArray].
type ModerationNewParamsInputUnion interface {
	ImplementsModerationNewParamsInputUnion()
}

type ModerationNewParamsInputArray []string

func (r ModerationNewParamsInputArray) ImplementsModerationNewParamsInputUnion() {}

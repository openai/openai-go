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

// Classifies if text and/or image inputs are potentially harmful. Learn more in
// the [moderation guide](https://platform.openai.com/docs/guides/moderation).
func (r *ModerationService) New(ctx context.Context, body ModerationNewParams, opts ...option.RequestOption) (res *ModerationNewResponse, err error) {
	opts = append(r.Options[:], opts...)
	path := "moderations"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

type Moderation struct {
	// A list of the categories, and whether they are flagged or not.
	Categories ModerationCategories `json:"categories,required"`
	// A list of the categories along with the input type(s) that the score applies to.
	CategoryAppliedInputTypes ModerationCategoryAppliedInputTypes `json:"category_applied_input_types,required"`
	// A list of the categories along with their scores as predicted by model.
	CategoryScores ModerationCategoryScores `json:"category_scores,required"`
	// Whether any of the below categories are flagged.
	Flagged bool           `json:"flagged,required"`
	JSON    moderationJSON `json:"-"`
}

// moderationJSON contains the JSON metadata for the struct [Moderation]
type moderationJSON struct {
	Categories                apijson.Field
	CategoryAppliedInputTypes apijson.Field
	CategoryScores            apijson.Field
	Flagged                   apijson.Field
	raw                       string
	ExtraFields               map[string]apijson.Field
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
	// Content that includes instructions or advice that facilitate the planning or
	// execution of wrongdoing, or that gives advice or instruction on how to commit
	// illicit acts. For example, "how to shoplift" would fit this category.
	Illicit bool `json:"illicit,required,nullable"`
	// Content that includes instructions or advice that facilitate the planning or
	// execution of wrongdoing that also includes violence, or that gives advice or
	// instruction on the procurement of any weapon.
	IllicitViolent bool `json:"illicit/violent,required,nullable"`
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
	Illicit               apijson.Field
	IllicitViolent        apijson.Field
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

// A list of the categories along with the input type(s) that the score applies to.
type ModerationCategoryAppliedInputTypes struct {
	// The applied input type(s) for the category 'harassment'.
	Harassment []ModerationCategoryAppliedInputTypesHarassment `json:"harassment,required"`
	// The applied input type(s) for the category 'harassment/threatening'.
	HarassmentThreatening []ModerationCategoryAppliedInputTypesHarassmentThreatening `json:"harassment/threatening,required"`
	// The applied input type(s) for the category 'hate'.
	Hate []ModerationCategoryAppliedInputTypesHate `json:"hate,required"`
	// The applied input type(s) for the category 'hate/threatening'.
	HateThreatening []ModerationCategoryAppliedInputTypesHateThreatening `json:"hate/threatening,required"`
	// The applied input type(s) for the category 'illicit'.
	Illicit []ModerationCategoryAppliedInputTypesIllicit `json:"illicit,required"`
	// The applied input type(s) for the category 'illicit/violent'.
	IllicitViolent []ModerationCategoryAppliedInputTypesIllicitViolent `json:"illicit/violent,required"`
	// The applied input type(s) for the category 'self-harm'.
	SelfHarm []ModerationCategoryAppliedInputTypesSelfHarm `json:"self-harm,required"`
	// The applied input type(s) for the category 'self-harm/instructions'.
	SelfHarmInstructions []ModerationCategoryAppliedInputTypesSelfHarmInstruction `json:"self-harm/instructions,required"`
	// The applied input type(s) for the category 'self-harm/intent'.
	SelfHarmIntent []ModerationCategoryAppliedInputTypesSelfHarmIntent `json:"self-harm/intent,required"`
	// The applied input type(s) for the category 'sexual'.
	Sexual []ModerationCategoryAppliedInputTypesSexual `json:"sexual,required"`
	// The applied input type(s) for the category 'sexual/minors'.
	SexualMinors []ModerationCategoryAppliedInputTypesSexualMinor `json:"sexual/minors,required"`
	// The applied input type(s) for the category 'violence'.
	Violence []ModerationCategoryAppliedInputTypesViolence `json:"violence,required"`
	// The applied input type(s) for the category 'violence/graphic'.
	ViolenceGraphic []ModerationCategoryAppliedInputTypesViolenceGraphic `json:"violence/graphic,required"`
	JSON            moderationCategoryAppliedInputTypesJSON              `json:"-"`
}

// moderationCategoryAppliedInputTypesJSON contains the JSON metadata for the
// struct [ModerationCategoryAppliedInputTypes]
type moderationCategoryAppliedInputTypesJSON struct {
	Harassment            apijson.Field
	HarassmentThreatening apijson.Field
	Hate                  apijson.Field
	HateThreatening       apijson.Field
	Illicit               apijson.Field
	IllicitViolent        apijson.Field
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

func (r *ModerationCategoryAppliedInputTypes) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r moderationCategoryAppliedInputTypesJSON) RawJSON() string {
	return r.raw
}

type ModerationCategoryAppliedInputTypesHarassment string

const (
	ModerationCategoryAppliedInputTypesHarassmentText ModerationCategoryAppliedInputTypesHarassment = "text"
)

func (r ModerationCategoryAppliedInputTypesHarassment) IsKnown() bool {
	switch r {
	case ModerationCategoryAppliedInputTypesHarassmentText:
		return true
	}
	return false
}

type ModerationCategoryAppliedInputTypesHarassmentThreatening string

const (
	ModerationCategoryAppliedInputTypesHarassmentThreateningText ModerationCategoryAppliedInputTypesHarassmentThreatening = "text"
)

func (r ModerationCategoryAppliedInputTypesHarassmentThreatening) IsKnown() bool {
	switch r {
	case ModerationCategoryAppliedInputTypesHarassmentThreateningText:
		return true
	}
	return false
}

type ModerationCategoryAppliedInputTypesHate string

const (
	ModerationCategoryAppliedInputTypesHateText ModerationCategoryAppliedInputTypesHate = "text"
)

func (r ModerationCategoryAppliedInputTypesHate) IsKnown() bool {
	switch r {
	case ModerationCategoryAppliedInputTypesHateText:
		return true
	}
	return false
}

type ModerationCategoryAppliedInputTypesHateThreatening string

const (
	ModerationCategoryAppliedInputTypesHateThreateningText ModerationCategoryAppliedInputTypesHateThreatening = "text"
)

func (r ModerationCategoryAppliedInputTypesHateThreatening) IsKnown() bool {
	switch r {
	case ModerationCategoryAppliedInputTypesHateThreateningText:
		return true
	}
	return false
}

type ModerationCategoryAppliedInputTypesIllicit string

const (
	ModerationCategoryAppliedInputTypesIllicitText ModerationCategoryAppliedInputTypesIllicit = "text"
)

func (r ModerationCategoryAppliedInputTypesIllicit) IsKnown() bool {
	switch r {
	case ModerationCategoryAppliedInputTypesIllicitText:
		return true
	}
	return false
}

type ModerationCategoryAppliedInputTypesIllicitViolent string

const (
	ModerationCategoryAppliedInputTypesIllicitViolentText ModerationCategoryAppliedInputTypesIllicitViolent = "text"
)

func (r ModerationCategoryAppliedInputTypesIllicitViolent) IsKnown() bool {
	switch r {
	case ModerationCategoryAppliedInputTypesIllicitViolentText:
		return true
	}
	return false
}

type ModerationCategoryAppliedInputTypesSelfHarm string

const (
	ModerationCategoryAppliedInputTypesSelfHarmText  ModerationCategoryAppliedInputTypesSelfHarm = "text"
	ModerationCategoryAppliedInputTypesSelfHarmImage ModerationCategoryAppliedInputTypesSelfHarm = "image"
)

func (r ModerationCategoryAppliedInputTypesSelfHarm) IsKnown() bool {
	switch r {
	case ModerationCategoryAppliedInputTypesSelfHarmText, ModerationCategoryAppliedInputTypesSelfHarmImage:
		return true
	}
	return false
}

type ModerationCategoryAppliedInputTypesSelfHarmInstruction string

const (
	ModerationCategoryAppliedInputTypesSelfHarmInstructionText  ModerationCategoryAppliedInputTypesSelfHarmInstruction = "text"
	ModerationCategoryAppliedInputTypesSelfHarmInstructionImage ModerationCategoryAppliedInputTypesSelfHarmInstruction = "image"
)

func (r ModerationCategoryAppliedInputTypesSelfHarmInstruction) IsKnown() bool {
	switch r {
	case ModerationCategoryAppliedInputTypesSelfHarmInstructionText, ModerationCategoryAppliedInputTypesSelfHarmInstructionImage:
		return true
	}
	return false
}

type ModerationCategoryAppliedInputTypesSelfHarmIntent string

const (
	ModerationCategoryAppliedInputTypesSelfHarmIntentText  ModerationCategoryAppliedInputTypesSelfHarmIntent = "text"
	ModerationCategoryAppliedInputTypesSelfHarmIntentImage ModerationCategoryAppliedInputTypesSelfHarmIntent = "image"
)

func (r ModerationCategoryAppliedInputTypesSelfHarmIntent) IsKnown() bool {
	switch r {
	case ModerationCategoryAppliedInputTypesSelfHarmIntentText, ModerationCategoryAppliedInputTypesSelfHarmIntentImage:
		return true
	}
	return false
}

type ModerationCategoryAppliedInputTypesSexual string

const (
	ModerationCategoryAppliedInputTypesSexualText  ModerationCategoryAppliedInputTypesSexual = "text"
	ModerationCategoryAppliedInputTypesSexualImage ModerationCategoryAppliedInputTypesSexual = "image"
)

func (r ModerationCategoryAppliedInputTypesSexual) IsKnown() bool {
	switch r {
	case ModerationCategoryAppliedInputTypesSexualText, ModerationCategoryAppliedInputTypesSexualImage:
		return true
	}
	return false
}

type ModerationCategoryAppliedInputTypesSexualMinor string

const (
	ModerationCategoryAppliedInputTypesSexualMinorText ModerationCategoryAppliedInputTypesSexualMinor = "text"
)

func (r ModerationCategoryAppliedInputTypesSexualMinor) IsKnown() bool {
	switch r {
	case ModerationCategoryAppliedInputTypesSexualMinorText:
		return true
	}
	return false
}

type ModerationCategoryAppliedInputTypesViolence string

const (
	ModerationCategoryAppliedInputTypesViolenceText  ModerationCategoryAppliedInputTypesViolence = "text"
	ModerationCategoryAppliedInputTypesViolenceImage ModerationCategoryAppliedInputTypesViolence = "image"
)

func (r ModerationCategoryAppliedInputTypesViolence) IsKnown() bool {
	switch r {
	case ModerationCategoryAppliedInputTypesViolenceText, ModerationCategoryAppliedInputTypesViolenceImage:
		return true
	}
	return false
}

type ModerationCategoryAppliedInputTypesViolenceGraphic string

const (
	ModerationCategoryAppliedInputTypesViolenceGraphicText  ModerationCategoryAppliedInputTypesViolenceGraphic = "text"
	ModerationCategoryAppliedInputTypesViolenceGraphicImage ModerationCategoryAppliedInputTypesViolenceGraphic = "image"
)

func (r ModerationCategoryAppliedInputTypesViolenceGraphic) IsKnown() bool {
	switch r {
	case ModerationCategoryAppliedInputTypesViolenceGraphicText, ModerationCategoryAppliedInputTypesViolenceGraphicImage:
		return true
	}
	return false
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
	// The score for the category 'illicit'.
	Illicit float64 `json:"illicit,required"`
	// The score for the category 'illicit/violent'.
	IllicitViolent float64 `json:"illicit/violent,required"`
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
	Illicit               apijson.Field
	IllicitViolent        apijson.Field
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

// An object describing an image to classify.
type ModerationImageURLInputParam struct {
	// Contains either an image URL or a data URL for a base64 encoded image.
	ImageURL param.Field[ModerationImageURLInputImageURLParam] `json:"image_url,required"`
	// Always `image_url`.
	Type param.Field[ModerationImageURLInputType] `json:"type,required"`
}

func (r ModerationImageURLInputParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ModerationImageURLInputParam) implementsModerationMultiModalInputUnionParam() {}

// Contains either an image URL or a data URL for a base64 encoded image.
type ModerationImageURLInputImageURLParam struct {
	// Either a URL of the image or the base64 encoded image data.
	URL param.Field[string] `json:"url,required" format:"uri"`
}

func (r ModerationImageURLInputImageURLParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Always `image_url`.
type ModerationImageURLInputType string

const (
	ModerationImageURLInputTypeImageURL ModerationImageURLInputType = "image_url"
)

func (r ModerationImageURLInputType) IsKnown() bool {
	switch r {
	case ModerationImageURLInputTypeImageURL:
		return true
	}
	return false
}

type ModerationModel = string

const (
	ModerationModelOmniModerationLatest     ModerationModel = "omni-moderation-latest"
	ModerationModelOmniModeration2024_09_26 ModerationModel = "omni-moderation-2024-09-26"
	ModerationModelTextModerationLatest     ModerationModel = "text-moderation-latest"
	ModerationModelTextModerationStable     ModerationModel = "text-moderation-stable"
)

// An object describing an image to classify.
type ModerationMultiModalInputParam struct {
	// Always `image_url`.
	Type     param.Field[ModerationMultiModalInputType] `json:"type,required"`
	ImageURL param.Field[interface{}]                   `json:"image_url"`
	// A string of text to classify.
	Text param.Field[string] `json:"text"`
}

func (r ModerationMultiModalInputParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ModerationMultiModalInputParam) implementsModerationMultiModalInputUnionParam() {}

// An object describing an image to classify.
//
// Satisfied by [ModerationImageURLInputParam], [ModerationTextInputParam],
// [ModerationMultiModalInputParam].
type ModerationMultiModalInputUnionParam interface {
	implementsModerationMultiModalInputUnionParam()
}

// Always `image_url`.
type ModerationMultiModalInputType string

const (
	ModerationMultiModalInputTypeImageURL ModerationMultiModalInputType = "image_url"
	ModerationMultiModalInputTypeText     ModerationMultiModalInputType = "text"
)

func (r ModerationMultiModalInputType) IsKnown() bool {
	switch r {
	case ModerationMultiModalInputTypeImageURL, ModerationMultiModalInputTypeText:
		return true
	}
	return false
}

// An object describing text to classify.
type ModerationTextInputParam struct {
	// A string of text to classify.
	Text param.Field[string] `json:"text,required"`
	// Always `text`.
	Type param.Field[ModerationTextInputType] `json:"type,required"`
}

func (r ModerationTextInputParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ModerationTextInputParam) implementsModerationMultiModalInputUnionParam() {}

// Always `text`.
type ModerationTextInputType string

const (
	ModerationTextInputTypeText ModerationTextInputType = "text"
)

func (r ModerationTextInputType) IsKnown() bool {
	switch r {
	case ModerationTextInputTypeText:
		return true
	}
	return false
}

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
	// Input (or inputs) to classify. Can be a single string, an array of strings, or
	// an array of multi-modal input objects similar to other models.
	Input param.Field[ModerationNewParamsInputUnion] `json:"input,required"`
	// The content moderation model you would like to use. Learn more in
	// [the moderation guide](https://platform.openai.com/docs/guides/moderation), and
	// learn about available models
	// [here](https://platform.openai.com/docs/models#moderation).
	Model param.Field[ModerationModel] `json:"model"`
}

func (r ModerationNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Input (or inputs) to classify. Can be a single string, an array of strings, or
// an array of multi-modal input objects similar to other models.
//
// Satisfied by [shared.UnionString], [ModerationNewParamsInputArray],
// [ModerationNewParamsInputModerationMultiModalArray].
type ModerationNewParamsInputUnion interface {
	ImplementsModerationNewParamsInputUnion()
}

type ModerationNewParamsInputArray []string

func (r ModerationNewParamsInputArray) ImplementsModerationNewParamsInputUnion() {}

type ModerationNewParamsInputModerationMultiModalArray []ModerationMultiModalInputUnionParam

func (r ModerationNewParamsInputModerationMultiModalArray) ImplementsModerationNewParamsInputUnion() {
}

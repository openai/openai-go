// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai

import (
	"encoding/json"

	"github.com/openai/openai-go/internal/apijson"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/param"
	"github.com/openai/openai-go/packages/respjson"
	"github.com/openai/openai-go/shared/constant"
)

// FineTuningMethodService contains methods and other services that help with
// interacting with the openai API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewFineTuningMethodService] method instead.
type FineTuningMethodService struct {
	Options []option.RequestOption
}

// NewFineTuningMethodService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewFineTuningMethodService(opts ...option.RequestOption) (r FineTuningMethodService) {
	r = FineTuningMethodService{}
	r.Options = opts
	return
}

// Configuration for the DPO fine-tuning method.
type DpoMethod struct {
	// The hyperparameters used for the fine-tuning job.
	Hyperparameters DpoMethodHyperparameters `json:"hyperparameters"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Hyperparameters respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DpoMethod) RawJSON() string { return r.JSON.raw }
func (r *DpoMethod) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this DpoMethod to a DpoMethodParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// DpoMethodParam.Overrides()
func (r DpoMethod) ToParam() DpoMethodParam {
	return param.Override[DpoMethodParam](json.RawMessage(r.RawJSON()))
}

// The hyperparameters used for the fine-tuning job.
type DpoMethodHyperparameters struct {
	// Number of examples in each batch. A larger batch size means that model
	// parameters are updated less frequently, but with lower variance.
	BatchSize DpoMethodHyperparametersBatchSizeUnion `json:"batch_size"`
	// The beta value for the DPO method. A higher beta value will increase the weight
	// of the penalty between the policy and reference model.
	Beta DpoMethodHyperparametersBetaUnion `json:"beta"`
	// Scaling factor for the learning rate. A smaller learning rate may be useful to
	// avoid overfitting.
	LearningRateMultiplier DpoMethodHyperparametersLearningRateMultiplierUnion `json:"learning_rate_multiplier"`
	// The number of epochs to train the model for. An epoch refers to one full cycle
	// through the training dataset.
	NEpochs DpoMethodHyperparametersNEpochsUnion `json:"n_epochs"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		BatchSize              respjson.Field
		Beta                   respjson.Field
		LearningRateMultiplier respjson.Field
		NEpochs                respjson.Field
		ExtraFields            map[string]respjson.Field
		raw                    string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DpoMethodHyperparameters) RawJSON() string { return r.JSON.raw }
func (r *DpoMethodHyperparameters) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// DpoMethodHyperparametersBatchSizeUnion contains all possible properties and
// values from [constant.Auto], [int64].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfAuto OfInt]
type DpoMethodHyperparametersBatchSizeUnion struct {
	// This field will be present if the value is a [constant.Auto] instead of an
	// object.
	OfAuto constant.Auto `json:",inline"`
	// This field will be present if the value is a [int64] instead of an object.
	OfInt int64 `json:",inline"`
	JSON  struct {
		OfAuto respjson.Field
		OfInt  respjson.Field
		raw    string
	} `json:"-"`
}

func (u DpoMethodHyperparametersBatchSizeUnion) AsAuto() (v constant.Auto) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u DpoMethodHyperparametersBatchSizeUnion) AsInt() (v int64) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u DpoMethodHyperparametersBatchSizeUnion) RawJSON() string { return u.JSON.raw }

func (r *DpoMethodHyperparametersBatchSizeUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// DpoMethodHyperparametersBetaUnion contains all possible properties and values
// from [constant.Auto], [float64].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfAuto OfFloat]
type DpoMethodHyperparametersBetaUnion struct {
	// This field will be present if the value is a [constant.Auto] instead of an
	// object.
	OfAuto constant.Auto `json:",inline"`
	// This field will be present if the value is a [float64] instead of an object.
	OfFloat float64 `json:",inline"`
	JSON    struct {
		OfAuto  respjson.Field
		OfFloat respjson.Field
		raw     string
	} `json:"-"`
}

func (u DpoMethodHyperparametersBetaUnion) AsAuto() (v constant.Auto) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u DpoMethodHyperparametersBetaUnion) AsFloat() (v float64) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u DpoMethodHyperparametersBetaUnion) RawJSON() string { return u.JSON.raw }

func (r *DpoMethodHyperparametersBetaUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// DpoMethodHyperparametersLearningRateMultiplierUnion contains all possible
// properties and values from [constant.Auto], [float64].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfAuto OfFloat]
type DpoMethodHyperparametersLearningRateMultiplierUnion struct {
	// This field will be present if the value is a [constant.Auto] instead of an
	// object.
	OfAuto constant.Auto `json:",inline"`
	// This field will be present if the value is a [float64] instead of an object.
	OfFloat float64 `json:",inline"`
	JSON    struct {
		OfAuto  respjson.Field
		OfFloat respjson.Field
		raw     string
	} `json:"-"`
}

func (u DpoMethodHyperparametersLearningRateMultiplierUnion) AsAuto() (v constant.Auto) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u DpoMethodHyperparametersLearningRateMultiplierUnion) AsFloat() (v float64) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u DpoMethodHyperparametersLearningRateMultiplierUnion) RawJSON() string { return u.JSON.raw }

func (r *DpoMethodHyperparametersLearningRateMultiplierUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// DpoMethodHyperparametersNEpochsUnion contains all possible properties and values
// from [constant.Auto], [int64].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfAuto OfInt]
type DpoMethodHyperparametersNEpochsUnion struct {
	// This field will be present if the value is a [constant.Auto] instead of an
	// object.
	OfAuto constant.Auto `json:",inline"`
	// This field will be present if the value is a [int64] instead of an object.
	OfInt int64 `json:",inline"`
	JSON  struct {
		OfAuto respjson.Field
		OfInt  respjson.Field
		raw    string
	} `json:"-"`
}

func (u DpoMethodHyperparametersNEpochsUnion) AsAuto() (v constant.Auto) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u DpoMethodHyperparametersNEpochsUnion) AsInt() (v int64) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u DpoMethodHyperparametersNEpochsUnion) RawJSON() string { return u.JSON.raw }

func (r *DpoMethodHyperparametersNEpochsUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for the DPO fine-tuning method.
type DpoMethodParam struct {
	// The hyperparameters used for the fine-tuning job.
	Hyperparameters DpoMethodHyperparametersParam `json:"hyperparameters,omitzero"`
	paramObj
}

func (r DpoMethodParam) MarshalJSON() (data []byte, err error) {
	type shadow DpoMethodParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DpoMethodParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The hyperparameters used for the fine-tuning job.
type DpoMethodHyperparametersParam struct {
	// Number of examples in each batch. A larger batch size means that model
	// parameters are updated less frequently, but with lower variance.
	BatchSize DpoMethodHyperparametersBatchSizeUnionParam `json:"batch_size,omitzero"`
	// The beta value for the DPO method. A higher beta value will increase the weight
	// of the penalty between the policy and reference model.
	Beta DpoMethodHyperparametersBetaUnionParam `json:"beta,omitzero"`
	// Scaling factor for the learning rate. A smaller learning rate may be useful to
	// avoid overfitting.
	LearningRateMultiplier DpoMethodHyperparametersLearningRateMultiplierUnionParam `json:"learning_rate_multiplier,omitzero"`
	// The number of epochs to train the model for. An epoch refers to one full cycle
	// through the training dataset.
	NEpochs DpoMethodHyperparametersNEpochsUnionParam `json:"n_epochs,omitzero"`
	paramObj
}

func (r DpoMethodHyperparametersParam) MarshalJSON() (data []byte, err error) {
	type shadow DpoMethodHyperparametersParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DpoMethodHyperparametersParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DpoMethodHyperparametersBatchSizeUnionParam struct {
	// Construct this variant with constant.ValueOf[constant.Auto]()
	OfAuto constant.Auto    `json:",omitzero,inline"`
	OfInt  param.Opt[int64] `json:",omitzero,inline"`
	paramUnion
}

func (u DpoMethodHyperparametersBatchSizeUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfAuto, u.OfInt)
}
func (u *DpoMethodHyperparametersBatchSizeUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *DpoMethodHyperparametersBatchSizeUnionParam) asAny() any {
	if !param.IsOmitted(u.OfAuto) {
		return &u.OfAuto
	} else if !param.IsOmitted(u.OfInt) {
		return &u.OfInt.Value
	}
	return nil
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DpoMethodHyperparametersBetaUnionParam struct {
	// Construct this variant with constant.ValueOf[constant.Auto]()
	OfAuto  constant.Auto      `json:",omitzero,inline"`
	OfFloat param.Opt[float64] `json:",omitzero,inline"`
	paramUnion
}

func (u DpoMethodHyperparametersBetaUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfAuto, u.OfFloat)
}
func (u *DpoMethodHyperparametersBetaUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *DpoMethodHyperparametersBetaUnionParam) asAny() any {
	if !param.IsOmitted(u.OfAuto) {
		return &u.OfAuto
	} else if !param.IsOmitted(u.OfFloat) {
		return &u.OfFloat.Value
	}
	return nil
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DpoMethodHyperparametersLearningRateMultiplierUnionParam struct {
	// Construct this variant with constant.ValueOf[constant.Auto]()
	OfAuto  constant.Auto      `json:",omitzero,inline"`
	OfFloat param.Opt[float64] `json:",omitzero,inline"`
	paramUnion
}

func (u DpoMethodHyperparametersLearningRateMultiplierUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfAuto, u.OfFloat)
}
func (u *DpoMethodHyperparametersLearningRateMultiplierUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *DpoMethodHyperparametersLearningRateMultiplierUnionParam) asAny() any {
	if !param.IsOmitted(u.OfAuto) {
		return &u.OfAuto
	} else if !param.IsOmitted(u.OfFloat) {
		return &u.OfFloat.Value
	}
	return nil
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DpoMethodHyperparametersNEpochsUnionParam struct {
	// Construct this variant with constant.ValueOf[constant.Auto]()
	OfAuto constant.Auto    `json:",omitzero,inline"`
	OfInt  param.Opt[int64] `json:",omitzero,inline"`
	paramUnion
}

func (u DpoMethodHyperparametersNEpochsUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfAuto, u.OfInt)
}
func (u *DpoMethodHyperparametersNEpochsUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *DpoMethodHyperparametersNEpochsUnionParam) asAny() any {
	if !param.IsOmitted(u.OfAuto) {
		return &u.OfAuto
	} else if !param.IsOmitted(u.OfInt) {
		return &u.OfInt.Value
	}
	return nil
}

// Configuration for the supervised fine-tuning method.
type SupervisedMethod struct {
	// The hyperparameters used for the fine-tuning job.
	Hyperparameters SupervisedMethodHyperparameters `json:"hyperparameters"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Hyperparameters respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r SupervisedMethod) RawJSON() string { return r.JSON.raw }
func (r *SupervisedMethod) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this SupervisedMethod to a SupervisedMethodParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// SupervisedMethodParam.Overrides()
func (r SupervisedMethod) ToParam() SupervisedMethodParam {
	return param.Override[SupervisedMethodParam](json.RawMessage(r.RawJSON()))
}

// The hyperparameters used for the fine-tuning job.
type SupervisedMethodHyperparameters struct {
	// Number of examples in each batch. A larger batch size means that model
	// parameters are updated less frequently, but with lower variance.
	BatchSize SupervisedMethodHyperparametersBatchSizeUnion `json:"batch_size"`
	// Scaling factor for the learning rate. A smaller learning rate may be useful to
	// avoid overfitting.
	LearningRateMultiplier SupervisedMethodHyperparametersLearningRateMultiplierUnion `json:"learning_rate_multiplier"`
	// The number of epochs to train the model for. An epoch refers to one full cycle
	// through the training dataset.
	NEpochs SupervisedMethodHyperparametersNEpochsUnion `json:"n_epochs"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		BatchSize              respjson.Field
		LearningRateMultiplier respjson.Field
		NEpochs                respjson.Field
		ExtraFields            map[string]respjson.Field
		raw                    string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r SupervisedMethodHyperparameters) RawJSON() string { return r.JSON.raw }
func (r *SupervisedMethodHyperparameters) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// SupervisedMethodHyperparametersBatchSizeUnion contains all possible properties
// and values from [constant.Auto], [int64].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfAuto OfInt]
type SupervisedMethodHyperparametersBatchSizeUnion struct {
	// This field will be present if the value is a [constant.Auto] instead of an
	// object.
	OfAuto constant.Auto `json:",inline"`
	// This field will be present if the value is a [int64] instead of an object.
	OfInt int64 `json:",inline"`
	JSON  struct {
		OfAuto respjson.Field
		OfInt  respjson.Field
		raw    string
	} `json:"-"`
}

func (u SupervisedMethodHyperparametersBatchSizeUnion) AsAuto() (v constant.Auto) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u SupervisedMethodHyperparametersBatchSizeUnion) AsInt() (v int64) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u SupervisedMethodHyperparametersBatchSizeUnion) RawJSON() string { return u.JSON.raw }

func (r *SupervisedMethodHyperparametersBatchSizeUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// SupervisedMethodHyperparametersLearningRateMultiplierUnion contains all possible
// properties and values from [constant.Auto], [float64].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfAuto OfFloat]
type SupervisedMethodHyperparametersLearningRateMultiplierUnion struct {
	// This field will be present if the value is a [constant.Auto] instead of an
	// object.
	OfAuto constant.Auto `json:",inline"`
	// This field will be present if the value is a [float64] instead of an object.
	OfFloat float64 `json:",inline"`
	JSON    struct {
		OfAuto  respjson.Field
		OfFloat respjson.Field
		raw     string
	} `json:"-"`
}

func (u SupervisedMethodHyperparametersLearningRateMultiplierUnion) AsAuto() (v constant.Auto) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u SupervisedMethodHyperparametersLearningRateMultiplierUnion) AsFloat() (v float64) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u SupervisedMethodHyperparametersLearningRateMultiplierUnion) RawJSON() string {
	return u.JSON.raw
}

func (r *SupervisedMethodHyperparametersLearningRateMultiplierUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// SupervisedMethodHyperparametersNEpochsUnion contains all possible properties and
// values from [constant.Auto], [int64].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfAuto OfInt]
type SupervisedMethodHyperparametersNEpochsUnion struct {
	// This field will be present if the value is a [constant.Auto] instead of an
	// object.
	OfAuto constant.Auto `json:",inline"`
	// This field will be present if the value is a [int64] instead of an object.
	OfInt int64 `json:",inline"`
	JSON  struct {
		OfAuto respjson.Field
		OfInt  respjson.Field
		raw    string
	} `json:"-"`
}

func (u SupervisedMethodHyperparametersNEpochsUnion) AsAuto() (v constant.Auto) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u SupervisedMethodHyperparametersNEpochsUnion) AsInt() (v int64) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u SupervisedMethodHyperparametersNEpochsUnion) RawJSON() string { return u.JSON.raw }

func (r *SupervisedMethodHyperparametersNEpochsUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for the supervised fine-tuning method.
type SupervisedMethodParam struct {
	// The hyperparameters used for the fine-tuning job.
	Hyperparameters SupervisedMethodHyperparametersParam `json:"hyperparameters,omitzero"`
	paramObj
}

func (r SupervisedMethodParam) MarshalJSON() (data []byte, err error) {
	type shadow SupervisedMethodParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *SupervisedMethodParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The hyperparameters used for the fine-tuning job.
type SupervisedMethodHyperparametersParam struct {
	// Number of examples in each batch. A larger batch size means that model
	// parameters are updated less frequently, but with lower variance.
	BatchSize SupervisedMethodHyperparametersBatchSizeUnionParam `json:"batch_size,omitzero"`
	// Scaling factor for the learning rate. A smaller learning rate may be useful to
	// avoid overfitting.
	LearningRateMultiplier SupervisedMethodHyperparametersLearningRateMultiplierUnionParam `json:"learning_rate_multiplier,omitzero"`
	// The number of epochs to train the model for. An epoch refers to one full cycle
	// through the training dataset.
	NEpochs SupervisedMethodHyperparametersNEpochsUnionParam `json:"n_epochs,omitzero"`
	paramObj
}

func (r SupervisedMethodHyperparametersParam) MarshalJSON() (data []byte, err error) {
	type shadow SupervisedMethodHyperparametersParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *SupervisedMethodHyperparametersParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type SupervisedMethodHyperparametersBatchSizeUnionParam struct {
	// Construct this variant with constant.ValueOf[constant.Auto]()
	OfAuto constant.Auto    `json:",omitzero,inline"`
	OfInt  param.Opt[int64] `json:",omitzero,inline"`
	paramUnion
}

func (u SupervisedMethodHyperparametersBatchSizeUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfAuto, u.OfInt)
}
func (u *SupervisedMethodHyperparametersBatchSizeUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *SupervisedMethodHyperparametersBatchSizeUnionParam) asAny() any {
	if !param.IsOmitted(u.OfAuto) {
		return &u.OfAuto
	} else if !param.IsOmitted(u.OfInt) {
		return &u.OfInt.Value
	}
	return nil
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type SupervisedMethodHyperparametersLearningRateMultiplierUnionParam struct {
	// Construct this variant with constant.ValueOf[constant.Auto]()
	OfAuto  constant.Auto      `json:",omitzero,inline"`
	OfFloat param.Opt[float64] `json:",omitzero,inline"`
	paramUnion
}

func (u SupervisedMethodHyperparametersLearningRateMultiplierUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfAuto, u.OfFloat)
}
func (u *SupervisedMethodHyperparametersLearningRateMultiplierUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *SupervisedMethodHyperparametersLearningRateMultiplierUnionParam) asAny() any {
	if !param.IsOmitted(u.OfAuto) {
		return &u.OfAuto
	} else if !param.IsOmitted(u.OfFloat) {
		return &u.OfFloat.Value
	}
	return nil
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type SupervisedMethodHyperparametersNEpochsUnionParam struct {
	// Construct this variant with constant.ValueOf[constant.Auto]()
	OfAuto constant.Auto    `json:",omitzero,inline"`
	OfInt  param.Opt[int64] `json:",omitzero,inline"`
	paramUnion
}

func (u SupervisedMethodHyperparametersNEpochsUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfAuto, u.OfInt)
}
func (u *SupervisedMethodHyperparametersNEpochsUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *SupervisedMethodHyperparametersNEpochsUnionParam) asAny() any {
	if !param.IsOmitted(u.OfAuto) {
		return &u.OfAuto
	} else if !param.IsOmitted(u.OfInt) {
		return &u.OfInt.Value
	}
	return nil
}

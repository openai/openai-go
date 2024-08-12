// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package shared

type UnionString string

func (UnionString) ImplementsCompletionNewParamsPromptUnion()   {}
func (UnionString) ImplementsCompletionNewParamsStopUnion()     {}
func (UnionString) ImplementsChatCompletionNewParamsStopUnion() {}
func (UnionString) ImplementsEmbeddingNewParamsInputUnion()     {}
func (UnionString) ImplementsModerationNewParamsInputUnion()    {}

type UnionInt int64

func (UnionInt) ImplementsFineTuningJobHyperparametersNEpochsUnion()            {}
func (UnionInt) ImplementsFineTuningJobNewParamsHyperparametersBatchSizeUnion() {}
func (UnionInt) ImplementsFineTuningJobNewParamsHyperparametersNEpochsUnion()   {}

type UnionFloat float64

func (UnionFloat) ImplementsFineTuningJobNewParamsHyperparametersLearningRateMultiplierUnion() {}

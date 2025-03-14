// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package shared

type UnionString string

func (UnionString) ImplementsCompletionNewParamsPromptUnion()   {}
func (UnionString) ImplementsCompletionNewParamsStopUnion()     {}
func (UnionString) ImplementsChatCompletionNewParamsStopUnion() {}
func (UnionString) ImplementsEmbeddingNewParamsInputUnion()     {}
func (UnionString) ImplementsModerationNewParamsInputUnion()    {}

type UnionInt int64

func (UnionInt) ImplementsFineTuningJobHyperparametersBatchSizeUnion()                          {}
func (UnionInt) ImplementsFineTuningJobHyperparametersNEpochsUnion()                            {}
func (UnionInt) ImplementsFineTuningJobMethodDpoHyperparametersBatchSizeUnion()                 {}
func (UnionInt) ImplementsFineTuningJobMethodDpoHyperparametersNEpochsUnion()                   {}
func (UnionInt) ImplementsFineTuningJobMethodSupervisedHyperparametersBatchSizeUnion()          {}
func (UnionInt) ImplementsFineTuningJobMethodSupervisedHyperparametersNEpochsUnion()            {}
func (UnionInt) ImplementsFineTuningJobNewParamsHyperparametersBatchSizeUnion()                 {}
func (UnionInt) ImplementsFineTuningJobNewParamsHyperparametersNEpochsUnion()                   {}
func (UnionInt) ImplementsFineTuningJobNewParamsMethodDpoHyperparametersBatchSizeUnion()        {}
func (UnionInt) ImplementsFineTuningJobNewParamsMethodDpoHyperparametersNEpochsUnion()          {}
func (UnionInt) ImplementsFineTuningJobNewParamsMethodSupervisedHyperparametersBatchSizeUnion() {}
func (UnionInt) ImplementsFineTuningJobNewParamsMethodSupervisedHyperparametersNEpochsUnion()   {}

type UnionFloat float64

func (UnionFloat) ImplementsFineTuningJobHyperparametersLearningRateMultiplierUnion()          {}
func (UnionFloat) ImplementsFineTuningJobMethodDpoHyperparametersBetaUnion()                   {}
func (UnionFloat) ImplementsFineTuningJobMethodDpoHyperparametersLearningRateMultiplierUnion() {}
func (UnionFloat) ImplementsFineTuningJobMethodSupervisedHyperparametersLearningRateMultiplierUnion() {
}
func (UnionFloat) ImplementsFineTuningJobNewParamsHyperparametersLearningRateMultiplierUnion() {}
func (UnionFloat) ImplementsFineTuningJobNewParamsMethodDpoHyperparametersBetaUnion()          {}
func (UnionFloat) ImplementsFineTuningJobNewParamsMethodDpoHyperparametersLearningRateMultiplierUnion() {
}
func (UnionFloat) ImplementsFineTuningJobNewParamsMethodSupervisedHyperparametersLearningRateMultiplierUnion() {
}

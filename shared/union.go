// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package shared

type UnionString string

func (UnionString) ImplementsComparisonFilterValueUnionParam()                      {}
func (UnionString) ImplementsCompletionNewParamsPromptUnion()                       {}
func (UnionString) ImplementsCompletionNewParamsStopUnion()                         {}
func (UnionString) ImplementsChatCompletionAssistantMessageParamContentUnion()      {}
func (UnionString) ImplementsChatCompletionDeveloperMessageParamContentUnion()      {}
func (UnionString) ImplementsChatCompletionPredictionContentContentUnionParam()     {}
func (UnionString) ImplementsChatCompletionSystemMessageParamContentUnion()         {}
func (UnionString) ImplementsChatCompletionToolMessageParamContentUnion()           {}
func (UnionString) ImplementsChatCompletionUserMessageParamContentUnion()           {}
func (UnionString) ImplementsChatCompletionNewParamsStopUnion()                     {}
func (UnionString) ImplementsEmbeddingNewParamsInputUnion()                         {}
func (UnionString) ImplementsModerationNewParamsInputUnion()                        {}
func (UnionString) ImplementsVectorStoreSearchResponseAttributesUnion()             {}
func (UnionString) ImplementsVectorStoreSearchParamsQueryUnion()                    {}
func (UnionString) ImplementsVectorStoreFileAttributesUnion()                       {}
func (UnionString) ImplementsVectorStoreFileNewParamsAttributesUnion()              {}
func (UnionString) ImplementsVectorStoreFileUpdateParamsAttributesUnion()           {}
func (UnionString) ImplementsVectorStoreFileBatchNewParamsAttributesUnion()         {}
func (UnionString) ImplementsBetaThreadNewParamsMessagesContentUnion()              {}
func (UnionString) ImplementsBetaThreadNewAndRunParamsThreadMessagesContentUnion()  {}
func (UnionString) ImplementsBetaThreadRunNewParamsAdditionalMessagesContentUnion() {}
func (UnionString) ImplementsBetaThreadMessageNewParamsContentUnion()               {}

type UnionBool bool

func (UnionBool) ImplementsComparisonFilterValueUnionParam()              {}
func (UnionBool) ImplementsVectorStoreSearchResponseAttributesUnion()     {}
func (UnionBool) ImplementsVectorStoreFileAttributesUnion()               {}
func (UnionBool) ImplementsVectorStoreFileNewParamsAttributesUnion()      {}
func (UnionBool) ImplementsVectorStoreFileUpdateParamsAttributesUnion()   {}
func (UnionBool) ImplementsVectorStoreFileBatchNewParamsAttributesUnion() {}

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

func (UnionFloat) ImplementsComparisonFilterValueUnionParam()                                  {}
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
func (UnionFloat) ImplementsVectorStoreSearchResponseAttributesUnion()     {}
func (UnionFloat) ImplementsVectorStoreFileAttributesUnion()               {}
func (UnionFloat) ImplementsVectorStoreFileNewParamsAttributesUnion()      {}
func (UnionFloat) ImplementsVectorStoreFileUpdateParamsAttributesUnion()   {}
func (UnionFloat) ImplementsVectorStoreFileBatchNewParamsAttributesUnion() {}

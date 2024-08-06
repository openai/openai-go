// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package shared

type UnionString string

func (UnionString) ImplementsCompletionNewParamsPromptUnion()                       {}
func (UnionString) ImplementsCompletionNewParamsStopUnion()                         {}
func (UnionString) ImplementsChatCompletionAssistantMessageParamContentUnion()      {}
func (UnionString) ImplementsChatCompletionSystemMessageParamContentUnion()         {}
func (UnionString) ImplementsChatCompletionUserMessageParamContentUnion()           {}
func (UnionString) ImplementsChatCompletionToolMessageParamContentUnion()           {}
func (UnionString) ImplementsChatCompletionNewParamsStopUnion()                     {}
func (UnionString) ImplementsEmbeddingNewParamsInputUnion()                         {}
func (UnionString) ImplementsModerationNewParamsInputUnion()                        {}
func (UnionString) ImplementsBetaThreadNewParamsMessagesContentUnion()              {}
func (UnionString) ImplementsBetaThreadNewAndRunParamsThreadMessagesContentUnion()  {}
func (UnionString) ImplementsBetaThreadRunNewParamsAdditionalMessagesContentUnion() {}
func (UnionString) ImplementsBetaThreadMessageNewParamsContentUnion()               {}

type UnionInt int64

func (UnionInt) ImplementsFineTuningJobHyperparametersNEpochsUnion()            {}
func (UnionInt) ImplementsFineTuningJobNewParamsHyperparametersBatchSizeUnion() {}
func (UnionInt) ImplementsFineTuningJobNewParamsHyperparametersNEpochsUnion()   {}

type UnionFloat float64

func (UnionFloat) ImplementsFineTuningJobNewParamsHyperparametersLearningRateMultiplierUnion() {}

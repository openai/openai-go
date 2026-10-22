package openai_test

import openai "github.com/openai/openai-go/v3"

// Existing consumers may use the published positional field order.
var releasedChoice openai.ChatCompletionChunkChoice
var _ = openai.ChatCompletionChunkChoice{
	releasedChoice.Delta,
	releasedChoice.FinishReason,
	releasedChoice.Index,
	releasedChoice.Logprobs,
	releasedChoice.JSON,
}

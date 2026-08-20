package openai

import "github.com/openai/openai-go/v3/shared/constant"

// Helper to accumulate chunks from a stream
type ChatCompletionAccumulator struct {
	// The up-to-date accumulation of model's responses
	ChatCompletion
	choiceChatCompletionStates       []chatCompletionResponseState
	legacyChoiceChatCompletionStates []chatCompletionResponseState
	justFinished                     chatCompletionResponseState
	justFinishedByChoice             []chatCompletionResponseState
}

type FinishedChatCompletionToolCall struct {
	ChatCompletionMessageFunctionToolCallFunction
	Index int
	ID    string
}

type chatCompletionResponseState struct {
	state         chatCompletionResponseStateEnum
	choiceIndex   int
	toolCallIndex int
}

type chatCompletionResponseStateEnum int

const (
	emptyResponseState chatCompletionResponseStateEnum = iota
	contentResponseState
	refusalResponseState
	toolResponseState
	finishedResponseState
)

// maxStreamAccumulatorIndex limits zero-based protocol positions to 128 choices
// or tool calls. Stream indices are expected to be small and dense; bounding them
// prevents response data from controlling sparse allocation sizes and guarantees
// that the subsequent conversion fits in int on every Go architecture.
const maxStreamAccumulatorIndex int64 = 127

// AddChunk incorporates a chunk into the accumulation. Chunks must be added in order.
// Returns false if the chunk could not be successfully accumulated.
// Choice and tool-call indices must be between 0 and 127. For compatibility with
// providers that use -1 for a single tool call, that value is treated as index 0.
//
// The ChatCompletion field JSON does not get accumulated.
func (acc *ChatCompletionAccumulator) AddChunk(chunk ChatCompletionChunk) bool {
	if !validChatCompletionChunkIndices(chunk) {
		return false
	}
	// All conversions below consume values checked by the full preflight above.

	acc.justFinished = chatCompletionResponseState{}
	acc.justFinishedByChoice = acc.justFinishedByChoice[:0]
	if !acc.accumulateDelta(chunk) {
		return false
	}

	if len(chunk.Choices) > 0 {
		firstChoice := chunk.Choices[0]
		choiceIndex, _ := checkedStreamAccumulatorIndex(firstChoice.Index)
		acc.legacyChoiceChatCompletionStates = expandToFit(acc.legacyChoiceChatCompletionStates, choiceIndex)
		acc.justFinished = acc.legacyChoiceChatCompletionStates[choiceIndex].update(firstChoice)
	}

	for _, choice := range chunk.Choices {
		choiceIndex, _ := checkedStreamAccumulatorIndex(choice.Index)
		acc.choiceChatCompletionStates = expandToFit(acc.choiceChatCompletionStates, choiceIndex)
		justFinished := acc.choiceChatCompletionStates[choiceIndex].update(choice)
		if justFinished.state != emptyResponseState {
			acc.justFinishedByChoice = append(acc.justFinishedByChoice, justFinished)
		}
	}
	return true
}

// JustFinishedContent retrieves the chat completion content when it is known to have just been completed.
// The content is "just completed" when the last added chunk no longer contains a content
// delta. If the content is just completed, the content is returned and the boolean is true. Otherwise,
// an empty string is returned and the boolean will be false.
//
// This method preserves the legacy behavior of reporting events only for the first choice in the
// most recently added chunk. Use [ChatCompletionAccumulator.JustFinishedContentForChoice] to inspect
// another choice.
func (acc *ChatCompletionAccumulator) JustFinishedContent() (content string, ok bool) {
	if acc.justFinished.state == contentResponseState {
		return acc.Choices[acc.justFinished.choiceIndex].Message.Content, true
	}
	return "", false
}

// JustFinishedContentForChoice retrieves content that was just completed for the given choice.
func (acc *ChatCompletionAccumulator) JustFinishedContentForChoice(choiceIndex int) (content string, ok bool) {
	for _, justFinished := range acc.justFinishedByChoice {
		if justFinished.choiceIndex == choiceIndex && justFinished.state == contentResponseState {
			return acc.Choices[justFinished.choiceIndex].Message.Content, true
		}
	}
	return "", false
}

// JustFinishedRefusal retrieves the chat completion refusal when it is known to have just been completed.
// The refusal is "just completed" when the last added chunk no longer contains a refusal
// delta. If the refusal is just completed, the refusal is returned and the boolean is true. Otherwise,
// an empty string is returned and the boolean will be false.
//
// This method preserves the legacy behavior of reporting events only for the first choice in the
// most recently added chunk. Use [ChatCompletionAccumulator.JustFinishedRefusalForChoice] to inspect
// another choice.
func (acc *ChatCompletionAccumulator) JustFinishedRefusal() (refusal string, ok bool) {
	if acc.justFinished.state == refusalResponseState {
		return acc.Choices[acc.justFinished.choiceIndex].Message.Refusal, true
	}
	return "", false
}

// JustFinishedRefusalForChoice retrieves a refusal that was just completed for the given choice.
func (acc *ChatCompletionAccumulator) JustFinishedRefusalForChoice(choiceIndex int) (refusal string, ok bool) {
	for _, justFinished := range acc.justFinishedByChoice {
		if justFinished.choiceIndex == choiceIndex && justFinished.state == refusalResponseState {
			return acc.Choices[justFinished.choiceIndex].Message.Refusal, true
		}
	}
	return "", false
}

// JustFinishedToolCall retrieves a tool call when it is known to have just been completed.
// A tool call is "just completed" when the last added chunk no longer contains a tool call
// delta or contains a delta for a different tool call. If the tool call is just completed,
// a FinishedChatCompletionToolCall is returned and the boolean is true. Otherwise, an empty
// tool call is returned and the boolean will be false.
//
// You cannot rely on this with a stream that has ParallelToolCalls enabled.
//
// This method preserves the legacy behavior of reporting events only for the first choice in the
// most recently added chunk. Use [ChatCompletionAccumulator.JustFinishedToolCallForChoice] to inspect
// another choice.
func (acc *ChatCompletionAccumulator) JustFinishedToolCall() (toolcall FinishedChatCompletionToolCall, ok bool) {
	if acc.justFinished.state != toolResponseState {
		return FinishedChatCompletionToolCall{}, false
	}
	return acc.finishedToolCall(acc.justFinished), true
}

// JustFinishedToolCallForChoice retrieves a tool call that was just completed for the given choice.
//
// You cannot rely on this with a stream that has ParallelToolCalls enabled.
func (acc *ChatCompletionAccumulator) JustFinishedToolCallForChoice(choiceIndex int) (toolcall FinishedChatCompletionToolCall, ok bool) {
	for _, justFinished := range acc.justFinishedByChoice {
		if justFinished.choiceIndex == choiceIndex && justFinished.state == toolResponseState {
			return acc.finishedToolCall(justFinished), true
		}
	}
	return FinishedChatCompletionToolCall{}, false
}

func (acc *ChatCompletionAccumulator) finishedToolCall(justFinished chatCompletionResponseState) FinishedChatCompletionToolCall {
	toolCall := acc.Choices[justFinished.choiceIndex].Message.ToolCalls[justFinished.toolCallIndex]
	return FinishedChatCompletionToolCall{
		ID:    toolCall.ID,
		Index: justFinished.toolCallIndex,
		ChatCompletionMessageFunctionToolCallFunction: ChatCompletionMessageFunctionToolCallFunction{
			Name:      toolCall.Function.Name,
			Arguments: toolCall.Function.Arguments,
		},
	}
}

// Concatenates a ChatCompletionChunk onto a ChatCompletion. Returns false and
// does nothing if a mismatch is detected.
//
// Ignores the JSON field
func (cc *ChatCompletion) accumulateDelta(chunk ChatCompletionChunk) bool {
	if len(cc.ID) == 0 {
		cc.ID = chunk.ID
	} else if cc.ID != chunk.ID {
		return false
	}

	for _, delta := range chunk.Choices {
		choiceIndex, _ := checkedStreamAccumulatorIndex(delta.Index)
		cc.Choices = expandToFit(cc.Choices, choiceIndex)
		choice := &cc.Choices[choiceIndex]

		choice.Index = delta.Index
		choice.FinishReason = delta.FinishReason

		if delta.Delta.Role != "" {
			choice.Message.Role = constant.Assistant(delta.Delta.Role)
		}

		choice.Message.Content += delta.Delta.Content
		choice.Message.Refusal += delta.Delta.Refusal

		for j := range delta.Delta.ToolCalls {
			deltaTool := &delta.Delta.ToolCalls[j]
			toolIndex, _ := checkedToolCallIndex(deltaTool.Index)

			choice.Message.ToolCalls = expandToFit(choice.Message.ToolCalls, toolIndex)
			tool := &choice.Message.ToolCalls[toolIndex]

			if deltaTool.ID != "" {
				tool.ID = deltaTool.ID
			}
			if deltaTool.Type != "" {
				tool.Type = deltaTool.Type
			}
			tool.Function.Name += deltaTool.Function.Name
			tool.Function.Arguments += deltaTool.Function.Arguments
		}

		choice.Logprobs.Content = append(choice.Logprobs.Content, delta.Logprobs.Content...)
		choice.Logprobs.Refusal = append(choice.Logprobs.Refusal, delta.Logprobs.Refusal...)
	}

	cc.Usage.CompletionTokens += chunk.Usage.CompletionTokens
	cc.Usage.PromptTokens += chunk.Usage.PromptTokens
	cc.Usage.TotalTokens += chunk.Usage.TotalTokens

	cc.Usage.CompletionTokensDetails.AcceptedPredictionTokens += chunk.Usage.CompletionTokensDetails.AcceptedPredictionTokens
	cc.Usage.CompletionTokensDetails.AudioTokens += chunk.Usage.CompletionTokensDetails.AudioTokens
	cc.Usage.CompletionTokensDetails.ReasoningTokens += chunk.Usage.CompletionTokensDetails.ReasoningTokens
	cc.Usage.CompletionTokensDetails.RejectedPredictionTokens += chunk.Usage.CompletionTokensDetails.RejectedPredictionTokens

	cc.Usage.PromptTokensDetails.AudioTokens += chunk.Usage.PromptTokensDetails.AudioTokens
	cc.Usage.PromptTokensDetails.CachedTokens += chunk.Usage.PromptTokensDetails.CachedTokens

	cc.Model = chunk.Model
	cc.Created = chunk.Created
	cc.SystemFingerprint = chunk.SystemFingerprint
	cc.ServiceTier = ChatCompletionServiceTier(chunk.ServiceTier)
	if chunk.Object == chunk.Object.Default() {
		cc.Object = cc.Object.Default()
	}

	return true
}

// Updates the internal response state and returns the previous state if
// the state changed. This ensures that JustFinished events only fire once.
func (prev *chatCompletionResponseState) update(choice ChatCompletionChunkChoice) (justFinished chatCompletionResponseState) {
	delta := choice.Delta
	choiceIndex, _ := checkedStreamAccumulatorIndex(choice.Index)
	new := chatCompletionResponseState{choiceIndex: choiceIndex}
	switch {
	case len(delta.ToolCalls) > 0 && delta.Content == "":
		new.state = toolResponseState
		new.toolCallIndex, _ = checkedToolCallIndex(delta.ToolCalls[0].Index)
	case delta.JSON.Content.Valid():
		new.state = contentResponseState
	case delta.JSON.Refusal.Valid():
		new.state = refusalResponseState
	case len(delta.ToolCalls) > 0:
		new.state = toolResponseState
		new.toolCallIndex, _ = checkedToolCallIndex(delta.ToolCalls[0].Index)
	default:
		new.state = finishedResponseState
	}

	if *prev != new {
		justFinished = *prev
	}
	*prev = new

	return
}

func validChatCompletionChunkIndices(chunk ChatCompletionChunk) bool {
	for _, choice := range chunk.Choices {
		if _, ok := checkedStreamAccumulatorIndex(choice.Index); !ok {
			return false
		}
		for _, toolCall := range choice.Delta.ToolCalls {
			if _, ok := checkedToolCallIndex(toolCall.Index); !ok {
				return false
			}
		}
	}
	return true
}

func checkedStreamAccumulatorIndex(index int64) (int, bool) {
	if index < 0 || index > maxStreamAccumulatorIndex {
		return 0, false
	}
	return int(index), true
}

// checkedToolCallIndex handles providers like AWS Bedrock that return -1 for a
// single tool call while rejecting all other values outside the shared bound.
func checkedToolCallIndex(index int64) (int, bool) {
	if index == -1 {
		return 0, true
	}
	return checkedStreamAccumulatorIndex(index)
}

func expandToFit[T any](slice []T, index int) []T {
	if index < len(slice) {
		return slice
	}
	if index < cap(slice) {
		return slice[:index+1]
	}
	newSlice := make([]T, index+1)
	copy(newSlice, slice)
	return newSlice
}

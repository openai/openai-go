package openai

import (
	"unsafe"
	"weak"

	"github.com/openai/openai-go/v3/shared/constant"
)

const (
	chatCompletionAccumulatorInlineProjectionSlots = 1_024
	maxChatCompletionAccumulatorInt                = int(^uint(0) >> 1)
	chatCompletionAccumulatorChoiceWork            = 16
	chatCompletionAccumulatorToolWork              = 8
	// Conservatively covers a non-empty map header and its first backing bucket.
	chatCompletionAccumulatorMapOverheadBytes = 512
)

// Helper to accumulate chunks from a stream
type ChatCompletionAccumulator struct {
	// The up-to-date accumulation of model's responses
	ChatCompletion
	choiceChatCompletionStates       [maxStreamAccumulatorChoiceIndex + 1]chatCompletionResponseStateCode
	legacyChoiceChatCompletionStates [maxStreamAccumulatorChoiceIndex + 1]chatCompletionResponseStateCode
	justFinished                     chatCompletionResponseState
	justFinishedByChoice             [maxStreamAccumulatorChoiceIndex + 1]chatCompletionChoiceEvents
	stringState                      chatCompletionAccumulatorStringState
	logprobState                     chatCompletionAccumulatorLogprobState
	textBytes                        int
	textReconciliationWork           int
	logprobBytes                     int
	reconciliationWork               int
	privateStateOwner                weak.Pointer[ChatCompletionAccumulator]
	privateStateInitialized          bool
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

type chatCompletionResponseStateCode struct {
	state         chatCompletionResponseStateEnum
	toolCallIndex int
}

type chatCompletionChoiceEvents struct {
	states        uint8
	toolCallIndex int
}

type chatCompletionAccumulatorStringState struct {
	choices           []*chatCompletionChoiceStringState
	activeChoices     []int
	activeTools       int
	id                string
	model             string
	systemFingerprint string
	object            string
	serviceTier       string
}

type chatCompletionChoiceStringState struct {
	content         chatCompletionString
	refusal         chatCompletionString
	toolCalls       []weak.Pointer[chatCompletionToolCallStringState]
	toolOwners      *chatCompletionToolCallOwner
	toolOwnerDepth  uint8
	activeToolCalls []int
	reconcileCursor int
	finishReason    string
	role            string
}

type chatCompletionToolCallStringState struct {
	name      chatCompletionString
	arguments chatCompletionString
	id        string
	typeName  string
}

type chatCompletionString struct {
	buffer    []byte
	published string
	shared    bool
}

type chatCompletionResponseStateEnum int

const (
	emptyResponseState chatCompletionResponseStateEnum = iota
	contentResponseState
	refusalResponseState
	toolResponseState
	finishedResponseState
)

// maxStreamAccumulatorChoiceIndex limits zero-based protocol positions to the
// maximum 128 choices supported by the API. The bound also guarantees that the
// subsequent conversion fits in int on every Go architecture.
const maxStreamAccumulatorChoiceIndex = 127

// maxStreamAccumulatorToolCallGrowth preserves support for sparse provider
// indices while limiting one received tool-call entry to 128 newly allocated
// positions. The limit applies to growth, not the total number of tool calls,
// so dense streams can exceed 128 calls without allowing one sparse index to
// trigger an unbounded allocation.
const maxStreamAccumulatorToolCallGrowth = 128

// AddChunk incorporates a chunk into the accumulation. Chunks must be added in order.
// Returns false if the chunk could not be successfully accumulated. To bound work and
// memory for untrusted streams, an accumulator accepts choice indices from 0 through
// 127 and bounds the sparse growth caused by each individual tool-call index.
// Accumulated chunks, dense tool calls, content, refusal, tool-function text,
// metadata, and log probabilities retain their existing unlimited-size contract;
// text and logprob buffers grow geometrically with checked representable-int
// accounting.
// A tool-call index may grow its choice by at most 128 positions; dense sequences
// may contain more than 128 calls.
// For compatibility with providers that use -1 for a single tool call, that value is
// treated as index 0. A rejected chunk does not modify accumulated response data or
// resource budgets, but it clears any prior JustFinished event.
//
// While accumulation is in progress, callers may replace or clear accumulated
// top-level strings and string fields of choices and tool calls that the stream
// has populated, or replace an entire Content or Refusal logprob slice. ID may
// be cleared but cannot be replaced with a different nonempty value. Writing into
// sparse placeholder slots or mutating fields inside an accumulated logprob
// element is unsupported; copy the result after accumulation before making those
// edits.
//
// The ChatCompletion field JSON does not get accumulated.
func (acc *ChatCompletionAccumulator) AddChunk(chunk ChatCompletionChunk) bool {
	acc.justFinished = chatCompletionResponseState{}
	clear(acc.justFinishedByChoice[:])
	if !acc.preflightChunk(&chunk) {
		return false
	}

	acc.accumulateDelta(&chunk)
	acc.logprobState.acceptChunk(&acc.ChatCompletion, &chunk)
	if len(chunk.Choices) > 0 {
		firstChoice := chunk.Choices[0]
		choiceIndex, _ := checkedStreamAccumulatorChoiceIndex(firstChoice.Index)
		toolCallCount := len(acc.Choices[choiceIndex].Message.ToolCalls)
		acc.justFinished = acc.legacyChoiceChatCompletionStates[choiceIndex].update(firstChoice, toolCallCount)
	}

	for _, choice := range chunk.Choices {
		choiceIndex, _ := checkedStreamAccumulatorChoiceIndex(choice.Index)
		toolCallCount := len(acc.Choices[choiceIndex].Message.ToolCalls)
		justFinished := acc.choiceChatCompletionStates[choiceIndex].update(choice, toolCallCount)
		if justFinished.state != emptyResponseState {
			events := &acc.justFinishedByChoice[choiceIndex]
			flag := uint8(1) << justFinished.state
			if justFinished.state == toolResponseState && events.states&flag == 0 {
				events.toolCallIndex = justFinished.toolCallIndex
			}
			events.states |= flag
		}
	}
	return true
}

func (acc *ChatCompletionAccumulator) preflightChunk(chunk *ChatCompletionChunk) bool {
	canonicalID := acc.stringState.id
	if canonicalID != "" && canonicalID != chunk.ID {
		return false
	}
	if acc.ID != "" && !accumulatorStringUsesPublishedBacking(acc.ID, canonicalID) && acc.ID != chunk.ID {
		return false
	}
	if !acc.validChatCompletionChunkIndices(*chunk) {
		return false
	}
	projectedReconciliationWork, ok := acc.projectReconciliationWork(chunk)
	if !ok {
		return false
	}
	projectedTextBytes, ok := addChatCompletionChunkTextBytes(acc.textBytes, chunk)
	if !ok {
		return false
	}
	projectedTextReconciliationWork := acc.textReconciliationWork
	liveTextBytes, ok := acc.addChatCompletionTextBytes(0, &projectedTextReconciliationWork, chunk)
	if !ok {
		return false
	}
	if _, ok = addChatCompletionChunkTextBytes(liveTextBytes, chunk); !ok {
		return false
	}
	if !acc.chatCompletionMetadataWithinLimit(chunk, &projectedReconciliationWork) {
		return false
	}
	var logprobPlan chatCompletionLogprobReconcilePlan
	hasLogprobPlan, ok := acc.planLogprobReconciliation(&logprobPlan, &projectedReconciliationWork, chunk)
	if !ok {
		return false
	}
	logprobState := &acc.logprobState
	if hasLogprobPlan {
		logprobState = &logprobPlan.state
	}
	projectedLogprobBytes, ok := addChatCompletionChunkLogprobBytes(acc.logprobBytes, chunk)
	if !ok {
		return false
	}
	if !logprobState.chunkWithinLimit(chunk) {
		return false
	}

	acc.detachPrivateStateForCopy()
	if hasLogprobPlan {
		acc.applyLogprobReconciliation(&logprobPlan)
	}
	acc.reconcilePublicState(chunk)
	acc.textBytes = projectedTextBytes
	acc.textReconciliationWork = projectedTextReconciliationWork
	acc.logprobBytes = projectedLogprobBytes
	acc.reconciliationWork = projectedReconciliationWork
	acc.claimPrivateStateOwnership()
	return true
}

func (acc *ChatCompletionAccumulator) projectReconciliationWork(chunk *ChatCompletionChunk) (int, bool) {
	// Preflight and commit make seven passes over populated choices, including
	// the structural-copy projection, plus one
	// structural projection pass when the chunk has choices. Populated tools are
	// visited by text, metadata, and commit reconciliation. Charging those exact
	// passes keeps sparse detection work bounded without scanning placeholder slots.
	// Per-chunk choice and tool units conservatively cover the remaining fixed
	// validation, projection, accumulation, and state-update passes over new input.
	choicePasses := 7
	if len(chunk.Choices) > 0 {
		choicePasses++
	}
	work := len(acc.stringState.activeChoices) * choicePasses
	for _, choiceIndex := range acc.stringState.activeChoices {
		if choiceIndex >= len(acc.Choices) {
			continue
		}
		choiceState := acc.stringState.choices[choiceIndex]
		if !choiceState.visitReconciledTools(choiceIndex, acc.Choices[choiceIndex].Message.ToolCalls, chunk,
			func(_ int) bool { return addAccumulatorReconciliationWork(&work, 3) }) {
			return 0, false
		}
	}
	privateCopyWork := acc.privateStateCopyWork()
	if privateCopyWork > maxChatCompletionAccumulatorInt-work {
		return 0, false
	}
	work += privateCopyWork
	structuralCopyWork := acc.structuralReconciliationCopyWork()
	if structuralCopyWork > maxChatCompletionAccumulatorInt-work {
		return 0, false
	}
	work += structuralCopyWork
	if work > maxChatCompletionAccumulatorInt-acc.reconciliationWork {
		return 0, false
	}
	projected := acc.reconciliationWork + work
	if len(chunk.Choices) > (maxChatCompletionAccumulatorInt-projected)/chatCompletionAccumulatorChoiceWork {
		return 0, false
	}
	projected += len(chunk.Choices) * chatCompletionAccumulatorChoiceWork
	for i := range chunk.Choices {
		toolCount := len(chunk.Choices[i].Delta.ToolCalls)
		if toolCount > (maxChatCompletionAccumulatorInt-projected)/chatCompletionAccumulatorToolWork {
			return 0, false
		}
		projected += toolCount * chatCompletionAccumulatorToolWork
	}
	return projected, true
}

func (acc *ChatCompletionAccumulator) structuralReconciliationCopyWork() int {
	completion := &acc.ChatCompletion
	stringState := &acc.stringState
	choiceCount := len(completion.Choices)
	work := 0
	if choiceCount < len(stringState.choices) {
		// Public and private choice slices may each be copied before growing again,
		// while active indices are visited once to count and once to retain them.
		work += 4*choiceCount + 2*len(stringState.activeChoices)
	}
	for _, i := range stringState.activeChoices {
		if i >= choiceCount {
			continue
		}
		choiceState := stringState.choices[i]
		toolCallCount := len(completion.Choices[i].Message.ToolCalls)
		if toolCallCount >= len(choiceState.toolCalls) {
			continue
		}
		// Truncation copies the visible public and private prefixes. Re-activation
		// can then grow both exact-capacity slices, copying those prefixes again.
		// The bounded outer choice table may also detach at both stages.
		work += 4*toolCallCount + 2*len(choiceState.activeToolCalls) + 2*len(stringState.choices)
	}
	return work
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
	if choiceIndex >= 0 && choiceIndex < len(acc.justFinishedByChoice) &&
		acc.justFinishedByChoice[choiceIndex].states&(uint8(1)<<contentResponseState) != 0 {
		return acc.Choices[choiceIndex].Message.Content, true
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
	if choiceIndex >= 0 && choiceIndex < len(acc.justFinishedByChoice) &&
		acc.justFinishedByChoice[choiceIndex].states&(uint8(1)<<refusalResponseState) != 0 {
		return acc.Choices[choiceIndex].Message.Refusal, true
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
	if choiceIndex >= 0 && choiceIndex < len(acc.justFinishedByChoice) &&
		acc.justFinishedByChoice[choiceIndex].states&(uint8(1)<<toolResponseState) != 0 {
		return acc.finishedToolCall(chatCompletionResponseState{
			state:         toolResponseState,
			choiceIndex:   choiceIndex,
			toolCallIndex: acc.justFinishedByChoice[choiceIndex].toolCallIndex,
		}), true
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

// Concatenates a preflighted ChatCompletionChunk onto a ChatCompletion.
// Ignores the JSON field.
func (acc *ChatCompletionAccumulator) accumulateDelta(chunk *ChatCompletionChunk) {
	cc := &acc.ChatCompletion
	acc.stringState.detachToolActivationState(chunk)
	if len(cc.ID) == 0 {
		assignAccumulatorString(&acc.stringState.id, &cc.ID, chunk.ID)
	}

	for _, delta := range chunk.Choices {
		choiceIndex, _ := checkedStreamAccumulatorChoiceIndex(delta.Index)
		cc.Choices = expandToFit(cc.Choices, choiceIndex)
		choice := &cc.Choices[choiceIndex]
		choiceStrings := acc.stringState.choice(choiceIndex)
		if accumulatorStringWillGrow(&choiceStrings.content, delta.Delta.Content) ||
			accumulatorStringWillGrow(&choiceStrings.refusal, delta.Delta.Refusal) ||
			accumulatorMetadataWillChange(choiceStrings.finishReason, delta.FinishReason) ||
			accumulatorMetadataWillChange(choiceStrings.role, delta.Delta.Role) {
			choiceStrings = acc.stringState.detachChoiceState(choiceIndex)
		}

		choice.Index = delta.Index
		assignAccumulatorString(&choiceStrings.finishReason, &choice.FinishReason, delta.FinishReason)

		if delta.Delta.Role != "" {
			assignAccumulatorString(&choiceStrings.role, &choice.Message.Role, constant.Assistant(delta.Delta.Role))
		}

		choiceStrings.content.append(&choice.Message.Content, delta.Delta.Content)
		choiceStrings.refusal.append(&choice.Message.Refusal, delta.Delta.Refusal)

		toolCallCount := len(choice.Message.ToolCalls)
		for _, deltaTool := range delta.Delta.ToolCalls {
			toolIndex, _ := checkedToolCallIndex(deltaTool.Index, toolCallCount)
			if toolIndex >= toolCallCount {
				toolCallCount = toolIndex + 1
			}
		}
		if toolCallCount > len(choice.Message.ToolCalls) {
			choice.Message.ToolCalls = expandToFit(choice.Message.ToolCalls, toolCallCount-1)
		}

		previousToolCallCount := len(choiceStrings.toolCalls)
		for j := range delta.Delta.ToolCalls {
			deltaTool := &delta.Delta.ToolCalls[j]
			toolIndex, _ := checkedToolCallIndex(deltaTool.Index, len(choice.Message.ToolCalls))
			tool := &choice.Message.ToolCalls[toolIndex]
			toolStrings := choiceStrings.toolCall(toolIndex, &acc.stringState.activeTools)
			if toolIndex < previousToolCallCount &&
				(accumulatorStringWillGrow(&toolStrings.name, deltaTool.Function.Name) ||
					accumulatorStringWillGrow(&toolStrings.arguments, deltaTool.Function.Arguments) ||
					accumulatorMetadataWillChange(toolStrings.id, deltaTool.ID) ||
					accumulatorMetadataWillChange(toolStrings.typeName, deltaTool.Type)) {
				choiceStrings = acc.stringState.detachChoiceState(choiceIndex)
				toolStrings = choiceStrings.detachToolCallState(toolIndex)
			}

			if deltaTool.ID != "" {
				assignAccumulatorString(&toolStrings.id, &tool.ID, deltaTool.ID)
			}
			if deltaTool.Type != "" {
				assignAccumulatorString(&toolStrings.typeName, &tool.Type, deltaTool.Type)
			}
			toolStrings.name.append(&tool.Function.Name, deltaTool.Function.Name)
			toolStrings.arguments.append(&tool.Function.Arguments, deltaTool.Function.Arguments)
		}

		choice.Logprobs.Content = appendChatCompletionLogprobs(choice.Logprobs.Content, delta.Logprobs.Content)
		choice.Logprobs.Refusal = appendChatCompletionLogprobs(choice.Logprobs.Refusal, delta.Logprobs.Refusal)
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

	assignAccumulatorString(&acc.stringState.model, &cc.Model, chunk.Model)
	cc.Created = chunk.Created
	assignAccumulatorString(&acc.stringState.systemFingerprint, &cc.SystemFingerprint, chunk.SystemFingerprint)
	assignAccumulatorString(&acc.stringState.serviceTier, &cc.ServiceTier, ChatCompletionServiceTier(chunk.ServiceTier))
	if chunk.Object == chunk.Object.Default() {
		assignAccumulatorString(&acc.stringState.object, &cc.Object, cc.Object.Default())
	}
}

func (acc *chatCompletionAccumulatorStringState) detachToolActivationState(chunk *ChatCompletionChunk) {
	var detachChoice [maxStreamAccumulatorChoiceIndex + 1]bool
	var detachToolCalls [maxStreamAccumulatorChoiceIndex + 1]bool
	hasDetach := false
	for i := range chunk.Choices {
		choice := &chunk.Choices[i]
		choiceIndex := int(choice.Index)
		if choiceIndex >= len(acc.choices) || acc.choices[choiceIndex] == nil {
			continue
		}
		choiceState := acc.choices[choiceIndex]
		for j := range choice.Delta.ToolCalls {
			toolIndex := preflightedToolCallIndex(choice.Delta.ToolCalls[j].Index)
			if toolIndex < len(choiceState.toolCalls) && choiceState.toolCallState(toolIndex) != nil {
				continue
			}
			detachChoice[choiceIndex] = true
			detachToolCalls[choiceIndex] = detachToolCalls[choiceIndex] || toolIndex < len(choiceState.toolCalls)
			hasDetach = true
		}
	}
	if !hasDetach {
		return
	}

	acc.choices = cloneAccumulatorSlice(acc.choices)
	for choiceIndex, detach := range detachChoice {
		if !detach {
			continue
		}
		choiceState := *acc.choices[choiceIndex]
		if detachToolCalls[choiceIndex] {
			choiceState.toolCalls = cloneAccumulatorSlice(choiceState.toolCalls)
		}
		acc.choices[choiceIndex] = &choiceState
	}
}

func (acc *chatCompletionAccumulatorStringState) choice(index int) *chatCompletionChoiceStringState {
	previousCapacity := cap(acc.choices)
	acc.choices = expandToFit(acc.choices, index)
	if acc.choices[index] == nil {
		if index < previousCapacity {
			acc.choices = cloneAccumulatorSlice(acc.choices)
		}
		acc.choices[index] = &chatCompletionChoiceStringState{}
		acc.activeChoices = append(acc.activeChoices, index)
	}
	return acc.choices[index]
}

func (acc *chatCompletionAccumulatorStringState) detachChoiceState(index int) *chatCompletionChoiceStringState {
	acc.choices = cloneAccumulatorSlice(acc.choices)
	state := *acc.choices[index]
	acc.choices[index] = &state
	return &state
}

func (acc *ChatCompletionAccumulator) reconcilePublicState(chunk *ChatCompletionChunk) {
	completion := &acc.ChatCompletion
	stringState := &acc.stringState
	reconcileAccumulatorString(&stringState.id, &completion.ID)
	reconcileAccumulatorString(&stringState.model, &completion.Model)
	reconcileAccumulatorString(&stringState.systemFingerprint, &completion.SystemFingerprint)
	reconcileAccumulatorString(&stringState.object, &completion.Object)
	reconcileAccumulatorString(&stringState.serviceTier, &completion.ServiceTier)

	previousChoices := stringState.choices
	previousChoiceCount := len(previousChoices)
	completion.Choices = detachTruncatedTail(completion.Choices, previousChoiceCount)
	choiceCount := len(completion.Choices)
	choicesDetached := choiceCount < previousChoiceCount
	if choicesDetached {
		stringState.choices = cloneAccumulatorSlice(previousChoices[:choiceCount])
	}
	for _, i := range stringState.activeChoices {
		if i >= choiceCount {
			continue
		}
		message := &completion.Choices[i].Message
		previousToolCallCount := len(previousChoices[i].toolCalls)
		message.ToolCalls = detachTruncatedTail(message.ToolCalls, previousToolCallCount)
	}

	previousActiveChoices := stringState.activeChoices
	stringState.activeChoices = retainAccumulatorIndices(previousActiveChoices, choiceCount)
	for _, i := range previousActiveChoices {
		if i >= choiceCount {
			stringState.activeTools -= len(previousChoices[i].activeToolCalls)
			continue
		}
		choiceState := stringState.choices[i]

		message := &completion.Choices[i].Message
		toolCallsTruncated := len(message.ToolCalls) < len(choiceState.toolCalls)
		if toolCallsTruncated {
			if !choicesDetached {
				stringState.choices = cloneAccumulatorSlice(stringState.choices)
				choicesDetached = true
			}
			detachedChoiceState := *choiceState
			detachedChoiceState.toolCalls = cloneAccumulatorSlice(choiceState.toolCalls[:len(message.ToolCalls)])
			detachedChoiceState.activeToolCalls = cloneAccumulatorSlice(
				retainAccumulatorIndices(choiceState.activeToolCalls, len(message.ToolCalls)),
			)
			detachedChoiceState.rebuildToolCallOwnership()
			stringState.activeTools -= len(choiceState.activeToolCalls) - len(detachedChoiceState.activeToolCalls)
			stringState.choices[i] = &detachedChoiceState
			choiceState = &detachedChoiceState
		}
		if accumulatorStringNeedsReconciliationDetach(&choiceState.content, message.Content) ||
			accumulatorStringNeedsReconciliationDetach(&choiceState.refusal, message.Refusal) ||
			accumulatorMetadataNeedsReconciliationDetach(choiceState.finishReason, completion.Choices[i].FinishReason) ||
			accumulatorMetadataNeedsReconciliationDetach(choiceState.role, message.Role) {
			choiceState = stringState.detachChoiceState(i)
			choicesDetached = true
		}
		choiceState.content.reconcile(&message.Content)
		choiceState.refusal.reconcile(&message.Refusal)
		reconcileAccumulatorString(&choiceState.finishReason, &completion.Choices[i].FinishReason)
		reconcileAccumulatorString(&choiceState.role, &message.Role)

		if toolCallsTruncated {
			invalidateRemovedToolCallState(&acc.legacyChoiceChatCompletionStates, i, len(message.ToolCalls))
			invalidateRemovedToolCallState(&acc.choiceChatCompletionStates, i, len(message.ToolCalls))
		}
		choiceState.visitReconciledTools(i, message.ToolCalls, chunk, func(j int) bool {
			toolCallState := choiceState.toolCallState(j)
			function := &message.ToolCalls[j].Function
			if accumulatorStringNeedsReconciliationDetach(&toolCallState.name, function.Name) ||
				accumulatorStringNeedsReconciliationDetach(&toolCallState.arguments, function.Arguments) ||
				accumulatorMetadataNeedsReconciliationDetach(toolCallState.id, message.ToolCalls[j].ID) ||
				accumulatorMetadataNeedsReconciliationDetach(toolCallState.typeName, message.ToolCalls[j].Type) {
				choiceState = stringState.detachChoiceState(i)
				choicesDetached = true
				toolCallState = choiceState.detachToolCallState(j)
			}
			reconcileAccumulatorString(&toolCallState.id, &message.ToolCalls[j].ID)
			reconcileAccumulatorString(&toolCallState.typeName, &message.ToolCalls[j].Type)
			toolCallState.name.reconcile(&function.Name)
			toolCallState.arguments.reconcile(&function.Arguments)
			return true
		})
		if len(choiceState.activeToolCalls) > 0 {
			choiceState.reconcileCursor =
				(choiceState.reconcileCursor + choiceState.reconciliationSweepCount(i, chunk)) % len(choiceState.activeToolCalls)
		}
	}
	truncateResponseStates(&acc.legacyChoiceChatCompletionStates, choiceCount)
	truncateResponseStates(&acc.choiceChatCompletionStates, choiceCount)
}

func retainAccumulatorIndices(indices []int, limit int) []int {
	retainedCount := 0
	for _, index := range indices {
		if index < limit {
			retainedCount++
		}
	}
	if retainedCount == len(indices) {
		return indices
	}
	retained := make([]int, 0, retainedCount)
	for _, index := range indices {
		if index < limit {
			retained = append(retained, index)
		}
	}
	return retained
}

func invalidateRemovedToolCallState(states *[maxStreamAccumulatorChoiceIndex + 1]chatCompletionResponseStateCode, choiceIndex int, toolCallCount int) {
	state := states[choiceIndex].decode(choiceIndex)
	if state.state == toolResponseState && state.toolCallIndex >= toolCallCount {
		states[choiceIndex] = chatCompletionResponseStateCode{}
	}
}

func truncateResponseStates(states *[maxStreamAccumulatorChoiceIndex + 1]chatCompletionResponseStateCode, choiceCount int) {
	choiceCount = min(choiceCount, len(states))
	clear(states[choiceCount:])
}

func (choice *chatCompletionChoiceStringState) toolCall(index int, activeTools *int) *chatCompletionToolCallStringState {
	choice.toolCalls = expandToFit(choice.toolCalls, index)
	state := choice.toolCallState(index)
	if state == nil {
		state = &chatCompletionToolCallStringState{}
		choice.setToolCallOwner(index, state)
		choice.toolCalls[index] = weak.Make(state)
		choice.activeToolCalls = append(choice.activeToolCalls, index)
		(*activeTools)++
	}
	return state
}

func (choice *chatCompletionChoiceStringState) toolCallState(index int) *chatCompletionToolCallStringState {
	return chatCompletionToolCallOwnerAt(choice.toolOwners, choice.toolOwnerDepth, index)
}

func (choice *chatCompletionChoiceStringState) detachToolCallState(index int) *chatCompletionToolCallStringState {
	state := *choice.toolCallState(index)
	choice.setToolCallOwner(index, &state)
	choice.toolCalls[index] = weak.Make(&state)
	return &state
}

func (choice *chatCompletionChoiceStringState) rebuildToolCallOwnership() {
	previousOwners, previousDepth := choice.toolOwners, choice.toolOwnerDepth
	choice.toolOwners, choice.toolOwnerDepth = nil, 0
	for _, index := range choice.activeToolCalls {
		state := *chatCompletionToolCallOwnerAt(previousOwners, previousDepth, index)
		choice.setToolCallOwner(index, &state)
		choice.toolCalls[index] = weak.Make(&state)
	}
}

func (acc *chatCompletionString) append(current *string, fragment string) {
	if fragment == "" {
		return
	}
	acc.reconcile(current)
	required := len(acc.buffer) + len(fragment)
	if acc.shared || required > cap(acc.buffer) {
		capacity := projectedAccumulatorTextCapacity(cap(acc.buffer), required)
		grown := make([]byte, len(acc.buffer), capacity)
		copy(grown, acc.buffer)
		acc.buffer = grown
		acc.shared = false
	}
	acc.buffer = append(acc.buffer, fragment...)
	acc.published = accumulatorBufferString(acc.buffer)
	*current = acc.published
}

func accumulatorStringWillGrow(state *chatCompletionString, fragment string) bool {
	return fragment != "" && (state.shared || len(fragment) > cap(state.buffer)-len(state.buffer))
}

func accumulatorMetadataWillChange[T ~string](published string, next T) bool {
	value := string(next)
	return value != "" && value != published
}

func accumulatorStringNeedsReconciliationDetach(state *chatCompletionString, current string) bool {
	return current != state.published && !accumulatorStringUsesPublishedBacking(current, state.published)
}

func accumulatorMetadataNeedsReconciliationDetach[T ~string](published string, current T) bool {
	return string(current) != published && !accumulatorStringUsesPublishedBacking(string(current), published)
}

func (acc *chatCompletionString) reconcile(current *string) {
	if accumulatorStringUsesPublishedBacking(*current, acc.published) {
		*current = acc.published
		return
	}
	if *current == acc.published {
		*current = acc.published
		return
	}

	replacement := *current
	acc.published = ""
	acc.buffer = append(make([]byte, 0, len(replacement)), replacement...)
	acc.shared = false
	acc.published = accumulatorBufferString(acc.buffer)
	*current = acc.published
}

func projectedAccumulatorTextCapacity(current int, required int) int {
	capacity := max(1, current)
	for capacity < required {
		if capacity > maxChatCompletionAccumulatorInt/2 {
			return required
		}
		capacity *= 2
	}
	return capacity
}

func accumulatorBufferString(buffer []byte) string {
	if len(buffer) == 0 {
		return ""
	}
	// The accumulator only appends beyond the published length; it never mutates
	// bytes visible through an earlier string value.
	return unsafe.String(unsafe.SliceData(buffer), len(buffer))
}

func (acc *ChatCompletionAccumulator) addChatCompletionTextBytes(total int, work *int, chunk *ChatCompletionChunk) (int, bool) {
	if !acc.addCopiedTextBufferWork(work, chunk) {
		return 0, false
	}
	completion := &acc.ChatCompletion
	var appends *chatCompletionTextAppendProjection
	for _, i := range acc.stringState.activeChoices {
		if i >= len(completion.Choices) {
			continue
		}
		message := &completion.Choices[i].Message
		choiceState := acc.stringState.choices[i]
		contentAppend, refusalAppend := false, false
		if !accumulatorStringUsesPublishedBacking(message.Content, choiceState.content.published) ||
			!accumulatorStringUsesPublishedBacking(message.Refusal, choiceState.refusal.published) {
			appends = ensureChatCompletionTextAppendProjection(appends, chunk)
			contentAppend = appends.choiceContent(i)
			refusalAppend = appends.choiceRefusal(i)
		}
		if !addAccumulatorTextBytes(&total, message.Content) ||
			!addAccumulatorTextBytes(&total, message.Refusal) ||
			!addAccumulatorBufferReconciliationWork(work, message.Content, &choiceState.content, contentAppend) ||
			!addAccumulatorBufferReconciliationWork(work, message.Refusal, &choiceState.refusal, refusalAppend) {
			return 0, false
		}
		if !choiceState.visitReconciledTools(i, message.ToolCalls, chunk, func(j int) bool {
			function := &message.ToolCalls[j].Function
			toolState := choiceState.toolCallState(j)
			nameAppend, argumentsAppend := false, false
			if !accumulatorStringUsesPublishedBacking(function.Name, toolState.name.published) ||
				!accumulatorStringUsesPublishedBacking(function.Arguments, toolState.arguments.published) {
				appends = ensureChatCompletionTextAppendProjection(appends, chunk)
				nameAppend = appends.toolName(i, j)
				argumentsAppend = appends.toolArguments(i, j)
			}
			if !addAccumulatorTextBytes(&total, function.Name) ||
				!addAccumulatorTextBytes(&total, function.Arguments) ||
				!addAccumulatorBufferReconciliationWork(work, function.Name, &toolState.name, nameAppend) ||
				!addAccumulatorBufferReconciliationWork(work, function.Arguments, &toolState.arguments, argumentsAppend) {
				return false
			}
			return true
		}) {
			return 0, false
		}
	}
	return total, true
}

func addAccumulatorBufferReconciliationWork(work *int, current string, state *chatCompletionString, mayAppend bool) bool {
	if accumulatorStringUsesPublishedBacking(current, state.published) {
		return true
	}
	passes := 2
	if mayAppend && len(current) > 0 {
		// A changed public value is committed into exact-capacity backing. The
		// first same-chunk append must therefore copy that prefix while growing.
		passes++
	}
	for range passes {
		if !addAccumulatorTextBytes(work, current) {
			return false
		}
	}
	return true
}

func addChatCompletionChunkTextBytes(total int, chunk *ChatCompletionChunk) (int, bool) {
	for i := range chunk.Choices {
		delta := &chunk.Choices[i].Delta
		if !addAccumulatorTextBytes(&total, delta.Content) ||
			!addAccumulatorTextBytes(&total, delta.Refusal) {
			return 0, false
		}
		for j := range delta.ToolCalls {
			function := &delta.ToolCalls[j].Function
			if !addAccumulatorTextBytes(&total, function.Name) ||
				!addAccumulatorTextBytes(&total, function.Arguments) {
				return 0, false
			}
		}
	}
	return total, true
}

func (acc *ChatCompletionAccumulator) addChatCompletionMetadataBytes(total int, work *int, chunk *ChatCompletionChunk) (int, bool) {
	completion := &acc.ChatCompletion
	if !addAccumulatorMetadataBytes(&total, completion.ID) ||
		!addAccumulatorMetadataBytes(&total, completion.Model) ||
		!addAccumulatorMetadataBytes(&total, completion.SystemFingerprint) ||
		!addAccumulatorMetadataBytes(&total, string(completion.Object)) ||
		!addAccumulatorMetadataBytes(&total, string(completion.ServiceTier)) ||
		!addAccumulatorStringReconciliationWork(work, completion.ID, acc.stringState.id) ||
		!addAccumulatorStringReconciliationWork(work, completion.Model, acc.stringState.model) ||
		!addAccumulatorStringReconciliationWork(work, completion.SystemFingerprint, acc.stringState.systemFingerprint) ||
		!addAccumulatorStringReconciliationWork(work, completion.Object, acc.stringState.object) ||
		!addAccumulatorStringReconciliationWork(work, completion.ServiceTier, acc.stringState.serviceTier) {
		return 0, false
	}
	for _, i := range acc.stringState.activeChoices {
		if i >= len(completion.Choices) {
			continue
		}
		choice := &completion.Choices[i]
		choiceState := acc.stringState.choices[i]
		if !addAccumulatorMetadataBytes(&total, choice.FinishReason) ||
			!addAccumulatorMetadataBytes(&total, string(choice.Message.Role)) ||
			!addAccumulatorStringReconciliationWork(work, choice.FinishReason, choiceState.finishReason) ||
			!addAccumulatorStringReconciliationWork(work, choice.Message.Role, choiceState.role) {
			return 0, false
		}
		if !choiceState.visitReconciledTools(i, choice.Message.ToolCalls, chunk, func(j int) bool {
			toolCall := &choice.Message.ToolCalls[j]
			toolState := choiceState.toolCallState(j)
			if !addAccumulatorMetadataBytes(&total, toolCall.ID) ||
				!addAccumulatorMetadataBytes(&total, toolCall.Type) ||
				!addAccumulatorStringReconciliationWork(work, toolCall.ID, toolState.id) ||
				!addAccumulatorStringReconciliationWork(work, toolCall.Type, toolState.typeName) {
				return false
			}
			return true
		}) {
			return 0, false
		}
	}
	return total, true
}

func addAccumulatorStringReconciliationWork[T ~string](work *int, current T, published string) bool {
	value := string(current)
	if accumulatorStringUsesPublishedBacking(value, published) {
		return true
	}
	return addAccumulatorStringCopyWork(work, value, 2)
}

func addAccumulatorStringCopyWork[T ~string](work *int, value T, passes int) bool {
	for range passes {
		if !addAccumulatorReconciliationWork(work, len(value)) {
			return false
		}
	}
	return true
}

// update updates the internal response state and returns the previous state if
// the state changed. This ensures that JustFinished events only fire once.
func (prev *chatCompletionResponseStateCode) update(choice ChatCompletionChunkChoice, toolCallCount int) (justFinished chatCompletionResponseState) {
	delta := choice.Delta
	choiceIndex, _ := checkedStreamAccumulatorChoiceIndex(choice.Index)
	new := chatCompletionResponseState{choiceIndex: choiceIndex}
	switch {
	case len(delta.ToolCalls) > 0 && delta.Content == "":
		new.state = toolResponseState
		new.toolCallIndex, _ = checkedToolCallIndex(delta.ToolCalls[0].Index, toolCallCount)
	case delta.JSON.Content.Valid():
		new.state = contentResponseState
	case delta.JSON.Refusal.Valid():
		new.state = refusalResponseState
	case len(delta.ToolCalls) > 0:
		new.state = toolResponseState
		new.toolCallIndex, _ = checkedToolCallIndex(delta.ToolCalls[0].Index, toolCallCount)
	default:
		new.state = finishedResponseState
	}

	newCode := encodeChatCompletionResponseState(new)
	if *prev != newCode {
		justFinished = prev.decode(choiceIndex)
	}
	*prev = newCode

	return
}

func encodeChatCompletionResponseState(state chatCompletionResponseState) chatCompletionResponseStateCode {
	return chatCompletionResponseStateCode{state: state.state, toolCallIndex: state.toolCallIndex}
}

func (code chatCompletionResponseStateCode) decode(choiceIndex int) chatCompletionResponseState {
	return chatCompletionResponseState{
		state:         code.state,
		choiceIndex:   choiceIndex,
		toolCallIndex: code.toolCallIndex,
	}
}

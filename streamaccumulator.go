package openai

import (
	"slices"
	"unsafe"

	"github.com/openai/openai-go/v3/shared/constant"
)

const (
	maxChatCompletionAccumulatorChunks          = 100_000
	maxChatCompletionAccumulatorStructuralSlots = 1_024
	maxChatCompletionAccumulatorTextBytes       = 16 << 20
	// Deterministic doubling preserves capacity < 2*length for non-empty buffers.
	maxChatCompletionAccumulatorTextCapacity  = 2 * maxChatCompletionAccumulatorTextBytes
	maxChatCompletionAccumulatorMetadataBytes = 16 << 20
	maxChatCompletionAccumulatorLogprobBytes  = 16 << 20
	maxChatCompletionAccumulatorReconcileWork = 64_000_000
	chatCompletionAccumulatorChoiceWork       = 16
	chatCompletionAccumulatorToolWork         = 8
	// Conservatively covers a non-empty map header and its first backing bucket.
	chatCompletionAccumulatorMapOverheadBytes = 512
)

// Helper to accumulate chunks from a stream
type ChatCompletionAccumulator struct {
	// The up-to-date accumulation of model's responses
	ChatCompletion
	choiceChatCompletionStates       []chatCompletionResponseState
	legacyChoiceChatCompletionStates []chatCompletionResponseState
	justFinished                     chatCompletionResponseState
	justFinishedByChoice             []chatCompletionResponseState
	stringState                      chatCompletionAccumulatorStringState
	logprobState                     chatCompletionAccumulatorLogprobState
	chunkCount                       int
	textBytes                        int
	logprobBytes                     int
	reconciliationWork               int
	privateStateOwner                *ChatCompletionAccumulator
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
	toolCalls       []*chatCompletionToolCallStringState
	activeToolCalls []int
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
// 127, at most 100,000 chunks and 1,024 combined choice and tool-call slots, 16 MiB
// each of cumulative remote text, live combined content and tool function text,
// other retained string metadata, cumulative remote log probability data, and live
// retained log probability data, 32 MiB of retained text-buffer capacity, and 64
// million accumulation work units.
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
	acc.ensurePrivateStateOwned()
	acc.justFinished = chatCompletionResponseState{}
	acc.justFinishedByChoice = acc.justFinishedByChoice[:0]
	if !acc.preflightChunk(&chunk) {
		return false
	}

	acc.accumulateDelta(&chunk)
	acc.logprobState.acceptChunk(&acc.ChatCompletion, &chunk)
	acc.chunkCount++

	if len(chunk.Choices) > 0 {
		firstChoice := chunk.Choices[0]
		choiceIndex, _ := checkedStreamAccumulatorChoiceIndex(firstChoice.Index)
		acc.legacyChoiceChatCompletionStates = expandToFit(acc.legacyChoiceChatCompletionStates, choiceIndex)
		toolCallCount := len(acc.Choices[choiceIndex].Message.ToolCalls)
		acc.justFinished = acc.legacyChoiceChatCompletionStates[choiceIndex].update(firstChoice, toolCallCount)
	}

	for _, choice := range chunk.Choices {
		choiceIndex, _ := checkedStreamAccumulatorChoiceIndex(choice.Index)
		acc.choiceChatCompletionStates = expandToFit(acc.choiceChatCompletionStates, choiceIndex)
		toolCallCount := len(acc.Choices[choiceIndex].Message.ToolCalls)
		justFinished := acc.choiceChatCompletionStates[choiceIndex].update(choice, toolCallCount)
		if justFinished.state != emptyResponseState {
			acc.justFinishedByChoice = append(acc.justFinishedByChoice, justFinished)
		}
	}
	return true
}

func (acc *ChatCompletionAccumulator) ensurePrivateStateOwned() {
	if acc.privateStateOwner == acc {
		return
	}
	if acc.privateStateOwner == nil {
		acc.privateStateOwner = acc
		return
	}

	acc.choiceChatCompletionStates = cloneAccumulatorSlice(acc.choiceChatCompletionStates)
	acc.legacyChoiceChatCompletionStates = cloneAccumulatorSlice(acc.legacyChoiceChatCompletionStates)
	acc.justFinishedByChoice = cloneAccumulatorSlice(acc.justFinishedByChoice)
	acc.stringState = acc.stringState.cloneForAccumulatorCopy()
	acc.logprobState.choices = cloneAccumulatorSlice(acc.logprobState.choices)
	acc.privateStateOwner = acc
}

func (state chatCompletionAccumulatorStringState) cloneForAccumulatorCopy() chatCompletionAccumulatorStringState {
	state.choices = cloneAccumulatorSlice(state.choices)
	state.activeChoices = cloneAccumulatorSlice(state.activeChoices)
	for i, choice := range state.choices {
		if choice == nil {
			continue
		}
		choiceCopy := *choice
		choiceCopy.content.shared = len(choiceCopy.content.buffer) > 0
		choiceCopy.refusal.shared = len(choiceCopy.refusal.buffer) > 0
		choiceCopy.toolCalls = cloneAccumulatorSlice(choiceCopy.toolCalls)
		choiceCopy.activeToolCalls = cloneAccumulatorSlice(choiceCopy.activeToolCalls)
		for j, toolCall := range choiceCopy.toolCalls {
			if toolCall == nil {
				continue
			}
			toolCallCopy := *toolCall
			toolCallCopy.name.shared = len(toolCallCopy.name.buffer) > 0
			toolCallCopy.arguments.shared = len(toolCallCopy.arguments.buffer) > 0
			choiceCopy.toolCalls[j] = &toolCallCopy
		}
		state.choices[i] = &choiceCopy
	}
	return state
}

func (acc *ChatCompletionAccumulator) preflightChunk(chunk *ChatCompletionChunk) bool {
	if acc.chunkCount >= maxChatCompletionAccumulatorChunks {
		return false
	}
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
	if !acc.chatCompletionStructuralSlotsWithinLimit(chunk) {
		return false
	}

	projectedTextBytes, ok := addChatCompletionChunkTextBytes(acc.textBytes, chunk)
	if !ok {
		return false
	}
	liveTextBytes, ok := acc.addChatCompletionTextBytes(0, &projectedReconciliationWork, chunk)
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

	if hasLogprobPlan {
		acc.applyLogprobReconciliation(&logprobPlan)
	}
	acc.reconcilePublicState()
	acc.textBytes = projectedTextBytes
	acc.logprobBytes = projectedLogprobBytes
	acc.reconciliationWork = projectedReconciliationWork
	return true
}

func (acc *ChatCompletionAccumulator) projectReconciliationWork(chunk *ChatCompletionChunk) (int, bool) {
	// Preflight and commit make six passes over populated choices, plus one
	// structural projection pass when the chunk has choices. Populated tools are
	// visited by text, metadata, and commit reconciliation. Charging those exact
	// passes keeps sparse detection work bounded without scanning placeholder slots.
	// Per-chunk choice and tool units conservatively cover the remaining fixed
	// validation, projection, accumulation, and state-update passes over new input.
	choicePasses := 6
	if len(chunk.Choices) > 0 {
		choicePasses++
	}
	work := len(acc.stringState.activeChoices)*choicePasses + acc.stringState.activeTools*3
	if work > maxChatCompletionAccumulatorReconcileWork-acc.reconciliationWork {
		return 0, false
	}
	projected := acc.reconciliationWork + work
	if len(chunk.Choices) > (maxChatCompletionAccumulatorReconcileWork-projected)/chatCompletionAccumulatorChoiceWork {
		return 0, false
	}
	projected += len(chunk.Choices) * chatCompletionAccumulatorChoiceWork
	for i := range chunk.Choices {
		toolCount := len(chunk.Choices[i].Delta.ToolCalls)
		if toolCount > (maxChatCompletionAccumulatorReconcileWork-projected)/chatCompletionAccumulatorToolWork {
			return 0, false
		}
		projected += toolCount * chatCompletionAccumulatorToolWork
	}
	return projected, true
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

// Concatenates a preflighted ChatCompletionChunk onto a ChatCompletion.
// Ignores the JSON field.
func (acc *ChatCompletionAccumulator) accumulateDelta(chunk *ChatCompletionChunk) {
	cc := &acc.ChatCompletion
	if len(cc.ID) == 0 {
		assignAccumulatorString(&acc.stringState.id, &cc.ID, chunk.ID)
	}

	for _, delta := range chunk.Choices {
		choiceIndex, _ := checkedStreamAccumulatorChoiceIndex(delta.Index)
		cc.Choices = expandToFit(cc.Choices, choiceIndex)
		choice := &cc.Choices[choiceIndex]
		choiceStrings := acc.stringState.choice(choiceIndex)

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

		for j := range delta.Delta.ToolCalls {
			deltaTool := &delta.Delta.ToolCalls[j]
			toolIndex, _ := checkedToolCallIndex(deltaTool.Index, len(choice.Message.ToolCalls))
			tool := &choice.Message.ToolCalls[toolIndex]
			toolStrings := choiceStrings.toolCall(toolIndex, &acc.stringState.activeTools)

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

func (acc *chatCompletionAccumulatorStringState) choice(index int) *chatCompletionChoiceStringState {
	acc.choices = expandToFit(acc.choices, index)
	if acc.choices[index] == nil {
		acc.choices[index] = &chatCompletionChoiceStringState{}
		acc.activeChoices = append(acc.activeChoices, index)
	}
	return acc.choices[index]
}

func (acc *ChatCompletionAccumulator) reconcilePublicState() {
	completion := &acc.ChatCompletion
	stringState := &acc.stringState
	reconcileAccumulatorString(&stringState.id, &completion.ID)
	reconcileAccumulatorString(&stringState.model, &completion.Model)
	reconcileAccumulatorString(&stringState.systemFingerprint, &completion.SystemFingerprint)
	reconcileAccumulatorString(&stringState.object, &completion.Object)
	reconcileAccumulatorString(&stringState.serviceTier, &completion.ServiceTier)

	previousChoiceCount := len(stringState.choices)
	completion.Choices = detachTruncatedTail(completion.Choices, previousChoiceCount)
	for _, i := range stringState.activeChoices {
		if i >= len(completion.Choices) {
			continue
		}
		message := &completion.Choices[i].Message
		previousToolCallCount := len(stringState.choices[i].toolCalls)
		message.ToolCalls = detachTruncatedTail(message.ToolCalls, previousToolCallCount)
	}

	activeChoiceCount := 0
	for _, i := range stringState.activeChoices {
		if i >= len(completion.Choices) {
			stringState.activeTools -= len(stringState.choices[i].activeToolCalls)
			continue
		}
		stringState.activeChoices[activeChoiceCount] = i
		activeChoiceCount++
		choiceState := stringState.choices[i]

		message := &completion.Choices[i].Message
		choiceState.content.reconcile(&message.Content)
		choiceState.refusal.reconcile(&message.Refusal)
		reconcileAccumulatorString(&choiceState.finishReason, &completion.Choices[i].FinishReason)
		reconcileAccumulatorString(&choiceState.role, &message.Role)

		if len(message.ToolCalls) < len(choiceState.toolCalls) {
			invalidateRemovedToolCallState(acc.legacyChoiceChatCompletionStates, i, len(message.ToolCalls))
			invalidateRemovedToolCallState(acc.choiceChatCompletionStates, i, len(message.ToolCalls))
		}
		activeToolCallCount := 0
		for _, j := range choiceState.activeToolCalls {
			if j >= len(message.ToolCalls) {
				continue
			}
			choiceState.activeToolCalls[activeToolCallCount] = j
			activeToolCallCount++
			toolCallState := choiceState.toolCalls[j]
			reconcileAccumulatorString(&toolCallState.id, &message.ToolCalls[j].ID)
			reconcileAccumulatorString(&toolCallState.typeName, &message.ToolCalls[j].Type)
			function := &message.ToolCalls[j].Function
			toolCallState.name.reconcile(&function.Name)
			toolCallState.arguments.reconcile(&function.Arguments)
		}
		stringState.activeTools -= len(choiceState.activeToolCalls) - activeToolCallCount
		clear(choiceState.activeToolCalls[activeToolCallCount:])
		choiceState.activeToolCalls = choiceState.activeToolCalls[:activeToolCallCount]
		if len(message.ToolCalls) < len(choiceState.toolCalls) {
			clear(choiceState.toolCalls[len(message.ToolCalls):])
			choiceState.toolCalls = choiceState.toolCalls[:len(message.ToolCalls)]
		}
	}
	clear(stringState.activeChoices[activeChoiceCount:])
	stringState.activeChoices = stringState.activeChoices[:activeChoiceCount]
	choiceCount := len(completion.Choices)
	if choiceCount < previousChoiceCount {
		clear(stringState.choices[choiceCount:])
		stringState.choices = stringState.choices[:choiceCount]
	}
	acc.legacyChoiceChatCompletionStates = truncateResponseStates(acc.legacyChoiceChatCompletionStates, choiceCount)
	acc.choiceChatCompletionStates = truncateResponseStates(acc.choiceChatCompletionStates, choiceCount)
}

func invalidateRemovedToolCallState(states []chatCompletionResponseState, choiceIndex int, toolCallCount int) {
	if choiceIndex >= len(states) {
		return
	}
	state := &states[choiceIndex]
	if state.state == toolResponseState && state.toolCallIndex >= toolCallCount {
		*state = chatCompletionResponseState{}
	}
}

func truncateResponseStates(states []chatCompletionResponseState, choiceCount int) []chatCompletionResponseState {
	choiceCount = min(len(states), choiceCount)
	clear(states[choiceCount:])
	return states[:choiceCount]
}

func (choice *chatCompletionChoiceStringState) toolCall(index int, activeTools *int) *chatCompletionToolCallStringState {
	choice.toolCalls = expandToFit(choice.toolCalls, index)
	if choice.toolCalls[index] == nil {
		choice.toolCalls[index] = &chatCompletionToolCallStringState{}
		choice.activeToolCalls = append(choice.activeToolCalls, index)
		(*activeTools)++
	}
	return choice.toolCalls[index]
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
		if capacity > maxChatCompletionAccumulatorTextBytes/2 {
			return maxChatCompletionAccumulatorTextBytes
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

func (acc *ChatCompletionAccumulator) chatCompletionStructuralSlotsWithinLimit(chunk *ChatCompletionChunk) bool {
	completion := &acc.ChatCompletion
	choiceCount := len(completion.Choices)
	if choiceCount > maxChatCompletionAccumulatorStructuralSlots {
		return false
	}

	structuralSlots := choiceCount
	for _, i := range acc.stringState.activeChoices {
		if i >= len(completion.Choices) {
			continue
		}
		toolCallCount := len(completion.Choices[i].Message.ToolCalls)
		if toolCallCount > maxChatCompletionAccumulatorStructuralSlots-structuralSlots {
			return false
		}
		structuralSlots += toolCallCount
	}
	if len(chunk.Choices) == 0 {
		return true
	}
	// This fixed-size projection keeps duplicate choice deltas exact without a
	// per-chunk map allocation. The aggregate slot limit bounds its stack cost.
	var toolCallCounts [maxChatCompletionAccumulatorStructuralSlots]int
	for _, i := range acc.stringState.activeChoices {
		if i < len(completion.Choices) {
			toolCallCounts[i] = len(completion.Choices[i].Message.ToolCalls)
		}
	}

	chunkEntries := len(chunk.Choices)
	if chunkEntries > maxChatCompletionAccumulatorStructuralSlots {
		return false
	}
	for i := range chunk.Choices {
		toolCallCount := len(chunk.Choices[i].Delta.ToolCalls)
		if toolCallCount > maxChatCompletionAccumulatorStructuralSlots-chunkEntries {
			return false
		}
		chunkEntries += toolCallCount
	}

	for i := range chunk.Choices {
		delta := &chunk.Choices[i]
		if delta.Index < 0 || delta.Index >= maxChatCompletionAccumulatorStructuralSlots {
			return false
		}
		choiceIndex := int(delta.Index)
		if newChoiceCount := choiceIndex + 1; newChoiceCount > choiceCount {
			additionalChoices := newChoiceCount - choiceCount
			if additionalChoices > maxChatCompletionAccumulatorStructuralSlots-structuralSlots {
				return false
			}
			structuralSlots += additionalChoices
			choiceCount = newChoiceCount
		}

		for j := range delta.Delta.ToolCalls {
			toolIndex64 := delta.Delta.ToolCalls[j].Index
			toolIndex := preflightedToolCallIndex(toolIndex64)
			newToolCallCount := toolIndex + 1
			if newToolCallCount <= toolCallCounts[choiceIndex] {
				continue
			}
			additionalToolCalls := newToolCallCount - toolCallCounts[choiceIndex]
			if additionalToolCalls > maxChatCompletionAccumulatorStructuralSlots-structuralSlots {
				return false
			}
			structuralSlots += additionalToolCalls
			toolCallCounts[choiceIndex] = newToolCallCount
		}
	}
	return true
}

func (acc *ChatCompletionAccumulator) addChatCompletionTextBytes(total int, work *int, chunk *ChatCompletionChunk) (int, bool) {
	completion := &acc.ChatCompletion
	var appends *chatCompletionTextAppendProjection
	capacity := 0
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
		capacity += projectedAccumulatorBufferCapacity(message.Content, &choiceState.content)
		capacity += projectedAccumulatorBufferCapacity(message.Refusal, &choiceState.refusal)
		for _, j := range choiceState.activeToolCalls {
			if j >= len(message.ToolCalls) {
				continue
			}
			function := &message.ToolCalls[j].Function
			toolState := choiceState.toolCalls[j]
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
				return 0, false
			}
			capacity += projectedAccumulatorBufferCapacity(function.Name, &toolState.name)
			capacity += projectedAccumulatorBufferCapacity(function.Arguments, &toolState.arguments)
		}
	}
	return total, capacity <= maxChatCompletionAccumulatorTextCapacity
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
	return addAccumulatorStringCopyWork(work, current, passes)
}

func projectedAccumulatorBufferCapacity(current string, state *chatCompletionString) int {
	if !accumulatorStringUsesPublishedBacking(current, state.published) {
		return len(current)
	}
	return cap(state.buffer)
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

func (acc *ChatCompletionAccumulator) addChatCompletionMetadataBytes(total int, work *int) (int, bool) {
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
		for _, j := range choiceState.activeToolCalls {
			if j >= len(choice.Message.ToolCalls) {
				continue
			}
			toolCall := &choice.Message.ToolCalls[j]
			toolState := choiceState.toolCalls[j]
			if !addAccumulatorMetadataBytes(&total, toolCall.ID) ||
				!addAccumulatorMetadataBytes(&total, toolCall.Type) ||
				!addAccumulatorStringReconciliationWork(work, toolCall.ID, toolState.id) ||
				!addAccumulatorStringReconciliationWork(work, toolCall.Type, toolState.typeName) {
				return 0, false
			}
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
func (prev *chatCompletionResponseState) update(choice ChatCompletionChunkChoice, toolCallCount int) (justFinished chatCompletionResponseState) {
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

	if *prev != new {
		justFinished = *prev
	}
	*prev = new

	return
}

func (acc *ChatCompletionAccumulator) validChatCompletionChunkIndices(chunk ChatCompletionChunk) bool {
	var toolCallCounts [maxStreamAccumulatorChoiceIndex + 1]int
	var initializedToolCallCounts [maxStreamAccumulatorChoiceIndex + 1]bool

	for _, choice := range chunk.Choices {
		choiceIndex, ok := checkedStreamAccumulatorChoiceIndex(choice.Index)
		if !ok {
			return false
		}
		if !initializedToolCallCounts[choiceIndex] {
			if choiceIndex < len(acc.Choices) {
				toolCallCounts[choiceIndex] = len(acc.Choices[choiceIndex].Message.ToolCalls)
			}
			initializedToolCallCounts[choiceIndex] = true
		}
		for _, toolCall := range choice.Delta.ToolCalls {
			toolIndex, ok := checkedToolCallIndex(toolCall.Index, toolCallCounts[choiceIndex])
			if !ok {
				return false
			}
			if toolIndex >= toolCallCounts[choiceIndex] {
				toolCallCounts[choiceIndex] = toolIndex + 1
			}
		}
	}
	return true
}

func checkedStreamAccumulatorChoiceIndex(index int64) (int, bool) {
	if index < 0 || index > maxStreamAccumulatorChoiceIndex {
		return 0, false
	}
	return int(index), true
}

// checkedToolCallIndex handles providers like AWS Bedrock that return -1 for a
// single tool call. Tool calls have no protocol maximum, so the bound applies
// only to the growth caused by this index instead of the accumulated count.
func checkedToolCallIndex(index int64, toolCallCount int) (int, bool) {
	if index == -1 {
		index = 0
	}
	// The maximum int value cannot be used because slice growth needs index+1.
	if index < 0 || index >= int64(^uint(0)>>1) {
		return 0, false
	}
	if index >= int64(toolCallCount) && index-int64(toolCallCount) >= maxStreamAccumulatorToolCallGrowth {
		return 0, false
	}
	return int(index), true
}

func preflightedToolCallIndex(index int64) int {
	if index == -1 {
		return 0
	}
	return int(index)
}

func expandToFit[T any](slice []T, index int) []T {
	if index < len(slice) {
		return slice
	}
	return slices.Grow(slice, index+1-len(slice))[:index+1]
}

func detachTruncatedTail[T any](slice []T, previousLength int) []T {
	// A full slice expression can hide a retained tail while reporting len == cap.
	if slice == nil || (len(slice) == cap(slice) && len(slice) >= previousLength) {
		return slice
	}
	detached := make([]T, len(slice))
	copy(detached, slice)
	return detached
}

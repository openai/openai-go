package openai

import (
	"strings"

	"github.com/openai/openai-go/v3/shared/constant"
)

const (
	maxChatCompletionAccumulatorChunks          = 100_000
	maxChatCompletionAccumulatorStructuralSlots = 1_024
	maxChatCompletionAccumulatorTextBytes       = 16 << 20
	maxChatCompletionAccumulatorMetadataBytes   = 16 << 20
	maxChatCompletionAccumulatorLogprobBytes    = 16 << 20
	maxChatCompletionAccumulatorLogprobWork     = 64_000_000
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
	// Keep builders behind pointers because a non-zero strings.Builder must not be copied.
	choices       []*chatCompletionChoiceStringState
	activeChoices []int
}

type chatCompletionChoiceStringState struct {
	content         chatCompletionString
	refusal         chatCompletionString
	toolCalls       []*chatCompletionToolCallStringState
	activeToolCalls []int
}

type chatCompletionToolCallStringState struct {
	name      chatCompletionString
	arguments chatCompletionString
}

type chatCompletionString struct {
	builder   strings.Builder
	published string
}

type chatCompletionResponseStateEnum int

const (
	emptyResponseState chatCompletionResponseStateEnum = iota
	contentResponseState
	refusalResponseState
	toolResponseState
	finishedResponseState
)

// AddChunk incorporates a chunk into the accumulation. Chunks must be added in order.
// Returns false if the chunk could not be successfully accumulated. To bound work and
// memory for untrusted streams, an accumulator accepts at most 100,000 chunks, 1,024
// combined choice and tool-call slots, 16 MiB each of combined content and tool
// function text, other retained string metadata, and aggregate retained log
// probability data, and 64 million retained log-probability reconciliation steps.
// A rejected chunk does not modify the accumulator.
//
// The ChatCompletion field JSON does not get accumulated.
func (acc *ChatCompletionAccumulator) AddChunk(chunk ChatCompletionChunk) bool {
	if !acc.preflightChunk(&chunk) {
		return false
	}

	acc.justFinished = chatCompletionResponseState{}
	acc.justFinishedByChoice = acc.justFinishedByChoice[:0]
	acc.accumulateDelta(&chunk)
	acc.logprobState.acceptChunk(&acc.ChatCompletion, &chunk)
	acc.chunkCount++

	if len(chunk.Choices) > 0 {
		firstChoice := chunk.Choices[0]
		choiceIndex := int(firstChoice.Index)
		acc.legacyChoiceChatCompletionStates = expandToFit(acc.legacyChoiceChatCompletionStates, choiceIndex)
		acc.justFinished = acc.legacyChoiceChatCompletionStates[choiceIndex].update(firstChoice)
	}

	for _, choice := range chunk.Choices {
		choiceIndex := int(choice.Index)
		acc.choiceChatCompletionStates = expandToFit(acc.choiceChatCompletionStates, choiceIndex)
		justFinished := acc.choiceChatCompletionStates[choiceIndex].update(choice)
		if justFinished.state != emptyResponseState {
			acc.justFinishedByChoice = append(acc.justFinishedByChoice, justFinished)
		}
	}
	return true
}

func (acc *ChatCompletionAccumulator) preflightChunk(chunk *ChatCompletionChunk) bool {
	if acc.chunkCount >= maxChatCompletionAccumulatorChunks {
		return false
	}
	if acc.ID != "" && acc.ID != chunk.ID {
		return false
	}
	if !acc.chatCompletionStructuralSlotsWithinLimit(chunk) {
		return false
	}

	textBytes, ok := acc.addChatCompletionTextBytes(0)
	if !ok {
		return false
	}
	if _, ok = addChatCompletionChunkTextBytes(textBytes, chunk); !ok {
		return false
	}
	if !acc.chatCompletionMetadataWithinLimit(chunk) {
		return false
	}
	var logprobPlan chatCompletionLogprobReconcilePlan
	hasLogprobPlan, ok := acc.planLogprobReconciliation(&logprobPlan, chatCompletionChunkHasLogprobs(chunk))
	if !ok {
		return false
	}
	logprobState := &acc.logprobState
	if hasLogprobPlan {
		logprobState = &logprobPlan.state
	}
	if !logprobState.chunkWithinLimit(chunk) {
		return false
	}

	if hasLogprobPlan {
		acc.applyLogprobReconciliation(&logprobPlan)
	}
	acc.reconcilePublicState()
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

// Concatenates a preflighted ChatCompletionChunk onto a ChatCompletion.
// Ignores the JSON field.
func (acc *ChatCompletionAccumulator) accumulateDelta(chunk *ChatCompletionChunk) {
	cc := &acc.ChatCompletion
	if len(cc.ID) == 0 {
		assignAccumulatorString(&cc.ID, chunk.ID)
	}

	for _, delta := range chunk.Choices {
		cc.Choices = expandToFit(cc.Choices, int(delta.Index))
		choice := &cc.Choices[delta.Index]

		choice.Index = delta.Index
		assignAccumulatorString(&choice.FinishReason, delta.FinishReason)

		if delta.Delta.Role != "" {
			assignAccumulatorString(&choice.Message.Role, constant.Assistant(delta.Delta.Role))
		}

		choiceStrings := acc.stringState.choice(int(delta.Index))
		choiceStrings.content.append(&choice.Message.Content, delta.Delta.Content)
		choiceStrings.refusal.append(&choice.Message.Refusal, delta.Delta.Refusal)

		for j := range delta.Delta.ToolCalls {
			deltaTool := &delta.Delta.ToolCalls[j]
			// Clamp negative indices to 0 since the API may send -1 for single tool calls
			toolIndex := clampToZero(deltaTool.Index)

			choice.Message.ToolCalls = expandToFit(choice.Message.ToolCalls, toolIndex)
			tool := &choice.Message.ToolCalls[toolIndex]

			if deltaTool.ID != "" {
				assignAccumulatorString(&tool.ID, deltaTool.ID)
			}
			if deltaTool.Type != "" {
				assignAccumulatorString(&tool.Type, deltaTool.Type)
			}
			toolStrings := choiceStrings.toolCall(toolIndex)
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

	assignAccumulatorString(&cc.Model, chunk.Model)
	cc.Created = chunk.Created
	assignAccumulatorString(&cc.SystemFingerprint, chunk.SystemFingerprint)
	assignAccumulatorString(&cc.ServiceTier, ChatCompletionServiceTier(chunk.ServiceTier))
	if chunk.Object == chunk.Object.Default() {
		cc.Object = cc.Object.Default()
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
			continue
		}
		stringState.activeChoices[activeChoiceCount] = i
		activeChoiceCount++
		choiceState := stringState.choices[i]

		message := &completion.Choices[i].Message
		choiceState.content.reconcile(&message.Content)
		choiceState.refusal.reconcile(&message.Refusal)

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
			function := &message.ToolCalls[j].Function
			toolCallState.name.reconcile(&function.Name)
			toolCallState.arguments.reconcile(&function.Arguments)
		}
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

func (choice *chatCompletionChoiceStringState) toolCall(index int) *chatCompletionToolCallStringState {
	choice.toolCalls = expandToFit(choice.toolCalls, index)
	if choice.toolCalls[index] == nil {
		choice.toolCalls[index] = &chatCompletionToolCallStringState{}
		choice.activeToolCalls = append(choice.activeToolCalls, index)
	}
	return choice.toolCalls[index]
}

func (acc *chatCompletionString) append(current *string, fragment string) {
	if fragment == "" {
		return
	}
	acc.reconcile(current)
	_, _ = acc.builder.WriteString(fragment)
	acc.published = acc.builder.String()
	*current = acc.published
}

func (acc *chatCompletionString) reconcile(current *string) {
	if *current == acc.published {
		*current = acc.published
		return
	}

	replacement := *current
	acc.published = ""
	acc.builder.Reset()
	_, _ = acc.builder.WriteString(replacement)
	acc.published = acc.builder.String()
	*current = acc.published
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
			if toolIndex64 >= maxChatCompletionAccumulatorStructuralSlots {
				return false
			}
			toolIndex := clampToZero(toolIndex64)
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

func (acc *ChatCompletionAccumulator) addChatCompletionTextBytes(total int) (int, bool) {
	completion := &acc.ChatCompletion
	for _, i := range acc.stringState.activeChoices {
		if i >= len(completion.Choices) {
			continue
		}
		message := &completion.Choices[i].Message
		if !addAccumulatorTextBytes(&total, message.Content) ||
			!addAccumulatorTextBytes(&total, message.Refusal) {
			return 0, false
		}
		for _, j := range acc.stringState.choices[i].activeToolCalls {
			if j >= len(message.ToolCalls) {
				continue
			}
			function := &message.ToolCalls[j].Function
			if !addAccumulatorTextBytes(&total, function.Name) ||
				!addAccumulatorTextBytes(&total, function.Arguments) {
				return 0, false
			}
		}
	}
	return total, true
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

func (acc *ChatCompletionAccumulator) addChatCompletionMetadataBytes(total int) (int, bool) {
	completion := &acc.ChatCompletion
	if !addAccumulatorMetadataBytes(&total, completion.ID) ||
		!addAccumulatorMetadataBytes(&total, completion.Model) ||
		!addAccumulatorMetadataBytes(&total, completion.SystemFingerprint) ||
		!addAccumulatorMetadataBytes(&total, string(completion.Object)) ||
		!addAccumulatorMetadataBytes(&total, string(completion.ServiceTier)) {
		return 0, false
	}
	for _, i := range acc.stringState.activeChoices {
		if i >= len(completion.Choices) {
			continue
		}
		choice := &completion.Choices[i]
		if !addAccumulatorMetadataBytes(&total, choice.FinishReason) ||
			!addAccumulatorMetadataBytes(&total, string(choice.Message.Role)) {
			return 0, false
		}
		for _, j := range acc.stringState.choices[i].activeToolCalls {
			if j >= len(choice.Message.ToolCalls) {
				continue
			}
			toolCall := &choice.Message.ToolCalls[j]
			if !addAccumulatorMetadataBytes(&total, toolCall.ID) ||
				!addAccumulatorMetadataBytes(&total, toolCall.Type) {
				return 0, false
			}
		}
	}
	return total, true
}

type chatCompletionToolMetadataProjection struct {
	key    int
	fields uint8
	used   bool
}

const (
	projectedToolID uint8 = 1 << iota
	projectedToolType
)

func (acc *ChatCompletionAccumulator) chatCompletionMetadataWithinLimit(chunk *ChatCompletionChunk) bool {
	completion := &acc.ChatCompletion
	total, ok := acc.addChatCompletionMetadataBytes(0)
	if !ok {
		return false
	}
	delta := int64(0)
	if completion.ID == "" {
		delta += int64(len(chunk.ID))
	}
	delta += int64(len(chunk.Model) - len(completion.Model))
	delta += int64(len(chunk.SystemFingerprint) - len(completion.SystemFingerprint))
	delta += int64(len(chunk.ServiceTier) - len(completion.ServiceTier))
	if chunk.Object == chunk.Object.Default() {
		delta += int64(len(completion.Object.Default()) - len(completion.Object))
	}
	if len(chunk.Choices) == 0 {
		projected := int64(total) + delta
		return projected >= 0 && projected <= maxChatCompletionAccumulatorMetadataBytes
	}

	var finishReasonSeen [maxChatCompletionAccumulatorStructuralSlots / 64]uint64
	var roleSeen [maxChatCompletionAccumulatorStructuralSlots / 64]uint64
	for i := len(chunk.Choices) - 1; i >= 0; i-- {
		choiceDelta := &chunk.Choices[i]
		choiceIndex := int(choiceDelta.Index)
		var currentFinishReasonBytes, currentRoleBytes int
		if choiceIndex < len(completion.Choices) {
			choice := &completion.Choices[choiceIndex]
			currentFinishReasonBytes = len(choice.FinishReason)
			currentRoleBytes = len(choice.Message.Role)
		}
		if markChatCompletionMetadataProjected(&finishReasonSeen, choiceIndex) {
			deltaBytes := len(choiceDelta.FinishReason) - currentFinishReasonBytes
			delta += int64(deltaBytes)
		}
		if choiceDelta.Delta.Role != "" && markChatCompletionMetadataProjected(&roleSeen, choiceIndex) {
			deltaBytes := len(choiceDelta.Delta.Role) - currentRoleBytes
			delta += int64(deltaBytes)
		}
	}
	delta += chatCompletionToolMetadataDelta(completion, chunk)
	projected := int64(total) + delta
	return projected >= 0 && projected <= maxChatCompletionAccumulatorMetadataBytes
}

func markChatCompletionMetadataProjected(seen *[maxChatCompletionAccumulatorStructuralSlots / 64]uint64, index int) bool {
	word := index / 64
	mask := uint64(1) << (index % 64)
	first := seen[word]&mask == 0
	seen[word] |= mask
	return first
}

func chatCompletionToolMetadataDelta(completion *ChatCompletion, chunk *ChatCompletionChunk) int64 {
	hasMetadata := false
	for i := range chunk.Choices {
		for j := range chunk.Choices[i].Delta.ToolCalls {
			toolCall := &chunk.Choices[i].Delta.ToolCalls[j]
			hasMetadata = hasMetadata || toolCall.ID != "" || toolCall.Type != ""
		}
	}
	if !hasMetadata {
		return 0
	}
	return nonEmptyChatCompletionToolMetadataDelta(completion, chunk)
}

func nonEmptyChatCompletionToolMetadataDelta(completion *ChatCompletion, chunk *ChatCompletionChunk) int64 {
	var projections [maxChatCompletionAccumulatorStructuralSlots * 2]chatCompletionToolMetadataProjection
	deltaBytes := int64(0)
	for i := len(chunk.Choices) - 1; i >= 0; i-- {
		choice := &chunk.Choices[i]
		choiceIndex := int(choice.Index)
		for j := len(choice.Delta.ToolCalls) - 1; j >= 0; j-- {
			toolCall := &choice.Delta.ToolCalls[j]
			if toolCall.ID == "" && toolCall.Type == "" {
				continue
			}
			toolIndex := clampToZero(toolCall.Index)
			projection := lookupChatCompletionToolMetadataProjection(&projections, choiceIndex, toolIndex)
			var currentIDBytes, currentTypeBytes int
			if choiceIndex < len(completion.Choices) && toolIndex < len(completion.Choices[choiceIndex].Message.ToolCalls) {
				current := &completion.Choices[choiceIndex].Message.ToolCalls[toolIndex]
				currentIDBytes = len(current.ID)
				currentTypeBytes = len(current.Type)
			}
			if toolCall.ID != "" && projection.fields&projectedToolID == 0 {
				deltaBytes += int64(len(toolCall.ID) - currentIDBytes)
				projection.fields |= projectedToolID
			}
			if toolCall.Type != "" && projection.fields&projectedToolType == 0 {
				deltaBytes += int64(len(toolCall.Type) - currentTypeBytes)
				projection.fields |= projectedToolType
			}
		}
	}
	return deltaBytes
}

func lookupChatCompletionToolMetadataProjection(
	projections *[maxChatCompletionAccumulatorStructuralSlots * 2]chatCompletionToolMetadataProjection,
	choiceIndex int,
	toolIndex int,
) *chatCompletionToolMetadataProjection {
	key := choiceIndex*maxChatCompletionAccumulatorStructuralSlots + toolIndex
	hash := uint(key) * 2_654_435_761
	slot := int((hash ^ hash>>16) & uint(len(projections)-1))
	for {
		projection := &projections[slot]
		if !projection.used {
			projection.key = key
			projection.used = true
			return projection
		}
		if projection.key == key {
			return projection
		}
		slot = (slot + 1) & (len(projections) - 1)
	}
}

func addAccumulatorTextBytes(total *int, text string) bool {
	if len(text) > maxChatCompletionAccumulatorTextBytes-*total {
		return false
	}
	*total += len(text)
	return true
}

func addAccumulatorMetadataBytes(total *int, text string) bool {
	if len(text) > maxChatCompletionAccumulatorMetadataBytes-*total {
		return false
	}
	*total += len(text)
	return true
}

// the state changed. This ensures that JustFinished events only fire once.
func (prev *chatCompletionResponseState) update(choice ChatCompletionChunkChoice) (justFinished chatCompletionResponseState) {
	delta := choice.Delta
	new := chatCompletionResponseState{choiceIndex: int(choice.Index)}
	switch {
	case len(delta.ToolCalls) > 0 && delta.Content == "":
		new.state = toolResponseState
		new.toolCallIndex = clampToZero(delta.ToolCalls[0].Index)
	case delta.JSON.Content.Valid():
		new.state = contentResponseState
	case delta.JSON.Refusal.Valid():
		new.state = refusalResponseState
	case len(delta.ToolCalls) > 0:
		new.state = toolResponseState
		new.toolCallIndex = clampToZero(delta.ToolCalls[0].Index)
	default:
		new.state = finishedResponseState
	}

	if *prev != new {
		justFinished = *prev
	}
	*prev = new

	return
}

// clampToZero handles providers like AWS Bedrock that return tool call index -1.
func clampToZero(index int64) int {
	if index < 0 {
		return 0
	}
	return int(index)
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

func detachTruncatedTail[T any](slice []T, previousLength int) []T {
	// A full slice expression can hide a retained tail while reporting len == cap.
	if slice == nil || (len(slice) == cap(slice) && len(slice) >= previousLength) {
		return slice
	}
	detached := make([]T, len(slice))
	copy(detached, slice)
	return detached
}

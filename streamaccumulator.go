package openai

import (
	"strings"
	"unsafe"

	"github.com/openai/openai-go/v3/packages/respjson"
	"github.com/openai/openai-go/v3/shared/constant"
)

const (
	maxChatCompletionAccumulatorChunks          = 100_000
	maxChatCompletionAccumulatorStructuralSlots = 1_024
	maxChatCompletionAccumulatorTextBytes       = 16 << 20
	maxChatCompletionAccumulatorMetadataBytes   = 16 << 20
	maxChatCompletionAccumulatorLogprobBytes    = 16 << 20
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
	choices []*chatCompletionChoiceStringState
}

type chatCompletionChoiceStringState struct {
	content   chatCompletionString
	refusal   chatCompletionString
	toolCalls []*chatCompletionToolCallStringState
}

type chatCompletionToolCallStringState struct {
	name      chatCompletionString
	arguments chatCompletionString
}

type chatCompletionString struct {
	builder   strings.Builder
	published string
}

type chatCompletionAccumulatorLogprobState struct {
	choices []chatCompletionChoiceLogprobState
	bytes   int
}

type chatCompletionChoiceLogprobState struct {
	content chatCompletionLogprobSliceState
	refusal chatCompletionLogprobSliceState
}

type chatCompletionLogprobSliceState struct {
	// A uintptr fingerprints published storage without keeping cleared backing alive.
	data     uintptr
	length   int
	capacity int
	bytes    int
}

type chatCompletionLogprobProjection struct {
	contentCount int
	contentBytes int
	refusalCount int
	refusalBytes int
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
// probability data. A rejected chunk does not modify the accumulator.
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
	if !chatCompletionStructuralSlotsWithinLimit(&acc.ChatCompletion, chunk) {
		return false
	}

	textBytes, ok := addChatCompletionTextBytes(0, &acc.ChatCompletion)
	if !ok {
		return false
	}
	if _, ok = addChatCompletionChunkTextBytes(textBytes, chunk); !ok {
		return false
	}
	metadataBytes, ok := addChatCompletionMetadataBytes(0, &acc.ChatCompletion)
	if !ok {
		return false
	}
	if _, ok = addChatCompletionChunkMetadataBytes(metadataBytes, &acc.ChatCompletion, chunk); !ok {
		return false
	}
	if !acc.reconcileLogprobState() || !acc.logprobState.chunkWithinLimit(chunk) {
		return false
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
	}
	return acc.choices[index]
}

func (acc *ChatCompletionAccumulator) reconcilePublicState() {
	completion := &acc.ChatCompletion
	stringState := &acc.stringState

	previousChoiceCount := len(stringState.choices)
	completion.Choices = detachTruncatedTail(completion.Choices, previousChoiceCount)
	for i := range completion.Choices {
		message := &completion.Choices[i].Message
		previousToolCallCount := 0
		if i < previousChoiceCount && stringState.choices[i] != nil {
			previousToolCallCount = len(stringState.choices[i].toolCalls)
		}
		message.ToolCalls = detachTruncatedTail(message.ToolCalls, previousToolCallCount)
	}

	choiceCount := min(len(stringState.choices), len(completion.Choices))
	for i := range choiceCount {
		choiceState := stringState.choices[i]
		if choiceState == nil {
			continue
		}

		message := &completion.Choices[i].Message
		choiceState.content.reconcile(&message.Content)
		choiceState.refusal.reconcile(&message.Refusal)

		toolCallCount := min(len(choiceState.toolCalls), len(message.ToolCalls))
		if toolCallCount < len(choiceState.toolCalls) {
			invalidateRemovedToolCallState(acc.legacyChoiceChatCompletionStates, i, toolCallCount)
			invalidateRemovedToolCallState(acc.choiceChatCompletionStates, i, toolCallCount)
		}
		for j := range toolCallCount {
			toolCallState := choiceState.toolCalls[j]
			if toolCallState == nil {
				continue
			}
			function := &message.ToolCalls[j].Function
			toolCallState.name.reconcile(&function.Name)
			toolCallState.arguments.reconcile(&function.Arguments)
		}
		clear(choiceState.toolCalls[toolCallCount:])
		choiceState.toolCalls = choiceState.toolCalls[:toolCallCount]
	}
	clear(stringState.choices[choiceCount:])
	stringState.choices = stringState.choices[:choiceCount]
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

func chatCompletionStructuralSlotsWithinLimit(completion *ChatCompletion, chunk *ChatCompletionChunk) bool {
	choiceCount := len(completion.Choices)
	if choiceCount > maxChatCompletionAccumulatorStructuralSlots {
		return false
	}

	structuralSlots := choiceCount
	// This fixed-size projection keeps duplicate choice deltas exact without a
	// per-chunk map allocation. The aggregate slot limit bounds its stack cost.
	var toolCallCounts [maxChatCompletionAccumulatorStructuralSlots]int
	for i := range completion.Choices {
		toolCallCount := len(completion.Choices[i].Message.ToolCalls)
		if toolCallCount > maxChatCompletionAccumulatorStructuralSlots-structuralSlots {
			return false
		}
		structuralSlots += toolCallCount
		toolCallCounts[i] = toolCallCount
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

func addChatCompletionTextBytes(total int, completion *ChatCompletion) (int, bool) {
	for i := range completion.Choices {
		message := &completion.Choices[i].Message
		if !addAccumulatorTextBytes(&total, message.Content) ||
			!addAccumulatorTextBytes(&total, message.Refusal) {
			return 0, false
		}
		for j := range message.ToolCalls {
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

func addChatCompletionMetadataBytes(total int, completion *ChatCompletion) (int, bool) {
	if !addAccumulatorMetadataBytes(&total, completion.ID) ||
		!addAccumulatorMetadataBytes(&total, completion.Model) ||
		!addAccumulatorMetadataBytes(&total, completion.SystemFingerprint) ||
		!addAccumulatorMetadataBytes(&total, string(completion.Object)) ||
		!addAccumulatorMetadataBytes(&total, string(completion.ServiceTier)) {
		return 0, false
	}
	for i := range completion.Choices {
		choice := &completion.Choices[i]
		if !addAccumulatorMetadataBytes(&total, choice.FinishReason) ||
			!addAccumulatorMetadataBytes(&total, string(choice.Message.Role)) {
			return 0, false
		}
		for j := range choice.Message.ToolCalls {
			toolCall := &choice.Message.ToolCalls[j]
			if !addAccumulatorMetadataBytes(&total, toolCall.ID) ||
				!addAccumulatorMetadataBytes(&total, toolCall.Type) {
				return 0, false
			}
		}
	}
	return total, true
}

func addChatCompletionChunkMetadataBytes(total int, completion *ChatCompletion, chunk *ChatCompletionChunk) (int, bool) {
	if completion.ID == "" && !addAccumulatorMetadataBytes(&total, chunk.ID) {
		return 0, false
	}
	if !addAccumulatorMetadataBytes(&total, chunk.Model) ||
		!addAccumulatorMetadataBytes(&total, chunk.SystemFingerprint) ||
		!addAccumulatorMetadataBytes(&total, string(chunk.Object)) ||
		!addAccumulatorMetadataBytes(&total, string(chunk.ServiceTier)) {
		return 0, false
	}
	for i := range chunk.Choices {
		choice := &chunk.Choices[i]
		if !addAccumulatorMetadataBytes(&total, choice.FinishReason) ||
			!addAccumulatorMetadataBytes(&total, choice.Delta.Role) {
			return 0, false
		}
		for j := range choice.Delta.ToolCalls {
			toolCall := &choice.Delta.ToolCalls[j]
			if !addAccumulatorMetadataBytes(&total, toolCall.ID) ||
				!addAccumulatorMetadataBytes(&total, toolCall.Type) {
				return 0, false
			}
		}
	}
	return total, true
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

func (acc *ChatCompletionAccumulator) reconcileLogprobState() bool {
	state := &acc.logprobState
	choiceCount := len(acc.Choices)
	if choiceCount < len(state.choices) {
		for i := choiceCount; i < len(state.choices); i++ {
			state.bytes -= state.choices[i].content.bytes + state.choices[i].refusal.bytes
		}
		clear(state.choices[choiceCount:])
		state.choices = state.choices[:choiceCount]
	}
	for i := range acc.Choices {
		logprobs := &acc.Choices[i].Logprobs
		if i >= len(state.choices) {
			if cap(logprobs.Content) == 0 && cap(logprobs.Refusal) == 0 {
				continue
			}
			state.choices = expandToFit(state.choices, i)
		}
		choiceState := &state.choices[i]
		if !state.reconcileSlice(&choiceState.content, &logprobs.Content) ||
			!state.reconcileSlice(&choiceState.refusal, &logprobs.Refusal) {
			return false
		}
	}
	return true
}

func (state *chatCompletionAccumulatorLogprobState) reconcileSlice(current *chatCompletionLogprobSliceState, logprobs *[]ChatCompletionTokenLogprob) bool {
	header := chatCompletionLogprobHeader(*logprobs)
	if current.data == header.data && current.length == header.length && current.capacity == header.capacity {
		return true
	}
	if chatCompletionLogprobSliceRetainsBacking(*current, header) {
		*logprobs = detachChatCompletionLogprobs(*logprobs)
		header = chatCompletionLogprobHeader(*logprobs)
	}

	bytes := 0
	if !addAccumulatorLogprobStorage(&bytes, cap(*logprobs), int(unsafe.Sizeof(ChatCompletionTokenLogprob{}))) ||
		!addChatCompletionLogprobData(&bytes, (*logprobs)[:cap(*logprobs)]) {
		return false
	}
	projected := state.bytes - current.bytes
	if !addAccumulatorLogprobBytes(&projected, bytes) {
		return false
	}
	header.bytes = bytes
	*current = header
	state.bytes = projected
	return true
}

func chatCompletionLogprobSliceRetainsBacking(current chatCompletionLogprobSliceState, public chatCompletionLogprobSliceState) bool {
	if current.data == 0 || current.capacity == 0 || public.data == 0 {
		return false
	}
	allocationEnd := current.data + uintptr(current.capacity)*unsafe.Sizeof(ChatCompletionTokenLogprob{})
	withinCurrentAllocation := public.data >= current.data && public.data < allocationEnd
	return withinCurrentAllocation && (public.data != current.data || public.capacity < current.capacity)
}

func (state *chatCompletionAccumulatorLogprobState) chunkWithinLimit(chunk *ChatCompletionChunk) bool {
	for i := range chunk.Choices {
		logprobs := &chunk.Choices[i].Logprobs
		if len(logprobs.Content) > 0 || len(logprobs.Refusal) > 0 {
			return state.nonEmptyChunkWithinLimit(chunk)
		}
	}
	return true
}

func (state *chatCompletionAccumulatorLogprobState) nonEmptyChunkWithinLimit(chunk *ChatCompletionChunk) bool {
	var projections [maxChatCompletionAccumulatorStructuralSlots]chatCompletionLogprobProjection
	var touched [maxChatCompletionAccumulatorStructuralSlots]int
	touchedCount := 0
	for i := range chunk.Choices {
		choice := &chunk.Choices[i]
		projection := &projections[choice.Index]
		wasEmpty := projection.contentCount == 0 && projection.refusalCount == 0
		if !addChatCompletionLogprobProjection(&projection.contentCount, &projection.contentBytes, choice.Logprobs.Content) ||
			!addChatCompletionLogprobProjection(&projection.refusalCount, &projection.refusalBytes, choice.Logprobs.Refusal) {
			return false
		}
		if wasEmpty && (projection.contentCount > 0 || projection.refusalCount > 0) {
			touched[touchedCount] = int(choice.Index)
			touchedCount++
		}
	}

	total := state.bytes
	logprobSize := int(unsafe.Sizeof(ChatCompletionTokenLogprob{}))
	for _, i := range touched[:touchedCount] {
		projection := &projections[i]
		var current chatCompletionChoiceLogprobState
		if i < len(state.choices) {
			current = state.choices[i]
		}
		if !addProjectedLogprobSliceBytes(&total, current.content, projection.contentCount, projection.contentBytes, logprobSize) ||
			!addProjectedLogprobSliceBytes(&total, current.refusal, projection.refusalCount, projection.refusalBytes, logprobSize) {
			return false
		}
	}
	return true
}

func addChatCompletionLogprobProjection(count *int, bytes *int, logprobs []ChatCompletionTokenLogprob) bool {
	maxLogprobs := maxChatCompletionAccumulatorLogprobBytes / int(unsafe.Sizeof(ChatCompletionTokenLogprob{}))
	if len(logprobs) > maxLogprobs-*count {
		return false
	}
	*count += len(logprobs)
	return addChatCompletionLogprobData(bytes, logprobs)
}

func addProjectedLogprobSliceBytes(total *int, current chatCompletionLogprobSliceState, count int, dataBytes int, logprobSize int) bool {
	maxLogprobs := maxChatCompletionAccumulatorLogprobBytes / logprobSize
	if count > maxLogprobs-current.length {
		return false
	}
	capacity := projectedLogprobCapacity(current.capacity, current.length+count, maxLogprobs)
	if !addAccumulatorLogprobStorage(total, capacity-current.capacity, logprobSize) {
		return false
	}
	return addAccumulatorLogprobBytes(total, dataBytes)
}

func (state *chatCompletionAccumulatorLogprobState) acceptChunk(completion *ChatCompletion, chunk *ChatCompletionChunk) {
	hasLogprobs := false
	for i := range chunk.Choices {
		logprobs := &chunk.Choices[i].Logprobs
		if len(logprobs.Content) > 0 || len(logprobs.Refusal) > 0 {
			hasLogprobs = true
			break
		}
	}
	if !hasLogprobs {
		return
	}
	if len(completion.Choices) > len(state.choices) {
		state.choices = expandToFit(state.choices, len(completion.Choices)-1)
	}
	logprobSize := int(unsafe.Sizeof(ChatCompletionTokenLogprob{}))
	maxLogprobs := maxChatCompletionAccumulatorLogprobBytes / logprobSize
	for i := range chunk.Choices {
		choice := &chunk.Choices[i]
		choiceState := &state.choices[choice.Index]
		state.acceptSlice(&choiceState.content, choice.Logprobs.Content, logprobSize, maxLogprobs)
		state.acceptSlice(&choiceState.refusal, choice.Logprobs.Refusal, logprobSize, maxLogprobs)
	}
	for i := range chunk.Choices {
		choiceIndex := int(chunk.Choices[i].Index)
		logprobs := &completion.Choices[choiceIndex].Logprobs
		setChatCompletionLogprobHeader(&state.choices[choiceIndex].content, logprobs.Content)
		setChatCompletionLogprobHeader(&state.choices[choiceIndex].refusal, logprobs.Refusal)
	}
}

func (state *chatCompletionAccumulatorLogprobState) acceptSlice(current *chatCompletionLogprobSliceState, logprobs []ChatCompletionTokenLogprob, logprobSize int, maxLogprobs int) {
	if len(logprobs) == 0 {
		return
	}
	dataBytes := 0
	// preflightChunk ran this same bounded calculation before accumulation.
	_ = addChatCompletionLogprobData(&dataBytes, logprobs)
	capacity := projectedLogprobCapacity(current.capacity, current.length+len(logprobs), maxLogprobs)
	current.bytes += (capacity-current.capacity)*logprobSize + dataBytes
	state.bytes += (capacity-current.capacity)*logprobSize + dataBytes
	current.length += len(logprobs)
	current.capacity = capacity
}

func appendChatCompletionLogprobs(dst []ChatCompletionTokenLogprob, src []ChatCompletionTokenLogprob) []ChatCompletionTokenLogprob {
	if len(src) == 0 {
		return dst
	}
	required := len(dst) + len(src)
	if required > cap(dst) {
		maxLogprobs := maxChatCompletionAccumulatorLogprobBytes / int(unsafe.Sizeof(ChatCompletionTokenLogprob{}))
		capacity := projectedLogprobCapacity(cap(dst), required, maxLogprobs)
		grown := make([]ChatCompletionTokenLogprob, len(dst), capacity)
		copy(grown, dst)
		dst = grown
	}
	for i := range src {
		dst = append(dst, cloneChatCompletionTokenLogprob(src[i]))
	}
	return dst
}

func projectedLogprobCapacity(current int, required int, maximum int) int {
	if required <= current {
		return current
	}
	capacity := max(1, current)
	// Deterministic doubling keeps append amortized-linear while making the
	// retained logical capacity exactly predictable during preflight.
	for capacity < required {
		if capacity > maximum/2 {
			return required
		}
		capacity *= 2
	}
	return capacity
}

func chatCompletionLogprobHeader(logprobs []ChatCompletionTokenLogprob) chatCompletionLogprobSliceState {
	return chatCompletionLogprobSliceState{
		data:     uintptr(unsafe.Pointer(unsafe.SliceData(logprobs))),
		length:   len(logprobs),
		capacity: cap(logprobs),
	}
}

func setChatCompletionLogprobHeader(state *chatCompletionLogprobSliceState, logprobs []ChatCompletionTokenLogprob) {
	state.data = uintptr(unsafe.Pointer(unsafe.SliceData(logprobs)))
	state.length = len(logprobs)
	state.capacity = cap(logprobs)
}

func addChatCompletionLogprobData(total *int, logprobs []ChatCompletionTokenLogprob) bool {
	for i := range logprobs {
		logprob := &logprobs[i]
		if !addAccumulatorLogprobBytes(total, len(logprob.Token)) ||
			!addAccumulatorLogprobStorage(total, cap(logprob.Bytes), int(unsafe.Sizeof(int64(0)))) ||
			!addAccumulatorLogprobStorage(total, cap(logprob.TopLogprobs), int(unsafe.Sizeof(ChatCompletionTokenLogprobTopLogprob{}))) ||
			!addChatCompletionLogprobMetadata(total, logprob.RawJSON(), logprob.JSON.ExtraFields,
				logprob.JSON.Token, logprob.JSON.Bytes, logprob.JSON.Logprob, logprob.JSON.TopLogprobs) {
			return false
		}
		topLogprobs := logprob.TopLogprobs[:cap(logprob.TopLogprobs)]
		for j := range topLogprobs {
			topLogprob := &topLogprobs[j]
			if !addAccumulatorLogprobBytes(total, len(topLogprob.Token)) ||
				!addAccumulatorLogprobStorage(total, cap(topLogprob.Bytes), int(unsafe.Sizeof(int64(0)))) ||
				!addChatCompletionLogprobMetadata(total, topLogprob.RawJSON(), topLogprob.JSON.ExtraFields,
					topLogprob.JSON.Token, topLogprob.JSON.Bytes, topLogprob.JSON.Logprob) {
				return false
			}
		}
	}
	return true
}

func addChatCompletionLogprobMetadata(total *int, raw string, extraFields map[string]respjson.Field, fields ...respjson.Field) bool {
	if len(extraFields) > 0 && !addAccumulatorLogprobBytes(total, chatCompletionAccumulatorMapOverheadBytes) {
		return false
	}
	entrySize := int(unsafe.Sizeof(string(""))) + int(unsafe.Sizeof(respjson.Field{}))
	if !addAccumulatorLogprobStorage(total, len(extraFields), entrySize) {
		return false
	}
	if !addAccumulatorLogprobStorage(total, len(extraFields), entrySize) {
		return false
	}
	if !addAccumulatorLogprobBytes(total, len(raw)) {
		return false
	}
	for i := range fields {
		if !addAccumulatorLogprobBytes(total, len(fields[i].Raw())) {
			return false
		}
	}
	for name, field := range extraFields {
		if !addAccumulatorLogprobBytes(total, len(name)) ||
			!addAccumulatorLogprobBytes(total, len(field.Raw())) {
			return false
		}
	}
	return true
}

func addAccumulatorLogprobStorage(total *int, count int, size int) bool {
	if count > (maxChatCompletionAccumulatorLogprobBytes-*total)/size {
		return false
	}
	*total += count * size
	return true
}

func addAccumulatorLogprobBytes(total *int, count int) bool {
	if count > maxChatCompletionAccumulatorLogprobBytes-*total {
		return false
	}
	*total += count
	return true
}

// Updates the internal response state and returns the previous state if
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

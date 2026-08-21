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
	maxChatCompletionAccumulatorLogprobBytes    = 16 << 20
	chatCompletionAccumulatorMapOverheadBytes   = 512
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
	chunkCount                       int
	logprobBytes                     int
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
// combined choice and tool-call slots, 16 MiB of combined content, refusal, tool name,
// and tool argument text currently stored in the accumulator and the incoming chunk,
// and 16 MiB of aggregate retained log probability data. A rejected chunk does not
// modify the accumulator.
//
// The ChatCompletion field JSON does not get accumulated.
func (acc *ChatCompletionAccumulator) AddChunk(chunk ChatCompletionChunk) bool {
	logprobBytes, ok := acc.preflightChunk(&chunk)
	if !ok {
		return false
	}

	acc.justFinished = chatCompletionResponseState{}
	acc.justFinishedByChoice = acc.justFinishedByChoice[:0]
	acc.accumulateDelta(&chunk)
	acc.chunkCount++
	acc.logprobBytes = logprobBytes

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

func (acc *ChatCompletionAccumulator) preflightChunk(chunk *ChatCompletionChunk) (int, bool) {
	if acc.chunkCount >= maxChatCompletionAccumulatorChunks {
		return 0, false
	}
	if acc.ID != "" && acc.ID != chunk.ID {
		return 0, false
	}
	if !chatCompletionStructuralSlotsWithinLimit(&acc.ChatCompletion, chunk) {
		return 0, false
	}

	textBytes, ok := addChatCompletionTextBytes(0, &acc.ChatCompletion)
	if !ok {
		return 0, false
	}
	if _, ok = addChatCompletionChunkTextBytes(textBytes, chunk); !ok {
		return 0, false
	}
	logprobBytes, ok := addChatCompletionChunkLogprobBytes(acc.logprobBytes, chunk)
	if !ok {
		return 0, false
	}

	acc.reconcilePublicState()
	return logprobBytes, true
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
		cc.ID = chunk.ID
	}

	for _, delta := range chunk.Choices {
		cc.Choices = expandToFit(cc.Choices, int(delta.Index))
		choice := &cc.Choices[delta.Index]

		choice.Index = delta.Index
		choice.FinishReason = delta.FinishReason

		if delta.Delta.Role != "" {
			choice.Message.Role = constant.Assistant(delta.Delta.Role)
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
				tool.ID = deltaTool.ID
			}
			if deltaTool.Type != "" {
				tool.Type = deltaTool.Type
			}
			toolStrings := choiceStrings.toolCall(toolIndex)
			toolStrings.name.append(&tool.Function.Name, deltaTool.Function.Name)
			toolStrings.arguments.append(&tool.Function.Arguments, deltaTool.Function.Arguments)
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

func addAccumulatorTextBytes(total *int, text string) bool {
	if len(text) > maxChatCompletionAccumulatorTextBytes-*total {
		return false
	}
	*total += len(text)
	return true
}

func addChatCompletionChunkLogprobBytes(total int, chunk *ChatCompletionChunk) (int, bool) {
	chunkBytes := 0
	for i := range chunk.Choices {
		logprobs := &chunk.Choices[i].Logprobs
		if !addChatCompletionLogprobs(&chunkBytes, logprobs.Content) ||
			!addChatCompletionLogprobs(&chunkBytes, logprobs.Refusal) {
			return 0, false
		}
	}

	if chunkBytes == 0 {
		return total, true
	}
	if !addAccumulatorLogprobBytes(&total, chunkBytes) {
		return 0, false
	}
	return total, true
}

func addChatCompletionLogprobs(total *int, logprobs []ChatCompletionTokenLogprob) bool {
	// Appending can reserve more backing storage than the current length. Twice the
	// element storage is a conservative bound for the accumulator-owned slice.
	logprobSize := int(unsafe.Sizeof(ChatCompletionTokenLogprob{}))
	if !addAccumulatorLogprobStorage(total, len(logprobs), logprobSize) {
		return false
	}
	if !addAccumulatorLogprobStorage(total, len(logprobs), logprobSize) {
		return false
	}
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
	if raw != "" {
		return addAccumulatorLogprobBytes(total, len(raw))
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

package openai

import "slices"

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
	if slice == nil || len(slice) >= previousLength {
		return slice
	}
	detached := make([]T, len(slice))
	copy(detached, slice)
	return detached
}

type chatCompletionToolMetadataProjection struct {
	key    chatCompletionToolProjectionKey
	fields uint8
	used   bool
}

type chatCompletionToolProjectionKey struct {
	choice int
	tool   int
}

type chatCompletionToolProjectionTable struct {
	inline   [chatCompletionAccumulatorInlineProjectionSlots * 2]chatCompletionToolMetadataProjection
	overflow map[chatCompletionToolProjectionKey]*chatCompletionToolMetadataProjection
	count    int
}

type chatCompletionTextAppendProjection struct {
	content [2]uint64
	refusal [2]uint64
	tools   chatCompletionToolProjectionTable
}

const (
	projectedToolID uint8 = 1 << iota
	projectedToolType
)

const (
	projectedToolName uint8 = 1 << iota
	projectedToolArguments
)

func (projection *chatCompletionTextAppendProjection) addChunk(chunk *ChatCompletionChunk) {
	for i := range chunk.Choices {
		choice := &chunk.Choices[i]
		choiceIndex := int(choice.Index)
		if choice.Delta.Content != "" {
			markChatCompletionTextAppend(&projection.content, choiceIndex)
		}
		if choice.Delta.Refusal != "" {
			markChatCompletionTextAppend(&projection.refusal, choiceIndex)
		}
		for j := range choice.Delta.ToolCalls {
			tool := &choice.Delta.ToolCalls[j]
			toolIndex := preflightedToolCallIndex(tool.Index)
			fields := uint8(0)
			if tool.Function.Name != "" {
				fields |= projectedToolName
			}
			if tool.Function.Arguments != "" {
				fields |= projectedToolArguments
			}
			if fields != 0 {
				projection.lookupTool(choiceIndex, toolIndex).fields |= fields
			}
		}
	}
}

func ensureChatCompletionTextAppendProjection(
	projection *chatCompletionTextAppendProjection,
	chunk *ChatCompletionChunk,
) *chatCompletionTextAppendProjection {
	if projection == nil {
		projection = &chatCompletionTextAppendProjection{}
		projection.addChunk(chunk)
	}
	return projection
}

func (projection *chatCompletionTextAppendProjection) choiceContent(index int) bool {
	return chatCompletionTextAppendMarked(&projection.content, index)
}

func (projection *chatCompletionTextAppendProjection) choiceRefusal(index int) bool {
	return chatCompletionTextAppendMarked(&projection.refusal, index)
}

func (projection *chatCompletionTextAppendProjection) toolName(choiceIndex int, toolIndex int) bool {
	return projection.toolFields(choiceIndex, toolIndex)&projectedToolName != 0
}

func (projection *chatCompletionTextAppendProjection) toolArguments(choiceIndex int, toolIndex int) bool {
	return projection.toolFields(choiceIndex, toolIndex)&projectedToolArguments != 0
}

func markChatCompletionTextAppend(marked *[2]uint64, index int) {
	marked[index/64] |= uint64(1) << (index % 64)
}

func chatCompletionTextAppendMarked(marked *[2]uint64, index int) bool {
	return marked[index/64]&(uint64(1)<<(index%64)) != 0
}

func (projection *chatCompletionTextAppendProjection) toolFields(choiceIndex int, toolIndex int) uint8 {
	if tool := findChatCompletionToolMetadataProjection(&projection.tools, choiceIndex, toolIndex); tool != nil {
		return tool.fields
	}
	return 0
}

func (projection *chatCompletionTextAppendProjection) lookupTool(choiceIndex int, toolIndex int) *chatCompletionToolMetadataProjection {
	return lookupChatCompletionToolMetadataProjection(&projection.tools, choiceIndex, toolIndex)
}

func (acc *ChatCompletionAccumulator) chatCompletionMetadataWithinLimit(chunk *ChatCompletionChunk, work *int) bool {
	completion := &acc.ChatCompletion
	if !acc.addChatCompletionChunkMetadataWork(work, chunk) {
		return false
	}
	total, ok := acc.addChatCompletionMetadataBytes(0, work, chunk)
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
		return projected >= 0 && projected <= int64(maxChatCompletionAccumulatorInt)
	}

	var finishReasonSeen [2]uint64
	var roleSeen [2]uint64
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
	return projected >= 0 && projected <= int64(maxChatCompletionAccumulatorInt)
}

func (acc *ChatCompletionAccumulator) addChatCompletionChunkMetadataWork(work *int, chunk *ChatCompletionChunk) bool {
	completion := &acc.ChatCompletion
	idPasses := 1
	if acc.ID != "" || acc.stringState.id != "" {
		idPasses = 2
	}
	if !addAccumulatorStringCopyWork(work, chunk.ID, idPasses) {
		return false
	}
	if !addAccumulatorStringAssignmentAfterPublicReconciliation(work, chunk.Model, completion.Model, acc.stringState.model) ||
		!addAccumulatorStringAssignmentAfterPublicReconciliation(work, chunk.SystemFingerprint, completion.SystemFingerprint, acc.stringState.systemFingerprint) ||
		!addAccumulatorStringAssignmentAfterPublicReconciliation(work, chunk.ServiceTier, completion.ServiceTier, acc.stringState.serviceTier) {
		return false
	}
	if chunk.Object == chunk.Object.Default() &&
		!addAccumulatorStringAssignmentAfterPublicReconciliation(work, chunk.Object.Default(), completion.Object, acc.stringState.object) {
		return false
	}
	var finishReasonSeen [2]uint64
	var roleSeen [2]uint64
	for i := range chunk.Choices {
		delta := &chunk.Choices[i]
		choiceIndex := int(delta.Index)
		var choiceState *chatCompletionChoiceStringState
		if choiceIndex < len(acc.stringState.choices) {
			choiceState = acc.stringState.choices[choiceIndex]
		}
		finishReason, role := "", ""
		hasPublicChoice := choiceState != nil && choiceIndex < len(completion.Choices)
		if hasPublicChoice {
			finishReason = choiceState.finishReason
			role = choiceState.role
			choice := &completion.Choices[choiceIndex]
			if !addAccumulatorMetadataAssignmentAfterPublicReconciliation(
				work,
				delta.FinishReason,
				choice.FinishReason,
				finishReason,
				&finishReasonSeen,
				choiceIndex,
			) {
				return false
			}
			if delta.Delta.Role != "" && !addAccumulatorMetadataAssignmentAfterPublicReconciliation(
				work,
				delta.Delta.Role,
				choice.Message.Role,
				role,
				&roleSeen,
				choiceIndex,
			) {
				return false
			}
		}
		if !hasPublicChoice &&
			!addAccumulatorMetadataAssignmentWork(work, delta.FinishReason, finishReason, &finishReasonSeen, choiceIndex) {
			return false
		}
		if !hasPublicChoice && delta.Delta.Role != "" &&
			!addAccumulatorMetadataAssignmentWork(work, delta.Delta.Role, role, &roleSeen, choiceIndex) {
			return false
		}
	}
	return acc.addChatCompletionToolMetadataWork(work, chunk)
}

func addAccumulatorMetadataAssignmentWork[T ~string](
	work *int,
	value T,
	published string,
	seen *[2]uint64,
	index int,
) bool {
	word := index / 64
	mask := uint64(1) << (index % 64)
	if seen[word]&mask != 0 {
		return addAccumulatorStringCopyWork(work, value, 3)
	}
	seen[word] |= mask
	return addAccumulatorStringAssignmentWork(work, value, published)
}

func addAccumulatorMetadataAssignmentAfterPublicReconciliation[T ~string, U ~string](
	work *int,
	value T,
	current U,
	published string,
	seen *[2]uint64,
	index int,
) bool {
	word := index / 64
	mask := uint64(1) << (index % 64)
	if seen[word]&mask != 0 {
		return addAccumulatorStringCopyWork(work, value, 3)
	}
	seen[word] |= mask
	return addAccumulatorStringAssignmentAfterPublicReconciliation(work, value, current, published)
}

func (acc *ChatCompletionAccumulator) addChatCompletionToolMetadataWork(work *int, chunk *ChatCompletionChunk) bool {
	completion := &acc.ChatCompletion
	hasMetadata := false
	for i := range chunk.Choices {
		for j := range chunk.Choices[i].Delta.ToolCalls {
			tool := &chunk.Choices[i].Delta.ToolCalls[j]
			hasMetadata = hasMetadata || tool.ID != "" || tool.Type != ""
		}
	}
	if !hasMetadata {
		return true
	}

	var seen chatCompletionToolProjectionTable
	for i := range chunk.Choices {
		choice := &chunk.Choices[i]
		choiceIndex := int(choice.Index)
		var choiceState *chatCompletionChoiceStringState
		if choiceIndex < len(acc.stringState.choices) {
			choiceState = acc.stringState.choices[choiceIndex]
		}
		for j := range choice.Delta.ToolCalls {
			tool := &choice.Delta.ToolCalls[j]
			if tool.ID == "" && tool.Type == "" {
				continue
			}
			toolIndex := preflightedToolCallIndex(tool.Index)
			projection := lookupChatCompletionToolMetadataProjection(&seen, choiceIndex, toolIndex)
			id, typeName := "", ""
			hasPublicTool := choiceState != nil &&
				toolIndex < len(choiceState.toolCalls) && choiceState.toolCallState(toolIndex) != nil &&
				choiceIndex < len(completion.Choices) &&
				toolIndex < len(completion.Choices[choiceIndex].Message.ToolCalls)
			if hasPublicTool {
				toolState := choiceState.toolCallState(toolIndex)
				id = toolState.id
				typeName = toolState.typeName
				current := &completion.Choices[choiceIndex].Message.ToolCalls[toolIndex]
				if tool.ID != "" && projection.fields&projectedToolID == 0 &&
					!addAccumulatorStringAssignmentAfterPublicReconciliation(work, tool.ID, current.ID, id) {
					return false
				}
				if tool.Type != "" && projection.fields&projectedToolType == 0 &&
					!addAccumulatorStringAssignmentAfterPublicReconciliation(work, tool.Type, current.Type, typeName) {
					return false
				}
			}
			if tool.ID != "" {
				if projection.fields&projectedToolID == 0 {
					if !hasPublicTool && !addAccumulatorStringAssignmentWork(work, tool.ID, id) {
						return false
					}
				} else if !addAccumulatorStringCopyWork(work, tool.ID, 3) {
					return false
				}
				projection.fields |= projectedToolID
			}
			if tool.Type != "" {
				if projection.fields&projectedToolType == 0 {
					if !hasPublicTool && !addAccumulatorStringAssignmentWork(work, tool.Type, typeName) {
						return false
					}
				} else if !addAccumulatorStringCopyWork(work, tool.Type, 3) {
					return false
				}
				projection.fields |= projectedToolType
			}
		}
	}
	return true
}

func addAccumulatorStringAssignmentAfterPublicReconciliation[T ~string, U ~string](
	work *int,
	value T,
	current U,
	published string,
) bool {
	currentValue := string(current)
	if accumulatorStringUsesPublishedBacking(currentValue, published) {
		return addAccumulatorStringAssignmentWork(work, value, published)
	}

	text := string(value)
	if len(text) != len(currentValue) {
		return addAccumulatorStringCopyWork(work, text, 1)
	}
	// Reconciliation publishes either the prior equal value or a clone of the
	// changed public value. In both cases its backing differs from currentValue.
	if accumulatorStringUsesPublishedBacking(text, currentValue) {
		return addAccumulatorStringCopyWork(work, text, 1)
	}
	if !addAccumulatorStringCopyWork(work, text, 1) {
		return false
	}
	if text == currentValue {
		return addAccumulatorStringCopyWork(work, text, 1)
	}
	return addAccumulatorStringCopyWork(work, text, 2)
}

func addAccumulatorStringAssignmentWork[T ~string](work *int, value T, published string) bool {
	text := string(value)
	if accumulatorStringUsesPublishedBacking(text, published) {
		return true
	}
	if len(text) != len(published) {
		return addAccumulatorStringCopyWork(work, text, 1)
	}
	// Preflight performs this comparison once, and assignment repeats it after
	// every other check has succeeded. A changed value is then copied once.
	if !addAccumulatorStringCopyWork(work, text, 1) {
		return false
	}
	if text == published {
		return addAccumulatorStringCopyWork(work, text, 1)
	}
	return addAccumulatorStringCopyWork(work, text, 2)
}

func markChatCompletionMetadataProjected(seen *[2]uint64, index int) bool {
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
	var projections chatCompletionToolProjectionTable
	deltaBytes := int64(0)
	for i := len(chunk.Choices) - 1; i >= 0; i-- {
		choice := &chunk.Choices[i]
		choiceIndex := int(choice.Index)
		for j := len(choice.Delta.ToolCalls) - 1; j >= 0; j-- {
			toolCall := &choice.Delta.ToolCalls[j]
			if toolCall.ID == "" && toolCall.Type == "" {
				continue
			}
			toolIndex := preflightedToolCallIndex(toolCall.Index)
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
	projections *chatCompletionToolProjectionTable,
	choiceIndex int,
	toolIndex int,
) *chatCompletionToolMetadataProjection {
	key := chatCompletionToolProjectionKey{choice: choiceIndex, tool: toolIndex}
	if projections.overflow != nil {
		if projection := projections.overflow[key]; projection != nil {
			return projection
		}
		projection := &chatCompletionToolMetadataProjection{key: key, used: true}
		projections.overflow[key] = projection
		return projection
	}
	if projections.count >= len(projections.inline)/2 {
		projections.overflow = make(map[chatCompletionToolProjectionKey]*chatCompletionToolMetadataProjection, 2*projections.count)
		for i := range projections.inline {
			projection := &projections.inline[i]
			if projection.used {
				projections.overflow[projection.key] = projection
			}
		}
		return lookupChatCompletionToolMetadataProjection(projections, choiceIndex, toolIndex)
	}
	hash := uint(toolIndex)*2_654_435_761 ^ uint(choiceIndex)*4_046_345_921
	slot := int((hash ^ hash>>16) & uint(len(projections.inline)-1))
	for {
		projection := &projections.inline[slot]
		if !projection.used {
			projection.key = key
			projection.used = true
			projections.count++
			return projection
		}
		if projection.key == key {
			return projection
		}
		slot = (slot + 1) & (len(projections.inline) - 1)
	}
}

func findChatCompletionToolMetadataProjection(
	projections *chatCompletionToolProjectionTable,
	choiceIndex int,
	toolIndex int,
) *chatCompletionToolMetadataProjection {
	key := chatCompletionToolProjectionKey{choice: choiceIndex, tool: toolIndex}
	if projections.overflow != nil {
		return projections.overflow[key]
	}
	hash := uint(toolIndex)*2_654_435_761 ^ uint(choiceIndex)*4_046_345_921
	slot := int((hash ^ hash>>16) & uint(len(projections.inline)-1))
	for {
		projection := &projections.inline[slot]
		if !projection.used {
			return nil
		}
		if projection.key == key {
			return projection
		}
		slot = (slot + 1) & (len(projections.inline) - 1)
	}
}

func addAccumulatorTextBytes(total *int, text string) bool {
	return addAccumulatorTextWork(total, len(text))
}

func addAccumulatorTextWork(total *int, bytes int) bool {
	if bytes > maxChatCompletionAccumulatorInt-*total {
		return false
	}
	*total += bytes
	return true
}

func addAccumulatorMetadataBytes(total *int, text string) bool {
	if len(text) > maxChatCompletionAccumulatorInt-*total {
		return false
	}
	*total += len(text)
	return true
}

package openai

type chatCompletionToolMetadataProjection struct {
	key    int
	fields uint8
	used   bool
}

const (
	projectedToolID uint8 = 1 << iota
	projectedToolType
)

func (acc *ChatCompletionAccumulator) chatCompletionMetadataWithinLimit(chunk *ChatCompletionChunk, work *int) bool {
	completion := &acc.ChatCompletion
	if !acc.addChatCompletionChunkMetadataWork(work, chunk) {
		return false
	}
	total, ok := acc.addChatCompletionMetadataBytes(0, work)
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

func (acc *ChatCompletionAccumulator) addChatCompletionChunkMetadataWork(work *int, chunk *ChatCompletionChunk) bool {
	idPasses := 1
	if acc.ID != "" || acc.stringState.id != "" {
		idPasses = 2
	}
	if !addAccumulatorStringCopyWork(work, chunk.ID, idPasses) {
		return false
	}
	if !addAccumulatorStringAssignmentWork(work, chunk.Model, acc.stringState.model) ||
		!addAccumulatorStringAssignmentWork(work, chunk.SystemFingerprint, acc.stringState.systemFingerprint) ||
		!addAccumulatorStringAssignmentWork(work, chunk.ServiceTier, acc.stringState.serviceTier) {
		return false
	}
	if chunk.Object == chunk.Object.Default() &&
		!addAccumulatorStringAssignmentWork(work, chunk.Object.Default(), acc.stringState.object) {
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
		if choiceState != nil {
			finishReason = choiceState.finishReason
			role = choiceState.role
		}
		if !addAccumulatorMetadataAssignmentWork(work, delta.FinishReason, finishReason, &finishReasonSeen, choiceIndex) {
			return false
		}
		if delta.Delta.Role != "" &&
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

func (acc *ChatCompletionAccumulator) addChatCompletionToolMetadataWork(work *int, chunk *ChatCompletionChunk) bool {
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

	var seen [maxChatCompletionAccumulatorStructuralSlots * 2]chatCompletionToolMetadataProjection
	for i := range chunk.Choices {
		choice := &chunk.Choices[i]
		choiceIndex := int(choice.Index)
		var choiceState *chatCompletionChoiceStringState
		if choiceIndex < len(acc.stringState.choices) {
			choiceState = acc.stringState.choices[choiceIndex]
		}
		for j := range choice.Delta.ToolCalls {
			tool := &choice.Delta.ToolCalls[j]
			toolIndex := preflightedToolCallIndex(tool.Index)
			projection := lookupChatCompletionToolMetadataProjection(&seen, choiceIndex, toolIndex)
			id, typeName := "", ""
			if choiceState != nil && toolIndex < len(choiceState.toolCalls) && choiceState.toolCalls[toolIndex] != nil {
				id = choiceState.toolCalls[toolIndex].id
				typeName = choiceState.toolCalls[toolIndex].typeName
			}
			if tool.ID != "" {
				if projection.fields&projectedToolID == 0 {
					if !addAccumulatorStringAssignmentWork(work, tool.ID, id) {
						return false
					}
				} else if !addAccumulatorStringCopyWork(work, tool.ID, 3) {
					return false
				}
				projection.fields |= projectedToolID
			}
			if tool.Type != "" {
				if projection.fields&projectedToolType == 0 {
					if !addAccumulatorStringAssignmentWork(work, tool.Type, typeName) {
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

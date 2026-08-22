package openai

type chatCompletionToolMetadataProjection struct {
	key    int
	fields uint8
	used   bool
}

type chatCompletionTextAppendProjection struct {
	content [2]uint64
	refusal [2]uint64
	tools   [maxChatCompletionAccumulatorStructuralSlots * 2]chatCompletionToolTextAppendProjection
}

type chatCompletionToolTextAppendProjection struct {
	keyPlusOne uint32
	fields     uint8
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
	keyPlusOne := uint32(choiceIndex*maxChatCompletionAccumulatorStructuralSlots + toolIndex + 1)
	slot := chatCompletionTextAppendSlot(keyPlusOne, len(projection.tools))
	for {
		tool := &projection.tools[slot]
		if tool.keyPlusOne == 0 {
			return 0
		}
		if tool.keyPlusOne == keyPlusOne {
			return tool.fields
		}
		slot = (slot + 1) & (len(projection.tools) - 1)
	}
}

func (projection *chatCompletionTextAppendProjection) lookupTool(choiceIndex int, toolIndex int) *chatCompletionToolTextAppendProjection {
	keyPlusOne := uint32(choiceIndex*maxChatCompletionAccumulatorStructuralSlots + toolIndex + 1)
	slot := chatCompletionTextAppendSlot(keyPlusOne, len(projection.tools))
	for {
		tool := &projection.tools[slot]
		if tool.keyPlusOne == 0 {
			tool.keyPlusOne = keyPlusOne
			return tool
		}
		if tool.keyPlusOne == keyPlusOne {
			return tool
		}
		slot = (slot + 1) & (len(projection.tools) - 1)
	}
}

func chatCompletionTextAppendSlot(keyPlusOne uint32, slots int) int {
	hash := uint(keyPlusOne) * 2_654_435_761
	return int((hash ^ hash>>16) & uint(slots-1))
}

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
			hasPublicTool := choiceState != nil &&
				toolIndex < len(choiceState.toolCalls) && choiceState.toolCalls[toolIndex] != nil &&
				choiceIndex < len(completion.Choices) &&
				toolIndex < len(completion.Choices[choiceIndex].Message.ToolCalls)
			if hasPublicTool {
				id = choiceState.toolCalls[toolIndex].id
				typeName = choiceState.toolCalls[toolIndex].typeName
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

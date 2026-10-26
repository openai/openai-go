package openai

import "weak"

const chatCompletionToolOwnerRadixBits = 4

type chatCompletionToolCallOwner struct {
	children [1 << chatCompletionToolOwnerRadixBits]*chatCompletionToolCallOwner
	state    *chatCompletionToolCallStringState
}

func chatCompletionToolCallOwnerAt(
	owner *chatCompletionToolCallOwner,
	depth uint8,
	index int,
) *chatCompletionToolCallStringState {
	if depth == 0 || uint(index)>>(depth*chatCompletionToolOwnerRadixBits) != 0 {
		return nil
	}
	for depth > 0 && owner != nil {
		depth--
		digit := (uint(index) >> (depth * chatCompletionToolOwnerRadixBits)) &
			((1 << chatCompletionToolOwnerRadixBits) - 1)
		owner = owner.children[digit]
	}
	if owner == nil {
		return nil
	}
	return owner.state
}

func (choice *chatCompletionChoiceStringState) setToolCallOwner(index int, state *chatCompletionToolCallStringState) {
	if choice.toolOwnerDepth == 0 {
		choice.toolOwnerDepth = 1
	}
	for uint(index)>>(choice.toolOwnerDepth*chatCompletionToolOwnerRadixBits) != 0 {
		choice.toolOwners = &chatCompletionToolCallOwner{children: [1 << chatCompletionToolOwnerRadixBits]*chatCompletionToolCallOwner{
			0: choice.toolOwners,
		}}
		choice.toolOwnerDepth++
	}
	choice.toolOwners = setChatCompletionToolCallOwner(choice.toolOwners, choice.toolOwnerDepth, index, state)
}

func setChatCompletionToolCallOwner(
	owner *chatCompletionToolCallOwner,
	depth uint8,
	index int,
	state *chatCompletionToolCallStringState,
) *chatCompletionToolCallOwner {
	cloned := &chatCompletionToolCallOwner{}
	if owner != nil {
		*cloned = *owner
	}
	if depth == 0 {
		cloned.state = state
		return cloned
	}
	depth--
	digit := (uint(index) >> (depth * chatCompletionToolOwnerRadixBits)) &
		((1 << chatCompletionToolOwnerRadixBits) - 1)
	cloned.children[digit] = setChatCompletionToolCallOwner(cloned.children[digit], depth, index, state)
	return cloned
}

func (acc *ChatCompletionAccumulator) privateStateNeedsDetach() bool {
	return acc.privateStateInitialized && acc.privateStateOwner.Value() != acc
}

func (acc *ChatCompletionAccumulator) detachPrivateStateForCopy() {
	if !acc.privateStateNeedsDetach() {
		return
	}
	acc.stringState = acc.stringState.cloneForAccumulatorCopy()
	acc.logprobState.choices = cloneAccumulatorSlice(acc.logprobState.choices)
	acc.logprobState.markSharedSpareBacking()
}

func (acc *ChatCompletionAccumulator) claimPrivateStateOwnership() {
	if acc.privateStateInitialized && !acc.privateStateNeedsDetach() {
		return
	}
	acc.privateStateOwner = weak.Make(acc)
	acc.privateStateInitialized = true
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
		choiceCopy.toolOwners, choiceCopy.toolOwnerDepth = nil, 0
		for j := range choiceCopy.toolCalls {
			toolCall := choice.toolCallState(j)
			if toolCall == nil {
				continue
			}
			toolCallCopy := *toolCall
			toolCallCopy.name.shared = len(toolCallCopy.name.buffer) > 0
			toolCallCopy.arguments.shared = len(toolCallCopy.arguments.buffer) > 0
			choiceCopy.setToolCallOwner(j, &toolCallCopy)
			choiceCopy.toolCalls[j] = weak.Make(&toolCallCopy)
		}
		state.choices[i] = &choiceCopy
	}
	return state
}

func (acc *ChatCompletionAccumulator) privateStateCopyWork() int {
	if !acc.privateStateNeedsDetach() {
		return 0
	}
	// Charge both this projection scan and the accepted chunk's clone pass.
	work := 2*len(acc.stringState.choices) + len(acc.stringState.activeChoices) + len(acc.logprobState.choices)
	for _, choice := range acc.stringState.choices {
		if choice == nil {
			continue
		}
		work++
		work += 2*len(choice.toolCalls) + len(choice.activeToolCalls)
		for _, toolCall := range choice.toolCalls {
			if toolCall.Value() != nil {
				work++
			}
		}
	}
	return work
}

func (acc *ChatCompletionAccumulator) addCopiedTextBufferWork(work *int, chunk *ChatCompletionChunk) bool {
	copied := acc.privateStateNeedsDetach()
	var chargedContent, chargedRefusal [2]uint64
	var chargedTools *chatCompletionTextAppendProjection
	for i := range chunk.Choices {
		delta := &chunk.Choices[i]
		choiceIndex := int(delta.Index)
		if choiceIndex >= len(acc.Choices) || choiceIndex >= len(acc.stringState.choices) {
			continue
		}
		choiceState := acc.stringState.choices[choiceIndex]
		if choiceState == nil {
			continue
		}
		message := &acc.Choices[choiceIndex].Message
		if !addCopiedChoiceTextBufferWork(work, &chargedContent, choiceIndex, &choiceState.content,
			message.Content, delta.Delta.Content, copied) ||
			!addCopiedChoiceTextBufferWork(work, &chargedRefusal, choiceIndex, &choiceState.refusal,
				message.Refusal, delta.Delta.Refusal, copied) {
			return false
		}
		for j := range delta.Delta.ToolCalls {
			tool := &delta.Delta.ToolCalls[j]
			toolIndex := preflightedToolCallIndex(tool.Index)
			if toolIndex >= len(message.ToolCalls) || toolIndex >= len(choiceState.toolCalls) {
				continue
			}
			toolState := choiceState.toolCallState(toolIndex)
			if toolState == nil {
				continue
			}
			function := &message.ToolCalls[toolIndex].Function
			if !addCopiedToolTextBufferWork(work, &chargedTools, choiceIndex, toolIndex, projectedToolName,
				&toolState.name, function.Name, tool.Function.Name, copied) ||
				!addCopiedToolTextBufferWork(work, &chargedTools, choiceIndex, toolIndex, projectedToolArguments,
					&toolState.arguments, function.Arguments, tool.Function.Arguments, copied) {
				return false
			}
		}
	}
	return true
}

func addCopiedChoiceTextBufferWork(
	work *int,
	charged *[2]uint64,
	choiceIndex int,
	state *chatCompletionString,
	current string,
	fragment string,
	copied bool,
) bool {
	if !accumulatorBufferWillBeCopied(state, current, fragment, copied) ||
		chatCompletionTextAppendMarked(charged, choiceIndex) {
		return true
	}
	markChatCompletionTextAppend(charged, choiceIndex)
	return addAccumulatorTextWork(work, len(state.buffer))
}

func addCopiedToolTextBufferWork(
	work *int,
	charged **chatCompletionTextAppendProjection,
	choiceIndex int,
	toolIndex int,
	field uint8,
	state *chatCompletionString,
	current string,
	fragment string,
	copied bool,
) bool {
	if !accumulatorBufferWillBeCopied(state, current, fragment, copied) {
		return true
	}
	if *charged == nil {
		*charged = &chatCompletionTextAppendProjection{}
	}
	tool := (*charged).lookupTool(choiceIndex, toolIndex)
	if tool.fields&field != 0 {
		return true
	}
	tool.fields |= field
	return addAccumulatorTextWork(work, len(state.buffer))
}

func accumulatorBufferWillBeCopied(state *chatCompletionString, current string, fragment string, copied bool) bool {
	return fragment != "" && len(state.buffer) > 0 && (copied || state.shared) &&
		accumulatorStringUsesPublishedBacking(current, state.published)
}

package openai

// visitReconciledTools visits tools touched by the incoming delta plus enough
// rotating existing tools to match newly populated positions. Empty chunks
// still visit one tool, while dense growth cannot outrun cleared-buffer cleanup.
func (choice *chatCompletionChoiceStringState) visitReconciledTools(
	choiceIndex int,
	public []ChatCompletionMessageToolCallUnion,
	chunk *ChatCompletionChunk,
	visit func(int) bool,
) bool {
	if len(choice.activeToolCalls) == 0 {
		return true
	}

	cursor := choice.reconcileCursor % len(choice.activeToolCalls)
	swept := choice.activeToolCalls[cursor]
	for offset := range choice.reconciliationSweepCount(choiceIndex, chunk) {
		index := choice.activeToolCalls[(cursor+offset)%len(choice.activeToolCalls)]
		if index < len(public) && index < len(choice.toolCalls) && choice.toolCallState(index) != nil && !visit(index) {
			return false
		}
	}

	previous := -1
	for i := range chunk.Choices {
		delta := &chunk.Choices[i]
		if int(delta.Index) != choiceIndex {
			continue
		}
		for j := range delta.Delta.ToolCalls {
			index := preflightedToolCallIndex(delta.Delta.ToolCalls[j].Index)
			if index == swept || index == previous || index >= len(public) || index >= len(choice.toolCalls) ||
				choice.toolCallState(index) == nil {
				continue
			}
			if !visit(index) {
				return false
			}
			previous = index
		}
	}
	return true
}

func (choice *chatCompletionChoiceStringState) reconciliationSweepCount(
	choiceIndex int,
	chunk *ChatCompletionChunk,
) int {
	if len(choice.activeToolCalls) == 0 {
		return 0
	}
	count := 1
	for i := range chunk.Choices {
		delta := &chunk.Choices[i]
		if int(delta.Index) != choiceIndex {
			continue
		}
		for j := range delta.Delta.ToolCalls {
			index := preflightedToolCallIndex(delta.Delta.ToolCalls[j].Index)
			if index >= len(choice.toolCalls) || choice.toolCallState(index) == nil {
				count++
			}
		}
	}
	return min(count, len(choice.activeToolCalls))
}

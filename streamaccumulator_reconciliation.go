package openai

// visitReconciledTools visits tools touched by the incoming delta plus one
// rotating existing tool. The sweep eventually releases backing cleared
// through public fields without rescanning every dense tool on every chunk.
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
	if swept < len(public) && swept < len(choice.toolCalls) && choice.toolCalls[swept] != nil && !visit(swept) {
		return false
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
				choice.toolCalls[index] == nil {
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

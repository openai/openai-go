package openai

func chatCompletionMessageToAssistantParam(r ChatCompletionMessage) ChatCompletionAssistantMessageParam {
	var p ChatCompletionAssistantMessageParam

	// It is important to not rely on the JSON metadata property
	// here, it may be unset if the receiver was generated via a
	// [ChatCompletionAccumulator].
	//
	// Explicit null is intentionally elided from the response.
	if r.Content != "" {
		p.Content.OfString = String(r.Content)
	}
	if r.Refusal != "" {
		p.Refusal = String(r.Refusal)
	}

	p.Audio.ID = r.Audio.ID
	p.Role = r.Role
	p.FunctionCall.Arguments = r.FunctionCall.Arguments
	p.FunctionCall.Name = r.FunctionCall.Name

	if len(r.ToolCalls) > 0 {
		for _, v := range r.ToolCalls {
			u := ChatCompletionMessageToolCallUnionParam{}
			switch v.AsAny().(type) {
			case ChatCompletionMessageFunctionToolCall:
				u.OfFunction = &ChatCompletionMessageFunctionToolCallParam{
					ID: v.ID,
					Function: ChatCompletionMessageFunctionToolCallFunctionParam{
						Arguments: v.Function.Arguments,
						Name:      v.Function.Name,
					},
				}
			case ChatCompletionMessageCustomToolCall:
				u.OfCustom = &ChatCompletionMessageCustomToolCallParam{
					ID: v.ID,
					Custom: ChatCompletionMessageCustomToolCallCustomParam{
						Input: v.Custom.Input,
						Name:  v.Custom.Name,
					},
				}
			}

			p.ToolCalls = append(p.ToolCalls, u)
		}
	}
	return p
}

package openai

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/openai/openai-go/v3/packages/param"
)

func agentToolArguments(value any) (map[string]any, error) {
	var data []byte
	var err error
	if text, ok := value.(string); ok {
		data = []byte(text)
	} else {
		data, err = json.Marshal(value)
	}
	if err != nil {
		return nil, err
	}
	var args map[string]any
	// Round-trip objects too so nested handler mutations cannot change the event.
	if err := json.Unmarshal(data, &args); err != nil {
		return nil, err
	}
	if args == nil {
		return nil, errors.New("function arguments must be a JSON object")
	}
	return args, nil
}

func agentToolResult(ctx context.Context, call agentPendingCall) AgentSessionInputParamAgentSessionInputToolResult {
	result := AgentSessionInputParamAgentSessionInputToolResult{TurnID: call.turnID, CallID: call.callID}
	err := call.argumentErr
	if err == nil {
		var output any
		output, err = call.handler(ctx, call.arguments)
		if err == nil {
			result.Output, err = agentToolOutput(output)
		}
	}
	if err != nil {
		result.Output = AgentFunctionCallOutputParamUnion{}
		result.Error = param.NewOpt("Tool handler failed.")
	} else {
		result.Success = true
	}
	return result
}

func agentToolOutput(output any) (AgentFunctionCallOutputParamUnion, error) {
	var result AgentFunctionCallOutputParamUnion
	switch value := output.(type) {
	case nil:
		result = param.NullStruct[AgentFunctionCallOutputParamUnion]()
	case string:
		result.OfString = param.NewOpt(value)
	case map[string]any:
		data, err := json.Marshal(value)
		if err != nil {
			return result, err
		}
		result.OfString = param.NewOpt(string(data))
	case []InputContentParamUnion:
		result.OfArrayOfInputContents = value
		if value == nil {
			result.OfArrayOfInputContents = []InputContentParamUnion{}
		}
	default:
		return result, errors.New("unsupported tool output")
	}
	// Validate before posting so serialization failures produce a failed result.
	_, err := json.Marshal(result)
	return result, err
}

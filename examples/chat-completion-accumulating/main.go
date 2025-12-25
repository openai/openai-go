package main

import (
	"context"

	"github.com/Nordlys-Labs/openai-go/v3"
)

func main() {
	client := openai.NewClient()
	ctx := context.Background()

	sysprompt := "Share only a brief description of the place in 50 words. Then immediately make some tool calls and announce them."

	question := "Tell me about Greece's largest city."

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(sysprompt),
		openai.UserMessage(question),
	}

	print("> ")
	println(question)
	println()

	params := openai.ChatCompletionNewParams{
		Messages: messages,
		Seed:     openai.Int(0),
		Model:    openai.ChatModelGPT4o,
		Tools:    tools,
	}

	stream := client.Chat.Completions.NewStreaming(ctx, params)
	acc := openai.ChatCompletionAccumulator{}

	for stream.Next() {
		chunk := stream.Current()

		acc.AddChunk(chunk)

		// When this fires, the current chunk value will not contain content data
		if _, ok := acc.JustFinishedContent(); ok {
			println()
			println("finish-event: Content stream finished")
		}

		if refusal, ok := acc.JustFinishedRefusal(); ok {
			println()
			println("finish-event: refusal stream finished:", refusal)
			println()
		}

		if tool, ok := acc.JustFinishedToolCall(); ok {
			println("finish-event: tool call stream finished:", tool.Index, tool.Name, tool.Arguments)
		}

		// It's best to use chunks after handling JustFinished events.
		// Here we print the delta of the content, if it exists.
		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			print(chunk.Choices[0].Delta.Content)
		}
	}

	if err := stream.Err(); err != nil {
		panic(err)
	}

	if acc.Usage.TotalTokens > 0 {
		println("Total Tokens:", acc.Usage.TotalTokens)
	}
}

var tools = []openai.ChatCompletionToolUnionParam{
	openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
		Name:        "get_live_weather",
		Description: openai.String("Get weather at the given location"),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]any{
				"location": map[string]string{
					"type": "string",
				},
			},
			"required": []string{"location"},
		},
	}),
	openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
		Name:        "get_population",
		Description: openai.String("Get population of a given town"),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]any{
				"town": map[string]string{
					"type": "string",
				},
				"nation": map[string]string{
					"type": "string",
				},
				"rounding": map[string]string{
					"type":        "integer",
					"description": "Nearest base 10 to round to, e.g. 1000 or 1000000",
				},
			},
			"required": []string{"town", "nation"},
		},
	}),
}

// Mock function to simulate weather data retrieval
func getWeather(location string) string {
	// In a real implementation, this function would call a weather API
	return "Sunny, 25°C"
}

// Mock function to simulate population data retrieval
func getPopulation(town, nation string, rounding int) string {
	// In a real implementation, this function would call a population API
	return "Athens, Greece: 664,046"
}

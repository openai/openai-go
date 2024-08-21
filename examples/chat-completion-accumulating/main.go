package main

import (
	"context"

	"github.com/openai/openai-go"
)

// Mock function to simulate weather data retrieval
func getWeather(location string) string {
	// In a real implementation, this function would call a weather API
	return "Sunny, 25°C"
}

func main() {
	client := openai.NewClient()
	ctx := context.Background()

	question := "Begin a very brief introduction of Greece, then incorporate the local weather of a few towns"

	print("> ")
	println(question)
	println()

	params := openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(question),
		}),
		Seed:  openai.Int(0),
		Model: openai.F(openai.ChatModelGPT4o),
		Tools: openai.F([]openai.ChatCompletionToolParam{
			{
				Type: openai.F(openai.ChatCompletionToolTypeFunction),
				Function: openai.F(openai.FunctionDefinitionParam{
					Name:        openai.String("get_live_weather"),
					Description: openai.String("Get weather at the given location"),
					Parameters: openai.F(openai.FunctionParameters{
						"type": "object",
						"properties": map[string]interface{}{
							"location": map[string]string{
								"type": "string",
							},
						},
						"required": []string{"location"},
					}),
				}),
			},
		}),
	}

	stream := client.Chat.Completions.NewStreaming(ctx, params)

	acc := openai.ChatCompletionAccumulator{}

	for stream.Next() {
		chunk := stream.Current()
		acc.AddChunk(chunk)

		// When this fires, the current chunk value will not contain content data
		if content, ok := acc.JustFinishedContent(); ok {
			println("Content stream finished:", content)
			println()
		}

		if tool, ok := acc.JustFinishedToolCall(); ok {
			println("Tool call stream finished:", tool.Index, tool.Name, tool.Arguments)
			println()
		}

		if refusal, ok := acc.JustFinishedRefusal(); ok {
			println("Refusal stream finished:", refusal)
			println()
		}

		// It's best to use chunks after handling JustFinished events
		if len(chunk.Choices) > 0 {
			println(chunk.Choices[0].Delta.JSON.RawJSON())
		}
	}

	if err := stream.Err(); err != nil {
		panic(err)
	}

	// After the stream is finished, acc can be used like a ChatCompletion
	_ = acc.Choices[0].Message.Content

	println("Total Tokens:", acc.Usage.TotalTokens)
	println("Finish Reason:", acc.Choices[0].FinishReason)
}

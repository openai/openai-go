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

	questions := []string{"Describe Greece in 50 words.", "Grab the live weather of the largest city in Greece"}

	for _, question := range questions {
		print("> ")
		println(question)
		println()

		params := openai.ChatCompletionNewParams{
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.ChatCompletionMessageParamOfUser(question),
			},
			Seed:  openai.Int(0),
			Model: openai.ChatModelGPT4o,
			Tools: []openai.ChatCompletionToolParam{
				{
					Function: openai.FunctionDefinitionParam{
						Name:        "get_live_weather",
						Description: openai.String("Get weather at the given location"),
						Parameters: openai.FunctionParameters{
							"type": "object",
							"properties": map[string]interface{}{
								"location": map[string]string{
									"type": "string",
								},
							},
							"required": []string{"location"},
						},
					},
				},
			},
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
			if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
				print(chunk.Choices[0].Delta.Content)
			}
		}

		if err := stream.Err(); err != nil {
			panic(err)
		}

		if len(acc.Choices) == 0 {
			println("Total Tokens:", acc.Usage.TotalTokens)
			continue
		}

		if acc.Choices[0].Message.Content != "" {
			// After the stream is finished, acc can be used like a ChatCompletion
			println("Total text length:", len(acc.Choices[0].Message.Content))
		}

		if acc.Choices[0].Message.Refusal != "" {
			println("Refusal:", acc.Choices[0].Message.Refusal)
		}

		for _, toolCall := range acc.Choices[0].Message.ToolCalls {
			println("Tool call:", toolCall.Function.Name, toolCall.Function.Arguments)
		}

		if acc.Choices[0].FinishReason != "" {
			println("Finish Reason:", acc.Choices[0].FinishReason)
		}
	}

}

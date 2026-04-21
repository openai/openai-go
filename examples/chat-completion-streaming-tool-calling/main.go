package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/openai/openai-go/v3"
)

func main() {
	client := openai.NewClient()
	ctx := context.Background()

	question := "What's the weather in New York City, and what is its population?"

	fmt.Println("> " + question)
	fmt.Println()

	params := openai.ChatCompletionNewParams{
		Model: openai.ChatModelGPT4o,
		Seed:  openai.Int(0),
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("Answer briefly, but call tools first whenever you need external data."),
			openai.UserMessage(question),
		},
		Tools: []openai.ChatCompletionToolUnionParam{
			openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
				Name:        "get_weather",
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
				Description: openai.String("Get population for a city and country"),
				Parameters: openai.FunctionParameters{
					"type": "object",
					"properties": map[string]any{
						"city": map[string]string{
							"type": "string",
						},
						"country": map[string]string{
							"type": "string",
						},
					},
					"required": []string{"city", "country"},
				},
			}),
		},
	}

	stream := client.Chat.Completions.NewStreaming(ctx, params)
	acc := openai.ChatCompletionAccumulator{}

	for stream.Next() {
		chunk := stream.Current()
		acc.AddChunk(chunk)

		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			fmt.Print(chunk.Choices[0].Delta.Content)
		}

		if tool, ok := acc.JustFinishedToolCall(); ok {
			fmt.Printf("\nTool call ready: %s(%s)\n", tool.Name, tool.Arguments)
		}
	}

	if err := stream.Err(); err != nil {
		panic(err)
	}

	if len(acc.Choices) == 0 || len(acc.Choices[0].Message.ToolCalls) == 0 {
		fmt.Println("\nNo tool calls were requested.")
		return
	}

	params.Messages = append(params.Messages, acc.Choices[0].Message.ToParam())

	for _, toolCall := range acc.Choices[0].Message.ToolCalls {
		var args map[string]any
		if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err != nil {
			panic(err)
		}

		var output string
		switch toolCall.Function.Name {
		case "get_weather":
			output = getWeather(args["location"].(string))
		case "get_population":
			output = getPopulation(args["city"].(string), args["country"].(string))
		default:
			output = "unsupported tool"
		}

		fmt.Printf("Tool result: %s -> %s\n", toolCall.Function.Name, output)
		params.Messages = append(params.Messages, openai.ToolMessage(output, toolCall.ID))
	}

	fmt.Println("\nFinal answer:")

	finalStream := client.Chat.Completions.NewStreaming(ctx, params)
	for finalStream.Next() {
		chunk := finalStream.Current()
		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			fmt.Print(chunk.Choices[0].Delta.Content)
		}
	}

	if err := finalStream.Err(); err != nil {
		panic(err)
	}

	fmt.Println()
}

func getWeather(location string) string {
	return fmt.Sprintf("%s: sunny, 25C", location)
}

func getPopulation(city, country string) string {
	return fmt.Sprintf("%s, %s: population 8.3 million", city, country)
}

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

	question := "What is the weather in New York City? Write a short paragraph about it."

	print("> ")
	println(question)

	params := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
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
		},
		ParallelToolCalls: openai.Bool(false), // JustFinishedToolCall is only safe with parallel calls disabled
		Seed:              openai.Int(0),
		Model:             openai.ChatModelGPT4o,
	}

	stream := client.Chat.Completions.NewStreaming(ctx, params)
	acc := openai.ChatCompletionAccumulator{}

	fmt.Println("\nStreaming first response...")
	for stream.Next() {
		chunk := stream.Current()
		acc.AddChunk(chunk)

		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			print(chunk.Choices[0].Delta.Content)
		}

		if tool, ok := acc.JustFinishedToolCall(); ok {
			fmt.Printf("\nTool call detected: %s with arguments %s\n", tool.Name, tool.Arguments)
		}
	}
	println()

	if err := stream.Err(); err != nil {
		panic(err)
	}

	if len(acc.Choices) == 0 || len(acc.Choices[0].Message.ToolCalls) == 0 {
		fmt.Printf("No function call")
		return
	}

	message := acc.Choices[0].Message
	toolCalls := message.ToolCalls
	params.Messages = append(params.Messages, message.ToParam())
	for _, toolCall := range toolCalls {
		if toolCall.Function.Name == "get_weather" {
			var args map[string]any
			if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err != nil {
				panic(err)
			}

			location, ok := args["location"].(string)
			if !ok {
				fmt.Printf("Missing or invalid 'location' argument\n")
				return
			}
			weatherData := getWeather(location)
			fmt.Printf("Weather in %s: %s\n", location, weatherData)

			params.Messages = append(params.Messages, openai.ToolMessage(weatherData, toolCall.ID))
		}
	}

	// Disable tools for the second round so the model returns a final answer
	// instead of making additional tool calls we don't handle.
	params.Tools = []openai.ChatCompletionToolUnionParam{}

	responseStream := client.Chat.Completions.NewStreaming(ctx, params)

	fmt.Println("\nStreaming second response...")
	for responseStream.Next() {
		evt := responseStream.Current()
		_ = acc.AddChunk(evt) // best-effort; streaming accumulation for display only
		if len(evt.Choices) > 0 {
			print(evt.Choices[0].Delta.Content)
		}
	}
	println()

	if err := responseStream.Err(); err != nil {
		panic(err)
	}
}

func getWeather(location string) string {
	// In a real implementation, this function would call a weather API.
	return "Sunny, 25°C"
}
